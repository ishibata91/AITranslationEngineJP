package usecase

import (
	"context"
	"fmt"

	"aitranslationenginejp/internal/service"
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
	readModel, err := usecase.service.StartPhase(ctx, request.JobID)
	result := toBodyTranslationPhaseCommandResult(readModel)
	if err != nil {
		return result, fmt.Errorf("start body translation phase: %w", err)
	}
	return result, nil
}

// PauseBodyTranslationPhase returns the body phase command payload for pause handling.
func (usecase *BodyTranslationPhaseUsecase) PauseBodyTranslationPhase(
	ctx context.Context,
	request PauseBodyTranslationPhaseRequest,
) (BodyTranslationPhaseCommandResult, error) {
	readModel, err := usecase.service.PausePhase(ctx, request.JobID, request.PhaseRunID)
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
	readModel, err := usecase.service.ResumePhase(ctx, request.JobID, request.PhaseRunID)
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
	readModel, err := usecase.service.RetryPhase(ctx, request.JobID, request.PhaseRunID)
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
	readModel, err := usecase.service.CancelPhase(ctx, request.JobID, request.PhaseRunID)
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
