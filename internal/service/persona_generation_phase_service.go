package service

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"sort"
	"strings"
	"time"

	"aitranslationenginejp/internal/repository"
)

const (
	personaGenerationCurrentPhase              = "persona_generation"
	personaGenerationPhaseType                 = "persona_generation"
	personaGenerationTermPhaseType             = "term_translation"
	personaGenerationInitialPhaseType          = "translation"
	personaGenerationPhaseStateNotStarted      = "not_started"
	personaGenerationPhaseStateRunning         = "running"
	personaGenerationPhaseStateCompleted       = "completed"
	personaGenerationPhaseStatePaused          = "paused"
	personaGenerationPhaseStateRecoverableFail = "recoverable_failed"
	personaGenerationPhaseStateFailed          = "failed"
	personaGenerationPhaseStateCanceled        = "canceled"
	personaGenerationPhaseStateRejected        = "rejected"
	personaGenerationPhaseStateSnapshotMissing = "snapshot_missing"
	personaGenerationPhaseStateEmptyCompleted  = "empty_completed"
	personaGenerationJobStateCompleted         = "completed"
	personaGenerationJobStateFailed            = "failed"
	personaGenerationJobStateCanceled          = "canceled"
	personaGenerationJobStateRunning           = "running"
	personaGenerationInstructionKindDefault    = "default"
	personaGenerationSkippedReasonOrphan       = "orphan_npc_reference"
	personaGenerationSkippedReasonNoUtterance  = "no_source_utterance"
	personaGenerationPhaseLinkRoleSnapshot     = "snapshot_ref"
	personaGenerationPhaseLinkRoleApplied      = "applied"
	personaGenerationEvidenceRoleSource        = "source_utterance"
	personaGenerationPersonaLifecycleJob       = "job"
	personaGenerationPersonaScopeJob           = "job"
	personaGenerationPersonaSourceProvider     = "provider"
	personaGenerationErrorKindTermIncomplete   = "term_phase_incomplete"
	personaGenerationErrorKindTerminalJob      = "terminal_job"
	personaGenerationErrorKindActivePhase      = "active_phase_exists"
	personaGenerationErrorKindProviderFailure  = "provider_failure"
	personaGenerationErrorKindInvalidProvider  = "invalid_provider_response"
	personaGenerationErrorKindSaveFailed       = "save_failed"
	personaGenerationErrorKindSnapshotMissing  = "snapshot_missing"
	personaGenerationErrorKindBodyBlocked      = "body_readiness_blocked"
	personaGenerationErrorKindInputMissing     = "input_missing"
	personaGenerationTerminalJobReason         = "terminal job"
	personaGenerationRunLoadErrorMessage       = "load persona generation phase run: %w"
	personaGenerationRedactedPublicSummary     = "redacted public summary"
)

type personaGenerationPhaseJobLifecycleRepository interface {
	GetTranslationJobByID(ctx context.Context, id int64) (repository.TranslationJob, error)
	UpdateTranslationJob(ctx context.Context, id int64, draft repository.TranslationJobUpdateDraft) (repository.TranslationJob, error)
	CreateJobPhaseRun(ctx context.Context, draft repository.JobPhaseRunDraft) (repository.JobPhaseRun, error)
	FindJobPhaseRun(ctx context.Context, translationJobID int64, phaseType string) (repository.JobPhaseRun, error)
	UpdateJobPhaseRun(ctx context.Context, id int64, draft repository.JobPhaseRunUpdateDraft) (repository.JobPhaseRun, error)
	ListJobPhaseRunsByJobID(ctx context.Context, jobID int64) ([]repository.JobPhaseRun, error)
	CreatePhaseRunPersona(ctx context.Context, draft repository.PhaseRunPersonaDraft) (repository.PhaseRunPersona, error)
	ListPhaseRunPersonasByPhaseRunID(ctx context.Context, phaseRunID int64) ([]repository.PhaseRunPersona, error)
}

type personaGenerationPhaseFoundationDataRepository interface {
	CreatePersona(ctx context.Context, draft repository.PersonaDraft) (repository.Persona, error)
	GetPersonaByID(ctx context.Context, id int64) (repository.Persona, error)
	GetPersonaByNpcProfileID(ctx context.Context, npcProfileID int64) (repository.Persona, error)
	UpdatePersona(ctx context.Context, id int64, draft repository.PersonaUpdateDraft) (repository.Persona, error)
	CreatePersonaFieldEvidence(ctx context.Context, draft repository.PersonaFieldEvidenceDraft) (repository.PersonaFieldEvidence, error)
	ListPersonaFieldEvidenceByPersonaID(ctx context.Context, personaID int64) ([]repository.PersonaFieldEvidence, error)
}

type personaGenerationPhaseTranslationSourceRepository interface {
	GetXEditExtractedDataByID(ctx context.Context, id int64) (repository.XEditExtractedData, error)
	GetTranslationRecordByID(ctx context.Context, id int64) (repository.TranslationRecord, error)
	ListTranslationRecordsByXEditID(ctx context.Context, xEditID int64) ([]repository.TranslationRecord, error)
	GetNpcRecordByTranslationRecordID(ctx context.Context, translationRecordID int64) (repository.NpcRecord, error)
	GetNpcProfileByID(ctx context.Context, id int64) (repository.NpcProfile, error)
	ListTranslationFieldsByTranslationRecordID(ctx context.Context, translationRecordID int64) ([]repository.TranslationField, error)
	ListTranslationFieldRecordReferencesByFieldID(ctx context.Context, fieldID int64) ([]repository.TranslationFieldRecordReference, error)
}

// PersonaGenerationPhaseService orchestrates persona phase state and target snapshots.
type PersonaGenerationPhaseService struct {
	now                      func() time.Time
	jobLifecycleRepository   personaGenerationPhaseJobLifecycleRepository
	foundationDataRepository personaGenerationPhaseFoundationDataRepository
	translationSourceReader  personaGenerationPhaseTranslationSourceRepository
	transactor               repository.Transactor
	provider                 PersonaGenerationProvider
}

// PersonaGenerationPhaseProgressReadModel stores progress for one persona phase run.
type PersonaGenerationPhaseProgressReadModel struct {
	Percent        int
	ProcessedCount int
	TotalCount     int
	TargetCount    int
	CurrentStep    string
}

// PersonaGenerationTargetSummaryReadModel stores target extraction summary fields.
type PersonaGenerationTargetSummaryReadModel struct {
	TargetCount            int
	CommonPersonaHitCount  int
	CommonPersonaMissCount int
	SkippedCount           int
	SkippedReasons         []string
	TargetSnapshotID       *string
	TargetSnapshotDigest   string
}

// PersonaGenerationExecutionSummaryReadModel stores redacted execution settings.
type PersonaGenerationExecutionSummaryReadModel struct {
	CredentialRef string
	Provider      string
	Model         string
	ExecutionMode string
	PromptDigest  string
	InputCount    int
	OutputCount   int
	EvidenceRefs  []string
}

// PersonaGenerationPhaseResultSummaryReadModel stores snapshot result fields.
type PersonaGenerationPhaseResultSummaryReadModel struct {
	GeneratedCount          int
	FailedCount             int
	PersonaCount            int
	MissingCount            int
	SnapshotID              string
	SnapshotDigest          string
	SnapshotReferenceStatus string
	BodyReadiness           bool
}

// PersonaGenerationPhaseErrorSummaryReadModel stores redacted error fields.
type PersonaGenerationPhaseErrorSummaryReadModel struct {
	ErrorKind  string
	Reason     string
	Retryable  bool
	IsRedacted bool
}

// PersonaGenerationPhaseActionEnablementReadModel stores button enablement state.
type PersonaGenerationPhaseActionEnablementReadModel struct {
	CanStart               bool
	StartBlockedReason     *string
	CanPause               bool
	PauseBlockedReason     *string
	CanResume              bool
	ResumeBlockedReason    *string
	CanRetry               bool
	RetryBlockedReason     *string
	CanCancel              bool
	CancelBlockedReason    *string
	CanStartBodyPhase      bool
	BodyPhaseBlockedReason *string
}

// PersonaGenerationPhaseSummaryReadModel stores the Job Run summary payload.
type PersonaGenerationPhaseSummaryReadModel struct {
	JobID            int64
	CurrentPhase     string
	PhaseState       string
	PhaseRunID       *int64
	StartedAt        *time.Time
	FinishedAt       *time.Time
	Progress         PersonaGenerationPhaseProgressReadModel
	TargetSummary    PersonaGenerationTargetSummaryReadModel
	Execution        PersonaGenerationExecutionSummaryReadModel
	ResultSummary    *PersonaGenerationPhaseResultSummaryReadModel
	ErrorSummary     *PersonaGenerationPhaseErrorSummaryReadModel
	ActionEnablement PersonaGenerationPhaseActionEnablementReadModel
}

// PersonaGenerationPhaseCommandReadModel stores command response payload fields.
type PersonaGenerationPhaseCommandReadModel struct {
	JobID             int64
	CurrentPhase      string
	PhaseState        string
	PhaseRunID        *int64
	StartedAt         *time.Time
	FinishedAt        *time.Time
	Progress          PersonaGenerationPhaseProgressReadModel
	TargetSummary     PersonaGenerationTargetSummaryReadModel
	Execution         PersonaGenerationExecutionSummaryReadModel
	ResultSummary     *PersonaGenerationPhaseResultSummaryReadModel
	Retryable         bool
	CanStartBodyPhase bool
	ErrorSummary      *PersonaGenerationPhaseErrorSummaryReadModel
}

// PersonaGenerationBodyReadinessInputSummaryReadModel stores downstream input summary fields.
type PersonaGenerationBodyReadinessInputSummaryReadModel struct {
	PersonaCount   int
	MissingCount   int
	SnapshotID     string
	SnapshotDigest string
	EvidenceRefs   []string
}

// PersonaGenerationBodyReadinessReadModel stores downstream readiness fields.
type PersonaGenerationBodyReadinessReadModel struct {
	JobID         int64
	CurrentPhase  string
	PhaseState    string
	Ready         bool
	BlockedReason *string
	ErrorKind     string
	InputSummary  PersonaGenerationBodyReadinessInputSummaryReadModel
}

type personaGenerationTarget struct {
	record  repository.TranslationRecord
	npc     repository.NpcRecord
	profile repository.NpcProfile
	fields  []repository.TranslationField
	hit     bool
}

type personaGenerationTargetSnapshot struct {
	targets        []personaGenerationTarget
	targetCount    int
	commonHits     int
	commonMisses   int
	skippedReasons []string
	digest         string
}

type personaGenerationStartRejection struct {
	kind   string
	reason string
}

type personaGenerationRunSnapshot struct {
	links            []repository.PhaseRunPersona
	personasByID     map[int64]repository.Persona
	linkedByProfile  map[int64]repository.Persona
	totalLinkedCount int
	generatedCount   int
}

type personaGenerationExecutionFailure struct {
	kind      string
	reason    string
	retryable bool
}

// NewPersonaGenerationPhaseService creates the backend persona phase service.
func NewPersonaGenerationPhaseService(
	jobLifecycleRepository personaGenerationPhaseJobLifecycleRepository,
	foundationDataRepository personaGenerationPhaseFoundationDataRepository,
	translationSourceReader personaGenerationPhaseTranslationSourceRepository,
	transactor repository.Transactor,
) *PersonaGenerationPhaseService {
	return &PersonaGenerationPhaseService{
		now:                      func() time.Time { return time.Now().UTC() },
		jobLifecycleRepository:   jobLifecycleRepository,
		foundationDataRepository: foundationDataRepository,
		translationSourceReader:  translationSourceReader,
		transactor:               transactor,
	}
}

// WithPersonaGenerationProvider wires the provider adapter without changing the usecase boundary.
func (service *PersonaGenerationPhaseService) WithPersonaGenerationProvider(
	provider PersonaGenerationProvider,
) *PersonaGenerationPhaseService {
	service.provider = provider
	return service
}

// ReadSummary loads the current persona phase summary for one job.
func (service *PersonaGenerationPhaseService) ReadSummary(
	ctx context.Context,
	jobID int64,
) (PersonaGenerationPhaseSummaryReadModel, error) {
	job, run, termRun, snapshot, err := service.loadContext(ctx, jobID)
	if err != nil {
		return PersonaGenerationPhaseSummaryReadModel{}, err
	}
	state := personaGenerationPhaseStateNotStarted
	var phaseRunID *int64
	var startedAt *time.Time
	var finishedAt *time.Time
	if run != nil {
		state = normalizePersonaGenerationPhaseState(run.State, snapshot.targetCount)
		phaseRunID = &run.ID
		startedAt = clonePersonaTimePointer(run.StartedAt)
		finishedAt = clonePersonaTimePointer(run.FinishedAt)
	}
	progress := service.buildProgress(state, run, snapshot)
	resultSummary, errorSummary, canStartBodyPhase := service.buildResultState(ctx, run, snapshot)
	rejection := service.startRejection(job, run, termRun)
	actionEnablement := service.buildActionEnablement(job.State, run, state, rejection, canStartBodyPhase)
	execution := service.buildExecutionSummary(run, &termRun, snapshot, resultSummary)
	return PersonaGenerationPhaseSummaryReadModel{
		JobID:        job.ID,
		CurrentPhase: personaGenerationCurrentPhase,
		PhaseState:   state,
		PhaseRunID:   clonePersonaInt64Pointer(phaseRunID),
		StartedAt:    startedAt,
		FinishedAt:   finishedAt,
		Progress:     progress,
		TargetSummary: PersonaGenerationTargetSummaryReadModel{
			TargetCount:            snapshot.targetCount,
			CommonPersonaHitCount:  snapshot.commonHits,
			CommonPersonaMissCount: snapshot.commonMisses,
			SkippedCount:           len(snapshot.skippedReasons),
			SkippedReasons:         append([]string(nil), snapshot.skippedReasons...),
			TargetSnapshotDigest:   snapshot.digest,
		},
		Execution:        execution,
		ResultSummary:    resultSummary,
		ErrorSummary:     errorSummary,
		ActionEnablement: actionEnablement,
	}, nil
}

// StartPhase creates or resumes one persona phase run when start conditions pass.
func (service *PersonaGenerationPhaseService) StartPhase(
	ctx context.Context,
	jobID int64,
) (PersonaGenerationPhaseCommandReadModel, error) {
	if service.transactor == nil {
		return PersonaGenerationPhaseCommandReadModel{}, fmt.Errorf("start persona generation phase: transactor is not configured")
	}
	job, run, termRun, snapshot, err := service.loadContext(ctx, jobID)
	if err != nil {
		return PersonaGenerationPhaseCommandReadModel{}, err
	}
	if rejection := service.startRejection(job, run, termRun); rejection != nil {
		response := service.rejectedCommand(job, run, termRun, snapshot, *rejection)
		return response, errPersonaGenerationPhaseExecutionRejected
	}
	updatedJob, updatedRun, err := service.startPhaseRunTransaction(ctx, job, run, termRun, snapshot)
	if err != nil {
		return PersonaGenerationPhaseCommandReadModel{}, fmt.Errorf("start persona generation phase transaction: %w", err)
	}
	updatedRun, err = service.executePhaseRun(ctx, updatedJob, updatedRun, snapshot)
	if err != nil {
		return service.commandFromRun(ctx, updatedJob, updatedRun, snapshot, &termRun), err
	}
	return service.startPhaseCommandReadModel(ctx, updatedJob, updatedRun, snapshot), nil
}

// PausePhase pauses one running persona phase run.
func (service *PersonaGenerationPhaseService) PausePhase(
	ctx context.Context,
	jobID int64,
	phaseRunID int64,
) (PersonaGenerationPhaseCommandReadModel, error) {
	return service.mutatePhaseState(ctx, jobID, phaseRunID, personaGenerationPhaseStateRunning, personaGenerationPhaseStatePaused)
}

// ResumePhase resumes one paused or recoverable-failed persona phase run.
func (service *PersonaGenerationPhaseService) ResumePhase(
	ctx context.Context,
	jobID int64,
	phaseRunID int64,
) (PersonaGenerationPhaseCommandReadModel, error) {
	run, err := service.loadExistingPersonaRun(ctx, jobID, phaseRunID)
	if err != nil {
		return PersonaGenerationPhaseCommandReadModel{}, err
	}
	if run.State != personaGenerationPhaseStatePaused && run.State != personaGenerationPhaseStateRecoverableFail {
		return service.phaseMutationRejected(ctx, jobID, &run, personaGenerationErrorKindTermIncomplete, "persona phase is not resumable"), nil
	}
	command, err := service.mutatePhaseState(ctx, jobID, phaseRunID, run.State, personaGenerationPhaseStateRunning)
	if err != nil {
		return PersonaGenerationPhaseCommandReadModel{}, err
	}
	job, refreshedRun, termRun, snapshot, loadErr := service.loadContext(ctx, jobID)
	if loadErr != nil {
		return PersonaGenerationPhaseCommandReadModel{}, fmt.Errorf("reload persona generation phase after resume: %w", loadErr)
	}
	if refreshedRun == nil {
		return command, nil
	}
	updatedRun, execErr := service.executePhaseRun(ctx, job, *refreshedRun, snapshot)
	if execErr != nil {
		return service.commandFromRun(ctx, job, updatedRun, snapshot, &termRun), execErr
	}
	return service.commandFromRun(ctx, job, updatedRun, snapshot, &termRun), nil
}

// RetryPhase retries one persona phase run or returns its redacted failure.
func (service *PersonaGenerationPhaseService) RetryPhase(
	ctx context.Context,
	jobID int64,
	phaseRunID int64,
) (PersonaGenerationPhaseCommandReadModel, error) {
	run, err := service.loadExistingPersonaRun(ctx, jobID, phaseRunID)
	if err != nil {
		return PersonaGenerationPhaseCommandReadModel{}, err
	}
	job, _, _, snapshot, loadErr := service.loadContext(ctx, jobID)
	if loadErr != nil {
		return PersonaGenerationPhaseCommandReadModel{}, fmt.Errorf("load persona generation phase for retry guard: %w", loadErr)
	}
	if isPersonaGenerationTerminalJob(job.State) {
		return service.phaseMutationRejected(ctx, jobID, &run, personaGenerationErrorKindTerminalJob, personaGenerationTerminalJobReason), nil
	}
	if snapshotDriftReason := personaGenerationSnapshotDriftReason(&run, snapshot); snapshotDriftReason != nil {
		return service.phaseMutationRejected(ctx, jobID, &run, personaGenerationErrorKindSnapshotMissing, *snapshotDriftReason), nil
	}
	if !personaGenerationRetryAllowed(run.State, run.LatestError) {
		return service.phaseMutationRejected(ctx, jobID, &run, personaGenerationErrorKindTermIncomplete, "persona phase is not retryable"), nil
	}
	command, err := service.mutatePhaseState(ctx, jobID, phaseRunID, run.State, personaGenerationPhaseStateRunning)
	if err != nil {
		return PersonaGenerationPhaseCommandReadModel{}, err
	}
	var refreshedRun *repository.JobPhaseRun
	var termRun repository.JobPhaseRun
	job, refreshedRun, termRun, snapshot, loadErr = service.loadContext(ctx, jobID)
	if loadErr != nil {
		return PersonaGenerationPhaseCommandReadModel{}, fmt.Errorf("reload persona generation phase after retry: %w", loadErr)
	}
	if refreshedRun == nil {
		return command, nil
	}
	updatedRun, execErr := service.executePhaseRun(ctx, job, *refreshedRun, snapshot)
	if execErr != nil {
		return service.commandFromRun(ctx, job, updatedRun, snapshot, &termRun), execErr
	}
	return service.commandFromRun(ctx, job, updatedRun, snapshot, &termRun), nil
}

// CancelPhase cancels one persona phase run.
func (service *PersonaGenerationPhaseService) CancelPhase(
	ctx context.Context,
	jobID int64,
	phaseRunID int64,
) (PersonaGenerationPhaseCommandReadModel, error) {
	run, err := service.loadExistingPersonaRun(ctx, jobID, phaseRunID)
	if err != nil {
		return PersonaGenerationPhaseCommandReadModel{}, err
	}
	job, _, _, snapshot, loadErr := service.loadContext(ctx, jobID)
	if loadErr != nil {
		return PersonaGenerationPhaseCommandReadModel{}, fmt.Errorf("load persona generation phase for cancel guard: %w", loadErr)
	}
	if isPersonaGenerationTerminalJob(job.State) {
		return service.phaseMutationRejected(ctx, jobID, &run, personaGenerationErrorKindTerminalJob, personaGenerationTerminalJobReason), nil
	}
	if snapshotDriftReason := personaGenerationSnapshotDriftReason(&run, snapshot); snapshotDriftReason != nil {
		return service.phaseMutationRejected(ctx, jobID, &run, personaGenerationErrorKindSnapshotMissing, *snapshotDriftReason), nil
	}
	if !personaGenerationCancelAllowed(job.State, run.State) {
		return service.phaseMutationRejected(ctx, jobID, &run, personaGenerationErrorKindTermIncomplete, "persona phase is not cancelable"), nil
	}
	return service.mutatePhaseState(ctx, jobID, phaseRunID, "", personaGenerationPhaseStateCanceled)
}

// ReadBodyReadiness reports whether the body phase may start.
func (service *PersonaGenerationPhaseService) ReadBodyReadiness(
	ctx context.Context,
	jobID int64,
) (PersonaGenerationBodyReadinessReadModel, error) {
	job, run, _, snapshot, err := service.loadContext(ctx, jobID)
	if err != nil {
		return PersonaGenerationBodyReadinessReadModel{}, err
	}
	if isPersonaGenerationTerminalJob(job.State) {
		return PersonaGenerationBodyReadinessReadModel{
			JobID:        job.ID,
			CurrentPhase: personaGenerationCurrentPhase,
			PhaseState:   personaGenerationPhaseStateRejected,
			Ready:        false,
			ErrorKind:    personaGenerationErrorKindTerminalJob,
		}, errPersonaGenerationPhaseTerminalJob
	}
	state := personaGenerationPhaseStateNotStarted
	if run != nil {
		state = normalizePersonaGenerationPhaseState(run.State, snapshot.targetCount)
	}
	if run == nil || state != personaGenerationPhaseStateCompleted && state != personaGenerationPhaseStateEmptyCompleted {
		reason := "persona phase completion is required"
		return PersonaGenerationBodyReadinessReadModel{
			JobID:         job.ID,
			CurrentPhase:  personaGenerationCurrentPhase,
			PhaseState:    state,
			Ready:         false,
			BlockedReason: &reason,
			ErrorKind:     personaGenerationErrorKindBodyBlocked,
		}, nil
	}
	if snapshot.targetCount == 0 {
		return PersonaGenerationBodyReadinessReadModel{
			JobID:        job.ID,
			CurrentPhase: personaGenerationCurrentPhase,
			PhaseState:   personaGenerationPhaseStateCompleted,
			Ready:        true,
			InputSummary: PersonaGenerationBodyReadinessInputSummaryReadModel{
				PersonaCount:   0,
				MissingCount:   0,
				SnapshotDigest: snapshot.digest,
			},
		}, nil
	}
	phaseRunPersonas, err := service.jobLifecycleRepository.ListPhaseRunPersonasByPhaseRunID(ctx, run.ID)
	if err != nil {
		return PersonaGenerationBodyReadinessReadModel{}, fmt.Errorf("list phase run personas for body readiness: %w", err)
	}
	evidenceRefs := personaGenerationPhaseRunEvidenceRefs(phaseRunPersonas)
	if len(phaseRunPersonas) == 0 {
		reason := "persona snapshot reference is missing"
		return PersonaGenerationBodyReadinessReadModel{
			JobID:         job.ID,
			CurrentPhase:  personaGenerationCurrentPhase,
			PhaseState:    personaGenerationPhaseStateSnapshotMissing,
			Ready:         false,
			BlockedReason: &reason,
			ErrorKind:     personaGenerationErrorKindSnapshotMissing,
			InputSummary: PersonaGenerationBodyReadinessInputSummaryReadModel{
				MissingCount:   snapshot.targetCount,
				SnapshotDigest: snapshot.digest,
			},
		}, nil
	}
	if snapshotDriftReason := personaGenerationSnapshotDriftReason(run, snapshot); snapshotDriftReason != nil {
		distinctCount := service.distinctPhaseRunPersonaCount(ctx, phaseRunPersonas)
		missingCount := personaMax(0, snapshot.targetCount-distinctCount)
		return PersonaGenerationBodyReadinessReadModel{
			JobID:         job.ID,
			CurrentPhase:  personaGenerationCurrentPhase,
			PhaseState:    personaGenerationPhaseStateSnapshotMissing,
			Ready:         false,
			BlockedReason: snapshotDriftReason,
			ErrorKind:     personaGenerationErrorKindSnapshotMissing,
			InputSummary: PersonaGenerationBodyReadinessInputSummaryReadModel{
				PersonaCount:   distinctCount,
				MissingCount:   missingCount,
				SnapshotID:     personaGenerationSnapshotID(run.ID),
				SnapshotDigest: snapshot.digest,
				EvidenceRefs:   evidenceRefs,
			},
		}, nil
	}
	runSnapshot, snapshotAvailable := service.tryLoadRunSnapshot(ctx, run)
	if !snapshotAvailable {
		reason := "persona snapshot reference is missing"
		return PersonaGenerationBodyReadinessReadModel{
			JobID:         job.ID,
			CurrentPhase:  personaGenerationCurrentPhase,
			PhaseState:    personaGenerationPhaseStateSnapshotMissing,
			Ready:         false,
			BlockedReason: &reason,
			ErrorKind:     personaGenerationErrorKindSnapshotMissing,
			InputSummary: PersonaGenerationBodyReadinessInputSummaryReadModel{
				MissingCount:   snapshot.targetCount,
				SnapshotID:     personaGenerationSnapshotID(run.ID),
				SnapshotDigest: snapshot.digest,
				EvidenceRefs:   evidenceRefs,
			},
		}, nil
	}
	missingCount := personaMax(0, snapshot.targetCount-runSnapshot.totalLinkedCount)
	if missingCount > 0 {
		reason := "persona snapshot reference is not ready"
		return PersonaGenerationBodyReadinessReadModel{
			JobID:         job.ID,
			CurrentPhase:  personaGenerationCurrentPhase,
			PhaseState:    personaGenerationPhaseStateCompleted,
			Ready:         false,
			BlockedReason: &reason,
			ErrorKind:     personaGenerationErrorKindBodyBlocked,
			InputSummary: PersonaGenerationBodyReadinessInputSummaryReadModel{
				PersonaCount:   runSnapshot.totalLinkedCount,
				MissingCount:   missingCount,
				SnapshotID:     personaGenerationSnapshotID(run.ID),
				SnapshotDigest: snapshot.digest,
				EvidenceRefs:   evidenceRefs,
			},
		}, nil
	}
	return PersonaGenerationBodyReadinessReadModel{
		JobID:        job.ID,
		CurrentPhase: personaGenerationCurrentPhase,
		PhaseState:   personaGenerationPhaseStateCompleted,
		Ready:        true,
		InputSummary: PersonaGenerationBodyReadinessInputSummaryReadModel{
			PersonaCount:   runSnapshot.totalLinkedCount,
			MissingCount:   0,
			SnapshotID:     personaGenerationSnapshotID(run.ID),
			SnapshotDigest: snapshot.digest,
			EvidenceRefs:   evidenceRefs,
		},
	}, nil
}

func (service *PersonaGenerationPhaseService) mutatePhaseState(
	ctx context.Context,
	jobID int64,
	phaseRunID int64,
	requiredState string,
	nextState string,
) (PersonaGenerationPhaseCommandReadModel, error) {
	if service.transactor == nil {
		return PersonaGenerationPhaseCommandReadModel{}, fmt.Errorf("mutate persona generation phase: transactor is not configured")
	}
	job, run, _, snapshot, err := service.loadContext(ctx, jobID)
	if err != nil {
		return PersonaGenerationPhaseCommandReadModel{}, fmt.Errorf("load persona generation phase for mutation: %w", err)
	}
	if run == nil || run.ID != phaseRunID {
		return PersonaGenerationPhaseCommandReadModel{}, fmt.Errorf(personaGenerationRunLoadErrorMessage, repository.ErrNotFound)
	}
	if isPersonaGenerationTerminalJob(job.State) {
		return service.phaseMutationRejected(ctx, jobID, run, personaGenerationErrorKindTerminalJob, personaGenerationTerminalJobReason), nil
	}
	if snapshotDriftReason := personaGenerationSnapshotDriftReason(run, snapshot); snapshotDriftReason != nil {
		return service.phaseMutationRejected(ctx, jobID, run, personaGenerationErrorKindSnapshotMissing, *snapshotDriftReason), nil
	}
	if requiredState != "" && run.State != requiredState {
		return service.phaseMutationRejected(ctx, jobID, run, personaGenerationErrorKindTermIncomplete, "persona phase state does not allow this command"), nil
	}
	var updatedRun repository.JobPhaseRun
	err = service.transactor.WithTransaction(ctx, func(txCtx context.Context) error {
		nextRun, updateErr := service.jobLifecycleRepository.UpdateJobPhaseRun(txCtx, run.ID, repository.JobPhaseRunUpdateDraft{
			State:               nextState,
			ProgressPercent:     run.ProgressPercent,
			LatestExternalRunID: run.LatestExternalRunID,
			LatestError:         run.LatestError,
			StartedAt:           run.StartedAt,
			FinishedAt:          run.FinishedAt,
		})
		if updateErr != nil {
			return fmt.Errorf("update persona generation phase run state: %w", updateErr)
		}
		updatedRun = nextRun
		return nil
	})
	if err != nil {
		return PersonaGenerationPhaseCommandReadModel{}, fmt.Errorf("update persona generation phase state: %w", err)
	}
	state := normalizePersonaGenerationPhaseState(updatedRun.State, snapshot.targetCount)
	progress := service.buildProgress(state, &updatedRun, snapshot)
	resultSummary, _, canStartBodyPhase := service.buildResultState(ctx, &updatedRun, snapshot)
	return PersonaGenerationPhaseCommandReadModel{
		JobID:        job.ID,
		CurrentPhase: personaGenerationCurrentPhase,
		PhaseState:   state,
		PhaseRunID:   clonePersonaInt64Pointer(&updatedRun.ID),
		StartedAt:    clonePersonaTimePointer(updatedRun.StartedAt),
		FinishedAt:   clonePersonaTimePointer(updatedRun.FinishedAt),
		Progress:     progress,
		TargetSummary: PersonaGenerationTargetSummaryReadModel{
			TargetCount:            snapshot.targetCount,
			CommonPersonaHitCount:  snapshot.commonHits,
			CommonPersonaMissCount: snapshot.commonMisses,
			SkippedCount:           len(snapshot.skippedReasons),
			SkippedReasons:         append([]string(nil), snapshot.skippedReasons...),
			TargetSnapshotDigest:   snapshot.digest,
		},
		Execution:         service.buildExecutionSummary(&updatedRun, &updatedRun, snapshot, resultSummary),
		ResultSummary:     resultSummary,
		Retryable:         personaGenerationRetryAllowed(updatedRun.State, updatedRun.LatestError),
		CanStartBodyPhase: canStartBodyPhase,
	}, nil
}

func (service *PersonaGenerationPhaseService) phaseMutationRejected(
	ctx context.Context,
	jobID int64,
	run *repository.JobPhaseRun,
	errorKind string,
	reason string,
) PersonaGenerationPhaseCommandReadModel {
	job, _, termRun, snapshot, loadErr := service.loadContext(ctx, jobID)
	if loadErr != nil {
		return PersonaGenerationPhaseCommandReadModel{
			JobID:        jobID,
			CurrentPhase: personaGenerationCurrentPhase,
			PhaseState:   personaGenerationPhaseStateRejected,
			ErrorSummary: &PersonaGenerationPhaseErrorSummaryReadModel{
				ErrorKind:  errorKind,
				Reason:     reason,
				Retryable:  false,
				IsRedacted: true,
			},
		}
	}
	return PersonaGenerationPhaseCommandReadModel{
		JobID:        job.ID,
		CurrentPhase: personaGenerationCurrentPhase,
		PhaseState:   personaGenerationPhaseStateRejected,
		PhaseRunID:   clonePersonaInt64Pointer(&run.ID),
		StartedAt:    clonePersonaTimePointer(run.StartedAt),
		FinishedAt:   clonePersonaTimePointer(run.FinishedAt),
		Progress:     service.buildProgress(personaGenerationPhaseStateRejected, run, snapshot),
		TargetSummary: PersonaGenerationTargetSummaryReadModel{
			TargetCount:            snapshot.targetCount,
			CommonPersonaHitCount:  snapshot.commonHits,
			CommonPersonaMissCount: snapshot.commonMisses,
			SkippedCount:           len(snapshot.skippedReasons),
			SkippedReasons:         append([]string(nil), snapshot.skippedReasons...),
			TargetSnapshotDigest:   snapshot.digest,
		},
		Execution: service.buildExecutionSummary(run, &termRun, snapshot, nil),
		ErrorSummary: &PersonaGenerationPhaseErrorSummaryReadModel{
			ErrorKind:  normalizePersonaField(errorKind),
			Reason:     reason,
			Retryable:  false,
			IsRedacted: true,
		},
	}
}

func (service *PersonaGenerationPhaseService) loadContext(
	ctx context.Context,
	jobID int64,
) (repository.TranslationJob, *repository.JobPhaseRun, repository.JobPhaseRun, personaGenerationTargetSnapshot, error) {
	if service.jobLifecycleRepository == nil || service.foundationDataRepository == nil || service.translationSourceReader == nil {
		return repository.TranslationJob{}, nil, repository.JobPhaseRun{}, personaGenerationTargetSnapshot{}, fmt.Errorf("persona generation phase persistence is not configured")
	}
	job, err := service.jobLifecycleRepository.GetTranslationJobByID(ctx, jobID)
	if err != nil {
		return repository.TranslationJob{}, nil, repository.JobPhaseRun{}, personaGenerationTargetSnapshot{}, fmt.Errorf("load translation job: %w", err)
	}
	inputMetadata, inputMetadataErr := service.translationSourceReader.GetXEditExtractedDataByID(ctx, job.XEditExtractedDataID)
	if inputMetadataErr != nil {
		return repository.TranslationJob{}, nil, repository.JobPhaseRun{}, personaGenerationTargetSnapshot{}, fmt.Errorf("load translation input metadata: %w", inputMetadataErr)
	}
	_ = inputMetadata
	phases, err := service.jobLifecycleRepository.ListJobPhaseRunsByJobID(ctx, job.ID)
	if err != nil {
		return repository.TranslationJob{}, nil, repository.JobPhaseRun{}, personaGenerationTargetSnapshot{}, fmt.Errorf("list translation job phases: %w", err)
	}
	termRun, err := personaGenerationTermPhase(phases)
	if err != nil {
		return repository.TranslationJob{}, nil, repository.JobPhaseRun{}, personaGenerationTargetSnapshot{}, err
	}
	var run *repository.JobPhaseRun
	foundRun, err := service.jobLifecycleRepository.FindJobPhaseRun(ctx, job.ID, personaGenerationPhaseType)
	if err == nil {
		run = &foundRun
	} else if !errors.Is(err, repository.ErrNotFound) {
		return repository.TranslationJob{}, nil, repository.JobPhaseRun{}, personaGenerationTargetSnapshot{}, fmt.Errorf("find persona generation phase run: %w", err)
	}
	snapshot, err := service.buildTargetSnapshot(ctx, job.XEditExtractedDataID)
	if err != nil {
		return repository.TranslationJob{}, nil, repository.JobPhaseRun{}, personaGenerationTargetSnapshot{}, fmt.Errorf("build persona generation target snapshot: %w", err)
	}
	return job, run, termRun, snapshot, nil
}

func (service *PersonaGenerationPhaseService) buildTargetSnapshot(
	ctx context.Context,
	xEditID int64,
) (personaGenerationTargetSnapshot, error) {
	records, err := service.translationSourceReader.ListTranslationRecordsByXEditID(ctx, xEditID)
	if err != nil {
		return personaGenerationTargetSnapshot{}, fmt.Errorf("list translation records for persona targets: %w", err)
	}
	result := personaGenerationTargetSnapshot{
		targets:        make([]personaGenerationTarget, 0),
		skippedReasons: make([]string, 0),
	}
	for _, record := range records {
		if strings.TrimSpace(record.RecordType) != "NPC_" {
			continue
		}
		target, skippedReason, err := service.buildSnapshotTarget(ctx, record)
		if skippedReason != "" {
			result.skippedReasons = append(result.skippedReasons, skippedReason)
			continue
		}
		if err != nil {
			return personaGenerationTargetSnapshot{}, err
		}
		if target.hit {
			result.commonHits++
		} else {
			result.commonMisses++
		}
		result.targets = append(result.targets, target)
	}
	sort.Strings(result.skippedReasons)
	result.targetCount = len(result.targets)
	result.digest = buildPersonaGenerationTargetSnapshotDigest(result)
	return result, nil
}

func (service *PersonaGenerationPhaseService) buildProgress(
	state string,
	run *repository.JobPhaseRun,
	snapshot personaGenerationTargetSnapshot,
) PersonaGenerationPhaseProgressReadModel {
	processedCount := 0
	percent := 0
	currentStep := "not_started"
	if run != nil {
		percent = personaMax(0, run.ProgressPercent)
		currentStep = normalizePersonaField(run.State)
		if currentStep == "" {
			currentStep = state
		}
	}
	switch state {
	case personaGenerationPhaseStateCompleted, personaGenerationPhaseStateEmptyCompleted:
		processedCount = snapshot.targetCount
		percent = 100
		currentStep = "completed"
	case personaGenerationPhaseStateRunning:
		processedCount = personaMin(snapshot.targetCount, percentageCount(percent, snapshot.targetCount))
	case personaGenerationPhaseStatePaused:
		processedCount = personaMin(snapshot.targetCount, percentageCount(percent, snapshot.targetCount))
		currentStep = "paused"
	case personaGenerationPhaseStateRejected:
		currentStep = "rejected"
	case personaGenerationPhaseStateNotStarted:
		currentStep = "not_started"
	}
	return PersonaGenerationPhaseProgressReadModel{
		Percent:        percent,
		ProcessedCount: processedCount,
		TotalCount:     snapshot.targetCount,
		TargetCount:    snapshot.targetCount,
		CurrentStep:    currentStep,
	}
}

func (service *PersonaGenerationPhaseService) buildResultState(
	ctx context.Context,
	run *repository.JobPhaseRun,
	snapshot personaGenerationTargetSnapshot,
) (*PersonaGenerationPhaseResultSummaryReadModel, *PersonaGenerationPhaseErrorSummaryReadModel, bool) {
	if run == nil {
		return nil, nil, false
	}
	state := normalizePersonaGenerationPhaseState(run.State, snapshot.targetCount)
	if snapshot.targetCount == 0 && state == personaGenerationPhaseStateCompleted {
		return &PersonaGenerationPhaseResultSummaryReadModel{
			SnapshotDigest:          snapshot.digest,
			SnapshotReferenceStatus: "empty",
			BodyReadiness:           true,
		}, nil, true
	}
	runSnapshot, snapshotErr := service.loadRunSnapshot(ctx, run)
	if snapshotErr != nil {
		return &PersonaGenerationPhaseResultSummaryReadModel{
				MissingCount:            snapshot.targetCount,
				SnapshotDigest:          snapshot.digest,
				SnapshotReferenceStatus: "missing",
				BodyReadiness:           false,
			}, &PersonaGenerationPhaseErrorSummaryReadModel{
				ErrorKind:  personaGenerationErrorKindSnapshotMissing,
				Reason:     personaGenerationRedactedPublicSummary,
				Retryable:  false,
				IsRedacted: true,
			}, false
	}
	if snapshotDriftReason := personaGenerationSnapshotDriftReason(run, snapshot); snapshotDriftReason != nil {
		return &PersonaGenerationPhaseResultSummaryReadModel{
				PersonaCount:            runSnapshot.totalLinkedCount,
				MissingCount:            personaMax(0, snapshot.targetCount-runSnapshot.totalLinkedCount),
				SnapshotDigest:          snapshot.digest,
				SnapshotReferenceStatus: "missing",
				BodyReadiness:           false,
			}, &PersonaGenerationPhaseErrorSummaryReadModel{
				ErrorKind:  personaGenerationErrorKindSnapshotMissing,
				Reason:     *snapshotDriftReason,
				Retryable:  false,
				IsRedacted: true,
			}, false
	}
	missingCount := personaMax(0, snapshot.targetCount-runSnapshot.totalLinkedCount)
	if state == personaGenerationPhaseStateCompleted {
		if runSnapshot.totalLinkedCount > 0 || snapshot.targetCount == 0 {
			ready := missingCount == 0
			return &PersonaGenerationPhaseResultSummaryReadModel{
				GeneratedCount:          runSnapshot.generatedCount,
				PersonaCount:            runSnapshot.totalLinkedCount,
				MissingCount:            missingCount,
				SnapshotID:              personaGenerationSnapshotID(run.ID),
				SnapshotDigest:          snapshot.digest,
				SnapshotReferenceStatus: "available",
				BodyReadiness:           ready,
			}, nil, ready
		}
		return &PersonaGenerationPhaseResultSummaryReadModel{
				MissingCount:            snapshot.targetCount,
				SnapshotDigest:          snapshot.digest,
				SnapshotReferenceStatus: "missing",
				BodyReadiness:           false,
			}, &PersonaGenerationPhaseErrorSummaryReadModel{
				ErrorKind:  personaGenerationErrorKindSnapshotMissing,
				Reason:     personaGenerationRedactedPublicSummary,
				Retryable:  false,
				IsRedacted: true,
			}, false
	}
	if state == personaGenerationPhaseStateFailed || state == personaGenerationPhaseStateRecoverableFail {
		return &PersonaGenerationPhaseResultSummaryReadModel{
				GeneratedCount:          runSnapshot.generatedCount,
				FailedCount:             personaMax(1, missingCount),
				PersonaCount:            runSnapshot.totalLinkedCount,
				MissingCount:            missingCount,
				SnapshotDigest:          snapshot.digest,
				SnapshotReferenceStatus: "missing",
				BodyReadiness:           false,
			}, &PersonaGenerationPhaseErrorSummaryReadModel{
				ErrorKind:  normalizePersonaField(run.LatestError),
				Reason:     personaGenerationRedactedPublicSummary,
				Retryable:  true,
				IsRedacted: true,
			}, false
	}
	return &PersonaGenerationPhaseResultSummaryReadModel{
		GeneratedCount:          runSnapshot.generatedCount,
		PersonaCount:            runSnapshot.totalLinkedCount,
		MissingCount:            missingCount,
		SnapshotDigest:          snapshot.digest,
		SnapshotReferenceStatus: "pending",
		BodyReadiness:           false,
	}, nil, false
}

func (service *PersonaGenerationPhaseService) buildActionEnablement(
	jobState string,
	run *repository.JobPhaseRun,
	state string,
	rejection *personaGenerationStartRejection,
	canStartBodyPhase bool,
) PersonaGenerationPhaseActionEnablementReadModel {
	var startBlockedReason *string
	if rejection != nil {
		startBlockedReason = stringPersonaPointer(rejection.reason)
	}
	canPause := !isPersonaGenerationTerminalJob(jobState) && state == personaGenerationPhaseStateRunning
	canResume := !isPersonaGenerationTerminalJob(jobState) && (state == personaGenerationPhaseStatePaused || state == personaGenerationPhaseStateRecoverableFail)
	canRetry := run != nil && !isPersonaGenerationTerminalJob(jobState) && personaGenerationRetryAllowed(run.State, run.LatestError)
	canCancel := run != nil && personaGenerationCancelAllowed(jobState, run.State)
	bodyReason := stringPersonaPointer("persona snapshot reference is not ready")
	if canStartBodyPhase {
		bodyReason = nil
	}
	return PersonaGenerationPhaseActionEnablementReadModel{
		CanStart:               rejection == nil,
		StartBlockedReason:     clonePersonaStringPointer(startBlockedReason),
		CanPause:               canPause,
		PauseBlockedReason:     personaBlockedReason(canPause, "persona phase is not running"),
		CanResume:              canResume,
		ResumeBlockedReason:    personaBlockedReason(canResume, "persona phase is not resumable"),
		CanRetry:               canRetry,
		RetryBlockedReason:     personaBlockedReason(canRetry, "persona phase is not retryable"),
		CanCancel:              canCancel,
		CancelBlockedReason:    personaBlockedReason(canCancel, "persona phase is not cancelable"),
		CanStartBodyPhase:      canStartBodyPhase,
		BodyPhaseBlockedReason: clonePersonaStringPointer(bodyReason),
	}
}

func (service *PersonaGenerationPhaseService) startRejection(
	job repository.TranslationJob,
	run *repository.JobPhaseRun,
	termRun repository.JobPhaseRun,
) *personaGenerationStartRejection {
	if isPersonaGenerationTerminalJob(job.State) {
		return &personaGenerationStartRejection{
			kind:   personaGenerationErrorKindTerminalJob,
			reason: personaGenerationTerminalJobReason,
		}
	}
	if normalizePersonaField(termRun.State) != personaGenerationPhaseStateCompleted {
		return &personaGenerationStartRejection{
			kind:   personaGenerationErrorKindTermIncomplete,
			reason: "term phase must be completed",
		}
	}
	if run != nil && isPersonaGenerationActiveState(run.State) {
		return &personaGenerationStartRejection{
			kind:   personaGenerationErrorKindActivePhase,
			reason: "active persona phase run already exists",
		}
	}
	return nil
}

func (service *PersonaGenerationPhaseService) rejectedCommand(
	job repository.TranslationJob,
	run *repository.JobPhaseRun,
	termRun repository.JobPhaseRun,
	snapshot personaGenerationTargetSnapshot,
	rejection personaGenerationStartRejection,
) PersonaGenerationPhaseCommandReadModel {
	var phaseRunID *int64
	var startedAt *time.Time
	var finishedAt *time.Time
	if run != nil {
		phaseRunID = &run.ID
		startedAt = clonePersonaTimePointer(run.StartedAt)
		finishedAt = clonePersonaTimePointer(run.FinishedAt)
	}
	return PersonaGenerationPhaseCommandReadModel{
		JobID:        job.ID,
		CurrentPhase: personaGenerationCurrentPhase,
		PhaseState:   personaGenerationPhaseStateRejected,
		PhaseRunID:   clonePersonaInt64Pointer(phaseRunID),
		StartedAt:    startedAt,
		FinishedAt:   finishedAt,
		Progress: PersonaGenerationPhaseProgressReadModel{
			Percent:     0,
			TotalCount:  snapshot.targetCount,
			TargetCount: snapshot.targetCount,
			CurrentStep: "rejected",
		},
		TargetSummary: PersonaGenerationTargetSummaryReadModel{
			TargetCount:            snapshot.targetCount,
			CommonPersonaHitCount:  snapshot.commonHits,
			CommonPersonaMissCount: snapshot.commonMisses,
			SkippedCount:           len(snapshot.skippedReasons),
			SkippedReasons:         append([]string(nil), snapshot.skippedReasons...),
			TargetSnapshotDigest:   snapshot.digest,
		},
		Execution: PersonaGenerationExecutionSummaryReadModel{
			CredentialRef: strings.TrimSpace(termRun.CredentialRef),
			Provider:      strings.TrimSpace(termRun.AIProvider),
			Model:         strings.TrimSpace(termRun.ModelName),
			ExecutionMode: strings.TrimSpace(termRun.ExecutionMode),
			PromptDigest:  service.aggregatePromptDigest(run, snapshot),
			InputCount:    snapshot.targetCount,
		},
		CanStartBodyPhase: false,
		ErrorSummary: &PersonaGenerationPhaseErrorSummaryReadModel{
			ErrorKind:  rejection.kind,
			Reason:     rejection.reason,
			Retryable:  rejection.kind != personaGenerationErrorKindTerminalJob,
			IsRedacted: true,
		},
	}
}

func (service *PersonaGenerationPhaseService) loadExistingPersonaRun(
	ctx context.Context,
	jobID int64,
	phaseRunID int64,
) (repository.JobPhaseRun, error) {
	run, err := service.jobLifecycleRepository.FindJobPhaseRun(ctx, jobID, personaGenerationPhaseType)
	if err != nil {
		return repository.JobPhaseRun{}, fmt.Errorf(personaGenerationRunLoadErrorMessage, err)
	}
	if run.ID != phaseRunID {
		return repository.JobPhaseRun{}, fmt.Errorf(personaGenerationRunLoadErrorMessage, repository.ErrNotFound)
	}
	return run, nil
}

func (service *PersonaGenerationPhaseService) executePhaseRun(
	ctx context.Context,
	job repository.TranslationJob,
	run repository.JobPhaseRun,
	snapshot personaGenerationTargetSnapshot,
) (repository.JobPhaseRun, error) {
	if snapshot.targetCount == 0 {
		return run, nil
	}
	if isPersonaGenerationTerminalJob(job.State) {
		return run, errPersonaGenerationPhaseTerminalJob
	}
	runSnapshot, err := service.loadRunSnapshot(ctx, &run)
	if err != nil {
		return run, fmt.Errorf("load persona generation phase run snapshot: %w", err)
	}
	failures := make([]personaGenerationExecutionFailure, 0)
	for _, target := range snapshot.targets {
		nextSnapshot, failure, executeErr := service.executePhaseTarget(ctx, job, run, runSnapshot, target)
		if executeErr != nil {
			return run, executeErr
		}
		if failure != nil {
			failures = append(failures, *failure)
			continue
		}
		runSnapshot = nextSnapshot
	}
	return service.finishPhaseRun(ctx, run, snapshot, runSnapshot, failures)
}

func (service *PersonaGenerationPhaseService) startPhaseRunTransaction(
	ctx context.Context,
	job repository.TranslationJob,
	run *repository.JobPhaseRun,
	termRun repository.JobPhaseRun,
	snapshot personaGenerationTargetSnapshot,
) (repository.TranslationJob, repository.JobPhaseRun, error) {
	var updatedJob repository.TranslationJob
	var updatedRun repository.JobPhaseRun
	err := service.transactor.WithTransaction(ctx, func(txCtx context.Context) error {
		now := service.now()
		startedAt := personaGenerationStartedAt(job.StartedAt, now)
		runState, progressPercent, finishedAt := personaGenerationInitialRunStatus(snapshot, now)
		nextJob, updateErr := service.jobLifecycleRepository.UpdateTranslationJob(txCtx, job.ID, repository.TranslationJobUpdateDraft{
			JobName:         job.JobName,
			State:           personaGenerationJobStateRunning,
			ProgressPercent: 0,
			StartedAt:       startedAt,
			FinishedAt:      nil,
		})
		if updateErr != nil {
			return fmt.Errorf("update translation job for persona phase start: %w", updateErr)
		}
		updatedJob = nextJob
		currentRun, ensureErr := service.ensurePersonaPhaseRun(txCtx, job.ID, run, termRun, runState)
		if ensureErr != nil {
			return ensureErr
		}
		nextRun, updateErr := service.jobLifecycleRepository.UpdateJobPhaseRun(txCtx, currentRun.ID, repository.JobPhaseRunUpdateDraft{
			State:               runState,
			ProgressPercent:     progressPercent,
			LatestExternalRunID: snapshot.digest,
			LatestError:         "",
			StartedAt:           startedAt,
			FinishedAt:          finishedAt,
		})
		if updateErr != nil {
			return fmt.Errorf("update persona generation phase run: %w", updateErr)
		}
		updatedRun = nextRun
		return nil
	})
	if err != nil {
		return repository.TranslationJob{}, repository.JobPhaseRun{}, fmt.Errorf("start persona generation phase transaction: %w", err)
	}
	return updatedJob, updatedRun, nil
}

func (service *PersonaGenerationPhaseService) startPhaseCommandReadModel(
	ctx context.Context,
	job repository.TranslationJob,
	run repository.JobPhaseRun,
	snapshot personaGenerationTargetSnapshot,
) PersonaGenerationPhaseCommandReadModel {
	state := normalizePersonaGenerationPhaseState(run.State, snapshot.targetCount)
	progress := service.buildProgress(state, &run, snapshot)
	resultSummary, _, canStartBodyPhase := service.buildResultState(ctx, &run, snapshot)
	return PersonaGenerationPhaseCommandReadModel{
		JobID:        job.ID,
		CurrentPhase: personaGenerationCurrentPhase,
		PhaseState:   state,
		PhaseRunID:   clonePersonaInt64Pointer(&run.ID),
		StartedAt:    clonePersonaTimePointer(run.StartedAt),
		FinishedAt:   clonePersonaTimePointer(run.FinishedAt),
		Progress:     progress,
		TargetSummary: PersonaGenerationTargetSummaryReadModel{
			TargetCount:            snapshot.targetCount,
			CommonPersonaHitCount:  snapshot.commonHits,
			CommonPersonaMissCount: snapshot.commonMisses,
			SkippedCount:           len(snapshot.skippedReasons),
			SkippedReasons:         append([]string(nil), snapshot.skippedReasons...),
			TargetSnapshotID:       personaGenerationTargetSnapshotIDPointer(run.ID),
			TargetSnapshotDigest:   snapshot.digest,
		},
		Execution:         service.buildExecutionSummary(&run, &run, snapshot, resultSummary),
		ResultSummary:     resultSummary,
		Retryable:         false,
		CanStartBodyPhase: canStartBodyPhase,
	}
}

func personaGenerationStartedAt(existing *time.Time, now time.Time) *time.Time {
	startedAt := clonePersonaTimePointer(existing)
	if startedAt != nil {
		return startedAt
	}
	return &now
}

func personaGenerationInitialRunStatus(snapshot personaGenerationTargetSnapshot, now time.Time) (string, int, *time.Time) {
	if snapshot.targetCount == 0 {
		return personaGenerationPhaseStateCompleted, 100, &now
	}
	return personaGenerationPhaseStateRunning, 0, nil
}

func (service *PersonaGenerationPhaseService) ensurePersonaPhaseRun(
	ctx context.Context,
	jobID int64,
	run *repository.JobPhaseRun,
	termRun repository.JobPhaseRun,
	runState string,
) (repository.JobPhaseRun, error) {
	if run != nil {
		return *run, nil
	}
	phases, listErr := service.jobLifecycleRepository.ListJobPhaseRunsByJobID(ctx, jobID)
	if listErr != nil {
		return repository.JobPhaseRun{}, fmt.Errorf("list translation job phases for persona phase create: %w", listErr)
	}
	created, createErr := service.jobLifecycleRepository.CreateJobPhaseRun(ctx, repository.JobPhaseRunDraft{
		TranslationJobID: jobID,
		PhaseType:        personaGenerationPhaseType,
		State:            runState,
		ExecutionOrder:   len(phases) + 1,
		AIProvider:       termRun.AIProvider,
		ModelName:        termRun.ModelName,
		ExecutionMode:    termRun.ExecutionMode,
		CredentialRef:    termRun.CredentialRef,
		InstructionKind:  personaGenerationInstructionKindDefault,
	})
	if createErr != nil {
		return repository.JobPhaseRun{}, fmt.Errorf("create persona generation phase run: %w", createErr)
	}
	return created, nil
}

func personaGenerationTargetSnapshotIDPointer(runID int64) *string {
	if runID <= 0 {
		return nil
	}
	targetSnapshotID := fmt.Sprintf("persona-target-%d", runID)
	return clonePersonaStringPointer(&targetSnapshotID)
}

func personaGenerationSnapshotID(runID int64) string {
	return fmt.Sprintf("persona-snapshot-%d", runID)
}

func (service *PersonaGenerationPhaseService) buildSnapshotTarget(
	ctx context.Context,
	record repository.TranslationRecord,
) (personaGenerationTarget, string, error) {
	npcRecord, profile, skippedReason, err := service.loadSnapshotTargetContext(ctx, record)
	if skippedReason != "" || err != nil {
		return personaGenerationTarget{}, skippedReason, err
	}
	utteranceFields, skippedReason, err := service.loadSnapshotTargetFields(ctx, record.ID)
	if skippedReason != "" || err != nil {
		return personaGenerationTarget{}, skippedReason, err
	}
	target := personaGenerationTarget{
		record:  record,
		npc:     npcRecord,
		profile: profile,
		fields:  append([]repository.TranslationField(nil), utteranceFields...),
	}
	hit, hitErr := service.snapshotTargetHasCommonPersona(ctx, profile.ID)
	if hitErr != nil {
		return personaGenerationTarget{}, "", hitErr
	}
	target.hit = hit
	return target, "", nil
}

func (service *PersonaGenerationPhaseService) loadSnapshotTargetContext(
	ctx context.Context,
	record repository.TranslationRecord,
) (repository.NpcRecord, repository.NpcProfile, string, error) {
	npcRecord, err := service.translationSourceReader.GetNpcRecordByTranslationRecordID(ctx, record.ID)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return repository.NpcRecord{}, repository.NpcProfile{}, personaGenerationSkippedReasonOrphan, nil
		}
		return repository.NpcRecord{}, repository.NpcProfile{}, "", fmt.Errorf("load npc record %d: %w", record.ID, err)
	}
	profile, err := service.translationSourceReader.GetNpcProfileByID(ctx, npcRecord.NpcProfileID)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return repository.NpcRecord{}, repository.NpcProfile{}, personaGenerationSkippedReasonOrphan, nil
		}
		return repository.NpcRecord{}, repository.NpcProfile{}, "", fmt.Errorf("load npc profile %d: %w", npcRecord.NpcProfileID, err)
	}
	return npcRecord, profile, "", nil
}

func (service *PersonaGenerationPhaseService) loadSnapshotTargetFields(
	ctx context.Context,
	recordID int64,
) ([]repository.TranslationField, string, error) {
	fields, err := service.translationSourceReader.ListTranslationFieldsByTranslationRecordID(ctx, recordID)
	if err != nil {
		return nil, "", fmt.Errorf("list translation fields for persona target %d: %w", recordID, err)
	}
	utteranceFields := make([]repository.TranslationField, 0, len(fields))
	for _, field := range fields {
		includeField, skippedReason, fieldErr := service.snapshotTargetUtteranceField(ctx, field)
		if skippedReason != "" || fieldErr != nil {
			return nil, skippedReason, fieldErr
		}
		if includeField {
			utteranceFields = append(utteranceFields, field)
		}
	}
	if len(utteranceFields) == 0 {
		return nil, personaGenerationSkippedReasonNoUtterance, nil
	}
	return utteranceFields, "", nil
}

func (service *PersonaGenerationPhaseService) snapshotTargetUtteranceField(
	ctx context.Context,
	field repository.TranslationField,
) (bool, string, error) {
	if strings.TrimSpace(field.SourceText) == "" {
		return false, "", nil
	}
	references, err := service.translationSourceReader.ListTranslationFieldRecordReferencesByFieldID(ctx, field.ID)
	if err != nil {
		return false, "", fmt.Errorf("list translation field references for persona target %d: %w", field.ID, err)
	}
	for _, ref := range references {
		if _, referencedErr := service.translationSourceReader.GetTranslationRecordByID(ctx, ref.ReferencedTranslationRecordID); referencedErr != nil {
			if errors.Is(referencedErr, repository.ErrNotFound) {
				return false, personaGenerationSkippedReasonOrphan, nil
			}
			return false, "", fmt.Errorf("load referenced translation record %d: %w", ref.ReferencedTranslationRecordID, referencedErr)
		}
	}
	return true, "", nil
}

func (service *PersonaGenerationPhaseService) snapshotTargetHasCommonPersona(
	ctx context.Context,
	profileID int64,
) (bool, error) {
	persona, err := service.foundationDataRepository.GetPersonaByNpcProfileID(ctx, profileID)
	if err == nil {
		return persona.TranslationJobID == nil, nil
	}
	if errors.Is(err, repository.ErrNotFound) {
		return false, nil
	}
	return false, fmt.Errorf("load common persona by npc profile %d: %w", profileID, err)
}

func (service *PersonaGenerationPhaseService) executePhaseTarget(
	ctx context.Context,
	job repository.TranslationJob,
	run repository.JobPhaseRun,
	runSnapshot personaGenerationRunSnapshot,
	target personaGenerationTarget,
) (personaGenerationRunSnapshot, *personaGenerationExecutionFailure, error) {
	if _, alreadyLinked := runSnapshot.linkedByProfile[target.profile.ID]; alreadyLinked {
		return runSnapshot, nil, nil
	}
	if target.hit {
		return service.executeCommonPersonaTarget(ctx, job, run, target)
	}
	return service.executeGeneratedPersonaTarget(ctx, job, run, target)
}

func (service *PersonaGenerationPhaseService) executeCommonPersonaTarget(
	ctx context.Context,
	job repository.TranslationJob,
	run repository.JobPhaseRun,
	target personaGenerationTarget,
) (personaGenerationRunSnapshot, *personaGenerationExecutionFailure, error) {
	if !service.persistCommonPersonaTarget(ctx, job, run, target) {
		return personaGenerationRunSnapshot{}, personaGenerationSaveFailedFailure(), nil
	}
	runSnapshot, err := service.loadRunSnapshot(ctx, &run)
	if err != nil {
		return personaGenerationRunSnapshot{}, nil, fmt.Errorf("reload persona generation phase run snapshot after common link: %w", err)
	}
	return runSnapshot, nil, nil
}

func (service *PersonaGenerationPhaseService) executeGeneratedPersonaTarget(
	ctx context.Context,
	job repository.TranslationJob,
	run repository.JobPhaseRun,
	target personaGenerationTarget,
) (personaGenerationRunSnapshot, *personaGenerationExecutionFailure, error) {
	request, failure := service.buildProviderRequestForExecution(run, target)
	if failure != nil {
		return personaGenerationRunSnapshot{}, failure, nil
	}
	if service.provider == nil {
		failure := personaGenerationExecutionFailure{
			kind:      personaGenerationErrorKindProviderFailure,
			reason:    personaGenerationRedactedPublicSummary,
			retryable: true,
		}
		return personaGenerationRunSnapshot{}, &failure, nil
	}
	result := service.provider.GeneratePersona(ctx, request)
	if result.Failure != nil {
		failure := personaGenerationExecutionFailure{
			kind:      normalizePersonaField(string(result.Failure.Kind)),
			reason:    personaGenerationRedactedPublicSummary,
			retryable: result.Failure.Retryable,
		}
		return personaGenerationRunSnapshot{}, &failure, nil
	}
	if !personaGenerationProviderResultMatchesRequest(result, request) {
		failure := personaGenerationExecutionFailure{
			kind:      personaGenerationErrorKindInvalidProvider,
			reason:    personaGenerationRedactedPublicSummary,
			retryable: true,
		}
		return personaGenerationRunSnapshot{}, &failure, nil
	}
	if !service.persistGeneratedPersonaTarget(ctx, job, run, target, result) {
		return personaGenerationRunSnapshot{}, personaGenerationSaveFailedFailure(), nil
	}
	runSnapshot, err := service.loadRunSnapshot(ctx, &run)
	if err != nil {
		return personaGenerationRunSnapshot{}, nil, fmt.Errorf("reload persona generation phase run snapshot after persona save: %w", err)
	}
	return runSnapshot, nil, nil
}

func (service *PersonaGenerationPhaseService) persistCommonPersonaTarget(
	ctx context.Context,
	job repository.TranslationJob,
	run repository.JobPhaseRun,
	target personaGenerationTarget,
) bool {
	return service.persistCommonPersonaLink(ctx, job, run, target) == nil
}

func (service *PersonaGenerationPhaseService) buildProviderRequestForExecution(
	run repository.JobPhaseRun,
	target personaGenerationTarget,
) (PersonaGenerationProviderRequest, *personaGenerationExecutionFailure) {
	request, err := service.buildProviderRequest(run, target)
	if err == nil {
		return request, nil
	}
	return PersonaGenerationProviderRequest{}, &personaGenerationExecutionFailure{
		kind:      personaGenerationErrorKindInputMissing,
		reason:    personaGenerationRedactedPublicSummary,
		retryable: true,
	}
}

func (service *PersonaGenerationPhaseService) persistGeneratedPersonaTarget(
	ctx context.Context,
	job repository.TranslationJob,
	run repository.JobPhaseRun,
	target personaGenerationTarget,
	result PersonaGenerationProviderResult,
) bool {
	return service.persistGeneratedPersona(ctx, job, run, target, result) == nil
}

func personaGenerationSaveFailedFailure() *personaGenerationExecutionFailure {
	return &personaGenerationExecutionFailure{
		kind:      personaGenerationErrorKindSaveFailed,
		reason:    personaGenerationRedactedPublicSummary,
		retryable: false,
	}
}

func personaGenerationProviderResultMatchesRequest(
	result PersonaGenerationProviderResult,
	request PersonaGenerationProviderRequest,
) bool {
	return strings.TrimSpace(result.PersonaBody) != "" &&
		strings.TrimSpace(result.NPCCorrelationID) == request.NPCCorrelationID &&
		strings.TrimSpace(result.RequestUnitID) == request.RequestUnitID
}

func (service *PersonaGenerationPhaseService) buildProviderRequest(
	run repository.JobPhaseRun,
	target personaGenerationTarget,
) (PersonaGenerationProviderRequest, error) {
	provider := strings.TrimSpace(run.AIProvider)
	model := strings.TrimSpace(run.ModelName)
	credentialRef := strings.TrimSpace(run.CredentialRef)
	if provider == "" || model == "" || credentialRef == "" {
		return PersonaGenerationProviderRequest{}, fmt.Errorf("persona generation provider input is incomplete")
	}
	utterances := make([]string, 0, len(target.fields))
	for _, field := range target.fields {
		if text := strings.TrimSpace(field.SourceText); text != "" {
			utterances = append(utterances, text)
		}
	}
	if len(utterances) == 0 {
		return PersonaGenerationProviderRequest{}, fmt.Errorf("persona generation request has no utterance")
	}
	attributes := []string{
		fmt.Sprintf("display_name=%s", strings.TrimSpace(target.profile.DisplayName)),
		fmt.Sprintf("editor_id=%s", strings.TrimSpace(target.record.EditorID)),
		fmt.Sprintf("form_id=%s", strings.TrimSpace(target.record.FormID)),
		fmt.Sprintf("voice_type=%s", strings.TrimSpace(target.npc.VoiceType)),
		fmt.Sprintf("race=%s", strings.TrimSpace(personaNullableString(target.npc.Race))),
		fmt.Sprintf("sex=%s", strings.TrimSpace(personaNullableString(target.npc.Sex))),
		fmt.Sprintf("class=%s", strings.TrimSpace(personaNullableString(target.npc.NpcClass))),
	}
	requestUnitID := fmt.Sprintf("phase-run-%d-npc-%d", run.ID, target.profile.ID)
	return PersonaGenerationProviderRequest{
		Provider:                 provider,
		Model:                    model,
		ExecutionMode:            PersonaGenerationExecutionModeSingleRequest,
		CredentialRef:            credentialRef,
		RequestUnitID:            requestUnitID,
		NPCCorrelationID:         fmt.Sprintf("npc-profile-%d", target.profile.ID),
		NPCDisplayName:           strings.TrimSpace(target.profile.DisplayName),
		NPCEditorID:              strings.TrimSpace(target.record.EditorID),
		NPCFormID:                strings.TrimSpace(target.record.FormID),
		NPCAttributes:            attributes,
		ConversationContext:      utterances,
		RecentOriginalUtterances: utterances,
	}, nil
}

func (service *PersonaGenerationPhaseService) persistCommonPersonaLink(
	ctx context.Context,
	job repository.TranslationJob,
	run repository.JobPhaseRun,
	target personaGenerationTarget,
) error {
	if service.transactor == nil {
		return fmt.Errorf("persona generation transactor is not configured")
	}
	err := service.transactor.WithTransaction(ctx, func(txCtx context.Context) error {
		reloadedJob, err := service.jobLifecycleRepository.GetTranslationJobByID(txCtx, job.ID)
		if err != nil {
			return fmt.Errorf("reload translation job for common persona link: %w", err)
		}
		if isPersonaGenerationTerminalJob(reloadedJob.State) {
			return errPersonaGenerationPhaseTerminalJob
		}
		persona, err := service.foundationDataRepository.GetPersonaByNpcProfileID(txCtx, target.profile.ID)
		if err != nil {
			return fmt.Errorf("load common persona for snapshot link: %w", err)
		}
		if persona.TranslationJobID != nil {
			return fmt.Errorf("snapshot persona must be common-scoped: %w", repository.ErrConflict)
		}
		return service.ensurePhaseRunPersonaLink(txCtx, run.ID, persona.ID, personaGenerationPhaseLinkRoleSnapshot)
	})
	if err != nil {
		return fmt.Errorf("persist common persona link: %w", err)
	}
	return nil
}

func (service *PersonaGenerationPhaseService) persistGeneratedPersona(
	ctx context.Context,
	job repository.TranslationJob,
	run repository.JobPhaseRun,
	target personaGenerationTarget,
	result PersonaGenerationProviderResult,
) error {
	if service.transactor == nil {
		return fmt.Errorf("persona generation transactor is not configured")
	}
	err := service.transactor.WithTransaction(ctx, func(txCtx context.Context) error {
		reloadedJob, err := service.jobLifecycleRepository.GetTranslationJobByID(txCtx, job.ID)
		if err != nil {
			return fmt.Errorf("reload translation job for persona save: %w", err)
		}
		if isPersonaGenerationTerminalJob(reloadedJob.State) {
			return errPersonaGenerationPhaseTerminalJob
		}
		personaDraft := repository.PersonaDraft{
			NpcProfileID:           target.profile.ID,
			TranslationJobID:       &job.ID,
			PersonaLifecycle:       personaGenerationPersonaLifecycleJob,
			PersonaScope:           personaGenerationPersonaScopeJob,
			PersonaSource:          personaGenerationPersonaSourceProvider,
			PersonaDescription:     strings.TrimSpace(result.PersonaBody),
			SpeechStyle:            "",
			PersonalitySummary:     "",
			EvidenceUtteranceCount: len(target.fields),
		}
		persona, err := service.upsertJobScopedPersona(txCtx, job.ID, target.profile.ID, personaDraft)
		if err != nil {
			return err
		}
		if err := service.ensurePersonaFieldEvidence(txCtx, persona.ID, target.fields); err != nil {
			return err
		}
		if err := service.ensurePhaseRunPersonaLink(txCtx, run.ID, persona.ID, personaGenerationPhaseLinkRoleApplied); err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return fmt.Errorf("persist generated persona: %w", err)
	}
	return nil
}

func (service *PersonaGenerationPhaseService) upsertJobScopedPersona(
	ctx context.Context,
	jobID int64,
	npcProfileID int64,
	draft repository.PersonaDraft,
) (repository.Persona, error) {
	persona, err := service.foundationDataRepository.GetPersonaByNpcProfileID(ctx, npcProfileID)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			created, createErr := service.foundationDataRepository.CreatePersona(ctx, draft)
			if createErr != nil {
				return repository.Persona{}, fmt.Errorf("create job-scoped persona: %w", createErr)
			}
			return created, nil
		}
		return repository.Persona{}, fmt.Errorf("load job-scoped persona candidate: %w", err)
	}
	if persona.TranslationJobID == nil || *persona.TranslationJobID != jobID {
		return repository.Persona{}, fmt.Errorf("persona belongs to another scope: %w", repository.ErrConflict)
	}
	updated, updateErr := service.foundationDataRepository.UpdatePersona(ctx, persona.ID, repository.PersonaUpdateDraft{
		PersonaLifecycle:       draft.PersonaLifecycle,
		PersonaScope:           draft.PersonaScope,
		PersonaSource:          draft.PersonaSource,
		PersonaDescription:     draft.PersonaDescription,
		SpeechStyle:            draft.SpeechStyle,
		PersonalitySummary:     draft.PersonalitySummary,
		EvidenceUtteranceCount: draft.EvidenceUtteranceCount,
	})
	if updateErr != nil {
		return repository.Persona{}, fmt.Errorf("update job-scoped persona: %w", updateErr)
	}
	return updated, nil
}

func (service *PersonaGenerationPhaseService) ensurePersonaFieldEvidence(
	ctx context.Context,
	personaID int64,
	fields []repository.TranslationField,
) error {
	existing, err := service.foundationDataRepository.ListPersonaFieldEvidenceByPersonaID(ctx, personaID)
	if err != nil {
		return fmt.Errorf("list persona field evidence: %w", err)
	}
	existingFieldIDs := make(map[int64]struct{}, len(existing))
	for _, item := range existing {
		existingFieldIDs[item.TranslationFieldID] = struct{}{}
	}
	for _, field := range fields {
		if _, ok := existingFieldIDs[field.ID]; ok {
			continue
		}
		if _, createErr := service.foundationDataRepository.CreatePersonaFieldEvidence(ctx, repository.PersonaFieldEvidenceDraft{
			PersonaID:          personaID,
			TranslationFieldID: field.ID,
			EvidenceRole:       personaGenerationEvidenceRoleSource,
		}); createErr != nil {
			return fmt.Errorf("create persona field evidence: %w", createErr)
		}
	}
	return nil
}

func (service *PersonaGenerationPhaseService) ensurePhaseRunPersonaLink(
	ctx context.Context,
	phaseRunID int64,
	personaID int64,
	role string,
) error {
	existing, err := service.jobLifecycleRepository.ListPhaseRunPersonasByPhaseRunID(ctx, phaseRunID)
	if err != nil {
		return fmt.Errorf("list phase run personas before create: %w", err)
	}
	for _, link := range existing {
		if link.PersonaID == personaID {
			return nil
		}
	}
	if _, err := service.jobLifecycleRepository.CreatePhaseRunPersona(ctx, repository.PhaseRunPersonaDraft{
		PhaseRunID: phaseRunID,
		PersonaID:  personaID,
		Role:       role,
	}); err != nil {
		return fmt.Errorf("create phase run persona link: %w", err)
	}
	return nil
}

func (service *PersonaGenerationPhaseService) finishPhaseRun(
	ctx context.Context,
	run repository.JobPhaseRun,
	snapshot personaGenerationTargetSnapshot,
	runSnapshot personaGenerationRunSnapshot,
	failures []personaGenerationExecutionFailure,
) (repository.JobPhaseRun, error) {
	latestError := ""
	retryable := false
	if len(failures) > 0 {
		lastFailure := failures[len(failures)-1]
		latestError = normalizePersonaField(lastFailure.kind)
		retryable = lastFailure.retryable
	}
	missingCount := personaMax(0, snapshot.targetCount-runSnapshot.totalLinkedCount)
	nextState := personaGenerationPhaseStateCompleted
	progressPercent := 100
	var executionErr error
	if missingCount > 0 {
		progressPercent = runSnapshot.totalLinkedCount * 100 / personaMax(snapshot.targetCount, 1)
		if runSnapshot.totalLinkedCount > 0 || retryable {
			nextState = personaGenerationPhaseStateRecoverableFail
		} else {
			nextState = personaGenerationPhaseStateFailed
		}
		switch latestError {
		case personaGenerationErrorKindProviderFailure:
			executionErr = errPersonaGenerationPhaseProviderFailure
		case personaGenerationErrorKindInvalidProvider:
			executionErr = errPersonaGenerationPhaseInvalidProviderResponse
		case personaGenerationErrorKindInputMissing:
			executionErr = errPersonaGenerationPhaseInputMissing
		case personaGenerationErrorKindSaveFailed:
			executionErr = errPersonaGenerationPhaseSaveFailed
		}
	}
	now := service.now()
	updatedRun, err := service.jobLifecycleRepository.UpdateJobPhaseRun(ctx, run.ID, repository.JobPhaseRunUpdateDraft{
		State:               nextState,
		ProgressPercent:     progressPercent,
		LatestExternalRunID: run.LatestExternalRunID,
		LatestError:         latestError,
		StartedAt:           termPersonaStartedAt(run.StartedAt, now),
		FinishedAt:          &now,
	})
	if err != nil {
		return run, fmt.Errorf("finalize persona generation phase run: %w", err)
	}
	if executionErr != nil {
		return updatedRun, executionErr
	}
	return updatedRun, nil
}

func (service *PersonaGenerationPhaseService) loadRunSnapshot(
	ctx context.Context,
	run *repository.JobPhaseRun,
) (personaGenerationRunSnapshot, error) {
	if run == nil {
		return personaGenerationRunSnapshot{
			personasByID:    map[int64]repository.Persona{},
			linkedByProfile: map[int64]repository.Persona{},
		}, nil
	}
	links, err := service.jobLifecycleRepository.ListPhaseRunPersonasByPhaseRunID(ctx, run.ID)
	if err != nil {
		return personaGenerationRunSnapshot{}, fmt.Errorf("list phase run personas: %w", err)
	}
	snapshot := personaGenerationRunSnapshot{
		links:           append([]repository.PhaseRunPersona(nil), links...),
		personasByID:    make(map[int64]repository.Persona, len(links)),
		linkedByProfile: make(map[int64]repository.Persona, len(links)),
	}
	for _, link := range links {
		persona, loadErr := service.foundationDataRepository.GetPersonaByID(ctx, link.PersonaID)
		if loadErr != nil {
			return personaGenerationRunSnapshot{}, fmt.Errorf("load persona %d for phase snapshot: %w", link.PersonaID, loadErr)
		}
		if _, exists := snapshot.personasByID[persona.ID]; exists {
			continue
		}
		snapshot.personasByID[persona.ID] = persona
		snapshot.linkedByProfile[persona.NpcProfileID] = persona
		snapshot.totalLinkedCount++
		if persona.TranslationJobID != nil && run.TranslationJobID == *persona.TranslationJobID {
			snapshot.generatedCount++
		}
	}
	return snapshot, nil
}

func (service *PersonaGenerationPhaseService) commandFromRun(
	ctx context.Context,
	job repository.TranslationJob,
	run repository.JobPhaseRun,
	snapshot personaGenerationTargetSnapshot,
	termRun *repository.JobPhaseRun,
) PersonaGenerationPhaseCommandReadModel {
	state := normalizePersonaGenerationPhaseState(run.State, snapshot.targetCount)
	progress := service.buildProgress(state, &run, snapshot)
	resultSummary, errorSummary, canStartBodyPhase := service.buildResultState(ctx, &run, snapshot)
	command := PersonaGenerationPhaseCommandReadModel{
		JobID:        job.ID,
		CurrentPhase: personaGenerationCurrentPhase,
		PhaseState:   state,
		PhaseRunID:   clonePersonaInt64Pointer(&run.ID),
		StartedAt:    clonePersonaTimePointer(run.StartedAt),
		FinishedAt:   clonePersonaTimePointer(run.FinishedAt),
		Progress:     progress,
		TargetSummary: PersonaGenerationTargetSummaryReadModel{
			TargetCount:            snapshot.targetCount,
			CommonPersonaHitCount:  snapshot.commonHits,
			CommonPersonaMissCount: snapshot.commonMisses,
			SkippedCount:           len(snapshot.skippedReasons),
			SkippedReasons:         append([]string(nil), snapshot.skippedReasons...),
			TargetSnapshotDigest:   snapshot.digest,
		},
		ResultSummary:     resultSummary,
		Retryable:         personaGenerationRetryAllowed(run.State, run.LatestError),
		CanStartBodyPhase: canStartBodyPhase,
		ErrorSummary:      errorSummary,
	}
	if termRun != nil {
		command.Execution = service.buildExecutionSummary(&run, termRun, snapshot, resultSummary)
	}
	return command
}

func termPersonaStartedAt(startedAt *time.Time, now time.Time) *time.Time {
	if startedAt != nil {
		return startedAt
	}
	return &now
}

func personaGenerationExecutionOutputCount(resultSummary *PersonaGenerationPhaseResultSummaryReadModel) int {
	if resultSummary == nil {
		return 0
	}
	return resultSummary.PersonaCount
}

func personaGenerationRetryAllowed(state string, latestError string) bool {
	if normalizePersonaField(state) != personaGenerationPhaseStateRecoverableFail {
		return false
	}
	switch normalizePersonaField(latestError) {
	case personaGenerationErrorKindProviderFailure,
		personaGenerationErrorKindInvalidProvider,
		personaGenerationErrorKindInputMissing:
		return true
	default:
		return false
	}
}

func personaGenerationCancelAllowed(jobState string, runState string) bool {
	if isPersonaGenerationTerminalJob(jobState) {
		return false
	}
	switch normalizePersonaField(runState) {
	case personaGenerationPhaseStateRunning,
		personaGenerationPhaseStatePaused,
		personaGenerationPhaseStateRecoverableFail:
		return true
	default:
		return false
	}
}

func personaGenerationSnapshotDriftReason(
	run *repository.JobPhaseRun,
	snapshot personaGenerationTargetSnapshot,
) *string {
	if run == nil {
		return nil
	}
	savedDigest := strings.TrimSpace(run.LatestExternalRunID)
	if savedDigest == "" || savedDigest == snapshot.digest {
		return nil
	}
	reason := "persona target snapshot changed after phase run start"
	return &reason
}

func personaGenerationPhaseRunEvidenceRefs(phaseRunPersonas []repository.PhaseRunPersona) []string {
	evidenceRefs := make([]string, 0, len(phaseRunPersonas))
	for _, link := range phaseRunPersonas {
		evidenceRefs = append(evidenceRefs, fmt.Sprintf("phase-run-persona:%d", link.ID))
	}
	return evidenceRefs
}

func personaGenerationExecutionEvidenceRefs(run *repository.JobPhaseRun) []string {
	if run == nil {
		return nil
	}
	return []string{fmt.Sprintf("phase-run:%d", run.ID)}
}

func (service *PersonaGenerationPhaseService) buildExecutionSummary(
	executionRun *repository.JobPhaseRun,
	metadataRun *repository.JobPhaseRun,
	snapshot personaGenerationTargetSnapshot,
	resultSummary *PersonaGenerationPhaseResultSummaryReadModel,
) PersonaGenerationExecutionSummaryReadModel {
	if metadataRun == nil {
		metadataRun = executionRun
	}
	return PersonaGenerationExecutionSummaryReadModel{
		CredentialRef: strings.TrimSpace(personaGenerationExecutionField(metadataRun, func(run *repository.JobPhaseRun) string { return run.CredentialRef })),
		Provider:      strings.TrimSpace(personaGenerationExecutionField(metadataRun, func(run *repository.JobPhaseRun) string { return run.AIProvider })),
		Model:         strings.TrimSpace(personaGenerationExecutionField(metadataRun, func(run *repository.JobPhaseRun) string { return run.ModelName })),
		ExecutionMode: strings.TrimSpace(personaGenerationExecutionField(metadataRun, func(run *repository.JobPhaseRun) string { return run.ExecutionMode })),
		PromptDigest:  service.aggregatePromptDigest(executionRun, snapshot),
		InputCount:    snapshot.targetCount,
		OutputCount:   personaGenerationExecutionOutputCount(resultSummary),
		EvidenceRefs:  personaGenerationExecutionEvidenceRefs(executionRun),
	}
}

func personaGenerationExecutionField(run *repository.JobPhaseRun, pick func(*repository.JobPhaseRun) string) string {
	if run == nil {
		return ""
	}
	return pick(run)
}

func (service *PersonaGenerationPhaseService) aggregatePromptDigest(
	run *repository.JobPhaseRun,
	snapshot personaGenerationTargetSnapshot,
) string {
	if run == nil {
		return ""
	}
	digests := make([]string, 0, len(snapshot.targets))
	for _, target := range snapshot.targets {
		if target.hit {
			continue
		}
		request, err := service.buildProviderRequest(*run, target)
		if err != nil {
			return ""
		}
		prompt, err := BuildPersonaGenerationPrompt(request)
		if err != nil {
			return ""
		}
		digests = append(digests, personaGenerationPromptDigest(prompt))
	}
	sort.Strings(digests)
	builder := strings.Builder{}
	builder.WriteString("persona-generation-prompts|")
	for _, digest := range digests {
		builder.WriteString(digest)
		builder.WriteByte('|')
	}
	sum := sha256.Sum256([]byte(builder.String()))
	return "sha256:" + hex.EncodeToString(sum[:])
}

func (service *PersonaGenerationPhaseService) distinctPhaseRunPersonaCount(
	ctx context.Context,
	phaseRunPersonas []repository.PhaseRunPersona,
) int {
	seenPersonaIDs := make(map[int64]struct{}, len(phaseRunPersonas))
	seenProfileIDs := make(map[int64]struct{}, len(phaseRunPersonas))
	for _, link := range phaseRunPersonas {
		if _, exists := seenPersonaIDs[link.PersonaID]; exists {
			continue
		}
		seenPersonaIDs[link.PersonaID] = struct{}{}
		persona, err := service.foundationDataRepository.GetPersonaByID(ctx, link.PersonaID)
		if err != nil {
			continue
		}
		seenProfileIDs[persona.NpcProfileID] = struct{}{}
	}
	return len(seenProfileIDs)
}

func (service *PersonaGenerationPhaseService) tryLoadRunSnapshot(
	ctx context.Context,
	run *repository.JobPhaseRun,
) (personaGenerationRunSnapshot, bool) {
	snapshot, err := service.loadRunSnapshot(ctx, run)
	if err != nil {
		return personaGenerationRunSnapshot{}, false
	}
	return snapshot, true
}

func personaGenerationTermPhase(phases []repository.JobPhaseRun) (repository.JobPhaseRun, error) {
	for _, phase := range phases {
		if normalizePersonaField(phase.PhaseType) == personaGenerationTermPhaseType {
			return phase, nil
		}
	}
	for _, phase := range phases {
		if normalizePersonaField(phase.PhaseType) == personaGenerationInitialPhaseType {
			return phase, nil
		}
	}
	return repository.JobPhaseRun{}, fmt.Errorf("load persona generation term phase: %w", repository.ErrNotFound)
}

func buildPersonaGenerationTargetSnapshotDigest(snapshot personaGenerationTargetSnapshot) string {
	builder := strings.Builder{}
	builder.WriteString(fmt.Sprintf("targets:%d|hits:%d|misses:%d|", snapshot.targetCount, snapshot.commonHits, snapshot.commonMisses))
	for _, reason := range snapshot.skippedReasons {
		builder.WriteString(reason)
		builder.WriteByte('|')
	}
	targets := append([]personaGenerationTarget(nil), snapshot.targets...)
	sort.Slice(targets, func(i, j int) bool {
		return targets[i].profile.ID < targets[j].profile.ID
	})
	for _, target := range targets {
		builder.WriteString(fmt.Sprintf("%d:%s:%t|", target.profile.ID, strings.TrimSpace(target.record.EditorID), target.hit))
		fieldIDs := make([]int64, 0, len(target.fields))
		for _, field := range target.fields {
			fieldIDs = append(fieldIDs, field.ID)
		}
		sort.Slice(fieldIDs, func(i, j int) bool { return fieldIDs[i] < fieldIDs[j] })
		for _, fieldID := range fieldIDs {
			builder.WriteString(fmt.Sprintf("%d,", fieldID))
		}
		builder.WriteByte('|')
	}
	sum := sha256.Sum256([]byte(builder.String()))
	return "sha256:" + hex.EncodeToString(sum[:])
}

func normalizePersonaGenerationPhaseState(state string, targetCount int) string {
	normalized := normalizePersonaField(state)
	if normalized == personaGenerationPhaseStateCompleted && targetCount == 0 {
		return personaGenerationPhaseStateEmptyCompleted
	}
	if normalized == "" {
		return personaGenerationPhaseStateNotStarted
	}
	return normalized
}

func isPersonaGenerationActiveState(state string) bool {
	switch normalizePersonaField(state) {
	case personaGenerationPhaseStateRunning, personaGenerationPhaseStatePaused, personaGenerationPhaseStateRecoverableFail:
		return true
	default:
		return false
	}
}

func isPersonaGenerationTerminalJob(state string) bool {
	switch normalizePersonaField(state) {
	case personaGenerationJobStateCompleted, personaGenerationJobStateFailed, personaGenerationJobStateCanceled:
		return true
	default:
		return false
	}
}

func percentageCount(percent int, total int) int {
	if total <= 0 || percent <= 0 {
		return 0
	}
	return personaMin(total, int(float64(total)*float64(percent)/100.0))
}

func personaBlockedReason(enabled bool, reason string) *string {
	if enabled {
		return nil
	}
	return stringPersonaPointer(reason)
}

func normalizePersonaField(value string) string {
	return strings.ToLower(strings.TrimSpace(value))
}

func clonePersonaTimePointer(value *time.Time) *time.Time {
	if value == nil {
		return nil
	}
	cloned := *value
	return &cloned
}

func clonePersonaInt64Pointer(value *int64) *int64 {
	if value == nil {
		return nil
	}
	cloned := *value
	return &cloned
}

func clonePersonaStringPointer(value *string) *string {
	if value == nil {
		return nil
	}
	cloned := *value
	return &cloned
}

func personaNullableString(value *string) string {
	if value == nil {
		return ""
	}
	return *value
}

func stringPersonaPointer(value string) *string {
	if strings.TrimSpace(value) == "" {
		return nil
	}
	cloned := value
	return &cloned
}

func personaMin(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func personaMax(a, b int) int {
	if a > b {
		return a
	}
	return b
}

var errPersonaGenerationPhaseExecutionRejected = errors.New("persona generation phase execution rejected")
var (
	errPersonaGenerationPhaseTerminalJob             = errors.New("persona generation phase terminal job")
	errPersonaGenerationPhaseProviderFailure         = errors.New("persona generation provider failure")
	errPersonaGenerationPhaseInvalidProviderResponse = errors.New("persona generation invalid provider response")
	errPersonaGenerationPhaseInputMissing            = errors.New("persona generation input missing")
	errPersonaGenerationPhaseSaveFailed              = errors.New("persona generation save failed")
)
