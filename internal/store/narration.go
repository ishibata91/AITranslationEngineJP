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
// plugin が空でなければその対象 plugin の行だけに絞る（空なら全 plugin）。
func (s *Store) ListUntranslatedNarrations(ctx context.Context, plugin string) ([]model.Narration, error) {
	if plugin == "" {
		return s.queryNarrations(ctx, `SELECT `+narrationColumns+` FROM narration WHERE status = 0 ORDER BY id`)
	}
	return s.queryNarrations(ctx, `SELECT `+narrationColumns+` FROM narration WHERE status = 0 AND plugin = ? ORDER BY id`, plugin)
}

// NarrationsAfter は id が afterID より大きい叙述文を id 昇順で最大 limit 件返す（keyset ページング用）。
// afterID=0 で先頭から、最後に読んだ id を次回 afterID に渡して次ページを得る。
// plugin が空でなければその対象 plugin の行だけに絞る（空なら全 plugin）。
func (s *Store) NarrationsAfter(ctx context.Context, plugin string, afterID int64, limit int, untranslatedOnly bool) ([]model.Narration, error) {
	statusCondition := ""
	if untranslatedOnly {
		statusCondition = " AND status = 0"
	}
	if plugin == "" {
		return s.queryNarrations(ctx, `SELECT `+narrationColumns+` FROM narration WHERE id > ?`+statusCondition+` ORDER BY id LIMIT ?`, afterID, limit)
	}
	return s.queryNarrations(ctx, `SELECT `+narrationColumns+` FROM narration WHERE plugin = ? AND id > ?`+statusCondition+` ORDER BY id LIMIT ?`, plugin, afterID, limit)
}

// CountNarrations は叙述文の総件数を返す（ページャの総件数表示用）。plugin が空でなければその対象 plugin に絞る。
func (s *Store) CountNarrations(ctx context.Context, plugin string) (int, error) {
	if plugin == "" {
		return s.count(ctx, `SELECT COUNT(*) FROM narration`)
	}
	return s.count(ctx, `SELECT COUNT(*) FROM narration WHERE plugin = ?`, plugin)
}

func (s *Store) queryNarrations(ctx context.Context, query string, args ...any) ([]model.Narration, error) {
	var rows []model.Narration
	if err := s.db.SelectContext(ctx, &rows, query, args...); err != nil {
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
