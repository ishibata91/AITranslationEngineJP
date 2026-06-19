package store

import (
	"context"
	"fmt"

	"aitranslationenginejp/internal/model"
)

// masterTermColumns は master_term の SELECT 列。model.MasterTerm の db タグと対応する。
const masterTermColumns = `id, source, dest, category`

// ListMasterTerms はマスター辞書の全件を返す。engine が固有名の機械置換辞書を組むために読む。
// 件数は固有名の語彙数（数千〜数万）で、ページングせず一括で読む。
func (s *Store) ListMasterTerms(ctx context.Context) ([]model.MasterTerm, error) {
	var rows []model.MasterTerm
	if err := s.db.SelectContext(ctx, &rows,
		`SELECT `+masterTermColumns+` FROM master_term ORDER BY id`); err != nil {
		return nil, fmt.Errorf("master_term の取得: %w", err)
	}
	return rows, nil
}
