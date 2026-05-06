package apitest

import (
	"context"
	"encoding/json"
	"encoding/xml"
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	controllerwails "aitranslationenginejp/internal/controller/wails"
	"aitranslationenginejp/internal/repository"
	"aitranslationenginejp/internal/service"
	"aitranslationenginejp/internal/usecase"
)

func TestSCN_TOA_001_OutputReviewReadinessReturnsCompletedJobSummary(t *testing.T) {
	fixture := newTranslationOutputReviewReadinessAPIFixture(t)
	controller := fixture.controller()

	selectedJobID := int64(101)
	result, err := controller.GetTranslationOutputReview(
		controllerwails.GetTranslationOutputReviewRequestDTO{SelectedJobID: &selectedJobID},
	)
	if err != nil {
		t.Fatalf("SCN-TOA-001 expected Output Review summary: %v", err)
	}
	if len(result.CompletedJobs) != 2 {
		t.Fatalf("SCN-TOA-001 expected completed output candidates only, got %#v", result.CompletedJobs)
	}
	completedJob, ok := findTranslationOutputReadinessCompletedJob(result.CompletedJobs, 101)
	if !ok {
		t.Fatalf("SCN-TOA-001 expected selected completed job in list, got %#v", result.CompletedJobs)
	}
	if _, ok := findTranslationOutputReadinessCompletedJob(result.CompletedJobs, 202); ok {
		t.Fatalf("SCN-TOA-001 expected running job to stay out of completed list, got %#v", result.CompletedJobs)
	}
	if _, ok := findTranslationOutputReadinessCompletedJob(result.CompletedJobs, 404); ok {
		t.Fatalf("SCN-TOA-001 expected canceled job to stay out of completed list, got %#v", result.CompletedJobs)
	}
	if !completedJob.OutputReady || !result.OutputReadiness.Ready || !result.SelectedJob.OutputReady {
		t.Fatalf("SCN-TOA-001 expected completed job to be output ready, got %#v", result)
	}
	if result.SelectedJob.ResultSummary.TranslatedCount != 2 || result.SelectedJob.ResultSummary.RowCount != 2 {
		t.Fatalf("SCN-TOA-001 expected translated and row counts, got %#v", result.SelectedJob.ResultSummary)
	}
	if result.SelectedJob.ResultSummary.InputProvenance.InputSnapshotDigest != "sha256:body-input" {
		t.Fatalf("SCN-TOA-001 expected input provenance digest, got %#v", result.SelectedJob.ResultSummary.InputProvenance)
	}
	if result.ArtifactStatus.Status != "not_generated" || result.ArtifactStatus.RowCount != 0 {
		t.Fatalf("SCN-TOA-001 expected empty artifact status, got %#v", result.ArtifactStatus)
	}
	assertTranslationOutputReadinessHasNoForbiddenText(t, result)
}

func TestSCN_TOA_002_OutputReviewRejectsNotReadyJobsWithReason(t *testing.T) {
	fixture := newTranslationOutputReviewReadinessAPIFixture(t)
	controller := fixture.controller()

	selectedJobID := int64(202)
	result, err := controller.GetTranslationOutputReview(
		controllerwails.GetTranslationOutputReviewRequestDTO{SelectedJobID: &selectedJobID},
	)
	if err != nil {
		t.Fatalf("SCN-TOA-002 expected not-ready review payload: %v", err)
	}
	if result.OutputReadiness.Ready || result.SelectedJob.OutputReady {
		t.Fatalf("SCN-TOA-002 expected not-ready selected job, got %#v", result)
	}
	if result.OutputReadiness.RejectionKind != "not_completed" {
		t.Fatalf("SCN-TOA-002 expected not_completed rejection, got %#v", result.OutputReadiness)
	}
	if len(result.RejectionReasons) == 0 || result.RejectionReasons[0].ErrorKind != "not_completed" {
		t.Fatalf("SCN-TOA-002 expected redacted rejection reason, got %#v", result.RejectionReasons)
	}
	if fixture.store.translationArtifactCount != 0 || fixture.store.xtranslatorRowCount != 0 {
		t.Fatalf(
			"SCN-TOA-002 expected readiness rejection to create no artifact or row, got artifacts=%d rows=%d",
			fixture.store.translationArtifactCount,
			fixture.store.xtranslatorRowCount,
		)
	}
	assertTranslationOutputReadinessHasNoForbiddenText(t, result)
}

func TestSCN_TOA_008_ZeroTargetCompletedJobStaysOutputReady(t *testing.T) {
	fixture := newTranslationOutputReviewReadinessAPIFixture(t)
	controller := fixture.controller()

	selectedJobID := int64(303)
	result, err := controller.GetTranslationOutputReview(
		controllerwails.GetTranslationOutputReviewRequestDTO{SelectedJobID: &selectedJobID},
	)
	if err != nil {
		t.Fatalf("SCN-TOA-008 expected zero target Output Review summary: %v", err)
	}
	if !result.OutputReadiness.Ready || !result.SelectedJob.OutputReady {
		t.Fatalf("SCN-TOA-008 expected zero target completed job to be output ready, got %#v", result)
	}
	if result.SelectedJob.ResultSummary.TranslatedCount != 0 || result.SelectedJob.ResultSummary.RowCount != 0 {
		t.Fatalf("SCN-TOA-008 expected zero translated rows, got %#v", result.SelectedJob.ResultSummary)
	}
	if result.ArtifactStatus.Status != "not_generated" {
		t.Fatalf("SCN-TOA-008 expected zero-row job to still expose artifact status, got %#v", result.ArtifactStatus)
	}
	assertTranslationOutputReadinessHasNoForbiddenText(t, result)
}

func TestSCN_TOA_010_OutputReviewDoesNotResolveProviderSecretOrExposeRawPayload(t *testing.T) {
	fixture := newTranslationOutputReviewReadinessAPIFixture(t)
	controller := fixture.controller()

	selectedJobID := int64(101)
	result, err := controller.GetTranslationOutputReview(
		controllerwails.GetTranslationOutputReviewRequestDTO{SelectedJobID: &selectedJobID},
	)
	if err != nil {
		t.Fatalf("SCN-TOA-010 expected redacted Output Review summary: %v", err)
	}
	if fixture.store.providerCallCount != 0 || fixture.store.secretResolutionCallCount != 0 {
		t.Fatalf(
			"SCN-TOA-010 expected no provider or secret resolution call, got provider=%d secret=%d",
			fixture.store.providerCallCount,
			fixture.store.secretResolutionCallCount,
		)
	}
	if !result.OutputReadiness.Ready || result.ArtifactStatus.Status == "" {
		t.Fatalf("SCN-TOA-010 expected audit-safe readiness and artifact status, got %#v", result)
	}
	assertTranslationOutputReadinessHasNoForbiddenText(t, result)
}

func TestSCN_TOA_003_DiffPreviewBuildsRequiredRowsAndCachedStatus(t *testing.T) {
	fixture := newTranslationOutputReviewReadinessAPIFixture(t)
	controller := fixture.controller()

	result, err := controller.GetTranslationOutputDiffPreview(
		controllerwails.GetTranslationOutputDiffPreviewRequestDTO{JobID: 101, ArtifactID: 901},
	)
	if err != nil {
		t.Fatalf("SCN-TOA-003 expected diff preview rows: %v", err)
	}

	cachedRow, ok := findTranslationOutputDiffPreviewRow(result.Rows, 702)
	if !ok {
		t.Fatalf("SCN-TOA-003 expected cached field row, got %#v", result.Rows)
	}
	if cachedRow.EDID == "" || cachedRow.REC == "" || cachedRow.FIELD == "" || cachedRow.FORMID == "" {
		t.Fatalf("SCN-TOA-003 expected required xTranslator columns, got %#v", cachedRow)
	}
	if cachedRow.InternalOutputStatus != "cached" || cachedRow.XTranslatorStatus != 1 {
		t.Fatalf("SCN-TOA-003 expected cached to map to Status=1, got %#v", cachedRow)
	}
	if !strings.Contains(cachedRow.RowReflectionSummary, "cached") {
		t.Fatalf("SCN-TOA-003 expected dictionary-cache fact in internal summary, got %#v", cachedRow)
	}
	assertTranslationOutputReadinessHasNoForbiddenText(t, result)
}

func TestSCN_TOA_008_WailsResponseKeepsEmptyCompletedJobsAndRowsFields(t *testing.T) {
	fixture := newTranslationOutputReviewReadinessAPIFixture(t)
	fixture.store.jobs = []repository.TranslationJob{fixture.store.jobs[1]}
	controller := fixture.controller()
	selectedJobID := int64(202)

	review, err := controller.GetTranslationOutputReview(
		controllerwails.GetTranslationOutputReviewRequestDTO{SelectedJobID: &selectedJobID},
	)
	if err != nil {
		t.Fatalf("SCN-TOA-008 expected empty completedJobs review response: %v", err)
	}
	reviewJSON, err := json.Marshal(review)
	if err != nil {
		t.Fatalf("SCN-TOA-008 expected review JSON: %v", err)
	}
	if !strings.Contains(string(reviewJSON), `"completedJobs":[]`) {
		t.Fatalf("SCN-TOA-008 expected completedJobs empty array field to remain, got %s", string(reviewJSON))
	}

	diffPreview, err := controller.GetTranslationOutputDiffPreview(
		controllerwails.GetTranslationOutputDiffPreviewRequestDTO{JobID: 303, ArtifactID: 0},
	)
	if err != nil {
		t.Fatalf("SCN-TOA-008 expected empty rows diff preview response: %v", err)
	}
	diffJSON, err := json.Marshal(diffPreview)
	if err != nil {
		t.Fatalf("SCN-TOA-008 expected diff preview JSON: %v", err)
	}
	if !strings.Contains(string(diffJSON), `"rows":[]`) {
		t.Fatalf("SCN-TOA-008 expected rows empty array field to remain, got %s", string(diffJSON))
	}
}

func TestSCN_TOA_003_DiffPreviewRejectsUnknownOutputStatus(t *testing.T) {
	fixture := newTranslationOutputReviewReadinessAPIFixture(t)
	fixture.addCompletedJobWithOutput(505, 805, 701, "machine_pending")
	controller := fixture.controller()

	result, err := controller.GetTranslationOutputDiffPreview(
		controllerwails.GetTranslationOutputDiffPreviewRequestDTO{JobID: 505, ArtifactID: 905},
	)
	if err != nil {
		t.Fatalf("SCN-TOA-003 expected unknown status rejection summary: %v", err)
	}
	if result.CompatibilitySummary.Passed || result.CompatibilitySummary.RejectCount == 0 {
		t.Fatalf("SCN-TOA-003 expected unknown status to reject compatibility, got %#v", result)
	}
	if _, ok := findTranslationOutputDiffPreviewRow(result.Rows, 701); ok {
		t.Fatalf("SCN-TOA-003 expected unknown status not to enter success rows, got %#v", result.Rows)
	}
	assertTranslationOutputReadinessHasNoForbiddenText(t, result)
}

func TestSCN_TOA_004_GenerateArtifactWritesUTF8SkyrimXML(t *testing.T) {
	fixture := newTranslationOutputReviewReadinessAPIFixture(t)
	controller := fixture.controller()
	outputPath := filepath.Join(t.TempDir(), "skyrim-se-output.xml")

	result, err := controller.GenerateXTranslatorOutputArtifact(
		controllerwails.GenerateXTranslatorOutputArtifactRequestDTO{
			JobID:      101,
			TargetGame: "skyrim_se",
			OutputPath: outputPath,
		},
	)
	if err != nil {
		t.Fatalf("SCN-TOA-004 expected generated xTranslator XML artifact: %v", err)
	}
	if result.ArtifactStatus != "success" || result.RowCount != 2 {
		t.Fatalf("SCN-TOA-004 expected success artifact with two rows, got %#v", result)
	}
	assertTranslationOutputGeneratedXML(t, outputPath, "SSETranslator", 2)
	assertTranslationOutputReadinessHasNoForbiddenText(t, result)
}

func TestSCN_TOA_006_RegenerateArtifactDoesNotDuplicateArtifactOrRows(t *testing.T) {
	fixture := newTranslationOutputReviewReadinessAPIFixture(t)
	controller := fixture.controller()
	outputPath := filepath.Join(t.TempDir(), "reoutput.xml")

	generated, err := controller.GenerateXTranslatorOutputArtifact(
		controllerwails.GenerateXTranslatorOutputArtifactRequestDTO{
			JobID:      101,
			TargetGame: "skyrim_se",
			OutputPath: outputPath,
		},
	)
	if err != nil {
		t.Fatalf("SCN-TOA-006 expected initial artifact generation: %v", err)
	}
	regenerated, err := controller.RegenerateXTranslatorOutputArtifact(
		controllerwails.RegenerateXTranslatorOutputArtifactRequestDTO{
			JobID:      101,
			ArtifactID: generated.ArtifactID,
			TargetGame: "skyrim_se",
			OutputPath: outputPath,
		},
	)
	if err != nil {
		t.Fatalf("SCN-TOA-006 expected artifact regeneration: %v", err)
	}
	if regenerated.ArtifactID != generated.ArtifactID {
		t.Fatalf("SCN-TOA-006 expected current artifact update or replacement summary, got generated=%#v regenerated=%#v", generated, regenerated)
	}
	if regenerated.OperationSummary.DuplicateRowCreated || regenerated.RowCount != generated.RowCount {
		t.Fatalf("SCN-TOA-006 expected one row per field after regeneration, got %#v", regenerated)
	}
	assertTranslationOutputReadinessHasNoForbiddenText(t, regenerated)
}

func TestSCN_TOA_007_GenerateArtifactFailureIsNotPublishedAsSuccess(t *testing.T) {
	fixture := newTranslationOutputReviewReadinessAPIFixture(t)
	controller := fixture.controller()
	outputPath := t.TempDir()

	result, err := controller.GenerateXTranslatorOutputArtifact(
		controllerwails.GenerateXTranslatorOutputArtifactRequestDTO{
			JobID:      101,
			TargetGame: "skyrim_se",
			OutputPath: outputPath,
		},
	)
	if err != nil {
		t.Fatalf("SCN-TOA-007 expected redacted write failure response: %v", err)
	}
	if result.ArtifactStatus == "success" || result.ErrorSummary == nil {
		t.Fatalf("SCN-TOA-007 expected failed artifact summary instead of success publish, got %#v", result)
	}
	if result.ErrorSummary.ErrorKind != "file_write_failed" || !result.ErrorSummary.Retryable {
		t.Fatalf("SCN-TOA-007 expected retryable file_write_failed, got %#v", result.ErrorSummary)
	}
	assertTranslationOutputReadinessHasNoForbiddenText(t, result)
}

func TestSCN_TOA_008_GenerateArtifactAllowsZeroTargetCompletedJob(t *testing.T) {
	fixture := newTranslationOutputReviewReadinessAPIFixture(t)
	controller := fixture.controller()
	outputPath := filepath.Join(t.TempDir(), "zero-target.xml")

	result, err := controller.GenerateXTranslatorOutputArtifact(
		controllerwails.GenerateXTranslatorOutputArtifactRequestDTO{
			JobID:      303,
			TargetGame: "skyrim_se",
			OutputPath: outputPath,
		},
	)
	if err != nil {
		t.Fatalf("SCN-TOA-008 expected zero-row artifact generation: %v", err)
	}
	if result.ArtifactStatus != "success" || result.RowCount != 0 {
		t.Fatalf("SCN-TOA-008 expected successful zero-row artifact, got %#v", result)
	}
	assertTranslationOutputGeneratedXML(t, outputPath, "SSETranslator", 0)
	assertTranslationOutputReadinessHasNoForbiddenText(t, result)
}

func TestSCN_TOA_010_GenerateArtifactDoesNotResolveProviderSecretOrExposeRawPayload(t *testing.T) {
	fixture := newTranslationOutputReviewReadinessAPIFixture(t)
	controller := fixture.controller()
	outputPath := filepath.Join(t.TempDir(), "redacted-output.xml")

	result, err := controller.GenerateXTranslatorOutputArtifact(
		controllerwails.GenerateXTranslatorOutputArtifactRequestDTO{
			JobID:      101,
			TargetGame: "skyrim_se",
			OutputPath: outputPath,
		},
	)
	if err != nil {
		t.Fatalf("SCN-TOA-010 expected artifact generation without provider or secret dependency: %v", err)
	}
	if fixture.store.providerCallCount != 0 || fixture.store.secretResolutionCallCount != 0 {
		t.Fatalf(
			"SCN-TOA-010 expected no provider or secret resolution call, got provider=%d secret=%d",
			fixture.store.providerCallCount,
			fixture.store.secretResolutionCallCount,
		)
	}
	assertTranslationOutputReadinessHasNoForbiddenText(t, result)
}

func TestSCN_TOA_005_DiffPreviewReturnsRowDigestAndReoutputState(t *testing.T) {
	fixture := newTranslationOutputReviewReadinessAPIFixture(t)
	controller := fixture.controller()

	result, err := controller.GetTranslationOutputDiffPreview(
		controllerwails.GetTranslationOutputDiffPreviewRequestDTO{JobID: 101, ArtifactID: 901},
	)
	if err != nil {
		t.Fatalf("SCN-TOA-005 expected diff preview: %v", err)
	}

	row, ok := findTranslationOutputDiffPreviewRow(result.Rows, 701)
	if !ok {
		t.Fatalf("SCN-TOA-005 expected translated field row, got %#v", result.Rows)
	}
	if row.RowDigest == "" {
		t.Fatalf("SCN-TOA-005 expected row digest for diff preview, got %#v", row)
	}
	if row.SourceExcerpt == "" || row.DestExcerpt == "" || row.RowReflectionSummary == "" {
		t.Fatalf("SCN-TOA-005 expected Source, Dest, and reflection summary, got %#v", row)
	}
	if row.StaleReason != "" || row.CanRegenerate {
		t.Fatalf("SCN-TOA-005 expected current row not to be stale, got %#v", row)
	}
	assertTranslationOutputReadinessHasNoForbiddenText(t, result)
}

func TestSCN_TOA_005_DiffPreviewMarksStaleArtifactRowForRegenerate(t *testing.T) {
	fixture := newTranslationOutputReviewReadinessAPIFixture(t)
	fixture.store.seedTranslationArtifact(101, 901)
	fixture.store.rowsByArtifactID[901] = []repository.XTranslatorOutputRow{
		{
			ID:                    1,
			TranslationArtifactID: 901,
			JobTranslationFieldID: 801,
			EDID:                  "TranslatedNPC",
			REC:                   "NPC_",
			FIELD:                 "FULL",
			FORMID:                "00000001",
			Source:                "Hello",
			Dest:                  "outdated artifact row",
			Status:                0,
		},
	}
	controller := fixture.controllerWithPersistence()

	result, err := controller.GetTranslationOutputDiffPreview(
		controllerwails.GetTranslationOutputDiffPreviewRequestDTO{JobID: 101, ArtifactID: 901},
	)
	if err != nil {
		t.Fatalf("SCN-TOA-005 expected diff preview with stale artifact rows: %v", err)
	}

	row, ok := findTranslationOutputDiffPreviewRow(result.Rows, 701)
	if !ok {
		t.Fatalf("SCN-TOA-005 expected translated field row, got %#v", result.Rows)
	}
	if row.StaleReason == "" || !row.CanRegenerate {
		t.Fatalf("SCN-TOA-005 expected stale reason and regenerate enablement, got %#v", row)
	}
	assertTranslationOutputReadinessHasNoForbiddenText(t, result)
}

func TestSCN_TOA_009_DiffPreviewReportsCompatibilitySummary(t *testing.T) {
	fixture := newTranslationOutputReviewReadinessAPIFixture(t)
	fixture.addCompletedJobWithOutput(909, 809, 703, "translated")
	controller := fixture.controller()

	result, err := controller.GetTranslationOutputDiffPreview(
		controllerwails.GetTranslationOutputDiffPreviewRequestDTO{JobID: 909, ArtifactID: 909},
	)
	if err != nil {
		t.Fatalf("SCN-TOA-009 expected compatibility summary: %v", err)
	}
	if result.CompatibilitySummary.Passed {
		t.Fatalf("SCN-TOA-009 expected compatibility warning or rejection, got %#v", result.CompatibilitySummary)
	}
	if result.CompatibilitySummary.WarningCount == 0 && result.CompatibilitySummary.RejectCount == 0 {
		t.Fatalf("SCN-TOA-009 expected affected row count in compatibility summary, got %#v", result)
	}
	assertTranslationOutputReadinessHasNoForbiddenText(t, result)
}

func TestSCN_TOA_006_RegenerateMismatchedArtifactIDHasNoWriteOrPersistenceSideEffect(t *testing.T) {
	fixture := newTranslationOutputReviewReadinessAPIFixture(t)
	fixture.store.seedTranslationArtifact(101, 901)
	controller := fixture.controllerWithPersistence()
	outputPath := filepath.Join(t.TempDir(), "mismatched-artifact.xml")

	result, err := controller.RegenerateXTranslatorOutputArtifact(
		controllerwails.RegenerateXTranslatorOutputArtifactRequestDTO{
			JobID:      101,
			ArtifactID: 999,
			TargetGame: "skyrim_se",
			OutputPath: outputPath,
		},
	)
	if err != nil {
		t.Fatalf("SCN-TOA-006 expected redacted mismatch response: %v", err)
	}
	if result.ArtifactStatus == "success" || result.ErrorSummary == nil {
		t.Fatalf("SCN-TOA-006 expected mismatch failure before side effects, got %#v", result)
	}
	if fixture.store.translationArtifactCount != 1 || fixture.store.xtranslatorRowCount != 0 {
		t.Fatalf(
			"SCN-TOA-006 expected no new artifact or row side effect, got artifacts=%d rows=%d",
			fixture.store.translationArtifactCount,
			fixture.store.xtranslatorRowCount,
		)
	}
	if _, statErr := os.Stat(outputPath); !errors.Is(statErr, os.ErrNotExist) {
		t.Fatalf("SCN-TOA-006 expected no XML file write before mismatch failure, statErr=%v", statErr)
	}
	assertTranslationOutputReadinessHasNoForbiddenText(t, result)
}

func TestSCN_TOA_007_DBSaveFailureDoesNotLeaveSuccessXMLOrPublishedArtifact(t *testing.T) {
	fixture := newTranslationOutputReviewReadinessAPIFixture(t)
	fixture.store.failArtifactSave = true
	controller := fixture.controllerWithPersistence()
	outputPath := filepath.Join(t.TempDir(), "save-failure.xml")

	result, err := controller.GenerateXTranslatorOutputArtifact(
		controllerwails.GenerateXTranslatorOutputArtifactRequestDTO{
			JobID:      101,
			TargetGame: "skyrim_se",
			OutputPath: outputPath,
		},
	)
	if err != nil {
		t.Fatalf("SCN-TOA-007 expected redacted save failure response: %v", err)
	}
	if result.ArtifactStatus == "success" || result.ErrorSummary == nil {
		t.Fatalf("SCN-TOA-007 expected artifact_save_failed response, got %#v", result)
	}
	if result.ErrorSummary.ErrorKind != "artifact_save_failed" || !result.ErrorSummary.Retryable {
		t.Fatalf("SCN-TOA-007 expected retryable artifact_save_failed, got %#v", result.ErrorSummary)
	}
	if _, statErr := os.Stat(outputPath); !errors.Is(statErr, os.ErrNotExist) {
		t.Fatalf("SCN-TOA-007 expected no success XML left after DB failure, statErr=%v", statErr)
	}
	if fixture.store.translationArtifactCount != 0 || fixture.store.xtranslatorRowCount != 0 {
		t.Fatalf(
			"SCN-TOA-007 expected no published artifact or rows, got artifacts=%d rows=%d",
			fixture.store.translationArtifactCount,
			fixture.store.xtranslatorRowCount,
		)
	}
	assertTranslationOutputReadinessHasNoForbiddenText(t, result)
}

func TestSCN_TOA_007_PublishTemporaryFileFailureDoesNotLeaveSuccessArtifactOrRows(t *testing.T) {
	fixture := newTranslationOutputReviewReadinessAPIFixture(t)
	writer := &translationOutputReviewReadinessRecordingFileWriter{
		publishErr: errors.New("injected publish failure"),
	}
	controller := fixture.controllerWithPersistenceAndFileWriter(writer)
	outputPath := filepath.Join(t.TempDir(), "publish-failure.xml")

	result, err := controller.GenerateXTranslatorOutputArtifact(
		controllerwails.GenerateXTranslatorOutputArtifactRequestDTO{
			JobID:      101,
			TargetGame: "skyrim_se",
			OutputPath: outputPath,
		},
	)
	if err != nil {
		t.Fatalf("SCN-TOA-007 expected redacted publish failure response: %v", err)
	}
	if result.ArtifactStatus == "success" || result.ErrorSummary == nil {
		t.Fatalf("SCN-TOA-007 expected publish failure instead of success artifact, got %#v", result)
	}
	if result.ErrorSummary.ErrorKind != "file_write_failed" || !result.ErrorSummary.Retryable {
		t.Fatalf("SCN-TOA-007 expected retryable file_write_failed, got %#v", result.ErrorSummary)
	}
	if writer.writeTemporaryCount != 1 || writer.publishTemporaryCount != 1 {
		t.Fatalf(
			"SCN-TOA-007 expected temporary write then publish attempt, got tempWrites=%d publishes=%d",
			writer.writeTemporaryCount,
			writer.publishTemporaryCount,
		)
	}
	if _, statErr := os.Stat(outputPath); !errors.Is(statErr, os.ErrNotExist) {
		t.Errorf("SCN-TOA-007 expected no published XML after publish failure, statErr=%v", statErr)
	}
	if fixture.store.translationArtifactCount != 0 || fixture.store.xtranslatorRowCount != 0 {
		t.Errorf(
			"SCN-TOA-007 expected no success artifact or rows after publish failure, got artifacts=%d rows=%d",
			fixture.store.translationArtifactCount,
			fixture.store.xtranslatorRowCount,
		)
	}

	selectedJobID := int64(101)
	review, reviewErr := controller.GetTranslationOutputReview(
		controllerwails.GetTranslationOutputReviewRequestDTO{SelectedJobID: &selectedJobID},
	)
	if reviewErr != nil {
		t.Fatalf("SCN-TOA-007 expected next Output Review after publish failure: %v", reviewErr)
	}
	if review.ArtifactStatus.Status == "success" || review.ArtifactStatus.RowCount != 0 {
		t.Fatalf(
			"SCN-TOA-007 expected next Output Review not to misread failed publish as success, got %#v",
			review.ArtifactStatus,
		)
	}
	assertTranslationOutputReadinessHasNoForbiddenText(t, result)
	assertTranslationOutputReadinessHasNoForbiddenText(t, review)
}

func TestSCN_TOA_010_InvalidOutputPathNeverReachesFileWriter(t *testing.T) {
	for _, testCase := range []struct {
		name       string
		outputPath string
	}{
		{name: "non_xml", outputPath: filepath.Join(t.TempDir(), "output.txt")},
		{name: "relative", outputPath: "relative-output.xml"},
		{name: "readonly_like", outputPath: "/System/translation-output.xml"},
	} {
		t.Run(testCase.name, func(t *testing.T) {
			fixture := newTranslationOutputReviewReadinessAPIFixture(t)
			writer := &translationOutputReviewReadinessRecordingFileWriter{}
			controller := fixture.controllerWithPersistenceAndFileWriter(writer)

			result, err := controller.GenerateXTranslatorOutputArtifact(
				controllerwails.GenerateXTranslatorOutputArtifactRequestDTO{
					JobID:      101,
					TargetGame: "skyrim_se",
					OutputPath: testCase.outputPath,
				},
			)
			if err != nil {
				t.Fatalf("SCN-TOA-010 expected redacted invalid path response: %v", err)
			}
			if result.ArtifactStatus == "success" || result.ErrorSummary == nil {
				t.Fatalf("SCN-TOA-010 expected invalid path rejection, got %#v", result)
			}
			if writer.writeCount != 0 {
				t.Fatalf("SCN-TOA-010 expected invalid path not to reach file writer, writes=%d", writer.writeCount)
			}
			if fixture.store.translationArtifactCount != 0 || fixture.store.xtranslatorRowCount != 0 {
				t.Fatalf(
					"SCN-TOA-010 expected invalid path to create no artifact or rows, got artifacts=%d rows=%d",
					fixture.store.translationArtifactCount,
					fixture.store.xtranslatorRowCount,
				)
			}
			assertTranslationOutputReadinessHasNoForbiddenText(t, result)
		})
	}
}

func TestSCN_TOA_010_DiffPreviewDoesNotResolveProviderSecretOrExposeRawPayload(t *testing.T) {
	fixture := newTranslationOutputReviewReadinessAPIFixture(t)
	controller := fixture.controller()

	result, err := controller.GetTranslationOutputDiffPreview(
		controllerwails.GetTranslationOutputDiffPreviewRequestDTO{JobID: 101, ArtifactID: 901},
	)
	if err != nil {
		t.Fatalf("SCN-TOA-010 expected redacted diff preview: %v", err)
	}
	if fixture.store.providerCallCount != 0 || fixture.store.secretResolutionCallCount != 0 {
		t.Fatalf(
			"SCN-TOA-010 expected no provider or secret resolution call, got provider=%d secret=%d",
			fixture.store.providerCallCount,
			fixture.store.secretResolutionCallCount,
		)
	}
	assertTranslationOutputReadinessHasNoForbiddenText(t, result)
}

type translationOutputReviewReadinessAPIFixture struct {
	t     *testing.T
	store *translationOutputReviewReadinessAPIStore
}

func newTranslationOutputReviewReadinessAPIFixture(t *testing.T) *translationOutputReviewReadinessAPIFixture {
	t.Helper()

	now := time.Date(2026, 5, 2, 10, 0, 0, 0, time.UTC)
	store := &translationOutputReviewReadinessAPIStore{
		now:                      now,
		translationArtifactCount: 0,
		xtranslatorRowCount:      0,
		xeditByID: map[int64]repository.XEditExtractedData{
			301: {
				ID:                301,
				SourceFilePath:    "Data/SecretPlugin.esp",
				SourceContentHash: "sha256:source-plugin",
				SourceTool:        "xedit",
				TargetPluginName:  "SecretPlugin.esp",
				TargetPluginType:  "esp",
				RecordCount:       2,
				ImportedAt:        now,
			},
		},
		translationRecordsByID: map[int64]repository.TranslationRecord{
			601: {ID: 601, XEditExtractedDataID: 301, FormID: "00000001", EditorID: "TranslatedNPC", RecordType: "NPC_"},
			602: {ID: 602, XEditExtractedDataID: 301, FormID: "00000002", EditorID: "CachedNPC", RecordType: "NPC_"},
			603: {ID: 603, XEditExtractedDataID: 301, FormID: "00000003", EditorID: "RiskyRace", RecordType: "RACE"},
		},
		translationFieldsByID: map[int64]repository.TranslationField{
			701: {ID: 701, TranslationRecordID: 601, SubrecordType: "FULL", SourceText: "Hello", FieldOrder: 1},
			702: {ID: 702, TranslationRecordID: 602, SubrecordType: "FULL", SourceText: "Cached source", FieldOrder: 2},
			703: {ID: 703, TranslationRecordID: 603, SubrecordType: "FULL", SourceText: " Leading space risk", FieldOrder: 3},
		},
		jobs: []repository.TranslationJob{
			{
				ID:                   101,
				XEditExtractedDataID: 301,
				JobName:              "completed output ready",
				State:                "completed",
				ProgressPercent:      100,
				CreatedAt:            now,
				StartedAt:            &now,
				FinishedAt:           &now,
			},
			{
				ID:                   202,
				XEditExtractedDataID: 301,
				JobName:              "body phase still running",
				State:                "running",
				ProgressPercent:      70,
				CreatedAt:            now,
				StartedAt:            &now,
			},
			{
				ID:                   303,
				XEditExtractedDataID: 301,
				JobName:              "zero target completed",
				State:                "completed",
				ProgressPercent:      100,
				CreatedAt:            now,
				StartedAt:            &now,
				FinishedAt:           &now,
			},
			{
				ID:                   404,
				XEditExtractedDataID: 301,
				JobName:              "canceled body output",
				State:                "canceled",
				ProgressPercent:      100,
				CreatedAt:            now,
				StartedAt:            &now,
				FinishedAt:           &now,
			},
		},
		phaseRunsByJobID: map[int64][]repository.JobPhaseRun{
			101: {
				translationOutputReviewReadinessBodyRun(201, 101, "completed", 2, "sha256:body-input", now),
			},
			202: {
				translationOutputReviewReadinessBodyRun(202, 202, "running", 2, "sha256:running-input", now),
			},
			303: {
				translationOutputReviewReadinessBodyRun(203, 303, "completed", 0, "sha256:zero-target-input", now),
			},
			404: {
				translationOutputReviewReadinessBodyRun(204, 404, "canceled", 2, "sha256:canceled-input", now),
			},
		},
		outputsByJobID: map[int64][]repository.JobTranslationField{
			101: {
				translationOutputReviewReadinessOutputField(801, 101, 701, "translated", now),
				translationOutputReviewReadinessOutputField(802, 101, 702, "cached", now),
			},
			202: {
				translationOutputReviewReadinessOutputField(803, 202, 701, "translated", now),
			},
			404: {
				translationOutputReviewReadinessOutputField(804, 404, 701, "translated", now),
			},
		},
		artifactsByJobID: make(map[int64]repository.TranslationArtifact),
		artifactsByID:    make(map[int64]repository.TranslationArtifact),
		rowsByArtifactID: make(map[int64][]repository.XTranslatorOutputRow),
	}
	return &translationOutputReviewReadinessAPIFixture{t: t, store: store}
}

func (fixture *translationOutputReviewReadinessAPIFixture) addCompletedJobWithOutput(
	jobID int64,
	outputID int64,
	fieldID int64,
	status string,
) {
	fixture.t.Helper()
	now := fixture.store.now
	fixture.store.jobs = append(fixture.store.jobs, repository.TranslationJob{
		ID:                   jobID,
		XEditExtractedDataID: 301,
		JobName:              "diff preview fixture",
		State:                "completed",
		ProgressPercent:      100,
		CreatedAt:            now,
		StartedAt:            &now,
		FinishedAt:           &now,
	})
	fixture.store.phaseRunsByJobID[jobID] = []repository.JobPhaseRun{
		translationOutputReviewReadinessBodyRun(jobID+100, jobID, "completed", 1, "sha256:diff-input", now),
	}
	fixture.store.outputsByJobID[jobID] = []repository.JobTranslationField{
		translationOutputReviewReadinessOutputField(outputID, jobID, fieldID, status, now),
	}
}

func (fixture *translationOutputReviewReadinessAPIFixture) controller() *controllerwails.TranslationOutputArtifactController {
	fixture.t.Helper()
	outputService := service.NewTranslationOutputArtifactService(fixture.store, fixture.store, fixture.store)
	outputUsecase := usecase.NewTranslationOutputArtifactUsecase(outputService)
	return controllerwails.NewTranslationOutputArtifactController(outputUsecase)
}

func (fixture *translationOutputReviewReadinessAPIFixture) controllerWithPersistence() *controllerwails.TranslationOutputArtifactController {
	fixture.t.Helper()
	outputService := service.NewTranslationOutputArtifactService(
		fixture.store,
		fixture.store,
		fixture.store,
	).WithArtifactPersistence(fixture.store, fixture.store)
	outputUsecase := usecase.NewTranslationOutputArtifactUsecase(outputService)
	return controllerwails.NewTranslationOutputArtifactController(outputUsecase)
}

func (fixture *translationOutputReviewReadinessAPIFixture) controllerWithPersistenceAndFileWriter(
	writer *translationOutputReviewReadinessRecordingFileWriter,
) *controllerwails.TranslationOutputArtifactController {
	fixture.t.Helper()
	outputService := service.NewTranslationOutputArtifactService(
		fixture.store,
		fixture.store,
		fixture.store,
	).WithArtifactPersistence(fixture.store, fixture.store).WithFileWriter(writer)
	outputUsecase := usecase.NewTranslationOutputArtifactUsecase(outputService)
	return controllerwails.NewTranslationOutputArtifactController(outputUsecase)
}

type translationOutputReviewReadinessAPIStore struct {
	now                       time.Time
	jobs                      []repository.TranslationJob
	phaseRunsByJobID          map[int64][]repository.JobPhaseRun
	outputsByJobID            map[int64][]repository.JobTranslationField
	xeditByID                 map[int64]repository.XEditExtractedData
	translationRecordsByID    map[int64]repository.TranslationRecord
	translationFieldsByID     map[int64]repository.TranslationField
	translationArtifactCount  int
	xtranslatorRowCount       int
	providerCallCount         int
	secretResolutionCallCount int
	failArtifactSave          bool
	artifactsByJobID          map[int64]repository.TranslationArtifact
	artifactsByID             map[int64]repository.TranslationArtifact
	rowsByArtifactID          map[int64][]repository.XTranslatorOutputRow
}

type translationOutputReviewReadinessRecordingFileWriter struct {
	writeCount            int
	writeTemporaryCount   int
	publishTemporaryCount int
	removeCount           int
	lastPath              string
	lastBytes             []byte
	tempPath              string
	publishErr            error
}

func translationOutputReviewReadinessBodyRun(
	id int64,
	jobID int64,
	state string,
	targetCount int,
	inputDigest string,
	now time.Time,
) repository.JobPhaseRun {
	return repository.JobPhaseRun{
		ID:                  id,
		TranslationJobID:    jobID,
		PhaseType:           "body_translation",
		State:               state,
		ExecutionOrder:      2,
		ProgressPercent:     100,
		SnapshotFieldCount:  targetCount,
		ProviderTargetCount: targetCount,
		AIProvider:          "fake-provider-that-must-not-run",
		ModelName:           "fixture-model",
		ExecutionMode:       "single_request",
		CredentialRef:       "sk-live-secret-must-not-resolve",
		InstructionKind:     "body_translation",
		InputSnapshotDigest: inputDigest,
		LatestExternalRunID: "providerRawResponse:{\"token\":\"secret\"}",
		StartedAt:           &now,
		FinishedAt:          &now,
	}
}

func translationOutputReviewReadinessOutputField(
	id int64,
	jobID int64,
	fieldID int64,
	outputStatus string,
	now time.Time,
) repository.JobTranslationField {
	return repository.JobTranslationField{
		ID:                 id,
		TranslationJobID:   jobID,
		TranslationFieldID: fieldID,
		TranslatedText:     "translated fixture must not expose fullDestText",
		OutputStatus:       outputStatus,
		RetryCount:         0,
		UpdatedAt:          now,
	}
}

func (store *translationOutputReviewReadinessAPIStore) WithTransaction(
	ctx context.Context,
	fn func(ctx context.Context) error,
) error {
	artifactsByJobID := make(map[int64]repository.TranslationArtifact, len(store.artifactsByJobID))
	for jobID, artifact := range store.artifactsByJobID {
		artifactsByJobID[jobID] = artifact
	}
	artifactsByID := make(map[int64]repository.TranslationArtifact, len(store.artifactsByID))
	for artifactID, artifact := range store.artifactsByID {
		artifactsByID[artifactID] = artifact
	}
	rowsByArtifactID := make(map[int64][]repository.XTranslatorOutputRow, len(store.rowsByArtifactID))
	for artifactID, rows := range store.rowsByArtifactID {
		rowsByArtifactID[artifactID] = append([]repository.XTranslatorOutputRow(nil), rows...)
	}
	translationArtifactCount := store.translationArtifactCount
	xtranslatorRowCount := store.xtranslatorRowCount

	if err := fn(ctx); err != nil {
		store.artifactsByJobID = artifactsByJobID
		store.artifactsByID = artifactsByID
		store.rowsByArtifactID = rowsByArtifactID
		store.translationArtifactCount = translationArtifactCount
		store.xtranslatorRowCount = xtranslatorRowCount
		return err
	}
	return nil
}

func (store *translationOutputReviewReadinessAPIStore) seedTranslationArtifact(
	jobID int64,
	artifactID int64,
) {
	artifact := repository.TranslationArtifact{
		ID:               artifactID,
		TranslationJobID: jobID,
		ArtifactFormat:   "xtranslator_xml",
		TargetGame:       "skyrim_se",
		FilePath:         "/tmp/existing-output.xml",
		Status:           "success",
		GeneratedAt:      &store.now,
	}
	store.artifactsByJobID[jobID] = artifact
	store.artifactsByID[artifactID] = artifact
	store.translationArtifactCount = len(store.artifactsByID)
}

func (store *translationOutputReviewReadinessAPIStore) GetTranslationArtifactByID(
	_ context.Context,
	id int64,
) (repository.TranslationArtifact, error) {
	if artifact, ok := store.artifactsByID[id]; ok {
		return artifact, nil
	}
	return repository.TranslationArtifact{}, repository.ErrNotFound
}

func (store *translationOutputReviewReadinessAPIStore) GetTranslationArtifactByJobID(
	_ context.Context,
	jobID int64,
) (repository.TranslationArtifact, error) {
	if artifact, ok := store.artifactsByJobID[jobID]; ok {
		return artifact, nil
	}
	return repository.TranslationArtifact{}, repository.ErrNotFound
}

func (store *translationOutputReviewReadinessAPIStore) UpsertTranslationArtifact(
	_ context.Context,
	draft repository.TranslationArtifactDraft,
) (repository.TranslationArtifact, error) {
	if store.failArtifactSave {
		return repository.TranslationArtifact{}, errors.New("injected artifact save failure")
	}
	artifact, ok := store.artifactsByJobID[draft.TranslationJobID]
	if !ok {
		artifact.ID = int64(len(store.artifactsByID) + 1)
	}
	artifact.TranslationJobID = draft.TranslationJobID
	artifact.ArtifactFormat = draft.ArtifactFormat
	artifact.TargetGame = draft.TargetGame
	artifact.FilePath = draft.FilePath
	artifact.Status = draft.Status
	artifact.GeneratedAt = draft.GeneratedAt
	store.artifactsByJobID[draft.TranslationJobID] = artifact
	store.artifactsByID[artifact.ID] = artifact
	store.translationArtifactCount = len(store.artifactsByID)
	return artifact, nil
}

func (store *translationOutputReviewReadinessAPIStore) ReplaceXTranslatorOutputRows(
	_ context.Context,
	translationArtifactID int64,
	drafts []repository.XTranslatorOutputRowDraft,
) ([]repository.XTranslatorOutputRow, error) {
	rows := make([]repository.XTranslatorOutputRow, 0, len(drafts))
	for index, draft := range drafts {
		rows = append(rows, repository.XTranslatorOutputRow{
			ID:                    int64(index + 1),
			TranslationArtifactID: translationArtifactID,
			JobTranslationFieldID: draft.JobTranslationFieldID,
			EDID:                  draft.EDID,
			REC:                   draft.REC,
			FIELD:                 draft.FIELD,
			FORMID:                draft.FORMID,
			Source:                draft.Source,
			Dest:                  draft.Dest,
			Status:                draft.Status,
		})
	}
	store.rowsByArtifactID[translationArtifactID] = rows
	total := 0
	for _, artifactRows := range store.rowsByArtifactID {
		total += len(artifactRows)
	}
	store.xtranslatorRowCount = total
	return rows, nil
}

func (store *translationOutputReviewReadinessAPIStore) CountXTranslatorOutputRowsByArtifactID(
	_ context.Context,
	translationArtifactID int64,
) (int, error) {
	return len(store.rowsByArtifactID[translationArtifactID]), nil
}

func (writer *translationOutputReviewReadinessRecordingFileWriter) WriteFile(
	path string,
	payload []byte,
) error {
	writer.writeCount++
	writer.lastPath = path
	writer.lastBytes = append([]byte(nil), payload...)
	return nil
}

func (writer *translationOutputReviewReadinessRecordingFileWriter) WriteTemporaryFile(
	path string,
	payload []byte,
) (string, error) {
	writer.writeTemporaryCount++
	writer.lastPath = path
	writer.lastBytes = append([]byte(nil), payload...)
	writer.tempPath = path + ".tmp"
	if err := os.WriteFile(writer.tempPath, payload, 0o600); err != nil {
		return "", errors.Join(errors.New("write temporary output artifact file"), err)
	}
	return writer.tempPath, nil
}

func (writer *translationOutputReviewReadinessRecordingFileWriter) PublishTemporaryFile(
	_ string,
	_ string,
) error {
	writer.publishTemporaryCount++
	if writer.publishErr != nil {
		return writer.publishErr
	}
	return nil
}

func (writer *translationOutputReviewReadinessRecordingFileWriter) RemoveFile(path string) error {
	writer.removeCount++
	if err := os.Remove(path); err != nil && !errors.Is(err, os.ErrNotExist) {
		return errors.Join(errors.New("remove output artifact file"), err)
	}
	return nil
}

func (store *translationOutputReviewReadinessAPIStore) GetTranslationJobByID(
	_ context.Context,
	id int64,
) (repository.TranslationJob, error) {
	for _, job := range store.jobs {
		if job.ID == id {
			return job, nil
		}
	}
	return repository.TranslationJob{}, repository.ErrNotFound
}

func (store *translationOutputReviewReadinessAPIStore) ListJobPhaseRunsByJobID(
	_ context.Context,
	jobID int64,
) ([]repository.JobPhaseRun, error) {
	return append([]repository.JobPhaseRun(nil), store.phaseRunsByJobID[jobID]...), nil
}

func (store *translationOutputReviewReadinessAPIStore) FindJobPhaseRun(
	_ context.Context,
	translationJobID int64,
	phaseType string,
) (repository.JobPhaseRun, error) {
	for _, run := range store.phaseRunsByJobID[translationJobID] {
		if run.PhaseType == phaseType {
			return run, nil
		}
	}
	return repository.JobPhaseRun{}, repository.ErrNotFound
}

func (store *translationOutputReviewReadinessAPIStore) ListJobTranslationFieldsByJobID(
	_ context.Context,
	jobID int64,
) ([]repository.JobTranslationField, error) {
	return append([]repository.JobTranslationField(nil), store.outputsByJobID[jobID]...), nil
}

func (store *translationOutputReviewReadinessAPIStore) GetXEditExtractedDataByID(
	_ context.Context,
	id int64,
) (repository.XEditExtractedData, error) {
	if xedit, ok := store.xeditByID[id]; ok {
		return xedit, nil
	}
	return repository.XEditExtractedData{}, repository.ErrNotFound
}

func (store *translationOutputReviewReadinessAPIStore) DeleteXEditExtractedDataByID(
	_ context.Context,
	id int64,
) error {
	if _, ok := store.xeditByID[id]; !ok {
		return repository.ErrNotFound
	}
	delete(store.xeditByID, id)
	return nil
}

func (store *translationOutputReviewReadinessAPIStore) ResolveTranslationOutputSecret(
	context.Context,
	string,
) (string, error) {
	store.secretResolutionCallCount++
	return "sk-live-secret", nil
}

func (store *translationOutputReviewReadinessAPIStore) CallTranslationOutputProvider(context.Context) error {
	store.providerCallCount++
	return nil
}

func (store *translationOutputReviewReadinessAPIStore) CreateTranslationJob(
	context.Context,
	repository.TranslationJobDraft,
) (repository.TranslationJob, error) {
	return repository.TranslationJob{}, repository.ErrNotFound
}

func (store *translationOutputReviewReadinessAPIStore) UpdateTranslationJob(
	context.Context,
	int64,
	repository.TranslationJobUpdateDraft,
) (repository.TranslationJob, error) {
	return repository.TranslationJob{}, repository.ErrNotFound
}

func (store *translationOutputReviewReadinessAPIStore) CreateJobPhaseRun(
	context.Context,
	repository.JobPhaseRunDraft,
) (repository.JobPhaseRun, error) {
	return repository.JobPhaseRun{}, repository.ErrNotFound
}

func (store *translationOutputReviewReadinessAPIStore) GetJobPhaseRunByID(
	context.Context,
	int64,
) (repository.JobPhaseRun, error) {
	return repository.JobPhaseRun{}, repository.ErrNotFound
}

func (store *translationOutputReviewReadinessAPIStore) UpdateJobPhaseRun(
	context.Context,
	int64,
	repository.JobPhaseRunUpdateDraft,
) (repository.JobPhaseRun, error) {
	return repository.JobPhaseRun{}, repository.ErrNotFound
}

func (store *translationOutputReviewReadinessAPIStore) CreatePhaseRunTranslationField(
	context.Context,
	repository.PhaseRunTranslationFieldDraft,
) (repository.PhaseRunTranslationField, error) {
	return repository.PhaseRunTranslationField{}, repository.ErrNotFound
}

func (store *translationOutputReviewReadinessAPIStore) ListPhaseRunTranslationFieldsByPhaseRunID(
	context.Context,
	int64,
) ([]repository.PhaseRunTranslationField, error) {
	return nil, nil
}

func (store *translationOutputReviewReadinessAPIStore) CreatePhaseRunPersona(
	context.Context,
	repository.PhaseRunPersonaDraft,
) (repository.PhaseRunPersona, error) {
	return repository.PhaseRunPersona{}, repository.ErrNotFound
}

func (store *translationOutputReviewReadinessAPIStore) ListPhaseRunPersonasByPhaseRunID(
	context.Context,
	int64,
) ([]repository.PhaseRunPersona, error) {
	return nil, nil
}

func (store *translationOutputReviewReadinessAPIStore) CreatePhaseRunDictionaryEntry(
	context.Context,
	repository.PhaseRunDictionaryEntryDraft,
) (repository.PhaseRunDictionaryEntry, error) {
	return repository.PhaseRunDictionaryEntry{}, repository.ErrNotFound
}

func (store *translationOutputReviewReadinessAPIStore) ListPhaseRunDictionaryEntriesByPhaseRunID(
	context.Context,
	int64,
) ([]repository.PhaseRunDictionaryEntry, error) {
	return nil, nil
}

func (store *translationOutputReviewReadinessAPIStore) CreateJobTranslationField(
	context.Context,
	repository.JobTranslationFieldDraft,
) (repository.JobTranslationField, error) {
	return repository.JobTranslationField{}, repository.ErrNotFound
}

func (store *translationOutputReviewReadinessAPIStore) GetJobTranslationFieldByID(
	context.Context,
	int64,
) (repository.JobTranslationField, error) {
	return repository.JobTranslationField{}, repository.ErrNotFound
}

func (store *translationOutputReviewReadinessAPIStore) UpdateJobTranslationField(
	context.Context,
	int64,
	repository.JobTranslationFieldUpdateDraft,
) (repository.JobTranslationField, error) {
	return repository.JobTranslationField{}, repository.ErrNotFound
}

func (store *translationOutputReviewReadinessAPIStore) CreateXEditExtractedData(
	context.Context,
	repository.XEditExtractedDataDraft,
) (repository.XEditExtractedData, error) {
	return repository.XEditExtractedData{}, repository.ErrNotFound
}

func (store *translationOutputReviewReadinessAPIStore) CreateTranslationRecord(
	context.Context,
	repository.TranslationRecordDraft,
) (repository.TranslationRecord, error) {
	return repository.TranslationRecord{}, repository.ErrNotFound
}

func (store *translationOutputReviewReadinessAPIStore) GetTranslationRecordByID(
	_ context.Context,
	id int64,
) (repository.TranslationRecord, error) {
	if record, ok := store.translationRecordsByID[id]; ok {
		return record, nil
	}
	return repository.TranslationRecord{}, repository.ErrNotFound
}

func (store *translationOutputReviewReadinessAPIStore) ListTranslationRecordsByXEditID(
	context.Context,
	int64,
) ([]repository.TranslationRecord, error) {
	return nil, nil
}

func (store *translationOutputReviewReadinessAPIStore) UpsertNpcProfile(
	context.Context,
	repository.NpcProfileDraft,
) (repository.NpcProfile, error) {
	return repository.NpcProfile{}, repository.ErrNotFound
}

func (store *translationOutputReviewReadinessAPIStore) GetNpcProfileByID(
	context.Context,
	int64,
) (repository.NpcProfile, error) {
	return repository.NpcProfile{}, repository.ErrNotFound
}

func (store *translationOutputReviewReadinessAPIStore) CreateNpcRecord(
	context.Context,
	repository.NpcRecordDraft,
) (repository.NpcRecord, error) {
	return repository.NpcRecord{}, repository.ErrNotFound
}

func (store *translationOutputReviewReadinessAPIStore) GetNpcRecordByTranslationRecordID(
	context.Context,
	int64,
) (repository.NpcRecord, error) {
	return repository.NpcRecord{}, repository.ErrNotFound
}

func (store *translationOutputReviewReadinessAPIStore) CreateTranslationField(
	context.Context,
	repository.TranslationFieldDraft,
) (repository.TranslationField, error) {
	return repository.TranslationField{}, repository.ErrNotFound
}

func (store *translationOutputReviewReadinessAPIStore) GetTranslationFieldByID(
	_ context.Context,
	id int64,
) (repository.TranslationField, error) {
	if field, ok := store.translationFieldsByID[id]; ok {
		return field, nil
	}
	return repository.TranslationField{}, repository.ErrNotFound
}

func (store *translationOutputReviewReadinessAPIStore) ListTranslationFieldsByTranslationRecordID(
	context.Context,
	int64,
) ([]repository.TranslationField, error) {
	return nil, nil
}

func (store *translationOutputReviewReadinessAPIStore) CreateTranslationFieldRecordReference(
	context.Context,
	repository.TranslationFieldRecordReferenceDraft,
) (repository.TranslationFieldRecordReference, error) {
	return repository.TranslationFieldRecordReference{}, repository.ErrNotFound
}

func (store *translationOutputReviewReadinessAPIStore) ListTranslationFieldRecordReferencesByFieldID(
	context.Context,
	int64,
) ([]repository.TranslationFieldRecordReference, error) {
	return nil, nil
}

func assertTranslationOutputReadinessHasNoForbiddenText(t *testing.T, value any) {
	t.Helper()

	encoded, err := json.Marshal(value)
	if err != nil {
		t.Fatalf("expected response to marshal: %v", err)
	}
	serialized := string(encoded)
	for _, forbidden := range []string{
		"sk-live-secret",
		"api_key",
		"apiKey",
		"token",
		"authorizationHeader",
		"providerRawRequest",
		"providerRawResponse",
		"decrypted",
		"fullSourceText",
		"fullDestText",
		"translated fixture must not expose fullDestText",
	} {
		if strings.Contains(serialized, forbidden) {
			t.Fatalf("expected Output Review response not to expose %q: %s", forbidden, serialized)
		}
	}
}

func assertTranslationOutputGeneratedXML(t *testing.T, path string, wantRoot string, wantStringCount int) {
	t.Helper()

	content, err := os.ReadFile(path) // #nosec G304 -- path is created by the scenario test TempDir fixture.
	if err != nil {
		t.Fatalf("expected generated XML file to be readable: %v", err)
	}
	var document translationOutputArtifactXMLDocument
	if err := xml.Unmarshal(content, &document); err != nil {
		t.Fatalf("expected generated XML to parse as UTF-8 XML: %v", err)
	}
	if document.XMLName.Local != wantRoot {
		t.Fatalf("expected XML root %q, got %q", wantRoot, document.XMLName.Local)
	}
	if len(document.Strings) != wantStringCount {
		t.Fatalf("expected %d String rows, got %#v", wantStringCount, document.Strings)
	}
}

type translationOutputArtifactXMLDocument struct {
	XMLName xml.Name
	Strings []translationOutputArtifactXMLString `xml:"String"`
}

type translationOutputArtifactXMLString struct {
	EDID   string `xml:"EDID"`
	REC    string `xml:"REC"`
	FIELD  string `xml:"FIELD"`
	FORMID string `xml:"FORMID"`
	Source string `xml:"Source"`
	Dest   string `xml:"Dest"`
	Status int    `xml:"Status"`
}

func findTranslationOutputReadinessCompletedJob(
	jobs []controllerwails.TranslationOutputCompletedJobSummaryDTO,
	jobID int64,
) (controllerwails.TranslationOutputCompletedJobSummaryDTO, bool) {
	for _, job := range jobs {
		if job.JobID == jobID {
			return job, true
		}
	}
	return controllerwails.TranslationOutputCompletedJobSummaryDTO{}, false
}

func findTranslationOutputDiffPreviewRow(
	rows []controllerwails.TranslationOutputDiffPreviewRowDTO,
	fieldID int64,
) (controllerwails.TranslationOutputDiffPreviewRowDTO, bool) {
	for _, row := range rows {
		if row.FieldID == fieldID {
			return row, true
		}
	}
	return controllerwails.TranslationOutputDiffPreviewRowDTO{}, false
}

var _ repository.Transactor = (*translationOutputReviewReadinessAPIStore)(nil)
var _ repository.JobLifecycleRepository = (*translationOutputReviewReadinessAPIStore)(nil)
var _ repository.JobOutputRepository = (*translationOutputReviewReadinessAPIStore)(nil)
var _ repository.TranslationSourceRepository = (*translationOutputReviewReadinessAPIStore)(nil)
