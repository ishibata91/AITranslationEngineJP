using System.Diagnostics;
using Mutagen.Bethesda;
using Mutagen.Bethesda.Environments;
using Mutagen.Bethesda.Skyrim;

// 素振り: Mutagen で Skyrim SE の load order を読み、主要 record の件数と速度を測る。
// xEdit Pascal（Skyrim.esm で 18 分）との比較が目的。
try
{
    var sw = Stopwatch.StartNew();
    using var env = GameEnvironment.Typical.Skyrim(SkyrimRelease.SkyrimSE);
    Console.WriteLine($"[env] {env.LoadOrder.Count} plugins loaded in {sw.ElapsedMilliseconds} ms");
    foreach (var listing in env.LoadOrder.ListedOrder)
        Console.WriteLine($"   - {listing.ModKey} (enabled={listing.Enabled})");

    var books = env.LoadOrder.PriorityOrder.Book().WinningOverrides()
        .Count(b => !string.IsNullOrEmpty(b.Name?.String));
    Console.WriteLine($"[BOOK] name非空 {books}  ({sw.ElapsedMilliseconds} ms)");

    var npcs = env.LoadOrder.PriorityOrder.Npc().WinningOverrides()
        .Count(n => !string.IsNullOrEmpty(n.Name?.String));
    Console.WriteLine($"[NPC_] name非空 {npcs}  ({sw.ElapsedMilliseconds} ms)");

    var weaps = env.LoadOrder.PriorityOrder.Weapon().WinningOverrides()
        .Count(w => !string.IsNullOrEmpty(w.Name?.String));
    Console.WriteLine($"[WEAP] name非空 {weaps}  ({sw.ElapsedMilliseconds} ms)");

    Console.WriteLine($"[done] total {sw.ElapsedMilliseconds} ms");
}
catch (Exception e)
{
    Console.WriteLine($"[error] {e.GetType().Name}: {e.Message}");
    Console.WriteLine("自動検出に失敗した場合は、load order を明示構築する経路へ切り替える。");
}
