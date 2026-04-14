package wails

import (
	"context"
	"errors"
	"fmt"
	"testing"
	"time"

	"aitranslationenginejp/internal/usecase"
)

const controllerCanonicalXMLPath = "canonical.xml"

var controllerUpdatedAt = time.Date(2026, 4, 14, 11, 0, 0, 0, time.UTC)

type fakeMasterDictionaryUsecase struct {
	getPageFunc     func(ctx context.Context, query usecase.MasterDictionaryRefreshQuery, preferredID *int64) (usecase.MasterDictionaryPageState, error)
	getEntryFunc    func(ctx context.Context, id int64) (usecase.MasterDictionaryEntry, error)
	createEntryFunc func(ctx context.Context, input usecase.MasterDictionaryMutationInput, refreshQuery usecase.MasterDictionaryRefreshQuery) (usecase.MasterDictionaryMutationResult, error)
	updateEntryFunc func(ctx context.Context, id int64, input usecase.MasterDictionaryMutationInput, refreshQuery usecase.MasterDictionaryRefreshQuery) (usecase.MasterDictionaryMutationResult, error)
	deleteEntryFunc func(ctx context.Context, id int64, refreshQuery usecase.MasterDictionaryRefreshQuery) (usecase.MasterDictionaryMutationResult, error)
	importXMLFunc   func(ctx context.Context, xmlPath string, refreshQuery usecase.MasterDictionaryRefreshQuery) (usecase.MasterDictionaryImportResult, error)
}

func (fake fakeMasterDictionaryUsecase) GetPage(ctx context.Context, query usecase.MasterDictionaryRefreshQuery, preferredID *int64) (usecase.MasterDictionaryPageState, error) {
	if fake.getPageFunc == nil {
		return usecase.MasterDictionaryPageState{}, nil
	}
	return fake.getPageFunc(ctx, query, preferredID)
}

func (fake fakeMasterDictionaryUsecase) GetEntry(ctx context.Context, id int64) (usecase.MasterDictionaryEntry, error) {
	if fake.getEntryFunc == nil {
		return usecase.MasterDictionaryEntry{}, nil
	}
	return fake.getEntryFunc(ctx, id)
}

func (fake fakeMasterDictionaryUsecase) CreateEntry(ctx context.Context, input usecase.MasterDictionaryMutationInput, refreshQuery usecase.MasterDictionaryRefreshQuery) (usecase.MasterDictionaryMutationResult, error) {
	if fake.createEntryFunc == nil {
		return usecase.MasterDictionaryMutationResult{}, nil
	}
	return fake.createEntryFunc(ctx, input, refreshQuery)
}

func (fake fakeMasterDictionaryUsecase) UpdateEntry(ctx context.Context, id int64, input usecase.MasterDictionaryMutationInput, refreshQuery usecase.MasterDictionaryRefreshQuery) (usecase.MasterDictionaryMutationResult, error) {
	if fake.updateEntryFunc == nil {
		return usecase.MasterDictionaryMutationResult{}, nil
	}
	return fake.updateEntryFunc(ctx, id, input, refreshQuery)
}

func (fake fakeMasterDictionaryUsecase) DeleteEntry(ctx context.Context, id int64, refreshQuery usecase.MasterDictionaryRefreshQuery) (usecase.MasterDictionaryMutationResult, error) {
	if fake.deleteEntryFunc == nil {
		return usecase.MasterDictionaryMutationResult{}, nil
	}
	return fake.deleteEntryFunc(ctx, id, refreshQuery)
}

func (fake fakeMasterDictionaryUsecase) ImportXML(ctx context.Context, xmlPath string, refreshQuery usecase.MasterDictionaryRefreshQuery) (usecase.MasterDictionaryImportResult, error) {
	if fake.importXMLFunc == nil {
		return usecase.MasterDictionaryImportResult{}, nil
	}
	return fake.importXMLFunc(ctx, xmlPath, refreshQuery)
}

type fakeRuntimeEmitterSource struct {
	ctx context.Context
	ok  bool
}

type fakeRuntimeEmitterState struct {
	ctx context.Context
	ok  bool
}

func (fake fakeRuntimeEmitterSource) RuntimeEventContext() (context.Context, bool) {
	return fake.ctx, fake.ok
}

func (fake *fakeRuntimeEmitterState) RuntimeEventContext() (context.Context, bool) {
	return fake.ctx, fake.ok
}

func (fake *fakeRuntimeEmitterState) SetRuntimeContext(ctx context.Context) {
	fake.ctx = ctx
	fake.ok = ctx != nil
}

func (fake *fakeRuntimeEmitterState) ClearRuntimeContext() {
	fake.ctx = nil
	fake.ok = false
}

func TestMasterDictionaryControllerGetPageMapsUsecaseState(t *testing.T) {
	preferredID := int64(42)
	updatedAt := time.Date(2026, 4, 14, 10, 0, 0, 0, time.UTC)
	controller := NewMasterDictionaryController(fakeMasterDictionaryUsecase{
		getPageFunc: func(ctx context.Context, query usecase.MasterDictionaryRefreshQuery, preferred *int64) (usecase.MasterDictionaryPageState, error) {
			if ctx == nil {
				t.Fatal("expected request context")
			}
			if query.SearchTerm != "Auriel" || query.Category != "書籍" || query.Page != 2 || query.PageSize != 25 {
				t.Fatalf("unexpected refresh query: %#v", query)
			}
			if preferred == nil || *preferred != preferredID {
				t.Fatalf("expected preferred id %d, got %#v", preferredID, preferred)
			}
			return usecase.MasterDictionaryPageState{
				Items:      []usecase.MasterDictionaryEntry{{ID: preferredID, Source: "Auriel", Translation: "アーリエル", UpdatedAt: updatedAt}},
				TotalCount: 1,
				Page:       2,
				PageSize:   25,
				SelectedID: preferred,
			}, nil
		},
	}, nil)

	response, err := controller.MasterDictionaryGetPage(MasterDictionaryPageRequestDTO{
		Refresh:     MasterDictionaryRefreshQueryDTO{SearchTerm: "Auriel", Category: "書籍", Page: 2, PageSize: 25},
		PreferredID: &preferredID,
	})
	if err != nil {
		t.Fatalf("expected get page to succeed: %v", err)
	}
	if len(response.Page.Items) != 1 || response.Page.Items[0].ID != preferredID {
		t.Fatalf("unexpected page items: %#v", response.Page.Items)
	}
	if response.Page.SelectedID == nil || *response.Page.SelectedID != preferredID {
		t.Fatalf("expected selected id %d, got %#v", preferredID, response.Page.SelectedID)
	}
}

func TestMasterDictionaryControllerCreateEntryMapsFrontendContract(t *testing.T) {
	controller := NewMasterDictionaryController(fakeMasterDictionaryUsecase{
		createEntryFunc: func(_ context.Context, input usecase.MasterDictionaryMutationInput, refreshQuery usecase.MasterDictionaryRefreshQuery) (usecase.MasterDictionaryMutationResult, error) {
			assertControllerCreateInput(t, input, refreshQuery)
			entry := usecase.MasterDictionaryEntry{
				ID:          101,
				Source:      input.Source,
				Translation: input.Translation,
				Category:    input.Category,
				Origin:      input.Origin,
				REC:         "BOOK:FULL",
				EDID:        "BookAuriel",
				UpdatedAt:   controllerUpdatedAt,
			}
			return usecase.MasterDictionaryMutationResult{ChangedEntry: &entry, Page: usecase.MasterDictionaryPageState{Page: 3, PageSize: 1, SelectedID: &entry.ID}}, nil
		},
	}, nil)

	created, err := controller.CreateMasterDictionaryEntry(newControllerCreateRequest())
	if err != nil {
		t.Fatalf("expected create contract to succeed: %v", err)
	}
	if created.RefreshTargetID != "101" || created.Entry.Note != "REC: BOOK:FULL / EDID: BookAuriel" {
		t.Fatalf("unexpected create response: %#v", created)
	}
	if created.Page == nil || created.Page.Page != 3 || created.Page.PageSize != 1 {
		t.Fatalf("unexpected create page: %#v", created.Page)
	}
}

func TestMasterDictionaryControllerUpdateEntryMapsFrontendContract(t *testing.T) {
	controller := NewMasterDictionaryController(fakeMasterDictionaryUsecase{
		updateEntryFunc: func(_ context.Context, id int64, input usecase.MasterDictionaryMutationInput, refreshQuery usecase.MasterDictionaryRefreshQuery) (usecase.MasterDictionaryMutationResult, error) {
			assertControllerUpdateInput(t, id, input, refreshQuery)
			entry := usecase.MasterDictionaryEntry{ID: id, Source: input.Source, Translation: input.Translation, Category: input.Category, Origin: input.Origin, UpdatedAt: controllerUpdatedAt}
			return usecase.MasterDictionaryMutationResult{ChangedEntry: &entry, Page: usecase.MasterDictionaryPageState{Page: 4, PageSize: 1, SelectedID: &entry.ID}}, nil
		},
	}, nil)

	updated, err := controller.UpdateMasterDictionaryEntry(newControllerUpdateRequest())
	if err != nil {
		t.Fatalf("expected update contract to succeed: %v", err)
	}
	if updated.RefreshTargetID != "101" || updated.Entry.Translation != "更新訳語" {
		t.Fatalf("unexpected update response: %#v", updated)
	}
}

func TestMasterDictionaryControllerDeleteEntryMapsFrontendContract(t *testing.T) {
	controller := NewMasterDictionaryController(fakeMasterDictionaryUsecase{
		deleteEntryFunc: func(_ context.Context, id int64, refreshQuery usecase.MasterDictionaryRefreshQuery) (usecase.MasterDictionaryMutationResult, error) {
			if id != 101 || refreshQuery.Page != 5 || refreshQuery.PageSize != 1 {
				t.Fatalf("unexpected delete request: id=%d refresh=%#v", id, refreshQuery)
			}
			nextID := int64(102)
			return usecase.MasterDictionaryMutationResult{DeletedEntryID: &id, Page: usecase.MasterDictionaryPageState{Page: 5, PageSize: 1, SelectedID: &nextID}}, nil
		},
	}, nil)

	deleted, err := controller.DeleteMasterDictionaryEntry(DeleteMasterDictionaryEntryRequestDTO{
		ID:      "101",
		Refresh: &MasterDictionaryFrontendRefreshDTO{Page: 5, PageSize: 1},
	})
	if err != nil {
		t.Fatalf("expected delete contract to succeed: %v", err)
	}
	if deleted.DeletedID != "101" || deleted.NextSelectedID == nil || *deleted.NextSelectedID != "102" {
		t.Fatalf("unexpected delete response: %#v", deleted)
	}
}

func TestMasterDictionaryControllerGetMasterDictionaryEntryMapsNotFoundToNil(t *testing.T) {
	controller := NewMasterDictionaryController(fakeMasterDictionaryUsecase{
		getEntryFunc: func(_ context.Context, id int64) (usecase.MasterDictionaryEntry, error) {
			if id != 999 {
				t.Fatalf("expected id 999, got %d", id)
			}
			return usecase.MasterDictionaryEntry{}, errors.New("wrapped: " + usecase.ErrMasterDictionaryEntryNotFound.Error())
		},
	}, nil)

	_, err := controller.GetMasterDictionaryEntry(GetMasterDictionaryEntryRequestDTO{ID: "999"})
	if err == nil {
		t.Fatal("expected wrapped string error to remain an error")
	}

	notFoundController := NewMasterDictionaryController(fakeMasterDictionaryUsecase{
		getEntryFunc: func(_ context.Context, _ int64) (usecase.MasterDictionaryEntry, error) {
			return usecase.MasterDictionaryEntry{}, fmt.Errorf("lookup entry: %w", usecase.ErrMasterDictionaryEntryNotFound)
		},
	}, nil)

	response, err := notFoundController.GetMasterDictionaryEntry(GetMasterDictionaryEntryRequestDTO{ID: "999"})
	if err != nil {
		t.Fatalf("expected not found to be converted to nil entry: %v", err)
	}
	if response.Entry != nil {
		t.Fatalf("expected nil entry, got %#v", response.Entry)
	}
}

func TestMasterDictionaryControllerConstructorWrapsNonStateRuntimeSource(t *testing.T) {
	initialEmitter := &fakeRuntimeEventEmitter{}
	controller := NewMasterDictionaryController(fakeMasterDictionaryUsecase{}, fakeRuntimeEmitterSource{
		ctx: newRuntimeEventContext(initialEmitter),
		ok:  true,
	})

	runtimeCtx, ok := controller.runtimeEventContext()
	if !ok || runtimeCtx == nil {
		t.Fatal("expected constructor to preserve runtime context from non-state source")
	}
	resolvedEmitter, ok := extractRuntimeEventEmitter(runtimeCtx)
	if !ok || resolvedEmitter != initialEmitter {
		t.Fatal("expected preserved runtime emitter from non-state source")
	}

	nextCtx := newRuntimeEventContext(&fakeRuntimeEventEmitter{})
	controller.setRuntimeContext(nextCtx)

	updatedCtx, ok := controller.runtimeEventContext()
	if !ok || updatedCtx == nil {
		t.Fatal("expected lifecycle state updates to affect runtime event context")
	}
	resolvedEmitter, emitterOK := extractRuntimeEventEmitter(updatedCtx)
	if !emitterOK || resolvedEmitter == nil {
		t.Fatal("expected updated runtime event context to expose emitter")
	}
}

func TestMasterDictionaryControllerUsesInjectedRuntimeEmitterState(t *testing.T) {
	state := &fakeRuntimeEmitterState{}
	controller := NewMasterDictionaryController(fakeMasterDictionaryUsecase{}, state)
	runtimeCtx := newRuntimeEventContext(&fakeRuntimeEventEmitter{})

	controller.setRuntimeContext(runtimeCtx)

	if state.ctx != runtimeCtx || !state.ok {
		t.Fatal("expected injected runtime emitter state to receive lifecycle updates")
	}
	storedCtx, ok := controller.runtimeEventContext()
	if !ok || storedCtx != runtimeCtx {
		t.Fatalf("expected controller to read from injected state, got ok=%v ctx=%#v", ok, storedCtx)
	}

	controller.clearRuntimeContext()
	if state.ctx != nil || state.ok {
		t.Fatal("expected clear to reset injected runtime emitter state")
	}
}

func TestMasterDictionaryControllerImportAliasUsesFileReferenceAndPreservesPayload(t *testing.T) {
	controller := NewMasterDictionaryController(fakeMasterDictionaryUsecase{
		importXMLFunc: func(_ context.Context, xmlPath string, refreshQuery usecase.MasterDictionaryRefreshQuery) (usecase.MasterDictionaryImportResult, error) {
			if xmlPath != controllerCanonicalXMLPath {
				t.Fatalf("expected fileReference to be used, got %q", xmlPath)
			}
			if refreshQuery.Page != 6 || refreshQuery.PageSize != 2 {
				t.Fatalf("unexpected import refresh query: %#v", refreshQuery)
			}
			return usecase.MasterDictionaryImportResult{
				Page:    usecase.MasterDictionaryPageState{Page: 6, PageSize: 2},
				Summary: usecase.MasterDictionaryImportSummary{FileName: controllerCanonicalXMLPath, ImportedCount: 2},
			}, nil
		},
	}, fakeRuntimeEmitterSource{ctx: newRuntimeEventContext(&fakeRuntimeEventEmitter{}), ok: true})

	request := ImportMasterDictionaryXMLRequestDTO{
		FilePath:      "ignored.xml",
		FileReference: " " + controllerCanonicalXMLPath + " ",
		Refresh:       &MasterDictionaryFrontendRefreshDTO{Page: 6, PageSize: 2},
	}

	response, err := controller.ImportMasterDictionaryXML(request)
	if err != nil {
		t.Fatalf("expected import to succeed: %v", err)
	}
	if !response.Accepted || response.Page == nil || response.Summary == nil || response.Summary.FileName != controllerCanonicalXMLPath {
		t.Fatalf("unexpected import response: %#v", response)
	}

	alias, err := controller.ImportMasterDictionaryXml(request)
	if err != nil {
		t.Fatalf("expected alias import to succeed: %v", err)
	}
	if !alias.Accepted || alias.Summary == nil || alias.Summary.FileName != controllerCanonicalXMLPath {
		t.Fatalf("unexpected alias response: %#v", alias)
	}

	runtimeCtx, ok := controller.runtimeEventContext()
	if !ok || runtimeCtx == nil {
		t.Fatal("expected injected runtime event context to be preserved")
	}
}

func newControllerCreateRequest() CreateMasterDictionaryEntryRequestDTO {
	request := CreateMasterDictionaryEntryRequestDTO{}
	request.Payload.Source = "Source"
	request.Payload.Translation = "訳語"
	request.Payload.Category = "カテゴリ"
	request.Payload.Origin = "手動登録"
	request.Refresh = &MasterDictionaryFrontendRefreshDTO{Page: 3, PageSize: 1}
	return request
}

func newControllerUpdateRequest() UpdateMasterDictionaryEntryRequestDTO {
	request := UpdateMasterDictionaryEntryRequestDTO{ID: "101"}
	request.Payload.Source = "Source"
	request.Payload.Translation = "更新訳語"
	request.Payload.Category = "カテゴリ"
	request.Payload.Origin = "手動登録"
	request.Refresh = &MasterDictionaryFrontendRefreshDTO{Page: 4, PageSize: 1}
	return request
}

func assertControllerCreateInput(t *testing.T, input usecase.MasterDictionaryMutationInput, refreshQuery usecase.MasterDictionaryRefreshQuery) {
	t.Helper()

	if input.Source != "Source" || input.Translation != "訳語" || input.Category != "カテゴリ" || input.Origin != "手動登録" {
		t.Fatalf("unexpected create input: %#v", input)
	}
	if refreshQuery.Page != 3 || refreshQuery.PageSize != 1 {
		t.Fatalf("unexpected create refresh query: %#v", refreshQuery)
	}
}

func assertControllerUpdateInput(t *testing.T, id int64, input usecase.MasterDictionaryMutationInput, refreshQuery usecase.MasterDictionaryRefreshQuery) {
	t.Helper()

	if id != 101 || input.Translation != "更新訳語" {
		t.Fatalf("unexpected update request: id=%d input=%#v", id, input)
	}
	if refreshQuery.Page != 4 || refreshQuery.PageSize != 1 {
		t.Fatalf("unexpected update refresh query: %#v", refreshQuery)
	}
}
