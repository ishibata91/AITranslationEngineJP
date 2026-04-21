package repository

import (
	"context"
	"fmt"

	sqliteinfra "aitranslationenginejp/internal/infra/sqlite/dbinit"

	"github.com/jmoiron/sqlx"
)

// OpenSQLiteDatabase は SQLite データベースを開く。
// bootstrap とそのテストで raw *sqlx.DB が必要な場合に使用する。
func OpenSQLiteDatabase(ctx context.Context, databasePath string) (*sqlx.DB, error) {
	db, err := sqliteinfra.OpenMasterDictionaryDatabase(ctx, databasePath, nil)
	if err != nil {
		return nil, fmt.Errorf("open sqlite database: %w", err)
	}
	return db, nil
}

// NewSQLiteFoundationDataRepositoryFromPath は SQLite データベースを開き、
// FoundationDataRepository とクローザーを返す。
// bootstrap でのワイヤリングで使用する。
func NewSQLiteFoundationDataRepositoryFromPath(
	ctx context.Context,
	databasePath string,
) (FoundationDataRepository, func(context.Context) error, error) {
	db, err := sqliteinfra.OpenMasterDictionaryDatabase(ctx, databasePath, nil)
	if err != nil {
		return nil, nil, fmt.Errorf("open sqlite foundation data database: %w", err)
	}
	closer := func(context.Context) error {
		if closeErr := db.Close(); closeErr != nil {
			return fmt.Errorf("close sqlite foundation data database: %w", closeErr)
		}
		return nil
	}
	return &SQLiteFoundationDataRepository{db: db}, closer, nil
}
