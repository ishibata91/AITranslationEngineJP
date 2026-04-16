package ai

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"strings"
	"testing"
)

func assertGeminiRequest(t *testing.T, request *http.Request) {
	t.Helper()
	if request == nil {
		t.Fatalf("expected provider request capture")
	}
	if request.Method != http.MethodPost {
		t.Fatalf("expected POST request, got %s", request.Method)
	}
	if request.Header.Get("Content-Type") != "application/json" {
		t.Fatalf("expected json content type")
	}
	if !strings.Contains(request.URL.String(), "key=k") {
		t.Fatalf("expected api key query param in url: %s", request.URL.String())
	}
	bodyBytes, err := io.ReadAll(request.Body)
	if err != nil {
		t.Fatalf("expected request body read to succeed: %v", err)
	}
	if !strings.Contains(string(bodyBytes), "prompt body") {
		t.Fatalf("expected prompt in request body, got %s", string(bodyBytes))
	}
}

func assertOpenAIRequest(
	t *testing.T,
	request *http.Request,
	expectedURL string,
	expectedAPIKey string,
	expectedModel string,
	expectedPrompt string,
	expectAuthorization bool,
) {
	t.Helper()
	if request == nil {
		t.Fatalf("expected provider request capture")
	}
	if request.Method != http.MethodPost {
		t.Fatalf("expected POST request, got %s", request.Method)
	}
	if request.URL.String() != expectedURL {
		t.Fatalf("expected openai-compatible endpoint %q, got %q", expectedURL, request.URL.String())
	}
	if request.Header.Get("Content-Type") != "application/json" {
		t.Fatalf("expected json content type")
	}
	authorization := request.Header.Get("Authorization")
	if expectAuthorization {
		if authorization != "Bearer "+expectedAPIKey {
			t.Fatalf("expected authorization header with bearer api key, got %q", authorization)
		}
	} else if authorization != "" {
		t.Fatalf("expected no authorization header, got %q", authorization)
	}
	bodyBytes, err := io.ReadAll(request.Body)
	if err != nil {
		t.Fatalf("expected request body read to succeed: %v", err)
	}
	parsed := struct {
		Model    string `json:"model"`
		Messages []struct {
			Role    string `json:"role"`
			Content string `json:"content"`
		} `json:"messages"`
	}{}
	if err := json.Unmarshal(bodyBytes, &parsed); err != nil {
		t.Fatalf("expected request body json to parse: %v", err)
	}
	if parsed.Model != expectedModel {
		t.Fatalf("expected model %q, got %q", expectedModel, parsed.Model)
	}
	if len(parsed.Messages) != 1 {
		t.Fatalf("expected one prompt message, got %#v", parsed.Messages)
	}
	if parsed.Messages[0].Role != "user" || parsed.Messages[0].Content != expectedPrompt {
		t.Fatalf("unexpected prompt payload: %#v", parsed.Messages[0])
	}
}

type stubHTTPTransport struct {
	doFunc      func(req *http.Request) (*http.Response, error)
	lastRequest *http.Request
}

func (transport *stubHTTPTransport) Do(req *http.Request) (*http.Response, error) {
	transport.lastRequest = req
	if transport.doFunc == nil {
		return newHTTPJSONResponse(http.StatusOK, `{"candidates":[{"content":{"parts":[{"text":"ok"}]}}]}`), nil
	}
	return transport.doFunc(req)
}

func newHTTPJSONResponse(statusCode int, body string) *http.Response {
	return &http.Response{
		StatusCode: statusCode,
		Body:       io.NopCloser(bytes.NewBufferString(body)),
		Header:     make(http.Header),
	}
}
