package repository

import (
	"context"
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/99designs/keyring"
)

func TestMasterPersonaKeyringSecretStoreSaveLoadDelete(t *testing.T) {
	backend := keyring.NewArrayKeyring(nil)
	store, err := NewMasterPersonaKeyringSecretStoreWithKeyring(backend)
	if err != nil {
		t.Fatalf("expected keyring-backed secret store creation to succeed: %v", err)
	}

	if saveErr := store.Save(context.Background(), "master-persona:gemini", "saved-api-key"); saveErr != nil {
		t.Fatalf("expected keyring-backed secret save to succeed: %v", saveErr)
	}

	loaded, err := store.Load(context.Background(), "master-persona:gemini")
	if err != nil {
		t.Fatalf("expected keyring-backed secret load to succeed: %v", err)
	}
	if loaded != "saved-api-key" {
		t.Fatalf("expected saved api key to load from keyring backend, got %q", loaded)
	}

	if deleteErr := store.Delete(context.Background(), "master-persona:gemini"); deleteErr != nil {
		t.Fatalf("expected keyring-backed secret delete to succeed: %v", deleteErr)
	}
	if deleteErr := store.Delete(context.Background(), "master-persona:gemini"); deleteErr != nil {
		t.Fatalf("expected missing key delete to succeed: %v", deleteErr)
	}
	if deleteErr := store.Delete(context.Background(), " "); deleteErr != nil {
		t.Fatalf("expected blank key delete to succeed: %v", deleteErr)
	}

	loaded, err = store.Load(context.Background(), "master-persona:gemini")
	if err != nil {
		t.Fatalf("expected deleted keyring-backed secret load to succeed: %v", err)
	}
	if loaded != "" {
		t.Fatalf("expected deleted secret to return empty string, got %q", loaded)
	}
}

func TestMasterPersonaKeyringSecretStoreRequiresNonEmptyKeyOnSave(t *testing.T) {
	backend := keyring.NewArrayKeyring(nil)
	store, err := NewMasterPersonaKeyringSecretStoreWithKeyring(backend)
	if err != nil {
		t.Fatalf("expected keyring-backed secret store creation to succeed: %v", err)
	}

	if err := store.Save(context.Background(), " ", "saved-api-key"); err == nil {
		t.Fatalf("expected keyring-backed secret save with empty key to fail")
	}
}

func TestMasterPersonaKeyringSecretStoreLoadBlankKeyReturnsEmpty(t *testing.T) {
	backend := keyring.NewArrayKeyring(nil)
	store, err := NewMasterPersonaKeyringSecretStoreWithKeyring(backend)
	if err != nil {
		t.Fatalf("expected keyring-backed secret store creation to succeed: %v", err)
	}

	loaded, err := store.Load(context.Background(), " ")
	if err != nil {
		t.Fatalf("expected blank key load to succeed: %v", err)
	}
	if loaded != "" {
		t.Fatalf("expected blank key load to return empty string, got %q", loaded)
	}
}

func TestMasterPersonaKeyringSecretStoreWrapsBackendErrors(t *testing.T) {
	loadErr := errors.New("load failure")
	saveErr := errors.New("save failure")
	deleteErr := errors.New("delete failure")
	backend := &failingMasterPersonaKeyring{
		getErr:    loadErr,
		setErr:    saveErr,
		removeErr: deleteErr,
	}
	store, err := NewMasterPersonaKeyringSecretStoreWithKeyring(backend)
	if err != nil {
		t.Fatalf("expected keyring-backed secret store creation to succeed: %v", err)
	}

	_, err = store.Load(context.Background(), "master-persona:provider")
	if !errors.Is(err, loadErr) {
		t.Fatalf("expected load error to be wrapped, got %v", err)
	}

	if err := store.Save(context.Background(), "master-persona:provider", "value"); !errors.Is(err, saveErr) {
		t.Fatalf("expected save error to be wrapped, got %v", err)
	}

	if err := store.Delete(context.Background(), "master-persona:provider"); !errors.Is(err, deleteErr) {
		t.Fatalf("expected delete error to be wrapped, got %v", err)
	}
}

func TestNewMasterPersonaKeyringSecretStoreWithKeyringRejectsNilBackend(t *testing.T) {
	store, err := NewMasterPersonaKeyringSecretStoreWithKeyring(nil)
	if err == nil {
		t.Fatalf("expected nil backend constructor to fail")
	}
	if store != nil {
		t.Fatalf("expected nil store on constructor failure")
	}
}

func TestNewMasterPersonaKeyringSecretStoreRejectsNilOpener(t *testing.T) {
	store, err := newMasterPersonaKeyringSecretStore(nil, "darwin", nil)
	if err == nil {
		t.Fatalf("expected nil opener to fail")
	}
	if store != nil {
		t.Fatalf("expected nil store when opener is nil")
	}
}

func TestNewMasterPersonaKeyringSecretStoreDefaultBackendDarwin(t *testing.T) {
	capturedConfig := keyring.Config{}
	store, err := newMasterPersonaKeyringSecretStore(
		func(config keyring.Config) (keyring.Keyring, error) {
			capturedConfig = config
			return keyring.NewArrayKeyring(nil), nil
		},
		"darwin",
		envValues(map[string]string{}),
	)
	if err != nil {
		t.Fatalf("expected darwin default backend to be supported: %v", err)
	}
	if store == nil {
		t.Fatalf("expected store instance")
	}
	if len(capturedConfig.AllowedBackends) != 1 || capturedConfig.AllowedBackends[0] != keyring.KeychainBackend {
		t.Fatalf("expected keychain backend, got %#v", capturedConfig.AllowedBackends)
	}
}

func TestNewMasterPersonaKeyringSecretStoreDefaultBackendWindows(t *testing.T) {
	capturedConfig := keyring.Config{}
	store, err := newMasterPersonaKeyringSecretStore(
		func(config keyring.Config) (keyring.Keyring, error) {
			capturedConfig = config
			return keyring.NewArrayKeyring(nil), nil
		},
		"windows",
		envValues(map[string]string{}),
	)
	if err != nil {
		t.Fatalf("expected windows default backend to be supported: %v", err)
	}
	if store == nil {
		t.Fatalf("expected store instance")
	}
	if len(capturedConfig.AllowedBackends) != 1 || capturedConfig.AllowedBackends[0] != keyring.WinCredBackend {
		t.Fatalf("expected wincred backend, got %#v", capturedConfig.AllowedBackends)
	}
	if capturedConfig.WinCredPrefix != masterPersonaWinCredPrefix {
		t.Fatalf("expected wincred prefix %q, got %q", masterPersonaWinCredPrefix, capturedConfig.WinCredPrefix)
	}
}

func TestNewMasterPersonaKeyringSecretStoreUsesNilGetenvFallback(t *testing.T) {
	capturedConfig := keyring.Config{}
	store, err := newMasterPersonaKeyringSecretStore(
		func(config keyring.Config) (keyring.Keyring, error) {
			capturedConfig = config
			return keyring.NewArrayKeyring(nil), nil
		},
		"darwin",
		nil,
	)
	if err != nil {
		t.Fatalf("expected nil getenv fallback to behave as empty env: %v", err)
	}
	if store == nil {
		t.Fatalf("expected store instance")
	}
	if len(capturedConfig.AllowedBackends) != 1 || capturedConfig.AllowedBackends[0] != keyring.KeychainBackend {
		t.Fatalf("expected default keychain backend when getenv is nil, got %#v", capturedConfig.AllowedBackends)
	}
}

func TestNewMasterPersonaKeyringSecretStoreSupportsExplicitBackendOverrides(t *testing.T) {
	cases := []struct {
		name          string
		override      string
		wantBackend   keyring.BackendType
		wantWinPrefix string
	}{
		{name: "keychain", override: string(keyring.KeychainBackend), wantBackend: keyring.KeychainBackend},
		{name: "wincred", override: string(keyring.WinCredBackend), wantBackend: keyring.WinCredBackend, wantWinPrefix: masterPersonaWinCredPrefix},
	}

	for _, testCase := range cases {
		t.Run(testCase.name, func(t *testing.T) {
			capturedConfig := keyring.Config{}
			store, err := newMasterPersonaKeyringSecretStore(
				func(config keyring.Config) (keyring.Keyring, error) {
					capturedConfig = config
					return keyring.NewArrayKeyring(nil), nil
				},
				"linux",
				envValues(map[string]string{masterPersonaSecretBackendEnv: testCase.override}),
			)
			if err != nil {
				t.Fatalf("expected explicit backend override %q to succeed: %v", testCase.override, err)
			}
			if store == nil {
				t.Fatalf("expected store instance")
			}
			if len(capturedConfig.AllowedBackends) != 1 || capturedConfig.AllowedBackends[0] != testCase.wantBackend {
				t.Fatalf("expected backend %q, got %#v", testCase.wantBackend, capturedConfig.AllowedBackends)
			}
			if capturedConfig.WinCredPrefix != testCase.wantWinPrefix {
				t.Fatalf("expected wincred prefix %q, got %q", testCase.wantWinPrefix, capturedConfig.WinCredPrefix)
			}
		})
	}
}

func TestNewMasterPersonaKeyringSecretStoreRejectsUnsupportedDefaultOS(t *testing.T) {
	store, err := newMasterPersonaKeyringSecretStore(
		func(_ keyring.Config) (keyring.Keyring, error) {
			return keyring.NewArrayKeyring(nil), nil
		},
		"linux",
		envValues(map[string]string{}),
	)
	if err == nil {
		t.Fatalf("expected unsupported default os to fail")
	}
	if !strings.Contains(err.Error(), "requires darwin or windows") {
		t.Fatalf("expected unsupported os error, got %v", err)
	}
	if store != nil {
		t.Fatalf("expected nil store when default os is unsupported")
	}
}

func TestNewMasterPersonaKeyringSecretStoreRejectsUnsupportedBackendOverride(t *testing.T) {
	store, err := newMasterPersonaKeyringSecretStore(
		func(_ keyring.Config) (keyring.Keyring, error) {
			return keyring.NewArrayKeyring(nil), nil
		},
		"darwin",
		envValues(map[string]string{masterPersonaSecretBackendEnv: "invalid"}),
	)
	if err == nil {
		t.Fatalf("expected unsupported backend override to fail")
	}
	if !strings.Contains(err.Error(), "unsupported master persona keyring backend override") {
		t.Fatalf("expected unsupported backend override error, got %v", err)
	}
	if store != nil {
		t.Fatalf("expected nil store when backend override is unsupported")
	}
}

func TestNewMasterPersonaKeyringSecretStoreRequiresFileBackendDirectory(t *testing.T) {
	store, err := newMasterPersonaKeyringSecretStore(
		func(_ keyring.Config) (keyring.Keyring, error) {
			return keyring.NewArrayKeyring(nil), nil
		},
		"darwin",
		envValues(map[string]string{
			masterPersonaSecretBackendEnv:      "file",
			masterPersonaSecretFilePasswordEnv: "password",
		}),
	)
	if err == nil {
		t.Fatalf("expected missing file directory to fail")
	}
	if !strings.Contains(err.Error(), masterPersonaSecretFileDirEnv) {
		t.Fatalf("expected missing file directory env error, got %v", err)
	}
	if store != nil {
		t.Fatalf("expected nil store when file directory env is missing")
	}
}

func TestNewMasterPersonaKeyringSecretStoreRequiresFileBackendPassword(t *testing.T) {
	store, err := newMasterPersonaKeyringSecretStore(
		func(_ keyring.Config) (keyring.Keyring, error) {
			return keyring.NewArrayKeyring(nil), nil
		},
		"darwin",
		envValues(map[string]string{
			masterPersonaSecretBackendEnv: "file",
			masterPersonaSecretFileDirEnv: filepath.Join(t.TempDir(), "keyring"),
		}),
	)
	if err == nil {
		t.Fatalf("expected missing file password to fail")
	}
	if !strings.Contains(err.Error(), masterPersonaSecretFilePasswordEnv) {
		t.Fatalf("expected missing file password env error, got %v", err)
	}
	if store != nil {
		t.Fatalf("expected nil store when file password env is missing")
	}
}

func TestNewMasterPersonaKeyringSecretStoreReturnsOpenError(t *testing.T) {
	openerErr := errors.New("open failed")
	store, err := newMasterPersonaKeyringSecretStore(
		func(_ keyring.Config) (keyring.Keyring, error) {
			return nil, openerErr
		},
		"darwin",
		envValues(map[string]string{}),
	)
	if err == nil {
		t.Fatalf("expected opener failure to propagate")
	}
	if !errors.Is(err, openerErr) {
		t.Fatalf("expected opener failure to be wrapped, got %v", err)
	}
	if store != nil {
		t.Fatalf("expected nil store on opener failure")
	}
}

func TestNewMasterPersonaKeyringSecretStoreSupportsFileBackendOverride(t *testing.T) {
	keyringDirectory := filepath.Join(t.TempDir(), "keyring")
	capturedConfig := keyring.Config{}
	store, err := newMasterPersonaKeyringSecretStore(
		func(config keyring.Config) (keyring.Keyring, error) {
			capturedConfig = config
			return keyring.NewArrayKeyring(nil), nil
		},
		"linux",
		func(key string) string {
			switch key {
			case masterPersonaSecretBackendEnv:
				return "file"
			case masterPersonaSecretFileDirEnv:
				return keyringDirectory
			case masterPersonaSecretFilePasswordEnv:
				return "test-password"
			default:
				return ""
			}
		},
	)
	if err != nil {
		t.Fatalf("expected file backend override to create keyring-backed store: %v", err)
	}
	if store == nil {
		t.Fatalf("expected keyring-backed store instance")
	}
	if len(capturedConfig.AllowedBackends) != 1 || capturedConfig.AllowedBackends[0] != keyring.FileBackend {
		t.Fatalf("expected file backend config, got %#v", capturedConfig.AllowedBackends)
	}
	if capturedConfig.FileDir != keyringDirectory {
		t.Fatalf("expected file backend directory %q, got %q", keyringDirectory, capturedConfig.FileDir)
	}
	if _, err := os.Stat(keyringDirectory); err != nil {
		t.Fatalf("expected file backend directory to be created: %v", err)
	}
}

func TestNewMasterPersonaKeyringSecretStoreProductionConstructorSupportsFileBackend(t *testing.T) {
	keyringDirectory := filepath.Join(t.TempDir(), "keyring")
	t.Setenv(masterPersonaSecretBackendEnv, "file")
	t.Setenv(masterPersonaSecretFileDirEnv, keyringDirectory)
	t.Setenv(masterPersonaSecretFilePasswordEnv, "test-password")

	store, err := NewMasterPersonaKeyringSecretStore()
	if err != nil {
		t.Fatalf("expected production constructor with file backend override to succeed: %v", err)
	}

	if saveErr := store.Save(context.Background(), "master-persona:provider", "api-key"); saveErr != nil {
		t.Fatalf("expected save through production constructor store: %v", saveErr)
	}
	loaded, err := store.Load(context.Background(), "master-persona:provider")
	if err != nil {
		t.Fatalf("expected load through production constructor store: %v", err)
	}
	if loaded != "api-key" {
		t.Fatalf("expected stored api key value, got %q", loaded)
	}
}

func TestInMemorySecretStoreSaveLoadDelete(t *testing.T) {
	store := NewInMemorySecretStore()
	if err := store.Save(context.Background(), "master-persona:test", "value"); err != nil {
		t.Fatalf("expected in-memory save to succeed: %v", err)
	}
	loaded, err := store.Load(context.Background(), "master-persona:test")
	if err != nil {
		t.Fatalf("expected in-memory load to succeed: %v", err)
	}
	if loaded != "value" {
		t.Fatalf("expected in-memory loaded value, got %q", loaded)
	}
	if err := store.Delete(context.Background(), "master-persona:test"); err != nil {
		t.Fatalf("expected in-memory delete to succeed: %v", err)
	}
}

type failingMasterPersonaKeyring struct {
	getErr    error
	setErr    error
	removeErr error
}

func (keyringBackend *failingMasterPersonaKeyring) Get(_ string) (keyring.Item, error) {
	return keyring.Item{}, keyringBackend.getErr
}

func (keyringBackend *failingMasterPersonaKeyring) GetMetadata(_ string) (keyring.Metadata, error) {
	return keyring.Metadata{}, nil
}

func (keyringBackend *failingMasterPersonaKeyring) Set(_ keyring.Item) error {
	return keyringBackend.setErr
}

func (keyringBackend *failingMasterPersonaKeyring) Remove(_ string) error {
	return keyringBackend.removeErr
}

func (keyringBackend *failingMasterPersonaKeyring) Keys() ([]string, error) {
	return nil, nil
}

func envValues(values map[string]string) func(string) string {
	return func(key string) string {
		return values[key]
	}
}
