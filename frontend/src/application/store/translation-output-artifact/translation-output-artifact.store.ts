import type {
  TranslationOutputArtifactErrorKind,
  TranslationOutputCompletedJobSummary
} from "@application/gateway-contract/translation-output-artifact"

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

interface TranslationOutputArtifactActionDisablement {
  actionKind: "generate" | "regenerate"
  disabled: boolean
  reason?: string
  errorKind?: TranslationOutputArtifactErrorKind
}

interface TranslationOutputArtifactScreenState {
  phase: "idle" | "loading" | "ready" | "submitting"
  viewState:
    | "loading"
    | "empty"
    | "not_ready"
    | "ready"
    | "generating"
    | "success"
    | "failed"
    | "stale"
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

type Listener = (state: TranslationOutputArtifactScreenState) => void

function cloneCompletedJob(
  job: TranslationOutputCompletedJobSummary
): TranslationOutputCompletedJobSummary {
  return {
    ...job,
    outputStatusDistribution: job.outputStatusDistribution
      ? { ...job.outputStatusDistribution }
      : undefined
  }
}

function cloneReview(
  review: TranslationOutputReviewSnapshot | null
): TranslationOutputReviewSnapshot | null {
  if (!review) {
    return null
  }

  return {
    ...review,
    completedJobs: review.completedJobs.map(cloneCompletedJob),
    rejectionReasons: review.rejectionReasons.map((reason) => ({ ...reason }))
  }
}

function cloneDiffPreview(
  diffPreview: TranslationOutputDiffPreviewSnapshot | null
): TranslationOutputDiffPreviewSnapshot | null {
  if (!diffPreview) {
    return null
  }

  return {
    ...diffPreview,
    rows: diffPreview.rows.map((row) => ({ ...row }))
  }
}

function cloneLastCommand(
  lastCommand: TranslationOutputArtifactCommandSnapshot | null
): TranslationOutputArtifactCommandSnapshot | null {
  if (!lastCommand) {
    return null
  }

  return {
    ...lastCommand,
    affectedFieldIds: lastCommand.affectedFieldIds
      ? [...lastCommand.affectedFieldIds]
      : undefined
  }
}

function createInitialState(): TranslationOutputArtifactScreenState {
  return {
    phase: "idle",
    viewState: "loading",
    completedJobs: [],
    selectedJobId: null,
    selectedArtifactId: null,
    review: null,
    diffPreview: null,
    lastCommand: null,
    actionDisablements: [],
    refreshPending: false,
    targetGame: "skyrim_se",
    outputPath: "/tmp/translation-output-artifact.xml",
    pathState: "valid",
    pathReason: "",
    errorMessage: "",
    pendingAction: null,
    hasLoaded: false
  }
}

export class TranslationOutputArtifactStore {
  private state: TranslationOutputArtifactScreenState = createInitialState()

  private readonly listeners = new Set<Listener>()

  subscribe(listener: Listener): () => void {
    this.listeners.add(listener)
    listener(this.snapshot())
    return () => {
      this.listeners.delete(listener)
    }
  }

  snapshot(): TranslationOutputArtifactScreenState {
    return {
      ...this.state,
      completedJobs: this.state.completedJobs.map(cloneCompletedJob),
      review: cloneReview(this.state.review),
      diffPreview: cloneDiffPreview(this.state.diffPreview),
      lastCommand: cloneLastCommand(this.state.lastCommand),
      actionDisablements: this.state.actionDisablements.map((item) => ({
        ...item
      }))
    }
  }

  update(mutator: (draft: TranslationOutputArtifactScreenState) => void): void {
    const nextState = this.snapshot()
    mutator(nextState)
    this.state = nextState
    this.emit()
  }

  private emit(): void {
    const snapshot = this.snapshot()
    for (const listener of this.listeners) {
      listener(snapshot)
    }
  }
}
