package wails

import (
	"context"
	"errors"
	"strings"
	"testing"
	"time"

	"aitranslationenginejp/internal/usecase"
)

type fakeTermTranslationPhaseUsecase struct {
	getSummaryFunc            func(context.Context, usecase.GetTermTranslationPhaseSummaryRequest) (usecase.TermTranslationPhaseSummaryResult, error)
	startFunc                 func(context.Context, usecase.StartTermTranslationPhaseRequest) (usecase.TermTranslationPhaseCommandResult, error)
	pauseFunc                 func(context.Context, usecase.PauseTermTranslationPhaseRequest) (usecase.TermTranslationPhaseCommandResult, error)
	resumeFunc                func(context.Context, usecase.ResumeTermTranslationPhaseRequest) (usecase.TermTranslationPhaseCommandResult, error)
	retryFunc                 func(context.Context, usecase.RetryTermTranslationPhaseRequest) (usecase.TermTranslationPhaseCommandResult, error)
	getNextPhaseReadinessFunc func(context.Context, usecase.GetTermTranslationNextPhaseReadinessRequest) (usecase.TermTranslationNextPhaseReadinessResult, error)
}

func (fake fakeTermTranslationPhaseUsecase) GetTermTranslationPhaseSummary(ctx context.Context, request usecase.GetTermTranslationPhaseSummaryRequest) (usecase.TermTranslationPhaseSummaryResult, error) {
	if fake.getSummaryFunc != nil {
		return fake.getSummaryFunc(ctx, request)
	}
	return usecase.TermTranslationPhaseSummaryResult{}, nil
}
func (fake fakeTermTranslationPhaseUsecase) StartTermTranslationPhase(ctx context.Context, request usecase.StartTermTranslationPhaseRequest) (usecase.TermTranslationPhaseCommandResult, error) {
	if fake.startFunc != nil {
		return fake.startFunc(ctx, request)
	}
	return usecase.TermTranslationPhaseCommandResult{}, nil
}
func (fake fakeTermTranslationPhaseUsecase) PauseTermTranslationPhase(ctx context.Context, request usecase.PauseTermTranslationPhaseRequest) (usecase.TermTranslationPhaseCommandResult, error) {
	if fake.pauseFunc != nil {
		return fake.pauseFunc(ctx, request)
	}
	return usecase.TermTranslationPhaseCommandResult{}, nil
}
func (fake fakeTermTranslationPhaseUsecase) ResumeTermTranslationPhase(ctx context.Context, request usecase.ResumeTermTranslationPhaseRequest) (usecase.TermTranslationPhaseCommandResult, error) {
	if fake.resumeFunc != nil {
		return fake.resumeFunc(ctx, request)
	}
	return usecase.TermTranslationPhaseCommandResult{}, nil
}
func (fake fakeTermTranslationPhaseUsecase) RetryTermTranslationPhase(ctx context.Context, request usecase.RetryTermTranslationPhaseRequest) (usecase.TermTranslationPhaseCommandResult, error) {
	if fake.retryFunc != nil {
		return fake.retryFunc(ctx, request)
	}
	return usecase.TermTranslationPhaseCommandResult{}, nil
}
func (fake fakeTermTranslationPhaseUsecase) GetTermTranslationNextPhaseReadiness(ctx context.Context, request usecase.GetTermTranslationNextPhaseReadinessRequest) (usecase.TermTranslationNextPhaseReadinessResult, error) {
	if fake.getNextPhaseReadinessFunc != nil {
		return fake.getNextPhaseReadinessFunc(ctx, request)
	}
	return usecase.TermTranslationNextPhaseReadinessResult{}, nil
}

func TestTermTranslationPhaseControllerGetSummaryFormatsTimeAndNormalizesErrorKind(t *testing.T) {
	startedAt := time.Date(2026, 5, 1, 12, 0, 0, 0, time.FixedZone("JST", 9*60*60))
	controller := NewTermTranslationPhaseController(fakeTermTranslationPhaseUsecase{
		getSummaryFunc: func(context.Context, usecase.GetTermTranslationPhaseSummaryRequest) (usecase.TermTranslationPhaseSummaryResult, error) {
			return usecase.TermTranslationPhaseSummaryResult{
				JobID:        10,
				StartedAt:    &startedAt,
				ErrorSummary: &usecase.TermTranslationPhaseErrorSummary{ErrorKind: " PROVIDER_FAILURE "},
			}, nil
		},
	})

	response, err := controller.GetTermTranslationPhaseSummary(GetTermTranslationPhaseSummaryRequestDTO{JobID: 10})
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

func TestTermTranslationPhaseControllerStartForwardsRequest(t *testing.T) {
	var captured usecase.StartTermTranslationPhaseRequest
	controller := NewTermTranslationPhaseController(fakeTermTranslationPhaseUsecase{
		startFunc: func(_ context.Context, request usecase.StartTermTranslationPhaseRequest) (usecase.TermTranslationPhaseCommandResult, error) {
			captured = request
			return usecase.TermTranslationPhaseCommandResult{JobID: 77}, nil
		},
	})

	_, err := controller.StartTermTranslationPhase(StartTermTranslationPhaseRequestDTO{JobID: 77})
	if err != nil {
		t.Fatalf("expected success, got %v", err)
	}
	if captured.JobID != 77 {
		t.Fatalf("expected forwarded job id, got %#v", captured)
	}
}

func TestTermTranslationPhaseControllerGetNextPhaseReadinessWrapsUsecaseError(t *testing.T) {
	controller := NewTermTranslationPhaseController(fakeTermTranslationPhaseUsecase{
		getNextPhaseReadinessFunc: func(context.Context, usecase.GetTermTranslationNextPhaseReadinessRequest) (usecase.TermTranslationNextPhaseReadinessResult, error) {
			return usecase.TermTranslationNextPhaseReadinessResult{}, errors.New("storage down")
		},
	})

	_, err := controller.GetTermTranslationNextPhaseReadiness(GetTermTranslationNextPhaseReadinessRequestDTO{JobID: 10})
	if err == nil {
		t.Fatal("expected error")
	}
	if !strings.Contains(err.Error(), "get term translation next phase readiness") {
		t.Fatalf("expected wrapped error, got %v", err)
	}
}
