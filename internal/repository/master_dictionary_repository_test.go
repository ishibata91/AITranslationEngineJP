package repository

import (
	"context"
	"errors"
	"testing"
	"time"
)

func TestInMemoryMasterDictionaryRepositoryCRUDAndList(t *testing.T) {
	repository := NewInMemoryMasterDictionaryRepository(nil)
	now := time.Date(2026, 4, 12, 10, 0, 0, 0, time.UTC)

	created, err := repository.Create(context.Background(), MasterDictionaryDraft{
		Source:      "Whiterun",
		Translation: "ホワイトラン",
		Category:    "地名",
		Origin:      "手動登録",
		REC:         "LCTN:FULL",
		EDID:        "LocWhiterun",
		UpdatedAt:   now,
	})
	if err != nil {
		t.Fatalf("expected create to succeed: %v", err)
	}

	loaded, err := repository.GetByID(context.Background(), created.ID)
	if err != nil {
		t.Fatalf("expected get to succeed: %v", err)
	}
	if loaded.Source != "Whiterun" {
		t.Fatalf("expected created source, got %q", loaded.Source)
	}

	updated, err := repository.Update(context.Background(), created.ID, MasterDictionaryDraft{
		Source:      "Solitude",
		Translation: "ソリチュード",
		Category:    "地名",
		Origin:      "更新",
		REC:         "LCTN:FULL",
		EDID:        "LocSolitude",
		UpdatedAt:   now.Add(time.Minute),
	})
	if err != nil {
		t.Fatalf("expected update to succeed: %v", err)
	}
	if updated.Source != "Solitude" || updated.Origin != "更新" {
		t.Fatal("expected updated entry to reflect latest values")
	}

	listed, err := repository.List(context.Background(), MasterDictionaryListQuery{
		SearchTerm: "Solitude",
		Category:   "地名",
		Page:       1,
		PageSize:   30,
	})
	if err != nil {
		t.Fatalf("expected list to succeed: %v", err)
	}
	if listed.TotalCount != 1 || len(listed.Items) != 1 {
		t.Fatalf("expected one listed item, got total=%d items=%d", listed.TotalCount, len(listed.Items))
	}

	if err := repository.Delete(context.Background(), created.ID); err != nil {
		t.Fatalf("expected delete to succeed: %v", err)
	}
	if _, err := repository.GetByID(context.Background(), created.ID); !errors.Is(err, ErrMasterDictionaryEntryNotFound) {
		t.Fatalf("expected not found after delete, got %v", err)
	}
}

func TestInMemoryMasterDictionaryRepositoryUpsertBySourceAndREC(t *testing.T) {
	now := time.Date(2026, 4, 12, 10, 0, 0, 0, time.UTC)
	repository := NewInMemoryMasterDictionaryRepository([]MasterDictionaryEntry{{
		ID:          1,
		Source:      "Auriel's Bow",
		Translation: "旧訳",
		Category:    "装備",
		Origin:      "初期データ",
		REC:         "WEAP:FULL",
		EDID:        "DLC1AurielsBow",
		UpdatedAt:   now,
	}})

	entry, created, err := repository.UpsertBySourceAndREC(context.Background(), MasterDictionaryImportRecord{
		Source:      "Auriel's Bow",
		Translation: "アーリエルの弓",
		Category:    "装備",
		Origin:      "XML取込",
		REC:         "WEAP:FULL",
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

	entry, created, err = repository.UpsertBySourceAndREC(context.Background(), MasterDictionaryImportRecord{
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
