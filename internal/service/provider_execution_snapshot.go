package service

import (
	"context"
	"fmt"
	"log/slog"
	"strings"
	"time"

	"aitranslationenginejp/internal/repository"
)

const (
	phaseAISettingsLogEvent        = "phase_ai_settings_save"
	phaseAISettingsLogWhere        = "backend.service.phase_ai_settings"
	phaseAISettingsLogResultFailed = "failed"
)

type providerExecutionSnapshot struct {
	Provider        string
	Model           string
	ExecutionMode   string
	CredentialRef   string
	CredentialState string
	EndpointSummary *string
	RequestToken    *string
}

// PhaseAISettingsSelection carries the phase AI settings to save.
// PhaseType は "word_translation"/"npc_persona_generation"/"text_translation" のいずれか。
// 入力に job_id を含めない。保存操作はジョブと無関連で phase_type を主キーとして upsert する。
type PhaseAISettingsSelection struct {
	PhaseType     string
	Provider      string
	Model         string
	ExecutionMode string
	BatchMode     string
}

// PhaseAISettingsReadModel returns public phase AI settings state after save.
// PhaseType は保存対象の phase_type を示す。
type PhaseAISettingsReadModel struct {
	PhaseType     string
	Provider      string
	Model         string
	ExecutionMode string
	BatchMode     string
	UpdatedAt     time.Time
}

func providerExecutionUsesProviderSettings(providerID string) bool {
	switch strings.ToLower(strings.TrimSpace(providerID)) {
	case "gemini", "lm_studio", "xai":
		return true
	default:
		return false
	}
}

func providerExecutionEndpointSummary(endpoint *string) *string {
	if endpoint == nil {
		return nil
	}
	trimmed := strings.TrimSpace(*endpoint)
	if trimmed == "" {
		return nil
	}
	cloned := trimmed
	return &cloned
}

func providerExecutionOptionalString(value *string) string {
	if value == nil {
		return ""
	}
	return strings.TrimSpace(*value)
}

// savePhaseAISettings は JobPhaseAISettingsRepository へ AI 設定を upsert する。
// 入力の PhaseType を主キーとして保存する。credential_ref は保存しない。
func savePhaseAISettings(
	ctx context.Context,
	repo repository.JobPhaseAISettingsRepository,
	consumerID string,
	selection PhaseAISettingsSelection,
) (PhaseAISettingsReadModel, error) {
	phaseType := strings.TrimSpace(selection.PhaseType)
	providerID := strings.TrimSpace(selection.Provider)
	model := strings.TrimSpace(selection.Model)
	executionMode := strings.TrimSpace(selection.ExecutionMode)
	batchMode := strings.TrimSpace(selection.BatchMode)
	if batchMode == "" {
		batchMode = "disabled"
	}
	now := time.Now().UTC()
	saved, err := repo.Save(ctx, repository.JobPhaseAISettingsDraft{
		PhaseType:     phaseType,
		AIProvider:    providerID,
		ModelName:     model,
		ExecutionMode: executionMode,
		BatchMode:     batchMode,
		UpdatedAt:     now,
	})
	if err != nil {
		logPhaseAISettingsSave(ctx, phaseAISettingsLogEvent, phaseAISettingsLogWhere, phaseAISettingsLogResultFailed, consumerID, phaseType, providerID, model, executionMode, batchMode, "save_failed")
		return PhaseAISettingsReadModel{}, fmt.Errorf("save phase ai settings: %w", err)
	}
	logPhaseAISettingsSave(ctx, phaseAISettingsLogEvent, phaseAISettingsLogWhere, "accepted", consumerID, phaseType, providerID, model, executionMode, batchMode, "")
	return PhaseAISettingsReadModel{
		PhaseType:     saved.PhaseType,
		Provider:      saved.AIProvider,
		Model:         saved.ModelName,
		ExecutionMode: saved.ExecutionMode,
		BatchMode:     saved.BatchMode,
		UpdatedAt:     saved.UpdatedAt,
	}, nil
}

func logPhaseAISettingsSave(
	ctx context.Context,
	event string,
	where string,
	result string,
	consumerID string,
	phaseType string,
	providerID string,
	model string,
	executionMode string,
	batchMode string,
	reason string,
) {
	attrs := []slog.Attr{
		slog.String("event", strings.TrimSpace(event)),
		slog.String("where", strings.TrimSpace(where)),
		slog.String("result", strings.TrimSpace(result)),
		slog.String("consumer", strings.TrimSpace(consumerID)),
		slog.String("phase_type", strings.TrimSpace(phaseType)),
		slog.String("provider", strings.TrimSpace(providerID)),
		slog.String("model", strings.TrimSpace(model)),
		slog.String("executionMode", strings.TrimSpace(executionMode)),
		slog.String("batchMode", strings.TrimSpace(batchMode)),
	}
	if strings.TrimSpace(reason) != "" {
		attrs = append(attrs, slog.String("reason", strings.TrimSpace(reason)))
	}
	level := slog.LevelInfo
	if result == phaseAISettingsLogResultFailed {
		level = slog.LevelWarn
	}
	slog.LogAttrs(ctx, level, "phase ai settings save completed", attrs...)
}
