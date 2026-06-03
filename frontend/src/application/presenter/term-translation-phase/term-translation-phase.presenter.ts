import type {
  ProcessingTargetListPageState
} from "@application/gateway-contract/processing-target"

import type {
  TermTranslationExecutionConfigSummary,
  TermTranslationNextPhaseReadinessResponse,
  TermTranslationPhaseProjection,
  TermTranslationPhaseErrorSummary,
  TermTranslationPhaseErrorKind,
  TermTranslationPhaseProgressSummary,
  TermTranslationPhaseResultSummary,
  TermTranslationPhaseSummaryResponse
} from "@application/gateway-contract/term-translation-phase"

type TermTranslationPhaseActionKind =
  | "start"
  | "pause"
  | "resume"
  | "retry"

type TermTranslationPhaseViewState =
  | "loading"
  | "idle_ready"
  | "running"
  | "empty_completed"
  | "completed"
  | "paused"
  | "recoverable_failed"
  | "blocked"

interface TermTranslationPhaseScreenState {
  jobId: number | null
  phase: "idle" | "loading" | "ready" | "submitting"
  summary: TermTranslationPhaseSummaryResponse | null
  nextPhaseReadiness: TermTranslationNextPhaseReadinessResponse | null
  errorMessage: string
  pendingAction: TermTranslationPhaseActionKind | null
  hasLoaded: boolean
  initialFetchDone: boolean
  processingTargetPageState?: ProcessingTargetListPageState | null
}

interface TermTranslationPhaseActionCard {
  id: TermTranslationPhaseActionKind | "next-phase"
  label: string
  disabled: boolean
  blockedReason: string
  tone: "default" | "primary" | "warning"
}

interface SelectOption {
  value: string
  label: string
}

interface TermTranslationPhaseScreenViewModel extends TermTranslationPhaseScreenState {
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
  isExecutionConfigured: boolean
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
  providerOptions: SelectOption[]
  modelOptions: SelectOption[]
  executionOptions: SelectOption[]
  actionCards: TermTranslationPhaseActionCard[]
  lastErrorSummary: TermTranslationPhaseErrorSummary | null
  projection: TermTranslationPhaseProjection | null
  latestProgressSummary: TermTranslationPhaseProgressSummary | null
  latestResultSummary: TermTranslationPhaseResultSummary | null
  latestExecutionSummary: TermTranslationExecutionConfigSummary | null
  latestErrorKind: TermTranslationPhaseErrorKind | null
}

const PHASE_STATE_LABELS: Record<string, string> = {
  blocked: "開始待ち",
  completed: "完了",
  failed: "失敗",
  idle: "待機",
  idle_ready: "開始待ち",
  paused: "中断",
  pending: "開始待ち",
  ready: "開始可能",
  recoverable_failed: "再試行可能な失敗",
  retryable_failed: "再試行可能な失敗",
  running: "実行中"
}

const CURRENT_PHASE_LABELS: Record<string, string> = {
  term_translation: "単語翻訳"
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

function normalizePhaseState(phaseState: string | undefined): string {
  return phaseState?.trim().toLowerCase().replaceAll(" ", "_") ?? ""
}

function buildViewState(
  state: TermTranslationPhaseScreenState
): TermTranslationPhaseViewState {
  if (state.phase === "loading" && !state.summary) {
    return "loading"
  }

  const summary = state.summary
  if (!summary) {
    return state.jobId === null ? "blocked" : "loading"
  }

  const normalizedState = normalizePhaseState(summary.phaseState)
  if (normalizedState === "pending" || normalizedState === "idle_ready") {
    return "idle_ready"
  }

  if (normalizedState === "paused") {
    return "paused"
  }

  if (
    normalizedState === "recoverable_failed" ||
    normalizedState === "retryable_failed" ||
    summary.errorSummary?.retryable
  ) {
    return "recoverable_failed"
  }

  if (
    normalizedState === "running" ||
    normalizedState === "in_progress" ||
    normalizedState === "processing"
  ) {
    return "running"
  }

  if (
    normalizedState === "completed" ||
    normalizedState === "succeeded" ||
    normalizedState === "done"
  ) {
    return summary.aiTargetCount === 0 ? "empty_completed" : "completed"
  }

  if (isExecutionConfigured(summary) && !summary.phaseRunId) {
    return "idle_ready"
  }

  return "blocked"
}

function buildPhaseStateLabel(phaseState: string | undefined): string {
  if (!phaseState) {
    return "未取得"
  }

  const normalized = normalizePhaseState(phaseState)
  return PHASE_STATE_LABELS[normalized] ?? phaseState
}

function buildCurrentPhaseLabel(currentPhase: string | undefined): string {
  const normalized = normalizePhaseState(currentPhase)
  return CURRENT_PHASE_LABELS[normalized] ?? (currentPhase || "未開始")
}

const CURRENT_STEP_LABELS: Record<string, string> = {
  completed: "完了",
  idle: "待機",
  idle_ready: "開始待ち",
  pending: "開始待ち",
  provider_request: "AI 処理中",
  ready: "開始可能",
  retry: "再試行中",
  running: "実行中",
  starting: "開始中"
}

function buildCurrentStepLabel(currentStep: string | undefined): string {
  if (!currentStep) {
    return "-"
  }

  const normalized = normalizePhaseState(currentStep)
  return CURRENT_STEP_LABELS[normalized] ?? currentStep
}

function buildStatusCopy(
  state: TermTranslationPhaseScreenState,
  viewState: TermTranslationPhaseViewState
): { title: string; text: string } {
  if (state.jobId === null) {
    return {
      title: "ジョブ未選択",
      text: "ジョブIDを指定すると、単語翻訳段階の状態を確認できます。"
    }
  }

  if (state.phase === "loading" && !state.summary) {
    return {
      title: "summary を取得中",
      text: "単語翻訳段階の現在状態、進行状況、結果 summary を読み込んでいます。"
    }
  }

  const errorMessage = state.errorMessage.trim()
  if (errorMessage) {
    return {
      title: "更新に注意が必要",
      text: errorMessage
    }
  }

  switch (viewState) {
    case "idle_ready":
      return {
        title: "開始可能",
        text: "Ready ジョブで、実行中の処理がないため、単語翻訳段階を開始できます。"
      }
    case "running":
      return {
        title: "進行中",
        text: "進行状況と現在の作業を確認しながら、必要に応じて中断できます。"
      }
    case "empty_completed":
      return {
        title: "対象語なしで完了",
        text: "共通辞書で対象語が解決したため、provider 実行なしで Completed です。"
      }
    case "completed":
      return {
        title: "完了",
        text: "確定訳語とジョブ内辞書反映を確認して、次の翻訳段階へ進めるか判断できます。"
      }
    case "paused":
      return {
        title: "中断中",
        text: "同じ翻訳段階の処理を再開できます。進行状況は保持されています。"
      }
    case "recoverable_failed":
      return {
        title: "再試行可能な失敗",
        text: "error kind と短い理由を確認し、未処理 term だけを retry できます。"
      }
    case "blocked":
      return {
        title: "開始条件未達",
        text: "ジョブ状態または辞書参照条件を満たすまで操作を制限しています。"
      }
    default:
      return {
        title: "読み込み中",
        text: "翻訳段階の summary を更新しています。"
      }
  }
}

function buildProgressLabel(state: TermTranslationPhaseScreenState): string {
  if (!state.summary) {
    return "-"
  }

  return `${state.summary.progress.percent}%`
}

function buildProgressDetail(state: TermTranslationPhaseScreenState): string {
  if (!state.summary) {
    return "summary 未取得です。"
  }

  const progress = state.summary.progress
  return `${progress.processedCount.toLocaleString("ja-JP")} / ${progress.aiTargetCount.toLocaleString("ja-JP")} 件 / AI 翻訳対象語 ${progress.aiTargetCount.toLocaleString("ja-JP")} 件 / ${buildCurrentStepLabel(progress.currentStep)}`
}

function buildProviderSkippedLabel(
  state: TermTranslationPhaseScreenState
): string {
  if (!state.summary?.resultSummary) {
    return "-"
  }

  return state.summary.resultSummary.providerSkipped
    ? "provider 未実行"
    : "provider 実行あり"
}

function buildSnapshotLabel(state: TermTranslationPhaseScreenState): string {
  const execution = state.summary?.execution
  if (!execution) {
    return "-"
  }

  if (execution.snapshotVersion) {
    return execution.snapshotVersion
  }

  if (execution.snapshotDigest) {
    return execution.snapshotDigest
  }

  return "-"
}

// ドメイン状態射影 field に対する enum 集合定数（design-diff H 節の略号と対応）
const TERMINAL_JOB = new Set(["completed", "failed", "canceled"])
const RUNNING_PHASE = new Set(["running", "in_progress", "processing"])
const IDLE_READY_PHASE = new Set(["pending", "idle_ready", "ready", ""])
const PAUSED_PHASE = new Set(["paused"])
const RECOVERABLE_FAILED_PHASE = new Set(["recoverable_failed", "retryable_failed"])
const COMPLETED_PHASE = new Set(["completed", "succeeded", "done"])

function isExecutionConfigured(summary: TermTranslationPhaseSummaryResponse): boolean {
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

const EXECUTION_OPTIONS: SelectOption[] = [
  { value: "batch", label: "バッチ処理" },
  { value: "sync", label: "逐次処理" }
]

function buildProviderOptions(
  summary: TermTranslationPhaseSummaryResponse | null,
  availableProviders: SelectOption[]
): SelectOption[] {
  if (availableProviders.length > 0) {
    const currentProvider = summary?.aiSettings?.provider ?? summary?.execution?.provider
    if (currentProvider && !availableProviders.some((p) => p.value === currentProvider)) {
      return [{ value: "", label: "選んでください" }, ...availableProviders]
    }
    if (!currentProvider) {
      return [{ value: "", label: "選んでください" }, ...availableProviders]
    }
    return availableProviders
  }
  const fallbackLabel = summary?.aiSettings?.provider ?? summary?.execution?.provider ?? "設定未完了"
  return [{ value: fallbackLabel, label: fallbackLabel }]
}

function buildModelOptions(
  summary: TermTranslationPhaseSummaryResponse | null,
  availableModels: SelectOption[]
): SelectOption[] {
  if (availableModels.length > 0) {
    return availableModels
  }
  const fallbackLabel = summary?.aiSettings?.model ?? summary?.execution?.model
  if (fallbackLabel) {
    return [{ value: fallbackLabel, label: fallbackLabel }]
  }
  return [{ value: "", label: "モデル一覧を更新してください" }]
}

function deriveCanStartNextPhase(
  projection: TermTranslationPhaseProjection
): boolean {
  // H-5: not terminal ∧ phaseLifecycle ∈ COMPLETED_PHASE ∧ confirmedCount ≥ aiTargetCount
  if (TERMINAL_JOB.has(projection.jobLifecycle)) {
    return false
  }
  if (!COMPLETED_PHASE.has(projection.phaseLifecycle)) {
    return false
  }
  return projection.confirmedCount >= projection.aiTargetCount
}

function deriveNextPhaseBlockedReason(
  projection: TermTranslationPhaseProjection
): string {
  // H-5 BlockedReason 決定論理（上から評価）
  if (TERMINAL_JOB.has(projection.jobLifecycle)) {
    return "ジョブが終端状態のため次段階を開始できません。"
  }
  if (!COMPLETED_PHASE.has(projection.phaseLifecycle)) {
    return "単語翻訳段階が未完了のため次段階を開始できません。"
  }
  if (projection.confirmedCount < projection.aiTargetCount) {
    return "確定件数が AI 対象件数に達していないため次段階を開始できません。"
  }
  return ""
}

interface TermActionEnablementResult {
  canStart: boolean
  startBlockedReason?: string
  canPause: boolean
  pauseBlockedReason?: string
  canResume: boolean
  resumeBlockedReason?: string
  canRetry: boolean
  retryBlockedReason?: string
}

function deriveTermActionEnablement(
  projection: TermTranslationPhaseProjection
): TermActionEnablementResult {
  // H-1〜H-4: ドメイン状態射影だけを入力とする決定論的導出
  const terminal = TERMINAL_JOB.has(projection.jobLifecycle)
  const running = RUNNING_PHASE.has(projection.phaseLifecycle)
  const paused = PAUSED_PHASE.has(projection.phaseLifecycle)
  const recoverableFailed =
    RECOVERABLE_FAILED_PHASE.has(projection.phaseLifecycle) ||
    projection.errorKind === "recoverable"
  const idleReady = IDLE_READY_PHASE.has(projection.phaseLifecycle)
  const hasProcessingTarget = projection.aiTargetCount > 0

  // H-1: canStart = not terminal ∧ phaseLifecycle ∉ RUNNING_PHASE ∧ phaseLifecycle ∈ IDLE_READY_PHASE
  //               ∧ aiSettingsConfigured = true ∧ aiTargetCount > 0
  const canStart =
    !terminal &&
    !running &&
    idleReady &&
    projection.aiSettingsConfigured &&
    hasProcessingTarget

  // H-1 BlockedReason（上から評価）
  const startBlockedReason = terminal
    ? "ジョブが終端状態のため開始できません。"
    : running
      ? "実行中の翻訳段階があるため開始できません。"
      : !projection.aiSettingsConfigured
        ? "実行設定が未構成のため開始できません。"
        : !hasProcessingTarget
          ? "処理対象が 0 件のため開始できません。"
          : !idleReady
            ? "ジョブが開始可能状態ではありません。"
            : undefined

  // H-2: canPause = not terminal ∧ phaseLifecycle ∈ RUNNING_PHASE
  const canPause = !terminal && running
  const pauseBlockedReason = terminal
    ? "ジョブが終端状態のため中断できません。"
    : !running
      ? "フェーズが実行中ではありません。"
      : undefined

  // H-3: canResume = not terminal ∧ (phaseLifecycle ∈ PAUSED_PHASE ∨ RECOVERABLE_FAILED_PHASE ∨ errorKind = recoverable)
  const canResume = !terminal && (paused || recoverableFailed)
  const resumeBlockedReason = terminal
    ? "ジョブが終端状態のため再開できません。"
    : !paused && !recoverableFailed
      ? "フェーズが再開可能な状態ではありません。"
      : undefined

  // H-4: canRetry = not terminal ∧ (phaseLifecycle ∈ RECOVERABLE_FAILED_PHASE ∨ errorKind = recoverable)
  const canRetry = !terminal && recoverableFailed
  const retryBlockedReason = terminal
    ? "ジョブが終端状態のため再試行できません。"
    : !recoverableFailed
      ? "フェーズが再試行可能な状態ではありません。"
      : undefined

  return {
    canStart,
    startBlockedReason,
    canPause,
    pauseBlockedReason,
    canResume,
    resumeBlockedReason,
    canRetry,
    retryBlockedReason
  }
}

function buildActionCards(
  state: TermTranslationPhaseScreenState
): TermTranslationPhaseActionCard[] {
  const isBusy = state.phase === "loading" || state.phase === "submitting"

  if (!state.summary) {
    return [
      { id: "start", label: "開始", disabled: true, blockedReason: "", tone: "primary" },
      { id: "pause", label: "中断", disabled: true, blockedReason: "", tone: "warning" },
      { id: "resume", label: "再開", disabled: true, blockedReason: "", tone: "default" },
      { id: "retry", label: "リトライ", disabled: true, blockedReason: "", tone: "default" },
      { id: "next-phase", label: "次の翻訳段階へ進む", disabled: true, blockedReason: "", tone: "primary" }
    ]
  }

  const projection = state.summary.projection
  const derived = deriveTermActionEnablement(projection)

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
      id: "next-phase",
      label: "次の翻訳段階へ進む",
      disabled: isBusy || !deriveCanStartNextPhase(projection),
      blockedReason: deriveNextPhaseBlockedReason(projection) || "",
      tone: "primary"
    }
  ]
}

export class TermTranslationPhasePresenter {
  toViewModel(
    state: TermTranslationPhaseScreenState,
    isGatewayConnected: boolean,
    availableProviders: SelectOption[] = [],
    availableModels: SelectOption[] = []
  ): TermTranslationPhaseScreenViewModel {
    const viewState = buildViewState(state)
    const statusCopy = buildStatusCopy(state, viewState)
    const summary = state.summary
    const resultSummary = summary?.resultSummary
    const errorSummary = summary?.errorSummary ?? null
    const providerLabel =
      summary?.aiSettings?.provider ??
      summary?.execution?.provider ??
      "設定未完了"
    const modelLabel =
      summary?.aiSettings?.model ??
      summary?.execution?.model ??
      "設定未完了"
    const executionModeLabel =
      summary?.aiSettings?.executionMode ??
      summary?.execution?.executionMode ??
      "設定未完了"
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
      progressLabel: buildProgressLabel(state),
      progressDetail: buildProgressDetail(state),
      startedAtLabel: formatDate(summary?.startedAt),
      finishedAtLabel: formatDate(summary?.finishedAt),
      totalTermCountLabel: formatCount(summary?.totalTermCount),
      dictionaryHitCountLabel: formatCount(summary?.dictionaryHitCount),
      aiTargetCountLabel: formatCount(summary?.aiTargetCount),
      confirmedCountLabel: formatCount(resultSummary?.confirmedCount),
      jobDictionaryAppliedCountLabel: formatCount(
        resultSummary?.jobDictionaryAppliedCount
      ),
      replacementTargetCountLabel: formatCount(
        resultSummary?.replacementTargetCount
      ),
      unmatchedCountLabel: formatCount(resultSummary?.unmatchedCount),
      isExecutionConfigured: summary ? isExecutionConfigured(summary) : false,
      providerLabel,
      modelLabel,
      executionModeLabel,
      credentialRefLabel: summary?.execution?.credentialRef ?? "設定未完了",
      snapshotLabel: buildSnapshotLabel(state),
      errorKindLabel: errorSummary?.errorKind ?? "-",
      errorReasonLabel: errorSummary?.reason ?? "-",
      retryableLabel:
        errorSummary === null
          ? "-"
          : errorSummary.retryable
            ? "再試行可能"
            : "再試行不可",
      nextPhaseStatusLabel:
        (summary ? deriveCanStartNextPhase(summary.projection) : false)
          ? "開始可能"
          : "開始不可",
      nextPhaseBlockedReason:
        summary ? (deriveNextPhaseBlockedReason(summary.projection) || "") : "",
      providerSkippedLabel: buildProviderSkippedLabel(state),
      providerOptions: buildProviderOptions(summary, availableProviders),
      modelOptions: buildModelOptions(summary, availableModels),
      executionOptions: EXECUTION_OPTIONS,
      actionCards: buildActionCards(state),
      lastErrorSummary: errorSummary,
      projection: summary?.projection ?? null,
      latestProgressSummary: summary?.progress ?? null,
      latestResultSummary: resultSummary ?? null,
      latestExecutionSummary: summary?.execution ?? null,
      latestErrorKind: errorSummary?.errorKind ?? null
    }
  }
}
