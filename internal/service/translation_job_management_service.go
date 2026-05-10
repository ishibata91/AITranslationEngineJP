package service

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"path/filepath"
	"strings"
	"time"

	"aitranslationenginejp/internal/repository"
)

const (
	translationJobManagementReasonCacheMissing                   = "cache_missing"
	translationJobManagementReasonTerminalState                  = "terminal_state"
	translationJobManagementReasonStateProjectionInconsistent    = "state_projection_inconsistent"
	translationJobManagementReasonPhaseProgressAggregationFailed = "phase_progress_aggregation_failed"
	translationJobManagementReasonStaleSelection                 = "stale_selection"
	translationJobManagementReasonRunningDeleteBlocked           = "running_delete_blocked"
	translationJobManagementReasonStopRequested                  = "stop_requested"
	translationJobManagementReasonStopFailed                     = "stop_failed"
	translationJobManagementReasonDeleteFailed                   = "delete_failed"
	translationJobManagementReasonResumeFailed                   = "resume_failed"
)

const (
	translationJobManagementMessageStaleSelection     = "選択中の job は一覧に残っていません。再読込してください。"
	translationJobManagementJobStateReady             = "ready"
	translationJobManagementJobStateRunning           = "running"
	translationJobManagementJobStatePaused            = "paused"
	translationJobManagementJobStateRecoverableFailed = "recoverable_failed"
	translationJobManagementJobStateCompleted         = "completed"
	translationJobManagementJobStateFailed            = "failed"
	translationJobManagementJobStateCanceled          = "canceled"
)

// TranslationJobManagementListReadModel returns the incomplete-job list.
type TranslationJobManagementListReadModel struct {
	Jobs []TranslationJobManagementJobSummaryReadModel
}

// TranslationJobManagementInputSourceSummaryReadModel summarizes one input source.
type TranslationJobManagementInputSourceSummaryReadModel struct {
	InputSourceID        int64
	InputSourceLabel     string
	InputSourceKindLabel string
	SourcePath           string
	PluginName           string
	ExtractedJSONLabel   string
}

// TranslationJobManagementProgressSummaryReadModel summarizes one job progress state.
type TranslationJobManagementProgressSummaryReadModel struct {
	CurrentPhase      string
	CurrentPhaseLabel string
	Percent           int
	ProgressLabel     string
	LastUpdatedLabel  string
}

// TranslationJobManagementProtectedSettingSummaryReadModel exposes redacted runtime settings.
type TranslationJobManagementProtectedSettingSummaryReadModel struct {
	ProviderLabel        string
	ModelLabel           string
	ExecutionModeLabel   string
	CredentialState      string
	CredentialStateLabel string
}

// TranslationJobManagementOperationAvailabilityReadModel describes one operation guard.
type TranslationJobManagementOperationAvailabilityReadModel struct {
	Kind           string
	Enabled        bool
	Label          string
	HelperText     string
	ReasonCategory string
	ReasonText     string
}

// TranslationJobManagementBlockedReasonReadModel explains one blocked or warning state.
type TranslationJobManagementBlockedReasonReadModel struct {
	Category string
	Title    string
	Detail   string
}

// TranslationJobManagementJobSummaryReadModel is one incomplete-job list item.
type TranslationJobManagementJobSummaryReadModel struct {
	JobID              int64
	JobState           string
	JobStateLabel      string
	StateTone          string
	CanOpenPhase       bool
	OpenBlockedReason  *TranslationJobManagementBlockedReasonReadModel
	InputSource        TranslationJobManagementInputSourceSummaryReadModel
	Progress           TranslationJobManagementProgressSummaryReadModel
	StopAvailability   TranslationJobManagementOperationAvailabilityReadModel
	ResumeAvailability TranslationJobManagementOperationAvailabilityReadModel
	DeleteAvailability TranslationJobManagementOperationAvailabilityReadModel
}

// TranslationJobManagementJobDetailReadModel is one selected job detail.
type TranslationJobManagementJobDetailReadModel struct {
	TranslationJobManagementJobSummaryReadModel
	CacheState           string
	CacheStateLabel      string
	RuntimeSummary       TranslationJobManagementProtectedSettingSummaryReadModel
	ResumeBlockedReasons []TranslationJobManagementBlockedReasonReadModel
	Warnings             []TranslationJobManagementBlockedReasonReadModel
	DeleteImpactLines    []string
}

// TranslationJobManagementActionReadModel returns one public action outcome.
type TranslationJobManagementActionReadModel struct {
	Message        string
	Tone           string
	Detail         *TranslationJobManagementJobDetailReadModel
	DeletedJobID   *int64
	ReasonCategory string
}

type translationJobManagementLifecycleReadRepository interface {
	ListIncompleteTranslationJobs(ctx context.Context) ([]repository.TranslationJob, error)
	ListTranslationJobPhaseRuntimeSnapshots(ctx context.Context, translationJobID int64) ([]repository.TranslationJobPhaseRuntimeSnapshot, error)
	DeleteNonRunningTranslationJob(ctx context.Context, jobID int64) (repository.TranslationJobDeleteResult, error)
}

type translationJobManagementCacheReader interface {
	HasTranslationCacheByXEditID(ctx context.Context, xEditID int64) (bool, error)
}

// TranslationJobManagementService projects incomplete jobs for management views.
type TranslationJobManagementService struct {
	jobLifecycleRepository repository.JobLifecycleRepository
	managementRepository   translationJobManagementLifecycleReadRepository
	translationSourceRepo  repository.TranslationSourceRepository
	cacheReader            translationJobManagementCacheReader
	transactor             repository.Transactor
}

// NewTranslationJobManagementService builds the Job Management service.
func NewTranslationJobManagementService(
	jobLifecycleRepository repository.JobLifecycleRepository,
	translationSourceRepository repository.TranslationSourceRepository,
	transactor repository.Transactor,
) *TranslationJobManagementService {
	service := &TranslationJobManagementService{
		jobLifecycleRepository: jobLifecycleRepository,
		translationSourceRepo:  translationSourceRepository,
		transactor:             transactor,
	}
	if managementRepository, ok := jobLifecycleRepository.(translationJobManagementLifecycleReadRepository); ok {
		service.managementRepository = managementRepository
	}
	if cacheReader, ok := translationSourceRepository.(translationJobManagementCacheReader); ok {
		service.cacheReader = cacheReader
	}
	return service
}

// ListIncompleteJobs loads every non-completed job summary.
func (service *TranslationJobManagementService) ListIncompleteJobs(
	ctx context.Context,
) (TranslationJobManagementListReadModel, error) {
	if service.managementRepository == nil {
		return TranslationJobManagementListReadModel{}, fmt.Errorf("translation job management repository is not configured")
	}
	jobs, err := service.managementRepository.ListIncompleteTranslationJobs(ctx)
	if err != nil {
		return TranslationJobManagementListReadModel{}, fmt.Errorf("list incomplete translation jobs: %w", err)
	}
	result := make([]TranslationJobManagementJobSummaryReadModel, 0, len(jobs))
	for _, job := range jobs {
		detail, detailErr := service.buildJobDetail(ctx, job)
		if detailErr != nil {
			return TranslationJobManagementListReadModel{}, detailErr
		}
		result = append(result, detail.TranslationJobManagementJobSummaryReadModel)
	}
	return TranslationJobManagementListReadModel{Jobs: result}, nil
}

// GetJobDetail loads one selected incomplete job detail.
func (service *TranslationJobManagementService) GetJobDetail(
	ctx context.Context,
	jobID int64,
) (TranslationJobManagementJobDetailReadModel, error) {
	job, err := service.jobLifecycleRepository.GetTranslationJobByID(ctx, jobID)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return TranslationJobManagementJobDetailReadModel{}, fmt.Errorf("%s: %w", translationJobManagementReasonStaleSelection, err)
		}
		return TranslationJobManagementJobDetailReadModel{}, fmt.Errorf("get translation job by id: %w", err)
	}
	if normalizeTranslationJobManagementValue(job.State) == translationJobManagementJobStateCompleted {
		return TranslationJobManagementJobDetailReadModel{}, fmt.Errorf("%s: completed jobs are outside management scope", translationJobManagementReasonStaleSelection)
	}
	return service.buildJobDetail(ctx, job)
}

// DeleteJob deletes one non-running job and its job-owned rows.
func (service *TranslationJobManagementService) DeleteJob(
	ctx context.Context,
	jobID int64,
) (TranslationJobManagementActionReadModel, error) {
	if service.managementRepository == nil || service.transactor == nil {
		return TranslationJobManagementActionReadModel{}, fmt.Errorf("translation job management delete dependencies are not configured")
	}

	var deletionResult repository.TranslationJobDeleteResult
	var blockedAction *TranslationJobManagementActionReadModel
	err := service.transactor.WithTransaction(ctx, func(txCtx context.Context) error {
		result, action, txErr := service.executeDeleteJobTransaction(txCtx, jobID)
		deletionResult = result
		blockedAction = action
		return txErr
	})
	if err != nil {
		return TranslationJobManagementActionReadModel{}, fmt.Errorf("delete translation job transaction: %w", err)
	}
	if blockedAction != nil {
		return *blockedAction, nil
	}

	return service.buildDeleteJobAction(ctx, deletionResult)
}

func (service *TranslationJobManagementService) executeDeleteJobTransaction(
	ctx context.Context,
	jobID int64,
) (repository.TranslationJobDeleteResult, *TranslationJobManagementActionReadModel, error) {
	job, blockedAction, err := service.getDeleteCandidate(ctx, jobID)
	if err != nil || blockedAction != nil {
		if blockedAction != nil {
			logTranslationJobDeleteRejected(ctx, jobID, "", blockedAction.ReasonCategory)
		}
		return repository.TranslationJobDeleteResult{}, blockedAction, err
	}
	blockedAction, err = service.buildDeleteBlockedAction(ctx, job)
	if err != nil || blockedAction != nil {
		if blockedAction != nil {
			logTranslationJobDeleteRejected(ctx, jobID, job.State, blockedAction.ReasonCategory)
		}
		return repository.TranslationJobDeleteResult{}, blockedAction, err
	}
	result, err := service.managementRepository.DeleteNonRunningTranslationJob(ctx, jobID)
	if err != nil {
		return repository.TranslationJobDeleteResult{}, nil, fmt.Errorf("delete non-running translation job: %w", err)
	}
	logTranslationJobDeleteResult(ctx, jobID, job.State, result.Outcome)
	return result, nil, nil
}

func (service *TranslationJobManagementService) getDeleteCandidate(
	ctx context.Context,
	jobID int64,
) (repository.TranslationJob, *TranslationJobManagementActionReadModel, error) {
	job, err := service.jobLifecycleRepository.GetTranslationJobByID(ctx, jobID)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			action := buildStaleSelectionAction()
			return repository.TranslationJob{}, &action, nil
		}
		return repository.TranslationJob{}, nil, fmt.Errorf("get translation job by id: %w", err)
	}
	if normalizeTranslationJobManagementValue(job.State) == translationJobManagementJobStateCompleted {
		action := buildStaleSelectionAction()
		return repository.TranslationJob{}, &action, nil
	}
	return job, nil, nil
}

func (service *TranslationJobManagementService) buildDeleteBlockedAction(
	ctx context.Context,
	job repository.TranslationJob,
) (*TranslationJobManagementActionReadModel, error) {
	detail, err := service.buildJobDetail(ctx, job)
	if err != nil {
		return nil, err
	}
	if detail.DeleteAvailability.Enabled {
		return nil, nil
	}
	action := buildBlockedDeleteAction(detail)
	return &action, nil
}

func (service *TranslationJobManagementService) buildDeleteJobAction(
	ctx context.Context,
	deletionResult repository.TranslationJobDeleteResult,
) (TranslationJobManagementActionReadModel, error) {
	switch deletionResult.Outcome {
	case repository.TranslationJobDeleteOutcomeNotFound:
		return buildStaleSelectionAction(), nil
	case repository.TranslationJobDeleteOutcomeBlockedRunning:
		detail, detailErr := service.buildDetailFromDeleteResult(ctx, deletionResult)
		if detailErr != nil {
			return TranslationJobManagementActionReadModel{}, detailErr
		}
		return TranslationJobManagementActionReadModel{
			Message:        "実行中 job は削除できません。停止または中断後に再判定してください。",
			Tone:           "warning",
			Detail:         detail,
			ReasonCategory: translationJobManagementReasonRunningDeleteBlocked,
		}, nil
	case repository.TranslationJobDeleteOutcomeBlockedUnsafe:
		detail, detailErr := service.buildDetailFromDeleteResult(ctx, deletionResult)
		if detailErr != nil {
			return TranslationJobManagementActionReadModel{}, detailErr
		}
		action := buildBlockedDeleteActionForDetail(detail)
		return action, nil
	case repository.TranslationJobDeleteOutcomeDeleted:
		if deletionResult.Job == nil {
			return TranslationJobManagementActionReadModel{
				Message:        "job 情報だけを削除しました。",
				Tone:           "success",
				ReasonCategory: "",
			}, nil
		}
		deletedJobID := deletionResult.Job.ID
		return TranslationJobManagementActionReadModel{
			Message:      "job 本体と job 配下の DB 情報を削除しました。入力データと入力ファイルは残しています。",
			Tone:         "success",
			DeletedJobID: &deletedJobID,
		}, nil
	default:
		return TranslationJobManagementActionReadModel{
			Message:        "削除結果を判定できませんでした。",
			Tone:           "danger",
			ReasonCategory: translationJobManagementReasonDeleteFailed,
		}, nil
	}
}

func buildStaleSelectionAction() TranslationJobManagementActionReadModel {
	return TranslationJobManagementActionReadModel{
		Message:        translationJobManagementMessageStaleSelection,
		Tone:           "warning",
		ReasonCategory: translationJobManagementReasonStaleSelection,
	}
}

func buildBlockedDeleteAction(
	detail TranslationJobManagementJobDetailReadModel,
) TranslationJobManagementActionReadModel {
	return buildBlockedDeleteActionForDetail(&detail)
}

func buildBlockedDeleteActionForDetail(
	detail *TranslationJobManagementJobDetailReadModel,
) TranslationJobManagementActionReadModel {
	category := translationJobManagementReasonDeleteFailed
	if detail != nil && strings.TrimSpace(detail.DeleteAvailability.ReasonCategory) != "" {
		category = detail.DeleteAvailability.ReasonCategory
	}
	message := "現在状態では削除できません。表示中の削除不可理由を確認してください。"
	switch category {
	case translationJobManagementReasonRunningDeleteBlocked:
		message = "実行中 job は削除できません。停止または中断後に再判定してください。"
	case translationJobManagementReasonStopRequested:
		message = "停止要求中は削除できません。状態更新後に再判定してください。"
	case translationJobManagementReasonStateProjectionInconsistent:
		message = "状態不整合があるため削除できません。job と phase 状態を保持したまま再読込してください。"
	}
	return TranslationJobManagementActionReadModel{
		Message:        message,
		Tone:           "warning",
		Detail:         detail,
		ReasonCategory: category,
	}
}

func logTranslationJobDeleteResult(
	ctx context.Context,
	jobID int64,
	beforeState string,
	outcome repository.TranslationJobDeleteOutcome,
) {
	switch outcome {
	case repository.TranslationJobDeleteOutcomeDeleted:
		slog.InfoContext(ctx, "translation job delete state changed",
			slog.String("event", "translation_job_delete"),
			slog.String("where", "backend.service.translation_job_management.delete"),
			slog.String("result", "allowed"),
			slog.String("id", fmt.Sprintf("job:%d", jobID)),
			slog.String("before_state", normalizeTranslationJobManagementValue(beforeState)),
			slog.String("after_state", "deleted"),
		)
	case repository.TranslationJobDeleteOutcomeNotFound:
		logTranslationJobDeleteRejected(ctx, jobID, beforeState, translationJobManagementReasonStaleSelection)
	case repository.TranslationJobDeleteOutcomeBlockedRunning:
		logTranslationJobDeleteRejected(ctx, jobID, beforeState, translationJobManagementReasonRunningDeleteBlocked)
	case repository.TranslationJobDeleteOutcomeBlockedUnsafe:
		logTranslationJobDeleteRejected(ctx, jobID, beforeState, translationJobManagementReasonDeleteFailed)
	default:
		logTranslationJobDeleteRejected(ctx, jobID, beforeState, translationJobManagementReasonDeleteFailed)
	}
}

func logTranslationJobDeleteRejected(
	ctx context.Context,
	jobID int64,
	state string,
	reason string,
) {
	if strings.TrimSpace(reason) == "" {
		reason = translationJobManagementReasonDeleteFailed
	}
	normalizedState := normalizeTranslationJobManagementValue(state)
	if normalizedState == "" {
		normalizedState = "unknown"
	}
	slog.WarnContext(ctx, "translation job delete rejected",
		slog.String("event", "translation_job_delete"),
		slog.String("where", "backend.service.translation_job_management.delete"),
		slog.String("result", "rejected"),
		slog.String("id", fmt.Sprintf("job:%d", jobID)),
		slog.String("before_state", normalizedState),
		slog.String("after_state", normalizedState),
		slog.String("reason", reason),
	)
}

// RequestStop returns the stop-facing projection without executing stop control.
func (service *TranslationJobManagementService) RequestStop(
	ctx context.Context,
	jobID int64,
) (TranslationJobManagementActionReadModel, error) {
	detail, err := service.GetJobDetail(ctx, jobID)
	if err != nil {
		return TranslationJobManagementActionReadModel{}, err
	}
	if detail.StopAvailability.ReasonCategory == translationJobManagementReasonStopRequested {
		return TranslationJobManagementActionReadModel{
			Message:        "停止要求中です。状態更新後に削除可否を再判定してください。",
			Tone:           "info",
			Detail:         &detail,
			ReasonCategory: translationJobManagementReasonStopRequested,
		}, nil
	}
	return TranslationJobManagementActionReadModel{
		Message:        "停止要求の表示境界だけを提供しています。停止制御本体は後続 task の対象です。",
		Tone:           "warning",
		Detail:         &detail,
		ReasonCategory: translationJobManagementReasonStopFailed,
	}, nil
}

// ResumeJob returns the resume-facing projection without executing resume control.
func (service *TranslationJobManagementService) ResumeJob(
	ctx context.Context,
	jobID int64,
) (TranslationJobManagementActionReadModel, error) {
	detail, err := service.GetJobDetail(ctx, jobID)
	if err != nil {
		return TranslationJobManagementActionReadModel{}, err
	}
	if !detail.ResumeAvailability.Enabled {
		category := detail.ResumeAvailability.ReasonCategory
		if category == "" {
			category = translationJobManagementReasonResumeFailed
		}
		return TranslationJobManagementActionReadModel{
			Message:        "現在状態では再開を開始できません。表示中の再開不可理由を確認してください。",
			Tone:           "warning",
			Detail:         &detail,
			ReasonCategory: category,
		}, nil
	}
	return TranslationJobManagementActionReadModel{
		Message:        "再開入口の要約だけを提供しています。再開実行本体は後続 task の対象です。",
		Tone:           "warning",
		Detail:         &detail,
		ReasonCategory: translationJobManagementReasonResumeFailed,
	}, nil
}

func (service *TranslationJobManagementService) buildDetailFromDeleteResult(
	ctx context.Context,
	result repository.TranslationJobDeleteResult,
) (*TranslationJobManagementJobDetailReadModel, error) {
	if result.Job == nil {
		return nil, nil
	}
	detail, err := service.buildJobDetail(ctx, *result.Job)
	if err != nil {
		return nil, err
	}
	return &detail, nil
}

func (service *TranslationJobManagementService) buildJobDetail(
	ctx context.Context,
	job repository.TranslationJob,
) (TranslationJobManagementJobDetailReadModel, error) {
	inputSource, warnings, err := service.buildInputSourceSummary(ctx, job.XEditExtractedDataID)
	if err != nil {
		return TranslationJobManagementJobDetailReadModel{}, err
	}
	cacheState, cacheStateLabel, cacheAvailable, cacheWarnings, err := service.buildCacheState(ctx, job.XEditExtractedDataID)
	if err != nil {
		return TranslationJobManagementJobDetailReadModel{}, err
	}
	warnings = append(warnings, cacheWarnings...)

	phaseRuns, err := service.jobLifecycleRepository.ListJobPhaseRunsByJobID(ctx, job.ID)
	if err != nil {
		return TranslationJobManagementJobDetailReadModel{}, fmt.Errorf("list job phase runs by job id: %w", err)
	}
	snapshots, err := service.listSnapshots(ctx, job.ID)
	if err != nil {
		return TranslationJobManagementJobDetailReadModel{}, fmt.Errorf("list translation job phase runtime snapshots: %w", err)
	}

	progressSummary, progressWarnings := buildTranslationJobManagementProgress(job, phaseRuns)
	warnings = append(warnings, progressWarnings...)

	runtimeSummary, runtimeWarnings := buildTranslationJobManagementRuntimeSummary(snapshots)
	warnings = append(warnings, runtimeWarnings...)

	resumeBlockedReasons := buildTranslationJobManagementResumeBlockedReasons(job, cacheAvailable, warnings)
	stopAvailability := buildTranslationJobManagementStopAvailability(job, phaseRuns)
	resumeAvailability := buildTranslationJobManagementResumeAvailability(job, resumeBlockedReasons)
	deleteAvailability := buildTranslationJobManagementDeleteAvailability(job, warnings, phaseRuns)
	canOpenPhase, openBlockedReason := buildTranslationJobManagementPhaseNavigationAvailability(job, progressSummary, warnings)

	detail := TranslationJobManagementJobDetailReadModel{
		TranslationJobManagementJobSummaryReadModel: TranslationJobManagementJobSummaryReadModel{
			JobID:              job.ID,
			JobState:           publicTranslationJobManagementState(job.State),
			JobStateLabel:      publicTranslationJobManagementStateLabel(job.State),
			StateTone:          publicTranslationJobManagementStateTone(job.State),
			CanOpenPhase:       canOpenPhase,
			OpenBlockedReason:  openBlockedReason,
			InputSource:        inputSource,
			Progress:           progressSummary,
			StopAvailability:   stopAvailability,
			ResumeAvailability: resumeAvailability,
			DeleteAvailability: deleteAvailability,
		},
		CacheState:           cacheState,
		CacheStateLabel:      cacheStateLabel,
		RuntimeSummary:       runtimeSummary,
		ResumeBlockedReasons: resumeBlockedReasons,
		Warnings:             warnings,
		DeleteImpactLines: []string{
			"削除対象は TRANSLATION_JOB と job 配下の DB 情報です。",
			"X_EDIT_EXTRACTED_DATA、抽出 JSON、入力ファイルは削除しません。",
			"Running job は削除できません。Paused 収束後に再判定してください。",
		},
	}
	return detail, nil
}

func (service *TranslationJobManagementService) buildInputSourceSummary(
	ctx context.Context,
	xEditID int64,
) (TranslationJobManagementInputSourceSummaryReadModel, []TranslationJobManagementBlockedReasonReadModel, error) {
	source, err := service.translationSourceRepo.GetXEditExtractedDataByID(ctx, xEditID)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return TranslationJobManagementInputSourceSummaryReadModel{
					InputSourceID:        xEditID,
					InputSourceLabel:     fmt.Sprintf("Input #%d", xEditID),
					InputSourceKindLabel: "xEdit 抽出データ",
					SourcePath:           "",
					PluginName:           "",
					ExtractedJSONLabel:   fmt.Sprintf("抽出データ #%d", xEditID),
				},
				[]TranslationJobManagementBlockedReasonReadModel{
					{
						Category: translationJobManagementReasonStateProjectionInconsistent,
						Title:    "入力参照が不整合です",
						Detail:   "job が参照する入力データを読み出せませんでした。削除可否と再開可否を安全側に倒します。",
					},
				},
				nil
		}
		return TranslationJobManagementInputSourceSummaryReadModel{}, nil, fmt.Errorf("get x_edit_extracted_data by id: %w", err)
	}
	sourcePath := strings.TrimSpace(source.SourceFilePath)
	inputSourceLabel := filepath.Base(sourcePath)
	if inputSourceLabel == "." || inputSourceLabel == "" {
		inputSourceLabel = fmt.Sprintf("Input #%d", xEditID)
	}
	pluginName := strings.TrimSpace(source.TargetPluginName)
	if pluginName == "" {
		pluginName = "未設定"
	}
	return TranslationJobManagementInputSourceSummaryReadModel{
		InputSourceID:        source.ID,
		InputSourceLabel:     inputSourceLabel,
		InputSourceKindLabel: "xEdit 抽出データ",
		SourcePath:           sourcePath,
		PluginName:           pluginName,
		ExtractedJSONLabel:   fmt.Sprintf("抽出データ #%d", source.ID),
	}, nil, nil
}

func (service *TranslationJobManagementService) buildCacheState(
	ctx context.Context,
	xEditID int64,
) (string, string, bool, []TranslationJobManagementBlockedReasonReadModel, error) {
	if service.cacheReader == nil {
		return "missing", "確認不可", false, []TranslationJobManagementBlockedReasonReadModel{
			{
				Category: translationJobManagementReasonStateProjectionInconsistent,
				Title:    "入力キャッシュ状態を確認できません",
				Detail:   "cache reader が未設定のため、安全側で再開不可として扱います。",
			},
		}, nil
	}
	available, err := service.cacheReader.HasTranslationCacheByXEditID(ctx, xEditID)
	if err != nil {
		return "", "", false, nil, fmt.Errorf("check translation cache by x_edit_id: %w", err)
	}
	if available {
		return "available", "利用可能", true, nil, nil
	}
	return "missing", "欠落", false, []TranslationJobManagementBlockedReasonReadModel{
		{
			Category: translationJobManagementReasonCacheMissing,
			Title:    "入力キャッシュが欠落しています",
			Detail:   "再開前に入力キャッシュを再構築する必要があります。",
		},
	}, nil
}

func (service *TranslationJobManagementService) listSnapshots(
	ctx context.Context,
	jobID int64,
) ([]repository.TranslationJobPhaseRuntimeSnapshot, error) {
	if service.managementRepository == nil {
		return nil, fmt.Errorf("translation job management repository is not configured")
	}
	snapshots, err := service.managementRepository.ListTranslationJobPhaseRuntimeSnapshots(ctx, jobID)
	if err != nil {
		return nil, fmt.Errorf("list translation job phase runtime snapshots: %w", err)
	}
	return snapshots, nil
}

func buildTranslationJobManagementProgress(
	job repository.TranslationJob,
	phaseRuns []repository.JobPhaseRun,
) (TranslationJobManagementProgressSummaryReadModel, []TranslationJobManagementBlockedReasonReadModel) {
	currentPhase := "term_translation"
	currentPhaseLabel := "開始待ち"
	percent := clampTranslationJobManagementPercent(job.ProgressPercent)
	lastUpdated := translationJobManagementTimestamp(job.CreatedAt)
	warnings := make([]TranslationJobManagementBlockedReasonReadModel, 0)

	currentRun, hasRun := findTranslationJobManagementCurrentRun(phaseRuns)
	if hasRun {
		currentPhase = publicTranslationJobManagementPhase(currentRun.PhaseType)
		currentPhaseLabel = publicTranslationJobManagementPhaseLabel(currentRun.PhaseType)
		percent = clampTranslationJobManagementPercent(currentRun.ProgressPercent)
		lastUpdated = translationJobManagementMostRecentTimestamp(job, currentRun)
		if percent == 0 && normalizeTranslationJobManagementValue(job.State) != translationJobManagementJobStateReady {
			warnings = append(warnings, TranslationJobManagementBlockedReasonReadModel{
				Category: translationJobManagementReasonPhaseProgressAggregationFailed,
				Title:    "進捗を集約できません",
				Detail:   "phase progress が 0% のままで、進捗の補足情報が不足しています。",
			})
		}
	} else if normalizeTranslationJobManagementValue(job.State) != translationJobManagementJobStateReady {
		warnings = append(warnings, TranslationJobManagementBlockedReasonReadModel{
			Category: translationJobManagementReasonStateProjectionInconsistent,
			Title:    "phase 状態が見つかりません",
			Detail:   "job は未完了ですが phase run が存在しません。状態表示を安全側に固定します。",
		})
	}

	return TranslationJobManagementProgressSummaryReadModel{
		CurrentPhase:      currentPhase,
		CurrentPhaseLabel: currentPhaseLabel,
		Percent:           percent,
		ProgressLabel:     fmt.Sprintf("%d%%", percent),
		LastUpdatedLabel:  lastUpdated,
	}, warnings
}

func buildTranslationJobManagementPhaseNavigationAvailability(
	job repository.TranslationJob,
	progress TranslationJobManagementProgressSummaryReadModel,
	warnings []TranslationJobManagementBlockedReasonReadModel,
) (bool, *TranslationJobManagementBlockedReasonReadModel) {
	state := normalizeTranslationJobManagementValue(job.State)
	if state == translationJobManagementJobStateCompleted {
		reason := TranslationJobManagementBlockedReasonReadModel{
			Category: translationJobManagementReasonStaleSelection,
			Title:    "完了済み job です",
			Detail:   "完了済み job は未完了一覧ではなく出力管理で扱います。",
		}
		return false, &reason
	}
	for _, warning := range warnings {
		if warning.Category == translationJobManagementReasonStateProjectionInconsistent ||
			warning.Category == translationJobManagementReasonPhaseProgressAggregationFailed {
			reason := warning
			return false, &reason
		}
	}
	if strings.TrimSpace(progress.CurrentPhase) == "" {
		reason := TranslationJobManagementBlockedReasonReadModel{
			Category: translationJobManagementReasonStateProjectionInconsistent,
			Title:    "現在の翻訳段階を判定できません",
			Detail:   "job 状態を変えずに未完了一覧で理由だけを表示します。",
		}
		return false, &reason
	}
	return true, nil
}

func buildTranslationJobManagementRuntimeSummary(
	snapshots []repository.TranslationJobPhaseRuntimeSnapshot,
) (TranslationJobManagementProtectedSettingSummaryReadModel, []TranslationJobManagementBlockedReasonReadModel) {
	if len(snapshots) == 0 {
		return TranslationJobManagementProtectedSettingSummaryReadModel{
				ProviderLabel:        "未設定",
				ModelLabel:           "未設定",
				ExecutionModeLabel:   "未設定",
				CredentialState:      "missing",
				CredentialStateLabel: "未設定",
			},
			[]TranslationJobManagementBlockedReasonReadModel{
				{
					Category: translationJobManagementReasonStateProjectionInconsistent,
					Title:    "保存済み AI 設定要約が不足しています",
					Detail:   "runtime snapshot が存在しないため、再開可否を安全側で評価します。",
				},
			}
	}
	snapshot := chooseTranslationJobManagementSnapshot(snapshots)
	credentialState, credentialStateLabel := publicTranslationJobManagementCredentialState(snapshot.CredentialStatus)
	return TranslationJobManagementProtectedSettingSummaryReadModel{
		ProviderLabel:        strings.TrimSpace(snapshot.Provider),
		ModelLabel:           strings.TrimSpace(snapshot.ModelName),
		ExecutionModeLabel:   strings.TrimSpace(snapshot.ExecutionMode),
		CredentialState:      credentialState,
		CredentialStateLabel: credentialStateLabel,
	}, nil
}

func buildTranslationJobManagementResumeBlockedReasons(
	job repository.TranslationJob,
	cacheAvailable bool,
	warnings []TranslationJobManagementBlockedReasonReadModel,
) []TranslationJobManagementBlockedReasonReadModel {
	state := normalizeTranslationJobManagementValue(job.State)
	result := make([]TranslationJobManagementBlockedReasonReadModel, 0)
	if state == translationJobManagementJobStatePaused || state == translationJobManagementJobStateRecoverableFailed {
		if !cacheAvailable {
			result = append(result, TranslationJobManagementBlockedReasonReadModel{
				Category: translationJobManagementReasonCacheMissing,
				Title:    "入力キャッシュが欠落しています",
				Detail:   "再開前に入力キャッシュを用意する必要があります。",
			})
		}
	}
	if state == translationJobManagementJobStateFailed || state == translationJobManagementJobStateCanceled || state == translationJobManagementJobStateCompleted {
		result = append(result, TranslationJobManagementBlockedReasonReadModel{
			Category: translationJobManagementReasonTerminalState,
			Title:    "terminal state です",
			Detail:   "完了済みまたは回復不能な状態のため、再開できません。",
		})
	}
	for _, warning := range warnings {
		if warning.Category == translationJobManagementReasonStateProjectionInconsistent || warning.Category == translationJobManagementReasonPhaseProgressAggregationFailed {
			result = append(result, TranslationJobManagementBlockedReasonReadModel{
				Category: translationJobManagementReasonStateProjectionInconsistent,
				Title:    warning.Title,
				Detail:   warning.Detail,
			})
		}
	}
	return result
}

func buildTranslationJobManagementStopAvailability(
	job repository.TranslationJob,
	phaseRuns []repository.JobPhaseRun,
) TranslationJobManagementOperationAvailabilityReadModel {
	state := normalizeTranslationJobManagementValue(job.State)
	result := TranslationJobManagementOperationAvailabilityReadModel{
		Kind:       "stop",
		Label:      "停止",
		HelperText: "実行中の job だけが停止対象です。",
	}
	switch {
	case hasTranslationJobManagementLatestError(phaseRuns, translationJobManagementReasonStopRequested):
		result.HelperText = "停止要求を受け付けています。状態更新後に再読込してください。"
		result.ReasonCategory = translationJobManagementReasonStopRequested
		result.ReasonText = "停止要求中"
	case state == translationJobManagementJobStateRunning:
		result.Enabled = true
		result.HelperText = "停止要求の入口を表示します。"
		if hasTranslationJobManagementLatestError(phaseRuns, translationJobManagementReasonStopFailed) {
			result.ReasonCategory = translationJobManagementReasonStopFailed
			result.ReasonText = "前回の停止要求は失敗しました。"
		}
	default:
		result.ReasonText = "Running job だけが停止対象です。"
	}
	return result
}

func buildTranslationJobManagementResumeAvailability(
	job repository.TranslationJob,
	blocked []TranslationJobManagementBlockedReasonReadModel,
) TranslationJobManagementOperationAvailabilityReadModel {
	result := TranslationJobManagementOperationAvailabilityReadModel{
		Kind:       "resume",
		Label:      "再開",
		HelperText: "Paused または RecoverableFailed の job で再開可否を確認できます。",
	}
	state := normalizeTranslationJobManagementValue(job.State)
	if state == translationJobManagementJobStatePaused || state == translationJobManagementJobStateRecoverableFailed {
		if len(blocked) == 0 {
			result.Enabled = true
			result.HelperText = "保存済み設定要約と入力キャッシュ状態は再開可能です。"
			return result
		}
		result.ReasonCategory = blocked[0].Category
		result.ReasonText = blocked[0].Detail
		return result
	}
	if len(blocked) > 0 {
		result.ReasonCategory = blocked[0].Category
		result.ReasonText = blocked[0].Detail
	} else {
		result.ReasonText = "現在状態では再開できません。"
	}
	return result
}

func buildTranslationJobManagementDeleteAvailability(
	job repository.TranslationJob,
	warnings []TranslationJobManagementBlockedReasonReadModel,
	phaseRuns []repository.JobPhaseRun,
) TranslationJobManagementOperationAvailabilityReadModel {
	result := TranslationJobManagementOperationAvailabilityReadModel{
		Kind:       "delete",
		Label:      "削除",
		HelperText: "job 配下の DB 情報だけを削除します。",
	}
	state := normalizeTranslationJobManagementValue(job.State)
	switch {
	case hasTranslationJobManagementLatestError(phaseRuns, translationJobManagementReasonStopRequested):
		result.ReasonCategory = translationJobManagementReasonStopRequested
		result.ReasonText = "停止要求中は削除可否を再判定できません。"
		return result
	case state == translationJobManagementJobStateRunning:
		result.ReasonCategory = translationJobManagementReasonRunningDeleteBlocked
		result.ReasonText = "Running job は削除できません。"
		return result
	case hasTranslationJobManagementRunningEquivalentPhaseRun(phaseRuns):
		result.ReasonCategory = translationJobManagementReasonStateProjectionInconsistent
		result.ReasonText = "phase 実行状態と job 状態が不整合なため、削除できません。"
		return result
	}
	for _, warning := range warnings {
		if warning.Category == translationJobManagementReasonStateProjectionInconsistent {
			result.ReasonCategory = translationJobManagementReasonStateProjectionInconsistent
			result.ReasonText = "状態不整合があるため、一覧を再読込してから削除してください。"
			return result
		}
	}
	result.Enabled = true
	return result
}

func findTranslationJobManagementCurrentRun(
	phaseRuns []repository.JobPhaseRun,
) (repository.JobPhaseRun, bool) {
	if len(phaseRuns) == 0 {
		return repository.JobPhaseRun{}, false
	}
	activeIndex := -1
	for index, run := range phaseRuns {
		switch normalizeTranslationJobManagementValue(run.State) {
		case "running", "paused", "recoverable_failed":
			activeIndex = index
		}
	}
	if activeIndex >= 0 {
		return phaseRuns[activeIndex], true
	}
	for _, phaseType := range []string{"term_translation", "translation", "persona_generation", "body_translation"} {
		for _, run := range phaseRuns {
			if normalizeTranslationJobManagementValue(run.PhaseType) != phaseType {
				continue
			}
			if !translationJobManagementPhaseRunCompleted(run) {
				return run, true
			}
		}
	}
	return phaseRuns[len(phaseRuns)-1], true
}

func translationJobManagementPhaseRunCompleted(run repository.JobPhaseRun) bool {
	switch normalizeTranslationJobManagementValue(run.State) {
	case "completed", "empty_completed":
		return true
	default:
		return false
	}
}

func chooseTranslationJobManagementSnapshot(
	snapshots []repository.TranslationJobPhaseRuntimeSnapshot,
) repository.TranslationJobPhaseRuntimeSnapshot {
	priority := []string{"text_translation", "npc_persona_generation", "word_translation"}
	for _, phaseID := range priority {
		for _, snapshot := range snapshots {
			if normalizeTranslationJobManagementValue(snapshot.PhaseID) == phaseID {
				return snapshot
			}
		}
	}
	return snapshots[0]
}

func hasTranslationJobManagementLatestError(
	phaseRuns []repository.JobPhaseRun,
	want string,
) bool {
	for _, run := range phaseRuns {
		if normalizeTranslationJobManagementValue(run.LatestError) == want {
			return true
		}
	}
	return false
}

func hasTranslationJobManagementRunningEquivalentPhaseRun(
	phaseRuns []repository.JobPhaseRun,
) bool {
	for _, run := range phaseRuns {
		state := normalizeTranslationJobManagementValue(run.State)
		if state == "running" {
			return true
		}
	}
	return false
}

func publicTranslationJobManagementState(value string) string {
	switch normalizeTranslationJobManagementValue(value) {
	case translationJobManagementJobStateReady:
		return "Ready"
	case translationJobManagementJobStateRunning:
		return "Running"
	case translationJobManagementJobStatePaused:
		return "Paused"
	case translationJobManagementJobStateRecoverableFailed:
		return "RecoverableFailed"
	case translationJobManagementJobStateCompleted:
		return "Completed"
	case translationJobManagementJobStateFailed:
		return "Failed"
	case translationJobManagementJobStateCanceled:
		return "Canceled"
	default:
		return strings.TrimSpace(value)
	}
}

func publicTranslationJobManagementStateLabel(value string) string {
	switch normalizeTranslationJobManagementValue(value) {
	case translationJobManagementJobStateReady:
		return "実行前"
	case translationJobManagementJobStateRunning:
		return "実行中"
	case translationJobManagementJobStatePaused:
		return "中断中"
	case translationJobManagementJobStateRecoverableFailed:
		return "再開可能な失敗"
	case translationJobManagementJobStateCompleted:
		return "完了"
	case translationJobManagementJobStateFailed:
		return "回復不能な失敗"
	case translationJobManagementJobStateCanceled:
		return "キャンセル済み"
	default:
		return "不明"
	}
}

func publicTranslationJobManagementStateTone(value string) string {
	switch normalizeTranslationJobManagementValue(value) {
	case translationJobManagementJobStateRunning:
		return "info"
	case translationJobManagementJobStatePaused, translationJobManagementJobStateRecoverableFailed:
		return "warning"
	case translationJobManagementJobStateFailed, translationJobManagementJobStateCanceled:
		return "danger"
	default:
		return "neutral"
	}
}

func publicTranslationJobManagementPhaseLabel(value string) string {
	switch normalizeTranslationJobManagementValue(value) {
	case "translation", "term_translation":
		return "用語翻訳"
	case "persona_generation":
		return "ペルソナ生成"
	case "body_translation":
		return "本文翻訳"
	default:
		return "開始待ち"
	}
}

func publicTranslationJobManagementPhase(value string) string {
	switch normalizeTranslationJobManagementValue(value) {
	case "translation", "term_translation":
		return "term_translation"
	case "persona_generation":
		return "persona_generation"
	case "body_translation":
		return "body_translation"
	default:
		return "term_translation"
	}
}

func publicTranslationJobManagementCredentialState(value string) (string, string) {
	switch normalizeTranslationJobManagementValue(value) {
	case "configured":
		return "configured", "設定済み"
	case "missing":
		return "missing", "未設定"
	case "inaccessible":
		return "inaccessible", "参照不可"
	case "not_required":
		return "configured", "不要"
	default:
		return "missing", "未設定"
	}
}

func clampTranslationJobManagementPercent(value int) int {
	if value < 0 {
		return 0
	}
	if value > 100 {
		return 100
	}
	return value
}

func normalizeTranslationJobManagementValue(value string) string {
	return strings.ToLower(strings.TrimSpace(value))
}

func translationJobManagementTimestamp(value time.Time) string {
	if value.IsZero() {
		return "未記録"
	}
	return value.UTC().Format("2006-01-02 15:04 UTC")
}

func translationJobManagementMostRecentTimestamp(
	job repository.TranslationJob,
	run repository.JobPhaseRun,
) string {
	candidates := []time.Time{job.CreatedAt}
	if job.StartedAt != nil {
		candidates = append(candidates, *job.StartedAt)
	}
	if job.FinishedAt != nil {
		candidates = append(candidates, *job.FinishedAt)
	}
	if run.StartedAt != nil {
		candidates = append(candidates, *run.StartedAt)
	}
	if run.FinishedAt != nil {
		candidates = append(candidates, *run.FinishedAt)
	}
	latest := candidates[0]
	for _, candidate := range candidates[1:] {
		if candidate.After(latest) {
			latest = candidate
		}
	}
	return translationJobManagementTimestamp(latest)
}
