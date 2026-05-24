package repository

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"
	"strings"

	"github.com/jmoiron/sqlx"
)

const (
	processingTargetPhaseTermTranslation      = "term_translation"
	processingTargetPhasePersonaGeneration    = "persona_generation"
	processingTargetPhaseBodyTranslation      = "body_translation"
	processingTargetPhaseTranslationComplete  = "translation_complete"
	processingTargetMetadataSourceText        = "原文"
	processingTargetMetadataTranslatedText    = "訳文"
	processingTargetMetadataRecordType        = "レコード種別"
	processingTargetMetadataSubrecordType     = "項目種別"
	processingTargetMetadataFormID            = "FormID"
	processingTargetMetadataEditorID          = "EditorID"
	processingTargetMetadataTermKind          = "用語種別"
	processingTargetMetadataRace              = "種族"
	processingTargetMetadataSex               = "性別"
	processingTargetMetadataClass             = "クラス"
	processingTargetMetadataVoiceType         = "声種別"
	processingTargetMetadataEvidenceCount     = "原文発話数"
	processingTargetMetadataOutputStatus      = "出力状態"
	processingTargetMetadataAppliedPersona    = "適用ペルソナ"
	processingTargetSearchPatternEscapedSlash = `\`
	processingTargetListLogMessage            = "processing target list repository failed"
	processingTargetListLogEvent              = "processing_target_list_repository_failed"
	processingTargetListLogWhere              = "backend.repository.processing_target"
	processingTargetListLogResultFailed       = "failed"
	processingTargetListLogJobIDFormat        = "job:%d"
)

type processingTargetRow struct {
	ID              string         `db:"id"`
	Name            string         `db:"name"`
	Detail          string         `db:"detail"`
	TitlePart1      string         `db:"title_part_1"`
	TitlePart2      string         `db:"title_part_2"`
	TitlePart3      string         `db:"title_part_3"`
	MetadataValue1  string         `db:"metadata_value_1"`
	MetadataValue2  string         `db:"metadata_value_2"`
	MetadataValue3  string         `db:"metadata_value_3"`
	MetadataValue4  string         `db:"metadata_value_4"`
	MetadataValue5  sql.NullString `db:"metadata_value_5"`
	MetadataValue6  sql.NullString `db:"metadata_value_6"`
	MetadataValue7  sql.NullString `db:"metadata_value_7"`
	MetadataValue8  sql.NullString `db:"metadata_value_8"`
	MetadataValue9  sql.NullString `db:"metadata_value_9"`
	MetadataValue10 sql.NullString `db:"metadata_value_10"`
}

type processingTargetSQLSpec struct {
	countSQL       string
	listSQL        string
	metadataLabels []string
}

// ListProcessingTargets returns one processing target page from canonical job-owned tables.
func (r *SQLiteJobLifecycleRepository) ListProcessingTargets(
	ctx context.Context,
	query ProcessingTargetListQuery,
) (ProcessingTargetListResult, error) {
	spec, err := processingTargetSQLSpecForPhase(query.Phase)
	if err != nil {
		slog.WarnContext(ctx, processingTargetListLogMessage,
			slog.String("event", processingTargetListLogEvent),
			slog.String("where", processingTargetListLogWhere),
			slog.String("result", processingTargetListLogResultFailed),
			slog.String("id", fmt.Sprintf(processingTargetListLogJobIDFormat, query.JobID)),
			slog.String("reason", "unsupported_phase"),
		)
		return ProcessingTargetListResult{}, err
	}
	searchQuery := strings.TrimSpace(query.SearchQuery)
	searchPattern := processingTargetSearchPattern(searchQuery)
	page := query.Page
	if page < 1 {
		page = 1
	}
	pageSize := query.PageSize
	if pageSize < 1 {
		pageSize = 50
	}
	offset := (page - 1) * pageSize

	ext := extractTx(ctx, r.db)
	var totalCount int
	if err := sqlx.GetContext(ctx, ext, &totalCount, spec.countSQL, query.JobID, searchQuery, searchPattern); err != nil {
		slog.WarnContext(ctx, processingTargetListLogMessage,
			slog.String("event", processingTargetListLogEvent),
			slog.String("where", processingTargetListLogWhere),
			slog.String("result", processingTargetListLogResultFailed),
			slog.String("id", fmt.Sprintf(processingTargetListLogJobIDFormat, query.JobID)),
			slog.String("reason", "count_failed"),
		)
		return ProcessingTargetListResult{}, fmt.Errorf("count processing targets: %w", err)
	}
	rows := make([]processingTargetRow, 0, pageSize)
	if err := sqlx.SelectContext(ctx, ext, &rows, spec.listSQL, query.JobID, searchQuery, searchPattern, pageSize, offset); err != nil {
		slog.WarnContext(ctx, processingTargetListLogMessage,
			slog.String("event", processingTargetListLogEvent),
			slog.String("where", processingTargetListLogWhere),
			slog.String("result", processingTargetListLogResultFailed),
			slog.String("id", fmt.Sprintf(processingTargetListLogJobIDFormat, query.JobID)),
			slog.String("reason", "list_failed"),
		)
		return ProcessingTargetListResult{}, fmt.Errorf("list processing targets: %w", err)
	}
	items := make([]ProcessingTargetListItem, 0, len(rows))
	for _, row := range rows {
		items = append(items, row.toProcessingTargetItem(spec.metadataLabels))
	}
	return ProcessingTargetListResult{
		Items:       items,
		Page:        page,
		PageSize:    pageSize,
		TotalCount:  totalCount,
		SearchQuery: searchQuery,
	}, nil
}

func (row processingTargetRow) toProcessingTargetItem(labels []string) ProcessingTargetListItem {
	values := []string{
		row.MetadataValue1,
		row.MetadataValue2,
		row.MetadataValue3,
		row.MetadataValue4,
		processingTargetNullStringValue(row.MetadataValue5),
		processingTargetNullStringValue(row.MetadataValue6),
		processingTargetNullStringValue(row.MetadataValue7),
		processingTargetNullStringValue(row.MetadataValue8),
		processingTargetNullStringValue(row.MetadataValue9),
		processingTargetNullStringValue(row.MetadataValue10),
	}
	metadata := make([]ProcessingTargetListItemMetadata, 0, len(labels))
	for index, label := range labels {
		if index >= len(values) {
			break
		}
		value := strings.TrimSpace(values[index])
		if value == "" {
			continue
		}
		metadata = append(metadata, ProcessingTargetListItemMetadata{Label: label, Value: value})
	}
	return ProcessingTargetListItem{
		ID:         row.ID,
		Name:       row.Name,
		Detail:     row.Detail,
		TitleParts: nonBlankStrings(row.TitlePart1, row.TitlePart2, row.TitlePart3),
		Metadata:   metadata,
	}
}

func processingTargetSQLSpecForPhase(phase string) (processingTargetSQLSpec, error) {
	switch strings.TrimSpace(phase) {
	case processingTargetPhaseTermTranslation:
		return processingTargetSQLSpec{
			countSQL: processingTargetTermCountSQL,
			listSQL:  processingTargetTermListSQL,
			metadataLabels: []string{
				processingTargetMetadataTermKind,
				processingTargetMetadataTranslatedText,
			},
		}, nil
	case processingTargetPhasePersonaGeneration:
		return processingTargetSQLSpec{
			countSQL: processingTargetPersonaCountSQL,
			listSQL:  processingTargetPersonaListSQL,
			metadataLabels: []string{
				processingTargetMetadataFormID,
				processingTargetMetadataEditorID,
				processingTargetMetadataRecordType,
				processingTargetMetadataRace,
				processingTargetMetadataSex,
				processingTargetMetadataClass,
				processingTargetMetadataVoiceType,
				processingTargetMetadataEvidenceCount,
			},
		}, nil
	case processingTargetPhaseBodyTranslation:
		return processingTargetSQLSpec{
			countSQL: processingTargetBodyCountSQL,
			listSQL:  processingTargetBodyListSQL,
			metadataLabels: []string{
				processingTargetMetadataSourceText,
				processingTargetMetadataTranslatedText,
				processingTargetMetadataRecordType,
				processingTargetMetadataSubrecordType,
				processingTargetMetadataFormID,
				processingTargetMetadataEditorID,
				processingTargetMetadataAppliedPersona,
				processingTargetMetadataOutputStatus,
			},
		}, nil
	case processingTargetPhaseTranslationComplete:
		return processingTargetSQLSpec{
			countSQL: processingTargetCompleteCountSQL,
			listSQL:  processingTargetCompleteListSQL,
			metadataLabels: []string{
				processingTargetMetadataSourceText,
				processingTargetMetadataTranslatedText,
				processingTargetMetadataRecordType,
				processingTargetMetadataSubrecordType,
				processingTargetMetadataFormID,
				processingTargetMetadataEditorID,
				processingTargetMetadataOutputStatus,
			},
		}, nil
	default:
		return processingTargetSQLSpec{}, fmt.Errorf("unsupported processing target phase: %s", phase)
	}
}

func processingTargetSearchPattern(query string) string {
	escaped := strings.ReplaceAll(query, processingTargetSearchPatternEscapedSlash, `\\`)
	escaped = strings.ReplaceAll(escaped, `%`, `\%`)
	escaped = strings.ReplaceAll(escaped, `_`, `\_`)
	return "%" + strings.ToLower(escaped) + "%"
}

func processingTargetNullStringValue(value sql.NullString) string {
	if !value.Valid {
		return ""
	}
	return value.String
}

func nonBlankStrings(values ...string) []string {
	result := make([]string, 0, len(values))
	for _, value := range values {
		trimmed := strings.TrimSpace(value)
		if trimmed != "" {
			result = append(result, trimmed)
		}
	}
	return result
}

const processingTargetTermCountSQL = `
SELECT COUNT(1)
FROM DICTIONARY_ENTRY de
WHERE de.translation_job_id = ?
  AND de.dictionary_lifecycle = 'job'
  AND (
    ? = ''
    OR lower(de.source_term || ' ' || de.translated_term || ' ' || de.term_kind) LIKE ? ESCAPE '\'
  )`

const processingTargetTermListSQL = `
SELECT
  'dictionary-entry:' || de.id AS id,
  de.source_term AS name,
  'AI サービスへ送り、確定訳語として翻訳ジョブ内辞書へ保存する用語です。' AS detail,
  de.source_term AS title_part_1,
  de.term_kind AS title_part_2,
  de.translated_term AS title_part_3,
  de.term_kind AS metadata_value_1,
  de.translated_term AS metadata_value_2,
  '' AS metadata_value_3,
  '' AS metadata_value_4,
  NULL AS metadata_value_5,
  NULL AS metadata_value_6,
  NULL AS metadata_value_7,
  NULL AS metadata_value_8,
  NULL AS metadata_value_9,
  NULL AS metadata_value_10
FROM DICTIONARY_ENTRY de
WHERE de.translation_job_id = ?
  AND de.dictionary_lifecycle = 'job'
  AND (
    ? = ''
    OR lower(de.source_term || ' ' || de.translated_term || ' ' || de.term_kind) LIKE ? ESCAPE '\'
  )
ORDER BY de.source_term ASC, de.id ASC
LIMIT ? OFFSET ?`

const processingTargetPersonaCountSQL = `
SELECT COUNT(1)
FROM PERSONA p
INNER JOIN NPC_PROFILE np ON np.id = p.npc_profile_id
LEFT JOIN NPC_RECORD nr ON nr.translation_record_id = (
  SELECT MAX(nr2.translation_record_id) FROM NPC_RECORD nr2 WHERE nr2.npc_profile_id = np.id
)
WHERE p.translation_job_id = ?
  AND (
    ? = ''
    OR lower(np.display_name || ' ' || np.form_id || ' ' || np.editor_id || ' ' || np.record_type || ' ' ||
      COALESCE(nr.race, '') || ' ' || COALESCE(nr.sex, '') || ' ' || COALESCE(nr.npc_class, '') || ' ' || COALESCE(nr.voice_type, '')) LIKE ? ESCAPE '\'
  )`

const processingTargetPersonaListSQL = `
SELECT
  'persona:' || p.id AS id,
  CASE WHEN trim(np.display_name) = '' THEN np.editor_id ELSE np.display_name END AS name,
  'NPC の原文発話、NPC 属性、会話文脈、共通ペルソナ参照からペルソナ参照情報を作る入力です。' AS detail,
  CASE WHEN trim(np.display_name) = '' THEN np.editor_id ELSE np.display_name END AS title_part_1,
  np.form_id AS title_part_2,
  np.record_type AS title_part_3,
  np.form_id AS metadata_value_1,
  np.editor_id AS metadata_value_2,
  np.record_type AS metadata_value_3,
  COALESCE(nr.race, '') AS metadata_value_4,
  nr.sex AS metadata_value_5,
  nr.npc_class AS metadata_value_6,
  nr.voice_type AS metadata_value_7,
  CAST(p.evidence_utterance_count AS TEXT) AS metadata_value_8,
  NULL AS metadata_value_9,
  NULL AS metadata_value_10
FROM PERSONA p
INNER JOIN NPC_PROFILE np ON np.id = p.npc_profile_id
LEFT JOIN NPC_RECORD nr ON nr.translation_record_id = (
  SELECT MAX(nr2.translation_record_id) FROM NPC_RECORD nr2 WHERE nr2.npc_profile_id = np.id
)
WHERE p.translation_job_id = ?
  AND (
    ? = ''
    OR lower(np.display_name || ' ' || np.form_id || ' ' || np.editor_id || ' ' || np.record_type || ' ' ||
      COALESCE(nr.race, '') || ' ' || COALESCE(nr.sex, '') || ' ' || COALESCE(nr.npc_class, '') || ' ' || COALESCE(nr.voice_type, '')) LIKE ? ESCAPE '\'
  )
ORDER BY np.display_name ASC, np.editor_id ASC, p.id ASC
LIMIT ? OFFSET ?`

const processingTargetBodyCountSQL = `
SELECT COUNT(1)
FROM JOB_TRANSLATION_FIELD jtf
INNER JOIN TRANSLATION_FIELD tf ON tf.id = jtf.translation_field_id
INNER JOIN TRANSLATION_RECORD tr ON tr.id = tf.translation_record_id
LEFT JOIN PERSONA p ON p.id = jtf.applied_persona_id
LEFT JOIN NPC_PROFILE np ON np.id = p.npc_profile_id
WHERE jtf.translation_job_id = ?
  AND jtf.output_status != 'dictionary_exact_match'
  AND (
    ? = ''
    OR lower(tf.source_text || ' ' || jtf.translated_text || ' ' || tr.record_type || ' ' || tf.subrecord_type || ' ' ||
      tr.form_id || ' ' || tr.editor_id || ' ' || COALESCE(np.display_name, '') || ' ' || jtf.output_status) LIKE ? ESCAPE '\'
  )`

const processingTargetBodyListSQL = `
SELECT
  'job-translation-field:' || jtf.id AS id,
  CASE WHEN trim(tf.source_text) = '' THEN tr.editor_id ELSE tf.source_text END AS name,
  '辞書とペルソナを参照して訳文を作る翻訳項目です。' AS detail,
  tr.editor_id AS title_part_1,
  tf.subrecord_type AS title_part_2,
  tr.form_id AS title_part_3,
  tf.source_text AS metadata_value_1,
  jtf.translated_text AS metadata_value_2,
  tr.record_type AS metadata_value_3,
  tf.subrecord_type AS metadata_value_4,
  tr.form_id AS metadata_value_5,
  tr.editor_id AS metadata_value_6,
  COALESCE(np.display_name, '') AS metadata_value_7,
  jtf.output_status AS metadata_value_8,
  NULL AS metadata_value_9,
  NULL AS metadata_value_10
FROM JOB_TRANSLATION_FIELD jtf
INNER JOIN TRANSLATION_FIELD tf ON tf.id = jtf.translation_field_id
INNER JOIN TRANSLATION_RECORD tr ON tr.id = tf.translation_record_id
LEFT JOIN PERSONA p ON p.id = jtf.applied_persona_id
LEFT JOIN NPC_PROFILE np ON np.id = p.npc_profile_id
WHERE jtf.translation_job_id = ?
  AND jtf.output_status != 'dictionary_exact_match'
  AND (
    ? = ''
    OR lower(tf.source_text || ' ' || jtf.translated_text || ' ' || tr.record_type || ' ' || tf.subrecord_type || ' ' ||
      tr.form_id || ' ' || tr.editor_id || ' ' || COALESCE(np.display_name, '') || ' ' || jtf.output_status) LIKE ? ESCAPE '\'
  )
ORDER BY tr.id ASC, tf.field_order ASC, jtf.id ASC
LIMIT ? OFFSET ?`

const processingTargetCompleteCountSQL = `
SELECT COUNT(1)
FROM JOB_TRANSLATION_FIELD jtf
INNER JOIN TRANSLATION_FIELD tf ON tf.id = jtf.translation_field_id
INNER JOIN TRANSLATION_RECORD tr ON tr.id = tf.translation_record_id
WHERE jtf.translation_job_id = ?
  AND (
    ? = ''
    OR lower(tf.source_text || ' ' || jtf.translated_text || ' ' || tr.record_type || ' ' || tf.subrecord_type || ' ' ||
      tr.form_id || ' ' || tr.editor_id || ' ' || jtf.output_status) LIKE ? ESCAPE '\'
  )`

const processingTargetCompleteListSQL = `
SELECT
  'job-translation-field:' || jtf.id AS id,
  CASE WHEN trim(jtf.translated_text) = '' THEN tf.source_text ELSE jtf.translated_text END AS name,
  '本文翻訳で保持された訳文です。出力管理へ進む前に確認する翻訳項目です。' AS detail,
  tr.editor_id AS title_part_1,
  tf.subrecord_type AS title_part_2,
  tr.form_id AS title_part_3,
  tf.source_text AS metadata_value_1,
  jtf.translated_text AS metadata_value_2,
  tr.record_type AS metadata_value_3,
  tf.subrecord_type AS metadata_value_4,
  tr.form_id AS metadata_value_5,
  tr.editor_id AS metadata_value_6,
  jtf.output_status AS metadata_value_7,
  NULL AS metadata_value_8,
  NULL AS metadata_value_9,
  NULL AS metadata_value_10
FROM JOB_TRANSLATION_FIELD jtf
INNER JOIN TRANSLATION_FIELD tf ON tf.id = jtf.translation_field_id
INNER JOIN TRANSLATION_RECORD tr ON tr.id = tf.translation_record_id
WHERE jtf.translation_job_id = ?
  AND (
    ? = ''
    OR lower(tf.source_text || ' ' || jtf.translated_text || ' ' || tr.record_type || ' ' || tf.subrecord_type || ' ' ||
      tr.form_id || ' ' || tr.editor_id || ' ' || jtf.output_status) LIKE ? ESCAPE '\'
  )
ORDER BY tr.id ASC, tf.field_order ASC, jtf.id ASC
LIMIT ? OFFSET ?`
