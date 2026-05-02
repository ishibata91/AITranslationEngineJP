import { describe, expect, test, vi } from "vitest"

import type {
  TermTranslationNextPhaseReadinessResponse,
  TermTranslationPhaseGatewayContract,
  TermTranslationPhaseSummaryResponse
} from "@application/gateway-contract/term-translation-phase"

import { TermTranslationPhaseUseCase } from "./term-translation-phase.usecase"

function createSummary(overrides: Partial<TermTranslationPhaseSummaryResponse> = {}): TermTranslationPhaseSummaryResponse {
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
      canRefresh: true,
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
  pendingAction: "refresh" | "start" | "pause" | "resume" | "retry" | null
  hasLoaded: boolean
}

function cloneState(state: TermTranslationPhaseScreenStateLike): TermTranslationPhaseScreenStateLike {
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
      : null
  }
}

function createStore(initialState: Partial<TermTranslationPhaseScreenStateLike> = {}) {
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
  const spies = {
    getTermTranslationPhaseSummary: vi.fn(() => Promise.resolve(createSummary())),
    getTermTranslationNextPhaseReadiness: vi.fn(() =>
      Promise.resolve(createReadiness())
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

describe("TermTranslationPhaseUseCase", () => {
  test("gateway 未接続で refresh は接続エラーメッセージを設定する", async () => {
    const store = createStore({ jobId: 9, phase: "ready" })
    const useCase = new TermTranslationPhaseUseCase(null, store)

    await useCase.refresh()

    expect(store.snapshot().errorMessage).toBe("term-translation-phase gateway が未接続です。")
  })

  test("setJobId は summary を再取得し hasLoaded を true に更新する", async () => {
    const { gateway, spies } = createGatewaySpies()
    const store = createStore()
    const useCase = new TermTranslationPhaseUseCase(gateway, store)

    await useCase.setJobId(9)

    expect(spies.getTermTranslationPhaseSummary).toHaveBeenCalledWith({ jobId: 9 })
    expect(spies.getTermTranslationNextPhaseReadiness).toHaveBeenCalledWith({ jobId: 9 })
    expect(store.snapshot().hasLoaded).toBe(true)
  })

  test("summary refresh 失敗時は前回 snapshot を保持する", async () => {
    const { gateway, spies } = createGatewaySpies()
    spies.getTermTranslationPhaseSummary.mockRejectedValueOnce(new Error("timeout"))
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

    await useCase.refresh()

    const snapshot = store.snapshot()
    expect(snapshot.summary).toMatchObject(oldSummary)
    expect(snapshot.nextPhaseReadiness).toMatchObject(oldReadiness)
    expect(snapshot.errorMessage).toBe("単語翻訳フェーズ summary の取得に失敗しました。")
  })

  test("start 成功時は start command 実行後に summary refresh する", async () => {
    const { gateway, spies } = createGatewaySpies()
    const store = createStore({ jobId: 9, phase: "ready", summary: createSummary() })
    const useCase = new TermTranslationPhaseUseCase(gateway, store)

    await useCase.startPhase()

    expect(spies.startTermTranslationPhase).toHaveBeenCalledWith({ jobId: 9 })
    expect(spies.getTermTranslationPhaseSummary).toHaveBeenCalledTimes(1)
    expect(store.snapshot().pendingAction).toBeNull()
  })

  test("pause は phaseRunId 未確定の時に command を呼ばずエラーを返す", async () => {
    const { gateway, spies } = createGatewaySpies()
    const store = createStore({ jobId: 9, phase: "ready", summary: createSummary() })
    const useCase = new TermTranslationPhaseUseCase(gateway, store)

    await useCase.pausePhase()

    expect(spies.pauseTermTranslationPhase).not.toHaveBeenCalled()
    expect(store.snapshot().errorMessage).toBe("phase run が未確定のため操作できません。")
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
    spies.retryTermTranslationPhase.mockRejectedValueOnce(new Error("backend failure"))
    const store = createStore({
      jobId: 9,
      phase: "ready",
      summary: createSummary({ phaseRunId: 81 })
    })
    const useCase = new TermTranslationPhaseUseCase(gateway, store)

    await useCase.retryPhase()

    expect(store.snapshot().errorMessage).toBe("単語翻訳フェーズ操作に失敗しました。")
    expect(store.snapshot().phase).toBe("ready")
  })
})
