package service

import (
	"context"
	"time"
)

const (
	masterDictionaryImportProgressEventName  = "master-dictionary:import-progress"
	masterDictionaryImportCompletedEventName = "master-dictionary:import-completed"
	runtimeEventEmitterContextValueKey       = "events"
)

// ContextProvider returns runtime context for Wails event emission.
type ContextProvider func() (context.Context, bool)

// WailsMasterDictionaryRuntimeEventPublisher publishes events through Wails runtime context values.
type WailsMasterDictionaryRuntimeEventPublisher struct {
	contextProvider ContextProvider
}

// NewWailsMasterDictionaryRuntimeEventPublisher creates a runtime event publisher.
func NewWailsMasterDictionaryRuntimeEventPublisher(
	contextProvider ContextProvider,
) *WailsMasterDictionaryRuntimeEventPublisher {
	return &WailsMasterDictionaryRuntimeEventPublisher{contextProvider: contextProvider}
}

// PublishImportProgress emits runtime import progress events.
func (publisher *WailsMasterDictionaryRuntimeEventPublisher) PublishImportProgress(
	_ context.Context,
	progress int,
) {
	runtimeContext, ok := publisher.contextProvider()
	if !ok {
		return
	}
	publisher.emitRuntimeEvent(
		runtimeContext,
		masterDictionaryImportProgressEventName,
		masterDictionaryImportProgressEventDTO{Progress: progress},
	)
}

// PublishImportCompleted emits runtime import completed events.
func (publisher *WailsMasterDictionaryRuntimeEventPublisher) PublishImportCompleted(
	_ context.Context,
	payload MasterDictionaryImportCompletedPayload,
) {
	runtimeContext, ok := publisher.contextProvider()
	if !ok {
		return
	}

	items := make([]masterDictionaryEntryDTO, 0, len(payload.Page.Items))
	for _, entry := range payload.Page.Items {
		items = append(items, masterDictionaryEntryDTO{
			ID:          entry.ID,
			Source:      entry.Source,
			Translation: entry.Translation,
			Category:    entry.Category,
			Origin:      entry.Origin,
			REC:         entry.REC,
			EDID:        entry.EDID,
			UpdatedAt:   entry.UpdatedAt.Format(time.RFC3339),
		})
	}

	publisher.emitRuntimeEvent(
		runtimeContext,
		masterDictionaryImportCompletedEventName,
		masterDictionaryImportCompletedEventDTO{
			Page: masterDictionaryPageDTO{
				Items:      items,
				TotalCount: payload.Page.TotalCount,
				Page:       payload.Page.Page,
				PageSize:   payload.Page.PageSize,
				SelectedID: payload.Page.SelectedID,
			},
			Summary: masterDictionaryImportSummaryDTO{
				FilePath:      payload.Summary.FilePath,
				FileName:      payload.Summary.FileName,
				ImportedCount: payload.Summary.ImportedCount,
				UpdatedCount:  payload.Summary.UpdatedCount,
				SkippedCount:  payload.Summary.SkippedCount,
				SelectedREC:   payload.Summary.SelectedREC,
				LastEntryID:   payload.Summary.LastEntryID,
			},
			Refresh: masterDictionaryImportCompletedRefreshDTO{
				Query:           payload.Refresh.Query,
				Category:        payload.Refresh.Category,
				Page:            payload.Refresh.Page,
				PageSize:        payload.Refresh.PageSize,
				RefreshTargetID: payload.Refresh.RefreshTargetID,
			},
		},
	)
}

type wailsEventEmitter interface {
	Emit(eventName string, optionalData ...interface{})
}

func (publisher *WailsMasterDictionaryRuntimeEventPublisher) emitRuntimeEvent(
	runtimeContext context.Context,
	eventName string,
	payload interface{},
) {
	events, ok := runtimeContext.Value(runtimeEventEmitterContextValueKey).(wailsEventEmitter)
	if !ok || events == nil {
		return
	}
	events.Emit(eventName, payload)
}

type masterDictionaryImportProgressEventDTO struct {
	Progress int `json:"progress"`
}

type masterDictionaryEntryDTO struct {
	ID          int64  `json:"id"`
	Source      string `json:"source"`
	Translation string `json:"translation"`
	Category    string `json:"category"`
	Origin      string `json:"origin"`
	REC         string `json:"rec"`
	EDID        string `json:"edid"`
	UpdatedAt   string `json:"updatedAt"`
}

type masterDictionaryPageDTO struct {
	Items      []masterDictionaryEntryDTO `json:"items"`
	TotalCount int                        `json:"totalCount"`
	Page       int                        `json:"page"`
	PageSize   int                        `json:"pageSize"`
	SelectedID *int64                     `json:"selectedId,omitempty"`
}

type masterDictionaryImportSummaryDTO struct {
	FilePath      string   `json:"filePath"`
	FileName      string   `json:"fileName"`
	ImportedCount int      `json:"importedCount"`
	UpdatedCount  int      `json:"updatedCount"`
	SkippedCount  int      `json:"skippedCount"`
	SelectedREC   []string `json:"selectedRec"`
	LastEntryID   int64    `json:"lastEntryId"`
}

type masterDictionaryImportCompletedRefreshDTO struct {
	Query           string `json:"query"`
	Category        string `json:"category"`
	Page            int    `json:"page"`
	PageSize        int    `json:"pageSize"`
	RefreshTargetID *int64 `json:"refreshTargetId,omitempty"`
}

type masterDictionaryImportCompletedEventDTO struct {
	Page    masterDictionaryPageDTO                   `json:"page"`
	Summary masterDictionaryImportSummaryDTO          `json:"summary"`
	Refresh masterDictionaryImportCompletedRefreshDTO `json:"refresh"`
}
