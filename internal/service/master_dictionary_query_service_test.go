package service

import (
	"context"
	"errors"
	"testing"
)

const (
	queryServiceWhiterunSource = "Whiterun 01"
	errQuerySearchSucceeded    = "expected search to succeed: %v"
)

func TestMasterDictionaryQueryServiceSearchEntriesTrimsSearchTermBeforeRepositoryCall(t *testing.T) {
	repo := &repositoryStub{
		listFunc: func(_ context.Context, query MasterDictionaryQuery) (MasterDictionaryListResult, error) {
			return MasterDictionaryListResult{Page: query.Page, PageSize: query.PageSize}, nil
		},
	}
	service := NewMasterDictionaryQueryService(repo)

	_, err := service.SearchEntries(context.Background(), MasterDictionaryQuery{
		SearchTerm: "  " + queryServiceWhiterunSource + "  ",
		Category:   " 地名 ",
		Page:       1,
		PageSize:   10,
	})
	if err != nil {
		t.Fatalf(errQuerySearchSucceeded, err)
	}

	if repo.listQueries[0].SearchTerm != "Whiterun 01" {
		t.Fatalf("expected trimmed search term, got %q", repo.listQueries[0].SearchTerm)
	}
}

func TestMasterDictionaryQueryServiceSearchEntriesTrimsCategoryBeforeRepositoryCall(t *testing.T) {
	repo := &repositoryStub{
		listFunc: func(_ context.Context, query MasterDictionaryQuery) (MasterDictionaryListResult, error) {
			return MasterDictionaryListResult{Page: query.Page, PageSize: query.PageSize}, nil
		},
	}
	service := NewMasterDictionaryQueryService(repo)

	_, err := service.SearchEntries(context.Background(), MasterDictionaryQuery{
		SearchTerm: "  " + queryServiceWhiterunSource + "  ",
		Category:   " 地名 ",
		Page:       1,
		PageSize:   10,
	})
	if err != nil {
		t.Fatalf(errQuerySearchSucceeded, err)
	}

	if repo.listQueries[0].Category != "地名" {
		t.Fatalf("expected trimmed category, got %q", repo.listQueries[0].Category)
	}
}

func TestMasterDictionaryQueryServiceSearchEntriesReturnsRepositoryResult(t *testing.T) {
	repo := &repositoryStub{
		listFunc: func(_ context.Context, query MasterDictionaryQuery) (MasterDictionaryListResult, error) {
			return MasterDictionaryListResult{
				Items:      []MasterDictionaryEntry{{ID: 1, Source: queryServiceWhiterunSource}},
				TotalCount: 1,
				Page:       query.Page,
				PageSize:   query.PageSize,
			}, nil
		},
	}
	service := NewMasterDictionaryQueryService(repo)

	result, err := service.SearchEntries(context.Background(), MasterDictionaryQuery{
		SearchTerm: queryServiceWhiterunSource,
		Category:   "地名",
		Page:       1,
		PageSize:   10,
	})
	if err != nil {
		t.Fatalf(errQuerySearchSucceeded, err)
	}

	if result.Items[0].Source != queryServiceWhiterunSource {
		t.Fatalf("expected returned item to be preserved, got %q", result.Items[0].Source)
	}
}

func TestMasterDictionaryQueryServiceLoadEntryDetailRejectsInvalidID(t *testing.T) {
	service := NewMasterDictionaryQueryService(&repositoryStub{})

	_, err := service.LoadEntryDetail(context.Background(), 0)
	if err == nil {
		t.Fatal("expected invalid id to fail")
	}
}

func TestMasterDictionaryQueryServiceLoadEntryDetailMapsNotFound(t *testing.T) {
	repo := &repositoryStub{
		getByIDFunc: func(_ context.Context, _ int64) (MasterDictionaryEntry, error) {
			return MasterDictionaryEntry{}, ErrMasterDictionaryEntryNotFound
		},
	}
	service := NewMasterDictionaryQueryService(repo)

	_, err := service.LoadEntryDetail(context.Background(), 999)
	if err == nil {
		t.Fatal("expected missing entry to fail")
	}
	if !IsNotFoundError(err) {
		t.Fatal("expected missing entry to map to service not found error")
	}
	if !errors.Is(err, ErrMasterDictionaryEntryNotFound) {
		t.Fatal("expected wrapped error to preserve not found sentinel")
	}
}
