package repository

import (
	"context"
	"testing"
	"time"
)

func TestBuildMasterPersonaIdentityKeyUsesPluginFormIDRecordType(t *testing.T) {
	identityKey := BuildMasterPersonaIdentityKey(" FollowersPlus.esp ", " FE01A812 ", " NPC_ ")

	if identityKey != "FollowersPlus.esp:FE01A812:NPC_" {
		t.Fatalf("unexpected identity key: %q", identityKey)
	}
}

func TestInMemoryMasterPersonaRepositoryTreatsPluginAsIdentityBoundary(t *testing.T) {
	now := time.Date(2026, 4, 15, 12, 0, 0, 0, time.UTC)
	repository := NewInMemoryMasterPersonaRepository(nil)

	first, created, err := repository.UpsertIfAbsent(context.Background(), MasterPersonaDraft{
		IdentityKey:  BuildMasterPersonaIdentityKey("FollowersPlus.esp", "FE01A812", "NPC_"),
		TargetPlugin: "FollowersPlus.esp",
		FormID:       "FE01A812",
		RecordType:   "NPC_",
		DisplayName:  "Lys Maren",
		Dialogues:    []string{"one"},
		UpdatedAt:    now,
	})
	if err != nil || !created {
		t.Fatalf("expected first upsert to create entry: entry=%#v created=%v err=%v", first, created, err)
	}

	second, created, err := repository.UpsertIfAbsent(context.Background(), MasterPersonaDraft{
		IdentityKey:  BuildMasterPersonaIdentityKey("NightCourt.esp", "FE01A812", "NPC_"),
		TargetPlugin: "NightCourt.esp",
		FormID:       "FE01A812",
		RecordType:   "NPC_",
		DisplayName:  "Watcher Husk",
		Dialogues:    []string{"two"},
		UpdatedAt:    now,
	})
	if err != nil || !created {
		t.Fatalf("expected second upsert to create plugin-distinct entry: entry=%#v created=%v err=%v", second, created, err)
	}

	followers, err := repository.GetByIdentityKey(context.Background(), BuildMasterPersonaIdentityKey("FollowersPlus.esp", "FE01A812", "NPC_"))
	if err != nil {
		t.Fatalf("expected followers entry to exist: %v", err)
	}
	nightCourt, err := repository.GetByIdentityKey(context.Background(), BuildMasterPersonaIdentityKey("NightCourt.esp", "FE01A812", "NPC_"))
	if err != nil {
		t.Fatalf("expected night court entry to exist: %v", err)
	}
	if followers.TargetPlugin == nightCourt.TargetPlugin {
		t.Fatalf("expected target plugin to remain part of identity boundary: followers=%#v nightCourt=%#v", followers, nightCourt)
	}
}

func TestInMemoryMasterPersonaRepositoryListUsesPluginFilterOnlyForPluginGrouping(t *testing.T) {
	now := time.Date(2026, 4, 15, 12, 0, 0, 0, time.UTC)
	repository := NewInMemoryMasterPersonaRepository(DefaultMasterPersonaSeed(now))

	result, err := repository.List(context.Background(), MasterPersonaListQuery{
		Keyword:      "watcher",
		PluginFilter: "NightCourt.esp",
		Page:         1,
		PageSize:     30,
	})
	if err != nil {
		t.Fatalf("expected list to succeed: %v", err)
	}
	if len(result.Items) != 1 || result.Items[0].TargetPlugin != "NightCourt.esp" {
		t.Fatalf("unexpected filtered items: %#v", result.Items)
	}
	if len(result.PluginGroups) != 1 || result.PluginGroups[0].TargetPlugin != "NightCourt.esp" {
		t.Fatalf("unexpected plugin groups: %#v", result.PluginGroups)
	}
}
