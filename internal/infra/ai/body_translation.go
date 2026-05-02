package ai

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"
)

const bodyTranslationExecutionModeSingleRequest = "single_request"

type providerBodyTranslationItem struct {
	RequestUnitID       string `json:"request_unit_id"`
	FieldCorrelationKey string `json:"field_correlation_key"`
	TranslatedText      string `json:"translated_text"`
}

// BodyTranslationResponse defines one parsed provider response with correlated field items.
type BodyTranslationResponse struct {
	Items         []BodyTranslationItem
	ExecutionMode string
	PromptDigest  string
}

// BodyTranslationItem defines one correlated body translation result item.
type BodyTranslationItem struct {
	RequestUnitID       string
	FieldCorrelationKey string
	TranslatedText      string
}

// BodyTranslationErrorKind identifies infra-level provider response failure families.
type BodyTranslationErrorKind string

const (
	// BodyTranslationErrorKindProviderFailure identifies transport or provider execution failures.
	BodyTranslationErrorKindProviderFailure BodyTranslationErrorKind = "provider_failure"
	// BodyTranslationErrorKindInvalidProviderResponse identifies invalid response shapes or empty translations.
	BodyTranslationErrorKindInvalidProviderResponse BodyTranslationErrorKind = "invalid_provider_response"
)

// BodyTranslationError wraps one typed provider body-translation failure.
type BodyTranslationError struct {
	Kind      BodyTranslationErrorKind
	Retryable bool
	cause     error
}

func (err *BodyTranslationError) Error() string {
	if err == nil || err.cause == nil {
		return ""
	}
	return err.cause.Error()
}

// Unwrap returns the underlying transport or parsing failure.
func (err *BodyTranslationError) Unwrap() error {
	if err == nil {
		return nil
	}
	return err.cause
}

// FailureKind exposes the typed failure kind without importing this package-specific enum.
func (err *BodyTranslationError) FailureKind() string {
	if err == nil {
		return ""
	}
	return string(err.Kind)
}

// FailureRetryable exposes whether the provider failure is retryable.
func (err *BodyTranslationError) FailureRetryable() bool {
	return err != nil && err.Retryable
}

func newBodyTranslationError(kind BodyTranslationErrorKind, retryable bool, cause error) error {
	return &BodyTranslationError{
		Kind:      kind,
		Retryable: retryable,
		cause:     cause,
	}
}

func parseBodyTranslationResponse(responseText string) (BodyTranslationResponse, error) {
	trimmed := strings.TrimSpace(responseText)
	if trimmed == "" {
		return BodyTranslationResponse{}, newBodyTranslationError(
			BodyTranslationErrorKindInvalidProviderResponse,
			true,
			errors.New(providerResponseEmptyError),
		)
	}

	if items, ok, err := parseBodyTranslationObjectResponse(trimmed); ok || err != nil {
		if err != nil {
			return BodyTranslationResponse{}, err
		}
		return bodyTranslationResponseFromItems(items), nil
	}
	if items, ok, err := parseBodyTranslationArrayResponse(trimmed); ok || err != nil {
		if err != nil {
			return BodyTranslationResponse{}, err
		}
		return bodyTranslationResponseFromItems(items), nil
	}

	return BodyTranslationResponse{}, newBodyTranslationError(
		BodyTranslationErrorKindInvalidProviderResponse,
		true,
		fmt.Errorf("provider body translation response must be strict JSON with translations"),
	)
}

func parseBodyTranslationObjectResponse(trimmed string) ([]BodyTranslationItem, bool, error) {
	objectResponse := struct {
		Translations []providerBodyTranslationItem `json:"translations"`
	}{}
	if err := json.Unmarshal([]byte(trimmed), &objectResponse); err == nil && len(objectResponse.Translations) > 0 {
		items, validationErr := normalizeBodyTranslationItems(objectResponse.Translations)
		return items, true, validationErr
	}
	return nil, false, nil
}

func parseBodyTranslationArrayResponse(trimmed string) ([]BodyTranslationItem, bool, error) {
	arrayResponse := []providerBodyTranslationItem{}
	if err := json.Unmarshal([]byte(trimmed), &arrayResponse); err == nil && len(arrayResponse) > 0 {
		items, validationErr := normalizeBodyTranslationItems(arrayResponse)
		return items, true, validationErr
	}
	return nil, false, nil
}

func normalizeBodyTranslationItems(rawItems []providerBodyTranslationItem) ([]BodyTranslationItem, error) {
	items := make([]BodyTranslationItem, 0, len(rawItems))
	for _, item := range rawItems {
		requestUnitID := strings.TrimSpace(item.RequestUnitID)
		fieldCorrelationKey := strings.TrimSpace(item.FieldCorrelationKey)
		translatedText := strings.TrimSpace(item.TranslatedText)
		if requestUnitID == "" || fieldCorrelationKey == "" || translatedText == "" {
			return nil, newBodyTranslationError(
				BodyTranslationErrorKindInvalidProviderResponse,
				true,
				fmt.Errorf("provider body translation item must include request_unit_id, field_correlation_key, and translated_text"),
			)
		}
		items = append(items, BodyTranslationItem{
			RequestUnitID:       requestUnitID,
			FieldCorrelationKey: fieldCorrelationKey,
			TranslatedText:      translatedText,
		})
	}
	return items, nil
}

func bodyTranslationResponseFromItems(items []BodyTranslationItem) BodyTranslationResponse {
	return BodyTranslationResponse{
		Items:         items,
		ExecutionMode: bodyTranslationExecutionModeSingleRequest,
	}
}
