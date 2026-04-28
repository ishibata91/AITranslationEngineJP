package usecase

import (
	"context"
	"fmt"

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

// TranslationJobSetupUsecase implements the Job Setup Wails seam.
type TranslationJobSetupUsecase struct {
	service translationJobSetupServicePort
}

// NewTranslationJobSetupUsecase creates a Job Setup usecase.
func NewTranslationJobSetupUsecase(service translationJobSetupServicePort) *TranslationJobSetupUsecase {
	return &TranslationJobSetupUsecase{service: service}
}

// GetTranslationJobSetupOptions returns a not-implemented error for the unfinished read-only slice.
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
		InputCandidates:    toTranslationJobSetupInputCandidates(readModel.InputCandidates),
		ExistingJob:        toTranslationJobSetupExistingJob(readModel.ExistingJob),
		SharedDictionaries: toTranslationJobSetupDictionaryOptions(readModel.SharedDictionaries),
		SharedPersonas:     toTranslationJobSetupPersonaOptions(readModel.SharedPersonas),
		AIRuntimeOptions:   toTranslationJobSetupRuntimeOptions(readModel.AIRuntimeOptions),
		CredentialRefs:     toTranslationJobSetupCredentialReferences(readModel.CredentialRefs),
	}, nil
}

// ValidateTranslationJobSetup returns a not-implemented error for the unfinished validation slice.
func (usecase *TranslationJobSetupUsecase) ValidateTranslationJobSetup(
	ctx context.Context,
	request ValidateTranslationJobSetupRequest,
) (TranslationJobSetupValidationResult, error) {
	decision, err := usecase.service.ValidateRequest(ctx, jobsetupservice.TranslationJobSetupValidationRequest{
		InputSourceID: request.InputSourceID,
		Provider:      request.Runtime.Provider,
		Model:         request.Runtime.Model,
		ExecutionMode: request.Runtime.ExecutionMode,
		CredentialRef: request.CredentialRef,
	})
	if err != nil {
		return TranslationJobSetupValidationResult{}, fmt.Errorf("validate translation job setup request: %w", err)
	}
	return toTranslationJobSetupValidationResult(decision), nil
}

// CreateTranslationJob rejects create requests whose setup validation did not pass.
func (usecase *TranslationJobSetupUsecase) CreateTranslationJob(
	ctx context.Context,
	request CreateTranslationJobRequest,
) (CreateTranslationJobResult, error) {
	decision, err := usecase.service.EvaluateCreateRequest(ctx, jobsetupservice.TranslationJobSetupCreateRequest{
		InputSourceID:    request.InputSourceID,
		ValidationStatus: request.ValidationStatus,
		ValidatedAt:      request.ValidatedAt,
		Provider:         request.Runtime.Provider,
		Model:            request.Runtime.Model,
		ExecutionMode:    request.Runtime.ExecutionMode,
		CredentialRef:    request.CredentialRef,
	})
	if err != nil {
		return CreateTranslationJobResult{}, fmt.Errorf("evaluate translation job setup create request: %w", err)
	}
	if !decision.CanCreate {
		return CreateTranslationJobResult{ErrorKind: mapTranslationJobSetupCreateErrorKind(decision.ErrorKind)}, nil
	}

	creator, ok := usecase.service.(translationJobSetupCreateExecutor)
	if !ok {
		return CreateTranslationJobResult{}, errTranslationJobSetupNotImplemented
	}
	created, err := creator.CreateTranslationJob(ctx, jobsetupservice.TranslationJobSetupCreateRequest{
		InputSourceID:    request.InputSourceID,
		ValidationStatus: request.ValidationStatus,
		ValidatedAt:      request.ValidatedAt,
		Provider:         request.Runtime.Provider,
		Model:            request.Runtime.Model,
		ExecutionMode:    request.Runtime.ExecutionMode,
		CredentialRef:    request.CredentialRef,
	}, decision.ValidationPassSlices)
	if err != nil {
		return CreateTranslationJobResult{}, fmt.Errorf("create translation job: %w", err)
	}
	if created.ErrorKind != "" {
		return CreateTranslationJobResult{ErrorKind: mapTranslationJobSetupCreateErrorKind(created.ErrorKind)}, nil
	}
	return CreateTranslationJobResult{
		JobID:       created.JobID,
		JobState:    created.JobState,
		InputSource: created.InputSource,
		ExecutionSummary: TranslationJobExecutionSummary{
			Provider:      created.ExecutionSummary.Provider,
			Model:         created.ExecutionSummary.Model,
			ExecutionMode: created.ExecutionSummary.ExecutionMode,
		},
		ValidationPassSlices: append([]string(nil), created.ValidationPassSlices...),
	}, nil
}

// GetTranslationJobSetupSummary returns a not-implemented error for the unfinished summary slice.
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
		JobID:         readModel.JobID,
		JobState:      readModel.JobState,
		InputSource:   readModel.InputSource,
		CanStartPhase: readModel.CanStartPhase,
		ExecutionSummary: TranslationJobExecutionSummary{
			Provider:      readModel.ExecutionSummary.Provider,
			Model:         readModel.ExecutionSummary.Model,
			ExecutionMode: readModel.ExecutionSummary.ExecutionMode,
		},
		ValidationPassSlices: append([]string(nil), readModel.ValidationPassSlices...),
	}, nil
}

func toTranslationJobSetupValidationResult(
	decision jobsetupservice.TranslationJobSetupValidationDecision,
) TranslationJobSetupValidationResult {
	return TranslationJobSetupValidationResult{
		Status:                  decision.Status,
		BlockingFailureCategory: cloneOptionalString(decision.BlockingFailureCategory),
		TargetSlices:            append([]string(nil), decision.TargetSlices...),
		ValidatedAt:             decision.ValidatedAt,
		CanCreate:               decision.CanCreate,
		PassSlices:              append([]string(nil), decision.PassSlices...),
	}
}

func cloneOptionalString(value *string) *string {
	if value == nil {
		return nil
	}
	cloned := *value
	return &cloned
}

func mapTranslationJobSetupCreateErrorKind(kind string) TranslationJobSetupErrorKind {
	switch kind {
	case "", TranslationJobSetupErrorKindReadyRequired:
		return kind
	case "validation_failed":
		return TranslationJobSetupErrorKindReadyRequired
	case "duplicate_input":
		return TranslationJobSetupErrorKindDuplicateJobForInput
	case TranslationJobSetupErrorKindRequiredSettingMissing:
		return TranslationJobSetupErrorKindRequiredSettingMissing
	case TranslationJobSetupErrorKindInputNotFound:
		return TranslationJobSetupErrorKindInputNotFound
	case TranslationJobSetupErrorKindCacheMissing:
		return TranslationJobSetupErrorKindCacheMissing
	case TranslationJobSetupErrorKindFoundationRefMissing:
		return TranslationJobSetupErrorKindFoundationRefMissing
	case TranslationJobSetupErrorKindCredentialMissing:
		return TranslationJobSetupErrorKindCredentialMissing
	case TranslationJobSetupErrorKindProviderModeUnsupported:
		return TranslationJobSetupErrorKindProviderModeUnsupported
	case TranslationJobSetupErrorKindProviderUnreachable:
		return TranslationJobSetupErrorKindProviderUnreachable
	case TranslationJobSetupErrorKindValidationStale:
		return TranslationJobSetupErrorKindValidationStale
	case TranslationJobSetupErrorKindPartialCreateFailed:
		return TranslationJobSetupErrorKindPartialCreateFailed
	default:
		return kind
	}
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
		result = append(result, TranslationJobSetupDictionaryOption{
			ID:    option.ID,
			Label: option.Label,
		})
	}
	return result
}

func toTranslationJobSetupPersonaOptions(
	options []jobsetupservice.TranslationJobSetupPersonaOptionReadModel,
) []TranslationJobSetupPersonaOption {
	result := make([]TranslationJobSetupPersonaOption, 0, len(options))
	for _, option := range options {
		result = append(result, TranslationJobSetupPersonaOption{
			ID:    option.ID,
			Label: option.Label,
		})
	}
	return result
}

func toTranslationJobSetupRuntimeOptions(
	options []jobsetupservice.TranslationJobSetupRuntimeOptionReadModel,
) []TranslationJobSetupRuntimeOption {
	result := make([]TranslationJobSetupRuntimeOption, 0, len(options))
	for _, option := range options {
		result = append(result, TranslationJobSetupRuntimeOption{
			Provider: option.Provider,
			Model:    option.Model,
			Mode:     option.Mode,
		})
	}
	return result
}

func toTranslationJobSetupCredentialReferences(
	refs []jobsetupservice.TranslationJobSetupCredentialReferenceReadModel,
) []TranslationJobSetupCredentialReference {
	result := make([]TranslationJobSetupCredentialReference, 0, len(refs))
	for _, ref := range refs {
		result = append(result, TranslationJobSetupCredentialReference{
			Provider:        ref.Provider,
			CredentialRef:   ref.CredentialRef,
			IsConfigured:    ref.IsConfigured,
			IsMissingSecret: ref.IsMissingSecret,
		})
	}
	return result
}
