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

type openAICompatibleProvider struct {
	transport      HTTPTransport
	baseURL        string
	apiKeyOptional bool
}

func (provider openAICompatibleProvider) Generate(
	ctx context.Context,
	request ProviderRequest,
) (ProviderResponse, error) {
	requestBytes, err := newOpenAICompatibleRequestBytes(request.Model, request.Prompt)
	if err != nil {
		return ProviderResponse{}, err
	}
	httpRequest, err := newOpenAICompatibleRequest(
		ctx,
		provider.baseURL,
		request.APIKey,
		requestBytes,
		provider.apiKeyOptional,
	)
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
	text, err := readOpenAICompatibleResponse(httpResponse)
	if err != nil {
		return ProviderResponse{}, err
	}
	return ProviderResponse{Text: text}, nil
}

type openAICompatibleChatRequest struct {
	Model    string                    `json:"model"`
	Messages []openAICompatibleMessage `json:"messages"`
}

type openAICompatibleMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type openAICompatibleChatResponse struct {
	Choices []struct {
		Message struct {
			Content string `json:"content"`
		} `json:"message"`
	} `json:"choices"`
	Error *struct {
		Message string `json:"message"`
	} `json:"error,omitempty"`
}

func newOpenAICompatibleRequestBytes(model string, prompt string) ([]byte, error) {
	trimmedModel := strings.TrimSpace(model)
	if trimmedModel == "" {
		return nil, fmt.Errorf("model is required")
	}
	requestBytes, err := json.Marshal(openAICompatibleChatRequest{
		Model: trimmedModel,
		Messages: []openAICompatibleMessage{
			{Role: "user", Content: prompt},
		},
	})
	if err != nil {
		return nil, fmt.Errorf("marshal ai provider request: %w", err)
	}
	return requestBytes, nil
}

func newOpenAICompatibleRequest(
	ctx context.Context,
	baseURL string,
	apiKey string,
	requestBytes []byte,
	apiKeyOptional bool,
) (*http.Request, error) {
	endpoint, err := openAICompatibleEndpoint(baseURL)
	if err != nil {
		return nil, err
	}
	trimmedAPIKey := strings.TrimSpace(apiKey)
	if !apiKeyOptional && trimmedAPIKey == "" {
		return nil, fmt.Errorf("api key is required")
	}
	request, err := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		endpoint,
		bytes.NewReader(requestBytes),
	)
	if err != nil {
		return nil, fmt.Errorf("build ai provider request: %w", err)
	}
	request.Header.Set("Content-Type", "application/json")
	if trimmedAPIKey != "" {
		request.Header.Set("Authorization", "Bearer "+trimmedAPIKey)
	}
	return request, nil
}

func openAICompatibleEndpoint(baseURL string) (string, error) {
	trimmedBaseURL := strings.TrimSpace(baseURL)
	if trimmedBaseURL == "" {
		return "", fmt.Errorf("provider base url is required")
	}
	parsedURL, err := url.Parse(trimmedBaseURL)
	if err != nil {
		return "", fmt.Errorf("parse provider base url: %w", err)
	}
	if strings.TrimSpace(parsedURL.Scheme) == "" || strings.TrimSpace(parsedURL.Host) == "" {
		return "", fmt.Errorf("provider base url must be absolute: %s", trimmedBaseURL)
	}
	trimmedPath := strings.TrimSuffix(parsedURL.Path, "/")
	parsedURL.Path = trimmedPath + "/chat/completions"
	parsedURL.RawQuery = ""
	parsedURL.Fragment = ""
	return parsedURL.String(), nil
}

func readOpenAICompatibleResponse(response *http.Response) (string, error) {
	responseBytes, err := io.ReadAll(response.Body)
	if err != nil {
		return "", fmt.Errorf("read ai provider response: %w", err)
	}
	parsed := openAICompatibleChatResponse{}
	if err := json.Unmarshal(responseBytes, &parsed); err != nil {
		return "", fmt.Errorf("parse ai provider response: %w", err)
	}
	if parsed.Error != nil {
		return "", fmt.Errorf("ai provider response error: %s", strings.TrimSpace(parsed.Error.Message))
	}
	if response.StatusCode < http.StatusOK || response.StatusCode >= http.StatusMultipleChoices {
		return "", fmt.Errorf("ai provider request failed: status=%d", response.StatusCode)
	}
	text := extractOpenAICompatibleText(parsed)
	if text == "" {
		return "", errors.New(providerResponseEmptyError)
	}
	return text, nil
}

func extractOpenAICompatibleText(response openAICompatibleChatResponse) string {
	for _, choice := range response.Choices {
		text := strings.TrimSpace(choice.Message.Content)
		if text != "" {
			return text
		}
	}
	return ""
}
