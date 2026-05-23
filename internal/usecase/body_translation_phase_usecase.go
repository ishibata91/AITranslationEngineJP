package usecase

import (
	"context"
	"fmt"

	"aitranslationenginejp/internal/service"
	"aitranslationenginejp/internal/usecase/translationjobpolicy"
)

type bodyTranslationPhaseServicePort interface {
	ReadSummary(ctx context.Context, jobID int64) (service.BodyTranslationPhaseSummaryReadModel, error)
	StartPhase(ctx context.Context, jobID int64) (service.BodyTranslationPhaseCommandReadModel, error)
	PausePhase(ctx context.Context, jobID int64, phaseRunID int64) (service.BodyTranslationPhaseCommandReadModel, error)
	ResumePhase(ctx context.Context, jobID int64, phaseRunID int64) (service.BodyTranslationPhaseCommandReadModel, error)
	RetryPhase(ctx context.Context, jobID int64, phaseRunID int64) (service.BodyTranslationPhaseCommandReadModel, error)
	CancelPhase(ctx context.Context, jobID int64, phaseRunID int64) (service.BodyTranslationPhaseCommandReadModel, error)
	ReadOutputReadiness(ctx context.Context, jobID int64) (service.BodyTranslationOutputReadinessReadModel, error)
}

type bodyTranslationPhaseAISettingsSaver interface {
	SaveAISettings(ctx context.Context, selection service.PhaseAISettingsSelection) (service.PhaseAISettingsReadModel, error)
}

const bodyTranslationPhaseUsecaseLogWhere = "backend.usecase.body_translation_phase"

// BodyTranslationPhaseUsecase bridges the frozen body translation contract to the service layer.
type BodyTranslationPhaseUsecase struct {
	service bodyTranslationPhaseServicePort
}

// NewBodyTranslationPhaseUsecase creates the body translation phase usecase.
func NewBodyTranslationPhaseUsecase(service bodyTranslationPhaseServicePort) *BodyTranslationPhaseUsecase {
	return &BodyTranslationPhaseUsecase{service: service}
}

// GetBodyTranslationPhaseSummary returns the current body phase summary contract.
func (usecase *BodyTranslationPhaseUsecase) GetBodyTranslationPhaseSummary(
	ctx context.Context,
	request GetBodyTranslationPhaseSummaryRequest,
) (BodyTranslationPhaseSummaryResult, error) {
	readModel, err := usecase.service.ReadSummary(ctx, request.JobID)
	if err != nil {
		return BodyTranslationPhaseSummaryResult{}, fmt.Errorf("get body translation phase summary: %w", err)
	}
	return toBodyTranslationPhaseSummaryResult(readModel), nil
}

// StartBodyTranslationPhase starts one body phase run and returns the command contract.
func (usecase *BodyTranslationPhaseUsecase) StartBodyTranslationPhase(
	ctx context.Context,
	request StartBodyTranslationPhaseRequest,
) (BodyTranslationPhaseCommandResult, error) {
	if result, allowed := usecase.evaluateBodyPolicy(ctx, request.JobID, 0, translationjobpolicy.OperationStart); !allowed {
		return result, nil
	}
	readModel, err := usecase.service.StartPhase(ctx, request.JobID)
	logPhaseStateCommand(ctx, "body_translation_phase_start", bodyTranslationPhaseUsecaseLogWhere, request.JobID, readModel.PhaseRunID, readModel.BeforePhaseState, readModel.AfterPhaseState, bodyTranslationReadModelErrorReason(readModel.ErrorSummary), err)
	result := toBodyTranslationPhaseCommandResult(readModel)
	if err != nil {
		return result, fmt.Errorf("start body translation phase: %w", err)
	}
	return result, nil
}

// SaveBodyTranslationPhaseAISettings saves public AI settings for the body phase.
func (usecase *BodyTranslationPhaseUsecase) SaveBodyTranslationPhaseAISettings(
	ctx context.Context,
	request SaveBodyTranslationPhaseAISettingsRequest,
) (BodyTranslationPhaseAISettingsResult, error) {
	saver, ok := usecase.service.(bodyTranslationPhaseAISettingsSaver)
	if !ok {
		return BodyTranslationPhaseAISettingsResult{}, fmt.Errorf("save body translation phase ai settings: service is not configured")
	}
	readModel, err := saver.SaveAISettings(ctx, service.PhaseAISettingsSelection{
		JobID:         request.JobID,
		Provider:      request.Provider,
		Model:         request.Model,
		ExecutionMode: request.ExecutionMode,
		BatchMode:     request.BatchMode,
	})
	if err != nil {
		return BodyTranslationPhaseAISettingsResult{}, fmt.Errorf("save body translation phase ai settings: %w", err)
	}
	return BodyTranslationPhaseAISettingsResult(readModel), nil
}

// PauseBodyTranslationPhase returns the body phase command payload for pause handling.
func (usecase *BodyTranslationPhaseUsecase) PauseBodyTranslationPhase(
	ctx context.Context,
	request PauseBodyTranslationPhaseRequest,
) (BodyTranslationPhaseCommandResult, error) {
	if result, allowed := usecase.evaluateBodyPolicy(ctx, request.JobID, request.PhaseRunID, translationjobpolicy.OperationPause); !allowed {
		return result, nil
	}
	readModel, err := usecase.service.PausePhase(ctx, request.JobID, request.PhaseRunID)
	phaseRunID := request.PhaseRunID
	logPhaseStateCommand(ctx, "body_translation_phase_pause", bodyTranslationPhaseUsecaseLogWhere, request.JobID, &phaseRunID, readModel.BeforePhaseState, readModel.AfterPhaseState, bodyTranslationReadModelErrorReason(readModel.ErrorSummary), err)
	result := toBodyTranslationPhaseCommandResult(readModel)
	if err != nil {
		return result, fmt.Errorf("pause body translation phase: %w", err)
	}
	return result, nil
}

// ResumeBodyTranslationPhase returns the body phase command payload for resume handling.
func (usecase *BodyTranslationPhaseUsecase) ResumeBodyTranslationPhase(
	ctx context.Context,
	request ResumeBodyTranslationPhaseRequest,
) (BodyTranslationPhaseCommandResult, error) {
	if result, allowed := usecase.evaluateBodyPolicy(ctx, request.JobID, request.PhaseRunID, translationjobpolicy.OperationResume); !allowed {
		return result, nil
	}
	readModel, err := usecase.service.ResumePhase(ctx, request.JobID, request.PhaseRunID)
	phaseRunID := request.PhaseRunID
	logPhaseStateCommand(ctx, "body_translation_phase_resume", bodyTranslationPhaseUsecaseLogWhere, request.JobID, &phaseRunID, readModel.BeforePhaseState, readModel.AfterPhaseState, bodyTranslationReadModelErrorReason(readModel.ErrorSummary), err)
	result := toBodyTranslationPhaseCommandResult(readModel)
	if err != nil {
		return result, fmt.Errorf("resume body translation phase: %w", err)
	}
	return result, nil
}

// RetryBodyTranslationPhase returns the body phase command payload for retry handling.
func (usecase *BodyTranslationPhaseUsecase) RetryBodyTranslationPhase(
	ctx context.Context,
	request RetryBodyTranslationPhaseRequest,
) (BodyTranslationPhaseCommandResult, error) {
	if result, allowed := usecase.evaluateBodyPolicy(ctx, request.JobID, request.PhaseRunID, translationjobpolicy.OperationRetry); !allowed {
		return result, nil
	}
	readModel, err := usecase.service.RetryPhase(ctx, request.JobID, request.PhaseRunID)
	phaseRunID := request.PhaseRunID
	logPhaseStateCommand(ctx, "body_translation_phase_retry", bodyTranslationPhaseUsecaseLogWhere, request.JobID, &phaseRunID, readModel.BeforePhaseState, readModel.AfterPhaseState, bodyTranslationReadModelErrorReason(readModel.ErrorSummary), err)
	result := toBodyTranslationPhaseCommandResult(readModel)
	if err != nil {
		return result, fmt.Errorf("retry body translation phase: %w", err)
	}
	return result, nil
}

// CancelBodyTranslationPhase returns the body phase command payload for cancel handling.
func (usecase *BodyTranslationPhaseUsecase) CancelBodyTranslationPhase(
	ctx context.Context,
	request CancelBodyTranslationPhaseRequest,
) (BodyTranslationPhaseCommandResult, error) {
	if result, allowed := usecase.evaluateBodyPolicy(ctx, request.JobID, request.PhaseRunID, translationjobpolicy.OperationCancel); !allowed {
		return result, nil
	}
	readModel, err := usecase.service.CancelPhase(ctx, request.JobID, request.PhaseRunID)
	phaseRunID := request.PhaseRunID
	logPhaseStateCommand(ctx, "body_translation_phase_cancel", bodyTranslationPhaseUsecaseLogWhere, request.JobID, &phaseRunID, readModel.BeforePhaseState, readModel.AfterPhaseState, bodyTranslationReadModelErrorReason(readModel.ErrorSummary), err)
	result := toBodyTranslationPhaseCommandResult(readModel)
	if err != nil {
		return result, fmt.Errorf("cancel body translation phase: %w", err)
	}
	return result, nil
}

// GetBodyTranslationOutputReadiness returns whether downstream output is ready.
func (usecase *BodyTranslationPhaseUsecase) GetBodyTranslationOutputReadiness(
	ctx context.Context,
	request GetBodyTranslationOutputReadinessRequest,
) (BodyTranslationOutputReadinessResult, error) {
	readModel, err := usecase.service.ReadOutputReadiness(ctx, request.JobID)
	logPhaseReadiness(ctx, "body_translation_output_readiness", bodyTranslationPhaseUsecaseLogWhere, request.JobID, readModel.PhaseState, readModel.Ready, stringPointerIfNotBlank(readModel.BlockedReason), err)
	result := BodyTranslationOutputReadinessResult{
		JobID:               readModel.JobID,
		CurrentPhase:        readModel.CurrentPhase,
		PhaseState:          readModel.PhaseState,
		Ready:               readModel.Ready,
		BlockedReason:       readModel.BlockedReason,
		ErrorKind:           NormalizeBodyTranslationPhasePublicErrorKind(BodyTranslationPhaseErrorKind(readModel.ErrorKind)),
		CompletedFieldCount: readModel.CompletedFieldCount,
		StatusConsistent:    readModel.StatusConsistent,
		OutputCount:         readModel.OutputCount,
	}
	if err != nil {
		return result, fmt.Errorf("get body translation output readiness: %w", err)
	}
	return result, nil
}

func stringPointerIfNotBlank(value string) *string {
	if value == "" {
		return nil
	}
	return &value
}

func bodyTranslationReadModelErrorReason(readModel *service.BodyTranslationPhaseErrorSummaryReadModel) string {
	if readModel == nil {
		return ""
	}
	return readModel.ErrorKind
}

func (usecase *BodyTranslationPhaseUsecase) evaluateBodyPolicy(
	ctx context.Context,
	jobID int64,
	phaseRunID int64,
	operation translationjobpolicy.Operation,
) (BodyTranslationPhaseCommandResult, bool) {
	summary, err := usecase.service.ReadSummary(ctx, jobID)
	if err != nil {
		return BodyTranslationPhaseCommandResult{}, true
	}
	if bodyPolicySummaryMissing(summary) {
		return BodyTranslationPhaseCommandResult{}, true
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
		return BodyTranslationPhaseCommandResult{}, true
	}
	result := bodyCommandResultFromSummary(summary)
	result.ErrorSummary = bodyPolicyErrorSummary(decision, summary, operation)
	logPhasePolicyRejected(ctx, phasePolicyLogEvent("body_translation_phase", operation), bodyTranslationPhaseUsecaseLogWhere, jobID, phaseRunID, operation, result.ErrorSummary.Reason)
	return result, false
}

func bodyPolicySummaryMissing(summary service.BodyTranslationPhaseSummaryReadModel) bool {
	return summary.JobID == 0 && summary.CurrentPhase == "" && summary.PhaseState == ""
}

func bodyPolicyErrorSummary(
	decision translationjobpolicy.Result,
	summary service.BodyTranslationPhaseSummaryReadModel,
	operation translationjobpolicy.Operation,
) *BodyTranslationPhaseErrorSummary {
	return &BodyTranslationPhaseErrorSummary{
		ErrorKind:  BodyTranslationPhaseErrorKind(bodyPolicyErrorKind(decision)),
		Reason:     bodyPolicyReason(decision, summary, operation),
		Retryable:  false,
		IsRedacted: false,
	}
}

func bodyPolicyErrorKind(decision translationjobpolicy.Result) string {
	switch decision.Reason {
	case translationjobpolicy.ReasonTerminalJob:
		return string(BodyTranslationPhaseErrorKindTerminalJob)
	case translationjobpolicy.ReasonActivePhaseExists:
		return string(BodyTranslationPhaseErrorKindActivePhaseExists)
	case translationjobpolicy.ReasonStartPrerequisite:
		return string(BodyTranslationPhaseErrorKindPersonaPhaseIncomplete)
	default:
		return string(BodyTranslationPhaseErrorKindPersonaPhaseIncomplete)
	}
}

func bodyPolicyReason(
	decision translationjobpolicy.Result,
	summary service.BodyTranslationPhaseSummaryReadModel,
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

func bodyCommandResultFromSummary(summary service.BodyTranslationPhaseSummaryReadModel) BodyTranslationPhaseCommandResult {
	return BodyTranslationPhaseCommandResult{
		JobID:              summary.JobID,
		CurrentPhase:       summary.CurrentPhase,
		PhaseState:         summary.PhaseState,
		PhaseRunID:         cloneInt64Pointer(summary.PhaseRunID),
		StartedAt:          cloneTimePointer(summary.StartedAt),
		FinishedAt:         cloneTimePointer(summary.FinishedAt),
		Progress:           toBodyTranslationPhaseProgressSummary(summary.Progress),
		InputSummary:       toBodyTranslationPhaseInputSummary(summary.InputSummary),
		RequestSummary:     toBodyTranslationPhaseRequestSummary(summary.RequestSummary),
		Execution:          toBodyTranslationPhaseExecutionSummary(summary.Execution),
		FieldResultSummary: toOptionalBodyTranslationFieldResultSummary(summary.FieldResultSummary),
		ResultSummary:      toOptionalBodyTranslationFieldResultSummary(summary.ResultSummary),
		FieldResults:       toBodyTranslationFieldResultItems(summary.FieldResults),
		Retryable:          false,
		OutputReadiness:    toBodyTranslationOutputReadinessSummary(summary.OutputReadiness),
		ErrorSummary:       toOptionalBodyTranslationErrorSummary(summary.ErrorSummary),
	}
}

func toBodyTranslationPhaseSummaryResult(
	readModel service.BodyTranslationPhaseSummaryReadModel,
) BodyTranslationPhaseSummaryResult {
	return BodyTranslationPhaseSummaryResult{
		JobID:              readModel.JobID,
		CurrentPhase:       readModel.CurrentPhase,
		PhaseState:         readModel.PhaseState,
		PhaseRunID:         cloneInt64Pointer(readModel.PhaseRunID),
		StartedAt:          cloneTimePointer(readModel.StartedAt),
		FinishedAt:         cloneTimePointer(readModel.FinishedAt),
		Progress:           toBodyTranslationPhaseProgressSummary(readModel.Progress),
		InputSummary:       toBodyTranslationPhaseInputSummary(readModel.InputSummary),
		RequestSummary:     toBodyTranslationPhaseRequestSummary(readModel.RequestSummary),
		Execution:          toBodyTranslationPhaseExecutionSummary(readModel.Execution),
		FieldResultSummary: toOptionalBodyTranslationFieldResultSummary(readModel.FieldResultSummary),
		ResultSummary:      toOptionalBodyTranslationFieldResultSummary(readModel.ResultSummary),
		FieldResults:       toBodyTranslationFieldResultItems(readModel.FieldResults),
		ErrorSummary:       toOptionalBodyTranslationErrorSummary(readModel.ErrorSummary),
		ActionEnablement:   toBodyTranslationPhaseActionEnablement(readModel.ActionEnablement),
		OutputReadiness:    toBodyTranslationOutputReadinessSummary(readModel.OutputReadiness),
	}
}

func toBodyTranslationPhaseCommandResult(
	readModel service.BodyTranslationPhaseCommandReadModel,
) BodyTranslationPhaseCommandResult {
	return BodyTranslationPhaseCommandResult{
		JobID:               readModel.JobID,
		CurrentPhase:        readModel.CurrentPhase,
		PhaseState:          readModel.PhaseState,
		PhaseRunID:          cloneInt64Pointer(readModel.PhaseRunID),
		StartedAt:           cloneTimePointer(readModel.StartedAt),
		FinishedAt:          cloneTimePointer(readModel.FinishedAt),
		Progress:            toBodyTranslationPhaseProgressSummary(readModel.Progress),
		InputSnapshotDigest: readModel.InputSnapshotDigest,
		InputSummary:        toBodyTranslationPhaseInputSummary(readModel.InputSummary),
		RequestSummary:      toBodyTranslationPhaseRequestSummary(readModel.RequestSummary),
		Execution:           toBodyTranslationPhaseExecutionSummary(readModel.Execution),
		FieldResultSummary:  toOptionalBodyTranslationFieldResultSummary(readModel.FieldResultSummary),
		ResultSummary:       toOptionalBodyTranslationFieldResultSummary(readModel.ResultSummary),
		FieldResults:        toBodyTranslationFieldResultItems(readModel.FieldResults),
		Retryable:           readModel.Retryable,
		OutputReadiness:     toBodyTranslationOutputReadinessSummary(readModel.OutputReadiness),
		ErrorSummary:        toOptionalBodyTranslationErrorSummary(readModel.ErrorSummary),
	}
}

func toBodyTranslationPhaseProgressSummary(
	readModel service.BodyTranslationPhaseProgressReadModel,
) BodyTranslationPhaseProgressSummary {
	return BodyTranslationPhaseProgressSummary{
		Percent:         readModel.Percent,
		ProcessedCount:  readModel.ProcessedCount,
		TotalCount:      readModel.TotalCount,
		TargetCount:     readModel.TargetCount,
		TranslatedCount: readModel.TranslatedCount,
		SkippedCount:    readModel.SkippedCount,
		CurrentStep:     readModel.CurrentStep,
	}
}

func toBodyTranslationPhaseInputSummary(
	readModel service.BodyTranslationPhaseInputSummaryReadModel,
) BodyTranslationPhaseInputSummary {
	return BodyTranslationPhaseInputSummary{
		TargetCount:      readModel.TargetCount,
		SkippedReasons:   append([]string(nil), readModel.SkippedReasons...),
		InputSnapshotRef: cloneStringPointer(readModel.InputSnapshotRef),
		DictionaryDigest: readModel.DictionaryDigest,
		PersonaDigest:    readModel.PersonaDigest,
		MetadataDigest:   readModel.MetadataDigest,
		PromptDigest:     readModel.PromptDigest,
	}
}

func toBodyTranslationPhaseRequestSummary(
	readModel service.BodyTranslationPhaseRequestSummaryReadModel,
) BodyTranslationPhaseRequestSummary {
	return BodyTranslationPhaseRequestSummary{
		ProviderTargetCount:              readModel.ProviderTargetCount,
		ExactDictionaryExclusionCount:    readModel.ExactDictionaryExclusionCount,
		PartialDictionaryConstraintCount: readModel.PartialDictionaryConstraintCount,
	}
}

func toBodyTranslationPhaseExecutionSummary(
	readModel service.BodyTranslationPhaseExecutionSummaryReadModel,
) BodyTranslationPhaseExecutionSummary {
	return BodyTranslationPhaseExecutionSummary{
		CredentialRef:    readModel.CredentialRef,
		Provider:         readModel.Provider,
		Model:            readModel.Model,
		ExecutionMode:    readModel.ExecutionMode,
		RequestUnitCount: readModel.RequestUnitCount,
		OutputCount:      readModel.OutputCount,
	}
}

func toOptionalBodyTranslationFieldResultSummary(
	readModel *service.BodyTranslationPhaseFieldResultSummaryReadModel,
) *BodyTranslationPhaseFieldResultSummary {
	if readModel == nil {
		return nil
	}
	return &BodyTranslationPhaseFieldResultSummary{
		TranslatedCount:       readModel.TranslatedCount,
		FailedCount:           readModel.FailedCount,
		SkippedCount:          readModel.SkippedCount,
		ProtectionFailedCount: readModel.ProtectionFailedCount,
		OutputReadyCount:      readModel.OutputReadyCount,
		OutputCount:           readModel.OutputCount,
		FieldResults:          toBodyTranslationFieldResultItems(readModel.FieldResults),
	}
}

func toBodyTranslationFieldResultItems(
	readModels []service.BodyTranslationPhaseFieldResultItemReadModel,
) []BodyTranslationPhaseFieldResultItem {
	if len(readModels) == 0 {
		return nil
	}
	items := make([]BodyTranslationPhaseFieldResultItem, 0, len(readModels))
	for _, readModel := range readModels {
		items = append(items, BodyTranslationPhaseFieldResultItem{
			Identity: BodyTranslationPhaseFieldIdentity{
				TranslationFieldID:      readModel.Identity.TranslationFieldID,
				PhaseTranslationFieldID: readModel.Identity.PhaseTranslationFieldID,
				RecordType:              readModel.Identity.RecordType,
				FieldType:               readModel.Identity.FieldType,
				FormID:                  readModel.Identity.FormID,
				EditorID:                readModel.Identity.EditorID,
				FieldLabel:              readModel.Identity.FieldLabel,
			},
			FieldID:                     readModel.FieldID,
			FieldLabel:                  readModel.FieldLabel,
			SourceExcerpt:               readModel.SourceExcerpt,
			TranslatedText:              readModel.TranslatedText,
			OutputStatus:                readModel.OutputStatus,
			ProtectionValidationResult:  readModel.ProtectionValidationResult,
			ProtectionValidationSummary: readModel.ProtectionValidationSummary,
			RetryCount:                  readModel.RetryCount,
		})
	}
	return items
}

func toOptionalBodyTranslationErrorSummary(
	readModel *service.BodyTranslationPhaseErrorSummaryReadModel,
) *BodyTranslationPhaseErrorSummary {
	if readModel == nil {
		return nil
	}
	return &BodyTranslationPhaseErrorSummary{
		ErrorKind:  NormalizeBodyTranslationPhasePublicErrorKind(BodyTranslationPhaseErrorKind(readModel.ErrorKind)),
		Reason:     readModel.Reason,
		Retryable:  readModel.Retryable,
		IsRedacted: readModel.IsRedacted,
	}
}

func toBodyTranslationPhaseActionEnablement(
	readModel service.BodyTranslationPhaseActionEnablementReadModel,
) BodyTranslationPhaseActionEnablement {
	return BodyTranslationPhaseActionEnablement{
		CanStart:                     readModel.CanStart,
		StartBlockedReason:           cloneStringPointer(readModel.StartBlockedReason),
		CanPause:                     readModel.CanPause,
		PauseBlockedReason:           cloneStringPointer(readModel.PauseBlockedReason),
		CanResume:                    readModel.CanResume,
		ResumeBlockedReason:          cloneStringPointer(readModel.ResumeBlockedReason),
		CanRetry:                     readModel.CanRetry,
		RetryBlockedReason:           cloneStringPointer(readModel.RetryBlockedReason),
		CanCancel:                    readModel.CanCancel,
		CancelBlockedReason:          cloneStringPointer(readModel.CancelBlockedReason),
		CanCheckOutputReadiness:      readModel.CanCheckOutputReadiness,
		OutputReadinessBlockedReason: cloneStringPointer(readModel.OutputReadinessBlockedReason),
	}
}

func toBodyTranslationOutputReadinessSummary(
	readModel service.BodyTranslationOutputReadinessReadModel,
) BodyTranslationOutputReadinessSummary {
	return BodyTranslationOutputReadinessSummary{
		Ready:               readModel.Ready,
		BlockedReason:       readModel.BlockedReason,
		ErrorKind:           NormalizeBodyTranslationPhasePublicErrorKind(BodyTranslationPhaseErrorKind(readModel.ErrorKind)),
		CompletedFieldCount: readModel.CompletedFieldCount,
		StatusConsistent:    readModel.StatusConsistent,
	}
}
