// Package engine は翻訳手続きの本体を持つ。中心データを読み、ペルソナ生成を経て AI 翻訳し、
// 配置へ書き戻す純 Go の手続き。GUI から切り離して単体テストできる。
package engine

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"sort"

	"aitranslationenginejp/internal/core/dictionary"
	"aitranslationenginejp/internal/core/linefeatures"
	"aitranslationenginejp/internal/core/personatone"
	"aitranslationenginejp/internal/core/prompt"
	"aitranslationenginejp/internal/core/rolespeech"
	"aitranslationenginejp/internal/core/termxml"
	"aitranslationenginejp/internal/model"
	"aitranslationenginejp/internal/provider"
)

// statusProvisional は xTranslator の訳状態 3（仮）。AI 翻訳は仮訳として書き戻す。
const statusProvisional = 3

// NarrationStore は engine が叙述文の翻訳に使う中心データアクセス（使う分だけ宣言する）。
type NarrationStore interface {
	ListUntranslatedNarrations(ctx context.Context) ([]model.Narration, error)
	UpdateNarrationDest(ctx context.Context, id int64, dest string, status int) error
}

// LineStore は engine が台詞の翻訳に使う中心データアクセス（使う分だけ宣言する）。
type LineStore interface {
	ListUntranslatedLines(ctx context.Context) ([]model.Line, error)
	UpdateLineDest(ctx context.Context, id int64, dest string, status int) error
}

// PersonaStore は engine が口調ペルソナの一括生成と注入に使う中心データアクセス。
// 生成は台詞の言語特徴を line_analysis にキャッシュして集計し persona_character へ保存する。
// 注入は台詞 id 群へ生成済みの基底口調を引く。
type PersonaStore interface {
	ListSpeakerLineSources(ctx context.Context) ([]model.SpeakerLineSource, error)
	GetLineAnalyses(ctx context.Context, hashes []string) (map[string]model.LineAnalysis, error)
	UpsertLineAnalysis(ctx context.Context, a model.LineAnalysis) error
	UpsertPersonaCharacter(ctx context.Context, pc model.PersonaCharacter) error
	LoadLinePersonas(ctx context.Context, lineIDs []int64) (map[int64]model.LinePersonaInput, error)
}

// Persona は台詞 1 件へ与える口調指示文（全文）と、UI 表示用の口調メタデータ。話者を解決できた台詞だけ持つ。
// Directive・Label は翻訳と一覧チップ用。Cell 以降は結果行を展開したときに出す判定結果と根拠。
type Persona struct {
	Directive    string
	Label        string
	Cell         string // 基底口調セル名（判定結果）
	Trait        string // 基底口調の性質文（口調を普通の言葉で説明した一文）
	AttitudeBand int    // 対人段階 0尊大/1中立/2丁寧
	EmotionBand  int    // 感情段階 0抑制/1中/2激情
	Marked       int    // 印（信頼度の目安）
	DecisionPath string // 本文/voice/保留
}

// DictStore は engine が固有名の機械置換辞書（master_term）を読み書きするための中心データアクセス。
// ListMasterTerms は本文置換用の全件読み出し、InsertDerivedTerms は派生（人名の部分形）の追記。
type DictStore interface {
	ListMasterTerms(ctx context.Context) ([]model.MasterTerm, error)
	InsertDerivedTerms(ctx context.Context, terms []model.MasterTerm) (int, error)
}

// TemplateStore は engine がプロンプトテンプレート（base 指示）を読むための中心データアクセス。
type TemplateStore interface {
	GetPromptTemplate(ctx context.Context) (model.PromptTemplate, error)
}

// DirectiveStore は engine が指示文（directive）を読むための中心データアクセス。
// 本文フェーズは (rec, field) ごとの文体 directive と口調 directive の指示文を引いてプロンプトを組む。
type DirectiveStore interface {
	ListDirectives(ctx context.Context) ([]model.Directive, error)
	GetDirectiveInstruction(ctx context.Context, key string) (string, error)
}

// ProperNounDictStore は engine が本文機械置換辞書へ固有名（proper_noun）の訳を合流させるための読み出し。
type ProperNounDictStore interface {
	ListProperNouns(ctx context.Context) ([]model.ProperNoun, error)
}

// Store は engine が必要とする中心データアクセスをまとめる。concrete は internal/store が 1 つ実装する。
type Store interface {
	NarrationStore
	LineStore
	DictStore
	TemplateStore
	PersonaStore
	IngestStore
	ProperNounStore
	ProperNounDictStore
	DirectiveStore
}

// Engine は翻訳手続きを実行する。
type Engine struct {
	store      Store
	provider   provider.Translator
	lexicon    linefeatures.EmotionLexicon
	roleSpeech *rolespeech.Table
}

// New は engine を生成する。provider は AI 翻訳の port、lexicon は口調生成の感情辞書（差し替え可能な境界）、
// roleSpeech は注入時に引く一人称・語尾テンプレート（差し替え可能な参照データ。nil なら役割語を付けない）。
func New(store Store, p provider.Translator, lexicon linefeatures.EmotionLexicon, roleSpeech *rolespeech.Table) *Engine {
	return &Engine{store: store, provider: p, lexicon: lexicon, roleSpeech: roleSpeech}
}

// GeneratePersonas は台詞を持つ全話者の基底口調を一括生成し persona_character へ保存する。保存した話者数を返す。
// 最も重い行解析は line_analysis に本文ハッシュでキャッシュし、本文ごとに 1 度だけ prose を回す。
// 冪等で、手修正済み（hand_edited）の行は store が再生成から保護する。
func (e *Engine) GeneratePersonas(ctx context.Context) (int, error) {
	return NewPersonaGenerator(e.store, e.lexicon).Generate(ctx)
}

// Run は未訳の固有名・叙述文・定型句・台詞を順に翻訳し、訳文を書き戻す。翻訳できた合計件数を返す。
// 段階順序は固有名フェーズ（本文より先に固有名を確定）→ 本文フェーズ（叙述文・定型句・台詞）。
// プロンプトは Base 指示 ＋ その REC:FIELD に割り当てた directive の指示文（台詞は口調 directive の {traits} を埋める）で組む。
// onProgress が非 nil なら、固有名・本文を通した進捗（処理済み件数 done、総件数 total）を都度通知する。
func (e *Engine) Run(ctx context.Context, conn provider.Connection, model string, onProgress func(done, total int)) (int, error) { //nolint:gocognit // TODO(refactor): 翻訳手続きの段階連結（固有名→叙述文/定型句→台詞）。リファクタ本体で段階を関数へ分割する。
	propers, err := e.store.ListUntranslatedProperNouns(ctx)
	if err != nil {
		return 0, fmt.Errorf("未訳固有名の取得: %w", err)
	}
	narrations, err := e.store.ListUntranslatedNarrations(ctx)
	if err != nil {
		return 0, fmt.Errorf("未訳叙述文の取得: %w", err)
	}
	lines, err := e.store.ListUntranslatedLines(ctx)
	if err != nil {
		return 0, fmt.Errorf("未訳台詞の取得: %w", err)
	}

	// プロンプトテンプレート（base 指示）と指示文（directive）をループ前に 1 度だけ読む。
	tmpl, err := e.store.GetPromptTemplate(ctx)
	if err != nil {
		return 0, fmt.Errorf("プロンプトテンプレートの取得: %w", err)
	}
	instructionByKey, keyByRF, err := e.directiveLookups(ctx)
	if err != nil {
		return 0, err
	}

	// 既訳辞書（master_term の source→dest）を固有名フェーズの供給源選別に使う。
	authoritative, err := e.authoritativeTerms(ctx)
	if err != nil {
		return 0, err
	}

	// 口調ペルソナを一括生成して persona_character を最新化する。注入はこの生成結果を引く。
	if _, err = e.GeneratePersonas(ctx); err != nil {
		return 0, err
	}
	// 台詞ごとの口調指示をループ前に 1 度だけ一括で組む。口調 directive の指示文（{traits} を含む）へ性質を差し込む。
	personas, err := e.LinePersonas(ctx, lineIDsOf(lines), instructionByKey[directiveTone])
	if err != nil {
		return 0, err
	}

	total := len(propers) + len(narrations) + len(lines)
	done := 0
	report := func() {
		if onProgress != nil {
			onProgress(done, total)
		}
	}
	report()

	// 固有名フェーズ: 本文より先に固有名を確定する。既訳ありは権威訳、既訳なしは固有名 directive で AI 訳。
	if err = e.translateProperNouns(ctx, conn, model, propers, authoritative,
		tmpl.BaseDirective, instructionByKey[directiveProperNoun], func() { done++; report() }); err != nil {
		return done, err
	}

	// 本文フェーズの機械置換辞書は固有名フェーズ後に組む（master_term ∪ 確定した proper_noun）。
	dict, err := e.LoadDictionary(ctx)
	if err != nil {
		return done, err
	}

	for _, row := range narrations {
		// 本文中の固有名を辞書の確定訳語へ機械置換してから AI 翻訳する。AI は周りの英語だけを訳す。
		source, _ := dict.Apply(row.Source)
		// 叙述文・定型句は、その REC:FIELD に割り当てた文体・定型句 directive の指示文を base へ合成して訳す。
		instruction := instructionByKey[keyByRF[RecordKey{Rec: row.Rec, Field: row.Field}]]
		dest, err := e.provider.Translate(ctx, conn, model, prompt.ComposePrompt(tmpl.BaseDirective, instruction, source))
		if err != nil {
			return done, fmt.Errorf("叙述文の翻訳: %w", err)
		}
		if err := e.store.UpdateNarrationDest(ctx, row.ID, dest, statusProvisional); err != nil {
			return done, fmt.Errorf("叙述文の書き戻し: %w", err)
		}
		done++
		report()
	}

	for _, row := range lines {
		source, _ := dict.Apply(row.Source)
		// 台詞は口調 directive（{traits} を話者の性質で埋めたもの）を base へ合成して訳す。話者なしなら口調指示は空。
		dest, err := e.provider.Translate(ctx, conn, model, prompt.ComposePrompt(tmpl.BaseDirective, personas[row.ID].Directive, source))
		if err != nil {
			return done, fmt.Errorf("台詞の翻訳: %w", err)
		}
		if err := e.store.UpdateLineDest(ctx, row.ID, dest, statusProvisional); err != nil {
			return done, fmt.Errorf("台詞の書き戻し: %w", err)
		}
		done++
		report()
	}
	return done, nil
}

// directive のキー定数。本文フェーズが口調・固有名 directive を引くために使う（record_type_master.directive と一致）。
const (
	directiveTone       = "口調"
	directiveProperNoun = "固有名"
)

// directiveLookups は directive と record_type_master を読み、本文フェーズの 2 つの引きを返す。
// instructionByKey は directive キー → 指示文、keyByRF は (rec, field) → directive キー。
func (e *Engine) directiveLookups(ctx context.Context) (instructionByKey map[string]string, keyByRF map[RecordKey]string, err error) {
	directives, err := e.store.ListDirectives(ctx)
	if err != nil {
		return nil, nil, fmt.Errorf("directive の取得: %w", err)
	}
	instructionByKey = make(map[string]string, len(directives))
	for _, d := range directives {
		instructionByKey[d.Key] = d.Instruction
	}
	masterRows, err := e.store.ListRecordTypeMaster(ctx)
	if err != nil {
		return nil, nil, fmt.Errorf("record_type_master の取得: %w", err)
	}
	keyByRF = make(map[RecordKey]string, len(masterRows))
	for _, r := range masterRows {
		keyByRF[RecordKey{Rec: r.Rec, Field: r.Field}] = r.Directive
	}
	return instructionByKey, keyByRF, nil
}

// authoritativeTerms は master_term を source→dest の既訳辞書へ畳む（固有名フェーズの供給源選別に使う）。
func (e *Engine) authoritativeTerms(ctx context.Context) (map[string]string, error) {
	terms, err := e.store.ListMasterTerms(ctx)
	if err != nil {
		return nil, fmt.Errorf("master_term の取得: %w", err)
	}
	m := make(map[string]string, len(terms))
	for _, t := range terms {
		if _, ok := m[t.Source]; !ok {
			m[t.Source] = t.Dest
		}
	}
	return m, nil
}

// LoadDictionary は master_term と確定済みの proper_noun を合流し、固有名の機械置換辞書を組む。
// 翻訳実行（Run の本文フェーズ）と結果取得（api の内訳・実プロンプト再構成）が、ページ単位で 1 度だけ組んで使う。
// master_term（権威訳）を先に積み、同綴りの proper_noun より優先する（NewDictionary は先勝ち）。
func (e *Engine) LoadDictionary(ctx context.Context) (*dictionary.Dictionary, error) {
	terms, err := e.store.ListMasterTerms(ctx)
	if err != nil {
		return nil, fmt.Errorf("マスター辞書の取得: %w", err)
	}
	propers, err := e.store.ListProperNouns(ctx)
	if err != nil {
		return nil, fmt.Errorf("固有名の取得: %w", err)
	}
	pairs := make([]dictionary.Term, 0, len(terms)+len(propers))
	for _, term := range terms {
		pairs = append(pairs, dictionary.Term{Source: term.Source, Dest: term.Dest})
	}
	for _, pn := range propers {
		pairs = append(pairs, dictionary.Term{Source: pn.Source, Dest: pn.Dest})
	}
	return dictionary.NewDictionary(pairs), nil
}

// DeriveMasterTerms は xTranslator 英日 XML から人名の部分形（名のみ・短名）の確定訳語を派生し、
// master_term へ追記する。base 既出の原語との衝突を避け、派生行は category="derive:<種別>" で由来を残す。
// 翻訳 Run の抽出後に呼び、DB を作り直しても単独名（例 Aventus→アベンタス）が辞書へ入るようにする。
// 追記した件数を返す。INSERT OR IGNORE のため二重実行でも増えない（冪等）。
func (e *Engine) DeriveMasterTerms(ctx context.Context, xmlDir string) (int, error) {
	terms, err := e.store.ListMasterTerms(ctx)
	if err != nil {
		return 0, fmt.Errorf("base 辞書の取得: %w", err)
	}
	baseSources := make(map[string]bool, len(terms))
	for _, t := range terms {
		baseSources[t.Source] = true
	}
	files, err := readXMLDir(xmlDir)
	if err != nil {
		return 0, err
	}
	derived, err := termxml.DeriveTermsFromFiles(files, baseSources)
	if err != nil {
		return 0, fmt.Errorf("固有名の派生: %w", err)
	}
	rows := make([]model.MasterTerm, len(derived))
	for i, d := range derived {
		rows[i] = model.MasterTerm{Source: d.Source, Dest: d.Dest, Category: "derive:" + d.Kind}
	}
	return e.store.InsertDerivedTerms(ctx, rows)
}

// readXMLDir はディレクトリ内の全 XML を読み込み、ファイル名昇順で termxml.XMLFile 群にして返す。
// 純粋な派生（termxml.DeriveTermsFromFiles）へ渡す前段の os 読み。決定性のため名前順に並べる。
func readXMLDir(xmlDir string) ([]termxml.XMLFile, error) {
	entries, err := filepath.Glob(filepath.Join(xmlDir, "*.xml"))
	if err != nil {
		return nil, fmt.Errorf("XML の列挙: %w", err)
	}
	if len(entries) == 0 {
		return nil, fmt.Errorf("XML が無い: %s", xmlDir)
	}
	sort.Strings(entries)
	files := make([]termxml.XMLFile, 0, len(entries))
	for _, path := range entries {
		data, err := os.ReadFile(path)
		if err != nil {
			return nil, fmt.Errorf("%s の読み込み: %w", path, err)
		}
		files = append(files, termxml.XMLFile{Name: filepath.Base(path), Data: data})
	}
	return files, nil
}

// LinePersonas は台詞 id 群へ生成済みの基底口調を一括で引き、各台詞の口調指示文（全文）と一覧用の短い要約を map で返す。
// personaTemplate は口調指示の雛形（{traits} を含む）で、基底口調の性質文と種族訛りを差し込んで口調指示文を組む。
// 生成済みペルソナが無い台詞、および段階が範囲外で性質文が空になる台詞は map に現れない
// （呼び出し側は欠落を「口調なし」として扱う）。引きを 1 度の一括問い合わせにまとめ、台詞数に依存させない。
func (e *Engine) LinePersonas(ctx context.Context, lineIDs []int64, personaTemplate string) (map[int64]Persona, error) {
	inputs, err := e.store.LoadLinePersonas(ctx, lineIDs)
	if err != nil {
		return nil, fmt.Errorf("注入入力（生成ペルソナ）の一括取得: %w", err)
	}
	out := make(map[int64]Persona, len(inputs))
	for lineID, in := range inputs {
		directive := personatone.BuildToneDirective(personaTemplate, personatone.BuildToneTraits(in, e.roleSpeech))
		if directive == "" {
			continue // 段階が範囲外で性質文が組めない場合は口調なし扱い。
		}
		cell, trait := personatone.PersonaMetaOf(in)
		out[lineID] = Persona{
			Directive:    directive,
			Label:        personatone.BuildToneLabel(in),
			Cell:         cell,
			Trait:        trait,
			AttitudeBand: in.AttitudeBand,
			EmotionBand:  in.EmotionBand,
			Marked:       in.Marked,
			DecisionPath: in.DecisionPath,
		}
	}
	return out, nil
}

// lineIDsOf は台詞行から id 列を取り出す（一括取得の入力用）。
func lineIDsOf(lines []model.Line) []int64 {
	ids := make([]int64, len(lines))
	for i, row := range lines {
		ids[i] = row.ID
	}
	return ids
}
