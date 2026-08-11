package api

import (
	"errors"
	"fmt"

	"aitranslationenginejp/internal/model"
)

type TermDictionaryFilterView struct {
	Source       string `json:"source"`
	Destination  string `json:"destination"`
	PartOfSpeech string `json:"partOfSpeech"`
	Category     string `json:"category"`
}
type TermDictionaryEntryView struct {
	ID           int64    `json:"id"`
	Source       string   `json:"source"`
	Destination  string   `json:"destination"`
	PartOfSpeech string   `json:"partOfSpeech"`
	Revision     int64    `json:"revision"`
	Categories   []string `json:"categories"`
}
type TermDictionaryPageView struct {
	Entries    []TermDictionaryEntryView `json:"entries"`
	TotalCount int                       `json:"totalCount"`
	PageNumber int                       `json:"pageNumber"`
}
type TermDictionaryCreateRequest struct {
	Source       string `json:"source"`
	Destination  string `json:"destination"`
	PartOfSpeech string `json:"partOfSpeech"`
}
type TermDictionaryPatchRequest struct {
	ID           int64   `json:"id"`
	Revision     int64   `json:"revision"`
	Source       *string `json:"source,omitempty"`
	Destination  *string `json:"destination,omitempty"`
	PartOfSpeech *string `json:"partOfSpeech,omitempty"`
}
type TermDictionaryDeleteRequest struct {
	ID       int64 `json:"id"`
	Revision int64 `json:"revision"`
}

func (a *App) ListTermDictionary(filter TermDictionaryFilterView, pageNumber int) (TermDictionaryPageView, error) {
	editor, err := a.editor()
	if err != nil {
		return TermDictionaryPageView{}, err
	}
	page, err := editor.List(a.baseCtx(), model.TermDictionaryFilter{Source: filter.Source, Destination: filter.Destination, PartOfSpeech: filter.PartOfSpeech, Category: filter.Category}, pageNumber)
	if err != nil {
		return TermDictionaryPageView{}, fmt.Errorf("用語辞書の一覧取得: %w", err)
	}
	return termDictionaryPageView(page), nil
}
func (a *App) CreateTermDictionary(req TermDictionaryCreateRequest) (TermDictionaryEntryView, error) {
	editor, err := a.editor()
	if err != nil {
		return TermDictionaryEntryView{}, err
	}
	entry, err := editor.Create(a.baseCtx(), model.TermDictionaryCreate{Source: req.Source, Destination: req.Destination, PartOfSpeech: req.PartOfSpeech})
	if err != nil {
		return TermDictionaryEntryView{}, fmt.Errorf("用語辞書の作成: %w", err)
	}
	return termDictionaryEntryView(entry), nil
}
func (a *App) PatchTermDictionary(req TermDictionaryPatchRequest) (TermDictionaryEntryView, error) {
	editor, err := a.editor()
	if err != nil {
		return TermDictionaryEntryView{}, err
	}
	entry, err := editor.Patch(a.baseCtx(), model.TermDictionaryPatch{ID: req.ID, Revision: req.Revision, Source: req.Source, Destination: req.Destination, PartOfSpeech: req.PartOfSpeech})
	if err != nil {
		return TermDictionaryEntryView{}, fmt.Errorf("用語辞書の更新: %w", err)
	}
	return termDictionaryEntryView(entry), nil
}
func (a *App) DeleteTermDictionary(req TermDictionaryDeleteRequest) error {
	editor, err := a.editor()
	if err != nil {
		return err
	}
	if err := editor.Delete(a.baseCtx(), req.ID, req.Revision); err != nil {
		return fmt.Errorf("用語辞書の削除: %w", err)
	}
	return nil
}
func (a *App) editor() (TermDictionaryEditorStore, error) {
	if a.termDictionaryEditor == nil {
		return nil, errors.New("用語辞書編集APIは未初期化")
	}
	return a.termDictionaryEditor, nil
}
func termDictionaryPageView(page model.TermDictionaryPage) TermDictionaryPageView {
	entries := make([]TermDictionaryEntryView, len(page.Entries))
	for i, entry := range page.Entries {
		entries[i] = termDictionaryEntryView(entry)
	}
	return TermDictionaryPageView{Entries: entries, TotalCount: page.TotalCount, PageNumber: page.PageNumber}
}
func termDictionaryEntryView(entry model.TermDictionaryEntry) TermDictionaryEntryView {
	return TermDictionaryEntryView{ID: entry.ID, Source: entry.Source, Destination: entry.Destination, PartOfSpeech: entry.PartOfSpeech, Revision: entry.Revision, Categories: entry.Categories}
}
