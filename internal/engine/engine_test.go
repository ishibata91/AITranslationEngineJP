package engine

import (
	"context"
	"errors"
	"strings"
	"testing"

	"aitranslationenginejp/internal/model"
	"aitranslationenginejp/internal/provider"
)

// testBaseDirective・testPersonaTemplate はテスト用のプロンプトテンプレート既定値。
// 実際の seed 文面（migration 0004）に依存せず、テンプレート駆動の組み立てを確かめるために使う。
const (
	testBaseDirective   = "あなたは Skyrim Mod の翻訳者です。訳文だけを出力してください。"
	testPersonaTemplate = "この台詞の話者の人物像:\n{traits}\nこの人物像に合う口調と人称で訳すこと。"
)

type fakeStore struct {
	untranslated []model.Narration
	lines        []model.Line
	linePersonas map[int64]model.LinePersonaInput // lineID → 生成済み基底口調（無ければ口調なし）
	terms        []model.MasterTerm               // 固有名の機械置換辞書（無ければ置換なし）
	tmpl         model.PromptTemplate             // プロンプトテンプレート（未設定ならテスト用既定値を返す）
	genSources   []model.SpeakerLineSource        // 一括生成が読む入力（Run テストでは空で生成 no-op）
	updates      []update
	lineUpdates  []update
	// loadPersonasCalls は LoadLinePersonas の呼び出し回数。注入の引きが台詞数 N 非依存（N+1 廃止）の観測に使う。
	loadPersonasCalls int
}

func (f *fakeStore) GetPromptTemplate(_ context.Context) (model.PromptTemplate, error) {
	if f.tmpl.BaseDirective == "" && f.tmpl.PersonaTemplate == "" {
		return model.PromptTemplate{BaseDirective: testBaseDirective, PersonaTemplate: testPersonaTemplate}, nil
	}
	return f.tmpl, nil
}

type update struct {
	id     int64
	dest   string
	status int
}

func (f *fakeStore) ListUntranslatedNarrations(_ context.Context) ([]model.Narration, error) {
	return f.untranslated, nil
}

func (f *fakeStore) UpdateNarrationDest(_ context.Context, id int64, dest string, status int) error {
	f.updates = append(f.updates, update{id: id, dest: dest, status: status})
	return nil
}

func (f *fakeStore) ListUntranslatedLines(_ context.Context) ([]model.Line, error) {
	return f.lines, nil
}

func (f *fakeStore) UpdateLineDest(_ context.Context, id int64, dest string, status int) error {
	f.lineUpdates = append(f.lineUpdates, update{id: id, dest: dest, status: status})
	return nil
}

func (f *fakeStore) ListMasterTerms(_ context.Context) ([]model.MasterTerm, error) {
	return f.terms, nil
}

func (f *fakeStore) InsertDerivedTerms(_ context.Context, terms []model.MasterTerm) (int, error) {
	return len(terms), nil
}

// --- PersonaStore（生成入力・キャッシュ・保存・注入） ---

func (f *fakeStore) ListSpeakerLineSources(_ context.Context) ([]model.SpeakerLineSource, error) {
	return f.genSources, nil
}

func (f *fakeStore) GetLineAnalyses(_ context.Context, _ []string) (map[string]model.LineAnalysis, error) {
	return map[string]model.LineAnalysis{}, nil
}

func (f *fakeStore) UpsertLineAnalysis(_ context.Context, _ model.LineAnalysis) error { return nil }

func (f *fakeStore) UpsertPersonaCharacter(_ context.Context, _ model.PersonaCharacter) error {
	return nil
}

func (f *fakeStore) LoadLinePersonas(_ context.Context, lineIDs []int64) (map[int64]model.LinePersonaInput, error) {
	f.loadPersonasCalls++
	out := make(map[int64]model.LinePersonaInput)
	for _, id := range lineIDs {
		if in, ok := f.linePersonas[id]; ok {
			out[id] = in
		}
	}
	return out, nil
}

type fakeTranslator struct {
	out        map[string]string // user メッセージ（機械置換済み原文）→ 訳文
	gotModel   string
	gotPrompts map[string]provider.Prompt // user メッセージ → engine が組んで送った完成プロンプト
	err        error
}

func (f *fakeTranslator) Translate(_ context.Context, _ provider.Connection, model string, prompt provider.Prompt) (string, error) {
	f.gotModel = model
	if f.gotPrompts == nil {
		f.gotPrompts = map[string]provider.Prompt{}
	}
	f.gotPrompts[prompt.User] = prompt
	if f.err != nil {
		return "", f.err
	}
	return f.out[prompt.User], nil
}

func (f *fakeTranslator) ListModels(_ context.Context, _ provider.Connection) ([]string, error) {
	return nil, nil
}

// 未訳の叙述文を順に翻訳し、訳文を status=3（仮訳）で書き戻すこと。叙述文はペルソナ指示を付けない。
func TestRunTranslatesUntranslatedAsProvisional(t *testing.T) {
	store := &fakeStore{untranslated: []model.Narration{
		{ID: 1, Source: "halls"},
		{ID: 2, Source: "cairn"},
	}}
	tr := &fakeTranslator{out: map[string]string{"halls": "広間", "cairn": "ケルン"}}
	eng := New(store, tr, fakeLexicon{}, nil)

	count, err := eng.Run(context.Background(), provider.Connection{Endpoint: "http://x"}, "model-x", nil)
	if err != nil {
		t.Fatalf("Run: %v", err)
	}
	if count != 2 {
		t.Errorf("count = %d, want 2", count)
	}
	if tr.gotModel != "model-x" {
		t.Errorf("model passed = %q, want model-x", tr.gotModel)
	}
	// 叙述文は口調指示を持たないため、system はテンプレートの base 指示だけになる（口調指示の合成なし）。
	if sys := tr.gotPrompts["halls"].System; sys != testBaseDirective {
		t.Errorf("叙述文の system に口調指示が合成された: %q", sys)
	}
	want := []update{{1, "広間", 3}, {2, "ケルン", 3}}
	if len(store.updates) != len(want) {
		t.Fatalf("updates = %v, want %v", store.updates, want)
	}
	for i := range want {
		if store.updates[i] != want[i] {
			t.Errorf("update[%d] = %v, want %v", i, store.updates[i], want[i])
		}
	}
}

// 本文翻訳の前に、原文中の固有名を辞書の確定訳語へ機械置換してから AI へ渡すこと。
func TestRunReplacesTermsBeforeTranslate(t *testing.T) {
	store := &fakeStore{
		untranslated: []model.Narration{{ID: 1, Source: "The Riften guard waited."}},
		terms:        []model.MasterTerm{{Source: "Riften", Dest: "リフテン"}},
	}
	tr := &fakeTranslator{out: map[string]string{
		"The リフテン guard waited.": "リフテンの衛兵が待っていた。",
	}}
	eng := New(store, tr, fakeLexicon{}, nil)

	if _, err := eng.Run(context.Background(), provider.Connection{}, "m", nil); err != nil {
		t.Fatalf("Run: %v", err)
	}
	// AI へ渡した user メッセージ（原文）は固有名が置換済みであること（置換前の原文では渡らない）。
	if _, ok := tr.gotPrompts["The リフテン guard waited."]; !ok {
		t.Errorf("置換後の原文が AI へ渡っていない。gotPrompts=%v", tr.gotPrompts)
	}
	// 書き戻した訳文は置換済み原文に対する訳であること。
	want := []update{{1, "リフテンの衛兵が待っていた。", 3}}
	if len(store.updates) != 1 || store.updates[0] != want[0] {
		t.Errorf("updates = %v, want %v", store.updates, want)
	}
}

// provider がエラーを返したら、そのエラーを返すこと（壊れた訳を書き戻さない）。
func TestRunReturnsProviderError(t *testing.T) {
	store := &fakeStore{untranslated: []model.Narration{{ID: 1, Source: "x"}}}
	tr := &fakeTranslator{err: errors.New("connection refused")}
	eng := New(store, tr, fakeLexicon{}, nil)

	_, err := eng.Run(context.Background(), provider.Connection{}, "m", nil)
	if err == nil {
		t.Fatal("want error, got nil")
	}
	if len(store.updates) != 0 {
		t.Errorf("should not write dest on provider error, got %v", store.updates)
	}
}

// 生成済みペルソナを持つ台詞は、基底口調の性質文と種族訛りが翻訳プロンプトへ注入されること。
func TestRunTranslatesLinesWithPersonaDirective(t *testing.T) {
	store := &fakeStore{
		lines: []model.Line{{ID: 10, Source: "mother?"}},
		linePersonas: map[int64]model.LinePersonaInput{
			// 丁寧×中（物腰やわ）＋ Khajiit 種族訛り。
			10: {AttitudeBand: 2, EmotionBand: 1, RaceEDID: "KhajiitRace"},
		},
	}
	tr := &fakeTranslator{out: map[string]string{"mother?": "母さん？"}}
	eng := New(store, tr, fakeLexicon{}, nil)

	count, err := eng.Run(context.Background(), provider.Connection{}, "m", nil)
	if err != nil {
		t.Fatalf("Run: %v", err)
	}
	if count != 1 {
		t.Errorf("count = %d, want 1", count)
	}
	// 口調指示は完成プロンプトの system へ合成される（base 指示の後ろに続く）。
	system := tr.gotPrompts["mother?"].System
	if !strings.Contains(system, "柔らかく丁寧") {
		t.Errorf("基底口調の性質文が system に無い: %q", system)
	}
	if !strings.Contains(system, "三人称") {
		t.Errorf("種族訛り（カジート三人称）が system に無い: %q", system)
	}
	want := []update{{10, "母さん？", 3}}
	if len(store.lineUpdates) != 1 || store.lineUpdates[0] != want[0] {
		t.Errorf("lineUpdates = %v, want %v", store.lineUpdates, want)
	}
	// Run の台詞翻訳ループでも注入の引きは一括 1 回（ループ内の個別問い合わせを廃止）。
	if store.loadPersonasCalls != 1 {
		t.Errorf("Run の注入引き = %d 回, want 1（台詞数に非依存）", store.loadPersonasCalls)
	}
}

// 生成済みペルソナを持たない台詞は、口調指示を注入せず空 directive で訳すこと。
func TestRunTranslatesLineWithoutSpeaker(t *testing.T) {
	store := &fakeStore{lines: []model.Line{{ID: 20, Source: "door"}}}
	tr := &fakeTranslator{out: map[string]string{"door": "扉"}}
	eng := New(store, tr, fakeLexicon{}, nil)

	if _, err := eng.Run(context.Background(), provider.Connection{}, "m", nil); err != nil {
		t.Fatalf("Run: %v", err)
	}
	// ペルソナなしの台詞は口調指示を合成せず、system はテンプレートの base 指示だけになる。
	if sys := tr.gotPrompts["door"].System; sys != testBaseDirective {
		t.Errorf("ペルソナなしの台詞の system に口調指示が合成された: %q", sys)
	}
}

// 進捗 callback が、叙述文と台詞の合計を total に、処理ごとに done を増やして通知すること。
func TestRunReportsProgress(t *testing.T) {
	store := &fakeStore{
		untranslated: []model.Narration{{ID: 1, Source: "a"}},
		lines:        []model.Line{{ID: 10, Source: "b"}, {ID: 11, Source: "c"}},
	}
	tr := &fakeTranslator{out: map[string]string{"a": "あ", "b": "び", "c": "し"}}
	eng := New(store, tr, fakeLexicon{}, nil)

	var seen [][2]int
	_, err := eng.Run(context.Background(), provider.Connection{}, "m", func(done, total int) {
		seen = append(seen, [2]int{done, total})
	})
	if err != nil {
		t.Fatalf("Run: %v", err)
	}
	// 初回 0/3 と、3 件処理ごとの 1/3,2/3,3/3。
	want := [][2]int{{0, 3}, {1, 3}, {2, 3}, {3, 3}}
	if len(seen) != len(want) {
		t.Fatalf("progress = %v, want %v", seen, want)
	}
	for i := range want {
		if seen[i] != want[i] {
			t.Errorf("progress[%d] = %v, want %v", i, seen[i], want[i])
		}
	}
}

// LinePersonas は生成済みペルソナを持つ台詞の口調指示と短い要約を map で返し、無い台詞は map に現れないこと。
func TestLinePersonas(t *testing.T) {
	store := &fakeStore{
		linePersonas: map[int64]model.LinePersonaInput{
			// 尊大×抑制（冷然・見下し）。
			10: {AttitudeBand: 0, EmotionBand: 0, RaceEDID: "NordRace"},
		},
	}
	eng := New(store, &fakeTranslator{}, fakeLexicon{}, nil)

	personas, err := eng.LinePersonas(context.Background(), []int64{10, 99}, testPersonaTemplate) // 99 はペルソナなし
	if err != nil {
		t.Fatalf("LinePersonas: %v", err)
	}
	p, ok := personas[10]
	if !ok {
		t.Fatalf("ペルソナありの台詞 10 が persona map に無い")
	}
	if !strings.Contains(p.Directive, "相手を見下す冷たい口調") {
		t.Errorf("directive = %q", p.Directive)
	}
	if p.Label != "口調: 冷然・見下し" {
		t.Errorf("label = %q", p.Label)
	}
	if _, ok := personas[99]; ok {
		t.Errorf("ペルソナなしの台詞 99 が persona map に現れた")
	}
}

// 注入の引きの DB 呼び出しが台詞数 N に依存せず一括（定数 1 回）であること（N+1 廃止の観測）。
func TestLinePersonasBulkLoadsOnce(t *testing.T) {
	for _, n := range []int{2, 50} {
		store := &fakeStore{linePersonas: map[int64]model.LinePersonaInput{}}
		ids := make([]int64, n)
		for i := range n {
			id := int64(i + 1)
			ids[i] = id
			store.linePersonas[id] = model.LinePersonaInput{AttitudeBand: 1, EmotionBand: 1}
		}
		eng := New(store, &fakeTranslator{}, fakeLexicon{}, nil)
		if _, err := eng.LinePersonas(context.Background(), ids, testPersonaTemplate); err != nil {
			t.Fatalf("LinePersonas(N=%d): %v", n, err)
		}
		if store.loadPersonasCalls != 1 {
			t.Errorf("N=%d: LoadLinePersonas 呼び出し = %d, want 1（台詞数に非依存）", n, store.loadPersonasCalls)
		}
	}
}
