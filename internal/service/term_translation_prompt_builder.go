package service

import (
	"fmt"
	"strings"
)

// TermTranslationPromptInput fixes one source term as one provider request unit.
type TermTranslationPromptInput struct {
	RequestUnitID  string
	SourceTerm     string
	SourceLanguage string
	TargetLanguage string
	RequestShapeID string
}

// TermTranslationPromptBuilder builds the provider prompt for one term translation unit.
type TermTranslationPromptBuilder interface {
	Build(input TermTranslationPromptInput) (PromptEnvelope, error)
}

type defaultTermTranslationPromptBuilder struct{}

// NewTermTranslationPromptBuilder creates the default one-term prompt builder.
func NewTermTranslationPromptBuilder() TermTranslationPromptBuilder {
	return defaultTermTranslationPromptBuilder{}
}

func (builder defaultTermTranslationPromptBuilder) Build(input TermTranslationPromptInput) (PromptEnvelope, error) {
	normalized, err := normalizeTermTranslationPromptInput(input)
	if err != nil {
		return PromptEnvelope{}, err
	}
	prompt := strings.TrimSpace(strings.Join([]string{
		normalized.RequestShapeID,
		"Return strict JSON only.",
		`Use the exact shape {"translations":[{"source_term":"...","translated_term":"..."}]}.`,
		"Do not add markdown, commentary, or extra keys.",
		"input_count=1",
		"execution_mode=" + TermTranslationExecutionModeSingleRequest,
		"request_unit_id=" + normalized.RequestUnitID,
		"source_language=" + normalized.SourceLanguage,
		"target_language=" + normalized.TargetLanguage,
		"source_term=" + normalized.SourceTerm,
	}, "\n"))
	return NewPromptEnvelope(prompt, normalized.RequestShapeID, PromptSafeSummary{
		InputCount:     1,
		ExecutionMode:  TermTranslationExecutionModeSingleRequest,
		CorrelationIDs: []string{normalized.RequestUnitID},
	})
}

// BuildTermTranslationPrompt returns the strict JSON-only prompt for one source term request unit.
func BuildTermTranslationPrompt(request TermTranslationProviderRequest) (string, error) {
	envelope, err := BuildTermTranslationPromptEnvelope(request)
	if err != nil {
		return "", err
	}
	return envelope.RawPrompt, nil
}

// BuildTermTranslationPromptEnvelope returns the internal prompt handoff unit for one source term.
func BuildTermTranslationPromptEnvelope(request TermTranslationProviderRequest) (PromptEnvelope, error) {
	envelope, err := NewTermTranslationPromptBuilder().Build(TermTranslationPromptInput{
		RequestUnitID:  request.RequestUnitID,
		SourceTerm:     request.SourceTerm,
		SourceLanguage: request.SourceLanguage,
		TargetLanguage: request.TargetLanguage,
		RequestShapeID: request.RequestShapeID,
	})
	if err != nil {
		return PromptEnvelope{}, fmt.Errorf("build term translation prompt envelope: %w", err)
	}
	return envelope, nil
}

func normalizeTermTranslationPromptInput(input TermTranslationPromptInput) (TermTranslationPromptInput, error) {
	sourceTerm := strings.TrimSpace(input.SourceTerm)
	if sourceTerm == "" {
		return TermTranslationPromptInput{}, fmt.Errorf("term translation source term is required")
	}
	requestUnitID := strings.TrimSpace(input.RequestUnitID)
	if requestUnitID == "" {
		requestUnitID = sourceTerm
	}
	sourceLanguage := strings.TrimSpace(input.SourceLanguage)
	if sourceLanguage == "" {
		sourceLanguage = "source"
	}
	targetLanguage := strings.TrimSpace(input.TargetLanguage)
	if targetLanguage == "" {
		targetLanguage = "target"
	}
	requestShapeID := firstNonEmptyTermTranslationValue(input.RequestShapeID, TermTranslationRequestShapeV1)
	return TermTranslationPromptInput{
		RequestUnitID:  requestUnitID,
		SourceTerm:     sourceTerm,
		SourceLanguage: sourceLanguage,
		TargetLanguage: targetLanguage,
		RequestShapeID: requestShapeID,
	}, nil
}
