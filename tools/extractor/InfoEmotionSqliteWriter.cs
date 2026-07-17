using Microsoft.Data.Sqlite;
using Mutagen.Bethesda.Plugins;

namespace Extractor;

// INFO 応答の感情型（TRDT。制作者が付ける「その台詞を喋る時の表情演出」）を extracted_info_emotion へ書く。
// ordinal は非空応答の出現順で、ExtractedFieldSqliteWriter の INFO:NAM1 と一致させ、Go 取込段が line_emotion へ結ぶ。
// Neutral は書かない（翻訳プロンプトへ加算しないため）。stub INFO（master 同一）は extracted_field に現れないため揃えて除く。
// 書くのは事実（感情型）だけで、口調などの解釈は持たない（解釈は Go engine が与える）。
public static class InfoEmotionSqliteWriter
{
    // dbPath の SQLite に schema を ensure し、INFO 応答の感情型を extracted_info_emotion へ書く。書いた行数を返す。
    public static int Write(string dbPath, string migrationsDir, ExtractionResult result)
    {
        using var conn = new SqliteConnection($"Data Source={dbPath}");
        conn.Open();

        SchemaMigrator.Ensure(conn, migrationsDir);

        using var tx = conn.BeginTransaction();
        using var cmd = conn.CreateCommand();
        cmd.Transaction = tx;
        cmd.CommandText =
            """
            INSERT OR IGNORE INTO extracted_info_emotion (info_plugin, info_form_id, ordinal, emotion_type)
            VALUES ($plugin, $form_id, $ordinal, $emotion)
            """;
        var plugin = cmd.Parameters.Add("$plugin", SqliteType.Text);
        var formId = cmd.Parameters.Add("$form_id", SqliteType.Text);
        var ordinal = cmd.Parameters.Add("$ordinal", SqliteType.Integer);
        var emotion = cmd.Parameters.Add("$emotion", SqliteType.Text);
        plugin.Value = result.TargetPlugin; // plugin は result 全体で一定。

        var written = 0;
        foreach (var dlg in result.Dialogues)
        {
            foreach (var info in dlg.Infos)
            {
                if (info.IdenticalToMaster) continue; // stub INFO は extracted_field に現れないため揃える。
                var n = 0;
                foreach (var resp in info.Responses)
                {
                    if (resp.Text.Trim().Length == 0) continue; // extracted_field の空除外（AddRaw）と揃える。
                    if (resp.Emotion.Length > 0 && resp.Emotion != "Neutral")
                    {
                        formId.Value = HexFormId(info.Id);
                        ordinal.Value = n;
                        emotion.Value = resp.Emotion;
                        written += cmd.ExecuteNonQuery();
                    }
                    n++; // Neutral でも ordinal は進め、extracted_field の応答順（非空応答の出現順）と一致させる。
                }
            }
        }
        tx.Commit();
        return written;
    }

    private static string HexFormId(FormKey key) => $"0x{key.ID:X6}";
}
