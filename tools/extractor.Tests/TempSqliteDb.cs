using Microsoft.Data.Sqlite;

namespace Extractor.Tests;

// テストが writer へ渡す一時 SQLite ファイル。using で作って捨てる。
// 一時ファイルの作成と破棄を 1 箇所に集め、後片付けの順序を各テストへ散らさない。
public sealed class TempSqliteDb : IDisposable
{
    public string Path { get; }

    public string ConnectionString => $"Data Source={Path}";

    // prefix はどのテストが作ったファイルかを一時フォルダ上で見分けるために付ける。
    public TempSqliteDb(string prefix)
    {
        Path = System.IO.Path.Combine(System.IO.Path.GetTempPath(), $"{prefix}-{Guid.NewGuid():N}.sqlite3");
    }

    // 書かれた内容を読むための接続。呼び手が using で閉じる。
    public SqliteConnection OpenConnection()
    {
        var conn = new SqliteConnection(ConnectionString);
        conn.Open();
        return conn;
    }

    public void Dispose()
    {
        // Microsoft.Data.Sqlite は接続を閉じてもプールがファイルを開いたまま保持する。
        // 解放せずに消すと IOException（別プロセスが使用中）になるため、先にプールを解放する。
        SqliteConnection.ClearPool(new SqliteConnection(ConnectionString));
        if (File.Exists(Path)) File.Delete(Path);
    }
}
