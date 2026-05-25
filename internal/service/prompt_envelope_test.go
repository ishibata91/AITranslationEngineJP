package service

import (
	"strings"
	"testing"
)

func TestPromptEnvelopeSeparatesRawPromptDigestShapeAndSafeSummary(t *testing.T) {
	envelope, err := NewPromptEnvelope(" secret raw prompt ", TermTranslationRequestShapeV1, PromptSafeSummary{
		InputCount:     1,
		ExecutionMode:  "SINGLE_REQUEST",
		CorrelationIDs: []string{" term-1 ", ""},
		Counts: map[string]int{
			"outputs": 0,
			"bad":     -1,
		},
	})
	if err != nil {
		t.Fatalf("expected prompt envelope: %v", err)
	}

	if envelope.RawPrompt != "secret raw prompt" {
		t.Fatalf("expected raw prompt to stay internal, got %q", envelope.RawPrompt)
	}
	if envelope.Digest == "" || strings.Contains(string(envelope.Digest), "secret raw prompt") {
		t.Fatalf("expected non-reversible digest, got %q", envelope.Digest)
	}
	if envelope.RequestShapeID != TermTranslationRequestShapeV1 {
		t.Fatalf("expected request shape id, got %q", envelope.RequestShapeID)
	}
	if envelope.Summary.PromptDigest != envelope.Digest || envelope.Summary.RequestShapeID != envelope.RequestShapeID {
		t.Fatalf("expected summary identity fields to match envelope: %#v", envelope.Summary)
	}
	if envelope.Summary.ExecutionMode != TermTranslationExecutionModeSingleRequest {
		t.Fatalf("expected normalized execution mode, got %q", envelope.Summary.ExecutionMode)
	}
	if len(envelope.Summary.CorrelationIDs) != 1 || envelope.Summary.CorrelationIDs[0] != "term-1" {
		t.Fatalf("expected normalized correlation ids, got %#v", envelope.Summary.CorrelationIDs)
	}
	if _, ok := envelope.Summary.Counts["bad"]; ok {
		t.Fatalf("expected invalid count to be removed, got %#v", envelope.Summary.Counts)
	}
}

func TestRedactedPromptDiagnosticDoesNotExposeProtectedValue(t *testing.T) {
	diagnostic := RedactedPromptDiagnostic("request", `{"api_key":"secret","prompt":"raw"}`)
	if diagnostic == "" {
		t.Fatal("expected redacted diagnostic")
	}
	if strings.Contains(diagnostic, "secret") || strings.Contains(diagnostic, "raw") {
		t.Fatalf("expected protected value to stay hidden, got %q", diagnostic)
	}
	if !strings.HasPrefix(diagnostic, "sha256:request:") {
		t.Fatalf("expected labeled digest diagnostic, got %q", diagnostic)
	}
}

func TestPhasePromptEnvelopesUseRequestShapeIdentifiers(t *testing.T) {
	termEnvelope, err := BuildTermTranslationPromptEnvelope(TermTranslationProviderRequest{
		SourceTerm: "Dragonborn",
	})
	if err != nil {
		t.Fatalf("expected term prompt envelope: %v", err)
	}
	if termEnvelope.RequestShapeID != TermTranslationRequestShapeV1 {
		t.Fatalf("expected term request shape id, got %q", termEnvelope.RequestShapeID)
	}

	personaEnvelope, err := BuildPersonaGenerationPromptEnvelope(PersonaGenerationProviderRequest{
		RequestUnitID:    "unit-1",
		NPCCorrelationID: "npc-1",
		NPCDisplayName:   "Lydia",
	})
	if err != nil {
		t.Fatalf("expected persona prompt envelope: %v", err)
	}
	if personaEnvelope.RequestShapeID != PersonaGenerationRequestShapeV1 {
		t.Fatalf("expected persona request shape id, got %q", personaEnvelope.RequestShapeID)
	}

	bodyEnvelope, err := BuildBodyTranslationPromptEnvelope(BodyTranslationProviderRequest{
		RequestUnitID:       "unit-1",
		FieldCorrelationKey: "field-1",
		RecordType:          "INFO",
		FieldType:           "TEXT",
		SourceText:          "Hello",
	})
	if err != nil {
		t.Fatalf("expected body prompt envelope: %v", err)
	}
	if bodyEnvelope.RequestShapeID != BodyTranslationRequestShapeV1 {
		t.Fatalf("expected body request shape id, got %q", bodyEnvelope.RequestShapeID)
	}
}
