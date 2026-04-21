package dbinit

import (
	"context"
	"os"
	"path/filepath"
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
