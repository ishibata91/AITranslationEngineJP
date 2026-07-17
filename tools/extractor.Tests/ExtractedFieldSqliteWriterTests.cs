using Extractor;
using Microsoft.Data.Sqlite;
using Mutagen.Bethesda.Plugins;
using Xunit;

namespace Extractor.Tests;

// extracted_field writer の検証。全 REC:FIELD の原文を素朴に書き、同一 (form_id, rec, field) の複数値へ ordinal を採番すること。
// schema は db/migrations 全体を ensure する（extracted_field は 0006）。
public class ExtractedFieldSqliteWriterTests
{
    private static string RepoRoot()
    {
        var dir = new DirectoryInfo(AppContext.BaseDirectory);
        while (dir != null && !File.Exists(Path.Combine(dir.FullName, "go.mod")))
            dir = dir.Parent;
        Assert.NotNull(dir);
        return dir!.FullName;
    }

    // migrations ディレクトリ（extracted_field は 0006）。Writer が SchemaMigrator で ensure する。
    private static string MigrationsDir() => Path.Combine(RepoRoot(), "db", "migrations");

    private static ExtractionResult WithBook(string edid, uint id, string title, string body, string author)
    {
        var result = new ExtractionResult { TargetPlugin = "Test.esm" };
        result.Books.Add(new BookEntry
        {
            Id = new FormKey(ModKey.FromFileName("Test.esm"), id),
            EditorId = edid,
            Kind = "BOOK",
            Title = title,
            Body = body,
            Author = author,
        });
        return result;
    }

    [Fact]
    public void 本の_FULL_DESC_CNAM_を_extracted_field_へ書く()
    {
        var dbPath = Path.Combine(Path.GetTempPath(), $"ef-{Guid.NewGuid():N}.sqlite3");
        try
        {
            var written = ExtractedFieldSqliteWriter.Write(dbPath, MigrationsDir(),
                WithBook("TestBook", 0x800, "Lusty Argonian Maid", "Ancient Nord prose", "Crassius Curio"));
            Assert.Equal(3, written); // FULL + DESC + CNAM

            using var conn = new SqliteConnection($"Data Source={dbPath}");
            conn.Open();
            using var cmd = conn.CreateCommand();
            cmd.CommandText = "SELECT rec, field, ordinal, source FROM extracted_field ORDER BY field";
            using var reader = cmd.ExecuteReader();

            var got = new List<(string rec, string field, int ordinal, string source)>();
            while (reader.Read())
                got.Add((reader.GetString(0), reader.GetString(1), reader.GetInt32(2), reader.GetString(3)));

            Assert.Contains(("BOOK", "FULL", 0, "Lusty Argonian Maid"), got);
            Assert.Contains(("BOOK", "DESC", 0, "Ancient Nord prose"), got);
            Assert.Contains(("BOOK", "CNAM", 0, "Crassius Curio"), got);
        }
        finally
        {
            File.Delete(dbPath);
        }
    }

    [Fact]
    public void 同一フィールドの複数値へ_ordinal_を採番する()
    {
        // INFO の応答 3 件は同じ (form_id, INFO, NAM1)。ordinal 0,1,2 で一意キーが衝突しないこと。
        var result = new ExtractionResult { TargetPlugin = "Test.esp" };
        var info = new InfoNode
        {
            Id = new FormKey(ModKey.FromFileName("Test.esp"), 0x900),
            EditorId = "InfoA",
            Prompt = "",
        };
        info.Responses.Add(new ResponseLine(1, "First.", "Neutral"));
        info.Responses.Add(new ResponseLine(2, "Second.", "Neutral"));
        info.Responses.Add(new ResponseLine(3, "Third.", "Neutral"));
        result.Dialogues.Add(new DialogueTopic
        {
            Id = new FormKey(ModKey.FromFileName("Test.esp"), 0x901),
            EditorId = "Topic",
            Kind = "DIAL",
            Name = "",
            Category = "",
            Subtype = "",
            Infos = { info },
        });

        var dbPath = Path.Combine(Path.GetTempPath(), $"ef-{Guid.NewGuid():N}.sqlite3");
        try
        {
            var written = ExtractedFieldSqliteWriter.Write(dbPath, MigrationsDir(), result);
            Assert.Equal(3, written);

            using var conn = new SqliteConnection($"Data Source={dbPath}");
            conn.Open();
            using var cmd = conn.CreateCommand();
            cmd.CommandText = "SELECT ordinal, source FROM extracted_field WHERE rec='INFO' AND field='NAM1' ORDER BY ordinal";
            using var reader = cmd.ExecuteReader();

            Assert.True(reader.Read());
            Assert.Equal(0, reader.GetInt32(0));
            Assert.Equal("First.", reader.GetString(1));
            Assert.True(reader.Read());
            Assert.Equal(1, reader.GetInt32(0));
            Assert.True(reader.Read());
            Assert.Equal(2, reader.GetInt32(0));
            Assert.False(reader.Read());
        }
        finally
        {
            File.Delete(dbPath);
        }
    }

    [Fact]
    public void 同じ抽出を二度書いても重複しない()
    {
        var dbPath = Path.Combine(Path.GetTempPath(), $"ef-{Guid.NewGuid():N}.sqlite3");
        try
        {
            var r = WithBook("TestBook", 0x800, "Title", "prose", "Author");
            ExtractedFieldSqliteWriter.Write(dbPath, MigrationsDir(), r);
            ExtractedFieldSqliteWriter.Write(dbPath, MigrationsDir(), r);

            using var conn = new SqliteConnection($"Data Source={dbPath}");
            conn.Open();
            using var cmd = conn.CreateCommand();
            cmd.CommandText = "SELECT COUNT(*) FROM extracted_field";
            Assert.Equal(3L, (long)cmd.ExecuteScalar()!);
        }
        finally
        {
            File.Delete(dbPath);
        }
    }
}
