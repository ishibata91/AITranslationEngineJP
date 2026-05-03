package usecase

import (
	"context"
	"fmt"
	"slices"
	"strings"
	"testing"
)

type translationOutputArtifactContractAPI interface {
	GetTranslationOutputReview(
		context.Context,
		GetTranslationOutputReviewRequest,
	) (TranslationOutputReviewResult, error)
	GetTranslationOutputDiffPreview(
		context.Context,
		GetTranslationOutputDiffPreviewRequest,
	) (TranslationOutputDiffPreviewResult, error)
	GenerateXTranslatorOutputArtifact(
		context.Context,
		GenerateXTranslatorOutputArtifactRequest,
	) (XTranslatorOutputArtifactCommandResult, error)
	RegenerateXTranslatorOutputArtifact(
		context.Context,
		RegenerateXTranslatorOutputArtifactRequest,
	) (XTranslatorOutputArtifactCommandResult, error)
}

func TestSCN_TOA_001_OutputReviewPublicSeamFreezesReviewContract(t *testing.T) {
	var api translationOutputArtifactContractAPI
	_ = api

	result := TranslationOutputReviewResult{
		CompletedJobs: []TranslationOutputCompletedJobSummary{
			{
				JobID:           101,
				JobStatus:       "completed",
				ArtifactStatus:  "not_generated",
				OutputReady:     true,
				TranslatedCount: 2,
				OutputStatusDistribution: map[string]int{
					"translated": 1,
					"cached":     1,
				},
			},
		},
		SelectedJob: TranslationOutputSelectedJobSummary{
			JobID:           101,
			JobStatus:       "completed",
			BodyPhaseStatus: "completed",
			OutputReady:     true,
			ResultSummary: TranslationOutputResultSummary{
				TranslatedCount: 2,
				RowCount:        2,
				InputProvenance: TranslationOutputInputProvenanceSummary{
					InputSnapshotDigest: "sha256:input",
					SourceFileDigest:    "sha256:plugin",
				},
			},
		},
		OutputReadiness: TranslationOutputReadinessSummary{
			Ready:         true,
			Retryable:     false,
			RejectionKind: "",
		},
		ArtifactStatus: TranslationOutputArtifactStatusSummary{
			ArtifactID:     0,
			Status:         "not_generated",
			RowCount:       0,
			CurrentVersion: false,
		},
	}

	if len(result.CompletedJobs) != 1 {
		t.Fatalf("expected completed job list, got %#v", result.CompletedJobs)
	}
	if !result.SelectedJob.OutputReady || !result.OutputReadiness.Ready {
		t.Fatalf("expected completed job to be output ready, got %#v", result)
	}
	if result.SelectedJob.ResultSummary.InputProvenance.InputSnapshotDigest == "" {
		t.Fatal("expected input provenance digest")
	}
}

func TestSCN_TOA_002_GenerateContractRejectsNotReadyJobsWithRedactedReason(t *testing.T) {
	response := XTranslatorOutputArtifactCommandResult{
		JobID:          202,
		ArtifactID:     0,
		ArtifactStatus: "rejected",
		RowCount:       0,
		TargetGame:     "skyrim_se",
		ErrorSummary: &TranslationOutputArtifactErrorSummary{
			ErrorKind:  TranslationOutputArtifactErrorKindNotCompleted,
			Reason:     "job is not completed",
			Retryable:  false,
			IsRedacted: true,
		},
	}

	if response.ArtifactStatus == "success" {
		t.Fatalf("expected not-ready job to reject success artifact, got %#v", response)
	}
	if response.ErrorSummary == nil {
		t.Fatal("expected redacted rejection summary")
	}
	if response.ErrorSummary.ErrorKind != TranslationOutputArtifactErrorKindNotCompleted {
		t.Fatalf("expected not_completed, got %#v", response.ErrorSummary)
	}
	assertTranslationOutputArtifactHasNoForbiddenSecretFields(t, response)
}

func TestSCN_TOA_003_DiffPreviewContractFreezesRowsStatusMappingAndValidation(t *testing.T) {
	preview := TranslationOutputDiffPreviewResult{
		JobID:      303,
		ArtifactID: 404,
		Rows: []TranslationOutputDiffPreviewRow{
			{
				FieldID:              501,
				RowDigest:            "sha256:row-cached",
				EDID:                 "CachedNPC",
				REC:                  "NPC_",
				FIELD:                "FULL",
				FORMID:               "00000001",
				SourceExcerpt:        "Hello",
				DestExcerpt:          "こんにちは",
				XTranslatorStatus:    1,
				InternalOutputStatus: "cached",
				RowReflectionSummary: "cached output reflected as xTranslator status 1",
				StaleReason:          "",
				CanRegenerate:        false,
			},
		},
		CompatibilitySummary: TranslationOutputCompatibilitySummary{
			Passed:       true,
			WarningCount: 0,
			RejectCount:  0,
		},
	}

	if len(preview.Rows) != 1 {
		t.Fatalf("expected one preview row, got %#v", preview.Rows)
	}
	row := preview.Rows[0]
	if row.EDID == "" || row.REC == "" || row.FIELD == "" || row.FORMID == "" {
		t.Fatalf("expected xTranslator required columns, got %#v", row)
	}
	if row.InternalOutputStatus != "cached" || row.XTranslatorStatus != 1 {
		t.Fatalf("expected cached to map to Status=1, got %#v", row)
	}
}

func TestSCN_TOA_006_RegenerateContractKeepsSingleArtifactAndRowPerField(t *testing.T) {
	response := XTranslatorOutputArtifactCommandResult{
		JobID:          606,
		ArtifactID:     707,
		ArtifactStatus: "success",
		RowCount:       2,
		FilePath:       "/tmp/output/skyrim.xml",
		TargetGame:     "skyrim_se",
		OperationSummary: TranslationOutputOperationSummary{
			OperationKind:       "regenerate",
			ReplacedArtifactID:  707,
			AffectedFieldIDs:    []int64{801, 802},
			DuplicateRowCreated: false,
		},
	}

	if response.OperationSummary.OperationKind != "regenerate" {
		t.Fatalf("expected regenerate summary, got %#v", response.OperationSummary)
	}
	if response.OperationSummary.DuplicateRowCreated {
		t.Fatal("expected regenerate contract to forbid duplicate rows")
	}
	if response.OperationSummary.ReplacedArtifactID != response.ArtifactID {
		t.Fatalf("expected current artifact to be replaced or updated, got %#v", response)
	}
}

func TestSCN_TOA_010_ErrorKindsAndRedactionObligationAreFrozen(t *testing.T) {
	errorKinds := []TranslationOutputArtifactErrorKind{
		TranslationOutputArtifactErrorKindNotCompleted,
		TranslationOutputArtifactErrorKindCanceled,
		TranslationOutputArtifactErrorKindStatusMismatch,
		TranslationOutputArtifactErrorKindMissingRequiredRowField,
		TranslationOutputArtifactErrorKindUnknownOutputStatus,
		TranslationOutputArtifactErrorKindXMLSerializationFailed,
		TranslationOutputArtifactErrorKindFileWriteFailed,
		TranslationOutputArtifactErrorKindArtifactSaveFailed,
		TranslationOutputArtifactErrorKindCompatibilityRejected,
		TranslationOutputArtifactErrorKindSecretRedacted,
	}

	for _, expected := range []TranslationOutputArtifactErrorKind{
		"not_completed",
		"canceled",
		"status_mismatch",
		"missing_required_row_field",
		"unknown_output_status",
		"xml_serialization_failed",
		"file_write_failed",
		"artifact_save_failed",
		"compatibility_rejected",
		"secret_redacted",
	} {
		if !slices.Contains(errorKinds, expected) {
			t.Fatalf("expected public error kind %q in %#v", expected, errorKinds)
		}
	}

	obligation := TranslationOutputArtifactPublicRedactionFieldObligation()
	for _, allowed := range []string{
		"job_id",
		"artifact_id",
		"field_id",
		"row_digest",
		"count",
		"status",
		"file_path",
		"target_game",
		"error_kind",
		"retryable",
	} {
		if !slices.Contains(obligation.AllowedReferenceFields, allowed) {
			t.Fatalf("expected allowed reference field %q in %#v", allowed, obligation)
		}
	}
	for _, forbidden := range []string{
		"secret",
		"api_key",
		"token",
		"authorization_header",
		"provider_raw_request",
		"provider_raw_response",
		"decryptable_value",
		"full_source_text",
		"full_dest_text",
	} {
		if !slices.Contains(obligation.ForbiddenOutputFields, forbidden) {
			t.Fatalf("expected forbidden output field %q in %#v", forbidden, obligation)
		}
	}
}

func assertTranslationOutputArtifactHasNoForbiddenSecretFields(
	t *testing.T,
	response XTranslatorOutputArtifactCommandResult,
) {
	t.Helper()
	serialized := stringifyForContractAssertion(response)
	for _, forbidden := range []string{
		"sk-live-secret",
		"apiKey",
		"providerRawRequest",
		"providerRawResponse",
		"authorizationHeader",
		"decrypted",
		"fullSourceText",
		"fullDestText",
	} {
		if strings.Contains(serialized, forbidden) {
			t.Fatalf("expected response not to expose %q: %s", forbidden, serialized)
		}
	}
}

func stringifyForContractAssertion(value any) string {
	return fmt.Sprintf("%#v", value)
}
