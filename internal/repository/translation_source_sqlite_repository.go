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

// SQLiteTranslationSourceRepository は TranslationSourceRepository の SQLite 実装。
type SQLiteTranslationSourceRepository struct {
	db *sqlx.DB
}

// NewSQLiteTranslationSourceRepository は TranslationSourceRepository を返す。
func NewSQLiteTranslationSourceRepository(db *sqlx.DB) TranslationSourceRepository {
	return &SQLiteTranslationSourceRepository{db: db}
}

// ---------------------------------------------------------------------------
// 内部 row 型 (db タグで SQL カラムとマッピング)
// ---------------------------------------------------------------------------

type xEditExtractedDataRow struct {
	ID                int64  `db:"id"`
	SourceFilePath    string `db:"source_file_path"`
	SourceContentHash string `db:"source_content_hash"`
	SourceTool        string `db:"source_tool"`
	TargetPluginName  string `db:"target_plugin_name"`
	TargetPluginType  string `db:"target_plugin_type"`
	RecordCount       int    `db:"record_count"`
	ImportedAt        string `db:"imported_at"`
}

func (r xEditExtractedDataRow) toModel() (XEditExtractedData, error) {
	importedAt, err := time.Parse(time.RFC3339, r.ImportedAt)
	if err != nil {
		return XEditExtractedData{}, fmt.Errorf("parse imported_at: %w", err)
	}
	return XEditExtractedData{
		ID:                r.ID,
		SourceFilePath:    r.SourceFilePath,
		SourceContentHash: r.SourceContentHash,
		SourceTool:        r.SourceTool,
		TargetPluginName:  r.TargetPluginName,
		TargetPluginType:  r.TargetPluginType,
		RecordCount:       r.RecordCount,
		ImportedAt:        importedAt,
	}, nil
}

type translationRecordRow struct {
	ID                   int64  `db:"id"`
	XEditExtractedDataID int64  `db:"x_edit_extracted_data_id"`
	FormID               string `db:"form_id"`
	EditorID             string `db:"editor_id"`
	RecordType           string `db:"record_type"`
}

type existingTranslationJobRow struct {
	ID    int64  `db:"id"`
	State string `db:"state"`
}

func (r translationRecordRow) toModel() TranslationRecord {
	return TranslationRecord(r)
}

type npcProfileRow struct {
	ID               int64  `db:"id"`
	TargetPluginName string `db:"target_plugin_name"`
	FormID           string `db:"form_id"`
	RecordType       string `db:"record_type"`
	EditorID         string `db:"editor_id"`
	DisplayName      string `db:"display_name"`
	CreatedAt        string `db:"created_at"`
	UpdatedAt        string `db:"updated_at"`
}

func (r npcProfileRow) toModel() (NpcProfile, error) {
	createdAt, err := time.Parse(time.RFC3339, r.CreatedAt)
	if err != nil {
		return NpcProfile{}, fmt.Errorf("parse created_at: %w", err)
	}
	updatedAt, err := time.Parse(time.RFC3339, r.UpdatedAt)
	if err != nil {
		return NpcProfile{}, fmt.Errorf("parse updated_at: %w", err)
	}
	return NpcProfile{
		ID:               r.ID,
		TargetPluginName: r.TargetPluginName,
		FormID:           r.FormID,
		RecordType:       r.RecordType,
		EditorID:         r.EditorID,
		DisplayName:      r.DisplayName,
		CreatedAt:        createdAt,
		UpdatedAt:        updatedAt,
	}, nil
}

type npcRecordRow struct {
	TranslationRecordID int64   `db:"translation_record_id"`
	NpcProfileID        int64   `db:"npc_profile_id"`
	Race                *string `db:"race"`
	Sex                 *string `db:"sex"`
	NpcClass            *string `db:"npc_class"`
	VoiceType           string  `db:"voice_type"`
}

func (r npcRecordRow) toModel() NpcRecord {
	return NpcRecord(r)
}

type translationFieldRow struct {
	ID                           int64  `db:"id"`
	TranslationRecordID          int64  `db:"translation_record_id"`
	TranslationFieldDefinitionID *int64 `db:"translation_field_definition_id"`
	SubrecordType                string `db:"subrecord_type"`
	SourceText                   string `db:"source_text"`
	FieldOrder                   int    `db:"field_order"`
	PreviousTranslationFieldID   *int64 `db:"previous_translation_field_id"`
	NextTranslationFieldID       *int64 `db:"next_translation_field_id"`
}

func (r translationFieldRow) toModel() TranslationField {
	return TranslationField(r)
}

type translationFieldRecordReferenceRow struct {
	ID                            int64  `db:"id"`
	TranslationFieldID            int64  `db:"translation_field_id"`
	ReferencedTranslationRecordID int64  `db:"referenced_translation_record_id"`
	ReferenceRole                 string `db:"reference_role"`
}

func (r translationFieldRecordReferenceRow) toModel() TranslationFieldRecordReference {
	return TranslationFieldRecordReference(r)
}

// ---------------------------------------------------------------------------
// エラー変換ヘルパー
// ---------------------------------------------------------------------------

const labeledErrorFmt = "%s: %w"

func isUniqueConstraintError(err error) bool {
	return strings.Contains(err.Error(), "UNIQUE constraint failed")
}

func mapSQLError(err error, label string) error {
	if errors.Is(err, sql.ErrNoRows) {
		return fmt.Errorf(labeledErrorFmt, label, ErrNotFound)
	}
	if isUniqueConstraintError(err) {
		return fmt.Errorf(labeledErrorFmt, label, ErrConflict)
	}
	return fmt.Errorf(labeledErrorFmt, label, err)
}

// ---------------------------------------------------------------------------
// SQL 定数
// ---------------------------------------------------------------------------

const (
	insertXEditExtractedData = `
INSERT INTO X_EDIT_EXTRACTED_DATA
	(source_file_path, source_content_hash, source_tool, target_plugin_name, target_plugin_type, record_count, imported_at)
VALUES
	(:source_file_path, :source_content_hash, :source_tool, :target_plugin_name, :target_plugin_type, :record_count, :imported_at)`

	selectXEditExtractedDataByID = `
SELECT id, source_file_path, source_content_hash, source_tool, target_plugin_name, target_plugin_type, record_count, imported_at
FROM X_EDIT_EXTRACTED_DATA WHERE id = ?`

	selectAllXEditExtractedData = `
SELECT id, source_file_path, source_content_hash, source_tool, target_plugin_name, target_plugin_type, record_count, imported_at
FROM X_EDIT_EXTRACTED_DATA
ORDER BY imported_at DESC, id DESC`

	selectXEditExtractedDataBySourceContentHash = `
SELECT id, source_file_path, source_content_hash, source_tool, target_plugin_name, target_plugin_type, record_count, imported_at
FROM X_EDIT_EXTRACTED_DATA WHERE source_content_hash = ?`

	updateXEditExtractedDataMetadata = `
UPDATE X_EDIT_EXTRACTED_DATA
SET source_file_path = ?,
	source_content_hash = ?,
	source_tool = ?,
	target_plugin_name = ?,
	target_plugin_type = ?,
	record_count = ?
WHERE id = ?`

	deleteTranslationFieldRecordReferencesByXEditID = `
DELETE FROM TRANSLATION_FIELD_RECORD_REFERENCE
WHERE translation_field_id IN (
	SELECT translation_field.id
	FROM TRANSLATION_FIELD AS translation_field
	INNER JOIN TRANSLATION_RECORD AS translation_record
		ON translation_record.id = translation_field.translation_record_id
	WHERE translation_record.x_edit_extracted_data_id = ?
)`

	deletePersonaFieldEvidenceByXEditID = `
DELETE FROM PERSONA_FIELD_EVIDENCE
	WHERE persona_id IN (
		SELECT id
		FROM PERSONA
		WHERE translation_job_id IN (
			SELECT id FROM TRANSLATION_JOB WHERE x_edit_extracted_data_id = ?
		)
	)
	OR translation_field_id IN (
	SELECT translation_field.id
	FROM TRANSLATION_FIELD AS translation_field
	INNER JOIN TRANSLATION_RECORD AS translation_record
		ON translation_record.id = translation_field.translation_record_id
	WHERE translation_record.x_edit_extracted_data_id = ?
)`

	deleteXTranslatorOutputRowsByXEditID = `
DELETE FROM XTRANSLATOR_OUTPUT_ROW
	WHERE translation_artifact_id IN (
		SELECT id
		FROM TRANSLATION_ARTIFACT
		WHERE translation_job_id IN (
			SELECT id FROM TRANSLATION_JOB WHERE x_edit_extracted_data_id = ?
		)
)`

	deletePhaseRunTranslationFieldsByXEditID = `
DELETE FROM PHASE_RUN_TRANSLATION_FIELD
	WHERE phase_run_id IN (
		SELECT id
		FROM JOB_PHASE_RUN
		WHERE translation_job_id IN (
			SELECT id FROM TRANSLATION_JOB WHERE x_edit_extracted_data_id = ?
		)
)`

	deletePhaseRunPersonasByXEditID = `
	DELETE FROM PHASE_RUN_PERSONA
	WHERE phase_run_id IN (
		SELECT id
		FROM JOB_PHASE_RUN
		WHERE translation_job_id IN (
			SELECT id FROM TRANSLATION_JOB WHERE x_edit_extracted_data_id = ?
		)
	)`

	deletePhaseRunDictionaryEntriesByXEditID = `
	DELETE FROM PHASE_RUN_DICTIONARY_ENTRY
	WHERE phase_run_id IN (
		SELECT id
		FROM JOB_PHASE_RUN
		WHERE translation_job_id IN (
			SELECT id FROM TRANSLATION_JOB WHERE x_edit_extracted_data_id = ?
		)
	)`

	deleteJobPhaseRunsByXEditID = `
	DELETE FROM JOB_PHASE_RUN
	WHERE translation_job_id IN (
		SELECT id FROM TRANSLATION_JOB WHERE x_edit_extracted_data_id = ?
	)`

	deleteJobTranslationFieldsByXEditID = `
DELETE FROM JOB_TRANSLATION_FIELD
WHERE translation_job_id IN (
	SELECT id FROM TRANSLATION_JOB WHERE x_edit_extracted_data_id = ?
)`

	deleteTranslationArtifactsByXEditID = `
	DELETE FROM TRANSLATION_ARTIFACT
	WHERE translation_job_id IN (
		SELECT id FROM TRANSLATION_JOB WHERE x_edit_extracted_data_id = ?
	)`

	deletePersonasByXEditID = `
	DELETE FROM PERSONA
	WHERE translation_job_id IN (
		SELECT id FROM TRANSLATION_JOB WHERE x_edit_extracted_data_id = ?
	)`

	deleteDictionaryEntriesByXEditID = `
	DELETE FROM DICTIONARY_ENTRY
	WHERE translation_job_id IN (
		SELECT id FROM TRANSLATION_JOB WHERE x_edit_extracted_data_id = ?
	)`

	deleteTranslationJobsByXEditID = `
	DELETE FROM TRANSLATION_JOB WHERE x_edit_extracted_data_id = ?`

	deleteNpcRecordsByXEditID = `
DELETE FROM NPC_RECORD
WHERE translation_record_id IN (
	SELECT id FROM TRANSLATION_RECORD WHERE x_edit_extracted_data_id = ?
)`

	deleteTranslationFieldsByXEditID = `
DELETE FROM TRANSLATION_FIELD
WHERE translation_record_id IN (
	SELECT id FROM TRANSLATION_RECORD WHERE x_edit_extracted_data_id = ?
)`

	deleteTranslationRecordsByXEditID = `
DELETE FROM TRANSLATION_RECORD WHERE x_edit_extracted_data_id = ?`

	insertTranslationRecord = `
INSERT INTO TRANSLATION_RECORD
  (x_edit_extracted_data_id, form_id, editor_id, record_type)
VALUES
  (:x_edit_extracted_data_id, :form_id, :editor_id, :record_type)`

	selectTranslationRecordByID = `
SELECT id, x_edit_extracted_data_id, form_id, editor_id, record_type
FROM TRANSLATION_RECORD WHERE id = ?`

	selectTranslationRecordsByXEditID = `
SELECT id, x_edit_extracted_data_id, form_id, editor_id, record_type
FROM TRANSLATION_RECORD WHERE x_edit_extracted_data_id = ?`

	selectExistingTranslationJob = `
SELECT id, state
FROM TRANSLATION_JOB
WHERE x_edit_extracted_data_id = ?
ORDER BY created_at DESC, id DESC
LIMIT 1`

	selectTranslationCachePresenceByXEditID = `
SELECT EXISTS(
	SELECT 1
	FROM TRANSLATION_RECORD
	WHERE x_edit_extracted_data_id = ?
)`

	upsertNpcProfile = `
INSERT INTO NPC_PROFILE
  (target_plugin_name, form_id, record_type, editor_id, display_name, created_at, updated_at)
VALUES
  (:target_plugin_name, :form_id, :record_type, :editor_id, :display_name, :created_at, :updated_at)
ON CONFLICT(target_plugin_name, form_id, record_type) DO UPDATE SET
  editor_id    = excluded.editor_id,
  display_name = excluded.display_name,
  updated_at   = excluded.updated_at`

	selectNpcProfileByUnique = `
SELECT id, target_plugin_name, form_id, record_type, editor_id, display_name, created_at, updated_at
FROM NPC_PROFILE WHERE target_plugin_name = ? AND form_id = ? AND record_type = ?`

	selectNpcProfileByID = `
SELECT id, target_plugin_name, form_id, record_type, editor_id, display_name, created_at, updated_at
FROM NPC_PROFILE WHERE id = ?`

	insertNpcRecord = `
INSERT INTO NPC_RECORD
  (translation_record_id, npc_profile_id, race, sex, npc_class, voice_type)
VALUES
  (:translation_record_id, :npc_profile_id, :race, :sex, :npc_class, :voice_type)`

	selectNpcRecordByTranslationRecordID = `
SELECT translation_record_id, npc_profile_id, race, sex, npc_class, voice_type
FROM NPC_RECORD WHERE translation_record_id = ?`

	insertTranslationField = `
INSERT INTO TRANSLATION_FIELD
  (translation_record_id, translation_field_definition_id, subrecord_type, source_text,
   field_order, previous_translation_field_id, next_translation_field_id)
VALUES
  (:translation_record_id, :translation_field_definition_id, :subrecord_type, :source_text,
   :field_order, :previous_translation_field_id, :next_translation_field_id)`

	selectTranslationFieldByID = `
SELECT id, translation_record_id, translation_field_definition_id, subrecord_type, source_text,
       field_order, previous_translation_field_id, next_translation_field_id
FROM TRANSLATION_FIELD WHERE id = ?`

	selectTranslationFieldsByTranslationRecordID = `
SELECT id, translation_record_id, translation_field_definition_id, subrecord_type, source_text,
       field_order, previous_translation_field_id, next_translation_field_id
FROM TRANSLATION_FIELD WHERE translation_record_id = ?`

	insertTranslationFieldRecordReference = `
INSERT INTO TRANSLATION_FIELD_RECORD_REFERENCE
  (translation_field_id, referenced_translation_record_id, reference_role)
VALUES
  (:translation_field_id, :referenced_translation_record_id, :reference_role)`

	selectTranslationFieldRecordReferencesByFieldID = `
SELECT id, translation_field_id, referenced_translation_record_id, reference_role
FROM TRANSLATION_FIELD_RECORD_REFERENCE WHERE translation_field_id = ?`
)

// ---------------------------------------------------------------------------
// XEditExtractedData
// ---------------------------------------------------------------------------

// CreateXEditExtractedData は XEditExtractedData レコードを作成する。
func (r *SQLiteTranslationSourceRepository) CreateXEditExtractedData(
	ctx context.Context,
	draft XEditExtractedDataDraft,
) (XEditExtractedData, error) {
	ext := extractTx(ctx, r.db)
	row := xEditExtractedDataRow{
		SourceFilePath:    draft.SourceFilePath,
		SourceContentHash: draft.SourceContentHash,
		SourceTool:        draft.SourceTool,
		TargetPluginName:  draft.TargetPluginName,
		TargetPluginType:  draft.TargetPluginType,
		RecordCount:       draft.RecordCount,
		ImportedAt:        draft.ImportedAt.UTC().Format(time.RFC3339),
	}
	q, args, err := sqlx.Named(insertXEditExtractedData, row)
	if err != nil {
		return XEditExtractedData{}, fmt.Errorf("create x_edit_extracted_data named: %w", err)
	}
	result, err := ext.ExecContext(ctx, q, args...)
	if err != nil {
		return XEditExtractedData{}, mapSQLError(err, "create x_edit_extracted_data")
	}
	id, err := result.LastInsertId()
	if err != nil {
		return XEditExtractedData{}, fmt.Errorf("create x_edit_extracted_data last insert id: %w", err)
	}
	return r.GetXEditExtractedDataByID(ctx, id)
}

// GetXEditExtractedDataByID は ID で XEditExtractedData を取得する。
func (r *SQLiteTranslationSourceRepository) GetXEditExtractedDataByID(
	ctx context.Context,
	id int64,
) (XEditExtractedData, error) {
	ext := extractTx(ctx, r.db)
	var row xEditExtractedDataRow
	if err := sqlx.GetContext(ctx, ext, &row, selectXEditExtractedDataByID, id); err != nil {
		return XEditExtractedData{}, mapSQLError(err, "get x_edit_extracted_data by id")
	}
	return row.toModel()
}

// ListXEditExtractedData returns all imported translation inputs ordered by most recent import first.
func (r *SQLiteTranslationSourceRepository) ListXEditExtractedData(
	ctx context.Context,
) ([]XEditExtractedData, error) {
	ext := extractTx(ctx, r.db)
	rows := []xEditExtractedDataRow{}
	if err := sqlx.SelectContext(ctx, ext, &rows, selectAllXEditExtractedData); err != nil {
		return nil, mapSQLError(err, "list x_edit_extracted_data")
	}
	results := make([]XEditExtractedData, 0, len(rows))
	for _, row := range rows {
		model, err := row.toModel()
		if err != nil {
			return nil, err
		}
		results = append(results, model)
	}
	return results, nil
}

// GetExistingTranslationJob returns the newest persisted translation job for Job Setup display.
func (r *SQLiteTranslationSourceRepository) GetExistingTranslationJob(
	ctx context.Context,
	xEditID int64,
) (TranslationJob, error) {
	ext := extractTx(ctx, r.db)
	var row existingTranslationJobRow
	if err := sqlx.GetContext(ctx, ext, &row, selectExistingTranslationJob, xEditID); err != nil {
		return TranslationJob{}, mapSQLError(err, "get existing translation_job")
	}
	return TranslationJob{ID: row.ID, State: row.State}, nil
}

// HasTranslationCacheByXEditID reports whether at least one cached translation record remains for the input.
func (r *SQLiteTranslationSourceRepository) HasTranslationCacheByXEditID(
	ctx context.Context,
	xEditID int64,
) (bool, error) {
	ext := extractTx(ctx, r.db)
	var exists bool
	if err := sqlx.GetContext(ctx, ext, &exists, selectTranslationCachePresenceByXEditID, xEditID); err != nil {
		return false, mapSQLError(err, "check translation cache by x_edit_id")
	}
	return exists, nil
}

// FindXEditExtractedDataBySourceContentHash は source_content_hash で XEditExtractedData を取得する。
func (r *SQLiteTranslationSourceRepository) FindXEditExtractedDataBySourceContentHash(
	ctx context.Context,
	sourceContentHash string,
) (XEditExtractedData, error) {
	ext := extractTx(ctx, r.db)
	var row xEditExtractedDataRow
	if err := sqlx.GetContext(ctx, ext, &row, selectXEditExtractedDataBySourceContentHash, sourceContentHash); err != nil {
		return XEditExtractedData{}, mapSQLError(err, "get x_edit_extracted_data by source_content_hash")
	}
	return row.toModel()
}

// UpdateXEditExtractedDataMetadata は既存入力メタデータを更新する。
func (r *SQLiteTranslationSourceRepository) UpdateXEditExtractedDataMetadata(
	ctx context.Context,
	id int64,
	draft XEditExtractedDataDraft,
) (XEditExtractedData, error) {
	ext := extractTx(ctx, r.db)
	_, err := ext.ExecContext(
		ctx,
		updateXEditExtractedDataMetadata,
		draft.SourceFilePath,
		draft.SourceContentHash,
		draft.SourceTool,
		draft.TargetPluginName,
		draft.TargetPluginType,
		draft.RecordCount,
		id,
	)
	if err != nil {
		return XEditExtractedData{}, mapSQLError(err, "update x_edit_extracted_data metadata")
	}
	return r.GetXEditExtractedDataByID(ctx, id)
}

// DeleteTranslationCacheByXEditID は入力に紐づく翻訳キャッシュを削除する。
func (r *SQLiteTranslationSourceRepository) DeleteTranslationCacheByXEditID(
	ctx context.Context,
	xEditID int64,
) error {
	ext := extractTx(ctx, r.db)
	statements := []struct {
		query string
		label string
		args  []any
	}{
		{query: deleteTranslationFieldRecordReferencesByXEditID, label: "delete translation_field_record_reference by x_edit_id", args: []any{xEditID}},
		{query: deletePersonaFieldEvidenceByXEditID, label: "delete persona_field_evidence by x_edit_id", args: []any{xEditID, xEditID}},
		{query: deleteXTranslatorOutputRowsByXEditID, label: "delete xtranslator_output_row by x_edit_id", args: []any{xEditID}},
		{query: deletePhaseRunTranslationFieldsByXEditID, label: "delete phase_run_translation_field by x_edit_id", args: []any{xEditID}},
		{query: deletePhaseRunPersonasByXEditID, label: "delete phase_run_persona by x_edit_id", args: []any{xEditID}},
		{query: deletePhaseRunDictionaryEntriesByXEditID, label: "delete phase_run_dictionary_entry by x_edit_id", args: []any{xEditID}},
		{query: deleteJobPhaseRunsByXEditID, label: "delete job_phase_run by x_edit_id", args: []any{xEditID}},
		{query: deleteTranslationArtifactsByXEditID, label: "delete translation_artifact by x_edit_id", args: []any{xEditID}},
		{query: deleteJobTranslationFieldsByXEditID, label: "delete job_translation_field by x_edit_id", args: []any{xEditID}},
		{query: deletePersonasByXEditID, label: "delete persona by x_edit_id", args: []any{xEditID}},
		{query: deleteDictionaryEntriesByXEditID, label: "delete dictionary_entry by x_edit_id", args: []any{xEditID}},
		{query: deleteTranslationJobsByXEditID, label: "delete translation_job by x_edit_id", args: []any{xEditID}},
		{query: deleteNpcRecordsByXEditID, label: "delete npc_record by x_edit_id", args: []any{xEditID}},
		{query: deleteTranslationFieldsByXEditID, label: "delete translation_field by x_edit_id", args: []any{xEditID}},
		{query: deleteTranslationRecordsByXEditID, label: "delete translation_record by x_edit_id", args: []any{xEditID}},
	}
	for _, statement := range statements {
		if _, err := ext.ExecContext(ctx, statement.query, statement.args...); err != nil {
			return mapSQLError(err, statement.label)
		}
	}
	return nil
}

// ---------------------------------------------------------------------------
// TranslationRecord
// ---------------------------------------------------------------------------

// CreateTranslationRecord は TranslationRecord レコードを作成する。
func (r *SQLiteTranslationSourceRepository) CreateTranslationRecord(
	ctx context.Context,
	draft TranslationRecordDraft,
) (TranslationRecord, error) {
	ext := extractTx(ctx, r.db)
	row := translationRecordRow{
		XEditExtractedDataID: draft.XEditExtractedDataID,
		FormID:               draft.FormID,
		EditorID:             draft.EditorID,
		RecordType:           draft.RecordType,
	}
	q, args, err := sqlx.Named(insertTranslationRecord, row)
	if err != nil {
		return TranslationRecord{}, fmt.Errorf("create translation_record named: %w", err)
	}
	result, err := ext.ExecContext(ctx, q, args...)
	if err != nil {
		return TranslationRecord{}, mapSQLError(err, "create translation_record")
	}
	id, err := result.LastInsertId()
	if err != nil {
		return TranslationRecord{}, fmt.Errorf("create translation_record last insert id: %w", err)
	}
	return r.GetTranslationRecordByID(ctx, id)
}

// GetTranslationRecordByID は ID で TranslationRecord を取得する。
func (r *SQLiteTranslationSourceRepository) GetTranslationRecordByID(
	ctx context.Context,
	id int64,
) (TranslationRecord, error) {
	ext := extractTx(ctx, r.db)
	var row translationRecordRow
	if err := sqlx.GetContext(ctx, ext, &row, selectTranslationRecordByID, id); err != nil {
		return TranslationRecord{}, mapSQLError(err, "get translation_record by id")
	}
	return row.toModel(), nil
}

// ListTranslationRecordsByXEditID は XEditID に紐づく TranslationRecord 一覧を返す。
func (r *SQLiteTranslationSourceRepository) ListTranslationRecordsByXEditID(
	ctx context.Context,
	xEditID int64,
) ([]TranslationRecord, error) {
	ext := extractTx(ctx, r.db)
	var rows []translationRecordRow
	if err := sqlx.SelectContext(ctx, ext, &rows, selectTranslationRecordsByXEditID, xEditID); err != nil {
		return nil, mapSQLError(err, "list translation_records by x_edit_id")
	}
	result := make([]TranslationRecord, len(rows))
	for i, row := range rows {
		result[i] = row.toModel()
	}
	return result, nil
}

// ---------------------------------------------------------------------------
// NpcProfile
// ---------------------------------------------------------------------------

// UpsertNpcProfile は NpcProfile を upsert する。
func (r *SQLiteTranslationSourceRepository) UpsertNpcProfile(
	ctx context.Context,
	draft NpcProfileDraft,
) (NpcProfile, error) {
	ext := extractTx(ctx, r.db)
	now := time.Now().UTC().Format(time.RFC3339)
	row := npcProfileRow{
		TargetPluginName: draft.TargetPluginName,
		FormID:           draft.FormID,
		RecordType:       draft.RecordType,
		EditorID:         draft.EditorID,
		DisplayName:      draft.DisplayName,
		CreatedAt:        now,
		UpdatedAt:        now,
	}
	q, args, err := sqlx.Named(upsertNpcProfile, row)
	if err != nil {
		return NpcProfile{}, fmt.Errorf("upsert npc_profile named: %w", err)
	}
	if _, err := ext.ExecContext(ctx, q, args...); err != nil {
		return NpcProfile{}, mapSQLError(err, "upsert npc_profile")
	}
	var fetched npcProfileRow
	if err := sqlx.GetContext(ctx, ext, &fetched, selectNpcProfileByUnique,
		draft.TargetPluginName, draft.FormID, draft.RecordType,
	); err != nil {
		return NpcProfile{}, mapSQLError(err, "upsert npc_profile fetch")
	}
	return fetched.toModel()
}

// GetNpcProfileByID は ID で NpcProfile を取得する。
func (r *SQLiteTranslationSourceRepository) GetNpcProfileByID(
	ctx context.Context,
	id int64,
) (NpcProfile, error) {
	ext := extractTx(ctx, r.db)
	var row npcProfileRow
	if err := sqlx.GetContext(ctx, ext, &row, selectNpcProfileByID, id); err != nil {
		return NpcProfile{}, mapSQLError(err, "get npc_profile by id")
	}
	return row.toModel()
}

// ---------------------------------------------------------------------------
// NpcRecord
// ---------------------------------------------------------------------------

// CreateNpcRecord は NpcRecord レコードを作成する。
func (r *SQLiteTranslationSourceRepository) CreateNpcRecord(
	ctx context.Context,
	draft NpcRecordDraft,
) (NpcRecord, error) {
	ext := extractTx(ctx, r.db)
	row := npcRecordRow(draft)
	q, args, err := sqlx.Named(insertNpcRecord, row)
	if err != nil {
		return NpcRecord{}, fmt.Errorf("create npc_record named: %w", err)
	}
	if _, err := ext.ExecContext(ctx, q, args...); err != nil {
		return NpcRecord{}, mapSQLError(err, "create npc_record")
	}
	return row.toModel(), nil
}

// GetNpcRecordByTranslationRecordID は TranslationRecordID で NpcRecord を取得する。
func (r *SQLiteTranslationSourceRepository) GetNpcRecordByTranslationRecordID(
	ctx context.Context,
	translationRecordID int64,
) (NpcRecord, error) {
	ext := extractTx(ctx, r.db)
	var row npcRecordRow
	if err := sqlx.GetContext(ctx, ext, &row, selectNpcRecordByTranslationRecordID, translationRecordID); err != nil {
		return NpcRecord{}, mapSQLError(err, "get npc_record by translation_record_id")
	}
	return row.toModel(), nil
}

// ---------------------------------------------------------------------------
// TranslationField
// ---------------------------------------------------------------------------

// CreateTranslationField は TranslationField レコードを作成する。
func (r *SQLiteTranslationSourceRepository) CreateTranslationField(
	ctx context.Context,
	draft TranslationFieldDraft,
) (TranslationField, error) {
	ext := extractTx(ctx, r.db)
	row := translationFieldRow{
		TranslationRecordID:          draft.TranslationRecordID,
		TranslationFieldDefinitionID: draft.TranslationFieldDefinitionID,
		SubrecordType:                draft.SubrecordType,
		SourceText:                   draft.SourceText,
		FieldOrder:                   draft.FieldOrder,
		PreviousTranslationFieldID:   draft.PreviousTranslationFieldID,
		NextTranslationFieldID:       draft.NextTranslationFieldID,
	}
	q, args, err := sqlx.Named(insertTranslationField, row)
	if err != nil {
		return TranslationField{}, fmt.Errorf("create translation_field named: %w", err)
	}
	result, err := ext.ExecContext(ctx, q, args...)
	if err != nil {
		return TranslationField{}, mapSQLError(err, "create translation_field")
	}
	id, err := result.LastInsertId()
	if err != nil {
		return TranslationField{}, fmt.Errorf("create translation_field last insert id: %w", err)
	}
	return r.GetTranslationFieldByID(ctx, id)
}

// GetTranslationFieldByID は ID で TranslationField を取得する。
func (r *SQLiteTranslationSourceRepository) GetTranslationFieldByID(
	ctx context.Context,
	id int64,
) (TranslationField, error) {
	ext := extractTx(ctx, r.db)
	var row translationFieldRow
	if err := sqlx.GetContext(ctx, ext, &row, selectTranslationFieldByID, id); err != nil {
		return TranslationField{}, mapSQLError(err, "get translation_field by id")
	}
	return row.toModel(), nil
}

// ListTranslationFieldsByTranslationRecordID は TranslationRecordID に紐づく TranslationField 一覧を返す。
func (r *SQLiteTranslationSourceRepository) ListTranslationFieldsByTranslationRecordID(
	ctx context.Context,
	translationRecordID int64,
) ([]TranslationField, error) {
	ext := extractTx(ctx, r.db)
	var rows []translationFieldRow
	if err := sqlx.SelectContext(ctx, ext, &rows, selectTranslationFieldsByTranslationRecordID, translationRecordID); err != nil {
		return nil, mapSQLError(err, "list translation_fields by translation_record_id")
	}
	result := make([]TranslationField, len(rows))
	for i, row := range rows {
		result[i] = row.toModel()
	}
	return result, nil
}

// ---------------------------------------------------------------------------
// TranslationFieldRecordReference
// ---------------------------------------------------------------------------

// CreateTranslationFieldRecordReference は TranslationFieldRecordReference レコードを作成する。
func (r *SQLiteTranslationSourceRepository) CreateTranslationFieldRecordReference(
	ctx context.Context,
	draft TranslationFieldRecordReferenceDraft,
) (TranslationFieldRecordReference, error) {
	ext := extractTx(ctx, r.db)
	row := translationFieldRecordReferenceRow{
		TranslationFieldID:            draft.TranslationFieldID,
		ReferencedTranslationRecordID: draft.ReferencedTranslationRecordID,
		ReferenceRole:                 draft.ReferenceRole,
	}
	q, args, err := sqlx.Named(insertTranslationFieldRecordReference, row)
	if err != nil {
		return TranslationFieldRecordReference{}, fmt.Errorf("create translation_field_record_reference named: %w", err)
	}
	result, err := ext.ExecContext(ctx, q, args...)
	if err != nil {
		return TranslationFieldRecordReference{}, mapSQLError(err, "create translation_field_record_reference")
	}
	id, err := result.LastInsertId()
	if err != nil {
		return TranslationFieldRecordReference{}, fmt.Errorf("create translation_field_record_reference last insert id: %w", err)
	}
	return TranslationFieldRecordReference{
		ID:                            id,
		TranslationFieldID:            draft.TranslationFieldID,
		ReferencedTranslationRecordID: draft.ReferencedTranslationRecordID,
		ReferenceRole:                 draft.ReferenceRole,
	}, nil
}

// ListTranslationFieldRecordReferencesByFieldID は FieldID に紐づく TranslationFieldRecordReference 一覧を返す。
func (r *SQLiteTranslationSourceRepository) ListTranslationFieldRecordReferencesByFieldID(
	ctx context.Context,
	fieldID int64,
) ([]TranslationFieldRecordReference, error) {
	ext := extractTx(ctx, r.db)
	var rows []translationFieldRecordReferenceRow
	if err := sqlx.SelectContext(ctx, ext, &rows, selectTranslationFieldRecordReferencesByFieldID, fieldID); err != nil {
		return nil, mapSQLError(err, "list translation_field_record_references by field_id")
	}
	result := make([]TranslationFieldRecordReference, len(rows))
	for i, row := range rows {
		result[i] = row.toModel()
	}
	return result, nil
}
