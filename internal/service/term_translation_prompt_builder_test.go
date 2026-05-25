package service

import (
	"strings"
	"testing"
)

func TestTermTranslationPromptBuilderFixesOneTermRequestUnit(t *testing.T) {
	envelope, err := NewTermTranslationPromptBuilder().Build(TermTranslationPromptInput{
		RequestUnitID:  "INFO:Dragonborn",
		SourceTerm:     "Dragonborn",
		SourceLanguage: "English",
		TargetLanguage: "Japanese",
	})
	if err != nil {
		t.Fatalf("expected prompt envelope: %v", err)
	}

	for _, expected := range []string{
		TermTranslationRequestShapeV1,
		"input_count=1",
		"execution_mode=" + TermTranslationExecutionModeSingleRequest,
		"request_unit_id=INFO:Dragonborn",
		"source_language=English",
		"target_language=Japanese",
		"source_term=Dragonborn",
	} {
		if !strings.Contains(envelope.RawPrompt, expected) {
			t.Fatalf("expected prompt to contain %q, got %q", expected, envelope.RawPrompt)
		}
	}
	if envelope.RequestShapeID != TermTranslationRequestShapeV1 {
		t.Fatalf("expected request shape id, got %q", envelope.RequestShapeID)
	}
	if len(envelope.Summary.CorrelationIDs) != 1 || envelope.Summary.CorrelationIDs[0] != "INFO:Dragonborn" {
		t.Fatalf("expected request unit correlation id, got %#v", envelope.Summary.CorrelationIDs)
	}
	if envelope.Summary.InputCount != 1 || envelope.Summary.ExecutionMode != TermTranslationExecutionModeSingleRequest {
		t.Fatalf("expected one-term safe summary, got %#v", envelope.Summary)
	}
}

func TestTermTranslationPromptBuilderRequiresSourceTerm(t *testing.T) {
	_, err := NewTermTranslationPromptBuilder().Build(TermTranslationPromptInput{
		RequestUnitID: "unit-1",
	})
	if err == nil {
		t.Fatal("expected source term validation error")
	}
}
