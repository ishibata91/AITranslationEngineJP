package service

import (
	"context"
	"crypto/sha256"
	"errors"
	"fmt"
	"log/slog"
	"sort"
	"strings"
	"sync"
	"time"

	"aitranslationenginejp/internal/repository"
)

var (
	// ErrProviderSettingsValidation reports invalid provider settings input.
	ErrProviderSettingsValidation = errors.New("provider settings validation error")
)

const (
	providerSettingsRouteID          = "provider-settings"
	providerSettingsRouteLabel       = "AIサービス設定"
	providerSettingsRouteState       = "current"
	providerSettingsDashboardEntryID = "ai-provider-settings"

	providerSettingsCredentialStateMissing     = "missing"
	providerSettingsCredentialStateConfigured  = "configured"
	providerSettingsCredentialStateNotRequired = "not_required"

	providerSettingsValidationStateNotValidated = "not_validated"
	providerSettingsValidationStatePending      = "pending"
	providerSettingsValidationStateValidated    = "validated"
	providerSettingsValidationStateFailed       = "failed"

	providerSettingsSavedStateNotSaved   = "not_saved"
	providerSettingsSavedStatePartial    = "partial"
	providerSettingsSavedStateConfigured = "configured"

	providerSettingsModelListStateNotRequested = "not_requested"
	providerSettingsModelListStateReady        = "ready"
	providerSettingsModelListStateFailed       = "failed"
	//nolint:gosec // credential_missing is a fixed public contract literal, not a secret.
	providerSettingsModelListStateCredentialMissing = "credential_missing"
	providerSettingsModelListStateEndpointMissing   = "endpoint_missing"
	//nolint:gosec // credential_not_required is a fixed public contract literal, not a secret.
	providerSettingsModelListStateCredentialNotNeeded = "credential_not_required"

	//nolint:gosec // credential_missing is a fixed public contract literal, not a secret.
	providerSettingsErrorKindCredentialMissing   = "credential_missing"
	providerSettingsErrorKindEndpointMissing     = "endpoint_missing"
	providerSettingsErrorKindValidationFailed    = "validation_failed"
	providerSettingsErrorKindValidationStale     = "validation_stale"
	providerSettingsErrorKindProviderUnreachable = "provider_unreachable"
	providerSettingsErrorKindModelListFailed     = "model_list_failed"
	providerSettingsErrorKindSettingsNotSaved    = "settings_not_saved"
	providerSettingsServiceWhere                 = "provider_settings.service"
)

type providerSettingsProviderSpec struct {
	ID                 string
	Label              string
	CredentialRequired bool
	DefaultEndpoint    string
}

var providerSettingsProviderCatalog = []providerSettingsProviderSpec{
	{ID: "gemini", Label: "Gemini", CredentialRequired: true, DefaultEndpoint: "https://generativelanguage.googleapis.com/v1beta"},
	{ID: "lm_studio", Label: "LM Studio", CredentialRequired: false, DefaultEndpoint: "http://localhost:1234/v1"},
	{ID: "xai", Label: "xAI", CredentialRequired: true, DefaultEndpoint: "https://api.x.ai/v1"},
}

var providerSettingsMutationLocks sync.Map

// ProviderSettingsRoute stores the fixed AppShell route contract.
type ProviderSettingsRoute struct {
	RouteID           string
	Label             string
	CurrentRouteState string
	DashboardEntryID  string
}

// ProviderSettingsSummary stores one provider row visible to callers.
type ProviderSettingsSummary struct {
	ProviderID            string
	Label                 string
	Endpoint              *string
	CredentialState       string
	CredentialReferenceID *string
	ValidationState       string
	SavedState            string
	RequestToken          *string
	LastFailureKind       *string
}

// ProviderSettingsSaveInput stores one secret-aware save request.
type ProviderSettingsSaveInput struct {
	ProviderID         string
	Endpoint           *string
	APIKeyInputPresent bool
	APIKey             *string
}

// ProviderSettingsValidateInput stores one validation request.
type ProviderSettingsValidateInput struct {
	ProviderID            string
	Endpoint              *string
	CredentialState       string
	CredentialReferenceID *string
	RequestToken          string
}

// ProviderSettingsValidationResult stores one validation outcome.
type ProviderSettingsValidationResult struct {
	ProviderID      string
	ValidationState string
	RequestToken    string
	FailureKind     *string
}

// ProviderSettingsValidationProbe stores one provider validation transport request.
type ProviderSettingsValidationProbe struct {
	ProviderID string
	Endpoint   string
	APIKey     string
}

// ProviderSettingsModelListInput stores one provider model list request.
type ProviderSettingsModelListInput struct {
	ProviderID            string
	Endpoint              *string
	CredentialState       string
	CredentialReferenceID *string
	RequestToken          string
}

// ProviderSettingsModelOption stores one provider model choice.
type ProviderSettingsModelOption struct {
	ModelID string
	Label   string
}

// ProviderSettingsModelListResult stores one provider model list response.
type ProviderSettingsModelListResult struct {
	ProviderID      string
	Endpoint        *string
	CredentialState string
	RequestToken    string
	State           string
	Models          []ProviderSettingsModelOption
	FailureKind     *string
}

// ProviderSettingsResolveSelection stores one reference-side execution selection.
type ProviderSettingsResolveSelection struct {
	ProviderID      string
	Model           string
	ExecutionMethod string
	UseBatchAPI     bool
}

// ProviderSettingsResolveInput stores one resolve request.
type ProviderSettingsResolveInput struct {
	ConsumerID          string
	Selection           ProviderSettingsResolveSelection
	AllowSecretSnapshot bool
}

// ProviderSettingsResolveResult stores one resolved reference-side execution setting.
type ProviderSettingsResolveResult struct {
	ConsumerID            string
	ProviderID            string
	Model                 string
	ExecutionMethod       string
	UseBatchAPI           bool
	Endpoint              *string
	CredentialReferenceID *string
	CredentialState       string
	RequestToken          *string
	ErrorKind             *string
}

// ProviderSettingsRepository defines persistence dependencies.
type ProviderSettingsRepository interface {
	List(ctx context.Context) ([]repository.ProviderSettingsRecord, error)
	GetByProviderID(ctx context.Context, providerID string) (repository.ProviderSettingsRecord, error)
	Upsert(ctx context.Context, draft repository.ProviderSettingsUpsertDraft) (repository.ProviderSettingsRecord, error)
	UpdateValidationByRequestToken(
		ctx context.Context,
		providerID string,
		requestToken string,
		validationState string,
		lastFailureKind *string,
	) (repository.ProviderSettingsRecord, bool, error)
}

// ProviderSettingsService provides provider settings backend operations.
type ProviderSettingsService struct {
	repository      ProviderSettingsRepository
	secretStore     repository.ProviderSettingsSecretStore
	transactor      repository.Transactor
	modelListLoader providerSettingsModelListLoader
	validator       providerSettingsValidator
	now             func() time.Time
}

type providerSettingsModelListLoader interface {
	ListProviderModelsWithEndpoint(
		ctx context.Context,
		providerID string,
		endpoint string,
		apiKey string,
	) ([]ProviderSettingsModelOption, error)
}

type providerSettingsValidator interface {
	ValidateProviderSettings(ctx context.Context, request ProviderSettingsValidationProbe) error
}

// NewProviderSettingsService creates a provider settings service.
func NewProviderSettingsService(
	repository ProviderSettingsRepository,
	secretStore repository.ProviderSettingsSecretStore,
	transactor repository.Transactor,
	modelListLoader providerSettingsModelListLoader,
	validator providerSettingsValidator,
	now func() time.Time,
) *ProviderSettingsService {
	if now == nil {
		now = time.Now
	}
	return &ProviderSettingsService{
		repository:      repository,
		secretStore:     secretStore,
		transactor:      transactor,
		modelListLoader: modelListLoader,
		validator:       validator,
		now:             now,
	}
}

// ListProviderSettings returns the fixed provider settings read model.
func (service *ProviderSettingsService) ListProviderSettings(ctx context.Context) (ProviderSettingsRoute, []ProviderSettingsSummary, error) {
	rows, err := service.repository.List(providerSettingsContextOrBackground(ctx))
	if err != nil {
		return ProviderSettingsRoute{}, nil, fmt.Errorf("list provider settings rows: %w", err)
	}
	rowMap := make(map[string]repository.ProviderSettingsRecord, len(rows))
	for _, row := range rows {
		rowMap[strings.TrimSpace(row.ProviderID)] = row
	}
	summaries := make([]ProviderSettingsSummary, 0, len(providerSettingsProviderCatalog))
	for _, spec := range providerSettingsProviderCatalog {
		row, ok := rowMap[spec.ID]
		summary, buildErr := service.buildSummary(providerSettingsContextOrBackground(ctx), spec, row, ok)
		if buildErr != nil {
			return ProviderSettingsRoute{}, nil, buildErr
		}
		summaries = append(summaries, summary)
	}
	return ProviderSettingsRoute{
		RouteID:           providerSettingsRouteID,
		Label:             providerSettingsRouteLabel,
		CurrentRouteState: providerSettingsRouteState,
		DashboardEntryID:  providerSettingsDashboardEntryID,
	}, summaries, nil
}

// SaveProviderSettings stores one provider settings row and secret as one compensating unit.
func (service *ProviderSettingsService) SaveProviderSettings(
	ctx context.Context,
	input ProviderSettingsSaveInput,
) (ProviderSettingsSummary, error) {
	spec, err := providerSettingsSpec(strings.TrimSpace(input.ProviderID))
	if err != nil {
		return ProviderSettingsSummary{}, err
	}
	var saved ProviderSettingsSummary
	lockErr := withProviderSettingsMutationLock(spec.ID, func() error {
		summary, err := service.saveProviderSettingsLocked(providerSettingsContextOrBackground(ctx), spec, input)
		saved = summary
		return err
	})
	if lockErr != nil {
		return ProviderSettingsSummary{}, lockErr
	}
	return saved, nil
}

func (service *ProviderSettingsService) saveProviderSettingsLocked(
	ctx context.Context,
	spec providerSettingsProviderSpec,
	input ProviderSettingsSaveInput,
) (ProviderSettingsSummary, error) {
	current, currentExists, err := service.loadRecord(ctx, spec.ID)
	if err != nil {
		return ProviderSettingsSummary{}, err
	}
	revision := current.Revision + 1
	requestToken := providerSettingsRequestToken(spec.ID, revision)
	credentialRef, credentialState, secretValue, secretChanged, err := service.resolveSecretSavePlan(
		spec,
		current,
		currentExists,
		input,
		requestToken,
	)
	if err != nil {
		return ProviderSettingsSummary{}, err
	}
	secretKey := providerSettingsSecretKey(spec.ID)
	savedRecord, err := service.upsertProviderSettingsRecord(ctx, providerSettingsSaveDraft(
		spec,
		input,
		credentialRef,
		credentialState,
		requestToken,
		revision,
		service.now().UTC(),
	), "save provider settings")
	if err != nil {
		return ProviderSettingsSummary{}, err
	}
	if err := service.saveCanonicalProviderSettingsSecret(ctx, secretChanged, secretKey, secretValue, current, currentExists); err != nil {
		return ProviderSettingsSummary{}, err
	}
	return service.buildSummary(ctx, spec, savedRecord, true)
}

func providerSettingsSaveDraft(
	spec providerSettingsProviderSpec,
	input ProviderSettingsSaveInput,
	credentialRef *string,
	credentialState string,
	requestToken string,
	revision int64,
	now time.Time,
) repository.ProviderSettingsUpsertDraft {
	return repository.ProviderSettingsUpsertDraft{
		ProviderID:            spec.ID,
		Endpoint:              normalizeServiceOptionalString(input.Endpoint),
		CredentialReferenceID: credentialRef,
		CredentialState:       credentialState,
		ValidationState:       providerSettingsValidationStateNotValidated,
		RequestToken:          providerSettingsStringPointer(requestToken),
		LastFailureKind:       nil,
		Revision:              revision,
		UpdatedAt:             now,
	}
}

func (service *ProviderSettingsService) upsertProviderSettingsRecord(
	ctx context.Context,
	draft repository.ProviderSettingsUpsertDraft,
	operation string,
) (repository.ProviderSettingsRecord, error) {
	var savedRecord repository.ProviderSettingsRecord
	txErr := service.transactor.WithTransaction(ctx, func(txCtx context.Context) error {
		record, upsertErr := service.repository.Upsert(txCtx, draft)
		if upsertErr != nil {
			return fmt.Errorf("%s row: %w", operation, upsertErr)
		}
		savedRecord = record
		return nil
	})
	if txErr != nil {
		return repository.ProviderSettingsRecord{}, fmt.Errorf("%s: %w", operation, txErr)
	}
	return savedRecord, nil
}

func (service *ProviderSettingsService) saveCanonicalProviderSettingsSecret(
	ctx context.Context,
	secretChanged bool,
	secretKey string,
	secretValue string,
	current repository.ProviderSettingsRecord,
	currentExists bool,
) error {
	if !secretChanged {
		return nil
	}
	if err := service.secretStore.Save(ctx, secretKey, secretValue); err != nil {
		if rollbackErr := service.restoreProviderSettingsRecord(ctx, current, currentExists); rollbackErr != nil {
			return fmt.Errorf("save provider settings canonical secret: %w; rollback provider settings row: %w", err, rollbackErr)
		}
		return fmt.Errorf("save provider settings canonical secret: %w", err)
	}
	return nil
}

// ResetProviderSettings clears endpoint and credential state while keeping the row.
func (service *ProviderSettingsService) ResetProviderSettings(
	ctx context.Context,
	providerID string,
) (ProviderSettingsSummary, error) {
	spec, err := providerSettingsSpec(providerID)
	if err != nil {
		return ProviderSettingsSummary{}, err
	}
	var reset ProviderSettingsSummary
	lockErr := withProviderSettingsMutationLock(spec.ID, func() error {
		summary, err := service.resetProviderSettingsLocked(providerSettingsContextOrBackground(ctx), spec)
		reset = summary
		return err
	})
	if lockErr != nil {
		return ProviderSettingsSummary{}, lockErr
	}
	reset.Endpoint = nil
	return reset, nil
}

func (service *ProviderSettingsService) resetProviderSettingsLocked(
	ctx context.Context,
	spec providerSettingsProviderSpec,
) (ProviderSettingsSummary, error) {
	current, currentExists, err := service.loadRecord(ctx, spec.ID)
	if err != nil {
		return ProviderSettingsSummary{}, err
	}
	revision := current.Revision + 1
	savedRecord, err := service.upsertProviderSettingsRecord(ctx, providerSettingsResetDraft(
		spec,
		providerSettingsRequestToken(spec.ID, revision),
		revision,
		service.now().UTC(),
	), "reset provider settings")
	if err != nil {
		return ProviderSettingsSummary{}, err
	}
	if err := service.deleteProviderSettingsSecret(ctx, providerSettingsSecretKey(spec.ID), current, currentExists); err != nil {
		return ProviderSettingsSummary{}, err
	}
	return service.buildSummary(ctx, spec, savedRecord, true)
}

func providerSettingsResetDraft(
	spec providerSettingsProviderSpec,
	requestToken string,
	revision int64,
	now time.Time,
) repository.ProviderSettingsUpsertDraft {
	credentialState := providerSettingsCredentialStateMissing
	if !spec.CredentialRequired {
		credentialState = providerSettingsCredentialStateNotRequired
	}
	return repository.ProviderSettingsUpsertDraft{
		ProviderID:            spec.ID,
		Endpoint:              nil,
		CredentialReferenceID: nil,
		CredentialState:       credentialState,
		ValidationState:       providerSettingsValidationStateNotValidated,
		RequestToken:          providerSettingsStringPointer(requestToken),
		LastFailureKind:       nil,
		Revision:              revision,
		UpdatedAt:             now,
	}
}

func (service *ProviderSettingsService) deleteProviderSettingsSecret(
	ctx context.Context,
	secretKey string,
	current repository.ProviderSettingsRecord,
	currentExists bool,
) error {
	if err := service.secretStore.Delete(ctx, secretKey); err != nil {
		if rollbackErr := service.restoreProviderSettingsRecord(ctx, current, currentExists); rollbackErr != nil {
			return fmt.Errorf("delete provider settings secret: %w; rollback provider settings row: %w", err, rollbackErr)
		}
		return fmt.Errorf("delete provider settings secret: %w", err)
	}
	return nil
}

// ValidateProviderSettings validates one provider settings snapshot without leaking secrets.
func (service *ProviderSettingsService) ValidateProviderSettings(
	ctx context.Context,
	input ProviderSettingsValidateInput,
) (ProviderSettingsValidationResult, error) {
	spec, err := providerSettingsSpec(input.ProviderID)
	if err != nil {
		return ProviderSettingsValidationResult{}, err
	}
	var result ProviderSettingsValidationResult
	lockErr := withProviderSettingsMutationLock(spec.ID, func() error {
		validated, err := service.validateProviderSettingsLocked(providerSettingsContextOrBackground(ctx), spec, input)
		result = validated
		return err
	})
	if lockErr != nil {
		return ProviderSettingsValidationResult{}, lockErr
	}
	return result, nil
}

func (service *ProviderSettingsService) validateProviderSettingsLocked(
	ctx context.Context,
	spec providerSettingsProviderSpec,
	input ProviderSettingsValidateInput,
) (ProviderSettingsValidationResult, error) {
	summary, err := service.getCurrentSummary(ctx, spec.ID)
	if err != nil {
		return ProviderSettingsValidationResult{}, err
	}
	requestToken := strings.TrimSpace(input.RequestToken)
	if requestToken == "" {
		return ProviderSettingsValidationResult{}, fmt.Errorf("%w: request token is required", ErrProviderSettingsValidation)
	}
	if failureKind := providerSettingsValidationPreflightFailure(spec, summary, input, requestToken); failureKind != "" {
		logProviderBoundaryFailure(ctx, "provider_settings_validation", providerSettingsServiceWhere, spec.ID, failureKind)
		return providerSettingsValidationFailure(spec.ID, requestToken, failureKind), nil
	}
	apiKey, failureKind, err := service.resolveValidationSecret(ctx, spec, summary)
	if err != nil {
		logProviderBoundaryFailure(ctx, "provider_settings_validation", providerSettingsServiceWhere, spec.ID, "secret_store_failed")
		return ProviderSettingsValidationResult{}, err
	}
	if failureKind != nil {
		logProviderBoundaryFailure(ctx, "provider_settings_validation", providerSettingsServiceWhere, spec.ID, *failureKind)
		return providerSettingsValidationFailurePointer(spec.ID, requestToken, failureKind), nil
	}
	pendingStale, pendingErr := service.updateProviderSettingsValidation(ctx, spec.ID, requestToken, providerSettingsValidationStatePending, nil, "mark provider settings validation pending")
	if pendingErr != nil || pendingStale {
		return providerSettingsValidationStaleOrError(spec.ID, requestToken, pendingStale, pendingErr)
	}
	finalState, finalFailure := service.validateProviderSettingsProbe(ctx, spec, summary, apiKey)
	stale, err := service.updateProviderSettingsValidation(ctx, spec.ID, requestToken, finalState, finalFailure, "update provider settings validation result")
	if err != nil || stale {
		return providerSettingsValidationStaleOrError(spec.ID, requestToken, stale, err)
	}
	return ProviderSettingsValidationResult{
		ProviderID:      spec.ID,
		ValidationState: finalState,
		RequestToken:    requestToken,
		FailureKind:     finalFailure,
	}, nil
}

func providerSettingsValidationPreflightFailure(
	spec providerSettingsProviderSpec,
	summary ProviderSettingsSummary,
	input ProviderSettingsValidateInput,
	requestToken string,
) string {
	if !providerSettingsMatchesSnapshot(summary, input.Endpoint, input.CredentialState, input.CredentialReferenceID, requestToken) {
		return providerSettingsErrorKindValidationStale
	}
	if summary.Endpoint == nil {
		return providerSettingsErrorKindEndpointMissing
	}
	if spec.CredentialRequired && summary.CredentialReferenceID == nil {
		return providerSettingsErrorKindCredentialMissing
	}
	return ""
}

func (service *ProviderSettingsService) updateProviderSettingsValidation(
	ctx context.Context,
	providerID string,
	requestToken string,
	state string,
	failureKind *string,
	operation string,
) (bool, error) {
	_, matched, err := service.repository.UpdateValidationByRequestToken(ctx, providerID, requestToken, state, failureKind)
	if err != nil {
		return false, fmt.Errorf("%s: %w", operation, err)
	}
	return !matched, nil
}

func (service *ProviderSettingsService) validateProviderSettingsProbe(
	ctx context.Context,
	spec providerSettingsProviderSpec,
	summary ProviderSettingsSummary,
	apiKey string,
) (string, *string) {
	err := service.validator.ValidateProviderSettings(ctx, ProviderSettingsValidationProbe{
		ProviderID: spec.ID,
		Endpoint:   *summary.Endpoint,
		APIKey:     apiKey,
	})
	if err != nil {
		logProviderBoundaryFailure(ctx, "provider_settings_validation", providerSettingsServiceWhere, spec.ID, classifyProviderBoundaryError(err, providerSettingsErrorKindProviderUnreachable))
		return providerSettingsValidationStateFailed, providerSettingsStringPointer(providerSettingsErrorKindProviderUnreachable)
	}
	return providerSettingsValidationStateValidated, nil
}

func providerSettingsValidationStaleOrError(
	providerID string,
	requestToken string,
	stale bool,
	err error,
) (ProviderSettingsValidationResult, error) {
	if err != nil {
		return ProviderSettingsValidationResult{}, err
	}
	if stale {
		return providerSettingsValidationFailure(providerID, requestToken, providerSettingsErrorKindValidationStale), nil
	}
	return ProviderSettingsValidationResult{}, nil
}

func providerSettingsValidationFailure(providerID string, requestToken string, failureKind string) ProviderSettingsValidationResult {
	return providerSettingsValidationFailurePointer(providerID, requestToken, providerSettingsStringPointer(failureKind))
}

func providerSettingsValidationFailurePointer(providerID string, requestToken string, failureKind *string) ProviderSettingsValidationResult {
	return ProviderSettingsValidationResult{
		ProviderID:      providerID,
		ValidationState: providerSettingsValidationStateNotValidated,
		RequestToken:    requestToken,
		FailureKind:     failureKind,
	}
}

// ListProviderModels returns provider model options only when endpoint and credential state are ready.
func (service *ProviderSettingsService) ListProviderModels(
	ctx context.Context,
	input ProviderSettingsModelListInput,
) (ProviderSettingsModelListResult, error) {
	spec, err := providerSettingsSpec(input.ProviderID)
	if err != nil {
		return ProviderSettingsModelListResult{}, err
	}
	var result ProviderSettingsModelListResult
	lockErr := withProviderSettingsMutationLock(spec.ID, func() error {
		listed, err := service.listProviderModelsLocked(providerSettingsContextOrBackground(ctx), spec, input)
		result = listed
		return err
	})
	if lockErr != nil {
		return ProviderSettingsModelListResult{}, lockErr
	}
	return result, nil
}

func (service *ProviderSettingsService) listProviderModelsLocked(
	ctx context.Context,
	spec providerSettingsProviderSpec,
	input ProviderSettingsModelListInput,
) (ProviderSettingsModelListResult, error) {
	summary, err := service.getCurrentSummary(ctx, spec.ID)
	if err != nil {
		return ProviderSettingsModelListResult{}, err
	}
	requestToken := strings.TrimSpace(input.RequestToken)
	result := providerSettingsInitialModelListResult(spec.ID, summary, requestToken)
	if done := providerSettingsApplyModelListPreflight(&result, summary, input, requestToken); done {
		logProviderBoundaryFailure(ctx, "provider_model_list", providerSettingsServiceWhere, spec.ID, valueOrEmpty(result.FailureKind))
		return result, nil
	}
	apiKey, done, err := service.resolveModelListAPIKey(ctx, spec, summary, &result)
	if err != nil || done {
		if err != nil {
			logProviderBoundaryFailure(ctx, "provider_model_list", providerSettingsServiceWhere, spec.ID, "secret_store_failed")
		}
		if done {
			logProviderBoundaryFailure(ctx, "provider_model_list", providerSettingsServiceWhere, spec.ID, valueOrEmpty(result.FailureKind))
		}
		return result, err
	}
	models, loadErr := service.modelListLoader.ListProviderModelsWithEndpoint(ctx, spec.ID, *summary.Endpoint, apiKey)
	if loadErr != nil {
		result.State = providerSettingsModelListStateFailed
		result.FailureKind = providerSettingsStringPointer(providerSettingsErrorKindModelListFailed)
		logProviderBoundaryFailure(ctx, "provider_model_list", providerSettingsServiceWhere, spec.ID, classifyProviderBoundaryError(loadErr, providerSettingsErrorKindModelListFailed))
		_ = loadErr
		//nolint:nilerr // public contract intentionally returns a redacted failure instead of the transport error.
		return result, nil
	}
	result.State = providerSettingsModelListReadyState(spec)
	result.Models = providerSettingsModelOptions(models)
	return result, nil
}

func providerSettingsInitialModelListResult(
	providerID string,
	summary ProviderSettingsSummary,
	requestToken string,
) ProviderSettingsModelListResult {
	return ProviderSettingsModelListResult{
		ProviderID:      providerID,
		Endpoint:        summary.Endpoint,
		CredentialState: summary.CredentialState,
		RequestToken:    requestToken,
		State:           providerSettingsModelListStateNotRequested,
		Models:          []ProviderSettingsModelOption{},
	}
}

func providerSettingsApplyModelListPreflight(
	result *ProviderSettingsModelListResult,
	summary ProviderSettingsSummary,
	input ProviderSettingsModelListInput,
	requestToken string,
) bool {
	if !providerSettingsMatchesSnapshot(summary, input.Endpoint, input.CredentialState, input.CredentialReferenceID, requestToken) {
		result.State = providerSettingsModelListStateFailed
		result.FailureKind = providerSettingsStringPointer(providerSettingsErrorKindValidationStale)
		return true
	}
	if summary.Endpoint == nil {
		result.State = providerSettingsModelListStateEndpointMissing
		result.FailureKind = providerSettingsStringPointer(providerSettingsErrorKindEndpointMissing)
		return true
	}
	return false
}

func (service *ProviderSettingsService) resolveModelListAPIKey(
	ctx context.Context,
	spec providerSettingsProviderSpec,
	summary ProviderSettingsSummary,
	result *ProviderSettingsModelListResult,
) (string, bool, error) {
	if !spec.CredentialRequired {
		result.State = providerSettingsModelListStateCredentialNotNeeded
		return "", false, nil
	}
	if summary.CredentialReferenceID == nil {
		result.State = providerSettingsModelListStateCredentialMissing
		result.FailureKind = providerSettingsStringPointer(providerSettingsErrorKindCredentialMissing)
		return "", true, nil
	}
	loaded, failureKind, err := service.resolveValidationSecret(ctx, spec, summary)
	if err != nil {
		return "", false, err
	}
	if failureKind != nil {
		result.State = providerSettingsModelListStateCredentialMissing
		result.FailureKind = failureKind
		return "", true, nil
	}
	return loaded, false, nil
}

func providerSettingsModelListReadyState(spec providerSettingsProviderSpec) string {
	if !spec.CredentialRequired {
		return providerSettingsModelListStateCredentialNotNeeded
	}
	return providerSettingsModelListStateReady
}

func providerSettingsModelOptions(models []ProviderSettingsModelOption) []ProviderSettingsModelOption {
	options := make([]ProviderSettingsModelOption, 0, len(models))
	for _, model := range models {
		options = append(options, ProviderSettingsModelOption{
			ModelID: strings.TrimSpace(model.ModelID),
			Label:   strings.TrimSpace(model.Label),
		})
	}
	sort.Slice(options, func(left, right int) bool {
		return options[left].ModelID < options[right].ModelID
	})
	return options
}

// ResolveProviderExecutionSettings resolves provider endpoint and credential reference for downstream consumers.
func (service *ProviderSettingsService) ResolveProviderExecutionSettings(
	ctx context.Context,
	input ProviderSettingsResolveInput,
) (ProviderSettingsResolveResult, error) {
	spec, err := providerSettingsSpec(input.Selection.ProviderID)
	if err != nil {
		return ProviderSettingsResolveResult{}, err
	}
	resolve := func() (ProviderSettingsResolveResult, error) {
		return service.resolveProviderExecutionSettingsUnlocked(providerSettingsContextOrBackground(ctx), spec, input)
	}
	if !input.AllowSecretSnapshot {
		return resolve()
	}
	var result ProviderSettingsResolveResult
	lockErr := withProviderSettingsMutationLock(spec.ID, func() error {
		resolved, err := resolve()
		if err != nil {
			return err
		}
		result = resolved
		return nil
	})
	if lockErr != nil {
		return ProviderSettingsResolveResult{}, lockErr
	}
	return result, nil
}

func (service *ProviderSettingsService) resolveProviderExecutionSettingsUnlocked(
	ctx context.Context,
	spec providerSettingsProviderSpec,
	input ProviderSettingsResolveInput,
) (ProviderSettingsResolveResult, error) {
	summary, err := service.getCurrentSummary(ctx, spec.ID)
	if err != nil {
		return ProviderSettingsResolveResult{}, err
	}
	result := providerSettingsResolveResult(spec, input, summary)
	if err := service.applyProviderSettingsExecutionSnapshot(ctx, spec, input, summary, &result); err != nil {
		logProviderBoundaryFailure(ctx, "provider_execution_settings", providerSettingsServiceWhere, spec.ID, "secret_store_failed")
		return ProviderSettingsResolveResult{}, err
	}
	if result.ErrorKind != nil {
		logProviderBoundaryFailure(ctx, "provider_execution_settings", providerSettingsServiceWhere, spec.ID, *result.ErrorKind)
	}
	return result, nil
}

func logProviderBoundaryFailure(ctx context.Context, event string, where string, providerID string, reason string) {
	trimmedReason := strings.TrimSpace(reason)
	if trimmedReason == "" {
		trimmedReason = "unknown"
	}
	slog.WarnContext(providerSettingsContextOrBackground(ctx), "provider boundary failed",
		slog.String("event", strings.TrimSpace(event)),
		slog.String("where", strings.TrimSpace(where)),
		slog.String("result", "failed"),
		slog.String("provider", strings.TrimSpace(providerID)),
		slog.String("reason", trimmedReason),
	)
}

func logProviderBoundarySkipped(ctx context.Context, event string, where string, providerID string, reason string, count int) {
	trimmedReason := strings.TrimSpace(reason)
	if trimmedReason == "" {
		trimmedReason = "provider_skipped"
	}
	slog.InfoContext(providerSettingsContextOrBackground(ctx), "provider boundary skipped",
		slog.String("event", strings.TrimSpace(event)),
		slog.String("where", strings.TrimSpace(where)),
		slog.String("result", "skipped"),
		slog.String("provider", strings.TrimSpace(providerID)),
		slog.String("reason", trimmedReason),
		slog.Int("count", count),
	)
}

func classifyProviderBoundaryError(err error, fallback string) string {
	if errors.Is(err, context.DeadlineExceeded) {
		return "provider_timeout"
	}
	trimmedFallback := strings.TrimSpace(fallback)
	if trimmedFallback == "" {
		return "provider_failure"
	}
	return trimmedFallback
}

func normalizePersonaGenerationProviderFailureKind(kind string) string {
	switch strings.TrimSpace(kind) {
	case personaGenerationErrorKindInvalidProvider:
		return "invalid_provider_response"
	case "provider_timeout":
		return "provider_timeout"
	default:
		return personaGenerationErrorKindProviderFailure
	}
}

func providerSettingsResolveResult(
	spec providerSettingsProviderSpec,
	input ProviderSettingsResolveInput,
	summary ProviderSettingsSummary,
) ProviderSettingsResolveResult {
	result := ProviderSettingsResolveResult{
		ConsumerID:            strings.TrimSpace(input.ConsumerID),
		ProviderID:            spec.ID,
		Model:                 strings.TrimSpace(input.Selection.Model),
		ExecutionMethod:       strings.TrimSpace(input.Selection.ExecutionMethod),
		UseBatchAPI:           input.Selection.UseBatchAPI,
		Endpoint:              summary.Endpoint,
		CredentialReferenceID: summary.CredentialReferenceID,
		CredentialState:       summary.CredentialState,
		RequestToken:          summary.RequestToken,
	}
	switch {
	case summary.Endpoint == nil:
		result.ErrorKind = providerSettingsStringPointer(providerSettingsErrorKindEndpointMissing)
	case spec.CredentialRequired && summary.CredentialReferenceID == nil:
		result.ErrorKind = providerSettingsStringPointer(providerSettingsErrorKindCredentialMissing)
	}
	return result
}

func (service *ProviderSettingsService) applyProviderSettingsExecutionSnapshot(
	ctx context.Context,
	spec providerSettingsProviderSpec,
	input ProviderSettingsResolveInput,
	summary ProviderSettingsSummary,
	result *ProviderSettingsResolveResult,
) error {
	if !input.AllowSecretSnapshot || result.ErrorKind != nil || !spec.CredentialRequired || summary.CredentialReferenceID == nil {
		return nil
	}
	snapshotRef, err := service.snapshotExecutionSecret(ctx, strings.TrimSpace(input.ConsumerID), spec.ID, summary)
	if err != nil {
		return err
	}
	result.CredentialReferenceID = providerSettingsStringPointer(snapshotRef)
	return nil
}

func (service *ProviderSettingsService) snapshotExecutionSecret(
	ctx context.Context,
	consumerID string,
	providerID string,
	summary ProviderSettingsSummary,
) (string, error) {
	requestToken := valueOrEmpty(summary.RequestToken)
	sourceCredentialRef := valueOrEmpty(summary.CredentialReferenceID)
	sourceRef := strings.TrimSpace(sourceCredentialRef)
	if sourceRef == "" {
		return "", fmt.Errorf("provider settings credential reference is required")
	}
	if strings.TrimSpace(requestToken) == "" {
		return "", fmt.Errorf("provider settings request token is required")
	}
	secret, err := service.secretStore.Load(ctx, sourceRef)
	if err != nil {
		return "", fmt.Errorf("load provider settings execution secret: %w", err)
	}
	trimmedSecret := strings.TrimSpace(secret)
	if trimmedSecret == "" {
		return "", fmt.Errorf("provider settings execution secret is missing")
	}
	snapshotRef := providerSettingsExecutionSecretKey(consumerID, providerID, requestToken, sourceRef)
	if err := service.secretStore.Save(ctx, snapshotRef, trimmedSecret); err != nil {
		return "", fmt.Errorf("save provider settings execution secret snapshot: %w", err)
	}
	return snapshotRef, nil
}

func providerSettingsContextOrBackground(ctx context.Context) context.Context {
	if ctx != nil {
		return ctx
	}
	return context.TODO()
}

func providerSettingsSpec(providerID string) (providerSettingsProviderSpec, error) {
	normalized := strings.ToLower(strings.TrimSpace(providerID))
	for _, spec := range providerSettingsProviderCatalog {
		if spec.ID == normalized {
			return spec, nil
		}
	}
	return providerSettingsProviderSpec{}, fmt.Errorf("%w: unsupported provider: %s", ErrProviderSettingsValidation, normalized)
}

func withProviderSettingsMutationLock(providerID string, fn func() error) error {
	if fn == nil {
		return nil
	}
	lockKey := strings.ToLower(strings.TrimSpace(providerID))
	lockValue, _ := providerSettingsMutationLocks.LoadOrStore(lockKey, &sync.Mutex{})
	locker := lockValue.(*sync.Mutex) //nolint:forcetypeassert // package-owned sync.Map stores only *sync.Mutex.
	locker.Lock()
	defer locker.Unlock()
	return fn()
}

func normalizeServiceOptionalString(value *string) *string {
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

func providerSettingsStringPointer(value string) *string {
	trimmed := strings.TrimSpace(value)
	if trimmed == "" {
		return nil
	}
	cloned := trimmed
	return &cloned
}

func providerSettingsRequestToken(providerID string, revision int64) string {
	return fmt.Sprintf("%s|%d", strings.TrimSpace(providerID), revision)
}

func providerSettingsSecretKey(providerID string) string {
	return "provider-settings:" + strings.ToLower(strings.TrimSpace(providerID))
}

func providerSettingsExecutionSecretKey(consumerID string, providerID string, requestToken string, sourceRef string) string {
	keyParts := []string{
		"provider-execution",
		strings.ToLower(strings.TrimSpace(consumerID)),
		strings.ToLower(strings.TrimSpace(providerID)),
		strings.TrimSpace(requestToken),
		strings.TrimSpace(sourceRef),
	}
	sum := sha256.Sum256([]byte(strings.Join(keyParts, "|")))
	return fmt.Sprintf("provider-execution:%s:%s:%x", keyParts[1], keyParts[2], sum[:])
}

func (service *ProviderSettingsService) loadRecord(
	ctx context.Context,
	providerID string,
) (repository.ProviderSettingsRecord, bool, error) {
	record, err := service.repository.GetByProviderID(ctx, providerID)
	if err != nil {
		if errors.Is(err, repository.ErrProviderSettingsNotFound) {
			return repository.ProviderSettingsRecord{ProviderID: providerID}, false, nil
		}
		return repository.ProviderSettingsRecord{}, false, fmt.Errorf("load provider settings row: %w", err)
	}
	return record, true, nil
}

func (service *ProviderSettingsService) resolveSecretSavePlan(
	spec providerSettingsProviderSpec,
	current repository.ProviderSettingsRecord,
	currentExists bool,
	input ProviderSettingsSaveInput,
	_ string,
) (*string, string, string, bool, error) {
	if !spec.CredentialRequired {
		return nil, providerSettingsCredentialStateNotRequired, "", false, nil
	}
	if input.APIKeyInputPresent {
		if input.APIKey == nil {
			return nil, "", "", false, fmt.Errorf("%w: secret-aware save input is required when apiKeyInputPresent is true", ErrProviderSettingsValidation)
		}
		trimmed := strings.TrimSpace(*input.APIKey)
		if trimmed == "" {
			return nil, "", "", false, fmt.Errorf("%w: api key is required", ErrProviderSettingsValidation)
		}
		return providerSettingsStringPointer(providerSettingsSecretKey(spec.ID)), providerSettingsCredentialStateConfigured, trimmed, true, nil
	}
	if currentExists && current.CredentialReferenceID != nil &&
		normalizeProviderSettingsCredentialState(current.CredentialState) == providerSettingsCredentialStateConfigured {
		return providerSettingsStringPointer(*current.CredentialReferenceID), providerSettingsCredentialStateConfigured, "", false, nil
	}
	return nil, providerSettingsCredentialStateMissing, "", false, nil
}

func (service *ProviderSettingsService) restoreProviderSettingsRecord(
	ctx context.Context,
	current repository.ProviderSettingsRecord,
	currentExists bool,
) error {
	if !currentExists {
		return nil
	}
	if err := service.transactor.WithTransaction(ctx, func(txCtx context.Context) error {
		_, upsertErr := service.repository.Upsert(txCtx, repository.ProviderSettingsUpsertDraft{
			ProviderID:            current.ProviderID,
			Endpoint:              normalizeServiceOptionalString(current.Endpoint),
			CredentialReferenceID: normalizeServiceOptionalString(current.CredentialReferenceID),
			CredentialState:       normalizeProviderSettingsCredentialState(current.CredentialState),
			ValidationState:       normalizeProviderSettingsValidationState(current.ValidationState),
			RequestToken:          normalizeServiceOptionalString(current.RequestToken),
			LastFailureKind:       normalizeServiceOptionalString(current.LastFailureKind),
			Revision:              current.Revision,
			UpdatedAt:             current.UpdatedAt,
		})
		if upsertErr != nil {
			return fmt.Errorf("restore provider settings row: %w", upsertErr)
		}
		return nil
	}); err != nil {
		return fmt.Errorf("restore provider settings row transaction: %w", err)
	}
	return nil
}

func (service *ProviderSettingsService) getCurrentSummary(ctx context.Context, providerID string) (ProviderSettingsSummary, error) {
	spec, err := providerSettingsSpec(providerID)
	if err != nil {
		return ProviderSettingsSummary{}, err
	}
	record, exists, err := service.loadRecord(ctx, spec.ID)
	if err != nil {
		return ProviderSettingsSummary{}, err
	}
	return service.buildSummary(ctx, spec, record, exists)
}

func (service *ProviderSettingsService) buildSummary(
	_ context.Context,
	spec providerSettingsProviderSpec,
	record repository.ProviderSettingsRecord,
	recordExists bool,
) (ProviderSettingsSummary, error) {
	persistedEndpoint := normalizeServiceOptionalString(record.Endpoint)
	defaultEndpoint := normalizeServiceOptionalString(providerSettingsStringPointer(spec.DefaultEndpoint))
	summary := ProviderSettingsSummary{
		ProviderID:      spec.ID,
		Label:           spec.Label,
		Endpoint:        firstNonNilProviderSettingsString(persistedEndpoint, defaultEndpoint),
		ValidationState: providerSettingsValidationStateNotValidated,
		SavedState:      providerSettingsSavedStateNotSaved,
	}
	if recordExists {
		summary.ValidationState = normalizeProviderSettingsValidationState(record.ValidationState)
		summary.RequestToken = normalizeServiceOptionalString(record.RequestToken)
		summary.LastFailureKind = normalizeServiceOptionalString(record.LastFailureKind)
	}
	if !spec.CredentialRequired {
		summary.CredentialState = providerSettingsCredentialStateNotRequired
		if recordExists {
			summary.SavedState = providerSettingsSavedStatePartial
		}
		if persistedEndpoint != nil {
			summary.SavedState = providerSettingsSavedStateConfigured
		}
		return summary, nil
	}
	credentialRef := normalizeServiceOptionalString(record.CredentialReferenceID)
	if credentialRef == nil {
		summary.CredentialState = providerSettingsCredentialStateMissing
		if recordExists {
			summary.SavedState = providerSettingsSavedStatePartial
		}
		return summary, nil
	}
	credentialState := normalizeProviderSettingsCredentialState(record.CredentialState)
	if credentialState != providerSettingsCredentialStateConfigured {
		summary.CredentialState = providerSettingsCredentialStateMissing
		if recordExists {
			summary.SavedState = providerSettingsSavedStatePartial
		}
		return summary, nil
	}
	summary.CredentialReferenceID = credentialRef
	summary.CredentialState = providerSettingsCredentialStateConfigured
	if persistedEndpoint == nil {
		summary.SavedState = providerSettingsSavedStatePartial
		return summary, nil
	}
	summary.SavedState = providerSettingsSavedStateConfigured
	return summary, nil
}

func normalizeProviderSettingsCredentialState(value string) string {
	switch strings.TrimSpace(value) {
	case providerSettingsCredentialStateConfigured:
		return providerSettingsCredentialStateConfigured
	case providerSettingsCredentialStateNotRequired:
		return providerSettingsCredentialStateNotRequired
	default:
		return providerSettingsCredentialStateMissing
	}
}

func firstNonNilProviderSettingsString(values ...*string) *string {
	for _, value := range values {
		if normalized := normalizeServiceOptionalString(value); normalized != nil {
			return normalized
		}
	}
	return nil
}

func normalizeProviderSettingsValidationState(value string) string {
	switch strings.TrimSpace(value) {
	case providerSettingsValidationStatePending,
		providerSettingsValidationStateValidated,
		providerSettingsValidationStateFailed:
		return strings.TrimSpace(value)
	default:
		return providerSettingsValidationStateNotValidated
	}
}

func providerSettingsMatchesSnapshot(
	summary ProviderSettingsSummary,
	endpoint *string,
	credentialState string,
	credentialReferenceID *string,
	requestToken string,
) bool {
	if strings.TrimSpace(requestToken) == "" {
		return false
	}
	if !strings.EqualFold(strings.TrimSpace(valueOrEmpty(summary.RequestToken)), strings.TrimSpace(requestToken)) {
		return false
	}
	if strings.TrimSpace(valueOrEmpty(summary.Endpoint)) != strings.TrimSpace(valueOrEmpty(normalizeServiceOptionalString(endpoint))) {
		return false
	}
	if strings.TrimSpace(summary.CredentialState) != strings.TrimSpace(credentialState) {
		return false
	}
	return strings.TrimSpace(valueOrEmpty(summary.CredentialReferenceID)) ==
		strings.TrimSpace(valueOrEmpty(normalizeServiceOptionalString(credentialReferenceID)))
}

func valueOrEmpty(value *string) string {
	if value == nil {
		return ""
	}
	return strings.TrimSpace(*value)
}

func (service *ProviderSettingsService) resolveValidationSecret(
	ctx context.Context,
	spec providerSettingsProviderSpec,
	summary ProviderSettingsSummary,
) (string, *string, error) {
	if !spec.CredentialRequired {
		return "", nil, nil
	}
	if summary.CredentialReferenceID == nil {
		return "", providerSettingsStringPointer(providerSettingsErrorKindCredentialMissing), nil
	}
	loaded, err := service.secretStore.Load(ctx, *summary.CredentialReferenceID)
	if err != nil {
		return "", nil, fmt.Errorf("load provider settings validation secret: %w", err)
	}
	trimmed := strings.TrimSpace(loaded)
	if trimmed == "" {
		return "", providerSettingsStringPointer(providerSettingsErrorKindCredentialMissing), nil
	}
	return trimmed, nil, nil
}
