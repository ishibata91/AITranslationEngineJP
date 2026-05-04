package usecase

import (
	"context"
	"fmt"
	"strings"

	"aitranslationenginejp/internal/service"
)

// ProviderSettingsServicePort defines the service dependency used by the provider settings usecase.
type ProviderSettingsServicePort interface {
	ListProviderSettings(ctx context.Context) (service.ProviderSettingsRoute, []service.ProviderSettingsSummary, error)
	SaveProviderSettings(ctx context.Context, input service.ProviderSettingsSaveInput) (service.ProviderSettingsSummary, error)
	ResetProviderSettings(ctx context.Context, providerID string) (service.ProviderSettingsSummary, error)
	ValidateProviderSettings(ctx context.Context, input service.ProviderSettingsValidateInput) (service.ProviderSettingsValidationResult, error)
	ListProviderModels(ctx context.Context, input service.ProviderSettingsModelListInput) (service.ProviderSettingsModelListResult, error)
	ResolveProviderExecutionSettings(ctx context.Context, input service.ProviderSettingsResolveInput) (service.ProviderSettingsResolveResult, error)
}

// ProviderSettingsAppUsecase implements the frozen provider settings usecase contract.
type ProviderSettingsAppUsecase struct {
	service ProviderSettingsServicePort
}

// NewProviderSettingsUsecase creates a provider settings usecase implementation.
func NewProviderSettingsUsecase(service ProviderSettingsServicePort) *ProviderSettingsAppUsecase {
	return &ProviderSettingsAppUsecase{service: service}
}

// ListProviderSettings returns the frozen provider settings read model.
func (usecase *ProviderSettingsAppUsecase) ListProviderSettings(
	ctx context.Context,
	_ ListProviderSettingsRequest,
) (ListProviderSettingsResult, error) {
	route, providers, err := usecase.service.ListProviderSettings(ctx)
	if err != nil {
		return ListProviderSettingsResult{}, fmt.Errorf("list provider settings usecase: %w", err)
	}
	result := ListProviderSettingsResult{
		Route: ProviderSettingsRoute{
			RouteID:           route.RouteID,
			Label:             route.Label,
			CurrentRouteState: route.CurrentRouteState,
			DashboardEntryID:  route.DashboardEntryID,
		},
		Providers: make([]ProviderSettingsSummary, 0, len(providers)),
	}
	for _, provider := range providers {
		result.Providers = append(result.Providers, toProviderSettingsSummary(provider))
	}
	return result, nil
}

// SaveProviderSettings stores one provider settings row without exposing secrets through the public contract.
func (usecase *ProviderSettingsAppUsecase) SaveProviderSettings(
	ctx context.Context,
	request SaveProviderSettingsRequest,
) (SaveProviderSettingsResult, error) {
	saved, err := usecase.service.SaveProviderSettings(ctx, service.ProviderSettingsSaveInput{
		ProviderID:         string(request.ProviderID),
		Endpoint:           cloneUsecaseOptionalString(request.Endpoint),
		APIKeyInputPresent: request.APIKeyInputPresent,
	})
	if err != nil {
		return SaveProviderSettingsResult{}, fmt.Errorf("save provider settings usecase: %w", err)
	}
	return SaveProviderSettingsResult{Provider: toProviderSettingsSummary(saved)}, nil
}

// SaveProviderSettingsWithSecret stores one provider settings row and secret for backend-internal callers.
func (usecase *ProviderSettingsAppUsecase) SaveProviderSettingsWithSecret(
	ctx context.Context,
	request SaveProviderSettingsRequest,
	apiKey string,
) (SaveProviderSettingsResult, error) {
	saved, err := usecase.service.SaveProviderSettings(ctx, service.ProviderSettingsSaveInput{
		ProviderID:         string(request.ProviderID),
		Endpoint:           cloneUsecaseOptionalString(request.Endpoint),
		APIKeyInputPresent: request.APIKeyInputPresent,
		APIKey:             cloneUsecaseOptionalString(stringPointerUsecase(apiKey)),
	})
	if err != nil {
		return SaveProviderSettingsResult{}, fmt.Errorf("save provider settings with secret usecase: %w", err)
	}
	return SaveProviderSettingsResult{Provider: toProviderSettingsSummary(saved)}, nil
}

// ResetProviderSettings resets one provider settings row.
func (usecase *ProviderSettingsAppUsecase) ResetProviderSettings(
	ctx context.Context,
	request ResetProviderSettingsRequest,
) (ResetProviderSettingsResult, error) {
	reset, err := usecase.service.ResetProviderSettings(ctx, string(request.ProviderID))
	if err != nil {
		return ResetProviderSettingsResult{}, fmt.Errorf("reset provider settings usecase: %w", err)
	}
	return ResetProviderSettingsResult{Provider: toProviderSettingsSummary(reset)}, nil
}

// ValidateProviderSettings validates one provider settings snapshot.
func (usecase *ProviderSettingsAppUsecase) ValidateProviderSettings(
	ctx context.Context,
	request ValidateProviderSettingsRequest,
) (ValidateProviderSettingsResult, error) {
	validated, err := usecase.service.ValidateProviderSettings(ctx, service.ProviderSettingsValidateInput{
		ProviderID:            string(request.ProviderID),
		Endpoint:              cloneUsecaseOptionalString(request.Endpoint),
		CredentialState:       string(request.CredentialState),
		CredentialReferenceID: cloneUsecaseOptionalString(request.CredentialReferenceID),
		RequestToken:          strings.TrimSpace(request.RequestToken),
	})
	if err != nil {
		return ValidateProviderSettingsResult{}, fmt.Errorf("validate provider settings usecase: %w", err)
	}
	return ValidateProviderSettingsResult{
		ProviderID:      ProviderSettingsProviderID(validated.ProviderID),
		ValidationState: ProviderSettingsValidationState(validated.ValidationState),
		RequestToken:    validated.RequestToken,
		FailureKind:     toUsecaseErrorKind(validated.FailureKind),
	}, nil
}

// ListProviderModels lists provider models for one settings snapshot.
func (usecase *ProviderSettingsAppUsecase) ListProviderModels(
	ctx context.Context,
	request ListProviderModelsRequest,
) (ListProviderModelsResult, error) {
	listed, err := usecase.service.ListProviderModels(ctx, service.ProviderSettingsModelListInput{
		ProviderID:            string(request.ProviderID),
		Endpoint:              cloneUsecaseOptionalString(request.Endpoint),
		CredentialState:       string(request.CredentialState),
		CredentialReferenceID: cloneUsecaseOptionalString(request.CredentialReferenceID),
		RequestToken:          strings.TrimSpace(request.RequestToken),
	})
	if err != nil {
		return ListProviderModelsResult{}, fmt.Errorf("list provider models usecase: %w", err)
	}
	result := ListProviderModelsResult{
		ProviderID:      ProviderSettingsProviderID(listed.ProviderID),
		Endpoint:        cloneUsecaseOptionalString(listed.Endpoint),
		CredentialState: ProviderSettingsCredentialState(listed.CredentialState),
		RequestToken:    listed.RequestToken,
		State:           ProviderSettingsModelListState(listed.State),
		Models:          make([]ProviderModelOption, 0, len(listed.Models)),
		FailureKind:     toUsecaseErrorKind(listed.FailureKind),
	}
	for _, model := range listed.Models {
		result.Models = append(result.Models, ProviderModelOption{
			ModelID: model.ModelID,
			Label:   model.Label,
		})
	}
	return result, nil
}

// ResolveProviderExecutionSettings resolves provider settings for downstream consumers.
func (usecase *ProviderSettingsAppUsecase) ResolveProviderExecutionSettings(
	ctx context.Context,
	request ResolveProviderExecutionSettingsRequest,
) (ResolveProviderExecutionSettingsResult, error) {
	resolved, err := usecase.service.ResolveProviderExecutionSettings(ctx, service.ProviderSettingsResolveInput{
		ConsumerID: request.ConsumerID,
		Selection: service.ProviderSettingsResolveSelection{
			ProviderID:      string(request.Selection.ProviderID),
			Model:           request.Selection.Model,
			ExecutionMethod: request.Selection.ExecutionMethod,
			UseBatchAPI:     request.Selection.UseBatchAPI,
		},
	})
	if err != nil {
		return ResolveProviderExecutionSettingsResult{}, fmt.Errorf("resolve provider execution settings usecase: %w", err)
	}
	return ResolveProviderExecutionSettingsResult{
		ConsumerID:            resolved.ConsumerID,
		ProviderID:            ProviderSettingsProviderID(resolved.ProviderID),
		Model:                 resolved.Model,
		ExecutionMethod:       resolved.ExecutionMethod,
		UseBatchAPI:           resolved.UseBatchAPI,
		Endpoint:              cloneUsecaseOptionalString(resolved.Endpoint),
		CredentialReferenceID: cloneUsecaseOptionalString(resolved.CredentialReferenceID),
		CredentialState:       ProviderSettingsCredentialState(resolved.CredentialState),
		RequestToken:          cloneUsecaseOptionalString(resolved.RequestToken),
		ErrorKind:             toUsecaseErrorKind(resolved.ErrorKind),
	}, nil
}

func toProviderSettingsSummary(summary service.ProviderSettingsSummary) ProviderSettingsSummary {
	return ProviderSettingsSummary{
		ProviderID:            ProviderSettingsProviderID(summary.ProviderID),
		Label:                 summary.Label,
		Endpoint:              cloneUsecaseOptionalString(summary.Endpoint),
		CredentialState:       ProviderSettingsCredentialState(summary.CredentialState),
		CredentialReferenceID: cloneUsecaseOptionalString(summary.CredentialReferenceID),
		ValidationState:       ProviderSettingsValidationState(summary.ValidationState),
		SavedState:            ProviderSettingsSavedState(summary.SavedState),
		RequestToken:          cloneUsecaseOptionalString(summary.RequestToken),
		LastFailureKind:       toUsecaseErrorKind(summary.LastFailureKind),
	}
}

func cloneUsecaseOptionalString(value *string) *string {
	if value == nil {
		return nil
	}
	trimmed := strings.TrimSpace(*value)
	if trimmed == "" {
		return nil
	}
	cloned := trimmed
	return &cloned
}

func toUsecaseErrorKind(value *string) *ProviderSettingsErrorKind {
	if value == nil {
		return nil
	}
	kind := ProviderSettingsErrorKind(strings.TrimSpace(*value))
	if kind == "" {
		return nil
	}
	return &kind
}

func stringPointerUsecase(value string) *string {
	trimmed := strings.TrimSpace(value)
	if trimmed == "" {
		return nil
	}
	cloned := trimmed
	return &cloned
}
