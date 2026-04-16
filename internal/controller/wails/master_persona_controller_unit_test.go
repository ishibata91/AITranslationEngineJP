package wails

import (
	"context"
	"errors"
	"testing"
	"time"

	"aitranslationenginejp/internal/usecase"
)

type fakeMasterPersonaUsecase struct {
	getPageFunc             func(ctx context.Context, query usecase.MasterPersonaListQuery, preferredIdentityKey *string) (usecase.MasterPersonaPageState, error)
	getDetailFunc           func(ctx context.Context, identityKey string) (usecase.MasterPersonaEntry, error)
	getDialogueListFunc     func(ctx context.Context, identityKey string) (usecase.MasterPersonaDialogueList, error)
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
func (fake fakeMasterPersonaUsecase) GetDialogueList(ctx context.Context, identityKey string) (usecase.MasterPersonaDialogueList, error) {
	if fake.getDialogueListFunc == nil {
		return usecase.MasterPersonaDialogueList{}, nil
	}
	return fake.getDialogueListFunc(ctx, identityKey)
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

func TestMasterPersonaControllerGetDialogueListMapsCountAndRows(t *testing.T) {
	controller := NewMasterPersonaController(fakeMasterPersonaUsecase{
		getDialogueListFunc: func(_ context.Context, identityKey string) (usecase.MasterPersonaDialogueList, error) {
			if identityKey != "FollowersPlus.esp:FE01A812:NPC_" {
				t.Fatalf("unexpected identity key: %q", identityKey)
			}
			return usecase.MasterPersonaDialogueList{
				IdentityKey:   identityKey,
				DialogueCount: 2,
				Dialogues: []usecase.MasterPersonaDialogueLine{
					{Index: 1, Text: "line1"},
					{Index: 2, Text: "line2"},
				},
			}, nil
		},
	})

	response, err := controller.MasterPersonaGetDialogueList(MasterPersonaDetailRequestDTO{IdentityKey: "FollowersPlus.esp:FE01A812:NPC_"})
	if err != nil {
		t.Fatalf("expected get dialogue list to succeed: %v", err)
	}
	if response.DialogueCount != 2 {
		t.Fatalf("unexpected dialogue count: %#v", response)
	}
	if len(response.Dialogues) != 2 || response.Dialogues[1].Text != "line2" {
		t.Fatalf("unexpected dialogue rows: %#v", response.Dialogues)
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
