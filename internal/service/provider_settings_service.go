package service

import (
	"context"
	"crypto/sha256"
	"errors"
	"fmt"
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
		now := service.now().UTC()
		current, currentExists, err := service.loadRecord(providerSettingsContextOrBackground(ctx), spec.ID)
		if err != nil {
			return err
		}

		revision := current.Revision + 1
		requestToken := providerSettingsRequestToken(spec.ID, revision)
		secretKey := providerSettingsSecretKey(spec.ID)
		secretStageKey := providerSettingsStagedSecretKey(spec.ID, requestToken)

		credentialRef, credentialState, secretValue, secretChanged, err := service.resolveSecretSavePlan(
			spec,
			current,
			currentExists,
			input,
			requestToken,
		)
		if err != nil {
			return err
		}

		draft := repository.ProviderSettingsUpsertDraft{
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

		if secretChanged {
			if saveErr := service.secretStore.Save(providerSettingsContextOrBackground(ctx), secretStageKey, secretValue); saveErr != nil {
				return fmt.Errorf("stage provider settings secret: %w", saveErr)
			}
			defer func() {
				_ = service.secretStore.Delete(providerSettingsContextOrBackground(ctx), secretStageKey)
			}()
		}

		var savedRecord repository.ProviderSettingsRecord
		txErr := service.transactor.WithTransaction(providerSettingsContextOrBackground(ctx), func(txCtx context.Context) error {
			record, upsertErr := service.repository.Upsert(txCtx, draft)
			if upsertErr != nil {
				return fmt.Errorf("save provider settings row: %w", upsertErr)
			}
			savedRecord = record
			return nil
		})
		if txErr != nil {
			return fmt.Errorf("save provider settings: %w", txErr)
		}
		if secretChanged {
			if saveErr := service.secretStore.Save(providerSettingsContextOrBackground(ctx), secretKey, secretValue); saveErr != nil {
				if rollbackErr := service.restoreProviderSettingsRecord(providerSettingsContextOrBackground(ctx), current, currentExists); rollbackErr != nil {
					return fmt.Errorf("save provider settings canonical secret: %w; rollback provider settings row: %w", saveErr, rollbackErr)
				}
				return fmt.Errorf("save provider settings canonical secret: %w", saveErr)
			}
		}
		saved, err = service.buildSummary(providerSettingsContextOrBackground(ctx), spec, savedRecord, true)
		return err
	})
	if lockErr != nil {
		return ProviderSettingsSummary{}, lockErr
	}
	return saved, nil
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
		now := service.now().UTC()
		current, currentExists, err := service.loadRecord(providerSettingsContextOrBackground(ctx), spec.ID)
		if err != nil {
			return err
		}
		secretKey := providerSettingsSecretKey(spec.ID)

		credentialState := providerSettingsCredentialStateMissing
		if !spec.CredentialRequired {
			credentialState = providerSettingsCredentialStateNotRequired
		}
		revision := current.Revision + 1
		requestToken := providerSettingsRequestToken(spec.ID, revision)
		draft := repository.ProviderSettingsUpsertDraft{
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

		var savedRecord repository.ProviderSettingsRecord
		txErr := service.transactor.WithTransaction(providerSettingsContextOrBackground(ctx), func(txCtx context.Context) error {
			record, upsertErr := service.repository.Upsert(txCtx, draft)
			if upsertErr != nil {
				return fmt.Errorf("reset provider settings row: %w", upsertErr)
			}
			savedRecord = record
			return nil
		})
		if txErr != nil {
			return fmt.Errorf("reset provider settings: %w", txErr)
		}
		if deleteErr := service.secretStore.Delete(providerSettingsContextOrBackground(ctx), secretKey); deleteErr != nil {
			if rollbackErr := service.restoreProviderSettingsRecord(providerSettingsContextOrBackground(ctx), current, currentExists); rollbackErr != nil {
				return fmt.Errorf("delete provider settings secret: %w; rollback provider settings row: %w", deleteErr, rollbackErr)
			}
			return fmt.Errorf("delete provider settings secret: %w", deleteErr)
		}
		reset, err = service.buildSummary(providerSettingsContextOrBackground(ctx), spec, savedRecord, true)
		return err
	})
	if lockErr != nil {
		return ProviderSettingsSummary{}, lockErr
	}
	reset.Endpoint = nil
	return reset, nil
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
		summary, err := service.getCurrentSummary(providerSettingsContextOrBackground(ctx), spec.ID)
		if err != nil {
			return err
		}
		requestToken := strings.TrimSpace(input.RequestToken)
		if requestToken == "" {
			return fmt.Errorf("%w: request token is required", ErrProviderSettingsValidation)
		}
		if !providerSettingsMatchesSnapshot(summary, input.Endpoint, input.CredentialState, input.CredentialReferenceID, requestToken) {
			result = ProviderSettingsValidationResult{
				ProviderID:      spec.ID,
				ValidationState: providerSettingsValidationStateNotValidated,
				RequestToken:    requestToken,
				FailureKind:     providerSettingsStringPointer(providerSettingsErrorKindValidationStale),
			}
			return nil
		}
		if summary.Endpoint == nil {
			result = ProviderSettingsValidationResult{
				ProviderID:      spec.ID,
				ValidationState: providerSettingsValidationStateNotValidated,
				RequestToken:    requestToken,
				FailureKind:     providerSettingsStringPointer(providerSettingsErrorKindEndpointMissing),
			}
			return nil
		}
		if spec.CredentialRequired && summary.CredentialReferenceID == nil {
			result = ProviderSettingsValidationResult{
				ProviderID:      spec.ID,
				ValidationState: providerSettingsValidationStateNotValidated,
				RequestToken:    requestToken,
				FailureKind:     providerSettingsStringPointer(providerSettingsErrorKindCredentialMissing),
			}
			return nil
		}
		apiKey, failureKind, err := service.resolveValidationSecret(providerSettingsContextOrBackground(ctx), spec, summary)
		if err != nil {
			return err
		}
		if failureKind != nil {
			result = ProviderSettingsValidationResult{
				ProviderID:      spec.ID,
				ValidationState: providerSettingsValidationStateNotValidated,
				RequestToken:    requestToken,
				FailureKind:     failureKind,
			}
			return nil
		}

		if _, matched, updateErr := service.repository.UpdateValidationByRequestToken(
			providerSettingsContextOrBackground(ctx),
			spec.ID,
			requestToken,
			providerSettingsValidationStatePending,
			nil,
		); updateErr != nil {
			return fmt.Errorf("mark provider settings validation pending: %w", updateErr)
		} else if !matched {
			result = ProviderSettingsValidationResult{
				ProviderID:      spec.ID,
				ValidationState: providerSettingsValidationStateNotValidated,
				RequestToken:    requestToken,
				FailureKind:     providerSettingsStringPointer(providerSettingsErrorKindValidationStale),
			}
			return nil
		}

		validationErr := service.validator.ValidateProviderSettings(providerSettingsContextOrBackground(ctx), ProviderSettingsValidationProbe{
			ProviderID: spec.ID,
			Endpoint:   *summary.Endpoint,
			APIKey:     apiKey,
		})
		finalState := providerSettingsValidationStateValidated
		var finalFailure *string
		if validationErr != nil {
			finalState = providerSettingsValidationStateFailed
			finalFailure = providerSettingsStringPointer(providerSettingsErrorKindProviderUnreachable)
		}
		_, matched, updateErr := service.repository.UpdateValidationByRequestToken(
			providerSettingsContextOrBackground(ctx),
			spec.ID,
			requestToken,
			finalState,
			finalFailure,
		)
		if updateErr != nil {
			return fmt.Errorf("update provider settings validation result: %w", updateErr)
		}
		if !matched {
			result = ProviderSettingsValidationResult{
				ProviderID:      spec.ID,
				ValidationState: providerSettingsValidationStateNotValidated,
				RequestToken:    requestToken,
				FailureKind:     providerSettingsStringPointer(providerSettingsErrorKindValidationStale),
			}
			return nil
		}
		result = ProviderSettingsValidationResult{
			ProviderID:      spec.ID,
			ValidationState: finalState,
			RequestToken:    requestToken,
			FailureKind:     finalFailure,
		}
		return nil
	})
	if lockErr != nil {
		return ProviderSettingsValidationResult{}, lockErr
	}
	return result, nil
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
		summary, err := service.getCurrentSummary(providerSettingsContextOrBackground(ctx), spec.ID)
		if err != nil {
			return err
		}
		requestToken := strings.TrimSpace(input.RequestToken)
		result = ProviderSettingsModelListResult{
			ProviderID:      spec.ID,
			Endpoint:        summary.Endpoint,
			CredentialState: summary.CredentialState,
			RequestToken:    requestToken,
			State:           providerSettingsModelListStateNotRequested,
			Models:          []ProviderSettingsModelOption{},
		}
		if !providerSettingsMatchesSnapshot(summary, input.Endpoint, input.CredentialState, input.CredentialReferenceID, requestToken) {
			result.State = providerSettingsModelListStateFailed
			result.FailureKind = providerSettingsStringPointer(providerSettingsErrorKindValidationStale)
			return nil
		}
		if summary.Endpoint == nil {
			result.State = providerSettingsModelListStateEndpointMissing
			result.FailureKind = providerSettingsStringPointer(providerSettingsErrorKindEndpointMissing)
			return nil
		}
		apiKey := ""
		if spec.CredentialRequired {
			if summary.CredentialReferenceID == nil {
				result.State = providerSettingsModelListStateCredentialMissing
				result.FailureKind = providerSettingsStringPointer(providerSettingsErrorKindCredentialMissing)
				return nil
			}
			loaded, failureKind, loadErr := service.resolveValidationSecret(providerSettingsContextOrBackground(ctx), spec, summary)
			if loadErr != nil {
				return loadErr
			}
			if failureKind != nil {
				result.State = providerSettingsModelListStateCredentialMissing
				result.FailureKind = failureKind
				return nil
			}
			apiKey = loaded
		} else {
			result.State = providerSettingsModelListStateCredentialNotNeeded
		}
		models, loadErr := service.modelListLoader.ListProviderModelsWithEndpoint(
			providerSettingsContextOrBackground(ctx),
			spec.ID,
			*summary.Endpoint,
			apiKey,
		)
		if loadErr != nil {
			result.State = providerSettingsModelListStateFailed
			result.FailureKind = providerSettingsStringPointer(providerSettingsErrorKindModelListFailed)
			//nolint:nilerr // public contract intentionally returns a redacted failure instead of the transport error.
			return nil
		}
		result.State = providerSettingsModelListStateReady
		if !spec.CredentialRequired {
			result.State = providerSettingsModelListStateCredentialNotNeeded
		}
		result.Models = make([]ProviderSettingsModelOption, 0, len(models))
		for _, model := range models {
			result.Models = append(result.Models, ProviderSettingsModelOption{
				ModelID: strings.TrimSpace(model.ModelID),
				Label:   strings.TrimSpace(model.Label),
			})
		}
		sort.Slice(result.Models, func(left, right int) bool {
			return result.Models[left].ModelID < result.Models[right].ModelID
		})
		return nil
	})
	if lockErr != nil {
		return ProviderSettingsModelListResult{}, lockErr
	}
	return result, nil
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
		summary, err := service.getCurrentSummary(providerSettingsContextOrBackground(ctx), spec.ID)
		if err != nil {
			return ProviderSettingsResolveResult{}, err
		}
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
		if input.AllowSecretSnapshot && result.ErrorKind == nil && spec.CredentialRequired && summary.CredentialReferenceID != nil {
			snapshotRef, snapshotErr := service.snapshotExecutionSecret(
				providerSettingsContextOrBackground(ctx),
				strings.TrimSpace(input.ConsumerID),
				spec.ID,
				summary,
			)
			if snapshotErr != nil {
				return ProviderSettingsResolveResult{}, snapshotErr
			}
			result.CredentialReferenceID = providerSettingsStringPointer(snapshotRef)
		}
		return result, nil
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
	return context.Background()
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

func providerSettingsStagedSecretKey(providerID string, requestToken string) string {
	keyParts := []string{
		"provider-settings-stage",
		strings.ToLower(strings.TrimSpace(providerID)),
		strings.TrimSpace(requestToken),
	}
	sum := sha256.Sum256([]byte(strings.Join(keyParts, "|")))
	return fmt.Sprintf("provider-settings-stage:%s:%x", keyParts[1], sum[:])
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
