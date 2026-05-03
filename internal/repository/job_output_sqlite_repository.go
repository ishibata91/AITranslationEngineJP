package repository

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
)

// SQLiteJobOutputRepository は JobOutputRepository の SQLite 実装。
type SQLiteJobOutputRepository struct {
	db *sqlx.DB
}

// NewSQLiteJobOutputRepository は JobOutputRepository を返す。
func NewSQLiteJobOutputRepository(db *sqlx.DB) JobOutputRepository {
	return &SQLiteJobOutputRepository{db: db}
}

// NewSQLiteJobOutputArtifactRepository は JobOutputArtifactRepository を返す。
func NewSQLiteJobOutputArtifactRepository(db *sqlx.DB) JobOutputArtifactRepository {
	return &SQLiteJobOutputRepository{db: db}
}

// ---------------------------------------------------------------------------
// 内部 row 型
// ---------------------------------------------------------------------------

type jobTranslationFieldRow struct {
	ID                 int64  `db:"id"`
	TranslationJobID   int64  `db:"translation_job_id"`
	TranslationFieldID int64  `db:"translation_field_id"`
	AppliedPersonaID   *int64 `db:"applied_persona_id"`
	TranslatedText     string `db:"translated_text"`
	OutputStatus       string `db:"output_status"`
	RetryCount         int    `db:"retry_count"`
	UpdatedAt          string `db:"updated_at"`
}

func (r jobTranslationFieldRow) toModel() (JobTranslationField, error) {
	updatedAt, err := time.Parse(time.RFC3339, r.UpdatedAt)
	if err != nil {
		return JobTranslationField{}, fmt.Errorf("parse updated_at: %w", err)
	}
	return JobTranslationField{
		ID:                 r.ID,
		TranslationJobID:   r.TranslationJobID,
		TranslationFieldID: r.TranslationFieldID,
		AppliedPersonaID:   r.AppliedPersonaID,
		TranslatedText:     r.TranslatedText,
		OutputStatus:       r.OutputStatus,
		RetryCount:         r.RetryCount,
		UpdatedAt:          updatedAt,
	}, nil
}

type translationArtifactRow struct {
	ID               int64   `db:"id"`
	TranslationJobID int64   `db:"translation_job_id"`
	ArtifactFormat   string  `db:"artifact_format"`
	TargetGame       string  `db:"target_game"`
	FilePath         string  `db:"file_path"`
	Status           string  `db:"status"`
	GeneratedAt      *string `db:"generated_at"`
}

func (r translationArtifactRow) toModel() (TranslationArtifact, error) {
	var generatedAt *time.Time
	if r.GeneratedAt != nil && *r.GeneratedAt != "" {
		parsed, err := time.Parse(time.RFC3339, *r.GeneratedAt)
		if err != nil {
			return TranslationArtifact{}, fmt.Errorf("parse generated_at: %w", err)
		}
		generatedAt = &parsed
	}
	return TranslationArtifact{
		ID:               r.ID,
		TranslationJobID: r.TranslationJobID,
		ArtifactFormat:   r.ArtifactFormat,
		TargetGame:       r.TargetGame,
		FilePath:         r.FilePath,
		Status:           r.Status,
		GeneratedAt:      generatedAt,
	}, nil
}

type xTranslatorOutputRowRow struct {
	ID                    int64  `db:"id"`
	TranslationArtifactID int64  `db:"translation_artifact_id"`
	JobTranslationFieldID int64  `db:"job_translation_field_id"`
	EDID                  string `db:"edid"`
	REC                   string `db:"rec"`
	FIELD                 string `db:"field"`
	FORMID                string `db:"formid"`
	Source                string `db:"source"`
	Dest                  string `db:"dest"`
	Status                int    `db:"status"`
}

func (r xTranslatorOutputRowRow) toModel() XTranslatorOutputRow {
	return XTranslatorOutputRow(r)
}

// ---------------------------------------------------------------------------
// SQL 定数
// ---------------------------------------------------------------------------

const (
	insertJobTranslationField = `
INSERT INTO JOB_TRANSLATION_FIELD
  (translation_job_id, translation_field_id, applied_persona_id, translated_text,
   output_status, retry_count, updated_at)
VALUES
  (:translation_job_id, :translation_field_id, :applied_persona_id, :translated_text,
   :output_status, :retry_count, :updated_at)`

	selectJobTranslationFieldByID = `
SELECT id, translation_job_id, translation_field_id, applied_persona_id,
       translated_text, output_status, retry_count, updated_at
FROM JOB_TRANSLATION_FIELD WHERE id = ?`

	updateJobTranslationField = `
UPDATE JOB_TRANSLATION_FIELD SET
  applied_persona_id = :applied_persona_id,
  translated_text    = :translated_text,
  output_status      = :output_status,
  retry_count        = :retry_count,
  updated_at         = :updated_at
WHERE id = :id`

	selectJobTranslationFieldsByJobID = `
SELECT id, translation_job_id, translation_field_id, applied_persona_id,
       translated_text, output_status, retry_count, updated_at
FROM JOB_TRANSLATION_FIELD WHERE translation_job_id = ?`

	selectTranslationArtifactByID = `
SELECT id, translation_job_id, artifact_format, target_game, file_path, status, generated_at
FROM TRANSLATION_ARTIFACT WHERE id = ?`

	selectTranslationArtifactByJobID = `
SELECT id, translation_job_id, artifact_format, target_game, file_path, status, generated_at
FROM TRANSLATION_ARTIFACT WHERE translation_job_id = ?`

	insertTranslationArtifact = `
INSERT INTO TRANSLATION_ARTIFACT
  (translation_job_id, artifact_format, target_game, file_path, status, generated_at)
VALUES
  (:translation_job_id, :artifact_format, :target_game, :file_path, :status, :generated_at)`

	updateTranslationArtifactByJobID = `
UPDATE TRANSLATION_ARTIFACT SET
  artifact_format = :artifact_format,
  target_game     = :target_game,
  file_path       = :file_path,
  status          = :status,
  generated_at    = :generated_at
WHERE translation_job_id = :translation_job_id`

	deleteXTranslatorOutputRowsByArtifactID = `
DELETE FROM XTRANSLATOR_OUTPUT_ROW WHERE translation_artifact_id = ?`

	insertXTranslatorOutputRow = `
INSERT INTO XTRANSLATOR_OUTPUT_ROW
  (translation_artifact_id, job_translation_field_id, edid, rec, field, formid, source, dest, status)
VALUES
  (:translation_artifact_id, :job_translation_field_id, :edid, :rec, :field, :formid, :source, :dest, :status)`

	selectXTranslatorOutputRowsByArtifactID = `
SELECT id, translation_artifact_id, job_translation_field_id, edid, rec, field, formid, source, dest, status
FROM XTRANSLATOR_OUTPUT_ROW
WHERE translation_artifact_id = ?
ORDER BY job_translation_field_id ASC`

	countXTranslatorOutputRowsByArtifactID = `
SELECT COUNT(1)
FROM XTRANSLATOR_OUTPUT_ROW
WHERE translation_artifact_id = ?`
)

// ---------------------------------------------------------------------------
// JobTranslationField
// ---------------------------------------------------------------------------

// CreateJobTranslationField は JobTranslationField レコードを作成する。
func (r *SQLiteJobOutputRepository) CreateJobTranslationField(
	ctx context.Context,
	draft JobTranslationFieldDraft,
) (JobTranslationField, error) {
	ext := extractTx(ctx, r.db)
	now := time.Now().UTC().Format(time.RFC3339)
	row := jobTranslationFieldRow{
		TranslationJobID:   draft.TranslationJobID,
		TranslationFieldID: draft.TranslationFieldID,
		AppliedPersonaID:   draft.AppliedPersonaID,
		TranslatedText:     draft.TranslatedText,
		OutputStatus:       draft.OutputStatus,
		RetryCount:         draft.RetryCount,
		UpdatedAt:          now,
	}
	q, args, err := sqlx.Named(insertJobTranslationField, row)
	if err != nil {
		return JobTranslationField{}, fmt.Errorf("create job_translation_field named: %w", err)
	}
	result, err := ext.ExecContext(ctx, q, args...)
	if err != nil {
		return JobTranslationField{}, mapFoundationSQLError(err, "create job_translation_field")
	}
	id, err := result.LastInsertId()
	if err != nil {
		return JobTranslationField{}, fmt.Errorf("create job_translation_field last insert id: %w", err)
	}
	return r.GetJobTranslationFieldByID(ctx, id)
}

// GetJobTranslationFieldByID は ID で JobTranslationField を取得する。
func (r *SQLiteJobOutputRepository) GetJobTranslationFieldByID(
	ctx context.Context,
	id int64,
) (JobTranslationField, error) {
	ext := extractTx(ctx, r.db)
	var row jobTranslationFieldRow
	if err := sqlx.GetContext(ctx, ext, &row, selectJobTranslationFieldByID, id); err != nil {
		return JobTranslationField{}, mapSQLError(err, "get job_translation_field by id")
	}
	return row.toModel()
}

// UpdateJobTranslationField は JobTranslationField を更新する。
func (r *SQLiteJobOutputRepository) UpdateJobTranslationField(
	ctx context.Context,
	id int64,
	draft JobTranslationFieldUpdateDraft,
) (JobTranslationField, error) {
	ext := extractTx(ctx, r.db)
	now := time.Now().UTC().Format(time.RFC3339)
	args := map[string]interface{}{
		"id":                 id,
		"applied_persona_id": draft.AppliedPersonaID,
		"translated_text":    draft.TranslatedText,
		"output_status":      draft.OutputStatus,
		"retry_count":        draft.RetryCount,
		"updated_at":         now,
	}
	q, qArgs, err := sqlx.Named(updateJobTranslationField, args)
	if err != nil {
		return JobTranslationField{}, fmt.Errorf("update job_translation_field named: %w", err)
	}
	if _, err := ext.ExecContext(ctx, q, qArgs...); err != nil {
		return JobTranslationField{}, mapFoundationSQLError(err, "update job_translation_field")
	}
	return r.GetJobTranslationFieldByID(ctx, id)
}

// ListJobTranslationFieldsByJobID は JobID に紐づく JobTranslationField 一覧を返す。
func (r *SQLiteJobOutputRepository) ListJobTranslationFieldsByJobID(
	ctx context.Context,
	jobID int64,
) ([]JobTranslationField, error) {
	ext := extractTx(ctx, r.db)
	var rows []jobTranslationFieldRow
	if err := sqlx.SelectContext(ctx, ext, &rows, selectJobTranslationFieldsByJobID, jobID); err != nil {
		return nil, mapSQLError(err, "list job_translation_fields by job_id")
	}
	result := make([]JobTranslationField, len(rows))
	for i, row := range rows {
		m, err := row.toModel()
		if err != nil {
			return nil, err
		}
		result[i] = m
	}
	return result, nil
}

// GetTranslationArtifactByID は ID で TranslationArtifact を取得する。
func (r *SQLiteJobOutputRepository) GetTranslationArtifactByID(
	ctx context.Context,
	id int64,
) (TranslationArtifact, error) {
	ext := extractTx(ctx, r.db)
	var row translationArtifactRow
	if err := sqlx.GetContext(ctx, ext, &row, selectTranslationArtifactByID, id); err != nil {
		return TranslationArtifact{}, mapSQLError(err, "get translation_artifact by id")
	}
	return row.toModel()
}

// GetTranslationArtifactByJobID は JobID で TranslationArtifact を取得する。
func (r *SQLiteJobOutputRepository) GetTranslationArtifactByJobID(
	ctx context.Context,
	jobID int64,
) (TranslationArtifact, error) {
	ext := extractTx(ctx, r.db)
	var row translationArtifactRow
	if err := sqlx.GetContext(ctx, ext, &row, selectTranslationArtifactByJobID, jobID); err != nil {
		return TranslationArtifact{}, mapSQLError(err, "get translation_artifact by job_id")
	}
	return row.toModel()
}

// UpsertTranslationArtifact は translation_job_id 単位で artifact を作成または更新する。
func (r *SQLiteJobOutputRepository) UpsertTranslationArtifact(
	ctx context.Context,
	draft TranslationArtifactDraft,
) (TranslationArtifact, error) {
	ext := extractTx(ctx, r.db)
	var generatedAt *string
	if draft.GeneratedAt != nil {
		value := draft.GeneratedAt.UTC().Format(time.RFC3339)
		generatedAt = &value
	}
	row := translationArtifactRow{
		TranslationJobID: draft.TranslationJobID,
		ArtifactFormat:   draft.ArtifactFormat,
		TargetGame:       draft.TargetGame,
		FilePath:         draft.FilePath,
		Status:           draft.Status,
		GeneratedAt:      generatedAt,
	}

	if _, err := r.GetTranslationArtifactByJobID(ctx, draft.TranslationJobID); err != nil {
		if !errors.Is(err, ErrNotFound) {
			return TranslationArtifact{}, err
		}
		q, args, namedErr := sqlx.Named(insertTranslationArtifact, row)
		if namedErr != nil {
			return TranslationArtifact{}, fmt.Errorf("create translation_artifact named: %w", namedErr)
		}
		if _, execErr := ext.ExecContext(ctx, q, args...); execErr != nil {
			return TranslationArtifact{}, mapFoundationSQLError(execErr, "create translation_artifact")
		}
	} else {
		q, args, namedErr := sqlx.Named(updateTranslationArtifactByJobID, row)
		if namedErr != nil {
			return TranslationArtifact{}, fmt.Errorf("update translation_artifact named: %w", namedErr)
		}
		if _, execErr := ext.ExecContext(ctx, q, args...); execErr != nil {
			return TranslationArtifact{}, mapFoundationSQLError(execErr, "update translation_artifact")
		}
	}

	return r.GetTranslationArtifactByJobID(ctx, draft.TranslationJobID)
}

// ListXTranslatorOutputRowsByArtifactID は artifact 配下 row 一覧を返す。
func (r *SQLiteJobOutputRepository) ListXTranslatorOutputRowsByArtifactID(
	ctx context.Context,
	translationArtifactID int64,
) ([]XTranslatorOutputRow, error) {
	ext := extractTx(ctx, r.db)
	var rows []xTranslatorOutputRowRow
	if err := sqlx.SelectContext(ctx, ext, &rows, selectXTranslatorOutputRowsByArtifactID, translationArtifactID); err != nil {
		return nil, mapSQLError(err, "list xtranslator_output_row by artifact_id")
	}
	result := make([]XTranslatorOutputRow, len(rows))
	for i, row := range rows {
		result[i] = row.toModel()
	}
	return result, nil
}

// ReplaceXTranslatorOutputRows は artifact 配下の row を全置換する。
func (r *SQLiteJobOutputRepository) ReplaceXTranslatorOutputRows(
	ctx context.Context,
	translationArtifactID int64,
	drafts []XTranslatorOutputRowDraft,
) ([]XTranslatorOutputRow, error) {
	ext := extractTx(ctx, r.db)
	if _, err := ext.ExecContext(ctx, deleteXTranslatorOutputRowsByArtifactID, translationArtifactID); err != nil {
		return nil, mapFoundationSQLError(err, "delete xtranslator_output_row by artifact_id")
	}

	for _, draft := range drafts {
		row := map[string]any{
			"translation_artifact_id":  translationArtifactID,
			"job_translation_field_id": draft.JobTranslationFieldID,
			"edid":                     draft.EDID,
			"rec":                      draft.REC,
			"field":                    draft.FIELD,
			"formid":                   draft.FORMID,
			"source":                   draft.Source,
			"dest":                     draft.Dest,
			"status":                   draft.Status,
		}
		q, args, namedErr := sqlx.Named(insertXTranslatorOutputRow, row)
		if namedErr != nil {
			return nil, fmt.Errorf("create xtranslator_output_row named: %w", namedErr)
		}
		if _, execErr := ext.ExecContext(ctx, q, args...); execErr != nil {
			return nil, mapFoundationSQLError(execErr, "create xtranslator_output_row")
		}
	}

	return r.ListXTranslatorOutputRowsByArtifactID(ctx, translationArtifactID)
}

// CountXTranslatorOutputRowsByArtifactID は artifact 配下 row 数を返す。
func (r *SQLiteJobOutputRepository) CountXTranslatorOutputRowsByArtifactID(
	ctx context.Context,
	translationArtifactID int64,
) (int, error) {
	ext := extractTx(ctx, r.db)
	var count int
	if err := sqlx.GetContext(ctx, ext, &count, countXTranslatorOutputRowsByArtifactID, translationArtifactID); err != nil {
		return 0, mapSQLError(err, "count xtranslator_output_row by artifact_id")
	}
	return count, nil
}
