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
    projection: {
      phaseLifecycle: "ready",
      jobLifecycle: "running",
      errorKind: "none",
      aiSettingsConfigured: true,
      aiTargetCount: 7,
      confirmedCount: 0
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
    jobIsTerminal: false,
    totalCount: 10,
    confirmedCount: 0,
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
  initialFetchDone: boolean
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
          execution: state.summary.execution ? { ...state.summary.execution } : undefined,
          resultSummary: state.summary.resultSummary
            ? { ...state.summary.resultSummary }
            : undefined,
          errorSummary: state.summary.errorSummary
            ? { ...state.summary.errorSummary }
            : undefined,
          projection: { ...state.summary.projection }
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
    initialFetchDone: false,
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
        retryable: false
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
        retryable: true
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
        retryable: false
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
        retryable: false
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

  test("setJobId は summary を再取得し hasLoaded と initialFetchDone を true に更新する", async () => {
    const { gateway, spies } = createGatewaySpies()
    const store = createStore()
    const useCase = new TermTranslationPhaseUseCase(gateway, store)

    await useCase.setJobId(9)

    expect(spies.getTermTranslationPhaseSummary).toHaveBeenCalledWith({
      jobId: 9
    })
    expect(spies.getTermTranslationNextPhaseReadiness).not.toHaveBeenCalled()
    expect(store.snapshot().hasLoaded).toBe(true)
    expect(store.snapshot().initialFetchDone).toBe(true)
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
    const oldReadiness = createReadiness({ phaseState: "paused" })
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

  // UT-FETCH-001: setJobId で summary と processingTarget の両方が起動し最大2本に収まる
  test("setJobId は summary と processingTarget の両方を取得し可否専用取得を呼ばない", async () => {
    const { gateway, spies } = createGatewaySpies()
    const store = createStore()
    const useCase = new TermTranslationPhaseUseCase(gateway, store)

    await useCase.setJobId(9)

    // summary 取得が 1 回起動する
    expect(spies.getTermTranslationPhaseSummary).toHaveBeenCalledTimes(1)
    // processingTarget 取得が 1 回起動する（最大2本に収まる）
    expect(spies.getProcessingTargetList).toHaveBeenCalledTimes(1)
    // 可否専用取得は起動しない（フロント導出のため取得物でなくなった）
    expect(spies.getTermTranslationNextPhaseReadiness).not.toHaveBeenCalled()
    // 両取得完了後に initialFetchDone が true になる
    expect(store.snapshot().initialFetchDone).toBe(true)
  })

  // UT-ASYM-001: 後発取得（新しい sequence）が processingTarget を反映し空状態のまま残さない
  test("後発の processingTarget 取得が先行取得より遅れて解決しても反映され空状態にならない", async () => {
    const { gateway, spies } = createGatewaySpies()
    const summaryDeferred = { resolve: null as ((v: TermTranslationPhaseSummaryResponse) => void) | null }
    const processingTargetDeferred = createDeferredResponse()

    // summary は即時解決、processingTarget は遅延
    spies.getTermTranslationPhaseSummary.mockImplementationOnce(
      () => Promise.resolve(createSummary())
    )
    spies.getProcessingTargetList.mockImplementationOnce(
      () => processingTargetDeferred.promise
    )
    const store = createStore()
    const useCase = new TermTranslationPhaseUseCase(gateway, store)

    // setJobId 起動（processingTarget は保留中）
    const setJobIdPromise = useCase.setJobId(9)

    // processingTarget の遅延応答を解決する（後発取得）
    processingTargetDeferred.resolve({
      items: [
        {
          id: "term:1",
          name: "Dragon",
          detail: "原文: Dragon",
          titleParts: [{ text: "Dragon" }],
          metadata: []
        }
      ],
      metadata: [],
      page: 1,
      pageSize: 50,
      totalCount: 1,
      searchQuery: ""
    })

    await setJobIdPromise

    // processingTarget が反映され空状態のまま残らない
    expect(store.snapshot().processingTargetPageState?.items.length).toBeGreaterThan(0)
    expect(store.snapshot().initialFetchDone).toBe(true)
    void summaryDeferred
  })

  // UT-ASYM-003: summary 反映は processingTarget と独立（summary 失敗でも processingTarget は反映される）
  test("summary 取得失敗時も processingTarget の反映は独立して行われる", async () => {
    const { gateway, spies } = createGatewaySpies()
    spies.getTermTranslationPhaseSummary.mockRejectedValueOnce(
      new Error("summary timeout")
    )
    const store = createStore()
    const useCase = new TermTranslationPhaseUseCase(gateway, store)

    await useCase.setJobId(9)

    // summary は失敗するがエラーメッセージが設定される
    expect(store.snapshot().errorMessage).toContain("summary")
    // processingTarget は独立して反映される（summary 失敗を理由に取りこぼさない）
    expect(store.snapshot().processingTargetPageState?.items.length).toBeGreaterThan(0)
    expect(store.snapshot().initialFetchDone).toBe(true)
  })

  // UT-REOPEN-001: setJobId で旧 sequence を無効化して再取得し旧遅延応答を破棄する
  test("setJobId 呼び出し中の旧取得が遅れて解決しても setJobId 後の新 sequence に上書きされない", async () => {
    const { gateway, spies } = createGatewaySpies()
    const oldProcessingTargetDeferred = createDeferredResponse()

    // 1 回目の setJobId の processingTarget 取得は遅延させる
    spies.getProcessingTargetList.mockImplementationOnce(
      () => oldProcessingTargetDeferred.promise
    )

    const store = createStore()
    const useCase = new TermTranslationPhaseUseCase(gateway, store)

    // 1 回目の setJobId を開始（processingTarget は保留中のまま）
    const firstSetJobIdPromise = useCase.setJobId(9)

    // 2 回目の setJobId を呼ぶ（新しい jobId、新しい取得が開始される）
    // この時点で 1 回目の旧取得はまだ保留中
    await useCase.setJobId(10)

    // 2 回目の setJobId が完了し、新しい processingTarget が反映されている
    const stateAfterSecondSetJobId = store.snapshot()
    expect(stateAfterSecondSetJobId.jobId).toBe(10)
    expect(stateAfterSecondSetJobId.processingTargetPageState?.items.length).toBeGreaterThan(0)

    // 旧取得（1 回目の保留中応答）を遅れて解決する
    oldProcessingTargetDeferred.resolve({
      items: [
        {
          id: "old:1",
          name: "OldItem",
          detail: "旧取得の応答",
          titleParts: [{ text: "OldItem" }],
          metadata: []
        }
      ],
      metadata: [],
      page: 1,
      pageSize: 50,
      totalCount: 1,
      searchQuery: ""
    })

    await firstSetJobIdPromise

    // 旧取得の応答が新しい sequence の結果を上書きしていないことを確認
    const finalState = store.snapshot()
    const items = finalState.processingTargetPageState?.items ?? []
    const hasOldItem = items.some((item) => item.id === "old:1")
    expect(hasOldItem).toBe(false)
  })

  // UT-REOPEN-002: 再取得後に items.length>0（総件数 1 以上の場合）
  test("setJobId による再取得で総件数1以上の段階では items が反映され空状態にならない", async () => {
    const { gateway } = createGatewaySpies()
    const store = createStore()
    const useCase = new TermTranslationPhaseUseCase(gateway, store)

    // 既存スパイが items を含む応答を返す設定になっている
    await useCase.setJobId(9)

    // 再取得結果が反映され items が空でない
    const pageState = store.snapshot().processingTargetPageState
    expect(pageState).not.toBeNull()
    expect(pageState!.totalCount).toBeGreaterThan(0)
    expect(pageState!.items.length).toBeGreaterThan(0)
  })

  // UT-LOAD-001: summary と processingTarget の両取得完了後のみ initialFetchDone=true
  test("summary 取得のみ完了で processingTarget 未完了の間は initialFetchDone が false のまま", async () => {
    const { gateway, spies } = createGatewaySpies()
    const processingTargetDeferred = createDeferredResponse()
    spies.getProcessingTargetList.mockImplementationOnce(
      () => processingTargetDeferred.promise
    )
    const store = createStore()
    const useCase = new TermTranslationPhaseUseCase(gateway, store)

    // setJobId 起動（processingTarget は保留中）
    const setJobIdPromise = useCase.setJobId(9)

    // summary は即時解決するが processingTarget はまだ保留
    // この時点では initialFetchDone は false のまま
    await Promise.resolve() // microtask をフラッシュ
    expect(store.snapshot().initialFetchDone).toBe(false)

    // processingTarget の取得を完了させる
    processingTargetDeferred.resolve({
      items: [],
      metadata: [],
      page: 1,
      pageSize: 50,
      totalCount: 0,
      searchQuery: ""
    })
    await setJobIdPromise

    // 両取得完了後に initialFetchDone が true になる
    expect(store.snapshot().initialFetchDone).toBe(true)
  })

  // UT-DISP-001: initialFetchDone=true かつ総件数1以上で空状態判定にならない（items が反映される）
  test("setJobId 後に initialFetchDone=true かつ processingTarget の totalCount が 1 以上になる", async () => {
    const { gateway } = createGatewaySpies()
    const store = createStore()
    const useCase = new TermTranslationPhaseUseCase(gateway, store)

    // createGatewaySpies の processingTarget は totalCount=137 で items 1 件
    await useCase.setJobId(9)

    const state = store.snapshot()
    // initialFetchDone=true を評価条件とした derived 判定の前提が成立している
    expect(state.initialFetchDone).toBe(true)
    // 総件数が 1 以上で空状態判定にならない（件数分の items が保持される）
    expect(state.processingTargetPageState?.totalCount).toBeGreaterThan(0)
    expect(state.processingTargetPageState?.items.length).toBeGreaterThan(0)
  })

  // UT-DISP-002: 進捗母数と処理対象一覧の総件数は独立値として保持される
  test("summary の aiTargetCount と processingTarget の totalCount は独立した値として store に保持される", async () => {
    const { gateway, spies } = createGatewaySpies()
    // summary の aiTargetCount = 7（進捗母数）
    spies.getTermTranslationPhaseSummary.mockResolvedValueOnce(
      createSummary({ aiTargetCount: 7, progress: { percent: 0, processedCount: 0, totalCount: 10, aiTargetCount: 7, currentStep: "ready" } })
    )
    // processingTarget の totalCount = 137（処理対象一覧の総件数）
    // createGatewaySpies のデフォルトが totalCount=137 を返す
    const store = createStore()
    const useCase = new TermTranslationPhaseUseCase(gateway, store)

    await useCase.setJobId(9)

    const state = store.snapshot()
    // 進捗母数（summary.aiTargetCount = 7）と総件数（processingTarget.totalCount = 137）は別値
    expect(state.summary?.aiTargetCount).toBe(7)
    expect(state.processingTargetPageState?.totalCount).toBe(137)
    // 値が一致しなくても不具合ではない（独立した値として扱う）
    expect(state.summary?.aiTargetCount).not.toBe(state.processingTargetPageState?.totalCount)
  })
})
