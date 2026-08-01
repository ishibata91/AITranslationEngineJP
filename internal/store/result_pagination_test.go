package store

import (
	"context"
	"path/filepath"
	"testing"
)

// R-3-1, R-3-3: 結果ページの各区間は untranslatedOnly=true の場合だけ status=0 へ絞る。
// 固有名の機械派生行は、絞り込みの有無にかかわらず結果ページへ含めない。
func TestResultPagesFilterUntranslatedRows(t *testing.T) {
	dbPath := filepath.Join(t.TempDir(), "central.sqlite3")
	s, err := Open(dbPath)
	if err != nil {
		t.Fatalf("Open: %v", err)
	}
	defer func() { _ = s.Close() }()

	execSQL(t, dbPath, `INSERT INTO narration (plugin, form_id, edid, rec, field, ordinal, source, dest, status) VALUES
		('A.esp', '1', 'N1', 'BOOK', 'DESC', 0, 'n pending', '', 0),
		('A.esp', '2', 'N2', 'BOOK', 'DESC', 0, 'n done', '訳', 3)`)
	execSQL(t, dbPath, `INSERT INTO line (plugin, form_id, edid, rec, field, ordinal, source, dest, status) VALUES
		('A.esp', '3', 'L1', 'INFO', 'NAM1', 0, 'l pending', '', 0),
		('A.esp', '4', 'L2', 'INFO', 'NAM1', 0, 'l done', '訳', 3)`)
	execSQL(t, dbPath, `INSERT INTO proper_noun (plugin, source, category, dest, status, origin) VALUES
		('A.esp', 'p pending', 'NPC_', '', 0, ''),
		('A.esp', 'p done', 'NPC_', '訳', 3, ''),
		('A.esp', 'p derived', '', '', 0, 'derive')`)

	ctx := context.Background()
	narrations, err := s.NarrationsAfter(ctx, "A.esp", 0, 10, true)
	if err != nil || len(narrations) != 1 || narrations[0].Source != "n pending" {
		t.Fatalf("未訳叙述文ページ = rows:%+v err:%v", narrations, err)
	}
	lines, err := s.LinesAfter(ctx, "A.esp", 0, 10, true)
	if err != nil || len(lines) != 1 || lines[0].Source != "l pending" {
		t.Fatalf("未訳台詞ページ = rows:%+v err:%v", lines, err)
	}
	propers, err := s.ProperNounsAfter(ctx, "A.esp", 0, 10, true)
	if err != nil || len(propers) != 1 || propers[0].Source != "p pending" {
		t.Fatalf("未訳固有名ページ = rows:%+v err:%v", propers, err)
	}

	allNarrations, err := s.NarrationsAfter(ctx, "A.esp", 0, 10, false)
	if err != nil || len(allNarrations) != 2 {
		t.Fatalf("全叙述文ページ = rows:%+v err:%v", allNarrations, err)
	}
	allLines, err := s.LinesAfter(ctx, "A.esp", 0, 10, false)
	if err != nil || len(allLines) != 2 {
		t.Fatalf("全台詞ページ = rows:%+v err:%v", allLines, err)
	}
	allPropers, err := s.ProperNounsAfter(ctx, "A.esp", 0, 10, false)
	if err != nil || len(allPropers) != 2 {
		t.Fatalf("全固有名ページ = rows:%+v err:%v", allPropers, err)
	}
}
