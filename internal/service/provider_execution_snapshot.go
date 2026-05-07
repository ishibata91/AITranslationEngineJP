package service

import (
	"context"
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
