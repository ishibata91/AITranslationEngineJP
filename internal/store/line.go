package store

import (
	"context"
	"fmt"

	"aitranslationenginejp/internal/model"
)

// lineColumns は line の SELECT 列。model.Line の db タグと対応する。1 箇所に集約する。
const lineColumns = `id, source, dest, status, response_order, plugin, form_id, edid, rec, field, ordinal`

// ListUntranslatedLines は未訳（status=0）の台詞を id 昇順で返す。
func (s *Store) ListUntranslatedLines(ctx context.Context) ([]model.Line, error) {
	return s.queryLines(ctx, `SELECT `+lineColumns+` FROM line WHERE status = 0 ORDER BY id`)
}

// LinesAfter は id が afterID より大きい台詞を id 昇順で最大 limit 件返す（keyset ページング用）。
// afterID=0 で先頭から、最後に読んだ id を次回 afterID に渡して次ページを得る。
func (s *Store) LinesAfter(ctx context.Context, afterID int64, limit int) ([]model.Line, error) {
	return s.queryLines(ctx, `SELECT `+lineColumns+` FROM line WHERE id > ? ORDER BY id LIMIT ?`, afterID, limit)
}

// CountLines は台詞の総件数を返す（ページャの総件数表示用）。
func (s *Store) CountLines(ctx context.Context) (int, error) {
	return s.count(ctx, `SELECT COUNT(*) FROM line`)
}

func (s *Store) queryLines(ctx context.Context, query string, args ...any) ([]model.Line, error) {
	var rows []model.Line
	if err := s.db.SelectContext(ctx, &rows, query, args...); err != nil {
		return nil, fmt.Errorf("line の取得: %w", err)
	}
	return rows, nil
}

// UpdateLineDest は訳文と訳状態を書き戻す。
func (s *Store) UpdateLineDest(ctx context.Context, id int64, dest string, status int) error {
	if _, err := s.db.ExecContext(ctx,
		`UPDATE line SET dest = ?, status = ? WHERE id = ?`, dest, status, id); err != nil {
		return fmt.Errorf("line dest の更新: %w", err)
	}
	return nil
}
