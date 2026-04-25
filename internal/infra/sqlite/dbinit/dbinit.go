// Package dbinit provides SQLite database initialization that can be used
// without importing the parent sqlite package, breaking the import cycle between
// internal/infra/sqlite and internal/repository.
package dbinit

import (
	"context"
	"embed"
	"errors"
	"fmt"
	"io/fs"
	"net/url"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/jmoiron/sqlx"
	_ "modernc.org/sqlite" // Register the modernc SQLite driver for database/sql.
)

const (
	sqliteDriverName                = "sqlite"
	masterDictionaryTimestampLayout = time.RFC3339Nano
	countAllMasterDictionaryEntries = `SELECT COUNT(*) FROM DICTIONARY_ENTRY WHERE dictionary_lifecycle = 'master';`
	insertMasterDictionaryEntry     = `
INSERT INTO DICTIONARY_ENTRY (
  dictionary_lifecycle,
  dictionary_scope,
  dictionary_source,
  source_term,
  translated_term,
  term_kind,
  reusable,
  created_at,
  updated_at
) VALUES (
  'master',
  :dictionary_scope,
  :dictionary_source,
  :source_term,
  :translated_term,
  :term_kind,
  1,
  :created_at,
  :updated_at
);`
)

//go:embed migrations/*.sql
var migrationFiles embed.FS

// MasterDictionarySeedEntry describes one startup seed row.
type MasterDictionarySeedEntry struct {
	Source      string
	Translation string
	Category    string
	Origin      string
	REC         string
	EDID        string
	UpdatedAt   time.Time
}

type masterDictionarySeedRow struct {
	Source      string `db:"source_term"`
	Translation string `db:"translated_term"`
	Category    string `db:"term_kind"`
	Origin      string `db:"dictionary_source"`
	REC         string `db:"dictionary_scope"`
	UpdatedAt   string `db:"updated_at"`
	CreatedAt   string `db:"created_at"`
}

// OpenMasterDictionaryDatabase opens the SQLite database, reapplies migrations, and seeds an empty database.
func OpenMasterDictionaryDatabase(
	ctx context.Context,
	databasePath string,
	seed []MasterDictionarySeedEntry,
) (*sqlx.DB, error) {
	resolvedPath, err := ensureDatabasePath(databasePath)
	if err != nil {
		return nil, err
	}

	database, err := sqlx.Open(sqliteDriverName, buildSQLiteDSN(resolvedPath))
	if err != nil {
		return nil, fmt.Errorf("open sqlite database: %w", err)
	}
	database.SetMaxOpenConns(1)
	database.SetMaxIdleConns(1)

	if err := database.PingContext(ctx); err != nil {
		if closeErr := database.Close(); closeErr != nil {
			return nil, fmt.Errorf("ping sqlite database: %w", errors.Join(err, closeErr))
		}
		return nil, fmt.Errorf("ping sqlite database: %w", err)
	}
	if err := applyMigrations(ctx, database); err != nil {
		if closeErr := database.Close(); closeErr != nil {
			return nil, fmt.Errorf("apply sqlite migrations: %w", errors.Join(err, closeErr))
		}
		return nil, err
	}
	if err := seedMasterDictionaryEntries(ctx, database, seed); err != nil {
		if closeErr := database.Close(); closeErr != nil {
			return nil, fmt.Errorf("seed sqlite master dictionary: %w", errors.Join(err, closeErr))
		}
		return nil, err
	}
	return database, nil
}

func ensureDatabasePath(databasePath string) (string, error) {
	if strings.TrimSpace(databasePath) == "" {
		return "", fmt.Errorf("database path is required")
	}

	resolvedPath, err := filepath.Abs(databasePath)
	if err != nil {
		return "", fmt.Errorf("resolve sqlite database path: %w", err)
	}
	if err := os.MkdirAll(filepath.Dir(resolvedPath), 0o750); err != nil {
		return "", fmt.Errorf("create sqlite database directory: %w", err)
	}
	return resolvedPath, nil
}

func buildSQLiteDSN(databasePath string) string {
	query := url.Values{}
	query.Add("_pragma", "foreign_keys(1)")
	query.Add("_pragma", "journal_mode(WAL)")
	query.Add("_pragma", "busy_timeout(5000)")
	query.Set("_time_format", "sqlite")
	query.Set("_txlock", "immediate")
	query.Set("_timezone", "UTC")
	return (&url.URL{
		Scheme:   "file",
		Path:     databasePath,
		RawQuery: query.Encode(),
	}).String()
}

func applyMigrations(ctx context.Context, database *sqlx.DB) error {
	files, err := fs.Glob(migrationFiles, "migrations/*.sql")
	if err != nil {
		return fmt.Errorf("list sqlite migrations: %w", err)
	}
	sort.Strings(files)

	for _, migrationPath := range files {
		migrationSQL, err := fs.ReadFile(migrationFiles, migrationPath)
		if err != nil {
			return fmt.Errorf("read sqlite migration %s: %w", migrationPath, err)
		}
		if strings.TrimSpace(string(migrationSQL)) == "" {
			continue
		}
		if _, err := database.ExecContext(ctx, string(migrationSQL)); err != nil {
			if shouldIgnoreMigrationError(migrationPath, err) {
				continue
			}
			return fmt.Errorf("apply sqlite migration %s: %w", migrationPath, err)
		}
	}
	return nil
}

func shouldIgnoreMigrationError(migrationPath string, err error) bool {
	if migrationPath != "migrations/005_translation_input_source_hash.sql" {
		return false
	}

	return strings.Contains(err.Error(), "duplicate column name")
}

func seedMasterDictionaryEntries(
	ctx context.Context,
	database *sqlx.DB,
	seed []MasterDictionarySeedEntry,
) error {
	if len(seed) == 0 {
		return nil
	}

	var count int
	if err := database.GetContext(ctx, &count, countAllMasterDictionaryEntries); err != nil {
		return fmt.Errorf("count master dictionary seed entries: %w", err)
	}
	if count > 0 {
		return nil
	}

	for _, entry := range seed {
		ts := entry.UpdatedAt.UTC().Format(masterDictionaryTimestampLayout)
		row := masterDictionarySeedRow{
			Source:      entry.Source,
			Translation: entry.Translation,
			Category:    entry.Category,
			Origin:      entry.Origin,
			REC:         entry.REC,
			UpdatedAt:   ts,
			CreatedAt:   ts,
		}
		if _, err := database.NamedExecContext(ctx, insertMasterDictionaryEntry, row); err != nil {
			return fmt.Errorf("seed master dictionary entry: %w", err)
		}
	}
	return nil
}
