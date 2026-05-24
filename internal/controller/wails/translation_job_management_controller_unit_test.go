package wails

import (
	"context"
	"errors"
	"reflect"
	"strings"
	"testing"

	"aitranslationenginejp/internal/usecase"
)

type fakeTranslationJobManagementUsecase struct {
	listFunc      func(context.Context) (usecase.TranslationJobManagementListResult, error)
	getDetailFunc func(context.Context, usecase.TranslationJobManagementGetDetailRequest) (usecase.TranslationJobManagementJobDetail, error)
	deleteFunc    func(context.Context, usecase.TranslationJobManagementDeleteRequest) (usecase.TranslationJobManagementActionResult, error)
	stopFunc      func(context.Context, usecase.TranslationJobManagementActionRequest) (usecase.TranslationJobManagementActionResult, error)
	resumeFunc    func(context.Context, usecase.TranslationJobManagementActionRequest) (usecase.TranslationJobManagementActionResult, error)
}

func (fake fakeTranslationJobManagementUsecase) ListIncompleteJobs(ctx context.Context) (usecase.TranslationJobManagementListResult, error) {
	if fake.listFunc != nil {
		return fake.listFunc(ctx)
	}
	return usecase.TranslationJobManagementListResult{}, nil
}

func (fake fakeTranslationJobManagementUsecase) GetJobDetail(ctx context.Context, req usecase.TranslationJobManagementGetDetailRequest) (usecase.TranslationJobManagementJobDetail, error) {
	if fake.getDetailFunc != nil {
		return fake.getDetailFunc(ctx, req)
	}
	return usecase.TranslationJobManagementJobDetail{}, nil
}

func (fake fakeTranslationJobManagementUsecase) DeleteJob(ctx context.Context, req usecase.TranslationJobManagementDeleteRequest) (usecase.TranslationJobManagementActionResult, error) {
	if fake.deleteFunc != nil {
		return fake.deleteFunc(ctx, req)
	}
	return usecase.TranslationJobManagementActionResult{}, nil
}

func (fake fakeTranslationJobManagementUsecase) RequestStop(ctx context.Context, req usecase.TranslationJobManagementActionRequest) (usecase.TranslationJobManagementActionResult, error) {
	if fake.stopFunc != nil {
		return fake.stopFunc(ctx, req)
	}
	return usecase.TranslationJobManagementActionResult{}, nil
}

func (fake fakeTranslationJobManagementUsecase) ResumeJob(ctx context.Context, req usecase.TranslationJobManagementActionRequest) (usecase.TranslationJobManagementActionResult, error) {
	if fake.resumeFunc != nil {
		return fake.resumeFunc(ctx, req)
	}
	return usecase.TranslationJobManagementActionResult{}, nil
}

func TestTranslationJobManagementControllerGetJobDetailMapsRequestAndResponse(t *testing.T) {
	// GetJobDetail の request/response DTO 写像を証明する。
	var captured usecase.TranslationJobManagementGetDetailRequest
	controller := NewTranslationJobManagementController(fakeTranslationJobManagementUsecase{
		getDetailFunc: func(_ context.Context, req usecase.TranslationJobManagementGetDetailRequest) (usecase.TranslationJobManagementJobDetail, error) {
			captured = req
			return usecase.TranslationJobManagementJobDetail{
				TranslationJobManagementJobSummary: usecase.TranslationJobManagementJobSummary{
					JobID:         req.JobID,
					JobState:      "Paused",
					JobStateLabel: "一時停止",
					Progress: usecase.TranslationJobManagementProgressSummary{
						CurrentPhase: "term_translation",
						Percent:      45,
					},
				},
				CacheState:        "ready",
				CacheStateLabel:   "入力あり",
				DeleteImpactLines: []string{"入力データは残る"},
			}, nil
		},
	})

	response, err := controller.GetJobDetail(TranslationJobManagementGetDetailRequestDTO{JobID: 12})
	if err != nil {
		t.Fatalf("expected get job detail to succeed: %v", err)
	}
	if captured.JobID != 12 {
		t.Fatalf("expected forwarded job id, got %#v", captured)
	}
	if response.JobID != 12 || response.Progress.CurrentPhase != "term_translation" {
		t.Fatalf("expected response mapping, got %#v", response)
	}
	if len(response.DeleteImpactLines) != 1 || response.DeleteImpactLines[0] != "入力データは残る" {
		t.Fatalf("expected delete impact lines mapping, got %#v", response)
	}
}

func TestTranslationJobManagementControllerGetJobDetailWrapsError(t *testing.T) {
	// GetJobDetail の失敗時に method 境界の wrap を証明する。
	controller := NewTranslationJobManagementController(fakeTranslationJobManagementUsecase{
		getDetailFunc: func(context.Context, usecase.TranslationJobManagementGetDetailRequest) (usecase.TranslationJobManagementJobDetail, error) {
			return usecase.TranslationJobManagementJobDetail{}, errors.New("db unavailable")
		},
	})

	_, err := controller.GetJobDetail(TranslationJobManagementGetDetailRequestDTO{JobID: 3})
	if err == nil {
		t.Fatal("expected error")
	}
	if !strings.Contains(err.Error(), "get translation job management detail") {
		t.Fatalf("expected wrapped method context, got %v", err)
	}
}

func TestTranslationJobManagementControllerDeleteJobMapsReasonCategoryAndNullability(t *testing.T) {
	controller := NewTranslationJobManagementController(fakeTranslationJobManagementUsecase{
		deleteFunc: func(context.Context, usecase.TranslationJobManagementDeleteRequest) (usecase.TranslationJobManagementActionResult, error) {
			return usecase.TranslationJobManagementActionResult{
				Message:        "実行中 job は削除できません",
				Tone:           "warning",
				ReasonCategory: "running_delete_blocked",
				Detail: &usecase.TranslationJobManagementJobDetail{
					TranslationJobManagementJobSummary: usecase.TranslationJobManagementJobSummary{JobID: 1},
				},
			}, nil
		},
	})

	response, err := controller.DeleteJob(TranslationJobManagementDeleteRequestDTO{JobID: 1})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if response.ReasonCategory != "running_delete_blocked" {
		t.Fatalf("expected reason category mapping, got %#v", response)
	}
	if response.Detail == nil {
		t.Fatalf("expected detail to stay non-nil, got %#v", response)
	}
	if response.DeletedJobID != nil {
		t.Fatalf("expected deletedJobId to stay nil, got %#v", response)
	}
}

func TestTranslationJobManagementControllerRequestStopMapsRequestAndResponse(t *testing.T) {
	// RequestStop の request/response DTO 写像を証明する。
	var captured usecase.TranslationJobManagementActionRequest
	controller := NewTranslationJobManagementController(fakeTranslationJobManagementUsecase{
		stopFunc: func(_ context.Context, req usecase.TranslationJobManagementActionRequest) (usecase.TranslationJobManagementActionResult, error) {
			captured = req
			return usecase.TranslationJobManagementActionResult{
				Message:        "停止要求を受け付けました",
				Tone:           "info",
				ReasonCategory: "stop_requested",
			}, nil
		},
	})

	response, err := controller.RequestStop(TranslationJobManagementActionRequestDTO{JobID: 22})
	if err != nil {
		t.Fatalf("expected request stop to succeed: %v", err)
	}
	if captured.JobID != 22 {
		t.Fatalf("expected forwarded job id, got %#v", captured)
	}
	if response.Message != "停止要求を受け付けました" || response.ReasonCategory != "stop_requested" {
		t.Fatalf("expected response mapping, got %#v", response)
	}
}

func TestTranslationJobManagementControllerRequestStopWrapsError(t *testing.T) {
	// RequestStop の失敗時に method 境界の wrap を証明する。
	controller := NewTranslationJobManagementController(fakeTranslationJobManagementUsecase{
		stopFunc: func(context.Context, usecase.TranslationJobManagementActionRequest) (usecase.TranslationJobManagementActionResult, error) {
			return usecase.TranslationJobManagementActionResult{}, errors.New("stop queue unavailable")
		},
	})

	_, err := controller.RequestStop(TranslationJobManagementActionRequestDTO{JobID: 22})
	if err == nil {
		t.Fatal("expected error")
	}
	if !strings.Contains(err.Error(), "request translation job stop") {
		t.Fatalf("expected wrapped method context, got %v", err)
	}
}

func TestTranslationJobManagementControllerResumeJobMapsRequestAndResponse(t *testing.T) {
	// ResumeJob の request/response DTO 写像を証明する。
	var captured usecase.TranslationJobManagementActionRequest
	controller := NewTranslationJobManagementController(fakeTranslationJobManagementUsecase{
		resumeFunc: func(_ context.Context, req usecase.TranslationJobManagementActionRequest) (usecase.TranslationJobManagementActionResult, error) {
			captured = req
			return usecase.TranslationJobManagementActionResult{
				Message: "再開しました",
				Tone:    "success",
			}, nil
		},
	})

	response, err := controller.ResumeJob(TranslationJobManagementActionRequestDTO{JobID: 99})
	if err != nil {
		t.Fatalf("expected resume job to succeed: %v", err)
	}
	if captured.JobID != 99 {
		t.Fatalf("expected forwarded job id, got %#v", captured)
	}
	if response.Message != "再開しました" || response.Tone != "success" {
		t.Fatalf("expected response mapping, got %#v", response)
	}
}

func TestTranslationJobManagementControllerResumeJobWrapsError(t *testing.T) {
	// ResumeJob の失敗時に method 境界の wrap を証明する。
	controller := NewTranslationJobManagementController(fakeTranslationJobManagementUsecase{
		resumeFunc: func(context.Context, usecase.TranslationJobManagementActionRequest) (usecase.TranslationJobManagementActionResult, error) {
			return usecase.TranslationJobManagementActionResult{}, errors.New("resume blocked")
		},
	})

	_, err := controller.ResumeJob(TranslationJobManagementActionRequestDTO{JobID: 99})
	if err == nil {
		t.Fatal("expected error")
	}
	if !strings.Contains(err.Error(), "resume translation job") {
		t.Fatalf("expected wrapped method context, got %v", err)
	}
}

func TestTranslationJobManagementControllerListIncompleteJobsWrapsError(t *testing.T) {
	controller := NewTranslationJobManagementController(fakeTranslationJobManagementUsecase{
		listFunc: func(context.Context) (usecase.TranslationJobManagementListResult, error) {
			return usecase.TranslationJobManagementListResult{}, errors.New("boom")
		},
	})

	_, err := controller.ListIncompleteJobs()
	if err == nil {
		t.Fatal("expected error")
	}
	if !strings.Contains(err.Error(), "list translation job management jobs") {
		t.Fatalf("expected wrapped error context, got %v", err)
	}
}

func TestTranslationJobManagementActionResponseDTODeletedJobIDTag(t *testing.T) {
	field, ok := reflect.TypeOf(TranslationJobManagementActionResponseDTO{}).FieldByName("DeletedJobID")
	if !ok {
		t.Fatal("expected DeletedJobID field")
	}
	if got := field.Tag.Get("json"); got != "deletedJobId,omitempty" {
		t.Fatalf("expected deletedJobId json tag, got %q", got)
	}
}
