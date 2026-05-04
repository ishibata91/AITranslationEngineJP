package usecase

import (
	"context"
	"strings"
)

// ProviderSettingsProviderID identifies one real provider exposed to the user.
type ProviderSettingsProviderID string

// ProviderSettingsProviderID values are the only real providers exposed through provider settings.
const (
	ProviderSettingsProviderIDGemini   ProviderSettingsProviderID = "gemini"
	ProviderSettingsProviderIDLMStudio ProviderSettingsProviderID = "lm_studio"
	ProviderSettingsProviderIDXAI      ProviderSettingsProviderID = "xai"
)

// ProviderSettingsCredentialState identifies the public credential state.
type ProviderSettingsCredentialState string

// ProviderSettingsCredentialState values describe the public credential state.
const (
	ProviderSettingsCredentialStateMissing     ProviderSettingsCredentialState = "missing"
	ProviderSettingsCredentialStateConfigured  ProviderSettingsCredentialState = "configured"
	ProviderSettingsCredentialStateNotRequired ProviderSettingsCredentialState = "not_required"
)

// ProviderSettingsValidationState identifies the public validation lifecycle.
type ProviderSettingsValidationState string

// ProviderSettingsValidationState values describe the public validation lifecycle.
const (
	ProviderSettingsValidationStateNotValidated ProviderSettingsValidationState = "not_validated"
	ProviderSettingsValidationStatePending      ProviderSettingsValidationState = "pending"
	ProviderSettingsValidationStateValidated    ProviderSettingsValidationState = "validated"
	ProviderSettingsValidationStateFailed       ProviderSettingsValidationState = "failed"
)

// ProviderSettingsSavedState identifies whether the settings row is stored.
type ProviderSettingsSavedState string

// ProviderSettingsSavedState values describe whether the settings row is stored.
const (
	ProviderSettingsSavedStateNotSaved   ProviderSettingsSavedState = "not_saved"
	ProviderSettingsSavedStatePartial    ProviderSettingsSavedState = "partial"
	ProviderSettingsSavedStateConfigured ProviderSettingsSavedState = "configured"
)

// ProviderSettingsModelListState identifies one provider model list loading state.
type ProviderSettingsModelListState string

// ProviderSettingsModelListState values describe one provider model list loading state.
const (
	ProviderSettingsModelListStateNotRequested ProviderSettingsModelListState = "not_requested"
	ProviderSettingsModelListStateLoading      ProviderSettingsModelListState = "loading"
	ProviderSettingsModelListStateReady        ProviderSettingsModelListState = "ready"
	ProviderSettingsModelListStateFailed       ProviderSettingsModelListState = "failed"
	//nolint:gosec // credential_missing is a fixed public contract literal, not a secret.
	ProviderSettingsModelListStateCredentialMissing ProviderSettingsModelListState = "credential_missing"
	ProviderSettingsModelListStateEndpointMissing   ProviderSettingsModelListState = "endpoint_missing"
	//nolint:gosec // credential_not_required is a fixed public contract literal, not a secret.
	ProviderSettingsModelListStateCredentialNotNeeded ProviderSettingsModelListState = "credential_not_required"
)

// ProviderSettingsErrorKind identifies redacted public failures and warnings.
type ProviderSettingsErrorKind string

// ProviderSettingsErrorKind values identify redacted public failures and warnings.
const (
	//nolint:gosec // credential_missing is a fixed public contract literal, not a secret.
	ProviderSettingsErrorKindCredentialMissing   ProviderSettingsErrorKind = "credential_missing"
	ProviderSettingsErrorKindEndpointMissing     ProviderSettingsErrorKind = "endpoint_missing"
	ProviderSettingsErrorKindValidationFailed    ProviderSettingsErrorKind = "validation_failed"
	ProviderSettingsErrorKindValidationStale     ProviderSettingsErrorKind = "validation_stale"
	ProviderSettingsErrorKindProviderUnreachable ProviderSettingsErrorKind = "provider_unreachable"
	ProviderSettingsErrorKindModelListFailed     ProviderSettingsErrorKind = "model_list_failed"
	ProviderSettingsErrorKindSettingsNotSaved    ProviderSettingsErrorKind = "settings_not_saved"
)

// NormalizeProviderSettingsErrorKind trims spacing while preserving the frozen public literal.
func NormalizeProviderSettingsErrorKind(kind ProviderSettingsErrorKind) ProviderSettingsErrorKind {
	return ProviderSettingsErrorKind(strings.TrimSpace(string(kind)))
}

// ProviderSettingsRoute describes the AppShell route contract fixed for downstream use.
type ProviderSettingsRoute struct {
	RouteID           string
	Label             string
	CurrentRouteState string
	DashboardEntryID  string
}

// ProviderSettingsSummary returns one provider row visible to the user.
type ProviderSettingsSummary struct {
	ProviderID            ProviderSettingsProviderID
	Label                 string
	Endpoint              *string
	CredentialState       ProviderSettingsCredentialState
	CredentialReferenceID *string
	ValidationState       ProviderSettingsValidationState
	SavedState            ProviderSettingsSavedState
	RequestToken          *string
	LastFailureKind       *ProviderSettingsErrorKind
}

// ProviderSettingsReferenceExecutionSelection keeps model-side settings outside provider settings save DTO.
type ProviderSettingsReferenceExecutionSelection struct {
	ProviderID      ProviderSettingsProviderID
	Model           string
	ExecutionMethod string
	UseBatchAPI     bool
}

// ProviderModelOption is one selectable provider model.
type ProviderModelOption struct {
	ModelID string
	Label   string
}

// ListProviderSettingsRequest identifies a provider settings list request.
type ListProviderSettingsRequest struct{}

// ListProviderSettingsResult returns the frozen provider settings read model.
type ListProviderSettingsResult struct {
	Route     ProviderSettingsRoute
	Providers []ProviderSettingsSummary
}

// SaveProviderSettingsRequest carries only reference values and input presence.
type SaveProviderSettingsRequest struct {
	ProviderID         ProviderSettingsProviderID
	Endpoint           *string
	APIKeyInputPresent bool
}

// SaveProviderSettingsResult returns one saved provider summary.
type SaveProviderSettingsResult struct {
	Provider ProviderSettingsSummary
}

// ResetProviderSettingsRequest identifies which provider settings row to reset.
type ResetProviderSettingsRequest struct {
	ProviderID ProviderSettingsProviderID
}

// ResetProviderSettingsResult returns one reset provider summary.
type ResetProviderSettingsResult struct {
	Provider ProviderSettingsSummary
}

// ValidateProviderSettingsRequest carries the current request correlation data.
type ValidateProviderSettingsRequest struct {
	ProviderID            ProviderSettingsProviderID
	Endpoint              *string
	CredentialState       ProviderSettingsCredentialState
	CredentialReferenceID *string
	RequestToken          string
}

// ValidateProviderSettingsResult returns the redacted validation outcome.
type ValidateProviderSettingsResult struct {
	ProviderID      ProviderSettingsProviderID
	ValidationState ProviderSettingsValidationState
	RequestToken    string
	FailureKind     *ProviderSettingsErrorKind
}

// ListProviderModelsRequest carries the provider settings state required for model lookup.
type ListProviderModelsRequest struct {
	ProviderID            ProviderSettingsProviderID
	Endpoint              *string
	CredentialState       ProviderSettingsCredentialState
	CredentialReferenceID *string
	RequestToken          string
}

// ListProviderModelsResult returns the model list state for one provider lookup.
type ListProviderModelsResult struct {
	ProviderID      ProviderSettingsProviderID
	Endpoint        *string
	CredentialState ProviderSettingsCredentialState
	RequestToken    string
	State           ProviderSettingsModelListState
	Models          []ProviderModelOption
	FailureKind     *ProviderSettingsErrorKind
}

// ResolveProviderExecutionSettingsRequest identifies one reference-side resolution request.
type ResolveProviderExecutionSettingsRequest struct {
	ConsumerID string
	Selection  ProviderSettingsReferenceExecutionSelection
}

// ResolveProviderExecutionSettingsResult returns the resolved reference values without exposing secrets.
type ResolveProviderExecutionSettingsResult struct {
	ConsumerID            string
	ProviderID            ProviderSettingsProviderID
	Model                 string
	ExecutionMethod       string
	UseBatchAPI           bool
	Endpoint              *string
	CredentialReferenceID *string
	CredentialState       ProviderSettingsCredentialState
	RequestToken          *string
	ErrorKind             *ProviderSettingsErrorKind
}

// ProviderSettingsUsecase defines the frozen provider settings seam for Wails and downstream consumers.
type ProviderSettingsUsecase interface {
	ListProviderSettings(ctx context.Context, request ListProviderSettingsRequest) (ListProviderSettingsResult, error)
	SaveProviderSettings(ctx context.Context, request SaveProviderSettingsRequest) (SaveProviderSettingsResult, error)
	ResetProviderSettings(ctx context.Context, request ResetProviderSettingsRequest) (ResetProviderSettingsResult, error)
	ValidateProviderSettings(ctx context.Context, request ValidateProviderSettingsRequest) (ValidateProviderSettingsResult, error)
	ListProviderModels(ctx context.Context, request ListProviderModelsRequest) (ListProviderModelsResult, error)
	ResolveProviderExecutionSettings(ctx context.Context, request ResolveProviderExecutionSettingsRequest) (ResolveProviderExecutionSettingsResult, error)
}
