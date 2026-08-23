package engine

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"testing"

	"aitranslationenginejp/internal/model"
	"aitranslationenginejp/internal/provider"
)

// --- fakeStore の BatchStore 実装（batch 結合テスト専用。1 plugin 1 進行を想定） ---

func (f *fakeStore) StartBatchProgression(_ context.Context, plugin, providerName, modelName string) (int64, error) {
	if f.batchProg == nil {
		f.nextBatchID++
		f.batchProg = &model.BatchTranslation{ID: f.nextBatchID, Plugin: plugin, Provider: providerName, Model: modelName, Stage: model.BatchStageProperNoun}
	} else {
		// reset: 同じ行を使い回し、進行を固有名段へ戻し、外部 ID と送信行対応を空へ戻す（実 store と同じ）。
		f.batchProg.Provider = providerName
		f.batchProg.Model = modelName
		f.batchProg.Stage = model.BatchStageProperNoun
		f.batchProg.ProperBatchID = ""
		f.batchProg.BodyBatchID = ""
	}
	f.batchReqs = nil
	return f.batchProg.ID, nil
}

func (f *fakeStore) RecordBatchExternalID(_ context.Context, _ int64, stage, externalID string) error {
	switch stage {
	case model.BatchStageProperNoun:
		f.batchProg.ProperBatchID = externalID
	case model.BatchStageBody:
		f.batchProg.BodyBatchID = externalID
	default:
		return fmt.Errorf("外部 ID を持たない進行段: %q", stage)
	}
	return nil
}

func (f *fakeStore) AdvanceBatchStage(_ context.Context, _ int64, stage string) error {
	f.batchProg.Stage = stage
	return nil
}

func (f *fakeStore) GetBatchProgression(_ context.Context, plugin string) (model.BatchTranslation, bool, error) {
	if f.batchProg == nil || f.batchProg.Plugin != plugin {
		return model.BatchTranslation{}, false, nil
	}
	return *f.batchProg, true, nil
}

func (f *fakeStore) InsertBatchRequests(_ context.Context, rows []model.BatchRequest) (int, error) {
	f.batchReqs = append(f.batchReqs, rows...)
	return len(rows), nil
}

func (f *fakeStore) ListBatchRequests(_ context.Context, externalBatchID string) ([]model.BatchRequest, error) {
	var out []model.BatchRequest
	for _, r := range f.batchReqs {
		if r.ExternalBatchID == externalBatchID {
			out = append(out, r)
		}
	}
	return out, nil
}

func (f *fakeStore) ListBatchRequestsByStage(_ context.Context, batchID int64, stage string) ([]model.BatchRequest, error) {
	var out []model.BatchRequest
	for _, row := range f.batchReqs {
		if row.BatchID != batchID {
			continue
		}
		if stage == model.BatchStageProperNoun && row.Kind == model.BatchKindProper {
			out = append(out, row)
		}
		if stage == model.BatchStageBody && (row.Kind == model.BatchKindNarration || row.Kind == model.BatchKindLine) {
			out = append(out, row)
		}
	}
	return out, nil
}

func (f *fakeStore) UpdateBatchRequestSendState(_ context.Context, batchID int64, customIDs []string, state, externalID string) error {
	for i := range f.batchReqs {
		if f.batchReqs[i].BatchID != batchID {
			continue
		}
		for _, customID := range customIDs {
			if f.batchReqs[i].CustomID == customID {
				f.batchReqs[i].SendState = state
				f.batchReqs[i].ExternalBatchID = externalID
			}
		}
	}
	return nil
}

func (f *fakeStore) MarkSyncRetryReady(_ context.Context, _ string) error {
	f.syncRetryReady = true
	return nil
}

func (f *fakeStore) IsSyncRetryReady(_ context.Context, _ string) (bool, error) {
	return f.syncRetryReady, nil
}

// fakeBatchProvider は provider.BatchTranslator の in-memory 実装（実 xAI API に触れない・課金しない）。
// 送信した要求を外部 ID ごとに保持し、状態確認は常に終端（即完了）を返す。結果は out（user メッセージ→訳文）から引く。
// errByUser に載る user は失敗種別を返す（据え置きの再現）。fakeTranslator と同じ out を共有すれば、同期と同じ訳文を返す。
type fakeBatchProvider struct {
	out          map[string]string
	errByUser    map[string]error
	batches      map[string][]provider.BatchRequest
	nextID       int
	submits      int // SubmitBatch の呼び出し回数（固有名段・本文段の 2 回になることの観測）
	polls        int
	fetches      int
	failSubmitsN int // 先頭から何回の SubmitBatch を失敗させるか（送信 HTTP 失敗の再現。0 なら常に成功）
	failAtSubmit int // 指定回のSubmitBatchだけを失敗させる。本文段だけの失敗再開を再現する。
	attempts     int
	pollErr      error
}

func (f *fakeBatchProvider) SubmitBatch(_ context.Context, _ provider.Connection, _ string, requests []provider.BatchRequest) (string, error) {
	if f.batches == nil {
		f.batches = map[string][]provider.BatchRequest{}
	}
	f.attempts++
	if f.failAtSubmit == f.attempts {
		return "", fmt.Errorf("fake 送信失敗（HTTP 失敗の再現）")
	}
	if f.failSubmitsN > 0 {
		f.failSubmitsN--
		return "", fmt.Errorf("fake 送信失敗（HTTP 失敗の再現）")
	}
	f.nextID++
	f.submits++
	id := fmt.Sprintf("ext-batch-%d", f.nextID)
	f.batches[id] = requests
	return id, nil
}

func (f *fakeBatchProvider) PollBatch(_ context.Context, _ provider.Connection, externalBatchID string) (provider.BatchStatus, error) {
	f.polls++
	if f.pollErr != nil {
		return provider.BatchStatus{}, f.pollErr
	}
	reqs := f.batches[externalBatchID]
	return provider.BatchStatus{Total: len(reqs), Succeeded: len(reqs), Done: true}, nil
}

// R-1-1, R-1-4: 1001件は1000件と1件へ分かれ、現在の外部 batch が完了するまで次を送らない。
func TestBatchSkipsJapaneseRows(t *testing.T) {
	ctx := context.Background()
	const plugin = "japanese.esp"
	store := &fakeStore{
		proper:       []model.ProperNoun{{ID: 1, Plugin: plugin, Source: "剣"}},
		untranslated: []model.Narration{{ID: 2, Plugin: plugin, Source: "本"}},
		lines:        []model.Line{{ID: 3, Plugin: plugin, Source: "話す"}},
	}
	batch := &fakeBatchProvider{}
	runner := NewBatchRunner(New(store, &fakeTranslator{}, fakeLexicon{}, nil, nil), map[string]provider.BatchTranslator{provider.BatchProviderXAI: batch}, store)

	outcome, err := runner.SubmitBatch(ctx, provider.BatchProviderXAI, provider.Connection{}, "grok", plugin)
	if err != nil {
		t.Fatalf("SubmitBatch: %v", err)
	}
	if !outcome.CompletedWithoutExternalBatch || batch.submits != 0 {
		t.Errorf("batch結果 = %+v, submits = %d", outcome, batch.submits)
	}
	if len(store.properUpdates) != 1 || store.properUpdates[0].dest != "剣" || store.properUpdates[0].status != statusTranslated {
		t.Errorf("固有名の更新 = %+v", store.properUpdates)
	}
	if len(store.updates) != 1 || store.updates[0].dest != "本" || store.updates[0].status != statusTranslated {
		t.Errorf("叙述文の更新 = %+v", store.updates)
	}
	if len(store.lineUpdates) != 1 || store.lineUpdates[0].dest != "話す" || store.lineUpdates[0].status != statusTranslated {
		t.Errorf("台詞の更新 = %+v", store.lineUpdates)
	}
}

func TestBatchは1001件を最大1000件ずつ順番に送る(t *testing.T) {
	ctx := context.Background()
	const plugin = "chunked.esp"
	store := &fakeStore{}
	out := make(map[string]string, 1001)
	for i := 1; i <= 1001; i++ {
		source := fmt.Sprintf("line-%04d", i)
		store.lines = append(store.lines, model.Line{ID: int64(i), Plugin: plugin, Rec: "INFO", Field: "NAM1", Source: source})
		out[source] = fmt.Sprintf("訳-%04d", i)
	}
	fb := &fakeBatchProvider{out: out}
	runner := NewBatchRunner(New(store, &fakeTranslator{}, fakeLexicon{}, nil, nil), map[string]provider.BatchTranslator{provider.BatchProviderXAI: fb}, store)

	if _, err := runner.SubmitBatch(ctx, provider.BatchProviderXAI, provider.Connection{}, "grok", plugin); err != nil {
		t.Fatalf("初回送信: %v", err)
	}
	if fb.submits != 1 || len(fb.batches["ext-batch-1"]) != 1000 {
		t.Fatalf("初回外部 batch = submits:%d requests:%d", fb.submits, len(fb.batches["ext-batch-1"]))
	}
	progress, ok, err := runner.ProgressStatus(ctx, provider.BatchProviderXAI, provider.Connection{}, plugin)
	if err != nil || !ok || progress.Total != 1000 {
		t.Fatalf("初回進行 = %+v ok:%t err:%v", progress, ok, err)
	}
	if err := runner.RefreshPlugin(ctx, provider.BatchProviderXAI, provider.Connection{}, plugin); err != nil {
		t.Fatalf("1回目の反映: %v", err)
	}
	if fb.submits != 2 || len(fb.batches["ext-batch-2"]) != 1 || store.batchProg.Stage != model.BatchStageBody {
		t.Fatalf("次の外部 batch = submits:%d requests:%d progression:%+v", fb.submits, len(fb.batches["ext-batch-2"]), store.batchProg)
	}
	progress, ok, err = runner.ProgressStatus(ctx, provider.BatchProviderXAI, provider.Connection{}, plugin)
	if err != nil || !ok || progress.Total != 1 {
		t.Fatalf("2回目進行 = %+v ok:%t err:%v", progress, ok, err)
	}
	if err := runner.RefreshPlugin(ctx, provider.BatchProviderXAI, provider.Connection{}, plugin); err != nil {
		t.Fatalf("2回目の反映: %v", err)
	}
	if store.batchProg.Stage != model.BatchStageDone || len(store.lineUpdates) != 1001 {
		t.Fatalf("完了状態 = progression:%+v updates:%d", store.batchProg, len(store.lineUpdates))
	}
}

// R-1-2: ちょうど1000件は一つの外部 batch として送る。
func TestBatchは1000件を一つの外部Batchとして送る(t *testing.T) {
	ctx := context.Background()
	const plugin = "boundary-1000.esp"
	store := &fakeStore{}
	for i := 1; i <= 1000; i++ {
		store.lines = append(store.lines, model.Line{ID: int64(i), Plugin: plugin, Rec: "INFO", Field: "NAM1", Source: fmt.Sprintf("line-%04d", i)})
	}
	fb := &fakeBatchProvider{}
	runner := NewBatchRunner(New(store, &fakeTranslator{}, fakeLexicon{}, nil, nil), map[string]provider.BatchTranslator{provider.BatchProviderXAI: fb}, store)
	if _, err := runner.SubmitBatch(ctx, provider.BatchProviderXAI, provider.Connection{}, "grok", plugin); err != nil {
		t.Fatalf("送信: %v", err)
	}
	if fb.submits != 1 || len(fb.batches["ext-batch-1"]) != 1000 {
		t.Fatalf("外部 batch = submits:%d requests:%d", fb.submits, len(fb.batches["ext-batch-1"]))
	}
}

// R-1-3: 対象が0件なら外部 batch を作らず完了する。
func TestBatchは対象0件なら外部Batchを作らない(t *testing.T) {
	ctx := context.Background()
	const plugin = "empty.esp"
	store := &fakeStore{}
	fb := &fakeBatchProvider{}
	runner := NewBatchRunner(New(store, &fakeTranslator{}, fakeLexicon{}, nil, nil), map[string]provider.BatchTranslator{provider.BatchProviderXAI: fb}, store)
	outcome, err := runner.SubmitBatch(ctx, provider.BatchProviderXAI, provider.Connection{}, "grok", plugin)
	if err != nil {
		t.Fatalf("送信: %v", err)
	}
	if !outcome.CompletedWithoutExternalBatch || fb.submits != 0 || store.batchProg.Stage != model.BatchStageDone {
		t.Fatalf("完了状態 = outcome:%+v submits:%d progression:%+v", outcome, fb.submits, store.batchProg)
	}
}

// R-2-1, R-2-2: 状態確認が failed を返した場合は結果取得、次の送信、進行段更新を行わない。
func TestBatch状態確認失敗は取り込みと進行を止める(t *testing.T) {
	ctx := context.Background()
	const plugin = "failed-batch.esp"
	store := &fakeStore{lines: []model.Line{{ID: 1, Plugin: plugin, Rec: "INFO", Field: "NAM1", Source: "line"}}}
	fb := &fakeBatchProvider{pollErr: errors.New("OpenAI batch batch-failed failed: token_limit_exceeded: queued token limit")}
	runner := NewBatchRunner(New(store, &fakeTranslator{}, fakeLexicon{}, nil, nil), map[string]provider.BatchTranslator{provider.BatchProviderOpenAI: fb}, store)
	conn := provider.Connection{APIKey: "key"}
	if _, err := runner.SubmitBatch(ctx, provider.BatchProviderOpenAI, conn, "gpt", plugin); err != nil {
		t.Fatalf("送信: %v", err)
	}
	before := *store.batchProg
	if err := runner.RefreshPlugin(ctx, provider.BatchProviderOpenAI, conn, plugin); err == nil || !strings.Contains(err.Error(), "batch-failed") {
		t.Fatalf("取り込み error = %v", err)
	}
	if fb.fetches != 0 || fb.submits != 1 || *store.batchProg != before || len(store.lineUpdates) != 0 {
		t.Fatalf("失敗後に処理を進めた: fetches:%d submits:%d progression:%+v updates:%v", fb.fetches, fb.submits, store.batchProg, store.lineUpdates)
	}
}

func (f *fakeBatchProvider) FetchResults(_ context.Context, _ provider.Connection, externalBatchID string) ([]provider.BatchResult, error) {
	f.fetches++
	var out []provider.BatchResult
	for _, r := range f.batches[externalBatchID] {
		res := provider.BatchResult{CustomID: r.CustomID}
		if err, ok := f.errByUser[r.Prompt.User]; ok {
			res.Err = err
		} else {
			res.Translation = f.out[r.Prompt.User]
		}
		out = append(out, res)
	}
	return out, nil
}

// R-1-5: 保存済み OpenAI 進行を xAI として状態確認または取り込みできず、外部 API と進行を変更しないこと。
func Test進行中のOpenAIBatchをXAIとして状態確認または取り込みできない(t *testing.T) {
	ctx := context.Background()
	conn := provider.Connection{Endpoint: "http://x", APIKey: "key"}
	const plugin = "provider-boundary.esp"
	store := seedFor(plugin)
	fb := &fakeBatchProvider{out: translationOut()}
	runner := NewBatchRunner(New(store, &fakeTranslator{out: translationOut()}, fakeLexicon{}, nil, nil), map[string]provider.BatchTranslator{
		provider.BatchProviderOpenAI: fb,
		provider.BatchProviderXAI:    fb,
	}, store)
	if _, err := runner.SubmitBatch(ctx, provider.BatchProviderOpenAI, conn, "gpt-5.6-luna", plugin); err != nil {
		t.Fatalf("OpenAI batch 送信: %v", err)
	}
	before := *store.batchProg
	fb.polls, fb.fetches = 0, 0
	if _, _, err := runner.ProgressStatus(ctx, provider.BatchProviderXAI, conn, plugin); !errors.Is(err, errBatchProviderMismatch) {
		t.Fatalf("状態確認 error = %v", err)
	}
	if err := runner.RefreshPlugin(ctx, provider.BatchProviderXAI, conn, plugin); !errors.Is(err, errBatchProviderMismatch) {
		t.Fatalf("取り込み error = %v", err)
	}
	if fb.polls != 0 || fb.fetches != 0 {
		t.Errorf("不一致で外部 API を呼んだ: polls=%d fetches=%d", fb.polls, fb.fetches)
	}
	if *store.batchProg != before {
		t.Errorf("不一致で進行を変更した: before=%+v after=%+v", before, *store.batchProg)
	}
}

// R-1-6: OpenAI API キーが空なら送信、状態確認、取り込みのいずれも開始しないこと。
func TestOpenAIAPIキーが空ならBatch操作を開始しない(t *testing.T) {
	ctx := context.Background()
	const plugin = "no-key.esp"
	store := seedFor(plugin)
	fb := &fakeBatchProvider{out: translationOut()}
	runner := NewBatchRunner(New(store, &fakeTranslator{out: translationOut()}, fakeLexicon{}, nil, nil), map[string]provider.BatchTranslator{
		provider.BatchProviderOpenAI: fb,
	}, store)
	if _, err := runner.SubmitBatch(ctx, provider.BatchProviderOpenAI, provider.Connection{}, "gpt-5.6-luna", plugin); err == nil {
		t.Fatalf("API キーなしの送信が成功した")
	}
	if store.batchProg != nil || fb.submits != 0 {
		t.Errorf("API キーなしで送信状態を変更した: progression=%+v submits=%d", store.batchProg, fb.submits)
	}
	store.batchProg = &model.BatchTranslation{ID: 1, Plugin: plugin, Provider: provider.BatchProviderOpenAI, Model: "gpt-5.6-luna", Stage: model.BatchStageBody, BodyBatchID: "ext"}
	if _, _, err := runner.ProgressStatus(ctx, provider.BatchProviderOpenAI, provider.Connection{}, plugin); err == nil {
		t.Fatalf("API キーなしの状態確認が成功した")
	}
	if err := runner.RefreshPlugin(ctx, provider.BatchProviderOpenAI, provider.Connection{}, plugin); err == nil {
		t.Fatalf("API キーなしの取り込みが成功した")
	}
	if fb.polls != 0 || fb.fetches != 0 {
		t.Errorf("API キーなしで外部 API を呼んだ: polls=%d fetches=%d", fb.polls, fb.fetches)
	}
}

// R-1-7: 全行が失敗して成功訳が無くても進行を完了し、未訳だけを再送信できること。
func TestOpenAIBatchが全件失敗しても未訳を再送信できる(t *testing.T) {
	ctx := context.Background()
	conn := provider.Connection{APIKey: "key"}
	const plugin = "all-failed.esp"
	store := &fakeStore{lines: []model.Line{{ID: 1, Plugin: plugin, Rec: "INFO", Field: "NAM1", Source: "fail."}}}
	fb := &fakeBatchProvider{errByUser: map[string]error{"fail.": fmt.Errorf("%w: failed", provider.ErrServerTransient)}}
	runner := NewBatchRunner(New(store, &fakeTranslator{}, fakeLexicon{}, nil, nil), map[string]provider.BatchTranslator{provider.BatchProviderOpenAI: fb}, store)
	if _, err := runner.SubmitBatch(ctx, provider.BatchProviderOpenAI, conn, "gpt-5.6-luna", plugin); err != nil {
		t.Fatalf("初回送信: %v", err)
	}
	if err := runner.RefreshPlugin(ctx, provider.BatchProviderOpenAI, conn, plugin); err != nil {
		t.Fatalf("全件失敗の取り込み: %v", err)
	}
	if store.batchProg.Stage != model.BatchStageDone || len(store.lineUpdates) != 0 {
		t.Fatalf("全件失敗後 = progression=%+v updates=%v", store.batchProg, store.lineUpdates)
	}
	if _, err := runner.SubmitBatch(ctx, provider.BatchProviderOpenAI, conn, "gpt-5.6-luna", plugin); err != nil {
		t.Fatalf("未訳の再送信: %v", err)
	}
	if fb.submits != 2 {
		t.Errorf("送信回数 = %d, want 2", fb.submits)
	}
}

// R-2-1, R-2-2: 準備済みの再送信は保存済みの口調を使い、未訳 1 件だけを外部 batch へ載せる。
// 既に訳がある行は未訳一覧に入らないため、送信も書き換えも行わない。
func TestBatchRetryUsesOnlyPendingRowsWithoutPersonaRegeneration(t *testing.T) {
	ctx := context.Background()
	const plugin = "retry.esp"
	store := &fakeStore{
		syncRetryReady: true,
		untranslated: []model.Narration{
			{ID: 1, Plugin: plugin, Rec: "BOOK", Field: "DESC", Source: "done", Dest: "訳済み", Status: 3},
			{ID: 2, Plugin: plugin, Rec: "BOOK", Field: "DESC", Source: "pending", Status: 0},
		},
	}
	fb := &fakeBatchProvider{out: map[string]string{"pending": "未訳の訳"}}
	runner := NewBatchRunner(New(store, &fakeTranslator{}, fakeLexicon{}, nil, nil), map[string]provider.BatchTranslator{
		provider.BatchProviderXAI: fb,
	}, store)

	outcome, err := runner.SubmitBatch(ctx, provider.BatchProviderXAI, provider.Connection{}, "grok", plugin)
	if err != nil {
		t.Fatalf("準備済み batch の再送信: %v", err)
	}
	if outcome.CompletedWithoutExternalBatch {
		t.Fatal("外部 batch が必要な未訳を即時完了として返した")
	}
	if store.generateInputCalls != 0 {
		t.Fatalf("準備済みの再送信で口調を再集計した: %d 回", store.generateInputCalls)
	}
	if fb.submits != 1 || len(store.batchReqs) != 1 || store.batchReqs[0].Kind != model.BatchKindNarration || store.batchReqs[0].RowID != 2 {
		t.Fatalf("外部 batch の対象 = submits:%d requests:%+v", fb.submits, store.batchReqs)
	}
	if store.untranslated[0].Dest != "訳済み" || store.untranslated[0].Status != 3 {
		t.Fatalf("既訳行を変更した: %+v", store.untranslated[0])
	}
}

// R-2-4: 準備済みの未訳を既訳だけで解決できる場合は、外部 batch を作らず完了する。
func TestBatchRetryCompletesWithoutExternalBatchWhenReferencesResolveAllRows(t *testing.T) {
	ctx := context.Background()
	const plugin = "reference-only.esp"
	store := &fakeStore{
		syncRetryReady: true,
		untranslated:   []model.Narration{{ID: 1, Plugin: plugin, Rec: "BOOK", Field: "DESC", Source: "known"}},
		refs:           []model.ReferenceTranslation{{Rec: "BOOK", Field: "DESC", Source: "known", Dest: "既訳"}},
	}
	fb := &fakeBatchProvider{}
	runner := NewBatchRunner(New(store, &fakeTranslator{}, fakeLexicon{}, nil, nil), map[string]provider.BatchTranslator{
		provider.BatchProviderXAI: fb,
	}, store)

	outcome, err := runner.SubmitBatch(ctx, provider.BatchProviderXAI, provider.Connection{}, "grok", plugin)
	if err != nil {
		t.Fatalf("既訳だけの再送信: %v", err)
	}
	if !outcome.CompletedWithoutExternalBatch || fb.submits != 0 || store.batchProg.Stage != model.BatchStageDone {
		t.Fatalf("即時完了 = outcome:%+v submits:%d progression:%+v", outcome, fb.submits, store.batchProg)
	}
	if len(store.updates) != 1 || store.updates[0] != (update{1, "既訳", 1}) {
		t.Fatalf("既訳の適用 = %v", store.updates)
	}
	if store.generateInputCalls != 0 {
		t.Fatalf("即時完了の再送信で口調を再集計した: %d 回", store.generateInputCalls)
	}
}

// updatesByID は書き戻しログを id→update の map へ畳む（順序に依存せず集合として照合する）。
func updatesByID(logs []update) map[int64]update {
	m := make(map[int64]update, len(logs))
	for _, u := range logs {
		m[u.id] = u
	}
	return m
}

// seedFor は同期・batch で同一の未訳データを積んだ fakeStore を作る。
// 固有名（proper）→ 本文への機械置換合流、既存訳流用、口調なし台詞を 1 経路ずつ含める。
func seedFor(plugin string) *fakeStore {
	return &fakeStore{
		proper: []model.ProperNoun{
			{ID: 1, Plugin: plugin, Source: "Riften", Category: "NPC_"},
		},
		untranslated: []model.Narration{
			{ID: 1, Plugin: plugin, Rec: "BOOK", Field: "DESC", Source: "See Riften."},   // 固有名を機械置換してから AI
			{ID: 2, Plugin: plugin, Rec: "WEAP", Field: "DESC", Source: "A fine blade."}, // 既存訳流用（AI を呼ばない）
		},
		lines: []model.Line{
			{ID: 1, Plugin: plugin, Rec: "INFO", Field: "NAM1", Source: "Hello."}, // 口調なし台詞
		},
		refs: []model.ReferenceTranslation{
			{Rec: "WEAP", Field: "DESC", Source: "A fine blade.", Dest: "見事な刃。"},
		},
	}
}

// translationOut は同期・batch 共有の訳（user メッセージ→訳文）。
// "See Riften." は英語本文へ参考語を付けた本文promptに対して引く。
func translationOut() map[string]string {
	return map[string]string{
		"Riften":      "リフテン",
		"See Riften.": "リフテンを見よ。",
		"Hello.":      "こんにちは。",
	}
}

// batch 経路が、同期経路と同一の dest・訳状態を書き戻すことを表明する（外から見て同期と batch が変わらない）。
// 併せて、固有名 batch → 本文 batch の 2 段連鎖と、送信時点では dest を確定しない（反映まで遅延する）ことを確かめる。
func TestBatchMatchesSyncEndToEnd(t *testing.T) {
	ctx := context.Background()
	conn := provider.Connection{Endpoint: "http://x", APIKey: "k"}
	const plugin = "A.esp"
	out := translationOut()

	// --- 同期経路（基準） ---
	syncStore := seedFor(plugin)
	syncEng := New(syncStore, &fakeTranslator{out: out}, fakeLexicon{}, nil, nil)
	if _, err := syncEng.Run(ctx, conn, "grok", plugin, nil); err != nil {
		t.Fatalf("同期 Run: %v", err)
	}

	// --- batch 経路 ---
	batchStore := seedFor(plugin)
	fb := &fakeBatchProvider{out: out}
	runner := NewBatchRunner(New(batchStore, &fakeTranslator{out: out}, fakeLexicon{}, nil, nil), map[string]provider.BatchTranslator{provider.BatchProviderXAI: fb}, batchStore)

	// 送信: 固有名 batch を送るが、dest はまだ確定しない（2 時点の遅延）。
	if _, err := runner.SubmitBatch(ctx, provider.BatchProviderXAI, conn, "grok", plugin); err != nil {
		t.Fatalf("batch 送信: %v", err)
	}
	if len(batchStore.properUpdates) != 0 || len(batchStore.updates) != 0 || len(batchStore.lineUpdates) != 0 {
		t.Fatalf("送信時点で dest を確定した（反映まで遅延すべき）: proper=%v narr=%v line=%v",
			batchStore.properUpdates, batchStore.updates, batchStore.lineUpdates)
	}
	if batchStore.batchProg.Stage != model.BatchStageProperNoun || batchStore.batchProg.ProperBatchID == "" {
		t.Fatalf("送信後の進行が固有名段でない: %+v", batchStore.batchProg)
	}

	// 反映①: 固有名 batch が完了 → 固有名を確定し、本文 batch を送る（進行段=本文）。
	if err := runner.RefreshPlugin(ctx, provider.BatchProviderXAI, conn, plugin); err != nil {
		t.Fatalf("batch 反映①: %v", err)
	}
	if batchStore.batchProg.Stage != model.BatchStageBody || batchStore.batchProg.BodyBatchID == "" {
		t.Fatalf("反映①後の進行が本文段でない: %+v", batchStore.batchProg)
	}
	if len(batchStore.properUpdates) != 1 {
		t.Fatalf("固有名が確定していない: %v", batchStore.properUpdates)
	}
	if !batchStore.syncRetryReady {
		t.Fatal("本文計画の完了後に同期再実行の準備完了が保存されていない")
	}

	// 反映②: アプリ再起動相当（永続から続ける）。新しい engine・runner を組み直しても反映が続く。
	restartEng := New(batchStore, &fakeTranslator{out: out}, fakeLexicon{}, nil, nil)
	restartRunner := NewBatchRunner(restartEng, map[string]provider.BatchTranslator{provider.BatchProviderXAI: fb}, batchStore)
	if err := restartRunner.RefreshPlugin(ctx, provider.BatchProviderXAI, conn, plugin); err != nil {
		t.Fatalf("batch 反映②（再起動相当）: %v", err)
	}
	if batchStore.batchProg.Stage != model.BatchStageDone {
		t.Fatalf("反映②後に進行が完了していない: %+v", batchStore.batchProg)
	}

	// 固有名 batch と本文 batch の 2 回送信になったこと。
	if fb.submits != 2 {
		t.Errorf("batch 送信回数 = %d, want 2（固有名→本文）", fb.submits)
	}

	// 外から見て同期と batch が一致すること（dest・訳状態を集合として照合）。
	assertSameUpdates(t, "固有名", syncStore.properUpdates, batchStore.properUpdates)
	assertSameUpdates(t, "叙述文", syncStore.updates, batchStore.updates)
	assertSameUpdates(t, "台詞", syncStore.lineUpdates, batchStore.lineUpdates)

	// 具体値も固定する（固有名の機械置換・既存訳流用・AI 訳の 3 経路が batch でも同じ結果になる）。
	bn := updatesByID(batchStore.updates)
	if bn[1] != (update{1, "リフテンを見よ。", 3}) {
		t.Errorf("固有名機械置換つき叙述文の batch 結果 = %v", bn[1])
	}
	if bn[2] != (update{2, "見事な刃。", 1}) {
		t.Errorf("既存訳流用の batch 結果 = %v（status=1 で流用すべき）", bn[2])
	}
	if updatesByID(batchStore.properUpdates)[1] != (update{1, "リフテン", 3}) {
		t.Errorf("固有名 AI 訳の batch 結果 = %v", batchStore.properUpdates)
	}
}

// 個別リクエストの失敗（skippable）を、batch も同期と同じく未訳のまま据え置く（再送信で回収）ことを表明する。
func TestBatchLeavesUntranslatedOnFailureLikeSync(t *testing.T) {
	ctx := context.Background()
	conn := provider.Connection{Endpoint: "http://x"}
	const plugin = "B.esp"
	out := map[string]string{"ok.": "はい。"}
	// "fail." はサーバ一時失敗を返す（同期・batch とも据え置き）。
	failErr := map[string]error{"fail.": fmt.Errorf("%w: status 503", provider.ErrServerTransient)}

	seed := func() *fakeStore {
		return &fakeStore{lines: []model.Line{
			{ID: 1, Plugin: plugin, Rec: "INFO", Field: "NAM1", Source: "ok."},
			{ID: 2, Plugin: plugin, Rec: "INFO", Field: "NAM1", Source: "fail."},
		}}
	}

	// 同期。
	syncStore := seed()
	syncEng := New(syncStore, &fakeTranslator{out: out, errByUser: failErr}, fakeLexicon{}, nil, nil)
	if _, err := syncEng.Run(ctx, conn, "grok", plugin, nil); err != nil {
		t.Fatalf("同期 Run: %v", err)
	}

	// batch（固有名なしのため submit で本文 batch まで進み、反映 1 回で確定）。
	batchStore := seed()
	fb := &fakeBatchProvider{out: out, errByUser: failErr}
	runner := NewBatchRunner(New(batchStore, &fakeTranslator{out: out}, fakeLexicon{}, nil, nil), map[string]provider.BatchTranslator{provider.BatchProviderXAI: fb}, batchStore)
	if _, err := runner.SubmitBatch(ctx, provider.BatchProviderXAI, conn, "grok", plugin); err != nil {
		t.Fatalf("batch 送信: %v", err)
	}
	if err := runner.RefreshPlugin(ctx, provider.BatchProviderXAI, conn, plugin); err != nil {
		t.Fatalf("batch 反映: %v", err)
	}

	// 成功した ok.(ID=1) だけ確定し、失敗した fail.(ID=2) は据え置き。同期と一致すること。
	assertSameUpdates(t, "台詞（失敗据え置き）", syncStore.lineUpdates, batchStore.lineUpdates)
	bl := updatesByID(batchStore.lineUpdates)
	if _, applied := bl[2]; applied {
		t.Errorf("失敗した台詞 ID=2 を確定した（据え置くべき）: %v", batchStore.lineUpdates)
	}
	if bl[1] != (update{1, "はい。", 3}) {
		t.Errorf("成功した台詞 ID=1 の batch 結果 = %v", bl[1])
	}
}

// 固有名 batch の送信 HTTP 失敗で進行が半端に残っても、再送信が拒否されず reset して回復し、
// 反映まで通常どおり dest を確定できることを表明する（送信失敗の詰まりの復旧経路）。
func TestBatchResubmitRecoversFromSubmitFailure(t *testing.T) {
	ctx := context.Background()
	conn := provider.Connection{Endpoint: "http://x", APIKey: "k"}
	const plugin = "C.esp"
	out := translationOut()

	batchStore := seedFor(plugin)
	fb := &fakeBatchProvider{out: out, failSubmitsN: 1} // 最初の固有名送信を失敗させる。
	runner := NewBatchRunner(New(batchStore, &fakeTranslator{out: out}, fakeLexicon{}, nil, nil), map[string]provider.BatchTranslator{provider.BatchProviderXAI: fb}, batchStore)

	// 送信①: 固有名 batch の送信が失敗する。進行は固有名段のまま外部 ID が空で残る（半端な進行）。
	if _, err := runner.SubmitBatch(ctx, provider.BatchProviderXAI, conn, "grok", plugin); err == nil {
		t.Fatalf("送信失敗を模したが SubmitBatch が成功した")
	}
	if batchStore.batchProg == nil || batchStore.batchProg.Stage != model.BatchStageProperNoun || batchStore.batchProg.ProperBatchID != "" {
		t.Fatalf("送信失敗後の進行が想定外: %+v", batchStore.batchProg)
	}
	// 送信失敗時点で dest を確定していないこと。
	if len(batchStore.properUpdates) != 0 || len(batchStore.updates) != 0 || len(batchStore.lineUpdates) != 0 {
		t.Fatalf("送信失敗時点で dest を確定した: proper=%v narr=%v line=%v",
			batchStore.properUpdates, batchStore.updates, batchStore.lineUpdates)
	}

	// 送信②: 半端な進行は拒否されず reset されて成功する（BlocksResubmit が外部 ID 空を許すため）。
	if _, err := runner.SubmitBatch(ctx, provider.BatchProviderXAI, conn, "grok", plugin); err != nil {
		t.Fatalf("再送信が拒否・失敗した（回復できるべき）: %v", err)
	}
	if batchStore.batchProg.Stage != model.BatchStageProperNoun || batchStore.batchProg.ProperBatchID == "" {
		t.Fatalf("再送信後の進行が固有名段でない: %+v", batchStore.batchProg)
	}

	// 反映で通常どおり 2 段連鎖して完了し、dest が確定すること。
	if err := runner.RefreshPlugin(ctx, provider.BatchProviderXAI, conn, plugin); err != nil {
		t.Fatalf("batch 反映①: %v", err)
	}
	if err := runner.RefreshPlugin(ctx, provider.BatchProviderXAI, conn, plugin); err != nil {
		t.Fatalf("batch 反映②: %v", err)
	}
	if batchStore.batchProg.Stage != model.BatchStageDone {
		t.Fatalf("反映後に進行が完了していない: %+v", batchStore.batchProg)
	}
	if updatesByID(batchStore.properUpdates)[1] != (update{1, "リフテン", 3}) {
		t.Errorf("回復後の固有名 AI 訳 = %v", batchStore.properUpdates)
	}
	if updatesByID(batchStore.updates)[1] != (update{1, "リフテンを見よ。", 3}) {
		t.Errorf("回復後の叙述文 = %v", batchStore.updates)
	}
}

// 本文段の外部送信が失敗しても、進行を本文段に残し、同じcustom_idで再送できることを確かめる。
func TestBatchResubmitsFailedBodyStage(t *testing.T) {
	ctx := context.Background()
	conn := provider.Connection{Endpoint: "http://x", APIKey: "k"}
	const plugin = "body-failed.esp"
	store := seedFor(plugin)
	fb := &fakeBatchProvider{out: translationOut(), failAtSubmit: 2}
	reader := fakePrebuiltDictionary{references: []model.PrebuiltDictionaryReference{{Source: "Riften", Dest: "リフテン", PartOfSpeech: "noun", Meaning: "城塞を守る都市", SkyrimCategory: "city"}}}
	runner := NewBatchRunner(New(store, &fakeTranslator{out: translationOut()}, fakeLexicon{}, nil, nil, reader), map[string]provider.BatchTranslator{provider.BatchProviderXAI: fb}, store)

	if _, err := runner.SubmitBatch(ctx, provider.BatchProviderXAI, conn, "grok", plugin); err != nil {
		t.Fatalf("固有名段の送信: %v", err)
	}
	if err := runner.RefreshPlugin(ctx, provider.BatchProviderXAI, conn, plugin); err == nil {
		t.Fatal("本文段の送信失敗を返さなかった")
	}
	if store.batchProg.Stage != model.BatchStageBody || store.batchProg.BodyBatchID != "" {
		t.Fatalf("本文送信失敗後の進行 = %+v", store.batchProg)
	}
	if len(store.batchReqs) < 2 || store.batchReqs[len(store.batchReqs)-1].SendState != "failed" {
		t.Fatalf("本文requestがfailed状態でない: %+v", store.batchReqs)
	}
	customID := store.batchReqs[len(store.batchReqs)-1].CustomID
	referencesJSON := store.batchReqs[len(store.batchReqs)-1].ReferencesJSON
	promptHash := store.batchReqs[len(store.batchReqs)-1].PromptHash
	// 再起動相当のrunner再生成後も、仮状態を同じ進行から再送する。
	runner = NewBatchRunner(New(store, &fakeTranslator{out: translationOut()}, fakeLexicon{}, nil, nil, reader), map[string]provider.BatchTranslator{provider.BatchProviderXAI: fb}, store)
	if _, err := runner.SubmitBatch(ctx, provider.BatchProviderXAI, conn, "grok", plugin); err != nil {
		t.Fatalf("本文段の再送: %v", err)
	}
	if store.batchProg.Stage != model.BatchStageBody || store.batchProg.BodyBatchID == "" {
		t.Fatalf("本文段の再送後の進行 = %+v", store.batchProg)
	}
	if store.batchReqs[len(store.batchReqs)-1].CustomID != customID || store.batchReqs[len(store.batchReqs)-1].SendState != "submitted" {
		t.Fatalf("本文段の再送requestが変わった: %+v", store.batchReqs)
	}
	if store.batchReqs[len(store.batchReqs)-1].ReferencesJSON != referencesJSON || store.batchReqs[len(store.batchReqs)-1].PromptHash != promptHash {
		t.Fatalf("本文段の参考語snapshotまたはhashが変わった: %+v", store.batchReqs)
	}
}

func TestBatchStopsBeforeProviderWhenPrebuiltValidationFails(t *testing.T) {
	ctx := context.Background()
	store := seedFor("invalid-reader.esp")
	fb := &fakeBatchProvider{out: translationOut()}
	runner := NewBatchRunner(New(store, &fakeTranslator{out: translationOut()}, fakeLexicon{}, nil, nil, fakePrebuiltDictionary{err: errors.New("dictionary unavailable")}), map[string]provider.BatchTranslator{provider.BatchProviderXAI: fb}, store)
	if _, err := runner.SubmitBatch(ctx, provider.BatchProviderXAI, provider.Connection{APIKey: "k"}, "grok", "invalid-reader.esp"); err == nil {
		t.Fatal("reader検証失敗を返さなかった")
	}
	if fb.submits != 0 || store.batchProg != nil || len(store.batchReqs) != 0 {
		t.Fatalf("reader失敗後にbatch送信または保存を開始した: submits=%d progression=%+v requests=%+v", fb.submits, store.batchProg, store.batchReqs)
	}
}

// assertSameUpdates は 2 つの書き戻しログが集合（id→dest/status）として一致することを表明する。
func assertSameUpdates(t *testing.T, label string, want, got []update) {
	t.Helper()
	w, g := updatesByID(want), updatesByID(got)
	if len(w) != len(g) {
		t.Fatalf("%s: 書き戻し件数が違う sync=%v batch=%v", label, want, got)
	}
	for id, wu := range w {
		gu, ok := g[id]
		if !ok {
			t.Errorf("%s: batch に id=%d の書き戻しが無い（sync=%v）", label, id, wu)
			continue
		}
		if gu != wu {
			t.Errorf("%s: id=%d sync=%v batch=%v", label, id, wu, gu)
		}
	}
}
