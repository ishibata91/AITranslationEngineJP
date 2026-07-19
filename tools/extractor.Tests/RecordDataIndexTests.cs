using Extractor;
using Mutagen.Bethesda;
using Mutagen.Bethesda.Plugins;
using Mutagen.Bethesda.Plugins.Records;
using Mutagen.Bethesda.Strings;
using Xunit;

namespace Extractor.Tests;

// RecordDataIndex.Normalize の memo（初回計算・以降 cache 返却）を固定する。
// OwnsRecord は record 件数 × master 数だけ Normalize を呼ぶため、同じ FormKey の再計算を避ける。
// memo は最適化であり、返す正規化 data（抽出所有判定の入力）を変えてはならない。
public class RecordDataIndexTests
{
    // 存在する record は 2 回目も同一の正規化 data を返し、かつ同一参照（＝再計算していない）。
    [Fact]
    public void Normalizeは同じキーで同一結果を返しmemoする()
    {
        // Arrange: 合成 esm の自 mod index と、その mod が新規所有する record の FormKey を用意する。
        using var env = OracleInput.LoadEnv();
        var mod = env.TargetMod;
        var key = mod.EnumerateMajorRecords()
            .First(r => r.FormKey.ModKey == mod.ModKey)
            .FormKey;
        var index = RecordDataIndex.Build(env.DataFolder, mod.ModKey, Language.English);

        // Act: 同じ key を 2 回正規化する。
        var first = index.Normalize(key);
        var second = index.Normalize(key);

        // Assert: 非 null で byte 一致し、memo により同一配列参照を返す。
        Assert.NotNull(first);
        Assert.Equal(first, second);
        Assert.Same(first, second);
    }

    // 存在しない record は両呼びとも null（null 結果も memo する分岐を固定する）。
    [Fact]
    public void Normalizeは存在しないキーで両呼びともnullを返す()
    {
        // Arrange: 自 mod に存在しない FormID の FormKey を作る。
        using var env = OracleInput.LoadEnv();
        var mod = env.TargetMod;
        var index = RecordDataIndex.Build(env.DataFolder, mod.ModKey, Language.English);
        var missing = new FormKey(mod.ModKey, 0xFFFFFF);

        // Act & Assert: 初回・2 回目とも null。
        Assert.Null(index.Normalize(missing));
        Assert.Null(index.Normalize(missing));
    }
}
