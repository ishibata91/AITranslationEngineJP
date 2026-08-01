package harness

import (
	"context"
	"path/filepath"
	"strings"
	"testing"

	"aitranslationenginejp/internal/api"
	"aitranslationenginejp/internal/core/dictionary"
	"aitranslationenginejp/internal/core/rolespeech"
	"aitranslationenginejp/internal/engine"
	"aitranslationenginejp/internal/provider"
	"aitranslationenginejp/internal/store"

	"github.com/jmoiron/sqlx"
	_ "modernc.org/sqlite"
)

type recordingBatchProvider struct {
	requests []provider.BatchRequest
	fixed    map[string]string
}

func (p *recordingBatchProvider) SubmitBatch(_ context.Context, _ provider.Connection, _ string, requests []provider.BatchRequest) (string, error) {
	p.requests = append([]provider.BatchRequest(nil), requests...)
	return "batch-1", nil
}

func (p *recordingBatchProvider) PollBatch(_ context.Context, _ provider.Connection, _ string) (provider.BatchStatus, error) {
	return provider.BatchStatus{Total: len(p.requests), Succeeded: len(p.requests), Done: true}, nil
}

func (p *recordingBatchProvider) FetchResults(_ context.Context, _ provider.Connection, _ string) ([]provider.BatchResult, error) {
	results := make([]provider.BatchResult, len(p.requests))
	for i, request := range p.requests {
		results[i] = provider.BatchResult{CustomID: request.CustomID, Translation: p.fixed[request.Prompt.User]}
	}
	return results, nil
}

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

// R-2-1, R-2-2: 実 DB を使う batch 再送信も抽出系を繰り返さず、未訳 1 件だけを送って既訳を保つ。
func TestBatchRetriesOnlyUntranslatedRows(t *testing.T) {
	dir := t.TempDir()
	dbPath := filepath.Join(dir, "batch-retry.sqlite3")
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
	fixed := SyntheticAITranslations()
	translator := &RecordingProvider{Fixed: fixed}
	eng := engine.New(s, translator, fakeLexicon{}, roleSpeech, stoplist)
	if _, err := api.New(s, eng, nil, translator, extractor).RunExtractAndTranslate(api.RunRequest{PluginPath: pluginPath, Model: "fake-model"}); err != nil {
		t.Fatalf("初回実行: %v", err)
	}

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

	batchProvider := &recordingBatchProvider{fixed: fixed}
	runner := engine.NewBatchRunner(eng, map[string]provider.BatchTranslator{provider.BatchProviderXAI: batchProvider}, s)
	app := api.New(s, eng, runner, translator, extractor)
	outcome, err := app.SubmitBatchTranslation(api.RunRequest{
		PluginPath: pluginPath,
		Model:      "fake-model",
		Provider:   provider.BatchProviderXAI,
	})
	if err != nil {
		t.Fatalf("batch 再送信: %v", err)
	}
	if !outcome.ReusedPreparation || outcome.CompletedWithoutExternalBatch {
		t.Fatalf("batch 再送信結果 = %+v", outcome)
	}
	if len(extractor.Calls) != 2 {
		t.Fatalf("batch 再送信で抽出系を繰り返した: %v", extractor.Calls)
	}
	if len(batchProvider.requests) != 1 {
		t.Fatalf("batch 送信件数 = %d, want 1", len(batchProvider.requests))
	}
	if err := app.RefreshBatchTranslations(api.BatchPluginRequest{Plugin: fixture.PluginName, Provider: provider.BatchProviderXAI}); err != nil {
		t.Fatalf("batch 取り込み: %v", err)
	}
	var keptAfter string
	if err := db.Get(&keptAfter, `SELECT dest FROM narration WHERE id = ?`, keptID); err != nil {
		t.Fatalf("再送信後の既訳本文取得: %v", err)
	}
	if keptAfter != keptDest {
		t.Fatalf("既訳本文が変わった: before=%q after=%q", keptDest, keptAfter)
	}
	var retryStatus int
	if err := db.Get(&retryStatus, `SELECT status FROM narration WHERE id = ?`, retryID); err != nil {
		t.Fatalf("再送信行の状態取得: %v", err)
	}
	if retryStatus == 0 {
		t.Fatal("再送信した未訳へ訳が反映されていない")
	}
}
