import type {
  ProcessingTargetListPageState,
  ProcessingTargetListPageStatesByPhase
} from "@application/gateway-contract/processing-target"

import type {
  BodyTranslationOutputReadinessSummary,
  BodyTranslationPhaseActionEnablement,
  BodyTranslationPhaseErrorKind,
  BodyTranslationPhaseErrorSummary,
  BodyTranslationPhaseExecutionSummary,
  BodyTranslationPhaseFieldResultItem as BodyTranslationPhaseGatewayFieldResultItem,
  BodyTranslationPhaseFieldResultSummary,
  BodyTranslationPhaseInputSummary,
  BodyTranslationPhaseProgressSummary,
  BodyTranslationPhaseRequestSummary,
  BodyTranslationPhaseSummaryResponse,
  CancelBodyTranslationPhaseRequest,
  PauseBodyTranslationPhaseRequest,
  ResumeBodyTranslationPhaseRequest,
  RetryBodyTranslationPhaseRequest
} from "@application/gateway-contract/body-translation-phase"

export type BodyTranslationPhaseActionKind =
  | "start"
  | "pause"
  | "resume"
  | "retry"
  | "cancel"
  | "check-output-readiness"

export type BodyTranslationPhaseViewState =
  | "loading"
  | "not_ready"
  | "ready"
  | "running"
  | "paused"
  | "recoverable_failed"
  | "validation_failed"
  | "empty_completed"
  | "completed"
  | "canceled"
  | "failed"

export interface BodyTranslationPhaseScreenState {
  jobId: number | null
  phase: "idle" | "loading" | "ready" | "submitting"
  summary: BodyTranslationPhaseSummaryResponse | null
  outputReadiness: BodyTranslationOutputReadinessSummary | null
  errorMessage: string
  pendingAction: BodyTranslationPhaseActionKind | null
  hasLoaded: boolean
  initialFetchDone: boolean
  processingTargetPageState?: ProcessingTargetListPageState | null
  processingTargetPageStatesByPhase?: ProcessingTargetListPageStatesByPhase
}

export interface BodyTranslationPhaseActionCard {
  id: BodyTranslationPhaseActionKind
  label: string
  disabled: boolean
  blockedReason: string
  tone: "default" | "primary" | "warning"
}

export interface BodyTranslationFieldResultItem {
  fieldId: string
  fieldLabel: string
  recordTypeLabel: string
  fieldTypeLabel: string
  formIdLabel: string
  editorIdLabel: string
  sourceExcerpt: string
  translatedText: string
  outputStatus: string
  protectionValidation: string
  retryCountLabel: string
  rawItem: BodyTranslationPhaseGatewayFieldResultItem | null
}

export interface BodyTranslationPhaseScreenActionEnablement {
  canStart: boolean
  canPause: boolean
  canResume: boolean
  canRetry: boolean
  canCancel: boolean
  canCheckOutputReadiness: boolean
}

export interface BodyTranslationPhaseScreenViewModel extends BodyTranslationPhaseScreenState {
  gatewayStatus: string
  viewState: BodyTranslationPhaseViewState
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
  processedCountLabel: string
  translatedCountLabel: string
  failedCountLabel: string
  skippedCountLabel: string
  dictionaryDigestLabel: string
  personaDigestLabel: string
  metadataDigestLabel: string
  promptDigestLabel: string
  inputSnapshotRefLabel: string
  skippedReasonsLabel: string
  providerLabel: string
  modelLabel: string
  executionModeLabel: string
  credentialRefLabel: string
  providerTargetCountLabel: string
  exactDictionaryExclusionCountLabel: string
  partialDictionaryConstraintCountLabel: string
  requestUnitCountLabel: string
  outputCountLabel: string
  resultOutputCountLabel: string
  providerStateLabel: string
  lateResponseLabel: string
  outputReadinessLabel: string
  outputReadinessBlockedReason: string
  outputReadinessCompletedFieldCountLabel: string
  outputReadinessStatusLabel: string
  errorKindLabel: string
  errorReasonLabel: string
  retryableLabel: string
  fieldResultAvailabilityLabel: string
  fieldResultItems: BodyTranslationFieldResultItem[]
  actionCards: BodyTranslationPhaseActionCard[]
  screenActionEnablement: BodyTranslationPhaseScreenActionEnablement
  lastErrorSummary: BodyTranslationPhaseErrorSummary | null
  actionEnablement: BodyTranslationPhaseActionEnablement | null
  latestProgressSummary: BodyTranslationPhaseProgressSummary | null
  latestInputSummary: BodyTranslationPhaseInputSummary | null
  latestRequestSummary: BodyTranslationPhaseRequestSummary | null
  latestExecutionSummary: BodyTranslationPhaseExecutionSummary | null
  latestResultSummary: BodyTranslationPhaseFieldResultSummary | null
  latestErrorKind: BodyTranslationPhaseErrorKind | null
  latestOutputReadiness: BodyTranslationOutputReadinessSummary | null
  pauseRequestShape?: PauseBodyTranslationPhaseRequest
  resumeRequestShape?: ResumeBodyTranslationPhaseRequest
  retryRequestShape?: RetryBodyTranslationPhaseRequest
  cancelRequestShape?: CancelBodyTranslationPhaseRequest
  /** @deprecated 専用取得 endpoint 廃止済み。wave-4 でテスト更新後に完全廃止する。 */
  outputReadinessRequestShape?: { jobId: number }
}
