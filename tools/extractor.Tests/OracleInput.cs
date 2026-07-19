using Extractor;
using Microsoft.Data.Sqlite;
using Mutagen.Bethesda;
using Mutagen.Bethesda.Plugins;
using Mutagen.Bethesda.Plugins.Records;
using Mutagen.Bethesda.Skyrim;

namespace Extractor.Tests;

// オラクルテストの Arrange を支える状態なしヘルパ。共有する可変状態を持たず、各テストが自分で呼ぶ（テスト独立の規約）。
// 索き（EditorID→FormKey / INFO / staging）だけを担い、期待値の判断は持たない（それはテスト本体）。
public static class OracleInput
{
    private static string Root => OraclePaths.RepoRoot();

    // Arrange: 合成 esm を実抽出の入口へ渡せる形で読む（呼び手が using で捨てる）。
    public static PluginEnvironment LoadEnv() =>
        PluginEnvironment.Load(System.IO.Path.Combine(Root, "test-oracle", "fixture"), "Synthetic.esm");

    public static string MigrationsDir => System.IO.Path.Combine(Root, "db", "migrations");

    // EditorID で record の FormKey を索く（種別を問わない）。
    public static FormKey Key(ISkyrimModGetter mod, string edid) =>
        mod.EnumerateMajorRecords().Single(r => r.EditorID == edid).FormKey;

    // EditorID で INFO を索く。
    public static InfoNode Info(ExtractionResult result, string edid) =>
        result.Dialogues.SelectMany(d => d.Infos).Single(i => i.EditorId == edid);

    // INFO の staging 上の form_id（writer と同じ 0x 表記）。
    public static string InfoFormId(ExtractionResult result, string edid) =>
        $"0x{Info(result, edid).Id.ID:X6}";

    // Act で writer が書く staging を読む temp DB（using で消える）。可変状態を各テストへ閉じ込め、独立を保つ。
    public static TempDb NewDb() => new();

    public sealed class TempDb : IDisposable
    {
        public string Path { get; } =
            System.IO.Path.Combine(System.IO.Path.GetTempPath(), $"oracle-{Guid.NewGuid():N}.sqlite3");

        // 感情型 staging を INFO form_id で索く。
        public List<(int Ordinal, string Emotion)> Emotions(string infoFormId) =>
            Rows(infoFormId,
                "SELECT ordinal, emotion_type FROM extracted_info_emotion WHERE info_form_id=$f ORDER BY ordinal",
                r => (r.GetInt32(0), r.GetString(1)));

        // 条件由来の性別を INFO form_id で索く。行が無ければ null。
        public string? Sex(string infoFormId) =>
            Rows(infoFormId, "SELECT sex FROM extracted_info_condition WHERE info_form_id=$f",
                r => r.GetString(0)).FirstOrDefault();

        private List<T> Rows<T>(string formId, string sql, Func<SqliteDataReader, T> map)
        {
            using var conn = new SqliteConnection($"Data Source={Path}");
            conn.Open();
            using var cmd = conn.CreateCommand();
            cmd.CommandText = sql;
            cmd.Parameters.AddWithValue("$f", formId);
            using var reader = cmd.ExecuteReader();
            var got = new List<T>();
            while (reader.Read()) got.Add(map(reader));
            return got;
        }

        public void Dispose()
        {
            if (File.Exists(Path)) File.Delete(Path);
        }
    }
}
