package wails

import (
	"context"
	"fmt"
	"log/slog"

	"aitranslationenginejp/internal/usecase"
)

// ProcessingTargetReadModelUsecasePort defines the processing target read-model seam.
type ProcessingTargetReadModelUsecasePort interface {
	ListProcessingTargets(ctx context.Context, request usecase.ProcessingTargetListRequest) (usecase.ProcessingTargetListResult, error)
}

// ProcessingTargetController exposes Wails-bound processing target list entrypoints.
type ProcessingTargetController struct {
	processingTargetReadModelUsecase ProcessingTargetReadModelUsecasePort
}

// GetProcessingTargetListRequestDTO identifies one processing target list page.
type GetProcessingTargetListRequestDTO struct {
	JobID       int64  `json:"jobId"`
	Phase       string `json:"phase"`
	Page        int    `json:"page"`
	PageSize    int    `json:"pageSize"`
	SearchQuery string `json:"searchQuery"`
}

// ProcessingTargetListItemTitlePartDTO stores one title fragment.
type ProcessingTargetListItemTitlePartDTO struct {
	Text string `json:"text"`
}

// ProcessingTargetListItemMetadataDTO stores one display metadata pair.
type ProcessingTargetListItemMetadataDTO struct {
	Label string `json:"label"`
	Value string `json:"value"`
}

// ProcessingTargetListItemDTO stores one processing target item.
type ProcessingTargetListItemDTO struct {
	ID         string                                 `json:"id"`
	Name       string                                 `json:"name"`
	Detail     string                                 `json:"detail"`
	TitleParts []ProcessingTargetListItemTitlePartDTO `json:"titleParts"`
	Metadata   []ProcessingTargetListItemMetadataDTO  `json:"metadata"`
}

// GetProcessingTargetListResponseDTO returns one processing target list page.
type GetProcessingTargetListResponseDTO struct {
	Items       []ProcessingTargetListItemDTO         `json:"items"`
	Metadata    []ProcessingTargetListItemMetadataDTO `json:"metadata"`
	Page        int                                   `json:"page"`
	PageSize    int                                   `json:"pageSize"`
	TotalCount  int                                   `json:"totalCount"`
	SearchQuery string                                `json:"searchQuery"`
}

// NewProcessingTargetController creates a processing target controller.
func NewProcessingTargetController(usecase ProcessingTargetReadModelUsecasePort) *ProcessingTargetController {
	return &ProcessingTargetController{processingTargetReadModelUsecase: usecase}
}

// GetProcessingTargetList returns one processing target list page.
func (controller *ProcessingTargetController) GetProcessingTargetList(
	request GetProcessingTargetListRequestDTO,
) (GetProcessingTargetListResponseDTO, error) {
	if controller.processingTargetReadModelUsecase == nil {
		slog.WarnContext(context.Background(), "processing target list boundary failed",
			slog.String("event", "processing_target_list_boundary_failed"),
			slog.String("where", "backend.controller.wails.processing_target"),
			slog.String("result", "failed"),
			slog.String("id", fmt.Sprintf("job:%d", request.JobID)),
			slog.String("reason", "usecase_not_configured"),
		)
		return GetProcessingTargetListResponseDTO{}, fmt.Errorf("get processing target list: usecase is not configured")
	}
	result, err := controller.processingTargetReadModelUsecase.ListProcessingTargets(
		context.Background(),
		usecase.ProcessingTargetListRequest{
			JobID:       request.JobID,
			Phase:       request.Phase,
			Page:        request.Page,
			PageSize:    request.PageSize,
			SearchQuery: request.SearchQuery,
		},
	)
	if err != nil {
		slog.WarnContext(context.Background(), "processing target list boundary failed",
			slog.String("event", "processing_target_list_boundary_failed"),
			slog.String("where", "backend.controller.wails.processing_target"),
			slog.String("result", "failed"),
			slog.String("id", fmt.Sprintf("job:%d", request.JobID)),
			slog.String("reason", "usecase_failed"),
		)
		return GetProcessingTargetListResponseDTO{}, fmt.Errorf("get processing target list: %w", err)
	}
	return toGetProcessingTargetListResponseDTO(result), nil
}

func toGetProcessingTargetListResponseDTO(
	result usecase.ProcessingTargetListResult,
) GetProcessingTargetListResponseDTO {
	items := make([]ProcessingTargetListItemDTO, 0, len(result.Items))
	for _, item := range result.Items {
		items = append(items, toProcessingTargetListItemDTO(item))
	}
	return GetProcessingTargetListResponseDTO{
		Items:       items,
		Metadata:    []ProcessingTargetListItemMetadataDTO{},
		Page:        result.Page,
		PageSize:    result.PageSize,
		TotalCount:  result.TotalCount,
		SearchQuery: result.SearchQuery,
	}
}

func toProcessingTargetListItemDTO(
	item usecase.ProcessingTargetListItem,
) ProcessingTargetListItemDTO {
	titleParts := make([]ProcessingTargetListItemTitlePartDTO, 0, len(item.TitleParts))
	for _, part := range item.TitleParts {
		titleParts = append(titleParts, ProcessingTargetListItemTitlePartDTO{Text: part.Text})
	}
	metadata := make([]ProcessingTargetListItemMetadataDTO, 0, len(item.Metadata))
	for _, entry := range item.Metadata {
		metadata = append(metadata, ProcessingTargetListItemMetadataDTO{
			Label: entry.Label,
			Value: entry.Value,
		})
	}
	return ProcessingTargetListItemDTO{
		ID:         item.ID,
		Name:       item.Name,
		Detail:     item.Detail,
		TitleParts: titleParts,
		Metadata:   metadata,
	}
}
