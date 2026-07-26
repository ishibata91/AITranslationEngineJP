package store

import (
	"context"
	"path/filepath"
	"testing"

	"aitranslationenginejp/internal/model"

	"github.com/jmoiron/sqlx"
	_ "modernc.org/sqlite"
)

// テスト補助: 別接続で任意 SQL を実行する（抽出器や取込段が書く想定の行を直接入れる）。
func execSQL(t *testing.T, dbPath, query string) {
	t.Helper()
	db, err := sqlx.Connect("sqlite", dbPath)
	if err != nil {
		t.Fatalf("test connect: %v", err)
	}
	defer func() { _ = db.Close() }()
	if _, err := db.ExecContext(context.Background(), query); err != nil {
		t.Fatalf("test exec: %v", err)
	}
}

// テスト補助: 別接続で件数を数える。
func countRows(t *testing.T, dbPath, query string) int {
	t.Helper()
	db, err := sqlx.Connect("sqlite", dbPath)
	if err != nil {
		t.Fatalf("test connect: %v", err)
	}
	defer func() { _ = db.Close() }()
	var n int
	if err := db.Get(&n, query); err != nil {
		t.Fatalf("test count: %v", err)
	}
	return n
}

// 言及の投入が冪等であること。排他 2 列は 0 を NULL に写すため、SQLite の「UNIQUE 索引は NULL 同士を
// 別値と扱う」性質では重複を止められない。部分 UNIQUE 索引がその重複を止めることを証明する。
func TestInsertMentionsIdempotent(t *testing.T) {
	dbPath := filepath.Join(t.TempDir(), "central.sqlite3")
	s, err := Open(dbPath)
	if err != nil {
		t.Fatalf("Open: %v", err)
	}
	defer func() { _ = s.Close() }()

	narrRows := []model.NarrationMention{
		{NarrationID: 1, ProperNounID: 20}, // proper_noun 側の言及（master_term_id は NULL）
		{NarrationID: 1, MasterTermID: 10}, // master_term 側の言及（proper_noun_id は NULL）
	}
	lineRows := []model.LineMention{{LineID: 5, ProperNounID: 20}}

	n1, err := s.InsertNarrationMentions(context.Background(), narrRows)
	if err != nil {
		t.Fatalf("InsertNarrationMentions: %v", err)
	}
	l1, err := s.InsertLineMentions(context.Background(), lineRows)
	if err != nil {
		t.Fatalf("InsertLineMentions: %v", err)
	}
	if n1 != 2 || l1 != 1 {
		t.Fatalf("初回投入 = narration %d, line %d, want 2, 1", n1, l1)
	}

	// 同じ行の再投入は 1 件も増えない（再取込の冪等）。
	n2, err := s.InsertNarrationMentions(context.Background(), narrRows)
	if err != nil {
		t.Fatalf("InsertNarrationMentions 再投入: %v", err)
	}
	l2, err := s.InsertLineMentions(context.Background(), lineRows)
	if err != nil {
		t.Fatalf("InsertLineMentions 再投入: %v", err)
	}
	if n2 != 0 || l2 != 0 {
		t.Errorf("再投入 = narration %d, line %d, want 0, 0（冪等）", n2, l2)
	}
	if got := countRows(t, dbPath, `SELECT COUNT(*) FROM narration_mention`); got != 2 {
		t.Errorf("narration_mention 件数 = %d, want 2", got)
	}
	if got := countRows(t, dbPath, `SELECT COUNT(*) FROM line_mention`); got != 1 {
		t.Errorf("line_mention 件数 = %d, want 1", got)
	}
}

// 説明対象（e3）の解決が「box='叙述文' の行 × 同一レコードの FULL の固有名」だけを結ぶこと。
// 定型句の行と、FULL を持たないレコードの叙述文は結ばない。再実行は 1 件も増えない（冪等）。
func TestLinkNarrationDescribed(t *testing.T) {
	dbPath := filepath.Join(t.TempDir(), "central.sqlite3")
	s, err := Open(dbPath)
	if err != nil {
		t.Fatalf("Open: %v", err)
	}
	defer func() { _ = s.Close() }()

	// 同一レコード（P.esp/0x01 WEAP）の FULL（固有名）と DESC（叙述文）。この対だけが結ばれる。
	execSQL(t, dbPath, `INSERT INTO extracted_field (plugin, form_id, edid, rec, field, ordinal, source)
		VALUES ('P.esp', '0x01', 'Sword1', 'WEAP', 'FULL', 0, 'Dragonbane')`)
	// proper_noun.plugin は取込段が extracted_field.plugin を写す。位置解決は同一 plugin 内で結ぶ。
	execSQL(t, dbPath, `INSERT INTO proper_noun (plugin, source, category) VALUES ('P.esp', 'Dragonbane', 'WEAP')`)
	insertNarration(t, dbPath, model.Narration{
		Plugin: "P.esp", FormID: "0x01", EDID: "Sword1",
		Rec: "WEAP", Field: "DESC", Ordinal: 0, Source: "A blade of legend.", Status: 0,
	})
	// 定型句（ACTI:RNAM）は同一レコードに FULL（固有名）があっても説明対象を持たない。
	execSQL(t, dbPath, `INSERT INTO extracted_field (plugin, form_id, edid, rec, field, ordinal, source)
		VALUES ('P.esp', '0x02', 'Door1', 'ACTI', 'FULL', 0, 'Old Door')`)
	execSQL(t, dbPath, `INSERT INTO proper_noun (plugin, source, category) VALUES ('P.esp', 'Old Door', 'ACTI')`)
	insertNarration(t, dbPath, model.Narration{
		Plugin: "P.esp", FormID: "0x02", EDID: "Door1",
		Rec: "ACTI", Field: "RNAM", Ordinal: 0, Source: "Open", Status: 0,
	})
	// FULL を持たないレコード（LSCR）の叙述文は結ばれない。
	insertNarration(t, dbPath, model.Narration{
		Plugin: "P.esp", FormID: "0x03", EDID: "Load1",
		Rec: "LSCR", Field: "DESC", Ordinal: 0, Source: "The world is vast.", Status: 0,
	})

	linked, err := s.LinkNarrationDescribed(context.Background())
	if err != nil {
		t.Fatalf("LinkNarrationDescribed: %v", err)
	}
	if linked != 1 {
		t.Fatalf("解決件数 = %d, want 1（WEAP:DESC だけ）", linked)
	}
	// 結ばれた対が「WEAP:DESC の叙述文 → Dragonbane の固有名」であること。
	got := countRows(t, dbPath, `SELECT COUNT(*) FROM narration_described nd
		JOIN narration n ON n.id = nd.narration_id AND n.rec = 'WEAP' AND n.field = 'DESC'
		JOIN proper_noun pn ON pn.id = nd.proper_noun_id AND pn.source = 'Dragonbane'`)
	if got != 1 {
		t.Errorf("narration_described の対 = %d, want 1（WEAP:DESC → Dragonbane）", got)
	}

	again, err := s.LinkNarrationDescribed(context.Background())
	if err != nil {
		t.Fatalf("LinkNarrationDescribed 再実行: %v", err)
	}
	if again != 0 {
		t.Errorf("再実行の解決件数 = %d, want 0（冪等）", again)
	}
}
