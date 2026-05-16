package apitest

import (
	"context"
	"strings"
	"testing"
	"time"

	controllerwails "aitranslationenginejp/internal/controller/wails"
	"aitranslationenginejp/internal/repository"
	"aitranslationenginejp/internal/service"
	"aitranslationenginejp/internal/usecase"
)

func TestSCN_BTP_007_RecoverableProviderFailureIsNotPublishedAsSuccess(t *testing.T) {
	fixture := newBodyTranslationAPIFixture(t, bodyTranslationAPIFixtureOptions{
		JobState:     "running",
		BodyRunState: "recoverable_failed",
		LatestError:  "provider_failure",
		Outputs: []repository.JobTranslationField{
			bodyTranslationAPIOutputField(801, 701, "translated"),
		},
	})
	controller := fixture.controller()

	result, err := controller.GetBodyTranslationPhaseSummary(
		controllerwails.GetBodyTranslationPhaseSummaryRequestDTO{JobID: fixture.job.ID},
	)
	if err != nil {
		t.Fatalf("SCN-BTP-007 public summary returned error: %v", err)
	}
	if result.PhaseState == "completed" || result.OutputReadiness.Ready {
		t.Fatalf("SCN-BTP-007 failed provider result must not be completed or output-ready: %#v", result)
	}
	if result.ErrorSummary == nil {
		t.Fatalf("SCN-BTP-007 expected provider failure error summary, got nil")
	}
	if result.ErrorSummary.ErrorKind != "provider_failure" || !result.ErrorSummary.Retryable {
		t.Fatalf("SCN-BTP-007 expected retryable provider_failure, got %#v", result.ErrorSummary)
	}
	if result.ResultSummary == nil || result.ResultSummary.TranslatedCount != 1 || result.ResultSummary.FailedCount != 1 {
		t.Fatalf("SCN-BTP-007 expected one success retained and one failed target, got %#v", result.ResultSummary)
	}
}

func TestSCN_BTP_008_RetryReusesPhaseRunAndDoesNotDuplicatePublishedResults(t *testing.T) {
	fixture := newBodyTranslationAPIFixture(t, bodyTranslationAPIFixtureOptions{
		JobState:     "running",
		BodyRunState: "recoverable_failed",
		LatestError:  "invalid_provider_response",
		Outputs: []repository.JobTranslationField{
			bodyTranslationAPIOutputField(801, 701, "translated"),
		},
	})
	controller := fixture.controller()

	result, err := controller.RetryBodyTranslationPhase(
		controllerwails.RetryBodyTranslationPhaseRequestDTO{
			JobID:      fixture.job.ID,
			PhaseRunID: fixture.bodyRun.ID,
		},
	)
	if err != nil {
		t.Fatalf("SCN-BTP-008 public retry returned error: %v", err)
	}
	if result.PhaseRunID == nil || *result.PhaseRunID != fixture.bodyRun.ID {
		t.Fatalf("SCN-BTP-008 expected same phase run id %d, got %#v", fixture.bodyRun.ID, result.PhaseRunID)
	}
	if result.ResultSummary == nil || result.ResultSummary.TranslatedCount != 1 {
		t.Fatalf("SCN-BTP-008 expected retry to preserve one published success without duplication, got %#v", result.ResultSummary)
	}
	if result.Progress.ProcessedCount != 1 || result.Progress.TargetCount != 2 {
		t.Fatalf("SCN-BTP-008 expected retry progress to keep one processed field from two targets, got %#v", result.Progress)
	}
}

func TestSCN_BTP_008_RetryPersistsStateAndCompletesPendingFieldResults(t *testing.T) {
	fixture := newBodyTranslationAPIFixture(t, bodyTranslationAPIFixtureOptions{
		JobState:     "running",
		BodyRunState: "recoverable_failed",
		LatestError:  "provider_failure",
		AIProvider:   "xai",
		Provider:     bodyTranslationAPISuccessProvider{},
	})
	controller := fixture.controller()

	result, err := controller.RetryBodyTranslationPhase(
		controllerwails.RetryBodyTranslationPhaseRequestDTO{
			JobID:      fixture.job.ID,
			PhaseRunID: fixture.bodyRun.ID,
		},
	)
	if err != nil {
		t.Fatalf("SCN-BTP-008 public retry returned error: %v", err)
	}
	if result.PhaseState != "completed" {
		t.Fatalf("SCN-BTP-008 expected retry to complete pending fields, got %#v", result)
	}
	if result.ResultSummary == nil || result.ResultSummary.OutputCount != 2 || result.ResultSummary.OutputReadyCount != 2 {
		t.Fatalf("SCN-BTP-008 expected two ready outputs after retry, got %#v", result.ResultSummary)
	}
}

func TestSCN_BTP_008_RetryRejectsInputSnapshotDriftWithoutWritingResults(t *testing.T) {
	fixture := newBodyTranslationAPIFixture(t, bodyTranslationAPIFixtureOptions{
		JobState:     "running",
		BodyRunState: "recoverable_failed",
		LatestError:  "provider_failure",
		AIProvider:   "xai",
		Provider:     bodyTranslationAPISuccessProvider{},
	})
	fixture.store.phaseRuns[1].InputSnapshotDigest = "sha256:started-before-drift"
	controller := fixture.controller()

	result, err := controller.RetryBodyTranslationPhase(
		controllerwails.RetryBodyTranslationPhaseRequestDTO{
			JobID:      fixture.job.ID,
			PhaseRunID: fixture.bodyRun.ID,
		},
	)
	if err != nil {
		t.Fatalf("SCN-BTP-008 expected drift rejection payload without Wails error, got %v", err)
	}
	if result.ErrorSummary == nil || result.ErrorSummary.ErrorKind != "input_snapshot_failed" {
		t.Fatalf("SCN-BTP-008 expected input snapshot drift rejection, got %#v", result)
	}
	if len(fixture.store.outputs) != 0 {
		t.Fatalf("SCN-BTP-008 expected no field result writes after drift, got %#v", fixture.store.outputs)
	}
}

func TestSCN_TJSM_003_RetryRecoverableFailedReusesPhaseRunAndDoesNotDuplicateExistingResult(t *testing.T) {
	fixture := newBodyTranslationAPIFixture(t, bodyTranslationAPIFixtureOptions{
		JobState:     "running",
		BodyRunState: "recoverable_failed",
		LatestError:  "provider_failure",
		AIProvider:   "xai",
		Provider:     bodyTranslationAPISuccessProvider{},
		Outputs: []repository.JobTranslationField{
			bodyTranslationAPIOutputField(801, 701, "translated"),
		},
	})
	controller := fixture.controller()

	result, err := controller.RetryBodyTranslationPhase(
		controllerwails.RetryBodyTranslationPhaseRequestDTO{
			JobID:      fixture.job.ID,
			PhaseRunID: fixture.bodyRun.ID,
		},
	)
	if err != nil {
		t.Fatalf("SCN-TJSM-003 public retry returned error: %v", err)
	}
	if result.PhaseRunID == nil || *result.PhaseRunID != fixture.bodyRun.ID {
		t.Fatalf("SCN-TJSM-003 expected retry to reuse phase run %d, got %#v", fixture.bodyRun.ID, result.PhaseRunID)
	}
	if len(fixture.store.phaseRuns) != 2 {
		t.Fatalf("SCN-TJSM-003 expected retry not to create phase run, got %#v", fixture.store.phaseRuns)
	}
	if bodyTranslationAPIOutputCountForField(fixture.store.outputs, 701) != 1 {
		t.Fatalf("SCN-TJSM-003 expected retry not to duplicate existing field result, got %#v", fixture.store.outputs)
	}
}

func TestSCN_TJSM_003_ResumePausedReusesPhaseRun(t *testing.T) {
	fixture := newBodyTranslationAPIFixture(t, bodyTranslationAPIFixtureOptions{
		JobState:     "paused",
		BodyRunState: "paused",
	})
	controller := fixture.controller()

	result, err := controller.ResumeBodyTranslationPhase(
		controllerwails.ResumeBodyTranslationPhaseRequestDTO{
			JobID:      fixture.job.ID,
			PhaseRunID: fixture.bodyRun.ID,
		},
	)
	if err != nil {
		t.Fatalf("SCN-TJSM-003 public resume returned error: %v", err)
	}
	if result.PhaseRunID == nil || *result.PhaseRunID != fixture.bodyRun.ID {
		t.Fatalf("SCN-TJSM-003 expected resume to reuse phase run %d, got %#v", fixture.bodyRun.ID, result.PhaseRunID)
	}
	if result.PhaseState != "running" {
		t.Fatalf("SCN-TJSM-003 expected paused resume to continue as running, got %#v", result)
	}
	if len(fixture.store.phaseRuns) != 2 {
		t.Fatalf("SCN-TJSM-003 expected resume not to create phase run, got %#v", fixture.store.phaseRuns)
	}
}

func TestSCN_TJSM_003_ResumeRecoverableFailedIsRejectedWithoutMutation(t *testing.T) {
	fixture := newBodyTranslationAPIFixture(t, bodyTranslationAPIFixtureOptions{
		JobState:     "running",
		BodyRunState: "recoverable_failed",
		LatestError:  "provider_failure",
		Outputs: []repository.JobTranslationField{
			bodyTranslationAPIOutputField(801, 701, "translated"),
		},
	})
	controller := fixture.controller()

	result, err := controller.ResumeBodyTranslationPhase(
		controllerwails.ResumeBodyTranslationPhaseRequestDTO{
			JobID:      fixture.job.ID,
			PhaseRunID: fixture.bodyRun.ID,
		},
	)
	if err != nil {
		t.Fatalf("SCN-TJSM-003 expected RecoverableFailed resume rejection payload without Wails error, got %v", err)
	}
	assertBodyTranslationAPIRejectedCommand(t, "SCN-TJSM-003", result, "persona_phase_incomplete", "phase_not_paused")
	if result.PhaseState != "recoverable_failed" {
		t.Fatalf("SCN-TJSM-003 expected rejected resume to keep recoverable_failed, got %#v", result)
	}
	if len(fixture.store.phaseRuns) != 2 || len(fixture.store.outputs) != 1 || len(fixture.store.phaseLinks) != 1 {
		t.Fatalf("SCN-TJSM-003 expected rejected resume not to mutate rows, phaseRuns=%#v outputs=%#v links=%#v", fixture.store.phaseRuns, fixture.store.outputs, fixture.store.phaseLinks)
	}
}

func TestSCN_TJSM_003_StartResendWithActivePhaseRunIsRejectedWithoutNewRun(t *testing.T) {
	fixture := newBodyTranslationAPIFixture(t, bodyTranslationAPIFixtureOptions{
		JobState:     "running",
		BodyRunState: "running",
		Outputs: []repository.JobTranslationField{
			bodyTranslationAPIOutputField(801, 701, "translated"),
		},
	})
	controller := fixture.controller()

	result, err := controller.StartBodyTranslationPhase(
		controllerwails.StartBodyTranslationPhaseRequestDTO{JobID: fixture.job.ID},
	)
	if err != nil {
		t.Fatalf("SCN-TJSM-003 expected active start resend rejection payload without Wails error, got %v", err)
	}
	assertBodyTranslationAPIRejectedCommand(t, "SCN-TJSM-003", result, "active_phase_exists", "active_phase_exists")
	if result.PhaseRunID == nil || *result.PhaseRunID != fixture.bodyRun.ID {
		t.Fatalf("SCN-TJSM-003 expected active start resend to keep phase run %d, got %#v", fixture.bodyRun.ID, result.PhaseRunID)
	}
	if len(fixture.store.phaseRuns) != 2 || len(fixture.store.outputs) != 1 || len(fixture.store.phaseLinks) != 1 {
		t.Fatalf("SCN-TJSM-003 expected active start resend not to mutate rows, phaseRuns=%#v outputs=%#v links=%#v", fixture.store.phaseRuns, fixture.store.outputs, fixture.store.phaseLinks)
	}
}

func TestSCN_TJSR_001_ReadyJobSummaryDoesNotCreateBodyPhaseRunAndStartCreatesNonPendingRun(t *testing.T) {
	fixture := newBodyTranslationAPIFixture(t, bodyTranslationAPIFixtureOptions{
		JobState:    "ready",
		OmitBodyRun: true,
	})
	controller := fixture.controller()

	summary, summaryErr := controller.GetBodyTranslationPhaseSummary(
		controllerwails.GetBodyTranslationPhaseSummaryRequestDTO{JobID: fixture.job.ID},
	)
	if summaryErr != nil {
		t.Fatalf("SCN-TJSR-001 public summary returned error: %v", summaryErr)
	}
	if summary.PhaseState == "pending" || summary.PhaseRunID != nil {
		t.Fatalf("SCN-TJSR-001 expected read-only Ready summary without canonical pending phase run, got %#v", summary)
	}
	if len(fixture.store.phaseRuns) != 1 {
		t.Fatalf("SCN-TJSR-001 expected Ready query not to create body phase run, got %#v", fixture.store.phaseRuns)
	}

	started, startErr := controller.StartBodyTranslationPhase(
		controllerwails.StartBodyTranslationPhaseRequestDTO{JobID: fixture.job.ID},
	)
	if startErr != nil {
		t.Fatalf("SCN-TJSR-001 public start returned error: %v", startErr)
	}
	if started.PhaseRunID == nil || started.PhaseState == "pending" {
		t.Fatalf("SCN-TJSR-001 expected start to create observable non-pending phase run, got %#v", started)
	}
	if len(fixture.store.phaseRuns) != 2 || fixture.store.phaseRuns[1].State == "pending" {
		t.Fatalf("SCN-TJSR-001 expected exactly one new non-pending body phase run, got %#v", fixture.store.phaseRuns)
	}
}

func TestSCN_TJSR_002_BodyPhaseActionEnablementUsesCommonStateRules(t *testing.T) {
	testCases := []struct {
		name       string
		jobState   string
		phaseState string
		canPause   bool
		canResume  bool
		canRetry   bool
		canCancel  bool
	}{
		{name: "running", jobState: "running", phaseState: "running", canPause: true},
		{name: "paused", jobState: "paused", phaseState: "paused", canResume: true, canCancel: true},
		{name: "recoverable failed", jobState: "recoverable_failed", phaseState: "recoverable_failed", canRetry: true},
		{name: "terminal job", jobState: "completed", phaseState: "paused"},
	}
	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			fixture := newBodyTranslationAPIFixture(t, bodyTranslationAPIFixtureOptions{
				JobState:     testCase.jobState,
				BodyRunState: testCase.phaseState,
			})
			controller := fixture.controller()

			summary, err := controller.GetBodyTranslationPhaseSummary(
				controllerwails.GetBodyTranslationPhaseSummaryRequestDTO{JobID: fixture.job.ID},
			)
			if err != nil {
				t.Fatalf("SCN-TJSR-002 public summary returned error: %v", err)
			}
			action := summary.ActionEnablement
			if action.CanPause != testCase.canPause ||
				action.CanResume != testCase.canResume ||
				action.CanRetry != testCase.canRetry ||
				action.CanCancel != testCase.canCancel {
				t.Fatalf("SCN-TJSR-002 expected common action rules for %s, got %#v", testCase.name, action)
			}
		})
	}
}

func TestSCN_TFN_011_LateResponseForOldPhaseRunIsRejectedWithoutResultOrArtifactMutation(t *testing.T) {
	fixture := newBodyTranslationAPIFixture(t, bodyTranslationAPIFixtureOptions{
		JobState:     "running",
		BodyRunState: "running",
		Outputs: []repository.JobTranslationField{
			bodyTranslationAPIOutputField(801, 701, "translated"),
		},
	})
	phaseService := service.NewBodyTranslationPhaseService(
		fixture.store,
		fixture.store,
		fixture.store,
		fixture.store,
		fixture.store,
	)

	result, err := phaseService.PersistBodyTranslationFieldResults(context.Background(), service.BodyTranslationFieldResultPersistenceRequest{
		TranslationJobID: fixture.job.ID,
		PhaseRunID:       fixture.bodyRun.ID + 999,
		TargetFields: []service.BodyTranslationFieldResultTarget{
			{
				TranslationFieldID:    702,
				FieldCorrelationKey:   "field:702",
				OutputStatusCandidate: "ready",
			},
		},
		ProviderResults: []service.BodyTranslationProviderResult{
			{
				FieldCorrelationKey: "field:702",
				TranslatedCandidate: &service.BodyTranslationTranslatedCandidate{
					FieldCorrelationKey: "field:702",
					RecordType:          "NPC_",
					FieldType:           "DESC",
					TranslatedText:      "late translated fixture",
				},
				ProtectionValidationTarget: &service.BodyTranslationProtectionValidationTarget{
					FieldCorrelationKey: "field:702",
					TranslatedText:      "late translated fixture",
				},
			},
		},
	})
	if err != nil {
		t.Fatalf("SCN-TFN-011 expected late response rejection without persistence error, got %v", err)
	}
	if result.ErrorSummary == nil || result.ErrorSummary.ErrorKind != "late_response_rejected" || !result.ErrorSummary.IsRedacted {
		t.Fatalf("SCN-TFN-011 expected redacted late_response_rejected summary, got %#v", result.ErrorSummary)
	}
	if len(fixture.store.outputs) != 1 || len(fixture.store.phaseLinks) != 1 {
		t.Fatalf("SCN-TFN-011 expected no result or artifact row mutation, outputs=%#v links=%#v", fixture.store.outputs, fixture.store.phaseLinks)
	}
	if len(fixture.store.phaseRuns) != 2 {
		t.Fatalf("SCN-TFN-011 expected no phase run mutation, got %#v", fixture.store.phaseRuns)
	}
}

func TestSCN_BTP_009_PauseThenCancelPersistsTerminalState(t *testing.T) {
	fixture := newBodyTranslationAPIFixture(t, bodyTranslationAPIFixtureOptions{
		JobState:     "running",
		BodyRunState: "running",
	})
	controller := fixture.controller()

	paused, pauseErr := controller.PauseBodyTranslationPhase(
		controllerwails.PauseBodyTranslationPhaseRequestDTO{
			JobID:      fixture.job.ID,
			PhaseRunID: fixture.bodyRun.ID,
		},
	)
	if pauseErr != nil {
		t.Fatalf("SCN-BTP-009 public pause returned error: %v", pauseErr)
	}
	if paused.PhaseState != "paused" {
		t.Fatalf("SCN-BTP-009 expected paused state, got %#v", paused)
	}

	canceled, cancelErr := controller.CancelBodyTranslationPhase(
		controllerwails.CancelBodyTranslationPhaseRequestDTO{
			JobID:      fixture.job.ID,
			PhaseRunID: fixture.bodyRun.ID,
		},
	)
	if cancelErr != nil {
		t.Fatalf("SCN-BTP-009 public cancel returned error: %v", cancelErr)
	}
	if canceled.PhaseState != "canceled" || canceled.OutputReadiness.Ready {
		t.Fatalf("SCN-BTP-009 expected canceled state and blocked readiness, got %#v", canceled)
	}
}

func TestSCN_BTP_009_RunningCancelIsRejectedBeforeTerminalResultRewrite(t *testing.T) {
	fixture := newBodyTranslationAPIFixture(t, bodyTranslationAPIFixtureOptions{
		JobState:     "running",
		BodyRunState: "running",
		Outputs: []repository.JobTranslationField{
			bodyTranslationAPIOutputField(801, 701, "translated"),
		},
	})
	controller := fixture.controller()

	result, err := controller.CancelBodyTranslationPhase(
		controllerwails.CancelBodyTranslationPhaseRequestDTO{
			JobID:      fixture.job.ID,
			PhaseRunID: fixture.bodyRun.ID,
		},
	)
	if err != nil {
		t.Fatalf("SCN-BTP-009 expected Running cancel rejection payload without Wails error, got %v", err)
	}
	assertBodyTranslationAPIRejectedCommand(t, "SCN-TJSM-003", result, "persona_phase_incomplete", "body translation phase is not cancelable")
	summary, summaryErr := controller.GetBodyTranslationPhaseSummary(
		controllerwails.GetBodyTranslationPhaseSummaryRequestDTO{JobID: fixture.job.ID},
	)
	if summaryErr != nil {
		t.Fatalf("SCN-BTP-009 public summary after rejected cancel returned error: %v", summaryErr)
	}
	if summary.PhaseState != "running" || summary.OutputReadiness.Ready {
		t.Fatalf("SCN-BTP-009 expected running state and blocked readiness after rejected cancel, got %#v", summary)
	}
	if len(fixture.store.phaseRuns) != 2 || len(fixture.store.outputs) != 1 || len(fixture.store.phaseLinks) != 1 {
		t.Fatalf("SCN-TJSM-003 expected rejected Running cancel not to mutate rows, phaseRuns=%#v outputs=%#v links=%#v", fixture.store.phaseRuns, fixture.store.outputs, fixture.store.phaseLinks)
	}
}

func TestSCN_TJSM_008_TerminalJobRejectsStartWithoutCreatingPhaseRun(t *testing.T) {
	fixture := newBodyTranslationAPIFixture(t, bodyTranslationAPIFixtureOptions{
		JobState:     "completed",
		BodyRunState: "completed",
		Outputs: []repository.JobTranslationField{
			bodyTranslationAPIOutputField(801, 701, "ready"),
			bodyTranslationAPIOutputField(802, 702, "ready"),
		},
	})
	controller := fixture.controller()

	result, err := controller.StartBodyTranslationPhase(
		controllerwails.StartBodyTranslationPhaseRequestDTO{JobID: fixture.job.ID},
	)
	if err != nil {
		t.Fatalf("SCN-TJSM-008 expected terminal start rejection payload without Wails error, got %v", err)
	}
	assertBodyTranslationAPIRejectedCommand(t, "SCN-TJSM-008", result, "terminal_job", "terminal_job")
	if result.PhaseState != "completed" {
		t.Fatalf("SCN-TJSM-008 expected terminal start rejection to keep completed phase, got %#v", result)
	}
	if len(fixture.store.phaseRuns) != 2 || len(fixture.store.outputs) != 2 || len(fixture.store.phaseLinks) != 2 {
		t.Fatalf("SCN-TJSM-008 expected terminal start not to mutate rows, phaseRuns=%#v outputs=%#v links=%#v", fixture.store.phaseRuns, fixture.store.outputs, fixture.store.phaseLinks)
	}
}

func TestSCN_TJSM_008_TerminalJobRejectsCancelWithoutStateMutation(t *testing.T) {
	fixture := newBodyTranslationAPIFixture(t, bodyTranslationAPIFixtureOptions{
		JobState:     "completed",
		BodyRunState: "paused",
		Outputs: []repository.JobTranslationField{
			bodyTranslationAPIOutputField(801, 701, "translated"),
		},
	})
	controller := fixture.controller()

	result, err := controller.CancelBodyTranslationPhase(
		controllerwails.CancelBodyTranslationPhaseRequestDTO{
			JobID:      fixture.job.ID,
			PhaseRunID: fixture.bodyRun.ID,
		},
	)
	if err != nil {
		t.Fatalf("SCN-TJSM-008 expected terminal cancel rejection payload without Wails error, got %v", err)
	}
	assertBodyTranslationAPIRejectedCommand(t, "SCN-TJSM-008", result, "terminal_job", "terminal_job")
	if result.PhaseState != "paused" {
		t.Fatalf("SCN-TJSM-008 expected terminal cancel rejection to keep paused phase, got %#v", result)
	}
	if fixture.store.job.State != "completed" || fixture.store.phaseRuns[1].State != "paused" {
		t.Fatalf("SCN-TJSM-008 expected terminal cancel not to mutate state, job=%#v phaseRun=%#v", fixture.store.job, fixture.store.phaseRuns[1])
	}
	if len(fixture.store.phaseRuns) != 2 || len(fixture.store.outputs) != 1 || len(fixture.store.phaseLinks) != 1 {
		t.Fatalf("SCN-TJSM-008 expected terminal cancel not to mutate rows, phaseRuns=%#v outputs=%#v links=%#v", fixture.store.phaseRuns, fixture.store.outputs, fixture.store.phaseLinks)
	}
}

func TestSCN_BTP_010_CompletedConsistentResultEnablesOutputReadiness(t *testing.T) {
	fixture := newBodyTranslationAPIFixture(t, bodyTranslationAPIFixtureOptions{
		JobState:     "completed",
		BodyRunState: "completed",
		Outputs: []repository.JobTranslationField{
			bodyTranslationAPIOutputField(801, 701, "ready"),
			bodyTranslationAPIOutputField(802, 702, "translated"),
		},
	})
	controller := fixture.controller()

	result, err := controller.GetBodyTranslationOutputReadiness(
		controllerwails.GetBodyTranslationOutputReadinessRequestDTO{JobID: fixture.job.ID},
	)
	if err != nil {
		t.Fatalf("SCN-BTP-010 public readiness returned error: %v", err)
	}
	if !result.Ready || !result.StatusConsistent || result.CompletedFieldCount != 2 {
		t.Fatalf("SCN-BTP-010 expected downstream readiness for completed consistent result, got %#v", result)
	}
}

func TestSCN_BTP_010_StatusInconsistencyBlocksOutputReadiness(t *testing.T) {
	fixture := newBodyTranslationAPIFixture(t, bodyTranslationAPIFixtureOptions{
		JobState:     "completed",
		BodyRunState: "completed",
		Outputs: []repository.JobTranslationField{
			bodyTranslationAPIOutputField(801, 701, "ready"),
			bodyTranslationAPIOutputField(802, 702, "failed"),
		},
	})
	controller := fixture.controller()

	result, err := controller.GetBodyTranslationOutputReadiness(
		controllerwails.GetBodyTranslationOutputReadinessRequestDTO{JobID: fixture.job.ID},
	)
	if err != nil {
		t.Fatalf("SCN-BTP-010 public readiness returned error: %v", err)
	}
	if result.Ready || result.StatusConsistent {
		t.Fatalf("SCN-BTP-010 expected status inconsistency to block readiness, got %#v", result)
	}
	if result.ErrorKind != "output_readiness_blocked" {
		t.Fatalf("SCN-BTP-010 expected output_readiness_blocked, got %#v", result)
	}
}

func assertBodyTranslationAPIRejectedCommand(
	t *testing.T,
	scenarioID string,
	result controllerwails.BodyTranslationPhaseCommandResponseDTO,
	expectedErrorKind string,
	expectedReason string,
) {
	t.Helper()

	if result.ErrorSummary == nil {
		t.Fatalf("%s expected rejection payload, got %#v", scenarioID, result)
	}
	if result.ErrorSummary.ErrorKind != expectedErrorKind {
		t.Fatalf("%s expected error kind %q, got %#v", scenarioID, expectedErrorKind, result.ErrorSummary)
	}
	if result.ErrorSummary.Reason != expectedReason {
		t.Fatalf("%s expected reason %q, got %#v", scenarioID, expectedReason, result.ErrorSummary)
	}
	if result.ErrorSummary.Retryable || result.Retryable {
		t.Fatalf("%s expected non-retryable rejection, got %#v", scenarioID, result)
	}
}

func bodyTranslationAPIOutputCountForField(outputs []repository.JobTranslationField, translationFieldID int64) int {
	count := 0
	for _, output := range outputs {
		if output.TranslationFieldID == translationFieldID {
			count++
		}
	}
	return count
}

type bodyTranslationAPIFixtureOptions struct {
	JobState     string
	BodyRunState string
	LatestError  string
	AIProvider   string
	Outputs      []repository.JobTranslationField
	Provider     service.BodyTranslationProvider
	OmitBodyRun  bool
}

type bodyTranslationAPIFixture struct {
	t       *testing.T
	store   *bodyTranslationAPIStore
	job     repository.TranslationJob
	bodyRun repository.JobPhaseRun
}

func newBodyTranslationAPIFixture(
	t *testing.T,
	options bodyTranslationAPIFixtureOptions,
) *bodyTranslationAPIFixture {
	t.Helper()
	now := time.Date(2026, 5, 2, 9, 0, 0, 0, time.UTC)
	jobState := firstNonEmptyBodyTranslationAPIValue(options.JobState, "running")
	bodyRunState := firstNonEmptyBodyTranslationAPIValue(options.BodyRunState, "running")
	store := &bodyTranslationAPIStore{
		xedit: repository.XEditExtractedData{
			ID:                301,
			SourceFilePath:    "fixture.esp",
			SourceContentHash: "sha256:xedit",
			SourceTool:        "xedit",
			TargetPluginName:  "Fixture.esp",
			TargetPluginType:  "esp",
			RecordCount:       1,
			ImportedAt:        now,
		},
		records: []repository.TranslationRecord{
			{ID: 601, XEditExtractedDataID: 301, FormID: "000001", EditorID: "NPCFixture", RecordType: "NPC_"},
		},
		fieldsByRecordID: map[int64][]repository.TranslationField{
			601: {
				{ID: 701, TranslationRecordID: 601, SubrecordType: "FULL", SourceText: "Hello <Alias=Player>.", FieldOrder: 1},
				{ID: 702, TranslationRecordID: 601, SubrecordType: "DESC", SourceText: "A second line.", FieldOrder: 2},
			},
		},
		personasByJobID: map[int64][]repository.Persona{
			101: {
				{ID: 501, TranslationJobID: int64PointerForBodyTranslationAPI(101), PersonaLifecycle: "job", PersonaScope: "job", PersonaDescription: "calm speaker", SpeechStyle: "plain", PersonalitySummary: "steady", CreatedAt: now, UpdatedAt: now},
			},
		},
		dictionary: []repository.DictionaryEntry{
			{ID: 401, TranslationJobID: int64PointerForBodyTranslationAPI(101), DictionaryLifecycle: "job", DictionaryScope: "job", SourceTerm: "Hello", TranslatedTerm: "こんにちは", TermKind: "term", Reusable: true, CreatedAt: now, UpdatedAt: now},
		},
		provider: options.Provider,
	}
	store.job = repository.TranslationJob{
		ID:                   101,
		XEditExtractedDataID: 301,
		JobName:              "body translation api scenario",
		State:                jobState,
		ProgressPercent:      20,
		CreatedAt:            now,
		StartedAt:            &now,
	}
	if jobState == "completed" {
		store.job.ProgressPercent = 100
		store.job.FinishedAt = &now
	}
	store.phaseRuns = []repository.JobPhaseRun{
		{
			ID:               201,
			TranslationJobID: 101,
			PhaseType:        "persona_generation",
			State:            "completed",
			ExecutionOrder:   1,
			AIProvider:       "fake",
			ModelName:        "body-fixture-model",
			ExecutionMode:    "single_request",
			CredentialRef:    "credential:body-fixture",
			InstructionKind:  "persona_generation",
			StartedAt:        &now,
			FinishedAt:       &now,
		},
	}
	if !options.OmitBodyRun {
		store.phaseRuns = append(store.phaseRuns, repository.JobPhaseRun{
			ID:               202,
			TranslationJobID: 101,
			PhaseType:        "body_translation",
			State:            bodyRunState,
			ExecutionOrder:   2,
			ProgressPercent:  50,
			AIProvider:       firstNonEmptyBodyTranslationAPIValue(options.AIProvider, "fake"),
			ModelName:        "body-fixture-model",
			ExecutionMode:    "single_request",
			CredentialRef:    "credential:body-fixture",
			InstructionKind:  "body_translation",
			LatestError:      options.LatestError,
			StartedAt:        &now,
		})
	}
	if !options.OmitBodyRun && (bodyRunState == "completed" || bodyRunState == "canceled") {
		store.phaseRuns[1].ProgressPercent = 100
		store.phaseRuns[1].FinishedAt = &now
	}
	store.outputs = append([]repository.JobTranslationField(nil), options.Outputs...)
	for _, output := range store.outputs {
		store.phaseLinks = append(store.phaseLinks, repository.PhaseRunTranslationField{
			ID:                    int64(len(store.phaseLinks) + 901),
			PhaseRunID:            202,
			JobTranslationFieldID: output.ID,
			Role:                  "provider_result",
		})
	}
	return &bodyTranslationAPIFixture{
		t:       t,
		store:   store,
		job:     store.job,
		bodyRun: bodyTranslationAPIBodyRun(store.phaseRuns),
	}
}

func bodyTranslationAPIBodyRun(phaseRuns []repository.JobPhaseRun) repository.JobPhaseRun {
	for _, run := range phaseRuns {
		if run.PhaseType == "body_translation" {
			return run
		}
	}
	return repository.JobPhaseRun{}
}

func (fixture *bodyTranslationAPIFixture) controller() *controllerwails.BodyTranslationPhaseController {
	fixture.t.Helper()
	phaseService := service.NewBodyTranslationPhaseService(
		fixture.store,
		fixture.store,
		fixture.store,
		fixture.store,
		fixture.store,
	)
	if fixture.optionsProvider() != nil {
		phaseService.WithBodyTranslationProvider(fixture.optionsProvider())
	}
	phaseService.WithBodyTranslationProviderSettings(bodyTranslationAPIProviderSettings{})
	phaseUsecase := usecase.NewBodyTranslationPhaseUsecase(phaseService)
	return controllerwails.NewBodyTranslationPhaseController(phaseUsecase)
}

func (fixture *bodyTranslationAPIFixture) optionsProvider() service.BodyTranslationProvider {
	fixture.t.Helper()
	return fixture.store.provider
}

type bodyTranslationAPIProviderSettings struct{}

func (bodyTranslationAPIProviderSettings) ListProviderSettings(context.Context) (service.ProviderSettingsRoute, []service.ProviderSettingsSummary, error) {
	return service.ProviderSettingsRoute{}, nil, nil
}

func (bodyTranslationAPIProviderSettings) SaveProviderSettings(context.Context, service.ProviderSettingsSaveInput) (service.ProviderSettingsSummary, error) {
	return service.ProviderSettingsSummary{}, nil
}

func (bodyTranslationAPIProviderSettings) ListProviderModels(context.Context, service.ProviderSettingsModelListInput) (service.ProviderSettingsModelListResult, error) {
	return service.ProviderSettingsModelListResult{}, nil
}

func (bodyTranslationAPIProviderSettings) ResolveProviderExecutionSettings(_ context.Context, input service.ProviderSettingsResolveInput) (service.ProviderSettingsResolveResult, error) {
	endpoint := "https://body-provider.example.test"
	providerReference := "body-fixture-ref"
	return service.ProviderSettingsResolveResult{
		ConsumerID:            input.ConsumerID,
		ProviderID:            input.Selection.ProviderID,
		Model:                 input.Selection.Model,
		ExecutionMethod:       input.Selection.ExecutionMethod,
		UseBatchAPI:           input.Selection.UseBatchAPI,
		Endpoint:              &endpoint,
		CredentialReferenceID: &providerReference,
		CredentialState:       "configured",
	}, nil
}

func bodyTranslationAPIOutputField(
	id int64,
	translationFieldID int64,
	outputStatus string,
) repository.JobTranslationField {
	return repository.JobTranslationField{
		ID:                 id,
		TranslationJobID:   101,
		TranslationFieldID: translationFieldID,
		TranslatedText:     "translated fixture",
		OutputStatus:       outputStatus,
		RetryCount:         0,
		UpdatedAt:          time.Date(2026, 5, 2, 9, 1, 0, 0, time.UTC),
	}
}

type bodyTranslationAPIStore struct {
	job              repository.TranslationJob
	phaseRuns        []repository.JobPhaseRun
	phaseLinks       []repository.PhaseRunTranslationField
	xedit            repository.XEditExtractedData
	records          []repository.TranslationRecord
	fieldsByRecordID map[int64][]repository.TranslationField
	personasByJobID  map[int64][]repository.Persona
	dictionary       []repository.DictionaryEntry
	outputs          []repository.JobTranslationField
	provider         service.BodyTranslationProvider
	nextPhaseRunID   int64
	nextLinkID       int64
	nextOutputID     int64
}

type bodyTranslationAPISuccessProvider struct{}

func (bodyTranslationAPISuccessProvider) BodyTranslationProviderRequestsAreTestSafe() bool {
	return true
}

func (bodyTranslationAPISuccessProvider) TranslateBodyField(
	_ context.Context,
	request service.BodyTranslationProviderRequest,
) service.BodyTranslationProviderResult {
	return service.BodyTranslationProviderResult{
		FieldCorrelationKey: request.FieldCorrelationKey,
		TranslatedCandidate: &service.BodyTranslationTranslatedCandidate{
			FieldCorrelationKey: request.FieldCorrelationKey,
			RecordType:          request.RecordType,
			FieldType:           request.FieldType,
			TranslatedText:      "translated " + request.SourceText,
		},
		ProtectionValidationTarget: &service.BodyTranslationProtectionValidationTarget{
			FieldCorrelationKey: request.FieldCorrelationKey,
			TranslatedText:      "translated " + request.SourceText,
		},
	}
}

func (store *bodyTranslationAPIStore) WithTransaction(ctx context.Context, fn func(ctx context.Context) error) error {
	return fn(ctx)
}

func (store *bodyTranslationAPIStore) GetTranslationJobByID(
	_ context.Context,
	id int64,
) (repository.TranslationJob, error) {
	if store.job.ID == id {
		return store.job, nil
	}
	return repository.TranslationJob{}, repository.ErrNotFound
}

func (store *bodyTranslationAPIStore) UpdateTranslationJob(
	_ context.Context,
	id int64,
	draft repository.TranslationJobUpdateDraft,
) (repository.TranslationJob, error) {
	if store.job.ID != id {
		return repository.TranslationJob{}, repository.ErrNotFound
	}
	store.job.JobName = draft.JobName
	store.job.State = draft.State
	store.job.ProgressPercent = draft.ProgressPercent
	store.job.StartedAt = draft.StartedAt
	store.job.FinishedAt = draft.FinishedAt
	return store.job, nil
}

func (store *bodyTranslationAPIStore) CreateJobPhaseRun(
	_ context.Context,
	draft repository.JobPhaseRunDraft,
) (repository.JobPhaseRun, error) {
	if store.nextPhaseRunID == 0 {
		store.nextPhaseRunID = 300
	}
	store.nextPhaseRunID++
	run := repository.JobPhaseRun{
		ID:                     store.nextPhaseRunID,
		TranslationJobID:       draft.TranslationJobID,
		PhaseType:              draft.PhaseType,
		State:                  draft.State,
		ExecutionOrder:         draft.ExecutionOrder,
		AIProvider:             draft.AIProvider,
		ModelName:              draft.ModelName,
		ExecutionMode:          draft.ExecutionMode,
		CredentialRef:          draft.CredentialRef,
		InstructionKind:        draft.InstructionKind,
		SnapshotFieldCount:     draft.SnapshotFieldCount,
		ProviderTargetCount:    draft.ProviderTargetCount,
		ExactExclusionCount:    draft.ExactExclusionCount,
		PartialConstraintCount: draft.PartialConstraintCount,
		InputSnapshotDigest:    draft.InputSnapshotDigest,
		DictionaryDigest:       draft.DictionaryDigest,
		PersonaDigest:          draft.PersonaDigest,
		MetadataDigest:         draft.MetadataDigest,
		PromptDigest:           draft.PromptDigest,
	}
	store.phaseRuns = append(store.phaseRuns, run)
	return run, nil
}

func (store *bodyTranslationAPIStore) FindJobPhaseRun(
	_ context.Context,
	translationJobID int64,
	phaseType string,
) (repository.JobPhaseRun, error) {
	for _, run := range store.phaseRuns {
		if run.TranslationJobID == translationJobID && strings.TrimSpace(run.PhaseType) == strings.TrimSpace(phaseType) {
			return run, nil
		}
	}
	return repository.JobPhaseRun{}, repository.ErrNotFound
}

func (store *bodyTranslationAPIStore) UpdateJobPhaseRun(
	_ context.Context,
	id int64,
	draft repository.JobPhaseRunUpdateDraft,
) (repository.JobPhaseRun, error) {
	for index := range store.phaseRuns {
		if store.phaseRuns[index].ID != id {
			continue
		}
		store.phaseRuns[index].State = draft.State
		store.phaseRuns[index].ProgressPercent = draft.ProgressPercent
		store.phaseRuns[index].LatestExternalRunID = draft.LatestExternalRunID
		store.phaseRuns[index].LatestError = draft.LatestError
		store.phaseRuns[index].StartedAt = draft.StartedAt
		store.phaseRuns[index].FinishedAt = draft.FinishedAt
		return store.phaseRuns[index], nil
	}
	return repository.JobPhaseRun{}, repository.ErrNotFound
}

func (store *bodyTranslationAPIStore) UpdateJobPhaseRunWhenState(
	ctx context.Context,
	id int64,
	expectedState string,
	draft repository.JobPhaseRunUpdateDraft,
) (repository.JobPhaseRun, error) {
	for _, run := range store.phaseRuns {
		if run.ID == id && strings.TrimSpace(run.State) != strings.TrimSpace(expectedState) {
			return repository.JobPhaseRun{}, repository.ErrConflict
		}
	}
	return store.UpdateJobPhaseRun(ctx, id, draft)
}

func (store *bodyTranslationAPIStore) ListJobPhaseRunsByJobID(
	_ context.Context,
	jobID int64,
) ([]repository.JobPhaseRun, error) {
	var runs []repository.JobPhaseRun
	for _, run := range store.phaseRuns {
		if run.TranslationJobID == jobID {
			runs = append(runs, run)
		}
	}
	return runs, nil
}

func (store *bodyTranslationAPIStore) CreatePhaseRunTranslationField(
	_ context.Context,
	draft repository.PhaseRunTranslationFieldDraft,
) (repository.PhaseRunTranslationField, error) {
	if store.nextLinkID == 0 {
		store.nextLinkID = 950
	}
	store.nextLinkID++
	link := repository.PhaseRunTranslationField{
		ID:                    store.nextLinkID,
		PhaseRunID:            draft.PhaseRunID,
		JobTranslationFieldID: draft.JobTranslationFieldID,
		Role:                  draft.Role,
	}
	store.phaseLinks = append(store.phaseLinks, link)
	return link, nil
}

func (store *bodyTranslationAPIStore) ListPhaseRunTranslationFieldsByPhaseRunID(
	_ context.Context,
	phaseRunID int64,
) ([]repository.PhaseRunTranslationField, error) {
	var links []repository.PhaseRunTranslationField
	for _, link := range store.phaseLinks {
		if link.PhaseRunID == phaseRunID {
			links = append(links, link)
		}
	}
	return links, nil
}

func (store *bodyTranslationAPIStore) ListDictionaryEntries(
	_ context.Context,
	translationJobID *int64,
	_ string,
	_ string,
	_ string,
) ([]repository.DictionaryEntry, error) {
	if translationJobID == nil || *translationJobID != store.job.ID {
		return nil, nil
	}
	return append([]repository.DictionaryEntry(nil), store.dictionary...), nil
}

func (store *bodyTranslationAPIStore) ListPersonasByTranslationJobID(
	_ context.Context,
	translationJobID int64,
) ([]repository.Persona, error) {
	return append([]repository.Persona(nil), store.personasByJobID[translationJobID]...), nil
}

func (store *bodyTranslationAPIStore) GetXEditExtractedDataByID(
	_ context.Context,
	id int64,
) (repository.XEditExtractedData, error) {
	if store.xedit.ID == id {
		return store.xedit, nil
	}
	return repository.XEditExtractedData{}, repository.ErrNotFound
}

func (store *bodyTranslationAPIStore) ListTranslationRecordsByXEditID(
	_ context.Context,
	xEditID int64,
) ([]repository.TranslationRecord, error) {
	var records []repository.TranslationRecord
	for _, record := range store.records {
		if record.XEditExtractedDataID == xEditID {
			records = append(records, record)
		}
	}
	return records, nil
}

func (store *bodyTranslationAPIStore) ListTranslationFieldsByTranslationRecordID(
	_ context.Context,
	translationRecordID int64,
) ([]repository.TranslationField, error) {
	return append([]repository.TranslationField(nil), store.fieldsByRecordID[translationRecordID]...), nil
}

func (store *bodyTranslationAPIStore) CreateJobTranslationField(
	_ context.Context,
	draft repository.JobTranslationFieldDraft,
) (repository.JobTranslationField, error) {
	if store.nextOutputID == 0 {
		store.nextOutputID = 850
	}
	store.nextOutputID++
	output := repository.JobTranslationField{
		ID:                 store.nextOutputID,
		TranslationJobID:   draft.TranslationJobID,
		TranslationFieldID: draft.TranslationFieldID,
		AppliedPersonaID:   draft.AppliedPersonaID,
		TranslatedText:     draft.TranslatedText,
		OutputStatus:       draft.OutputStatus,
		RetryCount:         draft.RetryCount,
		UpdatedAt:          time.Date(2026, 5, 2, 9, 2, 0, 0, time.UTC),
	}
	store.outputs = append(store.outputs, output)
	return output, nil
}

func (store *bodyTranslationAPIStore) UpdateJobTranslationField(
	_ context.Context,
	id int64,
	draft repository.JobTranslationFieldUpdateDraft,
) (repository.JobTranslationField, error) {
	for index := range store.outputs {
		if store.outputs[index].ID != id {
			continue
		}
		store.outputs[index].AppliedPersonaID = draft.AppliedPersonaID
		store.outputs[index].TranslatedText = draft.TranslatedText
		store.outputs[index].OutputStatus = draft.OutputStatus
		store.outputs[index].RetryCount = draft.RetryCount
		store.outputs[index].UpdatedAt = time.Date(2026, 5, 2, 9, 3, 0, 0, time.UTC)
		return store.outputs[index], nil
	}
	return repository.JobTranslationField{}, repository.ErrNotFound
}

func (store *bodyTranslationAPIStore) ListJobTranslationFieldsByJobID(
	_ context.Context,
	jobID int64,
) ([]repository.JobTranslationField, error) {
	var outputs []repository.JobTranslationField
	for _, output := range store.outputs {
		if output.TranslationJobID == jobID {
			outputs = append(outputs, output)
		}
	}
	return outputs, nil
}

func firstNonEmptyBodyTranslationAPIValue(values ...string) string {
	for _, value := range values {
		if strings.TrimSpace(value) != "" {
			return strings.TrimSpace(value)
		}
	}
	return ""
}

func int64PointerForBodyTranslationAPI(value int64) *int64 {
	return &value
}

var _ repository.Transactor = (*bodyTranslationAPIStore)(nil)
