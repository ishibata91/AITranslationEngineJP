import { describe, expect, test, vi } from "vitest"

import type {
  CreateTranslationJobRequest,
  CreateTranslationJobResponse,
  GetTranslationJobSetupSummaryRequest,
  TranslationJobSetupGatewayContract,
  TranslationJobSetupOptionsResponse,
  TranslationJobSetupScreenState,
  TranslationJobSetupSummaryResponse,
  TranslationJobSetupValidationResponse
} from "@application/gateway-contract/translation-job-setup"

import { TranslationJobSetupUseCase } from "./translation-job-setup.usecase"

type StoreLike = {
  snapshot(): TranslationJobSetupScreenState
  update(mutator: (draft: TranslationJobSetupScreenState) => void): void
}

function createOptions(
  overrides: Partial<TranslationJobSetupOptionsResponse> = {}
): TranslationJobSetupOptionsResponse {
  return {
    inputCandidates: [
      {
        id: 41,
        label: "/mods/input-review.json",
        sourceKind: "xEdit extract",
        recordCount: 128,
        registeredAt: "2026-04-27T10:20:00Z"
      }
    ],
    sharedDictionaries: [],
    sharedPersonas: [],
    aiRuntimeOptions: [
      {
        provider: "openai",
        model: "gpt-5.4-mini",
        mode: "batch"
      }
    ],
    credentialRefs: [
      {
        provider: "openai",
        credentialRef: "openai-primary",
        isConfigured: true,
        isMissingSecret: false
      }
    ],
    ...overrides
  }
}

function createValidationResult(): TranslationJobSetupValidationResponse {
  return {
    status: "pass",
    targetSlices: ["input", "runtime", "credentials"],
    validatedAt: "2026-04-27T10:30:00Z",
    canCreate: true,
    passSlices: ["input", "runtime", "credentials"]
  }
}

function createSummary(): TranslationJobSetupSummaryResponse {
  return {
    jobId: 91,
    jobState: "ready",
    inputSource: "/mods/input-review.json",
    canStartPhase: true,
    executionSummary: {
      provider: "openai",
      model: "gpt-5.4-mini",
      executionMode: "batch"
    },
    validationPassSlices: ["input", "runtime", "credentials"]
  }
}

function createState(
  overrides: Partial<TranslationJobSetupScreenState> = {}
): TranslationJobSetupScreenState {
  return {
    phase: "ready",
    options: createOptions(),
    selectedInputSourceId: 41,
    selectedRuntimeKey: "openai::gpt-5.4-mini::batch",
    selectedCredentialRef: "openai-primary",
    validationResult: createValidationResult(),
    validationState: "fresh",
    dirty: false,
    errorMessage: "",
    createErrorKind: null,
    summary: null,
    ...overrides
  }
}

function clone<T>(value: T): T {
  return structuredClone(value)
}

function createStore(initialState: TranslationJobSetupScreenState = createState()): StoreLike {
  let state = clone(initialState)
  return {
    snapshot() {
      return clone(state)
    },
    update(mutator) {
      const draft = clone(state)
      mutator(draft)
      state = draft
    }
  }
}

function createGateway(): TranslationJobSetupGatewayContract & {
  createTranslationJob: ReturnType<typeof vi.fn>
  getTranslationJobSetupSummary: ReturnType<typeof vi.fn>
} {
  return {
    getTranslationJobSetupOptions: vi.fn(),
    validateTranslationJobSetup: vi.fn(),
    createTranslationJob: vi.fn<
      (request: CreateTranslationJobRequest) => Promise<CreateTranslationJobResponse>
    >().mockResolvedValue({
      jobId: 91,
      jobState: "ready",
      inputSource: "/mods/input-review.json",
      executionSummary: {
        provider: "openai",
        model: "gpt-5.4-mini",
        executionMode: "batch"
      },
      validationPassSlices: ["input", "runtime", "credentials"]
    }),
    getTranslationJobSetupSummary: vi.fn<
      (request: GetTranslationJobSetupSummaryRequest) => Promise<TranslationJobSetupSummaryResponse>
    >().mockResolvedValue(createSummary())
  }
}

describe("TranslationJobSetupUseCase", () => {
  test("createJob は validation freshness を create request へ転送する", async () => {
    const gateway = createGateway()
    const store = createStore()
    const usecase = new TranslationJobSetupUseCase(gateway, store)

    await usecase.createJob()

    expect(gateway.createTranslationJob).toHaveBeenCalledTimes(1)
    expect(gateway.createTranslationJob).toHaveBeenCalledWith({
      inputSourceId: 41,
      inputSource: "/mods/input-review.json",
      validationStatus: "pass",
      validatedAt: "2026-04-27T10:30:00Z",
      validationPassSlices: ["input", "runtime", "credentials"],
      runtime: {
        provider: "openai",
        model: "gpt-5.4-mini",
        executionMode: "batch"
      },
      credentialRef: "openai-primary"
    })
  })

  test("selectedInputSourceId が異なる existingJob は inputSource 表示名が一致しても create を無効化しない", async () => {
    const gateway = createGateway()
    const store = createStore(
      createState({
        options: createOptions({
          inputCandidates: [
            {
              id: 41,
              label: "/mods/input-review.json",
              sourceKind: "xEdit extract",
              recordCount: 128,
              registeredAt: "2026-04-27T10:20:00Z"
            },
            {
              id: 42,
              label: "/mods/other-input-review.json",
              sourceKind: "xEdit extract",
              recordCount: 64,
              registeredAt: "2026-04-27T10:25:00Z"
            }
          ],
          existingJob: {
            inputSourceId: 42,
            jobId: 300,
            status: "ready",
            inputSource: "/mods/input-review.json"
          }
        })
      })
    )
    const usecase = new TranslationJobSetupUseCase(gateway, store)

    await usecase.createJob()

    expect(gateway.createTranslationJob).toHaveBeenCalledTimes(1)
    expect(gateway.createTranslationJob).toHaveBeenCalledWith(
      expect.objectContaining({
        inputSourceId: 41,
        inputSource: "/mods/input-review.json"
      })
    )
  })

  test("selectedInputSourceId と一致する existingJob は inputSource 表示名が異なっても create を無効化する", async () => {
    const gateway = createGateway()
    const store = createStore(
      createState({
        options: createOptions({
          existingJob: {
            inputSourceId: 41,
            jobId: 300,
            status: "ready",
            inputSource: "Imported from another label"
          }
        })
      })
    )
    const usecase = new TranslationJobSetupUseCase(gateway, store)

    await usecase.createJob()

    expect(gateway.createTranslationJob).not.toHaveBeenCalled()
    expect(store.snapshot().errorMessage).toBe(
      "create 条件を満たしていません。validation と既存 job 状態を確認してください。"
    )
  })

  test("createJob 後の summary fetch は canStartPhase を保持して state へ反映する", async () => {
    const gateway = createGateway()
    const store = createStore()
    const usecase = new TranslationJobSetupUseCase(gateway, store)

    await usecase.createJob()

    expect(gateway.getTranslationJobSetupSummary).toHaveBeenCalledTimes(1)
    expect(gateway.getTranslationJobSetupSummary).toHaveBeenCalledWith({ jobId: 91 })
    expect(store.snapshot().summary).toEqual(
      expect.objectContaining({
        jobId: 91,
        canStartPhase: true
      })
    )
  })
})
