package service

import "context"

type fakePhaseProviderSettingsConsumer struct {
	resolveFunc func(context.Context, ProviderSettingsResolveInput) (ProviderSettingsResolveResult, error)
}

func (fake fakePhaseProviderSettingsConsumer) ListProviderSettings(context.Context) (ProviderSettingsRoute, []ProviderSettingsSummary, error) {
	return ProviderSettingsRoute{}, nil, nil
}

func (fake fakePhaseProviderSettingsConsumer) SaveProviderSettings(context.Context, ProviderSettingsSaveInput) (ProviderSettingsSummary, error) {
	return ProviderSettingsSummary{}, nil
}

func (fake fakePhaseProviderSettingsConsumer) ListProviderModels(context.Context, ProviderSettingsModelListInput) (ProviderSettingsModelListResult, error) {
	return ProviderSettingsModelListResult{}, nil
}

func (fake fakePhaseProviderSettingsConsumer) ResolveProviderExecutionSettings(
	ctx context.Context,
	input ProviderSettingsResolveInput,
) (ProviderSettingsResolveResult, error) {
	if fake.resolveFunc != nil {
		return fake.resolveFunc(ctx, input)
	}
	return ProviderSettingsResolveResult{}, nil
}
