import type {
  ProcessingTargetListRequest,
  ProcessingTargetListResponse
} from "@application/gateway-contract/processing-target"

export type PersonaGenerationPhaseErrorKind =
  | "term_phase_incomplete"
  | "terminal_job"
  | "active_phase_exists"
  | "target_snapshot_failed"
  | "provider_failure"
  | "invalid_provider_response"
  | "save_failed"
  | "snapshot_missing"
  | "body_readiness_blocked"
  | "secret_redacted"
  | "input_missing"

// ドメイン状態射影: UX 遷移可否導出にだけ使う、enum / 数値 / 真偽 の最小集合（design-diff C-2）
export interface PersonaGenerationPhaseProjection {
  phaseLifecycle: string
  jobLifecycle: string
  errorKind: string
  aiSettingsConfigured: boolean
  targetCount: number
  previousPhaseLifecycle: string
}

export interface GetPersonaGenerationPhaseSummaryRequest {
  jobId: number
}

export interface StartPersonaGenerationPhaseRequest {
  jobId: number
}

export interface PersonaGenerationPhaseAISettingsRequest {
  provider: string
  model: string
  executionMode: string
  batchMode: string
}

export interface PersonaGenerationPhaseAISettingsResponse {
  phaseType: string
  provider: string
  model: string
  executionMode: string
  batchMode: string
}

export interface PersonaGenerationPhaseAISettings {
  provider: string
  model: string
  executionMode: string
  batchMode: string
}

export interface PausePersonaGenerationPhaseRequest {
  jobId: number
  phaseRunId: number
}

export interface ResumePersonaGenerationPhaseRequest {
  jobId: number
  phaseRunId: number
}

export interface RetryPersonaGenerationPhaseRequest {
  jobId: number
  phaseRunId: number
}

export interface CancelPersonaGenerationPhaseRequest {
  jobId: number
  phaseRunId: number
}

export interface GetPersonaGenerationBodyReadinessRequest {
  jobId: number
}

export interface PersonaGenerationPhaseProgressSummary {
  percent: number
  processedCount: number
  totalCount: number
  targetCount: number
  currentStep: string
}

export interface PersonaGenerationTargetSummary {
  targetCount: number
  commonPersonaHitCount: number
  commonPersonaMissCount: number
  skippedCount: number
  skippedReasons: string[]
  targetSnapshotId?: string
  targetSnapshotDigest: string
}

export interface PersonaGenerationExecutionSummary {
  credentialRef: string
  provider: string
  model: string
  executionMode: string
  promptDigest: string
  inputCount: number
  outputCount: number
  evidenceRefs: string[]
}

export interface PersonaGenerationPhaseResultSummary {
  generatedCount: number
  failedCount: number
  personaCount: number
  missingCount: number
  snapshotId: string
  snapshotDigest: string
  snapshotReferenceStatus: string
  bodyReadiness: boolean
}

export interface PersonaGenerationPhaseErrorSummary {
  errorKind: PersonaGenerationPhaseErrorKind
  reason: string
  retryable: boolean
  isRedacted: boolean
}

export interface PersonaGenerationPhaseSummaryResponse {
  jobId: number
  currentPhase: string
  phaseState: string
  phaseRunId?: number
  startedAt?: string
  finishedAt?: string
  progress: PersonaGenerationPhaseProgressSummary
  targetSummary: PersonaGenerationTargetSummary
  aiSettings?: PersonaGenerationPhaseAISettings
  execution?: PersonaGenerationExecutionSummary
  resultSummary?: PersonaGenerationPhaseResultSummary
  errorSummary?: PersonaGenerationPhaseErrorSummary
}

// fetch 系 response: projection（UX 遷移可否判定入力）と summary（表示入力）の 2 系統を並べる（design-diff C-2 / G-4-b）
export interface PersonaGenerationPhaseFetchResponse {
  projection: PersonaGenerationPhaseProjection
  summary: PersonaGenerationPhaseSummaryResponse
}

export interface PersonaGenerationPhaseCommandResponse {
  jobId: number
  currentPhase: string
  phaseState: string
  phaseRunId?: number
  startedAt?: string
  finishedAt?: string
  progress: PersonaGenerationPhaseProgressSummary
  targetSummary: PersonaGenerationTargetSummary
  execution?: PersonaGenerationExecutionSummary
  resultSummary?: PersonaGenerationPhaseResultSummary
  retryable: boolean
  errorSummary?: PersonaGenerationPhaseErrorSummary
}

export interface PersonaGenerationBodyReadinessInputSummary {
  personaCount: number
  missingCount: number
  snapshotId: string
  snapshotDigest: string
  evidenceRefs: string[]
}

export interface PersonaGenerationBodyReadinessResponse {
  jobId: number
  currentPhase: string
  phaseState: string
  inputSummary: PersonaGenerationBodyReadinessInputSummary
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

export interface PersonaGenerationPhaseGatewayContract {
  listProviderModels?(
    request: PhaseProviderModelsRequest
  ): Promise<PhaseProviderModelsResponse>
  getProcessingTargetList?(
    request: ProcessingTargetListRequest
  ): Promise<ProcessingTargetListResponse>
  getPersonaGenerationPhaseSummary(
    request: GetPersonaGenerationPhaseSummaryRequest
  ): Promise<PersonaGenerationPhaseFetchResponse>
  startPersonaGenerationPhase(
    request: StartPersonaGenerationPhaseRequest
  ): Promise<PersonaGenerationPhaseCommandResponse>
  savePersonaGenerationPhaseAISettings?(
    request: PersonaGenerationPhaseAISettingsRequest
  ): Promise<PersonaGenerationPhaseAISettingsResponse>
  pausePersonaGenerationPhase(
    request: PausePersonaGenerationPhaseRequest
  ): Promise<PersonaGenerationPhaseCommandResponse>
  resumePersonaGenerationPhase(
    request: ResumePersonaGenerationPhaseRequest
  ): Promise<PersonaGenerationPhaseCommandResponse>
  retryPersonaGenerationPhase(
    request: RetryPersonaGenerationPhaseRequest
  ): Promise<PersonaGenerationPhaseCommandResponse>
  cancelPersonaGenerationPhase(
    request: CancelPersonaGenerationPhaseRequest
  ): Promise<PersonaGenerationPhaseCommandResponse>
  getPersonaGenerationBodyReadiness(
    request: GetPersonaGenerationBodyReadinessRequest
  ): Promise<PersonaGenerationBodyReadinessResponse>
}
