package dbinit

import (
	"context"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"testing"

	"github.com/jmoiron/sqlx"
)

func TestProviderSettingsMigrationCreatesProviderSettingsTable(t *testing.T) {
	db := openMasterDictionaryDatabaseForTest(t, filepath.Join(t.TempDir(), "db", "provider-settings-migration.sqlite3"), nil)
	assertTableExists(t, db, "PROVIDER_SETTINGS")
}

func TestProviderSettingsMigrationAppliesToMigratedSQLiteDatabase(t *testing.T) {
	databasePath := filepath.Join(t.TempDir(), "db", "provider-settings-migrated.sqlite3")
	if err := os.MkdirAll(filepath.Dir(databasePath), 0o750); err != nil {
		t.Fatalf("expected sqlite db directory creation to succeed: %v", err)
	}
	rawDB, err := sqlx.Open(sqliteDriverName, buildSQLiteDSN(databasePath))
	if err != nil {
		t.Fatalf("expected raw sqlite open to succeed: %v", err)
	}
	if err := rawDB.PingContext(context.Background()); err != nil {
		t.Fatalf("expected raw sqlite ping to succeed: %v", err)
	}
	if err := applySQLiteMigrationsUpToProviderSettings(context.Background(), rawDB); err != nil {
		t.Fatalf("expected pre-provider-settings migrations to succeed: %v", err)
	}
	if closeErr := rawDB.Close(); closeErr != nil {
		t.Fatalf("expected raw sqlite close to succeed: %v", closeErr)
	}

	reopened := openMasterDictionaryDatabaseForTest(t, databasePath, nil)
	assertTableExists(t, reopened, "PROVIDER_SETTINGS")
}

func applySQLiteMigrationsUpToProviderSettings(ctx context.Context, database *sqlx.DB) error {
	files, err := fs.Glob(migrationFiles, "migrations/*.sql")
	if err != nil {
		return fmt.Errorf("list sqlite migration files: %w", err)
	}
	sort.Strings(files)
	for _, migrationPath := range files {
		if strings.Contains(migrationPath, "010_provider_settings.sql") {
			continue
		}
		sqlBytes, readErr := fs.ReadFile(migrationFiles, migrationPath)
		if readErr != nil {
			return fmt.Errorf("read sqlite migration file %s: %w", migrationPath, readErr)
		}
		if _, execErr := database.ExecContext(ctx, string(sqlBytes)); execErr != nil && !shouldIgnoreMigrationError(migrationPath, execErr) {
			return fmt.Errorf("apply sqlite migration file %s: %w", migrationPath, execErr)
		}
	}
	return nil
}
