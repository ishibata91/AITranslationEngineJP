using Xunit;

namespace Extractor.Tests;

// 実データの場所を決める解決順の検証。実データそのものは要らず、一時フォルダを repo root に見立てて確かめる。
// 解決順は「repo root の .env のキーが実在する場所を指していればそれ、そうでなければ既定 <repo>/dictionaries/Data」。
public class TestPathsTests
{
    // Arrange 用。repo root に見立てた一時フォルダを作り、消えるまでを 1 つにまとめる。
    private sealed class TempRepoRoot : IDisposable
    {
        public string Path { get; }

        public TempRepoRoot()
        {
            Path = System.IO.Path.Combine(System.IO.Path.GetTempPath(), $"repo-{Guid.NewGuid():N}");
            Directory.CreateDirectory(Path);
        }

        public void WriteEnv(string content) => File.WriteAllText(System.IO.Path.Combine(Path, ".env"), content);

        // 実在する Data フォルダを作って、その絶対パスを返す。
        public string CreateDataFolder(string name)
        {
            var dir = System.IO.Path.Combine(Path, name);
            Directory.CreateDirectory(dir);
            return dir;
        }

        public string DefaultDataFolder => System.IO.Path.Combine(Path, "dictionaries", "Data");

        public void Dispose()
        {
            if (Directory.Exists(Path)) Directory.Delete(Path, recursive: true);
        }
    }

    [Fact]
    public void env_のキーが実在する場所を指していればその場所を使う()
    {
        // Arrange
        using var root = new TempRepoRoot();
        var real = root.CreateDataFolder("RealData");
        root.WriteEnv($"{TestPaths.DataDirKey}={real}\n");

        // Act
        var got = TestPaths.ResolveDataFolder(root.Path);

        // Assert
        Assert.Equal(real, got);
    }

    [Fact]
    public void env_が無ければ既定へ落ちる()
    {
        // Arrange
        using var root = new TempRepoRoot();

        // Act
        var got = TestPaths.ResolveDataFolder(root.Path);

        // Assert
        Assert.Equal(root.DefaultDataFolder, got);
    }

    [Fact]
    public void env_に対象のキーが無ければ既定へ落ちる()
    {
        // Arrange
        using var root = new TempRepoRoot();
        root.WriteEnv("AITRANSLATIONENGINEJP_TEST_MODE=false\n");

        // Act
        var got = TestPaths.ResolveDataFolder(root.Path);

        // Assert
        Assert.Equal(root.DefaultDataFolder, got);
    }

    [Fact]
    public void env_のキーが実在しない場所を指していれば既定へ落ちる()
    {
        // Arrange
        using var root = new TempRepoRoot();
        root.WriteEnv($"{TestPaths.DataDirKey}={System.IO.Path.Combine(root.Path, "NotExist")}\n");

        // Act
        var got = TestPaths.ResolveDataFolder(root.Path);

        // Assert
        Assert.Equal(root.DefaultDataFolder, got);
    }

    [Fact]
    public void 値の前後のクォートを剥がす()
    {
        // Arrange: 実データの場所は空白を含むことがあり（Skyrim Special Edition）、.env ではクォートで囲む。
        using var root = new TempRepoRoot();
        var real = root.CreateDataFolder("Real Data With Space");
        root.WriteEnv($"{TestPaths.DataDirKey}=\"{real}\"\n");

        // Act
        var got = TestPaths.ResolveDataFolder(root.Path);

        // Assert
        Assert.Equal(real, got);
    }

    [Fact]
    public void 注記行は読まない()
    {
        // Arrange: # で始まる行は設定ではない。既定へ落ちること。
        using var root = new TempRepoRoot();
        var real = root.CreateDataFolder("RealData");
        root.WriteEnv($"# {TestPaths.DataDirKey}={real}\n");

        // Act
        var got = TestPaths.ResolveDataFolder(root.Path);

        // Assert
        Assert.Equal(root.DefaultDataFolder, got);
    }
}
