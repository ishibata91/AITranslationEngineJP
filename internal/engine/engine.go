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
// LoadLineSpeaker は台詞に紐づく話者の事実上の識別子を返す。話者が無い台詞は found=false を返す。
type LineStore interface {
	ListUntranslatedLines(ctx context.Context) ([]model.Line, error)
	LoadLineSpeaker(ctx context.Context, lineID int64) (model.SpeakerIdentity, bool, error)
	UpdateLineDest(ctx context.Context, id int64, dest string, status int) error
}

// Store は engine が必要とする中心データアクセスをまとめる。concrete は internal/store が 1 つ実装する。
type Store interface {
	NarrationStore
	LineStore
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

	total := len(narrations) + len(lines)
	done := 0
	report := func() {
		if onProgress != nil {
			onProgress(done, total)
		}
	}
	report()

	for _, row := range narrations {
		dest, err := e.provider.Translate(ctx, conn, model, row.Source, "")
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
		directive, _, err := e.linePersona(ctx, row.ID)
		if err != nil {
			return done, err
		}
		dest, err := e.provider.Translate(ctx, conn, model, row.Source, directive)
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

// LineDirective は台詞のペルソナ口調指示文（全文）と、一覧表示用の短い口調要約を返す（画面表示用）。
// 話者が解決できない台詞は両方とも空文字を返す。
func (e *Engine) LineDirective(ctx context.Context, lineID int64) (directive, label string, err error) {
	return e.linePersona(ctx, lineID)
}

// linePersona は台詞の話者識別子を読み、ルールで口調 traits へ写し、口調指示文と短い要約を組む。
// 話者が無ければ両方空文字を返す。
func (e *Engine) linePersona(ctx context.Context, lineID int64) (directive, label string, err error) {
	id, found, err := e.store.LoadLineSpeaker(ctx, lineID)
	if err != nil {
		return "", "", fmt.Errorf("ペルソナ生成の話者識別子取得: %w", err)
	}
	if !found {
		return "", "", nil
	}
	persona := personaFromIdentity(id)
	return buildPersonaDirective(persona), buildPersonaLabel(persona), nil
}
