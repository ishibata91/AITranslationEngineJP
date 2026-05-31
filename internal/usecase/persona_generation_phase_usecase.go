package usecase

import (
	"context"
	"fmt"
	"time"

	"aitranslationenginejp/internal/service"
	"aitranslationenginejp/internal/usecase/translationjobpolicy"
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

type personaGenerationPhaseAISettingsSaver interface {
	SaveAISettings(ctx context.Context, selection service.PhaseAISettingsSelection) (service.PhaseAISettingsReadModel, error)
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
	if result, allowed := usecase.evaluatePersonaPolicy(ctx, request.JobID, 0, translationjobpolicy.OperationStart); !allowed {
		return result, nil
	}
	readModel, err := usecase.service.StartPhase(ctx, request.JobID)
	logPhaseStateCommand(ctx, "persona_generation_phase_start", personaGenerationPhaseUsecaseLogWhere, request.JobID, readModel.PhaseRunID, readModel.BeforePhaseState, readModel.AfterPhaseState, personaGenerationReadModelErrorReason(readModel.ErrorSummary), err)
	result := toPersonaGenerationPhaseCommandResult(readModel)
	if err != nil {
		return result, fmt.Errorf("start persona generation phase: %w", err)
	}
	return result, nil
}

// SavePersonaGenerationPhaseAISettings saves public AI settings for the persona phase.
func (usecase *PersonaGenerationPhaseUsecase) SavePersonaGenerationPhaseAISettings(
	ctx context.Context,
	request SavePersonaGenerationPhaseAISettingsRequest,
) (PersonaGenerationPhaseAISettingsResult, error) {
	saver, ok := usecase.service.(personaGenerationPhaseAISettingsSaver)
	if !ok {
		return PersonaGenerationPhaseAISettingsResult{}, fmt.Errorf("save persona generation phase ai settings: service is not configured")
	}
	readModel, err := saver.SaveAISettings(ctx, service.PhaseAISettingsSelection{
		JobID:         request.JobID,
		Provider:      request.Provider,
		Model:         request.Model,
		ExecutionMode: request.ExecutionMode,
		BatchMode:     request.BatchMode,
	})
	if err != nil {
		return PersonaGenerationPhaseAISettingsResult{}, fmt.Errorf("save persona generation phase ai settings: %w", err)
	}
	return PersonaGenerationPhaseAISettingsResult(readModel), nil
}

// PausePersonaGenerationPhase pauses one active persona phase run.
func (usecase *PersonaGenerationPhaseUsecase) PausePersonaGenerationPhase(
	ctx context.Context,
	request PausePersonaGenerationPhaseRequest,
) (PersonaGenerationPhaseCommandResult, error) {
	if result, allowed := usecase.evaluatePersonaPolicy(ctx, request.JobID, request.PhaseRunID, translationjobpolicy.OperationPause); !allowed {
		return result, nil
	}
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
	if result, allowed := usecase.evaluatePersonaPolicy(ctx, request.JobID, request.PhaseRunID, translationjobpolicy.OperationResume); !allowed {
		return result, nil
	}
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
	if result, allowed := usecase.evaluatePersonaPolicy(ctx, request.JobID, request.PhaseRunID, translationjobpolicy.OperationRetry); !allowed {
		return result, nil
	}
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
	if result, allowed := usecase.evaluatePersonaPolicy(ctx, request.JobID, request.PhaseRunID, translationjobpolicy.OperationCancel); !allowed {
		return result, nil
	}
	readModel, err := usecase.service.CancelPhase(ctx, request.JobID, request.PhaseRunID)
	phaseRunID := request.PhaseRunID
	logPhaseStateCommand(ctx, "persona_generation_phase_cancel", personaGenerationPhaseUsecaseLogWhere, request.JobID, &phaseRunID, readModel.BeforePhaseState, readModel.AfterPhaseState, personaGenerationReadModelErrorReason(readModel.ErrorSummary), err)
	if err != nil {
		return PersonaGenerationPhaseCommandResult{}, fmt.Errorf("cancel persona generation phase: %w", err)
	}
	return toPersonaGenerationPhaseCommandResult(readModel), nil
}

// GetPersonaGenerationBodyReadiness returns the downstream phase fact state for the body phase boundary.
func (usecase *PersonaGenerationPhaseUsecase) GetPersonaGenerationBodyReadiness(
	ctx context.Context,
	request GetPersonaGenerationBodyReadinessRequest,
) (PersonaGenerationBodyReadinessResult, error) {
	readModel, err := usecase.service.ReadBodyReadiness(ctx, request.JobID)
	if err != nil {
		return PersonaGenerationBodyReadinessResult{}, fmt.Errorf("get persona generation body readiness: %w", err)
	}
	return PersonaGenerationBodyReadinessResult{
		JobID:         readModel.JobID,
		CurrentPhase:  readModel.CurrentPhase,
		PhaseState:    readModel.PhaseState,
		JobIsTerminal: readModel.JobIsTerminal,
		InputSummary: PersonaGenerationBodyReadinessInputSummary{
			PersonaCount:   readModel.InputSummary.PersonaCount,
			MissingCount:   readModel.InputSummary.MissingCount,
			SnapshotID:     readModel.InputSummary.SnapshotID,
			SnapshotDigest: readModel.InputSummary.SnapshotDigest,
			EvidenceRefs:   append([]string(nil), readModel.InputSummary.EvidenceRefs...),
		},
	}, nil
}

func personaGenerationReadModelErrorReason(readModel *service.PersonaGenerationPhaseErrorSummaryReadModel) string {
	if readModel == nil {
		return ""
	}
	return readModel.ErrorKind
}

func (usecase *PersonaGenerationPhaseUsecase) evaluatePersonaPolicy(
	ctx context.Context,
	jobID int64,
	phaseRunID int64,
	operation translationjobpolicy.Operation,
) (PersonaGenerationPhaseCommandResult, bool) {
	summary, err := usecase.service.ReadSummary(ctx, jobID)
	if err != nil {
		return PersonaGenerationPhaseCommandResult{}, true
	}
	if personaPolicySummaryMissing(summary) {
		return PersonaGenerationPhaseCommandResult{}, true
	}
	decision := translationjobpolicy.Evaluate(phasePolicyInput(phasePolicyInputSource{
		Operation:            operation,
		JobState:             summary.JobState,
		PhaseState:           summary.PhaseState,
		PhaseRunID:           summary.PhaseRunID,
		RequestedPhaseRunID:  phaseRunID,
		StartPrerequisiteMet: summary.ActionEnablement.CanStart,
	}))
	if decision.Allowed {
		return PersonaGenerationPhaseCommandResult{}, true
	}
	result := personaCommandResultFromSummary(summary)
	result.ErrorSummary = personaPolicyErrorSummary(decision, summary, operation)
	logPhasePolicyRejected(ctx, phasePolicyLogEvent("persona_generation_phase", operation), personaGenerationPhaseUsecaseLogWhere, jobID, phaseRunID, operation, result.ErrorSummary.Reason)
	return result, false
}

func personaPolicySummaryMissing(summary service.PersonaGenerationPhaseSummaryReadModel) bool {
	return summary.JobID == 0 && summary.CurrentPhase == "" && summary.PhaseState == ""
}

func personaPolicyErrorSummary(
	decision translationjobpolicy.Result,
	summary service.PersonaGenerationPhaseSummaryReadModel,
	operation translationjobpolicy.Operation,
) *PersonaGenerationPhaseErrorSummary {
	return &PersonaGenerationPhaseErrorSummary{
		ErrorKind:  personaPolicyErrorKind(decision),
		Reason:     personaPolicyReason(decision, summary, operation),
		Retryable:  false,
		IsRedacted: false,
	}
}

func personaPolicyErrorKind(decision translationjobpolicy.Result) string {
	switch decision.Reason {
	case translationjobpolicy.ReasonTerminalJob:
		return PersonaGenerationPhaseErrorKindTerminalJob
	case translationjobpolicy.ReasonActivePhaseExists:
		return PersonaGenerationPhaseErrorKindActivePhaseExists
	case translationjobpolicy.ReasonStartPrerequisite:
		return PersonaGenerationPhaseErrorKindTermPhaseIncomplete
	default:
		return PersonaGenerationPhaseErrorKindTermPhaseIncomplete
	}
}

func personaPolicyReason(
	decision translationjobpolicy.Result,
	summary service.PersonaGenerationPhaseSummaryReadModel,
	operation translationjobpolicy.Operation,
) string {
	switch operation {
	case translationjobpolicy.OperationStart:
		return stringFromPointer(summary.ActionEnablement.StartBlockedReason, string(decision.Reason))
	case translationjobpolicy.OperationPause:
		return stringFromPointer(summary.ActionEnablement.PauseBlockedReason, string(decision.Reason))
	case translationjobpolicy.OperationResume:
		return stringFromPointer(summary.ActionEnablement.ResumeBlockedReason, string(decision.Reason))
	case translationjobpolicy.OperationRetry:
		return stringFromPointer(summary.ActionEnablement.RetryBlockedReason, string(decision.Reason))
	case translationjobpolicy.OperationCancel:
		return stringFromPointer(summary.ActionEnablement.CancelBlockedReason, string(decision.Reason))
	default:
		return string(decision.Reason)
	}
}

func personaCommandResultFromSummary(summary service.PersonaGenerationPhaseSummaryReadModel) PersonaGenerationPhaseCommandResult {
	return PersonaGenerationPhaseCommandResult{
		JobID:         summary.JobID,
		CurrentPhase:  summary.CurrentPhase,
		PhaseState:    summary.PhaseState,
		PhaseRunID:    clonePersonaGenerationInt64Pointer(summary.PhaseRunID),
		StartedAt:     clonePersonaGenerationTimePointer(summary.StartedAt),
		FinishedAt:    clonePersonaGenerationTimePointer(summary.FinishedAt),
		Progress:      toPersonaGenerationPhaseProgressSummary(summary.Progress),
		TargetSummary: toPersonaGenerationTargetSummary(summary.TargetSummary),
		Execution:     toPersonaGenerationExecutionSummary(summary.Execution),
		ResultSummary: toPersonaGenerationPhaseResultSummary(summary.ResultSummary),
		Retryable:     false,
		ErrorSummary:  toPersonaGenerationPhaseErrorSummary(summary.ErrorSummary),
	}
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
		ResultSummary: toPersonaGenerationPhaseResultSummary(readModel.ResultSummary),
		Retryable:     readModel.Retryable,
		ErrorSummary:  toPersonaGenerationPhaseErrorSummary(readModel.ErrorSummary),
	}
}

func toPersonaGenerationPhaseProgressSummary(
	readModel service.PersonaGenerationPhaseProgressReadModel,
) PersonaGenerationPhaseProgressSummary {
	return PersonaGenerationPhaseProgressSummary{
		Percent:        readModel.Percent,
		ProcessedCount: readModel.ProcessedCount,
		TotalCount:     readModel.TotalCount,
		TargetCount:    readModel.TargetCount,
		CurrentStep:    readModel.CurrentStep,
	}
}

func toPersonaGenerationTargetSummary(
	readModel service.PersonaGenerationTargetSummaryReadModel,
) PersonaGenerationTargetSummary {
	return PersonaGenerationTargetSummary{
		TargetCount:            readModel.TargetCount,
		CommonPersonaHitCount:  readModel.CommonPersonaHitCount,
		CommonPersonaMissCount: readModel.CommonPersonaMissCount,
		SkippedCount:           readModel.SkippedCount,
		SkippedReasons:         append([]string(nil), readModel.SkippedReasons...),
		TargetSnapshotID:       clonePersonaGenerationStringPointer(readModel.TargetSnapshotID),
		TargetSnapshotDigest:   readModel.TargetSnapshotDigest,
	}
}

func toPersonaGenerationExecutionSummary(
	readModel service.PersonaGenerationExecutionSummaryReadModel,
) PersonaGenerationExecutionSummary {
	return PersonaGenerationExecutionSummary{
		CredentialRef: readModel.CredentialRef,
		Provider:      readModel.Provider,
		Model:         readModel.Model,
		ExecutionMode: readModel.ExecutionMode,
		PromptDigest:  readModel.PromptDigest,
		InputCount:    readModel.InputCount,
		OutputCount:   readModel.OutputCount,
		EvidenceRefs:  append([]string(nil), readModel.EvidenceRefs...),
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
		CanStart:            readModel.CanStart,
		StartBlockedReason:  clonePersonaGenerationStringPointer(readModel.StartBlockedReason),
		CanPause:            readModel.CanPause,
		PauseBlockedReason:  clonePersonaGenerationStringPointer(readModel.PauseBlockedReason),
		CanResume:           readModel.CanResume,
		ResumeBlockedReason: clonePersonaGenerationStringPointer(readModel.ResumeBlockedReason),
		CanRetry:            readModel.CanRetry,
		RetryBlockedReason:  clonePersonaGenerationStringPointer(readModel.RetryBlockedReason),
		CanCancel:           readModel.CanCancel,
		CancelBlockedReason: clonePersonaGenerationStringPointer(readModel.CancelBlockedReason),
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
