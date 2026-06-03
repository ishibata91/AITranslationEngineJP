import { describe, expect, test, vi } from "vitest"

import type {
  TermTranslationPhaseScreenState,
  TermTranslationPhaseScreenViewModel
} from "@application/contract/term-translation-phase"

import { TermTranslationPhaseScreenController } from "./term-translation-phase-screen-controller"

function createState(
  overrides: Partial<TermTranslationPhaseScreenState> = {}
): TermTranslationPhaseScreenState {
  return {
    jobId: 7,
    phase: "ready",
    summary: null,
    nextPhaseReadiness: null,
    errorMessage: "",
    pendingAction: null,
    hasLoaded: true,
    initialFetchDone: false,
    ...overrides
  }
}

function createViewModel(
  state: TermTranslationPhaseScreenState
): TermTranslationPhaseScreenViewModel {
  return {
    ...state,
    gatewayStatus: "接続準備済み",
    viewState: "blocked",
    isLoading: false,
    isRefreshing: false,
    isSubmitting: false,
    hasJobSelection: true,
    currentPhaseLabel: "term_translation",
    phaseStateLabel: "開始待ち",
    statusTitle: "title",
    statusText: "text",
    progressPercent: 0,
    progressLabel: "0%",
    progressDetail: "detail",
    startedAtLabel: "-",
    finishedAtLabel: "-",
    totalTermCountLabel: "-",
    dictionaryHitCountLabel: "-",
    aiTargetCountLabel: "-",
    confirmedCountLabel: "-",
    jobDictionaryAppliedCountLabel: "-",
    replacementTargetCountLabel: "-",
    unmatchedCountLabel: "-",
    isExecutionConfigured: false,
    providerLabel: "設定未完了",
    modelLabel: "設定未完了",
    executionModeLabel: "設定未完了",
    credentialRefLabel: "設定未完了",
    snapshotLabel: "-",
    errorKindLabel: "-",
    errorReasonLabel: "-",
    retryableLabel: "-",
    nextPhaseStatusLabel: "開始不可",
    nextPhaseBlockedReason: "",
    providerSkippedLabel: "-",
    providerOptions: [],
    modelOptions: [],
    executionOptions: [],
    actionCards: [],
    lastErrorSummary: null,
    projection: null,
    latestProgressSummary: null,
    latestResultSummary: null,
    latestExecutionSummary: null,
    latestErrorKind: null
  }
}

interface FakeGateway {
  listProviderModels?: (request: {
    provider: string
    credentialStatus: string
    requestToken: string
  }) => Promise<{
    provider: string
    status: string
    models: { modelId: string; label: string }[]
  }>
}

function createHarness(gatewayOverride?: FakeGateway | null) {
  const state = createState()

  const store = {
    subscribe: vi.fn(
      (listener: (value: TermTranslationPhaseScreenState) => void) => {
        listener(state)
        return () => {}
      }
    ),
    snapshot: vi.fn(() => state)
  }

  const presenter = {
    toViewModel: vi.fn(
      (
        value: TermTranslationPhaseScreenState,
        isGatewayConnected: boolean
      ) => ({
        ...createViewModel(value),
        gatewayStatus: isGatewayConnected ? "接続準備済み" : "未接続"
      })
    )
  }

  const useCase = {
    load: vi.fn(async () => {}),
    setJobId: vi.fn((jobId: number | null) =>
      Promise.resolve(jobId).then(() => undefined)
    ),
    startPhase: vi.fn(async () => {}),
    pausePhase: vi.fn(async () => {}),
    resumePhase: vi.fn(async () => {}),
    retryPhase: vi.fn(async () => {})
  }

  const controller = new TermTranslationPhaseScreenController({
    isGatewayConnected: false,
    store,
    presenter,
    useCase,
    gateway: gatewayOverride ?? undefined
  })

  return { controller, presenter, useCase, state }
}

describe("TermTranslationPhaseScreenController", () => {
  test("mount は useCase.load を呼ぶ", async () => {
    const harness = createHarness()

    await harness.controller.mount()

    expect(harness.useCase.load).toHaveBeenCalledTimes(1)
  })

  test("subscribe は store state を presenter 経由で view model に変換する", () => {
    const harness = createHarness()
    const listener =
      vi.fn<(value: TermTranslationPhaseScreenViewModel) => void>()

    harness.controller.subscribe(listener)

    expect(harness.presenter.toViewModel).toHaveBeenCalledWith(
      harness.state,
      false,
      [],
      []
    )
    const firstViewModel = listener.mock.calls[0]?.[0]
    expect(firstViewModel?.gatewayStatus).toBe("未接続")
  })

  test("action dispatch は対応 useCase メソッドへ委譲する", async () => {
    const harness = createHarness()

    await harness.controller.setJobId(11)
    await harness.controller.startPhase()
    await harness.controller.pausePhase()
    await harness.controller.resumePhase()
    await harness.controller.retryPhase()

    expect(harness.useCase.setJobId).toHaveBeenCalledWith(11)
    expect(harness.useCase.startPhase).toHaveBeenCalledTimes(1)
    expect(harness.useCase.pausePhase).toHaveBeenCalledTimes(1)
    expect(harness.useCase.resumePhase).toHaveBeenCalledTimes(1)
    expect(harness.useCase.retryPhase).toHaveBeenCalledTimes(1)
  })

  test("refreshModelList は provider 値付きで gateway.listProviderModels を呼ぶ", async () => {
    const listProviderModels = vi.fn(() =>
      Promise.resolve({
        provider: "gemini",
        status: "ok",
        models: [
          { modelId: "gemini-1.5-flash", label: "Gemini 1.5 Flash" },
          { modelId: "gemini-1.5-pro", label: "Gemini 1.5 Pro" }
        ]
      })
    )
    const gateway = { listProviderModels }
    const harness = createHarness(gateway)

    await harness.controller.refreshModelList("gemini")

    expect(listProviderModels).toHaveBeenCalledTimes(1)
    expect(listProviderModels).toHaveBeenCalledWith(
      expect.objectContaining({ provider: "gemini" })
    )
  })

  test("refreshModelList は provider が空文字の場合 gateway を呼ばない", async () => {
    const listProviderModels = vi.fn(() =>
      Promise.resolve({
        provider: "",
        status: "ok",
        models: [] as { modelId: string; label: string }[]
      })
    )
    const gateway = { listProviderModels }
    const harness = createHarness(gateway)

    await harness.controller.refreshModelList("")

    expect(listProviderModels).not.toHaveBeenCalled()
  })

  test("refreshModelList 成功後に availableModels が更新されリスナーへ通知される", async () => {
    const listProviderModels = vi.fn(() =>
      Promise.resolve({
        provider: "gemini",
        status: "ok",
        models: [{ modelId: "gemini-1.5-flash", label: "Gemini 1.5 Flash" }]
      })
    )
    const gateway = { listProviderModels }
    const harness = createHarness(gateway)
    const listener = vi.fn<(value: TermTranslationPhaseScreenViewModel) => void>()
    harness.controller.subscribe(listener)
    listener.mockClear()

    await harness.controller.refreshModelList("gemini")

    expect(listener).toHaveBeenCalledTimes(1)
    const vmAfterRefresh = listener.mock.calls[0]?.[0]
    expect(vmAfterRefresh).toBeDefined()
  })
})
