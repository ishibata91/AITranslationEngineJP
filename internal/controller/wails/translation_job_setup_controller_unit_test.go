package wails

import (
	"context"
	"errors"
	"reflect"
	"slices"
	"testing"
	"time"

	"aitranslationenginejp/internal/usecase"
)

const (
	translationJobSetupRealProviderGemini   = "gemini"
	translationJobSetupRealProviderLMStudio = "lm_studio"
	translationJobSetupRealProviderXAI      = "xai"
)

type fakeTranslationJobSetupUsecase struct {
	getOptionsFunc func(ctx context.Context) (usecase.TranslationJobSetupOptionsResult, error)
	validateFunc   func(ctx context.Context, request usecase.ValidateTranslationJobSetupRequest) (usecase.TranslationJobSetupValidationResult, error)
	createJobFunc  func(ctx context.Context, request usecase.CreateTranslationJobRequest) (usecase.CreateTranslationJobResult, error)
	getSummaryFunc func(ctx context.Context, request usecase.GetTranslationJobSetupSummaryRequest) (usecase.TranslationJobSetupSummaryResult, error)
}

func (fake fakeTranslationJobSetupUsecase) GetTranslationJobSetupOptions(ctx context.Context) (usecase.TranslationJobSetupOptionsResult, error) {
	if fake.getOptionsFunc == nil {
		return usecase.TranslationJobSetupOptionsResult{}, nil
	}
	return fake.getOptionsFunc(ctx)
}

func (fake fakeTranslationJobSetupUsecase) ValidateTranslationJobSetup(ctx context.Context, request usecase.ValidateTranslationJobSetupRequest) (usecase.TranslationJobSetupValidationResult, error) {
	if fake.validateFunc == nil {
		return usecase.TranslationJobSetupValidationResult{}, nil
	}
	return fake.validateFunc(ctx, request)
}

func (fake fakeTranslationJobSetupUsecase) CreateTranslationJob(ctx context.Context, request usecase.CreateTranslationJobRequest) (usecase.CreateTranslationJobResult, error) {
	if fake.createJobFunc == nil {
		return usecase.CreateTranslationJobResult{}, nil
	}
	return fake.createJobFunc(ctx, request)
}

func (fake fakeTranslationJobSetupUsecase) GetTranslationJobSetupSummary(ctx context.Context, request usecase.GetTranslationJobSetupSummaryRequest) (usecase.TranslationJobSetupSummaryResult, error) {
	if fake.getSummaryFunc == nil {
		return usecase.TranslationJobSetupSummaryResult{}, nil
	}
	return fake.getSummaryFunc(ctx, request)
}

func TestTranslationJobSetupControllerGetOptionsMapsOptionsContract(t *testing.T) {
	controller := NewTranslationJobSetupController(fakeTranslationJobSetupUsecase{
		getOptionsFunc: func(ctx context.Context) (usecase.TranslationJobSetupOptionsResult, error) {
			if ctx == nil {
				t.Fatal("expected request context")
			}
			return usecase.TranslationJobSetupOptionsResult{
				InputCandidates: []usecase.TranslationJobSetupInputCandidate{
					{ID: 44, Label: "Dialogues import", SourceKind: "translation_input", RecordCount: 120},
				},
				ExistingJob: &usecase.TranslationJobSetupExistingJob{
					JobID:         88,
					Status:        "ready",
					InputSourceID: 44,
					InputSource:   "translation_input",
				},
				SharedDictionaries: []usecase.TranslationJobSetupDictionaryOption{
					{ID: "dict-core", Label: "Core Dictionary"},
				},
				SharedPersonas: []usecase.TranslationJobSetupPersonaOption{
					{ID: "persona-guard", Label: "Guard Persona"},
				},
				AIRuntimeOptions: []usecase.TranslationJobSetupRuntimeOption{
					{Provider: "openai", Model: "gpt-5.4-mini", Mode: "batch"},
					{Provider: "gemini", Model: "gemini-2.5-pro", Mode: "sync"},
				},
				CredentialRefs: []usecase.TranslationJobSetupCredentialReference{
					{Provider: "openai", CredentialRef: "openai-primary", IsConfigured: true, IsMissingSecret: false},
					{Provider: "gemini", CredentialRef: "gemini-missing", IsConfigured: false, IsMissingSecret: true},
				},
			}, nil
		},
	})
	response, err := controller.GetTranslationJobSetupOptions()
	if err != nil {
		t.Fatalf("expected options request to succeed: %v", err)
	}
	assertTranslationJobSetupOptionsContractResponse(t, response)
}

func TestTranslationJobSetupControllerGetOptionsIncludesInputCandidateRegisteredAtContract(t *testing.T) {
	field, ok := reflect.TypeOf(TranslationJobSetupInputCandidateDTO{}).FieldByName("RegisteredAt")
	if !ok {
		t.Fatalf("expected input candidate contract to expose RegisteredAt")
	}
	if field.Type != reflect.TypeOf("") {
		t.Fatalf("expected RegisteredAt to be string, got %v", field.Type)
	}
	if got := field.Tag.Get("json"); got != "registeredAt" {
		t.Fatalf("expected RegisteredAt json tag %q, got %q", "registeredAt", got)
	}
}

func TestTranslationJobSetupControllerGetOptionsIncludesExistingJobInputSourceIDContract(t *testing.T) {
	field, ok := reflect.TypeOf(TranslationJobSetupExistingJobDTO{}).FieldByName("InputSourceID")
	if !ok {
		t.Fatalf("expected existing job contract to expose InputSourceID")
	}
	if field.Type != reflect.TypeOf(int64(0)) {
		t.Fatalf("expected InputSourceID to be int64, got %v", field.Type)
	}
	if got := field.Tag.Get("json"); got != "inputSourceId" {
		t.Fatalf("expected InputSourceID json tag %q, got %q", "inputSourceId", got)
	}
}

func TestTranslationJobSetupControllerGetOptionsKeepsRealProviderListUserVisible(t *testing.T) {
	controller := NewTranslationJobSetupController(fakeTranslationJobSetupUsecase{
		getOptionsFunc: func(context.Context) (usecase.TranslationJobSetupOptionsResult, error) {
			return usecase.TranslationJobSetupOptionsResult{
				AIRuntimeOptions: []usecase.TranslationJobSetupRuntimeOption{
					{Provider: translationJobSetupRealProviderGemini, Model: "gemini-2.5-pro", Mode: "sync"},
					{Provider: translationJobSetupRealProviderLMStudio, Model: "local-model", Mode: "sync"},
					{Provider: translationJobSetupRealProviderXAI, Model: "grok-2", Mode: "sync"},
				},
			}, nil
		},
	})

	response, err := controller.GetTranslationJobSetupOptions()
	if err != nil {
		t.Fatalf("expected options request to succeed: %v", err)
	}
	providers := make([]string, 0, len(response.AIRuntimeOptions))
	for _, option := range response.AIRuntimeOptions {
		providers = append(providers, option.Provider)
	}
	want := []string{
		translationJobSetupRealProviderGemini,
		translationJobSetupRealProviderLMStudio,
		translationJobSetupRealProviderXAI,
	}
	if !slices.Equal(providers, want) {
		t.Fatalf("expected real provider list %#v, got %#v", want, providers)
	}
	if slices.Contains(providers, "fake-provider") {
		t.Fatalf("expected user-facing provider list to exclude fake providers, got %#v", providers)
	}
}

func TestTranslationJobSetupControllerValidateMapsValidationShape(t *testing.T) {
	validatedAt := time.Date(2026, 4, 27, 9, 30, 0, 0, time.UTC)
	controller := NewTranslationJobSetupController(fakeTranslationJobSetupUsecase{
		validateFunc: func(ctx context.Context, request usecase.ValidateTranslationJobSetupRequest) (usecase.TranslationJobSetupValidationResult, error) {
			if ctx == nil {
				t.Fatal("expected request context")
			}
			if request.InputSourceID != 44 {
				t.Fatalf("expected input source id to be forwarded, got %#v", request)
			}
			return usecase.TranslationJobSetupValidationResult{
				Status:                  "warning",
				BlockingFailureCategory: nil,
				TargetSlices:            []string{"input", "runtime", "credentials"},
				ValidatedAt:             validatedAt,
				CanCreate:               true,
				PassSlices:              []string{"input"},
			}, nil
		},
	})

	response, err := controller.ValidateTranslationJobSetup(ValidateTranslationJobSetupRequestDTO{
		InputSourceID: 44,
		Runtime: TranslationJobSetupRuntimeSelectionDTO{
			Provider:      "openai",
			Model:         "gpt-5.4-mini",
			ExecutionMode: "batch",
		},
		CredentialRef: "openai-primary",
	})
	if err != nil {
		t.Fatalf("expected validation request to succeed: %v", err)
	}
	if response.Status != "warning" || !response.CanCreate {
		t.Fatalf("expected warning validation response with create allowance, got %#v", response)
	}
	if response.ValidatedAt != validatedAt.Format(time.RFC3339) {
		t.Fatalf("expected RFC3339 validation timestamp, got %q", response.ValidatedAt)
	}
	if response.BlockingFailureCategory != nil {
		t.Fatalf("expected nil blocking failure category, got %#v", response.BlockingFailureCategory)
	}
	if len(response.TargetSlices) != 3 || response.TargetSlices[2] != "credentials" {
		t.Fatalf("expected target slices in response, got %#v", response.TargetSlices)
	}
	if len(response.PassSlices) != 1 || response.PassSlices[0] != "input" {
		t.Fatalf("expected pass slices in response, got %#v", response.PassSlices)
	}
}

func TestTranslationJobSetupControllerValidateReturnsBlockingFailureContract(t *testing.T) {
	testCases := []struct {
		name                 string
		request              ValidateTranslationJobSetupRequestDTO
		expectedCategory     string
		expectedTargetSlices []string
	}{
		{
			name: "必須設定不足は blocking failure として返す",
			request: ValidateTranslationJobSetupRequestDTO{
				InputSourceID: 0,
				Runtime: TranslationJobSetupRuntimeSelectionDTO{
					Provider:      "",
					Model:         "",
					ExecutionMode: "",
				},
				CredentialRef: "",
			},
			expectedCategory:     "required_setting_missing",
			expectedTargetSlices: []string{"input", "runtime", "credentials"},
		},
		{
			name: "共通基盤参照不能は blocking failure として返す",
			request: ValidateTranslationJobSetupRequestDTO{
				InputSourceID: 44,
				Runtime: TranslationJobSetupRuntimeSelectionDTO{
					Provider:      "openai",
					Model:         "gpt-5.4-mini",
					ExecutionMode: "batch",
				},
				CredentialRef: "foundation-ref-missing",
			},
			expectedCategory:     "foundation_ref_missing",
			expectedTargetSlices: []string{"foundation"},
		},
		{
			name: "provider と mode の不整合は blocking failure として返す",
			request: ValidateTranslationJobSetupRequestDTO{
				InputSourceID: 44,
				Runtime: TranslationJobSetupRuntimeSelectionDTO{
					Provider:      "lmstudio",
					Model:         "local-model",
					ExecutionMode: "batch",
				},
				CredentialRef: "lmstudio-local",
			},
			expectedCategory:     "provider_mode_unsupported",
			expectedTargetSlices: []string{"runtime"},
		},
		{
			name: "credential 参照不能は blocking failure として返す",
			request: ValidateTranslationJobSetupRequestDTO{
				InputSourceID: 44,
				Runtime: TranslationJobSetupRuntimeSelectionDTO{
					Provider:      "openai",
					Model:         "gpt-5.4-mini",
					ExecutionMode: "batch",
				},
				CredentialRef: "missing-credential-ref",
			},
			expectedCategory:     "credential_missing",
			expectedTargetSlices: []string{"credentials"},
		},
		{
			name: "cache 欠落は blocking failure として返す",
			request: ValidateTranslationJobSetupRequestDTO{
				InputSourceID: 999,
				Runtime: TranslationJobSetupRuntimeSelectionDTO{
					Provider:      "openai",
					Model:         "gpt-5.4-mini",
					ExecutionMode: "batch",
				},
				CredentialRef: "openai-primary",
			},
			expectedCategory:     "cache_missing",
			expectedTargetSlices: []string{"input"},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			controller := NewTranslationJobSetupController(usecase.NewTranslationJobSetupContractStub())

			response, err := controller.ValidateTranslationJobSetup(testCase.request)
			if err != nil {
				t.Fatalf("expected blocking validation response without transport error: %v", err)
			}
			assertBlockingFailureValidationResponse(t, response, testCase.expectedCategory, testCase.expectedTargetSlices)
		})
	}
}

func TestTranslationJobSetupControllerCreateMapsReadyJobResponse(t *testing.T) {
	validatedAt := time.Date(2026, 4, 27, 9, 45, 0, 0, time.UTC)
	controller := NewTranslationJobSetupController(fakeTranslationJobSetupUsecase{
		createJobFunc: func(ctx context.Context, request usecase.CreateTranslationJobRequest) (usecase.CreateTranslationJobResult, error) {
			if ctx == nil {
				t.Fatal("expected request context")
			}
			assertCreateTranslationJobRequestReadyContract(t, request, validatedAt)
			return usecase.CreateTranslationJobResult{
				JobID:       91,
				JobState:    "ready",
				InputSource: "translation_input",
				ExecutionSummary: usecase.TranslationJobExecutionSummary{
					Provider:      "openai",
					Model:         "gpt-5.4-mini",
					ExecutionMode: "batch",
				},
				ValidationPassSlices: []string{"input", "runtime"},
			}, nil
		},
	})

	request := CreateTranslationJobRequestDTO{
		InputSourceID:        44,
		InputSource:          "translation_input",
		ValidationStatus:     "pass",
		ValidationPassSlices: []string{"input", "runtime"},
		Runtime: TranslationJobSetupRuntimeSelectionDTO{
			Provider:      "openai",
			Model:         "gpt-5.4-mini",
			ExecutionMode: "batch",
		},
		CredentialRef: "openai-primary",
	}
	setTranslationJobSetupFreshness(t, &request, validatedAt)

	response, err := controller.CreateTranslationJob(request)
	if err != nil {
		t.Fatalf("expected create request to succeed: %v", err)
	}
	assertCreateTranslationJobReadyResponse(t, response)
}

func TestTranslationJobSetupControllerCreateRejectsStaleValidationWithValidationStaleErrorKind(t *testing.T) {
	controller := NewTranslationJobSetupController(usecase.NewTranslationJobSetupContractStub())
	request := CreateTranslationJobRequestDTO{
		InputSourceID:        44,
		InputSource:          "translation_input",
		ValidationStatus:     usecase.TranslationJobSetupValidationStatusPass,
		ValidationPassSlices: []string{"input", "runtime"},
		Runtime: TranslationJobSetupRuntimeSelectionDTO{
			Provider:      "openai",
			Model:         "gpt-5.4-mini",
			ExecutionMode: "batch",
		},
		CredentialRef: "openai-primary",
	}
	setTranslationJobSetupFreshness(t, &request, time.Date(2026, 4, 27, 8, 30, 0, 0, time.UTC))

	response, err := controller.CreateTranslationJob(request)
	if err != nil {
		t.Fatalf("expected stale validation rejection without transport error, got %v", err)
	}
	if response.ErrorKind != "validation_stale" {
		t.Fatalf("expected validation_stale error kind, got %#v", response)
	}
	if response.ExecutionSummary != nil {
		t.Fatalf("expected nil execution summary on stale rejection, got %#v", response.ExecutionSummary)
	}
	if len(response.ValidationPassSlices) != 0 {
		t.Fatalf("expected no validation pass slices on stale rejection, got %#v", response.ValidationPassSlices)
	}
}

func TestTranslationJobSetupControllerSummaryMapsReadOnlyResponse(t *testing.T) {
	controller := NewTranslationJobSetupController(fakeTranslationJobSetupUsecase{
		getSummaryFunc: func(ctx context.Context, request usecase.GetTranslationJobSetupSummaryRequest) (usecase.TranslationJobSetupSummaryResult, error) {
			if ctx == nil {
				t.Fatal("expected request context")
			}
			if request.JobID != 91 {
				t.Fatalf("expected job id to be forwarded, got %#v", request)
			}
			return usecase.TranslationJobSetupSummaryResult{
				JobID:         91,
				JobState:      "ready",
				InputSource:   "translation_input",
				CanStartPhase: true,
				ExecutionSummary: usecase.TranslationJobExecutionSummary{
					Provider:      "openai",
					Model:         "gpt-5.4-mini",
					ExecutionMode: "batch",
				},
				ValidationPassSlices: []string{"input", "runtime"},
			}, nil
		},
	})

	response, err := controller.GetTranslationJobSetupSummary(GetTranslationJobSetupSummaryRequestDTO{JobID: 91})
	if err != nil {
		t.Fatalf("expected summary request to succeed: %v", err)
	}
	if response.JobID != 91 || response.JobState != "ready" {
		t.Fatalf("expected read-only summary response, got %#v", response)
	}
	if response.InputSource != "translation_input" {
		t.Fatalf("expected input source to stay read-only, got %#v", response)
	}
	if !response.CanStartPhase {
		t.Fatalf("expected summary response to expose phase startability, got %#v", response)
	}
	if response.ExecutionSummary.Provider != "openai" || response.ExecutionSummary.Model != "gpt-5.4-mini" || response.ExecutionSummary.ExecutionMode != "batch" {
		t.Fatalf("expected execution summary mapping, got %#v", response.ExecutionSummary)
	}
	if len(response.ValidationPassSlices) != 2 || response.ValidationPassSlices[0] != "input" || response.ValidationPassSlices[1] != "runtime" {
		t.Fatalf("expected validation pass slices, got %#v", response.ValidationPassSlices)
	}
}

func TestTranslationJobSetupControllerCreateMapsRejectedResponseWithoutTransportError(t *testing.T) {
	testCases := []struct {
		name              string
		usecaseErrorKind  usecase.TranslationJobSetupErrorKind
		expectedErrorKind usecase.TranslationJobSetupErrorKind
	}{
		{
			name:              "validation_failed alias は ready_required に正規化する",
			usecaseErrorKind:  usecase.TranslationJobSetupErrorKindValidationFailed,
			expectedErrorKind: usecase.TranslationJobSetupErrorKindReadyRequired,
		},
		{
			name:              "duplicate_input alias は duplicate_job_for_input に正規化する",
			usecaseErrorKind:  usecase.TranslationJobSetupErrorKindDuplicateInput,
			expectedErrorKind: usecase.TranslationJobSetupErrorKindDuplicateJobForInput,
		},
		{
			name:              "validation_stale は public error kind を維持する",
			usecaseErrorKind:  usecase.TranslationJobSetupErrorKindValidationStale,
			expectedErrorKind: usecase.TranslationJobSetupErrorKindValidationStale,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			controller := NewTranslationJobSetupController(fakeTranslationJobSetupUsecase{
				createJobFunc: func(_ context.Context, _ usecase.CreateTranslationJobRequest) (usecase.CreateTranslationJobResult, error) {
					return usecase.CreateTranslationJobResult{ErrorKind: testCase.usecaseErrorKind}, nil
				},
			})

			response, err := controller.CreateTranslationJob(CreateTranslationJobRequestDTO{})
			if err != nil {
				t.Fatalf("expected rejected response without transport error: %v", err)
			}
			assertCreateTranslationJobRejectedResponse(t, response, testCase.expectedErrorKind)
		})
	}
}

func TestTranslationJobSetupControllerCreateRejectsCreateWhenValidationDidNotPass(t *testing.T) {
	controller := NewTranslationJobSetupController(usecase.NewTranslationJobSetupContractStub())

	response, err := controller.CreateTranslationJob(CreateTranslationJobRequestDTO{
		InputSourceID:        44,
		InputSource:          "translation_input",
		ValidationStatus:     usecase.TranslationJobSetupValidationStatusFail,
		ValidationPassSlices: []string{"input"},
		Runtime: TranslationJobSetupRuntimeSelectionDTO{
			Provider:      "openai",
			Model:         "gpt-5.4-mini",
			ExecutionMode: "batch",
		},
		CredentialRef: "openai-primary",
	})
	if err != nil {
		t.Fatalf("expected rejected create response without transport error, got %v", err)
	}
	if response.ErrorKind != usecase.TranslationJobSetupErrorKindReadyRequired {
		t.Fatalf("expected ready_required error kind, got %#v", response)
	}
	if response.JobID != 0 || response.JobState != "" || response.InputSource != "" {
		t.Fatalf("expected rejected response to keep ready fields empty, got %#v", response)
	}
	if response.ExecutionSummary != nil {
		t.Fatalf("expected nil execution summary on rejection, got %#v", response.ExecutionSummary)
	}
	if len(response.ValidationPassSlices) != 0 {
		t.Fatalf("expected no validation pass slices on rejection, got %#v", response.ValidationPassSlices)
	}
}

func TestTranslationJobSetupControllerValidateWrapsUsecaseError(t *testing.T) {
	boom := errors.New("boom")
	controller := NewTranslationJobSetupController(fakeTranslationJobSetupUsecase{
		validateFunc: func(_ context.Context, _ usecase.ValidateTranslationJobSetupRequest) (usecase.TranslationJobSetupValidationResult, error) {
			return usecase.TranslationJobSetupValidationResult{}, boom
		},
	})

	_, err := controller.ValidateTranslationJobSetup(ValidateTranslationJobSetupRequestDTO{})
	if err == nil {
		t.Fatal("expected wrapped error, got nil")
	}
	if !errors.Is(err, boom) {
		t.Fatalf("expected wrapped error to contain original cause, got %v", err)
	}
}

func assertTranslationJobSetupOptionsContractResponse(t *testing.T, response TranslationJobSetupOptionsResponseDTO) {
	t.Helper()

	assertTranslationJobSetupInputCandidate(t, response.InputCandidates)
	assertTranslationJobSetupCredentialRefs(t, response.CredentialRefs)
	assertTranslationJobSetupExistingJob(t, response.ExistingJob)
	assertTranslationJobSetupRuntimeOptions(t, response.AIRuntimeOptions)
}

func assertTranslationJobSetupInputCandidate(t *testing.T, candidates []TranslationJobSetupInputCandidateDTO) {
	t.Helper()

	if len(candidates) != 1 {
		t.Fatalf("expected one input candidate, got %#v", candidates)
	}
	if candidates[0].ID != 44 || candidates[0].RecordCount != 120 {
		t.Fatalf("expected input metadata to be preserved, got %#v", candidates[0])
	}
	if candidates[0].SourceKind != "translation_input" || candidates[0].Label != "Dialogues import" {
		t.Fatalf("expected input source identity to be preserved, got %#v", candidates[0])
	}
}

func assertTranslationJobSetupCredentialRefs(t *testing.T, credentialRefs []TranslationJobSetupCredentialReferenceDTO) {
	t.Helper()

	if len(credentialRefs) != 2 {
		t.Fatalf("expected credential references, got %#v", credentialRefs)
	}
	assertTranslationJobSetupCredentialRef(t, credentialRefs[0], "openai-primary", true, false)
	assertTranslationJobSetupCredentialRef(t, credentialRefs[1], "gemini-missing", false, true)
}

func assertTranslationJobSetupCredentialRef(t *testing.T, credentialRef TranslationJobSetupCredentialReferenceDTO, wantRef string, wantConfigured bool, wantMissingSecret bool) {
	t.Helper()

	if credentialRef.SecretPlaintext != nil {
		t.Fatalf("expected secret plaintext to stay omitted, got %#v", credentialRef.SecretPlaintext)
	}
	if credentialRef.CredentialRef != wantRef || credentialRef.IsConfigured != wantConfigured || credentialRef.IsMissingSecret != wantMissingSecret {
		t.Fatalf("expected credential reference state ref=%q configured=%t missingSecret=%t, got %#v", wantRef, wantConfigured, wantMissingSecret, credentialRef)
	}
}

func assertTranslationJobSetupExistingJob(t *testing.T, existingJob *TranslationJobSetupExistingJobDTO) {
	t.Helper()

	if existingJob == nil || existingJob.Status != "ready" {
		t.Fatalf("expected existing job summary, got %#v", existingJob)
	}
	if existingJob.JobID != 88 || existingJob.InputSourceID != 44 || existingJob.InputSource != "translation_input" {
		t.Fatalf("expected existing job state to be preserved, got %#v", existingJob)
	}
}

func assertTranslationJobSetupRuntimeOptions(t *testing.T, runtimeOptions []TranslationJobSetupRuntimeOptionDTO) {
	t.Helper()

	if len(runtimeOptions) != 2 {
		t.Fatalf("expected runtime options mapping, got %#v", runtimeOptions)
	}
	assertTranslationJobSetupRuntimeOption(t, runtimeOptions[0], "openai", "gpt-5.4-mini", "batch")
	assertTranslationJobSetupRuntimeOption(t, runtimeOptions[1], "gemini", "gemini-2.5-pro", "sync")
}

func assertTranslationJobSetupRuntimeOption(t *testing.T, runtimeOption TranslationJobSetupRuntimeOptionDTO, wantProvider string, wantModel string, wantMode string) {
	t.Helper()

	if runtimeOption.Provider != wantProvider || runtimeOption.Model != wantModel || runtimeOption.Mode != wantMode {
		t.Fatalf("expected runtime option provider=%q model=%q mode=%q, got %#v", wantProvider, wantModel, wantMode, runtimeOption)
	}
}

func assertBlockingFailureValidationResponse(t *testing.T, response TranslationJobSetupValidationResponseDTO, expectedCategory string, expectedTargetSlices []string) {
	t.Helper()

	if response.Status != usecase.TranslationJobSetupValidationStatusFail {
		t.Fatalf("expected fail status, got %#v", response)
	}
	if response.BlockingFailureCategory == nil || *response.BlockingFailureCategory != expectedCategory {
		t.Fatalf("expected blocking failure category %q, got %#v", expectedCategory, response.BlockingFailureCategory)
	}
	if response.CanCreate {
		t.Fatalf("expected create to stay blocked, got %#v", response)
	}
	if !slices.Equal(response.TargetSlices, expectedTargetSlices) {
		t.Fatalf("expected target slices %#v, got %#v", expectedTargetSlices, response.TargetSlices)
	}
}

func assertCreateTranslationJobRequestReadyContract(t *testing.T, request usecase.CreateTranslationJobRequest, validatedAt time.Time) {
	t.Helper()

	if request.InputSourceID != 44 || request.InputSource != "translation_input" {
		t.Fatalf("expected input source to be forwarded, got %#v", request)
	}
	if request.ValidationStatus != "pass" {
		t.Fatalf("expected validation status to be forwarded, got %#v", request)
	}
	if observed := readTranslationJobSetupFreshness(t, request); !observed.Equal(validatedAt) {
		t.Fatalf("expected validated freshness %s to be forwarded, got %s", validatedAt.Format(time.RFC3339), observed.Format(time.RFC3339))
	}
	if !slices.Equal(request.ValidationPassSlices, []string{"input", "runtime"}) {
		t.Fatalf("expected validation pass slices to be forwarded, got %#v", request)
	}
	if request.Runtime.Provider != "openai" || request.Runtime.Model != "gpt-5.4-mini" || request.Runtime.ExecutionMode != "batch" {
		t.Fatalf("expected runtime selection to be forwarded, got %#v", request)
	}
	if request.CredentialRef != "openai-primary" {
		t.Fatalf("expected credential ref to be forwarded, got %#v", request)
	}
}

func assertCreateTranslationJobReadyResponse(t *testing.T, response CreateTranslationJobResponseDTO) {
	t.Helper()

	if response.JobID != 91 || response.JobState != "ready" {
		t.Fatalf("expected ready job response, got %#v", response)
	}
	if response.InputSource != "translation_input" {
		t.Fatalf("expected input source summary, got %#v", response)
	}
	if response.ErrorKind != "" {
		t.Fatalf("expected ready response to keep error kind empty, got %#v", response)
	}
	if response.ExecutionSummary == nil {
		t.Fatal("expected execution summary on ready response")
	}
	if response.ExecutionSummary.Model != "gpt-5.4-mini" {
		t.Fatalf("expected execution summary mapping, got %#v", response.ExecutionSummary)
	}
	if len(response.ValidationPassSlices) != 2 || response.ValidationPassSlices[1] != "runtime" {
		t.Fatalf("expected validation pass slices, got %#v", response.ValidationPassSlices)
	}
}

func assertCreateTranslationJobRejectedResponse(t *testing.T, response CreateTranslationJobResponseDTO, expectedErrorKind usecase.TranslationJobSetupErrorKind) {
	t.Helper()

	if response.ErrorKind != expectedErrorKind {
		t.Fatalf("expected error kind %q, got %#v", expectedErrorKind, response)
	}
	if response.JobID != 0 || response.JobState != "" || response.InputSource != "" {
		t.Fatalf("expected rejected response to keep ready fields empty, got %#v", response)
	}
	if response.ExecutionSummary != nil {
		t.Fatalf("expected nil execution summary on rejection, got %#v", response.ExecutionSummary)
	}
	if len(response.ValidationPassSlices) != 0 {
		t.Fatalf("expected no validation pass slices on rejection, got %#v", response.ValidationPassSlices)
	}
}

func setTranslationJobSetupFreshness(t *testing.T, target any, validatedAt time.Time) {
	t.Helper()

	value := reflect.ValueOf(target)
	if value.Kind() != reflect.Pointer || value.Elem().Kind() != reflect.Struct {
		t.Fatalf("expected pointer to struct target, got %T", target)
	}

	field := value.Elem().FieldByName("ValidatedAt")
	if !field.IsValid() {
		t.Fatalf("expected %T to expose freshness field ValidatedAt", target)
	}
	if !field.CanSet() {
		t.Fatalf("expected freshness field ValidatedAt on %T to be settable", target)
	}

	switch {
	case field.Kind() == reflect.String:
		field.SetString(validatedAt.UTC().Format(time.RFC3339))
	case field.Type() == reflect.TypeOf(time.Time{}):
		field.Set(reflect.ValueOf(validatedAt.UTC()))
	default:
		t.Fatalf("expected freshness field ValidatedAt on %T to be string or time.Time, got %s", target, field.Type())
	}
}

func readTranslationJobSetupFreshness(t *testing.T, source any) time.Time {
	t.Helper()

	value := reflect.ValueOf(source)
	if value.Kind() == reflect.Pointer {
		value = value.Elem()
	}
	if value.Kind() != reflect.Struct {
		t.Fatalf("expected struct source, got %T", source)
	}

	field := value.FieldByName("ValidatedAt")
	if !field.IsValid() {
		t.Fatalf("expected %T to expose freshness field ValidatedAt", source)
	}

	switch {
	case field.Kind() == reflect.String:
		validatedAt, err := time.Parse(time.RFC3339, field.String())
		if err != nil {
			t.Fatalf("expected RFC3339 freshness string, got %q: %v", field.String(), err)
		}
		return validatedAt.UTC()
	case field.Type() == reflect.TypeOf(time.Time{}):
		return field.Interface().(time.Time).UTC()
	default:
		t.Fatalf("expected freshness field ValidatedAt on %T to be string or time.Time, got %s", source, field.Type())
		return time.Time{}
	}
}
