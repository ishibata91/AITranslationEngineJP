package repository

import (
	"context"
	"errors"
	"path/filepath"
	"strings"
	"testing"
	"time"

	sqliteinfra "aitranslationenginejp/internal/infra/sqlite/dbinit"

	"github.com/jmoiron/sqlx"
)

func TestSQLiteProviderSettingsRepositoryPersistsRowsAcrossReopen(t *testing.T) {
	databasePath := filepath.Join(t.TempDir(), "db", "provider-settings.sqlite3")
	db := openProviderSettingsSQLiteDatabase(t, databasePath)
	repository := NewSQLiteProviderSettingsRepository(db)
	now := time.Date(2026, 5, 4, 10, 0, 0, 0, time.UTC)

	saved, err := repository.Upsert(context.Background(), ProviderSettingsUpsertDraft{
		ProviderID:            "gemini",
		Endpoint:              stringPointerForProviderSettingsTest("https://gemini.example/v1"),
		CredentialReferenceID: stringPointerForProviderSettingsTest("provider-settings:gemini"),
		CredentialState:       "configured",
		ValidationState:       "validated",
		RequestToken:          stringPointerForProviderSettingsTest("gemini|1"),
		Revision:              1,
		UpdatedAt:             now,
	})
	if err != nil {
		t.Fatalf("expected provider settings upsert to succeed: %v", err)
	}
	if saved.Endpoint == nil || *saved.Endpoint != "https://gemini.example/v1" {
		t.Fatalf("expected endpoint to be saved, got %#v", saved)
	}

	if closeErr := db.Close(); closeErr != nil {
		t.Fatalf("expected sqlite close to succeed: %v", closeErr)
	}

	reopened := openProviderSettingsSQLiteDatabase(t, databasePath)
	reopenedRepository := NewSQLiteProviderSettingsRepository(reopened)
	loaded, err := reopenedRepository.GetByProviderID(context.Background(), "gemini")
	if err != nil {
		t.Fatalf("expected provider settings load after reopen: %v", err)
	}
	if loaded.Endpoint == nil || *loaded.Endpoint != "https://gemini.example/v1" {
		t.Fatalf("expected persisted endpoint after reopen, got %#v", loaded)
	}
	if loaded.CredentialReferenceID == nil || *loaded.CredentialReferenceID != "provider-settings:gemini" {
		t.Fatalf("expected persisted credential reference after reopen, got %#v", loaded)
	}
}

func TestSQLiteProviderSettingsRepositoryUpdateValidationByRequestTokenRejectsStaleToken(t *testing.T) {
	db := openProviderSettingsSQLiteDatabase(t, filepath.Join(t.TempDir(), "db", "provider-settings.sqlite3"))
	repository := NewSQLiteProviderSettingsRepository(db)

	_, err := repository.Upsert(context.Background(), ProviderSettingsUpsertDraft{
		ProviderID:            "xai",
		Endpoint:              stringPointerForProviderSettingsTest("https://api.x.ai/v1"),
		CredentialReferenceID: stringPointerForProviderSettingsTest("provider-settings:xai"),
		CredentialState:       "configured",
		ValidationState:       "not_validated",
		RequestToken:          stringPointerForProviderSettingsTest("xai|1"),
		Revision:              1,
		UpdatedAt:             time.Date(2026, 5, 4, 10, 5, 0, 0, time.UTC),
	})
	if err != nil {
		t.Fatalf("expected provider settings upsert to succeed: %v", err)
	}

	_, matched, err := repository.UpdateValidationByRequestToken(
		context.Background(),
		"xai",
		"xai|stale",
		"validated",
		nil,
	)
	if err != nil {
		t.Fatalf("expected stale validation update check to succeed: %v", err)
	}
	if matched {
		t.Fatal("expected stale request token update to be rejected")
	}

	loaded, err := repository.GetByProviderID(context.Background(), "xai")
	if err != nil {
		t.Fatalf("expected provider settings load after stale update: %v", err)
	}
	if loaded.ValidationState != "not_validated" {
		t.Fatalf("expected validation state to remain unchanged, got %#v", loaded)
	}
}

func TestSQLiteProviderSettingsRepositoryRowContainsNoPlainAPIKey(t *testing.T) {
	db := openProviderSettingsSQLiteDatabase(t, filepath.Join(t.TempDir(), "db", "provider-settings.sqlite3"))
	repository := NewSQLiteProviderSettingsRepository(db)

	_, err := repository.Upsert(context.Background(), ProviderSettingsUpsertDraft{
		ProviderID:      "gemini",
		Endpoint:        stringPointerForProviderSettingsTest("https://gemini.example/v1"),
		CredentialState: "missing",
		ValidationState: "not_validated",
		RequestToken:    stringPointerForProviderSettingsTest("gemini|1"),
		Revision:        1,
		UpdatedAt:       time.Date(2026, 5, 4, 10, 10, 0, 0, time.UTC),
	})
	if err != nil {
		t.Fatalf("expected provider settings upsert to succeed: %v", err)
	}

	rows, queryErr := db.QueryxContext(context.Background(), "PRAGMA table_info(PROVIDER_SETTINGS)")
	if queryErr != nil {
		t.Fatalf("expected table info query to succeed: %v", queryErr)
	}
	defer func() {
		_ = rows.Close()
	}()

	for rows.Next() {
		columns := map[string]any{}
		if scanErr := rows.MapScan(columns); scanErr != nil {
			t.Fatalf("expected table info scan to succeed: %v", scanErr)
		}
		nameValue, ok := columns["name"]
		if !ok {
			t.Fatalf("expected table info to include name column, got %#v", columns)
		}
		columnName, ok := nameValue.(string)
		if !ok {
			t.Fatalf("expected table info name as string, got %#v", nameValue)
		}
		lowerName := strings.ToLower(strings.TrimSpace(columnName))
		if strings.Contains(lowerName, "api_key") || strings.Contains(lowerName, "secret") ||
			(strings.Contains(lowerName, "token") && lowerName != "request_token") {
			t.Fatalf("expected provider settings table to avoid secret columns, got %q", columnName)
		}
	}
}

func TestSQLiteProviderSettingsRepositoryReturnsNotFound(t *testing.T) {
	db := openProviderSettingsSQLiteDatabase(t, filepath.Join(t.TempDir(), "db", "provider-settings.sqlite3"))
	repository := NewSQLiteProviderSettingsRepository(db)

	_, err := repository.GetByProviderID(context.Background(), "missing")
	if !errors.Is(err, ErrProviderSettingsNotFound) {
		t.Fatalf("expected provider settings not found error, got %v", err)
	}
}

func openProviderSettingsSQLiteDatabase(t *testing.T, databasePath string) *sqlx.DB {
	t.Helper()
	db, err := sqliteinfra.OpenMasterDictionaryDatabase(context.Background(), databasePath, nil)
	if err != nil {
		t.Fatalf("expected sqlite database open to succeed: %v", err)
	}
	t.Cleanup(func() {
		_ = db.Close()
	})
	return db
}

func stringPointerForProviderSettingsTest(value string) *string {
	cloned := value
	return &cloned
}
