package bootstrap

import (
	"bytes"
	"context"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	controllerwails "aitranslationenginejp/internal/controller/wails"
	ai "aitranslationenginejp/internal/infra/ai"
	"aitranslationenginejp/internal/repository"
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

type bootstrapProviderCaptureTransport struct {
	lastRequest *http.Request
}

func (transport *bootstrapProviderCaptureTransport) Do(req *http.Request) (*http.Response, error) {
	transport.lastRequest = req
	return &http.Response{
		StatusCode: http.StatusOK,
		Body:       io.NopCloser(bytes.NewBufferString(`{"choices":[{"message":{"content":"ok"}}]}`)),
		Header:     make(http.Header),
	}, nil
}

const (
	bootstrapCreatedSource               = "Auriel's Bow"
	bootstrapCreatedTranslation          = "アーリエルの弓"
	bootstrapUpdatedTranslation          = "更新済みアーリエルの弓"
	bootstrapPersistedSource             = "Relic Hammer"
	bootstrapPersistedTranslation        = "遺物の鎚"
	bootstrapSeedSearchTerm              = "Relic"
	bootstrapSeedCategory                = "武器"
	bootstrapSeedNewestSource            = "Relic Blade"
	bootstrapSeedNewestTranslation       = "遺物の剣"
	bootstrapImportProgressEvent         = "master-dictionary:import-progress"
	bootstrapImportCompletedEvent        = "master-dictionary:import-completed"
	bootstrapPageQuerySucceeded          = "expected bootstrap graph page query to succeed: %v"
	bootstrapMasterPersonaIdentityKey    = "FollowersPlus.esp:FE01A812:NPC_"
	bootstrapMasterPersonaPersistedBody  = "再生成後も保持されるペルソナ本文"
	bootstrapMasterPersonaPersistedModel = "persisted-master-persona-model"
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
	controller := newBootstrapTestControllerWithDatabasePath(t, databasePath)

	_, err := controller.MasterDictionaryGetPage(controllerwails.MasterDictionaryPageRequestDTO{
		Refresh: controllerwails.MasterDictionaryRefreshQueryDTO{
			SearchTerm: bootstrapSeedSearchTerm,
			Category:   bootstrapSeedCategory,
			Page:       1,
			PageSize:   5,
		},
	})
	if err != nil {
		t.Fatalf(bootstrapPageQuerySucceeded, err)
	}

	_, err = os.Stat(databasePath)
	if err != nil {
		t.Fatalf("expected sqlite database file to exist: %v", err)
	}
}

func TestNewAppControllerUsesProductionSeedByDefault(t *testing.T) {
	configureBootstrapTestDatabase(t)
	controller := NewAppController()

	page, err := controller.MasterDictionaryGetPage(controllerwails.MasterDictionaryPageRequestDTO{
		Refresh: controllerwails.MasterDictionaryRefreshQueryDTO{
			SearchTerm: "Whiterun",
			Category:   "地名",
			Page:       1,
			PageSize:   5,
		},
	})
	if err != nil {
		t.Fatalf("expected direct constructor query to succeed: %v", err)
	}

	if page.Page.TotalCount == 0 || len(page.Page.Items) == 0 {
		t.Fatalf("expected production-seeded page items, got %#v", page.Page)
	}
	if !strings.Contains(page.Page.Items[0].Source, "Whiterun") || page.Page.Items[0].Category != "地名" {
		t.Fatalf("expected production-seeded Whiterun location entry, got %#v", page.Page.Items[0])
	}
}

func TestNewAppControllerMasterPersonaProductionWiringDoesNotUseInMemoryConcrete(t *testing.T) {
	content, err := os.ReadFile("app_controller.go")
	if err != nil {
		t.Fatalf("expected app controller source to be readable: %v", err)
	}

	source := string(content)
	if strings.Contains(source, "NewInMemoryMasterPersonaRepository") {
		t.Fatalf("production wiring must not use in-memory master persona repository")
	}
	if strings.Contains(source, "NewInMemorySecretStore") {
		t.Fatalf("production wiring must not use in-memory master persona secret store")
	}
}

func TestNewAppControllerReturnsSeededPageItems(t *testing.T) {
	controller := newBootstrapTestController(t)

	page, err := controller.MasterDictionaryGetPage(controllerwails.MasterDictionaryPageRequestDTO{
		Refresh: controllerwails.MasterDictionaryRefreshQueryDTO{
			SearchTerm: bootstrapSeedSearchTerm,
			Category:   bootstrapSeedCategory,
			Page:       1,
			PageSize:   5,
		},
	})
	if err != nil {
		t.Fatalf(bootstrapPageQuerySucceeded, err)
	}

	if page.Page.TotalCount != len(bootstrapTestSeed()) || len(page.Page.Items) == 0 {
		t.Fatalf("expected seeded page items, got %#v", page.Page)
	}
}

func TestNewAppControllerSelectsFirstPageItemByDefault(t *testing.T) {
	controller := newBootstrapTestController(t)

	page, err := controller.MasterDictionaryGetPage(controllerwails.MasterDictionaryPageRequestDTO{
		Refresh: controllerwails.MasterDictionaryRefreshQueryDTO{
			SearchTerm: bootstrapSeedSearchTerm,
			Category:   bootstrapSeedCategory,
			Page:       1,
			PageSize:   5,
		},
	})
	if err != nil {
		t.Fatalf(bootstrapPageQuerySucceeded, err)
	}

	if page.Page.SelectedID == nil || *page.Page.SelectedID != page.Page.Items[0].ID {
		t.Fatalf("expected selected id to match first page item, got selected=%#v items=%#v", page.Page.SelectedID, page.Page.Items)
	}
}

func TestNewAppControllerReturnsDetailForSeededEntry(t *testing.T) {
	controller := newBootstrapTestController(t)
	page, err := controller.MasterDictionaryGetPage(controllerwails.MasterDictionaryPageRequestDTO{
		Refresh: controllerwails.MasterDictionaryRefreshQueryDTO{
			SearchTerm: bootstrapSeedSearchTerm,
			Category:   bootstrapSeedCategory,
			Page:       1,
			PageSize:   5,
		},
	})
	if err != nil {
		t.Fatalf(bootstrapPageQuerySucceeded, err)
	}

	detail, err := controller.MasterDictionaryGetDetail(controllerwails.MasterDictionaryDetailRequestDTO{ID: page.Page.Items[0].ID})
	if err != nil {
		t.Fatalf("expected detail query to succeed: %v", err)
	}

	if detail.Entry.Source != bootstrapSeedNewestSource || detail.Entry.Translation != bootstrapSeedNewestTranslation {
		t.Fatalf("expected newest seeded entry, got %#v", detail.Entry)
	}
}

func TestNewAppControllerStartupDoesNotEmitRuntimeEvents(t *testing.T) {
	controller := newBootstrapTestController(t)
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
	databasePath := configureBootstrapTestDatabase(t)
	firstController := newBootstrapTestControllerWithDatabasePath(t, databasePath)
	created := mustBootstrapCreatePersistedEntry(t, firstController)

	firstController.OnShutdown(context.Background())

	secondController := newBootstrapTestControllerWithDatabasePath(t, databasePath)
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

func TestMasterPersonaAIModeDefaultsToReal(t *testing.T) {
	t.Setenv(masterPersonaAIModeEnv, "")
	if got := masterPersonaAIMode(); got != masterPersonaAIModeReal {
		t.Fatalf("expected ai mode to default to real, got %q", got)
	}
}

func TestMasterPersonaAIModeFakeWhenConfigured(t *testing.T) {
	t.Setenv(masterPersonaAIModeEnv, masterPersonaAIModeFake)
	if got := masterPersonaAIMode(); got != masterPersonaAIModeFake {
		t.Fatalf("expected ai mode fake, got %q", got)
	}
}

func TestNewAIProviderClientFromMasterPersonaEnvUsesFakeModeResponseOverride(t *testing.T) {
	t.Setenv(masterPersonaAIModeEnv, masterPersonaAIModeFake)
	t.Setenv(masterPersonaFakeResponseEnv, "env fake response")
	client := newAIProviderClientFromMasterPersonaEnv()

	geminiResponse, err := client.GenerateText(context.Background(), "gemini", ai.ProviderRequest{Model: "gemini-2.5-pro", Prompt: "prompt"})
	if err != nil {
		t.Fatalf("expected fake mode gemini generation to succeed: %v", err)
	}
	if geminiResponse.Text != "env fake response" {
		t.Fatalf("expected fake response override text, got %q", geminiResponse.Text)
	}

	xaiResponse, err := client.GenerateText(context.Background(), "xai", ai.ProviderRequest{Model: "grok-2", Prompt: "prompt"})
	if err != nil {
		t.Fatalf("expected fake mode xai generation to succeed without api key: %v", err)
	}
	if xaiResponse.Text != "env fake response" {
		t.Fatalf("expected fake response override text for xai, got %q", xaiResponse.Text)
	}
}

func TestNewAIProviderClientFromMasterPersonaEnvAppliesBaseURLOverrides(t *testing.T) {
	t.Setenv(masterPersonaAIModeEnv, masterPersonaAIModeReal)
	t.Setenv(masterPersonaXAIBaseURLEnv, "https://gateway.example.com/custom/v1")
	t.Setenv(masterPersonaLMStudioBaseURLEnv, "http://127.0.0.1:1234/proxy/v1")

	xaiTransport := &bootstrapProviderCaptureTransport{}
	xaiClient := newAIProviderClientFromMasterPersonaEnvWithTransport(xaiTransport)
	if _, err := xaiClient.GenerateText(context.Background(), "xai", ai.ProviderRequest{Model: "grok-2", APIKey: "x-key", Prompt: "prompt"}); err != nil {
		t.Fatalf("expected xai generation with override base url to succeed: %v", err)
	}
	if xaiTransport.lastRequest == nil {
		t.Fatalf("expected xai provider request capture")
	}
	if got := xaiTransport.lastRequest.URL.String(); got != "https://gateway.example.com/custom/v1/chat/completions" {
		t.Fatalf("expected xai request url override, got %q", got)
	}
	if got := xaiTransport.lastRequest.Header.Get("Authorization"); got != "Bearer x-key" {
		t.Fatalf("expected xai authorization header, got %q", got)
	}

	lmStudioTransport := &bootstrapProviderCaptureTransport{}
	lmStudioClient := newAIProviderClientFromMasterPersonaEnvWithTransport(lmStudioTransport)
	if _, err := lmStudioClient.GenerateText(context.Background(), "lm_studio", ai.ProviderRequest{Model: "local-model", Prompt: "prompt"}); err != nil {
		t.Fatalf("expected lm studio generation with override base url to succeed: %v", err)
	}
	if lmStudioTransport.lastRequest == nil {
		t.Fatalf("expected lm studio provider request capture")
	}
	if got := lmStudioTransport.lastRequest.URL.String(); got != "http://127.0.0.1:1234/proxy/v1/chat/completions" {
		t.Fatalf("expected lm studio request url override, got %q", got)
	}
	if got := lmStudioTransport.lastRequest.Header.Get("Authorization"); got != "" {
		t.Fatalf("expected lm studio authorization header to be omitted when api key blank, got %q", got)
	}
}

func TestNewAppControllerProvidesMasterPersonaPage(t *testing.T) {
	controller := newBootstrapTestController(t)

	page, err := controller.MasterPersonaGetPage(controllerwails.MasterPersonaPageRequestDTO{
		Refresh: controllerwails.MasterPersonaListQueryDTO{PluginFilter: "FollowersPlus.esp", Page: 1, PageSize: 10},
	})
	if err != nil {
		t.Fatalf("expected master persona page query to succeed: %v", err)
	}
	if page.Page.TotalCount == 0 || len(page.Page.Items) == 0 {
		t.Fatalf("expected master persona seed entries, got %#v", page.Page)
	}
}

func TestNewAppControllerProvidesMasterPersonaAISettingsPersistence(t *testing.T) {
	controller := newBootstrapTestController(t)

	saved, err := controller.MasterPersonaSaveAISettings(controllerwails.MasterPersonaAISettingsDTO{Provider: "gemini", Model: "gemini-2.5-pro", APIKey: "test-key"})
	if err != nil {
		t.Fatalf("expected master persona ai settings save to succeed: %v", err)
	}
	if saved.Provider != "gemini" || saved.Model != "gemini-2.5-pro" {
		t.Fatalf("unexpected saved settings: %#v", saved)
	}

	loaded, err := controller.MasterPersonaLoadAISettings()
	if err != nil {
		t.Fatalf("expected master persona ai settings load to succeed: %v", err)
	}
	if loaded.APIKey != "test-key" {
		t.Fatalf("expected saved api key to load, got %#v", loaded)
	}
}

func TestNewAppControllerPersistsMasterPersonaEntryAcrossControllerRecreation(t *testing.T) {
	databasePath := configureBootstrapTestDatabase(t)
	firstController := newBootstrapTestControllerWithDatabasePath(t, databasePath)
	originalDetail, err := firstController.MasterPersonaGetDetail(controllerwails.MasterPersonaDetailRequestDTO{IdentityKey: bootstrapMasterPersonaIdentityKey})
	if err != nil {
		t.Fatalf("expected master persona detail lookup through first controller to succeed: %v", err)
	}

	_, err = firstController.MasterPersonaUpdate(controllerwails.MasterPersonaUpdateRequestDTO{
		IdentityKey: bootstrapMasterPersonaIdentityKey,
		Entry: controllerwails.MasterPersonaUpdateInputDTO{
			FormID:       originalDetail.Entry.FormID,
			EditorID:     originalDetail.Entry.EditorID,
			DisplayName:  originalDetail.Entry.DisplayName,
			Race:         originalDetail.Entry.Race,
			Sex:          originalDetail.Entry.Sex,
			VoiceType:    originalDetail.Entry.VoiceType,
			ClassName:    originalDetail.Entry.ClassName,
			SourcePlugin: originalDetail.Entry.SourcePlugin,
			PersonaBody:  bootstrapMasterPersonaPersistedBody,
		},
		Refresh: controllerwails.MasterPersonaListQueryDTO{PluginFilter: "FollowersPlus.esp", Page: 1, PageSize: 10},
	})
	if err != nil {
		t.Fatalf("expected master persona update through first controller to succeed: %v", err)
	}
	firstController.OnShutdown(context.Background())

	secondController := newBootstrapTestControllerWithDatabasePath(t, databasePath)
	persistedDetail, err := secondController.MasterPersonaGetDetail(controllerwails.MasterPersonaDetailRequestDTO{IdentityKey: bootstrapMasterPersonaIdentityKey})
	if err != nil {
		t.Fatalf("expected master persona detail lookup through second controller to succeed: %v", err)
	}
	if persistedDetail.Entry.PersonaBody != bootstrapMasterPersonaPersistedBody {
		t.Fatalf("expected updated master persona body to persist, got %#v", persistedDetail.Entry)
	}

	secondController.OnShutdown(context.Background())
}

func TestNewAppControllerPersistsMasterPersonaAISettingsAcrossControllerRecreation(t *testing.T) {
	databasePath := configureBootstrapTestDatabase(t)
	firstController := newBootstrapTestControllerWithDatabasePath(t, databasePath)

	_, err := firstController.MasterPersonaSaveAISettings(controllerwails.MasterPersonaAISettingsDTO{
		Provider: "gemini",
		Model:    bootstrapMasterPersonaPersistedModel,
		APIKey:   "volatile-key",
	})
	if err != nil {
		t.Fatalf("expected master persona ai settings save through first controller to succeed: %v", err)
	}
	firstController.OnShutdown(context.Background())

	secondController := newBootstrapTestControllerWithDatabasePath(t, databasePath)
	loaded, err := secondController.MasterPersonaLoadAISettings()
	if err != nil {
		t.Fatalf("expected master persona ai settings load through second controller to succeed: %v", err)
	}
	if loaded.Provider != "gemini" || loaded.Model != bootstrapMasterPersonaPersistedModel {
		t.Fatalf("expected provider/model settings to persist, got %#v", loaded)
	}
	if loaded.APIKey != "volatile-key" {
		t.Fatalf("expected api key to persist through keyring backend, got %#v", loaded)
	}

	secondController.OnShutdown(context.Background())
}

func TestNewAppControllerPersistsMasterPersonaRunStatusAcrossControllerRecreation(t *testing.T) {
	databasePath := configureBootstrapTestDatabase(t)
	firstController := newBootstrapTestControllerWithDatabasePath(t, databasePath)
	extractPath := writeBootstrapMasterPersonaExtractFixture(t, `{
  "target_plugin": "FollowersPlus.esp",
  "npcs": [
    {
      "form_id": "FE01AFF0",
      "record_type": "NPC_",
      "editor_id": "FP_Persist",
      "display_name": "Persist",
      "dialogues": ["hello"]
    }
  ]
}`)

	executed, err := firstController.MasterPersonaExecuteGeneration(controllerwails.MasterPersonaExecuteRequestDTO{
		FilePath: extractPath,
		AISettings: controllerwails.MasterPersonaAISettingsDTO{
			Provider: "gemini",
			Model:    "gemini-2.5-pro",
		},
	})
	if err != nil {
		t.Fatalf("expected master persona execute through first controller to succeed: %v", err)
	}
	if executed.RunState != "完了" || executed.SuccessCount != 1 {
		t.Fatalf("expected completed run status through first controller, got %#v", executed)
	}
	firstController.OnShutdown(context.Background())

	secondController := newBootstrapTestControllerWithDatabasePath(t, databasePath)
	loadedStatus, err := secondController.MasterPersonaGetRunStatus()
	if err != nil {
		t.Fatalf("expected master persona run status load through second controller to succeed: %v", err)
	}
	if loadedStatus.RunState != "完了" || loadedStatus.SuccessCount != 1 || loadedStatus.ProcessedCount != 1 {
		t.Fatalf("expected run status to persist after controller recreation, got %#v", loadedStatus)
	}

	persistedIdentityKey := repository.BuildMasterPersonaIdentityKey("FollowersPlus.esp", "FE01AFF0", "NPC_")
	persistedDetail, err := secondController.MasterPersonaGetDetail(controllerwails.MasterPersonaDetailRequestDTO{IdentityKey: persistedIdentityKey})
	if err != nil {
		t.Fatalf("expected generated master persona entry to persist across controller recreation: %v", err)
	}
	if persistedDetail.Entry.DisplayName != "Persist" {
		t.Fatalf("expected generated master persona entry to persist, got %#v", persistedDetail.Entry)
	}

	secondController.OnShutdown(context.Background())
}

func runBootstrapImport(t *testing.T) (controllerwails.MasterDictionaryImportResponseDTO, *recordingRuntimeEventEmitter, *controllerwails.AppController) {
	t.Helper()

	xmlPath := writeBootstrapImportFixture(t)
	controller := newBootstrapTestController(t)
	emitter := &recordingRuntimeEventEmitter{}
	controller.OnStartup(runtimeEventsTestContext{Context: context.Background(), emitter: emitter})

	result, err := controller.MasterDictionaryImportXML(controllerwails.MasterDictionaryImportRequestDTO{
		XMLPath: xmlPath,
		Refresh: controllerwails.MasterDictionaryRefreshQueryDTO{Category: "すべて", Page: 1, PageSize: 30},
	})
	if err != nil {
		t.Fatalf("expected import through bootstrap graph to succeed: %v", err)
	}

	return result, emitter, controller
}

func runBootstrapCreate(t *testing.T) (controllerwails.CreateMasterDictionaryEntryResponseDTO, *controllerwails.AppController) {
	t.Helper()

	controller := newBootstrapTestController(t)
	return mustBootstrapCreateEntry(t, controller), controller
}

func mustBootstrapCreateEntry(t *testing.T, controller *controllerwails.AppController) controllerwails.CreateMasterDictionaryEntryResponseDTO {
	t.Helper()

	createRequest := controllerwails.CreateMasterDictionaryEntryRequestDTO{}
	createRequest.Payload.Source = bootstrapCreatedSource
	createRequest.Payload.Translation = bootstrapCreatedTranslation
	createRequest.Payload.Category = bootstrapSeedCategory
	createRequest.Payload.Origin = "手動登録"
	createRequest.Refresh = &controllerwails.MasterDictionaryFrontendRefreshDTO{Query: bootstrapCreatedSource, Category: bootstrapSeedCategory, Page: 1, PageSize: 10}
	created, err := controller.CreateMasterDictionaryEntry(createRequest)
	if err != nil {
		t.Fatalf("expected create through bootstrap graph to succeed: %v", err)
	}
	return created
}

func runBootstrapUpdate(t *testing.T, controller *controllerwails.AppController, entryID string) controllerwails.UpdateMasterDictionaryEntryResponseDTO {
	t.Helper()

	updateRequest := controllerwails.UpdateMasterDictionaryEntryRequestDTO{ID: entryID}
	updateRequest.Payload.Source = bootstrapCreatedSource
	updateRequest.Payload.Translation = bootstrapUpdatedTranslation
	updateRequest.Payload.Category = bootstrapSeedCategory
	updateRequest.Payload.Origin = "手動登録"
	updateRequest.Refresh = &controllerwails.MasterDictionaryFrontendRefreshDTO{Query: bootstrapCreatedSource, Category: bootstrapSeedCategory, Page: 1, PageSize: 10}
	updated, err := controller.UpdateMasterDictionaryEntry(updateRequest)
	if err != nil {
		t.Fatalf("expected update through bootstrap graph to succeed: %v", err)
	}
	return updated
}

func runBootstrapDelete(t *testing.T, controller *controllerwails.AppController, entryID string) controllerwails.DeleteMasterDictionaryEntryResponseDTO {
	t.Helper()

	deleted, err := controller.DeleteMasterDictionaryEntry(controllerwails.DeleteMasterDictionaryEntryRequestDTO{
		ID:      entryID,
		Refresh: &controllerwails.MasterDictionaryFrontendRefreshDTO{Query: bootstrapCreatedSource, Category: bootstrapSeedCategory, Page: 1, PageSize: 10},
	})
	if err != nil {
		t.Fatalf("expected delete through bootstrap graph to succeed: %v", err)
	}
	return deleted
}

func mustBootstrapCreatePersistedEntry(t *testing.T, controller *controllerwails.AppController) controllerwails.CreateMasterDictionaryEntryResponseDTO {
	t.Helper()

	createRequest := controllerwails.CreateMasterDictionaryEntryRequestDTO{}
	createRequest.Payload.Source = bootstrapPersistedSource
	createRequest.Payload.Translation = bootstrapPersistedTranslation
	createRequest.Payload.Category = bootstrapSeedCategory
	createRequest.Payload.Origin = "手動登録"
	createRequest.Refresh = &controllerwails.MasterDictionaryFrontendRefreshDTO{Query: bootstrapPersistedSource, Category: bootstrapSeedCategory, Page: 1, PageSize: 10}
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

func writeBootstrapMasterPersonaExtractFixture(t *testing.T, content string) string {
	t.Helper()

	extractPath := filepath.Join(t.TempDir(), "extract.json")
	if err := os.WriteFile(extractPath, []byte(content), 0o600); err != nil {
		t.Fatalf("write master persona extract fixture: %v", err)
	}
	return extractPath
}

func newBootstrapTestController(t *testing.T) *controllerwails.AppController {
	t.Helper()

	databasePath := configureBootstrapTestDatabase(t)
	return newBootstrapTestControllerWithDatabasePath(t, databasePath)
}

func newBootstrapTestControllerWithDatabasePath(t *testing.T, databasePath string) *controllerwails.AppController {
	t.Helper()

	setBootstrapTestDatabasePath(t, databasePath)
	setBootstrapTestMasterPersonaSecretStore(t, databasePath)
	return newAppControllerWithMasterDictionarySeed(bootstrapTestSeed(), bootstrapTestNow)
}

func configureBootstrapTestDatabase(t *testing.T) string {
	t.Helper()

	databasePath := filepath.Join(t.TempDir(), "db", "master-dictionary.sqlite3")
	setBootstrapTestDatabasePath(t, databasePath)
	return databasePath
}

func setBootstrapTestDatabasePath(t *testing.T, databasePath string) {
	t.Helper()

	t.Setenv("AITRANSLATIONENGINEJP_MASTER_DICTIONARY_DB_PATH", databasePath)
}

func setBootstrapTestMasterPersonaSecretStore(t *testing.T, databasePath string) {
	t.Helper()

	t.Setenv("AITRANSLATIONENGINEJP_TEST_MODE", "true")
	t.Setenv(masterPersonaAIModeEnv, masterPersonaAIModeFake)
	t.Setenv("AITRANSLATIONENGINEJP_MASTER_PERSONA_SECRET_BACKEND", "file")
	t.Setenv(
		"AITRANSLATIONENGINEJP_MASTER_PERSONA_SECRET_FILE_DIR",
		filepath.Join(filepath.Dir(databasePath), "master-persona-keyring"),
	)
	t.Setenv("AITRANSLATIONENGINEJP_MASTER_PERSONA_SECRET_FILE_PASSWORD", "bootstrap-test-keyring-password")
}

func bootstrapTestNow() time.Time {
	return time.Date(2026, 4, 15, 12, 0, 0, 0, time.UTC)
}

func bootstrapTestSeed() []repository.MasterDictionaryEntry {
	baseTime := time.Date(2026, 4, 15, 9, 0, 0, 0, time.UTC)
	return []repository.MasterDictionaryEntry{
		{
			ID:          1,
			Source:      "Relic Bow",
			Translation: "遺物の弓",
			Category:    bootstrapSeedCategory,
			Origin:      "テスト初期データ",
			REC:         "WEAP:FULL",
			EDID:        "RelicBow",
			UpdatedAt:   baseTime,
		},
		{
			ID:          2,
			Source:      bootstrapSeedNewestSource,
			Translation: bootstrapSeedNewestTranslation,
			Category:    bootstrapSeedCategory,
			Origin:      "テスト初期データ",
			REC:         "WEAP:FULL",
			EDID:        "RelicBlade",
			UpdatedAt:   baseTime.Add(time.Minute),
		},
	}
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
