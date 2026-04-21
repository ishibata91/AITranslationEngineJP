package repository

import (
	"context"
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
