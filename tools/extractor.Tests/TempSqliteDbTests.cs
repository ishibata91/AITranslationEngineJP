using Extractor;
using Mutagen.Bethesda.Plugins;
using Xunit;

namespace Extractor.Tests;

// 一時 SQLite ファイルの後片付けの検証。
// writer が書いたあとの一時ファイルを破棄で消せることを確かめる。Microsoft.Data.Sqlite は接続を閉じても
// プールがファイルを開いたまま保持するため、解放せずに消そうとすると IOException になる。
public class TempSqliteDbTests
{
    private static string RepoRoot()
    {
        var dir = new DirectoryInfo(AppContext.BaseDirectory);
        while (dir != null && !File.Exists(Path.Combine(dir.FullName, "go.mod")))
            dir = dir.Parent;
        Assert.NotNull(dir);
        return dir!.FullName;
    }

    private static string MigrationsDir() => Path.Combine(RepoRoot(), "db", "migrations");

    private static ExtractionResult WithBook()
    {
        var result = new ExtractionResult { TargetPlugin = "Test.esm" };
        result.Books.Add(new BookEntry
        {
            Id = new FormKey(ModKey.FromFileName("Test.esm"), 0x800),
            EditorId = "TestBook",
            Kind = "BOOK",
            Title = "Title",
            Body = "prose",
            Author = "Author",
        });
        return result;
    }

    [Fact]
    public void writer_が書いた一時ファイルを破棄で消せる()
    {
        // Arrange
        var db = new TempSqliteDb("tmp");

        // Act: writer がプロダクトコード側の接続で書き込む。書き込み後に破棄する。
        ExtractedFieldSqliteWriter.Write(db.Path, MigrationsDir(), WithBook());
        db.Dispose();

        // Assert
        Assert.False(File.Exists(db.Path), "破棄した一時ファイルが残っている");
    }

    [Fact]
    public void テスト側で開いた接続のあとも破棄で消せる()
    {
        // Arrange
        var db = new TempSqliteDb("tmp");
        ExtractedFieldSqliteWriter.Write(db.Path, MigrationsDir(), WithBook());

        // Act: 書かれた内容を読む接続を開いて閉じ、そのあとで破棄する。
        using (var conn = db.OpenConnection())
        {
            using var cmd = conn.CreateCommand();
            cmd.CommandText = "SELECT COUNT(*) FROM extracted_field";
            Assert.Equal(3L, (long)cmd.ExecuteScalar()!);
        }
        db.Dispose();

        // Assert
        Assert.False(File.Exists(db.Path), "破棄した一時ファイルが残っている");
    }
}
