using Microsoft.Data.Sqlite;

namespace Extractor;

// 抽出結果の英日対（english と japanese の Strings を両方解決できた field）を中心 DB（SQLite）の
// reference_translation へ書き出す。既訳（人間が作った訳）の照合表で、Go 側の完全一致置換と
// 横断辞書（master_term）の派生が入力にする。
// 翻訳対象の原文を書く extracted_field とは器を分ける。extracted_field は利用者が選んだ 1 plugin の
// 翻訳対象だけを持ち、こちらは Data フォルダにある全 plugin の英日対を持つ。
// plugin 列を持たない。既訳は plugin をまたいで同一原文へ再利用するため、由来で絞らない。
public static class ReferenceTranslationSqliteWriter
{
    // dbPath の SQLite に schema を ensure し、result の英日対を reference_translation へ書く。追加した行数を返す。
    // 翻訳所有（stub 除外）で絞らない。stub の record が運ぶ既訳も同一原文への供給として使えるため。
    // UNIQUE (rec, field, source) と INSERT OR IGNORE で、同じ Data フォルダを 2 回走らせても増えない（冪等）。
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
            INSERT OR IGNORE INTO reference_translation (rec, field, source, dest)
            VALUES ($rec, $field, $source, $dest)
            """;
        var rec = cmd.Parameters.Add("$rec", SqliteType.Text);
        var field = cmd.Parameters.Add("$field", SqliteType.Text);
        var source = cmd.Parameters.Add("$source", SqliteType.Text);
        var dest = cmd.Parameters.Add("$dest", SqliteType.Text);

        var written = 0;
        foreach (var (key, japanese) in result.JapanesePairs)
        {
            // key.RecField は "NPC_:FULL" の形。照合キーの rec と field へ割る。
            var parts = key.RecField.Split(':', 2);
            if (parts.Length != 2) continue;
            rec.Value = parts[0];
            field.Value = parts[1];
            source.Value = key.English;
            dest.Value = japanese;
            written += cmd.ExecuteNonQuery();
        }
        tx.Commit();
        return written;
    }
}
