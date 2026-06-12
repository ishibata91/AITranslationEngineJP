using System.Xml;

namespace Extractor;

// xTranslator 辞書 XML（SSTXMLRessources）を REC:FIELD 件数へ集計する
// （scripts/compare_counts.py の count_xml + 非抽出対象除外の移植）。
public static class XTranslatorXml
{
    // record 種別ごと除外（内部名・非表示・ゲーム設定など。docs/mutagen-migration-plan.md §7）。
    public static readonly HashSet<string> NonTargetRecBases =
    [
        "AVIF", // actor value（内部）
        "BPTD", // body part data（内部）
        "CLAS", // class（話者推定の手掛かりに使うが翻訳しない）
        "CLFM", // color form（内部）
        "EXPL", // explosion（非表示）
        "GMST", // game setting（内部）
        "HDPT", // head part（内部名）
        "PROJ", // projectile（非表示）
        "REFR", // 配置参照の name override（保留）
        "COLL", // collision layer（内部）
        "WATR", // water（内部）
        "****", // script fragment 等、REC を持たない行
    ];

    // REC は対象だが特定 FIELD は非対象。
    public static readonly HashSet<string> NonTargetRecFields =
    [
        "PERK:EPF2", // perk entry point button label（特殊）
        "PERK:EPFD", // perk entry point description fragment（特殊）
        "WOOP:FULL", // 力の言葉の龍語綴り（Source=Dest 固定、翻訳禁止）
    ];

    public static bool IsNonTarget(string recField)
    {
        var basePart = recField.Split(':', 2)[0];
        return NonTargetRecBases.Contains(basePart) || NonTargetRecFields.Contains(recField);
    }

    // String 行を REC:FIELD で数える。Source（翻訳元）が空の行は翻訳対象テキストが無いので
    // 数えない（xTranslator は空 field も空行として登録するため）。
    public static SortedDictionary<string, int> CountAll(string xmlPath)
    {
        var counts = new SortedDictionary<string, int>(StringComparer.Ordinal);
        using var reader = XmlReader.Create(xmlPath, new XmlReaderSettings { DtdProcessing = DtdProcessing.Ignore });

        string rec = "", field = "", source = "";
        var inString = false;

        while (reader.Read())
        {
            if (reader.NodeType == XmlNodeType.Element && reader.Name == "String")
            {
                inString = true;
                rec = field = source = "";
            }
            else if (inString && reader.NodeType == XmlNodeType.Element)
            {
                switch (reader.Name)
                {
                    case "REC": rec = reader.ReadElementContentAsString().Trim().ToUpperInvariant(); break;
                    case "FIELD": field = reader.ReadElementContentAsString().Trim().ToUpperInvariant(); break;
                    case "Source": source = reader.ReadElementContentAsString(); break;
                }
            }
            else if (inString && reader.NodeType == XmlNodeType.EndElement && reader.Name == "String")
            {
                inString = false;
                if (source.Trim().Length == 0) continue;
                var rf = rec.Contains(':') ? rec
                    : rec.Length > 0 && field.Length > 0 ? $"{rec}:{field}"
                    : rec.Length > 0 ? rec : field;
                if (rf.Length == 0) continue;
                counts[rf] = counts.GetValueOrDefault(rf) + 1;
            }
        }
        return counts;
    }

    // 非抽出対象を除外した比較対象だけ返す。
    public static SortedDictionary<string, int> CountInScope(string xmlPath) => InScope(CountAll(xmlPath));

    // CountAll の結果から非抽出対象を除外する（全件と比較対象の両方が要る時は
    // CountAll → InScope の順で呼び、XML の再 parse を避ける）。
    public static SortedDictionary<string, int> InScope(SortedDictionary<string, int> all)
    {
        var inScope = new SortedDictionary<string, int>(StringComparer.Ordinal);
        foreach (var (rf, n) in all)
            if (!IsNonTarget(rf))
                inScope[rf] = n;
        return inScope;
    }
}
