package notification

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"log/slog"
	"strings"
	"testing"
)

func TestDispatcherSendsRedactedImportCompletedNotification(t *testing.T) {
	capture := captureNotificationDispatchLog(t)
	port := &recordingPort{}
	dispatcher := NewDispatcher(port)

	result := dispatcher.Dispatch(context.Background(), Fact{
		Kind: KindMasterDictionaryImportCompleted,
		Import: &MasterDictionaryImportFact{
			Summary: MasterDictionaryImportSummary{
				FilePath:      "/tmp/Dawnguard_english_japanese.xml",
				FileName:      "",
				ImportedCount: 3,
				LastEntryID:   88,
			},
		},
	})

	if !result.Sent || result.Err != nil {
		t.Fatalf("expected notification to be sent, got %#v", result)
	}
	if len(port.sent) != 1 {
		t.Fatalf("expected one notification, got %d", len(port.sent))
	}
	summary := port.sent[0].Import.Summary
	if summary.FilePath != "" {
		t.Fatalf("expected file path redaction, got %q", summary.FilePath)
	}
	if summary.FileName != "Dawnguard_english_japanese.xml" {
		t.Fatalf("expected safe file name, got %q", summary.FileName)
	}
	payload := capture.requireEvent(t, "sent", "")
	if payload["id"] != string(KindMasterDictionaryImportCompleted) {
		t.Fatalf("expected dispatch log to identify notification kind, got %#v", payload)
	}
	assertNotificationDispatchLogExcludesForbiddenValues(t, capture)
}

func TestDispatcherDispatchTableForRedactionAndSendability(t *testing.T) {
	tests := []struct {
		name         string
		port         Port
		fact         Fact
		wantSent     bool
		wantErr      error
		wantSupposed bool
		wantResult   string
		wantReason   string
	}{
		{
			name: "redaction failure suppresses and never calls transport",
			port: &recordingPort{},
			fact: Fact{
				Kind: KindMasterDictionaryImportCompleted,
				Import: &MasterDictionaryImportFact{
					Summary: MasterDictionaryImportSummary{FilePath: "/tmp/secret-provider-raw.xml"},
				},
			},
			wantErr:      ErrUnsafePayload,
			wantSupposed: true,
			wantResult:   "rejected",
			wantReason:   "unsafe_payload",
		},
		{
			name: "not sendable fact is skipped",
			port: &recordingPort{},
			fact: Fact{
				Kind: KindMasterDictionaryImportProgress,
			},
			wantSupposed: true,
			wantResult:   "skipped",
			wantReason:   "not_sendable",
		},
		{
			name: "missing transport is skipped",
			port: nil,
			fact: Fact{
				Kind:     KindMasterDictionaryImportProgress,
				Progress: &ProgressFact{Percent: 42},
			},
			wantSupposed: true,
			wantResult:   "skipped",
			wantReason:   "transport_unavailable",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			capture := captureNotificationDispatchLog(t)
			dispatcher := NewDispatcher(tt.port)
			recordingPort, _ := tt.port.(*recordingPort)

			result := dispatcher.Dispatch(context.Background(), tt.fact)

			if result.Sent != tt.wantSent {
				t.Fatalf("expected sent=%t, got %#v", tt.wantSent, result)
			}
			if result.Suppressed != tt.wantSupposed {
				t.Fatalf("expected suppressed=%t, got %#v", tt.wantSupposed, result)
			}
			if !errors.Is(result.Err, tt.wantErr) {
				t.Fatalf("expected err %v, got %#v", tt.wantErr, result)
			}
			if recordingPort != nil && len(recordingPort.sent) != 0 {
				t.Fatalf("expected no transport calls, got %d", len(recordingPort.sent))
			}
			capture.requireEvent(t, tt.wantResult, tt.wantReason)
			assertNotificationDispatchLogExcludesForbiddenValues(t, capture)
		})
	}
}

func TestDispatcherKeepsSendFailureLocal(t *testing.T) {
	capture := captureNotificationDispatchLog(t)
	sendErr := errors.New("send failed")
	dispatcher := NewDispatcher(failingPort{err: sendErr})

	result := dispatcher.Dispatch(context.Background(), Fact{
		Kind:     KindMasterDictionaryImportProgress,
		Progress: &ProgressFact{Percent: 120},
	})

	if result.Sent || !errors.Is(result.Err, sendErr) {
		t.Fatalf("expected local send failure result, got %#v", result)
	}
	capture.requireEvent(t, "failed", "transport_error")
	assertNotificationDispatchLogExcludesForbiddenValues(t, capture)
}

func TestDispatcherLogsSuppressedNotification(t *testing.T) {
	capture := captureNotificationDispatchLog(t)
	dispatcher := NewDispatcher(nil)

	result := dispatcher.Dispatch(context.Background(), Fact{
		Kind: KindMasterDictionaryImportProgress,
	})

	if !result.Suppressed || result.Err != nil {
		t.Fatalf("expected not-sendable suppression, got %#v", result)
	}
	capture.requireEvent(t, "skipped", "not_sendable")
	assertNotificationDispatchLogExcludesForbiddenValues(t, capture)
}

func TestDispatcherLogsUnavailableTransport(t *testing.T) {
	capture := captureNotificationDispatchLog(t)
	dispatcher := NewDispatcher(nil)

	result := dispatcher.Dispatch(context.Background(), Fact{
		Kind:     KindMasterDictionaryImportProgress,
		Progress: &ProgressFact{Percent: 42},
	})

	if !result.Suppressed || result.Err != nil {
		t.Fatalf("expected unavailable transport suppression, got %#v", result)
	}
	capture.requireEvent(t, "skipped", "transport_unavailable")
	assertNotificationDispatchLogExcludesForbiddenValues(t, capture)
}

type recordingPort struct {
	sent []Notification
}

func (port *recordingPort) Send(_ context.Context, notification Notification) error {
	port.sent = append(port.sent, notification)
	return nil
}

type failingPort struct {
	err error
}

func (port failingPort) Send(context.Context, Notification) error {
	return port.err
}

type notificationDispatchLogCapture struct {
	buffer *bytes.Buffer
}

func captureNotificationDispatchLog(t *testing.T) notificationDispatchLogCapture {
	t.Helper()

	var buffer bytes.Buffer
	previous := slog.Default()
	slog.SetDefault(slog.New(slog.NewJSONHandler(&buffer, nil)))
	t.Cleanup(func() { slog.SetDefault(previous) })
	return notificationDispatchLogCapture{buffer: &buffer}
}

func (capture notificationDispatchLogCapture) requireEvent(
	t *testing.T,
	result string,
	reason string,
) map[string]any {
	t.Helper()

	var payload map[string]any
	if err := json.Unmarshal(capture.buffer.Bytes(), &payload); err != nil {
		t.Fatalf("expected JSON notification dispatch log, got %q: %v", capture.buffer.String(), err)
	}
	if payload["event"] != "notification_dispatch" ||
		payload["where"] != "backend.notification.dispatcher" ||
		payload["result"] != result {
		t.Fatalf("expected notification dispatch result=%q, got %#v", result, payload)
	}
	if reason != "" && payload["reason"] != reason {
		t.Fatalf("expected notification dispatch reason=%q, got %#v", reason, payload)
	}
	return payload
}

func assertNotificationDispatchLogExcludesForbiddenValues(t *testing.T, capture notificationDispatchLogCapture) {
	t.Helper()

	logText := capture.buffer.String()
	for _, forbidden := range []string{
		"/tmp/Dawnguard_english_japanese.xml",
		"/tmp/secret-provider-raw.xml",
		"secret-provider-raw.xml",
		"send failed",
		"api key",
		"provider raw",
		"<strings>",
	} {
		if strings.Contains(logText, forbidden) {
			t.Fatalf("expected notification dispatch log to exclude forbidden value %q, got %s", forbidden, logText)
		}
	}
}
