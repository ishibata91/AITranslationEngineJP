package wails

import (
	"context"
	"errors"
	"testing"
	"time"

	"aitranslationenginejp/internal/usecase"
)

type fakeTranslationInputUsecase struct {
	importXEditJSONFunc   func(ctx context.Context, filePath string) (usecase.TranslationInputImportResult, error)
	rebuildInputCacheFunc func(ctx context.Context, inputID int64) (usecase.TranslationInputImportResult, error)
}

func (fake fakeTranslationInputUsecase) ImportXEditJSON(ctx context.Context, filePath string) (usecase.TranslationInputImportResult, error) {
	if fake.importXEditJSONFunc == nil {
		return usecase.TranslationInputImportResult{}, nil
	}
	return fake.importXEditJSONFunc(ctx, filePath)
}

func (fake fakeTranslationInputUsecase) RebuildInputCache(ctx context.Context, inputID int64) (usecase.TranslationInputImportResult, error) {
	if fake.rebuildInputCacheFunc == nil {
		return usecase.TranslationInputImportResult{}, nil
	}
	return fake.rebuildInputCacheFunc(ctx, inputID)
}

func TestTranslationInputControllerImportPassesFilePathToUsecase(t *testing.T) {
	controller := NewTranslationInputController(fakeTranslationInputUsecase{
		importXEditJSONFunc: func(ctx context.Context, filePath string) (usecase.TranslationInputImportResult, error) {
			if ctx == nil {
				t.Fatal("expected request context")
			}
			if filePath != "/tmp/dialogues.json" {
				t.Fatalf("expected file path to be forwarded, got %q", filePath)
			}
			return usecase.TranslationInputImportResult{}, nil
		},
	})

	_, err := controller.ImportTranslationInput(TranslationInputImportRequestDTO{FilePath: "/tmp/dialogues.json"})
	if err != nil {
		t.Fatalf("expected import request to succeed: %v", err)
	}
}

func TestTranslationInputControllerImportMapsSummaryToDTO(t *testing.T) {
	importedAt := time.Date(2026, 4, 25, 10, 0, 0, 0, time.UTC)
	controller := NewTranslationInputController(fakeTranslationInputUsecase{
		importXEditJSONFunc: func(_ context.Context, _ string) (usecase.TranslationInputImportResult, error) {
			summary := usecase.TranslationInputImportSummary{
				Input: usecase.TranslationInputImportedInput{
					ID:               44,
					SourceFilePath:   "/tmp/dialogues.json",
					SourceTool:       "xEdit",
					TargetPluginName: "Dialogues",
					TargetPluginType: "ESP",
					RecordCount:      2,
					ImportedAt:       importedAt,
				},
				TranslationRecordCount: 2,
				TranslationFieldCount:  2,
				Categories: []usecase.TranslationInputCategoryCount{
					{Category: "DIAL", RecordCount: 1, FieldCount: 1},
					{Category: "INFO", RecordCount: 1, FieldCount: 1},
				},
				SampleFields: []usecase.TranslationInputSampleField{
					{RecordType: "DIAL", SubrecordType: "FULL", FormID: "01000ABC", EditorID: "DialogueGreeting", SourceText: "Hello there", Translatable: true},
				},
				Warnings: []usecase.TranslationInputWarning{
					{Kind: "unknown_field_definition", RecordType: "INFO", SubrecordType: "NAM2", Message: "fallback definition used"},
				},
			}
			return usecase.TranslationInputImportResult{
				Summary:   &summary,
				Warnings:  summary.Warnings,
				ErrorKind: "",
			}, nil
		},
	})

	response, err := controller.ImportTranslationInput(TranslationInputImportRequestDTO{FilePath: "/tmp/dialogues.json"})
	if err != nil {
		t.Fatalf("expected import mapping to succeed: %v", err)
	}
	if !response.Accepted || response.ErrorKind != "" {
		t.Fatalf("expected accepted response, got %#v", response)
	}
	if response.Summary == nil {
		t.Fatal("expected summary DTO on success")
	}
	if response.Summary.Input.ImportedAt != importedAt.Format(time.RFC3339) {
		t.Fatalf("expected RFC3339 importedAt, got %q", response.Summary.Input.ImportedAt)
	}
	if response.Summary.TranslationRecordCount != 2 || response.Summary.TranslationFieldCount != 2 {
		t.Fatalf("unexpected summary counts: %#v", response.Summary)
	}
	if len(response.Summary.Categories) != 2 || response.Summary.Categories[0].Category != "DIAL" {
		t.Fatalf("unexpected categories: %#v", response.Summary.Categories)
	}
	if len(response.Warnings) != 1 || response.Warnings[0].Kind != "unknown_field_definition" {
		t.Fatalf("expected warnings DTOs, got %#v", response.Warnings)
	}
}

func TestTranslationInputControllerImportMapsErrorKindToRejectedResponse(t *testing.T) {
	controller := NewTranslationInputController(fakeTranslationInputUsecase{
		importXEditJSONFunc: func(_ context.Context, _ string) (usecase.TranslationInputImportResult, error) {
			return usecase.TranslationInputImportResult{ErrorKind: "invalid_json"}, nil
		},
	})

	response, err := controller.ImportTranslationInput(TranslationInputImportRequestDTO{FilePath: "/tmp/invalid.json"})
	if err != nil {
		t.Fatalf("expected rejected response without transport error: %v", err)
	}
	if response.Accepted {
		t.Fatalf("expected rejected response, got %#v", response)
	}
	if response.ErrorKind != "invalid_json" {
		t.Fatalf("expected invalid_json error kind, got %q", response.ErrorKind)
	}
	if response.Summary != nil {
		t.Fatalf("expected nil summary on rejection, got %#v", response.Summary)
	}
}

func TestTranslationInputControllerImportWrapsUsecaseError(t *testing.T) {
	boom := errors.New("boom")
	controller := NewTranslationInputController(fakeTranslationInputUsecase{
		importXEditJSONFunc: func(_ context.Context, _ string) (usecase.TranslationInputImportResult, error) {
			return usecase.TranslationInputImportResult{}, boom
		},
	})

	_, err := controller.ImportTranslationInput(TranslationInputImportRequestDTO{FilePath: "/tmp/dialogues.json"})
	if err == nil {
		t.Fatal("expected wrapped error, got nil")
	}
	if !errors.Is(err, boom) {
		t.Fatalf("expected wrapped error to contain original cause, got %v", err)
	}
}

func TestTranslationInputControllerRebuildPassesInputIDToUsecase(t *testing.T) {
	controller := NewTranslationInputController(fakeTranslationInputUsecase{
		rebuildInputCacheFunc: func(ctx context.Context, inputID int64) (usecase.TranslationInputImportResult, error) {
			if ctx == nil {
				t.Fatal("expected request context")
			}
			if inputID != 44 {
				t.Fatalf("expected inputID to be forwarded, got %d", inputID)
			}
			return usecase.TranslationInputImportResult{}, nil
		},
	})

	_, err := controller.RebuildTranslationInputCache(TranslationInputRebuildRequestDTO{InputID: 44})
	if err != nil {
		t.Fatalf("expected rebuild request to succeed: %v", err)
	}
}

func TestTranslationInputControllerRebuildMapsSummaryToDTO(t *testing.T) {
	importedAt := time.Date(2026, 4, 26, 11, 0, 0, 0, time.UTC)
	controller := NewTranslationInputController(fakeTranslationInputUsecase{
		rebuildInputCacheFunc: func(_ context.Context, _ int64) (usecase.TranslationInputImportResult, error) {
			summary := usecase.TranslationInputImportSummary{
				Input: usecase.TranslationInputImportedInput{
					ID:               44,
					SourceFilePath:   "/tmp/dialogues.json",
					SourceTool:       "xEdit",
					TargetPluginName: "Dialogues",
					TargetPluginType: "ESP",
					RecordCount:      2,
					ImportedAt:       importedAt,
				},
				TranslationRecordCount: 2,
				TranslationFieldCount:  2,
				Categories: []usecase.TranslationInputCategoryCount{
					{Category: "DIAL", RecordCount: 1, FieldCount: 1},
					{Category: "INFO", RecordCount: 1, FieldCount: 1},
				},
				Warnings: []usecase.TranslationInputWarning{
					{Kind: "unknown_field_definition", RecordType: "INFO", SubrecordType: "NAM2", Message: "fallback definition used"},
				},
			}
			return usecase.TranslationInputImportResult{Summary: &summary, Warnings: summary.Warnings}, nil
		},
	})

	response, err := controller.RebuildTranslationInputCache(TranslationInputRebuildRequestDTO{InputID: 44})
	if err != nil {
		t.Fatalf("expected rebuild mapping to succeed: %v", err)
	}
	if !response.Accepted || response.ErrorKind != "" {
		t.Fatalf("expected accepted rebuild response, got %#v", response)
	}
	if response.Summary == nil {
		t.Fatal("expected summary DTO on rebuild success")
	}
	if response.Summary.Input.ID != 44 || response.Summary.Input.ImportedAt != importedAt.Format(time.RFC3339) {
		t.Fatalf("unexpected rebuild summary input: %#v", response.Summary.Input)
	}
	if len(response.Warnings) != 1 || response.Warnings[0].SubrecordType != "NAM2" {
		t.Fatalf("expected rebuild warnings DTOs, got %#v", response.Warnings)
	}
}
