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

// InsertDerivedTerms は派生した固有名（人名の部分形）を master_term へ追記する。実際に追加した件数を返す。
// UNIQUE(category, source) と INSERT OR IGNORE で二重追記を防ぐ（再 Run でも冪等）。
func (s *Store) InsertDerivedTerms(ctx context.Context, terms []model.MasterTerm) (int, error) {
	if len(terms) == 0 {
		return 0, nil
	}
	tx, err := s.db.BeginTxx(ctx, nil)
	if err != nil {
		return 0, fmt.Errorf("派生辞書のトランザクション開始: %w", err)
	}
	defer tx.Rollback() //nolint:errcheck // commit 済みなら no-op。失敗時の後始末用。
	written := 0
	for _, t := range terms {
		res, err := tx.ExecContext(ctx,
			`INSERT OR IGNORE INTO master_term (source, dest, category) VALUES (?, ?, ?)`,
			t.Source, t.Dest, t.Category)
		if err != nil {
			return 0, fmt.Errorf("派生 %q の書き込み: %w", t.Source, err)
		}
		n, _ := res.RowsAffected()
		written += int(n)
	}
	if err := tx.Commit(); err != nil {
		return 0, fmt.Errorf("派生辞書のコミット: %w", err)
	}
	return written, nil
}
