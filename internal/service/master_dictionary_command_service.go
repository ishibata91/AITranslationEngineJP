package service

import (
	"context"
	"fmt"
	"time"
)

// MasterDictionaryCommandService provides create/update/delete operations.
type MasterDictionaryCommandService struct {
	repository RepositoryPort
	now        func() time.Time
}

// NewMasterDictionaryCommandService creates a command service.
func NewMasterDictionaryCommandService(
	repository RepositoryPort,
	now func() time.Time,
) *MasterDictionaryCommandService {
	return &MasterDictionaryCommandService{
		repository: repository,
		now:        normalizeClock(now),
	}
}

// CreateEntry inserts a dictionary entry.
func (service *MasterDictionaryCommandService) CreateEntry(
	ctx context.Context,
	input MasterDictionaryMutationInput,
) (MasterDictionaryEntry, error) {
	draft, err := validateMutationInput(input, service.now)
	if err != nil {
		return MasterDictionaryEntry{}, err
	}

	created, err := service.repository.Create(ctx, draft)
	if err != nil {
		return MasterDictionaryEntry{}, fmt.Errorf("create master dictionary entry: %w", err)
	}
	return created, nil
}

// UpdateEntry updates one dictionary entry.
func (service *MasterDictionaryCommandService) UpdateEntry(
	ctx context.Context,
	id int64,
	input MasterDictionaryMutationInput,
) (MasterDictionaryEntry, error) {
	if err := validateMasterDictionaryID(id); err != nil {
		return MasterDictionaryEntry{}, err
	}

	draft, err := validateMutationInput(input, service.now)
	if err != nil {
		return MasterDictionaryEntry{}, err
	}

	updated, err := service.repository.Update(ctx, id, draft)
	if err != nil {
		return MasterDictionaryEntry{}, fmt.Errorf("update master dictionary entry: %w", err)
	}
	return updated, nil
}

// DeleteEntry removes one dictionary entry.
func (service *MasterDictionaryCommandService) DeleteEntry(ctx context.Context, id int64) error {
	if err := validateMasterDictionaryID(id); err != nil {
		return err
	}

	if err := service.repository.Delete(ctx, id); err != nil {
		return fmt.Errorf("delete master dictionary entry: %w", err)
	}
	return nil
}
