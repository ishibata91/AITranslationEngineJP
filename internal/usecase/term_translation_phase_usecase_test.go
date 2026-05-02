package usecase

import (
	"context"
	"errors"
	"strings"
	"testing"
	"time"

	"aitranslationenginejp/internal/service"
)

type fakeTermTranslationPhaseService struct {
	readSummaryFunc            func(context.Context, int64) (service.TermTranslationPhaseSummaryReadModel, error)
	startPhaseFunc             func(context.Context, int64) (service.TermTranslationPhaseCommandReadModel, error)
	pausePhaseFunc             func(context.Context, int64, int64) (service.TermTranslationPhaseCommandReadModel, error)
	resumePhaseFunc            func(context.Context, int64, int64) (service.TermTranslationPhaseCommandReadModel, error)
	retryPhaseFunc             func(context.Context, int64, int64) (service.TermTranslationPhaseCommandReadModel, error)
	readNextPhaseReadinessFunc func(context.Context, int64) (service.TermTranslationNextPhaseReadinessReadModel, error)
}

func (fake fakeTermTranslationPhaseService) ReadSummary(ctx context.Context, jobID int64) (service.TermTranslationPhaseSummaryReadModel, error) {
	if fake.readSummaryFunc != nil {
		return fake.readSummaryFunc(ctx, jobID)
	}
	return service.TermTranslationPhaseSummaryReadModel{}, nil
}
func (fake fakeTermTranslationPhaseService) StartPhase(ctx context.Context, jobID int64) (service.TermTranslationPhaseCommandReadModel, error) {
	if fake.startPhaseFunc != nil {
		return fake.startPhaseFunc(ctx, jobID)
	}
	return service.TermTranslationPhaseCommandReadModel{}, nil
}
func (fake fakeTermTranslationPhaseService) PausePhase(ctx context.Context, jobID int64, phaseRunID int64) (service.TermTranslationPhaseCommandReadModel, error) {
	if fake.pausePhaseFunc != nil {
		return fake.pausePhaseFunc(ctx, jobID, phaseRunID)
	}
	return service.TermTranslationPhaseCommandReadModel{}, nil
}
func (fake fakeTermTranslationPhaseService) ResumePhase(ctx context.Context, jobID int64, phaseRunID int64) (service.TermTranslationPhaseCommandReadModel, error) {
	if fake.resumePhaseFunc != nil {
		return fake.resumePhaseFunc(ctx, jobID, phaseRunID)
	}
	return service.TermTranslationPhaseCommandReadModel{}, nil
}
func (fake fakeTermTranslationPhaseService) RetryPhase(ctx context.Context, jobID int64, phaseRunID int64) (service.TermTranslationPhaseCommandReadModel, error) {
	if fake.retryPhaseFunc != nil {
		return fake.retryPhaseFunc(ctx, jobID, phaseRunID)
	}
	return service.TermTranslationPhaseCommandReadModel{}, nil
}
func (fake fakeTermTranslationPhaseService) ReadNextPhaseReadiness(ctx context.Context, jobID int64) (service.TermTranslationNextPhaseReadinessReadModel, error) {
	if fake.readNextPhaseReadinessFunc != nil {
		return fake.readNextPhaseReadinessFunc(ctx, jobID)
	}
	return service.TermTranslationNextPhaseReadinessReadModel{}, nil
}

func TestTermTranslationPhaseUsecaseGetSummaryMapsPointersAndNormalizesErrorKind(t *testing.T) {
	startedAt := time.Date(2026, 5, 1, 10, 0, 0, 0, time.UTC)
	snapshot := "digest-v1"
	usecase := NewTermTranslationPhaseUsecase(fakeTermTranslationPhaseService{
		readSummaryFunc: func(context.Context, int64) (service.TermTranslationPhaseSummaryReadModel, error) {
			return service.TermTranslationPhaseSummaryReadModel{
				JobID:        20,
				CurrentPhase: "term_translation",
				PhaseState:   "running",
				StartedAt:    &startedAt,
				Execution: service.TermTranslationExecutionConfigReadModel{
					SnapshotDigest: &snapshot,
				},
				ErrorSummary: &service.TermTranslationPhaseErrorReadModel{ErrorKind: "  PROVIDER_FAILURE  "},
			}, nil
		},
	})

	result, err := usecase.GetTermTranslationPhaseSummary(context.Background(), GetTermTranslationPhaseSummaryRequest{JobID: 20})
	if err != nil {
		t.Fatalf("expected success, got %v", err)
	}
	if result.Execution.SnapshotDigest == nil || *result.Execution.SnapshotDigest != "digest-v1" {
		t.Fatalf("expected snapshot digest mapping, got %#v", result.Execution)
	}
	if result.ErrorSummary == nil || result.ErrorSummary.ErrorKind != "provider_failure" {
		t.Fatalf("expected normalized error kind, got %#v", result.ErrorSummary)
	}
	if result.StartedAt == &startedAt {
		t.Fatal("expected startedAt pointer clone")
	}
}

func TestTermTranslationPhaseUsecaseGetNextPhaseReadinessNormalizesErrorKind(t *testing.T) {
	reason := "term phase is not completed"
	usecase := NewTermTranslationPhaseUsecase(fakeTermTranslationPhaseService{
		readNextPhaseReadinessFunc: func(context.Context, int64) (service.TermTranslationNextPhaseReadinessReadModel, error) {
			return service.TermTranslationNextPhaseReadinessReadModel{
				JobID:             31,
				CanStartNextPhase: false,
				BlockedReason:     &reason,
				ErrorKind:         "  TERM_PHASE_INCOMPLETE ",
			}, nil
		},
	})

	result, err := usecase.GetTermTranslationNextPhaseReadiness(context.Background(), GetTermTranslationNextPhaseReadinessRequest{JobID: 31})
	if err != nil {
		t.Fatalf("expected success, got %v", err)
	}
	if result.ErrorKind != "term_phase_incomplete" {
		t.Fatalf("expected normalized error kind, got %#v", result)
	}
	if result.BlockedReason == &reason {
		t.Fatal("expected blocked reason pointer clone")
	}
}

func TestTermTranslationPhaseUsecaseStartWrapsServiceError(t *testing.T) {
	usecase := NewTermTranslationPhaseUsecase(fakeTermTranslationPhaseService{
		startPhaseFunc: func(context.Context, int64) (service.TermTranslationPhaseCommandReadModel, error) {
			return service.TermTranslationPhaseCommandReadModel{}, errors.New("provider down")
		},
	})

	_, err := usecase.StartTermTranslationPhase(context.Background(), StartTermTranslationPhaseRequest{JobID: 7})
	if err == nil {
		t.Fatal("expected error")
	}
	if !strings.Contains(err.Error(), "start term translation phase") {
		t.Fatalf("expected wrapped error, got %v", err)
	}
}
