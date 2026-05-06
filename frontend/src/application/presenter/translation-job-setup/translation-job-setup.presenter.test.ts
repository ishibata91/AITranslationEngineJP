import { describe, expect, test } from "vitest"

import type {
  TranslationJobSetupOptionsResponse,
  TranslationJobSetupScreenState,
  TranslationJobSetupValidationResponse
} from "@application/gateway-contract/translation-job-setup"

import { TranslationJobSetupPresenter } from "./translation-job-setup.presenter"

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
      },
      {
        id: 42,
        label: "/mods/other-input-review.json",
        sourceKind: "xEdit extract",
        recordCount: 64,
        registeredAt: "2026-04-27T10:25:00Z"
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

function createValidationResult(
  overrides: Partial<TranslationJobSetupValidationResponse> = {}
): TranslationJobSetupValidationResponse {
  return {
    status: "pass",
    targetSlices: ["input", "runtime", "credentials"],
    validatedAt: "2026-04-27T10:30:00Z",
    canCreate: true,
    passSlices: ["input", "runtime", "credentials"],
    ...overrides
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

describe("TranslationJobSetupPresenter", () => {
  test("selectedInputSourceId が異なる existingJob は inputSource 表示名が一致しても create を無効化しない", () => {
    const presenter = new TranslationJobSetupPresenter()

    const viewModel = presenter.toViewModel(
      createState({
        options: createOptions({
          existingJob: {
            inputSourceId: 42,
            jobId: 300,
            status: "ready",
            inputSource: "/mods/input-review.json"
          }
        })
      }),
      true
    )

    expect(viewModel.canCreate).toBe(true)
    expect(viewModel.blockedReasons).not.toContain(
      "既存 job があるため create を無効化しています。"
    )
  })

  test("selectedInputSourceId と一致する existingJob は参考表示だけを残し canCreate を維持する", () => {
    const presenter = new TranslationJobSetupPresenter()

    const viewModel = presenter.toViewModel(
      createState({
        options: createOptions({
          existingJob: {
            inputSourceId: 41,
            jobId: 300,
            status: "ready",
            inputSource: "Imported from another label"
          }
        })
      }),
      true
    )

    expect(viewModel.canCreate).toBe(true)
    expect(viewModel.blockedReasons).not.toContain(
      "既存 job があるため create を無効化しています。"
    )
    expect(viewModel.existingJobSummary).toContain("job #300")
  })

  test("phase-driven state では selectedInputSourceId と一致する existingJob があっても canCreate を false にしない", () => {
    const presenter = new TranslationJobSetupPresenter()

    const viewModel = presenter.toViewModel(
      createState({
        options: createOptions({
          existingJob: {
            inputSourceId: 41,
            jobId: 300,
            status: "ready",
            inputSource: "/mods/input-review.json"
          },
          providerCapabilities: [
            {
              provider: "openai",
              credentialRequirement: "required",
              supportedExecutionModes: ["batch"],
              supportsBatchMode: true
            }
          ]
        }),
        phaseRuntimeSelections: [
          {
            phaseId: "word_translation",
            provider: "openai",
            model: "gpt-5.4-mini",
            credentialRef: "openai-primary",
            credentialStatus: "configured",
            executionMode: "batch",
            batchMode: "enabled",
            modelListSourceToken: "word-1"
          },
          {
            phaseId: "npc_persona_generation",
            provider: "openai",
            model: "gpt-5.4-mini",
            credentialRef: "openai-primary",
            credentialStatus: "configured",
            executionMode: "batch",
            batchMode: "enabled",
            modelListSourceToken: "npc-1"
          },
          {
            phaseId: "text_translation",
            provider: "openai",
            model: "gpt-5.4-mini",
            credentialRef: "openai-primary",
            credentialStatus: "configured",
            executionMode: "batch",
            batchMode: "enabled",
            modelListSourceToken: "text-1"
          }
        ],
        providerModelLists: [
          {
            phaseId: "word_translation",
            provider: "openai",
            credentialStatus: "configured",
            requestToken: "word-request",
            sourceToken: "word-1",
            status: "success",
            models: [{ modelId: "gpt-5.4-mini", label: "GPT-5.4 Mini" }]
          },
          {
            phaseId: "npc_persona_generation",
            provider: "openai",
            credentialStatus: "configured",
            requestToken: "npc-request",
            sourceToken: "npc-1",
            status: "success",
            models: [{ modelId: "gpt-5.4-mini", label: "GPT-5.4 Mini" }]
          },
          {
            phaseId: "text_translation",
            provider: "openai",
            credentialStatus: "configured",
            requestToken: "text-request",
            sourceToken: "text-1",
            status: "success",
            models: [{ modelId: "gpt-5.4-mini", label: "GPT-5.4 Mini" }]
          }
        ],
        validationResult: createValidationResult({
          phaseResults: [
            {
              phaseId: "word_translation",
              status: "pass",
              canCreate: true,
              modelListState: "success",
              modelListSourceToken: "word-1",
              isModelSelectionStale: false
            },
            {
              phaseId: "npc_persona_generation",
              status: "pass",
              canCreate: true,
              modelListState: "success",
              modelListSourceToken: "npc-1",
              isModelSelectionStale: false
            },
            {
              phaseId: "text_translation",
              status: "pass",
              canCreate: true,
              modelListState: "success",
              modelListSourceToken: "text-1",
              isModelSelectionStale: false
            }
          ]
        })
      }),
      true
    )

    expect(viewModel.canCreate).toBe(true)
    expect(viewModel.blockedReasons).not.toContain(
      "同じ入力データの翻訳ジョブが既にあります。"
    )
    expect(viewModel.existingJobSummary).toContain("job #300")
  })
})
