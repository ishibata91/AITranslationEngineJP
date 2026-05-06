package ai

import (
	"context"
	"fmt"
	"strings"
)

// FakeModelID is the deterministic model id exposed by test-safe provider DI.
const FakeModelID = "fake-model"

type deterministicProvider struct {
	responseText          string
	useConfiguredResponse bool
}

func (provider deterministicProvider) Generate(
	_ context.Context,
	request ProviderRequest,
) (ProviderResponse, error) {
	requestBytes, err := newOpenAICompatibleRequestBytes(request.Model, request.Prompt)
	if err != nil {
		return ProviderResponse{}, err
	}
	text, err := provider.resolveResponseText(request.Prompt)
	if err != nil {
		return ProviderResponse{}, fmt.Errorf("build deterministic ai provider response: %w", err)
	}
	return ProviderResponse{
		Text: text,
		DebugLog: buildProviderDebugLog(
			request.Prompt,
			requestBytes,
			nil,
		),
	}, nil
}

func (provider deterministicProvider) ListModels(
	context.Context,
	string,
	string,
) ([]ProviderModelOption, error) {
	return []ProviderModelOption{{ModelID: FakeModelID, Label: FakeModelID}}, nil
}

func (provider deterministicProvider) resolveResponseText(prompt string) (string, error) {
	if provider.useConfiguredResponse {
		return strings.TrimSpace(provider.responseText), nil
	}
	if strings.TrimSpace(prompt) == "" {
		return defaultTestSafeText, nil
	}
	if strings.Contains(prompt, "BODY_TRANSLATION_REQUEST_V1") {
		requestUnitID := extractPromptField(prompt, "request_unit_id")
		fieldCorrelationKey := extractPromptField(prompt, "field_correlation_key")
		if requestUnitID == "" || fieldCorrelationKey == "" {
			return defaultTestSafeText, nil
		}
		return buildDeterministicBodyTranslationResponseText(requestUnitID, fieldCorrelationKey)
	}
	if strings.Contains(prompt, "PERSONA_GENERATION_REQUEST_V1") {
		requestUnitID := extractPromptField(prompt, "request_unit_id")
		npcCorrelationID := extractPromptField(prompt, "npc_correlation_id")
		if requestUnitID == "" || npcCorrelationID == "" {
			return defaultTestSafeText, nil
		}
		return buildDeterministicPersonaGenerationResponseText(requestUnitID, npcCorrelationID)
	}
	sourceTerm := extractPromptField(prompt, "source_term")
	if sourceTerm == "" {
		return defaultTestSafeText, nil
	}
	return buildDeterministicTermTranslationResponseText(sourceTerm)
}
