package service

import (
	"context"
	"errors"
	"testing"
)

func TestTermTranslationProviderAdapterReturnsRequestShapeIdentifierAndDigestAsInternalIdentity(t *testing.T) {
	// 単語翻訳の要求形状識別子と digest は生成規則版ではなく内部同一性情報として扱う。
	adapter := NewTermTranslationProviderAdapter(stubTermTranslationProviderClient{
		response: stubTermTranslationClientResponse{
			Items: []stubTermTranslationClientItem{{
				SourceTerm:     "Dragonborn",
				TranslatedTerm: "ドラゴンボーン",
			}},
			ExecutionMode: TermTranslationExecutionModeSingleRequest,
		},
	})

	result, err := adapter.TranslateTerm(context.Background(), TermTranslationProviderRequest{
		Provider:       TermTranslationProviderGemini,
		Model:          "gemini-2.5-pro",
		SourceTerm:     "Dragonborn",
		SourceLanguage: "English",
		TargetLanguage: "Japanese",
	})
	if err != nil {
		t.Fatalf("expected success, got %v", err)
	}
	if result.AuditSummary.RequestShapeID == nil || *result.AuditSummary.RequestShapeID != TermTranslationRequestShapeV1 {
		t.Fatalf("expected request shape identifier, got %#v", result.AuditSummary.RequestShapeID)
	}
	if result.AuditSummary.PromptDigest == nil || *result.AuditSummary.PromptDigest == "" {
		t.Fatalf("expected prompt digest identity, got %#v", result.AuditSummary.PromptDigest)
	}
}

func TestTermTranslationProviderAdapterReturnsInvalidResponseWhenProviderReturnsMissingItem(t *testing.T) {
	// 応答欠落は対象語単位の invalid response として扱う。
	adapter := NewTermTranslationProviderAdapter(stubTermTranslationProviderClient{
		response: stubTermTranslationClientResponse{
			Items:         []stubTermTranslationClientItem{},
			ExecutionMode: TermTranslationExecutionModeSingleRequest,
		},
	})

	_, err := adapter.TranslateTerm(context.Background(), TermTranslationProviderRequest{
		Provider:       TermTranslationProviderGemini,
		Model:          "gemini-2.5-pro",
		SourceTerm:     "Dragonborn",
		SourceLanguage: "English",
		TargetLanguage: "Japanese",
	})
	assertTermInvalidProviderResponse(t, err)
}

func TestTermTranslationProviderAdapterReturnsInvalidResponseWhenProviderReturnsExtraItem(t *testing.T) {
	// 余分な応答は対象語単位の invalid response として扱う。
	adapter := NewTermTranslationProviderAdapter(stubTermTranslationProviderClient{
		response: stubTermTranslationClientResponse{
			Items: []stubTermTranslationClientItem{
				{SourceTerm: "Dragonborn", TranslatedTerm: "ドラゴンボーン"},
				{SourceTerm: "Lydia", TranslatedTerm: "リディア"},
			},
			ExecutionMode: TermTranslationExecutionModeSingleRequest,
		},
	})

	_, err := adapter.TranslateTerm(context.Background(), TermTranslationProviderRequest{
		Provider:       TermTranslationProviderGemini,
		Model:          "gemini-2.5-pro",
		SourceTerm:     "Dragonborn",
		SourceLanguage: "English",
		TargetLanguage: "Japanese",
	})
	assertTermInvalidProviderResponse(t, err)
}

func TestTermTranslationProviderAdapterReturnsInvalidResponseWhenSourceTermMismatches(t *testing.T) {
	// 対象語不一致は対象語単位の invalid response として扱う。
	adapter := NewTermTranslationProviderAdapter(stubTermTranslationProviderClient{
		response: stubTermTranslationClientResponse{
			Items: []stubTermTranslationClientItem{{
				SourceTerm:     "Lydia",
				TranslatedTerm: "リディア",
			}},
			ExecutionMode: TermTranslationExecutionModeSingleRequest,
		},
	})

	_, err := adapter.TranslateTerm(context.Background(), TermTranslationProviderRequest{
		Provider:       TermTranslationProviderGemini,
		Model:          "gemini-2.5-pro",
		SourceTerm:     "Dragonborn",
		SourceLanguage: "English",
		TargetLanguage: "Japanese",
	})
	assertTermInvalidProviderResponse(t, err)
}

func TestTermTranslationProviderAdapterReturnsInvalidResponseWhenTranslatedTermIsEmpty(t *testing.T) {
	// 空訳語は対象語単位の invalid response として扱う。
	adapter := NewTermTranslationProviderAdapter(stubTermTranslationProviderClient{
		response: stubTermTranslationClientResponse{
			Items: []stubTermTranslationClientItem{{
				SourceTerm:     "Dragonborn",
				TranslatedTerm: " ",
			}},
			ExecutionMode: TermTranslationExecutionModeSingleRequest,
		},
	})

	_, err := adapter.TranslateTerm(context.Background(), TermTranslationProviderRequest{
		Provider:       TermTranslationProviderGemini,
		Model:          "gemini-2.5-pro",
		SourceTerm:     "Dragonborn",
		SourceLanguage: "English",
		TargetLanguage: "Japanese",
	})
	assertTermInvalidProviderResponse(t, err)
}

func assertTermInvalidProviderResponse(t *testing.T, err error) {
	t.Helper()
	if err == nil {
		t.Fatal("expected invalid provider response error")
	}
	var providerErr *TermTranslationProviderError
	if !errors.As(err, &providerErr) {
		t.Fatalf("expected term translation provider error, got %T", err)
	}
	if providerErr.Failure.Kind != TermTranslationProviderErrorKindInvalidProviderResponse || !providerErr.Failure.Retryable {
		t.Fatalf("unexpected provider failure: %#v", providerErr.Failure)
	}
}

type stubTermTranslationProviderClient struct {
	response stubTermTranslationClientResponse
	err      error
}

func (stub stubTermTranslationProviderClient) TranslateTerm(
	_ context.Context,
	_ string,
	_ string,
	_ string,
	_ string,
) (stubTermTranslationClientResponse, error) {
	if stub.err != nil {
		return stubTermTranslationClientResponse{}, stub.err
	}
	return stub.response, nil
}

type stubTermTranslationClientResponse struct {
	Items         []stubTermTranslationClientItem
	ExecutionMode string
}

type stubTermTranslationClientItem struct {
	SourceTerm     string
	TranslatedTerm string
}
