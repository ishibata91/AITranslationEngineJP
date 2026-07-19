package store

import (
	"context"
	"fmt"

	"aitranslationenginejp/internal/model"
)

// referenceTranslationColumns は reference_translation の SELECT 列。model.ReferenceTranslation の db タグと対応する。
const referenceTranslationColumns = `rec, field, source, dest`

// ListReferenceTranslations は参照訳（既存訳）の全件を返す。engine が完全一致置換の照合表を組むために読む。
// 件数は既訳の corpus 規模（base ゲームで数万件）で、ページングせず一括で読む（master_term と同じ扱い）。
func (s *Store) ListReferenceTranslations(ctx context.Context) ([]model.ReferenceTranslation, error) {
	var rows []model.ReferenceTranslation
	if err := s.db.SelectContext(ctx, &rows,
		`SELECT `+referenceTranslationColumns+` FROM reference_translation ORDER BY rec, field, source`); err != nil {
		return nil, fmt.Errorf("reference_translation の取得: %w", err)
	}
	return rows, nil
}

// InsertReferenceTranslations は取り込んだ参照訳を reference_translation へ追記する。実際に追加した件数を返す。
// UNIQUE(rec, field, source) と INSERT OR IGNORE で二重追記を防ぐ（同一 XML を毎 Run 取り込んでも冪等）。
// 同一 (rec, field, source) に異なる既訳が来た場合は先勝ち（後続を IGNORE）。
func (s *Store) InsertReferenceTranslations(ctx context.Context, refs []model.ReferenceTranslation) (int, error) {
	if len(refs) == 0 {
		return 0, nil
	}
	tx, err := s.db.BeginTxx(ctx, nil)
	if err != nil {
		return 0, fmt.Errorf("参照訳のトランザクション開始: %w", err)
	}
	defer tx.Rollback() //nolint:errcheck // commit 済みなら no-op。失敗時の後始末用。
	written := 0
	for _, r := range refs {
		res, err := tx.ExecContext(ctx,
			`INSERT OR IGNORE INTO reference_translation (rec, field, source, dest) VALUES (?, ?, ?, ?)`,
			r.Rec, r.Field, r.Source, r.Dest)
		if err != nil {
			return 0, fmt.Errorf("参照訳 %q の書き込み: %w", r.Source, err)
		}
		n, _ := res.RowsAffected()
		written += int(n)
	}
	if err := tx.Commit(); err != nil {
		return 0, fmt.Errorf("参照訳のコミット: %w", err)
	}
	return written, nil
}
