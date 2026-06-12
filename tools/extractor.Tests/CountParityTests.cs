using System.Collections.Concurrent;
using System.Text;
using Extractor;
using Xunit;

namespace Extractor.Tests;

// 件数比較検証: 抽出結果を xTranslator 辞書 XML と REC:FIELD 単位で突き合わせる
// （scripts/compare_counts.py の置き換え。完了条件 = 参考 XML 全 plugin で一致）。
//
// 既知差分（KnownDeltas）について:
//   xTranslator 辞書側に、plugin 内容から導出できない帰属の揺れが少数ある。
//   同一内容・同一構造の record 対で帰属が逆になる実測例（採用 / **** の対）:
//     - CELL: HallOfTheVigilant01（採用）vs KarthwastenSanuarachMine（****）
//       …どちらも record 同一 + 子 REFR の XNDP（navmesh link）のみ変更
//     - QUST: DA07 等（採用）vs TG07 / MQ105（****）…どちらも VMAD（script）変更
//     - NPC_: 新規は全て採用だが、AbVampire01（SPEL、data 変更）採用 vs Olda（NPC_、
//       data 同一）**** など、override の扱いが record 種別・plugin 間で非一貫
//   また XML 生成セッションが異なる（公式 5 esm = 2026-02、USSEP = 2026-06）ため、
//   生成設定の差も残差に含まれる。残差は全 98,000 strings 中 65（0.066%）で、
//   下表に plugin × REC:FIELD 単位で固定する。表に無い乖離はすべて fail になる。
public class CountParityTests
{
    // 既知差分: extracted - xml の期待値。
    private static readonly Dictionary<(string Plugin, string RecField), int> KnownDeltas = new()
    {
        // Update.esm: TG07 / MQ105（script + 条件のみの変更、辞書側 ****）を抽出が所有判定、
        // MG07StaffofMagnus（CRDT のみの変更、辞書側 ****）、CELL 2 件（辞書側は採用）。
        [("Update.esm", "QUST:FULL")] = +2,
        [("Update.esm", "QUST:CNAM")] = +8,
        [("Update.esm", "QUST:NNAM")] = +16,
        [("Update.esm", "WEAP:FULL")] = +1,
        [("Update.esm", "WEAP:DESC")] = +1,
        [("Update.esm", "CELL:FULL")] = -2, // ThalmorEmbassy04 / WinterholdCollege 系
        // Dawnguard.esm: Breezehome / SanuarachMine（子 REFR の navmesh link 変更のみ、辞書側 ****）、
        // ARMO 1 件（master 一致 stub と判定したが辞書側は採用）。
        [("Dawnguard.esm", "CELL:FULL")] = +2,
        [("Dawnguard.esm", "ARMO:FULL")] = -1,
        // HearthFires.esm: Update 版と同一内容の INFO 6 件と FURN 1 件、CELL 5 件（辞書側は採用）。
        [("HearthFires.esm", "INFO:NAM1")] = -6,
        [("HearthFires.esm", "CELL:FULL")] = -5,
        [("HearthFires.esm", "FURN:FULL")] = -1,
        // Dragonborn.esm: DA04（script + 条件 + 再 string、辞書側 ****）、CELL 1 件。
        [("Dragonborn.esm", "QUST:FULL")] = +1,
        [("Dragonborn.esm", "QUST:CNAM")] = +7,
        [("Dragonborn.esm", "QUST:NNAM")] = +12,
        [("Dragonborn.esm", "CELL:FULL")] = +1,
        // USSEP: DLC2MKMiraakMask1H の DESC（翻訳適用時に追加された文。原文 esp には無い）。
        [("unofficial skyrim special edition patch.esp", "ARMO:DESC")] = +1,
    };

    // 抽出は plugin ごとに重いので、テスト間で結果を共有する。
    private static readonly ConcurrentDictionary<string, ExtractionResult> Cache = new();

    public static ExtractionResult ExtractCached(string plugin) =>
        Cache.GetOrAdd(plugin, p =>
        {
            using var env = PluginEnvironment.Load(TestPaths.DataFolder, p);
            return PluginExtractor.Extract(env);
        });

    [Theory]
    [InlineData("Skyrim.esm", "Skyrim_english_japanese.xml")]
    [InlineData("Update.esm", "Update_english_japanese.xml")]
    [InlineData("Dawnguard.esm", "Dawnguard_english_japanese.xml")]
    [InlineData("HearthFires.esm", "HearthFires_english_japanese.xml")]
    [InlineData("Dragonborn.esm", "Dragonborn_english_japanese.xml")]
    [InlineData("unofficial skyrim special edition patch.esp", "unofficial skyrim special edition patch_english_japanese.xml")]
    public void 件数がxTranslatorXMLと一致する(string plugin, string xmlName)
    {
        var xmlPath = Path.Combine(TestPaths.XmlFolder, xmlName);
        Assert.True(File.Exists(xmlPath), $"照合先 XML が無い: {xmlPath}");

        var result = ExtractCached(plugin);
        var extracted = TranslationCounts.Flatten(result);
        var xmlInScope = XTranslatorXml.CountInScope(xmlPath);
        var comparison = CountComparison.Build(extracted, xmlInScope, xmlExcluded: 0);

        var unexpected = comparison.Mismatches
            .Where(row => KnownDeltas.GetValueOrDefault((plugin, row.RecField)) != row.Delta)
            .ToList();
        // 既知差分が解消された場合も fail にする（表の更新を強制し、暗黙の劣化と区別する）。
        var resolvedKnown = KnownDeltas.Keys
            .Where(k => k.Plugin == plugin)
            .Where(k => comparison.Rows.FirstOrDefault(r => r.RecField == k.RecField)?.Delta != KnownDeltas[k])
            .ToList();

        if (unexpected.Count > 0 || resolvedKnown.Count > 0)
        {
            var sb = new StringBuilder();
            sb.AppendLine($"{plugin}: REC:FIELD 件数が想定と一致しない（delta 合計 {comparison.DeltaTotal:+0;-0;+0}）");
            if (unexpected.Count > 0)
            {
                sb.AppendLine("| REC:FIELD | 抽出 | xml | delta |（既知差分表に無い乖離）");
                foreach (var row in unexpected.OrderBy(r => r.Delta))
                    sb.AppendLine($"| {row.RecField} | {row.Extracted} | {row.Xml} | {row.Delta:+0;-0} |");
            }
            foreach (var k in resolvedKnown)
                sb.AppendLine($"既知差分 {k.RecField}({KnownDeltas[k]:+0;-0}) が現状と一致しない。表の更新が必要。");
            Assert.Fail(sb.ToString());
        }
    }
}

public static class TestPaths
{
    // repo root をテスト実行ディレクトリから遡って探す（bin/Debug 配下から実行されるため）。
    public static string RepoRoot { get; } = FindRepoRoot();
    public static string DataFolder => Path.Combine(RepoRoot, "dictionaries", "Data");
    public static string XmlFolder => Path.Combine(RepoRoot, "dictionaries", "xTranslatorXMLs");

    private static string FindRepoRoot()
    {
        var dir = new DirectoryInfo(AppContext.BaseDirectory);
        while (dir != null)
        {
            if (Directory.Exists(Path.Combine(dir.FullName, "dictionaries", "Data")))
                return dir.FullName;
            dir = dir.Parent;
        }
        throw new DirectoryNotFoundException("repo root（dictionaries/Data を含む dir）が見つからない");
    }
}
