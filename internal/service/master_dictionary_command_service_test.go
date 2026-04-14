package service

import (
	"context"
	"testing"
	"time"
)

const commandServiceUpdatedSource = "Source B"

func TestMasterDictionaryCommandServiceCreateUsesValidatedRepositoryDraft(t *testing.T) {
	repo := &repositoryStub{}
	repo.createFunc = func(_ context.Context, draft MasterDictionaryDraft) (MasterDictionaryEntry, error) {
		return MasterDictionaryEntry{
			ID:          1,
			Source:      draft.Source,
			Translation: draft.Translation,
			Category:    draft.Category,
			Origin:      draft.Origin,
			REC:         draft.REC,
			EDID:        draft.EDID,
			UpdatedAt:   draft.UpdatedAt,
		}, nil
	}
	service := NewMasterDictionaryCommandService(repo, fixedMasterDictionaryNow)

	created, err := service.CreateEntry(context.Background(), MasterDictionaryMutationInput{
		Source:      "  Source A  ",
		Translation: "  訳語A  ",
		REC:         " BOOK:FULL ",
		EDID:        " EDID_A ",
	})
	if err != nil {
		t.Fatalf("expected create to succeed: %v", err)
	}
	if len(repo.createDrafts) != 1 {
		t.Fatalf("expected one create call, got %d", len(repo.createDrafts))
	}
	createdDraft := repo.createDrafts[0]
	if createdDraft.Source != "Source A" || createdDraft.Translation != "訳語A" {
		t.Fatal("expected create to trim source and translation")
	}
	if createdDraft.Category != "固有名詞" {
		t.Fatalf("expected default category, got %q", createdDraft.Category)
	}
	if createdDraft.Origin != "手動登録" {
		t.Fatalf("expected default origin, got %q", createdDraft.Origin)
	}
	if !createdDraft.UpdatedAt.Equal(fixedMasterDictionaryNow()) {
		t.Fatalf("expected created timestamp %s, got %s", fixedMasterDictionaryNow(), createdDraft.UpdatedAt)
	}
	if created.Category != createdDraft.Category || created.Origin != createdDraft.Origin {
		t.Fatal("expected create result to reflect validated draft")
	}
}

func TestMasterDictionaryCommandServiceUpdateUsesValidatedRepositoryDraft(t *testing.T) {
	repo := &repositoryStub{}
	repo.updateFunc = func(_ context.Context, id int64, draft MasterDictionaryDraft) (MasterDictionaryEntry, error) {
		return MasterDictionaryEntry{
			ID:          id,
			Source:      draft.Source,
			Translation: draft.Translation,
			Category:    draft.Category,
			Origin:      draft.Origin,
			UpdatedAt:   draft.UpdatedAt,
		}, nil
	}
	service := NewMasterDictionaryCommandService(repo, fixedMasterDictionaryNow)

	updated, err := service.UpdateEntry(context.Background(), 1, MasterDictionaryMutationInput{
		Source:      commandServiceUpdatedSource,
		Translation: "訳語B",
		Category:    "地名",
		Origin:      "更新",
	})
	if err != nil {
		t.Fatalf("expected update to succeed: %v", err)
	}
	if len(repo.updateCalls) != 1 {
		t.Fatalf("expected one update call, got %d", len(repo.updateCalls))
	}
	if repo.updateCalls[0].id != 1 {
		t.Fatalf("expected update id 1, got %d", repo.updateCalls[0].id)
	}
	updatedDraft := repo.updateCalls[0].draft
	if updatedDraft.Source != commandServiceUpdatedSource || updatedDraft.Translation != "訳語B" {
		t.Fatal("expected update draft to preserve provided values")
	}
	if updatedDraft.Category != "地名" || updatedDraft.Origin != "更新" {
		t.Fatal("expected update draft to preserve provided metadata")
	}
	if updated.Source != commandServiceUpdatedSource || updated.Translation != "訳語B" {
		t.Fatal("expected update result to be reflected")
	}
}

func TestMasterDictionaryCommandServiceDeleteCallsRepository(t *testing.T) {
	repo := &repositoryStub{}
	service := NewMasterDictionaryCommandService(repo, fixedMasterDictionaryNow)

	if err := service.DeleteEntry(context.Background(), 1); err != nil {
		t.Fatalf("expected delete to succeed: %v", err)
	}
	if len(repo.deleteCalls) != 1 || repo.deleteCalls[0] != 1 {
		t.Fatalf("expected delete to be called with id 1, got %v", repo.deleteCalls)
	}
}

func TestMasterDictionaryCommandServiceRejectsInvalidInput(t *testing.T) {
	service := NewMasterDictionaryCommandService(&repositoryStub{}, fixedMasterDictionaryNow)

	if _, err := service.CreateEntry(context.Background(), MasterDictionaryMutationInput{Source: "  ", Translation: "訳語"}); err == nil {
		t.Fatal("expected create with empty source to fail")
	}
	if _, err := service.CreateEntry(context.Background(), MasterDictionaryMutationInput{Source: "Source", Translation: "  "}); err == nil {
		t.Fatal("expected create with empty translation to fail")
	}
	if _, err := service.UpdateEntry(context.Background(), 0, MasterDictionaryMutationInput{Source: "S", Translation: "T"}); err == nil {
		t.Fatal("expected update with invalid id to fail")
	}
	if err := service.DeleteEntry(context.Background(), 0); err == nil {
		t.Fatal("expected delete with invalid id to fail")
	}
}

func fixedMasterDictionaryNow() time.Time {
	return time.Date(2026, 4, 12, 10, 0, 0, 0, time.UTC)
}
