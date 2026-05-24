import { describe, expect, test, vi } from "vitest"

import type {
  BodyTranslationOutputReadinessResponse,
  BodyTranslationPhaseCommandResponse,
  BodyTranslationPhaseGatewayContract,
  BodyTranslationPhaseSummaryResponse
} from "@application/gateway-contract/body-translation-phase"

import { BodyTranslationPhaseUseCase } from "./body-translation-phase.usecase"

type ScreenState = {
  jobId: number | null
  phase: "idle" | "loading" | "ready" | "submitting"
  summary: BodyTranslationPhaseSummaryResponse | null
  outputReadiness: BodyTranslationOutputReadinessResponse | null
  errorMessage: string
  pendingAction:
    | "start"
    | "pause"
    | "resume"
    | "retry"
    | "cancel"
    | "check-output-readiness"
    | null
  hasLoaded: boolean
}

type StoreLike = {
  snapshot: () => ScreenState
  update: (mutator: (draft: ScreenState) => void) => void
  getState: () => ScreenState
}

function createSummary(
  overrides: Partial<BodyTranslationPhaseSummaryResponse> = {}
): BodyTranslationPhaseSummaryResponse {
  return {
    jobId: 9,
    currentPhase: "body_translation",
    phaseState: "running",
    phaseRunId: 41,
    startedAt: "2026-05-02T01:00:00Z",
    progress: {
      percent: 45,
      processedCount: 9,
      totalCount: 20,
      targetCount: 20,
      translatedCount: 8,
      skippedCount: 1,
      currentStep: "provider_request"
    },
    inputSummary: {
      targetCount: 20,
      dictionaryDigest: "sha256:dictionary",
      personaDigest: "sha256:persona",
      metadataDigest: "sha256:metadata",
      promptDigest: "sha256:prompt"
    },
    execution: {
      credentialRef: "cred-1",
      provider: "openai-compatible",
      model: "gpt-4.1-mini",
      executionMode: "batch",
      requestUnitCount: 20,
      outputCount: 8
    },
    actionEnablement: {
      canStart: false,
      canPause: true,
      canResume: false,
      canRetry: false,
      canCancel: true,
      canCheckOutputReadiness: true
    },
    outputReadiness: {
      ready: false,
      blockedReason: "phase_running",
      completedFieldCount: 8,
      statusConsistent: true
    },
    ...overrides
  }
}

function createOutputReadiness(
  overrides: Partial<BodyTranslationOutputReadinessResponse> = {}
): BodyTranslationOutputReadinessResponse {
  return {
    jobId: 9,
    currentPhase: "body_translation",
    phaseState: "running",
    ready: false,
    blockedReason: "phase_running",
    completedFieldCount: 8,
    statusConsistent: true,
    outputCount: 8,
    ...overrides
  }
}

function createCommandResponse(
  overrides: Partial<BodyTranslationPhaseCommandResponse> = {}
): BodyTranslationPhaseCommandResponse {
  return {
    jobId: 9,
    currentPhase: "body_translation",
    phaseState: "paused",
    phaseRunId: 42,
    startedAt: "2026-05-02T01:00:00Z",
    finishedAt: "2026-05-02T01:10:00Z",
    progress: {
      percent: 100,
      processedCount: 20,
      totalCount: 20,
      targetCount: 20,
      translatedCount: 18,
      skippedCount: 2,
      currentStep: "completed"
    },
    inputSummary: {
      targetCount: 20,
      dictionaryDigest: "sha256:new-dictionary",
      personaDigest: "sha256:new-persona",
      metadataDigest: "sha256:new-metadata",
      promptDigest: "sha256:new-prompt"
    },
    execution: {
      credentialRef: "cred-2",
      provider: "openai-compatible",
      model: "gpt-4.1",
      executionMode: "batch",
      requestUnitCount: 20,
      outputCount: 18
    },
    retryable: true,
    outputReadiness: {
      ready: true,
      blockedReason: "",
      completedFieldCount: 20,
      statusConsistent: true
    },
    ...overrides
  }
}

function createState(overrides: Partial<ScreenState> = {}): ScreenState {
  return {
    jobId: 9,
    phase: "ready",
    summary: createSummary(),
    outputReadiness: createOutputReadiness(),
    errorMessage: "",
    pendingAction: null,
    hasLoaded: true,
    ...overrides
  }
}

function createStore(initial: ScreenState): StoreLike {
  let current = structuredClone(initial)

  return {
    snapshot: () => structuredClone(current),
    update: (mutator) => {
      const draft = structuredClone(current)
      mutator(draft)
      current = draft
    },
    getState: () => structuredClone(current)
  }
}

function createGateway() {
  const spies = {
    getBodyTranslationPhaseSummary: vi.fn(),
    getBodyTranslationOutputReadiness: vi.fn(),
    startBodyTranslationPhase: vi.fn(),
    pauseBodyTranslationPhase: vi.fn(),
    resumeBodyTranslationPhase: vi.fn(),
    retryBodyTranslationPhase: vi.fn(),
    cancelBodyTranslationPhase: vi.fn()
  }

  return {
    gateway: spies as BodyTranslationPhaseGatewayContract,
    spies
  }
}

describe("BodyTranslationPhaseUseCase", () => {
  test("load はジョブ未選択時に未選択エラーを設定する", async () => {
    const { gateway } = createGateway()
    const store = createStore(createState({ jobId: null }))
    const useCase = new BodyTranslationPhaseUseCase(gateway, store)

    await useCase.load()

    expect(store.getState().phase).toBe("ready")
    expect(store.getState().errorMessage).toBe(
      "ジョブIDを指定して summary を取得してください。"
    )
  })

  test("setJobId null は状態を初期化し未選択エラーを設定する", async () => {
    const { gateway, spies } = createGateway()
    const store = createStore(createState())
    const useCase = new BodyTranslationPhaseUseCase(gateway, store)

    await useCase.setJobId(null)

    expect(spies.getBodyTranslationPhaseSummary).not.toHaveBeenCalled()
    expect(store.getState()).toMatchObject({
      jobId: null,
      summary: null,
      outputReadiness: null,
      pendingAction: null,
      hasLoaded: false,
      phase: "ready",
      errorMessage: "ジョブIDを指定して summary を取得してください。"
    })
  })

  test("setJobId number は summary と output readiness を取得して反映する", async () => {
    const { gateway, spies } = createGateway()
    const store = createStore(
      createState({ hasLoaded: false, summary: null, outputReadiness: null })
    )
    const summary = createSummary({ phaseState: "ready" })
    const readiness = createOutputReadiness({ ready: true, outputCount: 20 })
    spies.getBodyTranslationPhaseSummary.mockResolvedValue(summary)
    spies.getBodyTranslationOutputReadiness.mockResolvedValue(readiness)
    const useCase = new BodyTranslationPhaseUseCase(gateway, store)

    await useCase.setJobId(55)

    expect(spies.getBodyTranslationPhaseSummary).toHaveBeenCalledWith({
      jobId: 55
    })
    expect(spies.getBodyTranslationOutputReadiness).toHaveBeenCalledWith({
      jobId: 55
    })
    expect(store.getState()).toMatchObject({
      jobId: 55,
      phase: "ready",
      summary,
      outputReadiness: readiness,
      pendingAction: null,
      errorMessage: "",
      hasLoaded: true
    })
  })

  test("load は gateway 未接続時にエラーを設定する", async () => {
    const store = createStore(createState())
    const useCase = new BodyTranslationPhaseUseCase(null, store)

    await useCase.load()

    expect(store.getState().errorMessage).toBe(
      "本文翻訳段階の gateway が未接続です。"
    )
    expect(store.getState().phase).toBe("ready")
  })

  test("summary 取得失敗時は既存 summary と readiness を保持して失敗文言を設定する", async () => {
    const { gateway, spies } = createGateway()
    const currentSummary = createSummary({ phaseState: "running" })
    const currentReadiness = createOutputReadiness({ ready: false })
    const store = createStore(
      createState({
        summary: currentSummary,
        outputReadiness: currentReadiness
      })
    )
    spies.getBodyTranslationPhaseSummary.mockRejectedValue(new Error("timeout"))
    spies.getBodyTranslationOutputReadiness.mockResolvedValue(
      createOutputReadiness({ ready: true })
    )
    const useCase = new BodyTranslationPhaseUseCase(gateway, store)

    await useCase.load()

    expect(store.getState().summary).toEqual(currentSummary)
    expect(store.getState().outputReadiness).toEqual(currentReadiness)
    expect(store.getState().errorMessage).toBe(
      "本文翻訳段階の summary 取得に失敗しました。"
    )
    expect(store.getState().pendingAction).toBeNull()
  })

  test("checkOutputReadiness 取得失敗で Wails 未接続エラー文を保持する", async () => {
    const { gateway, spies } = createGateway()
    const store = createStore(createState())
    spies.getBodyTranslationPhaseSummary.mockRejectedValue(
      new Error(
        "Wails binding is not wired yet: GetBodyTranslationPhaseSummary"
      )
    )
    spies.getBodyTranslationOutputReadiness.mockResolvedValue(
      createOutputReadiness({ ready: true })
    )
    const useCase = new BodyTranslationPhaseUseCase(gateway, store)

    await useCase.checkOutputReadiness()

    expect(store.getState().errorMessage).toBe(
      "Wails binding is not wired yet: GetBodyTranslationPhaseSummary"
    )
  })

  test("startPhase は jobId を渡して summary と readiness を更新する", async () => {
    const { gateway, spies } = createGateway()
    const store = createStore(createState())
    const command = createCommandResponse({
      phaseState: "running",
      retryable: false
    })
    spies.startBodyTranslationPhase.mockResolvedValue(command)
    const useCase = new BodyTranslationPhaseUseCase(gateway, store)

    await useCase.startPhase()

    expect(spies.startBodyTranslationPhase).toHaveBeenCalledWith({ jobId: 9 })
    expect(store.getState().summary).toMatchObject({
      phaseState: "running",
      phaseRunId: command.phaseRunId,
      actionEnablement: {
        canRetry: false
      }
    })
    expect(store.getState().outputReadiness).toMatchObject({
      jobId: command.jobId,
      ready: command.outputReadiness.ready,
      outputCount: command.execution.outputCount
    })
  })

  test.each([
    ["pausePhase", "pauseBodyTranslationPhase"],
    ["resumePhase", "resumeBodyTranslationPhase"],
    ["retryPhase", "retryBodyTranslationPhase"],
    ["cancelPhase", "cancelBodyTranslationPhase"]
  ] as const)(
    "%s は jobId と phaseRunId を gateway に渡す",
    async (method, gatewayMethod) => {
      const { gateway, spies } = createGateway()
      const store = createStore(
        createState({ jobId: 77, summary: createSummary({ phaseRunId: 888 }) })
      )
      const command = createCommandResponse({ jobId: 77, phaseRunId: 888 })
      spies[gatewayMethod].mockResolvedValue(command)
      const useCase = new BodyTranslationPhaseUseCase(gateway, store)

      await useCase[method]()

      expect(spies[gatewayMethod]).toHaveBeenCalledWith({
        jobId: 77,
        phaseRunId: 888
      })
      expect(store.getState().summary?.phaseRunId).toBe(888)
      expect(store.getState().outputReadiness?.outputCount).toBe(
        command.execution.outputCount
      )
    }
  )

  test("pausePhase は phaseRunId 未確定時に操作を拒否する", async () => {
    const { gateway, spies } = createGateway()
    const store = createStore(
      createState({ summary: createSummary({ phaseRunId: undefined }) })
    )
    const useCase = new BodyTranslationPhaseUseCase(gateway, store)

    await useCase.pausePhase()

    expect(spies.pauseBodyTranslationPhase).not.toHaveBeenCalled()
    expect(store.getState().errorMessage).toBe(
      "翻訳段階の実行情報が未確定のため操作できません。"
    )
  })

  test("resumePhase gateway 失敗時は操作失敗文言を設定する", async () => {
    const { gateway, spies } = createGateway()
    const store = createStore(createState())
    spies.resumeBodyTranslationPhase.mockRejectedValue(
      new Error("upstream failed")
    )
    const useCase = new BodyTranslationPhaseUseCase(gateway, store)

    await useCase.resumePhase()

    expect(store.getState().errorMessage).toBe(
      "本文翻訳段階の操作反映に失敗しました。"
    )
    expect(store.getState().phase).toBe("ready")
    expect(store.getState().pendingAction).toBeNull()
  })

  test("retryPhase gateway 失敗時は Wails 未接続エラー文を保持する", async () => {
    const { gateway, spies } = createGateway()
    const store = createStore(createState())
    spies.retryBodyTranslationPhase.mockRejectedValue(
      new Error("Wails binding is not wired yet: RetryBodyTranslationPhase")
    )
    const useCase = new BodyTranslationPhaseUseCase(gateway, store)

    await useCase.retryPhase()

    expect(store.getState().errorMessage).toBe(
      "Wails binding is not wired yet: RetryBodyTranslationPhase"
    )
  })

  test("cancelPhase はジョブ未選択時に未選択エラーを設定する", async () => {
    const { gateway, spies } = createGateway()
    const store = createStore(createState({ jobId: null }))
    const useCase = new BodyTranslationPhaseUseCase(gateway, store)

    await useCase.cancelPhase()

    expect(spies.cancelBodyTranslationPhase).not.toHaveBeenCalled()
    expect(store.getState().errorMessage).toBe(
      "ジョブIDを指定して summary を取得してください。"
    )
  })
})
