package bootstrap

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	controllerwails "aitranslationenginejp/internal/controller/wails"
)

type recordedRuntimeEvent struct {
	name    string
	payload []interface{}
}

type recordingRuntimeEventEmitter struct {
	events []recordedRuntimeEvent
}

type runtimeEventsTestContext struct {
	context.Context
	emitter *recordingRuntimeEventEmitter
}

const (
	bootstrapCreatedSource        = "Auriel's Bow"
	bootstrapCreatedTranslation   = "アーリエルの弓"
	bootstrapUpdatedTranslation   = "更新済みアーリエルの弓"
	bootstrapPersistedSource      = "Castle Volkihar"
	bootstrapImportProgressEvent  = "master-dictionary:import-progress"
	bootstrapImportCompletedEvent = "master-dictionary:import-completed"
	bootstrapPageQuerySucceeded   = "expected production graph page query to succeed: %v"
)

func (emitter *recordingRuntimeEventEmitter) Emit(eventName string, optionalData ...interface{}) {
	emitter.events = append(emitter.events, recordedRuntimeEvent{name: eventName, payload: optionalData})
}

func (ctx runtimeEventsTestContext) Value(key interface{}) interface{} {
	if key == "events" {
		return ctx.emitter
	}
	return ctx.Context.Value(key)
}

func TestNewAppControllerCreatesSQLiteDatabaseFile(t *testing.T) {
	databasePath := configureBootstrapTestDatabase(t)
	controller := NewAppController()

	_, err := controller.MasterDictionaryGetPage(controllerwails.MasterDictionaryPageRequestDTO{
		Refresh: controllerwails.MasterDictionaryRefreshQueryDTO{SearchTerm: "Whiterun", Category: "地名", Page: 1, PageSize: 5},
	})
	if err != nil {
		t.Fatalf(bootstrapPageQuerySucceeded, err)
	}

	_, err = os.Stat(databasePath)
	if err != nil {
		t.Fatalf("expected sqlite database file to exist: %v", err)
	}
}

func TestNewAppControllerReturnsSeededPageItems(t *testing.T) {
	configureBootstrapTestDatabase(t)
	controller := NewAppController()

	page, err := controller.MasterDictionaryGetPage(controllerwails.MasterDictionaryPageRequestDTO{
		Refresh: controllerwails.MasterDictionaryRefreshQueryDTO{SearchTerm: "Whiterun", Category: "地名", Page: 1, PageSize: 5},
	})
	if err != nil {
		t.Fatalf(bootstrapPageQuerySucceeded, err)
	}

	if page.Page.TotalCount == 0 || len(page.Page.Items) == 0 {
		t.Fatalf("expected seeded page items, got %#v", page.Page)
	}
}

func TestNewAppControllerSelectsFirstPageItemByDefault(t *testing.T) {
	configureBootstrapTestDatabase(t)
	controller := NewAppController()

	page, err := controller.MasterDictionaryGetPage(controllerwails.MasterDictionaryPageRequestDTO{
		Refresh: controllerwails.MasterDictionaryRefreshQueryDTO{SearchTerm: "Whiterun", Category: "地名", Page: 1, PageSize: 5},
	})
	if err != nil {
		t.Fatalf(bootstrapPageQuerySucceeded, err)
	}

	if page.Page.SelectedID == nil || *page.Page.SelectedID != page.Page.Items[0].ID {
		t.Fatalf("expected selected id to match first page item, got selected=%#v items=%#v", page.Page.SelectedID, page.Page.Items)
	}
}

func TestNewAppControllerReturnsDetailForSeededEntry(t *testing.T) {
	configureBootstrapTestDatabase(t)
	controller := NewAppController()
	page, err := controller.MasterDictionaryGetPage(controllerwails.MasterDictionaryPageRequestDTO{
		Refresh: controllerwails.MasterDictionaryRefreshQueryDTO{SearchTerm: "Whiterun", Category: "地名", Page: 1, PageSize: 5},
	})
	if err != nil {
		t.Fatalf(bootstrapPageQuerySucceeded, err)
	}

	detail, err := controller.MasterDictionaryGetDetail(controllerwails.MasterDictionaryDetailRequestDTO{ID: page.Page.Items[0].ID})
	if err != nil {
		t.Fatalf("expected detail query to succeed: %v", err)
	}

	if detail.Entry.Source == "" || detail.Entry.Translation == "" {
		t.Fatalf("expected populated detail entry, got %#v", detail.Entry)
	}
}

func TestNewAppControllerStartupDoesNotEmitRuntimeEvents(t *testing.T) {
	configureBootstrapTestDatabase(t)
	controller := NewAppController()
	emitter := &recordingRuntimeEventEmitter{}

	controller.OnStartup(runtimeEventsTestContext{Context: context.Background(), emitter: emitter})

	if len(emitter.events) != 0 {
		t.Fatalf("expected startup not to emit events by itself, got %#v", emitter.events)
	}
}

func TestNewAppControllerImportXMLReturnsImportSummary(t *testing.T) {
	result, _, controller := runBootstrapImport(t)

	if result.Summary.ImportedCount != 1 || result.Summary.LastEntryID == 0 {
		t.Fatalf("expected imported entry summary, got %#v", result.Summary)
	}

	controller.OnShutdown(context.Background())
}

func TestNewAppControllerImportXMLSelectsImportedEntry(t *testing.T) {
	result, _, controller := runBootstrapImport(t)

	if result.Page.SelectedID == nil || *result.Page.SelectedID != result.Summary.LastEntryID {
		t.Fatalf("expected import refresh to select last entry id, got page=%#v summary=%#v", result.Page, result.Summary)
	}

	controller.OnShutdown(context.Background())
}

func TestNewAppControllerImportXMLMakesImportedDetailQueryable(t *testing.T) {
	result, _, controller := runBootstrapImport(t)

	detail, err := controller.MasterDictionaryGetDetail(controllerwails.MasterDictionaryDetailRequestDTO{ID: result.Summary.LastEntryID})
	if err != nil {
		t.Fatalf("expected imported detail lookup to succeed: %v", err)
	}

	if detail.Entry.Source != "Auriel's Shield" || detail.Entry.Translation != "アーリエルの盾" {
		t.Fatalf("expected imported entry to be queryable, got %#v", detail.Entry)
	}

	controller.OnShutdown(context.Background())
}

func TestNewAppControllerImportXMLPublishesRuntimeEvents(t *testing.T) {
	_, emitter, controller := runBootstrapImport(t)

	assertEventNames(t, emitter.events,
		bootstrapImportProgressEvent,
		bootstrapImportProgressEvent,
		bootstrapImportProgressEvent,
		bootstrapImportCompletedEvent,
	)

	controller.OnShutdown(context.Background())
}

func TestNewAppControllerCreateMutationReturnsCreatedEntry(t *testing.T) {
	created, controller := runBootstrapCreate(t)

	if created.Entry.Source != bootstrapCreatedSource || created.Entry.Translation != bootstrapCreatedTranslation {
		t.Fatalf("unexpected created entry: %#v", created.Entry)
	}

	controller.OnShutdown(context.Background())
}

func TestNewAppControllerCreateMutationReturnsSelectedPage(t *testing.T) {
	created, controller := runBootstrapCreate(t)

	if created.Page == nil || created.Page.SelectedID == nil {
		t.Fatalf("expected create to return selected page, got %#v", created.Page)
	}

	controller.OnShutdown(context.Background())
}

func TestNewAppControllerCreateMutationPersistsDetail(t *testing.T) {
	created, controller := runBootstrapCreate(t)

	createdDetail, err := controller.GetMasterDictionaryEntry(controllerwails.GetMasterDictionaryEntryRequestDTO{ID: created.RefreshTargetID})
	if err != nil {
		t.Fatalf("expected created entry detail lookup to succeed: %v", err)
	}

	if createdDetail.Entry == nil || createdDetail.Entry.Translation != bootstrapCreatedTranslation {
		t.Fatalf("unexpected created detail payload: %#v", createdDetail.Entry)
	}

	controller.OnShutdown(context.Background())
}

func TestNewAppControllerUpdateMutationReturnsUpdatedEntry(t *testing.T) {
	created, controller := runBootstrapCreate(t)
	updated := runBootstrapUpdate(t, controller, created.RefreshTargetID)

	if updated.Entry.Translation != bootstrapUpdatedTranslation {
		t.Fatalf("unexpected updated entry payload: %#v", updated.Entry)
	}

	controller.OnShutdown(context.Background())
}

func TestNewAppControllerDeleteMutationReturnsDeletedID(t *testing.T) {
	created, controller := runBootstrapCreate(t)
	deleted := runBootstrapDelete(t, controller, created.RefreshTargetID)

	if deleted.DeletedID != created.RefreshTargetID {
		t.Fatalf("expected deleted id %q, got %#v", created.RefreshTargetID, deleted)
	}

	controller.OnShutdown(context.Background())
}

func TestNewAppControllerDeleteMutationMapsMissingDetailToNil(t *testing.T) {
	created, controller := runBootstrapCreate(t)
	runBootstrapDelete(t, controller, created.RefreshTargetID)

	deletedDetail, err := controller.GetMasterDictionaryEntry(controllerwails.GetMasterDictionaryEntryRequestDTO{ID: created.RefreshTargetID})
	if err != nil {
		t.Fatalf("expected deleted entry detail lookup to map not found to nil: %v", err)
	}

	if deletedDetail.Entry != nil {
		t.Fatalf("expected deleted entry detail to be nil, got %#v", deletedDetail.Entry)
	}

	controller.OnShutdown(context.Background())
}

func TestNewAppControllerPersistsEntriesAcrossControllerRecreation(t *testing.T) {
	configureBootstrapTestDatabase(t)
	firstController := NewAppController()
	created := mustBootstrapCreatePersistedEntry(t, firstController)

	firstController.OnShutdown(context.Background())

	secondController := NewAppController()
	loaded, err := secondController.GetMasterDictionaryEntry(controllerwails.GetMasterDictionaryEntryRequestDTO{ID: created.RefreshTargetID})
	if err != nil {
		t.Fatalf("expected persisted entry lookup through second controller to succeed: %v", err)
	}

	if loaded.Entry == nil || loaded.Entry.Source != bootstrapPersistedSource {
		t.Fatalf("expected recreated controller to read persisted entry, got %#v", loaded.Entry)
	}

	secondController.OnShutdown(context.Background())
}

func TestMasterDictionaryDatabasePathDefaultsToRepositoryRootDB(t *testing.T) {
	t.Setenv("AITRANSLATIONENGINEJP_MASTER_DICTIONARY_DB_PATH", "")

	repositoryRoot, err := repositoryRootDirectory()
	if err != nil {
		t.Fatalf("expected repository root directory to resolve: %v", err)
	}
	got := masterDictionaryDatabasePath()
	want := filepath.Join(repositoryRoot, "db", "master-dictionary.sqlite3")
	if got != want {
		t.Fatalf("expected database path %q, got %q", want, got)
	}
}

func runBootstrapImport(t *testing.T) (controllerwails.MasterDictionaryImportResponseDTO, *recordingRuntimeEventEmitter, *controllerwails.AppController) {
	t.Helper()

	configureBootstrapTestDatabase(t)
	xmlPath := writeBootstrapImportFixture(t)
	controller := NewAppController()
	emitter := &recordingRuntimeEventEmitter{}
	controller.OnStartup(runtimeEventsTestContext{Context: context.Background(), emitter: emitter})

	result, err := controller.MasterDictionaryImportXML(controllerwails.MasterDictionaryImportRequestDTO{
		XMLPath: xmlPath,
		Refresh: controllerwails.MasterDictionaryRefreshQueryDTO{Category: "すべて", Page: 1, PageSize: 30},
	})
	if err != nil {
		t.Fatalf("expected import through production graph to succeed: %v", err)
	}

	return result, emitter, controller
}

func runBootstrapCreate(t *testing.T) (controllerwails.CreateMasterDictionaryEntryResponseDTO, *controllerwails.AppController) {
	t.Helper()

	configureBootstrapTestDatabase(t)
	controller := NewAppController()
	return mustBootstrapCreateEntry(t, controller), controller
}

func mustBootstrapCreateEntry(t *testing.T, controller *controllerwails.AppController) controllerwails.CreateMasterDictionaryEntryResponseDTO {
	t.Helper()

	createRequest := controllerwails.CreateMasterDictionaryEntryRequestDTO{}
	createRequest.Payload.Source = bootstrapCreatedSource
	createRequest.Payload.Translation = bootstrapCreatedTranslation
	createRequest.Payload.Category = "武器"
	createRequest.Payload.Origin = "手動登録"
	createRequest.Refresh = &controllerwails.MasterDictionaryFrontendRefreshDTO{Query: bootstrapCreatedSource, Category: "武器", Page: 1, PageSize: 10}
	created, err := controller.CreateMasterDictionaryEntry(createRequest)
	if err != nil {
		t.Fatalf("expected create through production graph to succeed: %v", err)
	}
	return created
}

func runBootstrapUpdate(t *testing.T, controller *controllerwails.AppController, entryID string) controllerwails.UpdateMasterDictionaryEntryResponseDTO {
	t.Helper()

	updateRequest := controllerwails.UpdateMasterDictionaryEntryRequestDTO{ID: entryID}
	updateRequest.Payload.Source = bootstrapCreatedSource
	updateRequest.Payload.Translation = bootstrapUpdatedTranslation
	updateRequest.Payload.Category = "武器"
	updateRequest.Payload.Origin = "手動登録"
	updateRequest.Refresh = &controllerwails.MasterDictionaryFrontendRefreshDTO{Query: bootstrapCreatedSource, Category: "武器", Page: 1, PageSize: 10}
	updated, err := controller.UpdateMasterDictionaryEntry(updateRequest)
	if err != nil {
		t.Fatalf("expected update through production graph to succeed: %v", err)
	}
	return updated
}

func runBootstrapDelete(t *testing.T, controller *controllerwails.AppController, entryID string) controllerwails.DeleteMasterDictionaryEntryResponseDTO {
	t.Helper()

	deleted, err := controller.DeleteMasterDictionaryEntry(controllerwails.DeleteMasterDictionaryEntryRequestDTO{
		ID:      entryID,
		Refresh: &controllerwails.MasterDictionaryFrontendRefreshDTO{Query: bootstrapCreatedSource, Category: "武器", Page: 1, PageSize: 10},
	})
	if err != nil {
		t.Fatalf("expected delete through production graph to succeed: %v", err)
	}
	return deleted
}

func mustBootstrapCreatePersistedEntry(t *testing.T, controller *controllerwails.AppController) controllerwails.CreateMasterDictionaryEntryResponseDTO {
	t.Helper()

	createRequest := controllerwails.CreateMasterDictionaryEntryRequestDTO{}
	createRequest.Payload.Source = bootstrapPersistedSource
	createRequest.Payload.Translation = "ヴォルキハル城"
	createRequest.Payload.Category = "地名"
	createRequest.Payload.Origin = "手動登録"
	createRequest.Refresh = &controllerwails.MasterDictionaryFrontendRefreshDTO{Query: bootstrapPersistedSource, Category: "地名", Page: 1, PageSize: 10}
	created, err := controller.CreateMasterDictionaryEntry(createRequest)
	if err != nil {
		t.Fatalf("expected create through first controller to succeed: %v", err)
	}
	return created
}

func writeBootstrapImportFixture(t *testing.T) string {
	t.Helper()

	tmpDir := t.TempDir()
	xmlPath := filepath.Join(tmpDir, "bootstrap-import.xml")
	xmlContent := `<?xml version="1.0" encoding="utf-8"?>
<Root>
	<String>
		<REC>WEAP:FULL</REC>
		<EDID>DLC1AurielsShield</EDID>
		<Source>Auriel's Shield</Source>
		<Dest>アーリエルの盾</Dest>
	</String>
</Root>`
	if err := os.WriteFile(xmlPath, []byte(xmlContent), 0o600); err != nil {
		t.Fatalf("write xml fixture: %v", err)
	}
	return xmlPath
}

func configureBootstrapTestDatabase(t *testing.T) string {
	t.Helper()

	databasePath := filepath.Join(t.TempDir(), "db", "master-dictionary.sqlite3")
	t.Setenv("AITRANSLATIONENGINEJP_MASTER_DICTIONARY_DB_PATH", databasePath)
	return databasePath
}

func assertEventNames(t *testing.T, events []recordedRuntimeEvent, expected ...string) {
	t.Helper()

	if len(events) != len(expected) {
		t.Fatalf("expected %d runtime events, got %d: %#v", len(expected), len(events), events)
	}
	for index, eventName := range expected {
		if events[index].name != eventName {
			t.Fatalf("expected event[%d]=%q, got %q", index, eventName, events[index].name)
		}
	}
}
