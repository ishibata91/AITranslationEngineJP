package runtime

import (
	"context"
	"database/sql"
	"encoding/json"
	"log/slog"
	"path/filepath"
	"testing"
	"time"

	_ "modernc.org/sqlite"
)

func TestSQLiteDiagnosticHandlerWritesStructuredLog(t *testing.T) {
	ctx := context.Background()
	databasePath := filepath.Join(t.TempDir(), "diagnostic-log.sqlite")
	diagnosticLog, err := OpenDiagnosticLog(ctx, databasePath)
	if err != nil {
		t.Fatalf("open diagnostic log: %v", err)
	}
	defer func() {
		if err := diagnosticLog.Close(); err != nil {
			t.Fatalf("close diagnostic log: %v", err)
		}
	}()

	now := time.Date(2026, 5, 9, 1, 2, 3, 4, time.UTC)
	handler := diagnosticLog.Handler(
		WithDiagnosticLogSource("backend-test"),
		WithDiagnosticLogClock(func() time.Time { return now }),
	).WithAttrs([]slog.Attr{
		slog.String("trace_id", "trace-1"),
		slog.String("boundary", "unit-test"),
	})
	record := slog.NewRecord(now, slog.LevelInfo, "diagnostic sqlite write", 0)
	record.AddAttrs(
		slog.String("event_name", "diagnostic_sqlite_write"),
		slog.String("screen", "N/A"),
		slog.Int("count", 2),
	)
	if err := handler.Handle(ctx, record); err != nil {
		t.Fatalf("handle diagnostic log record: %v", err)
	}

	row := readDiagnosticLogRow(t, databasePath)
	if row.OccurredAt != "2026-05-09T01:02:03.000000004Z" {
		t.Fatalf("unexpected occurred_at: %q", row.OccurredAt)
	}
	if row.Level != "info" {
		t.Fatalf("unexpected level: %q", row.Level)
	}
	if row.Source != "backend-test" {
		t.Fatalf("unexpected source: %q", row.Source)
	}
	if row.TraceID.String != "trace-1" || !row.TraceID.Valid {
		t.Fatalf("unexpected trace_id: %#v", row.TraceID)
	}
	if row.EventName != "diagnostic_sqlite_write" {
		t.Fatalf("unexpected event_name: %q", row.EventName)
	}

	var attrs map[string]string
	if err := json.Unmarshal([]byte(row.AttrsJSON), &attrs); err != nil {
		t.Fatalf("unmarshal attrs_json: %v", err)
	}
	if attrs["boundary"] != "unit-test" {
		t.Fatalf("unexpected boundary attr: %#v", attrs)
	}
	if attrs["count"] != "2" {
		t.Fatalf("unexpected count attr: %#v", attrs)
	}
}

func TestSQLiteDiagnosticHandlerHonorsMinimumLevel(t *testing.T) {
	ctx := context.Background()
	databasePath := filepath.Join(t.TempDir(), "diagnostic-log.sqlite")
	diagnosticLog, err := OpenDiagnosticLog(ctx, databasePath)
	if err != nil {
		t.Fatalf("open diagnostic log: %v", err)
	}
	defer func() {
		if err := diagnosticLog.Close(); err != nil {
			t.Fatalf("close diagnostic log: %v", err)
		}
	}()

	logger := slog.New(diagnosticLog.Handler(WithDiagnosticLogLevel(slog.LevelWarn)))
	logger.InfoContext(ctx, "ignored info", "event_name", "ignored_info")
	logger.WarnContext(ctx, "kept warn", "event_name", "kept_warn")

	row := readDiagnosticLogRow(t, databasePath)
	if row.Level != "warn" {
		t.Fatalf("unexpected level: %q", row.Level)
	}
	if row.EventName != "kept_warn" {
		t.Fatalf("unexpected event_name: %q", row.EventName)
	}
}

type diagnosticLogRow struct {
	OccurredAt string
	Level      string
	Source     string
	TraceID    sql.NullString
	EventName  string
	AttrsJSON  string
}

func readDiagnosticLogRow(t *testing.T, databasePath string) diagnosticLogRow {
	t.Helper()
	db, err := sql.Open(sqliteDriverName, diagnosticSQLiteDSN(databasePath))
	if err != nil {
		t.Fatalf("open written diagnostic database: %v", err)
	}
	defer func() {
		if err := db.Close(); err != nil {
			t.Fatalf("close written diagnostic database: %v", err)
		}
	}()

	var row diagnosticLogRow
	if err := db.QueryRowContext(context.Background(), `
SELECT occurred_at, level, source, trace_id, event_name, attrs_json
FROM diagnostic_log
ORDER BY id
LIMIT 1;`).Scan(
		&row.OccurredAt,
		&row.Level,
		&row.Source,
		&row.TraceID,
		&row.EventName,
		&row.AttrsJSON,
	); err != nil {
		t.Fatalf("select diagnostic log row: %v", err)
	}
	return row
}
