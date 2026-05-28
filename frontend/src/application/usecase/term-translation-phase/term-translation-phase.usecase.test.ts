import { describe, expect, test, vi } from "vitest"

import type {
  ProcessingTargetListPageState,
  ProcessingTargetListResponse
} from "@application/gateway-contract/processing-target"
import type {
  TermTranslationNextPhaseReadinessResponse,
  TermTranslationPhaseGatewayContract,
  TermTranslationPhaseSummaryResponse
} from "@application/gateway-contract/term-translation-phase"

import { TermTranslationPhaseUseCase } from "./term-translation-phase.usecase"

function createSummary(
  overrides: Partial<TermTranslationPhaseSummaryResponse> = {}
): TermTranslationPhaseSummaryResponse {
  return {
    jobId: 9,
    currentPhase: "term_translation",
    phaseState: "ready",
    progress: {
      percent: 0,
      processedCount: 0,
      totalCount: 10,
      aiTargetCount: 7,
      currentStep: "ready"
    },
    totalTermCount: 10,
    dictionaryHitCount: 3,
    aiTargetCount: 7,
    execution: {
      credentialRef: "cred",
      provider: "openai-compatible",
      model: "gpt-4.1-mini",
      executionMode: "batch"
    },
    actionEnablement: {
      canStart: true,
      canPause: false,
      canResume: false,
      canRetry: false,
      canStartNextPhase: false
    },
    ...overrides
  }
}

function createReadiness(
  overrides: Partial<TermTranslationNextPhaseReadinessResponse> = {}
): TermTranslationNextPhaseReadinessResponse {
  return {
    jobId: 9,
    currentPhase: "term_translation",
    phaseState: "ready",
    canStartNextPhase: false,
    ...overrides
  }
}

interface TermTranslationPhaseScreenStateLike {
  jobId: number | null
  phase: "idle" | "loading" | "ready" | "submitting"
  summary: TermTranslationPhaseSummaryResponse | null
  nextPhaseReadiness: TermTranslationNextPhaseReadinessResponse | null
  errorMessage: string
  pendingAction: "start" | "pause" | "resume" | "retry" | null
  hasLoaded: boolean
  processingTargetPageState?: ProcessingTargetListPageState | null
}

function cloneState(
  state: TermTranslationPhaseScreenStateLike
): TermTranslationPhaseScreenStateLike {
  return {
    ...state,
    summary: state.summary
      ? {
          ...state.summary,
          progress: { ...state.summary.progress },
          execution: { ...state.summary.execution },
          resultSummary: state.summary.resultSummary
            ? { ...state.summary.resultSummary }
            : undefined,
          errorSummary: state.summary.errorSummary
            ? { ...state.summary.errorSummary }
            : undefined,
          actionEnablement: { ...state.summary.actionEnablement }
        }
      : null,
    nextPhaseReadiness: state.nextPhaseReadiness
      ? { ...state.nextPhaseReadiness }
      : null,
    processingTargetPageState: state.processingTargetPageState
      ? structuredClone(state.processingTargetPageState)
      : state.processingTargetPageState
  }
}

function createStore(
  initialState: Partial<TermTranslationPhaseScreenStateLike> = {}
) {
  let state: TermTranslationPhaseScreenStateLike = {
    jobId: null,
    phase: "idle",
    summary: null,
    nextPhaseReadiness: null,
    errorMessage: "",
    pendingAction: null,
    hasLoaded: false,
    ...initialState
  }

  return {
    snapshot() {
      return cloneState(state)
    },
    update(mutator: (draft: TermTranslationPhaseScreenStateLike) => void) {
      const draft = cloneState(state)
      mutator(draft)
      state = draft
    }
  }
}

function createGatewaySpies() {
  const processingTargetResponse: ProcessingTargetListResponse = {
    items: [
      {
        id: "term:1",
        name: "Dragon",
        detail: "原文: Dragon",
        titleParts: [{ text: "Dragon" }],
        metadata: [{ label: "候補", value: "ドラゴン" }]
      }
    ],
    metadata: [],
    page: 2,
    pageSize: 50,
    totalCount: 137,
    searchQuery: "Dragon"
  }
  const spies = {
    getTermTranslationPhaseSummary: vi.fn(() =>
      Promise.resolve(createSummary())
    ),
    getTermTranslationNextPhaseReadiness: vi.fn(() =>
      Promise.resolve(createReadiness())
    ),
    getProcessingTargetList: vi.fn(() =>
      Promise.resolve(processingTargetResponse)
    ),
    startTermTranslationPhase: vi.fn(() =>
      Promise.resolve({
        jobId: 9,
        currentPhase: "term_translation",
        phaseState: "running",
        phaseRunId: 81,
        progress: {
          percent: 0,
          processedCount: 0,
          totalCount: 10,
          aiTargetCount: 7,
          currentStep: "starting"
        },
        retryable: false,
        canStartNextPhase: false
      })
    ),
    pauseTermTranslationPhase: vi.fn(() =>
      Promise.resolve({
        jobId: 9,
        currentPhase: "term_translation",
        phaseState: "paused",
        phaseRunId: 81,
        progress: {
          percent: 10,
          processedCount: 1,
          totalCount: 10,
          aiTargetCount: 7,
          currentStep: "paused"
        },
        retryable: true,
        canStartNextPhase: false
      })
    ),
    resumeTermTranslationPhase: vi.fn(() =>
      Promise.resolve({
        jobId: 9,
        currentPhase: "term_translation",
        phaseState: "running",
        phaseRunId: 81,
        progress: {
          percent: 20,
          processedCount: 2,
          totalCount: 10,
          aiTargetCount: 7,
          currentStep: "resumed"
        },
        retryable: false,
        canStartNextPhase: false
      })
    ),
    retryTermTranslationPhase: vi.fn(() =>
      Promise.resolve({
        jobId: 9,
        currentPhase: "term_translation",
        phaseState: "running",
        phaseRunId: 81,
        progress: {
          percent: 30,
          processedCount: 3,
          totalCount: 10,
          aiTargetCount: 7,
          currentStep: "retry"
        },
        retryable: false,
        canStartNextPhase: false
      })
    )
  }

  return {
    gateway: spies as TermTranslationPhaseGatewayContract,
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

describe("TermTranslationPhaseUseCase", () => {
  test("gateway 未接続で load は接続エラーメッセージを設定する", async () => {
    const store = createStore({ jobId: 9, phase: "ready" })
    const useCase = new TermTranslationPhaseUseCase(null, store)

    await useCase.load()

    expect(store.snapshot().errorMessage).toBe(
      "単語翻訳段階の gateway が未接続です。"
    )
  })

  test("setJobId は summary を再取得し hasLoaded を true に更新する", async () => {
    const { gateway, spies } = createGatewaySpies()
    const store = createStore()
    const useCase = new TermTranslationPhaseUseCase(gateway, store)

    await useCase.setJobId(9)

    expect(spies.getTermTranslationPhaseSummary).toHaveBeenCalledWith({
      jobId: 9
    })
    expect(spies.getTermTranslationNextPhaseReadiness).toHaveBeenCalledWith({
      jobId: 9
    })
    expect(store.snapshot().hasLoaded).toBe(true)
  })

  test("load は processing target request を送り totalCount と searchQuery を page state に保持する", async () => {
    const { gateway, spies } = createGatewaySpies()
    const store = createStore({
      jobId: 9,
      processingTargetPageState: {
        items: [],
        metadata: [],
        page: 2,
        pageSize: 50,
        totalCount: 0,
        searchQuery: "Dragon",
        busy: false
      }
    })
    const useCase = new TermTranslationPhaseUseCase(gateway, store)

    await useCase.load()

    expect(spies.getProcessingTargetList).toHaveBeenCalledWith({
      jobId: 9,
      phase: "term_translation",
      page: 2,
      pageSize: 50,
      searchQuery: "Dragon"
    })
    expect(store.snapshot().processingTargetPageState).toMatchObject({
      items: [{ id: "term:1" }],
      page: 2,
      pageSize: 50,
      totalCount: 137,
      searchQuery: "Dragon",
      busy: false
    })
  })

  test("検索語変更は term_translation request を page 1 へ戻して送る", async () => {
    const { gateway, spies } = createGatewaySpies()
    const store = createStore({
      jobId: 9,
      processingTargetPageState: {
        items: [],
        metadata: [],
        page: 4,
        pageSize: 50,
        totalCount: 0,
        searchQuery: "",
        busy: false
      }
    })
    const useCase = new TermTranslationPhaseUseCase(gateway, store)

    await useCase.setProcessingTargetSearchQuery("Guard")

    expect(spies.getProcessingTargetList).toHaveBeenCalledWith({
      jobId: 9,
      phase: "term_translation",
      page: 1,
      pageSize: 50,
      searchQuery: "Guard"
    })
    expect(store.snapshot().processingTargetPageState).toMatchObject({
      page: 2,
      pageSize: 50,
      searchQuery: "Dragon",
      totalCount: 137
    })
  })

  test("検索応答の到着順が逆転しても最新検索結果だけを page state に反映する", async () => {
    const { gateway, spies } = createGatewaySpies()
    const first = createDeferredResponse()
    const second = createDeferredResponse()
    spies.getProcessingTargetList
      .mockImplementationOnce(() => first.promise)
      .mockImplementationOnce(() => second.promise)
    const store = createStore({
      jobId: 9,
      processingTargetPageState: {
        items: [],
        metadata: [],
        page: 4,
        pageSize: 50,
        totalCount: 0,
        searchQuery: "",
        busy: false
      }
    })
    const useCase = new TermTranslationPhaseUseCase(gateway, store)

    const firstSearch = useCase.setProcessingTargetSearchQuery("Dragon")
    const secondSearch = useCase.setProcessingTargetSearchQuery("Guard")

    second.resolve({
      items: [
        {
          id: "term:guard",
          name: "Guard",
          detail: "原文: Guard",
          titleParts: [{ text: "Guard" }],
          metadata: []
        }
      ],
      metadata: [],
      page: 1,
      pageSize: 50,
      totalCount: 1,
      searchQuery: "Guard"
    })
    await secondSearch

    first.resolve({
      items: [
        {
          id: "term:dragon",
          name: "Dragon",
          detail: "原文: Dragon",
          titleParts: [{ text: "Dragon" }],
          metadata: []
        }
      ],
      metadata: [],
      page: 1,
      pageSize: 50,
      totalCount: 9,
      searchQuery: "Dragon"
    })
    await firstSearch

    expect(store.snapshot().processingTargetPageState).toMatchObject({
      items: [{ id: "term:guard" }],
      page: 1,
      pageSize: 50,
      totalCount: 1,
      searchQuery: "Guard",
      busy: false
    })
  })

  test("load の一覧応答が検索応答より遅れても最新検索結果だけを page state に反映する", async () => {
    const { gateway, spies } = createGatewaySpies()
    const loadResponse = createDeferredResponse()
    const searchResponse = createDeferredResponse()
    spies.getProcessingTargetList
      .mockImplementationOnce(() => loadResponse.promise)
      .mockImplementationOnce(() => searchResponse.promise)
    const store = createStore({
      jobId: 9,
      processingTargetPageState: {
        items: [],
        metadata: [],
        page: 4,
        pageSize: 50,
        totalCount: 0,
        searchQuery: "",
        busy: false
      }
    })
    const useCase = new TermTranslationPhaseUseCase(gateway, store)

    const load = useCase.load()
    const search = useCase.setProcessingTargetSearchQuery("Guard")

    searchResponse.resolve({
      items: [
        {
          id: "term:guard",
          name: "Guard",
          detail: "原文: Guard",
          titleParts: [{ text: "Guard" }],
          metadata: []
        }
      ],
      metadata: [],
      page: 1,
      pageSize: 50,
      totalCount: 1,
      searchQuery: "Guard"
    })
    await search

    loadResponse.resolve({
      items: [
        {
          id: "term:old",
          name: "Old",
          detail: "原文: Old",
          titleParts: [{ text: "Old" }],
          metadata: []
        }
      ],
      metadata: [],
      page: 4,
      pageSize: 50,
      totalCount: 99,
      searchQuery: ""
    })
    await load

    expect(store.snapshot().processingTargetPageState).toMatchObject({
      items: [{ id: "term:guard" }],
      page: 1,
      pageSize: 50,
      totalCount: 1,
      searchQuery: "Guard",
      busy: false
    })
  })

  test("setJobId null は phase readiness を取得せず未選択エラーを設定する", async () => {
    const { gateway, spies } = createGatewaySpies()
    const store = createStore({
      jobId: 9,
      phase: "ready",
      summary: createSummary(),
      nextPhaseReadiness: createReadiness(),
      hasLoaded: true
    })
    const useCase = new TermTranslationPhaseUseCase(gateway, store)

    await useCase.setJobId(null)

    expect(spies.getTermTranslationPhaseSummary).not.toHaveBeenCalled()
    expect(spies.getTermTranslationNextPhaseReadiness).not.toHaveBeenCalled()
    expect(store.snapshot()).toMatchObject({
      jobId: null,
      summary: null,
      nextPhaseReadiness: null,
      hasLoaded: false,
      errorMessage: "ジョブIDを指定して summary を取得してください。"
    })
  })

  test("summary 取得失敗時は前回 snapshot を保持する", async () => {
    const { gateway, spies } = createGatewaySpies()
    spies.getTermTranslationPhaseSummary.mockRejectedValueOnce(
      new Error("timeout")
    )
    const oldSummary = createSummary({ phaseState: "paused", phaseRunId: 44 })
    const oldReadiness = createReadiness({ blockedReason: "old" })
    const store = createStore({
      jobId: 9,
      phase: "ready",
      summary: oldSummary,
      nextPhaseReadiness: oldReadiness,
      hasLoaded: true
    })
    const useCase = new TermTranslationPhaseUseCase(gateway, store)

    await useCase.load()

    const snapshot = store.snapshot()
    expect(snapshot.summary).toMatchObject(oldSummary)
    expect(snapshot.nextPhaseReadiness).toMatchObject(oldReadiness)
    expect(snapshot.errorMessage).toBe(
      "単語翻訳段階の summary 取得に失敗しました。"
    )
  })

  test("start 成功時は start command 実行後に summary を再取得する", async () => {
    const { gateway, spies } = createGatewaySpies()
    const store = createStore({
      jobId: 9,
      phase: "ready",
      summary: createSummary()
    })
    const useCase = new TermTranslationPhaseUseCase(gateway, store)

    await useCase.startPhase()

    expect(spies.startTermTranslationPhase).toHaveBeenCalledWith({ jobId: 9 })
    expect(spies.getTermTranslationPhaseSummary).toHaveBeenCalledTimes(1)
    expect(store.snapshot().pendingAction).toBeNull()
  })

  test("pause は phaseRunId 未確定の時に command を呼ばずエラーを返す", async () => {
    const { gateway, spies } = createGatewaySpies()
    const store = createStore({
      jobId: 9,
      phase: "ready",
      summary: createSummary()
    })
    const useCase = new TermTranslationPhaseUseCase(gateway, store)

    await useCase.pausePhase()

    expect(spies.pauseTermTranslationPhase).not.toHaveBeenCalled()
    expect(store.snapshot().errorMessage).toBe(
      "翻訳段階の実行情報が未確定のため操作できません。"
    )
  })

  test("resume 失敗時に Wails 未接続エラー文を優先する", async () => {
    const { gateway, spies } = createGatewaySpies()
    spies.resumeTermTranslationPhase.mockRejectedValueOnce(
      new Error("Wails binding is not wired yet: ResumeTermTranslationPhase")
    )
    const store = createStore({
      jobId: 9,
      phase: "ready",
      summary: createSummary({ phaseRunId: 81 })
    })
    const useCase = new TermTranslationPhaseUseCase(gateway, store)

    await useCase.resumePhase()

    expect(store.snapshot().errorMessage).toBe(
      "Wails binding is not wired yet: ResumeTermTranslationPhase"
    )
  })

  test("retry 失敗時は fallback エラーメッセージを返す", async () => {
    const { gateway, spies } = createGatewaySpies()
    spies.retryTermTranslationPhase.mockRejectedValueOnce(
      new Error("backend failure")
    )
    const store = createStore({
      jobId: 9,
      phase: "ready",
      summary: createSummary({ phaseRunId: 81 })
    })
    const useCase = new TermTranslationPhaseUseCase(gateway, store)

    await useCase.retryPhase()

    expect(store.snapshot().errorMessage).toBe(
      "単語翻訳段階の操作に失敗しました。"
    )
    expect(store.snapshot().phase).toBe("ready")
  })
})
