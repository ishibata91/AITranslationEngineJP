package repository

import (
	"context"
	"strings"
	"testing"

	"github.com/99designs/keyring"
)

func TestProviderSettingsKeyringSecretStoreSaveLoadDelete(t *testing.T) {
	backend := keyring.NewArrayKeyring(nil)
	store, err := NewProviderSettingsKeyringSecretStoreWithKeyring(backend)
	if err != nil {
		t.Fatalf("expected provider settings keyring store creation to succeed: %v", err)
	}

	if saveErr := store.Save(context.Background(), "provider-settings:gemini", "provider-settings-secret"); saveErr != nil {
		t.Fatalf("expected provider settings secret save to succeed: %v", saveErr)
	}
	loaded, err := store.Load(context.Background(), "provider-settings:gemini")
	if err != nil {
		t.Fatalf("expected provider settings secret load to succeed: %v", err)
	}
	if loaded != "provider-settings-secret" {
		t.Fatalf("expected stored provider settings secret, got %q", loaded)
	}
	if deleteErr := store.Delete(context.Background(), "provider-settings:gemini"); deleteErr != nil {
		t.Fatalf("expected provider settings secret delete to succeed: %v", deleteErr)
	}
	loaded, err = store.Load(context.Background(), "provider-settings:gemini")
	if err != nil {
		t.Fatalf("expected provider settings secret load after delete to succeed: %v", err)
	}
	if loaded != "" {
		t.Fatalf("expected deleted provider settings secret to be empty, got %q", loaded)
	}
}

func TestProviderSettingsKeyringConfigRejectsUnsupportedBackendOverride(t *testing.T) {
	_, err := newProviderSettingsKeyringConfig("darwin", func(string) string {
		return "unsupported-backend"
	})
	if err == nil {
		t.Fatalf("expected unsupported backend override to fail")
	}
	if !strings.Contains(err.Error(), "unsupported provider settings keyring backend override") {
		t.Fatalf("expected unsupported backend override error, got %v", err)
	}
}
