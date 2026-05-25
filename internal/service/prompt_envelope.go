package service

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"strings"
)

const (
	promptDigestScheme = "sha256:"

	// TermTranslationRequestShapeV1 identifies the provider request shape for one term translation unit.
	TermTranslationRequestShapeV1 = "TERM_TRANSLATION_REQUEST_V1"
	// PersonaGenerationRequestShapeV1 identifies the provider request shape for one NPC persona unit.
	PersonaGenerationRequestShapeV1 = "PERSONA_GENERATION_REQUEST_V1"
	// BodyTranslationRequestShapeV1 identifies the provider request shape for one body translation unit.
	BodyTranslationRequestShapeV1 = "BODY_TRANSLATION_REQUEST_V1"
)

// PromptDigest is a one-way identity value for prompt diagnostics.
type PromptDigest string

// PromptEnvelope is the internal handoff unit from a prompt builder to a provider adapter.
type PromptEnvelope struct {
	RawPrompt      string
	Digest         PromptDigest
	RequestShapeID string
	Summary        PromptSafeSummary
}

// PromptSafeSummary exposes provider request metadata without secrets, raw prompts, or raw provider payloads.
type PromptSafeSummary struct {
	RequestShapeID string
	PromptDigest   PromptDigest
	InputCount     int
	ExecutionMode  string
	CorrelationIDs []string
	Counts         map[string]int
}

// NewPromptEnvelope creates an internal provider prompt handoff unit.
func NewPromptEnvelope(rawPrompt string, requestShapeID string, summary PromptSafeSummary) (PromptEnvelope, error) {
	prompt := strings.TrimSpace(rawPrompt)
	if prompt == "" {
		return PromptEnvelope{}, fmt.Errorf("prompt raw text is required")
	}
	shapeID := strings.TrimSpace(requestShapeID)
	if shapeID == "" {
		return PromptEnvelope{}, fmt.Errorf("prompt request shape id is required")
	}
	safeSummary := normalizePromptSafeSummary(summary, shapeID, PromptDigestForRawPrompt(prompt))
	return PromptEnvelope{
		RawPrompt:      prompt,
		Digest:         safeSummary.PromptDigest,
		RequestShapeID: shapeID,
		Summary:        safeSummary,
	}, nil
}

// PromptDigestForRawPrompt returns a non-reversible digest for a raw prompt.
func PromptDigestForRawPrompt(rawPrompt string) PromptDigest {
	digest := sha256.Sum256([]byte(strings.TrimSpace(rawPrompt)))
	return PromptDigest(promptDigestScheme + hex.EncodeToString(digest[:]))
}

// PromptDigestString returns the digest string for boundaries that still store string values.
func PromptDigestString(rawPrompt string) string {
	return string(PromptDigestForRawPrompt(rawPrompt))
}

// RedactedPromptDiagnostic returns a diagnostic marker that cannot reveal the protected value.
func RedactedPromptDiagnostic(label string, value string) string {
	trimmed := strings.TrimSpace(value)
	if trimmed == "" {
		return ""
	}
	normalizedLabel := strings.Trim(strings.ToLower(strings.TrimSpace(label)), ":")
	if normalizedLabel == "" {
		normalizedLabel = "value"
	}
	return promptDigestScheme + normalizedLabel + ":" + strings.TrimPrefix(string(PromptDigestForRawPrompt(trimmed)), promptDigestScheme)
}

func normalizePromptSafeSummary(
	summary PromptSafeSummary,
	requestShapeID string,
	promptDigest PromptDigest,
) PromptSafeSummary {
	summary.RequestShapeID = strings.TrimSpace(requestShapeID)
	summary.PromptDigest = promptDigest
	summary.ExecutionMode = strings.ToLower(strings.TrimSpace(summary.ExecutionMode))
	summary.CorrelationIDs = normalizePromptSummaryValues(summary.CorrelationIDs)
	if summary.InputCount < 0 {
		summary.InputCount = 0
	}
	if summary.Counts != nil {
		counts := make(map[string]int, len(summary.Counts))
		for key, value := range summary.Counts {
			trimmedKey := strings.TrimSpace(key)
			if trimmedKey == "" || value < 0 {
				continue
			}
			counts[trimmedKey] = value
		}
		summary.Counts = counts
	}
	return summary
}

func normalizePromptSummaryValues(values []string) []string {
	normalized := make([]string, 0, len(values))
	for _, value := range values {
		trimmed := strings.TrimSpace(value)
		if trimmed == "" {
			continue
		}
		normalized = append(normalized, trimmed)
	}
	return normalized
}
