package ai

import (
	"context"
	"fmt"
	"strings"
)

// ValidateProviderSettingsRequest carries one provider settings validation probe request.
type ValidateProviderSettingsRequest struct {
	ProviderID string
	Endpoint   string
	APIKey     string
}

// ProviderSettingsValidator keeps provider settings validation transport inside the infra adapter.
type ProviderSettingsValidator struct {
	modelListLoader *ProviderModelListLoader
}

// NewProviderSettingsValidator creates the infra adapter used by provider settings validation.
func NewProviderSettingsValidator(modelListLoader *ProviderModelListLoader) *ProviderSettingsValidator {
	return &ProviderSettingsValidator{modelListLoader: modelListLoader}
}

// ValidateProviderSettings verifies that the configured provider endpoint and credential can complete one lightweight probe.
func (validator *ProviderSettingsValidator) ValidateProviderSettings(
	ctx context.Context,
	request ValidateProviderSettingsRequest,
) error {
	if validator == nil || validator.modelListLoader == nil {
		return fmt.Errorf("provider settings validator is required")
	}
	_, err := validator.modelListLoader.ListProviderModelsWithEndpoint(
		ctx,
		strings.TrimSpace(request.ProviderID),
		strings.TrimSpace(request.Endpoint),
		strings.TrimSpace(request.APIKey),
	)
	if err != nil {
		return fmt.Errorf("validate provider settings: %w", err)
	}
	return nil
}
