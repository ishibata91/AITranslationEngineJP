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
	ID                     int64   `db:"id"`
	TranslationJobID       int64   `db:"translation_job_id"`
	PhaseType              string  `db:"phase_type"`
	State                  string  `db:"state"`
	ExecutionOrder         int     `db:"execution_order"`
	ProgressPercent        int     `db:"progress_percent"`
	SnapshotFieldCount     int     `db:"snapshot_field_count"`
	ProviderTargetCount    int     `db:"provider_target_count"`
	ExactExclusionCount    int     `db:"exact_exclusion_count"`
	PartialConstraintCount int     `db:"partial_constraint_count"`
	AIProvider             string  `db:"ai_provider"`
	ModelName              string  `db:"model_name"`
	ExecutionMode          string  `db:"execution_mode"`
	CredentialRef          string  `db:"credential_ref"`
	InstructionKind        string  `db:"instruction_kind"`
	InputSnapshotDigest    string  `db:"input_snapshot_digest"`
	DictionaryDigest       string  `db:"dictionary_digest"`
	PersonaDigest          string  `db:"persona_digest"`
	MetadataDigest         string  `db:"metadata_digest"`
	PromptDigest           string  `db:"prompt_digest"`
	LatestExternalRunID    string  `db:"latest_external_run_id"`
	LatestError            string  `db:"latest_error"`
	StartedAt              *string `db:"started_at"`
	FinishedAt             *string `db:"finished_at"`
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
		ID:                     r.ID,
		TranslationJobID:       r.TranslationJobID,
		PhaseType:              r.PhaseType,
		State:                  r.State,
		ExecutionOrder:         r.ExecutionOrder,
		ProgressPercent:        r.ProgressPercent,
		SnapshotFieldCount:     r.SnapshotFieldCount,
		ProviderTargetCount:    r.ProviderTargetCount,
		ExactExclusionCount:    r.ExactExclusionCount,
		PartialConstraintCount: r.PartialConstraintCount,
		AIProvider:             r.AIProvider,
		ModelName:              r.ModelName,
		ExecutionMode:          r.ExecutionMode,
		CredentialRef:          r.CredentialRef,
		InstructionKind:        r.InstructionKind,
		InputSnapshotDigest:    r.InputSnapshotDigest,
		DictionaryDigest:       r.DictionaryDigest,
		PersonaDigest:          r.PersonaDigest,
		MetadataDigest:         r.MetadataDigest,
		PromptDigest:           r.PromptDigest,
		LatestExternalRunID:    r.LatestExternalRunID,
		LatestError:            r.LatestError,
		StartedAt:              startedAt,
		FinishedAt:             finishedAt,
	}, nil
}

type phaseRunTranslationFieldRow struct {
	ID                    int64  `db:"id"`
	PhaseRunID            int64  `db:"phase_run_id"`
	JobTranslationFieldID int64  `db:"job_translation_field_id"`
	Role                  string `db:"role"`
}

func (r phaseRunTranslationFieldRow) toModel() PhaseRunTranslationField {
	return PhaseRunTranslationField(r)
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

type translationJobPhaseRuntimeSnapshotRow struct {
	ID               int64  `db:"id"`
	TranslationJobID int64  `db:"translation_job_id"`
	PhaseID          string `db:"phase_id"`
	Provider         string `db:"provider"`
	ModelName        string `db:"model_name"`
	CredentialStatus string `db:"credential_status"`
	ExecutionMode    string `db:"execution_mode"`
	BatchMode        string `db:"batch_mode"`
	CreatedAt        string `db:"created_at"`
}

func (r translationJobPhaseRuntimeSnapshotRow) toModel() (TranslationJobPhaseRuntimeSnapshot, error) {
	createdAt, err := time.Parse(time.RFC3339, r.CreatedAt)
	if err != nil {
		return TranslationJobPhaseRuntimeSnapshot{}, fmt.Errorf("parse phase runtime created_at: %w", err)
	}
	return TranslationJobPhaseRuntimeSnapshot{
		ID:               r.ID,
		TranslationJobID: r.TranslationJobID,
		PhaseID:          r.PhaseID,
		Provider:         r.Provider,
		ModelName:        r.ModelName,
		CredentialStatus: r.CredentialStatus,
		ExecutionMode:    r.ExecutionMode,
		BatchMode:        r.BatchMode,
		CreatedAt:        createdAt,
	}, nil
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

	selectTranslationJobs = `
SELECT id, x_edit_extracted_data_id, job_name, state, progress_percent, created_at, started_at, finished_at
FROM TRANSLATION_JOB
ORDER BY id ASC`

	selectIncompleteTranslationJobs = `
SELECT id, x_edit_extracted_data_id, job_name, state, progress_percent, created_at, started_at, finished_at
FROM TRANSLATION_JOB
WHERE LOWER(TRIM(state)) <> 'completed'
ORDER BY created_at DESC, id DESC`

	updateTranslationJob = `
UPDATE TRANSLATION_JOB SET
  job_name         = :job_name,
  state            = :state,
  progress_percent = :progress_percent,
  started_at       = :started_at,
  finished_at      = :finished_at
WHERE id = :id`

	insertTranslationJobPhaseRuntimeSnapshot = `
INSERT INTO TRANSLATION_JOB_PHASE_RUNTIME_SNAPSHOT
  (translation_job_id, phase_id, provider, model_name, credential_status,
   execution_mode, batch_mode, created_at)
VALUES
  (:translation_job_id, :phase_id, :provider, :model_name, :credential_status,
   :execution_mode, :batch_mode, :created_at)
ON CONFLICT(translation_job_id, phase_id) DO UPDATE SET
  provider = excluded.provider,
  model_name = excluded.model_name,
  credential_status = excluded.credential_status,
  execution_mode = excluded.execution_mode,
  batch_mode = excluded.batch_mode,
  created_at = excluded.created_at`

	selectTranslationJobPhaseRuntimeSnapshotsByJobID = `
SELECT id, translation_job_id, phase_id, provider, model_name, credential_status,
       execution_mode, batch_mode, created_at
FROM TRANSLATION_JOB_PHASE_RUNTIME_SNAPSHOT
WHERE translation_job_id = ?
ORDER BY id ASC`

	selectTranslationJobPhaseRuntimeSnapshotByJobAndPhase = `
SELECT id, translation_job_id, phase_id, provider, model_name, credential_status,
       execution_mode, batch_mode, created_at
FROM TRANSLATION_JOB_PHASE_RUNTIME_SNAPSHOT
WHERE translation_job_id = ? AND phase_id = ?
LIMIT 1`

	selectTranslationJobForDelete = `
SELECT id, x_edit_extracted_data_id, job_name, state, progress_percent, created_at, started_at, finished_at
FROM TRANSLATION_JOB
WHERE id = ?`

	deletePersonaFieldEvidenceByJobID = `
DELETE FROM PERSONA_FIELD_EVIDENCE
WHERE persona_id IN (
	SELECT id
	FROM PERSONA
	WHERE translation_job_id = ?
)`

	deleteXTranslatorOutputRowsByJobID = `
DELETE FROM XTRANSLATOR_OUTPUT_ROW
WHERE translation_artifact_id IN (
	SELECT id
	FROM TRANSLATION_ARTIFACT
	WHERE translation_job_id = ?
)`

	deletePhaseRunTranslationFieldsByJobID = `
DELETE FROM PHASE_RUN_TRANSLATION_FIELD
WHERE phase_run_id IN (
	SELECT id
	FROM JOB_PHASE_RUN
	WHERE translation_job_id = ?
)`

	deletePhaseRunPersonasByJobID = `
DELETE FROM PHASE_RUN_PERSONA
WHERE phase_run_id IN (
	SELECT id
	FROM JOB_PHASE_RUN
	WHERE translation_job_id = ?
)`

	deletePhaseRunDictionaryEntriesByJobID = `
DELETE FROM PHASE_RUN_DICTIONARY_ENTRY
WHERE phase_run_id IN (
	SELECT id
	FROM JOB_PHASE_RUN
	WHERE translation_job_id = ?
)`

	deleteTranslationJobPhaseRuntimeSnapshotsByJobID = `
DELETE FROM TRANSLATION_JOB_PHASE_RUNTIME_SNAPSHOT
WHERE translation_job_id = ?`

	deleteJobPhaseRunsByJobID = `
DELETE FROM JOB_PHASE_RUN
WHERE translation_job_id = ?`

	deleteJobTranslationFieldsByJobID = `
DELETE FROM JOB_TRANSLATION_FIELD
WHERE translation_job_id = ?`

	deleteTranslationArtifactsByJobID = `
DELETE FROM TRANSLATION_ARTIFACT
WHERE translation_job_id = ?`

	deletePersonasByJobID = `
DELETE FROM PERSONA
WHERE translation_job_id = ?`

	deleteDictionaryEntriesByJobID = `
DELETE FROM DICTIONARY_ENTRY
WHERE translation_job_id = ?`

	deleteTranslationJobByID = `
DELETE FROM TRANSLATION_JOB
WHERE id = ?`

	insertJobPhaseRun = `
INSERT INTO JOB_PHASE_RUN
  (translation_job_id, phase_type, state, execution_order, progress_percent,
   snapshot_field_count, provider_target_count, exact_exclusion_count, partial_constraint_count,
   ai_provider, model_name, execution_mode, credential_ref, instruction_kind,
   input_snapshot_digest, dictionary_digest, persona_digest, metadata_digest, prompt_digest,
   latest_external_run_id, latest_error, started_at, finished_at)
VALUES
  (:translation_job_id, :phase_type, :state, :execution_order, :progress_percent,
   :snapshot_field_count, :provider_target_count, :exact_exclusion_count, :partial_constraint_count,
   :ai_provider, :model_name, :execution_mode, :credential_ref, :instruction_kind,
   :input_snapshot_digest, :dictionary_digest, :persona_digest, :metadata_digest, :prompt_digest,
   :latest_external_run_id, :latest_error, :started_at, :finished_at)`

	selectJobPhaseRunByID = `
SELECT id, translation_job_id, phase_type, state, execution_order, progress_percent,
       snapshot_field_count, provider_target_count, exact_exclusion_count, partial_constraint_count,
       ai_provider, model_name, execution_mode, credential_ref, instruction_kind,
       input_snapshot_digest, dictionary_digest, persona_digest, metadata_digest, prompt_digest,
       latest_external_run_id, latest_error, started_at, finished_at
FROM JOB_PHASE_RUN WHERE id = ?`

	updateJobPhaseRunSetClause = `
  state                 = :state,
  progress_percent      = :progress_percent,
  snapshot_field_count  = :snapshot_field_count,
  provider_target_count = :provider_target_count,
  exact_exclusion_count = :exact_exclusion_count,
  partial_constraint_count = :partial_constraint_count,
  ai_provider           = :ai_provider,
  model_name            = :model_name,
  execution_mode        = :execution_mode,
  credential_ref        = :credential_ref,
  instruction_kind      = :instruction_kind,
  input_snapshot_digest = :input_snapshot_digest,
  dictionary_digest     = :dictionary_digest,
  persona_digest        = :persona_digest,
  metadata_digest       = :metadata_digest,
  prompt_digest         = :prompt_digest,
  latest_external_run_id = :latest_external_run_id,
  latest_error          = :latest_error,
  started_at            = :started_at,
  finished_at           = :finished_at`

	updateJobPhaseRun = `
UPDATE JOB_PHASE_RUN SET
` + updateJobPhaseRunSetClause + `
WHERE id = :id`

	updateJobPhaseRunWhenState = `
UPDATE JOB_PHASE_RUN SET
` + updateJobPhaseRunSetClause + `
WHERE id = :id AND state = :expected_state`

	selectJobPhaseRunsByJobID = `
SELECT id, translation_job_id, phase_type, state, execution_order, progress_percent,
       snapshot_field_count, provider_target_count, exact_exclusion_count, partial_constraint_count,
       ai_provider, model_name, execution_mode, credential_ref, instruction_kind,
       input_snapshot_digest, dictionary_digest, persona_digest, metadata_digest, prompt_digest,
       latest_external_run_id, latest_error, started_at, finished_at
FROM JOB_PHASE_RUN WHERE translation_job_id = ? ORDER BY execution_order ASC`

	selectJobPhaseRunByJobAndType = `
SELECT id, translation_job_id, phase_type, state, execution_order, progress_percent,
       snapshot_field_count, provider_target_count, exact_exclusion_count, partial_constraint_count,
       ai_provider, model_name, execution_mode, credential_ref, instruction_kind,
       input_snapshot_digest, dictionary_digest, persona_digest, metadata_digest, prompt_digest,
       latest_external_run_id, latest_error, started_at, finished_at
FROM JOB_PHASE_RUN
WHERE translation_job_id = ? AND phase_type = ?
LIMIT 1`

	insertPhaseRunTranslationField = `
INSERT INTO PHASE_RUN_TRANSLATION_FIELD
  (phase_run_id, job_translation_field_id, role)
VALUES
  (:phase_run_id, :job_translation_field_id, :role)`

	selectPhaseRunTranslationFieldsByPhaseRunID = `
SELECT id, phase_run_id, job_translation_field_id, role
FROM PHASE_RUN_TRANSLATION_FIELD
WHERE phase_run_id = ?
ORDER BY id ASC`

	insertPhaseRunPersona = `
INSERT INTO PHASE_RUN_PERSONA
  (phase_run_id, persona_id, role)
VALUES
  (:phase_run_id, :persona_id, :role)`

	selectPhaseRunPersonasByPhaseRunID = `
SELECT id, phase_run_id, persona_id, role
FROM PHASE_RUN_PERSONA
WHERE phase_run_id = ?
ORDER BY id ASC`

	insertPhaseRunDictionaryEntry = `
INSERT INTO PHASE_RUN_DICTIONARY_ENTRY
  (phase_run_id, dictionary_entry_id, role)
VALUES
  (:phase_run_id, :dictionary_entry_id, :role)`

	selectPhaseRunDictionaryEntriesByPhaseRunID = `
SELECT id, phase_run_id, dictionary_entry_id, role
FROM PHASE_RUN_DICTIONARY_ENTRY
WHERE phase_run_id = ?
ORDER BY id ASC`
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

// ListTranslationJobs returns all translation jobs for read-model queries.
func (r *SQLiteJobLifecycleRepository) ListTranslationJobs(ctx context.Context) ([]TranslationJob, error) {
	ext := extractTx(ctx, r.db)
	rows := make([]translationJobRow, 0)
	if err := sqlx.SelectContext(ctx, ext, &rows, selectTranslationJobs); err != nil {
		return nil, mapSQLError(err, "list translation jobs")
	}

	jobs := make([]TranslationJob, 0, len(rows))
	for _, row := range rows {
		job, err := row.toModel()
		if err != nil {
			return nil, err
		}
		jobs = append(jobs, job)
	}
	return jobs, nil
}

// ListIncompleteTranslationJobs returns non-completed translation jobs for Job Management.
func (r *SQLiteJobLifecycleRepository) ListIncompleteTranslationJobs(
	ctx context.Context,
) ([]TranslationJob, error) {
	ext := extractTx(ctx, r.db)
	rows := make([]translationJobRow, 0)
	if err := sqlx.SelectContext(ctx, ext, &rows, selectIncompleteTranslationJobs); err != nil {
		return nil, mapSQLError(err, "list incomplete translation jobs")
	}

	jobs := make([]TranslationJob, 0, len(rows))
	for _, row := range rows {
		job, err := row.toModel()
		if err != nil {
			return nil, err
		}
		jobs = append(jobs, job)
	}
	return jobs, nil
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

// DeleteNonRunningTranslationJob deletes one non-running job and its job-owned rows.
func (r *SQLiteJobLifecycleRepository) DeleteNonRunningTranslationJob(
	ctx context.Context,
	jobID int64,
) (TranslationJobDeleteResult, error) {
	ext := extractTx(ctx, r.db)
	var row translationJobRow
	if err := sqlx.GetContext(ctx, ext, &row, selectTranslationJobForDelete, jobID); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return TranslationJobDeleteResult{Outcome: TranslationJobDeleteOutcomeNotFound}, nil
		}
		return TranslationJobDeleteResult{}, mapSQLError(err, "select translation_job for delete")
	}

	job, err := row.toModel()
	if err != nil {
		return TranslationJobDeleteResult{}, err
	}
	if normalizeTranslationJobManagementState(job.State) == "running" {
		return TranslationJobDeleteResult{
			Outcome: TranslationJobDeleteOutcomeBlockedRunning,
			Job:     &job,
		}, nil
	}
	phaseRuns, err := r.ListJobPhaseRunsByJobID(ctx, job.ID)
	if err != nil {
		return TranslationJobDeleteResult{}, fmt.Errorf("list job_phase_run for delete guard: %w", err)
	}
	if hasUnsafeDeletePhaseRun(phaseRuns) {
		return TranslationJobDeleteResult{
			Outcome: TranslationJobDeleteOutcomeBlockedUnsafe,
			Job:     &job,
		}, nil
	}

	statements := []struct {
		label string
		query string
	}{
		{label: "delete persona_field_evidence by job_id", query: deletePersonaFieldEvidenceByJobID},
		{label: "delete xtranslator_output_row by job_id", query: deleteXTranslatorOutputRowsByJobID},
		{label: "delete phase_run_translation_field by job_id", query: deletePhaseRunTranslationFieldsByJobID},
		{label: "delete phase_run_persona by job_id", query: deletePhaseRunPersonasByJobID},
		{label: "delete phase_run_dictionary_entry by job_id", query: deletePhaseRunDictionaryEntriesByJobID},
		{label: "delete translation_job_phase_runtime_snapshot by job_id", query: deleteTranslationJobPhaseRuntimeSnapshotsByJobID},
		{label: "delete job_phase_run by job_id", query: deleteJobPhaseRunsByJobID},
		{label: "delete job_translation_field by job_id", query: deleteJobTranslationFieldsByJobID},
		{label: "delete translation_artifact by job_id", query: deleteTranslationArtifactsByJobID},
		{label: "delete persona by job_id", query: deletePersonasByJobID},
		{label: "delete dictionary_entry by job_id", query: deleteDictionaryEntriesByJobID},
		{label: "delete translation_job by id", query: deleteTranslationJobByID},
	}
	for _, statement := range statements {
		if _, execErr := ext.ExecContext(ctx, statement.query, jobID); execErr != nil {
			return TranslationJobDeleteResult{}, mapFoundationSQLError(execErr, statement.label)
		}
	}
	return TranslationJobDeleteResult{
		Outcome: TranslationJobDeleteOutcomeDeleted,
		Job:     &job,
	}, nil
}

func normalizeTranslationJobManagementState(state string) string {
	return strings.ToLower(strings.TrimSpace(state))
}

func hasUnsafeDeletePhaseRun(phaseRuns []JobPhaseRun) bool {
	for _, phaseRun := range phaseRuns {
		state := normalizeTranslationJobManagementState(phaseRun.State)
		if state == "running" {
			return true
		}
		if normalizeTranslationJobManagementState(phaseRun.LatestError) == "stop_requested" {
			return true
		}
	}
	return false
}

// SaveTranslationJobPhaseRuntimeSnapshot は Job Setup の phase 別 runtime snapshot を保存する。
func (r *SQLiteJobLifecycleRepository) SaveTranslationJobPhaseRuntimeSnapshot(
	ctx context.Context,
	draft TranslationJobPhaseRuntimeSnapshotDraft,
) (TranslationJobPhaseRuntimeSnapshot, error) {
	ext := extractTx(ctx, r.db)
	row := translationJobPhaseRuntimeSnapshotRow{
		TranslationJobID: draft.TranslationJobID,
		PhaseID:          draft.PhaseID,
		Provider:         draft.Provider,
		ModelName:        draft.ModelName,
		CredentialStatus: draft.CredentialStatus,
		ExecutionMode:    draft.ExecutionMode,
		BatchMode:        draft.BatchMode,
		CreatedAt:        time.Now().UTC().Format(time.RFC3339),
	}
	q, args, err := sqlx.Named(insertTranslationJobPhaseRuntimeSnapshot, row)
	if err != nil {
		return TranslationJobPhaseRuntimeSnapshot{}, fmt.Errorf("create translation_job_phase_runtime_snapshot named: %w", err)
	}
	result, err := ext.ExecContext(ctx, q, args...)
	if err != nil {
		return TranslationJobPhaseRuntimeSnapshot{}, mapFoundationSQLError(err, "create translation_job_phase_runtime_snapshot")
	}
	if _, err := result.RowsAffected(); err != nil {
		return TranslationJobPhaseRuntimeSnapshot{}, fmt.Errorf("create translation_job_phase_runtime_snapshot rows affected: %w", err)
	}
	return r.GetTranslationJobPhaseRuntimeSnapshot(ctx, draft.TranslationJobID, draft.PhaseID)
}

// ListTranslationJobPhaseRuntimeSnapshots は translation_job_id に紐づく phase runtime snapshot を返す。
func (r *SQLiteJobLifecycleRepository) ListTranslationJobPhaseRuntimeSnapshots(
	ctx context.Context,
	translationJobID int64,
) ([]TranslationJobPhaseRuntimeSnapshot, error) {
	ext := extractTx(ctx, r.db)
	rows := make([]translationJobPhaseRuntimeSnapshotRow, 0)
	if err := sqlx.SelectContext(ctx, ext, &rows, selectTranslationJobPhaseRuntimeSnapshotsByJobID, translationJobID); err != nil {
		return nil, mapSQLError(err, "list translation_job_phase_runtime_snapshot by job_id")
	}
	result := make([]TranslationJobPhaseRuntimeSnapshot, 0, len(rows))
	for _, row := range rows {
		snapshot, err := row.toModel()
		if err != nil {
			return nil, err
		}
		result = append(result, snapshot)
	}
	return result, nil
}

// GetTranslationJobPhaseRuntimeSnapshot は translation_job_id と phase_id で phase runtime snapshot を取得する。
func (r *SQLiteJobLifecycleRepository) GetTranslationJobPhaseRuntimeSnapshot(
	ctx context.Context,
	translationJobID int64,
	phaseID string,
) (TranslationJobPhaseRuntimeSnapshot, error) {
	ext := extractTx(ctx, r.db)
	var row translationJobPhaseRuntimeSnapshotRow
	if err := sqlx.GetContext(ctx, ext, &row, selectTranslationJobPhaseRuntimeSnapshotByJobAndPhase, translationJobID, phaseID); err != nil {
		return TranslationJobPhaseRuntimeSnapshot{}, mapSQLError(err, "get translation_job_phase_runtime_snapshot by job and phase")
	}
	return row.toModel()
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
		TranslationJobID:       draft.TranslationJobID,
		PhaseType:              draft.PhaseType,
		State:                  draft.State,
		ExecutionOrder:         draft.ExecutionOrder,
		ProgressPercent:        0,
		SnapshotFieldCount:     draft.SnapshotFieldCount,
		ProviderTargetCount:    draft.ProviderTargetCount,
		ExactExclusionCount:    draft.ExactExclusionCount,
		PartialConstraintCount: draft.PartialConstraintCount,
		AIProvider:             draft.AIProvider,
		ModelName:              draft.ModelName,
		ExecutionMode:          draft.ExecutionMode,
		CredentialRef:          draft.CredentialRef,
		InstructionKind:        draft.InstructionKind,
		InputSnapshotDigest:    draft.InputSnapshotDigest,
		DictionaryDigest:       draft.DictionaryDigest,
		PersonaDigest:          draft.PersonaDigest,
		MetadataDigest:         draft.MetadataDigest,
		PromptDigest:           draft.PromptDigest,
		LatestExternalRunID:    "",
		LatestError:            "",
		StartedAt:              nil,
		FinishedAt:             nil,
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
	current, err := r.GetJobPhaseRunByID(ctx, id)
	if err != nil {
		return JobPhaseRun{}, err
	}
	completedDraft := completeJobPhaseRunUpdateDraft(current, draft)
	args := buildJobPhaseRunUpdateArgs(id, "", completedDraft)
	q, qArgs, err := sqlx.Named(updateJobPhaseRun, args)
	if err != nil {
		return JobPhaseRun{}, fmt.Errorf("update job_phase_run named: %w", err)
	}
	result, err := ext.ExecContext(ctx, q, qArgs...)
	if err != nil {
		return JobPhaseRun{}, mapFoundationSQLError(err, "update job_phase_run")
	}
	affected, err := result.RowsAffected()
	if err != nil {
		return JobPhaseRun{}, fmt.Errorf("update job_phase_run rows affected: %w", err)
	}
	if affected == 0 {
		return r.GetJobPhaseRunByID(ctx, id)
	}
	return r.GetJobPhaseRunByID(ctx, id)
}

// UpdateJobPhaseRunWhenState は現在状態が expectedState の場合だけ JobPhaseRun を更新する。
func (r *SQLiteJobLifecycleRepository) UpdateJobPhaseRunWhenState(
	ctx context.Context,
	id int64,
	expectedState string,
	draft JobPhaseRunUpdateDraft,
) (JobPhaseRun, error) {
	ext := extractTx(ctx, r.db)
	current, err := r.GetJobPhaseRunByID(ctx, id)
	if err != nil {
		return JobPhaseRun{}, err
	}
	completedDraft := completeJobPhaseRunUpdateDraft(current, draft)
	args := buildJobPhaseRunUpdateArgs(id, expectedState, completedDraft)
	q, qArgs, err := sqlx.Named(updateJobPhaseRunWhenState, args)
	if err != nil {
		return JobPhaseRun{}, fmt.Errorf("update job_phase_run when state named: %w", err)
	}
	result, err := ext.ExecContext(ctx, q, qArgs...)
	if err != nil {
		return JobPhaseRun{}, mapFoundationSQLError(err, "update job_phase_run when state")
	}
	affected, err := result.RowsAffected()
	if err != nil {
		return JobPhaseRun{}, fmt.Errorf("update job_phase_run when state rows affected: %w", err)
	}
	if affected == 0 {
		if _, err := r.GetJobPhaseRunByID(ctx, id); err != nil {
			return JobPhaseRun{}, err
		}
		return JobPhaseRun{}, fmt.Errorf("update job_phase_run expected state: %w", ErrConflict)
	}
	return r.GetJobPhaseRunByID(ctx, id)
}

func completeJobPhaseRunUpdateDraft(
	current JobPhaseRun,
	draft JobPhaseRunUpdateDraft,
) JobPhaseRunUpdateDraft {
	completed := draft
	if completed.SnapshotFieldCount == 0 {
		completed.SnapshotFieldCount = current.SnapshotFieldCount
	}
	if completed.ProviderTargetCount == 0 {
		completed.ProviderTargetCount = current.ProviderTargetCount
	}
	if completed.ExactExclusionCount == 0 {
		completed.ExactExclusionCount = current.ExactExclusionCount
	}
	if completed.PartialConstraintCount == 0 {
		completed.PartialConstraintCount = current.PartialConstraintCount
	}
	if completed.AIProvider == "" {
		completed.AIProvider = current.AIProvider
	}
	if completed.ModelName == "" {
		completed.ModelName = current.ModelName
	}
	if completed.ExecutionMode == "" {
		completed.ExecutionMode = current.ExecutionMode
	}
	if completed.CredentialRef == "" {
		completed.CredentialRef = current.CredentialRef
	}
	if completed.InstructionKind == "" {
		completed.InstructionKind = current.InstructionKind
	}
	if completed.InputSnapshotDigest == "" {
		completed.InputSnapshotDigest = current.InputSnapshotDigest
	}
	if completed.DictionaryDigest == "" {
		completed.DictionaryDigest = current.DictionaryDigest
	}
	if completed.PersonaDigest == "" {
		completed.PersonaDigest = current.PersonaDigest
	}
	if completed.MetadataDigest == "" {
		completed.MetadataDigest = current.MetadataDigest
	}
	if completed.PromptDigest == "" {
		completed.PromptDigest = current.PromptDigest
	}
	return completed
}

func buildJobPhaseRunUpdateArgs(
	id int64,
	expectedState string,
	draft JobPhaseRunUpdateDraft,
) map[string]interface{} {
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
	return map[string]interface{}{
		"id":                       id,
		"expected_state":           expectedState,
		"state":                    draft.State,
		"progress_percent":         draft.ProgressPercent,
		"snapshot_field_count":     draft.SnapshotFieldCount,
		"provider_target_count":    draft.ProviderTargetCount,
		"exact_exclusion_count":    draft.ExactExclusionCount,
		"partial_constraint_count": draft.PartialConstraintCount,
		"ai_provider":              draft.AIProvider,
		"model_name":               draft.ModelName,
		"execution_mode":           draft.ExecutionMode,
		"credential_ref":           draft.CredentialRef,
		"instruction_kind":         draft.InstructionKind,
		"input_snapshot_digest":    draft.InputSnapshotDigest,
		"dictionary_digest":        draft.DictionaryDigest,
		"persona_digest":           draft.PersonaDigest,
		"metadata_digest":          draft.MetadataDigest,
		"prompt_digest":            draft.PromptDigest,
		"latest_external_run_id":   draft.LatestExternalRunID,
		"latest_error":             draft.LatestError,
		"started_at":               startedAt,
		"finished_at":              finishedAt,
	}
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

// FindJobPhaseRun は translation_job_id と phase_type で JobPhaseRun を取得する。
func (r *SQLiteJobLifecycleRepository) FindJobPhaseRun(
	ctx context.Context,
	translationJobID int64,
	phaseType string,
) (JobPhaseRun, error) {
	ext := extractTx(ctx, r.db)
	var row jobPhaseRunRow
	if err := sqlx.GetContext(
		ctx,
		ext,
		&row,
		selectJobPhaseRunByJobAndType,
		translationJobID,
		phaseType,
	); err != nil {
		return JobPhaseRun{}, mapSQLError(err, "find job_phase_run by job and type")
	}
	return row.toModel()
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

// ListPhaseRunTranslationFieldsByPhaseRunID は PhaseRunID に紐づく field 関連一覧を返す。
func (r *SQLiteJobLifecycleRepository) ListPhaseRunTranslationFieldsByPhaseRunID(
	ctx context.Context,
	phaseRunID int64,
) ([]PhaseRunTranslationField, error) {
	ext := extractTx(ctx, r.db)
	var rows []phaseRunTranslationFieldRow
	if err := sqlx.SelectContext(ctx, ext, &rows, selectPhaseRunTranslationFieldsByPhaseRunID, phaseRunID); err != nil {
		return nil, mapSQLError(err, "list phase_run_translation_fields by phase_run_id")
	}
	result := make([]PhaseRunTranslationField, 0, len(rows))
	for _, row := range rows {
		result = append(result, row.toModel())
	}
	return result, nil
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

// ListPhaseRunPersonasByPhaseRunID は PhaseRunID に紐づく persona 関連一覧を返す。
func (r *SQLiteJobLifecycleRepository) ListPhaseRunPersonasByPhaseRunID(
	ctx context.Context,
	phaseRunID int64,
) ([]PhaseRunPersona, error) {
	ext := extractTx(ctx, r.db)
	var rows []phaseRunPersonaRow
	if err := sqlx.SelectContext(ctx, ext, &rows, selectPhaseRunPersonasByPhaseRunID, phaseRunID); err != nil {
		return nil, mapSQLError(err, "list phase_run_persona by phase_run_id")
	}
	result := make([]PhaseRunPersona, 0, len(rows))
	for _, row := range rows {
		result = append(result, PhaseRunPersona(row))
	}
	return result, nil
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

// ListPhaseRunDictionaryEntriesByPhaseRunID は PhaseRunID に紐づく辞書関連一覧を返す。
func (r *SQLiteJobLifecycleRepository) ListPhaseRunDictionaryEntriesByPhaseRunID(
	ctx context.Context,
	phaseRunID int64,
) ([]PhaseRunDictionaryEntry, error) {
	ext := extractTx(ctx, r.db)
	var rows []phaseRunDictionaryEntryRow
	if err := sqlx.SelectContext(ctx, ext, &rows, selectPhaseRunDictionaryEntriesByPhaseRunID, phaseRunID); err != nil {
		return nil, mapSQLError(err, "list phase_run_dictionary_entry by phase_run_id")
	}
	result := make([]PhaseRunDictionaryEntry, 0, len(rows))
	for _, row := range rows {
		result = append(result, PhaseRunDictionaryEntry(row))
	}
	return result, nil
}
