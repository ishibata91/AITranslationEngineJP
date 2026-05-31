package service

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"testing"
	"time"

	"aitranslationenginejp/internal/repository"
)

type fakeTermPhaseTransactor struct {
	jobRepo    *fakeTermPhaseJobLifecycleRepository
	foundation *fakeTermPhaseFoundationDataRepository
}

func (fake fakeTermPhaseTransactor) WithTransaction(ctx context.Context, fn func(context.Context) error) error {
	if fake.jobRepo == nil && fake.foundation == nil {
		return fn(ctx)
	}
	jobSnapshot := fake.jobRepo.clone()
	foundationSnapshot := fake.foundation.clone()
	err := fn(ctx)
	if err != nil {
		if fake.jobRepo != nil {
			fake.jobRepo.restore(jobSnapshot)
		}
		if fake.foundation != nil {
			fake.foundation.restore(foundationSnapshot)
		}
	}
	return err
}

type fakeTermPhaseProvider struct {
	translateFunc func(context.Context, TermTranslationProviderRequest) (TermTranslationProviderResult, error)
}

func (fake fakeTermPhaseProvider) TranslateTerm(ctx context.Context, request TermTranslationProviderRequest) (TermTranslationProviderResult, error) {
	if fake.translateFunc != nil {
		return fake.translateFunc(ctx, request)
	}
	return TermTranslationProviderResult{SourceTerm: request.SourceTerm, TranslatedTerm: request.SourceTerm + "_ja"}, nil
}

func (fake fakeTermPhaseProvider) TermTranslationProviderRequestsAreTestSafe() bool {
	return true
}

type fakeTermPhaseSecretStore struct {
	loadFunc func(context.Context, string) (string, error)
}

func (fake fakeTermPhaseSecretStore) Load(ctx context.Context, key string) (string, error) {
	if fake.loadFunc != nil {
		return fake.loadFunc(ctx, key)
	}
	return "", nil
}

type fakeTermPhaseJobLifecycleRepository struct {
	job                 repository.TranslationJob
	run                 *repository.JobPhaseRun
	phases              []repository.JobPhaseRun
	updateJobDrafts     []repository.TranslationJobUpdateDraft
	updatePhaseRunDraft []repository.JobPhaseRunUpdateDraft
	updateExpectedState []string
	expectedStateLog    *[]string
	phaseLinks          []repository.PhaseRunDictionaryEntry
	nextPhaseRunID      int64
	createLinkErr       error
	updatePhaseRunErr   error
}

func (fake *fakeTermPhaseJobLifecycleRepository) GetTranslationJobByID(context.Context, int64) (repository.TranslationJob, error) {
	return fake.job, nil
}
func (fake *fakeTermPhaseJobLifecycleRepository) UpdateTranslationJob(_ context.Context, _ int64, draft repository.TranslationJobUpdateDraft) (repository.TranslationJob, error) {
	fake.updateJobDrafts = append(fake.updateJobDrafts, draft)
	fake.job.State = draft.State
	fake.job.ProgressPercent = draft.ProgressPercent
	fake.job.StartedAt = draft.StartedAt
	fake.job.FinishedAt = draft.FinishedAt
	return fake.job, nil
}
func (fake *fakeTermPhaseJobLifecycleRepository) CreateJobPhaseRun(_ context.Context, draft repository.JobPhaseRunDraft) (repository.JobPhaseRun, error) {
	if fake.run != nil {
		return repository.JobPhaseRun{}, repository.ErrConflict
	}
	fake.nextPhaseRunID++
	run := repository.JobPhaseRun{
		ID:               fake.nextPhaseRunID,
		TranslationJobID: fake.job.ID,
		PhaseType:        draft.PhaseType,
		State:            draft.State,
		ExecutionOrder:   draft.ExecutionOrder,
		AIProvider:       draft.AIProvider,
		ModelName:        draft.ModelName,
		ExecutionMode:    draft.ExecutionMode,
		CredentialRef:    draft.CredentialRef,
		InstructionKind:  draft.InstructionKind,
	}
	fake.run = &run
	return run, nil
}
func (fake *fakeTermPhaseJobLifecycleRepository) FindJobPhaseRun(context.Context, int64, string) (repository.JobPhaseRun, error) {
	if fake.run == nil {
		return repository.JobPhaseRun{}, repository.ErrNotFound
	}
	return *fake.run, nil
}
func (fake *fakeTermPhaseJobLifecycleRepository) UpdateJobPhaseRun(_ context.Context, _ int64, draft repository.JobPhaseRunUpdateDraft) (repository.JobPhaseRun, error) {
	if fake.run == nil {
		return repository.JobPhaseRun{}, repository.ErrNotFound
	}
	if fake.updatePhaseRunErr != nil {
		return repository.JobPhaseRun{}, fake.updatePhaseRunErr
	}
	fake.updatePhaseRunDraft = append(fake.updatePhaseRunDraft, draft)
	fake.run.State = draft.State
	fake.run.ProgressPercent = draft.ProgressPercent
	fake.run.LatestError = draft.LatestError
	fake.run.AIProvider = draft.AIProvider
	fake.run.ModelName = draft.ModelName
	fake.run.ExecutionMode = draft.ExecutionMode
	fake.run.CredentialRef = draft.CredentialRef
	fake.run.InstructionKind = draft.InstructionKind
	fake.run.StartedAt = draft.StartedAt
	fake.run.FinishedAt = draft.FinishedAt
	return *fake.run, nil
}
func (fake *fakeTermPhaseJobLifecycleRepository) UpdateJobPhaseRunWhenState(
	ctx context.Context,
	id int64,
	expectedState string,
	draft repository.JobPhaseRunUpdateDraft,
) (repository.JobPhaseRun, error) {
	fake.updateExpectedState = append(fake.updateExpectedState, expectedState)
	if fake.expectedStateLog != nil {
		*fake.expectedStateLog = append(*fake.expectedStateLog, expectedState)
	}
	return fake.UpdateJobPhaseRun(ctx, id, draft)
}
func (fake *fakeTermPhaseJobLifecycleRepository) ListJobPhaseRunsByJobID(context.Context, int64) ([]repository.JobPhaseRun, error) {
	if len(fake.phases) != 0 {
		return fake.phases, nil
	}
	if fake.run != nil {
		return []repository.JobPhaseRun{*fake.run}, nil
	}
	return nil, nil
}
func (fake *fakeTermPhaseJobLifecycleRepository) CreatePhaseRunDictionaryEntry(_ context.Context, draft repository.PhaseRunDictionaryEntryDraft) (repository.PhaseRunDictionaryEntry, error) {
	if fake.createLinkErr != nil {
		return repository.PhaseRunDictionaryEntry{}, fake.createLinkErr
	}
	for _, link := range fake.phaseLinks {
		if link.PhaseRunID == draft.PhaseRunID && link.DictionaryEntryID == draft.DictionaryEntryID && link.Role == draft.Role {
			return repository.PhaseRunDictionaryEntry{}, repository.ErrConflict
		}
	}
	link := repository.PhaseRunDictionaryEntry{ID: int64(len(fake.phaseLinks) + 1), PhaseRunID: draft.PhaseRunID, DictionaryEntryID: draft.DictionaryEntryID, Role: draft.Role}
	fake.phaseLinks = append(fake.phaseLinks, link)
	return link, nil
}

func (fake *fakeTermPhaseJobLifecycleRepository) clone() *fakeTermPhaseJobLifecycleRepository {
	if fake == nil {
		return nil
	}
	cloned := *fake
	cloned.phases = append([]repository.JobPhaseRun(nil), fake.phases...)
	cloned.updateJobDrafts = append([]repository.TranslationJobUpdateDraft(nil), fake.updateJobDrafts...)
	cloned.updatePhaseRunDraft = append([]repository.JobPhaseRunUpdateDraft(nil), fake.updatePhaseRunDraft...)
	cloned.updateExpectedState = append([]string(nil), fake.updateExpectedState...)
	cloned.phaseLinks = append([]repository.PhaseRunDictionaryEntry(nil), fake.phaseLinks...)
	if fake.run != nil {
		runCopy := *fake.run
		cloned.run = &runCopy
	}
	return &cloned
}

func (fake *fakeTermPhaseJobLifecycleRepository) restore(snapshot *fakeTermPhaseJobLifecycleRepository) {
	if fake == nil || snapshot == nil {
		return
	}
	*fake = *snapshot
	if snapshot.run != nil {
		runCopy := *snapshot.run
		fake.run = &runCopy
	}
}
func (fake *fakeTermPhaseJobLifecycleRepository) ListPhaseRunDictionaryEntriesByPhaseRunID(_ context.Context, phaseRunID int64) ([]repository.PhaseRunDictionaryEntry, error) {
	result := make([]repository.PhaseRunDictionaryEntry, 0)
	for _, link := range fake.phaseLinks {
		if link.PhaseRunID == phaseRunID {
			result = append(result, link)
		}
	}
	return result, nil
}

type fakeTermPhaseFoundationDataRepository struct {
	entries           map[string]repository.DictionaryEntry
	nextID            int64
	createErrBySource map[string]error
}

func (fake *fakeTermPhaseFoundationDataRepository) key(jobID *int64, lifecycle string, scope string, sourceTerm string) string {
	jobPart := "nil"
	if jobID != nil {
		jobPart = fmt.Sprintf("%d", *jobID)
	}
	return strings.ToLower(strings.TrimSpace(jobPart + "|" + lifecycle + "|" + scope + "|" + sourceTerm))
}
func (fake *fakeTermPhaseFoundationDataRepository) idKey(id int64) string {
	return fmt.Sprintf("id:%d", id)
}
func (fake *fakeTermPhaseFoundationDataRepository) GetDictionaryEntryByID(_ context.Context, id int64) (repository.DictionaryEntry, error) {
	entry, ok := fake.entries[fake.idKey(id)]
	if !ok {
		return repository.DictionaryEntry{}, repository.ErrNotFound
	}
	return entry, nil
}
func (fake *fakeTermPhaseFoundationDataRepository) CreateDictionaryEntry(_ context.Context, draft repository.DictionaryEntryDraft) (repository.DictionaryEntry, error) {
	if err := fake.createErrBySource[strings.TrimSpace(draft.SourceTerm)]; err != nil {
		return repository.DictionaryEntry{}, err
	}
	key := fake.key(draft.TranslationJobID, draft.DictionaryLifecycle, draft.DictionaryScope, draft.SourceTerm)
	if _, ok := fake.entries[key]; ok {
		return repository.DictionaryEntry{}, repository.ErrConflict
	}
	fake.nextID++
	entry := repository.DictionaryEntry{
		ID:                  fake.nextID,
		TranslationJobID:    draft.TranslationJobID,
		DictionaryLifecycle: draft.DictionaryLifecycle,
		DictionaryScope:     draft.DictionaryScope,
		DictionarySource:    draft.DictionarySource,
		SourceTerm:          draft.SourceTerm,
		TranslatedTerm:      draft.TranslatedTerm,
		TermKind:            draft.TermKind,
		Reusable:            draft.Reusable,
		UpdatedAt:           time.Date(2026, 5, 1, 0, 0, int(fake.nextID), 0, time.UTC),
	}
	fake.entries[key] = entry
	fake.entries[fake.idKey(entry.ID)] = entry
	return entry, nil
}

func (fake *fakeTermPhaseFoundationDataRepository) clone() *fakeTermPhaseFoundationDataRepository {
	if fake == nil {
		return nil
	}
	cloned := &fakeTermPhaseFoundationDataRepository{
		entries:           make(map[string]repository.DictionaryEntry, len(fake.entries)),
		nextID:            fake.nextID,
		createErrBySource: make(map[string]error, len(fake.createErrBySource)),
	}
	for key, entry := range fake.entries {
		cloned.entries[key] = entry
	}
	for key, err := range fake.createErrBySource {
		cloned.createErrBySource[key] = err
	}
	return cloned
}

func (fake *fakeTermPhaseFoundationDataRepository) restore(snapshot *fakeTermPhaseFoundationDataRepository) {
	if fake == nil || snapshot == nil {
		return
	}
	fake.entries = make(map[string]repository.DictionaryEntry, len(snapshot.entries))
	for key, entry := range snapshot.entries {
		fake.entries[key] = entry
	}
	fake.nextID = snapshot.nextID
	fake.createErrBySource = make(map[string]error, len(snapshot.createErrBySource))
	for key, err := range snapshot.createErrBySource {
		fake.createErrBySource[key] = err
	}
}
func (fake *fakeTermPhaseFoundationDataRepository) UpdateDictionaryEntry(context.Context, int64, repository.DictionaryEntryUpdateDraft) (repository.DictionaryEntry, error) {
	return repository.DictionaryEntry{}, errors.New("not implemented")
}
func (fake *fakeTermPhaseFoundationDataRepository) ListDictionaryEntries(_ context.Context, translationJobID *int64, lifecycle string, scope string, sourceTerm string) ([]repository.DictionaryEntry, error) {
	result := make([]repository.DictionaryEntry, 0)
	for key, entry := range fake.entries {
		if !shouldIncludeDictionaryEntry(key, entry, translationJobID, lifecycle, scope, sourceTerm) {
			continue
		}
		result = append(result, entry)
	}
	return result, nil
}

func shouldIncludeDictionaryEntry(
	key string,
	entry repository.DictionaryEntry,
	translationJobID *int64,
	lifecycle string,
	scope string,
	sourceTerm string,
) bool {
	if strings.HasPrefix(key, "id:") {
		return false
	}
	if !matchesTranslationJobID(entry, translationJobID) {
		return false
	}
	if lifecycle != "" && entry.DictionaryLifecycle != lifecycle {
		return false
	}
	if scope != "" && entry.DictionaryScope != scope {
		return false
	}
	if sourceTerm != "" && entry.SourceTerm != sourceTerm {
		return false
	}
	return true
}

func matchesTranslationJobID(entry repository.DictionaryEntry, translationJobID *int64) bool {
	if translationJobID == nil {
		return entry.TranslationJobID == nil
	}
	if entry.TranslationJobID == nil {
		return false
	}
	return *entry.TranslationJobID == *translationJobID
}
func (fake *fakeTermPhaseFoundationDataRepository) FindDictionaryEntry(_ context.Context, translationJobID *int64, lifecycle string, scope string, sourceTerm string) (repository.DictionaryEntry, error) {
	entry, ok := fake.entries[fake.key(translationJobID, lifecycle, scope, sourceTerm)]
	if !ok {
		return repository.DictionaryEntry{}, repository.ErrNotFound
	}
	return entry, nil
}

type fakeTermPhaseTranslationSourceRepository struct {
	records []repository.TranslationRecord
	fields  map[int64][]repository.TranslationField
}

func (fake fakeTermPhaseTranslationSourceRepository) GetXEditExtractedDataByID(context.Context, int64) (repository.XEditExtractedData, error) {
	return repository.XEditExtractedData{ID: 1}, nil
}
func (fake fakeTermPhaseTranslationSourceRepository) ListTranslationRecordsByXEditID(context.Context, int64) ([]repository.TranslationRecord, error) {
	return fake.records, nil
}
func (fake fakeTermPhaseTranslationSourceRepository) ListTranslationFieldsByTranslationRecordID(_ context.Context, recordID int64) ([]repository.TranslationField, error) {
	return fake.fields[recordID], nil
}

func newTermTranslationPhaseServiceForTest(provider TermTranslationProvider) (*TermTranslationPhaseService, *fakeTermPhaseJobLifecycleRepository, *fakeTermPhaseFoundationDataRepository) {
	jobRepo := &fakeTermPhaseJobLifecycleRepository{
		job: repository.TranslationJob{ID: 1, XEditExtractedDataID: 1, JobName: "job-1", State: termTranslationJobStateReady},
		run: &repository.JobPhaseRun{
			ID:               200,
			TranslationJobID: 1,
			PhaseType:        termTranslationPhaseType,
			State:            termTranslationPhaseStatePending,
		},
		phases: []repository.JobPhaseRun{{
			ID:               100,
			TranslationJobID: 1,
			PhaseType:        termTranslationInitialPhaseType,
			State:            termTranslationPhaseStateCompleted,
			ExecutionMode:    "sync",
			AIProvider:       TermTranslationProviderGemini,
			ModelName:        "gemini-2.5-pro",
			CredentialRef:    "gemini-primary",
		}},
		nextPhaseRunID: 200,
	}
	foundation := &fakeTermPhaseFoundationDataRepository{
		entries:           map[string]repository.DictionaryEntry{},
		createErrBySource: map[string]error{},
	}
	source := fakeTermPhaseTranslationSourceRepository{
		records: []repository.TranslationRecord{{ID: 1, RecordType: "NPC_", XEditExtractedDataID: 1}},
		fields: map[int64][]repository.TranslationField{
			1: {
				{ID: 11, TranslationRecordID: 1, SubrecordType: "FULL", SourceText: "Iron Sword"},
				{ID: 12, TranslationRecordID: 1, SubrecordType: "FULL", SourceText: "Steel Sword"},
			},
		},
	}
	secretStore := repository.NewInMemorySecretStore()
	if err := secretStore.Save(context.Background(), termTranslationProviderSecretKey(TermTranslationProviderGemini), "gemini-secret"); err != nil {
		panic(err)
	}
	transactor := fakeTermPhaseTransactor{jobRepo: jobRepo, foundation: foundation}
	service := NewTermTranslationPhaseService(jobRepo, foundation, source, transactor, secretStore, provider)
	service.WithTermTranslationProviderSettings(fakePhaseProviderSettingsConsumer{
		resolveFunc: func(_ context.Context, input ProviderSettingsResolveInput) (ProviderSettingsResolveResult, error) {
			endpoint := "https://provider-settings.example.test"
			credentialRef := "gemini-primary"
			if input.Selection.ProviderID == TermTranslationProviderLMStudio {
				return ProviderSettingsResolveResult{
					ConsumerID:      input.ConsumerID,
					ProviderID:      input.Selection.ProviderID,
					Model:           input.Selection.Model,
					ExecutionMethod: input.Selection.ExecutionMethod,
					UseBatchAPI:     input.Selection.UseBatchAPI,
					Endpoint:        &endpoint,
					CredentialState: providerSettingsCredentialStateNotRequired,
				}, nil
			}
			return ProviderSettingsResolveResult{
				ConsumerID:            input.ConsumerID,
				ProviderID:            input.Selection.ProviderID,
				Model:                 input.Selection.Model,
				ExecutionMethod:       input.Selection.ExecutionMethod,
				UseBatchAPI:           input.Selection.UseBatchAPI,
				Endpoint:              &endpoint,
				CredentialReferenceID: &credentialRef,
				CredentialState:       providerSettingsCredentialStateConfigured,
			}, nil
		},
	})
	service.now = func() time.Time { return time.Date(2026, 5, 1, 12, 0, 0, 0, time.UTC) }
	return service, jobRepo, foundation
}

func TestTermTranslationPhaseServiceReadSummaryCountsSnapshotHitsAndConfirmedTerms(t *testing.T) {
	service, jobRepo, foundation := newTermTranslationPhaseServiceForTest(nil)
	run := repository.JobPhaseRun{ID: 210, TranslationJobID: 1, PhaseType: termTranslationPhaseType, State: termTranslationPhaseStateCompleted, ProgressPercent: 100}
	jobRepo.run = &run
	jobRepo.phases = append(jobRepo.phases, run)

	snapshot, err := foundation.CreateDictionaryEntry(context.Background(), repository.DictionaryEntryDraft{
		DictionaryLifecycle: termTranslationDictionaryLifecycleMaster,
		DictionaryScope:     "NPC_:FULL",
		DictionarySource:    "seed",
		SourceTerm:          "Iron Sword",
		TranslatedTerm:      "鉄の剣",
		Reusable:            true,
	})
	if err != nil {
		t.Fatalf("prepare snapshot entry: %v", err)
	}
	_, err = foundation.CreateDictionaryEntry(context.Background(), repository.DictionaryEntryDraft{
		TranslationJobID:    int64Pointer(1),
		DictionaryLifecycle: termTranslationDictionaryLifecycleJob,
		DictionaryScope:     "NPC_:FULL",
		DictionarySource:    termTranslationDictionarySourceSnapshot,
		SourceTerm:          "Iron Sword",
		TranslatedTerm:      "鉄の剣",
		Reusable:            true,
	})
	if err != nil {
		t.Fatalf("prepare job snapshot entry: %v", err)
	}
	_, err = foundation.CreateDictionaryEntry(context.Background(), repository.DictionaryEntryDraft{
		TranslationJobID:    int64Pointer(1),
		DictionaryLifecycle: termTranslationDictionaryLifecycleJob,
		DictionaryScope:     "NPC_:FULL",
		DictionarySource:    termTranslationDictionarySourceProvider,
		SourceTerm:          "Steel Sword",
		TranslatedTerm:      "鋼の剣",
		Reusable:            true,
	})
	if err != nil {
		t.Fatalf("prepare job provider entry: %v", err)
	}
	jobRepo.phaseLinks = append(jobRepo.phaseLinks, repository.PhaseRunDictionaryEntry{PhaseRunID: run.ID, DictionaryEntryID: snapshot.ID, Role: termTranslationLinkRoleSnapshot})

	summary, err := service.ReadSummary(context.Background(), 1)
	if err != nil {
		t.Fatalf("expected summary success, got %v", err)
	}
	if summary.TotalTermCount != 2 || summary.DictionaryHitCount != 1 || summary.AITargetCount != 1 {
		t.Fatalf("unexpected summary counts: %#v", summary)
	}
	if summary.ResultSummary == nil || summary.ResultSummary.ConfirmedCount != 2 {
		t.Fatalf("expected confirmed count=2, got %#v", summary.ResultSummary)
	}
	if summary.Execution.SnapshotDigest == nil || summary.Execution.SnapshotVersion == nil {
		t.Fatalf("expected snapshot metadata, got %#v", summary.Execution)
	}
}

func TestTermTranslationPhaseServiceReadSummaryKeepsReadyJobWithoutPhaseRunsIdle(t *testing.T) {
	service, jobRepo, _ := newTermTranslationPhaseServiceForTest(nil)
	jobRepo.run = nil
	jobRepo.phases = nil

	summary, err := service.ReadSummary(context.Background(), 1)
	if err != nil {
		t.Fatalf("expected summary success for ready job without phase runs: %v", err)
	}
	if summary.PhaseRunID != nil {
		t.Fatalf("expected no phase run id for ready job without phase runs, got %#v", summary.PhaseRunID)
	}
	if summary.PhaseState != termTranslationPhaseStateIdleReady {
		t.Fatalf("expected idle ready phase state, got %#v", summary)
	}
}

func TestTermTranslationPhaseServiceReadSummaryReturnsNotFoundForNonReadyJobWithoutPhaseRuns(t *testing.T) {
	service, jobRepo, _ := newTermTranslationPhaseServiceForTest(nil)
	jobRepo.job.State = termTranslationJobStateRunning
	jobRepo.run = nil
	jobRepo.phases = nil

	_, err := service.ReadSummary(context.Background(), 1)
	if err == nil {
		t.Fatal("expected service error for non-ready job without phase runs")
	}
	if !strings.Contains(err.Error(), "load initial execution phase: not found") {
		t.Fatalf("expected load initial execution phase not found error, got %v", err)
	}
}

func TestTermTranslationPhaseServiceReadSummaryUsesExistingRunExecutionConfig(t *testing.T) {
	service, jobRepo, _ := newTermTranslationPhaseServiceForTest(nil)
	run := repository.JobPhaseRun{
		ID:               220,
		TranslationJobID: 1,
		PhaseType:        termTranslationPhaseType,
		State:            termTranslationPhaseStatePaused,
		AIProvider:       TermTranslationProviderXAI,
		ModelName:        "grok-2-latest",
		ExecutionMode:    "batch",
		CredentialRef:    "xai-primary",
	}
	jobRepo.run = &run
	jobRepo.phases = append(jobRepo.phases, run)

	summary, err := service.ReadSummary(context.Background(), 1)
	if err != nil {
		t.Fatalf("expected summary success with existing run, got %v", err)
	}
	if summary.Execution.Provider != TermTranslationProviderXAI {
		t.Fatalf("expected execution provider from existing run, got %#v", summary.Execution)
	}
	if summary.Execution.Model != "grok-2-latest" {
		t.Fatalf("expected execution model from existing run, got %#v", summary.Execution)
	}
	if summary.Execution.ExecutionMode != "batch" {
		t.Fatalf("expected execution mode from existing run, got %#v", summary.Execution)
	}
	if summary.Execution.CredentialRef != "xai-primary" {
		t.Fatalf("expected execution credential ref from existing run, got %#v", summary.Execution)
	}
}

func TestTermTranslationPhaseServiceReadSummaryRejectsResumeAndAllowsRetryForRecoverableFailedRun(t *testing.T) {
	service, jobRepo, _ := newTermTranslationPhaseServiceForTest(nil)
	run := repository.JobPhaseRun{
		ID:               360,
		TranslationJobID: 1,
		PhaseType:        termTranslationPhaseType,
		State:            termTranslationPhaseStateRecoverableFail,
		ProgressPercent:  0,
		LatestError:      termTranslationErrorKindInvalidProvider,
	}
	jobRepo.job.State = termTranslationJobStateRunning
	jobRepo.run = &run
	jobRepo.phases = append(jobRepo.phases, run)

	summary, err := service.ReadSummary(context.Background(), 1)
	if err != nil {
		t.Fatalf("expected summary success with recoverable failed run, got %v", err)
	}
	if summary.ActionEnablement.CanResume {
		t.Fatalf("expected recoverable failed run to reject resume, got %#v", summary.ActionEnablement)
	}
	if !summary.ActionEnablement.CanRetry {
		t.Fatalf("expected recoverable failed run to be retryable, got %#v", summary.ActionEnablement)
	}
}

func TestTermTranslationPhaseServiceReadSummaryDisablesActionsForTerminalJob(t *testing.T) {
	cases := []struct {
		name  string
		state string
	}{
		{name: "running", state: termTranslationPhaseStateRunning},
		{name: "paused", state: termTranslationPhaseStatePaused},
		{name: "recoverable_failed", state: termTranslationPhaseStateRecoverableFail},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			service, jobRepo, _ := newTermTranslationPhaseServiceForTest(nil)
			finishedAt := time.Date(2026, 5, 1, 13, 0, 0, 0, time.UTC)
			run := repository.JobPhaseRun{
				ID:               370,
				TranslationJobID: 1,
				PhaseType:        termTranslationPhaseType,
				State:            tc.state,
				ProgressPercent:  50,
			}
			jobRepo.job.State = termTranslationJobStateCompleted
			jobRepo.job.FinishedAt = &finishedAt
			jobRepo.run = &run
			jobRepo.phases = append(jobRepo.phases, run)

			summary, err := service.ReadSummary(context.Background(), 1)
			if err != nil {
				t.Fatalf("expected terminal job summary success, got %v", err)
			}
			if summary.ActionEnablement.CanPause ||
				summary.ActionEnablement.CanResume ||
				summary.ActionEnablement.CanRetry {
				t.Fatalf("expected terminal job actions disabled, got %#v", summary.ActionEnablement)
			}
		})
	}
}

func TestTermTranslationPhaseServiceReadSummaryReturnsJapaneseBlockedReasons(t *testing.T) {
	service, jobRepo, _ := newTermTranslationPhaseServiceForTest(nil)
	run := repository.JobPhaseRun{
		ID:               361,
		TranslationJobID: 1,
		PhaseType:        termTranslationPhaseType,
		State:            termTranslationPhaseStateCompleted,
		ProgressPercent:  100,
	}
	jobRepo.job.State = termTranslationJobStateRunning
	jobRepo.run = &run
	jobRepo.phases = append(jobRepo.phases, run)

	summary, err := service.ReadSummary(context.Background(), 1)
	if err != nil {
		t.Fatalf("expected summary success with completed run, got %v", err)
	}
	if summary.ActionEnablement.PauseBlockedReason == nil ||
		*summary.ActionEnablement.PauseBlockedReason != termTranslationReasonPhaseNotRunning {
		t.Fatalf("expected Japanese pause blocked reason, got %#v", summary.ActionEnablement.PauseBlockedReason)
	}
	if summary.ActionEnablement.ResumeBlockedReason == nil ||
		*summary.ActionEnablement.ResumeBlockedReason != termTranslationReasonPhaseNotResumable {
		t.Fatalf("expected Japanese resume blocked reason, got %#v", summary.ActionEnablement.ResumeBlockedReason)
	}
	if summary.ActionEnablement.RetryBlockedReason == nil ||
		*summary.ActionEnablement.RetryBlockedReason != termTranslationReasonPhaseNotRetryable {
		t.Fatalf("expected Japanese retry blocked reason, got %#v", summary.ActionEnablement.RetryBlockedReason)
	}
}

func TestTermTranslationPhaseServiceStartPhaseReturnsProviderFailureWhenProviderIsNil(t *testing.T) {
	service, _, _ := newTermTranslationPhaseServiceForTest(nil)

	result, err := service.StartPhase(context.Background(), 1)
	if err != nil {
		t.Fatalf("expected command response without transport error, got %v", err)
	}
	if result.ErrorSummary == nil || result.ErrorSummary.ErrorKind != termTranslationErrorKindProviderFailure {
		t.Fatalf("expected provider failure error summary, got %#v", result)
	}
	if result.PhaseState != termTranslationPhaseStateRecoverableFail {
		t.Fatalf("expected recoverable failed phase state, got %#v", result)
	}
}

func TestTermTranslationPhaseServiceStartPhaseReturnsInvalidProviderResponseWhenTranslationIsInvalid(t *testing.T) {
	service, _, _ := newTermTranslationPhaseServiceForTest(fakeTermPhaseProvider{
		translateFunc: func(context.Context, TermTranslationProviderRequest) (TermTranslationProviderResult, error) {
			return TermTranslationProviderResult{SourceTerm: "different-source", TranslatedTerm: "訳語"}, nil
		},
	})

	result, err := service.StartPhase(context.Background(), 1)
	if err != nil {
		t.Fatalf("expected command response without transport error, got %v", err)
	}
	if result.ErrorSummary == nil || result.ErrorSummary.ErrorKind != termTranslationErrorKindInvalidProvider {
		t.Fatalf("expected invalid provider response error summary, got %#v", result)
	}
}

func TestTermTranslationPhaseServicePausePhaseRejectsNonRunningPhase(t *testing.T) {
	service, jobRepo, _ := newTermTranslationPhaseServiceForTest(nil)
	run := repository.JobPhaseRun{ID: 300, TranslationJobID: 1, PhaseType: termTranslationPhaseType, State: termTranslationPhaseStatePaused}
	jobRepo.run = &run
	jobRepo.phases = append(jobRepo.phases, run)

	result, err := service.PausePhase(context.Background(), 1, 300)
	if err != nil {
		t.Fatalf("expected rejected command without transport error, got %v", err)
	}
	if result.ErrorSummary == nil || result.ErrorSummary.ErrorKind != "term_phase_incomplete" {
		t.Fatalf("expected term_phase_incomplete rejection, got %#v", result)
	}
	if len(jobRepo.updateJobDrafts) != 0 || len(jobRepo.updatePhaseRunDraft) != 0 {
		t.Fatalf("expected no persistence mutation on rejection, got job=%d phase=%d", len(jobRepo.updateJobDrafts), len(jobRepo.updatePhaseRunDraft))
	}
}

func TestTermTranslationPhaseServicePausePhaseKeepsJobRunning(t *testing.T) {
	service, jobRepo, _ := newTermTranslationPhaseServiceForTest(nil)
	now := service.now()
	run := repository.JobPhaseRun{
		ID:               310,
		TranslationJobID: 1,
		PhaseType:        termTranslationPhaseType,
		State:            termTranslationPhaseStateRunning,
		ProgressPercent:  42,
		StartedAt:        &now,
	}
	jobRepo.job.State = termTranslationJobStateRunning
	jobRepo.job.StartedAt = &now
	jobRepo.run = &run
	jobRepo.phases = append(jobRepo.phases, run)

	result, err := service.PausePhase(context.Background(), 1, run.ID)
	if err != nil {
		t.Fatalf("expected pause success, got %v", err)
	}
	if result.PhaseState != termTranslationPhaseStatePaused {
		t.Fatalf("expected paused phase state, got %#v", result)
	}
	assertTermTranslationJobStateRunning(t, jobRepo)
}

func TestTermTranslationPhaseServiceRetryPhaseCompletesWhenSnapshotAlreadyConfirmsAllTerms(t *testing.T) {
	service, jobRepo, foundation := newTermTranslationPhaseServiceForTest(nil)
	run := repository.JobPhaseRun{ID: 320, TranslationJobID: 1, PhaseType: termTranslationPhaseType, State: termTranslationPhaseStateRecoverableFail, ProgressPercent: 10}
	jobRepo.run = &run
	jobRepo.phases = append(jobRepo.phases, run)

	masterEntry, err := foundation.CreateDictionaryEntry(context.Background(), repository.DictionaryEntryDraft{
		DictionaryLifecycle: termTranslationDictionaryLifecycleMaster,
		DictionaryScope:     "NPC_:FULL",
		DictionarySource:    "seed",
		SourceTerm:          "Iron Sword",
		TranslatedTerm:      "鉄の剣",
		Reusable:            true,
	})
	if err != nil {
		t.Fatalf("prepare snapshot master entry: %v", err)
	}
	secondMasterEntry, err := foundation.CreateDictionaryEntry(context.Background(), repository.DictionaryEntryDraft{
		DictionaryLifecycle: termTranslationDictionaryLifecycleMaster,
		DictionaryScope:     "NPC_:FULL",
		DictionarySource:    "seed",
		SourceTerm:          "Steel Sword",
		TranslatedTerm:      "鋼の剣",
		Reusable:            true,
	})
	if err != nil {
		t.Fatalf("prepare snapshot master entry 2: %v", err)
	}
	jobRepo.phaseLinks = append(jobRepo.phaseLinks, repository.PhaseRunDictionaryEntry{PhaseRunID: run.ID, DictionaryEntryID: masterEntry.ID, Role: termTranslationLinkRoleSnapshot})
	jobRepo.phaseLinks = append(jobRepo.phaseLinks, repository.PhaseRunDictionaryEntry{PhaseRunID: run.ID, DictionaryEntryID: secondMasterEntry.ID, Role: termTranslationLinkRoleSnapshot})

	result, err := service.RetryPhase(context.Background(), 1, 320)
	if err != nil {
		t.Fatalf("expected retry success, got %v", err)
	}
	if result.PhaseState != termTranslationPhaseStateCompleted {
		t.Fatalf("expected completed retry result, got %#v", result)
	}
	if result.ErrorSummary != nil {
		t.Fatalf("expected nil error summary, got %#v", result.ErrorSummary)
	}
}

func TestTermTranslationPhaseServiceRetryPhaseUsesPersistedSnapshotInsteadOfCurrentMasterDictionary(t *testing.T) {
	capturedSources := make([]string, 0)
	service, jobRepo, foundation := newTermTranslationPhaseServiceForTest(fakeTermPhaseProvider{
		translateFunc: func(_ context.Context, request TermTranslationProviderRequest) (TermTranslationProviderResult, error) {
			capturedSources = append(capturedSources, request.SourceTerm)
			return TermTranslationProviderResult{SourceTerm: request.SourceTerm, TranslatedTerm: request.SourceTerm + "_ja"}, nil
		},
	})
	run := repository.JobPhaseRun{ID: 330, TranslationJobID: 1, PhaseType: termTranslationPhaseType, State: termTranslationPhaseStateRecoverableFail}
	jobRepo.run = &run
	jobRepo.phases = append(jobRepo.phases, run)

	snapshot, err := foundation.CreateDictionaryEntry(context.Background(), repository.DictionaryEntryDraft{
		DictionaryLifecycle: termTranslationDictionaryLifecycleMaster,
		DictionaryScope:     "NPC_:FULL",
		DictionarySource:    "seed",
		SourceTerm:          "Iron Sword",
		TranslatedTerm:      "鉄の剣",
		Reusable:            true,
	})
	if err != nil {
		t.Fatalf("prepare persisted snapshot entry: %v", err)
	}
	jobRepo.phaseLinks = append(jobRepo.phaseLinks, repository.PhaseRunDictionaryEntry{PhaseRunID: run.ID, DictionaryEntryID: snapshot.ID, Role: termTranslationLinkRoleSnapshot})

	_, err = foundation.CreateDictionaryEntry(context.Background(), repository.DictionaryEntryDraft{
		DictionaryLifecycle: termTranslationDictionaryLifecycleMaster,
		DictionaryScope:     "NPC_:FULL",
		DictionarySource:    "seed",
		SourceTerm:          "Steel Sword",
		TranslatedTerm:      "鋼の剣",
		Reusable:            true,
	})
	if err != nil {
		t.Fatalf("prepare current master entry: %v", err)
	}

	result, err := service.RetryPhase(context.Background(), 1, run.ID)
	if err != nil {
		t.Fatalf("expected retry success, got %v", err)
	}
	if result.PhaseState != termTranslationPhaseStateCompleted {
		t.Fatalf("expected completed retry, got %#v", result)
	}
	if len(capturedSources) != 1 || capturedSources[0] != "Steel Sword" {
		t.Fatalf("expected retry to translate only non-snapshot term, got %#v", capturedSources)
	}
	if len(jobRepo.updateExpectedState) == 0 || jobRepo.updateExpectedState[0] != termTranslationPhaseStateRecoverableFail {
		t.Fatalf("expected retry to update from recoverable_failed state, got %#v", jobRepo.updateExpectedState)
	}
}

func TestTermTranslationPhaseServiceStartPhaseRetainsInitialSnapshotLinksAcrossProviderFailure(t *testing.T) {
	capturedSources := make([]string, 0)
	service, jobRepo, foundation := newTermTranslationPhaseServiceForTest(
		termTranslationProviderFailsOnceForSteelSword(&capturedSources),
	)

	snapshot := createTermTranslationSnapshotEntry(t, foundation, "Iron Sword", "鉄の剣")

	result, err := service.StartPhase(context.Background(), 1)
	if err != nil {
		t.Fatalf("expected provider failure command result, got %v", err)
	}
	assertTermTranslationFailedStartPreservesSnapshot(t, foundation, jobRepo, result, snapshot.ID)
	assertTermTranslationJobStateRunning(t, jobRepo)
	addTermTranslationMasterEntryAfterFailedStart(t, foundation, "Steel Sword", "鋼の剣")

	capturedSources = capturedSources[:0]
	result, err = service.RetryPhase(context.Background(), 1, *result.PhaseRunID)
	if err != nil {
		t.Fatalf("expected retry success, got %v", err)
	}
	assertTermTranslationRetryUsesRemainingAITarget(t, result, capturedSources, "Steel Sword")
}

func termTranslationProviderFailsOnceForSteelSword(capturedSources *[]string) fakeTermPhaseProvider {
	failuresBySource := map[string]int{}
	return fakeTermPhaseProvider{
		translateFunc: func(_ context.Context, request TermTranslationProviderRequest) (TermTranslationProviderResult, error) {
			*capturedSources = append(*capturedSources, request.SourceTerm)
			if request.SourceTerm == "Steel Sword" && failuresBySource[request.SourceTerm] == 0 {
				failuresBySource[request.SourceTerm]++
				return TermTranslationProviderResult{}, NewTermTranslationProviderError(
					TermTranslationProviderErrorKindProviderFailure,
					"provider request failed",
					true,
					errors.New("upstream timeout"),
				)
			}
			return TermTranslationProviderResult{
				SourceTerm:     request.SourceTerm,
				TranslatedTerm: request.SourceTerm + "_ja",
			}, nil
		},
	}
}

func createTermTranslationSnapshotEntry(
	t *testing.T,
	foundation *fakeTermPhaseFoundationDataRepository,
	sourceTerm string,
	translatedTerm string,
) repository.DictionaryEntry {
	t.Helper()
	snapshot, err := foundation.CreateDictionaryEntry(context.Background(), repository.DictionaryEntryDraft{
		DictionaryLifecycle: termTranslationDictionaryLifecycleMaster,
		DictionaryScope:     "NPC_:FULL",
		DictionarySource:    "seed",
		SourceTerm:          sourceTerm,
		TranslatedTerm:      translatedTerm,
		Reusable:            true,
	})
	if err != nil {
		t.Fatalf("prepare persisted snapshot entry: %v", err)
	}
	return snapshot
}

func assertTermTranslationFailedStartPreservesSnapshot(
	t *testing.T,
	foundation *fakeTermPhaseFoundationDataRepository,
	jobRepo *fakeTermPhaseJobLifecycleRepository,
	result TermTranslationPhaseCommandReadModel,
	snapshotID int64,
) {
	t.Helper()
	if result.ErrorSummary == nil || result.ErrorSummary.ErrorKind != termTranslationErrorKindProviderFailure {
		t.Fatalf("expected provider failure, got %#v", result)
	}
	if result.PhaseRunID == nil {
		t.Fatalf("expected recreated phase run id, got %#v", result)
	}
	if len(jobRepo.phaseLinks) != 1 {
		t.Fatalf("expected one persisted snapshot link, got %#v", jobRepo.phaseLinks)
	}
	link := jobRepo.phaseLinks[0]
	if link.DictionaryEntryID != snapshotID || link.Role != termTranslationLinkRoleSnapshot {
		t.Fatalf("expected persisted snapshot link only, got %#v", jobRepo.phaseLinks)
	}
	assertNoTermTranslationJobEntries(t, foundation)
}

func assertNoTermTranslationJobEntries(
	t *testing.T,
	foundation *fakeTermPhaseFoundationDataRepository,
) {
	t.Helper()
	jobEntries, loadErr := foundation.ListDictionaryEntries(context.Background(), int64Pointer(1), termTranslationDictionaryLifecycleJob, "", "")
	if loadErr != nil {
		t.Fatalf("load job dictionary entries: %v", loadErr)
	}
	if len(jobEntries) != 0 {
		t.Fatalf("expected provider failure to leave no job dictionary entries, got %#v", jobEntries)
	}
}

func addTermTranslationMasterEntryAfterFailedStart(
	t *testing.T,
	foundation *fakeTermPhaseFoundationDataRepository,
	sourceTerm string,
	translatedTerm string,
) {
	t.Helper()
	_, err := foundation.CreateDictionaryEntry(context.Background(), repository.DictionaryEntryDraft{
		DictionaryLifecycle: termTranslationDictionaryLifecycleMaster,
		DictionaryScope:     "NPC_:FULL",
		DictionarySource:    "seed",
		SourceTerm:          sourceTerm,
		TranslatedTerm:      translatedTerm,
		Reusable:            true,
	})
	if err != nil {
		t.Fatalf("prepare new master entry after failed start: %v", err)
	}
}

func assertTermTranslationRetryUsesRemainingAITarget(
	t *testing.T,
	result TermTranslationPhaseCommandReadModel,
	capturedSources []string,
	expectedSource string,
) {
	t.Helper()
	if result.PhaseState != termTranslationPhaseStateCompleted {
		t.Fatalf("expected completed retry, got %#v", result)
	}
	if len(capturedSources) != 1 || capturedSources[0] != expectedSource {
		t.Fatalf("expected retry to preserve initial snapshot and translate only remaining term, got %#v", capturedSources)
	}
}

func TestTermTranslationPhaseServiceStartPhaseRejectsNonReadyJobWithCompletedRun(t *testing.T) {
	service, jobRepo, _ := newTermTranslationPhaseServiceForTest(nil)
	jobRepo.job.State = termTranslationJobStateRunning
	run := repository.JobPhaseRun{ID: 340, TranslationJobID: 1, PhaseType: termTranslationPhaseType, State: termTranslationPhaseStateCompleted}
	jobRepo.run = &run
	jobRepo.phases = append(jobRepo.phases, run)

	result, err := service.StartPhase(context.Background(), 1)
	if err != nil {
		t.Fatalf("expected rejected start without transport error, got %v", err)
	}
	if result.ErrorSummary == nil || result.ErrorSummary.ErrorKind != "ready_required" {
		t.Fatalf("expected ready_required rejection, got %#v", result)
	}
}

func TestTermTranslationPhaseServiceStartPhasePromotesPreCreatedIdleReadyRunToRunning(t *testing.T) {
	service, jobRepo, _ := newTermTranslationPhaseServiceForTest(fakeTermPhaseProvider{})
	existingRun := repository.JobPhaseRun{
		ID:               260,
		TranslationJobID: 1,
		PhaseType:        termTranslationPhaseType,
		State:            termTranslationPhaseStatePending,
		ExecutionMode:    "sync",
		AIProvider:       TermTranslationProviderGemini,
		ModelName:        "gemini-2.5-pro",
		CredentialRef:    "gemini-primary",
	}
	jobRepo.run = &existingRun
	jobRepo.phases = append(jobRepo.phases, existingRun)

	result, err := service.StartPhase(context.Background(), 1)
	if err != nil {
		t.Fatalf("expected start success with pre-created unstarted run: %v", err)
	}
	if result.PhaseRunID == nil || *result.PhaseRunID != existingRun.ID {
		t.Fatalf("expected existing phase run id reuse, got %#v", result.PhaseRunID)
	}
	if len(jobRepo.updatePhaseRunDraft) == 0 {
		t.Fatalf("expected phase run updates, got %#v", jobRepo.updatePhaseRunDraft)
	}
	if jobRepo.updatePhaseRunDraft[0].State != termTranslationPhaseStateRunning {
		t.Fatalf("expected first phase run update to running, got %#v", jobRepo.updatePhaseRunDraft[0])
	}
	if jobRepo.updateExpectedState[0] != termTranslationPhaseStatePending {
		t.Fatalf("expected optimistic expected state pending, got %#v", jobRepo.updateExpectedState)
	}
}

func TestTermTranslationPhaseServiceStartPhaseCreatesNonPendingRunWithoutPrecreatedPhaseRun(t *testing.T) {
	requestCount := 0
	service, jobRepo, _ := newTermTranslationPhaseServiceForTest(fakeTermPhaseProvider{
		translateFunc: func(_ context.Context, _ TermTranslationProviderRequest) (TermTranslationProviderResult, error) {
			requestCount++
			return TermTranslationProviderResult{}, nil
		},
	})
	jobRepo.run = nil

	result, err := service.StartPhase(context.Background(), 1)
	if err != nil {
		t.Fatalf("expected start success without pre-created phase run: %v", err)
	}
	if result.PhaseRunID == nil || *result.PhaseRunID == 0 {
		t.Fatalf("expected created phase run id, got %#v", result.PhaseRunID)
	}
	if jobRepo.run == nil {
		t.Fatal("expected created term translation phase run")
	}
	if jobRepo.run.State == termTranslationPhaseStatePending {
		t.Fatalf("expected created run not to be pending, got %#v", jobRepo.run)
	}
	if requestCount != 1 {
		t.Fatalf("expected provider request after phase run creation, got %d", requestCount)
	}
}

func TestTermTranslationPhaseServiceStartPhaseStopsBeforeProviderWhenPhasePromotionConflicts(t *testing.T) {
	requestCount := 0
	service, jobRepo, foundation := newTermTranslationPhaseServiceForTest(fakeTermPhaseProvider{
		translateFunc: func(_ context.Context, _ TermTranslationProviderRequest) (TermTranslationProviderResult, error) {
			requestCount++
			return TermTranslationProviderResult{}, nil
		},
	})
	existingRun := repository.JobPhaseRun{
		ID:               261,
		TranslationJobID: 1,
		PhaseType:        termTranslationPhaseType,
		State:            termTranslationPhaseStatePending,
		ExecutionMode:    "sync",
		AIProvider:       TermTranslationProviderGemini,
		ModelName:        "gemini-2.5-pro",
		CredentialRef:    "gemini-primary",
	}
	jobRepo.run = &existingRun
	jobRepo.phases = append(jobRepo.phases, existingRun)
	jobRepo.updatePhaseRunErr = repository.ErrConflict

	_, err := service.StartPhase(context.Background(), 1)
	if !errors.Is(err, repository.ErrConflict) {
		t.Fatalf("expected conflict error, got %v", err)
	}
	if requestCount != 0 {
		t.Fatalf("expected no provider request on phase promotion conflict, got %d", requestCount)
	}
	assertNoTermTranslationJobEntries(t, foundation)
}

func TestTermTranslationPhaseServiceResumePhasePromotesPausedRunWithExpectedState(t *testing.T) {
	service, jobRepo, _ := newTermTranslationPhaseServiceForTest(fakeTermPhaseProvider{})
	run := repository.JobPhaseRun{
		ID:               262,
		TranslationJobID: 1,
		PhaseType:        termTranslationPhaseType,
		State:            termTranslationPhaseStatePaused,
		ProgressPercent:  50,
		ExecutionMode:    "sync",
		AIProvider:       TermTranslationProviderGemini,
		ModelName:        "gemini-2.5-pro",
		CredentialRef:    "gemini-primary",
	}
	jobRepo.run = &run
	jobRepo.phases = append(jobRepo.phases, run)

	result, err := service.ResumePhase(context.Background(), 1, run.ID)
	if err != nil {
		t.Fatalf("expected resume success, got %v", err)
	}
	if result.PhaseRunID == nil || *result.PhaseRunID != run.ID {
		t.Fatalf("expected same phase run id, got %#v", result.PhaseRunID)
	}
	if len(jobRepo.updateExpectedState) == 0 || jobRepo.updateExpectedState[0] != termTranslationPhaseStatePaused {
		t.Fatalf("expected resume to update from paused state, got %#v", jobRepo.updateExpectedState)
	}
	if len(jobRepo.updatePhaseRunDraft) == 0 || jobRepo.updatePhaseRunDraft[0].State != termTranslationPhaseStateRunning {
		t.Fatalf("expected resume to promote phase run to running, got %#v", jobRepo.updatePhaseRunDraft)
	}
}

func TestTermTranslationPhaseServiceResumeAndRetryStopBeforeProviderWhenStatePromotionConflicts(t *testing.T) {
	tests := []struct {
		name          string
		mode          termTranslationStartMode
		runState      string
		expectedState string
	}{
		{
			name:          "resume",
			mode:          termTranslationStartModeResume,
			runState:      termTranslationPhaseStatePaused,
			expectedState: termTranslationPhaseStatePaused,
		},
		{
			name:          "retry",
			mode:          termTranslationStartModeRetry,
			runState:      termTranslationPhaseStateRecoverableFail,
			expectedState: termTranslationPhaseStateRecoverableFail,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			requestCount := 0
			service, jobRepo, _ := newTermTranslationPhaseServiceForTest(fakeTermPhaseProvider{
				translateFunc: func(_ context.Context, _ TermTranslationProviderRequest) (TermTranslationProviderResult, error) {
					requestCount++
					return TermTranslationProviderResult{}, nil
				},
			})
			run := repository.JobPhaseRun{
				ID:               263,
				TranslationJobID: 1,
				PhaseType:        termTranslationPhaseType,
				State:            tt.runState,
				ExecutionMode:    "sync",
				AIProvider:       TermTranslationProviderGemini,
				ModelName:        "gemini-2.5-pro",
				CredentialRef:    "gemini-primary",
			}
			jobRepo.run = &run
			jobRepo.phases = append(jobRepo.phases, run)
			jobRepo.updatePhaseRunErr = repository.ErrConflict
			expectedStateLog := make([]string, 0)
			jobRepo.expectedStateLog = &expectedStateLog

			err := runTermTranslationResumeOrRetry(t, service, tt.mode, run.ID)

			if !errors.Is(err, repository.ErrConflict) {
				t.Fatalf("expected conflict error, got %v", err)
			}
			if requestCount != 0 {
				t.Fatalf("expected no provider request on state promotion conflict, got %d", requestCount)
			}
			if len(expectedStateLog) == 0 || expectedStateLog[0] != tt.expectedState {
				t.Fatalf("expected state %q, got %#v", tt.expectedState, expectedStateLog)
			}
			if jobRepo.run == nil || jobRepo.run.State == termTranslationPhaseStateRunning {
				t.Fatalf("expected phase run to stay non-running, got %#v", jobRepo.run)
			}
		})
	}
}

func runTermTranslationResumeOrRetry(
	t *testing.T,
	service *TermTranslationPhaseService,
	mode termTranslationStartMode,
	phaseRunID int64,
) error {
	t.Helper()
	switch mode {
	case termTranslationStartModeResume:
		_, err := service.ResumePhase(context.Background(), 1, phaseRunID)
		return err
	case termTranslationStartModeRetry:
		_, err := service.RetryPhase(context.Background(), 1, phaseRunID)
		return err
	default:
		t.Fatalf("unexpected mode %q", mode)
		return nil
	}
}

func TestTermTranslationPhaseServiceResumePhaseRejectsNonResumableState(t *testing.T) {
	service, jobRepo, _ := newTermTranslationPhaseServiceForTest(nil)
	run := repository.JobPhaseRun{ID: 350, TranslationJobID: 1, PhaseType: termTranslationPhaseType, State: termTranslationPhaseStateCompleted}
	jobRepo.run = &run
	jobRepo.phases = append(jobRepo.phases, run)

	result, err := service.ResumePhase(context.Background(), 1, run.ID)
	if err != nil {
		t.Fatalf("expected rejected resume without transport error, got %v", err)
	}
	if result.ErrorSummary == nil || result.ErrorSummary.ErrorKind != "term_phase_incomplete" {
		t.Fatalf("expected term_phase_incomplete rejection, got %#v", result)
	}
}

func TestTermTranslationPhaseServiceUsesSecretStoreValueInsteadOfCredentialRef(t *testing.T) {
	capturedAPIKeys := make([]string, 0)
	service, _, _ := newTermTranslationPhaseServiceForTest(fakeTermPhaseProvider{
		translateFunc: func(_ context.Context, request TermTranslationProviderRequest) (TermTranslationProviderResult, error) {
			capturedAPIKeys = append(capturedAPIKeys, request.APIKey)
			return TermTranslationProviderResult{SourceTerm: request.SourceTerm, TranslatedTerm: request.SourceTerm + "_ja"}, nil
		},
	})

	_, err := service.StartPhase(context.Background(), 1)
	if err != nil {
		t.Fatalf("expected start success, got %v", err)
	}
	if len(capturedAPIKeys) == 0 {
		t.Fatal("expected provider request capture")
	}
	for _, apiKey := range capturedAPIKeys {
		if apiKey != "gemini-secret" {
			t.Fatalf("expected provider auth to use secret store value, got %#v", capturedAPIKeys)
		}
	}
}

func TestTermTranslationPhaseServiceLoadsSecretBySnapshotCredentialRef(t *testing.T) {
	secretLoadKeys := make([]string, 0, 1)
	capturedAPIKeys := make([]string, 0, 1)
	service, _, _ := newTermTranslationPhaseServiceForTest(fakeTermPhaseProvider{
		translateFunc: func(_ context.Context, request TermTranslationProviderRequest) (TermTranslationProviderResult, error) {
			capturedAPIKeys = append(capturedAPIKeys, request.APIKey)
			return TermTranslationProviderResult{SourceTerm: request.SourceTerm, TranslatedTerm: request.SourceTerm + "_ja"}, nil
		},
	})
	service.secretStore = fakeTermPhaseSecretStore{
		loadFunc: func(_ context.Context, key string) (string, error) {
			secretLoadKeys = append(secretLoadKeys, key)
			if key == "gemini-primary" {
				return "gemini-secret", nil
			}
			return "", nil
		},
	}

	_, err := service.StartPhase(context.Background(), 1)
	if err != nil {
		t.Fatalf("expected start success: %v", err)
	}
	if len(secretLoadKeys) == 0 || secretLoadKeys[0] != "gemini-primary" {
		t.Fatalf("expected credential lookup by phase snapshot ref, got %#v", secretLoadKeys)
	}
	if len(capturedAPIKeys) == 0 || capturedAPIKeys[0] != "gemini-secret" {
		t.Fatalf("expected provider call to use secret loaded by credential ref, got %#v", capturedAPIKeys)
	}
}

func TestTermTranslationPhaseServiceStartPhaseReResolvesProviderSettingsBeforeExecution(t *testing.T) {
	capturedRequests := make([]TermTranslationProviderRequest, 0, 1)
	service, _, _ := newTermTranslationPhaseServiceForTest(fakeTermPhaseProvider{
		translateFunc: func(_ context.Context, request TermTranslationProviderRequest) (TermTranslationProviderResult, error) {
			capturedRequests = append(capturedRequests, request)
			return TermTranslationProviderResult{SourceTerm: request.SourceTerm, TranslatedTerm: request.SourceTerm + "_ja"}, nil
		},
	})
	service.secretStore = fakeTermPhaseSecretStore{
		loadFunc: func(_ context.Context, key string) (string, error) {
			if key == "provider-settings:gemini" {
				return "provider-settings-secret", nil
			}
			return "", nil
		},
	}
	service.WithTermTranslationProviderSettings(fakePhaseProviderSettingsConsumer{
		resolveFunc: func(_ context.Context, input ProviderSettingsResolveInput) (ProviderSettingsResolveResult, error) {
			return ProviderSettingsResolveResult{
				ConsumerID:            input.ConsumerID,
				ProviderID:            input.Selection.ProviderID,
				Model:                 input.Selection.Model,
				ExecutionMethod:       input.Selection.ExecutionMethod,
				Endpoint:              stringPointer("http://localhost:1234/v1"),
				CredentialReferenceID: stringPointer("provider-settings:gemini"),
				CredentialState:       providerSettingsCredentialStateConfigured,
			}, nil
		},
	})

	_, err := service.StartPhase(context.Background(), 1)
	if err != nil {
		t.Fatalf("expected start success: %v", err)
	}
	if len(capturedRequests) == 0 {
		t.Fatal("expected provider request capture")
	}
	if capturedRequests[0].APIKey != "provider-settings-secret" {
		t.Fatalf("expected provider request to use re-resolved secret, got %#v", capturedRequests[0])
	}
	if capturedRequests[0].EndpointSummary == nil || *capturedRequests[0].EndpointSummary != "http://localhost:1234/v1" {
		t.Fatalf("expected provider request to use re-resolved endpoint summary, got %#v", capturedRequests[0])
	}
	if capturedRequests[0].RequestUnitID == "" {
		t.Fatalf("expected provider request to carry request unit id, got %#v", capturedRequests[0])
	}
}

func TestTermTranslationPhaseServiceStartPhaseFallsBackWhenLMStudioSecretLoadStalls(t *testing.T) {
	blockedLoad := make(chan struct{})
	capturedAPIKeys := make([]string, 0, 2)
	service, jobRepo, _ := newTermTranslationPhaseServiceForTest(fakeTermPhaseProvider{
		translateFunc: func(_ context.Context, request TermTranslationProviderRequest) (TermTranslationProviderResult, error) {
			capturedAPIKeys = append(capturedAPIKeys, request.APIKey)
			return TermTranslationProviderResult{SourceTerm: request.SourceTerm, TranslatedTerm: request.SourceTerm + "_ja"}, nil
		},
	})
	jobRepo.phases[0].AIProvider = TermTranslationProviderLMStudio
	jobRepo.phases[0].ModelName = "lmstudio-community"
	jobRepo.phases[0].CredentialRef = "lmstudio-local"
	service.secretStore = fakeTermPhaseSecretStore{
		loadFunc: func(context.Context, string) (string, error) {
			<-blockedLoad
			return "", nil
		},
	}
	service.secretLoadTimeout = 5 * time.Millisecond

	startedAt := time.Now()
	result, err := service.StartPhase(context.Background(), 1)
	if err != nil {
		t.Fatalf("expected lm_studio start to fall back without error: %v", err)
	}
	if time.Since(startedAt) > 150*time.Millisecond {
		t.Fatalf("expected stalled lm_studio secret load to time out quickly, got %s", time.Since(startedAt))
	}
	if result.PhaseState != termTranslationPhaseStateCompleted {
		t.Fatalf("expected completed lm_studio start result, got %#v", result)
	}
	if len(capturedAPIKeys) == 0 {
		t.Fatal("expected provider request capture")
	}
	for _, apiKey := range capturedAPIKeys {
		if apiKey != "" {
			t.Fatalf("expected lm_studio provider calls without secret, got %#v", capturedAPIKeys)
		}
	}
}

func TestTermTranslationPhaseServiceStartPhaseFailsQuicklyWhenRequiredSecretLoadStalls(t *testing.T) {
	blockedLoad := make(chan struct{})
	providerCalled := false
	service, _, _ := newTermTranslationPhaseServiceForTest(fakeTermPhaseProvider{
		translateFunc: func(_ context.Context, request TermTranslationProviderRequest) (TermTranslationProviderResult, error) {
			providerCalled = true
			return TermTranslationProviderResult{SourceTerm: request.SourceTerm, TranslatedTerm: request.SourceTerm + "_ja"}, nil
		},
	})
	service.secretStore = fakeTermPhaseSecretStore{
		loadFunc: func(context.Context, string) (string, error) {
			<-blockedLoad
			return "", nil
		},
	}
	service.secretLoadTimeout = 5 * time.Millisecond

	startedAt := time.Now()
	result, err := service.StartPhase(context.Background(), 1)
	if err != nil {
		t.Fatalf("expected provider failure command result without transport error, got %v", err)
	}
	if time.Since(startedAt) > 150*time.Millisecond {
		t.Fatalf("expected stalled required secret load to fail quickly, got %s", time.Since(startedAt))
	}
	if providerCalled {
		t.Fatal("expected provider call to stay blocked when required secret is unavailable")
	}
	if result.ErrorSummary == nil || result.ErrorSummary.ErrorKind != termTranslationErrorKindProviderFailure {
		t.Fatalf("expected provider failure error summary, got %#v", result)
	}
	if result.PhaseState != termTranslationPhaseStateFailed {
		t.Fatalf("expected failed phase state for required secret stall, got %#v", result)
	}
}

func TestTermTranslationPhaseServiceReturnsSaveFailedAndRollsBackPartialDictionaryState(t *testing.T) {
	service, jobRepo, foundation := newTermTranslationPhaseServiceForTest(fakeTermPhaseProvider{
		translateFunc: func(_ context.Context, request TermTranslationProviderRequest) (TermTranslationProviderResult, error) {
			return TermTranslationProviderResult{SourceTerm: request.SourceTerm, TranslatedTerm: request.SourceTerm + "_ja"}, nil
		},
	})
	foundation.createErrBySource["Steel Sword"] = errors.New("disk full")

	result, err := service.StartPhase(context.Background(), 1)
	if err != nil {
		t.Fatalf("expected save failure command result, got %v", err)
	}
	if result.ErrorSummary == nil || result.ErrorSummary.ErrorKind != termTranslationErrorKindSaveFailed {
		t.Fatalf("expected save_failed error summary, got %#v", result)
	}
	if result.PhaseState != termTranslationPhaseStateFailed {
		t.Fatalf("expected failed phase state, got %#v", result)
	}
	assertTermTranslationJobStateRunning(t, jobRepo)
	jobEntries, loadErr := foundation.ListDictionaryEntries(context.Background(), int64Pointer(1), termTranslationDictionaryLifecycleJob, "", "")
	if loadErr != nil {
		t.Fatalf("load job dictionary entries: %v", loadErr)
	}
	if len(jobEntries) != 0 {
		t.Fatalf("expected failed attempt to leave no job dictionary entries, got %#v", jobEntries)
	}
	if len(jobRepo.phaseLinks) != 0 {
		t.Fatalf("expected failed attempt to leave no phase links, got %#v", jobRepo.phaseLinks)
	}
}

func TestTermTranslationPhaseServiceRetryPhaseRejectsNonRetryableFailureState(t *testing.T) {
	service, jobRepo, _ := newTermTranslationPhaseServiceForTest(nil)
	run := repository.JobPhaseRun{
		ID:               370,
		TranslationJobID: 1,
		PhaseType:        termTranslationPhaseType,
		State:            termTranslationPhaseStateFailed,
		LatestError:      termTranslationErrorKindSaveFailed,
	}
	jobRepo.job.State = termTranslationJobStateRunning
	jobRepo.run = &run
	jobRepo.phases = append(jobRepo.phases, run)

	result, err := service.RetryPhase(context.Background(), 1, run.ID)
	if err != nil {
		t.Fatalf("expected rejected retry without transport error, got %v", err)
	}
	if result.Retryable {
		t.Fatalf("expected retryable=false, got %#v", result)
	}
	if result.ErrorSummary == nil || result.ErrorSummary.ErrorKind != "term_phase_incomplete" {
		t.Fatalf("expected term_phase_incomplete rejection for failed run, got %#v", result)
	}
}

func assertTermTranslationJobStateRunning(t *testing.T, jobRepo *fakeTermPhaseJobLifecycleRepository) {
	t.Helper()
	if jobRepo.job.State != termTranslationJobStateRunning {
		t.Fatalf("expected translation job state running, got %q", jobRepo.job.State)
	}
	if len(jobRepo.updateJobDrafts) == 0 {
		t.Fatal("expected translation job update draft")
	}
	lastDraft := jobRepo.updateJobDrafts[len(jobRepo.updateJobDrafts)-1]
	if lastDraft.State != termTranslationJobStateRunning {
		t.Fatalf("expected last translation job draft state running, got %#v", lastDraft)
	}
}

func TestTermTranslationPhaseServiceReadNextPhaseReadinessUsesCandidateCoverage(t *testing.T) {
	service, jobRepo, foundation := newTermTranslationPhaseServiceForTest(nil)
	run := repository.JobPhaseRun{ID: 360, TranslationJobID: 1, PhaseType: termTranslationPhaseType, State: termTranslationPhaseStateCompleted, ProgressPercent: 100}
	jobRepo.run = &run
	jobRepo.phases = append(jobRepo.phases, run)

	_, err := foundation.CreateDictionaryEntry(context.Background(), repository.DictionaryEntryDraft{
		TranslationJobID:    int64Pointer(1),
		DictionaryLifecycle: termTranslationDictionaryLifecycleJob,
		DictionaryScope:     "MISC:FULL",
		DictionarySource:    termTranslationDictionarySourceProvider,
		SourceTerm:          "Unrelated Term",
		TranslatedTerm:      "無関係",
		Reusable:            true,
	})
	if err != nil {
		t.Fatalf("prepare unrelated job entry: %v", err)
	}
	_, err = foundation.CreateDictionaryEntry(context.Background(), repository.DictionaryEntryDraft{
		TranslationJobID:    int64Pointer(1),
		DictionaryLifecycle: termTranslationDictionaryLifecycleJob,
		DictionaryScope:     "NPC_:FULL",
		DictionarySource:    termTranslationDictionarySourceProvider,
		SourceTerm:          "Iron Sword",
		TranslatedTerm:      "鉄の剣",
		Reusable:            true,
	})
	if err != nil {
		t.Fatalf("prepare matched job entry: %v", err)
	}

	readiness, err := service.ReadNextPhaseReadiness(context.Background(), 1)
	if err != nil {
		t.Fatalf("expected readiness load success, got %v", err)
	}
	if readiness.ConfirmedCount >= readiness.TotalCount && readiness.TotalCount > 0 {
		t.Fatalf("expected candidate coverage to be incomplete, got %#v", readiness)
	}
}

func TestTermTranslationPhaseServiceReadNextPhaseReadinessKeepsReadyJobWithoutPhaseRunsBlocked(t *testing.T) {
	service, jobRepo, _ := newTermTranslationPhaseServiceForTest(nil)
	jobRepo.run = nil
	jobRepo.phases = nil

	readiness, err := service.ReadNextPhaseReadiness(context.Background(), 1)
	if err != nil {
		t.Fatalf("expected readiness success for ready job without phase runs: %v", err)
	}
	if readiness.PhaseState == termTranslationPhaseStateCompleted {
		t.Fatalf("expected phase state to indicate not completed, got %#v", readiness)
	}
}

func TestTermTranslationPhaseServiceReadNextPhaseReadinessReturnsNotFoundForNonReadyJobWithoutPhaseRuns(t *testing.T) {
	service, jobRepo, _ := newTermTranslationPhaseServiceForTest(nil)
	jobRepo.job.State = termTranslationJobStateRunning
	jobRepo.run = nil
	jobRepo.phases = nil

	_, err := service.ReadNextPhaseReadiness(context.Background(), 1)
	if err == nil {
		t.Fatal("expected service error for non-ready job without phase runs")
	}
	if !strings.Contains(err.Error(), "load initial execution phase: not found") {
		t.Fatalf("expected load initial execution phase not found error, got %v", err)
	}
}

// U-TTRC-002: collectCandidates が 13 種別だけを候補化し、13 種別外 REC を除外する
func TestCollectCandidates_filtersBy13RECTypes(t *testing.T) {
	// 13 種別 REC と 13 種別外 REC（DOOR:FULL / FLOR:FULL / FURN:FULL）を混在させた入力で
	// collectCandidates が 13 種別のみを返し、13 種別外を除外することを証明する。
	service, _, _ := newTermTranslationPhaseServiceForTest(nil)

	// 13 種別外 REC を含む入力: DOOR / FLOR / FURN は除外対象
	service.translationSourceReader = fakeTermPhaseTranslationSourceRepository{
		records: []repository.TranslationRecord{
			{ID: 1, RecordType: "NPC_", XEditExtractedDataID: 1},
			{ID: 2, RecordType: "BOOK", XEditExtractedDataID: 1},
			{ID: 3, RecordType: "DOOR", XEditExtractedDataID: 1},
			{ID: 4, RecordType: "FLOR", XEditExtractedDataID: 1},
			{ID: 5, RecordType: "FURN", XEditExtractedDataID: 1},
		},
		fields: map[int64][]repository.TranslationField{
			1: {{ID: 11, TranslationRecordID: 1, SubrecordType: "FULL", SourceText: "Nord Warrior"}},
			2: {{ID: 21, TranslationRecordID: 2, SubrecordType: "FULL", SourceText: "Ancient Tome"}},
			3: {{ID: 31, TranslationRecordID: 3, SubrecordType: "FULL", SourceText: "Old Door"}},
			4: {{ID: 41, TranslationRecordID: 4, SubrecordType: "FULL", SourceText: "Garden Plant"}},
			5: {{ID: 51, TranslationRecordID: 5, SubrecordType: "FULL", SourceText: "Wooden Chair"}},
		},
	}

	// Act
	candidates, err := service.collectCandidates(context.Background(), 1)
	if err != nil {
		t.Fatalf("unexpected error from collectCandidates: %v", err)
	}

	// Assert: 13 種別に該当する NPC_:FULL と BOOK:FULL だけが候補に残る
	if len(candidates) != 2 {
		t.Fatalf("REC 絞り込み: expected 2 candidates (NPC_:FULL, BOOK:FULL), got %d: %#v", len(candidates), candidates)
	}
	recTypes := make(map[string]bool, len(candidates))
	for _, c := range candidates {
		recTypes[c.RecordType] = true
	}
	if !recTypes["NPC_:FULL"] {
		t.Fatalf("REC 絞り込み: expected NPC_:FULL in candidates, got %#v", recTypes)
	}
	if !recTypes["BOOK:FULL"] {
		t.Fatalf("REC 絞り込み: expected BOOK:FULL in candidates, got %#v", recTypes)
	}
	// 13 種別外は除外されている
	if recTypes["DOOR:FULL"] || recTypes["FLOR:FULL"] || recTypes["FURN:FULL"] {
		t.Fatalf("REC 絞り込み欠落: expected DOOR/FLOR/FURN to be excluded, got %#v", recTypes)
	}
}

// U-TTRC-002 境界: 原語空のレコードは候補から除外される
func TestCollectCandidates_excludesEmptySourceTerm(t *testing.T) {
	// 原語が空文字のフィールドは候補化対象外になることを証明する。
	service, _, _ := newTermTranslationPhaseServiceForTest(nil)

	service.translationSourceReader = fakeTermPhaseTranslationSourceRepository{
		records: []repository.TranslationRecord{
			{ID: 1, RecordType: "NPC_", XEditExtractedDataID: 1},
		},
		fields: map[int64][]repository.TranslationField{
			1: {
				{ID: 11, TranslationRecordID: 1, SubrecordType: "FULL", SourceText: ""},
				{ID: 12, TranslationRecordID: 1, SubrecordType: "FULL", SourceText: "   "},
				{ID: 13, TranslationRecordID: 1, SubrecordType: "SHRT", SourceText: "Warrior"},
			},
		},
	}

	candidates, err := service.collectCandidates(context.Background(), 1)
	if err != nil {
		t.Fatalf("unexpected error from collectCandidates: %v", err)
	}

	// 原語空チェック欠落: 空文字と空白のみのフィールドは除外され、非空の SHRT だけが残る
	if len(candidates) != 1 {
		t.Fatalf("原語空チェック欠落: expected 1 candidate (SHRT with non-empty source), got %d: %#v", len(candidates), candidates)
	}
	if candidates[0].SourceTerm != "Warrior" {
		t.Fatalf("expected candidate source term 'Warrior', got %q", candidates[0].SourceTerm)
	}
}

// U-TTRC-003: candidate key が RECORD:FIELD 形式であることを証明する
func TestCollectCandidates_candidateKeyIsRECORDFieldFormat(t *testing.T) {
	// candidate の RecordType フィールドが "RECORD:FIELD" 形式（例: NPC_:FULL）で組み立てられることを証明する。
	service, _, _ := newTermTranslationPhaseServiceForTest(nil)

	service.translationSourceReader = fakeTermPhaseTranslationSourceRepository{
		records: []repository.TranslationRecord{
			{ID: 1, RecordType: "ARMO", XEditExtractedDataID: 1},
		},
		fields: map[int64][]repository.TranslationField{
			1: {{ID: 11, TranslationRecordID: 1, SubrecordType: "FULL", SourceText: "Iron Armor"}},
		},
	}

	candidates, err := service.collectCandidates(context.Background(), 1)
	if err != nil {
		t.Fatalf("unexpected error from collectCandidates: %v", err)
	}

	if len(candidates) != 1 {
		t.Fatalf("expected 1 candidate, got %d: %#v", len(candidates), candidates)
	}

	// candidate の RecordType が "ARMO:FULL" 形式であることを検証する
	if candidates[0].RecordType != "ARMO:FULL" {
		t.Fatalf("candidate key 形式違反: expected RecordType 'ARMO:FULL', got %q", candidates[0].RecordType)
	}
}

// U-TTRC-003: NPC_:FULL と NPC_:SHRT の同一原語が別候補として保持される
func TestCollectCandidates_NPC_FULLAndNPC_SHRTAreDistinctCandidates(t *testing.T) {
	// NPC_:FULL と NPC_:SHRT は異なる REC として扱われ、同一原語でも別候補になることを証明する。
	// 退行条件: record_type のみで集約すると 1 件になる。
	service, _, _ := newTermTranslationPhaseServiceForTest(nil)

	service.translationSourceReader = fakeTermPhaseTranslationSourceRepository{
		records: []repository.TranslationRecord{
			{ID: 1, RecordType: "NPC_", XEditExtractedDataID: 1},
		},
		fields: map[int64][]repository.TranslationField{
			1: {
				{ID: 11, TranslationRecordID: 1, SubrecordType: "FULL", SourceText: "Lydia"},
				{ID: 12, TranslationRecordID: 1, SubrecordType: "SHRT", SourceText: "Lydia"},
			},
		},
	}

	candidates, err := service.collectCandidates(context.Background(), 1)
	if err != nil {
		t.Fatalf("unexpected error from collectCandidates: %v", err)
	}

	// NPC_:FULL と NPC_:SHRT は別候補として 2 件返る
	if len(candidates) != 2 {
		t.Fatalf("NPC_ REC 分離: expected 2 candidates (NPC_:FULL and NPC_:SHRT with same source), got %d: %#v", len(candidates), candidates)
	}

	recTypes := make(map[string]bool, len(candidates))
	for _, c := range candidates {
		recTypes[c.RecordType] = true
	}
	if !recTypes["NPC_:FULL"] {
		t.Fatalf("NPC_ REC 分離: expected NPC_:FULL candidate, got %#v", recTypes)
	}
	if !recTypes["NPC_:SHRT"] {
		t.Fatalf("NPC_ REC 分離: expected NPC_:SHRT candidate, got %#v", recTypes)
	}
}
