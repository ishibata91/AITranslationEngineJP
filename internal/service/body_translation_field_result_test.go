package service

import (
	"testing"

	"aitranslationenginejp/internal/repository"
)

func TestValidateBodyTranslationProviderResultsReturnsProviderFailure(t *testing.T) {
	targets := []normalizedBodyTranslationFieldTarget{bodyTranslationTargetFixture()}
	results := []BodyTranslationProviderResult{
		{
			FieldCorrelationKey: "field:1",
			Failure:             &BodyTranslationProviderFailure{Kind: BodyTranslationProviderErrorKindProviderFailure, Retryable: true, IsRedacted: true},
		},
	}

	validated, failure, err := validateBodyTranslationProviderResults(targets, results)
	if err != nil {
		t.Fatalf("expected no direct error, got %v", err)
	}
	if len(validated) != 0 {
		t.Fatalf("expected no validated results, got %d", len(validated))
	}
	if failure == nil || failure.errorKind != "provider_failure" || !failure.retryable {
		t.Fatalf("expected provider failure summary, got %#v", failure)
	}
}

func TestValidateBodyTranslationProviderResultsReturnsInvalidResponseOnMissingCandidate(t *testing.T) {
	targets := []normalizedBodyTranslationFieldTarget{bodyTranslationTargetFixture()}
	results := []BodyTranslationProviderResult{{FieldCorrelationKey: "field:1"}}

	_, failure, err := validateBodyTranslationProviderResults(targets, results)
	if err != nil {
		t.Fatalf("expected no direct error, got %v", err)
	}
	if failure == nil || failure.errorKind != "invalid_provider_response" || !failure.redacted {
		t.Fatalf("expected invalid provider response failure, got %#v", failure)
	}
}

func TestValidateBodyTranslationProviderResultsReturnsProtectionValidationFailure(t *testing.T) {
	targets := []normalizedBodyTranslationFieldTarget{
		{
			TranslationFieldID:  2,
			FieldCorrelationKey: "field:2",
			OutputStatus:        bodyTranslationOutputStatusTranslated,
			ProtectedElements: []BodyTranslationProtectedElement{
				{ElementType: "tag", SourceText: "<name>", Digest: "sha256:tag"},
			},
			Field:      bodyTranslationFieldFixture(2),
			RecordType: "NPC_",
		},
	}
	results := []BodyTranslationProviderResult{
		{
			FieldCorrelationKey: "field:2",
			TranslatedCandidate: &BodyTranslationTranslatedCandidate{
				FieldCorrelationKey: "field:2",
				RecordType:          "NPC_",
				FieldType:           "FULL",
				TranslatedText:      "タグ欠落",
			},
			ProtectionValidationTarget: &BodyTranslationProtectionValidationTarget{
				FieldCorrelationKey: "field:2",
				TranslatedText:      "タグ欠落",
			},
		},
	}

	_, failure, err := validateBodyTranslationProviderResults(targets, results)
	if err != nil {
		t.Fatalf("expected no direct error, got %v", err)
	}
	if failure == nil || failure.errorKind != "protection_validation_failed" || failure.summaryAdd.ProtectionFailedCount != 1 {
		t.Fatalf("expected protection failure summary, got %#v", failure)
	}
}

func TestValidateBodyTranslationProviderResultsRejectsDuplicateCorrelationKey(t *testing.T) {
	targets := []normalizedBodyTranslationFieldTarget{bodyTranslationTargetFixture()}
	result := bodyTranslationProviderSuccessFixture("field:1")

	_, _, err := validateBodyTranslationProviderResults(targets, []BodyTranslationProviderResult{result, result})
	if err == nil {
		t.Fatal("expected duplicate correlation key error")
	}
}

func TestValidateBodyTranslationProviderResultsRejectsMissingCorrelationForTarget(t *testing.T) {
	targets := []normalizedBodyTranslationFieldTarget{bodyTranslationTargetFixture()}
	results := []BodyTranslationProviderResult{bodyTranslationProviderSuccessFixture("field:999")}

	_, _, err := validateBodyTranslationProviderResults(targets, results)
	if err == nil {
		t.Fatal("expected missing provider result error")
	}
}

func bodyTranslationTargetFixture() normalizedBodyTranslationFieldTarget {
	return normalizedBodyTranslationFieldTarget{
		TranslationFieldID:  1,
		FieldCorrelationKey: "field:1",
		OutputStatus:        bodyTranslationOutputStatusTranslated,
		Field:               bodyTranslationFieldFixture(1),
		RecordType:          "NPC_",
	}
}

func bodyTranslationFieldFixture(fieldID int64) repository.TranslationField {
	return repository.TranslationField{
		ID:            fieldID,
		SubrecordType: "FULL",
		SourceText:    "source",
	}
}

func bodyTranslationProviderSuccessFixture(correlationKey string) BodyTranslationProviderResult {
	return BodyTranslationProviderResult{
		FieldCorrelationKey: correlationKey,
		TranslatedCandidate: &BodyTranslationTranslatedCandidate{
			FieldCorrelationKey: correlationKey,
			RecordType:          "NPC_",
			FieldType:           "FULL",
			TranslatedText:      "翻訳済み <name>",
		},
		ProtectionValidationTarget: &BodyTranslationProtectionValidationTarget{
			FieldCorrelationKey: correlationKey,
			TranslatedText:      "翻訳済み <name>",
		},
	}
}
