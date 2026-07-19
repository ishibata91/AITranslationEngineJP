using System.Buffers.Binary;
using System.IO.Compression;
using System.Text;
using Mutagen.Bethesda.Plugins;
using Mutagen.Bethesda.Strings;

namespace Extractor;

// plugin file の record data を FormKey で引き、xTranslator 相当の「正規化 data」を返す index。
//
// 正規化 data = record の subrecord 列から
//   - VMAD（script。並び替えと fragment 名の揺れが翻訳と無関係に起きる）を除き、
//   - 翻訳対象 string field の string ID を strings table の本文に解決した
// バイト列。override の正規化 data が master（winning 版）と一致する record は、
// 「string ID の振り直しや script 再保存だけで、翻訳に関わる内容を何も変えていない stub」
// であり、xTranslator も辞書の record に紐づけない（**** 行 = 非対象になる）。
public sealed class RecordDataIndex
{
    private const uint CompressedFlag = 0x0004_0000;
    private const uint LocalizedFlag = 0x0000_0080;

    private readonly byte[] _data;
    private readonly bool _localized;
    // string ID → 本文（3 table 統合）。言語ごとに保持する。
    private readonly IReadOnlyList<Dictionary<uint, string>> _stringsByLanguage;
    private readonly Dictionary<FormKey, (int Offset, int DataSize, uint Flags)> _records = [];
    // 正規化 data の memo。OwnsRecord は record 件数 × master 数だけ Normalize を呼び、
    // 同じ FormKey を繰り返し再計算する（zlib 展開・2 回 Sort）。結果（null 含む）を初回で覚え、
    // 以降は再計算しない。master 依存が多い mod での再計算コストを抑える（抽出結果は変えない）。
    private readonly Dictionary<FormKey, byte[]?> _normalizedCache = [];

    public bool TryGetRecord(FormKey key, out (int Offset, int DataSize, uint Flags) loc)
        => _records.TryGetValue(key, out loc);

    private RecordDataIndex(byte[] data, bool localized, IReadOnlyList<Dictionary<uint, string>> stringsByLanguage)
    {
        _data = data;
        _localized = localized;
        _stringsByLanguage = stringsByLanguage;
    }

    public static RecordDataIndex Build(string dataFolder, ModKey ownKey, Language language)
    {
        var path = Path.Combine(dataFolder, ownKey.FileName);
        var data = File.ReadAllBytes(path);

        var tes4Flags = BinaryPrimitives.ReadUInt32LittleEndian(data.AsSpan(8));
        var localized = (tes4Flags & LocalizedFlag) != 0;
        // 所有判定はどちらかの言語で本文が変われば所有とする（union）。
        // xTranslator 辞書は英日 pair から生成されており、「英語だけ変わった record」
        //（DA16 等）も「日本語だけ改訳された record」も plugin 側の翻訳対象になる。
        var stringsByLanguage = new List<Dictionary<uint, string>>();
        if (localized)
        {
            // 所有判定の解決は英語（xTranslator 辞書の生成元言語）。翻訳先言語の改訳は
            // 所有を変えない（Update の日本語改訳 quest が **** である実測に基づく）。
            foreach (var lang in new[] { Language.English, language }.Distinct())
            {
                var table = LoadStrings(dataFolder, ownKey, lang);
                if (table.Count > 0) { stringsByLanguage.Add(table); break; }
            }
        }

        // TES4 header の master 列から FormID → FormKey の対応を作る。
        var masters = ReadMasters(data);

        var index = new RecordDataIndex(data, localized, stringsByLanguage);

        FormKey ToFormKey(uint rawFormId)
        {
            var masterIndex = (int)(rawFormId >> 24);
            var modKey = masterIndex < masters.Count ? masters[masterIndex] : ownKey;
            return new FormKey(modKey, rawFormId & 0x00FF_FFFF);
        }

        void Walk(int pos, int end)
        {
            while (pos < end)
            {
                if (data.AsSpan(pos, 4).SequenceEqual("GRUP"u8))
                {
                    var groupSize = (int)BinaryPrimitives.ReadUInt32LittleEndian(data.AsSpan(pos + 4));
                    Walk(pos + 24, pos + groupSize);
                    pos += groupSize;
                    continue;
                }
                var dataSize = (int)BinaryPrimitives.ReadUInt32LittleEndian(data.AsSpan(pos + 4));
                var flags = BinaryPrimitives.ReadUInt32LittleEndian(data.AsSpan(pos + 8));
                var rawFormId = BinaryPrimitives.ReadUInt32LittleEndian(data.AsSpan(pos + 12));
                index._records[ToFormKey(rawFormId)] = (pos, dataSize, flags);
                pos += 24 + dataSize;
            }
        }

        var tes4Size = (int)BinaryPrimitives.ReadUInt32LittleEndian(data.AsSpan(4));
        Walk(24 + tes4Size, data.Length);
        return index;
    }

    // record の正規化 data。record が無い場合は null。結果を memo し、同じ FormKey の再計算を避ける。
    public byte[]? Normalize(FormKey key)
    {
        // 格納済みなら（null でも）そのまま返す。TryGetValue が false の時だけ初回計算する。
        if (_normalizedCache.TryGetValue(key, out var cached)) return cached;
        var result = ComputeNormalized(key);
        _normalizedCache[key] = result;
        return result;
    }

    // 正規化 data を実際に計算する。memo 無しの本体。
    private byte[]? ComputeNormalized(FormKey key)
    {
        if (!_records.TryGetValue(key, out var loc)) return null;

        var recordType = Encoding.ASCII.GetString(_data, loc.Offset, 4);
        var body = (loc.Flags & CompressedFlag) != 0
            ? Inflate(_data.AsSpan(loc.Offset + 24, loc.DataSize))
            : _data.AsSpan(loc.Offset + 24, loc.DataSize).ToArray();

        var stringTags = StringFieldTags.GetValueOrDefault(recordType, []);
        var entries = new List<(string Tag, byte[] Payload)>();
        var pos = 0;
        var extendedSize = 0; // XXXX subrecord による次 field の拡張サイズ（64KB 超 field 用）
        while (pos + 6 <= body.Length)
        {
            var tag = Encoding.ASCII.GetString(body, pos, 4);
            var size = (int)BinaryPrimitives.ReadUInt16LittleEndian(body.AsSpan(pos + 4));
            if (tag == "XXXX")
            {
                extendedSize = (int)BinaryPrimitives.ReadUInt32LittleEndian(body.AsSpan(pos + 6));
                pos += 6 + size;
                continue;
            }
            if (extendedSize > 0)
            {
                size = extendedSize;
                extendedSize = 0;
            }
            var payload = body.AsSpan(pos + 6, size);
            pos += 6 + size;

            // CK が再保存時に書き換える翻訳非関連 field は差分にしない（実測）:
            //   XCLW = cell 水面高の sentinel 揺れ、XLCN = cell location の自動付与、
            //   XSCL = 配置 ref の scale（float 1ulp 揺れ）、XLOC = lock data（未使用 byte 揺れ）。
            if (tag is "XCLW" or "XLCN" or "XSCL" or "XLOC") continue;

            byte[] normalized;
            if (tag == "VMAD")
            {
                // script は再保存で並びだけ変わる（DA03 で実測）。byte を整列した正規形で
                // 比較し、並び替えを差分にせず、内容の実変更（fragment 追加等）は検出する。
                normalized = payload.ToArray();
                Array.Sort(normalized);
            }
            else if (_localized && size == 4 && stringTags.Contains(tag))
            {
                // string ID を全言語の本文へ解決して書く（ID 振り直しを差分にせず、
                // どちらかの言語の本文変更は差分として残す）。
                var id = BinaryPrimitives.ReadUInt32LittleEndian(payload);
                var text = id == 0 ? "" : string.Join('',
                    _stringsByLanguage.Select(t => t.GetValueOrDefault(id, $"<missing:{id:x}>")));
                normalized = Encoding.UTF8.GetBytes(text);
            }
            else
            {
                normalized = payload.ToArray();
            }
            entries.Add((tag, normalized));
        }

        // subrecord の並びも CK 再保存で揺れる（Breezehome の REFR で実測）ため、
        // tag + payload で整列した正規形で比較する。
        entries.Sort((a, b) =>
        {
            var c = string.CompareOrdinal(a.Tag, b.Tag);
            return c != 0 ? c : a.Payload.AsSpan().SequenceCompareTo(b.Payload);
        });

        using var output = new MemoryStream();
        foreach (var (tag, payload) in entries)
        {
            output.Write(Encoding.ASCII.GetBytes(tag));
            output.Write(BitConverter.GetBytes(payload.Length));
            output.Write(payload);
        }
        return output.ToArray();
    }

    private static List<ModKey> ReadMasters(byte[] data)
    {
        var masters = new List<ModKey>();
        var tes4Size = (int)BinaryPrimitives.ReadUInt32LittleEndian(data.AsSpan(4));
        var pos = 24;
        var end = 24 + tes4Size;
        while (pos + 6 <= end)
        {
            var tag = Encoding.ASCII.GetString(data, pos, 4);
            var size = (int)BinaryPrimitives.ReadUInt16LittleEndian(data.AsSpan(pos + 4));
            if (tag == "MAST")
                masters.Add(ModKey.FromFileName(Encoding.ASCII.GetString(data, pos + 6, size - 1)));
            pos += 6 + size;
        }
        return masters;
    }

    private static Dictionary<uint, string> LoadStrings(string dataFolder, ModKey modKey, Language language)
    {
        var entries = new Dictionary<uint, string>();
        var baseName = Path.GetFileNameWithoutExtension(modKey.FileName);
        var lang = language.ToString().ToLowerInvariant();

        foreach (var (ext, lengthPrefixed) in new[] { ("strings", false), ("dlstrings", true), ("ilstrings", true) })
        {
            var path = Path.Combine(dataFolder, "Strings", $"{baseName}_{lang}.{ext}");
            if (!File.Exists(path)) continue;

            var raw = File.ReadAllBytes(path);
            var count = (int)BinaryPrimitives.ReadUInt32LittleEndian(raw);
            var dataStart = 8 + count * 8;
            for (var i = 0; i < count; i++)
            {
                var id = BinaryPrimitives.ReadUInt32LittleEndian(raw.AsSpan(8 + i * 8));
                var offset = (int)BinaryPrimitives.ReadUInt32LittleEndian(raw.AsSpan(8 + i * 8 + 4));
                var pos = dataStart + offset;
                string text;
                if (lengthPrefixed)
                {
                    var len = (int)BinaryPrimitives.ReadUInt32LittleEndian(raw.AsSpan(pos));
                    text = Encoding.UTF8.GetString(raw, pos + 4, Math.Max(0, len - 1));
                }
                else
                {
                    var endPos = Array.IndexOf(raw, (byte)0, pos);
                    text = Encoding.UTF8.GetString(raw, pos, endPos - pos);
                }
                entries.TryAdd(id, text);
            }
        }
        return entries;
    }

    private static byte[] Inflate(ReadOnlySpan<byte> compressed)
    {
        var expectedSize = (int)BinaryPrimitives.ReadUInt32LittleEndian(compressed);
        using var input = new MemoryStream(compressed[4..].ToArray());
        using var zlib = new ZLibStream(input, CompressionMode.Decompress);
        var result = new byte[expectedSize];
        zlib.ReadExactly(result);
        return result;
    }

    // record type ごとの翻訳対象 string field tag（localized 時は 4 byte の string ID になる）。
    private static readonly Dictionary<string, HashSet<string>> StringFieldTags = new()
    {
        ["DIAL"] = ["FULL"],
        ["INFO"] = ["RNAM", "NAM1"],
        ["QUST"] = ["FULL", "CNAM", "NNAM"],
        ["RACE"] = ["FULL", "DESC"],
        ["FACT"] = ["FULL", "MNAM", "FNAM"],
        ["KEYM"] = ["FULL"], ["MISC"] = ["FULL"], ["LIGH"] = ["FULL"], ["CONT"] = ["FULL"],
        ["SLGM"] = ["FULL"], ["DOOR"] = ["FULL"], ["FURN"] = ["FULL"], ["APPA"] = ["FULL", "DESC"],
        ["HAZD"] = ["FULL"],
        ["ACTI"] = ["FULL", "RNAM"], ["FLOR"] = ["FULL", "RNAM", "ATTX"], ["TREE"] = ["FULL"],
        ["WEAP"] = ["FULL", "DESC"], ["ARMO"] = ["FULL", "DESC"], ["AMMO"] = ["FULL", "DESC"],
        ["ALCH"] = ["FULL", "DESC"], ["SCRL"] = ["FULL", "DESC"], ["INGR"] = ["FULL"],
        ["SPEL"] = ["FULL", "DESC"], ["ENCH"] = ["FULL"],
        ["MGEF"] = ["FULL", "DNAM"],
        ["SHOU"] = ["FULL", "DESC"], ["WOOP"] = ["FULL", "TNAM"],
        ["BOOK"] = ["FULL", "DESC", "CNAM"],
        ["LCTN"] = ["FULL"], ["WRLD"] = ["FULL"], ["CELL"] = ["FULL"],
        ["MESG"] = ["FULL", "DESC", "ITXT"],
        ["LSCR"] = ["DESC"],
        ["PERK"] = ["FULL", "DESC", "EPF2"],
        ["NPC_"] = ["FULL", "SHRT"],
        ["TACT"] = ["FULL"],
        ["SNCT"] = ["FULL"], ["EYES"] = ["FULL"], ["REGN"] = ["RDMP"],
        ["CLAS"] = ["FULL", "DESC"],
        ["AVIF"] = ["FULL", "DESC"],
        ["EXPL"] = ["FULL"], ["PROJ"] = ["FULL"], ["HDPT"] = ["FULL"],
        ["FLST"] = [], ["GMST"] = ["DATA"],
        ["REFR"] = ["FULL"],
        ["WATR"] = ["FULL"], ["COLL"] = ["DESC"], ["CLFM"] = ["FULL"],
        ["BPTD"] = [],
        ["DLBR"] = [], ["SCEN"] = [],
    };
}
