import type {
  ProcessingTargetListRequest,
  ProcessingTargetListResponse
} from "@application/gateway-contract/processing-target"

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

export interface BodyTranslationPhaseAISettingsRequest {
  provider: string
  model: string
  executionMode: string
  batchMode: string
}

export interface BodyTranslationPhaseAISettingsResponse {
  phaseType: string
  provider: string
  model: string
  executionMode: string
  batchMode: string
}

export interface BodyTranslationPhaseAISettings {
  provider: string
  model: string
  executionMode: string
  batchMode: string
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

export interface PersonaBodyReadiness {
  bodyReadiness: boolean
  snapshotReferenceStatus: string
}

export interface BodyTranslationPhaseProjection {
  phaseLifecycle: string
  jobLifecycle: string
  errorKind: string
  aiSettingsConfigured: boolean
  targetCount: number
  previousPhaseLifecycle: string
  personaBodyReadiness: PersonaBodyReadiness
}

export interface BodyTranslationOutputReadinessSummary {
  completedFieldCount: number
  statusConsistent: boolean
  outputCount: number
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
  aiSettings?: BodyTranslationPhaseAISettings
  execution?: BodyTranslationPhaseExecutionSummary
  fieldResults?: BodyTranslationPhaseFieldResultItem[]
  resultSummary?: BodyTranslationPhaseFieldResultSummary
  errorSummary?: BodyTranslationPhaseErrorSummary
  projection: BodyTranslationPhaseProjection
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
  execution?: BodyTranslationPhaseExecutionSummary
  fieldResults?: BodyTranslationPhaseFieldResultItem[]
  resultSummary?: BodyTranslationPhaseFieldResultSummary
  retryable: boolean
  outputReadiness: BodyTranslationOutputReadinessSummary
  errorSummary?: BodyTranslationPhaseErrorSummary
}

/**
 * @deprecated 専用取得 endpoint は廃止済み。段階要約取得（getBodyTranslationPhaseSummary）から
 * outputReadiness フィールドで事実状態（completedFieldCount / statusConsistent / outputCount）を参照する。
 * この型は既存テストの互換性のために残しており、wave-4 でテスト更新後に完全廃止する。
 */
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

export interface BodyTranslationPhaseGatewayContract {
  listProviderModels?(
    request: PhaseProviderModelsRequest
  ): Promise<PhaseProviderModelsResponse>
  getProcessingTargetList?(
    request: ProcessingTargetListRequest
  ): Promise<ProcessingTargetListResponse>
  getBodyTranslationPhaseSummary(
    request: GetBodyTranslationPhaseSummaryRequest
  ): Promise<BodyTranslationPhaseSummaryResponse>
  startBodyTranslationPhase(
    request: StartBodyTranslationPhaseRequest
  ): Promise<BodyTranslationPhaseCommandResponse>
  saveBodyTranslationPhaseAISettings?(
    request: BodyTranslationPhaseAISettingsRequest
  ): Promise<BodyTranslationPhaseAISettingsResponse>
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
}
