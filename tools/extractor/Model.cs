using Mutagen.Bethesda.Plugins;

namespace Extractor;

// Skyrim 構造体モデル v19（docs/skyrim-structure-model.md）に対応する抽出結果モデル。
// plugin 単位列挙のため、全 entry の source は target plugin になる（field としては持たない）。
// 参照は FormKey で保持し、参照先本体が target plugin に無い場合は ID のみ残る（dangling 許容）。

public sealed class ExtractionResult
{
    public required string TargetPlugin { get; init; }
    // plugin header（TES4）の作者（CNAM）と説明（SNAM）。
    // 非 localized plugin では xTranslator が翻訳対象として列挙する。
    public string HeaderAuthor { get; init; } = "";
    public string HeaderDescription { get; init; } = "";
    public List<DialogueTopic> Dialogues { get; } = [];
    public List<QuestEntry> Quests { get; } = [];
    public List<RaceEntry> Races { get; } = [];
    public List<FactionEntry> Factions { get; } = [];
    public List<NamedEntry> Items { get; } = [];          // KEYM/MISC/LIGH/CONT/SLGM/DOOR/FURN/APPA/HAZD
    public List<ActivatorEntry> Activators { get; } = []; // ACTI/FLOR/TREE
    public List<DescribedEntry> Equipment { get; } = [];  // WEAP/ARMO/AMMO
    public List<DescribedEntry> Consumables { get; } = [];// ALCH/SCRL/INGR
    public List<DescribedEntry> Magic { get; } = [];      // SPEL
    public List<DescribedEntry> Enchantments { get; } = [];// ENCH
    public List<DescribedEntry> MagicEffects { get; } = [];// MGEF（description = DNAM）
    public List<ShoutEntry> Shouts { get; } = [];         // SHOU
    public List<WordOfPowerEntry> Words { get; } = [];    // WOOP（v19: top-level。shout からは ID 参照）
    public List<BookEntry> Books { get; } = [];
    public List<NamedEntry> Locations { get; } = [];      // LCTN/WRLD/CELL
    public List<MessageEntry> Messages { get; } = [];
    public List<BodyEntry> LoadingScreens { get; } = [];  // LSCR（body = DESC）
    public List<DescribedEntry> Perks { get; } = [];
    public List<SpeakerEntry> Speakers { get; } = [];     // NPC_/TACT
    public List<VoiceTypeEntry> VoiceTypes { get; } = []; // VTYP（翻訳テキスト無し、識別子のみ）
    public List<NamedEntry> SoundCategories { get; } = [];// SNCT
    public List<NamedEntry> Eyes { get; } = [];           // EYES
    public List<RegionEntry> Regions { get; } = [];       // REGN（map_name = RDMP）
}

public abstract class RecordEntry
{
    public required FormKey Id { get; init; }
    public required string EditorId { get; init; }
    // REC base（"WEAP" 等）。flatten 時の REC:FIELD 構築に使う。
    public required string Kind { get; init; }
    // master と同一内容の override（参照用 stub）。record としては出すが、
    // 翻訳所有は master 側にあるため、この plugin の翻訳対象には数えない。
    public bool IdenticalToMaster { get; init; }
    // プレイヤーに表示されないと推定される text を持つ record の hint。
    // 対象: journal 文（CNAM/NNAM）を持たないクエスト名（script / scene 制御用）、
    // HideInUI flag の MGEF、NonPlayable flag の装備（WEAP/ARMO/AMMO）。
    // xTranslator 辞書はこれらも翻訳対象に数えるため件数比較には影響させず、
    // 読み込み側（翻訳エンジン）が翻訳要否の判断に使う。QUST は構造からの推定で
    // あり確定情報ではないため、除外ではなく hint として持つ。
    public bool NotPlayerFacing { get; init; }
}

// FULL のみ持つ record（item / location / SNCT / EYES）。
public sealed class NamedEntry : RecordEntry
{
    public required string Name { get; init; }
}

// FULL + DESC（または DNAM）を持つ record。
public sealed class DescribedEntry : RecordEntry
{
    public required string Name { get; init; }
    public required string Description { get; init; }
    public List<FormKey> MagicEffectIds { get; init; } = [];
    public FormKey? EnchantmentId { get; init; }
}

public sealed class BodyEntry : RecordEntry
{
    public required string Body { get; init; }
}

public sealed class ActivatorEntry : RecordEntry
{
    public required string Name { get; init; }
    public required string ActivateAction { get; init; } // RNAM
}

public sealed class RaceEntry : RecordEntry
{
    public required string Name { get; init; }
    public required string Description { get; init; }
}

public sealed class FactionEntry : RecordEntry
{
    public required string Name { get; init; }
    public List<RankTitle> Ranks { get; init; } = [];
}

// FACT 階級称号（MNAM = 男性、FNAM = 女性）。
public sealed record RankTitle(string Male, string Female);

public sealed class QuestEntry : RecordEntry
{
    public required string Name { get; init; }
    public required string QuestType { get; init; }
    public List<QuestStageEntry> Stages { get; init; } = [];
    public List<QuestObjectiveEntry> Objectives { get; init; } = [];
}

public sealed record QuestStageEntry(int StageIndex, List<string> LogEntries);
public sealed record QuestObjectiveEntry(int ObjectiveIndex, string DisplayText);

public sealed class BookEntry : RecordEntry
{
    public required string Title { get; init; }       // FULL
    public required string Body { get; init; }        // DESC（本文）
    public required string Author { get; init; }      // CNAM（SSEEdit 表示名 Description）
}

public sealed class MessageEntry : RecordEntry
{
    public required string Title { get; init; }       // FULL
    public required string Body { get; init; }        // DESC
    public FormKey? QuestId { get; init; }
    public List<string> Buttons { get; init; } = [];  // ITXT
}

public sealed class ShoutEntry : RecordEntry
{
    public required string Name { get; init; }
    public required string Description { get; init; }
    public List<FormKey> WordIds { get; init; } = [];
}

public sealed class WordOfPowerEntry : RecordEntry
{
    public required string DragonSpelling { get; init; } // FULL（龍語綴り、翻訳禁止）
    public required string Translation { get; init; }    // TNAM（意味の訳、翻訳対象）
}

// VTYP は翻訳テキストを持たない。EditorId がそのまま識別子になる。
public sealed class VoiceTypeEntry : RecordEntry;

public sealed class RegionEntry : RecordEntry
{
    public required string MapName { get; init; } // RDMP
}

public sealed class SpeakerEntry : RecordEntry
{
    public required string Name { get; init; }
    public required string PersonaName { get; init; }
    public required string SpeakerKind { get; init; } // 固有/形態/プール/声型代表/役割のみ/喋るオブジェクト
    public required string ShortName { get; init; }   // SHRT（NPC_ のみ）
    public required string Sex { get; init; }
    public FormKey? VoiceTypeId { get; init; }
    public FormKey? RaceId { get; init; }
    public FormKey? BaseSpeakerId { get; init; }      // TPLT
    public List<FormKey> FactionIds { get; init; } = [];
}

// 会話 2 階層（v19）: Dialogue（DIAL）→ InfoNode（INFO）→ ResponseLine（NAM1）。
public sealed class DialogueTopic : RecordEntry
{
    public required string Name { get; init; }        // FULL（プレイヤー選択肢の既定文）
    public required string Category { get; init; }
    public required string Subtype { get; init; }
    public FormKey? QuestId { get; init; }
    public List<InfoNode> Infos { get; init; } = [];
}

public sealed class InfoNode
{
    public required FormKey Id { get; init; }
    public required string EditorId { get; init; }
    // master と同一内容の override（stub）。翻訳対象に数えない。
    public bool IdenticalToMaster { get; init; }
    public required string Prompt { get; init; }      // RNAM（プレイヤー選択肢の条件別上書き）
    public List<ResponseLine> Responses { get; init; } = [];
    // 話者解決（r2a〜d）: 名指し / 派閥 / 種族 / 声型。
    public List<FormKey> SpeakerIds { get; init; } = [];
    public List<FormKey> FactionIds { get; init; } = [];
    public List<FormKey> RaceIds { get; init; } = [];
    public List<FormKey> VoiceTypeIds { get; init; } = [];
    // 前後関係（r7）: PNAM（単一）/ TCLT（複数）。
    public FormKey? PreviousId { get; init; }
    public List<FormKey> NextIds { get; init; } = [];
    public List<ConditionEntry> Conditions { get; init; } = [];
}

public sealed record ResponseLine(int Number, string Text, string Emotion);

public sealed record ConditionEntry(string Function, string Operator, string Value, FormKey? RefId);
