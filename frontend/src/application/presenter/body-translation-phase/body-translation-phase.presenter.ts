import type {
  ProcessingTargetListPageState,
  ProcessingTargetListPageStatesByPhase
} from "@application/gateway-contract/processing-target"

import type {
  BodyTranslationPhaseErrorKind,
  BodyTranslationPhaseErrorSummary,
  BodyTranslationPhaseExecutionSummary,
  BodyTranslationPhaseFieldResultItem as BodyTranslationPhaseGatewayFieldResultItem,
  BodyTranslationPhaseFieldResultSummary,
  BodyTranslationPhaseInputSummary,
  BodyTranslationPhaseProgressSummary,
  BodyTranslationPhaseProjection,
  BodyTranslationPhaseRequestSummary,
  BodyTranslationPhaseSummaryResponse
} from "@application/gateway-contract/body-translation-phase"

interface BodyTranslationEffectiveReadiness {
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

type BodyTranslationPhaseActionKind =
  | "start"
  | "pause"
  | "resume"
  | "retry"
  | "cancel"
  | "check-output-readiness"

type BodyTranslationPhaseViewState =
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

interface BodyTranslationPhaseScreenState {
  jobId: number | null
  phase: "idle" | "loading" | "ready" | "submitting"
  summary: BodyTranslationPhaseSummaryResponse | null
  outputReadiness: BodyTranslationEffectiveReadiness | null
  errorMessage: string
  pendingAction: BodyTranslationPhaseActionKind | null
  hasLoaded: boolean
  initialFetchDone: boolean
  processingTargetPageState?: ProcessingTargetListPageState | null
  processingTargetPageStatesByPhase?: ProcessingTargetListPageStatesByPhase
}

interface BodyTranslationFieldResultItem {
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

interface BodyTranslationPhaseScreenViewModel extends BodyTranslationPhaseScreenState {
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
  isExecutionConfigured: boolean
  providerLabel: string
  modelLabel: string
  executionModeLabel: string
  credentialRefLabel: string
  providerOptions: { value: string; label: string }[]
  modelOptions: { value: string; label: string }[]
  executionOptions: { value: string; label: string }[]
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
  actionCards: {
    id: BodyTranslationPhaseActionKind
    label: string
    disabled: boolean
    blockedReason: string
    tone: "default" | "primary" | "warning"
  }[]
  screenActionEnablement: {
    canStart: boolean
    canPause: boolean
    canResume: boolean
    canRetry: boolean
    canCancel: boolean
    canCheckOutputReadiness: boolean
  }
  lastErrorSummary: BodyTranslationPhaseErrorSummary | null
  projection: BodyTranslationPhaseProjection | null
  latestProgressSummary: BodyTranslationPhaseProgressSummary | null
  latestInputSummary: BodyTranslationPhaseInputSummary | null
  latestRequestSummary: BodyTranslationPhaseRequestSummary | null
  latestExecutionSummary: BodyTranslationPhaseExecutionSummary | null
  latestResultSummary: BodyTranslationPhaseFieldResultSummary | null
  latestErrorKind: BodyTranslationPhaseErrorKind | null
  latestOutputReadiness: BodyTranslationEffectiveReadiness | null
  pauseRequestShape?: { jobId: number; phaseRunId: number }
  resumeRequestShape?: { jobId: number; phaseRunId: number }
  retryRequestShape?: { jobId: number; phaseRunId: number }
  cancelRequestShape?: { jobId: number; phaseRunId: number }
  outputReadinessRequestShape?: { jobId: number }
}

// H 節 enum 集合定数（body phase 内に閉じる）
const TERMINAL_JOB = ["completed", "failed", "canceled"] as const
const RUNNING_PHASE = ["running", "in_progress", "processing"] as const
const IDLE_READY_PHASE = ["pending", "idle_ready", "ready", ""] as const
const PAUSED_PHASE = ["paused"] as const
const RECOVERABLE_FAILED_PHASE = ["recoverable_failed", "retryable_failed"] as const
const COMPLETED_PHASE = ["completed", "succeeded", "done"] as const

const PHASE_STATE_LABELS: Record<string, string> = {
  canceled: "キャンセル済み",
  cancelled: "キャンセル済み",
  completed: "完了",
  failed: "失敗",
  not_ready: "開始待ち",
  paused: "中断",
  pending: "開始待ち",
  ready: "開始可能",
  recoverable_failed: "再試行可能な失敗",
  running: "実行中",
  validation_failed: "検証失敗"
}

const CURRENT_PHASE_LABELS: Record<string, string> = {
  body_translation: "本文翻訳",
  translation_complete: "翻訳完了"
}

function formatDate(value: string | undefined): string {
  if (!value) {
    return "-"
  }

  const date = new Date(value)
  if (Number.isNaN(date.getTime())) {
    return value
  }

  return date.toLocaleString("ja-JP")
}

function formatCount(value: number | undefined): string {
  if (typeof value !== "number") {
    return "-"
  }

  return `${value.toLocaleString("ja-JP")} 件`
}

function formatReasonList(values: string[] | undefined): string {
  if (!values || values.length === 0) {
    return "-"
  }

  return values.join(" / ")
}

function normalizePhaseState(phaseState: string | undefined): string {
  return phaseState?.trim().toLowerCase().replaceAll(" ", "_") ?? ""
}

function isBodyExecutionConfigured(
  summary: BodyTranslationPhaseSummaryResponse
): boolean {
  if (summary.aiSettings) {
    return Boolean(summary.aiSettings.provider?.trim()) &&
      Boolean(summary.aiSettings.model?.trim()) &&
      Boolean(summary.aiSettings.executionMode?.trim())
  }
  if (!summary.execution) {
    return false
  }
  return Boolean(summary.execution.provider?.trim()) &&
    Boolean(summary.execution.model?.trim()) &&
    Boolean(summary.execution.executionMode?.trim())
}

const BODY_EXECUTION_OPTIONS = [
  { value: "batch", label: "バッチ処理" },
  { value: "sync", label: "逐次処理" }
]

function buildBodyProviderOptions(
  summary: BodyTranslationPhaseSummaryResponse | null,
  availableProviders: { value: string; label: string }[]
): { value: string; label: string }[] {
  if (availableProviders.length > 0) {
    const currentProvider = summary?.aiSettings?.provider ?? summary?.execution?.provider
    if (!currentProvider || !availableProviders.some((p) => p.value === currentProvider)) {
      return [{ value: "", label: "選んでください" }, ...availableProviders]
    }
    return availableProviders
  }
  const fallbackLabel = summary?.aiSettings?.provider ?? summary?.execution?.provider ?? "設定未完了"
  return [{ value: fallbackLabel, label: fallbackLabel }]
}

function buildBodyModelOptions(
  summary: BodyTranslationPhaseSummaryResponse | null,
  availableModels: { value: string; label: string }[]
): { value: string; label: string }[] {
  if (availableModels.length > 0) {
    return availableModels
  }
  const fallbackLabel = summary?.aiSettings?.model ?? summary?.execution?.model
  if (fallbackLabel) {
    return [{ value: fallbackLabel, label: fallbackLabel }]
  }
  return [{ value: "", label: "モデル一覧を更新してください" }]
}

function isBodyTerminalJob(jobLifecycle: string): boolean {
  return (TERMINAL_JOB as readonly string[]).includes(jobLifecycle)
}

function deriveBodyOutputReadinessReady(
  summary: BodyTranslationPhaseSummaryResponse
): boolean {
  const normalizedState = normalizePhaseState(summary.phaseState)
  const isCompleted =
    normalizedState === "completed" ||
    normalizedState === "succeeded" ||
    normalizedState === "done"
  if (!isCompleted) {
    return false
  }
  const completedFieldCount = summary.outputReadiness.completedFieldCount
  const statusConsistent = summary.outputReadiness.statusConsistent
  return statusConsistent === true && completedFieldCount >= 0
}

function deriveBodyOutputReadinessBlockedReason(
  summary: BodyTranslationPhaseSummaryResponse
): string {
  const normalizedState = normalizePhaseState(summary.phaseState)
  const isCompleted =
    normalizedState === "completed" ||
    normalizedState === "succeeded" ||
    normalizedState === "done"
  if (!isCompleted) {
    return "本文翻訳段階が完了していません。"
  }
  const statusConsistent = summary.outputReadiness.statusConsistent
  if (!statusConsistent) {
    return "翻訳出力の整合性が確認できません。"
  }
  return ""
}

interface BodyTranslationPhaseActionEnablement {
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
}

// H-12〜H-16 の論理式。入力は body projection のみ。summary は読まない。
function deriveBodyActionEnablement(
  projection: BodyTranslationPhaseProjection
): BodyTranslationPhaseActionEnablement {
  const terminal = isBodyTerminalJob(projection.jobLifecycle)

  const isRunning = (RUNNING_PHASE as readonly string[]).includes(projection.phaseLifecycle)
  const isPaused = (PAUSED_PHASE as readonly string[]).includes(projection.phaseLifecycle)
  const isRecoverableFailed =
    (RECOVERABLE_FAILED_PHASE as readonly string[]).includes(projection.phaseLifecycle) ||
    projection.errorKind === "recoverable"
  const isActive = isRunning || isPaused || isRecoverableFailed
  const isIdleReady = (IDLE_READY_PHASE as readonly string[]).includes(projection.phaseLifecycle)

  const previousPhaseCompleted = (COMPLETED_PHASE as readonly string[]).includes(projection.previousPhaseLifecycle)
  const hasProcessingTarget = projection.targetCount > 0
  const personaBodyReady =
    projection.personaBodyReadiness.bodyReadiness === true ||
    projection.personaBodyReadiness.snapshotReferenceStatus === "available"

  // H-12: 開始ボタン
  const canStart =
    !terminal &&
    !isRunning &&
    !isPaused &&
    !isRecoverableFailed &&
    isIdleReady &&
    projection.aiSettingsConfigured &&
    hasProcessingTarget &&
    previousPhaseCompleted &&
    personaBodyReady

  let startBlockedReason: string | undefined
  if (terminal) {
    startBlockedReason = "ジョブが終端状態のため開始できません。"
  } else if (isActive) {
    startBlockedReason = "実行中の翻訳段階があるため開始できません。"
  } else if (!previousPhaseCompleted) {
    startBlockedReason = "ペルソナ生成段階が完了していないため本文翻訳を開始できません。"
  } else if (!personaBodyReady) {
    startBlockedReason = "ペルソナ snapshot 参照が準備できていないため本文翻訳を開始できません。"
  } else if (!projection.aiSettingsConfigured) {
    startBlockedReason = "実行設定が未構成のため開始できません。"
  } else if (!hasProcessingTarget) {
    startBlockedReason = "処理対象が 0 件のため開始できません。"
  } else if (!isIdleReady) {
    startBlockedReason = "ジョブが開始可能状態ではありません。"
  }

  // H-13: 中断ボタン
  const canPause = !terminal && isRunning
  let pauseBlockedReason: string | undefined
  if (terminal) {
    pauseBlockedReason = "ジョブが終端状態のため中断できません。"
  } else if (!isRunning) {
    pauseBlockedReason = "フェーズが実行中ではありません。"
  }

  // H-14: 再開ボタン
  const canResume = !terminal && (isPaused || isRecoverableFailed)
  let resumeBlockedReason: string | undefined
  if (terminal) {
    resumeBlockedReason = "ジョブが終端状態のため再開できません。"
  } else if (!isPaused && !isRecoverableFailed) {
    resumeBlockedReason = "フェーズが再開可能な状態ではありません。"
  }

  // H-15: 再試行ボタン
  const canRetry = !terminal && isRecoverableFailed
  let retryBlockedReason: string | undefined
  if (terminal) {
    retryBlockedReason = "ジョブが終端状態のため再試行できません。"
  } else if (!isRecoverableFailed) {
    retryBlockedReason = "フェーズが再試行可能な状態ではありません。"
  }

  // H-16: キャンセルボタン
  const canCancel = !terminal && isActive
  let cancelBlockedReason: string | undefined
  if (terminal) {
    cancelBlockedReason = "ジョブが終端状態のためキャンセルできません。"
  } else if (!isActive) {
    cancelBlockedReason = "フェーズがキャンセル可能な状態ではありません。"
  }

  return {
    canStart,
    startBlockedReason,
    canPause,
    pauseBlockedReason,
    canResume,
    resumeBlockedReason,
    canRetry,
    retryBlockedReason,
    canCancel,
    cancelBlockedReason
  }
}

function getEffectiveReadiness(
  state: BodyTranslationPhaseScreenState
): BodyTranslationEffectiveReadiness | null {
  const summary = state.summary
  if (!summary) {
    return null
  }

  const ready = deriveBodyOutputReadinessReady(summary)
  const blockedReason = deriveBodyOutputReadinessBlockedReason(summary) || undefined

  return {
    jobId: summary.jobId,
    currentPhase: summary.currentPhase,
    phaseState: summary.phaseState,
    ready,
    blockedReason,
    completedFieldCount: summary.outputReadiness.completedFieldCount,
    statusConsistent: summary.outputReadiness.statusConsistent,
    outputCount: summary.outputReadiness.outputCount
  }
}

function buildViewState(
  state: BodyTranslationPhaseScreenState
): BodyTranslationPhaseViewState {
  if (state.phase === "loading" && !state.summary) {
    return "loading"
  }

  const summary = state.summary
  if (!summary) {
    return "not_ready"
  }

  if (summary.errorSummary?.errorKind === "protection_validation_failed") {
    return "validation_failed"
  }

  const normalizedState = normalizePhaseState(summary.phaseState)
  if (normalizedState === "pending") {
    return "ready"
  }

  if (normalizedState === "paused") {
    return "paused"
  }

  if (
    normalizedState === "recoverable_failed" ||
    summary.errorSummary?.retryable
  ) {
    return "recoverable_failed"
  }

  if (normalizedState === "validation_failed") {
    return "validation_failed"
  }

  if (normalizedState === "canceled" || normalizedState === "cancelled") {
    return "canceled"
  }

  if (normalizedState === "failed") {
    return "failed"
  }

  if (
    normalizedState === "running" ||
    normalizedState === "starting" ||
    normalizedState === "processing" ||
    normalizedState === "in_progress"
  ) {
    return "running"
  }

  if (
    normalizedState === "completed" ||
    normalizedState === "done" ||
    normalizedState === "succeeded"
  ) {
    return summary.progress.targetCount === 0 ? "empty_completed" : "completed"
  }

  const proj = summary.projection
  const enablement = deriveBodyActionEnablement(proj)
  if (enablement.canStart) {
    return "ready"
  }

  return "not_ready"
}

function buildPhaseStateLabel(phaseState: string | undefined): string {
  if (!phaseState) {
    return "未取得"
  }

  return PHASE_STATE_LABELS[normalizePhaseState(phaseState)] ?? phaseState
}

function buildCurrentPhaseLabel(currentPhase: string | undefined): string {
  const normalized = normalizePhaseState(currentPhase)
  return CURRENT_PHASE_LABELS[normalized] ?? (currentPhase || "本文翻訳")
}

const CURRENT_STEP_LABELS: Record<string, string> = {
  canceled: "キャンセル済み",
  completed: "完了",
  idle_ready: "開始待ち",
  paused: "中断",
  pending: "開始待ち",
  protection_validation_failed: "保護検証失敗",
  provider_request: "AI 処理中",
  ready: "開始可能",
  running: "実行中"
}

function buildCurrentStepLabel(currentStep: string | undefined): string {
  if (!currentStep) {
    return "-"
  }

  const normalized = normalizePhaseState(currentStep)
  return CURRENT_STEP_LABELS[normalized] ?? currentStep
}

function buildStatusCopy(
  state: BodyTranslationPhaseScreenState,
  viewState: BodyTranslationPhaseViewState
): { title: string; text: string } {
  if (state.jobId === null) {
    return {
      title: "ジョブ未選択",
      text: "ジョブIDを指定すると、本文翻訳段階の状態を確認できます。"
    }
  }

  if (state.phase === "loading" && !state.summary) {
    return {
      title: "summary を取得中",
      text: "現在の翻訳段階、進行状況、出力準備を読み込んでいます。"
    }
  }

  if (state.errorMessage.trim()) {
    return {
      title: "更新に注意が必要",
      text: state.errorMessage
    }
  }

  switch (viewState) {
    case "ready":
      return {
        title: "開始可能",
        text: "NPC ペルソナ生成段階の完了と前段参照成立が確認できたため、本文翻訳段階を開始できます。"
      }
    case "running":
      return {
        title: "実行中",
        text: "進行状況、件数、失敗要約を見ながら中断を判断できます。"
      }
    case "paused":
      return {
        title: "中断中",
        text: "同じ翻訳段階の処理を再開またはキャンセルできます。"
      }
    case "recoverable_failed":
      return {
        title: "再試行可能な失敗",
        text: "成功済み field を保持したまま、未処理または失敗 field を retry できます。"
      }
    case "validation_failed":
      return {
        title: "検証失敗",
        text: "保護要素検証に失敗したため、失敗訳文は保持せず再試行判断を待っています。"
      }
    case "empty_completed":
      return {
        title: "対象 0 件で完了",
        text: "provider を呼ばずに Completed となり、output readiness を確認できます。"
      }
    case "completed":
      return {
        title: "完了",
        text: "field result 整合と output readiness を確認し、後続出力へ進めます。"
      }
    case "canceled":
      return {
        title: "キャンセル済み",
        text: "途中成功結果は output readiness に使わず、再実行判断を待ちます。"
      }
    case "failed":
      return {
        title: "失敗",
        text: "retryable ではない失敗です。error kind と redacted summary を確認してください。"
      }
    case "not_ready":
      return {
        title: "開始条件未達",
        text: "NPC ペルソナ生成段階、ジョブ状態、参照 snapshot の条件が揃うまで操作を制限しています。"
      }
    default:
      return {
        title: "読み込み中",
        text: "翻訳段階の summary を更新しています。"
      }
  }
}

function buildProgressDetail(state: BodyTranslationPhaseScreenState): string {
  if (!state.summary) {
    return "summary 未取得です。"
  }

  const progress = state.summary.progress
  const providerTargetCount =
    state.summary.requestSummary?.providerTargetCount ?? progress.targetCount
  return `${progress.processedCount.toLocaleString("ja-JP")} / ${providerTargetCount.toLocaleString("ja-JP")} 件 / AI 送信対象 ${providerTargetCount.toLocaleString("ja-JP")} 件 / 成功 ${progress.translatedCount.toLocaleString("ja-JP")} 件 / スキップ ${progress.skippedCount.toLocaleString("ja-JP")} 件 / ${buildCurrentStepLabel(progress.currentStep)}`
}

function extractFieldResultItems(
  summary: BodyTranslationPhaseSummaryResponse | null
): BodyTranslationFieldResultItem[] {
  if (!summary) {
    return []
  }

  const candidates =
    summary.fieldResults ?? summary.resultSummary?.fieldResults ?? []

  return candidates.map((item) => {
    const identity = item.identity
    return {
      fieldId: readFieldIdentityLabel(item),
      fieldLabel:
        readOptionalString(identity?.fieldLabel) ??
        readOptionalString(item.fieldLabel) ??
        readOptionalString(identity?.editorId) ??
        readOptionalString(identity?.formId) ??
        "-",
      recordTypeLabel: readOptionalString(identity?.recordType) ?? "-",
      fieldTypeLabel: readOptionalString(identity?.fieldType) ?? "-",
      formIdLabel: readOptionalString(identity?.formId) ?? "-",
      editorIdLabel: readOptionalString(identity?.editorId) ?? "-",
      sourceExcerpt: readOptionalString(item.sourceExcerpt) ?? "-",
      translatedText: readOptionalString(item.translatedText) ?? "-",
      outputStatus: readOptionalString(item.outputStatus) ?? "-",
      protectionValidation:
        readOptionalString(item.protectionValidationResult) ??
        readOptionalString(item.protectionValidationSummary) ??
        "-",
      retryCountLabel: formatCountFromUnknown(item.retryCount),
      rawItem: item
    }
  })
}

function readFieldIdentityLabel(
  item: BodyTranslationPhaseGatewayFieldResultItem
): string {
  const identity = item.identity
  const preferredValues = [
    identity?.translationFieldId,
    item.fieldId,
    identity?.phaseTranslationFieldId
  ]

  for (const value of preferredValues) {
    if (typeof value === "string" && value.trim()) {
      return value
    }
    if (typeof value === "number") {
      return value.toLocaleString("ja-JP")
    }
  }

  return "-"
}

function readOptionalString(value: string | undefined): string | null {
  return typeof value === "string" && value.trim() ? value : null
}

function formatCountFromUnknown(value: unknown): string {
  return typeof value === "number" ? `${value.toLocaleString("ja-JP")} 回` : "-"
}

function buildProviderStateLabel(
  state: BodyTranslationPhaseScreenState
): string {
  const summary = state.summary
  if (!summary) {
    return "-"
  }

  const phaseState = normalizePhaseState(summary.phaseState)
  if (phaseState === "completed" && (summary.execution?.requestUnitCount ?? 0) === 0) {
    return "provider 未実行"
  }
  if (summary.errorSummary?.errorKind === "provider_failure") {
    return "provider failure"
  }
  if (summary.errorSummary?.errorKind === "late_response_rejected") {
    return "late response rejected"
  }
  if (phaseState === "running") {
    return "provider 実行中"
  }
  return (summary.execution?.requestUnitCount ?? 0) > 0 ? "provider 実行あり" : "-"
}

function buildLateResponseLabel(
  state: BodyTranslationPhaseScreenState
): string {
  return state.summary?.errorSummary?.errorKind === "late_response_rejected"
    ? "あり"
    : "なし"
}

function buildFieldResultAvailabilityLabel(
  state: BodyTranslationPhaseScreenState,
  fieldResultItems: BodyTranslationFieldResultItem[]
): string {
  if (fieldResultItems.length > 0) {
    return `${fieldResultItems.length.toLocaleString("ja-JP")} 件を表示中`
  }

  if (state.summary?.resultSummary) {
    return "public seam に field 単位一覧が未定義のため、件数要約のみ表示しています。"
  }

  return "summary 未取得です。"
}

function buildActionCards(
  state: BodyTranslationPhaseScreenState
): BodyTranslationPhaseScreenViewModel["actionCards"] {
  const isBusy = state.phase === "loading" || state.phase === "submitting"

  if (!state.summary) {
    return [
      { id: "start", label: "開始", disabled: true, blockedReason: "", tone: "primary" },
      { id: "pause", label: "中断", disabled: true, blockedReason: "", tone: "warning" },
      { id: "resume", label: "再開", disabled: true, blockedReason: "", tone: "default" },
      { id: "retry", label: "リトライ", disabled: true, blockedReason: "", tone: "default" },
      { id: "cancel", label: "キャンセル", disabled: true, blockedReason: "", tone: "warning" },
      { id: "check-output-readiness", label: "出力準備を確認", disabled: true, blockedReason: "", tone: "default" }
    ]
  }

  const derived = deriveBodyActionEnablement(state.summary.projection)

  return [
    {
      id: "start",
      label: "開始",
      disabled: isBusy || !derived.canStart,
      blockedReason: derived.startBlockedReason ?? "",
      tone: "primary"
    },
    {
      id: "pause",
      label: "中断",
      disabled: isBusy || !derived.canPause,
      blockedReason: derived.pauseBlockedReason ?? "",
      tone: "warning"
    },
    {
      id: "resume",
      label: "再開",
      disabled: isBusy || !derived.canResume,
      blockedReason: derived.resumeBlockedReason ?? "",
      tone: "default"
    },
    {
      id: "retry",
      label: "リトライ",
      disabled: isBusy || !derived.canRetry,
      blockedReason: derived.retryBlockedReason ?? "",
      tone: "default"
    },
    {
      id: "cancel",
      label: "キャンセル",
      disabled: isBusy || !derived.canCancel,
      blockedReason: derived.cancelBlockedReason ?? "",
      tone: "warning"
    },
    {
      id: "check-output-readiness",
      label: "出力準備を確認",
      disabled: isBusy || !deriveBodyOutputReadinessReady(state.summary),
      blockedReason: state.summary ? (deriveBodyOutputReadinessBlockedReason(state.summary) || "") : "",
      tone: state.summary && deriveBodyOutputReadinessReady(state.summary) ? "primary" : "default"
    }
  ]
}

function buildFailedCountLabel(
  summary: BodyTranslationPhaseSummaryResponse | null,
  resultSummary: NonNullable<
    BodyTranslationPhaseSummaryResponse["resultSummary"]
  > | null
): string {
  return formatCount(
    resultSummary?.failedCount ?? calculateDerivedFailedCount(summary)
  )
}

function calculateDerivedFailedCount(
  summary: BodyTranslationPhaseSummaryResponse | null
): number | undefined {
  if (
    typeof summary?.progress.processedCount !== "number" ||
    typeof summary.progress.translatedCount !== "number" ||
    typeof summary.progress.skippedCount !== "number"
  ) {
    return undefined
  }

  return Math.max(
    summary.progress.processedCount -
      summary.progress.translatedCount -
      summary.progress.skippedCount,
    0
  )
}

function buildScreenActionEnablement(
  summary: BodyTranslationPhaseSummaryResponse | null
): BodyTranslationPhaseScreenViewModel["screenActionEnablement"] {
  if (!summary) {
    return {
      canStart: false,
      canPause: false,
      canResume: false,
      canRetry: false,
      canCancel: false,
      canCheckOutputReadiness: false
    }
  }
  const derived = deriveBodyActionEnablement(summary.projection)
  return {
    canStart: derived.canStart,
    canPause: derived.canPause,
    canResume: derived.canResume,
    canRetry: derived.canRetry,
    canCancel: derived.canCancel,
    canCheckOutputReadiness: deriveBodyOutputReadinessReady(summary)
  }
}

function buildPhaseRunRequestShape(
  summary: BodyTranslationPhaseSummaryResponse | null
): { jobId: number; phaseRunId: number } | undefined {
  return summary && typeof summary.phaseRunId === "number"
    ? { jobId: summary.jobId, phaseRunId: summary.phaseRunId }
    : undefined
}

export class BodyTranslationPhasePresenter {
  toViewModel(
    state: BodyTranslationPhaseScreenState,
    isGatewayConnected: boolean,
    availableProviders: { value: string; label: string }[] = [],
    availableModels: { value: string; label: string }[] = []
  ): BodyTranslationPhaseScreenViewModel {
    const summary = state.summary
    const resultSummary = summary?.resultSummary ?? null
    const errorSummary = summary?.errorSummary ?? null
    const outputReadiness = getEffectiveReadiness(state)
    const fieldResultItems = extractFieldResultItems(summary)
    const viewState = buildViewState(state)
    const statusCopy = buildStatusCopy(state, viewState)
    const phaseRunRequestShape = buildPhaseRunRequestShape(summary)

    return {
      ...state,
      gatewayStatus: isGatewayConnected ? "接続準備済み" : "未接続",
      viewState,
      isLoading: state.phase === "loading" && !state.hasLoaded,
      isRefreshing: state.phase === "loading" && state.hasLoaded,
      isSubmitting: state.phase === "submitting",
      hasJobSelection: state.jobId !== null,
      currentPhaseLabel: buildCurrentPhaseLabel(summary?.currentPhase),
      phaseStateLabel: buildPhaseStateLabel(summary?.phaseState),
      statusTitle: statusCopy.title,
      statusText: statusCopy.text,
      progressPercent: summary?.progress.percent ?? 0,
      progressLabel: summary ? `${summary.progress.percent}%` : "-",
      progressDetail: buildProgressDetail(state),
      startedAtLabel: formatDate(summary?.startedAt),
      finishedAtLabel: formatDate(summary?.finishedAt),
      targetCountLabel: formatCount(summary?.progress.targetCount),
      processedCountLabel: formatCount(summary?.progress.processedCount),
      translatedCountLabel: formatCount(summary?.progress.translatedCount),
      failedCountLabel: buildFailedCountLabel(summary, resultSummary),
      skippedCountLabel: formatCount(summary?.progress.skippedCount),
      dictionaryDigestLabel: summary?.inputSummary.dictionaryDigest ?? "-",
      personaDigestLabel: summary?.inputSummary.personaDigest ?? "-",
      metadataDigestLabel: summary?.inputSummary.metadataDigest ?? "-",
      promptDigestLabel: summary?.inputSummary.promptDigest ?? "-",
      inputSnapshotRefLabel: summary?.inputSummary.inputSnapshotRef ?? "-",
      skippedReasonsLabel: formatReasonList(
        summary?.inputSummary.skippedReasons
      ),
      isExecutionConfigured: summary ? isBodyExecutionConfigured(summary) : false,
      providerLabel:
        summary?.aiSettings?.provider ??
        summary?.execution?.provider ??
        "設定未完了",
      modelLabel:
        summary?.aiSettings?.model ??
        summary?.execution?.model ??
        "設定未完了",
      executionModeLabel:
        summary?.aiSettings?.executionMode ??
        summary?.execution?.executionMode ??
        "設定未完了",
      credentialRefLabel: summary?.execution?.credentialRef ?? "設定未完了",
      providerOptions: buildBodyProviderOptions(summary, availableProviders),
      modelOptions: buildBodyModelOptions(summary, availableModels),
      executionOptions: BODY_EXECUTION_OPTIONS,
      providerTargetCountLabel: formatCount(
        summary?.requestSummary?.providerTargetCount
      ),
      exactDictionaryExclusionCountLabel: formatCount(
        summary?.requestSummary?.exactDictionaryExclusionCount
      ),
      partialDictionaryConstraintCountLabel: formatCount(
        summary?.requestSummary?.partialDictionaryConstraintCount
      ),
      requestUnitCountLabel: formatCount(summary?.execution?.requestUnitCount),
      outputCountLabel: formatCount(summary?.execution?.outputCount),
      resultOutputCountLabel: formatCount(resultSummary?.outputCount),
      providerStateLabel: buildProviderStateLabel(state),
      lateResponseLabel: buildLateResponseLabel(state),
      outputReadinessLabel:
        outputReadiness?.ready === true ? "ready" : "not ready",
      outputReadinessBlockedReason: outputReadiness?.blockedReason ?? "",
      outputReadinessCompletedFieldCountLabel: formatCount(
        outputReadiness?.completedFieldCount
      ),
      outputReadinessStatusLabel:
        outputReadiness?.statusConsistent === true ? "整合" : "不整合",
      errorKindLabel:
        errorSummary?.errorKind ?? "-",
      errorReasonLabel: errorSummary?.reason ?? "-",
      retryableLabel:
        errorSummary === null
          ? "-"
          : errorSummary.retryable
            ? "再試行可能"
            : "再試行不可",
      fieldResultAvailabilityLabel: buildFieldResultAvailabilityLabel(
        state,
        fieldResultItems
      ),
      fieldResultItems,
      actionCards: buildActionCards(state),
      screenActionEnablement: buildScreenActionEnablement(summary),
      lastErrorSummary: errorSummary,
      projection: summary?.projection ?? null,
      latestProgressSummary: summary?.progress ?? null,
      latestInputSummary: summary?.inputSummary ?? null,
      latestRequestSummary: summary?.requestSummary ?? null,
      latestExecutionSummary: summary?.execution ?? null,
      latestResultSummary: resultSummary,
      latestErrorKind:
        errorSummary?.errorKind ?? null,
      latestOutputReadiness: outputReadiness,
      pauseRequestShape: phaseRunRequestShape,
      resumeRequestShape: phaseRunRequestShape,
      retryRequestShape: phaseRunRequestShape,
      cancelRequestShape: phaseRunRequestShape,
      outputReadinessRequestShape:
        summary === null ? undefined : { jobId: summary.jobId }
    }
  }
}
