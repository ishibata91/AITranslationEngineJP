using Extractor;
using Xunit;

namespace Extractor.Tests;

// 抽出モデルの意味検証（scripts/validate_extraction.py の必須 field 検証の置き換え。
// 形式は C# の型が保証するため、型で表せない不変条件だけをここで検証する）。
public class ModelInvariantTests
{
    private static ExtractionResult Dawnguard => CountParityTests.ExtractCached("Dawnguard.esm");

    [Fact]
    public void 各カテゴリのidが重複しない()
    {
        var r = Dawnguard;
        AssertUniqueIds(r.Dialogues.Select(x => x.Id), "dialogues");
        AssertUniqueIds(r.Dialogues.SelectMany(d => d.Infos).Select(x => x.Id), "infos");
        AssertUniqueIds(r.Quests.Select(x => x.Id), "quests");
        AssertUniqueIds(r.Speakers.Select(x => x.Id), "speakers");
        AssertUniqueIds(r.Words.Select(x => x.Id), "words");
        AssertUniqueIds(r.Books.Select(x => x.Id), "books");
    }

    [Fact]
    public void 主要カテゴリが空でない()
    {
        var r = Dawnguard;
        Assert.NotEmpty(r.Dialogues);
        Assert.NotEmpty(r.Quests);
        Assert.NotEmpty(r.Speakers);
        Assert.NotEmpty(r.Books);
        Assert.NotEmpty(r.Shouts);
        Assert.NotEmpty(r.Words);
    }

    [Fact]
    public void Speakerはペルソナ名を必ず持つ()
    {
        // ペルソナ名解決順（FULL → TPLT 連鎖 → EditorID）により空にならないはず。
        Assert.All(Dawnguard.Speakers, s => Assert.False(
            s.PersonaName.Length == 0 && s.EditorId.Length > 0,
            $"persona 名が空: {s.Id} ({s.EditorId})"));
    }

    [Fact]
    public void InfoNodeの話者解決は名指しか集合のどちらかに分類される()
    {
        // r2a〜d のいずれも無い INFO は許容される（シーン会話等）が、
        // 解決された ID list に null 相当（FormKey.Null）が混入してはいけない。
        foreach (var info in Dawnguard.Dialogues.SelectMany(d => d.Infos))
        {
            Assert.All(info.SpeakerIds, k => Assert.False(k.IsNull));
            Assert.All(info.FactionIds, k => Assert.False(k.IsNull));
            Assert.All(info.RaceIds, k => Assert.False(k.IsNull));
            Assert.All(info.VoiceTypeIds, k => Assert.False(k.IsNull));
        }
    }

    [Fact]
    public void stubのrecordはモデルに出るが翻訳対象に数えない()
    {
        var r = Dawnguard;
        // DialogueFollowerWaitTopic は Skyrim と同一内容の運搬用 override（stub）。
        // record としては出す（レコードが出せない状態を作らない）が、翻訳対象には数えない。
        var stub = r.Dialogues.FirstOrDefault(d => d.EditorId == "DialogueFollowerWaitTopic");
        Assert.NotNull(stub);
        Assert.True(stub!.IdenticalToMaster);
        Assert.True(stub.Name.Length > 0);
        Assert.DoesNotContain(TranslationCounts.Enumerate(r),
            s => s.RecField == "DIAL:FULL" && s.EditorId == "DialogueFollowerWaitTopic");
    }

    [Fact]
    public void 龍語綴りは翻訳対象に数えない()
    {
        var counts = TranslationCounts.Flatten(Dawnguard);
        Assert.False(counts.ContainsKey("WOOP:FULL"), "WOOP:FULL（龍語綴り、翻訳禁止）が件数に混入している");
    }

    private static void AssertUniqueIds<T>(IEnumerable<T> ids, string category)
    {
        var list = ids.ToList();
        Assert.True(list.Count == list.Distinct().Count(), $"{category} に重複 id がある");
    }
}
