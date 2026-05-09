package usecase

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"log/slog"
	"strings"
	"testing"

	jobmanagementservice "aitranslationenginejp/internal/service"
)

type fakeTranslationJobManagementService struct {
	listFunc       func(context.Context) (jobmanagementservice.TranslationJobManagementListReadModel, error)
	getDetailFunc  func(context.Context, int64) (jobmanagementservice.TranslationJobManagementJobDetailReadModel, error)
	deleteFunc     func(context.Context, int64) (jobmanagementservice.TranslationJobManagementActionReadModel, error)
	requestStopFun func(context.Context, int64) (jobmanagementservice.TranslationJobManagementActionReadModel, error)
	resumeFunc     func(context.Context, int64) (jobmanagementservice.TranslationJobManagementActionReadModel, error)
}

func (fake fakeTranslationJobManagementService) ListIncompleteJobs(ctx context.Context) (jobmanagementservice.TranslationJobManagementListReadModel, error) {
	return fake.listFunc(ctx)
}
func (fake fakeTranslationJobManagementService) GetJobDetail(ctx context.Context, jobID int64) (jobmanagementservice.TranslationJobManagementJobDetailReadModel, error) {
	return fake.getDetailFunc(ctx, jobID)
}
func (fake fakeTranslationJobManagementService) DeleteJob(ctx context.Context, jobID int64) (jobmanagementservice.TranslationJobManagementActionReadModel, error) {
	return fake.deleteFunc(ctx, jobID)
}
func (fake fakeTranslationJobManagementService) RequestStop(ctx context.Context, jobID int64) (jobmanagementservice.TranslationJobManagementActionReadModel, error) {
	return fake.requestStopFun(ctx, jobID)
}
func (fake fakeTranslationJobManagementService) ResumeJob(ctx context.Context, jobID int64) (jobmanagementservice.TranslationJobManagementActionReadModel, error) {
	return fake.resumeFunc(ctx, jobID)
}

func TestTranslationJobManagementUsecaseDeleteJobMapsDTOAndNullability(t *testing.T) {
	deleted := int64(7)
	usecase := NewTranslationJobManagementUsecase(fakeTranslationJobManagementService{
		deleteFunc: func(context.Context, int64) (jobmanagementservice.TranslationJobManagementActionReadModel, error) {
			return jobmanagementservice.TranslationJobManagementActionReadModel{
				Message:        "ok",
				Tone:           "success",
				DeletedJobID:   &deleted,
				ReasonCategory: "",
			}, nil
		},
	})

	result, err := usecase.DeleteJob(context.Background(), TranslationJobManagementDeleteRequest{JobID: 7})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if result.DeletedJobID == nil || *result.DeletedJobID != 7 {
		t.Fatalf("expected deleted job id to be mapped, got %#v", result)
	}
	if result.Detail != nil {
		t.Fatalf("expected nil detail to stay nil, got %#v", result.Detail)
	}
}

func TestTranslationJobManagementUsecaseListIncompleteJobsWrapsError(t *testing.T) {
	usecase := NewTranslationJobManagementUsecase(fakeTranslationJobManagementService{
		listFunc: func(context.Context) (jobmanagementservice.TranslationJobManagementListReadModel, error) {
			return jobmanagementservice.TranslationJobManagementListReadModel{}, errors.New("db down")
		},
	})

	_, err := usecase.ListIncompleteJobs(context.Background())
	if err == nil {
		t.Fatal("expected error")
	}
	if !strings.Contains(err.Error(), "list incomplete translation jobs") {
		t.Fatalf("expected wrapped context, got %v", err)
	}
}

func TestTranslationJobManagementUsecaseGetDetailMapsReasonCategory(t *testing.T) {
	usecase := NewTranslationJobManagementUsecase(fakeTranslationJobManagementService{
		getDetailFunc: func(context.Context, int64) (jobmanagementservice.TranslationJobManagementJobDetailReadModel, error) {
			return jobmanagementservice.TranslationJobManagementJobDetailReadModel{
				TranslationJobManagementJobSummaryReadModel: jobmanagementservice.TranslationJobManagementJobSummaryReadModel{
					JobID:         44,
					JobState:      "Paused",
					JobStateLabel: "中断中",
				},
				RuntimeSummary: jobmanagementservice.TranslationJobManagementProtectedSettingSummaryReadModel{
					ProviderLabel:        "openai",
					ModelLabel:           "gpt-5",
					ExecutionModeLabel:   "batch",
					CredentialState:      "configured",
					CredentialStateLabel: "設定済み",
				},
				ResumeBlockedReasons: []jobmanagementservice.TranslationJobManagementBlockedReasonReadModel{{
					Category: "cache_missing",
					Title:    "入力キャッシュが欠落しています",
					Detail:   "再構築してください",
				}},
			}, nil
		},
	})

	result, err := usecase.GetJobDetail(context.Background(), TranslationJobManagementGetDetailRequest{JobID: 44})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if len(result.ResumeBlockedReasons) != 1 || result.ResumeBlockedReasons[0].Category != "cache_missing" {
		t.Fatalf("expected blocked reason category mapping, got %#v", result.ResumeBlockedReasons)
	}
}

func TestLogPhaseStateCommandUsesActualReadModelState(t *testing.T) {
	phaseRunID := int64(12)

	accepted := capturePhaseStateCommandLog(t, func(ctx context.Context) {
		logPhaseStateCommand(ctx, "phase_start", "backend.usecase.test", 7, &phaseRunID, "idle_ready", "running", "", nil)
	})
	if accepted["result"] != "accepted" {
		t.Fatalf("expected accepted result, got %#v", accepted)
	}
	if accepted["before_state"] != "idle_ready" || accepted["after_state"] != "running" {
		t.Fatalf("expected accepted state log to use actual before and after state, got %#v", accepted)
	}

	rejected := capturePhaseStateCommandLog(t, func(ctx context.Context) {
		logPhaseStateCommand(ctx, "phase_pause", "backend.usecase.test", 7, &phaseRunID, "completed", "completed", "invalid_phase_state", nil)
	})
	if rejected["result"] != "rejected" || rejected["reason"] != "invalid_phase_state" {
		t.Fatalf("expected rejected result and reason, got %#v", rejected)
	}
	if rejected["before_state"] != "completed" || rejected["after_state"] != "completed" {
		t.Fatalf("expected rejected state log to use actual unchanged state, got %#v", rejected)
	}
}

func TestLogPhaseStateCommandFallsBackToUnknownWhenActualStateMissing(t *testing.T) {
	entry := capturePhaseStateCommandLog(t, func(ctx context.Context) {
		logPhaseStateCommand(ctx, "phase_retry", "backend.usecase.test", 7, nil, "", "", "", errors.New("db down"))
	})
	if entry["result"] != "rejected" || entry["reason"] != "service_error" {
		t.Fatalf("expected service error rejection, got %#v", entry)
	}
	if entry["before_state"] != "unknown" || entry["after_state"] != "unknown" {
		t.Fatalf("expected unknown fallback for missing actual state, got %#v", entry)
	}
}

func capturePhaseStateCommandLog(t *testing.T, writeLog func(context.Context)) map[string]string {
	t.Helper()

	var buffer bytes.Buffer
	previousLogger := slog.Default()
	slog.SetDefault(slog.New(slog.NewJSONHandler(&buffer, nil)))
	t.Cleanup(func() {
		slog.SetDefault(previousLogger)
	})

	writeLog(context.Background())

	var raw map[string]any
	if err := json.Unmarshal(buffer.Bytes(), &raw); err != nil {
		t.Fatalf("expected JSON log entry, got %q: %v", buffer.String(), err)
	}
	entry := make(map[string]string, len(raw))
	for key, value := range raw {
		if text, ok := value.(string); ok {
			entry[key] = text
		}
	}
	return entry
}
