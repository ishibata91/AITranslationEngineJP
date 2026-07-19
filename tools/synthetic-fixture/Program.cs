using Mutagen.Bethesda.Plugins.Binary.Parameters;
using Mutagen.Bethesda.Plugins.Records;
using Mutagen.Bethesda.Skyrim;
using SyntheticFixture;

// 合成 fixture 生成スクリプト（独立実行・テストと build から切り離す）。
// 合成 esm を組んで disk へ書く。fixture 変更時だけ手動で回し、成果物を commit する（golden 再生成と同じ運用）。
//
// 使い方: dotnet run --project tools/synthetic-fixture -- --out test-oracle/fixture

var outDir = "test-oracle/fixture";
for (var i = 0; i < args.Length; i++)
{
    if (args[i] == "--out" && i + 1 < args.Length) outDir = args[++i];
}

Directory.CreateDirectory(outDir);
var esmPath = Path.Combine(outDir, SyntheticEsmBuilder.PluginName);

var mod = SyntheticEsmBuilder.Build();
mod.WriteToBinary(esmPath, new BinaryWriteParameters
{
    // master 無しの自己完結 plugin。全 record が自身の ModKey に属するため masters list は空になる。
    ModKey = ModKeyOption.NoCheck,
});

Console.WriteLine($"[synthetic-fixture] 合成 esm を書き出した: {esmPath}");
return 0;
