package usecase

import (
	"context"
	"errors"
	"testing"
	"time"

	"aitranslationenginejp/internal/service"
)

type fakePersonaGenerationPhaseServicePort struct {
	summary   service.PersonaGenerationPhaseSummaryReadModel
	command   service.PersonaGenerationPhaseCommandReadModel
	readiness service.PersonaGenerationBodyReadinessReadModel
	saveAI    service.PhaseAISettingsReadModel
	err       error
}

func (fake *fakePersonaGenerationPhaseServicePort) ReadSummary(context.Context, int64) (service.PersonaGenerationPhaseSummaryReadModel, error) {
	return fake.summary, fake.err
}
func (fake *fakePersonaGenerationPhaseServicePort) StartPhase(context.Context, int64) (service.PersonaGenerationPhaseCommandReadModel, error) {
	return fake.command, fake.err
}
func (fake *fakePersonaGenerationPhaseServicePort) PausePhase(context.Context, int64, int64) (service.PersonaGenerationPhaseCommandReadModel, error) {
	return fake.command, fake.err
}
func (fake *fakePersonaGenerationPhaseServicePort) ResumePhase(context.Context, int64, int64) (service.PersonaGenerationPhaseCommandReadModel, error) {
	return fake.command, fake.err
}
func (fake *fakePersonaGenerationPhaseServicePort) RetryPhase(context.Context, int64, int64) (service.PersonaGenerationPhaseCommandReadModel, error) {
	return fake.command, fake.err
}
func (fake *fakePersonaGenerationPhaseServicePort) CancelPhase(context.Context, int64, int64) (service.PersonaGenerationPhaseCommandReadModel, error) {
	return fake.command, fake.err
}
func (fake *fakePersonaGenerationPhaseServicePort) ReadBodyReadiness(context.Context, int64) (service.PersonaGenerationBodyReadinessReadModel, error) {
	return fake.readiness, fake.err
}
func (fake *fakePersonaGenerationPhaseServicePort) SaveAISettings(context.Context, service.PhaseAISettingsSelection) (service.PhaseAISettingsReadModel, error) {
	return fake.saveAI, fake.err
}

func TestPersonaGenerationPhaseUsecase_MapsSummaryDTO(t *testing.T) {
	now := time.Date(2026, 5, 2, 10, 0, 0, 0, time.UTC)
	phaseRunID := int64(9)
	targetSnapshotID := "snap"
	port := &fakePersonaGenerationPhaseServicePort{summary: service.PersonaGenerationPhaseSummaryReadModel{
		JobID: 1, CurrentPhase: "persona_generation", PhaseState: "running", PhaseRunID: &phaseRunID,
		StartedAt:        &now,
		TargetSummary:    service.PersonaGenerationTargetSummaryReadModel{TargetSnapshotID: &targetSnapshotID, SkippedReasons: []string{"orphan"}},
		Execution:        &service.PersonaGenerationExecutionSummaryReadModel{EvidenceRefs: []string{"phase-run:9"}},
		ActionEnablement: service.PersonaGenerationPhaseActionEnablementReadModel{CanStart: true},
	}}
	uc := NewPersonaGenerationPhaseUsecase(port)

	fetchResult, err := uc.GetPersonaGenerationPhaseSummary(context.Background(), GetPersonaGenerationPhaseSummaryRequest{JobID: 1})
	if err != nil {
		t.Fatalf("expected summary success: %v", err)
	}
	result := fetchResult.Summary
	if result.PhaseRunID == nil || *result.PhaseRunID != 9 {
		t.Fatalf("expected phase run id mapping, got %#v", result.PhaseRunID)
	}
	if len(result.TargetSummary.SkippedReasons) != 1 || len(result.Execution.EvidenceRefs) != 1 {
		t.Fatalf("expected list copy mapping, got %#v", result)
	}
	result.TargetSummary.SkippedReasons[0] = "mutated"
	if port.summary.TargetSummary.SkippedReasons[0] != "orphan" {
		t.Fatal("expected deep copy for skipped reasons")
	}
}

func TestPersonaGenerationPhaseUsecase_CommandAndBodyReadinessErrorWrapping(t *testing.T) {
	port := &fakePersonaGenerationPhaseServicePort{
		command:   service.PersonaGenerationPhaseCommandReadModel{JobID: 1, PhaseState: "rejected", ErrorSummary: &service.PersonaGenerationPhaseErrorSummaryReadModel{ErrorKind: "provider_failure", IsRedacted: true}},
		readiness: service.PersonaGenerationBodyReadinessReadModel{JobID: 1},
		err:       errors.New("boom"),
	}
	uc := NewPersonaGenerationPhaseUsecase(port)

	commandResult, commandErr := uc.StartPersonaGenerationPhase(context.Background(), StartPersonaGenerationPhaseRequest{JobID: 1})
	if commandErr == nil || commandResult.ErrorSummary == nil || !commandResult.ErrorSummary.IsRedacted {
		t.Fatalf("expected mapped command result and wrapped error, result=%#v err=%v", commandResult, commandErr)
	}

	_, pauseErr := uc.PausePersonaGenerationPhase(context.Background(), PausePersonaGenerationPhaseRequest{JobID: 1, PhaseRunID: 2})
	if pauseErr == nil {
		t.Fatal("expected pause error wrapping")
	}

	_, readinessErr := uc.GetPersonaGenerationBodyReadiness(context.Background(), GetPersonaGenerationBodyReadinessRequest{JobID: 1})
	if readinessErr == nil {
		t.Fatalf("expected readiness error wrapping, got nil err")
	}
}

// RAEF-UNIT-027: PersonaGenerationPhaseProjection 組み立て論理の境界値証明
// 根拠: design-diff.md G-5-b

// RAEF-UNIT-027: PreviousPhaseLifecycle が直前 term phase run の lifecycle 値を正しく返す。term phase run 未生成時は空文字
func TestPersonaGenerationPhaseProjection_PreviousPhaseLifecycleBoundary(t *testing.T) {
	// G-5-b: PreviousPhaseLifecycle は直前 term JOB_PHASE_RUN の lifecycle 値。phase run 未生成時は空文字。

	// term phase run 存在時は lifecycle 値を返す
	portWithPreviousPhase := &fakePersonaGenerationPhaseServicePort{
		summary: service.PersonaGenerationPhaseSummaryReadModel{
			JobID:                  1,
			JobState:               "running",
			CurrentPhase:           "persona_generation",
			PhaseState:             "idle_ready",
			PreviousPhaseLifecycle: "completed", // term phase run が completed
			TargetSummary:          service.PersonaGenerationTargetSummaryReadModel{TargetCount: 5},
		},
	}
	ucWithPrevious := NewPersonaGenerationPhaseUsecase(portWithPreviousPhase)

	resultWithPrevious, err := ucWithPrevious.GetPersonaGenerationPhaseSummary(context.Background(), GetPersonaGenerationPhaseSummaryRequest{JobID: 1})
	if err != nil {
		t.Fatalf("expected success, got %v", err)
	}
	// term phase run の lifecycle 値 "completed" が PreviousPhaseLifecycle に入ることを証明する
	if resultWithPrevious.Projection.PreviousPhaseLifecycle != "completed" {
		t.Fatalf("expected PreviousPhaseLifecycle=completed when term run exists, got %q", resultWithPrevious.Projection.PreviousPhaseLifecycle)
	}

	// term phase run 未生成時は空文字を返す
	portWithoutPreviousPhase := &fakePersonaGenerationPhaseServicePort{
		summary: service.PersonaGenerationPhaseSummaryReadModel{
			JobID:                  1,
			JobState:               "running",
			CurrentPhase:           "persona_generation",
			PhaseState:             "idle_ready",
			PreviousPhaseLifecycle: "", // term phase run 未生成
			TargetSummary:          service.PersonaGenerationTargetSummaryReadModel{TargetCount: 5},
		},
	}
	ucWithoutPrevious := NewPersonaGenerationPhaseUsecase(portWithoutPreviousPhase)

	resultWithoutPrevious, err := ucWithoutPrevious.GetPersonaGenerationPhaseSummary(context.Background(), GetPersonaGenerationPhaseSummaryRequest{JobID: 1})
	if err != nil {
		t.Fatalf("expected success, got %v", err)
	}
	// term phase run 未生成時は PreviousPhaseLifecycle が空文字になることを証明する
	if resultWithoutPrevious.Projection.PreviousPhaseLifecycle != "" {
		t.Fatalf("expected PreviousPhaseLifecycle=empty string when term run not generated, got %q", resultWithoutPrevious.Projection.PreviousPhaseLifecycle)
	}
}

func TestPersonaGenerationPhaseUsecase_SaveAISettingsMapsReadModel(t *testing.T) {
	port := &fakePersonaGenerationPhaseServicePort{
		saveAI: service.PhaseAISettingsReadModel{
			PhaseType:     "npc_persona_generation",
			Provider:      "xai",
			Model:         "grok-4",
			ExecutionMode: "sync",
			BatchMode:     "disabled",
		},
	}
	uc := NewPersonaGenerationPhaseUsecase(port)

	result, err := uc.SavePersonaGenerationPhaseAISettings(context.Background(), SavePersonaGenerationPhaseAISettingsRequest{
		Provider:      "xai",
		Model:         "grok-4",
		ExecutionMode: "sync",
		BatchMode:     "disabled",
	})
	if err != nil {
		t.Fatalf("expected save ai settings success: %v", err)
	}
	if result.Provider != "xai" || result.Model != "grok-4" {
		t.Fatalf("expected mapped public ai settings response, got %#v", result)
	}
}
