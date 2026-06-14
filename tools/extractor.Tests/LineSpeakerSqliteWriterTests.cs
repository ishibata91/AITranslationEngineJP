using Extractor;
using Microsoft.Data.Sqlite;
using Mutagen.Bethesda;
using Mutagen.Bethesda.Plugins;
using Mutagen.Bethesda.Plugins.Cache;
using Mutagen.Bethesda.Skyrim;
using Xunit;

namespace Extractor.Tests;

// 台詞 SQLite writer の検証。台詞（INFO:NAM1）を db/migrations の schema へ冪等に書くこと。
// 話者解決（NPC→種族/声型/勢力）の網羅は実 mod 抽出で確認するため、本テストは話者なし台詞の書込と
// 冪等性に絞り、LinkCache は空のものを渡す。
public class LineSpeakerSqliteWriterTests
{
    private static string RepoRoot()
    {
        var dir = new DirectoryInfo(AppContext.BaseDirectory);
        while (dir != null && !File.Exists(Path.Combine(dir.FullName, "go.mod")))
            dir = dir.Parent;
        Assert.NotNull(dir);
        return dir!.FullName;
    }

    // 0001 と 0002 を結合した schema（line/speaker テーブルは 0002）。
    private static string Schema()
    {
        var migrations = Path.Combine(RepoRoot(), "db", "migrations");
        return string.Join("\n",
            Directory.GetFiles(migrations, "*.sql").OrderBy(p => p, StringComparer.Ordinal).Select(File.ReadAllText));
    }

    private static ILinkCache EmptyLinkCache()
        => new SkyrimMod(ModKey.FromFileName("Empty.esm"), SkyrimRelease.SkyrimSE).ToImmutableLinkCache();

    private static ExtractionResult WithInfo(string edid, uint id, params string[] responses)
    {
        var result = new ExtractionResult { TargetPlugin = "Test.esp" };
        var info = new InfoNode
        {
            Id = new FormKey(ModKey.FromFileName("Test.esp"), id),
            EditorId = edid,
            Prompt = "",
        };
        for (var i = 0; i < responses.Length; i++)
            info.Responses.Add(new ResponseLine(i + 1, responses[i]));

        result.Dialogues.Add(new DialogueTopic
        {
            Id = new FormKey(ModKey.FromFileName("Test.esp"), id + 1),
            EditorId = "Topic",
            Kind = "DIAL",
            Name = "",
            Category = "",
            Subtype = "",
            Infos = { info },
        });
        return result;
    }

    [Fact]
    public void INFO_NAM1_を_line_へ書く()
    {
        var dbPath = Path.Combine(Path.GetTempPath(), $"lw-{Guid.NewGuid():N}.sqlite3");
        try
        {
            var count = LineSpeakerSqliteWriter.Write(dbPath, Schema(), WithInfo("InfoA", 0x800, "Hello.", "Goodbye."), EmptyLinkCache());
            Assert.Equal(2, count);

            using var conn = new SqliteConnection($"Data Source={dbPath}");
            conn.Open();
            using var cmd = conn.CreateCommand();
            cmd.CommandText = "SELECT plugin, form_id, edid, rec, field, ordinal, response_order, source, status FROM line ORDER BY ordinal";
            using var reader = cmd.ExecuteReader();

            Assert.True(reader.Read());
            Assert.Equal("Test.esp", reader.GetString(0));
            Assert.Equal("0x000800", reader.GetString(1));
            Assert.Equal("InfoA", reader.GetString(2));
            Assert.Equal("INFO", reader.GetString(3));
            Assert.Equal("NAM1", reader.GetString(4));
            Assert.Equal(1, reader.GetInt32(5));  // ordinal = response 番号
            Assert.Equal(1, reader.GetInt32(6));  // response_order
            Assert.Equal("Hello.", reader.GetString(7));
            Assert.Equal(0, reader.GetInt32(8));  // status = 未訳

            Assert.True(reader.Read());
            Assert.Equal("Goodbye.", reader.GetString(7));
            Assert.Equal(2, reader.GetInt32(5));
            Assert.False(reader.Read());
        }
        finally
        {
            File.Delete(dbPath);
        }
    }

    [Fact]
    public void 同じ台詞を二度書いても重複しない()
    {
        var dbPath = Path.Combine(Path.GetTempPath(), $"lw-{Guid.NewGuid():N}.sqlite3");
        try
        {
            LineSpeakerSqliteWriter.Write(dbPath, Schema(), WithInfo("InfoA", 0x800, "Hello."), EmptyLinkCache());
            LineSpeakerSqliteWriter.Write(dbPath, Schema(), WithInfo("InfoA", 0x800, "Hello."), EmptyLinkCache());

            using var conn = new SqliteConnection($"Data Source={dbPath}");
            conn.Open();
            using var cmd = conn.CreateCommand();
            cmd.CommandText = "SELECT COUNT(*) FROM line";
            Assert.Equal(1L, (long)cmd.ExecuteScalar()!);
        }
        finally
        {
            File.Delete(dbPath);
        }
    }
}
