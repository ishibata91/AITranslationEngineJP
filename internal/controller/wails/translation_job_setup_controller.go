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
	DeleteTranslationJobSetupInput(ctx context.Context, request usecase.DeleteTranslationJobSetupInputRequest) (usecase.DeleteTranslationJobSetupInputResult, error)
	GetTranslationJobSetupSummary(ctx context.Context, request usecase.GetTranslationJobSetupSummaryRequest) (usecase.TranslationJobSetupSummaryResult, error)
}

// TranslationJobSetupProviderModelListUsecasePort defines the optional provider model list seam.
type TranslationJobSetupProviderModelListUsecasePort interface {
	ListTranslationJobSetupProviderModels(
		ctx context.Context,
		request usecase.ListTranslationJobSetupProviderModelsRequest,
	) (usecase.ListTranslationJobSetupProviderModelsResult, error)
}

// TranslationJobSetupController exposes Wails-bound Job Setup entrypoints.
type TranslationJobSetupController struct {
	translationJobSetupUsecase TranslationJobSetupUsecasePort
}

// TranslationJobSetupInputCandidateDTO is one selectable translation input candidate.
type TranslationJobSetupInputCandidateDTO struct {
	ID           int64                              `json:"id"`
	Label        string                             `json:"label"`
	SourceKind   string                             `json:"sourceKind"`
	RecordCount  int                                `json:"recordCount"`
	RegisteredAt string                             `json:"registeredAt"`
	ExistingJob  *TranslationJobSetupExistingJobDTO `json:"existingJob,omitempty"`
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

// TranslationJobSetupProviderCapabilityDTO is one public provider capability.
type TranslationJobSetupProviderCapabilityDTO struct {
	Provider                string   `json:"provider"`
	CredentialRequirement   string   `json:"credentialRequirement"`
	SupportedExecutionModes []string `json:"supportedExecutionModes"`
	SupportsBatchMode       bool     `json:"supportsBatchMode"`
}

// TranslationJobSetupCredentialReferenceDTO exposes only credential reference state.
type TranslationJobSetupCredentialReferenceDTO struct {
	Provider        string  `json:"provider"`
	IsConfigured    bool    `json:"isConfigured"`
	IsMissingSecret bool    `json:"isMissingSecret"`
	SecretPlaintext *string `json:"-"`
}

// TranslationJobSetupPhaseRuntimeDraftDTO returns the current draft state for one phase.
type TranslationJobSetupPhaseRuntimeDraftDTO struct {
	PhaseID              string `json:"phaseId"`
	Provider             string `json:"provider"`
	Model                string `json:"model"`
	CredentialRef        string `json:"-"`
	CredentialStatus     string `json:"credentialStatus"`
	ExecutionMode        string `json:"executionMode"`
	BatchMode            string `json:"batchMode"`
	ModelListSourceToken string `json:"-"`
}

// TranslationJobSetupOptionsResponseDTO returns the read-only setup options.
type TranslationJobSetupOptionsResponseDTO struct {
	InputCandidates      []TranslationJobSetupInputCandidateDTO     `json:"inputCandidates"`
	ExistingJob          *TranslationJobSetupExistingJobDTO         `json:"existingJob,omitempty"`
	SharedDictionaries   []TranslationJobSetupDictionaryOptionDTO   `json:"sharedDictionaries"`
	SharedPersonas       []TranslationJobSetupPersonaOptionDTO      `json:"sharedPersonas"`
	AIRuntimeOptions     []TranslationJobSetupRuntimeOptionDTO      `json:"aiRuntimeOptions"`
	ProviderCapabilities []TranslationJobSetupProviderCapabilityDTO `json:"providerCapabilities"`
	PhaseRuntimeDrafts   []TranslationJobSetupPhaseRuntimeDraftDTO  `json:"phaseRuntimeDrafts"`
}

// TranslationJobSetupRuntimeSelectionDTO carries the runtime selection for validation and create.
type TranslationJobSetupRuntimeSelectionDTO struct {
	Provider      string `json:"provider"`
	Model         string `json:"model"`
	ExecutionMode string `json:"executionMode"`
}

// TranslationJobSetupPhaseRuntimeSelectionDTO carries the runtime selection for one phase.
type TranslationJobSetupPhaseRuntimeSelectionDTO struct {
	PhaseID          string `json:"phaseId"`
	Provider         string `json:"provider"`
	Model            string `json:"model"`
	CredentialStatus string `json:"credentialStatus"`
	ExecutionMode    string `json:"executionMode"`
	BatchMode        string `json:"batchMode"`
	FreshnessToken   string `json:"modelListFreshnessToken,omitempty"`
}

// ListTranslationJobSetupProviderModelsRequestDTO carries the provider model list request payload.
type ListTranslationJobSetupProviderModelsRequestDTO struct {
	PhaseID          string `json:"phaseId"`
	Provider         string `json:"provider"`
	CredentialStatus string `json:"credentialStatus"`
	RequestToken     string `json:"requestToken"`
}

// TranslationJobSetupProviderModelOptionDTO is one selectable provider model.
type TranslationJobSetupProviderModelOptionDTO struct {
	ModelID string `json:"modelId"`
	Label   string `json:"label"`
}

// ListTranslationJobSetupProviderModelsResponseDTO returns the provider model list state.
type ListTranslationJobSetupProviderModelsResponseDTO struct {
	PhaseID          string                                      `json:"phaseId"`
	Provider         string                                      `json:"provider"`
	CredentialStatus string                                      `json:"credentialStatus"`
	RequestToken     string                                      `json:"requestToken"`
	SourceToken      string                                      `json:"sourceToken"`
	Status           string                                      `json:"status"`
	Models           []TranslationJobSetupProviderModelOptionDTO `json:"models"`
	FailureKind      string                                      `json:"failureKind,omitempty"`
}

// ValidateTranslationJobSetupRequestDTO carries the frozen validation request payload.
type ValidateTranslationJobSetupRequestDTO struct {
	InputSourceID          int64                                         `json:"inputSourceId"`
	Runtime                TranslationJobSetupRuntimeSelectionDTO        `json:"runtime"`
	PhaseRuntimeSelections []TranslationJobSetupPhaseRuntimeSelectionDTO `json:"phaseRuntimeSelections"`
}

// TranslationJobSetupPhaseValidationResponseDTO returns the validation result for one phase.
type TranslationJobSetupPhaseValidationResponseDTO struct {
	PhaseID                 string  `json:"phaseId"`
	Status                  string  `json:"status"`
	BlockingFailureCategory *string `json:"blockingFailureCategory,omitempty"`
	CanCreate               bool    `json:"canCreate"`
	ModelListState          string  `json:"modelListState"`
	IsModelSelectionStale   bool    `json:"isModelSelectionStale"`
}

// TranslationJobSetupValidationResponseDTO returns the frozen validation result shape.
type TranslationJobSetupValidationResponseDTO struct {
	Status                  string                                          `json:"status"`
	BlockingFailureCategory *string                                         `json:"blockingFailureCategory,omitempty"`
	TargetSlices            []string                                        `json:"targetSlices"`
	ValidatedAt             string                                          `json:"validatedAt"`
	CanCreate               bool                                            `json:"canCreate"`
	PassSlices              []string                                        `json:"passSlices"`
	PhaseResults            []TranslationJobSetupPhaseValidationResponseDTO `json:"phaseResults"`
	StaleModelListPhaseIDs  []string                                        `json:"staleModelListPhaseIds"`
}

// CreateTranslationJobRequestDTO carries the frozen create request payload.
type CreateTranslationJobRequestDTO struct {
	InputSourceID          int64                                         `json:"inputSourceId"`
	InputSource            string                                        `json:"inputSource"`
	ValidationStatus       string                                        `json:"validationStatus"`
	ValidatedAt            string                                        `json:"validatedAt"`
	ValidationPassSlices   []string                                      `json:"validationPassSlices"`
	Runtime                TranslationJobSetupRuntimeSelectionDTO        `json:"runtime"`
	PhaseRuntimeSelections []TranslationJobSetupPhaseRuntimeSelectionDTO `json:"phaseRuntimeSelections"`
}

// TranslationJobExecutionSummaryDTO returns the runtime summary captured by a created job.
type TranslationJobExecutionSummaryDTO struct {
	Provider      string `json:"provider"`
	Model         string `json:"model"`
	ExecutionMode string `json:"executionMode"`
}

// TranslationJobSetupPhaseRuntimeSummaryDTO returns the runtime snapshot for one phase.
type TranslationJobSetupPhaseRuntimeSummaryDTO struct {
	PhaseID              string `json:"phaseId"`
	Provider             string `json:"provider"`
	Model                string `json:"model"`
	CredentialRef        string `json:"-"`
	CredentialStatus     string `json:"credentialStatus"`
	ExecutionMode        string `json:"executionMode"`
	BatchMode            string `json:"batchMode"`
	ModelListSourceToken string `json:"-"`
}

// CreateTranslationJobResponseDTO returns either a ready job or a rejected error kind.
type CreateTranslationJobResponseDTO struct {
	JobID                 int64                                       `json:"jobId"`
	JobState              string                                      `json:"jobState"`
	InputSource           string                                      `json:"inputSource"`
	ExecutionSummary      *TranslationJobExecutionSummaryDTO          `json:"executionSummary,omitempty"`
	ValidationPassSlices  []string                                    `json:"validationPassSlices"`
	ErrorKind             string                                      `json:"errorKind,omitempty"`
	PhaseRuntimeSummaries []TranslationJobSetupPhaseRuntimeSummaryDTO `json:"phaseRuntimeSummaries"`
}

// DeleteTranslationJobSetupInputRequestDTO identifies one Job Setup input delete target.
type DeleteTranslationJobSetupInputRequestDTO struct {
	InputSourceID int64 `json:"inputSourceId"`
}

// DeleteTranslationJobSetupInputResponseDTO returns one input delete outcome.
type DeleteTranslationJobSetupInputResponseDTO struct {
	DeletedInputSourceID *int64 `json:"deletedInputSourceId,omitempty"`
	ErrorKind            string `json:"errorKind,omitempty"`
}

// GetTranslationJobSetupSummaryRequestDTO identifies the requested created job.
type GetTranslationJobSetupSummaryRequestDTO struct {
	JobID int64 `json:"jobId"`
}

// TranslationJobSetupSummaryResponseDTO returns the frozen read-only job summary shape.
type TranslationJobSetupSummaryResponseDTO struct {
	JobID                 int64                                       `json:"jobId"`
	JobState              string                                      `json:"jobState"`
	InputSource           string                                      `json:"inputSource"`
	CanStartPhase         bool                                        `json:"canStartPhase"`
	ExecutionSummary      TranslationJobExecutionSummaryDTO           `json:"executionSummary"`
	ValidationPassSlices  []string                                    `json:"validationPassSlices"`
	PhaseRuntimeSummaries []TranslationJobSetupPhaseRuntimeSummaryDTO `json:"phaseRuntimeSummaries"`
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

// ListTranslationJobSetupProviderModels returns the frozen provider model list contract.
func (controller *TranslationJobSetupController) ListTranslationJobSetupProviderModels(
	request ListTranslationJobSetupProviderModelsRequestDTO,
) (ListTranslationJobSetupProviderModelsResponseDTO, error) {
	modelListUsecase, ok := controller.translationJobSetupUsecase.(TranslationJobSetupProviderModelListUsecasePort)
	if !ok {
		return ListTranslationJobSetupProviderModelsResponseDTO{}, fmt.Errorf(
			"list translation job setup provider models: %w",
			errTranslationJobSetupProviderModelListNotImplemented,
		)
	}

	result, err := modelListUsecase.ListTranslationJobSetupProviderModels(
		context.Background(),
		usecase.ListTranslationJobSetupProviderModelsRequest{
			PhaseID:          usecase.TranslationJobSetupPhaseID(request.PhaseID),
			Provider:         request.Provider,
			CredentialStatus: usecase.TranslationJobSetupCredentialStatus(request.CredentialStatus),
			RequestToken:     request.RequestToken,
		},
	)
	if err != nil {
		return ListTranslationJobSetupProviderModelsResponseDTO{}, fmt.Errorf(
			"list translation job setup provider models: %w",
			err,
		)
	}
	return toListTranslationJobSetupProviderModelsResponseDTO(result), nil
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
			PhaseRuntimeSelections: toTranslationJobSetupEffectivePhaseRuntimeSelections(
				request.Runtime,
				request.PhaseRuntimeSelections,
			),
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
			ValidationStatus:     usecase.TranslationJobSetupValidationStatus(request.ValidationStatus),
			ValidatedAt:          validatedAt,
			ValidationPassSlices: cloneStrings(request.ValidationPassSlices),
			Runtime:              toTranslationJobSetupRuntimeSelection(request.Runtime),
			PhaseRuntimeSelections: toTranslationJobSetupEffectivePhaseRuntimeSelections(
				request.Runtime,
				request.PhaseRuntimeSelections,
			),
		},
	)
	if err != nil {
		return CreateTranslationJobResponseDTO{}, fmt.Errorf("create translation job: %w", err)
	}
	return toCreateTranslationJobResponseDTO(result), nil
}

// DeleteTranslationJobSetupInput deletes one unreferenced input candidate.
func (controller *TranslationJobSetupController) DeleteTranslationJobSetupInput(
	request DeleteTranslationJobSetupInputRequestDTO,
) (DeleteTranslationJobSetupInputResponseDTO, error) {
	result, err := controller.translationJobSetupUsecase.DeleteTranslationJobSetupInput(
		context.Background(),
		usecase.DeleteTranslationJobSetupInputRequest{InputSourceID: request.InputSourceID},
	)
	if err != nil {
		return DeleteTranslationJobSetupInputResponseDTO{}, fmt.Errorf("delete translation job setup input: %w", err)
	}
	return toDeleteTranslationJobSetupInputResponseDTO(result), nil
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
		InputCandidates:      toTranslationJobSetupInputCandidateDTOs(result.InputCandidates),
		SharedDictionaries:   toTranslationJobSetupDictionaryOptionDTOs(result.SharedDictionaries),
		SharedPersonas:       toTranslationJobSetupPersonaOptionDTOs(result.SharedPersonas),
		AIRuntimeOptions:     toTranslationJobSetupRuntimeOptionDTOs(result.AIRuntimeOptions),
		ProviderCapabilities: toTranslationJobSetupProviderCapabilityDTOs(result.ProviderCapabilities),
		PhaseRuntimeDrafts:   toTranslationJobSetupPhaseRuntimeDraftDTOs(result.PhaseRuntimeDrafts),
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
		dto := TranslationJobSetupInputCandidateDTO{
			ID:           candidate.ID,
			Label:        candidate.Label,
			SourceKind:   candidate.SourceKind,
			RecordCount:  candidate.RecordCount,
			RegisteredAt: candidate.RegisteredAt.UTC().Format(time.RFC3339),
		}
		if candidate.ExistingJob != nil {
			existingJob := toTranslationJobSetupExistingJobDTO(*candidate.ExistingJob)
			dto.ExistingJob = &existingJob
		}
		results = append(results, dto)
	}
	return results
}

func toDeleteTranslationJobSetupInputResponseDTO(
	result usecase.DeleteTranslationJobSetupInputResult,
) DeleteTranslationJobSetupInputResponseDTO {
	response := DeleteTranslationJobSetupInputResponseDTO{
		ErrorKind: string(result.ErrorKind),
	}
	if result.DeletedInputSourceID != nil {
		deletedInputSourceID := *result.DeletedInputSourceID
		response.DeletedInputSourceID = &deletedInputSourceID
	}
	return response
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

func toTranslationJobSetupProviderCapabilityDTOs(
	capabilities []usecase.TranslationJobSetupProviderCapability,
) []TranslationJobSetupProviderCapabilityDTO {
	results := make([]TranslationJobSetupProviderCapabilityDTO, 0, len(capabilities))
	for _, capability := range capabilities {
		results = append(results, TranslationJobSetupProviderCapabilityDTO{
			Provider:                capability.Provider,
			CredentialRequirement:   string(capability.CredentialRequirement),
			SupportedExecutionModes: cloneStrings(capability.SupportedExecutionModes),
			SupportsBatchMode:       capability.SupportsBatchMode,
		})
	}
	return results
}

func toTranslationJobSetupPhaseRuntimeDraftDTOs(
	drafts []usecase.TranslationJobSetupPhaseRuntimeDraft,
) []TranslationJobSetupPhaseRuntimeDraftDTO {
	results := make([]TranslationJobSetupPhaseRuntimeDraftDTO, 0, len(drafts))
	for _, draft := range drafts {
		results = append(results, TranslationJobSetupPhaseRuntimeDraftDTO{
			PhaseID:          string(draft.PhaseID),
			Provider:         draft.Provider,
			Model:            draft.Model,
			CredentialStatus: string(draft.CredentialStatus),
			ExecutionMode:    draft.ExecutionMode,
			BatchMode:        string(draft.BatchMode),
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

func toTranslationJobSetupPhaseRuntimeSelections(
	selections []TranslationJobSetupPhaseRuntimeSelectionDTO,
) []usecase.TranslationJobSetupPhaseRuntimeSelection {
	if len(selections) == 0 {
		return []usecase.TranslationJobSetupPhaseRuntimeSelection{}
	}

	results := make([]usecase.TranslationJobSetupPhaseRuntimeSelection, 0, len(selections))
	for _, selection := range selections {
		results = append(results, usecase.TranslationJobSetupPhaseRuntimeSelection{
			PhaseID:          usecase.TranslationJobSetupPhaseID(selection.PhaseID),
			Provider:         selection.Provider,
			Model:            selection.Model,
			CredentialStatus: usecase.TranslationJobSetupCredentialStatus(selection.CredentialStatus),
			ExecutionMode:    selection.ExecutionMode,
			BatchMode:        usecase.TranslationJobSetupBatchMode(selection.BatchMode),
			FreshnessToken:   selection.FreshnessToken,
		})
	}
	return results
}

func toTranslationJobSetupEffectivePhaseRuntimeSelections(
	runtime TranslationJobSetupRuntimeSelectionDTO,
	selections []TranslationJobSetupPhaseRuntimeSelectionDTO,
) []usecase.TranslationJobSetupPhaseRuntimeSelection {
	convertedSelections := toTranslationJobSetupPhaseRuntimeSelections(selections)
	if len(convertedSelections) > 0 {
		return convertedSelections
	}

	credentialStatus := usecase.TranslationJobSetupCredentialStatusConfigured
	if runtime.Provider == "lm_studio" {
		credentialStatus = usecase.TranslationJobSetupCredentialStatusNotRequired
	}

	legacyPhaseSelections := []TranslationJobSetupPhaseRuntimeSelectionDTO{
		{
			PhaseID:          string(usecase.TranslationJobSetupPhaseIDWordTranslation),
			Provider:         runtime.Provider,
			Model:            runtime.Model,
			CredentialStatus: string(credentialStatus),
			ExecutionMode:    runtime.ExecutionMode,
			BatchMode:        string(usecase.TranslationJobSetupBatchModeUnsupported),
		},
		{
			PhaseID:          string(usecase.TranslationJobSetupPhaseIDNPCPersonaGeneration),
			Provider:         runtime.Provider,
			Model:            runtime.Model,
			CredentialStatus: string(credentialStatus),
			ExecutionMode:    runtime.ExecutionMode,
			BatchMode:        string(usecase.TranslationJobSetupBatchModeUnsupported),
		},
		{
			PhaseID:          string(usecase.TranslationJobSetupPhaseIDTextTranslation),
			Provider:         runtime.Provider,
			Model:            runtime.Model,
			CredentialStatus: string(credentialStatus),
			ExecutionMode:    runtime.ExecutionMode,
			BatchMode:        string(usecase.TranslationJobSetupBatchModeUnsupported),
		},
	}
	return toTranslationJobSetupPhaseRuntimeSelections(legacyPhaseSelections)
}

func toListTranslationJobSetupProviderModelsResponseDTO(
	result usecase.ListTranslationJobSetupProviderModelsResult,
) ListTranslationJobSetupProviderModelsResponseDTO {
	return ListTranslationJobSetupProviderModelsResponseDTO{
		PhaseID:          string(result.PhaseID),
		Provider:         result.Provider,
		CredentialStatus: string(result.CredentialStatus),
		RequestToken:     result.RequestToken,
		SourceToken:      result.SourceToken,
		Status:           string(result.Status),
		Models:           toTranslationJobSetupProviderModelOptionDTOs(result.Models),
		FailureKind:      string(usecase.NormalizeTranslationJobSetupPublicErrorKind(result.FailureKind)),
	}
}

func toTranslationJobSetupProviderModelOptionDTOs(
	models []usecase.TranslationJobSetupProviderModelOption,
) []TranslationJobSetupProviderModelOptionDTO {
	results := make([]TranslationJobSetupProviderModelOptionDTO, 0, len(models))
	for _, model := range models {
		results = append(results, TranslationJobSetupProviderModelOptionDTO{
			ModelID: model.ModelID,
			Label:   model.Label,
		})
	}
	return results
}

func toTranslationJobSetupValidationResponseDTO(result usecase.TranslationJobSetupValidationResult) TranslationJobSetupValidationResponseDTO {
	return TranslationJobSetupValidationResponseDTO{
		Status:                  string(result.Status),
		BlockingFailureCategory: usecase.NormalizeTranslationJobSetupPublicErrorCategory(result.BlockingFailureCategory),
		TargetSlices:            cloneStrings(result.TargetSlices),
		ValidatedAt:             result.ValidatedAt.UTC().Format(time.RFC3339),
		CanCreate:               result.CanCreate,
		PassSlices:              cloneStrings(result.PassSlices),
		PhaseResults:            toTranslationJobSetupPhaseValidationResponseDTOs(result.PhaseResults),
		StaleModelListPhaseIDs:  toTranslationJobSetupPhaseIDStrings(result.StaleModelListPhaseIDs),
	}
}

func toTranslationJobSetupPhaseValidationResponseDTOs(
	results []usecase.TranslationJobSetupPhaseValidationResult,
) []TranslationJobSetupPhaseValidationResponseDTO {
	converted := make([]TranslationJobSetupPhaseValidationResponseDTO, 0, len(results))
	for _, result := range results {
		converted = append(converted, TranslationJobSetupPhaseValidationResponseDTO{
			PhaseID:                 string(result.PhaseID),
			Status:                  string(result.Status),
			BlockingFailureCategory: usecase.NormalizeTranslationJobSetupPublicErrorCategory(result.BlockingFailureCategory),
			CanCreate:               result.CanCreate,
			ModelListState:          string(result.ModelListState),
			IsModelSelectionStale:   result.IsModelSelectionStale,
		})
	}
	return converted
}

func toCreateTranslationJobResponseDTO(result usecase.CreateTranslationJobResult) CreateTranslationJobResponseDTO {
	response := CreateTranslationJobResponseDTO{
		JobID:                 result.JobID,
		JobState:              result.JobState,
		InputSource:           result.InputSource,
		ValidationPassSlices:  cloneStrings(result.ValidationPassSlices),
		ErrorKind:             string(usecase.NormalizeTranslationJobSetupPublicErrorKind(result.ErrorKind)),
		PhaseRuntimeSummaries: toTranslationJobSetupPhaseRuntimeSummaryDTOs(result.PhaseRuntimeSummaries),
	}
	if result.ErrorKind == "" {
		executionSummary := toTranslationJobExecutionSummaryDTO(result.ExecutionSummary)
		response.ExecutionSummary = &executionSummary
	}
	return response
}

func toTranslationJobSetupSummaryResponseDTO(result usecase.TranslationJobSetupSummaryResult) TranslationJobSetupSummaryResponseDTO {
	return TranslationJobSetupSummaryResponseDTO{
		JobID:                 result.JobID,
		JobState:              result.JobState,
		InputSource:           result.InputSource,
		CanStartPhase:         result.CanStartPhase,
		ExecutionSummary:      toTranslationJobExecutionSummaryDTO(result.ExecutionSummary),
		ValidationPassSlices:  cloneStrings(result.ValidationPassSlices),
		PhaseRuntimeSummaries: toTranslationJobSetupPhaseRuntimeSummaryDTOs(result.PhaseRuntimeSummaries),
	}
}

func toTranslationJobExecutionSummaryDTO(summary usecase.TranslationJobExecutionSummary) TranslationJobExecutionSummaryDTO {
	return TranslationJobExecutionSummaryDTO{
		Provider:      summary.Provider,
		Model:         summary.Model,
		ExecutionMode: summary.ExecutionMode,
	}
}

func toTranslationJobSetupPhaseRuntimeSummaryDTOs(
	summaries []usecase.TranslationJobSetupPhaseRuntimeSummary,
) []TranslationJobSetupPhaseRuntimeSummaryDTO {
	results := make([]TranslationJobSetupPhaseRuntimeSummaryDTO, 0, len(summaries))
	for _, summary := range summaries {
		results = append(results, TranslationJobSetupPhaseRuntimeSummaryDTO{
			PhaseID:          string(summary.PhaseID),
			Provider:         summary.Provider,
			Model:            summary.Model,
			CredentialStatus: string(summary.CredentialStatus),
			ExecutionMode:    summary.ExecutionMode,
			BatchMode:        string(summary.BatchMode),
		})
	}
	return results
}

func toTranslationJobSetupPhaseIDStrings(phaseIDs []usecase.TranslationJobSetupPhaseID) []string {
	results := make([]string, 0, len(phaseIDs))
	for _, phaseID := range phaseIDs {
		results = append(results, string(phaseID))
	}
	return results
}

func cloneStrings(values []string) []string {
	if values == nil {
		return []string{}
	}
	return append([]string(nil), values...)
}

var errTranslationJobSetupProviderModelListNotImplemented = fmt.Errorf("translation job setup provider model list usecase is not implemented")
