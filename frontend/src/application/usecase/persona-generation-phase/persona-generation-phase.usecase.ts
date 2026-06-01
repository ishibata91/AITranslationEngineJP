import pino from "pino"

import type {
  ProcessingTargetListPageState,
  ProcessingTargetListResponse
} from "@application/gateway-contract/processing-target"

const logger = pino({ browser: { asObject: true } })

import type {
  CancelPersonaGenerationPhaseRequest,
  GetPersonaGenerationPhaseSummaryRequest,
  PersonaGenerationBodyReadinessResponse,
  PersonaGenerationPhaseAISettingsRequest,
  PersonaGenerationPhaseCommandResponse,
  PersonaGenerationPhaseGatewayContract,
  PersonaGenerationPhaseSummaryResponse,
  PausePersonaGenerationPhaseRequest,
  ResumePersonaGenerationPhaseRequest,
  RetryPersonaGenerationPhaseRequest,
  StartPersonaGenerationPhaseRequest
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
  processingTargetPageState?: ProcessingTargetListPageState | null
}

interface PersonaGenerationPhaseStoreLike {
  snapshot(): PersonaGenerationPhaseScreenState
  update(mutator: (draft: PersonaGenerationPhaseScreenState) => void): void
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
  return "NPC ペルソナ生成段階の gateway が未接続です。"
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

export class PersonaGenerationPhaseUseCase {
  private processingTargetListRequestSequence = 0

  constructor(
    private readonly gateway: PersonaGenerationPhaseGatewayContract | null,
    private readonly store: PersonaGenerationPhaseStoreLike
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
        where: "frontend.usecase.persona",
        result: "started",
        id: jobId !== null ? `job:${jobId}` : "job:null",
        seq
      },
      "persona phase setJobId"
    )
    this.store.update((draft) => {
      draft.jobId = jobId
      draft.summary = null
      draft.bodyReadiness = null
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
    await this.fetchProcessingTargetList({ page: 1, searchQuery })
  }

  async setProcessingTargetPage(page: number): Promise<void> {
    await this.fetchProcessingTargetList({ page })
  }

  async startPhase(): Promise<void> {
    await this.runCommand("start", (jobId) =>
      this.gateway!.startPersonaGenerationPhase({
        jobId
      } satisfies StartPersonaGenerationPhaseRequest)
    )
  }

  async pausePhase(): Promise<void> {
    await this.runCommand("pause", (jobId, phaseRunId) =>
      this.gateway!.pausePersonaGenerationPhase({
        jobId,
        phaseRunId
      } satisfies PausePersonaGenerationPhaseRequest)
    )
  }

  async resumePhase(): Promise<void> {
    await this.runCommand("resume", (jobId, phaseRunId) =>
      this.gateway!.resumePersonaGenerationPhase({
        jobId,
        phaseRunId
      } satisfies ResumePersonaGenerationPhaseRequest)
    )
  }

  async retryPhase(): Promise<void> {
    await this.runCommand("retry", (jobId, phaseRunId) =>
      this.gateway!.retryPersonaGenerationPhase({
        jobId,
        phaseRunId
      } satisfies RetryPersonaGenerationPhaseRequest)
    )
  }

  async saveAISettings(
    request: PersonaGenerationPhaseAISettingsRequest
  ): Promise<void> {
    const state = this.store.snapshot()
    if (state.jobId === null) {
      this.store.update((draft) => {
        draft.errorMessage = createNoJobSelectedMessage()
      })
      return
    }
    if (!this.gateway?.savePersonaGenerationPhaseAISettings) {
      this.store.update((draft) => {
        draft.errorMessage = createGatewayDisconnectedMessage()
      })
      return
    }
    try {
      await this.gateway.savePersonaGenerationPhaseAISettings(request)
      await this.fetchSummaryAndReadiness(state.jobId)
    } catch (error) {
      this.store.update((draft) => {
        draft.errorMessage = sanitizeErrorMessage(
          error,
          "NPC ペルソナ生成段階の AI 設定保存に失敗しました。"
        )
      })
    }
  }

  async cancelPhase(): Promise<void> {
    await this.runCommand("cancel", (jobId, phaseRunId) =>
      this.gateway!.cancelPersonaGenerationPhase({
        jobId,
        phaseRunId
      } satisfies CancelPersonaGenerationPhaseRequest)
    )
  }

  async checkBodyReadiness(): Promise<void> {
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

    await this.fetchSummaryAndReadiness(state.jobId, "check-body-readiness")
  }

  startBodyPhase(): Promise<void> {
    this.store.update((draft) => {
      draft.errorMessage =
        "本文翻訳開始の Wails 接続は integration-persona-phase-wails-gateway で実装します。"
    })
    return Promise.resolve()
  }

  private async fetchSummaryAndReadiness(
    jobId: number,
    pendingAction: PersonaGenerationPhaseActionKind | null = null
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
    const summaryPromise = this.gateway.getPersonaGenerationPhaseSummary({
      jobId
    } satisfies GetPersonaGenerationPhaseSummaryRequest)

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
          where: "frontend.usecase.persona",
          result: "reflected",
          id,
          seq,
          elapsedMs: summaryElapsedMs
        },
        "persona summary bridge done and store updated"
      )
    } catch (error) {
      const summaryElapsedMs = Math.round(performance.now() - summaryStart)
      logger.warn(
        {
          event: "bridge_summary_done",
          where: "frontend.usecase.persona",
          result: "failed",
          id,
          seq,
          elapsedMs: summaryElapsedMs
        },
        "persona summary bridge failed"
      )
      this.store.update((draft) => {
        draft.phase = "ready"
        draft.pendingAction = null
        draft.summary = before.summary
        draft.errorMessage = sanitizeErrorMessage(
          error,
          "NPC ペルソナ生成段階の summary 取得に失敗しました。"
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
              where: "frontend.usecase.persona",
              result: "reflected",
              id,
              seq,
              elapsedMs: processingTargetElapsedMs
            },
            "persona processingTarget bridge done and store updated"
          )
        } else {
          logger.warn(
            {
              event: "bridge_processing_target_done",
              where: "frontend.usecase.persona",
              result: "dropped",
              id,
              seq,
              currentSeq: this.processingTargetListRequestSequence,
              elapsedMs: processingTargetElapsedMs,
              reason: "stale_sequence"
            },
            "persona processingTarget dropped by stale sequence"
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
          where: "frontend.usecase.persona",
          result: "completed",
          id,
          seq
        },
        "persona initialFetchDone set to true"
      )
    } catch (error) {
      const processingTargetElapsedMs = Math.round(
        performance.now() - processingTargetStart
      )
      logger.warn(
        {
          event: "bridge_processing_target_done",
          where: "frontend.usecase.persona",
          result: "failed",
          id,
          seq,
          elapsedMs: processingTargetElapsedMs
        },
        "persona processingTarget bridge failed"
      )
      this.store.update((draft) => {
        draft.phase = "ready"
        draft.pendingAction = null
        draft.initialFetchDone = true
        draft.errorMessage = sanitizeErrorMessage(
          error,
          "NPC ペルソナ生成段階の処理対象一覧取得に失敗しました。"
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
        phase: "persona_generation",
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
          "NPC ペルソナ生成段階の処理対象一覧取得に失敗しました。"
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
      phase: "persona_generation",
      page: current.page,
      pageSize: current.pageSize,
      searchQuery: current.searchQuery
    })

    return toProcessingTargetPageState(response)
  }

  private async runCommand(
    action: Exclude<
      PersonaGenerationPhaseActionKind,
      "check-body-readiness" | "start-body-phase"
    >,
    execute: (
      jobId: number,
      phaseRunId: number
    ) => Promise<PersonaGenerationPhaseCommandResponse>
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
      if (action === "start") {
        await this.gateway.startPersonaGenerationPhase({
          jobId: state.jobId
        } satisfies StartPersonaGenerationPhaseRequest)
      } else {
        await execute(state.jobId, state.summary!.phaseRunId!)
      }

      await this.fetchSummaryAndReadiness(state.jobId, action)
    } catch (error) {
      this.store.update((draft) => {
        draft.phase = "ready"
        draft.pendingAction = null
        draft.errorMessage = sanitizeErrorMessage(
          error,
          "NPC ペルソナ生成段階の操作に失敗しました。"
        )
      })
    }
  }
}
