package usecase

import (
	"context"
	"reflect"
	"testing"
	"time"

	jobsetupservice "aitranslationenginejp/internal/service"
)

type fakeTranslationJobSetupService struct {
	validateRequestFunc    func(context.Context, jobsetupservice.TranslationJobSetupValidationRequest) (jobsetupservice.TranslationJobSetupValidationDecision, error)
	evaluateCreateFunc     func(context.Context, jobsetupservice.TranslationJobSetupCreateRequest) (jobsetupservice.TranslationJobSetupCreateDecision, error)
	createTranslationFunc  func(context.Context, jobsetupservice.TranslationJobSetupCreateRequest, []string) (jobsetupservice.TranslationJobSetupCreatedJobReadModel, error)
	readOptionsFunc        func(context.Context) (jobsetupservice.TranslationJobSetupOptionsReadModel, error)
	readSummaryFunc        func(context.Context, int64) (jobsetupservice.TranslationJobSetupSummaryReadModel, error)
	listProviderModelsFunc func(context.Context, jobsetupservice.ListTranslationJobSetupProviderModelsRequest) (jobsetupservice.ListTranslationJobSetupProviderModelsResult, error)
}

func (service fakeTranslationJobSetupService) ValidateRequest(ctx context.Context, request jobsetupservice.TranslationJobSetupValidationRequest) (jobsetupservice.TranslationJobSetupValidationDecision, error) {
	return service.validateRequestFunc(ctx, request)
}

func (service fakeTranslationJobSetupService) EvaluateCreateRequest(ctx context.Context, request jobsetupservice.TranslationJobSetupCreateRequest) (jobsetupservice.TranslationJobSetupCreateDecision, error) {
	return service.evaluateCreateFunc(ctx, request)
}

func (service fakeTranslationJobSetupService) CreateTranslationJob(ctx context.Context, request jobsetupservice.TranslationJobSetupCreateRequest, passSlices []string) (jobsetupservice.TranslationJobSetupCreatedJobReadModel, error) {
	return service.createTranslationFunc(ctx, request, passSlices)
}

func (service fakeTranslationJobSetupService) ReadOptions(ctx context.Context) (jobsetupservice.TranslationJobSetupOptionsReadModel, error) {
	return service.readOptionsFunc(ctx)
}

func (service fakeTranslationJobSetupService) ReadSummary(ctx context.Context, jobID int64) (jobsetupservice.TranslationJobSetupSummaryReadModel, error) {
	return service.readSummaryFunc(ctx, jobID)
}

func (service fakeTranslationJobSetupService) ListProviderModels(ctx context.Context, request jobsetupservice.ListTranslationJobSetupProviderModelsRequest) (jobsetupservice.ListTranslationJobSetupProviderModelsResult, error) {
	return service.listProviderModelsFunc(ctx, request)
}

func TestTranslationJobSetupUsecaseValidateForwardsPhaseRuntimes(t *testing.T) {
	validatedAt := time.Date(2026, 5, 4, 11, 0, 0, 0, time.UTC)
	var captured jobsetupservice.TranslationJobSetupValidationRequest
	usecase := NewTranslationJobSetupUsecase(fakeTranslationJobSetupService{
		validateRequestFunc: func(_ context.Context, request jobsetupservice.TranslationJobSetupValidationRequest) (jobsetupservice.TranslationJobSetupValidationDecision, error) {
			captured = request
			return jobsetupservice.TranslationJobSetupValidationDecision{
				Status:      "pass",
				ValidatedAt: validatedAt,
				CanCreate:   true,
				PassSlices:  []string{"input", "runtime", "credentials"},
				PhaseResults: []jobsetupservice.TranslationJobSetupPhaseValidationReadModel{{
					PhaseID:              "word_translation",
					Status:               "pass",
					CanCreate:            true,
					ModelListState:       "success",
					ModelListSourceToken: "word_translation|openai|openai-primary|req-1",
				}},
			}, nil
		},
	})

	got, err := usecase.ValidateTranslationJobSetup(context.Background(), ValidateTranslationJobSetupRequest{
		InputSourceID: 44,
		PhaseRuntimeSelections: []TranslationJobSetupPhaseRuntimeSelection{{
			PhaseID:          "word_translation",
			Provider:         "openai",
			Model:            "gpt-5.4-mini",
			CredentialStatus: "configured",
			ExecutionMode:    "sync",
			BatchMode:        "unsupported",
		}},
	})
	if err != nil {
		t.Fatalf("expected validation success: %v", err)
	}
	if captured.InputSourceID != 44 || len(captured.PhaseRuntimes) != 1 || captured.PhaseRuntimes[0].PhaseID != "word_translation" {
		t.Fatalf("expected phase runtime forwarding, got %#v", captured)
	}
	if got.Status != "pass" || !got.CanCreate || len(got.PhaseResults) != 1 {
		t.Fatalf("expected mapped validation result, got %#v", got)
	}
}

func TestTranslationJobSetupUsecaseListProviderModelsMapsResponse(t *testing.T) {
	usecase := NewTranslationJobSetupUsecase(fakeTranslationJobSetupService{
		listProviderModelsFunc: func(_ context.Context, request jobsetupservice.ListTranslationJobSetupProviderModelsRequest) (jobsetupservice.ListTranslationJobSetupProviderModelsResult, error) {
			if request.Provider != "gemini" {
				t.Fatalf("expected gemini request, got %#v", request)
			}
			return jobsetupservice.ListTranslationJobSetupProviderModelsResult{
				PhaseID:          request.PhaseID,
				Provider:         request.Provider,
				CredentialStatus: "configured",
				RequestToken:     request.RequestToken,
				SourceToken:      "npc_persona_generation|gemini|gemini-primary|req-2",
				Status:           "success",
				Models:           []jobsetupservice.TranslationJobSetupProviderModelOptionReadModel{{ModelID: "gemini-2.5-pro", Label: "Gemini 2.5 Pro"}},
			}, nil
		},
	})

	got, err := usecase.ListTranslationJobSetupProviderModels(context.Background(), ListTranslationJobSetupProviderModelsRequest{
		PhaseID:          "npc_persona_generation",
		Provider:         "gemini",
		CredentialStatus: "configured",
		RequestToken:     "req-2",
	})
	if err != nil {
		t.Fatalf("expected provider model list success: %v", err)
	}
	want := ListTranslationJobSetupProviderModelsResult{
		PhaseID:          "npc_persona_generation",
		Provider:         "gemini",
		CredentialStatus: "configured",
		RequestToken:     "req-2",
		SourceToken:      "npc_persona_generation|gemini|gemini-primary|req-2",
		Status:           "success",
		Models:           []TranslationJobSetupProviderModelOption{{ModelID: "gemini-2.5-pro", Label: "Gemini 2.5 Pro"}},
	}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("expected mapped provider model response %#v, got %#v", want, got)
	}
}

func TestTranslationJobSetupUsecaseCreateReturnsPhaseRuntimeSummaries(t *testing.T) {
	usecase := NewTranslationJobSetupUsecase(fakeTranslationJobSetupService{
		evaluateCreateFunc: func(_ context.Context, request jobsetupservice.TranslationJobSetupCreateRequest) (jobsetupservice.TranslationJobSetupCreateDecision, error) {
			if len(request.PhaseRuntimes) != 3 {
				t.Fatalf("expected three phase runtimes, got %#v", request.PhaseRuntimes)
			}
			return jobsetupservice.TranslationJobSetupCreateDecision{
				CanCreate:            true,
				ValidationPassSlices: []string{"input", "runtime", "credentials"},
			}, nil
		},
		createTranslationFunc: func(_ context.Context, _ jobsetupservice.TranslationJobSetupCreateRequest, passSlices []string) (jobsetupservice.TranslationJobSetupCreatedJobReadModel, error) {
			return jobsetupservice.TranslationJobSetupCreatedJobReadModel{
				JobID:                77,
				JobState:             "ready",
				InputSource:          "translation_input",
				ExecutionSummary:     jobsetupservice.TranslationJobSetupExecutionSummaryReadModel{Provider: "openai", Model: "gpt-5.4-mini", ExecutionMode: "sync"},
				ValidationPassSlices: passSlices,
				PhaseRuntimeSummaries: []jobsetupservice.TranslationJobSetupPhaseRuntimeSummaryReadModel{{
					PhaseID:              "text_translation",
					Provider:             "xai",
					Model:                "grok-4",
					CredentialRef:        "xai-primary",
					CredentialStatus:     "configured",
					ExecutionMode:        "batch",
					BatchMode:            "enabled",
					ModelListSourceToken: "text_translation|xai|xai-primary|req-3",
				}},
			}, nil
		},
	})

	got, err := usecase.CreateTranslationJob(context.Background(), CreateTranslationJobRequest{
		InputSourceID:    44,
		ValidationStatus: "pass",
		ValidatedAt:      time.Date(2026, 5, 4, 12, 0, 0, 0, time.UTC),
		PhaseRuntimeSelections: []TranslationJobSetupPhaseRuntimeSelection{
			{PhaseID: "word_translation", Provider: "openai", Model: "gpt-5.4-mini", CredentialStatus: "configured", ExecutionMode: "sync", BatchMode: "unsupported"},
			{PhaseID: "npc_persona_generation", Provider: "gemini", Model: "gemini-2.5-pro", CredentialStatus: "configured", ExecutionMode: "sync", BatchMode: "disabled"},
			{PhaseID: "text_translation", Provider: "xai", Model: "grok-4", CredentialStatus: "configured", ExecutionMode: "batch", BatchMode: "enabled"},
		},
	})
	if err != nil {
		t.Fatalf("expected create success: %v", err)
	}
	if got.JobID != 77 || len(got.PhaseRuntimeSummaries) != 1 || got.PhaseRuntimeSummaries[0].Provider != "xai" {
		t.Fatalf("expected phase runtime summaries in create result, got %#v", got)
	}
}
