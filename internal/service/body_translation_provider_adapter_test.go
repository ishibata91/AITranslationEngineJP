package service

import (
	"context"
	"reflect"
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

func TestBodyTranslationPromptInputExcludesProviderConnectionFields(t *testing.T) {
	inputType := reflect.TypeOf(BodyTranslationPromptInput{})
	for _, forbiddenField := range []string{"Provider", "Model", "CredentialRef", "EndpointSummary"} {
		if _, exists := inputType.FieldByName(forbiddenField); exists {
			t.Fatalf("expected BodyTranslationPromptInput to exclude provider connection field %s", forbiddenField)
		}
	}
}

func TestBodyTranslationPromptBuilderBuildsOneUnitEnvelope(t *testing.T) {
	envelope, err := NewBodyTranslationPromptBuilder().Build(BodyTranslationPromptInput{
		ExecutionMode:       "sync",
		RequestUnitID:       "unit-builder",
		FieldCorrelationKey: "field:42",
		RecordType:          "INFO",
		FieldType:           "FULL",
		SourceText:          "Hello <name>",
		PersonaSummary:      "warm speaker",
		ContextLines:        []string{"previous line"},
		CompleteMatchExclusions: []BodyTranslationDictionaryExactMatchExclusion{{
			SourceText:     "Dragonborn",
			TranslatedText: "ドラゴンボーン",
		}},
		PartialMatchConstraints: []BodyTranslationPartialMatchConstraint{{
			SourceText:          "Jarl",
			RequiredTranslation: "首長",
		}},
		ProtectedElements: []BodyTranslationProtectedElement{{
			ElementType: "tag",
			SourceText:  "<name>",
			Digest:      "sha256:protected-name",
		}},
	})
	if err != nil {
		t.Fatalf("expected prompt builder success, got %v", err)
	}
	if envelope.RequestShapeID != BodyTranslationRequestShapeV1 {
		t.Fatalf("expected body translation request shape, got %q", envelope.RequestShapeID)
	}
	if envelope.Summary.InputCount != 1 {
		t.Fatalf("expected one prompt input, got %#v", envelope.Summary)
	}
	if len(envelope.Summary.CorrelationIDs) != 2 ||
		envelope.Summary.CorrelationIDs[0] != "unit-builder" ||
		envelope.Summary.CorrelationIDs[1] != "field:42" {
		t.Fatalf("expected request identifiers in prompt summary, got %#v", envelope.Summary.CorrelationIDs)
	}
	if envelope.Summary.Counts["protected_elements"] != 1 ||
		envelope.Summary.Counts["complete_match_exclusions"] != 1 ||
		envelope.Summary.Counts["partial_match_constraints"] != 1 {
		t.Fatalf("expected prompt summary counts, got %#v", envelope.Summary.Counts)
	}
	if envelope.RawPrompt == "" || envelope.Digest == "" {
		t.Fatalf("expected raw prompt and digest in internal envelope, got %#v", envelope)
	}
}

func TestBodyTranslationProviderAdapterAuditSummaryDoesNotExposeRawPromptInputs(t *testing.T) {
	client := stubBodyTranslationProviderClient{
		response: stubBodyTranslationClientResponse{
			Items: []stubBodyTranslationClientItem{{
				RequestUnitID:       "unit-raw",
				FieldCorrelationKey: "field:7",
				TranslatedText:      "翻訳済み <tag>",
			}},
			ExecutionMode: BodyTranslationExecutionModeSingleRequest,
		},
	}
	adapter := NewBodyTranslationProviderAdapter(&client)

	result := adapter.TranslateBodyField(context.Background(), BodyTranslationProviderRequest{
		Provider:            BodyTranslationProviderLMStudio,
		Model:               "local-model",
		ExecutionMode:       BodyTranslationExecutionModeSingleRequest,
		CredentialRef:       "lmstudio-primary",
		RequestUnitID:       "unit-raw",
		FieldCorrelationKey: "field:7",
		RecordType:          "NPC_",
		FieldType:           "FULL",
		SourceText:          "Raw source text that must stay inside the provider prompt",
		CompleteMatchExclusions: []BodyTranslationDictionaryExactMatchExclusion{{
			SourceText:     "Raw exact source",
			TranslatedText: "Raw exact translation",
		}},
		PartialMatchConstraints: []BodyTranslationPartialMatchConstraint{{
			SourceText:          "Raw partial source",
			RequiredTranslation: "Raw partial translation",
		}},
		ProtectedElements: []BodyTranslationProtectedElement{{
			ElementType: "tag",
			SourceText:  "<tag>",
			Digest:      "sha256:protected-tag",
		}},
	})

	if result.Failure != nil {
		t.Fatalf("expected success, got failure %#v", result.Failure)
	}
	if client.capturedPrompt == "" {
		t.Fatal("expected provider client to receive raw prompt")
	}
	summary := result.AuditSummary.RequestSummary
	if summary.RequestUnitID != "unit-raw" || summary.FieldCorrelationKey != "field:7" {
		t.Fatalf("expected request identifiers in audit summary, got %#v", summary)
	}
	if summary.ProtectedElementCount != 1 ||
		summary.CompleteMatchExclusionCount != 1 ||
		summary.PartialMatchConstraintCount != 1 {
		t.Fatalf("expected request counts in audit summary, got %#v", summary)
	}
}

func TestBodyTranslationProviderAdapterRejectsMismatchedFieldCorrelationKey(t *testing.T) {
	// 翻訳項目識別子不一致は翻訳項目単位の invalid response として扱う。
	client := stubBodyTranslationProviderClient{
		response: stubBodyTranslationClientResponse{
			Items: []stubBodyTranslationClientItem{{
				RequestUnitID:       "unit-1",
				FieldCorrelationKey: "field:other",
				TranslatedText:      "翻訳済み",
			}},
			ExecutionMode: BodyTranslationExecutionModeSingleRequest,
		},
	}
	adapter := NewBodyTranslationProviderAdapter(&client)

	result := adapter.TranslateBodyField(context.Background(), BodyTranslationProviderRequest{
		Provider:            BodyTranslationProviderLMStudio,
		Model:               "local-model",
		ExecutionMode:       BodyTranslationExecutionModeSingleRequest,
		CredentialRef:       "lmstudio-primary",
		RequestUnitID:       "unit-1",
		FieldCorrelationKey: "field:1",
		RecordType:          "NPC_",
		FieldType:           "FULL",
		SourceText:          "Hello there",
	})

	if result.Failure == nil {
		t.Fatal("expected invalid provider response failure")
	}
	if result.Failure.Kind != BodyTranslationProviderErrorKindInvalidProviderResponse || !result.Failure.Retryable {
		t.Fatalf("unexpected failure: %#v", result.Failure)
	}
	if result.TranslatedCandidate != nil || result.ProtectionValidationTarget != nil {
		t.Fatalf("expected no translated candidates on invalid response, got %#v", result)
	}
}

func TestBodyTranslationProviderAdapterRejectsEmptyTranslatedText(t *testing.T) {
	// 空訳文は翻訳項目単位の invalid response として扱う。
	client := stubBodyTranslationProviderClient{
		response: stubBodyTranslationClientResponse{
			Items: []stubBodyTranslationClientItem{{
				RequestUnitID:       "unit-1",
				FieldCorrelationKey: "field:1",
				TranslatedText:      " ",
			}},
			ExecutionMode: BodyTranslationExecutionModeSingleRequest,
		},
	}
	adapter := NewBodyTranslationProviderAdapter(&client)

	result := adapter.TranslateBodyField(context.Background(), BodyTranslationProviderRequest{
		Provider:            BodyTranslationProviderLMStudio,
		Model:               "local-model",
		ExecutionMode:       BodyTranslationExecutionModeSingleRequest,
		CredentialRef:       "lmstudio-primary",
		RequestUnitID:       "unit-1",
		FieldCorrelationKey: "field:1",
		RecordType:          "NPC_",
		FieldType:           "FULL",
		SourceText:          "Hello there",
	})

	if result.Failure == nil {
		t.Fatal("expected invalid provider response failure")
	}
	if result.Failure.Kind != BodyTranslationProviderErrorKindInvalidProviderResponse || !result.Failure.Retryable {
		t.Fatalf("unexpected failure: %#v", result.Failure)
	}
	if result.TranslatedCandidate != nil || result.ProtectionValidationTarget != nil {
		t.Fatalf("expected no translated candidates on invalid response, got %#v", result)
	}
}

type stubBodyTranslationProviderClient struct {
	response              stubBodyTranslationClientResponse
	err                   error
	capturedExecutionMode string
	capturedPrompt        string
}

func (stub *stubBodyTranslationProviderClient) GenerateBodyTranslation(
	_ context.Context,
	_ string,
	_ string,
	executionMode string,
	_ string,
	_ string,
	prompt string,
) (stubBodyTranslationClientResponse, error) {
	stub.capturedExecutionMode = executionMode
	stub.capturedPrompt = prompt
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
