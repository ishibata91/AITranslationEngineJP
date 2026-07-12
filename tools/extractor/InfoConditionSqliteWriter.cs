using Microsoft.Data.Sqlite;
using Mutagen.Bethesda.Plugins;
using Mutagen.Bethesda.Skyrim;

namespace Extractor;

// INFO の条件から導いた性別（extracted_info_condition）を中心 DB へ書く。
// 話者を解決できない汎用台詞の一人称・語尾の根拠にする。話者解決の有無に依らず、性別を定められる INFO を全て書く。
// 性別は次の優先で取る: GetIsSex（極性を畳む）→ 単一声型 EDID の Male/Female 接頭 → 同性のみの声型リスト(FLST)。
// line 行は Go 取込段が作るため、INFO の安定キー（plugin・form_id）で一時保持し、Go 取込段が line_condition へ解決する。
// 書くのは事実（性別）だけで、口調などの解釈は持たない（解釈は Go engine が与える）。
public static class InfoConditionSqliteWriter
{
    // dbPath の SQLite に schema を ensure し、INFO→条件由来の性別を書く。書いた行数を返す。
    public static int Write(string dbPath, string migrationsDir, PluginEnvironment env, string targetPlugin)
    {
        using var conn = new SqliteConnection($"Data Source={dbPath}");
        conn.Open();

        SchemaMigrator.Ensure(conn, migrationsDir);

        using var tx = conn.BeginTransaction();
        var written = 0;
        foreach (var topic in env.TargetMod.DialogTopics)
        {
            foreach (var info in topic.Responses)
            {
                var sex = SexOf(env, info);
                if (sex.Length == 0) continue;
                written += UpsertCondition(conn, tx, targetPlugin, HexFormId(info.FormKey), sex);
            }
        }
        tx.Commit();
        return written;
    }

    private static int UpsertCondition(SqliteConnection conn, SqliteTransaction tx, string plugin, string formId, string sex)
    {
        using var cmd = conn.CreateCommand();
        cmd.Transaction = tx;
        cmd.CommandText = "INSERT OR IGNORE INTO extracted_info_condition (info_plugin, info_form_id, sex) VALUES ($p, $f, $s)";
        cmd.Parameters.AddWithValue("$p", plugin);
        cmd.Parameters.AddWithValue("$f", formId);
        cmd.Parameters.AddWithValue("$s", sex);
        return cmd.ExecuteNonQuery();
    }

    private static string HexFormId(FormKey key) => $"0x{key.ID:X6}";

    // SexOf は 1 INFO の条件から実効的な性別を返す（"Male"/"Female"/空）。
    // GetIsSex を最優先（極性を畳む）。無ければ声型（単一 VTYP の接頭 / 同性のみ FLST）から取る。
    private static string SexOf(PluginEnvironment env, IDialogResponsesGetter info)
    {
        string? fromSex = null;
        string? fromVoice = null;
        foreach (var cond in info.Conditions)
        {
            var asserts = AssertsTrue(cond);
            switch (cond.Data)
            {
                case IGetIsSexConditionDataGetter d:
                    // GetIsSex Male == 1 は男、== 0 は「男でない」= 女。極性を畳む。
                    var queried = d.MaleFemaleGender.ToString();
                    fromSex = asserts ? queried : (queried == "Male" ? "Female" : "Male");
                    break;
                case IGetIsVoiceTypeConditionDataGetter d when asserts && fromVoice == null:
                    fromVoice = VoiceSex(env, d.VoiceTypeOrList.Link.FormKey);
                    break;
            }
        }
        return fromSex ?? fromVoice ?? "";
    }

    // VoiceSex は声型条件（単一 VTYP かリスト FLST）から性別を取る。単一は EDID の Male/Female 接頭、
    // リストは構成員が全て同性ならその性別、混在または不明なら null。
    private static string? VoiceSex(PluginEnvironment env, FormKey key)
    {
        if (env.LinkCache.TryResolve<IVoiceTypeGetter>(key, out var vt))
            return SexOfEdid(vt.EditorID);
        if (env.LinkCache.TryResolve<IFormListGetter>(key, out var fl))
            return SexOfVoiceList(env, fl);
        return null;
    }

    // SexOfEdid は声型 EditorID の Male/Female 接頭から性別を取る。どちらでもなければ null。
    private static string? SexOfEdid(string? edid)
    {
        if (edid == null) return null;
        if (edid.StartsWith("Male", StringComparison.Ordinal)) return "Male";
        if (edid.StartsWith("Female", StringComparison.Ordinal)) return "Female";
        return null;
    }

    // SexOfVoiceList は声型リスト(FLST)の構成員が全て同性ならその性別を返す。混在または不明な構成員があれば null。
    private static string? SexOfVoiceList(PluginEnvironment env, IFormListGetter fl)
    {
        string? sex = null;
        foreach (var item in fl.Items)
        {
            if (!env.LinkCache.TryResolve<IVoiceTypeGetter>(item.FormKey, out var vt)) continue;
            var s = SexOfEdid(vt.EditorID);
            if (s == null) return null;     // 性別不明の構成員がいれば混在扱い
            if (sex == null) sex = s;
            else if (sex != s) return null; // 異性が混ざれば空
        }
        return sex;
    }

    // AssertsTrue は条件が「関数が真（=1）」を主張しているか（PluginExtractor.AssertsTrue と同等）。
    private static bool AssertsTrue(IConditionGetter cond)
    {
        if (cond is not IConditionFloatGetter f) return true;
        return cond.CompareOperator switch
        {
            CompareOperator.EqualTo => f.ComparisonValue != 0,
            CompareOperator.NotEqualTo => f.ComparisonValue == 0,
            CompareOperator.GreaterThan => f.ComparisonValue is >= 0 and < 1,
            CompareOperator.GreaterThanOrEqualTo => f.ComparisonValue is > 0 and <= 1,
            _ => false,
        };
    }
}
