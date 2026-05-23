package service

import (
	"context"
	"fmt"
	"log/slog"
	"strings"

	"aitranslationenginejp/internal/repository"
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

// PhaseAISettingsSelection carries public phase AI settings without secret material.
type PhaseAISettingsSelection struct {
	JobID         int64
	PhaseID       string
	Provider      string
	Model         string
	ExecutionMode string
	BatchMode     string
}

// PhaseAISettingsReadModel returns public phase AI settings state.
type PhaseAISettingsReadModel struct {
	JobID            int64
	PhaseID          string
	Provider         string
	Model            string
	CredentialStatus string
	ExecutionMode    string
	BatchMode        string
	ModelListStatus  string
}

type translationJobPhaseRuntimeSnapshotStore interface {
	SaveTranslationJobPhaseRuntimeSnapshot(
		ctx context.Context,
		draft repository.TranslationJobPhaseRuntimeSnapshotDraft,
	) (repository.TranslationJobPhaseRuntimeSnapshot, error)
	GetTranslationJobPhaseRuntimeSnapshot(
		ctx context.Context,
		translationJobID int64,
		phaseID string,
	) (repository.TranslationJobPhaseRuntimeSnapshot, error)
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

func providerExecutionSnapshotFromRuntimeSnapshot(
	snapshot repository.TranslationJobPhaseRuntimeSnapshot,
) providerExecutionSnapshot {
	return providerExecutionSnapshot{
		Provider:        strings.TrimSpace(snapshot.Provider),
		Model:           strings.TrimSpace(snapshot.ModelName),
		ExecutionMode:   strings.TrimSpace(snapshot.ExecutionMode),
		CredentialState: strings.TrimSpace(snapshot.CredentialStatus),
	}
}

func providerExecutionBatchMode(snapshot repository.TranslationJobPhaseRuntimeSnapshot) string {
	return strings.TrimSpace(snapshot.BatchMode)
}

func savePhaseAISettings(
	ctx context.Context,
	store translationJobPhaseRuntimeSnapshotStore,
	providerSettings ProviderSettingsConsumer,
	consumerID string,
	selection PhaseAISettingsSelection,
) (PhaseAISettingsReadModel, error) {
	phaseID := strings.TrimSpace(selection.PhaseID)
	providerID := strings.TrimSpace(selection.Provider)
	model := strings.TrimSpace(selection.Model)
	executionMode := strings.TrimSpace(selection.ExecutionMode)
	batchMode := strings.TrimSpace(selection.BatchMode)
	credentialStatus := ""
	modelListStatus := "not_updated"
	if batchMode == "" {
		batchMode = "disabled"
	}
	if providerSettings != nil && providerExecutionUsesProviderSettings(providerID) {
		resolved, err := providerSettings.ResolveProviderExecutionSettings(ctx, ProviderSettingsResolveInput{
			ConsumerID: consumerID,
			Selection: ProviderSettingsResolveSelection{
				ProviderID:      providerID,
				Model:           model,
				ExecutionMethod: executionMode,
				UseBatchAPI:     batchMode == "enabled",
			},
		})
		if err != nil {
			logPhaseAISettingsSave(ctx, "phase_ai_settings_save", "backend.service.phase_ai_settings", "failed", selection.JobID, phaseID, providerID, model, executionMode, batchMode, credentialStatus, modelListStatus, "provider_settings_resolve_failed")
			return PhaseAISettingsReadModel{}, fmt.Errorf("resolve phase ai settings: %w", err)
		}
		credentialStatus = strings.TrimSpace(resolved.CredentialState)
		if resolved.ErrorKind != nil {
			switch strings.TrimSpace(*resolved.ErrorKind) {
			case providerSettingsErrorKindCredentialMissing:
				credentialStatus = "missing"
				modelListStatus = "credential_missing"
			default:
				modelListStatus = "failed"
			}
		}
	}
	if credentialStatus == "" {
		if providerID == "lm_studio" {
			credentialStatus = "not_required"
			modelListStatus = "credential_not_required"
		} else {
			credentialStatus = "configured"
		}
	}
	if _, err := store.SaveTranslationJobPhaseRuntimeSnapshot(ctx, repository.TranslationJobPhaseRuntimeSnapshotDraft{
		TranslationJobID: selection.JobID,
		PhaseID:          phaseID,
		Provider:         providerID,
		ModelName:        model,
		CredentialStatus: credentialStatus,
		ExecutionMode:    executionMode,
		BatchMode:        batchMode,
	}); err != nil {
		logPhaseAISettingsSave(ctx, "phase_ai_settings_save", "backend.service.phase_ai_settings", "failed", selection.JobID, phaseID, providerID, model, executionMode, batchMode, credentialStatus, modelListStatus, "snapshot_save_failed")
		return PhaseAISettingsReadModel{}, fmt.Errorf("save phase ai settings: %w", err)
	}
	logPhaseAISettingsSave(ctx, "phase_ai_settings_save", "backend.service.phase_ai_settings", "accepted", selection.JobID, phaseID, providerID, model, executionMode, batchMode, credentialStatus, modelListStatus, "")
	return PhaseAISettingsReadModel{
		JobID:            selection.JobID,
		PhaseID:          phaseID,
		Provider:         providerID,
		Model:            model,
		CredentialStatus: credentialStatus,
		ExecutionMode:    executionMode,
		BatchMode:        batchMode,
		ModelListStatus:  modelListStatus,
	}, nil
}

func logPhaseAISettingsSave(
	ctx context.Context,
	event string,
	where string,
	result string,
	jobID int64,
	phaseID string,
	providerID string,
	model string,
	executionMode string,
	batchMode string,
	credentialStatus string,
	modelListStatus string,
	reason string,
) {
	attrs := []slog.Attr{
		slog.String("event", strings.TrimSpace(event)),
		slog.String("where", strings.TrimSpace(where)),
		slog.String("result", strings.TrimSpace(result)),
		slog.String("id", fmt.Sprintf("job:%d", jobID)),
		slog.String("phase", strings.TrimSpace(phaseID)),
		slog.String("provider", strings.TrimSpace(providerID)),
		slog.String("model", strings.TrimSpace(model)),
		slog.String("executionMode", strings.TrimSpace(executionMode)),
		slog.String("batchMode", strings.TrimSpace(batchMode)),
		slog.String("credentialStatus", strings.TrimSpace(credentialStatus)),
		slog.String("modelListStatus", strings.TrimSpace(modelListStatus)),
	}
	if strings.TrimSpace(reason) != "" {
		attrs = append(attrs, slog.String("reason", strings.TrimSpace(reason)))
	}
	level := slog.LevelInfo
	if result == "failed" {
		level = slog.LevelWarn
	}
	slog.LogAttrs(ctx, level, "phase ai settings save completed", attrs...)
}
