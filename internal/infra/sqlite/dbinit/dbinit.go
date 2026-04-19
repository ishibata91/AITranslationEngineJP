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
	countAllMasterDictionaryEntries = `SELECT COUNT(*) FROM master_dictionary_entries;`
	insertMasterDictionaryEntry     = `
INSERT INTO master_dictionary_entries (
  source,
  translation,
  category,
  origin,
  rec,
  edid,
  updated_at
) VALUES (
  :source,
  :translation,
  :category,
  :origin,
  :rec,
  :edid,
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
	Source      string `db:"source"`
	Translation string `db:"translation"`
	Category    string `db:"category"`
	Origin      string `db:"origin"`
	REC         string `db:"rec"`
	EDID        string `db:"edid"`
	UpdatedAt   string `db:"updated_at"`
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
			return fmt.Errorf("apply sqlite migration %s: %w", migrationPath, err)
		}
	}
	return nil
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
		return fmt.Errorf("count sqlite seed target rows: %w", err)
	}
	if count > 0 {
		return nil
	}

	transaction, err := database.BeginTxx(ctx, nil)
	if err != nil {
		return fmt.Errorf("begin sqlite seed transaction: %w", err)
	}

	for _, entry := range seed {
		if _, err := transaction.NamedExecContext(ctx, insertMasterDictionaryEntry, masterDictionarySeedRow{
			Source:      entry.Source,
			Translation: entry.Translation,
			Category:    entry.Category,
			Origin:      entry.Origin,
			REC:         entry.REC,
			EDID:        entry.EDID,
			UpdatedAt:   entry.UpdatedAt.UTC().Format(masterDictionaryTimestampLayout),
		}); err != nil {
			if rollbackErr := transaction.Rollback(); rollbackErr != nil {
				return fmt.Errorf("insert sqlite seed entry: %w", errors.Join(err, rollbackErr))
			}
			return fmt.Errorf("insert sqlite seed entry: %w", err)
		}
	}
	if err := transaction.Commit(); err != nil {
		return fmt.Errorf("commit sqlite seed transaction: %w", err)
	}
	return nil
}
