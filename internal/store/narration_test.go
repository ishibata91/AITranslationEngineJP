package store

import (
	"context"
	"path/filepath"
	"testing"

	"aitranslationenginejp/internal/model"

	"github.com/jmoiron/sqlx"
	_ "modernc.org/sqlite"
)

// 抽出器（C#）が書く想定の narration 行を、別接続から直接入れるテスト補助。
func insertNarration(t *testing.T, dbPath string, n model.Narration) {
	t.Helper()
	db, err := sqlx.Connect("sqlite", dbPath)
	if err != nil {
		t.Fatalf("test connect: %v", err)
	}
	defer func() { _ = db.Close() }()
	_, err = db.Exec(
		`INSERT INTO narration (plugin, form_id, edid, rec, field, ordinal, source, status)
		 VALUES (?, ?, ?, ?, ?, ?, ?, ?)`,
		n.Plugin, n.FormID, n.EDID, n.Rec, n.Field, n.Ordinal, n.Source, n.Status,
	)
	if err != nil {
		t.Fatalf("test insert: %v", err)
	}
}

// Open で schema を適用し、未訳（status=0）の narration を読み、dest を書き戻すと未訳から外れること。
func TestListUntranslatedAndUpdateDest(t *testing.T) {
	dbPath := filepath.Join(t.TempDir(), "central.sqlite3")
	s, err := Open(dbPath)
	if err != nil {
		t.Fatalf("Open: %v", err)
	}
	defer func() { _ = s.Close() }()

	insertNarration(t, dbPath, model.Narration{
		Plugin: "Dawnguard.esm", FormID: "0x004a", EDID: "DLC1BookSerana",
		Rec: "BOOK", Field: "DESC", Ordinal: 0, Source: "I have walked these halls", Status: 0,
	})
	insertNarration(t, dbPath, model.Narration{
		Plugin: "Dawnguard.esm", FormID: "0x004b", EDID: "DLC1NoteArvak",
		Rec: "BOOK", Field: "DESC", Ordinal: 0, Source: "the Soul Cairn takes", Status: 0,
	})

	rows, err := s.ListUntranslatedNarrations(context.Background())
	if err != nil {
		t.Fatalf("ListUntranslatedNarrations: %v", err)
	}
	if len(rows) != 2 {
		t.Fatalf("untranslated = %d, want 2", len(rows))
	}
	if rows[0].EDID == "" || rows[0].Source == "" {
		t.Errorf("row not hydrated: %+v", rows[0])
	}

	if err = s.UpdateNarrationDest(context.Background(), rows[0].ID, "私は広間を歩いてきた", 3); err != nil {
		t.Fatalf("UpdateNarrationDest: %v", err)
	}

	after, err := s.ListUntranslatedNarrations(context.Background())
	if err != nil {
		t.Fatalf("ListUntranslatedNarrations after update: %v", err)
	}
	if len(after) != 1 {
		t.Fatalf("untranslated after update = %d, want 1", len(after))
	}
}
