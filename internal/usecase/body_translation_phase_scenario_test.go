package usecase

import (
	"context"
	"encoding/json"
	"errors"
	"strings"
	"testing"

	"aitranslationenginejp/internal/service"
)

type fakeBodyTranslationScenarioService struct {
	startReadiness service.BodyTranslationPhaseCommandReadModel
	readiness      service.BodyTranslationOutputReadinessReadModel
	err            error
	capturedJobID  int64
}

func (fake *fakeBodyTranslationScenarioService) ReadSummary(context.Context, int64) (service.BodyTranslationPhaseSummaryReadModel, error) {
	return service.BodyTranslationPhaseSummaryReadModel{}, nil
}

func (fake *fakeBodyTranslationScenarioService) StartPhase(_ context.Context, jobID int64) (service.BodyTranslationPhaseCommandReadModel, error) {
	fake.capturedJobID = jobID
	return fake.startReadiness, fake.err
}

func (fake *fakeBodyTranslationScenarioService) PausePhase(context.Context, int64, int64) (service.BodyTranslationPhaseCommandReadModel, error) {
	return service.BodyTranslationPhaseCommandReadModel{}, nil
}

func (fake *fakeBodyTranslationScenarioService) ResumePhase(context.Context, int64, int64) (service.BodyTranslationPhaseCommandReadModel, error) {
	return service.BodyTranslationPhaseCommandReadModel{}, nil
}

func (fake *fakeBodyTranslationScenarioService) RetryPhase(context.Context, int64, int64) (service.BodyTranslationPhaseCommandReadModel, error) {
	return service.BodyTranslationPhaseCommandReadModel{}, nil
}

func (fake *fakeBodyTranslationScenarioService) CancelPhase(context.Context, int64, int64) (service.BodyTranslationPhaseCommandReadModel, error) {
	return service.BodyTranslationPhaseCommandReadModel{}, nil
}

func (fake *fakeBodyTranslationScenarioService) ReadOutputReadiness(context.Context, int64) (service.BodyTranslationOutputReadinessReadModel, error) {
	return fake.readiness, fake.err
}

func TestSCN_BTP_001_StartCreatesBodyRunWhenPersonaPhaseCompletedAndInputsResolve(t *testing.T) {
	phaseRunID := int64(501)
	fake := &fakeBodyTranslationScenarioService{
		startReadiness: service.BodyTranslationPhaseCommandReadModel{
			JobID:        1001,
			CurrentPhase: "body_translation",
			PhaseState:   "running",
			PhaseRunID:   &phaseRunID,
			InputSummary: service.BodyTranslationPhaseInputSummaryReadModel{
				TargetCount:      2,
				DictionaryDigest: "sha256:dictionary",
				PersonaDigest:    "sha256:persona",
				MetadataDigest:   "sha256:metadata",
				PromptDigest:     "sha256:prompt",
			},
		},
	}
	phaseUsecase := NewBodyTranslationPhaseUsecase(fake)

	result, err := phaseUsecase.StartBodyTranslationPhase(context.Background(), StartBodyTranslationPhaseRequest{JobID: 1001})
	if err != nil {
		t.Fatalf("SCN-BTP-001 expected body phase start success: %v", err)
	}
	if fake.capturedJobID != 1001 {
		t.Fatalf("SCN-BTP-001 expected job id 1001, got %d", fake.capturedJobID)
	}
	if result.PhaseRunID == nil || *result.PhaseRunID != phaseRunID {
		t.Fatalf("SCN-BTP-001 expected body phase run id, got %#v", result.PhaseRunID)
	}
	if result.CurrentPhase != "body_translation" || result.PhaseState != "running" {
		t.Fatalf("SCN-BTP-001 expected running body phase, got %#v", result)
	}
}

func TestSCN_BTP_001_StartRejectsBlockedStartConditions(t *testing.T) {
	cases := []struct {
		name      string
		errorKind BodyTranslationPhaseErrorKind
	}{
		{name: "persona incomplete", errorKind: BodyTranslationPhaseErrorKindPersonaPhaseIncomplete},
		{name: "active phase exists", errorKind: BodyTranslationPhaseErrorKindActivePhaseExists},
		{name: "terminal job", errorKind: BodyTranslationPhaseErrorKindTerminalJob},
		{name: "input snapshot failed", errorKind: BodyTranslationPhaseErrorKindInputSnapshotFailed},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			fake := &fakeBodyTranslationScenarioService{
				startReadiness: service.BodyTranslationPhaseCommandReadModel{
					JobID:        1002,
					CurrentPhase: "body_translation",
					PhaseState:   "rejected",
					ErrorSummary: &service.BodyTranslationPhaseErrorSummaryReadModel{
						ErrorKind:  strings.ToUpper(string(tc.errorKind)),
						Reason:     "redacted rejection",
						Retryable:  false,
						IsRedacted: true,
					},
				},
				err: errors.New("body phase start rejected"),
			}
			phaseUsecase := NewBodyTranslationPhaseUsecase(fake)

			result, err := phaseUsecase.StartBodyTranslationPhase(context.Background(), StartBodyTranslationPhaseRequest{JobID: 1002})
			if err == nil {
				t.Fatal("SCN-BTP-001 expected rejected start")
			}
			if result.PhaseRunID != nil {
				t.Fatalf("SCN-BTP-001 expected no phase run on rejection, got %#v", result.PhaseRunID)
			}
			if result.ErrorSummary == nil || result.ErrorSummary.ErrorKind != tc.errorKind || !result.ErrorSummary.IsRedacted {
				t.Fatalf("SCN-BTP-001 expected redacted error kind %q, got %#v", tc.errorKind, result.ErrorSummary)
			}
		})
	}
}

func TestSCN_BTP_002_InputSnapshotCarriesDigestAndDictionaryRequestSummary(t *testing.T) {
	phaseRunID := int64(502)
	fake := &fakeBodyTranslationScenarioService{
		startReadiness: service.BodyTranslationPhaseCommandReadModel{
			JobID:        1003,
			CurrentPhase: "body_translation",
			PhaseState:   "running",
			PhaseRunID:   &phaseRunID,
			InputSummary: service.BodyTranslationPhaseInputSummaryReadModel{
				TargetCount:      3,
				SkippedReasons:   []string{"exact_dictionary_match"},
				DictionaryDigest: "sha256:dictionary",
				PersonaDigest:    "sha256:persona",
				MetadataDigest:   "sha256:metadata",
				PromptDigest:     "sha256:prompt",
			},
			RequestSummary: service.BodyTranslationPhaseRequestSummaryReadModel{
				ProviderTargetCount:              2,
				ExactDictionaryExclusionCount:    1,
				PartialDictionaryConstraintCount: 1,
			},
		},
	}
	phaseUsecase := NewBodyTranslationPhaseUsecase(fake)

	result, err := phaseUsecase.StartBodyTranslationPhase(context.Background(), StartBodyTranslationPhaseRequest{JobID: 1003})
	if err != nil {
		t.Fatalf("SCN-BTP-002 expected body phase start success: %v", err)
	}
	if result.InputSummary.TargetCount != 3 {
		t.Fatalf("SCN-BTP-002 expected target count 3, got %d", result.InputSummary.TargetCount)
	}
	if len(result.InputSummary.SkippedReasons) == 0 {
		t.Fatal("SCN-BTP-002 expected skipped reasons to be part of input snapshot")
	}
	assertNonEmptyBodyDigest(t, "dictionary", result.InputSummary.DictionaryDigest)
	assertNonEmptyBodyDigest(t, "persona", result.InputSummary.PersonaDigest)
	assertNonEmptyBodyDigest(t, "metadata", result.InputSummary.MetadataDigest)
	assertNonEmptyBodyDigest(t, "prompt", result.InputSummary.PromptDigest)
	if result.RequestSummary.ProviderTargetCount != 2 {
		t.Fatalf("SCN-BTP-002 expected exact dictionary hit to be excluded from provider request, got %#v", result.RequestSummary)
	}
	if result.RequestSummary.PartialDictionaryConstraintCount != 1 {
		t.Fatalf("SCN-BTP-002 expected partial dictionary hit to remain as fixed term constraint, got %#v", result.RequestSummary)
	}
}

func TestSCN_BTP_008_StartPublicResultExposesOnlySafePromptSummary(t *testing.T) {
	phaseRunID := int64(504)
	fake := &fakeBodyTranslationScenarioService{
		startReadiness: service.BodyTranslationPhaseCommandReadModel{
			JobID:        1005,
			CurrentPhase: "body_translation",
			PhaseState:   "running",
			PhaseRunID:   &phaseRunID,
			InputSummary: service.BodyTranslationPhaseInputSummaryReadModel{
				TargetCount:      2,
				DictionaryDigest: "sha256:dictionary-safe-summary",
				PersonaDigest:    "sha256:persona-safe-summary",
				MetadataDigest:   "sha256:metadata-safe-summary",
				PromptDigest:     "sha256:prompt-safe-summary",
			},
			RequestSummary: service.BodyTranslationPhaseRequestSummaryReadModel{
				ProviderTargetCount:              2,
				ExactDictionaryExclusionCount:    1,
				PartialDictionaryConstraintCount: 1,
			},
			Execution: service.BodyTranslationPhaseExecutionSummaryReadModel{
				CredentialRef:    "credential:body:redacted",
				CredentialState:  "configured",
				Provider:         "fake",
				Model:            "body-model",
				ExecutionMode:    "single_request",
				RequestUnitCount: 2,
				OutputCount:      0,
			},
		},
	}
	phaseUsecase := NewBodyTranslationPhaseUsecase(fake)

	// SCN-01: backend 開始結果は生成指示全文ではなく digest と件数だけを公開する。
	result, err := phaseUsecase.StartBodyTranslationPhase(context.Background(), StartBodyTranslationPhaseRequest{JobID: 1005})
	if err != nil {
		t.Fatalf("SCN-01 expected safe public start summary: %v", err)
	}
	if result.InputSummary.PromptDigest != "sha256:prompt-safe-summary" {
		t.Fatalf("SCN-01 expected prompt digest in public result, got %#v", result.InputSummary)
	}
	if result.RequestSummary.ProviderTargetCount != 2 || result.Execution.RequestUnitCount != 2 {
		t.Fatalf("SCN-01 expected provider counts in public result, got request=%#v execution=%#v", result.RequestSummary, result.Execution)
	}
	assertBodyTranslationScenarioJSONOmitsForbiddenValues(t, result)
}

func TestSCN_BTP_008_StartInvalidProviderResponseReturnsRedactedFailureKind(t *testing.T) {
	phaseRunID := int64(505)
	fake := &fakeBodyTranslationScenarioService{
		startReadiness: service.BodyTranslationPhaseCommandReadModel{
			JobID:        1006,
			CurrentPhase: "body_translation",
			PhaseState:   "recoverable_failed",
			PhaseRunID:   &phaseRunID,
			InputSummary: service.BodyTranslationPhaseInputSummaryReadModel{
				TargetCount:  1,
				PromptDigest: "sha256:invalid-provider-response-prompt",
			},
			Execution: service.BodyTranslationPhaseExecutionSummaryReadModel{
				Provider:         "fake",
				Model:            "body-model",
				ExecutionMode:    "single_request",
				RequestUnitCount: 1,
				OutputCount:      0,
			},
			ErrorSummary: &service.BodyTranslationPhaseErrorSummaryReadModel{
				ErrorKind:  " INVALID_PROVIDER_RESPONSE ",
				Reason:     "provider response failed body translation validation",
				Retryable:  true,
				IsRedacted: true,
			},
		},
		err: errors.New("provider response failed validation"),
	}
	phaseUsecase := NewBodyTranslationPhaseUsecase(fake)

	// SCN-01: invalid provider response は本文翻訳フェーズの失敗分類だけを公開する。
	result, err := phaseUsecase.StartBodyTranslationPhase(context.Background(), StartBodyTranslationPhaseRequest{JobID: 1006})
	if err == nil {
		t.Fatal("SCN-01 expected invalid provider response failure")
	}
	if result.ErrorSummary == nil {
		t.Fatal("SCN-01 expected redacted error summary")
	}
	if result.ErrorSummary.ErrorKind != BodyTranslationPhaseErrorKindInvalidProviderResponse {
		t.Fatalf("SCN-01 expected invalid provider response kind, got %#v", result.ErrorSummary)
	}
	if !result.ErrorSummary.IsRedacted || !result.ErrorSummary.Retryable {
		t.Fatalf("SCN-01 expected retryable redacted failure, got %#v", result.ErrorSummary)
	}
	assertBodyTranslationScenarioJSONOmitsForbiddenValues(t, result)
}

func TestSCN_BTP_008_ZeroTargetCompletesWithoutProviderAndEnablesOutputReadiness(t *testing.T) {
	phaseRunID := int64(503)
	fake := &fakeBodyTranslationScenarioService{
		startReadiness: service.BodyTranslationPhaseCommandReadModel{
			JobID:        1004,
			CurrentPhase: "body_translation",
			PhaseState:   "completed",
			PhaseRunID:   &phaseRunID,
			InputSummary: service.BodyTranslationPhaseInputSummaryReadModel{
				TargetCount:      0,
				DictionaryDigest: "sha256:dictionary",
				PersonaDigest:    "sha256:persona",
				MetadataDigest:   "sha256:metadata",
				PromptDigest:     "sha256:prompt",
			},
			Execution: service.BodyTranslationPhaseExecutionSummaryReadModel{
				RequestUnitCount: 0,
				OutputCount:      0,
			},
			OutputReadiness: service.BodyTranslationOutputReadinessReadModel{
				JobID:            1004,
				CurrentPhase:     "body_translation",
				PhaseState:       "completed",
				Ready:            true,
				StatusConsistent: true,
			},
		},
		readiness: service.BodyTranslationOutputReadinessReadModel{
			JobID:            1004,
			CurrentPhase:     "body_translation",
			PhaseState:       "completed",
			Ready:            true,
			StatusConsistent: true,
		},
	}
	phaseUsecase := NewBodyTranslationPhaseUsecase(fake)

	result, err := phaseUsecase.StartBodyTranslationPhase(context.Background(), StartBodyTranslationPhaseRequest{JobID: 1004})
	if err != nil {
		t.Fatalf("SCN-BTP-008 expected zero target completed result: %v", err)
	}
	if result.PhaseState != "completed" {
		t.Fatalf("SCN-BTP-008 expected completed phase state, got %q", result.PhaseState)
	}
	if result.Execution.RequestUnitCount != 0 {
		t.Fatalf("SCN-BTP-008 expected no provider request, got %d", result.Execution.RequestUnitCount)
	}

	readiness, err := phaseUsecase.GetBodyTranslationOutputReadiness(context.Background(), GetBodyTranslationOutputReadinessRequest{JobID: 1004})
	if err != nil {
		t.Fatalf("SCN-BTP-008 expected output readiness success: %v", err)
	}
	if !readiness.Ready {
		t.Fatalf("SCN-BTP-008 expected output readiness, got %#v", readiness)
	}
}

func assertNonEmptyBodyDigest(t *testing.T, name string, digest string) {
	t.Helper()
	if strings.TrimSpace(digest) == "" {
		t.Fatalf("expected non-empty %s digest", name)
	}
}

func assertBodyTranslationScenarioJSONOmitsForbiddenValues(t *testing.T, value any) {
	t.Helper()
	payload, err := json.Marshal(value)
	if err != nil {
		t.Fatalf("expected scenario result to marshal: %v", err)
	}
	serialized := string(payload)
	for _, forbidden := range []string{
		"sk-live-secret",
		"raw prompt with protected source text",
		"provider raw response body",
		"api_key",
		"authorization",
		"token",
		"raw_prompt",
		"provider_raw_request",
		"provider_raw_response",
		"Whiterun protected source sentence",
		"conversation context full text",
	} {
		if strings.Contains(serialized, forbidden) {
			t.Fatalf("expected public scenario result to omit %q, got %s", forbidden, serialized)
		}
	}
}
