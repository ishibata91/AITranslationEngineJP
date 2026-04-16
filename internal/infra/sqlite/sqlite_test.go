package sqlite

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
	assertSeedCount(t, database, 1)

	err := database.Close()
	if err != nil {
		t.Fatalf(errInitialSQLiteClose, err)
	}

	reopenedDatabase := openMasterDictionaryDatabaseForTest(t, databasePath, []MasterDictionarySeedEntry{seedEntry("Solitude", secondSeedTime)})
	assertSeedCount(t, reopenedDatabase, 1)
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

	var source string
	err = reopenedDatabase.QueryRowContext(context.Background(), "SELECT source FROM master_dictionary_entries LIMIT 1").Scan(&source)
	if err != nil {
		t.Fatalf("expected to load seeded row after reopen: %v", err)
	}

	if source != "Whiterun" {
		t.Fatalf("expected original seed row to remain, got %q", source)
	}
}

func TestOpenMasterDictionaryDatabaseReappliesMigrationsOnExistingDatabase(t *testing.T) {
	databasePath := filepath.Join(t.TempDir(), "db", sqliteTestDatabaseFileName)
	database := openMasterDictionaryDatabaseForTest(t, databasePath, nil)

	_, err := database.ExecContext(
		context.Background(),
		"INSERT INTO master_dictionary_entries (source, translation, category, origin, rec, edid, updated_at) VALUES (?, ?, ?, ?, ?, ?, ?)",
		"Riften",
		"リフテン",
		"地名",
		"手動登録",
		"LCTN:FULL",
		"LocRiften",
		time.Date(2026, 4, 15, 10, 0, 0, 0, time.UTC).Format(masterDictionaryTimestampLayout),
	)
	if err != nil {
		t.Fatalf("expected insert before reopen to succeed: %v", err)
	}

	err = database.Close()
	if err != nil {
		t.Fatalf(errInitialSQLiteClose, err)
	}

	reopenedDatabase := openMasterDictionaryDatabaseForTest(t, databasePath, nil)
	assertSeedCount(t, reopenedDatabase, 1)
}

func TestOpenMasterDictionaryDatabaseCreatesMasterPersonaTables(t *testing.T) {
	databasePath := filepath.Join(t.TempDir(), "db", sqliteTestDatabaseFileName)
	database := openMasterDictionaryDatabaseForTest(t, databasePath, nil)

	assertTableExists(t, database, "master_persona_entries")
	assertTableExists(t, database, "master_persona_ai_settings")
	assertTableExists(t, database, "master_persona_run_status")
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

func assertSeedCount(t *testing.T, database *sqlx.DB, expected int) {
	t.Helper()

	var count int
	queryErr := database.QueryRowContext(context.Background(), "SELECT COUNT(*) FROM master_dictionary_entries").Scan(&count)
	if queryErr != nil {
		t.Fatalf("expected count query to succeed: %v", queryErr)
	}
	if count != expected {
		t.Fatalf("expected %d rows, got %d", expected, count)
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
