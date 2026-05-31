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

func TestTranslationInputImportServiceImportXEditJSONWithContentImportsNonDialogueRecords(t *testing.T) {
	repo := &translationInputRepositoryStub{}
	service := NewTranslationInputImportService(repo, translationInputTransactorStub{}, nil, fixedTranslationInputNow)

	summary, err := service.ImportXEditJSONWithContent(
		context.Background(),
		"/imports/non-dialogue.json",
		"non-dialogue.json",
		translationInputNonDialogueFixtureContent,
	)
	if err != nil {
		t.Fatalf("non-dialogue import should succeed: %v", err)
	}

	if summary.Input.TargetPluginName != "NonDialogue" || summary.Input.TargetPluginType != "ESP" {
		t.Fatalf("unexpected imported input metadata: %+v", summary.Input)
	}
	if summary.TranslationRecordCount != len(repo.recordDrafts) {
		t.Fatalf("summary record count = %d, persisted records = %d", summary.TranslationRecordCount, len(repo.recordDrafts))
	}
	if summary.TranslationFieldCount != len(repo.fieldDrafts) {
		t.Fatalf("summary field count = %d, persisted fields = %d", summary.TranslationFieldCount, len(repo.fieldDrafts))
	}

	persistedRECs := persistedTranslationInputRECs(repo)
	for _, rec := range []string{
		"BOOK:FULL",
		"NPC_:FULL",
		"NPC_:SHRT",
		"ARMO:FULL",
		"WEAP:FULL",
		"LCTN:FULL",
		"CELL:FULL",
		"CONT:FULL",
		"MISC:FULL",
		"ALCH:FULL",
		"RACE:FULL",
		"INGR:FULL",
		"SHOU:FULL",
	} {
		if !persistedRECs[rec] {
			t.Fatalf("expected non-dialogue import to persist REC %q, got %#v", rec, persistedRECs)
		}
	}

	persistedSources := persistedTranslationInputSources(repo)
	for _, sourceText := range []string{
		"Arcane Codex",
		"Steel Harness",
		"River Watch",
		"Non Dialogue NPC",
		"Short NPC",
		"Quest Name",
		"Quest stage text",
		"Quest objective text",
		"Message body",
		"Message title",
		"Load screen text",
		"Race name",
	} {
		if !persistedSources[sourceText] {
			t.Fatalf("expected source text %q to be persisted, got %#v", sourceText, persistedSources)
		}
	}
	if persistedSources[""] {
		t.Fatal("empty source text should not be persisted as a translation field")
	}
}

func TestTranslationInputImportServiceImportXEditJSONWithContentDoesNotRequireDialogueGroups(t *testing.T) {
	repo := &translationInputRepositoryStub{}
	service := NewTranslationInputImportService(repo, translationInputTransactorStub{}, nil, fixedTranslationInputNow)

	summary, err := service.ImportXEditJSONWithContent(
		context.Background(),
		"/imports/items-only.json",
		"items-only.json",
		`{
			"target_plugin": "ItemsOnly.esp",
			"items": [
				{
					"id": "01000001",
					"editor_id": "OnlyBook",
					"type": "BOOK FULL",
					"name": "Only Book"
				}
			]
		}`,
	)
	if err != nil {
		t.Fatalf("items-only import should not require dialogue_groups: %v", err)
	}
	if summary.TranslationRecordCount != 1 || summary.TranslationFieldCount != 1 {
		t.Fatalf("unexpected items-only import summary: %+v", summary)
	}
	if len(repo.recordDrafts) != 1 || repo.recordDrafts[0].RecordType != "BOOK" {
		t.Fatalf("expected one BOOK record, got %#v", repo.recordDrafts)
	}
}

func TestTranslationInputImportServiceImportXEditJSONWithContentRejectsEmptyImportableRecords(t *testing.T) {
	repo := &translationInputRepositoryStub{}
	service := NewTranslationInputImportService(repo, translationInputTransactorStub{}, nil, fixedTranslationInputNow)

	_, err := service.ImportXEditJSONWithContent(
		context.Background(),
		"/imports/empty.json",
		"empty.json",
		`{
			"target_plugin": "Empty.esp",
			"dialogue_groups": [],
			"items": [],
			"magic": [],
			"locations": [],
			"cells": [],
			"system": [],
			"messages": [],
			"load_screens": [],
			"npcs": {},
			"quests": []
		}`,
	)
	if err == nil {
		t.Fatal("empty importable records should fail")
	}
	if kind, ok := TranslationInputErrorKindOf(err); !ok || kind != TranslationInputErrorKindUnsupportedExtractShape {
		t.Fatalf("expected unsupported_extract_shape, got kind=%q ok=%v err=%v", kind, ok, err)
	}
	if len(repo.xEditDrafts) != 0 || len(repo.recordDrafts) != 0 || len(repo.fieldDrafts) != 0 {
		t.Fatalf("empty import should not persist drafts, got xedit=%d records=%d fields=%d", len(repo.xEditDrafts), len(repo.recordDrafts), len(repo.fieldDrafts))
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

func persistedTranslationInputRECs(repo *translationInputRepositoryStub) map[string]bool {
	recs := map[string]bool{}
	for _, field := range repo.fieldDrafts {
		recordIndex := int(field.TranslationRecordID) - 1
		if recordIndex < 0 || recordIndex >= len(repo.recordDrafts) {
			continue
		}
		record := repo.recordDrafts[recordIndex]
		recs[record.RecordType+":"+field.SubrecordType] = true
	}
	return recs
}

func persistedTranslationInputSources(repo *translationInputRepositoryStub) map[string]bool {
	sources := map[string]bool{}
	for _, field := range repo.fieldDrafts {
		sources[field.SourceText] = true
	}
	return sources
}

const translationInputNonDialogueFixtureContent = `{
	"target_plugin": "NonDialogue.esp",
	"items": [
		{
			"id": "01000001",
			"editor_id": "BookArcaneCodex",
			"type": "BOOK FULL",
			"name": "Arcane Codex",
			"description": ""
		},
		{
			"id": "01000002",
			"editor_id": "ArmorSteelHarness",
			"type": "ARMO FULL",
			"name": "Steel Harness"
		},
		{
			"id": "01000003",
			"editor_id": "WeaponGlassBlade",
			"type": "WEAP FULL",
			"name": "Glass Blade"
		},
		{
			"id": "01000004",
			"editor_id": "ContainerSupplyChest",
			"type": "CONT FULL",
			"name": "Supply Chest"
		},
		{
			"id": "01000005",
			"editor_id": "MiscDwemerGear",
			"type": "MISC FULL",
			"name": "Dwemer Gear"
		},
		{
			"id": "01000006",
			"editor_id": "PotionMoonSugar",
			"type": "ALCH FULL",
			"name": "Moon Sugar Draught"
		},
		{
			"id": "01000007",
			"editor_id": "IngredientFrostRoot",
			"type": "INGR FULL",
			"name": "Frost Root"
		}
	],
	"magic": [
		{
			"id": "01000008",
			"editor_id": "ShoutRiverCall",
			"type": "SHOU FULL",
			"name": "River Call",
			"description": "River call description"
		}
	],
	"locations": [
		{
			"id": "01000009",
			"editor_id": "LocationRiverWatch",
			"type": "LCTN FULL",
			"name": "River Watch"
		},
		{
			"id": "0100000A",
			"editor_id": "CellRiverWatchInterior",
			"type": "CELL FULL",
			"name": "River Watch Interior"
		}
	],
	"system": [
		{
			"id": "0100000B",
			"editor_id": "RaceRiverFolk",
			"type": "RACE FULL",
			"name": "Race name"
		},
		{
			"id": "0100000C",
			"editor_id": "ShortNPC",
			"type": "NPC_ SHRT",
			"name": "Short NPC"
		}
	],
	"messages": [
		{
			"id": "0100000D",
			"editor_id": "MessageWarning",
			"type": "MESG DESC",
			"text": "Message body",
			"title": "Message title"
		}
	],
	"load_screens": [
		{
			"id": "0100000E",
			"editor_id": "LoadRiverWatch",
			"type": "LSCR DESC",
			"text": "Load screen text"
		}
	],
	"npcs": {
		"0100000F": {
			"id": "0100000F",
			"editor_id": "NonDialogueNPC",
			"type": "NPC_ FULL",
			"name": "Non Dialogue NPC"
		}
	},
	"quests": [
		{
			"id": "01000010",
			"editor_id": "QuestRiverWatch",
			"type": "QUST FULL",
			"name": "Quest Name",
			"stages": [
				{
					"stage_index": 10,
					"log_index": 0,
					"type": "QUST CNAM",
					"parent_id": "01000010",
					"parent_editor_id": "QuestRiverWatch",
					"text": "Quest stage text"
				}
			],
			"objectives": [
				{
					"index": "20",
					"type": "QUST NNAM",
					"parent_id": "01000010",
					"parent_editor_id": "QuestRiverWatch",
					"text": "Quest objective text"
				}
			]
		}
	]
}`

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
