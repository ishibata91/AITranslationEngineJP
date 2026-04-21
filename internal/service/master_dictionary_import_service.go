package service

import (
	"context"
	"fmt"
	"strings"
	"time"
)

// MasterDictionaryImportService provides XML import operations.
type MasterDictionaryImportService struct {
	repository     RepositoryPort
	xmlFiles       XMLFilePort
	xmlRecords     XMLRecordReaderPort
	runtime        RuntimeContextPort
	foundationData FoundationDataPort
	now            func() time.Time
}

// WithFoundationData wires provenance persistence into the import service.
func (service *MasterDictionaryImportService) WithFoundationData(fd FoundationDataPort) *MasterDictionaryImportService {
	service.foundationData = fd
	return service
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

	// Buffer all records in one read pass so provenance can be created before any upsert.
	records, err := service.readAllStringRecords(resolvedPath)
	if err != nil {
		return MasterDictionaryImportSummary{}, err
	}

	// Create provenance BEFORE the import loop so entries can carry xtranslator_translation_xml_id.
	// If provenance creation fails here, no entries have been committed yet.
	provenanceID, err := service.persistXMLProvenance(ctx, resolvedPath, countImportableRecords(records))
	if err != nil {
		return MasterDictionaryImportSummary{}, err
	}
	var xmlID *int64
	if provenanceID != 0 {
		xmlID = &provenanceID
	}

	counters := masterDictionaryImportCounters{}
	for i, record := range records {
		if err := service.importXMLRecord(ctx, record, xmlID, &counters); err != nil {
			return MasterDictionaryImportSummary{}, fmt.Errorf("read xml records: %w", err)
		}
		service.emitImportProgress(ctx, normalizeImportProgress(i+1, totalRecords))
	}

	service.emitImportProgress(ctx, 100)
	return MasterDictionaryImportSummary{
		FilePath:      resolvedPath,
		FileName:      baseName(resolvedPath),
		ImportedCount: counters.importedCount,
		UpdatedCount:  counters.updatedCount,
		SkippedCount:  counters.skippedCount,
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
	xmlID *int64,
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
		Source:                      source,
		Translation:                 translation,
		REC:                         rec,
		EDID:                        strings.TrimSpace(record.EDID),
		Category:                    categoryFromREC(rec),
		Origin:                      masterDictionaryImportOrigin,
		UpdatedAt:                   service.now().UTC(),
		XTranslatorTranslationXMLID: xmlID,
	})
	if err != nil {
		return fmt.Errorf("upsert imported record: %w", err)
	}

	counters.lastEntryID = entry.ID
	if created {
		counters.trackImportedEntry(entry, record)
	} else {
		counters.trackUpdatedEntry(entry, record)
	}
	return nil
}

func (service *MasterDictionaryImportService) emitImportProgress(ctx context.Context, progress int) {
	if service.runtime == nil {
		return
	}
	service.runtime.EmitImportProgress(ctx, progress)
}

func (service *MasterDictionaryImportService) readAllStringRecords(resolvedPath string) ([]xmlStringRecord, error) {
	reader, err := service.xmlFiles.Open(resolvedPath)
	if err != nil {
		return nil, fmt.Errorf("open xml file: %w", err)
	}
	defer func() {
		_ = reader.Close()
	}()
	var records []xmlStringRecord
	if err := service.xmlRecords.ReadStringRecords(reader, func(record xmlStringRecord) error {
		records = append(records, record)
		return nil
	}); err != nil {
		return nil, fmt.Errorf("read xml records: %w", err)
	}
	return records, nil
}

func countImportableRecords(records []xmlStringRecord) int {
	count := 0
	for _, record := range records {
		rec := strings.TrimSpace(record.REC)
		if !isAllowedImportREC(rec) {
			continue
		}
		if strings.TrimSpace(record.Source) == "" || strings.TrimSpace(record.Dest) == "" {
			continue
		}
		count++
	}
	return count
}

func (service *MasterDictionaryImportService) persistXMLProvenance(
	ctx context.Context,
	resolvedPath string,
	termCount int,
) (int64, error) {
	if service.foundationData == nil {
		return 0, nil
	}
	draft := XMLProvenanceDraft{
		FilePath:         resolvedPath,
		TargetPluginName: pluginNameFromPath(resolvedPath),
		TargetPluginType: pluginTypeFromPath(resolvedPath),
		TermCount:        termCount,
		ImportedAt:       service.now().UTC(),
	}
	provenanceID, err := service.foundationData.CreateXTranslatorTranslationXML(ctx, draft)
	if err != nil {
		return 0, fmt.Errorf("persist xml provenance: %w", err)
	}
	return provenanceID, nil
}

func pluginNameFromPath(path string) string {
	name := baseName(path)
	if dot := strings.LastIndex(name, "."); dot > 0 {
		name = name[:dot]
	}
	for _, ext := range []string{".esm", ".esp", ".esl"} {
		if idx := strings.Index(strings.ToLower(name), ext); idx > 0 {
			return name[:idx]
		}
	}
	if idx := strings.Index(name, "_"); idx > 0 {
		return name[:idx]
	}
	return name
}

func pluginTypeFromPath(path string) string {
	name := strings.ToLower(baseName(path))
	for _, ext := range []string{"esm", "esp", "esl"} {
		if strings.Contains(name, "."+ext) {
			return ext
		}
	}
	return ""
}
