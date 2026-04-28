package usecase

import (
	"context"
	"reflect"
	"testing"
	"time"

	jobsetupservice "aitranslationenginejp/internal/service"
)

type fakeTranslationJobSetupService struct {
	validateRequestFunc       func(context.Context, jobsetupservice.TranslationJobSetupValidationRequest) (jobsetupservice.TranslationJobSetupValidationDecision, error)
	evaluateCreateRequestFunc func(context.Context, jobsetupservice.TranslationJobSetupCreateRequest) (jobsetupservice.TranslationJobSetupCreateDecision, error)
	readOptionsFunc           func(context.Context) (jobsetupservice.TranslationJobSetupOptionsReadModel, error)
}

func (service fakeTranslationJobSetupService) ValidateRequest(
	ctx context.Context,
	request jobsetupservice.TranslationJobSetupValidationRequest,
) (jobsetupservice.TranslationJobSetupValidationDecision, error) {
	if service.validateRequestFunc != nil {
		return service.validateRequestFunc(ctx, request)
	}
	return jobsetupservice.TranslationJobSetupValidationDecision{}, nil
}

func (service fakeTranslationJobSetupService) EvaluateCreateRequest(
	ctx context.Context,
	request jobsetupservice.TranslationJobSetupCreateRequest,
) (jobsetupservice.TranslationJobSetupCreateDecision, error) {
	if service.evaluateCreateRequestFunc != nil {
		return service.evaluateCreateRequestFunc(ctx, request)
	}
	return jobsetupservice.TranslationJobSetupCreateDecision{}, nil
}

func (service fakeTranslationJobSetupService) ReadOptions(
	ctx context.Context,
) (jobsetupservice.TranslationJobSetupOptionsReadModel, error) {
	if service.readOptionsFunc != nil {
		return service.readOptionsFunc(ctx)
	}
	return jobsetupservice.TranslationJobSetupOptionsReadModel{}, nil
}

func TestTranslationJobSetupUsecaseValidateTranslationJobSetupReturnsPassDecision(t *testing.T) {
	validatedAt := time.Date(2026, 4, 27, 13, 0, 0, 0, time.UTC)
	var captured jobsetupservice.TranslationJobSetupValidationRequest
	usecase := NewTranslationJobSetupUsecase(fakeTranslationJobSetupService{
		validateRequestFunc: func(
			_ context.Context,
			request jobsetupservice.TranslationJobSetupValidationRequest,
		) (jobsetupservice.TranslationJobSetupValidationDecision, error) {
			captured = request
			return jobsetupservice.TranslationJobSetupValidationDecision{
				Status:      "pass",
				ValidatedAt: validatedAt,
				CanCreate:   true,
				PassSlices:  []string{"input", "runtime", "credentials"},
			}, nil
		},
	})

	got, err := usecase.ValidateTranslationJobSetup(context.Background(), ValidateTranslationJobSetupRequest{
		InputSourceID: 44,
		Runtime: TranslationJobSetupRuntimeSelection{
			Provider:      "openai",
			Model:         "gpt-5.4-mini",
			ExecutionMode: "batch",
		},
		CredentialRef: "openai-primary",
	})
	if err != nil {
		t.Fatalf("expected validation request to succeed: %v", err)
	}

	wantCaptured := jobsetupservice.TranslationJobSetupValidationRequest{
		InputSourceID: 44,
		Provider:      "openai",
		Model:         "gpt-5.4-mini",
		ExecutionMode: "batch",
		CredentialRef: "openai-primary",
	}
	if !reflect.DeepEqual(captured, wantCaptured) {
		t.Fatalf("expected validation request %#v, got %#v", wantCaptured, captured)
	}

	want := TranslationJobSetupValidationResult{
		Status:      "pass",
		ValidatedAt: validatedAt,
		CanCreate:   true,
		PassSlices:  []string{"input", "runtime", "credentials"},
	}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("expected validation result %#v, got %#v", want, got)
	}
}

func TestTranslationJobSetupUsecaseGetOptionsMapsServerOwnedReadModelWithoutExistingJob(t *testing.T) {
	usecase := NewTranslationJobSetupUsecase(fakeTranslationJobSetupService{
		readOptionsFunc: func(context.Context) (jobsetupservice.TranslationJobSetupOptionsReadModel, error) {
			return jobsetupservice.TranslationJobSetupOptionsReadModel{
				InputCandidates: []jobsetupservice.TranslationJobSetupInputCandidateReadModel{
					{
						ID:          77,
						Label:       "Quest fragments import",
						SourceKind:  "translation_input",
						RecordCount: 32,
					},
				},
				ExistingJob: nil,
				SharedDictionaries: []jobsetupservice.TranslationJobSetupDictionaryOptionReadModel{
					{ID: "dict-quest", Label: "Quest Dictionary"},
				},
				SharedPersonas: []jobsetupservice.TranslationJobSetupPersonaOptionReadModel{
					{ID: "persona-mage", Label: "Mage Persona"},
				},
				AIRuntimeOptions: []jobsetupservice.TranslationJobSetupRuntimeOptionReadModel{
					{Provider: "xai", Model: "grok-4", Mode: "sync"},
					{Provider: "lm_studio", Model: "lmstudio-community", Mode: "sync"},
				},
				CredentialRefs: []jobsetupservice.TranslationJobSetupCredentialReferenceReadModel{
					{Provider: "xai", CredentialRef: "xai-primary", IsConfigured: true, IsMissingSecret: false},
					{Provider: "lm_studio", CredentialRef: "lmstudio-local", IsConfigured: true, IsMissingSecret: false},
				},
			}, nil
		},
	})

	got, err := usecase.GetTranslationJobSetupOptions(context.Background())
	if err != nil {
		t.Fatalf("expected options request to succeed: %v", err)
	}

	want := TranslationJobSetupOptionsResult{
		InputCandidates: []TranslationJobSetupInputCandidate{
			{ID: 77, Label: "Quest fragments import", SourceKind: "translation_input", RecordCount: 32},
		},
		ExistingJob: nil,
		SharedDictionaries: []TranslationJobSetupDictionaryOption{
			{ID: "dict-quest", Label: "Quest Dictionary"},
		},
		SharedPersonas: []TranslationJobSetupPersonaOption{
			{ID: "persona-mage", Label: "Mage Persona"},
		},
		AIRuntimeOptions: []TranslationJobSetupRuntimeOption{
			{Provider: "xai", Model: "grok-4", Mode: "sync"},
			{Provider: "lm_studio", Model: "lmstudio-community", Mode: "sync"},
		},
		CredentialRefs: []TranslationJobSetupCredentialReference{
			{Provider: "xai", CredentialRef: "xai-primary", IsConfigured: true, IsMissingSecret: false},
			{Provider: "lm_studio", CredentialRef: "lmstudio-local", IsConfigured: true, IsMissingSecret: false},
		},
	}

	if !reflect.DeepEqual(got, want) {
		t.Fatalf("expected assembled job setup options %#v, got %#v", want, got)
	}
}

func TestTranslationJobSetupUsecaseGetOptionsMapsExistingJobInputSourceID(t *testing.T) {
	usecase := NewTranslationJobSetupUsecase(fakeTranslationJobSetupService{
		readOptionsFunc: func(context.Context) (jobsetupservice.TranslationJobSetupOptionsReadModel, error) {
			return jobsetupservice.TranslationJobSetupOptionsReadModel{
				InputCandidates: []jobsetupservice.TranslationJobSetupInputCandidateReadModel{{
					ID:          77,
					Label:       "Quest fragments import",
					SourceKind:  "translation_input",
					RecordCount: 32,
				}},
				ExistingJob: &jobsetupservice.TranslationJobSetupExistingJobReadModel{
					InputSourceID: 77,
					JobID:         91,
					Status:        "ready",
					InputSource:   "translation_input",
				},
			}, nil
		},
	})

	got, err := usecase.GetTranslationJobSetupOptions(context.Background())
	if err != nil {
		t.Fatalf("expected options request to succeed: %v", err)
	}
	if got.ExistingJob == nil {
		t.Fatalf("expected existing job to be mapped, got %#v", got)
	}
	if got.ExistingJob.InputSourceID != 77 {
		t.Fatalf("expected existing job inputSourceId=77, got %#v", got.ExistingJob)
	}
	if got.ExistingJob.JobID != 91 || got.ExistingJob.Status != "ready" || got.ExistingJob.InputSource != "translation_input" {
		t.Fatalf("expected existing job fields to stay intact, got %#v", got.ExistingJob)
	}
}

func TestTranslationJobSetupUsecaseGetSummaryReturnsReadOnlyJob(t *testing.T) {
	usecase := NewTranslationJobSetupUsecase(fakeTranslationJobSetupService{})

	got, err := usecase.GetTranslationJobSetupSummary(context.Background(), GetTranslationJobSetupSummaryRequest{JobID: 91})
	if err != nil {
		t.Fatalf("expected summary request to succeed: %v", err)
	}

	want := TranslationJobSetupSummaryResult{
		JobID:       91,
		JobState:    "ready",
		InputSource: "translation_input",
		ExecutionSummary: TranslationJobExecutionSummary{
			Provider:      "openai",
			Model:         "gpt-5.4-mini",
			ExecutionMode: "batch",
		},
		ValidationPassSlices: []string{"input", "runtime", "credentials"},
	}

	if !reflect.DeepEqual(got, want) {
		t.Fatalf("expected read-only job setup summary %#v, got %#v", want, got)
	}
}

func TestTranslationJobSetupUsecaseCreateTranslationJobForwardsFullCreateRequest(t *testing.T) {
	validatedAt := time.Date(2026, 4, 27, 13, 30, 0, 0, time.UTC)
	var captured jobsetupservice.TranslationJobSetupCreateRequest
	usecase := NewTranslationJobSetupUsecase(fakeTranslationJobSetupService{
		evaluateCreateRequestFunc: func(
			_ context.Context,
			request jobsetupservice.TranslationJobSetupCreateRequest,
		) (jobsetupservice.TranslationJobSetupCreateDecision, error) {
			captured = request
			return jobsetupservice.TranslationJobSetupCreateDecision{
				CanCreate: false,
				ErrorKind: TranslationJobSetupErrorKindValidationFailed,
			}, nil
		},
	})

	got, err := usecase.CreateTranslationJob(context.Background(), CreateTranslationJobRequest{
		InputSourceID:        44,
		InputSource:          "translation_input",
		ValidationStatus:     TranslationJobSetupValidationStatusPass,
		ValidatedAt:          validatedAt,
		ValidationPassSlices: []string{"input", "runtime"},
		Runtime: TranslationJobSetupRuntimeSelection{
			Provider:      "openai",
			Model:         "gpt-5.4-mini",
			ExecutionMode: "batch",
		},
		CredentialRef: "openai-primary",
	})
	if err != nil {
		t.Fatalf("expected rejected create to return without transport error: %v", err)
	}

	wantCaptured := jobsetupservice.TranslationJobSetupCreateRequest{
		InputSourceID:    44,
		ValidationStatus: TranslationJobSetupValidationStatusPass,
		ValidatedAt:      validatedAt,
		Provider:         "openai",
		Model:            "gpt-5.4-mini",
		ExecutionMode:    "batch",
		CredentialRef:    "openai-primary",
	}
	if !reflect.DeepEqual(captured, wantCaptured) {
		t.Fatalf("expected create request %#v, got %#v", wantCaptured, captured)
	}
	if got.ErrorKind != TranslationJobSetupErrorKindValidationFailed {
		t.Fatalf("expected validation_failed result, got %#v", got)
	}
}

func TestTranslationJobSetupUsecaseCreateTranslationJobReturnsRejectedDecision(t *testing.T) {
	usecase := NewTranslationJobSetupUsecase(fakeTranslationJobSetupService{
		evaluateCreateRequestFunc: func(
			_ context.Context,
			_ jobsetupservice.TranslationJobSetupCreateRequest,
		) (jobsetupservice.TranslationJobSetupCreateDecision, error) {
			return jobsetupservice.TranslationJobSetupCreateDecision{
				CanCreate: false,
				ErrorKind: TranslationJobSetupErrorKindValidationStale,
			}, nil
		},
	})

	got, err := usecase.CreateTranslationJob(context.Background(), CreateTranslationJobRequest{
		ValidationStatus: TranslationJobSetupValidationStatusPass,
	})
	if err != nil {
		t.Fatalf("expected rejected create to return without transport error: %v", err)
	}
	if got.ErrorKind != TranslationJobSetupErrorKindValidationStale {
		t.Fatalf("expected validation_stale result, got %#v", got)
	}
}

func TestTranslationJobSetupUsecaseCreateTranslationJobPreservesProviderModeUnsupportedKind(t *testing.T) {
	usecase := NewTranslationJobSetupUsecase(fakeTranslationJobSetupService{
		evaluateCreateRequestFunc: func(
			_ context.Context,
			_ jobsetupservice.TranslationJobSetupCreateRequest,
		) (jobsetupservice.TranslationJobSetupCreateDecision, error) {
			return jobsetupservice.TranslationJobSetupCreateDecision{
				CanCreate: false,
				ErrorKind: TranslationJobSetupErrorKindProviderModeUnsupported,
			}, nil
		},
	})

	got, err := usecase.CreateTranslationJob(context.Background(), CreateTranslationJobRequest{
		ValidationStatus: TranslationJobSetupValidationStatusPass,
	})
	if err != nil {
		t.Fatalf("expected rejected create to preserve provider_mode_unsupported without transport error: %v", err)
	}
	if got.ErrorKind != TranslationJobSetupErrorKindProviderModeUnsupported {
		t.Fatalf("expected provider_mode_unsupported result, got %#v", got)
	}
}
