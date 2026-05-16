// Package notification accepts execution facts and dispatches redacted
// transport-independent notifications.
package notification

import (
	"context"
	"errors"
	"log/slog"
	"strings"
	"time"
)

// ErrUnsafePayload means a notification fact could not be reduced to a safe payload.
var ErrUnsafePayload = errors.New("notification unsafe payload")

// Kind identifies the notification fact type.
type Kind string

const (
	// KindMasterDictionaryImportProgress reports master dictionary import progress.
	KindMasterDictionaryImportProgress Kind = "master_dictionary_import_progress"
	// KindMasterDictionaryImportCompleted reports master dictionary import completion.
	KindMasterDictionaryImportCompleted Kind = "master_dictionary_import_completed"
	// KindMasterDictionaryImportDiscarded reports a discarded master dictionary import notification.
	KindMasterDictionaryImportDiscarded Kind = "master_dictionary_import_discarded"
)

// SinkPort is the execution-side entrypoint for notification facts.
type SinkPort interface {
	Notify(ctx context.Context, fact Fact)
}

// Port sends redacted notifications to a transport adapter.
type Port interface {
	Send(ctx context.Context, notification Notification) error
}

// Fact carries an execution fact before dispatch policy and redaction are applied.
type Fact struct {
	Kind     Kind
	Progress *ProgressFact
	Import   *MasterDictionaryImportFact
}

// ProgressFact carries progress that has already been confirmed by the execution side.
type ProgressFact struct {
	Percent int
}

// MasterDictionaryImportFact carries import completion or discard facts.
type MasterDictionaryImportFact struct {
	Page    MasterDictionaryPage
	Summary MasterDictionaryImportSummary
	Refresh MasterDictionaryImportRefresh
	Reason  string
}

// MasterDictionaryPage carries page state for import completion notifications.
type MasterDictionaryPage struct {
	Items      []MasterDictionaryEntry
	TotalCount int
	Page       int
	PageSize   int
	SelectedID *int64
}

// MasterDictionaryEntry carries redaction-safe dictionary entry fields.
type MasterDictionaryEntry struct {
	ID          int64
	Source      string
	Translation string
	Category    string
	Origin      string
	REC         string
	EDID        string
	UpdatedAt   time.Time
}

// MasterDictionaryImportSummary carries redaction-safe import summary fields.
type MasterDictionaryImportSummary struct {
	FilePath      string
	FileName      string
	ImportedCount int
	UpdatedCount  int
	SkippedCount  int
	SelectedREC   []string
	LastEntryID   int64
}

// MasterDictionaryImportRefresh carries approved refresh semantics for completion notifications.
type MasterDictionaryImportRefresh struct {
	Query           string
	Category        string
	Page            int
	PageSize        int
	RefreshTargetID *int64
}

// Notification is a redacted transport-independent notification.
type Notification struct {
	Kind     Kind
	Progress *ProgressNotification
	Import   *MasterDictionaryImportNotification
}

// ProgressNotification carries bounded progress for transport.
type ProgressNotification struct {
	Percent int
}

// MasterDictionaryImportNotification carries redacted import facts for transport.
type MasterDictionaryImportNotification struct {
	Page    MasterDictionaryPage
	Summary MasterDictionaryImportSummary
	Refresh MasterDictionaryImportRefresh
	Reason  string
}

// DispatchResult describes the local outcome of a notification dispatch.
type DispatchResult struct {
	Sent       bool
	Suppressed bool
	Err        error
}

// Dispatcher applies notification kind, redaction, sendability, and send failure policy.
type Dispatcher struct {
	port Port
}

// NewDispatcher creates a notification dispatcher.
func NewDispatcher(port Port) *Dispatcher {
	return &Dispatcher{port: port}
}

// Notify dispatches a notification fact and isolates send failures from application results.
func (dispatcher *Dispatcher) Notify(ctx context.Context, fact Fact) {
	_ = dispatcher.Dispatch(ctx, fact)
}

// Dispatch applies redaction and sendability policy before transport send.
func (dispatcher *Dispatcher) Dispatch(ctx context.Context, fact Fact) DispatchResult {
	redacted, ok, err := Redact(fact)
	if err != nil {
		result := DispatchResult{Suppressed: true, Err: err}
		logDispatchResult(ctx, fact.Kind, result, "unsafe_payload")
		return result
	}
	if !ok || dispatcher.port == nil {
		result := DispatchResult{Suppressed: true}
		logDispatchResult(ctx, fact.Kind, result, dispatchSuppressedReason(ok, dispatcher.port))
		return result
	}
	if err := dispatcher.port.Send(ctx, redacted); err != nil {
		result := DispatchResult{Err: err}
		logDispatchResult(ctx, fact.Kind, result, "transport_error")
		return result
	}
	result := DispatchResult{Sent: true}
	logDispatchResult(ctx, fact.Kind, result, "")
	return result
}

func dispatchSuppressedReason(sendable bool, port Port) string {
	if !sendable {
		return "not_sendable"
	}
	if port == nil {
		return "transport_unavailable"
	}
	return "unknown"
}

func logDispatchResult(ctx context.Context, kind Kind, result DispatchResult, reason string) {
	attrs := []slog.Attr{
		slog.String("event", "notification_dispatch"),
		slog.String("where", "backend.notification.dispatcher"),
		slog.String("result", dispatchLogResult(result)),
		slog.String("id", string(kind)),
	}
	if reason != "" {
		attrs = append(attrs, slog.String("reason", reason))
	}
	if result.Err != nil {
		slog.LogAttrs(ctx, slog.LevelWarn, "notification dispatch completed", attrs...)
		return
	}
	slog.LogAttrs(ctx, slog.LevelInfo, "notification dispatch completed", attrs...)
}

func dispatchLogResult(result DispatchResult) string {
	if result.Sent {
		return "sent"
	}
	if result.Suppressed && result.Err != nil {
		return "rejected"
	}
	if result.Suppressed {
		return "skipped"
	}
	if result.Err != nil {
		return "failed"
	}
	return "unknown"
}

// Redact converts an execution fact to a transport-independent safe notification.
func Redact(fact Fact) (Notification, bool, error) {
	switch fact.Kind {
	case KindMasterDictionaryImportProgress:
		if fact.Progress == nil {
			return Notification{}, false, nil
		}
		return Notification{
			Kind:     fact.Kind,
			Progress: &ProgressNotification{Percent: clampPercent(fact.Progress.Percent)},
		}, true, nil
	case KindMasterDictionaryImportCompleted:
		if fact.Import == nil {
			return Notification{}, false, nil
		}
		importFact, err := redactImportFact(*fact.Import)
		if err != nil {
			return Notification{}, false, err
		}
		return Notification{Kind: fact.Kind, Import: &importFact}, true, nil
	case KindMasterDictionaryImportDiscarded:
		if fact.Import == nil {
			return Notification{}, false, nil
		}
		importFact := MasterDictionaryImportNotification{
			Reason: redactReason(fact.Import.Reason),
		}
		return Notification{Kind: fact.Kind, Import: &importFact}, true, nil
	default:
		return Notification{}, false, nil
	}
}

func redactImportFact(fact MasterDictionaryImportFact) (MasterDictionaryImportNotification, error) {
	if containsUnsafeText(fact.Summary.FilePath) {
		return MasterDictionaryImportNotification{}, ErrUnsafePayload
	}
	sourcePath := fact.Summary.FilePath
	fact.Summary.FilePath = ""
	fact.Summary.FileName = baseName(fact.Summary.FileName)
	if fact.Summary.FileName == "" {
		fact.Summary.FileName = baseName(sourcePath)
	}
	return MasterDictionaryImportNotification{
		Page:    fact.Page,
		Summary: fact.Summary,
		Refresh: fact.Refresh,
		Reason:  redactReason(fact.Reason),
	}, nil
}

func clampPercent(progress int) int {
	if progress < 0 {
		return 0
	}
	if progress > 100 {
		return 100
	}
	return progress
}

func redactReason(reason string) string {
	reason = strings.TrimSpace(reason)
	if containsUnsafeText(reason) {
		return "redacted"
	}
	return reason
}

func containsUnsafeText(value string) bool {
	lower := strings.ToLower(value)
	for _, marker := range []string{"api key", "apikey", "secret", "credential", "provider raw", "<strings>", "<string>"} {
		if strings.Contains(lower, marker) {
			return true
		}
	}
	return false
}

func baseName(path string) string {
	path = strings.TrimSpace(path)
	path = strings.ReplaceAll(path, "\\", "/")
	if index := strings.LastIndex(path, "/"); index >= 0 {
		return path[index+1:]
	}
	return path
}
