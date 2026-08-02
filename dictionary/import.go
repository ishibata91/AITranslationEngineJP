package main

import (
	"context"
	"fmt"
	"net/url"
	"path/filepath"
	"strconv"

	"github.com/jmoiron/sqlx"
	_ "modernc.org/sqlite"
)

type importResult struct {
	Read      int `json:"read"`
	Created   int `json:"created"`
	Updated   int `json:"updated"`
	Unchanged int `json:"unchanged"`
}

type masterTerm struct {
	ID       int64  `db:"id"`
	Source   string `db:"source"`
	Dest     string `db:"dest"`
	Category string `db:"category"`
}

func importMasterTerms(ctx context.Context, sourcePath string, destination *store) (importResult, error) {
	abs, err := filepath.Abs(sourcePath)
	if err != nil {
		return importResult{}, fmt.Errorf("中心DBのpath解決: %w", err)
	}
	dsn := (&url.URL{Scheme: "file", Path: abs}).String() + "?mode=ro"
	source, err := sqlx.Open("sqlite", dsn)
	if err != nil {
		return importResult{}, fmt.Errorf("中心DBを読み取り専用で開く: %w", err)
	}
	defer source.Close()
	rows, err := source.QueryxContext(ctx,
		`SELECT id, source, dest, category FROM master_term ORDER BY id`)
	if err != nil {
		return importResult{}, fmt.Errorf("中心DBのmaster_term取得: %w", err)
	}
	defer rows.Close()

	tx, err := destination.db.BeginTxx(ctx, nil)
	if err != nil {
		return importResult{}, fmt.Errorf("移植transaction開始: %w", err)
	}
	defer tx.Rollback() //nolint:errcheck

	var result importResult
	for rows.Next() {
		var term masterTerm
		if err := rows.StructScan(&term); err != nil {
			return importResult{}, fmt.Errorf("master_termの読み取り: %w", err)
		}
		result.Read++
		reference := strconv.FormatInt(term.ID, 10)
		var entryID int64
		err := tx.GetContext(ctx, &entryID,
			`SELECT entry_id FROM dictionary_entry_source WHERE kind = 'master_term' AND reference = ?`, reference)
		switch {
		case err == nil:
			res, updateErr := tx.ExecContext(ctx, `
				UPDATE dictionary_entry
				SET source = ?, dest = ?, category = ?, revision = revision + 1, updated_at = CURRENT_TIMESTAMP
				WHERE id = ? AND (source <> ? OR dest <> ? OR category <> ?)`,
				term.Source, term.Dest, term.Category, entryID, term.Source, term.Dest, term.Category)
			if updateErr != nil {
				return importResult{}, fmt.Errorf("master_term %d の更新: %w", term.ID, updateErr)
			}
			n, _ := res.RowsAffected()
			if n == 0 {
				result.Unchanged++
			} else {
				result.Updated++
			}
		case isNoRows(err):
			res, insertErr := tx.ExecContext(ctx, `
				INSERT INTO dictionary_entry (source, dest, category)
				VALUES (?, ?, ?)
				ON CONFLICT(source, dest, category) DO NOTHING`, term.Source, term.Dest, term.Category)
			if insertErr != nil {
				return importResult{}, fmt.Errorf("master_term %d の追加: %w", term.ID, insertErr)
			}
			n, _ := res.RowsAffected()
			if n > 0 {
				entryID, err = res.LastInsertId()
				result.Created++
			} else {
				if err = tx.GetContext(ctx, &entryID,
					`SELECT id FROM dictionary_entry WHERE source = ? AND dest = ? AND category = ?`, term.Source, term.Dest, term.Category); err != nil {
					return importResult{}, fmt.Errorf("既存辞書項目の取得: %w", err)
				}
				result.Unchanged++
			}
			if _, err := tx.ExecContext(ctx,
				`INSERT INTO dictionary_entry_source (entry_id, kind, reference) VALUES (?, 'master_term', ?)`, entryID, reference); err != nil {
				return importResult{}, fmt.Errorf("master_term %d の出どころ追加: %w", term.ID, err)
			}
		default:
			return importResult{}, fmt.Errorf("master_term %d の既存確認: %w", term.ID, err)
		}
	}
	if err := rows.Err(); err != nil {
		return importResult{}, fmt.Errorf("master_termの走査: %w", err)
	}
	if err := tx.Commit(); err != nil {
		return importResult{}, fmt.Errorf("移植transactionのcommit: %w", err)
	}
	return result, nil
}
