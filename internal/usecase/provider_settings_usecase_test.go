package usecase

import (
	"context"
	"testing"

	"aitranslationenginejp/internal/service"
)

type fakeProviderSettingsService struct {
	save func(context.Context, service.ProviderSettingsSaveInput) (service.ProviderSettingsSummary, error)
}

func (fakeProviderSettingsService) ListProviderSettings(context.Context) (service.ProviderSettingsRoute, []service.ProviderSettingsSummary, error) {
	return service.ProviderSettingsRoute{}, nil, nil
}

func (serviceStub fakeProviderSettingsService) SaveProviderSettings(
	ctx context.Context,
	input service.ProviderSettingsSaveInput,
) (service.ProviderSettingsSummary, error) {
	return serviceStub.save(ctx, input)
}

func (fakeProviderSettingsService) ResetProviderSettings(context.Context, string) (service.ProviderSettingsSummary, error) {
	return service.ProviderSettingsSummary{}, nil
}

func (fakeProviderSettingsService) ValidateProviderSettings(
	context.Context,
	service.ProviderSettingsValidateInput,
) (service.ProviderSettingsValidationResult, error) {
	return service.ProviderSettingsValidationResult{}, nil
}

func (fakeProviderSettingsService) ListProviderModels(
	context.Context,
	service.ProviderSettingsModelListInput,
) (service.ProviderSettingsModelListResult, error) {
	return service.ProviderSettingsModelListResult{}, nil
}

func (fakeProviderSettingsService) ResolveProviderExecutionSettings(
	context.Context,
	service.ProviderSettingsResolveInput,
) (service.ProviderSettingsResolveResult, error) {
	return service.ProviderSettingsResolveResult{}, nil
}

func TestProviderSettingsUsecaseSaveWithSecretForwardsSecretOnlyToService(t *testing.T) {
	var captured service.ProviderSettingsSaveInput
	usecase := NewProviderSettingsUsecase(fakeProviderSettingsService{
		save: func(_ context.Context, input service.ProviderSettingsSaveInput) (service.ProviderSettingsSummary, error) {
			captured = input
			return service.ProviderSettingsSummary{
				ProviderID:            "gemini",
				Label:                 "Gemini",
				Endpoint:              input.Endpoint,
				CredentialState:       "configured",
				CredentialReferenceID: stringPointerForProviderSettingsUsecaseTest("provider-settings:gemini"),
				ValidationState:       "not_validated",
				SavedState:            "configured",
				RequestToken:          stringPointerForProviderSettingsUsecaseTest("gemini|1"),
			}, nil
		},
	})

	result, err := usecase.SaveProviderSettingsWithSecret(context.Background(), SaveProviderSettingsRequest{
		ProviderID:         ProviderSettingsProviderIDGemini,
		Endpoint:           stringPointerForProviderSettingsUsecaseTest(" https://gemini.example/v1 "),
		APIKeyInputPresent: true,
	}, " transport-secret ")
	if err != nil {
		t.Fatalf("expected usecase save with secret to succeed: %v", err)
	}
	if captured.APIKey == nil || *captured.APIKey != "transport-secret" {
		t.Fatalf("expected secret to be forwarded only to service input, got %#v", captured)
	}
	if result.Provider.CredentialReferenceID == nil || *result.Provider.CredentialReferenceID != "provider-settings:gemini" {
		t.Fatalf("expected saved provider summary, got %#v", result)
	}
}

func stringPointerForProviderSettingsUsecaseTest(value string) *string {
	cloned := value
	return &cloned
}
