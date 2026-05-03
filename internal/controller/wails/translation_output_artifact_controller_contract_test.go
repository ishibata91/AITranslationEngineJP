package wails

import (
	"context"
	"encoding/json"
	"strings"
	"testing"

	"aitranslationenginejp/internal/usecase"
)

type fakeTranslationOutputArtifactUsecase struct {
	getReviewFunc      func(context.Context, usecase.GetTranslationOutputReviewRequest) (usecase.TranslationOutputReviewResult, error)
	getDiffPreviewFunc func(context.Context, usecase.GetTranslationOutputDiffPreviewRequest) (usecase.TranslationOutputDiffPreviewResult, error)
	generateFunc       func(context.Context, usecase.GenerateXTranslatorOutputArtifactRequest) (usecase.XTranslatorOutputArtifactCommandResult, error)
	regenerateFunc     func(context.Context, usecase.RegenerateXTranslatorOutputArtifactRequest) (usecase.XTranslatorOutputArtifactCommandResult, error)
}

func (fake fakeTranslationOutputArtifactUsecase) GetTranslationOutputReview(
	ctx context.Context,
	request usecase.GetTranslationOutputReviewRequest,
) (usecase.TranslationOutputReviewResult, error) {
	if fake.getReviewFunc != nil {
		return fake.getReviewFunc(ctx, request)
	}
	return usecase.TranslationOutputReviewResult{}, nil
}

func (fake fakeTranslationOutputArtifactUsecase) GetTranslationOutputDiffPreview(
	ctx context.Context,
	request usecase.GetTranslationOutputDiffPreviewRequest,
) (usecase.TranslationOutputDiffPreviewResult, error) {
	if fake.getDiffPreviewFunc != nil {
		return fake.getDiffPreviewFunc(ctx, request)
	}
	return usecase.TranslationOutputDiffPreviewResult{}, nil
}

func (fake fakeTranslationOutputArtifactUsecase) GenerateXTranslatorOutputArtifact(
	ctx context.Context,
	request usecase.GenerateXTranslatorOutputArtifactRequest,
) (usecase.XTranslatorOutputArtifactCommandResult, error) {
	if fake.generateFunc != nil {
		return fake.generateFunc(ctx, request)
	}
	return usecase.XTranslatorOutputArtifactCommandResult{}, nil
}

func (fake fakeTranslationOutputArtifactUsecase) RegenerateXTranslatorOutputArtifact(
	ctx context.Context,
	request usecase.RegenerateXTranslatorOutputArtifactRequest,
) (usecase.XTranslatorOutputArtifactCommandResult, error) {
	if fake.regenerateFunc != nil {
		return fake.regenerateFunc(ctx, request)
	}
	return usecase.XTranslatorOutputArtifactCommandResult{}, nil
}

func TestSCN_TOA_001_GetTranslationOutputReviewControllerSeamMapsSummary(t *testing.T) {
	controller := NewTranslationOutputArtifactController(fakeTranslationOutputArtifactUsecase{
		getReviewFunc: func(_ context.Context, request usecase.GetTranslationOutputReviewRequest) (usecase.TranslationOutputReviewResult, error) {
			if request.SelectedJobID == nil || *request.SelectedJobID != 101 {
				t.Fatalf("expected selected job id 101, got %#v", request)
			}
			return usecase.TranslationOutputReviewResult{
				CompletedJobs: []usecase.TranslationOutputCompletedJobSummary{
					{JobID: 101, JobStatus: "completed", ArtifactStatus: "not_generated", OutputReady: true},
				},
				SelectedJob: usecase.TranslationOutputSelectedJobSummary{
					JobID:           101,
					JobStatus:       "completed",
					BodyPhaseStatus: "completed",
					OutputReady:     true,
				},
				OutputReadiness: usecase.TranslationOutputReadinessSummary{Ready: true},
			}, nil
		},
	})

	selectedJobID := int64(101)
	response, err := controller.GetTranslationOutputReview(GetTranslationOutputReviewRequestDTO{
		SelectedJobID: &selectedJobID,
	})
	if err != nil {
		t.Fatalf("expected review success, got %v", err)
	}
	if len(response.CompletedJobs) != 1 || response.CompletedJobs[0].JobID != 101 {
		t.Fatalf("expected completed job response, got %#v", response)
	}
	if !response.OutputReadiness.Ready {
		t.Fatalf("expected readiness mapping, got %#v", response.OutputReadiness)
	}
	assertTranslationOutputArtifactDTOHasNoForbiddenSecretFields(t, response)
}

func TestSCN_TOA_003_GetTranslationOutputDiffPreviewControllerSeamMapsRows(t *testing.T) {
	controller := NewTranslationOutputArtifactController(fakeTranslationOutputArtifactUsecase{
		getDiffPreviewFunc: func(_ context.Context, request usecase.GetTranslationOutputDiffPreviewRequest) (usecase.TranslationOutputDiffPreviewResult, error) {
			if request.JobID != 303 || request.ArtifactID != 404 {
				t.Fatalf("expected diff preview request, got %#v", request)
			}
			return usecase.TranslationOutputDiffPreviewResult{
				JobID:      303,
				ArtifactID: 404,
				Rows: []usecase.TranslationOutputDiffPreviewRow{
					{
						FieldID:              501,
						RowDigest:            "sha256:row",
						EDID:                 "CachedNPC",
						REC:                  "NPC_",
						FIELD:                "FULL",
						FORMID:               "00000001",
						SourceExcerpt:        "Hello",
						DestExcerpt:          "こんにちは",
						XTranslatorStatus:    1,
						InternalOutputStatus: "cached",
					},
				},
			}, nil
		},
	})

	response, err := controller.GetTranslationOutputDiffPreview(GetTranslationOutputDiffPreviewRequestDTO{
		JobID:      303,
		ArtifactID: 404,
	})
	if err != nil {
		t.Fatalf("expected diff preview success, got %v", err)
	}
	if len(response.Rows) != 1 || response.Rows[0].XTranslatorStatus != 1 {
		t.Fatalf("expected cached row Status=1 mapping, got %#v", response.Rows)
	}
	assertTranslationOutputArtifactDTOHasNoForbiddenSecretFields(t, response)
}

func TestSCN_TOA_002_006_010_GenerateAndRegenerateControllerCommandsMapContract(t *testing.T) {
	controller := NewTranslationOutputArtifactController(fakeTranslationOutputArtifactUsecase{
		generateFunc: func(_ context.Context, request usecase.GenerateXTranslatorOutputArtifactRequest) (usecase.XTranslatorOutputArtifactCommandResult, error) {
			if request.JobID != 202 || request.TargetGame != "skyrim_se" || request.OutputPath == "" {
				t.Fatalf("expected generate request, got %#v", request)
			}
			return usecase.XTranslatorOutputArtifactCommandResult{
				JobID:          202,
				ArtifactStatus: "rejected",
				TargetGame:     "skyrim_se",
				ErrorSummary: &usecase.TranslationOutputArtifactErrorSummary{
					ErrorKind:  usecase.TranslationOutputArtifactErrorKindStatusMismatch,
					Reason:     "status mismatch",
					Retryable:  false,
					IsRedacted: true,
				},
			}, nil
		},
		regenerateFunc: func(_ context.Context, request usecase.RegenerateXTranslatorOutputArtifactRequest) (usecase.XTranslatorOutputArtifactCommandResult, error) {
			if request.JobID != 606 || request.ArtifactID != 707 {
				t.Fatalf("expected regenerate request, got %#v", request)
			}
			return usecase.XTranslatorOutputArtifactCommandResult{
				JobID:          606,
				ArtifactID:     707,
				ArtifactStatus: "success",
				RowCount:       2,
				FilePath:       "/tmp/output/skyrim.xml",
				TargetGame:     "skyrim_se",
				OperationSummary: usecase.TranslationOutputOperationSummary{
					OperationKind:       "regenerate",
					ReplacedArtifactID:  707,
					AffectedFieldIDs:    []int64{801, 802},
					DuplicateRowCreated: false,
				},
			}, nil
		},
	})

	generated, err := controller.GenerateXTranslatorOutputArtifact(GenerateXTranslatorOutputArtifactRequestDTO{
		JobID:      202,
		TargetGame: "skyrim_se",
		OutputPath: "/tmp/output/skyrim.xml",
	})
	if err != nil {
		t.Fatalf("expected generate response, got %v", err)
	}
	if generated.ErrorSummary == nil || generated.ErrorSummary.ErrorKind != "status_mismatch" {
		t.Fatalf("expected redacted status mismatch, got %#v", generated.ErrorSummary)
	}

	regenerated, err := controller.RegenerateXTranslatorOutputArtifact(RegenerateXTranslatorOutputArtifactRequestDTO{
		JobID:      606,
		ArtifactID: 707,
		TargetGame: "skyrim_se",
		OutputPath: "/tmp/output/skyrim.xml",
	})
	if err != nil {
		t.Fatalf("expected regenerate response, got %v", err)
	}
	if regenerated.OperationSummary.DuplicateRowCreated {
		t.Fatalf("expected no duplicate row on regenerate, got %#v", regenerated.OperationSummary)
	}
	assertTranslationOutputArtifactDTOHasNoForbiddenSecretFields(t, generated)
	assertTranslationOutputArtifactDTOHasNoForbiddenSecretFields(t, regenerated)
}

func assertTranslationOutputArtifactDTOHasNoForbiddenSecretFields(t *testing.T, response any) {
	t.Helper()
	encoded, err := json.Marshal(response)
	if err != nil {
		t.Fatalf("expected DTO to marshal: %v", err)
	}
	serialized := string(encoded)
	for _, forbidden := range []string{
		"sk-live-secret",
		"apiKey",
		"token",
		"authorizationHeader",
		"providerRawRequest",
		"providerRawResponse",
		"decrypted",
		"fullSourceText",
		"fullDestText",
	} {
		if strings.Contains(serialized, forbidden) {
			t.Fatalf("expected DTO not to expose %q: %s", forbidden, serialized)
		}
	}
}
