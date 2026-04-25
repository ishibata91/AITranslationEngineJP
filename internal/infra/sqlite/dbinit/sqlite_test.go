package dbinit

import (
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/jmoiron/sqlx"
)

const (
	sqliteTestDatabaseFileName = "master-dictionary.sqlite3"
	errInitialSQLiteClose      = "expected initial sqlite close to succeed: %v"
)

func TestOpenMasterDictionaryDatabaseCreatesDatabaseFile(t *testing.T) {
	databasePath := filepath.Join(t.TempDir(), "db", sqliteTestDatabaseFileName)
	firstSeedTime := time.Date(2026, 4, 15, 9, 0, 0, 0, time.UTC)
	database := openMasterDictionaryDatabaseForTest(t, databasePath, []MasterDictionarySeedEntry{seedEntry("Whiterun", firstSeedTime)})

	err := database.Close()
	if err != nil {
		t.Fatalf(errInitialSQLiteClose, err)
	}

	_, err = os.Stat(databasePath)
	if err != nil {
		t.Fatalf("expected sqlite database file to exist: %v", err)
	}
}

func TestOpenMasterDictionaryDatabaseSeedsOnlyOnce(t *testing.T) {
	databasePath := filepath.Join(t.TempDir(), "db", sqliteTestDatabaseFileName)
	firstSeedTime := time.Date(2026, 4, 15, 9, 0, 0, 0, time.UTC)
	secondSeedTime := firstSeedTime.Add(24 * time.Hour)
	database := openMasterDictionaryDatabaseForTest(t, databasePath, []MasterDictionarySeedEntry{seedEntry("Whiterun", firstSeedTime)})
	assertTableNotExists(t, database, "master_dictionary_entries")

	err := database.Close()
	if err != nil {
		t.Fatalf(errInitialSQLiteClose, err)
	}

	reopenedDatabase := openMasterDictionaryDatabaseForTest(t, databasePath, []MasterDictionarySeedEntry{seedEntry("Solitude", secondSeedTime)})
	assertTableNotExists(t, reopenedDatabase, "master_dictionary_entries")
}

func TestOpenMasterDictionaryDatabasePreservesOriginalSeedDataOnReopen(t *testing.T) {
	databasePath := filepath.Join(t.TempDir(), "db", sqliteTestDatabaseFileName)
	firstSeedTime := time.Date(2026, 4, 15, 9, 0, 0, 0, time.UTC)
	secondSeedTime := firstSeedTime.Add(24 * time.Hour)
	database := openMasterDictionaryDatabaseForTest(t, databasePath, []MasterDictionarySeedEntry{seedEntry("Whiterun", firstSeedTime)})

	err := database.Close()
	if err != nil {
		t.Fatalf(errInitialSQLiteClose, err)
	}

	reopenedDatabase := openMasterDictionaryDatabaseForTest(t, databasePath, []MasterDictionarySeedEntry{seedEntry("Solitude", secondSeedTime)})
	assertTableNotExists(t, reopenedDatabase, "master_dictionary_entries")
	assertTableExists(t, reopenedDatabase, "PERSONA_GENERATION_SETTINGS")
}

func TestOpenMasterDictionaryDatabaseReappliesMigrationsOnExistingDatabase(t *testing.T) {
	databasePath := filepath.Join(t.TempDir(), "db", sqliteTestDatabaseFileName)
	database := openMasterDictionaryDatabaseForTest(t, databasePath, nil)

	err := database.Close()
	if err != nil {
		t.Fatalf(errInitialSQLiteClose, err)
	}

	reopenedDatabase := openMasterDictionaryDatabaseForTest(t, databasePath, nil)
	assertTableNotExists(t, reopenedDatabase, "master_dictionary_entries")
	assertTableExists(t, reopenedDatabase, "PERSONA_GENERATION_SETTINGS")
	assertColumnExists(t, reopenedDatabase, "X_EDIT_EXTRACTED_DATA", "source_content_hash")
	assertIndexExists(t, reopenedDatabase, "idx_x_edit_extracted_data_source_content_hash")
}

func TestOpenMasterDictionaryDatabaseCreatesSourceContentHashUniqueIndex(t *testing.T) {
	databasePath := filepath.Join(t.TempDir(), "db", sqliteTestDatabaseFileName)
	database := openMasterDictionaryDatabaseForTest(t, databasePath, nil)

	assertColumnExists(t, database, "X_EDIT_EXTRACTED_DATA", "source_content_hash")
	assertIndexExists(t, database, "idx_x_edit_extracted_data_source_content_hash")

	insertExtractedDataRow(t, database, "/mods/source-1.json", "hash-shared")

	_, err := database.ExecContext(context.Background(),
		`INSERT INTO X_EDIT_EXTRACTED_DATA (
			source_file_path,
			source_tool,
			target_plugin_name,
			target_plugin_type,
			record_count,
			imported_at,
			source_content_hash
		) VALUES (?, ?, ?, ?, ?, ?, ?)`,
		"/mods/source-2.json",
		"xEdit",
		"Skyrim.esm",
		"esm",
		1,
		"2026-04-26T09:31:00Z",
		"hash-shared",
	)
	if err == nil {
		t.Fatal("expected duplicate source_content_hash to fail due to unique index")
	}
	if !strings.Contains(err.Error(), "UNIQUE constraint failed") {
		t.Fatalf("expected UNIQUE constraint error, got: %v", err)
	}

	insertExtractedDataRow(t, database, "/mods/source-empty-1.json", "")
	insertExtractedDataRow(t, database, "/mods/source-empty-2.json", "")
}

func TestOpenMasterDictionaryDatabaseRecreatesSourceContentHashIndexWhenColumnAlreadyExists(t *testing.T) {
	databasePath := filepath.Join(t.TempDir(), "db", sqliteTestDatabaseFileName)
	database := openMasterDictionaryDatabaseForTest(t, databasePath, nil)

	err := database.Close()
	if err != nil {
		t.Fatalf(errInitialSQLiteClose, err)
	}

	rawDatabase, err := sqlx.Open(sqliteDriverName, buildSQLiteDSN(databasePath))
	if err != nil {
		t.Fatalf("expected raw sqlite open to succeed: %v", err)
	}

	_, err = rawDatabase.ExecContext(context.Background(), `DROP INDEX IF EXISTS idx_x_edit_extracted_data_source_content_hash`)
	if err != nil {
		t.Fatalf("expected source_content_hash index drop to succeed: %v", err)
	}

	if closeErr := rawDatabase.Close(); closeErr != nil {
		t.Fatalf("expected raw sqlite close to succeed: %v", closeErr)
	}

	reopenedDatabase := openMasterDictionaryDatabaseForTest(t, databasePath, nil)

	assertColumnExists(t, reopenedDatabase, "X_EDIT_EXTRACTED_DATA", "source_content_hash")
	assertIndexExists(t, reopenedDatabase, "idx_x_edit_extracted_data_source_content_hash")

	insertExtractedDataRow(t, reopenedDatabase, "/mods/reopened-1.json", "rehydrated-hash")

	_, err = reopenedDatabase.ExecContext(context.Background(),
		`INSERT INTO X_EDIT_EXTRACTED_DATA (
			source_file_path,
			source_tool,
			target_plugin_name,
			target_plugin_type,
			record_count,
			imported_at,
			source_content_hash
		) VALUES (?, ?, ?, ?, ?, ?, ?)`,
		"/mods/reopened-2.json",
		"xEdit",
		"Skyrim.esm",
		"esm",
		2,
		"2026-04-26T09:32:00Z",
		"rehydrated-hash",
	)
	if err == nil {
		t.Fatal("expected duplicate source_content_hash to fail after migrations reapply")
	}
	if !strings.Contains(err.Error(), "UNIQUE constraint failed") {
		t.Fatalf("expected UNIQUE constraint error after migrations reapply, got: %v", err)
	}
}

// TestSchemaCutoverDropsLegacyPersonaAndDictionaryTables は schema cutover 後に
// legacy master_* テーブルが存在しないことを検証する (completion_signal)。
func TestSchemaCutoverDropsLegacyPersonaAndDictionaryTables(t *testing.T) {
	databasePath := filepath.Join(t.TempDir(), "db", sqliteTestDatabaseFileName)
	database := openMasterDictionaryDatabaseForTest(t, databasePath, nil)

	assertTableNotExists(t, database, "master_dictionary_entries")
	assertTableNotExists(t, database, "master_persona_entries")
	assertTableNotExists(t, database, "master_persona_ai_settings")
	assertTableNotExists(t, database, "master_persona_run_status")
}

// TestSchemaCutoverCreatesPersonaGenerationSettingsTable は schema cutover 後に
// PERSONA_GENERATION_SETTINGS テーブルが存在することを検証する (completion_signal)。
func TestSchemaCutoverCreatesPersonaGenerationSettingsTable(t *testing.T) {
	databasePath := filepath.Join(t.TempDir(), "db", sqliteTestDatabaseFileName)
	database := openMasterDictionaryDatabaseForTest(t, databasePath, nil)

	assertTableExists(t, database, "PERSONA_GENERATION_SETTINGS")
}

func openMasterDictionaryDatabaseForTest(t *testing.T, databasePath string, seeds []MasterDictionarySeedEntry) *sqlx.DB {
	t.Helper()

	database, err := OpenMasterDictionaryDatabase(context.Background(), databasePath, seeds)
	if err != nil {
		t.Fatalf("expected sqlite open to succeed: %v", err)
	}

	t.Cleanup(func() {
		if closeErr := database.Close(); closeErr != nil {
			t.Fatalf("expected sqlite close to succeed: %v", closeErr)
		}
	})

	return database
}

func seedEntry(source string, updatedAt time.Time) MasterDictionarySeedEntry {
	return MasterDictionarySeedEntry{
		Source:      source,
		Translation: source + "訳",
		Category:    "地名",
		Origin:      "初期データ",
		REC:         "LCTN:FULL",
		EDID:        "Loc" + source,
		UpdatedAt:   updatedAt,
	}
}

func insertExtractedDataRow(t *testing.T, database *sqlx.DB, sourceFilePath string, sourceContentHash string) {
	t.Helper()

	_, err := database.ExecContext(context.Background(),
		`INSERT INTO X_EDIT_EXTRACTED_DATA (
			source_file_path,
			source_tool,
			target_plugin_name,
			target_plugin_type,
			record_count,
			imported_at,
			source_content_hash
		) VALUES (?, ?, ?, ?, ?, ?, ?)`,
		sourceFilePath,
		"xEdit",
		"Skyrim.esm",
		"esm",
		1,
		"2026-04-26T09:30:00Z",
		sourceContentHash,
	)
	if err != nil {
		t.Fatalf("expected X_EDIT_EXTRACTED_DATA insert to succeed: %v", err)
	}
}

func assertTableExists(t *testing.T, database *sqlx.DB, tableName string) {
	t.Helper()

	var count int
	queryErr := database.QueryRowContext(
		context.Background(),
		"SELECT COUNT(*) FROM sqlite_master WHERE type='table' AND name = ?",
		tableName,
	).Scan(&count)
	if queryErr != nil {
		t.Fatalf("expected sqlite table existence query to succeed: %v", queryErr)
	}
	if count != 1 {
		t.Fatalf("expected sqlite table %q to exist", tableName)
	}
}

func assertTableNotExists(t *testing.T, database *sqlx.DB, tableName string) {
	t.Helper()

	var count int
	queryErr := database.QueryRowContext(
		context.Background(),
		"SELECT COUNT(*) FROM sqlite_master WHERE type='table' AND name = ?",
		tableName,
	).Scan(&count)
	if queryErr != nil {
		t.Fatalf("expected sqlite table non-existence query to succeed: %v", queryErr)
	}
	if count != 0 {
		t.Fatalf("expected sqlite table %q to not exist, but it was found", tableName)
	}
}

func assertColumnExists(t *testing.T, database *sqlx.DB, tableName string, columnName string) {
	t.Helper()

	var count int
	queryErr := database.QueryRowContext(
		context.Background(),
		"SELECT COUNT(*) FROM pragma_table_info('"+tableName+"') WHERE name = ?",
		columnName,
	).Scan(&count)
	if queryErr != nil {
		t.Fatalf("expected sqlite column existence query to succeed: %v", queryErr)
	}
	if count != 1 {
		t.Fatalf("expected sqlite column %q to exist on table %q", columnName, tableName)
	}
}

func assertIndexExists(t *testing.T, database *sqlx.DB, indexName string) {
	t.Helper()

	var count int
	queryErr := database.QueryRowContext(
		context.Background(),
		"SELECT COUNT(*) FROM sqlite_master WHERE type='index' AND name = ?",
		indexName,
	).Scan(&count)
	if queryErr != nil {
		t.Fatalf("expected sqlite index existence query to succeed: %v", queryErr)
	}
	if count != 1 {
		t.Fatalf("expected sqlite index %q to exist", indexName)
	}
}
