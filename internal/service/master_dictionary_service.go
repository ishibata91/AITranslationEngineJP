package service

import (
	"context"
	"errors"
	"fmt"
	"io"
	"sort"
	"strings"
	"time"
)

var (
	// ErrMasterDictionaryValidation means the request payload is invalid.
	ErrMasterDictionaryValidation = errors.New("master dictionary validation error")
	// ErrMasterDictionaryEntryNotFound means the requested dictionary entry does not exist.
	ErrMasterDictionaryEntryNotFound = errors.New("master dictionary entry not found")
)

const (
	masterDictionaryDefaultCategory   = "固有名詞"
	masterDictionaryDefaultOrigin     = "手動登録"
	masterDictionaryImportOrigin      = "XML取込"
	masterDictionaryIDValidationError = "%w: id must be greater than zero"
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

// MasterDictionaryDraft is a validated payload for persistence.
type MasterDictionaryDraft struct {
	Source      string
	Translation string
	Category    string
	Origin      string
	REC         string
	EDID        string
	UpdatedAt   time.Time
}

// MasterDictionaryImportRecord describes one XML-derived dictionary record.
type MasterDictionaryImportRecord struct {
	Source      string
	Translation string
	REC         string
	EDID        string
	Category    string
	Origin      string
	UpdatedAt   time.Time
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

// RepositoryPort defines persistence operations for master dictionary CRUD.
type RepositoryPort interface {
	List(ctx context.Context, query MasterDictionaryQuery) (MasterDictionaryListResult, error)
	GetByID(ctx context.Context, id int64) (MasterDictionaryEntry, error)
	Create(ctx context.Context, draft MasterDictionaryDraft) (MasterDictionaryEntry, error)
	Update(ctx context.Context, id int64, draft MasterDictionaryDraft) (MasterDictionaryEntry, error)
	Delete(ctx context.Context, id int64) error
	UpsertBySourceAndREC(ctx context.Context, record MasterDictionaryImportRecord) (MasterDictionaryEntry, bool, error)
}

// XMLFilePort resolves and opens XML files for import.
type XMLFilePort interface {
	ResolvePath(rawPath string) (string, error)
	Open(path string) (io.ReadCloser, error)
}

// XMLRecordReaderPort counts and streams XML string records.
type XMLRecordReaderPort interface {
	CountStringRecords(reader io.Reader) (int, error)
	ReadStringRecords(reader io.Reader, handle func(xmlStringRecord) error) error
}

// RuntimeContextPort allows import service to emit runtime progress without knowing Wails.
type RuntimeContextPort interface {
	EmitImportProgress(ctx context.Context, progress int)
}

// IsNotFoundError reports whether the error means entry not found.
func IsNotFoundError(err error) bool {
	return errors.Is(err, ErrMasterDictionaryEntryNotFound)
}

func normalizeClock(now func() time.Time) func() time.Time {
	if now != nil {
		return now
	}
	return time.Now
}

func validateMasterDictionaryID(id int64) error {
	if id <= 0 {
		return fmt.Errorf(masterDictionaryIDValidationError, ErrMasterDictionaryValidation)
	}
	return nil
}

func validateMutationInput(input MasterDictionaryMutationInput, now func() time.Time) (MasterDictionaryDraft, error) {
	source := strings.TrimSpace(input.Source)
	if source == "" {
		return MasterDictionaryDraft{}, fmt.Errorf("%w: source is required", ErrMasterDictionaryValidation)
	}

	translation := strings.TrimSpace(input.Translation)
	if translation == "" {
		return MasterDictionaryDraft{}, fmt.Errorf("%w: translation is required", ErrMasterDictionaryValidation)
	}

	category := strings.TrimSpace(input.Category)
	if category == "" {
		category = masterDictionaryDefaultCategory
	}

	origin := strings.TrimSpace(input.Origin)
	if origin == "" {
		origin = masterDictionaryDefaultOrigin
	}

	return MasterDictionaryDraft{
		Source:      source,
		Translation: translation,
		Category:    category,
		Origin:      origin,
		REC:         strings.TrimSpace(input.REC),
		EDID:        strings.TrimSpace(input.EDID),
		UpdatedAt:   now().UTC(),
	}, nil
}

func isAllowedImportREC(rec string) bool {
	_, ok := allowedImportREC[rec]
	return ok
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

type masterDictionaryImportCounters struct {
	importedCount     int
	updatedCount      int
	skippedCount      int
	lastEntryID       int64
	importedEntryKeys map[string]struct{}
	updatedEntryKeys  map[string]struct{}
}

func (counters *masterDictionaryImportCounters) trackImportedEntry(
	entry MasterDictionaryEntry,
	record xmlStringRecord,
) {
	key := importEntryKey(entry.ID, record.Source, record.REC)
	if counters.importedEntryKeys == nil {
		counters.importedEntryKeys = map[string]struct{}{}
	}
	if _, exists := counters.importedEntryKeys[key]; exists {
		return
	}
	counters.importedEntryKeys[key] = struct{}{}
	counters.importedCount++
}

func (counters *masterDictionaryImportCounters) trackUpdatedEntry(
	entry MasterDictionaryEntry,
	record xmlStringRecord,
) {
	key := importEntryKey(entry.ID, record.Source, record.REC)
	if _, exists := counters.importedEntryKeys[key]; exists {
		return
	}
	if _, exists := counters.updatedEntryKeys[key]; exists {
		return
	}
	if counters.updatedEntryKeys == nil {
		counters.updatedEntryKeys = map[string]struct{}{}
	}
	counters.updatedEntryKeys[key] = struct{}{}
	counters.updatedCount++
}

func importEntryKey(entryID int64, source string, rec string) string {
	if entryID > 0 {
		return fmt.Sprintf("id:%d", entryID)
	}
	return strings.ToLower(strings.TrimSpace(source)) + "\x00" + strings.ToLower(strings.TrimSpace(rec))
}

type xmlStringRecord struct {
	EDID   string `xml:"EDID"`
	REC    string `xml:"REC"`
	Source string `xml:"Source"`
	Dest   string `xml:"Dest"`
}
