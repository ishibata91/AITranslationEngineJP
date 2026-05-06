package service

import (
	"context"
	"errors"
	"strings"
	"testing"
	"time"

	"aitranslationenginejp/internal/repository"
)

type fakeTranslationJobManagementLifecycleRepository struct {
	jobByID                    map[int64]repository.TranslationJob
	phaseRunsByJobID           map[int64][]repository.JobPhaseRun
	snapshotsByJobID           map[int64][]repository.TranslationJobPhaseRuntimeSnapshot
	deleteResultByJobID        map[int64]repository.TranslationJobDeleteResult
	listIncompleteJobsResponse []repository.TranslationJob
	deleteCalls                []int64
}

func (repo *fakeTranslationJobManagementLifecycleRepository) CreateTranslationJob(context.Context, repository.TranslationJobDraft) (repository.TranslationJob, error) {
	return repository.TranslationJob{}, errors.New("not used")
}
func (repo *fakeTranslationJobManagementLifecycleRepository) GetTranslationJobByID(_ context.Context, id int64) (repository.TranslationJob, error) {
	if job, ok := repo.jobByID[id]; ok {
		return job, nil
	}
	return repository.TranslationJob{}, repository.ErrNotFound
}
func (repo *fakeTranslationJobManagementLifecycleRepository) UpdateTranslationJob(context.Context, int64, repository.TranslationJobUpdateDraft) (repository.TranslationJob, error) {
	return repository.TranslationJob{}, errors.New("not used")
}
func (repo *fakeTranslationJobManagementLifecycleRepository) CreateJobPhaseRun(context.Context, repository.JobPhaseRunDraft) (repository.JobPhaseRun, error) {
	return repository.JobPhaseRun{}, errors.New("not used")
}
func (repo *fakeTranslationJobManagementLifecycleRepository) GetJobPhaseRunByID(context.Context, int64) (repository.JobPhaseRun, error) {
	return repository.JobPhaseRun{}, errors.New("not used")
}
func (repo *fakeTranslationJobManagementLifecycleRepository) UpdateJobPhaseRun(context.Context, int64, repository.JobPhaseRunUpdateDraft) (repository.JobPhaseRun, error) {
	return repository.JobPhaseRun{}, errors.New("not used")
}
func (repo *fakeTranslationJobManagementLifecycleRepository) ListJobPhaseRunsByJobID(_ context.Context, jobID int64) ([]repository.JobPhaseRun, error) {
	return append([]repository.JobPhaseRun(nil), repo.phaseRunsByJobID[jobID]...), nil
}
func (repo *fakeTranslationJobManagementLifecycleRepository) FindJobPhaseRun(context.Context, int64, string) (repository.JobPhaseRun, error) {
	return repository.JobPhaseRun{}, errors.New("not used")
}
func (repo *fakeTranslationJobManagementLifecycleRepository) CreatePhaseRunTranslationField(context.Context, repository.PhaseRunTranslationFieldDraft) (repository.PhaseRunTranslationField, error) {
	return repository.PhaseRunTranslationField{}, errors.New("not used")
}
func (repo *fakeTranslationJobManagementLifecycleRepository) ListPhaseRunTranslationFieldsByPhaseRunID(context.Context, int64) ([]repository.PhaseRunTranslationField, error) {
	return nil, errors.New("not used")
}
func (repo *fakeTranslationJobManagementLifecycleRepository) CreatePhaseRunPersona(context.Context, repository.PhaseRunPersonaDraft) (repository.PhaseRunPersona, error) {
	return repository.PhaseRunPersona{}, errors.New("not used")
}
func (repo *fakeTranslationJobManagementLifecycleRepository) ListPhaseRunPersonasByPhaseRunID(context.Context, int64) ([]repository.PhaseRunPersona, error) {
	return nil, errors.New("not used")
}
func (repo *fakeTranslationJobManagementLifecycleRepository) CreatePhaseRunDictionaryEntry(context.Context, repository.PhaseRunDictionaryEntryDraft) (repository.PhaseRunDictionaryEntry, error) {
	return repository.PhaseRunDictionaryEntry{}, errors.New("not used")
}
func (repo *fakeTranslationJobManagementLifecycleRepository) ListPhaseRunDictionaryEntriesByPhaseRunID(context.Context, int64) ([]repository.PhaseRunDictionaryEntry, error) {
	return nil, errors.New("not used")
}
func (repo *fakeTranslationJobManagementLifecycleRepository) ListIncompleteTranslationJobs(context.Context) ([]repository.TranslationJob, error) {
	return append([]repository.TranslationJob(nil), repo.listIncompleteJobsResponse...), nil
}
func (repo *fakeTranslationJobManagementLifecycleRepository) ListTranslationJobPhaseRuntimeSnapshots(_ context.Context, jobID int64) ([]repository.TranslationJobPhaseRuntimeSnapshot, error) {
	return append([]repository.TranslationJobPhaseRuntimeSnapshot(nil), repo.snapshotsByJobID[jobID]...), nil
}
func (repo *fakeTranslationJobManagementLifecycleRepository) DeleteNonRunningTranslationJob(_ context.Context, jobID int64) (repository.TranslationJobDeleteResult, error) {
	repo.deleteCalls = append(repo.deleteCalls, jobID)
	if result, ok := repo.deleteResultByJobID[jobID]; ok {
		return result, nil
	}
	return repository.TranslationJobDeleteResult{Outcome: repository.TranslationJobDeleteOutcomeNotFound}, nil
}

type fakeTranslationJobManagementSourceRepository struct {
	sourceByID map[int64]repository.XEditExtractedData
}

func (repo *fakeTranslationJobManagementSourceRepository) CreateXEditExtractedData(context.Context, repository.XEditExtractedDataDraft) (repository.XEditExtractedData, error) {
	return repository.XEditExtractedData{}, errors.New("not used")
}
func (repo *fakeTranslationJobManagementSourceRepository) GetXEditExtractedDataByID(_ context.Context, id int64) (repository.XEditExtractedData, error) {
	if source, ok := repo.sourceByID[id]; ok {
		return source, nil
	}
	return repository.XEditExtractedData{}, repository.ErrNotFound
}
func (repo *fakeTranslationJobManagementSourceRepository) DeleteXEditExtractedDataByID(context.Context, int64) error {
	return errors.New("not used")
}
func (repo *fakeTranslationJobManagementSourceRepository) CreateTranslationRecord(context.Context, repository.TranslationRecordDraft) (repository.TranslationRecord, error) {
	return repository.TranslationRecord{}, errors.New("not used")
}
func (repo *fakeTranslationJobManagementSourceRepository) GetTranslationRecordByID(context.Context, int64) (repository.TranslationRecord, error) {
	return repository.TranslationRecord{}, errors.New("not used")
}
func (repo *fakeTranslationJobManagementSourceRepository) ListTranslationRecordsByXEditID(context.Context, int64) ([]repository.TranslationRecord, error) {
	return nil, errors.New("not used")
}
func (repo *fakeTranslationJobManagementSourceRepository) UpsertNpcProfile(context.Context, repository.NpcProfileDraft) (repository.NpcProfile, error) {
	return repository.NpcProfile{}, errors.New("not used")
}
func (repo *fakeTranslationJobManagementSourceRepository) GetNpcProfileByID(context.Context, int64) (repository.NpcProfile, error) {
	return repository.NpcProfile{}, errors.New("not used")
}
func (repo *fakeTranslationJobManagementSourceRepository) CreateNpcRecord(context.Context, repository.NpcRecordDraft) (repository.NpcRecord, error) {
	return repository.NpcRecord{}, errors.New("not used")
}
func (repo *fakeTranslationJobManagementSourceRepository) GetNpcRecordByTranslationRecordID(context.Context, int64) (repository.NpcRecord, error) {
	return repository.NpcRecord{}, errors.New("not used")
}
func (repo *fakeTranslationJobManagementSourceRepository) CreateTranslationField(context.Context, repository.TranslationFieldDraft) (repository.TranslationField, error) {
	return repository.TranslationField{}, errors.New("not used")
}
func (repo *fakeTranslationJobManagementSourceRepository) GetTranslationFieldByID(context.Context, int64) (repository.TranslationField, error) {
	return repository.TranslationField{}, errors.New("not used")
}
func (repo *fakeTranslationJobManagementSourceRepository) ListTranslationFieldsByTranslationRecordID(context.Context, int64) ([]repository.TranslationField, error) {
	return nil, errors.New("not used")
}
func (repo *fakeTranslationJobManagementSourceRepository) CreateTranslationFieldRecordReference(context.Context, repository.TranslationFieldRecordReferenceDraft) (repository.TranslationFieldRecordReference, error) {
	return repository.TranslationFieldRecordReference{}, errors.New("not used")
}
func (repo *fakeTranslationJobManagementSourceRepository) ListTranslationFieldRecordReferencesByFieldID(context.Context, int64) ([]repository.TranslationFieldRecordReference, error) {
	return nil, errors.New("not used")
}

type fakeTranslationJobManagementCacheReader struct {
	available bool
}

func (reader fakeTranslationJobManagementCacheReader) HasTranslationCacheByXEditID(context.Context, int64) (bool, error) {
	return reader.available, nil
}

type fakeTranslationJobManagementTransactor struct{}

func (fakeTranslationJobManagementTransactor) WithTransaction(ctx context.Context, fn func(context.Context) error) error {
	return fn(ctx)
}

func TestTranslationJobManagementServiceGetJobDetailMarksCompletedAsStaleSelection(t *testing.T) {
	repo := &fakeTranslationJobManagementLifecycleRepository{
		jobByID: map[int64]repository.TranslationJob{1: {ID: 1, State: "completed"}},
	}
	service := NewTranslationJobManagementService(repo, &fakeTranslationJobManagementSourceRepository{}, fakeTranslationJobManagementTransactor{})

	_, err := service.GetJobDetail(context.Background(), 1)
	if err == nil {
		t.Fatal("expected stale selection error")
	}
	if !strings.Contains(err.Error(), "stale_selection") {
		t.Fatalf("expected stale_selection category, got %v", err)
	}
}

func TestTranslationJobManagementServiceDeleteJobRejectsRunning(t *testing.T) {
	now := time.Date(2026, 5, 6, 10, 0, 0, 0, time.UTC)
	running := repository.TranslationJob{ID: 20, XEditExtractedDataID: 5, State: "running", CreatedAt: now}
	repo := &fakeTranslationJobManagementLifecycleRepository{
		jobByID: map[int64]repository.TranslationJob{20: running},
		snapshotsByJobID: map[int64][]repository.TranslationJobPhaseRuntimeSnapshot{20: {
			{TranslationJobID: 20, PhaseID: "text_translation", Provider: "openai", ModelName: "gpt-5", ExecutionMode: "batch", CredentialStatus: "configured"},
		}},
		deleteResultByJobID: map[int64]repository.TranslationJobDeleteResult{
			20: {Outcome: repository.TranslationJobDeleteOutcomeBlockedRunning, Job: &running},
		},
	}
	sourceRepo := &fakeTranslationJobManagementSourceRepository{sourceByID: map[int64]repository.XEditExtractedData{5: {ID: 5, SourceFilePath: "/mods/input.json", TargetPluginName: "Skyrim.esm"}}}
	service := NewTranslationJobManagementService(repo, sourceRepo, fakeTranslationJobManagementTransactor{})
	service.cacheReader = fakeTranslationJobManagementCacheReader{available: true}

	result, err := service.DeleteJob(context.Background(), 20)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if result.ReasonCategory != "running_delete_blocked" {
		t.Fatalf("expected running_delete_blocked, got %#v", result)
	}
	if result.Detail == nil || result.Detail.DeleteAvailability.Enabled {
		t.Fatalf("expected blocked detail projection, got %#v", result.Detail)
	}
}

func TestTranslationJobManagementServiceDeleteJobDeletesNonRunning(t *testing.T) {
	now := time.Date(2026, 5, 6, 10, 0, 0, 0, time.UTC)
	job := repository.TranslationJob{ID: 31, XEditExtractedDataID: 7, State: "paused", CreatedAt: now}
	repo := &fakeTranslationJobManagementLifecycleRepository{
		jobByID: map[int64]repository.TranslationJob{31: job},
		phaseRunsByJobID: map[int64][]repository.JobPhaseRun{
			31: {{ID: 310, TranslationJobID: 31, PhaseType: "body_translation", State: "paused", ProgressPercent: 44}},
		},
		snapshotsByJobID: map[int64][]repository.TranslationJobPhaseRuntimeSnapshot{31: {
			{TranslationJobID: 31, PhaseID: "body_translation", Provider: "openai", ModelName: "gpt-5", ExecutionMode: "batch", CredentialStatus: "configured"},
		}},
		deleteResultByJobID: map[int64]repository.TranslationJobDeleteResult{
			31: {Outcome: repository.TranslationJobDeleteOutcomeDeleted, Job: &job},
		},
	}
	sourceRepo := &fakeTranslationJobManagementSourceRepository{sourceByID: map[int64]repository.XEditExtractedData{7: {ID: 7, SourceFilePath: "/mods/input.json", TargetPluginName: "Skyrim.esm"}}}
	service := NewTranslationJobManagementService(repo, sourceRepo, fakeTranslationJobManagementTransactor{})
	service.cacheReader = fakeTranslationJobManagementCacheReader{available: true}

	result, err := service.DeleteJob(context.Background(), 31)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if result.DeletedJobID == nil || *result.DeletedJobID != 31 {
		t.Fatalf("expected deleted job id 31, got %#v", result)
	}
	if result.Tone != "success" {
		t.Fatalf("expected success tone, got %#v", result)
	}
}

func TestTranslationJobManagementServiceDeleteJobRejectsProjectionInconsistentRunningPhase(t *testing.T) {
	now := time.Date(2026, 5, 6, 10, 0, 0, 0, time.UTC)
	job := repository.TranslationJob{ID: 41, XEditExtractedDataID: 8, State: "paused", CreatedAt: now}
	repo := &fakeTranslationJobManagementLifecycleRepository{
		jobByID: map[int64]repository.TranslationJob{41: job},
		phaseRunsByJobID: map[int64][]repository.JobPhaseRun{
			41: {{ID: 410, TranslationJobID: 41, PhaseType: "body_translation", State: "running", ProgressPercent: 52}},
		},
		snapshotsByJobID: map[int64][]repository.TranslationJobPhaseRuntimeSnapshot{41: {
			{TranslationJobID: 41, PhaseID: "body_translation", Provider: "openai", ModelName: "gpt-5", ExecutionMode: "batch", CredentialStatus: "configured"},
		}},
		deleteResultByJobID: map[int64]repository.TranslationJobDeleteResult{
			41: {Outcome: repository.TranslationJobDeleteOutcomeDeleted, Job: &job},
		},
	}
	sourceRepo := &fakeTranslationJobManagementSourceRepository{sourceByID: map[int64]repository.XEditExtractedData{8: {ID: 8, SourceFilePath: "/mods/inconsistent.json", TargetPluginName: "Skyrim.esm"}}}
	service := NewTranslationJobManagementService(repo, sourceRepo, fakeTranslationJobManagementTransactor{})
	service.cacheReader = fakeTranslationJobManagementCacheReader{available: true}

	result, err := service.DeleteJob(context.Background(), 41)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if result.ReasonCategory != translationJobManagementReasonStateProjectionInconsistent {
		t.Fatalf("expected state_projection_inconsistent, got %#v", result)
	}
	if result.Detail == nil || result.Detail.DeleteAvailability.Enabled {
		t.Fatalf("expected blocked delete detail, got %#v", result.Detail)
	}
	if len(repo.deleteCalls) != 0 {
		t.Fatalf("expected delete repository not to run, got %#v", repo.deleteCalls)
	}
}

func TestTranslationJobManagementServiceResumeJobReturnsCacheMissingReason(t *testing.T) {
	now := time.Date(2026, 5, 6, 10, 0, 0, 0, time.UTC)
	job := repository.TranslationJob{ID: 50, XEditExtractedDataID: 10, State: "paused", CreatedAt: now}
	repo := &fakeTranslationJobManagementLifecycleRepository{
		jobByID: map[int64]repository.TranslationJob{50: job},
		snapshotsByJobID: map[int64][]repository.TranslationJobPhaseRuntimeSnapshot{50: {
			{TranslationJobID: 50, PhaseID: "text_translation", Provider: "openai", ModelName: "gpt-5", ExecutionMode: "batch", CredentialStatus: "configured"},
		}},
	}
	sourceRepo := &fakeTranslationJobManagementSourceRepository{sourceByID: map[int64]repository.XEditExtractedData{10: {ID: 10, SourceFilePath: "/mods/long/long/path.json", TargetPluginName: "Skyrim.esm"}}}
	service := NewTranslationJobManagementService(repo, sourceRepo, fakeTranslationJobManagementTransactor{})
	service.cacheReader = fakeTranslationJobManagementCacheReader{available: false}

	result, err := service.ResumeJob(context.Background(), 50)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if result.ReasonCategory != "cache_missing" {
		t.Fatalf("expected cache_missing, got %#v", result)
	}
	if result.Detail == nil || result.Detail.RuntimeSummary.CredentialState == "" {
		t.Fatalf("expected redacted runtime summary, got %#v", result.Detail)
	}
}
