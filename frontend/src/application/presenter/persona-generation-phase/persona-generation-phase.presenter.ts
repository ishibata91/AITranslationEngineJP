import type { ProcessingTargetListPageState } from "@application/gateway-contract/processing-target"

import type {
  PersonaGenerationBodyReadinessInputSummary,
  PersonaGenerationBodyReadinessResponse,
  PersonaGenerationExecutionSummary,
  PersonaGenerationPhaseActionEnablement,
  PersonaGenerationPhaseErrorKind,
  PersonaGenerationPhaseErrorSummary,
  PersonaGenerationPhaseProgressSummary,
  PersonaGenerationPhaseResultSummary,
  PersonaGenerationPhaseSummaryResponse,
  PersonaGenerationTargetSummary
} from "@application/gateway-contract/persona-generation-phase"

type PersonaGenerationPhaseActionKind =
  | "start"
  | "pause"
  | "resume"
  | "retry"
  | "cancel"
  | "check-body-readiness"
  | "start-body-phase"

type PersonaGenerationPhaseViewState =
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

interface PersonaGenerationPhaseScreenState {
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

interface PersonaGenerationPhaseActionCard {
  id: PersonaGenerationPhaseActionKind
  label: string
  disabled: boolean
  blockedReason: string
  tone: "default" | "primary" | "warning"
}

interface PersonaGenerationPhaseScreenActionEnablement {
  canStart: boolean
  canPause: boolean
  canResume: boolean
  canRetry: boolean
  canCancel: boolean
  canCheckBodyReadiness: boolean
  canStartBodyPhase: boolean
}

interface PersonaGenerationPhaseScreenViewModel extends PersonaGenerationPhaseScreenState {
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
  providerOptions: PersonaSelectOption[]
  modelOptions: PersonaSelectOption[]
  executionOptions: PersonaSelectOption[]
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
}

const PHASE_STATE_LABELS: Record<string, string> = {
  blocked: "開始待ち",
  completed: "完了",
  failed: "失敗",
  idle: "未開始",
  not_started: "未開始",
  paused: "中断",
  pending: "開始待ち",
  ready: "開始可能",
  recoverable_failed: "再試行可能な失敗",
  retryable_failed: "再試行可能な失敗",
  running: "実行中",
  snapshot_missing: "snapshot 不足"
}

const PHASE_VIEW_STATE_BY_PHASE_STATE: Partial<
  Record<string, PersonaGenerationPhaseViewState>
> = {
  failed: "failed",
  in_progress: "running",
  paused: "paused",
  pending: "not_started",
  processing: "running",
  running: "running"
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

function formatList(values: string[] | undefined): string {
  if (!values || values.length === 0) {
    return "-"
  }

  return values.join(" / ")
}

function normalizePhaseState(phaseState: string | undefined): string {
  return phaseState?.trim().toLowerCase().replaceAll(" ", "_") ?? ""
}

function isSnapshotMissing(state: PersonaGenerationPhaseScreenState): boolean {
  const snapshotReferenceStatus =
    state.summary?.resultSummary?.snapshotReferenceStatus
      ?.trim()
      .toLowerCase() ?? ""

  return (
    state.summary?.errorSummary?.errorKind === "snapshot_missing" ||
    snapshotReferenceStatus.includes("missing")
  )
}

function isRecoverableFailureState(
  normalizedState: string,
  summary: PersonaGenerationPhaseSummaryResponse
): boolean {
  return (
    normalizedState === "recoverable_failed" ||
    normalizedState === "retryable_failed" ||
    summary.errorSummary?.retryable === true
  )
}

function isCompletedState(normalizedState: string): boolean {
  return (
    normalizedState === "completed" ||
    normalizedState === "succeeded" ||
    normalizedState === "done"
  )
}

function buildSummaryViewState(
  summary: PersonaGenerationPhaseSummaryResponse
): PersonaGenerationPhaseViewState {
  const normalizedState = normalizePhaseState(summary.phaseState)

  if (isRecoverableFailureState(normalizedState, summary)) {
    return "recoverable_failed"
  }

  const mappedState = PHASE_VIEW_STATE_BY_PHASE_STATE[normalizedState]
  if (mappedState) {
    return mappedState
  }

  if (isCompletedState(normalizedState)) {
    return summary.targetSummary.targetCount === 0
      ? "empty_completed"
      : "completed"
  }

  if (!summary.phaseRunId && summary.actionEnablement.canStart) {
    return "not_started"
  }

  return "blocked"
}

function buildViewState(
  state: PersonaGenerationPhaseScreenState
): PersonaGenerationPhaseViewState {
  if (state.phase === "loading" && !state.summary) {
    return "loading"
  }

  if (isSnapshotMissing(state)) {
    return "snapshot_missing"
  }

  const summary = state.summary
  if (!summary) {
    return state.jobId === null ? "blocked" : "loading"
  }

  return buildSummaryViewState(summary)
}

function buildPhaseStateLabel(phaseState: string | undefined): string {
  if (!phaseState) {
    return "未取得"
  }

  const normalized = normalizePhaseState(phaseState)
  return PHASE_STATE_LABELS[normalized] ?? phaseState
}

const CURRENT_STEP_LABELS: Record<string, string> = {
  completed: "完了",
  generating: "生成中",
  not_started: "未開始",
  paused: "中断",
  pending: "開始待ち",
  provider_request: "AI 処理中",
  rejected: "開始不可",
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
  state: PersonaGenerationPhaseScreenState,
  viewState: PersonaGenerationPhaseViewState
): { title: string; text: string } {
  if (state.jobId === null) {
    return {
      title: "ジョブ未選択",
      text: "ジョブIDを指定すると、NPC ペルソナ生成段階の状態を確認できます。"
    }
  }

  if (state.phase === "loading" && !state.summary) {
    return {
      title: "summary を取得中",
      text: "NPC ペルソナ生成段階の現在状態、進行状況、結果を読み込んでいます。"
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
    case "not_started":
      return {
        title: "開始可能",
        text: "単語翻訳段階の完了後なら、NPC ペルソナ生成を開始できます。"
      }
    case "running":
      return {
        title: "進行中",
        text: "target summary と生成件数を見ながら、中断または完了待ちを判断できます。"
      }
    case "paused":
      return {
        title: "中断中",
        text: "同じ翻訳段階の処理を再開できます。進行状況と対象 snapshot は維持されます。"
      }
    case "recoverable_failed":
      return {
        title: "再試行可能な失敗",
        text: "成功分は維持され、未処理 NPC だけを retry 対象として扱います。"
      }
    case "failed":
      return {
        title: "失敗",
        text: "error kind と短い理由を確認し、回復不能な失敗として扱います。"
      }
    case "completed":
      return {
        title: "完了",
        text: "snapshot 参照状態と本文翻訳の開始可否を確認して、次の翻訳段階へ進めるか判断できます。"
      }
    case "empty_completed":
      return {
        title: "対象 0 件で完了",
        text: "provider 未実行の empty completed として、snapshot 空の結果を表示しています。"
      }
    case "snapshot_missing":
      return {
        title: "snapshot 参照が不足",
        text: "snapshot missing のため、本文翻訳の開始可否と次の翻訳段階への開始を制限しています。"
      }
    case "blocked":
      return {
        title: "開始条件未達",
        text: "単語翻訳段階の完了、終端ジョブではないこと、実行中の翻訳段階がないことを満たすまで操作を制限しています。"
      }
    default:
      return {
        title: "読み込み中",
        text: "翻訳段階の summary を更新しています。"
      }
  }
}

function buildProgressLabel(state: PersonaGenerationPhaseScreenState): string {
  if (!state.summary) {
    return "-"
  }

  return `${state.summary.progress.percent}%`
}

function buildProgressDetail(state: PersonaGenerationPhaseScreenState): string {
  if (!state.summary) {
    return "summary 未取得です。"
  }

  const progress = state.summary.progress
  return `${progress.processedCount.toLocaleString("ja-JP")} / ${progress.totalCount.toLocaleString("ja-JP")} 件 / 対象 ${progress.targetCount.toLocaleString("ja-JP")} 件 / ${buildCurrentStepLabel(progress.currentStep)}`
}

function buildTargetSnapshotLabel(
  state: PersonaGenerationPhaseScreenState
): string {
  const summary = state.summary?.targetSummary
  if (!summary) {
    return "-"
  }

  return summary.targetSnapshotId || summary.targetSnapshotDigest || "-"
}

function buildSnapshotLabel(state: PersonaGenerationPhaseScreenState): string {
  const result = state.summary?.resultSummary
  if (!result) {
    return "-"
  }

  return result.snapshotId || result.snapshotDigest || "-"
}


function buildBodyReadinessInputSummaryLabel(
  state: PersonaGenerationPhaseScreenState
): string {
  const inputSummary = state.bodyReadiness?.inputSummary
  if (!inputSummary) {
    return "-"
  }

  return `${inputSummary.personaCount.toLocaleString("ja-JP")} 件 / missing ${inputSummary.missingCount.toLocaleString("ja-JP")} 件 / ${inputSummary.snapshotId || inputSummary.snapshotDigest || "-"}`
}

function isPersonaExecutionConfigured(
  summary: PersonaGenerationPhaseSummaryResponse
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

interface PersonaSelectOption {
  value: string
  label: string
}

const PERSONA_EXECUTION_OPTIONS: PersonaSelectOption[] = [
  { value: "batch", label: "バッチ処理" },
  { value: "sync", label: "逐次処理" }
]

function buildPersonaProviderOptions(
  summary: PersonaGenerationPhaseSummaryResponse | null,
  availableProviders: PersonaSelectOption[]
): PersonaSelectOption[] {
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

function buildPersonaModelOptions(
  summary: PersonaGenerationPhaseSummaryResponse | null,
  availableModels: PersonaSelectOption[]
): PersonaSelectOption[] {
  if (availableModels.length > 0) {
    return availableModels
  }
  const fallbackLabel = summary?.aiSettings?.model ?? summary?.execution?.model
  if (fallbackLabel) {
    return [{ value: fallbackLabel, label: fallbackLabel }]
  }
  return [{ value: "", label: "モデル一覧を更新してください" }]
}

function isPersonaTerminalJob(
  summary: PersonaGenerationPhaseSummaryResponse
): boolean {
  const startReason = summary.actionEnablement?.startBlockedReason ?? ""
  return startReason === "terminal_job"
}

function derivePersonaCanStartBodyPhase(
  summary: PersonaGenerationPhaseSummaryResponse
): boolean {
  if (isPersonaTerminalJob(summary)) {
    return false
  }
  const bodyReadiness = summary.resultSummary?.bodyReadiness
  return bodyReadiness === true
}

function derivePersonaBodyReadinessBlockedReason(
  summary: PersonaGenerationPhaseSummaryResponse
): string {
  if (isPersonaTerminalJob(summary)) {
    return "ジョブが終端状態のため本文翻訳を開始できません。"
  }
  const bodyReadiness = summary.resultSummary?.bodyReadiness
  if (!bodyReadiness) {
    return "ペルソナ snapshot 参照が準備できていません。"
  }
  return ""
}

function derivePersonaActionEnablement(
  summary: PersonaGenerationPhaseSummaryResponse
): PersonaGenerationPhaseActionEnablement {
  const normalizedState = normalizePhaseState(summary.phaseState)
  const terminal = isPersonaTerminalJob(summary)

  const isRunning =
    normalizedState === "running" ||
    normalizedState === "in_progress" ||
    normalizedState === "processing"
  const isPaused = normalizedState === "paused"
  const isRecoverableFailed =
    normalizedState === "recoverable_failed" ||
    normalizedState === "retryable_failed" ||
    summary.errorSummary?.retryable === true
  const isActive = isRunning || isPaused || isRecoverableFailed

  const canStart = !terminal && !isActive
  const startBlockedReason = terminal
    ? "ジョブが終端状態のため開始できません。"
    : isActive
      ? "実行中の翻訳段階があるため開始できません。"
      : undefined

  const canPause = !terminal && isRunning
  const pauseBlockedReason = terminal
    ? "ジョブが終端状態のため中断できません。"
    : !isRunning
      ? "フェーズが実行中ではありません。"
      : undefined

  const canResume = !terminal && (isPaused || isRecoverableFailed)
  const resumeBlockedReason = terminal
    ? "ジョブが終端状態のため再開できません。"
    : !isPaused && !isRecoverableFailed
      ? "フェーズが再開可能な状態ではありません。"
      : undefined

  const canRetry = !terminal && isRecoverableFailed
  const retryBlockedReason = terminal
    ? "ジョブが終端状態のため再試行できません。"
    : !isRecoverableFailed
      ? "フェーズが再試行可能な状態ではありません。"
      : undefined

  const canCancel = !terminal && isActive
  const cancelBlockedReason = terminal
    ? "ジョブが終端状態のためキャンセルできません。"
    : !isActive
      ? "フェーズがキャンセル可能な状態ではありません。"
      : undefined

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

function buildScreenActionEnablement(
  state: PersonaGenerationPhaseScreenState
): PersonaGenerationPhaseScreenActionEnablement {
  const isBusy = state.phase === "loading" || state.phase === "submitting"

  if (!state.summary) {
    return {
      canStart: false,
      canPause: false,
      canResume: false,
      canRetry: false,
      canCancel: false,
      canCheckBodyReadiness: !isBusy && state.jobId !== null,
      canStartBodyPhase: false
    }
  }

  const derived = derivePersonaActionEnablement(state.summary)

  return {
    canStart: !isBusy && derived.canStart,
    canPause: !isBusy && derived.canPause,
    canResume: !isBusy && derived.canResume,
    canRetry: !isBusy && derived.canRetry,
    canCancel: !isBusy && derived.canCancel,
    canCheckBodyReadiness: !isBusy && state.jobId !== null,
    canStartBodyPhase: !isBusy && derivePersonaCanStartBodyPhase(state.summary)
  }
}

function buildActionCards(
  state: PersonaGenerationPhaseScreenState
): PersonaGenerationPhaseActionCard[] {
  const screenEnablement = buildScreenActionEnablement(state)
  const derived = state.summary
    ? derivePersonaActionEnablement(state.summary)
    : null
  const bodyReadinessBlockedReason = state.summary
    ? derivePersonaBodyReadinessBlockedReason(state.summary)
    : ""

  return [
    {
      id: "start",
      label: "開始",
      disabled: !screenEnablement.canStart,
      blockedReason: derived?.startBlockedReason ?? "",
      tone: "primary"
    },
    {
      id: "pause",
      label: "中断",
      disabled: !screenEnablement.canPause,
      blockedReason: derived?.pauseBlockedReason ?? "",
      tone: "warning"
    },
    {
      id: "resume",
      label: "再開",
      disabled: !screenEnablement.canResume,
      blockedReason: derived?.resumeBlockedReason ?? "",
      tone: "default"
    },
    {
      id: "retry",
      label: "リトライ",
      disabled: !screenEnablement.canRetry,
      blockedReason: derived?.retryBlockedReason ?? "",
      tone: "default"
    },
    {
      id: "cancel",
      label: "キャンセル",
      disabled: !screenEnablement.canCancel,
      blockedReason: derived?.cancelBlockedReason ?? "",
      tone: "warning"
    },
    {
      id: "check-body-readiness",
      label: "本文翻訳の開始可否を確認",
      disabled: !screenEnablement.canCheckBodyReadiness,
      blockedReason: state.jobId === null ? "ジョブIDを指定してください。" : "",
      tone: "default"
    },
    {
      id: "start-body-phase",
      label: "本文翻訳を開始",
      disabled: !screenEnablement.canStartBodyPhase,
      blockedReason: bodyReadinessBlockedReason,
      tone: "primary"
    }
  ]
}

export class PersonaGenerationPhasePresenter {
  toViewModel(
    state: PersonaGenerationPhaseScreenState,
    isGatewayConnected: boolean,
    availableProviders: PersonaSelectOption[] = [],
    availableModels: PersonaSelectOption[] = []
  ): PersonaGenerationPhaseScreenViewModel {
    const viewState = buildViewState(state)
    const statusCopy = buildStatusCopy(state, viewState)
    const summary = state.summary
    const targetSummary = summary?.targetSummary
    const resultSummary = summary?.resultSummary
    const executionSummary = summary?.execution
    const errorSummary = summary?.errorSummary ?? null
    const screenActionEnablement = buildScreenActionEnablement(state)
    const providerLabel =
      summary?.aiSettings?.provider ??
      executionSummary?.provider ??
      "設定未完了"
    const modelLabel =
      summary?.aiSettings?.model ??
      executionSummary?.model ??
      "設定未完了"
    const executionModeLabel =
      summary?.aiSettings?.executionMode ??
      executionSummary?.executionMode ??
      "設定未完了"

    return {
      ...state,
      gatewayStatus: isGatewayConnected ? "接続準備済み" : "未接続",
      viewState,
      isLoading: state.phase === "loading" && !state.hasLoaded,
      isRefreshing: state.phase === "loading" && state.hasLoaded,
      isSubmitting: state.phase === "submitting",
      hasJobSelection: state.jobId !== null,
      currentPhaseLabel: "NPC ペルソナ生成",
      phaseStateLabel: buildPhaseStateLabel(summary?.phaseState),
      statusTitle: statusCopy.title,
      statusText: statusCopy.text,
      progressPercent: summary?.progress.percent ?? 0,
      progressLabel: buildProgressLabel(state),
      progressDetail: buildProgressDetail(state),
      startedAtLabel: formatDate(summary?.startedAt),
      finishedAtLabel: formatDate(summary?.finishedAt),
      targetCountLabel: formatCount(summary?.progress.targetCount),
      generatedCountLabel: formatCount(resultSummary?.generatedCount),
      failedCountLabel: formatCount(resultSummary?.failedCount),
      skippedCountLabel: formatCount(targetSummary?.skippedCount),
      npcCountLabel: formatCount(targetSummary?.targetCount),
      commonPersonaHitCountLabel: formatCount(
        targetSummary?.commonPersonaHitCount
      ),
      commonPersonaMissCountLabel: formatCount(
        targetSummary?.commonPersonaMissCount
      ),
      skippedReasonsLabel: formatList(targetSummary?.skippedReasons),
      targetSnapshotLabel: buildTargetSnapshotLabel(state),
      isExecutionConfigured: summary ? isPersonaExecutionConfigured(summary) : false,
      providerLabel,
      modelLabel,
      executionModeLabel,
      credentialRefLabel: executionSummary?.credentialRef ?? "設定未完了",
      providerOptions: buildPersonaProviderOptions(summary, availableProviders),
      modelOptions: buildPersonaModelOptions(summary, availableModels),
      executionOptions: PERSONA_EXECUTION_OPTIONS,
      inputCountLabel: formatCount(executionSummary?.inputCount),
      outputCountLabel: formatCount(executionSummary?.outputCount),
      evidenceRefsLabel: formatList(executionSummary?.evidenceRefs),
      promptDigestLabel: executionSummary?.promptDigest ?? "-",
      snapshotLabel: buildSnapshotLabel(state),
      snapshotReferenceStatusLabel:
        resultSummary?.snapshotReferenceStatus ?? "-",
      personaCountLabel: formatCount(resultSummary?.personaCount),
      missingCountLabel: formatCount(resultSummary?.missingCount),
      bodyReadinessLabel: summary
        ? (derivePersonaCanStartBodyPhase(summary) ? "Ready" : "Blocked")
        : "-",
      bodyReadinessBlockedReason: summary
        ? derivePersonaBodyReadinessBlockedReason(summary)
        : "",
      bodyReadinessInputSummaryLabel:
        buildBodyReadinessInputSummaryLabel(state),
      errorKindLabel:
        errorSummary?.errorKind ?? "-",
      errorReasonLabel: errorSummary?.reason ?? "-",
      retryableLabel:
        errorSummary === null
          ? "-"
          : errorSummary.retryable
            ? "再試行可能"
            : "再試行不可",
      actionCards: buildActionCards(state),
      screenActionEnablement,
      lastErrorSummary: errorSummary,
      actionEnablement: summary?.actionEnablement ?? null,
      latestProgressSummary: summary?.progress ?? null,
      latestTargetSummary: targetSummary ?? null,
      latestResultSummary: resultSummary ?? null,
      latestExecutionSummary: executionSummary ?? null,
      latestErrorKind:
        errorSummary?.errorKind ?? null,
      latestBodyReadiness: state.bodyReadiness,
      latestBodyReadinessInputSummary: state.bodyReadiness?.inputSummary ?? null
    }
  }
}
