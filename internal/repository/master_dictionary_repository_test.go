package repository

import (
	"context"
	"errors"
	"testing"
	"time"
)

const (
	inMemoryRepositoryLocationCategory    = "地名"
	inMemoryRepositoryWhiterunTranslation = "ホワイトラン"
	inMemoryRepositoryLocationREC         = "LCTN:FULL"
	inMemoryRepositoryAurielsBow          = "Auriel's Bow"
	inMemoryRepositoryWeaponREC           = "WEAP:FULL"
)

func TestInMemoryMasterDictionaryRepositoryCreateStoresEntry(t *testing.T) {
	repository := NewInMemoryMasterDictionaryRepository(nil)
	now := time.Date(2026, 4, 12, 10, 0, 0, 0, time.UTC)

	created, err := repository.Create(context.Background(), MasterDictionaryDraft{
		Source:      "Whiterun",
		Translation: inMemoryRepositoryWhiterunTranslation,
		Category:    inMemoryRepositoryLocationCategory,
		Origin:      "手動登録",
		REC:         inMemoryRepositoryLocationREC,
		EDID:        "LocWhiterun",
		UpdatedAt:   now,
	})
	if err != nil {
		t.Fatalf("expected create to succeed: %v", err)
	}

	if created.Source != "Whiterun" {
		t.Fatalf("expected created source, got %q", created.Source)
	}
}

func TestInMemoryMasterDictionaryRepositoryGetByIDReturnsCreatedEntry(t *testing.T) {
	repository := NewInMemoryMasterDictionaryRepository(nil)
	now := time.Date(2026, 4, 12, 10, 0, 0, 0, time.UTC)
	created := createInMemoryMasterDictionaryEntry(t, repository, MasterDictionaryDraft{
		Source:      "Whiterun",
		Translation: inMemoryRepositoryWhiterunTranslation,
		Category:    inMemoryRepositoryLocationCategory,
		Origin:      "手動登録",
		REC:         inMemoryRepositoryLocationREC,
		EDID:        "LocWhiterun",
		UpdatedAt:   now,
	})

	loaded, err := repository.GetByID(context.Background(), created.ID)
	if err != nil {
		t.Fatalf("expected get to succeed: %v", err)
	}

	if loaded.Source != "Whiterun" {
		t.Fatalf("expected created source, got %q", loaded.Source)
	}
}

func TestInMemoryMasterDictionaryRepositoryUpdateReplacesEntryValues(t *testing.T) {
	repository := NewInMemoryMasterDictionaryRepository(nil)
	now := time.Date(2026, 4, 12, 10, 0, 0, 0, time.UTC)
	created := createInMemoryMasterDictionaryEntry(t, repository, MasterDictionaryDraft{
		Source:      "Whiterun",
		Translation: inMemoryRepositoryWhiterunTranslation,
		Category:    inMemoryRepositoryLocationCategory,
		Origin:      "手動登録",
		REC:         inMemoryRepositoryLocationREC,
		EDID:        "LocWhiterun",
		UpdatedAt:   now,
	})

	updated, err := repository.Update(context.Background(), created.ID, MasterDictionaryDraft{
		Source:      "Solitude",
		Translation: "ソリチュード",
		Category:    inMemoryRepositoryLocationCategory,
		Origin:      "更新",
		REC:         inMemoryRepositoryLocationREC,
		EDID:        "LocSolitude",
		UpdatedAt:   now.Add(time.Minute),
	})
	if err != nil {
		t.Fatalf("expected update to succeed: %v", err)
	}

	if updated.Source != "Solitude" {
		t.Fatalf("expected updated source, got %q", updated.Source)
	}
}

func TestInMemoryMasterDictionaryRepositoryListReturnsFilteredEntries(t *testing.T) {
	repository := NewInMemoryMasterDictionaryRepository(nil)
	now := time.Date(2026, 4, 12, 10, 0, 0, 0, time.UTC)
	createInMemoryMasterDictionaryEntry(t, repository, MasterDictionaryDraft{
		Source:      "Solitude",
		Translation: "ソリチュード",
		Category:    inMemoryRepositoryLocationCategory,
		Origin:      "更新",
		REC:         inMemoryRepositoryLocationREC,
		EDID:        "LocSolitude",
		UpdatedAt:   now,
	})

	listed, err := repository.List(context.Background(), MasterDictionaryListQuery{
		SearchTerm: "Solitude",
		Category:   inMemoryRepositoryLocationCategory,
		Page:       1,
		PageSize:   30,
	})
	if err != nil {
		t.Fatalf("expected list to succeed: %v", err)
	}

	if listed.TotalCount != 1 {
		t.Fatalf("expected one listed item, got total=%d", listed.TotalCount)
	}
}

func TestInMemoryMasterDictionaryRepositoryDeleteRemovesEntry(t *testing.T) {
	repository := NewInMemoryMasterDictionaryRepository(nil)
	now := time.Date(2026, 4, 12, 10, 0, 0, 0, time.UTC)
	created := createInMemoryMasterDictionaryEntry(t, repository, MasterDictionaryDraft{
		Source:      "Whiterun",
		Translation: inMemoryRepositoryWhiterunTranslation,
		Category:    inMemoryRepositoryLocationCategory,
		Origin:      "手動登録",
		REC:         inMemoryRepositoryLocationREC,
		EDID:        "LocWhiterun",
		UpdatedAt:   now,
	})

	err := repository.Delete(context.Background(), created.ID)
	if err != nil {
		t.Fatalf("expected delete to succeed: %v", err)
	}

	_, err = repository.GetByID(context.Background(), created.ID)
	if !errors.Is(err, ErrMasterDictionaryEntryNotFound) {
		t.Fatalf("expected not found after delete, got %v", err)
	}
}

func TestInMemoryMasterDictionaryRepositoryUpsertBySourceAndRECUpdatesExistingRecord(t *testing.T) {
	now := time.Date(2026, 4, 12, 10, 0, 0, 0, time.UTC)
	repository := NewInMemoryMasterDictionaryRepository([]MasterDictionaryEntry{{
		ID:          1,
		Source:      inMemoryRepositoryAurielsBow,
		Translation: "旧訳",
		Category:    "装備",
		Origin:      "初期データ",
		REC:         inMemoryRepositoryWeaponREC,
		EDID:        "DLC1AurielsBow",
		UpdatedAt:   now,
	}})

	entry, created, err := repository.UpsertBySourceAndREC(context.Background(), MasterDictionaryImportRecord{
		Source:      inMemoryRepositoryAurielsBow,
		Translation: "アーリエルの弓",
		Category:    "装備",
		Origin:      "XML取込",
		REC:         inMemoryRepositoryWeaponREC,
		EDID:        "DLC1AurielsBow",
		UpdatedAt:   now.Add(time.Minute),
	})
	if err != nil {
		t.Fatalf("expected update upsert to succeed: %v", err)
	}

	if created {
		t.Fatal("expected existing record to be updated")
	}
	if entry.Translation != "アーリエルの弓" {
		t.Fatalf("expected translation to be updated, got %q", entry.Translation)
	}
}

func TestInMemoryMasterDictionaryRepositoryUpsertBySourceAndRECCreatesNewRecord(t *testing.T) {
	now := time.Date(2026, 4, 12, 10, 0, 0, 0, time.UTC)
	repository := NewInMemoryMasterDictionaryRepository([]MasterDictionaryEntry{{
		ID:          1,
		Source:      inMemoryRepositoryAurielsBow,
		Translation: "旧訳",
		Category:    "装備",
		Origin:      "初期データ",
		REC:         inMemoryRepositoryWeaponREC,
		EDID:        "DLC1AurielsBow",
		UpdatedAt:   now,
	}})

	entry, created, err := repository.UpsertBySourceAndREC(context.Background(), MasterDictionaryImportRecord{
		Source:      "Snow Elf",
		Translation: "スノーエルフ",
		Category:    "NPC",
		Origin:      "XML取込",
		REC:         "NPC_:FULL",
		EDID:        "DLC1SnowElf",
		UpdatedAt:   now.Add(2 * time.Minute),
	})
	if err != nil {
		t.Fatalf("expected create upsert to succeed: %v", err)
	}

	if !created {
		t.Fatal("expected new record to be created")
	}
	if entry.ID == 1 {
		t.Fatalf("expected a new id, got %d", entry.ID)
	}
}

func createInMemoryMasterDictionaryEntry(t *testing.T, repository *InMemoryMasterDictionaryRepository, draft MasterDictionaryDraft) MasterDictionaryEntry {
	t.Helper()

	entry, err := repository.Create(context.Background(), draft)
	if err != nil {
		t.Fatalf("expected create to succeed: %v", err)
	}

	return entry
}
