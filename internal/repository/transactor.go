package repository

import "context"

// Transactor は複数 repository を同一 transaction で扱う境界を定義する。
// 実装は infra/sqlite 層に置く。
type Transactor interface {
	// WithTransaction は fn を単一 DB transaction 内で実行する。
	// fn が error を返した場合はロールバックし、その error を返す。
	// fn が nil を返した場合はコミットする。
	WithTransaction(ctx context.Context, fn func(ctx context.Context) error) error
}
