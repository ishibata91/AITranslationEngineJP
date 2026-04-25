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
	sqliteRepositoryLocationCategory     = "地名"
	sqliteRepositoryWhiterunTranslation  = "ホワイトラン"
	sqliteRepositoryWindhelmTranslation  = "ウィンドヘルム"
	sqliteRepositoryWeaponCategory       = "装備"
	recLocationFull                      = "LCTN:FULL"
	sqliteRepositoryAurielsBow           = "Auriel's Bow"
	sqliteRepositoryWeaponREC            = "WEAP:FULL"
	errSQLiteRepositoryOpen              = "expected sqlite repository to open: %v"
)

func TestSQLiteMasterDictionaryRepositoryCreateStoresEntry(t *testing.T) {
	repository := newSQLiteMasterDictionaryRepositoryForTest(t, filepath.Join(t.TempDir(), "db", sqliteRepositoryTestDatabaseFileName), nil)
	baseTime := time.Date(2026, 4, 15, 10, 0, 0, 0, time.UTC)

	created, err := repository.Create(context.Background(), MasterDictionaryDraft{
		Source:      "Whiterun",
		Translation: sqliteRepositoryWhiterunTranslation,
		Category:    sqliteRepositoryLocationCategory,
		Origin:      "手動登録",
		REC:         recLocationFull,
		EDID:        "LocWhiterun",
		UpdatedAt:   baseTime,
	})
	if err != nil {
		t.Fatalf("expected create to succeed: %v", err)
	}

	if created.Source != "Whiterun" {
		t.Fatalf("expected created source, got %q", created.Source)
	}
}

func TestSQLiteMasterDictionaryRepositoryUpdateReplacesEntryValues(t *testing.T) {
	repository := newSQLiteMasterDictionaryRepositoryForTest(t, filepath.Join(t.TempDir(), "db", sqliteRepositoryTestDatabaseFileName), nil)
	baseTime := time.Date(2026, 4, 15, 10, 0, 0, 0, time.UTC)
	created := createSQLiteMasterDictionaryEntry(t, repository, MasterDictionaryDraft{
		Source:      "Whiterun",
		Translation: sqliteRepositoryWhiterunTranslation,
		Category:    sqliteRepositoryLocationCategory,
		Origin:      "手動登録",
		REC:         recLocationFull,
		EDID:        "LocWhiterun",
		UpdatedAt:   baseTime,
	})

	updated, err := repository.Update(context.Background(), created.ID, MasterDictionaryDraft{
		Source:      "Whiterun",
		Translation: "更新後ホワイトラン",
		Category:    sqliteRepositoryLocationCategory,
		Origin:      "更新",
		REC:         recLocationFull,
		EDID:        "LocWhiterunUpdated",
		UpdatedAt:   baseTime.Add(time.Minute),
	})
	if err != nil {
		t.Fatalf("expected update to succeed: %v", err)
	}

	if updated.Translation != "更新後ホワイトラン" {
		t.Fatalf("expected updated translation, got %q", updated.Translation)
	}
}

func TestSQLiteMasterDictionaryRepositoryListReturnsFilteredEntry(t *testing.T) {
	repository := newSQLiteMasterDictionaryRepositoryForTest(t, filepath.Join(t.TempDir(), "db", sqliteRepositoryTestDatabaseFileName), nil)
	baseTime := time.Date(2026, 4, 15, 10, 0, 0, 0, time.UTC)
	createSQLiteMasterDictionaryEntry(t, repository, MasterDictionaryDraft{
		Source:      "Whiterun",
		Translation: sqliteRepositoryWhiterunTranslation,
		Category:    sqliteRepositoryLocationCategory,
		Origin:      "手動登録",
		REC:         recLocationFull,
		EDID:        "LocWhiterun",
		UpdatedAt:   baseTime,
	})
	windhelm := createSQLiteMasterDictionaryEntry(t, repository, MasterDictionaryDraft{
		Source:      "Windhelm",
		Translation: sqliteRepositoryWindhelmTranslation,
		Category:    sqliteRepositoryLocationCategory,
		Origin:      "手動登録",
		REC:         recLocationFull,
		EDID:        "LocWindhelm",
		UpdatedAt:   baseTime.Add(time.Minute),
	})

	listed, err := repository.List(context.Background(), MasterDictionaryListQuery{
		SearchTerm: "helm",
		Category:   sqliteRepositoryLocationCategory,
		Page:       1,
		PageSize:   1,
	})
	if err != nil {
		t.Fatalf("expected filtered list to succeed: %v", err)
	}

	if listed.Items[0].ID != windhelm.ID {
		t.Fatalf("expected Windhelm entry to match search, got %#v", listed.Items[0])
	}
}

func TestSQLiteMasterDictionaryRepositoryListOrdersByUpdatedAtDescThenIDDesc(t *testing.T) {
	repository := newSQLiteMasterDictionaryRepositoryForTest(t, filepath.Join(t.TempDir(), "db", sqliteRepositoryTestDatabaseFileName), nil)
	baseTime := time.Date(2026, 4, 15, 10, 0, 0, 0, time.UTC)
	createSQLiteMasterDictionaryEntry(t, repository, MasterDictionaryDraft{
		Source:      "Whiterun",
		Translation: sqliteRepositoryWhiterunTranslation,
		Category:    sqliteRepositoryLocationCategory,
		Origin:      "手動登録",
		REC:         recLocationFull,
		EDID:        "LocWhiterun",
		UpdatedAt:   baseTime,
	})
	windhelm := createSQLiteMasterDictionaryEntry(t, repository, MasterDictionaryDraft{
		Source:      "Windhelm",
		Translation: sqliteRepositoryWindhelmTranslation,
		Category:    sqliteRepositoryLocationCategory,
		Origin:      "手動登録",
		REC:         recLocationFull,
		EDID:        "LocWindhelm",
		UpdatedAt:   baseTime.Add(time.Minute),
	})

	ordered, err := repository.List(context.Background(), MasterDictionaryListQuery{Page: 1, PageSize: 10})
	if err != nil {
		t.Fatalf("expected ordered list to succeed: %v", err)
	}

	if ordered.Items[0].ID != windhelm.ID {
		t.Fatalf("expected larger updated item first, got %#v", ordered.Items)
	}
}

func TestSQLiteMasterDictionaryRepositoryDeleteRemovesEntry(t *testing.T) {
	repository := newSQLiteMasterDictionaryRepositoryForTest(t, filepath.Join(t.TempDir(), "db", sqliteRepositoryTestDatabaseFileName), nil)
	baseTime := time.Date(2026, 4, 15, 10, 0, 0, 0, time.UTC)
	created := createSQLiteMasterDictionaryEntry(t, repository, MasterDictionaryDraft{
		Source:      "Windhelm",
		Translation: sqliteRepositoryWindhelmTranslation,
		Category:    sqliteRepositoryLocationCategory,
		Origin:      "手動登録",
		REC:         recLocationFull,
		EDID:        "LocWindhelm",
		UpdatedAt:   baseTime,
	})

	err := repository.Delete(context.Background(), created.ID)
	if err != nil {
		t.Fatalf("expected delete to succeed: %v", err)
	}

	_, err = repository.GetByID(context.Background(), created.ID)
	if !errors.Is(err, ErrMasterDictionaryEntryNotFound) {
		t.Fatalf("expected deleted entry to be not found, got %v", err)
	}
}

func TestSQLiteMasterDictionaryRepositoryUpsertBySourceAndRECUpdatesExistingRecord(t *testing.T) {
	repository := newSQLiteMasterDictionaryRepositoryForTest(t, filepath.Join(t.TempDir(), "db", sqliteRepositoryTestDatabaseFileName), []MasterDictionaryEntry{{
		ID:          1,
		Source:      sqliteRepositoryAurielsBow,
		Translation: "旧訳",
		Category:    sqliteRepositoryWeaponCategory,
		Origin:      "初期データ",
		REC:         sqliteRepositoryWeaponREC,
		EDID:        "DLC1AurielsBow",
		UpdatedAt:   time.Date(2026, 4, 15, 10, 0, 0, 0, time.UTC),
	}})

	updatedEntry, created, err := repository.UpsertBySourceAndREC(context.Background(), MasterDictionaryImportRecord{
		Source:      "  Auriel's Bow ",
		Translation: "アーリエルの弓",
		Category:    sqliteRepositoryWeaponCategory,
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
}

func TestSQLiteMasterDictionaryRepositoryUpsertBySourceAndRECFallsBackToSourceAndTranslationMatch(t *testing.T) {
	repository := newSQLiteMasterDictionaryRepositoryForTest(t, filepath.Join(t.TempDir(), "db", sqliteRepositoryTestDatabaseFileName), []MasterDictionaryEntry{{
		ID:          1,
		Source:      sqliteRepositoryAurielsBow,
		Translation: "アーリエルの弓",
		Category:    sqliteRepositoryWeaponCategory,
		Origin:      "初期データ",
		REC:         "WEAP:OLD",
		EDID:        "DLC1AurielsBow",
		UpdatedAt:   time.Date(2026, 4, 15, 10, 0, 0, 0, time.UTC),
	}})

	updatedEntry, created, err := repository.UpsertBySourceAndREC(context.Background(), MasterDictionaryImportRecord{
		Source:      "  Auriel's Bow ",
		Translation: " アーリエルの弓 ",
		Category:    sqliteRepositoryWeaponCategory,
		Origin:      "XML取込",
		REC:         sqliteRepositoryWeaponREC,
		EDID:        "DLC1AurielsBow",
		UpdatedAt:   time.Date(2026, 4, 15, 10, 1, 0, 0, time.UTC),
	})
	if err != nil {
		t.Fatalf("expected fallback update upsert to succeed: %v", err)
	}

	if created {
		t.Fatal("expected canonical fallback to update existing record")
	}
	if updatedEntry.ID != 1 {
		t.Fatalf("expected existing id to be reused, got %#v", updatedEntry)
	}
	if updatedEntry.REC != sqliteRepositoryWeaponREC {
		t.Fatalf("expected REC to be updated, got %q", updatedEntry.REC)
	}

	var totalCount int
	if err := repository.database.GetContext(context.Background(), &totalCount,
		"SELECT COUNT(*) FROM DICTIONARY_ENTRY WHERE dictionary_lifecycle = 'master'",
	); err != nil {
		t.Fatalf("expected master dictionary count query to succeed: %v", err)
	}
	if totalCount != 1 {
		t.Fatalf("expected fallback upsert to avoid create, got count %d", totalCount)
	}
}

func TestSQLiteMasterDictionaryRepositoryUpsertBySourceAndRECCreatesNewRecord(t *testing.T) {
	repository := newSQLiteMasterDictionaryRepositoryForTest(t, filepath.Join(t.TempDir(), "db", sqliteRepositoryTestDatabaseFileName), []MasterDictionaryEntry{{
		ID:          1,
		Source:      sqliteRepositoryAurielsBow,
		Translation: "旧訳",
		Category:    sqliteRepositoryWeaponCategory,
		Origin:      "初期データ",
		REC:         sqliteRepositoryWeaponREC,
		EDID:        "DLC1AurielsBow",
		UpdatedAt:   time.Date(2026, 4, 15, 10, 0, 0, 0, time.UTC),
	}})

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
	if createdEntry.ID == 1 {
		t.Fatalf("expected a new id, got %#v", createdEntry)
	}
}

func TestSQLiteMasterDictionaryRepositoryUpsertBySourceAndRECPersistsProvenanceIDOnCreate(t *testing.T) {
	repo := newSQLiteMasterDictionaryRepositoryForTest(t, filepath.Join(t.TempDir(), "db", sqliteRepositoryTestDatabaseFileName), nil)

	// Arrange: FK 制約を満たすため XTRANSLATOR_TRANSLATION_XML 親行を先行挿入する。
	xmlResult, err := repo.database.ExecContext(context.Background(),
		`INSERT INTO XTRANSLATOR_TRANSLATION_XML (file_path, target_plugin_name, target_plugin_type, term_count, imported_at)
		 VALUES (?, ?, ?, ?, ?)`,
		"skyrim_jp.xml", "Skyrim.esm", "ESM", 1, "2026-04-15T10:00:00Z",
	)
	if err != nil {
		t.Fatalf("expected parent XML row insert to succeed: %v", err)
	}
	xmlRowID, err := xmlResult.LastInsertId()
	if err != nil {
		t.Fatalf("expected to read last insert id: %v", err)
	}

	// Act: provenance ID を持つ import record を作成方向で upsert する。
	createdEntry, created, err := repo.UpsertBySourceAndREC(context.Background(), MasterDictionaryImportRecord{
		Source:                      "Iron Sword",
		Translation:                 "アイアンソード",
		Category:                    sqliteRepositoryWeaponCategory,
		Origin:                      "XML取込",
		REC:                         sqliteRepositoryWeaponREC,
		EDID:                        "WeapIronSword",
		UpdatedAt:                   time.Date(2026, 4, 15, 10, 0, 0, 0, time.UTC),
		XTranslatorTranslationXMLID: &xmlRowID,
	})
	if err != nil {
		t.Fatalf("expected upsert create to succeed: %v", err)
	}
	if !created {
		t.Fatal("expected new record to be created")
	}

	// Assert: DB の xtranslator_translation_xml_id が挿入した値と一致する。
	var storedXMLID *int64
	if err := repo.database.GetContext(context.Background(), &storedXMLID,
		"SELECT xtranslator_translation_xml_id FROM DICTIONARY_ENTRY WHERE id = ?", createdEntry.ID,
	); err != nil {
		t.Fatalf("expected to read xtranslator_translation_xml_id from DB: %v", err)
	}
	if storedXMLID == nil || *storedXMLID != xmlRowID {
		t.Fatalf("expected xtranslator_translation_xml_id=%d persisted, got %v", xmlRowID, storedXMLID)
	}
}

func TestSQLiteMasterDictionaryRepositoryUpsertBySourceAndRECPersistsProvenanceIDOnUpdate(t *testing.T) {
	repo := newSQLiteMasterDictionaryRepositoryForTest(t, filepath.Join(t.TempDir(), "db", sqliteRepositoryTestDatabaseFileName), []MasterDictionaryEntry{{
		ID:          1,
		Source:      sqliteRepositoryAurielsBow,
		Translation: "旧訳",
		Category:    sqliteRepositoryWeaponCategory,
		Origin:      "初期データ",
		REC:         sqliteRepositoryWeaponREC,
		EDID:        "DLC1AurielsBow",
		UpdatedAt:   time.Date(2026, 4, 15, 10, 0, 0, 0, time.UTC),
	}})

	// Arrange: FK 制約を満たすため XTRANSLATOR_TRANSLATION_XML 親行を先行挿入する。
	xmlResult, err := repo.database.ExecContext(context.Background(),
		`INSERT INTO XTRANSLATOR_TRANSLATION_XML (file_path, target_plugin_name, target_plugin_type, term_count, imported_at)
		 VALUES (?, ?, ?, ?, ?)`,
		"dawnguard_jp.xml", "Dawnguard.esm", "ESM", 5, "2026-04-15T10:00:00Z",
	)
	if err != nil {
		t.Fatalf("expected parent XML row insert to succeed: %v", err)
	}
	xmlRowID, err := xmlResult.LastInsertId()
	if err != nil {
		t.Fatalf("expected to read last insert id: %v", err)
	}

	// Act: provenance ID を持つ import record を更新方向で upsert する。
	updatedEntry, created, err := repo.UpsertBySourceAndREC(context.Background(), MasterDictionaryImportRecord{
		Source:                      sqliteRepositoryAurielsBow,
		Translation:                 "アーリエルの弓",
		Category:                    sqliteRepositoryWeaponCategory,
		Origin:                      "XML取込",
		REC:                         sqliteRepositoryWeaponREC,
		EDID:                        "DLC1AurielsBow",
		UpdatedAt:                   time.Date(2026, 4, 15, 10, 1, 0, 0, time.UTC),
		XTranslatorTranslationXMLID: &xmlRowID,
	})
	if err != nil {
		t.Fatalf("expected upsert update to succeed: %v", err)
	}
	if created {
		t.Fatal("expected existing record to be updated, not created")
	}

	// Assert: DB の xtranslator_translation_xml_id が挿入した値と一致する。
	var storedXMLID *int64
	if err := repo.database.GetContext(context.Background(), &storedXMLID,
		"SELECT xtranslator_translation_xml_id FROM DICTIONARY_ENTRY WHERE id = ?", updatedEntry.ID,
	); err != nil {
		t.Fatalf("expected to read xtranslator_translation_xml_id from DB: %v", err)
	}
	if storedXMLID == nil || *storedXMLID != xmlRowID {
		t.Fatalf("expected xtranslator_translation_xml_id=%d persisted, got %v", xmlRowID, storedXMLID)
	}
}

func TestSQLiteMasterDictionaryRepositorySeedsEmptyDatabase(t *testing.T) {
	databasePath := filepath.Join(t.TempDir(), "db", sqliteRepositoryTestDatabaseFileName)
	firstSeedTime := time.Date(2026, 4, 15, 9, 0, 0, 0, time.UTC)
	repository := newSQLiteMasterDictionaryRepositoryForTest(t, databasePath, DefaultMasterDictionarySeed(firstSeedTime))

	initialPage, err := repository.List(context.Background(), MasterDictionaryListQuery{Page: 1, PageSize: 100})
	if err != nil {
		t.Fatalf("expected first list to succeed: %v", err)
	}

	if initialPage.TotalCount != 40 {
		t.Fatalf("expected initial seed count 40, got %d", initialPage.TotalCount)
	}
}

func TestSQLiteMasterDictionaryRepositoryDoesNotReseedExistingDatabase(t *testing.T) {
	databasePath := filepath.Join(t.TempDir(), "db", sqliteRepositoryTestDatabaseFileName)
	firstSeedTime := time.Date(2026, 4, 15, 9, 0, 0, 0, time.UTC)
	secondSeedTime := firstSeedTime.Add(24 * time.Hour)
	repository := openSQLiteMasterDictionaryRepositoryWithoutCleanup(t, databasePath, DefaultMasterDictionarySeed(firstSeedTime))
	createdEntry := createSQLiteMasterDictionaryEntry(t, repository, MasterDictionaryDraft{
		Source:      "Forgotten Vale",
		Translation: "忘れられた谷",
		Category:    sqliteRepositoryLocationCategory,
		Origin:      "手動登録",
		REC:         recLocationFull,
		EDID:        "DLC1ForgottenVale",
		UpdatedAt:   firstSeedTime.Add(time.Hour),
	})

	closeSQLiteRepository(t, repository)
	reopenedRepository := newSQLiteMasterDictionaryRepositoryForTest(t, databasePath, DefaultMasterDictionarySeed(secondSeedTime))
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
	repository := openSQLiteMasterDictionaryRepositoryWithoutCleanup(t, databasePath, nil)
	baseTime := time.Date(2026, 4, 15, 11, 0, 0, 0, time.UTC)
	for index := 0; index < 4; index++ {
		createSQLiteMasterDictionaryEntry(t, repository, MasterDictionaryDraft{
			Source:      "Item " + string(rune('A'+index)),
			Translation: "訳",
			Category:    sqliteRepositoryWeaponCategory,
			Origin:      "手動登録",
			REC:         sqliteRepositoryWeaponREC,
			EDID:        "Item",
			UpdatedAt:   baseTime.Add(time.Duration(index) * time.Minute),
		})
	}

	closeSQLiteRepository(t, repository)
	reopenedRepository := newSQLiteMasterDictionaryRepositoryForTest(t, databasePath, nil)
	page, err := reopenedRepository.List(context.Background(), MasterDictionaryListQuery{
		Category: sqliteRepositoryWeaponCategory,
		Page:     2,
		PageSize: 2,
	})
	if err != nil {
		t.Fatalf("expected paged list after reopen to succeed: %v", err)
	}

	if page.TotalCount != 4 || page.Page != 2 || len(page.Items) != 2 {
		t.Fatalf("expected second page with 2 items, got %#v", page)
	}
}

func TestSQLiteMasterDictionaryRepositoryPreservesListOrderingAcrossRestart(t *testing.T) {
	databasePath := filepath.Join(t.TempDir(), "db", sqliteRepositoryTestDatabaseFileName)
	repository := openSQLiteMasterDictionaryRepositoryWithoutCleanup(t, databasePath, nil)
	baseTime := time.Date(2026, 4, 15, 11, 0, 0, 0, time.UTC)
	for index := 0; index < 4; index++ {
		createSQLiteMasterDictionaryEntry(t, repository, MasterDictionaryDraft{
			Source:      "Item " + string(rune('A'+index)),
			Translation: "訳",
			Category:    sqliteRepositoryWeaponCategory,
			Origin:      "手動登録",
			REC:         sqliteRepositoryWeaponREC,
			EDID:        "Item",
			UpdatedAt:   baseTime.Add(time.Duration(index) * time.Minute),
		})
	}

	closeSQLiteRepository(t, repository)
	reopenedRepository := newSQLiteMasterDictionaryRepositoryForTest(t, databasePath, nil)
	page, err := reopenedRepository.List(context.Background(), MasterDictionaryListQuery{
		Category: sqliteRepositoryWeaponCategory,
		Page:     2,
		PageSize: 2,
	})
	if err != nil {
		t.Fatalf("expected paged list after reopen to succeed: %v", err)
	}

	if !page.Items[0].UpdatedAt.After(page.Items[1].UpdatedAt) && page.Items[0].ID <= page.Items[1].ID {
		t.Fatalf("expected page items to remain ordered by updated_at desc/id desc, got %#v", page.Items)
	}
}

func newSQLiteMasterDictionaryRepositoryForTest(t *testing.T, databasePath string, seeds []MasterDictionaryEntry) *SQLiteMasterDictionaryRepository {
	t.Helper()

	repository := openSQLiteMasterDictionaryRepositoryWithoutCleanup(t, databasePath, seeds)
	t.Cleanup(func() {
		closeSQLiteRepository(t, repository)
	})
	return repository
}

func openSQLiteMasterDictionaryRepositoryWithoutCleanup(t *testing.T, databasePath string, seeds []MasterDictionaryEntry) *SQLiteMasterDictionaryRepository {
	t.Helper()

	repository, err := NewSQLiteMasterDictionaryRepository(context.Background(), databasePath, seeds)
	if err != nil {
		t.Fatalf(errSQLiteRepositoryOpen, err)
	}
	return repository
}

func createSQLiteMasterDictionaryEntry(t *testing.T, repository *SQLiteMasterDictionaryRepository, draft MasterDictionaryDraft) MasterDictionaryEntry {
	t.Helper()

	entry, err := repository.Create(context.Background(), draft)
	if err != nil {
		t.Fatalf("expected create to succeed: %v", err)
	}
	return entry
}

func closeSQLiteRepository(t *testing.T, repository *SQLiteMasterDictionaryRepository) {
	t.Helper()
	if err := repository.Close(); err != nil {
		t.Fatalf("expected sqlite repository close to succeed: %v", err)
	}
}
