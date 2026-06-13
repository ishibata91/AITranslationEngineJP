// Package engine は翻訳手続きの本体を持つ。中心データを読み、AI 翻訳し、配置へ書き戻す純 Go の手続き。
// GUI から切り離して単体テストできる。
package engine

import (
	"context"
	"fmt"

	"aitranslationenginejp/internal/model"
	"aitranslationenginejp/internal/provider"
)

// statusProvisional は xTranslator の訳状態 3（仮）。AI 翻訳は仮訳として書き戻す。
const statusProvisional = 3

// NarrationStore は engine が必要とする中心データアクセスの interface（使う分だけ宣言する）。
type NarrationStore interface {
	ListUntranslatedNarrations(ctx context.Context) ([]model.Narration, error)
	UpdateNarrationDest(ctx context.Context, id int64, dest string, status int) error
}

// Engine は翻訳手続きを実行する。
type Engine struct {
	store    NarrationStore
	provider provider.Translator
}

// New は engine を生成する。
func New(store NarrationStore, p provider.Translator) *Engine {
	return &Engine{store: store, provider: p}
}

// Run は未訳の叙述文を順に翻訳し、訳文を仮訳として書き戻す。翻訳できた件数を返す。
func (e *Engine) Run(ctx context.Context, conn provider.Connection, model string) (int, error) {
	rows, err := e.store.ListUntranslatedNarrations(ctx)
	if err != nil {
		return 0, fmt.Errorf("未訳の取得: %w", err)
	}
	translated := 0
	for _, row := range rows {
		dest, err := e.provider.Translate(ctx, conn, model, row.Source)
		if err != nil {
			return translated, fmt.Errorf("翻訳: %w", err)
		}
		if err := e.store.UpdateNarrationDest(ctx, row.ID, dest, statusProvisional); err != nil {
			return translated, fmt.Errorf("訳文の書き戻し: %w", err)
		}
		translated++
	}
	return translated, nil
}
