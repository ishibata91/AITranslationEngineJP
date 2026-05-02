package ai

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"
)

const termTranslationExecutionModeSingleRequest = "single_request"

type providerTermTranslationItem struct {
	SourceTerm     string `json:"source_term"`
	TranslatedTerm string `json:"translated_term"`
}

// TermTranslationResponse defines one parsed provider response with correlated items.
type TermTranslationResponse struct {
	Items         []TermTranslationItem
	ExecutionMode string
}

// TermTranslationItem defines one correlated source/translated term pair.
type TermTranslationItem struct {
	SourceTerm     string
	TranslatedTerm string
}

// TermTranslationErrorKind identifies infra-level provider response failure families.
type TermTranslationErrorKind string

const (
	// TermTranslationErrorKindProviderFailure identifies transport or provider execution failures.
	TermTranslationErrorKindProviderFailure TermTranslationErrorKind = "provider_failure"
	// TermTranslationErrorKindInvalidProviderResponse identifies invalid response shapes or empty terms.
	TermTranslationErrorKindInvalidProviderResponse TermTranslationErrorKind = "invalid_provider_response"
)

// TermTranslationError wraps one typed provider term-translation failure.
type TermTranslationError struct {
	Kind      TermTranslationErrorKind
	Retryable bool
	cause     error
}

func (err *TermTranslationError) Error() string {
	if err == nil || err.cause == nil {
		return ""
	}
	return err.cause.Error()
}

// Unwrap returns the underlying transport or parsing failure.
func (err *TermTranslationError) Unwrap() error {
	if err == nil {
		return nil
	}
	return err.cause
}

// FailureKind exposes the typed failure kind without importing this package-specific enum.
func (err *TermTranslationError) FailureKind() string {
	if err == nil {
		return ""
	}
	return string(err.Kind)
}

// FailureRetryable exposes whether the provider failure is retryable.
func (err *TermTranslationError) FailureRetryable() bool {
	return err != nil && err.Retryable
}

func newTermTranslationError(kind TermTranslationErrorKind, retryable bool, cause error) error {
	return &TermTranslationError{
		Kind:      kind,
		Retryable: retryable,
		cause:     cause,
	}
}

func parseTermTranslationResponse(responseText string) (TermTranslationResponse, error) {
	trimmed := strings.TrimSpace(responseText)
	if trimmed == "" {
		return TermTranslationResponse{}, newTermTranslationError(
			TermTranslationErrorKindInvalidProviderResponse,
			true,
			errors.New(providerResponseEmptyError),
		)
	}

	if items, ok, err := parseTermTranslationObjectResponse(trimmed); ok || err != nil {
		if err != nil {
			return TermTranslationResponse{}, err
		}
		return termTranslationResponseFromItems(items), nil
	}
	if items, ok, err := parseTermTranslationArrayResponse(trimmed); ok || err != nil {
		if err != nil {
			return TermTranslationResponse{}, err
		}
		return termTranslationResponseFromItems(items), nil
	}

	return TermTranslationResponse{}, newTermTranslationError(
		TermTranslationErrorKindInvalidProviderResponse,
		true,
		fmt.Errorf("provider term translation response must be strict JSON with translations"),
	)
}

func parseTermTranslationObjectResponse(trimmed string) ([]TermTranslationItem, bool, error) {
	objectResponse := struct {
		Translations []providerTermTranslationItem `json:"translations"`
	}{}
	if err := json.Unmarshal([]byte(trimmed), &objectResponse); err == nil && len(objectResponse.Translations) > 0 {
		items, validationErr := normalizeTermTranslationItems(objectResponse.Translations)
		return items, true, validationErr
	}
	return nil, false, nil
}

func parseTermTranslationArrayResponse(trimmed string) ([]TermTranslationItem, bool, error) {
	arrayResponse := []providerTermTranslationItem{}
	if err := json.Unmarshal([]byte(trimmed), &arrayResponse); err == nil && len(arrayResponse) > 0 {
		items, validationErr := normalizeTermTranslationItems(arrayResponse)
		return items, true, validationErr
	}
	return nil, false, nil
}

func normalizeTermTranslationItems(rawItems []providerTermTranslationItem) ([]TermTranslationItem, error) {
	items := make([]TermTranslationItem, 0, len(rawItems))
	for _, item := range rawItems {
		sourceTerm := strings.TrimSpace(item.SourceTerm)
		translatedTerm := strings.TrimSpace(item.TranslatedTerm)
		if sourceTerm == "" || translatedTerm == "" {
			return nil, newTermTranslationError(
				TermTranslationErrorKindInvalidProviderResponse,
				true,
				fmt.Errorf("provider term translation item must include source_term and translated_term"),
			)
		}
		items = append(items, TermTranslationItem{
			SourceTerm:     sourceTerm,
			TranslatedTerm: translatedTerm,
		})
	}
	return items, nil
}

func termTranslationResponseFromItems(items []TermTranslationItem) TermTranslationResponse {
	return TermTranslationResponse{
		Items:         items,
		ExecutionMode: termTranslationExecutionModeSingleRequest,
	}
}

func isProviderExecutionRetryable(err error) bool {
	lowerMessage := strings.ToLower(strings.TrimSpace(err.Error()))
	if strings.Contains(lowerMessage, "unsupported ai provider") {
		return false
	}
	if strings.Contains(lowerMessage, "model is required") {
		return false
	}
	if strings.Contains(lowerMessage, "api key is required") {
		return false
	}
	if strings.Contains(lowerMessage, "provider base url") {
		return false
	}
	return true
}
