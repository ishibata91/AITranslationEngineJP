package ai

import (
	"context"
	"net/http"
	"strings"
	"testing"
)

func TestOpenAICompatibleProviderGenerateReturnsCommonResponse(t *testing.T) {
	provider := openAICompatibleProvider{
		transport: &stubHTTPTransport{doFunc: func(_ *http.Request) (*http.Response, error) {
			return newHTTPJSONResponse(http.StatusOK, `{"choices":[{"message":{"content":"openai text"}}]}`), nil
		}},
		baseURL:        "http://localhost:1234/v1",
		apiKeyOptional: true,
	}

	response, err := provider.Generate(context.Background(), ProviderRequest{Model: "local-model", Prompt: "prompt"})

	if err != nil {
		t.Fatalf("expected openai-compatible concrete generate to succeed: %v", err)
	}
	if response.Text != "openai text" {
		t.Fatalf("expected common response text from openai-compatible concrete, got %q", response.Text)
	}
}

func TestProviderClientGenerateTextRejectsOpenAICompatibleInvalidJSON(t *testing.T) {
	client := NewProviderClient(&stubHTTPTransport{doFunc: func(_ *http.Request) (*http.Response, error) {
		return newHTTPJSONResponse(http.StatusOK, `{"choices":`), nil
	}})

	_, err := client.GenerateText(context.Background(), ProviderXAI, ProviderRequest{Model: "grok-2", APIKey: "x-key", Prompt: "prompt"})

	if err == nil || !strings.Contains(err.Error(), "parse ai provider response") {
		t.Fatalf("expected parse error, got %v", err)
	}
}

func TestProviderClientGenerateTextRejectsOpenAICompatibleProviderErrorPayload(t *testing.T) {
	client := NewProviderClient(&stubHTTPTransport{doFunc: func(_ *http.Request) (*http.Response, error) {
		return newHTTPJSONResponse(http.StatusOK, `{"error":{"message":"quota exhausted"}}`), nil
	}})

	_, err := client.GenerateText(context.Background(), ProviderLMStudio, ProviderRequest{Model: "local-model", Prompt: "prompt"})

	if err == nil || !strings.Contains(err.Error(), "quota exhausted") {
		t.Fatalf("expected provider error payload, got %v", err)
	}
}

func TestProviderClientGenerateTextRejectsOpenAICompatibleNon2xxStatus(t *testing.T) {
	client := NewProviderClient(&stubHTTPTransport{doFunc: func(_ *http.Request) (*http.Response, error) {
		return newHTTPJSONResponse(http.StatusInternalServerError, `{"choices":[{"message":{"content":"fallback"}}]}`), nil
	}})

	_, err := client.GenerateText(context.Background(), ProviderXAI, ProviderRequest{Model: "grok-2", APIKey: "x-key", Prompt: "prompt"})

	if err == nil || !strings.Contains(err.Error(), "status=500") {
		t.Fatalf("expected non-2xx status error, got %v", err)
	}
}

func TestProviderClientGenerateTextRejectsOpenAICompatibleEmptyResponse(t *testing.T) {
	client := NewProviderClient(&stubHTTPTransport{doFunc: func(_ *http.Request) (*http.Response, error) {
		return newHTTPJSONResponse(http.StatusOK, `{"choices":[{"message":{"content":"   "}}]}`), nil
	}})

	_, err := client.GenerateText(context.Background(), ProviderLMStudio, ProviderRequest{Model: "local-model", Prompt: "prompt"})

	if err == nil || !strings.Contains(err.Error(), "response is empty") {
		t.Fatalf("expected empty response error, got %v", err)
	}
}

func TestExtractOpenAICompatibleTextReturnsFirstNonEmptyText(t *testing.T) {
	response := openAICompatibleChatResponse{
		Choices: []struct {
			Message struct {
				Content string `json:"content"`
			} `json:"message"`
		}{
			{Message: struct {
				Content string `json:"content"`
			}{Content: " "}},
			{Message: struct {
				Content string `json:"content"`
			}{Content: "  openai text  "}},
		},
	}

	if text := extractOpenAICompatibleText(response); text != "openai text" {
		t.Fatalf("expected first non-empty openai message content, got %q", text)
	}
}
