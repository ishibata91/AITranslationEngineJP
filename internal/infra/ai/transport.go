package ai

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
)

// NewTestSafeHTTPTransport creates a deterministic test-safe transport seam with the default response text.
func NewTestSafeHTTPTransport() HTTPTransport {
	return &deterministicHTTPTransport{}
}

// NewTestSafeHTTPTransportWithResponse creates a deterministic test-safe transport seam with the given response text.
// Passing an empty string yields an empty-text response, which causes provider parsers to return an empty-response error.
func NewTestSafeHTTPTransportWithResponse(responseText string) HTTPTransport {
	return &deterministicHTTPTransport{
		responseText:          strings.TrimSpace(responseText),
		useConfiguredResponse: true,
	}
}

type deterministicHTTPTransport struct {
	responseText          string
	useConfiguredResponse bool
}

func (transport *deterministicHTTPTransport) Do(request *http.Request) (*http.Response, error) {
	responseText, err := resolveDeterministicResponseText(request, *transport)
	if err != nil {
		return nil, err
	}
	payload := map[string]interface{}{
		"candidates": []map[string]interface{}{
			{
				"content": map[string]interface{}{
					"parts": []map[string]string{{"text": responseText}},
				},
			},
		},
		"choices": []map[string]interface{}{
			{
				"message": map[string]string{"content": responseText},
			},
		},
	}
	responseBytes, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("marshal deterministic ai provider response: %w", err)
	}
	return &http.Response{
		StatusCode: http.StatusOK,
		Body:       io.NopCloser(bytes.NewReader(responseBytes)),
		Header:     make(http.Header),
	}, nil
}

func (transport *deterministicHTTPTransport) testSafeTransportMarker() {
	// Marker method designates this transport as test-safe for provider request DI.
}

func resolveDeterministicResponseText(request *http.Request, transport deterministicHTTPTransport) (string, error) {
	if transport.useConfiguredResponse {
		return strings.TrimSpace(transport.responseText), nil
	}
	prompt, err := extractPromptFromProviderRequest(request)
	if err != nil {
		return "", err
	}
	if prompt == "" {
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

func extractPromptFromProviderRequest(request *http.Request) (string, error) {
	if request == nil || request.Body == nil {
		return "", nil
	}
	requestBody, err := io.ReadAll(request.Body)
	if err != nil {
		return "", fmt.Errorf("read deterministic ai provider request: %w", err)
	}
	return extractPromptFromProviderRequestBody(requestBody), nil
}

func extractPromptFromProviderRequestBody(requestBody []byte) string {
	if len(requestBody) == 0 {
		return ""
	}
	if prompt := extractOpenAIProviderPrompt(requestBody); prompt != "" {
		return prompt
	}
	return extractGeminiProviderPrompt(requestBody)
}

func extractOpenAIProviderPrompt(requestBody []byte) string {
	openAIRequest := struct {
		Messages []struct {
			Content string `json:"content"`
		} `json:"messages"`
	}{}
	if err := json.Unmarshal(requestBody, &openAIRequest); err == nil {
		for _, message := range openAIRequest.Messages {
			trimmed := strings.TrimSpace(message.Content)
			if trimmed != "" {
				return trimmed
			}
		}
	}
	return ""
}

func extractGeminiProviderPrompt(requestBody []byte) string {
	geminiRequest := struct {
		Contents []struct {
			Parts []struct {
				Text string `json:"text"`
			} `json:"parts"`
		} `json:"contents"`
	}{}
	if err := json.Unmarshal(requestBody, &geminiRequest); err == nil {
		for _, content := range geminiRequest.Contents {
			for _, part := range content.Parts {
				trimmed := strings.TrimSpace(part.Text)
				if trimmed != "" {
					return trimmed
				}
			}
		}
	}
	return ""
}

func extractPromptField(prompt string, key string) string {
	fieldPrefix := strings.TrimSpace(key) + "="
	for _, line := range strings.Split(prompt, "\n") {
		trimmedLine := strings.TrimSpace(line)
		if !strings.HasPrefix(trimmedLine, fieldPrefix) {
			continue
		}
		return strings.TrimSpace(strings.TrimPrefix(trimmedLine, fieldPrefix))
	}
	return ""
}

func buildDeterministicTermTranslationResponseText(sourceTerm string) (string, error) {
	responseBytes, err := json.Marshal(map[string]any{
		"translations": []map[string]string{
			{
				"source_term":     sourceTerm,
				"translated_term": sourceTerm + "-translated",
			},
		},
	})
	if err != nil {
		return "", fmt.Errorf("marshal deterministic term translation response: %w", err)
	}
	return string(responseBytes), nil
}

func buildDeterministicPersonaGenerationResponseText(requestUnitID string, npcCorrelationID string) (string, error) {
	responseBytes, err := json.Marshal(map[string]any{
		"personas": []map[string]string{
			{
				"request_unit_id":    requestUnitID,
				"npc_correlation_id": npcCorrelationID,
				"persona_body":       "決定論的なペルソナ応答",
			},
		},
	})
	if err != nil {
		return "", fmt.Errorf("marshal deterministic persona generation response: %w", err)
	}
	return string(responseBytes), nil
}

func buildDeterministicBodyTranslationResponseText(requestUnitID string, fieldCorrelationKey string) (string, error) {
	responseBytes, err := json.Marshal(map[string]any{
		"translations": []map[string]string{
			{
				"request_unit_id":       requestUnitID,
				"field_correlation_key": fieldCorrelationKey,
				"translated_text":       "決定論的な本文翻訳応答",
			},
		},
	})
	if err != nil {
		return "", fmt.Errorf("marshal deterministic body translation response: %w", err)
	}
	return string(responseBytes), nil
}

func buildDeterministicBodyTranslationResponseFromPrompt(prompt string) (string, error) {
	requestUnitID := extractPromptField(prompt, "request_unit_id")
	fieldCorrelationKey := extractPromptField(prompt, "field_correlation_key")
	if requestUnitID == "" || fieldCorrelationKey == "" {
		return "", fmt.Errorf("body translation prompt must include request_unit_id and field_correlation_key")
	}
	return buildDeterministicBodyTranslationResponseText(requestUnitID, fieldCorrelationKey)
}

func buildProviderDebugLog(
	prompt string,
	requestBytes []byte,
	headers http.Header,
) ProviderDebugLog {
	return ProviderDebugLog{
		Prompt:         strings.TrimSpace(prompt),
		RequestBody:    strings.TrimSpace(string(requestBytes)),
		Headers:        redactProviderHeaders(headers),
		SecretRedacted: true,
	}
}

func redactProviderHeaders(headers http.Header) map[string]string {
	redacted := make(map[string]string, len(headers))
	for key, values := range headers {
		normalizedKey := http.CanonicalHeaderKey(strings.TrimSpace(key))
		if normalizedKey == "" {
			continue
		}
		switch normalizedKey {
		case "Authorization", "X-Goog-Api-Key":
			redacted[normalizedKey] = "[REDACTED]"
		default:
			redacted[normalizedKey] = strings.Join(values, ",")
		}
	}
	return redacted
}

func callProviderTransport(
	transport HTTPTransport,
	request *http.Request,
) (*http.Response, error) {
	if transport == nil {
		return nil, fmt.Errorf("ai provider transport is required")
	}
	response, err := transport.Do(request)
	if err != nil {
		return nil, fmt.Errorf("call ai provider transport: %w", err)
	}
	if response == nil || response.Body == nil {
		return nil, errors.New(providerResponseEmptyError)
	}
	return response, nil
}
