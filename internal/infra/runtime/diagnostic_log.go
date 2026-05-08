package runtime

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	// Register the modernc SQLite driver for database/sql.
	_ "modernc.org/sqlite"
)

const (
	// DefaultDiagnosticLogPath is the default temporary SQLite path for backend diagnostic logs.
	DefaultDiagnosticLogPath = "tmp/observability-logging-foundation/diagnostic-log.sqlite"
	sqliteDriverName         = "sqlite"
)

const createDiagnosticLogTableSQL = `
CREATE TABLE IF NOT EXISTS diagnostic_log (
  id INTEGER PRIMARY KEY,
  occurred_at TEXT NOT NULL,
  level TEXT NOT NULL,
  source TEXT NOT NULL,
  trace_id TEXT,
  screen TEXT,
  boundary TEXT,
  event_name TEXT NOT NULL,
  message TEXT NOT NULL,
  attrs_json TEXT NOT NULL
);`

const insertDiagnosticLogSQL = `
INSERT INTO diagnostic_log (
  occurred_at,
  level,
  source,
  trace_id,
  screen,
  boundary,
  event_name,
  message,
  attrs_json
) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?);`

// DiagnosticLog owns the SQLite database used by diagnostic slog handlers.
type DiagnosticLog struct {
	db *sql.DB
}

// DiagnosticLogOption customizes a diagnostic log or handler.
type DiagnosticLogOption func(*diagnosticLogConfig)

type diagnosticLogConfig struct {
	minLevel slog.Leveler
	source   string
	now      func() time.Time
}

// OpenDiagnosticLog opens the diagnostic SQLite database and prepares its table.
func OpenDiagnosticLog(ctx context.Context, databasePath string) (*DiagnosticLog, error) {
	resolvedPath, err := resolveDiagnosticLogPath(databasePath)
	if err != nil {
		return nil, err
	}

	db, err := sql.Open(sqliteDriverName, diagnosticSQLiteDSN(resolvedPath))
	if err != nil {
		return nil, fmt.Errorf("open diagnostic sqlite database: %w", err)
	}
	db.SetMaxOpenConns(1)
	db.SetMaxIdleConns(1)

	if err := db.PingContext(ctx); err != nil {
		return nil, closeDatabaseOnError(db, fmt.Errorf("ping diagnostic sqlite database: %w", err))
	}
	if _, err := db.ExecContext(ctx, createDiagnosticLogTableSQL); err != nil {
		return nil, closeDatabaseOnError(db, fmt.Errorf("create diagnostic log table: %w", err))
	}
	return &DiagnosticLog{db: db}, nil
}

// WithDiagnosticLogLevel sets the minimum slog level stored by the handler.
func WithDiagnosticLogLevel(level slog.Leveler) DiagnosticLogOption {
	return func(config *diagnosticLogConfig) {
		config.minLevel = level
	}
}

// WithDiagnosticLogSource sets the default diagnostic source value.
func WithDiagnosticLogSource(source string) DiagnosticLogOption {
	return func(config *diagnosticLogConfig) {
		config.source = strings.TrimSpace(source)
	}
}

// WithDiagnosticLogClock sets the fallback clock for records without a timestamp.
func WithDiagnosticLogClock(now func() time.Time) DiagnosticLogOption {
	return func(config *diagnosticLogConfig) {
		config.now = now
	}
}

// Handler returns a slog handler that writes structured diagnostic logs to SQLite.
func (log *DiagnosticLog) Handler(options ...DiagnosticLogOption) slog.Handler {
	config := diagnosticLogConfig{
		minLevel: slog.LevelInfo,
		source:   "backend",
		now:      func() time.Time { return time.Now().UTC() },
	}
	for _, option := range options {
		if option != nil {
			option(&config)
		}
	}
	if config.source == "" {
		config.source = "backend"
	}
	if config.now == nil {
		config.now = func() time.Time { return time.Now().UTC() }
	}
	return &sqliteDiagnosticHandler{
		db:       log.db,
		minLevel: config.minLevel,
		source:   config.source,
		now:      config.now,
		mu:       &sync.Mutex{},
	}
}

// Close closes the diagnostic SQLite database.
func (log *DiagnosticLog) Close() error {
	if log == nil || log.db == nil {
		return nil
	}
	if err := log.db.Close(); err != nil {
		return fmt.Errorf("close diagnostic sqlite database: %w", err)
	}
	return nil
}

type sqliteDiagnosticHandler struct {
	db       *sql.DB
	minLevel slog.Leveler
	source   string
	now      func() time.Time
	mu       *sync.Mutex
	attrs    []slog.Attr
	groups   []string
}

func (handler *sqliteDiagnosticHandler) Enabled(_ context.Context, level slog.Level) bool {
	if handler.minLevel == nil {
		return true
	}
	return level >= handler.minLevel.Level()
}

func (handler *sqliteDiagnosticHandler) Handle(ctx context.Context, record slog.Record) error {
	if !handler.Enabled(ctx, record.Level) {
		return nil
	}

	attrs := make([]slog.Attr, 0, len(handler.attrs)+record.NumAttrs())
	attrs = append(attrs, handler.attrs...)
	record.Attrs(func(attr slog.Attr) bool {
		attrs = append(attrs, handler.qualifyAttr(attr))
		return true
	})
	fields := diagnosticFieldsFromAttrs(attrs)
	attrsJSON, err := json.Marshal(fields)
	if err != nil {
		return fmt.Errorf("marshal diagnostic log attrs: %w", err)
	}

	source := firstNonEmpty(fields["source"], handler.source)
	eventName := firstNonEmpty(fields["event_name"], record.Message)
	occurredAt := record.Time
	if occurredAt.IsZero() && handler.now != nil {
		occurredAt = handler.now()
	}
	if occurredAt.IsZero() {
		occurredAt = time.Now().UTC()
	}

	handler.mu.Lock()
	defer handler.mu.Unlock()
	_, err = handler.db.ExecContext(
		ctx,
		insertDiagnosticLogSQL,
		occurredAt.UTC().Format(time.RFC3339Nano),
		strings.ToLower(record.Level.String()),
		source,
		nullableString(fields["trace_id"]),
		nullableString(fields["screen"]),
		nullableString(fields["boundary"]),
		eventName,
		record.Message,
		string(attrsJSON),
	)
	if err != nil {
		return fmt.Errorf("insert diagnostic log row: %w", err)
	}
	return nil
}

func (handler *sqliteDiagnosticHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	next := handler.clone()
	next.attrs = append(next.attrs, qualifyAttrs(handler.groups, attrs)...)
	return next
}

func (handler *sqliteDiagnosticHandler) WithGroup(name string) slog.Handler {
	if strings.TrimSpace(name) == "" {
		return handler
	}
	next := handler.clone()
	next.groups = append(next.groups, name)
	return next
}

func (handler *sqliteDiagnosticHandler) clone() *sqliteDiagnosticHandler {
	next := *handler
	next.attrs = append([]slog.Attr(nil), handler.attrs...)
	next.groups = append([]string(nil), handler.groups...)
	return &next
}

func (handler *sqliteDiagnosticHandler) qualifyAttr(attr slog.Attr) slog.Attr {
	attrs := qualifyAttrs(handler.groups, []slog.Attr{attr})
	if len(attrs) == 0 {
		return slog.Attr{}
	}
	return attrs[0]
}

func qualifyAttrs(groups []string, attrs []slog.Attr) []slog.Attr {
	if len(groups) == 0 {
		return append([]slog.Attr(nil), attrs...)
	}
	qualified := make([]slog.Attr, 0, len(attrs))
	prefix := strings.Join(groups, ".")
	for _, attr := range attrs {
		if attr.Key == "" {
			continue
		}
		attr.Key = prefix + "." + attr.Key
		qualified = append(qualified, attr)
	}
	return qualified
}

func diagnosticFieldsFromAttrs(attrs []slog.Attr) map[string]string {
	fields := make(map[string]string, len(attrs))
	for _, attr := range attrs {
		if attr.Key == "" {
			continue
		}
		fields[attr.Key] = diagnosticAttrValue(attr.Value)
	}
	return fields
}

func diagnosticAttrValue(value slog.Value) string {
	switch value.Kind() {
	case slog.KindString:
		return value.String()
	case slog.KindBool:
		if value.Bool() {
			return "true"
		}
		return "false"
	case slog.KindDuration:
		return value.Duration().String()
	case slog.KindFloat64:
		return fmt.Sprintf("%g", value.Float64())
	case slog.KindInt64:
		return fmt.Sprintf("%d", value.Int64())
	case slog.KindTime:
		return value.Time().UTC().Format(time.RFC3339Nano)
	case slog.KindUint64:
		return fmt.Sprintf("%d", value.Uint64())
	default:
		if value.Any() == nil {
			return ""
		}
		return fmt.Sprint(value.Any())
	}
}

func firstNonEmpty(values ...string) string {
	for _, value := range values {
		if strings.TrimSpace(value) != "" {
			return strings.TrimSpace(value)
		}
	}
	return ""
}

func nullableString(value string) any {
	if strings.TrimSpace(value) == "" {
		return nil
	}
	return strings.TrimSpace(value)
}

func resolveDiagnosticLogPath(databasePath string) (string, error) {
	path := strings.TrimSpace(databasePath)
	if path == "" {
		path = DefaultDiagnosticLogPath
	}
	resolvedPath, err := filepath.Abs(path)
	if err != nil {
		return "", fmt.Errorf("resolve diagnostic sqlite database path: %w", err)
	}
	if err := os.MkdirAll(filepath.Dir(resolvedPath), 0o750); err != nil {
		return "", fmt.Errorf("create diagnostic sqlite directory: %w", err)
	}
	return resolvedPath, nil
}

func diagnosticSQLiteDSN(databasePath string) string {
	query := url.Values{}
	query.Add("_pragma", "journal_mode(WAL)")
	query.Add("_pragma", "busy_timeout(5000)")
	query.Set("_time_format", "sqlite")
	query.Set("_txlock", "immediate")
	query.Set("_timezone", "UTC")
	return (&url.URL{
		Scheme:   "file",
		Path:     databasePath,
		RawQuery: query.Encode(),
	}).String()
}

func closeDatabaseOnError(db *sql.DB, cause error) error {
	if closeErr := db.Close(); closeErr != nil {
		return errors.Join(cause, closeErr)
	}
	return cause
}
