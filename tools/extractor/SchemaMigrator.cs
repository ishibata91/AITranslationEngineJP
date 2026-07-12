using Microsoft.Data.Sqlite;

namespace Extractor;

// 中心 DB のスキーマを db/migrations/*.sql から適用する。Go の db.Apply（PRAGMA user_version で適用済みを飛ばす）と同じ規約。
// C# 抽出器は毎抽出で呼ぶが、user_version 未満の適用済み migration は再実行しない。
// これにより 0009 の proper_noun 作り直し（DROP→CREATE）などの破壊的 DDL を毎抽出で再実行して
// 永続データ（対象 plugin 単位で残す翻訳成果）を消すことを避ける。Go 起動時に既に適用済みなら C# 側は何もしない。
public static class SchemaMigrator
{
    // migrationsDir の *.sql をファイル名昇順に、user_version 未満の分だけ順に適用し、user_version を本数へ進める。
    // Go 側と同じ本数（ファイル数）を user_version の到達点にし、両者で version の解釈をそろえる。
    public static void Ensure(SqliteConnection conn, string migrationsDir)
    {
        var files = Directory.GetFiles(migrationsDir, "*.sql")
            .OrderBy(p => p, StringComparer.Ordinal)
            .ToArray();

        var current = ReadUserVersion(conn);
        if (current == files.Length)
            return; // 適用済み（多くは Go 起動時の db.Apply で到達済み）。
        // DB がアプリの想定より新しい場合は、未知の schema を壊さないよう適用せずエラーにする（Go の db.Apply と同じ）。
        if (current > files.Length)
            throw new InvalidOperationException(
                $"DB schema version {current} がアプリの想定 {files.Length} より新しい");

        for (var i = current; i < files.Length; i++)
        {
            using var cmd = conn.CreateCommand();
            cmd.CommandText = File.ReadAllText(files[i]);
            cmd.ExecuteNonQuery();
        }
        SetUserVersion(conn, files.Length);
    }

    private static int ReadUserVersion(SqliteConnection conn)
    {
        using var cmd = conn.CreateCommand();
        cmd.CommandText = "PRAGMA user_version";
        return Convert.ToInt32(cmd.ExecuteScalar() ?? 0);
    }

    private static void SetUserVersion(SqliteConnection conn, int version)
    {
        using var cmd = conn.CreateCommand();
        // PRAGMA はパラメータ化できないため定数を埋め込む（version は migration 本数で外部入力ではない）。
        cmd.CommandText = $"PRAGMA user_version = {version}";
        cmd.ExecuteNonQuery();
    }
}
