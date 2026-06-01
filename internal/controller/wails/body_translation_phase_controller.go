package wails

import (
	"context"
	"fmt"

	"aitranslationenginejp/internal/usecase"
)

// BodyTranslationPhaseUsecasePort defines the frozen body translation phase usecase seam.
type BodyTranslationPhaseUsecasePort interface {
	GetBodyTranslationPhaseSummary(ctx context.Context, request usecase.GetBodyTranslationPhaseSummaryRequest) (usecase.BodyTranslationPhaseSummaryResult, error)
	StartBodyTranslationPhase(ctx context.Context, request usecase.StartBodyTranslationPhaseRequest) (usecase.BodyTranslationPhaseCommandResult, error)
	PauseBodyTranslationPhase(ctx context.Context, request usecase.PauseBodyTranslationPhaseRequest) (usecase.BodyTranslationPhaseCommandResult, error)
	ResumeBodyTranslationPhase(ctx context.Context, request usecase.ResumeBodyTranslationPhaseRequest) (usecase.BodyTranslationPhaseCommandResult, error)
	RetryBodyTranslationPhase(ctx context.Context, request usecase.RetryBodyTranslationPhaseRequest) (usecase.BodyTranslationPhaseCommandResult, error)
	CancelBodyTranslationPhase(ctx context.Context, request usecase.CancelBodyTranslationPhaseRequest) (usecase.BodyTranslationPhaseCommandResult, error)
}

type bodyTranslationPhaseAISettingsUsecasePort interface {
	SaveBodyTranslationPhaseAISettings(ctx context.Context, request usecase.SaveBodyTranslationPhaseAISettingsRequest) (usecase.BodyTranslationPhaseAISettingsResult, error)
}

// BodyTranslationPhaseController exposes Wails-bound body translation entrypoints.
type BodyTranslationPhaseController struct {
	bodyTranslationPhaseUsecase BodyTranslationPhaseUsecasePort
}

// GetBodyTranslationPhaseSummaryRequestDTO identifies the requested job.
type GetBodyTranslationPhaseSummaryRequestDTO struct {
	JobID int64 `json:"jobId"`
}

// StartBodyTranslationPhaseRequestDTO identifies the ready job to start.
type StartBodyTranslationPhaseRequestDTO struct {
	JobID int64 `json:"jobId"`
}

// PauseBodyTranslationPhaseRequestDTO identifies one active phase run to pause.
type PauseBodyTranslationPhaseRequestDTO struct {
	JobID      int64 `json:"jobId"`
	PhaseRunID int64 `json:"phaseRunId"`
}

// ResumeBodyTranslationPhaseRequestDTO identifies one paused or recoverable failed phase run to resume.
type ResumeBodyTranslationPhaseRequestDTO struct {
	JobID      int64 `json:"jobId"`
	PhaseRunID int64 `json:"phaseRunId"`
}

// RetryBodyTranslationPhaseRequestDTO identifies one retryable failed phase run to retry.
type RetryBodyTranslationPhaseRequestDTO struct {
	JobID      int64 `json:"jobId"`
	PhaseRunID int64 `json:"phaseRunId"`
}

// CancelBodyTranslationPhaseRequestDTO identifies one cancelable phase run.
type CancelBodyTranslationPhaseRequestDTO struct {
	JobID      int64 `json:"jobId"`
	PhaseRunID int64 `json:"phaseRunId"`
}

// SaveBodyTranslationPhaseAISettingsRequestDTO carries public AI settings.
type SaveBodyTranslationPhaseAISettingsRequestDTO struct {
	Provider      string `json:"provider"`
	Model         string `json:"model"`
	ExecutionMode string `json:"executionMode"`
	BatchMode     string `json:"batchMode"`
}

// BodyTranslationPhaseAISettingsResponseDTO returns saved AI settings state.
// credential_status、model_list_status は含まない。
type BodyTranslationPhaseAISettingsResponseDTO struct {
	PhaseType     string `json:"phaseType"`
	Provider      string `json:"provider"`
	Model         string `json:"model"`
	ExecutionMode string `json:"executionMode"`
	BatchMode     string `json:"batchMode"`
}

// BodyTranslationPhaseProgressSummaryDTO summarizes one phase run progress snapshot.
type BodyTranslationPhaseProgressSummaryDTO struct {
	Percent         int    `json:"percent"`
	ProcessedCount  int    `json:"processedCount"`
	TotalCount      int    `json:"totalCount"`
	TargetCount     int    `json:"targetCount"`
	TranslatedCount int    `json:"translatedCount"`
	SkippedCount    int    `json:"skippedCount"`
	CurrentStep     string `json:"currentStep"`
}

// BodyTranslationPhaseInputSummaryDTO summarizes the frozen body phase input snapshot.
type BodyTranslationPhaseInputSummaryDTO struct {
	TargetCount      int      `json:"targetCount"`
	SkippedReasons   []string `json:"skippedReasons"`
	InputSnapshotRef *string  `json:"inputSnapshotRef,omitempty"`
	DictionaryDigest string   `json:"dictionaryDigest"`
	PersonaDigest    string   `json:"personaDigest"`
	MetadataDigest   string   `json:"metadataDigest"`
	PromptDigest     string   `json:"promptDigest"`
}

// BodyTranslationPhaseRequestSummaryDTO summarizes provider request units after input filtering.
type BodyTranslationPhaseRequestSummaryDTO struct {
	ProviderTargetCount              int `json:"providerTargetCount"`
	ExactDictionaryExclusionCount    int `json:"exactDictionaryExclusionCount"`
	PartialDictionaryConstraintCount int `json:"partialDictionaryConstraintCount"`
}

// BodyTranslationPhaseExecutionSummaryDTO summarizes execution inputs without exposing secrets.
type BodyTranslationPhaseExecutionSummaryDTO struct {
	CredentialRef    string `json:"credentialRef"`
	Provider         string `json:"provider"`
	Model            string `json:"model"`
	ExecutionMode    string `json:"executionMode"`
	RequestUnitCount int    `json:"requestUnitCount"`
	OutputCount      int    `json:"outputCount"`
}

// BodyTranslationPhaseAISettingsDTO summarizes the saved AI settings for this phase.
// JOB_PHASE_AI_SETTINGS レコードが存在しない場合は省略する。
type BodyTranslationPhaseAISettingsDTO struct {
	Provider      string `json:"provider"`
	Model         string `json:"model"`
	ExecutionMode string `json:"executionMode"`
	BatchMode     string `json:"batchMode"`
}

// BodyTranslationPhaseFieldResultSummaryDTO summarizes stored field results.
type BodyTranslationPhaseFieldResultSummaryDTO struct {
	TranslatedCount       int                                      `json:"translatedCount"`
	FailedCount           int                                      `json:"failedCount"`
	SkippedCount          int                                      `json:"skippedCount"`
	ProtectionFailedCount int                                      `json:"protectionFailedCount"`
	OutputReadyCount      int                                      `json:"outputReadyCount"`
	OutputCount           int                                      `json:"outputCount"`
	FieldResults          []BodyTranslationPhaseFieldResultItemDTO `json:"fieldResults,omitempty"`
}

// BodyTranslationPhaseFieldIdentityDTO identifies one translated field.
type BodyTranslationPhaseFieldIdentityDTO struct {
	TranslationFieldID      int64  `json:"translationFieldId"`
	PhaseTranslationFieldID int64  `json:"phaseTranslationFieldId"`
	RecordType              string `json:"recordType,omitempty"`
	FieldType               string `json:"fieldType,omitempty"`
	FormID                  string `json:"formId,omitempty"`
	EditorID                string `json:"editorId,omitempty"`
	FieldLabel              string `json:"fieldLabel,omitempty"`
}

// BodyTranslationPhaseFieldResultItemDTO returns one public field result row.
type BodyTranslationPhaseFieldResultItemDTO struct {
	Identity                    BodyTranslationPhaseFieldIdentityDTO `json:"identity"`
	FieldID                     int64                                `json:"fieldId"`
	FieldLabel                  string                               `json:"fieldLabel,omitempty"`
	SourceExcerpt               string                               `json:"sourceExcerpt,omitempty"`
	TranslatedText              string                               `json:"translatedText,omitempty"`
	OutputStatus                string                               `json:"outputStatus,omitempty"`
	ProtectionValidationResult  string                               `json:"protectionValidationResult,omitempty"`
	ProtectionValidationSummary string                               `json:"protectionValidationSummary,omitempty"`
	RetryCount                  int                                  `json:"retryCount"`
}

// BodyTranslationPhaseErrorSummaryDTO exposes only redacted error information.
type BodyTranslationPhaseErrorSummaryDTO struct {
	ErrorKind  string `json:"errorKind"`
	Reason     string `json:"reason"`
	Retryable  bool   `json:"retryable"`
	IsRedacted bool   `json:"isRedacted"`
}

// BodyTranslationPhaseActionEnablementDTO summarizes Job Run button state.
type BodyTranslationPhaseActionEnablementDTO struct {
	CanStart            bool    `json:"canStart"`
	StartBlockedReason  *string `json:"startBlockedReason,omitempty"`
	CanPause            bool    `json:"canPause"`
	PauseBlockedReason  *string `json:"pauseBlockedReason,omitempty"`
	CanResume           bool    `json:"canResume"`
	ResumeBlockedReason *string `json:"resumeBlockedReason,omitempty"`
	CanRetry            bool    `json:"canRetry"`
	RetryBlockedReason  *string `json:"retryBlockedReason,omitempty"`
	CanCancel           bool    `json:"canCancel"`
	CancelBlockedReason *string `json:"cancelBlockedReason,omitempty"`
}

// BodyTranslationOutputReadinessSummaryDTO summarizes downstream output readiness.
type BodyTranslationOutputReadinessSummaryDTO struct {
	CompletedFieldCount int  `json:"completedFieldCount"`
	StatusConsistent    bool `json:"statusConsistent"`
	OutputCount         int  `json:"outputCount"`
}

// BodyTranslationPhaseSummaryResponseDTO returns the frozen summary response shape.
type BodyTranslationPhaseSummaryResponseDTO struct {
	JobID            int64                                      `json:"jobId"`
	CurrentPhase     string                                     `json:"currentPhase"`
	PhaseState       string                                     `json:"phaseState"`
	PhaseRunID       *int64                                     `json:"phaseRunId,omitempty"`
	StartedAt        *string                                    `json:"startedAt,omitempty"`
	FinishedAt       *string                                    `json:"finishedAt,omitempty"`
	Progress         BodyTranslationPhaseProgressSummaryDTO     `json:"progress"`
	InputSummary     BodyTranslationPhaseInputSummaryDTO        `json:"inputSummary"`
	RequestSummary   BodyTranslationPhaseRequestSummaryDTO      `json:"requestSummary"`
	AISettings       *BodyTranslationPhaseAISettingsDTO         `json:"aiSettings,omitempty"`
	Execution        *BodyTranslationPhaseExecutionSummaryDTO   `json:"execution,omitempty"`
	FieldResults     []BodyTranslationPhaseFieldResultItemDTO   `json:"fieldResults,omitempty"`
	ResultSummary    *BodyTranslationPhaseFieldResultSummaryDTO `json:"resultSummary,omitempty"`
	ErrorSummary     *BodyTranslationPhaseErrorSummaryDTO       `json:"errorSummary,omitempty"`
	ActionEnablement BodyTranslationPhaseActionEnablementDTO    `json:"actionEnablement"`
	OutputReadiness  BodyTranslationOutputReadinessSummaryDTO   `json:"outputReadiness"`
}

// BodyTranslationPhaseCommandResponseDTO returns the frozen write-seam response shape.
type BodyTranslationPhaseCommandResponseDTO struct {
	JobID               int64                                      `json:"jobId"`
	CurrentPhase        string                                     `json:"currentPhase"`
	PhaseState          string                                     `json:"phaseState"`
	PhaseRunID          *int64                                     `json:"phaseRunId,omitempty"`
	StartedAt           *string                                    `json:"startedAt,omitempty"`
	FinishedAt          *string                                    `json:"finishedAt,omitempty"`
	Progress            BodyTranslationPhaseProgressSummaryDTO     `json:"progress"`
	InputSnapshotDigest string                                     `json:"inputSnapshotDigest,omitempty"`
	InputSummary        BodyTranslationPhaseInputSummaryDTO        `json:"inputSummary"`
	RequestSummary      BodyTranslationPhaseRequestSummaryDTO      `json:"requestSummary"`
	Execution           *BodyTranslationPhaseExecutionSummaryDTO   `json:"execution,omitempty"`
	FieldResults        []BodyTranslationPhaseFieldResultItemDTO   `json:"fieldResults,omitempty"`
	ResultSummary       *BodyTranslationPhaseFieldResultSummaryDTO `json:"resultSummary,omitempty"`
	Retryable           bool                                       `json:"retryable"`
	OutputReadiness     BodyTranslationOutputReadinessSummaryDTO   `json:"outputReadiness"`
	ErrorSummary        *BodyTranslationPhaseErrorSummaryDTO       `json:"errorSummary,omitempty"`
}

// NewBodyTranslationPhaseController creates a body translation phase controller.
func NewBodyTranslationPhaseController(usecase BodyTranslationPhaseUsecasePort) *BodyTranslationPhaseController {
	return &BodyTranslationPhaseController{bodyTranslationPhaseUsecase: usecase}
}

// GetBodyTranslationPhaseSummary returns the frozen Job Run summary contract.
func (controller *BodyTranslationPhaseController) GetBodyTranslationPhaseSummary(
	request GetBodyTranslationPhaseSummaryRequestDTO,
) (BodyTranslationPhaseSummaryResponseDTO, error) {
	result, err := controller.bodyTranslationPhaseUsecase.GetBodyTranslationPhaseSummary(
		context.Background(),
		usecase.GetBodyTranslationPhaseSummaryRequest{JobID: request.JobID},
	)
	if err != nil {
		return BodyTranslationPhaseSummaryResponseDTO{}, fmt.Errorf("get body translation phase summary: %w", err)
	}
	return toBodyTranslationPhaseSummaryResponseDTO(result), nil
}

// StartBodyTranslationPhase starts one body translation phase run or returns a rejected result.
func (controller *BodyTranslationPhaseController) StartBodyTranslationPhase(
	request StartBodyTranslationPhaseRequestDTO,
) (BodyTranslationPhaseCommandResponseDTO, error) {
	result, err := controller.bodyTranslationPhaseUsecase.StartBodyTranslationPhase(
		context.Background(),
		usecase.StartBodyTranslationPhaseRequest{JobID: request.JobID},
	)
	return bodyTranslationCommandResponseOrError(result, err, "start body translation phase")
}

// PauseBodyTranslationPhase pauses one active body translation phase run.
func (controller *BodyTranslationPhaseController) PauseBodyTranslationPhase(
	request PauseBodyTranslationPhaseRequestDTO,
) (BodyTranslationPhaseCommandResponseDTO, error) {
	result, err := controller.bodyTranslationPhaseUsecase.PauseBodyTranslationPhase(
		context.Background(),
		usecase.PauseBodyTranslationPhaseRequest{JobID: request.JobID, PhaseRunID: request.PhaseRunID},
	)
	return bodyTranslationCommandResponseOrError(result, err, "pause body translation phase")
}

// ResumeBodyTranslationPhase resumes one paused or recoverable failed body translation phase run.
func (controller *BodyTranslationPhaseController) ResumeBodyTranslationPhase(
	request ResumeBodyTranslationPhaseRequestDTO,
) (BodyTranslationPhaseCommandResponseDTO, error) {
	result, err := controller.bodyTranslationPhaseUsecase.ResumeBodyTranslationPhase(
		context.Background(),
		usecase.ResumeBodyTranslationPhaseRequest{JobID: request.JobID, PhaseRunID: request.PhaseRunID},
	)
	return bodyTranslationCommandResponseOrError(result, err, "resume body translation phase")
}

// RetryBodyTranslationPhase retries one retryable failed body translation phase run.
func (controller *BodyTranslationPhaseController) RetryBodyTranslationPhase(
	request RetryBodyTranslationPhaseRequestDTO,
) (BodyTranslationPhaseCommandResponseDTO, error) {
	result, err := controller.bodyTranslationPhaseUsecase.RetryBodyTranslationPhase(
		context.Background(),
		usecase.RetryBodyTranslationPhaseRequest{JobID: request.JobID, PhaseRunID: request.PhaseRunID},
	)
	return bodyTranslationCommandResponseOrError(result, err, "retry body translation phase")
}

// CancelBodyTranslationPhase cancels one paused body translation phase run.
func (controller *BodyTranslationPhaseController) CancelBodyTranslationPhase(
	request CancelBodyTranslationPhaseRequestDTO,
) (BodyTranslationPhaseCommandResponseDTO, error) {
	result, err := controller.bodyTranslationPhaseUsecase.CancelBodyTranslationPhase(
		context.Background(),
		usecase.CancelBodyTranslationPhaseRequest{JobID: request.JobID, PhaseRunID: request.PhaseRunID},
	)
	return bodyTranslationCommandResponseOrError(result, err, "cancel body translation phase")
}

// SaveBodyTranslationPhaseAISettings saves public AI settings for the body phase.
func (controller *BodyTranslationPhaseController) SaveBodyTranslationPhaseAISettings(
	request SaveBodyTranslationPhaseAISettingsRequestDTO,
) (BodyTranslationPhaseAISettingsResponseDTO, error) {
	settingsUsecase, ok := controller.bodyTranslationPhaseUsecase.(bodyTranslationPhaseAISettingsUsecasePort)
	if !ok {
		return BodyTranslationPhaseAISettingsResponseDTO{}, fmt.Errorf("save body translation phase ai settings: usecase is not configured")
	}
	result, err := settingsUsecase.SaveBodyTranslationPhaseAISettings(
		context.Background(),
		usecase.SaveBodyTranslationPhaseAISettingsRequest{
			Provider:      request.Provider,
			Model:         request.Model,
			ExecutionMode: request.ExecutionMode,
			BatchMode:     request.BatchMode,
		},
	)
	if err != nil {
		return BodyTranslationPhaseAISettingsResponseDTO{}, fmt.Errorf("save body translation phase ai settings: %w", err)
	}
	return BodyTranslationPhaseAISettingsResponseDTO{
		PhaseType:     result.PhaseType,
		Provider:      result.Provider,
		Model:         result.Model,
		ExecutionMode: result.ExecutionMode,
		BatchMode:     result.BatchMode,
	}, nil
}

func toBodyTranslationPhaseSummaryResponseDTO(
	result usecase.BodyTranslationPhaseSummaryResult,
) BodyTranslationPhaseSummaryResponseDTO {
	executionDTO := toOptionalBodyTranslationPhaseExecutionSummaryDTO(result.Execution)
	var outputCount int
	if result.Execution != nil {
		outputCount = result.Execution.OutputCount
	}
	return BodyTranslationPhaseSummaryResponseDTO{
		JobID:            result.JobID,
		CurrentPhase:     result.CurrentPhase,
		PhaseState:       result.PhaseState,
		PhaseRunID:       cloneOptionalInt64(result.PhaseRunID),
		StartedAt:        formatOptionalTime(result.StartedAt),
		FinishedAt:       formatOptionalTime(result.FinishedAt),
		Progress:         toBodyTranslationPhaseProgressSummaryDTO(result.Progress),
		InputSummary:     toBodyTranslationPhaseInputSummaryDTO(result.InputSummary),
		RequestSummary:   toBodyTranslationPhaseRequestSummaryDTO(result.RequestSummary),
		AISettings:       toBodyTranslationPhaseAISettingsDTO(result.AISettings),
		Execution:        executionDTO,
		FieldResults:     toBodyTranslationPhaseFieldResultItemsDTO(result.FieldResults),
		ResultSummary:    toOptionalBodyTranslationPhaseFieldResultSummaryDTO(firstNonNilBodyTranslationFieldResult(result.ResultSummary, result.FieldResultSummary)),
		ErrorSummary:     toOptionalBodyTranslationPhaseErrorSummaryDTO(result.ErrorSummary),
		ActionEnablement: toBodyTranslationPhaseActionEnablementDTO(result.ActionEnablement),
		OutputReadiness:  toBodyTranslationOutputReadinessSummaryDTO(result.OutputReadiness, outputCount),
	}
}

func bodyTranslationCommandResponseOrError(
	result usecase.BodyTranslationPhaseCommandResult,
	err error,
	operation string,
) (BodyTranslationPhaseCommandResponseDTO, error) {
	response := toBodyTranslationPhaseCommandResponseDTO(result)
	if err == nil {
		return response, nil
	}
	if result.ErrorSummary != nil {
		return response, nil
	}
	return response, fmt.Errorf("%s: %w", operation, err)
}

func toBodyTranslationPhaseCommandResponseDTO(
	result usecase.BodyTranslationPhaseCommandResult,
) BodyTranslationPhaseCommandResponseDTO {
	executionDTO := toOptionalBodyTranslationPhaseExecutionSummaryDTO(result.Execution)
	var outputCount int
	if result.Execution != nil {
		outputCount = result.Execution.OutputCount
	}
	return BodyTranslationPhaseCommandResponseDTO{
		JobID:               result.JobID,
		CurrentPhase:        result.CurrentPhase,
		PhaseState:          result.PhaseState,
		PhaseRunID:          cloneOptionalInt64(result.PhaseRunID),
		StartedAt:           formatOptionalTime(result.StartedAt),
		FinishedAt:          formatOptionalTime(result.FinishedAt),
		Progress:            toBodyTranslationPhaseProgressSummaryDTO(result.Progress),
		InputSnapshotDigest: result.InputSnapshotDigest,
		InputSummary:        toBodyTranslationPhaseInputSummaryDTO(result.InputSummary),
		RequestSummary:      toBodyTranslationPhaseRequestSummaryDTO(result.RequestSummary),
		Execution:           executionDTO,
		FieldResults:        toBodyTranslationPhaseFieldResultItemsDTO(result.FieldResults),
		ResultSummary:       toOptionalBodyTranslationPhaseFieldResultSummaryDTO(firstNonNilBodyTranslationFieldResult(result.ResultSummary, result.FieldResultSummary)),
		Retryable:           result.Retryable,
		OutputReadiness:     toBodyTranslationOutputReadinessSummaryDTO(result.OutputReadiness, outputCount),
		ErrorSummary:        toOptionalBodyTranslationPhaseErrorSummaryDTO(result.ErrorSummary),
	}
}

func toBodyTranslationPhaseProgressSummaryDTO(
	summary usecase.BodyTranslationPhaseProgressSummary,
) BodyTranslationPhaseProgressSummaryDTO {
	return BodyTranslationPhaseProgressSummaryDTO{
		Percent:         summary.Percent,
		ProcessedCount:  summary.ProcessedCount,
		TotalCount:      summary.TotalCount,
		TargetCount:     summary.TargetCount,
		TranslatedCount: summary.TranslatedCount,
		SkippedCount:    summary.SkippedCount,
		CurrentStep:     summary.CurrentStep,
	}
}

func toBodyTranslationPhaseInputSummaryDTO(
	summary usecase.BodyTranslationPhaseInputSummary,
) BodyTranslationPhaseInputSummaryDTO {
	return BodyTranslationPhaseInputSummaryDTO{
		TargetCount:      summary.TargetCount,
		SkippedReasons:   append([]string{}, summary.SkippedReasons...),
		InputSnapshotRef: cloneOptionalString(summary.InputSnapshotRef),
		DictionaryDigest: summary.DictionaryDigest,
		PersonaDigest:    summary.PersonaDigest,
		MetadataDigest:   summary.MetadataDigest,
		PromptDigest:     summary.PromptDigest,
	}
}

func toBodyTranslationPhaseRequestSummaryDTO(
	summary usecase.BodyTranslationPhaseRequestSummary,
) BodyTranslationPhaseRequestSummaryDTO {
	return BodyTranslationPhaseRequestSummaryDTO{
		ProviderTargetCount:              summary.ProviderTargetCount,
		ExactDictionaryExclusionCount:    summary.ExactDictionaryExclusionCount,
		PartialDictionaryConstraintCount: summary.PartialDictionaryConstraintCount,
	}
}

func toBodyTranslationPhaseExecutionSummaryDTO(
	summary usecase.BodyTranslationPhaseExecutionSummary,
) BodyTranslationPhaseExecutionSummaryDTO {
	return BodyTranslationPhaseExecutionSummaryDTO{
		CredentialRef:    summary.CredentialRef,
		Provider:         summary.Provider,
		Model:            summary.Model,
		ExecutionMode:    summary.ExecutionMode,
		RequestUnitCount: summary.RequestUnitCount,
		OutputCount:      summary.OutputCount,
	}
}

func toOptionalBodyTranslationPhaseExecutionSummaryDTO(
	summary *usecase.BodyTranslationPhaseExecutionSummary,
) *BodyTranslationPhaseExecutionSummaryDTO {
	if summary == nil {
		return nil
	}
	v := toBodyTranslationPhaseExecutionSummaryDTO(*summary)
	return &v
}

func toBodyTranslationPhaseAISettingsDTO(
	settings *usecase.BodyTranslationPhaseAISettings,
) *BodyTranslationPhaseAISettingsDTO {
	if settings == nil {
		return nil
	}
	return &BodyTranslationPhaseAISettingsDTO{
		Provider:      settings.Provider,
		Model:         settings.Model,
		ExecutionMode: settings.ExecutionMode,
		BatchMode:     settings.BatchMode,
	}
}

func toOptionalBodyTranslationPhaseFieldResultSummaryDTO(
	summary *usecase.BodyTranslationPhaseFieldResultSummary,
) *BodyTranslationPhaseFieldResultSummaryDTO {
	if summary == nil {
		return nil
	}
	return &BodyTranslationPhaseFieldResultSummaryDTO{
		TranslatedCount:       summary.TranslatedCount,
		FailedCount:           summary.FailedCount,
		SkippedCount:          summary.SkippedCount,
		ProtectionFailedCount: summary.ProtectionFailedCount,
		OutputReadyCount:      firstNonZero(summary.OutputReadyCount, summary.OutputCount),
		OutputCount:           summary.OutputCount,
		FieldResults:          toBodyTranslationPhaseFieldResultItemsDTO(summary.FieldResults),
	}
}

func toBodyTranslationPhaseFieldResultItemsDTO(
	items []usecase.BodyTranslationPhaseFieldResultItem,
) []BodyTranslationPhaseFieldResultItemDTO {
	if len(items) == 0 {
		return nil
	}
	result := make([]BodyTranslationPhaseFieldResultItemDTO, 0, len(items))
	for _, item := range items {
		result = append(result, BodyTranslationPhaseFieldResultItemDTO{
			Identity: BodyTranslationPhaseFieldIdentityDTO{
				TranslationFieldID:      item.Identity.TranslationFieldID,
				PhaseTranslationFieldID: item.Identity.PhaseTranslationFieldID,
				RecordType:              item.Identity.RecordType,
				FieldType:               item.Identity.FieldType,
				FormID:                  item.Identity.FormID,
				EditorID:                item.Identity.EditorID,
				FieldLabel:              item.Identity.FieldLabel,
			},
			FieldID:                     item.FieldID,
			FieldLabel:                  item.FieldLabel,
			SourceExcerpt:               item.SourceExcerpt,
			TranslatedText:              item.TranslatedText,
			OutputStatus:                item.OutputStatus,
			ProtectionValidationResult:  item.ProtectionValidationResult,
			ProtectionValidationSummary: item.ProtectionValidationSummary,
			RetryCount:                  item.RetryCount,
		})
	}
	return result
}

func toOptionalBodyTranslationPhaseErrorSummaryDTO(
	summary *usecase.BodyTranslationPhaseErrorSummary,
) *BodyTranslationPhaseErrorSummaryDTO {
	if summary == nil {
		return nil
	}
	return &BodyTranslationPhaseErrorSummaryDTO{
		ErrorKind:  string(usecase.NormalizeBodyTranslationPhasePublicErrorKind(summary.ErrorKind)),
		Reason:     summary.Reason,
		Retryable:  summary.Retryable,
		IsRedacted: summary.IsRedacted,
	}
}

func toBodyTranslationPhaseActionEnablementDTO(
	summary usecase.BodyTranslationPhaseActionEnablement,
) BodyTranslationPhaseActionEnablementDTO {
	return BodyTranslationPhaseActionEnablementDTO{
		CanStart:            summary.CanStart,
		StartBlockedReason:  cloneOptionalString(summary.StartBlockedReason),
		CanPause:            summary.CanPause,
		PauseBlockedReason:  cloneOptionalString(summary.PauseBlockedReason),
		CanResume:           summary.CanResume,
		ResumeBlockedReason: cloneOptionalString(summary.ResumeBlockedReason),
		CanRetry:            summary.CanRetry,
		RetryBlockedReason:  cloneOptionalString(summary.RetryBlockedReason),
		CanCancel:           summary.CanCancel,
		CancelBlockedReason: cloneOptionalString(summary.CancelBlockedReason),
	}
}

func toBodyTranslationOutputReadinessSummaryDTO(
	summary usecase.BodyTranslationOutputReadinessSummary,
	outputCount int,
) BodyTranslationOutputReadinessSummaryDTO {
	return BodyTranslationOutputReadinessSummaryDTO{
		CompletedFieldCount: summary.CompletedFieldCount,
		StatusConsistent:    summary.StatusConsistent,
		OutputCount:         outputCount,
	}
}

func firstNonNilBodyTranslationFieldResult(
	primary *usecase.BodyTranslationPhaseFieldResultSummary,
	fallback *usecase.BodyTranslationFieldResultSummary,
) *usecase.BodyTranslationPhaseFieldResultSummary {
	if primary != nil {
		return primary
	}
	return fallback
}

func firstNonZero(primary, fallback int) int {
	if primary != 0 {
		return primary
	}
	return fallback
}
