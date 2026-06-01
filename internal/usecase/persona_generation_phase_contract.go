package usecase

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"
)

// PersonaGenerationPhaseErrorKind identifies contract-level rejected outcomes.
type PersonaGenerationPhaseErrorKind = string

const (
	// PersonaGenerationPhaseErrorKindTermPhaseIncomplete identifies incomplete term phase state.
	PersonaGenerationPhaseErrorKindTermPhaseIncomplete PersonaGenerationPhaseErrorKind = "term_phase_incomplete"
	// PersonaGenerationPhaseErrorKindTerminalJob identifies a terminal job rejection.
	PersonaGenerationPhaseErrorKindTerminalJob PersonaGenerationPhaseErrorKind = "terminal_job"
	// PersonaGenerationPhaseErrorKindActivePhaseExists identifies an already-active phase run.
	PersonaGenerationPhaseErrorKindActivePhaseExists PersonaGenerationPhaseErrorKind = "active_phase_exists"
	// PersonaGenerationPhaseErrorKindTargetSnapshotFailed identifies target snapshot creation failure.
	PersonaGenerationPhaseErrorKindTargetSnapshotFailed PersonaGenerationPhaseErrorKind = "target_snapshot_failed"
	// PersonaGenerationPhaseErrorKindProviderFailure identifies provider execution failure.
	PersonaGenerationPhaseErrorKindProviderFailure PersonaGenerationPhaseErrorKind = "provider_failure"
	// PersonaGenerationPhaseErrorKindInvalidProviderResponse identifies invalid provider response data.
	PersonaGenerationPhaseErrorKindInvalidProviderResponse PersonaGenerationPhaseErrorKind = "invalid_provider_response"
	// PersonaGenerationPhaseErrorKindSaveFailed identifies persona persistence failure.
	PersonaGenerationPhaseErrorKindSaveFailed PersonaGenerationPhaseErrorKind = "save_failed"
	// PersonaGenerationPhaseErrorKindSnapshotMissing identifies missing persona snapshot references.
	PersonaGenerationPhaseErrorKindSnapshotMissing PersonaGenerationPhaseErrorKind = "snapshot_missing"
	// PersonaGenerationPhaseErrorKindBodyReadinessBlocked identifies blocked body phase readiness.
	PersonaGenerationPhaseErrorKindBodyReadinessBlocked PersonaGenerationPhaseErrorKind = "body_readiness_blocked"
	// PersonaGenerationPhaseErrorKindSecretRedacted identifies redacted secret-related failures.
	PersonaGenerationPhaseErrorKindSecretRedacted PersonaGenerationPhaseErrorKind = "secret_redacted"
	// PersonaGenerationPhaseErrorKindInputMissing identifies missing redacted execution input.
	PersonaGenerationPhaseErrorKindInputMissing PersonaGenerationPhaseErrorKind = "input_missing"
)

const (
	personaGenerationFixtureEvidenceNPC001 = "evidence:npc:001"
)

// NormalizePersonaGenerationPhasePublicErrorKind trims and lowercases the frozen public error kinds.
func NormalizePersonaGenerationPhasePublicErrorKind(kind PersonaGenerationPhaseErrorKind) PersonaGenerationPhaseErrorKind {
	trimmedKind := strings.TrimSpace(kind)
	if trimmedKind == "" {
		return ""
	}
	return strings.ToLower(trimmedKind)
}

// PersonaGenerationPhaseRedactionFieldObligation freezes the public redaction boundary.
type PersonaGenerationPhaseRedactionFieldObligation struct {
	AllowedReferenceFields []string
	ForbiddenOutputFields  []string
}

// PersonaGenerationPhasePublicRedactionFieldObligation returns the frozen public redaction boundary.
func PersonaGenerationPhasePublicRedactionFieldObligation() PersonaGenerationPhaseRedactionFieldObligation {
	return PersonaGenerationPhaseRedactionFieldObligation{
		AllowedReferenceFields: []string{
			"credential_ref",
			"provider",
			"model",
			"execution_mode",
			"phase_run_id",
			"snapshot_id",
			"snapshot_digest",
			"prompt_digest",
			"evidence_ref",
			"input_count",
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
			"source_utterance_full_text",
			"conversation_context_full_text",
		},
	}
}

// GetPersonaGenerationPhaseSummaryRequest identifies the requested job.
type GetPersonaGenerationPhaseSummaryRequest struct {
	JobID int64
}

// StartPersonaGenerationPhaseRequest identifies the ready job to start.
type StartPersonaGenerationPhaseRequest struct {
	JobID int64
}

// PausePersonaGenerationPhaseRequest identifies one active phase run to pause.
type PausePersonaGenerationPhaseRequest struct {
	JobID      int64
	PhaseRunID int64
}

// ResumePersonaGenerationPhaseRequest identifies one paused or recoverable failed phase run to resume.
type ResumePersonaGenerationPhaseRequest struct {
	JobID      int64
	PhaseRunID int64
}

// RetryPersonaGenerationPhaseRequest identifies one retryable failed phase run to retry.
type RetryPersonaGenerationPhaseRequest struct {
	JobID      int64
	PhaseRunID int64
}

// CancelPersonaGenerationPhaseRequest identifies one cancelable phase run.
type CancelPersonaGenerationPhaseRequest struct {
	JobID      int64
	PhaseRunID int64
}

// GetPersonaGenerationBodyReadinessRequest identifies the job whose body readiness is requested.
type GetPersonaGenerationBodyReadinessRequest struct {
	JobID int64
}

// SavePersonaGenerationPhaseAISettingsRequest carries public AI settings for the persona phase.
type SavePersonaGenerationPhaseAISettingsRequest struct {
	JobID         int64
	Provider      string
	Model         string
	ExecutionMode string
	BatchMode     string
}

// PersonaGenerationPhaseAISettingsResult returns saved AI settings state.
// credential_status、model_list_status は含まない。
type PersonaGenerationPhaseAISettingsResult struct {
	PhaseType     string
	Provider      string
	Model         string
	ExecutionMode string
	BatchMode     string
}

// PersonaGenerationPhaseProgressSummary summarizes current progress for one phase run.
type PersonaGenerationPhaseProgressSummary struct {
	Percent        int
	ProcessedCount int
	TotalCount     int
	TargetCount    int
	CurrentStep    string
}

// PersonaGenerationTargetSummary summarizes target snapshot and skipped-target counts.
type PersonaGenerationTargetSummary struct {
	TargetCount            int
	CommonPersonaHitCount  int
	CommonPersonaMissCount int
	SkippedCount           int
	SkippedReasons         []string
	TargetSnapshotID       *string
	TargetSnapshotDigest   string
}

// PersonaGenerationExecutionSummary summarizes execution inputs without exposing secrets.
// JOB_PHASE_RUN が存在する場合だけ設定する。Ready 期は nil とする。
type PersonaGenerationExecutionSummary struct {
	CredentialRef string
	Provider      string
	Model         string
	ExecutionMode string
	PromptDigest  string
	InputCount    int
	OutputCount   int
	EvidenceRefs  []string
}

// PersonaGenerationPhaseAISettings holds the Ready 期 AI 選択値（JOB_PHASE_AI_SETTINGS から読む）。
// record 不在（AI 設定未保存）の場合は nil とする。
type PersonaGenerationPhaseAISettings struct {
	Provider      string
	Model         string
	ExecutionMode string
	BatchMode     string
}

// PersonaGenerationPhaseResultSummary summarizes the stored persona snapshot.
type PersonaGenerationPhaseResultSummary struct {
	GeneratedCount          int
	FailedCount             int
	PersonaCount            int
	MissingCount            int
	SnapshotID              string
	SnapshotDigest          string
	SnapshotReferenceStatus string
	BodyReadiness           bool
}

// PersonaGenerationPhaseErrorSummary exposes only redacted error information.
type PersonaGenerationPhaseErrorSummary struct {
	ErrorKind  PersonaGenerationPhaseErrorKind
	Reason     string
	Retryable  bool
	IsRedacted bool
}

// PersonaGenerationPhaseActionEnablement summarizes Job Run operation button state.
type PersonaGenerationPhaseActionEnablement struct {
	CanStart            bool
	StartBlockedReason  *string
	CanPause            bool
	PauseBlockedReason  *string
	CanResume           bool
	ResumeBlockedReason *string
	CanRetry            bool
	RetryBlockedReason  *string
	CanCancel           bool
	CancelBlockedReason *string
}

// PersonaGenerationPhaseSummaryResult is the frozen Job Run summary contract.
type PersonaGenerationPhaseSummaryResult struct {
	JobID         int64
	CurrentPhase  string
	PhaseState    string
	PhaseRunID    *int64
	StartedAt     *time.Time
	FinishedAt    *time.Time
	Progress      PersonaGenerationPhaseProgressSummary
	TargetSummary PersonaGenerationTargetSummary
	// AISettings は JOB_PHASE_AI_SETTINGS record が存在する場合だけ設定する。nil は未設定を意味する。
	AISettings *PersonaGenerationPhaseAISettings
	// Execution は JOB_PHASE_RUN が存在する場合だけ設定する。nil は未開始を意味する。
	Execution        *PersonaGenerationExecutionSummary
	ResultSummary    *PersonaGenerationPhaseResultSummary
	ErrorSummary     *PersonaGenerationPhaseErrorSummary
	ActionEnablement PersonaGenerationPhaseActionEnablement
}

// PersonaGenerationPhaseCommandResult is the frozen write-seam response contract.
type PersonaGenerationPhaseCommandResult struct {
	JobID         int64
	CurrentPhase  string
	PhaseState    string
	PhaseRunID    *int64
	StartedAt     *time.Time
	FinishedAt    *time.Time
	Progress      PersonaGenerationPhaseProgressSummary
	TargetSummary PersonaGenerationTargetSummary
	// Execution は JOB_PHASE_RUN が存在する場合だけ設定する。nil は未開始を意味する。
	Execution     *PersonaGenerationExecutionSummary
	ResultSummary *PersonaGenerationPhaseResultSummary
	Retryable     bool
	ErrorSummary  *PersonaGenerationPhaseErrorSummary
}

// PersonaGenerationBodyReadinessInputSummary summarizes the body phase inputs.
type PersonaGenerationBodyReadinessInputSummary struct {
	PersonaCount   int
	MissingCount   int
	SnapshotID     string
	SnapshotDigest string
	EvidenceRefs   []string
}

// PersonaGenerationBodyReadinessResult is the frozen downstream phase fact state contract.
type PersonaGenerationBodyReadinessResult struct {
	JobID         int64
	CurrentPhase  string
	PhaseState    string
	JobIsTerminal bool
	InputSummary  PersonaGenerationBodyReadinessInputSummary
}

// NewPersonaGenerationPhaseContractStub returns a temporary usecase stub for the frozen Wails seam.
func NewPersonaGenerationPhaseContractStub() PersonaGenerationPhaseContractStub {
	return PersonaGenerationPhaseContractStub{}
}

// PersonaGenerationPhaseContractStub is a temporary contract-only usecase used until the real phase usecase exists.
type PersonaGenerationPhaseContractStub struct{}

// GetPersonaGenerationPhaseSummary returns the frozen summary contract fixture.
func (PersonaGenerationPhaseContractStub) GetPersonaGenerationPhaseSummary(
	context.Context,
	GetPersonaGenerationPhaseSummaryRequest,
) (PersonaGenerationPhaseSummaryResult, error) {
	return personaGenerationPhaseSummaryFixture(), nil
}

// StartPersonaGenerationPhase returns one frozen start contract fixture.
func (PersonaGenerationPhaseContractStub) StartPersonaGenerationPhase(
	_ context.Context,
	request StartPersonaGenerationPhaseRequest,
) (PersonaGenerationPhaseCommandResult, error) {
	switch request.JobID {
	case 2010:
		result := personaGenerationStartedFixture(request.JobID, 0)
		result.PhaseRunID = nil
		result.PhaseState = "rejected"
		result.ErrorSummary = personaGenerationErrorSummary(PersonaGenerationPhaseErrorKindTerminalJob, false)
		return result, errPersonaGenerationPhaseTerminalJob
	default:
		phaseRunID := personaPhaseRunIDForJob(request.JobID)
		return personaGenerationStartedFixture(request.JobID, phaseRunID), nil
	}
}

// PausePersonaGenerationPhase returns one frozen pause contract fixture.
func (PersonaGenerationPhaseContractStub) PausePersonaGenerationPhase(
	_ context.Context,
	request PausePersonaGenerationPhaseRequest,
) (PersonaGenerationPhaseCommandResult, error) {
	result := personaGenerationStartedFixture(request.JobID, request.PhaseRunID)
	result.PhaseState = "paused"
	result.Progress.CurrentStep = "paused"
	return result, nil
}

// ResumePersonaGenerationPhase returns one frozen resume contract fixture.
func (PersonaGenerationPhaseContractStub) ResumePersonaGenerationPhase(
	_ context.Context,
	request ResumePersonaGenerationPhaseRequest,
) (PersonaGenerationPhaseCommandResult, error) {
	result := personaGenerationStartedFixture(request.JobID, request.PhaseRunID)
	result.PhaseState = "running"
	result.Progress.CurrentStep = "resumed"
	return result, nil
}

// RetryPersonaGenerationPhase returns one frozen retry contract fixture.
func (PersonaGenerationPhaseContractStub) RetryPersonaGenerationPhase(
	_ context.Context,
	request RetryPersonaGenerationPhaseRequest,
) (PersonaGenerationPhaseCommandResult, error) {
	switch request.JobID {
	case 2004:
		result := personaGenerationFailureFixture(request.JobID, request.PhaseRunID, PersonaGenerationPhaseErrorKindSaveFailed)
		return result, errPersonaGenerationPhaseSaveFailed
	case 2071:
		result := personaGenerationFailureFixture(request.JobID, request.PhaseRunID, PersonaGenerationPhaseErrorKindProviderFailure)
		return result, errPersonaGenerationPhaseProviderFailure
	case 2072:
		result := personaGenerationFailureFixture(request.JobID, request.PhaseRunID, PersonaGenerationPhaseErrorKindInvalidProviderResponse)
		return result, errPersonaGenerationPhaseInvalidProviderResponse
	case 2073:
		result := personaGenerationFailureFixture(request.JobID, request.PhaseRunID, PersonaGenerationPhaseErrorKindInputMissing)
		return result, errPersonaGenerationPhaseInputMissing
	case 2074:
		result := personaGenerationFailureFixture(request.JobID, request.PhaseRunID, PersonaGenerationPhaseErrorKindSaveFailed)
		return result, errPersonaGenerationPhaseSaveFailed
	default:
		result := personaGenerationStartedFixture(request.JobID, request.PhaseRunID)
		result.PhaseState = "completed"
		result.Progress = PersonaGenerationPhaseProgressSummary{
			Percent:        100,
			ProcessedCount: 5,
			TotalCount:     5,
			TargetCount:    5,
			CurrentStep:    "completed",
		}
		result.ResultSummary = &PersonaGenerationPhaseResultSummary{
			GeneratedCount:          4,
			FailedCount:             1,
			PersonaCount:            4,
			MissingCount:            1,
			SnapshotID:              "persona-snapshot-2008",
			SnapshotDigest:          "sha256:persona-snapshot-2008",
			SnapshotReferenceStatus: "available",
			BodyReadiness:           false,
		}
		result.ErrorSummary = nil
		return result, nil
	}
}

// CancelPersonaGenerationPhase returns one frozen cancel contract fixture.
func (PersonaGenerationPhaseContractStub) CancelPersonaGenerationPhase(
	_ context.Context,
	request CancelPersonaGenerationPhaseRequest,
) (PersonaGenerationPhaseCommandResult, error) {
	result := personaGenerationStartedFixture(request.JobID, request.PhaseRunID)
	result.PhaseState = "canceled"
	result.Progress.CurrentStep = "canceled"
	return result, nil
}

// GetPersonaGenerationBodyReadiness returns one frozen body readiness fact state fixture.
func (PersonaGenerationPhaseContractStub) GetPersonaGenerationBodyReadiness(
	_ context.Context,
	request GetPersonaGenerationBodyReadinessRequest,
) (PersonaGenerationBodyReadinessResult, error) {
	if request.JobID == 2010 {
		result := PersonaGenerationBodyReadinessResult{
			JobID:         request.JobID,
			CurrentPhase:  "persona_generation",
			PhaseState:    "rejected",
			JobIsTerminal: true,
		}
		return result, errPersonaGenerationPhaseTerminalJob
	}

	return PersonaGenerationBodyReadinessResult{
		JobID:        request.JobID,
		CurrentPhase: "persona_generation",
		PhaseState:   "completed",
		InputSummary: PersonaGenerationBodyReadinessInputSummary{
			PersonaCount:   4,
			MissingCount:   1,
			SnapshotID:     "persona-snapshot-2006",
			SnapshotDigest: "sha256:persona-snapshot-2006",
			EvidenceRefs:   []string{personaGenerationFixtureEvidenceNPC001},
		},
	}, nil
}

var (
	errPersonaGenerationPhaseTerminalJob             = errors.New("persona generation phase terminal job")
	errPersonaGenerationPhaseProviderFailure         = errors.New("persona generation provider failure")
	errPersonaGenerationPhaseInvalidProviderResponse = errors.New("persona generation invalid provider response")
	errPersonaGenerationPhaseInputMissing            = errors.New("persona generation input missing")
	errPersonaGenerationPhaseSaveFailed              = errors.New("persona generation save failed")
)

func personaGenerationPhaseSummaryFixture() PersonaGenerationPhaseSummaryResult {
	phaseRunID := int64(3001)
	startedAt := time.Date(2026, 5, 1, 12, 0, 0, 0, time.FixedZone("JST", 9*60*60))
	return PersonaGenerationPhaseSummaryResult{
		JobID:        101,
		CurrentPhase: "persona_generation",
		PhaseState:   "running",
		PhaseRunID:   &phaseRunID,
		StartedAt:    &startedAt,
		Progress: PersonaGenerationPhaseProgressSummary{
			Percent:        40,
			ProcessedCount: 2,
			TotalCount:     5,
			TargetCount:    5,
			CurrentStep:    "generating",
		},
		TargetSummary: PersonaGenerationTargetSummary{
			TargetCount:            5,
			CommonPersonaHitCount:  1,
			CommonPersonaMissCount: 4,
			SkippedCount:           1,
			SkippedReasons:         []string{"orphan_npc_reference"},
			TargetSnapshotDigest:   "sha256:target-snapshot",
		},
		Execution: &PersonaGenerationExecutionSummary{
			CredentialRef: "credential:persona:test",
			Provider:      "fake",
			Model:         "persona-model",
			ExecutionMode: "single",
			PromptDigest:  "sha256:prompt",
			InputCount:    5,
			OutputCount:   4,
			EvidenceRefs:  []string{personaGenerationFixtureEvidenceNPC001},
		},
		ResultSummary: &PersonaGenerationPhaseResultSummary{
			GeneratedCount:          4,
			FailedCount:             1,
			PersonaCount:            4,
			MissingCount:            1,
			SnapshotID:              "persona-snapshot-101",
			SnapshotDigest:          "sha256:persona-snapshot",
			SnapshotReferenceStatus: "available",
			BodyReadiness:           false,
		},
		ErrorSummary: personaGenerationErrorSummary(" PROVIDER_FAILURE ", true),
		ActionEnablement: PersonaGenerationPhaseActionEnablement{
			CanStart:  false,
			CanPause:  true,
			CanResume: false,
			CanRetry:  true,
			CanCancel: true,
		},
	}
}

func personaPhaseRunIDForJob(jobID int64) int64 {
	if jobID == 2008 {
		return 4008
	}
	return 4002
}

func personaGenerationStartedFixture(jobID, phaseRunID int64) PersonaGenerationPhaseCommandResult {
	startedAt := time.Date(2026, 5, 1, 15, 0, 0, 0, time.FixedZone("JST", 9*60*60))
	targetSnapshotID := fmt.Sprintf("target-snapshot-%d", jobID)
	result := PersonaGenerationPhaseCommandResult{
		JobID:        jobID,
		CurrentPhase: "persona_generation",
		PhaseState:   "running",
		StartedAt:    &startedAt,
		Progress: PersonaGenerationPhaseProgressSummary{
			Percent:        20,
			ProcessedCount: 1,
			TotalCount:     5,
			TargetCount:    5,
			CurrentStep:    "starting",
		},
		TargetSummary: PersonaGenerationTargetSummary{
			TargetCount:            5,
			CommonPersonaHitCount:  1,
			CommonPersonaMissCount: 4,
			SkippedCount:           1,
			SkippedReasons:         []string{"orphan_npc_reference"},
			TargetSnapshotID:       &targetSnapshotID,
			TargetSnapshotDigest:   fmt.Sprintf("sha256:target-snapshot-%d", jobID),
		},
		Execution: &PersonaGenerationExecutionSummary{
			CredentialRef: "credential:persona:test",
			Provider:      "fake",
			Model:         "persona-model",
			ExecutionMode: "single",
			PromptDigest:  fmt.Sprintf("sha256:prompt-%d", jobID),
			InputCount:    5,
			OutputCount:   0,
			EvidenceRefs:  []string{personaGenerationFixtureEvidenceNPC001},
		},
		ResultSummary: &PersonaGenerationPhaseResultSummary{
			GeneratedCount:          4,
			FailedCount:             1,
			PersonaCount:            4,
			MissingCount:            1,
			SnapshotReferenceStatus: "pending",
			BodyReadiness:           false,
		},
		Retryable: false,
	}
	if phaseRunID > 0 {
		result.PhaseRunID = &phaseRunID
	}
	return result
}

func personaGenerationFailureFixture(
	jobID, phaseRunID int64,
	errorKind PersonaGenerationPhaseErrorKind,
) PersonaGenerationPhaseCommandResult {
	result := personaGenerationStartedFixture(jobID, phaseRunID)
	result.PhaseState = "failed"
	result.Progress.CurrentStep = "failed"
	result.Retryable = true
	result.ResultSummary = &PersonaGenerationPhaseResultSummary{
		GeneratedCount:          0,
		FailedCount:             1,
		PersonaCount:            0,
		MissingCount:            1,
		SnapshotReferenceStatus: "missing",
		BodyReadiness:           false,
	}
	result.ErrorSummary = personaGenerationErrorSummary(errorKind, true)
	return result
}

func personaGenerationErrorSummary(
	errorKind PersonaGenerationPhaseErrorKind,
	retryable bool,
) *PersonaGenerationPhaseErrorSummary {
	return &PersonaGenerationPhaseErrorSummary{
		ErrorKind:  NormalizePersonaGenerationPhasePublicErrorKind(errorKind),
		Reason:     "redacted public summary",
		Retryable:  retryable,
		IsRedacted: true,
	}
}
