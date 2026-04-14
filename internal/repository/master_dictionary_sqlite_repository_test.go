package repository

import (
	"context"
	"errors"
	"path/filepath"
	"testing"
	"time"
)

const (
	sqliteRepositoryTestDatabaseFileName = "master-dictionary.sqlite3"
	recLocationFull                      = "LCTN:FULL"
	errSQLiteRepositoryOpen              = "expected sqlite repository to open: %v"
)

func TestSQLiteMasterDictionaryRepositoryCRUDAndList(t *testing.T) {
	databasePath := filepath.Join(t.TempDir(), "db", sqliteRepositoryTestDatabaseFileName)
	repository, err := NewSQLiteMasterDictionaryRepository(context.Background(), databasePath, nil)
	if err != nil {
		t.Fatalf(errSQLiteRepositoryOpen, err)
	}
	t.Cleanup(func() {
		closeSQLiteRepository(t, repository)
	})

	baseTime := time.Date(2026, 4, 15, 10, 0, 0, 0, time.UTC)
	createdWhiterun, err := repository.Create(context.Background(), MasterDictionaryDraft{
		Source:      "Whiterun",
		Translation: "ホワイトラン",
		Category:    "地名",
		Origin:      "手動登録",
		REC:         recLocationFull,
		EDID:        "LocWhiterun",
		UpdatedAt:   baseTime,
	})
	if err != nil {
		t.Fatalf("expected first create to succeed: %v", err)
	}
	createdWindhelm, err := repository.Create(context.Background(), MasterDictionaryDraft{
		Source:      "Windhelm",
		Translation: "ウィンドヘルム",
		Category:    "地名",
		Origin:      "手動登録",
		REC:         recLocationFull,
		EDID:        "LocWindhelm",
		UpdatedAt:   baseTime.Add(time.Minute),
	})
	if err != nil {
		t.Fatalf("expected second create to succeed: %v", err)
	}
	_, err = repository.Update(context.Background(), createdWhiterun.ID, MasterDictionaryDraft{
		Source:      "Whiterun",
		Translation: "更新後ホワイトラン",
		Category:    "地名",
		Origin:      "更新",
		REC:         recLocationFull,
		EDID:        "LocWhiterunUpdated",
		UpdatedAt:   baseTime.Add(time.Minute),
	})
	if err != nil {
		t.Fatalf("expected update to succeed: %v", err)
	}

	listed, err := repository.List(context.Background(), MasterDictionaryListQuery{
		SearchTerm: "helm",
		Category:   "地名",
		Page:       1,
		PageSize:   1,
	})
	if err != nil {
		t.Fatalf("expected filtered list to succeed: %v", err)
	}
	if listed.TotalCount != 1 || len(listed.Items) != 1 {
		t.Fatalf("expected one filtered item, got total=%d items=%d", listed.TotalCount, len(listed.Items))
	}
	if listed.Items[0].ID != createdWindhelm.ID {
		t.Fatalf("expected Windhelm entry to match search, got %#v", listed.Items[0])
	}

	ordered, err := repository.List(context.Background(), MasterDictionaryListQuery{Page: 1, PageSize: 10})
	if err != nil {
		t.Fatalf("expected ordered list to succeed: %v", err)
	}
	if len(ordered.Items) != 2 {
		t.Fatalf("expected 2 items, got %d", len(ordered.Items))
	}
	if ordered.Items[0].ID != createdWindhelm.ID {
		t.Fatalf("expected larger id to sort first when UpdatedAt ties, got %#v", ordered.Items)
	}

	if err := repository.Delete(context.Background(), createdWindhelm.ID); err != nil {
		t.Fatalf("expected delete to succeed: %v", err)
	}
	if _, err := repository.GetByID(context.Background(), createdWindhelm.ID); !errors.Is(err, ErrMasterDictionaryEntryNotFound) {
		t.Fatalf("expected deleted entry to be not found, got %v", err)
	}
}

func TestSQLiteMasterDictionaryRepositoryUpsertBySourceAndREC(t *testing.T) {
	databasePath := filepath.Join(t.TempDir(), "db", sqliteRepositoryTestDatabaseFileName)
	repository, err := NewSQLiteMasterDictionaryRepository(context.Background(), databasePath, []MasterDictionaryEntry{{
		ID:          1,
		Source:      "Auriel's Bow",
		Translation: "旧訳",
		Category:    "装備",
		Origin:      "初期データ",
		REC:         "WEAP:FULL",
		EDID:        "DLC1AurielsBow",
		UpdatedAt:   time.Date(2026, 4, 15, 10, 0, 0, 0, time.UTC),
	}})
	if err != nil {
		t.Fatalf(errSQLiteRepositoryOpen, err)
	}
	t.Cleanup(func() {
		closeSQLiteRepository(t, repository)
	})

	updatedEntry, created, err := repository.UpsertBySourceAndREC(context.Background(), MasterDictionaryImportRecord{
		Source:      "  Auriel's Bow ",
		Translation: "アーリエルの弓",
		Category:    "装備",
		Origin:      "XML取込",
		REC:         " weap:full ",
		EDID:        "DLC1AurielsBow",
		UpdatedAt:   time.Date(2026, 4, 15, 10, 1, 0, 0, time.UTC),
	})
	if err != nil {
		t.Fatalf("expected update upsert to succeed: %v", err)
	}
	if created {
		t.Fatal("expected existing record to be updated")
	}
	if updatedEntry.Translation != "アーリエルの弓" {
		t.Fatalf("expected translation to be updated, got %q", updatedEntry.Translation)
	}

	createdEntry, created, err := repository.UpsertBySourceAndREC(context.Background(), MasterDictionaryImportRecord{
		Source:      "Snow Elf",
		Translation: "スノーエルフ",
		Category:    "NPC",
		Origin:      "XML取込",
		REC:         "NPC_:FULL",
		EDID:        "DLC1SnowElf",
		UpdatedAt:   time.Date(2026, 4, 15, 10, 2, 0, 0, time.UTC),
	})
	if err != nil {
		t.Fatalf("expected create upsert to succeed: %v", err)
	}
	if !created {
		t.Fatal("expected new record to be created")
	}
	if createdEntry.ID == updatedEntry.ID {
		t.Fatalf("expected a new id, got %#v", createdEntry)
	}
}

func TestSQLiteMasterDictionaryRepositorySeedsOnlyEmptyDatabase(t *testing.T) {
	databasePath := filepath.Join(t.TempDir(), "db", sqliteRepositoryTestDatabaseFileName)
	firstSeedTime := time.Date(2026, 4, 15, 9, 0, 0, 0, time.UTC)
	secondSeedTime := firstSeedTime.Add(24 * time.Hour)

	repository, err := NewSQLiteMasterDictionaryRepository(context.Background(), databasePath, DefaultMasterDictionarySeed(firstSeedTime))
	if err != nil {
		t.Fatalf("expected first sqlite repository open to succeed: %v", err)
	}
	initialPage, err := repository.List(context.Background(), MasterDictionaryListQuery{Page: 1, PageSize: 100})
	if err != nil {
		t.Fatalf("expected first list to succeed: %v", err)
	}
	if initialPage.TotalCount != 40 {
		t.Fatalf("expected initial seed count 40, got %d", initialPage.TotalCount)
	}
	createdEntry, err := repository.Create(context.Background(), MasterDictionaryDraft{
		Source:      "Forgotten Vale",
		Translation: "忘れられた谷",
		Category:    "地名",
		Origin:      "手動登録",
		REC:         recLocationFull,
		EDID:        "DLC1ForgottenVale",
		UpdatedAt:   firstSeedTime.Add(time.Hour),
	})
	if err != nil {
		t.Fatalf("expected create before reopen to succeed: %v", err)
	}
	closeErr := repository.Close()
	if closeErr != nil {
		t.Fatalf("expected first sqlite repository close to succeed: %v", closeErr)
	}

	reopenedRepository, err := NewSQLiteMasterDictionaryRepository(context.Background(), databasePath, DefaultMasterDictionarySeed(secondSeedTime))
	if err != nil {
		t.Fatalf("expected reopened sqlite repository open to succeed: %v", err)
	}
	t.Cleanup(func() {
		closeSQLiteRepository(t, reopenedRepository)
	})

	reopenedPage, err := reopenedRepository.List(context.Background(), MasterDictionaryListQuery{Page: 1, PageSize: 100})
	if err != nil {
		t.Fatalf("expected reopened list to succeed: %v", err)
	}
	if reopenedPage.TotalCount != 41 {
		t.Fatalf("expected reopened count 41 without reseed duplication, got %d", reopenedPage.TotalCount)
	}
	loadedEntry, err := reopenedRepository.GetByID(context.Background(), createdEntry.ID)
	if err != nil {
		t.Fatalf("expected created entry to persist after reopen: %v", err)
	}
	if loadedEntry.Source != createdEntry.Source {
		t.Fatalf("expected persisted entry source %q, got %q", createdEntry.Source, loadedEntry.Source)
	}
}

func TestSQLiteMasterDictionaryRepositoryPreservesPaginationAcrossRestart(t *testing.T) {
	databasePath := filepath.Join(t.TempDir(), "db", sqliteRepositoryTestDatabaseFileName)
	repository, err := NewSQLiteMasterDictionaryRepository(context.Background(), databasePath, nil)
	if err != nil {
		t.Fatalf(errSQLiteRepositoryOpen, err)
	}

	baseTime := time.Date(2026, 4, 15, 11, 0, 0, 0, time.UTC)
	for index := 0; index < 4; index++ {
		_, createErr := repository.Create(context.Background(), MasterDictionaryDraft{
			Source:      "Item " + string(rune('A'+index)),
			Translation: "訳",
			Category:    "装備",
			Origin:      "手動登録",
			REC:         "WEAP:FULL",
			EDID:        "Item",
			UpdatedAt:   baseTime.Add(time.Duration(index) * time.Minute),
		})
		if createErr != nil {
			t.Fatalf("expected create %d to succeed: %v", index, createErr)
		}
	}
	closeErr := repository.Close()
	if closeErr != nil {
		t.Fatalf("expected sqlite repository close to succeed: %v", closeErr)
	}

	reopenedRepository, err := NewSQLiteMasterDictionaryRepository(context.Background(), databasePath, nil)
	if err != nil {
		t.Fatalf("expected reopened sqlite repository open to succeed: %v", err)
	}
	t.Cleanup(func() {
		closeSQLiteRepository(t, reopenedRepository)
	})

	page, err := reopenedRepository.List(context.Background(), MasterDictionaryListQuery{
		Category: "装備",
		Page:     2,
		PageSize: 2,
	})
	if err != nil {
		t.Fatalf("expected paged list after reopen to succeed: %v", err)
	}
	if page.TotalCount != 4 || page.Page != 2 || len(page.Items) != 2 {
		t.Fatalf("expected second page with 2 items, got %#v", page)
	}
	if !page.Items[0].UpdatedAt.After(page.Items[1].UpdatedAt) && page.Items[0].ID <= page.Items[1].ID {
		t.Fatalf("expected page items to remain ordered by updated_at desc/id desc, got %#v", page.Items)
	}
}

func closeSQLiteRepository(t *testing.T, repository *SQLiteMasterDictionaryRepository) {
	t.Helper()
	if err := repository.Close(); err != nil {
		t.Fatalf("expected sqlite repository close to succeed: %v", err)
	}
}
