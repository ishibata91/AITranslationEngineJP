package store

import (
	"context"
	"fmt"

	"aitranslationenginejp/internal/model"
)

// lineColumns は line の SELECT 列。model.Line の db タグと対応する。1 箇所に集約する。
const lineColumns = `id, source, dest, status, response_order, plugin, form_id, edid, rec, field, ordinal, emotion_type`

// ListUntranslatedLines は未訳（status=0）の台詞を id 昇順で返す。
// plugin が空でなければその対象 plugin の行だけに絞る（空なら全 plugin）。
func (s *Store) ListUntranslatedLines(ctx context.Context, plugin string) ([]model.Line, error) {
	if plugin == "" {
		return s.queryLines(ctx, `SELECT `+lineColumns+` FROM line WHERE status = 0 ORDER BY id`)
	}
	return s.queryLines(ctx, `SELECT `+lineColumns+` FROM line WHERE status = 0 AND plugin = ? ORDER BY id`, plugin)
}

// LinesAfter は id が afterID より大きい台詞を id 昇順で最大 limit 件返す（keyset ページング用）。
// afterID=0 で先頭から、最後に読んだ id を次回 afterID に渡して次ページを得る。
// plugin が空でなければその対象 plugin の行だけに絞る（空なら全 plugin）。
func (s *Store) LinesAfter(ctx context.Context, plugin string, afterID int64, limit int) ([]model.Line, error) {
	if plugin == "" {
		return s.queryLines(ctx, `SELECT `+lineColumns+` FROM line WHERE id > ? ORDER BY id LIMIT ?`, afterID, limit)
	}
	return s.queryLines(ctx, `SELECT `+lineColumns+` FROM line WHERE plugin = ? AND id > ? ORDER BY id LIMIT ?`, plugin, afterID, limit)
}

// CountLines は台詞の総件数を返す（ページャの総件数表示用）。plugin が空でなければその対象 plugin に絞る。
func (s *Store) CountLines(ctx context.Context, plugin string) (int, error) {
	if plugin == "" {
		return s.count(ctx, `SELECT COUNT(*) FROM line`)
	}
	return s.count(ctx, `SELECT COUNT(*) FROM line WHERE plugin = ?`, plugin)
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
