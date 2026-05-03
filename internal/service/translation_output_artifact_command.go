package service

import (
	"context"
	"errors"
	"fmt"
	"sort"
	"strings"
	"time"

	"aitranslationenginejp/internal/repository"
)

const (
	translationOutputArtifactFormatXML               = "xtranslator_xml"
	translationOutputArtifactStatusSuccess           = "success"
	translationOutputArtifactStatusFailed            = "failed"
	translationOutputArtifactStatusRejected          = "rejected"
	translationOutputArtifactReasonFileWriteFailed   = "output artifact file write failed"
	translationOutputArtifactReasonPersistenceFailed = "output artifact persistence failed"
)

// TranslationOutputArtifactCommandResult stores one generate or regenerate result.
type TranslationOutputArtifactCommandResult struct {
	Artifact            repository.TranslationArtifact
	RowCount            int
	OperationKind       string
	ReplacedArtifactID  int64
	AffectedFieldIDs    []int64
	DuplicateRowCreated bool
}

type translationOutputArtifactStatusLookup struct {
	ArtifactID     int64
	Status         string
	RowCount       int
	CurrentVersion bool
}

type xTranslatorArtifactRow struct {
	FieldID               int64
	JobTranslationFieldID int64
	EDID                  string
	REC                   string
	FIELD                 string
	FORMID                string
	Source                string
	Dest                  string
	Status                int
}

type xTranslatorArtifactBuildResult struct {
	rows             []xTranslatorArtifactRow
	affectedFieldIDs []int64
	warningCount     int
	rejectKind       string
}

// TranslationOutputArtifactFailure stores the redacted public failure summary.
type TranslationOutputArtifactFailure struct {
	artifactStatus string
	errorKind      string
	reason         string
	retryable      bool
}

// GenerateArtifact writes a new xTranslator-compatible XML artifact.
func (service *TranslationOutputArtifactService) GenerateArtifact(
	ctx context.Context,
	jobID int64,
	targetGame string,
	outputPath string,
) (TranslationOutputArtifactCommandResult, error) {
	return service.writeArtifact(ctx, translationOutputArtifactWriteRequest{
		jobID:              jobID,
		targetGame:         targetGame,
		outputPath:         outputPath,
		operationKind:      "generate",
		expectedArtifactID: 0,
	})
}

// RegenerateArtifact rewrites the current artifact for one completed job.
func (service *TranslationOutputArtifactService) RegenerateArtifact(
	ctx context.Context,
	jobID int64,
	artifactID int64,
	targetGame string,
	outputPath string,
) (TranslationOutputArtifactCommandResult, error) {
	return service.writeArtifact(ctx, translationOutputArtifactWriteRequest{
		jobID:              jobID,
		targetGame:         targetGame,
		outputPath:         outputPath,
		operationKind:      "regenerate",
		expectedArtifactID: artifactID,
	})
}

type translationOutputArtifactWriteRequest struct {
	jobID              int64
	targetGame         string
	outputPath         string
	operationKind      string
	expectedArtifactID int64
}

func (service *TranslationOutputArtifactService) writeArtifact(
	ctx context.Context,
	request translationOutputArtifactWriteRequest,
) (TranslationOutputArtifactCommandResult, error) {
	failureResult := newTranslationOutputArtifactFailureResult(request)
	loaded, err := service.loadJob(ctx, request.jobID)
	if err != nil {
		return TranslationOutputArtifactCommandResult{}, fmt.Errorf("load translation output artifact job %d: %w", request.jobID, err)
	}
	validatedOutputPath, failure, err := service.validateArtifactWriteRequest(loaded, request, failureResult)
	if err != nil {
		return failure, err
	}

	buildResult, err := service.buildArtifactRows(ctx, loaded.outputs)
	if err != nil {
		return TranslationOutputArtifactCommandResult{}, fmt.Errorf("build translation output rows: %w", err)
	}
	if buildResult.rejectKind != "" {
		failureResult.Artifact.Status = translationOutputArtifactStatusRejected
		failureResult.AffectedFieldIDs = normalizeOperationFieldIDs(buildResult.affectedFieldIDs)
		return failureResult, TranslationOutputArtifactFailure{
			artifactStatus: translationOutputArtifactStatusRejected,
			errorKind:      buildResult.rejectKind,
			reason:         buildArtifactRejectReason(buildResult.rejectKind),
			retryable:      false,
		}
	}
	failure, expectedArtifactErr := service.ensureExpectedArtifactMatches(
		ctx,
		request,
		buildResult.affectedFieldIDs,
		failureResult,
	)
	if expectedArtifactErr != nil {
		return failure, expectedArtifactErr
	}

	xmlPayload, err := service.xmlSerializer.Serialize(strings.TrimSpace(request.targetGame), buildResult.rows)
	if err != nil {
		failureResult.Artifact.Status = translationOutputArtifactStatusFailed
		failureResult.AffectedFieldIDs = normalizeOperationFieldIDs(buildResult.affectedFieldIDs)
		return failureResult, TranslationOutputArtifactFailure{
			artifactStatus: translationOutputArtifactStatusFailed,
			errorKind:      "xml_serialization_failed",
			reason:         "xtranslator xml serialization failed",
			retryable:      false,
		}
	}

	artifactDraft, rowDrafts := buildTranslationOutputArtifactPersistenceDrafts(request, validatedOutputPath, buildResult.rows)
	persisted, err := service.persistArtifactWithFile(ctx, validatedOutputPath, xmlPayload, artifactDraft, rowDrafts)
	if err != nil {
		failureResult.Artifact.Status = translationOutputArtifactStatusFailed
		failureResult.AffectedFieldIDs = normalizeOperationFieldIDs(buildResult.affectedFieldIDs)
		return failureResult, err
	}
	if request.expectedArtifactID > 0 && persisted.Artifact.ID != request.expectedArtifactID {
		if transactionalWriter, ok := service.fileWriter.(translationOutputArtifactTransactionalFileWriter); ok {
			_ = transactionalWriter.RemoveFile(validatedOutputPath)
		}
		failureResult.Artifact.Status = translationOutputArtifactStatusFailed
		failureResult.AffectedFieldIDs = normalizeOperationFieldIDs(buildResult.affectedFieldIDs)
		return failureResult, TranslationOutputArtifactFailure{
			artifactStatus: translationOutputArtifactStatusFailed,
			errorKind:      "artifact_save_failed",
			reason:         "requested artifact id does not match current output artifact",
			retryable:      false,
		}
	}

	persisted.OperationKind = request.operationKind
	persisted.ReplacedArtifactID = request.expectedArtifactID
	if request.expectedArtifactID == 0 {
		persisted.ReplacedArtifactID = persisted.Artifact.ID
	}
	persisted.AffectedFieldIDs = normalizeOperationFieldIDs(buildResult.affectedFieldIDs)
	persisted.DuplicateRowCreated = false

	return persisted, nil
}

func newTranslationOutputArtifactFailureResult(
	request translationOutputArtifactWriteRequest,
) TranslationOutputArtifactCommandResult {
	return TranslationOutputArtifactCommandResult{
		Artifact: repository.TranslationArtifact{
			ID:               request.expectedArtifactID,
			TranslationJobID: request.jobID,
			TargetGame:       strings.TrimSpace(request.targetGame),
			FilePath:         strings.TrimSpace(request.outputPath),
		},
		OperationKind:      request.operationKind,
		ReplacedArtifactID: request.expectedArtifactID,
	}
}

func (service *TranslationOutputArtifactService) validateArtifactWriteRequest(
	loaded translationOutputArtifactLoadedJob,
	request translationOutputArtifactWriteRequest,
	failureResult TranslationOutputArtifactCommandResult,
) (string, TranslationOutputArtifactCommandResult, error) {
	validatedOutputPath, pathErr := validateTranslationOutputArtifactPath(strings.TrimSpace(request.outputPath))
	if pathErr != nil {
		failure, err := buildTranslationOutputArtifactFailureResult(
			failureResult,
			translationOutputArtifactStatusFailed,
			"file_write_failed",
			translationOutputArtifactReasonFileWriteFailed,
			true,
			nil,
		)
		return "", failure, err
	}
	readiness := service.buildReadiness(loaded)
	if !readiness.Ready {
		failure, err := buildTranslationOutputArtifactFailureResult(
			failureResult,
			translationOutputArtifactStatusRejected,
			readiness.RejectionKind,
			buildReadinessReason(readiness.RejectionKind),
			readiness.Retryable,
			nil,
		)
		return "", failure, err
	}
	return validatedOutputPath, TranslationOutputArtifactCommandResult{}, nil
}

func (service *TranslationOutputArtifactService) ensureExpectedArtifactMatches(
	ctx context.Context,
	request translationOutputArtifactWriteRequest,
	affectedFieldIDs []int64,
	failureResult TranslationOutputArtifactCommandResult,
) (TranslationOutputArtifactCommandResult, error) {
	if request.expectedArtifactID == 0 {
		return TranslationOutputArtifactCommandResult{}, nil
	}
	currentArtifact, currentArtifactErr := service.lookupPersistedArtifact(ctx, request.jobID)
	if currentArtifactErr != nil {
		return buildTranslationOutputArtifactFailureResult(
			failureResult,
			translationOutputArtifactStatusFailed,
			"artifact_save_failed",
			translationOutputArtifactReasonPersistenceFailed,
			!errors.Is(currentArtifactErr, repository.ErrNotFound),
			affectedFieldIDs,
		)
	}
	if currentArtifact.ID != request.expectedArtifactID {
		return buildTranslationOutputArtifactFailureResult(
			failureResult,
			translationOutputArtifactStatusFailed,
			"artifact_save_failed",
			"requested artifact id does not match current output artifact",
			false,
			affectedFieldIDs,
		)
	}
	return TranslationOutputArtifactCommandResult{}, nil
}

func buildTranslationOutputArtifactPersistenceDrafts(
	request translationOutputArtifactWriteRequest,
	outputPath string,
	rows []xTranslatorArtifactRow,
) (repository.TranslationArtifactDraft, []repository.XTranslatorOutputRowDraft) {
	generatedAt := time.Now().UTC()
	artifactDraft := repository.TranslationArtifactDraft{
		TranslationJobID: request.jobID,
		ArtifactFormat:   translationOutputArtifactFormatXML,
		TargetGame:       strings.TrimSpace(request.targetGame),
		FilePath:         outputPath,
		Status:           translationOutputArtifactStatusSuccess,
		GeneratedAt:      &generatedAt,
	}
	rowDrafts := make([]repository.XTranslatorOutputRowDraft, 0, len(rows))
	for _, row := range rows {
		rowDrafts = append(rowDrafts, repository.XTranslatorOutputRowDraft{
			JobTranslationFieldID: row.JobTranslationFieldID,
			EDID:                  row.EDID,
			REC:                   row.REC,
			FIELD:                 row.FIELD,
			FORMID:                row.FORMID,
			Source:                row.Source,
			Dest:                  row.Dest,
			Status:                row.Status,
		})
	}
	return artifactDraft, rowDrafts
}

func buildTranslationOutputArtifactFailureResult(
	result TranslationOutputArtifactCommandResult,
	status string,
	errorKind string,
	reason string,
	retryable bool,
	affectedFieldIDs []int64,
) (TranslationOutputArtifactCommandResult, error) {
	result.Artifact.Status = status
	if affectedFieldIDs != nil {
		result.AffectedFieldIDs = normalizeOperationFieldIDs(affectedFieldIDs)
	}
	return result, TranslationOutputArtifactFailure{
		artifactStatus: status,
		errorKind:      errorKind,
		reason:         reason,
		retryable:      retryable,
	}
}

func (service *TranslationOutputArtifactService) buildArtifactRows(
	ctx context.Context,
	outputs []repository.JobTranslationField,
) (xTranslatorArtifactBuildResult, error) {
	result := xTranslatorArtifactBuildResult{}
	fieldCounts := make(map[int64]int, len(outputs))
	for _, output := range outputs {
		fieldCounts[output.TranslationFieldID]++
	}
	for _, output := range outputs {
		if fieldCounts[output.TranslationFieldID] > 1 {
			result.rejectKind = "compatibility_rejected"
			return result, nil
		}
		field, err := service.translationSourceReader.GetTranslationFieldByID(ctx, output.TranslationFieldID)
		if err != nil {
			return xTranslatorArtifactBuildResult{}, fmt.Errorf("get translation field %d: %w", output.TranslationFieldID, err)
		}
		record, err := service.translationSourceReader.GetTranslationRecordByID(ctx, field.TranslationRecordID)
		if err != nil {
			return xTranslatorArtifactBuildResult{}, fmt.Errorf("get translation record %d: %w", field.TranslationRecordID, err)
		}
		row := TranslationOutputDiffPreviewRowReadModel{
			FieldID:              field.ID,
			EDID:                 strings.TrimSpace(record.EditorID),
			REC:                  strings.TrimSpace(record.RecordType),
			FIELD:                strings.TrimSpace(field.SubrecordType),
			FORMID:               strings.TrimSpace(record.FormID),
			SourceExcerpt:        excerptForDiffPreview(field.SourceText),
			DestExcerpt:          excerptForDiffPreview(output.TranslatedText),
			InternalOutputStatus: strings.ToLower(strings.TrimSpace(output.OutputStatus)),
		}
		if !hasAllRequiredXTranslatorColumns(row) ||
			strings.TrimSpace(field.SourceText) == "" ||
			strings.TrimSpace(output.TranslatedText) == "" {
			result.rejectKind = "missing_required_row_field"
			return result, nil
		}
		status, _, ok := mapOutputStatusToXTranslator(row.InternalOutputStatus)
		if !ok {
			result.rejectKind = "unknown_output_status"
			return result, nil
		}
		warnings, rejects := validateXTranslatorOutputRow(row, field.SourceText, output.TranslatedText)
		result.warningCount += warnings
		if rejects > 0 {
			result.rejectKind = "compatibility_rejected"
			return result, nil
		}
		result.rows = append(result.rows, xTranslatorArtifactRow{
			FieldID:               field.ID,
			JobTranslationFieldID: output.ID,
			EDID:                  strings.TrimSpace(record.EditorID),
			REC:                   strings.TrimSpace(record.RecordType),
			FIELD:                 strings.TrimSpace(field.SubrecordType),
			FORMID:                strings.TrimSpace(record.FormID),
			Source:                field.SourceText,
			Dest:                  output.TranslatedText,
			Status:                status,
		})
		result.affectedFieldIDs = append(result.affectedFieldIDs, field.ID)
	}

	sort.Slice(result.rows, func(i, j int) bool {
		return result.rows[i].FieldID < result.rows[j].FieldID
	})
	sort.Slice(result.affectedFieldIDs, func(i, j int) bool {
		return result.affectedFieldIDs[i] < result.affectedFieldIDs[j]
	})
	return result, nil
}

func (service *TranslationOutputArtifactService) persistArtifact(
	ctx context.Context,
	artifactDraft repository.TranslationArtifactDraft,
	rowDrafts []repository.XTranslatorOutputRowDraft,
) (TranslationOutputArtifactCommandResult, error) {
	if service.persistenceRepository == nil {
		return TranslationOutputArtifactCommandResult{}, fmt.Errorf("translation output artifact persistence repository is not configured")
	}
	execute := func(txCtx context.Context) (TranslationOutputArtifactCommandResult, error) {
		artifact, err := service.persistenceRepository.UpsertTranslationArtifact(txCtx, artifactDraft)
		if err != nil {
			return TranslationOutputArtifactCommandResult{}, fmt.Errorf("upsert translation artifact: %w", err)
		}
		rows, err := service.persistenceRepository.ReplaceXTranslatorOutputRows(txCtx, artifact.ID, rowDrafts)
		if err != nil {
			return TranslationOutputArtifactCommandResult{}, fmt.Errorf("replace xtranslator output rows: %w", err)
		}
		return TranslationOutputArtifactCommandResult{
			Artifact: artifact,
			RowCount: len(rows),
		}, nil
	}
	if service.transactor == nil {
		return execute(ctx)
	}
	var result TranslationOutputArtifactCommandResult
	if err := service.transactor.WithTransaction(ctx, func(txCtx context.Context) error {
		persisted, err := execute(txCtx)
		if err != nil {
			return err
		}
		result = persisted
		return nil
	}); err != nil {
		return TranslationOutputArtifactCommandResult{}, fmt.Errorf("persist translation output artifact transaction: %w", err)
	}
	return result, nil
}

func (service *TranslationOutputArtifactService) persistArtifactWithFile(
	ctx context.Context,
	outputPath string,
	xmlPayload []byte,
	artifactDraft repository.TranslationArtifactDraft,
	rowDrafts []repository.XTranslatorOutputRowDraft,
) (TranslationOutputArtifactCommandResult, error) {
	if transactionalWriter, ok := service.fileWriter.(translationOutputArtifactTransactionalFileWriter); ok {
		tempPath, err := transactionalWriter.WriteTemporaryFile(outputPath, xmlPayload)
		if err != nil {
			return TranslationOutputArtifactCommandResult{}, TranslationOutputArtifactFailure{
				artifactStatus: translationOutputArtifactStatusFailed,
				errorKind:      "file_write_failed",
				reason:         translationOutputArtifactReasonFileWriteFailed,
				retryable:      true,
			}
		}
		if publishErr := transactionalWriter.PublishTemporaryFile(tempPath, outputPath); publishErr != nil {
			_ = transactionalWriter.RemoveFile(tempPath)
			_ = transactionalWriter.RemoveFile(outputPath)
			return TranslationOutputArtifactCommandResult{}, TranslationOutputArtifactFailure{
				artifactStatus: translationOutputArtifactStatusFailed,
				errorKind:      "file_write_failed",
				reason:         translationOutputArtifactReasonFileWriteFailed,
				retryable:      true,
			}
		}
		persisted, err := service.persistArtifact(ctx, artifactDraft, rowDrafts)
		if err != nil {
			_ = transactionalWriter.RemoveFile(outputPath)
			return TranslationOutputArtifactCommandResult{}, TranslationOutputArtifactFailure{
				artifactStatus: translationOutputArtifactStatusFailed,
				errorKind:      "artifact_save_failed",
				reason:         translationOutputArtifactReasonPersistenceFailed,
				retryable:      true,
			}
		}
		return persisted, nil
	}

	if err := service.fileWriter.WriteFile(outputPath, xmlPayload); err != nil {
		return TranslationOutputArtifactCommandResult{}, TranslationOutputArtifactFailure{
			artifactStatus: translationOutputArtifactStatusFailed,
			errorKind:      "file_write_failed",
			reason:         translationOutputArtifactReasonFileWriteFailed,
			retryable:      true,
		}
	}
	persisted, err := service.persistArtifact(ctx, artifactDraft, rowDrafts)
	if err != nil {
		return TranslationOutputArtifactCommandResult{}, TranslationOutputArtifactFailure{
			artifactStatus: translationOutputArtifactStatusFailed,
			errorKind:      "artifact_save_failed",
			reason:         translationOutputArtifactReasonPersistenceFailed,
			retryable:      true,
		}
	}
	return persisted, nil
}

func (service *TranslationOutputArtifactService) lookupArtifactStatus(
	ctx context.Context,
	jobID int64,
) translationOutputArtifactStatusLookup {
	artifact, err := service.lookupPersistedArtifact(ctx, jobID)
	if err == nil {
		rowCount, countErr := service.persistenceRepository.CountXTranslatorOutputRowsByArtifactID(ctx, artifact.ID)
		if countErr != nil {
			rowCount = 0
		}
		return translationOutputArtifactStatusLookup{
			ArtifactID:     artifact.ID,
			Status:         artifact.Status,
			RowCount:       rowCount,
			CurrentVersion: true,
		}
	}
	if err != nil && !errors.Is(err, repository.ErrNotFound) {
		return translationOutputArtifactStatusLookup{
			Status:         translationOutputArtifactStatusFailed,
			CurrentVersion: false,
		}
	}
	return translationOutputArtifactStatusLookup{
		Status:         translationOutputArtifactStatusNotGenerated,
		CurrentVersion: false,
	}
}

func (service *TranslationOutputArtifactService) lookupPersistedArtifact(
	ctx context.Context,
	jobID int64,
) (repository.TranslationArtifact, error) {
	if service.persistenceRepository == nil {
		return repository.TranslationArtifact{}, repository.ErrNotFound
	}
	artifact, err := service.persistenceRepository.GetTranslationArtifactByJobID(ctx, jobID)
	if err != nil {
		return repository.TranslationArtifact{}, fmt.Errorf("get translation artifact by job id: %w", err)
	}
	return artifact, nil
}

func buildArtifactRejectReason(kind string) string {
	switch kind {
	case "missing_required_row_field":
		return "required xtranslator row field is missing"
	case "unknown_output_status":
		return "unknown internal output status cannot map to xtranslator status"
	case "compatibility_rejected":
		return "xtranslator compatibility validation rejected the output rows"
	default:
		return "output artifact generation was rejected"
	}
}

func (failure TranslationOutputArtifactFailure) Error() string {
	return failure.reason
}

// IsTranslationOutputArtifactFailure exposes the redacted command failure summary.
func IsTranslationOutputArtifactFailure(err error) (TranslationOutputArtifactFailure, bool) {
	var failure TranslationOutputArtifactFailure
	if errors.As(err, &failure) {
		return failure, true
	}
	return TranslationOutputArtifactFailure{}, false
}

// ArtifactStatus returns the public artifact status for the failure.
func (failure TranslationOutputArtifactFailure) ArtifactStatus() string {
	return failure.artifactStatus
}

// ErrorKind returns the public error kind for the failure.
func (failure TranslationOutputArtifactFailure) ErrorKind() string {
	return failure.errorKind
}

// Reason returns the redacted reason for the failure.
func (failure TranslationOutputArtifactFailure) Reason() string {
	return failure.reason
}

// Retryable returns whether the failure can be retried.
func (failure TranslationOutputArtifactFailure) Retryable() bool {
	return failure.retryable
}

func dedupeAffectedFieldIDs(fieldIDs []int64) []int64 {
	if len(fieldIDs) == 0 {
		return nil
	}
	result := make([]int64, 0, len(fieldIDs))
	for _, fieldID := range fieldIDs {
		if len(result) == 0 || result[len(result)-1] != fieldID {
			result = append(result, fieldID)
		}
	}
	return result
}

func normalizeOperationFieldIDs(fieldIDs []int64) []int64 {
	if len(fieldIDs) == 0 {
		return nil
	}
	cloned := append([]int64(nil), fieldIDs...)
	sort.Slice(cloned, func(i, j int) bool { return cloned[i] < cloned[j] })
	return dedupeAffectedFieldIDs(cloned)
}
