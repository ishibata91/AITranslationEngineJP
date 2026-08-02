package main

import (
	"context"
	_ "embed"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/jmoiron/sqlx"
	_ "modernc.org/sqlite"
)

//go:embed schema.sql
var dictionarySchema string

var errRevisionConflict = errors.New("辞書項目は取得後に更新されている")

type store struct {
	db *sqlx.DB
}

type entry struct {
	ID        int64    `db:"id" json:"id"`
	Source    string   `db:"source" json:"source"`
	Dest      string   `db:"dest" json:"dest"`
	Category  string   `db:"category" json:"category"`
	Revision  int64    `db:"revision" json:"revision"`
	Sources   []origin `json:"sources,omitempty"`
	MatchKind string   `json:"match_kind,omitempty"`
}

type origin struct {
	Kind      string `db:"kind" json:"kind"`
	Reference string `db:"reference" json:"reference"`
}

func openStore(path string) (*store, error) {
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return nil, fmt.Errorf("辞書DBのdirectory作成: %w", err)
	}
	db, err := sqlx.Open("sqlite", path)
	if err != nil {
		return nil, fmt.Errorf("辞書DBを開く: %w", err)
	}
	db.SetMaxOpenConns(1)
	if _, err := db.Exec(`PRAGMA busy_timeout = 5000; PRAGMA foreign_keys = ON;`); err != nil {
		db.Close()
		return nil, fmt.Errorf("辞書DBの接続設定: %w", err)
	}
	if _, err := db.Exec(dictionarySchema); err != nil {
		db.Close()
		return nil, fmt.Errorf("辞書DBのschema適用: %w", err)
	}
	return &store{db: db}, nil
}

func (s *store) close() error {
	return s.db.Close()
}

func (s *store) get(ctx context.Context, id int64) (entry, error) {
	var e entry
	if err := s.db.GetContext(ctx, &e,
		`SELECT id, source, dest, category, revision FROM dictionary_entry WHERE id = ?`, id); err != nil {
		return entry{}, fmt.Errorf("辞書項目 %d の取得: %w", id, err)
	}
	if err := s.db.SelectContext(ctx, &e.Sources,
		`SELECT kind, reference FROM dictionary_entry_source WHERE entry_id = ? ORDER BY kind, reference`, id); err != nil {
		return entry{}, fmt.Errorf("辞書項目 %d の出どころ取得: %w", id, err)
	}
	return e, nil
}

func (s *store) add(ctx context.Context, source, dest, category string) (entry, error) {
	source, dest, category = strings.TrimSpace(source), strings.TrimSpace(dest), strings.TrimSpace(category)
	if source == "" || dest == "" {
		return entry{}, errors.New("原語と訳語は空にできない")
	}
	tx, err := s.db.BeginTxx(ctx, nil)
	if err != nil {
		return entry{}, fmt.Errorf("辞書項目追加のtransaction開始: %w", err)
	}
	defer tx.Rollback() //nolint:errcheck
	res, err := tx.ExecContext(ctx,
		`INSERT INTO dictionary_entry (source, dest, category) VALUES (?, ?, ?)`, source, dest, category)
	if err != nil {
		return entry{}, fmt.Errorf("辞書項目の追加: %w", err)
	}
	id, err := res.LastInsertId()
	if err != nil {
		return entry{}, fmt.Errorf("追加した辞書項目のid取得: %w", err)
	}
	if _, err := tx.ExecContext(ctx,
		`INSERT INTO dictionary_entry_source (entry_id, kind, reference) VALUES (?, 'manual', ?)`, id, fmt.Sprintf("entry:%d", id)); err != nil {
		return entry{}, fmt.Errorf("辞書項目の出どころ追加: %w", err)
	}
	if err := tx.Commit(); err != nil {
		return entry{}, fmt.Errorf("辞書項目追加のcommit: %w", err)
	}
	return s.get(ctx, id)
}

func (s *store) update(ctx context.Context, id, revision int64, source, dest, category string) (entry, error) {
	source, dest, category = strings.TrimSpace(source), strings.TrimSpace(dest), strings.TrimSpace(category)
	if source == "" || dest == "" {
		return entry{}, errors.New("原語と訳語は空にできない")
	}
	res, err := s.db.ExecContext(ctx, `
		UPDATE dictionary_entry
		SET source = ?, dest = ?, category = ?, revision = revision + 1, updated_at = ?
		WHERE id = ? AND revision = ?`, source, dest, category, time.Now().UTC().Format(time.RFC3339Nano), id, revision)
	if err != nil {
		return entry{}, fmt.Errorf("辞書項目 %d の更新: %w", id, err)
	}
	n, err := res.RowsAffected()
	if err != nil {
		return entry{}, fmt.Errorf("辞書項目 %d の更新件数取得: %w", id, err)
	}
	if n == 0 {
		return entry{}, errRevisionConflict
	}
	return s.get(ctx, id)
}

type status struct {
	Entries    int            `json:"entries"`
	Origins    map[string]int `json:"origins"`
	FTSEntries int            `json:"fts_entries"`
	Indexes    []string       `json:"indexes"`
}

func (s *store) status(ctx context.Context) (status, error) {
	var out status
	if err := s.db.GetContext(ctx, &out.Entries, `SELECT COUNT(*) FROM dictionary_entry`); err != nil {
		return status{}, fmt.Errorf("辞書項目数の取得: %w", err)
	}
	if err := s.db.GetContext(ctx, &out.FTSEntries, `SELECT COUNT(*) FROM dictionary_entry_fts`); err != nil {
		return status{}, fmt.Errorf("検索用index件数の取得: %w", err)
	}
	type originCount struct {
		Kind  string `db:"kind"`
		Count int    `db:"count"`
	}
	var counts []originCount
	if err := s.db.SelectContext(ctx, &counts,
		`SELECT kind, COUNT(*) AS count FROM dictionary_entry_source GROUP BY kind ORDER BY kind`); err != nil {
		return status{}, fmt.Errorf("出どころ別件数の取得: %w", err)
	}
	out.Origins = make(map[string]int, len(counts))
	for _, c := range counts {
		out.Origins[c.Kind] = c.Count
	}
	if err := s.db.SelectContext(ctx, &out.Indexes,
		`SELECT name FROM sqlite_schema WHERE type IN ('index', 'table') AND (name LIKE 'dictionary_entry_%_idx' OR name = 'dictionary_entry_fts') ORDER BY name`); err != nil {
		return status{}, fmt.Errorf("index状態の取得: %w", err)
	}
	return out, nil
}
