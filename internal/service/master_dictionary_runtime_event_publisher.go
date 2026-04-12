package service

import "context"

// MasterDictionaryImportRefreshPolicy carries approved refresh semantics for import completion.
type MasterDictionaryImportRefreshPolicy struct {
	Query           string
	Category        string
	Page            int
	PageSize        int
	RefreshTargetID *int64
}

// MasterDictionaryImportCompletedPage carries list/detail state for runtime completed events.
type MasterDictionaryImportCompletedPage struct {
	Items      []MasterDictionaryEntry
	TotalCount int
	Page       int
	PageSize   int
	SelectedID *int64
}

// MasterDictionaryImportCompletedPayload is emitted when XML import has completed.
type MasterDictionaryImportCompletedPayload struct {
	Page    MasterDictionaryImportCompletedPage
	Summary MasterDictionaryImportSummary
	Refresh MasterDictionaryImportRefreshPolicy
}

// MasterDictionaryRuntimeEventPublisher publishes import progress/completed events.
type MasterDictionaryRuntimeEventPublisher interface {
	PublishImportProgress(ctx context.Context, progress int)
	PublishImportCompleted(ctx context.Context, payload MasterDictionaryImportCompletedPayload)
}
