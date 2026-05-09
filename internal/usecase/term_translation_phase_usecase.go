package usecase

import (
	"context"
	"fmt"
	"time"

	"aitranslationenginejp/internal/service"
)

type termTranslationPhaseServicePort interface {
	ReadSummary(ctx context.Context, jobID int64) (service.TermTranslationPhaseSummaryReadModel, error)
	StartPhase(ctx context.Context, jobID int64) (service.TermTranslationPhaseCommandReadModel, error)
	PausePhase(ctx context.Context, jobID int64, phaseRunID int64) (service.TermTranslationPhaseCommandReadModel, error)
	ResumePhase(ctx context.Context, jobID int64, phaseRunID int64) (service.TermTranslationPhaseCommandReadModel, error)
	RetryPhase(ctx context.Context, jobID int64, phaseRunID int64) (service.TermTranslationPhaseCommandReadModel, error)
	ReadNextPhaseReadiness(ctx context.Context, jobID int64) (service.TermTranslationNextPhaseReadinessReadModel, error)
}

const termTranslationPhaseUsecaseLogWhere = "backend.usecase.term_translation_phase"

// TermTranslationPhaseUsecase bridges the frozen term phase contract to the service layer.
type TermTranslationPhaseUsecase struct {
	service termTranslationPhaseServicePort
}

// NewTermTranslationPhaseUsecase creates the frozen term phase usecase.
func NewTermTranslationPhaseUsecase(service termTranslationPhaseServicePort) *TermTranslationPhaseUsecase {
	return &TermTranslationPhaseUsecase{service: service}
}

// GetTermTranslationPhaseSummary returns the term phase summary contract.
func (usecase *TermTranslationPhaseUsecase) GetTermTranslationPhaseSummary(
	ctx context.Context,
	request GetTermTranslationPhaseSummaryRequest,
) (TermTranslationPhaseSummaryResult, error) {
	readModel, err := usecase.service.ReadSummary(ctx, request.JobID)
	if err != nil {
		return TermTranslationPhaseSummaryResult{}, fmt.Errorf("get term translation phase summary: %w", err)
	}
	return toTermTranslationPhaseSummaryResult(readModel), nil
}

// StartTermTranslationPhase starts one term phase run and returns the command contract.
func (usecase *TermTranslationPhaseUsecase) StartTermTranslationPhase(
	ctx context.Context,
	request StartTermTranslationPhaseRequest,
) (TermTranslationPhaseCommandResult, error) {
	readModel, err := usecase.service.StartPhase(ctx, request.JobID)
	logPhaseStateCommand(ctx, "term_translation_phase_start", termTranslationPhaseUsecaseLogWhere, request.JobID, readModel.PhaseRunID, readModel.BeforePhaseState, readModel.AfterPhaseState, termTranslationReadModelErrorReason(readModel.ErrorSummary), err)
	if err != nil {
		return TermTranslationPhaseCommandResult{}, fmt.Errorf("start term translation phase: %w", err)
	}
	return toTermTranslationPhaseCommandResult(readModel), nil
}

// PauseTermTranslationPhase pauses one active term phase run and returns the command contract.
func (usecase *TermTranslationPhaseUsecase) PauseTermTranslationPhase(
	ctx context.Context,
	request PauseTermTranslationPhaseRequest,
) (TermTranslationPhaseCommandResult, error) {
	readModel, err := usecase.service.PausePhase(ctx, request.JobID, request.PhaseRunID)
	phaseRunID := request.PhaseRunID
	logPhaseStateCommand(ctx, "term_translation_phase_pause", termTranslationPhaseUsecaseLogWhere, request.JobID, &phaseRunID, readModel.BeforePhaseState, readModel.AfterPhaseState, termTranslationReadModelErrorReason(readModel.ErrorSummary), err)
	if err != nil {
		return TermTranslationPhaseCommandResult{}, fmt.Errorf("pause term translation phase: %w", err)
	}
	return toTermTranslationPhaseCommandResult(readModel), nil
}

// ResumeTermTranslationPhase resumes one paused term phase run and returns the command contract.
func (usecase *TermTranslationPhaseUsecase) ResumeTermTranslationPhase(
	ctx context.Context,
	request ResumeTermTranslationPhaseRequest,
) (TermTranslationPhaseCommandResult, error) {
	readModel, err := usecase.service.ResumePhase(ctx, request.JobID, request.PhaseRunID)
	phaseRunID := request.PhaseRunID
	logPhaseStateCommand(ctx, "term_translation_phase_resume", termTranslationPhaseUsecaseLogWhere, request.JobID, &phaseRunID, readModel.BeforePhaseState, readModel.AfterPhaseState, termTranslationReadModelErrorReason(readModel.ErrorSummary), err)
	if err != nil {
		return TermTranslationPhaseCommandResult{}, fmt.Errorf("resume term translation phase: %w", err)
	}
	return toTermTranslationPhaseCommandResult(readModel), nil
}

// RetryTermTranslationPhase retries one recoverable-failed term phase run and returns the command contract.
func (usecase *TermTranslationPhaseUsecase) RetryTermTranslationPhase(
	ctx context.Context,
	request RetryTermTranslationPhaseRequest,
) (TermTranslationPhaseCommandResult, error) {
	readModel, err := usecase.service.RetryPhase(ctx, request.JobID, request.PhaseRunID)
	phaseRunID := request.PhaseRunID
	logPhaseStateCommand(ctx, "term_translation_phase_retry", termTranslationPhaseUsecaseLogWhere, request.JobID, &phaseRunID, readModel.BeforePhaseState, readModel.AfterPhaseState, termTranslationReadModelErrorReason(readModel.ErrorSummary), err)
	if err != nil {
		return TermTranslationPhaseCommandResult{}, fmt.Errorf("retry term translation phase: %w", err)
	}
	return toTermTranslationPhaseCommandResult(readModel), nil
}

// GetTermTranslationNextPhaseReadiness returns whether the next phase may start.
func (usecase *TermTranslationPhaseUsecase) GetTermTranslationNextPhaseReadiness(
	ctx context.Context,
	request GetTermTranslationNextPhaseReadinessRequest,
) (TermTranslationNextPhaseReadinessResult, error) {
	readModel, err := usecase.service.ReadNextPhaseReadiness(ctx, request.JobID)
	logPhaseReadiness(ctx, "term_translation_next_phase_readiness", termTranslationPhaseUsecaseLogWhere, request.JobID, readModel.PhaseState, readModel.CanStartNextPhase, readModel.BlockedReason, err)
	if err != nil {
		return TermTranslationNextPhaseReadinessResult{}, fmt.Errorf("get term translation next phase readiness: %w", err)
	}
	return TermTranslationNextPhaseReadinessResult{
		JobID:             readModel.JobID,
		CurrentPhase:      readModel.CurrentPhase,
		PhaseState:        readModel.PhaseState,
		CanStartNextPhase: readModel.CanStartNextPhase,
		BlockedReason:     cloneStringPointer(readModel.BlockedReason),
		ErrorKind:         NormalizeTermTranslationPhasePublicErrorKind(readModel.ErrorKind),
	}, nil
}

func termTranslationReadModelErrorReason(readModel *service.TermTranslationPhaseErrorReadModel) string {
	if readModel == nil {
		return ""
	}
	return readModel.ErrorKind
}

func toTermTranslationPhaseSummaryResult(
	readModel service.TermTranslationPhaseSummaryReadModel,
) TermTranslationPhaseSummaryResult {
	return TermTranslationPhaseSummaryResult{
		JobID:              readModel.JobID,
		CurrentPhase:       readModel.CurrentPhase,
		PhaseState:         readModel.PhaseState,
		PhaseRunID:         cloneInt64Pointer(readModel.PhaseRunID),
		StartedAt:          cloneTimePointer(readModel.StartedAt),
		FinishedAt:         cloneTimePointer(readModel.FinishedAt),
		Progress:           toTermTranslationPhaseProgressSummary(readModel.Progress),
		TotalTermCount:     readModel.TotalTermCount,
		DictionaryHitCount: readModel.DictionaryHitCount,
		AITargetCount:      readModel.AITargetCount,
		Execution: TermTranslationExecutionConfigSummary{
			CredentialRef:   readModel.Execution.CredentialRef,
			Provider:        readModel.Execution.Provider,
			Model:           readModel.Execution.Model,
			ExecutionMode:   readModel.Execution.ExecutionMode,
			SnapshotDigest:  cloneStringPointer(readModel.Execution.SnapshotDigest),
			SnapshotVersion: cloneStringPointer(readModel.Execution.SnapshotVersion),
		},
		ResultSummary:    toTermTranslationPhaseResultSummary(readModel.ResultSummary),
		ErrorSummary:     toTermTranslationPhaseErrorSummary(readModel.ErrorSummary),
		ActionEnablement: toTermTranslationPhaseActionEnablement(readModel.ActionEnablement),
	}
}

func toTermTranslationPhaseCommandResult(
	readModel service.TermTranslationPhaseCommandReadModel,
) TermTranslationPhaseCommandResult {
	return TermTranslationPhaseCommandResult{
		JobID:             readModel.JobID,
		CurrentPhase:      readModel.CurrentPhase,
		PhaseState:        readModel.PhaseState,
		PhaseRunID:        cloneInt64Pointer(readModel.PhaseRunID),
		StartedAt:         cloneTimePointer(readModel.StartedAt),
		FinishedAt:        cloneTimePointer(readModel.FinishedAt),
		Progress:          toTermTranslationPhaseProgressSummary(readModel.Progress),
		Retryable:         readModel.Retryable,
		CanStartNextPhase: readModel.CanStartNextPhase,
		ErrorSummary:      toTermTranslationPhaseErrorSummary(readModel.ErrorSummary),
	}
}

func toTermTranslationPhaseProgressSummary(
	readModel service.TermTranslationPhaseProgressReadModel,
) TermTranslationPhaseProgressSummary {
	return TermTranslationPhaseProgressSummary{
		Percent:        readModel.Percent,
		ProcessedCount: readModel.ProcessedCount,
		TotalCount:     readModel.TotalCount,
		AITargetCount:  readModel.AITargetCount,
		CurrentStep:    readModel.CurrentStep,
	}
}

func toTermTranslationPhaseResultSummary(
	readModel *service.TermTranslationPhaseResultReadModel,
) *TermTranslationPhaseResultSummary {
	if readModel == nil {
		return nil
	}
	return &TermTranslationPhaseResultSummary{
		ConfirmedCount:            readModel.ConfirmedCount,
		JobDictionaryAppliedCount: readModel.JobDictionaryAppliedCount,
		ReplacementTargetCount:    readModel.ReplacementTargetCount,
		UnmatchedCount:            readModel.UnmatchedCount,
		ProviderSkipped:           readModel.ProviderSkipped,
	}
}

func toTermTranslationPhaseErrorSummary(
	readModel *service.TermTranslationPhaseErrorReadModel,
) *TermTranslationPhaseErrorSummary {
	if readModel == nil {
		return nil
	}
	return &TermTranslationPhaseErrorSummary{
		ErrorKind:  NormalizeTermTranslationPhasePublicErrorKind(readModel.ErrorKind),
		Reason:     readModel.Reason,
		Retryable:  readModel.Retryable,
		IsRedacted: readModel.IsRedacted,
	}
}

func toTermTranslationPhaseActionEnablement(
	readModel service.TermTranslationPhaseActionEnablementReadModel,
) TermTranslationPhaseActionEnablement {
	return TermTranslationPhaseActionEnablement{
		CanStart:               readModel.CanStart,
		StartBlockedReason:     cloneStringPointer(readModel.StartBlockedReason),
		CanPause:               readModel.CanPause,
		PauseBlockedReason:     cloneStringPointer(readModel.PauseBlockedReason),
		CanResume:              readModel.CanResume,
		ResumeBlockedReason:    cloneStringPointer(readModel.ResumeBlockedReason),
		CanRetry:               readModel.CanRetry,
		RetryBlockedReason:     cloneStringPointer(readModel.RetryBlockedReason),
		CanRefresh:             readModel.CanRefresh,
		RefreshBlockedReason:   cloneStringPointer(readModel.RefreshBlockedReason),
		CanStartNextPhase:      readModel.CanStartNextPhase,
		NextPhaseBlockedReason: cloneStringPointer(readModel.NextPhaseBlockedReason),
	}
}

func cloneStringPointer(value *string) *string {
	if value == nil {
		return nil
	}
	cloned := *value
	return &cloned
}

func cloneInt64Pointer(value *int64) *int64 {
	if value == nil {
		return nil
	}
	cloned := *value
	return &cloned
}

func cloneTimePointer(value *time.Time) *time.Time {
	if value == nil {
		return nil
	}
	cloned := *value
	return &cloned
}
