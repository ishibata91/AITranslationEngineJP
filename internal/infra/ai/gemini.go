package ai

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
)

type geminiProvider struct {
	transport HTTPTransport
}

func (provider geminiProvider) Generate(
	ctx context.Context,
	request ProviderRequest,
) (ProviderResponse, error) {
	requestBytes, err := json.Marshal(geminiGenerateRequest{
		Contents: []geminiContent{
			{
				Role:  "user",
				Parts: []geminiPart{{Text: request.Prompt}},
			},
		},
	})
	if err != nil {
		return ProviderResponse{}, fmt.Errorf("marshal ai provider request: %w", err)
	}
	httpRequest, err := newGeminiRequest(ctx, request.Model, request.APIKey, requestBytes)
	if err != nil {
		return ProviderResponse{}, err
	}
	httpResponse, err := callProviderTransport(provider.transport, httpRequest)
	if err != nil {
		return ProviderResponse{}, err
	}
	defer func() {
		_ = httpResponse.Body.Close()
	}()
	text, err := readGeminiResponse(httpResponse)
	if err != nil {
		return ProviderResponse{}, err
	}
	return ProviderResponse{Text: text}, nil
}

type geminiGenerateRequest struct {
	Contents []geminiContent `json:"contents"`
}

type geminiContent struct {
	Role  string       `json:"role,omitempty"`
	Parts []geminiPart `json:"parts"`
}

type geminiPart struct {
	Text string `json:"text"`
}

type geminiGenerateResponse struct {
	Candidates []struct {
		Content struct {
			Parts []geminiPart `json:"parts"`
		} `json:"content"`
	} `json:"candidates"`
	Error *struct {
		Message string `json:"message"`
	} `json:"error,omitempty"`
}

func newGeminiRequest(
	ctx context.Context,
	model string,
	apiKey string,
	requestBytes []byte,
) (*http.Request, error) {
	trimmedModel := strings.TrimSpace(model)
	if trimmedModel == "" {
		return nil, fmt.Errorf("model is required")
	}
	endpoint := fmt.Sprintf(geminiEndpointTemplate, url.PathEscape(trimmedModel))
	query := url.Values{}
	query.Set("key", strings.TrimSpace(apiKey))
	request, err := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		endpoint+"?"+query.Encode(),
		bytes.NewReader(requestBytes),
	)
	if err != nil {
		return nil, fmt.Errorf("build ai provider request: %w", err)
	}
	request.Header.Set("Content-Type", "application/json")
	return request, nil
}

func readGeminiResponse(response *http.Response) (string, error) {
	responseBytes, err := io.ReadAll(response.Body)
	if err != nil {
		return "", fmt.Errorf("read ai provider response: %w", err)
	}

	parsed := geminiGenerateResponse{}
	if err := json.Unmarshal(responseBytes, &parsed); err != nil {
		return "", fmt.Errorf("parse ai provider response: %w", err)
	}
	if parsed.Error != nil {
		return "", fmt.Errorf("ai provider response error: %s", strings.TrimSpace(parsed.Error.Message))
	}
	if response.StatusCode < http.StatusOK || response.StatusCode >= http.StatusMultipleChoices {
		return "", fmt.Errorf("ai provider request failed: status=%d", response.StatusCode)
	}
	text := extractGeminiText(parsed)
	if text == "" {
		return "", errors.New(providerResponseEmptyError)
	}
	return text, nil
}

func extractGeminiText(response geminiGenerateResponse) string {
	for _, candidate := range response.Candidates {
		for _, part := range candidate.Content.Parts {
			text := strings.TrimSpace(part.Text)
			if text != "" {
				return text
			}
		}
	}
	return ""
}
