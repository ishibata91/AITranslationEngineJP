package dbinit

import (
	"context"
	"path/filepath"
	"strings"
	"testing"
)

const (
	canonicalMigrationTestDatabaseFileName = "canonical-migration.sqlite3"
)

// TestCanonicalSchemaCreatesNPCProfileTable は 003_canonical_er_v1_tables.sql 適用後に
// NPC_PROFILE テーブルが存在することを検証する (SCN-SMR-001)。
func TestCanonicalSchemaCreatesNPCProfileTable(t *testing.T) {
	db := openMasterDictionaryDatabaseForTest(t, filepath.Join(t.TempDir(), "db", canonicalMigrationTestDatabaseFileName), nil)
	assertTableExists(t, db, "NPC_PROFILE")
}

// TestCanonicalSchemaCreatesTranslationRecordTable は TRANSLATION_RECORD テーブルの存在を検証する。
func TestCanonicalSchemaCreatesTranslationRecordTable(t *testing.T) {
	db := openMasterDictionaryDatabaseForTest(t, filepath.Join(t.TempDir(), "db", canonicalMigrationTestDatabaseFileName), nil)
	assertTableExists(t, db, "TRANSLATION_RECORD")
}

// TestCanonicalSchemaCreatesTranslationFieldTable は TRANSLATION_FIELD テーブルの存在を検証する。
func TestCanonicalSchemaCreatesTranslationFieldTable(t *testing.T) {
	db := openMasterDictionaryDatabaseForTest(t, filepath.Join(t.TempDir(), "db", canonicalMigrationTestDatabaseFileName), nil)
	assertTableExists(t, db, "TRANSLATION_FIELD")
}

// TestCanonicalSchemaCreatesTranslationJobTable は TRANSLATION_JOB テーブルの存在を検証する。
func TestCanonicalSchemaCreatesTranslationJobTable(t *testing.T) {
	db := openMasterDictionaryDatabaseForTest(t, filepath.Join(t.TempDir(), "db", canonicalMigrationTestDatabaseFileName), nil)
	assertTableExists(t, db, "TRANSLATION_JOB")
}

// TestCanonicalSchemaCreatesJobPhaseRunTable は JOB_PHASE_RUN テーブルの存在を検証する。
func TestCanonicalSchemaCreatesJobPhaseRunTable(t *testing.T) {
	db := openMasterDictionaryDatabaseForTest(t, filepath.Join(t.TempDir(), "db", canonicalMigrationTestDatabaseFileName), nil)
	assertTableExists(t, db, "JOB_PHASE_RUN")
}

// TestCanonicalSchemaCreatesPersonaTable は PERSONA テーブルの存在を検証する。
func TestCanonicalSchemaCreatesPersonaTable(t *testing.T) {
	db := openMasterDictionaryDatabaseForTest(t, filepath.Join(t.TempDir(), "db", canonicalMigrationTestDatabaseFileName), nil)
	assertTableExists(t, db, "PERSONA")
}

// TestCanonicalSchemaCreatesDictionaryEntryTable は DICTIONARY_ENTRY テーブルの存在を検証する。
func TestCanonicalSchemaCreatesDictionaryEntryTable(t *testing.T) {
	db := openMasterDictionaryDatabaseForTest(t, filepath.Join(t.TempDir(), "db", canonicalMigrationTestDatabaseFileName), nil)
	assertTableExists(t, db, "DICTIONARY_ENTRY")
}

// TestCanonicalSchemaNPCProfileUniqueConstraintRejectsDuplicate は
// NPC_PROFILE の (target_plugin_name, form_id, record_type) UNIQUE 制約が機能することを検証する。
func TestCanonicalSchemaNPCProfileUniqueConstraintRejectsDuplicate(t *testing.T) {
	db := openMasterDictionaryDatabaseForTest(t, filepath.Join(t.TempDir(), "db", canonicalMigrationTestDatabaseFileName), nil)

	const insertNPCProfile = `
INSERT INTO NPC_PROFILE (target_plugin_name, form_id, record_type, editor_id, display_name, created_at, updated_at)
VALUES (?, ?, ?, ?, ?, ?, ?)`

	_, err := db.ExecContext(context.Background(), insertNPCProfile,
		"Skyrim.esm", "000001FF", "NPC_",
		"EditorNPC01", "Test NPC",
		"2026-01-01T00:00:00Z", "2026-01-01T00:00:00Z",
	)
	if err != nil {
		t.Fatalf("expected first NPC_PROFILE insert to succeed: %v", err)
	}

	_, err = db.ExecContext(context.Background(), insertNPCProfile,
		"Skyrim.esm", "000001FF", "NPC_",
		"EditorNPC01Dup", "Duplicate NPC",
		"2026-01-01T00:00:00Z", "2026-01-01T00:00:00Z",
	)
	if err == nil {
		t.Fatal("expected duplicate NPC_PROFILE insert to fail due to UNIQUE constraint")
	}
	if !strings.Contains(err.Error(), "UNIQUE constraint failed") {
		t.Fatalf("expected UNIQUE constraint error, got: %v", err)
	}
}

// TestCanonicalSchemaNPCRecordForeignKeyRejectsInvalidNPCProfileID は
// NPC_RECORD の npc_profile_id FK 制約が、存在しない NPC_PROFILE_ID を拒否することを検証する。
func TestCanonicalSchemaNPCRecordForeignKeyRejectsInvalidNPCProfileID(t *testing.T) {
	db := openMasterDictionaryDatabaseForTest(t, filepath.Join(t.TempDir(), "db", canonicalMigrationTestDatabaseFileName), nil)

	_, err := db.ExecContext(context.Background(),
		`INSERT INTO X_EDIT_EXTRACTED_DATA (source_file_path, source_tool, target_plugin_name, target_plugin_type, record_count, imported_at)
         VALUES (?, ?, ?, ?, ?, ?)`,
		"/tmp/data.json", "xEdit", "Skyrim.esm", "esm", 1, "2026-01-01T00:00:00Z",
	)
	if err != nil {
		t.Fatalf("expected X_EDIT_EXTRACTED_DATA insert to succeed: %v", err)
	}

	_, err = db.ExecContext(context.Background(),
		`INSERT INTO TRANSLATION_RECORD (x_edit_extracted_data_id, form_id, editor_id, record_type)
         VALUES (?, ?, ?, ?)`,
		1, "000001FF", "EditorNPC01", "NPC_",
	)
	if err != nil {
		t.Fatalf("expected TRANSLATION_RECORD insert to succeed: %v", err)
	}

	_, err = db.ExecContext(context.Background(),
		`INSERT INTO NPC_RECORD (translation_record_id, npc_profile_id, race, sex, npc_class, voice_type)
         VALUES (?, ?, ?, ?, ?, ?)`,
		1, 99999, "Nord", "Male", "CombatWarrior1H", "",
	)
	if err == nil {
		t.Fatal("expected NPC_RECORD insert with invalid npc_profile_id to fail due to FOREIGN KEY constraint")
	}
	if !strings.Contains(err.Error(), "FOREIGN KEY constraint failed") {
		t.Fatalf("expected FOREIGN KEY constraint error, got: %v", err)
	}
}

// TestCanonicalSchemaPersonaUniqueConstraintRejectsDuplicate は
// PERSONA の npc_profile_id UNIQUE 制約が機能することを検証する。
func TestCanonicalSchemaPersonaUniqueConstraintRejectsDuplicate(t *testing.T) {
	db := openMasterDictionaryDatabaseForTest(t, filepath.Join(t.TempDir(), "db", canonicalMigrationTestDatabaseFileName), nil)

	_, err := db.ExecContext(context.Background(),
		`INSERT INTO NPC_PROFILE (target_plugin_name, form_id, record_type, created_at, updated_at)
         VALUES (?, ?, ?, ?, ?)`,
		"Skyrim.esm", "000002FF", "NPC_",
		"2026-01-01T00:00:00Z", "2026-01-01T00:00:00Z",
	)
	if err != nil {
		t.Fatalf("expected NPC_PROFILE insert to succeed: %v", err)
	}

	_, err = db.ExecContext(context.Background(),
		`INSERT INTO PERSONA (npc_profile_id, persona_lifecycle, created_at, updated_at)
         VALUES (?, ?, ?, ?)`,
		1, "active", "2026-01-01T00:00:00Z", "2026-01-01T00:00:00Z",
	)
	if err != nil {
		t.Fatalf("expected first PERSONA insert to succeed: %v", err)
	}

	_, err = db.ExecContext(context.Background(),
		`INSERT INTO PERSONA (npc_profile_id, persona_lifecycle, created_at, updated_at)
         VALUES (?, ?, ?, ?)`,
		1, "active", "2026-01-01T00:00:00Z", "2026-01-01T00:00:00Z",
	)
	if err == nil {
		t.Fatal("expected duplicate PERSONA insert to fail due to UNIQUE constraint on npc_profile_id")
	}
	if !strings.Contains(err.Error(), "UNIQUE constraint failed") {
		t.Fatalf("expected UNIQUE constraint error, got: %v", err)
	}
}

// TestCanonicalSchemaTranslationArtifactUniqueConstraintRejectsDuplicate は
// TRANSLATION_ARTIFACT の translation_job_id UNIQUE 制約が機能することを検証する。
func TestCanonicalSchemaTranslationArtifactUniqueConstraintRejectsDuplicate(t *testing.T) {
	db := openMasterDictionaryDatabaseForTest(t, filepath.Join(t.TempDir(), "db", canonicalMigrationTestDatabaseFileName), nil)

	_, err := db.ExecContext(context.Background(),
		`INSERT INTO X_EDIT_EXTRACTED_DATA (source_file_path, source_tool, target_plugin_name, target_plugin_type, record_count, imported_at)
         VALUES (?, ?, ?, ?, ?, ?)`,
		"/tmp/artifact.json", "xEdit", "Skyrim.esm", "esm", 0, "2026-01-01T00:00:00Z",
	)
	if err != nil {
		t.Fatalf("expected X_EDIT_EXTRACTED_DATA insert to succeed: %v", err)
	}

	_, err = db.ExecContext(context.Background(),
		`INSERT INTO TRANSLATION_JOB (x_edit_extracted_data_id, state, created_at)
         VALUES (?, ?, ?)`,
		1, "pending", "2026-01-01T00:00:00Z",
	)
	if err != nil {
		t.Fatalf("expected TRANSLATION_JOB insert to succeed: %v", err)
	}

	_, err = db.ExecContext(context.Background(),
		`INSERT INTO TRANSLATION_ARTIFACT (translation_job_id, artifact_format)
         VALUES (?, ?)`,
		1, "xTranslator",
	)
	if err != nil {
		t.Fatalf("expected first TRANSLATION_ARTIFACT insert to succeed: %v", err)
	}

	_, err = db.ExecContext(context.Background(),
		`INSERT INTO TRANSLATION_ARTIFACT (translation_job_id, artifact_format)
         VALUES (?, ?)`,
		1, "xTranslator",
	)
	if err == nil {
		t.Fatal("expected duplicate TRANSLATION_ARTIFACT insert to fail due to UNIQUE constraint on translation_job_id")
	}
	if !strings.Contains(err.Error(), "UNIQUE constraint failed") {
		t.Fatalf("expected UNIQUE constraint error, got: %v", err)
	}
}

// TestCanonicalSchemaXtranslatorOutputRowUniqueConstraintRejectsDuplicate は
// XTRANSLATOR_OUTPUT_ROW の job_translation_field_id UNIQUE 制約が機能することを検証する。
func TestCanonicalSchemaXtranslatorOutputRowUniqueConstraintRejectsDuplicate(t *testing.T) {
	db := openMasterDictionaryDatabaseForTest(t, filepath.Join(t.TempDir(), "db", canonicalMigrationTestDatabaseFileName), nil)

	_, err := db.ExecContext(context.Background(),
		`INSERT INTO X_EDIT_EXTRACTED_DATA (source_file_path, source_tool, target_plugin_name, target_plugin_type, record_count, imported_at)
         VALUES (?, ?, ?, ?, ?, ?)`,
		"/tmp/output.json", "xEdit", "Skyrim.esm", "esm", 1, "2026-01-01T00:00:00Z",
	)
	if err != nil {
		t.Fatalf("expected X_EDIT_EXTRACTED_DATA insert to succeed: %v", err)
	}

	_, err = db.ExecContext(context.Background(),
		`INSERT INTO TRANSLATION_RECORD (x_edit_extracted_data_id, form_id, record_type)
         VALUES (?, ?, ?)`,
		1, "000003FF", "NPC_",
	)
	if err != nil {
		t.Fatalf("expected TRANSLATION_RECORD insert to succeed: %v", err)
	}

	_, err = db.ExecContext(context.Background(),
		`INSERT INTO TRANSLATION_FIELD (translation_record_id, subrecord_type)
         VALUES (?, ?)`,
		1, "FULL",
	)
	if err != nil {
		t.Fatalf("expected TRANSLATION_FIELD insert to succeed: %v", err)
	}

	_, err = db.ExecContext(context.Background(),
		`INSERT INTO TRANSLATION_JOB (x_edit_extracted_data_id, state, created_at)
         VALUES (?, ?, ?)`,
		1, "pending", "2026-01-01T00:00:00Z",
	)
	if err != nil {
		t.Fatalf("expected TRANSLATION_JOB insert to succeed: %v", err)
	}

	_, err = db.ExecContext(context.Background(),
		`INSERT INTO JOB_TRANSLATION_FIELD (translation_job_id, translation_field_id, updated_at)
         VALUES (?, ?, ?)`,
		1, 1, "2026-01-01T00:00:00Z",
	)
	if err != nil {
		t.Fatalf("expected JOB_TRANSLATION_FIELD insert to succeed: %v", err)
	}

	_, err = db.ExecContext(context.Background(),
		`INSERT INTO TRANSLATION_ARTIFACT (translation_job_id, artifact_format)
         VALUES (?, ?)`,
		1, "xTranslator",
	)
	if err != nil {
		t.Fatalf("expected TRANSLATION_ARTIFACT insert to succeed: %v", err)
	}

	_, err = db.ExecContext(context.Background(),
		`INSERT INTO XTRANSLATOR_OUTPUT_ROW (translation_artifact_id, job_translation_field_id)
         VALUES (?, ?)`,
		1, 1,
	)
	if err != nil {
		t.Fatalf("expected first XTRANSLATOR_OUTPUT_ROW insert to succeed: %v", err)
	}

	_, err = db.ExecContext(context.Background(),
		`INSERT INTO XTRANSLATOR_OUTPUT_ROW (translation_artifact_id, job_translation_field_id)
         VALUES (?, ?)`,
		1, 1,
	)
	if err == nil {
		t.Fatal("expected duplicate XTRANSLATOR_OUTPUT_ROW insert to fail due to UNIQUE constraint on job_translation_field_id")
	}
	if !strings.Contains(err.Error(), "UNIQUE constraint failed") {
		t.Fatalf("expected UNIQUE constraint error, got: %v", err)
	}
}

// TestCanonicalSchemaIndexesExist は 003_canonical_er_v1_tables.sql で定義された
// index が sqlite_master に存在することを検証する。
func TestCanonicalSchemaIndexesExist(t *testing.T) {
	db := openMasterDictionaryDatabaseForTest(t, filepath.Join(t.TempDir(), "db", canonicalMigrationTestDatabaseFileName), nil)

	expectedIndexes := []string{
		"idx_translation_record_x_edit",
		"idx_npc_record_npc_profile",
		"idx_translation_field_record",
		"idx_tf_record_reference_field",
		"idx_translation_job_x_edit",
		"idx_persona_translation_job",
		"idx_persona_field_evidence_persona",
		"idx_dictionary_entry_translation_job",
		"idx_job_translation_field_job",
		"idx_job_translation_field_field",
		"idx_job_phase_run_job",
		"idx_phase_run_translation_field_run",
		"idx_phase_run_persona_run",
		"idx_phase_run_dictionary_entry_run",
		"idx_xtranslator_output_row_artifact",
	}

	for _, indexName := range expectedIndexes {
		var count int
		queryErr := db.QueryRowContext(
			context.Background(),
			"SELECT COUNT(*) FROM sqlite_master WHERE type='index' AND name = ?",
			indexName,
		).Scan(&count)
		if queryErr != nil {
			t.Fatalf("expected index existence query to succeed for %q: %v", indexName, queryErr)
		}
		if count != 1 {
			t.Errorf("expected index %q to exist in sqlite_master", indexName)
		}
	}
}

// TestOpenMasterDictionaryDatabaseFailsWithEmptyPath は空パスを渡した場合に
// エラーが返ることを検証する (ensureDatabasePath の空チェックパスをカバー)。
func TestOpenMasterDictionaryDatabaseFailsWithEmptyPath(t *testing.T) {
	_, err := OpenMasterDictionaryDatabase(context.Background(), "", nil)
	if err == nil {
		t.Fatal("expected error for empty database path")
	}
	if !strings.Contains(err.Error(), "database path is required") {
		t.Fatalf("expected 'database path is required' error, got: %v", err)
	}
}

// TestOpenMasterDictionaryDatabaseFailsWhenDirectoryCannotBeCreated は
// 作成不可能なディレクトリパスを渡した場合にエラーが返ることを検証する
// (ensureDatabasePath の os.MkdirAll エラーパスをカバー)。
func TestOpenMasterDictionaryDatabaseFailsWhenDirectoryCannotBeCreated(t *testing.T) {
	_, err := OpenMasterDictionaryDatabase(context.Background(), "/dev/null/subdir/test.db", nil)
	if err == nil {
		t.Fatal("expected error when database parent directory cannot be created")
	}
}

// TestOpenMasterDictionaryDatabaseFailsWithCancelledContext は
// キャンセル済みコンテキストを渡した場合に PingContext エラーが返ることを検証する。
func TestOpenMasterDictionaryDatabaseFailsWithCancelledContext(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	dbPath := filepath.Join(t.TempDir(), "db", canonicalMigrationTestDatabaseFileName)
	_, err := OpenMasterDictionaryDatabase(ctx, dbPath, nil)
	if err == nil {
		t.Fatal("expected error when context is already cancelled")
	}
}

// TestCanonicalSchemaNPCProfileRequiresFormId は NPC_PROFILE の form_id NOT NULL 制約が
// 機能することを検証する。
func TestCanonicalSchemaNPCProfileRequiresFormId(t *testing.T) {
	db := openMasterDictionaryDatabaseForTest(t, filepath.Join(t.TempDir(), "db", canonicalMigrationTestDatabaseFileName), nil)

	_, err := db.ExecContext(context.Background(),
		`INSERT INTO NPC_PROFILE (target_plugin_name, form_id, record_type, created_at, updated_at)
         VALUES (?, NULL, ?, ?, ?)`,
		"Skyrim.esm", "NPC_",
		"2026-01-01T00:00:00Z", "2026-01-01T00:00:00Z",
	)
	if err == nil {
		t.Fatal("expected NPC_PROFILE insert with NULL form_id to fail due to NOT NULL constraint")
	}
	if !strings.Contains(err.Error(), "NOT NULL constraint failed") {
		t.Fatalf("expected NOT NULL constraint error, got: %v", err)
	}
}

// TestCanonicalSchemaTFRecordRefFKRejectsInvalidFieldId は
// TRANSLATION_FIELD_RECORD_REFERENCE の translation_field_id FK が
// 存在しない ID を拒否することを検証する。
// referenced_translation_record_id は有効な値を使い、translation_field_id だけを無効にする。
func TestCanonicalSchemaTFRecordRefFKRejectsInvalidFieldId(t *testing.T) {
	db := openMasterDictionaryDatabaseForTest(t, filepath.Join(t.TempDir(), "db", canonicalMigrationTestDatabaseFileName), nil)

	_, err := db.ExecContext(context.Background(),
		`INSERT INTO X_EDIT_EXTRACTED_DATA (source_file_path, source_tool, target_plugin_name, target_plugin_type, record_count, imported_at)
         VALUES (?, ?, ?, ?, ?, ?)`,
		"/tmp/tfref1.json", "xEdit", "Skyrim.esm", "esm", 1, "2026-01-01T00:00:00Z",
	)
	if err != nil {
		t.Fatalf("expected X_EDIT_EXTRACTED_DATA insert to succeed: %v", err)
	}

	_, err = db.ExecContext(context.Background(),
		`INSERT INTO TRANSLATION_RECORD (x_edit_extracted_data_id, form_id, record_type)
         VALUES (?, ?, ?)`,
		1, "000010FF", "NPC_",
	)
	if err != nil {
		t.Fatalf("expected TRANSLATION_RECORD insert to succeed: %v", err)
	}

	_, err = db.ExecContext(context.Background(),
		`INSERT INTO TRANSLATION_FIELD_RECORD_REFERENCE (translation_field_id, referenced_translation_record_id, reference_role)
         VALUES (?, ?, ?)`,
		9999, 1, "",
	)
	if err == nil {
		t.Fatal("expected TRANSLATION_FIELD_RECORD_REFERENCE insert with invalid translation_field_id to fail due to FOREIGN KEY constraint")
	}
	if !strings.Contains(err.Error(), "FOREIGN KEY constraint failed") {
		t.Fatalf("expected FOREIGN KEY constraint error, got: %v", err)
	}
}

// TestCanonicalSchemaTFRecordRefFKRejectsInvalidRecordRef は
// TRANSLATION_FIELD_RECORD_REFERENCE の referenced_translation_record_id FK が
// 存在しない ID を拒否することを検証する。
// translation_field_id は有効な値を使い、referenced_translation_record_id だけを無効にする。
func TestCanonicalSchemaTFRecordRefFKRejectsInvalidRecordRef(t *testing.T) {
	db := openMasterDictionaryDatabaseForTest(t, filepath.Join(t.TempDir(), "db", canonicalMigrationTestDatabaseFileName), nil)

	_, err := db.ExecContext(context.Background(),
		`INSERT INTO X_EDIT_EXTRACTED_DATA (source_file_path, source_tool, target_plugin_name, target_plugin_type, record_count, imported_at)
         VALUES (?, ?, ?, ?, ?, ?)`,
		"/tmp/tfref2.json", "xEdit", "Skyrim.esm", "esm", 1, "2026-01-01T00:00:00Z",
	)
	if err != nil {
		t.Fatalf("expected X_EDIT_EXTRACTED_DATA insert to succeed: %v", err)
	}

	_, err = db.ExecContext(context.Background(),
		`INSERT INTO TRANSLATION_RECORD (x_edit_extracted_data_id, form_id, record_type)
         VALUES (?, ?, ?)`,
		1, "000011FF", "NPC_",
	)
	if err != nil {
		t.Fatalf("expected TRANSLATION_RECORD insert to succeed: %v", err)
	}

	_, err = db.ExecContext(context.Background(),
		`INSERT INTO TRANSLATION_FIELD (translation_record_id, subrecord_type)
         VALUES (?, ?)`,
		1, "FULL",
	)
	if err != nil {
		t.Fatalf("expected TRANSLATION_FIELD insert to succeed: %v", err)
	}

	_, err = db.ExecContext(context.Background(),
		`INSERT INTO TRANSLATION_FIELD_RECORD_REFERENCE (translation_field_id, referenced_translation_record_id, reference_role)
         VALUES (?, ?, ?)`,
		1, 9999, "",
	)
	if err == nil {
		t.Fatal("expected TRANSLATION_FIELD_RECORD_REFERENCE insert with invalid referenced_translation_record_id to fail due to FOREIGN KEY constraint")
	}
	if !strings.Contains(err.Error(), "FOREIGN KEY constraint failed") {
		t.Fatalf("expected FOREIGN KEY constraint error, got: %v", err)
	}
}

// TestCanonicalSchemaPFEvidenceFKRejectsInvalidPersonaId は
// PERSONA_FIELD_EVIDENCE の persona_id FK が存在しない ID を拒否することを検証する。
// translation_field_id は有効な値を使い、persona_id だけを無効にする。
func TestCanonicalSchemaPFEvidenceFKRejectsInvalidPersonaId(t *testing.T) {
	db := openMasterDictionaryDatabaseForTest(t, filepath.Join(t.TempDir(), "db", canonicalMigrationTestDatabaseFileName), nil)

	_, err := db.ExecContext(context.Background(),
		`INSERT INTO X_EDIT_EXTRACTED_DATA (source_file_path, source_tool, target_plugin_name, target_plugin_type, record_count, imported_at)
         VALUES (?, ?, ?, ?, ?, ?)`,
		"/tmp/pfe1.json", "xEdit", "Skyrim.esm", "esm", 1, "2026-01-01T00:00:00Z",
	)
	if err != nil {
		t.Fatalf("expected X_EDIT_EXTRACTED_DATA insert to succeed: %v", err)
	}

	_, err = db.ExecContext(context.Background(),
		`INSERT INTO TRANSLATION_RECORD (x_edit_extracted_data_id, form_id, record_type)
         VALUES (?, ?, ?)`,
		1, "000020FF", "NPC_",
	)
	if err != nil {
		t.Fatalf("expected TRANSLATION_RECORD insert to succeed: %v", err)
	}

	_, err = db.ExecContext(context.Background(),
		`INSERT INTO TRANSLATION_FIELD (translation_record_id, subrecord_type)
         VALUES (?, ?)`,
		1, "FULL",
	)
	if err != nil {
		t.Fatalf("expected TRANSLATION_FIELD insert to succeed: %v", err)
	}

	_, err = db.ExecContext(context.Background(),
		`INSERT INTO PERSONA_FIELD_EVIDENCE (persona_id, translation_field_id, evidence_role)
         VALUES (?, ?, ?)`,
		9999, 1, "",
	)
	if err == nil {
		t.Fatal("expected PERSONA_FIELD_EVIDENCE insert with invalid persona_id to fail due to FOREIGN KEY constraint")
	}
	if !strings.Contains(err.Error(), "FOREIGN KEY constraint failed") {
		t.Fatalf("expected FOREIGN KEY constraint error, got: %v", err)
	}
}

// TestCanonicalSchemaPFEvidenceFKRejectsInvalidFieldId は
// PERSONA_FIELD_EVIDENCE の translation_field_id FK が存在しない ID を拒否することを検証する。
// persona_id は有効な値を使い、translation_field_id だけを無効にする。
func TestCanonicalSchemaPFEvidenceFKRejectsInvalidFieldId(t *testing.T) {
	db := openMasterDictionaryDatabaseForTest(t, filepath.Join(t.TempDir(), "db", canonicalMigrationTestDatabaseFileName), nil)

	_, err := db.ExecContext(context.Background(),
		`INSERT INTO NPC_PROFILE (target_plugin_name, form_id, record_type, created_at, updated_at)
         VALUES (?, ?, ?, ?, ?)`,
		"Skyrim.esm", "000021FF", "NPC_",
		"2026-01-01T00:00:00Z", "2026-01-01T00:00:00Z",
	)
	if err != nil {
		t.Fatalf("expected NPC_PROFILE insert to succeed: %v", err)
	}

	_, err = db.ExecContext(context.Background(),
		`INSERT INTO PERSONA (npc_profile_id, persona_lifecycle, created_at, updated_at)
         VALUES (?, ?, ?, ?)`,
		1, "active", "2026-01-01T00:00:00Z", "2026-01-01T00:00:00Z",
	)
	if err != nil {
		t.Fatalf("expected PERSONA insert to succeed: %v", err)
	}

	_, err = db.ExecContext(context.Background(),
		`INSERT INTO PERSONA_FIELD_EVIDENCE (persona_id, translation_field_id, evidence_role)
         VALUES (?, ?, ?)`,
		1, 9999, "",
	)
	if err == nil {
		t.Fatal("expected PERSONA_FIELD_EVIDENCE insert with invalid translation_field_id to fail due to FOREIGN KEY constraint")
	}
	if !strings.Contains(err.Error(), "FOREIGN KEY constraint failed") {
		t.Fatalf("expected FOREIGN KEY constraint error, got: %v", err)
	}
}

// TestCanonicalSchemaXOutputRowFKRejectsInvalidArtifactId は
// XTRANSLATOR_OUTPUT_ROW の translation_artifact_id FK が存在しない ID を拒否することを検証する。
// job_translation_field_id は有効な値を使い、translation_artifact_id だけを無効にする。
func TestCanonicalSchemaXOutputRowFKRejectsInvalidArtifactId(t *testing.T) {
	db := openMasterDictionaryDatabaseForTest(t, filepath.Join(t.TempDir(), "db", canonicalMigrationTestDatabaseFileName), nil)

	_, err := db.ExecContext(context.Background(),
		`INSERT INTO X_EDIT_EXTRACTED_DATA (source_file_path, source_tool, target_plugin_name, target_plugin_type, record_count, imported_at)
         VALUES (?, ?, ?, ?, ?, ?)`,
		"/tmp/xout1.json", "xEdit", "Skyrim.esm", "esm", 1, "2026-01-01T00:00:00Z",
	)
	if err != nil {
		t.Fatalf("expected X_EDIT_EXTRACTED_DATA insert to succeed: %v", err)
	}

	_, err = db.ExecContext(context.Background(),
		`INSERT INTO TRANSLATION_RECORD (x_edit_extracted_data_id, form_id, record_type)
         VALUES (?, ?, ?)`,
		1, "000030FF", "NPC_",
	)
	if err != nil {
		t.Fatalf("expected TRANSLATION_RECORD insert to succeed: %v", err)
	}

	_, err = db.ExecContext(context.Background(),
		`INSERT INTO TRANSLATION_FIELD (translation_record_id, subrecord_type)
         VALUES (?, ?)`,
		1, "FULL",
	)
	if err != nil {
		t.Fatalf("expected TRANSLATION_FIELD insert to succeed: %v", err)
	}

	_, err = db.ExecContext(context.Background(),
		`INSERT INTO TRANSLATION_JOB (x_edit_extracted_data_id, state, created_at)
         VALUES (?, ?, ?)`,
		1, "pending", "2026-01-01T00:00:00Z",
	)
	if err != nil {
		t.Fatalf("expected TRANSLATION_JOB insert to succeed: %v", err)
	}

	_, err = db.ExecContext(context.Background(),
		`INSERT INTO JOB_TRANSLATION_FIELD (translation_job_id, translation_field_id, updated_at)
         VALUES (?, ?, ?)`,
		1, 1, "2026-01-01T00:00:00Z",
	)
	if err != nil {
		t.Fatalf("expected JOB_TRANSLATION_FIELD insert to succeed: %v", err)
	}

	_, err = db.ExecContext(context.Background(),
		`INSERT INTO XTRANSLATOR_OUTPUT_ROW (translation_artifact_id, job_translation_field_id)
         VALUES (?, ?)`,
		9999, 1,
	)
	if err == nil {
		t.Fatal("expected XTRANSLATOR_OUTPUT_ROW insert with invalid translation_artifact_id to fail due to FOREIGN KEY constraint")
	}
	if !strings.Contains(err.Error(), "FOREIGN KEY constraint failed") {
		t.Fatalf("expected FOREIGN KEY constraint error, got: %v", err)
	}
}

// TestCanonicalSchemaXOutputRowFKRejectsInvalidFieldId は
// XTRANSLATOR_OUTPUT_ROW の job_translation_field_id FK が存在しない ID を拒否することを検証する。
// translation_artifact_id は有効な値を使い、job_translation_field_id だけを無効にする。
func TestCanonicalSchemaXOutputRowFKRejectsInvalidFieldId(t *testing.T) {
	db := openMasterDictionaryDatabaseForTest(t, filepath.Join(t.TempDir(), "db", canonicalMigrationTestDatabaseFileName), nil)

	_, err := db.ExecContext(context.Background(),
		`INSERT INTO X_EDIT_EXTRACTED_DATA (source_file_path, source_tool, target_plugin_name, target_plugin_type, record_count, imported_at)
         VALUES (?, ?, ?, ?, ?, ?)`,
		"/tmp/xout2.json", "xEdit", "Skyrim.esm", "esm", 0, "2026-01-01T00:00:00Z",
	)
	if err != nil {
		t.Fatalf("expected X_EDIT_EXTRACTED_DATA insert to succeed: %v", err)
	}

	_, err = db.ExecContext(context.Background(),
		`INSERT INTO TRANSLATION_JOB (x_edit_extracted_data_id, state, created_at)
         VALUES (?, ?, ?)`,
		1, "pending", "2026-01-01T00:00:00Z",
	)
	if err != nil {
		t.Fatalf("expected TRANSLATION_JOB insert to succeed: %v", err)
	}

	_, err = db.ExecContext(context.Background(),
		`INSERT INTO TRANSLATION_ARTIFACT (translation_job_id, artifact_format)
         VALUES (?, ?)`,
		1, "xTranslator",
	)
	if err != nil {
		t.Fatalf("expected TRANSLATION_ARTIFACT insert to succeed: %v", err)
	}

	_, err = db.ExecContext(context.Background(),
		`INSERT INTO XTRANSLATOR_OUTPUT_ROW (translation_artifact_id, job_translation_field_id)
         VALUES (?, ?)`,
		1, 9999,
	)
	if err == nil {
		t.Fatal("expected XTRANSLATOR_OUTPUT_ROW insert with invalid job_translation_field_id to fail due to FOREIGN KEY constraint")
	}
	if !strings.Contains(err.Error(), "FOREIGN KEY constraint failed") {
		t.Fatalf("expected FOREIGN KEY constraint error, got: %v", err)
	}
}

// TestCanonicalSchemaNPCRecordRequiresTranslationRecordId は
// NPC_RECORD.translation_record_id が TRANSLATION_RECORD を参照する FK 制約により
// 存在しない translation_record_id を持つ insert が失敗することを検証する。
func TestCanonicalSchemaNPCRecordRequiresTranslationRecordId(t *testing.T) {
	db := openMasterDictionaryDatabaseForTest(t, filepath.Join(t.TempDir(), "db", canonicalMigrationTestDatabaseFileName), nil)

	_, err := db.ExecContext(context.Background(),
		`INSERT INTO NPC_PROFILE (target_plugin_name, form_id, record_type, created_at, updated_at)
         VALUES (?, ?, ?, ?, ?)`,
		"Skyrim.esm", "000099FF", "NPC_",
		"2026-01-01T00:00:00Z", "2026-01-01T00:00:00Z",
	)
	if err != nil {
		t.Fatalf("expected NPC_PROFILE insert to succeed: %v", err)
	}

	_, err = db.ExecContext(context.Background(),
		`INSERT INTO NPC_RECORD (translation_record_id, npc_profile_id, race, sex, npc_class, voice_type)
         VALUES (?, ?, ?, ?, ?, ?)`,
		9999, 1, "Nord", "Male", "CombatWarrior1H", "",
	)
	if err == nil {
		t.Fatal("expected NPC_RECORD insert with invalid translation_record_id to fail due to FOREIGN KEY constraint")
	}
	if !strings.Contains(err.Error(), "FOREIGN KEY constraint failed") {
		t.Fatalf("expected FOREIGN KEY constraint error, got: %v", err)
	}
}

// TestCanonicalSchemaTranslationJobRequiresExtractedDataId は
// TRANSLATION_JOB.x_edit_extracted_data_id NOT NULL 制約が NULL を拒否することを検証する。
func TestCanonicalSchemaTranslationJobRequiresExtractedDataId(t *testing.T) {
	db := openMasterDictionaryDatabaseForTest(t, filepath.Join(t.TempDir(), "db", canonicalMigrationTestDatabaseFileName), nil)

	_, err := db.ExecContext(context.Background(),
		`INSERT INTO TRANSLATION_JOB (x_edit_extracted_data_id, state, created_at)
         VALUES (NULL, ?, ?)`,
		"pending", "2026-01-01T00:00:00Z",
	)
	if err == nil {
		t.Fatal("expected TRANSLATION_JOB insert with NULL x_edit_extracted_data_id to fail due to NOT NULL constraint")
	}
	if !strings.Contains(err.Error(), "NOT NULL constraint failed") {
		t.Fatalf("expected NOT NULL constraint error, got: %v", err)
	}
}

// ---------------------------------------------------------------------------
// Schema cutover completion_signal tests (handoff: schema-legacy-cutover)
// ---------------------------------------------------------------------------

// TestSchemaCutoverLegacyDictionaryTableAbsent は legacy master_dictionary_entries テーブルが
// migration 適用後に存在しないことを検証する (completion_signal: no legacy dictionary table)。
func TestSchemaCutoverLegacyDictionaryTableAbsent(t *testing.T) {
	db := openMasterDictionaryDatabaseForTest(t, filepath.Join(t.TempDir(), "db", canonicalMigrationTestDatabaseFileName), nil)
	assertTableNotExists(t, db, "master_dictionary_entries")
}

// TestSchemaCutoverLegacyPersonaEntriesTableAbsent は legacy master_persona_entries テーブルが
// migration 適用後に存在しないことを検証する。
func TestSchemaCutoverLegacyPersonaEntriesTableAbsent(t *testing.T) {
	db := openMasterDictionaryDatabaseForTest(t, filepath.Join(t.TempDir(), "db", canonicalMigrationTestDatabaseFileName), nil)
	assertTableNotExists(t, db, "master_persona_entries")
}

// TestSchemaCutoverLegacyPersonaAISettingsTableAbsent は legacy master_persona_ai_settings テーブルが
// migration 適用後に存在しないことを検証する。
func TestSchemaCutoverLegacyPersonaAISettingsTableAbsent(t *testing.T) {
	db := openMasterDictionaryDatabaseForTest(t, filepath.Join(t.TempDir(), "db", canonicalMigrationTestDatabaseFileName), nil)
	assertTableNotExists(t, db, "master_persona_ai_settings")
}

// TestSchemaCutoverLegacyPersonaRunStatusTableAbsent は legacy master_persona_run_status テーブルが
// migration 適用後に存在しないことを検証する (completion_signal: no persisted run status table)。
func TestSchemaCutoverLegacyPersonaRunStatusTableAbsent(t *testing.T) {
	db := openMasterDictionaryDatabaseForTest(t, filepath.Join(t.TempDir(), "db", canonicalMigrationTestDatabaseFileName), nil)
	assertTableNotExists(t, db, "master_persona_run_status")
}

// TestSchemaCutoverPersonaGenerationSettingsTableExists は PERSONA_GENERATION_SETTINGS テーブルが
// migration 適用後に存在することを検証する (completion_signal: PERSONA_GENERATION_SETTINGS exists)。
func TestSchemaCutoverPersonaGenerationSettingsTableExists(t *testing.T) {
	db := openMasterDictionaryDatabaseForTest(t, filepath.Join(t.TempDir(), "db", canonicalMigrationTestDatabaseFileName), nil)
	assertTableExists(t, db, "PERSONA_GENERATION_SETTINGS")
}

// TestSchemaCutoverPersonaGenerationSettingsSingletonConstraintEnforced は
// PERSONA_GENERATION_SETTINGS の id = 1 singleton 制約が機能することを検証する。
// id = 1 は挿入できるが、id = 2 は CHECK 制約で拒否されなければならない。
func TestSchemaCutoverPersonaGenerationSettingsSingletonConstraintEnforced(t *testing.T) {
	db := openMasterDictionaryDatabaseForTest(t, filepath.Join(t.TempDir(), "db", canonicalMigrationTestDatabaseFileName), nil)

	_, err := db.ExecContext(context.Background(),
		`INSERT INTO PERSONA_GENERATION_SETTINGS (id, provider, model) VALUES (1, 'openai', 'gpt-4o')`,
	)
	if err != nil {
		t.Fatalf("expected PERSONA_GENERATION_SETTINGS insert with id=1 to succeed: %v", err)
	}

	_, err = db.ExecContext(context.Background(),
		`INSERT INTO PERSONA_GENERATION_SETTINGS (id, provider, model) VALUES (2, 'openai', 'gpt-4o')`,
	)
	if err == nil {
		t.Fatal("expected PERSONA_GENERATION_SETTINGS insert with id=2 to fail due to CHECK constraint")
	}
	if !strings.Contains(err.Error(), "CHECK constraint failed") {
		t.Fatalf("expected CHECK constraint error for id=2, got: %v", err)
	}
}
