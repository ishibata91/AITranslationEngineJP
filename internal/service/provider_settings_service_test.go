package service

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
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

func TestLogProviderBoundaryFailureUsesSafePayloadOnly(t *testing.T) {
	var buffer bytes.Buffer
	previous := slog.Default()
	slog.SetDefault(slog.New(slog.NewJSONHandler(&buffer, nil)))
	t.Cleanup(func() { slog.SetDefault(previous) })

	logProviderBoundaryFailure(context.Background(), "provider_execution_settings", "provider_settings.service", "gemini", "")

	var payload map[string]any
	if err := json.Unmarshal(buffer.Bytes(), &payload); err != nil {
		t.Fatalf("unmarshal provider settings payload: %v", err)
	}
	if payload["event"] != "provider_execution_settings" || payload["where"] != "provider_settings.service" || payload["result"] != "failed" {
		t.Fatalf("unexpected provider settings payload: %#v", payload)
	}
	if payload["reason"] != "unknown" {
		t.Fatalf("expected normalized reason, got %#v", payload)
	}
	if payload["provider"] != "gemini" {
		t.Fatalf("expected provider id, got %#v", payload)
	}
	forbidden := []string{"api_key", "endpoint", "raw_request", "raw_response", "full_path", "trace_id"}
	for _, key := range forbidden {
		if _, ok := payload[key]; ok {
			t.Fatalf("forbidden key %q in payload: %#v", key, payload)
		}
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
	if len(loader.calls) != 1 || loader.calls[0].apiKey == "" {
		t.Fatalf("expected transport call with secret only after readiness gate, got %#v", loader.calls)
	}
	if strings.Contains(fmt.Sprintf("%#v", ready), "gemini-secret") {
		t.Fatalf("expected provider model list result to avoid secret output, got %#v", ready)
	}
}

func TestProviderSettingsServiceListProviderModelsDoesNotBypassCredentialAndEndpointForLoaderMarker(t *testing.T) {
	repo, secretStore, transactor := openProviderSettingsServiceDependencies(t)
	loader := &fakeProviderSettingsModelListLoader{
		models: []ProviderSettingsModelOption{{ModelID: "gemini-test-safe", Label: "Gemini Test Safe"}},
	}
	service := NewProviderSettingsService(
		repo,
		secretStore,
		transactor,
		loader,
		fakeProviderSettingsValidator{},
		func() time.Time { return time.Date(2026, 5, 6, 10, 0, 0, 0, time.UTC) },
	)

	listed, err := service.ListProviderModels(context.Background(), ProviderSettingsModelListInput{
		ProviderID:            "gemini",
		Endpoint:              nil,
		CredentialState:       "missing",
		CredentialReferenceID: nil,
		RequestToken:          "gemini|test-safe",
	})
	if err != nil {
		t.Fatalf("expected provider model list preflight to return a result: %v", err)
	}
	if listed.State != "failed" || listed.FailureKind == nil || *listed.FailureKind != "validation_stale" {
		t.Fatalf("expected normal stale preflight gate before loader call, got %#v", listed)
	}
	if len(loader.calls) != 0 {
		t.Fatalf("expected no loader call before normal preflight passes, got %#v", loader.calls)
	}
}

func TestProviderSettingsServiceApplyExecutionSnapshotDoesNotBypassCredentialAndEndpoint(t *testing.T) {
	repo, secretStore, transactor := openProviderSettingsServiceDependencies(t)
	service := NewProviderSettingsService(
		repo,
		secretStore,
		transactor,
		&fakeProviderSettingsModelListLoader{},
		fakeProviderSettingsValidator{},
		func() time.Time { return time.Date(2026, 5, 6, 10, 10, 0, 0, time.UTC) },
	)
	spec, err := providerSettingsSpec("gemini")
	if err != nil {
		t.Fatalf("expected provider settings spec to resolve: %v", err)
	}
	result := ProviderSettingsResolveResult{
		ProviderID:            "gemini",
		Endpoint:              nil,
		CredentialReferenceID: stringPointerForProviderSettingsServiceTest("provider-settings:gemini"),
		CredentialState:       "missing",
		ErrorKind:             stringPointerForProviderSettingsServiceTest("endpoint_missing"),
	}

	err = service.applyProviderSettingsExecutionSnapshot(
		context.Background(),
		spec,
		ProviderSettingsResolveInput{AllowSecretSnapshot: true},
		ProviderSettingsSummary{
			ProviderID:            "gemini",
			Endpoint:              nil,
			CredentialReferenceID: nil,
			CredentialState:       "missing",
		},
		&result,
	)
	if err != nil {
		t.Fatalf("expected provider execution snapshot apply to return without secret read: %v", err)
	}
	if result.Endpoint != nil {
		t.Fatalf("expected nil endpoint to remain nil, got %#v", result)
	}
	if result.CredentialReferenceID == nil {
		t.Fatalf("expected existing credential reference to remain, got %#v", result)
	}
	if result.CredentialState != "missing" {
		t.Fatalf("expected existing credential state to remain, got %#v", result)
	}
	if result.ErrorKind == nil || *result.ErrorKind != "endpoint_missing" {
		t.Fatalf("expected existing error kind to remain, got %#v", result)
	}
}

func TestProviderSettingsResolveResultKeepsRealProviderReadinessGates(t *testing.T) {
	spec, err := providerSettingsSpec("gemini")
	if err != nil {
		t.Fatalf("expected provider settings spec to resolve: %v", err)
	}
	missingEndpoint := providerSettingsResolveResult(
		spec,
		ProviderSettingsResolveInput{
			ConsumerID: "translation-job-setup",
			Selection: ProviderSettingsResolveSelection{
				ProviderID:      "gemini",
				Model:           "gemini-2.5-pro",
				ExecutionMethod: "standard",
				UseBatchAPI:     false,
			},
			AllowSecretSnapshot: true,
		},
		ProviderSettingsSummary{
			ProviderID:            "gemini",
			Endpoint:              nil,
			CredentialReferenceID: nil,
			CredentialState:       "missing",
		},
	)
	if missingEndpoint.ErrorKind == nil || *missingEndpoint.ErrorKind != "endpoint_missing" {
		t.Fatalf("expected endpoint_missing gate for real provider execution resolve, got %#v", missingEndpoint)
	}

	missingCredential := providerSettingsResolveResult(
		spec,
		ProviderSettingsResolveInput{
			ConsumerID: "translation-job-setup",
			Selection: ProviderSettingsResolveSelection{
				ProviderID:      "gemini",
				Model:           "gemini-2.5-pro",
				ExecutionMethod: "standard",
				UseBatchAPI:     false,
			},
			AllowSecretSnapshot: true,
		},
		ProviderSettingsSummary{
			ProviderID:            "gemini",
			Endpoint:              stringPointerForProviderSettingsServiceTest("https://gemini.example/v1"),
			CredentialReferenceID: nil,
			CredentialState:       "missing",
		},
	)
	if missingCredential.ErrorKind == nil || *missingCredential.ErrorKind != "credential_missing" {
		t.Fatalf("expected credential_missing gate for real provider execution resolve, got %#v", missingCredential)
	}
	if missingCredential.CredentialState != "missing" || missingCredential.CredentialReferenceID != nil {
		t.Fatalf("expected real provider missing credential state to remain gated, got %#v", missingCredential)
	}
}

func TestProviderSettingsServiceResolveExecutionDoesNotBypassMissingSettings(t *testing.T) {
	repo, secretStore, transactor := openProviderSettingsServiceDependencies(t)
	service := NewProviderSettingsService(
		repo,
		secretStore,
		transactor,
		&fakeProviderSettingsModelListLoader{},
		fakeProviderSettingsValidator{},
		func() time.Time { return time.Date(2026, 5, 6, 10, 20, 0, 0, time.UTC) },
	)

	resolved, err := service.ResolveProviderExecutionSettings(context.Background(), ProviderSettingsResolveInput{
		ConsumerID: "translation-job-setup",
		Selection: ProviderSettingsResolveSelection{
			ProviderID:      "gemini",
			Model:           "gemini-test-safe",
			ExecutionMethod: "standard",
			UseBatchAPI:     false,
		},
		AllowSecretSnapshot: true,
	})
	if err != nil {
		t.Fatalf("expected provider execution resolve to return gated result: %v", err)
	}
	if resolved.ErrorKind == nil || *resolved.ErrorKind != "credential_missing" {
		t.Fatalf("expected missing credential gate to remain, got %#v", resolved)
	}
}

func openProviderSettingsServiceDependencies(
	t *testing.T,
) (*repository.SQLiteProviderSettingsRepository, repository.ProviderSettingsSecretStore, repository.Transactor) {
	t.Helper()
	db, err := repository.OpenSQLiteDatabase(context.Background(), filepath.Join(t.TempDir(), "provider-settings.sqlite3"))
	if err != nil {
		t.Fatalf("expected sqlite database open to succeed: %v", err)
	}
	t.Cleanup(func() {
		_ = db.Close()
	})
	return repository.NewSQLiteProviderSettingsRepository(db), repository.NewFakeSecretStore(), repository.NewSQLiteTransactor(db)
}

func stringPointerForProviderSettingsServiceTest(value string) *string {
	cloned := value
	return &cloned
}
