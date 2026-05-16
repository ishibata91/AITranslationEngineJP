package apitest

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	controllerwails "aitranslationenginejp/internal/controller/wails"
	infraruntime "aitranslationenginejp/internal/infra/runtime"
	"aitranslationenginejp/internal/notification"
	"aitranslationenginejp/internal/repository"
	"aitranslationenginejp/internal/service"
	"aitranslationenginejp/internal/usecase"
)

const (
	nmdImportProgressEventName  = "master-dictionary:import-progress"
	nmdImportCompletedEventName = "master-dictionary:import-completed"
)

var _ notification.SinkPort = (*nmdRecordingNotificationSink)(nil)
var _ notification.Port = (*nmdRecordingPort)(nil)

func TestSCN_NMD_001_ExecutionSideUsesSinkPortWithoutControllerRuntimeRoute(t *testing.T) {
	source := readNMDSourceSet(t,
		"internal/usecase/master_dictionary_usecase.go",
		"internal/service/master_dictionary_import_service.go",
		"internal/controller/wails/master_dictionary_controller.go",
	)
	fixture := newNMDMasterDictionaryScenarioFixture(t, &nmdRecordingNotificationSink{})

	response, sink := fixture.importXMLThroughController(t)

	if response.Summary.ImportedCount != 1 || len(sink.facts) == 0 {
		t.Fatalf("SCN-NMD-001 expected command result and notification facts, got response=%#v facts=%#v", response, sink.facts)
	}
	assertNMDSourceExcludes(t, source["internal/usecase/master_dictionary_usecase.go"], []string{
		`"aitranslationenginejp/internal/infra/runtime"`,
		"Dispatcher",
		"RuntimeAdapter",
		"runtime.EventsEmit",
	})
	assertNMDSourceExcludes(t, source["internal/service/master_dictionary_import_service.go"], []string{
		`"aitranslationenginejp/internal/infra/runtime"`,
		"Dispatcher",
		"RuntimeAdapter",
		"runtime.EventsEmit",
	})
	assertNMDSourceExcludes(t, source["internal/controller/wails/master_dictionary_controller.go"], []string{
		`"aitranslationenginejp/internal/notification"`,
		"Dispatcher",
		"KindMasterDictionaryImportProgress",
	})
}

func TestSCN_NMD_002_MasterDictionaryCommandResponseStaysSeparateFromNotificationPayload(t *testing.T) {
	sink := &nmdRecordingNotificationSink{}
	fixture := newNMDMasterDictionaryScenarioFixture(t, sink)

	response, _ := fixture.importXMLThroughController(t)

	responseJSON := marshalNMDJSON(t, response)
	assertNMDTextExcludes(t, string(responseJSON), []string{
		nmdImportProgressEventName,
		nmdImportCompletedEventName,
		"eventName",
		"notification_failure",
		"transport_error",
	})
	if len(sink.facts) == 0 {
		t.Fatal("SCN-NMD-002 expected notification facts to be delivered outside the command response")
	}
	for _, fact := range sink.facts {
		if fact.Kind == "" {
			t.Fatalf("SCN-NMD-002 expected notification fact kind to be separate from DTO payload, got %#v", fact)
		}
	}
}

func TestSCN_NMD_003_MasterDictionaryNotificationFactsFollowPersistedImportState(t *testing.T) {
	sink := &nmdRecordingNotificationSink{}
	fixture := newNMDMasterDictionaryScenarioFixture(t, sink)

	response, _ := fixture.importXMLThroughController(t)
	persisted, err := fixture.repository.GetByID(context.Background(), response.Summary.LastEntryID)
	if err != nil {
		t.Fatalf("SCN-NMD-003 expected imported entry to be persisted: %v", err)
	}

	progress := sink.requireProgressFacts(t)
	completed := sink.requireCompletedFact(t)
	if progress[0] != 0 || progress[len(progress)-1] != 100 {
		t.Fatalf("SCN-NMD-003 expected progress facts to track confirmed processing, got %#v", progress)
	}
	if completed.Import.Summary.LastEntryID != persisted.ID || completed.Import.Page.SelectedID == nil || *completed.Import.Page.SelectedID != persisted.ID {
		t.Fatalf("SCN-NMD-003 expected completion fact to follow persisted state, fact=%#v persisted=%#v", completed.Import, persisted)
	}
}

func TestSCN_NMD_004_DispatcherRedactsPayloadAndRejectsUnsafePayload(t *testing.T) {
	capture := startObservabilityLogCapture(t)
	port := &nmdRecordingPort{}
	dispatcher := notification.NewDispatcher(port)

	safeResult := dispatcher.Dispatch(context.Background(), nmdCompletedFactWithSourcePathAndUnsafeReason())
	unsafeResult := dispatcher.Dispatch(context.Background(), nmdUnsafeCompletedFact())

	if !safeResult.Sent || len(port.sent) != 1 {
		t.Fatalf("SCN-NMD-004 expected redacted safe payload to be sent once, result=%#v sent=%#v", safeResult, port.sent)
	}
	payloadJSON := marshalNMDJSON(t, port.sent[0])
	assertNMDTextExcludes(t, string(payloadJSON), nmdForbiddenPayloadValues())
	if port.sent[0].Import.Summary.FilePath != "" || port.sent[0].Import.Summary.FileName != "source.xml" {
		t.Fatalf("SCN-NMD-004 expected file path redaction and basename payload, got %#v", port.sent[0].Import.Summary)
	}
	if !unsafeResult.Suppressed || !errors.Is(unsafeResult.Err, notification.ErrUnsafePayload) || len(port.sent) != 1 {
		t.Fatalf("SCN-NMD-004 expected unsafe payload rejection without transport send, result=%#v sent=%#v", unsafeResult, port.sent)
	}
	capture.requireEvent(t, "notification_dispatch", "rejected", "unsafe_payload")
	assertNMDTextExcludes(t, capture.buffer.String(), nmdForbiddenPayloadValues())
}

func TestSCN_NMD_005_NotificationSendFailureDoesNotRollbackCommandResponseOrDBState(t *testing.T) {
	capture := startObservabilityLogCapture(t)
	port := &nmdRecordingPort{err: errors.New("runtime send failed")}
	fixture := newNMDMasterDictionaryScenarioFixture(t, notification.NewDispatcher(port))

	response, _ := fixture.importXMLThroughController(t)
	persisted, err := fixture.repository.GetByID(context.Background(), response.Summary.LastEntryID)
	if err != nil {
		t.Fatalf("SCN-NMD-005 expected DB state to remain persisted after notification failure: %v", err)
	}

	if response.Summary.ImportedCount != 1 || persisted.Source != nmdImportedSource {
		t.Fatalf("SCN-NMD-005 expected command response and DB state to stay successful, response=%#v persisted=%#v", response, persisted)
	}
	capture.requireEvent(t, "notification_dispatch", "failed", "transport_error")
	assertNMDTextExcludes(t, capture.buffer.String(), nmdForbiddenPayloadValues())
}

func TestSCN_NMD_006_StateAndProviderJudgementStayOutsideNotificationModule(t *testing.T) {
	source := readNMDSourceSet(t, "internal/notification/notification.go")
	port := &nmdRecordingPort{}
	dispatcher := notification.NewDispatcher(port)

	first := dispatcher.Dispatch(context.Background(), nmdProgressFact(100))
	second := dispatcher.Dispatch(context.Background(), nmdProgressFact(100))

	if !first.Sent || !second.Sent || len(port.sent) != 2 {
		t.Fatalf("SCN-NMD-006 expected duplicate notification facts to remain transport facts only, first=%#v second=%#v sent=%#v", first, second, port.sent)
	}
	assertNMDSourceExcludes(t, source["internal/notification/notification.go"], []string{
		"TranslationJobPolicy",
		"terminal_job",
		"RecoverableFailed",
		"provider response",
		"invalid_provider_response",
		"late response",
		"phase complete",
	})
}

func TestSCN_NMD_007_RuntimeAdapterOwnsWailsEventMappingForPushOnlyNotification(t *testing.T) {
	source := readNMDSourceSet(t,
		"internal/usecase/master_dictionary_usecase.go",
		"internal/service/master_dictionary_import_service.go",
		"internal/notification/notification.go",
	)
	emitter := &nmdRuntimeEventEmitter{}
	adapter := infraruntime.NewNotificationAdapter(func() (context.Context, bool) {
		return nmdRuntimeEventContext{Context: context.Background(), emitter: emitter}, true
	})

	err := adapter.Send(context.Background(), notification.Notification{
		Kind:     notification.KindMasterDictionaryImportProgress,
		Progress: &notification.ProgressNotification{Percent: 42},
	})

	if err != nil {
		t.Fatalf("SCN-NMD-007 expected runtime adapter send to succeed: %v", err)
	}
	if len(emitter.events) != 1 || emitter.events[0].name != nmdImportProgressEventName {
		t.Fatalf("SCN-NMD-007 expected runtime adapter to map event name, got %#v", emitter.events)
	}
	payloadJSON := marshalNMDJSON(t, emitter.events[0].payload[0])
	if !strings.Contains(string(payloadJSON), `"progress":42`) {
		t.Fatalf("SCN-NMD-007 expected adapter transport payload to expose progress field, got %s", payloadJSON)
	}
	assertNMDSourceExcludes(t, source["internal/usecase/master_dictionary_usecase.go"], []string{nmdImportProgressEventName, "EventsEmit"})
	assertNMDSourceExcludes(t, source["internal/service/master_dictionary_import_service.go"], []string{nmdImportProgressEventName, "EventsEmit"})
	assertNMDSourceExcludes(t, source["internal/notification/notification.go"], []string{nmdImportProgressEventName, "EventsEmit"})
}

func TestSCN_NMD_008_NotificationResultIsNotPersistedAndLogsStayMinimal(t *testing.T) {
	capture := startObservabilityLogCapture(t)
	db := openNMDSchemaDatabase(t)
	port := &nmdRecordingPort{}
	dispatcher := notification.NewDispatcher(port)

	dispatcher.Dispatch(context.Background(), nmdProgressFact(10))
	dispatcher.Dispatch(context.Background(), notification.Fact{Kind: notification.KindMasterDictionaryImportCompleted})
	dispatcher.Dispatch(context.Background(), nmdUnsafeCompletedFact())
	dispatcher.Dispatch(context.Background(), nmdCompletedFactWithSourcePathAndUnsafeReason())

	assertNMDNotificationTablesAbsent(t, db)
	for _, expected := range []struct {
		result string
		reason string
	}{
		{result: "sent"},
		{result: "skipped", reason: "not_sendable"},
		{result: "rejected", reason: "unsafe_payload"},
	} {
		payload := capture.requireEvent(t, "notification_dispatch", expected.result, expected.reason)
		assertNMDLogHasOnlyMinimalCustomKeys(t, payload)
	}
	assertNMDTextExcludes(t, capture.buffer.String(), nmdForbiddenPayloadValues())
}

type nmdMasterDictionaryScenarioFixture struct {
	controller *controllerwails.MasterDictionaryController
	repository service.RepositoryPort
	xmlPath    string
	sink       *nmdRecordingNotificationSink
}

func newNMDMasterDictionaryScenarioFixture(
	t *testing.T,
	sink notification.SinkPort,
) nmdMasterDictionaryScenarioFixture {
	t.Helper()

	dbPath := filepath.Join(t.TempDir(), "nmd-master-dictionary.sqlite3")
	repositoryPort, err := service.NewSQLiteMasterDictionaryRepositoryPort(context.Background(), dbPath, nil)
	if err != nil {
		t.Fatalf("build sqlite master dictionary repository: %v", err)
	}
	t.Cleanup(func() {
		closeFn := service.SQLiteMasterDictionaryRepositoryPortCloser(repositoryPort)
		if closeFn != nil {
			_ = closeFn(context.Background())
		}
	})

	importService := service.NewMasterDictionaryImportService(
		repositoryPort,
		service.NewLocalMasterDictionaryXMLFilePort(),
		service.NewXMLDecoderMasterDictionaryRecordReader(),
		sink,
		nmdNow,
	)
	queryService := service.NewMasterDictionaryQueryService(repositoryPort)
	commandService := service.NewMasterDictionaryCommandService(repositoryPort, nmdNow)
	masterDictionaryUsecase := usecase.NewMasterDictionaryUsecase(queryService, commandService, importService, sink)
	controller := controllerwails.NewMasterDictionaryController(masterDictionaryUsecase, controllerwails.NewRuntimeEmitterState())

	recordingSink, _ := sink.(*nmdRecordingNotificationSink)
	return nmdMasterDictionaryScenarioFixture{
		controller: controller,
		repository: repositoryPort,
		xmlPath:    writeNMDImportXMLFixture(t),
		sink:       recordingSink,
	}
}

func (fixture nmdMasterDictionaryScenarioFixture) importXMLThroughController(
	t *testing.T,
) (controllerwails.MasterDictionaryImportResponseDTO, *nmdRecordingNotificationSink) {
	t.Helper()

	response, err := fixture.controller.MasterDictionaryImportXML(controllerwails.MasterDictionaryImportRequestDTO{
		XMLPath: fixture.xmlPath,
		Refresh: controllerwails.MasterDictionaryRefreshQueryDTO{
			Category: "すべて",
			Page:     1,
			PageSize: 30,
		},
	})
	if err != nil {
		t.Fatalf("import master dictionary XML through controller: %v", err)
	}
	return response, fixture.sink
}

type nmdRecordingNotificationSink struct {
	facts []notification.Fact
}

func (sink *nmdRecordingNotificationSink) Notify(_ context.Context, fact notification.Fact) {
	sink.facts = append(sink.facts, fact)
}

func (sink *nmdRecordingNotificationSink) requireProgressFacts(t *testing.T) []int {
	t.Helper()

	progress := []int{}
	for _, fact := range sink.facts {
		if fact.Kind != notification.KindMasterDictionaryImportProgress || fact.Progress == nil {
			continue
		}
		progress = append(progress, fact.Progress.Percent)
	}
	if len(progress) == 0 {
		t.Fatalf("expected progress facts, got %#v", sink.facts)
	}
	return progress
}

func (sink *nmdRecordingNotificationSink) requireCompletedFact(t *testing.T) notification.Fact {
	t.Helper()

	for _, fact := range sink.facts {
		if fact.Kind == notification.KindMasterDictionaryImportCompleted && fact.Import != nil {
			return fact
		}
	}
	t.Fatalf("expected completed fact, got %#v", sink.facts)
	return notification.Fact{}
}

type nmdRecordingPort struct {
	sent []notification.Notification
	err  error
}

func (port *nmdRecordingPort) Send(_ context.Context, event notification.Notification) error {
	if port.err != nil {
		return port.err
	}
	port.sent = append(port.sent, event)
	return nil
}

type nmdRuntimeEvent struct {
	name    string
	payload []interface{}
}

type nmdRuntimeEventEmitter struct {
	events []nmdRuntimeEvent
}

func (emitter *nmdRuntimeEventEmitter) Emit(eventName string, optionalData ...interface{}) {
	emitter.events = append(emitter.events, nmdRuntimeEvent{name: eventName, payload: optionalData})
}

type nmdRuntimeEventContext struct {
	context.Context
	emitter *nmdRuntimeEventEmitter
}

func (ctx nmdRuntimeEventContext) Value(key interface{}) interface{} {
	if key == "events" {
		return ctx.emitter
	}
	return ctx.Context.Value(key)
}

const (
	nmdImportedSource      = "NMD Silver Sword"
	nmdImportedTranslation = "NMD 銀の剣"
)

func writeNMDImportXMLFixture(t *testing.T) string {
	t.Helper()

	xmlPath := filepath.Join(t.TempDir(), "nmd-source.xml")
	xmlContent := fmt.Sprintf(`<?xml version="1.0" encoding="utf-8"?>
<Root>
	<String>
		<REC>WEAP:FULL</REC>
		<EDID>NMD_SilverSword</EDID>
		<Source>%s</Source>
		<Dest>%s</Dest>
	</String>
</Root>`, nmdImportedSource, nmdImportedTranslation)
	if err := os.WriteFile(xmlPath, []byte(xmlContent), 0o600); err != nil {
		t.Fatalf("write NMD XML fixture: %v", err)
	}
	return xmlPath
}

func nmdProgressFact(percent int) notification.Fact {
	return notification.Fact{
		Kind:     notification.KindMasterDictionaryImportProgress,
		Progress: &notification.ProgressFact{Percent: percent},
	}
}

func nmdCompletedFactWithSourcePathAndUnsafeReason() notification.Fact {
	return notification.Fact{
		Kind: notification.KindMasterDictionaryImportCompleted,
		Import: &notification.MasterDictionaryImportFact{
			Page: notification.MasterDictionaryPage{
				Items: []notification.MasterDictionaryEntry{{
					ID:          1001,
					Source:      "Safe source",
					Translation: "安全な訳文",
					Category:    "武器",
					Origin:      "XML取込",
					REC:         "WEAP:FULL",
					EDID:        "SafeSword",
					UpdatedAt:   nmdNow(),
				}},
				TotalCount: 1,
				Page:       1,
				PageSize:   30,
			},
			Summary: notification.MasterDictionaryImportSummary{
				FilePath:      "/tmp/source.xml",
				FileName:      "/tmp/source.xml",
				ImportedCount: 1,
				LastEntryID:   1001,
			},
			Reason: "provider raw response secret",
		},
	}
}

func nmdUnsafeCompletedFact() notification.Fact {
	return notification.Fact{
		Kind: notification.KindMasterDictionaryImportCompleted,
		Import: &notification.MasterDictionaryImportFact{
			Summary: notification.MasterDictionaryImportSummary{
				FilePath: "/tmp/secret-api-key-provider-raw.xml",
			},
		},
	}
}

func nmdForbiddenPayloadValues() []string {
	return []string{
		"secret-api-key",
		"provider raw",
		"provider_raw_request",
		"provider_raw_response",
		"credential/source.xml",
		"<Strings>",
		"<String>",
		"full prompt text",
	}
}

func marshalNMDJSON(t *testing.T, value interface{}) []byte {
	t.Helper()

	payload, err := json.Marshal(value)
	if err != nil {
		t.Fatalf("marshal JSON: %v", err)
	}
	return payload
}

func assertNMDTextExcludes(t *testing.T, text string, forbiddenValues []string) {
	t.Helper()

	lowerText := strings.ToLower(text)
	for _, forbidden := range forbiddenValues {
		if strings.Contains(lowerText, strings.ToLower(forbidden)) {
			t.Fatalf("expected text to exclude forbidden value %q, got %s", forbidden, text)
		}
	}
}

func readNMDSourceSet(t *testing.T, paths ...string) map[string]string {
	t.Helper()

	root := nmdRepositoryRoot(t)
	sources := map[string]string{}
	for _, path := range paths {
		content, err := os.ReadFile(nmdKnownSourcePath(t, root, path))
		if err != nil {
			t.Fatalf("read source %s: %v", path, err)
		}
		sources[path] = string(content)
	}
	return sources
}

func nmdKnownSourcePath(t *testing.T, root string, path string) string {
	t.Helper()

	switch path {
	case "internal/usecase/master_dictionary_usecase.go":
		return filepath.Join(root, "internal", "usecase", "master_dictionary_usecase.go")
	case "internal/service/master_dictionary_import_service.go":
		return filepath.Join(root, "internal", "service", "master_dictionary_import_service.go")
	case "internal/controller/wails/master_dictionary_controller.go":
		return filepath.Join(root, "internal", "controller", "wails", "master_dictionary_controller.go")
	case "internal/notification/notification.go":
		return filepath.Join(root, "internal", "notification", "notification.go")
	default:
		t.Fatalf("unknown NMD source fixture path %q", path)
		return ""
	}
}

func assertNMDSourceExcludes(t *testing.T, source string, forbiddenSnippets []string) {
	t.Helper()

	for _, forbidden := range forbiddenSnippets {
		if strings.Contains(source, forbidden) {
			t.Fatalf("expected source to exclude %q", forbidden)
		}
	}
}

func nmdRepositoryRoot(t *testing.T) string {
	t.Helper()

	dir, err := os.Getwd()
	if err != nil {
		t.Fatalf("get working directory: %v", err)
	}
	for {
		if _, err := os.Stat(filepath.Join(dir, "go.mod")); err == nil {
			return dir
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			t.Fatal("repository root with go.mod was not found")
		}
		dir = parent
	}
}

func openNMDSchemaDatabase(t *testing.T) *sql.DB {
	t.Helper()

	dbPath := filepath.Join(t.TempDir(), "nmd-schema.sqlite3")
	db, err := repository.OpenSQLiteDatabase(context.Background(), dbPath)
	if err != nil {
		t.Fatalf("open sqlite database: %v", err)
	}
	t.Cleanup(func() { _ = db.Close() })
	return db.DB
}

func assertNMDNotificationTablesAbsent(t *testing.T, db *sql.DB) {
	t.Helper()

	rows, err := db.QueryContext(context.Background(), "SELECT name FROM sqlite_master WHERE type='table'")
	if err != nil {
		t.Fatalf("query sqlite tables: %v", err)
	}
	defer func() {
		if err := rows.Close(); err != nil {
			t.Fatalf("close sqlite table rows: %v", err)
		}
	}()

	for rows.Next() {
		var tableName string
		if err := rows.Scan(&tableName); err != nil {
			t.Fatalf("scan sqlite table name: %v", err)
		}
		lowerName := strings.ToLower(tableName)
		if strings.Contains(lowerName, "notification") || strings.Contains(lowerName, "runtime_event") || strings.Contains(lowerName, "operation_summary") {
			t.Fatalf("expected notification result tables to be absent, found %q", tableName)
		}
	}
	if err := rows.Err(); err != nil {
		t.Fatalf("iterate sqlite tables: %v", err)
	}
}

func assertNMDLogHasOnlyMinimalCustomKeys(t *testing.T, payload map[string]any) {
	t.Helper()

	allowed := map[string]struct{}{
		"time":   {},
		"level":  {},
		"msg":    {},
		"event":  {},
		"where":  {},
		"result": {},
		"id":     {},
		"count":  {},
		"reason": {},
	}
	for key := range payload {
		if _, ok := allowed[key]; !ok {
			t.Fatalf("expected notification log to use only minimal keys, got key=%q payload=%#v", key, payload)
		}
	}
}

func nmdNow() time.Time {
	return time.Date(2026, 5, 16, 10, 0, 0, 0, time.UTC)
}
