package usecase

import (
	"context"
	"errors"
	"testing"
	"time"

	"aitranslationenginejp/internal/service"
)

const (
	usecaseBookREC       = "BOOK:FULL"
	usecaseDictionaryXML = "dictionary.xml"
)

type fakeQueryService struct {
	searchEntriesFunc   func(ctx context.Context, query service.MasterDictionaryQuery) (service.MasterDictionaryListResult, error)
	loadEntryDetailFunc func(ctx context.Context, id int64) (service.MasterDictionaryEntry, error)
}

func (fake fakeQueryService) SearchEntries(ctx context.Context, query service.MasterDictionaryQuery) (service.MasterDictionaryListResult, error) {
	if fake.searchEntriesFunc == nil {
		return service.MasterDictionaryListResult{}, nil
	}
	return fake.searchEntriesFunc(ctx, query)
}

func (fake fakeQueryService) LoadEntryDetail(ctx context.Context, id int64) (service.MasterDictionaryEntry, error) {
	if fake.loadEntryDetailFunc == nil {
		return service.MasterDictionaryEntry{}, nil
	}
	return fake.loadEntryDetailFunc(ctx, id)
}

type fakeCommandService struct {
	createEntryFunc func(ctx context.Context, input service.MasterDictionaryMutationInput) (service.MasterDictionaryEntry, error)
	updateEntryFunc func(ctx context.Context, id int64, input service.MasterDictionaryMutationInput) (service.MasterDictionaryEntry, error)
	deleteEntryFunc func(ctx context.Context, id int64) error
}

func (fake fakeCommandService) CreateEntry(ctx context.Context, input service.MasterDictionaryMutationInput) (service.MasterDictionaryEntry, error) {
	if fake.createEntryFunc == nil {
		return service.MasterDictionaryEntry{}, nil
	}
	return fake.createEntryFunc(ctx, input)
}

func (fake fakeCommandService) UpdateEntry(ctx context.Context, id int64, input service.MasterDictionaryMutationInput) (service.MasterDictionaryEntry, error) {
	if fake.updateEntryFunc == nil {
		return service.MasterDictionaryEntry{}, nil
	}
	return fake.updateEntryFunc(ctx, id, input)
}

func (fake fakeCommandService) DeleteEntry(ctx context.Context, id int64) error {
	if fake.deleteEntryFunc == nil {
		return nil
	}
	return fake.deleteEntryFunc(ctx, id)
}

type fakeImportService struct {
	importXMLFunc func(ctx context.Context, xmlPath string) (service.MasterDictionaryImportSummary, error)
}

func (fake fakeImportService) ImportXML(ctx context.Context, xmlPath string) (service.MasterDictionaryImportSummary, error) {
	if fake.importXMLFunc == nil {
		return service.MasterDictionaryImportSummary{}, nil
	}
	return fake.importXMLFunc(ctx, xmlPath)
}

type fakeRuntimeEventPublisher struct {
	publishedCompleted []service.MasterDictionaryImportCompletedPayload
}

func (fake *fakeRuntimeEventPublisher) PublishImportProgress(_ context.Context, _ int) {
	// Import progress is irrelevant for these usecase tests.
}

func (fake *fakeRuntimeEventPublisher) PublishImportCompleted(_ context.Context, payload service.MasterDictionaryImportCompletedPayload) {
	fake.publishedCompleted = append(fake.publishedCompleted, payload)
}

func TestMasterDictionaryUsecaseGetPageSelectsPreferredEntry(t *testing.T) {
	ctx := context.Background()
	preferredID := int64(22)
	searchCalled := false

	usecase := NewMasterDictionaryUsecase(
		fakeQueryService{
			searchEntriesFunc: func(callCtx context.Context, query service.MasterDictionaryQuery) (service.MasterDictionaryListResult, error) {
				searchCalled = true
				if callCtx != ctx {
					t.Fatal("expected query service to receive original context")
				}
				if query.SearchTerm != "auriel" || query.Category != "書籍" || query.Page != 2 || query.PageSize != 50 {
					t.Fatalf("unexpected query: %#v", query)
				}
				return service.MasterDictionaryListResult{
					Items:      []service.MasterDictionaryEntry{{ID: 11}, {ID: preferredID}},
					TotalCount: 2,
					Page:       2,
					PageSize:   50,
				}, nil
			},
		},
		fakeCommandService{},
		fakeImportService{},
		nil,
	)

	page, err := usecase.GetPage(ctx, MasterDictionaryRefreshQuery{
		SearchTerm: "auriel",
		Category:   "書籍",
		Page:       2,
		PageSize:   50,
	}, &preferredID)
	if err != nil {
		t.Fatalf("expected get page to succeed: %v", err)
	}
	if !searchCalled {
		t.Fatal("expected query service to be called")
	}
	if page.SelectedID == nil || *page.SelectedID != preferredID {
		t.Fatalf("expected selected id %d, got %#v", preferredID, page.SelectedID)
	}
}

func TestMasterDictionaryUsecaseCreateEntryRefreshesWithChangedEntryID(t *testing.T) {
	ctx := context.Background()
	createdAt := time.Date(2026, 4, 14, 12, 0, 0, 0, time.UTC)
	refreshCalled := false

	usecase := NewMasterDictionaryUsecase(
		fakeQueryService{
			searchEntriesFunc: func(_ context.Context, query service.MasterDictionaryQuery) (service.MasterDictionaryListResult, error) {
				refreshCalled = true
				if query.Page != 3 || query.PageSize != 10 {
					t.Fatalf("unexpected refresh query: %#v", query)
				}
				return service.MasterDictionaryListResult{
					Items:      []service.MasterDictionaryEntry{{ID: 41, UpdatedAt: createdAt}},
					TotalCount: 1,
					Page:       3,
					PageSize:   10,
				}, nil
			},
		},
		fakeCommandService{
			createEntryFunc: func(callCtx context.Context, input service.MasterDictionaryMutationInput) (service.MasterDictionaryEntry, error) {
				if callCtx != ctx {
					t.Fatal("expected create to receive original context")
				}
				if input.Source != "Source" || input.Translation != "訳語" || input.Category != "カテゴリ" || input.Origin != "手動登録" || input.REC != usecaseBookREC || input.EDID != "BookAuriel" {
					t.Fatalf("unexpected create input: %#v", input)
				}
				return service.MasterDictionaryEntry{ID: 41, Source: input.Source, Translation: input.Translation, UpdatedAt: createdAt}, nil
			},
		},
		fakeImportService{},
		nil,
	)

	result, err := usecase.CreateEntry(ctx, MasterDictionaryMutationInput{
		Source: "Source", Translation: "訳語", Category: "カテゴリ", Origin: "手動登録", REC: usecaseBookREC, EDID: "BookAuriel",
	}, MasterDictionaryRefreshQuery{Page: 3, PageSize: 10})
	if err != nil {
		t.Fatalf("expected create entry to succeed: %v", err)
	}
	if !refreshCalled {
		t.Fatal("expected create to trigger refresh")
	}
	if result.ChangedEntry == nil || result.ChangedEntry.ID != 41 {
		t.Fatalf("expected changed entry id 41, got %#v", result.ChangedEntry)
	}
	if result.Page.SelectedID == nil || *result.Page.SelectedID != 41 {
		t.Fatalf("expected selected id 41 after create, got %#v", result.Page.SelectedID)
	}
}

func TestMasterDictionaryUsecaseDeleteEntryReturnsDeletedID(t *testing.T) {
	ctx := context.Background()
	deletedID := int64(72)
	deleteCalled := false

	usecase := NewMasterDictionaryUsecase(
		fakeQueryService{
			searchEntriesFunc: func(_ context.Context, query service.MasterDictionaryQuery) (service.MasterDictionaryListResult, error) {
				if query.SearchTerm != "after delete" || query.Page != 4 || query.PageSize != 20 {
					t.Fatalf("unexpected delete refresh query: %#v", query)
				}
				return service.MasterDictionaryListResult{
					Items:      []service.MasterDictionaryEntry{{ID: 99}},
					TotalCount: 1,
					Page:       4,
					PageSize:   20,
				}, nil
			},
		},
		fakeCommandService{
			deleteEntryFunc: func(callCtx context.Context, id int64) error {
				deleteCalled = true
				if callCtx != ctx {
					t.Fatal("expected delete to receive original context")
				}
				if id != deletedID {
					t.Fatalf("expected delete id %d, got %d", deletedID, id)
				}
				return nil
			},
		},
		fakeImportService{},
		nil,
	)

	result, err := usecase.DeleteEntry(ctx, deletedID, MasterDictionaryRefreshQuery{SearchTerm: "after delete", Page: 4, PageSize: 20})
	if err != nil {
		t.Fatalf("expected delete entry to succeed: %v", err)
	}
	if !deleteCalled {
		t.Fatal("expected delete service to be called")
	}
	if result.DeletedEntryID == nil || *result.DeletedEntryID != deletedID {
		t.Fatalf("expected deleted id %d, got %#v", deletedID, result.DeletedEntryID)
	}
	if result.Page.SelectedID == nil || *result.Page.SelectedID != 99 {
		t.Fatalf("expected next selected id 99, got %#v", result.Page.SelectedID)
	}
}

func TestMasterDictionaryUsecaseImportXMLPublishesCompletedEvent(t *testing.T) {
	ctx := context.Background()
	publisher := &fakeRuntimeEventPublisher{}
	importedAt := time.Date(2026, 4, 14, 13, 0, 0, 0, time.UTC)
	searchCalls := newFakeImportSearchSequence(t, importedAt)

	usecase := NewMasterDictionaryUsecase(
		fakeQueryService{searchEntriesFunc: searchCalls.next},
		fakeCommandService{},
		fakeImportService{importXMLFunc: newFakeImportXMLFunc(ctx, t)},
		publisher,
	)

	result, err := usecase.ImportXML(ctx, usecaseDictionaryXML, MasterDictionaryRefreshQuery{SearchTerm: "imported", Category: "書籍", Page: 5, PageSize: 15})
	if err != nil {
		t.Fatalf("expected import xml to succeed: %v", err)
	}
	if searchCalls.count != 2 {
		t.Fatalf("expected two page refresh calls, got %d", searchCalls.count)
	}
	if result.Page.SelectedID == nil || *result.Page.SelectedID != 88 {
		t.Fatalf("expected selected id 88, got %#v", result.Page.SelectedID)
	}
	assertCompletedImportPayload(t, publisher.publishedCompleted)
}

func TestMasterDictionaryUsecaseImportXMLReturnsEventBuildError(t *testing.T) {
	usecase := NewMasterDictionaryUsecase(
		fakeQueryService{
			searchEntriesFunc: func(_ context.Context, _ service.MasterDictionaryQuery) (service.MasterDictionaryListResult, error) {
				return service.MasterDictionaryListResult{}, errors.New("refresh failed")
			},
		},
		fakeCommandService{},
		fakeImportService{
			importXMLFunc: func(_ context.Context, _ string) (service.MasterDictionaryImportSummary, error) {
				return service.MasterDictionaryImportSummary{LastEntryID: 10}, nil
			},
		},
		&fakeRuntimeEventPublisher{},
	)

	_, err := usecase.ImportXML(context.Background(), usecaseDictionaryXML, MasterDictionaryRefreshQuery{})
	if err == nil {
		t.Fatal("expected import refresh error")
	}
}

type fakeImportSearchSequence struct {
	t        *testing.T
	count    int
	results  []service.MasterDictionaryListResult
	expected []service.MasterDictionaryQuery
}

func newFakeImportSearchSequence(t *testing.T, importedAt time.Time) *fakeImportSearchSequence {
	return &fakeImportSearchSequence{
		t: t,
		expected: []service.MasterDictionaryQuery{
			{SearchTerm: "imported", Category: "書籍", Page: 5, PageSize: 15},
			{SearchTerm: "", Category: masterDictionaryDefaultImportCategory, Page: masterDictionaryDefaultImportPage, PageSize: masterDictionaryDefaultImportPageSize},
		},
		results: []service.MasterDictionaryListResult{
			{Items: []service.MasterDictionaryEntry{{ID: 88, UpdatedAt: importedAt}}, TotalCount: 1, Page: 5, PageSize: 15},
			{Items: []service.MasterDictionaryEntry{{ID: 88, Source: "Auriel", UpdatedAt: importedAt}}, TotalCount: 1, Page: masterDictionaryDefaultImportPage, PageSize: masterDictionaryDefaultImportPageSize},
		},
	}
}

func (sequence *fakeImportSearchSequence) next(_ context.Context, query service.MasterDictionaryQuery) (service.MasterDictionaryListResult, error) {
	if sequence.count >= len(sequence.expected) {
		sequence.t.Fatalf("unexpected search call count: %d", sequence.count+1)
	}
	expected := sequence.expected[sequence.count]
	if query != expected {
		sequence.t.Fatalf("unexpected import refresh query[%d]: %#v", sequence.count, query)
	}
	result := sequence.results[sequence.count]
	sequence.count++
	return result, nil
}

func newFakeImportXMLFunc(expectedCtx context.Context, t *testing.T) func(context.Context, string) (service.MasterDictionaryImportSummary, error) {
	return func(callCtx context.Context, xmlPath string) (service.MasterDictionaryImportSummary, error) {
		if callCtx != expectedCtx {
			t.Fatal("expected import to receive original context")
		}
		if xmlPath != usecaseDictionaryXML {
			t.Fatalf("expected xml path %s, got %q", usecaseDictionaryXML, xmlPath)
		}
		return service.MasterDictionaryImportSummary{
			FilePath:      usecaseDictionaryXML,
			FileName:      usecaseDictionaryXML,
			ImportedCount: 3,
			UpdatedCount:  1,
			SkippedCount:  2,
			SelectedREC:   []string{usecaseBookREC},
			LastEntryID:   88,
		}, nil
	}
}

func assertCompletedImportPayload(t *testing.T, payloads []service.MasterDictionaryImportCompletedPayload) {
	t.Helper()

	if len(payloads) != 1 {
		t.Fatalf("expected one completed event, got %d", len(payloads))
	}
	payload := payloads[0]
	if payload.Refresh.Category != masterDictionaryDefaultImportCategory || payload.Refresh.Page != masterDictionaryDefaultImportPage || payload.Refresh.PageSize != masterDictionaryDefaultImportPageSize {
		t.Fatalf("unexpected refresh payload: %#v", payload.Refresh)
	}
	if payload.Refresh.RefreshTargetID == nil || *payload.Refresh.RefreshTargetID != 88 {
		t.Fatalf("expected refresh target id 88, got %#v", payload.Refresh.RefreshTargetID)
	}
	if payload.Summary.FileName != usecaseDictionaryXML || payload.Page.SelectedID == nil || *payload.Page.SelectedID != 88 {
		t.Fatalf("unexpected completed payload: %#v", payload)
	}
}
