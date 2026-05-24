package wails

import (
	"context"
	"errors"
	"strings"
	"testing"

	"aitranslationenginejp/internal/usecase"
)

type fakeProviderSettingsUsecase struct {
	list     func(context.Context, usecase.ListProviderSettingsRequest) (usecase.ListProviderSettingsResult, error)
	save     func(context.Context, usecase.SaveProviderSettingsRequest) (usecase.SaveProviderSettingsResult, error)
	reset    func(context.Context, usecase.ResetProviderSettingsRequest) (usecase.ResetProviderSettingsResult, error)
	validate func(context.Context, usecase.ValidateProviderSettingsRequest) (usecase.ValidateProviderSettingsResult, error)
}

func (stub fakeProviderSettingsUsecase) ListProviderSettings(
	ctx context.Context,
	request usecase.ListProviderSettingsRequest,
) (usecase.ListProviderSettingsResult, error) {
	if stub.list != nil {
		return stub.list(ctx, request)
	}
	return usecase.ListProviderSettingsResult{}, nil
}

func (stub fakeProviderSettingsUsecase) SaveProviderSettings(
	ctx context.Context,
	request usecase.SaveProviderSettingsRequest,
) (usecase.SaveProviderSettingsResult, error) {
	if stub.save != nil {
		return stub.save(ctx, request)
	}
	return usecase.SaveProviderSettingsResult{}, nil
}

func (stub fakeProviderSettingsUsecase) ResetProviderSettings(
	ctx context.Context,
	request usecase.ResetProviderSettingsRequest,
) (usecase.ResetProviderSettingsResult, error) {
	if stub.reset != nil {
		return stub.reset(ctx, request)
	}
	return usecase.ResetProviderSettingsResult{}, nil
}

func (stub fakeProviderSettingsUsecase) ValidateProviderSettings(
	ctx context.Context,
	request usecase.ValidateProviderSettingsRequest,
) (usecase.ValidateProviderSettingsResult, error) {
	if stub.validate != nil {
		return stub.validate(ctx, request)
	}
	return usecase.ValidateProviderSettingsResult{}, nil
}

func TestProviderSettingsControllerListProviderSettingsMapsResponseDTO(t *testing.T) {
	// ListProviderSettings の公開応答で route と provider の DTO 写像を証明する。
	controller := NewProviderSettingsController(fakeProviderSettingsUsecase{
		list: func(context.Context, usecase.ListProviderSettingsRequest) (usecase.ListProviderSettingsResult, error) {
			failureKind := usecase.ProviderSettingsErrorKindValidationFailed
			return usecase.ListProviderSettingsResult{
				Route: usecase.ProviderSettingsRoute{
					RouteID:           "provider-settings",
					Label:             "AIサービス設定",
					CurrentRouteState: "active",
					DashboardEntryID:  "ai-settings",
				},
				Providers: []usecase.ProviderSettingsSummary{{
					ProviderID:            usecase.ProviderSettingsProviderIDXAI,
					Label:                 "xAI",
					Endpoint:              stringPointerForProviderSettingsControllerTest(" https://api.x.ai/v1 "),
					CredentialState:       usecase.ProviderSettingsCredentialStateConfigured,
					CredentialReferenceID: stringPointerForProviderSettingsControllerTest(" ref-01 "),
					ValidationState:       usecase.ProviderSettingsValidationStateFailed,
					SavedState:            usecase.ProviderSettingsSavedStateConfigured,
					RequestToken:          stringPointerForProviderSettingsControllerTest(" token-1 "),
					LastFailureKind:       &failureKind,
				}},
			}, nil
		},
	})

	response, err := controller.ListProviderSettings(ListProviderSettingsRequestDTO{})
	if err != nil {
		t.Fatalf("expected provider settings list to succeed: %v", err)
	}
	if response.Route.RouteID != "provider-settings" || response.Route.DashboardEntryID != "ai-settings" {
		t.Fatalf("expected route mapping, got %#v", response.Route)
	}
	if len(response.Providers) != 1 {
		t.Fatalf("expected one provider summary, got %#v", response.Providers)
	}
	provider := response.Providers[0]
	if provider.ProviderID != "xai" || provider.ValidationState != "failed" {
		t.Fatalf("expected provider id and validation state mapping, got %#v", provider)
	}
	if provider.Endpoint == nil || *provider.Endpoint != "https://api.x.ai/v1" {
		t.Fatalf("expected trimmed endpoint mapping, got %#v", provider)
	}
	if provider.CredentialReferenceID == nil || *provider.CredentialReferenceID != "ref-01" {
		t.Fatalf("expected trimmed credential reference id mapping, got %#v", provider)
	}
	if provider.LastFailureKind == nil || *provider.LastFailureKind != "validation_failed" {
		t.Fatalf("expected normalized last failure kind mapping, got %#v", provider)
	}
}

func TestProviderSettingsControllerListProviderSettingsWrapsError(t *testing.T) {
	// ListProviderSettings の失敗時に public method 名が分かる wrap を証明する。
	controller := NewProviderSettingsController(fakeProviderSettingsUsecase{
		list: func(context.Context, usecase.ListProviderSettingsRequest) (usecase.ListProviderSettingsResult, error) {
			return usecase.ListProviderSettingsResult{}, errors.New("backend down")
		},
	})

	_, err := controller.ListProviderSettings(ListProviderSettingsRequestDTO{})
	if err == nil {
		t.Fatal("expected error")
	}
	if !strings.Contains(err.Error(), "list provider settings") {
		t.Fatalf("expected wrapped method context, got %v", err)
	}
}

func TestProviderSettingsControllerSaveTrimsEndpointAndProviderID(t *testing.T) {
	// SaveProviderSettings の request DTO 正規化を証明する。
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

func TestProviderSettingsControllerResetProviderSettingsMapsRequestAndResponse(t *testing.T) {
	// ResetProviderSettings の request/response DTO 写像を証明する。
	var captured usecase.ResetProviderSettingsRequest
	controller := NewProviderSettingsController(fakeProviderSettingsUsecase{
		reset: func(_ context.Context, request usecase.ResetProviderSettingsRequest) (usecase.ResetProviderSettingsResult, error) {
			captured = request
			return usecase.ResetProviderSettingsResult{
				Provider: usecase.ProviderSettingsSummary{
					ProviderID:      request.ProviderID,
					Label:           "Gemini",
					CredentialState: usecase.ProviderSettingsCredentialStateMissing,
					ValidationState: usecase.ProviderSettingsValidationStateNotValidated,
					SavedState:      usecase.ProviderSettingsSavedStateNotSaved,
				},
			}, nil
		},
	})

	response, err := controller.ResetProviderSettings(ResetProviderSettingsRequestDTO{ProviderID: " gemini "})
	if err != nil {
		t.Fatalf("expected reset to succeed: %v", err)
	}
	if captured.ProviderID != usecase.ProviderSettingsProviderIDGemini {
		t.Fatalf("expected trimmed provider id, got %#v", captured)
	}
	if response.Provider.ProviderID != "gemini" || response.Provider.SavedState != "not_saved" {
		t.Fatalf("expected response mapping, got %#v", response)
	}
}

func TestProviderSettingsControllerResetProviderSettingsWrapsError(t *testing.T) {
	// ResetProviderSettings の失敗時に method 境界の wrap を証明する。
	controller := NewProviderSettingsController(fakeProviderSettingsUsecase{
		reset: func(context.Context, usecase.ResetProviderSettingsRequest) (usecase.ResetProviderSettingsResult, error) {
			return usecase.ResetProviderSettingsResult{}, errors.New("store unavailable")
		},
	})

	_, err := controller.ResetProviderSettings(ResetProviderSettingsRequestDTO{ProviderID: "gemini"})
	if err == nil {
		t.Fatal("expected error")
	}
	if !strings.Contains(err.Error(), "reset provider settings") {
		t.Fatalf("expected wrapped method context, got %v", err)
	}
}

func TestProviderSettingsControllerValidateProviderSettingsMapsRequestAndResponse(t *testing.T) {
	// ValidateProviderSettings の request/response DTO 写像を証明する。
	var captured usecase.ValidateProviderSettingsRequest
	controller := NewProviderSettingsController(fakeProviderSettingsUsecase{
		validate: func(_ context.Context, request usecase.ValidateProviderSettingsRequest) (usecase.ValidateProviderSettingsResult, error) {
			captured = request
			failureKind := usecase.ProviderSettingsErrorKindCredentialMissing
			return usecase.ValidateProviderSettingsResult{
				ProviderID:      request.ProviderID,
				ValidationState: usecase.ProviderSettingsValidationStateFailed,
				RequestToken:    request.RequestToken,
				FailureKind:     &failureKind,
			}, nil
		},
	})

	response, err := controller.ValidateProviderSettings(ValidateProviderSettingsRequestDTO{
		ProviderID:            " xai ",
		Endpoint:              stringPointerForProviderSettingsControllerTest(" https://api.x.ai/v1 "),
		CredentialState:       " configured ",
		CredentialReferenceID: stringPointerForProviderSettingsControllerTest(" cred-9 "),
		RequestToken:          " token-9 ",
	})
	if err != nil {
		t.Fatalf("expected validation to succeed: %v", err)
	}
	if captured.ProviderID != usecase.ProviderSettingsProviderIDXAI {
		t.Fatalf("expected trimmed provider id, got %#v", captured)
	}
	if captured.CredentialState != usecase.ProviderSettingsCredentialStateConfigured {
		t.Fatalf("expected trimmed credential state, got %#v", captured)
	}
	if captured.CredentialReferenceID == nil || *captured.CredentialReferenceID != "cred-9" {
		t.Fatalf("expected trimmed credential reference id, got %#v", captured)
	}
	if captured.RequestToken != "token-9" {
		t.Fatalf("expected trimmed request token, got %#v", captured)
	}
	if response.FailureKind == nil || *response.FailureKind != "credential_missing" {
		t.Fatalf("expected failure kind mapping, got %#v", response)
	}
}

func TestProviderSettingsControllerValidateProviderSettingsWrapsError(t *testing.T) {
	// ValidateProviderSettings の失敗時に method 境界の wrap を証明する。
	controller := NewProviderSettingsController(fakeProviderSettingsUsecase{
		validate: func(context.Context, usecase.ValidateProviderSettingsRequest) (usecase.ValidateProviderSettingsResult, error) {
			return usecase.ValidateProviderSettingsResult{}, errors.New("provider timeout")
		},
	})

	_, err := controller.ValidateProviderSettings(ValidateProviderSettingsRequestDTO{ProviderID: "xai", RequestToken: "token"})
	if err == nil {
		t.Fatal("expected error")
	}
	if !strings.Contains(err.Error(), "validate provider settings") {
		t.Fatalf("expected wrapped method context, got %v", err)
	}
}

func stringPointerForProviderSettingsControllerTest(value string) *string {
	cloned := value
	return &cloned
}
