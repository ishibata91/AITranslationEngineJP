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

string? dataFolder = null, plugin = null, xmlPath = null, dumpRecField = null, sqlitePath = null, schemaDir = null, termsXmlDir = null;
var language = Language.English;

for (var i = 0; i < args.Length; i++)
{
    switch (args[i])
    {
        case "--data": dataFolder = Next(ref i) ; break;
        case "--plugin": plugin = Next(ref i); break;
        case "--xml": xmlPath = Next(ref i); break;
        case "--dump": dumpRecField = Next(ref i).ToUpperInvariant(); break;
        case "--sqlite": sqlitePath = Next(ref i); break;
        case "--schema": schemaDir = Next(ref i); break;
        case "--terms-xml": termsXmlDir = Next(ref i); break;
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

if (sqlitePath != null)
{
    // 中心 DB（SQLite）へ翻訳対象を書く。schema は db/migrations の SQL を ensure する（C#↔Go 契約 1 本）。
    var dir = schemaDir ?? FindSchemaDir();
    var schemaSql = string.Join("\n",
        Directory.GetFiles(dir, "*.sql").OrderBy(p => p, StringComparer.Ordinal).Select(File.ReadAllText));
    // 全 REC:FIELD の原文を extracted_field へ素朴に書く（箱・directive の判定は持たない）。
    // 箱の振り分けは Go 取込段が record_type_master を引いて narration/proper_noun/line へ行う。
    var written = ExtractedFieldSqliteWriter.Write(sqlitePath, schemaSql, result);
    Console.WriteLine($"[sqlite] extracted_field {written} 件を {sqlitePath} へ書き込み（{sw.ElapsedMilliseconds} ms）");

    // 台詞の話者属性と INFO→speaker の橋渡し。話者は master 側 NPC を LinkCache で解決して書く。
    // 台詞本文は extracted_field 側に入り、line 行と line_speaker は Go 取込段が作る。
    var linkCount = SpeakerSqliteWriter.Write(sqlitePath, schemaSql, result, env.LinkCache);
    Console.WriteLine($"[sqlite] 話者属性と INFO→speaker 橋渡し {linkCount} 件を {sqlitePath} へ書き込み（{sw.ElapsedMilliseconds} ms）");

    // INFO の条件から導いた性別（汎用台詞の一人称・語尾の根拠）を書く。話者解決の有無に依らず性別を定められる INFO を書き、
    // line 行は Go 取込段が作るため安定キー（INFO の plugin・form_id）で一時保持する（line_condition へ Go 取込段が解決する）。
    var condCount = InfoConditionSqliteWriter.Write(sqlitePath, schemaSql, env, result.TargetPlugin);
    Console.WriteLine($"[sqlite] INFO→条件由来の性別 {condCount} 件を {sqlitePath} へ書き込み（{sw.ElapsedMilliseconds} ms）");

    // T3 の固有名辞書（master_term）。xTranslator 英日辞書 XML（公式日本語版既訳）から FULL を読み、
    // 叙述文・台詞の本文へ機械置換するための「原語 → 確定訳語」を書く。
    // 既定は data folder の兄弟 xTranslatorXMLs、--terms-xml で上書きできる。
    var xmlDir = termsXmlDir ?? Path.Combine(dataFolder, "..", "xTranslatorXMLs");
    if (Directory.Exists(xmlDir))
    {
        var termXmls = Directory.GetFiles(xmlDir, "*.xml");
        var termCount = MasterTermXmlWriter.Write(sqlitePath, schemaSql, termXmls);
        Console.WriteLine($"[sqlite] master_term {termCount} 件を {termXmls.Length} 個の XML から {sqlitePath} へ書き込み（{sw.ElapsedMilliseconds} ms）");
    }
    else
    {
        Console.WriteLine($"[sqlite] 固有名辞書の XML ディレクトリが無い（{xmlDir}）。master_term は作らない。");
    }
    return 0;
}

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

// db/migrations を、作業ディレクトリから repo root（go.mod がある所）まで遡って探す。
static string FindSchemaDir()
{
    var dir = new DirectoryInfo(Directory.GetCurrentDirectory());
    while (dir != null && !File.Exists(Path.Combine(dir.FullName, "go.mod")))
        dir = dir.Parent;
    if (dir == null) throw new InvalidOperationException("repo root（go.mod）が見つからない。--schema で migrations ディレクトリを指定してください。");
    return Path.Combine(dir.FullName, "db", "migrations");
}

static void Print(string label, int count) => Console.WriteLine($"- {label}: {count}");

static void PrintUsage()
{
    Console.WriteLine("""
        使い方: extractor --data <DataFolder> --plugin <name.esp> [--sqlite <db>] [--xml <xTranslator XML>] [--language Japanese]
          --data      Skyrim の Data 相当フォルダ（esm/esp/esl + Strings/）
          --plugin    抽出対象 plugin ファイル名（master は同フォルダから自動解決）
          --sqlite    中心 DB（SQLite）へ全 REC:FIELD の原文（extracted_field）と話者属性を書き込む。書込後に終了する
          --schema    --sqlite 用の migrations ディレクトリ（既定 repo の db/migrations を自動探索）
          --terms-xml 固有名辞書（master_term）の供給元 xTranslator 英日 XML ディレクトリ（既定 data の兄弟 xTranslatorXMLs）
          --xml       xTranslator 辞書 XML と REC:FIELD 件数比較を行う（不一致なら exit 2）
          --language  localized strings の言語（既定 English。翻訳元テキストを解決する）
        """);
}
