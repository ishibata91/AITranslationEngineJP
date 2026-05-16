package runtime

import (
	"context"
	"testing"
	"time"

	"aitranslationenginejp/internal/notification"
)

type recordingNotificationEmitter struct {
	events []recordedNotificationEvent
}

type recordedNotificationEvent struct {
	name    string
	payload []interface{}
}

type notificationRuntimeContext struct {
	context.Context
	emitter *recordingNotificationEmitter
}

func (emitter *recordingNotificationEmitter) Emit(eventName string, optionalData ...interface{}) {
	emitter.events = append(emitter.events, recordedNotificationEvent{name: eventName, payload: optionalData})
}

func (ctx notificationRuntimeContext) Value(key interface{}) interface{} {
	if key == "events" {
		return ctx.emitter
	}
	return ctx.Context.Value(key)
}

func TestNotificationAdapterSendsImportProgressRuntimeEvent(t *testing.T) {
	emitter := &recordingNotificationEmitter{}
	adapter := NewNotificationAdapter(func() (context.Context, bool) {
		return notificationRuntimeContext{Context: context.Background(), emitter: emitter}, true
	})

	err := adapter.Send(context.Background(), notification.Notification{
		Kind:     notification.KindMasterDictionaryImportProgress,
		Progress: &notification.ProgressNotification{Percent: 42},
	})

	if err != nil {
		t.Fatalf("expected progress send to succeed: %v", err)
	}
	if len(emitter.events) != 1 || emitter.events[0].name != masterDictionaryImportProgressEventName {
		t.Fatalf("expected progress event, got %#v", emitter.events)
	}
	payload, ok := emitter.events[0].payload[0].(masterDictionaryImportProgressEventDTO)
	if !ok || payload.Progress != 42 {
		t.Fatalf("expected progress payload, got %#v", emitter.events[0].payload)
	}
}

func TestNotificationAdapterSendsImportCompletedRuntimeEventContract(t *testing.T) {
	selectedID := int64(7)
	updatedAt := time.Date(2026, 5, 16, 1, 2, 3, 0, time.UTC)
	emitter := &recordingNotificationEmitter{}
	adapter := NewNotificationAdapter(func() (context.Context, bool) {
		return notificationRuntimeContext{Context: context.Background(), emitter: emitter}, true
	})

	err := adapter.Send(context.Background(), notification.Notification{
		Kind: notification.KindMasterDictionaryImportCompleted,
		Import: &notification.MasterDictionaryImportNotification{
			Page: notification.MasterDictionaryPage{
				Items: []notification.MasterDictionaryEntry{{
					ID:          selectedID,
					Source:      "Auriel",
					Translation: "アーリエル",
					Category:    "武器",
					Origin:      "import",
					REC:         "WEAP",
					EDID:        "AurielBow",
					UpdatedAt:   updatedAt,
				}},
				TotalCount: 1,
				Page:       1,
				PageSize:   30,
				SelectedID: &selectedID,
			},
			Summary: notification.MasterDictionaryImportSummary{
				FileName:      "dict.xml",
				ImportedCount: 1,
				UpdatedCount:  2,
				SkippedCount:  3,
				LastEntryID:   selectedID,
			},
			Refresh: notification.MasterDictionaryImportRefresh{
				Query:           "Auriel",
				Category:        "武器",
				Page:            1,
				PageSize:        30,
				RefreshTargetID: &selectedID,
			},
		},
	})

	if err != nil {
		t.Fatalf("expected completed send to succeed: %v", err)
	}
	if len(emitter.events) != 1 || emitter.events[0].name != masterDictionaryImportCompletedEventName {
		t.Fatalf("expected completed event, got %#v", emitter.events)
	}
	payload, ok := emitter.events[0].payload[0].(masterDictionaryImportCompletedEventDTO)
	if !ok {
		t.Fatalf("expected completed payload type, got %#v", emitter.events[0].payload)
	}
	if payload.Page.Items[0].UpdatedAt != "2026-05-16T01:02:03Z" || payload.Summary.FileName != "dict.xml" {
		t.Fatalf("expected existing completed event contract, got %#v", payload)
	}
	if payload.Refresh.RefreshTargetID == nil || *payload.Refresh.RefreshTargetID != selectedID {
		t.Fatalf("expected refresh target id, got %#v", payload.Refresh)
	}
}

func TestNotificationAdapterSkipsWhenRuntimeContextIsMissing(t *testing.T) {
	adapter := NewNotificationAdapter(func() (context.Context, bool) {
		return nil, false
	})

	err := adapter.Send(context.Background(), notification.Notification{
		Kind:     notification.KindMasterDictionaryImportProgress,
		Progress: &notification.ProgressNotification{Percent: 42},
	})

	if err != nil {
		t.Fatalf("expected missing runtime context to be ignored: %v", err)
	}
}

func TestNotificationAdapterDoesNotRedactOrJudgeState(t *testing.T) {
	emitter := &recordingNotificationEmitter{}
	adapter := NewNotificationAdapter(func() (context.Context, bool) {
		return notificationRuntimeContext{Context: context.Background(), emitter: emitter}, true
	})

	err := adapter.Send(context.Background(), notification.Notification{
		Kind: notification.KindMasterDictionaryImportCompleted,
		Import: &notification.MasterDictionaryImportNotification{
			Summary: notification.MasterDictionaryImportSummary{
				FilePath: "/already-redacted-by-dispatcher.xml",
				FileName: "already-redacted-by-dispatcher.xml",
			},
		},
	})

	if err != nil {
		t.Fatalf("expected completed send to succeed: %v", err)
	}
	if len(emitter.events) != 1 || emitter.events[0].name != masterDictionaryImportCompletedEventName {
		t.Fatalf("expected one completed event, got %#v", emitter.events)
	}
	payload, ok := emitter.events[0].payload[0].(masterDictionaryImportCompletedEventDTO)
	if !ok {
		t.Fatalf("expected completed payload type, got %#v", emitter.events[0].payload)
	}
	if payload.Summary.FilePath != "/already-redacted-by-dispatcher.xml" {
		t.Fatalf("expected adapter not to redact file path, got %#v", payload.Summary)
	}
}

func TestNotificationAdapterIgnoresDiscardedFactWithoutStateDecision(t *testing.T) {
	emitter := &recordingNotificationEmitter{}
	adapter := NewNotificationAdapter(func() (context.Context, bool) {
		return notificationRuntimeContext{Context: context.Background(), emitter: emitter}, true
	})

	err := adapter.Send(context.Background(), notification.Notification{
		Kind: notification.KindMasterDictionaryImportDiscarded,
		Import: &notification.MasterDictionaryImportNotification{
			Reason: "redacted",
		},
	})

	if err != nil {
		t.Fatalf("expected discarded send to be ignored: %v", err)
	}
	if len(emitter.events) != 0 {
		t.Fatalf("expected no runtime event for discarded kind, got %#v", emitter.events)
	}
}
