package service

import (
	"context"
	"fmt"
	"strings"

	"aitranslationenginejp/internal/repository"
)

const (
	processingTargetDefaultPageSize = 50
	processingTargetMaxPageSize     = 200
)

// ProcessingTargetListRequest describes one service read-model request.
type ProcessingTargetListRequest struct {
	JobID       int64
	Phase       string
	Page        int
	PageSize    int
	SearchQuery string
}

// ProcessingTargetListItemReadModel stores one processing target item.
type ProcessingTargetListItemReadModel struct {
	ID         string
	Name       string
	Detail     string
	TitleParts []ProcessingTargetListItemTitlePartReadModel
	Metadata   []ProcessingTargetListItemMetadataReadModel
}

// ProcessingTargetListItemTitlePartReadModel stores one title fragment.
type ProcessingTargetListItemTitlePartReadModel struct {
	Text string
}

// ProcessingTargetListItemMetadataReadModel stores one metadata pair.
type ProcessingTargetListItemMetadataReadModel struct {
	Label string
	Value string
}

// ProcessingTargetListReadModel stores a paged processing target list.
type ProcessingTargetListReadModel struct {
	Items       []ProcessingTargetListItemReadModel
	Page        int
	PageSize    int
	TotalCount  int
	SearchQuery string
}

type processingTargetReadRepository interface {
	ListProcessingTargets(ctx context.Context, query repository.ProcessingTargetListQuery) (repository.ProcessingTargetListResult, error)
}

// ProcessingTargetReadModelService loads processing target read models.
type ProcessingTargetReadModelService struct {
	repository processingTargetReadRepository
}

// NewProcessingTargetReadModelService builds the processing target read-model service.
func NewProcessingTargetReadModelService(repository processingTargetReadRepository) *ProcessingTargetReadModelService {
	return &ProcessingTargetReadModelService{repository: repository}
}

// ListProcessingTargets returns the requested processing target page.
func (service *ProcessingTargetReadModelService) ListProcessingTargets(
	ctx context.Context,
	request ProcessingTargetListRequest,
) (ProcessingTargetListReadModel, error) {
	if service.repository == nil {
		return ProcessingTargetListReadModel{}, fmt.Errorf("processing target repository is not configured")
	}
	normalized := normalizeProcessingTargetListRequest(request)
	result, err := service.repository.ListProcessingTargets(ctx, repository.ProcessingTargetListQuery{
		JobID:       normalized.JobID,
		Phase:       normalized.Phase,
		Page:        normalized.Page,
		PageSize:    normalized.PageSize,
		SearchQuery: normalized.SearchQuery,
	})
	if err != nil {
		return ProcessingTargetListReadModel{}, fmt.Errorf("list processing targets: %w", err)
	}
	return toProcessingTargetListReadModel(result), nil
}

func normalizeProcessingTargetListRequest(request ProcessingTargetListRequest) ProcessingTargetListRequest {
	page := request.Page
	if page < 1 {
		page = 1
	}
	pageSize := request.PageSize
	if pageSize < 1 {
		pageSize = processingTargetDefaultPageSize
	}
	if pageSize > processingTargetMaxPageSize {
		pageSize = processingTargetMaxPageSize
	}
	return ProcessingTargetListRequest{
		JobID:       request.JobID,
		Phase:       normalizeProcessingTargetPhase(request.Phase),
		Page:        page,
		PageSize:    pageSize,
		SearchQuery: strings.TrimSpace(request.SearchQuery),
	}
}

func normalizeProcessingTargetPhase(phase string) string {
	switch strings.TrimSpace(phase) {
	case "word_translation":
		return "term_translation"
	case "npc_persona_generation":
		return "persona_generation"
	case "translation_complete", "completion", "confirmation":
		return "translation_complete"
	default:
		return strings.TrimSpace(phase)
	}
}

func toProcessingTargetListReadModel(
	source repository.ProcessingTargetListResult,
) ProcessingTargetListReadModel {
	items := make([]ProcessingTargetListItemReadModel, 0, len(source.Items))
	for _, item := range source.Items {
		items = append(items, toProcessingTargetItemReadModel(item))
	}
	return ProcessingTargetListReadModel{
		Items:       items,
		Page:        source.Page,
		PageSize:    source.PageSize,
		TotalCount:  source.TotalCount,
		SearchQuery: source.SearchQuery,
	}
}

func toProcessingTargetItemReadModel(
	source repository.ProcessingTargetListItem,
) ProcessingTargetListItemReadModel {
	titleParts := make([]ProcessingTargetListItemTitlePartReadModel, 0, len(source.TitleParts))
	for _, part := range source.TitleParts {
		titleParts = append(titleParts, ProcessingTargetListItemTitlePartReadModel{Text: part})
	}
	metadata := make([]ProcessingTargetListItemMetadataReadModel, 0, len(source.Metadata))
	for _, item := range source.Metadata {
		metadata = append(metadata, ProcessingTargetListItemMetadataReadModel{
			Label: item.Label,
			Value: item.Value,
		})
	}
	return ProcessingTargetListItemReadModel{
		ID:         source.ID,
		Name:       source.Name,
		Detail:     source.Detail,
		TitleParts: titleParts,
		Metadata:   metadata,
	}
}
