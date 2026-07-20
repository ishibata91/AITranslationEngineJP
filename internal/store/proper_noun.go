package store

import (
	"context"
	"fmt"

	"aitranslationenginejp/internal/model"
)

// properNounColumns は proper_noun の SELECT 列。model.ProperNoun の db タグと対応する。
const properNounColumns = `id, plugin, source, category, dest, status`

// ListProperNouns は固有名の訳の単位を全件返す。本文フェーズの機械置換辞書（master_term と合流）に使う。
func (s *Store) ListProperNouns(ctx context.Context) ([]model.ProperNoun, error) {
	var rows []model.ProperNoun
	if err := s.db.SelectContext(ctx, &rows,
		`SELECT `+properNounColumns+` FROM proper_noun ORDER BY id`); err != nil {
		return nil, fmt.Errorf("proper_noun の取得: %w", err)
	}
	return rows, nil
}

// ListUntranslatedProperNouns は未訳（status=0）の固有名を id 昇順で返す。固有名フェーズが訳す対象。
// plugin が空でなければその対象 plugin の行だけに絞る（空なら全 plugin）。
func (s *Store) ListUntranslatedProperNouns(ctx context.Context, plugin string) ([]model.ProperNoun, error) {
	var rows []model.ProperNoun
	query := `SELECT ` + properNounColumns + ` FROM proper_noun WHERE status = 0 ORDER BY id`
	args := []any{}
	if plugin != "" {
		query = `SELECT ` + properNounColumns + ` FROM proper_noun WHERE status = 0 AND plugin = ? ORDER BY id`
		args = []any{plugin}
	}
	if err := s.db.SelectContext(ctx, &rows, query, args...); err != nil {
		return nil, fmt.Errorf("未訳 proper_noun の取得: %w", err)
	}
	return rows, nil
}

// CountProperNouns は固有名の総件数を返す（結果一覧の総件数表示用）。plugin が空でなければその対象 plugin に絞る。
func (s *Store) CountProperNouns(ctx context.Context, plugin string) (int, error) {
	if plugin == "" {
		return s.count(ctx, `SELECT COUNT(*) FROM proper_noun`)
	}
	return s.count(ctx, `SELECT COUNT(*) FROM proper_noun WHERE plugin = ?`, plugin)
}

// ProperNounsAfter は id が afterID より大きい固有名を id 昇順で最大 limit 件返す（keyset ページング用）。
// plugin が空でなければその対象 plugin の行だけに絞る（空なら全 plugin）。
func (s *Store) ProperNounsAfter(ctx context.Context, plugin string, afterID int64, limit int) ([]model.ProperNoun, error) {
	var rows []model.ProperNoun
	query := `SELECT ` + properNounColumns + ` FROM proper_noun WHERE id > ? ORDER BY id LIMIT ?`
	args := []any{afterID, limit}
	if plugin != "" {
		query = `SELECT ` + properNounColumns + ` FROM proper_noun WHERE plugin = ? AND id > ? ORDER BY id LIMIT ?`
		args = []any{plugin, afterID, limit}
	}
	if err := s.db.SelectContext(ctx, &rows, query, args...); err != nil {
		return nil, fmt.Errorf("proper_noun ページの取得: %w", err)
	}
	return rows, nil
}

// UpdateProperNounDest は固有名の訳文と訳状態を書き戻す。
// 既訳（master_term）流用も AI 訳も同じ宛先（proper_noun）へ書く。master_term へは書かない（方針A の不変境界）。
func (s *Store) UpdateProperNounDest(ctx context.Context, id int64, dest string, status int) error {
	if _, err := s.db.ExecContext(ctx,
		`UPDATE proper_noun SET dest = ?, status = ? WHERE id = ?`, dest, status, id); err != nil {
		return fmt.Errorf("proper_noun dest の更新: %w", err)
	}
	return nil
}
