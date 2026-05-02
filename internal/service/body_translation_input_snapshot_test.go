package service

import (
	"testing"
	"time"

	"aitranslationenginejp/internal/repository"
)

func TestBuildBodyTranslationInputSnapshotExactMatchExcludesProviderTarget(t *testing.T) {
	now := time.Date(2026, 5, 3, 9, 0, 0, 0, time.UTC)
	job := repository.TranslationJob{ID: 10, XEditExtractedDataID: 20, CreatedAt: now}
	records := []repository.TranslationRecord{{ID: 100, RecordType: "NPC_", FormID: "0001", EditorID: "EditorA"}}
	fieldsByRecordID := map[int64][]repository.TranslationField{
		100: {
			{ID: 1000, TranslationRecordID: 100, SubrecordType: "FULL", SourceText: "Hello", FieldOrder: 1},
			{ID: 1001, TranslationRecordID: 100, SubrecordType: "DESC", SourceText: "Foo Bar", FieldOrder: 2},
		},
	}
	dictionaryEntries := []repository.DictionaryEntry{
		{SourceTerm: "hello", TranslatedTerm: "こんにちは", TermKind: "exact"},
		{SourceTerm: "foo", TranslatedTerm: "フー", TermKind: "partial"},
	}
	persona := repository.Persona{ID: 30, PersonaDescription: "丁寧", SpeechStyle: "敬語", PersonalitySummary: "穏やか", EvidenceUtteranceCount: 2}
	execution := BodyTranslationPhaseExecutionSummaryReadModel{
		CredentialRef: "cred", Provider: "fake", Model: "model", ExecutionMode: "single_request",
	}

	snapshot, err := buildBodyTranslationInputSnapshot(job, records, fieldsByRecordID, dictionaryEntries, persona, execution)
	if err != nil {
		t.Fatalf("expected success, got %v", err)
	}
	if snapshot.ProviderTargetCount != 1 {
		t.Fatalf("expected provider target count 1, got %d", snapshot.ProviderTargetCount)
	}
	if snapshot.ExactExclusionCount != 1 {
		t.Fatalf("expected exact exclusion count 1, got %d", snapshot.ExactExclusionCount)
	}
	if snapshot.PartialConstraintCount != 1 {
		t.Fatalf("expected partial constraint count 1, got %d", snapshot.PartialConstraintCount)
	}
	if len(snapshot.SkippedReasons) != 1 || snapshot.SkippedReasons[0] != bodyTranslationSkippedReasonExactDictionary {
		t.Fatalf("expected exact dictionary skipped reason, got %#v", snapshot.SkippedReasons)
	}
	if snapshot.Fields[0].IncludedInProviderRequests {
		t.Fatalf("expected exact match field excluded from provider requests")
	}
	if !snapshot.Fields[1].IncludedInProviderRequests {
		t.Fatalf("expected non-exact field included in provider requests")
	}
}

func TestBuildBodyTranslationInputSnapshotDigestIsStableAcrossEntryOrder(t *testing.T) {
	job := repository.TranslationJob{ID: 12, XEditExtractedDataID: 22}
	records := []repository.TranslationRecord{{ID: 200, RecordType: "DIAL", FormID: "0010", EditorID: "E10"}}
	fieldsByRecordID := map[int64][]repository.TranslationField{
		200: {
			{ID: 2001, TranslationRecordID: 200, SubrecordType: "FULL", SourceText: "Beta", FieldOrder: 2},
			{ID: 2000, TranslationRecordID: 200, SubrecordType: "FULL", SourceText: "Alpha", FieldOrder: 1},
		},
	}
	persona := repository.Persona{ID: 44, PersonaDescription: "A", SpeechStyle: "B", PersonalitySummary: "C", EvidenceUtteranceCount: 1}
	execution := BodyTranslationPhaseExecutionSummaryReadModel{CredentialRef: "cred", Provider: "fake", Model: "model", ExecutionMode: "single_request"}
	entriesA := []repository.DictionaryEntry{
		{SourceTerm: "alpha", TranslatedTerm: "アルファ", TermKind: "exact", DictionaryLifecycle: "production", DictionaryScope: "global", DictionarySource: "xml"},
		{SourceTerm: "beta", TranslatedTerm: "ベータ", TermKind: "partial", DictionaryLifecycle: "production", DictionaryScope: "global", DictionarySource: "xml"},
	}
	entriesB := []repository.DictionaryEntry{entriesA[1], entriesA[0]}

	snapshotA, err := buildBodyTranslationInputSnapshot(job, records, fieldsByRecordID, entriesA, persona, execution)
	if err != nil {
		t.Fatalf("expected snapshotA success, got %v", err)
	}
	snapshotB, err := buildBodyTranslationInputSnapshot(job, records, fieldsByRecordID, entriesB, persona, execution)
	if err != nil {
		t.Fatalf("expected snapshotB success, got %v", err)
	}

	if snapshotA.DictionaryDigest != snapshotB.DictionaryDigest || snapshotA.PromptDigest != snapshotB.PromptDigest || snapshotA.InputSnapshotDigest != snapshotB.InputSnapshotDigest {
		t.Fatalf("expected stable digests across dictionary order, got A=%#v B=%#v", snapshotA, snapshotB)
	}
}
