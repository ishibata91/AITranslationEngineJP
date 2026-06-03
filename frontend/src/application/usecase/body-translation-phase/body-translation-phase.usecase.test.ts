import { describe, expect, test, vi } from "vitest"

import type {
  ProcessingTargetListPageState,
  ProcessingTargetListPageStatesByPhase,
  ProcessingTargetListResponse
} from "@application/gateway-contract/processing-target"
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
  initialFetchDone: boolean
  processingTargetPageState?: ProcessingTargetListPageState | null
  processingTargetPageStatesByPhase?: ProcessingTargetListPageStatesByPhase
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
    projection: {
      phaseLifecycle: "running",
      jobLifecycle: "running",
      errorKind: "none",
      aiSettingsConfigured: true,
      targetCount: 20,
      previousPhaseLifecycle: "completed",
      personaBodyReadiness: {
        bodyReadiness: true,
        snapshotReferenceStatus: "available"
      }
    },
    outputReadiness: {
      completedFieldCount: 8,
      statusConsistent: true,
      outputCount: 8
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
      completedFieldCount: 20,
      statusConsistent: true,
      outputCount: 20
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
    initialFetchDone: true,
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
  const processingTargetResponse: ProcessingTargetListResponse = {
    items: [
      {
        id: "field:1",
        name: "FULL",
        detail: "原文: Hello",
        titleParts: [{ text: "FULL" }],
        metadata: [{ label: "訳語", value: "こんにちは" }]
      }
    ],
    metadata: [],
    page: 4,
    pageSize: 40,
    totalCount: 256,
    searchQuery: "Hello"
  }
  const spies = {
    getBodyTranslationPhaseSummary: vi.fn(),
    getBodyTranslationOutputReadiness: vi.fn(),
    getProcessingTargetList: vi.fn(() =>
      Promise.resolve(processingTargetResponse)
    ),
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

function createDeferredResponse() {
  let resolve!: (response: ProcessingTargetListResponse) => void
  const promise = new Promise<ProcessingTargetListResponse>((resolver) => {
    resolve = resolver
  })

  return { promise, resolve }
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

  test("setJobId number は summary を取得して反映し、outputReadiness 専用取得は行わない", async () => {
    const { gateway, spies } = createGateway()
    const store = createStore(
      createState({ hasLoaded: false, summary: null, outputReadiness: null, initialFetchDone: false })
    )
    const summary = createSummary({ phaseState: "ready" })
    spies.getBodyTranslationPhaseSummary.mockResolvedValue(summary)
    const useCase = new BodyTranslationPhaseUseCase(gateway, store)

    await useCase.setJobId(55)

    expect(spies.getBodyTranslationPhaseSummary).toHaveBeenCalledWith({
      jobId: 55
    })
    expect(spies.getBodyTranslationOutputReadiness).not.toHaveBeenCalled()
    expect(store.getState()).toMatchObject({
      jobId: 55,
      phase: "ready",
      summary,
      pendingAction: null,
      errorMessage: "",
      hasLoaded: true,
      initialFetchDone: true
    })
  })

  test("load は body_translation の processing target request を送り page state に保持する", async () => {
    const { gateway, spies } = createGateway()
    const store = createStore(
      createState({
        hasLoaded: false,
        summary: null,
        outputReadiness: null,
        processingTargetPageState: {
          items: [],
          metadata: [],
          page: 4,
          pageSize: 40,
          totalCount: 0,
          searchQuery: "Hello",
          busy: false
        }
      })
    )
    const summary = createSummary({ phaseState: "ready" })
    spies.getBodyTranslationPhaseSummary.mockResolvedValue(summary)
    const useCase = new BodyTranslationPhaseUseCase(gateway, store)

    await useCase.load()

    expect(spies.getProcessingTargetList).toHaveBeenCalledWith({
      jobId: 9,
      phase: "body_translation",
      page: 4,
      pageSize: 40,
      searchQuery: "Hello"
    })
    expect(store.getState().processingTargetPageState).toMatchObject({
      items: [{ id: "field:1" }],
      page: 4,
      pageSize: 40,
      totalCount: 256,
      searchQuery: "Hello",
      busy: false
    })
    expect(
      store.getState().processingTargetPageStatesByPhase?.["body_translation"]
    ).toMatchObject({
      totalCount: 256,
      searchQuery: "Hello"
    })
  })

  test("検索語変更は body_translation request を page 1 へ戻して送る", async () => {
    const { gateway, spies } = createGateway()
    const store = createStore(
      createState({
        processingTargetPageState: {
          items: [],
          metadata: [],
          page: 6,
          pageSize: 40,
          totalCount: 0,
          searchQuery: "",
          busy: false
        },
        processingTargetPageStatesByPhase: {
          body_translation: {
            items: [],
            metadata: [],
            page: 6,
            pageSize: 40,
            totalCount: 0,
            searchQuery: "",
            busy: false
          }
        }
      })
    )
    const useCase = new BodyTranslationPhaseUseCase(gateway, store)

    await useCase.setProcessingTargetSearchQuery("Quest")

    expect(spies.getProcessingTargetList).toHaveBeenCalledWith({
      jobId: 9,
      phase: "body_translation",
      page: 1,
      pageSize: 40,
      searchQuery: "Quest"
    })
    expect(store.getState().processingTargetPageState).toMatchObject({
      page: 4,
      pageSize: 40,
      totalCount: 256,
      searchQuery: "Hello"
    })
  })

  test("検索応答の到着順が逆転しても最新検索結果だけを page state に反映する", async () => {
    const { gateway, spies } = createGateway()
    const first = createDeferredResponse()
    const second = createDeferredResponse()
    spies.getProcessingTargetList
      .mockImplementationOnce(() => first.promise)
      .mockImplementationOnce(() => second.promise)
    const store = createStore(
      createState({
        processingTargetPageState: {
          items: [],
          metadata: [],
          page: 6,
          pageSize: 40,
          totalCount: 0,
          searchQuery: "",
          busy: false
        },
        processingTargetPageStatesByPhase: {
          body_translation: {
            items: [],
            metadata: [],
            page: 6,
            pageSize: 40,
            totalCount: 0,
            searchQuery: "",
            busy: false
          }
        }
      })
    )
    const useCase = new BodyTranslationPhaseUseCase(gateway, store)

    const firstSearch = useCase.setProcessingTargetSearchQuery("Hello")
    const secondSearch = useCase.setProcessingTargetSearchQuery("Quest")

    second.resolve({
      items: [
        {
          id: "field:quest",
          name: "Quest",
          detail: "原文: Quest",
          titleParts: [{ text: "Quest" }],
          metadata: []
        }
      ],
      metadata: [],
      page: 1,
      pageSize: 40,
      totalCount: 3,
      searchQuery: "Quest"
    })
    await secondSearch

    first.resolve({
      items: [
        {
          id: "field:hello",
          name: "Hello",
          detail: "原文: Hello",
          titleParts: [{ text: "Hello" }],
          metadata: []
        }
      ],
      metadata: [],
      page: 1,
      pageSize: 40,
      totalCount: 11,
      searchQuery: "Hello"
    })
    await firstSearch

    expect(store.getState().processingTargetPageState).toMatchObject({
      items: [{ id: "field:quest" }],
      page: 1,
      pageSize: 40,
      totalCount: 3,
      searchQuery: "Quest",
      busy: false
    })
    expect(
      store.getState().processingTargetPageStatesByPhase?.["body_translation"]
    ).toMatchObject({
      totalCount: 3,
      searchQuery: "Quest"
    })
  })

  test("load の一覧応答が検索応答より遅れても最新検索結果だけを page state に反映する", async () => {
    const { gateway, spies } = createGateway()
    const loadResponse = createDeferredResponse()
    const searchResponse = createDeferredResponse()
    spies.getProcessingTargetList
      .mockImplementationOnce(() => loadResponse.promise)
      .mockImplementationOnce(() => searchResponse.promise)
    spies.getBodyTranslationPhaseSummary.mockResolvedValue(
      createSummary({ phaseState: "ready" })
    )
    const store = createStore(
      createState({
        processingTargetPageState: {
          items: [],
          metadata: [],
          page: 6,
          pageSize: 40,
          totalCount: 0,
          searchQuery: "",
          busy: false
        },
        processingTargetPageStatesByPhase: {
          body_translation: {
            items: [],
            metadata: [],
            page: 6,
            pageSize: 40,
            totalCount: 0,
            searchQuery: "",
            busy: false
          }
        }
      })
    )
    const useCase = new BodyTranslationPhaseUseCase(gateway, store)

    const load = useCase.load()
    const search = useCase.setProcessingTargetSearchQuery("Quest")

    searchResponse.resolve({
      items: [
        {
          id: "field:quest",
          name: "Quest",
          detail: "原文: Quest",
          titleParts: [{ text: "Quest" }],
          metadata: []
        }
      ],
      metadata: [],
      page: 1,
      pageSize: 40,
      totalCount: 3,
      searchQuery: "Quest"
    })
    await search

    loadResponse.resolve({
      items: [
        {
          id: "field:old",
          name: "Old",
          detail: "原文: Old",
          titleParts: [{ text: "Old" }],
          metadata: []
        }
      ],
      metadata: [],
      page: 6,
      pageSize: 40,
      totalCount: 99,
      searchQuery: ""
    })
    await load

    expect(store.getState().processingTargetPageState).toMatchObject({
      items: [{ id: "field:quest" }],
      page: 1,
      pageSize: 40,
      totalCount: 3,
      searchQuery: "Quest",
      busy: false
    })
    expect(
      store.getState().processingTargetPageStatesByPhase?.["body_translation"]
    ).toMatchObject({
      totalCount: 3,
      searchQuery: "Quest"
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

  test("summary 取得失敗時は既存 summary を保持して失敗文言を設定する", async () => {
    const { gateway, spies } = createGateway()
    const currentSummary = createSummary({ phaseState: "running" })
    const store = createStore(
      createState({
        summary: currentSummary,
        outputReadiness: null
      })
    )
    spies.getBodyTranslationPhaseSummary.mockRejectedValue(new Error("timeout"))
    const useCase = new BodyTranslationPhaseUseCase(gateway, store)

    await useCase.load()

    expect(store.getState().summary).toEqual(currentSummary)
    expect(store.getState().outputReadiness).toBeNull()
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
    const summary = store.getState().summary
    expect(summary?.phaseState).toBe("running")
    expect(summary?.phaseRunId).toBe(command.phaseRunId)
    expect(typeof summary?.projection.phaseLifecycle).toBe("string")
    expect(store.getState().outputReadiness).toBeNull()
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
      expect(store.getState().outputReadiness).toBeNull()
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
