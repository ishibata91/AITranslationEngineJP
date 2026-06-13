using Microsoft.Data.Sqlite;

namespace Extractor;

// 抽出結果の指定 REC:FIELD（叙述文）を中心 DB（SQLite）の narration テーブルへ書く。
// どの REC:FIELD を叙述文として書くかは呼び出し側が recFields で渡す（writer は分類知識を持たない再利用可能な機構）。
// schema は db/migrations の SQL を渡して冪等 ensure する（C#↔Go 契約 1 本）。
public static class NarrationSqliteWriter
{
    // dbPath の SQLite に schema を ensure し、recFields に含まれる叙述文を narration へ UPSERT する。書いた行数を返す。
    public static int Write(string dbPath, string schemaSql, ExtractionResult result, ISet<string> recFields)
    {
        using var conn = new SqliteConnection($"Data Source={dbPath}");
        conn.Open();

        using (var ddl = conn.CreateCommand())
        {
            ddl.CommandText = schemaSql;
            ddl.ExecuteNonQuery();
        }

        using var tx = conn.BeginTransaction();
        using var cmd = conn.CreateCommand();
        cmd.Transaction = tx;
        cmd.CommandText =
            """
            INSERT OR IGNORE INTO narration (plugin, form_id, edid, rec, field, ordinal, source, status)
            VALUES ($plugin, $form_id, $edid, $rec, $field, 0, $source, 0)
            """;
        // パラメータを 1 度だけ登録し、行ごとに値だけ差し替えて prepared statement を再利用する。
        var plugin = cmd.Parameters.Add("$plugin", SqliteType.Text);
        var formId = cmd.Parameters.Add("$form_id", SqliteType.Text);
        var edid = cmd.Parameters.Add("$edid", SqliteType.Text);
        var rec = cmd.Parameters.Add("$rec", SqliteType.Text);
        var field = cmd.Parameters.Add("$field", SqliteType.Text);
        var source = cmd.Parameters.Add("$source", SqliteType.Text);
        plugin.Value = result.TargetPlugin; // plugin は result 全体で一定。

        var written = 0;
        foreach (var s in TranslationCounts.Enumerate(result))
        {
            if (!recFields.Contains(s.RecField)) continue;
            var parts = s.RecField.Split(':', 2);
            formId.Value = $"0x{s.Id.ID:X6}";
            edid.Value = s.EditorId;
            rec.Value = parts[0];
            field.Value = parts[1];
            source.Value = s.Text;
            written += cmd.ExecuteNonQuery();
        }
        tx.Commit();
        return written;
    }
}
