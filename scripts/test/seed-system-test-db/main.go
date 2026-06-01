package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"aitranslationenginejp/internal/repository"
)

var fixedImportedAt = time.Date(2026, 5, 25, 9, 0, 0, 0, time.UTC)

func main() {
	databasePath := flag.String("db", "", "system-test sqlite database path")
	reset := flag.Bool("reset", false, "reset the system-test sqlite database before seeding")
	flag.Parse()

	if err := run(*databasePath, *reset); err != nil {
		fmt.Fprintf(os.Stderr, "seed system-test db: %v\n", err)
		os.Exit(1)
	}
}

func run(databasePath string, reset bool) error {
	trimmedPath := strings.TrimSpace(databasePath)
	if trimmedPath == "" {
		return fmt.Errorf("database path is required")
	}
	resolvedPath, err := filepath.Abs(trimmedPath)
	if err != nil {
		return fmt.Errorf("resolve database path: %w", err)
	}
	if resolvedPath == "" || filepath.Base(resolvedPath) == "." {
		return fmt.Errorf("database path is required")
	}
	if !strings.Contains(filepath.ToSlash(resolvedPath), "/test-results/") {
		return fmt.Errorf("database path must be under test-results")
	}
	if reset {
		if err := resetSQLiteFiles(resolvedPath); err != nil {
			return err
		}
	}

	ctx := context.Background()
	db, err := repository.OpenSQLiteDatabase(ctx, resolvedPath)
	if err != nil {
		return fmt.Errorf("open sqlite database: %w", err)
	}
	defer db.Close()

	sourceRepo := repository.NewSQLiteTranslationSourceRepository(db)
	jobRepo := repository.NewSQLiteJobLifecycleRepository(db)

	if _, err := createScenarioJob(ctx, sourceRepo, jobRepo, scenarioJobDraft{
		Suffix:          "running",
		JobState:        "running",
		JobProgress:     48,
		PhaseType:       "body_translation",
		PhaseState:      "running",
		PhaseProgress:   48,
		RuntimePhaseID:  "text_translation",
		RuntimeProvider: "system-test-provider",
		RuntimeModel:    "system-test-model",
	}); err != nil {
		return err
	}
	if _, err := createScenarioJob(ctx, sourceRepo, jobRepo, scenarioJobDraft{
		Suffix:          "ready",
		JobState:        "ready",
		JobProgress:     0,
		RuntimePhaseID:  "word_translation",
		RuntimeProvider: "system-test-provider",
		RuntimeModel:    "system-test-model",
	}); err != nil {
		return err
	}
	if _, err := createScenarioJob(ctx, sourceRepo, jobRepo, scenarioJobDraft{
		Suffix:      "completed",
		JobState:    "completed",
		JobProgress: 100,
	}); err != nil {
		return err
	}
	if _, err := createScenarioJob(ctx, sourceRepo, jobRepo, scenarioJobDraft{
		Suffix:          "paused",
		JobState:        "paused",
		JobProgress:     52,
		PhaseType:       "body_translation",
		PhaseState:      "paused",
		PhaseProgress:   52,
		RuntimePhaseID:  "text_translation",
		RuntimeProvider: "system-test-provider",
		RuntimeModel:    "system-test-model",
	}); err != nil {
		return err
	}
	if _, err := createScenarioJob(ctx, sourceRepo, jobRepo, scenarioJobDraft{
		Suffix:          "failed",
		JobState:        "failed",
		JobProgress:     68,
		PhaseType:       "body_translation",
		PhaseState:      "failed",
		PhaseProgress:   68,
		RuntimePhaseID:  "text_translation",
		RuntimeProvider: "system-test-provider",
		RuntimeModel:    "system-test-model",
	}); err != nil {
		return err
	}
	if _, err := createScenarioJob(ctx, sourceRepo, jobRepo, scenarioJobDraft{
		Suffix:          "canceled",
		JobState:        "canceled",
		JobProgress:     12,
		PhaseType:       "body_translation",
		PhaseState:      "canceled",
		PhaseProgress:   12,
		RuntimePhaseID:  "text_translation",
		RuntimeProvider: "system-test-provider",
		RuntimeModel:    "system-test-model",
	}); err != nil {
		return err
	}
	if _, err := createScenarioJob(ctx, sourceRepo, jobRepo, scenarioJobDraft{
		Suffix:          "term",
		JobState:        "ready",
		JobProgress:     0,
		PhaseType:       "term_translation",
		PhaseState:      "pending",
		PhaseProgress:   0,
		RuntimePhaseID:  "word_translation",
		RuntimeProvider: "system-test-provider",
		RuntimeModel:    "system-test-model",
	}); err != nil {
		return err
	}
	if _, err := createScenarioJob(ctx, sourceRepo, jobRepo, scenarioJobDraft{
		Suffix:          "persona",
		JobState:        "ready",
		JobProgress:     0,
		PhaseType:       "persona_generation",
		PhaseState:      "pending",
		PhaseProgress:   0,
		RuntimePhaseID:  "npc_persona_generation",
		RuntimeProvider: "system-test-provider",
		RuntimeModel:    "system-test-model",
	}); err != nil {
		return err
	}
	if _, err := createScenarioJob(ctx, sourceRepo, jobRepo, scenarioJobDraft{
		Suffix:          "body-pending",
		JobState:        "ready",
		JobProgress:     0,
		PhaseType:       "body_translation",
		PhaseState:      "pending",
		PhaseProgress:   0,
		RuntimePhaseID:  "text_translation",
		RuntimeProvider: "system-test-provider",
		RuntimeModel:    "system-test-model",
	}); err != nil {
		return err
	}
	if _, err := createScenarioJob(ctx, sourceRepo, jobRepo, scenarioJobDraft{
		Suffix:          "completed-term",
		JobState:        "ready",
		JobProgress:     100,
		PhaseType:       "term_translation",
		PhaseState:      "completed",
		PhaseProgress:   100,
		RuntimePhaseID:  "word_translation",
		RuntimeProvider: "system-test-provider",
		RuntimeModel:    "system-test-model",
	}); err != nil {
		return err
	}
	if _, err := createScenarioJob(ctx, sourceRepo, jobRepo, scenarioJobDraft{
		Suffix:          "body-ready-for-completion",
		JobState:        "ready",
		JobProgress:     100,
		PhaseType:       "body_translation",
		PhaseState:      "completed",
		PhaseProgress:   100,
		RuntimePhaseID:  "text_translation",
		RuntimeProvider: "system-test-provider",
		RuntimeModel:    "system-test-model",
	}); err != nil {
		return err
	}
	if _, err := createScenarioJob(ctx, sourceRepo, jobRepo, scenarioJobDraft{
		Suffix:          "body-failed",
		JobState:        "failed",
		JobProgress:     68,
		PhaseType:       "body_translation",
		PhaseState:      "failed",
		PhaseProgress:   68,
		RuntimePhaseID:  "text_translation",
		RuntimeProvider: "system-test-provider",
		RuntimeModel:    "system-test-model",
	}); err != nil {
		return err
	}
	if _, err := createScenarioJob(ctx, sourceRepo, jobRepo, scenarioJobDraft{
		Suffix:          "body-running",
		JobState:        "running",
		JobProgress:     48,
		PhaseType:       "body_translation",
		PhaseState:      "running",
		PhaseProgress:   48,
		RuntimePhaseID:  "text_translation",
		RuntimeProvider: "system-test-provider",
		RuntimeModel:    "system-test-model",
	}); err != nil {
		return err
	}
	return nil
}

func resetSQLiteFiles(databasePath string) error {
	for _, path := range []string{databasePath, databasePath + "-shm", databasePath + "-wal"} {
		if err := os.Remove(path); err != nil && !os.IsNotExist(err) {
			return fmt.Errorf("remove %s: %w", path, err)
		}
	}
	return nil
}

type scenarioJobDraft struct {
	Suffix          string
	JobState        string
	JobProgress     int
	PhaseType       string
	PhaseState      string
	PhaseProgress   int
	RuntimePhaseID  string
	RuntimeProvider string
	RuntimeModel    string
}

func createScenarioJob(
	ctx context.Context,
	sourceRepo repository.TranslationSourceRepository,
	jobRepo repository.JobLifecycleRepository,
	draft scenarioJobDraft,
) (repository.TranslationJob, error) {
	xEdit, err := sourceRepo.CreateXEditExtractedData(ctx, repository.XEditExtractedDataDraft{
		SourceFilePath:   fmt.Sprintf("system-test-%s.esp", draft.Suffix),
		SourceTool:       "xEdit",
		TargetPluginName: fmt.Sprintf("SystemTest%s.esm", draft.Suffix),
		TargetPluginType: "ESM",
		RecordCount:      1,
		ImportedAt:       fixedImportedAt,
	})
	if err != nil {
		return repository.TranslationJob{}, fmt.Errorf("create input %s: %w", draft.Suffix, err)
	}
	if _, err := sourceRepo.CreateTranslationRecord(ctx, repository.TranslationRecordDraft{
		XEditExtractedDataID: xEdit.ID,
		FormID:               fmt.Sprintf("FORM-%s", draft.Suffix),
		EditorID:             fmt.Sprintf("NPC_%s", draft.Suffix),
		RecordType:           "NPC_",
	}); err != nil {
		return repository.TranslationJob{}, fmt.Errorf("create translation record %s: %w", draft.Suffix, err)
	}

	job, err := jobRepo.CreateTranslationJob(ctx, repository.TranslationJobDraft{
		XEditExtractedDataID: xEdit.ID,
		JobName:              fmt.Sprintf("system-test-%s", draft.Suffix),
		State:                draft.JobState,
		ProgressPercent:      draft.JobProgress,
	})
	if err != nil {
		return repository.TranslationJob{}, fmt.Errorf("create job %s: %w", draft.Suffix, err)
	}

	if strings.TrimSpace(draft.PhaseType) != "" {
		phaseRun, err := jobRepo.CreateJobPhaseRun(ctx, repository.JobPhaseRunDraft{
			TranslationJobID: job.ID,
			PhaseType:        draft.PhaseType,
			State:            draft.PhaseState,
			ExecutionOrder:   1,
			AIProvider:       draft.RuntimeProvider,
			ModelName:        draft.RuntimeModel,
			ExecutionMode:    "batch",
			CredentialRef:    "system-test-secret",
			InstructionKind:  "default",
		})
		if err != nil {
			return repository.TranslationJob{}, fmt.Errorf("create phase run %s: %w", draft.Suffix, err)
		}
		if _, err := jobRepo.UpdateJobPhaseRun(ctx, phaseRun.ID, repository.JobPhaseRunUpdateDraft{
			State:           draft.PhaseState,
			ProgressPercent: draft.PhaseProgress,
			AIProvider:      draft.RuntimeProvider,
			ModelName:       draft.RuntimeModel,
			ExecutionMode:   "batch",
			CredentialRef:   "system-test-secret",
			InstructionKind: "default",
		}); err != nil {
			return repository.TranslationJob{}, fmt.Errorf("update phase run %s: %w", draft.Suffix, err)
		}
	}
	return job, nil
}
