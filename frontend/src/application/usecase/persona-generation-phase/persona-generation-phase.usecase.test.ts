import { describe, expect, test, vi } from "vitest"
import { PersonaGenerationPhaseUseCase } from "./persona-generation-phase.usecase"
import type {
  ProcessingTargetListPageState,
  ProcessingTargetListResponse
} from "@application/gateway-contract/processing-target"
import type {
  PersonaGenerationBodyReadinessResponse,
  PersonaGenerationPhaseCommandResponse,
  PersonaGenerationPhaseGatewayContract,
  PersonaGenerationPhaseSummaryResponse
} from "@application/gateway-contract/persona-generation-phase"

type ActionKind =
  | "start"
  | "pause"
  | "resume"
  | "retry"
  | "cancel"
  | "check-body-readiness"
  | "start-body-phase"

interface ScreenState {
  jobId: number | null
  phase: "idle" | "loading" | "ready" | "submitting"
  summary: PersonaGenerationPhaseSummaryResponse | null
  bodyReadiness: PersonaGenerationBodyReadinessResponse | null
  errorMessage: string
  pendingAction: ActionKind | null
  hasLoaded: boolean
  processingTargetPageState?: ProcessingTargetListPageState | null
}

interface StoreLike {
  snapshot(): ScreenState
  update(mutator: (draft: ScreenState) => void): void
}

function createStore(initial: ScreenState): StoreLike {
  let state = structuredClone(initial)
  return {
    snapshot: () => structuredClone(state),
    update: (mutator: (draft: ScreenState) => void) => {
      const next = structuredClone(state)
      mutator(next)
      state = next
    }
  }
}

interface GatewayWithSpies {
  gateway: PersonaGenerationPhaseGatewayContract
  getSummarySpy: ReturnType<typeof vi.fn>
  getProcessingTargetListSpy: ReturnType<typeof vi.fn>
  startSpy: ReturnType<typeof vi.fn>
  pauseSpy: ReturnType<typeof vi.fn>
  resumeSpy: ReturnType<typeof vi.fn>
  retrySpy: ReturnType<typeof vi.fn>
  cancelSpy: ReturnType<typeof vi.fn>
}

function createGateway(): GatewayWithSpies {
  const processingTargetResponse: ProcessingTargetListResponse = {
    items: [
      {
        id: "persona:1",
        name: "Aela",
        detail: "FormID: 0001A696",
        titleParts: [{ text: "Aela" }],
        metadata: [{ label: "種族", value: "Nord" }]
      }
    ],
    metadata: [],
    page: 3,
    pageSize: 25,
    totalCount: 88,
    searchQuery: "Nord"
  }
  const getSummarySpy = vi.fn(
    (request: {
      jobId: number
    }): Promise<PersonaGenerationPhaseSummaryResponse> =>
      Promise.resolve({
        jobId: request.jobId,
        currentPhase: "persona_generation",
        phaseState: "running",
        phaseRunId: 77,
        progress: {
          percent: 10,
          processedCount: 1,
          totalCount: 10,
          targetCount: 10,
          currentStep: "running"
        },
        targetSummary: {
          targetCount: 10,
          commonPersonaHitCount: 0,
          commonPersonaMissCount: 10,
          skippedCount: 0,
          skippedReasons: [],
          targetSnapshotDigest: "sha256:1"
        },
        execution: {
          credentialRef: "cred",
          provider: "fake",
          model: "m",
          executionMode: "single_request",
          promptDigest: "sha256:1",
          inputCount: 10,
          outputCount: 0,
          evidenceRefs: []
        },
        actionEnablement: {
          canStart: false,
          canPause: true,
          canResume: false,
          canRetry: false,
          canCancel: true,
          canStartBodyPhase: false
        }
      })
  )
  const getBodySpy = vi.fn((request: { jobId: number }) =>
    Promise.resolve({
      jobId: request.jobId,
      currentPhase: "persona_generation",
      phaseState: "running",
      ready: false,
      errorKind: "body_readiness_blocked" as const,
      inputSummary: {
        personaCount: 0,
        missingCount: 10,
        snapshotId: "",
        snapshotDigest: "sha256:1",
        evidenceRefs: []
      }
    })
  )
  const getProcessingTargetListSpy = vi.fn(() =>
    Promise.resolve(processingTargetResponse)
  )
  const command = (
    phaseState: string
  ): PersonaGenerationPhaseCommandResponse => ({
    jobId: 10,
    currentPhase: "persona_generation",
    phaseState,
    progress: {
      percent: 0,
      processedCount: 0,
      totalCount: 10,
      targetCount: 10,
      currentStep: "running"
    },
    targetSummary: {
      targetCount: 10,
      commonPersonaHitCount: 0,
      commonPersonaMissCount: 10,
      skippedCount: 0,
      skippedReasons: [],
      targetSnapshotDigest: "sha256:1"
    },
    execution: {
      credentialRef: "cred",
      provider: "fake",
      model: "m",
      executionMode: "single_request",
      promptDigest: "sha256:1",
      inputCount: 10,
      outputCount: 0,
      evidenceRefs: []
    },
    retryable: false,
    canStartBodyPhase: false
  })
  const startSpy = vi.fn(() => Promise.resolve(command("running")))
  const pauseSpy = vi.fn(() => Promise.resolve(command("paused")))
  const resumeSpy = vi.fn(() => Promise.resolve(command("running")))
  const retrySpy = vi.fn(() => Promise.resolve(command("running")))
  const cancelSpy = vi.fn(() => Promise.resolve(command("canceled")))

  return {
    gateway: {
      getPersonaGenerationPhaseSummary: getSummarySpy,
      getPersonaGenerationBodyReadiness: getBodySpy,
      getProcessingTargetList: getProcessingTargetListSpy,
      startPersonaGenerationPhase: startSpy,
      pausePersonaGenerationPhase: pauseSpy,
      resumePersonaGenerationPhase: resumeSpy,
      retryPersonaGenerationPhase: retrySpy,
      cancelPersonaGenerationPhase: cancelSpy
    },
    getSummarySpy,
    getProcessingTargetListSpy,
    startSpy,
    pauseSpy,
    resumeSpy,
    retrySpy,
    cancelSpy
  }
}

function baseSummary(): PersonaGenerationPhaseSummaryResponse {
  return {
    jobId: 10,
    currentPhase: "persona_generation",
    phaseState: "running",
    phaseRunId: 77,
    progress: {
      percent: 0,
      processedCount: 0,
      totalCount: 10,
      targetCount: 10,
      currentStep: "running"
    },
    targetSummary: {
      targetCount: 10,
      commonPersonaHitCount: 0,
      commonPersonaMissCount: 10,
      skippedCount: 0,
      skippedReasons: [],
      targetSnapshotDigest: "sha256:1"
    },
    execution: {
      credentialRef: "cred",
      provider: "fake",
      model: "m",
      executionMode: "single_request",
      promptDigest: "sha256:1",
      inputCount: 10,
      outputCount: 0,
      evidenceRefs: []
    },
    actionEnablement: {
      canStart: false,
      canPause: true,
      canResume: false,
      canRetry: false,
      canCancel: true,
      canStartBodyPhase: false
    }
  }
}

function createDeferredResponse() {
  let resolve!: (response: ProcessingTargetListResponse) => void
  const promise = new Promise<ProcessingTargetListResponse>((resolver) => {
    resolve = resolver
  })

  return { promise, resolve }
}

describe("PersonaGenerationPhaseUseCase", () => {
  test("gateway 未接続時はエラーを状態に出す", async () => {
    const store = createStore({
      jobId: 10,
      phase: "idle",
      summary: null,
      bodyReadiness: null,
      errorMessage: "",
      pendingAction: null,
      hasLoaded: false
    })
    const usecase = new PersonaGenerationPhaseUseCase(null, store)

    await usecase.load()

    expect(store.snapshot().errorMessage).toContain("gateway が未接続")
  })

  test("summary load と start/pause/resume/retry/cancel を実行できる", async () => {
    const gatewayBundle = createGateway()
    const store = createStore({
      jobId: 10,
      phase: "idle",
      summary: baseSummary(),
      bodyReadiness: null,
      errorMessage: "",
      pendingAction: null,
      hasLoaded: false
    })
    const usecase = new PersonaGenerationPhaseUseCase(
      gatewayBundle.gateway,
      store
    )

    await usecase.load()
    await usecase.startPhase()
    await usecase.pausePhase()
    await usecase.resumePhase()
    await usecase.retryPhase()
    await usecase.cancelPhase()

    expect(gatewayBundle.getSummarySpy).toHaveBeenCalled()
    expect(gatewayBundle.startSpy).toHaveBeenCalledWith({ jobId: 10 })
    expect(gatewayBundle.pauseSpy).toHaveBeenCalledWith({
      jobId: 10,
      phaseRunId: 77
    })
    expect(gatewayBundle.resumeSpy).toHaveBeenCalledWith({
      jobId: 10,
      phaseRunId: 77
    })
    expect(gatewayBundle.retrySpy).toHaveBeenCalledWith({
      jobId: 10,
      phaseRunId: 77
    })
    expect(gatewayBundle.cancelSpy).toHaveBeenCalledWith({
      jobId: 10,
      phaseRunId: 77
    })
  })

  test("load は processing target request を送り totalCount と searchQuery を page state に保持する", async () => {
    const gatewayBundle = createGateway()
    const store = createStore({
      jobId: 10,
      phase: "idle",
      summary: baseSummary(),
      bodyReadiness: null,
      errorMessage: "",
      pendingAction: null,
      hasLoaded: false,
      processingTargetPageState: {
        items: [],
        metadata: [],
        page: 3,
        pageSize: 25,
        totalCount: 0,
        searchQuery: "Nord",
        busy: false
      }
    })
    const usecase = new PersonaGenerationPhaseUseCase(
      gatewayBundle.gateway,
      store
    )

    await usecase.load()

    expect(gatewayBundle.getProcessingTargetListSpy).toHaveBeenCalledWith({
      jobId: 10,
      phase: "persona_generation",
      page: 3,
      pageSize: 25,
      searchQuery: "Nord"
    })
    expect(store.snapshot().processingTargetPageState).toMatchObject({
      items: [{ id: "persona:1" }],
      page: 3,
      pageSize: 25,
      totalCount: 88,
      searchQuery: "Nord",
      busy: false
    })
  })

  test("検索語変更は persona_generation request を page 1 へ戻して送る", async () => {
    const gatewayBundle = createGateway()
    const store = createStore({
      jobId: 10,
      phase: "ready",
      summary: baseSummary(),
      bodyReadiness: null,
      errorMessage: "",
      pendingAction: null,
      hasLoaded: true,
      processingTargetPageState: {
        items: [],
        metadata: [],
        page: 5,
        pageSize: 25,
        totalCount: 0,
        searchQuery: "",
        busy: false
      }
    })
    const usecase = new PersonaGenerationPhaseUseCase(
      gatewayBundle.gateway,
      store
    )

    await usecase.setProcessingTargetSearchQuery("Mage")

    expect(gatewayBundle.getProcessingTargetListSpy).toHaveBeenCalledWith({
      jobId: 10,
      phase: "persona_generation",
      page: 1,
      pageSize: 25,
      searchQuery: "Mage"
    })
    expect(store.snapshot().processingTargetPageState).toMatchObject({
      page: 3,
      pageSize: 25,
      totalCount: 88,
      searchQuery: "Nord"
    })
  })

  test("検索応答の到着順が逆転しても最新検索結果だけを page state に反映する", async () => {
    const gatewayBundle = createGateway()
    const first = createDeferredResponse()
    const second = createDeferredResponse()
    gatewayBundle.getProcessingTargetListSpy
      .mockReturnValueOnce(first.promise)
      .mockReturnValueOnce(second.promise)
    const store = createStore({
      jobId: 10,
      phase: "ready",
      summary: baseSummary(),
      bodyReadiness: null,
      errorMessage: "",
      pendingAction: null,
      hasLoaded: true,
      processingTargetPageState: {
        items: [],
        metadata: [],
        page: 5,
        pageSize: 25,
        totalCount: 0,
        searchQuery: "",
        busy: false
      }
    })
    const usecase = new PersonaGenerationPhaseUseCase(
      gatewayBundle.gateway,
      store
    )

    const firstSearch = usecase.setProcessingTargetSearchQuery("Nord")
    const secondSearch = usecase.setProcessingTargetSearchQuery("Mage")

    second.resolve({
      items: [
        {
          id: "persona:mage",
          name: "Mage",
          detail: "EditorID: Mage",
          titleParts: [{ text: "Mage" }],
          metadata: []
        }
      ],
      metadata: [],
      page: 1,
      pageSize: 25,
      totalCount: 2,
      searchQuery: "Mage"
    })
    await secondSearch

    first.resolve({
      items: [
        {
          id: "persona:nord",
          name: "Nord",
          detail: "Race: Nord",
          titleParts: [{ text: "Nord" }],
          metadata: []
        }
      ],
      metadata: [],
      page: 1,
      pageSize: 25,
      totalCount: 7,
      searchQuery: "Nord"
    })
    await firstSearch

    expect(store.snapshot().processingTargetPageState).toMatchObject({
      items: [{ id: "persona:mage" }],
      page: 1,
      pageSize: 25,
      totalCount: 2,
      searchQuery: "Mage",
      busy: false
    })
  })

  test("load の一覧応答が検索応答より遅れても最新検索結果だけを page state に反映する", async () => {
    const gatewayBundle = createGateway()
    const loadResponse = createDeferredResponse()
    const searchResponse = createDeferredResponse()
    gatewayBundle.getProcessingTargetListSpy
      .mockReturnValueOnce(loadResponse.promise)
      .mockReturnValueOnce(searchResponse.promise)
    const store = createStore({
      jobId: 10,
      phase: "ready",
      summary: baseSummary(),
      bodyReadiness: null,
      errorMessage: "",
      pendingAction: null,
      hasLoaded: true,
      processingTargetPageState: {
        items: [],
        metadata: [],
        page: 5,
        pageSize: 25,
        totalCount: 0,
        searchQuery: "",
        busy: false
      }
    })
    const usecase = new PersonaGenerationPhaseUseCase(
      gatewayBundle.gateway,
      store
    )

    const load = usecase.load()
    const search = usecase.setProcessingTargetSearchQuery("Mage")

    searchResponse.resolve({
      items: [
        {
          id: "persona:mage",
          name: "Mage",
          detail: "EditorID: Mage",
          titleParts: [{ text: "Mage" }],
          metadata: []
        }
      ],
      metadata: [],
      page: 1,
      pageSize: 25,
      totalCount: 2,
      searchQuery: "Mage"
    })
    await search

    loadResponse.resolve({
      items: [
        {
          id: "persona:old",
          name: "Old",
          detail: "EditorID: Old",
          titleParts: [{ text: "Old" }],
          metadata: []
        }
      ],
      metadata: [],
      page: 5,
      pageSize: 25,
      totalCount: 99,
      searchQuery: ""
    })
    await load

    expect(store.snapshot().processingTargetPageState).toMatchObject({
      items: [{ id: "persona:mage" }],
      page: 1,
      pageSize: 25,
      totalCount: 2,
      searchQuery: "Mage",
      busy: false
    })
  })

  test("Wails 未接続エラーは文言を保持して redacted エラー表示状態に載せる", async () => {
    const gatewayBundle = createGateway()
    gatewayBundle.getSummarySpy.mockRejectedValueOnce(
      new Error(
        "Wails binding is not wired yet: GetPersonaGenerationPhaseSummary"
      )
    )
    const store = createStore({
      jobId: 10,
      phase: "idle",
      summary: null,
      bodyReadiness: null,
      errorMessage: "",
      pendingAction: null,
      hasLoaded: false
    })
    const usecase = new PersonaGenerationPhaseUseCase(
      gatewayBundle.gateway,
      store
    )

    await usecase.load()

    expect(store.snapshot().errorMessage).toContain(
      "Wails binding is not wired yet"
    )
  })
})
