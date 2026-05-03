package usecase

import (
	"context"
	"fmt"

	"aitranslationenginejp/internal/service"
)

// TranslationInputImportServicePort defines the import service required by the usecase.
type TranslationInputImportServicePort interface {
	ImportXEditJSON(ctx context.Context, filePath string) (service.TranslationInputImportSummary, error)
}

type translationInputCacheRebuildServicePort interface {
	RebuildInputCache(ctx context.Context, inputID int64) (service.TranslationInputImportSummary, error)
}

type translationInputContentImportServicePort interface {
	ImportXEditJSONWithContent(
		ctx context.Context,
		filePath string,
		fileName string,
		fileContent string,
	) (service.TranslationInputImportSummary, error)
}

// TranslationInputUsecase orchestrates the translation input import boundary.
type TranslationInputUsecase struct {
	importService TranslationInputImportServicePort
}

// TranslationInputImportedInput is the imported input metadata exposed to controllers.
type TranslationInputImportedInput = service.TranslationInputImportedInput

// TranslationInputCategoryCount aggregates imported records and fields by category.
type TranslationInputCategoryCount = service.TranslationInputCategoryCount

// TranslationInputSampleField is one representative imported field exposed at the boundary.
type TranslationInputSampleField = service.TranslationInputSampleField

// TranslationInputWarning is one non-fatal import warning exposed at the boundary.
type TranslationInputWarning = service.TranslationInputWarning

// TranslationInputImportSummary is the successful import payload exposed at the boundary.
type TranslationInputImportSummary = service.TranslationInputImportSummary

// TranslationInputImportResult returns either a summary or a validation error kind.
type TranslationInputImportResult struct {
	Summary   *TranslationInputImportSummary
	ErrorKind string
	Warnings  []TranslationInputWarning
}

// NewTranslationInputUsecase creates a translation input usecase.
func NewTranslationInputUsecase(importService TranslationInputImportServicePort) *TranslationInputUsecase {
	return &TranslationInputUsecase{importService: importService}
}

// ImportXEditJSON imports one xEdit JSON file and maps validation failures to error kinds.
func (usecase *TranslationInputUsecase) ImportXEditJSON(
	ctx context.Context,
	filePath string,
) (TranslationInputImportResult, error) {
	summary, err := usecase.importService.ImportXEditJSON(ctx, filePath)
	return mapTranslationInputSummaryOrError(summary, err)
}

// ImportXEditJSONWithContent imports one xEdit JSON payload and maps validation failures to error kinds.
func (usecase *TranslationInputUsecase) ImportXEditJSONWithContent(
	ctx context.Context,
	filePath string,
	fileName string,
	fileContent string,
) (TranslationInputImportResult, error) {
	contentImportService, ok := usecase.importService.(translationInputContentImportServicePort)
	if !ok {
		return TranslationInputImportResult{}, fmt.Errorf("import translation input: usecase service does not support content import")
	}
	summary, err := contentImportService.ImportXEditJSONWithContent(ctx, filePath, fileName, fileContent)
	return mapTranslationInputSummaryOrError(summary, err)
}

// RebuildInputCache rebuilds translation records and fields for one input ID.
func (usecase *TranslationInputUsecase) RebuildInputCache(
	ctx context.Context,
	inputID int64,
) (TranslationInputImportResult, error) {
	rebuildService, ok := usecase.importService.(translationInputCacheRebuildServicePort)
	if !ok {
		return TranslationInputImportResult{}, fmt.Errorf("rebuild translation input cache: usecase service does not support rebuild")
	}
	summary, err := rebuildService.RebuildInputCache(ctx, inputID)
	return mapTranslationInputSummaryOrError(summary, err)
}

func mapTranslationInputSummaryOrError(
	summary TranslationInputImportSummary,
	err error,
) (TranslationInputImportResult, error) {
	if err != nil {
		if kind, ok := service.TranslationInputErrorKindOf(err); ok {
			return TranslationInputImportResult{ErrorKind: kind}, nil
		}
		return TranslationInputImportResult{}, fmt.Errorf("import translation input: %w", err)
	}
	return TranslationInputImportResult{
		Summary:   &summary,
		Warnings:  summary.Warnings,
		ErrorKind: "",
	}, nil
}
