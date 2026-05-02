import type {
  BodyTranslationPhaseCommandResponse,
  BodyTranslationPhaseGatewayContract,
  BodyTranslationPhaseSummaryResponse,
  CancelBodyTranslationPhaseRequest,
  GetBodyTranslationOutputReadinessRequest,
  GetBodyTranslationPhaseSummaryRequest,
  PauseBodyTranslationPhaseRequest,
  ResumeBodyTranslationPhaseRequest,
  RetryBodyTranslationPhaseRequest,
  StartBodyTranslationPhaseRequest
} from "@application/gateway-contract/body-translation-phase"

type BodyTranslationPhaseActionKind =
  | "refresh"
  | "start"
  | "pause"
  | "resume"
  | "retry"
  | "cancel"
  | "check-output-readiness"

interface BodyTranslationPhaseScreenState {
  jobId: number | null
  phase: "idle" | "loading" | "ready" | "submitting"
  summary: BodyTranslationPhaseSummaryResponse | null
  outputReadiness: {
    jobId: number
    currentPhase: string
    phaseState: string
    ready: boolean
    blockedReason?: string
    errorKind?: BodyTranslationPhaseSummaryResponse["errorSummary"] extends infer T
      ? T extends { errorKind: infer K }
        ? K
        : never
      : never
    completedFieldCount: number
    statusConsistent: boolean
    outputCount: number
  } | null
  errorMessage: string
  pendingAction: BodyTranslationPhaseActionKind | null
  hasLoaded: boolean
}

interface BodyTranslationPhaseStoreLike {
  snapshot(): BodyTranslationPhaseScreenState
  update(mutator: (draft: BodyTranslationPhaseScreenState) => void): void
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
  return "job id を入力して summary を取得してください。"
}

function createGatewayDisconnectedMessage(): string {
  return "body-translation-phase gateway が未接続です。"
}

function patchSummaryFromCommand(
  currentSummary: BodyTranslationPhaseSummaryResponse | null,
  response: BodyTranslationPhaseCommandResponse
): BodyTranslationPhaseSummaryResponse | null {
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
    inputSummary: response.inputSummary
      ? {
          ...response.inputSummary,
          skippedReasons: response.inputSummary.skippedReasons
            ? [...response.inputSummary.skippedReasons]
            : undefined
        }
      : currentSummary.inputSummary,
    requestSummary: response.requestSummary
      ? { ...response.requestSummary }
      : currentSummary.requestSummary,
    execution: response.execution
      ? { ...response.execution }
      : currentSummary.execution,
    fieldResults: response.fieldResults
      ? response.fieldResults.map((fieldResult) => ({
          ...fieldResult,
          identity: fieldResult.identity
            ? { ...fieldResult.identity }
            : undefined
        }))
      : currentSummary.fieldResults,
    resultSummary: response.resultSummary
      ? {
          ...response.resultSummary,
          fieldResults: response.resultSummary.fieldResults?.map(
            (fieldResult) => ({
              ...fieldResult,
              identity: fieldResult.identity
                ? { ...fieldResult.identity }
                : undefined
            })
          )
        }
      : currentSummary.resultSummary,
    errorSummary: response.errorSummary
      ? { ...response.errorSummary }
      : undefined,
    actionEnablement: {
      ...currentSummary.actionEnablement,
      canRetry: response.retryable
    },
    outputReadiness: { ...response.outputReadiness }
  }
}

export class BodyTranslationPhaseUseCase {
  constructor(
    private readonly gateway: BodyTranslationPhaseGatewayContract | null,
    private readonly store: BodyTranslationPhaseStoreLike
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
      draft.outputReadiness = null
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
      this.gateway!.startBodyTranslationPhase({
        jobId
      } satisfies StartBodyTranslationPhaseRequest)
    )
  }

  async pausePhase(): Promise<void> {
    await this.runCommand("pause", (jobId, phaseRunId) =>
      this.gateway!.pauseBodyTranslationPhase({
        jobId,
        phaseRunId
      } satisfies PauseBodyTranslationPhaseRequest)
    )
  }

  async resumePhase(): Promise<void> {
    await this.runCommand("resume", (jobId, phaseRunId) =>
      this.gateway!.resumeBodyTranslationPhase({
        jobId,
        phaseRunId
      } satisfies ResumeBodyTranslationPhaseRequest)
    )
  }

  async retryPhase(): Promise<void> {
    await this.runCommand("retry", (jobId, phaseRunId) =>
      this.gateway!.retryBodyTranslationPhase({
        jobId,
        phaseRunId
      } satisfies RetryBodyTranslationPhaseRequest)
    )
  }

  async cancelPhase(): Promise<void> {
    await this.runCommand("cancel", (jobId, phaseRunId) =>
      this.gateway!.cancelBodyTranslationPhase({
        jobId,
        phaseRunId
      } satisfies CancelBodyTranslationPhaseRequest)
    )
  }

  async checkOutputReadiness(): Promise<void> {
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

    await this.fetchSummaryAndReadiness(state.jobId, "check-output-readiness")
  }

  private async fetchSummaryAndReadiness(
    jobId: number,
    pendingAction: BodyTranslationPhaseActionKind
  ): Promise<void> {
    const before = this.store.snapshot()

    this.store.update((draft) => {
      draft.phase = "loading"
      draft.pendingAction = pendingAction
      draft.errorMessage = ""
    })

    try {
      const [summary, outputReadiness] = await Promise.all([
        this.gateway!.getBodyTranslationPhaseSummary({
          jobId
        } satisfies GetBodyTranslationPhaseSummaryRequest),
        this.gateway!.getBodyTranslationOutputReadiness({
          jobId
        } satisfies GetBodyTranslationOutputReadinessRequest)
      ])

      this.store.update((draft) => {
        draft.phase = "ready"
        draft.summary = summary
        draft.outputReadiness = outputReadiness
        draft.pendingAction = null
        draft.errorMessage = ""
        draft.hasLoaded = true
      })
    } catch (error) {
      this.store.update((draft) => {
        draft.phase = "ready"
        draft.pendingAction = null
        draft.summary = before.summary
        draft.outputReadiness = before.outputReadiness
        draft.errorMessage = sanitizeErrorMessage(
          error,
          "本文翻訳フェーズ summary の取得に失敗しました。"
        )
      })
    }
  }

  private async runCommand(
    action: Exclude<
      BodyTranslationPhaseActionKind,
      "refresh" | "check-output-readiness"
    >,
    execute: (
      jobId: number,
      phaseRunId: number
    ) => Promise<BodyTranslationPhaseCommandResponse>
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
        draft.errorMessage = "phase run が未確定のため操作できません。"
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
          ? await this.gateway.startBodyTranslationPhase({
              jobId: state.jobId
            } satisfies StartBodyTranslationPhaseRequest)
          : await execute(state.jobId, state.summary!.phaseRunId!)

      this.store.update((draft) => {
        draft.phase = "ready"
        draft.pendingAction = null
        draft.summary = patchSummaryFromCommand(draft.summary, response)
        draft.outputReadiness = {
          jobId: response.jobId,
          currentPhase: response.currentPhase,
          phaseState: response.phaseState,
          ready: response.outputReadiness.ready,
          blockedReason: response.outputReadiness.blockedReason,
          errorKind: response.outputReadiness.errorKind,
          completedFieldCount: response.outputReadiness.completedFieldCount,
          statusConsistent: response.outputReadiness.statusConsistent,
          outputCount: response.execution.outputCount
        }
        draft.errorMessage = ""
        draft.hasLoaded = true
      })
    } catch (error) {
      this.store.update((draft) => {
        draft.phase = "ready"
        draft.pendingAction = null
        draft.errorMessage = sanitizeErrorMessage(
          error,
          "本文翻訳フェーズ操作の反映に失敗しました。"
        )
      })
    }
  }
}
