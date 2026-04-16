package service

import (
	"context"
	"errors"
	"os"
	"path/filepath"
	"testing"
	"time"

	"aitranslationenginejp/internal/repository"
)

func TestMasterPersonaGenerationServicePreviewSkipsExistingAndZeroDialogue(t *testing.T) {
	now := func() time.Time { return time.Date(2026, 4, 15, 12, 0, 0, 0, time.UTC) }
	repo := repository.NewInMemoryMasterPersonaRepository(repository.DefaultMasterPersonaSeed(now()))
	secretStore := repository.NewInMemorySecretStore()
	service := NewMasterPersonaGenerationService(repo, repo, repo, secretStore, now, false)
	fixturePath := writeMasterPersonaExtractFixture(t, `{
  "target_plugin": "FollowersPlus.esp",
  "npcs": [
    {
      "form_id": "FE01A812",
      "record_type": "NPC_",
      "editor_id": "FP_LysMaren",
      "display_name": "Lys Maren",
      "dialogues": ["one"]
    },
    {
      "form_id": "FE01A900",
      "record_type": "NPC_",
      "editor_id": "FP_Zero",
      "display_name": "Zero",
      "dialogues": []
    },
    {
      "form_id": "FE01A901",
      "record_type": "NPC_",
      "editor_id": "FP_New",
      "display_name": "New One",
      "dialogues": ["hello", "world"]
    }
  ]
}`)

	result, err := service.Preview(context.Background(), fixturePath, MasterPersonaAISettings{Provider: "fake", Model: "fake-model"})
	if err != nil {
		t.Fatalf("expected preview to succeed: %v", err)
	}
	if result.ExistingSkipCount != 1 || result.ZeroDialogueSkipCount != 1 || result.GeneratableCount != 1 {
		t.Fatalf("unexpected preview result: %#v", result)
	}
}

func TestMasterPersonaGenerationServiceExecuteUsesIdentityKeyWithoutOverwrite(t *testing.T) {
	now := func() time.Time { return time.Date(2026, 4, 15, 12, 0, 0, 0, time.UTC) }
	repo := repository.NewInMemoryMasterPersonaRepository(repository.DefaultMasterPersonaSeed(now()))
	secretStore := repository.NewInMemorySecretStore()
	service := NewMasterPersonaGenerationService(repo, repo, repo, secretStore, now, false)
	fixturePath := writeMasterPersonaExtractFixture(t, `{
  "target_plugin": "FollowersPlus.esp",
  "npcs": [
    {
      "form_id": "FE01A812",
      "record_type": "NPC_",
      "editor_id": "FP_LysMaren",
      "display_name": "Lys Maren",
      "dialogues": ["one"]
    },
    {
      "form_id": "FE01A902",
      "record_type": "NPC_",
      "editor_id": "FP_New",
      "display_name": "Brand New",
      "dialogues": ["alpha"],
      "race": "Nord"
    }
  ]
}`)

	result, err := service.Execute(context.Background(), fixturePath, MasterPersonaAISettings{Provider: "fake", Model: "fake-model"})
	if err != nil {
		t.Fatalf("expected execute to succeed: %v", err)
	}
	if result.SuccessCount != 1 || result.ExistingSkipCount != 1 {
		t.Fatalf("unexpected execute result: %#v", result)
	}
	existing, err := repo.GetByIdentityKey(context.Background(), repository.BuildMasterPersonaIdentityKey("FollowersPlus.esp", "FE01A812", "NPC_"))
	if err != nil {
		t.Fatalf("expected existing entry to remain readable: %v", err)
	}
	if existing.DisplayName != "Lys Maren" {
		t.Fatalf("expected existing entry not to be overwritten: %#v", existing)
	}
}

func TestMasterPersonaGenerationServiceExecuteSkipsZeroDialogueWithoutCreatingEntry(t *testing.T) {
	now := func() time.Time { return time.Date(2026, 4, 15, 12, 0, 0, 0, time.UTC) }
	repo := repository.NewInMemoryMasterPersonaRepository(nil)
	service := NewMasterPersonaGenerationService(repo, repo, repo, repository.NewInMemorySecretStore(), now, false)
	fixturePath := writeMasterPersonaExtractFixture(t, `{
  "target_plugin": "FollowersPlus.esp",
  "npcs": [
    {
      "form_id": "FE01A903",
      "record_type": "NPC_",
      "editor_id": "FP_Zero",
      "display_name": "Zero",
      "dialogues": []
    },
    {
      "form_id": "FE01A904",
      "record_type": "NPC_",
      "editor_id": "FP_New",
      "display_name": "Brand New",
      "dialogues": ["alpha"]
    }
  ]
}`)

	result, err := service.Execute(context.Background(), fixturePath, MasterPersonaAISettings{Provider: "fake", Model: "fake-model"})
	if err != nil {
		t.Fatalf("expected execute to succeed: %v", err)
	}
	if result.ZeroDialogueSkipCount != 1 || result.SuccessCount != 1 {
		t.Fatalf("unexpected execute result: %#v", result)
	}
	_, err = repo.GetByIdentityKey(context.Background(), repository.BuildMasterPersonaIdentityKey("FollowersPlus.esp", "FE01A903", "NPC_"))
	if !errors.Is(err, repository.ErrMasterPersonaEntryNotFound) {
		t.Fatalf("expected zero dialogue entry not to be created, got %v", err)
	}
}

func TestMasterPersonaGenerationServiceExecuteMarksGenericNPCWithoutSurfacingBaselinePhrase(t *testing.T) {
	now := func() time.Time { return time.Date(2026, 4, 15, 12, 0, 0, 0, time.UTC) }
	repo := repository.NewInMemoryMasterPersonaRepository(nil)
	secretStore := repository.NewInMemorySecretStore()
	service := NewMasterPersonaGenerationService(repo, repo, repo, secretStore, now, false)
	fixturePath := writeMasterPersonaExtractFixture(t, `{
  "target_plugin": "NightCourt.esp",
  "npcs": [
    {
      "form_id": "FE01A999",
      "record_type": "NPC_",
      "editor_id": "NC_Generic",
      "display_name": "Generic",
      "dialogues": ["keep distance"]
    }
  ]
}`)

	result, err := service.Execute(context.Background(), fixturePath, MasterPersonaAISettings{Provider: "fake", Model: "fake-model"})
	if err != nil {
		t.Fatalf("expected execute to succeed: %v", err)
	}
	if result.GenericNPCCount != 1 {
		t.Fatalf("expected generic npc count to increment: %#v", result)
	}
	entry, err := repo.GetByIdentityKey(context.Background(), repository.BuildMasterPersonaIdentityKey("NightCourt.esp", "FE01A999", "NPC_"))
	if err != nil {
		t.Fatalf("expected generated entry to exist: %v", err)
	}
	if !entry.BaselineApplied {
		t.Fatalf("expected baseline flag to be applied: %#v", entry)
	}
	if entry.Race != nil || entry.Sex != nil {
		t.Fatalf("expected missing attributes to remain nil: %#v", entry)
	}
	if entry.PersonaBody == masterPersonaNeutralBaseline || entry.PersonaSummary == masterPersonaNeutralBaseline {
		t.Fatalf("expected baseline phrase not to surface in persona-facing text: %#v", entry)
	}
}

func TestMasterPersonaGenerationServiceRejectsUpdateDuringActiveRun(t *testing.T) {
	now := func() time.Time { return time.Date(2026, 4, 15, 12, 0, 0, 0, time.UTC) }
	repo := repository.NewInMemoryMasterPersonaRepository(repository.DefaultMasterPersonaSeed(now()))
	if err := repo.SaveRunStatus(context.Background(), repository.MasterPersonaRunStatusRecord{RunState: MasterPersonaStatusRunning}); err != nil {
		t.Fatalf("expected run status save to succeed: %v", err)
	}
	service := NewMasterPersonaGenerationService(repo, repo, repo, repository.NewInMemorySecretStore(), now, false)

	_, err := service.UpdateEntry(context.Background(), repository.BuildMasterPersonaIdentityKey("FollowersPlus.esp", "FE01A812", "NPC_"), MasterPersonaUpdateInput{FormID: "FE01A812"})
	if !errors.Is(err, ErrMasterPersonaActiveRun) {
		t.Fatalf("expected active run error, got %v", err)
	}
}

func TestMasterPersonaGenerationServiceRejectsDeleteDuringActiveRun(t *testing.T) {
	now := func() time.Time { return time.Date(2026, 4, 15, 12, 0, 0, 0, time.UTC) }
	repo := repository.NewInMemoryMasterPersonaRepository(repository.DefaultMasterPersonaSeed(now()))
	if err := repo.SaveRunStatus(context.Background(), repository.MasterPersonaRunStatusRecord{RunState: MasterPersonaStatusRunning}); err != nil {
		t.Fatalf("expected run status save to succeed: %v", err)
	}
	service := NewMasterPersonaGenerationService(repo, repo, repo, repository.NewInMemorySecretStore(), now, false)

	err := service.DeleteEntry(context.Background(), repository.BuildMasterPersonaIdentityKey("FollowersPlus.esp", "FE01A812", "NPC_"))
	if !errors.Is(err, ErrMasterPersonaActiveRun) {
		t.Fatalf("expected active run error, got %v", err)
	}
}

func TestMasterPersonaGenerationServiceRejectsRealProviderInTestMode(t *testing.T) {
	now := func() time.Time { return time.Date(2026, 4, 15, 12, 0, 0, 0, time.UTC) }
	repo := repository.NewInMemoryMasterPersonaRepository(nil)
	secretStore := repository.NewInMemorySecretStore()
	if err := secretStore.Save(context.Background(), "master-persona:gemini", "saved-real-key"); err != nil {
		t.Fatalf("expected secret save to succeed: %v", err)
	}
	service := NewMasterPersonaGenerationService(repo, repo, repo, secretStore, now, true)
	fixturePath := writeMasterPersonaExtractFixture(t, `{
  "target_plugin": "FollowersPlus.esp",
  "npcs": [
    {
      "form_id": "FE01A910",
      "record_type": "NPC_",
      "editor_id": "FP_Test",
      "display_name": "Test",
      "dialogues": ["hello"]
    }
  ]
}`)

	_, err := service.Preview(context.Background(), fixturePath, MasterPersonaAISettings{Provider: "gemini", Model: "gemini-2.5-pro"})
	if !errors.Is(err, ErrMasterPersonaRealProviderDenied) {
		t.Fatalf("expected real provider denial, got %v", err)
	}
}

func TestMasterPersonaGenerationServiceExecuteFakeProviderWithoutSavedAPIKey(t *testing.T) {
	now := func() time.Time { return time.Date(2026, 4, 15, 12, 0, 0, 0, time.UTC) }
	repo := repository.NewInMemoryMasterPersonaRepository(nil)
	service := NewMasterPersonaGenerationService(repo, repo, repo, repository.NewInMemorySecretStore(), now, true)
	fixturePath := writeMasterPersonaExtractFixture(t, `{
  "target_plugin": "FollowersPlus.esp",
  "npcs": [
    {
      "form_id": "FE01A911",
      "record_type": "NPC_",
      "editor_id": "FP_Fake",
      "display_name": "Fake Path",
      "dialogues": ["hello"]
    }
  ]
}`)

	result, err := service.Execute(context.Background(), fixturePath, MasterPersonaAISettings{Provider: "fake", Model: "fake-master-persona", APIKey: ""})
	if err != nil {
		t.Fatalf("expected fake provider execute to succeed without saved key: %v", err)
	}
	if result.RunState != MasterPersonaStatusCompleted || result.SuccessCount != 1 {
		t.Fatalf("unexpected fake provider execute result: %#v", result)
	}
}

func writeMasterPersonaExtractFixture(t *testing.T, content string) string {
	t.Helper()
	path := filepath.Join(t.TempDir(), "extract.json")
	if err := os.WriteFile(path, []byte(content), 0o600); err != nil {
		t.Fatalf("write fixture: %v", err)
	}
	return path
}
