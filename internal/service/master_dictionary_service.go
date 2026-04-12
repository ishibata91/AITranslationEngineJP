package service

import (
	"context"
	"encoding/xml"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"aitranslationenginejp/internal/repository"
)

var (
	// ErrMasterDictionaryValidation means the request payload is invalid.
	ErrMasterDictionaryValidation = errors.New("master dictionary validation error")
)

const (
	masterDictionaryDefaultCategory   = "固有名詞"
	masterDictionaryDefaultOrigin     = "手動登録"
	masterDictionaryImportOrigin      = "XML取込"
	masterDictionaryIDValidationError = "%w: id must be greater than zero"
	masterDictionaryOpenXMLFileError  = "open xml file: %w"
)

var allowedImportREC = map[string]struct{}{
	"BOOK:FULL": {},
	"NPC_:FULL": {},
	"NPC_:SHRT": {},
	"ARMO:FULL": {},
	"WEAP:FULL": {},
	"LCTN:FULL": {},
	"CELL:FULL": {},
	"CONT:FULL": {},
	"MISC:FULL": {},
	"ALCH:FULL": {},
	"FURN:FULL": {},
	"DOOR:FULL": {},
	"RACE:FULL": {},
	"INGR:FULL": {},
	"FLOR:FULL": {},
	"SHOU:FULL": {},
}

// MasterDictionaryEntry describes one dictionary record in service boundary.
type MasterDictionaryEntry struct {
	ID          int64
	Source      string
	Translation string
	Category    string
	Origin      string
	REC         string
	EDID        string
	UpdatedAt   time.Time
}

// MasterDictionaryQuery defines list/search conditions.
type MasterDictionaryQuery struct {
	SearchTerm string
	Category   string
	Page       int
	PageSize   int
}

// MasterDictionaryListResult is the service-layer list result.
type MasterDictionaryListResult struct {
	Items      []MasterDictionaryEntry
	TotalCount int
	Page       int
	PageSize   int
}

// MasterDictionaryMutationInput is an input payload for create or update.
type MasterDictionaryMutationInput struct {
	Source      string
	Translation string
	Category    string
	Origin      string
	REC         string
	EDID        string
}

// MasterDictionaryImportSummary returns import execution results.
type MasterDictionaryImportSummary struct {
	FilePath      string
	FileName      string
	ImportedCount int
	UpdatedCount  int
	SkippedCount  int
	SelectedREC   []string
	LastEntryID   int64
}

// MasterDictionaryService provides reusable master dictionary operations.
type MasterDictionaryService struct {
	repository         repository.MasterDictionaryRepository
	now                func() time.Time
	emitImportProgress func(context.Context, int)
}

// NewDefaultMasterDictionaryService creates a service with in-memory default seed.
func NewDefaultMasterDictionaryService(now func() time.Time) *MasterDictionaryService {
	clock := now
	if clock == nil {
		clock = time.Now
	}

	repo := repository.NewInMemoryMasterDictionaryRepository(
		repository.DefaultMasterDictionarySeed(clock().UTC()),
	)
	return NewMasterDictionaryService(repo, clock)
}

// NewMasterDictionaryService creates a service instance.
func NewMasterDictionaryService(repo repository.MasterDictionaryRepository, now func() time.Time) *MasterDictionaryService {
	clock := now
	if clock == nil {
		clock = time.Now
	}
	return &MasterDictionaryService{repository: repo, now: clock}
}

// SetImportProgressEmitter sets an optional runtime progress emitter for XML import.
func (service *MasterDictionaryService) SetImportProgressEmitter(emitter func(context.Context, int)) {
	service.emitImportProgress = emitter
}

// ListEntries returns filtered and paged dictionary entries.
func (service *MasterDictionaryService) ListEntries(ctx context.Context, query MasterDictionaryQuery) (MasterDictionaryListResult, error) {
	result, err := service.repository.List(ctx, repository.MasterDictionaryListQuery{
		SearchTerm: strings.TrimSpace(query.SearchTerm),
		Category:   strings.TrimSpace(query.Category),
		Page:       query.Page,
		PageSize:   query.PageSize,
	})
	if err != nil {
		return MasterDictionaryListResult{}, fmt.Errorf("list master dictionary entries: %w", err)
	}

	items := make([]MasterDictionaryEntry, 0, len(result.Items))
	for _, item := range result.Items {
		items = append(items, fromRepositoryEntry(item))
	}

	return MasterDictionaryListResult{
		Items:      items,
		TotalCount: result.TotalCount,
		Page:       result.Page,
		PageSize:   result.PageSize,
	}, nil
}

// GetEntry returns one dictionary entry.
func (service *MasterDictionaryService) GetEntry(ctx context.Context, id int64) (MasterDictionaryEntry, error) {
	if id <= 0 {
		return MasterDictionaryEntry{}, fmt.Errorf(masterDictionaryIDValidationError, ErrMasterDictionaryValidation)
	}

	entry, err := service.repository.GetByID(ctx, id)
	if err != nil {
		return MasterDictionaryEntry{}, fmt.Errorf("get master dictionary entry: %w", err)
	}
	return fromRepositoryEntry(entry), nil
}

// CreateEntry inserts a dictionary entry.
func (service *MasterDictionaryService) CreateEntry(ctx context.Context, input MasterDictionaryMutationInput) (MasterDictionaryEntry, error) {
	draft, err := service.validateMutationInput(input)
	if err != nil {
		return MasterDictionaryEntry{}, err
	}

	created, err := service.repository.Create(ctx, draft)
	if err != nil {
		return MasterDictionaryEntry{}, fmt.Errorf("create master dictionary entry: %w", err)
	}
	return fromRepositoryEntry(created), nil
}

// UpdateEntry updates a dictionary entry.
func (service *MasterDictionaryService) UpdateEntry(ctx context.Context, id int64, input MasterDictionaryMutationInput) (MasterDictionaryEntry, error) {
	if id <= 0 {
		return MasterDictionaryEntry{}, fmt.Errorf(masterDictionaryIDValidationError, ErrMasterDictionaryValidation)
	}

	draft, err := service.validateMutationInput(input)
	if err != nil {
		return MasterDictionaryEntry{}, err
	}

	updated, err := service.repository.Update(ctx, id, draft)
	if err != nil {
		return MasterDictionaryEntry{}, fmt.Errorf("update master dictionary entry: %w", err)
	}
	return fromRepositoryEntry(updated), nil
}

// DeleteEntry removes one dictionary entry.
func (service *MasterDictionaryService) DeleteEntry(ctx context.Context, id int64) error {
	if id <= 0 {
		return fmt.Errorf(masterDictionaryIDValidationError, ErrMasterDictionaryValidation)
	}

	if err := service.repository.Delete(ctx, id); err != nil {
		return fmt.Errorf("delete master dictionary entry: %w", err)
	}
	return nil
}

// ImportFromXML imports XML entries with REC allowlist filtering.
func (service *MasterDictionaryService) ImportFromXML(ctx context.Context, xmlPath string) (MasterDictionaryImportSummary, error) {
	resolvedPath, err := resolveMasterDictionaryXMLPath(xmlPath)
	if err != nil {
		return MasterDictionaryImportSummary{}, err
	}

	totalRecords, err := countXMLStringRecords(resolvedPath)
	if err != nil {
		return MasterDictionaryImportSummary{}, err
	}
	service.emitImportProgressIfEnabled(ctx, 0)

	file, err := os.Open(resolvedPath) // #nosec G304 -- resolvedPath is validated by resolveMasterDictionaryXMLPath.
	if err != nil {
		return MasterDictionaryImportSummary{}, fmt.Errorf(masterDictionaryOpenXMLFileError, err)
	}
	defer func() {
		_ = file.Close()
	}()

	decoder := xml.NewDecoder(file)
	counters := masterDictionaryImportCounters{}
	processedCount := 0
	for {
		record, done, readErr := readNextXMLStringRecord(decoder)
		if readErr != nil {
			return MasterDictionaryImportSummary{}, readErr
		}
		if done {
			break
		}
		if importErr := service.importXMLRecord(ctx, record, &counters); importErr != nil {
			return MasterDictionaryImportSummary{}, importErr
		}

		processedCount++
		service.emitImportProgressIfEnabled(ctx, normalizeImportProgress(processedCount, totalRecords))
	}

	service.emitImportProgressIfEnabled(ctx, 100)
	return MasterDictionaryImportSummary{
		FilePath:      resolvedPath,
		FileName:      filepath.Base(resolvedPath),
		ImportedCount: counters.importedCount,
		UpdatedCount:  counters.updatedCount,
		SkippedCount:  counters.skippedCount,
		SelectedREC:   allowedRECList(),
		LastEntryID:   counters.lastEntryID,
	}, nil
}

func resolveMasterDictionaryXMLPath(rawPath string) (string, error) {
	trimmedPath := strings.TrimSpace(rawPath)
	if trimmedPath == "" {
		return "", fmt.Errorf("%w: xml path is required", ErrMasterDictionaryValidation)
	}

	for _, candidate := range masterDictionaryXMLPathCandidates(trimmedPath) {
		info, statErr := os.Stat(candidate)
		if statErr != nil {
			if errors.Is(statErr, os.ErrNotExist) {
				continue
			}
			return "", fmt.Errorf("stat xml file: %w", statErr)
		}
		if info.IsDir() {
			continue
		}
		if !strings.EqualFold(filepath.Ext(candidate), ".xml") {
			return "", fmt.Errorf("%w: xml file extension is required", ErrMasterDictionaryValidation)
		}
		return candidate, nil
	}

	return "", fmt.Errorf(masterDictionaryOpenXMLFileError, os.ErrNotExist)
}

func masterDictionaryXMLPathCandidates(rawPath string) []string {
	cleanedPath := filepath.Clean(rawPath)
	baseName := filepath.Base(cleanedPath)

	candidates := make([]string, 0, 8)
	seen := make(map[string]struct{}, 8)
	appendCandidate := func(path string) {
		if strings.TrimSpace(path) == "" {
			return
		}
		cleaned := filepath.Clean(path)
		if _, exists := seen[cleaned]; exists {
			return
		}
		seen[cleaned] = struct{}{}
		candidates = append(candidates, cleaned)
	}

	appendCandidate(cleanedPath)
	appendCandidate(baseName)

	if cwd, err := os.Getwd(); err == nil {
		directory := cwd
		for depth := 0; depth < 6; depth++ {
			appendCandidate(filepath.Join(directory, "dictionaries", baseName))
			parent := filepath.Dir(directory)
			if parent == directory {
				break
			}
			directory = parent
		}
	}

	return candidates
}

func countXMLStringRecords(resolvedPath string) (int, error) {
	file, err := os.Open(resolvedPath) // #nosec G304 -- resolvedPath is validated by resolveMasterDictionaryXMLPath.
	if err != nil {
		return 0, fmt.Errorf(masterDictionaryOpenXMLFileError, err)
	}
	defer func() {
		_ = file.Close()
	}()

	decoder := xml.NewDecoder(file)
	total := 0
	for {
		token, tokenErr := decoder.Token()
		if tokenErr != nil {
			if errors.Is(tokenErr, io.EOF) {
				break
			}
			return 0, fmt.Errorf("read xml token: %w", tokenErr)
		}

		startElement, ok := token.(xml.StartElement)
		if ok && startElement.Name.Local == "String" {
			total++
		}
	}
	return total, nil
}

func normalizeImportProgress(processedCount, totalCount int) int {
	if totalCount <= 0 {
		return 100
	}
	if processedCount <= 0 {
		return 0
	}
	progress := int(float64(processedCount*100) / float64(totalCount))
	if progress < 0 {
		return 0
	}
	if progress > 100 {
		return 100
	}
	return progress
}

func (service *MasterDictionaryService) emitImportProgressIfEnabled(ctx context.Context, progress int) {
	if service.emitImportProgress == nil {
		return
	}
	service.emitImportProgress(ctx, progress)
}

// IsNotFoundError reports whether the error means entry not found.
func IsNotFoundError(err error) bool {
	return errors.Is(err, repository.ErrMasterDictionaryEntryNotFound)
}

func (service *MasterDictionaryService) validateMutationInput(input MasterDictionaryMutationInput) (repository.MasterDictionaryDraft, error) {
	source := strings.TrimSpace(input.Source)
	if source == "" {
		return repository.MasterDictionaryDraft{}, fmt.Errorf("%w: source is required", ErrMasterDictionaryValidation)
	}

	translation := strings.TrimSpace(input.Translation)
	if translation == "" {
		return repository.MasterDictionaryDraft{}, fmt.Errorf("%w: translation is required", ErrMasterDictionaryValidation)
	}

	category := strings.TrimSpace(input.Category)
	if category == "" {
		category = masterDictionaryDefaultCategory
	}

	origin := strings.TrimSpace(input.Origin)
	if origin == "" {
		origin = masterDictionaryDefaultOrigin
	}

	return repository.MasterDictionaryDraft{
		Source:      source,
		Translation: translation,
		Category:    category,
		Origin:      origin,
		REC:         strings.TrimSpace(input.REC),
		EDID:        strings.TrimSpace(input.EDID),
		UpdatedAt:   service.now().UTC(),
	}, nil
}

func (service *MasterDictionaryService) importXMLRecord(ctx context.Context, record xmlStringRecord, counters *masterDictionaryImportCounters) error {
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

	entry, created, err := service.repository.UpsertBySourceAndREC(ctx, repository.MasterDictionaryImportRecord{
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

func readNextXMLStringRecord(decoder *xml.Decoder) (xmlStringRecord, bool, error) {
	for {
		token, err := decoder.Token()
		if err != nil {
			if errors.Is(err, io.EOF) {
				return xmlStringRecord{}, true, nil
			}
			return xmlStringRecord{}, false, fmt.Errorf("read xml token: %w", err)
		}

		start, ok := token.(xml.StartElement)
		if !ok || start.Name.Local != "String" {
			continue
		}

		record := xmlStringRecord{}
		if err := decoder.DecodeElement(&record, &start); err != nil {
			return xmlStringRecord{}, false, fmt.Errorf("decode xml string record: %w", err)
		}
		return record, false, nil
	}
}

func isAllowedImportREC(rec string) bool {
	_, ok := allowedImportREC[rec]
	return ok
}

func fromRepositoryEntry(entry repository.MasterDictionaryEntry) MasterDictionaryEntry {
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

func categoryFromREC(rec string) string {
	switch strings.TrimSpace(rec) {
	case "BOOK:FULL":
		return "書籍"
	case "NPC_:FULL", "NPC_:SHRT", "RACE:FULL":
		return "NPC"
	case "ARMO:FULL", "WEAP:FULL":
		return "装備"
	case "LCTN:FULL", "CELL:FULL", "DOOR:FULL":
		return "地名"
	case "CONT:FULL", "MISC:FULL", "INGR:FULL", "FLOR:FULL", "ALCH:FULL":
		return "アイテム"
	case "FURN:FULL":
		return "設備"
	case "SHOU:FULL":
		return "シャウト"
	default:
		return "その他"
	}
}

func allowedRECList() []string {
	items := make([]string, 0, len(allowedImportREC))
	for rec := range allowedImportREC {
		items = append(items, rec)
	}
	sort.Strings(items)
	return items
}

type masterDictionaryImportCounters struct {
	importedCount int
	updatedCount  int
	skippedCount  int
	lastEntryID   int64
}

type xmlStringRecord struct {
	EDID   string `xml:"EDID"`
	REC    string `xml:"REC"`
	Source string `xml:"Source"`
	Dest   string `xml:"Dest"`
}
