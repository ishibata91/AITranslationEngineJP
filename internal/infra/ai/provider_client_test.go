package ai

import (
	"context"
	"errors"
	"net/http"
	"strings"
	"testing"
)

func TestProviderClientReportsTestSafety(t *testing.T) {
	testSafeClient := NewProviderClient(NewTestSafeHTTPTransport())

	if !testSafeClient.ProviderRequestsAreTestSafe() {
		t.Fatalf("expected deterministic transport client to be test-safe")
	}
}

func TestProviderClientReportsRealTransportAsNotTestSafe(t *testing.T) {
	realClient := NewProviderClient(&stubHTTPTransport{})

	if realClient.ProviderRequestsAreTestSafe() {
		t.Fatalf("expected regular transport client to be non-test-safe")
	}
}

func TestProviderClientGenerateTextRejectsUnsupportedProvider(t *testing.T) {
	client := NewProviderClient(&stubHTTPTransport{})

	_, err := client.GenerateText(context.Background(), "unknown", ProviderRequest{Model: "gemini-2.5-pro", APIKey: "k", Prompt: "prompt"})

	if err == nil || !strings.Contains(err.Error(), "unsupported ai provider") {
		t.Fatalf("expected unsupported provider error, got %v", err)
	}
}

func TestProviderClientGenerateTextRejectsMissingGeminiModel(t *testing.T) {
	client := NewProviderClient(&stubHTTPTransport{})

	_, err := client.GenerateText(context.Background(), ProviderGemini, ProviderRequest{Model: "", APIKey: "k", Prompt: "prompt"})

	if err == nil || !strings.Contains(err.Error(), "model is required") {
		t.Fatalf("expected model validation error, got %v", err)
	}
}

func TestProviderClientGenerateTextRejectsMissingXAIAPIKeyInRealMode(t *testing.T) {
	client := NewProviderClient(&stubHTTPTransport{})

	_, err := client.GenerateText(context.Background(), ProviderXAI, ProviderRequest{Model: "grok-2", Prompt: "prompt"})

	if err == nil || !strings.Contains(err.Error(), "api key is required") {
		t.Fatalf("expected api key required error, got %v", err)
	}
}

func TestProviderClientGenerateTextWrapsTransportError(t *testing.T) {
	transportErr := errors.New("network failed")
	client := NewProviderClient(&stubHTTPTransport{doFunc: func(_ *http.Request) (*http.Response, error) {
		return nil, transportErr
	}})

	_, err := client.GenerateText(context.Background(), ProviderGemini, ProviderRequest{Model: "gemini-2.5-pro", APIKey: "k", Prompt: "prompt"})

	if !errors.Is(err, transportErr) {
		t.Fatalf("expected transport error wrapping, got %v", err)
	}
}

func TestProviderClientGenerateTextSendsXAIRequest(t *testing.T) {
	transport := &stubHTTPTransport{doFunc: func(_ *http.Request) (*http.Response, error) {
		return newHTTPJSONResponse(http.StatusOK, `{"choices":[{"message":{"content":"  ai success  "}}]}`), nil
	}}
	client := NewProviderClient(transport)

	response, err := client.GenerateText(context.Background(), ProviderXAI, ProviderRequest{Model: "grok-2", APIKey: "x-key", Prompt: "prompt body"})

	if err != nil {
		t.Fatalf("expected success, got %v", err)
	}
	if response.Text != "ai success" {
		t.Fatalf("expected trimmed response text, got %q", response.Text)
	}
	assertOpenAIRequest(t, transport.lastRequest, "https://api.x.ai/v1/chat/completions", "x-key", "grok-2", "prompt body", true)
}

func TestProviderClientGenerateTextAppliesXAIBaseURLOverride(t *testing.T) {
	transport := &stubHTTPTransport{doFunc: func(_ *http.Request) (*http.Response, error) {
		return newHTTPJSONResponse(http.StatusOK, `{"choices":[{"message":{"content":"ok"}}]}`), nil
	}}
	client := NewProviderClient(transport, WithXAIBaseURL("https://gateway.example.com/custom/v1/"))

	_, err := client.GenerateText(context.Background(), ProviderXAI, ProviderRequest{Model: "grok-2", APIKey: "x-key", Prompt: "prompt"})

	if err != nil {
		t.Fatalf("expected success, got %v", err)
	}
	assertOpenAIRequest(t, transport.lastRequest, "https://gateway.example.com/custom/v1/chat/completions", "x-key", "grok-2", "prompt", true)
}

func TestProviderClientGenerateTextSendsLMStudioRequestWithoutAuthorization(t *testing.T) {
	transport := &stubHTTPTransport{doFunc: func(_ *http.Request) (*http.Response, error) {
		return newHTTPJSONResponse(http.StatusOK, `{"choices":[{"message":{"content":"lmstudio body"}}]}`), nil
	}}
	client := NewProviderClient(transport)

	response, err := client.GenerateText(context.Background(), ProviderLMStudio, ProviderRequest{Model: "local-model", Prompt: "prompt"})

	if err != nil {
		t.Fatalf("expected success, got %v", err)
	}
	if response.Text != "lmstudio body" {
		t.Fatalf("unexpected text: %q", response.Text)
	}
	assertOpenAIRequest(t, transport.lastRequest, "http://localhost:1234/v1/chat/completions", "", "local-model", "prompt", false)
}

func TestProviderClientGenerateTextAppliesLMStudioAuthorizationAndBaseURLOverride(t *testing.T) {
	transport := &stubHTTPTransport{doFunc: func(_ *http.Request) (*http.Response, error) {
		return newHTTPJSONResponse(http.StatusOK, `{"choices":[{"message":{"content":"lmstudio body"}}]}`), nil
	}}
	client := NewProviderClient(transport, WithLMStudioBaseURL("http://127.0.0.1:1234/proxy/v1"))

	_, err := client.GenerateText(context.Background(), ProviderLMStudio, ProviderRequest{Model: "local-model", APIKey: "optional-key", Prompt: "prompt"})

	if err != nil {
		t.Fatalf("expected success, got %v", err)
	}
	assertOpenAIRequest(t, transport.lastRequest, "http://127.0.0.1:1234/proxy/v1/chat/completions", "optional-key", "local-model", "prompt", true)
}

func TestProviderClientGenerateTextUsesProviderResponse(t *testing.T) {
	stub := &stubProvider{response: ProviderResponse{Text: "registry text"}}
	client := &ProviderClient{
		transport: &stubHTTPTransport{},
		providers: map[string]provider{ProviderGemini: stub},
	}

	response, err := client.GenerateText(context.Background(), ProviderGemini, ProviderRequest{Model: "model", APIKey: "key", Prompt: "prompt"})

	if err != nil {
		t.Fatalf("expected provider registry generate to succeed: %v", err)
	}
	if response.Text != "registry text" {
		t.Fatalf("expected text from common provider response, got %q", response.Text)
	}
	if stub.calls != 1 {
		t.Fatalf("expected provider generate to be called once, got %d", stub.calls)
	}
	if stub.request.Model != "model" || stub.request.APIKey != "key" || stub.request.Prompt != "prompt" {
		t.Fatalf("expected request to be bridged into common provider request, got %#v", stub.request)
	}
}

type stubProvider struct {
	response ProviderResponse
	err      error
	request  ProviderRequest
	calls    int
}

func (provider *stubProvider) Generate(
	_ context.Context,
	request ProviderRequest,
) (ProviderResponse, error) {
	provider.calls++
	provider.request = request
	if provider.err != nil {
		return ProviderResponse{}, provider.err
	}
	return provider.response, nil
}
