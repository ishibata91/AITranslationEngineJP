package usecase

import (
	"context"
	"errors"
	"reflect"
	"testing"
	"time"

	"aitranslationenginejp/internal/service"
)

type fakeMasterPersonaQueryService struct {
	searchEntriesFunc    func(ctx context.Context, query service.MasterPersonaListQuery) (service.MasterPersonaListResult, error)
	loadEntryDetailFunc  func(ctx context.Context, identityKey string) (service.MasterPersonaEntry, error)
	loadDialogueListFunc func(ctx context.Context, identityKey string) (service.MasterPersonaDialogueList, error)
}

func (fake fakeMasterPersonaQueryService) SearchEntries(ctx context.Context, query service.MasterPersonaListQuery) (service.MasterPersonaListResult, error) {
	if fake.searchEntriesFunc == nil {
		return service.MasterPersonaListResult{}, nil
	}
	return fake.searchEntriesFunc(ctx, query)
}
func (fake fakeMasterPersonaQueryService) LoadEntryDetail(ctx context.Context, identityKey string) (service.MasterPersonaEntry, error) {
	if fake.loadEntryDetailFunc == nil {
		return service.MasterPersonaEntry{}, nil
	}
	return fake.loadEntryDetailFunc(ctx, identityKey)
}
func (fake fakeMasterPersonaQueryService) LoadDialogueList(ctx context.Context, identityKey string) (service.MasterPersonaDialogueList, error) {
	if fake.loadDialogueListFunc == nil {
		return service.MasterPersonaDialogueList{}, nil
	}
	return fake.loadDialogueListFunc(ctx, identityKey)
}

type fakeMasterPersonaGenerationService struct {
	previewFunc     func(ctx context.Context, filePath string, requestSettings service.MasterPersonaAISettings) (service.MasterPersonaPreviewResult, error)
	executeFunc     func(ctx context.Context, filePath string, requestSettings service.MasterPersonaAISettings) (service.MasterPersonaRunStatus, error)
	updateEntryFunc func(ctx context.Context, identityKey string, input service.MasterPersonaUpdateInput) (service.MasterPersonaEntry, error)
	deleteEntryFunc func(ctx context.Context, identityKey string) error
}

func (fake fakeMasterPersonaGenerationService) LoadSettings(_ context.Context) (service.MasterPersonaAISettings, error) {
	return service.MasterPersonaAISettings{}, nil
}
func (fake fakeMasterPersonaGenerationService) SaveSettings(_ context.Context, settings service.MasterPersonaAISettings) (service.MasterPersonaAISettings, error) {
	return settings, nil
}
func (fake fakeMasterPersonaGenerationService) Preview(ctx context.Context, filePath string, requestSettings service.MasterPersonaAISettings) (service.MasterPersonaPreviewResult, error) {
	if fake.previewFunc == nil {
		return service.MasterPersonaPreviewResult{}, nil
	}
	return fake.previewFunc(ctx, filePath, requestSettings)
}
func (fake fakeMasterPersonaGenerationService) Execute(ctx context.Context, filePath string, requestSettings service.MasterPersonaAISettings) (service.MasterPersonaRunStatus, error) {
	if fake.executeFunc == nil {
		return service.MasterPersonaRunStatus{}, nil
	}
	return fake.executeFunc(ctx, filePath, requestSettings)
}
func (fake fakeMasterPersonaGenerationService) UpdateEntry(ctx context.Context, identityKey string, input service.MasterPersonaUpdateInput) (service.MasterPersonaEntry, error) {
	if fake.updateEntryFunc == nil {
		return service.MasterPersonaEntry{}, nil
	}
	return fake.updateEntryFunc(ctx, identityKey, input)
}
func (fake fakeMasterPersonaGenerationService) DeleteEntry(ctx context.Context, identityKey string) error {
	if fake.deleteEntryFunc == nil {
		return nil
	}
	return fake.deleteEntryFunc(ctx, identityKey)
}

type fakeMasterPersonaGenerationServiceWithLoad struct {
	loadSettingsFunc func(ctx context.Context) (service.MasterPersonaAISettings, error)
}

func (fake fakeMasterPersonaGenerationServiceWithLoad) LoadSettings(ctx context.Context) (service.MasterPersonaAISettings, error) {
	if fake.loadSettingsFunc == nil {
		return service.MasterPersonaAISettings{}, nil
	}
	return fake.loadSettingsFunc(ctx)
}
func (fake fakeMasterPersonaGenerationServiceWithLoad) SaveSettings(_ context.Context, settings service.MasterPersonaAISettings) (service.MasterPersonaAISettings, error) {
	return settings, nil
}
func (fake fakeMasterPersonaGenerationServiceWithLoad) Preview(_ context.Context, _ string, _ service.MasterPersonaAISettings) (service.MasterPersonaPreviewResult, error) {
	return service.MasterPersonaPreviewResult{}, nil
}
func (fake fakeMasterPersonaGenerationServiceWithLoad) Execute(_ context.Context, _ string, _ service.MasterPersonaAISettings) (service.MasterPersonaRunStatus, error) {
	return service.MasterPersonaRunStatus{}, nil
}
func (fake fakeMasterPersonaGenerationServiceWithLoad) UpdateEntry(_ context.Context, _ string, _ service.MasterPersonaUpdateInput) (service.MasterPersonaEntry, error) {
	return service.MasterPersonaEntry{}, nil
}
func (fake fakeMasterPersonaGenerationServiceWithLoad) DeleteEntry(_ context.Context, _ string) error {
	return nil
}

type fakeMasterPersonaRunStatusService struct{}

func (fakeMasterPersonaRunStatusService) GetStatus(_ context.Context) (service.MasterPersonaRunStatus, error) {
	return service.MasterPersonaRunStatus{}, nil
}
func (fakeMasterPersonaRunStatusService) Interrupt(_ context.Context) (service.MasterPersonaRunStatus, error) {
	return service.MasterPersonaRunStatus{}, nil
}
func (fakeMasterPersonaRunStatusService) Cancel(_ context.Context) (service.MasterPersonaRunStatus, error) {
	return service.MasterPersonaRunStatus{}, nil
}

func TestMasterPersonaUsecaseGetPageSelectsPreferredIdentityKey(t *testing.T) {
	preferred := "FollowersPlus.esp:FE01A812:NPC_"
	usecase := NewMasterPersonaUsecase(
		fakeMasterPersonaQueryService{
			searchEntriesFunc: func(_ context.Context, query service.MasterPersonaListQuery) (service.MasterPersonaListResult, error) {
				if query.PluginFilter != "FollowersPlus.esp" {
					t.Fatalf("unexpected query: %#v", query)
				}
				return service.MasterPersonaListResult{Items: []service.MasterPersonaEntry{{IdentityKey: preferred}, {IdentityKey: "other"}}, TotalCount: 2, Page: 1, PageSize: 30}, nil
			},
		},
		fakeMasterPersonaGenerationService{},
		fakeMasterPersonaRunStatusService{},
	)

	page, err := usecase.GetPage(context.Background(), MasterPersonaListQuery{PluginFilter: "FollowersPlus.esp"}, &preferred)
	if err != nil {
		t.Fatalf("expected get page to succeed: %v", err)
	}
	if page.SelectedIdentityKey == nil || *page.SelectedIdentityKey != preferred {
		t.Fatalf("unexpected selected identity key: %#v", page.SelectedIdentityKey)
	}
}

func TestMasterPersonaUsecasePreviewGenerationPassesThroughProviderRequest(t *testing.T) {
	usecase := NewMasterPersonaUsecase(
		fakeMasterPersonaQueryService{},
		fakeMasterPersonaGenerationService{
			previewFunc: func(_ context.Context, filePath string, requestSettings service.MasterPersonaAISettings) (service.MasterPersonaPreviewResult, error) {
				if filePath != "/tmp/sample.json" {
					t.Fatalf("unexpected file path: %q", filePath)
				}
				if requestSettings.Provider != "gemini" || requestSettings.APIKey != "" {
					t.Fatalf("unexpected preview settings: %#v", requestSettings)
				}
				return service.MasterPersonaPreviewResult{Status: service.MasterPersonaStatusReady}, nil
			},
		},
		fakeMasterPersonaRunStatusService{},
	)

	result, err := usecase.PreviewGeneration(context.Background(), "/tmp/sample.json", MasterPersonaAISettings{Provider: "gemini", Model: "gemini-2.5-pro", APIKey: ""})
	if err != nil {
		t.Fatalf("expected preview generation to succeed: %v", err)
	}
	if result.Status != service.MasterPersonaStatusReady {
		t.Fatalf("unexpected preview result: %#v", result)
	}
}

func TestMasterPersonaUsecaseDeleteEntryWrapsGenerationError(t *testing.T) {
	usecase := NewMasterPersonaUsecase(
		fakeMasterPersonaQueryService{},
		fakeMasterPersonaGenerationService{deleteEntryFunc: func(_ context.Context, _ string) error {
			return errors.New("delete failed")
		}},
		fakeMasterPersonaRunStatusService{},
	)

	_, err := usecase.DeleteEntry(context.Background(), "key", MasterPersonaListQuery{})
	if err == nil {
		t.Fatal("expected delete error")
	}
}

func TestMasterPersonaUsecaseGetDetailForwardsIdentityKey(t *testing.T) {
	const wantKey = "FollowersPlus.esp:FE01A812:NPC_"
	uc := NewMasterPersonaUsecase(
		fakeMasterPersonaQueryService{
			loadEntryDetailFunc: func(_ context.Context, identityKey string) (service.MasterPersonaEntry, error) {
				if identityKey != wantKey {
					t.Fatalf("unexpected identity key: %q", identityKey)
				}
				return service.MasterPersonaEntry{IdentityKey: wantKey, DisplayName: "Lys Maren"}, nil
			},
		},
		fakeMasterPersonaGenerationService{},
		fakeMasterPersonaRunStatusService{},
	)

	got, err := uc.GetDetail(context.Background(), wantKey)
	if err != nil {
		t.Fatalf("expected get detail to succeed: %v", err)
	}
	if got.IdentityKey != wantKey || got.DisplayName != "Lys Maren" {
		t.Fatalf("unexpected entry: %#v", got)
	}
}

func TestMasterPersonaUsecaseGetDetailPropagatesQueryServiceError(t *testing.T) {
	uc := NewMasterPersonaUsecase(
		fakeMasterPersonaQueryService{
			loadEntryDetailFunc: func(_ context.Context, _ string) (service.MasterPersonaEntry, error) {
				return service.MasterPersonaEntry{}, errors.New("db connection lost")
			},
		},
		fakeMasterPersonaGenerationService{},
		fakeMasterPersonaRunStatusService{},
	)

	_, err := uc.GetDetail(context.Background(), "FollowersPlus.esp:FE01A812:NPC_")
	if err == nil {
		t.Fatal("expected get detail to fail")
	}
}

func TestMasterPersonaUsecaseGetPageForwardsKeyword(t *testing.T) {
	uc := NewMasterPersonaUsecase(
		fakeMasterPersonaQueryService{
			searchEntriesFunc: func(_ context.Context, query service.MasterPersonaListQuery) (service.MasterPersonaListResult, error) {
				if query.Keyword != "lys" {
					t.Fatalf("unexpected keyword: %q", query.Keyword)
				}
				return service.MasterPersonaListResult{}, nil
			},
		},
		fakeMasterPersonaGenerationService{},
		fakeMasterPersonaRunStatusService{},
	)

	_, err := uc.GetPage(context.Background(), MasterPersonaListQuery{Keyword: "lys"}, nil)
	if err != nil {
		t.Fatalf("expected get page to succeed: %v", err)
	}
}

// persona-read-detail-cutover: GetDetail が identity snapshot / NPC_PROFILE join fields を転送することを証明する。
func TestMasterPersonaUsecasePersonaReadDetailCutoverGetDetailForwardsCanonicalFields(t *testing.T) {
	race := "Nord"
	sex := "Male"
	uc := NewMasterPersonaUsecase(
		fakeMasterPersonaQueryService{
			loadEntryDetailFunc: func(_ context.Context, identityKey string) (service.MasterPersonaEntry, error) {
				return service.MasterPersonaEntry{
					IdentityKey:  identityKey,
					FormID:       "FE01A813",
					RecordType:   "NPC_",
					EditorID:     "FP_KaelRuun",
					DisplayName:  "Kael Ruun",
					Race:         &race,
					Sex:          &sex,
					VoiceType:    "MaleCommander",
					ClassName:    "FPMercenaryClass",
					SourcePlugin: "FollowersPlus.esp",
					PersonaBody:  "判断を先に述べ、必要な指示だけを短く渡す。",
					UpdatedAt:    time.Date(2026, 4, 15, 0, 0, 0, 0, time.UTC),
				}, nil
			},
		},
		fakeMasterPersonaGenerationService{},
		fakeMasterPersonaRunStatusService{},
	)

	got, err := uc.GetDetail(context.Background(), "FollowersPlus.esp:FE01A813:NPC_")
	if err != nil {
		t.Fatalf("expected get detail to succeed: %v", err)
	}
	if got.EditorID != "FP_KaelRuun" {
		t.Fatalf("expected EditorID forwarded, got %q", got.EditorID)
	}
	if got.VoiceType != "MaleCommander" {
		t.Fatalf("expected VoiceType forwarded, got %q", got.VoiceType)
	}
	if got.SourcePlugin != "FollowersPlus.esp" {
		t.Fatalf("expected SourcePlugin forwarded, got %q", got.SourcePlugin)
	}
	if got.Race == nil || *got.Race != "Nord" {
		t.Fatalf("expected Race forwarded, got %#v", got.Race)
	}
	if got.Sex == nil || *got.Sex != "Male" {
		t.Fatalf("expected Sex forwarded, got %#v", got.Sex)
	}
	if got.PersonaBody != "判断を先に述べ、必要な指示だけを短く渡す。" {
		t.Fatalf("expected PersonaBody forwarded, got %q", got.PersonaBody)
	}
}

// persona-read-detail-cutover: MasterPersonaUsecase の public seam から GetDialogueList が除去されることを証明する。
// MasterPersonaUsecase.GetDialogueList が残っている間は失敗する。
func TestMasterPersonaUsecasePersonaReadDetailCutoverHasNoGetDialogueList(t *testing.T) {
	ucType := reflect.TypeOf(&MasterPersonaUsecase{})
	for i := 0; i < ucType.NumMethod(); i++ {
		if ucType.Method(i).Name == "GetDialogueList" {
			t.Fatal("MasterPersonaUsecase still exposes GetDialogueList; persona-read-detail-cutover requires removal from public seam")
		}
	}
}

// persona-read-detail-cutover: MasterPersonaQueryServicePort から LoadDialogueList が除去されることを証明する。
// MasterPersonaQueryServicePort インターフェースが LoadDialogueList を保持している間は失敗する。
func TestMasterPersonaUsecasePersonaReadDetailCutoverQueryServicePortHasNoLoadDialogueList(t *testing.T) {
	portType := reflect.TypeOf((*MasterPersonaQueryServicePort)(nil)).Elem()
	for i := 0; i < portType.NumMethod(); i++ {
		if portType.Method(i).Name == "LoadDialogueList" {
			t.Fatal("MasterPersonaQueryServicePort still defines LoadDialogueList; persona-read-detail-cutover requires removal from query service port seam")
		}
	}
}

// persona-ai-settings-restart-cutover: LoadAISettings は generation service から復元済み provider と model を返すことを証明する。
// persona-ai-settings-restart-cutover: LoadAISettings は generation service が返した provider と model をそのまま呼び出し元へ転送することを証明する。
func TestMasterPersonaUsecasePersonaAISettingsRestartCutoverLoadAISettingsReturnsRestoredSettingsFromService(t *testing.T) {
	uc := NewMasterPersonaUsecase(
		fakeMasterPersonaQueryService{},
		fakeMasterPersonaGenerationServiceWithLoad{
			loadSettingsFunc: func(_ context.Context) (service.MasterPersonaAISettings, error) {
				return service.MasterPersonaAISettings{Provider: "gemini", Model: "restored-model", APIKey: "restored-key"}, nil
			},
		},
		fakeMasterPersonaRunStatusService{},
	)

	settings, err := uc.LoadAISettings(context.Background())
	if err != nil {
		t.Fatalf("expected load ai settings to succeed: %v", err)
	}
	if settings.Provider != "gemini" || settings.Model != "restored-model" {
		t.Fatalf("expected restored provider/model forwarded, got %#v", settings)
	}
	if settings.APIKey != "restored-key" {
		t.Fatalf("expected restored api key forwarded, got %q", settings.APIKey)
	}
}

// persona-edit-delete-cutover: UpdateEntry が ChangedEntry を含む結果を返すことを証明する。
// usecase が ChangedEntry を返さない場合は失敗する。
func TestMasterPersonaUsecasePersonaEditDeleteCutoverUpdateEntryReturnsChangedEntry(t *testing.T) {
	const wantKey = "FollowersPlus.esp:FE01A812:NPC_"
	uc := NewMasterPersonaUsecase(
		fakeMasterPersonaQueryService{},
		fakeMasterPersonaGenerationService{
			updateEntryFunc: func(_ context.Context, identityKey string, input service.MasterPersonaUpdateInput) (service.MasterPersonaEntry, error) {
				return service.MasterPersonaEntry{
					IdentityKey: identityKey,
					PersonaBody: input.PersonaBody,
				}, nil
			},
		},
		fakeMasterPersonaRunStatusService{},
	)

	result, err := uc.UpdateEntry(context.Background(), wantKey, MasterPersonaUpdateInput{PersonaBody: "updated body"}, MasterPersonaListQuery{})

	if err != nil {
		t.Fatalf("expected update entry to succeed: %v", err)
	}
	if result.ChangedEntry == nil {
		t.Fatal("expected ChangedEntry to be populated in mutation result")
	}
	if result.ChangedEntry.IdentityKey != wantKey {
		t.Fatalf("expected ChangedEntry.IdentityKey=%q, got %q", wantKey, result.ChangedEntry.IdentityKey)
	}
	if result.ChangedEntry.PersonaBody != "updated body" {
		t.Fatalf("expected ChangedEntry.PersonaBody from input, got %q", result.ChangedEntry.PersonaBody)
	}
}

// persona-edit-delete-cutover: DeleteEntry が DeletedEntryID を含む結果を返すことを証明する。
// usecase が DeletedEntryID を返さない場合は失敗する。
func TestMasterPersonaUsecasePersonaEditDeleteCutoverDeleteEntryReturnsDeletedEntryID(t *testing.T) {
	const wantKey = "FollowersPlus.esp:FE01A812:NPC_"
	uc := NewMasterPersonaUsecase(
		fakeMasterPersonaQueryService{},
		fakeMasterPersonaGenerationService{},
		fakeMasterPersonaRunStatusService{},
	)

	result, err := uc.DeleteEntry(context.Background(), wantKey, MasterPersonaListQuery{})

	if err != nil {
		t.Fatalf("expected delete entry to succeed: %v", err)
	}
	if result.DeletedEntryID == nil {
		t.Fatal("expected DeletedEntryID to be populated in mutation result")
	}
	if *result.DeletedEntryID != wantKey {
		t.Fatalf("expected DeletedEntryID=%q, got %q", wantKey, *result.DeletedEntryID)
	}
}

// persona-edit-delete-cutover: UpdateEntry が service エラーをラップして返すことを証明する。
// usecase がエラーを握り潰す場合は失敗する。
func TestMasterPersonaUsecasePersonaEditDeleteCutoverUpdateEntryWrapsServiceError(t *testing.T) {
	uc := NewMasterPersonaUsecase(
		fakeMasterPersonaQueryService{},
		fakeMasterPersonaGenerationService{
			updateEntryFunc: func(_ context.Context, _ string, _ service.MasterPersonaUpdateInput) (service.MasterPersonaEntry, error) {
				return service.MasterPersonaEntry{}, errors.New("update failed in service")
			},
		},
		fakeMasterPersonaRunStatusService{},
	)

	_, err := uc.UpdateEntry(context.Background(), "FollowersPlus.esp:FE01A812:NPC_", MasterPersonaUpdateInput{}, MasterPersonaListQuery{})

	if err == nil {
		t.Fatal("expected update entry to propagate service error")
	}
}
