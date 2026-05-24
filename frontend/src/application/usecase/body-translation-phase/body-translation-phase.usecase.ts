import type {
  ProcessingTargetListPageState,
  ProcessingTargetListPageStatesByPhase,
  ProcessingTargetListResponse
} from "@application/gateway-contract/processing-target"

import type {
  BodyTranslationPhaseAISettingsRequest,
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
  processingTargetPageState?: ProcessingTargetListPageState | null
  processingTargetPageStatesByPhase?: ProcessingTargetListPageStatesByPhase
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
  return "ジョブIDを指定して summary を取得してください。"
}

function createGatewayDisconnectedMessage(): string {
  return "本文翻訳段階の gateway が未接続です。"
}

function createDefaultProcessingTargetPageState(): ProcessingTargetListPageState {
  return {
    items: [],
    metadata: [],
    page: 1,
    pageSize: 50,
    totalCount: 0,
    searchQuery: "",
    busy: false
  }
}

function toProcessingTargetPageState(
  response: ProcessingTargetListResponse
): ProcessingTargetListPageState {
  return {
    ...response,
    busy: false
  }
}

function getProcessingTargetPageState(
  state: BodyTranslationPhaseScreenState,
  phase: string
): ProcessingTargetListPageState {
  return (
    state.processingTargetPageStatesByPhase?.[phase] ??
    (phase === "body_translation" ? state.processingTargetPageState : null) ??
    createDefaultProcessingTargetPageState()
  )
}

function setProcessingTargetPageState(
  draft: BodyTranslationPhaseScreenState,
  phase: string,
  pageState: ProcessingTargetListPageState
): void {
  draft.processingTargetPageStatesByPhase = {
    ...(draft.processingTargetPageStatesByPhase ?? {}),
    [phase]: pageState
  }
  if (phase === "body_translation") {
    draft.processingTargetPageState = pageState
  }
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

    await this.fetchSummaryAndReadiness(state.jobId)
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
      draft.processingTargetPageState = createDefaultProcessingTargetPageState()
      draft.processingTargetPageStatesByPhase = {}
    })

    if (jobId !== null) {
      await this.fetchSummaryAndReadiness(jobId)
    }
  }

  async setProcessingTargetSearchQuery(
    searchQuery: string,
    phase = "body_translation"
  ): Promise<void> {
    await this.fetchProcessingTargetList({ page: 1, searchQuery, phase })
  }

  async setProcessingTargetPage(
    page: number,
    phase = "body_translation"
  ): Promise<void> {
    await this.fetchProcessingTargetList({ page, phase })
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

  async saveAISettings(
    request: Omit<BodyTranslationPhaseAISettingsRequest, "jobId">
  ): Promise<void> {
    const state = this.store.snapshot()
    if (state.jobId === null) {
      this.store.update((draft) => {
        draft.errorMessage = createNoJobSelectedMessage()
      })
      return
    }
    if (!this.gateway?.saveBodyTranslationPhaseAISettings) {
      this.store.update((draft) => {
        draft.errorMessage = createGatewayDisconnectedMessage()
      })
      return
    }
    try {
      await this.gateway.saveBodyTranslationPhaseAISettings({
        jobId: state.jobId,
        ...request
      })
      await this.fetchSummaryAndReadiness(state.jobId)
    } catch (error) {
      this.store.update((draft) => {
        draft.errorMessage = sanitizeErrorMessage(
          error,
          "本文翻訳段階の AI 設定保存に失敗しました。"
        )
      })
    }
  }

  private async fetchSummaryAndReadiness(
    jobId: number,
    pendingAction: BodyTranslationPhaseActionKind | null = null
  ): Promise<void> {
    if (!this.gateway) {
      this.store.update((draft) => {
        draft.phase = "ready"
        draft.pendingAction = null
        draft.errorMessage = createGatewayDisconnectedMessage()
      })
      return
    }

    const before = this.store.snapshot()

    this.store.update((draft) => {
      draft.phase = "loading"
      draft.pendingAction = pendingAction
      draft.errorMessage = ""
    })

    try {
      const [summary, outputReadiness, processingTargetPageState] =
        await Promise.all([
        this.gateway.getBodyTranslationPhaseSummary({
          jobId
        } satisfies GetBodyTranslationPhaseSummaryRequest),
        this.gateway.getBodyTranslationOutputReadiness({
          jobId
        } satisfies GetBodyTranslationOutputReadinessRequest),
        this.fetchProcessingTargetListResponse(
          jobId,
          getProcessingTargetPageState(before, "body_translation"),
          "body_translation"
        )
      ])

      this.store.update((draft) => {
        draft.phase = "ready"
        draft.summary = summary
        draft.outputReadiness = outputReadiness
        setProcessingTargetPageState(
          draft,
          "body_translation",
          processingTargetPageState
        )
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
          "本文翻訳段階の summary 取得に失敗しました。"
        )
      })
    }
  }

  private async fetchProcessingTargetList(next: {
    page?: number
    searchQuery?: string
    phase: string
  }): Promise<void> {
    const state = this.store.snapshot()
    if (state.jobId === null || !this.gateway?.getProcessingTargetList) {
      return
    }

    const current = getProcessingTargetPageState(state, next.phase)
    const page = Math.max(1, next.page ?? current.page)
    const searchQuery = next.searchQuery ?? current.searchQuery

    this.store.update((draft) => {
      setProcessingTargetPageState(draft, next.phase, {
        ...getProcessingTargetPageState(draft, next.phase),
        page,
        searchQuery,
        busy: true
      })
    })

    try {
      const response = await this.gateway.getProcessingTargetList({
        jobId: state.jobId,
        phase: next.phase,
        page,
        pageSize: current.pageSize,
        searchQuery
      })

      this.store.update((draft) => {
        setProcessingTargetPageState(
          draft,
          next.phase,
          toProcessingTargetPageState(response)
        )
      })
    } catch (error) {
      this.store.update((draft) => {
        setProcessingTargetPageState(draft, next.phase, {
          ...getProcessingTargetPageState(draft, next.phase),
          busy: false
        })
        draft.errorMessage = sanitizeErrorMessage(
          error,
          "本文翻訳段階の処理対象一覧取得に失敗しました。"
        )
      })
    }
  }

  private async fetchProcessingTargetListResponse(
    jobId: number,
    currentPageState: ProcessingTargetListPageState | null,
    phase: string
  ): Promise<ProcessingTargetListPageState> {
    const current = currentPageState ?? createDefaultProcessingTargetPageState()
    if (!this.gateway?.getProcessingTargetList) {
      return current
    }

    const response = await this.gateway.getProcessingTargetList({
      jobId,
      phase,
      page: current.page,
      pageSize: current.pageSize,
      searchQuery: current.searchQuery
    })

    return toProcessingTargetPageState(response)
  }

  private async runCommand(
    action: Exclude<
      BodyTranslationPhaseActionKind,
      "check-output-readiness"
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
          "本文翻訳段階の操作反映に失敗しました。"
        )
      })
    }
  }
}
