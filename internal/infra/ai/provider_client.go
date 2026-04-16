package ai

import (
	"context"
	"fmt"
	"net/http"
	"strings"
)

// ProviderClientOption configures provider client runtime values.
type ProviderClientOption func(client *ProviderClient)

// WithLMStudioBaseURL sets the LM Studio OpenAI-compatible base URL.
func WithLMStudioBaseURL(baseURL string) ProviderClientOption {
	return func(client *ProviderClient) {
		if client == nil {
			return
		}
		client.lmStudioBaseURL = strings.TrimSpace(baseURL)
	}
}

// WithXAIBaseURL sets the xAI OpenAI-compatible base URL.
func WithXAIBaseURL(baseURL string) ProviderClientOption {
	return func(client *ProviderClient) {
		if client == nil {
			return
		}
		client.xaiBaseURL = strings.TrimSpace(baseURL)
	}
}

// ProviderClient sends provider-backed AI text requests.
type ProviderClient struct {
	transport       HTTPTransport
	testSafe        bool
	lmStudioBaseURL string
	xaiBaseURL      string
	providers       map[string]provider
}

// NewProviderClient creates a provider client backed by an HTTP transport.
func NewProviderClient(
	transport HTTPTransport,
	options ...ProviderClientOption,
) *ProviderClient {
	if transport == nil {
		transport = http.DefaultClient
	}
	_, testSafe := transport.(testSafeHTTPTransport)
	client := &ProviderClient{
		transport:       transport,
		testSafe:        testSafe,
		lmStudioBaseURL: lmStudioDefaultBaseURL,
		xaiBaseURL:      xaiDefaultBaseURL,
	}
	for _, option := range options {
		if option == nil {
			continue
		}
		option(client)
	}
	client.lmStudioBaseURL = normalizeBaseURL(client.lmStudioBaseURL, lmStudioDefaultBaseURL)
	client.xaiBaseURL = normalizeBaseURL(client.xaiBaseURL, xaiDefaultBaseURL)
	client.providers = client.newProviderRegistry()
	return client
}

// ProviderRequestsAreTestSafe reports whether this client can avoid paid real AI APIs.
func (client *ProviderClient) ProviderRequestsAreTestSafe() bool {
	return client != nil && client.testSafe
}

// GenerateText sends a provider request and returns provider-agnostic generated text.
func (client *ProviderClient) GenerateText(
	ctx context.Context,
	providerID string,
	request ProviderRequest,
) (ProviderResponse, error) {
	if client == nil || client.transport == nil {
		return ProviderResponse{}, fmt.Errorf("ai provider transport is required")
	}
	registry := client.providerRegistry()
	providerKey := strings.ToLower(strings.TrimSpace(providerID))
	concreteProvider, ok := registry[providerKey]
	if !ok {
		return ProviderResponse{}, fmt.Errorf("unsupported ai provider: %s", providerID)
	}
	response, err := concreteProvider.Generate(ctx, request)
	if err != nil {
		return ProviderResponse{}, fmt.Errorf("generate ai provider response: %w", err)
	}
	return response, nil
}

func (client *ProviderClient) providerRegistry() map[string]provider {
	if len(client.providers) == 0 {
		client.providers = client.newProviderRegistry()
	}
	return client.providers
}

func (client *ProviderClient) newProviderRegistry() map[string]provider {
	return map[string]provider{
		ProviderGemini: geminiProvider{
			transport: client.transport,
		},
		ProviderLMStudio: openAICompatibleProvider{
			transport:      client.transport,
			baseURL:        client.lmStudioBaseURL,
			apiKeyOptional: true,
		},
		ProviderXAI: openAICompatibleProvider{
			transport:      client.transport,
			baseURL:        client.xaiBaseURL,
			apiKeyOptional: client.testSafe,
		},
	}
}

func normalizeBaseURL(candidate string, fallback string) string {
	trimmedCandidate := strings.TrimSpace(candidate)
	if trimmedCandidate != "" {
		return trimmedCandidate
	}
	return strings.TrimSpace(fallback)
}
