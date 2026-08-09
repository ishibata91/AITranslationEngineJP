package store

import (
	"context"
	"path/filepath"
	"testing"

	"aitranslationenginejp/internal/model"
)

func TestTranslationReferenceSnapshotRoundTripOmitsMeaning(t *testing.T) {
	s, err := Open(filepath.Join(t.TempDir(), "central.sqlite3"))
	if err != nil {
		t.Fatalf("Open: %v", err)
	}
	defer func() { _ = s.Close() }()
	ctx := context.Background()
	input := model.TranslationReferenceSnapshot{Plugin: "A.esp", Kind: model.BatchKindNarration, RowID: 7, PromptHash: "hash", References: []model.TranslationReference{{
		Source: "Riften", Dest: "リフテン", PartOfSpeech: "noun", SkyrimCategory: "city", Origin: "事前作成済み翻訳辞書",
	}}}
	if err := s.UpsertTranslationReferenceSnapshot(ctx, input); err != nil {
		t.Fatalf("UpsertTranslationReferenceSnapshot: %v", err)
	}
	got, ok, err := s.GetTranslationReferenceSnapshot(ctx, "A.esp", model.BatchKindNarration, 7)
	if err != nil || !ok {
		t.Fatalf("GetTranslationReferenceSnapshot: ok=%v err=%v", ok, err)
	}
	if got.PromptHash != input.PromptHash || len(got.References) != 1 || got.References[0] != input.References[0] {
		t.Errorf("snapshot = %+v, want %+v", got, input)
	}
}
