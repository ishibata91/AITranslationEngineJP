package runtime

import (
	"bytes"
	"context"
	"encoding/json"
	"log/slog"
	"testing"
)

func TestInstallDiagnosticLoggerWritesJSONToWriter(t *testing.T) {
	var buffer bytes.Buffer
	InstallDiagnosticLogger(&buffer, "backend-test", slog.LevelInfo)

	slog.InfoContext(context.Background(), "diagnostic event",
		slog.String("event", "diagnostic_event"),
		slog.String("where", "backend.runtime"),
		slog.String("result", "success"),
		slog.Int("count", 2),
	)

	var payload map[string]any
	if err := json.Unmarshal(buffer.Bytes(), &payload); err != nil {
		t.Fatalf("unmarshal diagnostic log payload: %v", err)
	}
	if payload["source"] != "backend-test" {
		t.Fatalf("unexpected source: %#v", payload)
	}
	if payload["event"] != "diagnostic_event" {
		t.Fatalf("unexpected event: %#v", payload)
	}
	if payload["where"] != "backend.runtime" {
		t.Fatalf("unexpected where: %#v", payload)
	}
	if payload["result"] != "success" {
		t.Fatalf("unexpected result: %#v", payload)
	}
	if payload["count"] != float64(2) {
		t.Fatalf("unexpected count: %#v", payload)
	}
}

func TestInstallDiagnosticLoggerHonorsMinimumLevel(t *testing.T) {
	var buffer bytes.Buffer
	InstallDiagnosticLogger(&buffer, "backend-test", slog.LevelWarn)

	slog.InfoContext(context.Background(), "ignored info",
		slog.String("event", "ignored_info"),
	)
	slog.WarnContext(context.Background(), "kept warn",
		slog.String("event", "kept_warn"),
	)

	var payload map[string]any
	if err := json.Unmarshal(buffer.Bytes(), &payload); err != nil {
		t.Fatalf("unmarshal diagnostic log payload: %v", err)
	}
	if payload["event"] != "kept_warn" {
		t.Fatalf("unexpected event: %#v", payload)
	}
}
