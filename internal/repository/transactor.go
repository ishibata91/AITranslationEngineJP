package repository

import "context"

// TxKey は context に *sqlx.Tx を格納する際に使うキー型。
// infra/sqlite/transactor.go が WithTransaction 時に ctx に tx を埋め込み、
// repository 層がこのキーで取り出して tx-in-context 操作を行う。
type TxKey struct{}

// Transactor は複数 repository を同一 transaction で扱う境界を定義する。
// 実装は infra/sqlite 層に置く。
type Transactor interface {
	// WithTransaction は fn を単一 DB transaction 内で実行する。
	// fn が error を返した場合はロールバックし、その error を返す。
	// fn が nil を返した場合はコミットする。
	WithTransaction(ctx context.Context, fn func(ctx context.Context) error) error
}
