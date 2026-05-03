package service

import (
	"context"
	"testing"
)

func TestBuildBodyTranslationPromptAcceptsSyncExecutionModeAlias(t *testing.T) {
	prompt, err := BuildBodyTranslationPrompt(BodyTranslationProviderRequest{
		ExecutionMode:       "sync",
		RequestUnitID:       "unit-1",
		FieldCorrelationKey: "field:1",
		RecordType:          "NPC_",
		FieldType:           "FULL",
		SourceText:          "Hello there",
	})
	if err != nil {
		t.Fatalf("expected sync alias prompt build success, got %v", err)
	}
	if prompt == "" {
		t.Fatal("expected prompt text")
	}
}

func TestBodyTranslationProviderAdapterNormalizesSyncExecutionMode(t *testing.T) {
	client := stubBodyTranslationProviderClient{
		response: stubBodyTranslationClientResponse{
			Items: []stubBodyTranslationClientItem{{
				RequestUnitID:       "unit-1",
				FieldCorrelationKey: "field:1",
				TranslatedText:      "こんにちは",
			}},
			ExecutionMode: BodyTranslationExecutionModeSingleRequest,
			PromptDigest:  "sha256:from-client",
		},
	}
	adapter := NewBodyTranslationProviderAdapter(&client)

	result := adapter.TranslateBodyField(context.Background(), BodyTranslationProviderRequest{
		Provider:            BodyTranslationProviderLMStudio,
		Model:               "local-model",
		ExecutionMode:       "sync",
		CredentialRef:       "lmstudio-primary",
		RequestUnitID:       "unit-1",
		FieldCorrelationKey: "field:1",
		RecordType:          "NPC_",
		FieldType:           "FULL",
		SourceText:          "Hello there",
	})

	if result.Failure != nil {
		t.Fatalf("expected success, got failure %#v", result.Failure)
	}
	if client.capturedExecutionMode != BodyTranslationExecutionModeSingleRequest {
		t.Fatalf("expected client execution mode %q, got %q", BodyTranslationExecutionModeSingleRequest, client.capturedExecutionMode)
	}
	if result.AuditSummary.ExecutionMode != BodyTranslationExecutionModeSingleRequest {
		t.Fatalf("expected normalized audit execution mode, got %#v", result.AuditSummary)
	}
}

type stubBodyTranslationProviderClient struct {
	response              stubBodyTranslationClientResponse
	err                   error
	capturedExecutionMode string
}

func (stub *stubBodyTranslationProviderClient) GenerateBodyTranslation(
	_ context.Context,
	_ string,
	_ string,
	executionMode string,
	_ string,
	_ string,
) (stubBodyTranslationClientResponse, error) {
	stub.capturedExecutionMode = executionMode
	if stub.err != nil {
		return stubBodyTranslationClientResponse{}, stub.err
	}
	return stub.response, nil
}

func (stub *stubBodyTranslationProviderClient) ProviderRequestsAreTestSafe() bool {
	return true
}

type stubBodyTranslationClientResponse struct {
	Items         []stubBodyTranslationClientItem
	ExecutionMode string
	PromptDigest  string
}

type stubBodyTranslationClientItem struct {
	RequestUnitID       string
	FieldCorrelationKey string
	TranslatedText      string
}
