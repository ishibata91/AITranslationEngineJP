package service

import (
	"context"
	"errors"
	"fmt"
	"sort"
	"strings"
	"time"

	"aitranslationenginejp/internal/repository"
)

const (
	bodyTranslationCurrentPhase              = "body_translation"
	bodyTranslationPhaseType                 = "body_translation"
	bodyTranslationPersonaPhaseType          = "persona_generation"
	bodyTranslationPhaseStateIdleReady       = "idle_ready"
	bodyTranslationPhaseStateRunning         = "running"
	bodyTranslationPhaseStateCompleted       = "completed"
	bodyTranslationPhaseStatePaused          = "paused"
	bodyTranslationPhaseStateRecoverableFail = "recoverable_failed"
	bodyTranslationPhaseStateCanceled        = "canceled"
	bodyTranslationJobStateRunning           = "running"
	bodyTranslationJobStateCompleted         = "completed"
	bodyTranslationJobStateFailed            = "failed"
	bodyTranslationJobStateCanceled          = "canceled"
	bodyTranslationStartRejectMessage        = "body translation phase execution rejected"
	bodyTranslationInputSnapshotFailedKind   = "input_snapshot_failed"
	bodyTranslationInputSnapshotDriftReason  = "body translation input snapshot changed after phase start"
)

var errBodyTranslationPhaseExecutionRejected = errors.New(bodyTranslationStartRejectMessage)

type bodyTranslationPhaseJobLifecycleRepository interface {
	GetTranslationJobByID(ctx context.Context, id int64) (repository.TranslationJob, error)
	UpdateTranslationJob(ctx context.Context, id int64, draft repository.TranslationJobUpdateDraft) (repository.TranslationJob, error)
	CreateJobPhaseRun(ctx context.Context, draft repository.JobPhaseRunDraft) (repository.JobPhaseRun, error)
	FindJobPhaseRun(ctx context.Context, translationJobID int64, phaseType string) (repository.JobPhaseRun, error)
	UpdateJobPhaseRun(ctx context.Context, id int64, draft repository.JobPhaseRunUpdateDraft) (repository.JobPhaseRun, error)
	ListJobPhaseRunsByJobID(ctx context.Context, jobID int64) ([]repository.JobPhaseRun, error)
	CreatePhaseRunTranslationField(ctx context.Context, draft repository.PhaseRunTranslationFieldDraft) (repository.PhaseRunTranslationField, error)
	ListPhaseRunTranslationFieldsByPhaseRunID(ctx context.Context, phaseRunID int64) ([]repository.PhaseRunTranslationField, error)
}

type bodyTranslationPhaseFoundationDataRepository interface {
	ListDictionaryEntries(ctx context.Context, translationJobID *int64, lifecycle string, scope string, sourceTerm string) ([]repository.DictionaryEntry, error)
	ListPersonasByTranslationJobID(ctx context.Context, translationJobID int64) ([]repository.Persona, error)
}

type bodyTranslationPhaseTranslationSourceRepository interface {
	GetXEditExtractedDataByID(ctx context.Context, id int64) (repository.XEditExtractedData, error)
	ListTranslationRecordsByXEditID(ctx context.Context, xEditID int64) ([]repository.TranslationRecord, error)
	ListTranslationFieldsByTranslationRecordID(ctx context.Context, translationRecordID int64) ([]repository.TranslationField, error)
}

type bodyTranslationPhaseJobOutputRepository interface {
	CreateJobTranslationField(ctx context.Context, draft repository.JobTranslationFieldDraft) (repository.JobTranslationField, error)
	UpdateJobTranslationField(ctx context.Context, id int64, draft repository.JobTranslationFieldUpdateDraft) (repository.JobTranslationField, error)
	ListJobTranslationFieldsByJobID(ctx context.Context, jobID int64) ([]repository.JobTranslationField, error)
}

type bodyTranslationPhaseRuntimeSnapshotReader interface {
	GetTranslationJobPhaseRuntimeSnapshot(
		ctx context.Context,
		translationJobID int64,
		phaseID string,
	) (repository.TranslationJobPhaseRuntimeSnapshot, error)
}

// BodyTranslationPhaseProgressReadModel stores body phase progress for the usecase boundary.
type BodyTranslationPhaseProgressReadModel struct {
	Percent         int
	ProcessedCount  int
	TotalCount      int
	TargetCount     int
	TranslatedCount int
	SkippedCount    int
	CurrentStep     string
}

// BodyTranslationPhaseInputSummaryReadModel stores the frozen input snapshot summary.
type BodyTranslationPhaseInputSummaryReadModel struct {
	TargetCount      int
	SkippedReasons   []string
	InputSnapshotRef *string
	DictionaryDigest string
	PersonaDigest    string
	MetadataDigest   string
	PromptDigest     string
}

// BodyTranslationPhaseRequestSummaryReadModel stores provider request counts after filtering.
type BodyTranslationPhaseRequestSummaryReadModel struct {
	ProviderTargetCount              int
	ExactDictionaryExclusionCount    int
	PartialDictionaryConstraintCount int
}

// BodyTranslationPhaseExecutionSummaryReadModel stores redacted execution settings.
type BodyTranslationPhaseExecutionSummaryReadModel struct {
	CredentialRef    string
	CredentialState  string
	EndpointSummary  *string
	Provider         string
	Model            string
	ExecutionMode    string
	RequestUnitCount int
	OutputCount      int
}

// BodyTranslationPhaseFieldResultSummaryReadModel stores persisted field result counts.
type BodyTranslationPhaseFieldResultSummaryReadModel struct {
	TranslatedCount       int
	FailedCount           int
	SkippedCount          int
	ProtectionFailedCount int
	OutputReadyCount      int
	OutputCount           int
	FieldResults          []BodyTranslationPhaseFieldResultItemReadModel
}

// BodyTranslationPhaseFieldIdentityReadModel identifies one translated field.
type BodyTranslationPhaseFieldIdentityReadModel struct {
	TranslationFieldID      int64
	PhaseTranslationFieldID int64
	RecordType              string
	FieldType               string
	FormID                  string
	EditorID                string
	FieldLabel              string
}

// BodyTranslationPhaseFieldResultItemReadModel stores one public field result row.
type BodyTranslationPhaseFieldResultItemReadModel struct {
	Identity                    BodyTranslationPhaseFieldIdentityReadModel
	FieldID                     int64
	FieldLabel                  string
	SourceExcerpt               string
	TranslatedText              string
	OutputStatus                string
	ProtectionValidationResult  string
	ProtectionValidationSummary string
	RetryCount                  int
}

// BodyTranslationPhaseErrorSummaryReadModel stores redacted error summary fields.
type BodyTranslationPhaseErrorSummaryReadModel struct {
	ErrorKind  string
	Reason     string
	Retryable  bool
	IsRedacted bool
}

// BodyTranslationPhaseActionEnablementReadModel stores action enablement for Job Run.
type BodyTranslationPhaseActionEnablementReadModel struct {
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

// BodyTranslationOutputReadinessReadModel stores downstream output readiness fields.
type BodyTranslationOutputReadinessReadModel struct {
	JobID               int64
	CurrentPhase        string
	PhaseState          string
	Ready               bool
	BlockedReason       string
	ErrorKind           string
	CompletedFieldCount int
	StatusConsistent    bool
	OutputCount         int
}

// BodyTranslationPhaseSummaryReadModel stores body phase summary payload data.
type BodyTranslationPhaseSummaryReadModel struct {
	JobID              int64
	CurrentPhase       string
	PhaseState         string
	PhaseRunID         *int64
	StartedAt          *time.Time
	FinishedAt         *time.Time
	Progress           BodyTranslationPhaseProgressReadModel
	InputSummary       BodyTranslationPhaseInputSummaryReadModel
	RequestSummary     BodyTranslationPhaseRequestSummaryReadModel
	Execution          BodyTranslationPhaseExecutionSummaryReadModel
	FieldResultSummary *BodyTranslationPhaseFieldResultSummaryReadModel
	ResultSummary      *BodyTranslationPhaseFieldResultSummaryReadModel
	FieldResults       []BodyTranslationPhaseFieldResultItemReadModel
	ErrorSummary       *BodyTranslationPhaseErrorSummaryReadModel
	ActionEnablement   BodyTranslationPhaseActionEnablementReadModel
	OutputReadiness    BodyTranslationOutputReadinessReadModel
}

// BodyTranslationPhaseCommandReadModel stores body phase command payload data.
type BodyTranslationPhaseCommandReadModel struct {
	JobID               int64
	CurrentPhase        string
	PhaseState          string
	PhaseRunID          *int64
	StartedAt           *time.Time
	FinishedAt          *time.Time
	Progress            BodyTranslationPhaseProgressReadModel
	InputSnapshotDigest string
	InputSummary        BodyTranslationPhaseInputSummaryReadModel
	RequestSummary      BodyTranslationPhaseRequestSummaryReadModel
	Execution           BodyTranslationPhaseExecutionSummaryReadModel
	FieldResultSummary  *BodyTranslationPhaseFieldResultSummaryReadModel
	ResultSummary       *BodyTranslationPhaseFieldResultSummaryReadModel
	FieldResults        []BodyTranslationPhaseFieldResultItemReadModel
	Retryable           bool
	OutputReadiness     BodyTranslationOutputReadinessReadModel
	ErrorSummary        *BodyTranslationPhaseErrorSummaryReadModel
}

type bodyTranslationStartRejection struct {
	errorKind string
	reason    string
}

type bodyTranslationLoadedContext struct {
	job                  repository.TranslationJob
	bodyRun              *repository.JobPhaseRun
	personaRun           *repository.JobPhaseRun
	records              []repository.TranslationRecord
	fieldsByRecord       map[int64][]repository.TranslationField
	dictionary           []repository.DictionaryEntry
	persona              repository.Persona
	execution            BodyTranslationPhaseExecutionSummaryReadModel
	snapshot             bodyTranslationInputSnapshot
	inputSnapshotDrifted bool
	outputFields         []repository.JobTranslationField
}

// BodyTranslationPhaseService orchestrates body phase startup and input snapshot creation.
type BodyTranslationPhaseService struct {
	now                      func() time.Time
	jobLifecycleRepository   bodyTranslationPhaseJobLifecycleRepository
	foundationDataRepository bodyTranslationPhaseFoundationDataRepository
	translationSourceReader  bodyTranslationPhaseTranslationSourceRepository
	jobOutputRepository      bodyTranslationPhaseJobOutputRepository
	bodyTranslationProvider  BodyTranslationProvider
	providerSettings         ProviderSettingsConsumer
	executionSnapshots       map[int64]providerExecutionSnapshot
	transactor               repository.Transactor
}

// NewBodyTranslationPhaseService creates the backend body phase service.
func NewBodyTranslationPhaseService(
	jobLifecycleRepository bodyTranslationPhaseJobLifecycleRepository,
	foundationDataRepository bodyTranslationPhaseFoundationDataRepository,
	translationSourceReader bodyTranslationPhaseTranslationSourceRepository,
	jobOutputRepository bodyTranslationPhaseJobOutputRepository,
	transactor repository.Transactor,
) *BodyTranslationPhaseService {
	return &BodyTranslationPhaseService{
		now:                      func() time.Time { return time.Now().UTC() },
		jobLifecycleRepository:   jobLifecycleRepository,
		foundationDataRepository: foundationDataRepository,
		translationSourceReader:  translationSourceReader,
		jobOutputRepository:      jobOutputRepository,
		executionSnapshots:       map[int64]providerExecutionSnapshot{},
		transactor:               transactor,
	}
}

// WithBodyTranslationProvider wires one provider adapter into the service.
func (service *BodyTranslationPhaseService) WithBodyTranslationProvider(provider BodyTranslationProvider) *BodyTranslationPhaseService {
	if service == nil {
		return nil
	}
	service.bodyTranslationProvider = provider
	return service
}

// WithBodyTranslationProviderSettings injects the provider settings resolver used at phase start.
func (service *BodyTranslationPhaseService) WithBodyTranslationProviderSettings(
	consumer ProviderSettingsConsumer,
) *BodyTranslationPhaseService {
	if service == nil {
		return nil
	}
	service.providerSettings = consumer
	return service
}

// ReadSummary returns the current body phase summary.
func (service *BodyTranslationPhaseService) ReadSummary(
	ctx context.Context,
	jobID int64,
) (BodyTranslationPhaseSummaryReadModel, error) {
	loaded, err := service.loadContext(ctx, jobID)
	if err != nil {
		return BodyTranslationPhaseSummaryReadModel{}, err
	}
	phaseState := bodyTranslationPhaseStateIdleReady
	var phaseRunID *int64
	var startedAt *time.Time
	var finishedAt *time.Time
	if loaded.bodyRun != nil {
		phaseState = strings.TrimSpace(loaded.bodyRun.State)
		phaseRunID = &loaded.bodyRun.ID
		startedAt = cloneBodyTranslationTimePointer(loaded.bodyRun.StartedAt)
		finishedAt = cloneBodyTranslationTimePointer(loaded.bodyRun.FinishedAt)
	}
	outputReadiness := service.buildOutputReadiness(loaded)
	resultSummary := service.buildLoadedFieldResultSummary(loaded)
	return BodyTranslationPhaseSummaryReadModel{
		JobID:              loaded.job.ID,
		CurrentPhase:       bodyTranslationCurrentPhase,
		PhaseState:         phaseState,
		PhaseRunID:         cloneBodyTranslationInt64Pointer(phaseRunID),
		StartedAt:          startedAt,
		FinishedAt:         finishedAt,
		Progress:           service.buildProgress(phaseState, loaded.snapshot, loaded.outputFields),
		InputSummary:       toBodyTranslationInputSummaryReadModel(loaded.snapshot),
		RequestSummary:     toBodyTranslationRequestSummaryReadModel(loaded.snapshot),
		Execution:          loaded.execution,
		FieldResultSummary: resultSummary,
		ResultSummary:      resultSummary,
		FieldResults:       service.buildFieldResultItems(loaded),
		ErrorSummary:       service.buildPhaseErrorSummary(loaded, outputReadiness),
		ActionEnablement:   service.buildActionEnablement(loaded, nil),
		OutputReadiness:    outputReadiness,
	}, nil
}

// StartPhase creates one body phase run when start conditions pass.
func (service *BodyTranslationPhaseService) StartPhase(
	ctx context.Context,
	jobID int64,
) (BodyTranslationPhaseCommandReadModel, error) {
	if service.transactor == nil {
		return BodyTranslationPhaseCommandReadModel{}, fmt.Errorf("start body translation phase: transactor is not configured")
	}
	loaded, err := service.loadContext(ctx, jobID)
	if err != nil {
		return BodyTranslationPhaseCommandReadModel{}, err
	}
	if loaded.bodyRun == nil {
		resolvedExecution, rejection, resolveErr := service.resolveExecutionSnapshotForStart(ctx, loaded.execution)
		if resolveErr != nil {
			return BodyTranslationPhaseCommandReadModel{}, resolveErr
		}
		if rejection != nil {
			result := service.rejectedCommand(loaded, *rejection)
			return result, errBodyTranslationPhaseExecutionRejected
		}
		loaded.execution = resolvedExecution
	}
	if rejection := service.startRejection(loaded); rejection != nil {
		result := service.rejectedCommand(loaded, *rejection)
		return result, errBodyTranslationPhaseExecutionRejected
	}

	var createdRun repository.JobPhaseRun
	var updatedJob repository.TranslationJob
	now := service.now()
	err = service.transactor.WithTransaction(ctx, func(txCtx context.Context) error {
		runState := bodyTranslationPhaseStateRunning
		progressPercent := 0
		finishedAt := (*time.Time)(nil)
		jobState := loaded.job.State
		jobProgress := loaded.job.ProgressPercent
		if loaded.snapshot.ProviderTargetCount == 0 {
			runState = bodyTranslationPhaseStateCompleted
			progressPercent = 100
			finishedAt = &now
			jobState = bodyTranslationJobStateCompleted
			jobProgress = 100
		}

		run, createErr := service.jobLifecycleRepository.CreateJobPhaseRun(txCtx, repository.JobPhaseRunDraft{
			TranslationJobID:       loaded.job.ID,
			PhaseType:              bodyTranslationPhaseType,
			State:                  runState,
			ExecutionOrder:         service.nextExecutionOrder(loaded),
			SnapshotFieldCount:     loaded.snapshot.TargetCount,
			ProviderTargetCount:    loaded.snapshot.ProviderTargetCount,
			ExactExclusionCount:    loaded.snapshot.ExactExclusionCount,
			PartialConstraintCount: loaded.snapshot.PartialConstraintCount,
			AIProvider:             loaded.execution.Provider,
			ModelName:              loaded.execution.Model,
			ExecutionMode:          loaded.execution.ExecutionMode,
			CredentialRef:          loaded.execution.CredentialRef,
			InstructionKind:        bodyTranslationPhaseType,
			InputSnapshotDigest:    loaded.snapshot.InputSnapshotDigest,
			DictionaryDigest:       loaded.snapshot.DictionaryDigest,
			PersonaDigest:          loaded.snapshot.PersonaDigest,
			MetadataDigest:         loaded.snapshot.MetadataDigest,
			PromptDigest:           loaded.snapshot.PromptDigest,
		})
		if createErr != nil {
			return fmt.Errorf("create body translation phase run: %w", createErr)
		}
		run, createErr = service.jobLifecycleRepository.UpdateJobPhaseRun(txCtx, run.ID, repository.JobPhaseRunUpdateDraft{
			State:           runState,
			ProgressPercent: progressPercent,
			StartedAt:       &now,
			FinishedAt:      finishedAt,
		})
		if createErr != nil {
			return fmt.Errorf("update body translation phase run: %w", createErr)
		}
		job, updateErr := service.jobLifecycleRepository.UpdateTranslationJob(txCtx, loaded.job.ID, repository.TranslationJobUpdateDraft{
			JobName:         loaded.job.JobName,
			State:           jobState,
			ProgressPercent: jobProgress,
			StartedAt:       loaded.job.StartedAt,
			FinishedAt:      finishedAt,
		})
		if updateErr != nil {
			return fmt.Errorf("update translation job for body translation phase: %w", updateErr)
		}
		startSnapshot := providerExecutionSnapshot{
			Provider:        loaded.execution.Provider,
			Model:           loaded.execution.Model,
			ExecutionMode:   loaded.execution.ExecutionMode,
			CredentialRef:   loaded.execution.CredentialRef,
			CredentialState: loaded.execution.CredentialState,
			EndpointSummary: providerExecutionEndpointSummary(loaded.execution.EndpointSummary),
		}
		if persistErr := service.persistRuntimeSnapshot(txCtx, loaded.job.ID, "text_translation", startSnapshot); persistErr != nil {
			return persistErr
		}
		createdRun = run
		updatedJob = job
		service.executionSnapshots[run.ID] = startSnapshot
		return nil
	})
	if err != nil {
		return BodyTranslationPhaseCommandReadModel{}, fmt.Errorf("start body translation phase transaction: %w", err)
	}

	loaded.job = updatedJob
	loaded.bodyRun = &createdRun
	loaded.execution.RequestUnitCount = loaded.snapshot.ProviderTargetCount
	if createdRun.State == bodyTranslationPhaseStateCompleted {
		return service.bodyTranslationCommandFromLoaded(loaded, nil, nil), nil
	}
	return service.executeBodyTranslationRun(ctx, loaded, createdRun)
}

const errLoadBodyTranslationPhaseRun = "load body translation phase run: %w"

// PausePhase returns the current phase payload for deferred pause handling.
func (service *BodyTranslationPhaseService) PausePhase(
	ctx context.Context,
	jobID int64,
	phaseRunID int64,
) (BodyTranslationPhaseCommandReadModel, error) {
	return service.transitionBodyTranslationRunState(
		ctx,
		jobID,
		phaseRunID,
		bodyTranslationPhaseStateRunning,
		bodyTranslationPhaseStatePaused,
		bodyTranslationPhaseErrorSummaryRejectPause,
	)
}

// ResumePhase returns the current phase payload for deferred resume handling.
func (service *BodyTranslationPhaseService) ResumePhase(
	ctx context.Context,
	jobID int64,
	phaseRunID int64,
) (BodyTranslationPhaseCommandReadModel, error) {
	return service.transitionBodyTranslationRunState(
		ctx,
		jobID,
		phaseRunID,
		"",
		bodyTranslationPhaseStateRunning,
		bodyTranslationPhaseErrorSummaryRejectResume,
	)
}

// RetryPhase returns the current phase payload for deferred retry handling.
func (service *BodyTranslationPhaseService) RetryPhase(
	ctx context.Context,
	jobID int64,
	phaseRunID int64,
) (BodyTranslationPhaseCommandReadModel, error) {
	loaded, run, err := service.persistBodyTranslationRunStateTransition(
		ctx,
		jobID,
		phaseRunID,
		bodyTranslationPhaseStateRecoverableFail,
		bodyTranslationPhaseStateRunning,
	)
	if err != nil {
		errorSummary := (*BodyTranslationPhaseErrorSummaryReadModel)(nil)
		if loaded.inputSnapshotDrifted {
			errorSummary = bodyTranslationInputSnapshotDriftErrorSummary()
		}
		return service.bodyTranslationCommandFromLoaded(loaded, nil, errorSummary), err
	}
	return service.executeBodyTranslationRun(ctx, loaded, run)
}

// CancelPhase returns the current phase payload for deferred cancel handling.
func (service *BodyTranslationPhaseService) CancelPhase(
	ctx context.Context,
	jobID int64,
	phaseRunID int64,
) (BodyTranslationPhaseCommandReadModel, error) {
	return service.transitionBodyTranslationRunState(
		ctx,
		jobID,
		phaseRunID,
		bodyTranslationPhaseStatePaused,
		bodyTranslationPhaseStateCanceled,
		bodyTranslationPhaseErrorSummaryRejectCancel,
	)
}

// ReadOutputReadiness returns whether body translation output is ready.
func (service *BodyTranslationPhaseService) ReadOutputReadiness(
	ctx context.Context,
	jobID int64,
) (BodyTranslationOutputReadinessReadModel, error) {
	loaded, err := service.loadContext(ctx, jobID)
	if err != nil {
		return BodyTranslationOutputReadinessReadModel{}, err
	}
	return service.buildOutputReadiness(loaded), nil
}

func (service *BodyTranslationPhaseService) loadContext(
	ctx context.Context,
	jobID int64,
) (bodyTranslationLoadedContext, error) {
	job, err := service.jobLifecycleRepository.GetTranslationJobByID(ctx, jobID)
	if err != nil {
		return bodyTranslationLoadedContext{}, fmt.Errorf("load body translation job: %w", err)
	}
	personaRun, err := service.findPhaseRun(ctx, jobID, bodyTranslationPersonaPhaseType)
	if err != nil {
		return bodyTranslationLoadedContext{}, fmt.Errorf("load persona generation phase run: %w", err)
	}
	bodyRun, err := service.findPhaseRun(ctx, jobID, bodyTranslationPhaseType)
	if err != nil {
		return bodyTranslationLoadedContext{}, fmt.Errorf(errLoadBodyTranslationPhaseRun, err)
	}
	records, fieldsByRecord, err := service.loadBodyTranslationSourceFields(ctx, job)
	if err != nil {
		return bodyTranslationLoadedContext{}, err
	}
	dictionary, err := service.foundationDataRepository.ListDictionaryEntries(ctx, &job.ID, "job", "", "")
	if err != nil {
		return bodyTranslationLoadedContext{}, fmt.Errorf("list body translation dictionary entries: %w", err)
	}
	personas, err := service.foundationDataRepository.ListPersonasByTranslationJobID(ctx, job.ID)
	if err != nil {
		return bodyTranslationLoadedContext{}, fmt.Errorf("list body translation personas: %w", err)
	}
	outputFields, err := service.jobOutputRepository.ListJobTranslationFieldsByJobID(ctx, job.ID)
	if err != nil {
		return bodyTranslationLoadedContext{}, fmt.Errorf("list body translation outputs: %w", err)
	}
	execution := bodyTranslationExecutionFromPhaseRuns(bodyRun, personaRun)
	if snapshotReader, ok := service.jobLifecycleRepository.(bodyTranslationPhaseRuntimeSnapshotReader); ok && bodyRun == nil {
		snapshot, snapshotErr := snapshotReader.GetTranslationJobPhaseRuntimeSnapshot(ctx, job.ID, "text_translation")
		switch {
		case snapshotErr == nil:
			execution.CredentialRef = strings.TrimSpace(snapshot.CredentialRef)
			execution.Provider = strings.TrimSpace(snapshot.Provider)
			execution.Model = strings.TrimSpace(snapshot.ModelName)
			execution.ExecutionMode = strings.TrimSpace(snapshot.ExecutionMode)
		case errors.Is(snapshotErr, repository.ErrNotFound):
			execution.CredentialRef = ""
			execution.Provider = ""
			execution.Model = ""
			execution.ExecutionMode = ""
		default:
			return bodyTranslationLoadedContext{}, fmt.Errorf("load body translation phase runtime snapshot: %w", snapshotErr)
		}
	}
	if bodyRun != nil {
		if snapshot, ok := service.executionSnapshots[bodyRun.ID]; ok {
			execution.CredentialState = snapshot.CredentialState
			execution.EndpointSummary = providerExecutionEndpointSummary(snapshot.EndpointSummary)
		} else if snapshotReader, ok := service.jobLifecycleRepository.(bodyTranslationPhaseRuntimeSnapshotReader); ok {
			snapshot, snapshotErr := snapshotReader.GetTranslationJobPhaseRuntimeSnapshot(ctx, job.ID, "text_translation")
			switch {
			case snapshotErr == nil:
				persisted := providerExecutionSnapshotFromRuntimeSnapshot(snapshot)
				execution.CredentialState = persisted.CredentialState
				execution.EndpointSummary = providerExecutionEndpointSummary(persisted.EndpointSummary)
			case errors.Is(snapshotErr, repository.ErrNotFound):
			default:
				return bodyTranslationLoadedContext{}, fmt.Errorf("load body translation phase runtime snapshot: %w", snapshotErr)
			}
		}
	}
	persona := firstBodyTranslationPersona(personas)
	snapshot, err := buildBodyTranslationInputSnapshot(job, records, fieldsByRecord, dictionary, persona, execution)
	if err != nil {
		return bodyTranslationLoadedContext{}, fmt.Errorf("build body translation input snapshot: %w", err)
	}
	inputSnapshotDrifted := bodyRun != nil && bodyTranslationRunSnapshotDrifted(*bodyRun, snapshot)
	if bodyRun != nil {
		snapshot = applyBodyTranslationRunSnapshot(*bodyRun, snapshot)
	}
	execution.RequestUnitCount = snapshot.ProviderTargetCount
	execution.OutputCount = len(outputFields)
	return bodyTranslationLoadedContext{
		job:                  job,
		bodyRun:              bodyRun,
		personaRun:           personaRun,
		records:              records,
		fieldsByRecord:       fieldsByRecord,
		dictionary:           dictionary,
		persona:              persona,
		execution:            execution,
		snapshot:             snapshot,
		inputSnapshotDrifted: inputSnapshotDrifted,
		outputFields:         outputFields,
	}, nil
}

func (service *BodyTranslationPhaseService) loadBodyTranslationSourceFields(
	ctx context.Context,
	job repository.TranslationJob,
) ([]repository.TranslationRecord, map[int64][]repository.TranslationField, error) {
	_, err := service.translationSourceReader.GetXEditExtractedDataByID(ctx, job.XEditExtractedDataID)
	if err != nil {
		return nil, nil, fmt.Errorf("load body translation xedit source: %w", err)
	}
	records, err := service.translationSourceReader.ListTranslationRecordsByXEditID(ctx, job.XEditExtractedDataID)
	if err != nil {
		return nil, nil, fmt.Errorf("list body translation records: %w", err)
	}
	sort.SliceStable(records, func(left int, right int) bool {
		return records[left].ID < records[right].ID
	})
	fieldsByRecord := make(map[int64][]repository.TranslationField, len(records))
	for _, record := range records {
		fields, fieldErr := service.translationSourceReader.ListTranslationFieldsByTranslationRecordID(ctx, record.ID)
		if fieldErr != nil {
			return nil, nil, fmt.Errorf("list body translation fields for record %d: %w", record.ID, fieldErr)
		}
		fieldsByRecord[record.ID] = append([]repository.TranslationField(nil), fields...)
	}
	return records, fieldsByRecord, nil
}

func bodyTranslationExecutionFromPhaseRuns(
	bodyRun *repository.JobPhaseRun,
	personaRun *repository.JobPhaseRun,
) BodyTranslationPhaseExecutionSummaryReadModel {
	execution := bodyTranslationExecutionSummaryFromRun(bodyRun)
	if execution.Provider == "" && execution.Model == "" && execution.ExecutionMode == "" && execution.CredentialRef == "" {
		return bodyTranslationExecutionSummaryFromRun(personaRun)
	}
	return execution
}

func firstBodyTranslationPersona(personas []repository.Persona) repository.Persona {
	if len(personas) == 0 {
		return repository.Persona{}
	}
	return personas[0]
}

func (service *BodyTranslationPhaseService) persistRuntimeSnapshot(
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
		return fmt.Errorf("persist body translation runtime snapshot: %w", err)
	}
	if _, err := store.SaveTranslationJobPhaseRuntimeSnapshot(ctx, repository.TranslationJobPhaseRuntimeSnapshotDraft{
		TranslationJobID:     jobID,
		PhaseID:              phaseID,
		Provider:             snapshot.Provider,
		ModelName:            snapshot.Model,
		CredentialRef:        snapshot.CredentialRef,
		CredentialStatus:     snapshot.CredentialState,
		EndpointSummary:      providerExecutionOptionalString(snapshot.EndpointSummary),
		ExecutionMode:        snapshot.ExecutionMode,
		BatchMode:            providerExecutionBatchMode(current),
		ModelListSourceToken: providerExecutionModelListSourceToken(current),
	}); err != nil {
		return fmt.Errorf("persist body translation runtime snapshot: %w", err)
	}
	return nil
}

func (service *BodyTranslationPhaseService) startRejection(
	loaded bodyTranslationLoadedContext,
) *bodyTranslationStartRejection {
	if loaded.personaRun == nil || strings.TrimSpace(loaded.personaRun.State) != bodyTranslationPhaseStateCompleted {
		return &bodyTranslationStartRejection{
			errorKind: "persona_phase_incomplete",
			reason:    "persona phase must be completed",
		}
	}
	if isBodyTranslationTerminalJob(loaded.job.State) {
		return &bodyTranslationStartRejection{
			errorKind: "terminal_job",
			reason:    "terminal job",
		}
	}
	if !bodyTranslationExecutionConfigured(loaded.execution) {
		return &bodyTranslationStartRejection{
			errorKind: "persona_phase_incomplete",
			reason:    "phase runtime snapshot is missing",
		}
	}
	if loaded.bodyRun != nil && isBodyTranslationActiveRunState(strings.TrimSpace(loaded.bodyRun.State)) {
		return &bodyTranslationStartRejection{
			errorKind: "active_phase_exists",
			reason:    "body translation phase run already exists",
		}
	}
	return nil
}

func bodyTranslationExecutionConfigured(execution BodyTranslationPhaseExecutionSummaryReadModel) bool {
	return strings.TrimSpace(execution.Provider) != "" &&
		strings.TrimSpace(execution.Model) != "" &&
		strings.TrimSpace(execution.ExecutionMode) != ""
}

func (service *BodyTranslationPhaseService) resolveExecutionSnapshotForStart(
	ctx context.Context,
	execution BodyTranslationPhaseExecutionSummaryReadModel,
) (BodyTranslationPhaseExecutionSummaryReadModel, *bodyTranslationStartRejection, error) {
	if service.providerSettings == nil || !providerExecutionUsesProviderSettings(execution.Provider) {
		return execution, nil, nil
	}
	resolved, err := service.providerSettings.ResolveProviderExecutionSettings(ctx, ProviderSettingsResolveInput{
		ConsumerID:          "body_translation_phase",
		AllowSecretSnapshot: true,
		Selection: ProviderSettingsResolveSelection{
			ProviderID:      execution.Provider,
			Model:           execution.Model,
			ExecutionMethod: execution.ExecutionMode,
			UseBatchAPI:     false,
		},
	})
	if err != nil {
		return BodyTranslationPhaseExecutionSummaryReadModel{}, nil, fmt.Errorf("resolve body translation provider settings: %w", err)
	}
	execution.CredentialRef = providerExecutionOptionalString(resolved.CredentialReferenceID)
	execution.CredentialState = strings.TrimSpace(resolved.CredentialState)
	execution.EndpointSummary = providerExecutionEndpointSummary(resolved.Endpoint)
	if resolved.ErrorKind == nil {
		return execution, nil, nil
	}
	switch strings.TrimSpace(*resolved.ErrorKind) {
	case providerSettingsErrorKindCredentialMissing:
		return execution, &bodyTranslationStartRejection{
			errorKind: "persona_phase_incomplete",
			reason:    "phase runtime snapshot is missing",
		}, nil
	case providerSettingsErrorKindEndpointMissing:
		return execution, &bodyTranslationStartRejection{
			errorKind: "persona_phase_incomplete",
			reason:    "phase runtime snapshot is missing",
		}, nil
	default:
		return execution, &bodyTranslationStartRejection{
			errorKind: "provider_failure",
			reason:    "provider execution failed",
		}, nil
	}
}

func (service *BodyTranslationPhaseService) rejectedCommand(
	loaded bodyTranslationLoadedContext,
	rejection bodyTranslationStartRejection,
) BodyTranslationPhaseCommandReadModel {
	outputReadiness := service.buildOutputReadiness(loaded)
	return BodyTranslationPhaseCommandReadModel{
		JobID:               loaded.job.ID,
		CurrentPhase:        bodyTranslationCurrentPhase,
		PhaseState:          bodyTranslationPhaseStateIdleReady,
		PhaseRunID:          nil,
		StartedAt:           nil,
		FinishedAt:          nil,
		Progress:            service.buildProgress(bodyTranslationPhaseStateIdleReady, loaded.snapshot, loaded.outputFields),
		InputSnapshotDigest: "",
		InputSummary:        toBodyTranslationInputSummaryReadModel(loaded.snapshot),
		RequestSummary:      toBodyTranslationRequestSummaryReadModel(loaded.snapshot),
		Execution:           loaded.execution,
		FieldResultSummary:  service.buildFieldResultSummary(loaded.outputFields),
		ResultSummary:       service.buildFieldResultSummary(loaded.outputFields),
		FieldResults:        service.buildFieldResultItems(loaded),
		Retryable:           false,
		OutputReadiness:     outputReadiness,
		ErrorSummary: &BodyTranslationPhaseErrorSummaryReadModel{
			ErrorKind:  rejection.errorKind,
			Reason:     rejection.reason,
			Retryable:  false,
			IsRedacted: true,
		},
	}
}

func (service *BodyTranslationPhaseService) buildProgress(
	phaseState string,
	snapshot bodyTranslationInputSnapshot,
	outputFields []repository.JobTranslationField,
) BodyTranslationPhaseProgressReadModel {
	processedCount := 0
	translatedCount := 0
	for _, outputField := range outputFields {
		processedCount++
		switch strings.TrimSpace(outputField.OutputStatus) {
		case bodyTranslationPhaseStateCompleted, bodyTranslationOutputStatusReady, bodyTranslationOutputStatusTranslated:
			translatedCount++
		}
	}
	percent := 0
	if snapshot.ProviderTargetCount == 0 && phaseState == bodyTranslationPhaseStateCompleted {
		percent = 100
	} else if snapshot.ProviderTargetCount > 0 {
		percent = (processedCount * 100) / snapshot.ProviderTargetCount
		if percent > 100 {
			percent = 100
		}
	}
	return BodyTranslationPhaseProgressReadModel{
		Percent:         percent,
		ProcessedCount:  processedCount,
		TotalCount:      snapshot.TargetCount,
		TargetCount:     snapshot.TargetCount,
		TranslatedCount: translatedCount,
		SkippedCount:    snapshot.ExactExclusionCount,
		CurrentStep:     phaseState,
	}
}

func (service *BodyTranslationPhaseService) buildFieldResultSummary(
	outputFields []repository.JobTranslationField,
) *BodyTranslationPhaseFieldResultSummaryReadModel {
	if len(outputFields) == 0 {
		return &BodyTranslationPhaseFieldResultSummaryReadModel{}
	}
	summary := &BodyTranslationPhaseFieldResultSummaryReadModel{}
	for _, outputField := range outputFields {
		summary.OutputCount++
		switch strings.TrimSpace(outputField.OutputStatus) {
		case bodyTranslationPhaseStateCompleted, bodyTranslationOutputStatusReady, bodyTranslationOutputStatusTranslated:
			summary.TranslatedCount++
			if strings.TrimSpace(outputField.OutputStatus) != bodyTranslationOutputStatusTranslated {
				summary.OutputReadyCount++
			}
		case "failed":
			summary.FailedCount++
		case "skipped":
			summary.SkippedCount++
		case "protection_failed":
			summary.ProtectionFailedCount++
		}
	}
	return summary
}

func (service *BodyTranslationPhaseService) buildFieldResultItems(
	loaded bodyTranslationLoadedContext,
) []BodyTranslationPhaseFieldResultItemReadModel {
	if len(loaded.outputFields) == 0 {
		return nil
	}
	outputFields := append([]repository.JobTranslationField(nil), loaded.outputFields...)
	sort.SliceStable(outputFields, func(left int, right int) bool {
		return outputFields[left].TranslationFieldID < outputFields[right].TranslationFieldID
	})
	items := make([]BodyTranslationPhaseFieldResultItemReadModel, 0, len(outputFields))
	for _, outputField := range outputFields {
		items = append(
			items,
			bodyTranslationFieldResultItemFromOutput(
				outputField,
				loaded.fieldByID(outputField.TranslationFieldID),
				loaded.recordByFieldID(outputField.TranslationFieldID),
			),
		)
	}
	return items
}

func bodyTranslationFieldResultItemFromOutput(
	outputField repository.JobTranslationField,
	field repository.TranslationField,
	record repository.TranslationRecord,
) BodyTranslationPhaseFieldResultItemReadModel {
	outputStatus := strings.TrimSpace(outputField.OutputStatus)
	fieldLabel := bodyTranslationFieldLabel(record, field)
	return BodyTranslationPhaseFieldResultItemReadModel{
		Identity: BodyTranslationPhaseFieldIdentityReadModel{
			TranslationFieldID:      outputField.TranslationFieldID,
			PhaseTranslationFieldID: outputField.ID,
			RecordType:              strings.TrimSpace(record.RecordType),
			FieldType:               strings.TrimSpace(field.SubrecordType),
			FormID:                  strings.TrimSpace(record.FormID),
			EditorID:                strings.TrimSpace(record.EditorID),
			FieldLabel:              fieldLabel,
		},
		FieldID:                     outputField.TranslationFieldID,
		FieldLabel:                  fieldLabel,
		SourceExcerpt:               bodyTranslationSourceExcerpt(field.SourceText),
		TranslatedText:              outputField.TranslatedText,
		OutputStatus:                outputStatus,
		ProtectionValidationResult:  bodyTranslationProtectionValidationResult(outputStatus),
		ProtectionValidationSummary: bodyTranslationProtectionValidationResult(outputStatus),
		RetryCount:                  outputField.RetryCount,
	}
}

func bodyTranslationFieldLabel(record repository.TranslationRecord, field repository.TranslationField) string {
	parts := []string{
		strings.TrimSpace(record.EditorID),
		strings.TrimSpace(record.FormID),
		strings.TrimSpace(field.SubrecordType),
	}
	labels := make([]string, 0, len(parts))
	for _, part := range parts {
		if part != "" {
			labels = append(labels, part)
		}
	}
	return strings.Join(labels, " / ")
}

func bodyTranslationSourceExcerpt(sourceText string) string {
	trimmed := strings.TrimSpace(sourceText)
	if trimmed == "" {
		return ""
	}
	runes := []rune(trimmed)
	if len(runes) <= 120 {
		return trimmed
	}
	return string(runes[:120])
}

func bodyTranslationProtectionValidationResult(outputStatus string) string {
	switch strings.TrimSpace(outputStatus) {
	case "protection_failed":
		return "failed"
	case bodyTranslationOutputStatusSkipped:
		return "not_applicable"
	case bodyTranslationOutputStatusReady, bodyTranslationOutputStatusTranslated:
		return "passed"
	default:
		return "unknown"
	}
}

func (service *BodyTranslationPhaseService) buildLoadedFieldResultSummary(
	loaded bodyTranslationLoadedContext,
) *BodyTranslationPhaseFieldResultSummaryReadModel {
	summary := service.buildFieldResultSummary(loaded.outputFields)
	if summary == nil {
		summary = &BodyTranslationPhaseFieldResultSummaryReadModel{}
	}
	summary.FieldResults = service.buildFieldResultItems(loaded)
	if loaded.bodyRun == nil || strings.TrimSpace(loaded.bodyRun.State) != bodyTranslationPhaseStateRecoverableFail {
		return summary
	}
	implicitFailedCount := loaded.snapshot.ProviderTargetCount - summary.OutputCount
	if implicitFailedCount <= 0 {
		return summary
	}
	cloned := *summary
	cloned.FieldResults = append([]BodyTranslationPhaseFieldResultItemReadModel(nil), summary.FieldResults...)
	cloned.FailedCount += implicitFailedCount
	return &cloned
}

func (service *BodyTranslationPhaseService) buildActionEnablement(
	loaded bodyTranslationLoadedContext,
	rejection *bodyTranslationStartRejection,
) BodyTranslationPhaseActionEnablementReadModel {
	result := BodyTranslationPhaseActionEnablementReadModel{
		CanStart:                rejection == nil,
		CanCheckOutputReadiness: loaded.bodyRun != nil,
	}
	if rejection != nil {
		result.StartBlockedReason = cloneBodyTranslationStringPointer(&rejection.reason)
	}
	if loaded.bodyRun == nil {
		return result
	}
	state := strings.TrimSpace(loaded.bodyRun.State)
	if state != bodyTranslationPhaseStateRunning {
		reason := "body translation phase is not running"
		result.PauseBlockedReason = &reason
	} else {
		result.CanPause = true
	}
	if state != bodyTranslationPhaseStatePaused && state != bodyTranslationPhaseStateRecoverableFail {
		reason := "body translation phase is not resumable"
		result.ResumeBlockedReason = &reason
	} else {
		result.CanResume = true
	}
	if state != bodyTranslationPhaseStateRecoverableFail {
		reason := "body translation phase is not retryable"
		result.RetryBlockedReason = &reason
	} else {
		result.CanRetry = true
	}
	if loaded.inputSnapshotDrifted {
		reason := bodyTranslationInputSnapshotDriftReason
		result.CanResume = false
		result.ResumeBlockedReason = &reason
		result.CanRetry = false
		result.RetryBlockedReason = &reason
	}
	if state != bodyTranslationPhaseStatePaused {
		reason := "body translation phase is not cancelable"
		result.CancelBlockedReason = &reason
	} else {
		result.CanCancel = true
	}
	return result
}

func (service *BodyTranslationPhaseService) buildOutputReadiness(
	loaded bodyTranslationLoadedContext,
) BodyTranslationOutputReadinessReadModel {
	result := BodyTranslationOutputReadinessReadModel{
		JobID:               loaded.job.ID,
		CurrentPhase:        bodyTranslationCurrentPhase,
		PhaseState:          bodyTranslationPhaseStateIdleReady,
		Ready:               false,
		BlockedReason:       "body phase is not completed",
		ErrorKind:           "output_readiness_blocked",
		CompletedFieldCount: 0,
		StatusConsistent:    false,
		OutputCount:         len(loaded.outputFields),
	}
	if loaded.bodyRun != nil {
		result.PhaseState = strings.TrimSpace(loaded.bodyRun.State)
	}
	if loaded.bodyRun == nil {
		return result
	}
	if strings.TrimSpace(loaded.bodyRun.State) != bodyTranslationPhaseStateCompleted {
		return result
	}
	if loaded.job.State != bodyTranslationJobStateCompleted {
		result.BlockedReason = "translation job is not completed"
		return result
	}
	completedFieldCount, statusConsistent := service.evaluateOutputStatusConsistency(loaded)
	result.CompletedFieldCount = completedFieldCount
	result.StatusConsistent = statusConsistent
	if !statusConsistent {
		result.BlockedReason = "body translation output status is inconsistent"
		return result
	}
	result.Ready = true
	result.BlockedReason = ""
	result.ErrorKind = ""
	return result
}

func (service *BodyTranslationPhaseService) buildPhaseErrorSummary(
	loaded bodyTranslationLoadedContext,
	outputReadiness BodyTranslationOutputReadinessReadModel,
) *BodyTranslationPhaseErrorSummaryReadModel {
	if loaded.bodyRun == nil {
		return nil
	}
	errorKind := strings.TrimSpace(loaded.bodyRun.LatestError)
	if errorKind == "" {
		return nil
	}
	retryable := false
	reason := "body translation phase reported an error"
	switch errorKind {
	case "provider_failure":
		retryable = true
		reason = bodyTranslationFieldResultReasonProviderFailure
	case "invalid_provider_response":
		retryable = true
		reason = bodyTranslationFieldResultReasonInvalidProviderResponse
	case "output_readiness_blocked":
		reason = outputReadiness.BlockedReason
	case "protection_validation_failed":
		reason = bodyTranslationFieldResultReasonProtectionValidationFailed
	case "save_failed":
		reason = "body translation result persistence failed"
	case "late_response_rejected":
		retryable = true
		reason = "late provider response was rejected"
	}
	return &BodyTranslationPhaseErrorSummaryReadModel{
		ErrorKind:  errorKind,
		Reason:     reason,
		Retryable:  retryable,
		IsRedacted: true,
	}
}

func (service *BodyTranslationPhaseService) evaluateOutputStatusConsistency(
	loaded bodyTranslationLoadedContext,
) (int, bool) {
	if loaded.snapshot.ProviderTargetCount == 0 {
		return 0, true
	}
	completedFieldCount := 0
	for _, outputField := range loaded.outputFields {
		switch strings.TrimSpace(outputField.OutputStatus) {
		case bodyTranslationOutputStatusReady, bodyTranslationOutputStatusTranslated, bodyTranslationOutputStatusSkipped:
			completedFieldCount++
		default:
			return completedFieldCount, false
		}
	}
	if completedFieldCount != loaded.snapshot.ProviderTargetCount {
		return completedFieldCount, false
	}
	return completedFieldCount, true
}

func (service *BodyTranslationPhaseService) transitionBodyTranslationRunState(
	ctx context.Context,
	jobID int64,
	phaseRunID int64,
	requiredState string,
	nextState string,
	reject func() *BodyTranslationPhaseErrorSummaryReadModel,
) (BodyTranslationPhaseCommandReadModel, error) {
	loaded, updatedRun, err := service.persistBodyTranslationRunStateTransition(
		ctx,
		jobID,
		phaseRunID,
		requiredState,
		nextState,
	)
	if err != nil {
		errorSummary := reject()
		if loaded.inputSnapshotDrifted {
			errorSummary = bodyTranslationInputSnapshotDriftErrorSummary()
		}
		return service.bodyTranslationCommandFromLoaded(loaded, nil, errorSummary), err
	}
	loaded.bodyRun = &updatedRun
	reloaded, reloadErr := service.loadContext(ctx, jobID)
	if reloadErr != nil {
		return BodyTranslationPhaseCommandReadModel{}, reloadErr
	}
	return service.bodyTranslationCommandFromLoaded(reloaded, nil, nil), nil
}

func (service *BodyTranslationPhaseService) persistBodyTranslationRunStateTransition(
	ctx context.Context,
	jobID int64,
	phaseRunID int64,
	requiredState string,
	nextState string,
) (bodyTranslationLoadedContext, repository.JobPhaseRun, error) {
	loaded, err := service.loadContext(ctx, jobID)
	if err != nil {
		return bodyTranslationLoadedContext{}, repository.JobPhaseRun{}, err
	}
	if loaded.bodyRun == nil || loaded.bodyRun.ID != phaseRunID {
		return loaded, repository.JobPhaseRun{}, fmt.Errorf(errLoadBodyTranslationPhaseRun, repository.ErrNotFound)
	}
	if nextState == bodyTranslationPhaseStateRunning && loaded.inputSnapshotDrifted {
		return loaded, *loaded.bodyRun, errors.New(bodyTranslationInputSnapshotDriftReason)
	}
	if !bodyTranslationRunStateTransitionAllowed(strings.TrimSpace(loaded.bodyRun.State), requiredState) {
		return loaded, *loaded.bodyRun, fmt.Errorf("body translation phase state transition rejected")
	}
	if service.transactor == nil {
		return loaded, repository.JobPhaseRun{}, fmt.Errorf("body translation phase state transition: transactor is not configured")
	}
	now := service.now()
	var updatedRun repository.JobPhaseRun
	var updatedJob repository.TranslationJob
	err = service.transactor.WithTransaction(ctx, func(txCtx context.Context) error {
		runDraft, jobDraft := bodyTranslationRunStateTransitionDrafts(loaded, nextState, now)
		run, updateErr := service.jobLifecycleRepository.UpdateJobPhaseRun(txCtx, loaded.bodyRun.ID, runDraft)
		if updateErr != nil {
			return fmt.Errorf("update body translation phase run state: %w", updateErr)
		}
		job, updateErr := service.jobLifecycleRepository.UpdateTranslationJob(txCtx, loaded.job.ID, jobDraft)
		if updateErr != nil {
			return fmt.Errorf("update body translation phase job state: %w", updateErr)
		}
		updatedRun = run
		updatedJob = job
		return nil
	})
	if err != nil {
		return loaded, repository.JobPhaseRun{}, fmt.Errorf("persist body translation phase state transition: %w", err)
	}
	loaded.job = updatedJob
	loaded.bodyRun = &updatedRun
	return loaded, updatedRun, nil
}

func bodyTranslationRunStateTransitionAllowed(currentState string, requiredState string) bool {
	if requiredState != "" {
		return currentState == requiredState
	}
	return currentState == bodyTranslationPhaseStatePaused || currentState == bodyTranslationPhaseStateRecoverableFail
}

func bodyTranslationRunStateTransitionDrafts(
	loaded bodyTranslationLoadedContext,
	nextState string,
	now time.Time,
) (repository.JobPhaseRunUpdateDraft, repository.TranslationJobUpdateDraft) {
	finishedAt := loaded.bodyRun.FinishedAt
	jobState := loaded.job.State
	jobFinishedAt := loaded.job.FinishedAt
	latestError := loaded.bodyRun.LatestError
	if nextState == bodyTranslationPhaseStateRunning {
		finishedAt = nil
		jobState = bodyTranslationJobStateRunning
		jobFinishedAt = nil
		latestError = ""
	}
	if nextState == bodyTranslationPhaseStateCanceled {
		finishedAt = &now
		jobState = bodyTranslationJobStateCanceled
		jobFinishedAt = &now
		latestError = ""
	}
	return repository.JobPhaseRunUpdateDraft{
			State:               nextState,
			ProgressPercent:     loaded.bodyRun.ProgressPercent,
			LatestExternalRunID: loaded.bodyRun.LatestExternalRunID,
			LatestError:         latestError,
			StartedAt:           loaded.bodyRun.StartedAt,
			FinishedAt:          finishedAt,
		}, repository.TranslationJobUpdateDraft{
			JobName:         loaded.job.JobName,
			State:           jobState,
			ProgressPercent: loaded.job.ProgressPercent,
			StartedAt:       loaded.job.StartedAt,
			FinishedAt:      jobFinishedAt,
		}
}

func bodyTranslationPhaseErrorSummaryRejectCancel() *BodyTranslationPhaseErrorSummaryReadModel {
	return &BodyTranslationPhaseErrorSummaryReadModel{
		ErrorKind:  "output_readiness_blocked",
		Reason:     "body translation phase is not cancelable",
		Retryable:  false,
		IsRedacted: true,
	}
}

func bodyTranslationPhaseErrorSummaryRejectPause() *BodyTranslationPhaseErrorSummaryReadModel {
	return &BodyTranslationPhaseErrorSummaryReadModel{
		ErrorKind:  "output_readiness_blocked",
		Reason:     "body translation phase is not running",
		Retryable:  false,
		IsRedacted: true,
	}
}

func bodyTranslationPhaseErrorSummaryRejectResume() *BodyTranslationPhaseErrorSummaryReadModel {
	return &BodyTranslationPhaseErrorSummaryReadModel{
		ErrorKind:  "output_readiness_blocked",
		Reason:     "body translation phase is not resumable",
		Retryable:  false,
		IsRedacted: true,
	}
}

func (service *BodyTranslationPhaseService) executeBodyTranslationRun(
	ctx context.Context,
	loaded bodyTranslationLoadedContext,
	run repository.JobPhaseRun,
) (BodyTranslationPhaseCommandReadModel, error) {
	if loaded.inputSnapshotDrifted {
		return service.bodyTranslationCommandFromLoaded(loaded, nil, bodyTranslationInputSnapshotDriftErrorSummary()), nil
	}
	if service.bodyTranslationProviderUnavailable() {
		updatedLoaded, err := service.persistBodyTranslationRunFailure(ctx, loaded, run, "provider_failure")
		if err != nil {
			return BodyTranslationPhaseCommandReadModel{}, err
		}
		return service.bodyTranslationCommandFromLoaded(updatedLoaded, nil, bodyTranslationProviderFailureSummary()), nil
	}
	pendingTargets := service.pendingBodyTranslationTargets(loaded)
	if len(pendingTargets) == 0 {
		return service.reloadedBodyTranslationCommand(ctx, loaded.job.ID)
	}
	for _, target := range pendingTargets {
		command, err := service.executeBodyTranslationTarget(ctx, loaded, run, target)
		if err != nil {
			return command, err
		}
		if command.PhaseState != bodyTranslationPhaseStateRunning && command.PhaseState != bodyTranslationPhaseStateCompleted {
			return command, nil
		}
		loaded, err = service.loadContext(ctx, loaded.job.ID)
		if err != nil {
			return BodyTranslationPhaseCommandReadModel{}, err
		}
		if loaded.bodyRun == nil {
			return BodyTranslationPhaseCommandReadModel{}, fmt.Errorf(errLoadBodyTranslationPhaseRun, repository.ErrNotFound)
		}
		run = *loaded.bodyRun
	}
	return service.reloadedBodyTranslationCommand(ctx, loaded.job.ID)
}

func (service *BodyTranslationPhaseService) bodyTranslationProviderUnavailable() bool {
	return service.bodyTranslationProvider == nil
}

func bodyTranslationProviderFailureSummary() *BodyTranslationPhaseErrorSummaryReadModel {
	return &BodyTranslationPhaseErrorSummaryReadModel{
		ErrorKind:  "provider_failure",
		Reason:     bodyTranslationFieldResultReasonProviderFailure,
		Retryable:  false,
		IsRedacted: true,
	}
}

func (service *BodyTranslationPhaseService) executeBodyTranslationTarget(
	ctx context.Context,
	loaded bodyTranslationLoadedContext,
	run repository.JobPhaseRun,
	target normalizedBodyTranslationFieldTarget,
) (BodyTranslationPhaseCommandReadModel, error) {
	providerRequest := buildBodyTranslationProviderRequest(loaded, target)
	providerResult := service.bodyTranslationProvider.TranslateBodyField(ctx, providerRequest)
	return service.PersistBodyTranslationFieldResults(ctx, BodyTranslationFieldResultPersistenceRequest{
		TranslationJobID: loaded.job.ID,
		PhaseRunID:       run.ID,
		TargetFields: []BodyTranslationFieldResultTarget{
			{
				TranslationFieldID:    target.TranslationFieldID,
				FieldCorrelationKey:   target.FieldCorrelationKey,
				OutputStatusCandidate: bodyTranslationOutputStatusReady,
				ProtectedElements:     target.ProtectedElements,
			},
		},
		ProviderResults: []BodyTranslationProviderResult{providerResult},
	})
}

func (service *BodyTranslationPhaseService) reloadedBodyTranslationCommand(
	ctx context.Context,
	jobID int64,
) (BodyTranslationPhaseCommandReadModel, error) {
	reloaded, err := service.loadContext(ctx, jobID)
	if err != nil {
		return BodyTranslationPhaseCommandReadModel{}, err
	}
	return service.bodyTranslationCommandFromLoaded(reloaded, nil, nil), nil
}

func (service *BodyTranslationPhaseService) pendingBodyTranslationTargets(
	loaded bodyTranslationLoadedContext,
) []normalizedBodyTranslationFieldTarget {
	existing := indexBodyTranslationOutputFields(loaded.outputFields)
	targets := make([]normalizedBodyTranslationFieldTarget, 0)
	for _, field := range loaded.snapshot.Fields {
		if !field.IncludedInProviderRequests {
			continue
		}
		if _, ok := existing[field.TranslationFieldID]; ok {
			continue
		}
		repositoryField := loaded.fieldByID(field.TranslationFieldID)
		if repositoryField.ID == 0 {
			continue
		}
		targets = append(targets, normalizedBodyTranslationFieldTarget{
			TranslationFieldID:  field.TranslationFieldID,
			FieldCorrelationKey: fmt.Sprintf("field:%d", field.TranslationFieldID),
			OutputStatus:        bodyTranslationOutputStatusReady,
			ProtectedElements:   bodyTranslationProtectedElementsFromSource(field.SourceText),
			Field:               repositoryField,
			RecordType:          field.RecordType,
		})
	}
	return targets
}

func buildBodyTranslationProviderRequest(
	loaded bodyTranslationLoadedContext,
	target normalizedBodyTranslationFieldTarget,
) BodyTranslationProviderRequest {
	snapshotField := loaded.snapshotFieldByID(target.TranslationFieldID)
	return BodyTranslationProviderRequest{
		Provider:                loaded.execution.Provider,
		Model:                   loaded.execution.Model,
		ExecutionMode:           loaded.execution.ExecutionMode,
		CredentialRef:           loaded.execution.CredentialRef,
		EndpointSummary:         providerExecutionEndpointSummary(loaded.execution.EndpointSummary),
		RequestUnitID:           fmt.Sprintf("body-field-%d", target.TranslationFieldID),
		FieldCorrelationKey:     target.FieldCorrelationKey,
		RecordType:              target.RecordType,
		FieldType:               target.Field.SubrecordType,
		SourceText:              bodyTranslationPromptSourceText(target.Field.SourceText),
		SourceLanguage:          bodyTranslationSourceLanguageDefaultValue,
		TargetLanguage:          bodyTranslationTargetLanguageDefaultValue,
		PersonaSummary:          bodyTranslationPersonaSummary(loaded.persona),
		CompleteMatchExclusions: append([]BodyTranslationDictionaryExactMatchExclusion(nil), snapshotField.CompleteMatchExclusions...),
		PartialMatchConstraints: append([]BodyTranslationPartialMatchConstraint(nil), snapshotField.PartialMatchConstraints...),
		ProtectedElements:       append([]BodyTranslationProtectedElement(nil), target.ProtectedElements...),
	}
}

func (loaded bodyTranslationLoadedContext) fieldByID(fieldID int64) repository.TranslationField {
	for _, fields := range loaded.fieldsByRecord {
		for _, field := range fields {
			if field.ID == fieldID {
				return field
			}
		}
	}
	return repository.TranslationField{}
}

func (loaded bodyTranslationLoadedContext) recordByFieldID(fieldID int64) repository.TranslationRecord {
	for _, record := range loaded.records {
		for _, field := range loaded.fieldsByRecord[record.ID] {
			if field.ID == fieldID {
				return record
			}
		}
	}
	return repository.TranslationRecord{}
}

func (loaded bodyTranslationLoadedContext) snapshotFieldByID(fieldID int64) bodyTranslationSnapshotField {
	for _, field := range loaded.snapshot.Fields {
		if field.TranslationFieldID == fieldID {
			return field
		}
	}
	return bodyTranslationSnapshotField{}
}

func bodyTranslationProtectedElementsFromSource(sourceText string) []BodyTranslationProtectedElement {
	matches := bodyTranslationProtectedElementPattern.FindAllString(sourceText, -1)
	elements := make([]BodyTranslationProtectedElement, 0, len(matches))
	for _, match := range matches {
		elementType := "token"
		if strings.HasPrefix(match, "<") {
			elementType = "tag"
		}
		if strings.HasPrefix(match, "{") {
			elementType = "placeholder"
		}
		elements = append(elements, BodyTranslationProtectedElement{
			ElementType: elementType,
			SourceText:  match,
			Digest:      bodyTranslationDigestLines([]string{elementType + ":" + match}),
		})
	}
	return elements
}

func applyBodyTranslationRunSnapshot(
	run repository.JobPhaseRun,
	snapshot bodyTranslationInputSnapshot,
) bodyTranslationInputSnapshot {
	if strings.TrimSpace(run.InputSnapshotDigest) == "" {
		return snapshot
	}
	snapshot.TargetCount = run.SnapshotFieldCount
	snapshot.ProviderTargetCount = run.ProviderTargetCount
	snapshot.ExactExclusionCount = run.ExactExclusionCount
	snapshot.PartialConstraintCount = run.PartialConstraintCount
	snapshot.InputSnapshotDigest = strings.TrimSpace(run.InputSnapshotDigest)
	snapshot.DictionaryDigest = strings.TrimSpace(run.DictionaryDigest)
	snapshot.PersonaDigest = strings.TrimSpace(run.PersonaDigest)
	snapshot.MetadataDigest = strings.TrimSpace(run.MetadataDigest)
	snapshot.PromptDigest = strings.TrimSpace(run.PromptDigest)
	snapshot.SkippedReasons = bodyTranslationSkippedReasonsFromCount(run.ExactExclusionCount)
	return snapshot
}

func bodyTranslationRunSnapshotDrifted(
	run repository.JobPhaseRun,
	snapshot bodyTranslationInputSnapshot,
) bool {
	persistedDigest := strings.TrimSpace(run.InputSnapshotDigest)
	if persistedDigest == "" {
		return false
	}
	return persistedDigest != strings.TrimSpace(snapshot.InputSnapshotDigest)
}

func bodyTranslationInputSnapshotDriftErrorSummary() *BodyTranslationPhaseErrorSummaryReadModel {
	return &BodyTranslationPhaseErrorSummaryReadModel{
		ErrorKind:  bodyTranslationInputSnapshotFailedKind,
		Reason:     bodyTranslationInputSnapshotDriftReason,
		Retryable:  false,
		IsRedacted: true,
	}
}

func bodyTranslationSkippedReasonsFromCount(count int) []string {
	if count <= 0 {
		return nil
	}
	reasons := make([]string, 0, count)
	for i := 0; i < count; i++ {
		reasons = append(reasons, bodyTranslationSkippedReasonExactDictionary)
	}
	return reasons
}

func (service *BodyTranslationPhaseService) nextExecutionOrder(loaded bodyTranslationLoadedContext) int {
	order := 0
	if loaded.personaRun != nil && loaded.personaRun.ExecutionOrder > order {
		order = loaded.personaRun.ExecutionOrder
	}
	if loaded.bodyRun != nil && loaded.bodyRun.ExecutionOrder > order {
		order = loaded.bodyRun.ExecutionOrder
	}
	return order + 1
}

func bodyTranslationExecutionSummaryFromRun(
	run *repository.JobPhaseRun,
) BodyTranslationPhaseExecutionSummaryReadModel {
	if run == nil {
		return BodyTranslationPhaseExecutionSummaryReadModel{}
	}
	return BodyTranslationPhaseExecutionSummaryReadModel{
		CredentialRef: strings.TrimSpace(run.CredentialRef),
		Provider:      strings.TrimSpace(run.AIProvider),
		Model:         strings.TrimSpace(run.ModelName),
		ExecutionMode: strings.TrimSpace(run.ExecutionMode),
	}
}

func toBodyTranslationInputSummaryReadModel(
	snapshot bodyTranslationInputSnapshot,
) BodyTranslationPhaseInputSummaryReadModel {
	return BodyTranslationPhaseInputSummaryReadModel{
		TargetCount:      snapshot.TargetCount,
		SkippedReasons:   append([]string(nil), snapshot.SkippedReasons...),
		InputSnapshotRef: nil,
		DictionaryDigest: snapshot.DictionaryDigest,
		PersonaDigest:    snapshot.PersonaDigest,
		MetadataDigest:   snapshot.MetadataDigest,
		PromptDigest:     snapshot.PromptDigest,
	}
}

func toBodyTranslationRequestSummaryReadModel(
	snapshot bodyTranslationInputSnapshot,
) BodyTranslationPhaseRequestSummaryReadModel {
	return BodyTranslationPhaseRequestSummaryReadModel{
		ProviderTargetCount:              snapshot.ProviderTargetCount,
		ExactDictionaryExclusionCount:    snapshot.ExactExclusionCount,
		PartialDictionaryConstraintCount: snapshot.PartialConstraintCount,
	}
}

func isBodyTranslationTerminalJob(state string) bool {
	switch strings.TrimSpace(state) {
	case bodyTranslationJobStateCompleted, bodyTranslationJobStateFailed, bodyTranslationJobStateCanceled:
		return true
	default:
		return false
	}
}

func isBodyTranslationActiveRunState(state string) bool {
	switch strings.TrimSpace(state) {
	case bodyTranslationPhaseStateRunning, bodyTranslationPhaseStatePaused, bodyTranslationPhaseStateRecoverableFail:
		return true
	default:
		return false
	}
}

func (service *BodyTranslationPhaseService) findPhaseRun(
	ctx context.Context,
	jobID int64,
	phaseType string,
) (*repository.JobPhaseRun, error) {
	run, err := service.jobLifecycleRepository.FindJobPhaseRun(ctx, jobID, phaseType)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, nil
		}
		return nil, fmt.Errorf("find body translation phase run %s: %w", phaseType, err)
	}
	return &run, nil
}

func cloneBodyTranslationInt64Pointer(value *int64) *int64 {
	if value == nil {
		return nil
	}
	cloned := *value
	return &cloned
}

func cloneBodyTranslationTimePointer(value *time.Time) *time.Time {
	if value == nil {
		return nil
	}
	cloned := *value
	return &cloned
}

func cloneBodyTranslationStringPointer(value *string) *string {
	if value == nil {
		return nil
	}
	cloned := *value
	return &cloned
}
