package store

import (
	"context"
	"path/filepath"
	"testing"

	"github.com/jmoiron/sqlx"
	_ "modernc.org/sqlite"
)

// UpsertTargetPlugin が冪等であること。同じ plugin を 2 回登録しても 1 行で、source_path は更新され、
// created_at（初回登録時刻）は保たれる。開始時 upsert を再実行しても登録が二重化しないことを保証する。
func TestUpsertTargetPluginIdempotent(t *testing.T) {
	dbPath := filepath.Join(t.TempDir(), "central.sqlite3")
	s, err := Open(dbPath)
	if err != nil {
		t.Fatalf("Open: %v", err)
	}
	defer func() { _ = s.Close() }()
	ctx := context.Background()

	if err := s.UpsertTargetPlugin(ctx, "Dawnguard.esm", "/one/Data/Dawnguard.esm"); err != nil {
		t.Fatalf("1 回目の upsert: %v", err)
	}
	// created_at を固定値で観測するため、初回の値を読む。
	first := countRows(t, dbPath, `SELECT COUNT(*) FROM target_plugin`)
	if first != 1 {
		t.Fatalf("初回登録後の行数 = %d, want 1", first)
	}
	var createdBefore, sourceBefore string
	readTargetPlugin(t, dbPath, "Dawnguard.esm", &createdBefore, &sourceBefore)
	if sourceBefore != "/one/Data/Dawnguard.esm" {
		t.Errorf("source_path = %q, want /one/Data/Dawnguard.esm", sourceBefore)
	}

	// 別パスで再登録。行は増えず、source_path は更新、created_at は据え置き。
	if err := s.UpsertTargetPlugin(ctx, "Dawnguard.esm", "/two/Data/Dawnguard.esm"); err != nil {
		t.Fatalf("2 回目の upsert: %v", err)
	}
	if again := countRows(t, dbPath, `SELECT COUNT(*) FROM target_plugin`); again != 1 {
		t.Errorf("再登録後の行数 = %d, want 1（冪等）", again)
	}
	var createdAfter, sourceAfter string
	readTargetPlugin(t, dbPath, "Dawnguard.esm", &createdAfter, &sourceAfter)
	if sourceAfter != "/two/Data/Dawnguard.esm" {
		t.Errorf("再登録後の source_path = %q, want /two/Data/Dawnguard.esm", sourceAfter)
	}
	if createdAfter != createdBefore {
		t.Errorf("created_at が上書きされた: before=%q after=%q", createdBefore, createdAfter)
	}
}

// 準備完了は明示保存後だけ true になり、同じ plugin の再準備開始で false へ戻ること。
func TestSyncRetryReadinessLifecycle(t *testing.T) {
	dbPath := filepath.Join(t.TempDir(), "central.sqlite3")
	s, err := Open(dbPath)
	if err != nil {
		t.Fatalf("Open: %v", err)
	}
	defer func() { _ = s.Close() }()
	ctx := context.Background()

	ready, err := s.IsSyncRetryReady(ctx, "A.esp")
	if err != nil || ready {
		t.Fatalf("未登録 plugin の準備状態 = %v, err=%v, want false", ready, err)
	}
	if err := s.UpsertTargetPlugin(ctx, "A.esp", "/data/A.esp"); err != nil {
		t.Fatalf("登録: %v", err)
	}
	if err := s.MarkSyncRetryReady(ctx, "A.esp"); err != nil {
		t.Fatalf("準備完了の保存: %v", err)
	}
	if ready, err = s.IsSyncRetryReady(ctx, "A.esp"); err != nil || !ready {
		t.Fatalf("保存後の準備状態 = %v, err=%v, want true", ready, err)
	}
	if err := s.UpsertTargetPlugin(ctx, "A.esp", "/new/A.esp"); err != nil {
		t.Fatalf("再準備開始: %v", err)
	}
	if ready, err = s.IsSyncRetryReady(ctx, "A.esp"); err != nil || ready {
		t.Fatalf("再準備開始後の準備状態 = %v, err=%v, want false", ready, err)
	}
}

// DeleteTargetPlugin が対象 plugin の翻訳成果と連関だけを消し、別 plugin の成果と共有資産を残すこと。
// 束ねと削除の核。FK cascade を使わず明示 DELETE で対象スコープを閉じることを保証する。
func TestDeleteTargetPluginScope(t *testing.T) {
	dbPath := filepath.Join(t.TempDir(), "central.sqlite3")
	s, err := Open(dbPath)
	if err != nil {
		t.Fatalf("Open: %v", err)
	}
	defer func() { _ = s.Close() }()
	ctx := context.Background()

	if err := s.UpsertTargetPlugin(ctx, "A.esp", "/data/A.esp"); err != nil {
		t.Fatalf("A の登録: %v", err)
	}
	if err := s.UpsertTargetPlugin(ctx, "B.esp", "/data/B.esp"); err != nil {
		t.Fatalf("B の登録: %v", err)
	}
	seedScopeFixture(t, dbPath)

	if err := s.DeleteTargetPlugin(ctx, "A.esp"); err != nil {
		t.Fatalf("DeleteTargetPlugin: %v", err)
	}

	// A の対象スコープ本体・staging・連関は消える。
	gone := map[string]string{
		"target_plugin（A）":            `SELECT COUNT(*) FROM target_plugin WHERE plugin = 'A.esp'`,
		"narration（A）":                `SELECT COUNT(*) FROM narration WHERE plugin = 'A.esp'`,
		"line（A）":                     `SELECT COUNT(*) FROM line WHERE plugin = 'A.esp'`,
		"proper_noun（A）":              `SELECT COUNT(*) FROM proper_noun WHERE plugin = 'A.esp'`,
		"extracted_field（A）":          `SELECT COUNT(*) FROM extracted_field WHERE plugin = 'A.esp'`,
		"extracted_info_speaker（A）":   `SELECT COUNT(*) FROM extracted_info_speaker WHERE info_plugin = 'A.esp'`,
		"extracted_info_condition（A）": `SELECT COUNT(*) FROM extracted_info_condition WHERE info_plugin = 'A.esp'`,
		"narration_mention（A）":        `SELECT COUNT(*) FROM narration_mention WHERE narration_id = 10`,
		"narration_described（A）":      `SELECT COUNT(*) FROM narration_described WHERE narration_id = 10`,
		"line_mention（A）":             `SELECT COUNT(*) FROM line_mention WHERE line_id = 20`,
		"line_speaker（A）":             `SELECT COUNT(*) FROM line_speaker WHERE line_id = 20`,
		"line_condition（A）":           `SELECT COUNT(*) FROM line_condition WHERE line_id = 20`,
	}
	for name, query := range gone {
		if n := countRows(t, dbPath, query); n != 0 {
			t.Errorf("%s が残った（%d 件）。対象 plugin の削除で消えるべき", name, n)
		}
	}

	// B の対象スコープと連関、共有資産（speaker・master_term）は残る。
	kept := map[string]string{
		"target_plugin（B）":          `SELECT COUNT(*) FROM target_plugin WHERE plugin = 'B.esp'`,
		"narration（B）":              `SELECT COUNT(*) FROM narration WHERE plugin = 'B.esp'`,
		"line（B）":                   `SELECT COUNT(*) FROM line WHERE plugin = 'B.esp'`,
		"proper_noun（B）":            `SELECT COUNT(*) FROM proper_noun WHERE plugin = 'B.esp'`,
		"extracted_field（B）":        `SELECT COUNT(*) FROM extracted_field WHERE plugin = 'B.esp'`,
		"extracted_info_speaker（B）": `SELECT COUNT(*) FROM extracted_info_speaker WHERE info_plugin = 'B.esp'`,
		"narration_mention（B）":      `SELECT COUNT(*) FROM narration_mention WHERE narration_id = 11`,
		"共有 speaker":                `SELECT COUNT(*) FROM speaker WHERE id = 1`,
		"共有 master_term":            `SELECT COUNT(*) FROM master_term WHERE id = 1`,
	}
	for name, query := range kept {
		if n := countRows(t, dbPath, query); n != 1 {
			t.Errorf("%s が %d 件（want 1）。対象外は残すべき", name, n)
		}
	}
}

// plugin 識別子の導出が抽出器の TargetPlugin（ModKey.FileName、拡張子込みファイル名）と一致すること。
// api の upsert キーは filepath.Base(pluginPath)。ズレると narration.plugin 等と束ねられず削除が効かない不変ルール。
func TestPluginKeyDerivation(t *testing.T) {
	cases := map[string]string{
		"/Users/me/Skyrim/Data/Dawnguard.esm": "Dawnguard.esm",
		"/data/Inigo.esp":                     "Inigo.esp",
		"/x/y/3DNPC.esl":                      "3DNPC.esl",
		"Falskaar.esm":                        "Falskaar.esm",
	}
	for path, want := range cases {
		if got := filepath.Base(path); got != want {
			t.Errorf("filepath.Base(%q) = %q, want %q", path, got, want)
		}
	}
}

// readTargetPlugin は plugin の created_at と source_path を別接続で読む。
func readTargetPlugin(t *testing.T, dbPath, plugin string, createdAt, sourcePath *string) {
	t.Helper()
	db, err := sqlx.Connect("sqlite", dbPath)
	if err != nil {
		t.Fatalf("test connect: %v", err)
	}
	defer func() { _ = db.Close() }()
	row := db.QueryRowContext(context.Background(),
		`SELECT created_at, source_path FROM target_plugin WHERE plugin = ?`, plugin)
	if err := row.Scan(createdAt, sourcePath); err != nil {
		t.Fatalf("target_plugin の読み取り: %v", err)
	}
}

// seedScopeFixture は削除スコープ検証用に、A・B 両 plugin の対象スコープ行・staging・連関と共有資産を入れる。
func seedScopeFixture(t *testing.T, dbPath string) {
	t.Helper()
	execSQL(t, dbPath, `
-- 共有資産（削除で残す）。speaker は出自 mod でスコープ、master_term は横断。
INSERT INTO speaker (id, plugin, form_id) VALUES (1, 'Skyrim.esm', '0x1');
INSERT INTO master_term (id, source, dest) VALUES (1, 'Foo', 'フー');

-- 対象スコープ本体（A=10/20/30/40、B=11/21/31/41）。
INSERT INTO narration (id, source, plugin, form_id, edid, rec, field) VALUES
  (10, 's', 'A.esp', '0x10', 'eA', 'BOOK', 'DESC'),
  (11, 's', 'B.esp', '0x11', 'eB', 'BOOK', 'DESC');
INSERT INTO line (id, source, plugin, form_id, edid, rec, field) VALUES
  (20, 's', 'A.esp', '0x20', 'eA', 'INFO', 'NAM1'),
  (21, 's', 'B.esp', '0x21', 'eB', 'INFO', 'NAM1');
INSERT INTO proper_noun (id, plugin, source, category) VALUES
  (30, 'A.esp', 'Foo', 'BOOK'),
  (31, 'B.esp', 'Bar', 'BOOK');
INSERT INTO extracted_field (id, plugin, form_id, edid, rec, field, source) VALUES
  (40, 'A.esp', '0x10', 'eA', 'BOOK', 'DESC', 's'),
  (41, 'B.esp', '0x11', 'eB', 'BOOK', 'DESC', 's');

-- staging（info_plugin は対象 plugin 値）。
INSERT INTO extracted_info_speaker (info_plugin, info_form_id, speaker_id) VALUES
  ('A.esp', '0x20', 1), ('B.esp', '0x21', 1);
INSERT INTO extracted_info_condition (info_plugin, info_form_id, sex) VALUES
  ('A.esp', '0x20', 'Female'), ('B.esp', '0x21', 'Male');

-- 連関（A の narration=10 / line=20 に紐づく分と、B の narration=11）。
INSERT INTO narration_mention (narration_id, proper_noun_id) VALUES (10, 30), (11, 31);
INSERT INTO narration_described (narration_id, proper_noun_id) VALUES (10, 30);
INSERT INTO line_mention (line_id, proper_noun_id) VALUES (20, 30);
INSERT INTO line_speaker (line_id, speaker_id) VALUES (20, 1);
INSERT INTO line_condition (line_id, sex) VALUES (20, 'Female');
`)
}

// 仕様: 実行が完了し未訳が複数残ったとき／1 件だけ残ったとき／0 件のとき、CountUntranslated がその件数を返すこと。
// 数え方は一覧の進捗（total - translated）と一致し、機械派生した人名の部分形は数えない。
func TestCountUntranslated(t *testing.T) {
	dbPath := filepath.Join(t.TempDir(), "central.sqlite3")
	s, err := Open(dbPath)
	if err != nil {
		t.Fatalf("Open: %v", err)
	}
	defer func() { _ = s.Close() }()
	ctx := context.Background()

	if err = s.UpsertTargetPlugin(ctx, "A.esp", "/data/A.esp"); err != nil {
		t.Fatalf("A の登録: %v", err)
	}
	if err = s.UpsertTargetPlugin(ctx, "B.esp", "/data/B.esp"); err != nil {
		t.Fatalf("B の登録: %v", err)
	}
	// A.esp: 未訳 3 件（叙述文 1・台詞 1・固有名 1）。訳済み 2 件と、派生した部分形（origin != ''）1 件も置く。
	// 派生の部分形は未訳（status=0）でも翻訳対象でないため数えない。B.esp の未訳は A の集計へ入らない。
	execSQL(t, dbPath, `
INSERT INTO narration (id, source, plugin, form_id, edid, rec, field, status) VALUES
  (10, 's', 'A.esp', '0x10', 'eA', 'BOOK', 'DESC', 0),
  (11, 's', 'A.esp', '0x11', 'eA', 'BOOK', 'DESC', 3),
  (12, 's', 'B.esp', '0x12', 'eB', 'BOOK', 'DESC', 0);
INSERT INTO line (id, source, plugin, form_id, edid, rec, field, status) VALUES
  (20, 's', 'A.esp', '0x20', 'eA', 'INFO', 'NAM1', 0),
  (21, 's', 'A.esp', '0x21', 'eA', 'INFO', 'NAM1', 1);
INSERT INTO proper_noun (id, plugin, source, category, status, origin) VALUES
  (30, 'A.esp', 'Inigo', 'NPC_', 0, ''),
  (31, 'A.esp', 'Sorine', 'NPC_', 0, 'derived');
`)

	// 未訳が複数（叙述文 10・台詞 20・固有名 30 の 3 件）。
	got, err := s.CountUntranslated(ctx, "A.esp")
	if err != nil {
		t.Fatalf("CountUntranslated: %v", err)
	}
	if got != 3 {
		t.Errorf("未訳件数 = %d, want 3（派生した部分形と別 plugin は数えない）", got)
	}
	// 一覧の進捗（total - translated）と同じ値になること。
	rows, err := s.ListTargetPlugins(ctx)
	if err != nil {
		t.Fatalf("ListTargetPlugins: %v", err)
	}
	for _, r := range rows {
		if r.Plugin != "A.esp" {
			continue
		}
		if diff := r.Total - r.Translated; diff != got {
			t.Errorf("一覧の未訳 = %d, CountUntranslated = %d（同じ数え方であるべき）", diff, got)
		}
	}

	// 未訳が 1 件だけ残った状態（叙述文と台詞を訳済みにする）。
	execSQL(t, dbPath, `UPDATE narration SET status = 3 WHERE id = 10; UPDATE line SET status = 3 WHERE id = 20;`)
	if got, err = s.CountUntranslated(ctx, "A.esp"); err != nil || got != 1 {
		t.Errorf("未訳 1 件の集計 = %d（err=%v）, want 1", got, err)
	}

	// 未訳が 0 件（固有名も訳済み）。
	execSQL(t, dbPath, `UPDATE proper_noun SET status = 3 WHERE id = 30;`)
	if got, err = s.CountUntranslated(ctx, "A.esp"); err != nil || got != 0 {
		t.Errorf("未訳 0 件の集計 = %d（err=%v）, want 0", got, err)
	}
}
