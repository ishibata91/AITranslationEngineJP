package ai

import (
	"context"
	"net/http"
	"strings"
	"testing"
)

func TestListProviderModelsOpenAICompatibleUsesModelsEndpoint(t *testing.T) {
	transport := &stubHTTPTransport{doFunc: func(_ *http.Request) (*http.Response, error) {
		return newHTTPJSONResponse(http.StatusOK, `{"data":[{"id":"gpt-5.4-mini"},{"id":"gpt-5.4"}]}`), nil
	}}

	models, err := ListProviderModels(context.Background(), transport, ListProviderModelsRequest{
		ProviderID:    ProviderXAI,
		APIKey:        "x-key",
		XAIBaseURL:    "https://api.x.ai/v1/",
		OpenAIBaseURL: "https://api.openai.com/v1",
	})
	if err != nil {
		t.Fatalf("expected model list success: %v", err)
	}
	if transport.lastRequest == nil || transport.lastRequest.Method != http.MethodGet {
		t.Fatalf("expected GET model list request, got %#v", transport.lastRequest)
	}
	if got := transport.lastRequest.URL.String(); got != "https://api.x.ai/v1/models" {
		t.Fatalf("expected xai model endpoint, got %q", got)
	}
	if got := transport.lastRequest.Header.Get("Authorization"); got != "Bearer x-key" {
		t.Fatalf("expected xai bearer auth, got %q", got)
	}
	if len(models) != 2 || models[0].ModelID != "gpt-5.4-mini" || models[1].ModelID != "gpt-5.4" {
		t.Fatalf("expected parsed model list, got %#v", models)
	}
}

func TestListProviderModelsGeminiUsesHeaderWithoutQuerySecret(t *testing.T) {
	transport := &stubHTTPTransport{doFunc: func(_ *http.Request) (*http.Response, error) {
		return newHTTPJSONResponse(http.StatusOK, `{"models":[{"name":"models/gemini-2.5-pro","displayName":"Gemini 2.5 Pro"}]}`), nil
	}}

	models, err := ListProviderModels(context.Background(), transport, ListProviderModelsRequest{
		ProviderID: ProviderGemini,
		APIKey:     "gemini-secret",
	})
	if err != nil {
		t.Fatalf("expected gemini model list success: %v", err)
	}
	if got := transport.lastRequest.URL.String(); strings.Contains(got, "gemini-secret") {
		t.Fatalf("expected gemini secret to stay out of url, got %q", got)
	}
	if got := transport.lastRequest.Header.Get("x-goog-api-key"); got != "gemini-secret" {
		t.Fatalf("expected x-goog-api-key header, got %q", got)
	}
	if len(models) != 1 || models[0].ModelID != "gemini-2.5-pro" {
		t.Fatalf("expected parsed gemini models, got %#v", models)
	}
}

func TestListProviderModelsLMStudioAllowsMissingAPIKey(t *testing.T) {
	transport := &stubHTTPTransport{doFunc: func(_ *http.Request) (*http.Response, error) {
		return newHTTPJSONResponse(http.StatusOK, `{"data":[{"id":"local-model"}]}`), nil
	}}

	models, err := ListProviderModels(context.Background(), transport, ListProviderModelsRequest{
		ProviderID:      ProviderLMStudio,
		LMStudioBaseURL: "http://localhost:1234/v1",
	})
	if err != nil {
		t.Fatalf("expected lm studio model list success: %v", err)
	}
	if got := transport.lastRequest.Header.Get("Authorization"); got != "" {
		t.Fatalf("expected no authorization header, got %q", got)
	}
	if len(models) != 1 || models[0].ModelID != "local-model" {
		t.Fatalf("expected parsed lm studio models, got %#v", models)
	}
}
