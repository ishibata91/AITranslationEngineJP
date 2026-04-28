package wails

import (
	"context"
	"fmt"
	"time"

	"aitranslationenginejp/internal/usecase"
)

// TranslationJobSetupUsecasePort defines the frozen Job Setup usecase seam.
type TranslationJobSetupUsecasePort interface {
	GetTranslationJobSetupOptions(ctx context.Context) (usecase.TranslationJobSetupOptionsResult, error)
	ValidateTranslationJobSetup(ctx context.Context, request usecase.ValidateTranslationJobSetupRequest) (usecase.TranslationJobSetupValidationResult, error)
	CreateTranslationJob(ctx context.Context, request usecase.CreateTranslationJobRequest) (usecase.CreateTranslationJobResult, error)
	GetTranslationJobSetupSummary(ctx context.Context, request usecase.GetTranslationJobSetupSummaryRequest) (usecase.TranslationJobSetupSummaryResult, error)
}

// TranslationJobSetupController exposes Wails-bound Job Setup entrypoints.
type TranslationJobSetupController struct {
	translationJobSetupUsecase TranslationJobSetupUsecasePort
}

// TranslationJobSetupInputCandidateDTO is one selectable translation input candidate.
type TranslationJobSetupInputCandidateDTO struct {
	ID           int64  `json:"id"`
	Label        string `json:"label"`
	SourceKind   string `json:"sourceKind"`
	RecordCount  int    `json:"recordCount"`
	RegisteredAt string `json:"registeredAt"`
}

// TranslationJobSetupExistingJobDTO summarizes one already prepared job.
type TranslationJobSetupExistingJobDTO struct {
	InputSourceID int64  `json:"inputSourceId"`
	JobID         int64  `json:"jobId"`
	Status        string `json:"status"`
	InputSource   string `json:"inputSource"`
}

// TranslationJobSetupDictionaryOptionDTO is one shared dictionary option.
type TranslationJobSetupDictionaryOptionDTO struct {
	ID    string `json:"id"`
	Label string `json:"label"`
}

// TranslationJobSetupPersonaOptionDTO is one shared persona option.
type TranslationJobSetupPersonaOptionDTO struct {
	ID    string `json:"id"`
	Label string `json:"label"`
}

// TranslationJobSetupRuntimeOptionDTO is one selectable runtime option.
type TranslationJobSetupRuntimeOptionDTO struct {
	Provider string `json:"provider"`
	Model    string `json:"model"`
	Mode     string `json:"mode"`
}

// TranslationJobSetupCredentialReferenceDTO exposes only credential reference state.
type TranslationJobSetupCredentialReferenceDTO struct {
	Provider        string  `json:"provider"`
	CredentialRef   string  `json:"credentialRef"`
	IsConfigured    bool    `json:"isConfigured"`
	IsMissingSecret bool    `json:"isMissingSecret"`
	SecretPlaintext *string `json:"-"`
}

// TranslationJobSetupOptionsResponseDTO returns the read-only setup options.
type TranslationJobSetupOptionsResponseDTO struct {
	InputCandidates    []TranslationJobSetupInputCandidateDTO      `json:"inputCandidates"`
	ExistingJob        *TranslationJobSetupExistingJobDTO          `json:"existingJob,omitempty"`
	SharedDictionaries []TranslationJobSetupDictionaryOptionDTO    `json:"sharedDictionaries"`
	SharedPersonas     []TranslationJobSetupPersonaOptionDTO       `json:"sharedPersonas"`
	AIRuntimeOptions   []TranslationJobSetupRuntimeOptionDTO       `json:"aiRuntimeOptions"`
	CredentialRefs     []TranslationJobSetupCredentialReferenceDTO `json:"credentialRefs"`
}

// TranslationJobSetupRuntimeSelectionDTO carries the runtime selection for validation and create.
type TranslationJobSetupRuntimeSelectionDTO struct {
	Provider      string `json:"provider"`
	Model         string `json:"model"`
	ExecutionMode string `json:"executionMode"`
}

// ValidateTranslationJobSetupRequestDTO carries the frozen validation request payload.
type ValidateTranslationJobSetupRequestDTO struct {
	InputSourceID int64                                  `json:"inputSourceId"`
	Runtime       TranslationJobSetupRuntimeSelectionDTO `json:"runtime"`
	CredentialRef string                                 `json:"credentialRef"`
}

// TranslationJobSetupValidationResponseDTO returns the frozen validation result shape.
type TranslationJobSetupValidationResponseDTO struct {
	Status                  string   `json:"status"`
	BlockingFailureCategory *string  `json:"blockingFailureCategory,omitempty"`
	TargetSlices            []string `json:"targetSlices"`
	ValidatedAt             string   `json:"validatedAt"`
	CanCreate               bool     `json:"canCreate"`
	PassSlices              []string `json:"passSlices"`
}

// CreateTranslationJobRequestDTO carries the frozen create request payload.
type CreateTranslationJobRequestDTO struct {
	InputSourceID        int64                                  `json:"inputSourceId"`
	InputSource          string                                 `json:"inputSource"`
	ValidationStatus     string                                 `json:"validationStatus"`
	ValidatedAt          string                                 `json:"validatedAt"`
	ValidationPassSlices []string                               `json:"validationPassSlices"`
	Runtime              TranslationJobSetupRuntimeSelectionDTO `json:"runtime"`
	CredentialRef        string                                 `json:"credentialRef"`
}

// TranslationJobExecutionSummaryDTO returns the runtime summary captured by a created job.
type TranslationJobExecutionSummaryDTO struct {
	Provider      string `json:"provider"`
	Model         string `json:"model"`
	ExecutionMode string `json:"executionMode"`
}

// CreateTranslationJobResponseDTO returns either a ready job or a rejected error kind.
type CreateTranslationJobResponseDTO struct {
	JobID                int64                              `json:"jobId"`
	JobState             string                             `json:"jobState"`
	InputSource          string                             `json:"inputSource"`
	ExecutionSummary     *TranslationJobExecutionSummaryDTO `json:"executionSummary,omitempty"`
	ValidationPassSlices []string                           `json:"validationPassSlices"`
	ErrorKind            string                             `json:"errorKind,omitempty"`
}

// GetTranslationJobSetupSummaryRequestDTO identifies the requested created job.
type GetTranslationJobSetupSummaryRequestDTO struct {
	JobID int64 `json:"jobId"`
}

// TranslationJobSetupSummaryResponseDTO returns the frozen read-only job summary shape.
type TranslationJobSetupSummaryResponseDTO struct {
	JobID                int64                             `json:"jobId"`
	JobState             string                            `json:"jobState"`
	InputSource          string                            `json:"inputSource"`
	CanStartPhase        bool                              `json:"canStartPhase"`
	ExecutionSummary     TranslationJobExecutionSummaryDTO `json:"executionSummary"`
	ValidationPassSlices []string                          `json:"validationPassSlices"`
}

// NewTranslationJobSetupController creates a Job Setup controller.
func NewTranslationJobSetupController(usecase TranslationJobSetupUsecasePort) *TranslationJobSetupController {
	return &TranslationJobSetupController{translationJobSetupUsecase: usecase}
}

// GetTranslationJobSetupOptions returns the frozen Job Setup option contract.
func (controller *TranslationJobSetupController) GetTranslationJobSetupOptions() (TranslationJobSetupOptionsResponseDTO, error) {
	result, err := controller.translationJobSetupUsecase.GetTranslationJobSetupOptions(context.Background())
	if err != nil {
		return TranslationJobSetupOptionsResponseDTO{}, fmt.Errorf("get translation job setup options: %w", err)
	}
	return toTranslationJobSetupOptionsResponseDTO(result), nil
}

// ValidateTranslationJobSetup validates one Job Setup request.
func (controller *TranslationJobSetupController) ValidateTranslationJobSetup(
	request ValidateTranslationJobSetupRequestDTO,
) (TranslationJobSetupValidationResponseDTO, error) {
	result, err := controller.translationJobSetupUsecase.ValidateTranslationJobSetup(
		context.Background(),
		usecase.ValidateTranslationJobSetupRequest{
			InputSourceID: request.InputSourceID,
			Runtime:       toTranslationJobSetupRuntimeSelection(request.Runtime),
			CredentialRef: request.CredentialRef,
		},
	)
	if err != nil {
		return TranslationJobSetupValidationResponseDTO{}, fmt.Errorf("validate translation job setup: %w", err)
	}
	return toTranslationJobSetupValidationResponseDTO(result), nil
}

// CreateTranslationJob creates one ready translation job or returns a rejected error kind.
func (controller *TranslationJobSetupController) CreateTranslationJob(
	request CreateTranslationJobRequestDTO,
) (CreateTranslationJobResponseDTO, error) {
	var validatedAt time.Time
	if request.ValidatedAt != "" {
		parsedValidatedAt, err := time.Parse(time.RFC3339, request.ValidatedAt)
		if err != nil {
			return CreateTranslationJobResponseDTO{}, fmt.Errorf("parse create translation job validation freshness: %w", err)
		}
		validatedAt = parsedValidatedAt.UTC()
	}

	result, err := controller.translationJobSetupUsecase.CreateTranslationJob(
		context.Background(),
		usecase.CreateTranslationJobRequest{
			InputSourceID:        request.InputSourceID,
			InputSource:          request.InputSource,
			ValidationStatus:     request.ValidationStatus,
			ValidatedAt:          validatedAt,
			ValidationPassSlices: cloneStrings(request.ValidationPassSlices),
			Runtime:              toTranslationJobSetupRuntimeSelection(request.Runtime),
			CredentialRef:        request.CredentialRef,
		},
	)
	if err != nil {
		return CreateTranslationJobResponseDTO{}, fmt.Errorf("create translation job: %w", err)
	}
	return toCreateTranslationJobResponseDTO(result), nil
}

// GetTranslationJobSetupSummary returns the frozen read-only job summary.
func (controller *TranslationJobSetupController) GetTranslationJobSetupSummary(
	request GetTranslationJobSetupSummaryRequestDTO,
) (TranslationJobSetupSummaryResponseDTO, error) {
	result, err := controller.translationJobSetupUsecase.GetTranslationJobSetupSummary(
		context.Background(),
		usecase.GetTranslationJobSetupSummaryRequest{JobID: request.JobID},
	)
	if err != nil {
		return TranslationJobSetupSummaryResponseDTO{}, fmt.Errorf("get translation job setup summary: %w", err)
	}
	return toTranslationJobSetupSummaryResponseDTO(result), nil
}

func toTranslationJobSetupOptionsResponseDTO(result usecase.TranslationJobSetupOptionsResult) TranslationJobSetupOptionsResponseDTO {
	response := TranslationJobSetupOptionsResponseDTO{
		InputCandidates:    toTranslationJobSetupInputCandidateDTOs(result.InputCandidates),
		SharedDictionaries: toTranslationJobSetupDictionaryOptionDTOs(result.SharedDictionaries),
		SharedPersonas:     toTranslationJobSetupPersonaOptionDTOs(result.SharedPersonas),
		AIRuntimeOptions:   toTranslationJobSetupRuntimeOptionDTOs(result.AIRuntimeOptions),
		CredentialRefs:     toTranslationJobSetupCredentialReferenceDTOs(result.CredentialRefs),
	}
	if result.ExistingJob != nil {
		existingJob := toTranslationJobSetupExistingJobDTO(*result.ExistingJob)
		response.ExistingJob = &existingJob
	}
	return response
}

func toTranslationJobSetupInputCandidateDTOs(candidates []usecase.TranslationJobSetupInputCandidate) []TranslationJobSetupInputCandidateDTO {
	results := make([]TranslationJobSetupInputCandidateDTO, 0, len(candidates))
	for _, candidate := range candidates {
		results = append(results, TranslationJobSetupInputCandidateDTO{
			ID:           candidate.ID,
			Label:        candidate.Label,
			SourceKind:   candidate.SourceKind,
			RecordCount:  candidate.RecordCount,
			RegisteredAt: candidate.RegisteredAt.UTC().Format(time.RFC3339),
		})
	}
	return results
}

func toTranslationJobSetupExistingJobDTO(existingJob usecase.TranslationJobSetupExistingJob) TranslationJobSetupExistingJobDTO {
	return TranslationJobSetupExistingJobDTO{
		InputSourceID: existingJob.InputSourceID,
		JobID:         existingJob.JobID,
		Status:        existingJob.Status,
		InputSource:   existingJob.InputSource,
	}
}

func toTranslationJobSetupDictionaryOptionDTOs(options []usecase.TranslationJobSetupDictionaryOption) []TranslationJobSetupDictionaryOptionDTO {
	results := make([]TranslationJobSetupDictionaryOptionDTO, 0, len(options))
	for _, option := range options {
		results = append(results, TranslationJobSetupDictionaryOptionDTO{ID: option.ID, Label: option.Label})
	}
	return results
}

func toTranslationJobSetupPersonaOptionDTOs(options []usecase.TranslationJobSetupPersonaOption) []TranslationJobSetupPersonaOptionDTO {
	results := make([]TranslationJobSetupPersonaOptionDTO, 0, len(options))
	for _, option := range options {
		results = append(results, TranslationJobSetupPersonaOptionDTO{ID: option.ID, Label: option.Label})
	}
	return results
}

func toTranslationJobSetupRuntimeOptionDTOs(options []usecase.TranslationJobSetupRuntimeOption) []TranslationJobSetupRuntimeOptionDTO {
	results := make([]TranslationJobSetupRuntimeOptionDTO, 0, len(options))
	for _, option := range options {
		results = append(results, TranslationJobSetupRuntimeOptionDTO{
			Provider: option.Provider,
			Model:    option.Model,
			Mode:     option.Mode,
		})
	}
	return results
}

func toTranslationJobSetupCredentialReferenceDTOs(refs []usecase.TranslationJobSetupCredentialReference) []TranslationJobSetupCredentialReferenceDTO {
	results := make([]TranslationJobSetupCredentialReferenceDTO, 0, len(refs))
	for _, ref := range refs {
		results = append(results, TranslationJobSetupCredentialReferenceDTO{
			Provider:        ref.Provider,
			CredentialRef:   ref.CredentialRef,
			IsConfigured:    ref.IsConfigured,
			IsMissingSecret: ref.IsMissingSecret,
		})
	}
	return results
}

func toTranslationJobSetupRuntimeSelection(runtime TranslationJobSetupRuntimeSelectionDTO) usecase.TranslationJobSetupRuntimeSelection {
	return usecase.TranslationJobSetupRuntimeSelection{
		Provider:      runtime.Provider,
		Model:         runtime.Model,
		ExecutionMode: runtime.ExecutionMode,
	}
}

func toTranslationJobSetupValidationResponseDTO(result usecase.TranslationJobSetupValidationResult) TranslationJobSetupValidationResponseDTO {
	return TranslationJobSetupValidationResponseDTO{
		Status:                  result.Status,
		BlockingFailureCategory: usecase.NormalizeTranslationJobSetupPublicErrorCategory(result.BlockingFailureCategory),
		TargetSlices:            cloneStrings(result.TargetSlices),
		ValidatedAt:             result.ValidatedAt.UTC().Format(time.RFC3339),
		CanCreate:               result.CanCreate,
		PassSlices:              cloneStrings(result.PassSlices),
	}
}

func toCreateTranslationJobResponseDTO(result usecase.CreateTranslationJobResult) CreateTranslationJobResponseDTO {
	response := CreateTranslationJobResponseDTO{
		JobID:                result.JobID,
		JobState:             result.JobState,
		InputSource:          result.InputSource,
		ValidationPassSlices: cloneStrings(result.ValidationPassSlices),
		ErrorKind:            usecase.NormalizeTranslationJobSetupPublicErrorKind(result.ErrorKind),
	}
	if result.ErrorKind == "" {
		executionSummary := toTranslationJobExecutionSummaryDTO(result.ExecutionSummary)
		response.ExecutionSummary = &executionSummary
	}
	return response
}

func toTranslationJobSetupSummaryResponseDTO(result usecase.TranslationJobSetupSummaryResult) TranslationJobSetupSummaryResponseDTO {
	return TranslationJobSetupSummaryResponseDTO{
		JobID:                result.JobID,
		JobState:             result.JobState,
		InputSource:          result.InputSource,
		CanStartPhase:        result.CanStartPhase,
		ExecutionSummary:     toTranslationJobExecutionSummaryDTO(result.ExecutionSummary),
		ValidationPassSlices: cloneStrings(result.ValidationPassSlices),
	}
}

func toTranslationJobExecutionSummaryDTO(summary usecase.TranslationJobExecutionSummary) TranslationJobExecutionSummaryDTO {
	return TranslationJobExecutionSummaryDTO{
		Provider:      summary.Provider,
		Model:         summary.Model,
		ExecutionMode: summary.ExecutionMode,
	}
}

func cloneStrings(values []string) []string {
	if values == nil {
		return nil
	}
	return append([]string(nil), values...)
}
