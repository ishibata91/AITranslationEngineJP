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
  providerLabel: string
  modelLabel: string
  executionModeLabel: string
  credentialRefLabel: string
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
    state.bodyReadiness?.errorKind === "snapshot_missing" ||
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

function buildBodyReadinessLabel(
  state: PersonaGenerationPhaseScreenState
): string {
  if (typeof state.bodyReadiness?.ready === "boolean") {
    return state.bodyReadiness.ready ? "Ready" : "Blocked"
  }

  if (typeof state.summary?.resultSummary?.bodyReadiness === "boolean") {
    return state.summary.resultSummary.bodyReadiness ? "Ready" : "Blocked"
  }

  return "-"
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

function buildScreenActionEnablement(
  state: PersonaGenerationPhaseScreenState
): PersonaGenerationPhaseScreenActionEnablement {
  const summaryEnablement = state.summary?.actionEnablement
  const isBusy = state.phase === "loading" || state.phase === "submitting"

  return {
    canStart: !isBusy && (summaryEnablement?.canStart ?? false),
    canPause: !isBusy && (summaryEnablement?.canPause ?? false),
    canResume: !isBusy && (summaryEnablement?.canResume ?? false),
    canRetry: !isBusy && (summaryEnablement?.canRetry ?? false),
    canCancel: !isBusy && (summaryEnablement?.canCancel ?? false),
    canCheckBodyReadiness: !isBusy && state.jobId !== null,
    canStartBodyPhase:
      !isBusy && (summaryEnablement?.canStartBodyPhase ?? false)
  }
}

function buildActionCards(
  state: PersonaGenerationPhaseScreenState
): PersonaGenerationPhaseActionCard[] {
  const enablement = state.summary?.actionEnablement
  const screenEnablement = buildScreenActionEnablement(state)
  const bodyReadinessBlockedReason =
    state.bodyReadiness?.blockedReason ??
    enablement?.bodyPhaseBlockedReason ??
    ""

  return [
    {
      id: "start",
      label: "開始",
      disabled: !screenEnablement.canStart,
      blockedReason: enablement?.startBlockedReason ?? "",
      tone: "primary"
    },
    {
      id: "pause",
      label: "中断",
      disabled: !screenEnablement.canPause,
      blockedReason: enablement?.pauseBlockedReason ?? "",
      tone: "warning"
    },
    {
      id: "resume",
      label: "再開",
      disabled: !screenEnablement.canResume,
      blockedReason: enablement?.resumeBlockedReason ?? "",
      tone: "default"
    },
    {
      id: "retry",
      label: "リトライ",
      disabled: !screenEnablement.canRetry,
      blockedReason: enablement?.retryBlockedReason ?? "",
      tone: "default"
    },
    {
      id: "cancel",
      label: "キャンセル",
      disabled: !screenEnablement.canCancel,
      blockedReason: enablement?.cancelBlockedReason ?? "",
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
    isGatewayConnected: boolean
  ): PersonaGenerationPhaseScreenViewModel {
    const viewState = buildViewState(state)
    const statusCopy = buildStatusCopy(state, viewState)
    const summary = state.summary
    const targetSummary = summary?.targetSummary
    const resultSummary = summary?.resultSummary
    const executionSummary = summary?.execution
    const errorSummary = summary?.errorSummary ?? null
    const screenActionEnablement = buildScreenActionEnablement(state)

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
      providerLabel: executionSummary?.provider ?? "-",
      modelLabel: executionSummary?.model ?? "-",
      executionModeLabel: executionSummary?.executionMode ?? "-",
      credentialRefLabel: executionSummary?.credentialRef ?? "-",
      inputCountLabel: formatCount(executionSummary?.inputCount),
      outputCountLabel: formatCount(executionSummary?.outputCount),
      evidenceRefsLabel: formatList(executionSummary?.evidenceRefs),
      promptDigestLabel: executionSummary?.promptDigest ?? "-",
      snapshotLabel: buildSnapshotLabel(state),
      snapshotReferenceStatusLabel:
        resultSummary?.snapshotReferenceStatus ?? "-",
      personaCountLabel: formatCount(resultSummary?.personaCount),
      missingCountLabel: formatCount(resultSummary?.missingCount),
      bodyReadinessLabel: buildBodyReadinessLabel(state),
      bodyReadinessBlockedReason: state.bodyReadiness?.blockedReason ?? "",
      bodyReadinessInputSummaryLabel:
        buildBodyReadinessInputSummaryLabel(state),
      errorKindLabel:
        state.bodyReadiness?.errorKind ?? errorSummary?.errorKind ?? "-",
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
        errorSummary?.errorKind ?? state.bodyReadiness?.errorKind ?? null,
      latestBodyReadiness: state.bodyReadiness,
      latestBodyReadinessInputSummary: state.bodyReadiness?.inputSummary ?? null
    }
  }
}
