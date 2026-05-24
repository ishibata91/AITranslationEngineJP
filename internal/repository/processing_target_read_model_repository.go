package repository

import "context"

// ProcessingTargetListQuery describes one backend read-model page request.
type ProcessingTargetListQuery struct {
	JobID       int64
	Phase       string
	Page        int
	PageSize    int
	SearchQuery string
}

// ProcessingTargetListItem stores one processing target row.
type ProcessingTargetListItem struct {
	ID         string
	Name       string
	Detail     string
	TitleParts []string
	Metadata   []ProcessingTargetListItemMetadata
}

// ProcessingTargetListItemMetadata stores one display metadata pair.
type ProcessingTargetListItemMetadata struct {
	Label string
	Value string
}

// ProcessingTargetListResult stores a paged processing target list.
type ProcessingTargetListResult struct {
	Items       []ProcessingTargetListItem
	Page        int
	PageSize    int
	TotalCount  int
	SearchQuery string
}

// ProcessingTargetReadModelRepository defines processing target read-only persistence.
type ProcessingTargetReadModelRepository interface {
	ListProcessingTargets(ctx context.Context, query ProcessingTargetListQuery) (ProcessingTargetListResult, error)
}
