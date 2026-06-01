import type { ProcessingTargetListPageState } from "@application/gateway-contract/processing-target"

import type {
  PersonaGenerationBodyReadinessResponse,
  PersonaGenerationPhaseSummaryResponse
} from "@application/gateway-contract/persona-generation-phase"

type PersonaGenerationPhaseActionKind =
  | "start"
  | "pause"
  | "resume"
  | "retry"
  | "cancel"
  | "check-body-readiness"
  | "start-body-phase"

interface PersonaGenerationPhaseScreenState {
  jobId: number | null
  phase: "idle" | "loading" | "ready" | "submitting"
  summary: PersonaGenerationPhaseSummaryResponse | null
  bodyReadiness: PersonaGenerationBodyReadinessResponse | null
  errorMessage: string
  pendingAction: PersonaGenerationPhaseActionKind | null
  hasLoaded: boolean
  initialFetchDone: boolean
  processingTargetPageState: ProcessingTargetListPageState | null
}

type Listener = (state: PersonaGenerationPhaseScreenState) => void

function cloneSummary(
  summary: PersonaGenerationPhaseSummaryResponse | null
): PersonaGenerationPhaseSummaryResponse | null {
  if (!summary) {
    return null
  }

  return {
    ...summary,
    progress: { ...summary.progress },
    targetSummary: {
      ...summary.targetSummary,
      skippedReasons: [...summary.targetSummary.skippedReasons]
    },
    execution: summary.execution
      ? {
          ...summary.execution,
          evidenceRefs: summary.execution.evidenceRefs
            ? [...summary.execution.evidenceRefs]
            : []
        }
      : undefined,
    resultSummary: summary.resultSummary
      ? { ...summary.resultSummary }
      : undefined,
    errorSummary: summary.errorSummary
      ? { ...summary.errorSummary }
      : undefined,
    actionEnablement: { ...summary.actionEnablement }
  }
}

function cloneBodyReadiness(
  readiness: PersonaGenerationBodyReadinessResponse | null
): PersonaGenerationBodyReadinessResponse | null {
  if (!readiness) {
    return null
  }

  return {
    ...readiness,
    inputSummary: {
      ...readiness.inputSummary,
      evidenceRefs: [...readiness.inputSummary.evidenceRefs]
    }
  }
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

function createInitialState(): PersonaGenerationPhaseScreenState {
  return {
    jobId: null,
    phase: "idle",
    summary: null,
    bodyReadiness: null,
    errorMessage: "",
    pendingAction: null,
    hasLoaded: false,
    initialFetchDone: false,
    processingTargetPageState: null
  }
}

export class PersonaGenerationPhaseStore {
  private state: PersonaGenerationPhaseScreenState = createInitialState()

  private readonly listeners = new Set<Listener>()

  subscribe(listener: Listener): () => void {
    this.listeners.add(listener)
    listener(this.snapshot())
    return () => {
      this.listeners.delete(listener)
    }
  }

  snapshot(): PersonaGenerationPhaseScreenState {
    return {
      ...this.state,
      summary: cloneSummary(this.state.summary),
      bodyReadiness: cloneBodyReadiness(this.state.bodyReadiness),
      processingTargetPageState: cloneProcessingTargetPageState(
        this.state.processingTargetPageState
      )
    }
  }

  update(mutator: (draft: PersonaGenerationPhaseScreenState) => void): void {
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
