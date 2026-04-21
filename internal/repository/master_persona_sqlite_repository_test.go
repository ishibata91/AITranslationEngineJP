package repository

import (
	"context"
	"database/sql"
	"path/filepath"
	"strings"
	"testing"
	"time"

	sqliteinfra "aitranslationenginejp/internal/infra/sqlite/dbinit"
)

const (
	sqliteMasterPersonaTestDatabaseFileName = "master-persona.sqlite3"
)

func TestSQLiteMasterPersonaRepositoriesPersistEntriesAcrossReopen(t *testing.T) {
	databasePath := filepath.Join(t.TempDir(), "db", sqliteMasterPersonaTestDatabaseFileName)
	repositories := openSQLiteMasterPersonaRepositoriesWithoutCleanup(t, databasePath, nil)
	createdAt := time.Date(2026, 4, 16, 9, 0, 0, 0, time.UTC)
	identityKey := BuildMasterPersonaIdentityKey("FollowersPlus.esp", "FE01AFF0", "NPC_")

	createdEntry, created, err := repositories.EntryRepository.UpsertIfAbsent(context.Background(), MasterPersonaDraft{
		IdentityKey:  identityKey,
		TargetPlugin: "FollowersPlus.esp",
		FormID:       "FE01AFF0",
		RecordType:   "NPC_",
		EditorID:     "FP_Persist",
		DisplayName:  "Persist Entry",
		PersonaBody:  "再起動後も保持される本文",
		Dialogues:    []string{"line"},
		UpdatedAt:    createdAt,
	})
	if err != nil || !created {
		t.Fatalf("expected entry upsert to create record: entry=%#v created=%v err=%v", createdEntry, created, err)
	}

	closeSQLiteMasterPersonaRepositories(t, repositories)
	reopenedRepositories := newSQLiteMasterPersonaRepositoriesForTest(t, databasePath, nil)
	loadedEntry, err := reopenedRepositories.EntryRepository.GetByIdentityKey(context.Background(), identityKey)
	if err != nil {
		t.Fatalf("expected entry to load after reopen: %v", err)
	}
	if loadedEntry.DisplayName != "Persist Entry" || loadedEntry.PersonaBody != "再起動後も保持される本文" {
		t.Fatalf("expected persisted entry to match created value, got %#v", loadedEntry)
	}
}

func TestSQLiteMasterPersonaRepositoriesPersistAISettingsAcrossReopen(t *testing.T) {
	databasePath := filepath.Join(t.TempDir(), "db", sqliteMasterPersonaTestDatabaseFileName)
	repositories := openSQLiteMasterPersonaRepositoriesWithoutCleanup(t, databasePath, nil)

	err := repositories.AISettingsRepository.SaveAISettings(context.Background(), MasterPersonaAISettingsRecord{
		Provider: "gemini",
		Model:    "persisted-model",
	})
	if err != nil {
		t.Fatalf("expected ai settings save to succeed: %v", err)
	}

	closeSQLiteMasterPersonaRepositories(t, repositories)
	reopenedRepositories := newSQLiteMasterPersonaRepositoriesForTest(t, databasePath, nil)
	loadedSettings, err := reopenedRepositories.AISettingsRepository.LoadAISettings(context.Background())
	if err != nil {
		t.Fatalf("expected ai settings load after reopen to succeed: %v", err)
	}
	if loadedSettings.Provider != "gemini" || loadedSettings.Model != "persisted-model" {
		t.Fatalf("expected persisted ai settings, got %#v", loadedSettings)
	}
}

// persona-ai-settings-restart-cutover: run status は DB に保存されないことを証明する。
// SaveRunStatus 後に再オープンしても、LoadRunStatus は初期状態 (入力待ち) を返す。
func TestSQLiteMasterPersonaRepositoriesPersistRunStatusAcrossReopen(t *testing.T) {
	databasePath := filepath.Join(t.TempDir(), "db", sqliteMasterPersonaTestDatabaseFileName)
	repositories := openSQLiteMasterPersonaRepositoriesWithoutCleanup(t, databasePath, nil)
	startedAt := time.Date(2026, 4, 16, 9, 10, 0, 0, time.UTC)
	finishedAt := startedAt.Add(time.Minute)

	// Arrange: run status を保存する。
	err := repositories.RunStatusRepository.SaveRunStatus(context.Background(), MasterPersonaRunStatusRecord{
		RunState:              "完了",
		TargetPlugin:          "FollowersPlus.esp",
		ProcessedCount:        4,
		SuccessCount:          3,
		ExistingSkipCount:     1,
		ZeroDialogueSkipCount: 2,
		GenericNPCCount:       1,
		CurrentActorLabel:     "Persist Actor",
		Message:               "persisted status",
		StartedAt:             &startedAt,
		FinishedAt:            &finishedAt,
	})
	if err != nil {
		t.Fatalf("expected run status save to succeed: %v", err)
	}

	// Act: DB を再オープンする (アプリ再起動をシミュレートする)。
	closeSQLiteMasterPersonaRepositories(t, repositories)
	reopenedRepositories := newSQLiteMasterPersonaRepositoriesForTest(t, databasePath, nil)
	loadedStatus, err := reopenedRepositories.RunStatusRepository.LoadRunStatus(context.Background())
	if err != nil {
		t.Fatalf("expected run status load after reopen to succeed: %v", err)
	}

	// Assert: run status は永続化されないため、初期状態 (入力待ち) に戻る。
	if loadedStatus.RunState != "入力待ち" {
		t.Fatalf("expected run status not to be persisted after reopen (expected RunState=入力待ち), got %#v", loadedStatus)
	}
	if loadedStatus.StartedAt != nil || loadedStatus.FinishedAt != nil {
		t.Fatalf("expected no timestamps after reopen (run state must not be persisted), got %#v", loadedStatus)
	}
}

func TestSQLiteMasterPersonaRepositoriesSeedOnlyWhenDatabaseIsEmpty(t *testing.T) {
	databasePath := filepath.Join(t.TempDir(), "db", sqliteMasterPersonaTestDatabaseFileName)
	seedTime := time.Date(2026, 4, 16, 9, 0, 0, 0, time.UTC)
	seed := DefaultMasterPersonaSeed(seedTime)
	repositories := openSQLiteMasterPersonaRepositoriesWithoutCleanup(t, databasePath, seed)

	_, _, err := repositories.EntryRepository.UpsertIfAbsent(context.Background(), MasterPersonaDraft{
		IdentityKey:  BuildMasterPersonaIdentityKey("FollowersPlus.esp", "FE01AFF1", "NPC_"),
		TargetPlugin: "FollowersPlus.esp",
		FormID:       "FE01AFF1",
		RecordType:   "NPC_",
		DisplayName:  "New Persisted",
		Dialogues:    []string{"line"},
		UpdatedAt:    seedTime.Add(time.Minute),
	})
	if err != nil {
		t.Fatalf("expected new entry creation to succeed: %v", err)
	}

	closeSQLiteMasterPersonaRepositories(t, repositories)
	reopenedRepositories := newSQLiteMasterPersonaRepositoriesForTest(t, databasePath, DefaultMasterPersonaSeed(seedTime.Add(24*time.Hour)))
	listed, err := reopenedRepositories.EntryRepository.List(context.Background(), MasterPersonaListQuery{Page: 1, PageSize: 100})
	if err != nil {
		t.Fatalf("expected list after reopen to succeed: %v", err)
	}
	if listed.TotalCount != len(seed)+1 {
		t.Fatalf("expected seed to run only once and keep created entry, got total=%d", listed.TotalCount)
	}
}

func TestSQLiteMasterPersonaEntryRepositoryListKeepsKeywordPluginGroupsBeforeFilter(t *testing.T) {
	now := time.Date(2026, 4, 16, 9, 0, 0, 0, time.UTC)
	databasePath := filepath.Join(t.TempDir(), "db", sqliteMasterPersonaTestDatabaseFileName)
	repositories := newSQLiteMasterPersonaRepositoriesForTest(t, databasePath, DefaultMasterPersonaSeed(now))

	result, err := repositories.EntryRepository.List(context.Background(), MasterPersonaListQuery{
		Keyword:      "watcher",
		PluginFilter: "NightCourt.esp",
		Page:         1,
		PageSize:     30,
	})
	if err != nil {
		t.Fatalf("expected list query to succeed: %v", err)
	}
	if len(result.Items) != 1 || result.Items[0].TargetPlugin != "NightCourt.esp" {
		t.Fatalf("unexpected filtered items: %#v", result.Items)
	}
	if len(result.PluginGroups) != 1 || result.PluginGroups[0].TargetPlugin != "NightCourt.esp" {
		t.Fatalf("unexpected plugin groups: %#v", result.PluginGroups)
	}
}

func newSQLiteMasterPersonaRepositoriesForTest(
	t *testing.T,
	databasePath string,
	seed []MasterPersonaEntry,
) *SQLiteMasterPersonaRepositories {
	t.Helper()

	repositories := openSQLiteMasterPersonaRepositoriesWithoutCleanup(t, databasePath, seed)
	t.Cleanup(func() {
		closeSQLiteMasterPersonaRepositories(t, repositories)
	})
	return repositories
}

func openSQLiteMasterPersonaRepositoriesWithoutCleanup(
	t *testing.T,
	databasePath string,
	seed []MasterPersonaEntry,
) *SQLiteMasterPersonaRepositories {
	t.Helper()

	repositories, err := NewSQLiteMasterPersonaRepositories(context.Background(), databasePath, seed)
	if err != nil {
		t.Fatalf("expected sqlite master persona repositories to open: %v", err)
	}
	return repositories
}

func closeSQLiteMasterPersonaRepositories(t *testing.T, repositories *SQLiteMasterPersonaRepositories) {
	t.Helper()
	if err := repositories.Close(); err != nil {
		t.Fatalf("expected sqlite master persona repositories close to succeed: %v", err)
	}
}

// persona-ai-settings-restart-cutover: provider と model が PERSONA_GENERATION_SETTINGS に保存され、
// DB 再オープン後に復元されることを証明する。
func TestSQLiteMasterPersonaRepositoriesPersonaAISettingsRestartCutoverProviderModelRestoredAfterReopen(t *testing.T) {
	// Arrange: 空 DB を開き PERSONA_GENERATION_SETTINGS に provider/model を書き込む。
	databasePath := filepath.Join(t.TempDir(), "db", sqliteMasterPersonaTestDatabaseFileName)
	db1, err := sqliteinfra.OpenMasterDictionaryDatabase(context.Background(), databasePath, nil)
	if err != nil {
		t.Fatalf("expected db open to succeed: %v", err)
	}
	_, err = db1.ExecContext(context.Background(),
		"INSERT OR REPLACE INTO PERSONA_GENERATION_SETTINGS (id, provider, model) VALUES (1, 'gemini', 'restart-cutover-model')")
	if err != nil {
		_ = db1.Close()
		t.Fatalf("expected PERSONA_GENERATION_SETTINGS insert to succeed: %v", err)
	}
	if closeErr := db1.Close(); closeErr != nil {
		t.Fatalf("expected db close to succeed: %v", closeErr)
	}

	// Act: DB を再オープンする。
	db2, err := sqliteinfra.OpenMasterDictionaryDatabase(context.Background(), databasePath, nil)
	if err != nil {
		t.Fatalf("expected db reopen to succeed: %v", err)
	}
	defer func() {
		if closeErr := db2.Close(); closeErr != nil {
			t.Fatalf("expected db2 close to succeed: %v", closeErr)
		}
	}()

	// Assert: provider と model が復元されている。
	var provider, model string
	if err := db2.QueryRowContext(context.Background(),
		"SELECT provider, model FROM PERSONA_GENERATION_SETTINGS WHERE id = 1").Scan(&provider, &model); err != nil {
		t.Fatalf("expected PERSONA_GENERATION_SETTINGS load after reopen to succeed: %v", err)
	}
	if provider != "gemini" || model != "restart-cutover-model" {
		t.Fatalf("expected provider/model restored after reopen, got provider=%q model=%q", provider, model)
	}
}

// persona-ai-settings-restart-cutover: PERSONA_GENERATION_SETTINGS テーブルは api_key 列を持たないことを証明する。
// API key は DB の外 (keyring) に保管される契約を schema 境界で確認する。
func TestSQLiteMasterPersonaRepositoriesPersonaAISettingsRestartCutoverAISettingsRecordHasNoAPIKeyField(t *testing.T) {
	// Arrange: migration 適用済み DB を開く。
	databasePath := filepath.Join(t.TempDir(), "db", sqliteMasterPersonaTestDatabaseFileName)
	db, err := sqliteinfra.OpenMasterDictionaryDatabase(context.Background(), databasePath, nil)
	if err != nil {
		t.Fatalf("expected db open to succeed: %v", err)
	}
	defer func() {
		if closeErr := db.Close(); closeErr != nil {
			t.Fatalf("expected db close to succeed: %v", closeErr)
		}
	}()

	// Act: PERSONA_GENERATION_SETTINGS の列情報を取得する。
	rows, err := db.QueryContext(context.Background(), "PRAGMA table_info(PERSONA_GENERATION_SETTINGS)")
	if err != nil {
		t.Fatalf("expected PRAGMA table_info query to succeed: %v", err)
	}
	defer func() {
		if closeErr := rows.Close(); closeErr != nil {
			t.Fatalf("expected rows.Close to succeed: %v", closeErr)
		}
	}()

	// Assert: api_key 列は存在しない。
	var foundAPIKeyColumn bool
	for rows.Next() {
		var cid, notNull, pk int
		var name, colType string
		var defaultValue sql.NullString
		if err := rows.Scan(&cid, &name, &colType, &notNull, &defaultValue, &pk); err != nil {
			t.Fatalf("expected column row scan to succeed: %v", err)
		}
		if strings.Contains(strings.ToLower(name), "api_key") {
			foundAPIKeyColumn = true
		}
	}
	if err := rows.Err(); err != nil {
		t.Fatalf("expected table_info rows iteration to succeed: %v", err)
	}
	if foundAPIKeyColumn {
		t.Fatalf("expected PERSONA_GENERATION_SETTINGS to have no api_key column (API key stays outside DB)")
	}
}

// persona-ai-settings-restart-cutover: run state は DB に保存されないことを証明する。
// master_persona_run_status および PERSONA_GENERATION_RUN_STATUS テーブルが存在しないことで確認する。
func TestSQLiteMasterPersonaRepositoriesPersonaAISettingsRestartCutoverRunStatusIsInputWaitingOnFreshDB(t *testing.T) {
	// Arrange: migration 適用済み DB を開く。
	databasePath := filepath.Join(t.TempDir(), "db", sqliteMasterPersonaTestDatabaseFileName)
	db, err := sqliteinfra.OpenMasterDictionaryDatabase(context.Background(), databasePath, nil)
	if err != nil {
		t.Fatalf("expected db open to succeed: %v", err)
	}
	defer func() {
		if closeErr := db.Close(); closeErr != nil {
			t.Fatalf("expected db close to succeed: %v", closeErr)
		}
	}()

	// Assert: run state 永続化テーブルが存在しない (run state は保存されない)。
	for _, tableName := range []string{"master_persona_run_status", "PERSONA_GENERATION_RUN_STATUS"} {
		var count int
		if err := db.QueryRowContext(context.Background(),
			"SELECT COUNT(*) FROM sqlite_master WHERE type='table' AND name=?", tableName).Scan(&count); err != nil {
			t.Fatalf("expected table existence check to succeed for %q: %v", tableName, err)
		}
		if count != 0 {
			t.Fatalf("expected no run status table %q (run state must not be persisted), but table exists", tableName)
		}
	}
}

// persona-ai-settings-restart-cutover: actual repository path を通じて provider/model が再起動後に復元されることを証明する。
// このテストは master_persona_ai_settings ではなく実際のリポジトリ実装 (PERSONA_GENERATION_SETTINGS) を経由することを要求する。
func TestSQLiteMasterPersonaRepositoriesPersonaAISettingsRestartCutoverProviderModelRestoredThroughRepositoryPath(t *testing.T) {
	// Arrange: actual repository path で AI settings を書き込む。
	databasePath := filepath.Join(t.TempDir(), "db", sqliteMasterPersonaTestDatabaseFileName)
	repos1, err := NewSQLiteMasterPersonaRepositories(context.Background(), databasePath, nil)
	if err != nil {
		t.Fatalf("expected repository open to succeed (repository must not depend on legacy master_persona_ai_settings or master_persona_run_status): %v", err)
	}
	saveErr := repos1.AISettingsRepository.SaveAISettings(context.Background(), MasterPersonaAISettingsRecord{
		Provider: "gemini",
		Model:    "restart-cutover-model",
	})
	if saveErr != nil {
		_ = repos1.Close()
		t.Fatalf("expected ai settings save to succeed: %v", saveErr)
	}
	if closeErr := repos1.Close(); closeErr != nil {
		t.Fatalf("expected repos1 close to succeed: %v", closeErr)
	}

	// Act: DB 再オープン (アプリ再起動をシミュレートする)。
	repos2, err := NewSQLiteMasterPersonaRepositories(context.Background(), databasePath, nil)
	if err != nil {
		t.Fatalf("expected repository reopen to succeed: %v", err)
	}
	defer func() {
		if closeErr := repos2.Close(); closeErr != nil {
			t.Fatalf("expected repos2 close to succeed: %v", closeErr)
		}
	}()

	// Assert: actual repository path を通じて provider/model が復元されている。
	loaded, loadErr := repos2.AISettingsRepository.LoadAISettings(context.Background())
	if loadErr != nil {
		t.Fatalf("expected ai settings load after reopen to succeed: %v", loadErr)
	}
	if loaded.Provider != "gemini" || loaded.Model != "restart-cutover-model" {
		t.Fatalf("expected provider/model restored through repository path, got provider=%q model=%q", loaded.Provider, loaded.Model)
	}
}

// persona-ai-settings-restart-cutover: run state が actual repository path を通じて永続化されないことを証明する。
// DB を再オープンすると run state はデフォルト値 (入力待ち) に戻る。
func TestSQLiteMasterPersonaRepositoriesPersonaAISettingsRestartCutoverRunStateNotPersistedThroughRepositoryPath(t *testing.T) {
	// Arrange: actual repository path で run state を変更する。
	databasePath := filepath.Join(t.TempDir(), "db", sqliteMasterPersonaTestDatabaseFileName)
	repos1, err := NewSQLiteMasterPersonaRepositories(context.Background(), databasePath, nil)
	if err != nil {
		t.Fatalf("expected repository open to succeed (run status must not depend on master_persona_run_status table): %v", err)
	}
	startedAt := time.Date(2026, 4, 21, 10, 0, 0, 0, time.UTC)
	saveErr := repos1.RunStatusRepository.SaveRunStatus(context.Background(), MasterPersonaRunStatusRecord{
		RunState:  "生成中",
		StartedAt: &startedAt,
	})
	if saveErr != nil {
		_ = repos1.Close()
		t.Fatalf("expected run status save to succeed: %v", saveErr)
	}
	if closeErr := repos1.Close(); closeErr != nil {
		t.Fatalf("expected repos1 close to succeed: %v", closeErr)
	}

	// Act: DB 再オープン (アプリ再起動をシミュレートする)。
	repos2, err := NewSQLiteMasterPersonaRepositories(context.Background(), databasePath, nil)
	if err != nil {
		t.Fatalf("expected repository reopen to succeed: %v", err)
	}
	defer func() {
		if closeErr := repos2.Close(); closeErr != nil {
			t.Fatalf("expected repos2 close to succeed: %v", closeErr)
		}
	}()

	// Assert: run state は永続化されていない (再起動後にデフォルト値に戻る)。
	loaded, loadErr := repos2.RunStatusRepository.LoadRunStatus(context.Background())
	if loadErr != nil {
		t.Fatalf("expected run status load after reopen to succeed: %v", loadErr)
	}
	if loaded.RunState == "生成中" {
		t.Fatalf("expected run state to reset after restart (run state must not be persisted to DB), got %q", loaded.RunState)
	}
}

// persona-generation-cutover: migration 002 によって master_persona_entries テーブルが削除されることを証明する。
// 生成の write sink は canonical NPC_PROFILE + PERSONA であるべきであり、legacy テーブルは存在しない。
func TestSQLiteMasterPersonaEntryRepositoryPersonaGenerationCutoverLegacyTableAbsent(t *testing.T) {
	// Arrange: migration 適用済み DB を開く。
	databasePath := filepath.Join(t.TempDir(), "db", sqliteMasterPersonaTestDatabaseFileName)
	db, err := sqliteinfra.OpenMasterDictionaryDatabase(context.Background(), databasePath, nil)
	if err != nil {
		t.Fatalf("expected db open to succeed: %v", err)
	}
	defer func() {
		if closeErr := db.Close(); closeErr != nil {
			t.Fatalf("expected db close to succeed: %v", closeErr)
		}
	}()

	// Assert: master_persona_entries テーブルは schema に存在しない (migration 002 で削除済み)。
	var count int
	if err := db.QueryRowContext(context.Background(),
		"SELECT COUNT(*) FROM sqlite_master WHERE type='table' AND name='master_persona_entries'").Scan(&count); err != nil {
		t.Fatalf("expected table existence check to succeed: %v", err)
	}
	if count != 0 {
		t.Fatalf("expected master_persona_entries to be absent from schema after migration (write sink must be canonical NPC_PROFILE + PERSONA), but table was found")
	}
}

// persona-generation-cutover: migration 003 によって NPC_PROFILE テーブルが作成されることを証明する。
// canonical generation write path の write 先として NPC_PROFILE が schema に存在しなければならない。
func TestSQLiteMasterPersonaEntryRepositoryPersonaGenerationCutoverNPCProfileTablePresent(t *testing.T) {
	// Arrange: migration 適用済み DB を開く。
	databasePath := filepath.Join(t.TempDir(), "db", sqliteMasterPersonaTestDatabaseFileName)
	db, err := sqliteinfra.OpenMasterDictionaryDatabase(context.Background(), databasePath, nil)
	if err != nil {
		t.Fatalf("expected db open to succeed: %v", err)
	}
	defer func() {
		if closeErr := db.Close(); closeErr != nil {
			t.Fatalf("expected db close to succeed: %v", closeErr)
		}
	}()

	// Assert: NPC_PROFILE テーブルが存在する (migration 003 で作成済み)。
	var count int
	if err := db.QueryRowContext(context.Background(),
		"SELECT COUNT(*) FROM sqlite_master WHERE type='table' AND name='NPC_PROFILE'").Scan(&count); err != nil {
		t.Fatalf("expected NPC_PROFILE existence check to succeed: %v", err)
	}
	if count == 0 {
		t.Fatalf("expected NPC_PROFILE table to exist after migration (canonical write target must be present)")
	}
}

// persona-generation-cutover: migration 003 によって PERSONA テーブルが作成されることを証明する。
// canonical generation write path の write 先として PERSONA が schema に存在しなければならない。
func TestSQLiteMasterPersonaEntryRepositoryPersonaGenerationCutoverPersonaTablePresent(t *testing.T) {
	// Arrange: migration 適用済み DB を開く。
	databasePath := filepath.Join(t.TempDir(), "db", sqliteMasterPersonaTestDatabaseFileName)
	db, err := sqliteinfra.OpenMasterDictionaryDatabase(context.Background(), databasePath, nil)
	if err != nil {
		t.Fatalf("expected db open to succeed: %v", err)
	}
	defer func() {
		if closeErr := db.Close(); closeErr != nil {
			t.Fatalf("expected db close to succeed: %v", closeErr)
		}
	}()

	// Assert: PERSONA テーブルが存在する (migration 003 で作成済み)。
	var count int
	if err := db.QueryRowContext(context.Background(),
		"SELECT COUNT(*) FROM sqlite_master WHERE type='table' AND name='PERSONA'").Scan(&count); err != nil {
		t.Fatalf("expected PERSONA table existence check to succeed: %v", err)
	}
	if count == 0 {
		t.Fatalf("expected PERSONA table to exist after migration (canonical persona write target must be present)")
	}
}

// persona-generation-cutover: UpsertIfAbsent が canonical NPC_PROFILE 行を書き込むことを証明する。
// 現在は master_persona_entries を参照しているため FAIL する。
// canonical write path 実装後に PASS となる (RED テスト)。
func TestSQLiteMasterPersonaEntryRepositoryPersonaGenerationCutoverUpsertWritesCanonicalNPCProfile(t *testing.T) {
	// Arrange: seed なしで repository を開く (master_persona_entries への seed アクセスを回避する)。
	databasePath := filepath.Join(t.TempDir(), "db", sqliteMasterPersonaTestDatabaseFileName)
	repos, err := NewSQLiteMasterPersonaRepositories(context.Background(), databasePath, nil)
	if err != nil {
		t.Fatalf("expected repository open to succeed: %v", err)
	}
	defer func() {
		if closeErr := repos.Close(); closeErr != nil {
			t.Fatalf("expected repository close to succeed: %v", closeErr)
		}
	}()

	draft := MasterPersonaDraft{
		IdentityKey:  BuildMasterPersonaIdentityKey("FollowersPlus.esp", "FE01A812", "NPC_"),
		TargetPlugin: "FollowersPlus.esp",
		FormID:       "FE01A812",
		RecordType:   "NPC_",
		EditorID:     "FP_LysMaren",
		DisplayName:  "Lys Maren",
		PersonaBody:  "generation-cutover-persona-body",
		Dialogues:    []string{"line one"},
		UpdatedAt:    time.Date(2026, 4, 21, 0, 0, 0, 0, time.UTC),
	}

	// Act: canonical write path を通じてペルソナを書き込む。
	_, created, upsertErr := repos.EntryRepository.UpsertIfAbsent(context.Background(), draft)
	if upsertErr != nil {
		t.Fatalf("expected UpsertIfAbsent to succeed via canonical NPC_PROFILE write path, got error: %v (canonical write path not yet implemented)", upsertErr)
	}
	if !created {
		t.Fatalf("expected entry to be created, but UpsertIfAbsent returned created=false")
	}

	// Assert: NPC_PROFILE に canonical 行が書き込まれている。
	var npcCount int
	if err := repos.EntryRepository.database.QueryRowContext(context.Background(),
		"SELECT COUNT(*) FROM NPC_PROFILE WHERE target_plugin_name = ? AND form_id = ? AND record_type = ?",
		"FollowersPlus.esp", "FE01A812", "NPC_").Scan(&npcCount); err != nil {
		t.Fatalf("expected NPC_PROFILE count query to succeed: %v", err)
	}
	if npcCount != 1 {
		t.Fatalf("expected 1 NPC_PROFILE row for generated persona, got %d (canonical write path must write to NPC_PROFILE, not master_persona_entries)", npcCount)
	}
}

// persona-generation-cutover: UpsertIfAbsent が canonical PERSONA 行を書き込むことを証明する。
// 現在は master_persona_entries を参照しているため FAIL する。
// canonical write path 実装後に PASS となる (RED テスト)。
func TestSQLiteMasterPersonaEntryRepositoryPersonaGenerationCutoverUpsertWritesCanonicalPersona(t *testing.T) {
	// Arrange: seed なしで repository を開く。
	databasePath := filepath.Join(t.TempDir(), "db", sqliteMasterPersonaTestDatabaseFileName)
	repos, err := NewSQLiteMasterPersonaRepositories(context.Background(), databasePath, nil)
	if err != nil {
		t.Fatalf("expected repository open to succeed: %v", err)
	}
	defer func() {
		if closeErr := repos.Close(); closeErr != nil {
			t.Fatalf("expected repository close to succeed: %v", closeErr)
		}
	}()

	draft := MasterPersonaDraft{
		IdentityKey:  BuildMasterPersonaIdentityKey("FollowersPlus.esp", "FE01A812", "NPC_"),
		TargetPlugin: "FollowersPlus.esp",
		FormID:       "FE01A812",
		RecordType:   "NPC_",
		EditorID:     "FP_LysMaren",
		DisplayName:  "Lys Maren",
		PersonaBody:  "generation-cutover-persona-body",
		Dialogues:    []string{"line one"},
		UpdatedAt:    time.Date(2026, 4, 21, 0, 0, 0, 0, time.UTC),
	}

	// Act: canonical write path を通じてペルソナを書き込む。
	_, _, upsertErr := repos.EntryRepository.UpsertIfAbsent(context.Background(), draft)
	if upsertErr != nil {
		t.Fatalf("expected UpsertIfAbsent to succeed via canonical PERSONA write path, got error: %v (canonical write path not yet implemented)", upsertErr)
	}

	// Assert: PERSONA に NPC_PROFILE に紐づく canonical 行が書き込まれている。
	var personaCount int
	if err := repos.EntryRepository.database.QueryRowContext(context.Background(),
		`SELECT COUNT(*) FROM PERSONA p
		 JOIN NPC_PROFILE np ON p.npc_profile_id = np.id
		 WHERE np.target_plugin_name = ? AND np.form_id = ? AND np.record_type = ?`,
		"FollowersPlus.esp", "FE01A812", "NPC_").Scan(&personaCount); err != nil {
		t.Fatalf("expected PERSONA join NPC_PROFILE count query to succeed: %v", err)
	}
	if personaCount != 1 {
		t.Fatalf("expected 1 PERSONA row linked to NPC_PROFILE for generated persona, got %d (canonical write path must write to PERSONA, not master_persona_entries)", personaCount)
	}
}

// persona-generation-cutover: master_persona_entries は shipped sqlite repository path の write sink でないことを証明する。
// migration 002 でテーブルが削除されており、canonical write は NPC_PROFILE + PERSONA へ行く。
func TestSQLiteMasterPersonaEntryRepositoryPersonaGenerationCutoverMasterPersonaEntriesNotWriteSink(t *testing.T) {
	// Arrange: seed なしで repository を開く。
	databasePath := filepath.Join(t.TempDir(), "db", sqliteMasterPersonaTestDatabaseFileName)
	repos, err := NewSQLiteMasterPersonaRepositories(context.Background(), databasePath, nil)
	if err != nil {
		t.Fatalf("expected repository open to succeed: %v", err)
	}
	defer func() {
		if closeErr := repos.Close(); closeErr != nil {
			t.Fatalf("expected repository close to succeed: %v", closeErr)
		}
	}()

	// Assert: master_persona_entries テーブルが schema に存在しない (write sink は canonical テーブルのみ)。
	var tableCount int
	if err := repos.EntryRepository.database.QueryRowContext(context.Background(),
		"SELECT COUNT(*) FROM sqlite_master WHERE type='table' AND name='master_persona_entries'").Scan(&tableCount); err != nil {
		t.Fatalf("expected table existence check to succeed: %v", err)
	}
	if tableCount != 0 {
		t.Fatalf("expected master_persona_entries to be absent from schema (generation must write to canonical NPC_PROFILE + PERSONA, not legacy table)")
	}
}

// persona-generation-cutover: UpsertIfAbsent が失敗したとき、partial canonical rows が残らないことを証明する。
// canonical write path は NPC_PROFILE と PERSONA をアトミックに書き込む必要がある。
func TestSQLiteMasterPersonaEntryRepositoryPersonaGenerationCutoverFailureLeavesNoPartialCanonicalRows(t *testing.T) {
	// Arrange: seed なしで repository を開く。
	databasePath := filepath.Join(t.TempDir(), "db", sqliteMasterPersonaTestDatabaseFileName)
	repos, err := NewSQLiteMasterPersonaRepositories(context.Background(), databasePath, nil)
	if err != nil {
		t.Fatalf("expected repository open to succeed: %v", err)
	}
	defer func() {
		if closeErr := repos.Close(); closeErr != nil {
			t.Fatalf("expected repository close to succeed: %v", closeErr)
		}
	}()

	draft := MasterPersonaDraft{
		IdentityKey:  BuildMasterPersonaIdentityKey("FollowersPlus.esp", "FE01A812", "NPC_"),
		TargetPlugin: "FollowersPlus.esp",
		FormID:       "FE01A812",
		RecordType:   "NPC_",
		EditorID:     "FP_LysMaren",
		DisplayName:  "Lys Maren",
		PersonaBody:  "cutover-failure-test-body",
		Dialogues:    []string{"line"},
		UpdatedAt:    time.Date(2026, 4, 21, 0, 0, 0, 0, time.UTC),
	}

	// Act: UpsertIfAbsent を試みる。現在は master_persona_entries が存在しないため失敗する。
	_, _, upsertErr := repos.EntryRepository.UpsertIfAbsent(context.Background(), draft)

	// Assert: 失敗後に partial canonical rows が残らない (write はアトミックであるべき)。
	if upsertErr != nil {
		var npcCount int
		if err := repos.EntryRepository.database.QueryRowContext(context.Background(),
			"SELECT COUNT(*) FROM NPC_PROFILE").Scan(&npcCount); err != nil {
			t.Fatalf("expected NPC_PROFILE count query to succeed: %v", err)
		}
		if npcCount != 0 {
			t.Fatalf("expected no partial NPC_PROFILE rows after failed write (canonical write must be atomic), got %d", npcCount)
		}
		var personaCount int
		if err := repos.EntryRepository.database.QueryRowContext(context.Background(),
			"SELECT COUNT(*) FROM PERSONA").Scan(&personaCount); err != nil {
			t.Fatalf("expected PERSONA count query to succeed: %v", err)
		}
		if personaCount != 0 {
			t.Fatalf("expected no partial PERSONA rows after failed write (canonical write must be atomic), got %d", personaCount)
		}
	}
}
