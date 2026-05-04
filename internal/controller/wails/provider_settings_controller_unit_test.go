package wails

import (
	"context"
	"testing"

	"aitranslationenginejp/internal/usecase"
)

type fakeProviderSettingsUsecase struct {
	save func(context.Context, usecase.SaveProviderSettingsRequest) (usecase.SaveProviderSettingsResult, error)
}

func (fakeProviderSettingsUsecase) ListProviderSettings(context.Context, usecase.ListProviderSettingsRequest) (usecase.ListProviderSettingsResult, error) {
	return usecase.ListProviderSettingsResult{}, nil
}

func (stub fakeProviderSettingsUsecase) SaveProviderSettings(
	ctx context.Context,
	request usecase.SaveProviderSettingsRequest,
) (usecase.SaveProviderSettingsResult, error) {
	return stub.save(ctx, request)
}

func (fakeProviderSettingsUsecase) ResetProviderSettings(context.Context, usecase.ResetProviderSettingsRequest) (usecase.ResetProviderSettingsResult, error) {
	return usecase.ResetProviderSettingsResult{}, nil
}

func (fakeProviderSettingsUsecase) ValidateProviderSettings(context.Context, usecase.ValidateProviderSettingsRequest) (usecase.ValidateProviderSettingsResult, error) {
	return usecase.ValidateProviderSettingsResult{}, nil
}

func TestProviderSettingsControllerSaveTrimsEndpointAndProviderID(t *testing.T) {
	var captured usecase.SaveProviderSettingsRequest
	controller := NewProviderSettingsController(fakeProviderSettingsUsecase{
		save: func(_ context.Context, request usecase.SaveProviderSettingsRequest) (usecase.SaveProviderSettingsResult, error) {
			captured = request
			return usecase.SaveProviderSettingsResult{
				Provider: usecase.ProviderSettingsSummary{
					ProviderID:      usecase.ProviderSettingsProviderIDGemini,
					Label:           "Gemini",
					Endpoint:        request.Endpoint,
					CredentialState: usecase.ProviderSettingsCredentialStateConfigured,
					ValidationState: usecase.ProviderSettingsValidationStateNotValidated,
					SavedState:      usecase.ProviderSettingsSavedStateConfigured,
				},
			}, nil
		},
	})

	response, err := controller.SaveProviderSettings(SaveProviderSettingsRequestDTO{
		ProviderID:         " gemini ",
		Endpoint:           stringPointerForProviderSettingsControllerTest(" https://gemini.example/v1 "),
		APIKeyInputPresent: false,
	})
	if err != nil {
		t.Fatalf("expected provider settings controller save to succeed: %v", err)
	}
	if captured.ProviderID != usecase.ProviderSettingsProviderIDGemini {
		t.Fatalf("expected trimmed provider id, got %#v", captured)
	}
	if captured.Endpoint == nil || *captured.Endpoint != "https://gemini.example/v1" {
		t.Fatalf("expected trimmed endpoint, got %#v", captured)
	}
	if response.Provider.Endpoint == nil || *response.Provider.Endpoint != "https://gemini.example/v1" {
		t.Fatalf("expected response endpoint, got %#v", response)
	}
}

func stringPointerForProviderSettingsControllerTest(value string) *string {
	cloned := value
	return &cloned
}
