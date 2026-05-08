package wails

import (
	"context"
	"fmt"

	"aitranslationenginejp/internal/usecase"
)

// TranslationOutputArtifactUsecasePort defines the frozen output artifact usecase seam.
type TranslationOutputArtifactUsecasePort interface {
	GetTranslationOutputReview(context.Context, usecase.GetTranslationOutputReviewRequest) (usecase.TranslationOutputReviewResult, error)
	GetTranslationOutputDiffPreview(context.Context, usecase.GetTranslationOutputDiffPreviewRequest) (usecase.TranslationOutputDiffPreviewResult, error)
	GenerateXTranslatorOutputArtifact(context.Context, usecase.GenerateXTranslatorOutputArtifactRequest) (usecase.XTranslatorOutputArtifactCommandResult, error)
	RegenerateXTranslatorOutputArtifact(context.Context, usecase.RegenerateXTranslatorOutputArtifactRequest) (usecase.XTranslatorOutputArtifactCommandResult, error)
}

// TranslationOutputArtifactController exposes Wails-bound output artifact entrypoints.
type TranslationOutputArtifactController struct {
	usecase TranslationOutputArtifactUsecasePort
}

// GetTranslationOutputReviewRequestDTO identifies the requested review target.
type GetTranslationOutputReviewRequestDTO struct {
	SelectedJobID *int64 `json:"selectedJobId,omitempty"`
}

// GetTranslationOutputDiffPreviewRequestDTO identifies the requested preview target.
type GetTranslationOutputDiffPreviewRequestDTO struct {
	JobID      int64 `json:"jobId"`
	ArtifactID int64 `json:"artifactId"`
}

// GenerateXTranslatorOutputArtifactRequestDTO identifies one generation request.
type GenerateXTranslatorOutputArtifactRequestDTO struct {
	JobID      int64  `json:"jobId"`
	TargetGame string `json:"targetGame"`
	OutputPath string `json:"outputPath"`
}

// RegenerateXTranslatorOutputArtifactRequestDTO identifies one regeneration request.
type RegenerateXTranslatorOutputArtifactRequestDTO struct {
	JobID      int64  `json:"jobId"`
	ArtifactID int64  `json:"artifactId"`
	TargetGame string `json:"targetGame"`
	OutputPath string `json:"outputPath"`
}

// TranslationOutputReviewResponseDTO returns the frozen Output Review response shape.
type TranslationOutputReviewResponseDTO struct {
	CompletedJobs    []TranslationOutputCompletedJobSummaryDTO  `json:"completedJobs"`
	HasSelectedJob   bool                                       `json:"hasSelectedJob"`
	SelectedJob      TranslationOutputSelectedJobSummaryDTO     `json:"selectedJob"`
	OutputReadiness  TranslationOutputReadinessSummaryDTO       `json:"outputReadiness"`
	ArtifactStatus   TranslationOutputArtifactStatusSummaryDTO  `json:"artifactStatus"`
	RejectionReasons []TranslationOutputArtifactErrorSummaryDTO `json:"rejectionReasons,omitempty"`
}

// TranslationOutputCompletedJobSummaryDTO summarizes one completed job candidate.
type TranslationOutputCompletedJobSummaryDTO struct {
	JobID                    int64          `json:"jobId"`
	JobStatus                string         `json:"jobStatus"`
	ArtifactStatus           string         `json:"artifactStatus"`
	OutputReady              bool           `json:"outputReady"`
	TranslatedCount          int            `json:"translatedCount"`
	OutputStatusDistribution map[string]int `json:"outputStatusDistribution,omitempty"`
}

// TranslationOutputSelectedJobSummaryDTO summarizes one selected job.
type TranslationOutputSelectedJobSummaryDTO struct {
	JobID           int64                             `json:"jobId"`
	JobStatus       string                            `json:"jobStatus"`
	BodyPhaseStatus string                            `json:"bodyPhaseStatus"`
	OutputReady     bool                              `json:"outputReady"`
	ResultSummary   TranslationOutputResultSummaryDTO `json:"resultSummary"`
}

// TranslationOutputResultSummaryDTO summarizes counts and provenance.
type TranslationOutputResultSummaryDTO struct {
	TranslatedCount int                                        `json:"translatedCount"`
	RowCount        int                                        `json:"rowCount"`
	InputProvenance TranslationOutputInputProvenanceSummaryDTO `json:"inputProvenance"`
}

// TranslationOutputInputProvenanceSummaryDTO summarizes provenance digests.
type TranslationOutputInputProvenanceSummaryDTO struct {
	InputSnapshotDigest string `json:"inputSnapshotDigest"`
	SourceFileDigest    string `json:"sourceFileDigest"`
}

// TranslationOutputReadinessSummaryDTO summarizes output readiness.
type TranslationOutputReadinessSummaryDTO struct {
	Ready         bool   `json:"ready"`
	Retryable     bool   `json:"retryable"`
	RejectionKind string `json:"rejectionKind,omitempty"`
}

// TranslationOutputArtifactStatusSummaryDTO summarizes artifact state.
type TranslationOutputArtifactStatusSummaryDTO struct {
	ArtifactID     int64  `json:"artifactId"`
	Status         string `json:"status"`
	RowCount       int    `json:"rowCount"`
	CurrentVersion bool   `json:"currentVersion"`
}

// TranslationOutputDiffPreviewResponseDTO returns the frozen diff preview response shape.
type TranslationOutputDiffPreviewResponseDTO struct {
	JobID                int64                                    `json:"jobId"`
	ArtifactID           int64                                    `json:"artifactId"`
	Rows                 []TranslationOutputDiffPreviewRowDTO     `json:"rows"`
	CompatibilitySummary TranslationOutputCompatibilitySummaryDTO `json:"compatibilitySummary"`
}

// TranslationOutputDiffPreviewRowDTO summarizes one preview row.
type TranslationOutputDiffPreviewRowDTO struct {
	FieldID              int64  `json:"fieldId"`
	RowDigest            string `json:"rowDigest"`
	EDID                 string `json:"edid"`
	REC                  string `json:"rec"`
	FIELD                string `json:"field"`
	FORMID               string `json:"formId"`
	SourceExcerpt        string `json:"sourceExcerpt"`
	DestExcerpt          string `json:"destExcerpt"`
	XTranslatorStatus    int    `json:"xTranslatorStatus"`
	InternalOutputStatus string `json:"internalOutputStatus"`
	RowReflectionSummary string `json:"rowReflectionSummary"`
	StaleReason          string `json:"staleReason,omitempty"`
	CanRegenerate        bool   `json:"canRegenerate"`
}

// TranslationOutputCompatibilitySummaryDTO summarizes compatibility counts.
type TranslationOutputCompatibilitySummaryDTO struct {
	Passed       bool `json:"passed"`
	WarningCount int  `json:"warningCount"`
	RejectCount  int  `json:"rejectCount"`
}

// TranslationOutputArtifactCommandResponseDTO returns the frozen command response shape.
type TranslationOutputArtifactCommandResponseDTO struct {
	JobID            int64                                     `json:"jobId"`
	ArtifactID       int64                                     `json:"artifactId"`
	ArtifactStatus   string                                    `json:"artifactStatus"`
	RowCount         int                                       `json:"rowCount"`
	FilePath         string                                    `json:"filePath,omitempty"`
	TargetGame       string                                    `json:"targetGame"`
	ErrorSummary     *TranslationOutputArtifactErrorSummaryDTO `json:"errorSummary,omitempty"`
	OperationSummary TranslationOutputOperationSummaryDTO      `json:"operationSummary"`
}

// TranslationOutputArtifactErrorSummaryDTO exposes only redacted error information.
type TranslationOutputArtifactErrorSummaryDTO struct {
	ErrorKind  string `json:"errorKind"`
	Reason     string `json:"reason"`
	Retryable  bool   `json:"retryable"`
	IsRedacted bool   `json:"isRedacted"`
}

// TranslationOutputOperationSummaryDTO summarizes write effects.
type TranslationOutputOperationSummaryDTO struct {
	OperationKind       string  `json:"operationKind"`
	ReplacedArtifactID  int64   `json:"replacedArtifactId"`
	AffectedFieldIDs    []int64 `json:"affectedFieldIds,omitempty"`
	DuplicateRowCreated bool    `json:"duplicateRowCreated"`
}

// NewTranslationOutputArtifactController creates an output artifact controller.
func NewTranslationOutputArtifactController(usecasePort TranslationOutputArtifactUsecasePort) *TranslationOutputArtifactController {
	return &TranslationOutputArtifactController{usecase: usecasePort}
}

// GetTranslationOutputReview returns the frozen Output Review contract.
func (controller *TranslationOutputArtifactController) GetTranslationOutputReview(
	request GetTranslationOutputReviewRequestDTO,
) (TranslationOutputReviewResponseDTO, error) {
	result, err := controller.usecase.GetTranslationOutputReview(
		context.Background(),
		usecase.GetTranslationOutputReviewRequest{SelectedJobID: request.SelectedJobID},
	)
	if err != nil {
		return TranslationOutputReviewResponseDTO{}, fmt.Errorf("get translation output review: %w", err)
	}
	return toTranslationOutputReviewResponseDTO(result), nil
}

// GetTranslationOutputDiffPreview returns the frozen diff preview contract.
func (controller *TranslationOutputArtifactController) GetTranslationOutputDiffPreview(
	request GetTranslationOutputDiffPreviewRequestDTO,
) (TranslationOutputDiffPreviewResponseDTO, error) {
	result, err := controller.usecase.GetTranslationOutputDiffPreview(
		context.Background(),
		usecase.GetTranslationOutputDiffPreviewRequest{
			JobID:      request.JobID,
			ArtifactID: request.ArtifactID,
		},
	)
	if err != nil {
		return TranslationOutputDiffPreviewResponseDTO{}, fmt.Errorf("get translation output diff preview: %w", err)
	}
	return toTranslationOutputDiffPreviewResponseDTO(result), nil
}

// GenerateXTranslatorOutputArtifact returns the frozen artifact generation contract.
func (controller *TranslationOutputArtifactController) GenerateXTranslatorOutputArtifact(
	request GenerateXTranslatorOutputArtifactRequestDTO,
) (TranslationOutputArtifactCommandResponseDTO, error) {
	result, err := controller.usecase.GenerateXTranslatorOutputArtifact(
		context.Background(),
		usecase.GenerateXTranslatorOutputArtifactRequest{
			JobID:      request.JobID,
			TargetGame: request.TargetGame,
			OutputPath: request.OutputPath,
		},
	)
	if err != nil {
		return TranslationOutputArtifactCommandResponseDTO{}, fmt.Errorf("generate xtranslator output artifact: %w", err)
	}
	return toTranslationOutputArtifactCommandResponseDTO(result), nil
}

// RegenerateXTranslatorOutputArtifact returns the frozen artifact regeneration contract.
func (controller *TranslationOutputArtifactController) RegenerateXTranslatorOutputArtifact(
	request RegenerateXTranslatorOutputArtifactRequestDTO,
) (TranslationOutputArtifactCommandResponseDTO, error) {
	result, err := controller.usecase.RegenerateXTranslatorOutputArtifact(
		context.Background(),
		usecase.RegenerateXTranslatorOutputArtifactRequest{
			JobID:      request.JobID,
			ArtifactID: request.ArtifactID,
			TargetGame: request.TargetGame,
			OutputPath: request.OutputPath,
		},
	)
	if err != nil {
		return TranslationOutputArtifactCommandResponseDTO{}, fmt.Errorf("regenerate xtranslator output artifact: %w", err)
	}
	return toTranslationOutputArtifactCommandResponseDTO(result), nil
}

func toTranslationOutputReviewResponseDTO(result usecase.TranslationOutputReviewResult) TranslationOutputReviewResponseDTO {
	completedJobs := make([]TranslationOutputCompletedJobSummaryDTO, 0, len(result.CompletedJobs))
	for _, job := range result.CompletedJobs {
		completedJobs = append(completedJobs, TranslationOutputCompletedJobSummaryDTO{
			JobID:                    job.JobID,
			JobStatus:                job.JobStatus,
			ArtifactStatus:           job.ArtifactStatus,
			OutputReady:              job.OutputReady,
			TranslatedCount:          job.TranslatedCount,
			OutputStatusDistribution: job.OutputStatusDistribution,
		})
	}

	rejectionReasons := make([]TranslationOutputArtifactErrorSummaryDTO, 0, len(result.RejectionReasons))
	for _, reason := range result.RejectionReasons {
		rejectionReasons = append(rejectionReasons, toTranslationOutputArtifactErrorSummaryDTO(&reason))
	}

	return TranslationOutputReviewResponseDTO{
		CompletedJobs:  completedJobs,
		HasSelectedJob: result.HasSelectedJob,
		SelectedJob: TranslationOutputSelectedJobSummaryDTO{
			JobID:           result.SelectedJob.JobID,
			JobStatus:       result.SelectedJob.JobStatus,
			BodyPhaseStatus: result.SelectedJob.BodyPhaseStatus,
			OutputReady:     result.SelectedJob.OutputReady,
			ResultSummary: TranslationOutputResultSummaryDTO{
				TranslatedCount: result.SelectedJob.ResultSummary.TranslatedCount,
				RowCount:        result.SelectedJob.ResultSummary.RowCount,
				InputProvenance: TranslationOutputInputProvenanceSummaryDTO{
					InputSnapshotDigest: result.SelectedJob.ResultSummary.InputProvenance.InputSnapshotDigest,
					SourceFileDigest:    result.SelectedJob.ResultSummary.InputProvenance.SourceFileDigest,
				},
			},
		},
		OutputReadiness: TranslationOutputReadinessSummaryDTO{
			Ready:         result.OutputReadiness.Ready,
			Retryable:     result.OutputReadiness.Retryable,
			RejectionKind: result.OutputReadiness.RejectionKind,
		},
		ArtifactStatus: TranslationOutputArtifactStatusSummaryDTO{
			ArtifactID:     result.ArtifactStatus.ArtifactID,
			Status:         result.ArtifactStatus.Status,
			RowCount:       result.ArtifactStatus.RowCount,
			CurrentVersion: result.ArtifactStatus.CurrentVersion,
		},
		RejectionReasons: rejectionReasons,
	}
}

func toTranslationOutputDiffPreviewResponseDTO(result usecase.TranslationOutputDiffPreviewResult) TranslationOutputDiffPreviewResponseDTO {
	rows := make([]TranslationOutputDiffPreviewRowDTO, 0, len(result.Rows))
	for _, row := range result.Rows {
		rows = append(rows, TranslationOutputDiffPreviewRowDTO{
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

	return TranslationOutputDiffPreviewResponseDTO{
		JobID:      result.JobID,
		ArtifactID: result.ArtifactID,
		Rows:       rows,
		CompatibilitySummary: TranslationOutputCompatibilitySummaryDTO{
			Passed:       result.CompatibilitySummary.Passed,
			WarningCount: result.CompatibilitySummary.WarningCount,
			RejectCount:  result.CompatibilitySummary.RejectCount,
		},
	}
}

func toTranslationOutputArtifactCommandResponseDTO(result usecase.XTranslatorOutputArtifactCommandResult) TranslationOutputArtifactCommandResponseDTO {
	return TranslationOutputArtifactCommandResponseDTO{
		JobID:            result.JobID,
		ArtifactID:       result.ArtifactID,
		ArtifactStatus:   result.ArtifactStatus,
		RowCount:         result.RowCount,
		FilePath:         result.FilePath,
		TargetGame:       result.TargetGame,
		ErrorSummary:     toTranslationOutputArtifactErrorSummaryPtrDTO(result.ErrorSummary),
		OperationSummary: toTranslationOutputOperationSummaryDTO(result.OperationSummary),
	}
}

func toTranslationOutputArtifactErrorSummaryPtrDTO(summary *usecase.TranslationOutputArtifactErrorSummary) *TranslationOutputArtifactErrorSummaryDTO {
	if summary == nil {
		return nil
	}

	dto := toTranslationOutputArtifactErrorSummaryDTO(summary)
	return &dto
}

func toTranslationOutputArtifactErrorSummaryDTO(summary *usecase.TranslationOutputArtifactErrorSummary) TranslationOutputArtifactErrorSummaryDTO {
	return TranslationOutputArtifactErrorSummaryDTO{
		ErrorKind:  summary.ErrorKind,
		Reason:     summary.Reason,
		Retryable:  summary.Retryable,
		IsRedacted: summary.IsRedacted,
	}
}

func toTranslationOutputOperationSummaryDTO(summary usecase.TranslationOutputOperationSummary) TranslationOutputOperationSummaryDTO {
	return TranslationOutputOperationSummaryDTO{
		OperationKind:       summary.OperationKind,
		ReplacedArtifactID:  summary.ReplacedArtifactID,
		AffectedFieldIDs:    summary.AffectedFieldIDs,
		DuplicateRowCreated: summary.DuplicateRowCreated,
	}
}
