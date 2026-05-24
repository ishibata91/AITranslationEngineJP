package usecase

import (
	"context"
	"fmt"

	processingtargetservice "aitranslationenginejp/internal/service"
)

type processingTargetReadModelServicePort interface {
	ListProcessingTargets(
		ctx context.Context,
		request processingtargetservice.ProcessingTargetListRequest,
	) (processingtargetservice.ProcessingTargetListReadModel, error)
}

// ProcessingTargetReadModelUsecase adapts processing target pages for backend callers.
type ProcessingTargetReadModelUsecase struct {
	service processingTargetReadModelServicePort
}

// NewProcessingTargetReadModelUsecase builds the processing target read-model usecase.
func NewProcessingTargetReadModelUsecase(service processingTargetReadModelServicePort) *ProcessingTargetReadModelUsecase {
	return &ProcessingTargetReadModelUsecase{service: service}
}

// ListProcessingTargets returns one processing target page.
func (usecase *ProcessingTargetReadModelUsecase) ListProcessingTargets(
	ctx context.Context,
	request ProcessingTargetListRequest,
) (ProcessingTargetListResult, error) {
	result, err := usecase.service.ListProcessingTargets(ctx, processingtargetservice.ProcessingTargetListRequest{
		JobID:       request.JobID,
		Phase:       request.Phase,
		Page:        request.Page,
		PageSize:    request.PageSize,
		SearchQuery: request.SearchQuery,
	})
	if err != nil {
		return ProcessingTargetListResult{}, fmt.Errorf("list processing targets: %w", err)
	}
	return toProcessingTargetListResult(result), nil
}

func toProcessingTargetListResult(
	source processingtargetservice.ProcessingTargetListReadModel,
) ProcessingTargetListResult {
	items := make([]ProcessingTargetListItem, 0, len(source.Items))
	for _, item := range source.Items {
		items = append(items, toProcessingTargetListItem(item))
	}
	return ProcessingTargetListResult{
		Items:       items,
		Page:        source.Page,
		PageSize:    source.PageSize,
		TotalCount:  source.TotalCount,
		SearchQuery: source.SearchQuery,
	}
}

func toProcessingTargetListItem(
	source processingtargetservice.ProcessingTargetListItemReadModel,
) ProcessingTargetListItem {
	titleParts := make([]ProcessingTargetListItemTitlePart, 0, len(source.TitleParts))
	for _, part := range source.TitleParts {
		titleParts = append(titleParts, ProcessingTargetListItemTitlePart{Text: part.Text})
	}
	metadata := make([]ProcessingTargetListItemMetadata, 0, len(source.Metadata))
	for _, item := range source.Metadata {
		metadata = append(metadata, ProcessingTargetListItemMetadata{
			Label: item.Label,
			Value: item.Value,
		})
	}
	return ProcessingTargetListItem{
		ID:         source.ID,
		Name:       source.Name,
		Detail:     source.Detail,
		TitleParts: titleParts,
		Metadata:   metadata,
	}
}
