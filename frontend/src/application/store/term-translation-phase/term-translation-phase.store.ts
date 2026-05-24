import type { ProcessingTargetListPageState } from "@application/gateway-contract/processing-target"

import type {
  TermTranslationNextPhaseReadinessResponse,
  TermTranslationPhaseSummaryResponse
} from "@application/gateway-contract/term-translation-phase"

type TermTranslationPhaseActionKind =
  | "refresh"
  | "start"
  | "pause"
  | "resume"
  | "retry"

interface TermTranslationPhaseScreenState {
  jobId: number | null
  phase: "idle" | "loading" | "ready" | "submitting"
  summary: TermTranslationPhaseSummaryResponse | null
  nextPhaseReadiness: TermTranslationNextPhaseReadinessResponse | null
  errorMessage: string
  pendingAction: TermTranslationPhaseActionKind | null
  hasLoaded: boolean
  processingTargetPageState: ProcessingTargetListPageState | null
}

type Listener = (state: TermTranslationPhaseScreenState) => void

function cloneSummary(
  summary: TermTranslationPhaseSummaryResponse | null
): TermTranslationPhaseSummaryResponse | null {
  if (!summary) {
    return null
  }

  return {
    ...summary,
    progress: { ...summary.progress },
    execution: { ...summary.execution },
    resultSummary: summary.resultSummary
      ? { ...summary.resultSummary }
      : undefined,
    errorSummary: summary.errorSummary
      ? { ...summary.errorSummary }
      : undefined,
    actionEnablement: { ...summary.actionEnablement }
  }
}

function cloneNextPhaseReadiness(
  readiness: TermTranslationNextPhaseReadinessResponse | null
): TermTranslationNextPhaseReadinessResponse | null {
  if (!readiness) {
    return null
  }

  return { ...readiness }
}

function cloneProcessingTargetPageState(
  pageState: ProcessingTargetListPageState | null
): ProcessingTargetListPageState | null {
  if (!pageState) {
    return null
  }

  return {
    ...pageState,
    items: pageState.items.map((item) => ({
      ...item,
      titleParts: item.titleParts.map((part) => ({ ...part })),
      metadata: item.metadata.map((metadata) => ({ ...metadata }))
    })),
    metadata: pageState.metadata.map((metadata) => ({ ...metadata }))
  }
}

function createInitialState(): TermTranslationPhaseScreenState {
  return {
    jobId: null,
    phase: "idle",
    summary: null,
    nextPhaseReadiness: null,
    errorMessage: "",
    pendingAction: null,
    hasLoaded: false,
    processingTargetPageState: null
  }
}

export class TermTranslationPhaseStore {
  private state: TermTranslationPhaseScreenState = createInitialState()

  private readonly listeners = new Set<Listener>()

  subscribe(listener: Listener): () => void {
    this.listeners.add(listener)
    listener(this.snapshot())
    return () => {
      this.listeners.delete(listener)
    }
  }

  snapshot(): TermTranslationPhaseScreenState {
    return {
      ...this.state,
      summary: cloneSummary(this.state.summary),
      nextPhaseReadiness: cloneNextPhaseReadiness(
        this.state.nextPhaseReadiness
      ),
      processingTargetPageState: cloneProcessingTargetPageState(
        this.state.processingTargetPageState
      )
    }
  }

  update(mutator: (draft: TermTranslationPhaseScreenState) => void): void {
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
