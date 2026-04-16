package repository

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	sqliteinfra "aitranslationenginejp/internal/infra/sqlite"
	"github.com/jmoiron/sqlx"
)

const (
	masterPersonaTimestampLayout      = time.RFC3339Nano
	masterPersonaDefaultRunState      = "入力待ち"
	masterPersonaIdentityKeyErrFormat = "%w: identity_key=%s"

	countMasterPersonaEntriesSQL = `
SELECT COUNT(*)
FROM master_persona_entries
WHERE (
  ? = ''
  OR lower(display_name || ' ' || form_id || ' ' || editor_id || ' ' || ifnull(race, '') || ' ' || voice_type) LIKE ?
)
AND (
  ? = ''
  OR target_plugin = ?
);`
	listMasterPersonaEntriesSQL = `
SELECT
  identity_key,
  target_plugin,
  form_id,
  record_type,
  editor_id,
  display_name,
  race,
  sex,
  voice_type,
  class_name,
  source_plugin,
  persona_summary,
  persona_body,
  generation_source_json,
  baseline_applied,
  dialogue_count,
  dialogues_json,
  updated_at
FROM master_persona_entries
WHERE (
  ? = ''
  OR lower(display_name || ' ' || form_id || ' ' || editor_id || ' ' || ifnull(race, '') || ' ' || voice_type) LIKE ?
)
AND (
  ? = ''
  OR target_plugin = ?
)
ORDER BY updated_at DESC, identity_key ASC
LIMIT ?
OFFSET ?;`
	listMasterPersonaPluginGroupsSQL = `
SELECT target_plugin, COUNT(*) AS count
FROM master_persona_entries
WHERE (
  ? = ''
  OR lower(display_name || ' ' || form_id || ' ' || editor_id || ' ' || ifnull(race, '') || ' ' || voice_type) LIKE ?
)
GROUP BY target_plugin
ORDER BY target_plugin ASC;`
	getMasterPersonaEntryByIdentityKeySQL = `
SELECT
  identity_key,
  target_plugin,
  form_id,
  record_type,
  editor_id,
  display_name,
  race,
  sex,
  voice_type,
  class_name,
  source_plugin,
  persona_summary,
  persona_body,
  generation_source_json,
  baseline_applied,
  dialogue_count,
  dialogues_json,
  updated_at
FROM master_persona_entries
WHERE identity_key = ?
LIMIT 1;`
	insertMasterPersonaEntryIfAbsentSQL = `
INSERT INTO master_persona_entries (
  identity_key,
  target_plugin,
  form_id,
  record_type,
  editor_id,
  display_name,
  race,
  sex,
  voice_type,
  class_name,
  source_plugin,
  persona_summary,
  persona_body,
  generation_source_json,
  baseline_applied,
  dialogue_count,
  dialogues_json,
  updated_at
) VALUES (
  :identity_key,
  :target_plugin,
  :form_id,
  :record_type,
  :editor_id,
  :display_name,
  :race,
  :sex,
  :voice_type,
  :class_name,
  :source_plugin,
  :persona_summary,
  :persona_body,
  :generation_source_json,
  :baseline_applied,
  :dialogue_count,
  :dialogues_json,
  :updated_at
)
ON CONFLICT(identity_key) DO NOTHING;`
	updateMasterPersonaEntrySQL = `
UPDATE master_persona_entries
SET identity_key = :next_identity_key,
    target_plugin = :target_plugin,
    form_id = :form_id,
    record_type = :record_type,
    editor_id = :editor_id,
    display_name = :display_name,
    race = :race,
    sex = :sex,
    voice_type = :voice_type,
    class_name = :class_name,
    source_plugin = :source_plugin,
    persona_summary = :persona_summary,
    persona_body = :persona_body,
    generation_source_json = :generation_source_json,
    baseline_applied = :baseline_applied,
    dialogue_count = :dialogue_count,
    dialogues_json = :dialogues_json,
    updated_at = :updated_at
WHERE identity_key = :current_identity_key;`
	deleteMasterPersonaEntrySQL = `
DELETE FROM master_persona_entries
WHERE identity_key = ?;`
	countMasterPersonaSeedRowsSQL = `
SELECT COUNT(*)
FROM master_persona_entries;`

	loadMasterPersonaAISettingsSQL = `
SELECT provider, model
FROM master_persona_ai_settings
WHERE id = 1
LIMIT 1;`
	upsertMasterPersonaAISettingsSQL = `
INSERT INTO master_persona_ai_settings (id, provider, model)
VALUES (1, ?, ?)
ON CONFLICT(id) DO UPDATE
SET provider = excluded.provider,
    model = excluded.model;`

	loadMasterPersonaRunStatusSQL = `
SELECT
  run_state,
  target_plugin,
  processed_count,
  success_count,
  existing_skip_count,
  zero_dialogue_skip_count,
  generic_npc_count,
  current_actor_label,
  message,
  started_at,
  finished_at
FROM master_persona_run_status
WHERE id = 1
LIMIT 1;`
	upsertMasterPersonaRunStatusSQL = `
INSERT INTO master_persona_run_status (
  id,
  run_state,
  target_plugin,
  processed_count,
  success_count,
  existing_skip_count,
  zero_dialogue_skip_count,
  generic_npc_count,
  current_actor_label,
  message,
  started_at,
  finished_at
) VALUES (
  1,
  ?,
  ?,
  ?,
  ?,
  ?,
  ?,
  ?,
  ?,
  ?,
  ?,
  ?
)
ON CONFLICT(id) DO UPDATE
SET run_state = excluded.run_state,
    target_plugin = excluded.target_plugin,
    processed_count = excluded.processed_count,
    success_count = excluded.success_count,
    existing_skip_count = excluded.existing_skip_count,
    zero_dialogue_skip_count = excluded.zero_dialogue_skip_count,
    generic_npc_count = excluded.generic_npc_count,
    current_actor_label = excluded.current_actor_label,
    message = excluded.message,
    started_at = excluded.started_at,
    finished_at = excluded.finished_at;`
	insertMasterPersonaDefaultRunStatusSQL = `
INSERT INTO master_persona_run_status (id, run_state)
VALUES (1, ?)
ON CONFLICT(id) DO NOTHING;`
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

type sqliteMasterPersonaEntryRow struct {
	IdentityKey          string         `db:"identity_key"`
	TargetPlugin         string         `db:"target_plugin"`
	FormID               string         `db:"form_id"`
	RecordType           string         `db:"record_type"`
	EditorID             string         `db:"editor_id"`
	DisplayName          string         `db:"display_name"`
	Race                 sql.NullString `db:"race"`
	Sex                  sql.NullString `db:"sex"`
	VoiceType            string         `db:"voice_type"`
	ClassName            string         `db:"class_name"`
	SourcePlugin         string         `db:"source_plugin"`
	PersonaSummary       string         `db:"persona_summary"`
	PersonaBody          string         `db:"persona_body"`
	GenerationSourceJSON string         `db:"generation_source_json"`
	BaselineApplied      int            `db:"baseline_applied"`
	DialogueCount        int            `db:"dialogue_count"`
	DialoguesJSON        string         `db:"dialogues_json"`
	UpdatedAt            string         `db:"updated_at"`
}

type sqliteMasterPersonaEntryMutationParams struct {
	IdentityKey          string  `db:"identity_key"`
	CurrentIdentityKey   string  `db:"current_identity_key"`
	NextIdentityKey      string  `db:"next_identity_key"`
	TargetPlugin         string  `db:"target_plugin"`
	FormID               string  `db:"form_id"`
	RecordType           string  `db:"record_type"`
	EditorID             string  `db:"editor_id"`
	DisplayName          string  `db:"display_name"`
	Race                 *string `db:"race"`
	Sex                  *string `db:"sex"`
	VoiceType            string  `db:"voice_type"`
	ClassName            string  `db:"class_name"`
	SourcePlugin         string  `db:"source_plugin"`
	PersonaSummary       string  `db:"persona_summary"`
	PersonaBody          string  `db:"persona_body"`
	GenerationSourceJSON string  `db:"generation_source_json"`
	BaselineApplied      int     `db:"baseline_applied"`
	DialogueCount        int     `db:"dialogue_count"`
	DialoguesJSON        string  `db:"dialogues_json"`
	UpdatedAt            string  `db:"updated_at"`
}

type sqliteMasterPersonaPluginGroupRow struct {
	TargetPlugin string `db:"target_plugin"`
	Count        int    `db:"count"`
}

type sqliteMasterPersonaAISettingsRow struct {
	Provider string `db:"provider"`
	Model    string `db:"model"`
}

type sqliteMasterPersonaRunStatusRow struct {
	RunState              string         `db:"run_state"`
	TargetPlugin          string         `db:"target_plugin"`
	ProcessedCount        int            `db:"processed_count"`
	SuccessCount          int            `db:"success_count"`
	ExistingSkipCount     int            `db:"existing_skip_count"`
	ZeroDialogueSkipCount int            `db:"zero_dialogue_skip_count"`
	GenericNPCCount       int            `db:"generic_npc_count"`
	CurrentActorLabel     string         `db:"current_actor_label"`
	Message               string         `db:"message"`
	StartedAt             sql.NullString `db:"started_at"`
	FinishedAt            sql.NullString `db:"finished_at"`
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
	if err := runStatusRepository.ensureDefaultStatus(ctx); err != nil {
		if closeErr := database.Close(); closeErr != nil {
			return nil, fmt.Errorf("initialize sqlite master persona run status: %w", errors.Join(err, closeErr))
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

// List returns filtered and paginated master persona entries.
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
		countMasterPersonaEntriesSQL,
		keyword,
		searchPattern,
		pluginFilter,
		pluginFilter,
	); err != nil {
		return MasterPersonaListResult{}, fmt.Errorf("count master persona entries: %w", err)
	}

	page, pageSize := normalizeMasterPersonaPagination(query.Page, query.PageSize, totalCount)
	rows := []sqliteMasterPersonaEntryRow{}
	if err := repository.database.SelectContext(
		ctx,
		&rows,
		listMasterPersonaEntriesSQL,
		keyword,
		searchPattern,
		pluginFilter,
		pluginFilter,
		pageSize,
		(page-1)*pageSize,
	); err != nil {
		return MasterPersonaListResult{}, fmt.Errorf("list master persona entries: %w", err)
	}

	items := make([]MasterPersonaEntry, 0, len(rows))
	for _, row := range rows {
		entry, err := fromSQLiteMasterPersonaEntryRow(row)
		if err != nil {
			return MasterPersonaListResult{}, err
		}
		items = append(items, entry)
	}

	pluginGroupRows := []sqliteMasterPersonaPluginGroupRow{}
	if err := repository.database.SelectContext(
		ctx,
		&pluginGroupRows,
		listMasterPersonaPluginGroupsSQL,
		keyword,
		searchPattern,
	); err != nil {
		return MasterPersonaListResult{}, fmt.Errorf("list master persona plugin groups: %w", err)
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

// GetByIdentityKey loads one master persona entry by identity key.
func (repository *SQLiteMasterPersonaEntryRepository) GetByIdentityKey(
	ctx context.Context,
	identityKey string,
) (MasterPersonaEntry, error) {
	trimmedIdentityKey := strings.TrimSpace(identityKey)
	row := sqliteMasterPersonaEntryRow{}
	if err := repository.database.GetContext(ctx, &row, getMasterPersonaEntryByIdentityKeySQL, trimmedIdentityKey); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return MasterPersonaEntry{}, fmt.Errorf(masterPersonaIdentityKeyErrFormat, ErrMasterPersonaEntryNotFound, identityKey)
		}
		return MasterPersonaEntry{}, fmt.Errorf("get master persona entry by identity key: %w", err)
	}
	return fromSQLiteMasterPersonaEntryRow(row)
}

// UpsertIfAbsent inserts one entry only when identity key does not exist.
func (repository *SQLiteMasterPersonaEntryRepository) UpsertIfAbsent(
	ctx context.Context,
	draft MasterPersonaDraft,
) (MasterPersonaEntry, bool, error) {
	normalized := entryFromDraft(draft)
	params, err := toSQLiteMasterPersonaEntryMutationParams(normalized.IdentityKey, normalized)
	if err != nil {
		return MasterPersonaEntry{}, false, err
	}
	result, err := repository.database.NamedExecContext(ctx, insertMasterPersonaEntryIfAbsentSQL, params)
	if err != nil {
		return MasterPersonaEntry{}, false, fmt.Errorf("insert master persona entry if absent: %w", err)
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return MasterPersonaEntry{}, false, fmt.Errorf("read inserted master persona rows affected: %w", err)
	}
	if rowsAffected == 0 {
		existing, loadErr := repository.GetByIdentityKey(ctx, normalized.IdentityKey)
		if loadErr != nil {
			return MasterPersonaEntry{}, false, loadErr
		}
		return existing, false, nil
	}
	created, loadErr := repository.GetByIdentityKey(ctx, normalized.IdentityKey)
	if loadErr != nil {
		return MasterPersonaEntry{}, false, loadErr
	}
	return created, true, nil
}

// Update replaces one existing entry by identity key.
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

	params, err := toSQLiteMasterPersonaEntryMutationParams(currentIdentityKey, normalized)
	if err != nil {
		return MasterPersonaEntry{}, err
	}
	result, err := repository.database.NamedExecContext(ctx, updateMasterPersonaEntrySQL, params)
	if err != nil {
		return MasterPersonaEntry{}, fmt.Errorf("update master persona entry: %w", err)
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return MasterPersonaEntry{}, fmt.Errorf("read updated master persona rows affected: %w", err)
	}
	if rowsAffected == 0 {
		return MasterPersonaEntry{}, fmt.Errorf(masterPersonaIdentityKeyErrFormat, ErrMasterPersonaEntryNotFound, identityKey)
	}
	return repository.GetByIdentityKey(ctx, normalized.IdentityKey)
}

// Delete removes one master persona entry by identity key.
func (repository *SQLiteMasterPersonaEntryRepository) Delete(ctx context.Context, identityKey string) error {
	result, err := repository.database.ExecContext(ctx, deleteMasterPersonaEntrySQL, strings.TrimSpace(identityKey))
	if err != nil {
		return fmt.Errorf("delete master persona entry: %w", err)
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("read deleted master persona rows affected: %w", err)
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

// LoadRunStatus returns persisted generation run status.
func (repository *SQLiteMasterPersonaRunStatusRepository) LoadRunStatus(
	ctx context.Context,
) (MasterPersonaRunStatusRecord, error) {
	row := sqliteMasterPersonaRunStatusRow{}
	if err := repository.database.GetContext(ctx, &row, loadMasterPersonaRunStatusSQL); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return MasterPersonaRunStatusRecord{RunState: masterPersonaDefaultRunState}, nil
		}
		return MasterPersonaRunStatusRecord{}, fmt.Errorf("load master persona run status: %w", err)
	}
	startedAt, err := parseMasterPersonaNullableTimestamp(row.StartedAt)
	if err != nil {
		return MasterPersonaRunStatusRecord{}, err
	}
	finishedAt, err := parseMasterPersonaNullableTimestamp(row.FinishedAt)
	if err != nil {
		return MasterPersonaRunStatusRecord{}, err
	}
	runState := strings.TrimSpace(row.RunState)
	if runState == "" {
		runState = masterPersonaDefaultRunState
	}
	return MasterPersonaRunStatusRecord{
		RunState:              runState,
		TargetPlugin:          strings.TrimSpace(row.TargetPlugin),
		ProcessedCount:        row.ProcessedCount,
		SuccessCount:          row.SuccessCount,
		ExistingSkipCount:     row.ExistingSkipCount,
		ZeroDialogueSkipCount: row.ZeroDialogueSkipCount,
		GenericNPCCount:       row.GenericNPCCount,
		CurrentActorLabel:     strings.TrimSpace(row.CurrentActorLabel),
		Message:               strings.TrimSpace(row.Message),
		StartedAt:             startedAt,
		FinishedAt:            finishedAt,
	}, nil
}

// SaveRunStatus stores generation run status.
func (repository *SQLiteMasterPersonaRunStatusRepository) SaveRunStatus(
	ctx context.Context,
	status MasterPersonaRunStatusRecord,
) error {
	normalized := cloneMasterPersonaRunStatus(status)
	normalized.RunState = strings.TrimSpace(normalized.RunState)
	if normalized.RunState == "" {
		normalized.RunState = masterPersonaDefaultRunState
	}
	if _, err := repository.database.ExecContext(
		ctx,
		upsertMasterPersonaRunStatusSQL,
		normalized.RunState,
		strings.TrimSpace(normalized.TargetPlugin),
		normalized.ProcessedCount,
		normalized.SuccessCount,
		normalized.ExistingSkipCount,
		normalized.ZeroDialogueSkipCount,
		normalized.GenericNPCCount,
		strings.TrimSpace(normalized.CurrentActorLabel),
		strings.TrimSpace(normalized.Message),
		masterPersonaNullableTimestamp(normalized.StartedAt),
		masterPersonaNullableTimestamp(normalized.FinishedAt),
	); err != nil {
		return fmt.Errorf("save master persona run status: %w", err)
	}
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
	if err := repository.database.GetContext(ctx, &count, countMasterPersonaSeedRowsSQL); err != nil {
		return fmt.Errorf("count master persona seed target rows: %w", err)
	}
	if count > 0 {
		return nil
	}

	transaction, err := repository.database.BeginTxx(ctx, nil)
	if err != nil {
		return fmt.Errorf("begin master persona seed transaction: %w", err)
	}

	for _, seedEntry := range seed {
		if err := insertMasterPersonaSeedEntry(ctx, transaction, seedEntry); err != nil {
			if rollbackErr := transaction.Rollback(); rollbackErr != nil {
				return fmt.Errorf("insert master persona seed entry: %w", errors.Join(err, rollbackErr))
			}
			return err
		}
	}
	if err := transaction.Commit(); err != nil {
		return fmt.Errorf("commit master persona seed transaction: %w", err)
	}
	return nil
}

func insertMasterPersonaSeedEntry(
	ctx context.Context,
	transaction *sqlx.Tx,
	seedEntry MasterPersonaEntry,
) error {
	normalized := entryFromDraft(masterPersonaDraftFromEntry(seedEntry))
	params, err := toSQLiteMasterPersonaEntryMutationParams(normalized.IdentityKey, normalized)
	if err != nil {
		return fmt.Errorf("prepare master persona seed entry: %w", err)
	}
	if _, err := transaction.NamedExecContext(ctx, insertMasterPersonaEntryIfAbsentSQL, params); err != nil {
		return fmt.Errorf("insert master persona seed entry: %w", err)
	}
	return nil
}

func masterPersonaDraftFromEntry(entry MasterPersonaEntry) MasterPersonaDraft {
	return MasterPersonaDraft{
		IdentityKey:          entry.IdentityKey,
		TargetPlugin:         entry.TargetPlugin,
		FormID:               entry.FormID,
		RecordType:           entry.RecordType,
		EditorID:             entry.EditorID,
		DisplayName:          entry.DisplayName,
		Race:                 entry.Race,
		Sex:                  entry.Sex,
		VoiceType:            entry.VoiceType,
		ClassName:            entry.ClassName,
		SourcePlugin:         entry.SourcePlugin,
		PersonaSummary:       entry.PersonaSummary,
		PersonaBody:          entry.PersonaBody,
		GenerationSourceJSON: entry.GenerationSourceJSON,
		BaselineApplied:      entry.BaselineApplied,
		Dialogues:            append([]string(nil), entry.Dialogues...),
		UpdatedAt:            entry.UpdatedAt,
	}
}

func (repository *SQLiteMasterPersonaRunStatusRepository) ensureDefaultStatus(ctx context.Context) error {
	if _, err := repository.database.ExecContext(ctx, insertMasterPersonaDefaultRunStatusSQL, masterPersonaDefaultRunState); err != nil {
		return fmt.Errorf("insert default master persona run status: %w", err)
	}
	return nil
}

func fromSQLiteMasterPersonaEntryRow(row sqliteMasterPersonaEntryRow) (MasterPersonaEntry, error) {
	updatedAt, err := time.Parse(masterPersonaTimestampLayout, row.UpdatedAt)
	if err != nil {
		return MasterPersonaEntry{}, fmt.Errorf("parse updated_at for master persona entry %s: %w", row.IdentityKey, err)
	}
	dialogues := []string{}
	if err := json.Unmarshal([]byte(row.DialoguesJSON), &dialogues); err != nil {
		return MasterPersonaEntry{}, fmt.Errorf("decode dialogues_json for master persona entry %s: %w", row.IdentityKey, err)
	}
	dialogueCount := row.DialogueCount
	if dialogueCount <= 0 {
		dialogueCount = len(dialogues)
	}
	return MasterPersonaEntry{
		IdentityKey:          row.IdentityKey,
		TargetPlugin:         row.TargetPlugin,
		FormID:               row.FormID,
		RecordType:           row.RecordType,
		EditorID:             row.EditorID,
		DisplayName:          row.DisplayName,
		Race:                 masterPersonaNullableString(row.Race),
		Sex:                  masterPersonaNullableString(row.Sex),
		VoiceType:            row.VoiceType,
		ClassName:            row.ClassName,
		SourcePlugin:         row.SourcePlugin,
		PersonaSummary:       row.PersonaSummary,
		PersonaBody:          row.PersonaBody,
		GenerationSourceJSON: row.GenerationSourceJSON,
		BaselineApplied:      row.BaselineApplied != 0,
		DialogueCount:        dialogueCount,
		Dialogues:            dialogues,
		UpdatedAt:            updatedAt.UTC(),
	}, nil
}

func toSQLiteMasterPersonaEntryMutationParams(
	currentIdentityKey string,
	entry MasterPersonaEntry,
) (sqliteMasterPersonaEntryMutationParams, error) {
	dialoguesJSON, err := json.Marshal(entry.Dialogues)
	if err != nil {
		return sqliteMasterPersonaEntryMutationParams{}, fmt.Errorf("encode dialogues for master persona entry %s: %w", entry.IdentityKey, err)
	}
	baselineApplied := 0
	if entry.BaselineApplied {
		baselineApplied = 1
	}
	return sqliteMasterPersonaEntryMutationParams{
		IdentityKey:          entry.IdentityKey,
		CurrentIdentityKey:   currentIdentityKey,
		NextIdentityKey:      entry.IdentityKey,
		TargetPlugin:         entry.TargetPlugin,
		FormID:               entry.FormID,
		RecordType:           entry.RecordType,
		EditorID:             entry.EditorID,
		DisplayName:          entry.DisplayName,
		Race:                 entry.Race,
		Sex:                  entry.Sex,
		VoiceType:            entry.VoiceType,
		ClassName:            entry.ClassName,
		SourcePlugin:         entry.SourcePlugin,
		PersonaSummary:       entry.PersonaSummary,
		PersonaBody:          entry.PersonaBody,
		GenerationSourceJSON: entry.GenerationSourceJSON,
		BaselineApplied:      baselineApplied,
		DialogueCount:        entry.DialogueCount,
		DialoguesJSON:        string(dialoguesJSON),
		UpdatedAt:            entry.UpdatedAt.UTC().Format(masterPersonaTimestampLayout),
	}, nil
}

func masterPersonaNullableString(value sql.NullString) *string {
	if !value.Valid {
		return nil
	}
	trimmed := strings.TrimSpace(value.String)
	if trimmed == "" {
		return nil
	}
	return &trimmed
}

func parseMasterPersonaNullableTimestamp(value sql.NullString) (*time.Time, error) {
	if !value.Valid {
		return nil, nil
	}
	trimmed := strings.TrimSpace(value.String)
	if trimmed == "" {
		return nil, nil
	}
	parsed, err := time.Parse(masterPersonaTimestampLayout, trimmed)
	if err != nil {
		return nil, fmt.Errorf("parse master persona nullable timestamp: %w", err)
	}
	utc := parsed.UTC()
	return &utc, nil
}

func masterPersonaNullableTimestamp(value *time.Time) interface{} {
	if value == nil {
		return nil
	}
	return value.UTC().Format(masterPersonaTimestampLayout)
}
