package service

import (
	"context"
	"fmt"
	"strings"
)

// ProviderSettingsReader defines the backend-internal read seam for provider settings summaries and model lists.
type ProviderSettingsReader interface {
	ListProviderSettings(ctx context.Context) (ProviderSettingsRoute, []ProviderSettingsSummary, error)
	ListProviderModels(ctx context.Context, input ProviderSettingsModelListInput) (ProviderSettingsModelListResult, error)
}

// ProviderSettingsResolver defines the backend-internal execution-resolution seam for downstream consumers.
type ProviderSettingsResolver interface {
	ResolveProviderExecutionSettings(ctx context.Context, input ProviderSettingsResolveInput) (ProviderSettingsResolveResult, error)
}

// ProviderSettingsConsumer defines the backend-internal read-only provider settings seam shared by downstream consumers.
type ProviderSettingsConsumer interface {
	ProviderSettingsReader
	ProviderSettingsResolver
}

func providerSettingsSummaryMap(
	ctx context.Context,
	reader ProviderSettingsReader,
) (map[string]ProviderSettingsSummary, error) {
	if reader == nil {
		return nil, nil
	}
	_, summaries, err := reader.ListProviderSettings(ctx)
	if err != nil {
		return nil, fmt.Errorf("list provider settings for consumer: %w", err)
	}
	result := make(map[string]ProviderSettingsSummary, len(summaries))
	for _, summary := range summaries {
		providerID := strings.TrimSpace(summary.ProviderID)
		if providerID == "" {
			continue
		}
		result[providerID] = summary
	}
	return result, nil
}

func providerSettingsSummaryForProvider(
	ctx context.Context,
	reader ProviderSettingsReader,
	providerID string,
) (ProviderSettingsSummary, bool, error) {
	summaries, err := providerSettingsSummaryMap(ctx, reader)
	if err != nil {
		return ProviderSettingsSummary{}, false, err
	}
	if summaries == nil {
		return ProviderSettingsSummary{}, false, nil
	}
	summary, ok := summaries[strings.TrimSpace(providerID)]
	return summary, ok, nil
}
