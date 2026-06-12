using System.Text;

namespace Extractor;

// 抽出件数と xTranslator XML 件数の REC:FIELD 比較（scripts/compare_counts.py の report 相当）。
public sealed record ComparisonRow(string RecField, int Extracted, int Xml, int Delta);

public sealed class CountComparison
{
    public required IReadOnlyList<ComparisonRow> Rows { get; init; }
    public required int XmlExcluded { get; init; }

    public int DeltaTotal => Rows.Sum(r => r.Delta);
    public bool IsMatch => Rows.All(r => r.Delta == 0);
    public IEnumerable<ComparisonRow> Mismatches => Rows.Where(r => r.Delta != 0);

    public static CountComparison Build(
        IReadOnlyDictionary<string, int> extracted,
        IReadOnlyDictionary<string, int> xmlInScope,
        int xmlExcluded)
    {
        var fields = extracted.Keys.Union(xmlInScope.Keys).OrderBy(k => k, StringComparer.Ordinal);
        var rows = fields
            .Select(rf =>
            {
                var e = extracted.GetValueOrDefault(rf);
                var x = xmlInScope.GetValueOrDefault(rf);
                return new ComparisonRow(rf, e, x, e - x);
            })
            .ToList();
        return new CountComparison { Rows = rows, XmlExcluded = xmlExcluded };
    }

    public string ToMarkdown()
    {
        var sb = new StringBuilder();
        sb.AppendLine("## 合計");
        sb.AppendLine($"- XML 比較対象: {Rows.Sum(r => r.Xml)}（非抽出対象 {XmlExcluded} 行を除外）");
        sb.AppendLine($"- 抽出: {Rows.Sum(r => r.Extracted)}");
        sb.AppendLine($"- 差分: {DeltaTotal:+0;-0;+0}");

        var mismatches = Mismatches.ToList();
        sb.AppendLine();
        sb.AppendLine($"## 不一致（{mismatches.Count} 種）");
        if (mismatches.Count > 0)
        {
            sb.AppendLine("| REC:FIELD | 抽出 | xml | delta |");
            sb.AppendLine("| --- | ---: | ---: | ---: |");
            foreach (var r in mismatches.OrderBy(r => r.Delta))
                sb.AppendLine($"| {r.RecField} | {r.Extracted} | {r.Xml} | {r.Delta:+0;-0} |");
        }
        else
        {
            sb.AppendLine("(なし)");
        }

        var matches = Rows.Where(r => r.Delta == 0 && r.Xml > 0).ToList();
        sb.AppendLine();
        sb.AppendLine($"## 一致（{matches.Count} 種）");
        sb.AppendLine("| REC:FIELD | 件数 |");
        sb.AppendLine("| --- | ---: |");
        foreach (var r in matches)
            sb.AppendLine($"| {r.RecField} | {r.Xml} |");

        return sb.ToString();
    }
}
