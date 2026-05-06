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
	return fake.listFunc(ctx)
}
func (fake fakeTranslationJobManagementUsecase) GetJobDetail(ctx context.Context, req usecase.TranslationJobManagementGetDetailRequest) (usecase.TranslationJobManagementJobDetail, error) {
	return fake.getDetailFunc(ctx, req)
}
func (fake fakeTranslationJobManagementUsecase) DeleteJob(ctx context.Context, req usecase.TranslationJobManagementDeleteRequest) (usecase.TranslationJobManagementActionResult, error) {
	return fake.deleteFunc(ctx, req)
}
func (fake fakeTranslationJobManagementUsecase) RequestStop(ctx context.Context, req usecase.TranslationJobManagementActionRequest) (usecase.TranslationJobManagementActionResult, error) {
	return fake.stopFunc(ctx, req)
}
func (fake fakeTranslationJobManagementUsecase) ResumeJob(ctx context.Context, req usecase.TranslationJobManagementActionRequest) (usecase.TranslationJobManagementActionResult, error) {
	return fake.resumeFunc(ctx, req)
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
