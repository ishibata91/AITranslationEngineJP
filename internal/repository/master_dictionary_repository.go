package repository

import (
	"context"
	"errors"
	"fmt"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"
)

var (
	// ErrMasterDictionaryEntryNotFound means the requested dictionary entry does not exist.
	ErrMasterDictionaryEntryNotFound = errors.New("master dictionary entry not found")
)

const (
	masterDictionaryRecLocationFull = "LCTN:FULL"
	masterDictionaryErrIDFormat     = "%w: id=%d"
)

// MasterDictionaryEntry represents one master dictionary record.
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

// MasterDictionaryListQuery describes list and search conditions.
type MasterDictionaryListQuery struct {
	SearchTerm string
	Category   string
	Page       int
	PageSize   int
}

// MasterDictionaryListResult is a paged query result.
type MasterDictionaryListResult struct {
	Items      []MasterDictionaryEntry
	TotalCount int
	Page       int
	PageSize   int
}

// MasterDictionaryDraft is an input payload for create or update.
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

// MasterDictionaryRepository defines persistence operations for master dictionary CRUD.
type MasterDictionaryRepository interface {
	List(ctx context.Context, query MasterDictionaryListQuery) (MasterDictionaryListResult, error)
	GetByID(ctx context.Context, id int64) (MasterDictionaryEntry, error)
	Create(ctx context.Context, draft MasterDictionaryDraft) (MasterDictionaryEntry, error)
	Update(ctx context.Context, id int64, draft MasterDictionaryDraft) (MasterDictionaryEntry, error)
	Delete(ctx context.Context, id int64) error
	UpsertBySourceAndREC(ctx context.Context, record MasterDictionaryImportRecord) (MasterDictionaryEntry, bool, error)
}

// InMemoryMasterDictionaryRepository is an in-memory implementation for the current phase.
type InMemoryMasterDictionaryRepository struct {
	mutex   sync.RWMutex
	entries []MasterDictionaryEntry
	nextID  int64
}

// NewInMemoryMasterDictionaryRepository creates a repository seeded with initial records.
func NewInMemoryMasterDictionaryRepository(seed []MasterDictionaryEntry) *InMemoryMasterDictionaryRepository {
	entries := make([]MasterDictionaryEntry, len(seed))
	copy(entries, seed)
	maxID := int64(0)
	for _, entry := range entries {
		if entry.ID > maxID {
			maxID = entry.ID
		}
	}

	return &InMemoryMasterDictionaryRepository{
		entries: entries,
		nextID:  maxID + 1,
	}
}

// DefaultMasterDictionarySeed returns deterministic initial records used by the backend shell.
func DefaultMasterDictionarySeed(now time.Time) []MasterDictionaryEntry {
	base := []MasterDictionaryDraft{
		{Source: "Whiterun", Translation: "ホワイトラン", Category: "地名", Origin: "初期データ", REC: masterDictionaryRecLocationFull, EDID: "LocWhiterun"},
		{Source: "Riften", Translation: "リフテン", Category: "地名", Origin: "初期データ", REC: masterDictionaryRecLocationFull, EDID: "LocRiften"},
		{Source: "Windhelm", Translation: "ウィンドヘルム", Category: "地名", Origin: "初期データ", REC: masterDictionaryRecLocationFull, EDID: "LocWindhelm"},
		{Source: "Solitude", Translation: "ソリチュード", Category: "地名", Origin: "初期データ", REC: masterDictionaryRecLocationFull, EDID: "LocSolitude"},
		{Source: "Markarth", Translation: "マルカルス", Category: "地名", Origin: "初期データ", REC: masterDictionaryRecLocationFull, EDID: "LocMarkarth"},
		{Source: "Morthal", Translation: "モーサル", Category: "地名", Origin: "初期データ", REC: masterDictionaryRecLocationFull, EDID: "LocMorthal"},
		{Source: "Dawnstar", Translation: "ドーンスター", Category: "地名", Origin: "初期データ", REC: masterDictionaryRecLocationFull, EDID: "LocDawnstar"},
		{Source: "Falkreath", Translation: "ファルクリース", Category: "地名", Origin: "初期データ", REC: masterDictionaryRecLocationFull, EDID: "LocFalkreath"},
		{Source: "Winterhold", Translation: "ウィンターホールド", Category: "地名", Origin: "初期データ", REC: masterDictionaryRecLocationFull, EDID: "LocWinterhold"},
		{Source: "Riverwood", Translation: "リバーウッド", Category: "地名", Origin: "初期データ", REC: "CELL:FULL", EDID: "CellRiverwood"},
	}

	entries := make([]MasterDictionaryEntry, 0, 40)
	for index := 0; index < 40; index++ {
		seed := base[index%len(base)]
		entries = append(entries, MasterDictionaryEntry{
			ID:          int64(index + 1),
			Source:      fmt.Sprintf("%s %02d", seed.Source, index+1),
			Translation: fmt.Sprintf("%s %02d", seed.Translation, index+1),
			Category:    seed.Category,
			Origin:      seed.Origin,
			REC:         seed.REC,
			EDID:        fmt.Sprintf("%s_%02d", seed.EDID, index+1),
			UpdatedAt:   now.Add(-time.Duration(index) * time.Minute),
		})
	}
	return entries
}

// List returns filtered and paginated dictionary entries.
func (repository *InMemoryMasterDictionaryRepository) List(_ context.Context, query MasterDictionaryListQuery) (MasterDictionaryListResult, error) {
	repository.mutex.RLock()
	defer repository.mutex.RUnlock()

	filtered := repository.filter(query)
	total := len(filtered)
	page, pageSize := normalizePagination(query.Page, query.PageSize, total)

	start := (page - 1) * pageSize
	end := start + pageSize
	if start > total {
		start = total
	}
	if end > total {
		end = total
	}

	items := make([]MasterDictionaryEntry, end-start)
	copy(items, filtered[start:end])

	return MasterDictionaryListResult{
		Items:      items,
		TotalCount: total,
		Page:       page,
		PageSize:   pageSize,
	}, nil
}

// GetByID returns one dictionary entry by identifier.
func (repository *InMemoryMasterDictionaryRepository) GetByID(_ context.Context, id int64) (MasterDictionaryEntry, error) {
	repository.mutex.RLock()
	defer repository.mutex.RUnlock()

	entry, ok := repository.findByID(id)
	if !ok {
		return MasterDictionaryEntry{}, fmt.Errorf(masterDictionaryErrIDFormat, ErrMasterDictionaryEntryNotFound, id)
	}
	return entry, nil
}

// Create inserts a dictionary entry.
func (repository *InMemoryMasterDictionaryRepository) Create(_ context.Context, draft MasterDictionaryDraft) (MasterDictionaryEntry, error) {
	repository.mutex.Lock()
	defer repository.mutex.Unlock()

	entry := MasterDictionaryEntry{
		ID:          repository.nextID,
		Source:      draft.Source,
		Translation: draft.Translation,
		Category:    draft.Category,
		Origin:      draft.Origin,
		REC:         draft.REC,
		EDID:        draft.EDID,
		UpdatedAt:   draft.UpdatedAt,
	}
	repository.nextID++
	repository.entries = append(repository.entries, entry)

	return entry, nil
}

// Update changes an existing dictionary entry.
func (repository *InMemoryMasterDictionaryRepository) Update(_ context.Context, id int64, draft MasterDictionaryDraft) (MasterDictionaryEntry, error) {
	repository.mutex.Lock()
	defer repository.mutex.Unlock()

	entryIndex := repository.findIndexByID(id)
	if entryIndex == -1 {
		return MasterDictionaryEntry{}, fmt.Errorf(masterDictionaryErrIDFormat, ErrMasterDictionaryEntryNotFound, id)
	}

	repository.entries[entryIndex].Source = draft.Source
	repository.entries[entryIndex].Translation = draft.Translation
	repository.entries[entryIndex].Category = draft.Category
	repository.entries[entryIndex].Origin = draft.Origin
	repository.entries[entryIndex].REC = draft.REC
	repository.entries[entryIndex].EDID = draft.EDID
	repository.entries[entryIndex].UpdatedAt = draft.UpdatedAt

	return repository.entries[entryIndex], nil
}

// Delete removes an existing dictionary entry.
func (repository *InMemoryMasterDictionaryRepository) Delete(_ context.Context, id int64) error {
	repository.mutex.Lock()
	defer repository.mutex.Unlock()

	entryIndex := repository.findIndexByID(id)
	if entryIndex == -1 {
		return fmt.Errorf(masterDictionaryErrIDFormat, ErrMasterDictionaryEntryNotFound, id)
	}

	repository.entries = append(repository.entries[:entryIndex], repository.entries[entryIndex+1:]...)
	return nil
}

// UpsertBySourceAndREC creates or updates an XML-derived record identified by source + REC.
func (repository *InMemoryMasterDictionaryRepository) UpsertBySourceAndREC(_ context.Context, record MasterDictionaryImportRecord) (MasterDictionaryEntry, bool, error) {
	repository.mutex.Lock()
	defer repository.mutex.Unlock()

	for index := range repository.entries {
		entry := repository.entries[index]
		if !strings.EqualFold(strings.TrimSpace(entry.Source), strings.TrimSpace(record.Source)) {
			continue
		}
		if !strings.EqualFold(strings.TrimSpace(entry.REC), strings.TrimSpace(record.REC)) {
			continue
		}

		repository.entries[index].Translation = record.Translation
		repository.entries[index].Category = record.Category
		repository.entries[index].Origin = record.Origin
		repository.entries[index].EDID = record.EDID
		repository.entries[index].UpdatedAt = record.UpdatedAt
		return repository.entries[index], false, nil
	}

	entry := MasterDictionaryEntry{
		ID:          repository.nextID,
		Source:      record.Source,
		Translation: record.Translation,
		Category:    record.Category,
		Origin:      record.Origin,
		REC:         record.REC,
		EDID:        record.EDID,
		UpdatedAt:   record.UpdatedAt,
	}
	repository.nextID++
	repository.entries = append(repository.entries, entry)
	return entry, true, nil
}

func (repository *InMemoryMasterDictionaryRepository) findByID(id int64) (MasterDictionaryEntry, bool) {
	for _, entry := range repository.entries {
		if entry.ID == id {
			return entry, true
		}
	}
	return MasterDictionaryEntry{}, false
}

func (repository *InMemoryMasterDictionaryRepository) findIndexByID(id int64) int {
	for index, entry := range repository.entries {
		if entry.ID == id {
			return index
		}
	}
	return -1
}

func (repository *InMemoryMasterDictionaryRepository) filter(query MasterDictionaryListQuery) []MasterDictionaryEntry {
	items := make([]MasterDictionaryEntry, 0, len(repository.entries))
	needle := strings.ToLower(strings.TrimSpace(query.SearchTerm))
	category := strings.TrimSpace(query.Category)

	for _, entry := range repository.entries {
		if category != "" && category != "すべて" && entry.Category != category {
			continue
		}

		if needle != "" {
			haystack := strings.ToLower(entry.Source + " " + entry.Translation + " " + entry.EDID + " " + strconv.FormatInt(entry.ID, 10))
			if !strings.Contains(haystack, needle) {
				continue
			}
		}

		items = append(items, entry)
	}

	sort.SliceStable(items, func(left, right int) bool {
		if items[left].UpdatedAt.Equal(items[right].UpdatedAt) {
			return items[left].ID > items[right].ID
		}
		return items[left].UpdatedAt.After(items[right].UpdatedAt)
	})
	return items
}

func normalizePagination(page, pageSize, total int) (int, int) {
	normalizedPageSize := pageSize
	if normalizedPageSize <= 0 {
		normalizedPageSize = 30
	}
	if normalizedPageSize > 100 {
		normalizedPageSize = 100
	}

	normalizedPage := page
	if normalizedPage <= 0 {
		normalizedPage = 1
	}

	maxPage := 1
	if total > 0 {
		maxPage = ((total - 1) / normalizedPageSize) + 1
	}
	if normalizedPage > maxPage {
		normalizedPage = maxPage
	}

	return normalizedPage, normalizedPageSize
}
