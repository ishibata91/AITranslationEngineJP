package service

import (
	"context"
	"errors"
	"fmt"

	"aitranslationenginejp/internal/repository"
)

// NewSQLiteMasterDictionaryRepositoryPort builds the SQLite-backed repository port for wiring.
func NewSQLiteMasterDictionaryRepositoryPort(
	ctx context.Context,
	databasePath string,
	seed []repository.MasterDictionaryEntry,
) (RepositoryPort, error) {
	repositoryAdapter, err := repository.NewSQLiteMasterDictionaryRepository(ctx, databasePath, seed)
	if err != nil {
		return nil, fmt.Errorf("build sqlite repository port: %w", err)
	}
	return sqliteMasterDictionaryRepositoryPort{repository: repositoryAdapter}, nil
}

// SQLiteMasterDictionaryRepositoryPortCloser extracts a shutdown closer from the SQLite-backed port.
func SQLiteMasterDictionaryRepositoryPortCloser(port RepositoryPort) func(context.Context) error {
	closer, ok := port.(interface{ Close() error })
	if !ok {
		return nil
	}
	return func(context.Context) error {
		if err := closer.Close(); err != nil {
			return fmt.Errorf("close sqlite repository port: %w", err)
		}
		return nil
	}
}

type sqliteMasterDictionaryRepositoryPort struct {
	repository interface {
		repository.MasterDictionaryRepository
		Close() error
	}
}

func (adapter sqliteMasterDictionaryRepositoryPort) List(
	ctx context.Context,
	query MasterDictionaryQuery,
) (MasterDictionaryListResult, error) {
	result, err := adapter.repository.List(ctx, repository.MasterDictionaryListQuery{
		SearchTerm: query.SearchTerm,
		Category:   query.Category,
		Page:       query.Page,
		PageSize:   query.PageSize,
	})
	if err != nil {
		return MasterDictionaryListResult{}, fmt.Errorf("list repository entries: %w", err)
	}
	return toServiceListResult(result), nil
}

func (adapter sqliteMasterDictionaryRepositoryPort) GetByID(
	ctx context.Context,
	id int64,
) (MasterDictionaryEntry, error) {
	entry, err := adapter.repository.GetByID(ctx, id)
	if err != nil {
		return MasterDictionaryEntry{}, mapRepositoryError(err, id)
	}
	return toServiceEntry(entry), nil
}

func (adapter sqliteMasterDictionaryRepositoryPort) Create(
	ctx context.Context,
	draft MasterDictionaryDraft,
) (MasterDictionaryEntry, error) {
	entry, err := adapter.repository.Create(ctx, repository.MasterDictionaryDraft{
		Source:      draft.Source,
		Translation: draft.Translation,
		Category:    draft.Category,
		Origin:      draft.Origin,
		REC:         draft.REC,
		EDID:        draft.EDID,
		UpdatedAt:   draft.UpdatedAt,
	})
	if err != nil {
		return MasterDictionaryEntry{}, fmt.Errorf("create repository entry: %w", err)
	}
	return toServiceEntry(entry), nil
}

func (adapter sqliteMasterDictionaryRepositoryPort) Update(
	ctx context.Context,
	id int64,
	draft MasterDictionaryDraft,
) (MasterDictionaryEntry, error) {
	entry, err := adapter.repository.Update(ctx, id, repository.MasterDictionaryDraft{
		Source:      draft.Source,
		Translation: draft.Translation,
		Category:    draft.Category,
		Origin:      draft.Origin,
		REC:         draft.REC,
		EDID:        draft.EDID,
		UpdatedAt:   draft.UpdatedAt,
	})
	if err != nil {
		return MasterDictionaryEntry{}, fmt.Errorf("update repository entry: %w", err)
	}
	return toServiceEntry(entry), nil
}

func (adapter sqliteMasterDictionaryRepositoryPort) Delete(
	ctx context.Context,
	id int64,
) error {
	if err := adapter.repository.Delete(ctx, id); err != nil {
		return mapRepositoryError(err, id)
	}
	return nil
}

func (adapter sqliteMasterDictionaryRepositoryPort) UpsertBySourceAndREC(
	ctx context.Context,
	record MasterDictionaryImportRecord,
) (MasterDictionaryEntry, bool, error) {
	entry, created, err := adapter.repository.UpsertBySourceAndREC(ctx, repository.MasterDictionaryImportRecord{
		Source:      record.Source,
		Translation: record.Translation,
		REC:         record.REC,
		EDID:        record.EDID,
		Category:    record.Category,
		Origin:      record.Origin,
		UpdatedAt:   record.UpdatedAt,
	})
	if err != nil {
		return MasterDictionaryEntry{}, false, fmt.Errorf("upsert repository import record: %w", err)
	}
	return toServiceEntry(entry), created, nil
}

func (adapter sqliteMasterDictionaryRepositoryPort) Close() error {
	if err := adapter.repository.Close(); err != nil {
		return fmt.Errorf("close repository entry adapter: %w", err)
	}
	return nil
}

func mapRepositoryError(err error, id int64) error {
	if errors.Is(err, repository.ErrMasterDictionaryEntryNotFound) {
		return fmt.Errorf("%w: id=%d", ErrMasterDictionaryEntryNotFound, id)
	}
	return err
}

func toServiceListResult(result repository.MasterDictionaryListResult) MasterDictionaryListResult {
	items := make([]MasterDictionaryEntry, 0, len(result.Items))
	for _, item := range result.Items {
		items = append(items, toServiceEntry(item))
	}
	return MasterDictionaryListResult{
		Items:      items,
		TotalCount: result.TotalCount,
		Page:       result.Page,
		PageSize:   result.PageSize,
	}
}

func toServiceEntry(entry repository.MasterDictionaryEntry) MasterDictionaryEntry {
	return MasterDictionaryEntry{
		ID:          entry.ID,
		Source:      entry.Source,
		Translation: entry.Translation,
		Category:    entry.Category,
		Origin:      entry.Origin,
		REC:         entry.REC,
		EDID:        entry.EDID,
		UpdatedAt:   entry.UpdatedAt,
	}
}
