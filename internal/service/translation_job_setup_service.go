package service

import (
	"context"
	"errors"
	"fmt"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"time"

	"aitranslationenginejp/internal/repository"
)

const (
	translationJobSetupValidationStatusPass        = "pass"
	translationJobSetupValidationStatusFail        = "fail"
	translationJobSetupErrorKindValidationStale    = "validation_stale"
	translationJobSetupErrorKindDuplicateInput     = "duplicate_input"
	translationJobSetupErrorKindInputNotFound      = "input_not_found"
	translationJobSetupErrorKindInputDeleteBlocked = "input_delete_blocked"
	translationJobSetupErrorKindReadyRequired      = "ready_required"
	translationJobSetupErrorKindPartialCreateFail  = "partial_create_failed"

	translationJobSetupValidationFreshnessCutoffHourUTC = 9

	translationJobSetupJobStateReady       = "ready"
	translationJobSetupPhaseStatePending   = "pending"
	translationJobSetupInputSource         = "translation_input"
	translationJobSetupInstructionKindWord = "default"

	translationJobSetupProviderOpenAI = "openai"
	translationJobSetupProviderGemini = "gemini"
	translationJobSetupProviderLM     = "lm_studio"
	translationJobSetupProviderXAI    = "xai"

	translationJobSetupCredentialRefOpenAIPrimary = "openai-primary"
	translationJobSetupCredentialRefGeminiPrimary = "gemini-primary"
	translationJobSetupModelGPT54Mini             = "gpt-5.4-mini"
	translationJobSetupSecretTimeout              = 250 * time.Millisecond
)

var translationJobSetupAllSlices = []string{"input", "runtime", "credentials"}
var translationJobSetupDefaultContext = context.Background()

var translationJobSetupPhaseOrder = []string{
	"word_translation",
	"npc_persona_generation",
	"text_translation",
}

var translationJobSetupUserFacingProviderIDs = []string{
	translationJobSetupProviderGemini,
	translationJobSetupProviderLM,
	translationJobSetupProviderXAI,
}

var translationJobSetupProviderCatalog = map[string]translationJobSetupProviderSpec{
	translationJobSetupProviderOpenAI: {
		ID:                   translationJobSetupProviderOpenAI,
		DefaultModel:         translationJobSetupModelGPT54Mini,
		CredentialRequired:   true,
		SupportedModes:       []string{"sync"},
		SupportsBatchMode:    false,
		DefaultCredentialRef: translationJobSetupCredentialRefOpenAIPrimary,
	},
	translationJobSetupProviderGemini: {
		ID:                   translationJobSetupProviderGemini,
		DefaultModel:         "gemini-2.5-pro",
		CredentialRequired:   true,
		SupportedModes:       []string{"sync", "batch"},
		SupportsBatchMode:    true,
		DefaultCredentialRef: translationJobSetupCredentialRefGeminiPrimary,
	},
	translationJobSetupProviderLM: {
		ID:                   translationJobSetupProviderLM,
		DefaultModel:         "lmstudio-community",
		CredentialRequired:   false,
		SupportedModes:       []string{"sync"},
		SupportsBatchMode:    false,
		DefaultCredentialRef: "lmstudio-local",
	},
	translationJobSetupProviderXAI: {
		ID:                   translationJobSetupProviderXAI,
		DefaultModel:         "grok-4",
		CredentialRequired:   true,
		SupportedModes:       []string{"sync", "batch"},
		SupportsBatchMode:    true,
		DefaultCredentialRef: "xai-primary",
	},
}

type translationJobSetupProviderSpec struct {
	ID                   string
	DefaultModel         string
	CredentialRequired   bool
	SupportedModes       []string
	SupportsBatchMode    bool
	DefaultCredentialRef string
}

// TranslationJobSetupValidationRequest carries the validation inputs needed by the service layer.
type TranslationJobSetupValidationRequest struct {
	InputSourceID int64
	PhaseRuntimes []TranslationJobSetupPhaseRuntimeDraftReadModel
}

// TranslationJobSetupValidationDecision returns the backend validation outcome.
type TranslationJobSetupValidationDecision struct {
	Status                  string
	BlockingFailureCategory *string
	TargetSlices            []string
	ValidatedAt             time.Time
	CanCreate               bool
	PassSlices              []string
	PhaseResults            []TranslationJobSetupPhaseValidationReadModel
	StaleModelListPhaseIDs  []string
}

// TranslationJobSetupPhaseValidationReadModel returns one phase validation outcome.
type TranslationJobSetupPhaseValidationReadModel struct {
	PhaseID                 string
	Status                  string
	BlockingFailureCategory *string
	CanCreate               bool
	ModelListState          string
	ModelListSourceToken    string
	IsModelSelectionStale   bool
}

// TranslationJobSetupCreateRequest carries the create gating inputs needed by the service layer.
type TranslationJobSetupCreateRequest struct {
	InputSourceID        int64
	ValidationStatus     string
	ValidatedAt          time.Time
	PhaseRuntimes        []TranslationJobSetupPhaseRuntimeDraftReadModel
	ValidationPassSlices []string
}

// TranslationJobSetupCreateDecision returns whether create may proceed.
type TranslationJobSetupCreateDecision struct {
	CanCreate            bool
	ErrorKind            string
	ValidationPassSlices []string
}

// TranslationJobSetupDeleteInputDecision returns one input delete outcome.
type TranslationJobSetupDeleteInputDecision struct {
	DeletedInputSourceID *int64
	ErrorKind            string
}

// TranslationJobSetupCreatedJobReadModel stores the created ready-job response.
type TranslationJobSetupCreatedJobReadModel struct {
	JobID                 int64
	JobState              string
	InputSource           string
	ExecutionSummary      TranslationJobSetupExecutionSummaryReadModel
	ValidationPassSlices  []string
	PhaseRuntimeSummaries []TranslationJobSetupPhaseRuntimeSummaryReadModel
	ErrorKind             string
}

// TranslationJobSetupOptionsReadModel stores the read-only page data needed by Job Setup.
type TranslationJobSetupOptionsReadModel struct {
	InputCandidates      []TranslationJobSetupInputCandidateReadModel
	ExistingJob          *TranslationJobSetupExistingJobReadModel
	SharedDictionaries   []TranslationJobSetupDictionaryOptionReadModel
	SharedPersonas       []TranslationJobSetupPersonaOptionReadModel
	AIRuntimeOptions     []TranslationJobSetupRuntimeOptionReadModel
	CredentialRefs       []TranslationJobSetupCredentialReferenceReadModel
	ProviderCapabilities []TranslationJobSetupProviderCapabilityReadModel
	PhaseRuntimeDrafts   []TranslationJobSetupPhaseRuntimeDraftReadModel
}

// TranslationJobSetupInputCandidateReadModel is one selectable translation input source.
type TranslationJobSetupInputCandidateReadModel struct {
	ID           int64
	Label        string
	SourceKind   string
	RecordCount  int
	RegisteredAt time.Time
	ExistingJob  *TranslationJobSetupExistingJobReadModel
}

// TranslationJobSetupExistingJobReadModel summarizes one already prepared job.
type TranslationJobSetupExistingJobReadModel struct {
	InputSourceID int64
	JobID         int64
	Status        string
	InputSource   string
}

// TranslationJobSetupDictionaryOptionReadModel is one shared dictionary choice.
type TranslationJobSetupDictionaryOptionReadModel struct {
	ID    string
	Label string
}

// TranslationJobSetupPersonaOptionReadModel is one shared persona choice.
type TranslationJobSetupPersonaOptionReadModel struct {
	ID    string
	Label string
}

// TranslationJobSetupRuntimeOptionReadModel is one selectable AI runtime option.
type TranslationJobSetupRuntimeOptionReadModel struct {
	Provider string
	Model    string
	Mode     string
}

// TranslationJobSetupCredentialReferenceReadModel exposes only credential reference state.
type TranslationJobSetupCredentialReferenceReadModel struct {
	Provider        string
	CredentialRef   string
	IsConfigured    bool
	IsMissingSecret bool
}

// TranslationJobSetupProviderCapabilityReadModel describes one provider capability.
type TranslationJobSetupProviderCapabilityReadModel struct {
	Provider                string
	CredentialRequirement   string
	SupportedExecutionModes []string
	SupportsBatchMode       bool
}

// TranslationJobSetupPhaseRuntimeDraftReadModel stores one phase runtime draft or snapshot.
type TranslationJobSetupPhaseRuntimeDraftReadModel struct {
	PhaseID              string
	Provider             string
	Model                string
	CredentialRef        string
	CredentialStatus     string
	ExecutionMode        string
	BatchMode            string
	ModelListSourceToken string
}

// TranslationJobSetupPhaseRuntimeSummaryReadModel stores one persisted phase runtime snapshot.
type TranslationJobSetupPhaseRuntimeSummaryReadModel = TranslationJobSetupPhaseRuntimeDraftReadModel

// ListTranslationJobSetupProviderModelsRequest carries the provider model-list request input.
type ListTranslationJobSetupProviderModelsRequest struct {
	PhaseID          string
	Provider         string
	CredentialRef    string
	CredentialStatus string
	RequestToken     string
}

// TranslationJobSetupProviderModelOptionReadModel stores one model option.
type TranslationJobSetupProviderModelOptionReadModel struct {
	ModelID string
	Label   string
}

// ListTranslationJobSetupProviderModelsResult returns the provider model-list state.
type ListTranslationJobSetupProviderModelsResult struct {
	PhaseID          string
	Provider         string
	CredentialStatus string
	RequestToken     string
	SourceToken      string
	Status           string
	Models           []TranslationJobSetupProviderModelOptionReadModel
	FailureKind      string
}

// TranslationJobSetupSummaryReadModel stores the read-only job display.
type TranslationJobSetupSummaryReadModel struct {
	JobID                 int64
	JobState              string
	InputSource           string
	CanStartPhase         bool
	ExecutionSummary      TranslationJobSetupExecutionSummaryReadModel
	ValidationPassSlices  []string
	PhaseRuntimeSummaries []TranslationJobSetupPhaseRuntimeSummaryReadModel
}

// TranslationJobSetupExecutionSummaryReadModel stores runtime fields captured by one job.
type TranslationJobSetupExecutionSummaryReadModel struct {
	Provider      string
	Model         string
	ExecutionMode string
}

// TranslationJobSetupService evaluates backend job-setup rules before persistence.
type TranslationJobSetupService struct {
	now                     func() time.Time
	jobLifecycleRepository  translationJobSetupJobLifecycleRepository
	translationSourceRepo   translationJobSetupTranslationSourceRepository
	masterDictionaryRepo    translationJobSetupMasterDictionaryRepository
	masterPersonaRepo       translationJobSetupMasterPersonaRepository
	secretStore             translationJobSetupSecretStore
	secretLoadTimeout       time.Duration
	providerSettings        ProviderSettingsConsumer
	providerModelListLoader TranslationJobSetupProviderModelListLoader
	transactor              repository.Transactor
	optionsReadModel        TranslationJobSetupOptionsReadModel
}

type translationJobSetupJobLifecycleRepository interface {
	CreateTranslationJob(ctx context.Context, draft repository.TranslationJobDraft) (repository.TranslationJob, error)
	GetTranslationJobByID(ctx context.Context, id int64) (repository.TranslationJob, error)
	CreateJobPhaseRun(ctx context.Context, draft repository.JobPhaseRunDraft) (repository.JobPhaseRun, error)
	ListJobPhaseRunsByJobID(ctx context.Context, jobID int64) ([]repository.JobPhaseRun, error)
}

type translationJobSetupPhaseRuntimeSnapshotStore interface {
	SaveTranslationJobPhaseRuntimeSnapshot(ctx context.Context, draft repository.TranslationJobPhaseRuntimeSnapshotDraft) (repository.TranslationJobPhaseRuntimeSnapshot, error)
	ListTranslationJobPhaseRuntimeSnapshots(ctx context.Context, translationJobID int64) ([]repository.TranslationJobPhaseRuntimeSnapshot, error)
}

type translationJobSetupTranslationSourceRepository interface {
	GetXEditExtractedDataByID(ctx context.Context, id int64) (repository.XEditExtractedData, error)
}

type translationJobSetupTranslationSourceLister interface {
	ListXEditExtractedData(ctx context.Context) ([]repository.XEditExtractedData, error)
}

type translationJobSetupExistingJobLoader interface {
	GetExistingTranslationJob(ctx context.Context, xEditID int64) (repository.TranslationJob, error)
}

type translationJobSetupInputDeleteRepository interface {
	DeleteXEditExtractedDataByID(ctx context.Context, id int64) error
}

type translationJobSetupMasterDictionaryRepository interface {
	List(ctx context.Context, query MasterDictionaryQuery) (MasterDictionaryListResult, error)
}

type translationJobSetupMasterPersonaRepository interface {
	List(ctx context.Context, query repository.MasterPersonaListQuery) (repository.MasterPersonaListResult, error)
}

type translationJobSetupSecretStore interface {
	Load(ctx context.Context, key string) (string, error)
}

// TranslationJobSetupProviderModelListLoader provides provider model lists without exposing transport details to the service core.
type TranslationJobSetupProviderModelListLoader interface {
	ListProviderModels(
		ctx context.Context,
		providerID string,
		apiKey string,
	) ([]TranslationJobSetupProviderModelOptionReadModel, error)
}

// TranslationJobSetupProviderModelListLoaderFunc adapts a plain function into the loader interface.
type TranslationJobSetupProviderModelListLoaderFunc func(
	ctx context.Context,
	providerID string,
	apiKey string,
) ([]TranslationJobSetupProviderModelOptionReadModel, error)

// ListProviderModels calls the wrapped function.
func (fn TranslationJobSetupProviderModelListLoaderFunc) ListProviderModels(
	ctx context.Context,
	providerID string,
	apiKey string,
) ([]TranslationJobSetupProviderModelOptionReadModel, error) {
	return fn(ctx, providerID, apiKey)
}

// TranslationJobSetupServiceOption configures optional runtime dependencies for the service.
type TranslationJobSetupServiceOption func(service *TranslationJobSetupService)

// WithTranslationJobSetupProviderSettings injects the provider settings consumer seam.
func WithTranslationJobSetupProviderSettings(
	consumer ProviderSettingsConsumer,
) TranslationJobSetupServiceOption {
	return func(service *TranslationJobSetupService) {
		if service != nil {
			service.providerSettings = consumer
		}
	}
}

// WithTranslationJobSetupProviderModelListLoader injects an adapter-owned provider model list loader.
func WithTranslationJobSetupProviderModelListLoader(
	loader TranslationJobSetupProviderModelListLoader,
) TranslationJobSetupServiceOption {
	return func(service *TranslationJobSetupService) {
		if service != nil {
			service.providerModelListLoader = loader
		}
	}
}

// NewTranslationJobSetupService creates a Job Setup service.
func NewTranslationJobSetupService() *TranslationJobSetupService {
	return &TranslationJobSetupService{now: func() time.Time { return time.Now().UTC() }}
}

// NewPersistentTranslationJobSetupService creates a Job Setup service backed by repositories.
func NewPersistentTranslationJobSetupService(
	jobLifecycleRepository translationJobSetupJobLifecycleRepository,
	translationSourceRepository translationJobSetupTranslationSourceRepository,
	masterDictionaryRepository translationJobSetupMasterDictionaryRepository,
	masterPersonaRepository translationJobSetupMasterPersonaRepository,
	_ translationJobSetupMasterPersonaAISettingsRepository,
	secretStore translationJobSetupSecretStore,
	transactor repository.Transactor,
	options ...TranslationJobSetupServiceOption,
) *TranslationJobSetupService {
	service := &TranslationJobSetupService{
		now:                    func() time.Time { return time.Now().UTC() },
		jobLifecycleRepository: jobLifecycleRepository,
		translationSourceRepo:  translationSourceRepository,
		masterDictionaryRepo:   masterDictionaryRepository,
		masterPersonaRepo:      masterPersonaRepository,
		secretStore:            secretStore,
		secretLoadTimeout:      translationJobSetupSecretTimeout,
		transactor:             transactor,
	}
	for _, option := range options {
		if option != nil {
			option(service)
		}
	}
	return service
}

type translationJobSetupMasterPersonaAISettingsRepository interface {
	LoadAISettings(ctx context.Context) (repository.MasterPersonaAISettingsRecord, error)
}

// ReadOptions returns the current Job Setup read model from server-owned state.
func (service *TranslationJobSetupService) ReadOptions(ctx context.Context) (TranslationJobSetupOptionsReadModel, error) {
	readModel, err := service.currentOptionsReadModel(requestContextOrBackground(ctx))
	if err != nil {
		return TranslationJobSetupOptionsReadModel{}, err
	}
	return cloneTranslationJobSetupOptionsReadModel(readModel), nil
}

// DeleteInputSource deletes one unreferenced input by removing the parent row.
func (service *TranslationJobSetupService) DeleteInputSource(
	ctx context.Context,
	inputSourceID int64,
) (TranslationJobSetupDeleteInputDecision, error) {
	if service.transactor == nil {
		return TranslationJobSetupDeleteInputDecision{}, fmt.Errorf("delete translation job setup input: transactor is not configured")
	}
	deleteRepository, ok := service.translationSourceRepo.(translationJobSetupInputDeleteRepository)
	if !ok {
		return TranslationJobSetupDeleteInputDecision{}, fmt.Errorf("delete translation job setup input: repository does not support delete")
	}

	var decision TranslationJobSetupDeleteInputDecision
	err := service.transactor.WithTransaction(requestContextOrBackground(ctx), func(txCtx context.Context) error {
		decisionResult, err := service.deleteInputSourceInTransaction(txCtx, deleteRepository, inputSourceID)
		if err != nil {
			return err
		}
		decision = decisionResult
		return nil
	})
	if err != nil {
		return TranslationJobSetupDeleteInputDecision{}, fmt.Errorf("delete translation job setup input transaction: %w", err)
	}
	return decision, nil
}

func (service *TranslationJobSetupService) deleteInputSourceInTransaction(
	ctx context.Context,
	deleteRepository translationJobSetupInputDeleteRepository,
	inputSourceID int64,
) (TranslationJobSetupDeleteInputDecision, error) {
	if err := service.ensureInputSourceExists(ctx, inputSourceID); err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return TranslationJobSetupDeleteInputDecision{
				ErrorKind: translationJobSetupErrorKindInputNotFound,
			}, nil
		}
		return TranslationJobSetupDeleteInputDecision{}, err
	}

	blocked, err := service.isInputSourceDeleteBlocked(ctx, inputSourceID)
	if err != nil {
		return TranslationJobSetupDeleteInputDecision{}, err
	}
	if blocked {
		return TranslationJobSetupDeleteInputDecision{
			ErrorKind: translationJobSetupErrorKindInputDeleteBlocked,
		}, nil
	}

	if err := deleteRepository.DeleteXEditExtractedDataByID(ctx, inputSourceID); err != nil {
		return TranslationJobSetupDeleteInputDecision{}, fmt.Errorf("delete translation job setup input source: %w", err)
	}
	deletedInputSourceID := inputSourceID
	return TranslationJobSetupDeleteInputDecision{
		DeletedInputSourceID: &deletedInputSourceID,
	}, nil
}

func (service *TranslationJobSetupService) ensureInputSourceExists(ctx context.Context, inputSourceID int64) error {
	if _, err := service.translationSourceRepo.GetXEditExtractedDataByID(ctx, inputSourceID); err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return repository.ErrNotFound
		}
		return fmt.Errorf("load translation job setup input delete target: %w", err)
	}
	return nil
}

func (service *TranslationJobSetupService) isInputSourceDeleteBlocked(ctx context.Context, inputSourceID int64) (bool, error) {
	loader, ok := service.translationSourceRepo.(translationJobSetupExistingJobLoader)
	if !ok {
		return false, nil
	}

	existingJob, err := loader.GetExistingTranslationJob(ctx, inputSourceID)
	if err == nil {
		return existingJob.ID != 0, nil
	}
	if errors.Is(err, repository.ErrNotFound) {
		return false, nil
	}
	return false, fmt.Errorf("load translation job setup input delete guard: %w", err)
}

// ListProviderModels loads one phase-specific provider model list after credential gating.
func (service *TranslationJobSetupService) ListProviderModels(
	ctx context.Context,
	request ListTranslationJobSetupProviderModelsRequest,
) (ListTranslationJobSetupProviderModelsResult, error) {
	spec, ok := translationJobSetupProviderCatalog[normalizeTranslationJobSetupField(request.Provider)]
	if !ok {
		return ListTranslationJobSetupProviderModelsResult{
			PhaseID:          strings.TrimSpace(request.PhaseID),
			Provider:         strings.TrimSpace(request.Provider),
			CredentialStatus: strings.TrimSpace(request.CredentialStatus),
			RequestToken:     strings.TrimSpace(request.RequestToken),
			Status:           "failed",
			FailureKind:      "provider_mode_unsupported",
		}, nil
	}

	result := ListTranslationJobSetupProviderModelsResult{
		PhaseID:          strings.TrimSpace(request.PhaseID),
		Provider:         spec.ID,
		CredentialStatus: translationJobSetupCredentialStatus(spec, request.CredentialStatus, request.CredentialRef),
		RequestToken:     strings.TrimSpace(request.RequestToken),
		Status:           "not_updated",
	}
	if service.providerSettings != nil && translationJobSetupProviderUsesProviderSettings(spec.ID) {
		return service.listProviderModelsViaProviderSettings(requestContextOrBackground(ctx), spec, result)
	}
	return service.listProviderModelsDirect(ctx, spec, request, result)
}

func (service *TranslationJobSetupService) listProviderModelsViaProviderSettings(
	ctx context.Context,
	spec translationJobSetupProviderSpec,
	result ListTranslationJobSetupProviderModelsResult,
) (ListTranslationJobSetupProviderModelsResult, error) {
	summary, ok, err := providerSettingsSummaryForProvider(ctx, service.providerSettings, spec.ID)
	if err != nil {
		return ListTranslationJobSetupProviderModelsResult{}, err
	}
	if !ok {
		result.CredentialStatus = "missing"
		result.Status = "failed"
		result.FailureKind = "required_setting_missing"
		return result, nil
	}
	credentialRef := strings.TrimSpace(pointerStringValue(summary.CredentialReferenceID))
	settingsRequestToken := strings.TrimSpace(pointerStringValue(summary.RequestToken))
	result.CredentialStatus = strings.TrimSpace(summary.CredentialState)
	result.SourceToken = translationJobSetupModelListSourceToken(result.PhaseID, spec.ID, credentialRef, settingsRequestToken)
	listed, err := service.providerSettings.ListProviderModels(ctx, ProviderSettingsModelListInput{
		ProviderID:            spec.ID,
		Endpoint:              cloneTranslationJobSetupOptionalString(summary.Endpoint),
		CredentialState:       strings.TrimSpace(summary.CredentialState),
		CredentialReferenceID: cloneTranslationJobSetupOptionalString(summary.CredentialReferenceID),
		RequestToken:          settingsRequestToken,
	})
	if err != nil {
		return ListTranslationJobSetupProviderModelsResult{}, fmt.Errorf("list translation job setup provider models via provider settings: %w", err)
	}
	result.CredentialStatus = strings.TrimSpace(listed.CredentialState)
	settingsRequestToken = strings.TrimSpace(listed.RequestToken)
	result.SourceToken = translationJobSetupModelListSourceToken(result.PhaseID, spec.ID, credentialRef, settingsRequestToken)
	result.Status = translationJobSetupMapProviderSettingsModelListState(listed.State)
	result.FailureKind = translationJobSetupMapProviderSettingsFailureKind(listed.FailureKind)
	result.Models = translationJobSetupProviderModelOptions(listed.Models)
	return result, nil
}

func (service *TranslationJobSetupService) listProviderModelsDirect(
	ctx context.Context,
	spec translationJobSetupProviderSpec,
	request ListTranslationJobSetupProviderModelsRequest,
	result ListTranslationJobSetupProviderModelsResult,
) (ListTranslationJobSetupProviderModelsResult, error) {
	credentialRef := translationJobSetupNormalizeCredentialRef(spec, request.CredentialRef)
	result.SourceToken = translationJobSetupModelListSourceToken(result.PhaseID, spec.ID, credentialRef, result.RequestToken)
	if !spec.CredentialRequired {
		models, listErr := service.requestProviderModels(ctx, spec.ID, "")
		if listErr != nil {
			result.Status = "failed"
			result.FailureKind = "model_list_failed"
			//nolint:nilerr // public contract intentionally returns a redacted failure instead of the transport error.
			return result, nil
		}
		result.Status = "credential_not_required"
		result.Models = models
		return result, nil
	}
	if credentialRef == "" || normalizeTranslationJobSetupField(result.CredentialStatus) == "missing" {
		result.CredentialStatus = "missing"
		result.Status = "credential_missing"
		result.FailureKind = "model_list_credential_missing"
		return result, nil
	}
	apiKey, resolved, resolveErr := service.loadCredentialSecret(ctx, credentialRef)
	if resolveErr != nil {
		result.Status = "failed"
		result.FailureKind = "model_list_failed"
		//nolint:nilerr // public contract intentionally returns a redacted failure instead of the secret-store error.
		return result, nil
	}
	if !resolved || strings.TrimSpace(apiKey) == "" {
		result.Status = "credential_missing"
		result.FailureKind = "model_list_credential_missing"
		return result, nil
	}
	models, listErr := service.requestProviderModels(ctx, spec.ID, apiKey)
	if listErr != nil {
		result.Status = "failed"
		result.FailureKind = "model_list_failed"
		//nolint:nilerr // public contract intentionally returns a redacted failure instead of the transport error.
		return result, nil
	}
	result.Status = "success"
	result.Models = models
	return result, nil
}

func translationJobSetupProviderModelOptions(models []ProviderSettingsModelOption) []TranslationJobSetupProviderModelOptionReadModel {
	options := make([]TranslationJobSetupProviderModelOptionReadModel, 0, len(models))
	for _, model := range models {
		options = append(options, TranslationJobSetupProviderModelOptionReadModel{
			ModelID: strings.TrimSpace(model.ModelID),
			Label:   strings.TrimSpace(model.Label),
		})
	}
	return options
}

// ValidateRequest classifies one setup request into blocking or creatable states.
func (service *TranslationJobSetupService) ValidateRequest(
	ctx context.Context,
	request TranslationJobSetupValidationRequest,
) (TranslationJobSetupValidationDecision, error) {
	validatedAt := service.now().UTC()
	if err := service.validateInputSource(ctx, request.InputSourceID); err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return blockingTranslationJobSetupValidation(validatedAt, "input_not_found", []string{"input"}), nil
		}
		return TranslationJobSetupValidationDecision{}, err
	}

	phaseRuntimes := normalizeTranslationJobSetupPhaseRuntimes(request.PhaseRuntimes)
	phaseResults := make([]TranslationJobSetupPhaseValidationReadModel, 0, len(translationJobSetupPhaseOrder))
	stalePhaseIDs := make([]string, 0)

	var firstFailure string
	var targetSlices []string
	for _, phaseID := range translationJobSetupPhaseOrder {
		runtime := phaseRuntimes[phaseID]
		result := service.validatePhaseRuntime(requestContextOrBackground(ctx), runtime)
		phaseResults = append(phaseResults, result)
		if result.IsModelSelectionStale {
			stalePhaseIDs = append(stalePhaseIDs, phaseID)
		}
		if result.BlockingFailureCategory != nil && firstFailure == "" {
			firstFailure = *result.BlockingFailureCategory
			targetSlices = translationJobSetupTargetSlicesForFailure(firstFailure)
		}
	}

	if firstFailure != "" {
		decision := blockingTranslationJobSetupValidation(validatedAt, firstFailure, targetSlices)
		decision.PhaseResults = phaseResults
		decision.StaleModelListPhaseIDs = stalePhaseIDs
		return decision, nil
	}

	return TranslationJobSetupValidationDecision{
		Status:                 translationJobSetupValidationStatusPass,
		ValidatedAt:            validatedAt,
		CanCreate:              true,
		PassSlices:             append([]string(nil), translationJobSetupAllSlices...),
		TargetSlices:           []string{},
		PhaseResults:           phaseResults,
		StaleModelListPhaseIDs: stalePhaseIDs,
	}, nil
}

// EvaluateCreateRequest blocks create until setup validation has passed.
func (service *TranslationJobSetupService) EvaluateCreateRequest(
	ctx context.Context,
	request TranslationJobSetupCreateRequest,
) (TranslationJobSetupCreateDecision, error) {
	if normalizeTranslationJobSetupField(request.ValidationStatus) != translationJobSetupValidationStatusPass {
		return TranslationJobSetupCreateDecision{CanCreate: false, ErrorKind: translationJobSetupErrorKindReadyRequired}, nil
	}
	if translationJobSetupValidationIsStale(service.now().UTC(), request.ValidatedAt.UTC()) {
		return TranslationJobSetupCreateDecision{CanCreate: false, ErrorKind: translationJobSetupErrorKindValidationStale}, nil
	}

	decision, err := service.ValidateRequest(ctx, TranslationJobSetupValidationRequest{
		InputSourceID: request.InputSourceID,
		PhaseRuntimes: request.PhaseRuntimes,
	})
	if err != nil {
		return TranslationJobSetupCreateDecision{}, err
	}
	if !decision.CanCreate {
		if decision.BlockingFailureCategory == nil {
			return TranslationJobSetupCreateDecision{CanCreate: false, ErrorKind: translationJobSetupErrorKindReadyRequired}, nil
		}
		return TranslationJobSetupCreateDecision{CanCreate: false, ErrorKind: *decision.BlockingFailureCategory}, nil
	}
	return TranslationJobSetupCreateDecision{
		CanCreate:            true,
		ValidationPassSlices: append([]string(nil), translationJobSetupAllSlices...),
	}, nil
}

// CreateTranslationJob creates a ready job and its runtime snapshots inside one transaction.
func (service *TranslationJobSetupService) CreateTranslationJob(
	ctx context.Context,
	request TranslationJobSetupCreateRequest,
	validationPassSlices []string,
) (TranslationJobSetupCreatedJobReadModel, error) {
	if service.jobLifecycleRepository == nil || service.translationSourceRepo == nil || service.transactor == nil {
		return TranslationJobSetupCreatedJobReadModel{}, fmt.Errorf("create translation job: persistence is not configured")
	}

	var created TranslationJobSetupCreatedJobReadModel
	err := service.transactor.WithTransaction(requestContextOrBackground(ctx), func(txCtx context.Context) error {
		result, err := service.createTranslationJobInTransaction(txCtx, request, validationPassSlices)
		if err != nil {
			return err
		}
		created = result
		return nil
	})
	if err != nil {
		return TranslationJobSetupCreatedJobReadModel{}, fmt.Errorf("create translation job transaction: %w", err)
	}
	return created, nil
}

func (service *TranslationJobSetupService) createTranslationJobInTransaction(
	txCtx context.Context,
	request TranslationJobSetupCreateRequest,
	validationPassSlices []string,
) (TranslationJobSetupCreatedJobReadModel, error) {
	if _, err := service.translationSourceRepo.GetXEditExtractedDataByID(txCtx, request.InputSourceID); err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return TranslationJobSetupCreatedJobReadModel{ErrorKind: translationJobSetupErrorKindInputNotFound}, nil
		}
		return TranslationJobSetupCreatedJobReadModel{}, fmt.Errorf("load translation input metadata: %w", err)
	}

	job, err := service.createReadyTranslationJob(txCtx, request.InputSourceID)
	if err != nil {
		return TranslationJobSetupCreatedJobReadModel{}, err
	}

	phaseRuntimes := normalizeTranslationJobSetupPhaseRuntimes(request.PhaseRuntimes)
	if service.providerSettings != nil {
		resolvedRuntimes, resolveErr := service.resolvePhaseRuntimeMapAgainstProviderSettings(txCtx, phaseRuntimes)
		if resolveErr != nil {
			return TranslationJobSetupCreatedJobReadModel{}, resolveErr
		}
		phaseRuntimes = resolvedRuntimes
	}
	wordRuntime := phaseRuntimes["word_translation"]
	if phaseRunErr := service.createInitialTranslationPhaseRun(txCtx, job.ID, wordRuntime); phaseRunErr != nil {
		return TranslationJobSetupCreatedJobReadModel{}, phaseRunErr
	}

	summaries, err := service.savePhaseRuntimeSnapshots(txCtx, job.ID, phaseRuntimes)
	if err != nil {
		return TranslationJobSetupCreatedJobReadModel{}, err
	}

	return TranslationJobSetupCreatedJobReadModel{
		JobID:                job.ID,
		JobState:             job.State,
		InputSource:          translationJobSetupInputSource,
		ErrorKind:            "",
		ValidationPassSlices: append([]string(nil), validationPassSlices...),
		ExecutionSummary: TranslationJobSetupExecutionSummaryReadModel{
			Provider:      wordRuntime.Provider,
			Model:         wordRuntime.Model,
			ExecutionMode: wordRuntime.ExecutionMode,
		},
		PhaseRuntimeSummaries: summaries,
	}, nil
}

func (service *TranslationJobSetupService) createReadyTranslationJob(
	ctx context.Context,
	inputSourceID int64,
) (repository.TranslationJob, error) {
	job, err := service.jobLifecycleRepository.CreateTranslationJob(ctx, repository.TranslationJobDraft{
		XEditExtractedDataID: inputSourceID,
		JobName:              translationJobSetupJobName(inputSourceID),
		State:                translationJobSetupJobStateReady,
		ProgressPercent:      0,
	})
	if err != nil {
		return repository.TranslationJob{}, fmt.Errorf("create translation job: %w", err)
	}
	return job, nil
}

func (service *TranslationJobSetupService) createInitialTranslationPhaseRun(
	ctx context.Context,
	jobID int64,
	wordRuntime TranslationJobSetupPhaseRuntimeDraftReadModel,
) error {
	if _, err := service.jobLifecycleRepository.CreateJobPhaseRun(ctx, repository.JobPhaseRunDraft{
		TranslationJobID: jobID,
		PhaseType:        "translation",
		State:            translationJobSetupPhaseStatePending,
		ExecutionOrder:   1,
		AIProvider:       wordRuntime.Provider,
		ModelName:        wordRuntime.Model,
		ExecutionMode:    wordRuntime.ExecutionMode,
		CredentialRef:    wordRuntime.CredentialRef,
		InstructionKind:  translationJobSetupInstructionKindWord,
	}); err != nil {
		return fmt.Errorf("create translation setup initial phase: %w", err)
	}
	return nil
}

func (service *TranslationJobSetupService) savePhaseRuntimeSnapshots(
	ctx context.Context,
	jobID int64,
	phaseRuntimes map[string]TranslationJobSetupPhaseRuntimeDraftReadModel,
) ([]TranslationJobSetupPhaseRuntimeSummaryReadModel, error) {
	snapshotStore, ok := service.jobLifecycleRepository.(translationJobSetupPhaseRuntimeSnapshotStore)
	if !ok {
		return nil, fmt.Errorf("create translation job phase runtime snapshot: snapshot store is not configured")
	}

	summaries := make([]TranslationJobSetupPhaseRuntimeSummaryReadModel, 0, len(translationJobSetupPhaseOrder))
	for _, phaseID := range translationJobSetupPhaseOrder {
		runtime := sanitizeTranslationJobSetupPhaseRuntime(phaseRuntimes[phaseID])
		runtime.PhaseID = phaseID
		if _, err := snapshotStore.SaveTranslationJobPhaseRuntimeSnapshot(ctx, repository.TranslationJobPhaseRuntimeSnapshotDraft{
			TranslationJobID:     jobID,
			PhaseID:              runtime.PhaseID,
			Provider:             runtime.Provider,
			ModelName:            runtime.Model,
			CredentialRef:        runtime.CredentialRef,
			CredentialStatus:     runtime.CredentialStatus,
			ExecutionMode:        runtime.ExecutionMode,
			BatchMode:            runtime.BatchMode,
			ModelListSourceToken: runtime.ModelListSourceToken,
		}); err != nil {
			return nil, fmt.Errorf("create translation job phase runtime snapshot: %w", err)
		}
		summaries = append(summaries, runtime)
	}
	return summaries, nil
}

// ReadSummary loads the ready job summary from the persisted job and snapshots.
func (service *TranslationJobSetupService) ReadSummary(
	ctx context.Context,
	jobID int64,
) (TranslationJobSetupSummaryReadModel, error) {
	if service.jobLifecycleRepository == nil {
		return TranslationJobSetupSummaryReadModel{}, fmt.Errorf("read translation job setup summary: persistence is not configured")
	}
	job, err := service.jobLifecycleRepository.GetTranslationJobByID(requestContextOrBackground(ctx), jobID)
	if err != nil {
		return TranslationJobSetupSummaryReadModel{}, fmt.Errorf("load translation job: %w", err)
	}
	snapshotStore, ok := service.jobLifecycleRepository.(translationJobSetupPhaseRuntimeSnapshotStore)
	if !ok {
		return TranslationJobSetupSummaryReadModel{}, fmt.Errorf("list translation job phase runtime snapshots: snapshot store is not configured")
	}
	snapshots, err := snapshotStore.ListTranslationJobPhaseRuntimeSnapshots(requestContextOrBackground(ctx), jobID)
	if err != nil {
		return TranslationJobSetupSummaryReadModel{}, fmt.Errorf("list translation job phase runtime snapshots: %w", err)
	}
	summaries := make([]TranslationJobSetupPhaseRuntimeSummaryReadModel, 0, len(snapshots))
	for _, snapshot := range snapshots {
		summaries = append(summaries, TranslationJobSetupPhaseRuntimeSummaryReadModel{
			PhaseID:              snapshot.PhaseID,
			Provider:             snapshot.Provider,
			Model:                snapshot.ModelName,
			CredentialRef:        snapshot.CredentialRef,
			CredentialStatus:     snapshot.CredentialStatus,
			ExecutionMode:        snapshot.ExecutionMode,
			BatchMode:            snapshot.BatchMode,
			ModelListSourceToken: snapshot.ModelListSourceToken,
		})
	}
	sort.SliceStable(summaries, func(i, j int) bool {
		return translationJobSetupPhaseIndex(summaries[i].PhaseID) < translationJobSetupPhaseIndex(summaries[j].PhaseID)
	})
	snapshotComplete := translationJobSetupHasAllPhaseSnapshots(summaries)
	wordRuntime := TranslationJobSetupPhaseRuntimeSummaryReadModel{}
	if snapshotComplete {
		wordRuntime = firstTranslationJobSetupPhaseSummary(summaries, "word_translation")
	}
	return TranslationJobSetupSummaryReadModel{
		JobID:         job.ID,
		JobState:      job.State,
		InputSource:   translationJobSetupInputSource,
		CanStartPhase: normalizeTranslationJobSetupField(job.State) == translationJobSetupJobStateReady && snapshotComplete,
		ExecutionSummary: TranslationJobSetupExecutionSummaryReadModel{
			Provider:      wordRuntime.Provider,
			Model:         wordRuntime.Model,
			ExecutionMode: wordRuntime.ExecutionMode,
		},
		ValidationPassSlices:  append([]string(nil), translationJobSetupAllSlices...),
		PhaseRuntimeSummaries: summaries,
	}, nil
}

// TranslationJobSetupPassSlices returns the canonical passing slices for Job Setup.
func TranslationJobSetupPassSlices() []string {
	return append([]string(nil), translationJobSetupAllSlices...)
}

func (service *TranslationJobSetupService) validateInputSource(ctx context.Context, inputSourceID int64) error {
	if inputSourceID <= 0 || service.translationSourceRepo == nil {
		return nil
	}
	_, err := service.translationSourceRepo.GetXEditExtractedDataByID(requestContextOrBackground(ctx), inputSourceID)
	if err != nil {
		return fmt.Errorf("load translation input metadata: %w", err)
	}
	return nil
}

func (service *TranslationJobSetupService) validatePhaseRuntime(
	ctx context.Context,
	runtime TranslationJobSetupPhaseRuntimeDraftReadModel,
) TranslationJobSetupPhaseValidationReadModel {
	sanitized := sanitizeTranslationJobSetupPhaseRuntime(runtime)
	result := TranslationJobSetupPhaseValidationReadModel{
		PhaseID:              sanitized.PhaseID,
		Status:               translationJobSetupValidationStatusPass,
		CanCreate:            true,
		ModelListState:       "success",
		ModelListSourceToken: sanitized.ModelListSourceToken,
	}

	spec, ok := translationJobSetupProviderCatalog[sanitized.Provider]
	if !ok || sanitized.Provider == "" || sanitized.Model == "" {
		result.Status = translationJobSetupValidationStatusFail
		result.CanCreate = false
		result.BlockingFailureCategory = stringPointer("phase_runtime_missing")
		result.ModelListState = "not_updated"
		return result
	}

	if service.providerSettings != nil {
		resolved, err := service.resolvePhaseRuntimeAgainstProviderSettings(ctx, sanitized)
		if err != nil {
			result.Status = translationJobSetupValidationStatusFail
			result.CanCreate = false
			result.BlockingFailureCategory = stringPointer("required_setting_missing")
			result.ModelListState = "failed"
			return result
		}
		sanitized = resolved
		result.ModelListSourceToken = sanitized.ModelListSourceToken
	}

	if spec.CredentialRequired &&
		(sanitized.CredentialRef == "" || normalizeTranslationJobSetupField(sanitized.CredentialStatus) == "missing") {
		result.Status = translationJobSetupValidationStatusFail
		result.CanCreate = false
		result.BlockingFailureCategory = stringPointer("credential_missing")
		result.ModelListState = "credential_missing"
		return result
	}

	if !spec.CredentialRequired {
		result.ModelListState = "credential_not_required"
	}

	if sanitized.ModelListSourceToken == "" || !strings.HasPrefix(sanitized.ModelListSourceToken, translationJobSetupModelListSourcePrefix(sanitized.PhaseID, sanitized.Provider, sanitized.CredentialRef)) {
		result.Status = translationJobSetupValidationStatusFail
		result.CanCreate = false
		result.BlockingFailureCategory = stringPointer("model_selection_stale")
		result.IsModelSelectionStale = true
		if spec.CredentialRequired {
			result.ModelListState = "failed"
		}
		return result
	}

	return result
}

func blockingTranslationJobSetupValidation(validatedAt time.Time, failure string, targetSlices []string) TranslationJobSetupValidationDecision {
	return TranslationJobSetupValidationDecision{
		Status:                  translationJobSetupValidationStatusFail,
		BlockingFailureCategory: stringPointer(failure),
		TargetSlices:            append([]string(nil), targetSlices...),
		ValidatedAt:             validatedAt,
		CanCreate:               false,
		PassSlices:              []string{},
	}
}

func translationJobSetupTargetSlicesForFailure(failure string) []string {
	switch normalizeTranslationJobSetupField(failure) {
	case "credential_missing", "model_list_credential_missing":
		return []string{"credentials"}
	case "input_not_found":
		return []string{"input"}
	default:
		return []string{"runtime"}
	}
}

func (service *TranslationJobSetupService) currentOptionsReadModel(
	ctx context.Context,
) (TranslationJobSetupOptionsReadModel, error) {
	if translationJobSetupOptionsOverrideExists(service.optionsReadModel) {
		return cloneTranslationJobSetupOptionsReadModel(service.optionsReadModel), nil
	}
	inputCandidates, existingJob, err := service.loadTranslationInputCandidates(ctx)
	if err != nil {
		return TranslationJobSetupOptionsReadModel{}, err
	}
	sharedDictionaries, err := service.loadSharedDictionaryOptions(ctx)
	if err != nil {
		return TranslationJobSetupOptionsReadModel{}, err
	}
	sharedPersonas, err := service.loadSharedPersonaOptions(ctx)
	if err != nil {
		return TranslationJobSetupOptionsReadModel{}, err
	}
	return TranslationJobSetupOptionsReadModel{
		InputCandidates:      inputCandidates,
		ExistingJob:          existingJob,
		SharedDictionaries:   sharedDictionaries,
		SharedPersonas:       sharedPersonas,
		AIRuntimeOptions:     translationJobSetupRuntimeOptions(),
		CredentialRefs:       service.translationJobSetupCredentialRefs(ctx),
		ProviderCapabilities: translationJobSetupProviderCapabilities(),
		PhaseRuntimeDrafts:   translationJobSetupEmptyPhaseRuntimeDrafts(),
	}, nil
}

func (service *TranslationJobSetupService) loadTranslationInputCandidates(
	ctx context.Context,
) ([]TranslationJobSetupInputCandidateReadModel, *TranslationJobSetupExistingJobReadModel, error) {
	lister, ok := service.translationSourceRepo.(translationJobSetupTranslationSourceLister)
	if !ok {
		return nil, nil, nil
	}
	inputs, err := lister.ListXEditExtractedData(ctx)
	if err != nil {
		return nil, nil, fmt.Errorf("list translation job setup input candidates: %w", err)
	}
	result := make([]TranslationJobSetupInputCandidateReadModel, 0, len(inputs))
	var existingJob *TranslationJobSetupExistingJobReadModel
	for _, input := range inputs {
		jobForInput, err := service.loadExistingJobForInput(ctx, input.ID)
		if err != nil {
			return nil, nil, err
		}
		if jobForInput != nil {
			if existingJob == nil {
				existingJob = cloneTranslationJobSetupExistingJobReadModel(jobForInput)
			}
			continue
		}
		result = append(result, TranslationJobSetupInputCandidateReadModel{
			ID:           input.ID,
			Label:        translationJobSetupInputCandidateLabel(input),
			SourceKind:   translationJobSetupInputSource,
			RecordCount:  input.RecordCount,
			RegisteredAt: input.ImportedAt,
		})
	}
	return result, existingJob, nil
}

func (service *TranslationJobSetupService) loadExistingJobForInput(
	ctx context.Context,
	inputSourceID int64,
) (*TranslationJobSetupExistingJobReadModel, error) {
	loader, ok := service.translationSourceRepo.(translationJobSetupExistingJobLoader)
	if !ok {
		return nil, nil
	}
	job, err := loader.GetExistingTranslationJob(ctx, inputSourceID)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, nil
		}
		return nil, fmt.Errorf("load existing translation job: %w", err)
	}
	return &TranslationJobSetupExistingJobReadModel{
		InputSourceID: inputSourceID,
		JobID:         job.ID,
		Status:        job.State,
		InputSource:   translationJobSetupInputSource,
	}, nil
}

func translationJobSetupInputCandidateLabel(input repository.XEditExtractedData) string {
	if strings.TrimSpace(input.TargetPluginName) != "" {
		return strings.TrimSpace(input.TargetPluginName)
	}
	if strings.TrimSpace(input.SourceFilePath) != "" {
		return filepath.Base(strings.TrimSpace(input.SourceFilePath))
	}
	return fmt.Sprintf("input-%d", input.ID)
}

func (service *TranslationJobSetupService) loadSharedDictionaryOptions(
	ctx context.Context,
) ([]TranslationJobSetupDictionaryOptionReadModel, error) {
	if service.masterDictionaryRepo == nil {
		return nil, nil
	}
	result, err := service.masterDictionaryRepo.List(ctx, MasterDictionaryQuery{Page: 1, PageSize: 100})
	if err != nil {
		return nil, fmt.Errorf("list translation job setup shared dictionaries: %w", err)
	}
	options := make([]TranslationJobSetupDictionaryOptionReadModel, 0, len(result.Items))
	for _, entry := range result.Items {
		label := strings.TrimSpace(entry.Source)
		if label == "" {
			label = strconv.FormatInt(entry.ID, 10)
		}
		options = append(options, TranslationJobSetupDictionaryOptionReadModel{ID: strconv.FormatInt(entry.ID, 10), Label: label})
	}
	return options, nil
}

func (service *TranslationJobSetupService) loadSharedPersonaOptions(
	ctx context.Context,
) ([]TranslationJobSetupPersonaOptionReadModel, error) {
	if service.masterPersonaRepo == nil {
		return nil, nil
	}
	result, err := service.masterPersonaRepo.List(ctx, repository.MasterPersonaListQuery{Page: 1, PageSize: 100})
	if err != nil {
		return nil, fmt.Errorf("list translation job setup shared personas: %w", err)
	}
	options := make([]TranslationJobSetupPersonaOptionReadModel, 0, len(result.Items))
	for _, entry := range result.Items {
		label := strings.TrimSpace(entry.DisplayName)
		if label == "" {
			label = strings.TrimSpace(entry.IdentityKey)
		}
		options = append(options, TranslationJobSetupPersonaOptionReadModel{ID: entry.IdentityKey, Label: label})
	}
	return options, nil
}

func translationJobSetupRuntimeOptions() []TranslationJobSetupRuntimeOptionReadModel {
	result := make([]TranslationJobSetupRuntimeOptionReadModel, 0, len(translationJobSetupUserFacingProviderIDs))
	for _, providerID := range translationJobSetupUserFacingProviderIDs {
		spec := translationJobSetupProviderCatalog[providerID]
		result = append(result, TranslationJobSetupRuntimeOptionReadModel{
			Provider: spec.ID,
			Model:    spec.DefaultModel,
			Mode:     spec.SupportedModes[0],
		})
	}
	return result
}

func translationJobSetupProviderCapabilities() []TranslationJobSetupProviderCapabilityReadModel {
	result := make([]TranslationJobSetupProviderCapabilityReadModel, 0, len(translationJobSetupUserFacingProviderIDs))
	for _, providerID := range translationJobSetupUserFacingProviderIDs {
		spec := translationJobSetupProviderCatalog[providerID]
		requirement := "required"
		if !spec.CredentialRequired {
			requirement = "not_required"
		}
		result = append(result, TranslationJobSetupProviderCapabilityReadModel{
			Provider:                spec.ID,
			CredentialRequirement:   requirement,
			SupportedExecutionModes: append([]string(nil), spec.SupportedModes...),
			SupportsBatchMode:       spec.SupportsBatchMode,
		})
	}
	return result
}

func (service *TranslationJobSetupService) translationJobSetupCredentialRefs(
	ctx context.Context,
) []TranslationJobSetupCredentialReferenceReadModel {
	if service.providerSettings != nil {
		summaries, err := providerSettingsSummaryMap(requestContextOrBackground(ctx), service.providerSettings)
		if err == nil && summaries != nil {
			result := make([]TranslationJobSetupCredentialReferenceReadModel, 0, len(translationJobSetupUserFacingProviderIDs))
			for _, providerID := range translationJobSetupUserFacingProviderIDs {
				spec := translationJobSetupProviderCatalog[providerID]
				summary, ok := summaries[providerID]
				configured := !spec.CredentialRequired
				missingSecret := spec.CredentialRequired
				credentialRef := spec.DefaultCredentialRef
				if ok {
					credentialRef = strings.TrimSpace(pointerStringValue(summary.CredentialReferenceID))
					configured = strings.TrimSpace(summary.CredentialState) == providerSettingsCredentialStateConfigured ||
						strings.TrimSpace(summary.CredentialState) == providerSettingsCredentialStateNotRequired
					missingSecret = strings.TrimSpace(summary.CredentialState) == providerSettingsCredentialStateMissing
				}
				result = append(result, TranslationJobSetupCredentialReferenceReadModel{
					Provider:        spec.ID,
					CredentialRef:   credentialRef,
					IsConfigured:    configured,
					IsMissingSecret: missingSecret,
				})
			}
			return result
		}
	}
	result := make([]TranslationJobSetupCredentialReferenceReadModel, 0, len(translationJobSetupUserFacingProviderIDs))
	for _, providerID := range translationJobSetupUserFacingProviderIDs {
		spec := translationJobSetupProviderCatalog[providerID]
		result = append(result, TranslationJobSetupCredentialReferenceReadModel{
			Provider:        spec.ID,
			CredentialRef:   spec.DefaultCredentialRef,
			IsConfigured:    !spec.CredentialRequired,
			IsMissingSecret: spec.CredentialRequired,
		})
	}
	return result
}

func (service *TranslationJobSetupService) loadCredentialSecret(
	ctx context.Context,
	key string,
) (string, bool, error) {
	if service == nil || service.secretStore == nil || strings.TrimSpace(key) == "" {
		return "", false, nil
	}
	timeout := service.secretLoadTimeout
	if timeout <= 0 {
		timeout = translationJobSetupSecretTimeout
	}
	requestContext := requestContextOrBackground(ctx)
	resultChannel := make(chan struct {
		value string
		err   error
	}, 1)
	go func() {
		value, err := service.secretStore.Load(requestContext, key)
		resultChannel <- struct {
			value string
			err   error
		}{value: value, err: err}
	}()
	timer := time.NewTimer(timeout)
	defer timer.Stop()
	select {
	case <-requestContext.Done():
		return "", false, nil
	case <-timer.C:
		return "", false, nil
	case result := <-resultChannel:
		return result.value, true, result.err
	}
}

func normalizeTranslationJobSetupPhaseRuntimes(
	phaseRuntimes []TranslationJobSetupPhaseRuntimeDraftReadModel,
) map[string]TranslationJobSetupPhaseRuntimeDraftReadModel {
	result := make(map[string]TranslationJobSetupPhaseRuntimeDraftReadModel, len(translationJobSetupPhaseOrder))
	for _, phaseID := range translationJobSetupPhaseOrder {
		result[phaseID] = TranslationJobSetupPhaseRuntimeDraftReadModel{PhaseID: phaseID, BatchMode: "unsupported"}
	}
	for _, runtime := range phaseRuntimes {
		phaseID := normalizeTranslationJobSetupField(runtime.PhaseID)
		if phaseID == "" {
			continue
		}
		sanitized := sanitizeTranslationJobSetupPhaseRuntime(runtime)
		sanitized.PhaseID = phaseID
		result[phaseID] = sanitized
	}
	return result
}

func sanitizeTranslationJobSetupPhaseRuntime(runtime TranslationJobSetupPhaseRuntimeDraftReadModel) TranslationJobSetupPhaseRuntimeDraftReadModel {
	sanitized := TranslationJobSetupPhaseRuntimeDraftReadModel{
		PhaseID:              normalizeTranslationJobSetupField(runtime.PhaseID),
		Provider:             normalizeTranslationJobSetupField(runtime.Provider),
		Model:                strings.TrimSpace(runtime.Model),
		CredentialRef:        strings.TrimSpace(runtime.CredentialRef),
		CredentialStatus:     normalizeTranslationJobSetupField(runtime.CredentialStatus),
		ExecutionMode:        normalizeTranslationJobSetupField(runtime.ExecutionMode),
		BatchMode:            normalizeTranslationJobSetupField(runtime.BatchMode),
		ModelListSourceToken: strings.TrimSpace(runtime.ModelListSourceToken),
	}
	spec, ok := translationJobSetupProviderCatalog[sanitized.Provider]
	if !ok {
		return sanitized
	}
	sanitized.CredentialRef = translationJobSetupNormalizeCredentialRef(spec, sanitized.CredentialRef)
	sanitized.CredentialStatus = translationJobSetupCredentialStatus(spec, sanitized.CredentialStatus, sanitized.CredentialRef)
	if !spec.SupportsBatchMode {
		sanitized.BatchMode = "unsupported"
		sanitized.ExecutionMode = "sync"
	} else {
		if sanitized.BatchMode == "enabled" {
			sanitized.ExecutionMode = "batch"
		} else {
			sanitized.BatchMode = "disabled"
			sanitized.ExecutionMode = "sync"
		}
	}
	if sanitized.ExecutionMode == "" {
		sanitized.ExecutionMode = spec.SupportedModes[0]
	}
	if !spec.CredentialRequired {
		sanitized.CredentialStatus = "not_required"
		sanitized.CredentialRef = ""
	}
	return sanitized
}

func translationJobSetupCredentialStatus(spec translationJobSetupProviderSpec, requested string, credentialRef string) string {
	if !spec.CredentialRequired {
		return "not_required"
	}
	if normalizeTranslationJobSetupField(requested) == "configured" &&
		translationJobSetupNormalizeCredentialRef(spec, credentialRef) != "" {
		return "configured"
	}
	return "missing"
}

func translationJobSetupNormalizeCredentialRef(
	spec translationJobSetupProviderSpec,
	credentialRef string,
) string {
	if !spec.CredentialRequired {
		return ""
	}
	normalizedRef := strings.TrimSpace(credentialRef)
	if normalizedRef == "" {
		return ""
	}
	if normalizedRef != strings.TrimSpace(spec.DefaultCredentialRef) {
		return ""
	}
	return normalizedRef
}

func translationJobSetupPhaseIndex(phaseID string) int {
	for index, candidate := range translationJobSetupPhaseOrder {
		if candidate == normalizeTranslationJobSetupField(phaseID) {
			return index
		}
	}
	return len(translationJobSetupPhaseOrder)
}

func firstTranslationJobSetupPhaseSummary(
	summaries []TranslationJobSetupPhaseRuntimeSummaryReadModel,
	phaseID string,
) TranslationJobSetupPhaseRuntimeSummaryReadModel {
	for _, summary := range summaries {
		if normalizeTranslationJobSetupField(summary.PhaseID) == normalizeTranslationJobSetupField(phaseID) {
			return summary
		}
	}
	return TranslationJobSetupPhaseRuntimeSummaryReadModel{}
}

func translationJobSetupJobName(inputSourceID int64) string {
	return fmt.Sprintf("translation-job-%d", inputSourceID)
}

func translationJobSetupValidationIsStale(now time.Time, validatedAt time.Time) bool {
	if validatedAt.IsZero() {
		return true
	}
	return validatedAt.Before(translationJobSetupValidationFreshnessCutoff(now))
}

func translationJobSetupValidationFreshnessCutoff(now time.Time) time.Time {
	nowUTC := now.UTC()
	cutoff := time.Date(nowUTC.Year(), nowUTC.Month(), nowUTC.Day(), translationJobSetupValidationFreshnessCutoffHourUTC, 0, 0, 0, time.UTC)
	if nowUTC.Before(cutoff) {
		return cutoff.AddDate(0, 0, -1)
	}
	return cutoff
}

func translationJobSetupModelListSourcePrefix(phaseID string, provider string, credentialRef string) string {
	return fmt.Sprintf("%s|%s|%s|", normalizeTranslationJobSetupField(phaseID), normalizeTranslationJobSetupField(provider), strings.TrimSpace(credentialRef))
}

func translationJobSetupModelListSourceToken(phaseID string, provider string, credentialRef string, requestToken string) string {
	return translationJobSetupModelListSourcePrefix(phaseID, provider, credentialRef) + strings.TrimSpace(requestToken)
}

func translationJobSetupEmptyPhaseRuntimeDrafts() []TranslationJobSetupPhaseRuntimeDraftReadModel {
	result := make([]TranslationJobSetupPhaseRuntimeDraftReadModel, 0, len(translationJobSetupPhaseOrder))
	for _, phaseID := range translationJobSetupPhaseOrder {
		result = append(result, TranslationJobSetupPhaseRuntimeDraftReadModel{
			PhaseID:          phaseID,
			CredentialStatus: "missing",
			ExecutionMode:    "sync",
			BatchMode:        "unsupported",
		})
	}
	return result
}

func normalizeTranslationJobSetupField(value string) string {
	return strings.ToLower(strings.TrimSpace(value))
}

func translationJobSetupOptionsOverrideExists(readModel TranslationJobSetupOptionsReadModel) bool {
	return len(readModel.InputCandidates) > 0 ||
		readModel.ExistingJob != nil ||
		len(readModel.SharedDictionaries) > 0 ||
		len(readModel.SharedPersonas) > 0 ||
		len(readModel.AIRuntimeOptions) > 0 ||
		len(readModel.CredentialRefs) > 0 ||
		len(readModel.ProviderCapabilities) > 0 ||
		len(readModel.PhaseRuntimeDrafts) > 0
}

func cloneTranslationJobSetupOptionsReadModel(readModel TranslationJobSetupOptionsReadModel) TranslationJobSetupOptionsReadModel {
	inputCandidates := make([]TranslationJobSetupInputCandidateReadModel, 0, len(readModel.InputCandidates))
	for _, candidate := range readModel.InputCandidates {
		clonedCandidate := candidate
		clonedCandidate.ExistingJob = cloneTranslationJobSetupExistingJobReadModel(candidate.ExistingJob)
		inputCandidates = append(inputCandidates, clonedCandidate)
	}
	cloned := TranslationJobSetupOptionsReadModel{
		InputCandidates:      inputCandidates,
		SharedDictionaries:   append([]TranslationJobSetupDictionaryOptionReadModel(nil), readModel.SharedDictionaries...),
		SharedPersonas:       append([]TranslationJobSetupPersonaOptionReadModel(nil), readModel.SharedPersonas...),
		AIRuntimeOptions:     append([]TranslationJobSetupRuntimeOptionReadModel(nil), readModel.AIRuntimeOptions...),
		CredentialRefs:       append([]TranslationJobSetupCredentialReferenceReadModel(nil), readModel.CredentialRefs...),
		ProviderCapabilities: append([]TranslationJobSetupProviderCapabilityReadModel(nil), readModel.ProviderCapabilities...),
		PhaseRuntimeDrafts:   append([]TranslationJobSetupPhaseRuntimeDraftReadModel(nil), readModel.PhaseRuntimeDrafts...),
	}
	if readModel.ExistingJob != nil {
		existingJob := *readModel.ExistingJob
		cloned.ExistingJob = &existingJob
	}
	return cloned
}

func cloneTranslationJobSetupExistingJobReadModel(
	existingJob *TranslationJobSetupExistingJobReadModel,
) *TranslationJobSetupExistingJobReadModel {
	if existingJob == nil {
		return nil
	}
	cloned := *existingJob
	return &cloned
}

// TranslationJobSetupReadOptions returns the current read-only Job Setup page model.
func TranslationJobSetupReadOptions() TranslationJobSetupOptionsReadModel {
	return TranslationJobSetupOptionsReadModel{
		InputCandidates: []TranslationJobSetupInputCandidateReadModel{{
			ID:          44,
			Label:       "Dialogues import",
			SourceKind:  translationJobSetupInputSource,
			RecordCount: 120,
		}},
		ExistingJob: &TranslationJobSetupExistingJobReadModel{
			InputSourceID: 999,
			JobID:         88,
			Status:        translationJobSetupJobStateReady,
			InputSource:   translationJobSetupInputSource,
		},
		SharedDictionaries: []TranslationJobSetupDictionaryOptionReadModel{{ID: "dict-core", Label: "Core Dictionary"}},
		SharedPersonas:     []TranslationJobSetupPersonaOptionReadModel{{ID: "persona-guard", Label: "Guard Persona"}},
		AIRuntimeOptions:   translationJobSetupRuntimeOptions(),
		CredentialRefs: []TranslationJobSetupCredentialReferenceReadModel{
			{Provider: translationJobSetupProviderGemini, CredentialRef: translationJobSetupCredentialRefGeminiPrimary, IsConfigured: false, IsMissingSecret: true},
			{Provider: translationJobSetupProviderLM, CredentialRef: "lmstudio-local", IsConfigured: true, IsMissingSecret: false},
			{Provider: translationJobSetupProviderXAI, CredentialRef: "xai-primary", IsConfigured: true, IsMissingSecret: false},
		},
		ProviderCapabilities: translationJobSetupProviderCapabilities(),
		PhaseRuntimeDrafts:   translationJobSetupEmptyPhaseRuntimeDrafts(),
	}
}

// TranslationJobSetupReadSummary returns the read-only re-display for one created job.
func TranslationJobSetupReadSummary(jobID int64) TranslationJobSetupSummaryReadModel {
	word := TranslationJobSetupPhaseRuntimeSummaryReadModel{
		PhaseID:              "word_translation",
		Provider:             translationJobSetupProviderGemini,
		Model:                "gemini-2.5-pro",
		CredentialRef:        translationJobSetupCredentialRefGeminiPrimary,
		CredentialStatus:     "configured",
		ExecutionMode:        "sync",
		BatchMode:            "disabled",
		ModelListSourceToken: translationJobSetupModelListSourceToken("word_translation", translationJobSetupProviderGemini, translationJobSetupCredentialRefGeminiPrimary, "bootstrap"),
	}
	return TranslationJobSetupSummaryReadModel{
		JobID:                 jobID,
		JobState:              translationJobSetupJobStateReady,
		InputSource:           translationJobSetupInputSource,
		CanStartPhase:         true,
		ExecutionSummary:      TranslationJobSetupExecutionSummaryReadModel{Provider: word.Provider, Model: word.Model, ExecutionMode: word.ExecutionMode},
		ValidationPassSlices:  append([]string(nil), translationJobSetupAllSlices...),
		PhaseRuntimeSummaries: []TranslationJobSetupPhaseRuntimeSummaryReadModel{word},
	}
}

func requestContextOrBackground(ctx context.Context) context.Context {
	if ctx != nil {
		return ctx
	}
	return translationJobSetupDefaultContext
}

func stringPointer(value string) *string {
	if strings.TrimSpace(value) == "" {
		return nil
	}
	cloned := value
	return &cloned
}

func (service *TranslationJobSetupService) requestProviderModels(
	ctx context.Context,
	providerID string,
	apiKey string,
) ([]TranslationJobSetupProviderModelOptionReadModel, error) {
	if service == nil || service.providerModelListLoader == nil {
		return nil, fmt.Errorf("provider model list loader is not configured")
	}
	models, err := service.providerModelListLoader.ListProviderModels(requestContextOrBackground(ctx), providerID, apiKey)
	if err != nil {
		return nil, fmt.Errorf("load provider model list: %w", err)
	}
	result := make([]TranslationJobSetupProviderModelOptionReadModel, 0, len(models))
	for _, model := range models {
		result = append(result, TranslationJobSetupProviderModelOptionReadModel{
			ModelID: strings.TrimSpace(model.ModelID),
			Label:   strings.TrimSpace(model.Label),
		})
	}
	return result, nil
}

func (service *TranslationJobSetupService) resolvePhaseRuntimeMapAgainstProviderSettings(
	ctx context.Context,
	phaseRuntimes map[string]TranslationJobSetupPhaseRuntimeDraftReadModel,
) (map[string]TranslationJobSetupPhaseRuntimeDraftReadModel, error) {
	resolved := make(map[string]TranslationJobSetupPhaseRuntimeDraftReadModel, len(phaseRuntimes))
	for phaseID, runtime := range phaseRuntimes {
		nextRuntime := sanitizeTranslationJobSetupPhaseRuntime(runtime)
		var err error
		if translationJobSetupProviderUsesProviderSettings(nextRuntime.Provider) {
			nextRuntime, err = service.resolvePhaseRuntimeAgainstProviderSettings(ctx, runtime)
		}
		if err != nil {
			return nil, err
		}
		resolved[phaseID] = nextRuntime
	}
	return resolved, nil
}

func (service *TranslationJobSetupService) resolvePhaseRuntimeAgainstProviderSettings(
	ctx context.Context,
	runtime TranslationJobSetupPhaseRuntimeDraftReadModel,
) (TranslationJobSetupPhaseRuntimeDraftReadModel, error) {
	if service.providerSettings == nil {
		return sanitizeTranslationJobSetupPhaseRuntime(runtime), nil
	}
	sanitized := sanitizeTranslationJobSetupPhaseRuntime(runtime)
	spec, ok := translationJobSetupProviderCatalog[sanitized.Provider]
	if !ok {
		return sanitized, nil
	}
	resolved, err := service.providerSettings.ResolveProviderExecutionSettings(requestContextOrBackground(ctx), ProviderSettingsResolveInput{
		ConsumerID: "translation_job_setup",
		Selection: ProviderSettingsResolveSelection{
			ProviderID:      spec.ID,
			Model:           sanitized.Model,
			ExecutionMethod: sanitized.ExecutionMode,
			UseBatchAPI:     sanitized.BatchMode == "enabled",
		},
	})
	if err != nil {
		return TranslationJobSetupPhaseRuntimeDraftReadModel{}, fmt.Errorf("resolve translation job setup provider settings: %w", err)
	}
	sanitized.CredentialRef = strings.TrimSpace(pointerStringValue(resolved.CredentialReferenceID))
	sanitized.CredentialStatus = strings.TrimSpace(resolved.CredentialState)
	sanitized.ModelListSourceToken = translationJobSetupModelListSourceToken(
		sanitized.PhaseID,
		sanitized.Provider,
		sanitized.CredentialRef,
		pointerStringValue(resolved.RequestToken),
	)
	if resolved.ErrorKind != nil {
		switch strings.TrimSpace(*resolved.ErrorKind) {
		case providerSettingsErrorKindCredentialMissing:
			sanitized.CredentialStatus = "missing"
		case providerSettingsErrorKindEndpointMissing:
			return sanitized, fmt.Errorf("translation job setup provider settings endpoint is missing")
		}
	}
	return sanitized, nil
}

func translationJobSetupMapProviderSettingsModelListState(state string) string {
	switch strings.TrimSpace(state) {
	case providerSettingsModelListStateCredentialNotNeeded:
		return "credential_not_required"
	case providerSettingsModelListStateCredentialMissing:
		return "credential_missing"
	case providerSettingsModelListStateReady:
		return "success"
	case providerSettingsModelListStateFailed:
		return "failed"
	default:
		return "not_updated"
	}
}

func translationJobSetupMapProviderSettingsFailureKind(kind *string) string {
	if kind == nil {
		return ""
	}
	switch strings.TrimSpace(*kind) {
	case providerSettingsErrorKindCredentialMissing:
		return "model_list_credential_missing"
	case providerSettingsErrorKindEndpointMissing:
		return "required_setting_missing"
	case providerSettingsErrorKindModelListFailed:
		return "model_list_failed"
	default:
		return strings.TrimSpace(*kind)
	}
}

func translationJobSetupProviderUsesProviderSettings(providerID string) bool {
	switch normalizeTranslationJobSetupField(providerID) {
	case translationJobSetupProviderGemini, translationJobSetupProviderLM, translationJobSetupProviderXAI:
		return true
	default:
		return false
	}
}

func cloneTranslationJobSetupOptionalString(value *string) *string {
	if value == nil {
		return nil
	}
	trimmed := strings.TrimSpace(*value)
	if trimmed == "" {
		return nil
	}
	cloned := trimmed
	return &cloned
}

func pointerStringValue(value *string) string {
	if value == nil {
		return ""
	}
	return strings.TrimSpace(*value)
}

func translationJobSetupHasAllPhaseSnapshots(
	summaries []TranslationJobSetupPhaseRuntimeSummaryReadModel,
) bool {
	if len(summaries) < len(translationJobSetupPhaseOrder) {
		return false
	}
	seen := make(map[string]struct{}, len(summaries))
	for _, summary := range summaries {
		phaseID := strings.TrimSpace(summary.PhaseID)
		if phaseID == "" {
			continue
		}
		seen[phaseID] = struct{}{}
	}
	for _, phaseID := range translationJobSetupPhaseOrder {
		if _, ok := seen[phaseID]; !ok {
			return false
		}
	}
	return true
}
