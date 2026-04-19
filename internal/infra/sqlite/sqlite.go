// Package sqlite owns SQLite startup initialization for repo-managed databases.
package sqlite

import (
	"context"
	"time"

	"github.com/jmoiron/sqlx"

	dbinit "aitranslationenginejp/internal/infra/sqlite/dbinit"
)

// masterDictionaryTimestampLayout は sqlite_test.go で参照されるため sqlite パッケージに置く。
const masterDictionaryTimestampLayout = time.RFC3339Nano

// MasterDictionarySeedEntry describes one startup seed row.
// This is a type alias for dbinit.MasterDictionarySeedEntry for backward compatibility.
type MasterDictionarySeedEntry = dbinit.MasterDictionarySeedEntry

// OpenMasterDictionaryDatabase opens the SQLite database, reapplies migrations, and seeds an empty database.
func OpenMasterDictionaryDatabase(
	ctx context.Context,
	databasePath string,
	seed []MasterDictionarySeedEntry,
) (*sqlx.DB, error) {
	return dbinit.OpenMasterDictionaryDatabase(ctx, databasePath, seed)
}
