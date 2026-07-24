// Package engine は翻訳手続きの本体を持つ。中心データを読み、ペルソナ生成を経て AI 翻訳し、
// 配置へ書き戻す純 Go の手続き。GUI から切り離して単体テストできる。
package engine

import (
	"context"
	"fmt"
	"log/slog"
	"sort"
	"strings"

	"aitranslationenginejp/internal/core/batchplan"
	"aitranslationenginejp/internal/core/dictionary"
	"aitranslationenginejp/internal/core/linefeatures"
	"aitranslationenginejp/internal/core/personatone"
	"aitranslationenginejp/internal/core/prompt"
	"aitranslationenginejp/internal/core/rolespeech"
	"aitranslationenginejp/internal/core/runtimetag"
	"aitranslationenginejp/internal/core/termderive"
	"aitranslationenginejp/internal/core/termusage"
	"aitranslationenginejp/internal/core/tone"
	"aitranslationenginejp/internal/model"
	"aitranslationenginejp/internal/provider"
)

// 台詞のレコード種別キー。注入の経路振り分けに使う（record_type_master の台詞 3 種に対応）。
const (
	recInfo   = "INFO"
	recDial   = "DIAL"
	fieldNam1 = "NAM1" // NPC の返答
	fieldRnam = "RNAM" // 選択肢の条件別上書き（PC 発話）
	fieldFull = "FULL" // 選択肢の既定文（PC 発話）
)

// statusProvisional は xTranslator の訳状態 3（仮）。AI 翻訳は仮訳として書き戻す。
const statusProvisional = 3

// NarrationStore は engine が叙述文の翻訳に使う中心データアクセス（使う分だけ宣言する）。
type NarrationStore interface {
	ListUntranslatedNarrations(ctx context.Context, plugin string) ([]model.Narration, error)
	UpdateNarrationDest(ctx context.Context, id int64, dest string, status int) error
}

// LineStore は engine が台詞の翻訳に使う中心データアクセス（使う分だけ宣言する）。
type LineStore interface {
	ListUntranslatedLines(ctx context.Context, plugin string) ([]model.Line, error)
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
	LoadLineConditions(ctx context.Context, lineIDs []int64) (map[int64]string, error)
}

// Persona は台詞 1 件へ与える口調指示文（全文）と、UI 表示用の口調メタデータ。口調が付いた台詞だけ持つ。
// Directive・Label は翻訳と一覧チップ用。Cell 以降は結果行を展開したときに出す判定結果と根拠。
// 名指し話者は Cell・AttitudeBand・Marked を持つ。汎用台詞・PC 発話はそれらを持たず、EmotionBand・Sex を持つ。
type Persona struct {
	Directive    string
	Label        string
	Cell         string // 基底口調セル名（判定結果）。名指し話者だけ。汎用・PC は空。
	Trait        string // 基底口調の性質文（口調を普通の言葉で説明した一文）。名指し話者だけ。
	AttitudeBand int    // 対人段階 0尊大/1中立/2丁寧。名指し話者だけ。
	EmotionBand  int    // 感情段階 0抑制/1中/2激情
	Marked       int    // 印（信頼度の目安）。名指し話者だけ。
	DecisionPath string // 本文/voice/保留/汎用/PC
	Sex          string // 性別（Female/Male/空）。汎用・PC の一人称・語尾の根拠。名指し話者は空（話者欄が出す）。
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
	MentionStore
	ProperNounStore
	ProperNounDictStore
	DirectiveStore
	ExportStore
	ReferenceStore
}

// Engine は翻訳手続きを実行する。
type Engine struct {
	store      Store
	provider   provider.Translator
	lexicon    linefeatures.EmotionLexicon
	roleSpeech *rolespeech.Table
	stoplist   *dictionary.Stoplist
	classifier *tone.Classifier
}

// New は engine を生成する。provider は AI 翻訳の port、lexicon は口調生成の感情辞書（差し替え可能な境界）、
// roleSpeech は注入時に引く一人称・語尾テンプレート（差し替え可能な参照データ。nil なら役割語を付けない）。
// stoplist は機械置換辞書・言及語彙の供給から一般語を除く選別規則（nil なら選別しない）。
// classifier は口調分類の不変ルール。汎用台詞・PC 発話の感情段階を本文 1 行から出すのに使う。
func New(store Store, p provider.Translator, lexicon linefeatures.EmotionLexicon, roleSpeech *rolespeech.Table, stoplist *dictionary.Stoplist) *Engine {
	return &Engine{store: store, provider: p, lexicon: lexicon, roleSpeech: roleSpeech, stoplist: stoplist, classifier: tone.NewClassifier()}
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
// plugin が空でなければその対象 plugin の未訳行だけを翻訳する（空なら全 plugin）。抽出した plugin の
// 実行が他 plugin の未訳を巻き込まないよう、固有名・叙述文・台詞の 3 段すべてを対象 plugin へ絞る。
func (e *Engine) Run(ctx context.Context, conn provider.Connection, model string, plugin string, onProgress func(done, total int)) (int, error) {
	propers, err := e.store.ListUntranslatedProperNouns(ctx, plugin)
	if err != nil {
		return 0, fmt.Errorf("未訳固有名の取得: %w", err)
	}
	narrations, err := e.store.ListUntranslatedNarrations(ctx, plugin)
	if err != nil {
		return 0, fmt.Errorf("未訳叙述文の取得: %w", err)
	}
	lines, err := e.store.ListUntranslatedLines(ctx, plugin)
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
	// 話者なし台詞（汎用・PC）は prompt_template 由来の自由記述口調へ感情と性別を重ねる。
	personas, err := e.LinePersonas(ctx, lines, instructionByKey[directiveTone], ToneDefaults{
		Generic: tmpl.GenericToneText,
		PC:      tmpl.PcToneText,
		PcSex:   tmpl.PcSex,
	})
	if err != nil {
		return 0, err
	}

	p := &runProgress{total: len(propers) + len(narrations) + len(lines), onProgress: onProgress}
	p.notify()

	// 固有名フェーズ: 本文より先に固有名を確定する。既訳ありは権威訳、既訳なしは固有名 directive で AI 訳。
	if err = e.translateProperNouns(ctx, conn, model, propers, authoritative,
		tmpl.BaseDirective, instructionByKey[directiveProperNoun], p.step); err != nil {
		return p.done, err
	}

	// 本文フェーズの機械置換辞書は固有名フェーズ後に組む（master_term ∪ 確定した proper_noun）。
	dict, err := e.LoadDictionary(ctx)
	if err != nil {
		return p.done, err
	}

	// 既存訳（参照訳）の照合表を組む。原文が既訳と完全一致する本文は AI を呼ばず既訳を流用する。
	refIndex, err := e.referenceIndex(ctx)
	if err != nil {
		return p.done, err
	}

	// 叙述文・定型句フェーズ。
	if err = e.translateNarrations(ctx, conn, model, narrations, dict, refIndex, tmpl, instructionByKey, keyByRF, p); err != nil {
		return p.done, err
	}
	// 台詞フェーズ。
	if err = e.translateLines(ctx, conn, model, lines, dict, refIndex, tmpl, personas, p); err != nil {
		return p.done, err
	}
	return p.done, nil
}

// runProgress は Run の進捗（処理済み done・総数 total）を集計し onProgress へ通知する。
// 段階メソッド（translateProperNouns・translateNarrations・translateLines）へ渡し、各段階が 1 件処理ごとに step を呼ぶ。
type runProgress struct {
	done, total int
	onProgress  func(done, total int)
}

// step は処理済み件数を 1 増やし、進捗を通知する。
func (p *runProgress) step() {
	p.done++
	p.notify()
}

// notify は onProgress が非 nil なら現在の進捗を通知する。
func (p *runProgress) notify() {
	if p.onProgress != nil {
		p.onProgress(p.done, p.total)
	}
}

// translateNarrations は叙述文・定型句を翻訳し、書き戻す。1 件ごとに進捗を 1 歩進める。
// 原文が既存訳（参照訳）と完全一致する行は AI を呼ばず既訳を確定訳（statusTranslated）で流用する。
// それ以外は実行時タグを機械置換の間だけ退避して生タグへ戻し、AI 翻訳して仮訳（statusProvisional）で書き戻す。
// AI 出力にタグが原形で残らなかった行は未訳のまま残す。各行は REC:FIELD の文体・定型句 directive を base へ合成して訳す。
func (e *Engine) translateNarrations(ctx context.Context, conn provider.Connection, modelName string, narrations []model.Narration, dict *dictionary.Dictionary, refIndex map[referenceKey]string, tmpl model.PromptTemplate, instructionByKey map[string]string, keyByRF map[RecordKey]string, p *runProgress) error {
	lostTags := 0
	var skips lineSkips
	for _, row := range narrations {
		// 原文が既存訳と完全一致するなら AI を呼ばず既訳を確定訳として流用する（known-issues 項目7）。
		if dest, ok := refIndex[referenceKey{Rec: row.Rec, Field: row.Field, Source: row.Source}]; ok {
			if err := e.store.UpdateNarrationDest(ctx, row.ID, dest, statusTranslated); err != nil {
				return fmt.Errorf("叙述文の既訳流用: %w", err)
			}
			p.step()
			continue
		}
		instruction := instructionByKey[keyByRF[RecordKey{Rec: row.Rec, Field: row.Field}]]
		bodyPrompt, tags := e.composeBodyPrompt(tmpl.BaseDirective, instruction, row.Source, dict)
		dest, err := e.provider.Translate(ctx, conn, modelName, bodyPrompt)
		// 結果適用の判断は batch 反映と共有する純粋規則（DecideApply）を通す。
		// 成功・タグ欠落・skippable 失敗は同期でも batch でも同じ扱い。非 skippable（Fatal）だけ同期は run を止める。
		missing := runtimetag.CountMissing(dest, tags)
		switch batchplan.DecideApply(dest, err, missing).Kind {
		case batchplan.ApplyConfirm:
			if writeErr := e.store.UpdateNarrationDest(ctx, row.ID, dest, statusProvisional); writeErr != nil {
				return fmt.Errorf("叙述文の書き戻し: %w", writeErr)
			}
			p.step()
		case batchplan.ApplySkipTagLost:
			lostTags += missing
		case batchplan.ApplySkipStructuredParse:
			skips.structuredParse++
		case batchplan.ApplySkipResponseUnreadable:
			skips.responseUnreadable++
		case batchplan.ApplySkipServerTransient:
			skips.serverTransient++
		case batchplan.ApplyFatal:
			return fmt.Errorf("叙述文の翻訳: %w", err)
		}
	}
	logLostRuntimeTags(ctx, "translateNarrations", lostTags)
	skips.log(ctx, "translateNarrations")
	return nil
}

// prepareSource は AI へ送る原文を組む。実行時タグを機械置換の間だけ退避し、置換後に生タグへ戻す。
// タグ内部の語を辞書機械置換が書き換えるのを防ぎつつ、AI には意味の分かる生タグを見せる。退避したタグ列も返す（欠落照合用）。
func (e *Engine) prepareSource(rawSource string, dict *dictionary.Dictionary) (source string, tags []string) {
	masked, tags := runtimetag.Mask(rawSource)
	replaced, _ := dict.Apply(masked)
	return runtimetag.Restore(replaced, tags), tags
}

// composeBodyPrompt は本文 1 件の送信プロンプトを組む（同期の本文フェーズと batch 送信が共有する文面構築）。
// 実行時タグを機械置換の間だけ退避して生タグへ戻し、base 指示＋directive（文体または口調）＋機械置換済み原文を合成する。
// 退避したタグ列も返す（AI 出力のタグ欠落照合に使う）。同期と batch が同じ文面を組むための 1 箇所。
func (e *Engine) composeBodyPrompt(base, directive, rawSource string, dict *dictionary.Dictionary) (provider.Prompt, []string) {
	source, tags := e.prepareSource(rawSource, dict)
	return prompt.ComposePrompt(base, directive, source), tags
}

// translateLines は台詞を翻訳し、書き戻す。1 件ごとに進捗を 1 歩進める。
// 原文が既存訳（参照訳）と完全一致する台詞は AI を呼ばず既訳を確定訳（statusTranslated）で流用する（叙述文と同じ）。
// それ以外は実行時タグを機械置換の間だけ退避して生タグへ戻し、口調指示つきで AI 翻訳して仮訳で書き戻す。
// AI 出力にタグが原形で残らなかった台詞は未訳のまま残す。話者なしなら口調指示は空。
func (e *Engine) translateLines(ctx context.Context, conn provider.Connection, modelName string, lines []model.Line, dict *dictionary.Dictionary, refIndex map[referenceKey]string, tmpl model.PromptTemplate, personas map[int64]Persona, p *runProgress) error {
	lostTags := 0
	var skips lineSkips
	for _, row := range lines {
		// 原文が既存訳と完全一致するなら AI を呼ばず既訳を確定訳として流用する（叙述文と同じ、known-issues 項目7）。
		if dest, ok := refIndex[referenceKey{Rec: row.Rec, Field: row.Field, Source: row.Source}]; ok {
			if err := e.store.UpdateLineDest(ctx, row.ID, dest, statusTranslated); err != nil {
				return fmt.Errorf("台詞の既訳流用: %w", err)
			}
			p.step()
			continue
		}
		bodyPrompt, tags := e.composeBodyPrompt(tmpl.BaseDirective, personas[row.ID].Directive, row.Source, dict)
		dest, err := e.provider.Translate(ctx, conn, modelName, bodyPrompt)
		// 結果適用の判断は batch 反映と共有する純粋規則（DecideApply）を通す（叙述文と同じ）。
		missing := runtimetag.CountMissing(dest, tags)
		switch batchplan.DecideApply(dest, err, missing).Kind {
		case batchplan.ApplyConfirm:
			if writeErr := e.store.UpdateLineDest(ctx, row.ID, dest, statusProvisional); writeErr != nil {
				return fmt.Errorf("台詞の書き戻し: %w", writeErr)
			}
			p.step()
		case batchplan.ApplySkipTagLost:
			lostTags += missing
		case batchplan.ApplySkipStructuredParse:
			skips.structuredParse++
		case batchplan.ApplySkipResponseUnreadable:
			skips.responseUnreadable++
		case batchplan.ApplySkipServerTransient:
			skips.serverTransient++
		case batchplan.ApplyFatal:
			return fmt.Errorf("台詞の翻訳: %w", err)
		}
	}
	logLostRuntimeTags(ctx, "translateLines", lostTags)
	skips.log(ctx, "translateLines")
	return nil
}

// logLostRuntimeTags は本文フェーズで失われた実行時タグの件数を、集約して 1 度だけ観測ログへ出す。
// where は失われた段（translateNarrations / translateLines）。lost が 0 なら何も出さない。
// 該当行は未訳のまま残す（再実行で再翻訳）ため、result は skipped とする。
// loop 内の 1 件ごとには出さず、フェーズ末で件数を集約する（観測ログ規約の大量処理集約に従う）。
func logLostRuntimeTags(ctx context.Context, where string, lost int) {
	if lost == 0 {
		return
	}
	slog.WarnContext(ctx, "実行時タグが失われた行を未訳として残した",
		slog.String("event", "runtime_tag_lost"),
		slog.String("where", "engine."+where),
		slog.String("result", "skipped"),
		slog.Int("count", lost),
	)
}

// lineSkips は本文フェーズで未訳のまま飛ばした行を、飛ばした理由別に数える。フェーズ末で集約して観測ログへ出す。
// 対象は provider の skippable な失敗（構造化出力の空・スキーマ違反、応答エンベロープの読み取り失敗、サーバ一時失敗）。
type lineSkips struct {
	structuredParse    int // 構造化出力の空・スキーマ違反（provider.ErrStructuredParse）
	responseUnreadable int // 応答エンベロープの読み取り失敗（provider.ErrResponseUnreadable）
	serverTransient    int // 非 200 の 429・5xx（provider.ErrServerTransient）
}

// log は未訳のまま飛ばした行の件数を、飛ばした理由別に集約して観測ログへ出す。
// where は飛ばした段（translateNarrations / translateLines）。件数 0 の理由は出さない。
// 該当行は未訳のまま残す（再実行で再翻訳）ため result は skipped とする。
// loop 内の 1 件ごとには出さず、フェーズ末で件数を集約する（観測ログ規約の大量処理集約に従う）。
func (s lineSkips) log(ctx context.Context, where string) {
	logSkippedLines(ctx, where, "structured_parse_failed", "構造化出力の解析に失敗した行を未訳として残した", s.structuredParse)
	logSkippedLines(ctx, where, "response_unreadable", "応答エンベロープを読めず未訳として残した行", s.responseUnreadable)
	logSkippedLines(ctx, where, "server_transient", "サーバ一時失敗（429・5xx）で未訳として残した行", s.serverTransient)
}

// logSkippedLines は 1 理由分の飛ばした行数を観測ログへ出す。count が 0 なら何も出さない。
func logSkippedLines(ctx context.Context, where, event, msg string, count int) {
	if count == 0 {
		return
	}
	slog.WarnContext(ctx, msg,
		slog.String("event", event),
		slog.String("where", "engine."+where),
		slog.String("result", "skipped"),
		slog.Int("count", count),
	)
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
// 供給源は言及検出（mentionDetector）と共通の translationVocabulary で読み、一般語 stoplist の選別を通す。
func (e *Engine) LoadDictionary(ctx context.Context) (*dictionary.Dictionary, error) {
	terms, propers, err := e.translationVocabulary(ctx)
	if err != nil {
		return nil, err
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

// translationVocabulary は機械置換辞書（LoadDictionary）と言及検出（mentionDetector）が共有する語彙の供給源。
// master_term → proper_noun の順で読み、一般語 stoplist の選別をこの 1 箇所で通す。
// 言及レコードは注入の忠実な記録という不変条件を保つため、片側だけに効く除外を作らない。
func (e *Engine) translationVocabulary(ctx context.Context) ([]model.MasterTerm, []model.ProperNoun, error) {
	terms, err := e.store.ListMasterTerms(ctx)
	if err != nil {
		return nil, nil, fmt.Errorf("置換・言及語彙（master_term）の取得: %w", err)
	}
	propers, err := e.store.ListProperNouns(ctx)
	if err != nil {
		return nil, nil, fmt.Errorf("置換・言及語彙（proper_noun）の取得: %w", err)
	}
	keptTerms := make([]model.MasterTerm, 0, len(terms))
	for _, t := range terms {
		if e.stoplist.Blocks(t.Source) {
			continue
		}
		keptTerms = append(keptTerms, t)
	}
	keptPropers := make([]model.ProperNoun, 0, len(propers))
	for _, pn := range propers {
		if e.stoplist.Blocks(pn.Source) {
			continue
		}
		keptPropers = append(keptPropers, pn)
	}
	return keptTerms, keptPropers, nil
}

// DeriveMasterTerms は extracted_field の英日対（dest 非空行）から固有名の確定訳語を master_term へ組む。
// 翻訳 Run の抽出後に呼ぶ。追記した件数を返す。INSERT OR IGNORE のため二重実行でも増えない（冪等）。
// 手順は依存があるため固定する。
//  1. record_type_master で箱＝固有名の FULL（source・dest 非空）を (source, dest, category=rec) で書く
//     （旧 MasterTermXmlWriter の役目を畳んだもの）。
//  2. master_term を baseSources として読む（手順 1 で書いた FULL を含む）。派生の衝突除外に使う。
//  3. NPC_:FULL / NPC_:SHRT・INFO:NAM1 から termderive で人名の部分形（名のみ・短名）を派生して追記する。
//     派生行は category="derive:<種別>" で由来を残す。
func (e *Engine) DeriveMasterTerms(ctx context.Context) (int, error) {
	fields, err := e.store.ListExtractedFields(ctx)
	if err != nil {
		return 0, fmt.Errorf("確定訳語の供給元（extracted_field）の取得: %w", err)
	}
	masterRows, err := e.store.ListRecordTypeMaster(ctx)
	if err != nil {
		return 0, fmt.Errorf("record_type_master の取得: %w", err)
	}
	master := recordMasterMap(masterRows)

	// 手順 1: 箱＝固有名の FULL 完全形を確定訳語として書く。原語と訳語が同一の行は辞書にしない。
	fulls := make([]model.MasterTerm, 0, len(fields))
	for _, f := range fields {
		rt, ok := master[RecordKey{Rec: f.Rec, Field: f.Field}]
		if !ok || rt.Box != boxProperNoun || f.Field != fieldFull {
			continue
		}
		src, dst := strings.TrimSpace(f.Source), strings.TrimSpace(f.Dest)
		if src == "" || dst == "" || src == dst {
			continue
		}
		fulls = append(fulls, model.MasterTerm{Source: src, Dest: dst, Category: f.Rec})
	}
	insertedFulls, err := e.store.InsertDerivedTerms(ctx, fulls)
	if err != nil {
		return 0, fmt.Errorf("固有名の完全形の追記: %w", err)
	}

	// 手順 2: 手順 1 の FULL を含む既存 master_term の原語集合。派生が既出原語と衝突しないための除外に使う。
	terms, err := e.store.ListMasterTerms(ctx)
	if err != nil {
		return insertedFulls, fmt.Errorf("base 辞書の取得: %w", err)
	}
	baseSources := make(map[string]bool, len(terms))
	for _, t := range terms {
		baseSources[t.Source] = true
	}

	// 手順 3: termderive の入力を extracted_field から組む。姓名分割（two）の可否は英日対の有無で決める。
	// dest 非空＝英日 Strings が揃う行だけが対になるため、対を組めた行はすべて two を許す
	//（base ゲーム名（XML ファイル名接頭）による限定は廃止した）。
	var fullPairs, shrtPairs []termderive.NamePair
	var dialogues []string
	for _, f := range fields {
		src, dst := strings.TrimSpace(f.Source), strings.TrimSpace(f.Dest)
		switch {
		case f.Rec == "NPC_" && f.Field == fieldFull:
			if src != "" && dst != "" {
				fullPairs = append(fullPairs, termderive.NamePair{Source: src, Dest: dst, BaseGame: true})
			}
		case f.Rec == "NPC_" && f.Field == "SHRT":
			if src != "" && dst != "" {
				shrtPairs = append(shrtPairs, termderive.NamePair{Source: src, Dest: dst, BaseGame: true})
			}
		case f.Rec == recInfo && f.Field == fieldNam1:
			// 会話文の英語原文は用法分布（一般語か固有名か）の材料。既訳が無くても使える。
			if src != "" {
				dialogues = append(dialogues, f.Source)
			}
		}
	}
	usage := termusage.BuildUsage(dialogues)
	derived := termderive.DeriveTerms(fullPairs, shrtPairs, usage, baseSources, termderive.DefaultDeriveConfig())
	rows := make([]model.MasterTerm, len(derived))
	for i, d := range derived {
		rows[i] = model.MasterTerm{Source: d.Source, Dest: d.Dest, Category: "derive:" + d.Kind}
	}
	insertedDerived, err := e.store.InsertDerivedTerms(ctx, rows)
	if err != nil {
		return insertedFulls + insertedDerived, fmt.Errorf("派生固有名の追記: %w", err)
	}
	return insertedFulls + insertedDerived, nil
}

// ToneDefaults は話者なし台詞（汎用・PC）の口調設定。利用者が編集する自由記述の口調と PC 性別。
type ToneDefaults struct {
	Generic string // 汎用台詞の自由記述口調
	PC      string // PC 発話の自由記述口調
	PcSex   string // PC の性別（Female/Male/空）。一人称・語尾の根拠
}

// LinePersonas は台詞群へ口調指示文（全文）と一覧用の短い要約を組み、lineID→口調 の map で返す。
// 経路は 3 つ。①名指し話者は生成済みの基底口調（話者集計）、②汎用台詞（INFO・話者なし）は自由記述の
// 汎用口調へ本文 1 行の感情段階と条件由来の性別を重ねる、③PC 発話（選択肢）は自由記述の PC 口調へ
// 本文 1 行の感情段階と利用者選択の PC 性別を重ねる。口調指示が空になる台詞は map に現れない（口調なし）。
// personaTemplate は口調指示の雛形（{traits} を含む）で、3 経路すべての性質列を差し込む共通の枠。
func (e *Engine) LinePersonas(ctx context.Context, lines []model.Line, personaTemplate string, defaults ToneDefaults) (map[int64]Persona, error) {
	ids := lineIDsOf(lines)
	resolved, err := e.store.LoadLinePersonas(ctx, ids)
	if err != nil {
		return nil, fmt.Errorf("注入入力（生成ペルソナ）の一括取得: %w", err)
	}
	conds, err := e.store.LoadLineConditions(ctx, ids)
	if err != nil {
		return nil, fmt.Errorf("条件由来の性別の一括取得: %w", err)
	}
	// 話者を解決できない台詞（汎用・PC）の本文特徴をまとめて確保する（line_analysis に本文ハッシュでキャッシュ）。
	features, err := e.ensureFeatures(ctx, pendingLines(lines, resolved))
	if err != nil {
		return nil, err
	}

	out := make(map[int64]Persona, len(lines))
	for _, l := range lines {
		if p, ok := e.personaForLine(l, resolved, conds, features, personaTemplate, defaults); ok {
			out[l.ID] = p
		}
	}
	return out, nil
}

// pendingLines は自由記述の口調へ回す台詞（本文特徴が要る行）を抜き出す。
// PC 発話（選択肢）は INFO 由来で NPC 話者が結ばれ得るが、話者を無視して PC 経路へ回すため常に含める。
// それ以外は話者を解決できない台詞（汎用経路）だけ含める。
func pendingLines(lines []model.Line, resolved map[int64]model.LinePersonaInput) []model.Line {
	var pending []model.Line
	for _, l := range lines {
		if isPCRecord(l.Rec, l.Field) {
			pending = append(pending, l)
			continue
		}
		if _, ok := resolved[l.ID]; !ok {
			pending = append(pending, l)
		}
	}
	return pending
}

// personaForLine は 1 台詞の口調を経路①②③で決める。
// PC 発話（選択肢）は NPC 話者が結ばれていても PC 既定の口調を優先する（③、話者を無視）。
// それ以外は ①名指し話者なら集計済みペルソナ、②話者なし汎用台詞なら本文特徴から自由記述の口調を組む。
// 口調が付かないなら ok=false。
func (e *Engine) personaForLine(l model.Line, resolved map[int64]model.LinePersonaInput, conds map[int64]string, features map[string]tone.Features, personaTemplate string, defaults ToneDefaults) (Persona, bool) {
	if !isPCRecord(l.Rec, l.Field) {
		if in, ok := resolved[l.ID]; ok {
			in.EmotionType = l.EmotionType // 台詞の TRDT 感情を注入入力へ重ねる（ペルソナ基底は不変）。
			return e.namedPersona(in, personaTemplate)
		}
	}
	f, ok := features[linefeatures.SourceHash(l.Source)]
	if !ok {
		return Persona{}, false // 本文が空などで特徴が無い → 口調なし。
	}
	return e.freeTonePersona(l, f, conds[l.ID], personaTemplate, defaults)
}

// namedPersona は名指し話者（経路①）の口調指示と要約を組む。段階が範囲外で性質文が空なら built=false。
func (e *Engine) namedPersona(in model.LinePersonaInput, personaTemplate string) (Persona, bool) {
	directive := personatone.BuildToneDirective(personaTemplate, personatone.BuildToneTraits(in, e.roleSpeech))
	if directive == "" {
		return Persona{}, false
	}
	cell, trait := personatone.PersonaMetaOf(in)
	return Persona{
		Directive:    directive,
		Label:        personatone.BuildToneLabel(in),
		Cell:         cell,
		Trait:        trait,
		AttitudeBand: in.AttitudeBand,
		EmotionBand:  in.EmotionBand,
		Marked:       in.Marked,
		DecisionPath: in.DecisionPath,
	}, true
}

// freeTonePersona は話者なし台詞（経路②③）の口調指示と要約を組む。
// (rec, field) で汎用・PC を振り分け、自由記述の口調へ本文 1 行の感情と性別の一人称・語尾を重ねる。
// 台詞の 3 種以外、または自由記述が空なら built=false。
func (e *Engine) freeTonePersona(l model.Line, f tone.Features, condSex, personaTemplate string, defaults ToneDefaults) (Persona, bool) {
	var label, path, sex, baseText string
	switch {
	case isPCRecord(l.Rec, l.Field):
		label, path, sex, baseText = "口調: PC発話", tone.PathPC, defaults.PcSex, defaults.PC
	case l.Rec == recInfo && l.Field == fieldNam1:
		label, path, sex, baseText = "口調: 汎用台詞", tone.PathGeneric, condSex, defaults.Generic
	default:
		return Persona{}, false // 台詞の 3 種以外は来ない想定。来ても口調なし。
	}
	emotion := e.classifier.EmotionBandOfLine(f)
	directive := personatone.BuildToneDirective(personaTemplate, personatone.BuildFreeToneTraits(baseText, emotion, sex, e.roleSpeech, l.EmotionType))
	if directive == "" {
		return Persona{}, false // 自由記述が空 → 口調なし。
	}
	return Persona{
		Directive:    directive,
		Label:        label,
		EmotionBand:  emotion,
		DecisionPath: path,
		Sex:          sex,
	}, true
}

// isPCRecord はプレイヤーの選択肢の台詞か（PC 発話）を返す。
// DIAL:FULL（選択肢の既定文）と INFO:RNAM（選択肢の条件別上書き）が PC 発話。
func isPCRecord(rec, field string) bool {
	return (rec == recDial && field == fieldFull) || (rec == recInfo && field == fieldRnam)
}

// ensureFeatures は台詞群の本文を、line_analysis に無い分だけ prose で解析して保存し、source_hash→Features を返す。
// 汎用・PC 台詞の感情段階を出すために使う。重複本文は 1 度だけ解析する（プール台詞も 1 回で済む）。
func (e *Engine) ensureFeatures(ctx context.Context, lines []model.Line) (map[string]tone.Features, error) {
	hashToText := make(map[string]string)
	for _, l := range lines {
		if strings.TrimSpace(l.Source) == "" {
			continue
		}
		hashToText[linefeatures.SourceHash(l.Source)] = l.Source
	}
	// hash 群を昇順に固定し、UpsertLineAnalysis の書き込み順を決定的にする（採番 id の非決定を避ける）。
	hashes := make([]string, 0, len(hashToText))
	for h := range hashToText {
		hashes = append(hashes, h)
	}
	sort.Strings(hashes)
	cached, err := e.store.GetLineAnalyses(ctx, hashes)
	if err != nil {
		return nil, fmt.Errorf("行解析キャッシュの取得: %w", err)
	}
	out := make(map[string]tone.Features, len(hashes))
	for _, h := range hashes {
		if err := ctx.Err(); err != nil { // 解析ループは長時間になり得るため、本文ごとにキャンセルを確認する。
			return nil, err
		}
		if row, ok := cached[h]; ok {
			out[h] = featuresFromAnalysis(row)
			continue
		}
		f := linefeatures.ExtractFeatures(hashToText[h], e.lexicon) // prose。キャッシュ未命中の本文だけ実行する。
		if err := e.store.UpsertLineAnalysis(ctx, analysisFromFeatures(h, f)); err != nil {
			return nil, fmt.Errorf("行解析キャッシュの保存: %w", err)
		}
		out[h] = f
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
