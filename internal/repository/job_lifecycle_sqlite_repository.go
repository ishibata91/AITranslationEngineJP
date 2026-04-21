package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
)

// SQLiteJobLifecycleRepository は JobLifecycleRepository の SQLite 実装。
type SQLiteJobLifecycleRepository struct {
	db *sqlx.DB
}

// NewSQLiteJobLifecycleRepository は JobLifecycleRepository を返す。
func NewSQLiteJobLifecycleRepository(db *sqlx.DB) JobLifecycleRepository {
	return &SQLiteJobLifecycleRepository{db: db}
}

// ---------------------------------------------------------------------------
// 内部 row 型
// ---------------------------------------------------------------------------

type translationJobRow struct {
	ID                   int64   `db:"id"`
	XEditExtractedDataID int64   `db:"x_edit_extracted_data_id"`
	JobName              string  `db:"job_name"`
	State                string  `db:"state"`
	ProgressPercent      int     `db:"progress_percent"`
	CreatedAt            string  `db:"created_at"`
	StartedAt            *string `db:"started_at"`
	FinishedAt           *string `db:"finished_at"`
}

func (r translationJobRow) toModel() (TranslationJob, error) {
	createdAt, err := time.Parse(time.RFC3339, r.CreatedAt)
	if err != nil {
		return TranslationJob{}, fmt.Errorf("parse created_at: %w", err)
	}
	var startedAt *time.Time
	if r.StartedAt != nil {
		t, err := time.Parse(time.RFC3339, *r.StartedAt)
		if err != nil {
			return TranslationJob{}, fmt.Errorf("parse started_at: %w", err)
		}
		startedAt = &t
	}
	var finishedAt *time.Time
	if r.FinishedAt != nil {
		t, err := time.Parse(time.RFC3339, *r.FinishedAt)
		if err != nil {
			return TranslationJob{}, fmt.Errorf("parse finished_at: %w", err)
		}
		finishedAt = &t
	}
	return TranslationJob{
		ID:                   r.ID,
		XEditExtractedDataID: r.XEditExtractedDataID,
		JobName:              r.JobName,
		State:                r.State,
		ProgressPercent:      r.ProgressPercent,
		CreatedAt:            createdAt,
		StartedAt:            startedAt,
		FinishedAt:           finishedAt,
	}, nil
}

type jobPhaseRunRow struct {
	ID                  int64   `db:"id"`
	TranslationJobID    int64   `db:"translation_job_id"`
	PhaseType           string  `db:"phase_type"`
	State               string  `db:"state"`
	ExecutionOrder      int     `db:"execution_order"`
	ProgressPercent     int     `db:"progress_percent"`
	AIProvider          string  `db:"ai_provider"`
	ModelName           string  `db:"model_name"`
	ExecutionMode       string  `db:"execution_mode"`
	CredentialRef       string  `db:"credential_ref"`
	InstructionKind     string  `db:"instruction_kind"`
	LatestExternalRunID string  `db:"latest_external_run_id"`
	LatestError         string  `db:"latest_error"`
	StartedAt           *string `db:"started_at"`
	FinishedAt          *string `db:"finished_at"`
}

func (r jobPhaseRunRow) toModel() (JobPhaseRun, error) {
	var startedAt *time.Time
	if r.StartedAt != nil {
		t, err := time.Parse(time.RFC3339, *r.StartedAt)
		if err != nil {
			return JobPhaseRun{}, fmt.Errorf("parse started_at: %w", err)
		}
		startedAt = &t
	}
	var finishedAt *time.Time
	if r.FinishedAt != nil {
		t, err := time.Parse(time.RFC3339, *r.FinishedAt)
		if err != nil {
			return JobPhaseRun{}, fmt.Errorf("parse finished_at: %w", err)
		}
		finishedAt = &t
	}
	return JobPhaseRun{
		ID:                  r.ID,
		TranslationJobID:    r.TranslationJobID,
		PhaseType:           r.PhaseType,
		State:               r.State,
		ExecutionOrder:      r.ExecutionOrder,
		ProgressPercent:     r.ProgressPercent,
		AIProvider:          r.AIProvider,
		ModelName:           r.ModelName,
		ExecutionMode:       r.ExecutionMode,
		CredentialRef:       r.CredentialRef,
		InstructionKind:     r.InstructionKind,
		LatestExternalRunID: r.LatestExternalRunID,
		LatestError:         r.LatestError,
		StartedAt:           startedAt,
		FinishedAt:          finishedAt,
	}, nil
}

type phaseRunTranslationFieldRow struct {
	ID                    int64  `db:"id"`
	PhaseRunID            int64  `db:"phase_run_id"`
	JobTranslationFieldID int64  `db:"job_translation_field_id"`
	Role                  string `db:"role"`
}

type phaseRunPersonaRow struct {
	ID         int64  `db:"id"`
	PhaseRunID int64  `db:"phase_run_id"`
	PersonaID  int64  `db:"persona_id"`
	Role       string `db:"role"`
}

type phaseRunDictionaryEntryRow struct {
	ID                int64  `db:"id"`
	PhaseRunID        int64  `db:"phase_run_id"`
	DictionaryEntryID int64  `db:"dictionary_entry_id"`
	Role              string `db:"role"`
}

// ---------------------------------------------------------------------------
// SQL 定数
// ---------------------------------------------------------------------------

const (
	insertTranslationJob = `
INSERT INTO TRANSLATION_JOB
  (x_edit_extracted_data_id, job_name, state, progress_percent, created_at, started_at, finished_at)
VALUES
  (:x_edit_extracted_data_id, :job_name, :state, :progress_percent, :created_at, :started_at, :finished_at)`

	selectTranslationJobByID = `
SELECT id, x_edit_extracted_data_id, job_name, state, progress_percent, created_at, started_at, finished_at
FROM TRANSLATION_JOB WHERE id = ?`

	updateTranslationJob = `
UPDATE TRANSLATION_JOB SET
  job_name         = :job_name,
  state            = :state,
  progress_percent = :progress_percent,
  started_at       = :started_at,
  finished_at      = :finished_at
WHERE id = :id`

	insertJobPhaseRun = `
INSERT INTO JOB_PHASE_RUN
  (translation_job_id, phase_type, state, execution_order, progress_percent,
   ai_provider, model_name, execution_mode, credential_ref, instruction_kind,
   latest_external_run_id, latest_error, started_at, finished_at)
VALUES
  (:translation_job_id, :phase_type, :state, :execution_order, :progress_percent,
   :ai_provider, :model_name, :execution_mode, :credential_ref, :instruction_kind,
   :latest_external_run_id, :latest_error, :started_at, :finished_at)`

	selectJobPhaseRunByID = `
SELECT id, translation_job_id, phase_type, state, execution_order, progress_percent,
       ai_provider, model_name, execution_mode, credential_ref, instruction_kind,
       latest_external_run_id, latest_error, started_at, finished_at
FROM JOB_PHASE_RUN WHERE id = ?`

	updateJobPhaseRun = `
UPDATE JOB_PHASE_RUN SET
  state                 = :state,
  progress_percent      = :progress_percent,
  latest_external_run_id = :latest_external_run_id,
  latest_error          = :latest_error,
  started_at            = :started_at,
  finished_at           = :finished_at
WHERE id = :id`

	selectJobPhaseRunsByJobID = `
SELECT id, translation_job_id, phase_type, state, execution_order, progress_percent,
       ai_provider, model_name, execution_mode, credential_ref, instruction_kind,
       latest_external_run_id, latest_error, started_at, finished_at
FROM JOB_PHASE_RUN WHERE translation_job_id = ? ORDER BY execution_order ASC`

	insertPhaseRunTranslationField = `
INSERT INTO PHASE_RUN_TRANSLATION_FIELD
  (phase_run_id, job_translation_field_id, role)
VALUES
  (:phase_run_id, :job_translation_field_id, :role)`

	insertPhaseRunPersona = `
INSERT INTO PHASE_RUN_PERSONA
  (phase_run_id, persona_id, role)
VALUES
  (:phase_run_id, :persona_id, :role)`

	insertPhaseRunDictionaryEntry = `
INSERT INTO PHASE_RUN_DICTIONARY_ENTRY
  (phase_run_id, dictionary_entry_id, role)
VALUES
  (:phase_run_id, :dictionary_entry_id, :role)`
)

// ---------------------------------------------------------------------------
// TranslationJob
// ---------------------------------------------------------------------------

// CreateTranslationJob は TranslationJob レコードを作成する。
func (r *SQLiteJobLifecycleRepository) CreateTranslationJob(
	ctx context.Context,
	draft TranslationJobDraft,
) (TranslationJob, error) {
	ext := extractTx(ctx, r.db)
	now := time.Now().UTC().Format(time.RFC3339)
	row := translationJobRow{
		XEditExtractedDataID: draft.XEditExtractedDataID,
		JobName:              draft.JobName,
		State:                draft.State,
		ProgressPercent:      draft.ProgressPercent,
		CreatedAt:            now,
		StartedAt:            nil,
		FinishedAt:           nil,
	}
	q, args, err := sqlx.Named(insertTranslationJob, row)
	if err != nil {
		return TranslationJob{}, fmt.Errorf("create translation_job named: %w", err)
	}
	result, err := ext.ExecContext(ctx, q, args...)
	if err != nil {
		return TranslationJob{}, mapFoundationSQLError(err, "create translation_job")
	}
	id, err := result.LastInsertId()
	if err != nil {
		return TranslationJob{}, fmt.Errorf("create translation_job last insert id: %w", err)
	}
	return r.GetTranslationJobByID(ctx, id)
}

// GetTranslationJobByID は ID で TranslationJob を取得する。
func (r *SQLiteJobLifecycleRepository) GetTranslationJobByID(
	ctx context.Context,
	id int64,
) (TranslationJob, error) {
	ext := extractTx(ctx, r.db)
	var row translationJobRow
	if err := sqlx.GetContext(ctx, ext, &row, selectTranslationJobByID, id); err != nil {
		return TranslationJob{}, mapSQLError(err, "get translation_job by id")
	}
	return row.toModel()
}

// UpdateTranslationJob は TranslationJob を更新する。
func (r *SQLiteJobLifecycleRepository) UpdateTranslationJob(
	ctx context.Context,
	id int64,
	draft TranslationJobUpdateDraft,
) (TranslationJob, error) {
	ext := extractTx(ctx, r.db)
	var startedAt *string
	if draft.StartedAt != nil {
		s := draft.StartedAt.UTC().Format(time.RFC3339)
		startedAt = &s
	}
	var finishedAt *string
	if draft.FinishedAt != nil {
		s := draft.FinishedAt.UTC().Format(time.RFC3339)
		finishedAt = &s
	}
	args := map[string]interface{}{
		"id":               id,
		"job_name":         draft.JobName,
		"state":            draft.State,
		"progress_percent": draft.ProgressPercent,
		"started_at":       startedAt,
		"finished_at":      finishedAt,
	}
	q, qArgs, err := sqlx.Named(updateTranslationJob, args)
	if err != nil {
		return TranslationJob{}, fmt.Errorf("update translation_job named: %w", err)
	}
	if _, err := ext.ExecContext(ctx, q, qArgs...); err != nil {
		return TranslationJob{}, mapFoundationSQLError(err, "update translation_job")
	}
	return r.GetTranslationJobByID(ctx, id)
}

// ---------------------------------------------------------------------------
// JobPhaseRun
// ---------------------------------------------------------------------------

// CreateJobPhaseRun は JobPhaseRun レコードを作成する。
func (r *SQLiteJobLifecycleRepository) CreateJobPhaseRun(
	ctx context.Context,
	draft JobPhaseRunDraft,
) (JobPhaseRun, error) {
	ext := extractTx(ctx, r.db)
	row := jobPhaseRunRow{
		TranslationJobID:    draft.TranslationJobID,
		PhaseType:           draft.PhaseType,
		State:               draft.State,
		ExecutionOrder:      draft.ExecutionOrder,
		ProgressPercent:     0,
		AIProvider:          draft.AIProvider,
		ModelName:           draft.ModelName,
		ExecutionMode:       draft.ExecutionMode,
		CredentialRef:       draft.CredentialRef,
		InstructionKind:     draft.InstructionKind,
		LatestExternalRunID: "",
		LatestError:         "",
		StartedAt:           nil,
		FinishedAt:          nil,
	}
	q, args, err := sqlx.Named(insertJobPhaseRun, row)
	if err != nil {
		return JobPhaseRun{}, fmt.Errorf("create job_phase_run named: %w", err)
	}
	result, err := ext.ExecContext(ctx, q, args...)
	if err != nil {
		return JobPhaseRun{}, mapFoundationSQLError(err, "create job_phase_run")
	}
	id, err := result.LastInsertId()
	if err != nil {
		return JobPhaseRun{}, fmt.Errorf("create job_phase_run last insert id: %w", err)
	}
	return r.GetJobPhaseRunByID(ctx, id)
}

// GetJobPhaseRunByID は ID で JobPhaseRun を取得する。
func (r *SQLiteJobLifecycleRepository) GetJobPhaseRunByID(
	ctx context.Context,
	id int64,
) (JobPhaseRun, error) {
	ext := extractTx(ctx, r.db)
	var row jobPhaseRunRow
	if err := sqlx.GetContext(ctx, ext, &row, selectJobPhaseRunByID, id); err != nil {
		return JobPhaseRun{}, mapSQLError(err, "get job_phase_run by id")
	}
	return row.toModel()
}

// UpdateJobPhaseRun は JobPhaseRun を更新する。
func (r *SQLiteJobLifecycleRepository) UpdateJobPhaseRun(
	ctx context.Context,
	id int64,
	draft JobPhaseRunUpdateDraft,
) (JobPhaseRun, error) {
	ext := extractTx(ctx, r.db)
	var startedAt *string
	if draft.StartedAt != nil {
		s := draft.StartedAt.UTC().Format(time.RFC3339)
		startedAt = &s
	}
	var finishedAt *string
	if draft.FinishedAt != nil {
		s := draft.FinishedAt.UTC().Format(time.RFC3339)
		finishedAt = &s
	}
	args := map[string]interface{}{
		"id":                     id,
		"state":                  draft.State,
		"progress_percent":       draft.ProgressPercent,
		"latest_external_run_id": draft.LatestExternalRunID,
		"latest_error":           draft.LatestError,
		"started_at":             startedAt,
		"finished_at":            finishedAt,
	}
	q, qArgs, err := sqlx.Named(updateJobPhaseRun, args)
	if err != nil {
		return JobPhaseRun{}, fmt.Errorf("update job_phase_run named: %w", err)
	}
	if _, err := ext.ExecContext(ctx, q, qArgs...); err != nil {
		return JobPhaseRun{}, mapFoundationSQLError(err, "update job_phase_run")
	}
	return r.GetJobPhaseRunByID(ctx, id)
}

// ListJobPhaseRunsByJobID は JobID に紐づく JobPhaseRun 一覧を返す。
func (r *SQLiteJobLifecycleRepository) ListJobPhaseRunsByJobID(
	ctx context.Context,
	jobID int64,
) ([]JobPhaseRun, error) {
	ext := extractTx(ctx, r.db)
	var rows []jobPhaseRunRow
	if err := sqlx.SelectContext(ctx, ext, &rows, selectJobPhaseRunsByJobID, jobID); err != nil {
		return nil, mapSQLError(err, "list job_phase_runs by job_id")
	}
	result := make([]JobPhaseRun, len(rows))
	for i, row := range rows {
		m, err := row.toModel()
		if err != nil {
			return nil, err
		}
		result[i] = m
	}
	return result, nil
}

// ---------------------------------------------------------------------------
// PhaseRunTranslationField
// ---------------------------------------------------------------------------

// CreatePhaseRunTranslationField は PhaseRunTranslationField レコードを作成する。
func (r *SQLiteJobLifecycleRepository) CreatePhaseRunTranslationField(
	ctx context.Context,
	draft PhaseRunTranslationFieldDraft,
) (PhaseRunTranslationField, error) {
	ext := extractTx(ctx, r.db)
	row := phaseRunTranslationFieldRow{
		PhaseRunID:            draft.PhaseRunID,
		JobTranslationFieldID: draft.JobTranslationFieldID,
		Role:                  draft.Role,
	}
	q, args, err := sqlx.Named(insertPhaseRunTranslationField, row)
	if err != nil {
		return PhaseRunTranslationField{}, fmt.Errorf("create phase_run_translation_field named: %w", err)
	}
	result, err := ext.ExecContext(ctx, q, args...)
	if err != nil {
		return PhaseRunTranslationField{}, mapFoundationSQLError(err, "create phase_run_translation_field")
	}
	id, err := result.LastInsertId()
	if err != nil {
		return PhaseRunTranslationField{}, fmt.Errorf("create phase_run_translation_field last insert id: %w", err)
	}
	return PhaseRunTranslationField{
		ID:                    id,
		PhaseRunID:            draft.PhaseRunID,
		JobTranslationFieldID: draft.JobTranslationFieldID,
		Role:                  draft.Role,
	}, nil
}

// ---------------------------------------------------------------------------
// PhaseRunPersona
// ---------------------------------------------------------------------------

// CreatePhaseRunPersona は PhaseRunPersona レコードを作成する。
func (r *SQLiteJobLifecycleRepository) CreatePhaseRunPersona(
	ctx context.Context,
	draft PhaseRunPersonaDraft,
) (PhaseRunPersona, error) {
	ext := extractTx(ctx, r.db)
	row := phaseRunPersonaRow{
		PhaseRunID: draft.PhaseRunID,
		PersonaID:  draft.PersonaID,
		Role:       draft.Role,
	}
	q, args, err := sqlx.Named(insertPhaseRunPersona, row)
	if err != nil {
		return PhaseRunPersona{}, fmt.Errorf("create phase_run_persona named: %w", err)
	}
	result, err := ext.ExecContext(ctx, q, args...)
	if err != nil {
		return PhaseRunPersona{}, mapFoundationSQLError(err, "create phase_run_persona")
	}
	id, err := result.LastInsertId()
	if err != nil {
		return PhaseRunPersona{}, fmt.Errorf("create phase_run_persona last insert id: %w", err)
	}
	return PhaseRunPersona{
		ID:         id,
		PhaseRunID: draft.PhaseRunID,
		PersonaID:  draft.PersonaID,
		Role:       draft.Role,
	}, nil
}

// ---------------------------------------------------------------------------
// PhaseRunDictionaryEntry
// ---------------------------------------------------------------------------

// CreatePhaseRunDictionaryEntry は PhaseRunDictionaryEntry レコードを作成する。
func (r *SQLiteJobLifecycleRepository) CreatePhaseRunDictionaryEntry(
	ctx context.Context,
	draft PhaseRunDictionaryEntryDraft,
) (PhaseRunDictionaryEntry, error) {
	ext := extractTx(ctx, r.db)
	row := phaseRunDictionaryEntryRow{
		PhaseRunID:        draft.PhaseRunID,
		DictionaryEntryID: draft.DictionaryEntryID,
		Role:              draft.Role,
	}
	q, args, err := sqlx.Named(insertPhaseRunDictionaryEntry, row)
	if err != nil {
		return PhaseRunDictionaryEntry{}, fmt.Errorf("create phase_run_dictionary_entry named: %w", err)
	}
	result, err := ext.ExecContext(ctx, q, args...)
	if err != nil {
		return PhaseRunDictionaryEntry{}, mapFoundationSQLError(err, "create phase_run_dictionary_entry")
	}
	id, err := result.LastInsertId()
	if err != nil {
		return PhaseRunDictionaryEntry{}, fmt.Errorf("create phase_run_dictionary_entry last insert id: %w", err)
	}
	return PhaseRunDictionaryEntry{
		ID:                id,
		PhaseRunID:        draft.PhaseRunID,
		DictionaryEntryID: draft.DictionaryEntryID,
		Role:              draft.Role,
	}, nil
}
