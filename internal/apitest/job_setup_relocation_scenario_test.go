package apitest

import (
	"context"
	"path/filepath"
	"testing"
	"time"

	wails "aitranslationenginejp/internal/controller/wails"
	"aitranslationenginejp/internal/infra/sqlite/dbinit"
	"aitranslationenginejp/internal/repository"
	"aitranslationenginejp/internal/service"
	"aitranslationenginejp/internal/usecase"

	"github.com/jmoiron/sqlx"
)

func TestSCN_JSR_001_InputReviewCreatesJobWithoutPhaseAISettings(t *testing.T) {
	harness := newJobSetupRelocationAPIHarness(t)
	input := harness.createInputWithTranslationField(t, "no-ai-settings")
	inputController := harness.newTranslationInputController()

	response, err := inputController.CreateTranslationJobFromInput(wails.TranslationInputCreateJobRequestDTO{
		InputID: input.ID,
	})
	if err != nil {
		t.Fatalf("SCN-JSR-001 expected input review job creation API to succeed: %v", err)
	}

	if !response.Accepted || response.JobID == nil || response.JobState != "ready" || response.CurrentPhase != "term_translation" {
		t.Fatalf("SCN-JSR-001 expected accepted ready job for term translation, got %#v", response)
	}
}

func TestSCN_JSR_002_TermPhaseStartRejectsMissingAISettingsThenStartsAfterSave(t *testing.T) {
	harness := newJobSetupRelocationAPIHarness(t)
	input := harness.createInputWithTranslationField(t, "term-start")
	inputController := harness.newTranslationInputController()
	termController := harness.newTermTranslationPhaseController()

	created, err := inputController.CreateTranslationJobFromInput(wails.TranslationInputCreateJobRequestDTO{
		InputID: input.ID,
	})
	if err != nil {
		t.Fatalf("SCN-JSR-002 expected input review job creation API to succeed: %v", err)
	}
	if created.JobID == nil {
		t.Fatalf("SCN-JSR-002 expected created job id, got %#v", created)
	}
	jobID := *created.JobID

	rejected, err := termController.StartTermTranslationPhase(wails.StartTermTranslationPhaseRequestDTO{
		JobID: jobID,
	})
	if err != nil {
		t.Fatalf("SCN-JSR-002 expected missing AI settings to return public rejection response: %v", err)
	}
	if rejected.ErrorSummary == nil || rejected.ErrorSummary.Reason != "phase runtime snapshot is missing" {
		t.Fatalf("SCN-JSR-002 expected missing phase AI settings rejection, got %#v", rejected)
	}

	saved, err := termController.SaveTermTranslationPhaseAISettings(wails.SaveTermTranslationPhaseAISettingsRequestDTO{
		Provider:      "gemini",
		Model:         "gemini-2.5-pro",
		ExecutionMode: "single_request",
		BatchMode:     "disabled",
	})
	if err != nil {
		t.Fatalf("SCN-JSR-002 expected term phase AI settings save to succeed: %v", err)
	}
	if saved.PhaseType != "word_translation" || saved.Provider != "gemini" {
		t.Fatalf("SCN-JSR-002 expected saved term phase AI settings, got %#v", saved)
	}

	started, err := termController.StartTermTranslationPhase(wails.StartTermTranslationPhaseRequestDTO{
		JobID: jobID,
	})
	if err != nil {
		t.Fatalf("SCN-JSR-002 expected configured term phase start to succeed: %v", err)
	}
	if started.ErrorSummary != nil || started.PhaseRunID == nil || started.PhaseState != "completed" {
		t.Fatalf("SCN-JSR-002 expected configured term phase to complete via fake provider, got %#v", started)
	}
}

type jobSetupRelocationAPIHarness struct {
	ctx         context.Context
	db          *sqlx.DB
	sourceRepo  repository.TranslationSourceRepository
	jobRepo     repository.JobLifecycleRepository
	foundation  repository.FoundationDataRepository
	transactor  repository.Transactor
	secretStore repository.ProviderSettingsSecretStore
}

func newJobSetupRelocationAPIHarness(t *testing.T) jobSetupRelocationAPIHarness {
	t.Helper()
	ctx := context.Background()
	db, err := dbinit.OpenMasterDictionaryDatabase(ctx, filepath.Join(t.TempDir(), "job-setup-relocation.sqlite3"), nil)
	if err != nil {
		t.Fatalf("open job setup relocation api database: %v", err)
	}
	t.Cleanup(func() { _ = db.Close() })

	secretStore := repository.NewFakeSecretStore()

	jobRepo := repository.NewSQLiteJobLifecycleRepository(db)
	return jobSetupRelocationAPIHarness{
		ctx:         ctx,
		db:          db,
		sourceRepo:  repository.NewSQLiteTranslationSourceRepository(db),
		jobRepo:     jobRepo,
		foundation:  repository.NewSQLiteFoundationDataRepository(db),
		transactor:  repository.NewSQLiteTransactor(db),
		secretStore: secretStore,
	}
}

func (h jobSetupRelocationAPIHarness) createInputWithTranslationField(t *testing.T, suffix string) repository.XEditExtractedData {
	t.Helper()
	input, err := h.sourceRepo.CreateXEditExtractedData(h.ctx, repository.XEditExtractedDataDraft{
		SourceFilePath:    "/tmp/job-setup-relocation-" + suffix + ".json",
		SourceContentHash: "job-setup-relocation-" + suffix,
		SourceTool:        "scenario-test",
		TargetPluginName:  "ScenarioPlugin.esm",
		TargetPluginType:  "esm",
		RecordCount:       1,
		ImportedAt:        time.Date(2026, 5, 22, 9, 0, 0, 0, time.UTC),
	})
	if err != nil {
		t.Fatalf("create scenario input: %v", err)
	}
	record, err := h.sourceRepo.CreateTranslationRecord(h.ctx, repository.TranslationRecordDraft{
		XEditExtractedDataID: input.ID,
		FormID:               "000JSR01",
		EditorID:             "SCN_JSR_" + suffix,
		RecordType:           "NPC_",
	})
	if err != nil {
		t.Fatalf("create scenario translation record: %v", err)
	}
	if _, err := h.sourceRepo.CreateTranslationField(h.ctx, repository.TranslationFieldDraft{
		TranslationRecordID: record.ID,
		SubrecordType:       "FULL",
		SourceText:          "Iron Sword " + suffix,
		FieldOrder:          0,
	}); err != nil {
		t.Fatalf("create scenario translation field: %v", err)
	}
	return input
}

func (h jobSetupRelocationAPIHarness) newTranslationInputController() *wails.TranslationInputController {
	setupService := service.NewPersistentTranslationJobSetupService(
		h.jobRepo,
		h.sourceRepo,
		jobSetupRelocationDictionaryRepository{},
		jobSetupRelocationPersonaRepository{},
		nil,
		h.secretStore,
		h.transactor,
	)
	setupUsecase := usecase.NewTranslationJobSetupUsecase(setupService)
	return wails.NewTranslationInputControllerWithJobCreate(jobSetupRelocationInputUsecase{}, setupUsecase)
}

func (h jobSetupRelocationAPIHarness) newTermTranslationPhaseController() *wails.TermTranslationPhaseController {
	termService := service.NewTermTranslationPhaseService(
		h.jobRepo,
		h.foundation,
		h.sourceRepo,
		h.transactor,
		h.secretStore,
		jobSetupRelocationTermProvider{},
	)
	termService.WithTermTranslationProviderSettings(jobSetupRelocationProviderSettings{})
	termService.WithTermTranslationJobPhaseAISettings(repository.NewSQLiteJobPhaseAISettingsRepository(h.db))
	termUsecase := usecase.NewTermTranslationPhaseUsecase(termService)
	return wails.NewTermTranslationPhaseController(termUsecase)
}

type jobSetupRelocationInputUsecase struct{}

func (jobSetupRelocationInputUsecase) ImportXEditJSON(context.Context, string) (usecase.TranslationInputImportResult, error) {
	return usecase.TranslationInputImportResult{}, nil
}

type jobSetupRelocationDictionaryRepository struct{}

func (jobSetupRelocationDictionaryRepository) List(context.Context, service.MasterDictionaryQuery) (service.MasterDictionaryListResult, error) {
	return service.MasterDictionaryListResult{}, nil
}

type jobSetupRelocationPersonaRepository struct{}

func (jobSetupRelocationPersonaRepository) List(context.Context, repository.MasterPersonaListQuery) (repository.MasterPersonaListResult, error) {
	return repository.MasterPersonaListResult{}, nil
}

type jobSetupRelocationProviderSettings struct{}

func (jobSetupRelocationProviderSettings) ListProviderSettings(context.Context) (service.ProviderSettingsRoute, []service.ProviderSettingsSummary, error) {
	return service.ProviderSettingsRoute{}, nil, nil
}

func (jobSetupRelocationProviderSettings) SaveProviderSettings(context.Context, service.ProviderSettingsSaveInput) (service.ProviderSettingsSummary, error) {
	return service.ProviderSettingsSummary{}, nil
}

func (jobSetupRelocationProviderSettings) ListProviderModels(context.Context, service.ProviderSettingsModelListInput) (service.ProviderSettingsModelListResult, error) {
	return service.ProviderSettingsModelListResult{}, nil
}

func (jobSetupRelocationProviderSettings) ResolveProviderExecutionSettings(_ context.Context, input service.ProviderSettingsResolveInput) (service.ProviderSettingsResolveResult, error) {
	endpoint := "https://scenario-provider.example.test"
	credentialRef := "provider-settings:gemini"
	return service.ProviderSettingsResolveResult{
		ConsumerID:            input.ConsumerID,
		ProviderID:            input.Selection.ProviderID,
		Model:                 input.Selection.Model,
		ExecutionMethod:       input.Selection.ExecutionMethod,
		UseBatchAPI:           input.Selection.UseBatchAPI,
		Endpoint:              &endpoint,
		CredentialReferenceID: &credentialRef,
		CredentialState:       "configured",
	}, nil
}

type jobSetupRelocationTermProvider struct{}

func (jobSetupRelocationTermProvider) TranslateTerm(_ context.Context, request service.TermTranslationProviderRequest) (service.TermTranslationProviderResult, error) {
	return service.TermTranslationProviderResult{
		SourceTerm:     request.SourceTerm,
		TranslatedTerm: "訳語: " + request.SourceTerm,
		Confirmed:      true,
		AuditSummary: service.TermTranslationProviderAuditSummary{
			Provider:      request.Provider,
			Model:         request.Model,
			ExecutionMode: "single_request",
			InputCount:    1,
			OutputCount:   1,
		},
	}, nil
}
