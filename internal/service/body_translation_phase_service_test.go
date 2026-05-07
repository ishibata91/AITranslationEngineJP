package service

import (
	"context"
	"testing"
	"time"

	"aitranslationenginejp/internal/repository"
)

type fakeBodyPhaseTransactor struct {
	jobRepo    *fakeBodyPhaseJobLifecycleRepository
	outputRepo *fakeBodyPhaseOutputRepository
}

func (fake fakeBodyPhaseTransactor) WithTransaction(ctx context.Context, fn func(context.Context) error) error {
	jobSnapshot := fake.jobRepo.clone()
	outputSnapshot := fake.outputRepo.clone()
	err := fn(ctx)
	if err != nil {
		fake.jobRepo.restore(jobSnapshot)
		fake.outputRepo.restore(outputSnapshot)
	}
	return err
}

type fakeBodyPhaseProvider struct {
	translateFunc func(context.Context, BodyTranslationProviderRequest) BodyTranslationProviderResult
}

func (fake fakeBodyPhaseProvider) TranslateBodyField(
	ctx context.Context,
	request BodyTranslationProviderRequest,
) BodyTranslationProviderResult {
	if fake.translateFunc != nil {
		return fake.translateFunc(ctx, request)
	}
	return BodyTranslationProviderResult{}
}

func (fake fakeBodyPhaseProvider) BodyTranslationProviderRequestsAreTestSafe() bool {
	return true
}

type fakeBodyPhaseJobLifecycleRepository struct {
	job                     repository.TranslationJob
	runsByPhaseType         map[string]repository.JobPhaseRun
	updateJobDrafts         []repository.TranslationJobUpdateDraft
	updatePhaseRunDrafts    []repository.JobPhaseRunUpdateDraft
	phaseRunFieldLinks      []repository.PhaseRunTranslationField
	nextPhaseRunID          int64
	nextPhaseRunFieldLinkID int64
}

func (fake *fakeBodyPhaseJobLifecycleRepository) GetTranslationJobByID(context.Context, int64) (repository.TranslationJob, error) {
	return fake.job, nil
}

func (fake *fakeBodyPhaseJobLifecycleRepository) UpdateTranslationJob(
	_ context.Context,
	_ int64,
	draft repository.TranslationJobUpdateDraft,
) (repository.TranslationJob, error) {
	fake.updateJobDrafts = append(fake.updateJobDrafts, draft)
	fake.job.JobName = draft.JobName
	fake.job.State = draft.State
	fake.job.ProgressPercent = draft.ProgressPercent
	fake.job.StartedAt = draft.StartedAt
	fake.job.FinishedAt = draft.FinishedAt
	return fake.job, nil
}

func (fake *fakeBodyPhaseJobLifecycleRepository) CreateJobPhaseRun(
	_ context.Context,
	draft repository.JobPhaseRunDraft,
) (repository.JobPhaseRun, error) {
	if _, exists := fake.runsByPhaseType[draft.PhaseType]; exists {
		return repository.JobPhaseRun{}, repository.ErrConflict
	}
	fake.nextPhaseRunID++
	run := repository.JobPhaseRun{
		ID:                     fake.nextPhaseRunID,
		TranslationJobID:       draft.TranslationJobID,
		PhaseType:              draft.PhaseType,
		State:                  draft.State,
		ExecutionOrder:         draft.ExecutionOrder,
		SnapshotFieldCount:     draft.SnapshotFieldCount,
		ProviderTargetCount:    draft.ProviderTargetCount,
		ExactExclusionCount:    draft.ExactExclusionCount,
		PartialConstraintCount: draft.PartialConstraintCount,
		AIProvider:             draft.AIProvider,
		ModelName:              draft.ModelName,
		ExecutionMode:          draft.ExecutionMode,
		CredentialRef:          draft.CredentialRef,
		InstructionKind:        draft.InstructionKind,
		InputSnapshotDigest:    draft.InputSnapshotDigest,
		DictionaryDigest:       draft.DictionaryDigest,
		PersonaDigest:          draft.PersonaDigest,
		MetadataDigest:         draft.MetadataDigest,
		PromptDigest:           draft.PromptDigest,
	}
	fake.runsByPhaseType[draft.PhaseType] = run
	return run, nil
}

func (fake *fakeBodyPhaseJobLifecycleRepository) FindJobPhaseRun(
	_ context.Context,
	_ int64,
	phaseType string,
) (repository.JobPhaseRun, error) {
	run, ok := fake.runsByPhaseType[phaseType]
	if !ok {
		return repository.JobPhaseRun{}, repository.ErrNotFound
	}
	return run, nil
}

func (fake *fakeBodyPhaseJobLifecycleRepository) UpdateJobPhaseRun(
	_ context.Context,
	id int64,
	draft repository.JobPhaseRunUpdateDraft,
) (repository.JobPhaseRun, error) {
	for phaseType, run := range fake.runsByPhaseType {
		if run.ID != id {
			continue
		}
		fake.updatePhaseRunDrafts = append(fake.updatePhaseRunDrafts, draft)
		run.State = draft.State
		run.ProgressPercent = draft.ProgressPercent
		run.LatestExternalRunID = draft.LatestExternalRunID
		run.LatestError = draft.LatestError
		run.StartedAt = draft.StartedAt
		run.FinishedAt = draft.FinishedAt
		fake.runsByPhaseType[phaseType] = run
		return run, nil
	}
	return repository.JobPhaseRun{}, repository.ErrNotFound
}

func (fake *fakeBodyPhaseJobLifecycleRepository) ListJobPhaseRunsByJobID(context.Context, int64) ([]repository.JobPhaseRun, error) {
	result := make([]repository.JobPhaseRun, 0, len(fake.runsByPhaseType))
	for _, run := range fake.runsByPhaseType {
		result = append(result, run)
	}
	return result, nil
}

func (fake *fakeBodyPhaseJobLifecycleRepository) CreatePhaseRunTranslationField(
	_ context.Context,
	draft repository.PhaseRunTranslationFieldDraft,
) (repository.PhaseRunTranslationField, error) {
	for _, link := range fake.phaseRunFieldLinks {
		if link.PhaseRunID == draft.PhaseRunID &&
			link.JobTranslationFieldID == draft.JobTranslationFieldID &&
			link.Role == draft.Role {
			return repository.PhaseRunTranslationField{}, repository.ErrConflict
		}
	}
	fake.nextPhaseRunFieldLinkID++
	link := repository.PhaseRunTranslationField{
		ID:                    fake.nextPhaseRunFieldLinkID,
		PhaseRunID:            draft.PhaseRunID,
		JobTranslationFieldID: draft.JobTranslationFieldID,
		Role:                  draft.Role,
	}
	fake.phaseRunFieldLinks = append(fake.phaseRunFieldLinks, link)
	return link, nil
}

func (fake *fakeBodyPhaseJobLifecycleRepository) ListPhaseRunTranslationFieldsByPhaseRunID(
	_ context.Context,
	phaseRunID int64,
) ([]repository.PhaseRunTranslationField, error) {
	result := make([]repository.PhaseRunTranslationField, 0)
	for _, link := range fake.phaseRunFieldLinks {
		if link.PhaseRunID == phaseRunID {
			result = append(result, link)
		}
	}
	return result, nil
}

func (fake *fakeBodyPhaseJobLifecycleRepository) clone() *fakeBodyPhaseJobLifecycleRepository {
	cloned := *fake
	cloned.runsByPhaseType = make(map[string]repository.JobPhaseRun, len(fake.runsByPhaseType))
	for phaseType, run := range fake.runsByPhaseType {
		cloned.runsByPhaseType[phaseType] = run
	}
	cloned.updateJobDrafts = append([]repository.TranslationJobUpdateDraft(nil), fake.updateJobDrafts...)
	cloned.updatePhaseRunDrafts = append([]repository.JobPhaseRunUpdateDraft(nil), fake.updatePhaseRunDrafts...)
	cloned.phaseRunFieldLinks = append([]repository.PhaseRunTranslationField(nil), fake.phaseRunFieldLinks...)
	return &cloned
}

func (fake *fakeBodyPhaseJobLifecycleRepository) restore(snapshot *fakeBodyPhaseJobLifecycleRepository) {
	*fake = *snapshot
	fake.runsByPhaseType = make(map[string]repository.JobPhaseRun, len(snapshot.runsByPhaseType))
	for phaseType, run := range snapshot.runsByPhaseType {
		fake.runsByPhaseType[phaseType] = run
	}
	fake.updateJobDrafts = append([]repository.TranslationJobUpdateDraft(nil), snapshot.updateJobDrafts...)
	fake.updatePhaseRunDrafts = append([]repository.JobPhaseRunUpdateDraft(nil), snapshot.updatePhaseRunDrafts...)
	fake.phaseRunFieldLinks = append([]repository.PhaseRunTranslationField(nil), snapshot.phaseRunFieldLinks...)
}

type fakeBodyPhaseFoundationDataRepository struct {
	dictionary []repository.DictionaryEntry
	personas   []repository.Persona
}

func (fake fakeBodyPhaseFoundationDataRepository) ListDictionaryEntries(
	context.Context,
	*int64,
	string,
	string,
	string,
) ([]repository.DictionaryEntry, error) {
	return append([]repository.DictionaryEntry(nil), fake.dictionary...), nil
}

func (fake fakeBodyPhaseFoundationDataRepository) ListPersonasByTranslationJobID(
	context.Context,
	int64,
) ([]repository.Persona, error) {
	return append([]repository.Persona(nil), fake.personas...), nil
}

type fakeBodyPhaseTranslationSourceRepository struct {
	xedit   repository.XEditExtractedData
	records []repository.TranslationRecord
	fields  map[int64][]repository.TranslationField
}

func (fake fakeBodyPhaseTranslationSourceRepository) GetXEditExtractedDataByID(
	context.Context,
	int64,
) (repository.XEditExtractedData, error) {
	return fake.xedit, nil
}

func (fake fakeBodyPhaseTranslationSourceRepository) ListTranslationRecordsByXEditID(
	context.Context,
	int64,
) ([]repository.TranslationRecord, error) {
	return append([]repository.TranslationRecord(nil), fake.records...), nil
}

func (fake fakeBodyPhaseTranslationSourceRepository) ListTranslationFieldsByTranslationRecordID(
	_ context.Context,
	recordID int64,
) ([]repository.TranslationField, error) {
	return append([]repository.TranslationField(nil), fake.fields[recordID]...), nil
}

type fakeBodyPhaseOutputRepository struct {
	fields  []repository.JobTranslationField
	nextID  int64
	nowFunc func() time.Time
}

func (fake *fakeBodyPhaseOutputRepository) CreateJobTranslationField(
	_ context.Context,
	draft repository.JobTranslationFieldDraft,
) (repository.JobTranslationField, error) {
	for _, field := range fake.fields {
		if field.TranslationJobID == draft.TranslationJobID &&
			field.TranslationFieldID == draft.TranslationFieldID {
			return repository.JobTranslationField{}, repository.ErrConflict
		}
	}
	fake.nextID++
	field := repository.JobTranslationField{
		ID:                 fake.nextID,
		TranslationJobID:   draft.TranslationJobID,
		TranslationFieldID: draft.TranslationFieldID,
		AppliedPersonaID:   draft.AppliedPersonaID,
		TranslatedText:     draft.TranslatedText,
		OutputStatus:       draft.OutputStatus,
		RetryCount:         draft.RetryCount,
		UpdatedAt:          fake.now(),
	}
	fake.fields = append(fake.fields, field)
	return field, nil
}

func (fake *fakeBodyPhaseOutputRepository) UpdateJobTranslationField(
	_ context.Context,
	id int64,
	draft repository.JobTranslationFieldUpdateDraft,
) (repository.JobTranslationField, error) {
	for index, field := range fake.fields {
		if field.ID != id {
			continue
		}
		field.AppliedPersonaID = draft.AppliedPersonaID
		field.TranslatedText = draft.TranslatedText
		field.OutputStatus = draft.OutputStatus
		field.RetryCount = draft.RetryCount
		field.UpdatedAt = fake.now()
		fake.fields[index] = field
		return field, nil
	}
	return repository.JobTranslationField{}, repository.ErrNotFound
}

func (fake *fakeBodyPhaseOutputRepository) ListJobTranslationFieldsByJobID(
	_ context.Context,
	jobID int64,
) ([]repository.JobTranslationField, error) {
	result := make([]repository.JobTranslationField, 0)
	for _, field := range fake.fields {
		if field.TranslationJobID == jobID {
			result = append(result, field)
		}
	}
	return result, nil
}

func (fake *fakeBodyPhaseOutputRepository) clone() *fakeBodyPhaseOutputRepository {
	cloned := *fake
	cloned.fields = append([]repository.JobTranslationField(nil), fake.fields...)
	return &cloned
}

func (fake *fakeBodyPhaseOutputRepository) restore(snapshot *fakeBodyPhaseOutputRepository) {
	*fake = *snapshot
	fake.fields = append([]repository.JobTranslationField(nil), snapshot.fields...)
}

func (fake *fakeBodyPhaseOutputRepository) now() time.Time {
	if fake.nowFunc != nil {
		return fake.nowFunc()
	}
	return time.Date(2026, 5, 3, 12, 0, 0, 0, time.UTC)
}

func newBodyTranslationPhaseServiceForTest(
	provider BodyTranslationProvider,
) (*BodyTranslationPhaseService, *fakeBodyPhaseJobLifecycleRepository, *fakeBodyPhaseOutputRepository) {
	jobRepo := &fakeBodyPhaseJobLifecycleRepository{
		job: repository.TranslationJob{
			ID:                   1,
			XEditExtractedDataID: 10,
			JobName:              "job-1",
			State:                bodyTranslationPhaseStateIdleReady,
		},
		runsByPhaseType: map[string]repository.JobPhaseRun{
			bodyTranslationPersonaPhaseType: {
				ID:               100,
				TranslationJobID: 1,
				PhaseType:        bodyTranslationPersonaPhaseType,
				State:            bodyTranslationPhaseStateCompleted,
				ExecutionMode:    "sync",
				AIProvider:       BodyTranslationProviderLMStudio,
				ModelName:        "local-model",
				CredentialRef:    "lmstudio-local",
			},
		},
		nextPhaseRunID: 200,
	}
	outputRepo := &fakeBodyPhaseOutputRepository{}
	foundation := fakeBodyPhaseFoundationDataRepository{}
	source := fakeBodyPhaseTranslationSourceRepository{
		xedit: repository.XEditExtractedData{ID: 10},
		records: []repository.TranslationRecord{
			{ID: 11, XEditExtractedDataID: 10, RecordType: "INFO", FormID: "000AAA", EditorID: "LineOne"},
		},
		fields: map[int64][]repository.TranslationField{
			11: {
				{ID: 101, TranslationRecordID: 11, SubrecordType: "FULL", SourceText: "Hello there", FieldOrder: 1},
				{ID: 102, TranslationRecordID: 11, SubrecordType: "DESC", SourceText: "Safe travels", FieldOrder: 2},
			},
		},
	}
	transactor := fakeBodyPhaseTransactor{jobRepo: jobRepo, outputRepo: outputRepo}
	service := NewBodyTranslationPhaseService(jobRepo, foundation, source, outputRepo, transactor).
		WithBodyTranslationProvider(provider)
	service.WithBodyTranslationProviderSettings(fakePhaseProviderSettingsConsumer{
		resolveFunc: func(_ context.Context, input ProviderSettingsResolveInput) (ProviderSettingsResolveResult, error) {
			endpoint := "http://localhost:1234"
			return ProviderSettingsResolveResult{
				ConsumerID:      input.ConsumerID,
				ProviderID:      input.Selection.ProviderID,
				Model:           input.Selection.Model,
				ExecutionMethod: input.Selection.ExecutionMethod,
				UseBatchAPI:     input.Selection.UseBatchAPI,
				Endpoint:        &endpoint,
				CredentialState: providerSettingsCredentialStateNotRequired,
			}, nil
		},
	})
	service.now = func() time.Time { return time.Date(2026, 5, 3, 12, 0, 0, 0, time.UTC) }
	outputRepo.nowFunc = service.now
	return service, jobRepo, outputRepo
}

func TestBodyTranslationPhaseServiceStartPhaseReResolvesProviderSettingsBeforeExecution(t *testing.T) {
	capturedRequests := make([]BodyTranslationProviderRequest, 0, 1)
	provider := fakeBodyPhaseProvider{
		translateFunc: func(_ context.Context, request BodyTranslationProviderRequest) BodyTranslationProviderResult {
			capturedRequests = append(capturedRequests, request)
			return BodyTranslationProviderResult{
				RequestUnitID:       request.RequestUnitID,
				FieldCorrelationKey: request.FieldCorrelationKey,
				RecordType:          request.RecordType,
				FieldType:           request.FieldType,
				TranslatedCandidate: &BodyTranslationTranslatedCandidate{
					RequestUnitID:       request.RequestUnitID,
					FieldCorrelationKey: request.FieldCorrelationKey,
					RecordType:          request.RecordType,
					FieldType:           request.FieldType,
					TranslatedText:      request.SourceText + "_ja",
				},
				ProtectionValidationTarget: &BodyTranslationProtectionValidationTarget{
					RequestUnitID:       request.RequestUnitID,
					FieldCorrelationKey: request.FieldCorrelationKey,
					TranslatedText:      request.SourceText + "_ja",
				},
			}
		},
	}
	service, jobRepo, _ := newBodyTranslationPhaseServiceForTest(provider)
	jobRepo.runsByPhaseType[bodyTranslationPersonaPhaseType] = repository.JobPhaseRun{
		ID:               100,
		TranslationJobID: 1,
		PhaseType:        bodyTranslationPersonaPhaseType,
		State:            bodyTranslationPhaseStateCompleted,
		ExecutionMode:    "sync",
		AIProvider:       BodyTranslationProviderXAI,
		ModelName:        "grok-4",
		CredentialRef:    "stale-ref",
	}
	service.WithBodyTranslationProviderSettings(fakePhaseProviderSettingsConsumer{
		resolveFunc: func(_ context.Context, input ProviderSettingsResolveInput) (ProviderSettingsResolveResult, error) {
			return ProviderSettingsResolveResult{
				ConsumerID:            input.ConsumerID,
				ProviderID:            input.Selection.ProviderID,
				Model:                 input.Selection.Model,
				ExecutionMethod:       input.Selection.ExecutionMethod,
				Endpoint:              stringPointer("https://api.x.ai/v1"),
				CredentialReferenceID: stringPointer("provider-settings:xai"),
				CredentialState:       providerSettingsCredentialStateConfigured,
			}, nil
		},
	})

	_, err := service.StartPhase(context.Background(), 1)
	if err != nil {
		t.Fatalf("expected body phase start success, got %v", err)
	}
	if len(capturedRequests) == 0 {
		t.Fatal("expected provider request capture")
	}
	if capturedRequests[0].CredentialRef != "provider-settings:xai" {
		t.Fatalf("expected re-resolved credential ref, got %#v", capturedRequests[0])
	}
	if capturedRequests[0].EndpointSummary == nil || *capturedRequests[0].EndpointSummary != "https://api.x.ai/v1" {
		t.Fatalf("expected re-resolved endpoint summary, got %#v", capturedRequests[0])
	}
}

func TestBodyTranslationPhaseServiceStartPhaseCompletesWithoutPersonaSnapshot(t *testing.T) {
	provider := fakeBodyPhaseProvider{
		translateFunc: func(_ context.Context, request BodyTranslationProviderRequest) BodyTranslationProviderResult {
			return BodyTranslationProviderResult{
				RequestUnitID:       request.RequestUnitID,
				FieldCorrelationKey: request.FieldCorrelationKey,
				RecordType:          request.RecordType,
				FieldType:           request.FieldType,
				TranslatedCandidate: &BodyTranslationTranslatedCandidate{
					RequestUnitID:       request.RequestUnitID,
					FieldCorrelationKey: request.FieldCorrelationKey,
					RecordType:          request.RecordType,
					FieldType:           request.FieldType,
					TranslatedText:      request.SourceText + "_ja",
				},
				ProtectionValidationTarget: &BodyTranslationProtectionValidationTarget{
					RequestUnitID:       request.RequestUnitID,
					FieldCorrelationKey: request.FieldCorrelationKey,
					TranslatedText:      request.SourceText + "_ja",
				},
			}
		},
	}
	service, jobRepo, outputRepo := newBodyTranslationPhaseServiceForTest(provider)

	result, err := service.StartPhase(context.Background(), 1)
	if err != nil {
		t.Fatalf("expected body phase start success, got %v", err)
	}
	if result.ErrorSummary != nil {
		t.Fatalf("expected no error summary, got %#v", result.ErrorSummary)
	}
	assertBodyPhaseStartPhaseState(t, result)
	assertBodyPhaseCompletedOutput(t, result, outputRepo.fields)

	bodyRun, ok := jobRepo.runsByPhaseType[bodyTranslationPhaseType]
	if !ok {
		t.Fatal("expected body translation phase run to be created")
	}
	if bodyRun.PersonaDigest == "" || bodyRun.InputSnapshotDigest == "" {
		t.Fatalf("expected persisted snapshot digests, got %#v", bodyRun)
	}
	assertBodyPhaseCompletedJobState(t, result, jobRepo.job)
}

func assertBodyPhaseStartPhaseState(t *testing.T, result BodyTranslationPhaseCommandReadModel) {
	t.Helper()
	if result.PhaseState != bodyTranslationPhaseStateRunning && result.PhaseState != bodyTranslationPhaseStateCompleted {
		t.Fatalf("expected running or completed phase state, got %#v", result)
	}
}

func assertBodyPhaseCompletedOutput(
	t *testing.T,
	result BodyTranslationPhaseCommandReadModel,
	fields []repository.JobTranslationField,
) {
	t.Helper()
	if result.PhaseState != bodyTranslationPhaseStateCompleted {
		return
	}
	if !result.OutputReadiness.Ready {
		t.Fatalf("expected output readiness ready, got %#v", result.OutputReadiness)
	}
	if result.Progress.TargetCount != 2 || result.Progress.TranslatedCount != 2 {
		t.Fatalf("expected two translated targets, got %#v", result.Progress)
	}
	if len(fields) != 2 {
		t.Fatalf("expected two persisted output fields, got %#v", fields)
	}
	assertBodyPhaseOutputFieldsWithoutPersonaSnapshot(t, fields)
}

func assertBodyPhaseCompletedJobState(
	t *testing.T,
	result BodyTranslationPhaseCommandReadModel,
	job repository.TranslationJob,
) {
	t.Helper()
	if result.PhaseState != bodyTranslationPhaseStateCompleted {
		return
	}
	if job.State != bodyTranslationJobStateCompleted {
		t.Fatalf("expected completed job state, got %#v", job)
	}
}

func assertBodyPhaseOutputFieldsWithoutPersonaSnapshot(
	t *testing.T,
	fields []repository.JobTranslationField,
) {
	t.Helper()
	for _, field := range fields {
		if field.AppliedPersonaID != nil {
			t.Fatalf("expected nil applied persona id for empty persona snapshot, got %#v", fields)
		}
		if field.OutputStatus != bodyTranslationOutputStatusReady {
			t.Fatalf("expected ready output status, got %#v", fields)
		}
	}
}
