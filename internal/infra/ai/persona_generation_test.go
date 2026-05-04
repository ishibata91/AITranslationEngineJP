package ai

import (
	"context"
	"errors"
	"net/http"
	"strings"
	"testing"
)

func TestParsePersonaGenerationResponseObject(t *testing.T) {
	response, err := parsePersonaGenerationResponse(`{"personas":[{"request_unit_id":"unit-1","npc_correlation_id":"npc-1","persona_body":" body "}]}`)
	if err != nil {
		t.Fatalf("expected parser success: %v", err)
	}
	if len(response.Items) != 1 {
		t.Fatalf("expected one item, got %d", len(response.Items))
	}
	if response.Items[0].RequestUnitID != "unit-1" || response.Items[0].NPCCorrelationID != "npc-1" || response.Items[0].PersonaBody != "body" {
		t.Fatalf("unexpected parsed item: %#v", response.Items[0])
	}
}

func TestParsePersonaGenerationResponseRejectsMissingCorrelation(t *testing.T) {
	_, err := parsePersonaGenerationResponse(`{"personas":[{"request_unit_id":"unit-1","persona_body":"body"}]}`)
	if err == nil {
		t.Fatal("expected invalid response error")
	}
	var personaErr *PersonaGenerationError
	if !errors.As(err, &personaErr) {
		t.Fatalf("expected typed persona generation error, got %T", err)
	}
	if personaErr.Kind != PersonaGenerationErrorKindInvalidProviderResponse || !personaErr.Retryable {
		t.Fatalf("unexpected error metadata: %#v", personaErr)
	}
}

func TestProviderClientGeneratePersonaUsesResolverAndRedactsDebugLog(t *testing.T) {
	transport := &stubHTTPTransport{doFunc: func(_ *http.Request) (*http.Response, error) {
		return newHTTPJSONResponse(http.StatusOK, `{"choices":[{"message":{"content":"{\"personas\":[{\"request_unit_id\":\"unit-1\",\"npc_correlation_id\":\"npc-1\",\"persona_body\":\"guardian\"}]}"}}]}`), nil
	}}
	resolver := stubProviderCredentialResolver{
		resolveFunc: func(_ context.Context, providerID string, credentialRef string) (string, error) {
			if providerID != ProviderXAI || credentialRef != "cred-1" {
				t.Fatalf("unexpected resolver input: provider=%s credentialRef=%s", providerID, credentialRef)
			}
			return "super-secret-key", nil
		},
	}
	client := NewProviderClient(transport, WithProviderCredentialResolver(resolver))

	response, err := client.GeneratePersona(context.Background(), ProviderXAI, "grok-2", personaGenerationExecutionModeSingleRequest, "cred-1", "", "PERSONA_GENERATION_REQUEST_V1\nrequest_unit_id=unit-1\nnpc_correlation_id=npc-1")
	if err != nil {
		t.Fatalf("expected provider client success: %v", err)
	}
	if response.PromptDigest == "" || !strings.HasPrefix(response.PromptDigest, "sha256:") {
		t.Fatalf("expected prompt digest, got %q", response.PromptDigest)
	}
	if len(response.Items) != 1 || response.Items[0].PersonaBody != "guardian" {
		t.Fatalf("unexpected parsed persona response: %#v", response.Items)
	}
	if response.DebugLog.Headers["Authorization"] != "[REDACTED]" {
		t.Fatalf("expected authorization header to be redacted, got %#v", response.DebugLog.Headers)
	}
	if strings.Contains(response.DebugLog.RequestBody, "super-secret-key") {
		t.Fatalf("request body must not include secret: %s", response.DebugLog.RequestBody)
	}
	if !response.DebugLog.SecretRedacted {
		t.Fatal("expected debug log to report redacted secret state")
	}
}

func TestProviderClientGeneratePersonaSupportsFakeProviderWithFixedResponse(t *testing.T) {
	fixedResponse := `{"personas":[{"request_unit_id":"unit-fixed","npc_correlation_id":"npc-fixed","persona_body":"fixed persona"}]}`
	client := NewProviderClient(NewTestSafeHTTPTransportWithResponse(fixedResponse))

	response, err := client.GeneratePersona(
		context.Background(),
		ProviderFake,
		"fake-model",
		personaGenerationExecutionModeSingleRequest,
		"",
		"",
		"PERSONA_GENERATION_REQUEST_V1\nrequest_unit_id=unit-fixed\nnpc_correlation_id=npc-fixed",
	)
	if err != nil {
		t.Fatalf("expected fake provider success: %v", err)
	}
	if response.Items[0].RequestUnitID != "unit-fixed" || response.Items[0].NPCCorrelationID != "npc-fixed" {
		t.Fatalf("unexpected fake provider response correlation: %#v", response.Items[0])
	}
	if response.Items[0].PersonaBody != "fixed persona" {
		t.Fatalf("unexpected persona body: %q", response.Items[0].PersonaBody)
	}
}

type stubProviderCredentialResolver struct {
	resolveFunc func(ctx context.Context, providerID string, credentialRef string) (string, error)
}

func (stub stubProviderCredentialResolver) ResolveProviderCredential(
	ctx context.Context,
	providerID string,
	credentialRef string,
) (string, error) {
	if stub.resolveFunc == nil {
		return "", nil
	}
	return stub.resolveFunc(ctx, providerID, credentialRef)
}
