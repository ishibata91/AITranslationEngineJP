package service

import (
	"context"
	"errors"
	"strings"
	"testing"
	"time"

	"aitranslationenginejp/internal/repository"
)

type fakePersonaPhaseTransactor struct{}

func (fakePersonaPhaseTransactor) WithTransaction(ctx context.Context, fn func(context.Context) error) error {
	return fn(ctx)
}

type fakePersonaPhaseJobLifecycleRepository struct {
	job              repository.TranslationJob
	termRun          repository.JobPhaseRun
	personaRun       *repository.JobPhaseRun
	phaseRunPersonas []repository.PhaseRunPersona
	updateDrafts     []repository.JobPhaseRunUpdateDraft
	updateStates     []string
	updateErr        error
}

func (fake *fakePersonaPhaseJobLifecycleRepository) GetTranslationJobByID(context.Context, int64) (repository.TranslationJob, error) {
	return fake.job, nil
}

func (fake *fakePersonaPhaseJobLifecycleRepository) UpdateTranslationJob(_ context.Context, _ int64, _ repository.TranslationJobUpdateDraft) (repository.TranslationJob, error) {
	return fake.job, nil
}

func (fake *fakePersonaPhaseJobLifecycleRepository) CreateJobPhaseRun(_ context.Context, draft repository.JobPhaseRunDraft) (repository.JobPhaseRun, error) {
	run := repository.JobPhaseRun{
		ID:               99,
		TranslationJobID: fake.job.ID,
		PhaseType:        draft.PhaseType,
		State:            draft.State,
		AIProvider:       draft.AIProvider,
		ModelName:        draft.ModelName,
		ExecutionMode:    draft.ExecutionMode,
		CredentialRef:    draft.CredentialRef,
		InstructionKind:  draft.InstructionKind,
	}
	fake.personaRun = &run
	return run, nil
}

func (fake *fakePersonaPhaseJobLifecycleRepository) FindJobPhaseRun(_ context.Context, _ int64, phaseType string) (repository.JobPhaseRun, error) {
	if phaseType != personaGenerationPhaseType || fake.personaRun == nil {
		return repository.JobPhaseRun{}, repository.ErrNotFound
	}
	return *fake.personaRun, nil
}

func (fake *fakePersonaPhaseJobLifecycleRepository) UpdateJobPhaseRun(_ context.Context, id int64, draft repository.JobPhaseRunUpdateDraft) (repository.JobPhaseRun, error) {
	if fake.personaRun == nil || fake.personaRun.ID != id {
		return repository.JobPhaseRun{}, repository.ErrNotFound
	}
	if fake.updateErr != nil {
		return repository.JobPhaseRun{}, fake.updateErr
	}
	fake.updateDrafts = append(fake.updateDrafts, draft)
	fake.personaRun.State = draft.State
	fake.personaRun.ProgressPercent = draft.ProgressPercent
	fake.personaRun.LatestError = draft.LatestError
	fake.personaRun.StartedAt = draft.StartedAt
	fake.personaRun.FinishedAt = draft.FinishedAt
	return *fake.personaRun, nil
}
func (fake *fakePersonaPhaseJobLifecycleRepository) UpdateJobPhaseRunWhenState(
	ctx context.Context,
	id int64,
	expectedState string,
	draft repository.JobPhaseRunUpdateDraft,
) (repository.JobPhaseRun, error) {
	fake.updateStates = append(fake.updateStates, expectedState)
	return fake.UpdateJobPhaseRun(ctx, id, draft)
}

func (fake *fakePersonaPhaseJobLifecycleRepository) ListJobPhaseRunsByJobID(context.Context, int64) ([]repository.JobPhaseRun, error) {
	phases := []repository.JobPhaseRun{fake.termRun}
	if fake.personaRun != nil {
		phases = append(phases, *fake.personaRun)
	}
	return phases, nil
}

func (fake *fakePersonaPhaseJobLifecycleRepository) CreatePhaseRunPersona(_ context.Context, draft repository.PhaseRunPersonaDraft) (repository.PhaseRunPersona, error) {
	link := repository.PhaseRunPersona{ID: int64(len(fake.phaseRunPersonas) + 1), PhaseRunID: draft.PhaseRunID, PersonaID: draft.PersonaID, Role: draft.Role}
	fake.phaseRunPersonas = append(fake.phaseRunPersonas, link)
	return link, nil
}

func (fake *fakePersonaPhaseJobLifecycleRepository) ListPhaseRunPersonasByPhaseRunID(_ context.Context, phaseRunID int64) ([]repository.PhaseRunPersona, error) {
	result := make([]repository.PhaseRunPersona, 0)
	for _, link := range fake.phaseRunPersonas {
		if link.PhaseRunID == phaseRunID {
			result = append(result, link)
		}
	}
	return result, nil
}

type fakePersonaPhaseFoundationDataRepository struct {
	personaByID                map[int64]repository.Persona
	personaFieldEvidenceByID   map[int64][]repository.PersonaFieldEvidence
	nextPersonaID              int64
	nextPersonaFieldEvidenceID int64
}

func (fake *fakePersonaPhaseFoundationDataRepository) CreatePersona(_ context.Context, draft repository.PersonaDraft) (repository.Persona, error) {
	if fake.personaByID == nil {
		fake.personaByID = map[int64]repository.Persona{}
	}
	fake.nextPersonaID++
	persona := repository.Persona{
		ID:                     fake.nextPersonaID,
		NpcProfileID:           draft.NpcProfileID,
		TranslationJobID:       draft.TranslationJobID,
		PersonaLifecycle:       draft.PersonaLifecycle,
		PersonaScope:           draft.PersonaScope,
		PersonaSource:          draft.PersonaSource,
		PersonaDescription:     draft.PersonaDescription,
		SpeechStyle:            draft.SpeechStyle,
		PersonalitySummary:     draft.PersonalitySummary,
		EvidenceUtteranceCount: draft.EvidenceUtteranceCount,
	}
	fake.personaByID[persona.ID] = persona
	return persona, nil
}
func (fake *fakePersonaPhaseFoundationDataRepository) GetPersonaByID(_ context.Context, id int64) (repository.Persona, error) {
	persona, ok := fake.personaByID[id]
	if !ok {
		return repository.Persona{}, repository.ErrNotFound
	}
	return persona, nil
}
func (fake *fakePersonaPhaseFoundationDataRepository) GetPersonaByNpcProfileID(_ context.Context, npcProfileID int64) (repository.Persona, error) {
	for _, persona := range fake.personaByID {
		if persona.NpcProfileID == npcProfileID {
			return persona, nil
		}
	}
	return repository.Persona{}, repository.ErrNotFound
}
func (fake *fakePersonaPhaseFoundationDataRepository) UpdatePersona(_ context.Context, id int64, draft repository.PersonaUpdateDraft) (repository.Persona, error) {
	persona, ok := fake.personaByID[id]
	if !ok {
		return repository.Persona{}, repository.ErrNotFound
	}
	persona.PersonaLifecycle = draft.PersonaLifecycle
	persona.PersonaScope = draft.PersonaScope
	persona.PersonaSource = draft.PersonaSource
	persona.PersonaDescription = draft.PersonaDescription
	persona.SpeechStyle = draft.SpeechStyle
	persona.PersonalitySummary = draft.PersonalitySummary
	persona.EvidenceUtteranceCount = draft.EvidenceUtteranceCount
	fake.personaByID[id] = persona
	return persona, nil
}
func (fake *fakePersonaPhaseFoundationDataRepository) CreatePersonaFieldEvidence(_ context.Context, draft repository.PersonaFieldEvidenceDraft) (repository.PersonaFieldEvidence, error) {
	fake.nextPersonaFieldEvidenceID++
	evidence := repository.PersonaFieldEvidence{
		ID:                 fake.nextPersonaFieldEvidenceID,
		PersonaID:          draft.PersonaID,
		TranslationFieldID: draft.TranslationFieldID,
		EvidenceRole:       draft.EvidenceRole,
	}
	if fake.personaFieldEvidenceByID == nil {
		fake.personaFieldEvidenceByID = map[int64][]repository.PersonaFieldEvidence{}
	}
	fake.personaFieldEvidenceByID[draft.PersonaID] = append(fake.personaFieldEvidenceByID[draft.PersonaID], evidence)
	return evidence, nil
}
func (fake *fakePersonaPhaseFoundationDataRepository) ListPersonaFieldEvidenceByPersonaID(_ context.Context, personaID int64) ([]repository.PersonaFieldEvidence, error) {
	return append([]repository.PersonaFieldEvidence(nil), fake.personaFieldEvidenceByID[personaID]...), nil
}

type fakePersonaPhaseTranslationSourceRepository struct {
	records        []repository.TranslationRecord
	npcByRecord    map[int64]repository.NpcRecord
	profileByID    map[int64]repository.NpcProfile
	fieldsByRecord map[int64][]repository.TranslationField
}

type fakePersonaGenerationProvider struct{}

type capturingPersonaGenerationProvider struct {
	requests []PersonaGenerationProviderRequest
}

func (fakePersonaGenerationProvider) GeneratePersona(_ context.Context, request PersonaGenerationProviderRequest) PersonaGenerationProviderResult {
	prompt, err := BuildPersonaGenerationPrompt(request)
	if err != nil {
		return PersonaGenerationProviderResult{
			RequestUnitID:    request.RequestUnitID,
			NPCCorrelationID: request.NPCCorrelationID,
			Failure: &PersonaGenerationProviderFailure{
				Kind:       PersonaGenerationProviderErrorKindInvalidProviderResponse,
				Reason:     "invalid test request",
				Retryable:  false,
				IsRedacted: true,
			},
		}
	}
	return PersonaGenerationProviderResult{
		RequestUnitID:    request.RequestUnitID,
		NPCCorrelationID: request.NPCCorrelationID,
		PersonaBody:      "persona body for " + request.NPCCorrelationID,
		AuditSummary: PersonaGenerationProviderAuditSummary{
			CredentialRef: request.CredentialRef,
			Provider:      request.Provider,
			Model:         request.Model,
			ExecutionMode: request.ExecutionMode,
			PromptDigest:  personaGenerationPromptDigest(prompt),
			InputCount:    1,
			OutputCount:   1,
		},
	}
}

func (fakePersonaGenerationProvider) PersonaGenerationProviderRequestsAreTestSafe() bool {
	return true
}

func (provider *capturingPersonaGenerationProvider) GeneratePersona(_ context.Context, request PersonaGenerationProviderRequest) PersonaGenerationProviderResult {
	provider.requests = append(provider.requests, request)
	return fakePersonaGenerationProvider{}.GeneratePersona(context.Background(), request)
}

func (provider *capturingPersonaGenerationProvider) PersonaGenerationProviderRequestsAreTestSafe() bool {
	return true
}

func (fake *fakePersonaPhaseTranslationSourceRepository) GetXEditExtractedDataByID(context.Context, int64) (repository.XEditExtractedData, error) {
	return repository.XEditExtractedData{ID: 1}, nil
}
func (fake *fakePersonaPhaseTranslationSourceRepository) GetTranslationRecordByID(_ context.Context, id int64) (repository.TranslationRecord, error) {
	for _, record := range fake.records {
		if record.ID == id {
			return record, nil
		}
	}
	return repository.TranslationRecord{}, repository.ErrNotFound
}
func (fake *fakePersonaPhaseTranslationSourceRepository) ListTranslationRecordsByXEditID(context.Context, int64) ([]repository.TranslationRecord, error) {
	return append([]repository.TranslationRecord(nil), fake.records...), nil
}
func (fake *fakePersonaPhaseTranslationSourceRepository) GetNpcRecordByTranslationRecordID(_ context.Context, translationRecordID int64) (repository.NpcRecord, error) {
	npc, ok := fake.npcByRecord[translationRecordID]
	if !ok {
		return repository.NpcRecord{}, repository.ErrNotFound
	}
	return npc, nil
}
func (fake *fakePersonaPhaseTranslationSourceRepository) GetNpcProfileByID(_ context.Context, id int64) (repository.NpcProfile, error) {
	profile, ok := fake.profileByID[id]
	if !ok {
		return repository.NpcProfile{}, repository.ErrNotFound
	}
	return profile, nil
}
func (fake *fakePersonaPhaseTranslationSourceRepository) ListTranslationFieldsByTranslationRecordID(_ context.Context, translationRecordID int64) ([]repository.TranslationField, error) {
	return append([]repository.TranslationField(nil), fake.fieldsByRecord[translationRecordID]...), nil
}
func (fake *fakePersonaPhaseTranslationSourceRepository) ListTranslationFieldRecordReferencesByFieldID(context.Context, int64) ([]repository.TranslationFieldRecordReference, error) {
	return nil, nil
}

func newPersonaPhaseServiceForTest(jobState string, termState string, personaRun *repository.JobPhaseRun, sourceRecords []repository.TranslationRecord) (*PersonaGenerationPhaseService, *fakePersonaPhaseJobLifecycleRepository) {
	jobRepo := &fakePersonaPhaseJobLifecycleRepository{
		job:        repository.TranslationJob{ID: 1, State: jobState, XEditExtractedDataID: 1},
		termRun:    repository.JobPhaseRun{ID: 10, TranslationJobID: 1, PhaseType: personaGenerationTermPhaseType, State: termState, AIProvider: "fake", ModelName: "m", ExecutionMode: "single_request", CredentialRef: "cred"},
		personaRun: personaRun,
	}
	npcByRecord := make(map[int64]repository.NpcRecord, len(sourceRecords))
	profileByID := make(map[int64]repository.NpcProfile, len(sourceRecords))
	fieldsByRecord := make(map[int64][]repository.TranslationField, len(sourceRecords))
	for index, record := range sourceRecords {
		profileID := int64(300 + index)
		fieldID := int64(400 + index)
		npcByRecord[record.ID] = repository.NpcRecord{TranslationRecordID: record.ID, NpcProfileID: profileID}
		profileByID[profileID] = repository.NpcProfile{ID: profileID, DisplayName: "Lydia"}
		fieldsByRecord[record.ID] = []repository.TranslationField{{ID: fieldID, SourceText: "line"}}
	}
	source := &fakePersonaPhaseTranslationSourceRepository{
		records:        sourceRecords,
		npcByRecord:    npcByRecord,
		profileByID:    profileByID,
		fieldsByRecord: fieldsByRecord,
	}
	foundation := &fakePersonaPhaseFoundationDataRepository{
		personaByID:              map[int64]repository.Persona{},
		personaFieldEvidenceByID: map[int64][]repository.PersonaFieldEvidence{},
		nextPersonaID:            1000,
	}
	service := NewPersonaGenerationPhaseService(jobRepo, foundation, source, fakePersonaPhaseTransactor{})
	service.now = func() time.Time { return time.Date(2026, 5, 2, 0, 0, 0, 0, time.UTC) }
	return service, jobRepo
}

func preCreatedPersonaGenerationRun() *repository.JobPhaseRun {
	return &repository.JobPhaseRun{
		ID:               88,
		TranslationJobID: 1,
		PhaseType:        personaGenerationPhaseType,
		State:            personaGenerationPhaseStatePending,
		AIProvider:       "fake",
		ModelName:        "m",
		ExecutionMode:    PersonaGenerationExecutionModeSingleRequest,
		CredentialRef:    "cred",
	}
}

func TestPersonaGenerationPhaseServiceStartPhaseReResolvesProviderSettingsBeforeExecution(t *testing.T) {
	sourceRecords := []repository.TranslationRecord{{ID: 21, RecordType: "NPC_", EditorID: "Lydia", FormID: "000ABC"}}
	service, repo := newPersonaPhaseServiceForTest(personaGenerationJobStateRunning, personaGenerationPhaseStateCompleted, preCreatedPersonaGenerationRun(), sourceRecords)
	repo.termRun = repository.JobPhaseRun{
		ID:               10,
		TranslationJobID: 1,
		PhaseType:        personaGenerationTermPhaseType,
		State:            personaGenerationPhaseStateCompleted,
		AIProvider:       PersonaGenerationProviderXAI,
		ModelName:        "grok-4",
		ExecutionMode:    PersonaGenerationExecutionModeSingleRequest,
		CredentialRef:    "stale-ref",
	}
	capturingProvider := &capturingPersonaGenerationProvider{}
	service.WithPersonaGenerationProvider(capturingProvider)
	service.WithPersonaGenerationProviderSettings(fakePhaseProviderSettingsConsumer{
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
		t.Fatalf("expected persona phase start success, got %v", err)
	}
	if len(capturingProvider.requests) == 0 {
		t.Fatal("expected provider request capture")
	}
	request := capturingProvider.requests[0]
	if request.CredentialRef != "provider-settings:xai" {
		t.Fatalf("expected re-resolved credential ref, got %#v", request)
	}
	if request.EndpointSummary == nil || *request.EndpointSummary != "https://api.x.ai/v1" {
		t.Fatalf("expected re-resolved endpoint summary, got %#v", request)
	}
}

func TestPersonaGenerationBuildProviderRequestAllowsEmptyCredentialRefForLMStudio(t *testing.T) {
	service := NewPersonaGenerationPhaseService(
		&fakePersonaPhaseJobLifecycleRepository{},
		&fakePersonaPhaseFoundationDataRepository{},
		&fakePersonaPhaseTranslationSourceRepository{},
		fakePersonaPhaseTransactor{},
	)
	race := "Nord"
	sex := "F"
	className := "Warrior"
	request, err := service.buildProviderRequest(
		repository.JobPhaseRun{
			ID:            77,
			AIProvider:    "lm_studio",
			ModelName:     "lmstudio-community",
			ExecutionMode: "sync",
			CredentialRef: "",
		},
		personaGenerationTarget{
			record:  repository.TranslationRecord{EditorID: "NPC_01", FormID: "0001"},
			npc:     repository.NpcRecord{VoiceType: "FemaleEvenToned", Race: &race, Sex: &sex, NpcClass: &className},
			profile: repository.NpcProfile{ID: 42, DisplayName: "Lydia"},
			fields:  []repository.TranslationField{{SourceText: "I am sworn to carry your burdens."}},
		},
	)
	if err != nil {
		t.Fatalf("expected lm studio provider request without credentialRef: %v", err)
	}
	if request.Provider != "lm_studio" || request.CredentialRef != "" {
		t.Fatalf("expected lm studio request to keep empty credentialRef, got %#v", request)
	}
}

func TestPersonaGenerationBuildProviderRequestRejectsEmptyCredentialRefForCredentialRequiredProvider(t *testing.T) {
	service := NewPersonaGenerationPhaseService(
		&fakePersonaPhaseJobLifecycleRepository{},
		&fakePersonaPhaseFoundationDataRepository{},
		&fakePersonaPhaseTranslationSourceRepository{},
		fakePersonaPhaseTransactor{},
	)
	_, err := service.buildProviderRequest(
		repository.JobPhaseRun{
			ID:            77,
			AIProvider:    "gemini",
			ModelName:     "gemini-2.5-pro",
			ExecutionMode: "sync",
			CredentialRef: "",
		},
		personaGenerationTarget{
			record:  repository.TranslationRecord{EditorID: "NPC_01", FormID: "0001"},
			npc:     repository.NpcRecord{VoiceType: "FemaleEvenToned"},
			profile: repository.NpcProfile{ID: 42, DisplayName: "Lydia"},
			fields:  []repository.TranslationField{{SourceText: "I am sworn to carry your burdens."}},
		},
	)
	if err == nil {
		t.Fatal("expected credential-required provider to reject empty credentialRef")
	}
}

func TestPersonaGenerationPhaseService_StartPhaseRejectsTerminalJob(t *testing.T) {
	service, _ := newPersonaPhaseServiceForTest(personaGenerationJobStateCompleted, personaGenerationPhaseStateCompleted, nil, nil)

	result, err := service.StartPhase(context.Background(), 1)
	if !errors.Is(err, errPersonaGenerationPhaseExecutionRejected) {
		t.Fatalf("expected execution rejected error, got %v", err)
	}
	if result.ErrorSummary == nil || result.ErrorSummary.ErrorKind != personaGenerationErrorKindTerminalJob {
		t.Fatalf("expected terminal_job rejection, got %#v", result.ErrorSummary)
	}
}

func TestPersonaGenerationPhaseService_StartPhaseRejectsWhenTermIncomplete(t *testing.T) {
	service, _ := newPersonaPhaseServiceForTest(personaGenerationJobStateRunning, personaGenerationPhaseStateRunning, nil, nil)

	result, err := service.StartPhase(context.Background(), 1)
	if !errors.Is(err, errPersonaGenerationPhaseExecutionRejected) {
		t.Fatalf("expected execution rejected error, got %v", err)
	}
	if result.ErrorSummary == nil || result.ErrorSummary.ErrorKind != personaGenerationErrorKindTermIncomplete {
		t.Fatalf("expected term_phase_incomplete, got %#v", result.ErrorSummary)
	}
}

func TestPersonaGenerationPhaseService_StartPhaseRejectedCommandPromptDigestDiffersFromTargetSnapshotDigest(t *testing.T) {
	sourceRecords := []repository.TranslationRecord{{ID: 100, RecordType: "NPC_", EditorID: "NPC_A", FormID: "0001"}}
	service, _ := newPersonaPhaseServiceForTest(personaGenerationJobStateRunning, personaGenerationPhaseStateRunning, nil, sourceRecords)

	result, err := service.StartPhase(context.Background(), 1)
	if !errors.Is(err, errPersonaGenerationPhaseExecutionRejected) {
		t.Fatalf("expected execution rejected error, got %v", err)
	}
	if result.Execution.PromptDigest != "" {
		t.Fatalf("expected empty prompt digest for rejected command without execution run, got %q", result.Execution.PromptDigest)
	}
	if result.TargetSummary.TargetSnapshotDigest == result.Execution.PromptDigest {
		t.Fatalf("expected rejected command prompt digest to differ from target snapshot digest, got summary=%#v execution=%#v", result.TargetSummary, result.Execution)
	}
}

func TestPersonaGenerationPhaseService_StartPhaseRejectsWhenActiveRunExists(t *testing.T) {
	run := &repository.JobPhaseRun{ID: 20, TranslationJobID: 1, PhaseType: personaGenerationPhaseType, State: personaGenerationPhaseStateRunning}
	service, _ := newPersonaPhaseServiceForTest(personaGenerationJobStateRunning, personaGenerationPhaseStateCompleted, run, nil)

	result, err := service.StartPhase(context.Background(), 1)
	if !errors.Is(err, errPersonaGenerationPhaseExecutionRejected) {
		t.Fatalf("expected execution rejected error, got %v", err)
	}
	if result.ErrorSummary == nil || result.ErrorSummary.ErrorKind != personaGenerationErrorKindActivePhase {
		t.Fatalf("expected active_phase_exists, got %#v", result.ErrorSummary)
	}
}

func TestPersonaGenerationPhaseService_StartPhaseTargetZeroCompletesImmediately(t *testing.T) {
	service, repo := newPersonaPhaseServiceForTest(personaGenerationJobStateRunning, personaGenerationPhaseStateCompleted, preCreatedPersonaGenerationRun(), nil)

	result, err := service.StartPhase(context.Background(), 1)
	if err != nil {
		t.Fatalf("expected start success, got %v", err)
	}
	if result.PhaseState != personaGenerationPhaseStateEmptyCompleted {
		t.Fatalf("expected empty_completed state, got %s", result.PhaseState)
	}
	if result.Progress.Percent != 100 || result.TargetSummary.TargetCount != 0 {
		t.Fatalf("unexpected progress/target: %#v %#v", result.Progress, result.TargetSummary)
	}
	if len(repo.updateDrafts) == 0 {
		t.Fatal("expected phase run update")
	}
}

func TestPersonaGenerationPhaseServiceStartPhasePromotesPreCreatedPendingRunToRunning(t *testing.T) {
	sourceRecords := []repository.TranslationRecord{{ID: 21, RecordType: "NPC_", EditorID: "Lydia", FormID: "000ABC"}}
	preCreatedRun := preCreatedPersonaGenerationRun()
	service, repo := newPersonaPhaseServiceForTest(personaGenerationJobStateRunning, personaGenerationPhaseStateCompleted, preCreatedRun, sourceRecords)
	service.WithPersonaGenerationProvider(fakePersonaGenerationProvider{})

	result, err := service.StartPhase(context.Background(), 1)
	if err != nil {
		t.Fatalf("expected start success with pre-created run: %v", err)
	}
	if result.PhaseRunID == nil || *result.PhaseRunID != preCreatedRun.ID {
		t.Fatalf("expected existing run id reuse, got %#v", result.PhaseRunID)
	}
	if len(repo.updateDrafts) == 0 {
		t.Fatalf("expected update drafts for pre-created run, got %#v", repo.updateDrafts)
	}
	if repo.updateDrafts[0].State != personaGenerationPhaseStateRunning {
		t.Fatalf("expected first update to running, got %#v", repo.updateDrafts[0])
	}
	if repo.updateStates[0] != personaGenerationPhaseStatePending {
		t.Fatalf("expected optimistic expected state pending, got %#v", repo.updateStates)
	}
}

func TestPersonaGenerationPhaseServiceStartPhaseCreatesNonPendingRunWithoutPrecreatedPhaseRun(t *testing.T) {
	sourceRecords := []repository.TranslationRecord{{ID: 21, RecordType: "NPC_", EditorID: "Lydia", FormID: "000ABC"}}
	service, repo := newPersonaPhaseServiceForTest(personaGenerationJobStateRunning, personaGenerationPhaseStateCompleted, nil, sourceRecords)
	capturingProvider := &capturingPersonaGenerationProvider{}
	service.WithPersonaGenerationProvider(capturingProvider)

	result, err := service.StartPhase(context.Background(), 1)
	if err != nil {
		t.Fatalf("expected start success without pre-created phase run: %v", err)
	}
	if result.PhaseRunID == nil || *result.PhaseRunID == 0 {
		t.Fatalf("expected created persona phase run id, got %#v", result.PhaseRunID)
	}
	if repo.personaRun == nil {
		t.Fatal("expected created persona generation phase run")
	}
	if repo.personaRun.State == personaGenerationPhaseStatePending {
		t.Fatalf("expected created run not to be pending, got %#v", repo.personaRun)
	}
	if len(capturingProvider.requests) != 1 {
		t.Fatalf("expected one provider request after phase run creation, got %#v", capturingProvider.requests)
	}
}

func TestPersonaGenerationPhaseServiceStartPhaseStopsBeforeProviderWhenPhasePromotionConflicts(t *testing.T) {
	sourceRecords := []repository.TranslationRecord{{ID: 21, RecordType: "NPC_", EditorID: "Lydia", FormID: "000ABC"}}
	preCreatedRun := &repository.JobPhaseRun{
		ID:               89,
		TranslationJobID: 1,
		PhaseType:        personaGenerationPhaseType,
		State:            personaGenerationPhaseStatePending,
		AIProvider:       "fake",
		ModelName:        "m",
		ExecutionMode:    PersonaGenerationExecutionModeSingleRequest,
		CredentialRef:    "cred",
	}
	provider := &capturingPersonaGenerationProvider{}
	service, repo := newPersonaPhaseServiceForTest(personaGenerationJobStateRunning, personaGenerationPhaseStateCompleted, preCreatedRun, sourceRecords)
	service.WithPersonaGenerationProvider(provider)
	repo.updateErr = repository.ErrConflict

	_, err := service.StartPhase(context.Background(), 1)
	if !errors.Is(err, repository.ErrConflict) {
		t.Fatalf("expected conflict error, got %v", err)
	}
	if len(provider.requests) != 0 {
		t.Fatalf("expected no provider request on phase promotion conflict, got %#v", provider.requests)
	}
}

func TestPersonaGenerationPhaseService_CommandMutationsFollowCommonOperationPolicy(t *testing.T) {
	run := &repository.JobPhaseRun{ID: 20, TranslationJobID: 1, PhaseType: personaGenerationPhaseType, State: personaGenerationPhaseStateRunning, ProgressPercent: 25}
	sourceRecords := []repository.TranslationRecord{
		{ID: 100, RecordType: "NPC_"},
		{ID: 101, RecordType: "NPC_"},
	}
	service, repo := newPersonaPhaseServiceForTest(personaGenerationJobStateRunning, personaGenerationPhaseStateCompleted, run, sourceRecords)

	paused, err := service.PausePhase(context.Background(), 1, 20)
	if err != nil || paused.PhaseState != personaGenerationPhaseStatePaused {
		t.Fatalf("pause failed: state=%s err=%v", paused.PhaseState, err)
	}
	resumed, err := service.ResumePhase(context.Background(), 1, 20)
	if err == nil {
		t.Fatalf("resume should return execution error for missing provider: %#v", resumed)
	}
	retried, err := service.RetryPhase(context.Background(), 1, 20)
	if err == nil {
		t.Fatalf("retry should return execution error for missing provider: %#v", retried)
	}
	canceled, err := service.CancelPhase(context.Background(), 1, 20)
	if err != nil {
		t.Fatalf("cancel should return rejected response without transport error: %v", err)
	}
	if canceled.PhaseState != personaGenerationPhaseStateRejected || canceled.ErrorSummary == nil || canceled.ErrorSummary.ErrorKind != personaGenerationErrorKindTermIncomplete {
		t.Fatalf("expected cancel rejection from non-paused state, got %#v", canceled)
	}
	_ = repo
}

func TestPersonaGenerationPhaseService_ReadBodyReadinessBranches(t *testing.T) {
	run := &repository.JobPhaseRun{ID: 20, TranslationJobID: 1, PhaseType: personaGenerationPhaseType, State: personaGenerationPhaseStateRunning, ProgressPercent: 25}
	sourceRecords := []repository.TranslationRecord{
		{ID: 100, RecordType: "NPC_"},
		{ID: 101, RecordType: "NPC_"},
	}
	service, repo := newPersonaPhaseServiceForTest(personaGenerationJobStateRunning, personaGenerationPhaseStateCompleted, run, sourceRecords)
	repo.personaRun.State = personaGenerationPhaseStateCompleted
	repo.phaseRunPersonas = nil
	missing, missingErr := service.ReadBodyReadiness(context.Background(), 1)
	if missingErr != nil {
		t.Fatalf("expected snapshot missing response without error, got %v", missingErr)
	}
	if missing.InputSummary.MissingCount == 0 {
		t.Fatalf("expected snapshot missing count > 0, got %#v", missing)
	}

	repo.phaseRunPersonas = []repository.PhaseRunPersona{{ID: 1, PhaseRunID: 20, PersonaID: 777, Role: personaGenerationPhaseLinkRoleApplied}}
	jobID := int64(1)
	service.foundationDataRepository = &fakePersonaPhaseFoundationDataRepository{
		personaByID: map[int64]repository.Persona{
			777: {ID: 777, NpcProfileID: 300, TranslationJobID: &jobID},
		},
		personaFieldEvidenceByID: map[int64][]repository.PersonaFieldEvidence{},
		nextPersonaID:            1000,
	}
	ready, readyErr := service.ReadBodyReadiness(context.Background(), 1)
	if readyErr != nil {
		t.Fatalf("expected ready response, got %v", readyErr)
	}
	if ready.InputSummary.MissingCount != 1 {
		t.Fatalf("expected missing count 1 for one uncovered persona, got %#v", ready)
	}

	repo.phaseRunPersonas = append(repo.phaseRunPersonas, repository.PhaseRunPersona{ID: 2, PhaseRunID: 20, PersonaID: 778, Role: personaGenerationPhaseLinkRoleApplied})
	service.foundationDataRepository = &fakePersonaPhaseFoundationDataRepository{
		personaByID: map[int64]repository.Persona{
			777: {ID: 777, NpcProfileID: 300, TranslationJobID: &jobID},
			778: {ID: 778, NpcProfileID: 301, TranslationJobID: &jobID},
		},
		personaFieldEvidenceByID: map[int64][]repository.PersonaFieldEvidence{},
		nextPersonaID:            1000,
	}
	ready, readyErr = service.ReadBodyReadiness(context.Background(), 1)
	if readyErr != nil {
		t.Fatalf("expected ready response after complete snapshot, got %v", readyErr)
	}
	if ready.InputSummary.PersonaCount != 2 || ready.InputSummary.MissingCount != 0 {
		t.Fatalf("expected persona count 2 and missing count 0 when all personas linked, got %#v", ready)
	}
}

func TestPersonaGenerationPhaseService_ReadBodyReadinessCountsDistinctPersonaCoverage(t *testing.T) {
	run := &repository.JobPhaseRun{ID: 20, TranslationJobID: 1, PhaseType: personaGenerationPhaseType, State: personaGenerationPhaseStateCompleted}
	sourceRecords := []repository.TranslationRecord{
		{ID: 100, RecordType: "NPC_"},
		{ID: 101, RecordType: "NPC_"},
	}
	service, repo := newPersonaPhaseServiceForTest(personaGenerationJobStateRunning, personaGenerationPhaseStateCompleted, run, sourceRecords)
	jobID := int64(1)
	repo.phaseRunPersonas = []repository.PhaseRunPersona{
		{ID: 1, PhaseRunID: 20, PersonaID: 777, Role: personaGenerationPhaseLinkRoleApplied},
		{ID: 2, PhaseRunID: 20, PersonaID: 777, Role: personaGenerationPhaseLinkRoleApplied},
	}
	service.foundationDataRepository = &fakePersonaPhaseFoundationDataRepository{
		personaByID: map[int64]repository.Persona{
			777: {ID: 777, NpcProfileID: 300, TranslationJobID: &jobID},
		},
		personaFieldEvidenceByID: map[int64][]repository.PersonaFieldEvidence{},
		nextPersonaID:            1000,
	}

	readiness, err := service.ReadBodyReadiness(context.Background(), 1)
	if err != nil {
		t.Fatalf("expected readiness success, got %v", err)
	}
	if readiness.InputSummary.MissingCount != 1 || readiness.InputSummary.PersonaCount != 1 {
		t.Fatalf("expected distinct persona coverage counts, got %#v", readiness.InputSummary)
	}
}

func TestPersonaGenerationPhaseService_RetryRejectsCompletedRunWithoutMutation(t *testing.T) {
	run := &repository.JobPhaseRun{ID: 20, TranslationJobID: 1, PhaseType: personaGenerationPhaseType, State: personaGenerationPhaseStateCompleted, LatestError: personaGenerationErrorKindProviderFailure}
	service, repo := newPersonaPhaseServiceForTest(personaGenerationJobStateRunning, personaGenerationPhaseStateCompleted, run, []repository.TranslationRecord{{ID: 100, RecordType: "NPC_"}})

	result, err := service.RetryPhase(context.Background(), 1, 20)
	if err != nil {
		t.Fatalf("expected rejected result without error, got %v", err)
	}
	if result.PhaseState != personaGenerationPhaseStateRejected || result.ErrorSummary == nil || result.ErrorSummary.ErrorKind != personaGenerationErrorKindTermIncomplete {
		t.Fatalf("unexpected retry rejection: %#v", result)
	}
	if len(repo.updateDrafts) != 0 || repo.personaRun.State != personaGenerationPhaseStateCompleted {
		t.Fatalf("expected no mutation for completed run, got drafts=%#v state=%s", repo.updateDrafts, repo.personaRun.State)
	}
}

func TestPersonaGenerationPhaseService_RetryUsesRecoverableFailedWithoutLatestErrorDependency(t *testing.T) {
	run := &repository.JobPhaseRun{ID: 20, TranslationJobID: 1, PhaseType: personaGenerationPhaseType, State: personaGenerationPhaseStateRecoverableFail, LatestError: personaGenerationErrorKindSaveFailed}
	service, repo := newPersonaPhaseServiceForTest(personaGenerationJobStateRunning, personaGenerationPhaseStateCompleted, run, []repository.TranslationRecord{{ID: 100, RecordType: "NPC_"}})

	_, err := service.RetryPhase(context.Background(), 1, 20)
	if err == nil || !strings.Contains(err.Error(), "persona generation input missing") {
		t.Fatalf("expected retry execution error caused by missing input, got %v", err)
	}
	if len(repo.updateDrafts) < 2 || repo.updateDrafts[0].State != personaGenerationPhaseStateRunning {
		t.Fatalf("expected retry execution path updates, got %#v", repo.updateDrafts)
	}
	if repo.personaRun.State != personaGenerationPhaseStateRecoverableFail || repo.personaRun.LatestError != personaGenerationErrorKindInputMissing {
		t.Fatalf("expected recoverable_failed with input_missing after retry execution error, got %#v", repo.personaRun)
	}
}

func TestPersonaGenerationPhaseService_RetryRejectsTerminalJobWithoutMutation(t *testing.T) {
	run := &repository.JobPhaseRun{ID: 20, TranslationJobID: 1, PhaseType: personaGenerationPhaseType, State: personaGenerationPhaseStateRecoverableFail, LatestError: personaGenerationErrorKindProviderFailure}
	service, repo := newPersonaPhaseServiceForTest(personaGenerationJobStateCompleted, personaGenerationPhaseStateCompleted, run, []repository.TranslationRecord{{ID: 100, RecordType: "NPC_"}})

	result, err := service.RetryPhase(context.Background(), 1, 20)
	if err != nil {
		t.Fatalf("expected rejected result without error, got %v", err)
	}
	if result.PhaseState != personaGenerationPhaseStateRejected || result.ErrorSummary == nil || result.ErrorSummary.ErrorKind != personaGenerationErrorKindTerminalJob {
		t.Fatalf("unexpected terminal retry rejection: %#v", result)
	}
	if len(repo.updateDrafts) != 0 || repo.personaRun.State != personaGenerationPhaseStateRecoverableFail {
		t.Fatalf("expected no mutation for terminal job retry, got drafts=%#v state=%s", repo.updateDrafts, repo.personaRun.State)
	}
}

func TestPersonaGenerationPhaseService_RetryRejectsSnapshotDriftWithoutMutation(t *testing.T) {
	run := &repository.JobPhaseRun{
		ID:                  20,
		TranslationJobID:    1,
		PhaseType:           personaGenerationPhaseType,
		State:               personaGenerationPhaseStateRecoverableFail,
		LatestError:         personaGenerationErrorKindProviderFailure,
		LatestExternalRunID: "sha256:saved-snapshot",
	}
	sourceRecords := []repository.TranslationRecord{{ID: 100, RecordType: "NPC_"}}
	service, repo := newPersonaPhaseServiceForTest(personaGenerationJobStateRunning, personaGenerationPhaseStateCompleted, run, sourceRecords)

	result, err := service.RetryPhase(context.Background(), 1, 20)
	if err != nil {
		t.Fatalf("expected rejected result without error, got %v", err)
	}
	if result.ErrorSummary == nil || result.ErrorSummary.ErrorKind != personaGenerationErrorKindSnapshotMissing {
		t.Fatalf("expected snapshot drift rejection, got %#v", result)
	}
	if len(repo.updateDrafts) != 0 {
		t.Fatalf("expected no mutation on drift, got %#v", repo.updateDrafts)
	}
}

func TestPersonaGenerationPhaseService_CancelRejectsNonCancelableStateAndTerminalJobWithoutMutation(t *testing.T) {
	sourceRecords := []repository.TranslationRecord{{ID: 100, RecordType: "NPC_"}}
	testCases := []struct {
		name        string
		runState    string
		jobState    string
		expectedErr string
	}{
		{
			name:        "completed run",
			runState:    personaGenerationPhaseStateCompleted,
			jobState:    personaGenerationJobStateRunning,
			expectedErr: personaGenerationErrorKindTermIncomplete,
		},
		{
			name:        "terminal job",
			runState:    personaGenerationPhaseStateRunning,
			jobState:    personaGenerationJobStateCanceled,
			expectedErr: personaGenerationErrorKindTerminalJob,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			run := &repository.JobPhaseRun{ID: 20, TranslationJobID: 1, PhaseType: personaGenerationPhaseType, State: tc.runState}
			service, repo := newPersonaPhaseServiceForTest(tc.jobState, personaGenerationPhaseStateCompleted, run, sourceRecords)

			result, err := service.CancelPhase(context.Background(), 1, 20)
			if err != nil {
				t.Fatalf("expected rejected result without error, got %v", err)
			}
			if result.PhaseState != personaGenerationPhaseStateRejected || result.ErrorSummary == nil || result.ErrorSummary.ErrorKind != tc.expectedErr {
				t.Fatalf("unexpected cancel rejection: %#v", result)
			}
			if len(repo.updateDrafts) != 0 {
				t.Fatalf("expected no mutation, got %#v", repo.updateDrafts)
			}
		})
	}
}

func TestPersonaGenerationPhaseService_ReadSummaryAndReadinessBlockSnapshotDrift(t *testing.T) {
	run := &repository.JobPhaseRun{
		ID:                  20,
		TranslationJobID:    1,
		PhaseType:           personaGenerationPhaseType,
		State:               personaGenerationPhaseStateCompleted,
		LatestExternalRunID: "sha256:saved-snapshot",
	}
	sourceRecords := []repository.TranslationRecord{{ID: 100, RecordType: "NPC_"}}
	service, repo := newPersonaPhaseServiceForTest(personaGenerationJobStateRunning, personaGenerationPhaseStateCompleted, run, sourceRecords)
	repo.phaseRunPersonas = []repository.PhaseRunPersona{{ID: 1, PhaseRunID: 20, PersonaID: 777, Role: personaGenerationPhaseLinkRoleApplied}}
	service.foundationDataRepository = &fakePersonaPhaseFoundationDataRepository{
		personaByID: map[int64]repository.Persona{777: {ID: 777, NpcProfileID: 300, TranslationJobID: func() *int64 { jobID := int64(1); return &jobID }()}},
	}

	summary, err := service.ReadSummary(context.Background(), 1)
	if err != nil {
		t.Fatalf("expected summary success, got %v", err)
	}
	if summary.ResultSummary == nil || summary.ResultSummary.BodyReadiness || summary.ErrorSummary == nil || summary.ErrorSummary.ErrorKind != personaGenerationErrorKindSnapshotMissing {
		t.Fatalf("expected summary drift block, got %#v", summary)
	}
	if summary.ActionEnablement.CanRetry || summary.ActionEnablement.CanCancel {
		t.Fatalf("expected actions blocked on drift, got %#v", summary.ActionEnablement)
	}

	readiness, readinessErr := service.ReadBodyReadiness(context.Background(), 1)
	if readinessErr != nil {
		t.Fatalf("expected readiness block without error, got %v", readinessErr)
	}
	if readiness.PhaseState != personaGenerationPhaseStateSnapshotMissing {
		t.Fatalf("expected snapshot_missing phase state on drift, got %#v", readiness)
	}
}

func TestPersonaGenerationPhaseService_FinishPhaseRunSaveFailureAndLoadRunSnapshotDedup(t *testing.T) {
	run := repository.JobPhaseRun{ID: 20, TranslationJobID: 1, PhaseType: personaGenerationPhaseType, State: personaGenerationPhaseStateRunning}
	repo := &fakePersonaPhaseJobLifecycleRepository{job: repository.TranslationJob{ID: 1}, termRun: repository.JobPhaseRun{ID: 10, PhaseType: personaGenerationTermPhaseType, State: personaGenerationPhaseStateCompleted}, personaRun: &run}
	repo.phaseRunPersonas = []repository.PhaseRunPersona{{ID: 1, PhaseRunID: 20, PersonaID: 501}, {ID: 2, PhaseRunID: 20, PersonaID: 501}}
	jobID := int64(1)
	foundation := &fakePersonaPhaseFoundationDataRepository{personaByID: map[int64]repository.Persona{501: {ID: 501, NpcProfileID: 300, TranslationJobID: &jobID}}}
	service := NewPersonaGenerationPhaseService(repo, foundation, &fakePersonaPhaseTranslationSourceRepository{}, fakePersonaPhaseTransactor{})
	service.now = func() time.Time { return time.Date(2026, 5, 2, 0, 0, 0, 0, time.UTC) }

	snapshot, err := service.loadRunSnapshot(context.Background(), &run)
	if err != nil {
		t.Fatalf("expected loadRunSnapshot success, got %v", err)
	}
	if snapshot.totalLinkedCount != 1 || snapshot.generatedCount != 1 {
		t.Fatalf("expected deduplicated run snapshot, got %#v", snapshot)
	}

	updated, finishErr := service.finishPhaseRun(context.Background(), run, personaGenerationTargetSnapshot{targetCount: 2, digest: "sha256:test"}, snapshot, []personaGenerationExecutionFailure{{kind: personaGenerationErrorKindSaveFailed, retryable: false}})
	if !errors.Is(finishErr, errPersonaGenerationPhaseSaveFailed) {
		t.Fatalf("expected save failed error, got %v", finishErr)
	}
	if updated.State != personaGenerationPhaseStateRecoverableFail {
		t.Fatalf("expected recoverable_failed, got %s", updated.State)
	}
}

func TestPersonaGenerationPhaseService_ExecutionPromptDigestUsesAggregatePromptDigest(t *testing.T) {
	sourceRecords := []repository.TranslationRecord{{ID: 100, RecordType: "NPC_", EditorID: "NPC_A", FormID: "0001"}}
	service, repo := newPersonaPhaseServiceForTest(personaGenerationJobStateRunning, personaGenerationPhaseStateCompleted, preCreatedPersonaGenerationRun(), sourceRecords)
	service.provider = fakePersonaGenerationProvider{}

	result, err := service.StartPhase(context.Background(), 1)
	if err != nil {
		t.Fatalf("expected start success, got %v", err)
	}
	if result.TargetSummary.TargetSnapshotDigest == result.Execution.PromptDigest {
		t.Fatalf("expected prompt digest to differ from target snapshot digest, got %#v", result.Execution)
	}

	summary, summaryErr := service.ReadSummary(context.Background(), 1)
	if summaryErr != nil {
		t.Fatalf("expected summary success, got %v", summaryErr)
	}
	if summary.Execution.PromptDigest == "" || summary.Execution.PromptDigest != result.Execution.PromptDigest {
		t.Fatalf("expected summary prompt digest to match command response, got summary=%#v result=%#v", summary.Execution, result.Execution)
	}
	if len(repo.phaseRunPersonas) != 1 {
		t.Fatalf("expected one persisted persona link, got %#v", repo.phaseRunPersonas)
	}
}

func TestPersonaGenerationPhaseService_ExecutionPromptDigestChangesWhenPromptInputChanges(t *testing.T) {
	runA := &repository.JobPhaseRun{
		ID:               20,
		TranslationJobID: 1,
		PhaseType:        personaGenerationPhaseType,
		State:            personaGenerationPhaseStateCompleted,
		AIProvider:       "fake",
		ModelName:        "m",
		ExecutionMode:    "single_request",
		CredentialRef:    "cred",
	}
	recordsA := []repository.TranslationRecord{{ID: 100, RecordType: "NPC_", EditorID: "NPC_A", FormID: "0001"}}
	serviceA, repoA := newPersonaPhaseServiceForTest(personaGenerationJobStateRunning, personaGenerationPhaseStateCompleted, runA, recordsA)
	jobID := int64(1)
	repoA.phaseRunPersonas = []repository.PhaseRunPersona{{ID: 1, PhaseRunID: 20, PersonaID: 777, Role: personaGenerationPhaseLinkRoleApplied}}
	serviceA.foundationDataRepository = &fakePersonaPhaseFoundationDataRepository{
		personaByID: map[int64]repository.Persona{
			777: {ID: 777, NpcProfileID: 300, TranslationJobID: &jobID},
		},
		personaFieldEvidenceByID: map[int64][]repository.PersonaFieldEvidence{},
		nextPersonaID:            1000,
	}

	summaryA, err := serviceA.ReadSummary(context.Background(), 1)
	if err != nil {
		t.Fatalf("expected first summary success, got %v", err)
	}

	runB := &repository.JobPhaseRun{
		ID:               20,
		TranslationJobID: 1,
		PhaseType:        personaGenerationPhaseType,
		State:            personaGenerationPhaseStateCompleted,
		AIProvider:       "fake",
		ModelName:        "m",
		ExecutionMode:    "single_request",
		CredentialRef:    "cred",
	}
	recordsB := []repository.TranslationRecord{{ID: 100, RecordType: "NPC_", EditorID: "NPC_A", FormID: "0001"}}
	serviceB, repoB := newPersonaPhaseServiceForTest(personaGenerationJobStateRunning, personaGenerationPhaseStateCompleted, runB, recordsB)
	repoB.phaseRunPersonas = []repository.PhaseRunPersona{{ID: 1, PhaseRunID: 20, PersonaID: 777, Role: personaGenerationPhaseLinkRoleApplied}}
	serviceB.foundationDataRepository = &fakePersonaPhaseFoundationDataRepository{
		personaByID: map[int64]repository.Persona{
			777: {ID: 777, NpcProfileID: 300, TranslationJobID: &jobID},
		},
		personaFieldEvidenceByID: map[int64][]repository.PersonaFieldEvidence{},
		nextPersonaID:            1000,
	}
	sourceB := serviceB.translationSourceReader.(*fakePersonaPhaseTranslationSourceRepository)
	sourceB.fieldsByRecord[100] = []repository.TranslationField{{ID: 400, TranslationRecordID: 100, SourceText: "changed utterance"}}

	summaryB, err := serviceB.ReadSummary(context.Background(), 1)
	if err != nil {
		t.Fatalf("expected second summary success, got %v", err)
	}
	if summaryA.TargetSummary.TargetSnapshotDigest != summaryB.TargetSummary.TargetSnapshotDigest {
		t.Fatalf("expected target snapshot digest to remain stable, got %q and %q", summaryA.TargetSummary.TargetSnapshotDigest, summaryB.TargetSummary.TargetSnapshotDigest)
	}
	if summaryA.Execution.PromptDigest == summaryB.Execution.PromptDigest {
		t.Fatalf("expected prompt digest to change with prompt input, got %q", summaryA.Execution.PromptDigest)
	}
}

// ---------------------------------------------------------------------------
// U-DTO-002: AI 設定 record 不在時に AISettings field が nil を返す（ペルソナ生成）
// 根拠: persona-generation-phase-REQ-002「JOB_PHASE_AI_SETTINGS record 不在は応答の AI 設定 field 不在で表現する」
// ---------------------------------------------------------------------------

func TestPersonaGenerationPhaseServiceReadSummaryReturnsNilAISettingsWhenRecordAbsent(t *testing.T) {
	// AI 設定 record が存在しない状態で ReadSummary を呼ぶと AISettings が nil を返すことを証明する。
	// 空文字代理表現ではなく field 不在（nil）で「未設定」を表す仕様の確認。
	sourceRecords := []repository.TranslationRecord{{ID: 21, RecordType: "NPC_", EditorID: "Lydia", FormID: "000ABC"}}
	service, _ := newPersonaPhaseServiceForTest(personaGenerationJobStateRunning, personaGenerationPhaseStateCompleted, nil, sourceRecords)
	// AI 設定 repository を record 空の InMemory に差し替える（record を保存しない）。
	service.phaseAISettingsRepository = repository.NewInMemoryJobPhaseAISettingsRepository()

	// Arrange: AI 設定 record が存在しない状態

	// Act
	summary, err := service.ReadSummary(context.Background(), 1)
	if err != nil {
		t.Fatalf("expected summary success when AI settings record absent: %v", err)
	}

	// Assert: AISettings field が nil であることを確認する（空文字 record ではない）
	if summary.AISettings != nil {
		t.Errorf("expected AISettings=nil when record absent, got %#v", summary.AISettings)
	}
}

// ---------------------------------------------------------------------------
// U-DTO-003: AI 設定 record 存在時に AISettings field が値を返す（ペルソナ生成）
// 根拠: persona-generation-phase-REQ-002「record 存在時は保持 3 値を返す」
// ---------------------------------------------------------------------------

func TestPersonaGenerationPhaseServiceReadSummaryReturnsAISettingsWhenRecordExists(t *testing.T) {
	// AI 設定 record が存在する状態で ReadSummary を呼ぶと AISettings が値を返すことを証明する。
	sourceRecords := []repository.TranslationRecord{{ID: 21, RecordType: "NPC_", EditorID: "Lydia", FormID: "000ABC"}}
	service, _ := newPersonaPhaseServiceForTest(personaGenerationJobStateRunning, personaGenerationPhaseStateCompleted, nil, sourceRecords)

	// Arrange: AI 設定 record を明示的に保存する
	repo := repository.NewInMemoryJobPhaseAISettingsRepository()
	if _, err := repo.Save(context.Background(), repository.JobPhaseAISettingsDraft{
		PhaseType:     "npc_persona_generation",
		AIProvider:    "xai",
		ModelName:     "grok-3",
		ExecutionMode: "single_request",
		BatchMode:     "disabled",
	}); err != nil {
		t.Fatalf("prepare AI settings record: %v", err)
	}
	service.WithPersonaGenerationPhaseAISettingsRepository(repo)

	// Act
	summary, err := service.ReadSummary(context.Background(), 1)
	if err != nil {
		t.Fatalf("expected summary success when AI settings record exists: %v", err)
	}

	// Assert: AISettings field が nil でなく、保存値を返すことを確認する
	if summary.AISettings == nil {
		t.Fatal("expected AISettings to be non-nil when record exists")
	}
	if summary.AISettings.Provider != "xai" {
		t.Errorf("expected AISettings.Provider=xai, got %q", summary.AISettings.Provider)
	}
	if summary.AISettings.Model != "grok-3" {
		t.Errorf("expected AISettings.Model=grok-3, got %q", summary.AISettings.Model)
	}
}
