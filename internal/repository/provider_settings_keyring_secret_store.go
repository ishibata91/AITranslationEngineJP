package repository

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/99designs/keyring"
)

const (
	providerSettingsKeyringServiceName = "AITranslationEngineJP.ProviderSettings"
	providerSettingsWinCredPrefix      = "AITEJP:"

	providerSettingsSecretBackendEnv      = "AITRANSLATIONENGINEJP_PROVIDER_SETTINGS_SECRET_BACKEND"
	providerSettingsSecretFileDirEnv      = "AITRANSLATIONENGINEJP_PROVIDER_SETTINGS_SECRET_FILE_DIR"
	providerSettingsSecretFilePasswordEnv = "AITRANSLATIONENGINEJP_PROVIDER_SETTINGS_SECRET_FILE_PASSWORD"
)

// ProviderSettingsKeyringSecretStore stores provider settings API keys in OS keyring backends.
type ProviderSettingsKeyringSecretStore struct {
	keyring keyring.Keyring
}

// NewProviderSettingsKeyringSecretStore creates the production keyring-backed provider settings secret store.
func NewProviderSettingsKeyringSecretStore() (*ProviderSettingsKeyringSecretStore, error) {
	return newProviderSettingsKeyringSecretStore(keyring.Open, runtime.GOOS, os.Getenv)
}

// NewProviderSettingsKeyringSecretStoreWithKeyring creates a keyring-backed provider settings secret store from an injected backend.
func NewProviderSettingsKeyringSecretStoreWithKeyring(
	backend keyring.Keyring,
) (*ProviderSettingsKeyringSecretStore, error) {
	if backend == nil {
		return nil, fmt.Errorf("provider settings keyring backend is required")
	}
	return &ProviderSettingsKeyringSecretStore{keyring: backend}, nil
}

func newProviderSettingsKeyringSecretStore(
	openKeyring keyringOpenFunc,
	goos string,
	getenv func(string) string,
) (*ProviderSettingsKeyringSecretStore, error) {
	if openKeyring == nil {
		return nil, fmt.Errorf("provider settings keyring opener is required")
	}
	if getenv == nil {
		getenv = func(string) string { return "" }
	}

	config, err := newProviderSettingsKeyringConfig(goos, getenv)
	if err != nil {
		return nil, err
	}
	if ensureErr := ensureProviderSettingsFileBackendDirectory(config); ensureErr != nil {
		return nil, ensureErr
	}
	backend, err := openKeyring(config)
	if err != nil {
		return nil, fmt.Errorf("open provider settings keyring backend: %w", err)
	}
	return &ProviderSettingsKeyringSecretStore{keyring: backend}, nil
}

func newProviderSettingsKeyringConfig(goos string, getenv func(string) string) (keyring.Config, error) {
	config := keyring.Config{ServiceName: providerSettingsKeyringServiceName}
	requestedBackend := strings.ToLower(strings.TrimSpace(getenv(providerSettingsSecretBackendEnv)))

	switch requestedBackend {
	case "", "default":
		switch strings.ToLower(strings.TrimSpace(goos)) {
		case "darwin":
			config.AllowedBackends = []keyring.BackendType{keyring.KeychainBackend}
			config.KeychainTrustApplication = true
		case "windows":
			config.AllowedBackends = []keyring.BackendType{keyring.WinCredBackend}
			config.WinCredPrefix = providerSettingsWinCredPrefix
		default:
			return keyring.Config{}, fmt.Errorf(
				"provider settings keyring backend requires darwin or windows (goos=%s); set %s=file for test-safe backend",
				strings.TrimSpace(goos),
				providerSettingsSecretBackendEnv,
			)
		}
	case string(keyring.KeychainBackend):
		config.AllowedBackends = []keyring.BackendType{keyring.KeychainBackend}
		config.KeychainTrustApplication = true
	case string(keyring.WinCredBackend):
		config.AllowedBackends = []keyring.BackendType{keyring.WinCredBackend}
		config.WinCredPrefix = providerSettingsWinCredPrefix
	case string(keyring.FileBackend):
		fileDirectory := strings.TrimSpace(getenv(providerSettingsSecretFileDirEnv))
		if fileDirectory == "" {
			return keyring.Config{}, fmt.Errorf(
				"provider settings keyring file backend requires %s",
				providerSettingsSecretFileDirEnv,
			)
		}
		filePassword := strings.TrimSpace(getenv(providerSettingsSecretFilePasswordEnv))
		if filePassword == "" {
			return keyring.Config{}, fmt.Errorf(
				"provider settings keyring file backend requires %s",
				providerSettingsSecretFilePasswordEnv,
			)
		}
		config.AllowedBackends = []keyring.BackendType{keyring.FileBackend}
		config.FileDir = fileDirectory
		config.FilePasswordFunc = keyring.FixedStringPrompt(filePassword)
	default:
		return keyring.Config{}, fmt.Errorf("unsupported provider settings keyring backend override")
	}

	return config, nil
}

func ensureProviderSettingsFileBackendDirectory(config keyring.Config) error {
	if len(config.AllowedBackends) != 1 || config.AllowedBackends[0] != keyring.FileBackend {
		return nil
	}
	if err := os.MkdirAll(filepath.Clean(config.FileDir), 0o700); err != nil {
		return fmt.Errorf("create provider settings keyring file directory: %w", err)
	}
	return nil
}

// Load returns one provider settings secret by key.
func (store *ProviderSettingsKeyringSecretStore) Load(_ context.Context, key string) (string, error) {
	trimmedKey := strings.TrimSpace(key)
	if trimmedKey == "" {
		return "", nil
	}
	item, err := store.keyring.Get(trimmedKey)
	if errors.Is(err, keyring.ErrKeyNotFound) {
		return "", nil
	}
	if err != nil {
		return "", fmt.Errorf("load provider settings secret from keyring: %w", err)
	}
	return string(item.Data), nil
}

// Save stores one provider settings secret by key.
func (store *ProviderSettingsKeyringSecretStore) Save(_ context.Context, key string, value string) error {
	trimmedKey := strings.TrimSpace(key)
	if trimmedKey == "" {
		return fmt.Errorf("provider settings secret key is required")
	}
	if err := store.keyring.Set(keyring.Item{Key: trimmedKey, Data: []byte(value)}); err != nil {
		return fmt.Errorf("save provider settings secret to keyring: %w", err)
	}
	return nil
}

// Delete removes one provider settings secret by key.
func (store *ProviderSettingsKeyringSecretStore) Delete(_ context.Context, key string) error {
	trimmedKey := strings.TrimSpace(key)
	if trimmedKey == "" {
		return nil
	}
	if err := store.keyring.Remove(trimmedKey); err != nil && !errors.Is(err, keyring.ErrKeyNotFound) {
		return fmt.Errorf("delete provider settings secret from keyring: %w", err)
	}
	return nil
}
