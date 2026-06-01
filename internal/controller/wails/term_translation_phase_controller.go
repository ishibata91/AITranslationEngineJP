package wails

import (
	"context"
	"fmt"
	"time"

	"aitranslationenginejp/internal/usecase"
)

// TermTranslationPhaseUsecasePort defines the frozen term translation phase usecase seam.
type TermTranslationPhaseUsecasePort interface {
	GetTermTranslationPhaseSummary(ctx context.Context, request usecase.GetTermTranslationPhaseSummaryRequest) (usecase.TermTranslationPhaseSummaryResult, error)
	StartTermTranslationPhase(ctx context.Context, request usecase.StartTermTranslationPhaseRequest) (usecase.TermTranslationPhaseCommandResult, error)
	PauseTermTranslationPhase(ctx context.Context, request usecase.PauseTermTranslationPhaseRequest) (usecase.TermTranslationPhaseCommandResult, error)
	ResumeTermTranslationPhase(ctx context.Context, request usecase.ResumeTermTranslationPhaseRequest) (usecase.TermTranslationPhaseCommandResult, error)
	RetryTermTranslationPhase(ctx context.Context, request usecase.RetryTermTranslationPhaseRequest) (usecase.TermTranslationPhaseCommandResult, error)
	GetTermTranslationNextPhaseReadiness(ctx context.Context, request usecase.GetTermTranslationNextPhaseReadinessRequest) (usecase.TermTranslationNextPhaseReadinessResult, error)
}

type termTranslationPhaseAISettingsUsecasePort interface {
	SaveTermTranslationPhaseAISettings(ctx context.Context, request usecase.SaveTermTranslationPhaseAISettingsRequest) (usecase.TermTranslationPhaseAISettingsResult, error)
}

// TermTranslationPhaseController exposes Wails-bound term translation entrypoints.
type TermTranslationPhaseController struct {
	termTranslationPhaseUsecase TermTranslationPhaseUsecasePort
}

// GetTermTranslationPhaseSummaryRequestDTO identifies the requested translation job.
type GetTermTranslationPhaseSummaryRequestDTO struct {
	JobID int64 `json:"jobId"`
}

// StartTermTranslationPhaseRequestDTO identifies the ready job that wants to create a phase run.
type StartTermTranslationPhaseRequestDTO struct {
	JobID int64 `json:"jobId"`
}

// PauseTermTranslationPhaseRequestDTO identifies one active phase run to pause.
type PauseTermTranslationPhaseRequestDTO struct {
	JobID      int64 `json:"jobId"`
	PhaseRunID int64 `json:"phaseRunId"`
}

// ResumeTermTranslationPhaseRequestDTO identifies one paused or recoverable failed phase run to resume.
type ResumeTermTranslationPhaseRequestDTO struct {
	JobID      int64 `json:"jobId"`
	PhaseRunID int64 `json:"phaseRunId"`
}

// RetryTermTranslationPhaseRequestDTO identifies one retryable failed phase run to retry.
type RetryTermTranslationPhaseRequestDTO struct {
	JobID      int64 `json:"jobId"`
	PhaseRunID int64 `json:"phaseRunId"`
}

// GetTermTranslationNextPhaseReadinessRequestDTO identifies the downstream readiness request target.
type GetTermTranslationNextPhaseReadinessRequestDTO struct {
	JobID int64 `json:"jobId"`
}

// SaveTermTranslationPhaseAISettingsRequestDTO carries public AI settings.
// job_id を含めない。保存操作はジョブと無関連で phase_type = "word_translation" へ upsert する。
type SaveTermTranslationPhaseAISettingsRequestDTO struct {
	Provider      string `json:"provider"`
	Model         string `json:"model"`
	ExecutionMode string `json:"executionMode"`
	BatchMode     string `json:"batchMode"`
}

// TermTranslationPhaseAISettingsResponseDTO returns saved AI settings state.
// credential_status、model_list_status は含まない。
type TermTranslationPhaseAISettingsResponseDTO struct {
	PhaseType     string `json:"phaseType"`
	Provider      string `json:"provider"`
	Model         string `json:"model"`
	ExecutionMode string `json:"executionMode"`
	BatchMode     string `json:"batchMode"`
}

// TermTranslationPhaseAISettingsDTO holds the Ready 期 AI 選択値。
// JOB_PHASE_AI_SETTINGS record が存在する場合だけ応答に含める。record 不在は SummaryResponseDTO の aiSettings field を omit して表現する。
type TermTranslationPhaseAISettingsDTO struct {
	Provider      string `json:"provider"`
	Model         string `json:"model"`
	ExecutionMode string `json:"executionMode"`
	BatchMode     string `json:"batchMode"`
}

// TermTranslationExecutionConfigSummaryDTO summarizes execution inputs without exposing secrets.
// JOB_PHASE_RUN が存在する場合だけ応答に含める。未開始（run 未作成）は SummaryResponseDTO の execution field を omit して表現する。
type TermTranslationExecutionConfigSummaryDTO struct {
	CredentialRef   string  `json:"credentialRef"`
	Provider        string  `json:"provider"`
	Model           string  `json:"model"`
	ExecutionMode   string  `json:"executionMode"`
	SnapshotDigest  *string `json:"snapshotDigest,omitempty"`
	SnapshotVersion *string `json:"snapshotVersion,omitempty"`
}

// TermTranslationPhaseProgressSummaryDTO summarizes one phase run progress snapshot.
type TermTranslationPhaseProgressSummaryDTO struct {
	Percent        int    `json:"percent"`
	ProcessedCount int    `json:"processedCount"`
	TotalCount     int    `json:"totalCount"`
	AITargetCount  int    `json:"aiTargetCount"`
	CurrentStep    string `json:"currentStep"`
}

// TermTranslationPhaseResultSummaryDTO summarizes one phase result snapshot.
type TermTranslationPhaseResultSummaryDTO struct {
	ConfirmedCount            int  `json:"confirmedCount"`
	JobDictionaryAppliedCount int  `json:"jobDictionaryAppliedCount"`
	ReplacementTargetCount    int  `json:"replacementTargetCount"`
	UnmatchedCount            int  `json:"unmatchedCount"`
	ProviderSkipped           bool `json:"providerSkipped"`
}

// TermTranslationPhaseErrorSummaryDTO exposes only redacted error information.
type TermTranslationPhaseErrorSummaryDTO struct {
	ErrorKind  string `json:"errorKind"`
	Reason     string `json:"reason"`
	Retryable  bool   `json:"retryable"`
	IsRedacted bool   `json:"isRedacted"`
}

// TermTranslationPhaseActionEnablementDTO summarizes Job Run operation button state.
type TermTranslationPhaseActionEnablementDTO struct {
	CanStart            bool    `json:"canStart"`
	StartBlockedReason  *string `json:"startBlockedReason,omitempty"`
	CanPause            bool    `json:"canPause"`
	PauseBlockedReason  *string `json:"pauseBlockedReason,omitempty"`
	CanResume           bool    `json:"canResume"`
	ResumeBlockedReason *string `json:"resumeBlockedReason,omitempty"`
	CanRetry            bool    `json:"canRetry"`
	RetryBlockedReason  *string `json:"retryBlockedReason,omitempty"`
}

// TermTranslationPhaseSummaryResponseDTO returns the frozen Job Run summary shape.
type TermTranslationPhaseSummaryResponseDTO struct {
	JobID              int64                                  `json:"jobId"`
	CurrentPhase       string                                 `json:"currentPhase"`
	PhaseState         string                                 `json:"phaseState"`
	PhaseRunID         *int64                                 `json:"phaseRunId,omitempty"`
	StartedAt          *string                                `json:"startedAt,omitempty"`
	FinishedAt         *string                                `json:"finishedAt,omitempty"`
	Progress           TermTranslationPhaseProgressSummaryDTO `json:"progress"`
	TotalTermCount     int                                    `json:"totalTermCount"`
	DictionaryHitCount int                                    `json:"dictionaryHitCount"`
	AITargetCount      int                                    `json:"aiTargetCount"`
	// AISettings は JOB_PHASE_AI_SETTINGS record が存在する場合だけ含める。nil は omit される。
	AISettings *TermTranslationPhaseAISettingsDTO `json:"aiSettings,omitempty"`
	// Execution は JOB_PHASE_RUN が存在する場合だけ含める。nil は omit される。
	Execution        *TermTranslationExecutionConfigSummaryDTO `json:"execution,omitempty"`
	ResultSummary    *TermTranslationPhaseResultSummaryDTO     `json:"resultSummary,omitempty"`
	ErrorSummary     *TermTranslationPhaseErrorSummaryDTO      `json:"errorSummary,omitempty"`
	ActionEnablement TermTranslationPhaseActionEnablementDTO   `json:"actionEnablement"`
}

// TermTranslationPhaseCommandResponseDTO returns the frozen write-seam response shape.
type TermTranslationPhaseCommandResponseDTO struct {
	JobID        int64                                  `json:"jobId"`
	CurrentPhase string                                 `json:"currentPhase"`
	PhaseState   string                                 `json:"phaseState"`
	PhaseRunID   *int64                                 `json:"phaseRunId,omitempty"`
	StartedAt    *string                                `json:"startedAt,omitempty"`
	FinishedAt   *string                                `json:"finishedAt,omitempty"`
	Progress     TermTranslationPhaseProgressSummaryDTO `json:"progress"`
	Retryable    bool                                   `json:"retryable"`
	ErrorSummary *TermTranslationPhaseErrorSummaryDTO   `json:"errorSummary,omitempty"`
}

// TermTranslationNextPhaseReadinessResponseDTO returns the frozen downstream phase fact state shape.
type TermTranslationNextPhaseReadinessResponseDTO struct {
	JobID          int64  `json:"jobId"`
	CurrentPhase   string `json:"currentPhase"`
	PhaseState     string `json:"phaseState"`
	JobIsTerminal  bool   `json:"jobIsTerminal"`
	TotalCount     int    `json:"totalCount"`
	ConfirmedCount int    `json:"confirmedCount"`
	ErrorKind      string `json:"errorKind,omitempty"`
}

// NewTermTranslationPhaseController creates a term translation phase controller.
func NewTermTranslationPhaseController(usecase TermTranslationPhaseUsecasePort) *TermTranslationPhaseController {
	return &TermTranslationPhaseController{termTranslationPhaseUsecase: usecase}
}

// GetTermTranslationPhaseSummary returns the frozen Job Run summary contract.
func (controller *TermTranslationPhaseController) GetTermTranslationPhaseSummary(
	request GetTermTranslationPhaseSummaryRequestDTO,
) (TermTranslationPhaseSummaryResponseDTO, error) {
	result, err := controller.termTranslationPhaseUsecase.GetTermTranslationPhaseSummary(
		context.Background(),
		usecase.GetTermTranslationPhaseSummaryRequest{JobID: request.JobID},
	)
	if err != nil {
		return TermTranslationPhaseSummaryResponseDTO{}, fmt.Errorf("get term translation phase summary: %w", err)
	}
	return toTermTranslationPhaseSummaryResponseDTO(result), nil
}

// StartTermTranslationPhase starts one term translation phase run or returns a rejected result.
func (controller *TermTranslationPhaseController) StartTermTranslationPhase(
	request StartTermTranslationPhaseRequestDTO,
) (TermTranslationPhaseCommandResponseDTO, error) {
	result, err := controller.termTranslationPhaseUsecase.StartTermTranslationPhase(
		context.Background(),
		usecase.StartTermTranslationPhaseRequest{JobID: request.JobID},
	)
	if err != nil {
		return TermTranslationPhaseCommandResponseDTO{}, fmt.Errorf("start term translation phase: %w", err)
	}
	return toTermTranslationPhaseCommandResponseDTO(result), nil
}

// PauseTermTranslationPhase pauses one active term translation phase run.
func (controller *TermTranslationPhaseController) PauseTermTranslationPhase(
	request PauseTermTranslationPhaseRequestDTO,
) (TermTranslationPhaseCommandResponseDTO, error) {
	result, err := controller.termTranslationPhaseUsecase.PauseTermTranslationPhase(
		context.Background(),
		usecase.PauseTermTranslationPhaseRequest{
			JobID:      request.JobID,
			PhaseRunID: request.PhaseRunID,
		},
	)
	if err != nil {
		return TermTranslationPhaseCommandResponseDTO{}, fmt.Errorf("pause term translation phase: %w", err)
	}
	return toTermTranslationPhaseCommandResponseDTO(result), nil
}

// ResumeTermTranslationPhase resumes one paused or recoverable failed term translation phase run.
func (controller *TermTranslationPhaseController) ResumeTermTranslationPhase(
	request ResumeTermTranslationPhaseRequestDTO,
) (TermTranslationPhaseCommandResponseDTO, error) {
	result, err := controller.termTranslationPhaseUsecase.ResumeTermTranslationPhase(
		context.Background(),
		usecase.ResumeTermTranslationPhaseRequest{
			JobID:      request.JobID,
			PhaseRunID: request.PhaseRunID,
		},
	)
	if err != nil {
		return TermTranslationPhaseCommandResponseDTO{}, fmt.Errorf("resume term translation phase: %w", err)
	}
	return toTermTranslationPhaseCommandResponseDTO(result), nil
}

// RetryTermTranslationPhase retries one retryable failed term translation phase run.
func (controller *TermTranslationPhaseController) RetryTermTranslationPhase(
	request RetryTermTranslationPhaseRequestDTO,
) (TermTranslationPhaseCommandResponseDTO, error) {
	result, err := controller.termTranslationPhaseUsecase.RetryTermTranslationPhase(
		context.Background(),
		usecase.RetryTermTranslationPhaseRequest{
			JobID:      request.JobID,
			PhaseRunID: request.PhaseRunID,
		},
	)
	if err != nil {
		return TermTranslationPhaseCommandResponseDTO{}, fmt.Errorf("retry term translation phase: %w", err)
	}
	return toTermTranslationPhaseCommandResponseDTO(result), nil
}

// GetTermTranslationNextPhaseReadiness returns whether the downstream phase can start.
func (controller *TermTranslationPhaseController) GetTermTranslationNextPhaseReadiness(
	request GetTermTranslationNextPhaseReadinessRequestDTO,
) (TermTranslationNextPhaseReadinessResponseDTO, error) {
	result, err := controller.termTranslationPhaseUsecase.GetTermTranslationNextPhaseReadiness(
		context.Background(),
		usecase.GetTermTranslationNextPhaseReadinessRequest{JobID: request.JobID},
	)
	if err != nil {
		return TermTranslationNextPhaseReadinessResponseDTO{}, fmt.Errorf("get term translation next phase readiness: %w", err)
	}
	return toTermTranslationNextPhaseReadinessResponseDTO(result), nil
}

// SaveTermTranslationPhaseAISettings saves public AI settings for the term phase.
func (controller *TermTranslationPhaseController) SaveTermTranslationPhaseAISettings(
	request SaveTermTranslationPhaseAISettingsRequestDTO,
) (TermTranslationPhaseAISettingsResponseDTO, error) {
	settingsUsecase, ok := controller.termTranslationPhaseUsecase.(termTranslationPhaseAISettingsUsecasePort)
	if !ok {
		return TermTranslationPhaseAISettingsResponseDTO{}, fmt.Errorf("save term translation phase ai settings: usecase is not configured")
	}
	result, err := settingsUsecase.SaveTermTranslationPhaseAISettings(
		context.Background(),
		usecase.SaveTermTranslationPhaseAISettingsRequest{
			Provider:      request.Provider,
			Model:         request.Model,
			ExecutionMode: request.ExecutionMode,
			BatchMode:     request.BatchMode,
		},
	)
	if err != nil {
		return TermTranslationPhaseAISettingsResponseDTO{}, fmt.Errorf("save term translation phase ai settings: %w", err)
	}
	return TermTranslationPhaseAISettingsResponseDTO{
		PhaseType:     result.PhaseType,
		Provider:      result.Provider,
		Model:         result.Model,
		ExecutionMode: result.ExecutionMode,
		BatchMode:     result.BatchMode,
	}, nil
}

func toTermTranslationPhaseSummaryResponseDTO(
	result usecase.TermTranslationPhaseSummaryResult,
) TermTranslationPhaseSummaryResponseDTO {
	return TermTranslationPhaseSummaryResponseDTO{
		JobID:              result.JobID,
		CurrentPhase:       result.CurrentPhase,
		PhaseState:         result.PhaseState,
		PhaseRunID:         cloneOptionalInt64(result.PhaseRunID),
		StartedAt:          formatOptionalTime(result.StartedAt),
		FinishedAt:         formatOptionalTime(result.FinishedAt),
		Progress:           toTermTranslationPhaseProgressSummaryDTO(result.Progress),
		TotalTermCount:     result.TotalTermCount,
		DictionaryHitCount: result.DictionaryHitCount,
		AITargetCount:      result.AITargetCount,
		AISettings:         toOptionalTermTranslationPhaseAISettingsDTO(result.AISettings),
		Execution:          toOptionalTermTranslationExecutionConfigSummaryDTO(result.Execution),
		ResultSummary:      toOptionalTermTranslationPhaseResultSummaryDTO(result.ResultSummary),
		ErrorSummary:       toOptionalTermTranslationPhaseErrorSummaryDTO(result.ErrorSummary),
		ActionEnablement:   toTermTranslationPhaseActionEnablementDTO(result.ActionEnablement),
	}
}

func toTermTranslationPhaseCommandResponseDTO(
	result usecase.TermTranslationPhaseCommandResult,
) TermTranslationPhaseCommandResponseDTO {
	return TermTranslationPhaseCommandResponseDTO{
		JobID:        result.JobID,
		CurrentPhase: result.CurrentPhase,
		PhaseState:   result.PhaseState,
		PhaseRunID:   cloneOptionalInt64(result.PhaseRunID),
		StartedAt:    formatOptionalTime(result.StartedAt),
		FinishedAt:   formatOptionalTime(result.FinishedAt),
		Progress:     toTermTranslationPhaseProgressSummaryDTO(result.Progress),
		Retryable:    result.Retryable,
		ErrorSummary: toOptionalTermTranslationPhaseErrorSummaryDTO(result.ErrorSummary),
	}
}

func toTermTranslationNextPhaseReadinessResponseDTO(
	result usecase.TermTranslationNextPhaseReadinessResult,
) TermTranslationNextPhaseReadinessResponseDTO {
	return TermTranslationNextPhaseReadinessResponseDTO{
		JobID:          result.JobID,
		CurrentPhase:   result.CurrentPhase,
		PhaseState:     result.PhaseState,
		JobIsTerminal:  result.JobIsTerminal,
		TotalCount:     result.TotalCount,
		ConfirmedCount: result.ConfirmedCount,
		ErrorKind:      usecase.NormalizeTermTranslationPhasePublicErrorKind(result.ErrorKind),
	}
}

func toOptionalTermTranslationPhaseAISettingsDTO(
	settings *usecase.TermTranslationPhaseAISettings,
) *TermTranslationPhaseAISettingsDTO {
	if settings == nil {
		return nil
	}
	return &TermTranslationPhaseAISettingsDTO{
		Provider:      settings.Provider,
		Model:         settings.Model,
		ExecutionMode: settings.ExecutionMode,
		BatchMode:     settings.BatchMode,
	}
}

func toOptionalTermTranslationExecutionConfigSummaryDTO(
	summary *usecase.TermTranslationExecutionConfigSummary,
) *TermTranslationExecutionConfigSummaryDTO {
	if summary == nil {
		return nil
	}
	return &TermTranslationExecutionConfigSummaryDTO{
		CredentialRef:   summary.CredentialRef,
		Provider:        summary.Provider,
		Model:           summary.Model,
		ExecutionMode:   summary.ExecutionMode,
		SnapshotDigest:  cloneOptionalString(summary.SnapshotDigest),
		SnapshotVersion: cloneOptionalString(summary.SnapshotVersion),
	}
}

func toTermTranslationPhaseProgressSummaryDTO(
	summary usecase.TermTranslationPhaseProgressSummary,
) TermTranslationPhaseProgressSummaryDTO {
	return TermTranslationPhaseProgressSummaryDTO{
		Percent:        summary.Percent,
		ProcessedCount: summary.ProcessedCount,
		TotalCount:     summary.TotalCount,
		AITargetCount:  summary.AITargetCount,
		CurrentStep:    summary.CurrentStep,
	}
}

func toOptionalTermTranslationPhaseResultSummaryDTO(
	summary *usecase.TermTranslationPhaseResultSummary,
) *TermTranslationPhaseResultSummaryDTO {
	if summary == nil {
		return nil
	}
	return &TermTranslationPhaseResultSummaryDTO{
		ConfirmedCount:            summary.ConfirmedCount,
		JobDictionaryAppliedCount: summary.JobDictionaryAppliedCount,
		ReplacementTargetCount:    summary.ReplacementTargetCount,
		UnmatchedCount:            summary.UnmatchedCount,
		ProviderSkipped:           summary.ProviderSkipped,
	}
}

func toOptionalTermTranslationPhaseErrorSummaryDTO(
	summary *usecase.TermTranslationPhaseErrorSummary,
) *TermTranslationPhaseErrorSummaryDTO {
	if summary == nil {
		return nil
	}
	return &TermTranslationPhaseErrorSummaryDTO{
		ErrorKind:  usecase.NormalizeTermTranslationPhasePublicErrorKind(summary.ErrorKind),
		Reason:     summary.Reason,
		Retryable:  summary.Retryable,
		IsRedacted: summary.IsRedacted,
	}
}

func toTermTranslationPhaseActionEnablementDTO(
	enablement usecase.TermTranslationPhaseActionEnablement,
) TermTranslationPhaseActionEnablementDTO {
	return TermTranslationPhaseActionEnablementDTO{
		CanStart:            enablement.CanStart,
		StartBlockedReason:  cloneOptionalString(enablement.StartBlockedReason),
		CanPause:            enablement.CanPause,
		PauseBlockedReason:  cloneOptionalString(enablement.PauseBlockedReason),
		CanResume:           enablement.CanResume,
		ResumeBlockedReason: cloneOptionalString(enablement.ResumeBlockedReason),
		CanRetry:            enablement.CanRetry,
		RetryBlockedReason:  cloneOptionalString(enablement.RetryBlockedReason),
	}
}

func formatOptionalTime(value *time.Time) *string {
	if value == nil {
		return nil
	}
	formatted := value.UTC().Format(time.RFC3339)
	return &formatted
}

func cloneOptionalInt64(value *int64) *int64 {
	if value == nil {
		return nil
	}
	cloned := *value
	return &cloned
}

func cloneOptionalString(value *string) *string {
	if value == nil {
		return nil
	}
	cloned := *value
	return &cloned
}
