package repository

import "errors"

var (
	// ErrNotFound は対象レコードが存在しない場合に返す。
	ErrNotFound = errors.New("not found")

	// ErrConflict は Unique 制約違反など衝突が発生した場合に返す。
	ErrConflict = errors.New("conflict")
)
