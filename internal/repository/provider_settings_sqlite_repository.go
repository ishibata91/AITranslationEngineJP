package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/jmoiron/sqlx"
)

const (
	providerSettingsTimestampLayout = time.RFC3339Nano

	listProviderSettingsSQL = `
SELECT
  provider_id,
  endpoint,
  credential_reference_id,
  credential_state,
  validation_state,
  request_token,
  last_failure_kind,
  revision,
  updated_at
FROM PROVIDER_SETTINGS
ORDER BY provider_id ASC;`
	getProviderSettingsByProviderIDSQL = `
SELECT
  provider_id,
  endpoint,
  credential_reference_id,
  credential_state,
  validation_state,
  request_token,
  last_failure_kind,
  revision,
  updated_at
FROM PROVIDER_SETTINGS
WHERE provider_id = ?
LIMIT 1;`
	upsertProviderSettingsSQL = `
INSERT INTO PROVIDER_SETTINGS (
  provider_id,
  endpoint,
  credential_reference_id,
  credential_state,
  validation_state,
  request_token,
  last_failure_kind,
  revision,
  updated_at
) VALUES (
  :provider_id,
  :endpoint,
  :credential_reference_id,
  :credential_state,
  :validation_state,
  :request_token,
  :last_failure_kind,
  :revision,
  :updated_at
)
ON CONFLICT(provider_id) DO UPDATE
SET endpoint = excluded.endpoint,
    credential_reference_id = excluded.credential_reference_id,
    credential_state = excluded.credential_state,
    validation_state = excluded.validation_state,
    request_token = excluded.request_token,
    last_failure_kind = excluded.last_failure_kind,
    revision = excluded.revision,
    updated_at = excluded.updated_at;`
	updateProviderSettingsValidationByRequestTokenSQL = `
UPDATE PROVIDER_SETTINGS
SET validation_state = ?,
    last_failure_kind = ?,
    updated_at = ?
WHERE provider_id = ?
  AND request_token = ?;`
)

// SQLiteProviderSettingsRepository persists provider settings rows to SQLite.
type SQLiteProviderSettingsRepository struct {
	db *sqlx.DB
}

type sqliteProviderSettingsRow struct {
	ProviderID            string         `db:"provider_id"`
	Endpoint              sql.NullString `db:"endpoint"`
	CredentialReferenceID sql.NullString `db:"credential_reference_id"`
	CredentialState       string         `db:"credential_state"`
	ValidationState       string         `db:"validation_state"`
	RequestToken          sql.NullString `db:"request_token"`
	LastFailureKind       sql.NullString `db:"last_failure_kind"`
	Revision              int64          `db:"revision"`
	UpdatedAt             string         `db:"updated_at"`
}

type sqliteProviderSettingsUpsertRow struct {
	ProviderID            string  `db:"provider_id"`
	Endpoint              *string `db:"endpoint"`
	CredentialReferenceID *string `db:"credential_reference_id"`
	CredentialState       string  `db:"credential_state"`
	ValidationState       string  `db:"validation_state"`
	RequestToken          *string `db:"request_token"`
	LastFailureKind       *string `db:"last_failure_kind"`
	Revision              int64   `db:"revision"`
	UpdatedAt             string  `db:"updated_at"`
}

// NewSQLiteProviderSettingsRepository creates a SQLite-backed provider settings repository.
func NewSQLiteProviderSettingsRepository(db *sqlx.DB) *SQLiteProviderSettingsRepository {
	return &SQLiteProviderSettingsRepository{db: db}
}

// List returns every persisted provider settings row.
func (repository *SQLiteProviderSettingsRepository) List(ctx context.Context) ([]ProviderSettingsRecord, error) {
	rows := []sqliteProviderSettingsRow{}
	if err := repository.db.SelectContext(ctx, &rows, listProviderSettingsSQL); err != nil {
		return nil, fmt.Errorf("list provider settings: %w", err)
	}
	result := make([]ProviderSettingsRecord, 0, len(rows))
	for _, row := range rows {
		record, err := row.toRecord()
		if err != nil {
			return nil, err
		}
		result = append(result, record)
	}
	return result, nil
}

// GetByProviderID returns one persisted provider settings row.
func (repository *SQLiteProviderSettingsRepository) GetByProviderID(
	ctx context.Context,
	providerID string,
) (ProviderSettingsRecord, error) {
	row := sqliteProviderSettingsRow{}
	if err := sqlx.GetContext(
		ctx,
		extractTx(ctx, repository.db),
		&row,
		getProviderSettingsByProviderIDSQL,
		strings.TrimSpace(providerID),
	); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return ProviderSettingsRecord{}, fmt.Errorf("%w: provider_id=%s", ErrProviderSettingsNotFound, providerID)
		}
		return ProviderSettingsRecord{}, fmt.Errorf("get provider settings by provider id: %w", err)
	}
	return row.toRecord()
}

// Upsert stores one provider settings row without secrets.
func (repository *SQLiteProviderSettingsRepository) Upsert(
	ctx context.Context,
	draft ProviderSettingsUpsertDraft,
) (ProviderSettingsRecord, error) {
	upsertRow := sqliteProviderSettingsUpsertRow{
		ProviderID:            strings.TrimSpace(draft.ProviderID),
		Endpoint:              normalizeOptionalProviderSettingsString(draft.Endpoint),
		CredentialReferenceID: normalizeOptionalProviderSettingsString(draft.CredentialReferenceID),
		CredentialState:       strings.TrimSpace(draft.CredentialState),
		ValidationState:       strings.TrimSpace(draft.ValidationState),
		RequestToken:          normalizeOptionalProviderSettingsString(draft.RequestToken),
		LastFailureKind:       normalizeOptionalProviderSettingsString(draft.LastFailureKind),
		Revision:              draft.Revision,
		UpdatedAt:             draft.UpdatedAt.UTC().Format(providerSettingsTimestampLayout),
	}
	if _, err := sqlx.NamedExecContext(ctx, extractTx(ctx, repository.db), upsertProviderSettingsSQL, upsertRow); err != nil {
		return ProviderSettingsRecord{}, fmt.Errorf("upsert provider settings: %w", err)
	}
	return repository.GetByProviderID(ctx, upsertRow.ProviderID)
}

// UpdateValidationByRequestToken updates validation fields only when the request token still matches.
func (repository *SQLiteProviderSettingsRepository) UpdateValidationByRequestToken(
	ctx context.Context,
	providerID string,
	requestToken string,
	validationState string,
	lastFailureKind *string,
) (ProviderSettingsRecord, bool, error) {
	updatedAt := time.Now().UTC().Format(providerSettingsTimestampLayout)
	execResult, err := extractTx(ctx, repository.db).ExecContext(
		ctx,
		updateProviderSettingsValidationByRequestTokenSQL,
		strings.TrimSpace(validationState),
		nullStringValue(normalizeOptionalProviderSettingsString(lastFailureKind)),
		updatedAt,
		strings.TrimSpace(providerID),
		strings.TrimSpace(requestToken),
	)
	if err != nil {
		return ProviderSettingsRecord{}, false, fmt.Errorf("update provider settings validation: %w", err)
	}
	affected, err := execResult.RowsAffected()
	if err != nil {
		return ProviderSettingsRecord{}, false, fmt.Errorf("read provider settings validation rows affected: %w", err)
	}
	if affected == 0 {
		return ProviderSettingsRecord{}, false, nil
	}
	record, err := repository.GetByProviderID(ctx, providerID)
	if err != nil {
		return ProviderSettingsRecord{}, false, err
	}
	return record, true, nil
}

func (row sqliteProviderSettingsRow) toRecord() (ProviderSettingsRecord, error) {
	updatedAt, err := time.Parse(providerSettingsTimestampLayout, strings.TrimSpace(row.UpdatedAt))
	if err != nil {
		return ProviderSettingsRecord{}, fmt.Errorf("parse provider settings updated_at: %w", err)
	}
	return ProviderSettingsRecord{
		ProviderID:            strings.TrimSpace(row.ProviderID),
		Endpoint:              nullStringPointer(row.Endpoint),
		CredentialReferenceID: nullStringPointer(row.CredentialReferenceID),
		CredentialState:       strings.TrimSpace(row.CredentialState),
		ValidationState:       strings.TrimSpace(row.ValidationState),
		RequestToken:          nullStringPointer(row.RequestToken),
		LastFailureKind:       nullStringPointer(row.LastFailureKind),
		Revision:              row.Revision,
		UpdatedAt:             updatedAt.UTC(),
	}, nil
}

func nullStringPointer(value sql.NullString) *string {
	if !value.Valid {
		return nil
	}
	trimmed := strings.TrimSpace(value.String)
	if trimmed == "" {
		return nil
	}
	cloned := trimmed
	return &cloned
}

func nullStringValue(value *string) any {
	if value == nil {
		return nil
	}
	return *value
}
