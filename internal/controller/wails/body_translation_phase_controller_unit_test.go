package wails

import (
	"context"
	"encoding/json"
	"errors"
	"strings"
	"testing"
	"time"

	"aitranslationenginejp/internal/usecase"
)

type fakeBodyTranslationPhaseUsecase struct {
	getSummaryFunc func(context.Context, usecase.GetBodyTranslationPhaseSummaryRequest) (usecase.BodyTranslationPhaseFetchResult, error)
	startFunc      func(context.Context, usecase.StartBodyTranslationPhaseRequest) (usecase.BodyTranslationPhaseCommandResult, error)
	pauseFunc      func(context.Context, usecase.PauseBodyTranslationPhaseRequest) (usecase.BodyTranslationPhaseCommandResult, error)
	resumeFunc     func(context.Context, usecase.ResumeBodyTranslationPhaseRequest) (usecase.BodyTranslationPhaseCommandResult, error)
	retryFunc      func(context.Context, usecase.RetryBodyTranslationPhaseRequest) (usecase.BodyTranslationPhaseCommandResult, error)
	cancelFunc     func(context.Context, usecase.CancelBodyTranslationPhaseRequest) (usecase.BodyTranslationPhaseCommandResult, error)
}

func (fake fakeBodyTranslationPhaseUsecase) GetBodyTranslationPhaseSummary(ctx context.Context, request usecase.GetBodyTranslationPhaseSummaryRequest) (usecase.BodyTranslationPhaseFetchResult, error) {
	if fake.getSummaryFunc != nil {
		return fake.getSummaryFunc(ctx, request)
	}
	return usecase.BodyTranslationPhaseFetchResult{}, nil
}

func (fake fakeBodyTranslationPhaseUsecase) StartBodyTranslationPhase(ctx context.Context, request usecase.StartBodyTranslationPhaseRequest) (usecase.BodyTranslationPhaseCommandResult, error) {
	if fake.startFunc != nil {
		return fake.startFunc(ctx, request)
	}
	return usecase.BodyTranslationPhaseCommandResult{}, nil
}

func (fake fakeBodyTranslationPhaseUsecase) PauseBodyTranslationPhase(ctx context.Context, request usecase.PauseBodyTranslationPhaseRequest) (usecase.BodyTranslationPhaseCommandResult, error) {
	if fake.pauseFunc != nil {
		return fake.pauseFunc(ctx, request)
	}
	return usecase.BodyTranslationPhaseCommandResult{}, nil
}

func (fake fakeBodyTranslationPhaseUsecase) ResumeBodyTranslationPhase(ctx context.Context, request usecase.ResumeBodyTranslationPhaseRequest) (usecase.BodyTranslationPhaseCommandResult, error) {
	if fake.resumeFunc != nil {
		return fake.resumeFunc(ctx, request)
	}
	return usecase.BodyTranslationPhaseCommandResult{}, nil
}

func (fake fakeBodyTranslationPhaseUsecase) RetryBodyTranslationPhase(ctx context.Context, request usecase.RetryBodyTranslationPhaseRequest) (usecase.BodyTranslationPhaseCommandResult, error) {
	if fake.retryFunc != nil {
		return fake.retryFunc(ctx, request)
	}
	return usecase.BodyTranslationPhaseCommandResult{}, nil
}

func (fake fakeBodyTranslationPhaseUsecase) CancelBodyTranslationPhase(ctx context.Context, request usecase.CancelBodyTranslationPhaseRequest) (usecase.BodyTranslationPhaseCommandResult, error) {
	if fake.cancelFunc != nil {
		return fake.cancelFunc(ctx, request)
	}
	return usecase.BodyTranslationPhaseCommandResult{}, nil
}

func TestBodyTranslationPhaseControllerGetSummaryMapsPublicSeamAndRedactsSecrets(t *testing.T) {
	startedAt := time.Date(2026, 5, 2, 13, 0, 0, 0, time.FixedZone("JST", 9*60*60))
	phaseRunID := int64(901)
	controller := NewBodyTranslationPhaseController(fakeBodyTranslationPhaseUsecase{
		getSummaryFunc: func(context.Context, usecase.GetBodyTranslationPhaseSummaryRequest) (usecase.BodyTranslationPhaseFetchResult, error) {
			return usecase.BodyTranslationPhaseFetchResult{
				Summary: usecase.BodyTranslationPhaseSummaryResult{
					JobID:        501,
					CurrentPhase: "body_translation",
					PhaseState:   "running",
					PhaseRunID:   &phaseRunID,
					StartedAt:    &startedAt,
					Progress: usecase.BodyTranslationPhaseProgressSummary{
						Percent:         40,
						ProcessedCount:  4,
						TotalCount:      10,
						TargetCount:     10,
						TranslatedCount: 3,
						SkippedCount:    1,
						CurrentStep:     "provider_request",
					},
					InputSummary: usecase.BodyTranslationInputSummary{
						TargetCount:      10,
						DictionaryDigest: "sha256:dictionary",
						PersonaDigest:    "sha256:persona",
						MetadataDigest:   "sha256:metadata",
						PromptDigest:     "sha256:prompt",
					},
					Execution: &usecase.BodyTranslationExecutionSummary{
						CredentialRef:    "credential:body:test",
						Provider:         "fake",
						Model:            "body-model",
						ExecutionMode:    "batch",
						RequestUnitCount: 10,
						OutputCount:      3,
					},
					FieldResultSummary: &usecase.BodyTranslationFieldResultSummary{
						TranslatedCount: 3,
						FailedCount:     1,
						SkippedCount:    1,
						OutputCount:     3,
					},
					FieldResults: []usecase.BodyTranslationPhaseFieldResultItem{
						{
							Identity: usecase.BodyTranslationPhaseFieldIdentity{
								TranslationFieldID:      701,
								PhaseTranslationFieldID: 801,
								RecordType:              "NPC_",
								FieldType:               "FULL",
								FormID:                  "000001",
								EditorID:                "FixtureNPC",
								FieldLabel:              "FixtureNPC / FULL",
							},
							FieldID:                    701,
							FieldLabel:                 "FixtureNPC / FULL",
							SourceExcerpt:              "Hello",
							TranslatedText:             "こんにちは",
							OutputStatus:               "ready",
							ProtectionValidationResult: "passed",
							RetryCount:                 1,
						},
					},
					ErrorSummary: &usecase.BodyTranslationPhaseErrorSummary{
						ErrorKind:  " SECRET_REDACTED ",
						Reason:     "redacted public summary",
						Retryable:  true,
						IsRedacted: true,
					},
					OutputReadiness: usecase.BodyTranslationOutputReadinessSummary{
						CompletedFieldCount: 3,
						StatusConsistent:    true,
					},
				},
				Projection: usecase.BodyTranslationPhaseProjectionResult{
					PhaseLifecycle:       "running",
					JobLifecycle:         "running",
					ErrorKind:            "recoverable",
					AISettingsConfigured: false,
					TargetCount:          10,
				},
			}, nil
		},
	})

	response, err := controller.GetBodyTranslationPhaseSummary(GetBodyTranslationPhaseSummaryRequestDTO{JobID: 501})
	if err != nil {
		t.Fatalf("expected success, got %v", err)
	}
	if response.StartedAt == nil || *response.StartedAt != "2026-05-02T04:00:00Z" {
		t.Fatalf("expected UTC startedAt, got %#v", response.StartedAt)
	}
	if response.ErrorSummary == nil || response.ErrorSummary.ErrorKind != "secret_redacted" || !response.ErrorSummary.IsRedacted {
		t.Fatalf("expected normalized redacted error summary, got %#v", response.ErrorSummary)
	}
	if response.Execution.CredentialRef != "credential:body:test" || response.Execution.Provider != "fake" {
		t.Fatalf("expected redacted execution references, got %#v", response.Execution)
	}
	if len(response.FieldResults) != 1 || response.FieldResults[0].Identity.TranslationFieldID != 701 {
		t.Fatalf("expected field result list mapping, got %#v", response.FieldResults)
	}
	assertBodyTranslationDTOHasNoForbiddenSecretFields(t, response)
}

func TestBodyTranslationPhaseControllerGetSummaryReturnsEmptySkippedReasonsArray(t *testing.T) {
	controller := NewBodyTranslationPhaseController(fakeBodyTranslationPhaseUsecase{
		getSummaryFunc: func(context.Context, usecase.GetBodyTranslationPhaseSummaryRequest) (usecase.BodyTranslationPhaseFetchResult, error) {
			return usecase.BodyTranslationPhaseFetchResult{}, nil
		},
	})

	response, err := controller.GetBodyTranslationPhaseSummary(GetBodyTranslationPhaseSummaryRequestDTO{JobID: 501})
	if err != nil {
		t.Fatalf("expected success, got %v", err)
	}
	payload, err := json.Marshal(response)
	if err != nil {
		t.Fatalf("expected marshal success, got %v", err)
	}
	if strings.Contains(string(payload), `"skippedReasons":null`) {
		t.Fatalf("expected skippedReasons array, got %s", payload)
	}
	if !strings.Contains(string(payload), `"skippedReasons":[]`) {
		t.Fatalf("expected empty skippedReasons array, got %s", payload)
	}
}

func TestBodyTranslationPhaseControllerStartCommandSeamForwardsJobID(t *testing.T) {
	controller := NewBodyTranslationPhaseController(fakeBodyTranslationPhaseUsecase{
		startFunc: func(_ context.Context, request usecase.StartBodyTranslationPhaseRequest) (usecase.BodyTranslationPhaseCommandResult, error) {
			if request.JobID != 777 {
				t.Fatalf("expected start job id 777, got %#v", request)
			}
			return bodyTranslationCommandFixture(), nil
		},
	})

	response, err := controller.StartBodyTranslationPhase(StartBodyTranslationPhaseRequestDTO{JobID: 777})
	if err != nil {
		t.Fatalf("expected start success, got %v", err)
	}
	if response.PhaseRunID == nil || *response.PhaseRunID != 41 {
		t.Fatalf("expected phase run id 41, got %#v", response.PhaseRunID)
	}
}

func TestBodyTranslationPhaseControllerStartReturnsOnlySafeSummaryForInvalidProviderResponse(t *testing.T) {
	controller := NewBodyTranslationPhaseController(fakeBodyTranslationPhaseUsecase{
		startFunc: func(context.Context, usecase.StartBodyTranslationPhaseRequest) (usecase.BodyTranslationPhaseCommandResult, error) {
			phaseRunID := int64(42)
			return usecase.BodyTranslationPhaseCommandResult{
				JobID:        778,
				CurrentPhase: "body_translation",
				PhaseState:   "recoverable_failed",
				PhaseRunID:   &phaseRunID,
				InputSummary: usecase.BodyTranslationPhaseInputSummary{
					TargetCount:  1,
					PromptDigest: "sha256:public-prompt-digest",
				},
				RequestSummary: usecase.BodyTranslationPhaseRequestSummary{
					ProviderTargetCount: 1,
				},
				Execution: &usecase.BodyTranslationPhaseExecutionSummary{
					CredentialRef:    "credential:body:redacted",
					Provider:         "fake",
					Model:            "body-model",
					ExecutionMode:    "single_request",
					RequestUnitCount: 1,
					OutputCount:      0,
				},
				ErrorSummary: &usecase.BodyTranslationPhaseErrorSummary{
					ErrorKind:  usecase.BodyTranslationPhaseErrorKindInvalidProviderResponse,
					Reason:     "provider response failed body translation validation",
					Retryable:  true,
					IsRedacted: true,
				},
			}, errors.New("provider raw response body included invalid field correlation")
		},
	})

	// SCN-01: Wails 公開接点は invalid provider response を分類だけで返す。
	response, err := controller.StartBodyTranslationPhase(StartBodyTranslationPhaseRequestDTO{JobID: 778})
	if err != nil {
		t.Fatalf("expected structured invalid response payload without raw provider error: %v", err)
	}
	if response.ErrorSummary == nil || response.ErrorSummary.ErrorKind != "invalid_provider_response" {
		t.Fatalf("expected invalid provider response kind, got %#v", response.ErrorSummary)
	}
	if !response.ErrorSummary.IsRedacted || !response.ErrorSummary.Retryable {
		t.Fatalf("expected retryable redacted summary, got %#v", response.ErrorSummary)
	}
	if response.InputSummary.PromptDigest != "sha256:public-prompt-digest" {
		t.Fatalf("expected prompt digest without raw prompt, got %#v", response.InputSummary)
	}
	if response.RequestSummary.ProviderTargetCount != 1 || response.Execution.RequestUnitCount != 1 {
		t.Fatalf("expected public counts, got request=%#v execution=%#v", response.RequestSummary, response.Execution)
	}
	assertBodyTranslationCommandDTOHasNoForbiddenRawPayload(t, response)
}

func TestBodyTranslationPhaseControllerPauseCommandSeamForwardsPhaseRunID(t *testing.T) {
	controller := NewBodyTranslationPhaseController(fakeBodyTranslationPhaseUsecase{
		pauseFunc: func(_ context.Context, request usecase.PauseBodyTranslationPhaseRequest) (usecase.BodyTranslationPhaseCommandResult, error) {
			if request.JobID != 777 || request.PhaseRunID != 41 {
				t.Fatalf("expected pause request, got %#v", request)
			}
			return bodyTranslationCommandFixture(), nil
		},
	})

	response, err := controller.PauseBodyTranslationPhase(PauseBodyTranslationPhaseRequestDTO{JobID: 777, PhaseRunID: 41})
	if err != nil {
		t.Fatalf("expected pause success, got %v", err)
	}
	if response.PhaseRunID == nil || *response.PhaseRunID != 41 {
		t.Fatalf("expected phase run id 41, got %#v", response.PhaseRunID)
	}
}

func TestBodyTranslationPhaseControllerResumeCommandSeamForwardsPhaseRunID(t *testing.T) {
	controller := NewBodyTranslationPhaseController(fakeBodyTranslationPhaseUsecase{
		resumeFunc: func(_ context.Context, request usecase.ResumeBodyTranslationPhaseRequest) (usecase.BodyTranslationPhaseCommandResult, error) {
			if request.JobID != 777 || request.PhaseRunID != 41 {
				t.Fatalf("expected resume request, got %#v", request)
			}
			return bodyTranslationCommandFixture(), nil
		},
	})

	response, err := controller.ResumeBodyTranslationPhase(ResumeBodyTranslationPhaseRequestDTO{JobID: 777, PhaseRunID: 41})
	if err != nil {
		t.Fatalf("expected resume success, got %v", err)
	}
	if response.PhaseRunID == nil || *response.PhaseRunID != 41 {
		t.Fatalf("expected phase run id 41, got %#v", response.PhaseRunID)
	}
}

func TestBodyTranslationPhaseControllerRetryCommandSeamForwardsPhaseRunID(t *testing.T) {
	controller := NewBodyTranslationPhaseController(fakeBodyTranslationPhaseUsecase{
		retryFunc: func(_ context.Context, request usecase.RetryBodyTranslationPhaseRequest) (usecase.BodyTranslationPhaseCommandResult, error) {
			if request.JobID != 777 || request.PhaseRunID != 41 {
				t.Fatalf("expected retry request, got %#v", request)
			}
			return bodyTranslationCommandFixture(), nil
		},
	})

	response, err := controller.RetryBodyTranslationPhase(RetryBodyTranslationPhaseRequestDTO{JobID: 777, PhaseRunID: 41})
	if err != nil {
		t.Fatalf("expected retry success, got %v", err)
	}
	if response.PhaseRunID == nil || *response.PhaseRunID != 41 {
		t.Fatalf("expected phase run id 41, got %#v", response.PhaseRunID)
	}
}

func TestBodyTranslationPhaseControllerCancelCommandSeamBlocksOutputReadiness(t *testing.T) {
	controller := NewBodyTranslationPhaseController(fakeBodyTranslationPhaseUsecase{
		cancelFunc: func(_ context.Context, request usecase.CancelBodyTranslationPhaseRequest) (usecase.BodyTranslationPhaseCommandResult, error) {
			if request.JobID != 777 || request.PhaseRunID != 41 {
				t.Fatalf("expected cancel request, got %#v", request)
			}
			result := bodyTranslationCommandFixture()
			result.PhaseState = "canceled"
			result.OutputReadiness.StatusConsistent = false
			return result, nil
		},
	})

	response, err := controller.CancelBodyTranslationPhase(CancelBodyTranslationPhaseRequestDTO{JobID: 777, PhaseRunID: 41})
	if err != nil {
		t.Fatalf("expected cancel success, got %v", err)
	}
	if response.OutputReadiness.StatusConsistent {
		t.Fatalf("expected canceled command to mark status inconsistent, got %#v", response.OutputReadiness)
	}
}

func TestBodyTranslationPhaseControllerCommandBusinessRejectionReturnsPayloadWithoutError(t *testing.T) {
	controller := NewBodyTranslationPhaseController(fakeBodyTranslationPhaseUsecase{
		cancelFunc: func(context.Context, usecase.CancelBodyTranslationPhaseRequest) (usecase.BodyTranslationPhaseCommandResult, error) {
			result := bodyTranslationCommandFixture()
			result.ErrorSummary = &usecase.BodyTranslationPhaseErrorSummary{
				ErrorKind:  usecase.BodyTranslationPhaseErrorKindOutputReadinessBlocked,
				Reason:     "body translation phase is not cancelable",
				Retryable:  false,
				IsRedacted: true,
			}
			return result, errors.New("body translation phase state transition rejected")
		},
	})

	response, err := controller.CancelBodyTranslationPhase(CancelBodyTranslationPhaseRequestDTO{JobID: 777, PhaseRunID: 41})
	if err != nil {
		t.Fatalf("expected business rejection payload without Wails error, got %v", err)
	}
	if response.ErrorSummary == nil || response.ErrorSummary.ErrorKind != "output_readiness_blocked" {
		t.Fatalf("expected structured rejection payload, got %#v", response)
	}
}

// TestBodyTranslationPhaseControllerOutputReadinessEndpointRemoved は廃止済み endpoint の確認テスト。
// GetBodyTranslationOutputReadiness endpoint は INT-body-summary-merge で廃止された。
// 出力確認可否は段階要約取得の outputReadiness フィールドから事実状態として取得する。

func bodyTranslationCommandFixture() usecase.BodyTranslationPhaseCommandResult {
	phaseRunID := int64(41)
	return usecase.BodyTranslationPhaseCommandResult{
		JobID:        777,
		CurrentPhase: "body_translation",
		PhaseState:   "running",
		PhaseRunID:   &phaseRunID,
		Progress: usecase.BodyTranslationPhaseProgressSummary{
			Percent:         50,
			ProcessedCount:  5,
			TotalCount:      10,
			TargetCount:     10,
			TranslatedCount: 4,
			SkippedCount:    1,
			CurrentStep:     "running",
		},
		InputSnapshotDigest: "sha256:input-snapshot",
		Retryable:           true,
		OutputReadiness: usecase.BodyTranslationOutputReadinessSummary{
			CompletedFieldCount: 4,
			StatusConsistent:    true,
		},
		ErrorSummary: &usecase.BodyTranslationPhaseErrorSummary{
			ErrorKind:  usecase.BodyTranslationPhaseErrorKindProviderFailure,
			Reason:     "redacted provider failure",
			Retryable:  true,
			IsRedacted: true,
		},
	}
}

func assertBodyTranslationDTOHasNoForbiddenSecretFields(t *testing.T, response BodyTranslationPhaseSummaryResponseDTO) {
	t.Helper()
	payload, err := json.Marshal(response)
	if err != nil {
		t.Fatalf("expected response to marshal, got %v", err)
	}
	serialized := string(payload)
	for _, forbidden := range []string{
		"apiKey",
		"authorization",
		"providerRawRequest",
		"providerRawResponse",
		"rawPrompt",
		"decrypted",
		"sk-live-secret",
	} {
		if strings.Contains(serialized, forbidden) {
			t.Fatalf("expected body translation DTO to omit %q, got %s", forbidden, serialized)
		}
	}
}

func assertBodyTranslationCommandDTOHasNoForbiddenRawPayload(t *testing.T, response BodyTranslationPhaseCommandResponseDTO) {
	t.Helper()
	payload, err := json.Marshal(response)
	if err != nil {
		t.Fatalf("expected response to marshal, got %v", err)
	}
	serialized := string(payload)
	for _, forbidden := range []string{
		"apiKey",
		"authorization",
		"providerRawRequest",
		"providerRawResponse",
		"rawPrompt",
		"provider raw response body",
		"raw prompt with protected source text",
		"sk-live-secret",
		"Whiterun protected source sentence",
		"conversation context full text",
	} {
		if strings.Contains(serialized, forbidden) {
			t.Fatalf("expected body translation command DTO to omit %q, got %s", forbidden, serialized)
		}
	}
}
