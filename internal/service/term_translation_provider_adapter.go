package service

import (
	"context"
	"errors"
	"fmt"
	"reflect"
	"sort"
	"strings"
)

const (
	// TermTranslationProviderGemini defines the supported Gemini provider id.
	TermTranslationProviderGemini = "gemini"
	// TermTranslationProviderLMStudio defines the supported LM Studio provider id.
	TermTranslationProviderLMStudio = "lm_studio"
	// TermTranslationProviderXAI defines the supported xAI provider id.
	TermTranslationProviderXAI = "xai"

	// TermTranslationExecutionModeSingleRequest identifies the one-term-per-request mode.
	TermTranslationExecutionModeSingleRequest = "single_request"

	termTranslationProviderResponseInvalidReason = "provider response is invalid"
)

var termTranslationSupportedProviderSet = map[string]struct{}{
	TermTranslationProviderGemini:   {},
	TermTranslationProviderLMStudio: {},
	TermTranslationProviderXAI:      {},
}

// TermTranslationProviderErrorKind identifies provider adapter failure families.
type TermTranslationProviderErrorKind string

const (
	// TermTranslationProviderErrorKindProviderFailure identifies provider execution failures.
	TermTranslationProviderErrorKindProviderFailure TermTranslationProviderErrorKind = "provider_failure"
	// TermTranslationProviderErrorKindInvalidProviderResponse identifies invalid provider response shapes.
	TermTranslationProviderErrorKindInvalidProviderResponse TermTranslationProviderErrorKind = "invalid_provider_response"
)

// TermTranslationProvider defines the provider-agnostic one-term translation port.
type TermTranslationProvider interface {
	TranslateTerm(ctx context.Context, request TermTranslationProviderRequest) (TermTranslationProviderResult, error)
	TermTranslationProviderRequestsAreTestSafe() bool
}

// TermTranslationProviderClientResponse defines the lower-level provider client reply shape.
type TermTranslationProviderClientResponse struct {
	Items         []TermTranslationProviderClientItem
	ExecutionMode string
}

// TermTranslationProviderClientItem defines one correlated provider client reply item.
type TermTranslationProviderClientItem struct {
	SourceTerm     string
	TranslatedTerm string
}

// TermTranslationProviderRequest defines one source term translation request unit.
type TermTranslationProviderRequest struct {
	Provider        string
	Model           string
	APIKey          string
	EndpointSummary *string
	SourceTerm      string
	SourceLanguage  string
	TargetLanguage  string
	PromptVersion   string
	PromptDigest    string
}

// TermTranslationProviderAuditSummary exposes provider execution metadata without secrets.
type TermTranslationProviderAuditSummary struct {
	Provider      string
	Model         string
	ExecutionMode string
	InputCount    int
	OutputCount   int
	PromptVersion *string
	PromptDigest  *string
}

// TermTranslationProviderResult defines one correlated translation result.
type TermTranslationProviderResult struct {
	SourceTerm      string
	TranslatedTerm  string
	Confirmed       bool
	Failure         *TermTranslationProviderFailure
	AuditSummary    TermTranslationProviderAuditSummary
	ProviderSkipped bool
}

// TermTranslationProviderFailure exposes only redacted provider failure information.
type TermTranslationProviderFailure struct {
	Kind       TermTranslationProviderErrorKind
	Reason     string
	Retryable  bool
	IsRedacted bool
}

// TermTranslationProviderError wraps one provider adapter failure with redacted public details.
type TermTranslationProviderError struct {
	Failure TermTranslationProviderFailure
	cause   error
}

type termTranslationFailureMetadata interface {
	error
	FailureKind() string
	FailureRetryable() bool
}

func (err *TermTranslationProviderError) Error() string {
	if err == nil {
		return ""
	}
	return err.Failure.Reason
}

// Unwrap returns the underlying cause for internal diagnostics.
func (err *TermTranslationProviderError) Unwrap() error {
	if err == nil {
		return nil
	}
	return err.cause
}

// NewTermTranslationProviderError creates a typed provider adapter error.
func NewTermTranslationProviderError(
	kind TermTranslationProviderErrorKind,
	reason string,
	retryable bool,
	cause error,
) error {
	return &TermTranslationProviderError{
		Failure: TermTranslationProviderFailure{
			Kind:       kind,
			Reason:     strings.TrimSpace(reason),
			Retryable:  retryable,
			IsRedacted: true,
		},
		cause: cause,
	}
}

// NewTermTranslationProviderAdapter creates the one-term provider adapter on the service boundary.
func NewTermTranslationProviderAdapter(client any) TermTranslationProvider {
	return termTranslationProviderAdapter{client: client}
}

// NormalizeTermTranslationProvider validates and normalizes the provider id.
func NormalizeTermTranslationProvider(provider string) (string, error) {
	normalized := strings.ToLower(strings.TrimSpace(provider))
	if normalized == "" {
		return "", fmt.Errorf("term translation provider is required")
	}
	if _, ok := termTranslationSupportedProviderSet[normalized]; !ok {
		return "", fmt.Errorf("unsupported term translation provider: %s", normalized)
	}
	return normalized, nil
}

// TermTranslationSupportedProviders returns the backend-supported term translation providers.
func TermTranslationSupportedProviders() []string {
	providers := make([]string, 0, len(termTranslationSupportedProviderSet))
	for provider := range termTranslationSupportedProviderSet {
		providers = append(providers, provider)
	}
	sort.Strings(providers)
	return providers
}

// BuildTermTranslationPrompt returns the strict JSON-only prompt for one source term request unit.
func BuildTermTranslationPrompt(request TermTranslationProviderRequest) (string, error) {
	sourceTerm := strings.TrimSpace(request.SourceTerm)
	if sourceTerm == "" {
		return "", fmt.Errorf("term translation source term is required")
	}
	sourceLanguage := strings.TrimSpace(request.SourceLanguage)
	if sourceLanguage == "" {
		sourceLanguage = "source"
	}
	targetLanguage := strings.TrimSpace(request.TargetLanguage)
	if targetLanguage == "" {
		targetLanguage = "target"
	}
	return strings.TrimSpace(strings.Join([]string{
		"TERM_TRANSLATION_REQUEST_V1",
		"Return strict JSON only.",
		`Use the exact shape {"translations":[{"source_term":"...","translated_term":"..."}]}.`,
		"Do not add markdown, commentary, or extra keys.",
		"input_count=1",
		"execution_mode=" + TermTranslationExecutionModeSingleRequest,
		"source_language=" + sourceLanguage,
		"target_language=" + targetLanguage,
		"source_term=" + sourceTerm,
	}, "\n")), nil
}

type termTranslationProviderAdapter struct {
	client any
}

func (adapter termTranslationProviderAdapter) TermTranslationProviderRequestsAreTestSafe() bool {
	if adapter.client == nil {
		return false
	}
	method := reflect.ValueOf(adapter.client).MethodByName("ProviderRequestsAreTestSafe")
	if !method.IsValid() || method.Type().NumIn() != 0 || method.Type().NumOut() != 1 || method.Type().Out(0).Kind() != reflect.Bool {
		return false
	}
	return method.Call(nil)[0].Bool()
}

func (adapter termTranslationProviderAdapter) TranslateTerm(
	ctx context.Context,
	request TermTranslationProviderRequest,
) (TermTranslationProviderResult, error) {
	if adapter.client == nil {
		return TermTranslationProviderResult{}, NewTermTranslationProviderError(
			TermTranslationProviderErrorKindProviderFailure,
			"provider request could not start",
			false,
			fmt.Errorf("term translation provider client is required"),
		)
	}

	providerID, err := NormalizeTermTranslationProvider(request.Provider)
	if err != nil {
		return TermTranslationProviderResult{}, NewTermTranslationProviderError(
			TermTranslationProviderErrorKindProviderFailure,
			"provider configuration is invalid",
			false,
			err,
		)
	}

	model := strings.TrimSpace(request.Model)
	if model == "" {
		return TermTranslationProviderResult{}, NewTermTranslationProviderError(
			TermTranslationProviderErrorKindProviderFailure,
			"provider configuration is invalid",
			false,
			fmt.Errorf("term translation model is required"),
		)
	}

	sourceTerm := strings.TrimSpace(request.SourceTerm)
	if sourceTerm == "" {
		return TermTranslationProviderResult{}, NewTermTranslationProviderError(
			TermTranslationProviderErrorKindInvalidProviderResponse,
			termTranslationProviderResponseInvalidReason,
			false,
			fmt.Errorf("term translation source term is required"),
		)
	}

	prompt, err := BuildTermTranslationPrompt(request)
	if err != nil {
		return TermTranslationProviderResult{}, NewTermTranslationProviderError(
			TermTranslationProviderErrorKindProviderFailure,
			"provider request could not start",
			false,
			err,
		)
	}

	clientResponse, err := invokeTermTranslationClientTranslateTerm(
		ctx,
		adapter.client,
		providerID,
		model,
		strings.TrimSpace(request.APIKey),
		providerExecutionOptionalString(request.EndpointSummary),
		prompt,
	)
	if err != nil {
		return TermTranslationProviderResult{}, mapTermTranslationProviderError(err)
	}
	if len(clientResponse.Items) != 1 {
		return TermTranslationProviderResult{}, NewTermTranslationProviderError(
			TermTranslationProviderErrorKindInvalidProviderResponse,
			termTranslationProviderResponseInvalidReason,
			true,
			fmt.Errorf("provider returned %d translation items for one source term", len(clientResponse.Items)),
		)
	}

	item := clientResponse.Items[0]
	if strings.TrimSpace(item.SourceTerm) != sourceTerm {
		return TermTranslationProviderResult{}, NewTermTranslationProviderError(
			TermTranslationProviderErrorKindInvalidProviderResponse,
			termTranslationProviderResponseInvalidReason,
			true,
			fmt.Errorf("provider returned mismatched source term"),
		)
	}
	if strings.TrimSpace(item.TranslatedTerm) == "" {
		return TermTranslationProviderResult{}, NewTermTranslationProviderError(
			TermTranslationProviderErrorKindInvalidProviderResponse,
			termTranslationProviderResponseInvalidReason,
			true,
			fmt.Errorf("provider returned empty translated term"),
		)
	}

	return TermTranslationProviderResult{
		SourceTerm:     sourceTerm,
		TranslatedTerm: strings.TrimSpace(item.TranslatedTerm),
		Confirmed:      true,
		AuditSummary: TermTranslationProviderAuditSummary{
			Provider:      providerID,
			Model:         model,
			ExecutionMode: strings.TrimSpace(clientResponse.ExecutionMode),
			InputCount:    1,
			OutputCount:   len(clientResponse.Items),
			PromptVersion: optionalStringPointer(request.PromptVersion),
			PromptDigest:  optionalStringPointer(request.PromptDigest),
		},
	}, nil
}

func invokeTermTranslationClientTranslateTerm(
	ctx context.Context,
	client any,
	providerID string,
	model string,
	apiKey string,
	endpointSummary string,
	prompt string,
) (TermTranslationProviderClientResponse, error) {
	if client == nil {
		return TermTranslationProviderClientResponse{}, fmt.Errorf("term translation provider client is required")
	}
	method := reflect.ValueOf(client).MethodByName("TranslateTerm")
	if !method.IsValid() {
		return TermTranslationProviderClientResponse{}, fmt.Errorf("term translation provider client does not implement TranslateTerm")
	}
	if (method.Type().NumIn() != 5 && method.Type().NumIn() != 6) || method.Type().NumOut() != 2 {
		return TermTranslationProviderClientResponse{}, fmt.Errorf("term translation provider client has incompatible TranslateTerm signature")
	}
	args := []reflect.Value{
		reflect.ValueOf(ctx),
		reflect.ValueOf(providerID),
		reflect.ValueOf(model),
		reflect.ValueOf(apiKey),
	}
	if method.Type().NumIn() == 6 {
		args = append(args, reflect.ValueOf(strings.TrimSpace(endpointSummary)))
	}
	args = append(args, reflect.ValueOf(prompt))
	results := method.Call(args)
	if errValue := results[1]; !errValue.IsNil() {
		err, _ := errValue.Interface().(error)
		return TermTranslationProviderClientResponse{}, err
	}
	return mapTermTranslationClientResponse(results[0])
}

func mapTermTranslationClientResponse(
	value reflect.Value,
) (TermTranslationProviderClientResponse, error) {
	if !value.IsValid() {
		return TermTranslationProviderClientResponse{}, fmt.Errorf("term translation provider response is invalid")
	}
	if value.Kind() == reflect.Pointer {
		if value.IsNil() {
			return TermTranslationProviderClientResponse{}, fmt.Errorf("term translation provider response is nil")
		}
		value = value.Elem()
	}
	itemsField := value.FieldByName("Items")
	executionModeField := value.FieldByName("ExecutionMode")
	if !itemsField.IsValid() || itemsField.Kind() != reflect.Slice || !executionModeField.IsValid() || executionModeField.Kind() != reflect.String {
		return TermTranslationProviderClientResponse{}, fmt.Errorf("term translation provider response has incompatible shape")
	}
	items := make([]TermTranslationProviderClientItem, 0, itemsField.Len())
	for index := 0; index < itemsField.Len(); index++ {
		itemValue := itemsField.Index(index)
		if itemValue.Kind() == reflect.Pointer {
			if itemValue.IsNil() {
				return TermTranslationProviderClientResponse{}, fmt.Errorf("term translation provider item is nil")
			}
			itemValue = itemValue.Elem()
		}
		sourceTermField := itemValue.FieldByName("SourceTerm")
		translatedTermField := itemValue.FieldByName("TranslatedTerm")
		if !sourceTermField.IsValid() || sourceTermField.Kind() != reflect.String || !translatedTermField.IsValid() || translatedTermField.Kind() != reflect.String {
			return TermTranslationProviderClientResponse{}, fmt.Errorf("term translation provider item has incompatible shape")
		}
		items = append(items, TermTranslationProviderClientItem{
			SourceTerm:     sourceTermField.String(),
			TranslatedTerm: translatedTermField.String(),
		})
	}
	return TermTranslationProviderClientResponse{
		Items:         items,
		ExecutionMode: executionModeField.String(),
	}, nil
}

func mapTermTranslationProviderError(err error) error {
	var metadata termTranslationFailureMetadata
	if !errors.As(err, &metadata) {
		return NewTermTranslationProviderError(
			TermTranslationProviderErrorKindProviderFailure,
			"provider request failed",
			true,
			err,
		)
	}
	switch metadata.FailureKind() {
	case string(TermTranslationProviderErrorKindInvalidProviderResponse):
		return NewTermTranslationProviderError(
			TermTranslationProviderErrorKindInvalidProviderResponse,
			termTranslationProviderResponseInvalidReason,
			metadata.FailureRetryable(),
			err,
		)
	default:
		return NewTermTranslationProviderError(
			TermTranslationProviderErrorKindProviderFailure,
			redactedTermTranslationProviderFailureReason(err),
			metadata.FailureRetryable(),
			err,
		)
	}
}

func redactedTermTranslationProviderFailureReason(err error) string {
	lowerMessage := strings.ToLower(strings.TrimSpace(err.Error()))
	if strings.Contains(lowerMessage, "api key") || strings.Contains(lowerMessage, "authorization") {
		return "provider credential is unavailable"
	}
	return "provider request failed"
}

func optionalStringPointer(value string) *string {
	trimmed := strings.TrimSpace(value)
	if trimmed == "" {
		return nil
	}
	return &trimmed
}
