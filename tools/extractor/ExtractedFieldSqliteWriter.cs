using Microsoft.Data.Sqlite;

namespace Extractor;

// 抽出結果の翻訳対象文字列を中心 DB（SQLite）の extracted_field へ素朴に書き出す。
// 箱（叙述文/固有名/定型句/台詞）や directive の判定は持たず、原文・日本語既訳（英日対がある field のみ）と
// (plugin, form_id, edid, rec, field, ordinal) だけを書く。
// 箱の振り分けは Go 取込段が record_type_master を引いて行う（C#↔Go 契約の責務分離）。
// stub（master と同一内容）は TranslationCounts.Enumerate が除くため、ここには現れない。
public static class ExtractedFieldSqliteWriter
{
    // dbPath の SQLite に schema を ensure し、翻訳対象文字列を extracted_field へ UPSERT する。書いた行数を返す。
    // 同一 (form_id, rec, field) に複数値が出るフィールド（INFO:NAM1 の応答、QUST:CNAM のログ、MESG:ITXT のボタン、
    // FACT:MNAM/FNAM の階級）には出現順で ordinal を採番し、一意キーの衝突を避ける。
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
            INSERT OR IGNORE INTO extracted_field (plugin, form_id, edid, rec, field, ordinal, source, dest)
            VALUES ($plugin, $form_id, $edid, $rec, $field, $ordinal, $source, $dest)
            """;
        var plugin = cmd.Parameters.Add("$plugin", SqliteType.Text);
        var formId = cmd.Parameters.Add("$form_id", SqliteType.Text);
        var edid = cmd.Parameters.Add("$edid", SqliteType.Text);
        var rec = cmd.Parameters.Add("$rec", SqliteType.Text);
        var field = cmd.Parameters.Add("$field", SqliteType.Text);
        var ordinal = cmd.Parameters.Add("$ordinal", SqliteType.Integer);
        var source = cmd.Parameters.Add("$source", SqliteType.Text);
        var dest = cmd.Parameters.Add("$dest", SqliteType.Text);
        plugin.Value = result.TargetPlugin; // plugin は result 全体で一定。

        // (form_id, rec, field) ごとの出現回数。同一キーの 2 値目以降に 1, 2, ... を採番する。
        var ordinals = new Dictionary<(string, string, string), int>();

        var written = 0;
        foreach (var s in TranslationCounts.Enumerate(result))
        {
            var parts = s.RecField.Split(':', 2);
            var rf = ($"0x{s.Id.ID:X6}", parts[0], parts[1]);
            var n = ordinals.GetValueOrDefault(rf);
            ordinals[rf] = n + 1;

            formId.Value = rf.Item1;
            edid.Value = s.EditorId;
            rec.Value = parts[0];
            field.Value = parts[1];
            ordinal.Value = n;
            source.Value = s.Text;
            // 英日 Strings が両方解決できた field だけ日本語本文を書く（片側欠けは空 = 対なし）。
            dest.Value = result.JapanesePairs.GetValueOrDefault((s.Id, s.RecField, s.Text), "");
            written += cmd.ExecuteNonQuery();
        }
        tx.Commit();
        return written;
    }
}
