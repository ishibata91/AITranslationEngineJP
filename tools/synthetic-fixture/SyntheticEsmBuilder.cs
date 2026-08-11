using Mutagen.Bethesda;
using Mutagen.Bethesda.Plugins;
using Mutagen.Bethesda.Plugins.Records;
using Mutagen.Bethesda.Skyrim;

namespace SyntheticFixture;

// 合成 esm（master 無しの自己完結 plugin）を Mutagen で組む。C# 抽出テストと Go 派生 seed が読む唯一の合成入力。
// 各 record の EditorID は specs.json の given が名指せるよう固定名にする。著作物は含まない自作データ。
//
// master・localized 形式が要る 3 spec（override 同一判定・container 所有・localized header 除外）は
// この master 無し esm では再現しない。既存の実 esm 単体テストへ委ねる（specs.json の coverage=existing-unit）。
public static class SyntheticEsmBuilder
{
    public const string PluginName = "Synthetic.esm";

    public static SkyrimMod Build()
    {
        var mod = new SkyrimMod(ModKey.FromFileName(PluginName), SkyrimRelease.SkyrimSE);

        // 種族: 話者の年齢区分・種族訛りの引き元。老年区分を通すため ElderRace を置く。
        var nordRace = mod.Races.AddNew();
        nordRace.EditorID = "NordRace";
        nordRace.Name = "Nord";
        var imperialRace = mod.Races.AddNew();
        imperialRace.EditorID = "ImperialRace";
        imperialRace.Name = "Imperial";

        // 声型: 対人 prior の引き元。1 話者 1 声型。
        var maleNord = mod.VoiceTypes.AddNew();
        maleNord.EditorID = "MaleNord";
        var femaleCondescending = mod.VoiceTypes.AddNew();
        femaleCondescending.EditorID = "FemaleCondescending";

        // 話者（NPC_）: 女性フラグ・種族・声型を持つ。名前解決と性別確定を観測する。
        var townGuard = mod.Npcs.AddNew();
        townGuard.EditorID = "TownGuard";
        townGuard.Name = "Town Guard";
        townGuard.Race.SetTo(nordRace);
        townGuard.Voice.SetTo(maleNord);

        var marketWoman = mod.Npcs.AddNew();
        marketWoman.EditorID = "MarketWoman";
        marketWoman.Name = "Market Woman";
        marketWoman.Configuration.Flags |= NpcConfiguration.Flag.Female; // 性別を女性フラグから確定する
        marketWoman.Race.SetTo(imperialRace);
        marketWoman.Voice.SetTo(femaleCondescending);

        // 固有名（NPC_:FULL）: 本文へ機械置換する訳の単位。姓名分割の派生元にもなる。
        var aventus = mod.Npcs.AddNew();
        aventus.EditorID = "AventusNPC";
        aventus.Name = "Aventus Aretino";
        aventus.Race.SetTo(imperialRace);
        aventus.Voice.SetTo(maleNord);

        // 継承元あり・自名なしの話者（NPC_:TPLT）: 名前を 自名 → 継承元の名 → 役割名 の順で解決する。
        // 自名を持たず Template を aventus（名前あり）へ向け、継承元の名で解決されることを観測する。
        var inheritedNpc = mod.Npcs.AddNew();
        inheritedNpc.EditorID = "InheritedGuard";
        inheritedNpc.Template.SetTo(aventus); // 名前は継承元（aventus）から採る
        inheritedNpc.Race.SetTo(nordRace);
        inheritedNpc.Voice.SetTo(maleNord);

        // 勢力（FACT）: INFO の GetInFaction 条件で話者を絞る引き元。
        var townFaction = mod.Factions.AddNew();
        townFaction.EditorID = "TownGuardFaction";
        townFaction.Name = "Town Guard Faction";

        // 声型リスト（FLST）: 声型条件がリスト指定のとき。voice-vtyp-only 用は単一構成、混在性別用は男女混在。
        var singleVoicePool = mod.FormLists.AddNew();
        singleVoicePool.EditorID = "SingleVoicePool";
        singleVoicePool.Items.Add(new FormLink<ISkyrimMajorRecordGetter>(maleNord.FormKey));

        var mixedVoicePool = mod.FormLists.AddNew();
        mixedVoicePool.EditorID = "MixedVoicePool";
        mixedVoicePool.Items.Add(new FormLink<ISkyrimMajorRecordGetter>(maleNord.FormKey));
        mixedVoicePool.Items.Add(new FormLink<ISkyrimMajorRecordGetter>(femaleCondescending.FormKey));

        // 叙述文（説明体 WEAP:DESC）: 本文中の固有名 Aventus に機械置換が当たる。
        var sword = mod.Weapons.AddNew();
        sword.EditorID = "TestSword";
        sword.Name = "Test Sword";
        sword.Description = "A blade once held by Aventus Aretino.";

        // 叙述文（文頭 No の WEAP:DESC）: stoplist 選別後、原文のままプロンプトへ渡る。
        var dagger = mod.Weapons.AddNew();
        dagger.EditorID = "TestDagger";
        dagger.Name = "Test Dagger";
        dagger.Description = "No matter what the weather, this blade holds its edge.";

        // 叙述文（書物体 BOOK:DESC）: 文体 directive が説明体と分かれる。末尾 Grelod で二つ名前部派生を観測する。
        var book = mod.Books.AddNew();
        book.EditorID = "TestBook";
        book.Name = "Test Book";
        book.BookText = "A long account of the old kings of the north, told by Grelod.";
        book.Description = ""; // CNAM（著者）は置かない

        // 定型句（ACTI:RNAM）: 固有名でも台詞でもない箱。
        var activator = mod.Activators.AddNew();
        activator.EditorID = "TestActi";
        activator.Name = "Test Door";
        activator.ActivateTextOverride = "Open";

        // 無訳片（WOOP）: 龍語の綴り（FULL）は翻訳禁止、語義（TNAM）だけ訳対象の定型句。
        var word = mod.WordsOfPower.AddNew();
        word.EditorID = "TestWord";
        word.Name = "FUS"; // 龍語綴り（翻訳禁止）
        word.Translation = "Force"; // 語義（訳対象）

        // 話者付き会話（DIAL→INFO→NAM1）: 名指し話者（ANAM）と応答順（ordinal）を観測する。
        var greetTopic = mod.DialogTopics.AddNew();
        greetTopic.EditorID = "GuardGreetTopic";
        var greetInfo = new DialogResponses(mod);
        greetInfo.EditorID = "GuardGreet";
        greetInfo.Speaker.SetTo(townGuard); // 名指し話者を発話者として結ぶ
        greetInfo.Responses.Add(new DialogResponse
        {
            ResponseNumber = 1,
            Text = "Have you come to help me with the trouble in town?",
            Emotion = Emotion.Neutral, // Neutral は感情型として記録しない
        });
        greetTopic.Responses.Add(greetInfo);

        // 感情型（TRDT）付き台詞: 非 Neutral の感情を台詞単位で記録する。本文 fear で感情段階も観測する。
        var warnTopic = mod.DialogTopics.AddNew();
        warnTopic.EditorID = "GuardWarnTopic";
        var warnInfo = new DialogResponses(mod);
        warnInfo.EditorID = "GuardWarn";
        warnInfo.Speaker.SetTo(townGuard);
        warnInfo.Responses.Add(new DialogResponse
        {
            ResponseNumber = 1,
            Text = "I fear Aventus will not come back to town.",
            Emotion = Emotion.Fear, // 非 Neutral の感情型
        });
        warnTopic.Responses.Add(warnInfo);

        // 条件で話者を絞る会話（INFO の CTDA）: 真を主張する GetIs 条件から話者属性を採る。
        //  - GetIsID(NPC)==1 → 話者、GetIsRace==1 → 種族、GetInFaction==1 → 勢力、GetIsVoiceType(VTYP)==1 → 声型。
        //  - GetIsVoiceType(FLST)==1 は声型プール扱いで声型に採らない（voice-vtyp-only-not-flst）。
        var condTopic = mod.DialogTopics.AddNew();
        condTopic.EditorID = "ConditionSpeakerTopic";
        var condInfo = new DialogResponses(mod);
        condInfo.EditorID = "ConditionSpeaker";
        condInfo.Conditions.Add(GetIsIdCond(marketWoman, asserts: true));      // 話者
        condInfo.Conditions.Add(GetIsRaceCond(nordRace));                       // 種族
        condInfo.Conditions.Add(GetInFactionCond(townFaction));                // 勢力
        condInfo.Conditions.Add(GetIsVoiceTypeCond(maleNord.FormKey));         // 声型（VTYP 直指定→採る）
        condInfo.Conditions.Add(GetIsVoiceTypeCond(singleVoicePool.FormKey));  // 声型（FLST→プール扱いで採らない）
        condInfo.Responses.Add(new DialogResponse
        {
            ResponseNumber = 1,
            Text = "The guard watches the market gate.",
            Emotion = Emotion.Neutral,
        });
        condTopic.Responses.Add(condInfo);

        // 否定条件を含む会話: 否定形（GetIsID X==0）は「X ではない」の主張なので話者解決に採らない。
        var negTopic = mod.DialogTopics.AddNew();
        negTopic.EditorID = "NegatedConditionTopic";
        var negInfo = new DialogResponses(mod);
        negInfo.EditorID = "NegatedCondition";
        negInfo.Conditions.Add(GetIsIdCond(townGuard, asserts: true));   // 肯定→話者に採る
        negInfo.Conditions.Add(GetIsIdCond(aventus, asserts: false));    // 否定→話者に採らない
        negInfo.Responses.Add(new DialogResponse
        {
            ResponseNumber = 1,
            Text = "Only the guard speaks here.",
            Emotion = Emotion.Neutral,
        });
        negTopic.Responses.Add(negInfo);

        // 性別条件（GetIsSex）の会話: 実効性別を導き、否定形は極性を畳んで反転する。
        // GetIsSex(Male)==0 は「男でない」= 女。話者解決の有無に依らず性別を定める。
        var sexTopic = mod.DialogTopics.AddNew();
        sexTopic.EditorID = "ConditionSexTopic";
        var sexInfo = new DialogResponses(mod);
        sexInfo.EditorID = "ConditionSex";
        sexInfo.Conditions.Add(GetIsSexCond(MaleFemaleGender.Male, asserts: false)); // 男でない→女へ畳む
        sexInfo.Responses.Add(new DialogResponse
        {
            ResponseNumber = 1,
            Text = "A voice from the shadows warns you.",
            Emotion = Emotion.Neutral,
        });
        sexTopic.Responses.Add(sexInfo);

        // 男性条件だけの会話: 男性だけなら性別を定める。
        var maleSexTopic = mod.DialogTopics.AddNew();
        maleSexTopic.EditorID = "ConditionSexMaleTopic";
        var maleSexInfo = new DialogResponses(mod);
        maleSexInfo.EditorID = "ConditionSexMale";
        maleSexInfo.Conditions.Add(GetIsSexCond(MaleFemaleGender.Male, asserts: true));
        maleSexInfo.Responses.Add(new DialogResponse
        {
            ResponseNumber = 1,
            Text = "A man calls out from the road.",
            Emotion = Emotion.Neutral,
        });
        maleSexTopic.Responses.Add(maleSexInfo);

        // 男女の性別条件を同じ INFO に持つ会話: 両方が現れるため性別を定めない。
        var bothSexTopic = mod.DialogTopics.AddNew();
        bothSexTopic.EditorID = "ConditionSexBothTopic";
        var bothSexInfo = new DialogResponses(mod);
        bothSexInfo.EditorID = "ConditionSexBoth";
        bothSexInfo.Conditions.Add(GetIsSexCond(MaleFemaleGender.Male, asserts: true));
        bothSexInfo.Conditions.Add(GetIsSexCond(MaleFemaleGender.Female, asserts: true));
        bothSexInfo.Responses.Add(new DialogResponse
        {
            ResponseNumber = 1,
            Text = "Someone calls out from the road.",
            Emotion = Emotion.Neutral,
        });
        bothSexTopic.Responses.Add(bothSexInfo);

        // 性別条件が無い会話: 条件由来の性別を定めない。
        var noSexTopic = mod.DialogTopics.AddNew();
        noSexTopic.EditorID = "ConditionSexNoneTopic";
        var noSexInfo = new DialogResponses(mod);
        noSexInfo.EditorID = "ConditionSexNone";
        noSexInfo.Responses.Add(new DialogResponse
        {
            ResponseNumber = 1,
            Text = "A distant voice echoes nearby.",
            Emotion = Emotion.Neutral,
        });
        noSexTopic.Responses.Add(noSexInfo);

        // 声型リストに異性が混在する会話: 混在なら性別を定めない（同性のみで揃うリストだけ性別を採る）。
        var mixTopic = mod.DialogTopics.AddNew();
        mixTopic.EditorID = "ConditionSexMixedTopic";
        var mixInfo = new DialogResponses(mod);
        mixInfo.EditorID = "ConditionSexMixed";
        mixInfo.Conditions.Add(GetIsVoiceTypeCond(mixedVoicePool.FormKey)); // 男女混在 FLST→性別を定めない
        mixInfo.Responses.Add(new DialogResponse
        {
            ResponseNumber = 1,
            Text = "An unseen speaker mutters nearby.",
            Emotion = Emotion.Neutral,
        });
        mixTopic.Responses.Add(mixInfo);

        // 同一 INFO に複数応答: 出現順で ordinal を採番し、感情・話者・性別の結線が正しい値へ対応できるようにする。
        var multiTopic = mod.DialogTopics.AddNew();
        multiTopic.EditorID = "MultiResponseTopic";
        var multiInfo = new DialogResponses(mod);
        multiInfo.EditorID = "MultiResponse";
        multiInfo.Speaker.SetTo(townGuard);
        multiInfo.Responses.Add(new DialogResponse
        {
            ResponseNumber = 1,
            Text = "Halt. State your business.",
            Emotion = Emotion.Anger, // 応答 1 に感情型
        });
        multiInfo.Responses.Add(new DialogResponse
        {
            ResponseNumber = 2,
            Text = "Move along, then.",
            Emotion = Emotion.Neutral, // 応答 2 は Neutral（ordinal は進むが感情は記録しない）
        });
        multiTopic.Responses.Add(multiInfo);

        return mod;
    }

    // ===== CTDA（INFO 条件）ビルダ =====
    // asserts=true は「関数が真（==1）」、false は否定（==0）。抽出側 AssertsTrue と対で解釈する。

    private static ConditionFloat Cond(ConditionData data, bool asserts) => new()
    {
        Data = data,
        CompareOperator = CompareOperator.EqualTo,
        ComparisonValue = asserts ? 1f : 0f,
    };

    private static ConditionFloat GetIsIdCond(ISkyrimMajorRecordGetter target, bool asserts)
    {
        var data = new GetIsIDConditionData();
        data.Object.Link.SetTo(target.FormKey);
        return Cond(data, asserts);
    }

    private static ConditionFloat GetIsRaceCond(IRaceGetter race)
    {
        var data = new GetIsRaceConditionData();
        data.Race.Link.SetTo(race.FormKey);
        return Cond(data, asserts: true);
    }

    private static ConditionFloat GetInFactionCond(IFactionGetter faction)
    {
        var data = new GetInFactionConditionData();
        data.Faction.Link.SetTo(faction.FormKey);
        return Cond(data, asserts: true);
    }

    private static ConditionFloat GetIsVoiceTypeCond(FormKey voiceOrList)
    {
        var data = new GetIsVoiceTypeConditionData();
        data.VoiceTypeOrList.Link.SetTo(voiceOrList);
        return Cond(data, asserts: true);
    }

    private static ConditionFloat GetIsSexCond(MaleFemaleGender gender, bool asserts)
    {
        var data = new GetIsSexConditionData { MaleFemaleGender = gender };
        return Cond(data, asserts);
    }
}
