package wails

import (
	"context"
	"testing"

	"aitranslationenginejp/internal/usecase"
)

type processingTargetReadModelUsecaseSpy struct {
	request usecase.ProcessingTargetListRequest
}

func (spy *processingTargetReadModelUsecaseSpy) ListProcessingTargets(
	_ context.Context,
	request usecase.ProcessingTargetListRequest,
) (usecase.ProcessingTargetListResult, error) {
	spy.request = request
	return usecase.ProcessingTargetListResult{
		Items: []usecase.ProcessingTargetListItem{
			{
				ID:     "term:1",
				Name:   "Dragon",
				Detail: "原文: Dragon",
				TitleParts: []usecase.ProcessingTargetListItemTitlePart{
					{Text: "Dragon"},
				},
				Metadata: []usecase.ProcessingTargetListItemMetadata{
					{Label: "候補", Value: "ドラゴン"},
				},
			},
		},
		Page:        request.Page,
		PageSize:    request.PageSize,
		TotalCount:  137,
		SearchQuery: request.SearchQuery,
	}, nil
}

func TestProcessingTargetControllerGetProcessingTargetListPreservesBoundaryFields(t *testing.T) {
	spy := &processingTargetReadModelUsecaseSpy{}
	controller := NewProcessingTargetController(spy)

	response, err := controller.GetProcessingTargetList(GetProcessingTargetListRequestDTO{
		JobID:       9,
		Phase:       "term_translation",
		Page:        2,
		PageSize:    50,
		SearchQuery: "Dragon",
	})
	if err != nil {
		t.Fatalf("GetProcessingTargetList returned error: %v", err)
	}

	if spy.request.JobID != 9 {
		t.Fatalf("JobID = %d, want 9", spy.request.JobID)
	}
	if spy.request.Phase != "term_translation" {
		t.Fatalf("Phase = %q, want term_translation", spy.request.Phase)
	}
	if spy.request.Page != 2 {
		t.Fatalf("Page = %d, want 2", spy.request.Page)
	}
	if spy.request.PageSize != 50 {
		t.Fatalf("PageSize = %d, want 50", spy.request.PageSize)
	}
	if spy.request.SearchQuery != "Dragon" {
		t.Fatalf("SearchQuery = %q, want Dragon", spy.request.SearchQuery)
	}
	if response.TotalCount != 137 {
		t.Fatalf("TotalCount = %d, want 137", response.TotalCount)
	}
	if response.Page != 2 {
		t.Fatalf("response Page = %d, want 2", response.Page)
	}
	if response.PageSize != 50 {
		t.Fatalf("response PageSize = %d, want 50", response.PageSize)
	}
	if response.SearchQuery != "Dragon" {
		t.Fatalf("response SearchQuery = %q, want Dragon", response.SearchQuery)
	}
	if len(response.Items) != 1 {
		t.Fatalf("len(response.Items) = %d, want 1", len(response.Items))
	}
	if response.Items[0].ID != "term:1" {
		t.Fatalf("response item ID = %q, want term:1", response.Items[0].ID)
	}
}
