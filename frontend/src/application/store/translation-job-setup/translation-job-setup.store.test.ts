import { describe, expect, it } from "vitest"

import type {
  TranslationJobSetupOptionsResponse,
  TranslationJobSetupSummaryResponse,
  TranslationJobSetupValidationResponse
} from "@application/gateway-contract/translation-job-setup"

import { TranslationJobSetupStore } from "./translation-job-setup.store"

function createOptions(): TranslationJobSetupOptionsResponse {
  return {
    inputCandidates: [
      {
        id: 1,
        label: "Test Input A",
        sourceKind: "xedit-json",
        registeredAt: "2026-05-05T00:00:00Z",
        recordCount: 2
      }
    ],
    existingJob: {
      inputSourceId: 1,
      jobId: 10,
      status: "ready",
      inputSource: "TestInputPluginA.esp"
    },
    sharedDictionaries: [
      {
        id: "dict-a",
        label: "Test Dictionary A"
      }
    ],
    sharedPersonas: [
      {
        id: "persona-a",
        label: "Test Persona A"
      }
    ],
    aiRuntimeOptions: [
      {
        provider: "fake",
        model: "test-model-a",
        mode: "sync"
      }
    ],
    credentialRefs: [
      {
        provider: "fake",
        credentialRef: "test-ref-a",
        isConfigured: true,
        isMissingSecret: false
      }
    ],
    providerCapabilities: [
      {
        provider: "fake",
        credentialRequirement: "not_required",
        supportedExecutionModes: ["sync"],
        supportsBatchMode: false
      }
    ],
    phaseRuntimeDrafts: [
      {
        phaseId: "word_translation",
        provider: "fake",
        model: "test-model-a",
        credentialRef: "test-ref-a",
        credentialStatus: "configured",
        executionMode: "sync",
        batchMode: "disabled",
        modelListSourceToken: "token-a"
      }
    ]
  }
}

function createValidation(): TranslationJobSetupValidationResponse {
  return {
    status: "ready",
    targetSlices: ["input"],
    validatedAt: "2026-05-05T00:00:00Z",
    canCreate: true,
    passSlices: ["input"],
    phaseResults: [
      {
        phaseId: "word_translation",
        status: "ready",
        canCreate: true,
        modelListState: "success",
        modelListSourceToken: "token-a",
        isModelSelectionStale: false
      }
    ],
    staleModelListPhaseIds: ["text_translation"]
  }
}

function createSummary(): TranslationJobSetupSummaryResponse {
  return {
    jobId: 10,
    jobState: "ready",
    inputSource: "TestInputPluginA.esp",
    executionSummary: {
      provider: "fake",
      model: "test-model-a",
      executionMode: "sync"
    },
    validationPassSlices: ["input"],
    canStartPhase: true,
    phaseRuntimeSummaries: [
      {
        phaseId: "word_translation",
        provider: "fake",
        model: "test-model-a",
        credentialRef: "test-ref-a",
        credentialStatus: "configured",
        executionMode: "sync",
        batchMode: "disabled",
        modelListSourceToken: "token-a"
      }
    ]
  }
}

describe("TranslationJobSetupStore", () => {
  it("subscribe immediately notifies current state", () => {
    const store = new TranslationJobSetupStore()
    const snapshots: unknown[] = []

    store.subscribe((state) => {
      snapshots.push(state)
    })

    expect(snapshots).toHaveLength(1)
    expect(store.snapshot()).toMatchObject({
      phase: "idle",
      options: null,
      selectedInputSourceId: null,
      selectedRuntimeKey: null,
      selectedCredentialRef: "",
      phaseRuntimeSelections: [],
      providerModelLists: [],
      validationResult: null,
      validationState: "not-run",
      dirty: false,
      errorMessage: "",
      createErrorKind: null,
      summary: null
    })
  })

  it("unsubscribe removes listener from future updates", () => {
    const store = new TranslationJobSetupStore()
    let callCount = 0

    const unsubscribe = store.subscribe(() => {
      callCount += 1
    })
    unsubscribe()

    store.update((draft) => {
      draft.phase = "loading"
    })

    expect(callCount).toBe(1)
  })

  it("snapshot returns defensive copies", () => {
    const store = new TranslationJobSetupStore()

    store.update((draft) => {
      draft.options = createOptions()
      draft.phaseRuntimeSelections = createOptions().phaseRuntimeDrafts
      draft.providerModelLists = [
        {
          phaseId: "word_translation",
          provider: "fake",
          credentialStatus: "configured",
          requestToken: "request-a",
          sourceToken: "source-a",
          status: "success",
          models: [
            {
              modelId: "test-model-a",
              label: "Test Model A"
            }
          ]
        }
      ]
      draft.validationResult = createValidation()
      draft.summary = createSummary()
    })

    const snapshot = store.snapshot()
    snapshot.options!.inputCandidates[0].label = "changed"
    snapshot.options!.existingJob!.status = "changed"
    snapshot.options!.sharedDictionaries[0].label = "changed"
    snapshot.options!.sharedPersonas[0].label = "changed"
    snapshot.options!.aiRuntimeOptions[0].model = "changed"
    snapshot.options!.credentialRefs[0].credentialRef = "changed"
    snapshot.options!.providerCapabilities![0].supportedExecutionModes[0] =
      "changed"
    snapshot.options!.phaseRuntimeDrafts![0].model = "changed"
    snapshot.phaseRuntimeSelections![0].model = "changed"
    snapshot.providerModelLists![0].models[0].label = "changed"
    snapshot.validationResult!.targetSlices[0] = "changed"
    snapshot.validationResult!.passSlices[0] = "changed"
    snapshot.validationResult!.phaseResults![0].status = "changed"
    snapshot.validationResult!.staleModelListPhaseIds![0] = "word_translation"
    snapshot.summary!.executionSummary.model = "changed"
    snapshot.summary!.validationPassSlices[0] = "changed"
    snapshot.summary!.phaseRuntimeSummaries![0].model = "changed"

    const nextSnapshot = store.snapshot()
    expect(nextSnapshot.options?.inputCandidates[0]?.label).toBe("Test Input A")
    expect(nextSnapshot.options?.existingJob?.status).toBe("ready")
    expect(nextSnapshot.options?.sharedDictionaries[0]?.label).toBe(
      "Test Dictionary A"
    )
    expect(nextSnapshot.options?.sharedPersonas[0]?.label).toBe(
      "Test Persona A"
    )
    expect(nextSnapshot.options?.aiRuntimeOptions[0]?.model).toBe(
      "test-model-a"
    )
    expect(nextSnapshot.options?.credentialRefs[0]?.credentialRef).toBe(
      "test-ref-a"
    )
    expect(
      nextSnapshot.options?.providerCapabilities?.[0]?.supportedExecutionModes[0]
    ).toBe("sync")
    expect(nextSnapshot.options?.phaseRuntimeDrafts?.[0]?.model).toBe(
      "test-model-a"
    )
    expect(nextSnapshot.phaseRuntimeSelections?.[0]?.model).toBe("test-model-a")
    expect(nextSnapshot.providerModelLists?.[0]?.models[0]?.label).toBe(
      "Test Model A"
    )
    expect(nextSnapshot.validationResult?.targetSlices[0]).toBe("input")
    expect(nextSnapshot.validationResult?.passSlices[0]).toBe("input")
    expect(nextSnapshot.validationResult?.phaseResults?.[0]?.status).toBe(
      "ready"
    )
    expect(nextSnapshot.validationResult?.staleModelListPhaseIds?.[0]).toBe(
      "text_translation"
    )
    expect(nextSnapshot.summary?.executionSummary.model).toBe("test-model-a")
    expect(nextSnapshot.summary?.validationPassSlices[0]).toBe("input")
    expect(nextSnapshot.summary?.phaseRuntimeSummaries?.[0]?.model).toBe(
      "test-model-a"
    )
  })
})
