package sqlite

import (
	"context"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/jmoiron/sqlx"
)

const sqliteTestDatabaseFileName = "master-dictionary.sqlite3"

func TestOpenMasterDictionaryDatabaseCreatesDatabaseAndSeedsOnce(t *testing.T) {
	databasePath := filepath.Join(t.TempDir(), "db", sqliteTestDatabaseFileName)
	firstSeedTime := time.Date(2026, 4, 15, 9, 0, 0, 0, time.UTC)
	secondSeedTime := firstSeedTime.Add(24 * time.Hour)

	database, err := OpenMasterDictionaryDatabase(context.Background(), databasePath, []MasterDictionarySeedEntry{seedEntry("Whiterun", firstSeedTime)})
	if err != nil {
		t.Fatalf("expected initial sqlite open to succeed: %v", err)
	}
	assertSeedCount(t, database, 1)
	closeErr := database.Close()
	if closeErr != nil {
		t.Fatalf("expected initial sqlite close to succeed: %v", closeErr)
	}
	_, statErr := os.Stat(databasePath)
	if statErr != nil {
		t.Fatalf("expected sqlite database file to exist: %v", statErr)
	}

	reopenedDatabase, err := OpenMasterDictionaryDatabase(context.Background(), databasePath, []MasterDictionarySeedEntry{seedEntry("Solitude", secondSeedTime)})
	if err != nil {
		t.Fatalf("expected reopened sqlite open to succeed: %v", err)
	}
	t.Cleanup(func() {
		reopenCloseErr := reopenedDatabase.Close()
		if reopenCloseErr != nil {
			t.Fatalf("expected reopened sqlite close to succeed: %v", reopenCloseErr)
		}
	})
	assertSeedCount(t, reopenedDatabase, 1)

	var source string
	queryErr := reopenedDatabase.QueryRowContext(context.Background(), "SELECT source FROM master_dictionary_entries LIMIT 1").Scan(&source)
	if queryErr != nil {
		t.Fatalf("expected to load seeded row after reopen: %v", queryErr)
	}
	if source != "Whiterun" {
		t.Fatalf("expected original seed row to remain, got %q", source)
	}
}

func TestOpenMasterDictionaryDatabaseReappliesMigrationsOnExistingDatabase(t *testing.T) {
	databasePath := filepath.Join(t.TempDir(), "db", sqliteTestDatabaseFileName)
	database, err := OpenMasterDictionaryDatabase(context.Background(), databasePath, nil)
	if err != nil {
		t.Fatalf("expected initial sqlite open to succeed: %v", err)
	}
	_, execErr := database.ExecContext(
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
	if execErr != nil {
		t.Fatalf("expected insert before reopen to succeed: %v", execErr)
	}
	closeErr := database.Close()
	if closeErr != nil {
		t.Fatalf("expected initial sqlite close to succeed: %v", closeErr)
	}

	reopenedDatabase, err := OpenMasterDictionaryDatabase(context.Background(), databasePath, nil)
	if err != nil {
		t.Fatalf("expected reopened sqlite open to succeed: %v", err)
	}
	t.Cleanup(func() {
		reopenCloseErr := reopenedDatabase.Close()
		if reopenCloseErr != nil {
			t.Fatalf("expected reopened sqlite close to succeed: %v", reopenCloseErr)
		}
	})

	assertSeedCount(t, reopenedDatabase, 1)
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
