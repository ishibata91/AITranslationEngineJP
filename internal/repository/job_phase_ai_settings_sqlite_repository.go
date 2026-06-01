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
	jobPhaseAISettingsTimestampLayout = time.RFC3339Nano

	upsertJobPhaseAISettingsSQL = `
INSERT INTO JOB_PHASE_AI_SETTINGS (
  phase_type,
  ai_provider,
  model_name,
  execution_mode,
  batch_mode,
  updated_at
) VALUES (
  :phase_type,
  :ai_provider,
  :model_name,
  :execution_mode,
  :batch_mode,
  :updated_at
)
ON CONFLICT(phase_type) DO UPDATE
SET ai_provider    = excluded.ai_provider,
    model_name     = excluded.model_name,
    execution_mode = excluded.execution_mode,
    batch_mode     = excluded.batch_mode,
    updated_at     = excluded.updated_at;`

	selectJobPhaseAISettingsByPhaseTypeSQL = `
SELECT
  phase_type,
  ai_provider,
  model_name,
  execution_mode,
  batch_mode,
  updated_at
FROM JOB_PHASE_AI_SETTINGS
WHERE phase_type = ?
LIMIT 1;`
)

// SQLiteJobPhaseAISettingsRepository は JobPhaseAISettingsRepository の SQLite 実装。
type SQLiteJobPhaseAISettingsRepository struct {
	db *sqlx.DB
}

// NewSQLiteJobPhaseAISettingsRepository は JobPhaseAISettingsRepository を返す。
func NewSQLiteJobPhaseAISettingsRepository(db *sqlx.DB) JobPhaseAISettingsRepository {
	return &SQLiteJobPhaseAISettingsRepository{db: db}
}

type jobPhaseAISettingsRow struct {
	PhaseType     string `db:"phase_type"`
	AIProvider    string `db:"ai_provider"`
	ModelName     string `db:"model_name"`
	ExecutionMode string `db:"execution_mode"`
	BatchMode     string `db:"batch_mode"`
	UpdatedAt     string `db:"updated_at"`
}

func (row jobPhaseAISettingsRow) toModel() (JobPhaseAISettings, error) {
	updatedAt, err := time.Parse(jobPhaseAISettingsTimestampLayout, row.UpdatedAt)
	if err != nil {
		return JobPhaseAISettings{}, fmt.Errorf("parse job_phase_ai_settings updated_at: %w", err)
	}
	return JobPhaseAISettings{
		PhaseType:     strings.TrimSpace(row.PhaseType),
		AIProvider:    strings.TrimSpace(row.AIProvider),
		ModelName:     strings.TrimSpace(row.ModelName),
		ExecutionMode: strings.TrimSpace(row.ExecutionMode),
		BatchMode:     strings.TrimSpace(row.BatchMode),
		UpdatedAt:     updatedAt,
	}, nil
}

// Save は draft を phase_type を主キーとして upsert する。
func (r *SQLiteJobPhaseAISettingsRepository) Save(
	ctx context.Context,
	draft JobPhaseAISettingsDraft,
) (JobPhaseAISettings, error) {
	row := jobPhaseAISettingsRow{
		PhaseType:     strings.TrimSpace(draft.PhaseType),
		AIProvider:    strings.TrimSpace(draft.AIProvider),
		ModelName:     strings.TrimSpace(draft.ModelName),
		ExecutionMode: strings.TrimSpace(draft.ExecutionMode),
		BatchMode:     strings.TrimSpace(draft.BatchMode),
		UpdatedAt:     draft.UpdatedAt.UTC().Format(jobPhaseAISettingsTimestampLayout),
	}
	if _, err := sqlx.NamedExecContext(ctx, extractTx(ctx, r.db), upsertJobPhaseAISettingsSQL, row); err != nil {
		return JobPhaseAISettings{}, fmt.Errorf("save job_phase_ai_settings: %w", err)
	}
	return r.LoadByPhaseType(ctx, draft.PhaseType)
}

// LoadByPhaseType は phase_type に対応する record を返す。
// record が存在しない場合は ErrJobPhaseAISettingsNotFound を返す。
func (r *SQLiteJobPhaseAISettingsRepository) LoadByPhaseType(
	ctx context.Context,
	phaseType string,
) (JobPhaseAISettings, error) {
	var row jobPhaseAISettingsRow
	if err := sqlx.GetContext(
		ctx,
		extractTx(ctx, r.db),
		&row,
		selectJobPhaseAISettingsByPhaseTypeSQL,
		strings.TrimSpace(phaseType),
	); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return JobPhaseAISettings{}, fmt.Errorf("%w: phase_type=%s", ErrJobPhaseAISettingsNotFound, phaseType)
		}
		return JobPhaseAISettings{}, fmt.Errorf("load job_phase_ai_settings by phase_type: %w", err)
	}
	return row.toModel()
}
