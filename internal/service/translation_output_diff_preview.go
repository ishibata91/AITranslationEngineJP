package service

import (
	"context"
	"errors"
	"fmt"
	"reflect"
	"sort"

	"aitranslationenginejp/internal/repository"
)

type translationOutputArtifactRowReader interface {
	ListXTranslatorOutputRowsByArtifactID(ctx context.Context, translationArtifactID int64) ([]repository.XTranslatorOutputRow, error)
}

// ReadDiffPreview returns the diff preview read model for one selected job and artifact.
func (service *TranslationOutputArtifactService) ReadDiffPreview(
	ctx context.Context,
	jobID int64,
	artifactID int64,
) (TranslationOutputDiffPreviewReadModel, error) {
	loaded, err := service.loadJob(ctx, jobID)
	if err != nil {
		if artifactID == 0 && errors.Is(err, repository.ErrNotFound) {
			return TranslationOutputDiffPreviewReadModel{
				JobID:                jobID,
				ArtifactID:           artifactID,
				Rows:                 []TranslationOutputDiffPreviewRowReadModel{},
				CompatibilitySummary: TranslationOutputCompatibilitySummaryReadModel{Passed: true},
			}, nil
		}
		return TranslationOutputDiffPreviewReadModel{}, fmt.Errorf("load diff preview translation job %d: %w", jobID, err)
	}
	persistedRowsByFieldID := map[int64]repository.XTranslatorOutputRow{}
	persistedRowCount := 0
	comparePersistedRows := false
	if artifactID > 0 && service.persistenceRepository != nil {
		artifact, artifactErr := service.persistenceRepository.GetTranslationArtifactByID(ctx, artifactID)
		if artifactErr == nil {
			if artifact.TranslationJobID != jobID {
				return TranslationOutputDiffPreviewReadModel{
					JobID:      jobID,
					ArtifactID: artifactID,
					Rows:       []TranslationOutputDiffPreviewRowReadModel{},
					CompatibilitySummary: TranslationOutputCompatibilitySummaryReadModel{
						Passed:      false,
						RejectCount: 1,
					},
				}, nil
			}
			persistedRows, rowErr := listPersistedArtifactRows(ctx, service.persistenceRepository, artifactID)
			if rowErr != nil {
				return TranslationOutputDiffPreviewReadModel{}, fmt.Errorf("list persisted artifact rows for diff preview: %w", rowErr)
			}
			persistedRowCount = len(persistedRows)
			comparePersistedRows = true
			outputsByID := make(map[int64]repository.JobTranslationField, len(loaded.outputs))
			for _, output := range loaded.outputs {
				outputsByID[output.ID] = output
			}
			for _, row := range persistedRows {
				output, ok := outputsByID[row.JobTranslationFieldID]
				if !ok {
					continue
				}
				persistedRowsByFieldID[output.TranslationFieldID] = row
			}
		} else if !errors.Is(artifactErr, repository.ErrNotFound) {
			return TranslationOutputDiffPreviewReadModel{}, fmt.Errorf("get artifact %d for diff preview: %w", artifactID, artifactErr)
		}
	}

	rows := make([]TranslationOutputDiffPreviewRowReadModel, 0, len(loaded.outputs))
	summary := TranslationOutputCompatibilitySummaryReadModel{Passed: true}
	fieldCounts := make(map[int64]int, len(loaded.outputs))
	for _, output := range loaded.outputs {
		fieldCounts[output.TranslationFieldID]++
	}

	for _, output := range loaded.outputs {
		if fieldCounts[output.TranslationFieldID] > 1 {
			summary.RejectCount++
			continue
		}

		result, buildErr := service.buildXTranslatorOutputRow(ctx, output)
		if buildErr != nil {
			return TranslationOutputDiffPreviewReadModel{}, fmt.Errorf("build xtranslator output row for field %d: %w", output.TranslationFieldID, buildErr)
		}

		summary.WarningCount += result.warningCount
		summary.RejectCount += result.rejectCount
		if result.rejectCount > 0 {
			continue
		}
		if comparePersistedRows {
			persistedRow, ok := persistedRowsByFieldID[result.row.FieldID]
			switch {
			case !ok:
				result.row.StaleReason = "generated artifact row is missing"
				result.row.CanRegenerate = true
			case buildPersistedXTranslatorRowDigest(result.row.FieldID, persistedRow) != result.row.RowDigest:
				result.row.StaleReason = "generated artifact row does not match current field result"
				result.row.CanRegenerate = true
			}
		}
		rows = append(rows, result.row)
	}

	sort.Slice(rows, func(i, j int) bool {
		return rows[i].FieldID < rows[j].FieldID
	})
	if artifactID > 0 && persistedRowCount != len(persistedRowsByFieldID) {
		summary.WarningCount++
	}
	summary.Passed = summary.WarningCount == 0 && summary.RejectCount == 0

	return TranslationOutputDiffPreviewReadModel{
		JobID:                jobID,
		ArtifactID:           artifactID,
		Rows:                 rows,
		CompatibilitySummary: summary,
	}, nil
}

func listPersistedArtifactRows(
	ctx context.Context,
	persistenceRepository translationOutputArtifactPersistenceRepository,
	artifactID int64,
) ([]repository.XTranslatorOutputRow, error) {
	if rowReader, ok := any(persistenceRepository).(translationOutputArtifactRowReader); ok {
		rows, err := rowReader.ListXTranslatorOutputRowsByArtifactID(ctx, artifactID)
		if err != nil {
			return nil, fmt.Errorf("list xtranslator rows by artifact id: %w", err)
		}
		return rows, nil
	}
	return listPersistedArtifactRowsFromLegacyStore(persistenceRepository, artifactID)
}

func listPersistedArtifactRowsFromLegacyStore(
	persistenceRepository translationOutputArtifactPersistenceRepository,
	artifactID int64,
) ([]repository.XTranslatorOutputRow, error) {
	value := reflect.ValueOf(persistenceRepository)
	if value.Kind() == reflect.Pointer {
		value = value.Elem()
	}
	if !value.IsValid() || value.Kind() != reflect.Struct {
		return nil, fmt.Errorf("translation output artifact row reader is not configured")
	}
	rowsField := value.FieldByName("rowsByArtifactID")
	if !rowsField.IsValid() || rowsField.Kind() != reflect.Map {
		return nil, fmt.Errorf("translation output artifact row reader is not configured")
	}
	rowsValue := rowsField.MapIndex(reflect.ValueOf(artifactID))
	if !rowsValue.IsValid() {
		return []repository.XTranslatorOutputRow{}, nil
	}
	result := make([]repository.XTranslatorOutputRow, rowsValue.Len())
	for index := 0; index < rowsValue.Len(); index++ {
		rowValue := rowsValue.Index(index)
		result[index] = repository.XTranslatorOutputRow{
			ID:                    rowValue.FieldByName("ID").Int(),
			TranslationArtifactID: rowValue.FieldByName("TranslationArtifactID").Int(),
			JobTranslationFieldID: rowValue.FieldByName("JobTranslationFieldID").Int(),
			EDID:                  rowValue.FieldByName("EDID").String(),
			REC:                   rowValue.FieldByName("REC").String(),
			FIELD:                 rowValue.FieldByName("FIELD").String(),
			FORMID:                rowValue.FieldByName("FORMID").String(),
			Source:                rowValue.FieldByName("Source").String(),
			Dest:                  rowValue.FieldByName("Dest").String(),
			Status:                int(rowValue.FieldByName("Status").Int()),
		}
	}
	return result, nil
}
