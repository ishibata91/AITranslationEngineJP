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
		PhaseType:     "text_translation",
		Provider:      "openai",
		Model:         "gpt-5.4-mini",
		ExecutionMode: "sync",
		BatchMode:     "unsupported",
	}}
	uc := NewBodyTranslationPhaseUsecase(port)

	result, err := uc.SaveBodyTranslationPhaseAISettings(context.Background(), SaveBodyTranslationPhaseAISettingsRequest{
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
}

// RAEF-UNIT-028〜029: BodyTranslationPhaseProjection 組み立て論理の境界値証明
// 根拠: design-diff.md G-5-c

// fakePersonaReadinessPort は resolvePersonaState の personaReadiness port の fake 実装。
type fakePersonaReadinessPort struct {
	readiness service.PersonaGenerationBodyReadinessReadModel
	err       error
}

func (f *fakePersonaReadinessPort) ReadBodyReadiness(_ context.Context, _ int64) (service.PersonaGenerationBodyReadinessReadModel, error) {
	return f.readiness, f.err
}

// body usecase の projection テスト用コンストラクタ（personaReadiness port を注入できる版）
// WithPersonaReadiness を使って persona readiness port を注入する
func newBodyUsecaseWithPersonaReadiness(
	port *fakeBodyTranslationPhaseServicePort,
	personaPort *fakePersonaReadinessPort,
) *BodyTranslationPhaseUsecase {
	return NewBodyTranslationPhaseUsecase(port).WithPersonaReadiness(personaPort)
}

// RAEF-UNIT-028: persona snapshot 無し時に PersonaBodyReadiness が zero 値（bodyReadiness=false / snapshotReferenceStatus=""）になる
func TestBodyTranslationPhaseProjection_PersonaBodyReadinessIsZeroWhenNoSnapshot(t *testing.T) {
	// G-5-c: PersonaBodyReadiness は snapshot 無し時に zero 値（BodyReadiness=false、SnapshotReferenceStatus=""）
	port := &fakeBodyTranslationPhaseServicePort{
		summary: service.BodyTranslationPhaseSummaryReadModel{
			JobID:      1,
			JobState:   "running",
			PhaseState: "idle_ready",
			Progress:   service.BodyTranslationPhaseProgressReadModel{TargetCount: 5},
		},
	}
	// persona readiness: SnapshotID が空 → snapshot 無し
	personaPort := &fakePersonaReadinessPort{
		readiness: service.PersonaGenerationBodyReadinessReadModel{
			PhaseState: "completed",
			InputSummary: service.PersonaGenerationBodyReadinessInputSummaryReadModel{
				PersonaCount: 0,
				MissingCount: 0,
				SnapshotID:   "", // snapshot 無し
			},
		},
	}
	uc := newBodyUsecaseWithPersonaReadiness(port, personaPort)

	result, err := uc.GetBodyTranslationPhaseSummary(context.Background(), GetBodyTranslationPhaseSummaryRequest{JobID: 1})
	if err != nil {
		t.Fatalf("expected success, got %v", err)
	}

	// snapshot 無し時は PersonaBodyReadiness が zero 値になることを証明する
	if result.Projection.PersonaBodyReadiness.BodyReadiness {
		t.Fatal("expected PersonaBodyReadiness.BodyReadiness=false when no snapshot")
	}
	if result.Projection.PersonaBodyReadiness.SnapshotReferenceStatus != "" {
		t.Fatalf("expected PersonaBodyReadiness.SnapshotReferenceStatus empty, got %q", result.Projection.PersonaBodyReadiness.SnapshotReferenceStatus)
	}
}

// RAEF-UNIT-029: persona snapshot 存在時に PersonaBodyReadiness.BodyReadiness=true かつ SnapshotReferenceStatus=available になる
func TestBodyTranslationPhaseProjection_PersonaBodyReadinessIsAvailableWhenSnapshot(t *testing.T) {
	// G-5-c: persona snapshot 存在時（PersonaCount>0 ∧ MissingCount=0 ∧ SnapshotID 非空）は
	// PersonaBodyReadiness.BodyReadiness=true かつ SnapshotReferenceStatus=available になる
	port := &fakeBodyTranslationPhaseServicePort{
		summary: service.BodyTranslationPhaseSummaryReadModel{
			JobID:      1,
			JobState:   "running",
			PhaseState: "idle_ready",
			Progress:   service.BodyTranslationPhaseProgressReadModel{TargetCount: 5},
		},
	}
	// persona readiness: snapshot 存在（PersonaCount>0 ∧ MissingCount=0 ∧ SnapshotID 非空）
	personaPort := &fakePersonaReadinessPort{
		readiness: service.PersonaGenerationBodyReadinessReadModel{
			PhaseState: "completed",
			InputSummary: service.PersonaGenerationBodyReadinessInputSummaryReadModel{
				PersonaCount: 5,
				MissingCount: 0,
				SnapshotID:   "snapshot-uuid-001", // snapshot 存在
			},
		},
	}
	uc := newBodyUsecaseWithPersonaReadiness(port, personaPort)

	result, err := uc.GetBodyTranslationPhaseSummary(context.Background(), GetBodyTranslationPhaseSummaryRequest{JobID: 1})
	if err != nil {
		t.Fatalf("expected success, got %v", err)
	}

	// snapshot 存在時は BodyReadiness=true かつ SnapshotReferenceStatus=available を証明する
	if !result.Projection.PersonaBodyReadiness.BodyReadiness {
		t.Fatal("expected PersonaBodyReadiness.BodyReadiness=true when snapshot is available")
	}
	if result.Projection.PersonaBodyReadiness.SnapshotReferenceStatus != "available" {
		t.Fatalf("expected PersonaBodyReadiness.SnapshotReferenceStatus=available, got %q", result.Projection.PersonaBodyReadiness.SnapshotReferenceStatus)
	}
}

func TestBodyTranslationPhaseUsecaseSaveAISettingsWrapsServiceError(t *testing.T) {
	port := &fakeBodyTranslationPhaseServicePort{err: errors.New("save failed")}
	uc := NewBodyTranslationPhaseUsecase(port)

	_, err := uc.SaveBodyTranslationPhaseAISettings(context.Background(), SaveBodyTranslationPhaseAISettingsRequest{
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
