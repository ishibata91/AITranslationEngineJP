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

type translationOutputPersistedPreviewState struct {
	rowsByFieldID map[int64]repository.XTranslatorOutputRow
	rowCount      int
	compareRows   bool
}

// ReadDiffPreview returns the diff preview read model for one selected job and artifact.
func (service *TranslationOutputArtifactService) ReadDiffPreview(
	ctx context.Context,
	jobID int64,
	artifactID int64,
) (TranslationOutputDiffPreviewReadModel, error) {
	loaded, emptyResult, ok, err := service.loadDiffPreviewJob(ctx, jobID, artifactID)
	if err != nil {
		return TranslationOutputDiffPreviewReadModel{}, err
	}
	if ok {
		return emptyResult, nil
	}
	persistedState, artifactMismatchResult, ok, err := service.loadPersistedPreviewState(ctx, loaded.outputs, jobID, artifactID)
	if err != nil {
		return TranslationOutputDiffPreviewReadModel{}, err
	}
	if ok {
		return artifactMismatchResult, nil
	}
	rows, summary, err := service.buildDiffPreviewRows(ctx, loaded.outputs, persistedState)
	if err != nil {
		return TranslationOutputDiffPreviewReadModel{}, err
	}

	return TranslationOutputDiffPreviewReadModel{
		JobID:                jobID,
		ArtifactID:           artifactID,
		Rows:                 rows,
		CompatibilitySummary: summary,
	}, nil
}

func (service *TranslationOutputArtifactService) loadDiffPreviewJob(
	ctx context.Context,
	jobID int64,
	artifactID int64,
) (translationOutputArtifactLoadedJob, TranslationOutputDiffPreviewReadModel, bool, error) {
	loaded, err := service.loadJob(ctx, jobID)
	if err == nil {
		return loaded, TranslationOutputDiffPreviewReadModel{}, false, nil
	}
	if artifactID == 0 && errors.Is(err, repository.ErrNotFound) {
		return translationOutputArtifactLoadedJob{}, TranslationOutputDiffPreviewReadModel{
			JobID:                jobID,
			ArtifactID:           artifactID,
			Rows:                 []TranslationOutputDiffPreviewRowReadModel{},
			CompatibilitySummary: TranslationOutputCompatibilitySummaryReadModel{Passed: true},
		}, true, nil
	}
	return translationOutputArtifactLoadedJob{}, TranslationOutputDiffPreviewReadModel{}, false, fmt.Errorf(
		"load diff preview translation job %d: %w",
		jobID,
		err,
	)
}

func (service *TranslationOutputArtifactService) loadPersistedPreviewState(
	ctx context.Context,
	outputs []repository.JobTranslationField,
	jobID int64,
	artifactID int64,
) (translationOutputPersistedPreviewState, TranslationOutputDiffPreviewReadModel, bool, error) {
	state := translationOutputPersistedPreviewState{
		rowsByFieldID: map[int64]repository.XTranslatorOutputRow{},
	}
	if artifactID == 0 || service.persistenceRepository == nil {
		return state, TranslationOutputDiffPreviewReadModel{}, false, nil
	}
	artifact, artifactErr := service.persistenceRepository.GetTranslationArtifactByID(ctx, artifactID)
	if artifactErr != nil {
		if errors.Is(artifactErr, repository.ErrNotFound) {
			return state, TranslationOutputDiffPreviewReadModel{}, false, nil
		}
		return state, TranslationOutputDiffPreviewReadModel{}, false, fmt.Errorf(
			"get artifact %d for diff preview: %w",
			artifactID,
			artifactErr,
		)
	}
	if artifact.TranslationJobID != jobID {
		return state, TranslationOutputDiffPreviewReadModel{
			JobID:      jobID,
			ArtifactID: artifactID,
			Rows:       []TranslationOutputDiffPreviewRowReadModel{},
			CompatibilitySummary: TranslationOutputCompatibilitySummaryReadModel{
				Passed:      false,
				RejectCount: 1,
			},
		}, true, nil
	}
	persistedRows, rowErr := listPersistedArtifactRows(ctx, service.persistenceRepository, artifactID)
	if rowErr != nil {
		return state, TranslationOutputDiffPreviewReadModel{}, false, fmt.Errorf(
			"list persisted artifact rows for diff preview: %w",
			rowErr,
		)
	}
	state.rowCount = len(persistedRows)
	state.compareRows = true
	state.rowsByFieldID = buildPersistedRowsByFieldID(outputs, persistedRows)
	return state, TranslationOutputDiffPreviewReadModel{}, false, nil
}

func buildPersistedRowsByFieldID(
	outputs []repository.JobTranslationField,
	persistedRows []repository.XTranslatorOutputRow,
) map[int64]repository.XTranslatorOutputRow {
	outputsByID := make(map[int64]repository.JobTranslationField, len(outputs))
	for _, output := range outputs {
		outputsByID[output.ID] = output
	}
	rowsByFieldID := make(map[int64]repository.XTranslatorOutputRow, len(persistedRows))
	for _, row := range persistedRows {
		output, ok := outputsByID[row.JobTranslationFieldID]
		if !ok {
			continue
		}
		rowsByFieldID[output.TranslationFieldID] = row
	}
	return rowsByFieldID
}

func (service *TranslationOutputArtifactService) buildDiffPreviewRows(
	ctx context.Context,
	outputs []repository.JobTranslationField,
	persistedState translationOutputPersistedPreviewState,
) ([]TranslationOutputDiffPreviewRowReadModel, TranslationOutputCompatibilitySummaryReadModel, error) {
	rows := make([]TranslationOutputDiffPreviewRowReadModel, 0, len(outputs))
	summary := TranslationOutputCompatibilitySummaryReadModel{Passed: true}
	fieldCounts := make(map[int64]int, len(outputs))
	for _, output := range outputs {
		fieldCounts[output.TranslationFieldID]++
	}
	for _, output := range outputs {
		if fieldCounts[output.TranslationFieldID] > 1 {
			summary.RejectCount++
			continue
		}
		row, includeRow, err := service.buildDiffPreviewRow(ctx, output, persistedState, &summary)
		if err != nil {
			return nil, TranslationOutputCompatibilitySummaryReadModel{}, err
		}
		if includeRow {
			rows = append(rows, row)
		}
	}
	sort.Slice(rows, func(i, j int) bool {
		return rows[i].FieldID < rows[j].FieldID
	})
	if persistedState.compareRows && persistedState.rowCount != len(persistedState.rowsByFieldID) {
		summary.WarningCount++
	}
	summary.Passed = summary.WarningCount == 0 && summary.RejectCount == 0
	return rows, summary, nil
}

func (service *TranslationOutputArtifactService) buildDiffPreviewRow(
	ctx context.Context,
	output repository.JobTranslationField,
	persistedState translationOutputPersistedPreviewState,
	summary *TranslationOutputCompatibilitySummaryReadModel,
) (TranslationOutputDiffPreviewRowReadModel, bool, error) {
	result, buildErr := service.buildXTranslatorOutputRow(ctx, output)
	if buildErr != nil {
		return TranslationOutputDiffPreviewRowReadModel{}, false, fmt.Errorf(
			"build xtranslator output row for field %d: %w",
			output.TranslationFieldID,
			buildErr,
		)
	}
	summary.WarningCount += result.warningCount
	summary.RejectCount += result.rejectCount
	if result.rejectCount > 0 {
		return TranslationOutputDiffPreviewRowReadModel{}, false, nil
	}
	markDiffPreviewRowStale(&result.row, persistedState)
	return result.row, true, nil
}

func markDiffPreviewRowStale(
	row *TranslationOutputDiffPreviewRowReadModel,
	persistedState translationOutputPersistedPreviewState,
) {
	if !persistedState.compareRows {
		return
	}
	persistedRow, ok := persistedState.rowsByFieldID[row.FieldID]
	switch {
	case !ok:
		row.StaleReason = "generated artifact row is missing"
		row.CanRegenerate = true
	case buildPersistedXTranslatorRowDigest(row.FieldID, persistedRow) != row.RowDigest:
		row.StaleReason = "generated artifact row does not match current field result"
		row.CanRegenerate = true
	}
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
