import type {
  TranslationJobManagementJobDetail,
  TranslationJobManagementJobSummary,
  TranslationJobManagementOperationAvailability,
  TranslationJobManagementReasonCategory
} from "@application/gateway-contract/translation-job-management"

type TranslationJobManagementStoreFilterId =
  | "all"
  | "Ready"
  | "Running"
  | "Paused"
  | "RecoverableFailed"
  | "Failed"
  | "Canceled"

interface TranslationJobManagementFeedbackState {
  tone: "info" | "success" | "warning" | "danger"
  title: string
  message: string
  category?:
    | TranslationJobManagementReasonCategory
    | "stop_requested"
    | "delete_success"
    | "resume_success"
}

interface TranslationJobManagementScreenState {
  phase: "idle" | "loading" | "ready" | "empty" | "error"
  jobs: TranslationJobManagementJobSummary[]
  selectedJobId: number | null
  selectedJobDetail: TranslationJobManagementJobDetail | null
  detailPhase: "idle" | "loading" | "ready" | "stale"
  filterId: TranslationJobManagementStoreFilterId
  searchQuery: string
  isReloading: boolean
  activeOperation: TranslationJobManagementOperationAvailability["kind"] | "reload" | null
  isDeleteConfirmationOpen: boolean
  feedback: TranslationJobManagementFeedbackState | null
}

type Listener = (state: TranslationJobManagementScreenState) => void

function cloneOperation(
  value: TranslationJobManagementOperationAvailability
): TranslationJobManagementOperationAvailability {
  return { ...value }
}

function cloneSummary(
  summary: TranslationJobManagementJobSummary
): TranslationJobManagementJobSummary {
  return {
    ...summary,
    inputSource: { ...summary.inputSource },
    progress: { ...summary.progress },
    stopAvailability: cloneOperation(summary.stopAvailability),
    resumeAvailability: cloneOperation(summary.resumeAvailability),
    deleteAvailability: cloneOperation(summary.deleteAvailability)
  }
}

function cloneDetail(
  detail: TranslationJobManagementJobDetail | null
): TranslationJobManagementJobDetail | null {
  if (!detail) {
    return null
  }

  return {
    ...cloneSummary(detail),
    cacheState: detail.cacheState,
    cacheStateLabel: detail.cacheStateLabel,
    runtimeSummary: { ...detail.runtimeSummary },
    resumeBlockedReasons: detail.resumeBlockedReasons.map((reason) => ({
      ...reason
    })),
    warnings: detail.warnings.map((reason) => ({ ...reason })),
    deleteImpactLines: [...detail.deleteImpactLines]
  }
}

function cloneFeedback(
  feedback: TranslationJobManagementFeedbackState | null
): TranslationJobManagementFeedbackState | null {
  return feedback ? { ...feedback } : null
}

function createInitialState(): TranslationJobManagementScreenState {
  return {
    phase: "idle",
    jobs: [],
    selectedJobId: null,
    selectedJobDetail: null,
    detailPhase: "idle",
    filterId: "all",
    searchQuery: "",
    isReloading: false,
    activeOperation: null,
    isDeleteConfirmationOpen: false,
    feedback: null
  }
}

export class TranslationJobManagementStore {
  private state: TranslationJobManagementScreenState = createInitialState()

  private readonly listeners = new Set<Listener>()

  subscribe(listener: Listener): () => void {
    this.listeners.add(listener)
    listener(this.snapshot())
    return () => {
      this.listeners.delete(listener)
    }
  }

  snapshot(): TranslationJobManagementScreenState {
    return {
      ...this.state,
      jobs: this.state.jobs.map((job) => cloneSummary(job)),
      selectedJobDetail: cloneDetail(this.state.selectedJobDetail),
      feedback: cloneFeedback(this.state.feedback)
    }
  }

  update(mutator: (draft: TranslationJobManagementScreenState) => void): void {
    const draft = this.snapshot()
    mutator(draft)
    this.state = draft
    this.emit()
  }

  private emit(): void {
    const snapshot = this.snapshot()
    for (const listener of this.listeners) {
      listener(snapshot)
    }
  }
}
