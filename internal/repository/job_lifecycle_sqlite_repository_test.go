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

func TestSQLiteJobLifecycleRepositoryDeleteNonRunningTranslationJobTreatsIdleReadyAsSafe(t *testing.T) {
	repo := newSQLiteJobLifecycleRepositoryForTest(t, filepath.Join(t.TempDir(), "db", "job-lifecycle.sqlite3"))
	job := createTranslationJobForLifecycleTest(t, repo, 501)

	_, err := repo.CreateJobPhaseRun(context.Background(), JobPhaseRunDraft{
		TranslationJobID: job.ID,
		PhaseType:        "term_translation",
		State:            "idle_ready",
		ExecutionOrder:   1,
		AIProvider:       "openai",
		ModelName:        "gpt-5.4-mini",
		ExecutionMode:    "sync",
		CredentialRef:    "openai-primary",
		InstructionKind:  "default",
	})
	if err != nil {
		t.Fatalf("expected phase run create success: %v", err)
	}

	result, err := repo.DeleteNonRunningTranslationJob(context.Background(), job.ID)
	if err != nil {
		t.Fatalf("expected delete result success: %v", err)
	}
	if result.Outcome != TranslationJobDeleteOutcomeDeleted {
		t.Fatalf("expected delete outcome deleted for idle_ready row, got %#v", result)
	}
}

func TestSQLiteJobLifecycleRepositoryUpdateJobPhaseRunWhenStateSucceedsWhenExpectedStateMatches(t *testing.T) {
	repo := newSQLiteJobLifecycleRepositoryForTest(t, filepath.Join(t.TempDir(), "db", "job-lifecycle.sqlite3"))
	job := createTranslationJobForLifecycleTest(t, repo, 503)

	phaseRun, err := repo.CreateJobPhaseRun(context.Background(), JobPhaseRunDraft{
		TranslationJobID: job.ID,
		PhaseType:        "term_translation",
		State:            "idle_ready",
		ExecutionOrder:   1,
		AIProvider:       "openai",
		ModelName:        "gpt-5.4-mini",
		ExecutionMode:    "sync",
		CredentialRef:    "openai-primary",
		InstructionKind:  "default",
	})
	if err != nil {
		t.Fatalf("expected phase run create success: %v", err)
	}

	updated, err := repo.UpdateJobPhaseRunWhenState(context.Background(), phaseRun.ID, "idle_ready", JobPhaseRunUpdateDraft{
		State:           "running",
		ProgressPercent: 0,
	})
	if err != nil {
		t.Fatalf("expected update success when expected state matches: %v", err)
	}
	if updated.State != "running" {
		t.Fatalf("expected state running, got %#v", updated)
	}
}

func TestSQLiteJobLifecycleRepositoryUpdateJobPhaseRunWhenStateReturnsConflictWhenExpectedStateMismatches(t *testing.T) {
	repo := newSQLiteJobLifecycleRepositoryForTest(t, filepath.Join(t.TempDir(), "db", "job-lifecycle.sqlite3"))
	job := createTranslationJobForLifecycleTest(t, repo, 504)

	phaseRun, err := repo.CreateJobPhaseRun(context.Background(), JobPhaseRunDraft{
		TranslationJobID: job.ID,
		PhaseType:        "term_translation",
		State:            "idle_ready",
		ExecutionOrder:   1,
		AIProvider:       "openai",
		ModelName:        "gpt-5.4-mini",
		ExecutionMode:    "sync",
		CredentialRef:    "openai-primary",
		InstructionKind:  "default",
	})
	if err != nil {
		t.Fatalf("expected phase run create success: %v", err)
	}
	_, err = repo.UpdateJobPhaseRun(context.Background(), phaseRun.ID, JobPhaseRunUpdateDraft{
		State:           "running",
		ProgressPercent: 0,
	})
	if err != nil {
		t.Fatalf("prepare running state failed: %v", err)
	}

	_, err = repo.UpdateJobPhaseRunWhenState(context.Background(), phaseRun.ID, "idle_ready", JobPhaseRunUpdateDraft{
		State:           "completed",
		ProgressPercent: 100,
	})
	if !errors.Is(err, ErrConflict) {
		t.Fatalf("expected ErrConflict for expected state mismatch, got %v", err)
	}

	reloaded, err := repo.GetJobPhaseRunByID(context.Background(), phaseRun.ID)
	if err != nil {
		t.Fatalf("reload phase run failed: %v", err)
	}
	if reloaded.State != "running" {
		t.Fatalf("expected state unchanged by conflicted update, got %#v", reloaded)
	}
}

func TestSQLiteJobLifecycleRepositoryDeleteNonRunningTranslationJobAllowsPendingRun(t *testing.T) {
	repo := newSQLiteJobLifecycleRepositoryForTest(t, filepath.Join(t.TempDir(), "db", "job-lifecycle.sqlite3"))
	job := createTranslationJobForLifecycleTest(t, repo, 502)

	_, err := repo.CreateJobPhaseRun(context.Background(), JobPhaseRunDraft{
		TranslationJobID: job.ID,
		PhaseType:        "term_translation",
		State:            "pending",
		ExecutionOrder:   1,
		AIProvider:       "openai",
		ModelName:        "gpt-5.4-mini",
		ExecutionMode:    "sync",
		CredentialRef:    "openai-primary",
		InstructionKind:  "default",
	})
	if err != nil {
		t.Fatalf("expected phase run create success: %v", err)
	}

	result, err := repo.DeleteNonRunningTranslationJob(context.Background(), job.ID)
	if err != nil {
		t.Fatalf("expected delete result success: %v", err)
	}
	if result.Outcome != TranslationJobDeleteOutcomeDeleted {
		t.Fatalf("expected pending row to be deleted, got %#v", result)
	}
}

func createTranslationJobForLifecycleTest(t *testing.T, repo *SQLiteJobLifecycleRepository, xEditExtractedDataID int64) TranslationJob {
	t.Helper()
	_, seedErr := repo.db.ExecContext(
		context.Background(),
		`INSERT INTO X_EDIT_EXTRACTED_DATA (id, source_file_path, source_content_hash, source_tool, target_plugin_name, target_plugin_type, record_count, imported_at)
		 VALUES (?, ?, ?, ?, ?, ?, ?, ?)`,
		xEditExtractedDataID,
		"/tmp/test.json",
		"hash",
		"xedit",
		"Skyrim.esm",
		"esm",
		0,
		time.Date(2026, 5, 10, 0, 0, 0, 0, time.UTC).Format(time.RFC3339),
	)
	if seedErr != nil {
		t.Fatalf("expected x_edit_extracted_data seed success: %v", seedErr)
	}
	job, err := repo.CreateTranslationJob(context.Background(), TranslationJobDraft{
		XEditExtractedDataID: xEditExtractedDataID,
		JobName:              filepath.Base(t.TempDir()),
		State:                "ready",
		ProgressPercent:      0,
	})
	if err != nil {
		t.Fatalf("expected translation job create success: %v", err)
	}
	return job
}

func newSQLiteJobLifecycleRepositoryForTest(t *testing.T, databasePath string) *SQLiteJobLifecycleRepository {
	t.Helper()
	db := openSQLiteJobLifecycleDatabase(t, databasePath)
	repo, ok := NewSQLiteJobLifecycleRepository(db).(*SQLiteJobLifecycleRepository)
	if !ok {
		t.Fatal("expected sqlite job lifecycle repository concrete type")
	}
	return repo
}

func openSQLiteJobLifecycleDatabase(t *testing.T, databasePath string) *sqlx.DB {
	t.Helper()
	db, err := sqliteinfra.OpenMasterDictionaryDatabase(context.Background(), databasePath, nil)
	if err != nil {
		t.Fatalf("expected sqlite database open to succeed: %v", err)
	}
	t.Cleanup(func() {
		_ = db.Close()
	})
	return db
}
