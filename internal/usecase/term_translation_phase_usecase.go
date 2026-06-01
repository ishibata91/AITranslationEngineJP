package usecase

import (
	"context"
	"fmt"
	"time"

	"aitranslationenginejp/internal/service"
	"aitranslationenginejp/internal/usecase/translationjobpolicy"
)

type termTranslationPhaseServicePort interface {
	ReadSummary(ctx context.Context, jobID int64) (service.TermTranslationPhaseSummaryReadModel, error)
	StartPhase(ctx context.Context, jobID int64) (service.TermTranslationPhaseCommandReadModel, error)
	PausePhase(ctx context.Context, jobID int64, phaseRunID int64) (service.TermTranslationPhaseCommandReadModel, error)
	ResumePhase(ctx context.Context, jobID int64, phaseRunID int64) (service.TermTranslationPhaseCommandReadModel, error)
	RetryPhase(ctx context.Context, jobID int64, phaseRunID int64) (service.TermTranslationPhaseCommandReadModel, error)
	ReadNextPhaseReadiness(ctx context.Context, jobID int64) (service.TermTranslationNextPhaseReadinessReadModel, error)
}

type termTranslationPhaseAISettingsSaver interface {
	SaveAISettings(ctx context.Context, selection service.PhaseAISettingsSelection) (service.PhaseAISettingsReadModel, error)
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
	if result, allowed := usecase.evaluateTermPolicy(ctx, request.JobID, 0, translationjobpolicy.OperationStart); !allowed {
		return result, nil
	}
	readModel, err := usecase.service.StartPhase(ctx, request.JobID)
	logPhaseStateCommand(ctx, "term_translation_phase_start", termTranslationPhaseUsecaseLogWhere, request.JobID, readModel.PhaseRunID, readModel.BeforePhaseState, readModel.AfterPhaseState, termTranslationReadModelErrorReason(readModel.ErrorSummary), err)
	if err != nil {
		return TermTranslationPhaseCommandResult{}, fmt.Errorf("start term translation phase: %w", err)
	}
	return toTermTranslationPhaseCommandResult(readModel), nil
}

// SaveTermTranslationPhaseAISettings saves public AI settings for the term phase.
// 入力に job_id を含めない。JOB_PHASE_AI_SETTINGS の phase_type = "word_translation" へ upsert する。
func (usecase *TermTranslationPhaseUsecase) SaveTermTranslationPhaseAISettings(
	ctx context.Context,
	request SaveTermTranslationPhaseAISettingsRequest,
) (TermTranslationPhaseAISettingsResult, error) {
	saver, ok := usecase.service.(termTranslationPhaseAISettingsSaver)
	if !ok {
		return TermTranslationPhaseAISettingsResult{}, fmt.Errorf("save term translation phase ai settings: service is not configured")
	}
	readModel, err := saver.SaveAISettings(ctx, service.PhaseAISettingsSelection{
		Provider:      request.Provider,
		Model:         request.Model,
		ExecutionMode: request.ExecutionMode,
		BatchMode:     request.BatchMode,
	})
	if err != nil {
		return TermTranslationPhaseAISettingsResult{}, fmt.Errorf("save term translation phase ai settings: %w", err)
	}
	return TermTranslationPhaseAISettingsResult{
		PhaseType:     readModel.PhaseType,
		Provider:      readModel.Provider,
		Model:         readModel.Model,
		ExecutionMode: readModel.ExecutionMode,
		BatchMode:     readModel.BatchMode,
	}, nil
}

// PauseTermTranslationPhase pauses one active term phase run and returns the command contract.
func (usecase *TermTranslationPhaseUsecase) PauseTermTranslationPhase(
	ctx context.Context,
	request PauseTermTranslationPhaseRequest,
) (TermTranslationPhaseCommandResult, error) {
	if result, allowed := usecase.evaluateTermPolicy(ctx, request.JobID, request.PhaseRunID, translationjobpolicy.OperationPause); !allowed {
		return result, nil
	}
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
	if result, allowed := usecase.evaluateTermPolicy(ctx, request.JobID, request.PhaseRunID, translationjobpolicy.OperationResume); !allowed {
		return result, nil
	}
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
	if result, allowed := usecase.evaluateTermPolicy(ctx, request.JobID, request.PhaseRunID, translationjobpolicy.OperationRetry); !allowed {
		return result, nil
	}
	readModel, err := usecase.service.RetryPhase(ctx, request.JobID, request.PhaseRunID)
	phaseRunID := request.PhaseRunID
	logPhaseStateCommand(ctx, "term_translation_phase_retry", termTranslationPhaseUsecaseLogWhere, request.JobID, &phaseRunID, readModel.BeforePhaseState, readModel.AfterPhaseState, termTranslationReadModelErrorReason(readModel.ErrorSummary), err)
	if err != nil {
		return TermTranslationPhaseCommandResult{}, fmt.Errorf("retry term translation phase: %w", err)
	}
	return toTermTranslationPhaseCommandResult(readModel), nil
}

// GetTermTranslationNextPhaseReadiness returns the downstream phase fact state for the term phase boundary.
func (usecase *TermTranslationPhaseUsecase) GetTermTranslationNextPhaseReadiness(
	ctx context.Context,
	request GetTermTranslationNextPhaseReadinessRequest,
) (TermTranslationNextPhaseReadinessResult, error) {
	readModel, err := usecase.service.ReadNextPhaseReadiness(ctx, request.JobID)
	if err != nil {
		return TermTranslationNextPhaseReadinessResult{}, fmt.Errorf("get term translation next phase readiness: %w", err)
	}
	return TermTranslationNextPhaseReadinessResult{
		JobID:          readModel.JobID,
		CurrentPhase:   readModel.CurrentPhase,
		PhaseState:     readModel.PhaseState,
		JobIsTerminal:  readModel.JobIsTerminal,
		TotalCount:     readModel.TotalCount,
		ConfirmedCount: readModel.ConfirmedCount,
		ErrorKind:      NormalizeTermTranslationPhasePublicErrorKind(readModel.ErrorKind),
	}, nil
}

func termTranslationReadModelErrorReason(readModel *service.TermTranslationPhaseErrorReadModel) string {
	if readModel == nil {
		return ""
	}
	return readModel.ErrorKind
}

func (usecase *TermTranslationPhaseUsecase) evaluateTermPolicy(
	ctx context.Context,
	jobID int64,
	phaseRunID int64,
	operation translationjobpolicy.Operation,
) (TermTranslationPhaseCommandResult, bool) {
	summary, err := usecase.service.ReadSummary(ctx, jobID)
	if err != nil {
		return TermTranslationPhaseCommandResult{}, true
	}
	if termPolicySummaryMissing(summary) {
		return TermTranslationPhaseCommandResult{}, true
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
		return TermTranslationPhaseCommandResult{}, true
	}
	result := termCommandResultFromSummary(summary)
	result.ErrorSummary = termPolicyErrorSummary(decision, summary, operation)
	logPhasePolicyRejected(ctx, phasePolicyLogEvent("term_translation_phase", operation), termTranslationPhaseUsecaseLogWhere, jobID, phaseRunID, operation, result.ErrorSummary.Reason)
	return result, false
}

func termPolicySummaryMissing(summary service.TermTranslationPhaseSummaryReadModel) bool {
	return summary.JobID == 0 && summary.CurrentPhase == "" && summary.PhaseState == ""
}

func termPolicyErrorSummary(
	decision translationjobpolicy.Result,
	summary service.TermTranslationPhaseSummaryReadModel,
	operation translationjobpolicy.Operation,
) *TermTranslationPhaseErrorSummary {
	reason := termPolicyReason(decision, summary, operation)
	return &TermTranslationPhaseErrorSummary{
		ErrorKind:  termPolicyErrorKind(decision),
		Reason:     reason,
		Retryable:  false,
		IsRedacted: false,
	}
}

func termPolicyErrorKind(decision translationjobpolicy.Result) string {
	switch decision.Reason {
	case translationjobpolicy.ReasonTerminalJob:
		return TermTranslationPhaseErrorKindTerminalJob
	case translationjobpolicy.ReasonActivePhaseExists:
		return TermTranslationPhaseErrorKindActivePhaseExists
	case translationjobpolicy.ReasonStartPrerequisite:
		return TermTranslationPhaseErrorKindReadyRequired
	default:
		return TermTranslationPhaseErrorKindTermPhaseIncomplete
	}
}

func termPolicyReason(
	decision translationjobpolicy.Result,
	summary service.TermTranslationPhaseSummaryReadModel,
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
	default:
		return string(decision.Reason)
	}
}

func termCommandResultFromSummary(summary service.TermTranslationPhaseSummaryReadModel) TermTranslationPhaseCommandResult {
	return TermTranslationPhaseCommandResult{
		JobID:        summary.JobID,
		CurrentPhase: summary.CurrentPhase,
		PhaseState:   summary.PhaseState,
		PhaseRunID:   cloneInt64Pointer(summary.PhaseRunID),
		StartedAt:    cloneTimePointer(summary.StartedAt),
		FinishedAt:   cloneTimePointer(summary.FinishedAt),
		Progress:     toTermTranslationPhaseProgressSummary(summary.Progress),
		Retryable:    false,
		ErrorSummary: toTermTranslationPhaseErrorSummary(summary.ErrorSummary),
	}
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
		AISettings:         toTermTranslationPhaseAISettings(readModel.AISettings),
		Execution:          toTermTranslationExecutionConfigSummary(readModel.Execution),
		ResultSummary:      toTermTranslationPhaseResultSummary(readModel.ResultSummary),
		ErrorSummary:       toTermTranslationPhaseErrorSummary(readModel.ErrorSummary),
		ActionEnablement:   toTermTranslationPhaseActionEnablement(readModel.ActionEnablement),
	}
}

func toTermTranslationPhaseAISettings(
	readModel *service.TermTranslationPhaseAISettingsReadModel,
) *TermTranslationPhaseAISettings {
	if readModel == nil {
		return nil
	}
	return &TermTranslationPhaseAISettings{
		Provider:      readModel.Provider,
		Model:         readModel.Model,
		ExecutionMode: readModel.ExecutionMode,
		BatchMode:     readModel.BatchMode,
	}
}

func toTermTranslationExecutionConfigSummary(
	readModel *service.TermTranslationExecutionConfigReadModel,
) *TermTranslationExecutionConfigSummary {
	if readModel == nil {
		return nil
	}
	return &TermTranslationExecutionConfigSummary{
		CredentialRef:   readModel.CredentialRef,
		Provider:        readModel.Provider,
		Model:           readModel.Model,
		ExecutionMode:   readModel.ExecutionMode,
		SnapshotDigest:  cloneStringPointer(readModel.SnapshotDigest),
		SnapshotVersion: cloneStringPointer(readModel.SnapshotVersion),
	}
}

func toTermTranslationPhaseCommandResult(
	readModel service.TermTranslationPhaseCommandReadModel,
) TermTranslationPhaseCommandResult {
	return TermTranslationPhaseCommandResult{
		JobID:        readModel.JobID,
		CurrentPhase: readModel.CurrentPhase,
		PhaseState:   readModel.PhaseState,
		PhaseRunID:   cloneInt64Pointer(readModel.PhaseRunID),
		StartedAt:    cloneTimePointer(readModel.StartedAt),
		FinishedAt:   cloneTimePointer(readModel.FinishedAt),
		Progress:     toTermTranslationPhaseProgressSummary(readModel.Progress),
		Retryable:    readModel.Retryable,
		ErrorSummary: toTermTranslationPhaseErrorSummary(readModel.ErrorSummary),
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
		CanStart:            readModel.CanStart,
		StartBlockedReason:  cloneStringPointer(readModel.StartBlockedReason),
		CanPause:            readModel.CanPause,
		PauseBlockedReason:  cloneStringPointer(readModel.PauseBlockedReason),
		CanResume:           readModel.CanResume,
		ResumeBlockedReason: cloneStringPointer(readModel.ResumeBlockedReason),
		CanRetry:            readModel.CanRetry,
		RetryBlockedReason:  cloneStringPointer(readModel.RetryBlockedReason),
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
