package service

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"strings"
	"unicode/utf8"

	"aitranslationenginejp/internal/repository"
)

const (
	xTranslatorStatusTranslated       = 0
	xTranslatorStatusCached           = 1
	xTranslatorRowExcerptLimit        = 32
	xTranslatorCompatibilityTextLimit = 4096
)

type xTranslatorOutputRowBuildResult struct {
	row          TranslationOutputDiffPreviewRowReadModel
	warningCount int
	rejectCount  int
}

func (service *TranslationOutputArtifactService) buildXTranslatorOutputRow(
	ctx context.Context,
	output repository.JobTranslationField,
) (xTranslatorOutputRowBuildResult, error) {
	field, err := service.translationSourceReader.GetTranslationFieldByID(ctx, output.TranslationFieldID)
	if err != nil {
		return xTranslatorOutputRowBuildResult{}, fmt.Errorf("get translation field %d: %w", output.TranslationFieldID, err)
	}
	record, err := service.translationSourceReader.GetTranslationRecordByID(ctx, field.TranslationRecordID)
	if err != nil {
		return xTranslatorOutputRowBuildResult{}, fmt.Errorf("get translation record %d: %w", field.TranslationRecordID, err)
	}

	row := TranslationOutputDiffPreviewRowReadModel{
		FieldID:              output.TranslationFieldID,
		EDID:                 strings.TrimSpace(record.EditorID),
		REC:                  strings.TrimSpace(record.RecordType),
		FIELD:                strings.TrimSpace(field.SubrecordType),
		FORMID:               strings.TrimSpace(record.FormID),
		SourceExcerpt:        excerptForDiffPreview(field.SourceText),
		DestExcerpt:          excerptForDiffPreview(output.TranslatedText),
		InternalOutputStatus: strings.ToLower(strings.TrimSpace(output.OutputStatus)),
	}

	if !hasAllRequiredXTranslatorColumns(row) {
		return xTranslatorOutputRowBuildResult{rejectCount: 1}, nil
	}

	status, reflectionSummary, ok := mapOutputStatusToXTranslator(row.InternalOutputStatus)
	if !ok {
		return xTranslatorOutputRowBuildResult{rejectCount: 1}, nil
	}
	row.XTranslatorStatus = status
	row.RowReflectionSummary = reflectionSummary
	row.RowDigest = buildXTranslatorRowDigest(row, field.SourceText, output.TranslatedText)

	warningCount, rejectCount := validateXTranslatorOutputRow(row, field.SourceText, output.TranslatedText)
	return xTranslatorOutputRowBuildResult{
		row:          row,
		warningCount: warningCount,
		rejectCount:  rejectCount,
	}, nil
}

func hasAllRequiredXTranslatorColumns(row TranslationOutputDiffPreviewRowReadModel) bool {
	return row.EDID != "" &&
		row.REC != "" &&
		row.FIELD != "" &&
		row.FORMID != "" &&
		row.SourceExcerpt != "" &&
		row.DestExcerpt != ""
}

func mapOutputStatusToXTranslator(status string) (int, string, bool) {
	switch status {
	case "translated":
		return xTranslatorStatusTranslated, "translated output reflected as xTranslator status 0", true
	case "cached":
		return xTranslatorStatusCached, "cached output reflected as xTranslator status 1; dictionary replacement kept in internal summary", true
	default:
		return 0, "", false
	}
}

func validateXTranslatorOutputRow(
	row TranslationOutputDiffPreviewRowReadModel,
	sourceText string,
	destText string,
) (warningCount int, rejectCount int) {
	if utf8.RuneCountInString(sourceText) > xTranslatorCompatibilityTextLimit ||
		utf8.RuneCountInString(destText) > xTranslatorCompatibilityTextLimit {
		warningCount++
	}
	if strings.EqualFold(row.FIELD, "WOOP") {
		warningCount++
	}
	if hasEdgeWhitespace(sourceText) || hasEdgeWhitespace(destText) {
		warningCount++
	}
	return warningCount, rejectCount
}

func hasEdgeWhitespace(value string) bool {
	if value == "" {
		return false
	}
	return value != strings.TrimSpace(value)
}

func excerptForDiffPreview(value string) string {
	trimmed := strings.TrimSpace(value)
	if trimmed == "" {
		return ""
	}
	runes := []rune(trimmed)
	if len(runes) <= xTranslatorRowExcerptLimit {
		return trimmed
	}
	return string(runes[:xTranslatorRowExcerptLimit]) + "..."
}

func buildXTranslatorRowDigest(
	row TranslationOutputDiffPreviewRowReadModel,
	sourceText string,
	destText string,
) string {
	sum := sha256.Sum256([]byte(strings.Join([]string{
		fmt.Sprint(row.FieldID),
		row.EDID,
		row.REC,
		row.FIELD,
		row.FORMID,
		row.InternalOutputStatus,
		sourceText,
		destText,
	}, "|")))
	return "sha256:" + hex.EncodeToString(sum[:])
}

func buildPersistedXTranslatorRowDigest(
	fieldID int64,
	row repository.XTranslatorOutputRow,
) string {
	internalStatus := ""
	switch row.Status {
	case xTranslatorStatusTranslated:
		internalStatus = "translated"
	case xTranslatorStatusCached:
		internalStatus = "cached"
	}
	return buildXTranslatorRowDigest(
		TranslationOutputDiffPreviewRowReadModel{
			FieldID:              fieldID,
			EDID:                 strings.TrimSpace(row.EDID),
			REC:                  strings.TrimSpace(row.REC),
			FIELD:                strings.TrimSpace(row.FIELD),
			FORMID:               strings.TrimSpace(row.FORMID),
			InternalOutputStatus: internalStatus,
		},
		row.Source,
		row.Dest,
	)
}
