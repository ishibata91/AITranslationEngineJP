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
// LoadLineSpeakers は台詞 id 群に紐づく話者の事実上の識別子を一括で返す。話者が無い台詞は map に現れない。
type LineStore interface {
	ListUntranslatedLines(ctx context.Context) ([]model.Line, error)
	LoadLineSpeakers(ctx context.Context, lineIDs []int64) (map[int64]model.SpeakerIdentity, error)
	UpdateLineDest(ctx context.Context, id int64, dest string, status int) error
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
}

// Engine は翻訳手続きを実行する。
type Engine struct {
	store    Store
	provider provider.Translator
}

// New は engine を生成する。provider は AI 翻訳の port。
func New(store Store, p provider.Translator) *Engine {
	return &Engine{store: store, provider: p}
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

	// 台詞の話者属性をループ前に 1 度だけ一括取得し、ループ内の個別 DB 問い合わせ（N+1）を避ける。
	// 口調指示テンプレートの {traits} へ話者の性質列を差し込んで口調指示文を組む。
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

// LinePersonas は台詞 id 群の話者を一括取得し、各台詞の口調指示文（全文）と一覧用の短い要約を map で返す。
// personaTemplate は口調指示の雛形（{traits} を含む）で、話者の性質列を差し込んで口調指示文を組む。
// 話者が解決できない台詞、および既知属性が無く口調 traits が空になる台詞は map に現れない
// （呼び出し側は欠落を「口調なし」として扱う）。話者取得を 1 度の一括問い合わせにまとめ、台詞数に依存させない。
func (e *Engine) LinePersonas(ctx context.Context, lineIDs []int64, personaTemplate string) (map[int64]Persona, error) {
	speakers, err := e.store.LoadLineSpeakers(ctx, lineIDs)
	if err != nil {
		return nil, fmt.Errorf("ペルソナ生成の話者識別子一括取得: %w", err)
	}
	out := make(map[int64]Persona, len(speakers))
	for lineID, identity := range speakers {
		persona := personaFromIdentity(identity)
		directive := buildPersonaDirective(personaTemplate, persona)
		if directive == "" {
			continue // 既知属性が無く口調が組めない話者は口調なし扱い。
		}
		out[lineID] = Persona{Directive: directive, Label: buildPersonaLabel(persona)}
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
