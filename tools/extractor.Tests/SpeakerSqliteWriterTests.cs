using Extractor;
using Microsoft.Data.Sqlite;
using Mutagen.Bethesda;
using Mutagen.Bethesda.Plugins;
using Mutagen.Bethesda.Plugins.Cache;
using Mutagen.Bethesda.Skyrim;
using Xunit;

namespace Extractor.Tests;

// 話者 writer の検証。INFO→speaker の橋渡し（extracted_info_speaker）を db/migrations の schema へ冪等に書くこと。
// 話者解決（NPC→種族/声型/勢力）の網羅は実 mod 抽出で確認するため、本テストは話者が解決できない場合の
// 振る舞い（橋渡しを書かない・落ちない）に絞り、LinkCache は空のものを渡す。台詞本文は ExtractedFieldSqliteWriter が書く。
public class SpeakerSqliteWriterTests
{
    private static string RepoRoot()
    {
        var dir = new DirectoryInfo(AppContext.BaseDirectory);
        while (dir != null && !File.Exists(Path.Combine(dir.FullName, "go.mod")))
            dir = dir.Parent;
        Assert.NotNull(dir);
        return dir!.FullName;
    }

    // migrations ディレクトリ。Writer が SchemaMigrator で ensure する。
    private static string MigrationsDir() => Path.Combine(RepoRoot(), "db", "migrations");

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
    public void 話者が解決できない_INFO_は橋渡しを書かない()
    {
        var dbPath = Path.Combine(Path.GetTempPath(), $"sp-{Guid.NewGuid():N}.sqlite3");
        try
        {
            // 空 LinkCache では話者 NPC を解決できない（info.SpeakerIds も空）。橋渡し 0 件で落ちないこと。
            var count = SpeakerSqliteWriter.Write(dbPath, MigrationsDir(), WithInfo("InfoA", 0x800, "Hello.", "Goodbye."), EmptyLinkCache());
            Assert.Equal(0, count);

            using var conn = new SqliteConnection($"Data Source={dbPath}");
            conn.Open();
            using var cmd = conn.CreateCommand();
            cmd.CommandText = "SELECT COUNT(*) FROM extracted_info_speaker";
            Assert.Equal(0L, (long)cmd.ExecuteScalar()!);
        }
        finally
        {
            File.Delete(dbPath);
        }
    }

    [Fact]
    public void 二度書いても落ちない_冪等()
    {
        var dbPath = Path.Combine(Path.GetTempPath(), $"sp-{Guid.NewGuid():N}.sqlite3");
        try
        {
            SpeakerSqliteWriter.Write(dbPath, MigrationsDir(), WithInfo("InfoA", 0x800, "Hello."), EmptyLinkCache());
            SpeakerSqliteWriter.Write(dbPath, MigrationsDir(), WithInfo("InfoA", 0x800, "Hello."), EmptyLinkCache());

            using var conn = new SqliteConnection($"Data Source={dbPath}");
            conn.Open();
            using var cmd = conn.CreateCommand();
            cmd.CommandText = "SELECT COUNT(*) FROM extracted_info_speaker";
            Assert.Equal(0L, (long)cmd.ExecuteScalar()!);
        }
        finally
        {
            File.Delete(dbPath);
        }
    }
}
