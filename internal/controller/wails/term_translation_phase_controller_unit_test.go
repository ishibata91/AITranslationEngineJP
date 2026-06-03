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
	getSummaryFunc            func(context.Context, usecase.GetTermTranslationPhaseSummaryRequest) (usecase.TermTranslationPhaseFetchResult, error)
	startFunc                 func(context.Context, usecase.StartTermTranslationPhaseRequest) (usecase.TermTranslationPhaseCommandResult, error)
	pauseFunc                 func(context.Context, usecase.PauseTermTranslationPhaseRequest) (usecase.TermTranslationPhaseCommandResult, error)
	resumeFunc                func(context.Context, usecase.ResumeTermTranslationPhaseRequest) (usecase.TermTranslationPhaseCommandResult, error)
	retryFunc                 func(context.Context, usecase.RetryTermTranslationPhaseRequest) (usecase.TermTranslationPhaseCommandResult, error)
	getNextPhaseReadinessFunc func(context.Context, usecase.GetTermTranslationNextPhaseReadinessRequest) (usecase.TermTranslationNextPhaseReadinessResult, error)
	saveAISettingsFunc        func(context.Context, usecase.SaveTermTranslationPhaseAISettingsRequest) (usecase.TermTranslationPhaseAISettingsResult, error)
}

func (fake fakeTermTranslationPhaseUsecase) GetTermTranslationPhaseSummary(ctx context.Context, request usecase.GetTermTranslationPhaseSummaryRequest) (usecase.TermTranslationPhaseFetchResult, error) {
	if fake.getSummaryFunc != nil {
		return fake.getSummaryFunc(ctx, request)
	}
	return usecase.TermTranslationPhaseFetchResult{}, nil
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

func (fake fakeTermTranslationPhaseUsecase) SaveTermTranslationPhaseAISettings(ctx context.Context, request usecase.SaveTermTranslationPhaseAISettingsRequest) (usecase.TermTranslationPhaseAISettingsResult, error) {
	if fake.saveAISettingsFunc != nil {
		return fake.saveAISettingsFunc(ctx, request)
	}
	return usecase.TermTranslationPhaseAISettingsResult{}, nil
}

func TestTermTranslationPhaseControllerGetSummaryFormatsTimeAndNormalizesErrorKind(t *testing.T) {
	startedAt := time.Date(2026, 5, 1, 12, 0, 0, 0, time.FixedZone("JST", 9*60*60))
	controller := NewTermTranslationPhaseController(fakeTermTranslationPhaseUsecase{
		getSummaryFunc: func(context.Context, usecase.GetTermTranslationPhaseSummaryRequest) (usecase.TermTranslationPhaseFetchResult, error) {
			return usecase.TermTranslationPhaseFetchResult{
				Summary: usecase.TermTranslationPhaseSummaryResult{
					JobID:        10,
					StartedAt:    &startedAt,
					ErrorSummary: &usecase.TermTranslationPhaseErrorSummary{ErrorKind: " PROVIDER_FAILURE "},
				},
			}, nil
		},
	})

	response, err := controller.GetTermTranslationPhaseSummary(GetTermTranslationPhaseSummaryRequestDTO{JobID: 10})
	if err != nil {
		t.Fatalf("expected success, got %v", err)
	}
	if response.Summary.StartedAt == nil || *response.Summary.StartedAt != "2026-05-01T03:00:00Z" {
		t.Fatalf("expected UTC RFC3339 time, got %#v", response.Summary.StartedAt)
	}
	if response.Summary.ErrorSummary == nil || response.Summary.ErrorSummary.ErrorKind != "provider_failure" {
		t.Fatalf("expected normalized error kind, got %#v", response.Summary.ErrorSummary)
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

func TestTermTranslationPhaseControllerPauseMapsRequestAndResponse(t *testing.T) {
	// PauseTermTranslationPhase の request/response DTO 境界を証明する。
	var captured usecase.PauseTermTranslationPhaseRequest
	controller := NewTermTranslationPhaseController(fakeTermTranslationPhaseUsecase{
		pauseFunc: func(_ context.Context, request usecase.PauseTermTranslationPhaseRequest) (usecase.TermTranslationPhaseCommandResult, error) {
			captured = request
			return usecase.TermTranslationPhaseCommandResult{
				JobID:      request.JobID,
				PhaseState: "Paused",
			}, nil
		},
	})

	response, err := controller.PauseTermTranslationPhase(PauseTermTranslationPhaseRequestDTO{JobID: 8, PhaseRunID: 11})
	if err != nil {
		t.Fatalf("expected pause to succeed: %v", err)
	}
	if captured.JobID != 8 || captured.PhaseRunID != 11 {
		t.Fatalf("expected forwarded request, got %#v", captured)
	}
	if response.JobID != 8 || response.PhaseState != "Paused" {
		t.Fatalf("expected response mapping, got %#v", response)
	}
}

func TestTermTranslationPhaseControllerPauseWrapsError(t *testing.T) {
	// PauseTermTranslationPhase の失敗時に method 境界の wrap を証明する。
	controller := NewTermTranslationPhaseController(fakeTermTranslationPhaseUsecase{
		pauseFunc: func(context.Context, usecase.PauseTermTranslationPhaseRequest) (usecase.TermTranslationPhaseCommandResult, error) {
			return usecase.TermTranslationPhaseCommandResult{}, errors.New("pause rejected")
		},
	})

	_, err := controller.PauseTermTranslationPhase(PauseTermTranslationPhaseRequestDTO{JobID: 8, PhaseRunID: 11})
	if err == nil {
		t.Fatal("expected error")
	}
	if !strings.Contains(err.Error(), "pause term translation phase") {
		t.Fatalf("expected wrapped method context, got %v", err)
	}
}

func TestTermTranslationPhaseControllerResumeMapsRequestAndResponse(t *testing.T) {
	// ResumeTermTranslationPhase の request/response DTO 境界を証明する。
	var captured usecase.ResumeTermTranslationPhaseRequest
	controller := NewTermTranslationPhaseController(fakeTermTranslationPhaseUsecase{
		resumeFunc: func(_ context.Context, request usecase.ResumeTermTranslationPhaseRequest) (usecase.TermTranslationPhaseCommandResult, error) {
			captured = request
			return usecase.TermTranslationPhaseCommandResult{
				JobID:      request.JobID,
				PhaseState: "Running",
			}, nil
		},
	})

	response, err := controller.ResumeTermTranslationPhase(ResumeTermTranslationPhaseRequestDTO{JobID: 9, PhaseRunID: 21})
	if err != nil {
		t.Fatalf("expected resume to succeed: %v", err)
	}
	if captured.JobID != 9 || captured.PhaseRunID != 21 {
		t.Fatalf("expected forwarded request, got %#v", captured)
	}
	if response.JobID != 9 || response.PhaseState != "Running" {
		t.Fatalf("expected response mapping, got %#v", response)
	}
}

func TestTermTranslationPhaseControllerResumeWrapsError(t *testing.T) {
	// ResumeTermTranslationPhase の失敗時に method 境界の wrap を証明する。
	controller := NewTermTranslationPhaseController(fakeTermTranslationPhaseUsecase{
		resumeFunc: func(context.Context, usecase.ResumeTermTranslationPhaseRequest) (usecase.TermTranslationPhaseCommandResult, error) {
			return usecase.TermTranslationPhaseCommandResult{}, errors.New("resume failed")
		},
	})

	_, err := controller.ResumeTermTranslationPhase(ResumeTermTranslationPhaseRequestDTO{JobID: 9, PhaseRunID: 21})
	if err == nil {
		t.Fatal("expected error")
	}
	if !strings.Contains(err.Error(), "resume term translation phase") {
		t.Fatalf("expected wrapped method context, got %v", err)
	}
}

func TestTermTranslationPhaseControllerRetryMapsRequestAndResponse(t *testing.T) {
	// RetryTermTranslationPhase の request/response DTO 境界を証明する。
	var captured usecase.RetryTermTranslationPhaseRequest
	controller := NewTermTranslationPhaseController(fakeTermTranslationPhaseUsecase{
		retryFunc: func(_ context.Context, request usecase.RetryTermTranslationPhaseRequest) (usecase.TermTranslationPhaseCommandResult, error) {
			captured = request
			return usecase.TermTranslationPhaseCommandResult{
				JobID:      request.JobID,
				PhaseState: "Running",
				Retryable:  false,
			}, nil
		},
	})

	response, err := controller.RetryTermTranslationPhase(RetryTermTranslationPhaseRequestDTO{JobID: 10, PhaseRunID: 33})
	if err != nil {
		t.Fatalf("expected retry to succeed: %v", err)
	}
	if captured.JobID != 10 || captured.PhaseRunID != 33 {
		t.Fatalf("expected forwarded request, got %#v", captured)
	}
	if response.JobID != 10 || response.Retryable {
		t.Fatalf("expected response mapping, got %#v", response)
	}
}

func TestTermTranslationPhaseControllerRetryWrapsError(t *testing.T) {
	// RetryTermTranslationPhase の失敗時に method 境界の wrap を証明する。
	controller := NewTermTranslationPhaseController(fakeTermTranslationPhaseUsecase{
		retryFunc: func(context.Context, usecase.RetryTermTranslationPhaseRequest) (usecase.TermTranslationPhaseCommandResult, error) {
			return usecase.TermTranslationPhaseCommandResult{}, errors.New("retry queue blocked")
		},
	})

	_, err := controller.RetryTermTranslationPhase(RetryTermTranslationPhaseRequestDTO{JobID: 10, PhaseRunID: 33})
	if err == nil {
		t.Fatal("expected error")
	}
	if !strings.Contains(err.Error(), "retry term translation phase") {
		t.Fatalf("expected wrapped method context, got %v", err)
	}
}

func TestTermTranslationPhaseControllerSaveAISettingsMapsRequestAndResponse(t *testing.T) {
	// SaveTermTranslationPhaseAISettings の request/response DTO 境界を証明する。
	var captured usecase.SaveTermTranslationPhaseAISettingsRequest
	controller := NewTermTranslationPhaseController(fakeTermTranslationPhaseUsecase{
		saveAISettingsFunc: func(_ context.Context, request usecase.SaveTermTranslationPhaseAISettingsRequest) (usecase.TermTranslationPhaseAISettingsResult, error) {
			captured = request
			return usecase.TermTranslationPhaseAISettingsResult{
				PhaseType:     "word_translation",
				Provider:      request.Provider,
				Model:         request.Model,
				ExecutionMode: request.ExecutionMode,
				BatchMode:     request.BatchMode,
			}, nil
		},
	})

	response, err := controller.SaveTermTranslationPhaseAISettings(SaveTermTranslationPhaseAISettingsRequestDTO{
		Provider:      "xAI",
		Model:         "grok-4",
		ExecutionMode: "batch",
		BatchMode:     "enabled",
	})
	if err != nil {
		t.Fatalf("expected save ai settings to succeed: %v", err)
	}
	if captured.Provider != "xAI" || captured.Model != "grok-4" {
		t.Fatalf("expected forwarded request, got %#v", captured)
	}
	if response.Provider != "xAI" || response.Model != "grok-4" {
		t.Fatalf("expected response mapping, got %#v", response)
	}
}

func TestTermTranslationPhaseControllerSaveAISettingsWrapsError(t *testing.T) {
	// SaveTermTranslationPhaseAISettings の失敗時に method 境界の wrap を証明する。
	controller := NewTermTranslationPhaseController(fakeTermTranslationPhaseUsecase{
		saveAISettingsFunc: func(context.Context, usecase.SaveTermTranslationPhaseAISettingsRequest) (usecase.TermTranslationPhaseAISettingsResult, error) {
			return usecase.TermTranslationPhaseAISettingsResult{}, errors.New("provider unavailable")
		},
	})

	_, err := controller.SaveTermTranslationPhaseAISettings(SaveTermTranslationPhaseAISettingsRequestDTO{})
	if err == nil {
		t.Fatal("expected error")
	}
	if !strings.Contains(err.Error(), "save term translation phase ai settings") {
		t.Fatalf("expected wrapped method context, got %v", err)
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
