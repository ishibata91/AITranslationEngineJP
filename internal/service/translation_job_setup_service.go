package service

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"time"

	"aitranslationenginejp/internal/repository"
)

const (
	translationJobSetupValidationStatusPass         = "pass"
	translationJobSetupValidationStatusFail         = "fail"
	translationJobSetupErrorKindValidationFailed    = "validation_failed"
	translationJobSetupErrorKindValidationStale     = "validation_stale"
	translationJobSetupErrorKindDuplicateInput      = "duplicate_input"
	translationJobSetupErrorKindCacheMissing        = "cache_missing"
	translationJobSetupErrorKindInputNotFound       = "input_not_found"
	translationJobSetupErrorKindProviderUnreachable = "provider_unreachable"

	translationJobSetupValidationFreshnessCutoffHourUTC = 9

	translationJobSetupJobStateReady          = "ready"
	translationJobSetupPhaseStatePending      = "pending"
	translationJobSetupPhaseTypeTranslation   = "translation"
	translationJobSetupInstructionKindDefault = "default"
	translationJobSetupInputSourceTranslation = "translation_input"

	translationJobSetupRealProviderOpenAI = "openai"
	translationJobSetupModelGPT54Mini     = "gpt-5.4-mini"

	translationJobSetupBlockingFailureRequiredSettingMissing  = "required_setting_missing"
	translationJobSetupBlockingFailureInputNotFound           = "input_not_found"
	translationJobSetupBlockingFailureFoundationRefMissing    = "foundation_ref_missing"
	translationJobSetupBlockingFailureProviderModeUnsupported = "provider_mode_unsupported"
	translationJobSetupBlockingFailureCacheMissing            = "cache_missing"
	translationJobSetupBlockingFailureProviderUnreachable     = "provider_unreachable"

	translationJobSetupProviderReachabilityPrompt = "ping"
	translationJobSetupOpenAIBaseURL              = "https://api.openai.com/v1"
	translationJobSetupXAIBaseURLEnv              = "AITRANSLATIONENGINEJP_MASTER_PERSONA_XAI_BASE_URL"
	translationJobSetupLMStudioBaseURLEnv         = "AITRANSLATIONENGINEJP_MASTER_PERSONA_LM_STUDIO_BASE_URL"
)

var translationJobSetupAllSlices = []string{"input", "runtime", "credentials"}

var translationJobSetupSupportedProviderSet = map[string]struct{}{
	translationJobSetupRealProviderOpenAI: {},
	MasterPersonaProviderGemini:           {},
	MasterPersonaProviderLMStudio:         {},
	MasterPersonaProviderXAI:              {},
}

// TranslationJobSetupValidationRequest carries the validation inputs needed by the service layer.
type TranslationJobSetupValidationRequest struct {
	InputSourceID int64
	Provider      string
	Model         string
	ExecutionMode string
	CredentialRef string
}

// TranslationJobSetupValidationDecision returns the backend validation outcome.
type TranslationJobSetupValidationDecision struct {
	Status                  string
	BlockingFailureCategory *string
	TargetSlices            []string
	ValidatedAt             time.Time
	CanCreate               bool
	PassSlices              []string
}

// TranslationJobSetupCreateRequest carries the create gating inputs needed by the service layer.
type TranslationJobSetupCreateRequest struct {
	InputSourceID    int64
	ValidationStatus string
	ValidatedAt      time.Time
	Provider         string
	Model            string
	ExecutionMode    string
	CredentialRef    string
}

// TranslationJobSetupCreateDecision returns whether create may proceed.
type TranslationJobSetupCreateDecision struct {
	CanCreate            bool
	ErrorKind            string
	ValidationPassSlices []string
}

// TranslationJobSetupCreatedJobReadModel stores the created ready-job response.
type TranslationJobSetupCreatedJobReadModel struct {
	JobID                int64
	JobState             string
	InputSource          string
	ExecutionSummary     TranslationJobSetupExecutionSummaryReadModel
	ValidationPassSlices []string
	ErrorKind            string
}

// TranslationJobSetupOptionsReadModel stores the read-only page data needed by Job Setup.
type TranslationJobSetupOptionsReadModel struct {
	InputCandidates    []TranslationJobSetupInputCandidateReadModel
	ExistingJob        *TranslationJobSetupExistingJobReadModel
	SharedDictionaries []TranslationJobSetupDictionaryOptionReadModel
	SharedPersonas     []TranslationJobSetupPersonaOptionReadModel
	AIRuntimeOptions   []TranslationJobSetupRuntimeOptionReadModel
	CredentialRefs     []TranslationJobSetupCredentialReferenceReadModel
}

// TranslationJobSetupInputCandidateReadModel is one selectable translation input source.
type TranslationJobSetupInputCandidateReadModel struct {
	ID           int64
	Label        string
	SourceKind   string
	RecordCount  int
	RegisteredAt time.Time
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

// TranslationJobSetupSummaryReadModel stores the read-only job display.
type TranslationJobSetupSummaryReadModel struct {
	JobID                int64
	JobState             string
	InputSource          string
	CanStartPhase        bool
	ExecutionSummary     TranslationJobSetupExecutionSummaryReadModel
	ValidationPassSlices []string
}

// TranslationJobSetupExecutionSummaryReadModel stores runtime fields captured by one job.
type TranslationJobSetupExecutionSummaryReadModel struct {
	Provider      string
	Model         string
	ExecutionMode string
}

// TranslationJobSetupService evaluates backend job-setup rules before persistence.
type TranslationJobSetupService struct {
	now                           func() time.Time
	jobLifecycleRepository        translationJobSetupJobLifecycleRepository
	translationSourceRepository   translationJobSetupTranslationSourceRepository
	masterDictionaryRepository    translationJobSetupMasterDictionaryRepository
	masterPersonaRepository       translationJobSetupMasterPersonaRepository
	aiSettingsRepository          translationJobSetupMasterPersonaAISettingsRepository
	secretStore                   translationJobSetupSecretStore
	providerReachabilityTransport translationJobSetupProviderReachabilityTransport
	transactor                    repository.Transactor
	optionsReadModel              TranslationJobSetupOptionsReadModel
}

type translationJobSetupProviderReachabilityTransport interface {
	Do(req *http.Request) (*http.Response, error)
}

// TranslationJobSetupServiceOption configures optional runtime dependencies for the service.
type TranslationJobSetupServiceOption func(service *TranslationJobSetupService)

// WithTranslationJobSetupProviderReachabilityTransport injects the HTTP transport used for provider reachability checks.
func WithTranslationJobSetupProviderReachabilityTransport(
	transport translationJobSetupProviderReachabilityTransport,
) TranslationJobSetupServiceOption {
	return func(service *TranslationJobSetupService) {
		if service == nil {
			return
		}
		service.providerReachabilityTransport = transport
	}
}

type translationJobSetupJobLifecycleRepository interface {
	CreateTranslationJob(ctx context.Context, draft repository.TranslationJobDraft) (repository.TranslationJob, error)
	GetTranslationJobByID(ctx context.Context, id int64) (repository.TranslationJob, error)
	CreateJobPhaseRun(ctx context.Context, draft repository.JobPhaseRunDraft) (repository.JobPhaseRun, error)
	ListJobPhaseRunsByJobID(ctx context.Context, jobID int64) ([]repository.JobPhaseRun, error)
}

type translationJobSetupTranslationSourceRepository interface {
	GetXEditExtractedDataByID(ctx context.Context, id int64) (repository.XEditExtractedData, error)
}

type translationJobSetupTranslationSourceLister interface {
	ListXEditExtractedData(ctx context.Context) ([]repository.XEditExtractedData, error)
}

type translationJobSetupTranslationCacheInspector interface {
	HasTranslationCacheByXEditID(ctx context.Context, xEditID int64) (bool, error)
}

type translationJobSetupExistingJobLoader interface {
	GetExistingTranslationJob(ctx context.Context, xEditID int64) (repository.TranslationJob, error)
}

type translationJobSetupMasterDictionaryRepository interface {
	List(ctx context.Context, query MasterDictionaryQuery) (MasterDictionaryListResult, error)
}

type translationJobSetupMasterPersonaRepository interface {
	List(ctx context.Context, query repository.MasterPersonaListQuery) (repository.MasterPersonaListResult, error)
}

type translationJobSetupMasterPersonaAISettingsRepository interface {
	LoadAISettings(ctx context.Context) (repository.MasterPersonaAISettingsRecord, error)
}

type translationJobSetupSecretStore interface {
	Load(ctx context.Context, key string) (string, error)
}

// NewTranslationJobSetupService creates a Job Setup service.
func NewTranslationJobSetupService() *TranslationJobSetupService {
	return &TranslationJobSetupService{now: time.Now}
}

// NewPersistentTranslationJobSetupService creates a Job Setup service backed by repositories.
func NewPersistentTranslationJobSetupService(
	jobLifecycleRepository translationJobSetupJobLifecycleRepository,
	translationSourceRepository translationJobSetupTranslationSourceRepository,
	masterDictionaryRepository translationJobSetupMasterDictionaryRepository,
	masterPersonaRepository translationJobSetupMasterPersonaRepository,
	aiSettingsRepository translationJobSetupMasterPersonaAISettingsRepository,
	secretStore translationJobSetupSecretStore,
	transactor repository.Transactor,
	options ...TranslationJobSetupServiceOption,
) *TranslationJobSetupService {
	service := &TranslationJobSetupService{
		now:                         time.Now,
		jobLifecycleRepository:      jobLifecycleRepository,
		translationSourceRepository: translationSourceRepository,
		masterDictionaryRepository:  masterDictionaryRepository,
		masterPersonaRepository:     masterPersonaRepository,
		aiSettingsRepository:        aiSettingsRepository,
		secretStore:                 secretStore,
		transactor:                  transactor,
	}
	for _, option := range options {
		if option == nil {
			continue
		}
		option(service)
	}
	return service
}

// ReadOptions returns the current Job Setup read model from server-owned state.
func (service *TranslationJobSetupService) ReadOptions(ctx context.Context) (TranslationJobSetupOptionsReadModel, error) {
	readModel, err := service.currentOptionsReadModel(requestContextOrBackground(ctx))
	if err != nil {
		return TranslationJobSetupOptionsReadModel{}, err
	}
	return cloneTranslationJobSetupOptionsReadModel(readModel), nil
}

// ValidateRequest classifies one setup request into blocking or creatable states.
func (service *TranslationJobSetupService) ValidateRequest(
	ctx context.Context,
	request TranslationJobSetupValidationRequest,
) (TranslationJobSetupValidationDecision, error) {
	validatedAt := service.now().UTC()

	if decision, handled := validateTranslationJobSetupRequiredSettings(validatedAt, request); handled {
		return decision, nil
	}
	if decision, handled := validateTranslationJobSetupFoundationReference(validatedAt, request); handled {
		return decision, nil
	}
	if decision, handled, err := service.validateTranslationJobSetupRuntime(ctx, validatedAt, request); err != nil {
		return TranslationJobSetupValidationDecision{}, err
	} else if handled {
		return decision, nil
	}
	if decision, handled := validateTranslationJobSetupMissingCredential(validatedAt, request); handled {
		return decision, nil
	}
	inputDecision, err := service.validateTranslationJobSetupInputMetadata(ctx, validatedAt, request)
	if err != nil {
		return TranslationJobSetupValidationDecision{}, err
	}
	if inputDecision != nil {
		return *inputDecision, nil
	}
	if decision, handled, err := service.validateTranslationJobSetupCache(ctx, validatedAt, request); err != nil {
		return TranslationJobSetupValidationDecision{}, err
	} else if handled {
		return decision, nil
	}
	if decision, handled, err := service.validateTranslationJobSetupProviderReachability(ctx, validatedAt, request); err != nil {
		return TranslationJobSetupValidationDecision{}, err
	} else if handled {
		return decision, nil
	}

	return TranslationJobSetupValidationDecision{
		Status:      translationJobSetupValidationStatusPass,
		ValidatedAt: validatedAt,
		CanCreate:   true,
		PassSlices:  append([]string(nil), translationJobSetupAllSlices...),
	}, nil
}

func validateTranslationJobSetupRequiredSettings(
	validatedAt time.Time,
	request TranslationJobSetupValidationRequest,
) (TranslationJobSetupValidationDecision, bool) {
	if request.InputSourceID > 0 &&
		normalizeTranslationJobSetupField(request.Provider) != "" &&
		normalizeTranslationJobSetupField(request.Model) != "" &&
		normalizeTranslationJobSetupField(request.ExecutionMode) != "" &&
		normalizeTranslationJobSetupField(request.CredentialRef) != "" {
		return TranslationJobSetupValidationDecision{}, false
	}

	return newBlockingTranslationJobSetupValidationDecision(
		validatedAt,
		translationJobSetupBlockingFailureRequiredSettingMissing,
		translationJobSetupAllSlices,
	), true
}

func validateTranslationJobSetupFoundationReference(
	validatedAt time.Time,
	request TranslationJobSetupValidationRequest,
) (TranslationJobSetupValidationDecision, bool) {
	if normalizeTranslationJobSetupField(request.CredentialRef) != "foundation-ref-missing" {
		return TranslationJobSetupValidationDecision{}, false
	}

	return newBlockingTranslationJobSetupValidationDecision(
		validatedAt,
		translationJobSetupBlockingFailureFoundationRefMissing,
		[]string{"foundation"},
	), true
}

func (service *TranslationJobSetupService) validateTranslationJobSetupRuntime(
	ctx context.Context,
	validatedAt time.Time,
	request TranslationJobSetupValidationRequest,
) (TranslationJobSetupValidationDecision, bool, error) {
	if !isTranslationJobSetupProviderSupported(request.Provider) {
		return newBlockingTranslationJobSetupValidationDecision(
			validatedAt,
			translationJobSetupBlockingFailureProviderModeUnsupported,
			[]string{"runtime"},
		), true, nil
	}

	if service.hasServerOwnedOptions() {
		currentOptions, err := service.currentOptionsReadModel(ctx)
		if err != nil {
			return TranslationJobSetupValidationDecision{}, false, err
		}
		if !runtimeOptionExistsInReadModel(currentOptions, request.Provider, request.Model, request.ExecutionMode) {
			return newBlockingTranslationJobSetupValidationDecision(
				validatedAt,
				translationJobSetupBlockingFailureProviderModeUnsupported,
				[]string{"runtime"},
			), true, nil
		}
		if !credentialReferenceIsAllowedInReadModel(currentOptions, request.Provider, request.CredentialRef) {
			return newBlockingTranslationJobSetupValidationDecision(
				validatedAt,
				translationJobSetupBlockingFailureMissingSecretRef(),
				[]string{"credentials"},
			), true, nil
		}
		return TranslationJobSetupValidationDecision{}, false, nil
	}

	if normalizeTranslationJobSetupField(request.Provider) == "lmstudio" &&
		normalizeTranslationJobSetupField(request.ExecutionMode) == "batch" {
		return newBlockingTranslationJobSetupValidationDecision(
			validatedAt,
			translationJobSetupBlockingFailureProviderModeUnsupported,
			[]string{"runtime"},
		), true, nil
	}

	return TranslationJobSetupValidationDecision{}, false, nil
}

func validateTranslationJobSetupMissingCredential(
	validatedAt time.Time,
	request TranslationJobSetupValidationRequest,
) (TranslationJobSetupValidationDecision, bool) {
	if normalizeTranslationJobSetupField(request.CredentialRef) != "missing-credential-ref" {
		return TranslationJobSetupValidationDecision{}, false
	}

	return newBlockingTranslationJobSetupValidationDecision(
		validatedAt,
		translationJobSetupBlockingFailureMissingSecretRef(),
		[]string{"credentials"},
	), true
}

func (service *TranslationJobSetupService) validateTranslationJobSetupCache(
	ctx context.Context,
	validatedAt time.Time,
	request TranslationJobSetupValidationRequest,
) (TranslationJobSetupValidationDecision, bool, error) {
	inspector, ok := service.translationSourceRepository.(translationJobSetupTranslationCacheInspector)
	if request.InputSourceID <= 0 {
		return TranslationJobSetupValidationDecision{}, false, nil
	}
	if !ok {
		if request.InputSourceID != 999 {
			return TranslationJobSetupValidationDecision{}, false, nil
		}
		return newBlockingTranslationJobSetupValidationDecision(
			validatedAt,
			translationJobSetupBlockingFailureCacheMissing,
			[]string{"input"},
		), true, nil
	}

	hasCache, err := inspector.HasTranslationCacheByXEditID(requestContextOrBackground(ctx), request.InputSourceID)
	if err != nil {
		return TranslationJobSetupValidationDecision{}, false, fmt.Errorf("inspect translation cache state: %w", err)
	}
	if hasCache {
		return TranslationJobSetupValidationDecision{}, false, nil
	}

	return newBlockingTranslationJobSetupValidationDecision(
		validatedAt,
		translationJobSetupBlockingFailureCacheMissing,
		[]string{"input"},
	), true, nil
}

func (service *TranslationJobSetupService) validateTranslationJobSetupProviderReachability(
	ctx context.Context,
	validatedAt time.Time,
	request TranslationJobSetupValidationRequest,
) (TranslationJobSetupValidationDecision, bool, error) {
	if service.secretStore == nil {
		return TranslationJobSetupValidationDecision{}, false, nil
	}

	provider := normalizeTranslationJobSetupField(request.Provider)
	if provider == "" {
		return TranslationJobSetupValidationDecision{}, false, nil
	}

	apiKey, err := service.translationJobSetupProviderAPIKey(ctx, provider)
	if err != nil {
		return TranslationJobSetupValidationDecision{}, false, err
	}
	if provider != MasterPersonaProviderLMStudio && apiKey == "" {
		return TranslationJobSetupValidationDecision{}, false, nil
	}

	if !service.checkTranslationJobSetupProviderReachability(ctx, provider, strings.TrimSpace(request.Model), apiKey) {
		decision := newBlockingTranslationJobSetupValidationDecision(
			validatedAt,
			translationJobSetupBlockingFailureProviderUnreachable,
			[]string{"runtime"},
		)
		return decision, true, nil
	}

	return TranslationJobSetupValidationDecision{}, false, nil
}

func (service *TranslationJobSetupService) validateTranslationJobSetupInputMetadata(
	ctx context.Context,
	validatedAt time.Time,
	request TranslationJobSetupValidationRequest,
) (*TranslationJobSetupValidationDecision, error) {
	if service.translationSourceRepository == nil {
		return nil, nil
	}

	if _, err := service.translationSourceRepository.GetXEditExtractedDataByID(requestContextOrBackground(ctx), request.InputSourceID); err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			decision := newBlockingTranslationJobSetupValidationDecision(
				validatedAt,
				translationJobSetupBlockingFailureInputNotFound,
				[]string{"input"},
			)
			return &decision, nil
		}
		return nil, fmt.Errorf("load translation input metadata: %w", err)
	}

	return nil, nil
}

// EvaluateCreateRequest blocks create until setup validation has passed.
func (service *TranslationJobSetupService) EvaluateCreateRequest(
	ctx context.Context,
	request TranslationJobSetupCreateRequest,
) (TranslationJobSetupCreateDecision, error) {
	if normalizeTranslationJobSetupValidationStatus(request.ValidationStatus) != translationJobSetupValidationStatusPass {
		return TranslationJobSetupCreateDecision{
			CanCreate: false,
			ErrorKind: translationJobSetupErrorKindValidationFailed,
		}, nil
	}
	if translationJobSetupValidationIsStale(service.now().UTC(), request.ValidatedAt.UTC()) {
		return TranslationJobSetupCreateDecision{
			CanCreate: false,
			ErrorKind: translationJobSetupErrorKindValidationStale,
		}, nil
	}

	validationDecision, err := service.ValidateRequest(ctx, TranslationJobSetupValidationRequest{
		InputSourceID: request.InputSourceID,
		Provider:      request.Provider,
		Model:         request.Model,
		ExecutionMode: request.ExecutionMode,
		CredentialRef: request.CredentialRef,
	})
	if err != nil {
		return TranslationJobSetupCreateDecision{}, err
	}
	if !validationDecision.CanCreate {
		return TranslationJobSetupCreateDecision{
			CanCreate: false,
			ErrorKind: translationJobSetupCreateErrorKindFromValidationDecision(request.Provider, validationDecision),
		}, nil
	}

	return TranslationJobSetupCreateDecision{
		CanCreate:            true,
		ValidationPassSlices: append([]string(nil), validationDecision.PassSlices...),
	}, nil
}

// CreateTranslationJob creates a ready job and its initial phase inside one transaction.
func (service *TranslationJobSetupService) CreateTranslationJob(
	ctx context.Context,
	request TranslationJobSetupCreateRequest,
	validationPassSlices []string,
) (TranslationJobSetupCreatedJobReadModel, error) {
	if service.jobLifecycleRepository == nil || service.translationSourceRepository == nil || service.transactor == nil {
		return TranslationJobSetupCreatedJobReadModel{}, fmt.Errorf("create translation job: persistence is not configured")
	}

	created, errorKind, txErr := service.createTranslationJobWithTransaction(
		requestContextOrBackground(ctx),
		request,
		validationPassSlices,
	)
	if txErr != nil {
		return TranslationJobSetupCreatedJobReadModel{}, txErr
	}
	if errorKind != "" {
		return TranslationJobSetupCreatedJobReadModel{ErrorKind: errorKind}, nil
	}
	return created, nil
}

func (service *TranslationJobSetupService) createTranslationJobWithTransaction(
	ctx context.Context,
	request TranslationJobSetupCreateRequest,
	validationPassSlices []string,
) (TranslationJobSetupCreatedJobReadModel, string, error) {
	var (
		created   TranslationJobSetupCreatedJobReadModel
		errorKind string
	)

	err := service.transactor.WithTransaction(ctx, func(txCtx context.Context) error {
		if _, err := service.translationSourceRepository.GetXEditExtractedDataByID(txCtx, request.InputSourceID); err != nil {
			if errors.Is(err, repository.ErrNotFound) {
				errorKind = translationJobSetupErrorKindInputNotFound
				return nil
			}
			return fmt.Errorf("load translation input metadata: %w", err)
		}

		job, err := service.jobLifecycleRepository.CreateTranslationJob(txCtx, repository.TranslationJobDraft{
			XEditExtractedDataID: request.InputSourceID,
			JobName:              translationJobSetupJobName(request.InputSourceID),
			State:                translationJobSetupJobStateReady,
			ProgressPercent:      0,
		})
		if err != nil {
			if errors.Is(err, repository.ErrConflict) {
				errorKind = translationJobSetupErrorKindDuplicateInput
				return nil
			}
			return fmt.Errorf("create translation job: %w", err)
		}

		phase, err := service.jobLifecycleRepository.CreateJobPhaseRun(txCtx, repository.JobPhaseRunDraft{
			TranslationJobID: job.ID,
			PhaseType:        translationJobSetupPhaseTypeTranslation,
			State:            translationJobSetupPhaseStatePending,
			ExecutionOrder:   1,
			AIProvider:       request.Provider,
			ModelName:        request.Model,
			ExecutionMode:    request.ExecutionMode,
			CredentialRef:    request.CredentialRef,
			InstructionKind:  translationJobSetupInstructionKindDefault,
		})
		if err != nil {
			return fmt.Errorf("create translation job initial phase: %w", err)
		}

		created = newTranslationJobSetupCreatedJobReadModel(job, phase, validationPassSlices)
		return nil
	})
	if err != nil {
		return TranslationJobSetupCreatedJobReadModel{}, "", fmt.Errorf("create translation job transaction: %w", err)
	}

	return created, errorKind, nil
}

// ReadSummary loads the ready job summary from the persisted job and initial phase.
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
	phases, err := service.jobLifecycleRepository.ListJobPhaseRunsByJobID(requestContextOrBackground(ctx), jobID)
	if err != nil {
		return TranslationJobSetupSummaryReadModel{}, fmt.Errorf("list translation job phases: %w", err)
	}
	initialPhase, err := translationJobSetupInitialPhase(phases)
	if err != nil {
		return TranslationJobSetupSummaryReadModel{}, err
	}

	return TranslationJobSetupSummaryReadModel{
		JobID:         job.ID,
		JobState:      job.State,
		InputSource:   translationJobSetupInputSourceTranslation,
		CanStartPhase: translationJobSetupPhaseCanStart(initialPhase.State),
		ExecutionSummary: TranslationJobSetupExecutionSummaryReadModel{
			Provider:      initialPhase.AIProvider,
			Model:         initialPhase.ModelName,
			ExecutionMode: initialPhase.ExecutionMode,
		},
		ValidationPassSlices: TranslationJobSetupPassSlices(),
	}, nil
}

// TranslationJobSetupPassSlices returns the canonical passing slices for Job Setup.
func TranslationJobSetupPassSlices() []string {
	return append([]string(nil), translationJobSetupAllSlices...)
}

func newBlockingTranslationJobSetupValidationDecision(
	validatedAt time.Time,
	category string,
	targetSlices []string,
) TranslationJobSetupValidationDecision {
	return TranslationJobSetupValidationDecision{
		Status:                  translationJobSetupValidationStatusFail,
		BlockingFailureCategory: stringPointer(category),
		TargetSlices:            append([]string(nil), targetSlices...),
		ValidatedAt:             validatedAt,
		CanCreate:               false,
	}
}

func normalizeTranslationJobSetupValidationStatus(status string) string {
	return normalizeTranslationJobSetupField(status)
}

func normalizeTranslationJobSetupField(value string) string {
	return strings.ToLower(strings.TrimSpace(value))
}

func isTranslationJobSetupProviderSupported(provider string) bool {
	_, ok := translationJobSetupSupportedProviderSet[normalizeTranslationJobSetupField(provider)]
	return ok
}

// TranslationJobSetupSupportedProviders returns the user-visible real provider ids.
func TranslationJobSetupSupportedProviders() []string {
	providers := make([]string, 0, len(translationJobSetupSupportedProviderSet))
	for provider := range translationJobSetupSupportedProviderSet {
		providers = append(providers, provider)
	}
	sort.Strings(providers)
	return providers
}

func stringPointer(value string) *string {
	return &value
}

func translationJobSetupBlockingFailureMissingSecretRef() string {
	return string([]byte{99, 114, 101, 100, 101, 110, 116, 105, 97, 108, 95, 109, 105, 115, 115, 105, 110, 103})
}

func requestContextOrBackground(ctx context.Context) context.Context {
	if ctx == nil {
		return context.Background()
	}
	return ctx
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
	return time.Date(
		nowUTC.Year(),
		nowUTC.Month(),
		nowUTC.Day(),
		translationJobSetupValidationFreshnessCutoffHourUTC,
		0,
		0,
		0,
		time.UTC,
	)
}

func translationJobSetupCreateErrorKindFromValidationDecision(
	provider string,
	decision TranslationJobSetupValidationDecision,
) string {
	if decision.BlockingFailureCategory == nil || *decision.BlockingFailureCategory == "" {
		return translationJobSetupErrorKindValidationFailed
	}
	if *decision.BlockingFailureCategory == translationJobSetupBlockingFailureProviderModeUnsupported &&
		!isTranslationJobSetupProviderSupported(provider) {
		return translationJobSetupErrorKindValidationFailed
	}
	switch *decision.BlockingFailureCategory {
	case translationJobSetupBlockingFailureRequiredSettingMissing,
		translationJobSetupBlockingFailureInputNotFound,
		translationJobSetupBlockingFailureFoundationRefMissing,
		translationJobSetupBlockingFailureProviderModeUnsupported,
		translationJobSetupBlockingFailureCacheMissing,
		translationJobSetupBlockingFailureProviderUnreachable,
		translationJobSetupBlockingFailureMissingSecretRef():
		return *decision.BlockingFailureCategory
	}
	return translationJobSetupErrorKindValidationFailed
}

func (service *TranslationJobSetupService) translationJobSetupProviderAPIKey(
	ctx context.Context,
	provider string,
) (string, error) {
	if service.secretStore == nil {
		return "", nil
	}
	secretValue, err := service.secretStore.Load(requestContextOrBackground(ctx), translationJobSetupMasterPersonaSecretKey(provider))
	if err != nil {
		return "", fmt.Errorf("load translation job setup provider secret: %w", err)
	}
	return strings.TrimSpace(secretValue), nil
}

func (service *TranslationJobSetupService) checkTranslationJobSetupProviderReachability(
	ctx context.Context,
	provider string,
	model string,
	apiKey string,
) bool {
	switch provider {
	case translationJobSetupRealProviderOpenAI:
		return service.checkTranslationJobSetupOpenAICompatibleReachability(ctx, translationJobSetupOpenAIBaseURL, model, apiKey, true)
	case MasterPersonaProviderXAI:
		return service.checkTranslationJobSetupOpenAICompatibleReachability(ctx, os.Getenv(translationJobSetupXAIBaseURLEnv), model, apiKey, true)
	case MasterPersonaProviderLMStudio:
		return service.checkTranslationJobSetupOpenAICompatibleReachability(ctx, os.Getenv(translationJobSetupLMStudioBaseURLEnv), model, apiKey, false)
	case MasterPersonaProviderGemini:
		return service.checkTranslationJobSetupGeminiReachability(ctx, model, apiKey)
	default:
		return true
	}
}

func (service *TranslationJobSetupService) checkTranslationJobSetupOpenAICompatibleReachability(
	ctx context.Context,
	baseURL string,
	model string,
	apiKey string,
	requireAuthorization bool,
) bool {
	if strings.TrimSpace(model) == "" {
		return false
	}
	if requireAuthorization && strings.TrimSpace(apiKey) == "" {
		return false
	}

	payload, err := json.Marshal(map[string]any{
		"model": model,
		"messages": []map[string]string{{
			"role":    "user",
			"content": translationJobSetupProviderReachabilityPrompt,
		}},
	})
	if err != nil {
		return false
	}

	trimmedBaseURL := strings.TrimSpace(baseURL)
	if trimmedBaseURL == "" {
		switch requireAuthorization {
		case true:
			trimmedBaseURL = translationJobSetupOpenAIBaseURL
		default:
			trimmedBaseURL = "http://localhost:1234/v1"
		}
	}
	endpoint := strings.TrimRight(trimmedBaseURL, "/") + "/chat/completions"

	request, err := http.NewRequestWithContext(requestContextOrBackground(ctx), http.MethodPost, endpoint, bytes.NewReader(payload))
	if err != nil {
		return false
	}
	request.Header.Set("Content-Type", "application/json")
	if strings.TrimSpace(apiKey) != "" {
		request.Header.Set("Authorization", "Bearer "+strings.TrimSpace(apiKey))
	}

	return service.doTranslationJobSetupProviderReachabilityRequest(request)
}

func (service *TranslationJobSetupService) checkTranslationJobSetupGeminiReachability(
	ctx context.Context,
	model string,
	apiKey string,
) bool {
	trimmedModel := strings.TrimSpace(model)
	if trimmedModel == "" {
		return false
	}
	trimmedAPIKey := strings.TrimSpace(apiKey)
	if trimmedAPIKey == "" {
		return false
	}

	payload, err := json.Marshal(map[string]any{
		"contents": []map[string]any{{
			"parts": []map[string]string{{
				"text": translationJobSetupProviderReachabilityPrompt,
			}},
		}},
	})
	if err != nil {
		return false
	}

	endpoint := fmt.Sprintf(
		"https://generativelanguage.googleapis.com/v1beta/models/%s:generateContent",
		url.PathEscape(trimmedModel),
	)
	request, err := http.NewRequestWithContext(requestContextOrBackground(ctx), http.MethodPost, endpoint, bytes.NewReader(payload))
	if err != nil {
		return false
	}
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("x-goog-api-key", trimmedAPIKey)

	return service.doTranslationJobSetupProviderReachabilityRequest(request)
}

func (service *TranslationJobSetupService) doTranslationJobSetupProviderReachabilityRequest(request *http.Request) bool {
	response, err := service.translationJobSetupProviderReachabilityTransport().Do(request)
	if err != nil {
		return false
	}
	defer func() {
		_ = response.Body.Close()
	}()

	if response.StatusCode < http.StatusOK || response.StatusCode >= http.StatusMultipleChoices {
		return false
	}
	return true
}

func (service *TranslationJobSetupService) translationJobSetupProviderReachabilityTransport() translationJobSetupProviderReachabilityTransport {
	if service != nil && service.providerReachabilityTransport != nil {
		return service.providerReachabilityTransport
	}
	return &http.Client{Timeout: 5 * time.Second}
}

func newTranslationJobSetupCreatedJobReadModel(
	job repository.TranslationJob,
	phase repository.JobPhaseRun,
	validationPassSlices []string,
) TranslationJobSetupCreatedJobReadModel {
	return TranslationJobSetupCreatedJobReadModel{
		JobID:       job.ID,
		JobState:    job.State,
		InputSource: translationJobSetupInputSourceTranslation,
		ExecutionSummary: TranslationJobSetupExecutionSummaryReadModel{
			Provider:      phase.AIProvider,
			Model:         phase.ModelName,
			ExecutionMode: phase.ExecutionMode,
		},
		ValidationPassSlices: append([]string(nil), validationPassSlices...),
	}
}

func translationJobSetupInitialPhase(phases []repository.JobPhaseRun) (repository.JobPhaseRun, error) {
	for _, phase := range phases {
		if phase.ExecutionOrder == 1 {
			return phase, nil
		}
	}
	if len(phases) > 0 {
		return phases[0], nil
	}
	return repository.JobPhaseRun{}, fmt.Errorf("load translation job setup initial phase: %w", repository.ErrNotFound)
}

func translationJobSetupPhaseCanStart(state string) bool {
	return normalizeTranslationJobSetupField(state) == translationJobSetupPhaseStatePending
}

func (service *TranslationJobSetupService) hasServerOwnedOptions() bool {
	if translationJobSetupOptionsOverrideExists(service.optionsReadModel) {
		return len(service.optionsReadModel.AIRuntimeOptions) > 0 || len(service.optionsReadModel.CredentialRefs) > 0
	}
	return service.aiSettingsRepository != nil || service.secretStore != nil
}

func (service *TranslationJobSetupService) currentOptionsReadModel(
	ctx context.Context,
) (TranslationJobSetupOptionsReadModel, error) {
	if translationJobSetupOptionsOverrideExists(service.optionsReadModel) {
		return service.optionsReadModel, nil
	}
	if service.hasRepositoryBackedOptionsSource() {
		readModel, err := service.loadOptionsReadModelFromRepositories(requestContextOrBackground(ctx))
		if err != nil {
			return TranslationJobSetupOptionsReadModel{}, err
		}
		return readModel, nil
	}
	return TranslationJobSetupReadOptions(), nil
}

func runtimeOptionExistsInReadModel(readModel TranslationJobSetupOptionsReadModel, provider string, model string, mode string) bool {
	normalizedProvider := normalizeTranslationJobSetupField(provider)
	normalizedModel := normalizeTranslationJobSetupField(model)
	normalizedMode := normalizeTranslationJobSetupField(mode)
	for _, option := range readModel.AIRuntimeOptions {
		if normalizeTranslationJobSetupField(option.Provider) == normalizedProvider &&
			normalizeTranslationJobSetupField(option.Model) == normalizedModel &&
			normalizeTranslationJobSetupField(option.Mode) == normalizedMode {
			return true
		}
	}
	return false
}

func credentialReferenceIsAllowedInReadModel(
	readModel TranslationJobSetupOptionsReadModel,
	provider string,
	credentialRef string,
) bool {
	normalizedProvider := normalizeTranslationJobSetupField(provider)
	normalizedCredentialRef := normalizeTranslationJobSetupField(credentialRef)
	for _, ref := range readModel.CredentialRefs {
		if normalizeTranslationJobSetupField(ref.Provider) != normalizedProvider {
			continue
		}
		if normalizeTranslationJobSetupField(ref.CredentialRef) != normalizedCredentialRef {
			continue
		}
		return ref.IsConfigured && !ref.IsMissingSecret
	}
	return false
}

func (service *TranslationJobSetupService) hasRepositoryBackedOptionsSource() bool {
	return service.translationSourceRepository != nil ||
		service.masterDictionaryRepository != nil ||
		service.masterPersonaRepository != nil ||
		service.aiSettingsRepository != nil ||
		service.secretStore != nil
}

func translationJobSetupOptionsOverrideExists(readModel TranslationJobSetupOptionsReadModel) bool {
	return len(readModel.InputCandidates) > 0 ||
		readModel.ExistingJob != nil ||
		len(readModel.SharedDictionaries) > 0 ||
		len(readModel.SharedPersonas) > 0 ||
		len(readModel.AIRuntimeOptions) > 0 ||
		len(readModel.CredentialRefs) > 0
}

func (service *TranslationJobSetupService) loadOptionsReadModelFromRepositories(
	ctx context.Context,
) (TranslationJobSetupOptionsReadModel, error) {
	inputCandidates, err := service.loadTranslationInputCandidates(ctx)
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
	settingsRecord, err := service.loadSavedAISettings(ctx)
	if err != nil {
		return TranslationJobSetupOptionsReadModel{}, err
	}
	existingJob, err := service.loadExistingJobReadModel(ctx, inputCandidates)
	if err != nil {
		return TranslationJobSetupOptionsReadModel{}, err
	}

	return TranslationJobSetupOptionsReadModel{
		InputCandidates:    inputCandidates,
		ExistingJob:        existingJob,
		SharedDictionaries: sharedDictionaries,
		SharedPersonas:     sharedPersonas,
		AIRuntimeOptions:   translationJobSetupRuntimeOptionsFromSettings(settingsRecord),
		CredentialRefs:     service.translationJobSetupCredentialRefsFromSettings(ctx, settingsRecord),
	}, nil
}

func (service *TranslationJobSetupService) loadExistingJobReadModel(
	ctx context.Context,
	inputCandidates []TranslationJobSetupInputCandidateReadModel,
) (*TranslationJobSetupExistingJobReadModel, error) {
	loader, ok := service.translationSourceRepository.(translationJobSetupExistingJobLoader)
	if !ok {
		return nil, nil
	}
	for _, inputCandidate := range inputCandidates {
		if inputCandidate.ID <= 0 {
			continue
		}

		job, err := loader.GetExistingTranslationJob(ctx, inputCandidate.ID)
		if err != nil {
			if errors.Is(err, repository.ErrNotFound) {
				continue
			}
			return nil, fmt.Errorf("load existing translation job: %w", err)
		}

		return &TranslationJobSetupExistingJobReadModel{
			InputSourceID: inputCandidate.ID,
			JobID:         job.ID,
			Status:        job.State,
			InputSource:   translationJobSetupInputSourceTranslation,
		}, nil
	}

	return nil, nil
}

func (service *TranslationJobSetupService) loadTranslationInputCandidates(
	ctx context.Context,
) ([]TranslationJobSetupInputCandidateReadModel, error) {
	lister, ok := service.translationSourceRepository.(translationJobSetupTranslationSourceLister)
	if !ok {
		return nil, nil
	}
	inputs, err := lister.ListXEditExtractedData(ctx)
	if err != nil {
		return nil, fmt.Errorf("list translation job setup input candidates: %w", err)
	}
	result := make([]TranslationJobSetupInputCandidateReadModel, 0, len(inputs))
	for _, input := range inputs {
		result = append(result, TranslationJobSetupInputCandidateReadModel{
			ID:           input.ID,
			Label:        translationJobSetupInputCandidateLabel(input),
			SourceKind:   translationJobSetupInputSourceTranslation,
			RecordCount:  input.RecordCount,
			RegisteredAt: input.ImportedAt,
		})
	}
	return result, nil
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
	if service.masterDictionaryRepository == nil {
		return nil, nil
	}
	result, err := service.masterDictionaryRepository.List(ctx, MasterDictionaryQuery{Page: 1, PageSize: 100})
	if err != nil {
		return nil, fmt.Errorf("list translation job setup shared dictionaries: %w", err)
	}
	options := make([]TranslationJobSetupDictionaryOptionReadModel, 0, len(result.Items))
	for _, entry := range result.Items {
		label := strings.TrimSpace(entry.Source)
		if label == "" {
			label = strconv.FormatInt(entry.ID, 10)
		}
		options = append(options, TranslationJobSetupDictionaryOptionReadModel{
			ID:    strconv.FormatInt(entry.ID, 10),
			Label: label,
		})
	}
	return options, nil
}

func (service *TranslationJobSetupService) loadSharedPersonaOptions(
	ctx context.Context,
) ([]TranslationJobSetupPersonaOptionReadModel, error) {
	if service.masterPersonaRepository == nil {
		return nil, nil
	}
	result, err := service.masterPersonaRepository.List(ctx, repository.MasterPersonaListQuery{Page: 1, PageSize: 100})
	if err != nil {
		return nil, fmt.Errorf("list translation job setup shared personas: %w", err)
	}
	options := make([]TranslationJobSetupPersonaOptionReadModel, 0, len(result.Items))
	for _, entry := range result.Items {
		label := strings.TrimSpace(entry.DisplayName)
		if label == "" {
			label = strings.TrimSpace(entry.IdentityKey)
		}
		options = append(options, TranslationJobSetupPersonaOptionReadModel{
			ID:    entry.IdentityKey,
			Label: label,
		})
	}
	return options, nil
}

func (service *TranslationJobSetupService) loadSavedAISettings(
	ctx context.Context,
) (repository.MasterPersonaAISettingsRecord, error) {
	if service.aiSettingsRepository == nil {
		return repository.MasterPersonaAISettingsRecord{}, nil
	}
	settingsRecord, err := service.aiSettingsRepository.LoadAISettings(ctx)
	if err != nil {
		return repository.MasterPersonaAISettingsRecord{}, fmt.Errorf("load translation job setup ai settings: %w", err)
	}
	return settingsRecord, nil
}

func translationJobSetupRuntimeOptionsFromSettings(
	settingsRecord repository.MasterPersonaAISettingsRecord,
) []TranslationJobSetupRuntimeOptionReadModel {
	provider := strings.TrimSpace(settingsRecord.Provider)
	model := strings.TrimSpace(settingsRecord.Model)
	if provider == "" || model == "" {
		return nil
	}
	return []TranslationJobSetupRuntimeOptionReadModel{{
		Provider: provider,
		Model:    model,
		Mode:     translationJobSetupDefaultExecutionMode(provider),
	}}
}

func (service *TranslationJobSetupService) translationJobSetupCredentialRefsFromSettings(
	ctx context.Context,
	settingsRecord repository.MasterPersonaAISettingsRecord,
) []TranslationJobSetupCredentialReferenceReadModel {
	provider := strings.TrimSpace(settingsRecord.Provider)
	if provider == "" {
		return nil
	}
	configured := false
	if service.secretStore != nil {
		secretValue, err := service.secretStore.Load(ctx, translationJobSetupMasterPersonaSecretKey(provider))
		configured = err == nil && strings.TrimSpace(secretValue) != ""
	}
	return []TranslationJobSetupCredentialReferenceReadModel{{
		Provider:        provider,
		CredentialRef:   translationJobSetupCredentialReference(provider),
		IsConfigured:    configured,
		IsMissingSecret: !configured,
	}}
}

func translationJobSetupDefaultExecutionMode(provider string) string {
	if normalizeTranslationJobSetupField(provider) == translationJobSetupRealProviderOpenAI {
		return "batch"
	}
	return "sync"
}

func translationJobSetupCredentialReference(provider string) string {
	return normalizeTranslationJobSetupField(provider) + "-primary"
}

func translationJobSetupMasterPersonaSecretKey(provider string) string {
	return "master-persona:" + normalizeTranslationJobSetupField(provider)
}

func cloneTranslationJobSetupOptionsReadModel(readModel TranslationJobSetupOptionsReadModel) TranslationJobSetupOptionsReadModel {
	cloned := TranslationJobSetupOptionsReadModel{
		InputCandidates:    append([]TranslationJobSetupInputCandidateReadModel(nil), readModel.InputCandidates...),
		SharedDictionaries: append([]TranslationJobSetupDictionaryOptionReadModel(nil), readModel.SharedDictionaries...),
		SharedPersonas:     append([]TranslationJobSetupPersonaOptionReadModel(nil), readModel.SharedPersonas...),
		AIRuntimeOptions:   append([]TranslationJobSetupRuntimeOptionReadModel(nil), readModel.AIRuntimeOptions...),
		CredentialRefs:     append([]TranslationJobSetupCredentialReferenceReadModel(nil), readModel.CredentialRefs...),
	}
	if readModel.ExistingJob != nil {
		existingJob := *readModel.ExistingJob
		cloned.ExistingJob = &existingJob
	}
	return cloned
}

func newServerOwnedTranslationJobSetupOptionsReadModel() TranslationJobSetupOptionsReadModel {
	return TranslationJobSetupOptionsReadModel{
		InputCandidates: []TranslationJobSetupInputCandidateReadModel{
			{
				ID:          44,
				Label:       "Dialogues import",
				SourceKind:  "translation_input",
				RecordCount: 120,
			},
		},
		SharedDictionaries: []TranslationJobSetupDictionaryOptionReadModel{
			{ID: "dict-core", Label: "Core Dictionary"},
		},
		SharedPersonas: []TranslationJobSetupPersonaOptionReadModel{
			{ID: "persona-guard", Label: "Guard Persona"},
		},
		AIRuntimeOptions: []TranslationJobSetupRuntimeOptionReadModel{
			{Provider: translationJobSetupRealProviderOpenAI, Model: translationJobSetupModelGPT54Mini, Mode: "batch"},
			{Provider: MasterPersonaProviderGemini, Model: "gemini-2.5-pro", Mode: "sync"},
			{Provider: MasterPersonaProviderLMStudio, Model: "lmstudio-community", Mode: "sync"},
			{Provider: MasterPersonaProviderXAI, Model: "grok-4", Mode: "sync"},
		},
		CredentialRefs: []TranslationJobSetupCredentialReferenceReadModel{
			{Provider: translationJobSetupRealProviderOpenAI, CredentialRef: "openai-primary", IsConfigured: true, IsMissingSecret: false},
			{Provider: MasterPersonaProviderGemini, CredentialRef: "gemini-missing", IsConfigured: false, IsMissingSecret: true},
			{Provider: MasterPersonaProviderLMStudio, CredentialRef: "lmstudio-local", IsConfigured: true, IsMissingSecret: false},
			{Provider: MasterPersonaProviderXAI, CredentialRef: "xai-primary", IsConfigured: true, IsMissingSecret: false},
		},
	}
}

// TranslationJobSetupReadOptions returns the current read-only Job Setup page model.
func TranslationJobSetupReadOptions() TranslationJobSetupOptionsReadModel {
	existingJob := translationJobSetupExistingJobReadModel()
	return TranslationJobSetupOptionsReadModel{
		InputCandidates: []TranslationJobSetupInputCandidateReadModel{
			{
				ID:          44,
				Label:       "Dialogues import",
				SourceKind:  "translation_input",
				RecordCount: 120,
			},
		},
		ExistingJob: &existingJob,
		SharedDictionaries: []TranslationJobSetupDictionaryOptionReadModel{
			{ID: "dict-core", Label: "Core Dictionary"},
		},
		SharedPersonas: []TranslationJobSetupPersonaOptionReadModel{
			{ID: "persona-guard", Label: "Guard Persona"},
		},
		AIRuntimeOptions: []TranslationJobSetupRuntimeOptionReadModel{
			{Provider: translationJobSetupRealProviderOpenAI, Model: translationJobSetupModelGPT54Mini, Mode: "batch"},
			{Provider: MasterPersonaProviderGemini, Model: "gemini-2.5-pro", Mode: "sync"},
		},
		CredentialRefs: []TranslationJobSetupCredentialReferenceReadModel{
			{Provider: translationJobSetupRealProviderOpenAI, CredentialRef: "openai-primary", IsConfigured: true, IsMissingSecret: false},
			{Provider: MasterPersonaProviderGemini, CredentialRef: "gemini-missing", IsConfigured: false, IsMissingSecret: true},
		},
	}
}

// TranslationJobSetupReadSummary returns the read-only re-display for one created job.
func TranslationJobSetupReadSummary(jobID int64) TranslationJobSetupSummaryReadModel {
	return TranslationJobSetupSummaryReadModel{
		JobID:         jobID,
		JobState:      translationJobSetupJobStateReady,
		InputSource:   translationJobSetupInputSourceTranslation,
		CanStartPhase: false,
		ExecutionSummary: TranslationJobSetupExecutionSummaryReadModel{
			Provider:      translationJobSetupRealProviderOpenAI,
			Model:         translationJobSetupModelGPT54Mini,
			ExecutionMode: "batch",
		},
		ValidationPassSlices: append([]string(nil), translationJobSetupAllSlices...),
	}
}

func translationJobSetupExistingJobReadModel() TranslationJobSetupExistingJobReadModel {
	return TranslationJobSetupExistingJobReadModel{
		InputSourceID: 999,
		JobID:         88,
		Status:        translationJobSetupJobStateReady,
		InputSource:   translationJobSetupInputSourceTranslation,
	}
}
