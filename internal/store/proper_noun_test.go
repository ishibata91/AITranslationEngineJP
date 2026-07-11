package store

import (
	"context"
	"path/filepath"
	"testing"

	"aitranslationenginejp/internal/model"
)

// IngestProperNouns の重複排除が plugin スコープになること（UNIQUE(plugin, category, source) の効き）。
// 別 plugin の同綴り固有名は別行に残し（横断の非共有）、同一 plugin 内の同綴りは 1 つにまとめる。
// 再投入は 1 件も増えない（plugin 内の冪等）。
func TestIngestProperNounsPluginScoped(t *testing.T) {
	dbPath := filepath.Join(t.TempDir(), "central.sqlite3")
	s, err := Open(dbPath)
	if err != nil {
		t.Fatalf("Open: %v", err)
	}
	defer func() { _ = s.Close() }()

	rows := []model.ProperNoun{
		{Plugin: "A.esp", Source: "Frost", Category: "WEAP"}, // A の Frost
		{Plugin: "B.esp", Source: "Frost", Category: "WEAP"}, // B の Frost（同綴りだが別 plugin → 別行）
		{Plugin: "A.esp", Source: "Frost", Category: "WEAP"}, // A の Frost 再掲（plugin 内で重複 → まとまる）
	}

	written, err := s.IngestProperNouns(context.Background(), rows)
	if err != nil {
		t.Fatalf("IngestProperNouns: %v", err)
	}
	// 追加できたのは A/Frost と B/Frost の 2 件（3 件目は A/Frost の重複）。
	if written != 2 {
		t.Fatalf("初回投入の追加件数 = %d, want 2（A/Frost・B/Frost、A/Frost 重複は除く）", written)
	}
	if got := countRows(t, dbPath, `SELECT COUNT(*) FROM proper_noun`); got != 2 {
		t.Errorf("proper_noun 総件数 = %d, want 2", got)
	}
	// 別 plugin の同綴りが別行で残ること（横断の非共有）。
	if got := countRows(t, dbPath, `SELECT COUNT(*) FROM proper_noun WHERE plugin = 'A.esp'`); got != 1 {
		t.Errorf("A.esp の件数 = %d, want 1", got)
	}
	if got := countRows(t, dbPath, `SELECT COUNT(*) FROM proper_noun WHERE plugin = 'B.esp'`); got != 1 {
		t.Errorf("B.esp の件数 = %d, want 1", got)
	}

	// 同じ行の再投入は 1 件も増えない（plugin 内の冪等）。
	again, err := s.IngestProperNouns(context.Background(), rows)
	if err != nil {
		t.Fatalf("IngestProperNouns 再投入: %v", err)
	}
	if again != 0 {
		t.Errorf("再投入の追加件数 = %d, want 0（冪等）", again)
	}
	if got := countRows(t, dbPath, `SELECT COUNT(*) FROM proper_noun`); got != 2 {
		t.Errorf("再投入後の proper_noun 総件数 = %d, want 2", got)
	}
}

// ProperNounPlacementsForExport が pn.plugin = ef.plugin で plugin 別に位置解決すること。
// 2 つの plugin に同綴りの固有名（別訳）があるとき、各 plugin の書き出しは自分の訳だけを結ぶ。
// plugin 結合が無いと A の Frost に B の proper_noun も当たり、訳の取り違え・重複が起きる。
func TestProperNounPlacementsForExportPluginScoped(t *testing.T) {
	dbPath := filepath.Join(t.TempDir(), "central.sqlite3")
	s, err := Open(dbPath)
	if err != nil {
		t.Fatalf("Open: %v", err)
	}
	defer func() { _ = s.Close() }()

	// 2 plugin が同綴り "Frost"（WEAP:FULL＝box '固有名'）を持ち、別々の確定訳を付けた状態を作る。
	execSQL(t, dbPath, `INSERT INTO extracted_field (plugin, form_id, edid, rec, field, ordinal, source)
		VALUES ('A.esp', '0x01', 'FrostA', 'WEAP', 'FULL', 0, 'Frost')`)
	execSQL(t, dbPath, `INSERT INTO extracted_field (plugin, form_id, edid, rec, field, ordinal, source)
		VALUES ('B.esp', '0x02', 'FrostB', 'WEAP', 'FULL', 0, 'Frost')`)
	execSQL(t, dbPath, `INSERT INTO proper_noun (plugin, source, category, dest, status)
		VALUES ('A.esp', 'Frost', 'WEAP', '霜A', 3)`)
	execSQL(t, dbPath, `INSERT INTO proper_noun (plugin, source, category, dest, status)
		VALUES ('B.esp', 'Frost', 'WEAP', '霜B', 3)`)

	// A の書き出しは A の訳（霜A）だけを 1 件結ぶ（B の proper_noun を混ぜない）。
	a, err := s.ProperNounPlacementsForExport(context.Background(), "A.esp")
	if err != nil {
		t.Fatalf("ProperNounPlacementsForExport(A.esp): %v", err)
	}
	if len(a) != 1 {
		t.Fatalf("A.esp の位置解決件数 = %d, want 1（自分の訳だけ）", len(a))
	}
	if a[0].Plugin != "A.esp" || a[0].Source != "Frost" || a[0].Dest != "霜A" {
		t.Errorf("A.esp の位置解決 = %+v, want plugin=A.esp source=Frost dest=霜A", a[0])
	}

	// B の書き出しは B の訳（霜B）だけを結ぶ。
	b, err := s.ProperNounPlacementsForExport(context.Background(), "B.esp")
	if err != nil {
		t.Fatalf("ProperNounPlacementsForExport(B.esp): %v", err)
	}
	if len(b) != 1 {
		t.Fatalf("B.esp の位置解決件数 = %d, want 1", len(b))
	}
	if b[0].Dest != "霜B" {
		t.Errorf("B.esp の位置解決の訳 = %q, want 霜B", b[0].Dest)
	}
}
