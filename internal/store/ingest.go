package store

import (
	"context"
	"fmt"

	"aitranslationenginejp/internal/model"

	"github.com/jmoiron/sqlx"
)

// extractedFieldColumns は extracted_field の SELECT 列。model.ExtractedField の db タグと対応する。
const extractedFieldColumns = `id, plugin, form_id, edid, rec, field, ordinal, source`

// ListExtractedFields は C# が素朴吸い出しした原文を全件返す。取込段が record_type_master で振り分ける入力。
func (s *Store) ListExtractedFields(ctx context.Context) ([]model.ExtractedField, error) {
	var rows []model.ExtractedField
	if err := s.db.SelectContext(ctx, &rows,
		`SELECT `+extractedFieldColumns+` FROM extracted_field ORDER BY id`); err != nil {
		return nil, fmt.Errorf("extracted_field の取得: %w", err)
	}
	return rows, nil
}

// IngestNarrations は取込段が振り分けた叙述文・定型句を narration へ一括投入する。実際に追加した件数を返す。
// UNIQUE(plugin, form_id, rec, field, ordinal) と INSERT OR IGNORE で再取込でも冪等。
// style には割り当て directive のキーを入れる（人間が中身を見るための由来。翻訳時の正は record_type_master）。
func (s *Store) IngestNarrations(ctx context.Context, rows []model.Narration) (int, error) {
	return batchInsert(ctx, s, rows, func(tx *sqlx.Tx, r model.Narration) (int64, error) {
		res, err := tx.ExecContext(ctx,
			`INSERT OR IGNORE INTO narration (source, style, plugin, form_id, edid, rec, field, ordinal)
			 VALUES (?, ?, ?, ?, ?, ?, ?, ?)`,
			r.Source, r.Style, r.Plugin, r.FormID, r.EDID, r.Rec, r.Field, r.Ordinal)
		if err != nil {
			return 0, fmt.Errorf("narration 行の投入: %w", err)
		}
		return res.RowsAffected()
	})
}

// IngestProperNouns は取込段が振り分けた固有名を proper_noun へ一括投入する。実際に追加した件数を返す。
// UNIQUE(category, source) と INSERT OR IGNORE で同一固有名は 1 つにまとまる（重複排除）。
func (s *Store) IngestProperNouns(ctx context.Context, rows []model.ProperNoun) (int, error) {
	return batchInsert(ctx, s, rows, func(tx *sqlx.Tx, r model.ProperNoun) (int64, error) {
		res, err := tx.ExecContext(ctx,
			`INSERT OR IGNORE INTO proper_noun (source, category) VALUES (?, ?)`,
			r.Source, r.Category)
		if err != nil {
			return 0, fmt.Errorf("proper_noun 行の投入: %w", err)
		}
		return res.RowsAffected()
	})
}

// IngestLines は取込段が振り分けた台詞を line へ一括投入する。実際に追加した件数を返す。
// 話者連関（line_speaker）は line 投入後に LinkLineSpeakersFromStaging で解決する。
func (s *Store) IngestLines(ctx context.Context, rows []model.Line) (int, error) {
	return batchInsert(ctx, s, rows, func(tx *sqlx.Tx, r model.Line) (int64, error) {
		res, err := tx.ExecContext(ctx,
			`INSERT OR IGNORE INTO line (source, response_order, plugin, form_id, edid, rec, field, ordinal)
			 VALUES (?, ?, ?, ?, ?, ?, ?, ?)`,
			r.Source, r.ResponseOrder, r.Plugin, r.FormID, r.EDID, r.Rec, r.Field, r.Ordinal)
		if err != nil {
			return 0, fmt.Errorf("line 行の投入: %w", err)
		}
		return res.RowsAffected()
	})
}

// LinkLineSpeakersFromStaging は extracted_info_speaker（INFO→speaker.id の橋渡し）から line_speaker（e6）を解決する。
// 取込段が line を作った後に 1 度呼ぶ。INFO 由来の line（rec='INFO'）を form_id 一致で staging に結び、
// 解決済み speaker.id を line.id へ結ぶ。INSERT OR IGNORE で再実行でも冪等。
func (s *Store) LinkLineSpeakersFromStaging(ctx context.Context) error {
	if _, err := s.db.ExecContext(ctx,
		`INSERT OR IGNORE INTO line_speaker (line_id, speaker_id)
		 SELECT l.id, st.speaker_id
		 FROM line l
		 JOIN extracted_info_speaker st
		   ON l.plugin = st.info_plugin AND l.form_id = st.info_form_id
		 WHERE l.rec = 'INFO'`); err != nil {
		return fmt.Errorf("line_speaker の解決: %w", err)
	}
	return nil
}

// batchInsert は rows を 1 トランザクションで insertOne に順に通し、追加できた合計件数を返す。
// 取込段の narration/proper_noun/line 投入で共通に使う（INSERT OR IGNORE で冪等）。
func batchInsertGeneric[T any](ctx context.Context, db *sqlx.DB, rows []T, insertOne func(*sqlx.Tx, T) (int64, error)) (int, error) {
	if len(rows) == 0 {
		return 0, nil
	}
	tx, err := db.BeginTxx(ctx, nil)
	if err != nil {
		return 0, fmt.Errorf("取込のトランザクション開始: %w", err)
	}
	defer tx.Rollback() //nolint:errcheck // commit 済みなら no-op。失敗時の後始末用。
	written := 0
	for _, r := range rows {
		n, err := insertOne(tx, r)
		if err != nil {
			return 0, fmt.Errorf("取込の書き込み: %w", err)
		}
		written += int(n)
	}
	if err := tx.Commit(); err != nil {
		return 0, fmt.Errorf("取込のコミット: %w", err)
	}
	return written, nil
}

// batchInsert は s.db を batchInsertGeneric へ渡す薄いラッパ。
func batchInsert[T any](ctx context.Context, s *Store, rows []T, insertOne func(*sqlx.Tx, T) (int64, error)) (int, error) {
	return batchInsertGeneric(ctx, s.db, rows, insertOne)
}
