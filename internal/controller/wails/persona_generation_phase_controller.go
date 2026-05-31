package wails

import (
	"context"
	"fmt"

	"aitranslationenginejp/internal/usecase"
)

// PersonaGenerationPhaseUsecasePort defines the frozen persona generation phase usecase seam.
type PersonaGenerationPhaseUsecasePort interface {
	GetPersonaGenerationPhaseSummary(ctx context.Context, request usecase.GetPersonaGenerationPhaseSummaryRequest) (usecase.PersonaGenerationPhaseSummaryResult, error)
	StartPersonaGenerationPhase(ctx context.Context, request usecase.StartPersonaGenerationPhaseRequest) (usecase.PersonaGenerationPhaseCommandResult, error)
	PausePersonaGenerationPhase(ctx context.Context, request usecase.PausePersonaGenerationPhaseRequest) (usecase.PersonaGenerationPhaseCommandResult, error)
	ResumePersonaGenerationPhase(ctx context.Context, request usecase.ResumePersonaGenerationPhaseRequest) (usecase.PersonaGenerationPhaseCommandResult, error)
	RetryPersonaGenerationPhase(ctx context.Context, request usecase.RetryPersonaGenerationPhaseRequest) (usecase.PersonaGenerationPhaseCommandResult, error)
	CancelPersonaGenerationPhase(ctx context.Context, request usecase.CancelPersonaGenerationPhaseRequest) (usecase.PersonaGenerationPhaseCommandResult, error)
	GetPersonaGenerationBodyReadiness(ctx context.Context, request usecase.GetPersonaGenerationBodyReadinessRequest) (usecase.PersonaGenerationBodyReadinessResult, error)
}

type personaGenerationPhaseAISettingsUsecasePort interface {
	SavePersonaGenerationPhaseAISettings(ctx context.Context, request usecase.SavePersonaGenerationPhaseAISettingsRequest) (usecase.PersonaGenerationPhaseAISettingsResult, error)
}

// PersonaGenerationPhaseController exposes Wails-bound persona generation entrypoints.
type PersonaGenerationPhaseController struct {
	personaGenerationPhaseUsecase PersonaGenerationPhaseUsecasePort
}

// GetPersonaGenerationPhaseSummaryRequestDTO identifies the requested job.
type GetPersonaGenerationPhaseSummaryRequestDTO struct {
	JobID int64 `json:"jobId"`
}

// StartPersonaGenerationPhaseRequestDTO identifies the ready job to start.
type StartPersonaGenerationPhaseRequestDTO struct {
	JobID int64 `json:"jobId"`
}

// PausePersonaGenerationPhaseRequestDTO identifies one active phase run to pause.
type PausePersonaGenerationPhaseRequestDTO struct {
	JobID      int64 `json:"jobId"`
	PhaseRunID int64 `json:"phaseRunId"`
}

// ResumePersonaGenerationPhaseRequestDTO identifies one paused or recoverable failed phase run to resume.
type ResumePersonaGenerationPhaseRequestDTO struct {
	JobID      int64 `json:"jobId"`
	PhaseRunID int64 `json:"phaseRunId"`
}

// RetryPersonaGenerationPhaseRequestDTO identifies one retryable failed phase run to retry.
type RetryPersonaGenerationPhaseRequestDTO struct {
	JobID      int64 `json:"jobId"`
	PhaseRunID int64 `json:"phaseRunId"`
}

// CancelPersonaGenerationPhaseRequestDTO identifies one cancelable phase run.
type CancelPersonaGenerationPhaseRequestDTO struct {
	JobID      int64 `json:"jobId"`
	PhaseRunID int64 `json:"phaseRunId"`
}

// GetPersonaGenerationBodyReadinessRequestDTO identifies the job whose body readiness is requested.
type GetPersonaGenerationBodyReadinessRequestDTO struct {
	JobID int64 `json:"jobId"`
}

// SavePersonaGenerationPhaseAISettingsRequestDTO carries public AI settings.
type SavePersonaGenerationPhaseAISettingsRequestDTO struct {
	JobID         int64  `json:"jobId"`
	Provider      string `json:"provider"`
	Model         string `json:"model"`
	ExecutionMode string `json:"executionMode"`
	BatchMode     string `json:"batchMode"`
}

// PersonaGenerationPhaseAISettingsResponseDTO returns public AI settings state.
type PersonaGenerationPhaseAISettingsResponseDTO struct {
	JobID            int64  `json:"jobId"`
	PhaseID          string `json:"phaseId"`
	Provider         string `json:"provider"`
	Model            string `json:"model"`
	CredentialStatus string `json:"credentialStatus"`
	ExecutionMode    string `json:"executionMode"`
	BatchMode        string `json:"batchMode"`
	ModelListStatus  string `json:"modelListStatus"`
}

// PersonaGenerationPhaseProgressSummaryDTO summarizes one phase run progress snapshot.
type PersonaGenerationPhaseProgressSummaryDTO struct {
	Percent        int    `json:"percent"`
	ProcessedCount int    `json:"processedCount"`
	TotalCount     int    `json:"totalCount"`
	TargetCount    int    `json:"targetCount"`
	CurrentStep    string `json:"currentStep"`
}

// PersonaGenerationTargetSummaryDTO summarizes target snapshot and skipped-target counts.
type PersonaGenerationTargetSummaryDTO struct {
	TargetCount            int      `json:"targetCount"`
	CommonPersonaHitCount  int      `json:"commonPersonaHitCount"`
	CommonPersonaMissCount int      `json:"commonPersonaMissCount"`
	SkippedCount           int      `json:"skippedCount"`
	SkippedReasons         []string `json:"skippedReasons"`
	TargetSnapshotID       *string  `json:"targetSnapshotId,omitempty"`
	TargetSnapshotDigest   string   `json:"targetSnapshotDigest"`
}

// PersonaGenerationExecutionSummaryDTO summarizes execution inputs without exposing secrets.
type PersonaGenerationExecutionSummaryDTO struct {
	CredentialRef string   `json:"credentialRef"`
	Provider      string   `json:"provider"`
	Model         string   `json:"model"`
	ExecutionMode string   `json:"executionMode"`
	PromptDigest  string   `json:"promptDigest"`
	InputCount    int      `json:"inputCount"`
	OutputCount   int      `json:"outputCount"`
	EvidenceRefs  []string `json:"evidenceRefs"`
}

// PersonaGenerationPhaseResultSummaryDTO summarizes the stored persona snapshot.
type PersonaGenerationPhaseResultSummaryDTO struct {
	GeneratedCount          int    `json:"generatedCount"`
	FailedCount             int    `json:"failedCount"`
	PersonaCount            int    `json:"personaCount"`
	MissingCount            int    `json:"missingCount"`
	SnapshotID              string `json:"snapshotId"`
	SnapshotDigest          string `json:"snapshotDigest"`
	SnapshotReferenceStatus string `json:"snapshotReferenceStatus"`
	BodyReadiness           bool   `json:"bodyReadiness"`
}

// PersonaGenerationPhaseErrorSummaryDTO exposes only redacted error information.
type PersonaGenerationPhaseErrorSummaryDTO struct {
	ErrorKind  string `json:"errorKind"`
	Reason     string `json:"reason"`
	Retryable  bool   `json:"retryable"`
	IsRedacted bool   `json:"isRedacted"`
}

// PersonaGenerationPhaseActionEnablementDTO summarizes Job Run operation button state.
type PersonaGenerationPhaseActionEnablementDTO struct {
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

// PersonaGenerationPhaseSummaryResponseDTO returns the frozen summary response shape.
type PersonaGenerationPhaseSummaryResponseDTO struct {
	JobID            int64                                     `json:"jobId"`
	CurrentPhase     string                                    `json:"currentPhase"`
	PhaseState       string                                    `json:"phaseState"`
	PhaseRunID       *int64                                    `json:"phaseRunId,omitempty"`
	StartedAt        *string                                   `json:"startedAt,omitempty"`
	FinishedAt       *string                                   `json:"finishedAt,omitempty"`
	Progress         PersonaGenerationPhaseProgressSummaryDTO  `json:"progress"`
	TargetSummary    PersonaGenerationTargetSummaryDTO         `json:"targetSummary"`
	Execution        PersonaGenerationExecutionSummaryDTO      `json:"execution"`
	ResultSummary    *PersonaGenerationPhaseResultSummaryDTO   `json:"resultSummary,omitempty"`
	ErrorSummary     *PersonaGenerationPhaseErrorSummaryDTO    `json:"errorSummary,omitempty"`
	ActionEnablement PersonaGenerationPhaseActionEnablementDTO `json:"actionEnablement"`
}

// PersonaGenerationPhaseCommandResponseDTO returns the frozen write-seam response shape.
type PersonaGenerationPhaseCommandResponseDTO struct {
	JobID         int64                                    `json:"jobId"`
	CurrentPhase  string                                   `json:"currentPhase"`
	PhaseState    string                                   `json:"phaseState"`
	PhaseRunID    *int64                                   `json:"phaseRunId,omitempty"`
	StartedAt     *string                                  `json:"startedAt,omitempty"`
	FinishedAt    *string                                  `json:"finishedAt,omitempty"`
	Progress      PersonaGenerationPhaseProgressSummaryDTO `json:"progress"`
	TargetSummary PersonaGenerationTargetSummaryDTO        `json:"targetSummary"`
	Execution     PersonaGenerationExecutionSummaryDTO     `json:"execution"`
	ResultSummary *PersonaGenerationPhaseResultSummaryDTO  `json:"resultSummary,omitempty"`
	Retryable     bool                                     `json:"retryable"`
	ErrorSummary  *PersonaGenerationPhaseErrorSummaryDTO   `json:"errorSummary,omitempty"`
}

// PersonaGenerationBodyReadinessInputSummaryDTO summarizes the body phase inputs.
type PersonaGenerationBodyReadinessInputSummaryDTO struct {
	PersonaCount   int      `json:"personaCount"`
	MissingCount   int      `json:"missingCount"`
	SnapshotID     string   `json:"snapshotId"`
	SnapshotDigest string   `json:"snapshotDigest"`
	EvidenceRefs   []string `json:"evidenceRefs"`
}

// PersonaGenerationBodyReadinessResponseDTO returns the frozen body phase fact state shape.
type PersonaGenerationBodyReadinessResponseDTO struct {
	JobID         int64                                         `json:"jobId"`
	CurrentPhase  string                                        `json:"currentPhase"`
	PhaseState    string                                        `json:"phaseState"`
	JobIsTerminal bool                                          `json:"jobIsTerminal"`
	InputSummary  PersonaGenerationBodyReadinessInputSummaryDTO `json:"inputSummary"`
}

// NewPersonaGenerationPhaseController creates a persona generation phase controller.
func NewPersonaGenerationPhaseController(usecase PersonaGenerationPhaseUsecasePort) *PersonaGenerationPhaseController {
	return &PersonaGenerationPhaseController{personaGenerationPhaseUsecase: usecase}
}

// GetPersonaGenerationPhaseSummary returns the frozen Job Run summary contract.
func (controller *PersonaGenerationPhaseController) GetPersonaGenerationPhaseSummary(
	request GetPersonaGenerationPhaseSummaryRequestDTO,
) (PersonaGenerationPhaseSummaryResponseDTO, error) {
	result, err := controller.personaGenerationPhaseUsecase.GetPersonaGenerationPhaseSummary(
		context.Background(),
		usecase.GetPersonaGenerationPhaseSummaryRequest{JobID: request.JobID},
	)
	if err != nil {
		return PersonaGenerationPhaseSummaryResponseDTO{}, fmt.Errorf("get persona generation phase summary: %w", err)
	}
	return toPersonaGenerationPhaseSummaryResponseDTO(result), nil
}

// StartPersonaGenerationPhase starts one persona generation phase run or returns a rejected result.
func (controller *PersonaGenerationPhaseController) StartPersonaGenerationPhase(
	request StartPersonaGenerationPhaseRequestDTO,
) (PersonaGenerationPhaseCommandResponseDTO, error) {
	result, err := controller.personaGenerationPhaseUsecase.StartPersonaGenerationPhase(
		context.Background(),
		usecase.StartPersonaGenerationPhaseRequest{JobID: request.JobID},
	)
	if err != nil {
		return toPersonaGenerationPhaseCommandResponseDTO(result), fmt.Errorf("start persona generation phase: %w", err)
	}
	return toPersonaGenerationPhaseCommandResponseDTO(result), nil
}

// PausePersonaGenerationPhase pauses one active persona generation phase run.
func (controller *PersonaGenerationPhaseController) PausePersonaGenerationPhase(
	request PausePersonaGenerationPhaseRequestDTO,
) (PersonaGenerationPhaseCommandResponseDTO, error) {
	result, err := controller.personaGenerationPhaseUsecase.PausePersonaGenerationPhase(
		context.Background(),
		usecase.PausePersonaGenerationPhaseRequest{JobID: request.JobID, PhaseRunID: request.PhaseRunID},
	)
	if err != nil {
		return PersonaGenerationPhaseCommandResponseDTO{}, fmt.Errorf("pause persona generation phase: %w", err)
	}
	return toPersonaGenerationPhaseCommandResponseDTO(result), nil
}

// ResumePersonaGenerationPhase resumes one paused or recoverable failed persona generation phase run.
func (controller *PersonaGenerationPhaseController) ResumePersonaGenerationPhase(
	request ResumePersonaGenerationPhaseRequestDTO,
) (PersonaGenerationPhaseCommandResponseDTO, error) {
	result, err := controller.personaGenerationPhaseUsecase.ResumePersonaGenerationPhase(
		context.Background(),
		usecase.ResumePersonaGenerationPhaseRequest{JobID: request.JobID, PhaseRunID: request.PhaseRunID},
	)
	if err != nil {
		return PersonaGenerationPhaseCommandResponseDTO{}, fmt.Errorf("resume persona generation phase: %w", err)
	}
	return toPersonaGenerationPhaseCommandResponseDTO(result), nil
}

// RetryPersonaGenerationPhase retries one retryable failed persona generation phase run.
func (controller *PersonaGenerationPhaseController) RetryPersonaGenerationPhase(
	request RetryPersonaGenerationPhaseRequestDTO,
) (PersonaGenerationPhaseCommandResponseDTO, error) {
	result, err := controller.personaGenerationPhaseUsecase.RetryPersonaGenerationPhase(
		context.Background(),
		usecase.RetryPersonaGenerationPhaseRequest{JobID: request.JobID, PhaseRunID: request.PhaseRunID},
	)
	if err != nil {
		return toPersonaGenerationPhaseCommandResponseDTO(result), fmt.Errorf("retry persona generation phase: %w", err)
	}
	return toPersonaGenerationPhaseCommandResponseDTO(result), nil
}

// CancelPersonaGenerationPhase cancels one active persona generation phase run.
func (controller *PersonaGenerationPhaseController) CancelPersonaGenerationPhase(
	request CancelPersonaGenerationPhaseRequestDTO,
) (PersonaGenerationPhaseCommandResponseDTO, error) {
	result, err := controller.personaGenerationPhaseUsecase.CancelPersonaGenerationPhase(
		context.Background(),
		usecase.CancelPersonaGenerationPhaseRequest{JobID: request.JobID, PhaseRunID: request.PhaseRunID},
	)
	if err != nil {
		return PersonaGenerationPhaseCommandResponseDTO{}, fmt.Errorf("cancel persona generation phase: %w", err)
	}
	return toPersonaGenerationPhaseCommandResponseDTO(result), nil
}

// GetPersonaGenerationBodyReadiness returns whether the body phase can start.
func (controller *PersonaGenerationPhaseController) GetPersonaGenerationBodyReadiness(
	request GetPersonaGenerationBodyReadinessRequestDTO,
) (PersonaGenerationBodyReadinessResponseDTO, error) {
	result, err := controller.personaGenerationPhaseUsecase.GetPersonaGenerationBodyReadiness(
		context.Background(),
		usecase.GetPersonaGenerationBodyReadinessRequest{JobID: request.JobID},
	)
	if err != nil {
		return toPersonaGenerationBodyReadinessResponseDTO(result), fmt.Errorf("get persona generation body readiness: %w", err)
	}
	return toPersonaGenerationBodyReadinessResponseDTO(result), nil
}

// SavePersonaGenerationPhaseAISettings saves public AI settings for the persona phase.
func (controller *PersonaGenerationPhaseController) SavePersonaGenerationPhaseAISettings(
	request SavePersonaGenerationPhaseAISettingsRequestDTO,
) (PersonaGenerationPhaseAISettingsResponseDTO, error) {
	settingsUsecase, ok := controller.personaGenerationPhaseUsecase.(personaGenerationPhaseAISettingsUsecasePort)
	if !ok {
		return PersonaGenerationPhaseAISettingsResponseDTO{}, fmt.Errorf("save persona generation phase ai settings: usecase is not configured")
	}
	result, err := settingsUsecase.SavePersonaGenerationPhaseAISettings(
		context.Background(),
		usecase.SavePersonaGenerationPhaseAISettingsRequest{
			JobID:         request.JobID,
			Provider:      request.Provider,
			Model:         request.Model,
			ExecutionMode: request.ExecutionMode,
			BatchMode:     request.BatchMode,
		},
	)
	if err != nil {
		return PersonaGenerationPhaseAISettingsResponseDTO{}, fmt.Errorf("save persona generation phase ai settings: %w", err)
	}
	return PersonaGenerationPhaseAISettingsResponseDTO(result), nil
}

func toPersonaGenerationPhaseSummaryResponseDTO(
	result usecase.PersonaGenerationPhaseSummaryResult,
) PersonaGenerationPhaseSummaryResponseDTO {
	return PersonaGenerationPhaseSummaryResponseDTO{
		JobID:            result.JobID,
		CurrentPhase:     result.CurrentPhase,
		PhaseState:       result.PhaseState,
		PhaseRunID:       cloneOptionalInt64(result.PhaseRunID),
		StartedAt:        formatOptionalTime(result.StartedAt),
		FinishedAt:       formatOptionalTime(result.FinishedAt),
		Progress:         toPersonaGenerationPhaseProgressSummaryDTO(result.Progress),
		TargetSummary:    toPersonaGenerationTargetSummaryDTO(result.TargetSummary),
		Execution:        toPersonaGenerationExecutionSummaryDTO(result.Execution),
		ResultSummary:    toOptionalPersonaGenerationPhaseResultSummaryDTO(result.ResultSummary),
		ErrorSummary:     toOptionalPersonaGenerationPhaseErrorSummaryDTO(result.ErrorSummary),
		ActionEnablement: toPersonaGenerationPhaseActionEnablementDTO(result.ActionEnablement),
	}
}

func toPersonaGenerationPhaseCommandResponseDTO(
	result usecase.PersonaGenerationPhaseCommandResult,
) PersonaGenerationPhaseCommandResponseDTO {
	return PersonaGenerationPhaseCommandResponseDTO{
		JobID:         result.JobID,
		CurrentPhase:  result.CurrentPhase,
		PhaseState:    result.PhaseState,
		PhaseRunID:    cloneOptionalInt64(result.PhaseRunID),
		StartedAt:     formatOptionalTime(result.StartedAt),
		FinishedAt:    formatOptionalTime(result.FinishedAt),
		Progress:      toPersonaGenerationPhaseProgressSummaryDTO(result.Progress),
		TargetSummary: toPersonaGenerationTargetSummaryDTO(result.TargetSummary),
		Execution:     toPersonaGenerationExecutionSummaryDTO(result.Execution),
		ResultSummary: toOptionalPersonaGenerationPhaseResultSummaryDTO(result.ResultSummary),
		Retryable:     result.Retryable,
		ErrorSummary:  toOptionalPersonaGenerationPhaseErrorSummaryDTO(result.ErrorSummary),
	}
}

func toPersonaGenerationBodyReadinessResponseDTO(
	result usecase.PersonaGenerationBodyReadinessResult,
) PersonaGenerationBodyReadinessResponseDTO {
	return PersonaGenerationBodyReadinessResponseDTO{
		JobID:         result.JobID,
		CurrentPhase:  result.CurrentPhase,
		PhaseState:    result.PhaseState,
		JobIsTerminal: result.JobIsTerminal,
		InputSummary: PersonaGenerationBodyReadinessInputSummaryDTO{
			PersonaCount:   result.InputSummary.PersonaCount,
			MissingCount:   result.InputSummary.MissingCount,
			SnapshotID:     result.InputSummary.SnapshotID,
			SnapshotDigest: result.InputSummary.SnapshotDigest,
			EvidenceRefs:   toPersonaGenerationStringArrayDTO(result.InputSummary.EvidenceRefs),
		},
	}
}

func toPersonaGenerationPhaseProgressSummaryDTO(
	summary usecase.PersonaGenerationPhaseProgressSummary,
) PersonaGenerationPhaseProgressSummaryDTO {
	return PersonaGenerationPhaseProgressSummaryDTO{
		Percent:        summary.Percent,
		ProcessedCount: summary.ProcessedCount,
		TotalCount:     summary.TotalCount,
		TargetCount:    summary.TargetCount,
		CurrentStep:    summary.CurrentStep,
	}
}

func toPersonaGenerationTargetSummaryDTO(
	summary usecase.PersonaGenerationTargetSummary,
) PersonaGenerationTargetSummaryDTO {
	return PersonaGenerationTargetSummaryDTO{
		TargetCount:            summary.TargetCount,
		CommonPersonaHitCount:  summary.CommonPersonaHitCount,
		CommonPersonaMissCount: summary.CommonPersonaMissCount,
		SkippedCount:           summary.SkippedCount,
		SkippedReasons:         toPersonaGenerationStringArrayDTO(summary.SkippedReasons),
		TargetSnapshotID:       cloneOptionalString(summary.TargetSnapshotID),
		TargetSnapshotDigest:   summary.TargetSnapshotDigest,
	}
}

func toPersonaGenerationExecutionSummaryDTO(
	summary usecase.PersonaGenerationExecutionSummary,
) PersonaGenerationExecutionSummaryDTO {
	return PersonaGenerationExecutionSummaryDTO{
		CredentialRef: summary.CredentialRef,
		Provider:      summary.Provider,
		Model:         summary.Model,
		ExecutionMode: summary.ExecutionMode,
		PromptDigest:  summary.PromptDigest,
		InputCount:    summary.InputCount,
		OutputCount:   summary.OutputCount,
		EvidenceRefs:  toPersonaGenerationStringArrayDTO(summary.EvidenceRefs),
	}
}

func toPersonaGenerationStringArrayDTO(values []string) []string {
	if values == nil {
		return []string{}
	}
	return append([]string{}, values...)
}

func toOptionalPersonaGenerationPhaseResultSummaryDTO(
	summary *usecase.PersonaGenerationPhaseResultSummary,
) *PersonaGenerationPhaseResultSummaryDTO {
	if summary == nil {
		return nil
	}
	return &PersonaGenerationPhaseResultSummaryDTO{
		GeneratedCount:          summary.GeneratedCount,
		FailedCount:             summary.FailedCount,
		PersonaCount:            summary.PersonaCount,
		MissingCount:            summary.MissingCount,
		SnapshotID:              summary.SnapshotID,
		SnapshotDigest:          summary.SnapshotDigest,
		SnapshotReferenceStatus: summary.SnapshotReferenceStatus,
		BodyReadiness:           summary.BodyReadiness,
	}
}

func toOptionalPersonaGenerationPhaseErrorSummaryDTO(
	summary *usecase.PersonaGenerationPhaseErrorSummary,
) *PersonaGenerationPhaseErrorSummaryDTO {
	if summary == nil {
		return nil
	}
	return &PersonaGenerationPhaseErrorSummaryDTO{
		ErrorKind:  usecase.NormalizePersonaGenerationPhasePublicErrorKind(summary.ErrorKind),
		Reason:     summary.Reason,
		Retryable:  summary.Retryable,
		IsRedacted: summary.IsRedacted,
	}
}

func toPersonaGenerationPhaseActionEnablementDTO(
	enablement usecase.PersonaGenerationPhaseActionEnablement,
) PersonaGenerationPhaseActionEnablementDTO {
	return PersonaGenerationPhaseActionEnablementDTO{
		CanStart:            enablement.CanStart,
		StartBlockedReason:  cloneOptionalString(enablement.StartBlockedReason),
		CanPause:            enablement.CanPause,
		PauseBlockedReason:  cloneOptionalString(enablement.PauseBlockedReason),
		CanResume:           enablement.CanResume,
		ResumeBlockedReason: cloneOptionalString(enablement.ResumeBlockedReason),
		CanRetry:            enablement.CanRetry,
		RetryBlockedReason:  cloneOptionalString(enablement.RetryBlockedReason),
		CanCancel:           enablement.CanCancel,
		CancelBlockedReason: cloneOptionalString(enablement.CancelBlockedReason),
	}
}
