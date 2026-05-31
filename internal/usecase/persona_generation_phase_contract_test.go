package usecase

import (
	"context"
	"testing"
)

type personaGenerationPhaseContractAPI interface {
	GetPersonaGenerationPhaseSummary(
		context.Context,
		GetPersonaGenerationPhaseSummaryRequest,
	) (PersonaGenerationPhaseSummaryResult, error)
	StartPersonaGenerationPhase(
		context.Context,
		StartPersonaGenerationPhaseRequest,
	) (PersonaGenerationPhaseCommandResult, error)
	PausePersonaGenerationPhase(
		context.Context,
		PausePersonaGenerationPhaseRequest,
	) (PersonaGenerationPhaseCommandResult, error)
	ResumePersonaGenerationPhase(
		context.Context,
		ResumePersonaGenerationPhaseRequest,
	) (PersonaGenerationPhaseCommandResult, error)
	RetryPersonaGenerationPhase(
		context.Context,
		RetryPersonaGenerationPhaseRequest,
	) (PersonaGenerationPhaseCommandResult, error)
	CancelPersonaGenerationPhase(
		context.Context,
		CancelPersonaGenerationPhaseRequest,
	) (PersonaGenerationPhaseCommandResult, error)
	GetPersonaGenerationBodyReadiness(
		context.Context,
		GetPersonaGenerationBodyReadinessRequest,
	) (PersonaGenerationBodyReadinessResult, error)
}

func newPersonaGenerationPhaseContractAPI(t *testing.T) personaGenerationPhaseContractAPI {
	t.Helper()
	return NewPersonaGenerationPhaseContractStub()
}

func requirePersonaGenerationPhaseRunID(t *testing.T, phaseRunID *int64) int64 {
	t.Helper()
	if phaseRunID == nil {
		t.Fatal("expected phase run ID")
	}
	return *phaseRunID
}

func TestPersonaGenerationContract_PublicSeamsAndRedactionDTOShape(t *testing.T) {
	api := newPersonaGenerationPhaseContractAPI(t)
	_ = api

	secretLikeFieldsMustNotExist := PersonaGenerationPhaseSummaryResult{
		JobID:        101,
		CurrentPhase: "persona_generation",
		PhaseState:   "running",
		Progress: PersonaGenerationPhaseProgressSummary{
			Percent:        40,
			ProcessedCount: 2,
			TotalCount:     5,
			TargetCount:    5,
		},
		TargetSummary: PersonaGenerationTargetSummary{
			TargetCount:            5,
			CommonPersonaHitCount:  1,
			CommonPersonaMissCount: 4,
			SkippedCount:           1,
			SkippedReasons:         []string{"orphan_npc_reference"},
			TargetSnapshotDigest:   "sha256:target-snapshot",
		},
		Execution: PersonaGenerationExecutionSummary{
			CredentialRef: "credential:persona:test",
			Provider:      "fake",
			Model:         "persona-model",
			ExecutionMode: "single",
			PromptDigest:  "sha256:prompt",
			InputCount:    5,
			OutputCount:   4,
			EvidenceRefs:  []string{"evidence:npc:001"},
		},
		ResultSummary: &PersonaGenerationPhaseResultSummary{
			GeneratedCount:          4,
			FailedCount:             1,
			PersonaCount:            4,
			MissingCount:            1,
			SnapshotID:              "persona-snapshot-101",
			SnapshotDigest:          "sha256:persona-snapshot",
			SnapshotReferenceStatus: "available",
			BodyReadiness:           false,
		},
		ErrorSummary: &PersonaGenerationPhaseErrorSummary{
			ErrorKind:  PersonaGenerationPhaseErrorKindProviderFailure,
			Reason:     "provider failure",
			Retryable:  true,
			IsRedacted: true,
		},
		ActionEnablement: PersonaGenerationPhaseActionEnablement{
			CanStart:  false,
			CanPause:  true,
			CanResume: false,
			CanRetry:  true,
			CanCancel: true,
		},
	}

	if secretLikeFieldsMustNotExist.Execution.CredentialRef == "" {
		t.Fatal("expected credential ref to be exposed as a reference value")
	}
	if secretLikeFieldsMustNotExist.Execution.PromptDigest == "" {
		t.Fatal("expected prompt digest to be exposed instead of raw prompt")
	}
	if !secretLikeFieldsMustNotExist.ErrorSummary.IsRedacted {
		t.Fatal("expected public error summary to carry redaction status")
	}
}

func TestPersonaGenerationContract_TargetSnapshotCommonPersonaHitMissAndSkippedReasons(t *testing.T) {
	api := newPersonaGenerationPhaseContractAPI(t)

	result, err := api.StartPersonaGenerationPhase(context.Background(), StartPersonaGenerationPhaseRequest{
		JobID: 2002,
	})
	if err != nil {
		t.Fatalf("StartPersonaGenerationPhase failed: %v", err)
	}

	if result.TargetSummary.TargetSnapshotDigest == "" {
		t.Fatal("expected target snapshot digest")
	}
	if result.TargetSummary.TargetCount != result.Progress.TotalCount {
		t.Fatalf("expected target count to match progress total, got target=%d total=%d",
			result.TargetSummary.TargetCount,
			result.Progress.TotalCount,
		)
	}
	if result.TargetSummary.CommonPersonaHitCount == 0 {
		t.Fatal("expected common persona hit count")
	}
	if result.TargetSummary.CommonPersonaMissCount == 0 {
		t.Fatal("expected common persona miss count")
	}
	if len(result.TargetSummary.SkippedReasons) == 0 {
		t.Fatal("expected skipped target reasons")
	}
}

func TestPersonaGenerationContract_PersonaPersistenceRejectsPartialSaveFailure(t *testing.T) {
	api := newPersonaGenerationPhaseContractAPI(t)

	result, err := api.RetryPersonaGenerationPhase(context.Background(), RetryPersonaGenerationPhaseRequest{
		JobID:      2004,
		PhaseRunID: 4004,
	})
	if err == nil {
		t.Fatal("expected save failure to reject persona phase completion")
	}
	if result.PhaseState == "completed" {
		t.Fatal("expected partial save failure not to complete persona phase")
	}
	if result.ErrorSummary == nil {
		t.Fatal("expected redacted save failure summary")
	}
	if result.ErrorSummary.ErrorKind != PersonaGenerationPhaseErrorKindSaveFailed {
		t.Fatalf("expected save_failed, got %q", result.ErrorSummary.ErrorKind)
	}
}

func TestPersonaGenerationContract_BodyReadinessRequiresCompletedPersonaSnapshot(t *testing.T) {
	api := newPersonaGenerationPhaseContractAPI(t)

	result, err := api.GetPersonaGenerationBodyReadiness(
		context.Background(),
		GetPersonaGenerationBodyReadinessRequest{JobID: 2006},
	)
	if err != nil {
		t.Fatalf("GetPersonaGenerationBodyReadiness failed: %v", err)
	}

	if result.InputSummary.PersonaCount == 0 {
		t.Fatal("expected body input summary persona count")
	}
	if result.InputSummary.SnapshotDigest == "" {
		t.Fatal("expected body input summary snapshot digest")
	}
}

func TestPersonaGenerationContract_FailuresAreNotPublishedAsSuccessfulPersonaSnapshot(t *testing.T) {
	api := newPersonaGenerationPhaseContractAPI(t)

	for _, tc := range []struct {
		name      string
		jobID     int64
		errorKind PersonaGenerationPhaseErrorKind
	}{
		{name: "provider failure", jobID: 2071, errorKind: PersonaGenerationPhaseErrorKindProviderFailure},
		{name: "invalid response", jobID: 2072, errorKind: PersonaGenerationPhaseErrorKindInvalidProviderResponse},
		{name: "input missing", jobID: 2073, errorKind: PersonaGenerationPhaseErrorKindInputMissing},
		{name: "save failure", jobID: 2074, errorKind: PersonaGenerationPhaseErrorKindSaveFailed},
	} {
		t.Run(tc.name, func(t *testing.T) {
			result, err := api.RetryPersonaGenerationPhase(context.Background(), RetryPersonaGenerationPhaseRequest{
				JobID:      tc.jobID,
				PhaseRunID: 4007,
			})
			if err == nil {
				t.Fatal("expected failed fixture not to be treated as success")
			}
			_ = result
			if result.ResultSummary != nil && result.ResultSummary.SnapshotReferenceStatus == "available" {
				t.Fatal("expected failed target not to publish successful snapshot")
			}
			if result.ErrorSummary == nil || result.ErrorSummary.ErrorKind != tc.errorKind {
				t.Fatalf("expected error kind %q, got %#v", tc.errorKind, result.ErrorSummary)
			}
		})
	}
}

func TestPersonaGenerationContract_RetryResumeAndResendReuseSamePhaseRunWithoutDuplicates(t *testing.T) {
	api := newPersonaGenerationPhaseContractAPI(t)

	started, err := api.StartPersonaGenerationPhase(context.Background(), StartPersonaGenerationPhaseRequest{
		JobID: 2008,
	})
	if err != nil {
		t.Fatalf("StartPersonaGenerationPhase failed: %v", err)
	}
	startedPhaseRunID := requirePersonaGenerationPhaseRunID(t, started.PhaseRunID)
	resumed, err := api.ResumePersonaGenerationPhase(context.Background(), ResumePersonaGenerationPhaseRequest{
		JobID:      2008,
		PhaseRunID: startedPhaseRunID,
	})
	if err != nil {
		t.Fatalf("ResumePersonaGenerationPhase failed: %v", err)
	}
	retried, err := api.RetryPersonaGenerationPhase(context.Background(), RetryPersonaGenerationPhaseRequest{
		JobID:      2008,
		PhaseRunID: startedPhaseRunID,
	})
	if err != nil {
		t.Fatalf("RetryPersonaGenerationPhase failed: %v", err)
	}

	resumedPhaseRunID := requirePersonaGenerationPhaseRunID(t, resumed.PhaseRunID)
	retriedPhaseRunID := requirePersonaGenerationPhaseRunID(t, retried.PhaseRunID)
	if startedPhaseRunID != resumedPhaseRunID || startedPhaseRunID != retriedPhaseRunID {
		t.Fatalf("expected same phase run, got start=%d resume=%d retry=%d",
			startedPhaseRunID,
			resumedPhaseRunID,
			retriedPhaseRunID,
		)
	}
	if started.ResultSummary == nil || retried.ResultSummary == nil {
		t.Fatal("expected persona result summaries")
	}
	if retried.ResultSummary.PersonaCount != started.ResultSummary.PersonaCount {
		t.Fatalf("expected no duplicate persona, got start=%d retry=%d",
			started.ResultSummary.PersonaCount,
			retried.ResultSummary.PersonaCount,
		)
	}
	if retried.TargetSummary.TargetSnapshotDigest != started.TargetSummary.TargetSnapshotDigest {
		t.Fatal("expected stable target snapshot digest across retry")
	}
}

func TestPersonaGenerationContract_TerminalJobRejectsPersonaAndBodyReadinessWrites(t *testing.T) {
	api := newPersonaGenerationPhaseContractAPI(t)

	result, err := api.StartPersonaGenerationPhase(context.Background(), StartPersonaGenerationPhaseRequest{
		JobID: 2010,
	})
	if err == nil {
		t.Fatal("expected terminal job to reject persona phase start")
	}
	if result.PhaseRunID != nil {
		t.Fatalf("expected no new phase run, got %d", *result.PhaseRunID)
	}
	if result.ErrorSummary == nil {
		t.Fatal("expected terminal job error summary")
	}
	if result.ErrorSummary.ErrorKind != PersonaGenerationPhaseErrorKindTerminalJob {
		t.Fatalf("expected terminal_job, got %q", result.ErrorSummary.ErrorKind)
	}

	readiness, readinessErr := api.GetPersonaGenerationBodyReadiness(
		context.Background(),
		GetPersonaGenerationBodyReadinessRequest{JobID: 2010},
	)
	if readinessErr == nil {
		t.Fatal("expected terminal job to reject body readiness update")
	}
	if !readiness.JobIsTerminal {
		t.Fatal("expected terminal job to set JobIsTerminal in body readiness")
	}
}

func TestPersonaGenerationContract_CancelFixtureUsesCanceledSpelling(t *testing.T) {
	api := newPersonaGenerationPhaseContractAPI(t)

	result, err := api.CancelPersonaGenerationPhase(context.Background(), CancelPersonaGenerationPhaseRequest{
		JobID:      2008,
		PhaseRunID: 4008,
	})
	if err != nil {
		t.Fatalf("CancelPersonaGenerationPhase failed: %v", err)
	}
	if result.PhaseState != "canceled" {
		t.Fatalf("expected canceled phase state, got %q", result.PhaseState)
	}
	if result.Progress.CurrentStep != "canceled" {
		t.Fatalf("expected canceled progress step, got %q", result.Progress.CurrentStep)
	}
}
