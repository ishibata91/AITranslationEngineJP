// Package store は SQLite への薄いデータアクセスを持つ。sqlx を使い、driver 固有 API を上位へ漏らさない。
package store

import (
	"context"
	"fmt"

	"aitranslationenginejp/db"
	"aitranslationenginejp/internal/model"

	"github.com/jmoiron/sqlx"
	_ "modernc.org/sqlite" // SQLite driver（pure-Go、CGO 不要）を "sqlite" 名で登録する。
)

// Store は中心 DB への接続を持つ。
type Store struct {
	db *sqlx.DB
}

// Open は指定パスの SQLite を開き、schema migration（db パッケージの責務）を適用する。
func Open(path string) (*Store, error) {
	conn, err := sqlx.Connect("sqlite", path)
	if err != nil {
		return nil, fmt.Errorf("SQLite を開けない: %w", err)
	}
	if err := db.Apply(conn); err != nil {
		_ = conn.Close()
		return nil, fmt.Errorf("schema migration の適用: %w", err)
	}
	return &Store{db: conn}, nil
}

// Close は接続を閉じる。
func (s *Store) Close() error {
	if err := s.db.Close(); err != nil {
		return fmt.Errorf("接続のクローズ: %w", err)
	}
	return nil
}

// narrationColumns は narration の SELECT 列。model.Narration の db タグと対応する。1 箇所に集約する。
const narrationColumns = `id, source, dest, status, style, plugin, form_id, edid, rec, field, ordinal`

// ListUntranslatedNarrations は未訳（status=0）の叙述文を id 昇順で返す。
func (s *Store) ListUntranslatedNarrations(ctx context.Context) ([]model.Narration, error) {
	return s.queryNarrations(ctx, `SELECT `+narrationColumns+` FROM narration WHERE status = 0 ORDER BY id`)
}

// ListNarrations は全ての叙述文を id 昇順で返す（画面の結果一覧表示用）。
func (s *Store) ListNarrations(ctx context.Context) ([]model.Narration, error) {
	return s.queryNarrations(ctx, `SELECT `+narrationColumns+` FROM narration ORDER BY id`)
}

func (s *Store) queryNarrations(ctx context.Context, query string) ([]model.Narration, error) {
	var rows []model.Narration
	if err := s.db.SelectContext(ctx, &rows, query); err != nil {
		return nil, fmt.Errorf("narration の取得: %w", err)
	}
	return rows, nil
}

// UpdateNarrationDest は訳文と訳状態を書き戻す。
func (s *Store) UpdateNarrationDest(ctx context.Context, id int64, dest string, status int) error {
	_, err := s.db.ExecContext(ctx,
		`UPDATE narration SET dest = ?, status = ? WHERE id = ?`, dest, status, id)
	if err != nil {
		return fmt.Errorf("narration dest の更新: %w", err)
	}
	return nil
}
