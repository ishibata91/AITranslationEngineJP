package engine

import (
	"context"
	"fmt"

	"aitranslationenginejp/internal/core/japanesetext"
	"aitranslationenginejp/internal/model"
)

func (e *Engine) completeJapaneseProperNouns(ctx context.Context, rows []model.ProperNoun) ([]model.ProperNoun, int, error) {
	pending := make([]model.ProperNoun, 0, len(rows))
	completed := 0
	for _, row := range rows {
		if !japanesetext.Contains(row.Source) {
			pending = append(pending, row)
			continue
		}
		if err := e.store.UpdateProperNounDest(ctx, row.ID, row.Source, statusTranslated); err != nil {
			return nil, 0, fmt.Errorf("日本語の固有名の書き戻し: %w", err)
		}
		completed++
	}
	return pending, completed, nil
}

func (e *Engine) completeJapaneseNarrations(ctx context.Context, rows []model.Narration) ([]model.Narration, int, error) {
	pending := make([]model.Narration, 0, len(rows))
	completed := 0
	for _, row := range rows {
		if !japanesetext.Contains(row.Source) {
			pending = append(pending, row)
			continue
		}
		if err := e.store.UpdateNarrationDest(ctx, row.ID, row.Source, statusTranslated); err != nil {
			return nil, 0, fmt.Errorf("日本語の叙述文の書き戻し: %w", err)
		}
		completed++
	}
	return pending, completed, nil
}

func (e *Engine) completeJapaneseLines(ctx context.Context, rows []model.Line) ([]model.Line, int, error) {
	pending := make([]model.Line, 0, len(rows))
	completed := 0
	for _, row := range rows {
		if !japanesetext.Contains(row.Source) {
			pending = append(pending, row)
			continue
		}
		if err := e.store.UpdateLineDest(ctx, row.ID, row.Source, statusTranslated); err != nil {
			return nil, 0, fmt.Errorf("日本語の台詞の書き戻し: %w", err)
		}
		completed++
	}
	return pending, completed, nil
}
