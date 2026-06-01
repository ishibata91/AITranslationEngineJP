package repository

import (
	"context"
	"fmt"
	"strings"
	"sync"
	"time"
)

// InMemoryJobPhaseAISettingsRepository は JobPhaseAISettingsRepository のインメモリ実装。
// テストおよびローカルワイヤリング用のシーム実装。
// 3 フェーズ種別共通で phase_type をキーとして保持する。最大 3 件のみ存在し得る。
type InMemoryJobPhaseAISettingsRepository struct {
	mu       sync.RWMutex
	settings map[string]JobPhaseAISettings
}

// NewInMemoryJobPhaseAISettingsRepository は InMemoryJobPhaseAISettingsRepository を返す。
func NewInMemoryJobPhaseAISettingsRepository() JobPhaseAISettingsRepository {
	return &InMemoryJobPhaseAISettingsRepository{
		settings: make(map[string]JobPhaseAISettings),
	}
}

// Save は draft を phase_type を主キーとして upsert する。
func (r *InMemoryJobPhaseAISettingsRepository) Save(
	_ context.Context,
	draft JobPhaseAISettingsDraft,
) (JobPhaseAISettings, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	record := JobPhaseAISettings{
		PhaseType:     strings.TrimSpace(draft.PhaseType),
		AIProvider:    strings.TrimSpace(draft.AIProvider),
		ModelName:     strings.TrimSpace(draft.ModelName),
		ExecutionMode: strings.TrimSpace(draft.ExecutionMode),
		BatchMode:     strings.TrimSpace(draft.BatchMode),
		UpdatedAt:     draft.UpdatedAt.UTC().Truncate(time.Millisecond),
	}
	r.settings[record.PhaseType] = record
	return record, nil
}

// LoadByPhaseType は phase_type に対応する record を返す。
// record が存在しない場合は ErrJobPhaseAISettingsNotFound を返す。
func (r *InMemoryJobPhaseAISettingsRepository) LoadByPhaseType(
	_ context.Context,
	phaseType string,
) (JobPhaseAISettings, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	key := strings.TrimSpace(phaseType)
	record, ok := r.settings[key]
	if !ok {
		return JobPhaseAISettings{}, fmt.Errorf("%w: phase_type=%s", ErrJobPhaseAISettingsNotFound, phaseType)
	}
	return record, nil
}
