package store

import (
	"context"
	"fmt"

	"aitranslationenginejp/internal/model"

	"github.com/jmoiron/sqlx"
)

// properNounColumns は proper_noun の SELECT 列。model.ProperNoun の db タグと対応する。
const properNounColumns = `id, plugin, source, category, dest, status, origin`

// 機械派生した人名の部分形を翻訳対象から外す絞り。origin が空の行だけが抽出由来＝翻訳対象。
// 条件をここ 1 箇所で組み、同じ条件文を各 SQL へ書かない。
// translationTargetProperNoun は別名を付けない proper_noun 用（一覧・件数）、
// translationTargetProperNounAsP は別名 p を付けた proper_noun 用（target_plugin の進捗集計）。
const (
	translationTargetProperNoun    = `origin = ''`
	translationTargetProperNounAsP = `p.origin = ''`
)

// ListProperNouns は固有名の訳の単位を全件返す。本文フェーズの機械置換辞書（master_term と合流）に使う。
// 機械派生した人名の部分形も返す。部分形は機械置換辞書の材料そのもので、外すと本文の置換が効かなくなる。
// 言及（注入の忠実な記録）も同じ供給源を読むため、片側だけに効く除外を作らない。
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
// 機械派生した人名の部分形は翻訳対象でないため返さない。
func (s *Store) ListUntranslatedProperNouns(ctx context.Context, plugin string) ([]model.ProperNoun, error) {
	var rows []model.ProperNoun
	query := `SELECT ` + properNounColumns + ` FROM proper_noun
		 WHERE status = 0 AND ` + translationTargetProperNoun + ` ORDER BY id`
	args := []any{}
	if plugin != "" {
		query = `SELECT ` + properNounColumns + ` FROM proper_noun
			 WHERE status = 0 AND plugin = ? AND ` + translationTargetProperNoun + ` ORDER BY id`
		args = []any{plugin}
	}
	if err := s.db.SelectContext(ctx, &rows, query, args...); err != nil {
		return nil, fmt.Errorf("未訳 proper_noun の取得: %w", err)
	}
	return rows, nil
}

// CountProperNouns は翻訳対象の固有名の総件数を返す（結果一覧の総件数表示用）。
// plugin が空でなければその対象 plugin に絞る。機械派生した人名の部分形は数えない。
func (s *Store) CountProperNouns(ctx context.Context, plugin string) (int, error) {
	if plugin == "" {
		return s.count(ctx, `SELECT COUNT(*) FROM proper_noun WHERE `+translationTargetProperNoun)
	}
	return s.count(ctx,
		`SELECT COUNT(*) FROM proper_noun WHERE plugin = ? AND `+translationTargetProperNoun, plugin)
}

// ProperNounsAfter は id が afterID より大きい翻訳対象の固有名を id 昇順で最大 limit 件返す（keyset ページング用）。
// plugin が空でなければその対象 plugin の行だけに絞る（空なら全 plugin）。
// 機械派生した人名の部分形は結果一覧に出さないため返さない（CountProperNouns の総件数と一致させる）。
func (s *Store) ProperNounsAfter(ctx context.Context, plugin string, afterID int64, limit int) ([]model.ProperNoun, error) {
	var rows []model.ProperNoun
	query := `SELECT ` + properNounColumns + ` FROM proper_noun
		 WHERE id > ? AND ` + translationTargetProperNoun + ` ORDER BY id LIMIT ?`
	args := []any{afterID, limit}
	if plugin != "" {
		query = `SELECT ` + properNounColumns + ` FROM proper_noun
			 WHERE plugin = ? AND id > ? AND ` + translationTargetProperNoun + ` ORDER BY id LIMIT ?`
		args = []any{plugin, afterID, limit}
	}
	if err := s.db.SelectContext(ctx, &rows, query, args...); err != nil {
		return nil, fmt.Errorf("proper_noun ページの取得: %w", err)
	}
	return rows, nil
}

// ListConfirmedNPCNames は指定 plugin で訳が確定した NPC の氏名（FULL）・短縮名（SHRT）を、原文の field つきで返す。
// proper_noun は氏名と短縮名を分ける列を持たないため、原文（extracted_field）へ
// (plugin, category = rec, source) で結んで field を取り戻す（mention・export の位置解決と同じ結合）。
// 既訳流用で確定した行と AI 訳で確定した行を区別しない（どちらも実行内で確定した訳）。
// 人名の部分形を派生する入力になるため、並び順を (field, source) で固定して派生結果を決定的にする。
func (s *Store) ListConfirmedNPCNames(ctx context.Context, plugin string) ([]model.ConfirmedName, error) {
	var rows []model.ConfirmedName
	if err := s.db.SelectContext(ctx, &rows,
		`SELECT DISTINCT ef.field AS field, pn.source AS source, pn.dest AS dest
		 FROM proper_noun pn
		 JOIN extracted_field ef
		   ON ef.plugin = pn.plugin AND ef.rec = pn.category AND ef.source = pn.source
		 WHERE pn.plugin = ? AND pn.dest <> '' AND pn.category = 'NPC_' AND ef.field IN ('FULL', 'SHRT')
		 ORDER BY ef.field, pn.source`, plugin); err != nil {
		return nil, fmt.Errorf("確定した NPC 名の取得: %w", err)
	}
	return rows, nil
}

// InsertDerivedProperNouns は派生した人名の部分形を proper_noun へ追記する。実際に追加した件数を返す。
// UNIQUE(plugin, category, source) と INSERT OR IGNORE で二重実行でも増えない（冪等）。
// 横断永続辞書 master_term へは書かない（方針A の不変境界。派生元が実行内の AI 訳を含むため）。
func (s *Store) InsertDerivedProperNouns(ctx context.Context, rows []model.ProperNoun) (int, error) {
	return batchInsert(ctx, s, rows, func(tx *sqlx.Tx, r model.ProperNoun) (int64, error) {
		res, err := tx.ExecContext(ctx,
			`INSERT OR IGNORE INTO proper_noun (plugin, source, category, dest, status, origin) VALUES (?, ?, ?, ?, ?, ?)`,
			r.Plugin, r.Source, r.Category, r.Dest, r.Status, r.Origin)
		if err != nil {
			return 0, fmt.Errorf("派生した人名の部分形の投入: %w", err)
		}
		return res.RowsAffected()
	})
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
