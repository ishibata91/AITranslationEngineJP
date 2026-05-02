package service

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"sort"
	"strings"
)

const (
	// BodyTranslationExecutionModeSingleRequest identifies the one-field-per-request mode.
	BodyTranslationExecutionModeSingleRequest = "single_request"

	bodyTranslationPromptFallbackNone         = "- none"
	bodyTranslationPromptVersionV1            = "BODY_TRANSLATION_REQUEST_V1"
	bodyTranslationTargetLanguageDefaultValue = "Japanese"
	bodyTranslationSourceLanguageDefaultValue = "source"
)

// BodyTranslationDictionaryExactMatchExclusion summarizes one exact-match exclusion kept out of provider execution.
type BodyTranslationDictionaryExactMatchExclusion struct {
	SourceText     string
	TranslatedText string
}

// BodyTranslationPartialMatchConstraint summarizes one partial-match fixed translation constraint.
type BodyTranslationPartialMatchConstraint struct {
	SourceText          string
	RequiredTranslation string
}

// BodyTranslationProtectedElement summarizes one protected source element that must survive translation.
type BodyTranslationProtectedElement struct {
	ElementType string
	SourceText  string
	Digest      string
}

// BodyTranslationProviderRequestSummary exposes provider request metadata without secrets or raw prompts.
type BodyTranslationProviderRequestSummary struct {
	RequestUnitID               string
	FieldCorrelationKey         string
	RecordType                  string
	FieldType                   string
	ProtectionSourceDigest      string
	ProtectedElementDigests     []string
	CompleteMatchExclusions     []BodyTranslationDictionaryExactMatchExclusion
	PartialMatchConstraints     []BodyTranslationPartialMatchConstraint
	ProtectedElementCount       int
	CompleteMatchExclusionCount int
	PartialMatchConstraintCount int
}

// BuildBodyTranslationPrompt returns the strict JSON-only prompt for one body translation request unit.
func BuildBodyTranslationPrompt(request BodyTranslationProviderRequest) (string, error) {
	requestUnitID := strings.TrimSpace(request.RequestUnitID)
	if requestUnitID == "" {
		return "", fmt.Errorf("body translation request unit id is required")
	}
	fieldCorrelationKey := strings.TrimSpace(request.FieldCorrelationKey)
	if fieldCorrelationKey == "" {
		return "", fmt.Errorf("body translation field correlation key is required")
	}
	recordType := strings.TrimSpace(request.RecordType)
	if recordType == "" {
		return "", fmt.Errorf("body translation record type is required")
	}
	fieldType := strings.TrimSpace(request.FieldType)
	if fieldType == "" {
		return "", fmt.Errorf("body translation field type is required")
	}
	sourceText := strings.TrimSpace(request.SourceText)
	if sourceText == "" {
		return "", fmt.Errorf("body translation source text is required")
	}

	executionMode := strings.ToLower(strings.TrimSpace(request.ExecutionMode))
	if executionMode == "" {
		executionMode = BodyTranslationExecutionModeSingleRequest
	}
	if executionMode != BodyTranslationExecutionModeSingleRequest {
		return "", fmt.Errorf("unsupported body translation execution mode: %s", request.ExecutionMode)
	}

	sourceLanguage := strings.TrimSpace(request.SourceLanguage)
	if sourceLanguage == "" {
		sourceLanguage = bodyTranslationSourceLanguageDefaultValue
	}
	targetLanguage := strings.TrimSpace(request.TargetLanguage)
	if targetLanguage == "" {
		targetLanguage = bodyTranslationTargetLanguageDefaultValue
	}

	personaSummary := strings.TrimSpace(request.PersonaSummary)
	if personaSummary == "" {
		personaSummary = "none"
	}

	completeMatchExclusions, err := normalizeBodyTranslationExactMatchExclusions(request.CompleteMatchExclusions)
	if err != nil {
		return "", err
	}
	partialMatchConstraints, err := normalizeBodyTranslationPartialMatchConstraints(request.PartialMatchConstraints)
	if err != nil {
		return "", err
	}
	protectedElements, err := normalizeBodyTranslationProtectedElements(request.ProtectedElements)
	if err != nil {
		return "", err
	}

	contextLines := normalizeBodyTranslationPromptLines(request.ContextLines, bodyTranslationPromptFallbackNone)
	instructionLines := buildBodyTranslationInstructionLines(recordType, fieldType)

	return strings.TrimSpace(strings.Join([]string{
		bodyTranslationPromptVersionV1,
		"Return strict JSON only.",
		`Use the exact shape {"translations":[{"request_unit_id":"...","field_correlation_key":"...","translated_text":"..."}]}.`,
		"Do not add markdown, commentary, or extra keys.",
		"input_count=1",
		"execution_mode=" + executionMode,
		"request_unit_id=" + requestUnitID,
		"field_correlation_key=" + fieldCorrelationKey,
		"record_type=" + recordType,
		"field_type=" + fieldType,
		"source_language=" + sourceLanguage,
		"target_language=" + targetLanguage,
		"protection_source_digest=" + buildBodyTranslationProtectionSourceDigest(protectedElements),
		"instructions:",
		strings.Join(instructionLines, "\n"),
		"persona_summary=" + personaSummary,
		"context_lines:",
		strings.Join(contextLines, "\n"),
		"complete_match_exclusions:",
		strings.Join(renderBodyTranslationExactMatchExclusions(completeMatchExclusions), "\n"),
		"partial_match_constraints:",
		strings.Join(renderBodyTranslationPartialMatchConstraints(partialMatchConstraints), "\n"),
		"protected_elements:",
		strings.Join(renderBodyTranslationProtectedElements(protectedElements), "\n"),
		"source_text=" + sourceText,
	}, "\n")), nil
}

func bodyTranslationPromptDigest(prompt string) string {
	digest := sha256.Sum256([]byte(strings.TrimSpace(prompt)))
	return "sha256:" + hex.EncodeToString(digest[:])
}

func buildBodyTranslationRequestSummary(request BodyTranslationProviderRequest) BodyTranslationProviderRequestSummary {
	exclusions, _ := normalizeBodyTranslationExactMatchExclusions(request.CompleteMatchExclusions)
	constraints, _ := normalizeBodyTranslationPartialMatchConstraints(request.PartialMatchConstraints)
	protectedElements, _ := normalizeBodyTranslationProtectedElements(request.ProtectedElements)
	digests := make([]string, 0, len(protectedElements))
	for _, element := range protectedElements {
		digests = append(digests, element.Digest)
	}

	return BodyTranslationProviderRequestSummary{
		RequestUnitID:               strings.TrimSpace(request.RequestUnitID),
		FieldCorrelationKey:         strings.TrimSpace(request.FieldCorrelationKey),
		RecordType:                  strings.TrimSpace(request.RecordType),
		FieldType:                   strings.TrimSpace(request.FieldType),
		ProtectionSourceDigest:      buildBodyTranslationProtectionSourceDigest(protectedElements),
		ProtectedElementDigests:     digests,
		CompleteMatchExclusions:     exclusions,
		PartialMatchConstraints:     constraints,
		ProtectedElementCount:       len(protectedElements),
		CompleteMatchExclusionCount: len(exclusions),
		PartialMatchConstraintCount: len(constraints),
	}
}

func normalizeBodyTranslationPromptLines(lines []string, fallback string) []string {
	normalized := make([]string, 0, len(lines))
	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if trimmed == "" {
			continue
		}
		normalized = append(normalized, "- "+trimmed)
	}
	if len(normalized) == 0 {
		return []string{fallback}
	}
	return normalized
}

func buildBodyTranslationInstructionLines(recordType string, fieldType string) []string {
	normalizedRecordType := strings.ToUpper(strings.TrimSpace(recordType))
	normalizedFieldType := strings.ToUpper(strings.TrimSpace(fieldType))

	lines := []string{
		"- Preserve all protected elements exactly as provided.",
		"- Respect partial match constraints exactly.",
		"- Keep the translation natural for Skyrim mod text.",
	}

	switch {
	case normalizedFieldType == "FULL" || strings.Contains(normalizedFieldType, "NAME"):
		lines = append(lines, "- Treat the field as a visible label or proper name.")
	case strings.Contains(normalizedFieldType, "DESC") || strings.Contains(normalizedFieldType, "TEXT"):
		lines = append(lines, "- Treat the field as descriptive prose and preserve sentence structure.")
	case normalizedRecordType == "INFO" || normalizedRecordType == "DIAL":
		lines = append(lines, "- Treat the field as spoken dialogue and keep the voice natural.")
	default:
		lines = append(lines, "- Preserve the intent implied by the record type and field type.")
	}

	return lines
}

func normalizeBodyTranslationExactMatchExclusions(
	exclusions []BodyTranslationDictionaryExactMatchExclusion,
) ([]BodyTranslationDictionaryExactMatchExclusion, error) {
	normalized := make([]BodyTranslationDictionaryExactMatchExclusion, 0, len(exclusions))
	for _, exclusion := range exclusions {
		sourceText := strings.TrimSpace(exclusion.SourceText)
		translatedText := strings.TrimSpace(exclusion.TranslatedText)
		if sourceText == "" || translatedText == "" {
			return nil, fmt.Errorf("body translation exact match exclusion must include source text and translated text")
		}
		normalized = append(normalized, BodyTranslationDictionaryExactMatchExclusion{
			SourceText:     sourceText,
			TranslatedText: translatedText,
		})
	}
	sort.SliceStable(normalized, func(left int, right int) bool {
		if normalized[left].SourceText == normalized[right].SourceText {
			return normalized[left].TranslatedText < normalized[right].TranslatedText
		}
		return normalized[left].SourceText < normalized[right].SourceText
	})
	return normalized, nil
}

func normalizeBodyTranslationPartialMatchConstraints(
	constraints []BodyTranslationPartialMatchConstraint,
) ([]BodyTranslationPartialMatchConstraint, error) {
	normalized := make([]BodyTranslationPartialMatchConstraint, 0, len(constraints))
	for _, constraint := range constraints {
		sourceText := strings.TrimSpace(constraint.SourceText)
		requiredTranslation := strings.TrimSpace(constraint.RequiredTranslation)
		if sourceText == "" || requiredTranslation == "" {
			return nil, fmt.Errorf("body translation partial match constraint must include source text and required translation")
		}
		normalized = append(normalized, BodyTranslationPartialMatchConstraint{
			SourceText:          sourceText,
			RequiredTranslation: requiredTranslation,
		})
	}
	sort.SliceStable(normalized, func(left int, right int) bool {
		if normalized[left].SourceText == normalized[right].SourceText {
			return normalized[left].RequiredTranslation < normalized[right].RequiredTranslation
		}
		return normalized[left].SourceText < normalized[right].SourceText
	})
	return normalized, nil
}

func normalizeBodyTranslationProtectedElements(
	elements []BodyTranslationProtectedElement,
) ([]BodyTranslationProtectedElement, error) {
	normalized := make([]BodyTranslationProtectedElement, 0, len(elements))
	for _, element := range elements {
		elementType := strings.TrimSpace(element.ElementType)
		sourceText := strings.TrimSpace(element.SourceText)
		digest := strings.TrimSpace(element.Digest)
		if elementType == "" || sourceText == "" || digest == "" {
			return nil, fmt.Errorf("body translation protected element must include element type, source text, and digest")
		}
		normalized = append(normalized, BodyTranslationProtectedElement{
			ElementType: elementType,
			SourceText:  sourceText,
			Digest:      digest,
		})
	}
	return normalized, nil
}

func renderBodyTranslationExactMatchExclusions(
	exclusions []BodyTranslationDictionaryExactMatchExclusion,
) []string {
	if len(exclusions) == 0 {
		return []string{bodyTranslationPromptFallbackNone}
	}
	lines := make([]string, 0, len(exclusions))
	for _, exclusion := range exclusions {
		lines = append(lines, "- "+exclusion.SourceText+" => "+exclusion.TranslatedText)
	}
	return lines
}

func renderBodyTranslationPartialMatchConstraints(
	constraints []BodyTranslationPartialMatchConstraint,
) []string {
	if len(constraints) == 0 {
		return []string{bodyTranslationPromptFallbackNone}
	}
	lines := make([]string, 0, len(constraints))
	for _, constraint := range constraints {
		lines = append(lines, "- "+constraint.SourceText+" => "+constraint.RequiredTranslation)
	}
	return lines
}

func renderBodyTranslationProtectedElements(elements []BodyTranslationProtectedElement) []string {
	if len(elements) == 0 {
		return []string{bodyTranslationPromptFallbackNone}
	}
	lines := make([]string, 0, len(elements))
	for _, element := range elements {
		lines = append(lines, "- type="+element.ElementType+"; digest="+element.Digest+"; source="+element.SourceText)
	}
	return lines
}

func buildBodyTranslationProtectionSourceDigest(elements []BodyTranslationProtectedElement) string {
	if len(elements) == 0 {
		return ""
	}
	if len(elements) == 1 {
		return strings.TrimSpace(elements[0].Digest)
	}
	parts := make([]string, 0, len(elements))
	for _, element := range elements {
		digest := strings.TrimSpace(element.Digest)
		if digest == "" {
			continue
		}
		parts = append(parts, digest)
	}
	if len(parts) == 0 {
		return ""
	}
	digest := sha256.Sum256([]byte(strings.Join(parts, "\n")))
	return "sha256:" + hex.EncodeToString(digest[:])
}
