package usecase

import (
	"context"
	"errors"
	"testing"

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

func TestMasterPersonaUsecaseGetDialogueListReturnsDialogueCountAndRows(t *testing.T) {
	usecase := NewMasterPersonaUsecase(
		fakeMasterPersonaQueryService{
			loadDialogueListFunc: func(_ context.Context, identityKey string) (service.MasterPersonaDialogueList, error) {
				if identityKey != "FollowersPlus.esp:FE01A812:NPC_" {
					t.Fatalf("unexpected identity key: %q", identityKey)
				}
				return service.MasterPersonaDialogueList{
					IdentityKey:   identityKey,
					DialogueCount: 2,
					Dialogues:     []service.MasterPersonaDialogueLine{{Index: 1, Text: "line1"}, {Index: 2, Text: "line2"}},
				}, nil
			},
		},
		fakeMasterPersonaGenerationService{},
		fakeMasterPersonaRunStatusService{},
	)

	result, err := usecase.GetDialogueList(context.Background(), "FollowersPlus.esp:FE01A812:NPC_")
	if err != nil {
		t.Fatalf("expected get dialogue list to succeed: %v", err)
	}
	if result.DialogueCount != 2 || len(result.Dialogues) != 2 {
		t.Fatalf("unexpected dialogue list: %#v", result)
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
