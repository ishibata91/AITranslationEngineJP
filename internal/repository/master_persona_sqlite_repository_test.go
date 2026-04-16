package repository

import (
	"context"
	"path/filepath"
	"testing"
	"time"
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

func TestSQLiteMasterPersonaRepositoriesPersistRunStatusAcrossReopen(t *testing.T) {
	databasePath := filepath.Join(t.TempDir(), "db", sqliteMasterPersonaTestDatabaseFileName)
	repositories := openSQLiteMasterPersonaRepositoriesWithoutCleanup(t, databasePath, nil)
	startedAt := time.Date(2026, 4, 16, 9, 10, 0, 0, time.UTC)
	finishedAt := startedAt.Add(time.Minute)

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

	closeSQLiteMasterPersonaRepositories(t, repositories)
	reopenedRepositories := newSQLiteMasterPersonaRepositoriesForTest(t, databasePath, nil)
	loadedStatus, err := reopenedRepositories.RunStatusRepository.LoadRunStatus(context.Background())
	if err != nil {
		t.Fatalf("expected run status load after reopen to succeed: %v", err)
	}
	if loadedStatus.RunState != "完了" || loadedStatus.SuccessCount != 3 || loadedStatus.ProcessedCount != 4 {
		t.Fatalf("expected persisted run status values, got %#v", loadedStatus)
	}
	if loadedStatus.StartedAt == nil || loadedStatus.FinishedAt == nil {
		t.Fatalf("expected persisted run status timestamps, got %#v", loadedStatus)
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
