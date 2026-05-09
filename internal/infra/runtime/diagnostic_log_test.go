package runtime

import (
	"bytes"
	"context"
	"encoding/json"
	"log/slog"
	"strings"
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

func TestInstallDiagnosticLoggerPayloadHasNoForbiddenOrUndefinedLikeValues(t *testing.T) {
	var buffer bytes.Buffer
	InstallDiagnosticLogger(&buffer, "backend-test", slog.LevelInfo)

	slog.InfoContext(context.Background(), "provider boundary skipped",
		slog.String("event", "provider_execution_settings"),
		slog.String("where", "provider_settings.service"),
		slog.String("result", "skipped"),
		slog.String("id", "job:12"),
		slog.String("reason", "provider_skipped"),
		slog.Int("count", 3),
	)

	var payload map[string]any
	if err := json.Unmarshal(buffer.Bytes(), &payload); err != nil {
		t.Fatalf("unmarshal diagnostic log payload: %v", err)
	}

	requiredKeys := []string{"event", "where", "result", "id", "reason", "count"}
	for _, key := range requiredKeys {
		if _, ok := payload[key]; !ok {
			t.Fatalf("expected key %q in payload: %#v", key, payload)
		}
	}

	forbiddenKeys := []string{
		"api_key",
		"apikey",
		"endpoint",
		"raw_request",
		"raw_response",
		"prompt",
		"full_text",
		"xml",
		"dto",
		"full_path",
		"trace_id",
	}
	for _, key := range forbiddenKeys {
		if _, exists := payload[key]; exists {
			t.Fatalf("forbidden key %q found in payload: %#v", key, payload)
		}
	}

	raw := buffer.String()
	forbiddenFragments := []string{
		"undefined",
		"null",
		"api_key",
		"trace_id",
		"/Users/",
		"https://",
	}
	for _, fragment := range forbiddenFragments {
		if strings.Contains(raw, fragment) {
			t.Fatalf("forbidden fragment %q found in payload: %s", fragment, raw)
		}
	}
}
