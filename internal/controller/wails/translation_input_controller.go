package wails

import (
	"context"
	"fmt"
	"time"

	"aitranslationenginejp/internal/usecase"
)

// TranslationInputUsecasePort defines the translation input import usecase contract.
type TranslationInputUsecasePort interface {
	ImportXEditJSON(ctx context.Context, filePath string) (usecase.TranslationInputImportResult, error)
}

type translationInputCacheRebuildUsecasePort interface {
	RebuildInputCache(ctx context.Context, inputID int64) (usecase.TranslationInputImportResult, error)
}

// TranslationInputController exposes Wails-bound translation input entrypoints.
type TranslationInputController struct {
	translationInputUsecase TranslationInputUsecasePort
}

// TranslationInputImportRequestDTO requests one xEdit JSON import.
type TranslationInputImportRequestDTO struct {
	FilePath string `json:"filePath"`
}

// TranslationInputRebuildRequestDTO requests one cache rebuild for an imported input.
type TranslationInputRebuildRequestDTO struct {
	InputID int64 `json:"inputId"`
}

// TranslationInputImportedInputDTO is the transport DTO for imported input metadata.
type TranslationInputImportedInputDTO struct {
	ID               int64  `json:"id"`
	SourceFilePath   string `json:"sourceFilePath"`
	SourceTool       string `json:"sourceTool"`
	TargetPluginName string `json:"targetPluginName"`
	TargetPluginType string `json:"targetPluginType"`
	RecordCount      int    `json:"recordCount"`
	ImportedAt       string `json:"importedAt"`
}

// TranslationInputCategoryCountDTO is the transport DTO for category aggregates.
type TranslationInputCategoryCountDTO struct {
	Category    string `json:"category"`
	RecordCount int    `json:"recordCount"`
	FieldCount  int    `json:"fieldCount"`
}

// TranslationInputSampleFieldDTO is the transport DTO for representative imported fields.
type TranslationInputSampleFieldDTO struct {
	RecordType    string `json:"recordType"`
	SubrecordType string `json:"subrecordType"`
	FormID        string `json:"formId"`
	EditorID      string `json:"editorId"`
	SourceText    string `json:"sourceText"`
	Translatable  bool   `json:"translatable"`
}

// TranslationInputWarningDTO is the transport DTO for non-fatal warnings.
type TranslationInputWarningDTO struct {
	Kind          string `json:"kind"`
	RecordType    string `json:"recordType"`
	SubrecordType string `json:"subrecordType"`
	Message       string `json:"message"`
}

// TranslationInputImportSummaryDTO is the transport DTO for a successful import summary.
type TranslationInputImportSummaryDTO struct {
	Input                  TranslationInputImportedInputDTO   `json:"input"`
	TranslationRecordCount int                                `json:"translationRecordCount"`
	TranslationFieldCount  int                                `json:"translationFieldCount"`
	Categories             []TranslationInputCategoryCountDTO `json:"categories"`
	SampleFields           []TranslationInputSampleFieldDTO   `json:"sampleFields"`
	Warnings               []TranslationInputWarningDTO       `json:"warnings"`
}

// TranslationInputImportResponseDTO is the Wails response for one import request.
type TranslationInputImportResponseDTO struct {
	Accepted  bool                              `json:"accepted"`
	Summary   *TranslationInputImportSummaryDTO `json:"summary,omitempty"`
	ErrorKind string                            `json:"errorKind,omitempty"`
	Warnings  []TranslationInputWarningDTO      `json:"warnings,omitempty"`
}

// NewTranslationInputController creates a translation input controller.
func NewTranslationInputController(usecase TranslationInputUsecasePort) *TranslationInputController {
	return &TranslationInputController{translationInputUsecase: usecase}
}

// ImportTranslationInput imports one xEdit JSON file through the usecase boundary.
func (controller *TranslationInputController) ImportTranslationInput(
	request TranslationInputImportRequestDTO,
) (TranslationInputImportResponseDTO, error) {
	result, err := controller.translationInputUsecase.ImportXEditJSON(context.Background(), request.FilePath)
	if err != nil {
		return TranslationInputImportResponseDTO{}, fmt.Errorf("import translation input: %w", err)
	}
	return toTranslationInputImportResponseDTO(result), nil
}

// RebuildTranslationInputCache rebuilds the persisted input cache from the canonical source JSON.
func (controller *TranslationInputController) RebuildTranslationInputCache(
	request TranslationInputRebuildRequestDTO,
) (TranslationInputImportResponseDTO, error) {
	rebuildUsecase, ok := controller.translationInputUsecase.(translationInputCacheRebuildUsecasePort)
	if !ok {
		return TranslationInputImportResponseDTO{}, fmt.Errorf("rebuild translation input cache: usecase does not support rebuild")
	}
	result, err := rebuildUsecase.RebuildInputCache(context.Background(), request.InputID)
	if err != nil {
		return TranslationInputImportResponseDTO{}, fmt.Errorf("rebuild translation input cache: %w", err)
	}
	return toTranslationInputImportResponseDTO(result), nil
}

func toTranslationInputImportResponseDTO(result usecase.TranslationInputImportResult) TranslationInputImportResponseDTO {
	response := TranslationInputImportResponseDTO{
		Accepted:  result.ErrorKind == "",
		ErrorKind: result.ErrorKind,
		Warnings:  toTranslationInputWarningDTOs(result.Warnings),
	}
	if result.Summary != nil {
		summary := toTranslationInputImportSummaryDTO(*result.Summary)
		response.Summary = &summary
	}
	return response
}

func toTranslationInputImportSummaryDTO(summary usecase.TranslationInputImportSummary) TranslationInputImportSummaryDTO {
	return TranslationInputImportSummaryDTO{
		Input: TranslationInputImportedInputDTO{
			ID:               summary.Input.ID,
			SourceFilePath:   summary.Input.SourceFilePath,
			SourceTool:       summary.Input.SourceTool,
			TargetPluginName: summary.Input.TargetPluginName,
			TargetPluginType: summary.Input.TargetPluginType,
			RecordCount:      summary.Input.RecordCount,
			ImportedAt:       summary.Input.ImportedAt.UTC().Format(time.RFC3339),
		},
		TranslationRecordCount: summary.TranslationRecordCount,
		TranslationFieldCount:  summary.TranslationFieldCount,
		Categories:             toTranslationInputCategoryDTOs(summary.Categories),
		SampleFields:           toTranslationInputSampleFieldDTOs(summary.SampleFields),
		Warnings:               toTranslationInputWarningDTOs(summary.Warnings),
	}
}

func toTranslationInputCategoryDTOs(categories []usecase.TranslationInputCategoryCount) []TranslationInputCategoryCountDTO {
	results := make([]TranslationInputCategoryCountDTO, 0, len(categories))
	for _, category := range categories {
		results = append(results, TranslationInputCategoryCountDTO{
			Category:    category.Category,
			RecordCount: category.RecordCount,
			FieldCount:  category.FieldCount,
		})
	}
	return results
}

func toTranslationInputSampleFieldDTOs(fields []usecase.TranslationInputSampleField) []TranslationInputSampleFieldDTO {
	results := make([]TranslationInputSampleFieldDTO, 0, len(fields))
	for _, field := range fields {
		results = append(results, TranslationInputSampleFieldDTO{
			RecordType:    field.RecordType,
			SubrecordType: field.SubrecordType,
			FormID:        field.FormID,
			EditorID:      field.EditorID,
			SourceText:    field.SourceText,
			Translatable:  field.Translatable,
		})
	}
	return results
}

func toTranslationInputWarningDTOs(warnings []usecase.TranslationInputWarning) []TranslationInputWarningDTO {
	results := make([]TranslationInputWarningDTO, 0, len(warnings))
	for _, warning := range warnings {
		results = append(results, TranslationInputWarningDTO{
			Kind:          warning.Kind,
			RecordType:    warning.RecordType,
			SubrecordType: warning.SubrecordType,
			Message:       warning.Message,
		})
	}
	return results
}
