package wails

import (
	"context"
	"fmt"
	"strings"

	"aitranslationenginejp/internal/usecase"
)

// ProviderSettingsUsecasePort defines the frozen provider settings usecase seam for Wails.
type ProviderSettingsUsecasePort interface {
	ListProviderSettings(ctx context.Context, request usecase.ListProviderSettingsRequest) (usecase.ListProviderSettingsResult, error)
	SaveProviderSettings(ctx context.Context, request usecase.SaveProviderSettingsRequest) (usecase.SaveProviderSettingsResult, error)
	ResetProviderSettings(ctx context.Context, request usecase.ResetProviderSettingsRequest) (usecase.ResetProviderSettingsResult, error)
	ValidateProviderSettings(ctx context.Context, request usecase.ValidateProviderSettingsRequest) (usecase.ValidateProviderSettingsResult, error)
}

type providerSettingsSecretAwareUsecasePort interface {
	SaveProviderSettingsWithSecret(
		ctx context.Context,
		request usecase.SaveProviderSettingsRequest,
		apiKey string,
	) (usecase.SaveProviderSettingsResult, error)
}

// ProviderSettingsController exposes Wails-bound provider settings entrypoints.
type ProviderSettingsController struct {
	providerSettingsUsecase ProviderSettingsUsecasePort
}

// ProviderSettingsRouteDTO returns the fixed AppShell route contract.
type ProviderSettingsRouteDTO struct {
	RouteID           string `json:"routeId"`
	Label             string `json:"label"`
	CurrentRouteState string `json:"currentRouteState"`
	DashboardEntryID  string `json:"dashboardEntryId"`
}

// ProviderSettingsSummaryDTO returns one provider row visible to the frontend.
type ProviderSettingsSummaryDTO struct {
	ProviderID            string  `json:"providerId"`
	Label                 string  `json:"label"`
	Endpoint              *string `json:"endpoint,omitempty"`
	CredentialState       string  `json:"credentialState"`
	CredentialReferenceID *string `json:"credentialReferenceId,omitempty"`
	ValidationState       string  `json:"validationState"`
	SavedState            string  `json:"savedState"`
	RequestToken          *string `json:"requestToken,omitempty"`
	LastFailureKind       *string `json:"lastFailureKind,omitempty"`
}

// ListProviderSettingsRequestDTO identifies a provider settings list request.
type ListProviderSettingsRequestDTO struct{}

// ListProviderSettingsResponseDTO returns the frozen provider settings read model.
type ListProviderSettingsResponseDTO struct {
	Route     ProviderSettingsRouteDTO     `json:"route"`
	Providers []ProviderSettingsSummaryDTO `json:"providers"`
}

// SaveProviderSettingsRequestDTO carries only reference values and transient credential input presence.
type SaveProviderSettingsRequestDTO struct {
	ProviderID         string  `json:"providerId"`
	Endpoint           *string `json:"endpoint,omitempty"`
	APIKeyInputPresent bool    `json:"apiKeyInputPresent"`
	CredentialInput    *string `json:"credentialInput,omitempty"`
}

// SaveProviderSettingsResponseDTO returns one saved provider summary.
type SaveProviderSettingsResponseDTO struct {
	Provider ProviderSettingsSummaryDTO `json:"provider"`
}

// ResetProviderSettingsRequestDTO identifies the provider settings row to reset.
type ResetProviderSettingsRequestDTO struct {
	ProviderID string `json:"providerId"`
}

// ResetProviderSettingsResponseDTO returns one reset provider summary.
type ResetProviderSettingsResponseDTO struct {
	Provider ProviderSettingsSummaryDTO `json:"provider"`
}

// ValidateProviderSettingsRequestDTO carries the current validation correlation data.
type ValidateProviderSettingsRequestDTO struct {
	ProviderID            string  `json:"providerId"`
	Endpoint              *string `json:"endpoint,omitempty"`
	CredentialState       string  `json:"credentialState"`
	CredentialReferenceID *string `json:"credentialReferenceId,omitempty"`
	RequestToken          string  `json:"requestToken"`
}

// ValidateProviderSettingsResponseDTO returns the redacted validation outcome.
type ValidateProviderSettingsResponseDTO struct {
	ProviderID      string  `json:"providerId"`
	ValidationState string  `json:"validationState"`
	RequestToken    string  `json:"requestToken"`
	FailureKind     *string `json:"failureKind,omitempty"`
}

// NewProviderSettingsController creates a provider settings Wails controller.
func NewProviderSettingsController(providerSettingsUsecase ProviderSettingsUsecasePort) *ProviderSettingsController {
	return &ProviderSettingsController{providerSettingsUsecase: providerSettingsUsecase}
}

// ListProviderSettings returns the frozen provider settings read model.
func (controller *ProviderSettingsController) ListProviderSettings(_ ListProviderSettingsRequestDTO) (ListProviderSettingsResponseDTO, error) {
	result, err := controller.providerSettingsUsecase.ListProviderSettings(context.Background(), usecase.ListProviderSettingsRequest{})
	if err != nil {
		return ListProviderSettingsResponseDTO{}, fmt.Errorf("list provider settings: %w", err)
	}
	return toListProviderSettingsResponseDTO(result), nil
}

// SaveProviderSettings stores provider settings reference values.
func (controller *ProviderSettingsController) SaveProviderSettings(request SaveProviderSettingsRequestDTO) (SaveProviderSettingsResponseDTO, error) {
	saveRequest := usecase.SaveProviderSettingsRequest{
		ProviderID:         usecase.ProviderSettingsProviderID(strings.TrimSpace(request.ProviderID)),
		Endpoint:           trimOptionalString(request.Endpoint),
		APIKeyInputPresent: request.APIKeyInputPresent || trimOptionalString(request.CredentialInput) != nil,
	}
	var (
		result usecase.SaveProviderSettingsResult
		err    error
	)
	if credentialInput := trimOptionalString(request.CredentialInput); credentialInput != nil {
		secretAwareUsecase, ok := controller.providerSettingsUsecase.(providerSettingsSecretAwareUsecasePort)
		if !ok {
			return SaveProviderSettingsResponseDTO{}, fmt.Errorf("save provider settings: secret-aware usecase is required")
		}
		result, err = secretAwareUsecase.SaveProviderSettingsWithSecret(context.Background(), saveRequest, *credentialInput)
	} else {
		result, err = controller.providerSettingsUsecase.SaveProviderSettings(context.Background(), saveRequest)
	}
	if err != nil {
		return SaveProviderSettingsResponseDTO{}, fmt.Errorf("save provider settings: %w", err)
	}
	return SaveProviderSettingsResponseDTO{Provider: toProviderSettingsSummaryDTO(result.Provider)}, nil
}

// ResetProviderSettings resets the provider settings row while keeping the row identity.
func (controller *ProviderSettingsController) ResetProviderSettings(request ResetProviderSettingsRequestDTO) (ResetProviderSettingsResponseDTO, error) {
	result, err := controller.providerSettingsUsecase.ResetProviderSettings(context.Background(), usecase.ResetProviderSettingsRequest{
		ProviderID: usecase.ProviderSettingsProviderID(strings.TrimSpace(request.ProviderID)),
	})
	if err != nil {
		return ResetProviderSettingsResponseDTO{}, fmt.Errorf("reset provider settings: %w", err)
	}
	return ResetProviderSettingsResponseDTO{Provider: toProviderSettingsSummaryDTO(result.Provider)}, nil
}

// ValidateProviderSettings validates one provider settings snapshot.
func (controller *ProviderSettingsController) ValidateProviderSettings(request ValidateProviderSettingsRequestDTO) (ValidateProviderSettingsResponseDTO, error) {
	result, err := controller.providerSettingsUsecase.ValidateProviderSettings(context.Background(), usecase.ValidateProviderSettingsRequest{
		ProviderID:            usecase.ProviderSettingsProviderID(strings.TrimSpace(request.ProviderID)),
		Endpoint:              trimOptionalString(request.Endpoint),
		CredentialState:       usecase.ProviderSettingsCredentialState(strings.TrimSpace(request.CredentialState)),
		CredentialReferenceID: trimOptionalString(request.CredentialReferenceID),
		RequestToken:          strings.TrimSpace(request.RequestToken),
	})
	if err != nil {
		return ValidateProviderSettingsResponseDTO{}, fmt.Errorf("validate provider settings: %w", err)
	}
	return toValidateProviderSettingsResponseDTO(result), nil
}

func toListProviderSettingsResponseDTO(result usecase.ListProviderSettingsResult) ListProviderSettingsResponseDTO {
	providers := make([]ProviderSettingsSummaryDTO, 0, len(result.Providers))
	for _, provider := range result.Providers {
		providers = append(providers, toProviderSettingsSummaryDTO(provider))
	}
	return ListProviderSettingsResponseDTO{
		Route: ProviderSettingsRouteDTO{
			RouteID:           result.Route.RouteID,
			Label:             result.Route.Label,
			CurrentRouteState: result.Route.CurrentRouteState,
			DashboardEntryID:  result.Route.DashboardEntryID,
		},
		Providers: providers,
	}
}

func toProviderSettingsSummaryDTO(summary usecase.ProviderSettingsSummary) ProviderSettingsSummaryDTO {
	return ProviderSettingsSummaryDTO{
		ProviderID:            string(summary.ProviderID),
		Label:                 summary.Label,
		Endpoint:              trimOptionalString(summary.Endpoint),
		CredentialState:       string(summary.CredentialState),
		CredentialReferenceID: trimOptionalString(summary.CredentialReferenceID),
		ValidationState:       string(summary.ValidationState),
		SavedState:            string(summary.SavedState),
		RequestToken:          trimOptionalString(summary.RequestToken),
		LastFailureKind:       toOptionalString(summary.LastFailureKind),
	}
}

func toValidateProviderSettingsResponseDTO(result usecase.ValidateProviderSettingsResult) ValidateProviderSettingsResponseDTO {
	return ValidateProviderSettingsResponseDTO{
		ProviderID:      string(result.ProviderID),
		ValidationState: string(result.ValidationState),
		RequestToken:    result.RequestToken,
		FailureKind:     toOptionalString(result.FailureKind),
	}
}

func trimOptionalString(value *string) *string {
	if value == nil {
		return nil
	}
	trimmed := strings.TrimSpace(*value)
	return &trimmed
}

func toOptionalString[T ~string](value *T) *string {
	if value == nil {
		return nil
	}
	stringValue := strings.TrimSpace(string(*value))
	return &stringValue
}
