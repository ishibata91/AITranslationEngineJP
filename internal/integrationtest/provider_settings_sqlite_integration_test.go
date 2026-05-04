package integrationtest

import (
	"context"
	"path/filepath"
	"testing"
	"time"

	"aitranslationenginejp/internal/infra/sqlite/dbinit"
	"aitranslationenginejp/internal/repository"

	"github.com/jmoiron/sqlx"
)

func TestProviderSettingsSQLiteIntegrationPersistsConfiguredStateAcrossReopen(t *testing.T) {
	databasePath := filepath.Join(t.TempDir(), "provider-settings.sqlite3")
	firstDB := openProviderSettingsIntegrationDatabase(t, databasePath)
	firstRepo := repository.NewSQLiteProviderSettingsRepository(firstDB)

	saved, err := firstRepo.Upsert(context.Background(), repository.ProviderSettingsUpsertDraft{
		ProviderID:            "gemini",
		Endpoint:              providerSettingsIntegrationStringPointer("https://gemini.example/v1"),
		CredentialReferenceID: providerSettingsIntegrationStringPointer("provider-settings:gemini"),
		CredentialState:       "configured",
		ValidationState:       "validated",
		RequestToken:          providerSettingsIntegrationStringPointer("gemini|1"),
		Revision:              1,
		UpdatedAt:             time.Date(2026, 5, 4, 12, 0, 0, 0, time.UTC),
	})
	if err != nil {
		t.Fatalf("expected provider settings integration upsert to succeed: %v", err)
	}
	if saved.RequestToken == nil {
		t.Fatalf("expected request token after upsert, got %#v", saved)
	}

	reopenedDB := openProviderSettingsIntegrationDatabase(t, databasePath)
	reopenedRepo := repository.NewSQLiteProviderSettingsRepository(reopenedDB)
	providers, err := reopenedRepo.List(context.Background())
	if err != nil {
		t.Fatalf("expected provider settings integration list after reopen: %v", err)
	}
	if len(providers) != 1 {
		t.Fatalf("expected one persisted provider settings row after reopen, got %#v", providers)
	}
	gemini := providers[0]
	if gemini.ProviderID != "gemini" || gemini.Endpoint == nil || *gemini.Endpoint != "https://gemini.example/v1" || gemini.CredentialState != "configured" {
		t.Fatalf("expected persisted provider settings row after reopen, got %#v", gemini)
	}
}

func openProviderSettingsIntegrationDatabase(t *testing.T, databasePath string) *sqlx.DB {
	t.Helper()
	db, err := dbinit.OpenMasterDictionaryDatabase(context.Background(), databasePath, nil)
	if err != nil {
		t.Fatalf("expected integration sqlite open to succeed: %v", err)
	}
	t.Cleanup(func() {
		_ = db.Close()
	})
	return db
}

func providerSettingsIntegrationStringPointer(value string) *string {
	cloned := value
	return &cloned
}
