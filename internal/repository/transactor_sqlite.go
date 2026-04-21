package repository

import (
	"context"
	"fmt"

	"github.com/jmoiron/sqlx"
)

// SQLiteTransactor は Transactor の SQLite 実装。
type SQLiteTransactor struct {
	db *sqlx.DB
}

// NewSQLiteTransactor は Transactor を返す。
func NewSQLiteTransactor(db *sqlx.DB) Transactor {
	return &SQLiteTransactor{db: db}
}

// WithTransaction は fn を単一 DB transaction 内で実行する。
// fn が error を返した場合はロールバックし、その error を返す。
// fn が nil を返した場合はコミットする。
func (t *SQLiteTransactor) WithTransaction(ctx context.Context, fn func(ctx context.Context) error) error {
	tx, err := t.db.BeginTxx(ctx, nil)
	if err != nil {
		return fmt.Errorf("begin transaction: %w", err)
	}
	ctxWithTx := context.WithValue(ctx, TxKey{}, tx)
	if fnErr := fn(ctxWithTx); fnErr != nil {
		if rbErr := tx.Rollback(); rbErr != nil {
			return fmt.Errorf("%w; rollback: %s", fnErr, rbErr.Error())
		}
		return fnErr
	}
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("commit transaction: %w", err)
	}
	return nil
}

// extractTx は ctx に *sqlx.Tx が格納されていればそれを返し、
// なければ db を返す。
func extractTx(ctx context.Context, db *sqlx.DB) sqlx.ExtContext {
	if tx, ok := ctx.Value(TxKey{}).(*sqlx.Tx); ok {
		return tx
	}
	return db
}
