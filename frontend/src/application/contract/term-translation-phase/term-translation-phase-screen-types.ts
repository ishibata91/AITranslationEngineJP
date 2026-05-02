import type {
  PauseTermTranslationPhaseRequest,
  ResumeTermTranslationPhaseRequest,
  TermTranslationExecutionConfigSummary,
  TermTranslationNextPhaseReadinessResponse,
  TermTranslationPhaseActionEnablement,
  TermTranslationPhaseErrorSummary,
  TermTranslationPhaseErrorKind,
  TermTranslationPhaseProgressSummary,
  TermTranslationPhaseResultSummary,
  TermTranslationPhaseSummaryResponse
} from "@application/gateway-contract/term-translation-phase"

export type TermTranslationPhaseActionKind =
  | "refresh"
  | "start"
  | "pause"
  | "resume"
  | "retry"

export type TermTranslationPhaseViewState =
  | "loading"
  | "idle_ready"
  | "running"
  | "empty_completed"
  | "completed"
  | "paused"
  | "recoverable_failed"
  | "blocked"

export interface TermTranslationPhaseScreenState {
  jobId: number | null
  phase: "idle" | "loading" | "ready" | "submitting"
  summary: TermTranslationPhaseSummaryResponse | null
  nextPhaseReadiness: TermTranslationNextPhaseReadinessResponse | null
  errorMessage: string
  pendingAction: TermTranslationPhaseActionKind | null
  hasLoaded: boolean
}

export interface TermTranslationPhaseActionCard {
  id: TermTranslationPhaseActionKind | "next-phase"
  label: string
  disabled: boolean
  blockedReason: string
  tone: "default" | "primary" | "warning"
}

export interface TermTranslationPhaseScreenViewModel extends TermTranslationPhaseScreenState {
  gatewayStatus: string
  viewState: TermTranslationPhaseViewState
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
  totalTermCountLabel: string
  dictionaryHitCountLabel: string
  aiTargetCountLabel: string
  confirmedCountLabel: string
  jobDictionaryAppliedCountLabel: string
  replacementTargetCountLabel: string
  unmatchedCountLabel: string
  providerLabel: string
  modelLabel: string
  executionModeLabel: string
  credentialRefLabel: string
  snapshotLabel: string
  errorKindLabel: string
  errorReasonLabel: string
  retryableLabel: string
  nextPhaseStatusLabel: string
  nextPhaseBlockedReason: string
  providerSkippedLabel: string
  actionCards: TermTranslationPhaseActionCard[]
  lastErrorSummary: TermTranslationPhaseErrorSummary | null
  actionEnablement: TermTranslationPhaseActionEnablement | null
  latestProgressSummary: TermTranslationPhaseProgressSummary | null
  latestResultSummary: TermTranslationPhaseResultSummary | null
  latestExecutionSummary: TermTranslationExecutionConfigSummary | null
  latestErrorKind: TermTranslationPhaseErrorKind | null
  pauseRequestShape?: PauseTermTranslationPhaseRequest
  resumeRequestShape?: ResumeTermTranslationPhaseRequest
}
