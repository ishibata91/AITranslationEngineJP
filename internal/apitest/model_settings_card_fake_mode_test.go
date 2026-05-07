package apitest

import (
	"context"
	"fmt"
	"path/filepath"
	"testing"
	"time"

	"aitranslationenginejp/internal/infra/ai"
	"aitranslationenginejp/internal/infra/sqlite/dbinit"
	"aitranslationenginejp/internal/repository"
	"aitranslationenginejp/internal/service"
	"aitranslationenginejp/internal/usecase"

	"github.com/jmoiron/sqlx"
)

func TestSCN_MSCC_003_ModelSettingsFakeModeListsFakeModelUnderRealProvider(t *testing.T) {
	ctx := context.Background()
	db := openModelSettingsCardAPITestDatabase(t)
	repo := repository.NewSQLiteProviderSettingsRepository(db)
	secretStore := repository.NewInMemorySecretStore()
	modelListLoader := modelSettingsCardFakeProviderModelListAdapter{
		loader: ai.NewProviderModelListLoader(
			ai.NewTestSafeHTTPTransport(),
			ai.WithProviderModelListDeterministicProviders(),
		),
	}
	providerSettingsService := service.NewProviderSettingsService(
		repo,
		secretStore,
		repository.NewSQLiteTransactor(db),
		modelListLoader,
		modelSettingsCardNoopProviderSettingsValidator{},
		func() time.Time { return time.Date(2026, 5, 7, 10, 0, 0, 0, time.UTC) },
	)
	providerSettings := usecase.NewProviderSettingsUsecase(providerSettingsService)

	listed, err := providerSettings.ListProviderSettings(ctx, usecase.ListProviderSettingsRequest{})
	if err != nil {
		t.Fatalf("SCN-MSCC-003 expected provider settings list to load: %v", err)
	}
	assertModelSettingsCardProviderListExcludesFake(t, listed.Providers)

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

func assertModelSettingsCardProviderListExcludesFake(t *testing.T, providers []usecase.ProviderSettingsSummary) {
	t.Helper()
	if len(providers) != 3 {
		t.Fatalf("SCN-MSCC-003 expected three user-facing providers, got %#v", providers)
	}
	for _, provider := range providers {
		if provider.ProviderID == "fake" {
			t.Fatalf("SCN-MSCC-003 expected fake provider to stay hidden, got %#v", providers)
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
