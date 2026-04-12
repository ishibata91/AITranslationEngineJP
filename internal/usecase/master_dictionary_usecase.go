package usecase

import (
	"context"
	"fmt"
	"time"

	"aitranslationenginejp/internal/service"
)

const (
	masterDictionaryDefaultImportCategory = "すべて"
	masterDictionaryDefaultImportPage     = 1
	masterDictionaryDefaultImportPageSize = 30
)

// MasterDictionaryUsecase orchestrates operations for master dictionary management.
type MasterDictionaryUsecase struct {
	queryService   *service.MasterDictionaryQueryService
	commandService *service.MasterDictionaryCommandService
	importService  *service.MasterDictionaryImportService
	events         service.MasterDictionaryRuntimeEventPublisher
}

// NewDefaultMasterDictionaryUsecase creates a usecase with default backend wiring.
func NewDefaultMasterDictionaryUsecase(
	now func() time.Time,
	events service.MasterDictionaryRuntimeEventPublisher,
) *MasterDictionaryUsecase {
	core := service.NewDefaultMasterDictionaryService(now)
	queryService := service.NewMasterDictionaryQueryService(core)
	commandService := service.NewMasterDictionaryCommandService(core)
	importService := service.NewMasterDictionaryImportService(core)
	if events != nil {
		importService.SetProgressEmitter(func(ctx context.Context, progress int) {
			events.PublishImportProgress(ctx, progress)
		})
	}

	return NewMasterDictionaryUsecase(
		queryService,
		commandService,
		importService,
		events,
	)
}

// NewMasterDictionaryUsecase creates a new usecase.
func NewMasterDictionaryUsecase(
	queryService *service.MasterDictionaryQueryService,
	commandService *service.MasterDictionaryCommandService,
	importService *service.MasterDictionaryImportService,
	events service.MasterDictionaryRuntimeEventPublisher,
) *MasterDictionaryUsecase {
	return &MasterDictionaryUsecase{
		queryService:   queryService,
		commandService: commandService,
		importService:  importService,
		events:         events,
	}
}

// MasterDictionaryMutationInput is the usecase boundary for create/update payload.
type MasterDictionaryMutationInput struct {
	Source      string
	Translation string
	Category    string
	Origin      string
	REC         string
	EDID        string
}

// MasterDictionaryRefreshQuery describes same-page refresh conditions.
type MasterDictionaryRefreshQuery struct {
	SearchTerm string
	Category   string
	Page       int
	PageSize   int
}

// MasterDictionaryEntry is one dictionary record at usecase boundary.
type MasterDictionaryEntry = service.MasterDictionaryEntry

// MasterDictionaryPageState represents page-ready payload for list/detail sync.
type MasterDictionaryPageState struct {
	Items      []MasterDictionaryEntry
	TotalCount int
	Page       int
	PageSize   int
	SelectedID *int64
}

// MasterDictionaryMutationResult is the response payload for create/update/delete.
type MasterDictionaryMutationResult struct {
	Page           MasterDictionaryPageState
	ChangedEntry   *MasterDictionaryEntry
	DeletedEntryID *int64
}

// MasterDictionaryImportResult contains import summary and refreshed page payload.
type MasterDictionaryImportResult struct {
	Page    MasterDictionaryPageState
	Summary service.MasterDictionaryImportSummary
}

// GetPage returns list and selected entry state.
func (usecase *MasterDictionaryUsecase) GetPage(
	ctx context.Context,
	query MasterDictionaryRefreshQuery,
	preferredID *int64,
) (MasterDictionaryPageState, error) {
	listResult, err := usecase.queryService.SearchEntries(ctx, service.MasterDictionaryQuery{
		SearchTerm: query.SearchTerm,
		Category:   query.Category,
		Page:       query.Page,
		PageSize:   query.PageSize,
	})
	if err != nil {
		return MasterDictionaryPageState{}, fmt.Errorf("list page: %w", err)
	}

	selectedID := selectEntryID(listResult.Items, preferredID)
	return MasterDictionaryPageState{
		Items:      listResult.Items,
		TotalCount: listResult.TotalCount,
		Page:       listResult.Page,
		PageSize:   listResult.PageSize,
		SelectedID: selectedID,
	}, nil
}

// GetEntry returns detail payload for one entry.
func (usecase *MasterDictionaryUsecase) GetEntry(ctx context.Context, id int64) (MasterDictionaryEntry, error) {
	entry, err := usecase.queryService.LoadEntryDetail(ctx, id)
	if err != nil {
		return MasterDictionaryEntry{}, fmt.Errorf("get entry detail: %w", err)
	}
	return entry, nil
}

// CreateEntry creates an entry and returns refreshed page payload.
func (usecase *MasterDictionaryUsecase) CreateEntry(
	ctx context.Context,
	input MasterDictionaryMutationInput,
	refreshQuery MasterDictionaryRefreshQuery,
) (MasterDictionaryMutationResult, error) {
	created, err := usecase.commandService.CreateEntry(ctx, toServiceMutationInput(input))
	if err != nil {
		return MasterDictionaryMutationResult{}, fmt.Errorf("create entry: %w", err)
	}

	page, err := usecase.GetPage(ctx, refreshQuery, &created.ID)
	if err != nil {
		return MasterDictionaryMutationResult{}, fmt.Errorf("refresh page after create: %w", err)
	}
	return MasterDictionaryMutationResult{Page: page, ChangedEntry: &created}, nil
}

// UpdateEntry updates an entry and returns refreshed page payload.
func (usecase *MasterDictionaryUsecase) UpdateEntry(
	ctx context.Context,
	id int64,
	input MasterDictionaryMutationInput,
	refreshQuery MasterDictionaryRefreshQuery,
) (MasterDictionaryMutationResult, error) {
	updated, err := usecase.commandService.UpdateEntry(ctx, id, toServiceMutationInput(input))
	if err != nil {
		return MasterDictionaryMutationResult{}, fmt.Errorf("update entry: %w", err)
	}

	page, err := usecase.GetPage(ctx, refreshQuery, &updated.ID)
	if err != nil {
		return MasterDictionaryMutationResult{}, fmt.Errorf("refresh page after update: %w", err)
	}
	return MasterDictionaryMutationResult{Page: page, ChangedEntry: &updated}, nil
}

// DeleteEntry deletes an entry and returns refreshed page payload.
func (usecase *MasterDictionaryUsecase) DeleteEntry(
	ctx context.Context,
	id int64,
	refreshQuery MasterDictionaryRefreshQuery,
) (MasterDictionaryMutationResult, error) {
	if err := usecase.commandService.DeleteEntry(ctx, id); err != nil {
		return MasterDictionaryMutationResult{}, fmt.Errorf("delete entry: %w", err)
	}

	page, err := usecase.GetPage(ctx, refreshQuery, nil)
	if err != nil {
		return MasterDictionaryMutationResult{}, fmt.Errorf("refresh page after delete: %w", err)
	}

	deletedID := id
	return MasterDictionaryMutationResult{
		Page:           page,
		DeletedEntryID: &deletedID,
	}, nil
}

// ImportXML imports an XML file, returns API response payload, and emits completion event.
func (usecase *MasterDictionaryUsecase) ImportXML(
	ctx context.Context,
	xmlPath string,
	refreshQuery MasterDictionaryRefreshQuery,
) (MasterDictionaryImportResult, error) {
	summary, err := usecase.importService.ImportXML(ctx, xmlPath)
	if err != nil {
		return MasterDictionaryImportResult{}, fmt.Errorf("import xml: %w", err)
	}

	preferredID := (*int64)(nil)
	if summary.LastEntryID > 0 {
		preferredID = &summary.LastEntryID
	}

	responsePage, err := usecase.GetPage(ctx, refreshQuery, preferredID)
	if err != nil {
		return MasterDictionaryImportResult{}, fmt.Errorf("refresh page after import: %w", err)
	}

	if usecase.events != nil {
		completedPayload, buildErr := usecase.buildImportCompletedPayload(ctx, summary)
		if buildErr != nil {
			return MasterDictionaryImportResult{}, fmt.Errorf("build import completed event: %w", buildErr)
		}
		usecase.events.PublishImportCompleted(ctx, completedPayload)
	}

	return MasterDictionaryImportResult{Page: responsePage, Summary: summary}, nil
}

func (usecase *MasterDictionaryUsecase) buildImportCompletedPayload(
	ctx context.Context,
	summary service.MasterDictionaryImportSummary,
) (service.MasterDictionaryImportCompletedPayload, error) {
	preferredID := (*int64)(nil)
	if summary.LastEntryID > 0 {
		preferredID = &summary.LastEntryID
	}

	refresh := service.MasterDictionaryImportRefreshPolicy{
		Query:           "",
		Category:        masterDictionaryDefaultImportCategory,
		Page:            masterDictionaryDefaultImportPage,
		PageSize:        masterDictionaryDefaultImportPageSize,
		RefreshTargetID: preferredID,
	}

	page, err := usecase.GetPage(ctx, MasterDictionaryRefreshQuery{
		SearchTerm: refresh.Query,
		Category:   refresh.Category,
		Page:       refresh.Page,
		PageSize:   refresh.PageSize,
	}, preferredID)
	if err != nil {
		return service.MasterDictionaryImportCompletedPayload{}, err
	}

	return service.MasterDictionaryImportCompletedPayload{
		Page: service.MasterDictionaryImportCompletedPage{
			Items:      page.Items,
			TotalCount: page.TotalCount,
			Page:       page.Page,
			PageSize:   page.PageSize,
			SelectedID: page.SelectedID,
		},
		Summary: summary,
		Refresh: refresh,
	}, nil
}

// IsNotFoundError reports whether the given error means entry not found.
func IsNotFoundError(err error) bool {
	return service.IsNotFoundError(err)
}

func toServiceMutationInput(input MasterDictionaryMutationInput) service.MasterDictionaryMutationInput {
	return service.MasterDictionaryMutationInput{
		Source:      input.Source,
		Translation: input.Translation,
		Category:    input.Category,
		Origin:      input.Origin,
		REC:         input.REC,
		EDID:        input.EDID,
	}
}

func selectEntryID(items []MasterDictionaryEntry, preferredID *int64) *int64 {
	if len(items) == 0 {
		return nil
	}

	if preferredID != nil {
		for _, item := range items {
			if item.ID == *preferredID {
				selected := item.ID
				return &selected
			}
		}
	}

	selected := items[0].ID
	return &selected
}
