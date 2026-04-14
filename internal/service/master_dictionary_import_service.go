package service

import (
	"context"
	"fmt"
	"strings"
	"time"
)

// MasterDictionaryImportService provides XML import operations.
type MasterDictionaryImportService struct {
	repository RepositoryPort
	xmlFiles   XMLFilePort
	xmlRecords XMLRecordReaderPort
	runtime    RuntimeContextPort
	now        func() time.Time
}

// NewMasterDictionaryImportService creates an import service.
func NewMasterDictionaryImportService(
	repository RepositoryPort,
	xmlFiles XMLFilePort,
	xmlRecords XMLRecordReaderPort,
	runtime RuntimeContextPort,
	now func() time.Time,
) *MasterDictionaryImportService {
	return &MasterDictionaryImportService{
		repository: repository,
		xmlFiles:   xmlFiles,
		xmlRecords: xmlRecords,
		runtime:    runtime,
		now:        normalizeClock(now),
	}
}

// ImportXML imports an XML file to master dictionary.
func (service *MasterDictionaryImportService) ImportXML(
	ctx context.Context,
	xmlPath string,
) (MasterDictionaryImportSummary, error) {
	resolvedPath, err := service.xmlFiles.ResolvePath(xmlPath)
	if err != nil {
		return MasterDictionaryImportSummary{}, fmt.Errorf("resolve xml path: %w", err)
	}

	totalRecords, err := service.countStringRecords(resolvedPath)
	if err != nil {
		return MasterDictionaryImportSummary{}, err
	}
	service.emitImportProgress(ctx, 0)

	reader, err := service.xmlFiles.Open(resolvedPath)
	if err != nil {
		return MasterDictionaryImportSummary{}, fmt.Errorf("open xml file: %w", err)
	}
	defer func() {
		_ = reader.Close()
	}()

	counters := masterDictionaryImportCounters{}
	processedCount := 0
	readErr := service.xmlRecords.ReadStringRecords(reader, func(record xmlStringRecord) error {
		if err := service.importXMLRecord(ctx, record, &counters); err != nil {
			return err
		}
		processedCount++
		service.emitImportProgress(ctx, normalizeImportProgress(processedCount, totalRecords))
		return nil
	})
	if readErr != nil {
		return MasterDictionaryImportSummary{}, fmt.Errorf("read xml records: %w", readErr)
	}

	service.emitImportProgress(ctx, 100)
	return MasterDictionaryImportSummary{
		FilePath:      resolvedPath,
		FileName:      baseName(resolvedPath),
		ImportedCount: counters.importedCount,
		UpdatedCount:  counters.updatedCount,
		SkippedCount:  counters.skippedCount,
		SelectedREC:   allowedRECList(),
		LastEntryID:   counters.lastEntryID,
	}, nil
}

func (service *MasterDictionaryImportService) countStringRecords(resolvedPath string) (int, error) {
	reader, err := service.xmlFiles.Open(resolvedPath)
	if err != nil {
		return 0, fmt.Errorf("open xml file: %w", err)
	}
	defer func() {
		_ = reader.Close()
	}()

	totalRecords, err := service.xmlRecords.CountStringRecords(reader)
	if err != nil {
		return 0, fmt.Errorf("count xml string records: %w", err)
	}
	return totalRecords, nil
}

func (service *MasterDictionaryImportService) importXMLRecord(
	ctx context.Context,
	record xmlStringRecord,
	counters *masterDictionaryImportCounters,
) error {
	rec := strings.TrimSpace(record.REC)
	if !isAllowedImportREC(rec) {
		counters.skippedCount++
		return nil
	}

	source := strings.TrimSpace(record.Source)
	translation := strings.TrimSpace(record.Dest)
	if source == "" || translation == "" {
		counters.skippedCount++
		return nil
	}

	entry, created, err := service.repository.UpsertBySourceAndREC(ctx, MasterDictionaryImportRecord{
		Source:      source,
		Translation: translation,
		REC:         rec,
		EDID:        strings.TrimSpace(record.EDID),
		Category:    categoryFromREC(rec),
		Origin:      masterDictionaryImportOrigin,
		UpdatedAt:   service.now().UTC(),
	})
	if err != nil {
		return fmt.Errorf("upsert imported record: %w", err)
	}

	counters.lastEntryID = entry.ID
	if created {
		counters.importedCount++
	} else {
		counters.updatedCount++
	}
	return nil
}

func (service *MasterDictionaryImportService) emitImportProgress(ctx context.Context, progress int) {
	if service.runtime == nil {
		return
	}
	service.runtime.EmitImportProgress(ctx, progress)
}
