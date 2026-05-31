import pino from "pino"

import type {
  ProcessingTargetListPageState,
  ProcessingTargetListResponse
} from "@application/gateway-contract/processing-target"

import type {
  GetTermTranslationPhaseSummaryRequest,
  PauseTermTranslationPhaseRequest,
  TermTranslationNextPhaseReadinessResponse,
  RetryTermTranslationPhaseRequest,
  ResumeTermTranslationPhaseRequest,
  StartTermTranslationPhaseRequest,
  TermTranslationPhaseAISettingsRequest,
  TermTranslationPhaseCommandResponse,
  TermTranslationPhaseGatewayContract,
  TermTranslationPhaseSummaryResponse
} from "@application/gateway-contract/term-translation-phase"

const logger = pino({ browser: { asObject: true } })

type TermTranslationPhaseActionKind =
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
  initialFetchDone: boolean
  processingTargetPageState?: ProcessingTargetListPageState | null
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
      canRetry: response.retryable
    }
  }
}

export class TermTranslationPhaseUseCase {
  private processingTargetListRequestSequence = 0

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

    await this.fetchSummaryAndReadiness(state.jobId)
  }

  async setJobId(jobId: number | null): Promise<void> {
    this.processingTargetListRequestSequence++
    const seq = this.processingTargetListRequestSequence
    logger.info(
      {
        event: "phase_fetch_start",
        where: "frontend.usecase.term",
        result: "started",
        id: jobId !== null ? `job:${jobId}` : "job:null",
        seq
      },
      "term phase setJobId"
    )
    this.store.update((draft) => {
      draft.jobId = jobId
      draft.summary = null
      draft.nextPhaseReadiness = null
      draft.hasLoaded = false
      draft.initialFetchDone = false
      draft.pendingAction = null
      draft.errorMessage = jobId === null ? createNoJobSelectedMessage() : ""
      draft.phase = "ready"
      draft.processingTargetPageState = createDefaultProcessingTargetPageState()
    })

    if (jobId !== null) {
      await this.fetchSummaryAndReadiness(jobId)
    }
  }

  async setProcessingTargetSearchQuery(searchQuery: string): Promise<void> {
    await this.fetchProcessingTargetList({
      page: 1,
      searchQuery
    })
  }

  async setProcessingTargetPage(page: number): Promise<void> {
    await this.fetchProcessingTargetList({ page })
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

  async saveAISettings(
    request: Omit<TermTranslationPhaseAISettingsRequest, "jobId">
  ): Promise<void> {
    const state = this.store.snapshot()
    if (state.jobId === null) {
      this.store.update((draft) => {
        draft.errorMessage = createNoJobSelectedMessage()
      })
      return
    }
    const gateway = this.gateway
    if (!gateway?.saveTermTranslationPhaseAISettings) {
      this.store.update((draft) => {
        draft.errorMessage = createGatewayDisconnectedMessage()
      })
      return
    }
    try {
      await gateway.saveTermTranslationPhaseAISettings({
        jobId: state.jobId,
        ...request
      })
      await this.fetchSummaryAndReadiness(state.jobId)
    } catch (error) {
      this.store.update((draft) => {
        draft.errorMessage = sanitizeErrorMessage(
          error,
          "単語翻訳段階の AI 設定保存に失敗しました。"
        )
      })
    }
  }

  private async fetchSummaryAndReadiness(
    jobId: number,
    pendingAction: TermTranslationPhaseActionKind | null = null
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
    const processingTargetListRequestSequence =
      ++this.processingTargetListRequestSequence
    const id = `job:${jobId}`
    const seq = processingTargetListRequestSequence

    this.store.update((draft) => {
      draft.phase = "loading"
      draft.pendingAction = pendingAction
      draft.errorMessage = ""
    })

    const summaryStart = performance.now()
    const summaryPromise = this.gateway.getTermTranslationPhaseSummary({
      jobId
    } satisfies GetTermTranslationPhaseSummaryRequest)

    const processingTargetStart = performance.now()
    const processingTargetPromise = this.fetchProcessingTargetListResponse(
      jobId,
      before.processingTargetPageState ?? null
    )

    try {
      const summary = await summaryPromise
      const summaryElapsedMs = Math.round(performance.now() - summaryStart)
      this.store.update((draft) => {
        draft.summary = summary
        draft.errorMessage = ""
      })
      logger.info(
        {
          event: "bridge_summary_done",
          where: "frontend.usecase.term",
          result: "reflected",
          id,
          seq,
          elapsedMs: summaryElapsedMs
        },
        "term summary bridge done and store updated"
      )
    } catch (error) {
      const summaryElapsedMs = Math.round(performance.now() - summaryStart)
      logger.warn(
        {
          event: "bridge_summary_done",
          where: "frontend.usecase.term",
          result: "failed",
          id,
          seq,
          elapsedMs: summaryElapsedMs
        },
        "term summary bridge failed"
      )
      this.store.update((draft) => {
        draft.phase = "ready"
        draft.pendingAction = null
        draft.summary = before.summary
        draft.errorMessage = sanitizeErrorMessage(
          error,
          "単語翻訳段階の summary 取得に失敗しました。"
        )
      })
    }

    try {
      const processingTargetPageState = await processingTargetPromise
      const processingTargetElapsedMs = Math.round(
        performance.now() - processingTargetStart
      )
      this.store.update((draft) => {
        if (
          processingTargetListRequestSequence ===
          this.processingTargetListRequestSequence
        ) {
          draft.processingTargetPageState = processingTargetPageState
          logger.info(
            {
              event: "bridge_processing_target_done",
              where: "frontend.usecase.term",
              result: "reflected",
              id,
              seq,
              elapsedMs: processingTargetElapsedMs
            },
            "term processingTarget bridge done and store updated"
          )
        } else {
          logger.warn(
            {
              event: "bridge_processing_target_done",
              where: "frontend.usecase.term",
              result: "dropped",
              id,
              seq,
              currentSeq: this.processingTargetListRequestSequence,
              elapsedMs: processingTargetElapsedMs,
              reason: "stale_sequence"
            },
            "term processingTarget dropped by stale sequence"
          )
        }
        draft.phase = "ready"
        draft.pendingAction = null
        draft.hasLoaded = true
        draft.initialFetchDone = true
      })
      logger.info(
        {
          event: "initial_fetch_done",
          where: "frontend.usecase.term",
          result: "completed",
          id,
          seq
        },
        "term initialFetchDone set to true"
      )
    } catch (error) {
      const processingTargetElapsedMs = Math.round(
        performance.now() - processingTargetStart
      )
      logger.warn(
        {
          event: "bridge_processing_target_done",
          where: "frontend.usecase.term",
          result: "failed",
          id,
          seq,
          elapsedMs: processingTargetElapsedMs
        },
        "term processingTarget bridge failed"
      )
      this.store.update((draft) => {
        draft.phase = "ready"
        draft.pendingAction = null
        draft.initialFetchDone = true
        draft.errorMessage = sanitizeErrorMessage(
          error,
          "単語翻訳段階の処理対象一覧取得に失敗しました。"
        )
      })
    }
  }

  private async fetchProcessingTargetList(
    next: { page?: number; searchQuery?: string }
  ): Promise<void> {
    const state = this.store.snapshot()
    if (state.jobId === null || !this.gateway?.getProcessingTargetList) {
      return
    }

    const current =
      state.processingTargetPageState ?? createDefaultProcessingTargetPageState()
    const page = Math.max(1, next.page ?? current.page)
    const searchQuery = next.searchQuery ?? current.searchQuery
    const requestSequence = ++this.processingTargetListRequestSequence

    this.store.update((draft) => {
      draft.processingTargetPageState = {
        ...(draft.processingTargetPageState ??
          createDefaultProcessingTargetPageState()),
        page,
        searchQuery,
        busy: true
      }
    })

    try {
      const response = await this.gateway.getProcessingTargetList({
        jobId: state.jobId,
        phase: "term_translation",
        page,
        pageSize: current.pageSize,
        searchQuery
      })

      this.store.update((draft) => {
        if (requestSequence !== this.processingTargetListRequestSequence) {
          return
        }
        draft.processingTargetPageState = toProcessingTargetPageState(response)
      })
    } catch (error) {
      this.store.update((draft) => {
        if (requestSequence !== this.processingTargetListRequestSequence) {
          return
        }
        draft.processingTargetPageState = {
          ...(draft.processingTargetPageState ??
            createDefaultProcessingTargetPageState()),
          busy: false
        }
        draft.errorMessage = sanitizeErrorMessage(
          error,
          "単語翻訳段階の処理対象一覧取得に失敗しました。"
        )
      })
    }
  }

  private async fetchProcessingTargetListResponse(
    jobId: number,
    currentPageState: ProcessingTargetListPageState | null
  ): Promise<ProcessingTargetListPageState> {
    const current = currentPageState ?? createDefaultProcessingTargetPageState()
    if (!this.gateway?.getProcessingTargetList) {
      return current
    }

    const response = await this.gateway.getProcessingTargetList({
      jobId,
      phase: "term_translation",
      page: current.page,
      pageSize: current.pageSize,
      searchQuery: current.searchQuery
    })

    return toProcessingTargetPageState(response)
  }

  private async runCommand(
    action: TermTranslationPhaseActionKind,
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

      await this.fetchSummaryAndReadiness(state.jobId, action)
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
