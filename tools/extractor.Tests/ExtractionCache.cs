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

    // 実データの場所を指定する .env のキー。プロセスの環境変数は読まず、repo root の .env だけを見る。
    public const string DataDirKey = "AITRANSLATIONENGINEJP_SKYRIM_DATA_DIR";

    public static string DataFolder => ResolveDataFolder(RepoRoot);

    // 実データの場所を決める。.env のキーが実在する場所を指していればそれを使い、そうでなければ既定へ落ちる。
    internal static string ResolveDataFolder(string repoRoot)
    {
        var configured = ReadEnvValue(Path.Combine(repoRoot, ".env"), DataDirKey);
        if (configured.Length > 0 && Directory.Exists(configured)) return configured;
        return Path.Combine(repoRoot, "dictionaries", "Data");
    }

    // .env から 1 つのキーの値を読む。dotnet test は scripts/dev/run-wails.sh を経由せず .env が環境変数に
    // ならないため、テスト側でファイルを直接読む。読むのは KEY=VALUE 行だけで、# で始まる行は注記として飛ばす。
    // 値の前後のクォートは剥がす（実データの場所は空白を含むことがあり、.env ではクォートで囲むため）。
    private static string ReadEnvValue(string envPath, string key)
    {
        if (!File.Exists(envPath)) return "";

        foreach (var raw in File.ReadAllLines(envPath))
        {
            var line = raw.Trim();
            if (line.Length == 0 || line.StartsWith('#')) continue;

            var sep = line.IndexOf('=');
            if (sep <= 0) continue;
            if (line[..sep].Trim() != key) continue;

            return line[(sep + 1)..].Trim().Trim('"', '\'');
        }
        return "";
    }
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
