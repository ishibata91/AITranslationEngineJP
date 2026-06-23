package main

import (
	"os"
	"path/filepath"
	"testing"

	"aitranslationenginejp/db"

	"github.com/jmoiron/sqlx"
	_ "modernc.org/sqlite"
)

// run が XML から人名の部分形を派生し、base 衝突を除いて master_term へ追記すること（結合テスト）。
// XML 解析・派生の本体は engine.DeriveTermsFromXMLDir（engine の単体テストで担保）で、本テストは
// run の I/O 配線（base 読み出し・派生追記・base 衝突 skip）を temp DB で確かめる。
// base に既出の原語（Maro）は派生せず、byname（Grelod）と two（Mercer / Frey）を category 印付きで足す。
func TestRunDerivesAndWritesMasterTerm(t *testing.T) {
	dir := t.TempDir()
	xmlDir := filepath.Join(dir, "xml")
	if err := os.MkdirAll(xmlDir, 0o755); err != nil {
		t.Fatal(err)
	}
	// ファイル名を base ゲーム接頭（Skyrim）にして two を有効にする。
	const fixture = `<?xml version="1.0" encoding="UTF-8"?>
<SSTXMLRessources>
  <Content>
    <String><REC>NPC_:FULL</REC><Source>Grelod the Kind</Source><Dest>親切者のグレロッド</Dest></String>
    <String><REC>NPC_:SHRT</REC><Source>Maro</Source><Dest>マロ</Dest></String>
    <String><REC>NPC_:FULL</REC><Source>Mercer Frey</Source><Dest>メルセル・フレイ</Dest></String>
    <String><REC>INFO:NAM1</REC><Source>I met Mercer once.</Source><Dest>x</Dest></String>
  </Content>
</SSTXMLRessources>`
	if err := os.WriteFile(filepath.Join(xmlDir, "Skyrim_english_japanese.xml"), []byte(fixture), 0o644); err != nil {
		t.Fatal(err)
	}

	dbPath := filepath.Join(dir, "central.sqlite3")
	// base 辞書を先に用意する: Maro を base として入れ、派生側で衝突 skip されることを確かめる。
	seedBaseTerm(t, dbPath, "Maro", "マロ", "NPC_")

	if err := run(dbPath, xmlDir); err != nil {
		t.Fatalf("run error: %v", err)
	}

	got := readMasterTerms(t, dbPath)
	// byname と two の派生が category 印付きで入る。
	assertHasTerm(t, got, "Grelod", "グレロッド", "derive:byname")
	assertHasTerm(t, got, "Mercer", "メルセル", "derive:two")
	assertHasTerm(t, got, "Frey", "フレイ", "derive:two")
	// base 既出の Maro は派生（derive:shrt）として追記されない。
	for _, r := range got {
		if r.Source == "Maro" && r.Category == "derive:shrt" {
			t.Errorf("Maro が派生として追記された（base 衝突を skip していない）: %+v", r)
		}
	}
}

type termRow struct {
	Source   string `db:"source"`
	Dest     string `db:"dest"`
	Category string `db:"category"`
}

func seedBaseTerm(t *testing.T, dbPath, source, dest, category string) {
	t.Helper()
	conn, err := sqlx.Connect("sqlite", dbPath)
	if err != nil {
		t.Fatal(err)
	}
	defer conn.Close()
	if err := db.Apply(conn); err != nil {
		t.Fatal(err)
	}
	if _, err := conn.Exec(`INSERT INTO master_term (source, dest, category) VALUES (?, ?, ?)`, source, dest, category); err != nil {
		t.Fatal(err)
	}
}

func readMasterTerms(t *testing.T, dbPath string) []termRow {
	t.Helper()
	conn, err := sqlx.Connect("sqlite", dbPath)
	if err != nil {
		t.Fatal(err)
	}
	defer conn.Close()
	var rows []termRow
	if err := conn.Select(&rows, `SELECT source, dest, category FROM master_term`); err != nil {
		t.Fatal(err)
	}
	return rows
}

func assertHasTerm(t *testing.T, rows []termRow, source, dest, category string) {
	t.Helper()
	for _, r := range rows {
		if r.Source == source && r.Dest == dest && r.Category == category {
			return
		}
	}
	t.Errorf("master_term に %q=>%q (%s) が無い。rows=%+v", source, dest, category, rows)
}
