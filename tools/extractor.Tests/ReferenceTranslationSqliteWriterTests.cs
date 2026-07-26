using Extractor;
using Mutagen.Bethesda.Plugins;
using Xunit;

namespace Extractor.Tests;

// reference_translation writer の検証。english と japanese の両方を解決できた field だけを既訳として書き、
// 同じ内容を 2 回書いても増えないこと（Data フォルダを毎回走査するため冪等が要る）。
// schema は db/migrations 全体を ensure する（reference_translation は 0012）。
public class ReferenceTranslationSqliteWriterTests
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

    // 英語だけ解決できた field（DESC）と、英日対のある field（FULL）を持つ本 1 冊。
    private static ExtractionResult WithPairedBook()
    {
        var result = new ExtractionResult { TargetPlugin = "Test.esm" };
        var id = new FormKey(ModKey.FromFileName("Test.esm"), 0x800);
        result.Books.Add(new BookEntry
        {
            Id = id,
            EditorId = "TestBook",
            Kind = "BOOK",
            Title = "Lusty Argonian Maid",
            Body = "Ancient Nord prose",
            Author = "Crassius Curio",
        });
        result.JapanesePairs[(id, "BOOK:FULL", "Lusty Argonian Maid")] = "アルゴニアンの侍女";
        return result;
    }

    private static List<(string rec, string field, string source, string dest)> ReadAll(TempSqliteDb db)
    {
        using var conn = db.OpenConnection();
        using var cmd = conn.CreateCommand();
        cmd.CommandText = "SELECT rec, field, source, dest FROM reference_translation ORDER BY rec, field, source";
        using var reader = cmd.ExecuteReader();
        var got = new List<(string, string, string, string)>();
        while (reader.Read())
            got.Add((reader.GetString(0), reader.GetString(1), reader.GetString(2), reader.GetString(3)));
        return got;
    }

    // R-1-3: 日本語側を解決できなかった record を、既訳として供給しないこと。
    [Fact]
    public void 日本語を解決できた_field_だけを既訳として書く()
    {
        using var db = new TempSqliteDb("ref");

        var written = ReferenceTranslationSqliteWriter.Write(db.Path, MigrationsDir(), WithPairedBook());

        Assert.Equal(1, written);
        var got = ReadAll(db);
        Assert.Equal(new[] { ("BOOK", "FULL", "Lusty Argonian Maid", "アルゴニアンの侍女") }, got);
    }

    // R-1-4: 英日対を集める段を 2 回続けて走らせても、英日対の件数が増えないこと。
    // Data フォルダは翻訳のたびに走査するため、同じ内容の再走査で既訳が膨らまないことが要る。
    [Fact]
    public void 同じ英日対を二度書いても件数が増えない()
    {
        using var db = new TempSqliteDb("ref-idem");

        var first = ReferenceTranslationSqliteWriter.Write(db.Path, MigrationsDir(), WithPairedBook());
        var countAfterFirst = ReadAll(db).Count;
        var second = ReferenceTranslationSqliteWriter.Write(db.Path, MigrationsDir(), WithPairedBook());

        Assert.Equal(1, first);
        Assert.Equal(0, second); // INSERT OR IGNORE で 2 回目は 1 件も足さない
        Assert.Equal(countAfterFirst, ReadAll(db).Count);
    }

    // 既訳は plugin をまたいで同一原文へ再利用するため、由来 plugin で行を分けない。
    // 別 plugin の抽出結果が同じ (rec, field, source) を運んでも 1 行のまま（先に入れた訳を保つ）。
    [Fact]
    public void 別_plugin_の同じ原文は行を増やさない()
    {
        using var db = new TempSqliteDb("ref-cross");

        ReferenceTranslationSqliteWriter.Write(db.Path, MigrationsDir(), WithPairedBook());

        var other = new ExtractionResult { TargetPlugin = "Other.esp" };
        var otherId = new FormKey(ModKey.FromFileName("Other.esp"), 0x900);
        other.Books.Add(new BookEntry
        {
            Id = otherId, EditorId = "OtherBook", Kind = "BOOK",
            Title = "Lusty Argonian Maid", Body = "", Author = "",
        });
        other.JapanesePairs[(otherId, "BOOK:FULL", "Lusty Argonian Maid")] = "別訳";

        Assert.Equal(0, ReferenceTranslationSqliteWriter.Write(db.Path, MigrationsDir(), other));
        Assert.Equal(new[] { ("BOOK", "FULL", "Lusty Argonian Maid", "アルゴニアンの侍女") }, ReadAll(db));
    }
}
