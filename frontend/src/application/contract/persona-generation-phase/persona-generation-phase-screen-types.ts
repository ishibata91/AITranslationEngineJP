import type {
  ProcessingTargetListPageState
} from "@application/gateway-contract/processing-target"

import type {
  GetPersonaGenerationBodyReadinessRequest,
  PausePersonaGenerationPhaseRequest,
  PersonaGenerationBodyReadinessInputSummary,
  PersonaGenerationBodyReadinessResponse,
  PersonaGenerationExecutionSummary,
  PersonaGenerationPhaseActionEnablement,
  PersonaGenerationPhaseErrorKind,
  PersonaGenerationPhaseErrorSummary,
  PersonaGenerationPhaseProgressSummary,
  PersonaGenerationPhaseResultSummary,
  PersonaGenerationPhaseSummaryResponse,
  PersonaGenerationTargetSummary,
  ResumePersonaGenerationPhaseRequest
} from "@application/gateway-contract/persona-generation-phase"

export type PersonaGenerationPhaseActionKind =
  | "start"
  | "pause"
  | "resume"
  | "retry"
  | "cancel"
  | "check-body-readiness"
  | "start-body-phase"

export type PersonaGenerationPhaseViewState =
  | "loading"
  | "not_started"
  | "running"
  | "paused"
  | "recoverable_failed"
  | "failed"
  | "completed"
  | "empty_completed"
  | "blocked"
  | "snapshot_missing"

export interface PersonaGenerationPhaseScreenState {
  jobId: number | null
  phase: "idle" | "loading" | "ready" | "submitting"
  summary: PersonaGenerationPhaseSummaryResponse | null
  bodyReadiness: PersonaGenerationBodyReadinessResponse | null
  errorMessage: string
  pendingAction: PersonaGenerationPhaseActionKind | null
  hasLoaded: boolean
  initialFetchDone: boolean
  processingTargetPageState?: ProcessingTargetListPageState | null
}

export interface PersonaGenerationPhaseActionCard {
  id: PersonaGenerationPhaseActionKind
  label: string
  disabled: boolean
  blockedReason: string
  tone: "default" | "primary" | "warning"
}

export interface PersonaGenerationPhaseScreenActionEnablement {
  canStart: boolean
  canPause: boolean
  canResume: boolean
  canRetry: boolean
  canCancel: boolean
  canCheckBodyReadiness: boolean
  canStartBodyPhase: boolean
}

export interface PersonaGenerationPhaseSelectOption {
  value: string
  label: string
}

export interface PersonaGenerationPhaseScreenViewModel extends PersonaGenerationPhaseScreenState {
  gatewayStatus: string
  viewState: PersonaGenerationPhaseViewState
  isLoading: boolean
  isRefreshing: boolean
  isSubmitting: boolean
  hasJobSelection: boolean
  currentPhaseLabel: string
  phaseStateLabel: string
  statusTitle: string
  statusText: string
  progressPercent: number
  progressLabel: string
  progressDetail: string
  startedAtLabel: string
  finishedAtLabel: string
  targetCountLabel: string
  generatedCountLabel: string
  failedCountLabel: string
  skippedCountLabel: string
  npcCountLabel: string
  commonPersonaHitCountLabel: string
  commonPersonaMissCountLabel: string
  skippedReasonsLabel: string
  targetSnapshotLabel: string
  isExecutionConfigured: boolean
  providerLabel: string
  modelLabel: string
  executionModeLabel: string
  credentialRefLabel: string
  providerOptions: PersonaGenerationPhaseSelectOption[]
  modelOptions: PersonaGenerationPhaseSelectOption[]
  executionOptions: PersonaGenerationPhaseSelectOption[]
  inputCountLabel: string
  outputCountLabel: string
  evidenceRefsLabel: string
  promptDigestLabel: string
  snapshotLabel: string
  snapshotReferenceStatusLabel: string
  personaCountLabel: string
  missingCountLabel: string
  bodyReadinessLabel: string
  bodyReadinessBlockedReason: string
  bodyReadinessInputSummaryLabel: string
  errorKindLabel: string
  errorReasonLabel: string
  retryableLabel: string
  actionCards: PersonaGenerationPhaseActionCard[]
  screenActionEnablement: PersonaGenerationPhaseScreenActionEnablement
  lastErrorSummary: PersonaGenerationPhaseErrorSummary | null
  actionEnablement: PersonaGenerationPhaseActionEnablement | null
  latestProgressSummary: PersonaGenerationPhaseProgressSummary | null
  latestTargetSummary: PersonaGenerationTargetSummary | null
  latestResultSummary: PersonaGenerationPhaseResultSummary | null
  latestExecutionSummary: PersonaGenerationExecutionSummary | null
  latestErrorKind: PersonaGenerationPhaseErrorKind | null
  latestBodyReadiness: PersonaGenerationBodyReadinessResponse | null
  latestBodyReadinessInputSummary: PersonaGenerationBodyReadinessInputSummary | null
  pauseRequestShape?: PausePersonaGenerationPhaseRequest
  resumeRequestShape?: ResumePersonaGenerationPhaseRequest
  bodyReadinessRequestShape?: GetPersonaGenerationBodyReadinessRequest
}
