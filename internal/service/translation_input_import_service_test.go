package service

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"testing"
	"time"

	"aitranslationenginejp/internal/repository"
)

const (
	translationInputBareFilenameFixture = "Lucien.esp_Export.json"
	translationInputFixtureContent      = `{
		"target_plugin": "Lucien.esp",
		"dialogue_groups": [
			{
				"id": "01000ABC",
				"editor_id": "LucienGreeting",
				"type": "DIAL FULL",
				"player_text": "Hello",
				"responses": [
					{
						"id": "01000ABD",
						"editor_id": "LucienGreetingResponse",
						"type": "INFO NAM1",
						"text": "Need something?",
						"order": 0
					}
				]
			}
		]
	}`
)

func TestReadTranslationInputFileResolvesBareFilenameFromDictionaries(t *testing.T) {
	workingDir := writeTranslationInputFixture(t)
	changeWorkingDirectory(t, workingDir)

	validatedPath, err := validateTranslationInputPath(translationInputBareFilenameFixture)
	if err != nil {
		t.Fatalf("expected bare filename to validate: %v", err)
	}

	content, err := readTranslationInputFile(validatedPath)
	if err != nil {
		t.Fatalf("expected dictionaries fallback to resolve bare filename: %v", err)
	}

	if string(content) != translationInputFixtureContent {
		t.Fatalf("expected fallback content to match fixture, got %q", string(content))
	}
}

func TestTranslationInputImportServiceImportXEditJSONAcceptsBareFilenameFromDictionaries(t *testing.T) {
	workingDir := writeTranslationInputFixture(t)
	changeWorkingDirectory(t, workingDir)

	repo := &translationInputRepositoryStub{}
	service := NewTranslationInputImportService(repo, translationInputTransactorStub{}, nil, fixedTranslationInputNow)

	summary, err := service.ImportXEditJSON(context.Background(), translationInputBareFilenameFixture)
	if err != nil {
		if kind, ok := TranslationInputErrorKindOf(err); ok && kind == TranslationInputErrorKindSourceFileMissing {
			t.Fatalf("expected bare filename import to pass read stage, got %v", err)
		}
		t.Fatalf("expected import to succeed after dictionaries fallback: %v", err)
	}

	if summary.Input.SourceFilePath != translationInputBareFilenameFixture {
		t.Fatalf("expected import summary to retain bare filename, got %q", summary.Input.SourceFilePath)
	}
	if summary.Input.TargetPluginName != "Lucien" || summary.Input.TargetPluginType != "ESP" {
		t.Fatalf("unexpected imported input metadata: %+v", summary.Input)
	}
	if summary.TranslationRecordCount != 2 || summary.TranslationFieldCount != 2 {
		t.Fatalf("expected import summary to persist decoded records, got %+v", summary)
	}
	if len(repo.xEditDrafts) != 1 || len(repo.recordDrafts) != 2 || len(repo.fieldDrafts) != 2 {
		t.Fatalf("expected persistence after read stage, got xedit=%d records=%d fields=%d", len(repo.xEditDrafts), len(repo.recordDrafts), len(repo.fieldDrafts))
	}
	if repo.xEditDrafts[0].SourceFilePath != translationInputBareFilenameFixture {
		t.Fatalf("expected xEdit draft to keep bare filename, got %+v", repo.xEditDrafts[0])
	}
}

func TestTranslationInputImportServiceImportXEditJSONWithContentAllowsDuplicateSourceContentHash(t *testing.T) {
	databasePath := filepath.Join(t.TempDir(), "db", "translation-input.sqlite3")
	db, err := repository.OpenSQLiteDatabase(context.Background(), databasePath)
	if err != nil {
		t.Fatalf("open sqlite database: %v", err)
	}
	t.Cleanup(func() {
		if closeErr := db.Close(); closeErr != nil {
			t.Fatalf("close sqlite database: %v", closeErr)
		}
	})

	sourceRepo := repository.NewSQLiteTranslationSourceRepository(db)
	sqliteSourceRepo, ok := sourceRepo.(*repository.SQLiteTranslationSourceRepository)
	if !ok {
		t.Fatal("expected SQLite translation source repository concrete type")
	}
	service := NewTranslationInputImportService(
		sourceRepo,
		repository.NewSQLiteTransactor(db),
		nil,
		fixedTranslationInputNow,
	)

	firstSummary, err := service.ImportXEditJSONWithContent(
		context.Background(),
		"/imports/duplicate-a.json",
		"duplicate-a.json",
		translationInputFixtureContent,
	)
	if err != nil {
		t.Fatalf("first import should succeed: %v", err)
	}

	secondSummary, err := service.ImportXEditJSONWithContent(
		context.Background(),
		"/imports/duplicate-b.json",
		"duplicate-b.json",
		translationInputFixtureContent,
	)
	if err != nil {
		t.Fatalf("second import with same content should succeed: %v", err)
	}

	if firstSummary.Input.ID == secondSummary.Input.ID {
		t.Fatalf("expected duplicate content imports to create distinct input ids, got %d", firstSummary.Input.ID)
	}
	if firstSummary.Input.SourceFilePath == secondSummary.Input.SourceFilePath {
		t.Fatalf("expected each import to retain its own source path, got %q", firstSummary.Input.SourceFilePath)
	}

	importedInputs, err := sqliteSourceRepo.ListXEditExtractedData(context.Background())
	if err != nil {
		t.Fatalf("list imported inputs: %v", err)
	}
	if len(importedInputs) != 2 {
		t.Fatalf("expected two imported inputs, got %d", len(importedInputs))
	}
	if importedInputs[0].SourceContentHash != importedInputs[1].SourceContentHash {
		t.Fatalf("expected duplicate content imports to share source_content_hash, got %q and %q", importedInputs[0].SourceContentHash, importedInputs[1].SourceContentHash)
	}
}

func writeTranslationInputFixture(t *testing.T) string {
	t.Helper()

	rootDir := t.TempDir()
	dictionariesDir := filepath.Join(rootDir, "dictionaries")
	if err := os.MkdirAll(dictionariesDir, 0o750); err != nil {
		t.Fatalf("create dictionaries dir: %v", err)
	}
	if err := os.WriteFile(filepath.Join(dictionariesDir, translationInputBareFilenameFixture), []byte(translationInputFixtureContent), 0o600); err != nil {
		t.Fatalf("write translation input fixture: %v", err)
	}

	workingDir := filepath.Join(rootDir, "internal", "service")
	if err := os.MkdirAll(workingDir, 0o750); err != nil {
		t.Fatalf("create working dir: %v", err)
	}

	return workingDir
}

func changeWorkingDirectory(t *testing.T, dir string) {
	t.Helper()

	originalDir, err := os.Getwd()
	if err != nil {
		t.Fatalf("get working directory: %v", err)
	}
	if err := os.Chdir(dir); err != nil {
		t.Fatalf("change working directory: %v", err)
	}
	t.Cleanup(func() {
		if err := os.Chdir(originalDir); err != nil {
			panic(fmt.Sprintf("restore working directory: %v", err))
		}
	})
}

func fixedTranslationInputNow() time.Time {
	return time.Date(2026, 4, 26, 12, 0, 0, 0, time.UTC)
}

type translationInputTransactorStub struct{}

func (translationInputTransactorStub) WithTransaction(ctx context.Context, fn func(ctx context.Context) error) error {
	return fn(ctx)
}

type translationInputRepositoryStub struct {
	xEditDrafts  []repository.XEditExtractedDataDraft
	recordDrafts []repository.TranslationRecordDraft
	fieldDrafts  []repository.TranslationFieldDraft
}

func (stub *translationInputRepositoryStub) CreateXEditExtractedData(_ context.Context, draft repository.XEditExtractedDataDraft) (repository.XEditExtractedData, error) {
	stub.xEditDrafts = append(stub.xEditDrafts, draft)
	return repository.XEditExtractedData{
		ID:                int64(len(stub.xEditDrafts)),
		SourceFilePath:    draft.SourceFilePath,
		SourceContentHash: draft.SourceContentHash,
		SourceTool:        draft.SourceTool,
		TargetPluginName:  draft.TargetPluginName,
		TargetPluginType:  draft.TargetPluginType,
		RecordCount:       draft.RecordCount,
		ImportedAt:        draft.ImportedAt,
	}, nil
}

func (stub *translationInputRepositoryStub) GetXEditExtractedDataByID(context.Context, int64) (repository.XEditExtractedData, error) {
	panic("unexpected GetXEditExtractedDataByID call")
}

func (stub *translationInputRepositoryStub) DeleteXEditExtractedDataByID(context.Context, int64) error {
	panic("unexpected DeleteXEditExtractedDataByID call")
}

func (stub *translationInputRepositoryStub) CreateTranslationRecord(_ context.Context, draft repository.TranslationRecordDraft) (repository.TranslationRecord, error) {
	stub.recordDrafts = append(stub.recordDrafts, draft)
	return repository.TranslationRecord{
		ID:                   int64(len(stub.recordDrafts)),
		XEditExtractedDataID: draft.XEditExtractedDataID,
		FormID:               draft.FormID,
		EditorID:             draft.EditorID,
		RecordType:           draft.RecordType,
	}, nil
}

func TestLogTranslationInputImportBulkSummaryUsesAggregateSafePayload(t *testing.T) {
	var buffer bytes.Buffer
	previous := slog.Default()
	slog.SetDefault(slog.New(slog.NewJSONHandler(&buffer, nil)))
	t.Cleanup(func() { slog.SetDefault(previous) })

	logTranslationInputImportBulkSummary(context.Background(), "import", TranslationInputImportSummary{
		TranslationRecordCount: 2,
		TranslationFieldCount:  5,
		Warnings: []TranslationInputWarning{
			{Kind: TranslationInputWarningKindUnknownFieldDefinition},
		},
	})

	var payload map[string]any
	if err := json.Unmarshal(buffer.Bytes(), &payload); err != nil {
		t.Fatalf("unmarshal translation input payload: %v", err)
	}
	if payload["event"] != "translation_input_import_bulk_summary" || payload["result"] != "completed" {
		t.Fatalf("unexpected translation input payload: %#v", payload)
	}
	if payload["input_count"] != float64(2) || payload["output_count"] != float64(5) || payload["skipped_count"] != float64(1) || payload["failed_count"] != float64(0) {
		t.Fatalf("unexpected aggregate counts: %#v", payload)
	}
	forbidden := []string{"api_key", "endpoint", "raw_request", "raw_response", "full_path", "trace_id", "dto"}
	for _, key := range forbidden {
		if _, ok := payload[key]; ok {
			t.Fatalf("forbidden key %q in payload: %#v", key, payload)
		}
	}
}

func (stub *translationInputRepositoryStub) GetTranslationRecordByID(context.Context, int64) (repository.TranslationRecord, error) {
	panic("unexpected GetTranslationRecordByID call")
}

func (stub *translationInputRepositoryStub) ListTranslationRecordsByXEditID(context.Context, int64) ([]repository.TranslationRecord, error) {
	panic("unexpected ListTranslationRecordsByXEditID call")
}

func (stub *translationInputRepositoryStub) UpsertNpcProfile(context.Context, repository.NpcProfileDraft) (repository.NpcProfile, error) {
	panic("unexpected UpsertNpcProfile call")
}

func (stub *translationInputRepositoryStub) GetNpcProfileByID(context.Context, int64) (repository.NpcProfile, error) {
	panic("unexpected GetNpcProfileByID call")
}

func (stub *translationInputRepositoryStub) CreateNpcRecord(context.Context, repository.NpcRecordDraft) (repository.NpcRecord, error) {
	panic("unexpected CreateNpcRecord call")
}

func (stub *translationInputRepositoryStub) GetNpcRecordByTranslationRecordID(context.Context, int64) (repository.NpcRecord, error) {
	panic("unexpected GetNpcRecordByTranslationRecordID call")
}

func (stub *translationInputRepositoryStub) CreateTranslationField(_ context.Context, draft repository.TranslationFieldDraft) (repository.TranslationField, error) {
	stub.fieldDrafts = append(stub.fieldDrafts, draft)
	return repository.TranslationField{
		ID:                           int64(len(stub.fieldDrafts)),
		TranslationRecordID:          draft.TranslationRecordID,
		TranslationFieldDefinitionID: draft.TranslationFieldDefinitionID,
		SubrecordType:                draft.SubrecordType,
		SourceText:                   draft.SourceText,
		FieldOrder:                   draft.FieldOrder,
		PreviousTranslationFieldID:   draft.PreviousTranslationFieldID,
		NextTranslationFieldID:       draft.NextTranslationFieldID,
	}, nil
}

func (stub *translationInputRepositoryStub) GetTranslationFieldByID(context.Context, int64) (repository.TranslationField, error) {
	panic("unexpected GetTranslationFieldByID call")
}

func (stub *translationInputRepositoryStub) ListTranslationFieldsByTranslationRecordID(context.Context, int64) ([]repository.TranslationField, error) {
	panic("unexpected ListTranslationFieldsByTranslationRecordID call")
}

func (stub *translationInputRepositoryStub) CreateTranslationFieldRecordReference(context.Context, repository.TranslationFieldRecordReferenceDraft) (repository.TranslationFieldRecordReference, error) {
	panic("unexpected CreateTranslationFieldRecordReference call")
}

func (stub *translationInputRepositoryStub) ListTranslationFieldRecordReferencesByFieldID(context.Context, int64) ([]repository.TranslationFieldRecordReference, error) {
	panic("unexpected ListTranslationFieldRecordReferencesByFieldID call")
}
