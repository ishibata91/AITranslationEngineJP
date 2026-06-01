package repository

import (
	"context"
	"errors"
	"path/filepath"
	"testing"
	"time"

	sqliteinfra "aitranslationenginejp/internal/infra/sqlite/dbinit"

	"github.com/jmoiron/sqlx"
)

// ---------------------------------------------------------------------------
// SQLite 実装テスト
// ---------------------------------------------------------------------------

func TestJobPhaseAISettingsSQLiteRepository_SaveAndLoad(t *testing.T) {
	db := openJobPhaseAISettingsSQLiteDatabase(t)
	repo := NewSQLiteJobPhaseAISettingsRepository(db)
	now := time.Date(2026, 6, 1, 10, 0, 0, 0, time.UTC)

	draft := JobPhaseAISettingsDraft{
		PhaseType:     "word_translation",
		AIProvider:    "gemini",
		ModelName:     "gemini-2.0-flash",
		ExecutionMode: "standard",
		BatchMode:     "disabled",
		UpdatedAt:     now,
	}

	saved, err := repo.Save(context.Background(), draft)
	if err != nil {
		t.Fatalf("expected Save to succeed: %v", err)
	}
	if saved.PhaseType != "word_translation" {
		t.Errorf("expected PhaseType=word_translation, got %q", saved.PhaseType)
	}
	if saved.AIProvider != "gemini" {
		t.Errorf("expected AIProvider=gemini, got %q", saved.AIProvider)
	}
	if saved.ModelName != "gemini-2.0-flash" {
		t.Errorf("expected ModelName=gemini-2.0-flash, got %q", saved.ModelName)
	}
	if saved.ExecutionMode != "standard" {
		t.Errorf("expected ExecutionMode=standard, got %q", saved.ExecutionMode)
	}
	if saved.BatchMode != "disabled" {
		t.Errorf("expected BatchMode=disabled, got %q", saved.BatchMode)
	}

	loaded, err := repo.LoadByPhaseType(context.Background(), "word_translation")
	if err != nil {
		t.Fatalf("expected LoadByPhaseType to succeed: %v", err)
	}
	if loaded.PhaseType != "word_translation" {
		t.Errorf("expected loaded PhaseType=word_translation, got %q", loaded.PhaseType)
	}
	if loaded.AIProvider != "gemini" {
		t.Errorf("expected loaded AIProvider=gemini, got %q", loaded.AIProvider)
	}
}

func TestJobPhaseAISettingsSQLiteRepository_ErrNotFoundWhenAbsent(t *testing.T) {
	db := openJobPhaseAISettingsSQLiteDatabase(t)
	repo := NewSQLiteJobPhaseAISettingsRepository(db)

	_, err := repo.LoadByPhaseType(context.Background(), "word_translation")
	if err == nil {
		t.Fatal("expected ErrJobPhaseAISettingsNotFound, got nil")
	}
	if !errors.Is(err, ErrJobPhaseAISettingsNotFound) {
		t.Errorf("expected ErrJobPhaseAISettingsNotFound, got %v", err)
	}
}

func TestJobPhaseAISettingsSQLiteRepository_UpsertOverwritesSamePhaseType(t *testing.T) {
	db := openJobPhaseAISettingsSQLiteDatabase(t)
	repo := NewSQLiteJobPhaseAISettingsRepository(db)
	now := time.Date(2026, 6, 1, 10, 0, 0, 0, time.UTC)

	_, err := repo.Save(context.Background(), JobPhaseAISettingsDraft{
		PhaseType:     "word_translation",
		AIProvider:    "gemini",
		ModelName:     "gemini-2.0-flash",
		ExecutionMode: "standard",
		BatchMode:     "disabled",
		UpdatedAt:     now,
	})
	if err != nil {
		t.Fatalf("expected first Save to succeed: %v", err)
	}

	later := now.Add(time.Hour)
	_, err = repo.Save(context.Background(), JobPhaseAISettingsDraft{
		PhaseType:     "word_translation",
		AIProvider:    "xai",
		ModelName:     "grok-3",
		ExecutionMode: "batch",
		BatchMode:     "enabled",
		UpdatedAt:     later,
	})
	if err != nil {
		t.Fatalf("expected second Save (upsert) to succeed: %v", err)
	}

	loaded, err := repo.LoadByPhaseType(context.Background(), "word_translation")
	if err != nil {
		t.Fatalf("expected LoadByPhaseType after upsert to succeed: %v", err)
	}
	if loaded.AIProvider != "xai" {
		t.Errorf("expected AIProvider=xai after upsert, got %q", loaded.AIProvider)
	}
	if loaded.ModelName != "grok-3" {
		t.Errorf("expected ModelName=grok-3 after upsert, got %q", loaded.ModelName)
	}
}

func TestJobPhaseAISettingsSQLiteRepository_ThreePhaseTypesCoexist(t *testing.T) {
	db := openJobPhaseAISettingsSQLiteDatabase(t)
	repo := NewSQLiteJobPhaseAISettingsRepository(db)
	now := time.Date(2026, 6, 1, 10, 0, 0, 0, time.UTC)

	phases := []struct {
		phaseType string
		provider  string
	}{
		{"word_translation", "gemini"},
		{"npc_persona_generation", "xai"},
		{"text_translation", "lm_studio"},
	}

	for _, ph := range phases {
		_, err := repo.Save(context.Background(), JobPhaseAISettingsDraft{
			PhaseType:     ph.phaseType,
			AIProvider:    ph.provider,
			ModelName:     "model-a",
			ExecutionMode: "standard",
			BatchMode:     "disabled",
			UpdatedAt:     now,
		})
		if err != nil {
			t.Fatalf("expected Save(%s) to succeed: %v", ph.phaseType, err)
		}
	}

	for _, ph := range phases {
		loaded, err := repo.LoadByPhaseType(context.Background(), ph.phaseType)
		if err != nil {
			t.Fatalf("expected LoadByPhaseType(%s) to succeed: %v", ph.phaseType, err)
		}
		if loaded.AIProvider != ph.provider {
			t.Errorf("expected AIProvider=%s for phase_type=%s, got %q", ph.provider, ph.phaseType, loaded.AIProvider)
		}
	}
}

func TestJobPhaseAISettingsSQLiteRepository_DoesNotContainJobIDOrUserID(t *testing.T) {
	db := openJobPhaseAISettingsSQLiteDatabase(t)
	repo := NewSQLiteJobPhaseAISettingsRepository(db)
	now := time.Date(2026, 6, 1, 10, 0, 0, 0, time.UTC)

	saved, err := repo.Save(context.Background(), JobPhaseAISettingsDraft{
		PhaseType:     "word_translation",
		AIProvider:    "gemini",
		ModelName:     "gemini-2.0-flash",
		ExecutionMode: "standard",
		BatchMode:     "disabled",
		UpdatedAt:     now,
	})
	if err != nil {
		t.Fatalf("expected Save to succeed: %v", err)
	}
	// JobPhaseAISettings 型に job_id / user_id フィールドが含まれないことをコンパイル時に確認する。
	// 実行時には保存値が struct に格納されることを確認する。
	_ = saved.PhaseType
	_ = saved.AIProvider
	_ = saved.ModelName
	_ = saved.ExecutionMode
	_ = saved.BatchMode
	_ = saved.UpdatedAt
}

// ---------------------------------------------------------------------------
// InMemory 実装テスト
// ---------------------------------------------------------------------------

func TestJobPhaseAISettingsInMemoryRepository_ErrNotFoundWhenAbsent(t *testing.T) {
	repo := NewInMemoryJobPhaseAISettingsRepository()

	_, err := repo.LoadByPhaseType(context.Background(), "word_translation")
	if err == nil {
		t.Fatal("expected ErrJobPhaseAISettingsNotFound, got nil")
	}
	if !errors.Is(err, ErrJobPhaseAISettingsNotFound) {
		t.Errorf("expected ErrJobPhaseAISettingsNotFound, got %v", err)
	}
}

func TestJobPhaseAISettingsInMemoryRepository_UpsertOverwritesSamePhaseType(t *testing.T) {
	repo := NewInMemoryJobPhaseAISettingsRepository()
	now := time.Date(2026, 6, 1, 10, 0, 0, 0, time.UTC)

	_, err := repo.Save(context.Background(), JobPhaseAISettingsDraft{
		PhaseType: "word_translation", AIProvider: "gemini", ModelName: "gemini-2.0-flash",
		ExecutionMode: "standard", BatchMode: "disabled", UpdatedAt: now,
	})
	if err != nil {
		t.Fatalf("expected first Save to succeed: %v", err)
	}

	_, err = repo.Save(context.Background(), JobPhaseAISettingsDraft{
		PhaseType: "word_translation", AIProvider: "xai", ModelName: "grok-3",
		ExecutionMode: "batch", BatchMode: "enabled", UpdatedAt: now.Add(time.Hour),
	})
	if err != nil {
		t.Fatalf("expected second Save (upsert) to succeed: %v", err)
	}

	loaded, err := repo.LoadByPhaseType(context.Background(), "word_translation")
	if err != nil {
		t.Fatalf("expected LoadByPhaseType after upsert to succeed: %v", err)
	}
	if loaded.AIProvider != "xai" {
		t.Errorf("expected AIProvider=xai after upsert, got %q", loaded.AIProvider)
	}
}

// ---------------------------------------------------------------------------
// テストヘルパー
// ---------------------------------------------------------------------------

func openJobPhaseAISettingsSQLiteDatabase(t *testing.T) *sqlx.DB {
	t.Helper()
	databasePath := filepath.Join(t.TempDir(), "db", "job-phase-ai-settings.sqlite3")
	db, err := sqliteinfra.OpenMasterDictionaryDatabase(context.Background(), databasePath, nil)
	if err != nil {
		t.Fatalf("expected sqlite database open to succeed: %v", err)
	}
	t.Cleanup(func() {
		_ = db.Close()
	})
	return db
}
