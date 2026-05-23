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

export interface GetPersonaGenerationPhaseSummaryRequest {
  jobId: number
}

export interface StartPersonaGenerationPhaseRequest {
  jobId: number
}

export interface PersonaGenerationPhaseAISettingsRequest {
  jobId: number
  provider: string
  model: string
  executionMode: string
  batchMode: string
}

export interface PersonaGenerationPhaseAISettingsResponse extends PersonaGenerationPhaseAISettingsRequest {
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

export interface PersonaGenerationPhaseActionEnablement {
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
  canStartBodyPhase: boolean
  bodyPhaseBlockedReason?: string
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
  execution: PersonaGenerationExecutionSummary
  resultSummary?: PersonaGenerationPhaseResultSummary
  errorSummary?: PersonaGenerationPhaseErrorSummary
  actionEnablement: PersonaGenerationPhaseActionEnablement
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
  execution: PersonaGenerationExecutionSummary
  resultSummary?: PersonaGenerationPhaseResultSummary
  retryable: boolean
  canStartBodyPhase: boolean
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
  ready: boolean
  blockedReason?: string
  errorKind?: PersonaGenerationPhaseErrorKind
  inputSummary: PersonaGenerationBodyReadinessInputSummary
}

export interface PersonaGenerationPhaseGatewayContract {
  getPersonaGenerationPhaseSummary(
    request: GetPersonaGenerationPhaseSummaryRequest
  ): Promise<PersonaGenerationPhaseSummaryResponse>
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
