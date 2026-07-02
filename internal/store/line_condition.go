package store

import (
	"context"
	"fmt"
)

// LinkLineConditionsFromStaging は extracted_info_condition（INFO→条件由来の性別）から line_condition を解決する。
// 取込段が line を作った後に 1 度呼ぶ。INFO 由来の line（rec='INFO'）を form_id 一致で staging に結び、
// 条件由来の性別を line.id へ結ぶ。INSERT OR IGNORE で再実行でも冪等（LinkLineSpeakersFromStaging と対称）。
func (s *Store) LinkLineConditionsFromStaging(ctx context.Context) error {
	if _, err := s.db.ExecContext(ctx,
		`INSERT OR IGNORE INTO line_condition (line_id, sex)
		 SELECT l.id, st.sex
		 FROM line l
		 JOIN extracted_info_condition st
		   ON l.plugin = st.info_plugin AND l.form_id = st.info_form_id
		 WHERE l.rec = 'INFO' AND TRIM(st.sex) <> ''`); err != nil {
		return fmt.Errorf("line_condition の解決: %w", err)
	}
	return nil
}

// LoadLineConditions は台詞 id 群に紐づく条件由来の性別を一括で読み、lineID→sex の map を返す。
// 条件で性別を定められた台詞だけ現れる（性別なしの台詞は欠落）。汎用台詞の一人称・語尾の根拠に使う。
func (s *Store) LoadLineConditions(ctx context.Context, lineIDs []int64) (map[int64]string, error) {
	out := make(map[int64]string, len(lineIDs))
	if len(lineIDs) == 0 {
		return out, nil
	}
	for _, chunk := range chunkInt64(lineIDs, speakerChunkSize) {
		marks, args := placeholders(chunk)
		type row struct {
			LineID int64  `db:"line_id"`
			Sex    string `db:"sex"`
		}
		var rows []row
		query := `SELECT line_id, sex FROM line_condition WHERE line_id IN (` + marks + `)`
		if err := s.db.SelectContext(ctx, &rows, query, args...); err != nil {
			return nil, fmt.Errorf("条件由来の性別の一括取得: %w", err)
		}
		for _, r := range rows {
			out[r.LineID] = r.Sex
		}
	}
	return out, nil
}
