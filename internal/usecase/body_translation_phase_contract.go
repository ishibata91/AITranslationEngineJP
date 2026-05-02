package usecase

import (
	"context"
	"errors"
	"strings"
	"time"
)

// BodyTranslationPhaseErrorKind identifies contract-level rejected outcomes.
type BodyTranslationPhaseErrorKind string

const (
	// BodyTranslationPhaseErrorKindPersonaPhaseIncomplete identifies incomplete persona phase state.
	BodyTranslationPhaseErrorKindPersonaPhaseIncomplete BodyTranslationPhaseErrorKind = "persona_phase_incomplete"
	// BodyTranslationPhaseErrorKindTerminalJob identifies a terminal job rejection.
	BodyTranslationPhaseErrorKindTerminalJob BodyTranslationPhaseErrorKind = "terminal_job"
	// BodyTranslationPhaseErrorKindActivePhaseExists identifies an already-active phase run.
	BodyTranslationPhaseErrorKindActivePhaseExists BodyTranslationPhaseErrorKind = "active_phase_exists"
	// BodyTranslationPhaseErrorKindInputSnapshotFailed identifies input snapshot creation failure.
	BodyTranslationPhaseErrorKindInputSnapshotFailed BodyTranslationPhaseErrorKind = "input_snapshot_failed"
	// BodyTranslationPhaseErrorKindProviderFailure identifies provider execution failure.
	BodyTranslationPhaseErrorKindProviderFailure BodyTranslationPhaseErrorKind = "provider_failure"
	// BodyTranslationPhaseErrorKindInvalidProviderResponse identifies invalid provider response data.
	BodyTranslationPhaseErrorKindInvalidProviderResponse BodyTranslationPhaseErrorKind = "invalid_provider_response"
	// BodyTranslationPhaseErrorKindProtectionValidationFailed identifies translation protection validation failure.
	BodyTranslationPhaseErrorKindProtectionValidationFailed BodyTranslationPhaseErrorKind = "protection_validation_failed"
	// BodyTranslationPhaseErrorKindSaveFailed identifies persistence failure.
	BodyTranslationPhaseErrorKindSaveFailed BodyTranslationPhaseErrorKind = "save_failed"
	// BodyTranslationPhaseErrorKindOutputReadinessBlocked identifies blocked downstream readiness.
	BodyTranslationPhaseErrorKindOutputReadinessBlocked BodyTranslationPhaseErrorKind = "output_readiness_blocked"
	// BodyTranslationPhaseErrorKindLateResponseRejected identifies rejected late provider responses.
	BodyTranslationPhaseErrorKindLateResponseRejected BodyTranslationPhaseErrorKind = "late_response_rejected"
	// BodyTranslationPhaseErrorKindSecretRedacted identifies redacted secret-related failures.
	BodyTranslationPhaseErrorKindSecretRedacted BodyTranslationPhaseErrorKind = "secret_redacted"
)

// NormalizeBodyTranslationPhasePublicErrorKind trims and lowercases the frozen public error kinds.
func NormalizeBodyTranslationPhasePublicErrorKind(kind BodyTranslationPhaseErrorKind) BodyTranslationPhaseErrorKind {
	trimmedKind := strings.TrimSpace(string(kind))
	if trimmedKind == "" {
		return ""
	}
	return BodyTranslationPhaseErrorKind(strings.ToLower(trimmedKind))
}

// BodyTranslationPhaseRedactionFieldObligation freezes the public redaction boundary.
type BodyTranslationPhaseRedactionFieldObligation struct {
	AllowedReferenceFields []string
	ForbiddenOutputFields  []string
}

// BodyTranslationPhasePublicRedactionFieldObligation returns the frozen public redaction boundary.
func BodyTranslationPhasePublicRedactionFieldObligation() BodyTranslationPhaseRedactionFieldObligation {
	return BodyTranslationPhaseRedactionFieldObligation{
		AllowedReferenceFields: []string{
			"credential_ref",
			"provider",
			"model",
			"execution_mode",
			"phase_run_id",
			"input_snapshot_digest",
			"target_count",
			"dictionary_digest",
			"persona_digest",
			"metadata_digest",
			"prompt_digest",
			"request_unit_count",
			"output_count",
			"error_kind",
		},
		ForbiddenOutputFields: []string{
			"secret",
			"api_key",
			"token",
			"authorization_header",
			"provider_raw_request",
			"provider_raw_response",
			"raw_prompt",
			"decryptable_value",
		},
	}
}

// GetBodyTranslationPhaseSummaryRequest identifies the requested job.
type GetBodyTranslationPhaseSummaryRequest struct {
	JobID int64
}

// StartBodyTranslationPhaseRequest identifies the ready job to start.
type StartBodyTranslationPhaseRequest struct {
	JobID int64
}

// PauseBodyTranslationPhaseRequest identifies one active phase run to pause.
type PauseBodyTranslationPhaseRequest struct {
	JobID      int64
	PhaseRunID int64
}

// ResumeBodyTranslationPhaseRequest identifies one paused or recoverable failed phase run to resume.
type ResumeBodyTranslationPhaseRequest struct {
	JobID      int64
	PhaseRunID int64
}

// RetryBodyTranslationPhaseRequest identifies one retryable failed phase run to retry.
type RetryBodyTranslationPhaseRequest struct {
	JobID      int64
	PhaseRunID int64
}

// CancelBodyTranslationPhaseRequest identifies one cancelable phase run.
type CancelBodyTranslationPhaseRequest struct {
	JobID      int64
	PhaseRunID int64
}

// GetBodyTranslationOutputReadinessRequest identifies the job whose output readiness is requested.
type GetBodyTranslationOutputReadinessRequest struct {
	JobID int64
}

// BodyTranslationPhaseProgressSummary summarizes current progress for one phase run.
type BodyTranslationPhaseProgressSummary struct {
	Percent         int
	ProcessedCount  int
	TotalCount      int
	TargetCount     int
	TranslatedCount int
	SkippedCount    int
	CurrentStep     string
}

// BodyTranslationPhaseInputSummary summarizes the frozen body phase input snapshot.
type BodyTranslationPhaseInputSummary struct {
	TargetCount      int
	SkippedReasons   []string
	InputSnapshotRef *string
	DictionaryDigest string
	PersonaDigest    string
	MetadataDigest   string
	PromptDigest     string
}

// BodyTranslationInputSummary is a compatibility alias for the frozen public input summary.
type BodyTranslationInputSummary = BodyTranslationPhaseInputSummary

// BodyTranslationPhaseExecutionSummary summarizes execution inputs without exposing secrets.
type BodyTranslationPhaseExecutionSummary struct {
	CredentialRef    string
	Provider         string
	Model            string
	ExecutionMode    string
	RequestUnitCount int
	OutputCount      int
}

// BodyTranslationExecutionSummary is a compatibility alias for the frozen public execution summary.
type BodyTranslationExecutionSummary = BodyTranslationPhaseExecutionSummary

// BodyTranslationPhaseRequestSummary summarizes provider request units after input filtering.
type BodyTranslationPhaseRequestSummary struct {
	ProviderTargetCount              int
	ExactDictionaryExclusionCount    int
	PartialDictionaryConstraintCount int
}

// BodyTranslationPhaseFieldResultSummary summarizes stored field results.
type BodyTranslationPhaseFieldResultSummary struct {
	TranslatedCount       int
	FailedCount           int
	SkippedCount          int
	ProtectionFailedCount int
	OutputReadyCount      int
	OutputCount           int
	FieldResults          []BodyTranslationPhaseFieldResultItem
}

// BodyTranslationFieldResultSummary is a compatibility alias for the frozen public field result summary.
type BodyTranslationFieldResultSummary = BodyTranslationPhaseFieldResultSummary

// BodyTranslationPhaseFieldIdentity identifies one translated field.
type BodyTranslationPhaseFieldIdentity struct {
	TranslationFieldID      int64
	PhaseTranslationFieldID int64
	RecordType              string
	FieldType               string
	FormID                  string
	EditorID                string
	FieldLabel              string
}

// BodyTranslationPhaseFieldResultItem stores one public field result row.
type BodyTranslationPhaseFieldResultItem struct {
	Identity                    BodyTranslationPhaseFieldIdentity
	FieldID                     int64
	FieldLabel                  string
	SourceExcerpt               string
	TranslatedText              string
	OutputStatus                string
	ProtectionValidationResult  string
	ProtectionValidationSummary string
	RetryCount                  int
}

// BodyTranslationPhaseErrorSummary exposes only redacted error information.
type BodyTranslationPhaseErrorSummary struct {
	ErrorKind  BodyTranslationPhaseErrorKind
	Reason     string
	Retryable  bool
	IsRedacted bool
}

// BodyTranslationPhaseActionEnablement summarizes Job Run button state.
type BodyTranslationPhaseActionEnablement struct {
	CanStart                     bool
	StartBlockedReason           *string
	CanPause                     bool
	PauseBlockedReason           *string
	CanResume                    bool
	ResumeBlockedReason          *string
	CanRetry                     bool
	RetryBlockedReason           *string
	CanCancel                    bool
	CancelBlockedReason          *string
	CanCheckOutputReadiness      bool
	OutputReadinessBlockedReason *string
}

// BodyTranslationOutputReadinessSummary summarizes downstream output readiness.
type BodyTranslationOutputReadinessSummary struct {
	Ready               bool
	BlockedReason       string
	ErrorKind           BodyTranslationPhaseErrorKind
	CompletedFieldCount int
	StatusConsistent    bool
}

// BodyTranslationPhaseSummaryResult is the frozen Job Run summary contract.
type BodyTranslationPhaseSummaryResult struct {
	JobID              int64
	CurrentPhase       string
	PhaseState         string
	PhaseRunID         *int64
	StartedAt          *time.Time
	FinishedAt         *time.Time
	Progress           BodyTranslationPhaseProgressSummary
	InputSummary       BodyTranslationPhaseInputSummary
	RequestSummary     BodyTranslationPhaseRequestSummary
	Execution          BodyTranslationPhaseExecutionSummary
	FieldResultSummary *BodyTranslationFieldResultSummary
	ResultSummary      *BodyTranslationPhaseFieldResultSummary
	FieldResults       []BodyTranslationPhaseFieldResultItem
	ErrorSummary       *BodyTranslationPhaseErrorSummary
	ActionEnablement   BodyTranslationPhaseActionEnablement
	OutputReadiness    BodyTranslationOutputReadinessSummary
}

// BodyTranslationPhaseCommandResult is the frozen write-seam response contract.
type BodyTranslationPhaseCommandResult struct {
	JobID               int64
	CurrentPhase        string
	PhaseState          string
	PhaseRunID          *int64
	StartedAt           *time.Time
	FinishedAt          *time.Time
	Progress            BodyTranslationPhaseProgressSummary
	InputSnapshotDigest string
	InputSummary        BodyTranslationPhaseInputSummary
	RequestSummary      BodyTranslationPhaseRequestSummary
	Execution           BodyTranslationPhaseExecutionSummary
	FieldResultSummary  *BodyTranslationFieldResultSummary
	ResultSummary       *BodyTranslationPhaseFieldResultSummary
	FieldResults        []BodyTranslationPhaseFieldResultItem
	Retryable           bool
	OutputReadiness     BodyTranslationOutputReadinessSummary
	ErrorSummary        *BodyTranslationPhaseErrorSummary
}

// BodyTranslationOutputReadinessResult is the frozen downstream readiness contract.
type BodyTranslationOutputReadinessResult struct {
	JobID               int64
	CurrentPhase        string
	PhaseState          string
	Ready               bool
	BlockedReason       string
	ErrorKind           BodyTranslationPhaseErrorKind
	CompletedFieldCount int
	StatusConsistent    bool
	OutputCount         int
}

// NewBodyTranslationPhaseContractStub returns a temporary usecase stub for the frozen Wails seam.
func NewBodyTranslationPhaseContractStub() BodyTranslationPhaseContractStub {
	return BodyTranslationPhaseContractStub{}
}

// BodyTranslationPhaseContractStub is a temporary contract-only usecase used until the real phase usecase exists.
type BodyTranslationPhaseContractStub struct{}

// GetBodyTranslationPhaseSummary returns a not-implemented error for the frozen contract seam.
func (BodyTranslationPhaseContractStub) GetBodyTranslationPhaseSummary(
	context.Context,
	GetBodyTranslationPhaseSummaryRequest,
) (BodyTranslationPhaseSummaryResult, error) {
	return BodyTranslationPhaseSummaryResult{}, errBodyTranslationPhaseNotImplemented
}

// StartBodyTranslationPhase returns a not-implemented error for the frozen contract seam.
func (BodyTranslationPhaseContractStub) StartBodyTranslationPhase(
	context.Context,
	StartBodyTranslationPhaseRequest,
) (BodyTranslationPhaseCommandResult, error) {
	return BodyTranslationPhaseCommandResult{}, errBodyTranslationPhaseNotImplemented
}

// PauseBodyTranslationPhase returns a not-implemented error for the frozen contract seam.
func (BodyTranslationPhaseContractStub) PauseBodyTranslationPhase(
	context.Context,
	PauseBodyTranslationPhaseRequest,
) (BodyTranslationPhaseCommandResult, error) {
	return BodyTranslationPhaseCommandResult{}, errBodyTranslationPhaseNotImplemented
}

// ResumeBodyTranslationPhase returns a not-implemented error for the frozen contract seam.
func (BodyTranslationPhaseContractStub) ResumeBodyTranslationPhase(
	context.Context,
	ResumeBodyTranslationPhaseRequest,
) (BodyTranslationPhaseCommandResult, error) {
	return BodyTranslationPhaseCommandResult{}, errBodyTranslationPhaseNotImplemented
}

// RetryBodyTranslationPhase returns a not-implemented error for the frozen contract seam.
func (BodyTranslationPhaseContractStub) RetryBodyTranslationPhase(
	context.Context,
	RetryBodyTranslationPhaseRequest,
) (BodyTranslationPhaseCommandResult, error) {
	return BodyTranslationPhaseCommandResult{}, errBodyTranslationPhaseNotImplemented
}

// CancelBodyTranslationPhase returns a not-implemented error for the frozen contract seam.
func (BodyTranslationPhaseContractStub) CancelBodyTranslationPhase(
	context.Context,
	CancelBodyTranslationPhaseRequest,
) (BodyTranslationPhaseCommandResult, error) {
	return BodyTranslationPhaseCommandResult{}, errBodyTranslationPhaseNotImplemented
}

// GetBodyTranslationOutputReadiness returns a not-implemented error for the frozen contract seam.
func (BodyTranslationPhaseContractStub) GetBodyTranslationOutputReadiness(
	context.Context,
	GetBodyTranslationOutputReadinessRequest,
) (BodyTranslationOutputReadinessResult, error) {
	return BodyTranslationOutputReadinessResult{}, errBodyTranslationPhaseNotImplemented
}

var errBodyTranslationPhaseNotImplemented = errors.New("body translation phase usecase is not implemented")
