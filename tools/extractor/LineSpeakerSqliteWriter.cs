using Microsoft.Data.Sqlite;
using Mutagen.Bethesda.Plugins;
using Mutagen.Bethesda.Plugins.Cache;
using Mutagen.Bethesda.Skyrim;

namespace Extractor;

// 抽出結果の台詞（INFO:NAM1）と、その話者属性（speaker / race / faction / voice_type）を中心 DB へ書く。
// 話者は esp に無く master 側 NPC にいるため、INFO の話者 FormKey を LinkCache で解決し、
// NPC の種族・声型・所属勢力を引いて書く。書くのは識別子と事実（EditorID / FormKey）に限り、
// 口調などの解釈は持たない（解釈は Go engine のルールが与える）。schema は db/migrations を冪等 ensure する。
public static class LineSpeakerSqliteWriter
{
    // dbPath の SQLite に schema を ensure し、台詞と話者属性を書く。書いた台詞の行数を返す。
    public static int Write(string dbPath, string schemaSql, ExtractionResult result, ILinkCache linkCache)
    {
        using var conn = new SqliteConnection($"Data Source={dbPath}");
        conn.Open();

        using (var ddl = conn.CreateCommand())
        {
            ddl.CommandText = schemaSql;
            ddl.ExecuteNonQuery();
        }

        using var tx = conn.BeginTransaction();
        var w = new Writer(conn, tx, result.TargetPlugin);

        var lineCount = 0;
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

                foreach (var response in info.Responses)
                {
                    var lineId = w.UpsertLine(info, response);
                    lineCount++;
                    foreach (var speakerId in speakerIds)
                        w.UpsertLineSpeaker(lineId, speakerId);
                }
            }
        }

        tx.Commit();
        return lineCount;
    }

    // Writer は prepared statement をまとめ、台詞・話者の UPSERT と id 取得を行う。
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

        // UpsertLine は台詞 1 件（INFO:NAM1）を line へ入れ、id を返す。複数 response は ordinal=Number で分ける。
        public long UpsertLine(InfoNode info, ResponseLine response)
        {
            var formId = HexFormId(info.Id);
            Exec(
                """
                INSERT OR IGNORE INTO line (source, response_order, plugin, form_id, edid, rec, field, ordinal)
                VALUES ($source, $order, $plugin, $form_id, $edid, 'INFO', 'NAM1', $ordinal)
                """,
                ("$source", response.Text), ("$order", response.Number),
                ("$plugin", _targetPlugin), ("$form_id", formId),
                ("$edid", info.EditorId), ("$ordinal", response.Number));
            return ScalarId(
                "SELECT id FROM line WHERE plugin=$plugin AND form_id=$form_id AND rec='INFO' AND field='NAM1' AND ordinal=$ordinal",
                ("$plugin", _targetPlugin), ("$form_id", formId), ("$ordinal", response.Number));
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

        // UpsertLineSpeaker は台詞と話者の対応（e6）を入れる。
        public void UpsertLineSpeaker(long lineId, long speakerId)
            => Exec("INSERT OR IGNORE INTO line_speaker (line_id, speaker_id) VALUES ($l, $s)",
                ("$l", lineId), ("$s", speakerId));

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

        // Command は transaction 付きのパラメータ化コマンドを組む。Exec / ScalarId 共通。
        private SqliteCommand Command(string sql, (string Name, object Value)[] ps)
        {
            var cmd = _conn.CreateCommand();
            cmd.Transaction = _tx;
            cmd.CommandText = sql;
            foreach (var (name, value) in ps) cmd.Parameters.AddWithValue(name, value);
            return cmd;
        }

        private static string HexFormId(FormKey key) => $"0x{key.ID:X6}";
    }
}
