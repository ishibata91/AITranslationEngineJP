package service

import (
	"context"
	"fmt"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"aitranslationenginejp/internal/repository"
)

type fakeProviderSettingsValidator struct {
	validate func(context.Context, ProviderSettingsValidationProbe) error
}

func (validator fakeProviderSettingsValidator) ValidateProviderSettings(
	ctx context.Context,
	request ProviderSettingsValidationProbe,
) error {
	if validator.validate == nil {
		return nil
	}
	return validator.validate(ctx, request)
}

type fakeProviderSettingsModelListLoader struct {
	calls []struct {
		providerID string
		endpoint   string
		apiKey     string
	}
	models []ProviderSettingsModelOption
	err    error
}

func (loader *fakeProviderSettingsModelListLoader) ListProviderModelsWithEndpoint(
	_ context.Context,
	providerID string,
	endpoint string,
	apiKey string,
) ([]ProviderSettingsModelOption, error) {
	loader.calls = append(loader.calls, struct {
		providerID string
		endpoint   string
		apiKey     string
	}{providerID: providerID, endpoint: endpoint, apiKey: apiKey})
	if loader.err != nil {
		return nil, loader.err
	}
	return append([]ProviderSettingsModelOption(nil), loader.models...), nil
}

func TestProviderSettingsServiceSavePreservesSecretBoundary(t *testing.T) {
	repo, secretStore, transactor := openProviderSettingsServiceDependencies(t)
	service := NewProviderSettingsService(
		repo,
		secretStore,
		transactor,
		&fakeProviderSettingsModelListLoader{},
		fakeProviderSettingsValidator{},
		func() time.Time { return time.Date(2026, 5, 4, 11, 0, 0, 0, time.UTC) },
	)
	secretValue := string([]byte{112, 114, 111, 118, 45, 115, 101, 99, 114, 101, 116})

	saved, err := service.SaveProviderSettings(context.Background(), ProviderSettingsSaveInput{
		ProviderID:         "gemini",
		Endpoint:           stringPointerForProviderSettingsServiceTest("https://gemini.example/v1"),
		APIKeyInputPresent: true,
		APIKey:             stringPointerForProviderSettingsServiceTest(secretValue),
	})
	if err != nil {
		t.Fatalf("expected provider settings save to succeed: %v", err)
	}
	if saved.CredentialState != "configured" || saved.CredentialReferenceID == nil {
		t.Fatalf("expected configured provider settings summary, got %#v", saved)
	}
	if strings.Contains(fmt.Sprintf("%#v", saved), secretValue) {
		t.Fatalf("expected provider settings summary to redact secret, got %#v", saved)
	}

	loadedRow, err := repo.GetByProviderID(context.Background(), "gemini")
	if err != nil {
		t.Fatalf("expected saved provider settings row to load: %v", err)
	}
	if loadedRow.CredentialReferenceID == nil || strings.Contains(fmt.Sprintf("%#v", loadedRow), secretValue) {
		t.Fatalf("expected provider settings row without plaintext secret, got %#v", loadedRow)
	}

	loadedSecret, err := secretStore.Load(context.Background(), "provider-settings:gemini")
	if err != nil {
		t.Fatalf("expected provider settings secret load to succeed: %v", err)
	}
	if loadedSecret != secretValue {
		t.Fatalf("expected provider settings secret in secret store, got %q", loadedSecret)
	}
}

func TestProviderSettingsServiceSavePreservesExistingSecretWhenInputMissing(t *testing.T) {
	repo, secretStore, transactor := openProviderSettingsServiceDependencies(t)
	service := NewProviderSettingsService(
		repo,
		secretStore,
		transactor,
		&fakeProviderSettingsModelListLoader{},
		fakeProviderSettingsValidator{},
		func() time.Time { return time.Date(2026, 5, 4, 11, 10, 0, 0, time.UTC) },
	)

	_, err := service.SaveProviderSettings(context.Background(), ProviderSettingsSaveInput{
		ProviderID:         "xai",
		Endpoint:           stringPointerForProviderSettingsServiceTest("https://api.x.ai/v1"),
		APIKeyInputPresent: true,
		APIKey:             stringPointerForProviderSettingsServiceTest("xai-secret"),
	})
	if err != nil {
		t.Fatalf("expected first provider settings save to succeed: %v", err)
	}
	saved, err := service.SaveProviderSettings(context.Background(), ProviderSettingsSaveInput{
		ProviderID:         "xai",
		Endpoint:           stringPointerForProviderSettingsServiceTest("https://api.x.ai/v1/custom"),
		APIKeyInputPresent: false,
	})
	if err != nil {
		t.Fatalf("expected provider settings save without secret input to succeed: %v", err)
	}
	if saved.CredentialState != "configured" {
		t.Fatalf("expected credential state to remain configured, got %#v", saved)
	}

	loadedSecret, err := secretStore.Load(context.Background(), "provider-settings:xai")
	if err != nil {
		t.Fatalf("expected provider settings secret reload to succeed: %v", err)
	}
	if loadedSecret != "xai-secret" {
		t.Fatalf("expected existing secret to remain stored, got %q", loadedSecret)
	}
}

func TestProviderSettingsServiceResetKeepsRowAndDeletesSecret(t *testing.T) {
	repo, secretStore, transactor := openProviderSettingsServiceDependencies(t)
	service := NewProviderSettingsService(
		repo,
		secretStore,
		transactor,
		&fakeProviderSettingsModelListLoader{},
		fakeProviderSettingsValidator{},
		func() time.Time { return time.Date(2026, 5, 4, 11, 20, 0, 0, time.UTC) },
	)

	_, err := service.SaveProviderSettings(context.Background(), ProviderSettingsSaveInput{
		ProviderID:         "gemini",
		Endpoint:           stringPointerForProviderSettingsServiceTest("https://gemini.example/v1"),
		APIKeyInputPresent: true,
		APIKey:             stringPointerForProviderSettingsServiceTest("gemini-secret"),
	})
	if err != nil {
		t.Fatalf("expected provider settings save to succeed: %v", err)
	}

	reset, err := service.ResetProviderSettings(context.Background(), "gemini")
	if err != nil {
		t.Fatalf("expected provider settings reset to succeed: %v", err)
	}
	if reset.Endpoint != nil || reset.CredentialState != "missing" {
		t.Fatalf("expected reset provider settings summary, got %#v", reset)
	}
	if reset.SavedState != "partial" {
		t.Fatalf("expected row-retained saved state after reset, got %#v", reset)
	}

	loadedRow, err := repo.GetByProviderID(context.Background(), "gemini")
	if err != nil {
		t.Fatalf("expected reset provider settings row to remain readable: %v", err)
	}
	if loadedRow.Endpoint != nil || loadedRow.CredentialReferenceID != nil {
		t.Fatalf("expected reset row to keep no endpoint and no credential reference, got %#v", loadedRow)
	}
	loadedSecret, err := secretStore.Load(context.Background(), "provider-settings:gemini")
	if err != nil {
		t.Fatalf("expected provider settings secret load after reset to succeed: %v", err)
	}
	if loadedSecret != "" {
		t.Fatalf("expected provider settings secret to be deleted after reset, got %q", loadedSecret)
	}
}

func TestProviderSettingsServiceValidateRejectsDelayedResponseByRequestToken(t *testing.T) {
	repo, secretStore, transactor := openProviderSettingsServiceDependencies(t)
	service := NewProviderSettingsService(
		repo,
		secretStore,
		transactor,
		&fakeProviderSettingsModelListLoader{},
		fakeProviderSettingsValidator{
			validate: func(_ context.Context, request ProviderSettingsValidationProbe) error {
				row, err := repo.GetByProviderID(context.Background(), request.ProviderID)
				if err != nil {
					return fmt.Errorf("load provider settings row during delayed validation test: %w", err)
				}
				_, err = repo.Upsert(context.Background(), repository.ProviderSettingsUpsertDraft{
					ProviderID:            request.ProviderID,
					Endpoint:              stringPointerForProviderSettingsServiceTest(request.Endpoint + "/next"),
					CredentialReferenceID: row.CredentialReferenceID,
					CredentialState:       row.CredentialState,
					ValidationState:       "not_validated",
					RequestToken:          stringPointerForProviderSettingsServiceTest("gemini|2"),
					Revision:              2,
					UpdatedAt:             time.Date(2026, 5, 4, 11, 31, 0, 0, time.UTC),
				})
				if err != nil {
					return fmt.Errorf("save newer provider settings row during delayed validation test: %w", err)
				}
				return nil
			},
		},
		func() time.Time { return time.Date(2026, 5, 4, 11, 30, 0, 0, time.UTC) },
	)

	saved, err := service.SaveProviderSettings(context.Background(), ProviderSettingsSaveInput{
		ProviderID:         "gemini",
		Endpoint:           stringPointerForProviderSettingsServiceTest("https://gemini.example/v1"),
		APIKeyInputPresent: true,
		APIKey:             stringPointerForProviderSettingsServiceTest("gemini-secret"),
	})
	if err != nil {
		t.Fatalf("expected provider settings save to succeed: %v", err)
	}

	validated, err := service.ValidateProviderSettings(context.Background(), ProviderSettingsValidateInput{
		ProviderID:            "gemini",
		Endpoint:              saved.Endpoint,
		CredentialState:       saved.CredentialState,
		CredentialReferenceID: saved.CredentialReferenceID,
		RequestToken:          *saved.RequestToken,
	})
	if err != nil {
		t.Fatalf("expected provider settings validate to finish with stale result: %v", err)
	}
	if validated.FailureKind == nil || *validated.FailureKind != "validation_stale" {
		t.Fatalf("expected delayed validation to be rejected as stale, got %#v", validated)
	}

	current, err := repo.GetByProviderID(context.Background(), "gemini")
	if err != nil {
		t.Fatalf("expected provider settings row after delayed validation: %v", err)
	}
	if current.RequestToken == nil || *current.RequestToken != "gemini|2" {
		t.Fatalf("expected current provider settings row to keep newer request token, got %#v", current)
	}
	if current.ValidationState != "not_validated" {
		t.Fatalf("expected delayed validation not to overwrite current validation state, got %#v", current)
	}
}

func TestProviderSettingsServiceListProviderModelsGatesCredentialAndEndpoint(t *testing.T) {
	repo, secretStore, transactor := openProviderSettingsServiceDependencies(t)
	loader := &fakeProviderSettingsModelListLoader{
		models: []ProviderSettingsModelOption{{ModelID: "gemini-2.5-pro", Label: "Gemini 2.5 Pro"}},
	}
	service := NewProviderSettingsService(
		repo,
		secretStore,
		transactor,
		loader,
		fakeProviderSettingsValidator{},
		func() time.Time { return time.Date(2026, 5, 4, 11, 40, 0, 0, time.UTC) },
	)

	blocked, err := service.ListProviderModels(context.Background(), ProviderSettingsModelListInput{
		ProviderID:            "gemini",
		Endpoint:              nil,
		CredentialState:       "missing",
		CredentialReferenceID: nil,
		RequestToken:          "",
	})
	if err != nil {
		t.Fatalf("expected blocked provider model list call to succeed: %v", err)
	}
	if blocked.State != "failed" && blocked.State != "not_requested" {
		t.Fatalf("expected blocked provider model list to avoid transport, got %#v", blocked)
	}
	if len(loader.calls) != 0 {
		t.Fatalf("expected blocked provider model list to avoid transport call, got %#v", loader.calls)
	}

	saved, err := service.SaveProviderSettings(context.Background(), ProviderSettingsSaveInput{
		ProviderID:         "gemini",
		Endpoint:           stringPointerForProviderSettingsServiceTest("https://gemini.example/v1"),
		APIKeyInputPresent: true,
		APIKey:             stringPointerForProviderSettingsServiceTest("gemini-secret"),
	})
	if err != nil {
		t.Fatalf("expected provider settings save to succeed: %v", err)
	}
	ready, err := service.ListProviderModels(context.Background(), ProviderSettingsModelListInput{
		ProviderID:            "gemini",
		Endpoint:              saved.Endpoint,
		CredentialState:       saved.CredentialState,
		CredentialReferenceID: saved.CredentialReferenceID,
		RequestToken:          *saved.RequestToken,
	})
	if err != nil {
		t.Fatalf("expected provider model list to succeed: %v", err)
	}
	if ready.State != "ready" || len(ready.Models) != 1 {
		t.Fatalf("expected ready provider model list, got %#v", ready)
	}
	if len(loader.calls) != 1 || loader.calls[0].apiKey != "gemini-secret" {
		t.Fatalf("expected transport call with secret only after readiness gate, got %#v", loader.calls)
	}
	if strings.Contains(fmt.Sprintf("%#v", ready), "gemini-secret") {
		t.Fatalf("expected provider model list result to avoid secret output, got %#v", ready)
	}
}

func openProviderSettingsServiceDependencies(
	t *testing.T,
) (*repository.SQLiteProviderSettingsRepository, *repository.InMemorySecretStore, repository.Transactor) {
	t.Helper()
	db, err := repository.OpenSQLiteDatabase(context.Background(), filepath.Join(t.TempDir(), "provider-settings.sqlite3"))
	if err != nil {
		t.Fatalf("expected sqlite database open to succeed: %v", err)
	}
	t.Cleanup(func() {
		_ = db.Close()
	})
	return repository.NewSQLiteProviderSettingsRepository(db), repository.NewInMemorySecretStore(), repository.NewSQLiteTransactor(db)
}

func stringPointerForProviderSettingsServiceTest(value string) *string {
	cloned := value
	return &cloned
}
