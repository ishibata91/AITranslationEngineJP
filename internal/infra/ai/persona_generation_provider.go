package ai

import (
	"bytes"
	"context"
	"fmt"
	"net/http"
)

type deterministicProvider struct {
	transport HTTPTransport
	testSafe  bool
}

func (provider deterministicProvider) Generate(
	ctx context.Context,
	request ProviderRequest,
) (ProviderResponse, error) {
	if !provider.testSafe {
		return ProviderResponse{}, fmt.Errorf("fake provider requires test-safe transport")
	}
	requestBytes, err := newOpenAICompatibleRequestBytes(request.Model, request.Prompt)
	if err != nil {
		return ProviderResponse{}, err
	}
	httpRequest, err := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		"https://fake.provider.local/chat/completions",
		bytes.NewReader(requestBytes),
	)
	if err != nil {
		return ProviderResponse{}, fmt.Errorf("build ai provider request: %w", err)
	}
	httpRequest.Header.Set("Content-Type", "application/json")
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
	return ProviderResponse{
		Text: text,
		DebugLog: buildProviderDebugLog(
			request.Prompt,
			requestBytes,
			httpRequest.Header,
		),
	}, nil
}
