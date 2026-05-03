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
	// BodyTranslationProviderGemini defines the supported Gemini provider id.
	BodyTranslationProviderGemini = "gemini"
	// BodyTranslationProviderLMStudio defines the supported LM Studio provider id.
	BodyTranslationProviderLMStudio = "lm_studio"
	// BodyTranslationProviderXAI defines the supported xAI provider id.
	BodyTranslationProviderXAI = "xai"
	// BodyTranslationProviderFake defines the supported test-only fake provider id.
	BodyTranslationProviderFake = "fake"

	bodyTranslationProviderResponseInvalidReason = "provider response is invalid"
	bodyTranslationInvalidConfigurationReason    = "provider configuration is invalid"
)

var bodyTranslationSupportedProviderSet = map[string]struct{}{
	BodyTranslationProviderFake:     {},
	BodyTranslationProviderGemini:   {},
	BodyTranslationProviderLMStudio: {},
	BodyTranslationProviderXAI:      {},
}

// BodyTranslationProviderErrorKind identifies provider adapter failure families.
type BodyTranslationProviderErrorKind string

const (
	// BodyTranslationProviderErrorKindProviderFailure identifies provider execution failures.
	BodyTranslationProviderErrorKindProviderFailure BodyTranslationProviderErrorKind = "provider_failure"
	// BodyTranslationProviderErrorKindInvalidProviderResponse identifies invalid response shapes.
	BodyTranslationProviderErrorKindInvalidProviderResponse BodyTranslationProviderErrorKind = "invalid_provider_response"
)

// BodyTranslationProvider defines the provider-agnostic one-field translation port.
type BodyTranslationProvider interface {
	TranslateBodyField(ctx context.Context, request BodyTranslationProviderRequest) BodyTranslationProviderResult
	BodyTranslationProviderRequestsAreTestSafe() bool
}

// BodyTranslationProviderRequest defines one body translation request unit.
type BodyTranslationProviderRequest struct {
	Provider                string
	Model                   string
	ExecutionMode           string
	CredentialRef           string
	RequestUnitID           string
	FieldCorrelationKey     string
	RecordType              string
	FieldType               string
	SourceText              string
	SourceLanguage          string
	TargetLanguage          string
	PersonaSummary          string
	ContextLines            []string
	CompleteMatchExclusions []BodyTranslationDictionaryExactMatchExclusion
	PartialMatchConstraints []BodyTranslationPartialMatchConstraint
	ProtectedElements       []BodyTranslationProtectedElement
}

// BodyTranslationProviderAuditSummary exposes provider execution metadata without secrets or raw prompts.
type BodyTranslationProviderAuditSummary struct {
	CredentialRef  string
	Provider       string
	Model          string
	ExecutionMode  string
	RequestUnitID  string
	PromptDigest   string
	InputCount     int
	OutputCount    int
	RequestSummary BodyTranslationProviderRequestSummary
}

// BodyTranslationProviderFailure exposes only redacted provider failure information.
type BodyTranslationProviderFailure struct {
	Kind       BodyTranslationProviderErrorKind
	Reason     string
	Retryable  bool
	IsRedacted bool
}

// BodyTranslationTranslatedCandidate is the translated output candidate that may be persisted later.
type BodyTranslationTranslatedCandidate struct {
	RequestUnitID       string
	FieldCorrelationKey string
	RecordType          string
	FieldType           string
	TranslatedText      string
}

// BodyTranslationProtectionValidationTarget identifies the payload later protection validation should inspect.
type BodyTranslationProtectionValidationTarget struct {
	RequestUnitID          string
	FieldCorrelationKey    string
	ProtectionSourceDigest string
	TranslatedText         string
}

// BodyTranslationProviderResult defines one correlated body translation result.
type BodyTranslationProviderResult struct {
	RequestUnitID              string
	FieldCorrelationKey        string
	RecordType                 string
	FieldType                  string
	TranslatedCandidate        *BodyTranslationTranslatedCandidate
	ProtectionValidationTarget *BodyTranslationProtectionValidationTarget
	Failure                    *BodyTranslationProviderFailure
	AuditSummary               BodyTranslationProviderAuditSummary
}

// NewBodyTranslationProviderAdapter creates the one-field provider adapter on the service boundary.
func NewBodyTranslationProviderAdapter(client any) BodyTranslationProvider {
	return bodyTranslationProviderAdapter{client: client}
}

// NormalizeBodyTranslationProvider validates and normalizes the provider id.
func NormalizeBodyTranslationProvider(provider string) (string, error) {
	normalized := strings.ToLower(strings.TrimSpace(provider))
	if normalized == "" {
		return "", fmt.Errorf("body translation provider is required")
	}
	if _, ok := bodyTranslationSupportedProviderSet[normalized]; !ok {
		return "", fmt.Errorf("unsupported body translation provider: %s", normalized)
	}
	return normalized, nil
}

// BodyTranslationSupportedProviders returns the backend-supported body translation providers.
func BodyTranslationSupportedProviders() []string {
	providers := make([]string, 0, len(bodyTranslationSupportedProviderSet))
	for provider := range bodyTranslationSupportedProviderSet {
		providers = append(providers, provider)
	}
	sort.Strings(providers)
	return providers
}

type bodyTranslationProviderAdapter struct {
	client any
}

func (adapter bodyTranslationProviderAdapter) BodyTranslationProviderRequestsAreTestSafe() bool {
	if adapter.client == nil {
		return false
	}
	method := reflectValueOfClientMethod(adapter.client, "ProviderRequestsAreTestSafe")
	if !method.IsValid() || method.Type().NumIn() != 0 || method.Type().NumOut() != 1 || method.Type().Out(0).Kind() != reflect.Bool {
		return false
	}
	return method.Call(nil)[0].Bool()
}

func (adapter bodyTranslationProviderAdapter) TranslateBodyField(
	ctx context.Context,
	request BodyTranslationProviderRequest,
) BodyTranslationProviderResult {
	baseResult := BodyTranslationProviderResult{
		RequestUnitID:       strings.TrimSpace(request.RequestUnitID),
		FieldCorrelationKey: strings.TrimSpace(request.FieldCorrelationKey),
		RecordType:          strings.TrimSpace(request.RecordType),
		FieldType:           strings.TrimSpace(request.FieldType),
		AuditSummary: BodyTranslationProviderAuditSummary{
			CredentialRef:  strings.TrimSpace(request.CredentialRef),
			Provider:       strings.ToLower(strings.TrimSpace(request.Provider)),
			Model:          strings.TrimSpace(request.Model),
			ExecutionMode:  strings.ToLower(strings.TrimSpace(request.ExecutionMode)),
			RequestUnitID:  strings.TrimSpace(request.RequestUnitID),
			InputCount:     1,
			RequestSummary: buildBodyTranslationRequestSummary(request),
		},
	}

	if adapter.client == nil {
		return bodyTranslationProviderFailureResult(
			baseResult,
			BodyTranslationProviderErrorKindProviderFailure,
			"provider request could not start",
			false,
		)
	}

	providerID, err := NormalizeBodyTranslationProvider(request.Provider)
	if err != nil {
		return bodyTranslationProviderFailureResult(
			baseResult,
			BodyTranslationProviderErrorKindProviderFailure,
			bodyTranslationInvalidConfigurationReason,
			false,
		)
	}
	baseResult.AuditSummary.Provider = providerID

	model := strings.TrimSpace(request.Model)
	if model == "" {
		return bodyTranslationProviderFailureResult(
			baseResult,
			BodyTranslationProviderErrorKindProviderFailure,
			bodyTranslationInvalidConfigurationReason,
			false,
		)
	}
	baseResult.AuditSummary.Model = model

	executionMode, err := normalizeBodyTranslationExecutionMode(request.ExecutionMode)
	if err != nil {
		return bodyTranslationProviderFailureResult(
			baseResult,
			BodyTranslationProviderErrorKindProviderFailure,
			bodyTranslationInvalidConfigurationReason,
			false,
		)
	}
	baseResult.AuditSummary.ExecutionMode = executionMode

	prompt, err := BuildBodyTranslationPrompt(request)
	if err != nil {
		return bodyTranslationProviderFailureResult(
			baseResult,
			BodyTranslationProviderErrorKindProviderFailure,
			bodyTranslationInvalidConfigurationReason,
			false,
		)
	}
	baseResult.AuditSummary.PromptDigest = bodyTranslationPromptDigest(prompt)

	clientResponse, err := invokeBodyTranslationClientGenerateBodyTranslation(
		ctx,
		adapter.client,
		providerID,
		model,
		executionMode,
		strings.TrimSpace(request.CredentialRef),
		prompt,
	)
	if err != nil {
		return mapBodyTranslationProviderFailure(baseResult, err)
	}
	baseResult.AuditSummary.ExecutionMode = firstNonEmptyBodyTranslationValue(clientResponse.ExecutionMode, executionMode)
	baseResult.AuditSummary.PromptDigest = firstNonEmptyBodyTranslationValue(clientResponse.PromptDigest, baseResult.AuditSummary.PromptDigest)
	baseResult.AuditSummary.OutputCount = len(clientResponse.Items)
	return mapBodyTranslationProviderResponse(baseResult, clientResponse)
}

func mapBodyTranslationProviderFailure(
	baseResult BodyTranslationProviderResult,
	err error,
) BodyTranslationProviderResult {
	var metadata bodyTranslationFailureMetadata
	if !errors.As(err, &metadata) {
		return bodyTranslationProviderFailureResult(
			baseResult,
			BodyTranslationProviderErrorKindProviderFailure,
			redactedBodyTranslationProviderFailureReason(err),
			true,
		)
	}
	switch metadata.FailureKind() {
	case string(BodyTranslationProviderErrorKindInvalidProviderResponse):
		return bodyTranslationProviderFailureResult(
			baseResult,
			BodyTranslationProviderErrorKindInvalidProviderResponse,
			bodyTranslationProviderResponseInvalidReason,
			metadata.FailureRetryable(),
		)
	default:
		return bodyTranslationProviderFailureResult(
			baseResult,
			BodyTranslationProviderErrorKindProviderFailure,
			redactedBodyTranslationProviderFailureReason(err),
			metadata.FailureRetryable(),
		)
	}
}

func redactedBodyTranslationProviderFailureReason(err error) string {
	lowerMessage := strings.ToLower(strings.TrimSpace(err.Error()))
	if strings.Contains(lowerMessage, "api key") ||
		strings.Contains(lowerMessage, "authorization") ||
		strings.Contains(lowerMessage, "credential") ||
		strings.Contains(lowerMessage, "secret") {
		return "provider credential is unavailable"
	}
	return "provider request failed"
}

func bodyTranslationProviderFailureResult(
	baseResult BodyTranslationProviderResult,
	kind BodyTranslationProviderErrorKind,
	reason string,
	retryable bool,
) BodyTranslationProviderResult {
	baseResult.TranslatedCandidate = nil
	baseResult.ProtectionValidationTarget = nil
	baseResult.Failure = &BodyTranslationProviderFailure{
		Kind:       kind,
		Reason:     strings.TrimSpace(reason),
		Retryable:  retryable,
		IsRedacted: true,
	}
	return baseResult
}

func firstNonEmptyBodyTranslationValue(values ...string) string {
	for _, value := range values {
		trimmed := strings.TrimSpace(value)
		if trimmed != "" {
			return trimmed
		}
	}
	return ""
}
