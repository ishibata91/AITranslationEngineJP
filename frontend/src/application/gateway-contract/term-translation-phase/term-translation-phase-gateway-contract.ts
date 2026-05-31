import type {
  ProcessingTargetListRequest,
  ProcessingTargetListResponse
} from "@application/gateway-contract/processing-target"

export type TermTranslationPhaseErrorKind =
  | "ready_required"
  | "terminal_job"
  | "active_phase_exists"
  | "dictionary_snapshot_failed"
  | "provider_failure"
  | "invalid_provider_response"
  | "save_failed"
  | "term_phase_incomplete"
  | "secret_redacted"

export interface GetTermTranslationPhaseSummaryRequest {
  jobId: number
}

export interface StartTermTranslationPhaseRequest {
  jobId: number
}

export interface TermTranslationPhaseAISettingsRequest {
  jobId: number
  provider: string
  model: string
  executionMode: string
  batchMode: string
}

export interface TermTranslationPhaseAISettingsResponse extends TermTranslationPhaseAISettingsRequest {
  phaseId: string
  credentialStatus: "configured" | "missing" | "not_required"
  modelListStatus:
    | "not_updated"
    | "loading"
    | "success"
    | "failed"
    | "credential_missing"
    | "credential_not_required"
}

export interface PauseTermTranslationPhaseRequest {
  jobId: number
  phaseRunId: number
}

export interface ResumeTermTranslationPhaseRequest {
  jobId: number
  phaseRunId: number
}

export interface RetryTermTranslationPhaseRequest {
  jobId: number
  phaseRunId: number
}

export interface GetTermTranslationNextPhaseReadinessRequest {
  jobId: number
}

export interface TermTranslationExecutionConfigSummary {
  credentialRef: string
  provider: string
  model: string
  executionMode: string
  snapshotDigest?: string
  snapshotVersion?: string
}

export interface TermTranslationPhaseProgressSummary {
  percent: number
  processedCount: number
  totalCount: number
  aiTargetCount: number
  currentStep: string
}

export interface TermTranslationPhaseResultSummary {
  confirmedCount: number
  jobDictionaryAppliedCount: number
  replacementTargetCount: number
  unmatchedCount: number
  providerSkipped: boolean
}

export interface TermTranslationPhaseErrorSummary {
  errorKind: TermTranslationPhaseErrorKind
  reason: string
  retryable: boolean
  isRedacted: boolean
}

export interface TermTranslationPhaseActionEnablement {
  canStart: boolean
  startBlockedReason?: string
  canPause: boolean
  pauseBlockedReason?: string
  canResume: boolean
  resumeBlockedReason?: string
  canRetry: boolean
  retryBlockedReason?: string
}

export interface TermTranslationPhaseSummaryResponse {
  jobId: number
  currentPhase: string
  phaseState: string
  phaseRunId?: number
  startedAt?: string
  finishedAt?: string
  progress: TermTranslationPhaseProgressSummary
  totalTermCount: number
  dictionaryHitCount: number
  aiTargetCount: number
  execution: TermTranslationExecutionConfigSummary
  resultSummary?: TermTranslationPhaseResultSummary
  errorSummary?: TermTranslationPhaseErrorSummary
  actionEnablement: TermTranslationPhaseActionEnablement
}

export interface TermTranslationPhaseCommandResponse {
  jobId: number
  currentPhase: string
  phaseState: string
  phaseRunId?: number
  startedAt?: string
  finishedAt?: string
  progress: TermTranslationPhaseProgressSummary
  retryable: boolean
  canStartNextPhase: boolean
  errorSummary?: TermTranslationPhaseErrorSummary
}

export interface TermTranslationNextPhaseReadinessResponse {
  jobId: number
  currentPhase: string
  phaseState: string
  errorKind?: TermTranslationPhaseErrorKind
}

export interface TermTranslationPhaseGatewayContract {
  getProcessingTargetList?(
    request: ProcessingTargetListRequest
  ): Promise<ProcessingTargetListResponse>
  getTermTranslationPhaseSummary(
    request: GetTermTranslationPhaseSummaryRequest
  ): Promise<TermTranslationPhaseSummaryResponse>
  startTermTranslationPhase(
    request: StartTermTranslationPhaseRequest
  ): Promise<TermTranslationPhaseCommandResponse>
  saveTermTranslationPhaseAISettings?(
    request: TermTranslationPhaseAISettingsRequest
  ): Promise<TermTranslationPhaseAISettingsResponse>
  pauseTermTranslationPhase(
    request: PauseTermTranslationPhaseRequest
  ): Promise<TermTranslationPhaseCommandResponse>
  resumeTermTranslationPhase(
    request: ResumeTermTranslationPhaseRequest
  ): Promise<TermTranslationPhaseCommandResponse>
  retryTermTranslationPhase(
    request: RetryTermTranslationPhaseRequest
  ): Promise<TermTranslationPhaseCommandResponse>
  getTermTranslationNextPhaseReadiness(
    request: GetTermTranslationNextPhaseReadinessRequest
  ): Promise<TermTranslationNextPhaseReadinessResponse>
}
