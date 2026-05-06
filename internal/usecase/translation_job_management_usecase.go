package usecase

import (
	"context"
	"fmt"

	jobmanagementservice "aitranslationenginejp/internal/service"
)

type translationJobManagementServicePort interface {
	ListIncompleteJobs(ctx context.Context) (jobmanagementservice.TranslationJobManagementListReadModel, error)
	GetJobDetail(ctx context.Context, jobID int64) (jobmanagementservice.TranslationJobManagementJobDetailReadModel, error)
	DeleteJob(ctx context.Context, jobID int64) (jobmanagementservice.TranslationJobManagementActionReadModel, error)
	RequestStop(ctx context.Context, jobID int64) (jobmanagementservice.TranslationJobManagementActionReadModel, error)
	ResumeJob(ctx context.Context, jobID int64) (jobmanagementservice.TranslationJobManagementActionReadModel, error)
}

// TranslationJobManagementUsecase adapts the service layer to the Wails-facing contract.
type TranslationJobManagementUsecase struct {
	service translationJobManagementServicePort
}

// NewTranslationJobManagementUsecase builds the Job Management usecase seam.
func NewTranslationJobManagementUsecase(service translationJobManagementServicePort) *TranslationJobManagementUsecase {
	return &TranslationJobManagementUsecase{service: service}
}

// ListIncompleteJobs returns non-completed jobs for the management screen.
func (usecase *TranslationJobManagementUsecase) ListIncompleteJobs(
	ctx context.Context,
) (TranslationJobManagementListResult, error) {
	result, err := usecase.service.ListIncompleteJobs(ctx)
	if err != nil {
		return TranslationJobManagementListResult{}, fmt.Errorf("list incomplete translation jobs: %w", err)
	}
	return TranslationJobManagementListResult{Jobs: toTranslationJobManagementJobSummaries(result.Jobs)}, nil
}

// GetJobDetail returns one selected incomplete job detail.
func (usecase *TranslationJobManagementUsecase) GetJobDetail(
	ctx context.Context,
	request TranslationJobManagementGetDetailRequest,
) (TranslationJobManagementJobDetail, error) {
	result, err := usecase.service.GetJobDetail(ctx, request.JobID)
	if err != nil {
		return TranslationJobManagementJobDetail{}, fmt.Errorf("get translation job management detail: %w", err)
	}
	return toTranslationJobManagementJobDetail(result), nil
}

// DeleteJob deletes one non-running job and returns the public outcome.
func (usecase *TranslationJobManagementUsecase) DeleteJob(
	ctx context.Context,
	request TranslationJobManagementDeleteRequest,
) (TranslationJobManagementActionResult, error) {
	result, err := usecase.service.DeleteJob(ctx, request.JobID)
	if err != nil {
		return TranslationJobManagementActionResult{}, fmt.Errorf("delete translation job management job: %w", err)
	}
	return toTranslationJobManagementActionResult(result), nil
}

// RequestStop returns the current stop-facing projection without mutating execution state.
func (usecase *TranslationJobManagementUsecase) RequestStop(
	ctx context.Context,
	request TranslationJobManagementActionRequest,
) (TranslationJobManagementActionResult, error) {
	result, err := usecase.service.RequestStop(ctx, request.JobID)
	if err != nil {
		return TranslationJobManagementActionResult{}, fmt.Errorf("request translation job stop: %w", err)
	}
	return toTranslationJobManagementActionResult(result), nil
}

// ResumeJob returns the current resume-facing projection without mutating execution state.
func (usecase *TranslationJobManagementUsecase) ResumeJob(
	ctx context.Context,
	request TranslationJobManagementActionRequest,
) (TranslationJobManagementActionResult, error) {
	result, err := usecase.service.ResumeJob(ctx, request.JobID)
	if err != nil {
		return TranslationJobManagementActionResult{}, fmt.Errorf("resume translation job: %w", err)
	}
	return toTranslationJobManagementActionResult(result), nil
}

func toTranslationJobManagementJobSummaries(
	source []jobmanagementservice.TranslationJobManagementJobSummaryReadModel,
) []TranslationJobManagementJobSummary {
	result := make([]TranslationJobManagementJobSummary, 0, len(source))
	for _, item := range source {
		result = append(result, TranslationJobManagementJobSummary{
			JobID:              item.JobID,
			JobState:           item.JobState,
			JobStateLabel:      item.JobStateLabel,
			StateTone:          item.StateTone,
			InputSource:        toTranslationJobManagementInputSourceSummary(item.InputSource),
			Progress:           toTranslationJobManagementProgressSummary(item.Progress),
			StopAvailability:   toTranslationJobManagementOperationAvailability(item.StopAvailability),
			ResumeAvailability: toTranslationJobManagementOperationAvailability(item.ResumeAvailability),
			DeleteAvailability: toTranslationJobManagementOperationAvailability(item.DeleteAvailability),
		})
	}
	return result
}

func toTranslationJobManagementJobDetail(
	source jobmanagementservice.TranslationJobManagementJobDetailReadModel,
) TranslationJobManagementJobDetail {
	return TranslationJobManagementJobDetail{
		TranslationJobManagementJobSummary: TranslationJobManagementJobSummary{
			JobID:              source.JobID,
			JobState:           source.JobState,
			JobStateLabel:      source.JobStateLabel,
			StateTone:          source.StateTone,
			InputSource:        toTranslationJobManagementInputSourceSummary(source.InputSource),
			Progress:           toTranslationJobManagementProgressSummary(source.Progress),
			StopAvailability:   toTranslationJobManagementOperationAvailability(source.StopAvailability),
			ResumeAvailability: toTranslationJobManagementOperationAvailability(source.ResumeAvailability),
			DeleteAvailability: toTranslationJobManagementOperationAvailability(source.DeleteAvailability),
		},
		CacheState:           source.CacheState,
		CacheStateLabel:      source.CacheStateLabel,
		RuntimeSummary:       toTranslationJobManagementProtectedSettingSummary(source.RuntimeSummary),
		ResumeBlockedReasons: toTranslationJobManagementBlockedReasons(source.ResumeBlockedReasons),
		Warnings:             toTranslationJobManagementBlockedReasons(source.Warnings),
		DeleteImpactLines:    append([]string(nil), source.DeleteImpactLines...),
	}
}

func toTranslationJobManagementActionResult(
	source jobmanagementservice.TranslationJobManagementActionReadModel,
) TranslationJobManagementActionResult {
	result := TranslationJobManagementActionResult{
		Message:        source.Message,
		Tone:           source.Tone,
		ReasonCategory: source.ReasonCategory,
	}
	if source.Detail != nil {
		detail := toTranslationJobManagementJobDetail(*source.Detail)
		result.Detail = &detail
	}
	if source.DeletedJobID != nil {
		deletedJobID := *source.DeletedJobID
		result.DeletedJobID = &deletedJobID
	}
	return result
}

func toTranslationJobManagementInputSourceSummary(
	source jobmanagementservice.TranslationJobManagementInputSourceSummaryReadModel,
) TranslationJobManagementInputSourceSummary {
	return TranslationJobManagementInputSourceSummary{
		InputSourceID:        source.InputSourceID,
		InputSourceLabel:     source.InputSourceLabel,
		InputSourceKindLabel: source.InputSourceKindLabel,
		SourcePath:           source.SourcePath,
		PluginName:           source.PluginName,
		ExtractedJSONLabel:   source.ExtractedJSONLabel,
	}
}

func toTranslationJobManagementProgressSummary(
	source jobmanagementservice.TranslationJobManagementProgressSummaryReadModel,
) TranslationJobManagementProgressSummary {
	return TranslationJobManagementProgressSummary{
		CurrentPhaseLabel: source.CurrentPhaseLabel,
		Percent:           source.Percent,
		ProgressLabel:     source.ProgressLabel,
		LastUpdatedLabel:  source.LastUpdatedLabel,
	}
}

func toTranslationJobManagementProtectedSettingSummary(
	source jobmanagementservice.TranslationJobManagementProtectedSettingSummaryReadModel,
) TranslationJobManagementProtectedSettingSummary {
	return TranslationJobManagementProtectedSettingSummary{
		ProviderLabel:        source.ProviderLabel,
		ModelLabel:           source.ModelLabel,
		ExecutionModeLabel:   source.ExecutionModeLabel,
		CredentialState:      source.CredentialState,
		CredentialStateLabel: source.CredentialStateLabel,
	}
}

func toTranslationJobManagementOperationAvailability(
	source jobmanagementservice.TranslationJobManagementOperationAvailabilityReadModel,
) TranslationJobManagementOperationAvailability {
	return TranslationJobManagementOperationAvailability{
		Kind:           source.Kind,
		Enabled:        source.Enabled,
		Label:          source.Label,
		HelperText:     source.HelperText,
		ReasonCategory: source.ReasonCategory,
		ReasonText:     source.ReasonText,
	}
}

func toTranslationJobManagementBlockedReasons(
	source []jobmanagementservice.TranslationJobManagementBlockedReasonReadModel,
) []TranslationJobManagementBlockedReason {
	result := make([]TranslationJobManagementBlockedReason, 0, len(source))
	for _, item := range source {
		result = append(result, TranslationJobManagementBlockedReason{
			Category: item.Category,
			Title:    item.Title,
			Detail:   item.Detail,
		})
	}
	return result
}
