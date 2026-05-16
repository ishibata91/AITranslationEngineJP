package service

import (
	"context"
	"errors"
	"fmt"
	"reflect"
	"strings"
	"testing"
	"time"

	"aitranslationenginejp/internal/repository"
)

var translationJobSetupProviderSettingsPlainSecret = string([]byte{
	116, 106, 45, 115, 112, 112, 115, 45, 115, 101, 99, 114, 101, 116, 45, 118, 97, 108, 117, 101,
})

func fixedTranslationJobSetupValidationNow() time.Time {
	return time.Date(2026, 5, 4, 13, 0, 0, 0, time.UTC)
}

type fakeTranslationJobSetupSourceRepository struct {
	getByID         func(context.Context, int64) (repository.XEditExtractedData, error)
	listAll         func(context.Context) ([]repository.XEditExtractedData, error)
	getExistingJob  func(context.Context, int64) (repository.TranslationJob, error)
	deleteCacheByID func(context.Context, int64) error
	deleteInputByID func(context.Context, int64) error
}

func (repo fakeTranslationJobSetupSourceRepository) GetXEditExtractedDataByID(ctx context.Context, id int64) (repository.XEditExtractedData, error) {
	if repo.getByID == nil {
		return repository.XEditExtractedData{}, nil
	}
	return repo.getByID(ctx, id)
}

func (repo fakeTranslationJobSetupSourceRepository) ListXEditExtractedData(ctx context.Context) ([]repository.XEditExtractedData, error) {
	if repo.listAll == nil {
		return nil, nil
	}
	return repo.listAll(ctx)
}

func (repo fakeTranslationJobSetupSourceRepository) GetExistingTranslationJob(ctx context.Context, xEditID int64) (repository.TranslationJob, error) {
	if repo.getExistingJob == nil {
		return repository.TranslationJob{}, repository.ErrNotFound
	}
	return repo.getExistingJob(ctx, xEditID)
}

func (repo fakeTranslationJobSetupSourceRepository) DeleteTranslationCacheByXEditID(ctx context.Context, xEditID int64) error {
	if repo.deleteCacheByID == nil {
		return nil
	}
	return repo.deleteCacheByID(ctx, xEditID)
}

func (repo fakeTranslationJobSetupSourceRepository) DeleteXEditExtractedDataByID(ctx context.Context, id int64) error {
	if repo.deleteInputByID == nil {
		return nil
	}
	return repo.deleteInputByID(ctx, id)
}

type fakeTranslationJobSetupDictionaryRepository struct{}

func (fakeTranslationJobSetupDictionaryRepository) List(context.Context, MasterDictionaryQuery) (MasterDictionaryListResult, error) {
	return MasterDictionaryListResult{}, nil
}

type fakeTranslationJobSetupPersonaRepository struct{}

func (fakeTranslationJobSetupPersonaRepository) List(context.Context, repository.MasterPersonaListQuery) (repository.MasterPersonaListResult, error) {
	return repository.MasterPersonaListResult{}, nil
}

type fakeTranslationJobSetupSecretStore struct {
	load func(context.Context, string) (string, error)
}

func (store fakeTranslationJobSetupSecretStore) Load(ctx context.Context, key string) (string, error) {
	return store.load(ctx, key)
}

type fakeTranslationJobSetupProviderSettingsConsumer struct {
	listProviderSettings          func(context.Context) (ProviderSettingsRoute, []ProviderSettingsSummary, error)
	listProviderModels            func(context.Context, ProviderSettingsModelListInput) (ProviderSettingsModelListResult, error)
	resolveProviderExecutionInput func(context.Context, ProviderSettingsResolveInput) (ProviderSettingsResolveResult, error)
}

func (consumer fakeTranslationJobSetupProviderSettingsConsumer) ListProviderSettings(
	ctx context.Context,
) (ProviderSettingsRoute, []ProviderSettingsSummary, error) {
	if consumer.listProviderSettings == nil {
		return ProviderSettingsRoute{}, nil, nil
	}
	return consumer.listProviderSettings(ctx)
}

func (consumer fakeTranslationJobSetupProviderSettingsConsumer) ListProviderModels(
	ctx context.Context,
	input ProviderSettingsModelListInput,
) (ProviderSettingsModelListResult, error) {
	if consumer.listProviderModels == nil {
		return ProviderSettingsModelListResult{}, nil
	}
	return consumer.listProviderModels(ctx, input)
}

func (consumer fakeTranslationJobSetupProviderSettingsConsumer) ResolveProviderExecutionSettings(
	ctx context.Context,
	input ProviderSettingsResolveInput,
) (ProviderSettingsResolveResult, error) {
	if consumer.resolveProviderExecutionInput == nil {
		return ProviderSettingsResolveResult{}, nil
	}
	return consumer.resolveProviderExecutionInput(ctx, input)
}

type fakeTranslationJobSetupTransactor struct{}

func (fakeTranslationJobSetupTransactor) WithTransaction(ctx context.Context, callback func(context.Context) error) error {
	return callback(ctx)
}

type fakeTranslationJobSetupJobLifecycleRepository struct {
	createdJobs      []repository.TranslationJobDraft
	createdPhaseRuns []repository.JobPhaseRunDraft
	savedSnapshots   []repository.TranslationJobPhaseRuntimeSnapshotDraft
	jobByID          repository.TranslationJob
	summarySnapshots []repository.TranslationJobPhaseRuntimeSnapshot
}

func (repo *fakeTranslationJobSetupJobLifecycleRepository) CreateTranslationJob(_ context.Context, draft repository.TranslationJobDraft) (repository.TranslationJob, error) {
	repo.createdJobs = append(repo.createdJobs, draft)
	return repository.TranslationJob{
		ID:                   91,
		XEditExtractedDataID: draft.XEditExtractedDataID,
		JobName:              draft.JobName,
		State:                draft.State,
	}, nil
}

func (repo *fakeTranslationJobSetupJobLifecycleRepository) GetTranslationJobByID(context.Context, int64) (repository.TranslationJob, error) {
	return repo.jobByID, nil
}

func (repo *fakeTranslationJobSetupJobLifecycleRepository) CreateJobPhaseRun(_ context.Context, draft repository.JobPhaseRunDraft) (repository.JobPhaseRun, error) {
	repo.createdPhaseRuns = append(repo.createdPhaseRuns, draft)
	return repository.JobPhaseRun{ID: 11, TranslationJobID: draft.TranslationJobID, PhaseType: draft.PhaseType}, nil
}

func (repo *fakeTranslationJobSetupJobLifecycleRepository) ListJobPhaseRunsByJobID(context.Context, int64) ([]repository.JobPhaseRun, error) {
	return nil, nil
}

func (repo *fakeTranslationJobSetupJobLifecycleRepository) SaveTranslationJobPhaseRuntimeSnapshot(_ context.Context, draft repository.TranslationJobPhaseRuntimeSnapshotDraft) (repository.TranslationJobPhaseRuntimeSnapshot, error) {
	repo.savedSnapshots = append(repo.savedSnapshots, draft)
	return repository.TranslationJobPhaseRuntimeSnapshot{
		ID:               int64(len(repo.savedSnapshots)),
		TranslationJobID: draft.TranslationJobID,
		PhaseID:          draft.PhaseID,
		Provider:         draft.Provider,
		ModelName:        draft.ModelName,
		CredentialStatus: draft.CredentialStatus,
		ExecutionMode:    draft.ExecutionMode,
		BatchMode:        draft.BatchMode,
		CreatedAt:        time.Date(2026, 5, 4, 12, 0, 0, 0, time.UTC),
	}, nil
}

func (repo *fakeTranslationJobSetupJobLifecycleRepository) ListTranslationJobPhaseRuntimeSnapshots(context.Context, int64) ([]repository.TranslationJobPhaseRuntimeSnapshot, error) {
	return append([]repository.TranslationJobPhaseRuntimeSnapshot(nil), repo.summarySnapshots...), nil
}

func TestTJSPPS002TranslationJobSetupServiceIgnoresMasterPersonaProviderSettings(t *testing.T) {
	secretLoadCalls := []string{}
	loadCalls := 0
	service := NewPersistentTranslationJobSetupService(
		&fakeTranslationJobSetupJobLifecycleRepository{},
		fakeTranslationJobSetupSourceRepository{
			getByID: func(context.Context, int64) (repository.XEditExtractedData, error) {
				return repository.XEditExtractedData{ID: 44}, nil
			},
		},
		fakeTranslationJobSetupDictionaryRepository{},
		fakeTranslationJobSetupPersonaRepository{},
		nil,
		fakeTranslationJobSetupSecretStore{
			load: func(_ context.Context, key string) (string, error) {
				secretLoadCalls = append(secretLoadCalls, key)
				return translationJobSetupProviderSettingsPlainSecret, nil
			},
		},
		fakeTranslationJobSetupTransactor{},
		WithTranslationJobSetupProviderModelListLoader(
			TranslationJobSetupProviderModelListLoaderFunc(
				func(context.Context, string, string) ([]TranslationJobSetupProviderModelOptionReadModel, error) {
					loadCalls++
					return nil, errors.New("SCN-TJSPPS-002: missing phase runtime must not call provider model list loader")
				},
			),
		),
	)

	decision, err := service.ValidateRequest(context.Background(), TranslationJobSetupValidationRequest{
		InputSourceID: 44,
		PhaseRuntimes: []TranslationJobSetupPhaseRuntimeDraftReadModel{
			{PhaseID: "word_translation"},
			{PhaseID: "npc_persona_generation"},
			{PhaseID: "text_translation"},
		},
	})
	if err != nil {
		t.Fatalf("expected validation classification, got error: %v", err)
	}
	if decision.BlockingFailureCategory == nil || *decision.BlockingFailureCategory != "phase_runtime_missing" {
		t.Fatalf("SCN-TJSPPS-002: expected phase_runtime_missing, got %#v", decision)
	}
	if decision.CanCreate {
		t.Fatalf("SCN-TJSPPS-002: expected missing Job Setup phase settings to block create, got %#v", decision)
	}
	if containsTranslationJobSetupString(secretLoadCalls, "master-persona:openai") {
		t.Fatalf("SCN-TJSPPS-002: Job Setup must not resolve master-persona secret namespace, got %#v", secretLoadCalls)
	}
	if len(secretLoadCalls) != 0 {
		t.Fatalf("SCN-TJSPPS-002: missing phase runtime must not resolve any credential, got %#v", secretLoadCalls)
	}
	if loadCalls != 0 {
		t.Fatalf("SCN-TJSPPS-002: missing phase runtime must not call provider model list loader, got %d", loadCalls)
	}
}

func TestTJSPPS003TranslationJobSetupServiceProviderModelListRunsOnlyAfterCredentialGate(t *testing.T) {
	loadCalls := []struct {
		provider string
		apiKey   string
	}{}
	service := NewPersistentTranslationJobSetupService(
		&fakeTranslationJobSetupJobLifecycleRepository{},
		fakeTranslationJobSetupSourceRepository{},
		fakeTranslationJobSetupDictionaryRepository{},
		fakeTranslationJobSetupPersonaRepository{},
		nil,
		fakeTranslationJobSetupSecretStore{load: func(_ context.Context, key string) (string, error) {
			if key != "openai-primary" && key != "gemini-primary" {
				t.Fatalf("SCN-TJSPPS-003: unexpected credential key lookup: %q", key)
			}
			return translationJobSetupProviderSettingsPlainSecret, nil
		}},
		fakeTranslationJobSetupTransactor{},
		WithTranslationJobSetupProviderModelListLoader(
			TranslationJobSetupProviderModelListLoaderFunc(
				func(_ context.Context, providerID string, apiKey string) ([]TranslationJobSetupProviderModelOptionReadModel, error) {
					loadCalls = append(loadCalls, struct {
						provider string
						apiKey   string
					}{provider: providerID, apiKey: apiKey})
					if apiKey == "" {
						return nil, errors.New("SCN-TJSPPS-003: configured credential must be sent to fake loader")
					}
					return []TranslationJobSetupProviderModelOptionReadModel{{ModelID: "gemini-2.5-pro", Label: "gemini-2.5-pro"}}, nil
				},
			),
		),
	)

	missingCredential, err := service.ListProviderModels(context.Background(), ListTranslationJobSetupProviderModelsRequest{
		PhaseID:          "word_translation",
		Provider:         "openai",
		CredentialRef:    "openai-primary",
		CredentialStatus: "missing",
		RequestToken:     "req-1",
	})
	if err != nil {
		t.Fatalf("SCN-TJSPPS-003: expected missing credential classification: %v", err)
	}
	if missingCredential.Status != "credential_missing" || missingCredential.FailureKind != "model_list_credential_missing" {
		t.Fatalf("SCN-TJSPPS-003: expected credential-missing result, got %#v", missingCredential)
	}
	if len(loadCalls) != 0 {
		t.Fatalf("SCN-TJSPPS-003: expected no loader call before credential gate, got %#v", loadCalls)
	}

	configured, err := service.ListProviderModels(context.Background(), ListTranslationJobSetupProviderModelsRequest{
		PhaseID:          "word_translation",
		Provider:         "gemini",
		CredentialRef:    "gemini-primary",
		CredentialStatus: "configured",
		RequestToken:     "req-2",
	})
	if err != nil {
		t.Fatalf("SCN-TJSPPS-003: expected configured credential model list to succeed: %v", err)
	}
	if configured.Status != "success" {
		t.Fatalf("SCN-TJSPPS-003: expected success status, got %#v", configured)
	}
	if len(configured.Models) != 1 || configured.Models[0].ModelID != "gemini-2.5-pro" {
		t.Fatalf("SCN-TJSPPS-003: expected fake model response to populate model options, got %#v", configured.Models)
	}
	if len(loadCalls) != 1 || loadCalls[0].provider != "gemini" || loadCalls[0].apiKey == "" {
		t.Fatalf("SCN-TJSPPS-003: expected exactly one fake loader call after credential gate, got %#v", loadCalls)
	}
}

func TestTranslationJobSetupServiceRejectsSharedSecretNamespaceCredentialRefs(t *testing.T) {
	loadCalls := 0
	secretLoadCalls := []string{}
	service := NewPersistentTranslationJobSetupService(
		&fakeTranslationJobSetupJobLifecycleRepository{},
		fakeTranslationJobSetupSourceRepository{
			getByID: func(context.Context, int64) (repository.XEditExtractedData, error) {
				return repository.XEditExtractedData{ID: 44}, nil
			},
		},
		fakeTranslationJobSetupDictionaryRepository{},
		fakeTranslationJobSetupPersonaRepository{},
		nil,
		fakeTranslationJobSetupSecretStore{
			load: func(_ context.Context, key string) (string, error) {
				secretLoadCalls = append(secretLoadCalls, key)
				return translationJobSetupProviderSettingsPlainSecret, nil
			},
		},
		fakeTranslationJobSetupTransactor{},
		WithTranslationJobSetupProviderModelListLoader(
			TranslationJobSetupProviderModelListLoaderFunc(
				func(context.Context, string, string) ([]TranslationJobSetupProviderModelOptionReadModel, error) {
					loadCalls++
					return []TranslationJobSetupProviderModelOptionReadModel{{ModelID: "blocked", Label: "blocked"}}, nil
				},
			),
		),
	)

	models, err := service.ListProviderModels(context.Background(), ListTranslationJobSetupProviderModelsRequest{
		PhaseID:          "word_translation",
		Provider:         "gemini",
		CredentialRef:    "master-persona:gemini",
		CredentialStatus: "configured",
		RequestToken:     "req-shared-secret",
	})
	if err != nil {
		t.Fatalf("expected blocked shared namespace result, got error: %v", err)
	}
	if models.Status != "credential_missing" || models.CredentialStatus != "missing" {
		t.Fatalf("expected shared namespace credential to be blocked before model list, got %#v", models)
	}
	if len(secretLoadCalls) != 0 || loadCalls != 0 {
		t.Fatalf("expected shared namespace credential to avoid secret load and loader calls, got loads=%#v loaderCalls=%d", secretLoadCalls, loadCalls)
	}

	validation, err := service.ValidateRequest(context.Background(), TranslationJobSetupValidationRequest{
		InputSourceID: 44,
		PhaseRuntimes: []TranslationJobSetupPhaseRuntimeDraftReadModel{
			{PhaseID: "word_translation", Provider: "openai", Model: "gpt-5.4-mini", CredentialRef: "openai-primary", CredentialStatus: "configured", ExecutionMode: "sync", BatchMode: "unsupported", ModelListSourceToken: "word_translation|openai|openai-primary|req-1"},
			{PhaseID: "npc_persona_generation", Provider: "gemini", Model: "gemini-2.5-pro", CredentialRef: "master-persona:gemini", CredentialStatus: "configured", ExecutionMode: "sync", BatchMode: "disabled", ModelListSourceToken: "npc_persona_generation|gemini|master-persona:gemini|req-shared-secret"},
			{PhaseID: "text_translation", Provider: "xai", Model: "grok-4", CredentialRef: "xai-primary", CredentialStatus: "configured", ExecutionMode: "sync", BatchMode: "disabled", ModelListSourceToken: "text_translation|xai|xai-primary|req-3"},
		},
	})
	if err != nil {
		t.Fatalf("expected validation result for shared namespace credential: %v", err)
	}
	if validation.BlockingFailureCategory == nil || *validation.BlockingFailureCategory != "credential_missing" {
		t.Fatalf("expected shared namespace credential to fail validation as missing, got %#v", validation)
	}
}

func TestTJSPPS003TranslationJobSetupServiceLMStudioModelListSkipsSecretLoad(t *testing.T) {
	secretLoadCalls := []string{}
	service := NewPersistentTranslationJobSetupService(
		&fakeTranslationJobSetupJobLifecycleRepository{},
		fakeTranslationJobSetupSourceRepository{},
		fakeTranslationJobSetupDictionaryRepository{},
		fakeTranslationJobSetupPersonaRepository{},
		nil,
		fakeTranslationJobSetupSecretStore{load: func(_ context.Context, key string) (string, error) {
			secretLoadCalls = append(secretLoadCalls, key)
			return translationJobSetupProviderSettingsPlainSecret, nil
		}},
		fakeTranslationJobSetupTransactor{},
		WithTranslationJobSetupProviderModelListLoader(
			TranslationJobSetupProviderModelListLoaderFunc(
				func(context.Context, string, string) ([]TranslationJobSetupProviderModelOptionReadModel, error) {
					return []TranslationJobSetupProviderModelOptionReadModel{
						{ModelID: "lmstudio-community", Label: "lmstudio-community"},
					}, nil
				},
			),
		),
	)

	result, err := service.ListProviderModels(context.Background(), ListTranslationJobSetupProviderModelsRequest{
		PhaseID:          "text_translation",
		Provider:         "lm_studio",
		CredentialRef:    "",
		CredentialStatus: "not_required",
		RequestToken:     "req-lmstudio",
	})
	if err != nil {
		t.Fatalf("expected lm studio model list success: %v", err)
	}
	if result.Status != "credential_not_required" {
		t.Fatalf("expected credential_not_required for lm studio, got %#v", result)
	}
	if len(secretLoadCalls) != 0 {
		t.Fatalf("expected no secret lookup for lm studio model list, got %#v", secretLoadCalls)
	}
}

func TestTJSPPS007TranslationJobSetupServiceCreateCapturesOnlyTargetPhaseRuntimeSettings(t *testing.T) {
	jobRepo := &fakeTranslationJobSetupJobLifecycleRepository{}
	service := NewPersistentTranslationJobSetupService(
		jobRepo,
		fakeTranslationJobSetupSourceRepository{
			getByID: func(context.Context, int64) (repository.XEditExtractedData, error) {
				return repository.XEditExtractedData{ID: 44}, nil
			},
		},
		fakeTranslationJobSetupDictionaryRepository{},
		fakeTranslationJobSetupPersonaRepository{},
		nil,
		fakeTranslationJobSetupSecretStore{load: func(context.Context, string) (string, error) { return "configured", nil }},
		fakeTranslationJobSetupTransactor{},
	)
	service.now = fixedTranslationJobSetupValidationNow

	created, err := service.CreateTranslationJob(context.Background(), TranslationJobSetupCreateRequest{
		InputSourceID:    44,
		ValidationStatus: "pass",
		ValidatedAt:      time.Date(2026, 5, 4, 12, 0, 0, 0, time.UTC),
		PhaseRuntimes: []TranslationJobSetupPhaseRuntimeDraftReadModel{
			{PhaseID: "word_translation", Provider: "openai", Model: "gpt-5.4-mini", CredentialRef: "openai-primary", CredentialStatus: "configured", ExecutionMode: "sync", BatchMode: "unsupported", ModelListSourceToken: "word_translation|openai|openai-primary|req-1"},
			{PhaseID: "npc_persona_generation", Provider: "gemini", Model: "gemini-2.5-pro", CredentialRef: "gemini-primary", CredentialStatus: "configured", ExecutionMode: "sync", BatchMode: "disabled", ModelListSourceToken: "npc_persona_generation|gemini|gemini-primary|req-2"},
			{PhaseID: "text_translation", Provider: "openai", Model: "gpt-5.4-mini", CredentialRef: "openai-primary", CredentialStatus: "configured", ExecutionMode: "sync", BatchMode: "enabled", ModelListSourceToken: "text_translation|openai|openai-primary|req-3"},
		},
	}, []string{"input", "runtime", "credentials"})
	if err != nil {
		t.Fatalf("SCN-TJSPPS-007: expected create success: %v", err)
	}
	if len(jobRepo.savedSnapshots) != 3 {
		t.Fatalf("SCN-TJSPPS-007: expected three saved snapshots, got %#v", jobRepo.savedSnapshots)
	}
	if len(jobRepo.createdPhaseRuns) != 0 {
		t.Fatalf("SCN-TJSPPS-007: expected no pre-created JOB_PHASE_RUN placeholders, got %#v", jobRepo.createdPhaseRuns)
	}
	if jobRepo.savedSnapshots[2].BatchMode != "unsupported" {
		t.Fatalf("SCN-TJSPPS-007: expected stale batch mode to be stripped for openai, got %#v", jobRepo.savedSnapshots[2])
	}
	want := []TranslationJobSetupPhaseRuntimeSummaryReadModel{
		{PhaseID: "word_translation", Provider: "openai", Model: "gpt-5.4-mini", CredentialRef: "openai-primary", CredentialStatus: "configured", ExecutionMode: "sync", BatchMode: "unsupported", ModelListSourceToken: "word_translation|openai|openai-primary|req-1"},
		{PhaseID: "npc_persona_generation", Provider: "gemini", Model: "gemini-2.5-pro", CredentialRef: "gemini-primary", CredentialStatus: "configured", ExecutionMode: "sync", BatchMode: "disabled", ModelListSourceToken: "npc_persona_generation|gemini|gemini-primary|req-2"},
		{PhaseID: "text_translation", Provider: "openai", Model: "gpt-5.4-mini", CredentialRef: "openai-primary", CredentialStatus: "configured", ExecutionMode: "sync", BatchMode: "unsupported", ModelListSourceToken: "text_translation|openai|openai-primary|req-3"},
	}
	if !reflect.DeepEqual(created.PhaseRuntimeSummaries, want) {
		t.Fatalf("SCN-TJSPPS-007: expected phase-only runtime summaries %#v, got %#v", want, created.PhaseRuntimeSummaries)
	}
}

func TestTJSPPS008TranslationJobSetupServiceProviderSettingsExposeNoSecretAndUseFakeTransportOnly(t *testing.T) {
	secretLoadCalls := []string{}
	loadCalls := []struct {
		provider string
		apiKey   string
	}{}
	service := NewPersistentTranslationJobSetupService(
		&fakeTranslationJobSetupJobLifecycleRepository{},
		fakeTranslationJobSetupSourceRepository{},
		fakeTranslationJobSetupDictionaryRepository{},
		fakeTranslationJobSetupPersonaRepository{},
		nil,
		fakeTranslationJobSetupSecretStore{load: func(_ context.Context, key string) (string, error) {
			secretLoadCalls = append(secretLoadCalls, key)
			return translationJobSetupProviderSettingsPlainSecret, nil
		}},
		fakeTranslationJobSetupTransactor{},
		WithTranslationJobSetupProviderModelListLoader(
			TranslationJobSetupProviderModelListLoaderFunc(
				func(_ context.Context, providerID string, apiKey string) ([]TranslationJobSetupProviderModelOptionReadModel, error) {
					loadCalls = append(loadCalls, struct {
						provider string
						apiKey   string
					}{provider: providerID, apiKey: apiKey})
					return []TranslationJobSetupProviderModelOptionReadModel{{ModelID: "gemini-2.5-pro", Label: "gemini-2.5-pro"}}, nil
				},
			),
		),
	)

	options, err := service.ReadOptions(context.Background())
	if err != nil {
		t.Fatalf("SCN-TJSPPS-008: expected provider settings options to load: %v", err)
	}
	for _, option := range options.AIRuntimeOptions {
		if option.Provider == "fake" {
			t.Fatalf("SCN-TJSPPS-008: fake provider must not be user-facing, got %#v", options.AIRuntimeOptions)
		}
	}
	for _, ref := range options.CredentialRefs {
		if strings.Contains(fmt.Sprintf("%#v", ref), translationJobSetupProviderSettingsPlainSecret) {
			t.Fatalf("SCN-TJSPPS-008: credential read model must not expose secret, got %#v", ref)
		}
	}
	secretLoadCalls = nil

	result, err := service.ListProviderModels(context.Background(), ListTranslationJobSetupProviderModelsRequest{
		PhaseID:          "text_translation",
		Provider:         "gemini",
		CredentialRef:    "gemini-primary",
		CredentialStatus: "configured",
		RequestToken:     "body-gemini-configured",
	})
	if err != nil {
		t.Fatalf("SCN-TJSPPS-008: expected model list to use fake transport without paid API: %v", err)
	}
	if len(loadCalls) != 1 || loadCalls[0].provider != "gemini" || loadCalls[0].apiKey == "" {
		t.Fatalf("SCN-TJSPPS-008: expected one fake loader call, got %#v", loadCalls)
	}
	if len(secretLoadCalls) != 1 || secretLoadCalls[0] != "gemini-primary" {
		t.Fatalf("SCN-TJSPPS-008: expected fake secret store to be the credential source, got %#v", secretLoadCalls)
	}
	if strings.Contains(fmt.Sprintf("%#v", result), translationJobSetupProviderSettingsPlainSecret) {
		t.Fatalf("SCN-TJSPPS-008: provider model result must not expose secret, got %#v", result)
	}
}

func TestTranslationJobSetupServiceProviderSettingsTestSafeModelListAllowsMissingCredential(t *testing.T) {
	providerSettingsModelListCalls := []ProviderSettingsModelListInput{}
	service := NewPersistentTranslationJobSetupService(
		&fakeTranslationJobSetupJobLifecycleRepository{},
		fakeTranslationJobSetupSourceRepository{},
		fakeTranslationJobSetupDictionaryRepository{},
		fakeTranslationJobSetupPersonaRepository{},
		nil,
		fakeTranslationJobSetupSecretStore{
			load: func(context.Context, string) (string, error) {
				t.Fatal("expected secret load to be skipped when provider settings handles test-safe model list")
				return "", nil
			},
		},
		fakeTranslationJobSetupTransactor{},
		WithTranslationJobSetupProviderSettings(
			fakeTranslationJobSetupProviderSettingsConsumer{
				listProviderSettings: func(context.Context) (ProviderSettingsRoute, []ProviderSettingsSummary, error) {
					return ProviderSettingsRoute{}, []ProviderSettingsSummary{
						{
							ProviderID:            "gemini",
							Label:                 "Gemini",
							Endpoint:              nil,
							CredentialState:       "missing",
							CredentialReferenceID: nil,
							ValidationState:       "not_validated",
							SavedState:            "partial",
							RequestToken:          stringPointer("gemini|test-safe"),
							LastFailureKind:       nil,
						},
					}, nil
				},
				listProviderModels: func(_ context.Context, input ProviderSettingsModelListInput) (ProviderSettingsModelListResult, error) {
					providerSettingsModelListCalls = append(providerSettingsModelListCalls, input)
					return ProviderSettingsModelListResult{
						ProviderID:      "gemini",
						Endpoint:        nil,
						CredentialState: "not_required",
						RequestToken:    "gemini|test-safe",
						State:           "ready",
						Models: []ProviderSettingsModelOption{
							{ModelID: "gemini-test-safe", Label: "Gemini Test Safe"},
						},
						FailureKind: nil,
					}, nil
				},
				resolveProviderExecutionInput: func(context.Context, ProviderSettingsResolveInput) (ProviderSettingsResolveResult, error) {
					return ProviderSettingsResolveResult{}, nil
				},
			},
		),
	)

	result, err := service.ListProviderModels(context.Background(), ListTranslationJobSetupProviderModelsRequest{
		PhaseID:          "word_translation",
		Provider:         "gemini",
		CredentialStatus: "missing",
		RequestToken:     "ui-req-1",
	})
	if err != nil {
		t.Fatalf("expected provider settings test-safe model list success: %v", err)
	}
	if result.Status != "success" {
		t.Fatalf("expected success status via provider settings test-safe model list, got %#v", result)
	}
	if result.CredentialStatus != "not_required" {
		t.Fatalf("expected credential status to switch to not_required via provider settings test-safe model list, got %#v", result)
	}
	if len(result.Models) != 1 || result.Models[0].ModelID != "gemini-test-safe" {
		t.Fatalf("expected model list response via provider settings test-safe model list, got %#v", result.Models)
	}
	if len(providerSettingsModelListCalls) != 1 {
		t.Fatalf("expected one provider settings model list call, got %#v", providerSettingsModelListCalls)
	}
	if providerSettingsModelListCalls[0].RequestToken != "gemini|test-safe" {
		t.Fatalf("expected provider settings request token from summary, got %#v", providerSettingsModelListCalls[0])
	}
	if providerSettingsModelListCalls[0].CredentialState != "missing" ||
		providerSettingsModelListCalls[0].CredentialReferenceID != nil ||
		providerSettingsModelListCalls[0].Endpoint != nil {
		t.Fatalf("expected provider settings input to preserve missing snapshot, got %#v", providerSettingsModelListCalls[0])
	}
	if result.RequestToken != "ui-req-1" {
		t.Fatalf("expected response request token to keep ui token, got %#v", result)
	}
}

func TestTranslationJobSetupServiceReadSummaryReturnsPersistedPhaseRuntimeSnapshots(t *testing.T) {
	jobRepo := &fakeTranslationJobSetupJobLifecycleRepository{
		jobByID: repository.TranslationJob{ID: 91, State: "ready"},
		summarySnapshots: []repository.TranslationJobPhaseRuntimeSnapshot{
			{PhaseID: "word_translation", Provider: "openai", ModelName: "gpt-5.4-mini", CredentialStatus: "configured", ExecutionMode: "sync", BatchMode: "unsupported"},
			{PhaseID: "npc_persona_generation", Provider: "gemini", ModelName: "gemini-2.5-pro", CredentialStatus: "configured", ExecutionMode: "batch", BatchMode: "enabled"},
			{PhaseID: "text_translation", Provider: "xai", ModelName: "grok-4", CredentialStatus: "configured", ExecutionMode: "sync", BatchMode: "disabled"},
		},
	}
	service := NewPersistentTranslationJobSetupService(
		jobRepo,
		fakeTranslationJobSetupSourceRepository{},
		fakeTranslationJobSetupDictionaryRepository{},
		fakeTranslationJobSetupPersonaRepository{},
		nil,
		nil,
		fakeTranslationJobSetupTransactor{},
	)

	summary, err := service.ReadSummary(context.Background(), 91)
	if err != nil {
		t.Fatalf("expected summary success: %v", err)
	}
	if !summary.CanStartPhase || len(summary.PhaseRuntimeSummaries) != 3 {
		t.Fatalf("expected persisted runtime snapshots in summary, got %#v", summary)
	}
	wantExecution := TranslationJobSetupExecutionSummaryReadModel{Provider: "openai", Model: "gpt-5.4-mini", ExecutionMode: "sync"}
	if !reflect.DeepEqual(summary.ExecutionSummary, wantExecution) {
		t.Fatalf("expected summary execution %#v, got %#v", wantExecution, summary.ExecutionSummary)
	}
}

func TestTranslationJobSetupServiceReadOptionsExcludesInputsReferencedByExistingJobsButKeepsExistingJobSummary(t *testing.T) {
	service := NewPersistentTranslationJobSetupService(
		&fakeTranslationJobSetupJobLifecycleRepository{},
		fakeTranslationJobSetupSourceRepository{
			listAll: func(context.Context) ([]repository.XEditExtractedData, error) {
				return []repository.XEditExtractedData{
					{ID: 41, TargetPluginName: "Input A", RecordCount: 12, ImportedAt: time.Date(2026, 5, 5, 10, 0, 0, 0, time.UTC)},
					{ID: 52, TargetPluginName: "Input B", RecordCount: 7, ImportedAt: time.Date(2026, 5, 5, 11, 0, 0, 0, time.UTC)},
				}, nil
			},
			getExistingJob: func(_ context.Context, xEditID int64) (repository.TranslationJob, error) {
				if xEditID == 41 {
					return repository.TranslationJob{ID: 900, XEditExtractedDataID: 41, State: "ready"}, nil
				}
				return repository.TranslationJob{}, repository.ErrNotFound
			},
		},
		fakeTranslationJobSetupDictionaryRepository{},
		fakeTranslationJobSetupPersonaRepository{},
		nil,
		fakeTranslationJobSetupSecretStore{load: func(context.Context, string) (string, error) { return "", repository.ErrNotFound }},
		fakeTranslationJobSetupTransactor{},
	)

	options, err := service.ReadOptions(context.Background())
	if err != nil {
		t.Fatalf("expected options to load: %v", err)
	}
	if len(options.InputCandidates) != 1 {
		t.Fatalf("expected only unreferenced input candidate, got %#v", options.InputCandidates)
	}
	if options.InputCandidates[0].ID != 52 {
		t.Fatalf("expected input 52 to remain visible, got %#v", options.InputCandidates[0])
	}
	wantExistingJob := &TranslationJobSetupExistingJobReadModel{
		InputSourceID: 41,
		JobID:         900,
		Status:        "ready",
		InputSource:   translationJobSetupInputSource,
	}
	if !reflect.DeepEqual(options.ExistingJob, wantExistingJob) {
		t.Fatalf("expected existing job summary %#v, got %#v", wantExistingJob, options.ExistingJob)
	}
}

func TestTranslationJobSetupServiceDeleteInputSourceRejectsReferencedInput(t *testing.T) {
	service := NewPersistentTranslationJobSetupService(
		&fakeTranslationJobSetupJobLifecycleRepository{},
		fakeTranslationJobSetupSourceRepository{
			getByID: func(context.Context, int64) (repository.XEditExtractedData, error) {
				return repository.XEditExtractedData{ID: 41}, nil
			},
			getExistingJob: func(context.Context, int64) (repository.TranslationJob, error) {
				return repository.TranslationJob{ID: 900, XEditExtractedDataID: 41, State: "ready"}, nil
			},
			deleteCacheByID: func(context.Context, int64) error {
				t.Fatal("expected delete cache to be skipped for referenced input")
				return nil
			},
			deleteInputByID: func(context.Context, int64) error {
				t.Fatal("expected delete source to be skipped for referenced input")
				return nil
			},
		},
		fakeTranslationJobSetupDictionaryRepository{},
		fakeTranslationJobSetupPersonaRepository{},
		nil,
		fakeTranslationJobSetupSecretStore{load: func(context.Context, string) (string, error) { return "", repository.ErrNotFound }},
		fakeTranslationJobSetupTransactor{},
	)

	decision, err := service.DeleteInputSource(context.Background(), 41)
	if err != nil {
		t.Fatalf("expected delete rejection without transport error: %v", err)
	}
	if decision.ErrorKind != translationJobSetupErrorKindInputDeleteBlocked {
		t.Fatalf("expected input_delete_blocked, got %#v", decision)
	}
	if decision.DeletedInputSourceID != nil {
		t.Fatalf("expected no deleted id on rejection, got %#v", decision)
	}
}

func TestTranslationJobSetupServiceReadSummaryRejectsStartWhenPhaseRuntimeSnapshotsAreIncomplete(t *testing.T) {
	jobRepo := &fakeTranslationJobSetupJobLifecycleRepository{
		jobByID: repository.TranslationJob{ID: 91, State: "ready"},
		summarySnapshots: []repository.TranslationJobPhaseRuntimeSnapshot{
			{PhaseID: "word_translation", Provider: "openai", ModelName: "gpt-5.4-mini", CredentialStatus: "configured", ExecutionMode: "sync", BatchMode: "unsupported"},
			{PhaseID: "npc_persona_generation", Provider: "gemini", ModelName: "gemini-2.5-pro", CredentialStatus: "configured", ExecutionMode: "sync", BatchMode: "disabled"},
		},
	}
	service := NewPersistentTranslationJobSetupService(
		jobRepo,
		fakeTranslationJobSetupSourceRepository{},
		fakeTranslationJobSetupDictionaryRepository{},
		fakeTranslationJobSetupPersonaRepository{},
		nil,
		nil,
		fakeTranslationJobSetupTransactor{},
	)

	summary, err := service.ReadSummary(context.Background(), 91)
	if err != nil {
		t.Fatalf("expected summary success: %v", err)
	}
	if summary.CanStartPhase {
		t.Fatalf("expected canStartPhase=false when phase runtime snapshots are incomplete, got %#v", summary)
	}
	if len(summary.PhaseRuntimeSummaries) != 2 {
		t.Fatalf("expected only persisted snapshot summaries, got %#v", summary.PhaseRuntimeSummaries)
	}
}

func TestTranslationJobSetupServiceValidateAcceptsLMStudioSourceTokenInValidationAndCreate(t *testing.T) {
	service := NewPersistentTranslationJobSetupService(
		&fakeTranslationJobSetupJobLifecycleRepository{},
		fakeTranslationJobSetupSourceRepository{
			getByID: func(context.Context, int64) (repository.XEditExtractedData, error) {
				return repository.XEditExtractedData{ID: 44}, nil
			},
		},
		fakeTranslationJobSetupDictionaryRepository{},
		fakeTranslationJobSetupPersonaRepository{},
		nil,
		nil,
		fakeTranslationJobSetupTransactor{},
	)
	service.now = fixedTranslationJobSetupValidationNow

	phaseRuntimes := []TranslationJobSetupPhaseRuntimeDraftReadModel{
		{PhaseID: "word_translation", Provider: "openai", Model: "gpt-5.4-mini", CredentialRef: "openai-primary", CredentialStatus: "configured", ExecutionMode: "sync", BatchMode: "unsupported", ModelListSourceToken: "word_translation|openai|openai-primary|req-1"},
		{PhaseID: "npc_persona_generation", Provider: "gemini", Model: "gemini-2.5-pro", CredentialRef: "gemini-primary", CredentialStatus: "configured", ExecutionMode: "sync", BatchMode: "disabled", ModelListSourceToken: "npc_persona_generation|gemini|gemini-primary|req-2"},
		{PhaseID: "text_translation", Provider: "lm_studio", Model: "lmstudio-community", CredentialRef: "lmstudio-local", CredentialStatus: "not_required", ExecutionMode: "sync", BatchMode: "unsupported", ModelListSourceToken: "text_translation|lm_studio||req-lm-1"},
	}

	validation, err := service.ValidateRequest(context.Background(), TranslationJobSetupValidationRequest{
		InputSourceID: 44,
		PhaseRuntimes: phaseRuntimes,
	})
	if err != nil {
		t.Fatalf("expected validation success with lm studio token: %v", err)
	}
	if !validation.CanCreate || validation.BlockingFailureCategory != nil {
		t.Fatalf("expected no stale failure for lm studio sourceToken, got %#v", validation)
	}

	createDecision, err := service.EvaluateCreateRequest(context.Background(), TranslationJobSetupCreateRequest{
		InputSourceID:    44,
		ValidationStatus: "pass",
		ValidatedAt:      time.Date(2026, 5, 4, 12, 0, 0, 0, time.UTC),
		PhaseRuntimes:    phaseRuntimes,
	})
	if err != nil {
		t.Fatalf("expected create evaluation success with lm studio token: %v", err)
	}
	if !createDecision.CanCreate {
		t.Fatalf("expected create to stay allowed with lm studio sourceToken, got %#v", createDecision)
	}
}

func TestTranslationJobSetupServiceListProviderModelsNormalizesLMStudioCredentialRefInSourceToken(t *testing.T) {
	service := NewPersistentTranslationJobSetupService(
		&fakeTranslationJobSetupJobLifecycleRepository{},
		fakeTranslationJobSetupSourceRepository{},
		fakeTranslationJobSetupDictionaryRepository{},
		fakeTranslationJobSetupPersonaRepository{},
		nil,
		fakeTranslationJobSetupSecretStore{
			load: func(context.Context, string) (string, error) {
				return translationJobSetupProviderSettingsPlainSecret, nil
			},
		},
		fakeTranslationJobSetupTransactor{},
		WithTranslationJobSetupProviderModelListLoader(
			TranslationJobSetupProviderModelListLoaderFunc(
				func(context.Context, string, string) ([]TranslationJobSetupProviderModelOptionReadModel, error) {
					return []TranslationJobSetupProviderModelOptionReadModel{{ModelID: "local-model", Label: "local-model"}}, nil
				},
			),
		),
	)

	result, err := service.ListProviderModels(context.Background(), ListTranslationJobSetupProviderModelsRequest{
		PhaseID:          "text_translation",
		Provider:         "lm_studio",
		CredentialRef:    "lmstudio-local",
		CredentialStatus: "configured",
		RequestToken:     "req-lm-token",
	})
	if err != nil {
		t.Fatalf("expected lm studio list success: %v", err)
	}
	if result.Status != "credential_not_required" {
		t.Fatalf("expected credential_not_required for lm studio, got %#v", result)
	}
	if result.SourceToken != "text_translation|lm_studio||req-lm-token" {
		t.Fatalf("expected lm studio source token to drop credential ref, got %#v", result.SourceToken)
	}
}

func TestTranslationJobSetupServiceListProviderModelsUsesLoaderWithoutTransportDetails(t *testing.T) {
	loadCalls := []struct {
		provider string
		apiKey   string
	}{}
	service := NewPersistentTranslationJobSetupService(
		&fakeTranslationJobSetupJobLifecycleRepository{},
		fakeTranslationJobSetupSourceRepository{},
		fakeTranslationJobSetupDictionaryRepository{},
		fakeTranslationJobSetupPersonaRepository{},
		nil,
		fakeTranslationJobSetupSecretStore{
			load: func(context.Context, string) (string, error) {
				return translationJobSetupProviderSettingsPlainSecret, nil
			},
		},
		fakeTranslationJobSetupTransactor{},
		WithTranslationJobSetupProviderModelListLoader(
			TranslationJobSetupProviderModelListLoaderFunc(
				func(_ context.Context, providerID string, apiKey string) ([]TranslationJobSetupProviderModelOptionReadModel, error) {
					loadCalls = append(loadCalls, struct {
						provider string
						apiKey   string
					}{provider: providerID, apiKey: apiKey})
					return []TranslationJobSetupProviderModelOptionReadModel{{ModelID: "gemini-2.5-pro", Label: "Gemini 2.5 Pro"}}, nil
				},
			),
		),
	)

	result, err := service.ListProviderModels(context.Background(), ListTranslationJobSetupProviderModelsRequest{
		PhaseID:          "word_translation",
		Provider:         "gemini",
		CredentialRef:    "gemini-primary",
		CredentialStatus: "configured",
		RequestToken:     "req-transport-hidden",
	})
	if err != nil {
		t.Fatalf("expected loader-based model list success: %v", err)
	}
	if result.Status != "success" || len(result.Models) != 1 {
		t.Fatalf("expected loader result mapping, got %#v", result)
	}
	if len(loadCalls) != 1 || loadCalls[0].provider != "gemini" {
		t.Fatalf("expected loader call with provider id, got %#v", loadCalls)
	}
}

func containsTranslationJobSetupString(values []string, target string) bool {
	for _, value := range values {
		if value == target {
			return true
		}
	}
	return false
}
