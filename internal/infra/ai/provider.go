package ai

import (
	"context"
	"net/http"
)

const (
	// ProviderGemini defines the supported Gemini provider id.
	ProviderGemini = "gemini"
	// ProviderLMStudio defines the supported LM Studio provider id.
	ProviderLMStudio = "lm_studio"
	// ProviderXAI defines the supported xAI provider id.
	ProviderXAI = "xai"
	// ProviderFake defines the supported test-only fake provider id.
	ProviderFake = "fake"

	geminiEndpointTemplate = "https://generativelanguage.googleapis.com/v1beta/models/%s:generateContent"
	lmStudioDefaultBaseURL = "http://localhost:1234/v1"
	xaiDefaultBaseURL      = "https://api.x.ai/v1"
	defaultTestSafeText    = "テスト用の決定論的な AI 応答です。"

	providerResponseEmptyError = "ai provider response is empty"
)

// HTTPTransport defines the low-level HTTP request seam for AI provider calls.
type HTTPTransport interface {
	Do(req *http.Request) (*http.Response, error)
}

type testSafeHTTPTransport interface {
	HTTPTransport
	testSafeTransportMarker()
}

// ProviderRequest defines the provider-agnostic AI text request contract.
type ProviderRequest struct {
	Model  string
	APIKey string
	Prompt string
}

// ProviderDebugLog exposes redacted prompt/request diagnostics for debugging.
type ProviderDebugLog struct {
	Prompt         string
	RequestBody    string
	Headers        map[string]string
	SecretRedacted bool
}

// ProviderResponse defines the provider-agnostic AI text response contract.
type ProviderResponse struct {
	Text     string
	DebugLog ProviderDebugLog
}

type provider interface {
	Generate(ctx context.Context, request ProviderRequest) (ProviderResponse, error)
}
