package ai

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"strings"
)

const (
	providerModelsXAIBaseURLEnv      = "AITRANSLATIONENGINEJP_MASTER_PERSONA_XAI_BASE_URL"
	providerModelsLMStudioBaseURLEnv = "AITRANSLATIONENGINEJP_MASTER_PERSONA_LM_STUDIO_BASE_URL"
)

// ProviderModelOption is one public model option returned by provider model-list APIs.
type ProviderModelOption struct {
	ModelID string
	Label   string
}

// ListProviderModelsRequest carries the provider model-list request input.
type ListProviderModelsRequest struct {
	ProviderID      string
	APIKey          string
	GeminiBaseURL   string
	OpenAIBaseURL   string
	LMStudioBaseURL string
	XAIBaseURL      string
}

// ProviderModelListLoader keeps model-list transport and environment resolution inside the infra adapter.
type ProviderModelListLoader struct {
	transport       HTTPTransport
	openAIBaseURL   string
	lmStudioBaseURL string
	xaiBaseURL      string
	providers       map[string]provider
}

// ProviderModelListLoaderOption configures optional base-URL overrides.
type ProviderModelListLoaderOption func(loader *ProviderModelListLoader)

// WithProviderModelListOpenAIBaseURL overrides the OpenAI-compatible base URL.
func WithProviderModelListOpenAIBaseURL(baseURL string) ProviderModelListLoaderOption {
	return func(loader *ProviderModelListLoader) {
		if loader != nil {
			loader.openAIBaseURL = strings.TrimSpace(baseURL)
		}
	}
}

// WithProviderModelListLMStudioBaseURL overrides the LM Studio base URL.
func WithProviderModelListLMStudioBaseURL(baseURL string) ProviderModelListLoaderOption {
	return func(loader *ProviderModelListLoader) {
		if loader != nil {
			loader.lmStudioBaseURL = strings.TrimSpace(baseURL)
		}
	}
}

// WithProviderModelListXAIBaseURL overrides the XAI base URL.
func WithProviderModelListXAIBaseURL(baseURL string) ProviderModelListLoaderOption {
	return func(loader *ProviderModelListLoader) {
		if loader != nil {
			loader.xaiBaseURL = strings.TrimSpace(baseURL)
		}
	}
}

// WithProviderModelListDeterministicProviders replaces model-list providers with the deterministic fake provider.
func WithProviderModelListDeterministicProviders() ProviderModelListLoaderOption {
	return func(loader *ProviderModelListLoader) {
		if loader != nil {
			loader.providers = deterministicProviderRegistry("", false)
		}
	}
}

// NewProviderModelListLoader creates the infra adapter used by backend services.
func NewProviderModelListLoader(
	transport HTTPTransport,
	options ...ProviderModelListLoaderOption,
) *ProviderModelListLoader {
	loader := &ProviderModelListLoader{transport: transport}
	for _, option := range options {
		if option != nil {
			option(loader)
		}
	}
	return loader
}

// ListProviderModels resolves provider model lists through the configured transport and env-backed base URLs.
func (loader *ProviderModelListLoader) ListProviderModels(
	ctx context.Context,
	providerID string,
	apiKey string,
) ([]ProviderModelOption, error) {
	return loader.ListProviderModelsWithEndpoint(ctx, providerID, "", apiKey)
}

// ListProviderModelsWithEndpoint resolves provider model lists through the configured transport and one explicit endpoint.
func (loader *ProviderModelListLoader) ListProviderModelsWithEndpoint(
	ctx context.Context,
	providerID string,
	endpoint string,
	apiKey string,
) ([]ProviderModelOption, error) {
	if loader == nil {
		return nil, fmt.Errorf("provider model list loader is required")
	}
	trimmedEndpoint := strings.TrimSpace(endpoint)
	request := ListProviderModelsRequest{
		ProviderID:      providerID,
		APIKey:          apiKey,
		GeminiBaseURL:   firstNonEmptyProviderModelsURL(trimmedEndpoint),
		OpenAIBaseURL:   firstNonEmptyProviderModelsURL(trimmedEndpoint, strings.TrimSpace(loader.openAIBaseURL)),
		LMStudioBaseURL: firstNonEmptyProviderModelsURL(trimmedEndpoint, loader.lmStudioBaseURL, os.Getenv(providerModelsLMStudioBaseURLEnv)),
		XAIBaseURL:      firstNonEmptyProviderModelsURL(trimmedEndpoint, loader.xaiBaseURL, os.Getenv(providerModelsXAIBaseURLEnv)),
	}
	if len(loader.providers) != 0 {
		return listProviderModelsWithRegistry(ctx, cloneProviderRegistry(loader.providers), request)
	}
	return ListProviderModels(ctx, loader.transport, request)
}

// ListProviderModels loads provider models through the shared HTTP transport seam.
func ListProviderModels(
	ctx context.Context,
	transport HTTPTransport,
	request ListProviderModelsRequest,
) ([]ProviderModelOption, error) {
	if transport == nil {
		return nil, fmt.Errorf("ai provider transport is required")
	}
	providers := providerModelListRegistry(transport, request)
	return listProviderModelsWithRegistry(ctx, providers, request)
}

func providerModelListRegistry(transport HTTPTransport, request ListProviderModelsRequest) map[string]provider {
	return map[string]provider{
		ProviderGemini: geminiProvider{transport: transport},
		ProviderXAI: openAICompatibleProvider{
			transport:      transport,
			baseURL:        normalizeBaseURL(request.XAIBaseURL, xaiDefaultBaseURL),
			apiKeyOptional: false,
		},
		ProviderLMStudio: openAICompatibleProvider{
			transport:      transport,
			baseURL:        normalizeBaseURL(request.LMStudioBaseURL, lmStudioDefaultBaseURL),
			apiKeyOptional: true,
		},
		ProviderFake: deterministicProvider{},
	}
}

func listProviderModelsWithRegistry(
	ctx context.Context,
	providers map[string]provider,
	request ListProviderModelsRequest,
) ([]ProviderModelOption, error) {
	switch strings.ToLower(strings.TrimSpace(request.ProviderID)) {
	case ProviderGemini:
		return listModelsFromProvider(ctx, providers[ProviderGemini], request.APIKey, normalizeBaseURL(request.GeminiBaseURL, geminiDefaultBaseURL))
	case ProviderXAI:
		return listModelsFromProvider(ctx, providers[ProviderXAI], request.APIKey, normalizeBaseURL(request.XAIBaseURL, xaiDefaultBaseURL))
	case ProviderLMStudio:
		return listModelsFromProvider(ctx, providers[ProviderLMStudio], request.APIKey, normalizeBaseURL(request.LMStudioBaseURL, lmStudioDefaultBaseURL))
	default:
		return listModelsFromProvider(ctx, openAICompatibleProvider{
			transport:      providerModelListTransport(providers),
			baseURL:        normalizeBaseURL(request.OpenAIBaseURL, "https://api.openai.com/v1"),
			apiKeyOptional: false,
		}, request.APIKey, normalizeBaseURL(request.OpenAIBaseURL, "https://api.openai.com/v1"))
	}
}

func listModelsFromProvider(
	ctx context.Context,
	implementation provider,
	apiKey string,
	endpointSummary string,
) ([]ProviderModelOption, error) {
	if implementation == nil {
		return nil, fmt.Errorf("provider model list implementation is required")
	}
	models, err := implementation.ListModels(ctx, apiKey, endpointSummary)
	if err != nil {
		return nil, fmt.Errorf("list provider models: %w", err)
	}
	return models, nil
}

func providerModelListTransport(providers map[string]provider) HTTPTransport {
	if gemini, ok := providers[ProviderGemini].(geminiProvider); ok {
		return gemini.transport
	}
	return nil
}

type openAICompatibleModelsResponse struct {
	Data []struct {
		ID string `json:"id"`
	} `json:"data"`
}

func listOpenAICompatibleModels(
	ctx context.Context,
	transport HTTPTransport,
	baseURL string,
	apiKey string,
	apiKeyOptional bool,
) ([]ProviderModelOption, error) {
	endpoint, err := openAICompatibleModelsEndpoint(baseURL)
	if err != nil {
		return nil, err
	}
	request, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint, nil)
	if err != nil {
		return nil, fmt.Errorf("build provider model list request: %w", err)
	}
	if trimmedAPIKey := strings.TrimSpace(apiKey); trimmedAPIKey != "" {
		request.Header.Set("Authorization", "Bearer "+trimmedAPIKey)
	} else if !apiKeyOptional {
		return nil, fmt.Errorf("api key is required")
	}
	response, err := callProviderTransport(transport, request)
	if err != nil {
		return nil, err
	}
	defer func() {
		_ = response.Body.Close()
	}()
	parsed := openAICompatibleModelsResponse{}
	if err := json.NewDecoder(response.Body).Decode(&parsed); err != nil {
		return nil, fmt.Errorf("parse provider model list response: %w", err)
	}
	if response.StatusCode < http.StatusOK || response.StatusCode >= http.StatusMultipleChoices {
		return nil, fmt.Errorf("provider model list request failed: status=%d", response.StatusCode)
	}
	result := make([]ProviderModelOption, 0, len(parsed.Data))
	for _, item := range parsed.Data {
		modelID := strings.TrimSpace(item.ID)
		if modelID == "" {
			continue
		}
		result = append(result, ProviderModelOption{ModelID: modelID, Label: modelID})
	}
	return result, nil
}

func openAICompatibleModelsEndpoint(baseURL string) (string, error) {
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
	parsedURL.Path = trimmedPath + "/models"
	parsedURL.RawQuery = ""
	parsedURL.Fragment = ""
	return parsedURL.String(), nil
}

type geminiModelsResponse struct {
	Models []struct {
		Name        string `json:"name"`
		DisplayName string `json:"displayName"`
	} `json:"models"`
}

func listGeminiModels(
	ctx context.Context,
	transport HTTPTransport,
	baseURL string,
	apiKey string,
) ([]ProviderModelOption, error) {
	endpoint, err := openAICompatibleModelsEndpoint(baseURL)
	if err != nil {
		return nil, fmt.Errorf("build gemini model list endpoint: %w", err)
	}
	request, err := http.NewRequestWithContext(
		ctx,
		http.MethodGet,
		endpoint,
		nil,
	)
	if err != nil {
		return nil, fmt.Errorf("build gemini model list request: %w", err)
	}
	if trimmedAPIKey := strings.TrimSpace(apiKey); trimmedAPIKey != "" {
		request.Header.Set("x-goog-api-key", trimmedAPIKey)
	} else {
		return nil, fmt.Errorf("api key is required")
	}
	response, err := callProviderTransport(transport, request)
	if err != nil {
		return nil, err
	}
	defer func() {
		_ = response.Body.Close()
	}()
	parsed := geminiModelsResponse{}
	if err := json.NewDecoder(response.Body).Decode(&parsed); err != nil {
		return nil, fmt.Errorf("parse gemini model list response: %w", err)
	}
	if response.StatusCode < http.StatusOK || response.StatusCode >= http.StatusMultipleChoices {
		return nil, fmt.Errorf("gemini model list request failed: status=%d", response.StatusCode)
	}
	result := make([]ProviderModelOption, 0, len(parsed.Models))
	for _, item := range parsed.Models {
		modelID := strings.TrimPrefix(strings.TrimSpace(item.Name), "models/")
		if modelID == "" {
			continue
		}
		label := strings.TrimSpace(item.DisplayName)
		if label == "" {
			label = modelID
		}
		result = append(result, ProviderModelOption{ModelID: modelID, Label: label})
	}
	return result, nil
}

func firstNonEmptyProviderModelsURL(values ...string) string {
	for _, value := range values {
		trimmed := strings.TrimSpace(value)
		if trimmed != "" {
			return trimmed
		}
	}
	return ""
}
