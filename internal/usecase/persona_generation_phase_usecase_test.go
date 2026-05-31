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
		Execution:        service.PersonaGenerationExecutionSummaryReadModel{EvidenceRefs: []string{"phase-run:9"}},
		ActionEnablement: service.PersonaGenerationPhaseActionEnablementReadModel{CanStart: true},
	}}
	uc := NewPersonaGenerationPhaseUsecase(port)

	result, err := uc.GetPersonaGenerationPhaseSummary(context.Background(), GetPersonaGenerationPhaseSummaryRequest{JobID: 1})
	if err != nil {
		t.Fatalf("expected summary success: %v", err)
	}
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

func TestPersonaGenerationPhaseUsecase_SaveAISettingsMapsReadModel(t *testing.T) {
	port := &fakePersonaGenerationPhaseServicePort{
		saveAI: service.PhaseAISettingsReadModel{
			JobID:            3,
			PhaseID:          "npc_persona_generation",
			Provider:         "xai",
			Model:            "grok-4",
			CredentialStatus: "configured",
			ExecutionMode:    "sync",
			BatchMode:        "disabled",
			ModelListStatus:  "success",
		},
	}
	uc := NewPersonaGenerationPhaseUsecase(port)

	result, err := uc.SavePersonaGenerationPhaseAISettings(context.Background(), SavePersonaGenerationPhaseAISettingsRequest{
		JobID:         3,
		Provider:      "xai",
		Model:         "grok-4",
		ExecutionMode: "sync",
		BatchMode:     "disabled",
	})
	if err != nil {
		t.Fatalf("expected save ai settings success: %v", err)
	}
	if result.CredentialStatus != "configured" || result.ModelListStatus != "success" {
		t.Fatalf("expected public ai settings status mapping, got %#v", result)
	}
}
