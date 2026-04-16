package ai

import (
	"context"
	"io"
	"net/http"
	"strings"
	"testing"
)

func TestNewTestSafeHTTPTransportReturnsDeterministicResponse(t *testing.T) {
	transport := NewTestSafeHTTPTransport()
	request, err := http.NewRequestWithContext(context.Background(), http.MethodPost, "https://example.com", nil)
	if err != nil {
		t.Fatalf("expected request build to succeed: %v", err)
	}

	response, err := transport.Do(request)

	if err != nil {
		t.Fatalf("expected deterministic transport to succeed: %v", err)
	}
	defer func() {
		_ = response.Body.Close()
	}()
	if response.StatusCode != http.StatusOK {
		t.Fatalf("expected status 200, got %d", response.StatusCode)
	}
	bodyBytes, err := io.ReadAll(response.Body)
	if err != nil {
		t.Fatalf("expected response body read to succeed: %v", err)
	}
	body := string(bodyBytes)
	if !strings.Contains(body, "テスト用の決定論的な AI 応答") {
		t.Fatalf("expected deterministic response body text, got %s", body)
	}
	if !strings.Contains(body, "\"choices\"") || !strings.Contains(body, "\"candidates\"") {
		t.Fatalf("expected deterministic response payload to contain openai and gemini fields, got %s", body)
	}
}

func TestNewTestSafeHTTPTransportWithResponseOverridesPayload(t *testing.T) {
	transport := NewTestSafeHTTPTransportWithResponse("custom response")
	request, err := http.NewRequestWithContext(context.Background(), http.MethodPost, "https://example.com", nil)
	if err != nil {
		t.Fatalf("expected request build to succeed: %v", err)
	}

	response, err := transport.Do(request)

	if err != nil {
		t.Fatalf("expected deterministic transport to succeed: %v", err)
	}
	defer func() {
		_ = response.Body.Close()
	}()
	bodyBytes, err := io.ReadAll(response.Body)
	if err != nil {
		t.Fatalf("expected response body read to succeed: %v", err)
	}
	if !strings.Contains(string(bodyBytes), "custom response") {
		t.Fatalf("expected custom response text in payload, got %s", string(bodyBytes))
	}
}
