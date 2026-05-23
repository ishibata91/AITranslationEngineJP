import type {
  CancelPersonaGenerationPhaseRequest,
  GetPersonaGenerationBodyReadinessRequest,
  GetPersonaGenerationPhaseSummaryRequest,
  PausePersonaGenerationPhaseRequest,
  PersonaGenerationBodyReadinessResponse,
  PersonaGenerationPhaseAISettingsRequest,
  PersonaGenerationPhaseCommandResponse,
  PersonaGenerationPhaseGatewayContract,
  PersonaGenerationPhaseSummaryResponse,
  ResumePersonaGenerationPhaseRequest,
  RetryPersonaGenerationPhaseRequest,
  StartPersonaGenerationPhaseRequest
} from "@application/gateway-contract/persona-generation-phase"

type PersonaGenerationPhaseActionKind =
  | "refresh"
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

export class PersonaGenerationPhaseUseCase {
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

    await this.refresh()
  }

  async setJobId(jobId: number | null): Promise<void> {
    this.store.update((draft) => {
      draft.jobId = jobId
      draft.summary = null
      draft.bodyReadiness = null
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
    request: Omit<PersonaGenerationPhaseAISettingsRequest, "jobId">
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
      await this.gateway.savePersonaGenerationPhaseAISettings({
        jobId: state.jobId,
        ...request
      })
      await this.refresh()
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
    pendingAction: PersonaGenerationPhaseActionKind
  ): Promise<void> {
    const before = this.store.snapshot()

    this.store.update((draft) => {
      draft.phase = "loading"
      draft.pendingAction = pendingAction
      draft.errorMessage = ""
    })

    try {
      const [summary, bodyReadiness] = await Promise.all([
        this.gateway!.getPersonaGenerationPhaseSummary({
          jobId
        } satisfies GetPersonaGenerationPhaseSummaryRequest),
        this.gateway!.getPersonaGenerationBodyReadiness({
          jobId
        } satisfies GetPersonaGenerationBodyReadinessRequest)
      ])

      this.store.update((draft) => {
        draft.phase = "ready"
        draft.summary = summary
        draft.bodyReadiness = bodyReadiness
        draft.pendingAction = null
        draft.errorMessage = ""
        draft.hasLoaded = true
      })
    } catch (error) {
      this.store.update((draft) => {
        draft.phase = "ready"
        draft.pendingAction = null
        draft.summary = before.summary
        draft.bodyReadiness = before.bodyReadiness
        draft.errorMessage = sanitizeErrorMessage(
          error,
          "NPC ペルソナ生成段階の summary 取得に失敗しました。"
        )
      })
    }
  }

  private async runCommand(
    action: Exclude<
      PersonaGenerationPhaseActionKind,
      "refresh" | "check-body-readiness" | "start-body-phase"
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

      await this.fetchSummaryAndReadiness(state.jobId, "refresh")
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
