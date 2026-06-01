package usecase

import (
	"context"
	"errors"
	"strings"
	"time"
)

// TermTranslationPhaseErrorKind identifies contract-level rejected outcomes.
type TermTranslationPhaseErrorKind = string

const (
	// TermTranslationPhaseErrorKindReadyRequired identifies a rejected outcome caused by a non-ready job.
	TermTranslationPhaseErrorKindReadyRequired TermTranslationPhaseErrorKind = "ready_required"
	// TermTranslationPhaseErrorKindTerminalJob identifies a rejected outcome caused by a terminal job.
	TermTranslationPhaseErrorKindTerminalJob TermTranslationPhaseErrorKind = "terminal_job"
	// TermTranslationPhaseErrorKindActivePhaseExists identifies a rejected outcome caused by an existing active phase run.
	TermTranslationPhaseErrorKindActivePhaseExists TermTranslationPhaseErrorKind = "active_phase_exists"
	// TermTranslationPhaseErrorKindDictionarySnapshotFail identifies a rejected outcome caused by dictionary snapshot creation failure.
	TermTranslationPhaseErrorKindDictionarySnapshotFail TermTranslationPhaseErrorKind = "dictionary_snapshot_failed"
	// TermTranslationPhaseErrorKindProviderFailure identifies a rejected outcome caused by provider execution failure.
	TermTranslationPhaseErrorKindProviderFailure TermTranslationPhaseErrorKind = "provider_failure"
	// TermTranslationPhaseErrorKindInvalidProviderReply identifies a rejected outcome caused by invalid provider response data.
	TermTranslationPhaseErrorKindInvalidProviderReply TermTranslationPhaseErrorKind = "invalid_provider_response"
	// TermTranslationPhaseErrorKindSaveFailed identifies a rejected outcome caused by save failure.
	TermTranslationPhaseErrorKindSaveFailed TermTranslationPhaseErrorKind = "save_failed"
	// TermTranslationPhaseErrorKindTermPhaseIncomplete identifies a rejected outcome caused by incomplete term phase state.
	TermTranslationPhaseErrorKindTermPhaseIncomplete TermTranslationPhaseErrorKind = "term_phase_incomplete"
	// TermTranslationPhaseErrorKindSecretRedacted identifies a redacted secret-related failure summary.
	TermTranslationPhaseErrorKindSecretRedacted TermTranslationPhaseErrorKind = "secret_redacted"
)

// NormalizeTermTranslationPhasePublicErrorKind trims and lowercases the frozen public error kinds.
func NormalizeTermTranslationPhasePublicErrorKind(kind TermTranslationPhaseErrorKind) TermTranslationPhaseErrorKind {
	trimmedKind := strings.TrimSpace(kind)
	if trimmedKind == "" {
		return ""
	}
	return strings.ToLower(trimmedKind)
}

// GetTermTranslationPhaseSummaryRequest identifies the requested translation job.
type GetTermTranslationPhaseSummaryRequest struct {
	JobID int64
}

// StartTermTranslationPhaseRequest identifies the ready job that wants to create a phase run.
type StartTermTranslationPhaseRequest struct {
	JobID int64
}

// PauseTermTranslationPhaseRequest identifies one active phase run to pause.
type PauseTermTranslationPhaseRequest struct {
	JobID      int64
	PhaseRunID int64
}

// ResumeTermTranslationPhaseRequest identifies one paused or recoverable failed phase run to resume.
type ResumeTermTranslationPhaseRequest struct {
	JobID      int64
	PhaseRunID int64
}

// RetryTermTranslationPhaseRequest identifies one retryable failed phase run to retry.
type RetryTermTranslationPhaseRequest struct {
	JobID      int64
	PhaseRunID int64
}

// GetTermTranslationNextPhaseReadinessRequest identifies the job whose downstream readiness is requested.
type GetTermTranslationNextPhaseReadinessRequest struct {
	JobID int64
}

// SaveTermTranslationPhaseAISettingsRequest carries public AI settings for the term phase.
// 入力に job_id を含めない。保存操作はジョブと無関連で phase_type を主キーとして upsert する。
type SaveTermTranslationPhaseAISettingsRequest struct {
	Provider      string
	Model         string
	ExecutionMode string
	BatchMode     string
}

// TermTranslationPhaseAISettingsResult returns the saved AI settings state.
// credential_ref、credential_status、model_list_status は含まない。
type TermTranslationPhaseAISettingsResult struct {
	PhaseType     string
	Provider      string
	Model         string
	ExecutionMode string
	BatchMode     string
}

// TermTranslationExecutionConfigSummary summarizes the execution inputs without exposing secrets.
// JOB_PHASE_RUN が存在する場合だけ設定する。Ready 期は nil とする。
type TermTranslationExecutionConfigSummary struct {
	CredentialRef   string
	Provider        string
	Model           string
	ExecutionMode   string
	SnapshotDigest  *string
	SnapshotVersion *string
}

// TermTranslationPhaseAISettings holds the Ready 期 AI 選択値（JOB_PHASE_AI_SETTINGS から読む）。
// record 不在（AI 設定未保存）の場合は nil とする。
type TermTranslationPhaseAISettings struct {
	Provider      string
	Model         string
	ExecutionMode string
	BatchMode     string
}

// TermTranslationPhaseProgressSummary summarizes current progress for one phase run.
type TermTranslationPhaseProgressSummary struct {
	Percent        int
	ProcessedCount int
	TotalCount     int
	AITargetCount  int
	CurrentStep    string
}

// TermTranslationPhaseResultSummary summarizes the phase output shown by Job Run.
type TermTranslationPhaseResultSummary struct {
	ConfirmedCount            int
	JobDictionaryAppliedCount int
	ReplacementTargetCount    int
	UnmatchedCount            int
	ProviderSkipped           bool
}

// TermTranslationPhaseErrorSummary exposes only redacted error information.
type TermTranslationPhaseErrorSummary struct {
	ErrorKind  TermTranslationPhaseErrorKind
	Reason     string
	Retryable  bool
	IsRedacted bool
}

// TermTranslationPhaseActionEnablement summarizes Job Run button state.
type TermTranslationPhaseActionEnablement struct {
	CanStart            bool
	StartBlockedReason  *string
	CanPause            bool
	PauseBlockedReason  *string
	CanResume           bool
	ResumeBlockedReason *string
	CanRetry            bool
	RetryBlockedReason  *string
}

// TermTranslationPhaseSummaryResult is the frozen Job Run summary contract.
type TermTranslationPhaseSummaryResult struct {
	JobID              int64
	CurrentPhase       string
	PhaseState         string
	PhaseRunID         *int64
	StartedAt          *time.Time
	FinishedAt         *time.Time
	Progress           TermTranslationPhaseProgressSummary
	TotalTermCount     int
	DictionaryHitCount int
	AITargetCount      int
	// AISettings は JOB_PHASE_AI_SETTINGS record が存在する場合だけ設定する。nil は未設定を意味する。
	AISettings *TermTranslationPhaseAISettings
	// Execution は JOB_PHASE_RUN が存在し AI 設定転写済みの場合だけ設定する。nil は未開始を意味する。
	Execution        *TermTranslationExecutionConfigSummary
	ResultSummary    *TermTranslationPhaseResultSummary
	ErrorSummary     *TermTranslationPhaseErrorSummary
	ActionEnablement TermTranslationPhaseActionEnablement
}

// TermTranslationPhaseCommandResult is the frozen write-seam response contract.
type TermTranslationPhaseCommandResult struct {
	JobID        int64
	CurrentPhase string
	PhaseState   string
	PhaseRunID   *int64
	StartedAt    *time.Time
	FinishedAt   *time.Time
	Progress     TermTranslationPhaseProgressSummary
	Retryable    bool
	ErrorSummary *TermTranslationPhaseErrorSummary
}

// TermTranslationNextPhaseReadinessResult is the frozen downstream phase fact state contract.
type TermTranslationNextPhaseReadinessResult struct {
	JobID          int64
	CurrentPhase   string
	PhaseState     string
	JobIsTerminal  bool
	TotalCount     int
	ConfirmedCount int
	ErrorKind      TermTranslationPhaseErrorKind
}

// NewTermTranslationPhaseContractStub returns a temporary usecase stub for the frozen Wails seam.
func NewTermTranslationPhaseContractStub() TermTranslationPhaseContractStub {
	return TermTranslationPhaseContractStub{}
}

// TermTranslationPhaseContractStub is a temporary contract-only usecase used until the real phase usecase exists.
type TermTranslationPhaseContractStub struct{}

// GetTermTranslationPhaseSummary returns a not-implemented error for the frozen contract seam.
func (TermTranslationPhaseContractStub) GetTermTranslationPhaseSummary(
	context.Context,
	GetTermTranslationPhaseSummaryRequest,
) (TermTranslationPhaseSummaryResult, error) {
	return TermTranslationPhaseSummaryResult{}, errTermTranslationPhaseNotImplemented
}

// StartTermTranslationPhase returns a not-implemented error for the frozen contract seam.
func (TermTranslationPhaseContractStub) StartTermTranslationPhase(
	context.Context,
	StartTermTranslationPhaseRequest,
) (TermTranslationPhaseCommandResult, error) {
	return TermTranslationPhaseCommandResult{}, errTermTranslationPhaseNotImplemented
}

// PauseTermTranslationPhase returns a not-implemented error for the frozen contract seam.
func (TermTranslationPhaseContractStub) PauseTermTranslationPhase(
	context.Context,
	PauseTermTranslationPhaseRequest,
) (TermTranslationPhaseCommandResult, error) {
	return TermTranslationPhaseCommandResult{}, errTermTranslationPhaseNotImplemented
}

// ResumeTermTranslationPhase returns a not-implemented error for the frozen contract seam.
func (TermTranslationPhaseContractStub) ResumeTermTranslationPhase(
	context.Context,
	ResumeTermTranslationPhaseRequest,
) (TermTranslationPhaseCommandResult, error) {
	return TermTranslationPhaseCommandResult{}, errTermTranslationPhaseNotImplemented
}

// RetryTermTranslationPhase returns a not-implemented error for the frozen contract seam.
func (TermTranslationPhaseContractStub) RetryTermTranslationPhase(
	context.Context,
	RetryTermTranslationPhaseRequest,
) (TermTranslationPhaseCommandResult, error) {
	return TermTranslationPhaseCommandResult{}, errTermTranslationPhaseNotImplemented
}

// GetTermTranslationNextPhaseReadiness returns a not-implemented error for the frozen contract seam.
func (TermTranslationPhaseContractStub) GetTermTranslationNextPhaseReadiness(
	context.Context,
	GetTermTranslationNextPhaseReadinessRequest,
) (TermTranslationNextPhaseReadinessResult, error) {
	return TermTranslationNextPhaseReadinessResult{}, errTermTranslationPhaseNotImplemented
}

var errTermTranslationPhaseNotImplemented = errors.New("term translation phase usecase is not implemented")
