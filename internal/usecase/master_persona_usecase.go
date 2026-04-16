package usecase

import (
	"context"
	"errors"
	"fmt"

	"aitranslationenginejp/internal/service"
)

// MasterPersonaListQuery describes list conditions at the usecase boundary.
type MasterPersonaListQuery struct {
	Keyword      string
	PluginFilter string
	Page         int
	PageSize     int
}

// MasterPersonaEntry aliases the service-layer master persona entry.
type MasterPersonaEntry = service.MasterPersonaEntry

// MasterPersonaPluginGroup aliases the service-layer plugin group summary.
type MasterPersonaPluginGroup = service.MasterPersonaPluginGroup

// MasterPersonaPageState stores page-level list state for the screen.
type MasterPersonaPageState struct {
	Items               []MasterPersonaEntry
	PluginGroups        []MasterPersonaPluginGroup
	TotalCount          int
	Page                int
	PageSize            int
	SelectedIdentityKey *string
}

// MasterPersonaDialogueLine aliases the service-layer dialogue line.
type MasterPersonaDialogueLine = service.MasterPersonaDialogueLine

// MasterPersonaDialogueList aliases the service-layer dialogue list.
type MasterPersonaDialogueList = service.MasterPersonaDialogueList

// MasterPersonaAISettings aliases the service-layer page-local AI settings.
type MasterPersonaAISettings = service.MasterPersonaAISettings

// MasterPersonaPreviewResult aliases the service-layer preview result.
type MasterPersonaPreviewResult = service.MasterPersonaPreviewResult

// MasterPersonaRunStatus aliases the service-layer run status.
type MasterPersonaRunStatus = service.MasterPersonaRunStatus

// MasterPersonaUpdateInput aliases the service-layer update input.
type MasterPersonaUpdateInput = service.MasterPersonaUpdateInput

// MasterPersonaMutationResult stores refreshed page state after update or delete.
type MasterPersonaMutationResult struct {
	Page           MasterPersonaPageState
	ChangedEntry   *MasterPersonaEntry
	DeletedEntryID *string
}

// MasterPersonaQueryServicePort defines read-only master persona usecase dependencies.
type MasterPersonaQueryServicePort interface {
	SearchEntries(ctx context.Context, query service.MasterPersonaListQuery) (service.MasterPersonaListResult, error)
	LoadEntryDetail(ctx context.Context, identityKey string) (service.MasterPersonaEntry, error)
	LoadDialogueList(ctx context.Context, identityKey string) (service.MasterPersonaDialogueList, error)
}

// MasterPersonaGenerationServicePort defines master persona generation and mutation dependencies.
type MasterPersonaGenerationServicePort interface {
	LoadSettings(ctx context.Context) (service.MasterPersonaAISettings, error)
	SaveSettings(ctx context.Context, settings service.MasterPersonaAISettings) (service.MasterPersonaAISettings, error)
	Preview(ctx context.Context, filePath string, requestSettings service.MasterPersonaAISettings) (service.MasterPersonaPreviewResult, error)
	Execute(ctx context.Context, filePath string, requestSettings service.MasterPersonaAISettings) (service.MasterPersonaRunStatus, error)
	UpdateEntry(ctx context.Context, identityKey string, input service.MasterPersonaUpdateInput) (service.MasterPersonaEntry, error)
	DeleteEntry(ctx context.Context, identityKey string) error
}

// MasterPersonaRunStatusServicePort defines run status dependencies.
type MasterPersonaRunStatusServicePort interface {
	GetStatus(ctx context.Context) (service.MasterPersonaRunStatus, error)
	Interrupt(ctx context.Context) (service.MasterPersonaRunStatus, error)
	Cancel(ctx context.Context) (service.MasterPersonaRunStatus, error)
}

// MasterPersonaUsecase orchestrates master persona page actions.
type MasterPersonaUsecase struct {
	queryService      MasterPersonaQueryServicePort
	generationService MasterPersonaGenerationServicePort
	runStatusService  MasterPersonaRunStatusServicePort
}

// NewMasterPersonaUsecase creates a master persona usecase.
func NewMasterPersonaUsecase(
	queryService MasterPersonaQueryServicePort,
	generationService MasterPersonaGenerationServicePort,
	runStatusService MasterPersonaRunStatusServicePort,
) *MasterPersonaUsecase {
	return &MasterPersonaUsecase{
		queryService:      queryService,
		generationService: generationService,
		runStatusService:  runStatusService,
	}
}

// GetPage returns page-level list state and selected identity key.
func (usecase *MasterPersonaUsecase) GetPage(
	ctx context.Context,
	query MasterPersonaListQuery,
	preferredIdentityKey *string,
) (MasterPersonaPageState, error) {
	result, err := usecase.queryService.SearchEntries(ctx, service.MasterPersonaListQuery{
		Keyword:      query.Keyword,
		PluginFilter: query.PluginFilter,
		Page:         query.Page,
		PageSize:     query.PageSize,
	})
	if err != nil {
		return MasterPersonaPageState{}, fmt.Errorf("list master persona page: %w", err)
	}
	selectedIdentityKey := selectMasterPersonaIdentityKey(result.Items, preferredIdentityKey)
	return MasterPersonaPageState{
		Items:               result.Items,
		PluginGroups:        result.PluginGroups,
		TotalCount:          result.TotalCount,
		Page:                result.Page,
		PageSize:            result.PageSize,
		SelectedIdentityKey: selectedIdentityKey,
	}, nil
}

// GetDetail returns one persona detail payload.
func (usecase *MasterPersonaUsecase) GetDetail(ctx context.Context, identityKey string) (MasterPersonaEntry, error) {
	entry, err := usecase.queryService.LoadEntryDetail(ctx, identityKey)
	if err != nil {
		return MasterPersonaEntry{}, fmt.Errorf("get master persona detail: %w", err)
	}
	return entry, nil
}

// GetDialogueList returns one persona dialogue list payload.
func (usecase *MasterPersonaUsecase) GetDialogueList(ctx context.Context, identityKey string) (MasterPersonaDialogueList, error) {
	result, err := usecase.queryService.LoadDialogueList(ctx, identityKey)
	if err != nil {
		return MasterPersonaDialogueList{}, fmt.Errorf("get master persona dialogue list: %w", err)
	}
	return result, nil
}

// LoadAISettings loads page-local AI settings.
func (usecase *MasterPersonaUsecase) LoadAISettings(ctx context.Context) (MasterPersonaAISettings, error) {
	result, err := usecase.generationService.LoadSettings(ctx)
	if err != nil {
		return MasterPersonaAISettings{}, fmt.Errorf("load master persona ai settings: %w", err)
	}
	return result, nil
}

// SaveAISettings saves page-local AI settings.
func (usecase *MasterPersonaUsecase) SaveAISettings(ctx context.Context, settings MasterPersonaAISettings) (MasterPersonaAISettings, error) {
	result, err := usecase.generationService.SaveSettings(ctx, settings)
	if err != nil {
		return MasterPersonaAISettings{}, fmt.Errorf("save master persona ai settings: %w", err)
	}
	return result, nil
}

// PreviewGeneration calculates preview counts before execution.
func (usecase *MasterPersonaUsecase) PreviewGeneration(
	ctx context.Context,
	filePath string,
	settings MasterPersonaAISettings,
) (MasterPersonaPreviewResult, error) {
	result, err := usecase.generationService.Preview(ctx, filePath, settings)
	if err != nil {
		return MasterPersonaPreviewResult{}, fmt.Errorf("preview master persona generation: %w", err)
	}
	return result, nil
}

// ExecuteGeneration starts generation execution.
func (usecase *MasterPersonaUsecase) ExecuteGeneration(
	ctx context.Context,
	filePath string,
	settings MasterPersonaAISettings,
) (MasterPersonaRunStatus, error) {
	result, err := usecase.generationService.Execute(ctx, filePath, settings)
	if err != nil {
		return MasterPersonaRunStatus{}, fmt.Errorf("execute master persona generation: %w", err)
	}
	return result, nil
}

// GetRunStatus returns the current generation run status.
func (usecase *MasterPersonaUsecase) GetRunStatus(ctx context.Context) (MasterPersonaRunStatus, error) {
	result, err := usecase.runStatusService.GetStatus(ctx)
	if err != nil {
		return MasterPersonaRunStatus{}, fmt.Errorf("get master persona run status: %w", err)
	}
	return result, nil
}

// InterruptGeneration interrupts the current run when possible.
func (usecase *MasterPersonaUsecase) InterruptGeneration(ctx context.Context) (MasterPersonaRunStatus, error) {
	result, err := usecase.runStatusService.Interrupt(ctx)
	if err != nil {
		return MasterPersonaRunStatus{}, fmt.Errorf("interrupt master persona generation: %w", err)
	}
	return result, nil
}

// CancelGeneration cancels the current run when possible.
func (usecase *MasterPersonaUsecase) CancelGeneration(ctx context.Context) (MasterPersonaRunStatus, error) {
	result, err := usecase.runStatusService.Cancel(ctx)
	if err != nil {
		return MasterPersonaRunStatus{}, fmt.Errorf("cancel master persona generation: %w", err)
	}
	return result, nil
}

// UpdateEntry updates one persona entry and returns refreshed page state.
func (usecase *MasterPersonaUsecase) UpdateEntry(
	ctx context.Context,
	identityKey string,
	input MasterPersonaUpdateInput,
	refreshQuery MasterPersonaListQuery,
) (MasterPersonaMutationResult, error) {
	updated, err := usecase.generationService.UpdateEntry(ctx, identityKey, input)
	if err != nil {
		return MasterPersonaMutationResult{}, fmt.Errorf("update master persona entry: %w", err)
	}
	page, err := usecase.GetPage(ctx, refreshQuery, &updated.IdentityKey)
	if err != nil {
		return MasterPersonaMutationResult{}, fmt.Errorf("refresh master persona page after update: %w", err)
	}
	return MasterPersonaMutationResult{Page: page, ChangedEntry: &updated}, nil
}

// DeleteEntry deletes one persona entry and returns refreshed page state.
func (usecase *MasterPersonaUsecase) DeleteEntry(
	ctx context.Context,
	identityKey string,
	refreshQuery MasterPersonaListQuery,
) (MasterPersonaMutationResult, error) {
	if err := usecase.generationService.DeleteEntry(ctx, identityKey); err != nil {
		return MasterPersonaMutationResult{}, fmt.Errorf("delete master persona entry: %w", err)
	}
	page, err := usecase.GetPage(ctx, refreshQuery, nil)
	if err != nil {
		return MasterPersonaMutationResult{}, fmt.Errorf("refresh master persona page after delete: %w", err)
	}
	deletedIdentityKey := identityKey
	return MasterPersonaMutationResult{Page: page, DeletedEntryID: &deletedIdentityKey}, nil
}

// IsMasterPersonaNotFoundError reports whether the error means persona not found.
func IsMasterPersonaNotFoundError(err error) bool {
	return service.IsMasterPersonaNotFoundError(err)
}

// IsMasterPersonaActiveRunError reports whether the error means mutations are locked by an active run.
func IsMasterPersonaActiveRunError(err error) bool {
	return errors.Is(err, service.ErrMasterPersonaActiveRun)
}

func selectMasterPersonaIdentityKey(items []MasterPersonaEntry, preferredIdentityKey *string) *string {
	if len(items) == 0 {
		return nil
	}
	if preferredIdentityKey != nil {
		for _, item := range items {
			if item.IdentityKey == *preferredIdentityKey {
				selected := item.IdentityKey
				return &selected
			}
		}
	}
	selected := items[0].IdentityKey
	return &selected
}
