import type {
  TranslationOutputArtifactErrorKind,
  TranslationOutputArtifactStatusSummary,
  TranslationOutputCompletedJobSummary
} from "@application/gateway-contract/translation-output-artifact"

type TranslationOutputArtifactViewState =
  | "loading"
  | "empty"
  | "not_ready"
  | "ready"
  | "generating"
  | "success"
  | "failed"
  | "stale"

interface TranslationOutputArtifactCommandSnapshot {
  jobId: number
  artifactId: number
  artifactStatus: string
  rowCount: number
  filePath?: string
  targetGame: string
  operationKind?: string
  replacedArtifactId?: number
  affectedFieldIds?: number[]
  duplicateRowCreated?: boolean
  errorKind?: TranslationOutputArtifactErrorKind
  errorReason?: string
  retryable?: boolean
  isRedacted?: boolean
}

interface TranslationOutputReviewSnapshot {
  completedJobs: TranslationOutputCompletedJobSummary[]
  selectedJobId: number
  selectedJobStatus: string
  bodyPhaseStatus: string
  outputReady: boolean
  translatedCount: number
  rowCount: number
  inputSnapshotDigest: string
  sourceFileDigest: string
  readiness: boolean
  retryable: boolean
  artifactId: number
  artifactStatus: string
  artifactRowCount: number
  currentVersion: boolean
  rejectionKind?: TranslationOutputArtifactErrorKind
  rejectionReasons: {
    errorKind: TranslationOutputArtifactErrorKind
    reason: string
    retryable: boolean
    isRedacted: boolean
  }[]
}

interface TranslationOutputDiffPreviewSnapshot {
  artifactId: number
  rows: {
    fieldId: number
    rowDigest: string
    edid: string
    rec: string
    field: string
    formId: string
    sourceExcerpt: string
    destExcerpt: string
    xTranslatorStatus: number
    internalOutputStatus: string
    rowReflectionSummary: string
    staleReason?: string
    canRegenerate: boolean
  }[]
  compatibilityPassed: boolean
  compatibilityWarningCount: number
  compatibilityRejectCount: number
}

interface TranslationOutputArtifactScreenState {
  phase: "idle" | "loading" | "ready" | "submitting"
  viewState: TranslationOutputArtifactViewState
  completedJobs: TranslationOutputCompletedJobSummary[]
  selectedJobId: number | null
  selectedArtifactId: number | null
  review: TranslationOutputReviewSnapshot | null
  diffPreview: TranslationOutputDiffPreviewSnapshot | null
  lastCommand: TranslationOutputArtifactCommandSnapshot | null
  actionDisablements: TranslationOutputArtifactActionDisablement[]
  refreshPending: boolean
  targetGame: string
  outputPath: string
  pathState: "empty" | "valid" | "invalid" | "readonly"
  pathReason: string
  errorMessage: string
  pendingAction: "generate" | "regenerate" | "refresh" | null
  hasLoaded: boolean
}

interface TranslationOutputArtifactActionDisablement {
  actionKind: "generate" | "regenerate"
  disabled: boolean
  reason?: string
  errorKind?: TranslationOutputArtifactErrorKind
}

interface TranslationOutputArtifactScreenViewModel extends TranslationOutputArtifactScreenState {
  canGenerate: boolean
  canRegenerate: boolean
  primaryAction: "generate" | "regenerate"
  disabledReason?: string
  selectedJobStatus?: string
  selectedArtifactStatus?: string
  gatewayStatus: string
  isLoading: boolean
  isSubmitting: boolean
  statusTitle: string
  statusText: string
  artifactStatusSummary: TranslationOutputArtifactStatusSummary | null
  compatibilitySummaryText: string
}

function getDisablement(
  items: TranslationOutputArtifactActionDisablement[],
  actionKind: "generate" | "regenerate"
): TranslationOutputArtifactActionDisablement | null {
  return items.find((item) => item.actionKind === actionKind) ?? null
}

function buildStatusTitle(state: TranslationOutputArtifactScreenState): string {
  switch (state.viewState) {
    case "loading":
      return "Output Review を読み込み中です。"
    case "empty":
      return "出力対象の completed job はありません。"
    case "not_ready":
      return "出力前の確認が必要です。"
    case "ready":
      return "出力条件を満たしています。"
    case "generating":
      return "xTranslator XML を生成しています。"
    case "success":
      return "最新 artifact を確認できます。"
    case "failed":
      return "出力に失敗しました。"
    case "stale":
      return "artifact と差分の再確認が必要です。"
  }
}

function buildStatusText(state: TranslationOutputArtifactScreenState): string {
  if (state.errorMessage) {
    return state.errorMessage
  }

  const review = state.review
  if (!review) {
    return "completed job、result summary、diff preview を読み込むと状態が更新されます。"
  }

  if (review.rejectionReasons.length > 0) {
    return review.rejectionReasons.map((reason) => reason.reason).join(" / ")
  }

  if (state.lastCommand?.errorReason) {
    return state.lastCommand.errorReason
  }

  if (state.viewState === "success" && state.lastCommand?.filePath) {
    return `出力先: ${state.lastCommand.filePath}`
  }

  if (state.viewState === "stale") {
    return "stale reason と canRegenerate を確認し、必要なら再出力してください。"
  }

  return "job list、summary、diff preview、action rail を同じ画面で確認できます。"
}

export class TranslationOutputArtifactPresenter {
  toViewModel(
    state: TranslationOutputArtifactScreenState,
    isGatewayConnected: boolean
  ): TranslationOutputArtifactScreenViewModel {
    const generateDisablement = getDisablement(
      state.actionDisablements,
      "generate"
    )
    const regenerateDisablement = getDisablement(
      state.actionDisablements,
      "regenerate"
    )
    const artifactStatusSummary = state.review
      ? {
          artifactId: state.review.artifactId,
          status: state.review.artifactStatus,
          rowCount: state.review.artifactRowCount,
          currentVersion: state.review.currentVersion
        }
      : null

    return {
      ...state,
      canGenerate: generateDisablement?.disabled !== true,
      canRegenerate: regenerateDisablement?.disabled !== true,
      primaryAction:
        artifactStatusSummary && artifactStatusSummary.artifactId > 0
          ? "regenerate"
          : "generate",
      disabledReason:
        artifactStatusSummary && artifactStatusSummary.artifactId > 0
          ? regenerateDisablement?.reason
          : generateDisablement?.reason,
      selectedJobStatus: state.review?.selectedJobStatus,
      selectedArtifactStatus: state.review?.artifactStatus,
      gatewayStatus: isGatewayConnected ? "接続準備済み" : "未接続",
      isLoading: state.phase === "loading",
      isSubmitting: state.phase === "submitting",
      statusTitle: buildStatusTitle(state),
      statusText: buildStatusText(state),
      artifactStatusSummary,
      compatibilitySummaryText: state.diffPreview
        ? `warning ${state.diffPreview.compatibilityWarningCount} 件 / reject ${state.diffPreview.compatibilityRejectCount} 件`
        : "未取得"
    }
  }
}
