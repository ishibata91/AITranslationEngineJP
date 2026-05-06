package wails

import (
	"context"
	"fmt"

	"aitranslationenginejp/internal/usecase"
)

// TranslationJobManagementUsecasePort defines the frozen Job Management usecase seam.
type TranslationJobManagementUsecasePort interface {
	ListIncompleteJobs(ctx context.Context) (usecase.TranslationJobManagementListResult, error)
	GetJobDetail(ctx context.Context, request usecase.TranslationJobManagementGetDetailRequest) (usecase.TranslationJobManagementJobDetail, error)
	DeleteJob(ctx context.Context, request usecase.TranslationJobManagementDeleteRequest) (usecase.TranslationJobManagementActionResult, error)
	RequestStop(ctx context.Context, request usecase.TranslationJobManagementActionRequest) (usecase.TranslationJobManagementActionResult, error)
	ResumeJob(ctx context.Context, request usecase.TranslationJobManagementActionRequest) (usecase.TranslationJobManagementActionResult, error)
}

// TranslationJobManagementController exposes Wails-bound Job Management entrypoints.
type TranslationJobManagementController struct {
	usecase TranslationJobManagementUsecasePort
}

// TranslationJobManagementInputSourceSummaryDTO returns one input source summary.
type TranslationJobManagementInputSourceSummaryDTO struct {
	InputSourceID        int64  `json:"inputSourceId"`
	InputSourceLabel     string `json:"inputSourceLabel"`
	InputSourceKindLabel string `json:"inputSourceKindLabel"`
	SourcePath           string `json:"sourcePath"`
	PluginName           string `json:"pluginName"`
	ExtractedJSONLabel   string `json:"extractedJsonLabel"`
}

// TranslationJobManagementProgressSummaryDTO returns one progress summary.
type TranslationJobManagementProgressSummaryDTO struct {
	CurrentPhaseLabel string `json:"currentPhaseLabel"`
	Percent           int    `json:"percent"`
	ProgressLabel     string `json:"progressLabel"`
	LastUpdatedLabel  string `json:"lastUpdatedLabel"`
}

// TranslationJobManagementProtectedSettingSummaryDTO returns redacted runtime settings.
type TranslationJobManagementProtectedSettingSummaryDTO struct {
	ProviderLabel        string `json:"providerLabel"`
	ModelLabel           string `json:"modelLabel"`
	ExecutionModeLabel   string `json:"executionModeLabel"`
	CredentialState      string `json:"credentialState"`
	CredentialStateLabel string `json:"credentialStateLabel"`
}

// TranslationJobManagementOperationAvailabilityDTO returns one operation guard.
type TranslationJobManagementOperationAvailabilityDTO struct {
	Kind           string `json:"kind"`
	Enabled        bool   `json:"enabled"`
	Label          string `json:"label"`
	HelperText     string `json:"helperText"`
	ReasonCategory string `json:"reasonCategory,omitempty"`
	ReasonText     string `json:"reasonText,omitempty"`
}

// TranslationJobManagementBlockedReasonDTO returns one blocked or warning reason.
type TranslationJobManagementBlockedReasonDTO struct {
	Category string `json:"category"`
	Title    string `json:"title"`
	Detail   string `json:"detail"`
}

// TranslationJobManagementJobSummaryDTO returns one list item.
type TranslationJobManagementJobSummaryDTO struct {
	JobID              int64                                            `json:"jobId"`
	JobState           string                                           `json:"jobState"`
	JobStateLabel      string                                           `json:"jobStateLabel"`
	StateTone          string                                           `json:"stateTone"`
	InputSource        TranslationJobManagementInputSourceSummaryDTO    `json:"inputSource"`
	Progress           TranslationJobManagementProgressSummaryDTO       `json:"progress"`
	StopAvailability   TranslationJobManagementOperationAvailabilityDTO `json:"stopAvailability"`
	ResumeAvailability TranslationJobManagementOperationAvailabilityDTO `json:"resumeAvailability"`
	DeleteAvailability TranslationJobManagementOperationAvailabilityDTO `json:"deleteAvailability"`
}

// TranslationJobManagementJobDetailDTO returns one selected job detail.
type TranslationJobManagementJobDetailDTO struct {
	TranslationJobManagementJobSummaryDTO
	CacheState           string                                             `json:"cacheState"`
	CacheStateLabel      string                                             `json:"cacheStateLabel"`
	RuntimeSummary       TranslationJobManagementProtectedSettingSummaryDTO `json:"runtimeSummary"`
	ResumeBlockedReasons []TranslationJobManagementBlockedReasonDTO         `json:"resumeBlockedReasons"`
	Warnings             []TranslationJobManagementBlockedReasonDTO         `json:"warnings"`
	DeleteImpactLines    []string                                           `json:"deleteImpactLines"`
}

// TranslationJobManagementListResponseDTO returns the incomplete-job list payload.
type TranslationJobManagementListResponseDTO struct {
	Jobs []TranslationJobManagementJobSummaryDTO `json:"jobs"`
}

// TranslationJobManagementGetDetailRequestDTO identifies one detail request.
type TranslationJobManagementGetDetailRequestDTO struct {
	JobID int64 `json:"jobId"`
}

// TranslationJobManagementDeleteRequestDTO identifies one delete request.
type TranslationJobManagementDeleteRequestDTO struct {
	JobID int64 `json:"jobId"`
}

// TranslationJobManagementActionRequestDTO identifies one stop or resume request.
type TranslationJobManagementActionRequestDTO struct {
	JobID int64 `json:"jobId"`
}

// TranslationJobManagementActionResponseDTO returns one action outcome.
type TranslationJobManagementActionResponseDTO struct {
	Message        string                                `json:"message"`
	Tone           string                                `json:"tone"`
	Detail         *TranslationJobManagementJobDetailDTO `json:"detail,omitempty"`
	DeletedJobID   *int64                                `json:"deletedJobId,omitempty"`
	ReasonCategory string                                `json:"reasonCategory,omitempty"`
}

// NewTranslationJobManagementController creates a Job Management controller.
func NewTranslationJobManagementController(
	usecase TranslationJobManagementUsecasePort,
) *TranslationJobManagementController {
	return &TranslationJobManagementController{usecase: usecase}
}

// ListIncompleteJobs returns non-completed jobs for the management screen.
func (controller *TranslationJobManagementController) ListIncompleteJobs() (TranslationJobManagementListResponseDTO, error) {
	result, err := controller.usecase.ListIncompleteJobs(context.Background())
	if err != nil {
		return TranslationJobManagementListResponseDTO{}, fmt.Errorf("list translation job management jobs: %w", err)
	}
	jobs := make([]TranslationJobManagementJobSummaryDTO, 0, len(result.Jobs))
	for _, item := range result.Jobs {
		jobs = append(jobs, toTranslationJobManagementJobSummaryDTO(item))
	}
	return TranslationJobManagementListResponseDTO{Jobs: jobs}, nil
}

// GetJobDetail returns one selected incomplete job detail.
func (controller *TranslationJobManagementController) GetJobDetail(
	request TranslationJobManagementGetDetailRequestDTO,
) (TranslationJobManagementJobDetailDTO, error) {
	result, err := controller.usecase.GetJobDetail(context.Background(), usecase.TranslationJobManagementGetDetailRequest{
		JobID: request.JobID,
	})
	if err != nil {
		return TranslationJobManagementJobDetailDTO{}, fmt.Errorf("get translation job management detail: %w", err)
	}
	return toTranslationJobManagementJobDetailDTO(result), nil
}

// DeleteJob deletes one non-running job and returns the public outcome.
func (controller *TranslationJobManagementController) DeleteJob(
	request TranslationJobManagementDeleteRequestDTO,
) (TranslationJobManagementActionResponseDTO, error) {
	result, err := controller.usecase.DeleteJob(context.Background(), usecase.TranslationJobManagementDeleteRequest{
		JobID: request.JobID,
	})
	if err != nil {
		return TranslationJobManagementActionResponseDTO{}, fmt.Errorf("delete translation job management job: %w", err)
	}
	return toTranslationJobManagementActionResponseDTO(result), nil
}

// RequestStop returns the stop-facing projection without mutating execution state.
func (controller *TranslationJobManagementController) RequestStop(
	request TranslationJobManagementActionRequestDTO,
) (TranslationJobManagementActionResponseDTO, error) {
	result, err := controller.usecase.RequestStop(context.Background(), usecase.TranslationJobManagementActionRequest{
		JobID: request.JobID,
	})
	if err != nil {
		return TranslationJobManagementActionResponseDTO{}, fmt.Errorf("request translation job stop: %w", err)
	}
	return toTranslationJobManagementActionResponseDTO(result), nil
}

// ResumeJob returns the resume-facing projection without mutating execution state.
func (controller *TranslationJobManagementController) ResumeJob(
	request TranslationJobManagementActionRequestDTO,
) (TranslationJobManagementActionResponseDTO, error) {
	result, err := controller.usecase.ResumeJob(context.Background(), usecase.TranslationJobManagementActionRequest{
		JobID: request.JobID,
	})
	if err != nil {
		return TranslationJobManagementActionResponseDTO{}, fmt.Errorf("resume translation job: %w", err)
	}
	return toTranslationJobManagementActionResponseDTO(result), nil
}

func toTranslationJobManagementJobSummaryDTO(
	source usecase.TranslationJobManagementJobSummary,
) TranslationJobManagementJobSummaryDTO {
	return TranslationJobManagementJobSummaryDTO{
		JobID:              source.JobID,
		JobState:           source.JobState,
		JobStateLabel:      source.JobStateLabel,
		StateTone:          source.StateTone,
		InputSource:        toTranslationJobManagementInputSourceSummaryDTO(source.InputSource),
		Progress:           toTranslationJobManagementProgressSummaryDTO(source.Progress),
		StopAvailability:   toTranslationJobManagementOperationAvailabilityDTO(source.StopAvailability),
		ResumeAvailability: toTranslationJobManagementOperationAvailabilityDTO(source.ResumeAvailability),
		DeleteAvailability: toTranslationJobManagementOperationAvailabilityDTO(source.DeleteAvailability),
	}
}

func toTranslationJobManagementJobDetailDTO(
	source usecase.TranslationJobManagementJobDetail,
) TranslationJobManagementJobDetailDTO {
	return TranslationJobManagementJobDetailDTO{
		TranslationJobManagementJobSummaryDTO: toTranslationJobManagementJobSummaryDTO(source.TranslationJobManagementJobSummary),
		CacheState:                            source.CacheState,
		CacheStateLabel:                       source.CacheStateLabel,
		RuntimeSummary:                        toTranslationJobManagementProtectedSettingSummaryDTO(source.RuntimeSummary),
		ResumeBlockedReasons:                  toTranslationJobManagementBlockedReasonDTOs(source.ResumeBlockedReasons),
		Warnings:                              toTranslationJobManagementBlockedReasonDTOs(source.Warnings),
		DeleteImpactLines:                     append([]string(nil), source.DeleteImpactLines...),
	}
}

func toTranslationJobManagementActionResponseDTO(
	source usecase.TranslationJobManagementActionResult,
) TranslationJobManagementActionResponseDTO {
	result := TranslationJobManagementActionResponseDTO{
		Message:        source.Message,
		Tone:           source.Tone,
		DeletedJobID:   source.DeletedJobID,
		ReasonCategory: source.ReasonCategory,
	}
	if source.Detail != nil {
		detail := toTranslationJobManagementJobDetailDTO(*source.Detail)
		result.Detail = &detail
	}
	return result
}

func toTranslationJobManagementInputSourceSummaryDTO(
	source usecase.TranslationJobManagementInputSourceSummary,
) TranslationJobManagementInputSourceSummaryDTO {
	return TranslationJobManagementInputSourceSummaryDTO(source)
}

func toTranslationJobManagementProgressSummaryDTO(
	source usecase.TranslationJobManagementProgressSummary,
) TranslationJobManagementProgressSummaryDTO {
	return TranslationJobManagementProgressSummaryDTO(source)
}

func toTranslationJobManagementProtectedSettingSummaryDTO(
	source usecase.TranslationJobManagementProtectedSettingSummary,
) TranslationJobManagementProtectedSettingSummaryDTO {
	return TranslationJobManagementProtectedSettingSummaryDTO(source)
}

func toTranslationJobManagementOperationAvailabilityDTO(
	source usecase.TranslationJobManagementOperationAvailability,
) TranslationJobManagementOperationAvailabilityDTO {
	return TranslationJobManagementOperationAvailabilityDTO(source)
}

func toTranslationJobManagementBlockedReasonDTOs(
	source []usecase.TranslationJobManagementBlockedReason,
) []TranslationJobManagementBlockedReasonDTO {
	result := make([]TranslationJobManagementBlockedReasonDTO, 0, len(source))
	for _, item := range source {
		result = append(result, TranslationJobManagementBlockedReasonDTO(item))
	}
	return result
}
