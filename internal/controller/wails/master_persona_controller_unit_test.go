package wails

import (
	"context"
	"encoding/json"
	"errors"
	"reflect"
	"strings"
	"testing"
	"time"

	"aitranslationenginejp/internal/usecase"
)

type fakeMasterPersonaUsecase struct {
	getPageFunc             func(ctx context.Context, query usecase.MasterPersonaListQuery, preferredIdentityKey *string) (usecase.MasterPersonaPageState, error)
	getDetailFunc           func(ctx context.Context, identityKey string) (usecase.MasterPersonaEntry, error)
	loadAISettingsFunc      func(ctx context.Context) (usecase.MasterPersonaAISettings, error)
	saveAISettingsFunc      func(ctx context.Context, settings usecase.MasterPersonaAISettings) (usecase.MasterPersonaAISettings, error)
	previewGenerationFunc   func(ctx context.Context, filePath string, settings usecase.MasterPersonaAISettings) (usecase.MasterPersonaPreviewResult, error)
	executeGenerationFunc   func(ctx context.Context, filePath string, settings usecase.MasterPersonaAISettings) (usecase.MasterPersonaRunStatus, error)
	getRunStatusFunc        func(ctx context.Context) (usecase.MasterPersonaRunStatus, error)
	interruptGenerationFunc func(ctx context.Context) (usecase.MasterPersonaRunStatus, error)
	cancelGenerationFunc    func(ctx context.Context) (usecase.MasterPersonaRunStatus, error)
	updateEntryFunc         func(ctx context.Context, identityKey string, input usecase.MasterPersonaUpdateInput, refreshQuery usecase.MasterPersonaListQuery) (usecase.MasterPersonaMutationResult, error)
	deleteEntryFunc         func(ctx context.Context, identityKey string, refreshQuery usecase.MasterPersonaListQuery) (usecase.MasterPersonaMutationResult, error)
}

func (fake fakeMasterPersonaUsecase) GetPage(ctx context.Context, query usecase.MasterPersonaListQuery, preferred *string) (usecase.MasterPersonaPageState, error) {
	if fake.getPageFunc == nil {
		return usecase.MasterPersonaPageState{}, nil
	}
	return fake.getPageFunc(ctx, query, preferred)
}
func (fake fakeMasterPersonaUsecase) GetDetail(ctx context.Context, identityKey string) (usecase.MasterPersonaEntry, error) {
	if fake.getDetailFunc == nil {
		return usecase.MasterPersonaEntry{}, nil
	}
	return fake.getDetailFunc(ctx, identityKey)
}
func (fake fakeMasterPersonaUsecase) LoadAISettings(ctx context.Context) (usecase.MasterPersonaAISettings, error) {
	if fake.loadAISettingsFunc == nil {
		return usecase.MasterPersonaAISettings{}, nil
	}
	return fake.loadAISettingsFunc(ctx)
}
func (fake fakeMasterPersonaUsecase) SaveAISettings(ctx context.Context, settings usecase.MasterPersonaAISettings) (usecase.MasterPersonaAISettings, error) {
	if fake.saveAISettingsFunc == nil {
		return usecase.MasterPersonaAISettings{}, nil
	}
	return fake.saveAISettingsFunc(ctx, settings)
}
func (fake fakeMasterPersonaUsecase) PreviewGeneration(ctx context.Context, filePath string, settings usecase.MasterPersonaAISettings) (usecase.MasterPersonaPreviewResult, error) {
	if fake.previewGenerationFunc == nil {
		return usecase.MasterPersonaPreviewResult{}, nil
	}
	return fake.previewGenerationFunc(ctx, filePath, settings)
}
func (fake fakeMasterPersonaUsecase) ExecuteGeneration(ctx context.Context, filePath string, settings usecase.MasterPersonaAISettings) (usecase.MasterPersonaRunStatus, error) {
	if fake.executeGenerationFunc == nil {
		return usecase.MasterPersonaRunStatus{}, nil
	}
	return fake.executeGenerationFunc(ctx, filePath, settings)
}
func (fake fakeMasterPersonaUsecase) GetRunStatus(ctx context.Context) (usecase.MasterPersonaRunStatus, error) {
	if fake.getRunStatusFunc == nil {
		return usecase.MasterPersonaRunStatus{}, nil
	}
	return fake.getRunStatusFunc(ctx)
}
func (fake fakeMasterPersonaUsecase) InterruptGeneration(ctx context.Context) (usecase.MasterPersonaRunStatus, error) {
	if fake.interruptGenerationFunc == nil {
		return usecase.MasterPersonaRunStatus{}, nil
	}
	return fake.interruptGenerationFunc(ctx)
}
func (fake fakeMasterPersonaUsecase) CancelGeneration(ctx context.Context) (usecase.MasterPersonaRunStatus, error) {
	if fake.cancelGenerationFunc == nil {
		return usecase.MasterPersonaRunStatus{}, nil
	}
	return fake.cancelGenerationFunc(ctx)
}
func (fake fakeMasterPersonaUsecase) UpdateEntry(ctx context.Context, identityKey string, input usecase.MasterPersonaUpdateInput, refreshQuery usecase.MasterPersonaListQuery) (usecase.MasterPersonaMutationResult, error) {
	if fake.updateEntryFunc == nil {
		return usecase.MasterPersonaMutationResult{}, nil
	}
	return fake.updateEntryFunc(ctx, identityKey, input, refreshQuery)
}
func (fake fakeMasterPersonaUsecase) DeleteEntry(ctx context.Context, identityKey string, refreshQuery usecase.MasterPersonaListQuery) (usecase.MasterPersonaMutationResult, error) {
	if fake.deleteEntryFunc == nil {
		return usecase.MasterPersonaMutationResult{}, nil
	}
	return fake.deleteEntryFunc(ctx, identityKey, refreshQuery)
}

func TestMasterPersonaControllerGetPageMapsPluginFilter(t *testing.T) {
	controller := NewMasterPersonaController(fakeMasterPersonaUsecase{
		getPageFunc: func(_ context.Context, query usecase.MasterPersonaListQuery, preferred *string) (usecase.MasterPersonaPageState, error) {
			if query.PluginFilter != "FollowersPlus.esp" || query.Keyword != "lys" || query.Page != 2 || query.PageSize != 20 {
				t.Fatalf("unexpected query: %#v", query)
			}
			if preferred == nil || *preferred != "FollowersPlus.esp:FE01A812:NPC_" {
				t.Fatalf("unexpected preferred identity key: %#v", preferred)
			}
			return usecase.MasterPersonaPageState{}, nil
		},
	})

	preferred := "FollowersPlus.esp:FE01A812:NPC_"
	_, err := controller.MasterPersonaGetPage(MasterPersonaPageRequestDTO{Refresh: MasterPersonaListQueryDTO{Keyword: "lys", PluginFilter: "FollowersPlus.esp", Page: 2, PageSize: 20}, PreferredIdentityKey: &preferred})
	if err != nil {
		t.Fatalf("expected get page to succeed: %v", err)
	}
}

func TestMasterPersonaControllerGetDetailMapsRunLockReason(t *testing.T) {
	controller := NewMasterPersonaController(fakeMasterPersonaUsecase{
		getDetailFunc: func(_ context.Context, identityKey string) (usecase.MasterPersonaEntry, error) {
			if identityKey != "FollowersPlus.esp:FE01A812:NPC_" {
				t.Fatalf("unexpected identity key: %q", identityKey)
			}
			return usecase.MasterPersonaEntry{IdentityKey: identityKey, DisplayName: "Lys Maren", UpdatedAt: time.Date(2026, 4, 15, 0, 0, 0, 0, time.UTC)}, nil
		},
		getRunStatusFunc: func(_ context.Context) (usecase.MasterPersonaRunStatus, error) {
			return usecase.MasterPersonaRunStatus{RunState: "生成中"}, nil
		},
	})

	response, err := controller.MasterPersonaGetDetail(MasterPersonaDetailRequestDTO{IdentityKey: "FollowersPlus.esp:FE01A812:NPC_"})
	if err != nil {
		t.Fatalf("expected get detail to succeed: %v", err)
	}
	if response.Entry.RunLockReason != "更新と削除を行えません" {
		t.Fatalf("expected run lock reason, got %#v", response.Entry)
	}
}

func TestMasterPersonaControllerPreviewGenerationMapsResponse(t *testing.T) {
	controller := NewMasterPersonaController(fakeMasterPersonaUsecase{
		previewGenerationFunc: func(_ context.Context, filePath string, settings usecase.MasterPersonaAISettings) (usecase.MasterPersonaPreviewResult, error) {
			if filePath != "/tmp/sample.json" || settings.Provider != "gemini" || settings.Model != "gemini-2.5-pro" {
				t.Fatalf("unexpected preview request: path=%q settings=%#v", filePath, settings)
			}
			return usecase.MasterPersonaPreviewResult{FileName: "sample.json", TargetPlugin: "FollowersPlus.esp", TotalNPCCount: 4, GeneratableCount: 2, ExistingSkipCount: 1, ZeroDialogueSkipCount: 1, GenericNPCCount: 1, Status: "生成可能"}, nil
		},
	})

	response, err := controller.MasterPersonaPreviewGeneration(MasterPersonaPreviewRequestDTO{FilePath: "/tmp/sample.json", AISettings: MasterPersonaAISettingsDTO{Provider: "gemini", Model: "gemini-2.5-pro"}})
	if err != nil {
		t.Fatalf("expected preview generation to succeed: %v", err)
	}
	if response.ExistingSkipCount != 1 || response.ZeroDialogueSkipCount != 1 || response.GenericNPCCount != 1 {
		t.Fatalf("unexpected preview response: %#v", response)
	}
}

func TestMasterPersonaControllerUpdatePropagatesError(t *testing.T) {
	controller := NewMasterPersonaController(fakeMasterPersonaUsecase{
		updateEntryFunc: func(_ context.Context, _ string, _ usecase.MasterPersonaUpdateInput, _ usecase.MasterPersonaListQuery) (usecase.MasterPersonaMutationResult, error) {
			return usecase.MasterPersonaMutationResult{}, errors.New("update failed")
		},
	})

	_, err := controller.MasterPersonaUpdate(MasterPersonaUpdateRequestDTO{IdentityKey: "FollowersPlus.esp:FE01A812:NPC_"})
	if err == nil {
		t.Fatal("expected update error")
	}
}

func TestMasterPersonaControllerGetDetailMapsNPCProfileFields(t *testing.T) {
	race := "Nord"
	sex := "Female"
	controller := NewMasterPersonaController(fakeMasterPersonaUsecase{
		getDetailFunc: func(_ context.Context, identityKey string) (usecase.MasterPersonaEntry, error) {
			return usecase.MasterPersonaEntry{
				IdentityKey:    identityKey,
				TargetPlugin:   "FollowersPlus.esp",
				FormID:         "FE01A812",
				RecordType:     "NPC_",
				EditorID:       "FP_LysMaren",
				DisplayName:    "Lys Maren",
				Race:           &race,
				Sex:            &sex,
				VoiceType:      "FemaleYoungEager",
				ClassName:      "FPScoutClass",
				SourcePlugin:   "FollowersPlus.esp",
				PersonaSummary: "人情家で裏表がない",
				PersonaBody:    "正直で直感的な行動派",
				UpdatedAt:      time.Date(2026, 4, 15, 0, 0, 0, 0, time.UTC),
			}, nil
		},
		getRunStatusFunc: func(_ context.Context) (usecase.MasterPersonaRunStatus, error) {
			return usecase.MasterPersonaRunStatus{}, nil
		},
	})

	response, err := controller.MasterPersonaGetDetail(MasterPersonaDetailRequestDTO{IdentityKey: "FollowersPlus.esp:FE01A812:NPC_"})
	if err != nil {
		t.Fatalf("expected get detail to succeed: %v", err)
	}
	entry := response.Entry
	if entry.IdentityKey != "FollowersPlus.esp:FE01A812:NPC_" {
		t.Fatalf("unexpected identity key: %q", entry.IdentityKey)
	}
	if entry.Race == nil || *entry.Race != "Nord" {
		t.Fatalf("unexpected race: %#v", entry.Race)
	}
	if entry.Sex == nil || *entry.Sex != "Female" {
		t.Fatalf("unexpected sex: %#v", entry.Sex)
	}
	if entry.PersonaBody != "正直で直感的な行動派" {
		t.Fatalf("unexpected persona body: %q", entry.PersonaBody)
	}
	if entry.SourcePlugin != "FollowersPlus.esp" {
		t.Fatalf("unexpected source plugin: %q", entry.SourcePlugin)
	}
}

func TestMasterPersonaControllerGetDetailPropagatesUsecaseError(t *testing.T) {
	controller := NewMasterPersonaController(fakeMasterPersonaUsecase{
		getDetailFunc: func(_ context.Context, _ string) (usecase.MasterPersonaEntry, error) {
			return usecase.MasterPersonaEntry{}, errors.New("entry not found")
		},
	})

	_, err := controller.MasterPersonaGetDetail(MasterPersonaDetailRequestDTO{IdentityKey: "FollowersPlus.esp:FE01A812:NPC_"})
	if err == nil {
		t.Fatal("expected get detail to fail")
	}
}

func TestMasterPersonaControllerGetPagePropagatesError(t *testing.T) {
	controller := NewMasterPersonaController(fakeMasterPersonaUsecase{
		getPageFunc: func(_ context.Context, _ usecase.MasterPersonaListQuery, _ *string) (usecase.MasterPersonaPageState, error) {
			return usecase.MasterPersonaPageState{}, errors.New("database error")
		},
	})

	_, err := controller.MasterPersonaGetPage(MasterPersonaPageRequestDTO{})
	if err == nil {
		t.Fatal("expected get page to fail")
	}
}

// RED test: cutover で GenerationSourceJSON を detail DTO から除去することを証明する。
// toMasterPersonaDetailDTO が GenerationSourceJSON のマッピングを停止するまで失敗する。
func TestMasterPersonaControllerGetDetailDoesNotExposeGenerationSourceJSON(t *testing.T) {
	controller := NewMasterPersonaController(fakeMasterPersonaUsecase{
		getDetailFunc: func(_ context.Context, identityKey string) (usecase.MasterPersonaEntry, error) {
			return usecase.MasterPersonaEntry{
				IdentityKey:          identityKey,
				GenerationSourceJSON: "sample.json",
				UpdatedAt:            time.Date(2026, 4, 15, 0, 0, 0, 0, time.UTC),
			}, nil
		},
		getRunStatusFunc: func(_ context.Context) (usecase.MasterPersonaRunStatus, error) {
			return usecase.MasterPersonaRunStatus{}, nil
		},
	})

	response, err := controller.MasterPersonaGetDetail(MasterPersonaDetailRequestDTO{IdentityKey: "FollowersPlus.esp:FE01A812:NPC_"})
	if err != nil {
		t.Fatalf("expected get detail to succeed: %v", err)
	}
	if response.Entry.GenerationSourceJSON != "" {
		t.Fatalf("expected GenerationSourceJSON to be empty in detail DTO, got %q", response.Entry.GenerationSourceJSON)
	}
}

// RED test: cutover で BaselineApplied を detail DTO から除去することを証明する。
// toMasterPersonaDetailDTO が BaselineApplied のマッピングを停止するまで失敗する。
func TestMasterPersonaControllerGetDetailDoesNotExposeBaselineApplied(t *testing.T) {
	controller := NewMasterPersonaController(fakeMasterPersonaUsecase{
		getDetailFunc: func(_ context.Context, identityKey string) (usecase.MasterPersonaEntry, error) {
			return usecase.MasterPersonaEntry{
				IdentityKey:     identityKey,
				BaselineApplied: true,
				UpdatedAt:       time.Date(2026, 4, 15, 0, 0, 0, 0, time.UTC),
			}, nil
		},
		getRunStatusFunc: func(_ context.Context) (usecase.MasterPersonaRunStatus, error) {
			return usecase.MasterPersonaRunStatus{}, nil
		},
	})

	response, err := controller.MasterPersonaGetDetail(MasterPersonaDetailRequestDTO{IdentityKey: "FollowersPlus.esp:FE01A812:NPC_"})
	if err != nil {
		t.Fatalf("expected get detail to succeed: %v", err)
	}
	if response.Entry.BaselineApplied {
		t.Fatal("expected BaselineApplied to be absent from detail DTO, got true")
	}
}

// RED test: cutover で DialogueCount を list item DTO から除去することを証明する。
// toMasterPersonaListItemDTO が DialogueCount のマッピングを停止するまで失敗する。
func TestMasterPersonaControllerGetPageDoesNotExposeDialogueCount(t *testing.T) {
	controller := NewMasterPersonaController(fakeMasterPersonaUsecase{
		getPageFunc: func(_ context.Context, _ usecase.MasterPersonaListQuery, _ *string) (usecase.MasterPersonaPageState, error) {
			return usecase.MasterPersonaPageState{
				Items: []usecase.MasterPersonaEntry{
					{
						IdentityKey:   "FollowersPlus.esp:FE01A812:NPC_",
						DialogueCount: 44,
						UpdatedAt:     time.Date(2026, 4, 15, 0, 0, 0, 0, time.UTC),
					},
				},
				TotalCount: 1,
				Page:       1,
				PageSize:   30,
			}, nil
		},
	})

	response, err := controller.MasterPersonaGetPage(MasterPersonaPageRequestDTO{})
	if err != nil {
		t.Fatalf("expected get page to succeed: %v", err)
	}
	if len(response.Page.Items) == 0 {
		t.Fatal("expected at least one item in page response")
	}
	if response.Page.Items[0].DialogueCount != 0 {
		t.Fatalf("expected DialogueCount to be absent from list item DTO, got %d", response.Page.Items[0].DialogueCount)
	}
}

// persona-read-detail-cutover: detail DTO に canonical NPC_PROFILE join fields が含まれることを証明する。
func TestMasterPersonaControllerPersonaReadDetailCutoverDetailHasCanonicalNPCProfileFields(t *testing.T) {
	race := "Breton"
	sex := "Female"
	controller := NewMasterPersonaController(fakeMasterPersonaUsecase{
		getDetailFunc: func(_ context.Context, identityKey string) (usecase.MasterPersonaEntry, error) {
			return usecase.MasterPersonaEntry{
				IdentityKey:  identityKey,
				TargetPlugin: "FollowersPlus.esp",
				FormID:       "FE01A812",
				RecordType:   "NPC_",
				EditorID:     "FP_LysMaren",
				DisplayName:  "Lys Maren",
				Race:         &race,
				Sex:          &sex,
				VoiceType:    "FemaleYoungEager",
				ClassName:    "FPScoutClass",
				SourcePlugin: "FollowersPlus.esp",
				PersonaBody:  "短く本音を置く。",
				UpdatedAt:    time.Date(2026, 4, 15, 0, 0, 0, 0, time.UTC),
			}, nil
		},
		getRunStatusFunc: func(_ context.Context) (usecase.MasterPersonaRunStatus, error) {
			return usecase.MasterPersonaRunStatus{}, nil
		},
	})

	response, err := controller.MasterPersonaGetDetail(MasterPersonaDetailRequestDTO{IdentityKey: "FollowersPlus.esp:FE01A812:NPC_"})
	if err != nil {
		t.Fatalf("expected get detail to succeed: %v", err)
	}
	entry := response.Entry
	if entry.VoiceType != "FemaleYoungEager" {
		t.Fatalf("expected VoiceType in detail DTO, got %q", entry.VoiceType)
	}
	if entry.ClassName != "FPScoutClass" {
		t.Fatalf("expected ClassName in detail DTO, got %q", entry.ClassName)
	}
	if entry.SourcePlugin != "FollowersPlus.esp" {
		t.Fatalf("expected SourcePlugin in detail DTO, got %q", entry.SourcePlugin)
	}
	if entry.Race == nil || *entry.Race != "Breton" {
		t.Fatalf("expected Race in detail DTO, got %#v", entry.Race)
	}
	if entry.Sex == nil || *entry.Sex != "Female" {
		t.Fatalf("expected Sex in detail DTO, got %#v", entry.Sex)
	}
	if entry.PersonaBody != "短く本音を置く。" {
		t.Fatalf("expected PersonaBody in detail DTO, got %q", entry.PersonaBody)
	}
}

// persona-read-detail-cutover: detail DTO から GenerationSourceJSON が除去されることを証明する。
func TestMasterPersonaControllerPersonaReadDetailCutoverDetailDoesNotExposeGenerationSourceJSON(t *testing.T) {
	controller := NewMasterPersonaController(fakeMasterPersonaUsecase{
		getDetailFunc: func(_ context.Context, identityKey string) (usecase.MasterPersonaEntry, error) {
			return usecase.MasterPersonaEntry{
				IdentityKey:          identityKey,
				GenerationSourceJSON: "sample.json",
				UpdatedAt:            time.Date(2026, 4, 15, 0, 0, 0, 0, time.UTC),
			}, nil
		},
		getRunStatusFunc: func(_ context.Context) (usecase.MasterPersonaRunStatus, error) {
			return usecase.MasterPersonaRunStatus{}, nil
		},
	})

	response, err := controller.MasterPersonaGetDetail(MasterPersonaDetailRequestDTO{IdentityKey: "FollowersPlus.esp:FE01A812:NPC_"})
	if err != nil {
		t.Fatalf("expected get detail to succeed: %v", err)
	}
	if response.Entry.GenerationSourceJSON != "" {
		t.Fatalf("expected GenerationSourceJSON to be empty in detail DTO, got %q", response.Entry.GenerationSourceJSON)
	}
}

// persona-read-detail-cutover: detail DTO から BaselineApplied が除去されることを証明する。
func TestMasterPersonaControllerPersonaReadDetailCutoverDetailDoesNotExposeBaselineApplied(t *testing.T) {
	controller := NewMasterPersonaController(fakeMasterPersonaUsecase{
		getDetailFunc: func(_ context.Context, identityKey string) (usecase.MasterPersonaEntry, error) {
			return usecase.MasterPersonaEntry{
				IdentityKey:     identityKey,
				BaselineApplied: true,
				UpdatedAt:       time.Date(2026, 4, 15, 0, 0, 0, 0, time.UTC),
			}, nil
		},
		getRunStatusFunc: func(_ context.Context) (usecase.MasterPersonaRunStatus, error) {
			return usecase.MasterPersonaRunStatus{}, nil
		},
	})

	response, err := controller.MasterPersonaGetDetail(MasterPersonaDetailRequestDTO{IdentityKey: "FollowersPlus.esp:FE01A812:NPC_"})
	if err != nil {
		t.Fatalf("expected get detail to succeed: %v", err)
	}
	if response.Entry.BaselineApplied {
		t.Fatal("expected BaselineApplied to be absent from detail DTO")
	}
}

// persona-read-detail-cutover: list item DTO から DialogueCount が除去されることを証明する。
func TestMasterPersonaControllerPersonaReadDetailCutoverListDoesNotExposeDialogueCount(t *testing.T) {
	controller := NewMasterPersonaController(fakeMasterPersonaUsecase{
		getPageFunc: func(_ context.Context, _ usecase.MasterPersonaListQuery, _ *string) (usecase.MasterPersonaPageState, error) {
			return usecase.MasterPersonaPageState{
				Items: []usecase.MasterPersonaEntry{
					{
						IdentityKey:   "FollowersPlus.esp:FE01A812:NPC_",
						DisplayName:   "Lys Maren",
						DialogueCount: 44,
						UpdatedAt:     time.Date(2026, 4, 15, 0, 0, 0, 0, time.UTC),
					},
				},
				TotalCount: 1,
				Page:       1,
				PageSize:   30,
			}, nil
		},
	})

	response, err := controller.MasterPersonaGetPage(MasterPersonaPageRequestDTO{})
	if err != nil {
		t.Fatalf("expected get page to succeed: %v", err)
	}
	if len(response.Page.Items) == 0 {
		t.Fatal("expected at least one item in page response")
	}
	if response.Page.Items[0].DialogueCount != 0 {
		t.Fatalf("expected DialogueCount to be absent from list item DTO, got %d", response.Page.Items[0].DialogueCount)
	}
}

// persona-read-detail-cutover: detail JSON に generationSourceJson キーが存在しないことを証明する。
// MasterPersonaDetailDTO が json:"generationSourceJson" タグを保持している間は失敗する。
func TestMasterPersonaControllerPersonaReadDetailCutoverDetailJSONOmitsGenerationSourceJSON(t *testing.T) {
	controller := NewMasterPersonaController(fakeMasterPersonaUsecase{
		getDetailFunc: func(_ context.Context, identityKey string) (usecase.MasterPersonaEntry, error) {
			return usecase.MasterPersonaEntry{IdentityKey: identityKey, UpdatedAt: time.Date(2026, 4, 15, 0, 0, 0, 0, time.UTC)}, nil
		},
		getRunStatusFunc: func(_ context.Context) (usecase.MasterPersonaRunStatus, error) {
			return usecase.MasterPersonaRunStatus{}, nil
		},
	})

	response, err := controller.MasterPersonaGetDetail(MasterPersonaDetailRequestDTO{IdentityKey: "FollowersPlus.esp:FE01A812:NPC_"})
	if err != nil {
		t.Fatalf("expected get detail to succeed: %v", err)
	}
	raw, err := json.Marshal(response)
	if err != nil {
		t.Fatalf("expected JSON marshal to succeed: %v", err)
	}
	if strings.Contains(string(raw), `"generationSourceJson"`) {
		t.Fatalf("JSON still contains generationSourceJson key; payload: %s", raw)
	}
}

// persona-read-detail-cutover: detail JSON に baselineApplied キーが存在しないことを証明する。
// MasterPersonaDetailDTO が json:"baselineApplied" タグを保持している間は失敗する。
func TestMasterPersonaControllerPersonaReadDetailCutoverDetailJSONOmitsBaselineApplied(t *testing.T) {
	controller := NewMasterPersonaController(fakeMasterPersonaUsecase{
		getDetailFunc: func(_ context.Context, identityKey string) (usecase.MasterPersonaEntry, error) {
			return usecase.MasterPersonaEntry{IdentityKey: identityKey, UpdatedAt: time.Date(2026, 4, 15, 0, 0, 0, 0, time.UTC)}, nil
		},
		getRunStatusFunc: func(_ context.Context) (usecase.MasterPersonaRunStatus, error) {
			return usecase.MasterPersonaRunStatus{}, nil
		},
	})

	response, err := controller.MasterPersonaGetDetail(MasterPersonaDetailRequestDTO{IdentityKey: "FollowersPlus.esp:FE01A812:NPC_"})
	if err != nil {
		t.Fatalf("expected get detail to succeed: %v", err)
	}
	raw, err := json.Marshal(response)
	if err != nil {
		t.Fatalf("expected JSON marshal to succeed: %v", err)
	}
	if strings.Contains(string(raw), `"baselineApplied"`) {
		t.Fatalf("JSON still contains baselineApplied key; payload: %s", raw)
	}
}

// persona-read-detail-cutover: list item JSON に dialogueCount キーが存在しないことを証明する。
// MasterPersonaListItemDTO が json:"dialogueCount" タグを保持している間は失敗する。
func TestMasterPersonaControllerPersonaReadDetailCutoverListItemJSONOmitsDialogueCount(t *testing.T) {
	controller := NewMasterPersonaController(fakeMasterPersonaUsecase{
		getPageFunc: func(_ context.Context, _ usecase.MasterPersonaListQuery, _ *string) (usecase.MasterPersonaPageState, error) {
			return usecase.MasterPersonaPageState{
				Items: []usecase.MasterPersonaEntry{
					{
						IdentityKey: "FollowersPlus.esp:FE01A812:NPC_",
						DisplayName: "Lys Maren",
						UpdatedAt:   time.Date(2026, 4, 15, 0, 0, 0, 0, time.UTC),
					},
				},
				TotalCount: 1,
				Page:       1,
				PageSize:   30,
			}, nil
		},
	})

	response, err := controller.MasterPersonaGetPage(MasterPersonaPageRequestDTO{})
	if err != nil {
		t.Fatalf("expected get page to succeed: %v", err)
	}
	raw, err := json.Marshal(response)
	if err != nil {
		t.Fatalf("expected JSON marshal to succeed: %v", err)
	}
	if strings.Contains(string(raw), `"dialogueCount"`) {
		t.Fatalf("JSON still contains dialogueCount key in list item; payload: %s", raw)
	}
}

// persona-read-detail-cutover: MasterPersonaUsecasePort から GetDialogueList が除去されることを証明する。
// MasterPersonaUsecasePort インターフェースが GetDialogueList を保持している間は失敗する。
func TestMasterPersonaControllerPersonaReadDetailCutoverUsecasePortHasNoGetDialogueList(t *testing.T) {
	portType := reflect.TypeOf((*MasterPersonaUsecasePort)(nil)).Elem()
	for i := 0; i < portType.NumMethod(); i++ {
		if portType.Method(i).Name == "GetDialogueList" {
			t.Fatal("MasterPersonaUsecasePort still exposes GetDialogueList; persona-read-detail-cutover requires removal from read/detail public seam")
		}
	}
}

// persona-read-detail-cutover: MasterPersonaController から MasterPersonaGetDialogueList エンドポイントが除去されることを証明する。
// MasterPersonaController が MasterPersonaGetDialogueList メソッドを保持している間は失敗する。
func TestMasterPersonaControllerPersonaReadDetailCutoverControllerHasNoGetDialogueListEndpoint(t *testing.T) {
	controllerType := reflect.TypeOf(&MasterPersonaController{})
	for i := 0; i < controllerType.NumMethod(); i++ {
		if controllerType.Method(i).Name == "MasterPersonaGetDialogueList" {
			t.Fatal("MasterPersonaController still exposes MasterPersonaGetDialogueList; persona-read-detail-cutover requires removal")
		}
	}
}

// persona-read-detail-cutover: update save regression - MasterPersonaUpdateInputDTO の JSON に
// formId が含まれないことを証明する。
// MasterPersonaUpdateInputDTO が json:"formId" タグを保持している間は失敗する。
func TestMasterPersonaControllerPersonaReadDetailCutoverUpdateInputDTOHasNoFormID(t *testing.T) {
	dto := MasterPersonaUpdateInputDTO{
		DisplayName: "Lys Maren",
		PersonaBody: "edited persona body",
	}

	raw, err := json.Marshal(dto)
	if err != nil {
		t.Fatalf("expected JSON marshal to succeed: %v", err)
	}

	if strings.Contains(string(raw), `"formId"`) {
		t.Fatalf("JSON still contains formId key in update input DTO; payload: %s", raw)
	}
}

// persona-json-preview-cutover: preview JSON が candidateCount/newlyAddableCount/existingCount を公開することを証明する。
// MasterPersonaPreviewResponseDTO が json タグを変更するまで失敗する。
func TestMasterPersonaControllerPersonaJSONPreviewCutoverExposesNewFields(t *testing.T) {
	controller := NewMasterPersonaController(fakeMasterPersonaUsecase{
		previewGenerationFunc: func(_ context.Context, _ string, _ usecase.MasterPersonaAISettings) (usecase.MasterPersonaPreviewResult, error) {
			return usecase.MasterPersonaPreviewResult{
				FileName:          "FollowersPlus.json",
				TargetPlugin:      "FollowersPlus.esp",
				TotalNPCCount:     840,
				GeneratableCount:  228,
				ExistingSkipCount: 612,
				Status:            "生成可能",
			}, nil
		},
	})

	response, err := controller.MasterPersonaPreviewGeneration(MasterPersonaPreviewRequestDTO{
		FilePath:   "/tmp/FollowersPlus.json",
		AISettings: MasterPersonaAISettingsDTO{Provider: "gemini", Model: "gemini-2.5-pro"},
	})
	if err != nil {
		t.Fatalf("expected preview generation to succeed: %v", err)
	}

	raw, err := json.Marshal(response)
	if err != nil {
		t.Fatalf("expected json.Marshal to succeed: %v", err)
	}
	var m map[string]interface{}
	if err := json.Unmarshal(raw, &m); err != nil {
		t.Fatalf("expected json.Unmarshal to succeed: %v", err)
	}

	candidateCount, ok := m["candidateCount"]
	if !ok {
		t.Fatalf("expected candidateCount in preview JSON; got keys: %v", raw)
	}
	if candidateCount != float64(840) {
		t.Fatalf("expected candidateCount=840, got %v", candidateCount)
	}

	newlyAddableCount, ok := m["newlyAddableCount"]
	if !ok {
		t.Fatalf("expected newlyAddableCount in preview JSON; got keys: %v", raw)
	}
	if newlyAddableCount != float64(228) {
		t.Fatalf("expected newlyAddableCount=228, got %v", newlyAddableCount)
	}

	existingCount, ok := m["existingCount"]
	if !ok {
		t.Fatalf("expected existingCount in preview JSON; got keys: %v", raw)
	}
	if existingCount != float64(612) {
		t.Fatalf("expected existingCount=612, got %v", existingCount)
	}
}

// persona-json-preview-cutover: preview JSON から zeroDialogueSkipCount と genericNpcCount が省略されることを証明する。
// MasterPersonaPreviewResponseDTO がこれらフィールドに json:"-" を付与するまで失敗する。
func TestMasterPersonaControllerPersonaJSONPreviewCutoverOmitsLegacyFields(t *testing.T) {
	controller := NewMasterPersonaController(fakeMasterPersonaUsecase{
		previewGenerationFunc: func(_ context.Context, _ string, _ usecase.MasterPersonaAISettings) (usecase.MasterPersonaPreviewResult, error) {
			return usecase.MasterPersonaPreviewResult{
				FileName:              "FollowersPlus.json",
				TargetPlugin:          "FollowersPlus.esp",
				TotalNPCCount:         840,
				GeneratableCount:      228,
				ExistingSkipCount:     612,
				ZeroDialogueSkipCount: 99,
				GenericNPCCount:       77,
				Status:                "生成可能",
			}, nil
		},
	})

	response, err := controller.MasterPersonaPreviewGeneration(MasterPersonaPreviewRequestDTO{
		FilePath:   "/tmp/FollowersPlus.json",
		AISettings: MasterPersonaAISettingsDTO{Provider: "gemini", Model: "gemini-2.5-pro"},
	})
	if err != nil {
		t.Fatalf("expected preview generation to succeed: %v", err)
	}

	raw, err := json.Marshal(response)
	if err != nil {
		t.Fatalf("expected json.Marshal to succeed: %v", err)
	}

	if strings.Contains(string(raw), `"zeroDialogueSkipCount"`) {
		t.Fatalf("JSON still contains zeroDialogueSkipCount key; payload: %s", raw)
	}
	if strings.Contains(string(raw), `"genericNpcCount"`) {
		t.Fatalf("JSON still contains genericNpcCount key; payload: %s", raw)
	}
}

// persona-ai-settings-restart-cutover: MasterPersonaLoadAISettings は usecase から復元済み provider と model を返すことを証明する。
func TestMasterPersonaControllerPersonaAISettingsRestartCutoverLoadAISettingsReturnsRestoredProviderModel(t *testing.T) {
	controller := NewMasterPersonaController(fakeMasterPersonaUsecase{
		loadAISettingsFunc: func(_ context.Context) (usecase.MasterPersonaAISettings, error) {
			return usecase.MasterPersonaAISettings{Provider: "gemini", Model: "restart-model", APIKey: "restored-key"}, nil
		},
	})

	dto, err := controller.MasterPersonaLoadAISettings()
	if err != nil {
		t.Fatalf("expected load ai settings to succeed: %v", err)
	}
	if dto.Provider != "gemini" || dto.Model != "restart-model" {
		t.Fatalf("expected restored provider/model in dto, got %#v", dto)
	}
	if dto.APIKey != "restored-key" {
		t.Fatalf("expected restored api key in dto, got %q", dto.APIKey)
	}
}

// persona-ai-settings-restart-cutover: MasterPersonaSaveAISettings は DTO の provider、model、apiKey をすべて usecase へ転送することを証明する。
func TestMasterPersonaControllerPersonaAISettingsRestartCutoverSaveAISettingsForwardsAllFieldsToUsecase(t *testing.T) {
	var capturedSettings usecase.MasterPersonaAISettings
	controller := NewMasterPersonaController(fakeMasterPersonaUsecase{
		saveAISettingsFunc: func(_ context.Context, settings usecase.MasterPersonaAISettings) (usecase.MasterPersonaAISettings, error) {
			capturedSettings = settings
			return settings, nil
		},
	})

	_, err := controller.MasterPersonaSaveAISettings(MasterPersonaAISettingsDTO{Provider: "gemini", Model: "cutover-model", APIKey: "cutover-key"})
	if err != nil {
		t.Fatalf("expected save ai settings to succeed: %v", err)
	}
	if capturedSettings.Provider != "gemini" || capturedSettings.Model != "cutover-model" {
		t.Fatalf("expected provider/model forwarded to usecase, got %#v", capturedSettings)
	}
	if capturedSettings.APIKey != "cutover-key" {
		t.Fatalf("expected api key forwarded to usecase, got %q", capturedSettings.APIKey)
	}
}

// persona-read-detail-cutover: update save regression - MasterPersonaUpdateInputDTO の JSON に
// generic identity fields (editorId / voiceType / className / sourcePlugin) が含まれないことを証明する。
// MasterPersonaUpdateInputDTO が generic identity フィールドを保持している間は失敗する。
func TestMasterPersonaControllerPersonaReadDetailCutoverUpdateInputDTOHasNoGenericIdentityFields(t *testing.T) {
	dto := MasterPersonaUpdateInputDTO{
		DisplayName: "Lys Maren",
		PersonaBody: "edited persona body",
	}

	raw, err := json.Marshal(dto)
	if err != nil {
		t.Fatalf("expected JSON marshal to succeed: %v", err)
	}
	rawStr := string(raw)
	for _, field := range []string{`"editorId"`, `"voiceType"`, `"className"`, `"sourcePlugin"`} {
		if strings.Contains(rawStr, field) {
			t.Fatalf("JSON still contains generic identity field %s in update input DTO; payload: %s", field, raw)
		}
	}
}

// persona-ai-settings-restart-cutover: 新セッション (再起動後) の GetRunStatus は "入力待ち" を返すことを
// controller → usecase path で証明する。
func TestMasterPersonaControllerPersonaAISettingsRestartCutoverGetRunStatusReturnsInputWaitingOnFreshStart(t *testing.T) {
	controller := NewMasterPersonaController(fakeMasterPersonaUsecase{
		getRunStatusFunc: func(_ context.Context) (usecase.MasterPersonaRunStatus, error) {
			return usecase.MasterPersonaRunStatus{
				RunState: "入力待ち",
				Message:  "入力ファイルを選ぶと状態を表示します。",
			}, nil
		},
	})

	// Act
	status, err := controller.MasterPersonaGetRunStatus()

	// Assert: 再起動後の GetRunStatus は "入力待ち" を返す。
	if err != nil {
		t.Fatalf("expected GetRunStatus to succeed on fresh start: %v", err)
	}
	if status.RunState != "入力待ち" {
		t.Fatalf("expected GetRunStatus to return 入力待ち on fresh start (run state must not persist), got %q", status.RunState)
	}
}

// persona-ai-settings-restart-cutover: 同一セッション中に live GetRunStatus が実行状態を返すことを
// controller → usecase path で証明する。
func TestMasterPersonaControllerPersonaAISettingsRestartCutoverGetRunStatusReturnsLiveStateWithinSession(t *testing.T) {
	controller := NewMasterPersonaController(fakeMasterPersonaUsecase{
		getRunStatusFunc: func(_ context.Context) (usecase.MasterPersonaRunStatus, error) {
			return usecase.MasterPersonaRunStatus{
				RunState:       "生成中",
				TargetPlugin:   "FollowersPlus.esp",
				ProcessedCount: 2,
				Message:        "生成中...",
			}, nil
		},
	})

	// Act
	status, err := controller.MasterPersonaGetRunStatus()

	// Assert: 同一セッション中の GetRunStatus は live 状態を返す。
	if err != nil {
		t.Fatalf("expected GetRunStatus to succeed within same session: %v", err)
	}
	if status.RunState != "生成中" {
		t.Fatalf("expected GetRunStatus to return live state 生成中 within session, got %q", status.RunState)
	}
	if status.ProcessedCount != 2 {
		t.Fatalf("expected GetRunStatus to return live processed count, got %d", status.ProcessedCount)
	}
}

// RED test: persona-json-preview-cutover - preview response JSON に zeroDialogueSkipCount キーが
// 存在しないことを証明する。MasterPersonaPreviewResponseDTO から除去されるまで失敗する。
func TestPersonaJSONPreviewCutoverControllerPreviewDTOJsonOmitsZeroDialogueSkipCount(t *testing.T) {
	controller := NewMasterPersonaController(fakeMasterPersonaUsecase{
		previewGenerationFunc: func(_ context.Context, _ string, _ usecase.MasterPersonaAISettings) (usecase.MasterPersonaPreviewResult, error) {
			return usecase.MasterPersonaPreviewResult{
				FileName:              "extract.json",
				TargetPlugin:          "CutoverTest.esp",
				TotalNPCCount:         2,
				GeneratableCount:      1,
				ExistingSkipCount:     0,
				ZeroDialogueSkipCount: 1,
				Status:                "生成可能",
			}, nil
		},
	})

	response, err := controller.MasterPersonaPreviewGeneration(MasterPersonaPreviewRequestDTO{
		FilePath:   "/tmp/extract.json",
		AISettings: MasterPersonaAISettingsDTO{Provider: "gemini", Model: "gemini-2.5-pro"},
	})
	if err != nil {
		t.Fatalf("expected preview generation to succeed: %v", err)
	}
	raw, marshalErr := json.Marshal(response)
	if marshalErr != nil {
		t.Fatalf("expected JSON marshal to succeed: %v", marshalErr)
	}
	if strings.Contains(string(raw), `"zeroDialogueSkipCount"`) {
		t.Fatalf("JSON still contains zeroDialogueSkipCount key; payload: %s", raw)
	}
}

// RED test: persona-json-preview-cutover - preview response JSON に genericNpcCount キーが
// 存在しないことを証明する。MasterPersonaPreviewResponseDTO から除去されるまで失敗する。
func TestPersonaJSONPreviewCutoverControllerPreviewDTOJsonOmitsGenericNpcCount(t *testing.T) {
	controller := NewMasterPersonaController(fakeMasterPersonaUsecase{
		previewGenerationFunc: func(_ context.Context, _ string, _ usecase.MasterPersonaAISettings) (usecase.MasterPersonaPreviewResult, error) {
			return usecase.MasterPersonaPreviewResult{
				FileName:         "extract.json",
				TargetPlugin:     "CutoverTest.esp",
				TotalNPCCount:    1,
				GeneratableCount: 1,
				GenericNPCCount:  1,
				Status:           "生成可能",
			}, nil
		},
	})

	response, err := controller.MasterPersonaPreviewGeneration(MasterPersonaPreviewRequestDTO{
		FilePath:   "/tmp/extract.json",
		AISettings: MasterPersonaAISettingsDTO{Provider: "gemini", Model: "gemini-2.5-pro"},
	})
	if err != nil {
		t.Fatalf("expected preview generation to succeed: %v", err)
	}
	raw, marshalErr := json.Marshal(response)
	if marshalErr != nil {
		t.Fatalf("expected JSON marshal to succeed: %v", marshalErr)
	}
	if strings.Contains(string(raw), `"genericNpcCount"`) {
		t.Fatalf("JSON still contains genericNpcCount key; payload: %s", raw)
	}
}

// persona-json-preview-cutover: preview response JSON に candidateCount キーが存在することを証明する。
// MasterPersonaPreviewResponseDTO が json:"candidateCount" タグを持つまで失敗する。
func TestPersonaJSONPreviewCutoverControllerPreviewDTOJsonExposesCandidateCount(t *testing.T) {
	controller := NewMasterPersonaController(fakeMasterPersonaUsecase{
		previewGenerationFunc: func(_ context.Context, _ string, _ usecase.MasterPersonaAISettings) (usecase.MasterPersonaPreviewResult, error) {
			return usecase.MasterPersonaPreviewResult{
				FileName:      "FollowersPlus.json",
				TargetPlugin:  "FollowersPlus.esp",
				TotalNPCCount: 840,
				Status:        "生成可能",
			}, nil
		},
	})

	response, err := controller.MasterPersonaPreviewGeneration(MasterPersonaPreviewRequestDTO{
		FilePath:   "/tmp/FollowersPlus.json",
		AISettings: MasterPersonaAISettingsDTO{Provider: "gemini", Model: "gemini-2.5-pro"},
	})
	if err != nil {
		t.Fatalf("expected preview generation to succeed: %v", err)
	}
	raw, marshalErr := json.Marshal(response)
	if marshalErr != nil {
		t.Fatalf("expected JSON marshal to succeed: %v", marshalErr)
	}
	if !strings.Contains(string(raw), `"candidateCount"`) {
		t.Fatalf("JSON does not contain candidateCount key; payload: %s", raw)
	}
}

// persona-json-preview-cutover: preview response JSON に newlyAddableCount キーが存在することを証明する。
// MasterPersonaPreviewResponseDTO が json:"newlyAddableCount" タグを持つまで失敗する。
func TestPersonaJSONPreviewCutoverControllerPreviewDTOJsonExposesNewlyAddableCount(t *testing.T) {
	controller := NewMasterPersonaController(fakeMasterPersonaUsecase{
		previewGenerationFunc: func(_ context.Context, _ string, _ usecase.MasterPersonaAISettings) (usecase.MasterPersonaPreviewResult, error) {
			return usecase.MasterPersonaPreviewResult{
				FileName:         "FollowersPlus.json",
				TargetPlugin:     "FollowersPlus.esp",
				GeneratableCount: 228,
				Status:           "生成可能",
			}, nil
		},
	})

	response, err := controller.MasterPersonaPreviewGeneration(MasterPersonaPreviewRequestDTO{
		FilePath:   "/tmp/FollowersPlus.json",
		AISettings: MasterPersonaAISettingsDTO{Provider: "gemini", Model: "gemini-2.5-pro"},
	})
	if err != nil {
		t.Fatalf("expected preview generation to succeed: %v", err)
	}
	raw, marshalErr := json.Marshal(response)
	if marshalErr != nil {
		t.Fatalf("expected JSON marshal to succeed: %v", marshalErr)
	}
	if !strings.Contains(string(raw), `"newlyAddableCount"`) {
		t.Fatalf("JSON does not contain newlyAddableCount key; payload: %s", raw)
	}
}

// persona-json-preview-cutover: preview response JSON に existingCount キーが存在することを証明する。
// MasterPersonaPreviewResponseDTO が json:"existingCount" タグを持つまで失敗する。
func TestPersonaJSONPreviewCutoverControllerPreviewDTOJsonExposesExistingCount(t *testing.T) {
	controller := NewMasterPersonaController(fakeMasterPersonaUsecase{
		previewGenerationFunc: func(_ context.Context, _ string, _ usecase.MasterPersonaAISettings) (usecase.MasterPersonaPreviewResult, error) {
			return usecase.MasterPersonaPreviewResult{
				FileName:          "FollowersPlus.json",
				TargetPlugin:      "FollowersPlus.esp",
				ExistingSkipCount: 612,
				Status:            "生成可能",
			}, nil
		},
	})

	response, err := controller.MasterPersonaPreviewGeneration(MasterPersonaPreviewRequestDTO{
		FilePath:   "/tmp/FollowersPlus.json",
		AISettings: MasterPersonaAISettingsDTO{Provider: "gemini", Model: "gemini-2.5-pro"},
	})
	if err != nil {
		t.Fatalf("expected preview generation to succeed: %v", err)
	}
	raw, marshalErr := json.Marshal(response)
	if marshalErr != nil {
		t.Fatalf("expected JSON marshal to succeed: %v", marshalErr)
	}
	if !strings.Contains(string(raw), `"existingCount"`) {
		t.Fatalf("JSON does not contain existingCount key; payload: %s", raw)
	}
}

// persona-edit-delete-cutover: MasterPersonaUpdateInputDTO の JSON に identity/snapshot fields が含まれないことを証明する。
// editorId / voiceType / className / sourcePlugin が json タグを保持している間は失敗する。
func TestMasterPersonaControllerPersonaEditDeleteCutoverUpdateInputDTOExcludesAllIdentityFields(t *testing.T) {
	dto := MasterPersonaUpdateInputDTO{
		FormID:       "FE01A812",
		EditorID:     "FP_LysMaren",
		DisplayName:  "Lys Maren",
		VoiceType:    "FemaleYoungEager",
		ClassName:    "FPScoutClass",
		SourcePlugin: "FollowersPlus.esp",
		PersonaBody:  "edited persona body",
	}

	raw, err := json.Marshal(dto)
	if err != nil {
		t.Fatalf("expected JSON marshal to succeed: %v", err)
	}

	payload := string(raw)
	for _, forbidden := range []string{`"formId"`, `"editorId"`, `"voiceType"`, `"className"`, `"sourcePlugin"`} {
		if strings.Contains(payload, forbidden) {
			t.Fatalf("JSON contains identity field key %s; payload: %s", forbidden, payload)
		}
	}
}

// persona-edit-delete-cutover: MasterPersonaUpdate が PersonaBody を usecase input に渡すことを証明する。
// PersonaBody が usecase に届かない場合は失敗する。
func TestMasterPersonaControllerPersonaEditDeleteCutoverUpdateMapsPersonaBodyToUsecaseInput(t *testing.T) {
	var capturedInput usecase.MasterPersonaUpdateInput
	controller := NewMasterPersonaController(fakeMasterPersonaUsecase{
		updateEntryFunc: func(_ context.Context, _ string, input usecase.MasterPersonaUpdateInput, _ usecase.MasterPersonaListQuery) (usecase.MasterPersonaMutationResult, error) {
			capturedInput = input
			return usecase.MasterPersonaMutationResult{}, nil
		},
	})

	_, err := controller.MasterPersonaUpdate(MasterPersonaUpdateRequestDTO{
		IdentityKey: "FollowersPlus.esp:FE01A812:NPC_",
		Entry:       MasterPersonaUpdateInputDTO{PersonaBody: "cutover persona body"},
	})

	if err != nil {
		t.Fatalf("expected update to succeed: %v", err)
	}
	if capturedInput.PersonaBody != "cutover persona body" {
		t.Fatalf("expected PersonaBody forwarded to usecase, got %q", capturedInput.PersonaBody)
	}
}

// persona-edit-delete-cutover: MasterPersonaDelete が identityKey を usecase に転送することを証明する。
// identityKey が usecase に届かない場合は失敗する。
func TestMasterPersonaControllerPersonaEditDeleteCutoverDeleteForwardsIdentityKey(t *testing.T) {
	const wantKey = "FollowersPlus.esp:FE01A812:NPC_"
	var capturedKey string
	controller := NewMasterPersonaController(fakeMasterPersonaUsecase{
		deleteEntryFunc: func(_ context.Context, identityKey string, _ usecase.MasterPersonaListQuery) (usecase.MasterPersonaMutationResult, error) {
			capturedKey = identityKey
			return usecase.MasterPersonaMutationResult{}, nil
		},
	})

	_, err := controller.MasterPersonaDelete(MasterPersonaDeleteRequestDTO{IdentityKey: wantKey})

	if err != nil {
		t.Fatalf("expected delete to succeed: %v", err)
	}
	if capturedKey != wantKey {
		t.Fatalf("expected identityKey=%q forwarded to usecase, got %q", wantKey, capturedKey)
	}
}

// persona-edit-delete-cutover: MasterPersonaDelete が usecase エラーをラップして返すことを証明する。
// エラーが握り潰される場合は失敗する。
func TestMasterPersonaControllerPersonaEditDeleteCutoverDeletePropagatesError(t *testing.T) {
	controller := NewMasterPersonaController(fakeMasterPersonaUsecase{
		deleteEntryFunc: func(_ context.Context, _ string, _ usecase.MasterPersonaListQuery) (usecase.MasterPersonaMutationResult, error) {
			return usecase.MasterPersonaMutationResult{}, errors.New("delete failed in usecase")
		},
	})

	_, err := controller.MasterPersonaDelete(MasterPersonaDeleteRequestDTO{IdentityKey: "FollowersPlus.esp:FE01A812:NPC_"})

	if err == nil {
		t.Fatal("expected delete to propagate usecase error")
	}
}
