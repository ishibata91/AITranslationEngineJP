package repository

import (
	"context"
	"errors"
	"time"
)

var (
	// ErrJobPhaseAISettingsNotFound は対象 phase_type の JOB_PHASE_AI_SETTINGS record が存在しない場合に返す。
	ErrJobPhaseAISettingsNotFound = errors.New("job phase ai settings not found")
)

// JobPhaseAISettings は JOB_PHASE_AI_SETTINGS テーブルの 1 レコードを表す。
// 3 フェーズ種別（word_translation、npc_persona_generation、text_translation）共通の汎用モデル。
// 主キーは phase_type のみ。job_id・user_id は持たない。
type JobPhaseAISettings struct {
	PhaseType     string
	AIProvider    string
	ModelName     string
	ExecutionMode string
	BatchMode     string
	UpdatedAt     time.Time
}

// JobPhaseAISettingsDraft は JOB_PHASE_AI_SETTINGS の upsert ペイロードを表す。
// 入力に job_id を含めない。保存操作はジョブと無関連で phase_type を主キーとして upsert する。
type JobPhaseAISettingsDraft struct {
	PhaseType     string
	AIProvider    string
	ModelName     string
	ExecutionMode string
	BatchMode     string
	UpdatedAt     time.Time
}

// JobPhaseAISettingsRepository は JOB_PHASE_AI_SETTINGS の永続化操作を定義する。
// 3 フェーズ共通で phase_type を受け取る汎用 interface。
type JobPhaseAISettingsRepository interface {
	// Save は draft を phase_type を主キーとして upsert する。
	// 同一 phase_type への再保存は上書きとして処理する。明示的な削除 API は持たない。
	Save(ctx context.Context, draft JobPhaseAISettingsDraft) (JobPhaseAISettings, error)

	// LoadByPhaseType は phase_type に対応する record を返す。
	// record が存在しない場合は ErrJobPhaseAISettingsNotFound を返す。
	LoadByPhaseType(ctx context.Context, phaseType string) (JobPhaseAISettings, error)
}
