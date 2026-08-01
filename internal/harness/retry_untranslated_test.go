package harness

import (
	"path/filepath"
	"strings"
	"testing"

	"aitranslationenginejp/internal/api"
	"aitranslationenginejp/internal/core/dictionary"
	"aitranslationenginejp/internal/core/rolespeech"
	"aitranslationenginejp/internal/engine"
	"aitranslationenginejp/internal/store"

	"github.com/jmoiron/sqlx"
	_ "modernc.org/sqlite"
)

// R-2-1 / R-2-3: 初回は全準備を行い、準備完了後の再実行は抽出系を繰り返さず未訳だけを書き戻すこと。
func TestRunExtractAndTranslateRetriesOnlyUntranslatedRows(t *testing.T) {
	dir := t.TempDir()
	dbPath := filepath.Join(dir, "retry.sqlite3")
	fixture := SyntheticFixture()
	pluginPath := filepath.Join("SyntheticData", fixture.PluginName)
	extractor := &SeedExtractor{DBPath: dbPath, Fixture: fixture}
	roleSpeech, err := rolespeech.ParseRoleSpeech(strings.NewReader(syntheticRoleSpeech))
	if err != nil {
		t.Fatalf("役割語: %v", err)
	}
	stoplist, err := dictionary.ParseStoplist(strings.NewReader(syntheticStopwords))
	if err != nil {
		t.Fatalf("stoplist: %v", err)
	}
	s, err := store.Open(dbPath)
	if err != nil {
		t.Fatalf("store.Open: %v", err)
	}
	defer func() { _ = s.Close() }()
	translator := &RecordingProvider{Fixed: SyntheticAITranslations()}
	eng := engine.New(s, translator, fakeLexicon{}, roleSpeech, stoplist)
	app := api.New(s, eng, nil, translator, extractor)

	if _, err := app.RunExtractAndTranslate(api.RunRequest{PluginPath: pluginPath, Model: "fake-model"}); err != nil {
		t.Fatalf("初回実行: %v", err)
	}
	if got := extractor.Calls; len(got) != 2 {
		t.Fatalf("初回の抽出系呼び出し = %v, want 2 件", got)
	}

	// 既訳 1 件を保持したまま別の 1 件だけを未訳へ戻し、再実行の対象を作る。
	db, err := sqlx.Connect("sqlite", dbPath)
	if err != nil {
		t.Fatalf("観測用 DB: %v", err)
	}
	defer func() { _ = db.Close() }()
	var keptID, retryID int64
	if err := db.Get(&keptID, `SELECT id FROM narration WHERE plugin = ? AND status != 0 ORDER BY id LIMIT 1`, fixture.PluginName); err != nil {
		t.Fatalf("保持する既訳の取得: %v", err)
	}
	if err := db.Get(&retryID, `SELECT id FROM narration WHERE plugin = ? AND status != 0 AND id != ? ORDER BY id LIMIT 1`, fixture.PluginName, keptID); err != nil {
		t.Fatalf("未訳へ戻す行の取得: %v", err)
	}
	var keptDest string
	if err := db.Get(&keptDest, `SELECT dest FROM narration WHERE id = ?`, keptID); err != nil {
		t.Fatalf("既訳本文の取得: %v", err)
	}
	if _, err := db.Exec(`UPDATE narration SET dest = '', status = 0 WHERE id = ?`, retryID); err != nil {
		t.Fatalf("未訳状態の作成: %v", err)
	}

	result, err := app.RunExtractAndTranslate(api.RunRequest{PluginPath: pluginPath, Model: "fake-model"})
	if err != nil {
		t.Fatalf("再実行: %v", err)
	}
	if result.TranslatedCount != 1 {
		t.Errorf("再実行の翻訳件数 = %d, want 1", result.TranslatedCount)
	}
	if got := extractor.Calls; len(got) != 2 {
		t.Errorf("再実行で抽出系を繰り返した: %v", got)
	}
	var keptAfter string
	if err := db.Get(&keptAfter, `SELECT dest FROM narration WHERE id = ?`, keptID); err != nil {
		t.Fatalf("再実行後の既訳本文取得: %v", err)
	}
	if keptAfter != keptDest {
		t.Errorf("既訳本文が変わった: before=%q after=%q", keptDest, keptAfter)
	}
}
