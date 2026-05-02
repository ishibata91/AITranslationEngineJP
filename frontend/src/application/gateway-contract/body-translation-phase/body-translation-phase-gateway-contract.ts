export type BodyTranslationPhaseErrorKind =
  | "persona_phase_incomplete"
  | "terminal_job"
  | "active_phase_exists"
  | "input_snapshot_failed"
  | "provider_failure"
  | "invalid_provider_response"
  | "protection_validation_failed"
  | "save_failed"
  | "output_readiness_blocked"
  | "late_response_rejected"
  | "secret_redacted"

export interface GetBodyTranslationPhaseSummaryRequest {
  jobId: number
}

export interface StartBodyTranslationPhaseRequest {
  jobId: number
}

export interface PauseBodyTranslationPhaseRequest {
  jobId: number
  phaseRunId: number
}

export interface ResumeBodyTranslationPhaseRequest {
  jobId: number
  phaseRunId: number
}

export interface RetryBodyTranslationPhaseRequest {
  jobId: number
  phaseRunId: number
}

export interface CancelBodyTranslationPhaseRequest {
  jobId: number
  phaseRunId: number
}

export interface GetBodyTranslationOutputReadinessRequest {
  jobId: number
}

export interface BodyTranslationPhaseProgressSummary {
  percent: number
  processedCount: number
  totalCount: number
  targetCount: number
  translatedCount: number
  skippedCount: number
  currentStep: string
}

export interface BodyTranslationPhaseInputSummary {
  targetCount: number
  skippedReasons?: string[]
  inputSnapshotRef?: string
  dictionaryDigest: string
  personaDigest: string
  metadataDigest: string
  promptDigest: string
}

export interface BodyTranslationPhaseRequestSummary {
  providerTargetCount?: number
  exactDictionaryExclusionCount?: number
  partialDictionaryConstraintCount?: number
}

export interface BodyTranslationPhaseExecutionSummary {
  credentialRef: string
  provider: string
  model: string
  executionMode: string
  requestUnitCount: number
  outputCount: number
}

export interface BodyTranslationPhaseFieldIdentity {
  translationFieldId?: number | string
  phaseTranslationFieldId?: number | string
  recordType?: string
  fieldType?: string
  formId?: string
  editorId?: string
  fieldLabel?: string
}

export interface BodyTranslationPhaseFieldResultItem {
  identity?: BodyTranslationPhaseFieldIdentity
  fieldId?: number | string
  fieldLabel?: string
  sourceExcerpt?: string
  translatedText?: string
  outputStatus?: string
  protectionValidationResult?: string
  protectionValidationSummary?: string
  retryCount?: number
}

export interface BodyTranslationPhaseFieldResultSummary {
  translatedCount: number
  failedCount: number
  skippedCount: number
  protectionFailedCount: number
  outputReadyCount: number
  outputCount?: number
  fieldResults?: BodyTranslationPhaseFieldResultItem[]
}

export interface BodyTranslationPhaseErrorSummary {
  errorKind: BodyTranslationPhaseErrorKind
  reason: string
  retryable: boolean
  isRedacted: boolean
}

export interface BodyTranslationPhaseActionEnablement {
  canStart: boolean
  startBlockedReason?: string
  canPause: boolean
  pauseBlockedReason?: string
  canResume: boolean
  resumeBlockedReason?: string
  canRetry: boolean
  retryBlockedReason?: string
  canCancel: boolean
  cancelBlockedReason?: string
  canCheckOutputReadiness: boolean
  outputReadinessBlockedReason?: string
}

export interface BodyTranslationOutputReadinessSummary {
  ready: boolean
  blockedReason?: string
  errorKind?: BodyTranslationPhaseErrorKind
  completedFieldCount: number
  statusConsistent: boolean
}

export interface BodyTranslationPhaseSummaryResponse {
  jobId: number
  currentPhase: string
  phaseState: string
  phaseRunId?: number
  startedAt?: string
  finishedAt?: string
  progress: BodyTranslationPhaseProgressSummary
  inputSummary: BodyTranslationPhaseInputSummary
  requestSummary?: BodyTranslationPhaseRequestSummary
  execution: BodyTranslationPhaseExecutionSummary
  fieldResults?: BodyTranslationPhaseFieldResultItem[]
  resultSummary?: BodyTranslationPhaseFieldResultSummary
  errorSummary?: BodyTranslationPhaseErrorSummary
  actionEnablement: BodyTranslationPhaseActionEnablement
  outputReadiness: BodyTranslationOutputReadinessSummary
}

export interface BodyTranslationPhaseCommandResponse {
  jobId: number
  currentPhase: string
  phaseState: string
  phaseRunId?: number
  startedAt?: string
  finishedAt?: string
  progress: BodyTranslationPhaseProgressSummary
  inputSnapshotDigest?: string
  inputSummary: BodyTranslationPhaseInputSummary
  requestSummary?: BodyTranslationPhaseRequestSummary
  execution: BodyTranslationPhaseExecutionSummary
  fieldResults?: BodyTranslationPhaseFieldResultItem[]
  resultSummary?: BodyTranslationPhaseFieldResultSummary
  retryable: boolean
  outputReadiness: BodyTranslationOutputReadinessSummary
  errorSummary?: BodyTranslationPhaseErrorSummary
}

export interface BodyTranslationOutputReadinessResponse {
  jobId: number
  currentPhase: string
  phaseState: string
  ready: boolean
  blockedReason?: string
  errorKind?: BodyTranslationPhaseErrorKind
  completedFieldCount: number
  statusConsistent: boolean
  outputCount: number
}

export interface BodyTranslationPhaseGatewayContract {
  getBodyTranslationPhaseSummary(
    request: GetBodyTranslationPhaseSummaryRequest
  ): Promise<BodyTranslationPhaseSummaryResponse>
  startBodyTranslationPhase(
    request: StartBodyTranslationPhaseRequest
  ): Promise<BodyTranslationPhaseCommandResponse>
  pauseBodyTranslationPhase(
    request: PauseBodyTranslationPhaseRequest
  ): Promise<BodyTranslationPhaseCommandResponse>
  resumeBodyTranslationPhase(
    request: ResumeBodyTranslationPhaseRequest
  ): Promise<BodyTranslationPhaseCommandResponse>
  retryBodyTranslationPhase(
    request: RetryBodyTranslationPhaseRequest
  ): Promise<BodyTranslationPhaseCommandResponse>
  cancelBodyTranslationPhase(
    request: CancelBodyTranslationPhaseRequest
  ): Promise<BodyTranslationPhaseCommandResponse>
  getBodyTranslationOutputReadiness(
    request: GetBodyTranslationOutputReadinessRequest
  ): Promise<BodyTranslationOutputReadinessResponse>
}
