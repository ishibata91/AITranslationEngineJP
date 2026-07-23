using System.Collections.Concurrent;
using Extractor;

namespace Extractor.Tests;

// 実データ（gitignore 配下の dictionaries/Data）を使うテストの共有パス。
public static class TestPaths
{
    public static string RepoRoot
    {
        get
        {
            var dir = new DirectoryInfo(AppContext.BaseDirectory);
            while (dir != null && !File.Exists(Path.Combine(dir.FullName, "go.mod")))
                dir = dir.Parent;
            return dir?.FullName ?? throw new InvalidOperationException("repo root（go.mod）が見つからない");
        }
    }

    public static string DataFolder => Path.Combine(RepoRoot, "dictionaries", "Data");
}

// 抽出は plugin ごとに重いので、テスト間で結果を共有する（読み取り専用の Arrange ヘルパ）。
public static class ExtractionCache
{
    private static readonly ConcurrentDictionary<string, ExtractionResult> Cache = new();

    public static ExtractionResult ExtractCached(string plugin) =>
        Cache.GetOrAdd(plugin, p =>
        {
            using var env = PluginEnvironment.Load(TestPaths.DataFolder, p);
            return PluginExtractor.Extract(env);
        });
}
