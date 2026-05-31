package repository

import (
	"context"
	"path/filepath"
	"testing"
	"time"
)

func TestSQLiteProcessingTargetListTermTranslationUsesAITargetTermPopulation(t *testing.T) {
	repo := newSQLiteJobLifecycleRepositoryForTest(t, filepath.Join(t.TempDir(), "db", "processing-target.sqlite3"))
	ctx := context.Background()
	jobID := seedProcessingTargetTranslationJob(t, repo, 8101)
	npcRecordID := seedProcessingTargetTranslationRecord(t, repo, 8101, "000001", "RiverwoodGuard", "NPC_")
	weaponRecordID := seedProcessingTargetTranslationRecord(t, repo, 8101, "000002", "IronSword", "WEAP")
	seedProcessingTargetTranslationField(t, repo, npcRecordID, "FULL", "Guard", 1)
	seedProcessingTargetTranslationField(t, repo, npcRecordID, "DESC", "Guard", 2)
	seedProcessingTargetTranslationField(t, repo, weaponRecordID, "FULL", "Iron Sword", 1)
	seedProcessingTargetDictionaryEntry(t, repo, nil, "master", "NPC_:FULL", "Guard", "衛兵", "NPC name")
	seedProcessingTargetDictionaryEntry(t, repo, &jobID, "job", "WEAP:FULL", "Iron Sword", "鉄の剣", "weapon")

	result, err := repo.ListProcessingTargets(ctx, ProcessingTargetListQuery{
		JobID:    jobID,
		Phase:    "term_translation",
		Page:     1,
		PageSize: 10,
	})
	if err != nil {
		t.Fatalf("list term processing targets failed: %v", err)
	}

	if result.TotalCount != 1 {
		t.Fatalf("expected only shared-dictionary-miss term target, got total=%d items=%#v", result.TotalCount, result.Items)
	}
	if len(result.Items) != 1 {
		t.Fatalf("expected one term target row, got %#v", result.Items)
	}
	item := result.Items[0]
	if item.Name != "Iron Sword" {
		t.Fatalf("expected AI target source term name, got %#v", item)
	}
	if !processingTargetTestStringSlicesEqual(item.TitleParts, []string{"Iron Sword", "鉄の剣", "WEAP:FULL"}) {
		t.Fatalf("expected term title parts to be source, translation, record type, got %#v", item.TitleParts)
	}
	if !processingTargetTestHasMetadata(item, processingTargetMetadataTermKind, "WEAP:FULL") {
		t.Fatalf("expected record type metadata for term target, got %#v", item.Metadata)
	}
	if !processingTargetTestHasMetadata(item, processingTargetMetadataTranslatedText, "鉄の剣") {
		t.Fatalf("expected job dictionary translated candidate metadata, got %#v", item.Metadata)
	}
}

func TestSQLiteProcessingTargetListPersonaGenerationKeepsPersonaPopulation(t *testing.T) {
	repo := newSQLiteJobLifecycleRepositoryForTest(t, filepath.Join(t.TempDir(), "db", "processing-target.sqlite3"))
	ctx := context.Background()
	jobID := seedProcessingTargetTranslationJob(t, repo, 8201)
	profileID := seedProcessingTargetNPCProfile(t, repo, "000101", "WhiterunGuard", "Whiterun Guard")
	seedProcessingTargetPersona(t, repo, jobID, profileID, 3)

	result, err := repo.ListProcessingTargets(ctx, ProcessingTargetListQuery{
		JobID:    jobID,
		Phase:    "persona_generation",
		Page:     1,
		PageSize: 10,
	})
	if err != nil {
		t.Fatalf("list persona processing targets failed: %v", err)
	}

	if result.TotalCount != 1 || len(result.Items) != 1 {
		t.Fatalf("expected persona target population to remain one row, got total=%d items=%#v", result.TotalCount, result.Items)
	}
	if result.Items[0].Name != "Whiterun Guard" {
		t.Fatalf("expected persona display name, got %#v", result.Items[0])
	}
}

func TestSQLiteProcessingTargetListBodyTranslationExcludesDictionaryExactMatches(t *testing.T) {
	repo := newSQLiteJobLifecycleRepositoryForTest(t, filepath.Join(t.TempDir(), "db", "processing-target.sqlite3"))
	ctx := context.Background()
	jobID := seedProcessingTargetTranslationJob(t, repo, 8301)
	recordID := seedProcessingTargetTranslationRecord(t, repo, 8301, "000201", "QuestLine", "QUST")
	exactFieldID := seedProcessingTargetTranslationField(t, repo, recordID, "FULL", "Known text", 1)
	providerFieldID := seedProcessingTargetTranslationField(t, repo, recordID, "DESC", "Unknown text", 2)
	seedProcessingTargetJobTranslationField(t, repo, jobID, exactFieldID, "既知", "dictionary_exact_match")
	seedProcessingTargetJobTranslationField(t, repo, jobID, providerFieldID, "", "pending")

	result, err := repo.ListProcessingTargets(ctx, ProcessingTargetListQuery{
		JobID:    jobID,
		Phase:    "body_translation",
		Page:     1,
		PageSize: 10,
	})
	if err != nil {
		t.Fatalf("list body processing targets failed: %v", err)
	}

	if result.TotalCount != 1 || len(result.Items) != 1 {
		t.Fatalf("expected only provider target body row, got total=%d items=%#v", result.TotalCount, result.Items)
	}
	if result.Items[0].Name != "Unknown text" {
		t.Fatalf("expected non-dictionary-exact body target, got %#v", result.Items[0])
	}
}

func TestSQLiteProcessingTargetListSearchesTermTranslationByTargetNameSourceAndTranslatedCandidate(t *testing.T) {
	repo := newSQLiteJobLifecycleRepositoryForTest(t, filepath.Join(t.TempDir(), "db", "processing-target.sqlite3"))
	ctx := context.Background()
	jobID := seedProcessingTargetTranslationJob(t, repo, 8401)
	recordID := seedProcessingTargetTranslationRecord(t, repo, 8401, "000301", "SearchTerms", "MISC")
	seedProcessingTargetTranslationField(t, repo, recordID, "FULL", "Moon Sugar", 1)
	seedProcessingTargetTranslationField(t, repo, recordID, "FULL", "Dwemer Lever", 2)
	seedProcessingTargetDictionaryEntry(t, repo, &jobID, "job", "MISC:FULL", "Moon Sugar", "月の砂糖", "item")
	seedProcessingTargetDictionaryEntry(t, repo, &jobID, "job", "MISC:FULL", "Dwemer Lever", "ドゥーマーのレバー", "item")

	allTargets, err := repo.ListProcessingTargets(ctx, ProcessingTargetListQuery{
		JobID:    jobID,
		Phase:    "term_translation",
		Page:     1,
		PageSize: 10,
	})
	if err != nil {
		t.Fatalf("list all term processing targets failed: %v", err)
	}
	if allTargets.TotalCount != 2 || len(allTargets.Items) != 2 {
		t.Fatalf("expected two term targets without search, got total=%d items=%#v", allTargets.TotalCount, allTargets.Items)
	}

	matchedByName, err := repo.ListProcessingTargets(ctx, ProcessingTargetListQuery{
		JobID:       jobID,
		Phase:       "term_translation",
		Page:        1,
		PageSize:    10,
		SearchQuery: "Moon Sugar",
	})
	if err != nil {
		t.Fatalf("search term processing targets by target name failed: %v", err)
	}
	if matchedByName.TotalCount != 1 || len(matchedByName.Items) != 1 || matchedByName.Items[0].Name != "Moon Sugar" {
		t.Fatalf("expected one term target matched by target name, got total=%d items=%#v", matchedByName.TotalCount, matchedByName.Items)
	}

	matchedByTranslatedCandidate, err := repo.ListProcessingTargets(ctx, ProcessingTargetListQuery{
		JobID:       jobID,
		Phase:       "term_translation",
		Page:        1,
		PageSize:    10,
		SearchQuery: "ドゥーマー",
	})
	if err != nil {
		t.Fatalf("search term processing targets by translated candidate failed: %v", err)
	}
	if matchedByTranslatedCandidate.TotalCount != 1 || len(matchedByTranslatedCandidate.Items) != 1 || matchedByTranslatedCandidate.Items[0].Name != "Dwemer Lever" {
		t.Fatalf("expected one term target matched by translated candidate, got total=%d items=%#v", matchedByTranslatedCandidate.TotalCount, matchedByTranslatedCandidate.Items)
	}

	noTargets, err := repo.ListProcessingTargets(ctx, ProcessingTargetListQuery{
		JobID:       jobID,
		Phase:       "term_translation",
		Page:        1,
		PageSize:    10,
		SearchQuery: "not-found-term",
	})
	if err != nil {
		t.Fatalf("search term processing targets with no match failed: %v", err)
	}
	if noTargets.TotalCount != 0 || len(noTargets.Items) != 0 {
		t.Fatalf("expected no term targets for unmatched search, got total=%d items=%#v", noTargets.TotalCount, noTargets.Items)
	}
}

func TestSQLiteProcessingTargetListSearchesPersonaGenerationByNameIdentityAndNPCAttributes(t *testing.T) {
	repo := newSQLiteJobLifecycleRepositoryForTest(t, filepath.Join(t.TempDir(), "db", "processing-target.sqlite3"))
	ctx := context.Background()
	jobID := seedProcessingTargetTranslationJob(t, repo, 8501)
	guardRecordID := seedProcessingTargetTranslationRecord(t, repo, 8501, "000401", "WhiterunGuard", "NPC_")
	mageRecordID := seedProcessingTargetTranslationRecord(t, repo, 8501, "000402", "WinterholdMage", "NPC_")
	guardProfileID := seedProcessingTargetNPCProfile(t, repo, "000401", "WhiterunGuard", "Whiterun Guard")
	mageProfileID := seedProcessingTargetNPCProfile(t, repo, "000402", "WinterholdMage", "Winterhold Mage")
	seedProcessingTargetNPCRecord(t, repo, guardRecordID, guardProfileID, "Nord", "Male", "Guard", "MaleNord")
	seedProcessingTargetNPCRecord(t, repo, mageRecordID, mageProfileID, "Breton", "Female", "Mage", "FemaleEvenToned")
	seedProcessingTargetPersona(t, repo, jobID, guardProfileID, 2)
	seedProcessingTargetPersona(t, repo, jobID, mageProfileID, 4)

	allTargets, err := repo.ListProcessingTargets(ctx, ProcessingTargetListQuery{
		JobID:    jobID,
		Phase:    "persona_generation",
		Page:     1,
		PageSize: 10,
	})
	if err != nil {
		t.Fatalf("list all persona processing targets failed: %v", err)
	}
	if allTargets.TotalCount != 2 || len(allTargets.Items) != 2 {
		t.Fatalf("expected two persona targets without search, got total=%d items=%#v", allTargets.TotalCount, allTargets.Items)
	}

	matchedByIdentity, err := repo.ListProcessingTargets(ctx, ProcessingTargetListQuery{
		JobID:       jobID,
		Phase:       "persona_generation",
		Page:        1,
		PageSize:    10,
		SearchQuery: "000401",
	})
	if err != nil {
		t.Fatalf("search persona processing targets by FormID failed: %v", err)
	}
	if matchedByIdentity.TotalCount != 1 || len(matchedByIdentity.Items) != 1 || matchedByIdentity.Items[0].Name != "Whiterun Guard" {
		t.Fatalf("expected one persona target matched by FormID, got total=%d items=%#v", matchedByIdentity.TotalCount, matchedByIdentity.Items)
	}

	matchedByAttribute, err := repo.ListProcessingTargets(ctx, ProcessingTargetListQuery{
		JobID:       jobID,
		Phase:       "persona_generation",
		Page:        1,
		PageSize:    10,
		SearchQuery: "FemaleEvenToned",
	})
	if err != nil {
		t.Fatalf("search persona processing targets by NPC attribute failed: %v", err)
	}
	if matchedByAttribute.TotalCount != 1 || len(matchedByAttribute.Items) != 1 || matchedByAttribute.Items[0].Name != "Winterhold Mage" {
		t.Fatalf("expected one persona target matched by NPC attribute, got total=%d items=%#v", matchedByAttribute.TotalCount, matchedByAttribute.Items)
	}

	noTargets, err := repo.ListProcessingTargets(ctx, ProcessingTargetListQuery{
		JobID:       jobID,
		Phase:       "persona_generation",
		Page:        1,
		PageSize:    10,
		SearchQuery: "not-found-persona",
	})
	if err != nil {
		t.Fatalf("search persona processing targets with no match failed: %v", err)
	}
	if noTargets.TotalCount != 0 || len(noTargets.Items) != 0 {
		t.Fatalf("expected no persona targets for unmatched search, got total=%d items=%#v", noTargets.TotalCount, noTargets.Items)
	}
}

func TestSQLiteProcessingTargetListSearchesBodyTranslationByNameSourceAndTranslatedText(t *testing.T) {
	repo := newSQLiteJobLifecycleRepositoryForTest(t, filepath.Join(t.TempDir(), "db", "processing-target.sqlite3"))
	ctx := context.Background()
	jobID := seedProcessingTargetTranslationJob(t, repo, 8601)
	bookRecordID := seedProcessingTargetTranslationRecord(t, repo, 8601, "000501", "BookOfDawn", "BOOK")
	questRecordID := seedProcessingTargetTranslationRecord(t, repo, 8601, "000502", "QuestMarkerName", "QUST")
	bookFieldID := seedProcessingTargetTranslationField(t, repo, bookRecordID, "DESC", "The sun rises over Skyrim.", 1)
	blankFieldID := seedProcessingTargetTranslationField(t, repo, questRecordID, "FULL", "", 1)
	seedProcessingTargetJobTranslationField(t, repo, jobID, bookFieldID, "太陽がスカイリムに昇る。", "translated")
	seedProcessingTargetJobTranslationField(t, repo, jobID, blankFieldID, "クエスト目印", "translated")

	allTargets, err := repo.ListProcessingTargets(ctx, ProcessingTargetListQuery{
		JobID:    jobID,
		Phase:    "body_translation",
		Page:     1,
		PageSize: 10,
	})
	if err != nil {
		t.Fatalf("list all body processing targets failed: %v", err)
	}
	if allTargets.TotalCount != 2 || len(allTargets.Items) != 2 {
		t.Fatalf("expected two body targets without search, got total=%d items=%#v", allTargets.TotalCount, allTargets.Items)
	}

	matchedBySource, err := repo.ListProcessingTargets(ctx, ProcessingTargetListQuery{
		JobID:       jobID,
		Phase:       "body_translation",
		Page:        1,
		PageSize:    10,
		SearchQuery: "sun rises",
	})
	if err != nil {
		t.Fatalf("search body processing targets by source text failed: %v", err)
	}
	if matchedBySource.TotalCount != 1 || len(matchedBySource.Items) != 1 || matchedBySource.Items[0].Name != "The sun rises over Skyrim." {
		t.Fatalf("expected one body target matched by source text, got total=%d items=%#v", matchedBySource.TotalCount, matchedBySource.Items)
	}

	matchedByName, err := repo.ListProcessingTargets(ctx, ProcessingTargetListQuery{
		JobID:       jobID,
		Phase:       "body_translation",
		Page:        1,
		PageSize:    10,
		SearchQuery: "QuestMarkerName",
	})
	if err != nil {
		t.Fatalf("search body processing targets by row name failed: %v", err)
	}
	if matchedByName.TotalCount != 1 || len(matchedByName.Items) != 1 || matchedByName.Items[0].Name != "QuestMarkerName" {
		t.Fatalf("expected one body target matched by row name, got total=%d items=%#v", matchedByName.TotalCount, matchedByName.Items)
	}

	matchedByTranslatedText, err := repo.ListProcessingTargets(ctx, ProcessingTargetListQuery{
		JobID:       jobID,
		Phase:       "body_translation",
		Page:        1,
		PageSize:    10,
		SearchQuery: "スカイリム",
	})
	if err != nil {
		t.Fatalf("search body processing targets by translated text failed: %v", err)
	}
	if matchedByTranslatedText.TotalCount != 1 || len(matchedByTranslatedText.Items) != 1 || matchedByTranslatedText.Items[0].Name != "The sun rises over Skyrim." {
		t.Fatalf("expected one body target matched by translated text, got total=%d items=%#v", matchedByTranslatedText.TotalCount, matchedByTranslatedText.Items)
	}

	noTargets, err := repo.ListProcessingTargets(ctx, ProcessingTargetListQuery{
		JobID:       jobID,
		Phase:       "body_translation",
		Page:        1,
		PageSize:    10,
		SearchQuery: "not-found-body",
	})
	if err != nil {
		t.Fatalf("search body processing targets with no match failed: %v", err)
	}
	if noTargets.TotalCount != 0 || len(noTargets.Items) != 0 {
		t.Fatalf("expected no body targets for unmatched search, got total=%d items=%#v", noTargets.TotalCount, noTargets.Items)
	}
}

func seedProcessingTargetTranslationJob(t *testing.T, repo *SQLiteJobLifecycleRepository, xEditID int64) int64 {
	t.Helper()
	now := time.Date(2026, 5, 28, 0, 0, 0, 0, time.UTC).Format(time.RFC3339)
	if _, err := repo.db.ExecContext(
		context.Background(),
		`INSERT INTO X_EDIT_EXTRACTED_DATA (id, source_file_path, source_content_hash, source_tool, target_plugin_name, target_plugin_type, record_count, imported_at)
		 VALUES (?, ?, ?, ?, ?, ?, ?, ?)`,
		xEditID,
		"/tmp/processing-target.json",
		"processing-target-hash",
		"xedit",
		"Skyrim.esm",
		"esm",
		0,
		now,
	); err != nil {
		t.Fatalf("seed X_EDIT_EXTRACTED_DATA failed: %v", err)
	}
	result, err := repo.db.ExecContext(
		context.Background(),
		`INSERT INTO TRANSLATION_JOB (x_edit_extracted_data_id, job_name, state, progress_percent, created_at)
		 VALUES (?, ?, ?, ?, ?)`,
		xEditID,
		"processing target test",
		"ready",
		0,
		now,
	)
	if err != nil {
		t.Fatalf("seed TRANSLATION_JOB failed: %v", err)
	}
	id, err := result.LastInsertId()
	if err != nil {
		t.Fatalf("read seeded TRANSLATION_JOB id failed: %v", err)
	}
	return id
}

func seedProcessingTargetTranslationRecord(
	t *testing.T,
	repo *SQLiteJobLifecycleRepository,
	xEditID int64,
	formID string,
	editorID string,
	recordType string,
) int64 {
	t.Helper()
	result, err := repo.db.ExecContext(
		context.Background(),
		`INSERT INTO TRANSLATION_RECORD (x_edit_extracted_data_id, form_id, editor_id, record_type)
		 VALUES (?, ?, ?, ?)`,
		xEditID,
		formID,
		editorID,
		recordType,
	)
	if err != nil {
		t.Fatalf("seed TRANSLATION_RECORD failed: %v", err)
	}
	id, err := result.LastInsertId()
	if err != nil {
		t.Fatalf("read seeded TRANSLATION_RECORD id failed: %v", err)
	}
	return id
}

func seedProcessingTargetTranslationField(
	t *testing.T,
	repo *SQLiteJobLifecycleRepository,
	recordID int64,
	subrecordType string,
	sourceText string,
	fieldOrder int,
) int64 {
	t.Helper()
	result, err := repo.db.ExecContext(
		context.Background(),
		`INSERT INTO TRANSLATION_FIELD (translation_record_id, subrecord_type, source_text, field_order)
		 VALUES (?, ?, ?, ?)`,
		recordID,
		subrecordType,
		sourceText,
		fieldOrder,
	)
	if err != nil {
		t.Fatalf("seed TRANSLATION_FIELD failed: %v", err)
	}
	id, err := result.LastInsertId()
	if err != nil {
		t.Fatalf("read seeded TRANSLATION_FIELD id failed: %v", err)
	}
	return id
}

func seedProcessingTargetDictionaryEntry(
	t *testing.T,
	repo *SQLiteJobLifecycleRepository,
	jobID *int64,
	lifecycle string,
	scope string,
	source string,
	translation string,
	kind string,
) {
	t.Helper()
	now := time.Date(2026, 5, 28, 0, 0, 0, 0, time.UTC).Format(time.RFC3339)
	if _, err := repo.db.ExecContext(
		context.Background(),
		`INSERT INTO DICTIONARY_ENTRY (translation_job_id, dictionary_lifecycle, dictionary_scope, dictionary_source, source_term, translated_term, term_kind, reusable, created_at, updated_at)
		 VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		jobID,
		lifecycle,
		scope,
		"processing_target_test",
		source,
		translation,
		kind,
		1,
		now,
		now,
	); err != nil {
		t.Fatalf("seed DICTIONARY_ENTRY failed: %v", err)
	}
}

func seedProcessingTargetNPCProfile(t *testing.T, repo *SQLiteJobLifecycleRepository, formID string, editorID string, displayName string) int64 {
	t.Helper()
	now := time.Date(2026, 5, 28, 0, 0, 0, 0, time.UTC).Format(time.RFC3339)
	result, err := repo.db.ExecContext(
		context.Background(),
		`INSERT INTO NPC_PROFILE (target_plugin_name, form_id, record_type, editor_id, display_name, created_at, updated_at)
		 VALUES (?, ?, ?, ?, ?, ?, ?)`,
		"Skyrim.esm",
		formID,
		"NPC_",
		editorID,
		displayName,
		now,
		now,
	)
	if err != nil {
		t.Fatalf("seed NPC_PROFILE failed: %v", err)
	}
	id, err := result.LastInsertId()
	if err != nil {
		t.Fatalf("read seeded NPC_PROFILE id failed: %v", err)
	}
	return id
}

func seedProcessingTargetPersona(t *testing.T, repo *SQLiteJobLifecycleRepository, jobID int64, profileID int64, evidenceCount int) {
	t.Helper()
	now := time.Date(2026, 5, 28, 0, 0, 0, 0, time.UTC).Format(time.RFC3339)
	if _, err := repo.db.ExecContext(
		context.Background(),
		`INSERT INTO PERSONA (translation_job_id, npc_profile_id, persona_lifecycle, persona_scope, persona_source, persona_description, speech_style, personality_summary, evidence_utterance_count, created_at, updated_at)
		 VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		jobID,
		profileID,
		"job",
		"NPC_",
		"processing_target_test",
		"Guard persona",
		"formal",
		"City guard",
		evidenceCount,
		now,
		now,
	); err != nil {
		t.Fatalf("seed PERSONA failed: %v", err)
	}
}

func seedProcessingTargetNPCRecord(
	t *testing.T,
	repo *SQLiteJobLifecycleRepository,
	recordID int64,
	profileID int64,
	race string,
	sex string,
	npcClass string,
	voiceType string,
) {
	t.Helper()
	if _, err := repo.db.ExecContext(
		context.Background(),
		`INSERT INTO NPC_RECORD (translation_record_id, npc_profile_id, race, sex, npc_class, voice_type)
		 VALUES (?, ?, ?, ?, ?, ?)`,
		recordID,
		profileID,
		race,
		sex,
		npcClass,
		voiceType,
	); err != nil {
		t.Fatalf("seed NPC_RECORD failed: %v", err)
	}
}

func seedProcessingTargetJobTranslationField(
	t *testing.T,
	repo *SQLiteJobLifecycleRepository,
	jobID int64,
	fieldID int64,
	translatedText string,
	outputStatus string,
) {
	t.Helper()
	now := time.Date(2026, 5, 28, 0, 0, 0, 0, time.UTC).Format(time.RFC3339)
	if _, err := repo.db.ExecContext(
		context.Background(),
		`INSERT INTO JOB_TRANSLATION_FIELD (translation_job_id, translation_field_id, translated_text, output_status, retry_count, updated_at)
		 VALUES (?, ?, ?, ?, ?, ?)`,
		jobID,
		fieldID,
		translatedText,
		outputStatus,
		0,
		now,
	); err != nil {
		t.Fatalf("seed JOB_TRANSLATION_FIELD failed: %v", err)
	}
}

// TestProcessingTargetTermSQL_filtersBy13RECTypes は
// processingTargetTermCountSQL と processingTargetTermListSQL が
// 13 種別 REC だけを集計・一覧化し、13 種別外 REC（DOOR:FULL / FLOR:FULL / FURN:FULL）を
// 件数にも一覧にも含まないことを証明する。
func TestProcessingTargetTermSQL_filtersBy13RECTypes(t *testing.T) {
	// Arrange: 13 種別 REC と 13 種別外 REC を両方 seed する
	repo := newSQLiteJobLifecycleRepositoryForTest(t, filepath.Join(t.TempDir(), "db", "rec-filter.sqlite3"))
	ctx := context.Background()
	jobID := seedProcessingTargetTranslationJob(t, repo, 9001)

	// 13 種別内: WEAP:FULL
	weaponRecordID := seedProcessingTargetTranslationRecord(t, repo, 9001, "W00001", "IronSword", "WEAP")
	seedProcessingTargetTranslationField(t, repo, weaponRecordID, "FULL", "Iron Sword", 1)

	// 13 種別内: MISC:FULL
	miscRecordID := seedProcessingTargetTranslationRecord(t, repo, 9001, "M00001", "MoonSugar", "MISC")
	seedProcessingTargetTranslationField(t, repo, miscRecordID, "FULL", "Moon Sugar", 1)

	// 13 種別外: DOOR:FULL（除外対象）
	doorRecordID := seedProcessingTargetTranslationRecord(t, repo, 9001, "D00001", "OakDoor", "DOOR")
	seedProcessingTargetTranslationField(t, repo, doorRecordID, "FULL", "Oak Door", 1)

	// 13 種別外: FLOR:FULL（除外対象）
	florRecordID := seedProcessingTargetTranslationRecord(t, repo, 9001, "F00001", "YellowFlower", "FLOR")
	seedProcessingTargetTranslationField(t, repo, florRecordID, "FULL", "Yellow Flower", 1)

	// 13 種別外: FURN:FULL（除外対象）
	furnRecordID := seedProcessingTargetTranslationRecord(t, repo, 9001, "FU0001", "WoodChair", "FURN")
	seedProcessingTargetTranslationField(t, repo, furnRecordID, "FULL", "Wood Chair", 1)

	// Act: term_translation フェーズで一覧を取得する
	result, err := repo.ListProcessingTargets(ctx, ProcessingTargetListQuery{
		JobID:    jobID,
		Phase:    "term_translation",
		Page:     1,
		PageSize: 20,
	})
	if err != nil {
		t.Fatalf("ListProcessingTargets failed: %v", err)
	}

	// Assert: count と list の両方が 13 種別 REC の 2 件だけを返す
	// DOOR:FULL / FLOR:FULL / FURN:FULL が結果に出れば REC 絞り込み欠落の退行
	if result.TotalCount != 2 {
		t.Errorf("REC 絞り込み欠落: 13 種別 REC 2 件だけを期待, got TotalCount=%d items=%#v", result.TotalCount, result.Items)
	}
	if len(result.Items) != 2 {
		t.Errorf("REC 絞り込み欠落: 13 種別 REC 2 件だけを期待, got len(Items)=%d items=%#v", len(result.Items), result.Items)
	}

	names := make(map[string]struct{}, len(result.Items))
	for _, item := range result.Items {
		names[item.Name] = struct{}{}
	}
	if _, ok := names["Iron Sword"]; !ok {
		t.Errorf("13 種別 REC WEAP:FULL が一覧に含まれていない: items=%#v", result.Items)
	}
	if _, ok := names["Moon Sugar"]; !ok {
		t.Errorf("13 種別 REC MISC:FULL が一覧に含まれていない: items=%#v", result.Items)
	}
	for _, name := range []string{"Oak Door", "Yellow Flower", "Wood Chair"} {
		if _, ok := names[name]; ok {
			t.Errorf("13 種別外 REC が一覧に含まれている（REC 絞り込み欠落）: name=%q items=%#v", name, result.Items)
		}
	}
}

// TestProcessingTargetTermSQL_excludesByDictionaryScopeRECORDFIELD は
// dictionary_scope を RECORD:FIELD 形式（例: NPC_:FULL）で投入したデータに対し、
// 共通辞書一致除外（master lifecycle、同一原語、同一 dictionary_scope）が成立することを証明する。
// dictionary_scope が RECORD:FIELD でなく record_type 単体で比較されている場合は除外が成立しない退行が発生する。
func TestProcessingTargetTermSQL_excludesByDictionaryScopeRECORDFIELD(t *testing.T) {
	// Arrange: NPC_:FULL の原語を共通辞書に登録し、除外が成立するか確認する
	repo := newSQLiteJobLifecycleRepositoryForTest(t, filepath.Join(t.TempDir(), "db", "dict-scope.sqlite3"))
	ctx := context.Background()
	jobID := seedProcessingTargetTranslationJob(t, repo, 9101)

	// NPC_:FULL レコード: "RiverwoodGuard" の FULL フィールド原語 "Guard"
	npcRecordID := seedProcessingTargetTranslationRecord(t, repo, 9101, "N00001", "RiverwoodGuard", "NPC_")
	seedProcessingTargetTranslationField(t, repo, npcRecordID, "FULL", "Guard", 1)

	// WEAP:FULL レコード: 共通辞書に登録なし（AI 対象として残るはず）
	weaponRecordID := seedProcessingTargetTranslationRecord(t, repo, 9101, "W00002", "IronDagger", "WEAP")
	seedProcessingTargetTranslationField(t, repo, weaponRecordID, "FULL", "Iron Dagger", 1)

	// 共通辞書: NPC_:FULL スコープで "Guard" を登録する（RECORD:FIELD 形式）
	// この entry に対し dictionary_scope が RECORD:FIELD で比較されれば "Guard" は除外される
	seedProcessingTargetDictionaryEntry(t, repo, nil, "master", "NPC_:FULL", "Guard", "衛兵", "NPC name")

	// Act: term_translation フェーズで一覧を取得する
	result, err := repo.ListProcessingTargets(ctx, ProcessingTargetListQuery{
		JobID:    jobID,
		Phase:    "term_translation",
		Page:     1,
		PageSize: 10,
	})
	if err != nil {
		t.Fatalf("ListProcessingTargets failed: %v", err)
	}

	// Assert: "Guard" は共通辞書一致で除外され、"Iron Dagger" だけが AI 対象として返る
	// TotalCount = 2 なら dictionary_scope が RECORD:FIELD でなく比較されている退行（NOT EXISTS 節の REC 比較欠落）
	if result.TotalCount != 1 {
		t.Errorf("dictionary_scope が RECORD:FIELD で比較されていない退行: 共通辞書一致 1 件除外後に 1 件を期待, got TotalCount=%d items=%#v", result.TotalCount, result.Items)
	}
	if len(result.Items) != 1 {
		t.Errorf("dictionary_scope が RECORD:FIELD で比較されていない退行: 1 件を期待, got len(Items)=%d items=%#v", len(result.Items), result.Items)
	}
	if len(result.Items) == 1 && result.Items[0].Name != "Iron Dagger" {
		t.Errorf("共通辞書一致除外後に残る項目が期待と異なる: want=Iron Dagger got=%q", result.Items[0].Name)
	}
	if len(result.Items) == 1 && !processingTargetTestStringSlicesEqual(result.Items[0].TitleParts, []string{"Iron Dagger", "", "WEAP:FULL"}) {
		t.Errorf("訳語未登録の title part が原語、空訳語、レコード種別の順ではない: got=%#v", result.Items[0].TitleParts)
	}
}

// TestProcessingTargetTermSQL_separatesNPCFullAndSHRT は
// NPC_:FULL と NPC_:SHRT が同一原語を持つ場合でも別 candidate として返ることを証明する。
// NPC_:FULL と NPC_:SHRT が同一キーで集約されると TotalCount=1 かつ Items の rec が 1 種別になる退行が発生する。
func TestProcessingTargetTermSQL_separatesNPCFullAndSHRT(t *testing.T) {
	// Arrange: 同一 NPC_ レコードに FULL と SHRT フィールドを持つ seed データ
	repo := newSQLiteJobLifecycleRepositoryForTest(t, filepath.Join(t.TempDir(), "db", "npc-sep.sqlite3"))
	ctx := context.Background()
	jobID := seedProcessingTargetTranslationJob(t, repo, 9201)

	npcRecordID := seedProcessingTargetTranslationRecord(t, repo, 9201, "N00002", "WhiterunGuard", "NPC_")
	// NPC_:FULL フィールド: 原語 "Whiterun Guard"
	seedProcessingTargetTranslationField(t, repo, npcRecordID, "FULL", "Whiterun Guard", 1)
	// NPC_:SHRT フィールド: 同一原語 "Whiterun Guard"
	seedProcessingTargetTranslationField(t, repo, npcRecordID, "SHRT", "Whiterun Guard", 2)

	// Act: term_translation フェーズで一覧を取得する
	result, err := repo.ListProcessingTargets(ctx, ProcessingTargetListQuery{
		JobID:    jobID,
		Phase:    "term_translation",
		Page:     1,
		PageSize: 10,
	})
	if err != nil {
		t.Fatalf("ListProcessingTargets failed: %v", err)
	}

	// Assert: NPC_:FULL と NPC_:SHRT が別 candidate（2 件）として返る
	// TotalCount=1 なら record_type のみで集約されている退行（RECORD:FIELD 分離欠落）
	if result.TotalCount != 2 {
		t.Errorf("NPC_:FULL と NPC_:SHRT が別 candidate として返っていない（RECORD:FIELD 分離欠落）: want TotalCount=2 got=%d items=%#v", result.TotalCount, result.Items)
	}
	if len(result.Items) != 2 {
		t.Errorf("NPC_:FULL と NPC_:SHRT が別 candidate として返っていない（RECORD:FIELD 分離欠落）: want len=2 got=%d items=%#v", len(result.Items), result.Items)
	}

	recValues := make(map[string]struct{}, len(result.Items))
	for _, item := range result.Items {
		for _, m := range item.Metadata {
			if m.Label == processingTargetMetadataTermKind {
				recValues[m.Value] = struct{}{}
			}
		}
	}
	if _, ok := recValues["NPC_:FULL"]; !ok {
		t.Errorf("NPC_:FULL が candidate として含まれていない: recValues=%v items=%#v", recValues, result.Items)
	}
	if _, ok := recValues["NPC_:SHRT"]; !ok {
		t.Errorf("NPC_:SHRT が candidate として含まれていない: recValues=%v items=%#v", recValues, result.Items)
	}
}

func processingTargetTestHasMetadata(item ProcessingTargetListItem, label string, value string) bool {
	for _, metadata := range item.Metadata {
		if metadata.Label == label && metadata.Value == value {
			return true
		}
	}
	return false
}

func processingTargetTestStringSlicesEqual(actual []string, expected []string) bool {
	if len(actual) != len(expected) {
		return false
	}
	for index := range actual {
		if actual[index] != expected[index] {
			return false
		}
	}
	return true
}
