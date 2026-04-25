package usecase

import (
	"context"
	"errors"
	"os"
	"path/filepath"
	"testing"
	"time"

	"aitranslationenginejp/internal/service"
)

type fakeTranslationInputImportService struct {
	importXEditJSONFunc   func(ctx context.Context, filePath string) (service.TranslationInputImportSummary, error)
	rebuildInputCacheFunc func(ctx context.Context, inputID int64) (service.TranslationInputImportSummary, error)
}

func (fake fakeTranslationInputImportService) ImportXEditJSON(ctx context.Context, filePath string) (service.TranslationInputImportSummary, error) {
	if fake.importXEditJSONFunc == nil {
		return service.TranslationInputImportSummary{}, nil
	}
	return fake.importXEditJSONFunc(ctx, filePath)
}

func (fake fakeTranslationInputImportService) RebuildInputCache(ctx context.Context, inputID int64) (service.TranslationInputImportSummary, error) {
	if fake.rebuildInputCacheFunc == nil {
		return service.TranslationInputImportSummary{}, nil
	}
	return fake.rebuildInputCacheFunc(ctx, inputID)
}

func TestTranslationInputUsecaseImportXEditJSONReturnsSummaryAndWarnings(t *testing.T) {
	ctx := context.Background()
	importedAt := time.Date(2026, 4, 25, 9, 30, 0, 0, time.UTC)
	usecase := NewTranslationInputUsecase(fakeTranslationInputImportService{
		importXEditJSONFunc: func(callCtx context.Context, filePath string) (service.TranslationInputImportSummary, error) {
			if callCtx != ctx {
				t.Fatal("expected import service to receive original context")
			}
			if filePath != "/tmp/dialogues.json" {
				t.Fatalf("expected file path to be forwarded, got %q", filePath)
			}
			return service.TranslationInputImportSummary{
				Input: service.TranslationInputImportedInput{
					ID:               12,
					SourceFilePath:   filePath,
					SourceTool:       "xEdit",
					TargetPluginName: "Dialogues",
					TargetPluginType: "ESP",
					RecordCount:      2,
					ImportedAt:       importedAt,
				},
				TranslationRecordCount: 2,
				TranslationFieldCount:  2,
				Categories: []service.TranslationInputCategoryCount{
					{Category: "DIAL", RecordCount: 1, FieldCount: 1},
					{Category: "INFO", RecordCount: 1, FieldCount: 1},
				},
				SampleFields: []service.TranslationInputSampleField{
					{RecordType: "DIAL", SubrecordType: "FULL", FormID: "01000ABC", EditorID: "DialogueGreeting", SourceText: "Hello there", Translatable: true},
				},
				Warnings: []service.TranslationInputWarning{
					{Kind: service.TranslationInputWarningKindUnknownFieldDefinition, RecordType: "INFO", SubrecordType: "NAM2", Message: "fallback definition used"},
				},
			}, nil
		},
	})

	result, err := usecase.ImportXEditJSON(ctx, "/tmp/dialogues.json")
	if err != nil {
		t.Fatalf("expected import to succeed: %v", err)
	}
	if result.ErrorKind != "" {
		t.Fatalf("expected empty error kind on success, got %q", result.ErrorKind)
	}
	if result.Summary == nil {
		t.Fatal("expected summary on success")
	}
	if result.Summary.TranslationRecordCount != 2 || result.Summary.TranslationFieldCount != 2 {
		t.Fatalf("unexpected summary counts: %#v", result.Summary)
	}
	if len(result.Warnings) != 1 || result.Warnings[0].Kind != service.TranslationInputWarningKindUnknownFieldDefinition {
		t.Fatalf("expected warnings to be forwarded, got %#v", result.Warnings)
	}
}

func TestTranslationInputUsecaseImportXEditJSONMapsValidationErrorsToErrorKind(t *testing.T) {
	fixtureDir := t.TempDir()
	runImportValidationErrorCase(
		t,
		filepath.Join(fixtureDir, "invalid.json"),
		"{",
		service.TranslationInputErrorKindInvalidJSON,
	)
}

func TestTranslationInputUsecaseImportXEditJSONMapsUnsupportedExtractShapeToErrorKind(t *testing.T) {
	fixtureDir := t.TempDir()
	runImportValidationErrorCase(
		t,
		filepath.Join(fixtureDir, "unsupported.json"),
		`{"target_plugin":"Dialogues.esp","dialogue_groups":[]}`,
		service.TranslationInputErrorKindUnsupportedExtractShape,
	)
}

func TestTranslationInputUsecaseImportXEditJSONMapsMissingRequiredFieldToErrorKind(t *testing.T) {
	assertImportValidationErrorKind(t, "   ", service.TranslationInputErrorKindMissingRequiredField)
}

func TestTranslationInputUsecaseImportXEditJSONWrapsUnexpectedError(t *testing.T) {
	boom := errors.New("boom")
	usecase := NewTranslationInputUsecase(fakeTranslationInputImportService{
		importXEditJSONFunc: func(_ context.Context, _ string) (service.TranslationInputImportSummary, error) {
			return service.TranslationInputImportSummary{}, boom
		},
	})

	_, err := usecase.ImportXEditJSON(context.Background(), "/tmp/dialogues.json")
	if err == nil {
		t.Fatal("expected wrapped error, got nil")
	}
	if !errors.Is(err, boom) {
		t.Fatalf("expected wrapped error to contain original cause, got %v", err)
	}
}

func runImportValidationErrorCase(t *testing.T, filePath string, content string, expected string) {
	t.Helper()
	if err := os.WriteFile(filePath, []byte(content), 0o600); err != nil {
		t.Fatalf("write fixture: %v", err)
	}
	assertImportValidationErrorKind(t, filePath, expected)
}

func assertImportValidationErrorKind(t *testing.T, filePath string, expected string) {
	t.Helper()
	ctx := context.Background()
	actualService := service.NewTranslationInputImportService(nil, nil, nil, nil)
	usecase := NewTranslationInputUsecase(fakeTranslationInputImportService{
		importXEditJSONFunc: func(callCtx context.Context, targetFilePath string) (service.TranslationInputImportSummary, error) {
			return actualService.ImportXEditJSON(callCtx, targetFilePath)
		},
	})

	result, err := usecase.ImportXEditJSON(ctx, filePath)
	if err != nil {
		t.Fatalf("expected validation failure to be mapped, got error: %v", err)
	}
	if result.Summary != nil {
		t.Fatalf("expected no summary on validation failure, got %#v", result.Summary)
	}
	if result.ErrorKind != expected {
		t.Fatalf("expected error kind %q, got %q", expected, result.ErrorKind)
	}
}

func TestTranslationInputUsecaseRebuildInputCacheReturnsSummaryAndWarnings(t *testing.T) {
	ctx := context.Background()
	importedAt := time.Date(2026, 4, 26, 10, 0, 0, 0, time.UTC)
	usecase := NewTranslationInputUsecase(fakeTranslationInputImportService{
		rebuildInputCacheFunc: func(callCtx context.Context, inputID int64) (service.TranslationInputImportSummary, error) {
			if callCtx != ctx {
				t.Fatal("expected rebuild service to receive original context")
			}
			if inputID != 44 {
				t.Fatalf("expected inputID to be forwarded, got %d", inputID)
			}
			return service.TranslationInputImportSummary{
				Input: service.TranslationInputImportedInput{
					ID:               inputID,
					SourceFilePath:   "/tmp/dialogues.json",
					SourceTool:       "xEdit",
					TargetPluginName: "Dialogues",
					TargetPluginType: "ESP",
					RecordCount:      2,
					ImportedAt:       importedAt,
				},
				TranslationRecordCount: 2,
				TranslationFieldCount:  2,
				Categories: []service.TranslationInputCategoryCount{
					{Category: "DIAL", RecordCount: 1, FieldCount: 1},
					{Category: "INFO", RecordCount: 1, FieldCount: 1},
				},
				Warnings: []service.TranslationInputWarning{
					{Kind: service.TranslationInputWarningKindUnknownFieldDefinition, RecordType: "INFO", SubrecordType: "NAM2", Message: "fallback definition used"},
				},
			}, nil
		},
	})

	result, err := usecase.RebuildInputCache(ctx, 44)
	if err != nil {
		t.Fatalf("expected rebuild to succeed: %v", err)
	}
	if result.ErrorKind != "" {
		t.Fatalf("expected empty error kind on success, got %q", result.ErrorKind)
	}
	if result.Summary == nil {
		t.Fatal("expected summary on success")
	}
	if result.Summary.Input.ID != 44 {
		t.Fatalf("expected summary to keep input ID, got %#v", result.Summary.Input)
	}
	if len(result.Warnings) != 1 || result.Warnings[0].Kind != service.TranslationInputWarningKindUnknownFieldDefinition {
		t.Fatalf("expected warnings to be forwarded, got %#v", result.Warnings)
	}
}

func TestTranslationInputUsecaseRebuildInputCacheMapsValidationErrorsToErrorKind(t *testing.T) {
	ctx := context.Background()
	actualService := service.NewTranslationInputImportService(nil, nil, nil, nil)
	usecase := NewTranslationInputUsecase(fakeTranslationInputImportService{
		rebuildInputCacheFunc: func(callCtx context.Context, inputID int64) (service.TranslationInputImportSummary, error) {
			return actualService.RebuildInputCache(callCtx, inputID)
		},
	})

	result, err := usecase.RebuildInputCache(ctx, 0)
	if err != nil {
		t.Fatalf("expected validation failure to be mapped, got error: %v", err)
	}
	if result.Summary != nil {
		t.Fatalf("expected no summary on validation failure, got %#v", result.Summary)
	}
	if result.ErrorKind != service.TranslationInputErrorKindMissingRequiredField {
		t.Fatalf("expected missing_required_field, got %q", result.ErrorKind)
	}
}

func TestTranslationInputUsecaseRebuildInputCacheWrapsUnexpectedError(t *testing.T) {
	boom := errors.New("boom")
	usecase := NewTranslationInputUsecase(fakeTranslationInputImportService{
		rebuildInputCacheFunc: func(_ context.Context, _ int64) (service.TranslationInputImportSummary, error) {
			return service.TranslationInputImportSummary{}, boom
		},
	})

	_, err := usecase.RebuildInputCache(context.Background(), 1)
	if err == nil {
		t.Fatal("expected wrapped error when rebuild returns unexpected failure")
	}
	if !errors.Is(err, boom) {
		t.Fatalf("expected wrapped error to contain original cause, got %v", err)
	}
}
