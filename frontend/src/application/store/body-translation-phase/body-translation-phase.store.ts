import type {
  BodyTranslationOutputReadinessResponse,
  BodyTranslationPhaseFieldResultItem,
  BodyTranslationPhaseSummaryResponse
} from "@application/gateway-contract/body-translation-phase"

interface BodyTranslationPhaseScreenState {
  jobId: number | null
  phase: "idle" | "loading" | "ready" | "submitting"
  summary: BodyTranslationPhaseSummaryResponse | null
  outputReadiness: BodyTranslationOutputReadinessResponse | null
  errorMessage: string
  pendingAction:
    | "refresh"
    | "start"
    | "pause"
    | "resume"
    | "retry"
    | "cancel"
    | "check-output-readiness"
    | null
  hasLoaded: boolean
}

type Listener = (state: BodyTranslationPhaseScreenState) => void

function cloneFieldResult(
  fieldResult: BodyTranslationPhaseFieldResultItem
): BodyTranslationPhaseFieldResultItem {
  return {
    ...fieldResult,
    identity: fieldResult.identity ? { ...fieldResult.identity } : undefined
  }
}

function cloneSummary(
  summary: BodyTranslationPhaseSummaryResponse | null
): BodyTranslationPhaseSummaryResponse | null {
  if (!summary) {
    return null
  }

  return {
    ...summary,
    progress: { ...summary.progress },
    inputSummary: {
      ...summary.inputSummary,
      skippedReasons: summary.inputSummary.skippedReasons
        ? [...summary.inputSummary.skippedReasons]
        : undefined
    },
    requestSummary: summary.requestSummary
      ? { ...summary.requestSummary }
      : undefined,
    execution: { ...summary.execution },
    fieldResults: summary.fieldResults?.map(cloneFieldResult),
    resultSummary: summary.resultSummary
      ? {
          ...summary.resultSummary,
          fieldResults: summary.resultSummary.fieldResults?.map(cloneFieldResult)
        }
      : undefined,
    errorSummary: summary.errorSummary
      ? { ...summary.errorSummary }
      : undefined,
    actionEnablement: { ...summary.actionEnablement },
    outputReadiness: { ...summary.outputReadiness }
  }
}

function cloneOutputReadiness(
  readiness: BodyTranslationOutputReadinessResponse | null
): BodyTranslationOutputReadinessResponse | null {
  if (!readiness) {
    return null
  }

  return { ...readiness }
}

function createInitialState(): BodyTranslationPhaseScreenState {
  return {
    jobId: null,
    phase: "idle",
    summary: null,
    outputReadiness: null,
    errorMessage: "",
    pendingAction: null,
    hasLoaded: false
  }
}

export class BodyTranslationPhaseStore {
  private state: BodyTranslationPhaseScreenState = createInitialState()

  private readonly listeners = new Set<Listener>()

  subscribe(listener: Listener): () => void {
    this.listeners.add(listener)
    listener(this.snapshot())
    return () => {
      this.listeners.delete(listener)
    }
  }

  snapshot(): BodyTranslationPhaseScreenState {
    return {
      ...this.state,
      summary: cloneSummary(this.state.summary),
      outputReadiness: cloneOutputReadiness(this.state.outputReadiness)
    }
  }

  update(mutator: (draft: BodyTranslationPhaseScreenState) => void): void {
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
