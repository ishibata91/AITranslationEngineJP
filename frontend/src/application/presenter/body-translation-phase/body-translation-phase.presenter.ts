import type {
  ProcessingTargetListPageState,
  ProcessingTargetListPageStatesByPhase
} from "@application/gateway-contract/processing-target"

import type {
  BodyTranslationPhaseActionEnablement,
  BodyTranslationPhaseErrorKind,
  BodyTranslationPhaseErrorSummary,
  BodyTranslationPhaseExecutionSummary,
  BodyTranslationOutputReadinessResponse,
  BodyTranslationPhaseFieldResultItem as BodyTranslationPhaseGatewayFieldResultItem,
  BodyTranslationPhaseFieldResultSummary,
  BodyTranslationPhaseInputSummary,
  BodyTranslationPhaseProgressSummary,
  BodyTranslationPhaseRequestSummary,
  BodyTranslationPhaseSummaryResponse
} from "@application/gateway-contract/body-translation-phase"

type BodyTranslationPhaseActionKind =
  | "refresh"
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
  outputReadiness: BodyTranslationOutputReadinessResponse | null
  errorMessage: string
  pendingAction: BodyTranslationPhaseActionKind | null
  hasLoaded: boolean
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
  actionCards: {
    id: BodyTranslationPhaseActionKind
    label: string
    disabled: boolean
    blockedReason: string
    tone: "default" | "primary" | "warning"
  }[]
  screenActionEnablement: {
    canRefresh: boolean
    canStart: boolean
    canPause: boolean
    canResume: boolean
    canRetry: boolean
    canCancel: boolean
    canCheckOutputReadiness: boolean
  }
  lastErrorSummary: BodyTranslationPhaseErrorSummary | null
  actionEnablement: BodyTranslationPhaseActionEnablement | null
  latestProgressSummary: BodyTranslationPhaseProgressSummary | null
  latestInputSummary: BodyTranslationPhaseInputSummary | null
  latestRequestSummary: BodyTranslationPhaseRequestSummary | null
  latestExecutionSummary: BodyTranslationPhaseExecutionSummary | null
  latestResultSummary: BodyTranslationPhaseFieldResultSummary | null
  latestErrorKind: BodyTranslationPhaseErrorKind | null
  latestOutputReadiness: BodyTranslationOutputReadinessResponse | null
  pauseRequestShape?: { jobId: number; phaseRunId: number }
  resumeRequestShape?: { jobId: number; phaseRunId: number }
  retryRequestShape?: { jobId: number; phaseRunId: number }
  cancelRequestShape?: { jobId: number; phaseRunId: number }
  outputReadinessRequestShape?: { jobId: number }
}

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

function getEffectiveReadiness(
  state: BodyTranslationPhaseScreenState
): BodyTranslationOutputReadinessResponse | null {
  if (state.outputReadiness) {
    return state.outputReadiness
  }

  const summary = state.summary
  if (!summary) {
    return null
  }

  return {
    jobId: summary.jobId,
    currentPhase: summary.currentPhase,
    phaseState: summary.phaseState,
    ready: summary.outputReadiness.ready,
    blockedReason: summary.outputReadiness.blockedReason,
    errorKind: summary.outputReadiness.errorKind,
    completedFieldCount: summary.outputReadiness.completedFieldCount,
    statusConsistent: summary.outputReadiness.statusConsistent,
    outputCount: summary.execution.outputCount
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

  if (summary.actionEnablement.canStart) {
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
  return `${progress.processedCount.toLocaleString("ja-JP")} / ${progress.totalCount.toLocaleString("ja-JP")} 件 / 成功 ${progress.translatedCount.toLocaleString("ja-JP")} 件 / スキップ ${progress.skippedCount.toLocaleString("ja-JP")} 件 / ${buildCurrentStepLabel(progress.currentStep)}`
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
  if (phaseState === "completed" && summary.execution.requestUnitCount === 0) {
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
  return summary.execution.requestUnitCount > 0 ? "provider 実行あり" : "-"
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
  const enablement = state.summary?.actionEnablement
  const readiness = getEffectiveReadiness(state)
  const isBusy = state.phase === "loading" || state.phase === "submitting"

  return [
    {
      id: "start",
      label: "開始",
      disabled: isBusy || !(enablement?.canStart ?? false),
      blockedReason: enablement?.startBlockedReason ?? "",
      tone: "primary"
    },
    {
      id: "pause",
      label: "中断",
      disabled: isBusy || !(enablement?.canPause ?? false),
      blockedReason: enablement?.pauseBlockedReason ?? "",
      tone: "warning"
    },
    {
      id: "resume",
      label: "再開",
      disabled: isBusy || !(enablement?.canResume ?? false),
      blockedReason: enablement?.resumeBlockedReason ?? "",
      tone: "default"
    },
    {
      id: "retry",
      label: "リトライ",
      disabled: isBusy || !(enablement?.canRetry ?? false),
      blockedReason: enablement?.retryBlockedReason ?? "",
      tone: "default"
    },
    {
      id: "cancel",
      label: "キャンセル",
      disabled: isBusy || !(enablement?.canCancel ?? false),
      blockedReason: enablement?.cancelBlockedReason ?? "",
      tone: "warning"
    },
    {
      id: "check-output-readiness",
      label: "出力準備を確認",
      disabled: isBusy || !(enablement?.canCheckOutputReadiness ?? false),
      blockedReason:
        enablement?.outputReadinessBlockedReason ??
        readiness?.blockedReason ??
        "",
      tone: readiness?.ready ? "primary" : "default"
    },
    {
      id: "refresh",
      label: "更新",
      disabled: isBusy || state.jobId === null,
      blockedReason: "",
      tone: "default"
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
  state: BodyTranslationPhaseScreenState,
  summary: BodyTranslationPhaseSummaryResponse | null
): BodyTranslationPhaseScreenViewModel["screenActionEnablement"] {
  return {
    canRefresh: state.jobId !== null,
    canStart: summary?.actionEnablement.canStart ?? false,
    canPause: summary?.actionEnablement.canPause ?? false,
    canResume: summary?.actionEnablement.canResume ?? false,
    canRetry: summary?.actionEnablement.canRetry ?? false,
    canCancel: summary?.actionEnablement.canCancel ?? false,
    canCheckOutputReadiness:
      summary?.actionEnablement.canCheckOutputReadiness ?? false
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
    isGatewayConnected: boolean
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
      currentPhaseLabel: summary?.currentPhase ?? "body_translation",
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
      providerLabel: summary?.execution.provider ?? "-",
      modelLabel: summary?.execution.model ?? "-",
      executionModeLabel: summary?.execution.executionMode ?? "-",
      credentialRefLabel: summary?.execution.credentialRef ?? "-",
      providerTargetCountLabel: formatCount(
        summary?.requestSummary?.providerTargetCount
      ),
      exactDictionaryExclusionCountLabel: formatCount(
        summary?.requestSummary?.exactDictionaryExclusionCount
      ),
      partialDictionaryConstraintCountLabel: formatCount(
        summary?.requestSummary?.partialDictionaryConstraintCount
      ),
      requestUnitCountLabel: formatCount(summary?.execution.requestUnitCount),
      outputCountLabel: formatCount(summary?.execution.outputCount),
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
        errorSummary?.errorKind ?? outputReadiness?.errorKind ?? "-",
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
      screenActionEnablement: buildScreenActionEnablement(state, summary),
      lastErrorSummary: errorSummary,
      actionEnablement: summary?.actionEnablement ?? null,
      latestProgressSummary: summary?.progress ?? null,
      latestInputSummary: summary?.inputSummary ?? null,
      latestRequestSummary: summary?.requestSummary ?? null,
      latestExecutionSummary: summary?.execution ?? null,
      latestResultSummary: resultSummary,
      latestErrorKind:
        errorSummary?.errorKind ?? outputReadiness?.errorKind ?? null,
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
