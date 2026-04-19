package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"time"

	sqliteinfra "aitranslationenginejp/internal/infra/sqlite/dbinit"

	"github.com/jmoiron/sqlx"
)

const (
	masterDictionaryTimestampLayout = time.RFC3339Nano
	countMasterDictionaryEntriesSQL = `
SELECT COUNT(*)
FROM master_dictionary_entries
WHERE (
  ? = ''
  OR lower(source || ' ' || translation || ' ' || edid || ' ' || CAST(id AS TEXT)) LIKE ?
)
AND (
  ? = ''
  OR ? = 'すべて'
  OR category = ?
);`
	listMasterDictionaryEntriesSQL = `
SELECT id, source, translation, category, origin, rec, edid, updated_at
FROM master_dictionary_entries
WHERE (
  ? = ''
  OR lower(source || ' ' || translation || ' ' || edid || ' ' || CAST(id AS TEXT)) LIKE ?
)
AND (
  ? = ''
  OR ? = 'すべて'
  OR category = ?
)
ORDER BY updated_at DESC, id DESC
LIMIT ?
OFFSET ?;`
	getMasterDictionaryEntrySQL = `
SELECT id, source, translation, category, origin, rec, edid, updated_at
FROM master_dictionary_entries
WHERE id = ?
LIMIT 1;`
	findMasterDictionaryEntryBySourceAndRECSQL = `
SELECT id, source, translation, category, origin, rec, edid, updated_at
FROM master_dictionary_entries
WHERE lower(trim(source)) = lower(trim(?))
  AND lower(trim(rec)) = lower(trim(?))
LIMIT 1;`
	createMasterDictionaryEntrySQL = `
INSERT INTO master_dictionary_entries (
  source,
  translation,
  category,
  origin,
  rec,
  edid,
  updated_at
) VALUES (
  :source,
  :translation,
  :category,
  :origin,
  :rec,
  :edid,
  :updated_at
);`
	updateMasterDictionaryEntrySQL = `
UPDATE master_dictionary_entries
SET source = :source,
    translation = :translation,
    category = :category,
    origin = :origin,
    rec = :rec,
    edid = :edid,
    updated_at = :updated_at
WHERE id = :id;`
	updateImportedMasterDictionaryEntrySQL = `
UPDATE master_dictionary_entries
SET translation = :translation,
    category = :category,
    origin = :origin,
    edid = :edid,
    updated_at = :updated_at
WHERE id = :id;`
	deleteMasterDictionaryEntrySQL = `
DELETE FROM master_dictionary_entries
WHERE id = ?;`
)

// SQLiteMasterDictionaryRepository persists master dictionary entries into SQLite.
type SQLiteMasterDictionaryRepository struct {
	database *sqlx.DB
}

type sqliteMasterDictionaryRow struct {
	ID          int64  `db:"id"`
	Source      string `db:"source"`
	Translation string `db:"translation"`
	Category    string `db:"category"`
	Origin      string `db:"origin"`
	REC         string `db:"rec"`
	EDID        string `db:"edid"`
	UpdatedAt   string `db:"updated_at"`
}

type sqliteMasterDictionaryMutationParams struct {
	ID          int64  `db:"id"`
	Source      string `db:"source"`
	Translation string `db:"translation"`
	Category    string `db:"category"`
	Origin      string `db:"origin"`
	REC         string `db:"rec"`
	EDID        string `db:"edid"`
	UpdatedAt   string `db:"updated_at"`
}

// NewSQLiteMasterDictionaryRepository opens the SQLite database after startup initialization has run.
func NewSQLiteMasterDictionaryRepository(
	ctx context.Context,
	databasePath string,
	seed []MasterDictionaryEntry,
) (*SQLiteMasterDictionaryRepository, error) {
	database, err := sqliteinfra.OpenMasterDictionaryDatabase(ctx, databasePath, toSQLiteSeedEntries(seed))
	if err != nil {
		return nil, fmt.Errorf("open sqlite master dictionary database: %w", err)
	}
	return &SQLiteMasterDictionaryRepository{database: database}, nil
}

// Close releases the underlying SQLite database handle.
func (repository *SQLiteMasterDictionaryRepository) Close() error {
	if err := repository.database.Close(); err != nil {
		return fmt.Errorf("close sqlite master dictionary database: %w", err)
	}
	return nil
}

// List returns filtered and paginated dictionary entries.
func (repository *SQLiteMasterDictionaryRepository) List(
	ctx context.Context,
	query MasterDictionaryListQuery,
) (result MasterDictionaryListResult, err error) {
	searchTerm := strings.TrimSpace(query.SearchTerm)
	searchPattern := "%" + strings.ToLower(searchTerm) + "%"
	category := strings.TrimSpace(query.Category)

	var totalCount int
	if err := repository.database.GetContext(
		ctx,
		&totalCount,
		countMasterDictionaryEntriesSQL,
		searchTerm,
		searchPattern,
		category,
		category,
		category,
	); err != nil {
		return MasterDictionaryListResult{}, fmt.Errorf("count master dictionary entries: %w", err)
	}

	page, pageSize := normalizePagination(query.Page, query.PageSize, totalCount)
	rows := []sqliteMasterDictionaryRow{}
	if err := repository.database.SelectContext(
		ctx,
		&rows,
		listMasterDictionaryEntriesSQL,
		searchTerm,
		searchPattern,
		category,
		category,
		category,
		pageSize,
		(page-1)*pageSize,
	); err != nil {
		return MasterDictionaryListResult{}, fmt.Errorf("list master dictionary entries: %w", err)
	}

	items := make([]MasterDictionaryEntry, 0, len(rows))
	for _, row := range rows {
		entry, err := fromSQLiteRow(row)
		if err != nil {
			return MasterDictionaryListResult{}, err
		}
		items = append(items, entry)
	}

	return MasterDictionaryListResult{
		Items:      items,
		TotalCount: totalCount,
		Page:       page,
		PageSize:   pageSize,
	}, nil
}

// GetByID returns one dictionary entry by identifier.
func (repository *SQLiteMasterDictionaryRepository) GetByID(ctx context.Context, id int64) (MasterDictionaryEntry, error) {
	row := sqliteMasterDictionaryRow{}
	if err := repository.database.GetContext(ctx, &row, getMasterDictionaryEntrySQL, id); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return MasterDictionaryEntry{}, fmt.Errorf(masterDictionaryErrIDFormat, ErrMasterDictionaryEntryNotFound, id)
		}
		return MasterDictionaryEntry{}, fmt.Errorf("get master dictionary entry: %w", err)
	}
	return fromSQLiteRow(row)
}

// Create inserts a dictionary entry.
func (repository *SQLiteMasterDictionaryRepository) Create(ctx context.Context, draft MasterDictionaryDraft) (MasterDictionaryEntry, error) {
	result, err := repository.database.NamedExecContext(ctx, createMasterDictionaryEntrySQL, mutationParamsFromDraft(0, draft))
	if err != nil {
		return MasterDictionaryEntry{}, fmt.Errorf("create master dictionary entry: %w", err)
	}
	createdID, err := result.LastInsertId()
	if err != nil {
		return MasterDictionaryEntry{}, fmt.Errorf("read created master dictionary id: %w", err)
	}
	return repository.GetByID(ctx, createdID)
}

// Update changes an existing dictionary entry.
func (repository *SQLiteMasterDictionaryRepository) Update(
	ctx context.Context,
	id int64,
	draft MasterDictionaryDraft,
) (MasterDictionaryEntry, error) {
	result, err := repository.database.NamedExecContext(ctx, updateMasterDictionaryEntrySQL, mutationParamsFromDraft(id, draft))
	if err != nil {
		return MasterDictionaryEntry{}, fmt.Errorf("update master dictionary entry: %w", err)
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return MasterDictionaryEntry{}, fmt.Errorf("read updated master dictionary rows affected: %w", err)
	}
	if rowsAffected == 0 {
		return MasterDictionaryEntry{}, fmt.Errorf(masterDictionaryErrIDFormat, ErrMasterDictionaryEntryNotFound, id)
	}
	return repository.GetByID(ctx, id)
}

// Delete removes an existing dictionary entry.
func (repository *SQLiteMasterDictionaryRepository) Delete(ctx context.Context, id int64) error {
	result, err := repository.database.ExecContext(ctx, deleteMasterDictionaryEntrySQL, id)
	if err != nil {
		return fmt.Errorf("delete master dictionary entry: %w", err)
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("read deleted master dictionary rows affected: %w", err)
	}
	if rowsAffected == 0 {
		return fmt.Errorf(masterDictionaryErrIDFormat, ErrMasterDictionaryEntryNotFound, id)
	}
	return nil
}

// UpsertBySourceAndREC creates or updates an XML-derived record identified by source + REC.
func (repository *SQLiteMasterDictionaryRepository) UpsertBySourceAndREC(
	ctx context.Context,
	record MasterDictionaryImportRecord,
) (MasterDictionaryEntry, bool, error) {
	row := sqliteMasterDictionaryRow{}
	err := repository.database.GetContext(
		ctx,
		&row,
		findMasterDictionaryEntryBySourceAndRECSQL,
		record.Source,
		record.REC,
	)
	if err == nil {
		result, updateErr := repository.database.NamedExecContext(ctx, updateImportedMasterDictionaryEntrySQL, sqliteMasterDictionaryMutationParams{
			ID:          row.ID,
			Translation: record.Translation,
			Category:    record.Category,
			Origin:      record.Origin,
			EDID:        record.EDID,
			UpdatedAt:   record.UpdatedAt.UTC().Format(masterDictionaryTimestampLayout),
		})
		if updateErr != nil {
			return MasterDictionaryEntry{}, false, fmt.Errorf("update imported master dictionary entry: %w", updateErr)
		}
		rowsAffected, rowsErr := result.RowsAffected()
		if rowsErr != nil {
			return MasterDictionaryEntry{}, false, fmt.Errorf("read updated imported master dictionary rows affected: %w", rowsErr)
		}
		if rowsAffected == 0 {
			return MasterDictionaryEntry{}, false, fmt.Errorf(masterDictionaryErrIDFormat, ErrMasterDictionaryEntryNotFound, row.ID)
		}
		entry, getErr := repository.GetByID(ctx, row.ID)
		if getErr != nil {
			return MasterDictionaryEntry{}, false, getErr
		}
		return entry, false, nil
	}
	if !errors.Is(err, sql.ErrNoRows) {
		return MasterDictionaryEntry{}, false, fmt.Errorf("find master dictionary entry by source and REC: %w", err)
	}

	entry, createErr := repository.Create(ctx, MasterDictionaryDraft{
		Source:      strings.TrimSpace(record.Source),
		Translation: record.Translation,
		Category:    record.Category,
		Origin:      record.Origin,
		REC:         strings.TrimSpace(record.REC),
		EDID:        record.EDID,
		UpdatedAt:   record.UpdatedAt,
	})
	if createErr != nil {
		return MasterDictionaryEntry{}, true, createErr
	}
	return entry, true, nil
}

func fromSQLiteRow(row sqliteMasterDictionaryRow) (MasterDictionaryEntry, error) {
	updatedAt, err := time.Parse(masterDictionaryTimestampLayout, row.UpdatedAt)
	if err != nil {
		return MasterDictionaryEntry{}, fmt.Errorf("parse updated_at for master dictionary entry %d: %w", row.ID, err)
	}
	return MasterDictionaryEntry{
		ID:          row.ID,
		Source:      row.Source,
		Translation: row.Translation,
		Category:    row.Category,
		Origin:      row.Origin,
		REC:         row.REC,
		EDID:        row.EDID,
		UpdatedAt:   updatedAt.UTC(),
	}, nil
}

func mutationParamsFromDraft(id int64, draft MasterDictionaryDraft) sqliteMasterDictionaryMutationParams {
	return sqliteMasterDictionaryMutationParams{
		ID:          id,
		Source:      draft.Source,
		Translation: draft.Translation,
		Category:    draft.Category,
		Origin:      draft.Origin,
		REC:         draft.REC,
		EDID:        draft.EDID,
		UpdatedAt:   draft.UpdatedAt.UTC().Format(masterDictionaryTimestampLayout),
	}
}

func toSQLiteSeedEntries(seed []MasterDictionaryEntry) []sqliteinfra.MasterDictionarySeedEntry {
	items := make([]sqliteinfra.MasterDictionarySeedEntry, 0, len(seed))
	for _, entry := range seed {
		items = append(items, sqliteinfra.MasterDictionarySeedEntry{
			Source:      entry.Source,
			Translation: entry.Translation,
			Category:    entry.Category,
			Origin:      entry.Origin,
			REC:         entry.REC,
			EDID:        entry.EDID,
			UpdatedAt:   entry.UpdatedAt,
		})
	}
	return items
}
