package usecase

import (
	"context"
	"errors"
	"strings"
	"time"

	jobsetupservice "aitranslationenginejp/internal/service"
)

// TranslationJobSetupErrorKind identifies contract-level rejected outcomes.
type TranslationJobSetupErrorKind string

const (
	// TranslationJobSetupErrorKindPhaseRuntimeMissing identifies a rejected outcome caused by one or more missing phase runtime settings.
	TranslationJobSetupErrorKindPhaseRuntimeMissing TranslationJobSetupErrorKind = "phase_runtime_missing"
	// TranslationJobSetupErrorKindRequiredSettingMissing identifies a rejected outcome caused by missing required setup state.
	TranslationJobSetupErrorKindRequiredSettingMissing TranslationJobSetupErrorKind = "required_setting_missing"
	// TranslationJobSetupErrorKindInputNotFound identifies a rejected outcome caused by a missing input source.
	TranslationJobSetupErrorKindInputNotFound TranslationJobSetupErrorKind = "input_not_found"
	// TranslationJobSetupErrorKindCacheMissing identifies a rejected outcome caused by a missing input cache.
	TranslationJobSetupErrorKindCacheMissing TranslationJobSetupErrorKind = "cache_missing"
	// TranslationJobSetupErrorKindFoundationRefMissing identifies a rejected outcome caused by a missing foundation reference.
	TranslationJobSetupErrorKindFoundationRefMissing TranslationJobSetupErrorKind = "foundation_ref_missing"
	//nolint:gosec // credential_missing is a fixed public error-kind literal, not a secret.
	// TranslationJobSetupErrorKindCredentialMissing identifies a rejected outcome caused by a missing credential reference.
	TranslationJobSetupErrorKindCredentialMissing TranslationJobSetupErrorKind = "credential_missing"
	//nolint:gosec // fixed public error-kind literal, not a secret.
	// TranslationJobSetupErrorKindModelListCredentialMissing identifies a rejected outcome caused by a missing credential during provider model listing.
	TranslationJobSetupErrorKindModelListCredentialMissing TranslationJobSetupErrorKind = "model_list_credential_missing"
	// TranslationJobSetupErrorKindModelListFailed identifies a rejected outcome caused by provider model list retrieval failure.
	TranslationJobSetupErrorKindModelListFailed TranslationJobSetupErrorKind = "model_list_failed"
	// TranslationJobSetupErrorKindModelSelectionStale identifies a rejected outcome caused by a stale provider model selection.
	TranslationJobSetupErrorKindModelSelectionStale TranslationJobSetupErrorKind = "model_selection_stale"
	// TranslationJobSetupErrorKindProviderModeUnsupported identifies a rejected outcome caused by an unsupported provider/mode combination.
	TranslationJobSetupErrorKindProviderModeUnsupported TranslationJobSetupErrorKind = "provider_mode_unsupported"
	// TranslationJobSetupErrorKindProviderUnreachable identifies a rejected outcome caused by provider reachability failure.
	TranslationJobSetupErrorKindProviderUnreachable TranslationJobSetupErrorKind = "provider_unreachable"
	// TranslationJobSetupErrorKindDuplicateJobForInput identifies a rejected outcome caused by an existing job for the same input.
	TranslationJobSetupErrorKindDuplicateJobForInput TranslationJobSetupErrorKind = "duplicate_job_for_input"
	// TranslationJobSetupErrorKindValidationStale identifies one rejected create response caused by stale setup validation.
	TranslationJobSetupErrorKindValidationStale TranslationJobSetupErrorKind = "validation_stale"
	// TranslationJobSetupErrorKindPartialCreateFailed identifies a rejected outcome caused by a partial create failure.
	TranslationJobSetupErrorKindPartialCreateFailed TranslationJobSetupErrorKind = "partial_create_failed"
	// TranslationJobSetupErrorKindReadyRequired identifies a rejected outcome caused by create or follow-up work before setup is ready.
	TranslationJobSetupErrorKindReadyRequired TranslationJobSetupErrorKind = "ready_required"
	// TranslationJobSetupErrorKindInputDeleteBlocked identifies a rejected input-delete outcome caused by an existing job reference.
	TranslationJobSetupErrorKindInputDeleteBlocked TranslationJobSetupErrorKind = "input_delete_blocked"

	// TranslationJobSetupErrorKindValidationFailed remains as a compatibility alias during downstream alignment.
	TranslationJobSetupErrorKindValidationFailed TranslationJobSetupErrorKind = TranslationJobSetupErrorKindReadyRequired
	// TranslationJobSetupErrorKindDuplicateInput remains as a compatibility alias during downstream alignment.
	TranslationJobSetupErrorKindDuplicateInput TranslationJobSetupErrorKind = TranslationJobSetupErrorKindDuplicateJobForInput
)

// NormalizeTranslationJobSetupPublicErrorKind collapses internal compatibility aliases to the frozen public kinds.
func NormalizeTranslationJobSetupPublicErrorKind(kind TranslationJobSetupErrorKind) TranslationJobSetupErrorKind {
	trimmedKind := strings.TrimSpace(string(kind))
	switch strings.ToLower(trimmedKind) {
	case "":
		return TranslationJobSetupErrorKind("")
	case "validation_failed", string(TranslationJobSetupErrorKindReadyRequired):
		return TranslationJobSetupErrorKindReadyRequired
	case "duplicate_input", string(TranslationJobSetupErrorKindDuplicateJobForInput):
		return TranslationJobSetupErrorKindDuplicateJobForInput
	default:
		return TranslationJobSetupErrorKind(trimmedKind)
	}
}

// NormalizeTranslationJobSetupPublicErrorCategory normalizes optional public error categories.
func NormalizeTranslationJobSetupPublicErrorCategory(category *string) *string {
	if category == nil {
		return nil
	}
	normalized := string(NormalizeTranslationJobSetupPublicErrorKind(TranslationJobSetupErrorKind(*category)))
	return &normalized
}

const translationJobSetupValidationFreshnessCutoffHourUTC = 9

// TranslationJobSetupPhaseID identifies one translation stage inside Job Setup.
type TranslationJobSetupPhaseID string

const (
	// TranslationJobSetupPhaseIDWordTranslation identifies the term translation runtime selection.
	TranslationJobSetupPhaseIDWordTranslation TranslationJobSetupPhaseID = "word_translation"
	// TranslationJobSetupPhaseIDNPCPersonaGeneration identifies the persona generation runtime selection.
	TranslationJobSetupPhaseIDNPCPersonaGeneration TranslationJobSetupPhaseID = "npc_persona_generation"
	// TranslationJobSetupPhaseIDTextTranslation identifies the body translation runtime selection.
	TranslationJobSetupPhaseIDTextTranslation TranslationJobSetupPhaseID = "text_translation"
)

// TranslationJobSetupCredentialRequirement identifies whether a provider requires a credential reference.
type TranslationJobSetupCredentialRequirement string

const (
	// TranslationJobSetupCredentialRequirementRequired identifies providers that need a credential reference.
	TranslationJobSetupCredentialRequirementRequired TranslationJobSetupCredentialRequirement = "required"
	// TranslationJobSetupCredentialRequirementNotRequired identifies providers that do not need a credential reference.
	TranslationJobSetupCredentialRequirementNotRequired TranslationJobSetupCredentialRequirement = "not_required"
)

// TranslationJobSetupCredentialStatus identifies the credential state that can be exposed publicly.
type TranslationJobSetupCredentialStatus string

const (
	// TranslationJobSetupCredentialStatusConfigured identifies a resolved credential reference.
	TranslationJobSetupCredentialStatusConfigured TranslationJobSetupCredentialStatus = "configured"
	// TranslationJobSetupCredentialStatusMissing identifies a missing or unresolved credential reference.
	TranslationJobSetupCredentialStatusMissing TranslationJobSetupCredentialStatus = "missing"
	// TranslationJobSetupCredentialStatusNotRequired identifies providers that do not need credentials.
	TranslationJobSetupCredentialStatusNotRequired TranslationJobSetupCredentialStatus = "not_required"
)

// TranslationJobSetupBatchMode identifies the public batch-mode state for one phase.
type TranslationJobSetupBatchMode string

const (
	// TranslationJobSetupBatchModeDisabled identifies supported providers with batch mode turned off.
	TranslationJobSetupBatchModeDisabled TranslationJobSetupBatchMode = "disabled"
	// TranslationJobSetupBatchModeEnabled identifies supported providers with batch mode turned on.
	TranslationJobSetupBatchModeEnabled TranslationJobSetupBatchMode = "enabled"
	// TranslationJobSetupBatchModeUnsupported identifies providers that do not support batch mode.
	TranslationJobSetupBatchModeUnsupported TranslationJobSetupBatchMode = "unsupported"
)

// TranslationJobSetupProviderModelListStatus identifies one provider model list loading state.
type TranslationJobSetupProviderModelListStatus string

const (
	// TranslationJobSetupProviderModelListStatusNotUpdated identifies a phase whose model list was not requested yet.
	TranslationJobSetupProviderModelListStatusNotUpdated TranslationJobSetupProviderModelListStatus = "not_updated"
	// TranslationJobSetupProviderModelListStatusLoading identifies an in-flight model-list request.
	TranslationJobSetupProviderModelListStatusLoading TranslationJobSetupProviderModelListStatus = "loading"
	// TranslationJobSetupProviderModelListStatusSuccess identifies a successful model-list fetch.
	TranslationJobSetupProviderModelListStatusSuccess TranslationJobSetupProviderModelListStatus = "success"
	// TranslationJobSetupProviderModelListStatusFailed identifies a redacted model-list failure.
	TranslationJobSetupProviderModelListStatusFailed TranslationJobSetupProviderModelListStatus = "failed"
	//nolint:gosec // fixed public status literal, not a secret.
	// TranslationJobSetupProviderModelListStatusCredentialMissing identifies missing credentials before request dispatch.
	TranslationJobSetupProviderModelListStatusCredentialMissing TranslationJobSetupProviderModelListStatus = "credential_missing"
	//nolint:gosec // fixed public status literal, not a secret.
	// TranslationJobSetupProviderModelListStatusCredentialNotNeeded identifies providers that do not need credentials.
	TranslationJobSetupProviderModelListStatusCredentialNotNeeded TranslationJobSetupProviderModelListStatus = "credential_not_required"
)

// TranslationJobSetupValidationStatus identifies the outcome of one setup validation.
type TranslationJobSetupValidationStatus string

const (
	// TranslationJobSetupValidationStatusPass identifies a fully passing setup validation.
	TranslationJobSetupValidationStatusPass TranslationJobSetupValidationStatus = "pass"
	// TranslationJobSetupValidationStatusFail identifies a blocking setup validation failure.
	TranslationJobSetupValidationStatusFail TranslationJobSetupValidationStatus = "fail"
	// TranslationJobSetupValidationStatusWarning identifies a non-blocking setup validation result.
	TranslationJobSetupValidationStatusWarning TranslationJobSetupValidationStatus = "warning"
)

// TranslationJobSetupOptionsResult returns the read-only inputs required to start job setup.
type TranslationJobSetupOptionsResult struct {
	InputCandidates      []TranslationJobSetupInputCandidate
	ExistingJob          *TranslationJobSetupExistingJob
	SharedDictionaries   []TranslationJobSetupDictionaryOption
	SharedPersonas       []TranslationJobSetupPersonaOption
	AIRuntimeOptions     []TranslationJobSetupRuntimeOption
	CredentialRefs       []TranslationJobSetupCredentialReference
	ProviderCapabilities []TranslationJobSetupProviderCapability
	PhaseRuntimeDrafts   []TranslationJobSetupPhaseRuntimeDraft
}

// TranslationJobSetupInputCandidate is one selectable translation input source.
type TranslationJobSetupInputCandidate struct {
	ID           int64
	Label        string
	SourceKind   string
	RecordCount  int
	RegisteredAt time.Time
	ExistingJob  *TranslationJobSetupExistingJob
}

// TranslationJobSetupExistingJob summarizes one already prepared job.
type TranslationJobSetupExistingJob struct {
	InputSourceID int64
	JobID         int64
	Status        string
	InputSource   string
}

// TranslationJobSetupDictionaryOption is one shared dictionary choice.
type TranslationJobSetupDictionaryOption struct {
	ID    string
	Label string
}

// TranslationJobSetupPersonaOption is one shared persona choice.
type TranslationJobSetupPersonaOption struct {
	ID    string
	Label string
}

// TranslationJobSetupRuntimeOption is one selectable AI runtime option.
type TranslationJobSetupRuntimeOption struct {
	Provider string
	Model    string
	Mode     string
}

// TranslationJobSetupCredentialReference exposes only credential reference state.
type TranslationJobSetupCredentialReference struct {
	Provider        string
	CredentialRef   string
	IsConfigured    bool
	IsMissingSecret bool
}

// TranslationJobSetupProviderCapability describes one public provider capability.
type TranslationJobSetupProviderCapability struct {
	Provider                string
	CredentialRequirement   TranslationJobSetupCredentialRequirement
	SupportedExecutionModes []string
	SupportsBatchMode       bool
}

// TranslationJobSetupPhaseRuntimeDraft returns the current draft state for one phase.
type TranslationJobSetupPhaseRuntimeDraft struct {
	PhaseID              TranslationJobSetupPhaseID
	Provider             string
	Model                string
	CredentialRef        string
	CredentialStatus     TranslationJobSetupCredentialStatus
	ExecutionMode        string
	BatchMode            TranslationJobSetupBatchMode
	ModelListSourceToken string
}

// TranslationJobSetupRuntimeSelection is the selected runtime configuration.
type TranslationJobSetupRuntimeSelection struct {
	Provider      string
	Model         string
	ExecutionMode string
}

// TranslationJobSetupPhaseRuntimeSelection is the selected runtime configuration for one phase.
type TranslationJobSetupPhaseRuntimeSelection struct {
	PhaseID              TranslationJobSetupPhaseID
	Provider             string
	Model                string
	CredentialRef        string
	CredentialStatus     TranslationJobSetupCredentialStatus
	ExecutionMode        string
	BatchMode            TranslationJobSetupBatchMode
	ModelListSourceToken string
}

// ListTranslationJobSetupProviderModelsRequest carries the transport-stable provider model list input.
type ListTranslationJobSetupProviderModelsRequest struct {
	PhaseID          TranslationJobSetupPhaseID
	Provider         string
	CredentialRef    string
	CredentialStatus TranslationJobSetupCredentialStatus
	RequestToken     string
}

// TranslationJobSetupProviderModelOption is one selectable provider model.
type TranslationJobSetupProviderModelOption struct {
	ModelID string
	Label   string
}

// ListTranslationJobSetupProviderModelsResult returns the public provider model list state.
type ListTranslationJobSetupProviderModelsResult struct {
	PhaseID          TranslationJobSetupPhaseID
	Provider         string
	CredentialStatus TranslationJobSetupCredentialStatus
	RequestToken     string
	SourceToken      string
	Status           TranslationJobSetupProviderModelListStatus
	Models           []TranslationJobSetupProviderModelOption
	FailureKind      TranslationJobSetupErrorKind
}

// ValidateTranslationJobSetupRequest carries the transport-stable validation input.
type ValidateTranslationJobSetupRequest struct {
	InputSourceID          int64
	Runtime                TranslationJobSetupRuntimeSelection
	CredentialRef          string
	PhaseRuntimeSelections []TranslationJobSetupPhaseRuntimeSelection
}

// TranslationJobSetupPhaseValidationResult returns the validation state for one phase.
type TranslationJobSetupPhaseValidationResult struct {
	PhaseID                 TranslationJobSetupPhaseID
	Status                  TranslationJobSetupValidationStatus
	BlockingFailureCategory *string
	CanCreate               bool
	ModelListState          TranslationJobSetupProviderModelListStatus
	ModelListSourceToken    string
	IsModelSelectionStale   bool
}

// TranslationJobSetupValidationResult returns the validation decision and affected slices.
type TranslationJobSetupValidationResult struct {
	Status                  TranslationJobSetupValidationStatus
	BlockingFailureCategory *string
	TargetSlices            []string
	ValidatedAt             time.Time
	CanCreate               bool
	PassSlices              []string
	PhaseResults            []TranslationJobSetupPhaseValidationResult
	StaleModelListPhaseIDs  []TranslationJobSetupPhaseID
}

// CreateTranslationJobRequest carries the frozen job creation contract.
type CreateTranslationJobRequest struct {
	InputSourceID          int64
	InputSource            string
	ValidationStatus       TranslationJobSetupValidationStatus
	ValidatedAt            time.Time
	ValidationPassSlices   []string
	Runtime                TranslationJobSetupRuntimeSelection
	CredentialRef          string
	PhaseRuntimeSelections []TranslationJobSetupPhaseRuntimeSelection
}

// TranslationJobExecutionSummary returns the runtime configuration captured by a job.
type TranslationJobExecutionSummary struct {
	Provider      string
	Model         string
	ExecutionMode string
}

// TranslationJobSetupPhaseRuntimeSummary returns the runtime snapshot captured per phase.
type TranslationJobSetupPhaseRuntimeSummary struct {
	PhaseID              TranslationJobSetupPhaseID
	Provider             string
	Model                string
	CredentialRef        string
	CredentialStatus     TranslationJobSetupCredentialStatus
	ExecutionMode        string
	BatchMode            TranslationJobSetupBatchMode
	ModelListSourceToken string
}

// CreateTranslationJobResult returns either a ready job summary or a rejected error kind.
type CreateTranslationJobResult struct {
	JobID                 int64
	JobState              string
	InputSource           string
	ExecutionSummary      TranslationJobExecutionSummary
	ValidationPassSlices  []string
	PhaseRuntimeSummaries []TranslationJobSetupPhaseRuntimeSummary
	ErrorKind             TranslationJobSetupErrorKind
}

// GetTranslationJobSetupSummaryRequest identifies one created job.
type GetTranslationJobSetupSummaryRequest struct {
	JobID int64
}

// TranslationJobSetupSummaryResult returns the read-only job summary contract.
type TranslationJobSetupSummaryResult struct {
	JobID                 int64
	JobState              string
	InputSource           string
	CanStartPhase         bool
	ExecutionSummary      TranslationJobExecutionSummary
	ValidationPassSlices  []string
	PhaseRuntimeSummaries []TranslationJobSetupPhaseRuntimeSummary
}

// DeleteTranslationJobSetupInputRequest identifies one Job Setup input delete target.
type DeleteTranslationJobSetupInputRequest struct {
	InputSourceID int64
}

// DeleteTranslationJobSetupInputResult returns one delete outcome.
type DeleteTranslationJobSetupInputResult struct {
	DeletedInputSourceID *int64
	ErrorKind            TranslationJobSetupErrorKind
}

// NewTranslationJobSetupContractStub returns a temporary usecase stub for the frozen Wails seam.
func NewTranslationJobSetupContractStub() TranslationJobSetupContractStub {
	return TranslationJobSetupContractStub{}
}

// TranslationJobSetupContractStub is a temporary contract-only usecase used until the real Job Setup usecase exists.
type TranslationJobSetupContractStub struct{}

// GetTranslationJobSetupOptions returns a not-implemented error for the frozen contract seam.
func (TranslationJobSetupContractStub) GetTranslationJobSetupOptions(
	context.Context,
) (TranslationJobSetupOptionsResult, error) {
	return TranslationJobSetupOptionsResult{}, errTranslationJobSetupNotImplemented
}

// ValidateTranslationJobSetup returns a not-implemented error for the frozen contract seam.
func (TranslationJobSetupContractStub) ValidateTranslationJobSetup(
	ctx context.Context,
	request ValidateTranslationJobSetupRequest,
) (TranslationJobSetupValidationResult, error) {
	decision, err := jobsetupservice.NewTranslationJobSetupService().ValidateRequest(ctx, jobsetupservice.TranslationJobSetupValidationRequest{
		InputSourceID: request.InputSourceID,
		PhaseRuntimes: []jobsetupservice.TranslationJobSetupPhaseRuntimeDraftReadModel{{
			PhaseID:              string(TranslationJobSetupPhaseIDWordTranslation),
			Provider:             request.Runtime.Provider,
			Model:                request.Runtime.Model,
			CredentialRef:        request.CredentialRef,
			CredentialStatus:     string(TranslationJobSetupCredentialStatusConfigured),
			ExecutionMode:        request.Runtime.ExecutionMode,
			BatchMode:            string(TranslationJobSetupBatchModeUnsupported),
			ModelListSourceToken: "stub",
		}},
	})
	if err != nil {
		return TranslationJobSetupValidationResult{}, errors.New("validate translation job setup request")
	}
	return toTranslationJobSetupValidationResult(decision), nil
}

// ListTranslationJobSetupProviderModels returns a not-implemented error for the frozen provider model list seam.
func (TranslationJobSetupContractStub) ListTranslationJobSetupProviderModels(
	context.Context,
	ListTranslationJobSetupProviderModelsRequest,
) (ListTranslationJobSetupProviderModelsResult, error) {
	return ListTranslationJobSetupProviderModelsResult{}, errTranslationJobSetupNotImplemented
}

// DeleteTranslationJobSetupInput returns a not-implemented error for the frozen contract seam.
func (TranslationJobSetupContractStub) DeleteTranslationJobSetupInput(
	context.Context,
	DeleteTranslationJobSetupInputRequest,
) (DeleteTranslationJobSetupInputResult, error) {
	return DeleteTranslationJobSetupInputResult{}, errTranslationJobSetupNotImplemented
}

// CreateTranslationJob returns a not-implemented error for the frozen contract seam.
func (TranslationJobSetupContractStub) CreateTranslationJob(
	_ context.Context,
	request CreateTranslationJobRequest,
) (CreateTranslationJobResult, error) {
	if strings.ToLower(strings.TrimSpace(string(request.ValidationStatus))) != string(TranslationJobSetupValidationStatusPass) {
		return CreateTranslationJobResult{ErrorKind: TranslationJobSetupErrorKindValidationFailed}, nil
	}
	if translationJobSetupValidationIsStale(time.Now().UTC(), request.ValidatedAt.UTC()) {
		return CreateTranslationJobResult{ErrorKind: TranslationJobSetupErrorKindValidationStale}, nil
	}
	return CreateTranslationJobResult{}, errTranslationJobSetupNotImplemented
}

// GetTranslationJobSetupSummary returns a not-implemented error for the frozen contract seam.
func (TranslationJobSetupContractStub) GetTranslationJobSetupSummary(
	context.Context,
	GetTranslationJobSetupSummaryRequest,
) (TranslationJobSetupSummaryResult, error) {
	return TranslationJobSetupSummaryResult{}, errTranslationJobSetupNotImplemented
}

var errTranslationJobSetupNotImplemented = errors.New("translation job setup usecase is not implemented")

func translationJobSetupValidationIsStale(now time.Time, validatedAt time.Time) bool {
	if validatedAt.IsZero() {
		return true
	}
	return validatedAt.Before(translationJobSetupValidationFreshnessCutoff(now))
}

func translationJobSetupValidationFreshnessCutoff(now time.Time) time.Time {
	nowUTC := now.UTC()
	cutoff := time.Date(
		nowUTC.Year(),
		nowUTC.Month(),
		nowUTC.Day(),
		translationJobSetupValidationFreshnessCutoffHourUTC,
		0,
		0,
		0,
		time.UTC,
	)
	if nowUTC.Before(cutoff) {
		return cutoff.AddDate(0, 0, -1)
	}
	return cutoff
}
