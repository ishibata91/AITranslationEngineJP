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
	// zero dialogue NPC は parse-time 除外なので ZeroDialogueSkipCount は 0
	if result.ExistingSkipCount != 1 || result.ZeroDialogueSkipCount != 0 || result.GeneratableCount != 1 {
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
	// zero dialogue NPC は parse-time 除外なので ZeroDialogueSkipCount は 0
	if result.ZeroDialogueSkipCount != 0 || result.SuccessCount != 1 {
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
	// zero dialogue NPC は parse-time 除外: TotalNPCCount は除外後の数、ZeroDialogueSkipCount は 0
	want := MasterPersonaPreviewResult{
		FileName:              filepath.Base(fixturePath),
		TargetPlugin:          "FollowersPlus.esp",
		TotalNPCCount:         2,
		GeneratableCount:      1,
		ExistingSkipCount:     1,
		ZeroDialogueSkipCount: 0,
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

// persona-read-detail-cutover: update save regression - identityKey のみ指定し FormID を省略した
// narrow update input が成功することを証明する。
// UpdateEntry が input.FormID を必須検証している間は失敗する。
func TestMasterPersonaGenerationServicePersonaReadDetailCutoverUpdateSucceedsWithNarrowInput(t *testing.T) {
	now := func() time.Time { return time.Date(2026, 4, 15, 12, 0, 0, 0, time.UTC) }
	repo := repository.NewInMemoryMasterPersonaRepository(repository.DefaultMasterPersonaSeed(now()))
	service := NewMasterPersonaGenerationService(
		repo,
		repo,
		repo,
		repository.NewInMemorySecretStore(),
		now,
		false,
		WithMasterPersonaBodyGenerator(newTestSafeMasterPersonaBodyGenerator()),
	)

	identityKey := repository.BuildMasterPersonaIdentityKey("FollowersPlus.esp", "FE01A812", "NPC_")

	// Act: FormID を省略した narrow update input (frontend cutover 後の送信形式)
	updated, err := service.UpdateEntry(context.Background(), identityKey, MasterPersonaUpdateInput{
		DisplayName: "Updated Lys",
		PersonaBody: "updated persona body",
	})

	// Assert: identityKey から identity を補完し FormID なしで成功するべき
	if err != nil {
		t.Fatalf("expected update to succeed with narrow input (no FormID); got: %v", err)
	}
	// DisplayName は read-only のため元の値が保持される (cutover 後の approved behavior)
	if updated.DisplayName != "Lys Maren" {
		t.Fatalf("expected DisplayName to remain 'Lys Maren' (read-only after cutover), got %q", updated.DisplayName)
	}
	if updated.PersonaBody != "updated persona body" {
		t.Fatalf("expected PersonaBody to be updated, got %q", updated.PersonaBody)
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

// persona-ai-settings-restart-cutover: SaveSettings は API キーを DB レコードに格納せず、
// secret store にだけ保存することを証明する。
func TestMasterPersonaGenerationServicePersonaAISettingsRestartCutoverSaveSettingsKeepsAPIKeyInSecretStoreOnly(t *testing.T) {
	now := func() time.Time { return time.Date(2026, 4, 16, 0, 0, 0, 0, time.UTC) }
	repo := repository.NewInMemoryMasterPersonaRepository(nil)
	secretStore := repository.NewInMemorySecretStore()
	svc := NewMasterPersonaGenerationService(repo, repo, repo, secretStore, now, false)

	_, err := svc.SaveSettings(context.Background(), MasterPersonaAISettings{
		Provider: "gemini",
		Model:    "cutover-model",
		APIKey:   "secret-api-key",
	})
	if err != nil {
		t.Fatalf("expected save settings to succeed: %v", err)
	}

	dbRecord, err := repo.LoadAISettings(context.Background())
	if err != nil {
		t.Fatalf("expected db record load to succeed: %v", err)
	}
	if dbRecord.Provider != "gemini" || dbRecord.Model != "cutover-model" {
		t.Fatalf("expected provider/model in db record, got %#v", dbRecord)
	}
	// MasterPersonaAISettingsRecord has no APIKey field: API key is not in the DB record.
	storedKey, err := secretStore.Load(context.Background(), "master-persona:gemini")
	if err != nil {
		t.Fatalf("expected secret store load to succeed: %v", err)
	}
	if storedKey != "secret-api-key" {
		t.Fatalf("expected api key in secret store only, got %q", storedKey)
	}
}

// persona-ai-settings-restart-cutover: LoadSettings は再起動後に secret store から API キーを復元することを証明する。
func TestMasterPersonaGenerationServicePersonaAISettingsRestartCutoverLoadSettingsRestoresAPIKeyFromSecretStore(t *testing.T) {
	now := func() time.Time { return time.Date(2026, 4, 16, 0, 0, 0, 0, time.UTC) }
	repo := repository.NewInMemoryMasterPersonaRepository(nil)
	secretStore := repository.NewInMemorySecretStore()

	// Arrange: DB に provider+model のみ保存、secret store に API キーを格納（再起動後の状態を模擬）
	if err := repo.SaveAISettings(context.Background(), repository.MasterPersonaAISettingsRecord{
		Provider: "gemini",
		Model:    "restored-model",
	}); err != nil {
		t.Fatalf("expected db record setup to succeed: %v", err)
	}
	if err := secretStore.Save(context.Background(), "master-persona:gemini", "restored-api-key"); err != nil {
		t.Fatalf("expected secret store setup to succeed: %v", err)
	}
	svc := NewMasterPersonaGenerationService(repo, repo, repo, secretStore, now, false)

	settings, err := svc.LoadSettings(context.Background())
	if err != nil {
		t.Fatalf("expected load settings to succeed: %v", err)
	}
	if settings.Provider != "gemini" || settings.Model != "restored-model" {
		t.Fatalf("expected provider/model restored, got %#v", settings)
	}
	if settings.APIKey != "restored-api-key" {
		t.Fatalf("expected api key restored from secret store, got %q", settings.APIKey)
	}
}

// RED test: persona-json-preview-cutover - zero-dialogue NPC が parse 時点で除外され
// preview の ZeroDialogueSkipCount が 0 になることを証明する。
// analyzePreview がゼロ会話 NPC を parse-time filter で除去するまで失敗する。
func TestPersonaJSONPreviewCutoverServicePreviewExcludesZeroDialogueAtParseTime(t *testing.T) {
	now := func() time.Time { return time.Date(2026, 4, 15, 12, 0, 0, 0, time.UTC) }
	repo := repository.NewInMemoryMasterPersonaRepository(nil)
	service := NewMasterPersonaGenerationService(
		repo, repo, repo, repository.NewInMemorySecretStore(), now, false,
		WithMasterPersonaBodyGenerator(newTestSafeMasterPersonaBodyGenerator()),
	)
	fixturePath := writeMasterPersonaExtractFixture(t, `{
  "target_plugin": "CutoverTest.esp",
  "npcs": [
    {
      "form_id": "FE02A001",
      "record_type": "NPC_",
      "editor_id": "CT_WithDialogue",
      "display_name": "With Dialogue",
      "dialogues": ["hello"]
    },
    {
      "form_id": "FE02A002",
      "record_type": "NPC_",
      "editor_id": "CT_NoDialogue",
      "display_name": "No Dialogue",
      "dialogues": []
    }
  ]
}`)

	result, err := service.Preview(context.Background(), fixturePath, MasterPersonaAISettings{Provider: "gemini", Model: "gemini-2.5-pro", APIKey: "transport-key"})
	if err != nil {
		t.Fatalf("expected preview to succeed: %v", err)
	}
	if result.ZeroDialogueSkipCount != 0 {
		t.Fatalf("expected zero-dialogue NPCs to be excluded at parse time (ZeroDialogueSkipCount=0), got %d", result.ZeroDialogueSkipCount)
	}
}

// persona-json-preview-cutover: all-zero-dialogue JSON は validation error ではなく 0 件 preview を返すこと。
// parseMasterPersonaExtractNPCList が parse-time filter で全 NPC を除外した場合、
// Preview は err=nil で ZeroDialogueSkipCount>0 / GeneratableCount=0 を返す必要がある。
// product code が validation error を返す間はこの test は FAIL する (RED)。
func TestPersonaJSONPreviewCutoverAllZeroDialogueReturnsZeroCountNotValidationError(t *testing.T) {
	now := func() time.Time { return time.Date(2026, 4, 15, 12, 0, 0, 0, time.UTC) }
	repo := repository.NewInMemoryMasterPersonaRepository(nil)
	service := NewMasterPersonaGenerationService(
		repo, repo, repo, repository.NewInMemorySecretStore(), now, false,
		WithMasterPersonaBodyGenerator(newTestSafeMasterPersonaBodyGenerator()),
	)
	fixturePath := writeMasterPersonaExtractFixture(t, `{
  "target_plugin": "AllZero.esp",
  "npcs": [
    {
      "form_id": "FE03A001",
      "record_type": "NPC_",
      "editor_id": "AZ_First",
      "display_name": "First Zero",
      "dialogues": []
    },
    {
      "form_id": "FE03A002",
      "record_type": "NPC_",
      "editor_id": "AZ_Second",
      "display_name": "Second Zero",
      "dialogues": []
    }
  ]
}`)

	result, err := service.Preview(context.Background(), fixturePath, MasterPersonaAISettings{Provider: "gemini", Model: "gemini-2.5-pro", APIKey: "transport-key"})
	if err != nil {
		t.Fatalf("expected Preview to succeed with all-zero-dialogue JSON (not a validation error), got: %v", err)
	}
	if result.ZeroDialogueSkipCount != 2 {
		t.Fatalf("expected ZeroDialogueSkipCount=2 for 2 zero-dialogue NPCs, got %d", result.ZeroDialogueSkipCount)
	}
	if result.GeneratableCount != 0 {
		t.Fatalf("expected GeneratableCount=0 when all NPCs are zero-dialogue, got %d", result.GeneratableCount)
	}
}

// persona-json-preview-cutover: identity boundary が plugin+form_id+record_type で維持されることを証明する。
// cutover 後も既存 NPC の identity key 解決が plugin+form_id+record_type に依存することを確認する。
func TestPersonaJSONPreviewCutoverServiceExistingIdentityBoundaryUsesPluginFormIDRecordType(t *testing.T) {
	now := func() time.Time { return time.Date(2026, 4, 15, 12, 0, 0, 0, time.UTC) }
	repo := repository.NewInMemoryMasterPersonaRepository(repository.DefaultMasterPersonaSeed(now()))
	service := NewMasterPersonaGenerationService(
		repo, repo, repo, repository.NewInMemorySecretStore(), now, false,
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
    }
  ]
}`)

	result, err := service.Preview(context.Background(), fixturePath, MasterPersonaAISettings{Provider: "gemini", Model: "gemini-2.5-pro", APIKey: "transport-key"})
	if err != nil {
		t.Fatalf("expected preview to succeed: %v", err)
	}
	if result.ExistingSkipCount != 1 {
		t.Fatalf("expected NPC with plugin+form_id+record_type identity to be counted as existing: %#v", result)
	}
}

// persona-generation-cutover: 既存 identity は generation によって上書きされないことを証明する。
// UpsertIfAbsent は ON CONFLICT DO NOTHING に相当し、既存行は変更されない。
func TestMasterPersonaGenerationCutoverExistingIdentityIsNotOverwritten(t *testing.T) {
	now := func() time.Time { return time.Date(2026, 4, 21, 12, 0, 0, 0, time.UTC) }
	repo := repository.NewInMemoryMasterPersonaRepository(repository.DefaultMasterPersonaSeed(now()))
	service := NewMasterPersonaGenerationService(
		repo, repo, repo, repository.NewInMemorySecretStore(), now, false,
		WithMasterPersonaBodyGenerator(&stubMasterPersonaBodyGenerator{body: "cutover-generated-body"}),
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
      "form_id": "FE04A001",
      "record_type": "NPC_",
      "editor_id": "FP_CutoverNew",
      "display_name": "Cutover New",
      "dialogues": ["hello", "world"]
    }
  ]
}`)

	// Act
	runResult, err := service.Execute(context.Background(), fixturePath, MasterPersonaAISettings{Provider: "gemini", Model: "gemini-2.5-pro", APIKey: "transport-key"})

	// Assert: existing entry is skipped and new entry is created
	if err != nil {
		t.Fatalf("expected execute to succeed: %v", err)
	}
	if runResult.ExistingSkipCount != 1 || runResult.SuccessCount != 1 {
		t.Fatalf("expected ExistingSkipCount=1 SuccessCount=1, got %#v", runResult)
	}

	// Assert: existing persona identity fields are not overwritten
	existingKey := repository.BuildMasterPersonaIdentityKey("FollowersPlus.esp", "FE01A812", "NPC_")
	existing, fetchErr := repo.GetByIdentityKey(context.Background(), existingKey)
	if fetchErr != nil {
		t.Fatalf("expected existing entry to remain readable: %v", fetchErr)
	}
	if existing.PersonaBody != "口調は丁寧語へ寄せず、中性的な温度を保つ。会話の主導権は急いで取らず、相手の出方を見てから短く返す。" {
		t.Fatalf("expected existing persona body to be unchanged after generation, got %q", existing.PersonaBody)
	}
	if existing.DisplayName != "Lys Maren" {
		t.Fatalf("expected existing display name to be unchanged, got %q", existing.DisplayName)
	}
}

// persona-generation-cutover: 生成成功時は NPC_PROFILE フィールド (identity) と PERSONA フィールド (AI 本文) の
// 両方が同一レコードに書き込まれることを証明する。body 生成完了後にのみ UpsertIfAbsent が呼ばれる設計を確認する。
func TestMasterPersonaGenerationCutoverSuccessWritesCanonicalNPCProfileAndPersona(t *testing.T) {
	now := func() time.Time { return time.Date(2026, 4, 21, 12, 0, 0, 0, time.UTC) }
	repo := repository.NewInMemoryMasterPersonaRepository(nil)
	const wantPersonaBody = "persona-generation-cutover canonical persona body"
	service := NewMasterPersonaGenerationService(
		repo, repo, repo, repository.NewInMemorySecretStore(), now, false,
		WithMasterPersonaBodyGenerator(&stubMasterPersonaBodyGenerator{body: wantPersonaBody}),
	)
	fixturePath := writeMasterPersonaExtractFixture(t, `{
  "target_plugin": "PersonaGenCutover.esp",
  "npcs": [
    {
      "form_id": "FE04A002",
      "record_type": "NPC_",
      "editor_id": "FP_CutoverCanonical",
      "display_name": "Cutover Canonical",
      "dialogues": ["dialogue one", "dialogue two"]
    }
  ]
}`)

	// Act
	runResult, err := service.Execute(context.Background(), fixturePath, MasterPersonaAISettings{Provider: "gemini", Model: "gemini-2.5-pro", APIKey: "transport-key"})

	// Assert: run completed with one success
	if err != nil {
		t.Fatalf("expected execute to succeed: %v", err)
	}
	if runResult.RunState != MasterPersonaStatusCompleted || runResult.SuccessCount != 1 {
		t.Fatalf("expected Completed with SuccessCount=1, got %#v", runResult)
	}

	// Assert: canonical entry has both NPC_PROFILE identity fields AND PERSONA body fields populated
	key := repository.BuildMasterPersonaIdentityKey("PersonaGenCutover.esp", "FE04A002", "NPC_")
	entry, fetchErr := repo.GetByIdentityKey(context.Background(), key)
	if fetchErr != nil {
		t.Fatalf("expected generated entry to be readable after complete AI output: %v", fetchErr)
	}

	// NPC_PROFILE fields (identity)
	if entry.TargetPlugin != "PersonaGenCutover.esp" {
		t.Fatalf("expected TargetPlugin to be set, got %q", entry.TargetPlugin)
	}
	if entry.FormID != "FE04A002" {
		t.Fatalf("expected FormID to be set, got %q", entry.FormID)
	}
	if entry.RecordType != "NPC_" {
		t.Fatalf("expected RecordType to be set, got %q", entry.RecordType)
	}
	if entry.EditorID != "FP_CutoverCanonical" {
		t.Fatalf("expected EditorID to be set, got %q", entry.EditorID)
	}

	// PERSONA fields (AI output)
	if entry.PersonaBody != wantPersonaBody {
		t.Fatalf("expected PersonaBody from AI output, got %q", entry.PersonaBody)
	}
	if entry.PersonaSummary == "" {
		t.Fatalf("expected PersonaSummary to be derived and non-empty after generation")
	}
}

// persona-generation-cutover: AI body 生成が失敗した場合は partial row を残さないことを証明する。
// generatePersonaBody が error を返すと UpsertIfAbsent は呼ばれず、entry は repo に存在しない。
func TestMasterPersonaGenerationCutoverBodyFailureLeavesNoPartialRow(t *testing.T) {
	now := func() time.Time { return time.Date(2026, 4, 21, 12, 0, 0, 0, time.UTC) }
	repo := repository.NewInMemoryMasterPersonaRepository(nil)
	bodyErr := errors.New("ai body generation failed")
	service := NewMasterPersonaGenerationService(
		repo, repo, repo, repository.NewInMemorySecretStore(), now, false,
		WithMasterPersonaBodyGenerator(&stubMasterPersonaBodyGenerator{err: bodyErr}),
	)
	fixturePath := writeMasterPersonaExtractFixture(t, `{
  "target_plugin": "PersonaGenCutover.esp",
  "npcs": [
    {
      "form_id": "FE04A003",
      "record_type": "NPC_",
      "editor_id": "FP_CutoverFailure",
      "display_name": "Cutover Failure",
      "dialogues": ["hello"]
    }
  ]
}`)

	// Act
	runResult, err := service.Execute(context.Background(), fixturePath, MasterPersonaAISettings{Provider: "gemini", Model: "gemini-2.5-pro", APIKey: "transport-key"})

	// Assert: execute propagates the body generation error
	if err == nil {
		t.Fatalf("expected execute to return an error when body generation fails")
	}
	if runResult.RunState != MasterPersonaStatusFailed {
		t.Fatalf("expected RunState=失敗 after body failure, got %q", runResult.RunState)
	}

	// Assert: no partial row is left in the repository
	key := repository.BuildMasterPersonaIdentityKey("PersonaGenCutover.esp", "FE04A003", "NPC_")
	_, fetchErr := repo.GetByIdentityKey(context.Background(), key)
	if !errors.Is(fetchErr, repository.ErrMasterPersonaEntryNotFound) {
		t.Fatalf("expected no partial row after body failure, got %v", fetchErr)
	}
}

// persona-generation-cutover: run state は in-memory かつ非永続であることを証明する。
// 別インスタンスの InMemoryMasterPersonaRepository は常にデフォルト値 ("入力待ち") から始まる。
func TestMasterPersonaGenerationCutoverRunStateIsInMemoryAndNonPersistent(t *testing.T) {
	now := func() time.Time { return time.Date(2026, 4, 21, 12, 0, 0, 0, time.UTC) }

	// Arrange: repo1 に "完了" の run status を書き込む (1 セッション完了を模擬)
	repo1 := repository.NewInMemoryMasterPersonaRepository(nil)
	finishedAt := now()
	if err := repo1.SaveRunStatus(context.Background(), repository.MasterPersonaRunStatusRecord{
		RunState:       MasterPersonaStatusCompleted,
		SuccessCount:   3,
		ProcessedCount: 3,
		FinishedAt:     &finishedAt,
	}); err != nil {
		t.Fatalf("expected run status save on repo1 to succeed: %v", err)
	}
	loaded1, err := repo1.LoadRunStatus(context.Background())
	if err != nil {
		t.Fatalf("expected run status load from repo1 to succeed: %v", err)
	}
	if loaded1.RunState != MasterPersonaStatusCompleted {
		t.Fatalf("expected repo1 to hold Completed run state, got %q", loaded1.RunState)
	}

	// Act: 別の InMemoryMasterPersonaRepository インスタンスを作成する (再起動を模擬)
	repo2 := repository.NewInMemoryMasterPersonaRepository(nil)
	loaded2, err := repo2.LoadRunStatus(context.Background())

	// Assert: repo2 は repo1 の run state を継承せず "入力待ち" から始まる
	if err != nil {
		t.Fatalf("expected run status load from repo2 to succeed: %v", err)
	}
	if loaded2.RunState == MasterPersonaStatusCompleted {
		t.Fatalf("expected run state to be non-persistent (repo2 must not carry repo1 state), got %q", loaded2.RunState)
	}
	if loaded2.RunState != "入力待ち" {
		t.Fatalf("expected fresh repo to start with default run state '入力待ち', got %q", loaded2.RunState)
	}
}

// persona-edit-delete-cutover: UpdateEntry が identity linkage を fetched entry から解決することを証明する。
// service が FormID / EditorID / VoiceType / ClassName / SourcePlugin を input から適用する場合は失敗する。
func TestMasterPersonaGenerationServicePersonaEditDeleteCutoverUpdateUsesIdentityLinkageFromFetchedEntry(t *testing.T) {
	now := func() time.Time { return time.Date(2026, 4, 15, 12, 0, 0, 0, time.UTC) }
	repo := repository.NewInMemoryMasterPersonaRepository(repository.DefaultMasterPersonaSeed(now()))
	svc := NewMasterPersonaGenerationService(
		repo, repo, repo,
		repository.NewInMemorySecretStore(),
		now, false,
		WithMasterPersonaBodyGenerator(newTestSafeMasterPersonaBodyGenerator()),
	)

	identityKey := repository.BuildMasterPersonaIdentityKey("FollowersPlus.esp", "FE01A812", "NPC_")

	// Arrange: identity fields に fetched entry と異なる値を渡す
	updated, err := svc.UpdateEntry(context.Background(), identityKey, MasterPersonaUpdateInput{
		FormID:       "XX999999",
		EditorID:     "FP_Fake",
		VoiceType:    "MaleGuard",
		ClassName:    "SomeOtherClass",
		SourcePlugin: "OtherMod.esp",
		PersonaBody:  "updated body",
	})

	// Assert: 結果の identity fields は fetched entry の値であること
	if err != nil {
		t.Fatalf("expected update to succeed: %v", err)
	}
	if updated.FormID != "FE01A812" {
		t.Fatalf("expected FormID preserved from fetched entry, got %q", updated.FormID)
	}
	if updated.EditorID != "FP_LysMaren" {
		t.Fatalf("expected EditorID preserved from fetched entry, got %q", updated.EditorID)
	}
	if updated.VoiceType != "FemaleYoungEager" {
		t.Fatalf("expected VoiceType preserved from fetched entry, got %q", updated.VoiceType)
	}
	if updated.ClassName != "FPScoutClass" {
		t.Fatalf("expected ClassName preserved from fetched entry, got %q", updated.ClassName)
	}
	if updated.SourcePlugin != "FollowersPlus.esp" {
		t.Fatalf("expected SourcePlugin preserved from fetched entry, got %q", updated.SourcePlugin)
	}
}

// persona-edit-delete-cutover: UpdateEntry が PersonaBody を input から書き込むことを証明する。
// service が input.PersonaBody を無視する場合は失敗する。
func TestMasterPersonaGenerationServicePersonaEditDeleteCutoverUpdateWritesPersonaBodyFromInput(t *testing.T) {
	now := func() time.Time { return time.Date(2026, 4, 15, 12, 0, 0, 0, time.UTC) }
	repo := repository.NewInMemoryMasterPersonaRepository(repository.DefaultMasterPersonaSeed(now()))
	svc := NewMasterPersonaGenerationService(
		repo, repo, repo,
		repository.NewInMemorySecretStore(),
		now, false,
		WithMasterPersonaBodyGenerator(newTestSafeMasterPersonaBodyGenerator()),
	)

	identityKey := repository.BuildMasterPersonaIdentityKey("FollowersPlus.esp", "FE01A812", "NPC_")

	updated, err := svc.UpdateEntry(context.Background(), identityKey, MasterPersonaUpdateInput{
		DisplayName: "Lys Maren",
		PersonaBody: "新しいペルソナ本文",
	})

	if err != nil {
		t.Fatalf("expected update to succeed: %v", err)
	}
	if updated.PersonaBody != "新しいペルソナ本文" {
		t.Fatalf("expected PersonaBody from input to be stored, got %q", updated.PersonaBody)
	}
}

// persona-edit-delete-cutover: DeleteEntry がエントリをリポジトリから削除することを証明する。
// DeleteEntry が実際に削除しない場合は失敗する。
func TestMasterPersonaGenerationServicePersonaEditDeleteCutoverDeleteRemovesEntry(t *testing.T) {
	now := func() time.Time { return time.Date(2026, 4, 15, 12, 0, 0, 0, time.UTC) }
	repo := repository.NewInMemoryMasterPersonaRepository(repository.DefaultMasterPersonaSeed(now()))
	svc := NewMasterPersonaGenerationService(
		repo, repo, repo,
		repository.NewInMemorySecretStore(),
		now, false,
		WithMasterPersonaBodyGenerator(newTestSafeMasterPersonaBodyGenerator()),
	)

	identityKey := repository.BuildMasterPersonaIdentityKey("FollowersPlus.esp", "FE01A812", "NPC_")

	err := svc.DeleteEntry(context.Background(), identityKey)

	if err != nil {
		t.Fatalf("expected delete to succeed: %v", err)
	}
	_, lookupErr := repo.GetByIdentityKey(context.Background(), identityKey)
	if !errors.Is(lookupErr, repository.ErrMasterPersonaEntryNotFound) {
		t.Fatalf("expected entry to be absent after delete, got: %v", lookupErr)
	}
}

// persona-edit-delete-cutover: DeleteEntry が空の identityKey を validation error で拒否することを証明する。
// validation が存在しない場合は失敗する。
func TestMasterPersonaGenerationServicePersonaEditDeleteCutoverDeleteRejectsEmptyIdentityKey(t *testing.T) {
	now := func() time.Time { return time.Date(2026, 4, 15, 12, 0, 0, 0, time.UTC) }
	repo := repository.NewInMemoryMasterPersonaRepository(nil)
	svc := NewMasterPersonaGenerationService(
		repo, repo, repo,
		repository.NewInMemorySecretStore(),
		now, false,
		WithMasterPersonaBodyGenerator(newTestSafeMasterPersonaBodyGenerator()),
	)

	err := svc.DeleteEntry(context.Background(), "")

	if !errors.Is(err, ErrMasterPersonaValidation) {
		t.Fatalf("expected validation error for empty identity key, got: %v", err)
	}
}
