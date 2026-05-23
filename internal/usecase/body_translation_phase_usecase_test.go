package usecase

import (
	"context"
	"errors"
	"strings"
	"testing"

	"aitranslationenginejp/internal/service"
)

type fakeBodyTranslationPhaseServicePort struct {
	summary         service.BodyTranslationPhaseSummaryReadModel
	command         service.BodyTranslationPhaseCommandReadModel
	outputReadiness service.BodyTranslationOutputReadinessReadModel
	saveAI          service.PhaseAISettingsReadModel
	err             error
}

func (fake *fakeBodyTranslationPhaseServicePort) ReadSummary(context.Context, int64) (service.BodyTranslationPhaseSummaryReadModel, error) {
	return fake.summary, fake.err
}
func (fake *fakeBodyTranslationPhaseServicePort) StartPhase(context.Context, int64) (service.BodyTranslationPhaseCommandReadModel, error) {
	return fake.command, fake.err
}
func (fake *fakeBodyTranslationPhaseServicePort) PausePhase(context.Context, int64, int64) (service.BodyTranslationPhaseCommandReadModel, error) {
	return fake.command, fake.err
}
func (fake *fakeBodyTranslationPhaseServicePort) ResumePhase(context.Context, int64, int64) (service.BodyTranslationPhaseCommandReadModel, error) {
	return fake.command, fake.err
}
func (fake *fakeBodyTranslationPhaseServicePort) RetryPhase(context.Context, int64, int64) (service.BodyTranslationPhaseCommandReadModel, error) {
	return fake.command, fake.err
}
func (fake *fakeBodyTranslationPhaseServicePort) CancelPhase(context.Context, int64, int64) (service.BodyTranslationPhaseCommandReadModel, error) {
	return fake.command, fake.err
}
func (fake *fakeBodyTranslationPhaseServicePort) ReadOutputReadiness(context.Context, int64) (service.BodyTranslationOutputReadinessReadModel, error) {
	return fake.outputReadiness, fake.err
}
func (fake *fakeBodyTranslationPhaseServicePort) SaveAISettings(context.Context, service.PhaseAISettingsSelection) (service.PhaseAISettingsReadModel, error) {
	return fake.saveAI, fake.err
}

func TestBodyTranslationPhaseUsecaseSaveAISettingsMapsReadModel(t *testing.T) {
	port := &fakeBodyTranslationPhaseServicePort{saveAI: service.PhaseAISettingsReadModel{
		JobID:            12,
		PhaseID:          "text_translation",
		Provider:         "openai",
		Model:            "gpt-5.4-mini",
		CredentialStatus: "configured",
		ExecutionMode:    "sync",
		BatchMode:        "unsupported",
		ModelListStatus:  "success",
	}}
	uc := NewBodyTranslationPhaseUsecase(port)

	result, err := uc.SaveBodyTranslationPhaseAISettings(context.Background(), SaveBodyTranslationPhaseAISettingsRequest{
		JobID:         12,
		Provider:      "openai",
		Model:         "gpt-5.4-mini",
		ExecutionMode: "sync",
		BatchMode:     "unsupported",
	})
	if err != nil {
		t.Fatalf("expected save ai settings success: %v", err)
	}
	if result.Provider != "openai" || result.Model != "gpt-5.4-mini" {
		t.Fatalf("expected mapped public ai settings response, got %#v", result)
	}
	if result.CredentialStatus != "configured" || result.ModelListStatus != "success" {
		t.Fatalf("expected status fields from service read model, got %#v", result)
	}
}

func TestBodyTranslationPhaseUsecaseSaveAISettingsWrapsServiceError(t *testing.T) {
	port := &fakeBodyTranslationPhaseServicePort{err: errors.New("save failed")}
	uc := NewBodyTranslationPhaseUsecase(port)

	_, err := uc.SaveBodyTranslationPhaseAISettings(context.Background(), SaveBodyTranslationPhaseAISettingsRequest{
		JobID:         12,
		Provider:      "openai",
		Model:         "gpt-5.4-mini",
		ExecutionMode: "sync",
		BatchMode:     "unsupported",
	})
	if err == nil {
		t.Fatal("expected wrapped error")
	}
	if !strings.Contains(err.Error(), "save body translation phase ai settings") {
		t.Fatalf("expected wrapped error context, got %v", err)
	}
}
