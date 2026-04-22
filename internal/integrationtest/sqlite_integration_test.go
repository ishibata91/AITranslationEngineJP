package integrationtest

// integration_test.go は SCN-SMR-002〜005 の SQLite repository integration test を提供する。
// 実装ファイル (transactor.go, *_repository.go) が追加されるとコンパイルが通る。
// 現時点では product code が未実装のため compile error になることが想定されている。

import (
	"context"
	"fmt"
	"path/filepath"
	"testing"
	"time"

	sqlitedbinit "aitranslationenginejp/internal/infra/sqlite/dbinit"
	"aitranslationenginejp/internal/repository"

	"github.com/jmoiron/sqlx"
)

// fixedNow は integration test 全体で共通の決定的タイムスタンプ。
var fixedNow = time.Date(2026, 4, 19, 12, 0, 0, 0, time.UTC)

// openIntegrationDB は integration test 用の一時 SQLite DB を返す。
func openIntegrationDB(t *testing.T) *sqlx.DB {
	t.Helper()
	db, err := sqlitedbinit.OpenMasterDictionaryDatabase(
		context.Background(),
		filepath.Join(t.TempDir(), "integration.sqlite3"),
		nil,
	)
	if err != nil {
		t.Fatalf("openIntegrationDB failed: %v", err)
	}
	t.Cleanup(func() {
		if closeErr := db.Close(); closeErr != nil {
			t.Fatalf("openIntegrationDB cleanup close failed: %v", closeErr)
		}
	})
	return db
}

// ---------------------------------------------------------------------------
// SCN-SMR-002 共通/job-local dictionary と persona
// ---------------------------------------------------------------------------

// TestSCN_SMR_002_CreateCommonDictionaryEntry は FoundationDataRepository で
// 共通 (global) dictionary entry を保存・取得できることを検証する。
func TestSCN_SMR_002_CreateCommonDictionaryEntry(t *testing.T) {
	ctx := context.Background()
	db := openIntegrationDB(t)
	repo := repository.NewSQLiteFoundationDataRepository(db)

	draft := repository.DictionaryEntryDraft{
		XTranslatorTranslationXMLID: nil,
		TranslationJobID:            nil,
		DictionaryLifecycle:         "permanent",
		DictionaryScope:             "global",
		DictionarySource:            "manual",
		SourceTerm:                  "Whiterun",
		TranslatedTerm:              "ホワイトラン",
		TermKind:                    "proper_noun",
		Reusable:                    true,
	}

	created, err := repo.CreateDictionaryEntry(ctx, draft)
	if err != nil {
		t.Fatalf("CreateDictionaryEntry failed: %v", err)
	}
	if created.ID == 0 {
		t.Fatal("expected non-zero ID after create")
	}

	got, err := repo.GetDictionaryEntryByID(ctx, created.ID)
	if err != nil {
		t.Fatalf("GetDictionaryEntryByID failed: %v", err)
	}
	if got.SourceTerm != "Whiterun" {
		t.Fatalf("expected SourceTerm=Whiterun, got %q", got.SourceTerm)
	}
	if got.TranslatedTerm != "ホワイトラン" {
		t.Fatalf("expected TranslatedTerm=ホワイトラン, got %q", got.TranslatedTerm)
	}
	if got.DictionaryScope != "global" {
		t.Fatalf("expected DictionaryScope=global, got %q", got.DictionaryScope)
	}
	if got.TranslationJobID != nil {
		t.Fatalf("expected nil TranslationJobID for global entry, got %v", got.TranslationJobID)
	}
}

// TestSCN_SMR_002_CreateJobLocalDictionaryEntry は FoundationDataRepository で
// job-local dictionary entry を保存・取得できることを検証する。
func TestSCN_SMR_002_CreateJobLocalDictionaryEntry(t *testing.T) {
	ctx := context.Background()
	db := openIntegrationDB(t)
	sourceRepo := repository.NewSQLiteTranslationSourceRepository(db)
	jobRepo := repository.NewSQLiteJobLifecycleRepository(db)
	foundRepo := repository.NewSQLiteFoundationDataRepository(db)

	xEdit, err := sourceRepo.CreateXEditExtractedData(ctx, repository.XEditExtractedDataDraft{
		SourceFilePath:   "Skyrim.esp",
		SourceTool:       "xEdit",
		TargetPluginName: "Skyrim.esm",
		TargetPluginType: "ESM",
		RecordCount:      100,
		ImportedAt:       fixedNow,
	})
	if err != nil {
		t.Fatalf("CreateXEditExtractedData failed: %v", err)
	}

	job, err := jobRepo.CreateTranslationJob(ctx, repository.TranslationJobDraft{
		XEditExtractedDataID: xEdit.ID,
		JobName:              "job-local-dict-test",
		State:                "pending",
		ProgressPercent:      0,
	})
	if err != nil {
		t.Fatalf("CreateTranslationJob failed: %v", err)
	}

	draft := repository.DictionaryEntryDraft{
		XTranslatorTranslationXMLID: nil,
		TranslationJobID:            &job.ID,
		DictionaryLifecycle:         "job",
		DictionaryScope:             "job_local",
		DictionarySource:            "ai",
		SourceTerm:                  "Jarl",
		TranslatedTerm:              "ヤール",
		TermKind:                    "title",
		Reusable:                    false,
	}

	created, err := foundRepo.CreateDictionaryEntry(ctx, draft)
	if err != nil {
		t.Fatalf("CreateDictionaryEntry (job-local) failed: %v", err)
	}

	got, err := foundRepo.GetDictionaryEntryByID(ctx, created.ID)
	if err != nil {
		t.Fatalf("GetDictionaryEntryByID failed: %v", err)
	}
	if got.TranslationJobID == nil {
		t.Fatal("expected non-nil TranslationJobID for job-local entry")
	}
	if *got.TranslationJobID != job.ID {
		t.Fatalf("expected TranslationJobID=%d, got %d", job.ID, *got.TranslationJobID)
	}
}

// TestSCN_SMR_002_DuplicatePersonaRejected は同一 npc_profile_id への
// PERSONA の二重保持が拒否されることを検証する。
func TestSCN_SMR_002_DuplicatePersonaRejected(t *testing.T) {
	ctx := context.Background()
	db := openIntegrationDB(t)
	sourceRepo := repository.NewSQLiteTranslationSourceRepository(db)
	foundRepo := repository.NewSQLiteFoundationDataRepository(db)

	profile, err := sourceRepo.UpsertNpcProfile(ctx, repository.NpcProfileDraft{
		TargetPluginName: "Skyrim.esm",
		FormID:           "000A2C8E",
		RecordType:       "NPC_",
		EditorID:         "Jarl",
		DisplayName:      "ヤール",
	})
	if err != nil {
		t.Fatalf("UpsertNpcProfile failed: %v", err)
	}

	firstDraft := repository.PersonaDraft{
		NpcProfileID:     profile.ID,
		TranslationJobID: nil,
		PersonaLifecycle: "permanent",
		PersonaScope:     "global",
		PersonaSource:    "manual",
	}
	_, err = foundRepo.CreatePersona(ctx, firstDraft)
	if err != nil {
		t.Fatalf("first CreatePersona failed: %v", err)
	}

	secondDraft := repository.PersonaDraft{
		NpcProfileID:     profile.ID,
		TranslationJobID: nil,
		PersonaLifecycle: "permanent",
		PersonaScope:     "global",
		PersonaSource:    "ai",
	}
	_, err = foundRepo.CreatePersona(ctx, secondDraft)
	if err == nil {
		t.Fatal("expected error on duplicate persona for same npc_profile_id, got nil")
	}
}

// ---------------------------------------------------------------------------
// SCN-SMR-003 Translation source persistence
// ---------------------------------------------------------------------------

// TestSCN_SMR_003_SaveTranslationSourceAllTables は TranslationSourceRepository で
// 全テーブル (X_EDIT_EXTRACTED_DATA, TRANSLATION_RECORD, NPC_PROFILE, NPC_RECORD,
// TRANSLATION_FIELD) への保存と取得が正常に動作することを検証する。
func TestSCN_SMR_003_SaveTranslationSourceAllTables(t *testing.T) {
	ctx := context.Background()
	db := openIntegrationDB(t)
	repo := repository.NewSQLiteTranslationSourceRepository(db)

	xEdit, err := repo.CreateXEditExtractedData(ctx, repository.XEditExtractedDataDraft{
		SourceFilePath:   "Dragonborn.esp",
		SourceTool:       "xEdit",
		TargetPluginName: "Dragonborn.esm",
		TargetPluginType: "ESM",
		RecordCount:      50,
		ImportedAt:       fixedNow,
	})
	if err != nil {
		t.Fatalf("CreateXEditExtractedData failed: %v", err)
	}
	if xEdit.SourceFilePath != "Dragonborn.esp" {
		t.Fatalf("expected SourceFilePath=Dragonborn.esp, got %q", xEdit.SourceFilePath)
	}

	rec, err := repo.CreateTranslationRecord(ctx, repository.TranslationRecordDraft{
		XEditExtractedDataID: xEdit.ID,
		FormID:               "020179A7",
		EditorID:             "NPC_Miraak",
		RecordType:           "NPC_",
	})
	if err != nil {
		t.Fatalf("CreateTranslationRecord failed: %v", err)
	}
	if rec.FormID != "020179A7" {
		t.Fatalf("expected FormID=020179A7, got %q", rec.FormID)
	}

	profile, err := repo.UpsertNpcProfile(ctx, repository.NpcProfileDraft{
		TargetPluginName: "Dragonborn.esm",
		FormID:           "020179A7",
		RecordType:       "NPC_",
		EditorID:         "NPC_Miraak",
		DisplayName:      "ミラーク",
	})
	if err != nil {
		t.Fatalf("UpsertNpcProfile failed: %v", err)
	}
	if profile.DisplayName != "ミラーク" {
		t.Fatalf("expected DisplayName=ミラーク, got %q", profile.DisplayName)
	}

	npcRec, err := repo.CreateNpcRecord(ctx, repository.NpcRecordDraft{
		TranslationRecordID: rec.ID,
		NpcProfileID:        profile.ID,
		VoiceType:           "MaleUniqueGhost",
	})
	if err != nil {
		t.Fatalf("CreateNpcRecord failed: %v", err)
	}
	if npcRec.NpcProfileID != profile.ID {
		t.Fatalf("expected NpcProfileID=%d, got %d", profile.ID, npcRec.NpcProfileID)
	}

	field, err := repo.CreateTranslationField(ctx, repository.TranslationFieldDraft{
		TranslationRecordID: rec.ID,
		SubrecordType:       "FULL",
		SourceText:          "Miraak",
		FieldOrder:          0,
	})
	if err != nil {
		t.Fatalf("CreateTranslationField failed: %v", err)
	}

	gotField, err := repo.GetTranslationFieldByID(ctx, field.ID)
	if err != nil {
		t.Fatalf("GetTranslationFieldByID failed: %v", err)
	}
	if gotField.SourceText != "Miraak" {
		t.Fatalf("expected SourceText=Miraak, got %q", gotField.SourceText)
	}

	fields, err := repo.ListTranslationFieldsByTranslationRecordID(ctx, rec.ID)
	if err != nil {
		t.Fatalf("ListTranslationFieldsByTranslationRecordID failed: %v", err)
	}
	if len(fields) != 1 {
		t.Fatalf("expected 1 field, got %d", len(fields))
	}
}

// TestSCN_SMR_003_ReopenPersistence は DB を reopen した後でも
// 同じ translation source を読み込めることを検証する。
func TestSCN_SMR_003_ReopenPersistence(t *testing.T) {
	ctx := context.Background()
	dbPath := filepath.Join(t.TempDir(), "reopen.sqlite3")

	var savedID int64

	// --- 第1オープン: データ保存 ---
	db1, err := sqlitedbinit.OpenMasterDictionaryDatabase(ctx, dbPath, nil)
	if err != nil {
		t.Fatalf("first open failed: %v", err)
	}
	repo1 := repository.NewSQLiteTranslationSourceRepository(db1)
	saved, err := repo1.CreateXEditExtractedData(ctx, repository.XEditExtractedDataDraft{
		SourceFilePath:   "skyrim.esp",
		SourceTool:       "xEdit",
		TargetPluginName: "Skyrim.esm",
		TargetPluginType: "ESM",
		RecordCount:      10,
		ImportedAt:       fixedNow,
	})
	if err != nil {
		_ = db1.Close()
		t.Fatalf("CreateXEditExtractedData on first open failed: %v", err)
	}
	savedID = saved.ID
	if closeErr := db1.Close(); closeErr != nil {
		t.Fatalf("close db1 failed: %v", closeErr)
	}

	// --- 第2オープン: データ読み込み ---
	db2, err := sqlitedbinit.OpenMasterDictionaryDatabase(ctx, dbPath, nil)
	if err != nil {
		t.Fatalf("reopen failed: %v", err)
	}
	t.Cleanup(func() { _ = db2.Close() })

	repo2 := repository.NewSQLiteTranslationSourceRepository(db2)
	got, err := repo2.GetXEditExtractedDataByID(ctx, savedID)
	if err != nil {
		t.Fatalf("GetXEditExtractedDataByID after reopen failed: %v", err)
	}
	if got.SourceFilePath != "skyrim.esp" {
		t.Fatalf("expected SourceFilePath=skyrim.esp, got %q", got.SourceFilePath)
	}
	if got.RecordCount != 10 {
		t.Fatalf("expected RecordCount=10, got %d", got.RecordCount)
	}
}

// ---------------------------------------------------------------------------
// SCN-SMR-004 Job lifecycle と output
// ---------------------------------------------------------------------------

// TestSCN_SMR_004_CreateJobAndPhaseRun は JobLifecycleRepository で
// TranslationJob と JobPhaseRun を作成し、一覧取得できることを検証する。
func TestSCN_SMR_004_CreateJobAndPhaseRun(t *testing.T) {
	ctx := context.Background()
	db := openIntegrationDB(t)
	sourceRepo := repository.NewSQLiteTranslationSourceRepository(db)
	jobRepo := repository.NewSQLiteJobLifecycleRepository(db)

	xEdit, err := sourceRepo.CreateXEditExtractedData(ctx, repository.XEditExtractedDataDraft{
		SourceFilePath:   "HearthFires.esp",
		SourceTool:       "xEdit",
		TargetPluginName: "HearthFires.esm",
		TargetPluginType: "ESM",
		RecordCount:      20,
		ImportedAt:       fixedNow,
	})
	if err != nil {
		t.Fatalf("CreateXEditExtractedData failed: %v", err)
	}

	job, err := jobRepo.CreateTranslationJob(ctx, repository.TranslationJobDraft{
		XEditExtractedDataID: xEdit.ID,
		JobName:              "phase-run-test-job",
		State:                "pending",
		ProgressPercent:      0,
	})
	if err != nil {
		t.Fatalf("CreateTranslationJob failed: %v", err)
	}
	if job.JobName != "phase-run-test-job" {
		t.Fatalf("expected JobName=phase-run-test-job, got %q", job.JobName)
	}

	phase, err := jobRepo.CreateJobPhaseRun(ctx, repository.JobPhaseRunDraft{
		TranslationJobID: job.ID,
		PhaseType:        "persona_generation",
		State:            "pending",
		ExecutionOrder:   1,
		AIProvider:       "openai",
		ModelName:        "gpt-4o",
		ExecutionMode:    "batch",
		CredentialRef:    "openai-key",
		InstructionKind:  "default",
	})
	if err != nil {
		t.Fatalf("CreateJobPhaseRun failed: %v", err)
	}
	if phase.TranslationJobID != job.ID {
		t.Fatalf("expected TranslationJobID=%d, got %d", job.ID, phase.TranslationJobID)
	}

	phases, err := jobRepo.ListJobPhaseRunsByJobID(ctx, job.ID)
	if err != nil {
		t.Fatalf("ListJobPhaseRunsByJobID failed: %v", err)
	}
	if len(phases) != 1 {
		t.Fatalf("expected 1 phase run, got %d", len(phases))
	}
	if phases[0].PhaseType != "persona_generation" {
		t.Fatalf("expected PhaseType=persona_generation, got %q", phases[0].PhaseType)
	}
}

// TestSCN_SMR_004_SaveJobTranslationField は JobOutputRepository で
// JobTranslationField を保存・取得できることを検証する。
func TestSCN_SMR_004_SaveJobTranslationField(t *testing.T) {
	ctx := context.Background()
	db := openIntegrationDB(t)
	sourceRepo := repository.NewSQLiteTranslationSourceRepository(db)
	jobRepo := repository.NewSQLiteJobLifecycleRepository(db)
	outputRepo := repository.NewSQLiteJobOutputRepository(db)

	xEdit, err := sourceRepo.CreateXEditExtractedData(ctx, repository.XEditExtractedDataDraft{
		SourceFilePath:   "Dawnguard.esp",
		SourceTool:       "xEdit",
		TargetPluginName: "Dawnguard.esm",
		TargetPluginType: "ESM",
		RecordCount:      30,
		ImportedAt:       fixedNow,
	})
	if err != nil {
		t.Fatalf("CreateXEditExtractedData failed: %v", err)
	}

	rec, err := sourceRepo.CreateTranslationRecord(ctx, repository.TranslationRecordDraft{
		XEditExtractedDataID: xEdit.ID,
		FormID:               "02003B79",
		EditorID:             "DLC1HunterBase",
		RecordType:           "NPC_",
	})
	if err != nil {
		t.Fatalf("CreateTranslationRecord failed: %v", err)
	}

	field, err := sourceRepo.CreateTranslationField(ctx, repository.TranslationFieldDraft{
		TranslationRecordID: rec.ID,
		SubrecordType:       "FULL",
		SourceText:          "Isran",
		FieldOrder:          0,
	})
	if err != nil {
		t.Fatalf("CreateTranslationField failed: %v", err)
	}

	job, err := jobRepo.CreateTranslationJob(ctx, repository.TranslationJobDraft{
		XEditExtractedDataID: xEdit.ID,
		JobName:              "output-field-test-job",
		State:                "pending",
		ProgressPercent:      0,
	})
	if err != nil {
		t.Fatalf("CreateTranslationJob failed: %v", err)
	}

	jtf, err := outputRepo.CreateJobTranslationField(ctx, repository.JobTranslationFieldDraft{
		TranslationJobID:   job.ID,
		TranslationFieldID: field.ID,
		AppliedPersonaID:   nil,
		TranslatedText:     "イスラン",
		OutputStatus:       "translated",
		RetryCount:         0,
	})
	if err != nil {
		t.Fatalf("CreateJobTranslationField failed: %v", err)
	}

	got, err := outputRepo.GetJobTranslationFieldByID(ctx, jtf.ID)
	if err != nil {
		t.Fatalf("GetJobTranslationFieldByID failed: %v", err)
	}
	if got.TranslatedText != "イスラン" {
		t.Fatalf("expected TranslatedText=イスラン, got %q", got.TranslatedText)
	}
	if got.TranslationJobID != job.ID {
		t.Fatalf("expected TranslationJobID=%d, got %d", job.ID, got.TranslationJobID)
	}

	list, err := outputRepo.ListJobTranslationFieldsByJobID(ctx, job.ID)
	if err != nil {
		t.Fatalf("ListJobTranslationFieldsByJobID failed: %v", err)
	}
	if len(list) != 1 {
		t.Fatalf("expected 1 job translation field, got %d", len(list))
	}
}

// TestSCN_SMR_004_JobSingleSourceFK は TRANSLATION_JOB が存在しない
// x_edit_extracted_data_id を参照しようとした場合に FK constraint error が
// 返ることを検証する (job は 1 translation source だけを参照する)。
func TestSCN_SMR_004_JobSingleSourceFK(t *testing.T) {
	ctx := context.Background()
	db := openIntegrationDB(t)
	jobRepo := repository.NewSQLiteJobLifecycleRepository(db)

	_, err := jobRepo.CreateTranslationJob(ctx, repository.TranslationJobDraft{
		XEditExtractedDataID: 99999, // 存在しない ID
		JobName:              "fk-test-job",
		State:                "pending",
		ProgressPercent:      0,
	})
	if err == nil {
		t.Fatal("expected FK constraint error for nonexistent x_edit_extracted_data_id, got nil")
	}
}

// ---------------------------------------------------------------------------
// SCN-SMR-005 Transaction rollback
// ---------------------------------------------------------------------------

// TestSCN_SMR_005_TransactionRollbackOnFKError は Transactor.WithTransaction で
// fn が FK error を返した場合に、それまでの insert が rollback されることを検証する。
func TestSCN_SMR_005_TransactionRollbackOnFKError(t *testing.T) {
	ctx := context.Background()
	db := openIntegrationDB(t)
	transactor := repository.NewSQLiteTransactor(db)
	sourceRepo := repository.NewSQLiteTranslationSourceRepository(db)
	jobRepo := repository.NewSQLiteJobLifecycleRepository(db)

	err := transactor.WithTransaction(ctx, func(txCtx context.Context) error {
		// 有効な insert: X_EDIT_EXTRACTED_DATA
		_, insertErr := sourceRepo.CreateXEditExtractedData(txCtx, repository.XEditExtractedDataDraft{
			SourceFilePath:   "rollback-test.esp",
			SourceTool:       "xEdit",
			TargetPluginName: "RollbackPlugin.esm",
			TargetPluginType: "ESM",
			RecordCount:      1,
			ImportedAt:       fixedNow,
		})
		if insertErr != nil {
			return fmt.Errorf("CreateXEditExtractedData: %w", insertErr)
		}

		// FK 違反の insert: 存在しない x_edit_extracted_data_id を参照
		_, insertErr = jobRepo.CreateTranslationJob(txCtx, repository.TranslationJobDraft{
			XEditExtractedDataID: 99999, // 存在しない ID → FK error
			JobName:              "should-rollback",
			State:                "pending",
			ProgressPercent:      0,
		})
		if insertErr != nil {
			return fmt.Errorf("CreateTranslationJob: %w", insertErr)
		}
		return nil // FK error を返して rollback をトリガー
	})

	if err == nil {
		t.Fatal("expected FK constraint error from WithTransaction, got nil")
	}

	// rollback 検証: X_EDIT_EXTRACTED_DATA に行が残っていないこと
	var count int
	if scanErr := db.QueryRowContext(
		ctx,
		"SELECT COUNT(*) FROM X_EDIT_EXTRACTED_DATA",
	).Scan(&count); scanErr != nil {
		t.Fatalf("count query after rollback failed: %v", scanErr)
	}
	if count != 0 {
		t.Fatalf("expected 0 rows in X_EDIT_EXTRACTED_DATA after rollback, got %d", count)
	}
}

// ---------------------------------------------------------------------------
// SCN-SMR-003 補強: TranslationFieldRecordReference ラウンドトリップ
// ---------------------------------------------------------------------------

// refRoundTripIDs は setupRefRoundTripDB1Phase が返す ID 群。
type refRoundTripIDs struct {
	fieldID  int64
	field2ID int64
	refRecID int64
}

// setupRefRoundTripDB1Phase は TestSCN_SMR_003_TranslationFieldRecordReferenceRoundTrip の
// 第1オープン保存フェーズを切り出す。認知複雑度を分散させるためのヘルパー。
func setupRefRoundTripDB1Phase(t *testing.T, ctx context.Context, dbPath string) refRoundTripIDs {
	t.Helper()
	db1, err := sqlitedbinit.OpenMasterDictionaryDatabase(ctx, dbPath, nil)
	if err != nil {
		t.Fatalf("first open failed: %v", err)
	}
	defer func() {
		if closeErr := db1.Close(); closeErr != nil {
			t.Fatalf("close db1 failed: %v", closeErr)
		}
	}()
	repo1 := repository.NewSQLiteTranslationSourceRepository(db1)

	xEdit, err := repo1.CreateXEditExtractedData(ctx, repository.XEditExtractedDataDraft{
		SourceFilePath:   "ref-test.esp",
		SourceTool:       "xEdit",
		TargetPluginName: "RefTest.esm",
		TargetPluginType: "ESM",
		RecordCount:      2,
		ImportedAt:       fixedNow,
	})
	if err != nil {
		t.Fatalf("CreateXEditExtractedData failed: %v", err)
	}
	rec1, err := repo1.CreateTranslationRecord(ctx, repository.TranslationRecordDraft{
		XEditExtractedDataID: xEdit.ID,
		FormID:               "000001AA",
		EditorID:             "SourceNPC",
		RecordType:           "NPC_",
	})
	if err != nil {
		t.Fatalf("CreateTranslationRecord (rec1) failed: %v", err)
	}
	rec2, err := repo1.CreateTranslationRecord(ctx, repository.TranslationRecordDraft{
		XEditExtractedDataID: xEdit.ID,
		FormID:               "000002BB",
		EditorID:             "ReferencedNPC",
		RecordType:           "NPC_",
	})
	if err != nil {
		t.Fatalf("CreateTranslationRecord (rec2) failed: %v", err)
	}
	field1, err := repo1.CreateTranslationField(ctx, repository.TranslationFieldDraft{
		TranslationRecordID: rec1.ID,
		SubrecordType:       "FULL",
		SourceText:          "Dragon",
		FieldOrder:          0,
	})
	if err != nil {
		t.Fatalf("CreateTranslationField (field1) failed: %v", err)
	}
	field2, err := repo1.CreateTranslationField(ctx, repository.TranslationFieldDraft{
		TranslationRecordID:        rec1.ID,
		SubrecordType:              "DESC",
		SourceText:                 "Dragon description",
		FieldOrder:                 1,
		PreviousTranslationFieldID: &field1.ID,
	})
	if err != nil {
		t.Fatalf("CreateTranslationField (field2) failed: %v", err)
	}
	_, err = repo1.CreateTranslationFieldRecordReference(ctx, repository.TranslationFieldRecordReferenceDraft{
		TranslationFieldID:            field1.ID,
		ReferencedTranslationRecordID: rec2.ID,
		ReferenceRole:                 "target",
	})
	if err != nil {
		t.Fatalf("CreateTranslationFieldRecordReference failed: %v", err)
	}
	return refRoundTripIDs{fieldID: field1.ID, field2ID: field2.ID, refRecID: rec2.ID}
}

// TestSCN_SMR_003_TranslationFieldRecordReferenceRoundTrip は
// ordered field (previous_translation_field_id) と別 record への reference を保存し、
// DB reopen 後に ListTranslationFieldRecordReferencesByFieldID で復元できることを検証する。
func TestSCN_SMR_003_TranslationFieldRecordReferenceRoundTrip(t *testing.T) {
	ctx := context.Background()
	dbPath := filepath.Join(t.TempDir(), "ref_roundtrip.sqlite3")

	// --- 第1オープン: データ保存 ---
	ids := setupRefRoundTripDB1Phase(t, ctx, dbPath)

	// --- 第2オープン: 読み込み確認 ---
	db2, err := sqlitedbinit.OpenMasterDictionaryDatabase(ctx, dbPath, nil)
	if err != nil {
		t.Fatalf("reopen failed: %v", err)
	}
	t.Cleanup(func() { _ = db2.Close() })

	repo2 := repository.NewSQLiteTranslationSourceRepository(db2)

	refs, err := repo2.ListTranslationFieldRecordReferencesByFieldID(ctx, ids.fieldID)
	if err != nil {
		t.Fatalf("ListTranslationFieldRecordReferencesByFieldID failed: %v", err)
	}
	if len(refs) != 1 {
		t.Fatalf("expected 1 record reference, got %d", len(refs))
	}
	if refs[0].ReferencedTranslationRecordID != ids.refRecID {
		t.Fatalf("expected ReferencedTranslationRecordID=%d, got %d", ids.refRecID, refs[0].ReferencedTranslationRecordID)
	}
	if refs[0].ReferenceRole != "target" {
		t.Fatalf("expected ReferenceRole=target, got %q", refs[0].ReferenceRole)
	}
	if refs[0].TranslationFieldID != ids.fieldID {
		t.Fatalf("expected TranslationFieldID=%d, got %d", ids.fieldID, refs[0].TranslationFieldID)
	}

	// ordered field の previous link 復元確認
	gotField2, err := repo2.GetTranslationFieldByID(ctx, ids.field2ID)
	if err != nil {
		t.Fatalf("GetTranslationFieldByID (field2) after reopen failed: %v", err)
	}
	if gotField2.PreviousTranslationFieldID == nil {
		t.Error("expected non-nil PreviousTranslationFieldID for field2 after reopen, got nil")
	} else if *gotField2.PreviousTranslationFieldID != ids.fieldID {
		t.Errorf("expected PreviousTranslationFieldID=%d, got %d", ids.fieldID, *gotField2.PreviousTranslationFieldID)
	}
}

// TestSCN_SMR_003_NpcProfileAndRecord は GetNpcProfileByID、
// GetNpcRecordByTranslationRecordID、ListTranslationRecordsByXEditID を検証する。
func TestSCN_SMR_003_NpcProfileAndRecord(t *testing.T) {
	ctx := context.Background()
	db := openIntegrationDB(t)
	repo := repository.NewSQLiteTranslationSourceRepository(db)

	xEdit, err := repo.CreateXEditExtractedData(ctx, repository.XEditExtractedDataDraft{
		SourceFilePath:   "npc-lookup.esp",
		SourceTool:       "xEdit",
		TargetPluginName: "NpcLookup.esm",
		TargetPluginType: "ESM",
		RecordCount:      2,
		ImportedAt:       fixedNow,
	})
	if err != nil {
		t.Fatalf("CreateXEditExtractedData failed: %v", err)
	}

	rec1, err := repo.CreateTranslationRecord(ctx, repository.TranslationRecordDraft{
		XEditExtractedDataID: xEdit.ID,
		FormID:               "000A1111",
		EditorID:             "NPC_Ulfric",
		RecordType:           "NPC_",
	})
	if err != nil {
		t.Fatalf("CreateTranslationRecord (rec1) failed: %v", err)
	}

	_, err = repo.CreateTranslationRecord(ctx, repository.TranslationRecordDraft{
		XEditExtractedDataID: xEdit.ID,
		FormID:               "000A2222",
		EditorID:             "NPC_Galmar",
		RecordType:           "NPC_",
	})
	if err != nil {
		t.Fatalf("CreateTranslationRecord (rec2) failed: %v", err)
	}

	profile, err := repo.UpsertNpcProfile(ctx, repository.NpcProfileDraft{
		TargetPluginName: "NpcLookup.esm",
		FormID:           "000A1111",
		RecordType:       "NPC_",
		EditorID:         "NPC_Ulfric",
		DisplayName:      "ウルフリック",
	})
	if err != nil {
		t.Fatalf("UpsertNpcProfile failed: %v", err)
	}

	_, err = repo.CreateNpcRecord(ctx, repository.NpcRecordDraft{
		TranslationRecordID: rec1.ID,
		NpcProfileID:        profile.ID,
		VoiceType:           "MaleUniqueUlfric",
	})
	if err != nil {
		t.Fatalf("CreateNpcRecord failed: %v", err)
	}

	// GetNpcProfileByID
	gotProfile, err := repo.GetNpcProfileByID(ctx, profile.ID)
	if err != nil {
		t.Fatalf("GetNpcProfileByID failed: %v", err)
	}
	if gotProfile.DisplayName != "ウルフリック" {
		t.Fatalf("expected DisplayName=ウルフリック, got %q", gotProfile.DisplayName)
	}
	if gotProfile.FormID != "000A1111" {
		t.Fatalf("expected FormID=000A1111, got %q", gotProfile.FormID)
	}

	// GetNpcRecordByTranslationRecordID
	gotNpcRec, err := repo.GetNpcRecordByTranslationRecordID(ctx, rec1.ID)
	if err != nil {
		t.Fatalf("GetNpcRecordByTranslationRecordID failed: %v", err)
	}
	if gotNpcRec.VoiceType != "MaleUniqueUlfric" {
		t.Fatalf("expected VoiceType=MaleUniqueUlfric, got %q", gotNpcRec.VoiceType)
	}
	if gotNpcRec.NpcProfileID != profile.ID {
		t.Fatalf("expected NpcProfileID=%d, got %d", profile.ID, gotNpcRec.NpcProfileID)
	}

	// ListTranslationRecordsByXEditID
	recs, err := repo.ListTranslationRecordsByXEditID(ctx, xEdit.ID)
	if err != nil {
		t.Fatalf("ListTranslationRecordsByXEditID failed: %v", err)
	}
	if len(recs) != 2 {
		t.Fatalf("expected 2 translation records, got %d", len(recs))
	}
}

// ---------------------------------------------------------------------------
// SCN-SMR-002 補強: Persona ライフサイクル・PersonaFieldEvidence・DictionaryEntry 更新/削除
// ---------------------------------------------------------------------------

// TestSCN_SMR_002_PersonaLifecycle は UpdatePersona、GetPersonaByNpcProfileID、
// CreatePersonaFieldEvidence、ListPersonaFieldEvidenceByPersonaID を検証する。
func TestSCN_SMR_002_PersonaLifecycle(t *testing.T) {
	ctx := context.Background()
	db := openIntegrationDB(t)
	sourceRepo := repository.NewSQLiteTranslationSourceRepository(db)
	foundRepo := repository.NewSQLiteFoundationDataRepository(db)

	xEdit, err := sourceRepo.CreateXEditExtractedData(ctx, repository.XEditExtractedDataDraft{
		SourceFilePath:   "persona-lifecycle.esp",
		SourceTool:       "xEdit",
		TargetPluginName: "PersonaLifecycle.esm",
		TargetPluginType: "ESM",
		RecordCount:      1,
		ImportedAt:       fixedNow,
	})
	if err != nil {
		t.Fatalf("CreateXEditExtractedData failed: %v", err)
	}

	rec, err := sourceRepo.CreateTranslationRecord(ctx, repository.TranslationRecordDraft{
		XEditExtractedDataID: xEdit.ID,
		FormID:               "000B1234",
		EditorID:             "NPC_Lydia",
		RecordType:           "NPC_",
	})
	if err != nil {
		t.Fatalf("CreateTranslationRecord failed: %v", err)
	}

	field, err := sourceRepo.CreateTranslationField(ctx, repository.TranslationFieldDraft{
		TranslationRecordID: rec.ID,
		SubrecordType:       "FULL",
		SourceText:          "Lydia",
		FieldOrder:          0,
	})
	if err != nil {
		t.Fatalf("CreateTranslationField failed: %v", err)
	}

	profile, err := sourceRepo.UpsertNpcProfile(ctx, repository.NpcProfileDraft{
		TargetPluginName: "PersonaLifecycle.esm",
		FormID:           "000B1234",
		RecordType:       "NPC_",
		EditorID:         "NPC_Lydia",
		DisplayName:      "リディア",
	})
	if err != nil {
		t.Fatalf("UpsertNpcProfile failed: %v", err)
	}

	persona, err := foundRepo.CreatePersona(ctx, repository.PersonaDraft{
		NpcProfileID:       profile.ID,
		PersonaLifecycle:   "permanent",
		PersonaScope:       "global",
		PersonaSource:      "manual",
		PersonaDescription: "忠実な従者",
		SpeechStyle:        "formal",
		PersonalitySummary: "誠実で勤勉",
	})
	if err != nil {
		t.Fatalf("CreatePersona failed: %v", err)
	}

	// UpdatePersona
	updated, err := foundRepo.UpdatePersona(ctx, persona.ID, repository.PersonaUpdateDraft{
		PersonaLifecycle:       "permanent",
		PersonaScope:           "global",
		PersonaSource:          "ai",
		PersonaDescription:     "更新された従者",
		SpeechStyle:            "formal",
		PersonalitySummary:     "更新されたサマリー",
		EvidenceUtteranceCount: 3,
	})
	if err != nil {
		t.Fatalf("UpdatePersona failed: %v", err)
	}
	if updated.PersonaDescription != "更新された従者" {
		t.Fatalf("expected PersonaDescription=更新された従者, got %q", updated.PersonaDescription)
	}
	if updated.EvidenceUtteranceCount != 3 {
		t.Fatalf("expected EvidenceUtteranceCount=3, got %d", updated.EvidenceUtteranceCount)
	}

	// GetPersonaByNpcProfileID
	gotPersona, err := foundRepo.GetPersonaByNpcProfileID(ctx, profile.ID)
	if err != nil {
		t.Fatalf("GetPersonaByNpcProfileID failed: %v", err)
	}
	if gotPersona.ID != persona.ID {
		t.Fatalf("expected persona ID=%d, got %d", persona.ID, gotPersona.ID)
	}

	// CreatePersonaFieldEvidence
	evidence, err := foundRepo.CreatePersonaFieldEvidence(ctx, repository.PersonaFieldEvidenceDraft{
		PersonaID:          persona.ID,
		TranslationFieldID: field.ID,
		EvidenceRole:       "voice_line",
	})
	if err != nil {
		t.Fatalf("CreatePersonaFieldEvidence failed: %v", err)
	}
	if evidence.EvidenceRole != "voice_line" {
		t.Fatalf("expected EvidenceRole=voice_line, got %q", evidence.EvidenceRole)
	}

	// ListPersonaFieldEvidenceByPersonaID
	evidences, err := foundRepo.ListPersonaFieldEvidenceByPersonaID(ctx, persona.ID)
	if err != nil {
		t.Fatalf("ListPersonaFieldEvidenceByPersonaID failed: %v", err)
	}
	if len(evidences) != 1 {
		t.Fatalf("expected 1 evidence, got %d", len(evidences))
	}
	if evidences[0].EvidenceRole != "voice_line" {
		t.Fatalf("expected EvidenceRole=voice_line, got %q", evidences[0].EvidenceRole)
	}
}

// TestSCN_SMR_002_DictionaryEntryUpdateDelete は UpdateDictionaryEntry と
// DeleteDictionaryEntry を検証する。
func TestSCN_SMR_002_DictionaryEntryUpdateDelete(t *testing.T) {
	ctx := context.Background()
	db := openIntegrationDB(t)
	foundRepo := repository.NewSQLiteFoundationDataRepository(db)

	entry, err := foundRepo.CreateDictionaryEntry(ctx, repository.DictionaryEntryDraft{
		DictionaryLifecycle: "permanent",
		DictionaryScope:     "global",
		DictionarySource:    "manual",
		SourceTerm:          "Nord",
		TranslatedTerm:      "ノルド",
		TermKind:            "race",
		Reusable:            true,
	})
	if err != nil {
		t.Fatalf("CreateDictionaryEntry failed: %v", err)
	}

	// UpdateDictionaryEntry
	updated, err := foundRepo.UpdateDictionaryEntry(ctx, entry.ID, repository.DictionaryEntryUpdateDraft{
		DictionaryLifecycle: "permanent",
		DictionaryScope:     "global",
		DictionarySource:    "manual",
		SourceTerm:          "Nord",
		TranslatedTerm:      "ノルド人",
		TermKind:            "race",
		Reusable:            true,
	})
	if err != nil {
		t.Fatalf("UpdateDictionaryEntry failed: %v", err)
	}
	if updated.TranslatedTerm != "ノルド人" {
		t.Fatalf("expected TranslatedTerm=ノルド人, got %q", updated.TranslatedTerm)
	}

	// DeleteDictionaryEntry
	if delErr := foundRepo.DeleteDictionaryEntry(ctx, entry.ID); delErr != nil {
		t.Fatalf("DeleteDictionaryEntry failed: %v", delErr)
	}

	// 削除後は ErrNotFound が返ること
	_, err = foundRepo.GetDictionaryEntryByID(ctx, entry.ID)
	if err == nil {
		t.Fatal("expected error after delete, got nil")
	}
}

// TestSCN_SMR_002_XTranslatorTranslationXML は CreateXTranslatorTranslationXML と
// GetXTranslatorTranslationXMLByID を検証する。
func TestSCN_SMR_002_XTranslatorTranslationXML(t *testing.T) {
	ctx := context.Background()
	db := openIntegrationDB(t)
	foundRepo := repository.NewSQLiteFoundationDataRepository(db)

	xml, err := foundRepo.CreateXTranslatorTranslationXML(ctx, repository.XTranslatorTranslationXMLDraft{
		FilePath:         "Skyrim_Dialogs.xml",
		TargetPluginName: "Skyrim.esm",
		TargetPluginType: "ESM",
		TermCount:        1500,
		ImportedAt:       fixedNow,
	})
	if err != nil {
		t.Fatalf("CreateXTranslatorTranslationXML failed: %v", err)
	}
	if xml.ID == 0 {
		t.Fatal("expected non-zero ID after create")
	}

	got, err := foundRepo.GetXTranslatorTranslationXMLByID(ctx, xml.ID)
	if err != nil {
		t.Fatalf("GetXTranslatorTranslationXMLByID failed: %v", err)
	}
	if got.FilePath != "Skyrim_Dialogs.xml" {
		t.Fatalf("expected FilePath=Skyrim_Dialogs.xml, got %q", got.FilePath)
	}
	if got.TermCount != 1500 {
		t.Fatalf("expected TermCount=1500, got %d", got.TermCount)
	}
	if got.TargetPluginName != "Skyrim.esm" {
		t.Fatalf("expected TargetPluginName=Skyrim.esm, got %q", got.TargetPluginName)
	}
}

// ---------------------------------------------------------------------------
// SCN-SMR-004 補強: PhaseRun associations と metadata 保持
// ---------------------------------------------------------------------------

// phaseRunFixture は arrangePhaseRunFixture が返す Arrange 結果。
type phaseRunFixture struct {
	phase   repository.JobPhaseRun
	jtf     repository.JobTranslationField
	persona repository.Persona
	entry   repository.DictionaryEntry
}

// arrangePhaseRunFixture は TestSCN_SMR_004_PhaseRunAssociations の Arrange フェーズを切り出す。
// 認知複雑度を分散させるためのヘルパー。
func arrangePhaseRunFixture(t *testing.T, ctx context.Context, db *sqlx.DB) phaseRunFixture {
	t.Helper()
	sourceRepo := repository.NewSQLiteTranslationSourceRepository(db)
	jobRepo := repository.NewSQLiteJobLifecycleRepository(db)
	outputRepo := repository.NewSQLiteJobOutputRepository(db)
	foundRepo := repository.NewSQLiteFoundationDataRepository(db)

	xEdit, err := sourceRepo.CreateXEditExtractedData(ctx, repository.XEditExtractedDataDraft{
		SourceFilePath:   "phase-assoc.esp",
		SourceTool:       "xEdit",
		TargetPluginName: "PhaseAssoc.esm",
		TargetPluginType: "ESM",
		RecordCount:      5,
		ImportedAt:       fixedNow,
	})
	if err != nil {
		t.Fatalf("CreateXEditExtractedData failed: %v", err)
	}
	rec, err := sourceRepo.CreateTranslationRecord(ctx, repository.TranslationRecordDraft{
		XEditExtractedDataID: xEdit.ID,
		FormID:               "01000ABC",
		EditorID:             "TestNPC",
		RecordType:           "NPC_",
	})
	if err != nil {
		t.Fatalf("CreateTranslationRecord failed: %v", err)
	}
	field, err := sourceRepo.CreateTranslationField(ctx, repository.TranslationFieldDraft{
		TranslationRecordID: rec.ID,
		SubrecordType:       "FULL",
		SourceText:          "Hello",
		FieldOrder:          0,
	})
	if err != nil {
		t.Fatalf("CreateTranslationField failed: %v", err)
	}
	job, err := jobRepo.CreateTranslationJob(ctx, repository.TranslationJobDraft{
		XEditExtractedDataID: xEdit.ID,
		JobName:              "phase-assoc-job",
		State:                "pending",
		ProgressPercent:      0,
	})
	if err != nil {
		t.Fatalf("CreateTranslationJob failed: %v", err)
	}
	jtf, err := outputRepo.CreateJobTranslationField(ctx, repository.JobTranslationFieldDraft{
		TranslationJobID:   job.ID,
		TranslationFieldID: field.ID,
		TranslatedText:     "こんにちは",
		OutputStatus:       "translated",
		RetryCount:         0,
	})
	if err != nil {
		t.Fatalf("CreateJobTranslationField failed: %v", err)
	}
	profile, err := sourceRepo.UpsertNpcProfile(ctx, repository.NpcProfileDraft{
		TargetPluginName: "PhaseAssoc.esm",
		FormID:           "01000ABC",
		RecordType:       "NPC_",
		EditorID:         "TestNPC",
		DisplayName:      "テストNPC",
	})
	if err != nil {
		t.Fatalf("UpsertNpcProfile failed: %v", err)
	}
	persona, err := foundRepo.CreatePersona(ctx, repository.PersonaDraft{
		NpcProfileID:     profile.ID,
		PersonaLifecycle: "job",
		PersonaScope:     "job_local",
		PersonaSource:    "ai",
	})
	if err != nil {
		t.Fatalf("CreatePersona failed: %v", err)
	}
	entry, err := foundRepo.CreateDictionaryEntry(ctx, repository.DictionaryEntryDraft{
		DictionaryLifecycle: "job",
		DictionaryScope:     "job_local",
		DictionarySource:    "ai",
		SourceTerm:          "Guard",
		TranslatedTerm:      "衛兵",
		TermKind:            "common_noun",
		Reusable:            true,
	})
	if err != nil {
		t.Fatalf("CreateDictionaryEntry failed: %v", err)
	}
	phase, err := jobRepo.CreateJobPhaseRun(ctx, repository.JobPhaseRunDraft{
		TranslationJobID: job.ID,
		PhaseType:        "translation",
		State:            "running",
		ExecutionOrder:   1,
		AIProvider:       "openai",
		ModelName:        "gpt-4o",
		ExecutionMode:    "batch",
		CredentialRef:    "openai-key",
		InstructionKind:  "default",
	})
	if err != nil {
		t.Fatalf("CreateJobPhaseRun failed: %v", err)
	}
	return phaseRunFixture{phase: phase, jtf: jtf, persona: persona, entry: entry}
}

// TestSCN_SMR_004_PhaseRunAssociations は CreatePhaseRunTranslationField、
// CreatePhaseRunPersona、CreatePhaseRunDictionaryEntry、UpdateJobPhaseRun の
// metadata (LatestExternalRunID, LatestError) 保持を検証する。
func TestSCN_SMR_004_PhaseRunAssociations(t *testing.T) {
	ctx := context.Background()
	db := openIntegrationDB(t)
	jobRepo := repository.NewSQLiteJobLifecycleRepository(db)

	// --- Arrange ---
	fix := arrangePhaseRunFixture(t, ctx, db)

	// --- Act ---
	prtf, err := jobRepo.CreatePhaseRunTranslationField(ctx, repository.PhaseRunTranslationFieldDraft{
		PhaseRunID:            fix.phase.ID,
		JobTranslationFieldID: fix.jtf.ID,
		Role:                  "target",
	})
	if err != nil {
		t.Fatalf("CreatePhaseRunTranslationField failed: %v", err)
	}

	prp, err := jobRepo.CreatePhaseRunPersona(ctx, repository.PhaseRunPersonaDraft{
		PhaseRunID: fix.phase.ID,
		PersonaID:  fix.persona.ID,
		Role:       "primary",
	})
	if err != nil {
		t.Fatalf("CreatePhaseRunPersona failed: %v", err)
	}

	prde, err := jobRepo.CreatePhaseRunDictionaryEntry(ctx, repository.PhaseRunDictionaryEntryDraft{
		PhaseRunID:        fix.phase.ID,
		DictionaryEntryID: fix.entry.ID,
		Role:              "context",
	})
	if err != nil {
		t.Fatalf("CreatePhaseRunDictionaryEntry failed: %v", err)
	}

	updated, err := jobRepo.UpdateJobPhaseRun(ctx, fix.phase.ID, repository.JobPhaseRunUpdateDraft{
		State:               "failed",
		ProgressPercent:     10,
		LatestExternalRunID: "batch-run-001",
		LatestError:         "AI timeout",
	})
	if err != nil {
		t.Fatalf("UpdateJobPhaseRun failed: %v", err)
	}

	// --- Assert ---
	if prtf.PhaseRunID != fix.phase.ID {
		t.Fatalf("expected PhaseRunID=%d, got %d", fix.phase.ID, prtf.PhaseRunID)
	}
	if prtf.JobTranslationFieldID != fix.jtf.ID {
		t.Fatalf("expected JobTranslationFieldID=%d, got %d", fix.jtf.ID, prtf.JobTranslationFieldID)
	}
	if prp.PersonaID != fix.persona.ID {
		t.Fatalf("expected PersonaID=%d, got %d", fix.persona.ID, prp.PersonaID)
	}
	if prde.DictionaryEntryID != fix.entry.ID {
		t.Fatalf("expected DictionaryEntryID=%d, got %d", fix.entry.ID, prde.DictionaryEntryID)
	}
	if updated.State != "failed" {
		t.Fatalf("expected State=failed, got %q", updated.State)
	}
	if updated.LatestError != "AI timeout" {
		t.Fatalf("expected LatestError=AI timeout, got %q", updated.LatestError)
	}
	if updated.LatestExternalRunID != "batch-run-001" {
		t.Fatalf("expected LatestExternalRunID=batch-run-001, got %q", updated.LatestExternalRunID)
	}

	// GetJobPhaseRunByID で永続化確認
	got, err := jobRepo.GetJobPhaseRunByID(ctx, fix.phase.ID)
	if err != nil {
		t.Fatalf("GetJobPhaseRunByID failed: %v", err)
	}
	if got.LatestError != "AI timeout" {
		t.Fatalf("GetJobPhaseRunByID: expected LatestError=AI timeout, got %q", got.LatestError)
	}
	if got.LatestExternalRunID != "batch-run-001" {
		t.Fatalf("GetJobPhaseRunByID: expected LatestExternalRunID=batch-run-001, got %q", got.LatestExternalRunID)
	}
}

// TestSCN_SMR_004_JobUpdateAndOutput は UpdateTranslationJob と
// UpdateJobTranslationField を検証する。
func TestSCN_SMR_004_JobUpdateAndOutput(t *testing.T) {
	ctx := context.Background()
	db := openIntegrationDB(t)
	sourceRepo := repository.NewSQLiteTranslationSourceRepository(db)
	jobRepo := repository.NewSQLiteJobLifecycleRepository(db)
	outputRepo := repository.NewSQLiteJobOutputRepository(db)

	xEdit, err := sourceRepo.CreateXEditExtractedData(ctx, repository.XEditExtractedDataDraft{
		SourceFilePath:   "job-update.esp",
		SourceTool:       "xEdit",
		TargetPluginName: "JobUpdate.esm",
		TargetPluginType: "ESM",
		RecordCount:      3,
		ImportedAt:       fixedNow,
	})
	if err != nil {
		t.Fatalf("CreateXEditExtractedData failed: %v", err)
	}

	rec, err := sourceRepo.CreateTranslationRecord(ctx, repository.TranslationRecordDraft{
		XEditExtractedDataID: xEdit.ID,
		FormID:               "000C5678",
		EditorID:             "NPC_Serana",
		RecordType:           "NPC_",
	})
	if err != nil {
		t.Fatalf("CreateTranslationRecord failed: %v", err)
	}

	field, err := sourceRepo.CreateTranslationField(ctx, repository.TranslationFieldDraft{
		TranslationRecordID: rec.ID,
		SubrecordType:       "FULL",
		SourceText:          "Serana",
		FieldOrder:          0,
	})
	if err != nil {
		t.Fatalf("CreateTranslationField failed: %v", err)
	}

	job, err := jobRepo.CreateTranslationJob(ctx, repository.TranslationJobDraft{
		XEditExtractedDataID: xEdit.ID,
		JobName:              "update-test-job",
		State:                "pending",
		ProgressPercent:      0,
	})
	if err != nil {
		t.Fatalf("CreateTranslationJob failed: %v", err)
	}

	// UpdateTranslationJob
	startedAt := fixedNow
	updatedJob, err := jobRepo.UpdateTranslationJob(ctx, job.ID, repository.TranslationJobUpdateDraft{
		JobName:         "update-test-job",
		State:           "running",
		ProgressPercent: 50,
		StartedAt:       &startedAt,
	})
	if err != nil {
		t.Fatalf("UpdateTranslationJob failed: %v", err)
	}
	if updatedJob.State != "running" {
		t.Fatalf("expected State=running, got %q", updatedJob.State)
	}
	if updatedJob.ProgressPercent != 50 {
		t.Fatalf("expected ProgressPercent=50, got %d", updatedJob.ProgressPercent)
	}

	jtf, err := outputRepo.CreateJobTranslationField(ctx, repository.JobTranslationFieldDraft{
		TranslationJobID:   job.ID,
		TranslationFieldID: field.ID,
		TranslatedText:     "セラナ",
		OutputStatus:       "translated",
		RetryCount:         0,
	})
	if err != nil {
		t.Fatalf("CreateJobTranslationField failed: %v", err)
	}

	// UpdateJobTranslationField
	updatedJtf, err := outputRepo.UpdateJobTranslationField(ctx, jtf.ID, repository.JobTranslationFieldUpdateDraft{
		TranslatedText: "セラナ (revised)",
		OutputStatus:   "revised",
		RetryCount:     1,
	})
	if err != nil {
		t.Fatalf("UpdateJobTranslationField failed: %v", err)
	}
	if updatedJtf.TranslatedText != "セラナ (revised)" {
		t.Fatalf("expected TranslatedText=セラナ (revised), got %q", updatedJtf.TranslatedText)
	}
	if updatedJtf.RetryCount != 1 {
		t.Fatalf("expected RetryCount=1, got %d", updatedJtf.RetryCount)
	}
	if updatedJtf.OutputStatus != "revised" {
		t.Fatalf("expected OutputStatus=revised, got %q", updatedJtf.OutputStatus)
	}
}
