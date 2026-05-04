package repository

import (
	"context"
	"errors"
	"strings"
	"time"
)

var (
	// ErrProviderSettingsNotFound reports that one provider settings row does not exist.
	ErrProviderSettingsNotFound = errors.New("provider settings not found")
)

// ProviderSettingsRecord stores one persisted provider settings row without secrets.
type ProviderSettingsRecord struct {
	ProviderID            string
	Endpoint              *string
	CredentialReferenceID *string
	CredentialState       string
	ValidationState       string
	RequestToken          *string
	LastFailureKind       *string
	Revision              int64
	UpdatedAt             time.Time
}

// ProviderSettingsUpsertDraft stores the persistence payload for one provider settings row.
type ProviderSettingsUpsertDraft struct {
	ProviderID            string
	Endpoint              *string
	CredentialReferenceID *string
	CredentialState       string
	ValidationState       string
	RequestToken          *string
	LastFailureKind       *string
	Revision              int64
	UpdatedAt             time.Time
}

// ProviderSettingsRepository defines provider settings persistence without secret values.
type ProviderSettingsRepository interface {
	List(ctx context.Context) ([]ProviderSettingsRecord, error)
	GetByProviderID(ctx context.Context, providerID string) (ProviderSettingsRecord, error)
	Upsert(ctx context.Context, draft ProviderSettingsUpsertDraft) (ProviderSettingsRecord, error)
	UpdateValidationByRequestToken(
		ctx context.Context,
		providerID string,
		requestToken string,
		validationState string,
		lastFailureKind *string,
	) (ProviderSettingsRecord, bool, error)
}

// ProviderSettingsSecretStore defines provider settings secret access.
type ProviderSettingsSecretStore interface {
	Load(ctx context.Context, key string) (string, error)
	Save(ctx context.Context, key string, value string) error
	Delete(ctx context.Context, key string) error
}

func normalizeOptionalProviderSettingsString(value *string) *string {
	if value == nil {
		return nil
	}
	trimmed := strings.TrimSpace(*value)
	if trimmed == "" {
		return nil
	}
	cloned := trimmed
	return &cloned
}
