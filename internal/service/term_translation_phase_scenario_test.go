package service

// ST1-processing-target-and-execution-rec-parity シナリオテスト
//
// 証明対象:
//   - 同一 xEdit 抽出データに対し、repository.SQLiteJobLifecycleRepository.ListProcessingTargets が返す
//     REC 集合と、service.collectCandidates が返す候補 REC 集合が一致することを検証する。
//   - 13 種別外 REC（DOOR:FULL / FLOR:FULL / FURN:FULL）のみを含む入力では、
//     両者とも対応行が存在しないことを検証する。
//
// 実行種別: APIテスト相当（backend 内部の公開接点を起点にする）
// UI 入口: 本 task に画面変更なし。
//
// 検証コマンド:
//   go test ./internal/service/... -run Scenario_ProcessingTargetAndExecutionRECParity -count=1

import (
	"context"
	"fmt"
	"path/filepath"
	"sort"
	"testing"
	"time"

	sqliteinfra "aitranslationenginejp/internal/infra/sqlite/dbinit"
	"aitranslationenginejp/internal/repository"

	"github.com/jmoiron/sqlx"
)

// TestScenario_ProcessingTargetAndExecutionRECParity_MatchesOn13Types は、
// 13 種別 REC を含む xEdit データに対して repository と service の REC 集合が一致することを証明する。
func TestScenario_ProcessingTargetAndExecutionRECParity_MatchesOn13Types(t *testing.T) {
	// 証明: 同一 xEdit データで repository と service が同じ REC 集合を返す（成功経路）

	// Arrange: in-file SQLite を立て、13 種別 REC + 除外対象の混在データを投入する
	db := openScenarioSQLiteDB(t)
	xEditID := int64(9101)
	seedScenarioXEditExtractedData(t, db, xEditID)

	rec13Types := []struct {
		recordType    string
		subrecordType string
		sourceTerm    string
	}{
		{"BOOK", "FULL", "The Lusty Argonian Maid"},
		{"NPC_", "FULL", "Guard"},
		{"NPC_", "SHRT", "Guard"},
		{"ARMO", "FULL", "Iron Armor"},
		{"WEAP", "FULL", "Iron Sword"},
		{"LCTN", "FULL", "Whiterun"},
		{"CELL", "FULL", "Riverwood"},
		{"CONT", "FULL", "Chest"},
		{"MISC", "FULL", "Dwemer Cog"},
		{"ALCH", "FULL", "Skooma"},
		{"RACE", "FULL", "Nord"},
		{"INGR", "FULL", "Wheat"},
		{"SHOU", "FULL", "Unrelenting Force"},
	}
	// 13 種別外（除外対象: DOOR:FULL / FLOR:FULL / FURN:FULL）も投入して混在させる
	rec13TypesOutside := []struct {
		recordType    string
		subrecordType string
		sourceTerm    string
	}{
		{"DOOR", "FULL", "Riverwood Door"},
		{"FLOR", "FULL", "Lavender"},
		{"FURN", "FULL", "Throne"},
	}

	for i, r := range rec13Types {
		recID := seedScenarioTranslationRecord(t, db, xEditID, r.recordType, i+1)
		seedScenarioTranslationField(t, db, recID, r.subrecordType, r.sourceTerm, 1)
	}
	for i, r := range rec13TypesOutside {
		recID := seedScenarioTranslationRecord(t, db, xEditID, r.recordType, len(rec13Types)+i+1)
		seedScenarioTranslationField(t, db, recID, r.subrecordType, r.sourceTerm, 1)
	}

	jobID := seedScenarioTranslationJob(t, db, xEditID)

	repo := repository.NewSQLiteJobLifecycleRepository(db)
	processingTargetRepo, ok := repo.(repository.ProcessingTargetReadModelRepository)
	if !ok {
		t.Fatal("expected SQLiteJobLifecycleRepository to implement ProcessingTargetReadModelRepository")
	}

	// Act: repository の ListProcessingTargets を呼び、REC 集合を取得する
	// PageSize を十分大きくして全件取得する
	listResult, err := processingTargetRepo.ListProcessingTargets(context.Background(), repository.ProcessingTargetListQuery{
		JobID:    jobID,
		Phase:    "term_translation",
		Page:     1,
		PageSize: 100,
	})
	if err != nil {
		t.Fatalf("ListProcessingTargets failed: %v", err)
	}

	// Act: service の collectCandidates を呼び、候補 REC 集合を取得する
	service := newScenarioTermTranslationPhaseService(t, db)
	candidates, err := service.collectCandidates(context.Background(), xEditID)
	if err != nil {
		t.Fatalf("collectCandidates failed: %v", err)
	}

	// Assert: repository 側の REC 集合を構築する
	// ListProcessingTargets の metadata から REC（term_kind metadata）を取得する
	repoRECSet := buildRECSetFromProcessingTargetItems(t, listResult.Items)

	// Assert: service 側の REC 集合を構築する
	serviceRECSet := buildRECSetFromCandidates(candidates)

	// Assert: 両者の REC 集合が一致することを検証する
	if !equalStringSet(repoRECSet, serviceRECSet) {
		t.Errorf("REC set mismatch between repository and service:\n  repository=%v\n  service=%v",
			sortedKeys(repoRECSet), sortedKeys(serviceRECSet))
	}

	// Assert: 13 種別外 REC（DOOR:FULL / FLOR:FULL / FURN:FULL）が両者の結果に含まれていないことを検証する
	outsideRECs := []string{"DOOR:FULL", "FLOR:FULL", "FURN:FULL"}
	for _, rec := range outsideRECs {
		if repoRECSet[rec] {
			t.Errorf("repository result contains outside REC %q but should not", rec)
		}
		if serviceRECSet[rec] {
			t.Errorf("service result contains outside REC %q but should not", rec)
		}
	}
}

// TestScenario_ProcessingTargetAndExecutionRECParity_EmptyWhenOnly13TypesOutside は、
// 13 種別外 REC のみを含む xEdit データで repository と service の両者が空を返すことを証明する。
func TestScenario_ProcessingTargetAndExecutionRECParity_EmptyWhenOnly13TypesOutside(t *testing.T) {
	// 証明: 13 種別外 REC（DOOR:FULL / FLOR:FULL / FURN:FULL）のみ含む入力で両者とも対象行が存在しない（境界値）

	// Arrange: 13 種別外 REC のみを含む xEdit データを投入する
	db := openScenarioSQLiteDB(t)
	xEditID := int64(9201)
	seedScenarioXEditExtractedData(t, db, xEditID)

	outsideRecords := []struct {
		recordType    string
		subrecordType string
		sourceTerm    string
	}{
		{"DOOR", "FULL", "Riverwood Door"},
		{"FLOR", "FULL", "Lavender"},
		{"FURN", "FULL", "Throne"},
	}
	for i, r := range outsideRecords {
		recID := seedScenarioTranslationRecord(t, db, xEditID, r.recordType, i+1)
		seedScenarioTranslationField(t, db, recID, r.subrecordType, r.sourceTerm, 1)
	}

	jobID := seedScenarioTranslationJob(t, db, xEditID)

	repo := repository.NewSQLiteJobLifecycleRepository(db)
	processingTargetRepo, ok := repo.(repository.ProcessingTargetReadModelRepository)
	if !ok {
		t.Fatal("expected SQLiteJobLifecycleRepository to implement ProcessingTargetReadModelRepository")
	}

	// Act: repository の ListProcessingTargets を呼ぶ
	listResult, err := processingTargetRepo.ListProcessingTargets(context.Background(), repository.ProcessingTargetListQuery{
		JobID:    jobID,
		Phase:    "term_translation",
		Page:     1,
		PageSize: 100,
	})
	if err != nil {
		t.Fatalf("ListProcessingTargets failed: %v", err)
	}

	// Act: service の collectCandidates を呼ぶ
	service := newScenarioTermTranslationPhaseService(t, db)
	candidates, err := service.collectCandidates(context.Background(), xEditID)
	if err != nil {
		t.Fatalf("collectCandidates failed: %v", err)
	}

	// Assert: repository の結果が 0 件であることを検証する
	if listResult.TotalCount != 0 {
		t.Errorf("expected repository to return 0 items for outside REC only input, got %d", listResult.TotalCount)
	}
	if len(listResult.Items) != 0 {
		t.Errorf("expected repository to return empty items for outside REC only input, got %d items", len(listResult.Items))
	}

	// Assert: service の結果が 0 件であることを検証する
	if len(candidates) != 0 {
		t.Errorf("expected service to return 0 candidates for outside REC only input, got %d", len(candidates))
	}
}

// --- シナリオテスト補助 ---

// openScenarioSQLiteDB はシナリオテスト専用の一時 SQLite DB を開く。
func openScenarioSQLiteDB(t *testing.T) *sqlx.DB {
	t.Helper()
	dbPath := filepath.Join(t.TempDir(), "scenario-parity.sqlite3")
	db, err := sqliteinfra.OpenMasterDictionaryDatabase(context.Background(), dbPath, nil)
	if err != nil {
		t.Fatalf("open scenario SQLite DB: %v", err)
	}
	t.Cleanup(func() { _ = db.Close() })
	return db
}

// seedScenarioXEditExtractedData は xEdit 抽出データの親行を挿入する。
func seedScenarioXEditExtractedData(t *testing.T, db *sqlx.DB, xEditID int64) {
	t.Helper()
	now := time.Date(2026, 5, 31, 0, 0, 0, 0, time.UTC).Format(time.RFC3339)
	_, err := db.ExecContext(
		context.Background(),
		`INSERT INTO X_EDIT_EXTRACTED_DATA (id, source_file_path, source_content_hash, source_tool, target_plugin_name, target_plugin_type, record_count, imported_at)
		 VALUES (?, ?, ?, ?, ?, ?, ?, ?)`,
		xEditID,
		"/tmp/scenario-parity.json",
		"scenario-parity-hash",
		"xedit",
		"Skyrim.esm",
		"esm",
		0,
		now,
	)
	if err != nil {
		t.Fatalf("seed X_EDIT_EXTRACTED_DATA: %v", err)
	}
}

// seedScenarioTranslationJob は TRANSLATION_JOB を挿入し jobID を返す。
func seedScenarioTranslationJob(t *testing.T, db *sqlx.DB, xEditID int64) int64 {
	t.Helper()
	now := time.Date(2026, 5, 31, 0, 0, 0, 0, time.UTC).Format(time.RFC3339)
	result, err := db.ExecContext(
		context.Background(),
		`INSERT INTO TRANSLATION_JOB (x_edit_extracted_data_id, job_name, state, progress_percent, created_at)
		 VALUES (?, ?, ?, ?, ?)`,
		xEditID,
		"scenario parity test",
		"ready",
		0,
		now,
	)
	if err != nil {
		t.Fatalf("seed TRANSLATION_JOB: %v", err)
	}
	id, err := result.LastInsertId()
	if err != nil {
		t.Fatalf("read seeded TRANSLATION_JOB id: %v", err)
	}
	return id
}

// seedScenarioTranslationRecord は TRANSLATION_RECORD を挿入し recordID を返す。
func seedScenarioTranslationRecord(t *testing.T, db *sqlx.DB, xEditID int64, recordType string, seq int) int64 {
	t.Helper()
	formID := "SCEN" + recordType + fmt.Sprintf("%04d", seq)
	result, err := db.ExecContext(
		context.Background(),
		`INSERT INTO TRANSLATION_RECORD (x_edit_extracted_data_id, form_id, editor_id, record_type)
		 VALUES (?, ?, ?, ?)`,
		xEditID,
		formID,
		"ScenarioRec"+formID,
		recordType,
	)
	if err != nil {
		t.Fatalf("seed TRANSLATION_RECORD (recordType=%s): %v", recordType, err)
	}
	id, err := result.LastInsertId()
	if err != nil {
		t.Fatalf("read seeded TRANSLATION_RECORD id: %v", err)
	}
	return id
}

// seedScenarioTranslationField は TRANSLATION_FIELD を挿入する。
func seedScenarioTranslationField(t *testing.T, db *sqlx.DB, recordID int64, subrecordType string, sourceText string, fieldOrder int) {
	t.Helper()
	_, err := db.ExecContext(
		context.Background(),
		`INSERT INTO TRANSLATION_FIELD (translation_record_id, subrecord_type, source_text, field_order)
		 VALUES (?, ?, ?, ?)`,
		recordID,
		subrecordType,
		sourceText,
		fieldOrder,
	)
	if err != nil {
		t.Fatalf("seed TRANSLATION_FIELD (subrecordType=%s): %v", subrecordType, err)
	}
}

// newScenarioTermTranslationPhaseService はシナリオテスト用の TermTranslationPhaseService を構築する。
// SQLiteTranslationSourceRepository を translationSourceReader として使い、
// collectCandidates が SQLite から直接 TranslationRecord / TranslationField を読む構造にする。
func newScenarioTermTranslationPhaseService(t *testing.T, db *sqlx.DB) *TermTranslationPhaseService {
	t.Helper()
	translationSourceReader := repository.NewSQLiteTranslationSourceRepository(db)

	jobRepo := &fakeTermPhaseJobLifecycleRepository{
		job: repository.TranslationJob{ID: 1, XEditExtractedDataID: 1},
	}
	foundation := &fakeTermPhaseFoundationDataRepository{
		entries:           map[string]repository.DictionaryEntry{},
		createErrBySource: map[string]error{},
	}
	secretStore := repository.NewFakeSecretStore()
	transactor := fakeTermPhaseTransactor{jobRepo: jobRepo, foundation: foundation}

	svc := &TermTranslationPhaseService{
		jobLifecycleRepository:   jobRepo,
		foundationDataRepository: foundation,
		translationSourceReader:  translationSourceReader,
		transactor:               transactor,
		secretStore:              secretStore,
		provider:                 nil,
		providerSettings:         nil,
		now:                      func() time.Time { return time.Date(2026, 5, 31, 0, 0, 0, 0, time.UTC) },
	}
	return svc
}

// buildRECSetFromProcessingTargetItems は ListProcessingTargets の結果から REC 集合（重複なし）を構築する。
// 各 item の metadata label "用語種別" の値（"RECORD:FIELD" 形式）を REC として取得する。
// "用語種別" は repository.processingTargetMetadataTermKind（非公開）の実際の文字列値。
func buildRECSetFromProcessingTargetItems(t *testing.T, items []repository.ProcessingTargetListItem) map[string]bool {
	t.Helper()
	set := make(map[string]bool)
	for _, item := range items {
		for _, meta := range item.Metadata {
			if meta.Label == "用語種別" {
				set[meta.Value] = true
			}
		}
	}
	return set
}

// buildRECSetFromCandidates は collectCandidates の結果から REC 集合（重複なし）を構築する。
func buildRECSetFromCandidates(candidates []termTranslationCandidate) map[string]bool {
	set := make(map[string]bool)
	for _, c := range candidates {
		set[c.RecordType] = true
	}
	return set
}

// equalStringSet は 2 つの string set が等しいことを検証する。
func equalStringSet(a, b map[string]bool) bool {
	if len(a) != len(b) {
		return false
	}
	for k := range a {
		if !b[k] {
			return false
		}
	}
	return true
}

// sortedKeys は map[string]bool の key を昇順に返す。失敗診断用。
func sortedKeys(m map[string]bool) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}
