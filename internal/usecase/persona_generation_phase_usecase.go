package usecase

import (
	"context"
	"fmt"
	"time"

	"aitranslationenginejp/internal/service"
)

type personaGenerationPhaseServicePort interface {
	ReadSummary(ctx context.Context, jobID int64) (service.PersonaGenerationPhaseSummaryReadModel, error)
	StartPhase(ctx context.Context, jobID int64) (service.PersonaGenerationPhaseCommandReadModel, error)
	PausePhase(ctx context.Context, jobID int64, phaseRunID int64) (service.PersonaGenerationPhaseCommandReadModel, error)
	ResumePhase(ctx context.Context, jobID int64, phaseRunID int64) (service.PersonaGenerationPhaseCommandReadModel, error)
	RetryPhase(ctx context.Context, jobID int64, phaseRunID int64) (service.PersonaGenerationPhaseCommandReadModel, error)
	CancelPhase(ctx context.Context, jobID int64, phaseRunID int64) (service.PersonaGenerationPhaseCommandReadModel, error)
	ReadBodyReadiness(ctx context.Context, jobID int64) (service.PersonaGenerationBodyReadinessReadModel, error)
}

const personaGenerationPhaseUsecaseLogWhere = "backend.usecase.persona_generation_phase"

// PersonaGenerationPhaseUsecase bridges the persona phase contract to the service layer.
type PersonaGenerationPhaseUsecase struct {
	service personaGenerationPhaseServicePort
}

// NewPersonaGenerationPhaseUsecase creates the persona phase usecase.
func NewPersonaGenerationPhaseUsecase(service personaGenerationPhaseServicePort) *PersonaGenerationPhaseUsecase {
	return &PersonaGenerationPhaseUsecase{service: service}
}

// GetPersonaGenerationPhaseSummary returns the persona phase summary contract.
func (usecase *PersonaGenerationPhaseUsecase) GetPersonaGenerationPhaseSummary(
	ctx context.Context,
	request GetPersonaGenerationPhaseSummaryRequest,
) (PersonaGenerationPhaseSummaryResult, error) {
	readModel, err := usecase.service.ReadSummary(ctx, request.JobID)
	if err != nil {
		return PersonaGenerationPhaseSummaryResult{}, fmt.Errorf("get persona generation phase summary: %w", err)
	}
	return toPersonaGenerationPhaseSummaryResult(readModel), nil
}

// StartPersonaGenerationPhase starts one persona phase run and returns the command contract.
func (usecase *PersonaGenerationPhaseUsecase) StartPersonaGenerationPhase(
	ctx context.Context,
	request StartPersonaGenerationPhaseRequest,
) (PersonaGenerationPhaseCommandResult, error) {
	readModel, err := usecase.service.StartPhase(ctx, request.JobID)
	logPhaseStateCommand(ctx, "persona_generation_phase_start", personaGenerationPhaseUsecaseLogWhere, request.JobID, readModel.PhaseRunID, readModel.BeforePhaseState, readModel.AfterPhaseState, personaGenerationReadModelErrorReason(readModel.ErrorSummary), err)
	result := toPersonaGenerationPhaseCommandResult(readModel)
	if err != nil {
		return result, fmt.Errorf("start persona generation phase: %w", err)
	}
	return result, nil
}

// PausePersonaGenerationPhase pauses one active persona phase run.
func (usecase *PersonaGenerationPhaseUsecase) PausePersonaGenerationPhase(
	ctx context.Context,
	request PausePersonaGenerationPhaseRequest,
) (PersonaGenerationPhaseCommandResult, error) {
	readModel, err := usecase.service.PausePhase(ctx, request.JobID, request.PhaseRunID)
	phaseRunID := request.PhaseRunID
	logPhaseStateCommand(ctx, "persona_generation_phase_pause", personaGenerationPhaseUsecaseLogWhere, request.JobID, &phaseRunID, readModel.BeforePhaseState, readModel.AfterPhaseState, personaGenerationReadModelErrorReason(readModel.ErrorSummary), err)
	if err != nil {
		return PersonaGenerationPhaseCommandResult{}, fmt.Errorf("pause persona generation phase: %w", err)
	}
	return toPersonaGenerationPhaseCommandResult(readModel), nil
}

// ResumePersonaGenerationPhase resumes one paused or recoverable failed persona phase run.
func (usecase *PersonaGenerationPhaseUsecase) ResumePersonaGenerationPhase(
	ctx context.Context,
	request ResumePersonaGenerationPhaseRequest,
) (PersonaGenerationPhaseCommandResult, error) {
	readModel, err := usecase.service.ResumePhase(ctx, request.JobID, request.PhaseRunID)
	phaseRunID := request.PhaseRunID
	logPhaseStateCommand(ctx, "persona_generation_phase_resume", personaGenerationPhaseUsecaseLogWhere, request.JobID, &phaseRunID, readModel.BeforePhaseState, readModel.AfterPhaseState, personaGenerationReadModelErrorReason(readModel.ErrorSummary), err)
	result := toPersonaGenerationPhaseCommandResult(readModel)
	if err != nil {
		return result, fmt.Errorf("resume persona generation phase: %w", err)
	}
	return result, nil
}

// RetryPersonaGenerationPhase retries one persona phase run.
func (usecase *PersonaGenerationPhaseUsecase) RetryPersonaGenerationPhase(
	ctx context.Context,
	request RetryPersonaGenerationPhaseRequest,
) (PersonaGenerationPhaseCommandResult, error) {
	readModel, err := usecase.service.RetryPhase(ctx, request.JobID, request.PhaseRunID)
	phaseRunID := request.PhaseRunID
	logPhaseStateCommand(ctx, "persona_generation_phase_retry", personaGenerationPhaseUsecaseLogWhere, request.JobID, &phaseRunID, readModel.BeforePhaseState, readModel.AfterPhaseState, personaGenerationReadModelErrorReason(readModel.ErrorSummary), err)
	result := toPersonaGenerationPhaseCommandResult(readModel)
	if err != nil {
		return result, fmt.Errorf("retry persona generation phase: %w", err)
	}
	return result, nil
}

// CancelPersonaGenerationPhase cancels one persona phase run.
func (usecase *PersonaGenerationPhaseUsecase) CancelPersonaGenerationPhase(
	ctx context.Context,
	request CancelPersonaGenerationPhaseRequest,
) (PersonaGenerationPhaseCommandResult, error) {
	readModel, err := usecase.service.CancelPhase(ctx, request.JobID, request.PhaseRunID)
	phaseRunID := request.PhaseRunID
	logPhaseStateCommand(ctx, "persona_generation_phase_cancel", personaGenerationPhaseUsecaseLogWhere, request.JobID, &phaseRunID, readModel.BeforePhaseState, readModel.AfterPhaseState, personaGenerationReadModelErrorReason(readModel.ErrorSummary), err)
	if err != nil {
		return PersonaGenerationPhaseCommandResult{}, fmt.Errorf("cancel persona generation phase: %w", err)
	}
	return toPersonaGenerationPhaseCommandResult(readModel), nil
}

// GetPersonaGenerationBodyReadiness returns whether the body phase may start.
func (usecase *PersonaGenerationPhaseUsecase) GetPersonaGenerationBodyReadiness(
	ctx context.Context,
	request GetPersonaGenerationBodyReadinessRequest,
) (PersonaGenerationBodyReadinessResult, error) {
	readModel, err := usecase.service.ReadBodyReadiness(ctx, request.JobID)
	logPhaseReadiness(ctx, "persona_generation_body_readiness", personaGenerationPhaseUsecaseLogWhere, request.JobID, readModel.PhaseState, readModel.Ready, readModel.BlockedReason, err)
	result := PersonaGenerationBodyReadinessResult{
		JobID:         readModel.JobID,
		CurrentPhase:  readModel.CurrentPhase,
		PhaseState:    readModel.PhaseState,
		Ready:         readModel.Ready,
		BlockedReason: clonePersonaGenerationStringPointer(readModel.BlockedReason),
		ErrorKind:     NormalizePersonaGenerationPhasePublicErrorKind(readModel.ErrorKind),
		InputSummary: PersonaGenerationBodyReadinessInputSummary{
			PersonaCount:   readModel.InputSummary.PersonaCount,
			MissingCount:   readModel.InputSummary.MissingCount,
			SnapshotID:     readModel.InputSummary.SnapshotID,
			SnapshotDigest: readModel.InputSummary.SnapshotDigest,
			EvidenceRefs:   append([]string(nil), readModel.InputSummary.EvidenceRefs...),
		},
	}
	if err != nil {
		return result, fmt.Errorf("get persona generation body readiness: %w", err)
	}
	return result, nil
}

func personaGenerationReadModelErrorReason(readModel *service.PersonaGenerationPhaseErrorSummaryReadModel) string {
	if readModel == nil {
		return ""
	}
	return readModel.ErrorKind
}

func toPersonaGenerationPhaseSummaryResult(
	readModel service.PersonaGenerationPhaseSummaryReadModel,
) PersonaGenerationPhaseSummaryResult {
	return PersonaGenerationPhaseSummaryResult{
		JobID:        readModel.JobID,
		CurrentPhase: readModel.CurrentPhase,
		PhaseState:   readModel.PhaseState,
		PhaseRunID:   clonePersonaGenerationInt64Pointer(readModel.PhaseRunID),
		StartedAt:    clonePersonaGenerationTimePointer(readModel.StartedAt),
		FinishedAt:   clonePersonaGenerationTimePointer(readModel.FinishedAt),
		Progress: PersonaGenerationPhaseProgressSummary{
			Percent:        readModel.Progress.Percent,
			ProcessedCount: readModel.Progress.ProcessedCount,
			TotalCount:     readModel.Progress.TotalCount,
			TargetCount:    readModel.Progress.TargetCount,
			CurrentStep:    readModel.Progress.CurrentStep,
		},
		TargetSummary: PersonaGenerationTargetSummary{
			TargetCount:            readModel.TargetSummary.TargetCount,
			CommonPersonaHitCount:  readModel.TargetSummary.CommonPersonaHitCount,
			CommonPersonaMissCount: readModel.TargetSummary.CommonPersonaMissCount,
			SkippedCount:           readModel.TargetSummary.SkippedCount,
			SkippedReasons:         append([]string(nil), readModel.TargetSummary.SkippedReasons...),
			TargetSnapshotID:       clonePersonaGenerationStringPointer(readModel.TargetSummary.TargetSnapshotID),
			TargetSnapshotDigest:   readModel.TargetSummary.TargetSnapshotDigest,
		},
		Execution: PersonaGenerationExecutionSummary{
			CredentialRef: readModel.Execution.CredentialRef,
			Provider:      readModel.Execution.Provider,
			Model:         readModel.Execution.Model,
			ExecutionMode: readModel.Execution.ExecutionMode,
			PromptDigest:  readModel.Execution.PromptDigest,
			InputCount:    readModel.Execution.InputCount,
			OutputCount:   readModel.Execution.OutputCount,
			EvidenceRefs:  append([]string(nil), readModel.Execution.EvidenceRefs...),
		},
		ResultSummary:    toPersonaGenerationPhaseResultSummary(readModel.ResultSummary),
		ErrorSummary:     toPersonaGenerationPhaseErrorSummary(readModel.ErrorSummary),
		ActionEnablement: toPersonaGenerationPhaseActionEnablement(readModel.ActionEnablement),
	}
}

func toPersonaGenerationPhaseCommandResult(
	readModel service.PersonaGenerationPhaseCommandReadModel,
) PersonaGenerationPhaseCommandResult {
	return PersonaGenerationPhaseCommandResult{
		JobID:        readModel.JobID,
		CurrentPhase: readModel.CurrentPhase,
		PhaseState:   readModel.PhaseState,
		PhaseRunID:   clonePersonaGenerationInt64Pointer(readModel.PhaseRunID),
		StartedAt:    clonePersonaGenerationTimePointer(readModel.StartedAt),
		FinishedAt:   clonePersonaGenerationTimePointer(readModel.FinishedAt),
		Progress: PersonaGenerationPhaseProgressSummary{
			Percent:        readModel.Progress.Percent,
			ProcessedCount: readModel.Progress.ProcessedCount,
			TotalCount:     readModel.Progress.TotalCount,
			TargetCount:    readModel.Progress.TargetCount,
			CurrentStep:    readModel.Progress.CurrentStep,
		},
		TargetSummary: PersonaGenerationTargetSummary{
			TargetCount:            readModel.TargetSummary.TargetCount,
			CommonPersonaHitCount:  readModel.TargetSummary.CommonPersonaHitCount,
			CommonPersonaMissCount: readModel.TargetSummary.CommonPersonaMissCount,
			SkippedCount:           readModel.TargetSummary.SkippedCount,
			SkippedReasons:         append([]string(nil), readModel.TargetSummary.SkippedReasons...),
			TargetSnapshotID:       clonePersonaGenerationStringPointer(readModel.TargetSummary.TargetSnapshotID),
			TargetSnapshotDigest:   readModel.TargetSummary.TargetSnapshotDigest,
		},
		Execution: PersonaGenerationExecutionSummary{
			CredentialRef: readModel.Execution.CredentialRef,
			Provider:      readModel.Execution.Provider,
			Model:         readModel.Execution.Model,
			ExecutionMode: readModel.Execution.ExecutionMode,
			PromptDigest:  readModel.Execution.PromptDigest,
			InputCount:    readModel.Execution.InputCount,
			OutputCount:   readModel.Execution.OutputCount,
			EvidenceRefs:  append([]string(nil), readModel.Execution.EvidenceRefs...),
		},
		ResultSummary:     toPersonaGenerationPhaseResultSummary(readModel.ResultSummary),
		Retryable:         readModel.Retryable,
		CanStartBodyPhase: readModel.CanStartBodyPhase,
		ErrorSummary:      toPersonaGenerationPhaseErrorSummary(readModel.ErrorSummary),
	}
}

func toPersonaGenerationPhaseResultSummary(
	readModel *service.PersonaGenerationPhaseResultSummaryReadModel,
) *PersonaGenerationPhaseResultSummary {
	if readModel == nil {
		return nil
	}
	return &PersonaGenerationPhaseResultSummary{
		GeneratedCount:          readModel.GeneratedCount,
		FailedCount:             readModel.FailedCount,
		PersonaCount:            readModel.PersonaCount,
		MissingCount:            readModel.MissingCount,
		SnapshotID:              readModel.SnapshotID,
		SnapshotDigest:          readModel.SnapshotDigest,
		SnapshotReferenceStatus: readModel.SnapshotReferenceStatus,
		BodyReadiness:           readModel.BodyReadiness,
	}
}

func toPersonaGenerationPhaseErrorSummary(
	readModel *service.PersonaGenerationPhaseErrorSummaryReadModel,
) *PersonaGenerationPhaseErrorSummary {
	if readModel == nil {
		return nil
	}
	return &PersonaGenerationPhaseErrorSummary{
		ErrorKind:  NormalizePersonaGenerationPhasePublicErrorKind(readModel.ErrorKind),
		Reason:     readModel.Reason,
		Retryable:  readModel.Retryable,
		IsRedacted: readModel.IsRedacted,
	}
}

func toPersonaGenerationPhaseActionEnablement(
	readModel service.PersonaGenerationPhaseActionEnablementReadModel,
) PersonaGenerationPhaseActionEnablement {
	return PersonaGenerationPhaseActionEnablement{
		CanStart:               readModel.CanStart,
		StartBlockedReason:     clonePersonaGenerationStringPointer(readModel.StartBlockedReason),
		CanPause:               readModel.CanPause,
		PauseBlockedReason:     clonePersonaGenerationStringPointer(readModel.PauseBlockedReason),
		CanResume:              readModel.CanResume,
		ResumeBlockedReason:    clonePersonaGenerationStringPointer(readModel.ResumeBlockedReason),
		CanRetry:               readModel.CanRetry,
		RetryBlockedReason:     clonePersonaGenerationStringPointer(readModel.RetryBlockedReason),
		CanCancel:              readModel.CanCancel,
		CancelBlockedReason:    clonePersonaGenerationStringPointer(readModel.CancelBlockedReason),
		CanStartBodyPhase:      readModel.CanStartBodyPhase,
		BodyPhaseBlockedReason: clonePersonaGenerationStringPointer(readModel.BodyPhaseBlockedReason),
	}
}

func clonePersonaGenerationStringPointer(value *string) *string {
	if value == nil {
		return nil
	}
	cloned := *value
	return &cloned
}

func clonePersonaGenerationInt64Pointer(value *int64) *int64 {
	if value == nil {
		return nil
	}
	cloned := *value
	return &cloned
}

func clonePersonaGenerationTimePointer(value *time.Time) *time.Time {
	if value == nil {
		return nil
	}
	cloned := *value
	return &cloned
}
