package service

import (
	"context"
	"errors"
	"strings"
	"testing"
)

func TestBuildPersonaGenerationPromptIncludesAttributesAndContext(t *testing.T) {
	prompt, err := BuildPersonaGenerationPrompt(PersonaGenerationProviderRequest{
		RequestUnitID:            "unit-1",
		NPCCorrelationID:         "npc-1",
		NPCDisplayName:           "Lydia",
		NPCEditorID:              "HousecarlWhiterun",
		NPCFormID:                "000A2C94",
		NPCAttributes:            []string{"housecarl", "loyal"},
		ConversationContext:      []string{"守る対象はドラゴンボーン。"},
		CommonPersonaSummary:     "使命感が強い。",
		RecentOriginalUtterances: []string{"I am sworn to carry your burdens."},
	})
	if err != nil {
		t.Fatalf("expected prompt builder success: %v", err)
	}
	if !strings.Contains(prompt, "npc_attributes:\n- housecarl\n- loyal") {
		t.Fatalf("expected prompt to include npc attributes: %s", prompt)
	}
	if !strings.Contains(prompt, "conversation_context:\n- 守る対象はドラゴンボーン。") {
		t.Fatalf("expected prompt to include conversation context: %s", prompt)
	}
	if !strings.Contains(prompt, "recent_original_utterances:\n- I am sworn to carry your burdens.") {
		t.Fatalf("expected prompt to include utterances: %s", prompt)
	}
}

func TestPersonaGenerationPromptBuilderBuildsOneNPCEnvelopeWithoutRawPromptSummary(t *testing.T) {
	builder := NewPersonaGenerationPromptBuilder()
	input := PersonaGenerationPromptInput{
		RequestUnitID:            "unit-1",
		NPCCorrelationID:         "npc-1",
		NPCDisplayName:           "Lydia",
		NPCEditorID:              "HousecarlWhiterun",
		NPCFormID:                "000A2C94",
		NPCAttributes:            []string{"housecarl", "loyal"},
		ConversationContext:      []string{"守る対象はドラゴンボーン。"},
		CommonPersonaSummary:     "使命感が強い。",
		RecentOriginalUtterances: []string{"I am sworn to carry your burdens."},
	}

	envelope, err := builder.Build(input)
	if err != nil {
		t.Fatalf("expected envelope builder success: %v", err)
	}

	if !strings.Contains(envelope.RawPrompt, "request_unit_id=unit-1") ||
		!strings.Contains(envelope.RawPrompt, "npc_correlation_id=npc-1") ||
		!strings.Contains(envelope.RawPrompt, "npc_display_name=Lydia") {
		t.Fatalf("expected prompt to include one NPC correlation fields: %s", envelope.RawPrompt)
	}
	if !strings.Contains(envelope.RawPrompt, "conversation_context:\n- 守る対象はドラゴンボーン。") ||
		!strings.Contains(envelope.RawPrompt, "recent_original_utterances:\n- I am sworn to carry your burdens.") {
		t.Fatalf("expected prompt to include protected source values internally: %s", envelope.RawPrompt)
	}
	if envelope.RequestShapeID != PersonaGenerationRequestShapeV1 {
		t.Fatalf("unexpected request shape id: %#v", envelope)
	}
	if envelope.Summary.InputCount != 1 ||
		envelope.Summary.ExecutionMode != PersonaGenerationExecutionModeSingleRequest {
		t.Fatalf("unexpected safe summary: %#v", envelope.Summary)
	}
	if strings.Contains(strings.Join(envelope.Summary.CorrelationIDs, "\n"), "I am sworn") ||
		strings.Contains(strings.Join(envelope.Summary.CorrelationIDs, "\n"), "守る対象") {
		t.Fatalf("safe summary must not include raw utterance or context: %#v", envelope.Summary)
	}
	if envelope.Summary.Counts["npc_attributes"] != 2 ||
		envelope.Summary.Counts["conversation_context"] != 1 ||
		envelope.Summary.Counts["recent_original_utterances"] != 1 ||
		envelope.Summary.Counts["common_persona_summary_unit"] != 1 {
		t.Fatalf("unexpected safe summary counts: %#v", envelope.Summary.Counts)
	}
}

func TestBuildPersonaGenerationPromptEnvelopeAdaptsProviderRequestToPromptInput(t *testing.T) {
	request := PersonaGenerationProviderRequest{
		RequestUnitID:            "unit-1",
		NPCCorrelationID:         "npc-1",
		NPCDisplayName:           "Lydia",
		NPCEditorID:              "HousecarlWhiterun",
		NPCFormID:                "000A2C94",
		NPCAttributes:            []string{"housecarl"},
		ConversationContext:      []string{"守る対象はドラゴンボーン。"},
		CommonPersonaSummary:     "使命感が強い。",
		RecentOriginalUtterances: []string{"I am sworn to carry your burdens."},
	}

	envelope, err := BuildPersonaGenerationPromptEnvelope(request)
	if err != nil {
		t.Fatalf("expected provider request wrapper success: %v", err)
	}

	if !strings.Contains(envelope.RawPrompt, "npc_editor_id=HousecarlWhiterun") ||
		!strings.Contains(envelope.RawPrompt, "common_persona_summary=使命感が強い。") {
		t.Fatalf("expected wrapper to preserve prompt input values: %s", envelope.RawPrompt)
	}
	if strings.Contains(strings.Join(envelope.Summary.CorrelationIDs, "\n"), "I am sworn") ||
		strings.Contains(strings.Join(envelope.Summary.CorrelationIDs, "\n"), "守る対象") {
		t.Fatalf("safe summary must not include raw utterance or context: %#v", envelope.Summary)
	}
}

func TestPersonaGenerationProviderAdapterMapsValidResponse(t *testing.T) {
	client := stubPersonaGenerationProviderClient{
		response: stubPersonaGenerationClientResponse{
			Items: []stubPersonaGenerationClientItem{{
				RequestUnitID:    "unit-1",
				NPCCorrelationID: "npc-1",
				PersonaBody:      "忠義を誇りにする衛士。",
			}},
			ExecutionMode: PersonaGenerationExecutionModeSingleRequest,
			PromptDigest:  "sha256:from-client",
			DebugLog: stubPersonaGenerationClientDebugLog{
				Prompt:         "PERSONA_GENERATION_REQUEST_V1",
				RequestBody:    "{\"messages\":[]}",
				Headers:        map[string]string{"Authorization": "[REDACTED]"},
				SecretRedacted: true,
			},
		},
	}
	adapter := NewPersonaGenerationProviderAdapter(client)

	result := adapter.GeneratePersona(context.Background(), PersonaGenerationProviderRequest{
		Provider:                 PersonaGenerationProviderGemini,
		Model:                    "gemini-model",
		ExecutionMode:            PersonaGenerationExecutionModeSingleRequest,
		CredentialRef:            "persona-ref",
		RequestUnitID:            "unit-1",
		NPCCorrelationID:         "npc-1",
		NPCDisplayName:           "Lydia",
		NPCAttributes:            []string{"housecarl"},
		ConversationContext:      []string{"守る対象はドラゴンボーン。"},
		RecentOriginalUtterances: []string{"I am sworn to carry your burdens."},
	})

	if result.Failure != nil {
		t.Fatalf("expected success, got failure %#v", result.Failure)
	}
	if result.PersonaBody != "忠義を誇りにする衛士。" {
		t.Fatalf("unexpected persona body: %q", result.PersonaBody)
	}
	if result.RequestUnitID != "unit-1" || result.NPCCorrelationID != "npc-1" {
		t.Fatalf("unexpected correlation mapping: %#v", result)
	}
	if result.AuditSummary.Provider != PersonaGenerationProviderGemini ||
		result.AuditSummary.Model != "gemini-model" ||
		result.AuditSummary.ExecutionMode != PersonaGenerationExecutionModeSingleRequest {
		t.Fatalf("unexpected audit summary: %#v", result.AuditSummary)
	}
	if result.AuditSummary.CredentialRef != "persona-ref" {
		t.Fatalf("expected credential ref passthrough, got %#v", result.AuditSummary)
	}
	if result.AuditSummary.RequestShapeID != PersonaGenerationRequestShapeV1 ||
		!strings.HasPrefix(result.AuditSummary.PromptDigest, "sha256:") {
		t.Fatalf("expected prompt digest and request shape id, got %#v", result.AuditSummary)
	}
	if !result.DebugLog.SecretRedacted || result.DebugLog.Headers["Authorization"] != "[REDACTED]" {
		t.Fatalf("expected redacted debug log, got %#v", result.DebugLog)
	}
	if strings.Contains(result.DebugLog.Prompt, PersonaGenerationRequestShapeV1) ||
		strings.Contains(result.DebugLog.RequestBody, "{\"messages\":[]}") {
		t.Fatalf("expected redacted debug log payloads, got %#v", result.DebugLog)
	}
}

func TestPersonaGenerationProviderAdapterRejectsMismatchedCorrelation(t *testing.T) {
	adapter := NewPersonaGenerationProviderAdapter(stubPersonaGenerationProviderClient{
		response: stubPersonaGenerationClientResponse{
			Items: []stubPersonaGenerationClientItem{{
				RequestUnitID:    "unit-1",
				NPCCorrelationID: "npc-other",
				PersonaBody:      "persona",
			}},
			ExecutionMode: PersonaGenerationExecutionModeSingleRequest,
			PromptDigest:  "sha256:from-client",
		},
	})

	result := adapter.GeneratePersona(context.Background(), PersonaGenerationProviderRequest{
		Provider:                 PersonaGenerationProviderGemini,
		Model:                    "gemini-model",
		ExecutionMode:            PersonaGenerationExecutionModeSingleRequest,
		RequestUnitID:            "unit-1",
		NPCCorrelationID:         "npc-1",
		NPCDisplayName:           "Lydia",
		ConversationContext:      []string{"ctx"},
		RecentOriginalUtterances: []string{"line"},
	})

	if result.Failure == nil {
		t.Fatal("expected invalid response failure")
	}
	if result.Failure.Kind != PersonaGenerationProviderErrorKindInvalidProviderResponse || !result.Failure.Retryable {
		t.Fatalf("unexpected failure: %#v", result.Failure)
	}
	if result.PersonaBody != "" {
		t.Fatalf("unexpected persona body on invalid response: %q", result.PersonaBody)
	}
}

func TestPersonaGenerationProviderAdapterRejectsEmptyPersonaBody(t *testing.T) {
	// 空ペルソナ本文は NPC 単位の invalid response として扱う。
	adapter := NewPersonaGenerationProviderAdapter(stubPersonaGenerationProviderClient{
		response: stubPersonaGenerationClientResponse{
			Items: []stubPersonaGenerationClientItem{{
				RequestUnitID:    "unit-1",
				NPCCorrelationID: "npc-1",
				PersonaBody:      " ",
			}},
			ExecutionMode: PersonaGenerationExecutionModeSingleRequest,
			PromptDigest:  "sha256:from-client",
		},
	})

	result := adapter.GeneratePersona(context.Background(), PersonaGenerationProviderRequest{
		Provider:                 PersonaGenerationProviderGemini,
		Model:                    "gemini-model",
		ExecutionMode:            PersonaGenerationExecutionModeSingleRequest,
		RequestUnitID:            "unit-1",
		NPCCorrelationID:         "npc-1",
		NPCDisplayName:           "Lydia",
		ConversationContext:      []string{"ctx"},
		RecentOriginalUtterances: []string{"line"},
	})

	if result.Failure == nil {
		t.Fatal("expected invalid response failure")
	}
	if result.Failure.Kind != PersonaGenerationProviderErrorKindInvalidProviderResponse || !result.Failure.Retryable {
		t.Fatalf("unexpected failure: %#v", result.Failure)
	}
	if result.PersonaBody != "" {
		t.Fatalf("unexpected persona body on invalid response: %q", result.PersonaBody)
	}
}

func TestPersonaGenerationProviderAdapterRedactsProviderCredentialFailure(t *testing.T) {
	adapter := NewPersonaGenerationProviderAdapter(stubPersonaGenerationProviderClient{
		err: stubPersonaGenerationClientError{
			kind:      "provider_failure",
			retryable: false,
			cause:     errors.New("provider credential resolver is required"),
		},
	})

	result := adapter.GeneratePersona(context.Background(), PersonaGenerationProviderRequest{
		Provider:                 PersonaGenerationProviderGemini,
		Model:                    "gemini-2.5-pro",
		ExecutionMode:            PersonaGenerationExecutionModeSingleRequest,
		RequestUnitID:            "unit-1",
		NPCCorrelationID:         "npc-1",
		NPCDisplayName:           "Lydia",
		ConversationContext:      []string{"ctx"},
		RecentOriginalUtterances: []string{"line"},
	})

	if result.Failure == nil {
		t.Fatal("expected provider failure")
	}
	if result.Failure.Reason != "provider credential is unavailable" {
		t.Fatalf("expected redacted credential failure, got %#v", result.Failure)
	}
	if !result.Failure.IsRedacted || result.Failure.Retryable {
		t.Fatalf("unexpected failure metadata: %#v", result.Failure)
	}
}

type stubPersonaGenerationProviderClient struct {
	response stubPersonaGenerationClientResponse
	err      error
	testSafe bool
}

func (stub stubPersonaGenerationProviderClient) GeneratePersona(
	_ context.Context,
	_ string,
	_ string,
	_ string,
	_ string,
	_ string,
) (stubPersonaGenerationClientResponse, error) {
	if stub.err != nil {
		return stubPersonaGenerationClientResponse{}, stub.err
	}
	return stub.response, nil
}

func (stub stubPersonaGenerationProviderClient) ProviderRequestsAreTestSafe() bool {
	return stub.testSafe
}

type stubPersonaGenerationClientResponse struct {
	Items         []stubPersonaGenerationClientItem
	ExecutionMode string
	PromptDigest  string
	DebugLog      stubPersonaGenerationClientDebugLog
}

type stubPersonaGenerationClientItem struct {
	RequestUnitID    string
	NPCCorrelationID string
	PersonaBody      string
}

type stubPersonaGenerationClientDebugLog struct {
	Prompt         string
	RequestBody    string
	Headers        map[string]string
	SecretRedacted bool
}

type stubPersonaGenerationClientError struct {
	kind      string
	retryable bool
	cause     error
}

func (err stubPersonaGenerationClientError) Error() string {
	if err.cause == nil {
		return ""
	}
	return err.cause.Error()
}

func (err stubPersonaGenerationClientError) Unwrap() error {
	return err.cause
}

func (err stubPersonaGenerationClientError) FailureKind() string {
	return err.kind
}

func (err stubPersonaGenerationClientError) FailureRetryable() bool {
	return err.retryable
}
