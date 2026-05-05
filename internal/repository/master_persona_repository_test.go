package repository

import (
	"context"
	"testing"
	"time"
)

func testMasterPersonaSeed(now time.Time) []MasterPersonaEntry {
	testNPCRace := "TestRace"
	testNPCSex := "TestSex"
	return []MasterPersonaEntry{
		{
			IdentityKey:    BuildMasterPersonaIdentityKey("TestPersonaPluginA.esp", "FE01A812", "NPC_"),
			TargetPlugin:   "TestPersonaPluginA.esp",
			FormID:         "FE01A812",
			RecordType:     "NPC_",
			EditorID:       "TEST_NPC_A",
			DisplayName:    "Test NPC A",
			Race:           &testNPCRace,
			Sex:            &testNPCSex,
			VoiceType:      "TestVoiceA",
			ClassName:      "TestClassA",
			SourcePlugin:   "TestPersonaPluginA.esp",
			PersonaSummary: "test summary a",
			PersonaBody:    "test body a",
			DialogueCount:  1,
			Dialogues:      []string{"test line a"},
			UpdatedAt:      now,
		},
		{
			IdentityKey:    BuildMasterPersonaIdentityKey("TestPersonaPluginB.esp", "FE01A814", "NPC_"),
			TargetPlugin:   "TestPersonaPluginB.esp",
			FormID:         "FE01A814",
			RecordType:     "NPC_",
			EditorID:       "TEST_NPC_B",
			DisplayName:    "Test NPC B",
			VoiceType:      "TestVoiceB",
			ClassName:      "TestClassB",
			SourcePlugin:   "TestPersonaPluginB.esp",
			PersonaSummary: "test summary b",
			PersonaBody:    "test body b",
			DialogueCount:  1,
			Dialogues:      []string{"test line b"},
			UpdatedAt:      now.Add(-time.Minute),
		},
	}
}

func TestBuildMasterPersonaIdentityKeyUsesPluginFormIDRecordType(t *testing.T) {
	identityKey := BuildMasterPersonaIdentityKey(" TestPersonaPluginA.esp ", " FE01A812 ", " NPC_ ")

	if identityKey != "TestPersonaPluginA.esp:FE01A812:NPC_" {
		t.Fatalf("unexpected identity key: %q", identityKey)
	}
}

func TestInMemoryMasterPersonaRepositoryTreatsPluginAsIdentityBoundary(t *testing.T) {
	now := time.Date(2026, 4, 15, 12, 0, 0, 0, time.UTC)
	repository := NewInMemoryMasterPersonaRepository(nil)

	first, created, err := repository.UpsertIfAbsent(context.Background(), MasterPersonaDraft{
		IdentityKey:  BuildMasterPersonaIdentityKey("TestPersonaPluginA.esp", "FE01A812", "NPC_"),
		TargetPlugin: "TestPersonaPluginA.esp",
		FormID:       "FE01A812",
		RecordType:   "NPC_",
		DisplayName:  "Test NPC A",
		Dialogues:    []string{"one"},
		UpdatedAt:    now,
	})
	if err != nil || !created {
		t.Fatalf("expected first upsert to create entry: entry=%#v created=%v err=%v", first, created, err)
	}

	second, created, err := repository.UpsertIfAbsent(context.Background(), MasterPersonaDraft{
		IdentityKey:  BuildMasterPersonaIdentityKey("TestPersonaPluginB.esp", "FE01A812", "NPC_"),
		TargetPlugin: "TestPersonaPluginB.esp",
		FormID:       "FE01A812",
		RecordType:   "NPC_",
		DisplayName:  "Test NPC C",
		Dialogues:    []string{"two"},
		UpdatedAt:    now,
	})
	if err != nil || !created {
		t.Fatalf("expected second upsert to create plugin-distinct entry: entry=%#v created=%v err=%v", second, created, err)
	}

	followers, err := repository.GetByIdentityKey(context.Background(), BuildMasterPersonaIdentityKey("TestPersonaPluginA.esp", "FE01A812", "NPC_"))
	if err != nil {
		t.Fatalf("expected followers entry to exist: %v", err)
	}
	nightCourt, err := repository.GetByIdentityKey(context.Background(), BuildMasterPersonaIdentityKey("TestPersonaPluginB.esp", "FE01A812", "NPC_"))
	if err != nil {
		t.Fatalf("expected night court entry to exist: %v", err)
	}
	if followers.TargetPlugin == nightCourt.TargetPlugin {
		t.Fatalf("expected target plugin to remain part of identity boundary: followers=%#v nightCourt=%#v", followers, nightCourt)
	}
}

func TestInMemoryMasterPersonaRepositoryListUsesPluginFilterOnlyForPluginGrouping(t *testing.T) {
	now := time.Date(2026, 4, 15, 12, 0, 0, 0, time.UTC)
	repository := NewInMemoryMasterPersonaRepository(testMasterPersonaSeed(now))

	result, err := repository.List(context.Background(), MasterPersonaListQuery{
		Keyword:      "test npc b",
		PluginFilter: "TestPersonaPluginB.esp",
		Page:         1,
		PageSize:     30,
	})
	if err != nil {
		t.Fatalf("expected list to succeed: %v", err)
	}
	if len(result.Items) != 1 || result.Items[0].TargetPlugin != "TestPersonaPluginB.esp" {
		t.Fatalf("unexpected filtered items: %#v", result.Items)
	}
	if len(result.PluginGroups) != 1 || result.PluginGroups[0].TargetPlugin != "TestPersonaPluginB.esp" {
		t.Fatalf("unexpected plugin groups: %#v", result.PluginGroups)
	}
}

func TestInMemoryMasterPersonaRepositoryUpdateExistingEntry(t *testing.T) {
	now := time.Date(2026, 4, 15, 12, 0, 0, 0, time.UTC)
	repo := NewInMemoryMasterPersonaRepository(nil)
	identityKey := BuildMasterPersonaIdentityKey("Test.esp", "AABBCC", "NPC_")

	_, _, err := repo.UpsertIfAbsent(context.Background(), MasterPersonaDraft{
		IdentityKey:  identityKey,
		TargetPlugin: "Test.esp",
		FormID:       "AABBCC",
		RecordType:   "NPC_",
		DisplayName:  "Before",
		UpdatedAt:    now,
	})
	if err != nil {
		t.Fatalf("upsert failed: %v", err)
	}

	updated, err := repo.Update(context.Background(), identityKey, MasterPersonaDraft{
		IdentityKey:  identityKey,
		TargetPlugin: "Test.esp",
		FormID:       "AABBCC",
		RecordType:   "NPC_",
		DisplayName:  "After",
		UpdatedAt:    now,
	})
	if err != nil {
		t.Fatalf("update failed: %v", err)
	}
	if updated.DisplayName != "After" {
		t.Errorf("expected DisplayName=After, got %q", updated.DisplayName)
	}
}

func TestInMemoryMasterPersonaRepositoryUpdateNotFoundReturnsError(t *testing.T) {
	repo := NewInMemoryMasterPersonaRepository(nil)
	_, err := repo.Update(context.Background(), "missing:key:NPC_", MasterPersonaDraft{
		IdentityKey: "missing:key:NPC_",
	})
	if err == nil {
		t.Fatal("expected error when updating non-existent entry")
	}
}

func TestInMemoryMasterPersonaRepositoryDeleteExistingEntry(t *testing.T) {
	now := time.Date(2026, 4, 15, 12, 0, 0, 0, time.UTC)
	repo := NewInMemoryMasterPersonaRepository(nil)
	identityKey := BuildMasterPersonaIdentityKey("Del.esp", "112233", "NPC_")

	_, _, err := repo.UpsertIfAbsent(context.Background(), MasterPersonaDraft{
		IdentityKey:  identityKey,
		TargetPlugin: "Del.esp",
		FormID:       "112233",
		RecordType:   "NPC_",
		DisplayName:  "ToDelete",
		UpdatedAt:    now,
	})
	if err != nil {
		t.Fatalf("upsert failed: %v", err)
	}

	if err := repo.Delete(context.Background(), identityKey); err != nil {
		t.Fatalf("delete failed: %v", err)
	}
	if _, err := repo.GetByIdentityKey(context.Background(), identityKey); err == nil {
		t.Fatal("expected error after deletion")
	}
}

func TestInMemoryMasterPersonaRepositoryDeleteNotFoundReturnsError(t *testing.T) {
	repo := NewInMemoryMasterPersonaRepository(nil)
	if err := repo.Delete(context.Background(), "ghost:key:NPC_"); err == nil {
		t.Fatal("expected error when deleting non-existent entry")
	}
}

func TestInMemoryMasterPersonaRepositoryLoadSaveAISettings(t *testing.T) {
	repo := NewInMemoryMasterPersonaRepository(nil)

	loaded, err := repo.LoadAISettings(context.Background())
	if err != nil {
		t.Fatalf("LoadAISettings failed: %v", err)
	}
	if loaded.Provider != "" || loaded.Model != "" {
		t.Fatalf("expected zero value, got %#v", loaded)
	}

	want := MasterPersonaAISettingsRecord{Provider: "openai", Model: "gpt-4o"}
	if saveErr := repo.SaveAISettings(context.Background(), want); saveErr != nil {
		t.Fatalf("SaveAISettings failed: %v", saveErr)
	}

	got, err := repo.LoadAISettings(context.Background())
	if err != nil {
		t.Fatalf("LoadAISettings after save failed: %v", err)
	}
	if got != want {
		t.Errorf("expected %#v, got %#v", want, got)
	}
}

func TestInMemoryMasterPersonaRepositoryLoadSaveRunStatus(t *testing.T) {
	repo := NewInMemoryMasterPersonaRepository(nil)

	loaded, err := repo.LoadRunStatus(context.Background())
	if err != nil {
		t.Fatalf("LoadRunStatus failed: %v", err)
	}
	if loaded.ProcessedCount != 0 {
		t.Fatalf("expected zero ProcessedCount on empty repo, got %#v", loaded)
	}

	want := MasterPersonaRunStatusRecord{RunState: "running", ProcessedCount: 5}
	if saveErr := repo.SaveRunStatus(context.Background(), want); saveErr != nil {
		t.Fatalf("SaveRunStatus failed: %v", saveErr)
	}

	got, err := repo.LoadRunStatus(context.Background())
	if err != nil {
		t.Fatalf("LoadRunStatus after save failed: %v", err)
	}
	if got.RunState != want.RunState || got.ProcessedCount != want.ProcessedCount {
		t.Errorf("expected %#v, got %#v", want, got)
	}
}
