using Xunit;

namespace Extractor.Tests;

// 実データ（Skyrim の Data フォルダ）が要るテストに付ける。
// 実データが無い機械では失敗ではなく skip として記録し、環境の未整備と抽出の不具合を実行結果で区別できるようにする。
// xunit 2.9 は動的 skip（Assert.Skip 系）を持たないため、属性の生成時に実在を見て Skip の理由を決める。
public sealed class RealDataFactAttribute : FactAttribute
{
    public RealDataFactAttribute()
    {
        var dataFolder = TestPaths.DataFolder;
        if (!Directory.Exists(dataFolder))
            Skip = $"実データが無いので skip する。repo root の .env に {TestPaths.DataDirKey} を書くと実行する（探した場所: {dataFolder}）";
    }
}
