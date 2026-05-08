import type {
  GetTermTranslationNextPhaseReadinessRequest,
  GetTermTranslationPhaseSummaryRequest,
  PauseTermTranslationPhaseRequest,
  TermTranslationNextPhaseReadinessResponse,
  RetryTermTranslationPhaseRequest,
  ResumeTermTranslationPhaseRequest,
  StartTermTranslationPhaseRequest,
  TermTranslationPhaseCommandResponse,
  TermTranslationPhaseGatewayContract,
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
}

interface TermTranslationPhaseStoreLike {
  snapshot(): TermTranslationPhaseScreenState
  update(mutator: (draft: TermTranslationPhaseScreenState) => void): void
}

function sanitizeErrorMessage(error: unknown, fallback: string): string {
  if (
    error instanceof Error &&
    error.message.startsWith("Wails binding is not wired yet:")
  ) {
    return error.message
  }

  return fallback
}

function createNoJobSelectedMessage(): string {
  return "ジョブIDを指定して summary を取得してください。"
}

function createGatewayDisconnectedMessage(): string {
  return "単語翻訳段階の gateway が未接続です。"
}

function patchSummaryFromCommand(
  currentSummary: TermTranslationPhaseSummaryResponse | null,
  response: TermTranslationPhaseCommandResponse
): TermTranslationPhaseSummaryResponse | null {
  if (!currentSummary) {
    return null
  }

  return {
    ...currentSummary,
    jobId: response.jobId,
    currentPhase: response.currentPhase,
    phaseState: response.phaseState,
    phaseRunId: response.phaseRunId,
    startedAt: response.startedAt,
    finishedAt: response.finishedAt,
    progress: { ...response.progress },
    errorSummary: response.errorSummary
      ? { ...response.errorSummary }
      : undefined,
    actionEnablement: {
      ...currentSummary.actionEnablement,
      canRetry: response.retryable,
      canStartNextPhase: response.canStartNextPhase
    }
  }
}

export class TermTranslationPhaseUseCase {
  constructor(
    private readonly gateway: TermTranslationPhaseGatewayContract | null,
    private readonly store: TermTranslationPhaseStoreLike
  ) {}

  async load(): Promise<void> {
    const state = this.store.snapshot()
    if (state.jobId === null) {
      this.store.update((draft) => {
        draft.phase = "ready"
        draft.errorMessage = createNoJobSelectedMessage()
      })
      return
    }

    await this.refresh()
  }

  async setJobId(jobId: number | null): Promise<void> {
    this.store.update((draft) => {
      draft.jobId = jobId
      draft.summary = null
      draft.nextPhaseReadiness = null
      draft.hasLoaded = false
      draft.pendingAction = null
      draft.errorMessage = jobId === null ? createNoJobSelectedMessage() : ""
      draft.phase = "ready"
    })

    if (jobId !== null) {
      await this.refresh()
    }
  }

  async refresh(): Promise<void> {
    const state = this.store.snapshot()
    if (state.jobId === null) {
      this.store.update((draft) => {
        draft.phase = "ready"
        draft.errorMessage = createNoJobSelectedMessage()
      })
      return
    }

    if (!this.gateway) {
      this.store.update((draft) => {
        draft.phase = "ready"
        draft.errorMessage = createGatewayDisconnectedMessage()
      })
      return
    }

    await this.fetchSummaryAndReadiness(state.jobId, "refresh")
  }

  async startPhase(): Promise<void> {
    await this.runCommand("start", (jobId) =>
      this.gateway!.startTermTranslationPhase({
        jobId
      } satisfies StartTermTranslationPhaseRequest)
    )
  }

  async pausePhase(): Promise<void> {
    await this.runCommand("pause", (jobId, phaseRunId) =>
      this.gateway!.pauseTermTranslationPhase({
        jobId,
        phaseRunId
      } satisfies PauseTermTranslationPhaseRequest)
    )
  }

  async resumePhase(): Promise<void> {
    await this.runCommand("resume", (jobId, phaseRunId) =>
      this.gateway!.resumeTermTranslationPhase({
        jobId,
        phaseRunId
      } satisfies ResumeTermTranslationPhaseRequest)
    )
  }

  async retryPhase(): Promise<void> {
    await this.runCommand("retry", (jobId, phaseRunId) =>
      this.gateway!.retryTermTranslationPhase({
        jobId,
        phaseRunId
      } satisfies RetryTermTranslationPhaseRequest)
    )
  }

  private async fetchSummaryAndReadiness(
    jobId: number,
    pendingAction: TermTranslationPhaseActionKind
  ): Promise<void> {
    const before = this.store.snapshot()

    this.store.update((draft) => {
      draft.phase = "loading"
      draft.pendingAction = pendingAction
      draft.errorMessage = ""
    })

    try {
      const [summary, nextPhaseReadiness] = await Promise.all([
        this.gateway!.getTermTranslationPhaseSummary({
          jobId
        } satisfies GetTermTranslationPhaseSummaryRequest),
        this.gateway!.getTermTranslationNextPhaseReadiness({
          jobId
        } satisfies GetTermTranslationNextPhaseReadinessRequest)
      ])

      this.store.update((draft) => {
        draft.phase = "ready"
        draft.summary = summary
        draft.nextPhaseReadiness = nextPhaseReadiness
        draft.pendingAction = null
        draft.errorMessage = ""
        draft.hasLoaded = true
      })
    } catch (error) {
      this.store.update((draft) => {
        draft.phase = "ready"
        draft.pendingAction = null
        draft.summary = before.summary
        draft.nextPhaseReadiness = before.nextPhaseReadiness
        draft.errorMessage = sanitizeErrorMessage(
          error,
          "単語翻訳段階の summary 取得に失敗しました。"
        )
      })
    }
  }

  private async runCommand(
    action: Exclude<TermTranslationPhaseActionKind, "refresh">,
    execute: (
      jobId: number,
      phaseRunId: number
    ) => Promise<TermTranslationPhaseCommandResponse>
  ): Promise<void> {
    const state = this.store.snapshot()
    if (state.jobId === null) {
      this.store.update((draft) => {
        draft.errorMessage = createNoJobSelectedMessage()
      })
      return
    }

    if (!this.gateway) {
      this.store.update((draft) => {
        draft.errorMessage = createGatewayDisconnectedMessage()
      })
      return
    }

    if (action !== "start" && typeof state.summary?.phaseRunId !== "number") {
      this.store.update((draft) => {
        draft.errorMessage = "翻訳段階の実行情報が未確定のため操作できません。"
      })
      return
    }

    this.store.update((draft) => {
      draft.phase = "submitting"
      draft.pendingAction = action
      draft.errorMessage = ""
    })

    try {
      const response =
        action === "start"
          ? await this.gateway.startTermTranslationPhase({
              jobId: state.jobId
            } satisfies StartTermTranslationPhaseRequest)
          : await execute(state.jobId, state.summary!.phaseRunId!)

      this.store.update((draft) => {
        draft.summary = patchSummaryFromCommand(draft.summary, response)
      })

      await this.fetchSummaryAndReadiness(state.jobId, "refresh")
    } catch (error) {
      this.store.update((draft) => {
        draft.phase = "ready"
        draft.pendingAction = null
        draft.errorMessage = sanitizeErrorMessage(
          error,
          "単語翻訳段階の操作に失敗しました。"
        )
      })
    }
  }
}
