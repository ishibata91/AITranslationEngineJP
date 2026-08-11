using System.Reflection;
using Extractor;
using Xunit;

namespace Extractor.Tests;

// 共有オラクル（test-oracle/specs.json）の extraction 段を照合する。書き方はフォルダの CLAUDE.md に従う:
//   1 オラクル 1 関数、[Oracle("<id>")] で id を引ける、AAA 必須、テスト独立、given は合成 esm 側。
public class OracleExtractionTests
{
    [Fact]
    [Oracle("rec-field-exhaustive-extraction")]
    public void 多様な_REC_FIELD_を箱判定せず素朴に吸い出す()
    {
        // Arrange: 多様な record を持つ合成 esm を読む。
        using var env = OracleInput.LoadEnv();
        // Act: 抽出の入口を叩く。
        var result = PluginExtractor.Extract(env);
        // Assert: 箱の分類をせず、種別の異なる原文が素朴に出る。
        Assert.Contains(result.Equipment, e => e.Description.Contains("Aventus"));
        Assert.Contains(result.Words, w => w.DragonSpelling == "FUS");
        Assert.True(result.Dialogues.Sum(d => d.Infos.Count) >= 5);
    }

    [Fact]
    [Oracle("dragon-word-untranslated")]
    public void 龍語は綴りを無訳片_語義を訳対象へ分ける()
    {
        // Arrange
        using var env = OracleInput.LoadEnv();
        // Act
        var result = PluginExtractor.Extract(env);
        // Assert: 綴り（FULL）と語義（TNAM）を別フィールドに持ち、綴りは翻訳禁止側にする。
        var word = result.Words.Single(w => w.EditorId == "TestWord");
        Assert.Equal(("FUS", "Force"), (word.DragonSpelling, word.Translation));
    }

    [Fact]
    [Oracle("speaker-from-anam")]
    public void 名指し話者を発話者に結ぶ()
    {
        // Arrange
        using var env = OracleInput.LoadEnv();
        // Act
        var result = PluginExtractor.Extract(env);
        // Assert: ANAM の話者がその台詞の発話者になる。
        Assert.Contains(OracleInput.Key(env.TargetMod, "TownGuard"), OracleInput.Info(result, "GuardGreet").SpeakerIds);
    }

    [Fact]
    [Oracle("speaker-from-condition")]
    public void 真を主張する条件から話者属性を採る()
    {
        // Arrange
        using var env = OracleInput.LoadEnv();
        // Act
        var result = PluginExtractor.Extract(env);
        // Assert: GetIs 系の肯定条件から、話者・種族・勢力・声型を採る。
        var info = OracleInput.Info(result, "ConditionSpeaker");
        Assert.Contains(OracleInput.Key(env.TargetMod, "MarketWoman"), info.SpeakerIds);
        Assert.Contains(OracleInput.Key(env.TargetMod, "NordRace"), info.RaceIds);
        Assert.Contains(OracleInput.Key(env.TargetMod, "TownGuardFaction"), info.FactionIds);
        Assert.Contains(OracleInput.Key(env.TargetMod, "MaleNord"), info.VoiceTypeIds);
    }

    [Fact]
    [Oracle("speaker-negated-condition-ignored")]
    public void 否定条件は話者解決に採らない()
    {
        // Arrange
        using var env = OracleInput.LoadEnv();
        // Act
        var result = PluginExtractor.Extract(env);
        // Assert: 肯定（TownGuard）は採り、否定（GetIsID Aventus==0）は採らない。
        var info = OracleInput.Info(result, "NegatedCondition");
        Assert.Contains(OracleInput.Key(env.TargetMod, "TownGuard"), info.SpeakerIds);
        Assert.DoesNotContain(OracleInput.Key(env.TargetMod, "AventusNPC"), info.SpeakerIds);
    }

    [Fact]
    [Oracle("voice-vtyp-only-not-flst")]
    public void 声型は_VTYP_直指定だけ採り_FLST_は採らない()
    {
        // Arrange
        using var env = OracleInput.LoadEnv();
        // Act
        var result = PluginExtractor.Extract(env);
        // Assert: VTYP 直指定を採り、FLST 指定はプール扱いで採らない。
        var info = OracleInput.Info(result, "ConditionSpeaker");
        Assert.Contains(OracleInput.Key(env.TargetMod, "MaleNord"), info.VoiceTypeIds);
        Assert.DoesNotContain(OracleInput.Key(env.TargetMod, "SingleVoicePool"), info.VoiceTypeIds);
    }

    [Fact]
    [Oracle("speaker-name-tplt-chain")]
    public void 自名なし話者は継承元の名で解決する()
    {
        // Arrange
        using var env = OracleInput.LoadEnv();
        // Act
        var result = PluginExtractor.Extract(env);
        // Assert: 自名なし NPC は TPLT 連鎖の本体（Aventus）名へ解決し、種別は「形態」。
        var speaker = result.Speakers.Single(s => s.EditorId == "InheritedGuard");
        Assert.Equal(("Aventus Aretino", "形態"), (speaker.PersonaName, speaker.SpeakerKind));
    }

    [Fact]
    [Oracle("persona-sex-from-female-flag")]
    public void 性別を女性フラグから確定する()
    {
        // Arrange
        using var env = OracleInput.LoadEnv();
        // Act
        var result = PluginExtractor.Extract(env);
        // Assert: 女性フラグの NPC は Female。
        Assert.Equal("Female", result.Speakers.Single(s => s.EditorId == "MarketWoman").Sex);
    }

    [Fact]
    [Oracle("line-emotion-extracted")]
    public void 非_Neutral_感情型を応答順で記録する()
    {
        // Arrange: 抽出結果と、感情型 staging を書く temp DB を用意する。
        using var env = OracleInput.LoadEnv();
        var result = PluginExtractor.Extract(env);
        using var db = OracleInput.NewDb();
        // Act: 感情型の書き出し入口を叩く。
        InfoEmotionSqliteWriter.Write(db.Path, OracleInput.MigrationsDir, result);
        // Assert: 非 Neutral が応答順（ordinal）と一致して記録される。
        Assert.Contains((0, "Fear"), db.Emotions(OracleInput.InfoFormId(result, "GuardWarn")));
        Assert.Contains((0, "Anger"), db.Emotions(OracleInput.InfoFormId(result, "MultiResponse")));
    }

    [Fact]
    [Oracle("line-emotion-neutral-skipped")]
    public void _Neutral_は記録せず応答順は進める()
    {
        // Arrange
        using var env = OracleInput.LoadEnv();
        var result = PluginExtractor.Extract(env);
        using var db = OracleInput.NewDb();
        // Act
        InfoEmotionSqliteWriter.Write(db.Path, OracleInput.MigrationsDir, result);
        // Assert: Neutral 単独は行なし。複数応答では Neutral（ordinal 1）を飛ばす。
        Assert.Empty(db.Emotions(OracleInput.InfoFormId(result, "GuardGreet")));
        Assert.DoesNotContain(db.Emotions(OracleInput.InfoFormId(result, "MultiResponse")), e => e.Ordinal == 1);
    }

    [Fact]
    [Oracle("condition-sex-from-getissex")]
    public void 性別条件の否定は極性を畳んで反転する()
    {
        // Arrange
        using var env = OracleInput.LoadEnv();
        var result = PluginExtractor.Extract(env);
        using var db = OracleInput.NewDb();
        // Act: 条件由来の性別の書き出し入口を叩く。
        InfoConditionSqliteWriter.Write(db.Path, OracleInput.MigrationsDir, env, result.TargetPlugin);
        // Assert: GetIsSex(Male)==0 は「男でない」= 女へ畳む。
        Assert.Equal("Female", db.Sex(OracleInput.InfoFormId(result, "ConditionSex")));
    }

    [Fact]
    [Oracle("condition-sex-male-only")]
    public void 性別条件が男性だけなら男性を記録する()
    {
        // Arrange
        using var env = OracleInput.LoadEnv();
        var result = PluginExtractor.Extract(env);
        using var db = OracleInput.NewDb();
        // Act
        InfoConditionSqliteWriter.Write(db.Path, OracleInput.MigrationsDir, env, result.TargetPlugin);
        // Assert
        Assert.Equal("Male", db.Sex(OracleInput.InfoFormId(result, "ConditionSexMale")));
    }

    [Fact]
    [Oracle("condition-sex-both-generic")]
    public void 性別条件に男女がある場合は性別を記録しない()
    {
        // Arrange
        using var env = OracleInput.LoadEnv();
        var result = PluginExtractor.Extract(env);
        using var db = OracleInput.NewDb();
        // Act
        InfoConditionSqliteWriter.Write(db.Path, OracleInput.MigrationsDir, env, result.TargetPlugin);
        // Assert
        Assert.Null(db.Sex(OracleInput.InfoFormId(result, "ConditionSexBoth")));
    }

    [Fact]
    [Oracle("condition-sex-absent")]
    public void 性別条件が無い場合は性別を記録しない()
    {
        // Arrange
        using var env = OracleInput.LoadEnv();
        var result = PluginExtractor.Extract(env);
        using var db = OracleInput.NewDb();
        // Act
        InfoConditionSqliteWriter.Write(db.Path, OracleInput.MigrationsDir, env, result.TargetPlugin);
        // Assert
        Assert.Null(db.Sex(OracleInput.InfoFormId(result, "ConditionSexNone")));
    }

    [Fact]
    [Oracle("condition-sex-voicelist-mixed")]
    public void 男女混在の声型リストは性別を定めない()
    {
        // Arrange
        using var env = OracleInput.LoadEnv();
        var result = PluginExtractor.Extract(env);
        using var db = OracleInput.NewDb();
        // Act
        InfoConditionSqliteWriter.Write(db.Path, OracleInput.MigrationsDir, env, result.TargetPlugin);
        // Assert: 混在なら性別を定めない（行を作らない）。
        Assert.Null(db.Sex(OracleInput.InfoFormId(result, "ConditionSexMixed")));
    }

    [Fact]
    [Oracle("multivalue-ordinal-numbering")]
    public void 複数応答は出現順で並ぶ()
    {
        // Arrange
        using var env = OracleInput.LoadEnv();
        // Act
        var result = PluginExtractor.Extract(env);
        // Assert: 同一 INFO の複数応答が出現順（ordinal 採番の素）で並ぶ。
        Assert.Equal(new[] { 1, 2 }, OracleInput.Info(result, "MultiResponse").Responses.Select(r => r.Number).ToArray());
    }

    // 網羅番人: [Oracle] を付けたテストの id 集合と、specs.json の extraction 非委任 id 集合が完全一致すること。
    // spec を足して関数を書き忘れたら（または委任印を外して関数を足さなかったら）ここで落ちる。
    [Fact]
    public void 網羅番人_登録テストと_specs_が一致する()
    {
        // Arrange: specs.json の extraction 非委任 id。
        var specIds = OracleSpecs.Load()
            .Where(s => s.Stage == "extraction" && s.Coverage != "existing-unit")
            .Select(s => s.Id).OrderBy(x => x).ToArray();
        // Act: この class の [Oracle] 付きテストが持つ id。
        var covered = typeof(OracleExtractionTests).GetMethods()
            .Select(m => m.GetCustomAttribute<OracleAttribute>())
            .Where(a => a != null)
            .Select(a => a!.Id).OrderBy(x => x).ToArray();
        // Assert: 完全一致（登録漏れ・余剰を落とす）。
        Assert.Equal(specIds, covered);
    }
}
