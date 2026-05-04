package repository

import (
	"context"
	"fmt"
	"strings"
	"sync"
)

// CachedProviderSettingsSecretStore keeps provider settings secrets in process memory after explicit reads.
type CachedProviderSettingsSecretStore struct {
	backend ProviderSettingsSecretStore
	mu      sync.RWMutex
	values  map[string]string
}

// NewCachedProviderSettingsSecretStore wraps a provider settings secret store with a process-local cache.
func NewCachedProviderSettingsSecretStore(backend ProviderSettingsSecretStore) (*CachedProviderSettingsSecretStore, error) {
	if backend == nil {
		return nil, fmt.Errorf("provider settings cached secret backend is required")
	}
	return &CachedProviderSettingsSecretStore{
		backend: backend,
		values:  make(map[string]string),
	}, nil
}

// Load returns a cached secret when one was already resolved in this process.
func (store *CachedProviderSettingsSecretStore) Load(ctx context.Context, key string) (string, error) {
	trimmedKey := strings.TrimSpace(key)
	if trimmedKey == "" {
		return "", nil
	}
	store.mu.RLock()
	cached, ok := store.values[trimmedKey]
	store.mu.RUnlock()
	if ok {
		return cached, nil
	}
	loaded, err := store.backend.Load(ctx, trimmedKey)
	if err != nil {
		return "", fmt.Errorf("load provider settings cached secret backend: %w", err)
	}
	store.mu.Lock()
	store.values[trimmedKey] = strings.TrimSpace(loaded)
	store.mu.Unlock()
	return strings.TrimSpace(loaded), nil
}

// Save stores one secret and refreshes the process-local cache.
func (store *CachedProviderSettingsSecretStore) Save(ctx context.Context, key string, value string) error {
	trimmedKey := strings.TrimSpace(key)
	if err := store.backend.Save(ctx, trimmedKey, value); err != nil {
		return fmt.Errorf("save provider settings cached secret backend: %w", err)
	}
	store.mu.Lock()
	store.values[trimmedKey] = strings.TrimSpace(value)
	store.mu.Unlock()
	return nil
}

// Delete removes one secret and clears the process-local cache entry.
func (store *CachedProviderSettingsSecretStore) Delete(ctx context.Context, key string) error {
	trimmedKey := strings.TrimSpace(key)
	if err := store.backend.Delete(ctx, trimmedKey); err != nil {
		return fmt.Errorf("delete provider settings cached secret backend: %w", err)
	}
	store.mu.Lock()
	delete(store.values, trimmedKey)
	store.mu.Unlock()
	return nil
}
