package repository

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/jmoiron/sqlx"
)

// SQLiteFoundationDataRepository は FoundationDataRepository の SQLite 実装。
type SQLiteFoundationDataRepository struct {
	db *sqlx.DB
}

// NewSQLiteFoundationDataRepository は FoundationDataRepository を返す。
func NewSQLiteFoundationDataRepository(db *sqlx.DB) FoundationDataRepository {
	return &SQLiteFoundationDataRepository{db: db}
}

// ---------------------------------------------------------------------------
// 内部 row 型
// ---------------------------------------------------------------------------

type xTranslatorTranslationXMLRow struct {
	ID               int64  `db:"id"`
	FilePath         string `db:"file_path"`
	TargetPluginName string `db:"target_plugin_name"`
	TargetPluginType string `db:"target_plugin_type"`
	TermCount        int    `db:"term_count"`
	ImportedAt       string `db:"imported_at"`
}

func (r xTranslatorTranslationXMLRow) toModel() (XTranslatorTranslationXML, error) {
	importedAt, err := time.Parse(time.RFC3339, r.ImportedAt)
	if err != nil {
		return XTranslatorTranslationXML{}, wrapParseError("imported_at", err)
	}
	return XTranslatorTranslationXML{
		ID:               r.ID,
		FilePath:         r.FilePath,
		TargetPluginName: r.TargetPluginName,
		TargetPluginType: r.TargetPluginType,
		TermCount:        r.TermCount,
		ImportedAt:       importedAt,
	}, nil
}

type personaRow struct {
	ID                     int64  `db:"id"`
	NpcProfileID           int64  `db:"npc_profile_id"`
	TranslationJobID       *int64 `db:"translation_job_id"`
	PersonaLifecycle       string `db:"persona_lifecycle"`
	PersonaScope           string `db:"persona_scope"`
	PersonaSource          string `db:"persona_source"`
	PersonaDescription     string `db:"persona_description"`
	SpeechStyle            string `db:"speech_style"`
	PersonalitySummary     string `db:"personality_summary"`
	EvidenceUtteranceCount int    `db:"evidence_utterance_count"`
	CreatedAt              string `db:"created_at"`
	UpdatedAt              string `db:"updated_at"`
}

func (r personaRow) toModel() (Persona, error) {
	createdAt, err := time.Parse(time.RFC3339, r.CreatedAt)
	if err != nil {
		return Persona{}, wrapParseError("created_at", err)
	}
	updatedAt, err := time.Parse(time.RFC3339, r.UpdatedAt)
	if err != nil {
		return Persona{}, wrapParseError("updated_at", err)
	}
	return Persona{
		ID:                     r.ID,
		NpcProfileID:           r.NpcProfileID,
		TranslationJobID:       r.TranslationJobID,
		PersonaLifecycle:       r.PersonaLifecycle,
		PersonaScope:           r.PersonaScope,
		PersonaSource:          r.PersonaSource,
		PersonaDescription:     r.PersonaDescription,
		SpeechStyle:            r.SpeechStyle,
		PersonalitySummary:     r.PersonalitySummary,
		EvidenceUtteranceCount: r.EvidenceUtteranceCount,
		CreatedAt:              createdAt,
		UpdatedAt:              updatedAt,
	}, nil
}

type personaFieldEvidenceRow struct {
	ID                 int64  `db:"id"`
	PersonaID          int64  `db:"persona_id"`
	TranslationFieldID int64  `db:"translation_field_id"`
	EvidenceRole       string `db:"evidence_role"`
}

func (r personaFieldEvidenceRow) toModel() PersonaFieldEvidence {
	return PersonaFieldEvidence(r)
}

type dictionaryEntryRow struct {
	ID                          int64  `db:"id"`
	XTranslatorTranslationXMLID *int64 `db:"xtranslator_translation_xml_id"`
	TranslationJobID            *int64 `db:"translation_job_id"`
	DictionaryLifecycle         string `db:"dictionary_lifecycle"`
	DictionaryScope             string `db:"dictionary_scope"`
	DictionarySource            string `db:"dictionary_source"`
	SourceTerm                  string `db:"source_term"`
	TranslatedTerm              string `db:"translated_term"`
	TermKind                    string `db:"term_kind"`
	Reusable                    bool   `db:"reusable"`
	CreatedAt                   string `db:"created_at"`
	UpdatedAt                   string `db:"updated_at"`
}

func (r dictionaryEntryRow) toModel() (DictionaryEntry, error) {
	createdAt, err := time.Parse(time.RFC3339, r.CreatedAt)
	if err != nil {
		return DictionaryEntry{}, wrapParseError("created_at", err)
	}
	updatedAt, err := time.Parse(time.RFC3339, r.UpdatedAt)
	if err != nil {
		return DictionaryEntry{}, wrapParseError("updated_at", err)
	}
	return DictionaryEntry{
		ID:                          r.ID,
		XTranslatorTranslationXMLID: r.XTranslatorTranslationXMLID,
		TranslationJobID:            r.TranslationJobID,
		DictionaryLifecycle:         r.DictionaryLifecycle,
		DictionaryScope:             r.DictionaryScope,
		DictionarySource:            r.DictionarySource,
		SourceTerm:                  r.SourceTerm,
		TranslatedTerm:              r.TranslatedTerm,
		TermKind:                    r.TermKind,
		Reusable:                    r.Reusable,
		CreatedAt:                   createdAt,
		UpdatedAt:                   updatedAt,
	}, nil
}

// ---------------------------------------------------------------------------
// SQL 定数
// ---------------------------------------------------------------------------

const (
	insertXTranslatorTranslationXML = `
INSERT INTO XTRANSLATOR_TRANSLATION_XML
  (file_path, target_plugin_name, target_plugin_type, term_count, imported_at)
VALUES
  (:file_path, :target_plugin_name, :target_plugin_type, :term_count, :imported_at)`

	selectXTranslatorTranslationXMLByID = `
SELECT id, file_path, target_plugin_name, target_plugin_type, term_count, imported_at
FROM XTRANSLATOR_TRANSLATION_XML WHERE id = ?`

	insertPersona = `
INSERT INTO PERSONA
  (npc_profile_id, translation_job_id, persona_lifecycle, persona_scope, persona_source,
   persona_description, speech_style, personality_summary, evidence_utterance_count,
   created_at, updated_at)
VALUES
  (:npc_profile_id, :translation_job_id, :persona_lifecycle, :persona_scope, :persona_source,
   :persona_description, :speech_style, :personality_summary, :evidence_utterance_count,
   :created_at, :updated_at)`

	selectPersonaByID = `
SELECT id, npc_profile_id, translation_job_id, persona_lifecycle, persona_scope, persona_source,
       persona_description, speech_style, personality_summary, evidence_utterance_count,
       created_at, updated_at
FROM PERSONA WHERE id = ?`

	selectPersonaByNpcProfileID = `
SELECT id, npc_profile_id, translation_job_id, persona_lifecycle, persona_scope, persona_source,
       persona_description, speech_style, personality_summary, evidence_utterance_count,
       created_at, updated_at
FROM PERSONA WHERE npc_profile_id = ?`

	updatePersona = `
UPDATE PERSONA SET
  persona_lifecycle       = :persona_lifecycle,
  persona_scope           = :persona_scope,
  persona_source          = :persona_source,
  persona_description     = :persona_description,
  speech_style            = :speech_style,
  personality_summary     = :personality_summary,
  evidence_utterance_count = :evidence_utterance_count,
  updated_at              = :updated_at
WHERE id = :id`

	insertPersonaFieldEvidence = `
INSERT INTO PERSONA_FIELD_EVIDENCE
  (persona_id, translation_field_id, evidence_role)
VALUES
  (:persona_id, :translation_field_id, :evidence_role)`

	selectPersonaFieldEvidenceByPersonaID = `
SELECT id, persona_id, translation_field_id, evidence_role
FROM PERSONA_FIELD_EVIDENCE WHERE persona_id = ?`

	insertDictionaryEntry = `
INSERT INTO DICTIONARY_ENTRY
  (xtranslator_translation_xml_id, translation_job_id, dictionary_lifecycle, dictionary_scope,
   dictionary_source, source_term, translated_term, term_kind, reusable, created_at, updated_at)
VALUES
  (:xtranslator_translation_xml_id, :translation_job_id, :dictionary_lifecycle, :dictionary_scope,
   :dictionary_source, :source_term, :translated_term, :term_kind, :reusable, :created_at, :updated_at)`

	selectDictionaryEntryByID = `
SELECT id, xtranslator_translation_xml_id, translation_job_id, dictionary_lifecycle, dictionary_scope,
       dictionary_source, source_term, translated_term, term_kind, reusable, created_at, updated_at
FROM DICTIONARY_ENTRY WHERE id = ?`

	updateDictionaryEntry = `
UPDATE DICTIONARY_ENTRY SET
  dictionary_lifecycle = :dictionary_lifecycle,
  dictionary_scope     = :dictionary_scope,
  dictionary_source    = :dictionary_source,
  source_term          = :source_term,
  translated_term      = :translated_term,
  term_kind            = :term_kind,
  reusable             = :reusable,
  updated_at           = :updated_at
WHERE id = :id`

	deleteDictionaryEntry = `DELETE FROM DICTIONARY_ENTRY WHERE id = ?`
)

// ---------------------------------------------------------------------------
// エラーヘルパー
// ---------------------------------------------------------------------------

func wrapParseError(field string, err error) error {
	return fmt.Errorf("parse %s: %w", field, err)
}

func isFKConstraintError(err error) bool {
	return strings.Contains(err.Error(), "FOREIGN KEY constraint")
}

func mapFoundationSQLError(err error, label string) error {
	if isFKConstraintError(err) {
		return fmt.Errorf("%s: %w", label, ErrConflict)
	}
	return mapSQLError(err, label)
}

// ---------------------------------------------------------------------------
// XTranslatorTranslationXML
// ---------------------------------------------------------------------------

// CreateXTranslatorTranslationXML は XTranslatorTranslationXML レコードを作成する。
func (r *SQLiteFoundationDataRepository) CreateXTranslatorTranslationXML(
	ctx context.Context,
	draft XTranslatorTranslationXMLDraft,
) (XTranslatorTranslationXML, error) {
	ext := extractTx(ctx, r.db)
	row := xTranslatorTranslationXMLRow{
		FilePath:         draft.FilePath,
		TargetPluginName: draft.TargetPluginName,
		TargetPluginType: draft.TargetPluginType,
		TermCount:        draft.TermCount,
		ImportedAt:       draft.ImportedAt.UTC().Format(time.RFC3339),
	}
	q, args, err := sqlx.Named(insertXTranslatorTranslationXML, row)
	if err != nil {
		return XTranslatorTranslationXML{}, fmt.Errorf("create xtranslator_translation_xml named: %w", err)
	}
	result, err := ext.ExecContext(ctx, q, args...)
	if err != nil {
		return XTranslatorTranslationXML{}, mapFoundationSQLError(err, "create xtranslator_translation_xml")
	}
	id, err := result.LastInsertId()
	if err != nil {
		return XTranslatorTranslationXML{}, fmt.Errorf("create xtranslator_translation_xml last insert id: %w", err)
	}
	return r.GetXTranslatorTranslationXMLByID(ctx, id)
}

// GetXTranslatorTranslationXMLByID は ID で XTranslatorTranslationXML を取得する。
func (r *SQLiteFoundationDataRepository) GetXTranslatorTranslationXMLByID(
	ctx context.Context,
	id int64,
) (XTranslatorTranslationXML, error) {
	ext := extractTx(ctx, r.db)
	var row xTranslatorTranslationXMLRow
	if err := sqlx.GetContext(ctx, ext, &row, selectXTranslatorTranslationXMLByID, id); err != nil {
		return XTranslatorTranslationXML{}, mapSQLError(err, "get xtranslator_translation_xml by id")
	}
	return row.toModel()
}

// ---------------------------------------------------------------------------
// Persona
// ---------------------------------------------------------------------------

// CreatePersona は Persona レコードを作成する。
func (r *SQLiteFoundationDataRepository) CreatePersona(
	ctx context.Context,
	draft PersonaDraft,
) (Persona, error) {
	ext := extractTx(ctx, r.db)
	now := time.Now().UTC().Format(time.RFC3339)
	row := personaRow{
		NpcProfileID:           draft.NpcProfileID,
		TranslationJobID:       draft.TranslationJobID,
		PersonaLifecycle:       draft.PersonaLifecycle,
		PersonaScope:           draft.PersonaScope,
		PersonaSource:          draft.PersonaSource,
		PersonaDescription:     draft.PersonaDescription,
		SpeechStyle:            draft.SpeechStyle,
		PersonalitySummary:     draft.PersonalitySummary,
		EvidenceUtteranceCount: draft.EvidenceUtteranceCount,
		CreatedAt:              now,
		UpdatedAt:              now,
	}
	q, args, err := sqlx.Named(insertPersona, row)
	if err != nil {
		return Persona{}, fmt.Errorf("create persona named: %w", err)
	}
	result, err := ext.ExecContext(ctx, q, args...)
	if err != nil {
		return Persona{}, mapFoundationSQLError(err, "create persona")
	}
	id, err := result.LastInsertId()
	if err != nil {
		return Persona{}, fmt.Errorf("create persona last insert id: %w", err)
	}
	return r.GetPersonaByID(ctx, id)
}

// GetPersonaByID は ID で Persona を取得する。
func (r *SQLiteFoundationDataRepository) GetPersonaByID(
	ctx context.Context,
	id int64,
) (Persona, error) {
	ext := extractTx(ctx, r.db)
	var row personaRow
	if err := sqlx.GetContext(ctx, ext, &row, selectPersonaByID, id); err != nil {
		return Persona{}, mapSQLError(err, "get persona by id")
	}
	return row.toModel()
}

// GetPersonaByNpcProfileID は NpcProfileID で Persona を取得する。
func (r *SQLiteFoundationDataRepository) GetPersonaByNpcProfileID(
	ctx context.Context,
	npcProfileID int64,
) (Persona, error) {
	ext := extractTx(ctx, r.db)
	var row personaRow
	if err := sqlx.GetContext(ctx, ext, &row, selectPersonaByNpcProfileID, npcProfileID); err != nil {
		return Persona{}, mapSQLError(err, "get persona by npc_profile_id")
	}
	return row.toModel()
}

// UpdatePersona は Persona を更新する。
func (r *SQLiteFoundationDataRepository) UpdatePersona(
	ctx context.Context,
	id int64,
	draft PersonaUpdateDraft,
) (Persona, error) {
	ext := extractTx(ctx, r.db)
	now := time.Now().UTC().Format(time.RFC3339)
	args := map[string]interface{}{
		"id":                       id,
		"persona_lifecycle":        draft.PersonaLifecycle,
		"persona_scope":            draft.PersonaScope,
		"persona_source":           draft.PersonaSource,
		"persona_description":      draft.PersonaDescription,
		"speech_style":             draft.SpeechStyle,
		"personality_summary":      draft.PersonalitySummary,
		"evidence_utterance_count": draft.EvidenceUtteranceCount,
		"updated_at":               now,
	}
	q, qArgs, err := sqlx.Named(updatePersona, args)
	if err != nil {
		return Persona{}, fmt.Errorf("update persona named: %w", err)
	}
	if _, err := ext.ExecContext(ctx, q, qArgs...); err != nil {
		return Persona{}, mapFoundationSQLError(err, "update persona")
	}
	return r.GetPersonaByID(ctx, id)
}

// ---------------------------------------------------------------------------
// PersonaFieldEvidence
// ---------------------------------------------------------------------------

// CreatePersonaFieldEvidence は PersonaFieldEvidence レコードを作成する。
func (r *SQLiteFoundationDataRepository) CreatePersonaFieldEvidence(
	ctx context.Context,
	draft PersonaFieldEvidenceDraft,
) (PersonaFieldEvidence, error) {
	ext := extractTx(ctx, r.db)
	row := personaFieldEvidenceRow{
		PersonaID:          draft.PersonaID,
		TranslationFieldID: draft.TranslationFieldID,
		EvidenceRole:       draft.EvidenceRole,
	}
	q, args, err := sqlx.Named(insertPersonaFieldEvidence, row)
	if err != nil {
		return PersonaFieldEvidence{}, fmt.Errorf("create persona_field_evidence named: %w", err)
	}
	result, err := ext.ExecContext(ctx, q, args...)
	if err != nil {
		return PersonaFieldEvidence{}, mapFoundationSQLError(err, "create persona_field_evidence")
	}
	id, err := result.LastInsertId()
	if err != nil {
		return PersonaFieldEvidence{}, fmt.Errorf("create persona_field_evidence last insert id: %w", err)
	}
	return PersonaFieldEvidence{
		ID:                 id,
		PersonaID:          draft.PersonaID,
		TranslationFieldID: draft.TranslationFieldID,
		EvidenceRole:       draft.EvidenceRole,
	}, nil
}

// ListPersonaFieldEvidenceByPersonaID は PersonaID に紐づく PersonaFieldEvidence 一覧を返す。
func (r *SQLiteFoundationDataRepository) ListPersonaFieldEvidenceByPersonaID(
	ctx context.Context,
	personaID int64,
) ([]PersonaFieldEvidence, error) {
	ext := extractTx(ctx, r.db)
	var rows []personaFieldEvidenceRow
	if err := sqlx.SelectContext(ctx, ext, &rows, selectPersonaFieldEvidenceByPersonaID, personaID); err != nil {
		return nil, mapSQLError(err, "list persona_field_evidence by persona_id")
	}
	result := make([]PersonaFieldEvidence, len(rows))
	for i, row := range rows {
		result[i] = row.toModel()
	}
	return result, nil
}

// ---------------------------------------------------------------------------
// DictionaryEntry
// ---------------------------------------------------------------------------

// CreateDictionaryEntry は DictionaryEntry レコードを作成する。
func (r *SQLiteFoundationDataRepository) CreateDictionaryEntry(
	ctx context.Context,
	draft DictionaryEntryDraft,
) (DictionaryEntry, error) {
	ext := extractTx(ctx, r.db)
	now := time.Now().UTC().Format(time.RFC3339)
	row := dictionaryEntryRow{
		XTranslatorTranslationXMLID: draft.XTranslatorTranslationXMLID,
		TranslationJobID:            draft.TranslationJobID,
		DictionaryLifecycle:         draft.DictionaryLifecycle,
		DictionaryScope:             draft.DictionaryScope,
		DictionarySource:            draft.DictionarySource,
		SourceTerm:                  draft.SourceTerm,
		TranslatedTerm:              draft.TranslatedTerm,
		TermKind:                    draft.TermKind,
		Reusable:                    draft.Reusable,
		CreatedAt:                   now,
		UpdatedAt:                   now,
	}
	q, args, err := sqlx.Named(insertDictionaryEntry, row)
	if err != nil {
		return DictionaryEntry{}, fmt.Errorf("create dictionary_entry named: %w", err)
	}
	result, err := ext.ExecContext(ctx, q, args...)
	if err != nil {
		return DictionaryEntry{}, mapFoundationSQLError(err, "create dictionary_entry")
	}
	id, err := result.LastInsertId()
	if err != nil {
		return DictionaryEntry{}, fmt.Errorf("create dictionary_entry last insert id: %w", err)
	}
	return r.GetDictionaryEntryByID(ctx, id)
}

// GetDictionaryEntryByID は ID で DictionaryEntry を取得する。
func (r *SQLiteFoundationDataRepository) GetDictionaryEntryByID(
	ctx context.Context,
	id int64,
) (DictionaryEntry, error) {
	ext := extractTx(ctx, r.db)
	var row dictionaryEntryRow
	if err := sqlx.GetContext(ctx, ext, &row, selectDictionaryEntryByID, id); err != nil {
		return DictionaryEntry{}, mapSQLError(err, "get dictionary_entry by id")
	}
	return row.toModel()
}

// UpdateDictionaryEntry は DictionaryEntry を更新する。
func (r *SQLiteFoundationDataRepository) UpdateDictionaryEntry(
	ctx context.Context,
	id int64,
	draft DictionaryEntryUpdateDraft,
) (DictionaryEntry, error) {
	ext := extractTx(ctx, r.db)
	now := time.Now().UTC().Format(time.RFC3339)
	args := map[string]interface{}{
		"id":                   id,
		"dictionary_lifecycle": draft.DictionaryLifecycle,
		"dictionary_scope":     draft.DictionaryScope,
		"dictionary_source":    draft.DictionarySource,
		"source_term":          draft.SourceTerm,
		"translated_term":      draft.TranslatedTerm,
		"term_kind":            draft.TermKind,
		"reusable":             draft.Reusable,
		"updated_at":           now,
	}
	q, qArgs, err := sqlx.Named(updateDictionaryEntry, args)
	if err != nil {
		return DictionaryEntry{}, fmt.Errorf("update dictionary_entry named: %w", err)
	}
	if _, err := ext.ExecContext(ctx, q, qArgs...); err != nil {
		return DictionaryEntry{}, mapFoundationSQLError(err, "update dictionary_entry")
	}
	return r.GetDictionaryEntryByID(ctx, id)
}

// DeleteDictionaryEntry は DictionaryEntry を削除する。
func (r *SQLiteFoundationDataRepository) DeleteDictionaryEntry(
	ctx context.Context,
	id int64,
) error {
	ext := extractTx(ctx, r.db)
	if _, err := ext.ExecContext(ctx, deleteDictionaryEntry, id); err != nil {
		return mapFoundationSQLError(err, "delete dictionary_entry")
	}
	return nil
}
