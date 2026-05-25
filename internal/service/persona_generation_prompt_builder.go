package service

import (
	"fmt"
	"strings"
)

const personaGenerationPromptFallbackNone = "- none"

// PersonaGenerationPromptInput contains only the values needed to build one NPC persona prompt.
type PersonaGenerationPromptInput struct {
	RequestUnitID            string
	NPCCorrelationID         string
	NPCDisplayName           string
	NPCEditorID              string
	NPCFormID                string
	NPCAttributes            []string
	ConversationContext      []string
	CommonPersonaSummary     string
	RecentOriginalUtterances []string
}

// PersonaGenerationPromptBuilder builds one NPC persona prompt envelope.
type PersonaGenerationPromptBuilder interface {
	Build(input PersonaGenerationPromptInput) (PromptEnvelope, error)
}

// DefaultPersonaGenerationPromptBuilder is the default one-NPC persona prompt builder.
type DefaultPersonaGenerationPromptBuilder struct{}

// NewPersonaGenerationPromptBuilder returns the default NPC persona prompt builder.
func NewPersonaGenerationPromptBuilder() PersonaGenerationPromptBuilder {
	return DefaultPersonaGenerationPromptBuilder{}
}

// Build returns the internal prompt handoff unit for one NPC.
func (builder DefaultPersonaGenerationPromptBuilder) Build(input PersonaGenerationPromptInput) (PromptEnvelope, error) {
	prompt, err := buildPersonaGenerationPromptText(input)
	if err != nil {
		return PromptEnvelope{}, err
	}
	return NewPromptEnvelope(prompt, PersonaGenerationRequestShapeV1, PromptSafeSummary{
		InputCount:     1,
		ExecutionMode:  PersonaGenerationExecutionModeSingleRequest,
		CorrelationIDs: []string{input.RequestUnitID, input.NPCCorrelationID},
		Counts: map[string]int{
			"npc_attributes":              len(input.NPCAttributes),
			"conversation_context":        len(input.ConversationContext),
			"recent_original_utterances":  len(input.RecentOriginalUtterances),
			"common_persona_summary_unit": personaGenerationNonEmptyCount(input.CommonPersonaSummary),
		},
	})
}

// BuildPersonaGenerationPrompt returns the strict JSON-only prompt for one NPC request unit.
func BuildPersonaGenerationPrompt(request PersonaGenerationProviderRequest) (string, error) {
	envelope, err := BuildPersonaGenerationPromptEnvelope(request)
	if err != nil {
		return "", err
	}
	return envelope.RawPrompt, nil
}

// BuildPersonaGenerationPromptEnvelope adapts provider request data to the dedicated prompt input.
func BuildPersonaGenerationPromptEnvelope(request PersonaGenerationProviderRequest) (PromptEnvelope, error) {
	envelope, err := NewPersonaGenerationPromptBuilder().Build(personaGenerationPromptInputFromProviderRequest(request))
	if err != nil {
		return PromptEnvelope{}, fmt.Errorf("build persona generation prompt envelope: %w", err)
	}
	return envelope, nil
}

func buildPersonaGenerationPromptText(input PersonaGenerationPromptInput) (string, error) {
	requestUnitID := strings.TrimSpace(input.RequestUnitID)
	if requestUnitID == "" {
		return "", fmt.Errorf("persona generation request unit id is required")
	}
	npcCorrelationID := strings.TrimSpace(input.NPCCorrelationID)
	if npcCorrelationID == "" {
		return "", fmt.Errorf("persona generation npc correlation id is required")
	}
	displayName := strings.TrimSpace(input.NPCDisplayName)
	if displayName == "" {
		displayName = "unknown"
	}
	editorID := strings.TrimSpace(input.NPCEditorID)
	if editorID == "" {
		editorID = "unknown"
	}
	formID := strings.TrimSpace(input.NPCFormID)
	if formID == "" {
		formID = "unknown"
	}
	attributes := normalizePersonaGenerationPromptLines(input.NPCAttributes, personaGenerationPromptFallbackNone)
	conversationContext := normalizePersonaGenerationPromptLines(input.ConversationContext, personaGenerationPromptFallbackNone)
	recentUtterances := normalizePersonaGenerationPromptLines(input.RecentOriginalUtterances, personaGenerationPromptFallbackNone)
	commonPersonaSummary := strings.TrimSpace(input.CommonPersonaSummary)
	if commonPersonaSummary == "" {
		commonPersonaSummary = "none"
	}

	return strings.TrimSpace(strings.Join([]string{
		PersonaGenerationRequestShapeV1,
		"Return strict JSON only.",
		`Use the exact shape {"personas":[{"request_unit_id":"...","npc_correlation_id":"...","persona_body":"..."}]}.`,
		"Do not add markdown, commentary, or extra keys.",
		"input_count=1",
		"execution_mode=" + PersonaGenerationExecutionModeSingleRequest,
		"request_unit_id=" + requestUnitID,
		"npc_correlation_id=" + npcCorrelationID,
		"npc_display_name=" + displayName,
		"npc_editor_id=" + editorID,
		"npc_form_id=" + formID,
		"common_persona_summary=" + commonPersonaSummary,
		"npc_attributes:",
		strings.Join(attributes, "\n"),
		"conversation_context:",
		strings.Join(conversationContext, "\n"),
		"recent_original_utterances:",
		strings.Join(recentUtterances, "\n"),
	}, "\n")), nil
}

func personaGenerationPromptInputFromProviderRequest(request PersonaGenerationProviderRequest) PersonaGenerationPromptInput {
	return PersonaGenerationPromptInput{
		RequestUnitID:            request.RequestUnitID,
		NPCCorrelationID:         request.NPCCorrelationID,
		NPCDisplayName:           request.NPCDisplayName,
		NPCEditorID:              request.NPCEditorID,
		NPCFormID:                request.NPCFormID,
		NPCAttributes:            append([]string(nil), request.NPCAttributes...),
		ConversationContext:      append([]string(nil), request.ConversationContext...),
		CommonPersonaSummary:     request.CommonPersonaSummary,
		RecentOriginalUtterances: append([]string(nil), request.RecentOriginalUtterances...),
	}
}

func personaGenerationPromptDigest(prompt string) string {
	return PromptDigestString(prompt)
}

func normalizePersonaGenerationPromptLines(lines []string, fallback string) []string {
	normalized := make([]string, 0, len(lines))
	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if trimmed == "" {
			continue
		}
		normalized = append(normalized, "- "+trimmed)
	}
	if len(normalized) == 0 {
		normalized = append(normalized, fallback)
	}
	return normalized
}

func personaGenerationNonEmptyCount(value string) int {
	if strings.TrimSpace(value) == "" {
		return 0
	}
	return 1
}
