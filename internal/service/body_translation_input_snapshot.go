package service

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"sort"
	"strings"

	"aitranslationenginejp/internal/repository"
)

const (
	bodyTranslationDictionaryExactKind   = "exact"
	bodyTranslationDictionaryPartialKind = "partial"

	bodyTranslationSkippedReasonExactDictionary = "exact_dictionary_match"
)

type bodyTranslationSnapshotField struct {
	TranslationFieldID         int64
	RecordType                 string
	FieldType                  string
	SourceText                 string
	CompleteMatchExclusions    []BodyTranslationDictionaryExactMatchExclusion
	PartialMatchConstraints    []BodyTranslationPartialMatchConstraint
	IncludedInProviderRequests bool
}

type bodyTranslationInputSnapshot struct {
	TargetCount            int
	Fields                 []bodyTranslationSnapshotField
	SkippedReasons         []string
	DictionaryDigest       string
	PersonaDigest          string
	MetadataDigest         string
	PromptDigest           string
	InputSnapshotDigest    string
	ProviderTargetCount    int
	ExactExclusionCount    int
	PartialConstraintCount int
}

func buildBodyTranslationInputSnapshot(
	job repository.TranslationJob,
	records []repository.TranslationRecord,
	fieldsByRecordID map[int64][]repository.TranslationField,
	dictionaryEntries []repository.DictionaryEntry,
	persona repository.Persona,
	execution BodyTranslationPhaseExecutionSummaryReadModel,
) (bodyTranslationInputSnapshot, error) {
	snapshot := bodyTranslationInputSnapshot{
		Fields:         make([]bodyTranslationSnapshotField, 0),
		SkippedReasons: make([]string, 0),
	}

	for _, record := range records {
		for _, field := range sortedBodyTranslationRecordFields(fieldsByRecordID[record.ID]) {
			snapshotField := buildBodyTranslationSnapshotField(record, field, dictionaryEntries)
			applyBodyTranslationSnapshotField(&snapshot, snapshotField)
		}
	}

	dictionaryDigest := bodyTranslationDigestLines(buildBodyTranslationDictionaryDigestLines(dictionaryEntries))
	personaDigest := bodyTranslationDigestLines([]string{
		fmt.Sprintf("persona_id=%d", persona.ID),
		"persona_description=" + strings.TrimSpace(persona.PersonaDescription),
		"speech_style=" + strings.TrimSpace(persona.SpeechStyle),
		"personality_summary=" + strings.TrimSpace(persona.PersonalitySummary),
		fmt.Sprintf("evidence_utterance_count=%d", persona.EvidenceUtteranceCount),
	})
	metadataDigest := bodyTranslationDigestLines(buildBodyTranslationMetadataDigestLines(job, records, snapshot.Fields, execution))
	promptDigest, err := buildBodyTranslationPromptDigest(snapshot.Fields, persona, execution)
	if err != nil {
		return bodyTranslationInputSnapshot{}, err
	}
	inputSnapshotDigest := bodyTranslationDigestLines([]string{
		dictionaryDigest,
		personaDigest,
		metadataDigest,
		promptDigest,
		fmt.Sprintf("provider_target_count=%d", snapshot.ProviderTargetCount),
		fmt.Sprintf("exact_exclusion_count=%d", snapshot.ExactExclusionCount),
		fmt.Sprintf("partial_constraint_count=%d", snapshot.PartialConstraintCount),
	})

	snapshot.DictionaryDigest = dictionaryDigest
	snapshot.PersonaDigest = personaDigest
	snapshot.MetadataDigest = metadataDigest
	snapshot.PromptDigest = promptDigest
	snapshot.InputSnapshotDigest = inputSnapshotDigest
	snapshot.TargetCount = len(snapshot.Fields)
	return snapshot, nil
}

func sortedBodyTranslationRecordFields(
	fields []repository.TranslationField,
) []repository.TranslationField {
	recordFields := append([]repository.TranslationField(nil), fields...)
	sort.SliceStable(recordFields, func(left int, right int) bool {
		return recordFields[left].FieldOrder < recordFields[right].FieldOrder
	})
	return recordFields
}

func buildBodyTranslationSnapshotField(
	record repository.TranslationRecord,
	field repository.TranslationField,
	dictionaryEntries []repository.DictionaryEntry,
) bodyTranslationSnapshotField {
	snapshotField := bodyTranslationSnapshotField{
		TranslationFieldID:         field.ID,
		RecordType:                 strings.TrimSpace(record.RecordType),
		FieldType:                  strings.TrimSpace(field.SubrecordType),
		SourceText:                 field.SourceText,
		CompleteMatchExclusions:    make([]BodyTranslationDictionaryExactMatchExclusion, 0),
		PartialMatchConstraints:    make([]BodyTranslationPartialMatchConstraint, 0),
		IncludedInProviderRequests: true,
	}
	addBodyTranslationDictionaryMatches(&snapshotField, dictionaryEntries)
	return snapshotField
}

func addBodyTranslationDictionaryMatches(
	snapshotField *bodyTranslationSnapshotField,
	dictionaryEntries []repository.DictionaryEntry,
) {
	sourceText := strings.TrimSpace(snapshotField.SourceText)
	normalizedSourceText := normalizeBodyTranslationDictionaryText(sourceText)
	for _, entry := range dictionaryEntries {
		addBodyTranslationDictionaryMatch(snapshotField, normalizedSourceText, entry)
	}
}

func addBodyTranslationDictionaryMatch(
	snapshotField *bodyTranslationSnapshotField,
	normalizedSourceText string,
	entry repository.DictionaryEntry,
) {
	sourceTerm := strings.TrimSpace(entry.SourceTerm)
	if sourceTerm == "" {
		return
	}
	normalizedSourceTerm := normalizeBodyTranslationDictionaryText(sourceTerm)
	switch strings.ToLower(strings.TrimSpace(entry.TermKind)) {
	case bodyTranslationDictionaryExactKind:
		if normalizedSourceText != "" && normalizedSourceText == normalizedSourceTerm {
			snapshotField.CompleteMatchExclusions = append(snapshotField.CompleteMatchExclusions, BodyTranslationDictionaryExactMatchExclusion{
				SourceText:     sourceTerm,
				TranslatedText: strings.TrimSpace(entry.TranslatedTerm),
			})
		}
	case bodyTranslationDictionaryPartialKind:
		if normalizedSourceText != "" && normalizedSourceTerm != "" &&
			strings.Contains(normalizedSourceText, normalizedSourceTerm) {
			snapshotField.PartialMatchConstraints = append(snapshotField.PartialMatchConstraints, BodyTranslationPartialMatchConstraint{
				SourceText:          sourceTerm,
				RequiredTranslation: strings.TrimSpace(entry.TranslatedTerm),
			})
		}
	}
}

func applyBodyTranslationSnapshotField(
	snapshot *bodyTranslationInputSnapshot,
	snapshotField bodyTranslationSnapshotField,
) {
	if len(snapshotField.CompleteMatchExclusions) > 0 {
		snapshotField.IncludedInProviderRequests = false
		snapshot.ExactExclusionCount += len(snapshotField.CompleteMatchExclusions)
		snapshot.SkippedReasons = append(snapshot.SkippedReasons, bodyTranslationSkippedReasonExactDictionary)
	} else {
		snapshot.ProviderTargetCount++
	}
	snapshot.PartialConstraintCount += len(snapshotField.PartialMatchConstraints)
	snapshot.Fields = append(snapshot.Fields, snapshotField)
	snapshot.TargetCount = len(snapshot.Fields)
}

func buildBodyTranslationDictionaryDigestLines(entries []repository.DictionaryEntry) []string {
	lines := make([]string, 0, len(entries))
	for _, entry := range entries {
		lines = append(lines, strings.Join([]string{
			strings.TrimSpace(entry.DictionaryLifecycle),
			strings.TrimSpace(entry.DictionaryScope),
			strings.TrimSpace(entry.DictionarySource),
			strings.TrimSpace(entry.SourceTerm),
			strings.TrimSpace(entry.TranslatedTerm),
			strings.TrimSpace(entry.TermKind),
		}, "|"))
	}
	sort.Strings(lines)
	return lines
}

func buildBodyTranslationMetadataDigestLines(
	job repository.TranslationJob,
	records []repository.TranslationRecord,
	fields []bodyTranslationSnapshotField,
	execution BodyTranslationPhaseExecutionSummaryReadModel,
) []string {
	lines := []string{
		fmt.Sprintf("job_id=%d", job.ID),
		fmt.Sprintf("xedit_id=%d", job.XEditExtractedDataID),
		"provider=" + strings.TrimSpace(execution.Provider),
		"model=" + strings.TrimSpace(execution.Model),
		"execution_mode=" + strings.TrimSpace(execution.ExecutionMode),
		"credential_ref=" + strings.TrimSpace(execution.CredentialRef),
	}
	for _, record := range records {
		lines = append(lines, strings.Join([]string{
			fmt.Sprintf("record_id=%d", record.ID),
			strings.TrimSpace(record.RecordType),
			strings.TrimSpace(record.FormID),
			strings.TrimSpace(record.EditorID),
		}, "|"))
	}
	for _, field := range fields {
		lines = append(lines, strings.Join([]string{
			fmt.Sprintf("field_id=%d", field.TranslationFieldID),
			field.RecordType,
			field.FieldType,
			strings.TrimSpace(field.SourceText),
		}, "|"))
	}
	sort.Strings(lines)
	return lines
}

func buildBodyTranslationPromptDigest(
	fields []bodyTranslationSnapshotField,
	persona repository.Persona,
	execution BodyTranslationPhaseExecutionSummaryReadModel,
) (string, error) {
	digests := make([]string, 0, len(fields))
	for _, field := range fields {
		if !field.IncludedInProviderRequests {
			continue
		}
		prompt, err := BuildBodyTranslationPrompt(BodyTranslationProviderRequest{
			Provider:                execution.Provider,
			Model:                   execution.Model,
			ExecutionMode:           execution.ExecutionMode,
			CredentialRef:           execution.CredentialRef,
			RequestUnitID:           fmt.Sprintf("body-field-%d", field.TranslationFieldID),
			FieldCorrelationKey:     fmt.Sprintf("field:%d", field.TranslationFieldID),
			RecordType:              field.RecordType,
			FieldType:               field.FieldType,
			SourceText:              bodyTranslationPromptSourceText(field.SourceText),
			SourceLanguage:          bodyTranslationSourceLanguageDefaultValue,
			TargetLanguage:          bodyTranslationTargetLanguageDefaultValue,
			PersonaSummary:          bodyTranslationPersonaSummary(persona),
			CompleteMatchExclusions: field.CompleteMatchExclusions,
			PartialMatchConstraints: field.PartialMatchConstraints,
		})
		if err != nil {
			return "", fmt.Errorf("build body translation prompt digest: %w", err)
		}
		digests = append(digests, bodyTranslationPromptDigest(prompt))
	}
	sort.Strings(digests)
	return bodyTranslationDigestLines(digests), nil
}

func bodyTranslationPersonaSummary(persona repository.Persona) string {
	parts := []string{
		strings.TrimSpace(persona.PersonaDescription),
		strings.TrimSpace(persona.SpeechStyle),
		strings.TrimSpace(persona.PersonalitySummary),
	}
	filtered := make([]string, 0, len(parts))
	for _, part := range parts {
		if part == "" {
			continue
		}
		filtered = append(filtered, part)
	}
	return strings.Join(filtered, "; ")
}

func bodyTranslationDigestLines(lines []string) string {
	sorted := append([]string(nil), lines...)
	sort.Strings(sorted)
	sum := sha256.Sum256([]byte(strings.Join(sorted, "\n")))
	return "sha256:" + hex.EncodeToString(sum[:])
}

func normalizeBodyTranslationDictionaryText(value string) string {
	return strings.ToLower(strings.TrimSpace(value))
}

func bodyTranslationPromptSourceText(value string) string {
	trimmed := strings.TrimSpace(value)
	if trimmed == "" {
		return "[empty_source_text]"
	}
	return value
}
