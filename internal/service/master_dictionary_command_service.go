package service

import "context"

// MasterDictionaryCommandService provides create/update/delete operations.
type MasterDictionaryCommandService struct {
	core *MasterDictionaryService
}

// NewMasterDictionaryCommandService creates a command service.
func NewMasterDictionaryCommandService(core *MasterDictionaryService) *MasterDictionaryCommandService {
	return &MasterDictionaryCommandService{core: core}
}

// CreateEntry inserts a dictionary entry.
func (service *MasterDictionaryCommandService) CreateEntry(
	ctx context.Context,
	input MasterDictionaryMutationInput,
) (MasterDictionaryEntry, error) {
	return service.core.CreateEntry(ctx, input)
}

// UpdateEntry updates one dictionary entry.
func (service *MasterDictionaryCommandService) UpdateEntry(
	ctx context.Context,
	id int64,
	input MasterDictionaryMutationInput,
) (MasterDictionaryEntry, error) {
	return service.core.UpdateEntry(ctx, id, input)
}

// DeleteEntry removes one dictionary entry.
func (service *MasterDictionaryCommandService) DeleteEntry(ctx context.Context, id int64) error {
	return service.core.DeleteEntry(ctx, id)
}
