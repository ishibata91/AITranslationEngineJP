package engine

import (
	"context"
	"errors"
	"fmt"
	"sort"
	"strings"
	"testing"

	"aitranslationenginejp/internal/core/rolespeech"
	"aitranslationenginejp/internal/core/tone"
	"aitranslationenginejp/internal/model"
	"aitranslationenginejp/internal/provider"
)

// testBaseDirective は base 指示、testPersonaTemplate は directive「口調」のテスト用既定値。
// 実際の seed 文面（migration 0004）に依存せず、テンプレート駆動の組み立てを確かめるために使う。
const (
	testBaseDirective   = "あなたは Skyrim Mod の翻訳者です。訳文だけを出力してください。"
	testPersonaTemplate = "この台詞の話者の人物像:\n{traits}\nこの人物像に合う口調と人称で訳すこと。"
)

type fakeStore struct {
	untranslated       []model.Narration
	lines              []model.Line
	proper             []model.ProperNoun               // 固有名（ListProperNouns / 未訳は status=0 を ListUntranslatedProperNouns）
	linePersonas       map[int64]model.LinePersonaInput // lineID → 生成済み基底口調（無ければ口調なし）
	terms              []model.MasterTerm               // 固有名の機械置換辞書（無ければ置換なし）
	tmpl               model.PromptTemplate             // プロンプトテンプレート（未設定ならテスト用既定値を返す）
	directives         []model.Directive                // 指示文（nil ならテスト用既定: 口調=testPersonaTemplate）
	recordTypes        []model.RecordType               // REC:FIELD 割り当て（既定は空。叙述文 directive テストで設定）
	extracted          []model.ExtractedField           // 取込段の入力（既定は空）
	genSources         []model.SpeakerLineSource        // 一括生成が読む入力（Run テストでは空で生成 no-op）
	generateInputCalls int
	updates            []update
	lineUpdates        []update
	properUpdates      []update           // UpdateProperNounDest の記録（固有名フェーズの観測）
	derivedPropers     []model.ProperNoun // InsertDerivedProperNouns で投入された人名の部分形（派生段の観測）
	insertedTerms      []model.MasterTerm // InsertDerivedTerms で投入された横断辞書の行（昇格しない不変境界の観測）
	ingestedNarr       []model.Narration  // IngestNarrations で投入された行
	ingestedPN         []model.ProperNoun // IngestProperNouns で投入された行
	ingestedLine       []model.Line       // IngestLines で投入された行
	linkCalled         bool               // LinkLineSpeakersFromStaging が呼ばれたか
	condLinkCalled     bool               // LinkLineConditionsFromStaging が呼ばれたか
	emoLinkCalled      bool               // LinkLineEmotionsFromStaging が呼ばれたか
	lineConditions     map[int64]string   // LoadLineConditions が返す条件由来の性別（lineID→sex）
	// 言及段の観測。allNarrations / allLines は ListNarrations / ListLines が返す全行（nil なら空）。
	allNarrations     []model.Narration
	allLines          []model.Line
	narrationMentions []model.NarrationMention // InsertNarrationMentions で投入された行
	lineMentions      []model.LineMention      // InsertLineMentions で投入された行
	describedCalled   bool                     // LinkNarrationDescribed が呼ばれたか
	// loadPersonasCalls は LoadLinePersonas の呼び出し回数。注入の引きが台詞数 N 非依存（N+1 廃止）の観測に使う。
	loadPersonasCalls int
	// refs は既存訳（参照訳）。ListReferenceTranslations が返す（既定は空＝流用なし）。
	refs []model.ReferenceTranslation
	// batch 進行の永続（BatchStore 実装用、batch 結合テストで使う）。1 plugin 1 進行を想定。
	batchProg      *model.BatchTranslation
	batchReqs      []model.BatchRequest
	nextBatchID    int64
	syncRetryReady bool
}

type fakePrebuiltDictionary struct {
	references []model.PrebuiltDictionaryReference
	err        error
}

func (f fakePrebuiltDictionary) ValidatePrebuiltDictionary(_ context.Context) error { return f.err }
func (f fakePrebuiltDictionary) References(_ context.Context) ([]model.PrebuiltDictionaryReference, error) {
	return f.references, f.err
}

func (f *fakeStore) UpsertTranslationReferenceSnapshot(_ context.Context, _ model.TranslationReferenceSnapshot) error {
	return nil
}

func (f *fakeStore) GetTranslationReferenceSnapshot(_ context.Context, _ string, _ string, _ int64) (model.TranslationReferenceSnapshot, bool, error) {
	return model.TranslationReferenceSnapshot{}, false, nil
}

// directivesOrDefault はテスト用の指示文集合を返す。未設定なら口調・固有名の最小集合を既定にする。
func (f *fakeStore) directivesOrDefault() []model.Directive {
	if f.directives != nil {
		return f.directives
	}
	return []model.Directive{
		{Key: "口調", Instruction: testPersonaTemplate},
		{Key: "固有名", Instruction: "これは固有名詞です。簡潔に訳すこと。"},
	}
}

func (f *fakeStore) ListDirectives(_ context.Context) ([]model.Directive, error) {
	return f.directivesOrDefault(), nil
}

func (f *fakeStore) GetDirectiveInstruction(_ context.Context, key string) (string, error) {
	for _, d := range f.directivesOrDefault() {
		if d.Key == key {
			return d.Instruction, nil
		}
	}
	return "", nil
}

func (f *fakeStore) ListRecordTypeMaster(_ context.Context) ([]model.RecordType, error) {
	return f.recordTypes, nil
}

func (f *fakeStore) ListExtractedFields(_ context.Context) ([]model.ExtractedField, error) {
	return f.extracted, nil
}

func (f *fakeStore) IngestNarrations(_ context.Context, rows []model.Narration) (int, error) {
	f.ingestedNarr = append(f.ingestedNarr, rows...)
	return len(rows), nil
}

func (f *fakeStore) IngestProperNouns(_ context.Context, rows []model.ProperNoun) (int, error) {
	f.ingestedPN = append(f.ingestedPN, rows...)
	return len(rows), nil
}

func (f *fakeStore) IngestLines(_ context.Context, rows []model.Line) (int, error) {
	f.ingestedLine = append(f.ingestedLine, rows...)
	return len(rows), nil
}

func (f *fakeStore) LinkLineSpeakersFromStaging(_ context.Context) error {
	f.linkCalled = true
	return nil
}

func (f *fakeStore) LinkLineConditionsFromStaging(_ context.Context) error {
	f.condLinkCalled = true
	return nil
}

func (f *fakeStore) LinkLineEmotionsFromStaging(_ context.Context) error {
	f.emoLinkCalled = true
	return nil
}

// --- MentionStore（言及段の全行取得・言及投入・説明対象解決） ---

func (f *fakeStore) ListNarrations(_ context.Context) ([]model.Narration, error) {
	return f.allNarrations, nil
}

func (f *fakeStore) ListLines(_ context.Context) ([]model.Line, error) {
	return f.allLines, nil
}

func (f *fakeStore) InsertNarrationMentions(_ context.Context, rows []model.NarrationMention) (int, error) {
	f.narrationMentions = append(f.narrationMentions, rows...)
	return len(rows), nil
}

func (f *fakeStore) InsertLineMentions(_ context.Context, rows []model.LineMention) (int, error) {
	f.lineMentions = append(f.lineMentions, rows...)
	return len(rows), nil
}

func (f *fakeStore) LinkNarrationDescribed(_ context.Context) (int, error) {
	f.describedCalled = true
	return 0, nil
}

func (f *fakeStore) ListProperNouns(_ context.Context) ([]model.ProperNoun, error) {
	return f.proper, nil
}

func (f *fakeStore) ListTranslatedProperNouns(_ context.Context, plugin string) ([]model.ProperNoun, error) {
	var out []model.ProperNoun
	for _, row := range f.proper {
		if row.Plugin == plugin && row.Status != 0 && row.Dest != "" {
			out = append(out, row)
		}
	}
	return out, nil
}

func (f *fakeStore) ListUntranslatedProperNouns(_ context.Context, plugin string) ([]model.ProperNoun, error) {
	var out []model.ProperNoun
	for _, p := range f.proper {
		if p.Status == 0 && (plugin == "" || p.Plugin == plugin) {
			out = append(out, p)
		}
	}
	return out, nil
}

// ListConfirmedNPCNames は実 DB の結合（proper_noun × extracted_field で field を取り戻す）を模す。
// 訳が確定した NPC_ の固有名を、原文の field（FULL / SHRT）つきで返す。並びは (field, source) で固定する。
func (f *fakeStore) ListConfirmedNPCNames(_ context.Context, plugin string) ([]model.ConfirmedName, error) {
	var out []model.ConfirmedName
	for _, pn := range f.proper {
		if pn.Dest == "" || pn.Category != "NPC_" || (plugin != "" && pn.Plugin != plugin) {
			continue
		}
		for _, ef := range f.extracted {
			if ef.Plugin != pn.Plugin || ef.Rec != pn.Category || ef.Source != pn.Source {
				continue
			}
			if ef.Field != "FULL" && ef.Field != "SHRT" {
				continue
			}
			out = append(out, model.ConfirmedName{Field: ef.Field, Source: pn.Source, Dest: pn.Dest})
		}
	}
	sort.Slice(out, func(i, j int) bool {
		if out[i].Field != out[j].Field {
			return out[i].Field < out[j].Field
		}
		return out[i].Source < out[j].Source
	})
	return out, nil
}

func (f *fakeStore) InsertDerivedProperNouns(_ context.Context, rows []model.ProperNoun) (int, error) {
	f.derivedPropers = append(f.derivedPropers, rows...)
	// 実 DB と同じく backing 行にも足す（派生した部分形が機械置換辞書へ合流する経路を再現する）。
	f.proper = append(f.proper, rows...)
	return len(rows), nil
}

func (f *fakeStore) UpdateProperNounDest(_ context.Context, id int64, dest string, status int) error {
	f.properUpdates = append(f.properUpdates, update{id: id, dest: dest, status: status})
	// 実 DB と同じく backing 行も更新する（未訳一覧の除外・proper_noun→辞書合流を batch 結合テストで再現する）。
	for i := range f.proper {
		if f.proper[i].ID == id {
			f.proper[i].Dest = dest
			f.proper[i].Status = status
		}
	}
	return nil
}

func (f *fakeStore) GetPromptTemplate(_ context.Context) (model.PromptTemplate, error) {
	if f.tmpl.BaseDirective == "" {
		return model.PromptTemplate{BaseDirective: testBaseDirective}, nil
	}
	return f.tmpl, nil
}

type update struct {
	id     int64
	dest   string
	status int
}

func (f *fakeStore) ListUntranslatedNarrations(_ context.Context, plugin string) ([]model.Narration, error) {
	var out []model.Narration
	for _, n := range f.untranslated {
		if n.Status == 0 && (plugin == "" || n.Plugin == plugin) {
			out = append(out, n)
		}
	}
	return out, nil
}

func (f *fakeStore) UpdateNarrationDest(_ context.Context, id int64, dest string, status int) error {
	f.updates = append(f.updates, update{id: id, dest: dest, status: status})
	// 実 DB と同じく backing 行も更新する（未訳一覧の除外を batch 結合テストで再現する）。
	for i := range f.untranslated {
		if f.untranslated[i].ID == id {
			f.untranslated[i].Dest = dest
			f.untranslated[i].Status = status
		}
	}
	return nil
}

func (f *fakeStore) ListUntranslatedLines(_ context.Context, plugin string) ([]model.Line, error) {
	var out []model.Line
	for _, l := range f.lines {
		if l.Status == 0 && (plugin == "" || l.Plugin == plugin) {
			out = append(out, l)
		}
	}
	return out, nil
}

func (f *fakeStore) UpdateLineDest(_ context.Context, id int64, dest string, status int) error {
	f.lineUpdates = append(f.lineUpdates, update{id: id, dest: dest, status: status})
	// 実 DB と同じく backing 行も更新する（未訳一覧の除外を batch 結合テストで再現する）。
	for i := range f.lines {
		if f.lines[i].ID == id {
			f.lines[i].Dest = dest
			f.lines[i].Status = status
		}
	}
	return nil
}

func (f *fakeStore) ListMasterTerms(_ context.Context) ([]model.MasterTerm, error) {
	return f.terms, nil
}

func (f *fakeStore) InsertDerivedTerms(_ context.Context, terms []model.MasterTerm) (int, error) {
	f.insertedTerms = append(f.insertedTerms, terms...)
	return len(terms), nil
}

func (f *fakeStore) ListReferenceTranslations(_ context.Context) ([]model.ReferenceTranslation, error) {
	return f.refs, nil
}

func (f *fakeStore) CountReferenceTranslations(_ context.Context) (int, error) {
	return len(f.refs), nil
}

// --- PersonaStore（生成入力・キャッシュ・保存・注入） ---

func (f *fakeStore) ListSpeakerLineSources(_ context.Context) ([]model.SpeakerLineSource, error) {
	f.generateInputCalls++
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

func (f *fakeStore) LoadLineConditions(_ context.Context, lineIDs []int64) (map[int64]string, error) {
	out := make(map[int64]string)
	for _, id := range lineIDs {
		if sex, ok := f.lineConditions[id]; ok {
			out[id] = sex
		}
	}
	return out, nil
}

// ExportStore（書き出し用読み出し）は engine 側の翻訳テストでは使わないため、空実装で interface を満たす。
// 書き出し経路は DB 込みのため単体では検証せず実画面 E2E に委ねる（テスト設計のとおり）。
func (f *fakeStore) ExportPlugins(_ context.Context) ([]string, error) { return nil, nil }

func (f *fakeStore) NarrationsForExport(_ context.Context, _ string) ([]model.Narration, error) {
	return nil, nil
}

func (f *fakeStore) LinesForExport(_ context.Context, _ string) ([]model.Line, error) {
	return nil, nil
}

func (f *fakeStore) ProperNounPlacementsForExport(_ context.Context, _ string) ([]model.ProperNounPlacement, error) {
	return nil, nil
}

type fakeTranslator struct {
	out        map[string]string // user メッセージ（機械置換済み原文）→ 訳文
	gotModel   string
	gotPrompts map[string]provider.Prompt // user メッセージ → engine が組んで送った完成プロンプト
	err        error                      // 全行共通で返すエラー（非 nil なら毎回失敗）
	errByUser  map[string]error           // user メッセージ → その行だけ返すエラー（行ごとの失敗注入用）
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
	if e, ok := f.errByUser[prompt.User]; ok {
		return "", e
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
	eng := New(store, tr, fakeLexicon{}, nil, nil)

	count, err := eng.Run(context.Background(), provider.Connection{Endpoint: "http://x"}, "model-x", "", nil)
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
		"The Riften guard waited.": "リフテンの衛兵が待っていた。",
	}}
	eng := New(store, tr, fakeLexicon{}, nil, nil)

	if _, err := eng.Run(context.Background(), provider.Connection{}, "m", "", nil); err != nil {
		t.Fatalf("Run: %v", err)
	}
	// AI へ渡した user メッセージは英語原文のままで、参考語だけをsystemへ追加すること。
	if got, ok := tr.gotPrompts["The Riften guard waited."]; !ok || !strings.Contains(got.System, "Riften -> リフテン") {
		t.Errorf("英語原文または参考語が AI へ渡っていない。gotPrompts=%v", tr.gotPrompts)
	}
	// 書き戻した訳文は置換済み原文に対する訳であること。
	want := []update{{1, "リフテンの衛兵が待っていた。", 3}}
	if len(store.updates) != 1 || store.updates[0] != want[0] {
		t.Errorf("updates = %v, want %v", store.updates, want)
	}
}

// readerを注入した本番経路ではmaster_termを本文参考語へ流さず、対象pluginの訳済み固有名だけを併記する。
func TestBodyReferencesUsesPrebuiltAndTargetPluginProperNouns(t *testing.T) {
	store := &fakeStore{
		terms: []model.MasterTerm{{Source: "Legacy", Dest: "旧辞書訳"}},
		proper: []model.ProperNoun{
			{Plugin: "A.esp", Source: "Inigo", Dest: "イニゴ", Status: statusTranslated},
			{Plugin: "B.esp", Source: "Lucien", Dest: "ルシエン", Status: statusTranslated},
		},
	}
	reader := fakePrebuiltDictionary{references: []model.PrebuiltDictionaryReference{{Source: "Riften", Dest: "リフテン", PartOfSpeech: "noun", Meaning: "城塞を守る都市", SkyrimCategory: "city"}}}
	refs, err := New(store, &fakeTranslator{}, fakeLexicon{}, nil, nil, reader).bodyReferences(context.Background(), "A.esp")
	if err != nil {
		t.Fatalf("bodyReferences: %v", err)
	}
	if len(refs) != 2 {
		t.Fatalf("参考語 = %+v", refs)
	}
	joined := fmt.Sprint(refs)
	for _, want := range []string{"Riften", "リフテン", "Inigo", "イニゴ"} {
		if !strings.Contains(joined, want) {
			t.Errorf("参考語に%qが無い: %s", want, joined)
		}
	}
	for _, forbidden := range []string{"Legacy", "Lucien", "城塞を守る都市"} {
		if strings.Contains(joined, forbidden) {
			t.Errorf("本文参考語に%qが混入した: %s", forbidden, joined)
		}
	}
}

// reader検証が失敗した場合は同期翻訳で固有名・本文のprovider送信を開始しない。
func TestRunStopsBeforeProviderWhenPrebuiltValidationFails(t *testing.T) {
	store := &fakeStore{untranslated: []model.Narration{{ID: 1, Source: "See Riften."}}, proper: []model.ProperNoun{{ID: 2, Source: "Riften"}}}
	translator := &fakeTranslator{out: map[string]string{"Riften": "リフテン", "See Riften.": "リフテンを見よ。"}}
	eng := New(store, translator, fakeLexicon{}, nil, nil, fakePrebuiltDictionary{err: errors.New("dictionary unavailable")})
	if _, err := eng.Run(context.Background(), provider.Connection{}, "m", "A.esp", nil); err == nil {
		t.Fatal("reader検証失敗を返さなかった")
	}
	if len(translator.gotPrompts) != 0 || len(store.updates) != 0 || len(store.properUpdates) != 0 {
		t.Fatalf("reader失敗後に送信または保存を開始した: prompts=%v narr=%v proper=%v", translator.gotPrompts, store.updates, store.properUpdates)
	}
}

// 原文が既存訳（参照訳）と完全一致する叙述文は、AI を呼ばず既訳を確定訳（status=1）で流用すること（known-issues 項目7）。
func TestRunReusesExistingTranslationWithoutCallingAI(t *testing.T) {
	store := &fakeStore{
		untranslated: []model.Narration{{ID: 1, Rec: "WEAP", Field: "DESC", Source: "A fine blade."}},
		refs:         []model.ReferenceTranslation{{Rec: "WEAP", Field: "DESC", Source: "A fine blade.", Dest: "見事な刃。"}},
	}
	// provider は訳を持たない。既訳流用で AI が呼ばれないことを、呼ばれたら訳が空になることで検出する。
	tr := &fakeTranslator{out: map[string]string{}}
	eng := New(store, tr, fakeLexicon{}, nil, nil)

	count, err := eng.Run(context.Background(), provider.Connection{}, "m", "", nil)
	if err != nil {
		t.Fatalf("Run: %v", err)
	}
	if count != 1 {
		t.Errorf("count = %d, want 1（流用も翻訳件数に数える）", count)
	}
	// 既訳を確定訳 status=1 で書き戻すこと（仮訳 status=3 でない）。
	want := []update{{1, "見事な刃。", 1}}
	if len(store.updates) != 1 || store.updates[0] != want[0] {
		t.Fatalf("updates = %v, want %v", store.updates, want)
	}
	// この原文は AI へ渡っていないこと（provider を呼ばずに流用した）。
	if _, called := tr.gotPrompts["A fine blade."]; called {
		t.Errorf("既訳一致の叙述文が AI へ渡った（流用されていない）。gotPrompts=%v", tr.gotPrompts)
	}
}

// 実行時タグが訳文から失われた行は、壊れた訳を確定させず未訳（status 更新なし）のまま残すこと（known-issues 項目8）。
// モデルが退避プレースホルダを落とした場合の扱い。翻訳件数にも数えない（再実行で再翻訳させる）。
func TestRunLeavesRowUntranslatedWhenRuntimeTagLost(t *testing.T) {
	store := &fakeStore{untranslated: []model.Narration{{ID: 1, Rec: "BOOK", Field: "DESC", Source: "See <Alias=Player> now."}}}
	// モデルが生タグ <Alias=Player> を落とした訳を返す（プロンプトの user は生タグ "See <Alias=Player> now."）。
	tr := &fakeTranslator{out: map[string]string{"See <Alias=Player> now.": "今すぐ見よ。"}}
	eng := New(store, tr, fakeLexicon{}, nil, nil)

	count, err := eng.Run(context.Background(), provider.Connection{}, "m", "", nil)
	if err != nil {
		t.Fatalf("Run: %v", err)
	}
	// 壊れた訳を書き戻していないこと（未訳のまま）。
	if len(store.updates) != 0 {
		t.Fatalf("タグ欠落行を書き戻した（未訳で残すべき）: %v", store.updates)
	}
	// 翻訳件数に数えていないこと。
	if count != 0 {
		t.Errorf("count = %d, want 0（タグ欠落行は数えない）", count)
	}
}

// 構造化パース失敗（provider.ErrStructuredParse）の行は未訳のまま skip し、翻訳全体を止めず他の行の翻訳を続けること。
// 実 LLM（7B）の非決定的な空応答で 1 行が失敗しても Run が停止しない回帰の再発防止。タグ欠落と同じ扱い。
func TestRunSkipsStructuredParseFailure(t *testing.T) {
	store := &fakeStore{untranslated: []model.Narration{
		{ID: 1, Source: "halls"},
		{ID: 2, Source: "cairn"},
		{ID: 3, Source: "keep"},
	}}
	tr := &fakeTranslator{
		out:       map[string]string{"halls": "広間", "keep": "砦"},
		errByUser: map[string]error{"cairn": provider.ErrStructuredParse},
	}
	eng := New(store, tr, fakeLexicon{}, nil, nil)

	count, err := eng.Run(context.Background(), provider.Connection{Endpoint: "http://x"}, "model-x", "", nil)
	if err != nil {
		t.Fatalf("Run が構造化パース失敗で停止した（skip すべき）: %v", err)
	}
	// 失敗した cairn(ID=2) は書き戻さず、成功した halls(ID=1)・keep(ID=3) だけを仮訳で書き戻すこと。
	want := []update{{1, "広間", 3}, {3, "砦", 3}}
	if len(store.updates) != len(want) {
		t.Fatalf("updates = %v, want %v", store.updates, want)
	}
	for i, u := range want {
		if store.updates[i] != u {
			t.Errorf("updates[%d] = %v, want %v", i, store.updates[i], u)
		}
	}
	// 失敗行は進捗に数えないこと（成功 2 件）。
	if count != 2 {
		t.Errorf("count = %d, want 2（失敗行は数えない）", count)
	}
}

// provider がエラーを返したら、そのエラーを返すこと（壊れた訳を書き戻さない）。
func TestRunReturnsProviderError(t *testing.T) {
	store := &fakeStore{untranslated: []model.Narration{{ID: 1, Source: "x"}}}
	tr := &fakeTranslator{err: errors.New("connection refused")}
	eng := New(store, tr, fakeLexicon{}, nil, nil)

	_, err := eng.Run(context.Background(), provider.Connection{}, "m", "", nil)
	if err == nil {
		t.Fatal("want error, got nil")
	}
	if len(store.updates) != 0 {
		t.Errorf("should not write dest on provider error, got %v", store.updates)
	}
}

// provider の skippable な失敗（サーバ一時失敗 429/5xx・応答エンベロープの読み取り失敗）の行は未訳のまま skip し、
// run 全体を止めず他の行の翻訳を続けること（known-issues 項目7。構造化パース失敗と同じ扱いへ広げた回帰の固定）。
func TestRunSkipsSkippableProviderFailures(t *testing.T) {
	store := &fakeStore{untranslated: []model.Narration{
		{ID: 1, Source: "halls"},
		{ID: 2, Source: "cairn"},
		{ID: 3, Source: "keep"},
		{ID: 4, Source: "gate"},
	}}
	tr := &fakeTranslator{
		out: map[string]string{"halls": "広間", "keep": "砦"},
		errByUser: map[string]error{
			"cairn": fmt.Errorf("%w: status 503", provider.ErrServerTransient),
			"gate":  fmt.Errorf("%w: choices が無い", provider.ErrResponseUnreadable),
		},
	}
	eng := New(store, tr, fakeLexicon{}, nil, nil)

	count, err := eng.Run(context.Background(), provider.Connection{Endpoint: "http://x"}, "model-x", "", nil)
	if err != nil {
		t.Fatalf("Run が skippable 失敗で停止した（skip すべき）: %v", err)
	}
	// 失敗した cairn(ID=2)・gate(ID=4) は書き戻さず、成功した halls(ID=1)・keep(ID=3) だけを仮訳で書き戻すこと。
	want := []update{{1, "広間", 3}, {3, "砦", 3}}
	if len(store.updates) != len(want) {
		t.Fatalf("updates = %v, want %v", store.updates, want)
	}
	for i, u := range want {
		if store.updates[i] != u {
			t.Errorf("updates[%d] = %v, want %v", i, store.updates[i], u)
		}
	}
	if count != 2 {
		t.Errorf("count = %d, want 2（skippable 失敗行は数えない）", count)
	}
}

// provider の fatal な失敗（skippable 番兵でラップされない失敗。通信断・認証や不正リクエストの 4xx・未知の失敗）は
// run を止め、失敗行より後の行を処理しないこと（既定で止める安全側の挙動）。
func TestRunStopsOnFatalProviderFailure(t *testing.T) {
	store := &fakeStore{untranslated: []model.Narration{
		{ID: 1, Source: "halls"},
		{ID: 2, Source: "cairn"},
		{ID: 3, Source: "keep"},
	}}
	tr := &fakeTranslator{
		out:       map[string]string{"halls": "広間", "keep": "砦"},
		errByUser: map[string]error{"cairn": errors.New("翻訳要求: status 401")},
	}
	eng := New(store, tr, fakeLexicon{}, nil, nil)

	_, err := eng.Run(context.Background(), provider.Connection{Endpoint: "http://x"}, "model-x", "", nil)
	if err == nil {
		t.Fatal("fatal 失敗で run が止まらなかった")
	}
	// fatal 行 cairn(ID=2) より前の halls(ID=1) だけ書き戻し、以降の keep(ID=3) は処理しないこと。
	want := []update{{1, "広間", 3}}
	if len(store.updates) != 1 || store.updates[0] != want[0] {
		t.Errorf("updates = %v, want %v（fatal 以降は処理しない）", store.updates, want)
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
	eng := New(store, tr, fakeLexicon{}, nil, nil)

	count, err := eng.Run(context.Background(), provider.Connection{}, "m", "", nil)
	if err != nil {
		t.Fatalf("Run: %v", err)
	}
	if count != 1 {
		t.Errorf("count = %d, want 1", count)
	}
	// 口調指示は完成プロンプトの system へ合成される（base 指示の後ろに続く）。
	system := tr.gotPrompts["mother?"].System
	if !strings.Contains(system, "穏やかに話す") {
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
	eng := New(store, tr, fakeLexicon{}, nil, nil)

	if _, err := eng.Run(context.Background(), provider.Connection{}, "m", "", nil); err != nil {
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
	eng := New(store, tr, fakeLexicon{}, nil, nil)

	var seen [][2]int
	_, err := eng.Run(context.Background(), provider.Connection{}, "m", "", func(done, total int) {
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

// R-2-1: 保存済みの準備結果を使う翻訳は口調を再集計せず、未訳だけを進捗総数にして処理すること。
func TestTranslateUntranslatedUsesOnlyPendingRowsWithoutPersonaRegeneration(t *testing.T) {
	store := &fakeStore{
		untranslated: []model.Narration{{ID: 1, Plugin: "A.esp", Source: "pending"}},
	}
	eng := New(store, &fakeTranslator{out: map[string]string{"pending": "未訳の訳"}}, fakeLexicon{}, nil, nil)

	var seen [][2]int
	count, err := eng.TranslateUntranslated(context.Background(), provider.Connection{}, "m", "A.esp", func(done, total int) {
		seen = append(seen, [2]int{done, total})
	})
	if err != nil {
		t.Fatalf("TranslateUntranslated: %v", err)
	}
	if count != 1 || len(store.updates) != 1 || store.updates[0].id != 1 {
		t.Fatalf("未訳の処理結果 = count:%d updates:%v", count, store.updates)
	}
	if store.generateInputCalls != 0 {
		t.Errorf("再実行で口調集計を開始した: %d 回", store.generateInputCalls)
	}
	want := [][2]int{{0, 1}, {1, 1}}
	if fmt.Sprint(seen) != fmt.Sprint(want) {
		t.Errorf("進捗 = %v, want %v", seen, want)
	}
}

// R-2-2: 全本文が未訳の場合も、その全件だけを進捗総数として処理すること。
func TestTranslateUntranslatedReportsAllPendingRows(t *testing.T) {
	store := &fakeStore{
		untranslated: []model.Narration{{ID: 1, Plugin: "A.esp", Source: "a"}},
		lines:        []model.Line{{ID: 2, Plugin: "A.esp", Source: "b"}},
	}
	eng := New(store, &fakeTranslator{out: map[string]string{"a": "あ", "b": "び"}}, fakeLexicon{}, nil, nil)

	var totals []int
	count, err := eng.TranslateUntranslated(context.Background(), provider.Connection{}, "m", "A.esp", func(_, total int) {
		totals = append(totals, total)
	})
	if err != nil {
		t.Fatalf("TranslateUntranslated: %v", err)
	}
	if count != 2 || len(totals) != 3 {
		t.Fatalf("全未訳の処理結果 = count:%d totals:%v", count, totals)
	}
	for _, total := range totals {
		if total != 2 {
			t.Errorf("進捗総数 = %d, want 2", total)
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
	eng := New(store, &fakeTranslator{}, fakeLexicon{}, nil, nil)

	// 10 はペルソナあり（経路①）。99 はペルソナなしで本文も rec/field も無いため口調なし。
	personas, err := eng.LinePersonas(context.Background(),
		[]model.Line{{ID: 10}, {ID: 99}}, testPersonaTemplate, ToneDefaults{})
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

// LinePersonas は話者を解決できない台詞を (rec, field) で汎用・PC へ振り分け、自由記述の口調へ
// 本文 1 行の感情と性別の一人称・語尾を重ねること。名指し話者は持たない経路②③の組み立てを確かめる。
func TestLinePersonasGenericAndPC(t *testing.T) {
	store := &fakeStore{
		// 40 は INFO:RNAM だが NPC 話者が結ばれている（line_speaker が INFO を解決するため起こる）。
		// それでも PC 発話として扱い、PC 既定の口調を優先する（名指しの口調を付けない）。
		linePersonas:   map[int64]model.LinePersonaInput{40: {AttitudeBand: 2, EmotionBand: 1, RaceEDID: "NordRace"}},
		lineConditions: map[int64]string{10: "Female"}, // 汎用台詞 10 は条件由来で女性
	}
	roles, err := rolespeech.ParseRoleSpeech(strings.NewReader("adult\tfemale\t*\tわたし\t女性らしく。\n"))
	if err != nil {
		t.Fatalf("ParseRoleSpeech: %v", err)
	}
	roles, err = rolespeech.ParseRoleSpeechExamples(roles, strings.NewReader(
		"adult\tdefault\tfemale\t*\tF1\t女性F1\n"+
			"adult\tdefault\tfemale\t*\tF2\t女性F2\n"+
			"adult\tdefault\tfemale\t*\tF3\t女性F3\n"))
	if err != nil {
		t.Fatalf("ParseRoleSpeechExamples: %v", err)
	}
	eng := New(store, &fakeTranslator{}, fakeLexicon{}, roles, nil)
	defaults := ToneDefaults{Generic: "衛兵の汎用台詞。", PC: "PCの選択肢。", PcSex: "Male"}
	lines := []model.Line{
		{ID: 10, Source: "Halt.", Rec: "INFO", Field: "NAM1"},  // 汎用（話者なし）
		{ID: 20, Source: "Yes.", Rec: "DIAL", Field: "FULL"},   // PC（選択肢の既定文）
		{ID: 30, Source: "Maybe.", Rec: "INFO", Field: "RNAM"}, // PC（選択肢の条件別上書き、話者なし）
		{ID: 40, Source: "Sure.", Rec: "INFO", Field: "RNAM"},  // PC（条件別上書き、NPC 話者が結ばれている）
	}
	personas, err := eng.LinePersonas(context.Background(), lines, testPersonaTemplate, defaults)
	if err != nil {
		t.Fatalf("LinePersonas: %v", err)
	}

	// 汎用: 決定経路=汎用、条件由来の女性、自由記述口調＋女性の一人称。
	g, ok := personas[10]
	if !ok {
		t.Fatal("汎用台詞 10 が persona map に無い")
	}
	if g.DecisionPath != tone.PathGeneric || g.Sex != "Female" || g.Label != "口調: 汎用台詞" {
		t.Errorf("汎用の口調メタが想定外: path=%q sex=%q label=%q", g.DecisionPath, g.Sex, g.Label)
	}
	if !strings.Contains(g.Directive, "衛兵の汎用台詞") || !strings.Contains(g.Directive, "わたし") {
		t.Errorf("汎用の directive が想定外: %q", g.Directive)
	}
	if strings.Count(g.Directive, "- 例:") != 3 || !strings.Contains(g.Directive, "女性F3") {
		t.Errorf("汎用の性別別3例が想定外: %q", g.Directive)
	}

	// PC（DIAL:FULL）: 決定経路=PC、利用者選択の男性。男性は役割語テンプレートに当たらず一人称なし。
	p, ok := personas[20]
	if !ok {
		t.Fatal("PC 発話 20 が persona map に無い")
	}
	if p.DecisionPath != tone.PathPC || p.Sex != "Male" || p.Label != "口調: PC発話" {
		t.Errorf("PC の口調メタが想定外: path=%q sex=%q label=%q", p.DecisionPath, p.Sex, p.Label)
	}
	if !strings.Contains(p.Directive, "PCの選択肢") {
		t.Errorf("PC の directive が想定外: %q", p.Directive)
	}
	if strings.Contains(p.Directive, "- 例:") {
		t.Errorf("PC の directive に例文が入った: %q", p.Directive)
	}

	// PC（INFO:RNAM・話者なし）も PC 経路へ振り分けること。
	if personas[30].DecisionPath != tone.PathPC {
		t.Errorf("INFO:RNAM（話者なし）の決定経路 = %q, want PC", personas[30].DecisionPath)
	}

	// PC（INFO:RNAM）に NPC 話者が結ばれていても、名指しでなく PC 既定の口調を優先すること。
	r, ok := personas[40]
	if !ok {
		t.Fatal("INFO:RNAM（話者あり）40 が persona map に無い")
	}
	if r.DecisionPath != tone.PathPC || r.Cell != "" {
		t.Errorf("INFO:RNAM（話者あり）の決定経路 = %q cell = %q, want PC・セルなし", r.DecisionPath, r.Cell)
	}
	if !strings.Contains(r.Directive, "PCの選択肢") {
		t.Errorf("INFO:RNAM（話者あり）の directive が PC 既定でない: %q", r.Directive)
	}
}

// 注入の引きの DB 呼び出しが台詞数 N に依存せず一括（定数 1 回）であること（N+1 廃止の観測）。
func TestLinePersonasBulkLoadsOnce(t *testing.T) {
	for _, n := range []int{2, 50} {
		store := &fakeStore{linePersonas: map[int64]model.LinePersonaInput{}}
		lines := make([]model.Line, n)
		for i := range n {
			id := int64(i + 1)
			lines[i] = model.Line{ID: id}
			store.linePersonas[id] = model.LinePersonaInput{AttitudeBand: 1, EmotionBand: 1}
		}
		eng := New(store, &fakeTranslator{}, fakeLexicon{}, nil, nil)
		if _, err := eng.LinePersonas(context.Background(), lines, testPersonaTemplate, ToneDefaults{}); err != nil {
			t.Fatalf("LinePersonas(N=%d): %v", n, err)
		}
		if store.loadPersonasCalls != 1 {
			t.Errorf("N=%d: LoadLinePersonas 呼び出し = %d, want 1（台詞数に非依存）", n, store.loadPersonasCalls)
		}
	}
}

// Run に対象 plugin を渡すと、その plugin の未訳台詞だけを翻訳し、別 plugin の未訳は触らないこと。
// 抽出した 1 plugin の実行が他 plugin の未訳を巻き込まない絞り込みを守る（実画面で発覚した既存挙動の修正）。
func TestRunScopesToTargetPlugin(t *testing.T) {
	store := &fakeStore{lines: []model.Line{
		{ID: 10, Source: "alpha", Plugin: "A.esp", FormID: "F1"},
		{ID: 20, Source: "bravo", Plugin: "B.esp", FormID: "F2"},
	}}
	tr := &fakeTranslator{out: map[string]string{"alpha": "ア", "bravo": "ブ"}}
	eng := New(store, tr, fakeLexicon{}, nil, nil)

	count, err := eng.Run(context.Background(), provider.Connection{}, "m", "A.esp", nil)
	if err != nil {
		t.Fatalf("Run: %v", err)
	}
	if count != 1 {
		t.Errorf("count = %d, want 1（A.esp の 1 行だけ）", count)
	}
	want := []update{{10, "ア", 3}}
	if len(store.lineUpdates) != 1 || store.lineUpdates[0] != want[0] {
		t.Errorf("lineUpdates = %v, want %v（B.esp は触らない）", store.lineUpdates, want)
	}
	if _, sent := tr.gotPrompts["bravo"]; sent {
		t.Errorf("B.esp の bravo を翻訳へ送った。対象 plugin 外は送らない")
	}
}
