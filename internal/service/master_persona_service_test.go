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
	service := NewMasterPersonaGenerationService(
		repo,
		repo,
		repo,
		secretStore,
		now,
		false,
		WithMasterPersonaBodyGenerator(newTestSafeMasterPersonaBodyGenerator()),
	)
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

	result, err := service.Preview(context.Background(), fixturePath, MasterPersonaAISettings{Provider: "gemini", Model: "gemini-2.5-pro", APIKey: "transport-key"})
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
	service := NewMasterPersonaGenerationService(
		repo,
		repo,
		repo,
		secretStore,
		now,
		false,
		WithMasterPersonaBodyGenerator(newTestSafeMasterPersonaBodyGenerator()),
	)
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

	result, err := service.Execute(context.Background(), fixturePath, MasterPersonaAISettings{Provider: "gemini", Model: "gemini-2.5-pro", APIKey: "transport-key"})
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
	service := NewMasterPersonaGenerationService(
		repo,
		repo,
		repo,
		repository.NewInMemorySecretStore(),
		now,
		false,
		WithMasterPersonaBodyGenerator(newTestSafeMasterPersonaBodyGenerator()),
	)
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

	result, err := service.Execute(context.Background(), fixturePath, MasterPersonaAISettings{Provider: "gemini", Model: "gemini-2.5-pro", APIKey: "transport-key"})
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
	service := NewMasterPersonaGenerationService(
		repo,
		repo,
		repo,
		secretStore,
		now,
		false,
		WithMasterPersonaBodyGenerator(newTestSafeMasterPersonaBodyGenerator()),
	)
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

	result, err := service.Execute(context.Background(), fixturePath, MasterPersonaAISettings{Provider: "gemini", Model: "gemini-2.5-pro", APIKey: "transport-key"})
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
	service := NewMasterPersonaGenerationService(
		repo,
		repo,
		repo,
		repository.NewInMemorySecretStore(),
		now,
		false,
		WithMasterPersonaBodyGenerator(newTestSafeMasterPersonaBodyGenerator()),
	)

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
	service := NewMasterPersonaGenerationService(
		repo,
		repo,
		repo,
		repository.NewInMemorySecretStore(),
		now,
		false,
		WithMasterPersonaBodyGenerator(newTestSafeMasterPersonaBodyGenerator()),
	)

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

	_, err := service.Preview(context.Background(), fixturePath, MasterPersonaAISettings{Provider: "gemini", Model: "gemini-2.5-pro", APIKey: "transport-key"})
	if !errors.Is(err, ErrMasterPersonaRealProviderDenied) {
		t.Fatalf("expected real provider denial, got %v", err)
	}
}

func TestMasterPersonaGenerationServiceExecuteWithFakeTransportDIWithoutSavedAPIKey(t *testing.T) {
	now := func() time.Time { return time.Date(2026, 4, 15, 12, 0, 0, 0, time.UTC) }
	repo := repository.NewInMemoryMasterPersonaRepository(nil)
	service := NewMasterPersonaGenerationService(
		repo,
		repo,
		repo,
		repository.NewInMemorySecretStore(),
		now,
		true,
		WithMasterPersonaBodyGenerator(newTestSafeMasterPersonaBodyGenerator()),
	)
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

	result, err := service.Execute(context.Background(), fixturePath, MasterPersonaAISettings{Provider: "gemini", Model: "gemini-2.5-pro", APIKey: ""})
	if err != nil {
		t.Fatalf("expected fake transport execute to succeed without saved key: %v", err)
	}
	if result.RunState != MasterPersonaStatusCompleted || result.SuccessCount != 1 {
		t.Fatalf("unexpected fake transport execute result: %#v", result)
	}
}

func TestMasterPersonaGenerationServiceExecuteAllowsLMStudioWithoutSavedAPIKey(t *testing.T) {
	now := func() time.Time { return time.Date(2026, 4, 16, 10, 0, 0, 0, time.UTC) }
	repo := repository.NewInMemoryMasterPersonaRepository(nil)
	generator := &stubMasterPersonaBodyGenerator{body: "lm studio transport body"}
	service := NewMasterPersonaGenerationService(
		repo,
		repo,
		repo,
		repository.NewInMemorySecretStore(),
		now,
		false,
		WithMasterPersonaBodyGenerator(generator),
	)
	fixturePath := writeMasterPersonaExtractFixture(t, `{
  "target_plugin": "FollowersPlus.esp",
  "npcs": [
    {
      "form_id": "FE01A916",
      "record_type": "NPC_",
      "editor_id": "FP_LMStudio",
      "display_name": "LM Studio Path",
      "dialogues": ["hello"]
    }
  ]
}`)

	result, err := service.Execute(context.Background(), fixturePath, MasterPersonaAISettings{Provider: "lm_studio", Model: "local-model", APIKey: ""})
	if err != nil {
		t.Fatalf("expected lm studio execute to succeed without saved key: %v", err)
	}
	if result.RunState != MasterPersonaStatusCompleted || result.SuccessCount != 1 {
		t.Fatalf("unexpected lm studio execute result: %#v", result)
	}
	if generator.calls != 1 {
		t.Fatalf("expected provider generator to be called once, got %d", generator.calls)
	}
	if generator.provider != MasterPersonaProviderLMStudio {
		t.Fatalf("expected lm studio provider id, got %q", generator.provider)
	}
	if generator.apiKey != "" {
		t.Fatalf("expected empty api key for lm studio optional auth, got %q", generator.apiKey)
	}
}

func TestMasterPersonaGenerationServicePreviewAggregatesWhenAISettingsIncomplete(t *testing.T) {
	now := func() time.Time { return time.Date(2026, 4, 16, 10, 0, 0, 0, time.UTC) }
	repo := repository.NewInMemoryMasterPersonaRepository([]repository.MasterPersonaEntry{{
		IdentityKey:  repository.BuildMasterPersonaIdentityKey("FollowersPlus.esp", "FE01A999", "NPC_"),
		TargetPlugin: "FollowersPlus.esp",
		FormID:       "FE01A999",
		RecordType:   "NPC_",
	}})
	service := NewMasterPersonaGenerationService(
		repo,
		repo,
		repo,
		repository.NewInMemorySecretStore(),
		now,
		false,
		WithMasterPersonaBodyGenerator(newTestSafeMasterPersonaBodyGenerator()),
	)
	fixturePath := writeMasterPersonaExtractFixture(t, `{
  "target_plugin": "FollowersPlus.esp",
  "npcs": [
    {
      "form_id": "FE01A912",
      "record_type": "NPC_",
      "editor_id": "FP_Aggregate",
      "display_name": "Aggregate Path",
      "dialogues": ["hello"]
    },
    {
      "form_id": "FE01A913",
      "record_type": "NPC_",
      "editor_id": "FP_Zero",
      "display_name": "Zero Path",
      "dialogues": []
    },
    {
      "form_id": "FE01A999",
      "record_type": "NPC_",
      "editor_id": "FP_Existing",
      "display_name": "Existing Path",
      "dialogues": ["hello"]
    }
  ]
}`)

	result, err := service.Preview(context.Background(), fixturePath, MasterPersonaAISettings{})
	if err != nil {
		t.Fatalf("expected preview aggregation to succeed without completed ai settings: %v", err)
	}
	want := MasterPersonaPreviewResult{
		FileName:              filepath.Base(fixturePath),
		TargetPlugin:          "FollowersPlus.esp",
		TotalNPCCount:         3,
		GeneratableCount:      1,
		ExistingSkipCount:     1,
		ZeroDialogueSkipCount: 1,
		GenericNPCCount:       1,
		Status:                MasterPersonaStatusSettingsIncomplete,
	}
	if result != want {
		t.Fatalf("expected aggregation with settings-incomplete status, got %#v", result)
	}
}

func TestMasterPersonaGenerationServiceExecutePersistsTransportResponseBody(t *testing.T) {
	now := func() time.Time { return time.Date(2026, 4, 16, 10, 0, 0, 0, time.UTC) }
	repo := repository.NewInMemoryMasterPersonaRepository(nil)
	service := NewMasterPersonaGenerationService(
		repo,
		repo,
		repo,
		repository.NewInMemorySecretStore(),
		now,
		false,
		WithMasterPersonaBodyGenerator(&stubMasterPersonaBodyGenerator{body: "transport persona body"}),
	)
	fixturePath := writeMasterPersonaExtractFixture(t, `{
  "target_plugin": "FollowersPlus.esp",
  "npcs": [
    {
      "form_id": "FE01A914",
      "record_type": "NPC_",
      "editor_id": "FP_Transport",
      "display_name": "Transport Path",
      "dialogues": ["hello"]
    }
  ]
}`)

	_, err := service.Execute(context.Background(), fixturePath, MasterPersonaAISettings{Provider: "gemini", Model: "gemini-2.5-pro", APIKey: "transport-key"})
	if err != nil {
		t.Fatalf("expected transport-backed execute to succeed: %v", err)
	}
	entry, err := repo.GetByIdentityKey(context.Background(), repository.BuildMasterPersonaIdentityKey("FollowersPlus.esp", "FE01A914", "NPC_"))
	if err != nil {
		t.Fatalf("expected generated entry to be readable: %v", err)
	}
	if entry.PersonaBody != "transport persona body" {
		t.Fatalf("expected persona body to come from transport response, got %#v", entry)
	}
}

func TestMasterPersonaGenerationServiceTestModeDeniesRealProviderBeforeHTTPCall(t *testing.T) {
	now := func() time.Time { return time.Date(2026, 4, 16, 10, 0, 0, 0, time.UTC) }
	repo := repository.NewInMemoryMasterPersonaRepository(nil)
	generator := &stubMasterPersonaBodyGenerator{err: errors.New("provider must not be called in test mode")}
	service := NewMasterPersonaGenerationService(
		repo,
		repo,
		repo,
		repository.NewInMemorySecretStore(),
		now,
		true,
		WithMasterPersonaBodyGenerator(generator),
	)
	fixturePath := writeMasterPersonaExtractFixture(t, `{
  "target_plugin": "FollowersPlus.esp",
  "npcs": [
    {
      "form_id": "FE01A915",
      "record_type": "NPC_",
      "editor_id": "FP_Denied",
      "display_name": "Denied Path",
      "dialogues": ["hello"]
    }
  ]
}`)

	_, err := service.Execute(context.Background(), fixturePath, MasterPersonaAISettings{Provider: "gemini", Model: "gemini-2.5-pro", APIKey: "real-key"})
	if !errors.Is(err, ErrMasterPersonaRealProviderDenied) {
		t.Fatalf("expected real provider denial before provider call, got %v", err)
	}
	if generator.calls != 0 {
		t.Fatalf("expected provider not to be called in test mode, got %d calls", generator.calls)
	}
}

func TestMasterPersonaGenerationServiceSaveSettingsRejectsUnsupportedProvider(t *testing.T) {
	now := func() time.Time { return time.Date(2026, 4, 15, 12, 0, 0, 0, time.UTC) }
	repo := repository.NewInMemoryMasterPersonaRepository(nil)
	service := NewMasterPersonaGenerationService(
		repo,
		repo,
		repo,
		repository.NewInMemorySecretStore(),
		now,
		false,
		WithMasterPersonaBodyGenerator(newTestSafeMasterPersonaBodyGenerator()),
	)

	_, err := service.SaveSettings(context.Background(), MasterPersonaAISettings{
		Provider: "fake",
		Model:    "gemini-2.5-pro",
		APIKey:   "ignored",
	})
	if !errors.Is(err, ErrMasterPersonaValidation) {
		t.Fatalf("expected unsupported provider validation error, got %v", err)
	}
}

func TestMasterPersonaGenerationServiceSaveSettingsKeepsSavedAPIKeyWhenInputIsBlank(t *testing.T) {
	now := func() time.Time { return time.Date(2026, 4, 15, 12, 0, 0, 0, time.UTC) }
	repo := repository.NewInMemoryMasterPersonaRepository(nil)
	secretStore := repository.NewInMemorySecretStore()
	service := NewMasterPersonaGenerationService(
		repo,
		repo,
		repo,
		secretStore,
		now,
		false,
		WithMasterPersonaBodyGenerator(newTestSafeMasterPersonaBodyGenerator()),
	)

	_, err := service.SaveSettings(context.Background(), MasterPersonaAISettings{
		Provider: "gemini",
		Model:    "gemini-2.5-pro",
		APIKey:   "saved-key",
	})
	if err != nil {
		t.Fatalf("expected first settings save with api key to succeed: %v", err)
	}

	saved, err := service.SaveSettings(context.Background(), MasterPersonaAISettings{
		Provider: "gemini",
		Model:    "gemini-2.5-pro",
		APIKey:   "",
	})
	if err != nil {
		t.Fatalf("expected second settings save without api key to keep persisted key: %v", err)
	}
	if saved.APIKey != "saved-key" {
		t.Fatalf("expected persisted api key to stay available after blank input save, got %#v", saved)
	}

	loaded, err := service.LoadSettings(context.Background())
	if err != nil {
		t.Fatalf("expected settings load to succeed: %v", err)
	}
	if loaded.APIKey != "saved-key" {
		t.Fatalf("expected saved api key to load without re-entry, got %#v", loaded)
	}
}

func TestMasterPersonaGenerationServiceLoadSettingsReturnsSecretLoadError(t *testing.T) {
	now := func() time.Time { return time.Date(2026, 4, 15, 12, 0, 0, 0, time.UTC) }
	repo := repository.NewInMemoryMasterPersonaRepository(nil)
	if err := repo.SaveAISettings(context.Background(), repository.MasterPersonaAISettingsRecord{
		Provider: "gemini",
		Model:    "gemini-2.5-pro",
	}); err != nil {
		t.Fatalf("expected ai settings fixture save to succeed: %v", err)
	}
	secretLoadErr := errors.New("secret load failed")
	service := NewMasterPersonaGenerationService(
		repo,
		repo,
		repo,
		&failingMasterPersonaSecretStore{loadErr: secretLoadErr},
		now,
		false,
	)

	_, err := service.LoadSettings(context.Background())
	if !errors.Is(err, secretLoadErr) {
		t.Fatalf("expected load settings to wrap secret load error, got %v", err)
	}
}

func TestMasterPersonaGenerationServiceSaveSettingsReturnsSecretSaveError(t *testing.T) {
	now := func() time.Time { return time.Date(2026, 4, 15, 12, 0, 0, 0, time.UTC) }
	repo := repository.NewInMemoryMasterPersonaRepository(nil)
	secretSaveErr := errors.New("secret save failed")
	service := NewMasterPersonaGenerationService(
		repo,
		repo,
		repo,
		&failingMasterPersonaSecretStore{saveErr: secretSaveErr},
		now,
		false,
	)

	_, err := service.SaveSettings(context.Background(), MasterPersonaAISettings{
		Provider: "gemini",
		Model:    "gemini-2.5-pro",
		APIKey:   "new-secret",
	})
	if !errors.Is(err, secretSaveErr) {
		t.Fatalf("expected save settings to wrap secret save error, got %v", err)
	}
}

func TestMasterPersonaGenerationServiceSaveSettingsReturnsPersistedSecretLoadError(t *testing.T) {
	now := func() time.Time { return time.Date(2026, 4, 15, 12, 0, 0, 0, time.UTC) }
	repo := repository.NewInMemoryMasterPersonaRepository(nil)
	secretLoadErr := errors.New("persisted secret load failed")
	service := NewMasterPersonaGenerationService(
		repo,
		repo,
		repo,
		&failingMasterPersonaSecretStore{loadErr: secretLoadErr},
		now,
		false,
	)

	_, err := service.SaveSettings(context.Background(), MasterPersonaAISettings{
		Provider: "gemini",
		Model:    "gemini-2.5-pro",
		APIKey:   "",
	})
	if !errors.Is(err, secretLoadErr) {
		t.Fatalf("expected save settings to wrap persisted secret load error, got %v", err)
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

type failingMasterPersonaSecretStore struct {
	loadErr   error
	saveErr   error
	loadValue string
}

func (store *failingMasterPersonaSecretStore) Load(_ context.Context, _ string) (string, error) {
	if store.loadErr != nil {
		return "", store.loadErr
	}
	return store.loadValue, nil
}

func (store *failingMasterPersonaSecretStore) Save(_ context.Context, _ string, _ string) error {
	if store.saveErr != nil {
		return store.saveErr
	}
	return nil
}
