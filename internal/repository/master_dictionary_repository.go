package repository

import (
	"context"
	"errors"
	"fmt"
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
	Source                      string
	Translation                 string
	Category                    string
	Origin                      string
	REC                         string
	EDID                        string
	UpdatedAt                   time.Time
	XTranslatorTranslationXMLID *int64
}

// MasterDictionaryImportRecord describes one XML-derived dictionary record.
type MasterDictionaryImportRecord struct {
	Source                      string
	Translation                 string
	REC                         string
	EDID                        string
	Category                    string
	Origin                      string
	UpdatedAt                   time.Time
	XTranslatorTranslationXMLID *int64
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
