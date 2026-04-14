package service

import (
	"context"
	"fmt"
	"strings"
)

// MasterDictionaryQueryService provides read-only operations for master dictionary entries.
type MasterDictionaryQueryService struct {
	repository RepositoryPort
}

// NewMasterDictionaryQueryService creates a query service.
func NewMasterDictionaryQueryService(repository RepositoryPort) *MasterDictionaryQueryService {
	return &MasterDictionaryQueryService{repository: repository}
}

// SearchEntries returns filtered and paged dictionary entries.
func (service *MasterDictionaryQueryService) SearchEntries(
	ctx context.Context,
	query MasterDictionaryQuery,
) (MasterDictionaryListResult, error) {
	result, err := service.repository.List(ctx, MasterDictionaryQuery{
		SearchTerm: strings.TrimSpace(query.SearchTerm),
		Category:   strings.TrimSpace(query.Category),
		Page:       query.Page,
		PageSize:   query.PageSize,
	})
	if err != nil {
		return MasterDictionaryListResult{}, fmt.Errorf("list master dictionary entries: %w", err)
	}
	return result, nil
}

// LoadEntryDetail returns one dictionary entry by id.
func (service *MasterDictionaryQueryService) LoadEntryDetail(
	ctx context.Context,
	id int64,
) (MasterDictionaryEntry, error) {
	if err := validateMasterDictionaryID(id); err != nil {
		return MasterDictionaryEntry{}, err
	}

	entry, err := service.repository.GetByID(ctx, id)
	if err != nil {
		return MasterDictionaryEntry{}, fmt.Errorf("get master dictionary entry: %w", err)
	}
	return entry, nil
}
