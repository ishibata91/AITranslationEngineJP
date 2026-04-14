// Package bootstrap wires the default backend graph outside the controller layer.
package bootstrap

import (
	"context"
	"errors"
	"fmt"
	"time"

	controllerwails "aitranslationenginejp/internal/controller/wails"
	"aitranslationenginejp/internal/repository"
	"aitranslationenginejp/internal/service"
	"aitranslationenginejp/internal/usecase"
)

// NewAppController builds the default backend graph for the desktop app.
func NewAppController() *controllerwails.AppController {
	now := func() time.Time { return time.Now().UTC() }
	runtimeEmitterState := controllerwails.NewRuntimeEmitterState()
	runtimePublisher := usecase.NewWailsMasterDictionaryRuntimeEventPublisher(runtimeEmitterState.RuntimeEventContext)
	repositoryAdapter := newMasterDictionaryServiceRepositoryAdapter(
		repository.NewInMemoryMasterDictionaryRepository(repository.DefaultMasterDictionarySeed(now())),
	)
	queryService := service.NewMasterDictionaryQueryService(repositoryAdapter)
	commandService := service.NewMasterDictionaryCommandService(repositoryAdapter, now)
	importService := service.NewMasterDictionaryImportService(
		repositoryAdapter,
		service.NewLocalMasterDictionaryXMLFilePort(),
		service.NewXMLDecoderMasterDictionaryRecordReader(),
		usecase.NewImportProgressEmitter(runtimePublisher),
		now,
	)
	masterDictionaryUsecase := usecase.NewMasterDictionaryUsecase(
		queryService,
		commandService,
		importService,
		runtimePublisher,
	)
	masterDictionaryController := controllerwails.NewMasterDictionaryController(
		masterDictionaryUsecase,
		runtimeEmitterState,
	)
	return controllerwails.NewAppController(masterDictionaryController)
}

type masterDictionaryServiceRepositoryAdapter struct {
	repository *repository.InMemoryMasterDictionaryRepository
}

func newMasterDictionaryServiceRepositoryAdapter(
	repository *repository.InMemoryMasterDictionaryRepository,
) service.RepositoryPort {
	return masterDictionaryServiceRepositoryAdapter{repository: repository}
}

func (adapter masterDictionaryServiceRepositoryAdapter) List(
	ctx context.Context,
	query service.MasterDictionaryQuery,
) (service.MasterDictionaryListResult, error) {
	result, err := adapter.repository.List(ctx, repository.MasterDictionaryListQuery{
		SearchTerm: query.SearchTerm,
		Category:   query.Category,
		Page:       query.Page,
		PageSize:   query.PageSize,
	})
	if err != nil {
		return service.MasterDictionaryListResult{}, fmt.Errorf("list repository entries: %w", err)
	}
	return toServiceListResult(result), nil
}

func (adapter masterDictionaryServiceRepositoryAdapter) GetByID(
	ctx context.Context,
	id int64,
) (service.MasterDictionaryEntry, error) {
	entry, err := adapter.repository.GetByID(ctx, id)
	if err != nil {
		return service.MasterDictionaryEntry{}, mapRepositoryError(err, id)
	}
	return toServiceEntry(entry), nil
}

func (adapter masterDictionaryServiceRepositoryAdapter) Create(
	ctx context.Context,
	draft service.MasterDictionaryDraft,
) (service.MasterDictionaryEntry, error) {
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
		return service.MasterDictionaryEntry{}, fmt.Errorf("create repository entry: %w", err)
	}
	return toServiceEntry(entry), nil
}

func (adapter masterDictionaryServiceRepositoryAdapter) Update(
	ctx context.Context,
	id int64,
	draft service.MasterDictionaryDraft,
) (service.MasterDictionaryEntry, error) {
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
		return service.MasterDictionaryEntry{}, fmt.Errorf("update repository entry: %w", err)
	}
	return toServiceEntry(entry), nil
}

func (adapter masterDictionaryServiceRepositoryAdapter) Delete(
	ctx context.Context,
	id int64,
) error {
	if err := adapter.repository.Delete(ctx, id); err != nil {
		return mapRepositoryError(err, id)
	}
	return nil
}

func (adapter masterDictionaryServiceRepositoryAdapter) UpsertBySourceAndREC(
	ctx context.Context,
	record service.MasterDictionaryImportRecord,
) (service.MasterDictionaryEntry, bool, error) {
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
		return service.MasterDictionaryEntry{}, false, fmt.Errorf("upsert repository import record: %w", err)
	}
	return toServiceEntry(entry), created, nil
}

func mapRepositoryError(err error, id int64) error {
	if errors.Is(err, repository.ErrMasterDictionaryEntryNotFound) {
		return fmt.Errorf("%w: id=%d", service.ErrMasterDictionaryEntryNotFound, id)
	}
	return err
}

func toServiceListResult(result repository.MasterDictionaryListResult) service.MasterDictionaryListResult {
	items := make([]service.MasterDictionaryEntry, 0, len(result.Items))
	for _, item := range result.Items {
		items = append(items, toServiceEntry(item))
	}
	return service.MasterDictionaryListResult{
		Items:      items,
		TotalCount: result.TotalCount,
		Page:       result.Page,
		PageSize:   result.PageSize,
	}
}

func toServiceEntry(entry repository.MasterDictionaryEntry) service.MasterDictionaryEntry {
	return service.MasterDictionaryEntry{
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
