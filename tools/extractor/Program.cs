using System.Diagnostics;
using Extractor;
using Mutagen.Bethesda.Strings;

// Mutagen 抽出 CLI。data folder と plugin を直指定して抽出し、カテゴリ別件数を表示する。
// --xml を渡すと xTranslator 辞書 XML との REC:FIELD 件数比較も行う（不一致なら exit 2）。
//
// 使い方:
//   dotnet run --project tools/extractor -- --data dictionaries/Data --plugin Dawnguard.esm
//   dotnet run --project tools/extractor -- --data dictionaries/Data --plugin Dawnguard.esm \
//       --xml dictionaries/xTranslatorXMLs/Dawnguard_english_japanese.xml

string? dataFolder = null, plugin = null, xmlPath = null, dumpRecField = null;
var language = Language.English;

for (var i = 0; i < args.Length; i++)
{
    switch (args[i])
    {
        case "--data": dataFolder = Next(ref i) ; break;
        case "--plugin": plugin = Next(ref i); break;
        case "--xml": xmlPath = Next(ref i); break;
        case "--dump": dumpRecField = Next(ref i).ToUpperInvariant(); break;
        case "--language": language = Enum.Parse<Language>(Next(ref i), ignoreCase: true); break;
        case "--help" or "-h": PrintUsage(); return 0;
        default:
            Console.Error.WriteLine($"不明な引数: {args[i]}");
            PrintUsage();
            return 1;
    }
}

if (dataFolder == null || plugin == null)
{
    Console.Error.WriteLine("--data と --plugin は必須。");
    PrintUsage();
    return 1;
}

var sw = Stopwatch.StartNew();
using var env = PluginEnvironment.Load(dataFolder, plugin, language);
Console.WriteLine($"[load] {env.LoadOrder.Count} plugins（master 連鎖込み）を {sw.ElapsedMilliseconds} ms で読み込み");

var result = PluginExtractor.Extract(env);
Console.WriteLine($"[extract] {plugin} を {sw.ElapsedMilliseconds} ms で抽出");

if (dumpRecField != null)
{
    // 調査用: 指定 REC:FIELD の翻訳文字列を tab 区切りで列挙する（EDID / FormKey / text 先頭 60 字）。
    foreach (var s in TranslationCounts.Enumerate(result).Where(s => s.RecField == dumpRecField))
        Console.WriteLine($"{s.EditorId}\t{s.Id}\t{s.Text[..Math.Min(60, s.Text.Length)].ReplaceLineEndings(" ")}");
    return 0;
}
Console.WriteLine();
Console.WriteLine("## カテゴリ別件数");
Print("dialogues (DIAL)", result.Dialogues.Count);
Print("  info nodes (INFO)", result.Dialogues.Sum(d => d.Infos.Count));
Print("  response lines (NAM1)", result.Dialogues.Sum(d => d.Infos.Sum(i => i.Responses.Count)));
Print("quests", result.Quests.Count);
Print("races", result.Races.Count);
Print("factions", result.Factions.Count);
Print("items", result.Items.Count);
Print("activators", result.Activators.Count);
Print("equipment", result.Equipment.Count);
Print("consumables", result.Consumables.Count);
Print("magic (SPEL)", result.Magic.Count);
Print("enchantments (ENCH)", result.Enchantments.Count);
Print("magic effects (MGEF)", result.MagicEffects.Count);
Print("shouts (SHOU)", result.Shouts.Count);
Print("words of power (WOOP)", result.Words.Count);
Print("books", result.Books.Count);
Print("locations", result.Locations.Count);
Print("messages", result.Messages.Count);
Print("loading screens (LSCR)", result.LoadingScreens.Count);
Print("perks", result.Perks.Count);
Print("speakers (NPC_/TACT)", result.Speakers.Count);
Print("voice types (VTYP)", result.VoiceTypes.Count);
Print("sound categories (SNCT)", result.SoundCategories.Count);
Print("eyes (EYES)", result.Eyes.Count);
Print("regions (REGN)", result.Regions.Count);

if (xmlPath == null) return 0;

var extracted = TranslationCounts.Flatten(result);
var xmlAll = XTranslatorXml.CountAll(xmlPath);
var xmlInScope = XTranslatorXml.InScope(xmlAll);
var excluded = xmlAll.Values.Sum() - xmlInScope.Values.Sum();
var comparison = CountComparison.Build(extracted, xmlInScope, excluded);

Console.WriteLine();
Console.WriteLine($"# 件数比較: {Path.GetFileName(xmlPath)}");
Console.WriteLine(comparison.ToMarkdown());
return comparison.IsMatch ? 0 : 2;

string Next(ref int i)
{
    if (i + 1 >= args.Length) throw new ArgumentException($"{args[i]} に値が無い");
    return args[++i];
}

static void Print(string label, int count) => Console.WriteLine($"- {label}: {count}");

static void PrintUsage()
{
    Console.WriteLine("""
        使い方: extractor --data <DataFolder> --plugin <name.esp> [--xml <xTranslator XML>] [--language Japanese]
          --data      Skyrim の Data 相当フォルダ（esm/esp/esl + Strings/）
          --plugin    抽出対象 plugin ファイル名（master は同フォルダから自動解決）
          --xml       xTranslator 辞書 XML と REC:FIELD 件数比較を行う（不一致なら exit 2）
          --language  localized strings の言語（既定 English。翻訳元テキストを解決する）
        """);
}
