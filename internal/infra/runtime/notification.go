package runtime

import (
	"context"
	"time"

	"aitranslationenginejp/internal/notification"
)

const (
	masterDictionaryImportProgressEventName  = "master-dictionary:import-progress"
	masterDictionaryImportCompletedEventName = "master-dictionary:import-completed"
	runtimeEventEmitterContextValueKey       = "events"
)

// NotificationContextProvider returns runtime context for Wails event emission.
type NotificationContextProvider func() (context.Context, bool)

// NotificationAdapter sends notifications through the Wails runtime event transport.
type NotificationAdapter struct {
	contextProvider NotificationContextProvider
}

// NewNotificationAdapter creates a Wails runtime notification adapter.
func NewNotificationAdapter(contextProvider NotificationContextProvider) *NotificationAdapter {
	return &NotificationAdapter{contextProvider: contextProvider}
}

// Send maps a transport-independent notification to the existing Wails runtime event contract.
func (adapter *NotificationAdapter) Send(_ context.Context, event notification.Notification) error {
	if adapter == nil || adapter.contextProvider == nil {
		return nil
	}
	runtimeContext, ok := adapter.contextProvider()
	if !ok {
		return nil
	}

	switch event.Kind {
	case notification.KindMasterDictionaryImportProgress:
		if event.Progress == nil {
			return nil
		}
		emitRuntimeEvent(runtimeContext, masterDictionaryImportProgressEventName, masterDictionaryImportProgressEventDTO{
			Progress: event.Progress.Percent,
		})
	case notification.KindMasterDictionaryImportCompleted:
		if event.Import == nil {
			return nil
		}
		emitRuntimeEvent(runtimeContext, masterDictionaryImportCompletedEventName, toImportCompletedEventDTO(*event.Import))
	}
	return nil
}

type wailsEventEmitter interface {
	Emit(eventName string, optionalData ...interface{})
}

func emitRuntimeEvent(runtimeContext context.Context, eventName string, payload interface{}) {
	events, ok := runtimeContext.Value(runtimeEventEmitterContextValueKey).(wailsEventEmitter)
	if !ok || events == nil {
		return
	}
	events.Emit(eventName, payload)
}

func toImportCompletedEventDTO(event notification.MasterDictionaryImportNotification) masterDictionaryImportCompletedEventDTO {
	items := make([]masterDictionaryEntryDTO, 0, len(event.Page.Items))
	for _, entry := range event.Page.Items {
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

	return masterDictionaryImportCompletedEventDTO{
		Page: masterDictionaryPageDTO{
			Items:      items,
			TotalCount: event.Page.TotalCount,
			Page:       event.Page.Page,
			PageSize:   event.Page.PageSize,
			SelectedID: event.Page.SelectedID,
		},
		Summary: masterDictionaryImportSummaryDTO{
			FilePath:      event.Summary.FilePath,
			FileName:      event.Summary.FileName,
			ImportedCount: event.Summary.ImportedCount,
			UpdatedCount:  event.Summary.UpdatedCount,
			SkippedCount:  event.Summary.SkippedCount,
			LastEntryID:   event.Summary.LastEntryID,
		},
		Refresh: masterDictionaryImportCompletedRefreshDTO{
			Query:           event.Refresh.Query,
			Category:        event.Refresh.Category,
			Page:            event.Refresh.Page,
			PageSize:        event.Refresh.PageSize,
			RefreshTargetID: event.Refresh.RefreshTargetID,
		},
	}
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
	FilePath      string `json:"filePath"`
	FileName      string `json:"fileName"`
	ImportedCount int    `json:"importedCount"`
	UpdatedCount  int    `json:"updatedCount"`
	SkippedCount  int    `json:"skippedCount"`
	LastEntryID   int64  `json:"lastEntryId"`
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
