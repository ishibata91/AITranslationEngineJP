package wails

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"sync"
	"time"

	"aitranslationenginejp/internal/usecase"
)

// MasterDictionaryUsecasePort defines the operations required by the Wails controller.
type MasterDictionaryUsecasePort interface {
	GetPage(ctx context.Context, query usecase.MasterDictionaryRefreshQuery, preferredID *int64) (usecase.MasterDictionaryPageState, error)
	GetEntry(ctx context.Context, id int64) (usecase.MasterDictionaryEntry, error)
	CreateEntry(ctx context.Context, input usecase.MasterDictionaryMutationInput, refreshQuery usecase.MasterDictionaryRefreshQuery) (usecase.MasterDictionaryMutationResult, error)
	UpdateEntry(ctx context.Context, id int64, input usecase.MasterDictionaryMutationInput, refreshQuery usecase.MasterDictionaryRefreshQuery) (usecase.MasterDictionaryMutationResult, error)
	DeleteEntry(ctx context.Context, id int64, refreshQuery usecase.MasterDictionaryRefreshQuery) (usecase.MasterDictionaryMutationResult, error)
	ImportXML(ctx context.Context, xmlPath string, refreshQuery usecase.MasterDictionaryRefreshQuery) (usecase.MasterDictionaryImportResult, error)
}

// RuntimeEmitterSource provides access to the current Wails runtime emitter.
type RuntimeEmitterSource interface {
	RuntimeEventContext() (context.Context, bool)
}

// RuntimeEmitterStatePort manages the current runtime emitter state.
type RuntimeEmitterStatePort interface {
	RuntimeEmitterSource
	SetRuntimeContext(ctx context.Context)
	ClearRuntimeContext()
}

type runtimeEmitterStatePort = RuntimeEmitterStatePort

// MasterDictionaryController exposes Wails-bound master dictionary entrypoints.
type MasterDictionaryController struct {
	masterDictionaryUsecase MasterDictionaryUsecasePort
	runtimeEmitterSource    RuntimeEmitterSource
	runtimeEmitterState     runtimeEmitterStatePort
}

// NewMasterDictionaryController builds a master dictionary controller.
func NewMasterDictionaryController(
	masterDictionaryUsecase MasterDictionaryUsecasePort,
	runtimeEmitterSource RuntimeEmitterSource,
) *MasterDictionaryController {
	runtimeEmitterState := resolveRuntimeEmitterState(runtimeEmitterSource)
	return &MasterDictionaryController{
		masterDictionaryUsecase: masterDictionaryUsecase,
		runtimeEmitterSource:    runtimeEmitterState,
		runtimeEmitterState:     runtimeEmitterState,
	}
}

// setRuntimeContext stores Wails runtime emitter for runtime.EventsEmit.
func (controller *MasterDictionaryController) setRuntimeContext(ctx context.Context) {
	controller.runtimeEmitterState.SetRuntimeContext(ctx)
}

// clearRuntimeContext clears stored Wails runtime emitter.
func (controller *MasterDictionaryController) clearRuntimeContext() {
	controller.runtimeEmitterState.ClearRuntimeContext()
}

// MasterDictionaryRefreshQueryDTO describes same-page refresh conditions.
type MasterDictionaryRefreshQueryDTO struct {
	SearchTerm string `json:"searchTerm"`
	Category   string `json:"category"`
	Page       int    `json:"page"`
	PageSize   int    `json:"pageSize"`
}

// MasterDictionaryPageRequestDTO describes list request parameters.
type MasterDictionaryPageRequestDTO struct {
	Refresh     MasterDictionaryRefreshQueryDTO `json:"refresh"`
	PreferredID *int64                          `json:"preferredId,omitempty"`
}

// MasterDictionaryEntryDTO is a transport DTO for one master dictionary entry.
// REC and EDID are excluded from JSON output; callers that need them use toEntryDetailDTO.
type MasterDictionaryEntryDTO struct {
	ID          int64  `json:"id"`
	Source      string `json:"source"`
	Translation string `json:"translation"`
	Category    string `json:"category"`
	Origin      string `json:"origin"`
	REC         string `json:"-"`
	EDID        string `json:"-"`
	UpdatedAt   string `json:"updatedAt"`
}

// MasterDictionaryPageDTO is a page-scoped payload shared by list/detail refreshes.
type MasterDictionaryPageDTO struct {
	Items      []MasterDictionaryEntryDTO `json:"items"`
	TotalCount int                        `json:"totalCount"`
	Page       int                        `json:"page"`
	PageSize   int                        `json:"pageSize"`
	SelectedID *int64                     `json:"selectedId,omitempty"`
}

// MasterDictionaryPageResponseDTO returns list and selection state.
type MasterDictionaryPageResponseDTO struct {
	Page MasterDictionaryPageDTO `json:"page"`
}

// MasterDictionaryDetailRequestDTO requests one entry detail.
type MasterDictionaryDetailRequestDTO struct {
	ID int64 `json:"id"`
}

// MasterDictionaryDetailResponseDTO returns one entry detail.
type MasterDictionaryDetailResponseDTO struct {
	Entry MasterDictionaryEntryDTO `json:"entry"`
}

// MasterDictionaryMutationInputDTO describes create/update payload.
type MasterDictionaryMutationInputDTO struct {
	Source      string `json:"source"`
	Translation string `json:"translation"`
	Category    string `json:"category"`
	Origin      string `json:"origin"`
	REC         string `json:"rec"`
	EDID        string `json:"edid"`
}

// MasterDictionaryCreateRequestDTO requests create + refresh.
type MasterDictionaryCreateRequestDTO struct {
	Entry   MasterDictionaryMutationInputDTO `json:"entry"`
	Refresh MasterDictionaryRefreshQueryDTO  `json:"refresh"`
}

// MasterDictionaryUpdateRequestDTO requests update + refresh.
type MasterDictionaryUpdateRequestDTO struct {
	ID      int64                            `json:"id"`
	Entry   MasterDictionaryMutationInputDTO `json:"entry"`
	Refresh MasterDictionaryRefreshQueryDTO  `json:"refresh"`
}

// MasterDictionaryDeleteRequestDTO requests delete + refresh.
type MasterDictionaryDeleteRequestDTO struct {
	ID      int64                           `json:"id"`
	Refresh MasterDictionaryRefreshQueryDTO `json:"refresh"`
}

// MasterDictionaryMutationResponseDTO returns refreshed list/detail state.
type MasterDictionaryMutationResponseDTO struct {
	Page           MasterDictionaryPageDTO   `json:"page"`
	ChangedEntry   *MasterDictionaryEntryDTO `json:"changedEntry,omitempty"`
	DeletedEntryID *int64                    `json:"deletedEntryId,omitempty"`
}

// MasterDictionaryImportRequestDTO requests XML import.
type MasterDictionaryImportRequestDTO struct {
	XMLPath string                          `json:"xmlPath"`
	Refresh MasterDictionaryRefreshQueryDTO `json:"refresh"`
}

// MasterDictionaryImportSummaryDTO reports XML import summary.
type MasterDictionaryImportSummaryDTO struct {
	FilePath      string `json:"filePath"`
	FileName      string `json:"fileName"`
	ImportedCount int    `json:"importedCount"`
	UpdatedCount  int    `json:"updatedCount"`
	SkippedCount  int    `json:"skippedCount"`
	LastEntryID   int64  `json:"lastEntryId"`
}

// MasterDictionaryImportResponseDTO returns import summary and refreshed page state.
type MasterDictionaryImportResponseDTO struct {
	Page    MasterDictionaryPageDTO          `json:"page"`
	Summary MasterDictionaryImportSummaryDTO `json:"summary"`
}

// ListMasterDictionaryEntriesRequestDTO is the frontend contract request payload.
type ListMasterDictionaryEntriesRequestDTO struct {
	Filters struct {
		Query    string `json:"query"`
		Category string `json:"category"`
		Page     int    `json:"page"`
		PageSize int    `json:"pageSize"`
	} `json:"filters"`
}

// MasterDictionaryEntrySummaryDTO is the frontend contract entry summary.
type MasterDictionaryEntrySummaryDTO struct {
	ID          string `json:"id"`
	Source      string `json:"source"`
	Translation string `json:"translation"`
	Category    string `json:"category"`
	Origin      string `json:"origin"`
	UpdatedAt   string `json:"updatedAt"`
}

// ListMasterDictionaryEntriesResponseDTO is the frontend contract list response.
type ListMasterDictionaryEntriesResponseDTO struct {
	Entries    []MasterDictionaryEntrySummaryDTO `json:"entries"`
	TotalCount int                               `json:"totalCount"`
	Page       int                               `json:"page"`
	PageSize   int                               `json:"pageSize"`
}

// GetMasterDictionaryEntryRequestDTO is the frontend contract detail request.
type GetMasterDictionaryEntryRequestDTO struct {
	ID string `json:"id"`
}

// MasterDictionaryEntryDetailDTO is the frontend contract detail payload.
type MasterDictionaryEntryDetailDTO struct {
	MasterDictionaryEntrySummaryDTO
	Note string `json:"note"`
}

// GetMasterDictionaryEntryResponseDTO is the frontend contract detail response.
type GetMasterDictionaryEntryResponseDTO struct {
	Entry *MasterDictionaryEntryDetailDTO `json:"entry"`
}

// MasterDictionaryFrontendRefreshDTO is the frontend contract refresh payload.
type MasterDictionaryFrontendRefreshDTO struct {
	Query    string `json:"query"`
	Category string `json:"category"`
	Page     int    `json:"page"`
	PageSize int    `json:"pageSize"`
}

// CreateMasterDictionaryEntryRequestDTO is the frontend contract create request.
type CreateMasterDictionaryEntryRequestDTO struct {
	Payload struct {
		Source      string `json:"source"`
		Translation string `json:"translation"`
		Category    string `json:"category"`
		Origin      string `json:"origin"`
	} `json:"payload"`
	Refresh *MasterDictionaryFrontendRefreshDTO `json:"refresh,omitempty"`
}

// CreateMasterDictionaryEntryResponseDTO is the frontend contract create response.
type CreateMasterDictionaryEntryResponseDTO struct {
	Entry           MasterDictionaryEntryDetailDTO `json:"entry"`
	RefreshTargetID string                         `json:"refreshTargetId"`
	Page            *MasterDictionaryPageDTO       `json:"page,omitempty"`
}

// UpdateMasterDictionaryEntryRequestDTO is the frontend contract update request.
type UpdateMasterDictionaryEntryRequestDTO struct {
	ID      string `json:"id"`
	Payload struct {
		Source      string `json:"source"`
		Translation string `json:"translation"`
		Category    string `json:"category"`
		Origin      string `json:"origin"`
	} `json:"payload"`
	Refresh *MasterDictionaryFrontendRefreshDTO `json:"refresh,omitempty"`
}

// UpdateMasterDictionaryEntryResponseDTO is the frontend contract update response.
type UpdateMasterDictionaryEntryResponseDTO struct {
	Entry           MasterDictionaryEntryDetailDTO `json:"entry"`
	RefreshTargetID string                         `json:"refreshTargetId"`
	Page            *MasterDictionaryPageDTO       `json:"page,omitempty"`
}

// DeleteMasterDictionaryEntryRequestDTO is the frontend contract delete request.
type DeleteMasterDictionaryEntryRequestDTO struct {
	ID      string                              `json:"id"`
	Refresh *MasterDictionaryFrontendRefreshDTO `json:"refresh,omitempty"`
}

// DeleteMasterDictionaryEntryResponseDTO is the frontend contract delete response.
type DeleteMasterDictionaryEntryResponseDTO struct {
	DeletedID      string                   `json:"deletedId"`
	NextSelectedID *string                  `json:"nextSelectedId"`
	Page           *MasterDictionaryPageDTO `json:"page,omitempty"`
}

// ImportMasterDictionaryXMLRequestDTO is the frontend contract import request.
// frontend contract keeps Xml casing.
//
//revive:disable-next-line var-naming
type ImportMasterDictionaryXMLRequestDTO struct {
	FilePath      string                              `json:"filePath"`
	FileReference string                              `json:"fileReference,omitempty"`
	Refresh       *MasterDictionaryFrontendRefreshDTO `json:"refresh,omitempty"`
}

// ImportMasterDictionaryXMLResponseDTO is the frontend contract import response.
// frontend contract keeps Xml casing.
//
//revive:disable-next-line var-naming
type ImportMasterDictionaryXMLResponseDTO struct {
	Accepted bool                              `json:"accepted"`
	Page     *MasterDictionaryPageDTO          `json:"page,omitempty"`
	Summary  *MasterDictionaryImportSummaryDTO `json:"summary,omitempty"`
}

// MasterDictionaryGetPage returns page-level list and selected entry state.
func (controller *MasterDictionaryController) MasterDictionaryGetPage(
	request MasterDictionaryPageRequestDTO,
) (MasterDictionaryPageResponseDTO, error) {
	pageState, err := controller.masterDictionaryUsecase.GetPage(
		controller.requestContext(),
		toRefreshQuery(request.Refresh),
		request.PreferredID,
	)
	if err != nil {
		return MasterDictionaryPageResponseDTO{}, fmt.Errorf("master dictionary get page: %w", err)
	}

	return MasterDictionaryPageResponseDTO{Page: toPageDTO(pageState)}, nil
}

// MasterDictionaryGetDetail returns one entry detail.
func (controller *MasterDictionaryController) MasterDictionaryGetDetail(
	request MasterDictionaryDetailRequestDTO,
) (MasterDictionaryDetailResponseDTO, error) {
	entry, err := controller.masterDictionaryUsecase.GetEntry(controller.requestContext(), request.ID)
	if err != nil {
		return MasterDictionaryDetailResponseDTO{}, fmt.Errorf("master dictionary get detail: %w", err)
	}

	return MasterDictionaryDetailResponseDTO{Entry: toEntryDTO(entry)}, nil
}

// MasterDictionaryCreate creates one entry and returns refreshed page state.
func (controller *MasterDictionaryController) MasterDictionaryCreate(
	request MasterDictionaryCreateRequestDTO,
) (MasterDictionaryMutationResponseDTO, error) {
	result, err := controller.masterDictionaryUsecase.CreateEntry(
		controller.requestContext(),
		toMutationInput(request.Entry),
		toRefreshQuery(request.Refresh),
	)
	if err != nil {
		return MasterDictionaryMutationResponseDTO{}, fmt.Errorf("master dictionary create: %w", err)
	}

	return toMutationResponseDTO(result), nil
}

// MasterDictionaryUpdate updates one entry and returns refreshed page state.
func (controller *MasterDictionaryController) MasterDictionaryUpdate(
	request MasterDictionaryUpdateRequestDTO,
) (MasterDictionaryMutationResponseDTO, error) {
	result, err := controller.masterDictionaryUsecase.UpdateEntry(
		controller.requestContext(),
		request.ID,
		toMutationInput(request.Entry),
		toRefreshQuery(request.Refresh),
	)
	if err != nil {
		return MasterDictionaryMutationResponseDTO{}, fmt.Errorf("master dictionary update: %w", err)
	}

	return toMutationResponseDTO(result), nil
}

// MasterDictionaryDelete deletes one entry and returns refreshed page state.
func (controller *MasterDictionaryController) MasterDictionaryDelete(
	request MasterDictionaryDeleteRequestDTO,
) (MasterDictionaryMutationResponseDTO, error) {
	result, err := controller.masterDictionaryUsecase.DeleteEntry(
		controller.requestContext(),
		request.ID,
		toRefreshQuery(request.Refresh),
	)
	if err != nil {
		return MasterDictionaryMutationResponseDTO{}, fmt.Errorf("master dictionary delete: %w", err)
	}

	return toMutationResponseDTO(result), nil
}

// MasterDictionaryImportXML imports XML data and returns refreshed page state.
func (controller *MasterDictionaryController) MasterDictionaryImportXML(
	request MasterDictionaryImportRequestDTO,
) (MasterDictionaryImportResponseDTO, error) {
	result, err := controller.masterDictionaryUsecase.ImportXML(
		controller.requestContext(),
		request.XMLPath,
		toRefreshQuery(request.Refresh),
	)
	if err != nil {
		return MasterDictionaryImportResponseDTO{}, fmt.Errorf("master dictionary import xml: %w", err)
	}

	return MasterDictionaryImportResponseDTO{
		Page:    toPageDTO(result.Page),
		Summary: toImportSummaryDTO(result),
	}, nil
}

// ListMasterDictionaryEntries bridges frontend contract to page response contract.
func (controller *MasterDictionaryController) ListMasterDictionaryEntries(
	request ListMasterDictionaryEntriesRequestDTO,
) (ListMasterDictionaryEntriesResponseDTO, error) {
	page, err := controller.MasterDictionaryGetPage(MasterDictionaryPageRequestDTO{
		Refresh: MasterDictionaryRefreshQueryDTO{
			SearchTerm: request.Filters.Query,
			Category:   request.Filters.Category,
			Page:       request.Filters.Page,
			PageSize:   request.Filters.PageSize,
		},
	})
	if err != nil {
		return ListMasterDictionaryEntriesResponseDTO{}, fmt.Errorf("list master dictionary entries: %w", err)
	}

	entries := make([]MasterDictionaryEntrySummaryDTO, 0, len(page.Page.Items))
	for _, item := range page.Page.Items {
		entries = append(entries, toEntrySummaryDTO(item))
	}

	return ListMasterDictionaryEntriesResponseDTO{
		Entries:    entries,
		TotalCount: page.Page.TotalCount,
		Page:       page.Page.Page,
		PageSize:   page.Page.PageSize,
	}, nil
}

// GetMasterDictionaryEntry bridges frontend contract to detail response contract.
func (controller *MasterDictionaryController) GetMasterDictionaryEntry(
	request GetMasterDictionaryEntryRequestDTO,
) (GetMasterDictionaryEntryResponseDTO, error) {
	entryID, err := parseStringID(request.ID)
	if err != nil {
		return GetMasterDictionaryEntryResponseDTO{}, fmt.Errorf("get master dictionary entry: %w", err)
	}

	detail, err := controller.MasterDictionaryGetDetail(MasterDictionaryDetailRequestDTO{ID: entryID})
	if err != nil {
		if usecase.IsNotFoundError(err) {
			return GetMasterDictionaryEntryResponseDTO{Entry: nil}, nil
		}
		return GetMasterDictionaryEntryResponseDTO{}, fmt.Errorf("get master dictionary entry detail: %w", err)
	}

	entry := toEntryDetailDTO(detail.Entry)
	return GetMasterDictionaryEntryResponseDTO{Entry: &entry}, nil
}

// CreateMasterDictionaryEntry bridges frontend contract to create response contract.
func (controller *MasterDictionaryController) CreateMasterDictionaryEntry(
	request CreateMasterDictionaryEntryRequestDTO,
) (CreateMasterDictionaryEntryResponseDTO, error) {
	result, err := controller.MasterDictionaryCreate(MasterDictionaryCreateRequestDTO{
		Entry: MasterDictionaryMutationInputDTO{
			Source:      request.Payload.Source,
			Translation: request.Payload.Translation,
			Category:    request.Payload.Category,
			Origin:      request.Payload.Origin,
		},
		Refresh: resolveFrontendRefreshQuery(request.Refresh),
	})
	if err != nil {
		return CreateMasterDictionaryEntryResponseDTO{}, fmt.Errorf("create master dictionary entry: %w", err)
	}
	if result.ChangedEntry == nil {
		return CreateMasterDictionaryEntryResponseDTO{}, fmt.Errorf("create master dictionary entry: changed entry is missing")
	}

	page := result.Page
	return CreateMasterDictionaryEntryResponseDTO{
		Entry:           toEntryDetailDTO(*result.ChangedEntry),
		RefreshTargetID: strconv.FormatInt(result.ChangedEntry.ID, 10),
		Page:            &page,
	}, nil
}

// UpdateMasterDictionaryEntry bridges frontend contract to update response contract.
func (controller *MasterDictionaryController) UpdateMasterDictionaryEntry(
	request UpdateMasterDictionaryEntryRequestDTO,
) (UpdateMasterDictionaryEntryResponseDTO, error) {
	entryID, err := parseStringID(request.ID)
	if err != nil {
		return UpdateMasterDictionaryEntryResponseDTO{}, fmt.Errorf("update master dictionary entry: %w", err)
	}

	result, err := controller.MasterDictionaryUpdate(MasterDictionaryUpdateRequestDTO{
		ID: entryID,
		Entry: MasterDictionaryMutationInputDTO{
			Source:      request.Payload.Source,
			Translation: request.Payload.Translation,
			Category:    request.Payload.Category,
			Origin:      request.Payload.Origin,
		},
		Refresh: resolveFrontendRefreshQuery(request.Refresh),
	})
	if err != nil {
		return UpdateMasterDictionaryEntryResponseDTO{}, fmt.Errorf("update master dictionary entry: %w", err)
	}
	if result.ChangedEntry == nil {
		return UpdateMasterDictionaryEntryResponseDTO{}, fmt.Errorf("update master dictionary entry: changed entry is missing")
	}

	page := result.Page
	return UpdateMasterDictionaryEntryResponseDTO{
		Entry:           toEntryDetailDTO(*result.ChangedEntry),
		RefreshTargetID: strconv.FormatInt(result.ChangedEntry.ID, 10),
		Page:            &page,
	}, nil
}

// DeleteMasterDictionaryEntry bridges frontend contract to delete response contract.
func (controller *MasterDictionaryController) DeleteMasterDictionaryEntry(
	request DeleteMasterDictionaryEntryRequestDTO,
) (DeleteMasterDictionaryEntryResponseDTO, error) {
	entryID, err := parseStringID(request.ID)
	if err != nil {
		return DeleteMasterDictionaryEntryResponseDTO{}, fmt.Errorf("delete master dictionary entry: %w", err)
	}

	result, err := controller.MasterDictionaryDelete(MasterDictionaryDeleteRequestDTO{
		ID:      entryID,
		Refresh: resolveFrontendRefreshQuery(request.Refresh),
	})
	if err != nil {
		return DeleteMasterDictionaryEntryResponseDTO{}, fmt.Errorf("delete master dictionary entry: %w", err)
	}

	var nextSelectedID *string
	if result.Page.SelectedID != nil {
		selected := strconv.FormatInt(*result.Page.SelectedID, 10)
		nextSelectedID = &selected
	}

	page := result.Page
	return DeleteMasterDictionaryEntryResponseDTO{
		DeletedID:      strconv.FormatInt(entryID, 10),
		NextSelectedID: nextSelectedID,
		Page:           &page,
	}, nil
}

// ImportMasterDictionaryXML bridges frontend contract to import response contract.
// frontend binding name is fixed by contract.
//
//revive:disable-next-line var-naming
func (controller *MasterDictionaryController) ImportMasterDictionaryXML(
	request ImportMasterDictionaryXMLRequestDTO,
) (ImportMasterDictionaryXMLResponseDTO, error) {
	xmlReference := resolveMasterDictionaryImportReference(request)
	result, err := controller.MasterDictionaryImportXML(MasterDictionaryImportRequestDTO{
		XMLPath: xmlReference,
		Refresh: resolveFrontendRefreshQuery(request.Refresh),
	})
	if err != nil {
		return ImportMasterDictionaryXMLResponseDTO{}, fmt.Errorf("import master dictionary xml: %w", err)
	}

	page := result.Page
	summary := result.Summary
	return ImportMasterDictionaryXMLResponseDTO{
		Accepted: true,
		Page:     &page,
		Summary:  &summary,
	}, nil
}

// ImportMasterDictionaryXml keeps the frontend contract binding name.
//
//nolint:staticcheck // Wails binding name is fixed by frontend contract.
//revive:disable-next-line var-naming
func (controller *MasterDictionaryController) ImportMasterDictionaryXml(
	request ImportMasterDictionaryXMLRequestDTO,
) (ImportMasterDictionaryXMLResponseDTO, error) {
	return controller.ImportMasterDictionaryXML(request)
}

func (controller *MasterDictionaryController) requestContext() context.Context {
	return context.Background()
}

// runtimeEventContext returns the current runtime event context for Wails event publication.
func (controller *MasterDictionaryController) runtimeEventContext() (context.Context, bool) {
	return controller.runtimeEmitterSource.RuntimeEventContext()
}

func extractRuntimeEventEmitter(ctx context.Context) (runtimeEventEmitter, bool) {
	if ctx == nil {
		return nil, false
	}
	if emitter, ok := ctx.Value(runtimeEventEmitterContextKey).(runtimeEventEmitter); ok && emitter != nil {
		return emitter, true
	}
	emitter, ok := ctx.Value(runtimeEventEmitterValueContextKey).(runtimeEventEmitter)
	if !ok || emitter == nil {
		return nil, false
	}
	return emitter, true
}

func newRuntimeEventContext(emitter runtimeEventEmitter) context.Context {
	base := context.WithValue(context.Background(), runtimeEventEmitterContextKey, emitter)
	return runtimeEventContext{
		Context: base,
		emitter: emitter,
	}
}

type runtimeEventEmitter interface {
	Emit(eventName string, optionalData ...interface{})
}

type runtimeEventContext struct {
	context.Context
	emitter runtimeEventEmitter
}

func (ctx runtimeEventContext) Value(key interface{}) interface{} {
	if key == runtimeEventEmitterValueContextKey {
		return ctx.emitter
	}
	return ctx.Context.Value(key)
}

type runtimeEventContextKey string

const (
	runtimeEventEmitterContextKey      runtimeEventContextKey = "events"
	runtimeEventEmitterValueContextKey string                 = "events"
)

type masterDictionaryRuntimeEmitterState struct {
	runtimeEventEmitterMu sync.RWMutex
	runtimeEventEmitter   runtimeEventEmitter
}

// NewRuntimeEmitterState creates a runtime emitter state shared by controller and publisher wiring.
func NewRuntimeEmitterState() RuntimeEmitterStatePort {
	return &masterDictionaryRuntimeEmitterState{}
}

func newMasterDictionaryRuntimeEmitterState() *masterDictionaryRuntimeEmitterState {
	return &masterDictionaryRuntimeEmitterState{}
}

func resolveRuntimeEmitterState(runtimeEmitterSource RuntimeEmitterSource) runtimeEmitterStatePort {
	if runtimeEmitterState, ok := runtimeEmitterSource.(runtimeEmitterStatePort); ok {
		return runtimeEmitterState
	}

	runtimeEmitterState := newMasterDictionaryRuntimeEmitterState()
	if runtimeEmitterSource == nil {
		return runtimeEmitterState
	}

	if runtimeCtx, ok := runtimeEmitterSource.RuntimeEventContext(); ok && runtimeCtx != nil {
		runtimeEmitterState.SetRuntimeContext(runtimeCtx)
	}
	return runtimeEmitterState
}

func (state *masterDictionaryRuntimeEmitterState) SetRuntimeContext(ctx context.Context) {
	emitter, ok := extractRuntimeEventEmitter(ctx)
	if !ok {
		state.ClearRuntimeContext()
		return
	}

	state.runtimeEventEmitterMu.Lock()
	defer state.runtimeEventEmitterMu.Unlock()
	state.runtimeEventEmitter = emitter
}

func (state *masterDictionaryRuntimeEmitterState) ClearRuntimeContext() {
	state.runtimeEventEmitterMu.Lock()
	defer state.runtimeEventEmitterMu.Unlock()
	state.runtimeEventEmitter = nil
}

func (state *masterDictionaryRuntimeEmitterState) RuntimeEventContext() (context.Context, bool) {
	state.runtimeEventEmitterMu.RLock()
	emitter := state.runtimeEventEmitter
	state.runtimeEventEmitterMu.RUnlock()
	if emitter == nil {
		return nil, false
	}
	return newRuntimeEventContext(emitter), true
}

func parseStringID(rawID string) (int64, error) {
	trimmed := strings.TrimSpace(rawID)
	if trimmed == "" {
		return 0, fmt.Errorf("id is required")
	}

	value, err := strconv.ParseInt(trimmed, 10, 64)
	if err != nil {
		return 0, fmt.Errorf("parse id %q: %w", rawID, err)
	}
	if value <= 0 {
		return 0, fmt.Errorf("id must be greater than zero")
	}
	return value, nil
}

func resolveMasterDictionaryImportReference(request ImportMasterDictionaryXMLRequestDTO) string {
	if reference := strings.TrimSpace(request.FileReference); reference != "" {
		return reference
	}
	return strings.TrimSpace(request.FilePath)
}

func toEntrySummaryDTO(entry MasterDictionaryEntryDTO) MasterDictionaryEntrySummaryDTO {
	return MasterDictionaryEntrySummaryDTO{
		ID:          strconv.FormatInt(entry.ID, 10),
		Source:      entry.Source,
		Translation: entry.Translation,
		Category:    entry.Category,
		Origin:      entry.Origin,
		UpdatedAt:   entry.UpdatedAt,
	}
}

func toEntryDetailDTO(entry MasterDictionaryEntryDTO) MasterDictionaryEntryDetailDTO {
	return MasterDictionaryEntryDetailDTO{
		MasterDictionaryEntrySummaryDTO: toEntrySummaryDTO(entry),
		Note:                            "マスター辞書エントリ",
	}
}

func toRefreshQuery(dto MasterDictionaryRefreshQueryDTO) usecase.MasterDictionaryRefreshQuery {
	return usecase.MasterDictionaryRefreshQuery{
		SearchTerm: dto.SearchTerm,
		Category:   dto.Category,
		Page:       dto.Page,
		PageSize:   dto.PageSize,
	}
}

func resolveFrontendRefreshQuery(refresh *MasterDictionaryFrontendRefreshDTO) MasterDictionaryRefreshQueryDTO {
	if refresh == nil {
		return defaultMasterDictionaryRefreshQueryDTO()
	}
	return MasterDictionaryRefreshQueryDTO{
		SearchTerm: refresh.Query,
		Category:   refresh.Category,
		Page:       refresh.Page,
		PageSize:   refresh.PageSize,
	}
}

func defaultMasterDictionaryRefreshQueryDTO() MasterDictionaryRefreshQueryDTO {
	return MasterDictionaryRefreshQueryDTO{Page: 1, PageSize: 30}
}

func toImportSummaryDTO(result usecase.MasterDictionaryImportResult) MasterDictionaryImportSummaryDTO {
	return MasterDictionaryImportSummaryDTO{
		FilePath:      result.Summary.FilePath,
		FileName:      result.Summary.FileName,
		ImportedCount: result.Summary.ImportedCount,
		UpdatedCount:  result.Summary.UpdatedCount,
		SkippedCount:  result.Summary.SkippedCount,
		LastEntryID:   result.Summary.LastEntryID,
	}
}

func toMutationInput(dto MasterDictionaryMutationInputDTO) usecase.MasterDictionaryMutationInput {
	return usecase.MasterDictionaryMutationInput{
		Source:      dto.Source,
		Translation: dto.Translation,
		Category:    dto.Category,
		Origin:      dto.Origin,
		REC:         dto.REC,
		EDID:        dto.EDID,
	}
}

func toPageDTO(page usecase.MasterDictionaryPageState) MasterDictionaryPageDTO {
	items := make([]MasterDictionaryEntryDTO, 0, len(page.Items))
	for _, item := range page.Items {
		items = append(items, toEntryDTO(item))
	}

	return MasterDictionaryPageDTO{
		Items:      items,
		TotalCount: page.TotalCount,
		Page:       page.Page,
		PageSize:   page.PageSize,
		SelectedID: page.SelectedID,
	}
}

func toEntryDTO(entry usecase.MasterDictionaryEntry) MasterDictionaryEntryDTO {
	return MasterDictionaryEntryDTO{
		ID:          entry.ID,
		Source:      entry.Source,
		Translation: entry.Translation,
		Category:    entry.Category,
		Origin:      entry.Origin,
		REC:         entry.REC,
		EDID:        entry.EDID,
		UpdatedAt:   entry.UpdatedAt.Format(time.RFC3339),
	}
}

func toMutationResponseDTO(result usecase.MasterDictionaryMutationResult) MasterDictionaryMutationResponseDTO {
	response := MasterDictionaryMutationResponseDTO{
		Page:           toPageDTO(result.Page),
		DeletedEntryID: result.DeletedEntryID,
	}
	if result.ChangedEntry != nil {
		entry := toEntryDTO(*result.ChangedEntry)
		response.ChangedEntry = &entry
	}
	return response
}
