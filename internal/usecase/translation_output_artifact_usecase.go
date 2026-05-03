package usecase

import (
	"context"
	"fmt"

	"aitranslationenginejp/internal/service"
)

type translationOutputArtifactServicePort interface {
	ReadReview(ctx context.Context, selectedJobID *int64) (service.TranslationOutputReviewReadModel, error)
	ReadDiffPreview(ctx context.Context, jobID int64, artifactID int64) (service.TranslationOutputDiffPreviewReadModel, error)
	GenerateArtifact(ctx context.Context, jobID int64, targetGame string, outputPath string) (service.TranslationOutputArtifactCommandResult, error)
	RegenerateArtifact(ctx context.Context, jobID int64, artifactID int64, targetGame string, outputPath string) (service.TranslationOutputArtifactCommandResult, error)
}

type translationOutputArtifactUsecase struct {
	service translationOutputArtifactServicePort
}

// NewTranslationOutputArtifactUsecase creates the output artifact usecase bridge.
func NewTranslationOutputArtifactUsecase(servicePort translationOutputArtifactServicePort) TranslationOutputArtifactUsecase {
	return &translationOutputArtifactUsecase{service: servicePort}
}

// GetTranslationOutputReview returns the current Output Review contract.
func (usecase *translationOutputArtifactUsecase) GetTranslationOutputReview(
	ctx context.Context,
	request GetTranslationOutputReviewRequest,
) (TranslationOutputReviewResult, error) {
	readModel, err := usecase.service.ReadReview(ctx, request.SelectedJobID)
	result := toTranslationOutputReviewResult(readModel)
	if err != nil {
		return result, fmt.Errorf("get translation output review: %w", err)
	}
	return result, nil
}

// GetTranslationOutputDiffPreview returns the current diff preview contract.
func (usecase *translationOutputArtifactUsecase) GetTranslationOutputDiffPreview(
	ctx context.Context,
	request GetTranslationOutputDiffPreviewRequest,
) (TranslationOutputDiffPreviewResult, error) {
	readModel, err := usecase.service.ReadDiffPreview(ctx, request.JobID, request.ArtifactID)
	result := toTranslationOutputDiffPreviewResult(readModel)
	if err != nil {
		return result, fmt.Errorf("get translation output diff preview: %w", err)
	}
	return result, nil
}

// GenerateXTranslatorOutputArtifact executes the frozen artifact generation command.
func (usecase *translationOutputArtifactUsecase) GenerateXTranslatorOutputArtifact(
	ctx context.Context,
	request GenerateXTranslatorOutputArtifactRequest,
) (XTranslatorOutputArtifactCommandResult, error) {
	result, err := usecase.service.GenerateArtifact(ctx, request.JobID, request.TargetGame, request.OutputPath)
	if err != nil {
		if _, ok := service.IsTranslationOutputArtifactFailure(err); !ok {
			return XTranslatorOutputArtifactCommandResult{}, fmt.Errorf("generate translation output artifact: %w", err)
		}
	}
	return toTranslationOutputArtifactCommandResult(result, "generate", 0, err), nil
}

// RegenerateXTranslatorOutputArtifact executes the frozen artifact regeneration command.
func (usecase *translationOutputArtifactUsecase) RegenerateXTranslatorOutputArtifact(
	ctx context.Context,
	request RegenerateXTranslatorOutputArtifactRequest,
) (XTranslatorOutputArtifactCommandResult, error) {
	result, err := usecase.service.RegenerateArtifact(
		ctx,
		request.JobID,
		request.ArtifactID,
		request.TargetGame,
		request.OutputPath,
	)
	if err != nil {
		if _, ok := service.IsTranslationOutputArtifactFailure(err); !ok {
			return XTranslatorOutputArtifactCommandResult{}, fmt.Errorf("regenerate translation output artifact: %w", err)
		}
	}
	return toTranslationOutputArtifactCommandResult(result, "regenerate", request.ArtifactID, err), nil
}

func toTranslationOutputReviewResult(
	readModel service.TranslationOutputReviewReadModel,
) TranslationOutputReviewResult {
	completedJobs := make([]TranslationOutputCompletedJobSummary, 0, len(readModel.CompletedJobs))
	for _, job := range readModel.CompletedJobs {
		completedJobs = append(completedJobs, TranslationOutputCompletedJobSummary{
			JobID:                    job.JobID,
			JobStatus:                job.JobStatus,
			ArtifactStatus:           job.ArtifactStatus,
			OutputReady:              job.OutputReady,
			TranslatedCount:          job.TranslatedCount,
			OutputStatusDistribution: cloneStringIntMap(job.OutputStatusDistribution),
		})
	}

	rejectionReasons := make([]TranslationOutputArtifactErrorSummary, 0, len(readModel.RejectionReasons))
	for _, reason := range readModel.RejectionReasons {
		rejectionReasons = append(rejectionReasons, TranslationOutputArtifactErrorSummary{
			ErrorKind:  normalizeTranslationOutputArtifactErrorKind(reason.ErrorKind),
			Reason:     reason.Reason,
			Retryable:  reason.Retryable,
			IsRedacted: reason.IsRedacted,
		})
	}

	return TranslationOutputReviewResult{
		CompletedJobs: completedJobs,
		SelectedJob: TranslationOutputSelectedJobSummary{
			JobID:           readModel.SelectedJob.JobID,
			JobStatus:       readModel.SelectedJob.JobStatus,
			BodyPhaseStatus: readModel.SelectedJob.BodyPhaseStatus,
			OutputReady:     readModel.SelectedJob.OutputReady,
			ResultSummary: TranslationOutputResultSummary{
				TranslatedCount: readModel.SelectedJob.ResultSummary.TranslatedCount,
				RowCount:        readModel.SelectedJob.ResultSummary.RowCount,
				InputProvenance: TranslationOutputInputProvenanceSummary{
					InputSnapshotDigest: readModel.SelectedJob.ResultSummary.InputProvenance.InputSnapshotDigest,
					SourceFileDigest:    readModel.SelectedJob.ResultSummary.InputProvenance.SourceFileDigest,
				},
			},
		},
		OutputReadiness: TranslationOutputReadinessSummary{
			Ready:         readModel.OutputReadiness.Ready,
			Retryable:     readModel.OutputReadiness.Retryable,
			RejectionKind: normalizeTranslationOutputArtifactErrorKind(readModel.OutputReadiness.RejectionKind),
		},
		ArtifactStatus: TranslationOutputArtifactStatusSummary{
			ArtifactID:     readModel.ArtifactStatus.ArtifactID,
			Status:         readModel.ArtifactStatus.Status,
			RowCount:       readModel.ArtifactStatus.RowCount,
			CurrentVersion: readModel.ArtifactStatus.CurrentVersion,
		},
		RejectionReasons: rejectionReasons,
	}
}

func toTranslationOutputDiffPreviewResult(
	readModel service.TranslationOutputDiffPreviewReadModel,
) TranslationOutputDiffPreviewResult {
	rows := make([]TranslationOutputDiffPreviewRow, 0, len(readModel.Rows))
	for _, row := range readModel.Rows {
		rows = append(rows, TranslationOutputDiffPreviewRow{
			FieldID:              row.FieldID,
			RowDigest:            row.RowDigest,
			EDID:                 row.EDID,
			REC:                  row.REC,
			FIELD:                row.FIELD,
			FORMID:               row.FORMID,
			SourceExcerpt:        row.SourceExcerpt,
			DestExcerpt:          row.DestExcerpt,
			XTranslatorStatus:    row.XTranslatorStatus,
			InternalOutputStatus: row.InternalOutputStatus,
			RowReflectionSummary: row.RowReflectionSummary,
			StaleReason:          row.StaleReason,
			CanRegenerate:        row.CanRegenerate,
		})
	}

	return TranslationOutputDiffPreviewResult{
		JobID:      readModel.JobID,
		ArtifactID: readModel.ArtifactID,
		Rows:       rows,
		CompatibilitySummary: TranslationOutputCompatibilitySummary{
			Passed:       readModel.CompatibilitySummary.Passed,
			WarningCount: readModel.CompatibilitySummary.WarningCount,
			RejectCount:  readModel.CompatibilitySummary.RejectCount,
		},
	}
}

func cloneStringIntMap(src map[string]int) map[string]int {
	if len(src) == 0 {
		return map[string]int{}
	}
	dst := make(map[string]int, len(src))
	for key, value := range src {
		dst[key] = value
	}
	return dst
}

func normalizeTranslationOutputArtifactErrorKind(kind string) TranslationOutputArtifactErrorKind {
	switch kind {
	case TranslationOutputArtifactErrorKindNotCompleted,
		TranslationOutputArtifactErrorKindCanceled,
		TranslationOutputArtifactErrorKindStatusMismatch,
		TranslationOutputArtifactErrorKindMissingRequiredRowField,
		TranslationOutputArtifactErrorKindUnknownOutputStatus,
		TranslationOutputArtifactErrorKindXMLSerializationFailed,
		TranslationOutputArtifactErrorKindFileWriteFailed,
		TranslationOutputArtifactErrorKindArtifactSaveFailed,
		TranslationOutputArtifactErrorKindCompatibilityRejected,
		TranslationOutputArtifactErrorKindSecretRedacted:
		return kind
	default:
		return ""
	}
}

func toTranslationOutputArtifactCommandResult(
	serviceResult service.TranslationOutputArtifactCommandResult,
	defaultOperationKind string,
	defaultReplacedArtifactID int64,
	err error,
) XTranslatorOutputArtifactCommandResult {
	result := XTranslatorOutputArtifactCommandResult{
		JobID:          serviceResult.Artifact.TranslationJobID,
		ArtifactID:     serviceResult.Artifact.ID,
		ArtifactStatus: serviceResult.Artifact.Status,
		RowCount:       serviceResult.RowCount,
		FilePath:       serviceResult.Artifact.FilePath,
		TargetGame:     serviceResult.Artifact.TargetGame,
		OperationSummary: TranslationOutputOperationSummary{
			OperationKind:       defaultOperationKind,
			ReplacedArtifactID:  defaultReplacedArtifactID,
			AffectedFieldIDs:    append([]int64(nil), serviceResult.AffectedFieldIDs...),
			DuplicateRowCreated: serviceResult.DuplicateRowCreated,
		},
	}
	if serviceResult.OperationKind != "" {
		result.OperationSummary.OperationKind = serviceResult.OperationKind
	}
	if serviceResult.ReplacedArtifactID != 0 {
		result.OperationSummary.ReplacedArtifactID = serviceResult.ReplacedArtifactID
	}
	if failure, ok := service.IsTranslationOutputArtifactFailure(err); ok {
		result.ArtifactStatus = failure.ArtifactStatus()
		result.ErrorSummary = &TranslationOutputArtifactErrorSummary{
			ErrorKind:  normalizeTranslationOutputArtifactErrorKind(failure.ErrorKind()),
			Reason:     failure.Reason(),
			Retryable:  failure.Retryable(),
			IsRedacted: true,
		}
	}
	return result
}
