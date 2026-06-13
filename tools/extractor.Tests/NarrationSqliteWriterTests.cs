using Extractor;
using Microsoft.Data.Sqlite;
using Mutagen.Bethesda.Plugins;
using Xunit;

namespace Extractor.Tests;

// SQLite writer の検証。叙述文 1 種（BOOK:DESC）を db/migrations の schema へ冪等に書くこと。
public class NarrationSqliteWriterTests
{
    private static readonly HashSet<string> BookDesc = ["BOOK:DESC"];

    private static string RepoRoot()
    {
        var dir = new DirectoryInfo(AppContext.BaseDirectory);
        while (dir != null && !File.Exists(Path.Combine(dir.FullName, "go.mod")))
            dir = dir.Parent;
        Assert.NotNull(dir);
        return dir!.FullName;
    }

    private static string Schema() =>
        File.ReadAllText(Path.Combine(RepoRoot(), "db", "migrations", "0001_init.sql"));

    private static ExtractionResult WithBook(string edid, uint id, string body)
    {
        var result = new ExtractionResult { TargetPlugin = "Test.esm" };
        result.Books.Add(new BookEntry
        {
            Id = new FormKey(ModKey.FromFileName("Test.esm"), id),
            EditorId = edid,
            Kind = "BOOK",
            Title = "Title",
            Body = body,
            Author = "Author",
        });
        return result;
    }

    [Fact]
    public void BOOK_DESC_を_narration_へ書く()
    {
        var dbPath = Path.Combine(Path.GetTempPath(), $"nw-{Guid.NewGuid():N}.sqlite3");
        try
        {
            var written = NarrationSqliteWriter.Write(dbPath, Schema(), WithBook("TestBook", 0x800, "Ancient Nord prose"), BookDesc);
            Assert.Equal(1, written);

            using var conn = new SqliteConnection($"Data Source={dbPath}");
            conn.Open();
            using var cmd = conn.CreateCommand();
            cmd.CommandText = "SELECT plugin, form_id, edid, rec, field, ordinal, source, status FROM narration";
            using var reader = cmd.ExecuteReader();
            Assert.True(reader.Read());
            Assert.Equal("Test.esm", reader.GetString(0));
            Assert.Equal("0x000800", reader.GetString(1));
            Assert.Equal("TestBook", reader.GetString(2));
            Assert.Equal("BOOK", reader.GetString(3));
            Assert.Equal("DESC", reader.GetString(4));
            Assert.Equal(0, reader.GetInt32(5));
            Assert.Equal("Ancient Nord prose", reader.GetString(6));
            Assert.Equal(0, reader.GetInt32(7));
            Assert.False(reader.Read());
        }
        finally
        {
            File.Delete(dbPath);
        }
    }

    [Fact]
    public void 同じ行を二度書いても重複しない()
    {
        var dbPath = Path.Combine(Path.GetTempPath(), $"nw-{Guid.NewGuid():N}.sqlite3");
        try
        {
            NarrationSqliteWriter.Write(dbPath, Schema(), WithBook("TestBook", 0x800, "prose"), BookDesc);
            NarrationSqliteWriter.Write(dbPath, Schema(), WithBook("TestBook", 0x800, "prose"), BookDesc);

            using var conn = new SqliteConnection($"Data Source={dbPath}");
            conn.Open();
            using var cmd = conn.CreateCommand();
            cmd.CommandText = "SELECT COUNT(*) FROM narration";
            Assert.Equal(1L, (long)cmd.ExecuteScalar()!);
        }
        finally
        {
            File.Delete(dbPath);
        }
    }
}
