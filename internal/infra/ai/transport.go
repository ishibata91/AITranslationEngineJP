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
	return &deterministicHTTPTransport{responseText: defaultTestSafeText}
}

// NewTestSafeHTTPTransportWithResponse creates a deterministic test-safe transport seam with the given response text.
// Passing an empty string yields an empty-text response, which causes provider parsers to return an empty-response error.
func NewTestSafeHTTPTransportWithResponse(responseText string) HTTPTransport {
	return &deterministicHTTPTransport{responseText: strings.TrimSpace(responseText)}
}

type deterministicHTTPTransport struct {
	responseText string
}

func (transport *deterministicHTTPTransport) Do(_ *http.Request) (*http.Response, error) {
	payload := map[string]interface{}{
		"candidates": []map[string]interface{}{
			{
				"content": map[string]interface{}{
					"parts": []map[string]string{{"text": transport.responseText}},
				},
			},
		},
		"choices": []map[string]interface{}{
			{
				"message": map[string]string{"content": transport.responseText},
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
