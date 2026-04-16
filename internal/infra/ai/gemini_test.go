package ai

import (
	"context"
	"net/http"
	"testing"
)

func TestGeminiProviderGenerateReturnsCommonResponse(t *testing.T) {
	provider := geminiProvider{transport: &stubHTTPTransport{doFunc: func(_ *http.Request) (*http.Response, error) {
		return newHTTPJSONResponse(http.StatusOK, `{"candidates":[{"content":{"parts":[{"text":"gemini text"}]}}]}`), nil
	}}}

	response, err := provider.Generate(context.Background(), ProviderRequest{Model: "gemini-2.5-pro", APIKey: "k", Prompt: "prompt"})

	if err != nil {
		t.Fatalf("expected gemini concrete generate to succeed: %v", err)
	}
	if response.Text != "gemini text" {
		t.Fatalf("expected common response text from gemini concrete, got %q", response.Text)
	}
}

func TestProviderClientGenerateTextSendsGeminiRequest(t *testing.T) {
	transport := &stubHTTPTransport{doFunc: func(_ *http.Request) (*http.Response, error) {
		return newHTTPJSONResponse(http.StatusOK, `{"candidates":[{"content":{"parts":[{"text":"  ai success  "}]}}]}`), nil
	}}
	client := NewProviderClient(transport)

	response, err := client.GenerateText(context.Background(), ProviderGemini, ProviderRequest{Model: "gemini-2.5-pro", APIKey: "k", Prompt: "prompt body"})

	if err != nil {
		t.Fatalf("expected success, got %v", err)
	}
	if response.Text != "ai success" {
		t.Fatalf("expected trimmed response text, got %q", response.Text)
	}
	assertGeminiRequest(t, transport.lastRequest)
}

func TestExtractGeminiTextReturnsFirstNonEmptyText(t *testing.T) {
	text := extractGeminiText(geminiGenerateResponse{
		Candidates: []struct {
			Content struct {
				Parts []geminiPart `json:"parts"`
			} `json:"content"`
		}{
			{Content: struct {
				Parts []geminiPart `json:"parts"`
			}{Parts: []geminiPart{{Text: " "}, {Text: "  first text  "}}}},
		},
	})

	if text != "first text" {
		t.Fatalf("expected first non-empty text, got %q", text)
	}
}

func TestExtractGeminiTextReturnsEmptyWhenNoTextExists(t *testing.T) {
	text := extractGeminiText(geminiGenerateResponse{
		Candidates: []struct {
			Content struct {
				Parts []geminiPart `json:"parts"`
			} `json:"content"`
		}{
			{Content: struct {
				Parts []geminiPart `json:"parts"`
			}{Parts: []geminiPart{{Text: " "}}}},
		},
	})

	if text != "" {
		t.Fatalf("expected empty text, got %q", text)
	}
}
