package store

import (
	"context"
	"fmt"

	"aitranslationenginejp/internal/model"

	"github.com/jmoiron/sqlx"
)

// ListNarrations は叙述文・定型句の全件を id 昇順で返す。取込段の言及検出が本文を走査する入力。
func (s *Store) ListNarrations(ctx context.Context) ([]model.Narration, error) {
	return s.queryNarrations(ctx, `SELECT `+narrationColumns+` FROM narration ORDER BY id`)
}

// ListLines は台詞の全件を id 昇順で返す。取込段の言及検出が本文を走査する入力。
func (s *Store) ListLines(ctx context.Context) ([]model.Line, error) {
	return s.queryLines(ctx, `SELECT `+lineColumns+` FROM line ORDER BY id`)
}

// InsertNarrationMentions は叙述文の本文中言及（e4）を narration_mention へ一括投入する。追加できた件数を返す。
// 排他 2 列（proper_noun_id / master_term_id）は 0 を NULL に写して書き、部分 UNIQUE 索引と
// INSERT OR IGNORE で再検出でも冪等。
func (s *Store) InsertNarrationMentions(ctx context.Context, rows []model.NarrationMention) (int, error) {
	return batchInsert(ctx, s, rows, func(tx *sqlx.Tx, r model.NarrationMention) (int64, error) {
		res, err := tx.ExecContext(ctx,
			`INSERT OR IGNORE INTO narration_mention (narration_id, proper_noun_id, master_term_id)
			 VALUES (?, NULLIF(?, 0), NULLIF(?, 0))`,
			r.NarrationID, r.ProperNounID, r.MasterTermID)
		if err != nil {
			return 0, fmt.Errorf("narration_mention 行の投入: %w", err)
		}
		return res.RowsAffected()
	})
}

// InsertLineMentions は台詞の本文中言及（e5）を line_mention へ一括投入する。追加できた件数を返す。
// InsertNarrationMentions と対称で、再検出でも冪等。
func (s *Store) InsertLineMentions(ctx context.Context, rows []model.LineMention) (int, error) {
	return batchInsert(ctx, s, rows, func(tx *sqlx.Tx, r model.LineMention) (int64, error) {
		res, err := tx.ExecContext(ctx,
			`INSERT OR IGNORE INTO line_mention (line_id, proper_noun_id, master_term_id)
			 VALUES (?, NULLIF(?, 0), NULLIF(?, 0))`,
			r.LineID, r.ProperNounID, r.MasterTermID)
		if err != nil {
			return 0, fmt.Errorf("line_mention 行の投入: %w", err)
		}
		return res.RowsAffected()
	})
}

// LinkNarrationDescribed は叙述文の説明対象（e3）を narration_described へ解決する。追加できた件数を返す。
// 説明対象は「同一レコード（plugin, form_id, rec）の FULL に載る固有名」で、box='叙述文' の行だけが持つ
// （定型句は説明対象を持たない）。extracted_field の FULL を proper_noun (category, source) へ結んで導く。
// PRIMARY KEY(narration_id) と INSERT OR IGNORE で叙述文 1 件あたり 0..1 に保ち、再実行でも冪等。
// SELECT の順序を固定し、万一 FULL が複数あっても先勝ちが決定的になるようにする。
func (s *Store) LinkNarrationDescribed(ctx context.Context) (int, error) {
	res, err := s.db.ExecContext(ctx,
		`INSERT OR IGNORE INTO narration_described (narration_id, proper_noun_id)
		 SELECT n.id, pn.id
		 FROM narration n
		 JOIN record_type_master rtm
		   ON rtm.rec = n.rec AND rtm.field = n.field AND rtm.box = '叙述文'
		 JOIN extracted_field ef
		   ON ef.plugin = n.plugin AND ef.form_id = n.form_id AND ef.rec = n.rec AND ef.field = 'FULL'
		 JOIN proper_noun pn
		   ON pn.category = ef.rec AND pn.source = ef.source
		 ORDER BY n.id, ef.ordinal, pn.id`)
	if err != nil {
		return 0, fmt.Errorf("narration_described の解決: %w", err)
	}
	n, err := res.RowsAffected()
	if err != nil {
		return 0, fmt.Errorf("narration_described 件数の取得: %w", err)
	}
	return int(n), nil
}
