package service

import (
	"context"
	"errors"
	"fmt"
	"slices"
	"sort"
	"strings"

	"aitranslationenginejp/internal/repository"
)

const (
	translationOutputArtifactJobStateCompleted    = "completed"
	translationOutputArtifactJobStateCanceled     = "canceled"
	translationOutputArtifactBodyPhaseType        = "body_translation"
	translationOutputArtifactPhaseStateCompleted  = "completed"
	translationOutputArtifactPhaseStateCanceled   = "canceled"
	translationOutputArtifactStatusNotGenerated   = "not_generated"
	translationOutputArtifactScanFallbackMaxJobID = int64(2048)
)

var translationOutputArtifactReadyStatuses = []string{"ready", "translated", "cached"}

type translationOutputArtifactJobLifecycleRepository interface {
	GetTranslationJobByID(ctx context.Context, id int64) (repository.TranslationJob, error)
	ListJobPhaseRunsByJobID(ctx context.Context, jobID int64) ([]repository.JobPhaseRun, error)
}

type translationOutputArtifactJobOutputRepository interface {
	ListJobTranslationFieldsByJobID(ctx context.Context, jobID int64) ([]repository.JobTranslationField, error)
}

type translationOutputArtifactPersistenceRepository interface {
	GetTranslationArtifactByID(ctx context.Context, id int64) (repository.TranslationArtifact, error)
	GetTranslationArtifactByJobID(ctx context.Context, jobID int64) (repository.TranslationArtifact, error)
	UpsertTranslationArtifact(
		ctx context.Context,
		draft repository.TranslationArtifactDraft,
	) (repository.TranslationArtifact, error)
	ReplaceXTranslatorOutputRows(
		ctx context.Context,
		translationArtifactID int64,
		drafts []repository.XTranslatorOutputRowDraft,
	) ([]repository.XTranslatorOutputRow, error)
	CountXTranslatorOutputRowsByArtifactID(ctx context.Context, translationArtifactID int64) (int, error)
}

type translationOutputArtifactTranslationSourceRepository interface {
	GetXEditExtractedDataByID(ctx context.Context, id int64) (repository.XEditExtractedData, error)
	GetTranslationRecordByID(ctx context.Context, id int64) (repository.TranslationRecord, error)
	GetTranslationFieldByID(ctx context.Context, id int64) (repository.TranslationField, error)
}

type translationOutputArtifactJobLister interface {
	ListTranslationJobs(ctx context.Context) ([]repository.TranslationJob, error)
}

// TranslationOutputReviewReadModel stores the Output Review read model.
type TranslationOutputReviewReadModel struct {
	CompletedJobs    []TranslationOutputCompletedJobSummaryReadModel
	SelectedJob      TranslationOutputSelectedJobSummaryReadModel
	OutputReadiness  TranslationOutputReadinessSummaryReadModel
	ArtifactStatus   TranslationOutputArtifactStatusSummaryReadModel
	RejectionReasons []TranslationOutputArtifactErrorSummaryReadModel
}

// TranslationOutputCompletedJobSummaryReadModel stores one completed job row.
type TranslationOutputCompletedJobSummaryReadModel struct {
	JobID                    int64
	JobStatus                string
	ArtifactStatus           string
	OutputReady              bool
	TranslatedCount          int
	OutputStatusDistribution map[string]int
}

// TranslationOutputSelectedJobSummaryReadModel stores one selected job summary.
type TranslationOutputSelectedJobSummaryReadModel struct {
	JobID           int64
	JobStatus       string
	BodyPhaseStatus string
	OutputReady     bool
	ResultSummary   TranslationOutputResultSummaryReadModel
}

// TranslationOutputResultSummaryReadModel stores output counts and provenance.
type TranslationOutputResultSummaryReadModel struct {
	TranslatedCount int
	RowCount        int
	InputProvenance TranslationOutputInputProvenanceSummaryReadModel
}

// TranslationOutputInputProvenanceSummaryReadModel stores safe provenance values.
type TranslationOutputInputProvenanceSummaryReadModel struct {
	InputSnapshotDigest string
	SourceFileDigest    string
}

// TranslationOutputReadinessSummaryReadModel stores readiness state.
type TranslationOutputReadinessSummaryReadModel struct {
	Ready         bool
	Retryable     bool
	RejectionKind string
}

// TranslationOutputArtifactStatusSummaryReadModel stores current artifact state.
type TranslationOutputArtifactStatusSummaryReadModel struct {
	ArtifactID     int64
	Status         string
	RowCount       int
	CurrentVersion bool
}

// TranslationOutputArtifactErrorSummaryReadModel stores safe rejection details.
type TranslationOutputArtifactErrorSummaryReadModel struct {
	ErrorKind  string
	Reason     string
	Retryable  bool
	IsRedacted bool
}

// TranslationOutputDiffPreviewReadModel stores the diff preview read model.
type TranslationOutputDiffPreviewReadModel struct {
	JobID                int64
	ArtifactID           int64
	Rows                 []TranslationOutputDiffPreviewRowReadModel
	CompatibilitySummary TranslationOutputCompatibilitySummaryReadModel
}

// TranslationOutputDiffPreviewRowReadModel stores one xTranslator-compatible preview row.
type TranslationOutputDiffPreviewRowReadModel struct {
	FieldID              int64
	RowDigest            string
	EDID                 string
	REC                  string
	FIELD                string
	FORMID               string
	SourceExcerpt        string
	DestExcerpt          string
	XTranslatorStatus    int
	InternalOutputStatus string
	RowReflectionSummary string
	StaleReason          string
	CanRegenerate        bool
}

// TranslationOutputCompatibilitySummaryReadModel stores compatibility counts.
type TranslationOutputCompatibilitySummaryReadModel struct {
	Passed       bool
	WarningCount int
	RejectCount  int
}

type translationOutputArtifactLoadedJob struct {
	job         repository.TranslationJob
	bodyRun     *repository.JobPhaseRun
	outputs     []repository.JobTranslationField
	inputSource repository.XEditExtractedData
}

// TranslationOutputArtifactService builds Output Review backend read models.
type TranslationOutputArtifactService struct {
	jobLifecycleRepository  translationOutputArtifactJobLifecycleRepository
	jobOutputRepository     translationOutputArtifactJobOutputRepository
	translationSourceReader translationOutputArtifactTranslationSourceRepository
	persistenceRepository   translationOutputArtifactPersistenceRepository
	transactor              repository.Transactor
	xmlSerializer           translationOutputArtifactXMLSerializer
	fileWriter              translationOutputArtifactFileWriter
}

// NewTranslationOutputArtifactService creates the output artifact review service.
func NewTranslationOutputArtifactService(
	jobLifecycleRepository translationOutputArtifactJobLifecycleRepository,
	jobOutputRepository translationOutputArtifactJobOutputRepository,
	translationSourceReader translationOutputArtifactTranslationSourceRepository,
) *TranslationOutputArtifactService {
	service := &TranslationOutputArtifactService{
		jobLifecycleRepository:  jobLifecycleRepository,
		jobOutputRepository:     jobOutputRepository,
		translationSourceReader: translationSourceReader,
		xmlSerializer:           NewXTranslatorOutputArtifactXMLSerializer(),
		fileWriter:              NewLocalTranslationOutputArtifactFileWriter(),
	}
	if persistenceRepository, ok := any(jobLifecycleRepository).(translationOutputArtifactPersistenceRepository); ok {
		service.persistenceRepository = persistenceRepository
	}
	return service
}

// WithArtifactPersistence configures DB-backed artifact persistence for command execution.
func (service *TranslationOutputArtifactService) WithArtifactPersistence(
	persistenceRepository translationOutputArtifactPersistenceRepository,
	transactor repository.Transactor,
) *TranslationOutputArtifactService {
	service.persistenceRepository = persistenceRepository
	service.transactor = transactor
	return service
}

// WithFileWriter configures the filesystem adapter used by generate and regenerate.
func (service *TranslationOutputArtifactService) WithFileWriter(
	fileWriter translationOutputArtifactFileWriter,
) *TranslationOutputArtifactService {
	if fileWriter != nil {
		service.fileWriter = fileWriter
	}
	return service
}

// WithXMLSerializer configures the XML serializer used by generate and regenerate.
func (service *TranslationOutputArtifactService) WithXMLSerializer(
	xmlSerializer translationOutputArtifactXMLSerializer,
) *TranslationOutputArtifactService {
	if xmlSerializer != nil {
		service.xmlSerializer = xmlSerializer
	}
	return service
}

// ReadReview returns the Output Review read model for one selected job.
func (service *TranslationOutputArtifactService) ReadReview(
	ctx context.Context,
	selectedJobID *int64,
) (TranslationOutputReviewReadModel, error) {
	jobs, err := service.listTranslationJobs(ctx, selectedJobID)
	if err != nil {
		return TranslationOutputReviewReadModel{}, fmt.Errorf("list translation jobs: %w", err)
	}
	if len(jobs) == 0 {
		return TranslationOutputReviewReadModel{}, nil
	}

	sort.Slice(jobs, func(i, j int) bool {
		return jobs[i].ID < jobs[j].ID
	})

	selectedID := service.resolveSelectedJobID(jobs, selectedJobID)
	selectedLoaded, err := service.loadJob(ctx, selectedID)
	if err != nil {
		return TranslationOutputReviewReadModel{}, fmt.Errorf("load selected translation job %d: %w", selectedID, err)
	}

	completedJobs := make([]TranslationOutputCompletedJobSummaryReadModel, 0)
	for _, job := range jobs {
		if !strings.EqualFold(strings.TrimSpace(job.State), translationOutputArtifactJobStateCompleted) {
			continue
		}
		loaded, loadErr := service.loadJob(ctx, job.ID)
		if loadErr != nil {
			return TranslationOutputReviewReadModel{}, fmt.Errorf("load completed translation job %d: %w", job.ID, loadErr)
		}
		if loaded.bodyRun == nil {
			continue
		}
		readiness := service.buildReadiness(loaded)
		artifactStatus := service.lookupArtifactStatus(ctx, loaded.job.ID)
		completedJobs = append(completedJobs, TranslationOutputCompletedJobSummaryReadModel{
			JobID:                    loaded.job.ID,
			JobStatus:                loaded.job.State,
			ArtifactStatus:           artifactStatus.Status,
			OutputReady:              readiness.Ready,
			TranslatedCount:          countReadyOutputs(loaded.outputs),
			OutputStatusDistribution: buildOutputStatusDistribution(loaded.outputs),
		})
	}

	selectedReadiness := service.buildReadiness(selectedLoaded)
	selectedArtifactStatus := service.lookupArtifactStatus(ctx, selectedLoaded.job.ID)
	rejectionReasons := []TranslationOutputArtifactErrorSummaryReadModel(nil)
	if !selectedReadiness.Ready {
		rejectionReasons = append(rejectionReasons, buildReadinessRejection(selectedReadiness.RejectionKind))
	}

	return TranslationOutputReviewReadModel{
		CompletedJobs: completedJobs,
		SelectedJob: TranslationOutputSelectedJobSummaryReadModel{
			JobID:           selectedLoaded.job.ID,
			JobStatus:       selectedLoaded.job.State,
			BodyPhaseStatus: bodyPhaseState(selectedLoaded.bodyRun),
			OutputReady:     selectedReadiness.Ready,
			ResultSummary: TranslationOutputResultSummaryReadModel{
				TranslatedCount: countReadyOutputs(selectedLoaded.outputs),
				RowCount:        len(selectedLoaded.outputs),
				InputProvenance: TranslationOutputInputProvenanceSummaryReadModel{
					InputSnapshotDigest: bodyPhaseInputDigest(selectedLoaded.bodyRun),
					SourceFileDigest:    selectedLoaded.inputSource.SourceContentHash,
				},
			},
		},
		OutputReadiness:  selectedReadiness,
		ArtifactStatus:   TranslationOutputArtifactStatusSummaryReadModel(selectedArtifactStatus),
		RejectionReasons: rejectionReasons,
	}, nil
}

func (service *TranslationOutputArtifactService) loadJob(
	ctx context.Context,
	jobID int64,
) (translationOutputArtifactLoadedJob, error) {
	job, err := service.jobLifecycleRepository.GetTranslationJobByID(ctx, jobID)
	if err != nil {
		return translationOutputArtifactLoadedJob{}, fmt.Errorf("get translation job: %w", err)
	}
	phaseRuns, err := service.jobLifecycleRepository.ListJobPhaseRunsByJobID(ctx, job.ID)
	if err != nil {
		return translationOutputArtifactLoadedJob{}, fmt.Errorf("list job phase runs: %w", err)
	}
	outputs, err := service.jobOutputRepository.ListJobTranslationFieldsByJobID(ctx, job.ID)
	if err != nil {
		return translationOutputArtifactLoadedJob{}, fmt.Errorf("list job translation fields: %w", err)
	}
	inputSource, err := service.translationSourceReader.GetXEditExtractedDataByID(ctx, job.XEditExtractedDataID)
	if err != nil {
		return translationOutputArtifactLoadedJob{}, fmt.Errorf("get input provenance: %w", err)
	}
	return translationOutputArtifactLoadedJob{
		job:         job,
		bodyRun:     findBodyPhaseRun(phaseRuns),
		outputs:     outputs,
		inputSource: inputSource,
	}, nil
}

func (service *TranslationOutputArtifactService) buildReadiness(
	loaded translationOutputArtifactLoadedJob,
) TranslationOutputReadinessSummaryReadModel {
	jobState := strings.ToLower(strings.TrimSpace(loaded.job.State))
	if jobState == translationOutputArtifactJobStateCanceled {
		return buildReadinessSummary("canceled")
	}
	if jobState != translationOutputArtifactJobStateCompleted {
		return buildReadinessSummary("not_completed")
	}
	if loaded.bodyRun == nil {
		return buildReadinessSummary("not_completed")
	}

	bodyState := strings.ToLower(strings.TrimSpace(loaded.bodyRun.State))
	if bodyState == translationOutputArtifactPhaseStateCanceled {
		return buildReadinessSummary("canceled")
	}
	if bodyState != translationOutputArtifactPhaseStateCompleted {
		return buildReadinessSummary("not_completed")
	}

	targetCount := max(loaded.bodyRun.ProviderTargetCount, loaded.bodyRun.SnapshotFieldCount)
	if targetCount == 0 {
		return TranslationOutputReadinessSummaryReadModel{Ready: true}
	}
	if len(loaded.outputs) != targetCount {
		return buildReadinessSummary("status_mismatch")
	}
	for _, output := range loaded.outputs {
		status := strings.ToLower(strings.TrimSpace(output.OutputStatus))
		if !slices.Contains(translationOutputArtifactReadyStatuses, status) {
			return buildReadinessSummary("status_mismatch")
		}
	}
	return TranslationOutputReadinessSummaryReadModel{Ready: true}
}

func (service *TranslationOutputArtifactService) listTranslationJobs(
	ctx context.Context,
	selectedJobID *int64,
) ([]repository.TranslationJob, error) {
	if lister, ok := service.jobLifecycleRepository.(translationOutputArtifactJobLister); ok {
		jobs, err := lister.ListTranslationJobs(ctx)
		if err != nil {
			return nil, fmt.Errorf("list translation jobs from repository: %w", err)
		}
		return append([]repository.TranslationJob(nil), jobs...), nil
	}

	upperBound := translationOutputArtifactScanFallbackMaxJobID
	if selectedJobID != nil && *selectedJobID > upperBound {
		upperBound = *selectedJobID + 1024
	}
	jobs := make([]repository.TranslationJob, 0)
	for id := int64(1); id <= upperBound; id++ {
		job, err := service.jobLifecycleRepository.GetTranslationJobByID(ctx, id)
		if err != nil {
			if errors.Is(err, repository.ErrNotFound) {
				continue
			}
			return nil, fmt.Errorf("get translation job by fallback scan: %w", err)
		}
		jobs = append(jobs, job)
	}
	return jobs, nil
}

func (service *TranslationOutputArtifactService) resolveSelectedJobID(
	jobs []repository.TranslationJob,
	selectedJobID *int64,
) int64 {
	if selectedJobID != nil {
		return *selectedJobID
	}
	for _, job := range jobs {
		if strings.EqualFold(strings.TrimSpace(job.State), translationOutputArtifactJobStateCompleted) {
			return job.ID
		}
	}
	return jobs[0].ID
}

func findBodyPhaseRun(phaseRuns []repository.JobPhaseRun) *repository.JobPhaseRun {
	for _, phaseRun := range phaseRuns {
		if phaseRun.PhaseType == translationOutputArtifactBodyPhaseType {
			run := phaseRun
			return &run
		}
	}
	return nil
}

func countReadyOutputs(outputs []repository.JobTranslationField) int {
	count := 0
	for _, output := range outputs {
		status := strings.ToLower(strings.TrimSpace(output.OutputStatus))
		if slices.Contains(translationOutputArtifactReadyStatuses, status) {
			count++
		}
	}
	return count
}

func buildOutputStatusDistribution(outputs []repository.JobTranslationField) map[string]int {
	if len(outputs) == 0 {
		return nil
	}
	distribution := make(map[string]int, len(outputs))
	for _, output := range outputs {
		status := strings.ToLower(strings.TrimSpace(output.OutputStatus))
		if status == "" {
			status = "unknown"
		}
		distribution[status]++
	}
	return distribution
}

func bodyPhaseState(bodyRun *repository.JobPhaseRun) string {
	if bodyRun == nil {
		return ""
	}
	return bodyRun.State
}

func bodyPhaseInputDigest(bodyRun *repository.JobPhaseRun) string {
	if bodyRun == nil {
		return ""
	}
	return bodyRun.InputSnapshotDigest
}

func buildReadinessSummary(kind string) TranslationOutputReadinessSummaryReadModel {
	return TranslationOutputReadinessSummaryReadModel{
		Ready:         false,
		Retryable:     false,
		RejectionKind: kind,
	}
}

func buildReadinessRejection(kind string) TranslationOutputArtifactErrorSummaryReadModel {
	return TranslationOutputArtifactErrorSummaryReadModel{
		ErrorKind:  kind,
		Reason:     buildReadinessReason(kind),
		Retryable:  false,
		IsRedacted: true,
	}
}

func buildReadinessReason(kind string) string {
	switch kind {
	case "canceled":
		return "job was canceled before output generation"
	case "status_mismatch":
		return "job output results are not aligned with completed body phase status"
	default:
		return "job is not completed for output generation"
	}
}
