package service

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"log/slog"
	"sort"
	"strings"
	"time"

	"aitranslationenginejp/internal/repository"
)

const (
	termTranslationPhaseType                 = "term_translation"
	termTranslationCurrentPhase              = "term_translation"
	termTranslationInitialPhaseType          = "translation"
	termTranslationPhaseStateIdleReady       = "idle_ready"
	termTranslationPhaseStatePending         = "pending"
	termTranslationPhaseStateRunning         = "running"
	termTranslationPhaseStateCompleted       = "completed"
	termTranslationPhaseStatePaused          = "paused"
	termTranslationPhaseStateRecoverableFail = "recoverable_failed"
	termTranslationPhaseStateFailed          = "failed"
	termTranslationJobStateReady             = "ready"
	termTranslationJobStateRunning           = "running"
	termTranslationJobStatePaused            = "paused"
	termTranslationJobStateRecoverableFail   = "recoverable_failed"
	termTranslationJobStateCompleted         = "completed"
	termTranslationJobStateFailed            = "failed"
	termTranslationJobStateCanceled          = "canceled"
	termTranslationDictionaryLifecycleMaster = "master"
	termTranslationDictionaryLifecycleJob    = "job"
	termTranslationDictionarySourceSnapshot  = "shared_dictionary_snapshot"
	termTranslationDictionarySourceProvider  = "term_translation_provider"
	termTranslationLinkRoleSnapshot          = "snapshot_hit"
	termTranslationLinkRoleApplied           = "job_applied"
	termTranslationInstructionKindDefault    = "default"
	termTranslationErrorKindProviderFailure  = "provider_failure"
	termTranslationErrorKindInvalidProvider  = "invalid_provider_response"
	termTranslationErrorKindSaveFailed       = "save_failed"
	termTranslationReasonProviderFailed      = "provider execution failed"
	termTranslationReasonProviderInvalid     = "provider response was invalid"
	termTranslationReasonSaveFailed          = "save failed"
	termTranslationReasonTerminalJob         = "terminal job"
	termTranslationReasonReadyRequired       = "ready job is required"
	termTranslationReasonPhaseRuntimeMissing = "phase runtime snapshot is missing"
	termTranslationReasonCredentialMissing   = "term translation provider credential is unavailable"
	termTranslationReasonActivePhaseExists   = "active phase run already exists"
	termTranslationReasonPhaseIncomplete     = "term phase state does not allow this command"
	termTranslationPhaseServiceWhere         = "term_translation_phase.service"
	termTranslationSecretLoadTimeout         = 250 * time.Millisecond
)

var (
	// ErrTermTranslationProviderUnavailable reports that the provider port is not configured.
	ErrTermTranslationProviderUnavailable = errors.New("term translation provider is not configured")
	// ErrTermTranslationInvalidProviderReply reports that the provider returned an invalid response.
	ErrTermTranslationInvalidProviderReply   = errors.New("term translation provider response is invalid")
	errTermTranslationPhaseExecutionRejected = errors.New("term translation phase execution rejected")
)

// TermTranslationPhaseService orchestrates term phase execution state and job-local dictionary persistence.
type TermTranslationPhaseService struct {
	now                      func() time.Time
	jobLifecycleRepository   termTranslationPhaseJobLifecycleRepository
	foundationDataRepository termTranslationPhaseFoundationDataRepository
	translationSourceReader  termTranslationPhaseTranslationSourceRepository
	transactor               repository.Transactor
	secretStore              termTranslationPhaseSecretStore
	provider                 TermTranslationProvider
	providerSettings         ProviderSettingsConsumer
	executionSnapshots       map[int64]providerExecutionSnapshot
	secretLoadTimeout        time.Duration
}

type termTranslationPhaseJobLifecycleRepository interface {
	GetTranslationJobByID(ctx context.Context, id int64) (repository.TranslationJob, error)
	UpdateTranslationJob(ctx context.Context, id int64, draft repository.TranslationJobUpdateDraft) (repository.TranslationJob, error)
	CreateJobPhaseRun(ctx context.Context, draft repository.JobPhaseRunDraft) (repository.JobPhaseRun, error)
	FindJobPhaseRun(ctx context.Context, translationJobID int64, phaseType string) (repository.JobPhaseRun, error)
	UpdateJobPhaseRun(ctx context.Context, id int64, draft repository.JobPhaseRunUpdateDraft) (repository.JobPhaseRun, error)
	UpdateJobPhaseRunWhenState(ctx context.Context, id int64, expectedState string, draft repository.JobPhaseRunUpdateDraft) (repository.JobPhaseRun, error)
	ListJobPhaseRunsByJobID(ctx context.Context, jobID int64) ([]repository.JobPhaseRun, error)
	CreatePhaseRunDictionaryEntry(ctx context.Context, draft repository.PhaseRunDictionaryEntryDraft) (repository.PhaseRunDictionaryEntry, error)
	ListPhaseRunDictionaryEntriesByPhaseRunID(ctx context.Context, phaseRunID int64) ([]repository.PhaseRunDictionaryEntry, error)
}

type termTranslationPhaseFoundationDataRepository interface {
	GetDictionaryEntryByID(ctx context.Context, id int64) (repository.DictionaryEntry, error)
	CreateDictionaryEntry(ctx context.Context, draft repository.DictionaryEntryDraft) (repository.DictionaryEntry, error)
	UpdateDictionaryEntry(ctx context.Context, id int64, draft repository.DictionaryEntryUpdateDraft) (repository.DictionaryEntry, error)
	ListDictionaryEntries(ctx context.Context, translationJobID *int64, lifecycle string, scope string, sourceTerm string) ([]repository.DictionaryEntry, error)
	FindDictionaryEntry(ctx context.Context, translationJobID *int64, lifecycle string, scope string, sourceTerm string) (repository.DictionaryEntry, error)
}

type termTranslationPhaseTranslationSourceRepository interface {
	GetXEditExtractedDataByID(ctx context.Context, id int64) (repository.XEditExtractedData, error)
	ListTranslationRecordsByXEditID(ctx context.Context, xEditID int64) ([]repository.TranslationRecord, error)
	ListTranslationFieldsByTranslationRecordID(ctx context.Context, translationRecordID int64) ([]repository.TranslationField, error)
}

type termTranslationPhaseSecretStore interface {
	Load(ctx context.Context, key string) (string, error)
}

type termTranslationPhaseRuntimeSnapshotReader interface {
	GetTranslationJobPhaseRuntimeSnapshot(
		ctx context.Context,
		translationJobID int64,
		phaseID string,
	) (repository.TranslationJobPhaseRuntimeSnapshot, error)
}

// TermTranslationExecutionConfigReadModel summarizes execution config for one term phase run.
type TermTranslationExecutionConfigReadModel struct {
	CredentialRef   string
	CredentialState string
	EndpointSummary *string
	Provider        string
	Model           string
	ExecutionMode   string
	SnapshotDigest  *string
	SnapshotVersion *string
}

// TermTranslationPhaseProgressReadModel summarizes persisted term phase progress.
type TermTranslationPhaseProgressReadModel struct {
	Percent        int
	ProcessedCount int
	TotalCount     int
	AITargetCount  int
	CurrentStep    string
}

// TermTranslationPhaseResultReadModel summarizes persisted term phase result counts.
type TermTranslationPhaseResultReadModel struct {
	ConfirmedCount            int
	JobDictionaryAppliedCount int
	ReplacementTargetCount    int
	UnmatchedCount            int
	ProviderSkipped           bool
}

// TermTranslationPhaseErrorReadModel exposes redacted term phase error information.
type TermTranslationPhaseErrorReadModel struct {
	ErrorKind  string
	Reason     string
	Retryable  bool
	IsRedacted bool
}

// TermTranslationPhaseActionEnablementReadModel summarizes action availability for one phase state.
type TermTranslationPhaseActionEnablementReadModel struct {
	CanStart               bool
	StartBlockedReason     *string
	CanPause               bool
	PauseBlockedReason     *string
	CanResume              bool
	ResumeBlockedReason    *string
	CanRetry               bool
	RetryBlockedReason     *string
	CanRefresh             bool
	RefreshBlockedReason   *string
	CanStartNextPhase      bool
	NextPhaseBlockedReason *string
}

// TermTranslationPhaseSummaryReadModel stores summary payload data for the usecase boundary.
type TermTranslationPhaseSummaryReadModel struct {
	JobID              int64
	CurrentPhase       string
	PhaseState         string
	PhaseRunID         *int64
	StartedAt          *time.Time
	FinishedAt         *time.Time
	Progress           TermTranslationPhaseProgressReadModel
	TotalTermCount     int
	DictionaryHitCount int
	AITargetCount      int
	Execution          TermTranslationExecutionConfigReadModel
	ResultSummary      *TermTranslationPhaseResultReadModel
	ErrorSummary       *TermTranslationPhaseErrorReadModel
	ActionEnablement   TermTranslationPhaseActionEnablementReadModel
}

// TermTranslationPhaseCommandReadModel stores command payload data for the usecase boundary.
type TermTranslationPhaseCommandReadModel struct {
	JobID             int64
	CurrentPhase      string
	BeforePhaseState  string
	AfterPhaseState   string
	PhaseState        string
	PhaseRunID        *int64
	StartedAt         *time.Time
	FinishedAt        *time.Time
	Progress          TermTranslationPhaseProgressReadModel
	Retryable         bool
	CanStartNextPhase bool
	ErrorSummary      *TermTranslationPhaseErrorReadModel
}

// TermTranslationNextPhaseReadinessReadModel stores downstream readiness for the next phase boundary.
type TermTranslationNextPhaseReadinessReadModel struct {
	JobID             int64
	CurrentPhase      string
	PhaseState        string
	CanStartNextPhase bool
	BlockedReason     *string
	ErrorKind         string
}

type termTranslationCandidate struct {
	RecordType       string
	SourceTerm       string
	NormalizedSource string
	FieldCount       int
}

type termTranslationExecutionPlan struct {
	Job              repository.TranslationJob
	PhaseRun         repository.JobPhaseRun
	BeforePhaseState string
	RunCreated       bool
	Execution        TermTranslationExecutionConfigReadModel
	Candidates       []termTranslationCandidate
	SnapshotHits     map[string]repository.DictionaryEntry
	ConfirmedTerms   map[string]repository.DictionaryEntry
	Rejection        *TermTranslationPhaseErrorReadModel
	StartSnapshot    *providerExecutionSnapshot
}

type termTranslationAttemptFailure struct {
	summary         TermTranslationPhaseErrorReadModel
	progressPercent int
	cause           error
}

func (failure *termTranslationAttemptFailure) Error() string {
	if failure == nil {
		return ""
	}
	if failure.cause == nil {
		return strings.TrimSpace(failure.summary.Reason)
	}
	return failure.cause.Error()
}

func (failure *termTranslationAttemptFailure) Unwrap() error {
	if failure == nil {
		return nil
	}
	return failure.cause
}

type termTranslationPersistedFailure struct {
	job              repository.TranslationJob
	run              repository.JobPhaseRun
	beforePhaseState string
	runCreated       bool
	snapshots        map[string]repository.DictionaryEntry
	failure          *termTranslationAttemptFailure
}

// NewTermTranslationPhaseService creates the backend service for the term translation phase.
func NewTermTranslationPhaseService(
	jobLifecycleRepository termTranslationPhaseJobLifecycleRepository,
	foundationDataRepository termTranslationPhaseFoundationDataRepository,
	translationSourceReader termTranslationPhaseTranslationSourceRepository,
	transactor repository.Transactor,
	secretStore termTranslationPhaseSecretStore,
	provider TermTranslationProvider,
) *TermTranslationPhaseService {
	return &TermTranslationPhaseService{
		now:                      func() time.Time { return time.Now().UTC() },
		jobLifecycleRepository:   jobLifecycleRepository,
		foundationDataRepository: foundationDataRepository,
		translationSourceReader:  translationSourceReader,
		transactor:               transactor,
		secretStore:              secretStore,
		provider:                 provider,
		executionSnapshots:       map[int64]providerExecutionSnapshot{},
		secretLoadTimeout:        termTranslationSecretLoadTimeout,
	}
}

// WithTermTranslationProviderSettings injects the provider settings resolver used at phase start.
func (service *TermTranslationPhaseService) WithTermTranslationProviderSettings(
	consumer ProviderSettingsConsumer,
) *TermTranslationPhaseService {
	if service != nil {
		service.providerSettings = consumer
	}
	return service
}

// ReadSummary loads the current term phase summary for one translation job.
func (service *TermTranslationPhaseService) ReadSummary(
	ctx context.Context,
	jobID int64,
) (TermTranslationPhaseSummaryReadModel, error) {
	job, run, execution, candidates, err := service.loadExecutionContext(termTranslationRequestContext(ctx), jobID)
	if err != nil {
		return TermTranslationPhaseSummaryReadModel{}, err
	}
	_, snapshotHits, resultSummary, progress, errorSummary, digest, version, err := service.buildRunState(
		termTranslationRequestContext(ctx),
		job.ID,
		run,
		candidates,
	)
	if err != nil {
		return TermTranslationPhaseSummaryReadModel{}, err
	}
	execution.SnapshotDigest = digest
	execution.SnapshotVersion = version
	state := termTranslationPhaseStateIdleReady
	var phaseRunID *int64
	var startedAt *time.Time
	var finishedAt *time.Time
	if run != nil {
		state = run.State
		phaseRunID = &run.ID
		startedAt = run.StartedAt
		finishedAt = run.FinishedAt
	}
	if state == "" {
		state = termTranslationPhaseStateIdleReady
	}
	total := len(candidates)
	hitCount := len(snapshotHits)
	aiTargetCount := total - hitCount
	if aiTargetCount < 0 {
		aiTargetCount = 0
	}
	confirmedCount := 0
	if resultSummary != nil {
		confirmedCount = resultSummary.ConfirmedCount
	}
	readiness := service.readinessFromState(job, run, total, confirmedCount, errorSummary)
	canStart := job.State == termTranslationJobStateReady && !termTranslationPhaseIsActive(run) && termTranslationExecutionConfigured(execution)
	return TermTranslationPhaseSummaryReadModel{
		JobID:              job.ID,
		CurrentPhase:       termTranslationCurrentPhase,
		PhaseState:         state,
		PhaseRunID:         cloneInt64Pointer(phaseRunID),
		StartedAt:          cloneTimePointer(startedAt),
		FinishedAt:         cloneTimePointer(finishedAt),
		Progress:           progress,
		TotalTermCount:     total,
		DictionaryHitCount: hitCount,
		AITargetCount:      aiTargetCount,
		Execution:          execution,
		ResultSummary:      resultSummary,
		ErrorSummary:       errorSummary,
		ActionEnablement: TermTranslationPhaseActionEnablementReadModel{
			CanStart:               canStart,
			StartBlockedReason:     termTranslationStartBlockedReason(job, run, execution),
			CanPause:               run != nil && run.State == termTranslationPhaseStateRunning,
			PauseBlockedReason:     termTranslationPauseBlockedReason(run),
			CanResume:              run != nil && run.State == termTranslationPhaseStatePaused,
			ResumeBlockedReason:    termTranslationResumeBlockedReason(run),
			CanRetry:               run != nil && run.State == termTranslationPhaseStateRecoverableFail,
			RetryBlockedReason:     termTranslationRetryBlockedReason(run),
			CanRefresh:             true,
			CanStartNextPhase:      readiness.CanStartNextPhase,
			NextPhaseBlockedReason: cloneStringPointer(readiness.BlockedReason),
		},
	}, nil
}

// StartPhase starts one term phase run for a ready job.
func (service *TermTranslationPhaseService) StartPhase(
	ctx context.Context,
	jobID int64,
) (TermTranslationPhaseCommandReadModel, error) {
	return service.executePhase(termTranslationRequestContext(ctx), jobID, termTranslationStartModeStart, 0)
}

// PausePhase pauses one active term phase run.
func (service *TermTranslationPhaseService) PausePhase(
	ctx context.Context,
	jobID int64,
	phaseRunID int64,
) (TermTranslationPhaseCommandReadModel, error) {
	if service.transactor == nil {
		return TermTranslationPhaseCommandReadModel{}, fmt.Errorf("pause term translation phase: transactor is not configured")
	}
	var response TermTranslationPhaseCommandReadModel
	err := service.transactor.WithTransaction(termTranslationRequestContext(ctx), func(txCtx context.Context) error {
		job, run, err := service.loadExistingRunForMutation(txCtx, jobID, phaseRunID)
		if err != nil {
			return err
		}
		if run.State != termTranslationPhaseStateRunning {
			response = termTranslationCommandFromRun(job, run, true, false, &TermTranslationPhaseErrorReadModel{
				ErrorKind:  "term_phase_incomplete",
				Reason:     "phase is not running",
				Retryable:  false,
				IsRedacted: false,
			})
			response.BeforePhaseState = run.State
			return nil
		}
		now := service.now()
		updatedRun, err := service.jobLifecycleRepository.UpdateJobPhaseRun(txCtx, run.ID, repository.JobPhaseRunUpdateDraft{
			State:               termTranslationPhaseStatePaused,
			ProgressPercent:     run.ProgressPercent,
			LatestExternalRunID: run.LatestExternalRunID,
			LatestError:         run.LatestError,
			StartedAt:           run.StartedAt,
			FinishedAt:          nil,
		})
		if err != nil {
			return fmt.Errorf("pause term translation phase run: %w", err)
		}
		updatedJob, err := service.jobLifecycleRepository.UpdateTranslationJob(txCtx, job.ID, repository.TranslationJobUpdateDraft{
			JobName:         job.JobName,
			State:           termTranslationJobStateRunning,
			ProgressPercent: updatedRun.ProgressPercent,
			StartedAt:       termTranslationJobStartedAt(job.StartedAt, now),
			FinishedAt:      nil,
		})
		if err != nil {
			return fmt.Errorf("pause translation job: %w", err)
		}
		response = termTranslationCommandFromRun(updatedJob, updatedRun, false, false, nil)
		response.BeforePhaseState = run.State
		return nil
	})
	if err != nil {
		return TermTranslationPhaseCommandReadModel{}, fmt.Errorf("pause term translation phase transaction: %w", err)
	}
	return response, nil
}

// ResumePhase resumes one paused term phase run.
func (service *TermTranslationPhaseService) ResumePhase(
	ctx context.Context,
	jobID int64,
	phaseRunID int64,
) (TermTranslationPhaseCommandReadModel, error) {
	return service.executePhase(termTranslationRequestContext(ctx), jobID, termTranslationStartModeResume, phaseRunID)
}

// RetryPhase retries one recoverable-failed term phase run.
func (service *TermTranslationPhaseService) RetryPhase(
	ctx context.Context,
	jobID int64,
	phaseRunID int64,
) (TermTranslationPhaseCommandReadModel, error) {
	return service.executePhase(termTranslationRequestContext(ctx), jobID, termTranslationStartModeRetry, phaseRunID)
}

// ReadNextPhaseReadiness loads whether the next phase may start.
func (service *TermTranslationPhaseService) ReadNextPhaseReadiness(
	ctx context.Context,
	jobID int64,
) (TermTranslationNextPhaseReadinessReadModel, error) {
	job, run, _, candidates, err := service.loadExecutionContext(termTranslationRequestContext(ctx), jobID)
	if err != nil {
		return TermTranslationNextPhaseReadinessReadModel{}, err
	}
	confirmedTerms, _, _, _, errorSummary, _, _, err := service.buildRunState(termTranslationRequestContext(ctx), job.ID, run, candidates)
	if err != nil {
		return TermTranslationNextPhaseReadinessReadModel{}, err
	}
	confirmedCount, _ := summarizeConfirmedCandidates(candidates, confirmedTerms)
	readiness := service.readinessFromState(job, run, len(candidates), confirmedCount, errorSummary)
	return readiness, nil
}

type termTranslationStartMode string

const (
	termTranslationStartModeStart  termTranslationStartMode = "start"
	termTranslationStartModeResume termTranslationStartMode = "resume"
	termTranslationStartModeRetry  termTranslationStartMode = "retry"
)

func (service *TermTranslationPhaseService) executePhase(
	ctx context.Context,
	jobID int64,
	mode termTranslationStartMode,
	phaseRunID int64,
) (TermTranslationPhaseCommandReadModel, error) {
	if service.transactor == nil {
		return TermTranslationPhaseCommandReadModel{}, fmt.Errorf("execute term translation phase: transactor is not configured")
	}
	var response TermTranslationPhaseCommandReadModel
	var persistedFailure *termTranslationPersistedFailure
	err := service.transactor.WithTransaction(ctx, func(txCtx context.Context) error {
		return service.executePhaseTransaction(txCtx, jobID, mode, phaseRunID, &response, &persistedFailure)
	})
	if persistedFailure != nil {
		failedRun, persistErr := service.persistAttemptFailure(ctx, *persistedFailure)
		if persistErr != nil {
			return TermTranslationPhaseCommandReadModel{}, fmt.Errorf("execute term translation phase transaction: %w", persistErr)
		}
		failureResponse := termTranslationCommandFromRun(
			persistedFailure.job,
			failedRun,
			persistedFailure.failure.summary.Retryable,
			false,
			&persistedFailure.failure.summary,
		)
		failureResponse.BeforePhaseState = persistedFailure.beforePhaseState
		return failureResponse, nil
	}
	if err != nil {
		return TermTranslationPhaseCommandReadModel{}, fmt.Errorf("execute term translation phase transaction: %w", err)
	}
	return response, nil
}

func (service *TermTranslationPhaseService) executePhaseTransaction(
	ctx context.Context,
	jobID int64,
	mode termTranslationStartMode,
	phaseRunID int64,
	response *TermTranslationPhaseCommandReadModel,
	persistedFailure **termTranslationPersistedFailure,
) error {
	plan, err := service.prepareExecutionPlan(ctx, jobID, mode, phaseRunID)
	if err != nil {
		if errors.Is(err, errTermTranslationPhaseExecutionRejected) {
			*response = termTranslationRejectedCommand(jobID, plan)
			return nil
		}
		return err
	}
	updatedRun, confirmedCount, errSummary, err := service.applyExecutionPlan(ctx, mode, plan)
	if err != nil {
		service.captureAttemptFailure(plan, err, persistedFailure)
		return err
	}
	*response = termTranslationCommandFromRun(
		plan.Job,
		updatedRun,
		errSummary != nil && errSummary.Retryable,
		confirmedCount == len(plan.Candidates),
		errSummary,
	)
	response.BeforePhaseState = plan.BeforePhaseState
	return nil
}

func (service *TermTranslationPhaseService) captureAttemptFailure(
	plan termTranslationExecutionPlan,
	err error,
	persistedFailure **termTranslationPersistedFailure,
) {
	var attemptFailure *termTranslationAttemptFailure
	if !errors.As(err, &attemptFailure) {
		return
	}
	*persistedFailure = &termTranslationPersistedFailure{
		job:              plan.Job,
		run:              plan.PhaseRun,
		beforePhaseState: plan.BeforePhaseState,
		runCreated:       plan.RunCreated,
		snapshots:        cloneTermTranslationDictionaryEntries(plan.SnapshotHits),
		failure:          attemptFailure,
	}
}

func (service *TermTranslationPhaseService) prepareExecutionPlan(
	ctx context.Context,
	jobID int64,
	mode termTranslationStartMode,
	phaseRunID int64,
) (termTranslationExecutionPlan, error) {
	job, run, execution, candidates, err := service.loadExecutionContext(ctx, jobID)
	if err != nil {
		return termTranslationExecutionPlan{}, err
	}
	beforePhaseState := termTranslationPhaseStateFromRun(run)
	if termTranslationJobIsTerminal(job.State) {
		return termTranslationExecutionPlan{
			Job:              job,
			PhaseRun:         valueOrEmptyRun(run),
			BeforePhaseState: beforePhaseState,
			Rejection: &TermTranslationPhaseErrorReadModel{
				ErrorKind:  "terminal_job",
				Reason:     termTranslationReasonTerminalJob,
				Retryable:  false,
				IsRedacted: false,
			},
		}, errTermTranslationPhaseExecutionRejected
	}
	if rejectedPlan, rejected := service.rejectExecutionPlan(job, run, mode, phaseRunID); rejected {
		return rejectedPlan, errTermTranslationPhaseExecutionRejected
	}
	execution, rejectedPlan, err := service.prepareExecutionPlanStartSnapshot(ctx, job, run, mode, execution)
	if err != nil {
		return termTranslationExecutionPlan{}, err
	}
	if rejectedPlan != nil {
		return *rejectedPlan, errTermTranslationPhaseExecutionRejected
	}
	if !termTranslationExecutionConfigured(execution) {
		return service.buildRejectedExecutionPlan(job, run, "term_phase_incomplete", termTranslationReasonPhaseRuntimeMissing), errTermTranslationPhaseExecutionRejected
	}
	job, run, execution, createdRun, err := service.ensureExecutionPlanRun(ctx, job, run, execution)
	if err != nil {
		return termTranslationExecutionPlan{}, err
	}
	if run == nil {
		return termTranslationExecutionPlan{Job: job, BeforePhaseState: beforePhaseState}, errTermTranslationPhaseExecutionRejected
	}
	snapshotHits, err := service.termTranslationSnapshotHitsForPlan(ctx, createdRun, run, candidates)
	if err != nil {
		return termTranslationExecutionPlan{
			Job:              job,
			PhaseRun:         *run,
			BeforePhaseState: beforePhaseState,
		}, err
	}
	confirmedEntries, err := service.loadJobDictionaryEntries(ctx, job.ID)
	if err != nil {
		return termTranslationExecutionPlan{}, err
	}
	return termTranslationExecutionPlan{
		Job:              job,
		PhaseRun:         *run,
		BeforePhaseState: beforePhaseState,
		RunCreated:       createdRun,
		Execution:        execution,
		Candidates:       candidates,
		SnapshotHits:     snapshotHits,
		ConfirmedTerms:   confirmedEntries,
		StartSnapshot:    termTranslationStartSnapshotForMode(mode, execution),
	}, nil
}

func termTranslationStartSnapshotForMode(
	mode termTranslationStartMode,
	execution TermTranslationExecutionConfigReadModel,
) *providerExecutionSnapshot {
	if mode != termTranslationStartModeRetry {
		return nil
	}
	snapshot := termTranslationProviderExecutionSnapshot(execution)
	return &snapshot
}

func (service *TermTranslationPhaseService) prepareExecutionPlanStartSnapshot(
	ctx context.Context,
	job repository.TranslationJob,
	run *repository.JobPhaseRun,
	mode termTranslationStartMode,
	execution TermTranslationExecutionConfigReadModel,
) (TermTranslationExecutionConfigReadModel, *termTranslationExecutionPlan, error) {
	if mode != termTranslationStartModeStart && mode != termTranslationStartModeRetry {
		return execution, nil, nil
	}
	resolvedExecution, rejection, err := service.resolveExecutionSnapshotForStart(ctx, execution)
	if err != nil {
		return TermTranslationExecutionConfigReadModel{}, nil, err
	}
	if rejection != nil {
		rejectedPlan := service.buildRejectedExecutionPlan(job, run, rejection.ErrorKind, rejection.Reason)
		return resolvedExecution, &rejectedPlan, nil
	}
	return resolvedExecution, nil, nil
}

func (service *TermTranslationPhaseService) ensureExecutionPlanRun(
	ctx context.Context,
	job repository.TranslationJob,
	run *repository.JobPhaseRun,
	execution TermTranslationExecutionConfigReadModel,
) (repository.TranslationJob, *repository.JobPhaseRun, TermTranslationExecutionConfigReadModel, bool, error) {
	if run != nil {
		if strings.TrimSpace(run.State) != termTranslationPhaseStatePending {
			return job, run, execution, false, nil
		}
		updatedJob, err := service.markExecutionPlanJobRunning(ctx, job)
		if err != nil {
			return repository.TranslationJob{}, nil, TermTranslationExecutionConfigReadModel{}, false, err
		}
		updatedRun, err := service.jobLifecycleRepository.UpdateJobPhaseRunWhenState(ctx, run.ID, termTranslationPhaseStatePending, repository.JobPhaseRunUpdateDraft{
			State:           termTranslationPhaseStateRunning,
			ProgressPercent: 0,
			AIProvider:      execution.Provider,
			ModelName:       execution.Model,
			ExecutionMode:   execution.ExecutionMode,
			CredentialRef:   execution.CredentialRef,
			InstructionKind: termTranslationInstructionKindDefault,
			StartedAt:       termTranslationRunStartedAt(run.StartedAt, service.now()),
			FinishedAt:      nil,
		})
		if err != nil {
			return repository.TranslationJob{}, nil, TermTranslationExecutionConfigReadModel{}, false, fmt.Errorf("update term translation phase run for start: %w", err)
		}
		execution.Provider = updatedRun.AIProvider
		execution.Model = updatedRun.ModelName
		execution.ExecutionMode = updatedRun.ExecutionMode
		startSnapshot := termTranslationProviderExecutionSnapshot(execution)
		if err := service.persistRuntimeSnapshot(ctx, updatedJob.ID, "word_translation", startSnapshot); err != nil {
			return repository.TranslationJob{}, nil, TermTranslationExecutionConfigReadModel{}, false, err
		}
		service.executionSnapshots[updatedRun.ID] = startSnapshot
		return updatedJob, &updatedRun, execution, true, nil
	}
	return repository.TranslationJob{}, nil, TermTranslationExecutionConfigReadModel{}, false, fmt.Errorf("find term translation phase run for start: %w", repository.ErrNotFound)
}

func (service *TermTranslationPhaseService) termTranslationSnapshotHitsForPlan(
	ctx context.Context,
	createdRun bool,
	run *repository.JobPhaseRun,
	candidates []termTranslationCandidate,
) (map[string]repository.DictionaryEntry, error) {
	if createdRun {
		snapshotHits, err := service.loadSnapshotHits(ctx, candidates)
		if err != nil {
			return nil, fmt.Errorf("load dictionary snapshot: %w", err)
		}
		return snapshotHits, nil
	}
	snapshotHits, _, _, err := service.loadRunSnapshotState(ctx, run)
	return snapshotHits, err
}

func (service *TermTranslationPhaseService) rejectExecutionPlan(
	job repository.TranslationJob,
	run *repository.JobPhaseRun,
	mode termTranslationStartMode,
	phaseRunID int64,
) (termTranslationExecutionPlan, bool) {
	switch mode {
	case termTranslationStartModeStart:
		return service.rejectStartExecutionPlan(job, run)
	case termTranslationStartModeResume:
		return service.rejectResumeExecutionPlan(job, run, phaseRunID)
	case termTranslationStartModeRetry:
		return service.rejectRetryExecutionPlan(job, run, phaseRunID)
	default:
		return service.buildRejectedExecutionPlan(job, run, "term_phase_incomplete", termTranslationReasonPhaseIncomplete), true
	}
}

func (service *TermTranslationPhaseService) rejectStartExecutionPlan(
	job repository.TranslationJob,
	run *repository.JobPhaseRun,
) (termTranslationExecutionPlan, bool) {
	if job.State != termTranslationJobStateReady {
		return service.buildRejectedExecutionPlan(job, run, "ready_required", termTranslationReasonReadyRequired), true
	}
	if run != nil && termTranslationPhaseIsActive(run) {
		return service.buildRejectedExecutionPlan(job, run, "active_phase_exists", termTranslationReasonActivePhaseExists), true
	}
	return termTranslationExecutionPlan{}, false
}

func (service *TermTranslationPhaseService) rejectResumeExecutionPlan(
	job repository.TranslationJob,
	run *repository.JobPhaseRun,
	phaseRunID int64,
) (termTranslationExecutionPlan, bool) {
	if run == nil || run.ID != phaseRunID {
		return service.buildRejectedExecutionPlan(job, run, "term_phase_incomplete", termTranslationReasonPhaseIncomplete), true
	}
	if run.State != termTranslationPhaseStatePaused && run.State != termTranslationPhaseStateRecoverableFail {
		return service.buildRejectedExecutionPlan(job, run, "term_phase_incomplete", "phase is not resumable"), true
	}
	return termTranslationExecutionPlan{}, false
}

func (service *TermTranslationPhaseService) rejectRetryExecutionPlan(
	job repository.TranslationJob,
	run *repository.JobPhaseRun,
	phaseRunID int64,
) (termTranslationExecutionPlan, bool) {
	if run == nil || run.ID != phaseRunID {
		return service.buildRejectedExecutionPlan(job, run, "term_phase_incomplete", termTranslationReasonPhaseIncomplete), true
	}
	if run.State != termTranslationPhaseStateRecoverableFail {
		return service.buildRejectedExecutionPlan(job, run, "term_phase_incomplete", "phase is not retryable"), true
	}
	return termTranslationExecutionPlan{}, false
}

func (service *TermTranslationPhaseService) buildRejectedExecutionPlan(
	job repository.TranslationJob,
	run *repository.JobPhaseRun,
	errorKind string,
	reason string,
) termTranslationExecutionPlan {
	return termTranslationExecutionPlan{
		Job:              job,
		PhaseRun:         valueOrEmptyRun(run),
		BeforePhaseState: termTranslationPhaseStateFromRun(run),
		Rejection: &TermTranslationPhaseErrorReadModel{
			ErrorKind:  errorKind,
			Reason:     reason,
			Retryable:  false,
			IsRedacted: false,
		},
	}
}

func (service *TermTranslationPhaseService) resolveExecutionSnapshotForStart(
	ctx context.Context,
	execution TermTranslationExecutionConfigReadModel,
) (TermTranslationExecutionConfigReadModel, *TermTranslationPhaseErrorReadModel, error) {
	if !providerExecutionUsesProviderSettings(execution.Provider) {
		logProviderBoundarySkipped(ctx, "term_translation_provider_settings", termTranslationPhaseServiceWhere, execution.Provider, "provider_skipped", 1)
		return execution, nil, nil
	}
	if service.providerSettings == nil {
		return execution, &TermTranslationPhaseErrorReadModel{
			ErrorKind:  "term_phase_incomplete",
			Reason:     termTranslationReasonPhaseRuntimeMissing,
			Retryable:  false,
			IsRedacted: true,
		}, nil
	}
	resolved, err := service.providerSettings.ResolveProviderExecutionSettings(ctx, ProviderSettingsResolveInput{
		ConsumerID:          "term_translation_phase",
		AllowSecretSnapshot: true,
		Selection: ProviderSettingsResolveSelection{
			ProviderID:      execution.Provider,
			Model:           execution.Model,
			ExecutionMethod: execution.ExecutionMode,
			UseBatchAPI:     false,
		},
	})
	if err != nil {
		logProviderBoundaryFailure(ctx, "term_translation_provider_settings", termTranslationPhaseServiceWhere, execution.Provider, "secret_store_failed")
		return TermTranslationExecutionConfigReadModel{}, nil, fmt.Errorf("resolve term translation provider settings: %w", err)
	}
	execution.CredentialRef = providerExecutionOptionalString(resolved.CredentialReferenceID)
	execution.CredentialState = strings.TrimSpace(resolved.CredentialState)
	execution.EndpointSummary = providerExecutionEndpointSummary(resolved.Endpoint)
	if resolved.ErrorKind == nil {
		return execution, nil, nil
	}
	logProviderBoundaryFailure(ctx, "term_translation_provider_settings", termTranslationPhaseServiceWhere, execution.Provider, *resolved.ErrorKind)
	switch strings.TrimSpace(*resolved.ErrorKind) {
	case providerSettingsErrorKindCredentialMissing:
		return execution, &TermTranslationPhaseErrorReadModel{
			ErrorKind:  "term_phase_incomplete",
			Reason:     termTranslationReasonCredentialMissing,
			Retryable:  false,
			IsRedacted: true,
		}, nil
	case providerSettingsErrorKindEndpointMissing:
		return execution, &TermTranslationPhaseErrorReadModel{
			ErrorKind:  "term_phase_incomplete",
			Reason:     termTranslationReasonPhaseRuntimeMissing,
			Retryable:  false,
			IsRedacted: true,
		}, nil
	default:
		return execution, &TermTranslationPhaseErrorReadModel{
			ErrorKind:  "term_phase_incomplete",
			Reason:     termTranslationReasonProviderFailed,
			Retryable:  false,
			IsRedacted: true,
		}, nil
	}
}

func termTranslationProviderExecutionSnapshot(
	execution TermTranslationExecutionConfigReadModel,
) providerExecutionSnapshot {
	return providerExecutionSnapshot{
		Provider:        execution.Provider,
		Model:           execution.Model,
		ExecutionMode:   execution.ExecutionMode,
		CredentialRef:   execution.CredentialRef,
		CredentialState: execution.CredentialState,
		EndpointSummary: providerExecutionEndpointSummary(execution.EndpointSummary),
	}
}

func (service *TermTranslationPhaseService) persistRuntimeSnapshot(
	ctx context.Context,
	jobID int64,
	phaseID string,
	snapshot providerExecutionSnapshot,
) error {
	store, ok := service.jobLifecycleRepository.(translationJobPhaseRuntimeSnapshotStore)
	if !ok {
		return nil
	}
	current, err := store.GetTranslationJobPhaseRuntimeSnapshot(ctx, jobID, phaseID)
	switch {
	case err == nil:
	case errors.Is(err, repository.ErrNotFound):
		current = repository.TranslationJobPhaseRuntimeSnapshot{}
	default:
		return fmt.Errorf("persist term translation runtime snapshot: %w", err)
	}
	if _, err := store.SaveTranslationJobPhaseRuntimeSnapshot(ctx, repository.TranslationJobPhaseRuntimeSnapshotDraft{
		TranslationJobID: jobID,
		PhaseID:          phaseID,
		Provider:         snapshot.Provider,
		ModelName:        snapshot.Model,
		CredentialStatus: snapshot.CredentialState,
		ExecutionMode:    snapshot.ExecutionMode,
		BatchMode:        providerExecutionBatchMode(current),
	}); err != nil {
		return fmt.Errorf("persist term translation runtime snapshot: %w", err)
	}
	return nil
}

func (service *TermTranslationPhaseService) markExecutionPlanJobRunning(
	ctx context.Context,
	job repository.TranslationJob,
) (repository.TranslationJob, error) {
	now := service.now()
	startedAt := cloneTimePointer(job.StartedAt)
	if startedAt == nil {
		startedAt = &now
	}
	updatedJob, err := service.jobLifecycleRepository.UpdateTranslationJob(ctx, job.ID, repository.TranslationJobUpdateDraft{
		JobName:         job.JobName,
		State:           termTranslationJobStateRunning,
		ProgressPercent: 0,
		StartedAt:       startedAt,
		FinishedAt:      nil,
	})
	if err != nil {
		return repository.TranslationJob{}, fmt.Errorf("update translation job running state: %w", err)
	}
	return updatedJob, nil
}

func (service *TermTranslationPhaseService) applyExecutionPlan(
	ctx context.Context,
	mode termTranslationStartMode,
	plan termTranslationExecutionPlan,
) (repository.JobPhaseRun, int, *TermTranslationPhaseErrorReadModel, error) {
	confirmedTerms := make(map[string]repository.DictionaryEntry, len(plan.ConfirmedTerms))
	for key, entry := range plan.ConfirmedTerms {
		confirmedTerms[key] = entry
	}
	baseConfirmedCount, _ := summarizeConfirmedCandidates(plan.Candidates, plan.ConfirmedTerms)
	if err := service.applySnapshotHits(ctx, plan, confirmedTerms, baseConfirmedCount); err != nil {
		return repository.JobPhaseRun{}, 0, nil, err
	}

	aiTargets := service.buildAITargets(plan.Candidates, confirmedTerms)
	if len(aiTargets) == 0 {
		return service.markPhaseRunCompleted(ctx, plan.Job, plan.PhaseRun, len(plan.Candidates))
	}

	updatedRun, err := service.resumeExecutionPlanRun(ctx, mode, plan.PhaseRun, len(confirmedTerms), len(plan.Candidates))
	if err != nil {
		return repository.JobPhaseRun{}, 0, nil, err
	}
	if plan.StartSnapshot != nil && strings.TrimSpace(updatedRun.State) == termTranslationPhaseStateRunning {
		if persistErr := service.persistRuntimeSnapshot(ctx, plan.Job.ID, "word_translation", *plan.StartSnapshot); persistErr != nil {
			return repository.JobPhaseRun{}, 0, nil, persistErr
		}
		service.executionSnapshots[updatedRun.ID] = *plan.StartSnapshot
	}
	plan.PhaseRun = updatedRun
	if service.provider == nil {
		logProviderBoundaryFailure(ctx, "term_translation_provider_execute", termTranslationPhaseServiceWhere, plan.Execution.Provider, "provider_unavailable")
		return repository.JobPhaseRun{}, 0, nil, newTermTranslationAttemptFailure(
			termTranslationErrorKindProviderFailure,
			termTranslationReasonProviderFailed,
			true,
			true,
			termTranslationProgressPercent(baseConfirmedCount, len(plan.Candidates)),
			ErrTermTranslationProviderUnavailable,
		)
	}
	apiKey, err := service.resolveProviderAPIKey(ctx, plan.Execution)
	if err != nil {
		logProviderBoundaryFailure(ctx, "term_translation_provider_secret", termTranslationPhaseServiceWhere, plan.Execution.Provider, classifyTermTranslationCredentialFailure(err))
		return repository.JobPhaseRun{}, 0, nil, newTermTranslationAttemptFailure(
			termTranslationErrorKindProviderFailure,
			redactedTermTranslationProviderFailureReason(err),
			false,
			true,
			termTranslationProgressPercent(baseConfirmedCount, len(plan.Candidates)),
			err,
		)
	}
	for _, candidate := range aiTargets {
		if err := service.applyProviderCandidate(ctx, plan, candidate, apiKey, confirmedTerms, baseConfirmedCount); err != nil {
			logTermTranslationProviderBulkSummary(ctx, len(plan.Candidates), len(confirmedTerms), 1, classifyTermTranslationBulkFailure(err))
			return repository.JobPhaseRun{}, 0, nil, err
		}
	}
	logTermTranslationProviderBulkSummary(ctx, len(plan.Candidates), len(confirmedTerms), 0, "")
	return service.markPhaseRunCompleted(ctx, plan.Job, plan.PhaseRun, len(plan.Candidates))
}

func logTermTranslationProviderBulkSummary(
	ctx context.Context,
	inputCount int,
	outputCount int,
	failedCount int,
	failureKind string,
) {
	attrs := []slog.Attr{
		slog.String("event", "term_translation_provider_bulk_summary"),
		slog.String("where", "backend.service.term_translation_phase.provider_execution"),
		slog.String("result", "completed"),
		slog.Int("input_count", inputCount),
		slog.Int("output_count", outputCount),
		slog.Int("skipped_count", maxInt(0, inputCount-outputCount-failedCount)),
		slog.Int("failed_count", failedCount),
	}
	if failureKind != "" {
		attrs = append(attrs,
			slog.String("first_failure_kind", failureKind),
			slog.String("last_failure_kind", failureKind),
		)
	}
	slog.LogAttrs(ctx, slog.LevelInfo, "term translation provider bulk summary", attrs...)
}

func classifyTermTranslationBulkFailure(err error) string {
	var attemptFailure *termTranslationAttemptFailure
	if errors.As(err, &attemptFailure) {
		return normalizeTermTranslationText(attemptFailure.summary.ErrorKind)
	}
	return "provider_execution_failed"
}

func (service *TermTranslationPhaseService) applySnapshotHits(
	ctx context.Context,
	plan termTranslationExecutionPlan,
	confirmedTerms map[string]repository.DictionaryEntry,
	baseConfirmedCount int,
) error {
	for key, entry := range plan.SnapshotHits {
		jobEntry, err := service.ensureJobDictionaryEntry(ctx, plan.Job.ID, entry.DictionaryScope, entry.SourceTerm, entry.TranslatedTerm, entry.TermKind, entry.Reusable, termTranslationDictionarySourceSnapshot)
		if err != nil {
			return newTermTranslationAttemptFailure(
				termTranslationErrorKindSaveFailed,
				termTranslationReasonSaveFailed,
				false,
				false,
				termTranslationProgressPercent(baseConfirmedCount, len(plan.Candidates)),
				fmt.Errorf("apply shared dictionary term: %w", err),
			)
		}
		if err := service.linkSnapshotHit(ctx, plan.PhaseRun.ID, entry.ID, jobEntry.ID); err != nil {
			return newTermTranslationAttemptFailure(
				termTranslationErrorKindSaveFailed,
				termTranslationReasonSaveFailed,
				false,
				false,
				termTranslationProgressPercent(baseConfirmedCount, len(plan.Candidates)),
				err,
			)
		}
		confirmedTerms[key] = jobEntry
	}
	return nil
}

func (service *TermTranslationPhaseService) linkSnapshotHit(
	ctx context.Context,
	phaseRunID int64,
	snapshotEntryID int64,
	jobEntryID int64,
) error {
	if err := service.linkPhaseRunDictionaryEntry(ctx, phaseRunID, snapshotEntryID, termTranslationLinkRoleSnapshot); err != nil {
		return fmt.Errorf("link shared dictionary snapshot: %w", err)
	}
	if err := service.linkPhaseRunDictionaryEntry(ctx, phaseRunID, jobEntryID, termTranslationLinkRoleApplied); err != nil {
		return fmt.Errorf("link job dictionary entry: %w", err)
	}
	return nil
}

func (service *TermTranslationPhaseService) buildAITargets(
	candidates []termTranslationCandidate,
	confirmedTerms map[string]repository.DictionaryEntry,
) []termTranslationCandidate {
	aiTargets := make([]termTranslationCandidate, 0, len(candidates))
	for _, candidate := range candidates {
		if _, ok := confirmedTerms[candidateKey(candidate.RecordType, candidate.SourceTerm)]; ok {
			continue
		}
		aiTargets = append(aiTargets, candidate)
	}
	return aiTargets
}

func (service *TermTranslationPhaseService) resumeExecutionPlanRun(
	ctx context.Context,
	mode termTranslationStartMode,
	run repository.JobPhaseRun,
	confirmedCount int,
	totalCount int,
) (repository.JobPhaseRun, error) {
	switch {
	case run.State == termTranslationPhaseStatePaused && mode == termTranslationStartModeResume:
		updatedRun, err := service.updateExecutionPlanRunState(ctx, run, run.ProgressPercent)
		if err != nil {
			return repository.JobPhaseRun{}, fmt.Errorf("resume phase run: %w", err)
		}
		return updatedRun, nil
	case run.State == termTranslationPhaseStateRecoverableFail && (mode == termTranslationStartModeRetry || mode == termTranslationStartModeResume):
		progress := confirmedCount * 100 / maxInt(totalCount, 1)
		updatedRun, err := service.updateExecutionPlanRunState(ctx, run, progress)
		if err != nil {
			return repository.JobPhaseRun{}, fmt.Errorf("retry phase run: %w", err)
		}
		return updatedRun, nil
	default:
		return run, nil
	}
}

func (service *TermTranslationPhaseService) updateExecutionPlanRunState(
	ctx context.Context,
	run repository.JobPhaseRun,
	progressPercent int,
) (repository.JobPhaseRun, error) {
	updatedRun, err := service.jobLifecycleRepository.UpdateJobPhaseRun(ctx, run.ID, repository.JobPhaseRunUpdateDraft{
		State:               termTranslationPhaseStateRunning,
		ProgressPercent:     progressPercent,
		LatestExternalRunID: run.LatestExternalRunID,
		LatestError:         "",
		StartedAt:           termTranslationRunStartedAt(run.StartedAt, service.now()),
		FinishedAt:          nil,
	})
	if err != nil {
		return repository.JobPhaseRun{}, fmt.Errorf("update execution plan run state: %w", err)
	}
	return updatedRun, nil
}

func (service *TermTranslationPhaseService) applyProviderCandidate(
	ctx context.Context,
	plan termTranslationExecutionPlan,
	candidate termTranslationCandidate,
	apiKey string,
	confirmedTerms map[string]repository.DictionaryEntry,
	baseConfirmedCount int,
) error {
	result, err := service.provider.TranslateTerm(ctx, TermTranslationProviderRequest{
		Provider:        plan.Execution.Provider,
		Model:           plan.Execution.Model,
		APIKey:          apiKey,
		EndpointSummary: providerExecutionEndpointSummary(plan.Execution.EndpointSummary),
		SourceTerm:      candidate.SourceTerm,
		SourceLanguage:  "English",
		TargetLanguage:  "Japanese",
		PromptVersion:   "term-translation-v1",
	})
	if err != nil {
		logProviderBoundaryFailure(ctx, "term_translation_provider_execute", termTranslationPhaseServiceWhere, plan.Execution.Provider, classifyProviderBoundaryError(err, termTranslationErrorKindProviderFailure))
		return service.failProviderCandidate(err, baseConfirmedCount, len(plan.Candidates))
	}
	if validateErr := validateProviderCandidateResult(candidate, result); validateErr != nil {
		logProviderBoundaryFailure(ctx, "term_translation_provider_execute", termTranslationPhaseServiceWhere, plan.Execution.Provider, "invalid_provider_response")
		return service.failInvalidProviderCandidate(validateErr, baseConfirmedCount, len(plan.Candidates))
	}
	jobEntry, err := service.ensureJobDictionaryEntry(
		ctx,
		plan.Job.ID,
		candidate.RecordType,
		result.SourceTerm,
		result.TranslatedTerm,
		"",
		true,
		termTranslationDictionarySourceProvider,
	)
	if err != nil {
		return newTermTranslationAttemptFailure(
			termTranslationErrorKindSaveFailed,
			termTranslationReasonSaveFailed,
			false,
			false,
			termTranslationProgressPercent(baseConfirmedCount, len(plan.Candidates)),
			fmt.Errorf("apply provider dictionary result: %w", err),
		)
	}
	if err := service.linkPhaseRunDictionaryEntry(ctx, plan.PhaseRun.ID, jobEntry.ID, termTranslationLinkRoleApplied); err != nil {
		return newTermTranslationAttemptFailure(
			termTranslationErrorKindSaveFailed,
			termTranslationReasonSaveFailed,
			false,
			false,
			termTranslationProgressPercent(baseConfirmedCount, len(plan.Candidates)),
			fmt.Errorf("link provider dictionary result: %w", err),
		)
	}
	confirmedTerms[candidateKey(candidate.RecordType, result.SourceTerm)] = jobEntry
	return nil
}

func (service *TermTranslationPhaseService) failProviderCandidate(
	cause error,
	baseConfirmedCount int,
	totalCount int,
) error {
	errorKind := termTranslationErrorKindProviderFailure
	reason := termTranslationReasonProviderFailed
	retryable := true
	isRedacted := true
	if errors.Is(cause, ErrTermTranslationInvalidProviderReply) {
		errorKind = termTranslationErrorKindInvalidProvider
		reason = termTranslationReasonProviderInvalid
	} else {
		var providerFailure *TermTranslationProviderError
		if errors.As(cause, &providerFailure) {
			if providerFailure.Failure.Kind == TermTranslationProviderErrorKindInvalidProviderResponse {
				errorKind = termTranslationErrorKindInvalidProvider
				reason = termTranslationReasonProviderInvalid
			}
			retryable = providerFailure.Failure.Retryable
			isRedacted = providerFailure.Failure.IsRedacted
		}
	}
	return newTermTranslationAttemptFailure(
		errorKind,
		reason,
		retryable,
		isRedacted,
		termTranslationProgressPercent(baseConfirmedCount, totalCount),
		cause,
	)
}

func (service *TermTranslationPhaseService) failInvalidProviderCandidate(
	cause error,
	baseConfirmedCount int,
	totalCount int,
) error {
	return newTermTranslationAttemptFailure(
		termTranslationErrorKindInvalidProvider,
		termTranslationReasonProviderInvalid,
		true,
		true,
		termTranslationProgressPercent(baseConfirmedCount, totalCount),
		cause,
	)
}

func validateProviderCandidateResult(
	candidate termTranslationCandidate,
	result TermTranslationProviderResult,
) error {
	if strings.TrimSpace(result.TranslatedTerm) == "" {
		return ErrTermTranslationInvalidProviderReply
	}
	if normalizeTermTranslationText(result.SourceTerm) != normalizeTermTranslationText(candidate.SourceTerm) {
		return ErrTermTranslationInvalidProviderReply
	}
	return nil
}

func (service *TermTranslationPhaseService) markPhaseRunCompleted(
	ctx context.Context,
	job repository.TranslationJob,
	run repository.JobPhaseRun,
	totalCount int,
) (repository.JobPhaseRun, int, *TermTranslationPhaseErrorReadModel, error) {
	now := service.now()
	updatedRun, err := service.jobLifecycleRepository.UpdateJobPhaseRun(ctx, run.ID, repository.JobPhaseRunUpdateDraft{
		State:               termTranslationPhaseStateCompleted,
		ProgressPercent:     100,
		LatestExternalRunID: run.LatestExternalRunID,
		LatestError:         "",
		StartedAt:           termTranslationRunStartedAt(run.StartedAt, now),
		FinishedAt:          &now,
	})
	if err != nil {
		return repository.JobPhaseRun{}, 0, nil, fmt.Errorf("complete phase run: %w", err)
	}
	if _, err := service.jobLifecycleRepository.UpdateTranslationJob(ctx, job.ID, repository.TranslationJobUpdateDraft{
		JobName:         job.JobName,
		State:           termTranslationJobStateRunning,
		ProgressPercent: updatedRun.ProgressPercent,
		StartedAt:       termTranslationJobStartedAt(job.StartedAt, now),
		FinishedAt:      nil,
	}); err != nil {
		return repository.JobPhaseRun{}, 0, nil, fmt.Errorf("update translation job after complete: %w", err)
	}
	return updatedRun, totalCount, nil, nil
}

func (service *TermTranslationPhaseService) markPhaseRunFailure(
	ctx context.Context,
	job repository.TranslationJob,
	run repository.JobPhaseRun,
	errorKind string,
	reason string,
	retryable bool,
	progressPercent int,
) (repository.JobPhaseRun, error) {
	now := service.now()
	runState := termTranslationPhaseStateRecoverableFail
	if !retryable {
		runState = termTranslationPhaseStateFailed
	}
	updatedRun, err := service.jobLifecycleRepository.UpdateJobPhaseRun(ctx, run.ID, repository.JobPhaseRunUpdateDraft{
		State:               runState,
		ProgressPercent:     progressPercent,
		LatestExternalRunID: run.LatestExternalRunID,
		LatestError:         errorKind,
		StartedAt:           termTranslationRunStartedAt(run.StartedAt, now),
		FinishedAt:          nil,
	})
	if err != nil {
		return repository.JobPhaseRun{}, fmt.Errorf("update term translation failure state: %w", err)
	}
	if _, err := service.jobLifecycleRepository.UpdateTranslationJob(ctx, job.ID, repository.TranslationJobUpdateDraft{
		JobName:         job.JobName,
		State:           termTranslationJobStateRunning,
		ProgressPercent: updatedRun.ProgressPercent,
		StartedAt:       termTranslationJobStartedAt(job.StartedAt, now),
		FinishedAt:      nil,
	}); err != nil {
		return repository.JobPhaseRun{}, fmt.Errorf("update translation job failure state: %w", err)
	}
	_ = reason
	return updatedRun, nil
}

func (service *TermTranslationPhaseService) loadExecutionContext(
	ctx context.Context,
	jobID int64,
) (repository.TranslationJob, *repository.JobPhaseRun, TermTranslationExecutionConfigReadModel, []termTranslationCandidate, error) {
	if service.jobLifecycleRepository == nil || service.foundationDataRepository == nil || service.translationSourceReader == nil {
		return repository.TranslationJob{}, nil, TermTranslationExecutionConfigReadModel{}, nil, fmt.Errorf("term translation phase persistence is not configured")
	}
	job, err := service.jobLifecycleRepository.GetTranslationJobByID(ctx, jobID)
	if err != nil {
		return repository.TranslationJob{}, nil, TermTranslationExecutionConfigReadModel{}, nil, fmt.Errorf("load translation job: %w", err)
	}
	if _, loadInputErr := service.translationSourceReader.GetXEditExtractedDataByID(ctx, job.XEditExtractedDataID); loadInputErr != nil {
		return repository.TranslationJob{}, nil, TermTranslationExecutionConfigReadModel{}, nil, fmt.Errorf("load translation input metadata: %w", loadInputErr)
	}
	var run *repository.JobPhaseRun
	foundRun, err := service.jobLifecycleRepository.FindJobPhaseRun(ctx, job.ID, termTranslationPhaseType)
	if err == nil {
		run = &foundRun
	} else if !errors.Is(err, repository.ErrNotFound) {
		return repository.TranslationJob{}, nil, TermTranslationExecutionConfigReadModel{}, nil, fmt.Errorf("find term translation phase run: %w", err)
	}
	phases, err := service.jobLifecycleRepository.ListJobPhaseRunsByJobID(ctx, job.ID)
	if err != nil {
		return repository.TranslationJob{}, nil, TermTranslationExecutionConfigReadModel{}, nil, fmt.Errorf("list translation job phases: %w", err)
	}
	initial, applyRuntimeSnapshot, err := termTranslationExecutionBasePhase(job.ID, job.State, run, phases)
	if err != nil {
		return repository.TranslationJob{}, nil, TermTranslationExecutionConfigReadModel{}, nil, err
	}
	initial, err = service.applyTermTranslationRuntimeSnapshot(ctx, job.ID, initial, applyRuntimeSnapshot)
	if err != nil {
		return repository.TranslationJob{}, nil, TermTranslationExecutionConfigReadModel{}, nil, err
	}
	candidates, err := service.collectCandidates(ctx, job.XEditExtractedDataID)
	if err != nil {
		return repository.TranslationJob{}, nil, TermTranslationExecutionConfigReadModel{}, nil, err
	}
	execution := TermTranslationExecutionConfigReadModel{
		CredentialRef: initial.CredentialRef,
		Provider:      initial.AIProvider,
		Model:         initial.ModelName,
		ExecutionMode: initial.ExecutionMode,
	}
	if run != nil {
		execution, err = service.applyTermTranslationRunExecution(ctx, job.ID, *run, execution)
		if err != nil {
			return repository.TranslationJob{}, nil, TermTranslationExecutionConfigReadModel{}, nil, err
		}
	}
	return job, run, execution, candidates, nil
}

func (service *TermTranslationPhaseService) applyTermTranslationRuntimeSnapshot(
	ctx context.Context,
	jobID int64,
	initial repository.JobPhaseRun,
	enabled bool,
) (repository.JobPhaseRun, error) {
	if !enabled {
		return initial, nil
	}
	snapshotReader, ok := service.jobLifecycleRepository.(termTranslationPhaseRuntimeSnapshotReader)
	if !ok {
		return initial, nil
	}
	snapshot, err := snapshotReader.GetTranslationJobPhaseRuntimeSnapshot(ctx, jobID, "word_translation")
	if errors.Is(err, repository.ErrNotFound) {
		initial.AIProvider = ""
		initial.ModelName = ""
		initial.ExecutionMode = ""
		initial.CredentialRef = ""
		return initial, nil
	}
	if err != nil {
		return repository.JobPhaseRun{}, fmt.Errorf("load term translation phase runtime snapshot: %w", err)
	}
	initial.AIProvider = snapshot.Provider
	initial.ModelName = snapshot.ModelName
	initial.ExecutionMode = snapshot.ExecutionMode
	initial.CredentialRef = ""
	return initial, nil
}

func (service *TermTranslationPhaseService) applyTermTranslationRunExecution(
	ctx context.Context,
	jobID int64,
	run repository.JobPhaseRun,
	execution TermTranslationExecutionConfigReadModel,
) (TermTranslationExecutionConfigReadModel, error) {
	execution = termTranslationExecutionFromRun(run, execution)
	if snapshot, ok := service.executionSnapshots[run.ID]; ok {
		execution.CredentialState = snapshot.CredentialState
		return execution, nil
	}
	snapshotReader, ok := service.jobLifecycleRepository.(termTranslationPhaseRuntimeSnapshotReader)
	if !ok {
		return execution, nil
	}
	snapshot, err := snapshotReader.GetTranslationJobPhaseRuntimeSnapshot(ctx, jobID, "word_translation")
	if errors.Is(err, repository.ErrNotFound) {
		return execution, nil
	}
	if err != nil {
		return TermTranslationExecutionConfigReadModel{}, fmt.Errorf("load term translation phase runtime snapshot: %w", err)
	}
	persisted := providerExecutionSnapshotFromRuntimeSnapshot(snapshot)
	execution.CredentialState = persisted.CredentialState
	return execution, nil
}

func termTranslationExecutionFromRun(
	run repository.JobPhaseRun,
	execution TermTranslationExecutionConfigReadModel,
) TermTranslationExecutionConfigReadModel {
	if trimmed := strings.TrimSpace(run.CredentialRef); trimmed != "" {
		execution.CredentialRef = trimmed
	}
	if trimmed := strings.TrimSpace(run.AIProvider); trimmed != "" {
		execution.Provider = trimmed
	}
	if trimmed := strings.TrimSpace(run.ModelName); trimmed != "" {
		execution.Model = trimmed
	}
	if trimmed := strings.TrimSpace(run.ExecutionMode); trimmed != "" {
		execution.ExecutionMode = trimmed
	}
	return execution
}

func (service *TermTranslationPhaseService) buildRunState(
	ctx context.Context,
	jobID int64,
	run *repository.JobPhaseRun,
	candidates []termTranslationCandidate,
) (
	map[string]repository.DictionaryEntry,
	map[string]repository.DictionaryEntry,
	*TermTranslationPhaseResultReadModel,
	TermTranslationPhaseProgressReadModel,
	*TermTranslationPhaseErrorReadModel,
	*string,
	*string,
	error,
) {
	confirmedTerms, err := service.loadJobDictionaryEntries(ctx, jobID)
	if err != nil {
		return nil, nil, nil, TermTranslationPhaseProgressReadModel{}, nil, nil, nil, err
	}
	snapshotHits, digest, version, err := service.loadRunSnapshotState(ctx, run)
	if err != nil {
		return nil, nil, nil, TermTranslationPhaseProgressReadModel{}, nil, nil, nil, err
	}
	confirmedCount, replacementTargetCount := summarizeConfirmedCandidates(candidates, confirmedTerms)
	resultSummary := buildTermTranslationResultSummary(run, candidates, snapshotHits, confirmedTerms, confirmedCount, replacementTargetCount)
	progress := buildTermTranslationProgress(run, candidates, snapshotHits, confirmedCount, resultSummary)
	errorSummary := buildTermTranslationErrorSummary(run)
	return confirmedTerms, snapshotHits, resultSummary, progress, errorSummary, digest, version, nil
}

func (service *TermTranslationPhaseService) loadRunSnapshotState(
	ctx context.Context,
	run *repository.JobPhaseRun,
) (map[string]repository.DictionaryEntry, *string, *string, error) {
	snapshotHits := map[string]repository.DictionaryEntry{}
	if run == nil {
		return snapshotHits, nil, nil, nil
	}
	links, err := service.jobLifecycleRepository.ListPhaseRunDictionaryEntriesByPhaseRunID(ctx, run.ID)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("list phase dictionary links: %w", err)
	}
	snapshotEntries := make([]repository.DictionaryEntry, 0)
	for _, link := range links {
		entry, err := service.foundationDataRepository.GetDictionaryEntryByID(ctx, link.DictionaryEntryID)
		if err != nil {
			return nil, nil, nil, fmt.Errorf("load linked dictionary entry: %w", err)
		}
		if link.Role != termTranslationLinkRoleSnapshot {
			continue
		}
		key := candidateKey(entry.DictionaryScope, entry.SourceTerm)
		snapshotHits[key] = entry
		snapshotEntries = append(snapshotEntries, entry)
	}
	digest, version := termTranslationSnapshotMetadata(snapshotEntries)
	return snapshotHits, digest, version, nil
}

func summarizeConfirmedCandidates(
	candidates []termTranslationCandidate,
	confirmedTerms map[string]repository.DictionaryEntry,
) (int, int) {
	confirmedCount := 0
	replacementTargetCount := 0
	for _, candidate := range candidates {
		if _, ok := confirmedTerms[candidateKey(candidate.RecordType, candidate.SourceTerm)]; !ok {
			continue
		}
		confirmedCount++
		replacementTargetCount += candidate.FieldCount
	}
	return confirmedCount, replacementTargetCount
}

func buildTermTranslationResultSummary(
	run *repository.JobPhaseRun,
	candidates []termTranslationCandidate,
	snapshotHits map[string]repository.DictionaryEntry,
	confirmedTerms map[string]repository.DictionaryEntry,
	confirmedCount int,
	replacementTargetCount int,
) *TermTranslationPhaseResultReadModel {
	if len(candidates) == 0 && run == nil {
		return nil
	}
	return &TermTranslationPhaseResultReadModel{
		ConfirmedCount:            confirmedCount,
		JobDictionaryAppliedCount: len(confirmedTerms),
		ReplacementTargetCount:    replacementTargetCount,
		UnmatchedCount:            len(candidates) - confirmedCount,
		ProviderSkipped:           len(candidates) == len(snapshotHits),
	}
}

func buildTermTranslationProgress(
	run *repository.JobPhaseRun,
	candidates []termTranslationCandidate,
	snapshotHits map[string]repository.DictionaryEntry,
	confirmedCount int,
	resultSummary *TermTranslationPhaseResultReadModel,
) TermTranslationPhaseProgressReadModel {
	progress := TermTranslationPhaseProgressReadModel{
		Percent:        0,
		ProcessedCount: confirmedCount,
		TotalCount:     len(candidates),
		AITargetCount:  maxInt(len(candidates)-len(snapshotHits), 0),
		CurrentStep:    termTranslationCurrentStep(run, resultSummary),
	}
	if run == nil {
		return progress
	}
	progress.Percent = run.ProgressPercent
	if run.ProgressPercent == 100 {
		progress.ProcessedCount = len(candidates)
	}
	return progress
}

func buildTermTranslationErrorSummary(run *repository.JobPhaseRun) *TermTranslationPhaseErrorReadModel {
	if run == nil || strings.TrimSpace(run.LatestError) == "" {
		return nil
	}
	retryable := run.State == termTranslationPhaseStateRecoverableFail
	isRedacted := false
	switch strings.TrimSpace(run.LatestError) {
	case termTranslationErrorKindProviderFailure, termTranslationErrorKindInvalidProvider:
		isRedacted = true
	}
	return &TermTranslationPhaseErrorReadModel{
		ErrorKind:  run.LatestError,
		Reason:     termTranslationReasonFromErrorKind(run.LatestError),
		Retryable:  retryable,
		IsRedacted: isRedacted,
	}
}

func (service *TermTranslationPhaseService) collectCandidates(
	ctx context.Context,
	xEditID int64,
) ([]termTranslationCandidate, error) {
	records, err := service.translationSourceReader.ListTranslationRecordsByXEditID(ctx, xEditID)
	if err != nil {
		return nil, fmt.Errorf("list translation records for term phase: %w", err)
	}
	candidateMap := make(map[string]termTranslationCandidate)
	for _, record := range records {
		fields, err := service.translationSourceReader.ListTranslationFieldsByTranslationRecordID(ctx, record.ID)
		if err != nil {
			return nil, fmt.Errorf("list translation fields for term phase: %w", err)
		}
		for _, field := range fields {
			sourceTerm := strings.TrimSpace(field.SourceText)
			if sourceTerm == "" {
				continue
			}
			key := candidateKey(record.RecordType, sourceTerm)
			candidate := candidateMap[key]
			if candidate.NormalizedSource == "" {
				candidate = termTranslationCandidate{
					RecordType:       record.RecordType,
					SourceTerm:       sourceTerm,
					NormalizedSource: normalizeTermTranslationText(sourceTerm),
				}
			}
			candidate.FieldCount++
			candidateMap[key] = candidate
		}
	}
	candidates := make([]termTranslationCandidate, 0, len(candidateMap))
	for _, candidate := range candidateMap {
		candidates = append(candidates, candidate)
	}
	sort.Slice(candidates, func(i, j int) bool {
		if candidates[i].RecordType == candidates[j].RecordType {
			return candidates[i].NormalizedSource < candidates[j].NormalizedSource
		}
		return candidates[i].RecordType < candidates[j].RecordType
	})
	return candidates, nil
}

func (service *TermTranslationPhaseService) loadSnapshotHits(
	ctx context.Context,
	candidates []termTranslationCandidate,
) (map[string]repository.DictionaryEntry, error) {
	entries, err := service.foundationDataRepository.ListDictionaryEntries(ctx, nil, termTranslationDictionaryLifecycleMaster, "", "")
	if err != nil {
		return nil, fmt.Errorf("list shared dictionary entries: %w", err)
	}
	index := make(map[string]repository.DictionaryEntry, len(entries))
	for _, entry := range entries {
		index[candidateKey(entry.DictionaryScope, entry.SourceTerm)] = entry
	}
	result := make(map[string]repository.DictionaryEntry)
	for _, candidate := range candidates {
		if entry, ok := index[candidateKey(candidate.RecordType, candidate.SourceTerm)]; ok {
			result[candidateKey(candidate.RecordType, candidate.SourceTerm)] = entry
		}
	}
	return result, nil
}

func (service *TermTranslationPhaseService) loadJobDictionaryEntries(
	ctx context.Context,
	jobID int64,
) (map[string]repository.DictionaryEntry, error) {
	entries, err := service.foundationDataRepository.ListDictionaryEntries(ctx, int64Pointer(jobID), termTranslationDictionaryLifecycleJob, "", "")
	if err != nil {
		return nil, fmt.Errorf("list job dictionary entries: %w", err)
	}
	result := make(map[string]repository.DictionaryEntry, len(entries))
	for _, entry := range entries {
		result[candidateKey(entry.DictionaryScope, entry.SourceTerm)] = entry
	}
	return result, nil
}

func (service *TermTranslationPhaseService) ensureJobDictionaryEntry(
	ctx context.Context,
	jobID int64,
	recordType string,
	sourceTerm string,
	translatedTerm string,
	termKind string,
	reusable bool,
	dictionarySource string,
) (repository.DictionaryEntry, error) {
	entry, err := service.foundationDataRepository.FindDictionaryEntry(
		ctx,
		int64Pointer(jobID),
		termTranslationDictionaryLifecycleJob,
		recordType,
		sourceTerm,
	)
	if err == nil {
		return entry, nil
	}
	if !errors.Is(err, repository.ErrNotFound) {
		return repository.DictionaryEntry{}, fmt.Errorf("find job dictionary entry: %w", err)
	}
	created, err := service.foundationDataRepository.CreateDictionaryEntry(ctx, repository.DictionaryEntryDraft{
		TranslationJobID:    int64Pointer(jobID),
		DictionaryLifecycle: termTranslationDictionaryLifecycleJob,
		DictionaryScope:     recordType,
		DictionarySource:    dictionarySource,
		SourceTerm:          sourceTerm,
		TranslatedTerm:      translatedTerm,
		TermKind:            strings.TrimSpace(termKind),
		Reusable:            reusable,
	})
	if err == nil {
		return created, nil
	}
	if !errors.Is(err, repository.ErrConflict) {
		return repository.DictionaryEntry{}, fmt.Errorf("create job dictionary entry: %w", err)
	}
	reloaded, reloadErr := service.foundationDataRepository.FindDictionaryEntry(
		ctx,
		int64Pointer(jobID),
		termTranslationDictionaryLifecycleJob,
		recordType,
		sourceTerm,
	)
	if reloadErr != nil {
		return repository.DictionaryEntry{}, fmt.Errorf("reload conflicted job dictionary entry: %w", reloadErr)
	}
	return reloaded, nil
}

func (service *TermTranslationPhaseService) linkPhaseRunDictionaryEntry(
	ctx context.Context,
	phaseRunID int64,
	dictionaryEntryID int64,
	role string,
) error {
	_, err := service.jobLifecycleRepository.CreatePhaseRunDictionaryEntry(ctx, repository.PhaseRunDictionaryEntryDraft{
		PhaseRunID:        phaseRunID,
		DictionaryEntryID: dictionaryEntryID,
		Role:              role,
	})
	if err == nil || errors.Is(err, repository.ErrConflict) {
		return nil
	}
	return fmt.Errorf("create phase_run_dictionary_entry: %w", err)
}

func (service *TermTranslationPhaseService) loadExistingRunForMutation(
	ctx context.Context,
	jobID int64,
	phaseRunID int64,
) (repository.TranslationJob, repository.JobPhaseRun, error) {
	job, run, _, _, err := service.loadExecutionContext(ctx, jobID)
	if err != nil {
		return repository.TranslationJob{}, repository.JobPhaseRun{}, err
	}
	if run == nil || run.ID != phaseRunID {
		return repository.TranslationJob{}, repository.JobPhaseRun{}, fmt.Errorf("load term translation phase run: %w", repository.ErrNotFound)
	}
	return job, *run, nil
}

func (service *TermTranslationPhaseService) readinessFromState(
	job repository.TranslationJob,
	run *repository.JobPhaseRun,
	totalCount int,
	confirmedCount int,
	errorSummary *TermTranslationPhaseErrorReadModel,
) TermTranslationNextPhaseReadinessReadModel {
	if termTranslationJobIsTerminal(job.State) {
		reason := termTranslationReasonTerminalJob
		return TermTranslationNextPhaseReadinessReadModel{
			JobID:             job.ID,
			CurrentPhase:      termTranslationCurrentPhase,
			PhaseState:        termTranslationPhaseStateFromRun(run),
			CanStartNextPhase: false,
			BlockedReason:     &reason,
			ErrorKind:         "terminal_job",
		}
	}
	if run == nil || run.State != termTranslationPhaseStateCompleted || confirmedCount < totalCount {
		reason := "term phase is not completed"
		if errorSummary != nil && strings.TrimSpace(errorSummary.Reason) != "" {
			reason = errorSummary.Reason
		}
		return TermTranslationNextPhaseReadinessReadModel{
			JobID:             job.ID,
			CurrentPhase:      termTranslationCurrentPhase,
			PhaseState:        termTranslationPhaseStateFromRun(run),
			CanStartNextPhase: false,
			BlockedReason:     &reason,
			ErrorKind:         "term_phase_incomplete",
		}
	}
	return TermTranslationNextPhaseReadinessReadModel{
		JobID:             job.ID,
		CurrentPhase:      termTranslationCurrentPhase,
		PhaseState:        run.State,
		CanStartNextPhase: true,
	}
}

func termTranslationInitialExecutionPhase(phases []repository.JobPhaseRun) (repository.JobPhaseRun, error) {
	for _, phase := range phases {
		if phase.PhaseType == termTranslationInitialPhaseType {
			return phase, nil
		}
	}
	return repository.JobPhaseRun{}, fmt.Errorf("load initial execution phase: %w", repository.ErrNotFound)
}

func termTranslationExecutionBasePhase(
	jobID int64,
	jobState string,
	run *repository.JobPhaseRun,
	phases []repository.JobPhaseRun,
) (repository.JobPhaseRun, bool, error) {
	initial, err := termTranslationInitialExecutionPhase(phases)
	if err == nil {
		return initial, run == nil, nil
	}
	if run != nil || !errors.Is(err, repository.ErrNotFound) {
		return repository.JobPhaseRun{}, false, err
	}
	_ = jobID
	_ = jobState
	return repository.JobPhaseRun{}, false, fmt.Errorf("load initial execution phase: %w", repository.ErrNotFound)
}

func termTranslationCurrentStep(run *repository.JobPhaseRun, result *TermTranslationPhaseResultReadModel) string {
	if run == nil {
		return "idle"
	}
	switch run.State {
	case termTranslationPhaseStateCompleted:
		if result != nil && result.ProviderSkipped {
			return "shared_dictionary_completed"
		}
		return "completed"
	case termTranslationPhaseStatePaused:
		return "paused"
	case termTranslationPhaseStateRecoverableFail:
		return "retry_waiting"
	case termTranslationPhaseStateFailed:
		return "failed"
	default:
		return "dictionary_snapshot"
	}
}

func termTranslationReasonFromErrorKind(errorKind string) string {
	switch strings.TrimSpace(errorKind) {
	case "ready_required":
		return termTranslationReasonReadyRequired
	case "terminal_job":
		return termTranslationReasonTerminalJob
	case "active_phase_exists":
		return termTranslationReasonActivePhaseExists
	case "dictionary_snapshot_failed":
		return "dictionary snapshot failed"
	case termTranslationErrorKindInvalidProvider:
		return termTranslationReasonProviderInvalid
	case termTranslationErrorKindSaveFailed:
		return termTranslationReasonSaveFailed
	case "term_phase_incomplete":
		return "term phase is not completed"
	default:
		return termTranslationReasonProviderFailed
	}
}

func termTranslationSnapshotMetadata(entries []repository.DictionaryEntry) (*string, *string) {
	if len(entries) == 0 {
		return nil, nil
	}
	sort.Slice(entries, func(i, j int) bool { return entries[i].ID < entries[j].ID })
	hasher := sha256.New()
	latestUpdatedAt := entries[0].UpdatedAt
	for _, entry := range entries {
		_, _ = fmt.Fprintf(hasher, "%d:%s:%s:%d;", entry.ID, entry.DictionaryScope, entry.SourceTerm, entry.UpdatedAt.UnixNano())
		if entry.UpdatedAt.After(latestUpdatedAt) {
			latestUpdatedAt = entry.UpdatedAt
		}
	}
	digest := hex.EncodeToString(hasher.Sum(nil))
	version := fmt.Sprintf("%d:%s", len(entries), latestUpdatedAt.UTC().Format(time.RFC3339))
	return &digest, &version
}

func termTranslationRejectedCommand(jobID int64, plan termTranslationExecutionPlan) TermTranslationPhaseCommandReadModel {
	errorSummary := plan.Rejection
	if errorSummary == nil {
		errorSummary = &TermTranslationPhaseErrorReadModel{
			ErrorKind:  "ready_required",
			Reason:     termTranslationReasonReadyRequired,
			Retryable:  false,
			IsRedacted: false,
		}
	}
	return TermTranslationPhaseCommandReadModel{
		JobID:             jobID,
		CurrentPhase:      termTranslationCurrentPhase,
		BeforePhaseState:  plan.BeforePhaseState,
		AfterPhaseState:   termTranslationPhaseStateFromRun(&plan.PhaseRun),
		PhaseState:        termTranslationPhaseStateFromRun(&plan.PhaseRun),
		PhaseRunID:        int64PointerOrNil(plan.PhaseRun.ID),
		StartedAt:         cloneTimePointer(plan.PhaseRun.StartedAt),
		FinishedAt:        cloneTimePointer(plan.PhaseRun.FinishedAt),
		Progress:          TermTranslationPhaseProgressReadModel{Percent: plan.PhaseRun.ProgressPercent},
		Retryable:         errorSummary.Retryable,
		CanStartNextPhase: false,
		ErrorSummary:      errorSummary,
	}
}

func termTranslationCommandFromRun(
	job repository.TranslationJob,
	run repository.JobPhaseRun,
	retryable bool,
	canStartNextPhase bool,
	errorSummary *TermTranslationPhaseErrorReadModel,
) TermTranslationPhaseCommandReadModel {
	return TermTranslationPhaseCommandReadModel{
		JobID:             job.ID,
		CurrentPhase:      termTranslationCurrentPhase,
		BeforePhaseState:  run.State,
		AfterPhaseState:   run.State,
		PhaseState:        run.State,
		PhaseRunID:        int64Pointer(run.ID),
		StartedAt:         cloneTimePointer(run.StartedAt),
		FinishedAt:        cloneTimePointer(run.FinishedAt),
		Progress:          TermTranslationPhaseProgressReadModel{Percent: run.ProgressPercent},
		Retryable:         retryable,
		CanStartNextPhase: canStartNextPhase && run.State == termTranslationPhaseStateCompleted,
		ErrorSummary:      errorSummary,
	}
}

func termTranslationJobIsTerminal(state string) bool {
	switch strings.TrimSpace(strings.ToLower(state)) {
	case termTranslationJobStateCompleted, termTranslationJobStateFailed, termTranslationJobStateCanceled:
		return true
	default:
		return false
	}
}

func termTranslationPhaseIsActive(run *repository.JobPhaseRun) bool {
	if run == nil {
		return false
	}
	switch strings.TrimSpace(strings.ToLower(run.State)) {
	case termTranslationPhaseStateRunning, termTranslationPhaseStatePaused, termTranslationPhaseStateRecoverableFail:
		return true
	default:
		return false
	}
}

func termTranslationPhaseStateFromRun(run *repository.JobPhaseRun) string {
	if run == nil || strings.TrimSpace(run.State) == "" {
		return termTranslationPhaseStateIdleReady
	}
	return run.State
}

func termTranslationStartBlockedReason(
	job repository.TranslationJob,
	run *repository.JobPhaseRun,
	execution TermTranslationExecutionConfigReadModel,
) *string {
	if termTranslationJobIsTerminal(job.State) {
		return termTranslationStringPointer(termTranslationReasonTerminalJob)
	}
	if run != nil && termTranslationPhaseIsActive(run) {
		return termTranslationStringPointer(termTranslationReasonActivePhaseExists)
	}
	if job.State != termTranslationJobStateReady {
		return termTranslationStringPointer(termTranslationReasonReadyRequired)
	}
	if !termTranslationExecutionConfigured(execution) {
		return termTranslationStringPointer(termTranslationReasonPhaseRuntimeMissing)
	}
	return nil
}

func termTranslationExecutionConfigured(execution TermTranslationExecutionConfigReadModel) bool {
	return strings.TrimSpace(execution.Provider) != "" &&
		strings.TrimSpace(execution.Model) != "" &&
		strings.TrimSpace(execution.ExecutionMode) != ""
}

func termTranslationPauseBlockedReason(run *repository.JobPhaseRun) *string {
	if run == nil || run.State == termTranslationPhaseStateRunning {
		return nil
	}
	return termTranslationStringPointer("phase is not running")
}

func termTranslationResumeBlockedReason(run *repository.JobPhaseRun) *string {
	if run == nil || run.State == termTranslationPhaseStatePaused {
		return nil
	}
	return termTranslationStringPointer("phase is not paused")
}

func termTranslationRetryBlockedReason(run *repository.JobPhaseRun) *string {
	if run == nil || run.State == termTranslationPhaseStateRecoverableFail {
		return nil
	}
	return termTranslationStringPointer("phase is not retryable")
}

func candidateKey(recordType string, sourceTerm string) string {
	return normalizeTermTranslationText(recordType) + "\x00" + normalizeTermTranslationText(sourceTerm)
}

func normalizeTermTranslationText(value string) string {
	return strings.ToLower(strings.TrimSpace(value))
}

func newTermTranslationAttemptFailure(
	errorKind string,
	reason string,
	retryable bool,
	isRedacted bool,
	progressPercent int,
	cause error,
) error {
	return &termTranslationAttemptFailure{
		summary: TermTranslationPhaseErrorReadModel{
			ErrorKind:  errorKind,
			Reason:     reason,
			Retryable:  retryable,
			IsRedacted: isRedacted,
		},
		progressPercent: progressPercent,
		cause:           cause,
	}
}

func termTranslationProgressPercent(confirmedCount int, totalCount int) int {
	if totalCount <= 0 {
		return 0
	}
	if confirmedCount <= 0 {
		return 0
	}
	return confirmedCount * 100 / totalCount
}

func (service *TermTranslationPhaseService) persistAttemptFailure(
	ctx context.Context,
	failure termTranslationPersistedFailure,
) (repository.JobPhaseRun, error) {
	var failedRun repository.JobPhaseRun
	err := service.transactor.WithTransaction(ctx, func(txCtx context.Context) error {
		run, runErr := service.loadPersistedFailureRun(txCtx, failure)
		if runErr != nil {
			return runErr
		}
		persistedRun, persistErr := service.persistMarkedAttemptFailure(txCtx, failure, run)
		if persistErr != nil {
			return persistErr
		}
		failedRun = persistedRun
		return nil
	})
	if err != nil {
		return repository.JobPhaseRun{}, fmt.Errorf("persist term translation failure state: %w", err)
	}
	return failedRun, nil
}

func (service *TermTranslationPhaseService) loadPersistedFailureRun(
	ctx context.Context,
	failure termTranslationPersistedFailure,
) (repository.JobPhaseRun, error) {
	if !failure.runCreated {
		return failure.run, nil
	}
	run, err := service.recreateFailurePhaseRun(ctx, failure.job.ID, failure.run)
	if err != nil {
		return repository.JobPhaseRun{}, err
	}
	return run, nil
}

func (service *TermTranslationPhaseService) recreateFailurePhaseRun(
	ctx context.Context,
	jobID int64,
	run repository.JobPhaseRun,
) (repository.JobPhaseRun, error) {
	recreatedRun, recreateErr := service.jobLifecycleRepository.CreateJobPhaseRun(ctx, repository.JobPhaseRunDraft{
		TranslationJobID: jobID,
		PhaseType:        run.PhaseType,
		State:            run.State,
		ExecutionOrder:   run.ExecutionOrder,
		AIProvider:       run.AIProvider,
		ModelName:        run.ModelName,
		ExecutionMode:    run.ExecutionMode,
		CredentialRef:    run.CredentialRef,
		InstructionKind:  run.InstructionKind,
	})
	if recreateErr == nil {
		return recreatedRun, nil
	}
	if !errors.Is(recreateErr, repository.ErrConflict) {
		return repository.JobPhaseRun{}, fmt.Errorf("recreate term translation phase run: %w", recreateErr)
	}
	reloadedRun, reloadErr := service.jobLifecycleRepository.FindJobPhaseRun(ctx, jobID, termTranslationPhaseType)
	if reloadErr != nil {
		return repository.JobPhaseRun{}, fmt.Errorf("reload term translation phase run after recreate conflict: %w", reloadErr)
	}
	return reloadedRun, nil
}

func (service *TermTranslationPhaseService) persistMarkedAttemptFailure(
	ctx context.Context,
	failure termTranslationPersistedFailure,
	run repository.JobPhaseRun,
) (repository.JobPhaseRun, error) {
	if err := service.persistFailureSnapshots(ctx, run.ID, failure.snapshots); err != nil {
		return repository.JobPhaseRun{}, err
	}
	failedRun, err := service.markPhaseRunFailure(
		ctx,
		failure.job,
		run,
		failure.failure.summary.ErrorKind,
		failure.failure.summary.Reason,
		failure.failure.summary.Retryable,
		failure.failure.progressPercent,
	)
	if err != nil {
		return repository.JobPhaseRun{}, err
	}
	return failedRun, nil
}

func (service *TermTranslationPhaseService) persistFailureSnapshots(
	ctx context.Context,
	phaseRunID int64,
	snapshots map[string]repository.DictionaryEntry,
) error {
	for _, snapshot := range snapshots {
		if err := service.linkPhaseRunDictionaryEntry(ctx, phaseRunID, snapshot.ID, termTranslationLinkRoleSnapshot); err != nil {
			return fmt.Errorf("persist failure snapshot link: %w", err)
		}
	}
	return nil
}

func cloneTermTranslationDictionaryEntries(
	source map[string]repository.DictionaryEntry,
) map[string]repository.DictionaryEntry {
	if len(source) == 0 {
		return nil
	}
	cloned := make(map[string]repository.DictionaryEntry, len(source))
	for key, entry := range source {
		cloned[key] = entry
	}
	return cloned
}

func (service *TermTranslationPhaseService) resolveProviderAPIKey(
	ctx context.Context,
	execution TermTranslationExecutionConfigReadModel,
) (string, error) {
	providerID, err := NormalizeTermTranslationProvider(execution.Provider)
	if err != nil {
		return "", err
	}
	if providerID == TermTranslationProviderLMStudio {
		if service.secretStore == nil {
			return "", nil
		}
		secretValue, loaded, loadErr := service.loadProviderSecret(ctx, strings.TrimSpace(execution.CredentialRef))
		if loadErr != nil {
			return "", fmt.Errorf("load term translation provider secret: %w", loadErr)
		}
		if !loaded {
			return "", nil
		}
		return strings.TrimSpace(secretValue), nil
	}
	if service.secretStore == nil {
		return "", fmt.Errorf("term translation provider secret store is not configured")
	}
	credentialRef := strings.TrimSpace(execution.CredentialRef)
	if credentialRef == "" {
		return "", errors.New(termTranslationReasonCredentialMissing)
	}
	secretValue, loaded, err := service.loadProviderSecret(ctx, credentialRef)
	if err != nil {
		return "", fmt.Errorf("load term translation provider secret: %w", err)
	}
	if !loaded {
		return "", errors.New(termTranslationReasonCredentialMissing)
	}
	trimmedSecret := strings.TrimSpace(secretValue)
	if trimmedSecret == "" {
		return "", errors.New(termTranslationReasonCredentialMissing)
	}
	return trimmedSecret, nil
}

func classifyTermTranslationCredentialFailure(err error) string {
	if err == nil {
		return "unknown"
	}
	if strings.Contains(err.Error(), termTranslationReasonCredentialMissing) {
		return providerSettingsErrorKindCredentialMissing
	}
	return "secret_store_failed"
}

type termTranslationSecretLoadResult struct {
	value string
	err   error
}

func (service *TermTranslationPhaseService) loadProviderSecret(
	ctx context.Context,
	key string,
) (string, bool, error) {
	if service == nil || service.secretStore == nil {
		return "", false, nil
	}

	timeout := service.secretLoadTimeout
	if timeout <= 0 {
		timeout = termTranslationSecretLoadTimeout
	}

	requestContext := termTranslationRequestContext(ctx)
	resultChannel := make(chan termTranslationSecretLoadResult, 1)
	go func() {
		value, err := service.secretStore.Load(requestContext, key)
		resultChannel <- termTranslationSecretLoadResult{value: value, err: err}
	}()

	timer := time.NewTimer(timeout)
	defer timer.Stop()

	select {
	case <-requestContext.Done():
		return "", false, fmt.Errorf("term translation secret load canceled: %w", requestContext.Err())
	case <-timer.C:
		return "", false, nil
	case result := <-resultChannel:
		if result.err != nil {
			return "", false, result.err
		}
		return result.value, true, nil
	}
}

func termTranslationProviderSecretKey(provider string) string {
	switch normalizeTermTranslationText(provider) {
	case normalizeTermTranslationText(string(TermTranslationProviderGemini)):
		return "gemini-primary"
	case normalizeTermTranslationText(string(TermTranslationProviderXAI)):
		return "xai-primary"
	case normalizeTermTranslationText(string(TermTranslationProviderLMStudio)):
		return "lmstudio-local"
	default:
		return normalizeTermTranslationText(provider)
	}
}

func int64Pointer(value int64) *int64 {
	return &value
}

func int64PointerOrNil(value int64) *int64 {
	if value == 0 {
		return nil
	}
	return &value
}

func termTranslationStringPointer(value string) *string {
	return &value
}

func cloneTimePointer(value *time.Time) *time.Time {
	if value == nil {
		return nil
	}
	cloned := *value
	return &cloned
}

func cloneInt64Pointer(value *int64) *int64 {
	if value == nil {
		return nil
	}
	cloned := *value
	return &cloned
}

func cloneStringPointer(value *string) *string {
	if value == nil {
		return nil
	}
	cloned := *value
	return &cloned
}

func valueOrEmptyRun(run *repository.JobPhaseRun) repository.JobPhaseRun {
	if run == nil {
		return repository.JobPhaseRun{}
	}
	return *run
}

func maxInt(a int, b int) int {
	if a > b {
		return a
	}
	return b
}

func termTranslationRequestContext(ctx context.Context) context.Context {
	if ctx == nil {
		return context.Background()
	}
	return ctx
}

func termTranslationJobStartedAt(startedAt *time.Time, now time.Time) *time.Time {
	if startedAt != nil {
		return startedAt
	}
	return &now
}

func termTranslationRunStartedAt(startedAt *time.Time, now time.Time) *time.Time {
	if startedAt != nil {
		return startedAt
	}
	return &now
}
