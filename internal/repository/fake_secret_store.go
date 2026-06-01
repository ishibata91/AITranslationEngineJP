package repository

import "context"

// FakeSecretStore provides a deterministic secret seam for tests and local fake-mode wiring.
// Load always returns a fixed non-empty value so that credential guards pass in fake environments.
// Save and Delete are no-ops; they neither store nor remove any value.
type FakeSecretStore struct{}

// NewFakeSecretStore creates a FakeSecretStore that satisfies SecretStore and ProviderSettingsSecretStore.
func NewFakeSecretStore() *FakeSecretStore {
	return &FakeSecretStore{}
}

// Load returns a fixed non-empty secret value regardless of key.
func (store *FakeSecretStore) Load(_ context.Context, _ string) (string, error) {
	return "fake-secret", nil
}

// Save is a no-op. FakeSecretStore does not persist any value.
func (store *FakeSecretStore) Save(_ context.Context, _ string, _ string) error {
	return nil
}

// Delete is a no-op. FakeSecretStore does not remove any value.
func (store *FakeSecretStore) Delete(_ context.Context, _ string) error {
	return nil
}
