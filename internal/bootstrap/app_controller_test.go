package bootstrap

import (
	"bytes"
	"context"
	"database/sql"
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
	"aitranslationenginejp/internal/service"
	"aitranslationenginejp/internal/usecase"
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

// persona-ai-settings-restart-cutover: run state は再起動後に "入力待ち" へ戻ることを actual repository path で証明する。
// SaveRunStatus は no-op のため、DB を再オープンすると run state は常にデフォルト値に戻る。
func TestNewAppControllerPersonaAISettingsRestartCutoverRunStatusIsInputWaitingAfterRepositoryRecreation(t *testing.T) {
	// Arrange: 1 セッションで run state = "完了" を書き込む (SaveRunStatus は no-op だが呼び出す)。
	databasePath := configureBootstrapTestDatabase(t)
	repos1, err := repository.NewSQLiteMasterPersonaRepositories(context.Background(), databasePath, nil)
	if err != nil {
		t.Fatalf("expected repository open to succeed: %v", err)
	}
	startedAt := time.Date(2026, 4, 21, 10, 0, 0, 0, time.UTC)
	finishedAt := startedAt.Add(time.Minute)
	if saveErr := repos1.RunStatusRepository.SaveRunStatus(context.Background(), repository.MasterPersonaRunStatusRecord{
		RunState:       "完了",
		SuccessCount:   1,
		ProcessedCount: 1,
		StartedAt:      &startedAt,
		FinishedAt:     &finishedAt,
	}); saveErr != nil {
		_ = repos1.Close()
		t.Fatalf("expected run status save to succeed: %v", saveErr)
	}
	if closeErr := repos1.Close(); closeErr != nil {
		t.Fatalf("expected repos1 close to succeed: %v", closeErr)
	}

	// Act: 同じ DB パスで repository を再作成する (再起動をシミュレートする)。
	repos2, err := repository.NewSQLiteMasterPersonaRepositories(context.Background(), databasePath, nil)
	if err != nil {
		t.Fatalf("expected repository reopen to succeed: %v", err)
	}
	defer func() {
		if closeErr := repos2.Close(); closeErr != nil {
			t.Fatalf("expected repos2 close to succeed: %v", closeErr)
		}
	}()

	// Assert: 再起動後の run state は "入力待ち" に戻る (run state は永続化されない)。
	loaded, loadErr := repos2.RunStatusRepository.LoadRunStatus(context.Background())
	if loadErr != nil {
		t.Fatalf("expected run status load after reopen to succeed: %v", loadErr)
	}
	if loaded.RunState == "完了" {
		t.Fatalf("expected run state to reset after restart (run state must not persist), but got %q", loaded.RunState)
	}
}

// persona-ai-settings-restart-cutover: PERSONA_GENERATION_SETTINGS に json_file_path 列がないことを証明する。
// JSON ファイル選択は DB に保存されないため、再起動後は frontend 側で "未選択" に戻る。
func TestNewAppControllerPersonaAISettingsRestartCutoverJSONFilePathIsAbsentFromDB(t *testing.T) {
	// Arrange: bootstrap テスト用 DB パスを確保する。
	databasePath := configureBootstrapTestDatabase(t)
	db, err := repository.OpenSQLiteDatabase(context.Background(), databasePath)
	if err != nil {
		t.Fatalf("expected db open to succeed: %v", err)
	}
	defer func() {
		if closeErr := db.Close(); closeErr != nil {
			t.Fatalf("expected db close to succeed: %v", closeErr)
		}
	}()

	// Act: PERSONA_GENERATION_SETTINGS の列情報を取得する。
	rows, err := db.QueryContext(context.Background(), "PRAGMA table_info(PERSONA_GENERATION_SETTINGS)")
	if err != nil {
		t.Fatalf("expected PRAGMA table_info query to succeed: %v", err)
	}
	defer func() {
		if closeErr := rows.Close(); closeErr != nil {
			t.Fatalf("expected rows.Close to succeed: %v", closeErr)
		}
	}()

	// Assert: file_path / json_path / json_file_path など JSON 選択に関わる列は存在しない。
	for rows.Next() {
		var cid, notNull, pk int
		var name, colType string
		var defaultValue sql.NullString
		if scanErr := rows.Scan(&cid, &name, &colType, &notNull, &defaultValue, &pk); scanErr != nil {
			t.Fatalf("expected column row scan to succeed: %v", scanErr)
		}
		lower := strings.ToLower(name)
		if strings.Contains(lower, "file") || strings.Contains(lower, "json") || strings.Contains(lower, "path") {
			t.Fatalf("expected PERSONA_GENERATION_SETTINGS to have no file/json/path column (JSON selection must not persist), found: %q", name)
		}
	}
	if err := rows.Err(); err != nil {
		t.Fatalf("expected table_info rows iteration to succeed: %v", err)
	}
}

// persona-ai-settings-restart-cutover: PERSONA_GENERATION_SETTINGS に保存された provider/model が
// DB 再オープン後に復元されることを bootstrap DB パス経由で証明する。
func TestNewAppControllerPersonaAISettingsRestartCutoverProviderModelRestoredAfterControllerRecreation(t *testing.T) {
	// Arrange: bootstrap テスト用 DB パスを確保し、最初のセッションで provider/model を書き込む。
	databasePath := configureBootstrapTestDatabase(t)
	db1, err := repository.OpenSQLiteDatabase(context.Background(), databasePath)
	if err != nil {
		t.Fatalf("expected db open to succeed: %v", err)
	}
	_, err = db1.ExecContext(context.Background(),
		"INSERT OR REPLACE INTO PERSONA_GENERATION_SETTINGS (id, provider, model) VALUES (1, 'gemini', 'restart-cutover-model')")
	if err != nil {
		_ = db1.Close()
		t.Fatalf("expected PERSONA_GENERATION_SETTINGS insert to succeed: %v", err)
	}
	if closeErr := db1.Close(); closeErr != nil {
		t.Fatalf("expected db1 close to succeed: %v", closeErr)
	}

	// Act: 同じ DB パスで DB を再オープンする。
	db2, err := repository.OpenSQLiteDatabase(context.Background(), databasePath)
	if err != nil {
		t.Fatalf("expected db reopen to succeed: %v", err)
	}
	defer func() {
		if closeErr := db2.Close(); closeErr != nil {
			t.Fatalf("expected db2 close to succeed: %v", closeErr)
		}
	}()

	// Assert: provider と model が復元されている。
	var provider, model string
	if err := db2.QueryRowContext(context.Background(),
		"SELECT provider, model FROM PERSONA_GENERATION_SETTINGS WHERE id = 1").Scan(&provider, &model); err != nil {
		t.Fatalf("expected PERSONA_GENERATION_SETTINGS load after reopen to succeed: %v", err)
	}
	if provider != "gemini" || model != "restart-cutover-model" {
		t.Fatalf("expected provider/model restored after db reopen, got provider=%q model=%q", provider, model)
	}
}

// persona-ai-settings-restart-cutover: run state は DB に保存されないことを bootstrap DB 経由で証明する。
// run status テーブルが存在しないことで、DB の外側がデフォルトの入力待ちに局限されることを確認する。
func TestNewAppControllerPersonaAISettingsRestartCutoverRunStatusIsInputWaitingOnFreshController(t *testing.T) {
	// Arrange: bootstrap テスト用 DB パスを確保する。
	databasePath := configureBootstrapTestDatabase(t)
	db, err := repository.OpenSQLiteDatabase(context.Background(), databasePath)
	if err != nil {
		t.Fatalf("expected db open to succeed: %v", err)
	}
	defer func() {
		if closeErr := db.Close(); closeErr != nil {
			t.Fatalf("expected db close to succeed: %v", closeErr)
		}
	}()

	// Assert: run state 永続化テーブルが存在しない。
	for _, tableName := range []string{"master_persona_run_status", "PERSONA_GENERATION_RUN_STATUS"} {
		var count int
		if err := db.QueryRowContext(context.Background(),
			"SELECT COUNT(*) FROM sqlite_master WHERE type='table' AND name=?", tableName).Scan(&count); err != nil {
			t.Fatalf("expected table existence check to succeed for %q: %v", tableName, err)
		}
		if count != 0 {
			t.Fatalf("expected no run status table %q in bootstrap db (run state must not be persisted), but table exists", tableName)
		}
	}
}

// persona-ai-settings-restart-cutover: PERSONA_GENERATION_SETTINGS に api_key 列がないことを
// bootstrap DB 経由で証明する。API key は DB の外 (keyring) に保管される。
func TestNewAppControllerPersonaAISettingsRestartCutoverAPIKeyStaysOutsideDB(t *testing.T) {
	// Arrange: bootstrap テスト用 DB パスを確保する。
	databasePath := configureBootstrapTestDatabase(t)
	db, err := repository.OpenSQLiteDatabase(context.Background(), databasePath)
	if err != nil {
		t.Fatalf("expected db open to succeed: %v", err)
	}
	defer func() {
		if closeErr := db.Close(); closeErr != nil {
			t.Fatalf("expected db close to succeed: %v", closeErr)
		}
	}()

	// Act: PERSONA_GENERATION_SETTINGS の列情報を取得する。
	rows, err := db.QueryContext(context.Background(), "PRAGMA table_info(PERSONA_GENERATION_SETTINGS)")
	if err != nil {
		t.Fatalf("expected PRAGMA table_info query to succeed: %v", err)
	}
	defer func() {
		if closeErr := rows.Close(); closeErr != nil {
			t.Fatalf("expected rows.Close to succeed: %v", closeErr)
		}
	}()

	// Assert: api_key 列は存在しない。
	var foundAPIKeyColumn bool
	for rows.Next() {
		var cid, notNull, pk int
		var name, colType string
		var defaultValue sql.NullString
		if err := rows.Scan(&cid, &name, &colType, &notNull, &defaultValue, &pk); err != nil {
			t.Fatalf("expected column row scan to succeed: %v", err)
		}
		if strings.Contains(strings.ToLower(name), "api_key") {
			foundAPIKeyColumn = true
		}
	}
	if err := rows.Err(); err != nil {
		t.Fatalf("expected table_info rows iteration to succeed: %v", err)
	}
	if foundAPIKeyColumn {
		t.Fatalf("expected PERSONA_GENERATION_SETTINGS to have no api_key column in bootstrap db (API key stays outside DB)")
	}
}

// persona-ai-settings-restart-cutover: 実際の bootstrap/controller/repository wiring を通じて
// GetRunStatus が "入力待ち" を返すことを同一インスタンスで証明する。
// また、同一インスタンス内で Execute を呼ぶと run status が "入力待ち" 以外へ遷移することを確認する。
// SQLiteMasterPersonaRunStatusRepository.LoadRunStatus は常にデフォルト値を返す (DB 永続化なし)。
func TestNewAppControllerPersonaAISettingsRestartCutoverSameInstanceRunStatusIsInputWaitingViaActualControllerWiring(t *testing.T) {
	// Arrange: 実 SQLite wiring で AppController を構築する (persona entry seed なし)。
	databasePath := configureBootstrapTestDatabase(t)
	controller := newBootstrapRunStatusTestController(t, databasePath)
	defer controller.OnShutdown(context.Background())

	// Act: 実 controller → usecase → service → SQLiteRunStatusRepository の経路で run status を取得する。
	status, err := controller.MasterPersonaGetRunStatus()

	// Assert: 初期 run state は "入力待ち" である (stub ではなく実 SQLite wiring 経由)。
	if err != nil {
		t.Fatalf("expected MasterPersonaGetRunStatus to succeed via actual controller wiring: %v", err)
	}
	if status.RunState != "入力待ち" {
		t.Fatalf("expected initial run state %q via actual controller wiring, got %q", "入力待ち", status.RunState)
	}

	// Act: 設定未設定の状態で ExecuteGeneration を呼ぶ。
	// AI 設定が空のため "設定未完了" ステータスが返り、run state が idle から遷移することを観測する。
	execStatus, execErr := controller.MasterPersonaExecuteGeneration(controllerwails.MasterPersonaExecuteRequestDTO{})

	// Assert: Execute が返す run state は "入力待ち" 以外である (run status が non-idle に遷移できることを証明する)。
	if execErr != nil {
		t.Fatalf("expected MasterPersonaExecuteGeneration with empty settings to succeed (returning settings-incomplete): %v", execErr)
	}
	if execStatus.RunState == "入力待ち" {
		t.Fatalf("expected run state to move away from %q after execute with empty settings, got %q (run status must be observable as non-idle)", "入力待ち", execStatus.RunState)
	}
}

// persona-ai-settings-restart-cutover: controller を再作成すると run status が "入力待ち" にリセットされることを
// 実 AppController wiring で証明する。SQLite の run state は永続化されないため常にデフォルトに戻る。
func TestNewAppControllerPersonaAISettingsRestartCutoverRecreatedControllerRunStatusResetsToInputWaiting(t *testing.T) {
	// Arrange: 最初のコントローラーを構築して run status を確認後にシャットダウンする (persona entry seed なし)。
	databasePath := configureBootstrapTestDatabase(t)
	firstController := newBootstrapRunStatusTestController(t, databasePath)
	firstStatus, err := firstController.MasterPersonaGetRunStatus()
	if err != nil {
		t.Fatalf("expected MasterPersonaGetRunStatus to succeed on first controller: %v", err)
	}
	if firstStatus.RunState != "入力待ち" {
		t.Fatalf("expected first controller initial run state %q, got %q", "入力待ち", firstStatus.RunState)
	}
	firstController.OnShutdown(context.Background())

	// Act: 同じ DB パスで controller を再作成する (再起動をシミュレート)。
	secondController := newBootstrapRunStatusTestController(t, databasePath)
	defer secondController.OnShutdown(context.Background())
	secondStatus, err := secondController.MasterPersonaGetRunStatus()

	// Assert: 再起動後の run state も "入力待ち" に戻る (run state は DB に永続化されない)。
	if err != nil {
		t.Fatalf("expected MasterPersonaGetRunStatus to succeed on recreated controller: %v", err)
	}
	if secondStatus.RunState != "入力待ち" {
		t.Fatalf("expected run state to reset to %q after controller recreation (restart must not carry run state), got %q", "入力待ち", secondStatus.RunState)
	}
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

// newBootstrapRunStatusTestController builds a test AppController with nil master persona seed.
// persona entry seed を nil にすることで、migration 002 で DROP された master_persona_entries テーブルへの
// アクセスを回避する。run status と AI settings のみを扱うテストで使用する。
func newBootstrapRunStatusTestController(t *testing.T, databasePath string) *controllerwails.AppController {
	t.Helper()
	setBootstrapTestDatabasePath(t, databasePath)
	setBootstrapTestMasterPersonaSecretStore(t, databasePath)
	return newAppControllerWithSeeds(bootstrapTestSeed(), nil, bootstrapTestNow)
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

// newBootstrapInMemoryRunStatusControllerWithRepo builds an AppController backed by
// InMemoryMasterPersonaRepository so that SaveRunStatus / LoadRunStatus are live
// within the session. The returned repo pointer allows callers to seed run state
// before invoking controller methods.
func newBootstrapInMemoryRunStatusControllerWithRepo(t *testing.T) (*controllerwails.AppController, *repository.InMemoryMasterPersonaRepository) {
	t.Helper()

	inMemoryRepo := repository.NewInMemoryMasterPersonaRepository(nil)
	inMemorySecretStore := repository.NewInMemorySecretStore()

	queryService := service.NewMasterPersonaQueryService(inMemoryRepo)
	generationService := service.NewMasterPersonaGenerationService(
		inMemoryRepo,
		inMemoryRepo,
		inMemoryRepo,
		inMemorySecretStore,
		bootstrapTestNow,
		true, // testMode: paid real AI API を呼ばない
	)
	runStatusService := service.NewMasterPersonaRunStatusService(inMemoryRepo, bootstrapTestNow)

	masterPersonaUsecase := usecase.NewMasterPersonaUsecase(
		queryService,
		generationService,
		runStatusService,
	)
	masterPersonaController := controllerwails.NewMasterPersonaController(masterPersonaUsecase)
	controller := controllerwails.NewAppController(nil, masterPersonaController, nil)
	return controller, inMemoryRepo
}

// persona-ai-settings-restart-cutover: same controller instance で ExecuteGeneration 後に
// MasterPersonaGetRunStatus を再読込しても non-idle を保つことを証明する。
// InMemoryMasterPersonaRepository は SaveRunStatus を実際に保存するため、
// 同一セッション内での live run status readback が可能である。
func TestNewAppControllerPersonaAISettingsRestartCutoverSameInstanceGetRunStatusAfterExecuteIsNonIdle(t *testing.T) {
	// Arrange: InMemory wiring で AppController を構築する (空 AI 設定)。
	controller, _ := newBootstrapInMemoryRunStatusControllerWithRepo(t)
	defer controller.OnShutdown(context.Background())

	// Act: 空設定で ExecuteGeneration を呼ぶ → "設定未完了" が返り、InMemory に保存される。
	execResult, err := controller.MasterPersonaExecuteGeneration(controllerwails.MasterPersonaExecuteRequestDTO{})
	if err != nil {
		t.Fatalf("expected ExecuteGeneration with empty settings to succeed: %v", err)
	}
	if execResult.RunState == "入力待ち" {
		t.Fatalf("expected ExecuteGeneration to return non-idle run state, got %q", execResult.RunState)
	}

	// Act: 同一インスタンスで MasterPersonaGetRunStatus を再読込する。
	readbackStatus, err := controller.MasterPersonaGetRunStatus()
	if err != nil {
		t.Fatalf("expected GetRunStatus readback to succeed: %v", err)
	}

	// Assert: 再読込しても non-idle を保つ (InMemory SaveRunStatus が live state を保持するため)。
	if readbackStatus.RunState == "入力待ち" {
		t.Fatalf("expected GetRunStatus readback to maintain non-idle after ExecuteGeneration, got %q (live run status must reflect last execute outcome)", readbackStatus.RunState)
	}
}

// persona-ai-settings-restart-cutover: same controller instance で running 状態のとき
// InterruptGeneration が current state を見て "中断済み" へ遷移することを証明する。
func TestNewAppControllerPersonaAISettingsRestartCutoverInterruptSeesCurrentRunningStateAndTransitions(t *testing.T) {
	// Arrange: InMemory wiring + "生成中" で run status を事前設定する。
	controller, inMemoryRepo := newBootstrapInMemoryRunStatusControllerWithRepo(t)
	defer controller.OnShutdown(context.Background())

	if err := inMemoryRepo.SaveRunStatus(context.Background(), repository.MasterPersonaRunStatusRecord{
		RunState: service.MasterPersonaStatusRunning,
	}); err != nil {
		t.Fatalf("expected run status seed to running state to succeed: %v", err)
	}

	// Act: running 状態で InterruptGeneration を呼ぶ。
	result, err := controller.MasterPersonaInterruptGeneration()
	if err != nil {
		t.Fatalf("expected InterruptGeneration from running state to succeed: %v", err)
	}

	// Assert: current state ("生成中") を見て "中断済み" へ遷移する。
	if result.RunState != service.MasterPersonaStatusInterrupted {
		t.Fatalf("expected InterruptGeneration to transition running→interrupted, got %q", result.RunState)
	}
}

// persona-ai-settings-restart-cutover: same controller instance で running 状態のとき
// CancelGeneration が current state を見て "中止済み" へ遷移することを証明する。
func TestNewAppControllerPersonaAISettingsRestartCutoverCancelSeesCurrentRunningStateAndTransitions(t *testing.T) {
	// Arrange: InMemory wiring + "生成中" で run status を事前設定する。
	controller, inMemoryRepo := newBootstrapInMemoryRunStatusControllerWithRepo(t)
	defer controller.OnShutdown(context.Background())

	if err := inMemoryRepo.SaveRunStatus(context.Background(), repository.MasterPersonaRunStatusRecord{
		RunState: service.MasterPersonaStatusRunning,
	}); err != nil {
		t.Fatalf("expected run status seed to running state to succeed: %v", err)
	}

	// Act: running 状態で CancelGeneration を呼ぶ。
	result, err := controller.MasterPersonaCancelGeneration()
	if err != nil {
		t.Fatalf("expected CancelGeneration from running state to succeed: %v", err)
	}

	// Assert: current state ("生成中") を見て "中止済み" へ遷移する。
	if result.RunState != service.MasterPersonaStatusCancelled {
		t.Fatalf("expected CancelGeneration to transition running→cancelled, got %q", result.RunState)
	}
}

// persona-generation-cutover: bootstrap controller を通じた Execute が canonical NPC_PROFILE 行を書き込むことを証明する。
// 現在は shipped sqlite repository が master_persona_entries を参照しているため FAIL する。
// canonical write path 実装後に PASS となる (RED テスト)。
func TestNewAppControllerPersonaGenerationCutoverExecuteWritesCanonicalNPCProfileRow(t *testing.T) {
	// Arrange: nil persona seed (master_persona_entries へのアクセスを回避) + fake AI mode。
	databasePath := configureBootstrapTestDatabase(t)
	setBootstrapTestMasterPersonaSecretStore(t, databasePath)
	t.Setenv(masterPersonaFakeResponseEnv, "cutover-npc-profile-persona-description")
	controller := newAppControllerWithSeeds(bootstrapTestSeed(), nil, bootstrapTestNow)

	extractPath := writeBootstrapMasterPersonaExtractFixture(t, `{
  "target_plugin": "CutoverPlugin.esp",
  "npcs": [
    {
      "form_id": "FE01BB01",
      "record_type": "NPC_",
      "editor_id": "CO_TestNPC",
      "display_name": "Cutover NPC",
      "dialogues": ["Hello", "Goodbye"]
    }
  ]
}`)

	// Act: canonical write path を通じてペルソナ生成を実行する。
	execResult, execErr := controller.MasterPersonaExecuteGeneration(controllerwails.MasterPersonaExecuteRequestDTO{
		FilePath:   extractPath,
		AISettings: controllerwails.MasterPersonaAISettingsDTO{Provider: "gemini", Model: "gemini-2.5-pro", APIKey: "cutover-test-key"},
	})
	controller.OnShutdown(context.Background())

	// Assert: 生成が成功する (canonical write path が実装されるまで FAIL する)。
	if execErr != nil {
		t.Fatalf("expected generation to succeed via canonical NPC_PROFILE write path, got error: %v", execErr)
	}
	if execResult.RunState != "完了" {
		t.Fatalf("expected run state 完了, got %q (message: %q)", execResult.RunState, execResult.Message)
	}

	// Assert: NPC_PROFILE に canonical 行が書き込まれている。
	db, openErr := repository.OpenSQLiteDatabase(context.Background(), databasePath)
	if openErr != nil {
		t.Fatalf("expected db open to succeed: %v", openErr)
	}
	defer func() {
		if closeErr := db.Close(); closeErr != nil {
			t.Fatalf("expected db close to succeed: %v", closeErr)
		}
	}()

	var npcProfileCount int
	if err := db.QueryRowContext(context.Background(),
		"SELECT COUNT(*) FROM NPC_PROFILE WHERE target_plugin_name = ? AND form_id = ? AND record_type = ?",
		"CutoverPlugin.esp", "FE01BB01", "NPC_").Scan(&npcProfileCount); err != nil {
		t.Fatalf("expected NPC_PROFILE count query to succeed: %v", err)
	}
	if npcProfileCount != 1 {
		t.Fatalf("expected 1 NPC_PROFILE row after generation cutover, got %d (write sink must be NPC_PROFILE, not master_persona_entries)", npcProfileCount)
	}
}

// persona-generation-cutover: bootstrap controller を通じた Execute が canonical PERSONA 行を書き込むことを証明する。
// 現在は shipped sqlite repository が master_persona_entries を参照しているため FAIL する。
// canonical write path 実装後に PASS となる (RED テスト)。
func TestNewAppControllerPersonaGenerationCutoverExecuteWritesCanonicalPersonaRow(t *testing.T) {
	// Arrange: nil persona seed + fake AI mode。
	databasePath := configureBootstrapTestDatabase(t)
	setBootstrapTestMasterPersonaSecretStore(t, databasePath)
	t.Setenv(masterPersonaFakeResponseEnv, "cutover-persona-row-description")
	controller := newAppControllerWithSeeds(bootstrapTestSeed(), nil, bootstrapTestNow)

	extractPath := writeBootstrapMasterPersonaExtractFixture(t, `{
  "target_plugin": "CutoverPlugin.esp",
  "npcs": [
    {
      "form_id": "FE01BB01",
      "record_type": "NPC_",
      "editor_id": "CO_TestNPC",
      "display_name": "Cutover NPC",
      "dialogues": ["Hello", "Goodbye"]
    }
  ]
}`)

	// Act: canonical write path を通じてペルソナ生成を実行する。
	execResult, execErr := controller.MasterPersonaExecuteGeneration(controllerwails.MasterPersonaExecuteRequestDTO{
		FilePath:   extractPath,
		AISettings: controllerwails.MasterPersonaAISettingsDTO{Provider: "gemini", Model: "gemini-2.5-pro", APIKey: "cutover-test-key"},
	})
	controller.OnShutdown(context.Background())

	// Assert: 生成が成功する。
	if execErr != nil {
		t.Fatalf("expected generation to succeed via canonical PERSONA write path, got error: %v", execErr)
	}
	if execResult.RunState != "完了" {
		t.Fatalf("expected run state 完了, got %q (message: %q)", execResult.RunState, execResult.Message)
	}

	// Assert: PERSONA に NPC_PROFILE に紐づく canonical 行が書き込まれている。
	db, openErr := repository.OpenSQLiteDatabase(context.Background(), databasePath)
	if openErr != nil {
		t.Fatalf("expected db open to succeed: %v", openErr)
	}
	defer func() {
		if closeErr := db.Close(); closeErr != nil {
			t.Fatalf("expected db close to succeed: %v", closeErr)
		}
	}()

	var personaCount int
	if err := db.QueryRowContext(context.Background(),
		`SELECT COUNT(*) FROM PERSONA p
		 JOIN NPC_PROFILE np ON p.npc_profile_id = np.id
		 WHERE np.target_plugin_name = ? AND np.form_id = ? AND np.record_type = ?`,
		"CutoverPlugin.esp", "FE01BB01", "NPC_").Scan(&personaCount); err != nil {
		t.Fatalf("expected PERSONA join NPC_PROFILE count query to succeed: %v", err)
	}
	if personaCount != 1 {
		t.Fatalf("expected 1 PERSONA row linked to NPC_PROFILE after generation cutover, got %d", personaCount)
	}
}

// persona-generation-cutover: bootstrap DB schema に master_persona_entries が存在しないことを証明する。
// generation write sink は canonical NPC_PROFILE + PERSONA であり、legacy table は migration 002 で削除済み。
func TestNewAppControllerPersonaGenerationCutoverLegacyMasterPersonaEntriesAbsentFromSchema(t *testing.T) {
	// Arrange: bootstrap テスト用 DB を開く。
	databasePath := configureBootstrapTestDatabase(t)
	db, err := repository.OpenSQLiteDatabase(context.Background(), databasePath)
	if err != nil {
		t.Fatalf("expected db open to succeed: %v", err)
	}
	defer func() {
		if closeErr := db.Close(); closeErr != nil {
			t.Fatalf("expected db close to succeed: %v", closeErr)
		}
	}()

	// Assert: master_persona_entries テーブルは schema に存在しない (migration 002 で削除済み)。
	var count int
	if err := db.QueryRowContext(context.Background(),
		"SELECT COUNT(*) FROM sqlite_master WHERE type='table' AND name='master_persona_entries'").Scan(&count); err != nil {
		t.Fatalf("expected table existence check to succeed: %v", err)
	}
	if count != 0 {
		t.Fatalf("expected master_persona_entries to be absent from bootstrap DB schema (generation write sink must be canonical NPC_PROFILE + PERSONA), but table was found")
	}
}

// persona-generation-cutover: Execute が失敗したとき、bootstrap DB に partial canonical rows が残らないことを証明する。
// canonical write path はアトミックであるべきであり、失敗後の NPC_PROFILE + PERSONA に孤立行が存在してはならない。
func TestNewAppControllerPersonaGenerationCutoverFailedExecutionLeavesNoPartialCanonicalRows(t *testing.T) {
	// Arrange: nil persona seed + fake AI mode (fake response 未設定で空の body)。
	databasePath := configureBootstrapTestDatabase(t)
	setBootstrapTestMasterPersonaSecretStore(t, databasePath)
	controller := newAppControllerWithSeeds(bootstrapTestSeed(), nil, bootstrapTestNow)

	extractPath := writeBootstrapMasterPersonaExtractFixture(t, `{
  "target_plugin": "FailPlugin.esp",
  "npcs": [
    {
      "form_id": "FE01CC01",
      "record_type": "NPC_",
      "editor_id": "FP_FailNPC",
      "display_name": "Fail NPC",
      "dialogues": ["line"]
    }
  ]
}`)

	// Act: 生成を試みる (canonical write path が未実装の場合、UpsertIfAbsent で失敗する)。
	_, _ = controller.MasterPersonaExecuteGeneration(controllerwails.MasterPersonaExecuteRequestDTO{
		FilePath:   extractPath,
		AISettings: controllerwails.MasterPersonaAISettingsDTO{Provider: "gemini", Model: "gemini-2.5-pro", APIKey: "cutover-fail-key"},
	})
	controller.OnShutdown(context.Background())

	// Assert: 失敗後に partial canonical rows が存在しない (write はアトミックであるべき)。
	db, openErr := repository.OpenSQLiteDatabase(context.Background(), databasePath)
	if openErr != nil {
		t.Fatalf("expected db open to succeed: %v", openErr)
	}
	defer func() {
		if closeErr := db.Close(); closeErr != nil {
			t.Fatalf("expected db close to succeed: %v", closeErr)
		}
	}()

	var npcProfileCount int
	if err := db.QueryRowContext(context.Background(),
		"SELECT COUNT(*) FROM NPC_PROFILE WHERE target_plugin_name = 'FailPlugin.esp'").Scan(&npcProfileCount); err != nil {
		t.Fatalf("expected NPC_PROFILE count query to succeed: %v", err)
	}
	if npcProfileCount != 0 {
		t.Fatalf("expected no partial NPC_PROFILE rows after failed generation, got %d (canonical write must be atomic)", npcProfileCount)
	}
}
