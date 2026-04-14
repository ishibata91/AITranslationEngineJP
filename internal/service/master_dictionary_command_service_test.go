package service

import (
	"context"
	"testing"
	"time"
)

const (
	commandServiceCreateSource  = "Source A"
	commandServiceUpdatedSource = "Source B"
	errCommandCreateSucceeded   = "expected create to succeed: %v"
)

func TestMasterDictionaryCommandServiceCreateTrimsMutationInputBeforeRepositoryCall(t *testing.T) {
	repo := &repositoryStub{}
	repo.createFunc = func(_ context.Context, draft MasterDictionaryDraft) (MasterDictionaryEntry, error) {
		return MasterDictionaryEntry{ID: 1, Source: draft.Source, Translation: draft.Translation}, nil
	}
	service := NewMasterDictionaryCommandService(repo, fixedMasterDictionaryNow)

	_, err := service.CreateEntry(context.Background(), MasterDictionaryMutationInput{
		Source:      "  " + commandServiceCreateSource + "  ",
		Translation: "  訳語A  ",
		REC:         " BOOK:FULL ",
		EDID:        " EDID_A ",
	})
	if err != nil {
		t.Fatalf(errCommandCreateSucceeded, err)
	}

	createdDraft := repo.createDrafts[0]
	if createdDraft.Source != commandServiceCreateSource {
		t.Fatalf("expected create to trim source, got %q", createdDraft.Source)
	}
	if createdDraft.Translation != "訳語A" {
		t.Fatalf("expected create to trim translation, got %q", createdDraft.Translation)
	}
}

func TestMasterDictionaryCommandServiceCreateAppliesDefaultMetadata(t *testing.T) {
	repo := &repositoryStub{}
	repo.createFunc = func(_ context.Context, draft MasterDictionaryDraft) (MasterDictionaryEntry, error) {
		return MasterDictionaryEntry{ID: 1, Category: draft.Category, Origin: draft.Origin, UpdatedAt: draft.UpdatedAt}, nil
	}
	service := NewMasterDictionaryCommandService(repo, fixedMasterDictionaryNow)

	_, err := service.CreateEntry(context.Background(), MasterDictionaryMutationInput{
		Source: commandServiceCreateSource, Translation: "訳語A",
	})
	if err != nil {
		t.Fatalf(errCommandCreateSucceeded, err)
	}

	createdDraft := repo.createDrafts[0]
	if createdDraft.Category != "固有名詞" {
		t.Fatalf("expected default category, got %q", createdDraft.Category)
	}
	if createdDraft.Origin != "手動登録" {
		t.Fatalf("expected default origin, got %q", createdDraft.Origin)
	}
	if !createdDraft.UpdatedAt.Equal(fixedMasterDictionaryNow()) {
		t.Fatalf("expected created timestamp %s, got %s", fixedMasterDictionaryNow(), createdDraft.UpdatedAt)
	}
}

func TestMasterDictionaryCommandServiceCreateReturnsRepositoryEntry(t *testing.T) {
	repo := &repositoryStub{}
	repo.createFunc = func(_ context.Context, draft MasterDictionaryDraft) (MasterDictionaryEntry, error) {
		return MasterDictionaryEntry{ID: 1, Category: draft.Category, Origin: draft.Origin}, nil
	}
	service := NewMasterDictionaryCommandService(repo, fixedMasterDictionaryNow)

	created, err := service.CreateEntry(context.Background(), MasterDictionaryMutationInput{
		Source: commandServiceCreateSource, Translation: "訳語A",
	})
	if err != nil {
		t.Fatalf(errCommandCreateSucceeded, err)
	}

	if created.Category != "固有名詞" {
		t.Fatalf("expected result category to match repository entry, got %q", created.Category)
	}
	if created.Origin != "手動登録" {
		t.Fatalf("expected result origin to match repository entry, got %q", created.Origin)
	}
}

func TestMasterDictionaryCommandServiceUpdatePassesValidatedDraftToRepository(t *testing.T) {
	repo := &repositoryStub{}
	repo.updateFunc = func(_ context.Context, id int64, draft MasterDictionaryDraft) (MasterDictionaryEntry, error) {
		return MasterDictionaryEntry{ID: id, Source: draft.Source, Translation: draft.Translation, Category: draft.Category, Origin: draft.Origin}, nil
	}
	service := NewMasterDictionaryCommandService(repo, fixedMasterDictionaryNow)

	_, err := service.UpdateEntry(context.Background(), 1, MasterDictionaryMutationInput{
		Source:      commandServiceUpdatedSource,
		Translation: "訳語B",
		Category:    "地名",
		Origin:      "更新",
	})
	if err != nil {
		t.Fatalf("expected update to succeed: %v", err)
	}

	if repo.updateCalls[0].id != 1 {
		t.Fatalf("expected update id 1, got %d", repo.updateCalls[0].id)
	}
	updatedDraft := repo.updateCalls[0].draft
	if updatedDraft.Source != commandServiceUpdatedSource {
		t.Fatalf("expected update draft source %q, got %q", commandServiceUpdatedSource, updatedDraft.Source)
	}
	if updatedDraft.Translation != "訳語B" {
		t.Fatalf("expected update draft translation 訳語B, got %q", updatedDraft.Translation)
	}
	if updatedDraft.Category != "地名" {
		t.Fatalf("expected update draft category 地名, got %q", updatedDraft.Category)
	}
	if updatedDraft.Origin != "更新" {
		t.Fatalf("expected update draft origin 更新, got %q", updatedDraft.Origin)
	}
}

func TestMasterDictionaryCommandServiceUpdateReturnsRepositoryEntry(t *testing.T) {
	repo := &repositoryStub{}
	repo.updateFunc = func(_ context.Context, id int64, draft MasterDictionaryDraft) (MasterDictionaryEntry, error) {
		return MasterDictionaryEntry{ID: id, Source: draft.Source, Translation: draft.Translation}, nil
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

	if updated.Source != commandServiceUpdatedSource {
		t.Fatalf("expected update result source %q, got %q", commandServiceUpdatedSource, updated.Source)
	}
	if updated.Translation != "訳語B" {
		t.Fatalf("expected update result translation 訳語B, got %q", updated.Translation)
	}
}

func TestMasterDictionaryCommandServiceDeleteCallsRepository(t *testing.T) {
	repo := &repositoryStub{}
	service := NewMasterDictionaryCommandService(repo, fixedMasterDictionaryNow)

	err := service.DeleteEntry(context.Background(), 1)
	if err != nil {
		t.Fatalf("expected delete to succeed: %v", err)
	}

	if len(repo.deleteCalls) != 1 || repo.deleteCalls[0] != 1 {
		t.Fatalf("expected delete to be called with id 1, got %v", repo.deleteCalls)
	}
}

func TestMasterDictionaryCommandServiceCreateRejectsEmptySource(t *testing.T) {
	service := NewMasterDictionaryCommandService(&repositoryStub{}, fixedMasterDictionaryNow)

	_, err := service.CreateEntry(context.Background(), MasterDictionaryMutationInput{Source: "  ", Translation: "訳語"})
	if err == nil {
		t.Fatal("expected create with empty source to fail")
	}
}

func TestMasterDictionaryCommandServiceCreateRejectsEmptyTranslation(t *testing.T) {
	service := NewMasterDictionaryCommandService(&repositoryStub{}, fixedMasterDictionaryNow)

	_, err := service.CreateEntry(context.Background(), MasterDictionaryMutationInput{Source: "Source", Translation: "  "})
	if err == nil {
		t.Fatal("expected create with empty translation to fail")
	}
}

func TestMasterDictionaryCommandServiceUpdateRejectsInvalidID(t *testing.T) {
	service := NewMasterDictionaryCommandService(&repositoryStub{}, fixedMasterDictionaryNow)

	_, err := service.UpdateEntry(context.Background(), 0, MasterDictionaryMutationInput{Source: "S", Translation: "T"})
	if err == nil {
		t.Fatal("expected update with invalid id to fail")
	}
}

func TestMasterDictionaryCommandServiceDeleteRejectsInvalidID(t *testing.T) {
	service := NewMasterDictionaryCommandService(&repositoryStub{}, fixedMasterDictionaryNow)

	err := service.DeleteEntry(context.Background(), 0)
	if err == nil {
		t.Fatal("expected delete with invalid id to fail")
	}
}

func fixedMasterDictionaryNow() time.Time {
	return time.Date(2026, 4, 12, 10, 0, 0, 0, time.UTC)
}
