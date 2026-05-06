package integrationtest

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"testing"

	"aitranslationenginejp/internal/repository"

	"github.com/jmoiron/sqlx"
)

type translationJobManagementScenarioHarness struct {
	ctx        context.Context
	db         *sqlx.DB
	sourceRepo *repository.SQLiteTranslationSourceRepository
	jobRepo    *repository.SQLiteJobLifecycleRepository
}

func newTranslationJobManagementScenarioHarness(t *testing.T) translationJobManagementScenarioHarness {
	t.Helper()
	ctx := context.Background()
	db := openIntegrationDB(t)
	sourceRepo := repository.NewSQLiteTranslationSourceRepository(db)
	jobRepo := repository.NewSQLiteJobLifecycleRepository(db)
	sqliteSourceRepo, ok := sourceRepo.(*repository.SQLiteTranslationSourceRepository)
	if !ok {
		t.Fatal("expected SQLiteTranslationSourceRepository")
	}
	sqliteJobRepo, ok := jobRepo.(*repository.SQLiteJobLifecycleRepository)
	if !ok {
		t.Fatal("expected SQLiteJobLifecycleRepository")
	}
	return translationJobManagementScenarioHarness{
		ctx:        ctx,
		db:         db,
		sourceRepo: sqliteSourceRepo,
		jobRepo:    sqliteJobRepo,
	}
}

func (h translationJobManagementScenarioHarness) createScenarioJob(
	t *testing.T,
	suffix string,
	jobState string,
	phaseState string,
	phaseProgress int,
	latestError string,
) repository.TranslationJob {
	t.Helper()
	xEdit, err := h.sourceRepo.CreateXEditExtractedData(h.ctx, repository.XEditExtractedDataDraft{
		SourceFilePath:   fmt.Sprintf("job-management-%s.esp", suffix),
		SourceTool:       "xEdit",
		TargetPluginName: fmt.Sprintf("JobManagement%s.esm", suffix),
		TargetPluginType: "ESM",
		RecordCount:      3,
		ImportedAt:       fixedNow,
	})
	if err != nil {
		t.Fatalf("CreateXEditExtractedData failed: %v", err)
	}
	_, err = h.sourceRepo.CreateTranslationRecord(h.ctx, repository.TranslationRecordDraft{
		XEditExtractedDataID: xEdit.ID,
		FormID:               fmt.Sprintf("FORM-%s", suffix),
		EditorID:             fmt.Sprintf("NPC_%s", suffix),
		RecordType:           "NPC_",
	})
	if err != nil {
		t.Fatalf("CreateTranslationRecord failed: %v", err)
	}

	job, err := h.jobRepo.CreateTranslationJob(h.ctx, repository.TranslationJobDraft{
		XEditExtractedDataID: xEdit.ID,
		JobName:              fmt.Sprintf("job-management-%s", suffix),
		State:                jobState,
		ProgressPercent:      phaseProgress,
	})
	if err != nil {
		t.Fatalf("CreateTranslationJob failed: %v", err)
	}

	phaseRun, err := h.jobRepo.CreateJobPhaseRun(h.ctx, repository.JobPhaseRunDraft{
		TranslationJobID: job.ID,
		PhaseType:        "body_translation",
		State:            phaseState,
		ExecutionOrder:   1,
		AIProvider:       "scenario-provider",
		ModelName:        "scenario-model",
		ExecutionMode:    "batch",
		CredentialRef:    "credential-ref",
		InstructionKind:  "default",
	})
	if err != nil {
		t.Fatalf("CreateJobPhaseRun failed: %v", err)
	}
	_, err = h.jobRepo.UpdateJobPhaseRun(h.ctx, phaseRun.ID, repository.JobPhaseRunUpdateDraft{
		State:           phaseState,
		ProgressPercent: phaseProgress,
		LatestError:     latestError,
	})
	if err != nil {
		t.Fatalf("UpdateJobPhaseRun failed: %v", err)
	}
	_, err = h.jobRepo.SaveTranslationJobPhaseRuntimeSnapshot(h.ctx, repository.TranslationJobPhaseRuntimeSnapshotDraft{
		TranslationJobID: job.ID,
		PhaseID:          "body_translation",
		Provider:         "scenario-provider",
		ModelName:        "scenario-model",
		CredentialStatus: "configured",
		ExecutionMode:    "batch",
	})
	if err != nil {
		t.Fatalf("SaveTranslationJobPhaseRuntimeSnapshot failed: %v", err)
	}
	return job
}

func TestSCN_TJM_004_RunningDeleteIsRejectedAndKeepsInput(t *testing.T) {
	h := newTranslationJobManagementScenarioHarness(t)
	job := h.createScenarioJob(t, "running", "running", "running", 48, "stop_requested")

	result, err := h.jobRepo.DeleteNonRunningTranslationJob(h.ctx, job.ID)
	if err != nil {
		t.Fatalf("DeleteNonRunningTranslationJob failed: %v", err)
	}
	if result.Outcome != repository.TranslationJobDeleteOutcomeBlockedRunning {
		t.Fatalf("expected blocked running delete, got %q", result.Outcome)
	}
	phaseRuns, err := h.jobRepo.ListJobPhaseRunsByJobID(h.ctx, job.ID)
	if err != nil {
		t.Fatalf("ListJobPhaseRunsByJobID failed: %v", err)
	}
	if len(phaseRuns) != 1 || phaseRuns[0].LatestError != "stop_requested" {
		t.Fatalf("expected stop_requested phase guard, got %#v", phaseRuns)
	}
	if _, err := h.jobRepo.GetTranslationJobByID(h.ctx, job.ID); err != nil {
		t.Fatalf("expected Running job to remain, got %v", err)
	}
	if _, err := h.sourceRepo.GetXEditExtractedDataByID(h.ctx, job.XEditExtractedDataID); err != nil {
		t.Fatalf("expected input source to remain, got %v", err)
	}
}

func TestSCN_TJM_006_NonRunningDeleteRemovesJobAndKeepsInput(t *testing.T) {
	h := newTranslationJobManagementScenarioHarness(t)
	job := h.createScenarioJob(t, "paused", "paused", "paused", 44, "")

	result, err := h.jobRepo.DeleteNonRunningTranslationJob(h.ctx, job.ID)
	if err != nil {
		t.Fatalf("DeleteNonRunningTranslationJob failed: %v", err)
	}
	if result.Outcome != repository.TranslationJobDeleteOutcomeDeleted || result.Job == nil || result.Job.ID != job.ID {
		t.Fatalf("expected deleted job id %d, got %#v", job.ID, result)
	}
	if _, getErr := h.jobRepo.GetTranslationJobByID(h.ctx, job.ID); !errors.Is(getErr, repository.ErrNotFound) {
		t.Fatalf("expected deleted job to be missing, got %v", getErr)
	}
	if _, getErr := h.sourceRepo.GetXEditExtractedDataByID(h.ctx, job.XEditExtractedDataID); getErr != nil {
		t.Fatalf("expected input source to remain, got %v", getErr)
	}
	snapshots, err := h.jobRepo.ListTranslationJobPhaseRuntimeSnapshots(h.ctx, job.ID)
	if err != nil {
		t.Fatalf("ListTranslationJobPhaseRuntimeSnapshots failed: %v", err)
	}
	if len(snapshots) != 0 {
		t.Fatalf("expected job-owned runtime snapshots to be deleted, got %d", len(snapshots))
	}
}

func TestSCN_TJM_008_ProjectionFailuresReturnSafeReasons(t *testing.T) {
	h := newTranslationJobManagementScenarioHarness(t)
	job := h.createScenarioJob(t, "projection", "failed", "failed", 0, "")

	_, err := h.jobRepo.GetTranslationJobByID(h.ctx, job.ID+9999)
	if !errors.Is(err, repository.ErrNotFound) {
		t.Fatalf("expected not found for stale selection, got %v", err)
	}

	phaseRuns, err := h.jobRepo.ListJobPhaseRunsByJobID(h.ctx, job.ID)
	if err != nil {
		t.Fatalf("ListJobPhaseRunsByJobID failed: %v", err)
	}
	if len(phaseRuns) != 1 {
		t.Fatalf("expected one phase run, got %d", len(phaseRuns))
	}
	if phaseRuns[0].ProgressPercent != 0 || strings.TrimSpace(phaseRuns[0].LatestError) != "" {
		t.Fatalf("expected aggregate-risk phase fixture, got %#v", phaseRuns[0])
	}
}

func TestSCN_TJM_008_ListLoadFailureIsNotEmptyList(t *testing.T) {
	h := newTranslationJobManagementScenarioHarness(t)
	if closeErr := h.db.Close(); closeErr != nil {
		t.Fatalf("Close failed: %v", closeErr)
	}
	jobs, err := h.jobRepo.ListIncompleteTranslationJobs(h.ctx)
	if err == nil {
		t.Fatal("expected list load failure, got nil")
	}
	if len(jobs) != 0 {
		t.Fatalf("expected no successful empty-list response, got %#v", jobs)
	}
	if !strings.Contains(err.Error(), "list incomplete translation jobs") {
		t.Fatalf("expected list load failure context, got %v", err)
	}
}

func TestSCN_TJM_005_InconsistentRunningPhaseDeleteIsRejectedAndKeepsRows(t *testing.T) {
	h := newTranslationJobManagementScenarioHarness(t)
	job := h.createScenarioJob(t, "inconsistent-running", "paused", "running", 52, "")

	result, err := h.jobRepo.DeleteNonRunningTranslationJob(h.ctx, job.ID)
	if err != nil {
		t.Fatalf("DeleteNonRunningTranslationJob failed: %v", err)
	}
	if result.Outcome != repository.TranslationJobDeleteOutcomeBlockedUnsafe {
		t.Fatalf("expected blocked unsafe delete, got %q", result.Outcome)
	}
	phaseRuns, err := h.jobRepo.ListJobPhaseRunsByJobID(h.ctx, job.ID)
	if err != nil {
		t.Fatalf("ListJobPhaseRunsByJobID failed: %v", err)
	}
	if len(phaseRuns) != 1 || phaseRuns[0].State != "running" {
		t.Fatalf("expected running phase run to remain, got %#v", phaseRuns)
	}
	if _, getErr := h.jobRepo.GetTranslationJobByID(h.ctx, job.ID); getErr != nil {
		t.Fatalf("expected paused job to remain, got %v", getErr)
	}
	if _, getErr := h.sourceRepo.GetXEditExtractedDataByID(h.ctx, job.XEditExtractedDataID); getErr != nil {
		t.Fatalf("expected input source to remain, got %v", getErr)
	}
	snapshots, err := h.jobRepo.ListTranslationJobPhaseRuntimeSnapshots(h.ctx, job.ID)
	if err != nil {
		t.Fatalf("ListTranslationJobPhaseRuntimeSnapshots failed: %v", err)
	}
	if len(snapshots) != 1 {
		t.Fatalf("expected runtime snapshot to remain, got %d", len(snapshots))
	}
}
