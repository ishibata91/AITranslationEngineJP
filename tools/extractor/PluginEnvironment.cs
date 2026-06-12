using Mutagen.Bethesda;
using Mutagen.Bethesda.Plugins;
using Mutagen.Bethesda.Plugins.Binary.Parameters;
using Mutagen.Bethesda.Plugins.Cache;
using Mutagen.Bethesda.Plugins.Records;
using Mutagen.Bethesda.Skyrim;
using Mutagen.Bethesda.Strings;
using Mutagen.Bethesda.Strings.DI;

namespace Extractor;

// data folder（Skyrim の Data 相当）から target plugin と master 連鎖を読み込む。
// macOS では GameEnvironment.Typical の自動検出（registry / plugins.txt）が使えないため、
// master 連鎖を plugin header から再帰解決して load order を明示構築する。
public sealed class PluginEnvironment : IDisposable
{
    private readonly List<ISkyrimModDisposableGetter> _mods = [];
    private readonly Dictionary<ModKey, RecordDataIndex> _dataIndexes;

    public required string DataFolder { get; init; }
    public required ISkyrimModGetter TargetMod { get; init; }
    // master が先、target が末尾（load order 相当の優先順）。
    public required IReadOnlyList<ISkyrimModGetter> LoadOrder { get; init; }
    public required ILinkCache<ISkyrimMod, ISkyrimModGetter> LinkCache { get; init; }

    private PluginEnvironment(Dictionary<ModKey, RecordDataIndex> dataIndexes)
        => _dataIndexes = dataIndexes;

    // record がこの plugin の翻訳対象か（翻訳所有）。
    // 新規 record、または正規化 data（VMAD 除外 + string ID を本文解決）が master の
    // winning 版と異なる override が対象。一致する override は参照用 stub であり、
    // xTranslator も辞書の record に紐づけない（**** 行 = 非対象になる）。
    public bool OwnsRecord(IMajorRecordGetter rec)
    {
        if (rec.FormKey.ModKey == TargetMod.ModKey) return true; // 新規 record

        var own = _dataIndexes[TargetMod.ModKey].Normalize(rec.FormKey);
        if (own == null) return true;

        // どれかの master 版と一致すれば stub（HearthFires が Update と同一の INFO を
        // 運ぶ場合など、winning 版に限らず一致した時点で「変更なし」とみなす）。
        for (var i = LoadOrder.Count - 2; i >= 0; i--)
        {
            var master = _dataIndexes[LoadOrder[i].ModKey].Normalize(rec.FormKey);
            if (master != null && own.AsSpan().SequenceEqual(master)) return false;
        }
        return true; // どの master 版とも一致しない（または master に無い）→ 所有
    }

    // 抽出（翻訳元テキストの解決）は英語を既定にする。翻訳エンジンの入力は英語 mod であり、
    // xTranslator 辞書の Source（照合対象）も英語のため。
    public static PluginEnvironment Load(string dataFolder, string pluginName, Language language = Language.English)
    {
        var readParams = new BinaryReadParameters
        {
            StringsParam = new StringsReadParameters
            {
                TargetLanguage = language,
                // 非 localized plugin の内蔵文字列は UTF-8 優先で decode する
                //（翻訳適用済み esp は日本語 UTF-8 を内蔵する。fallback は cp1252）。
                NonLocalizedEncodingOverride = MutagenEncoding._utf8_1252,
            },
        };

        // master 連鎖を深さ優先で解決する。同名 plugin は 1 回だけ読む。
        var loaded = new Dictionary<ModKey, ISkyrimModDisposableGetter>();
        var order = new List<ISkyrimModGetter>();

        ISkyrimModGetter? LoadRecursive(ModKey key, bool required)
        {
            if (loaded.TryGetValue(key, out var existing)) return existing;
            var path = Path.Combine(dataFolder, key.FileName);
            if (!File.Exists(path))
            {
                // master 欠落は許容する（参照は dangling、override の比較先なしは所有扱い）。
                if (!required) return null;
                throw new FileNotFoundException($"plugin が data folder に無い: {key.FileName}", path);
            }
            var mod = SkyrimMod.CreateFromBinaryOverlay(new ModPath(key, path), SkyrimRelease.SkyrimSE, readParams);
            loaded[key] = mod;
            foreach (var master in mod.ModHeader.MasterReferences)
                LoadRecursive(master.Master, required: false);
            order.Add(mod); // master を先、自身を後に積む（load order 相当の優先順）。
            return mod;
        }

        var targetKey = ModKey.FromFileName(pluginName);
        var target = LoadRecursive(targetKey, required: true)!;
        var linkCache = order.ToImmutableLinkCache<ISkyrimMod, ISkyrimModGetter>();

        // 翻訳所有判定用に、全 plugin の record data index（正規化比較）を作る。
        var dataIndexes = order.ToDictionary(
            m => m.ModKey,
            m => RecordDataIndex.Build(dataFolder, m.ModKey, language));

        var env = new PluginEnvironment(dataIndexes)
        {
            DataFolder = dataFolder,
            TargetMod = target,
            LoadOrder = order,
            LinkCache = linkCache,
        };
        env._mods.AddRange(loaded.Values);
        return env;
    }

    public void Dispose()
    {
        foreach (var mod in _mods) mod.Dispose();
        _mods.Clear();
    }
}
