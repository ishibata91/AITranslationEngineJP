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

// CountReferenceTranslations は既訳（参照訳）の総件数を返す。翻訳前区間が供給の観測ログへ出す。
// 全件を読まずに数えるため、既訳 corpus が数万件でもログのためだけに corpus を持ち上げない。
func (s *Store) CountReferenceTranslations(ctx context.Context) (int, error) {
	return s.count(ctx, `SELECT COUNT(*) FROM reference_translation`)
}
