package usecase

import (
	"context"
	"fmt"
	"log/slog"
	"strings"

	jobsetupservice "aitranslationenginejp/internal/service"
)

type translationJobSetupServicePort interface {
	ValidateRequest(ctx context.Context, request jobsetupservice.TranslationJobSetupValidationRequest) (jobsetupservice.TranslationJobSetupValidationDecision, error)
	EvaluateCreateRequest(ctx context.Context, request jobsetupservice.TranslationJobSetupCreateRequest) (jobsetupservice.TranslationJobSetupCreateDecision, error)
}

type translationJobSetupCreateExecutor interface {
	CreateTranslationJob(
		ctx context.Context,
		request jobsetupservice.TranslationJobSetupCreateRequest,
		validationPassSlices []string,
	) (jobsetupservice.TranslationJobSetupCreatedJobReadModel, error)
}

type translationJobSetupSummaryReader interface {
	ReadSummary(ctx context.Context, jobID int64) (jobsetupservice.TranslationJobSetupSummaryReadModel, error)
}

type translationJobSetupOptionsReader interface {
	ReadOptions(ctx context.Context) (jobsetupservice.TranslationJobSetupOptionsReadModel, error)
}

type translationJobSetupProviderModelListReader interface {
	ListProviderModels(
		ctx context.Context,
		request jobsetupservice.ListTranslationJobSetupProviderModelsRequest,
	) (jobsetupservice.ListTranslationJobSetupProviderModelsResult, error)
}

type translationJobSetupInputDeleteExecutor interface {
	DeleteInputSource(
		ctx context.Context,
		inputSourceID int64,
	) (jobsetupservice.TranslationJobSetupDeleteInputDecision, error)
}

// TranslationJobSetupUsecase implements the Job Setup Wails seam.
type TranslationJobSetupUsecase struct {
	service translationJobSetupServicePort
}

// NewTranslationJobSetupUsecase creates a Job Setup usecase.
func NewTranslationJobSetupUsecase(service translationJobSetupServicePort) *TranslationJobSetupUsecase {
	return &TranslationJobSetupUsecase{service: service}
}

// GetTranslationJobSetupOptions returns the read-only Job Setup options.
func (usecase *TranslationJobSetupUsecase) GetTranslationJobSetupOptions(
	ctx context.Context,
) (TranslationJobSetupOptionsResult, error) {
	readModel := jobsetupservice.TranslationJobSetupReadOptions()
	if reader, ok := usecase.service.(translationJobSetupOptionsReader); ok {
		persistedReadModel, err := reader.ReadOptions(ctx)
		if err != nil {
			return TranslationJobSetupOptionsResult{}, fmt.Errorf("read translation job setup options: %w", err)
		}
		readModel = persistedReadModel
	}
	return TranslationJobSetupOptionsResult{
		InputCandidates:      toTranslationJobSetupInputCandidates(readModel.InputCandidates),
		ExistingJob:          toTranslationJobSetupExistingJob(readModel.ExistingJob),
		SharedDictionaries:   toTranslationJobSetupDictionaryOptions(readModel.SharedDictionaries),
		SharedPersonas:       toTranslationJobSetupPersonaOptions(readModel.SharedPersonas),
		AIRuntimeOptions:     toTranslationJobSetupRuntimeOptions(readModel.AIRuntimeOptions),
		ProviderCapabilities: toTranslationJobSetupProviderCapabilities(readModel.ProviderCapabilities),
		PhaseRuntimeDrafts:   toTranslationJobSetupPhaseRuntimeDrafts(readModel.PhaseRuntimeDrafts),
	}, nil
}

// ListTranslationJobSetupProviderModels returns one provider model-list state.
func (usecase *TranslationJobSetupUsecase) ListTranslationJobSetupProviderModels(
	ctx context.Context,
	request ListTranslationJobSetupProviderModelsRequest,
) (ListTranslationJobSetupProviderModelsResult, error) {
	reader, ok := usecase.service.(translationJobSetupProviderModelListReader)
	if !ok {
		return ListTranslationJobSetupProviderModelsResult{}, errTranslationJobSetupNotImplemented
	}
	result, err := reader.ListProviderModels(ctx, jobsetupservice.ListTranslationJobSetupProviderModelsRequest{
		PhaseID:          string(request.PhaseID),
		Provider:         request.Provider,
		CredentialStatus: string(request.CredentialStatus),
		RequestToken:     request.RequestToken,
	})
	if err != nil {
		return ListTranslationJobSetupProviderModelsResult{}, fmt.Errorf("list translation job setup provider models: %w", err)
	}
	return ListTranslationJobSetupProviderModelsResult{
		PhaseID:          TranslationJobSetupPhaseID(result.PhaseID),
		Provider:         result.Provider,
		CredentialStatus: TranslationJobSetupCredentialStatus(result.CredentialStatus),
		RequestToken:     result.RequestToken,
		SourceToken:      result.SourceToken,
		Status:           TranslationJobSetupProviderModelListStatus(result.Status),
		Models:           toTranslationJobSetupProviderModels(result.Models),
		FailureKind:      NormalizeTranslationJobSetupPublicErrorKind(TranslationJobSetupErrorKind(result.FailureKind)),
	}, nil
}

// DeleteTranslationJobSetupInput deletes one unreferenced input from Job Setup.
func (usecase *TranslationJobSetupUsecase) DeleteTranslationJobSetupInput(
	ctx context.Context,
	request DeleteTranslationJobSetupInputRequest,
) (DeleteTranslationJobSetupInputResult, error) {
	executor, ok := usecase.service.(translationJobSetupInputDeleteExecutor)
	if !ok {
		return DeleteTranslationJobSetupInputResult{}, errTranslationJobSetupNotImplemented
	}
	decision, err := executor.DeleteInputSource(ctx, request.InputSourceID)
	if err != nil {
		return DeleteTranslationJobSetupInputResult{}, fmt.Errorf("delete translation job setup input: %w", err)
	}
	result := DeleteTranslationJobSetupInputResult{
		ErrorKind: NormalizeTranslationJobSetupPublicErrorKind(
			TranslationJobSetupErrorKind(decision.ErrorKind),
		),
	}
	if decision.DeletedInputSourceID != nil {
		deletedInputSourceID := *decision.DeletedInputSourceID
		result.DeletedInputSourceID = &deletedInputSourceID
	}
	return result, nil
}

// ValidateTranslationJobSetup validates the three phase runtime selections.
func (usecase *TranslationJobSetupUsecase) ValidateTranslationJobSetup(
	ctx context.Context,
	request ValidateTranslationJobSetupRequest,
) (TranslationJobSetupValidationResult, error) {
	decision, err := usecase.service.ValidateRequest(ctx, jobsetupservice.TranslationJobSetupValidationRequest{
		InputSourceID: request.InputSourceID,
		PhaseRuntimes: toServiceTranslationJobSetupPhaseRuntimes(request.PhaseRuntimeSelections),
	})
	if err != nil {
		return TranslationJobSetupValidationResult{}, fmt.Errorf("validate translation job setup request: %w", err)
	}
	return TranslationJobSetupValidationResult{
		Status:                  TranslationJobSetupValidationStatus(decision.Status),
		BlockingFailureCategory: cloneOptionalString(decision.BlockingFailureCategory),
		TargetSlices:            normalizeTranslationJobSetupStringSlice(decision.TargetSlices),
		ValidatedAt:             decision.ValidatedAt,
		CanCreate:               decision.CanCreate,
		PassSlices:              normalizeTranslationJobSetupStringSlice(decision.PassSlices),
		PhaseResults:            toTranslationJobSetupPhaseValidationResults(decision.PhaseResults),
		StaleModelListPhaseIDs:  toTranslationJobSetupPhaseIDs(decision.StaleModelListPhaseIDs),
	}, nil
}

// CreateTranslationJob creates one ready job from the validated phase runtimes.
func (usecase *TranslationJobSetupUsecase) CreateTranslationJob(
	ctx context.Context,
	request CreateTranslationJobRequest,
) (CreateTranslationJobResult, error) {
	slog.InfoContext(ctx, "translation job setup create requested",
		slog.String("event", "translation_job_setup_create_from_input"),
		slog.String("where", "backend.usecase.translation_job_setup.create"),
		slog.String("result", "started"),
		slog.String("id", fmt.Sprintf("input:%d", request.InputSourceID)),
	)
	serviceRequest := jobsetupservice.TranslationJobSetupCreateRequest{
		InputSourceID:        request.InputSourceID,
		ValidationStatus:     string(request.ValidationStatus),
		ValidatedAt:          request.ValidatedAt,
		PhaseRuntimes:        toServiceTranslationJobSetupPhaseRuntimes(request.PhaseRuntimeSelections),
		ValidationPassSlices: normalizeTranslationJobSetupStringSlice(request.ValidationPassSlices),
	}
	decision, err := usecase.service.EvaluateCreateRequest(ctx, serviceRequest)
	if err != nil {
		return CreateTranslationJobResult{}, fmt.Errorf("evaluate translation job setup create request: %w", err)
	}
	if !decision.CanCreate {
		slog.WarnContext(ctx, "translation job setup create rejected",
			slog.String("event", "translation_job_setup_create_from_input"),
			slog.String("where", "backend.usecase.translation_job_setup.create"),
			slog.String("result", "rejected"),
			slog.String("id", fmt.Sprintf("input:%d", request.InputSourceID)),
			slog.String("reason", strings.TrimSpace(decision.ErrorKind)),
		)
		return CreateTranslationJobResult{ErrorKind: mapTranslationJobSetupCreateErrorKind(decision.ErrorKind)}, nil
	}

	creator, ok := usecase.service.(translationJobSetupCreateExecutor)
	if !ok {
		return CreateTranslationJobResult{}, errTranslationJobSetupNotImplemented
	}
	created, err := creator.CreateTranslationJob(ctx, serviceRequest, decision.ValidationPassSlices)
	if err != nil {
		return CreateTranslationJobResult{}, fmt.Errorf("create translation job: %w", err)
	}
	if created.ErrorKind != "" {
		slog.WarnContext(ctx, "translation job setup create rejected",
			slog.String("event", "translation_job_setup_create_from_input"),
			slog.String("where", "backend.usecase.translation_job_setup.create"),
			slog.String("result", "rejected"),
			slog.String("id", fmt.Sprintf("input:%d", request.InputSourceID)),
			slog.String("reason", strings.TrimSpace(created.ErrorKind)),
		)
		return CreateTranslationJobResult{ErrorKind: mapTranslationJobSetupCreateErrorKind(created.ErrorKind)}, nil
	}
	slog.InfoContext(ctx, "translation job setup create accepted",
		slog.String("event", "translation_job_setup_create_from_input"),
		slog.String("where", "backend.usecase.translation_job_setup.create"),
		slog.String("result", "accepted"),
		slog.String("id", fmt.Sprintf("job:%d", created.JobID)),
		slog.Int("count", len(created.PhaseRuntimeSummaries)),
	)
	return CreateTranslationJobResult{
		JobID:                 created.JobID,
		JobState:              created.JobState,
		InputSource:           created.InputSource,
		ExecutionSummary:      TranslationJobExecutionSummary{Provider: created.ExecutionSummary.Provider, Model: created.ExecutionSummary.Model, ExecutionMode: created.ExecutionSummary.ExecutionMode},
		ValidationPassSlices:  normalizeTranslationJobSetupStringSlice(created.ValidationPassSlices),
		PhaseRuntimeSummaries: toTranslationJobSetupPhaseRuntimeSummaries(created.PhaseRuntimeSummaries),
	}, nil
}

// GetTranslationJobSetupSummary returns the read-only created job summary.
func (usecase *TranslationJobSetupUsecase) GetTranslationJobSetupSummary(
	ctx context.Context,
	request GetTranslationJobSetupSummaryRequest,
) (TranslationJobSetupSummaryResult, error) {
	readModel := jobsetupservice.TranslationJobSetupReadSummary(request.JobID)
	if reader, ok := usecase.service.(translationJobSetupSummaryReader); ok {
		persistedReadModel, err := reader.ReadSummary(ctx, request.JobID)
		if err != nil {
			return TranslationJobSetupSummaryResult{}, fmt.Errorf("read translation job setup summary: %w", err)
		}
		readModel = persistedReadModel
	}
	return TranslationJobSetupSummaryResult{
		JobID:                 readModel.JobID,
		JobState:              readModel.JobState,
		InputSource:           readModel.InputSource,
		CanStartPhase:         readModel.CanStartPhase,
		ExecutionSummary:      TranslationJobExecutionSummary{Provider: readModel.ExecutionSummary.Provider, Model: readModel.ExecutionSummary.Model, ExecutionMode: readModel.ExecutionSummary.ExecutionMode},
		ValidationPassSlices:  normalizeTranslationJobSetupStringSlice(readModel.ValidationPassSlices),
		PhaseRuntimeSummaries: toTranslationJobSetupPhaseRuntimeSummaries(readModel.PhaseRuntimeSummaries),
	}, nil
}

func toServiceTranslationJobSetupPhaseRuntimes(
	phaseRuntimes []TranslationJobSetupPhaseRuntimeSelection,
) []jobsetupservice.TranslationJobSetupPhaseRuntimeDraftReadModel {
	result := make([]jobsetupservice.TranslationJobSetupPhaseRuntimeDraftReadModel, 0, len(phaseRuntimes))
	for _, runtime := range phaseRuntimes {
		result = append(result, jobsetupservice.TranslationJobSetupPhaseRuntimeDraftReadModel{
			PhaseID:              string(runtime.PhaseID),
			Provider:             runtime.Provider,
			Model:                runtime.Model,
			CredentialStatus:     string(runtime.CredentialStatus),
			ExecutionMode:        runtime.ExecutionMode,
			BatchMode:            string(runtime.BatchMode),
			ModelListSourceToken: strings.TrimSpace(runtime.FreshnessToken),
		})
	}
	return result
}

func toTranslationJobSetupPhaseValidationResults(
	results []jobsetupservice.TranslationJobSetupPhaseValidationReadModel,
) []TranslationJobSetupPhaseValidationResult {
	mapped := make([]TranslationJobSetupPhaseValidationResult, 0, len(results))
	for _, result := range results {
		mapped = append(mapped, TranslationJobSetupPhaseValidationResult{
			PhaseID:                 TranslationJobSetupPhaseID(result.PhaseID),
			Status:                  TranslationJobSetupValidationStatus(result.Status),
			BlockingFailureCategory: cloneOptionalString(result.BlockingFailureCategory),
			CanCreate:               result.CanCreate,
			ModelListState:          TranslationJobSetupProviderModelListStatus(result.ModelListState),
			IsModelSelectionStale:   result.IsModelSelectionStale,
		})
	}
	return mapped
}

func toTranslationJobSetupProviderModels(
	models []jobsetupservice.TranslationJobSetupProviderModelOptionReadModel,
) []TranslationJobSetupProviderModelOption {
	result := make([]TranslationJobSetupProviderModelOption, 0, len(models))
	for _, model := range models {
		result = append(result, TranslationJobSetupProviderModelOption{ModelID: model.ModelID, Label: model.Label})
	}
	return result
}

func toTranslationJobSetupInputCandidates(
	inputCandidates []jobsetupservice.TranslationJobSetupInputCandidateReadModel,
) []TranslationJobSetupInputCandidate {
	result := make([]TranslationJobSetupInputCandidate, 0, len(inputCandidates))
	for _, candidate := range inputCandidates {
		result = append(result, TranslationJobSetupInputCandidate{
			ID:           candidate.ID,
			Label:        candidate.Label,
			SourceKind:   candidate.SourceKind,
			RecordCount:  candidate.RecordCount,
			RegisteredAt: candidate.RegisteredAt,
			ExistingJob:  toTranslationJobSetupExistingJob(candidate.ExistingJob),
		})
	}
	return result
}

func toTranslationJobSetupExistingJob(
	existingJob *jobsetupservice.TranslationJobSetupExistingJobReadModel,
) *TranslationJobSetupExistingJob {
	if existingJob == nil {
		return nil
	}
	return &TranslationJobSetupExistingJob{
		InputSourceID: existingJob.InputSourceID,
		JobID:         existingJob.JobID,
		Status:        existingJob.Status,
		InputSource:   existingJob.InputSource,
	}
}

func toTranslationJobSetupDictionaryOptions(
	options []jobsetupservice.TranslationJobSetupDictionaryOptionReadModel,
) []TranslationJobSetupDictionaryOption {
	result := make([]TranslationJobSetupDictionaryOption, 0, len(options))
	for _, option := range options {
		result = append(result, TranslationJobSetupDictionaryOption{ID: option.ID, Label: option.Label})
	}
	return result
}

func toTranslationJobSetupPersonaOptions(
	options []jobsetupservice.TranslationJobSetupPersonaOptionReadModel,
) []TranslationJobSetupPersonaOption {
	result := make([]TranslationJobSetupPersonaOption, 0, len(options))
	for _, option := range options {
		result = append(result, TranslationJobSetupPersonaOption{ID: option.ID, Label: option.Label})
	}
	return result
}

func toTranslationJobSetupRuntimeOptions(
	options []jobsetupservice.TranslationJobSetupRuntimeOptionReadModel,
) []TranslationJobSetupRuntimeOption {
	result := make([]TranslationJobSetupRuntimeOption, 0, len(options))
	for _, option := range options {
		result = append(result, TranslationJobSetupRuntimeOption{Provider: option.Provider, Model: option.Model, Mode: option.Mode})
	}
	return result
}

func toTranslationJobSetupProviderCapabilities(
	capabilities []jobsetupservice.TranslationJobSetupProviderCapabilityReadModel,
) []TranslationJobSetupProviderCapability {
	result := make([]TranslationJobSetupProviderCapability, 0, len(capabilities))
	for _, capability := range capabilities {
		result = append(result, TranslationJobSetupProviderCapability{
			Provider:                capability.Provider,
			CredentialRequirement:   TranslationJobSetupCredentialRequirement(capability.CredentialRequirement),
			SupportedExecutionModes: normalizeTranslationJobSetupStringSlice(capability.SupportedExecutionModes),
			SupportsBatchMode:       capability.SupportsBatchMode,
		})
	}
	return result
}

func toTranslationJobSetupPhaseRuntimeDrafts(
	drafts []jobsetupservice.TranslationJobSetupPhaseRuntimeDraftReadModel,
) []TranslationJobSetupPhaseRuntimeDraft {
	result := make([]TranslationJobSetupPhaseRuntimeDraft, 0, len(drafts))
	for _, draft := range drafts {
		result = append(result, TranslationJobSetupPhaseRuntimeDraft{
			PhaseID:          TranslationJobSetupPhaseID(draft.PhaseID),
			Provider:         draft.Provider,
			Model:            draft.Model,
			CredentialStatus: TranslationJobSetupCredentialStatus(draft.CredentialStatus),
			ExecutionMode:    draft.ExecutionMode,
			BatchMode:        TranslationJobSetupBatchMode(draft.BatchMode),
		})
	}
	return result
}

func toTranslationJobSetupPhaseRuntimeSummaries(
	summaries []jobsetupservice.TranslationJobSetupPhaseRuntimeSummaryReadModel,
) []TranslationJobSetupPhaseRuntimeSummary {
	result := make([]TranslationJobSetupPhaseRuntimeSummary, 0, len(summaries))
	for _, summary := range summaries {
		result = append(result, TranslationJobSetupPhaseRuntimeSummary{
			PhaseID:          TranslationJobSetupPhaseID(summary.PhaseID),
			Provider:         summary.Provider,
			Model:            summary.Model,
			CredentialStatus: TranslationJobSetupCredentialStatus(summary.CredentialStatus),
			ExecutionMode:    summary.ExecutionMode,
			BatchMode:        TranslationJobSetupBatchMode(summary.BatchMode),
		})
	}
	return result
}

func toTranslationJobSetupPhaseIDs(phaseIDs []string) []TranslationJobSetupPhaseID {
	result := make([]TranslationJobSetupPhaseID, 0, len(phaseIDs))
	for _, phaseID := range phaseIDs {
		result = append(result, TranslationJobSetupPhaseID(phaseID))
	}
	return result
}

func normalizeTranslationJobSetupStringSlice(values []string) []string {
	if values == nil {
		return []string{}
	}
	return append([]string(nil), values...)
}

func cloneOptionalString(value *string) *string {
	if value == nil {
		return nil
	}
	cloned := *value
	return &cloned
}

func mapTranslationJobSetupCreateErrorKind(kind string) TranslationJobSetupErrorKind {
	return NormalizeTranslationJobSetupPublicErrorKind(TranslationJobSetupErrorKind(kind))
}

func toTranslationJobSetupValidationResult(
	decision jobsetupservice.TranslationJobSetupValidationDecision,
) TranslationJobSetupValidationResult {
	return TranslationJobSetupValidationResult{
		Status:                  TranslationJobSetupValidationStatus(decision.Status),
		BlockingFailureCategory: cloneOptionalString(decision.BlockingFailureCategory),
		TargetSlices:            normalizeTranslationJobSetupStringSlice(decision.TargetSlices),
		ValidatedAt:             decision.ValidatedAt,
		CanCreate:               decision.CanCreate,
		PassSlices:              normalizeTranslationJobSetupStringSlice(decision.PassSlices),
		PhaseResults:            toTranslationJobSetupPhaseValidationResults(decision.PhaseResults),
		StaleModelListPhaseIDs:  toTranslationJobSetupPhaseIDs(decision.StaleModelListPhaseIDs),
	}
}
