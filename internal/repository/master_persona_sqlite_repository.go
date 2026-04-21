package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"time"

	sqliteinfra "aitranslationenginejp/internal/infra/sqlite/dbinit"

	"github.com/jmoiron/sqlx"
)

const (
	masterPersonaTimestampLayout      = time.RFC3339Nano
	masterPersonaDefaultRunState      = "入力待ち"
	masterPersonaIdentityKeyErrFormat = "%w: identity_key=%s"

	// canonical list / count / plugin-groups paths
	countCanonicalPersonasSQL = `
SELECT COUNT(*)
FROM NPC_PROFILE np
WHERE (
  ? = ''
  OR lower(np.display_name || ' ' || np.form_id || ' ' || np.editor_id || ' ' || np.target_plugin_name) LIKE ?
)
AND (
  ? = ''
  OR np.target_plugin_name = ?
);`
	listCanonicalPersonasSQL = `
SELECT
  np.target_plugin_name,
  np.form_id,
  np.record_type,
  np.editor_id,
  np.display_name,
  np.updated_at,
  COALESCE(p.persona_description, '') AS persona_description,
  COALESCE(p.personality_summary, '') AS persona_summary,
  COALESCE(p.speech_style, '') AS speech_style,
  nr.race,
  nr.sex,
  COALESCE(nr.npc_class, '') AS npc_class,
  COALESCE(nr.voice_type, '') AS voice_type
FROM NPC_PROFILE np
LEFT JOIN PERSONA p ON p.npc_profile_id = np.id
LEFT JOIN NPC_RECORD nr ON nr.translation_record_id = (
  SELECT MAX(nr2.translation_record_id) FROM NPC_RECORD nr2 WHERE nr2.npc_profile_id = np.id
)
WHERE (
  ? = ''
  OR lower(np.display_name || ' ' || np.form_id || ' ' || np.editor_id || ' ' || np.target_plugin_name) LIKE ?
)
AND (
  ? = ''
  OR np.target_plugin_name = ?
)
ORDER BY np.updated_at DESC, np.target_plugin_name ASC, np.form_id ASC
LIMIT ?
OFFSET ?;`
	listCanonicalPluginGroupsSQL = `
SELECT np.target_plugin_name AS target_plugin, COUNT(*) AS count
FROM NPC_PROFILE np
WHERE (
  ? = ''
  OR lower(np.display_name || ' ' || np.form_id || ' ' || np.editor_id || ' ' || np.target_plugin_name) LIKE ?
)
GROUP BY np.target_plugin_name
ORDER BY np.target_plugin_name ASC;`

	// canonical seed check
	countCanonicalNPCProfileRowsSQL = `SELECT COUNT(*) FROM NPC_PROFILE;`

	// canonical update path
	updateCanonicalNPCProfileSQL = `
UPDATE NPC_PROFILE
SET target_plugin_name = ?,
    form_id = ?,
    record_type = ?,
    editor_id = ?,
    display_name = ?,
    updated_at = ?
WHERE target_plugin_name = ? AND form_id = ? AND record_type = ?;`
	updateCanonicalPersonaSQL = `
UPDATE PERSONA
SET persona_description = ?,
    personality_summary = ?,
    speech_style = ?,
    updated_at = ?
WHERE npc_profile_id = (
  SELECT id FROM NPC_PROFILE
  WHERE target_plugin_name = ? AND form_id = ? AND record_type = ?
);`

	// canonical delete path
	deleteCanonicalPersonaSQL = `
DELETE FROM PERSONA
WHERE npc_profile_id = (
  SELECT id FROM NPC_PROFILE
  WHERE target_plugin_name = ? AND form_id = ? AND record_type = ?
);`
	deleteCanonicalNPCProfileSQL = `
DELETE FROM NPC_PROFILE
WHERE target_plugin_name = ? AND form_id = ? AND record_type = ?;`

	loadMasterPersonaAISettingsSQL = `
SELECT provider, model
FROM PERSONA_GENERATION_SETTINGS
WHERE id = 1
LIMIT 1;`
	upsertMasterPersonaAISettingsSQL = `
INSERT INTO PERSONA_GENERATION_SETTINGS (id, provider, model)
VALUES (1, ?, ?)
ON CONFLICT(id) DO UPDATE
SET provider = excluded.provider,
    model = excluded.model;`

	// canonical generation write path
	insertCanonicalNPCProfileIfAbsentSQL = `
INSERT OR IGNORE INTO NPC_PROFILE (
  target_plugin_name, form_id, record_type, editor_id, display_name, created_at, updated_at
) VALUES (?, ?, ?, ?, ?, ?, ?);`
	insertCanonicalPersonaSQL = `
INSERT INTO PERSONA (
  npc_profile_id, persona_description, personality_summary, speech_style, created_at, updated_at
) VALUES (?, ?, ?, ?, ?, ?);`
	getCanonicalMasterPersonaByIdentitySQL = `
SELECT
  np.target_plugin_name,
  np.form_id,
  np.record_type,
  np.editor_id,
  np.display_name,
  np.updated_at,
  COALESCE(p.persona_description, '') AS persona_description,
  COALESCE(p.personality_summary, '') AS persona_summary,
  COALESCE(p.speech_style, '') AS speech_style,
  nr.race,
  nr.sex,
  COALESCE(nr.npc_class, '') AS npc_class,
  COALESCE(nr.voice_type, '') AS voice_type
FROM NPC_PROFILE np
LEFT JOIN PERSONA p ON p.npc_profile_id = np.id
LEFT JOIN NPC_RECORD nr ON nr.translation_record_id = (
  SELECT MAX(nr2.translation_record_id) FROM NPC_RECORD nr2 WHERE nr2.npc_profile_id = np.id
)
WHERE np.target_plugin_name = ? AND np.form_id = ? AND np.record_type = ?
LIMIT 1;`
)

// SQLiteMasterPersonaRepositories bundles SQLite-backed concrete repositories by responsibility.
type SQLiteMasterPersonaRepositories struct {
	database             *sqlx.DB
	EntryRepository      *SQLiteMasterPersonaEntryRepository
	AISettingsRepository *SQLiteMasterPersonaAISettingsRepository
	RunStatusRepository  *SQLiteMasterPersonaRunStatusRepository
}

// SQLiteMasterPersonaEntryRepository persists master persona entries.
type SQLiteMasterPersonaEntryRepository struct {
	database *sqlx.DB
}

// SQLiteMasterPersonaAISettingsRepository persists page-local AI settings.
type SQLiteMasterPersonaAISettingsRepository struct {
	database *sqlx.DB
}

// SQLiteMasterPersonaRunStatusRepository persists generation run status.
type SQLiteMasterPersonaRunStatusRepository struct {
	database *sqlx.DB
}

type sqliteMasterPersonaPluginGroupRow struct {
	TargetPlugin string `db:"target_plugin"`
	Count        int    `db:"count"`
}

type sqliteMasterPersonaAISettingsRow struct {
	Provider string `db:"provider"`
	Model    string `db:"model"`
}

type canonicalNPCProfilePersonaRow struct {
	TargetPluginName   string  `db:"target_plugin_name"`
	FormID             string  `db:"form_id"`
	RecordType         string  `db:"record_type"`
	EditorID           string  `db:"editor_id"`
	DisplayName        string  `db:"display_name"`
	UpdatedAt          string  `db:"updated_at"`
	PersonaDescription string  `db:"persona_description"`
	PersonaSummary     string  `db:"persona_summary"`
	SpeechStyle        string  `db:"speech_style"`
	Race               *string `db:"race"`
	Sex                *string `db:"sex"`
	NPCClass           string  `db:"npc_class"`
	VoiceType          string  `db:"voice_type"`
}

// NewSQLiteMasterPersonaRepositories opens SQLite-backed master persona repositories.
func NewSQLiteMasterPersonaRepositories(
	ctx context.Context,
	databasePath string,
	seed []MasterPersonaEntry,
) (*SQLiteMasterPersonaRepositories, error) {
	database, err := sqliteinfra.OpenMasterDictionaryDatabase(ctx, databasePath, nil)
	if err != nil {
		return nil, fmt.Errorf("open sqlite database for master persona repositories: %w", err)
	}

	entryRepository := &SQLiteMasterPersonaEntryRepository{database: database}
	settingsRepository := &SQLiteMasterPersonaAISettingsRepository{database: database}
	runStatusRepository := &SQLiteMasterPersonaRunStatusRepository{database: database}

	if err := entryRepository.seedIfEmpty(ctx, seed); err != nil {
		if closeErr := database.Close(); closeErr != nil {
			return nil, fmt.Errorf("seed sqlite master persona entries: %w", errors.Join(err, closeErr))
		}
		return nil, err
	}
	return &SQLiteMasterPersonaRepositories{
		database:             database,
		EntryRepository:      entryRepository,
		AISettingsRepository: settingsRepository,
		RunStatusRepository:  runStatusRepository,
	}, nil
}

// Close releases the shared SQLite database handle.
func (repositories *SQLiteMasterPersonaRepositories) Close() error {
	if err := repositories.database.Close(); err != nil {
		return fmt.Errorf("close sqlite master persona repositories database: %w", err)
	}
	return nil
}

// Database returns the underlying *sqlx.DB so infrastructure components (e.g. transactor) can share the connection.
func (repositories *SQLiteMasterPersonaRepositories) Database() *sqlx.DB {
	return repositories.database
}

// List returns filtered and paginated master persona entries from the canonical NPC_PROFILE + PERSONA schema.
func (repository *SQLiteMasterPersonaEntryRepository) List(
	ctx context.Context,
	query MasterPersonaListQuery,
) (MasterPersonaListResult, error) {
	keyword := strings.ToLower(strings.TrimSpace(query.Keyword))
	searchPattern := "%" + keyword + "%"
	pluginFilter := strings.TrimSpace(query.PluginFilter)

	var totalCount int
	if err := repository.database.GetContext(
		ctx,
		&totalCount,
		countCanonicalPersonasSQL,
		keyword,
		searchPattern,
		pluginFilter,
		pluginFilter,
	); err != nil {
		return MasterPersonaListResult{}, fmt.Errorf("count canonical personas: %w", err)
	}

	page, pageSize := normalizeMasterPersonaPagination(query.Page, query.PageSize, totalCount)
	rows := []canonicalNPCProfilePersonaRow{}
	if err := repository.database.SelectContext(
		ctx,
		&rows,
		listCanonicalPersonasSQL,
		keyword,
		searchPattern,
		pluginFilter,
		pluginFilter,
		pageSize,
		(page-1)*pageSize,
	); err != nil {
		return MasterPersonaListResult{}, fmt.Errorf("list canonical personas: %w", err)
	}

	items := make([]MasterPersonaEntry, 0, len(rows))
	for _, row := range rows {
		identityKey := BuildMasterPersonaIdentityKey(row.TargetPluginName, row.FormID, row.RecordType)
		entry, err := fromCanonicalNPCProfilePersonaRow(identityKey, row)
		if err != nil {
			return MasterPersonaListResult{}, err
		}
		items = append(items, entry)
	}

	pluginGroupRows := []sqliteMasterPersonaPluginGroupRow{}
	if err := repository.database.SelectContext(
		ctx,
		&pluginGroupRows,
		listCanonicalPluginGroupsSQL,
		keyword,
		searchPattern,
	); err != nil {
		return MasterPersonaListResult{}, fmt.Errorf("list canonical persona plugin groups: %w", err)
	}

	pluginGroups := make([]MasterPersonaPluginGroup, 0, len(pluginGroupRows))
	for _, row := range pluginGroupRows {
		pluginGroups = append(pluginGroups, MasterPersonaPluginGroup(row))
	}

	return MasterPersonaListResult{
		Items:        items,
		TotalCount:   totalCount,
		Page:         page,
		PageSize:     pageSize,
		PluginGroups: pluginGroups,
	}, nil
}

// GetByIdentityKey loads one master persona entry by identity key from canonical NPC_PROFILE + PERSONA.
func (repository *SQLiteMasterPersonaEntryRepository) GetByIdentityKey(
	ctx context.Context,
	identityKey string,
) (MasterPersonaEntry, error) {
	trimmedKey := strings.TrimSpace(identityKey)
	targetPlugin, formID, recordType, ok := parseMasterPersonaIdentityKey(trimmedKey)
	if !ok {
		return MasterPersonaEntry{}, fmt.Errorf(masterPersonaIdentityKeyErrFormat, ErrMasterPersonaEntryNotFound, identityKey)
	}
	row := canonicalNPCProfilePersonaRow{}
	if err := repository.database.GetContext(ctx, &row, getCanonicalMasterPersonaByIdentitySQL, targetPlugin, formID, recordType); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return MasterPersonaEntry{}, fmt.Errorf(masterPersonaIdentityKeyErrFormat, ErrMasterPersonaEntryNotFound, identityKey)
		}
		return MasterPersonaEntry{}, fmt.Errorf("get canonical master persona entry by identity key: %w", err)
	}
	return fromCanonicalNPCProfilePersonaRow(trimmedKey, row)
}

// UpsertIfAbsent inserts NPC_PROFILE + PERSONA atomically only when the canonical identity does not exist.
// When an active *sqlx.Tx is found in ctx via TxKey, that transaction is reused so the caller's outer
// transaction covers all per-NPC writes.  Without an outer tx a self-contained transaction is started.
func (repository *SQLiteMasterPersonaEntryRepository) UpsertIfAbsent(
	ctx context.Context,
	draft MasterPersonaDraft,
) (MasterPersonaEntry, bool, error) {
	if tx, ok := ctx.Value(TxKey{}).(*sqlx.Tx); ok {
		return repository.upsertIfAbsentInTx(ctx, tx, draft)
	}
	return repository.upsertIfAbsentWithOwnTx(ctx, draft)
}

// upsertIfAbsentInTx runs the canonical write using a caller-owned *sqlx.Tx.
// It does not COMMIT or ROLLBACK; the caller is responsible for the transaction boundary.
func (repository *SQLiteMasterPersonaEntryRepository) upsertIfAbsentInTx(
	ctx context.Context,
	tx *sqlx.Tx,
	draft MasterPersonaDraft,
) (MasterPersonaEntry, bool, error) {
	normalized := entryFromDraft(draft)
	timestamp := normalized.UpdatedAt.UTC().Format(masterPersonaTimestampLayout)

	result, err := tx.ExecContext(ctx, insertCanonicalNPCProfileIfAbsentSQL,
		normalized.TargetPlugin,
		normalized.FormID,
		normalized.RecordType,
		normalized.EditorID,
		normalized.DisplayName,
		timestamp,
		timestamp,
	)
	if err != nil {
		return MasterPersonaEntry{}, false, fmt.Errorf("insert canonical npc profile if absent: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return MasterPersonaEntry{}, false, fmt.Errorf("read canonical npc profile rows affected: %w", err)
	}

	if rowsAffected == 0 {
		// NPC_PROFILE already exists (committed before this tx); load and return it without disturbing the outer tx.
		existing, loadErr := repository.GetByIdentityKey(ctx, normalized.IdentityKey)
		if loadErr != nil {
			return MasterPersonaEntry{}, false, loadErr
		}
		return existing, false, nil
	}

	npcProfileID, err := result.LastInsertId()
	if err != nil {
		return MasterPersonaEntry{}, false, fmt.Errorf("get canonical npc profile id after insert: %w", err)
	}

	if _, err := tx.ExecContext(ctx, insertCanonicalPersonaSQL,
		npcProfileID,
		normalized.PersonaBody,
		normalized.PersonaSummary,
		normalized.SpeechStyle,
		timestamp,
		timestamp,
	); err != nil {
		return MasterPersonaEntry{}, false, fmt.Errorf("insert canonical persona: %w", err)
	}

	return normalized, true, nil
}

// upsertIfAbsentWithOwnTx runs the canonical write inside a self-contained transaction.
func (repository *SQLiteMasterPersonaEntryRepository) upsertIfAbsentWithOwnTx(
	ctx context.Context,
	draft MasterPersonaDraft,
) (MasterPersonaEntry, bool, error) {
	normalized := entryFromDraft(draft)
	timestamp := normalized.UpdatedAt.UTC().Format(masterPersonaTimestampLayout)

	tx, err := repository.database.BeginTxx(ctx, nil)
	if err != nil {
		return MasterPersonaEntry{}, false, fmt.Errorf("begin canonical persona write transaction: %w", err)
	}

	result, err := tx.ExecContext(ctx, insertCanonicalNPCProfileIfAbsentSQL,
		normalized.TargetPlugin,
		normalized.FormID,
		normalized.RecordType,
		normalized.EditorID,
		normalized.DisplayName,
		timestamp,
		timestamp,
	)
	if err != nil {
		_ = tx.Rollback()
		return MasterPersonaEntry{}, false, fmt.Errorf("insert canonical npc profile if absent: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		_ = tx.Rollback()
		return MasterPersonaEntry{}, false, fmt.Errorf("read canonical npc profile rows affected: %w", err)
	}

	if rowsAffected == 0 {
		if rollbackErr := tx.Rollback(); rollbackErr != nil {
			return MasterPersonaEntry{}, false, fmt.Errorf("rollback canonical persona existing check: %w", rollbackErr)
		}
		existing, loadErr := repository.GetByIdentityKey(ctx, normalized.IdentityKey)
		if loadErr != nil {
			return MasterPersonaEntry{}, false, loadErr
		}
		return existing, false, nil
	}

	npcProfileID, err := result.LastInsertId()
	if err != nil {
		_ = tx.Rollback()
		return MasterPersonaEntry{}, false, fmt.Errorf("get canonical npc profile id after insert: %w", err)
	}

	if _, err := tx.ExecContext(ctx, insertCanonicalPersonaSQL,
		npcProfileID,
		normalized.PersonaBody,
		normalized.PersonaSummary,
		normalized.SpeechStyle,
		timestamp,
		timestamp,
	); err != nil {
		_ = tx.Rollback()
		return MasterPersonaEntry{}, false, fmt.Errorf("insert canonical persona: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return MasterPersonaEntry{}, false, fmt.Errorf("commit canonical persona write: %w", err)
	}
	return normalized, true, nil
}

// Update replaces one existing entry's canonical fields by identity key.
func (repository *SQLiteMasterPersonaEntryRepository) Update(
	ctx context.Context,
	identityKey string,
	draft MasterPersonaDraft,
) (MasterPersonaEntry, error) {
	currentIdentityKey := strings.TrimSpace(identityKey)
	_, err := repository.GetByIdentityKey(ctx, currentIdentityKey)
	if err != nil {
		return MasterPersonaEntry{}, err
	}

	normalized := entryFromDraft(draft)
	if normalized.IdentityKey != currentIdentityKey {
		_, duplicateErr := repository.GetByIdentityKey(ctx, normalized.IdentityKey)
		switch {
		case duplicateErr == nil:
			return MasterPersonaEntry{}, fmt.Errorf("master persona identity key already exists: %s", normalized.IdentityKey)
		case errors.Is(duplicateErr, ErrMasterPersonaEntryNotFound):
			// continue
		default:
			return MasterPersonaEntry{}, fmt.Errorf("check duplicate master persona identity key: %w", duplicateErr)
		}
	}

	currentPlugin, currentFormID, currentRecordType, ok := parseMasterPersonaIdentityKey(currentIdentityKey)
	if !ok {
		return MasterPersonaEntry{}, fmt.Errorf(masterPersonaIdentityKeyErrFormat, ErrMasterPersonaEntryNotFound, currentIdentityKey)
	}
	timestamp := normalized.UpdatedAt.UTC().Format(masterPersonaTimestampLayout)

	result, err := repository.database.ExecContext(ctx, updateCanonicalNPCProfileSQL,
		normalized.TargetPlugin,
		normalized.FormID,
		normalized.RecordType,
		normalized.EditorID,
		normalized.DisplayName,
		timestamp,
		currentPlugin,
		currentFormID,
		currentRecordType,
	)
	if err != nil {
		return MasterPersonaEntry{}, fmt.Errorf("update canonical npc profile: %w", err)
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return MasterPersonaEntry{}, fmt.Errorf("read updated canonical npc profile rows affected: %w", err)
	}
	if rowsAffected == 0 {
		return MasterPersonaEntry{}, fmt.Errorf(masterPersonaIdentityKeyErrFormat, ErrMasterPersonaEntryNotFound, identityKey)
	}

	if _, err := repository.database.ExecContext(ctx, updateCanonicalPersonaSQL,
		normalized.PersonaBody,
		normalized.PersonaSummary,
		normalized.SpeechStyle,
		timestamp,
		normalized.TargetPlugin,
		normalized.FormID,
		normalized.RecordType,
	); err != nil {
		return MasterPersonaEntry{}, fmt.Errorf("update canonical persona: %w", err)
	}

	return repository.GetByIdentityKey(ctx, normalized.IdentityKey)
}

// Delete removes one master persona entry from the canonical NPC_PROFILE + PERSONA schema.
func (repository *SQLiteMasterPersonaEntryRepository) Delete(ctx context.Context, identityKey string) error {
	trimmedKey := strings.TrimSpace(identityKey)
	targetPlugin, formID, recordType, ok := parseMasterPersonaIdentityKey(trimmedKey)
	if !ok {
		return fmt.Errorf(masterPersonaIdentityKeyErrFormat, ErrMasterPersonaEntryNotFound, identityKey)
	}

	// Delete PERSONA first to satisfy foreign key constraint before removing NPC_PROFILE.
	if _, err := repository.database.ExecContext(ctx, deleteCanonicalPersonaSQL, targetPlugin, formID, recordType); err != nil {
		return fmt.Errorf("delete canonical persona: %w", err)
	}

	result, err := repository.database.ExecContext(ctx, deleteCanonicalNPCProfileSQL, targetPlugin, formID, recordType)
	if err != nil {
		return fmt.Errorf("delete canonical npc profile: %w", err)
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("read deleted canonical npc profile rows affected: %w", err)
	}
	if rowsAffected == 0 {
		return fmt.Errorf(masterPersonaIdentityKeyErrFormat, ErrMasterPersonaEntryNotFound, identityKey)
	}
	return nil
}

// LoadAISettings returns page-local AI settings.
func (repository *SQLiteMasterPersonaAISettingsRepository) LoadAISettings(
	ctx context.Context,
) (MasterPersonaAISettingsRecord, error) {
	row := sqliteMasterPersonaAISettingsRow{}
	if err := repository.database.GetContext(ctx, &row, loadMasterPersonaAISettingsSQL); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return MasterPersonaAISettingsRecord{}, nil
		}
		return MasterPersonaAISettingsRecord{}, fmt.Errorf("load master persona ai settings: %w", err)
	}
	return MasterPersonaAISettingsRecord{
		Provider: strings.TrimSpace(row.Provider),
		Model:    strings.TrimSpace(row.Model),
	}, nil
}

// SaveAISettings saves page-local AI settings.
func (repository *SQLiteMasterPersonaAISettingsRepository) SaveAISettings(
	ctx context.Context,
	record MasterPersonaAISettingsRecord,
) error {
	if _, err := repository.database.ExecContext(
		ctx,
		upsertMasterPersonaAISettingsSQL,
		strings.TrimSpace(record.Provider),
		strings.TrimSpace(record.Model),
	); err != nil {
		return fmt.Errorf("save master persona ai settings: %w", err)
	}
	return nil
}

// LoadRunStatus returns the default idle run state. Run state is not persisted to DB.
func (repository *SQLiteMasterPersonaRunStatusRepository) LoadRunStatus(
	_ context.Context,
) (MasterPersonaRunStatusRecord, error) {
	return MasterPersonaRunStatusRecord{RunState: masterPersonaDefaultRunState}, nil
}

// SaveRunStatus is a no-op. Run state is not persisted to DB.
func (repository *SQLiteMasterPersonaRunStatusRepository) SaveRunStatus(
	_ context.Context,
	_ MasterPersonaRunStatusRecord,
) error {
	return nil
}

func (repository *SQLiteMasterPersonaEntryRepository) seedIfEmpty(
	ctx context.Context,
	seed []MasterPersonaEntry,
) error {
	if len(seed) == 0 {
		return nil
	}

	var count int
	if err := repository.database.GetContext(ctx, &count, countCanonicalNPCProfileRowsSQL); err != nil {
		return fmt.Errorf("count canonical npc profile seed target rows: %w", err)
	}
	if count > 0 {
		return nil
	}

	transaction, err := repository.database.BeginTxx(ctx, nil)
	if err != nil {
		return fmt.Errorf("begin master persona seed transaction: %w", err)
	}

	for _, seedEntry := range seed {
		draft := MasterPersonaDraft{
			IdentityKey:    seedEntry.IdentityKey,
			TargetPlugin:   seedEntry.TargetPlugin,
			FormID:         seedEntry.FormID,
			RecordType:     seedEntry.RecordType,
			EditorID:       seedEntry.EditorID,
			DisplayName:    seedEntry.DisplayName,
			PersonaBody:    seedEntry.PersonaBody,
			PersonaSummary: seedEntry.PersonaSummary,
			UpdatedAt:      seedEntry.UpdatedAt,
		}
		if _, _, err := repository.upsertIfAbsentInTx(ctx, transaction, draft); err != nil {
			if rollbackErr := transaction.Rollback(); rollbackErr != nil {
				return fmt.Errorf("insert master persona seed entry: %w", errors.Join(err, rollbackErr))
			}
			return fmt.Errorf("insert master persona seed entry: %w", err)
		}
	}
	if err := transaction.Commit(); err != nil {
		return fmt.Errorf("commit master persona seed transaction: %w", err)
	}
	return nil
}

func parseMasterPersonaIdentityKey(identityKey string) (targetPlugin, formID, recordType string, ok bool) {
	parts := strings.SplitN(identityKey, ":", 3)
	if len(parts) != 3 {
		return "", "", "", false
	}
	return parts[0], parts[1], parts[2], true
}

func fromCanonicalNPCProfilePersonaRow(identityKey string, row canonicalNPCProfilePersonaRow) (MasterPersonaEntry, error) {
	updatedAt, err := time.Parse(masterPersonaTimestampLayout, row.UpdatedAt)
	if err != nil {
		return MasterPersonaEntry{}, fmt.Errorf("parse updated_at for canonical master persona entry %s: %w", identityKey, err)
	}
	return MasterPersonaEntry{
		IdentityKey:    identityKey,
		TargetPlugin:   row.TargetPluginName,
		FormID:         row.FormID,
		RecordType:     row.RecordType,
		EditorID:       row.EditorID,
		DisplayName:    row.DisplayName,
		Race:           row.Race,
		Sex:            row.Sex,
		VoiceType:      row.VoiceType,
		ClassName:      row.NPCClass,
		SourcePlugin:   row.TargetPluginName,
		PersonaBody:    row.PersonaDescription,
		PersonaSummary: row.PersonaSummary,
		SpeechStyle:    row.SpeechStyle,
		Dialogues:      []string{},
		UpdatedAt:      updatedAt.UTC(),
	}, nil
}
