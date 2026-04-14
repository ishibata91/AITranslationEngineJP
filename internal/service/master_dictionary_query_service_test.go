package service

import (
	"context"
	"errors"
	"testing"
)

const queryServiceWhiterunSource = "Whiterun 01"

func TestMasterDictionaryQueryServiceSearchEntriesTrimsQueryBeforeRepositoryCall(t *testing.T) {
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
		SearchTerm: "  " + queryServiceWhiterunSource + "  ",
		Category:   " 地名 ",
		Page:       1,
		PageSize:   10,
	})
	if err != nil {
		t.Fatalf("expected search to succeed: %v", err)
	}
	if len(repo.listQueries) != 1 {
		t.Fatalf("expected one repository list call, got %d", len(repo.listQueries))
	}
	if repo.listQueries[0].SearchTerm != "Whiterun 01" {
		t.Fatalf("expected trimmed search term, got %q", repo.listQueries[0].SearchTerm)
	}
	if repo.listQueries[0].Category != "地名" {
		t.Fatalf("expected trimmed category, got %q", repo.listQueries[0].Category)
	}
	if result.TotalCount != 1 || len(result.Items) != 1 {
		t.Fatalf("expected one search result, got %+v", result)
	}
	if result.Items[0].Source != queryServiceWhiterunSource {
		t.Fatalf("expected returned item to be preserved, got %q", result.Items[0].Source)
	}
}

func TestMasterDictionaryQueryServiceLoadEntryDetailMapsNotFound(t *testing.T) {
	repo := &repositoryStub{
		getByIDFunc: func(_ context.Context, id int64) (MasterDictionaryEntry, error) {
			if id == 999 {
				return MasterDictionaryEntry{}, ErrMasterDictionaryEntryNotFound
			}
			return MasterDictionaryEntry{ID: id, Source: "entry"}, nil
		},
	}
	service := NewMasterDictionaryQueryService(repo)

	if _, err := service.LoadEntryDetail(context.Background(), 0); err == nil {
		t.Fatal("expected invalid id to fail")
	}

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
