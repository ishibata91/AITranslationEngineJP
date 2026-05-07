import { describe, expect, test, vi } from "vitest"

import type {
  CreateTranslationJobRequest,
  CreateTranslationJobResponse,
  GetTranslationJobSetupSummaryRequest,
  ListTranslationJobSetupProviderModelsRequest,
  ListTranslationJobSetupProviderModelsResponse,
  TranslationJobSetupGatewayContract,
  TranslationJobSetupOptionsResponse,
  TranslationJobSetupPhaseRuntimeSelection,
  TranslationJobSetupScreenState,
  TranslationJobSetupSummaryResponse,
  ValidateTranslationJobSetupRequest,
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

function createPhaseOptions(
  overrides: Partial<TranslationJobSetupOptionsResponse> = {}
): TranslationJobSetupOptionsResponse {
  return createOptions({
    aiRuntimeOptions: [],
    credentialRefs: [
      {
        provider: "gemini",
        credentialRef: "gemini-primary",
        isConfigured: true,
        isMissingSecret: false
      },
      {
        provider: "xai",
        credentialRef: "xai-primary",
        isConfigured: true,
        isMissingSecret: false
      },
      {
        provider: "lm_studio",
        credentialRef: "lmstudio-local",
        isConfigured: true,
        isMissingSecret: false
      }
    ],
    providerCapabilities: [
      {
        provider: "gemini",
        credentialRequirement: "required",
        supportedExecutionModes: ["sync"],
        supportsBatchMode: true
      },
      {
        provider: "xai",
        credentialRequirement: "required",
        supportedExecutionModes: ["sync"],
        supportsBatchMode: true
      },
      {
        provider: "lm_studio",
        credentialRequirement: "not_required",
        supportedExecutionModes: ["sync"],
        supportsBatchMode: false
      }
    ],
    phaseRuntimeDrafts: [
      {
        phaseId: "word_translation",
        provider: "gemini",
        model: "gemini-word-model",
        credentialRef: "gemini-primary",
        credentialStatus: "configured",
        executionMode: "sync",
        batchMode: "enabled",
        modelListSourceToken: "gemini-source-current"
      },
      {
        phaseId: "npc_persona_generation",
        provider: "xai",
        model: "xai-persona-model",
        credentialRef: "xai-primary",
        credentialStatus: "configured",
        executionMode: "sync",
        batchMode: "disabled",
        modelListSourceToken: "xai-source-current"
      },
      {
        phaseId: "text_translation",
        provider: "xai",
        model: "xai-text-model",
        credentialRef: "xai-primary",
        credentialStatus: "configured",
        executionMode: "sync",
        batchMode: "disabled",
        modelListSourceToken: "xai-source-text"
      }
    ],
    ...overrides
  })
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
    deletingInputSourceId: null,
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

function createPhaseDrivenState(
  overrides: Partial<TranslationJobSetupScreenState> = {}
): TranslationJobSetupScreenState {
  const options = createPhaseOptions()
  const phaseRuntimeSelections =
    options.phaseRuntimeDrafts?.map(
      (draft): TranslationJobSetupPhaseRuntimeSelection => ({ ...draft })
    ) ?? []
  const providerModelLists = phaseRuntimeSelections.map(
    (selection): ListTranslationJobSetupProviderModelsResponse => ({
      phaseId: selection.phaseId,
      provider: selection.provider,
      credentialStatus: selection.credentialStatus,
      requestToken: "",
      sourceToken: selection.modelListSourceToken,
      status: "success",
      models: [{ modelId: selection.model, label: selection.model }]
    })
  )

  return createState({
    options,
    selectedRuntimeKey: null,
    selectedCredentialRef: "",
    phaseRuntimeSelections,
    providerModelLists,
    validationResult: createValidationResult(),
    validationState: "fresh",
    dirty: false,
    ...overrides
  })
}

function createDeferred<T>(): {
  promise: Promise<T>
  resolve(value: T): void
  reject(reason?: unknown): void
} {
  let resolvePromise: (value: T) => void = () => {}
  let rejectPromise: (reason?: unknown) => void = () => {}
  const promise = new Promise<T>((resolve, reject) => {
    resolvePromise = resolve
    rejectPromise = reject
  })

  return {
    promise,
    resolve: resolvePromise,
    reject: rejectPromise
  }
}

function clone<T>(value: T): T {
  return structuredClone(value)
}

function createStore(
  initialState: TranslationJobSetupScreenState = createState()
): StoreLike {
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
  getTranslationJobSetupOptions: ReturnType<typeof vi.fn>
  listTranslationJobSetupProviderModels: ReturnType<typeof vi.fn>
  validateTranslationJobSetup: ReturnType<typeof vi.fn>
  createTranslationJob: ReturnType<typeof vi.fn>
  deleteTranslationJobSetupInput: ReturnType<typeof vi.fn>
  getTranslationJobSetupSummary: ReturnType<typeof vi.fn>
} {
  return {
    getTranslationJobSetupOptions: vi
      .fn<() => Promise<TranslationJobSetupOptionsResponse>>()
      .mockResolvedValue(createOptions()),
    listTranslationJobSetupProviderModels:
      vi.fn<
        (
          request: ListTranslationJobSetupProviderModelsRequest
        ) => Promise<ListTranslationJobSetupProviderModelsResponse>
      >(),
    validateTranslationJobSetup:
      vi.fn<
        (
          request: ValidateTranslationJobSetupRequest
        ) => Promise<TranslationJobSetupValidationResponse>
      >(),
    createTranslationJob: vi
      .fn<
        (
          request: CreateTranslationJobRequest
        ) => Promise<CreateTranslationJobResponse>
      >()
      .mockResolvedValue({
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
    deleteTranslationJobSetupInput: vi.fn().mockResolvedValue({
      deletedInputSourceId: 41
    }),
    getTranslationJobSetupSummary: vi
      .fn<
        (
          request: GetTranslationJobSetupSummaryRequest
        ) => Promise<TranslationJobSetupSummaryResponse>
      >()
      .mockResolvedValue(createSummary())
  }
}

describe("TranslationJobSetupUseCase", () => {
  test("runValidation は null の targetSlices と passSlices を受けても fresh 状態へ進める", async () => {
    const gateway = createGateway()
    gateway.validateTranslationJobSetup = vi.fn().mockResolvedValue({
      status: "pass",
      targetSlices: null,
      validatedAt: "2026-05-03T06:58:30Z",
      canCreate: true,
      passSlices: null
    } as unknown as TranslationJobSetupValidationResponse)
    const store = createStore(
      createState({
        validationResult: null,
        validationState: "not-run",
        dirty: false
      })
    )
    const usecase = new TranslationJobSetupUseCase(gateway, store)

    await usecase.runValidation()

    const state = store.snapshot()
    expect(state.phase).toBe("ready")
    expect(state.validationState).toBe("fresh")
    expect(state.validationResult).toEqual({
      status: "pass",
      targetSlices: null,
      validatedAt: "2026-05-03T06:58:30Z",
      canCreate: true,
      passSlices: null
    })
  })

  test("deleteInputSource は削除成功後に options 再読込なしで対象候補だけ除去し次候補を選ぶ", async () => {
    const gateway = createGateway()
    const store = createStore(
      createPhaseDrivenState({
        options: createPhaseOptions({
          inputCandidates: [
            {
              id: 41,
              label: "/mods/input-a.json",
              sourceKind: "translation_input",
              recordCount: 128,
              registeredAt: "2026-04-27T10:20:00Z"
            },
            {
              id: 52,
              label: "/mods/input-b.json",
              sourceKind: "translation_input",
              recordCount: 32,
              registeredAt: "2026-04-28T10:20:00Z"
            }
          ]
        }),
        selectedInputSourceId: 41
      })
    )
    const usecase = new TranslationJobSetupUseCase(gateway, store)

    await usecase.deleteInputSource(41)

    const state = store.snapshot()
    expect(gateway.deleteTranslationJobSetupInput).toHaveBeenCalledWith({
      inputSourceId: 41
    })
    expect(
      (
        gateway.getTranslationJobSetupOptions as ReturnType<
          typeof vi.fn<() => Promise<TranslationJobSetupOptionsResponse>>
        >
      ).mock.calls
    ).toHaveLength(0)
    expect(state.options?.inputCandidates).toHaveLength(1)
    expect(state.selectedInputSourceId).toBe(52)
    expect(state.options?.inputCandidates[0]?.id).toBe(52)
    expect(state.deletingInputSourceId).toBeNull()
  })

  test("deleteInputSource は削除中 id を保持し、失敗時に解除して errorMessage を残す", async () => {
    const deferredDelete =
      createDeferred<{ deletedInputSourceId?: number; errorKind?: never }>()
    const gateway = createGateway()
    gateway.deleteTranslationJobSetupInput = vi
      .fn()
      .mockImplementation(() => deferredDelete.promise)
    const store = createStore(
      createPhaseDrivenState({
        options: createPhaseOptions({
          inputCandidates: [
            {
              id: 41,
              label: "/mods/input-a.json",
              sourceKind: "translation_input",
              recordCount: 128,
              registeredAt: "2026-04-27T10:20:00Z"
            },
            {
              id: 52,
              label: "/mods/input-b.json",
              sourceKind: "translation_input",
              recordCount: 32,
              registeredAt: "2026-04-28T10:20:00Z"
            }
          ]
        }),
        selectedInputSourceId: 52
      })
    )
    const usecase = new TranslationJobSetupUseCase(gateway, store)

    const deletePromise = usecase.deleteInputSource(41)

    expect(store.snapshot().phase).toBe("ready")
    expect(store.snapshot().deletingInputSourceId).toBe(41)

    deferredDelete.reject(new Error("backend delete failed"))
    await deletePromise

    const state = store.snapshot()
    expect(state.deletingInputSourceId).toBeNull()
    expect(state.errorMessage).toBe("入力データの削除に失敗しました。")
  })

  test("deleteInputSource は削除対象が末尾でも残存候補を選ぶ", async () => {
    const gateway = createGateway()
    const store = createStore(
      createPhaseDrivenState({
        options: createPhaseOptions({
          inputCandidates: [
            {
              id: 30,
              label: "/mods/input-a.json",
              sourceKind: "translation_input",
              recordCount: 64,
              registeredAt: "2026-04-26T10:20:00Z"
            },
            {
              id: 41,
              label: "/mods/input-b.json",
              sourceKind: "translation_input",
              recordCount: 128,
              registeredAt: "2026-04-27T10:20:00Z"
            }
          ]
        }),
        selectedInputSourceId: 41
      })
    )
    const usecase = new TranslationJobSetupUseCase(gateway, store)

    await usecase.deleteInputSource(41)

    const state = store.snapshot()
    expect(state.options?.inputCandidates).toHaveLength(1)
    expect(state.options?.inputCandidates[0]?.id).toBe(30)
    expect(state.selectedInputSourceId).toBe(30)
    expect(state.deletingInputSourceId).toBeNull()
  })

  test("deleteInputSource は削除後に候補が空なら selectedInputSourceId を null にする", async () => {
    const gateway = createGateway()
    const store = createStore(
      createPhaseDrivenState({
        options: createPhaseOptions({
          inputCandidates: [
            {
              id: 41,
              label: "/mods/input-a.json",
              sourceKind: "translation_input",
              recordCount: 128,
              registeredAt: "2026-04-27T10:20:00Z"
            }
          ]
        }),
        selectedInputSourceId: 41
      })
    )
    const usecase = new TranslationJobSetupUseCase(gateway, store)

    await usecase.deleteInputSource(41)

    const state = store.snapshot()
    expect(state.options?.inputCandidates).toHaveLength(0)
    expect(state.selectedInputSourceId).toBeNull()
    expect(state.deletingInputSourceId).toBeNull()
  })

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

  test("selectedInputSourceId と一致する existingJob があっても create を gateway へ進める", async () => {
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

    expect(gateway.createTranslationJob).toHaveBeenCalledTimes(1)
    expect(gateway.createTranslationJob).toHaveBeenCalledWith(
      expect.objectContaining({
        inputSourceId: 41,
        inputSource: "/mods/input-review.json"
      })
    )
  })

  test("phase driven createJob は selectedInputSourceId と一致する existingJob があっても create を gateway へ進める", async () => {
    const gateway = createGateway()
    const store = createStore(
      createPhaseDrivenState({
        options: createPhaseOptions({
          existingJob: {
            inputSourceId: 41,
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
    const calls = gateway.createTranslationJob.mock
      .calls as CreateTranslationJobRequest[][]
    const [request] = calls[0] ?? []
    expect(request?.inputSourceId).toBe(41)
    expect(request?.inputSource).toBe("/mods/input-review.json")
    expect(request?.phaseRuntimeSelections).toHaveLength(3)
  })

  test("createJob 後の summary fetch は canStartPhase を保持して state へ反映する", async () => {
    const gateway = createGateway()
    const store = createStore()
    const usecase = new TranslationJobSetupUseCase(gateway, store)

    await usecase.createJob()

    expect(gateway.getTranslationJobSetupSummary).toHaveBeenCalledTimes(1)
    expect(gateway.getTranslationJobSetupSummary).toHaveBeenCalledWith({
      jobId: 91
    })
    expect(store.snapshot().summary).toEqual(
      expect.objectContaining({
        jobId: 91,
        canStartPhase: true
      })
    )
  })

  test("SCN-TJSPPS-006: provider 変更後は model 未選択へ戻し遅延 model list を混入させない", async () => {
    const delayedModels =
      createDeferred<ListTranslationJobSetupProviderModelsResponse>()
    const gateway = createGateway()
    gateway.validateTranslationJobSetup = vi
      .fn()
      .mockResolvedValue(createValidationResult())
    gateway.listTranslationJobSetupProviderModels = vi
      .fn()
      .mockImplementationOnce(() => delayedModels.promise)
      .mockResolvedValueOnce({
        phaseId: "word_translation",
        provider: "xai",
        credentialStatus: "configured",
        requestToken: "req-xai-2",
        sourceToken: "xai-latest-source",
        status: "success",
        models: [{ modelId: "xai-latest-model", label: "xai-latest-model" }]
      })
    const store = createStore(createPhaseDrivenState())
    const usecase = new TranslationJobSetupUseCase(gateway, store)

    const refreshPromise = usecase.refreshPhaseModels("word_translation")
    expect(
      store
        .snapshot()
        .providerModelLists?.find(
          (entry) => entry.phaseId === "word_translation"
        )
    ).toEqual(
      expect.objectContaining({
        provider: "gemini",
        status: "loading"
      })
    )

    usecase.selectPhaseProvider("word_translation", "xai")
    expect(
      store
        .snapshot()
        .phaseRuntimeSelections?.find(
          (selection) => selection.phaseId === "word_translation"
        )
    ).toEqual(
      expect.objectContaining({
        provider: "xai",
        model: "",
        modelListSourceToken: ""
      })
    )

    delayedModels.resolve({
      phaseId: "word_translation",
      provider: "gemini",
      credentialStatus: "configured",
      requestToken:
        store
          .snapshot()
          .providerModelLists?.find(
            (entry) => entry.phaseId === "word_translation"
          )?.requestToken ?? "",
      sourceToken: "gemini-delayed-source",
      status: "success",
      models: [
        { modelId: "gemini-delayed-model", label: "gemini-delayed-model" }
      ]
    })
    await refreshPromise

    const state = store.snapshot()
    expect(
      state.phaseRuntimeSelections?.find(
        (selection) => selection.phaseId === "word_translation"
      )
    ).toEqual(
      expect.objectContaining({
        provider: "xai",
        model: "",
        modelListSourceToken: ""
      })
    )
    expect(
      state.providerModelLists?.find(
        (entry) => entry.phaseId === "word_translation"
      )
    ).toEqual(
      expect.objectContaining({
        provider: "xai",
        status: "not_updated",
        models: []
      })
    )
    const cardAfterProviderSwitch = state.modelSettingsCards?.find(
      (card) => card.referenceId === "word_translation"
    )
    expect(cardAfterProviderSwitch?.provider).toBe("xai")
    expect(cardAfterProviderSwitch?.model).toBe("")
    expect(cardAfterProviderSwitch?.modelList.provider).toBe("xai")
    expect(cardAfterProviderSwitch?.modelList.status).toBe("not_updated")
    expect(cardAfterProviderSwitch?.modelList.models).toEqual([])

    await usecase.createJob()
    expect(gateway.createTranslationJob).not.toHaveBeenCalled()
    expect(store.snapshot().errorMessage).toBe(
      "翻訳ジョブを作成するには、3 つの翻訳段階で不足を解消してください。"
    )
  })

  test("SCN-TJSPPS-006: 古い failed 応答は最新 providerModelLists と model を巻き戻さない", async () => {
    const delayedModels =
      createDeferred<ListTranslationJobSetupProviderModelsResponse>()
    const gateway = createGateway()
    gateway.validateTranslationJobSetup = vi
      .fn()
      .mockResolvedValue(createValidationResult())
    gateway.listTranslationJobSetupProviderModels = vi
      .fn()
      .mockImplementationOnce(() => delayedModels.promise)
      .mockResolvedValue({
        phaseId: "word_translation",
        provider: "xai",
        credentialStatus: "configured",
        requestToken: "xai-fast",
        sourceToken: "word_translation|xai|xai-primary|xai-fast",
        status: "success",
        models: [{ modelId: "xai-word-model", label: "xai-word-model" }]
      })
    const store = createStore(createPhaseDrivenState())
    const usecase = new TranslationJobSetupUseCase(gateway, store)

    const staleRefresh = usecase.refreshPhaseModels("word_translation")
    usecase.selectPhaseProvider("word_translation", "xai")
    await usecase.refreshPhaseModels("word_translation")

    delayedModels.resolve({
      phaseId: "word_translation",
      provider: "gemini",
      credentialStatus: "configured",
      requestToken:
        store
          .snapshot()
          .providerModelLists?.find(
            (entry) => entry.phaseId === "word_translation"
          )?.requestToken ?? "",
      sourceToken: "gemini-stale-failed",
      status: "failed",
      models: []
    })
    await staleRefresh

    const state = store.snapshot()
    expect(
      state.phaseRuntimeSelections?.find(
        (selection) => selection.phaseId === "word_translation"
      )
    ).toEqual(
      expect.objectContaining({
        provider: "xai",
        model: "",
        modelListSourceToken: ""
      })
    )
    expect(
      state.providerModelLists?.find(
        (entry) => entry.phaseId === "word_translation"
      )
    ).toEqual(
      expect.objectContaining({
        provider: "xai"
      })
    )
    const cardAfterStaleFailed = state.modelSettingsCards?.find(
      (card) => card.referenceId === "word_translation"
    )
    expect(cardAfterStaleFailed?.provider).toBe("xai")
    expect(cardAfterStaleFailed?.model).toBe("")
    expect(cardAfterStaleFailed?.modelList.provider).toBe("xai")
    expect(state.validationState).toBe("stale")
  })

  test("state-invariant-001: provider 変更後の遅延 model list 失敗は現在 phase state を上書きしない", async () => {
    const delayedModels =
      createDeferred<ListTranslationJobSetupProviderModelsResponse>()
    const gateway = createGateway()
    gateway.validateTranslationJobSetup = vi
      .fn()
      .mockResolvedValue(createValidationResult())
    gateway.listTranslationJobSetupProviderModels = vi.fn(
      () => delayedModels.promise
    )
    const store = createStore(createPhaseDrivenState())
    const usecase = new TranslationJobSetupUseCase(gateway, store)

    const refreshPromise = usecase.refreshPhaseModels("word_translation")
    usecase.selectPhaseProvider("word_translation", "xai")

    delayedModels.reject(new Error("delayed model list failure"))
    await refreshPromise

    const state = store.snapshot()
    expect(
      state.phaseRuntimeSelections?.find(
        (selection) => selection.phaseId === "word_translation"
      )
    ).toEqual(
      expect.objectContaining({
        provider: "xai",
        model: "",
        modelListSourceToken: ""
      })
    )
    expect(
      state.providerModelLists?.find(
        (entry) => entry.phaseId === "word_translation"
      )
    ).toEqual(
      expect.objectContaining({
        provider: "xai",
        status: "not_updated",
        requestToken: "",
        models: []
      })
    )
    expect(state.errorMessage).toBe("")
  })

  test("LM Studio の credential_not_required でもモデル一覧取得後に model を選択できる", async () => {
    const gateway = createGateway()
    const listProviderModelsSpy = vi.fn().mockResolvedValue({
      phaseId: "text_translation",
      provider: "lm_studio",
      credentialStatus: "not_required",
      requestToken: "req-lm-1",
      sourceToken: "text_translation|lm_studio||req-lm-1",
      status: "credential_not_required",
      models: [{ modelId: "lmstudio-community", label: "LM Studio Community" }]
    })
    gateway.validateTranslationJobSetup = vi
      .fn()
      .mockResolvedValue(createValidationResult())
    gateway.listTranslationJobSetupProviderModels = listProviderModelsSpy
    const store = createStore(createPhaseDrivenState())
    const usecase = new TranslationJobSetupUseCase(gateway, store)

    usecase.selectPhaseProvider("text_translation", "lm_studio")
    await usecase.refreshPhaseModels("text_translation")
    usecase.selectPhaseModel("text_translation", "lmstudio-community")

    expect(listProviderModelsSpy).toHaveBeenCalledTimes(1)
    expect(listProviderModelsSpy).toHaveBeenCalledWith(
      expect.objectContaining({
        phaseId: "text_translation",
        provider: "lm_studio",
        credentialRef: "",
        credentialStatus: "not_required"
      })
    )
    expect(
      store
        .snapshot()
        .phaseRuntimeSelections?.find(
          (selection) => selection.phaseId === "text_translation"
        )
    ).toEqual(
      expect.objectContaining({
        provider: "lm_studio",
        credentialRef: "",
        credentialStatus: "not_required",
        model: "lmstudio-community",
        modelListSourceToken: "text_translation|lm_studio||req-lm-1"
      })
    )
  })

  test("credentialStatus が missing でも refreshPhaseModels は gateway 応答で model list と選択状態を更新する", async () => {
    const gateway = createGateway()
    const listProviderModelsSpy = vi.fn(
      (
        request: ListTranslationJobSetupProviderModelsRequest
      ): Promise<ListTranslationJobSetupProviderModelsResponse> =>
        Promise.resolve({
          phaseId: request.phaseId,
          provider: "gemini",
          credentialStatus: "not_required",
          requestToken: request.requestToken,
          sourceToken: "word_translation|gemini||req-test-safe",
          status: "success",
          models: [{ modelId: "gemini-test-safe", label: "Gemini Test Safe" }]
        })
    )
    gateway.listTranslationJobSetupProviderModels = listProviderModelsSpy
    const store = createStore(
      createPhaseDrivenState({
        phaseRuntimeSelections: [
          {
            phaseId: "word_translation",
            provider: "gemini",
            model: "",
            credentialRef: "gemini-primary",
            credentialStatus: "missing",
            executionMode: "sync",
            batchMode: "enabled",
            modelListSourceToken: ""
          },
          {
            phaseId: "npc_persona_generation",
            provider: "xai",
            model: "xai-persona-model",
            credentialRef: "xai-primary",
            credentialStatus: "configured",
            executionMode: "sync",
            batchMode: "disabled",
            modelListSourceToken: "xai-source-current"
          },
          {
            phaseId: "text_translation",
            provider: "xai",
            model: "xai-text-model",
            credentialRef: "xai-primary",
            credentialStatus: "configured",
            executionMode: "sync",
            batchMode: "disabled",
            modelListSourceToken: "xai-source-text"
          }
        ]
      })
    )
    const usecase = new TranslationJobSetupUseCase(gateway, store)

    await usecase.refreshPhaseModels("word_translation")

    expect(listProviderModelsSpy).toHaveBeenCalledTimes(1)
    expect(listProviderModelsSpy).toHaveBeenCalledWith(
      expect.objectContaining({
        phaseId: "word_translation",
        provider: "gemini",
        credentialStatus: "missing"
      })
    )
    expect(
      store
        .snapshot()
        .providerModelLists?.find(
          (entry) => entry.phaseId === "word_translation"
        )
    ).toEqual(
      expect.objectContaining({
        provider: "gemini",
        credentialStatus: "not_required",
        status: "success",
        sourceToken: "word_translation|gemini||req-test-safe",
        models: [{ modelId: "gemini-test-safe", label: "Gemini Test Safe" }]
      })
    )
    expect(
      store
        .snapshot()
        .phaseRuntimeSelections?.find(
          (selection) => selection.phaseId === "word_translation"
        )
    ).toEqual(
      expect.objectContaining({
        provider: "gemini",
        model: "gemini-test-safe",
        modelListSourceToken: "word_translation|gemini||req-test-safe"
      })
    )
  })

  test("単一モデル一覧は固定名に依存せず唯一の modelId を選択状態にする", async () => {
    const gateway = createGateway()
    gateway.listTranslationJobSetupProviderModels = vi.fn(
      (
        request: ListTranslationJobSetupProviderModelsRequest
      ): Promise<ListTranslationJobSetupProviderModelsResponse> =>
        Promise.resolve({
          phaseId: "text_translation",
          provider: "lm_studio",
          credentialStatus: "not_required",
          requestToken: request.requestToken,
          sourceToken: "text_translation|lm_studio||req-fake-1",
          status: "success",
          models: [{ modelId: "single-available-model", label: "single-available-model" }]
        })
    )
    const store = createStore(createPhaseDrivenState())
    const usecase = new TranslationJobSetupUseCase(gateway, store)

    usecase.selectPhaseProvider("text_translation", "lm_studio")
    await usecase.refreshPhaseModels("text_translation")

    expect(
      store
        .snapshot()
        .phaseRuntimeSelections?.find(
          (selection) => selection.phaseId === "text_translation"
        )
    ).toEqual(
      expect.objectContaining({
        provider: "lm_studio",
        credentialStatus: "not_required",
        model: "single-available-model",
        modelListSourceToken: "text_translation|lm_studio||req-fake-1"
      })
    )
  })

  test("空 provider の phase draft は先頭 provider capability へ正規化する", async () => {
    const gateway = createGateway()
    const store = createStore(createState({ options: null }))
    const usecase = new TranslationJobSetupUseCase(gateway, store)
    gateway.getTranslationJobSetupOptions.mockResolvedValueOnce(
      createPhaseOptions({
        phaseRuntimeDrafts: [
          {
            phaseId: "word_translation",
            provider: "",
            model: "",
            credentialRef: "",
            credentialStatus: "missing",
            executionMode: "sync",
            batchMode: "unsupported",
            modelListSourceToken: ""
          }
        ]
      })
    )

    await usecase.load()

    expect(
      store
        .snapshot()
        .phaseRuntimeSelections?.find(
          (selection) => selection.phaseId === "word_translation"
        )
    ).toEqual(expect.objectContaining({ provider: "gemini" }))
  })

  test("contract-004: LM Studio validation payload は空 credentialRef を使う", async () => {
    const gateway = createGateway()
    gateway.listTranslationJobSetupProviderModels = vi.fn().mockResolvedValue({
      phaseId: "text_translation",
      provider: "lm_studio",
      credentialStatus: "not_required",
      requestToken: "req-lm-1",
      sourceToken: "text_translation|lm_studio||req-lm-1",
      status: "credential_not_required",
      models: [{ modelId: "lmstudio-community", label: "LM Studio Community" }]
    })
    const validateSpy = vi.fn().mockResolvedValue(createValidationResult())
    gateway.validateTranslationJobSetup = validateSpy
    const store = createStore(createPhaseDrivenState())
    const usecase = new TranslationJobSetupUseCase(gateway, store)

    usecase.selectPhaseProvider("text_translation", "lm_studio")
    await usecase.refreshPhaseModels("text_translation")
    usecase.selectPhaseModel("text_translation", "lmstudio-community")
    await usecase.runValidation()

    const lastPayload = validateSpy.mock.calls.at(-1)?.[0] as
      | ValidateTranslationJobSetupRequest
      | undefined
    expect(lastPayload).toBeDefined()
    expect(lastPayload).toEqual(
      expect.objectContaining({
        credentialRef: "gemini-primary"
      })
    )
    expect(lastPayload?.phaseRuntimeSelections).toEqual(
      expect.arrayContaining([
        expect.objectContaining({
          phaseId: "text_translation",
          provider: "lm_studio",
          credentialRef: "",
          model: "lmstudio-community",
          modelListSourceToken: "text_translation|lm_studio||req-lm-1"
        })
      ])
    )
  })
})
