package apitest

import (
	"context"
	"fmt"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"aitranslationenginejp/internal/infra/ai"
	"aitranslationenginejp/internal/infra/sqlite/dbinit"
	"aitranslationenginejp/internal/repository"
	"aitranslationenginejp/internal/service"
	"aitranslationenginejp/internal/usecase"

	"github.com/jmoiron/sqlx"
)

func TestSCN_DFSS_005_FakeSecretStoreDoesNotAddFakeProviderToProviderList(t *testing.T) {
	ctx := context.Background()
	db := openModelSettingsCardAPITestDatabase(t)
	providerSettings := newModelSettingsCardProviderSettingsUsecase(t, db)

	listed, err := providerSettings.ListProviderSettings(ctx, usecase.ListProviderSettingsRequest{})
	if err != nil {
		t.Fatalf("SCN-DFSS-005 expected provider settings list to load: %v", err)
	}
	assertModelSettingsCardProviderListExcludesFake(t, "SCN-DFSS-005", listed.Providers)
	assertModelSettingsCardPublicTextExcludesSecretBackendValues(t, "SCN-DFSS-005", fmt.Sprintf("%#v", listed))
}

func TestSCN_DFSS_007_FakeSecretStoreEvidenceKeepsSecretValuesOutOfPublicResultAndLog(t *testing.T) {
	capture := startObservabilityLogCapture(t)
	ctx := context.Background()
	db := openModelSettingsCardAPITestDatabase(t)
	validator := &modelSettingsCardRecordingProviderSettingsValidator{}
	providerSettings := newModelSettingsCardProviderSettingsUsecaseWithValidator(t, db, validator)
	endpoint := "https://gemini.example/v1beta"

	saved, err := providerSettings.SaveProviderSettings(ctx, usecase.SaveProviderSettingsRequest{
		ProviderID: usecase.ProviderSettingsProviderIDGemini,
		Endpoint:   modelSettingsCardStringPointer(endpoint),
	})
	if err != nil {
		t.Fatalf("SCN-DFSS-007 expected provider settings save without credential input to succeed: %v", err)
	}
	validated, err := providerSettings.ValidateProviderSettings(ctx, usecase.ValidateProviderSettingsRequest{
		ProviderID:            saved.Provider.ProviderID,
		Endpoint:              saved.Provider.Endpoint,
		CredentialState:       saved.Provider.CredentialState,
		CredentialReferenceID: saved.Provider.CredentialReferenceID,
		RequestToken:          modelSettingsCardPointerValue(saved.Provider.RequestToken),
	})
	if err != nil {
		t.Fatalf("SCN-DFSS-007 expected validation to return redacted failure payload: %v", err)
	}
	if validated.FailureKind == nil || *validated.FailureKind != usecase.ProviderSettingsErrorKindCredentialMissing {
		t.Fatalf("SCN-DFSS-007 expected credential_missing failure kind, got %#v", validated)
	}
	if validator.called {
		t.Fatal("SCN-DFSS-007 expected missing credential to stop before provider request capture")
	}

	payload := capture.requireEvent(t, "provider_settings_validation", "failed", "credential_missing")
	if payload["where"] != "provider_settings.service" || payload["provider"] != "gemini" {
		t.Fatalf("SCN-DFSS-007 expected provider settings boundary log payload, got %#v", payload)
	}
	assertModelSettingsCardPublicTextExcludesSecretBackendValues(t, "SCN-DFSS-007", fmt.Sprintf("%#v\n%s", validated, capture.buffer.String()))
}

func TestSCN_MSCC_003_ModelSettingsFakeModeListsFakeModelUnderRealProvider(t *testing.T) {
	ctx := context.Background()
	db := openModelSettingsCardAPITestDatabase(t)
	repo := repository.NewSQLiteProviderSettingsRepository(db)
	providerSettings := newModelSettingsCardProviderSettingsUsecase(t, db)

	listed, err := providerSettings.ListProviderSettings(ctx, usecase.ListProviderSettingsRequest{})
	if err != nil {
		t.Fatalf("SCN-MSCC-003 expected provider settings list to load: %v", err)
	}
	assertModelSettingsCardProviderListExcludesFake(t, "SCN-MSCC-003", listed.Providers)

	endpoint := "https://gemini.example/v1beta"
	saved, err := providerSettings.SaveProviderSettingsWithSecret(ctx, usecase.SaveProviderSettingsRequest{
		ProviderID:         usecase.ProviderSettingsProviderIDGemini,
		Endpoint:           modelSettingsCardStringPointer(endpoint),
		APIKeyInputPresent: true,
	}, "test-safe-secret")
	if err != nil {
		t.Fatalf("SCN-MSCC-003 expected real provider settings save to succeed: %v", err)
	}
	if saved.Provider.ProviderID != usecase.ProviderSettingsProviderIDGemini {
		t.Fatalf("SCN-MSCC-003 expected saved provider to remain gemini, got %#v", saved.Provider)
	}

	models, err := providerSettings.ListProviderModels(ctx, usecase.ListProviderModelsRequest{
		ProviderID:            saved.Provider.ProviderID,
		Endpoint:              saved.Provider.Endpoint,
		CredentialState:       saved.Provider.CredentialState,
		CredentialReferenceID: saved.Provider.CredentialReferenceID,
		RequestToken:          modelSettingsCardPointerValue(saved.Provider.RequestToken),
	})
	if err != nil {
		t.Fatalf("SCN-MSCC-003 expected fake mode model list to load through real provider: %v", err)
	}
	if models.ProviderID != usecase.ProviderSettingsProviderIDGemini {
		t.Fatalf("SCN-MSCC-003 expected model list provider to remain gemini, got %#v", models)
	}
	if models.State != usecase.ProviderSettingsModelListStateReady {
		t.Fatalf("SCN-MSCC-003 expected ready model list state, got %#v", models)
	}
	if len(models.Models) != 1 || models.Models[0].ModelID != ai.FakeModelID {
		t.Fatalf("SCN-MSCC-003 expected fake-model under gemini response, got %#v", models.Models)
	}

	row, err := repo.GetByProviderID(ctx, string(usecase.ProviderSettingsProviderIDGemini))
	if err != nil {
		t.Fatalf("SCN-MSCC-003 expected saved provider settings row to load: %v", err)
	}
	if row.ProviderID != string(usecase.ProviderSettingsProviderIDGemini) {
		t.Fatalf("SCN-MSCC-003 expected saved row provider to remain gemini, got %#v", row)
	}
	if _, err := repo.GetByProviderID(ctx, "fake"); err == nil {
		t.Fatal("SCN-MSCC-003 expected fake provider settings row not to be saved")
	}
}

func newModelSettingsCardProviderSettingsUsecase(
	t *testing.T,
	db *sqlx.DB,
) *usecase.ProviderSettingsAppUsecase {
	t.Helper()
	return newModelSettingsCardProviderSettingsUsecaseWithValidator(
		t,
		db,
		modelSettingsCardNoopProviderSettingsValidator{},
	)
}

func newModelSettingsCardProviderSettingsUsecaseWithValidator(
	t *testing.T,
	db *sqlx.DB,
	validator modelSettingsCardProviderSettingsValidator,
) *usecase.ProviderSettingsAppUsecase {
	t.Helper()

	providerSettingsService := service.NewProviderSettingsService(
		repository.NewSQLiteProviderSettingsRepository(db),
		repository.NewInMemorySecretStore(),
		repository.NewSQLiteTransactor(db),
		modelSettingsCardFakeProviderModelListAdapter{
			loader: ai.NewProviderModelListLoader(
				ai.NewTestSafeHTTPTransport(),
				ai.WithProviderModelListDeterministicProviders(),
			),
		},
		validator,
		func() time.Time { return time.Date(2026, 5, 7, 10, 0, 0, 0, time.UTC) },
	)
	return usecase.NewProviderSettingsUsecase(providerSettingsService)
}

type modelSettingsCardProviderSettingsValidator interface {
	ValidateProviderSettings(context.Context, service.ProviderSettingsValidationProbe) error
}

type modelSettingsCardFakeProviderModelListAdapter struct {
	loader *ai.ProviderModelListLoader
}

func (adapter modelSettingsCardFakeProviderModelListAdapter) ListProviderModelsWithEndpoint(
	ctx context.Context,
	providerID string,
	endpoint string,
	apiKey string,
) ([]service.ProviderSettingsModelOption, error) {
	models, err := adapter.loader.ListProviderModelsWithEndpoint(ctx, providerID, endpoint, apiKey)
	if err != nil {
		return nil, fmt.Errorf("list model settings card fake provider models: %w", err)
	}
	result := make([]service.ProviderSettingsModelOption, 0, len(models))
	for _, model := range models {
		result = append(result, service.ProviderSettingsModelOption{
			ModelID: model.ModelID,
			Label:   model.Label,
		})
	}
	return result, nil
}

type modelSettingsCardNoopProviderSettingsValidator struct{}

func (modelSettingsCardNoopProviderSettingsValidator) ValidateProviderSettings(
	context.Context,
	service.ProviderSettingsValidationProbe,
) error {
	return nil
}

type modelSettingsCardRecordingProviderSettingsValidator struct {
	called bool
}

func (validator *modelSettingsCardRecordingProviderSettingsValidator) ValidateProviderSettings(
	context.Context,
	service.ProviderSettingsValidationProbe,
) error {
	validator.called = true
	return nil
}

func openModelSettingsCardAPITestDatabase(t *testing.T) *sqlx.DB {
	t.Helper()
	databasePath := filepath.Join(t.TempDir(), "model-settings-card.sqlite3")
	db, err := dbinit.OpenMasterDictionaryDatabase(context.Background(), databasePath, nil)
	if err != nil {
		t.Fatalf("expected model settings card api test database to open: %v", err)
	}
	t.Cleanup(func() {
		_ = db.Close()
	})
	return db
}

func assertModelSettingsCardProviderListExcludesFake(
	t *testing.T,
	scenarioID string,
	providers []usecase.ProviderSettingsSummary,
) {
	t.Helper()
	if len(providers) != 3 {
		t.Fatalf("%s expected three user-facing providers, got %#v", scenarioID, providers)
	}
	for _, provider := range providers {
		if provider.ProviderID == "fake" {
			t.Fatalf("%s expected fake provider to stay hidden, got %#v", scenarioID, providers)
		}
	}
}

func assertModelSettingsCardPublicTextExcludesSecretBackendValues(
	t *testing.T,
	scenarioID string,
	text string,
) {
	t.Helper()

	normalized := strings.ToLower(text)
	for _, forbidden := range []string{
		"in-memory",
		"provider-settings:",
		"secret",
		"apikey",
		"credentialinput",
		"rawrequest",
		"rawresponse",
	} {
		if strings.Contains(normalized, forbidden) {
			t.Fatalf("%s expected public API result and evidence text to exclude secret boundary value %q", scenarioID, forbidden)
		}
	}
}

func modelSettingsCardStringPointer(value string) *string {
	cloned := value
	return &cloned
}

func modelSettingsCardPointerValue(value *string) string {
	if value == nil {
		return ""
	}
	return *value
}
