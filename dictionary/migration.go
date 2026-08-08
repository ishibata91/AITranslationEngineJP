package main

import (
	"context"
	"fmt"
	"strings"

	"github.com/jmoiron/sqlx"
)

type legacyEntry struct {
	ID       int64  `db:"id"`
	Source   string `db:"source"`
	Dest     string `db:"dest"`
	Category string `db:"category"`
	Revision int64  `db:"revision"`
}

type legacyOrigin struct {
	Kind      string `db:"kind"`
	Reference string `db:"reference"`
}

func (s *store) migrateLegacy(ctx context.Context) error {
	entries, err := s.legacyEntries(ctx)
	if err != nil || len(entries) == 0 {
		return err
	}
	tx, err := s.db.BeginTxx(ctx, nil)
	if err != nil {
		return fmt.Errorf("旧辞書移行transaction開始: %w", err)
	}
	defer tx.Rollback() //nolint:errcheck
	for _, entry := range entries {
		if err := migrateLegacyEntry(ctx, tx, entry); err != nil {
			return err
		}
	}
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("旧辞書移行commit: %w", err)
	}
	return nil
}

func (s *store) legacyEntries(ctx context.Context) ([]legacyEntry, error) {
	var legacyExists int
	if err := s.db.GetContext(ctx, &legacyExists, `
		SELECT COUNT(*) FROM sqlite_schema WHERE type = 'table' AND name = 'dictionary_entry'`); err != nil {
		return nil, fmt.Errorf("旧schemaの確認: %w", err)
	}
	if legacyExists == 0 {
		return nil, nil
	}
	var currentSenses int
	if err := s.db.GetContext(ctx, &currentSenses, `SELECT COUNT(*) FROM dictionary_sense`); err != nil {
		return nil, fmt.Errorf("新schemaの移行済み確認: %w", err)
	}
	if currentSenses > 0 {
		return nil, nil
	}
	var entries []legacyEntry
	if err := s.db.SelectContext(ctx, &entries, `
		SELECT id, source, dest, category, revision FROM dictionary_entry ORDER BY id`); err != nil {
		return nil, fmt.Errorf("旧辞書項目の取得: %w", err)
	}
	return entries, nil
}

func migrateLegacyEntry(ctx context.Context, tx *sqlx.Tx, entry legacyEntry) error {
	termID, err := ensureTerm(ctx, tx, entry.Source)
	if err != nil {
		return err
	}
	if _, err := tx.ExecContext(ctx, `
		INSERT INTO dictionary_sense
		    (id, term_id, dest, classification_status, revision)
		VALUES (?, ?, ?, 'unclassified', ?)`, entry.ID, termID, entry.Dest, entry.Revision); err != nil {
		return fmt.Errorf("旧辞書項目 %d の意味移行: %w", entry.ID, err)
	}
	var origins []legacyOrigin
	if err := tx.SelectContext(ctx, &origins, `
		SELECT kind, reference FROM dictionary_entry_source
		WHERE entry_id = ? ORDER BY kind, reference`, entry.ID); err != nil {
		return fmt.Errorf("旧辞書項目 %d の出どころ取得: %w", entry.ID, err)
	}
	if len(origins) == 0 {
		origins = []legacyOrigin{{Kind: "legacy", Reference: fmt.Sprintf("entry:%d", entry.ID)}}
	}
	category, derivation := splitCategory(entry.Category)
	for _, origin := range origins {
		if _, err := tx.ExecContext(ctx, `
			INSERT INTO dictionary_occurrence
			    (term_id, sense_id, observed_dest, skyrim_category,
			     origin_kind, origin_reference, derivation_kind)
			VALUES (?, ?, ?, ?, ?, ?, ?)`,
			termID, entry.ID, entry.Dest, category,
			origin.Kind, origin.Reference, derivation); err != nil {
			return fmt.Errorf("旧辞書項目 %d の使用箇所移行: %w", entry.ID, err)
		}
	}
	return nil
}

func ensureTerm(ctx context.Context, tx *sqlx.Tx, source string) (int64, error) {
	if _, err := tx.ExecContext(ctx, `
		INSERT INTO dictionary_term (source) VALUES (?)
		ON CONFLICT(source) DO NOTHING`, source); err != nil {
		return 0, fmt.Errorf("英語表記の追加: %w", err)
	}
	var termID int64
	if err := tx.GetContext(ctx, &termID, `SELECT id FROM dictionary_term WHERE source = ?`, source); err != nil {
		return 0, fmt.Errorf("英語表記の取得: %w", err)
	}
	return termID, nil
}

func splitCategory(category string) (string, string) {
	if derivation, ok := strings.CutPrefix(category, "derive:"); ok {
		return "", derivation
	}
	return category, ""
}
