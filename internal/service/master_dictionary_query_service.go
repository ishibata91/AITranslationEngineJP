package service

import "context"

// MasterDictionaryQueryService provides read-only operations for master dictionary entries.
type MasterDictionaryQueryService struct {
	core *MasterDictionaryService
}

// NewMasterDictionaryQueryService creates a query service.
func NewMasterDictionaryQueryService(core *MasterDictionaryService) *MasterDictionaryQueryService {
	return &MasterDictionaryQueryService{core: core}
}

// SearchEntries returns filtered and paged dictionary entries.
func (service *MasterDictionaryQueryService) SearchEntries(
	ctx context.Context,
	query MasterDictionaryQuery,
) (MasterDictionaryListResult, error) {
	return service.core.ListEntries(ctx, query)
}

// LoadEntryDetail returns one dictionary entry by id.
func (service *MasterDictionaryQueryService) LoadEntryDetail(
	ctx context.Context,
	id int64,
) (MasterDictionaryEntry, error) {
	return service.core.GetEntry(ctx, id)
}
