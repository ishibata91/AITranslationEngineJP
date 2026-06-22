// Package engine は翻訳手続きの本体を持つ。中心データを読み、ペルソナ生成を経て AI 翻訳し、
// 配置へ書き戻す純 Go の手続き。GUI から切り離して単体テストできる。
package engine

import (
	"context"
	"fmt"

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

// Persona は台詞 1 件へ与える口調指示文（全文）と一覧表示用の短い要約。話者を解決できた台詞だけ持つ。
type Persona struct {
	Directive string
	Label     string
}

// DictStore は engine が固有名の機械置換辞書（master_term）を読むための中心データアクセス。
type DictStore interface {
	ListMasterTerms(ctx context.Context) ([]model.MasterTerm, error)
}

// TemplateStore は engine がプロンプトテンプレート（base 指示・口調指示の雛形）を読むための中心データアクセス。
type TemplateStore interface {
	GetPromptTemplate(ctx context.Context) (model.PromptTemplate, error)
}

// Store は engine が必要とする中心データアクセスをまとめる。concrete は internal/store が 1 つ実装する。
type Store interface {
	NarrationStore
	LineStore
	DictStore
	TemplateStore
	PersonaStore
}

// Engine は翻訳手続きを実行する。
type Engine struct {
	store    Store
	provider provider.Translator
	lexicon  EmotionLexicon
}

// New は engine を生成する。provider は AI 翻訳の port、lexicon は口調生成の感情辞書（差し替え可能な境界）。
func New(store Store, p provider.Translator, lexicon EmotionLexicon) *Engine {
	return &Engine{store: store, provider: p, lexicon: lexicon}
}

// GeneratePersonas は台詞を持つ全話者の基底口調を一括生成し persona_character へ保存する。保存した話者数を返す。
// 最も重い行解析は line_analysis に本文ハッシュでキャッシュし、本文ごとに 1 度だけ prose を回す。
// 冪等で、手修正済み（hand_edited）の行は store が再生成から保護する。
func (e *Engine) GeneratePersonas(ctx context.Context) (int, error) {
	return NewPersonaGenerator(e.store, e.lexicon).Generate(ctx)
}

// Run は未訳の叙述文と台詞を順に翻訳し、訳文を仮訳として書き戻す。翻訳できた合計件数を返す。
// 叙述文は base 指示だけで訳し、台詞は話者属性からのペルソナ口調指示を注入して訳す。
// onProgress が非 nil なら本文翻訳 phase の進捗（処理済み件数 done、総件数 total）を都度通知する。
func (e *Engine) Run(ctx context.Context, conn provider.Connection, model string, onProgress func(done, total int)) (int, error) {
	narrations, err := e.store.ListUntranslatedNarrations(ctx)
	if err != nil {
		return 0, fmt.Errorf("未訳叙述文の取得: %w", err)
	}
	lines, err := e.store.ListUntranslatedLines(ctx)
	if err != nil {
		return 0, fmt.Errorf("未訳台詞の取得: %w", err)
	}

	// プロンプトテンプレート（base 指示・口調指示の雛形）をループ前に 1 度だけ読む。
	// base 指示と口調指示の組み立てに使い、保存済みテンプレートを翻訳へ反映する。
	tmpl, err := e.store.GetPromptTemplate(ctx)
	if err != nil {
		return 0, fmt.Errorf("プロンプトテンプレートの取得: %w", err)
	}

	// 口調ペルソナを一括生成して persona_character を最新化する。最も重い行解析は line_analysis に
	// 本文ハッシュでキャッシュし、本文ごとに 1 度だけ prose を回す。注入はこの生成結果を引く。
	if _, err := e.GeneratePersonas(ctx); err != nil {
		return 0, err
	}

	// 台詞ごとの口調指示をループ前に 1 度だけ一括で組み、ループ内の個別 DB 問い合わせ（N+1）を避ける。
	// 生成済みの基底口調を性質文カタログへ写し、口調指示テンプレートの {traits} へ差し込む。
	personas, err := e.LinePersonas(ctx, lineIDsOf(lines), tmpl.PersonaTemplate)
	if err != nil {
		return 0, err
	}

	// 固有名の機械置換辞書をループ前に 1 度だけ組む。本文翻訳の前に原文の固有名を確定訳語へ置換する。
	dict, err := e.LoadDictionary(ctx)
	if err != nil {
		return 0, err
	}

	total := len(narrations) + len(lines)
	done := 0
	report := func() {
		if onProgress != nil {
			onProgress(done, total)
		}
	}
	report()

	for _, row := range narrations {
		// 本文中の固有名を辞書の確定訳語へ機械置換してから AI 翻訳する。AI は周りの英語だけを訳す。
		source, _ := dict.Apply(row.Source)
		// 叙述文は話者を持たないため口調指示なし。テンプレートの base 指示だけで完成プロンプトを組んで送る。
		dest, err := e.provider.Translate(ctx, conn, model, ComposePrompt(tmpl.BaseDirective, "", source))
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
		// 台詞は話者の口調指示をテンプレートの base 指示へ合成した完成プロンプトを組んで送る。話者なしなら口調指示は空。
		dest, err := e.provider.Translate(ctx, conn, model, ComposePrompt(tmpl.BaseDirective, personas[row.ID].Directive, source))
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

// LoadDictionary は master_term を読み、固有名の機械置換辞書を組む。
// 翻訳実行（Run）と結果取得（api の内訳・実プロンプト再構成）の両方が、ページ単位で 1 度だけ組んで使う。
func (e *Engine) LoadDictionary(ctx context.Context) (*Dictionary, error) {
	terms, err := e.store.ListMasterTerms(ctx)
	if err != nil {
		return nil, fmt.Errorf("マスター辞書の取得: %w", err)
	}
	pairs := make([]DictionaryTerm, len(terms))
	for i, term := range terms {
		pairs[i] = DictionaryTerm{Source: term.Source, Dest: term.Dest}
	}
	return NewDictionary(pairs), nil
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
		directive := buildToneDirective(personaTemplate, buildToneTraits(in))
		if directive == "" {
			continue // 段階が範囲外で性質文が組めない場合は口調なし扱い。
		}
		out[lineID] = Persona{Directive: directive, Label: buildToneLabel(in)}
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
