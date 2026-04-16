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
	masterPersonaKeyringServiceName = "AITranslationEngineJP.MasterPersona"
	masterPersonaWinCredPrefix      = "AITEJP:"

	masterPersonaSecretBackendEnv      = "AITRANSLATIONENGINEJP_MASTER_PERSONA_SECRET_BACKEND"
	masterPersonaSecretFileDirEnv      = "AITRANSLATIONENGINEJP_MASTER_PERSONA_SECRET_FILE_DIR"
	masterPersonaSecretFilePasswordEnv = "AITRANSLATIONENGINEJP_MASTER_PERSONA_SECRET_FILE_PASSWORD"
)

type keyringOpenFunc func(config keyring.Config) (keyring.Keyring, error)

// MasterPersonaKeyringSecretStore stores master persona API keys in OS keyring backends.
type MasterPersonaKeyringSecretStore struct {
	keyring keyring.Keyring
}

// NewMasterPersonaKeyringSecretStore creates the production keyring-backed secret store.
func NewMasterPersonaKeyringSecretStore() (*MasterPersonaKeyringSecretStore, error) {
	return newMasterPersonaKeyringSecretStore(keyring.Open, runtime.GOOS, os.Getenv)
}

// NewMasterPersonaKeyringSecretStoreWithKeyring creates a keyring-backed secret store from an injected backend.
func NewMasterPersonaKeyringSecretStoreWithKeyring(
	backend keyring.Keyring,
) (*MasterPersonaKeyringSecretStore, error) {
	if backend == nil {
		return nil, fmt.Errorf("master persona keyring backend is required")
	}
	return &MasterPersonaKeyringSecretStore{keyring: backend}, nil
}

func newMasterPersonaKeyringSecretStore(
	openKeyring keyringOpenFunc,
	goos string,
	getenv func(string) string,
) (*MasterPersonaKeyringSecretStore, error) {
	if openKeyring == nil {
		return nil, fmt.Errorf("master persona keyring opener is required")
	}
	if getenv == nil {
		getenv = func(string) string { return "" }
	}

	config, err := newMasterPersonaKeyringConfig(goos, getenv)
	if err != nil {
		return nil, err
	}
	if ensureErr := ensureMasterPersonaFileBackendDirectory(config); ensureErr != nil {
		return nil, ensureErr
	}
	backend, err := openKeyring(config)
	if err != nil {
		return nil, fmt.Errorf("open master persona keyring backend: %w", err)
	}
	return &MasterPersonaKeyringSecretStore{keyring: backend}, nil
}

func newMasterPersonaKeyringConfig(goos string, getenv func(string) string) (keyring.Config, error) {
	config := keyring.Config{ServiceName: masterPersonaKeyringServiceName}
	requestedBackend := strings.ToLower(strings.TrimSpace(getenv(masterPersonaSecretBackendEnv)))

	switch requestedBackend {
	case "", "default":
		switch strings.ToLower(strings.TrimSpace(goos)) {
		case "darwin":
			config.AllowedBackends = []keyring.BackendType{keyring.KeychainBackend}
		case "windows":
			config.AllowedBackends = []keyring.BackendType{keyring.WinCredBackend}
			config.WinCredPrefix = masterPersonaWinCredPrefix
		default:
			return keyring.Config{}, fmt.Errorf(
				"master persona keyring backend requires darwin or windows (goos=%s); set %s=file for test-safe backend",
				strings.TrimSpace(goos),
				masterPersonaSecretBackendEnv,
			)
		}
	case string(keyring.KeychainBackend):
		config.AllowedBackends = []keyring.BackendType{keyring.KeychainBackend}
	case string(keyring.WinCredBackend):
		config.AllowedBackends = []keyring.BackendType{keyring.WinCredBackend}
		config.WinCredPrefix = masterPersonaWinCredPrefix
	case string(keyring.FileBackend):
		fileDirectory := strings.TrimSpace(getenv(masterPersonaSecretFileDirEnv))
		if fileDirectory == "" {
			return keyring.Config{}, fmt.Errorf(
				"master persona keyring file backend requires %s",
				masterPersonaSecretFileDirEnv,
			)
		}
		filePassword := strings.TrimSpace(getenv(masterPersonaSecretFilePasswordEnv))
		if filePassword == "" {
			return keyring.Config{}, fmt.Errorf(
				"master persona keyring file backend requires %s",
				masterPersonaSecretFilePasswordEnv,
			)
		}
		config.AllowedBackends = []keyring.BackendType{keyring.FileBackend}
		config.FileDir = fileDirectory
		config.FilePasswordFunc = keyring.FixedStringPrompt(filePassword)
	default:
		return keyring.Config{}, fmt.Errorf(
			"unsupported master persona keyring backend override: %s",
			requestedBackend,
		)
	}

	return config, nil
}

func ensureMasterPersonaFileBackendDirectory(config keyring.Config) error {
	if len(config.AllowedBackends) != 1 || config.AllowedBackends[0] != keyring.FileBackend {
		return nil
	}
	if err := os.MkdirAll(filepath.Clean(config.FileDir), 0o700); err != nil {
		return fmt.Errorf("create master persona keyring file directory: %w", err)
	}
	return nil
}

// Load returns one secret value by key.
func (store *MasterPersonaKeyringSecretStore) Load(_ context.Context, key string) (string, error) {
	trimmedKey := strings.TrimSpace(key)
	if trimmedKey == "" {
		return "", nil
	}
	item, err := store.keyring.Get(trimmedKey)
	if errors.Is(err, keyring.ErrKeyNotFound) {
		return "", nil
	}
	if err != nil {
		return "", fmt.Errorf("load master persona secret from keyring: %w", err)
	}
	return string(item.Data), nil
}

// Save stores one secret value by key.
func (store *MasterPersonaKeyringSecretStore) Save(_ context.Context, key string, value string) error {
	trimmedKey := strings.TrimSpace(key)
	if trimmedKey == "" {
		return fmt.Errorf("master persona secret key is required")
	}
	if err := store.keyring.Set(keyring.Item{Key: trimmedKey, Data: []byte(value)}); err != nil {
		return fmt.Errorf("save master persona secret to keyring: %w", err)
	}
	return nil
}

// Delete removes one secret value by key.
func (store *MasterPersonaKeyringSecretStore) Delete(_ context.Context, key string) error {
	trimmedKey := strings.TrimSpace(key)
	if trimmedKey == "" {
		return nil
	}
	if err := store.keyring.Remove(trimmedKey); err != nil && !errors.Is(err, keyring.ErrKeyNotFound) {
		return fmt.Errorf("delete master persona secret from keyring: %w", err)
	}
	return nil
}
