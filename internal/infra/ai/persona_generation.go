package ai

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
)

const personaGenerationExecutionModeSingleRequest = "single_request"

type providerPersonaGenerationItem struct {
	RequestUnitID    string `json:"request_unit_id"`
	NPCCorrelationID string `json:"npc_correlation_id"`
	PersonaBody      string `json:"persona_body"`
}

// PersonaGenerationResponse defines one parsed provider response with correlated NPC persona items.
type PersonaGenerationResponse struct {
	Items         []PersonaGenerationItem
	ExecutionMode string
	PromptDigest  string
	DebugLog      ProviderDebugLog
}

// PersonaGenerationItem defines one correlated persona-generation result item.
type PersonaGenerationItem struct {
	RequestUnitID    string
	NPCCorrelationID string
	PersonaBody      string
}

// PersonaGenerationErrorKind identifies infra-level provider response failure families.
type PersonaGenerationErrorKind string

const (
	// PersonaGenerationErrorKindProviderFailure identifies transport or provider execution failures.
	PersonaGenerationErrorKindProviderFailure PersonaGenerationErrorKind = "provider_failure"
	// PersonaGenerationErrorKindInvalidProviderResponse identifies invalid response shapes or empty persona data.
	PersonaGenerationErrorKindInvalidProviderResponse PersonaGenerationErrorKind = "invalid_provider_response"
)

// PersonaGenerationError wraps one typed provider persona-generation failure.
type PersonaGenerationError struct {
	Kind      PersonaGenerationErrorKind
	Retryable bool
	cause     error
}

func (err *PersonaGenerationError) Error() string {
	if err == nil || err.cause == nil {
		return ""
	}
	return err.cause.Error()
}

// Unwrap returns the underlying transport or parsing failure.
func (err *PersonaGenerationError) Unwrap() error {
	if err == nil {
		return nil
	}
	return err.cause
}

// FailureKind exposes the typed failure kind without importing this package-specific enum.
func (err *PersonaGenerationError) FailureKind() string {
	if err == nil {
		return ""
	}
	return string(err.Kind)
}

// FailureRetryable exposes whether the provider failure is retryable.
func (err *PersonaGenerationError) FailureRetryable() bool {
	return err != nil && err.Retryable
}

func newPersonaGenerationError(kind PersonaGenerationErrorKind, retryable bool, cause error) error {
	return &PersonaGenerationError{
		Kind:      kind,
		Retryable: retryable,
		cause:     cause,
	}
}

func parsePersonaGenerationResponse(responseText string) (PersonaGenerationResponse, error) {
	trimmed := strings.TrimSpace(responseText)
	if trimmed == "" {
		return PersonaGenerationResponse{}, newPersonaGenerationError(
			PersonaGenerationErrorKindInvalidProviderResponse,
			true,
			errors.New(providerResponseEmptyError),
		)
	}

	if items, ok, err := parsePersonaGenerationObjectResponse(trimmed); ok || err != nil {
		if err != nil {
			return PersonaGenerationResponse{}, err
		}
		return personaGenerationResponseFromItems(items), nil
	}
	if items, ok, err := parsePersonaGenerationArrayResponse(trimmed); ok || err != nil {
		if err != nil {
			return PersonaGenerationResponse{}, err
		}
		return personaGenerationResponseFromItems(items), nil
	}

	return PersonaGenerationResponse{}, newPersonaGenerationError(
		PersonaGenerationErrorKindInvalidProviderResponse,
		true,
		fmt.Errorf("provider persona generation response must be strict JSON with personas"),
	)
}

func parsePersonaGenerationObjectResponse(trimmed string) ([]PersonaGenerationItem, bool, error) {
	objectResponse := struct {
		Personas []providerPersonaGenerationItem `json:"personas"`
	}{}
	if err := json.Unmarshal([]byte(trimmed), &objectResponse); err == nil && len(objectResponse.Personas) > 0 {
		items, validationErr := normalizePersonaGenerationItems(objectResponse.Personas)
		return items, true, validationErr
	}
	return nil, false, nil
}

func parsePersonaGenerationArrayResponse(trimmed string) ([]PersonaGenerationItem, bool, error) {
	arrayResponse := []providerPersonaGenerationItem{}
	if err := json.Unmarshal([]byte(trimmed), &arrayResponse); err == nil && len(arrayResponse) > 0 {
		items, validationErr := normalizePersonaGenerationItems(arrayResponse)
		return items, true, validationErr
	}
	return nil, false, nil
}

func normalizePersonaGenerationItems(rawItems []providerPersonaGenerationItem) ([]PersonaGenerationItem, error) {
	items := make([]PersonaGenerationItem, 0, len(rawItems))
	for _, item := range rawItems {
		requestUnitID := strings.TrimSpace(item.RequestUnitID)
		npcCorrelationID := strings.TrimSpace(item.NPCCorrelationID)
		personaBody := strings.TrimSpace(item.PersonaBody)
		if requestUnitID == "" || npcCorrelationID == "" || personaBody == "" {
			return nil, newPersonaGenerationError(
				PersonaGenerationErrorKindInvalidProviderResponse,
				true,
				fmt.Errorf("provider persona item must include request_unit_id, npc_correlation_id, and persona_body"),
			)
		}
		items = append(items, PersonaGenerationItem{
			RequestUnitID:    requestUnitID,
			NPCCorrelationID: npcCorrelationID,
			PersonaBody:      personaBody,
		})
	}
	return items, nil
}

func personaGenerationResponseFromItems(items []PersonaGenerationItem) PersonaGenerationResponse {
	return PersonaGenerationResponse{
		Items:         items,
		ExecutionMode: personaGenerationExecutionModeSingleRequest,
	}
}

func promptDigestSHA256(prompt string) string {
	digest := sha256.Sum256([]byte(prompt))
	return "sha256:" + hex.EncodeToString(digest[:])
}
