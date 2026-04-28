package usecase

import (
	"context"
	"errors"
	"strings"
	"time"

	jobsetupservice "aitranslationenginejp/internal/service"
)

// TranslationJobSetupErrorKind identifies contract-level rejected outcomes.
type TranslationJobSetupErrorKind = string

const (
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

	// TranslationJobSetupErrorKindValidationFailed remains as a compatibility alias during downstream alignment.
	TranslationJobSetupErrorKindValidationFailed TranslationJobSetupErrorKind = TranslationJobSetupErrorKindReadyRequired
	// TranslationJobSetupErrorKindDuplicateInput remains as a compatibility alias during downstream alignment.
	TranslationJobSetupErrorKindDuplicateInput TranslationJobSetupErrorKind = TranslationJobSetupErrorKindDuplicateJobForInput
)

// NormalizeTranslationJobSetupPublicErrorKind collapses internal compatibility aliases to the frozen public kinds.
func NormalizeTranslationJobSetupPublicErrorKind(kind TranslationJobSetupErrorKind) TranslationJobSetupErrorKind {
	trimmedKind := strings.TrimSpace(kind)
	switch strings.ToLower(trimmedKind) {
	case "":
		return ""
	case "validation_failed", TranslationJobSetupErrorKindReadyRequired:
		return TranslationJobSetupErrorKindReadyRequired
	case "duplicate_input", TranslationJobSetupErrorKindDuplicateJobForInput:
		return TranslationJobSetupErrorKindDuplicateJobForInput
	default:
		return trimmedKind
	}
}

// NormalizeTranslationJobSetupPublicErrorCategory normalizes optional public error categories.
func NormalizeTranslationJobSetupPublicErrorCategory(category *string) *string {
	if category == nil {
		return nil
	}
	normalized := NormalizeTranslationJobSetupPublicErrorKind(*category)
	return &normalized
}

const translationJobSetupValidationFreshnessCutoffHourUTC = 9

// TranslationJobSetupValidationStatus identifies the outcome of one setup validation.
type TranslationJobSetupValidationStatus = string

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
	InputCandidates    []TranslationJobSetupInputCandidate
	ExistingJob        *TranslationJobSetupExistingJob
	SharedDictionaries []TranslationJobSetupDictionaryOption
	SharedPersonas     []TranslationJobSetupPersonaOption
	AIRuntimeOptions   []TranslationJobSetupRuntimeOption
	CredentialRefs     []TranslationJobSetupCredentialReference
}

// TranslationJobSetupInputCandidate is one selectable translation input source.
type TranslationJobSetupInputCandidate struct {
	ID           int64
	Label        string
	SourceKind   string
	RecordCount  int
	RegisteredAt time.Time
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

// TranslationJobSetupRuntimeSelection is the selected runtime configuration.
type TranslationJobSetupRuntimeSelection struct {
	Provider      string
	Model         string
	ExecutionMode string
}

// ValidateTranslationJobSetupRequest carries the transport-stable validation input.
type ValidateTranslationJobSetupRequest struct {
	InputSourceID int64
	Runtime       TranslationJobSetupRuntimeSelection
	CredentialRef string
}

// TranslationJobSetupValidationResult returns the validation decision and affected slices.
type TranslationJobSetupValidationResult struct {
	Status                  TranslationJobSetupValidationStatus
	BlockingFailureCategory *string
	TargetSlices            []string
	ValidatedAt             time.Time
	CanCreate               bool
	PassSlices              []string
}

// CreateTranslationJobRequest carries the frozen job creation contract.
type CreateTranslationJobRequest struct {
	InputSourceID        int64
	InputSource          string
	ValidationStatus     TranslationJobSetupValidationStatus
	ValidatedAt          time.Time
	ValidationPassSlices []string
	Runtime              TranslationJobSetupRuntimeSelection
	CredentialRef        string
}

// TranslationJobExecutionSummary returns the runtime configuration captured by a job.
type TranslationJobExecutionSummary struct {
	Provider      string
	Model         string
	ExecutionMode string
}

// CreateTranslationJobResult returns either a ready job summary or a rejected error kind.
type CreateTranslationJobResult struct {
	JobID                int64
	JobState             string
	InputSource          string
	ExecutionSummary     TranslationJobExecutionSummary
	ValidationPassSlices []string
	ErrorKind            TranslationJobSetupErrorKind
}

// GetTranslationJobSetupSummaryRequest identifies one created job.
type GetTranslationJobSetupSummaryRequest struct {
	JobID int64
}

// TranslationJobSetupSummaryResult returns the read-only job summary contract.
type TranslationJobSetupSummaryResult struct {
	JobID                int64
	JobState             string
	InputSource          string
	CanStartPhase        bool
	ExecutionSummary     TranslationJobExecutionSummary
	ValidationPassSlices []string
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
		Provider:      request.Runtime.Provider,
		Model:         request.Runtime.Model,
		ExecutionMode: request.Runtime.ExecutionMode,
		CredentialRef: request.CredentialRef,
	})
	if err != nil {
		return TranslationJobSetupValidationResult{}, errors.New("validate translation job setup request")
	}
	return toTranslationJobSetupValidationResult(decision), nil
}

// CreateTranslationJob returns a not-implemented error for the frozen contract seam.
func (TranslationJobSetupContractStub) CreateTranslationJob(
	_ context.Context,
	request CreateTranslationJobRequest,
) (CreateTranslationJobResult, error) {
	if strings.ToLower(strings.TrimSpace(request.ValidationStatus)) != TranslationJobSetupValidationStatusPass {
		return CreateTranslationJobResult{ErrorKind: TranslationJobSetupErrorKindValidationFailed}, nil
	}
	if request.ValidatedAt.IsZero() || request.ValidatedAt.UTC().Hour() < translationJobSetupValidationFreshnessCutoffHourUTC {
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
