package service

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"path/filepath"
	"reflect"
	"slices"
	"strings"
	"testing"
	"time"

	"aitranslationenginejp/internal/repository"
)

type translationJobSetupServiceCacheCleaner interface {
	DeleteTranslationCacheByXEditID(context.Context, int64) error
}

type fakePersistentTranslationJobSetupSourceRepository struct {
	getByIDFunc            func(context.Context, int64) (repository.XEditExtractedData, error)
	listXEditExtractedData func(context.Context) ([]repository.XEditExtractedData, error)
	getExistingJobFunc     func(context.Context, int64) (repository.TranslationJob, error)
}

type fakeTranslationJobSetupSecretStore struct {
	loadFunc func(context.Context, string) (string, error)
}

type loggedTranslationJobSetupProviderReachabilityRequest struct {
	method        string
	url           string
	authorization string
	contentType   string
	body          string
}

type fakeTranslationJobSetupProviderReachabilityTransport struct {
	requests []loggedTranslationJobSetupProviderReachabilityRequest
	doFunc   func(*http.Request) (*http.Response, error)
}

type translationJobSetupProviderReachabilityTestCase struct {
	name                string
	provider            string
	model               string
	executionMode       string
	credentialRef       string
	secretValue         string
	expectedURL         string
	expectedAuth        string
	expectedBodySnippet string
	assertSecretFreeLog bool
}

func (transport *fakeTranslationJobSetupProviderReachabilityTransport) Do(req *http.Request) (*http.Response, error) {
	body, err := io.ReadAll(req.Body)
	if err != nil {
		return nil, fmt.Errorf("read provider reachability request body: %w", err)
	}
	if closeErr := req.Body.Close(); closeErr != nil {
		return nil, fmt.Errorf("close provider reachability request body: %w", closeErr)
	}

	transport.requests = append(transport.requests, loggedTranslationJobSetupProviderReachabilityRequest{
		method:        req.Method,
		url:           req.URL.String(),
		authorization: req.Header.Get("Authorization"),
		contentType:   req.Header.Get("Content-Type"),
		body:          string(body),
	})
	req.Body = io.NopCloser(bytes.NewReader(body))

	if transport.doFunc != nil {
		return transport.doFunc(req)
	}
	return &http.Response{
		StatusCode: http.StatusOK,
		Body:       io.NopCloser(bytes.NewReader([]byte(`{"ok":true}`))),
	}, nil
}

func (store fakeTranslationJobSetupSecretStore) Load(ctx context.Context, key string) (string, error) {
	if store.loadFunc != nil {
		return store.loadFunc(ctx, key)
	}
	return "", nil
}

func (sourceRepository fakePersistentTranslationJobSetupSourceRepository) GetXEditExtractedDataByID(
	ctx context.Context,
	id int64,
) (repository.XEditExtractedData, error) {
	if sourceRepository.getByIDFunc != nil {
		return sourceRepository.getByIDFunc(ctx, id)
	}
	return repository.XEditExtractedData{}, nil
}

func (sourceRepository fakePersistentTranslationJobSetupSourceRepository) ListXEditExtractedData(
	ctx context.Context,
) ([]repository.XEditExtractedData, error) {
	if sourceRepository.listXEditExtractedData != nil {
		return sourceRepository.listXEditExtractedData(ctx)
	}
	return nil, nil
}

func (sourceRepository fakePersistentTranslationJobSetupSourceRepository) GetExistingTranslationJob(
	ctx context.Context,
	xEditID int64,
) (repository.TranslationJob, error) {
	if sourceRepository.getExistingJobFunc != nil {
		return sourceRepository.getExistingJobFunc(ctx, xEditID)
	}
	return repository.TranslationJob{}, repository.ErrNotFound
}

func newServerOwnedTranslationJobSetupOptionsReadModelWithConfiguredGeminiCredential() TranslationJobSetupOptionsReadModel {
	readModel := newServerOwnedTranslationJobSetupOptionsReadModel()
	for index := range readModel.CredentialRefs {
		if readModel.CredentialRefs[index].Provider != MasterPersonaProviderGemini {
			continue
		}
		readModel.CredentialRefs[index].CredentialRef = "gemini-primary"
		readModel.CredentialRefs[index].IsConfigured = true
		readModel.CredentialRefs[index].IsMissingSecret = false
	}
	return readModel
}

func translationJobSetupProviderReachabilityCases() []translationJobSetupProviderReachabilityTestCase {
	return []translationJobSetupProviderReachabilityTestCase{
		{
			name:                "openai",
			provider:            translationJobSetupRealProviderOpenAI,
			model:               translationJobSetupModelGPT54Mini,
			executionMode:       "batch",
			credentialRef:       "openai-primary",
			secretValue:         "openai-test-key",
			expectedURL:         translationJobSetupOpenAIBaseURL + "/chat/completions",
			expectedAuth:        "Bearer openai-test-key",
			expectedBodySnippet: `"messages"`,
		},
		{
			name:                "xai",
			provider:            MasterPersonaProviderXAI,
			model:               "grok-4",
			executionMode:       "sync",
			credentialRef:       "xai-primary",
			secretValue:         "xai-test-key",
			expectedURL:         "https://xai.example/v1/chat/completions",
			expectedAuth:        "Bearer xai-test-key",
			expectedBodySnippet: `"messages"`,
		},
		{
			name:                "gemini",
			provider:            MasterPersonaProviderGemini,
			model:               "gemini-2.5-pro",
			executionMode:       "sync",
			credentialRef:       "gemini-primary",
			secretValue:         "gemini-test-key",
			expectedURL:         "https://generativelanguage.googleapis.com/v1beta/models/gemini-2.5-pro:generateContent",
			expectedAuth:        "",
			expectedBodySnippet: `"contents"`,
			assertSecretFreeLog: true,
		},
	}
}

func newTranslationJobSetupServiceForProviderReachabilityTest(
	t *testing.T,
	validatedAt time.Time,
	testCase translationJobSetupProviderReachabilityTestCase,
) (*TranslationJobSetupService, *fakeTranslationJobSetupProviderReachabilityTransport) {
	t.Helper()

	transport := &fakeTranslationJobSetupProviderReachabilityTransport{}
	optionsReadModel := newServerOwnedTranslationJobSetupOptionsReadModel()
	if testCase.provider == MasterPersonaProviderGemini {
		optionsReadModel = newServerOwnedTranslationJobSetupOptionsReadModelWithConfiguredGeminiCredential()
	}

	service := &TranslationJobSetupService{
		now:                           func() time.Time { return validatedAt },
		optionsReadModel:              optionsReadModel,
		providerReachabilityTransport: transport,
		secretStore: fakeTranslationJobSetupSecretStore{
			loadFunc: func(_ context.Context, key string) (string, error) {
				if key != translationJobSetupMasterPersonaSecretKey(testCase.provider) {
					t.Fatalf("expected %s secret key lookup, got %q", testCase.provider, key)
				}
				return testCase.secretValue, nil
			},
		},
	}

	return service, transport
}

func assertTranslationJobSetupProviderReachabilityRequest(
	t *testing.T,
	requests []loggedTranslationJobSetupProviderReachabilityRequest,
	testCase translationJobSetupProviderReachabilityTestCase,
) {
	t.Helper()

	if len(requests) != 1 {
		t.Fatalf("expected one logged reachability request, got %#v", requests)
	}

	logged := requests[0]
	if logged.method != http.MethodPost {
		t.Fatalf("expected POST reachability request, got %#v", logged)
	}
	if logged.url != testCase.expectedURL {
		t.Fatalf("expected reachability url %q, got %#v", testCase.expectedURL, logged)
	}
	if logged.authorization != testCase.expectedAuth {
		t.Fatalf("expected authorization %q, got %#v", testCase.expectedAuth, logged)
	}
	if logged.contentType != "application/json" {
		t.Fatalf("expected json reachability request, got %#v", logged)
	}
	if !strings.Contains(logged.body, testCase.expectedBodySnippet) {
		t.Fatalf("expected request body to contain %q, got %#v", testCase.expectedBodySnippet, logged)
	}
	if testCase.assertSecretFreeLog {
		if strings.Contains(logged.url, testCase.secretValue) {
			t.Fatalf("expected reachability url to exclude secret %q, got %#v", testCase.secretValue, logged)
		}
		if strings.Contains(fmt.Sprintf("%#v", logged), testCase.secretValue) {
			t.Fatalf("expected logged reachability request to exclude secret %q, got %#v", testCase.secretValue, logged)
		}
	}
}

func TestTranslationJobSetupServiceReadOptionsReturnsConfiguredLiveStateInsteadOfFallbackReadModel(t *testing.T) {
	service := &TranslationJobSetupService{
		optionsReadModel: TranslationJobSetupOptionsReadModel{
			InputCandidates: []TranslationJobSetupInputCandidateReadModel{
				{
					ID:           77,
					Label:        "Quest fragments import",
					SourceKind:   "translation_input",
					RecordCount:  32,
					RegisteredAt: time.Date(2026, 4, 20, 9, 0, 0, 0, time.UTC),
				},
			},
			ExistingJob: nil,
			SharedDictionaries: []TranslationJobSetupDictionaryOptionReadModel{
				{ID: "dict-quest", Label: "Quest Dictionary"},
			},
			SharedPersonas: []TranslationJobSetupPersonaOptionReadModel{
				{ID: "persona-mage", Label: "Mage Persona"},
			},
			AIRuntimeOptions: []TranslationJobSetupRuntimeOptionReadModel{
				{Provider: "xai", Model: "grok-4", Mode: "sync"},
			},
			CredentialRefs: []TranslationJobSetupCredentialReferenceReadModel{
				{Provider: "xai", CredentialRef: "xai-primary", IsConfigured: true, IsMissingSecret: false},
			},
		},
	}

	got, err := service.ReadOptions(context.Background())
	if err != nil {
		t.Fatalf("expected configured live options to be returned: %v", err)
	}

	want := service.optionsReadModel
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("expected live options %#v, got %#v", want, got)
	}
	if reflect.DeepEqual(got, TranslationJobSetupReadOptions()) {
		t.Fatalf("expected configured live options to differ from fallback read model, got %#v", got)
	}
}

func TestPersistentTranslationJobSetupServiceReadOptionsLoadsImportedInputsFromRepositoryLister(t *testing.T) {
	importedAt := time.Date(2026, 4, 27, 16, 0, 0, 0, time.UTC)
	service := NewPersistentTranslationJobSetupService(
		nil,
		fakePersistentTranslationJobSetupSourceRepository{
			listXEditExtractedData: func(context.Context) ([]repository.XEditExtractedData, error) {
				return []repository.XEditExtractedData{{
					ID:               44,
					TargetPluginName: "QuestFragments.esp",
					SourceFilePath:   "/imports/quest-fragments.json",
					RecordCount:      32,
					ImportedAt:       importedAt,
				}}, nil
			},
		},
		nil,
		nil,
		nil,
		nil,
		nil,
	)

	got, err := service.ReadOptions(context.Background())
	if err != nil {
		t.Fatalf("expected persistent read options to load imported inputs: %v", err)
	}

	want := []TranslationJobSetupInputCandidateReadModel{{
		ID:           44,
		Label:        "QuestFragments.esp",
		SourceKind:   "translation_input",
		RecordCount:  32,
		RegisteredAt: importedAt,
	}}
	if !reflect.DeepEqual(got.InputCandidates, want) {
		t.Fatalf("expected imported input candidates %#v, got %#v", want, got.InputCandidates)
	}
}

func TestPersistentTranslationJobSetupServiceReadOptionsReturnsExistingReadyJobFromRepositories(t *testing.T) {
	ctx := context.Background()
	service, sourceRepository, jobRepository, closeRepositories := newSQLiteBackedTranslationJobSetupServiceForTest(t)
	defer closeRepositories()

	input := createSQLiteTranslationJobSetupInputFixture(t, sourceRepository, repository.XEditExtractedDataDraft{
		SourceFilePath:    "/imports/job-setup.json",
		SourceContentHash: "job-setup-hash",
		SourceTool:        "xedit",
		TargetPluginName:  "JobSetup.esp",
		TargetPluginType:  "esp",
		RecordCount:       12,
		ImportedAt:        time.Date(2026, 4, 27, 16, 30, 0, 0, time.UTC),
	})
	job, err := jobRepository.CreateTranslationJob(ctx, repository.TranslationJobDraft{
		XEditExtractedDataID: input.ID,
		JobName:              "translation-job-setup-existing",
		State:                translationJobSetupJobStateReady,
		ProgressPercent:      0,
	})
	if err != nil {
		t.Fatalf("expected ready translation job fixture to be created: %v", err)
	}

	got, err := service.ReadOptions(ctx)
	if err != nil {
		t.Fatalf("expected repository-backed read options to succeed: %v", err)
	}

	if got.ExistingJob == nil {
		t.Fatalf("expected repository-backed read options to return existing job state, got %#v", got)
	}
	if got.ExistingJob.JobID != job.ID {
		t.Fatalf("expected existing job id %d, got %#v", job.ID, got.ExistingJob)
	}
	if got.ExistingJob.Status != translationJobSetupJobStateReady {
		t.Fatalf("expected existing job status %q, got %#v", translationJobSetupJobStateReady, got.ExistingJob)
	}
	if got.ExistingJob.InputSource != translationJobSetupInputSourceTranslation {
		t.Fatalf("expected existing job input source %q, got %#v", translationJobSetupInputSourceTranslation, got.ExistingJob)
	}
}

func TestTranslationJobSetupServiceReadOptionsDoesNotSurfaceExistingJobOutsideCurrentCandidates(t *testing.T) {
	service := NewPersistentTranslationJobSetupService(
		nil,
		fakePersistentTranslationJobSetupSourceRepository{
			listXEditExtractedData: func(context.Context) ([]repository.XEditExtractedData, error) {
				return []repository.XEditExtractedData{{
					ID:               44,
					TargetPluginName: "CurrentCandidate.esp",
					SourceFilePath:   "/imports/current-candidate.json",
					RecordCount:      12,
					ImportedAt:       time.Date(2026, 4, 27, 16, 45, 0, 0, time.UTC),
				}}, nil
			},
			getExistingJobFunc: func(_ context.Context, xEditID int64) (repository.TranslationJob, error) {
				if xEditID != 44 {
					return repository.TranslationJob{ID: 91, XEditExtractedDataID: 99, State: translationJobSetupJobStateReady}, nil
				}
				return repository.TranslationJob{}, repository.ErrNotFound
			},
		},
		nil,
		nil,
		nil,
		nil,
		nil,
	)

	got, err := service.ReadOptions(context.Background())
	if err != nil {
		t.Fatalf("expected read options to ignore existing jobs outside current candidates: %v", err)
	}

	if got.ExistingJob != nil {
		t.Fatalf("expected no existing job for current candidates, got %#v", got.ExistingJob)
	}
	if len(got.InputCandidates) != 1 || got.InputCandidates[0].ID != 44 {
		t.Fatalf("expected current candidates to stay intact, got %#v", got.InputCandidates)
	}
}

func TestTranslationJobSetupServiceValidateRequestUsesSameLiveOptionsAsReadOptions(t *testing.T) {
	service := &TranslationJobSetupService{
		now: func() time.Time { return time.Date(2026, 4, 27, 15, 0, 0, 0, time.UTC) },
		optionsReadModel: TranslationJobSetupOptionsReadModel{
			AIRuntimeOptions: []TranslationJobSetupRuntimeOptionReadModel{
				{Provider: "xai", Model: "grok-4", Mode: "sync"},
			},
			CredentialRefs: []TranslationJobSetupCredentialReferenceReadModel{
				{Provider: "xai", CredentialRef: "xai-primary", IsConfigured: true, IsMissingSecret: false},
			},
		},
	}

	options, err := service.ReadOptions(context.Background())
	if err != nil {
		t.Fatalf("expected live options read to succeed: %v", err)
	}
	if !reflect.DeepEqual(options.AIRuntimeOptions, service.optionsReadModel.AIRuntimeOptions) {
		t.Fatalf("expected runtime options %#v, got %#v", service.optionsReadModel.AIRuntimeOptions, options.AIRuntimeOptions)
	}
	if !reflect.DeepEqual(options.CredentialRefs, service.optionsReadModel.CredentialRefs) {
		t.Fatalf("expected credential refs %#v, got %#v", service.optionsReadModel.CredentialRefs, options.CredentialRefs)
	}

	decision, err := service.ValidateRequest(context.Background(), TranslationJobSetupValidationRequest{
		InputSourceID: 44,
		Provider:      options.AIRuntimeOptions[0].Provider,
		Model:         options.AIRuntimeOptions[0].Model,
		ExecutionMode: options.AIRuntimeOptions[0].Mode,
		CredentialRef: options.CredentialRefs[0].CredentialRef,
	})
	if err != nil {
		t.Fatalf("expected validation to succeed for option advertised by live state: %v", err)
	}
	if decision.Status != translationJobSetupValidationStatusPass || !decision.CanCreate {
		t.Fatalf("expected live option validation to pass, got %#v", decision)
	}
}

func TestTranslationJobSetupServiceValidateRequestRejectsCredentialNotPresentInLiveOptions(t *testing.T) {
	service := &TranslationJobSetupService{
		now: func() time.Time { return time.Date(2026, 4, 27, 15, 15, 0, 0, time.UTC) },
		optionsReadModel: TranslationJobSetupOptionsReadModel{
			AIRuntimeOptions: []TranslationJobSetupRuntimeOptionReadModel{
				{Provider: "xai", Model: "grok-4", Mode: "sync"},
			},
			CredentialRefs: []TranslationJobSetupCredentialReferenceReadModel{
				{Provider: "xai", CredentialRef: "xai-primary", IsConfigured: true, IsMissingSecret: false},
			},
		},
	}

	decision, err := service.ValidateRequest(context.Background(), TranslationJobSetupValidationRequest{
		InputSourceID: 44,
		Provider:      "xai",
		Model:         "grok-4",
		ExecutionMode: "sync",
		CredentialRef: "xai-secondary",
	})
	if err != nil {
		t.Fatalf("expected missing live credential to be classified without transport error: %v", err)
	}
	if decision.Status != translationJobSetupValidationStatusFail {
		t.Fatalf("expected missing live credential to fail validation, got %#v", decision)
	}
	if decision.BlockingFailureCategory == nil || *decision.BlockingFailureCategory != translationJobSetupBlockingFailureMissingSecretRef() {
		t.Fatalf("expected credential_missing category, got %#v", decision.BlockingFailureCategory)
	}
	if !slices.Equal(decision.TargetSlices, []string{"credentials"}) {
		t.Fatalf("expected credentials slice failure, got %#v", decision.TargetSlices)
	}
	if decision.CanCreate {
		t.Fatalf("expected missing live credential to block create, got %#v", decision)
	}
}

func TestPersistentTranslationJobSetupServiceValidateRequestReturnsCacheMissingFromActualCacheState(t *testing.T) {
	ctx := context.Background()
	service, sourceRepository, _, closeRepositories := newSQLiteBackedTranslationJobSetupServiceForTest(t)
	defer closeRepositories()

	input := createSQLiteTranslationJobSetupInputFixture(t, sourceRepository, repository.XEditExtractedDataDraft{
		SourceFilePath:    "/imports/cache-state.json",
		SourceContentHash: "cache-state-hash",
		SourceTool:        "xedit",
		TargetPluginName:  "CacheState.esp",
		TargetPluginType:  "esp",
		RecordCount:       1,
		ImportedAt:        time.Date(2026, 4, 27, 17, 0, 0, 0, time.UTC),
	})
	_, err := sourceRepository.CreateTranslationRecord(ctx, repository.TranslationRecordDraft{
		XEditExtractedDataID: input.ID,
		FormID:               "0000ABCD",
		EditorID:             "NPC_CacheState",
		RecordType:           "NPC_",
	})
	if err != nil {
		t.Fatalf("expected translation cache fixture to be created: %v", err)
	}

	cacheCleaner, ok := sourceRepository.(translationJobSetupServiceCacheCleaner)
	if !ok {
		t.Fatal("expected sqlite translation source repository to support cache deletion")
	}
	if deleteErr := cacheCleaner.DeleteTranslationCacheByXEditID(ctx, input.ID); deleteErr != nil {
		t.Fatalf("expected translation cache fixture to be deleted: %v", deleteErr)
	}

	decision, err := service.ValidateRequest(ctx, TranslationJobSetupValidationRequest{
		InputSourceID: input.ID,
		Provider:      "xai",
		Model:         "grok-4",
		ExecutionMode: "sync",
		CredentialRef: "configured-ref",
	})
	if err != nil {
		t.Fatalf("expected cache-missing validation to classify without transport error: %v", err)
	}
	if decision.Status != translationJobSetupValidationStatusFail {
		t.Fatalf("expected deleted cache to fail validation, got %#v", decision)
	}
	if decision.BlockingFailureCategory == nil || *decision.BlockingFailureCategory != translationJobSetupBlockingFailureCacheMissing {
		t.Fatalf("expected cache_missing category, got %#v", decision.BlockingFailureCategory)
	}
	if !slices.Equal(decision.TargetSlices, []string{"input"}) {
		t.Fatalf("expected input slice failure for missing cache, got %#v", decision.TargetSlices)
	}
	if decision.CanCreate {
		t.Fatalf("expected missing cache to block create, got %#v", decision)
	}
}

func TestTranslationJobSetupServiceValidateRequestAllowsRealProviders(t *testing.T) {
	validatedAt := time.Date(2026, 4, 27, 12, 0, 0, 0, time.UTC)
	service := &TranslationJobSetupService{now: func() time.Time { return validatedAt }}

	testCases := []struct {
		name     string
		provider string
	}{
		{name: "gemini", provider: "gemini"},
		{name: "xai", provider: "xai"},
		{name: "lm_studio", provider: "lm_studio"},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			decision, err := service.ValidateRequest(context.Background(), TranslationJobSetupValidationRequest{
				InputSourceID: 44,
				Provider:      testCase.provider,
				Model:         "test-model",
				ExecutionMode: "sync",
				CredentialRef: "configured-ref",
			})
			if err != nil {
				t.Fatalf("expected validation to succeed for real provider %q: %v", testCase.provider, err)
			}
			if decision.Status != translationJobSetupValidationStatusPass || !decision.CanCreate {
				t.Fatalf("expected passing validation for real provider %q, got %#v", testCase.provider, decision)
			}
			if !decision.ValidatedAt.Equal(validatedAt) {
				t.Fatalf("expected validatedAt %s, got %s", validatedAt, decision.ValidatedAt)
			}
			if !slices.Equal(decision.PassSlices, translationJobSetupAllSlices) {
				t.Fatalf("expected pass slices %#v, got %#v", translationJobSetupAllSlices, decision.PassSlices)
			}
		})
	}
}

func TestTranslationJobSetupServiceValidateRequestRejectsUnsupportedProvider(t *testing.T) {
	validatedAt := time.Date(2026, 4, 27, 12, 30, 0, 0, time.UTC)
	service := &TranslationJobSetupService{now: func() time.Time { return validatedAt }}

	decision, err := service.ValidateRequest(context.Background(), TranslationJobSetupValidationRequest{
		InputSourceID: 44,
		Provider:      "fake-provider",
		Model:         "test-model",
		ExecutionMode: "sync",
		CredentialRef: "configured-ref",
	})
	if err != nil {
		t.Fatalf("expected validation to classify unsupported provider without transport error: %v", err)
	}
	if decision.Status != translationJobSetupValidationStatusFail {
		t.Fatalf("expected unsupported provider to fail validation, got %#v", decision)
	}
	if decision.CanCreate {
		t.Fatalf("expected unsupported provider to block create, got %#v", decision)
	}
	if decision.BlockingFailureCategory == nil || *decision.BlockingFailureCategory == "" {
		t.Fatalf("expected unsupported provider failure category, got %#v", decision)
	}
	if !slices.Equal(decision.TargetSlices, []string{"runtime"}) {
		t.Fatalf("expected runtime slice failure for unsupported provider, got %#v", decision.TargetSlices)
	}
	if !decision.ValidatedAt.Equal(validatedAt) {
		t.Fatalf("expected validatedAt %s, got %s", validatedAt, decision.ValidatedAt)
	}
}

func TestTranslationJobSetupServiceValidateRequestRejectsRuntimeMismatchAgainstServerOwnedCatalog(t *testing.T) {
	validatedAt := time.Date(2026, 4, 27, 12, 45, 0, 0, time.UTC)
	service := &TranslationJobSetupService{
		now:              func() time.Time { return validatedAt },
		optionsReadModel: newServerOwnedTranslationJobSetupOptionsReadModel(),
	}

	decision, err := service.ValidateRequest(context.Background(), TranslationJobSetupValidationRequest{
		InputSourceID: 44,
		Provider:      "xai",
		Model:         "gpt-5.4-mini",
		ExecutionMode: "sync",
		CredentialRef: "xai-primary",
	})
	if err != nil {
		t.Fatalf("expected catalog mismatch to be classified without transport error: %v", err)
	}
	if decision.Status != translationJobSetupValidationStatusFail {
		t.Fatalf("expected runtime mismatch to fail validation, got %#v", decision)
	}
	if decision.BlockingFailureCategory == nil || *decision.BlockingFailureCategory != translationJobSetupBlockingFailureProviderModeUnsupported {
		t.Fatalf("expected provider_mode_unsupported category, got %#v", decision.BlockingFailureCategory)
	}
	if !slices.Equal(decision.TargetSlices, []string{"runtime"}) {
		t.Fatalf("expected runtime slice failure, got %#v", decision.TargetSlices)
	}
	if decision.CanCreate {
		t.Fatalf("expected catalog mismatch to block create, got %#v", decision)
	}
	if !decision.ValidatedAt.Equal(validatedAt) {
		t.Fatalf("expected validatedAt %s, got %s", validatedAt, decision.ValidatedAt)
	}
}

func TestTranslationJobSetupServiceValidateRequestRejectsCredentialReferenceFromDifferentProvider(t *testing.T) {
	service := &TranslationJobSetupService{
		now:              func() time.Time { return time.Date(2026, 4, 27, 13, 0, 0, 0, time.UTC) },
		optionsReadModel: newServerOwnedTranslationJobSetupOptionsReadModel(),
	}

	decision, err := service.ValidateRequest(context.Background(), TranslationJobSetupValidationRequest{
		InputSourceID: 44,
		Provider:      "openai",
		Model:         "gpt-5.4-mini",
		ExecutionMode: "batch",
		CredentialRef: "xai-primary",
	})
	if err != nil {
		t.Fatalf("expected provider-bound credential mismatch to be classified without transport error: %v", err)
	}
	if decision.Status != translationJobSetupValidationStatusFail {
		t.Fatalf("expected mismatched credential ref to fail validation, got %#v", decision)
	}
	if decision.BlockingFailureCategory == nil || *decision.BlockingFailureCategory != translationJobSetupBlockingFailureMissingSecretRef() {
		t.Fatalf("expected credential_missing category, got %#v", decision.BlockingFailureCategory)
	}
	if !slices.Equal(decision.TargetSlices, []string{"credentials"}) {
		t.Fatalf("expected credentials slice failure, got %#v", decision.TargetSlices)
	}
	if decision.CanCreate {
		t.Fatalf("expected mismatched credential ref to block create, got %#v", decision)
	}
}

func TestTranslationJobSetupServiceValidateRequestRejectsProviderUnreachableFromServerOwnedReachability(t *testing.T) {
	validatedAt := time.Date(2026, 4, 27, 13, 15, 0, 0, time.UTC)
	t.Setenv(translationJobSetupXAIBaseURLEnv, "http://127.0.0.1:1")
	service := &TranslationJobSetupService{
		now:              func() time.Time { return validatedAt },
		optionsReadModel: newServerOwnedTranslationJobSetupOptionsReadModel(),
		secretStore: fakeTranslationJobSetupSecretStore{
			loadFunc: func(_ context.Context, key string) (string, error) {
				if key != translationJobSetupMasterPersonaSecretKey(MasterPersonaProviderXAI) {
					t.Fatalf("expected xai secret key lookup, got %q", key)
				}
				return "xai-test-key", nil
			},
		},
	}

	decision, err := service.ValidateRequest(context.Background(), TranslationJobSetupValidationRequest{
		InputSourceID: 44,
		Provider:      MasterPersonaProviderXAI,
		Model:         "grok-4",
		ExecutionMode: "sync",
		CredentialRef: "xai-primary",
	})
	if err != nil {
		t.Fatalf("expected provider reachability failure to be classified without transport error: %v", err)
	}
	if decision.Status != translationJobSetupValidationStatusFail {
		t.Fatalf("expected provider_unreachable to fail validation, got %#v", decision)
	}
	if decision.BlockingFailureCategory == nil || *decision.BlockingFailureCategory != translationJobSetupBlockingFailureProviderUnreachable {
		t.Fatalf("expected provider_unreachable category, got %#v", decision.BlockingFailureCategory)
	}
	if !slices.Equal(decision.TargetSlices, []string{"runtime"}) {
		t.Fatalf("expected runtime slice failure, got %#v", decision.TargetSlices)
	}
	if decision.CanCreate {
		t.Fatalf("expected provider_unreachable to block create, got %#v", decision)
	}
	if !decision.ValidatedAt.Equal(validatedAt) {
		t.Fatalf("expected validatedAt %s, got %s", validatedAt, decision.ValidatedAt)
	}
}

func TestTranslationJobSetupServiceValidateRequestUsesInjectedTransportForProviderReachability(t *testing.T) {
	t.Setenv(translationJobSetupXAIBaseURLEnv, "https://xai.example/v1")

	for _, testCase := range translationJobSetupProviderReachabilityCases() {
		t.Run(testCase.name, func(t *testing.T) {
			validatedAt := time.Date(2026, 4, 27, 13, 20, 0, 0, time.UTC)
			service, transport := newTranslationJobSetupServiceForProviderReachabilityTest(t, validatedAt, testCase)

			decision, err := service.ValidateRequest(context.Background(), TranslationJobSetupValidationRequest{
				InputSourceID: 44,
				Provider:      testCase.provider,
				Model:         testCase.model,
				ExecutionMode: testCase.executionMode,
				CredentialRef: testCase.credentialRef,
			})
			if err != nil {
				t.Fatalf("expected injected transport validation to succeed without real network: %v", err)
			}
			if decision.Status != translationJobSetupValidationStatusPass || !decision.CanCreate {
				t.Fatalf("expected injected transport validation to pass, got %#v", decision)
			}
			if !decision.ValidatedAt.Equal(validatedAt) {
				t.Fatalf("expected validatedAt %s, got %s", validatedAt, decision.ValidatedAt)
			}
			if !slices.Equal(decision.PassSlices, translationJobSetupAllSlices) {
				t.Fatalf("expected pass slices %#v, got %#v", translationJobSetupAllSlices, decision.PassSlices)
			}
			assertTranslationJobSetupProviderReachabilityRequest(t, transport.requests, testCase)
		})
	}
}

func TestTranslationJobSetupServiceEvaluateCreateRequestRejectsBlockingValidationFailure(t *testing.T) {
	validatedAt := time.Date(2026, 4, 27, 14, 0, 0, 0, time.UTC)
	service := &TranslationJobSetupService{now: func() time.Time { return validatedAt }}

	decision, err := service.EvaluateCreateRequest(context.Background(), TranslationJobSetupCreateRequest{
		InputSourceID:    44,
		ValidationStatus: translationJobSetupValidationStatusPass,
		ValidatedAt:      validatedAt,
		Provider:         "fake-provider",
		Model:            "test-model",
		ExecutionMode:    "sync",
		CredentialRef:    "configured-ref",
	})
	if err != nil {
		t.Fatalf("expected create evaluation to classify blocking failure without transport error: %v", err)
	}
	if decision.CanCreate {
		t.Fatalf("expected blocking validation failure to reject create, got %#v", decision)
	}
	if decision.ErrorKind != translationJobSetupErrorKindValidationFailed {
		t.Fatalf("expected validation_failed error kind, got %#v", decision)
	}
}

func TestTranslationJobSetupServiceEvaluateCreateRequestKeepsValidationFailedDistinctFromStale(t *testing.T) {
	service := &TranslationJobSetupService{now: func() time.Time { return time.Date(2026, 4, 27, 14, 30, 0, 0, time.UTC) }}

	decision, err := service.EvaluateCreateRequest(context.Background(), TranslationJobSetupCreateRequest{
		ValidationStatus: translationJobSetupValidationStatusFail,
		ValidatedAt:      time.Date(2026, 4, 27, 8, 30, 0, 0, time.UTC),
	})
	if err != nil {
		t.Fatalf("expected failed validation to reject create without transport error: %v", err)
	}
	if decision.CanCreate {
		t.Fatalf("expected failed validation to reject create, got %#v", decision)
	}
	if decision.ErrorKind != translationJobSetupErrorKindValidationFailed {
		t.Fatalf("expected validation_failed to win over stale rejection, got %#v", decision)
	}
	if decision.ErrorKind == translationJobSetupErrorKindValidationStale {
		t.Fatalf("expected validation_failed and validation_stale to stay distinct, got %#v", decision)
	}
}

func TestTranslationJobSetupServiceEvaluateCreateRequestRejectsStalePass(t *testing.T) {
	service := &TranslationJobSetupService{
		now:              func() time.Time { return time.Date(2026, 4, 27, 14, 30, 0, 0, time.UTC) },
		optionsReadModel: newServerOwnedTranslationJobSetupOptionsReadModel(),
	}

	decision, err := service.EvaluateCreateRequest(context.Background(), TranslationJobSetupCreateRequest{
		InputSourceID:    44,
		ValidationStatus: translationJobSetupValidationStatusPass,
		ValidatedAt:      time.Date(2026, 4, 27, 8, 59, 0, 0, time.UTC),
		Provider:         "openai",
		Model:            "gpt-5.4-mini",
		ExecutionMode:    "batch",
		CredentialRef:    "openai-primary",
	})
	if err != nil {
		t.Fatalf("expected stale pass to be rejected without transport error: %v", err)
	}
	if decision.CanCreate {
		t.Fatalf("expected stale validation pass to block create, got %#v", decision)
	}
	if decision.ErrorKind != translationJobSetupErrorKindValidationStale {
		t.Fatalf("expected validation_stale error kind, got %#v", decision)
	}
}

func TestTranslationJobSetupServiceEvaluateCreateRequestUsesInjectedTransportForRevalidation(t *testing.T) {
	t.Setenv(translationJobSetupXAIBaseURLEnv, "https://xai.example/v1")

	for _, testCase := range translationJobSetupProviderReachabilityCases() {
		t.Run(testCase.name, func(t *testing.T) {
			validatedAt := time.Date(2026, 4, 27, 14, 33, 0, 0, time.UTC)
			service, transport := newTranslationJobSetupServiceForProviderReachabilityTest(t, validatedAt, testCase)

			decision, err := service.EvaluateCreateRequest(context.Background(), TranslationJobSetupCreateRequest{
				InputSourceID:    44,
				ValidationStatus: translationJobSetupValidationStatusPass,
				ValidatedAt:      validatedAt,
				Provider:         testCase.provider,
				Model:            testCase.model,
				ExecutionMode:    testCase.executionMode,
				CredentialRef:    testCase.credentialRef,
			})
			if err != nil {
				t.Fatalf("expected create revalidation to use injected transport without real network: %v", err)
			}
			if !decision.CanCreate {
				t.Fatalf("expected create revalidation to pass with injected transport, got %#v", decision)
			}
			if !slices.Equal(decision.ValidationPassSlices, translationJobSetupAllSlices) {
				t.Fatalf("expected validation pass slices %#v, got %#v", translationJobSetupAllSlices, decision.ValidationPassSlices)
			}
			assertTranslationJobSetupProviderReachabilityRequest(t, transport.requests, testCase)
		})
	}
}

func TestTranslationJobSetupServiceEvaluateCreateRequestPreservesProviderModeUnsupportedErrorKind(t *testing.T) {
	validatedAt := time.Date(2026, 4, 27, 14, 34, 0, 0, time.UTC)
	service := &TranslationJobSetupService{
		now:              func() time.Time { return validatedAt },
		optionsReadModel: newServerOwnedTranslationJobSetupOptionsReadModel(),
	}

	decision, err := service.EvaluateCreateRequest(context.Background(), TranslationJobSetupCreateRequest{
		InputSourceID:    44,
		ValidationStatus: translationJobSetupValidationStatusPass,
		ValidatedAt:      validatedAt,
		Provider:         MasterPersonaProviderXAI,
		Model:            translationJobSetupModelGPT54Mini,
		ExecutionMode:    "sync",
		CredentialRef:    "xai-primary",
	})
	if err != nil {
		t.Fatalf("expected create revalidation to classify runtime mismatch without transport error: %v", err)
	}
	if decision.CanCreate {
		t.Fatalf("expected provider_mode_unsupported revalidation to reject create, got %#v", decision)
	}
	if decision.ErrorKind != translationJobSetupBlockingFailureProviderModeUnsupported {
		t.Fatalf("expected provider_mode_unsupported error kind to be preserved, got %#v", decision)
	}
}

func TestTranslationJobSetupServiceEvaluateCreateRequestPreservesProviderUnreachableErrorKind(t *testing.T) {
	validatedAt := time.Date(2026, 4, 27, 14, 35, 0, 0, time.UTC)
	t.Setenv(translationJobSetupXAIBaseURLEnv, "http://127.0.0.1:1")
	service := &TranslationJobSetupService{
		now:              func() time.Time { return validatedAt },
		optionsReadModel: newServerOwnedTranslationJobSetupOptionsReadModel(),
		secretStore: fakeTranslationJobSetupSecretStore{
			loadFunc: func(_ context.Context, key string) (string, error) {
				if key != translationJobSetupMasterPersonaSecretKey(MasterPersonaProviderXAI) {
					t.Fatalf("expected xai secret key lookup, got %q", key)
				}
				return "xai-test-key", nil
			},
		},
	}

	decision, err := service.EvaluateCreateRequest(context.Background(), TranslationJobSetupCreateRequest{
		InputSourceID:    44,
		ValidationStatus: translationJobSetupValidationStatusPass,
		ValidatedAt:      validatedAt,
		Provider:         MasterPersonaProviderXAI,
		Model:            "grok-4",
		ExecutionMode:    "sync",
		CredentialRef:    "xai-primary",
	})
	if err != nil {
		t.Fatalf("expected create re-validation to classify provider reachability without transport error: %v", err)
	}
	if decision.CanCreate {
		t.Fatalf("expected provider_unreachable re-validation to reject create, got %#v", decision)
	}
	if decision.ErrorKind != translationJobSetupErrorKindProviderUnreachable {
		t.Fatalf("expected provider_unreachable error kind to be preserved, got %#v", decision)
	}
}

func TestPersistentTranslationJobSetupServiceEvaluateCreateRequestAllowsDifferentInputWithExistingJob(t *testing.T) {
	ctx := context.Background()
	service, sourceRepository, jobRepository, closeRepositories := newSQLiteBackedTranslationJobSetupServiceForTest(t)
	defer closeRepositories()

	existingInput := createSQLiteTranslationJobSetupInputFixture(t, sourceRepository, repository.XEditExtractedDataDraft{
		SourceFilePath:    "/imports/existing-job.json",
		SourceContentHash: "existing-job-hash",
		SourceTool:        "xedit",
		TargetPluginName:  "ExistingJob.esp",
		TargetPluginType:  "esp",
		RecordCount:       3,
		ImportedAt:        time.Date(2026, 4, 27, 18, 0, 0, 0, time.UTC),
	})
	_, err := sourceRepository.CreateTranslationRecord(ctx, repository.TranslationRecordDraft{
		XEditExtractedDataID: existingInput.ID,
		FormID:               "0000AAAA",
		EditorID:             "NPC_ExistingJob",
		RecordType:           "NPC_",
	})
	if err != nil {
		t.Fatalf("expected existing input cache fixture to be created: %v", err)
	}
	_, err = jobRepository.CreateTranslationJob(ctx, repository.TranslationJobDraft{
		XEditExtractedDataID: existingInput.ID,
		JobName:              "translation-job-existing",
		State:                translationJobSetupJobStateReady,
		ProgressPercent:      0,
	})
	if err != nil {
		t.Fatalf("expected existing translation job fixture to be created: %v", err)
	}

	newInput := createSQLiteTranslationJobSetupInputFixture(t, sourceRepository, repository.XEditExtractedDataDraft{
		SourceFilePath:    "/imports/new-job.json",
		SourceContentHash: "new-job-hash",
		SourceTool:        "xedit",
		TargetPluginName:  "NewJob.esp",
		TargetPluginType:  "esp",
		RecordCount:       5,
		ImportedAt:        time.Date(2026, 4, 27, 18, 15, 0, 0, time.UTC),
	})
	_, err = sourceRepository.CreateTranslationRecord(ctx, repository.TranslationRecordDraft{
		XEditExtractedDataID: newInput.ID,
		FormID:               "0000BBBB",
		EditorID:             "NPC_NewJob",
		RecordType:           "NPC_",
	})
	if err != nil {
		t.Fatalf("expected new input cache fixture to be created: %v", err)
	}

	options, err := service.ReadOptions(ctx)
	if err != nil {
		t.Fatalf("expected read options to succeed with an unrelated existing job: %v", err)
	}
	if options.ExistingJob == nil || options.ExistingJob.JobID == 0 {
		t.Fatalf("expected read options to surface an existing job, got %#v", options.ExistingJob)
	}

	decision, err := service.EvaluateCreateRequest(ctx, TranslationJobSetupCreateRequest{
		InputSourceID:    newInput.ID,
		ValidationStatus: translationJobSetupValidationStatusPass,
		ValidatedAt:      time.Date(2026, 4, 27, 18, 20, 0, 0, time.UTC),
		Provider:         "xai",
		Model:            "grok-4",
		ExecutionMode:    "sync",
		CredentialRef:    "configured-ref",
	})
	if err != nil {
		t.Fatalf("expected create evaluation for a different input to succeed: %v", err)
	}
	if !decision.CanCreate {
		t.Fatalf("expected unrelated existing job to not block create, got %#v", decision)
	}

	created, err := service.CreateTranslationJob(ctx, TranslationJobSetupCreateRequest{
		InputSourceID:    newInput.ID,
		ValidationStatus: translationJobSetupValidationStatusPass,
		ValidatedAt:      time.Date(2026, 4, 27, 18, 20, 0, 0, time.UTC),
		Provider:         "xai",
		Model:            "grok-4",
		ExecutionMode:    "sync",
		CredentialRef:    "configured-ref",
	}, decision.ValidationPassSlices)
	if err != nil {
		t.Fatalf("expected create to succeed for a different input even when an existing job is present: %v", err)
	}
	if created.ErrorKind != "" {
		t.Fatalf("expected successful create for different input, got %#v", created)
	}
	if created.JobID == 0 {
		t.Fatalf("expected created job id to be assigned, got %#v", created)
	}
}

func TestPersistentTranslationJobSetupServiceValidateRequestPrefersInputNotFoundOverCacheMissing(t *testing.T) {
	ctx := context.Background()
	service, _, _, closeRepositories := newSQLiteBackedTranslationJobSetupServiceForTest(t)
	defer closeRepositories()

	decision, err := service.ValidateRequest(ctx, TranslationJobSetupValidationRequest{
		InputSourceID: 9999,
		Provider:      "xai",
		Model:         "grok-4",
		ExecutionMode: "sync",
		CredentialRef: "configured-ref",
	})
	if err != nil {
		t.Fatalf("expected missing input validation to classify without transport error: %v", err)
	}
	if decision.Status != translationJobSetupValidationStatusFail {
		t.Fatalf("expected missing input to fail validation, got %#v", decision)
	}
	if decision.BlockingFailureCategory == nil || *decision.BlockingFailureCategory != translationJobSetupBlockingFailureInputNotFound {
		t.Fatalf("expected input_not_found to win over cache_missing, got %#v", decision.BlockingFailureCategory)
	}
	if !slices.Equal(decision.TargetSlices, []string{"input"}) {
		t.Fatalf("expected input slice failure for missing input, got %#v", decision.TargetSlices)
	}
	if decision.CanCreate {
		t.Fatalf("expected missing input to block create, got %#v", decision)
	}
}

func TestPersistentTranslationJobSetupServiceEvaluateCreateRequestPreservesCacheMissingErrorKind(t *testing.T) {
	ctx := context.Background()
	service, sourceRepository, _, closeRepositories := newSQLiteBackedTranslationJobSetupServiceForTest(t)
	defer closeRepositories()

	input := createSQLiteTranslationJobSetupInputFixture(t, sourceRepository, repository.XEditExtractedDataDraft{
		SourceFilePath:    "/imports/cache-missing-create.json",
		SourceContentHash: "cache-missing-create-hash",
		SourceTool:        "xedit",
		TargetPluginName:  "CacheMissingCreate.esp",
		TargetPluginType:  "esp",
		RecordCount:       1,
		ImportedAt:        time.Date(2026, 4, 27, 18, 30, 0, 0, time.UTC),
	})
	_, err := sourceRepository.CreateTranslationRecord(ctx, repository.TranslationRecordDraft{
		XEditExtractedDataID: input.ID,
		FormID:               "0000CCCC",
		EditorID:             "NPC_CacheMissingCreate",
		RecordType:           "NPC_",
	})
	if err != nil {
		t.Fatalf("expected cache fixture to be created: %v", err)
	}

	cacheCleaner, ok := sourceRepository.(translationJobSetupServiceCacheCleaner)
	if !ok {
		t.Fatal("expected sqlite translation source repository to support cache deletion")
	}
	if deleteErr := cacheCleaner.DeleteTranslationCacheByXEditID(ctx, input.ID); deleteErr != nil {
		t.Fatalf("expected cache fixture to be deleted: %v", deleteErr)
	}

	decision, err := service.EvaluateCreateRequest(ctx, TranslationJobSetupCreateRequest{
		InputSourceID:    input.ID,
		ValidationStatus: translationJobSetupValidationStatusPass,
		ValidatedAt:      time.Date(2026, 4, 27, 18, 35, 0, 0, time.UTC),
		Provider:         "xai",
		Model:            "grok-4",
		ExecutionMode:    "sync",
		CredentialRef:    "configured-ref",
	})
	if err != nil {
		t.Fatalf("expected create re-validation to classify cache missing without transport error: %v", err)
	}
	if decision.CanCreate {
		t.Fatalf("expected cache_missing re-validation to reject create, got %#v", decision)
	}
	if decision.ErrorKind != "cache_missing" {
		t.Fatalf("expected cache_missing error kind to be preserved, got %#v", decision)
	}
}

func newSQLiteBackedTranslationJobSetupServiceForTest(
	t *testing.T,
) (*TranslationJobSetupService, repository.TranslationSourceRepository, repository.JobLifecycleRepository, func()) {
	t.Helper()

	databasePath := filepath.Join(t.TempDir(), "db", "translation-job-setup-service.sqlite3")
	db, err := repository.OpenSQLiteDatabase(context.Background(), databasePath)
	if err != nil {
		t.Fatalf("expected sqlite job setup service database to open: %v", err)
	}

	sourceRepository := repository.NewSQLiteTranslationSourceRepository(db)
	jobRepository := repository.NewSQLiteJobLifecycleRepository(db)
	service := NewPersistentTranslationJobSetupService(
		jobRepository,
		sourceRepository,
		nil,
		nil,
		nil,
		nil,
		repository.NewSQLiteTransactor(db),
	)
	service.now = func() time.Time { return time.Date(2026, 4, 27, 17, 30, 0, 0, time.UTC) }

	return service, sourceRepository, jobRepository, func() {
		if closeErr := db.Close(); closeErr != nil {
			t.Fatalf("expected sqlite job setup service database close to succeed: %v", closeErr)
		}
	}
}

func createSQLiteTranslationJobSetupInputFixture(
	t *testing.T,
	sourceRepository repository.TranslationSourceRepository,
	draft repository.XEditExtractedDataDraft,
) repository.XEditExtractedData {
	t.Helper()

	created, err := sourceRepository.CreateXEditExtractedData(context.Background(), draft)
	if err != nil {
		t.Fatalf("expected translation job setup input fixture to be created: %v", err)
	}
	return created
}
