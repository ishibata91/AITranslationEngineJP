package ai

import (
	"context"
	"fmt"
	"net/http"
	"strings"
)

// ProviderCredentialResolver resolves one redacted credential reference to a provider secret.
type ProviderCredentialResolver interface {
	ResolveProviderCredential(ctx context.Context, providerID string, credentialRef string) (string, error)
}

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

// WithProviderCredentialResolver sets the credential resolver used by provider-backed request helpers.
func WithProviderCredentialResolver(resolver ProviderCredentialResolver) ProviderClientOption {
	return func(client *ProviderClient) {
		if client == nil {
			return
		}
		client.credentialResolver = resolver
	}
}

// WithProviderCredentialLoader wires one credential loader as the provider credential resolver.
func WithProviderCredentialLoader(loader func(context.Context, string) (string, error)) ProviderClientOption {
	return func(client *ProviderClient) {
		if client == nil || loader == nil {
			return
		}
		client.credentialResolver = providerCredentialLoaderResolver{load: loader}
	}
}

// WithDeterministicProviderRegistry replaces every supported runtime provider with the deterministic fake provider.
func WithDeterministicProviderRegistry(responseText string) ProviderClientOption {
	return func(client *ProviderClient) {
		if client == nil {
			return
		}
		client.providerOverride = deterministicProviderRegistry(responseText, true)
	}
}

// ProviderClient sends provider-backed AI text requests.
type ProviderClient struct {
	transport          HTTPTransport
	testSafe           bool
	lmStudioBaseURL    string
	xaiBaseURL         string
	credentialResolver ProviderCredentialResolver
	providers          map[string]provider
	providerOverride   map[string]provider
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

type providerCredentialLoaderResolver struct {
	load func(context.Context, string) (string, error)
}

func (resolver providerCredentialLoaderResolver) ResolveProviderCredential(
	ctx context.Context,
	_ string,
	credentialRef string,
) (string, error) {
	if resolver.load == nil {
		return "", nil
	}
	loaded, err := resolver.load(ctx, strings.TrimSpace(credentialRef))
	if err != nil {
		return "", fmt.Errorf("load provider credential: %w", err)
	}
	return strings.TrimSpace(loaded), nil
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

// ListProviderModelsWithEndpoint loads models through the provider registry selected by DI.
func (client *ProviderClient) ListProviderModelsWithEndpoint(
	ctx context.Context,
	providerID string,
	endpoint string,
	apiKey string,
) ([]ProviderModelOption, error) {
	if client == nil || client.transport == nil {
		return nil, fmt.Errorf("ai provider transport is required")
	}
	registry := client.providerRegistry()
	providerKey := strings.ToLower(strings.TrimSpace(providerID))
	concreteProvider, ok := registry[providerKey]
	if !ok {
		return nil, fmt.Errorf("unsupported ai provider: %s", providerID)
	}
	models, err := concreteProvider.ListModels(ctx, strings.TrimSpace(apiKey), strings.TrimSpace(endpoint))
	if err != nil {
		return nil, fmt.Errorf("list ai provider models: %w", err)
	}
	return models, nil
}

// TranslateTerm sends one provider request unit and parses a correlated JSON term translation response.
func (client *ProviderClient) TranslateTerm(
	ctx context.Context,
	providerID string,
	model string,
	apiKey string,
	endpointSummary string,
	prompt string,
) (TermTranslationResponse, error) {
	response, err := client.GenerateText(ctx, providerID, ProviderRequest{
		Model:           model,
		APIKey:          apiKey,
		Prompt:          prompt,
		EndpointSummary: strings.TrimSpace(endpointSummary),
	})
	if err != nil {
		return TermTranslationResponse{}, newTermTranslationError(
			TermTranslationErrorKindProviderFailure,
			isProviderExecutionRetryable(err),
			err,
		)
	}
	return parseTermTranslationResponse(response.Text)
}

// GeneratePersona sends one correlated NPC persona request unit and parses a strict JSON persona response.
func (client *ProviderClient) GeneratePersona(
	ctx context.Context,
	providerID string,
	model string,
	executionMode string,
	credentialRef string,
	endpointSummary string,
	prompt string,
) (PersonaGenerationResponse, error) {
	normalizedExecutionMode := strings.ToLower(strings.TrimSpace(executionMode))
	if normalizedExecutionMode == "" {
		return PersonaGenerationResponse{}, newPersonaGenerationError(
			PersonaGenerationErrorKindProviderFailure,
			false,
			fmt.Errorf("persona generation execution mode is required"),
		)
	}
	if normalizedExecutionMode != personaGenerationExecutionModeSingleRequest {
		return PersonaGenerationResponse{}, newPersonaGenerationError(
			PersonaGenerationErrorKindProviderFailure,
			false,
			fmt.Errorf("unsupported persona generation execution mode: %s", executionMode),
		)
	}

	apiKey, err := client.resolveProviderCredential(ctx, providerID, credentialRef)
	if err != nil {
		return PersonaGenerationResponse{}, newPersonaGenerationError(
			PersonaGenerationErrorKindProviderFailure,
			isProviderExecutionRetryable(err),
			err,
		)
	}
	response, err := client.GenerateText(ctx, providerID, ProviderRequest{
		Model:           model,
		APIKey:          apiKey,
		Prompt:          prompt,
		EndpointSummary: strings.TrimSpace(endpointSummary),
	})
	if err != nil {
		return PersonaGenerationResponse{}, newPersonaGenerationError(
			PersonaGenerationErrorKindProviderFailure,
			isProviderExecutionRetryable(err),
			err,
		)
	}
	parsed, err := parsePersonaGenerationResponse(response.Text)
	if err != nil {
		return PersonaGenerationResponse{}, err
	}
	parsed.ExecutionMode = normalizedExecutionMode
	parsed.PromptDigest = promptDigestSHA256(prompt)
	parsed.DebugLog = response.DebugLog
	return parsed, nil
}

// GenerateBodyTranslation sends one correlated field translation request unit and parses a strict JSON response.
func (client *ProviderClient) GenerateBodyTranslation(
	ctx context.Context,
	providerID string,
	model string,
	executionMode string,
	credentialRef string,
	endpointSummary string,
	prompt string,
) (BodyTranslationResponse, error) {
	normalizedExecutionMode := strings.ToLower(strings.TrimSpace(executionMode))
	if normalizedExecutionMode == "" {
		return BodyTranslationResponse{}, newBodyTranslationError(
			BodyTranslationErrorKindProviderFailure,
			false,
			fmt.Errorf("body translation execution mode is required"),
		)
	}
	if normalizedExecutionMode != bodyTranslationExecutionModeSingleRequest {
		return BodyTranslationResponse{}, newBodyTranslationError(
			BodyTranslationErrorKindProviderFailure,
			false,
			fmt.Errorf("unsupported body translation execution mode: %s", executionMode),
		)
	}

	apiKey, err := client.resolveProviderCredential(ctx, providerID, credentialRef)
	if err != nil {
		return BodyTranslationResponse{}, newBodyTranslationError(
			BodyTranslationErrorKindProviderFailure,
			isProviderExecutionRetryable(err),
			err,
		)
	}
	response, err := client.GenerateText(ctx, providerID, ProviderRequest{
		Model:           model,
		APIKey:          apiKey,
		Prompt:          prompt,
		EndpointSummary: strings.TrimSpace(endpointSummary),
	})
	if err != nil {
		return BodyTranslationResponse{}, newBodyTranslationError(
			BodyTranslationErrorKindProviderFailure,
			isProviderExecutionRetryable(err),
			err,
		)
	}
	parsed, err := parseBodyTranslationResponse(response.Text)
	if err != nil {
		return BodyTranslationResponse{}, err
	}
	parsed.ExecutionMode = normalizedExecutionMode
	parsed.PromptDigest = promptDigestSHA256(prompt)
	return parsed, nil
}

func (client *ProviderClient) providerRegistry() map[string]provider {
	if len(client.providers) == 0 {
		client.providers = client.newProviderRegistry()
	}
	return client.providers
}

func (client *ProviderClient) newProviderRegistry() map[string]provider {
	if len(client.providerOverride) != 0 {
		return cloneProviderRegistry(client.providerOverride)
	}
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
		ProviderFake: deterministicProviderFromTransport(client.transport),
	}
}

func deterministicProviderFromTransport(transport HTTPTransport) deterministicProvider {
	deterministicTransport, ok := transport.(*deterministicHTTPTransport)
	if !ok || deterministicTransport == nil {
		return deterministicProvider{}
	}
	return deterministicProvider{
		responseText:          deterministicTransport.responseText,
		useConfiguredResponse: deterministicTransport.useConfiguredResponse,
	}
}

func deterministicProviderRegistry(responseText string, useConfiguredResponse bool) map[string]provider {
	deterministic := deterministicProvider{
		responseText:          strings.TrimSpace(responseText),
		useConfiguredResponse: useConfiguredResponse,
	}
	return map[string]provider{
		ProviderGemini:   deterministic,
		ProviderLMStudio: deterministic,
		ProviderXAI:      deterministic,
		ProviderFake:     deterministic,
	}
}

func cloneProviderRegistry(source map[string]provider) map[string]provider {
	cloned := make(map[string]provider, len(source))
	for key, value := range source {
		cloned[strings.ToLower(strings.TrimSpace(key))] = value
	}
	return cloned
}

func normalizeBaseURL(candidate string, fallback string) string {
	trimmedCandidate := strings.TrimSpace(candidate)
	if trimmedCandidate != "" {
		return trimmedCandidate
	}
	return strings.TrimSpace(fallback)
}

func (client *ProviderClient) resolveProviderCredential(
	ctx context.Context,
	providerID string,
	credentialRef string,
) (string, error) {
	normalizedProvider := strings.ToLower(strings.TrimSpace(providerID))
	if client != nil && client.testSafe {
		return "", nil
	}
	switch normalizedProvider {
	case ProviderFake:
		return "", nil
	case ProviderLMStudio:
		if client == nil || client.credentialResolver == nil {
			return "", nil
		}
	default:
		if client == nil || client.credentialResolver == nil {
			return "", fmt.Errorf("provider credential resolver is required")
		}
	}
	resolved, err := client.credentialResolver.ResolveProviderCredential(ctx, normalizedProvider, strings.TrimSpace(credentialRef))
	if err != nil {
		return "", fmt.Errorf("resolve provider credential: %w", err)
	}
	if normalizedProvider == ProviderLMStudio {
		return strings.TrimSpace(resolved), nil
	}
	trimmed := strings.TrimSpace(resolved)
	if trimmed == "" {
		return "", fmt.Errorf("provider credential is unavailable")
	}
	return trimmed, nil
}
