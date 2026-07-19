using System.Text.Json;

namespace Extractor.Tests;

// 共有オラクル（test-oracle/specs.json）を C# 側で読む。Go 側と同じ 1 ファイルを parse し、id を join キーにする。
public static class OracleSpecs
{
    // 1 spec = 段 × 属性の判定基準。coverage は例外時のみ（既存単体へ委任した印）。
    public sealed record Spec(string Id, string Stage, string Attribute, string Category, string Given, string SpecText, string? Coverage);

    public static List<Spec> Load()
    {
        var path = Path.Combine(OraclePaths.RepoRoot(), "test-oracle", "specs.json");
        using var doc = JsonDocument.Parse(File.ReadAllText(path));
        var specs = new List<Spec>();
        foreach (var e in doc.RootElement.GetProperty("specs").EnumerateArray())
        {
            specs.Add(new Spec(
                e.GetProperty("id").GetString()!,
                e.GetProperty("stage").GetString()!,
                e.GetProperty("attribute").GetString()!,
                e.GetProperty("category").GetString()!,
                e.GetProperty("given").GetString()!,
                e.GetProperty("spec").GetString()!,
                e.TryGetProperty("coverage", out var c) ? c.GetString() : null));
        }
        return specs;
    }
}

// リポジトリ root を go.mod の位置で探す（test の実行 dir から遡る）。
public static class OraclePaths
{
    public static string RepoRoot()
    {
        var dir = new DirectoryInfo(AppContext.BaseDirectory);
        while (dir != null && !File.Exists(Path.Combine(dir.FullName, "go.mod")))
            dir = dir.Parent;
        if (dir == null) throw new DirectoryNotFoundException("go.mod を含む repo root が見つからない");
        return dir.FullName;
    }
}
