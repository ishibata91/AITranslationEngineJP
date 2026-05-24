package usecase

// ProcessingTargetListRequest identifies one processing target list page.
type ProcessingTargetListRequest struct {
	JobID       int64
	Phase       string
	Page        int
	PageSize    int
	SearchQuery string
}

// ProcessingTargetListItem stores one processing target item.
type ProcessingTargetListItem struct {
	ID         string
	Name       string
	Detail     string
	TitleParts []ProcessingTargetListItemTitlePart
	Metadata   []ProcessingTargetListItemMetadata
}

// ProcessingTargetListItemTitlePart stores one title fragment.
type ProcessingTargetListItemTitlePart struct {
	Text string
}

// ProcessingTargetListItemMetadata stores one metadata pair.
type ProcessingTargetListItemMetadata struct {
	Label string
	Value string
}

// ProcessingTargetListResult stores one processing target list page.
type ProcessingTargetListResult struct {
	Items       []ProcessingTargetListItem
	Page        int
	PageSize    int
	TotalCount  int
	SearchQuery string
}
