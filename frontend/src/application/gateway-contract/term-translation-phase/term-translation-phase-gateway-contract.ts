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
  provider: string
  model: string
  executionMode: string
  batchMode: string
}

export interface TermTranslationPhaseAISettingsResponse {
  phaseType: string
  provider: string
  model: string
  executionMode: string
  batchMode: string
}

export interface TermTranslationPhaseAISettings {
  provider: string
  model: string
  executionMode: string
  batchMode: string
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

// ドメイン状態射影: UX 遷移可否導出にだけ使う、enum / 数値 / 真偽 の最小集合（design-diff C-1）
export interface TermTranslationPhaseProjection {
  phaseLifecycle: string
  jobLifecycle: string
  errorKind: string
  aiSettingsConfigured: boolean
  aiTargetCount: number
  confirmedCount: number
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
  aiSettings?: TermTranslationPhaseAISettings
  execution?: TermTranslationExecutionConfigSummary
  resultSummary?: TermTranslationPhaseResultSummary
  errorSummary?: TermTranslationPhaseErrorSummary
  projection: TermTranslationPhaseProjection
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
  errorSummary?: TermTranslationPhaseErrorSummary
}

export interface TermTranslationNextPhaseReadinessResponse {
  jobId: number
  currentPhase: string
  phaseState: string
  jobIsTerminal: boolean
  totalCount: number
  confirmedCount: number
  errorKind?: string
}

export interface PhaseProviderModelsRequest {
  provider: string
  credentialStatus: string
  requestToken: string
}

export interface PhaseProviderModelOption {
  modelId: string
  label: string
}

export interface PhaseProviderModelsResponse {
  provider: string
  status: string
  models: PhaseProviderModelOption[]
  failureKind?: string
}

export interface TermTranslationPhaseGatewayContract {
  listProviderModels?(
    request: PhaseProviderModelsRequest
  ): Promise<PhaseProviderModelsResponse>
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
