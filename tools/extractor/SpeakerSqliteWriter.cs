using Microsoft.Data.Sqlite;
using Mutagen.Bethesda.Plugins;
using Mutagen.Bethesda.Plugins.Cache;
using Mutagen.Bethesda.Skyrim;

namespace Extractor;

// 台詞の話者属性（speaker / race / faction / voice_type）と、INFO→speaker の橋渡し（extracted_info_speaker）を中心 DB へ書く。
// 台詞本文（INFO:NAM1/RNAM）は ExtractedFieldSqliteWriter が extracted_field へ素朴に書き、line 行は Go 取込段が作る。
// そのため本 writer は line/line_speaker を書かず、INFO の安定キー（plugin・form_id）へ解決済み speaker.id を結ぶ。
// 話者は esp に無く master 側 NPC にいるため、INFO の話者 FormKey を LinkCache で解決して書く。
// 書くのは識別子と事実（EditorID / FormKey）に限り、口調などの解釈は持たない（解釈は Go engine が与える）。
public static class SpeakerSqliteWriter
{
    // dbPath の SQLite に schema を ensure し、話者属性と INFO→speaker の橋渡しを書く。書いた橋渡し行数を返す。
    public static int Write(string dbPath, string migrationsDir, ExtractionResult result, ILinkCache linkCache)
    {
        using var conn = new SqliteConnection($"Data Source={dbPath}");
        conn.Open();

        SchemaMigrator.Ensure(conn, migrationsDir);

        using var tx = conn.BeginTransaction();
        var w = new Writer(conn, tx, result.TargetPlugin);

        var linkCount = 0;
        foreach (var dialogue in result.Dialogues)
        {
            foreach (var info in dialogue.Infos)
            {
                // INFO の話者 NPC を 1 度だけ解決し、speaker 行 id を得る。
                var speakerIds = new List<long>();
                foreach (var speakerKey in info.SpeakerIds)
                {
                    if (linkCache.TryResolve<INpcGetter>(speakerKey, out var npc))
                        speakerIds.Add(w.UpsertSpeaker(npc, linkCache));
                }

                // INFO（plugin・form_id）→ speaker.id の橋渡しを書く。Go 取込段が line を作った後に line_speaker へ解決する。
                var infoFormId = HexFormId(info.Id);
                foreach (var speakerId in speakerIds)
                    linkCount += w.UpsertInfoSpeaker(infoFormId, speakerId);
            }
        }

        tx.Commit();
        return linkCount;
    }

    private static string HexFormId(FormKey key) => $"0x{key.ID:X6}";

    // Writer は prepared statement をまとめ、話者の UPSERT・id 取得と、INFO→speaker 橋渡しの書き込みを行う。
    private sealed class Writer
    {
        private readonly SqliteConnection _conn;
        private readonly SqliteTransaction _tx;
        private readonly string _targetPlugin;

        public Writer(SqliteConnection conn, SqliteTransaction tx, string targetPlugin)
        {
            _conn = conn;
            _tx = tx;
            _targetPlugin = targetPlugin;
        }

        // UpsertSpeaker は NPC とその種族・声型・所属勢力を解決して書き、speaker の id を返す。
        public long UpsertSpeaker(INpcGetter npc, ILinkCache linkCache)
        {
            long? raceId = null;
            if (linkCache.TryResolve<IRaceGetter>(npc.Race.FormKey, out var race))
                raceId = UpsertNamed("race", "INSERT OR IGNORE INTO race (edid, form_id, plugin) VALUES ($edid, $form_id, $plugin)", race.EditorID, race.FormKey);

            long? voiceId = null;
            if (linkCache.TryResolve<IVoiceTypeGetter>(npc.Voice.FormKey, out var voice))
                voiceId = UpsertNamed("voice_type", "INSERT OR IGNORE INTO voice_type (edid, voice_id, form_id, plugin) VALUES ($edid, $edid, $form_id, $plugin)", voice.EditorID, voice.FormKey);

            var formId = HexFormId(npc.FormKey);
            var plugin = npc.FormKey.ModKey.FileName.ToString();
            // 性別は NPC の Female フラグ（実体メタ）から取る。声型の Female/Male 接頭はユニーク NPC で当てにならず使わない。
            var sex = npc.Configuration.Flags.HasFlag(NpcConfiguration.Flag.Female) ? "Female" : "Male";
            Exec(
                """
                INSERT OR IGNORE INTO speaker (edid, form_id, plugin, sex, race_id, voice_type_id)
                VALUES ($edid, $form_id, $plugin, $sex, $race_id, $voice_id)
                """,
                ("$edid", npc.EditorID ?? ""), ("$form_id", formId), ("$plugin", plugin), ("$sex", sex),
                ("$race_id", (object?)raceId ?? DBNull.Value), ("$voice_id", (object?)voiceId ?? DBNull.Value));
            var speakerId = ScalarId(
                "SELECT id FROM speaker WHERE plugin=$plugin AND form_id=$form_id",
                ("$plugin", plugin), ("$form_id", formId));

            foreach (var rank in npc.Factions)
            {
                if (!linkCache.TryResolve<IFactionGetter>(rank.Faction.FormKey, out var faction)) continue;
                var factionId = UpsertNamed("faction", "INSERT OR IGNORE INTO faction (edid, form_id, plugin) VALUES ($edid, $form_id, $plugin)", faction.EditorID, faction.FormKey);
                Exec("INSERT OR IGNORE INTO speaker_faction (speaker_id, faction_id) VALUES ($s, $f)",
                    ("$s", speakerId), ("$f", factionId));
            }
            return speakerId;
        }

        // UpsertInfoSpeaker は INFO（plugin・form_id）→ speaker.id の橋渡しを書く。書けた行数（0/1）を返す。
        public int UpsertInfoSpeaker(string infoFormId, long speakerId)
        {
            using var cmd = Command(
                "INSERT OR IGNORE INTO extracted_info_speaker (info_plugin, info_form_id, speaker_id) VALUES ($p, $f, $s)",
                [("$p", _targetPlugin), ("$f", infoFormId), ("$s", speakerId)]);
            return cmd.ExecuteNonQuery();
        }

        // UpsertNamed は (edid, form_id, plugin) を持つ素材テーブル（race/faction/voice_type）へ入れ、id を返す。
        private long UpsertNamed(string table, string insertSql, string? edid, FormKey key)
        {
            var formId = HexFormId(key);
            var plugin = key.ModKey.FileName.ToString();
            Exec(insertSql, ("$edid", edid ?? ""), ("$form_id", formId), ("$plugin", plugin));
            return ScalarId($"SELECT id FROM {table} WHERE plugin=$plugin AND form_id=$form_id",
                ("$plugin", plugin), ("$form_id", formId));
        }

        private void Exec(string sql, params (string Name, object Value)[] ps)
        {
            using var cmd = Command(sql, ps);
            cmd.ExecuteNonQuery();
        }

        private long ScalarId(string sql, params (string Name, object Value)[] ps)
        {
            using var cmd = Command(sql, ps);
            return Convert.ToInt64(cmd.ExecuteScalar());
        }

        // Command は transaction 付きのパラメータ化コマンドを組む。Exec / ScalarId / UpsertInfoSpeaker 共通。
        private SqliteCommand Command(string sql, (string Name, object Value)[] ps)
        {
            var cmd = _conn.CreateCommand();
            cmd.Transaction = _tx;
            cmd.CommandText = sql;
            foreach (var (name, value) in ps) cmd.Parameters.AddWithValue(name, value);
            return cmd;
        }
    }
}
