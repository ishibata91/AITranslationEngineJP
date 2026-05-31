package wails

import (
	"context"
	"encoding/json"
	"errors"
	"strings"
	"testing"
	"time"

	"aitranslationenginejp/internal/usecase"
)

type fakePersonaGenerationPhaseUsecase struct {
	getSummaryFunc       func(context.Context, usecase.GetPersonaGenerationPhaseSummaryRequest) (usecase.PersonaGenerationPhaseSummaryResult, error)
	startFunc            func(context.Context, usecase.StartPersonaGenerationPhaseRequest) (usecase.PersonaGenerationPhaseCommandResult, error)
	pauseFunc            func(context.Context, usecase.PausePersonaGenerationPhaseRequest) (usecase.PersonaGenerationPhaseCommandResult, error)
	resumeFunc           func(context.Context, usecase.ResumePersonaGenerationPhaseRequest) (usecase.PersonaGenerationPhaseCommandResult, error)
	retryFunc            func(context.Context, usecase.RetryPersonaGenerationPhaseRequest) (usecase.PersonaGenerationPhaseCommandResult, error)
	cancelFunc           func(context.Context, usecase.CancelPersonaGenerationPhaseRequest) (usecase.PersonaGenerationPhaseCommandResult, error)
	getBodyReadinessFunc func(context.Context, usecase.GetPersonaGenerationBodyReadinessRequest) (usecase.PersonaGenerationBodyReadinessResult, error)
}

func (fake fakePersonaGenerationPhaseUsecase) GetPersonaGenerationPhaseSummary(ctx context.Context, request usecase.GetPersonaGenerationPhaseSummaryRequest) (usecase.PersonaGenerationPhaseSummaryResult, error) {
	if fake.getSummaryFunc != nil {
		return fake.getSummaryFunc(ctx, request)
	}
	return usecase.PersonaGenerationPhaseSummaryResult{}, nil
}
func (fake fakePersonaGenerationPhaseUsecase) StartPersonaGenerationPhase(ctx context.Context, request usecase.StartPersonaGenerationPhaseRequest) (usecase.PersonaGenerationPhaseCommandResult, error) {
	if fake.startFunc != nil {
		return fake.startFunc(ctx, request)
	}
	return usecase.PersonaGenerationPhaseCommandResult{}, nil
}
func (fake fakePersonaGenerationPhaseUsecase) PausePersonaGenerationPhase(ctx context.Context, request usecase.PausePersonaGenerationPhaseRequest) (usecase.PersonaGenerationPhaseCommandResult, error) {
	if fake.pauseFunc != nil {
		return fake.pauseFunc(ctx, request)
	}
	return usecase.PersonaGenerationPhaseCommandResult{}, nil
}
func (fake fakePersonaGenerationPhaseUsecase) ResumePersonaGenerationPhase(ctx context.Context, request usecase.ResumePersonaGenerationPhaseRequest) (usecase.PersonaGenerationPhaseCommandResult, error) {
	if fake.resumeFunc != nil {
		return fake.resumeFunc(ctx, request)
	}
	return usecase.PersonaGenerationPhaseCommandResult{}, nil
}
func (fake fakePersonaGenerationPhaseUsecase) RetryPersonaGenerationPhase(ctx context.Context, request usecase.RetryPersonaGenerationPhaseRequest) (usecase.PersonaGenerationPhaseCommandResult, error) {
	if fake.retryFunc != nil {
		return fake.retryFunc(ctx, request)
	}
	return usecase.PersonaGenerationPhaseCommandResult{}, nil
}
func (fake fakePersonaGenerationPhaseUsecase) CancelPersonaGenerationPhase(ctx context.Context, request usecase.CancelPersonaGenerationPhaseRequest) (usecase.PersonaGenerationPhaseCommandResult, error) {
	if fake.cancelFunc != nil {
		return fake.cancelFunc(ctx, request)
	}
	return usecase.PersonaGenerationPhaseCommandResult{}, nil
}
func (fake fakePersonaGenerationPhaseUsecase) GetPersonaGenerationBodyReadiness(ctx context.Context, request usecase.GetPersonaGenerationBodyReadinessRequest) (usecase.PersonaGenerationBodyReadinessResult, error) {
	if fake.getBodyReadinessFunc != nil {
		return fake.getBodyReadinessFunc(ctx, request)
	}
	return usecase.PersonaGenerationBodyReadinessResult{}, nil
}

func TestPersonaGenerationPhaseControllerGetSummaryFormatsTimeAndNormalizesErrorKind(t *testing.T) {
	startedAt := time.Date(2026, 5, 1, 12, 0, 0, 0, time.FixedZone("JST", 9*60*60))
	controller := NewPersonaGenerationPhaseController(fakePersonaGenerationPhaseUsecase{
		getSummaryFunc: func(context.Context, usecase.GetPersonaGenerationPhaseSummaryRequest) (usecase.PersonaGenerationPhaseSummaryResult, error) {
			return usecase.PersonaGenerationPhaseSummaryResult{
				JobID:     10,
				StartedAt: &startedAt,
				ErrorSummary: &usecase.PersonaGenerationPhaseErrorSummary{
					ErrorKind: " PROVIDER_FAILURE ",
				},
			}, nil
		},
	})

	response, err := controller.GetPersonaGenerationPhaseSummary(GetPersonaGenerationPhaseSummaryRequestDTO{JobID: 10})
	if err != nil {
		t.Fatalf("expected success, got %v", err)
	}
	if response.StartedAt == nil || *response.StartedAt != "2026-05-01T03:00:00Z" {
		t.Fatalf("expected UTC RFC3339 time, got %#v", response.StartedAt)
	}
	if response.ErrorSummary == nil || response.ErrorSummary.ErrorKind != "provider_failure" {
		t.Fatalf("expected normalized error kind, got %#v", response.ErrorSummary)
	}
}

func TestPersonaGenerationPhaseControllerGetSummaryReturnsEmptySkippedReasonsArray(t *testing.T) {
	controller := NewPersonaGenerationPhaseController(fakePersonaGenerationPhaseUsecase{
		getSummaryFunc: func(context.Context, usecase.GetPersonaGenerationPhaseSummaryRequest) (usecase.PersonaGenerationPhaseSummaryResult, error) {
			return usecase.PersonaGenerationPhaseSummaryResult{}, nil
		},
	})

	response, err := controller.GetPersonaGenerationPhaseSummary(GetPersonaGenerationPhaseSummaryRequestDTO{JobID: 10})
	if err != nil {
		t.Fatalf("expected success, got %v", err)
	}
	payload, err := json.Marshal(response)
	if err != nil {
		t.Fatalf("expected marshal success, got %v", err)
	}
	if strings.Contains(string(payload), `"skippedReasons":null`) {
		t.Fatalf("expected skippedReasons array, got %s", payload)
	}
	if !strings.Contains(string(payload), `"skippedReasons":[]`) {
		t.Fatalf("expected empty skippedReasons array, got %s", payload)
	}
	if strings.Contains(string(payload), `"evidenceRefs":null`) {
		t.Fatalf("expected evidenceRefs array, got %s", payload)
	}
	if !strings.Contains(string(payload), `"evidenceRefs":[]`) {
		t.Fatalf("expected empty evidenceRefs array, got %s", payload)
	}
}

func TestPersonaGenerationPhaseControllerGetBodyReadinessReturnsEmptyEvidenceRefsArray(t *testing.T) {
	controller := NewPersonaGenerationPhaseController(fakePersonaGenerationPhaseUsecase{
		getBodyReadinessFunc: func(context.Context, usecase.GetPersonaGenerationBodyReadinessRequest) (usecase.PersonaGenerationBodyReadinessResult, error) {
			return usecase.PersonaGenerationBodyReadinessResult{}, nil
		},
	})

	response, err := controller.GetPersonaGenerationBodyReadiness(GetPersonaGenerationBodyReadinessRequestDTO{JobID: 10})
	if err != nil {
		t.Fatalf("expected success, got %v", err)
	}
	payload, err := json.Marshal(response)
	if err != nil {
		t.Fatalf("expected marshal success, got %v", err)
	}
	if strings.Contains(string(payload), `"evidenceRefs":null`) {
		t.Fatalf("expected evidenceRefs array, got %s", payload)
	}
	if !strings.Contains(string(payload), `"evidenceRefs":[]`) {
		t.Fatalf("expected empty evidenceRefs array, got %s", payload)
	}
}

func TestPersonaGenerationPhaseControllerCancelForwardsRequest(t *testing.T) {
	var captured usecase.CancelPersonaGenerationPhaseRequest
	controller := NewPersonaGenerationPhaseController(fakePersonaGenerationPhaseUsecase{
		cancelFunc: func(_ context.Context, request usecase.CancelPersonaGenerationPhaseRequest) (usecase.PersonaGenerationPhaseCommandResult, error) {
			captured = request
			return usecase.PersonaGenerationPhaseCommandResult{JobID: 77}, nil
		},
	})

	_, err := controller.CancelPersonaGenerationPhase(CancelPersonaGenerationPhaseRequestDTO{
		JobID:      77,
		PhaseRunID: 91,
	})
	if err != nil {
		t.Fatalf("expected success, got %v", err)
	}
	if captured.JobID != 77 || captured.PhaseRunID != 91 {
		t.Fatalf("expected forwarded request, got %#v", captured)
	}
}

func TestPersonaGenerationPhaseControllerStartReturnsRedactedRejectedResultOnUsecaseError(t *testing.T) {
	controller := NewPersonaGenerationPhaseController(fakePersonaGenerationPhaseUsecase{
		startFunc: func(context.Context, usecase.StartPersonaGenerationPhaseRequest) (usecase.PersonaGenerationPhaseCommandResult, error) {
			return usecase.PersonaGenerationPhaseCommandResult{
				JobID:      2010,
				PhaseState: "rejected",
				ErrorSummary: &usecase.PersonaGenerationPhaseErrorSummary{
					ErrorKind:  usecase.PersonaGenerationPhaseErrorKindTerminalJob,
					IsRedacted: true,
				},
			}, errors.New("terminal job")
		},
	})

	response, err := controller.StartPersonaGenerationPhase(StartPersonaGenerationPhaseRequestDTO{JobID: 2010})
	if err == nil {
		t.Fatal("expected error")
	}
	if response.ErrorSummary == nil || response.ErrorSummary.ErrorKind != "terminal_job" {
		t.Fatalf("expected rejected response with normalized error kind, got %#v", response.ErrorSummary)
	}
	if !strings.Contains(err.Error(), "start persona generation phase") {
		t.Fatalf("expected wrapped error, got %v", err)
	}
}

func TestPersonaGenerationPhaseControllerStartMapsCommandResponseDTO(t *testing.T) {
	startedAt := time.Date(2026, 5, 2, 8, 30, 0, 0, time.FixedZone("JST", 9*60*60))
	finishedAt := startedAt.Add(3 * time.Minute)
	controller := NewPersonaGenerationPhaseController(fakePersonaGenerationPhaseUsecase{
		startFunc: func(context.Context, usecase.StartPersonaGenerationPhaseRequest) (usecase.PersonaGenerationPhaseCommandResult, error) {
			phaseRunID := int64(44)
			targetSnapshotID := "persona-target-44"
			return usecase.PersonaGenerationPhaseCommandResult{
				JobID:        3001,
				CurrentPhase: "persona_generation",
				PhaseState:   "running",
				PhaseRunID:   &phaseRunID,
				StartedAt:    &startedAt,
				FinishedAt:   &finishedAt,
				Progress: usecase.PersonaGenerationPhaseProgressSummary{
					Percent:        50,
					ProcessedCount: 3,
					TotalCount:     6,
					TargetCount:    6,
					CurrentStep:    "generating",
				},
				TargetSummary: usecase.PersonaGenerationTargetSummary{
					TargetCount:            6,
					CommonPersonaHitCount:  1,
					CommonPersonaMissCount: 5,
					SkippedCount:           1,
					SkippedReasons:         []string{"orphan_npc_reference"},
					TargetSnapshotID:       &targetSnapshotID,
					TargetSnapshotDigest:   "sha256:target",
				},
				Execution: usecase.PersonaGenerationExecutionSummary{
					CredentialRef: "cred-1",
					Provider:      "xai",
					Model:         "grok-2",
					ExecutionMode: "single_request",
					PromptDigest:  "sha256:prompt",
					InputCount:    6,
					OutputCount:   3,
					EvidenceRefs:  []string{"evidence:npc:1"},
				},
				ResultSummary: &usecase.PersonaGenerationPhaseResultSummary{
					GeneratedCount:          3,
					FailedCount:             0,
					PersonaCount:            3,
					MissingCount:            0,
					SnapshotID:              "snapshot-44",
					SnapshotDigest:          "sha256:snapshot",
					SnapshotReferenceStatus: "ready",
					BodyReadiness:           true,
				},
				Retryable: false,
				ErrorSummary: &usecase.PersonaGenerationPhaseErrorSummary{
					ErrorKind:  " SECRET_REDACTED ",
					Reason:     "redacted public summary",
					Retryable:  false,
					IsRedacted: true,
				},
			}, nil
		},
	})

	response, err := controller.StartPersonaGenerationPhase(StartPersonaGenerationPhaseRequestDTO{JobID: 3001})
	if err != nil {
		t.Fatalf("expected success, got %v", err)
	}
	if response.StartedAt == nil || *response.StartedAt != "2026-05-01T23:30:00Z" {
		t.Fatalf("expected UTC startedAt, got %#v", response.StartedAt)
	}
	if response.FinishedAt == nil || *response.FinishedAt != "2026-05-01T23:33:00Z" {
		t.Fatalf("expected UTC finishedAt, got %#v", response.FinishedAt)
	}
	if response.TargetSummary.TargetSnapshotID == nil || *response.TargetSummary.TargetSnapshotID != "persona-target-44" {
		t.Fatalf("expected target snapshot id, got %#v", response.TargetSummary.TargetSnapshotID)
	}
	if response.Execution.CredentialRef != "cred-1" || response.Execution.Provider != "xai" {
		t.Fatalf("expected execution summary, got %#v", response.Execution)
	}
	if response.ResultSummary == nil || response.ResultSummary.SnapshotID != "snapshot-44" || !response.ResultSummary.BodyReadiness {
		t.Fatalf("expected result summary, got %#v", response.ResultSummary)
	}
	if response.ErrorSummary == nil || response.ErrorSummary.ErrorKind != "secret_redacted" || !response.ErrorSummary.IsRedacted {
		t.Fatalf("expected normalized redacted error summary, got %#v", response.ErrorSummary)
	}
	if response.Retryable {
		t.Fatalf("expected retryable to be false, got %#v", response)
	}
}

func TestPersonaGenerationPhaseControllerGetBodyReadinessWrapsUsecaseErrorAndKeepsPayload(t *testing.T) {
	controller := NewPersonaGenerationPhaseController(fakePersonaGenerationPhaseUsecase{
		getBodyReadinessFunc: func(context.Context, usecase.GetPersonaGenerationBodyReadinessRequest) (usecase.PersonaGenerationBodyReadinessResult, error) {
			return usecase.PersonaGenerationBodyReadinessResult{
				JobID:      2010,
				PhaseState: "recoverable_failed",
				InputSummary: usecase.PersonaGenerationBodyReadinessInputSummary{
					PersonaCount: 3,
					MissingCount: 1,
				},
			}, errors.New("terminal job")
		},
	})

	response, err := controller.GetPersonaGenerationBodyReadiness(GetPersonaGenerationBodyReadinessRequestDTO{JobID: 2010})
	if err == nil {
		t.Fatal("expected error")
	}
	if response.JobID != 2010 {
		t.Fatalf("expected body readiness payload with jobId, got %#v", response)
	}
	if !strings.Contains(err.Error(), "get persona generation body readiness") {
		t.Fatalf("expected wrapped error, got %v", err)
	}
}
