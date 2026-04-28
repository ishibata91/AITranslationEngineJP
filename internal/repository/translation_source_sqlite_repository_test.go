package repository

import (
	"context"
	"path/filepath"
	"reflect"
	"testing"
	"time"
)

const sqliteTranslationSourceTestDatabaseFileName = "translation-source.sqlite3"

func TestTranslationJobSetupSQLiteTranslationSourceRepositoryListXEditExtractedDataReturnsImportedInputsInNewestFirstOrder(t *testing.T) {
	repository, closeRepository := newSQLiteTranslationSourceRepositoryForTest(
		t,
		filepath.Join(t.TempDir(), "db", sqliteTranslationSourceTestDatabaseFileName),
	)
	defer closeRepository()

	sameImportedAt := time.Date(2026, 4, 27, 16, 0, 0, 0, time.UTC)
	olderImported := createSQLiteTranslationSourceInput(t, repository, XEditExtractedDataDraft{
		SourceFilePath:    "/imports/older.json",
		SourceContentHash: "older-hash",
		SourceTool:        "xedit",
		TargetPluginName:  "OlderPlugin.esp",
		TargetPluginType:  "esp",
		RecordCount:       10,
		ImportedAt:        sameImportedAt.Add(-time.Hour),
	})
	firstSameTime := createSQLiteTranslationSourceInput(t, repository, XEditExtractedDataDraft{
		SourceFilePath:    "/imports/first-same.json",
		SourceContentHash: "first-same-hash",
		SourceTool:        "xedit",
		TargetPluginName:  "FirstSamePlugin.esp",
		TargetPluginType:  "esp",
		RecordCount:       20,
		ImportedAt:        sameImportedAt,
	})
	secondSameTime := createSQLiteTranslationSourceInput(t, repository, XEditExtractedDataDraft{
		SourceFilePath:    "/imports/second-same.json",
		SourceContentHash: "second-same-hash",
		SourceTool:        "xedit",
		TargetPluginName:  "SecondSamePlugin.esp",
		TargetPluginType:  "esp",
		RecordCount:       30,
		ImportedAt:        sameImportedAt,
	})

	got, err := repository.ListXEditExtractedData(context.Background())
	if err != nil {
		t.Fatalf("expected imported input list to succeed: %v", err)
	}

	want := []XEditExtractedData{secondSameTime, firstSameTime, olderImported}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("expected imported inputs %#v, got %#v", want, got)
	}
}

func newSQLiteTranslationSourceRepositoryForTest(
	t *testing.T,
	databasePath string,
) (*SQLiteTranslationSourceRepository, func()) {
	t.Helper()

	db, err := OpenSQLiteDatabase(context.Background(), databasePath)
	if err != nil {
		t.Fatalf("expected sqlite translation source repository to open: %v", err)
	}

	repository := &SQLiteTranslationSourceRepository{db: db}
	return repository, func() {
		if closeErr := db.Close(); closeErr != nil {
			t.Fatalf("expected sqlite translation source database close to succeed: %v", closeErr)
		}
	}
}

func createSQLiteTranslationSourceInput(
	t *testing.T,
	repository *SQLiteTranslationSourceRepository,
	draft XEditExtractedDataDraft,
) XEditExtractedData {
	t.Helper()

	created, err := repository.CreateXEditExtractedData(context.Background(), draft)
	if err != nil {
		t.Fatalf("expected x_edit_extracted_data create to succeed: %v", err)
	}
	return created
}
