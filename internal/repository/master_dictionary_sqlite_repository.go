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
FROM DICTIONARY_ENTRY
WHERE dictionary_lifecycle = 'master'
AND (
  ? = ''
  OR lower(source_term || ' ' || translated_term || ' ' || CAST(id AS TEXT)) LIKE ?
)
AND (
  ? = ''
  OR ? = 'すべて'
  OR term_kind = ?
);`
	listMasterDictionaryEntriesSQL = `
SELECT id, source_term, translated_term, term_kind, dictionary_source, dictionary_scope, updated_at
FROM DICTIONARY_ENTRY
WHERE dictionary_lifecycle = 'master'
AND (
  ? = ''
  OR lower(source_term || ' ' || translated_term || ' ' || CAST(id AS TEXT)) LIKE ?
)
AND (
  ? = ''
  OR ? = 'すべて'
  OR term_kind = ?
)
ORDER BY updated_at DESC, id DESC
LIMIT ?
OFFSET ?;`
	getMasterDictionaryEntrySQL = `
SELECT id, source_term, translated_term, term_kind, dictionary_source, dictionary_scope, updated_at
FROM DICTIONARY_ENTRY
WHERE dictionary_lifecycle = 'master'
  AND id = ?
LIMIT 1;`
	findMasterDictionaryEntryBySourceAndRECSQL = `
SELECT id, source_term, translated_term, term_kind, dictionary_source, dictionary_scope, updated_at
FROM DICTIONARY_ENTRY
WHERE dictionary_lifecycle = 'master'
  AND lower(trim(source_term)) = lower(trim(?))
  AND lower(trim(dictionary_scope)) = lower(trim(?))
LIMIT 1;`
	createMasterDictionaryEntrySQL = `
INSERT INTO DICTIONARY_ENTRY (
  xtranslator_translation_xml_id,
  dictionary_lifecycle,
  dictionary_scope,
  dictionary_source,
  source_term,
  translated_term,
  term_kind,
  reusable,
  created_at,
  updated_at
) VALUES (
  :xtranslator_translation_xml_id,
  'master',
  :dictionary_scope,
  :dictionary_source,
  :source_term,
  :translated_term,
  :term_kind,
  1,
  :created_at,
  :updated_at
);`
	updateMasterDictionaryEntrySQL = `
UPDATE DICTIONARY_ENTRY
SET source_term = :source_term,
    translated_term = :translated_term,
    term_kind = :term_kind,
    dictionary_source = :dictionary_source,
    dictionary_scope = :dictionary_scope,
    updated_at = :updated_at
WHERE id = :id
  AND dictionary_lifecycle = 'master';`
	updateImportedMasterDictionaryEntrySQL = `
UPDATE DICTIONARY_ENTRY
SET xtranslator_translation_xml_id = :xtranslator_translation_xml_id,
    translated_term = :translated_term,
    term_kind = :term_kind,
    dictionary_source = :dictionary_source,
    updated_at = :updated_at
WHERE id = :id
  AND dictionary_lifecycle = 'master';`
	deleteMasterDictionaryEntrySQL = `
DELETE FROM DICTIONARY_ENTRY
WHERE id = ?
  AND dictionary_lifecycle = 'master';`
)

// SQLiteMasterDictionaryRepository persists master dictionary entries into SQLite.
type SQLiteMasterDictionaryRepository struct {
	database *sqlx.DB
}

type sqliteMasterDictionaryRow struct {
	ID          int64  `db:"id"`
	Source      string `db:"source_term"`
	Translation string `db:"translated_term"`
	Category    string `db:"term_kind"`
	Origin      string `db:"dictionary_source"`
	REC         string `db:"dictionary_scope"`
	UpdatedAt   string `db:"updated_at"`
}

type sqliteMasterDictionaryMutationParams struct {
	ID                          int64  `db:"id"`
	Source                      string `db:"source_term"`
	Translation                 string `db:"translated_term"`
	Category                    string `db:"term_kind"`
	Origin                      string `db:"dictionary_source"`
	REC                         string `db:"dictionary_scope"`
	UpdatedAt                   string `db:"updated_at"`
	CreatedAt                   string `db:"created_at"`
	XTranslatorTranslationXMLID *int64 `db:"xtranslator_translation_xml_id"`
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
			ID:                          row.ID,
			Translation:                 record.Translation,
			Category:                    record.Category,
			Origin:                      record.Origin,
			UpdatedAt:                   record.UpdatedAt.UTC().Format(masterDictionaryTimestampLayout),
			XTranslatorTranslationXMLID: record.XTranslatorTranslationXMLID,
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
		Source:                      strings.TrimSpace(record.Source),
		Translation:                 record.Translation,
		Category:                    record.Category,
		Origin:                      record.Origin,
		REC:                         strings.TrimSpace(record.REC),
		EDID:                        record.EDID,
		UpdatedAt:                   record.UpdatedAt,
		XTranslatorTranslationXMLID: record.XTranslatorTranslationXMLID,
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
		EDID:        "",
		UpdatedAt:   updatedAt.UTC(),
	}, nil
}

func mutationParamsFromDraft(id int64, draft MasterDictionaryDraft) sqliteMasterDictionaryMutationParams {
	ts := draft.UpdatedAt.UTC().Format(masterDictionaryTimestampLayout)
	return sqliteMasterDictionaryMutationParams{
		ID:                          id,
		Source:                      draft.Source,
		Translation:                 draft.Translation,
		Category:                    draft.Category,
		Origin:                      draft.Origin,
		REC:                         draft.REC,
		UpdatedAt:                   ts,
		CreatedAt:                   ts,
		XTranslatorTranslationXMLID: draft.XTranslatorTranslationXMLID,
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
