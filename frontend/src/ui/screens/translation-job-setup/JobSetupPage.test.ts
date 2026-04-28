import { render, screen, waitFor } from "@testing-library/svelte"
import userEvent from "@testing-library/user-event"
import { describe, expect, test, vi } from "vitest"

import type {
  TranslationJobSetupScreenControllerContract,
  TranslationJobSetupScreenViewModelListener
} from "@application/contract/translation-job-setup/translation-job-setup-screen-contract"
import {
  createTranslationJobSetupRuntimeKey,
  TranslationJobSetupOptionsResponse,
  TranslationJobSetupScreenState,
  TranslationJobSetupScreenViewModel,
  TranslationJobSetupSummaryResponse,
  TranslationJobSetupValidationResponse
} from "@application/gateway-contract/translation-job-setup"
import { TranslationJobSetupPresenter } from "@application/presenter/translation-job-setup"
import JobSetupPage from "@ui/screens/translation-job-setup/JobSetupPage.svelte"

function createOptions(
  overrides: Partial<TranslationJobSetupOptionsResponse> = {}
): TranslationJobSetupOptionsResponse {
  return {
    inputCandidates: [
      {
        id: 41,
        label: "/mods/very/long/path/translation/input-review-export.json",
        sourceKind: "xEdit extract",
        recordCount: 128
      }
    ],
    existingJob: undefined,
    sharedDictionaries: [
      { id: "dict-core", label: "Shared Dictionary / Foundation Core" }
    ],
    sharedPersonas: [
      { id: "persona-core", label: "Foundation Persona / Translation Main" }
    ],
    aiRuntimeOptions: [
      {
        provider: "openai-compatible",
        model: "gpt-4.1-mini-preview-with-a-very-long-name",
        mode: "batch"
      },
      {
        provider: "anthropic",
        model: "claude-3-7-sonnet-with-a-very-long-name",
        mode: "sync"
      }
    ],
    credentialRefs: [
      {
        provider: "openai-compatible",
        credentialRef: "cred-main",
        isConfigured: true,
        isMissingSecret: false
      },
      {
        provider: "anthropic",
        credentialRef: "cred-missing",
        isConfigured: true,
        isMissingSecret: true
      }
    ],
    ...overrides
  }
}

function createValidationResult(
  overrides: Partial<TranslationJobSetupValidationResponse> = {}
): TranslationJobSetupValidationResponse {
  return {
    status: "warning",
    blockingFailureCategory: "cache missing",
    targetSlices: ["credential", "runtime"],
    validatedAt: "invalid-timestamp",
    canCreate: false,
    passSlices: ["input", "foundation"],
    ...overrides
  }
}

function createSummary(
  overrides: Partial<TranslationJobSetupSummaryResponse> = {}
): TranslationJobSetupSummaryResponse {
  return {
    jobId: 501,
    jobState: "ready",
    inputSource: "/mods/very/long/path/translation/input-review-export.json",
    canStartPhase: true,
    executionSummary: {
      provider: "openai-compatible",
      model: "gpt-4.1-mini-preview-with-a-very-long-name",
      executionMode: "batch"
    },
    validationPassSlices: ["input", "runtime", "credential"],
    ...overrides
  }
}

function createState(
  overrides: Partial<TranslationJobSetupScreenState> = {}
): TranslationJobSetupScreenState {
  const options = overrides.options ?? createOptions()
  const selectedRuntimeOption = options.aiRuntimeOptions[0] ?? null

  return {
    phase: "ready",
    options,
    selectedInputSourceId: options.inputCandidates[0]?.id ?? null,
    selectedRuntimeKey: selectedRuntimeOption
      ? createTranslationJobSetupRuntimeKey(selectedRuntimeOption)
      : null,
    selectedCredentialRef: options.credentialRefs[0]?.credentialRef ?? "",
    validationResult: createValidationResult(),
    validationState: "stale",
    dirty: true,
    errorMessage: "",
    createErrorKind: null,
    summary: null,
    ...overrides
  }
}

function createPresentedViewModel(
  overrides: Partial<TranslationJobSetupScreenState> = {}
): TranslationJobSetupScreenViewModel {
  const presenter = new TranslationJobSetupPresenter()
  return presenter.toViewModel(createState(overrides), true)
}

function createViewModel(
  overrides: Partial<TranslationJobSetupScreenViewModel> = {}
): TranslationJobSetupScreenViewModel {
  const options = overrides.options ?? createOptions()
  const validationResult = overrides.validationResult ?? createValidationResult()
  const summary = overrides.summary ?? null

  return {
    phase: summary ? "summary" : "ready",
    options,
    selectedInputSourceId: 41,
    selectedRuntimeKey: "openai-compatible::gpt-4.1-mini-preview-with-a-very-long-name::batch",
    selectedCredentialRef: "cred-main",
    validationResult,
    validationState: summary ? "fresh" : "stale",
    dirty: !summary,
    errorMessage: "",
    createErrorKind: null,
    summary,
    gatewayStatus: "接続準備済み",
    selectedInputCandidate: options.inputCandidates[0] ?? null,
    selectedRuntimeOption: options.aiRuntimeOptions[0] ?? null,
    availableCredentialRefs: options.credentialRefs.filter(
      (credential) => credential.provider === "openai-compatible"
    ),
    selectedInputLabel: options.inputCandidates[0]?.label ?? "未選択",
    selectedInputSourceKind: options.inputCandidates[0]?.sourceKind ?? "-",
    selectedInputRecordCountLabel: "128 件",
    selectedInputRegisteredAtLabel: "2026/4/27 9:30:00",
    existingJobSummary: options.existingJob
      ? `job #${options.existingJob.jobId} / ${options.existingJob.status} / ${options.existingJob.inputSource}`
      : "既存 job はありません。",
    dictionaryLabels: options.sharedDictionaries.map((item) => item.label),
    personaLabels: options.sharedPersonas.map((item) => item.label),
    validationStatusLabel: summary ? "validation pass" : "validation warning",
    validationStatusText: summary
      ? "validation pass / 対象断面: input / runtime / credential"
      : "設定を変更したため validation が失効しました。create 前に再実行が必要です。",
    createStatusText: summary
      ? "create 成功済みです。ready job summary を read-only で表示しています。"
      : "validation が fresh かつ create 可能な時だけ job を作成できます。",
    blockedReasons: summary
      ? []
      : [
          "validation が失効しています。",
          "blocking failure を解消するまで create できません。"
        ],
    canValidate: !summary,
    canCreate: false,
    isLoading: false,
    isValidating: false,
    isCreating: false,
    hasExistingJob: Boolean(options.existingJob),
    showCacheMissingGuidance: !summary,
    credentialStateText: "credential 参照は設定済みです。",
    ...overrides
  }
}

class TranslationJobSetupScreenControllerFake
  implements TranslationJobSetupScreenControllerContract
{
  private viewModel: TranslationJobSetupScreenViewModel

  private readonly listeners = new Set<TranslationJobSetupScreenViewModelListener>()

  readonly mount = vi.fn(async () => {})
  readonly dispose = vi.fn(() => {})
  readonly selectInputSource = vi.fn(() => {})
  readonly selectRuntime = vi.fn(() => {})
  readonly selectCredentialRef = vi.fn(() => {})
  readonly runValidation = vi.fn(async () => {})
  readonly createJob = vi.fn(async () => {})

  constructor(initialViewModel = createViewModel()) {
    this.viewModel = initialViewModel
  }

  subscribe(listener: TranslationJobSetupScreenViewModelListener): () => void {
    this.listeners.add(listener)
    return () => {
      this.listeners.delete(listener)
    }
  }

  getViewModel(): TranslationJobSetupScreenViewModel {
    return this.viewModel
  }

  pushViewModel(nextViewModel: TranslationJobSetupScreenViewModel): void {
    this.viewModel = nextViewModel
    for (const listener of this.listeners) {
      listener(nextViewModel)
    }
  }
}

describe("JobSetupPage", () => {
  test("input metadata の registeredAt supplied 値を表示する", () => {
    const registeredAt = "2026-04-27T00:30:00.000Z"
    const controller = new TranslationJobSetupScreenControllerFake(
      createPresentedViewModel({
        options: createOptions({
          inputCandidates: [
            {
              id: 41,
              label: "/mods/very/long/path/translation/input-review-export.json",
              sourceKind: "xEdit extract",
              registeredAt,
              recordCount: 128
            }
          ]
        })
      })
    )

    render(JobSetupPage, {
      props: {
        createController: () => controller
      }
    })

    expect(screen.getByText(new Date(registeredAt).toLocaleString("ja-JP"))).toBeInTheDocument()
  })

  test("入力、基盤参照、validation 状態、create 無効条件、cache missing 戻り導線を表示する", async () => {
    const user = userEvent.setup()
    const onReturnToInputReview = vi.fn()
    const controller = new TranslationJobSetupScreenControllerFake()

    render(JobSetupPage, {
      props: {
        createController: () => controller,
        onReturnToInputReview
      }
    })

    expect(screen.getByRole("heading", { level: 2, name: "Job Setup" })).toBeInTheDocument()
    expect(
      screen.getAllByText("/mods/very/long/path/translation/input-review-export.json").length
    ).toBeGreaterThan(0)
    expect(screen.getByText("xEdit extract")).toBeInTheDocument()
    expect(screen.getByText("2026/4/27 9:30:00")).toBeInTheDocument()
    expect(screen.getByText("128 件")).toBeInTheDocument()
    expect(screen.getByText("既存 job はありません。")).toBeInTheDocument()
    expect(screen.getByText("Shared Dictionary / Foundation Core")).toBeInTheDocument()
    expect(screen.getByText("Foundation Persona / Translation Main")).toBeInTheDocument()
    expect(screen.getByText("credential 参照は設定済みです。")).toBeInTheDocument()
    expect(screen.getAllByText("validation warning").length).toBeGreaterThan(0)
    expect(screen.getByText("invalid-timestamp")).toBeInTheDocument()
    expect(screen.getByText("cache missing")).toBeInTheDocument()
    expect(screen.getByText("dirty")).toBeInTheDocument()
    expect(screen.getAllByText("credential").length).toBeGreaterThan(0)
    expect(screen.getAllByText("runtime").length).toBeGreaterThan(0)
    expect(screen.getAllByText("input").length).toBeGreaterThan(0)
    expect(screen.getAllByText("foundation").length).toBeGreaterThan(0)
    expect(screen.getByRole("button", { name: "ready job を作成" })).toBeDisabled()
    expect(screen.getByText("validation が失効しています。")).toBeInTheDocument()
    expect(screen.getByText("blocking failure を解消するまで create できません。")).toBeInTheDocument()
    expect(
      screen.getByText("cache missing は Job Setup で再構築しません。Input Review の再構築導線へ戻ってください。")
    ).toBeInTheDocument()

    await user.click(screen.getByRole("button", { name: "Input Review へ戻る" }))

    expect(onReturnToInputReview).toHaveBeenCalledTimes(1)
    await waitFor(() => {
      expect(controller.mount).toHaveBeenCalledTimes(1)
    })
  })

  test("入力、runtime、credential の選択と validation/create action を controller へ委譲する", async () => {
    const user = userEvent.setup()
    const controller = new TranslationJobSetupScreenControllerFake(
      createViewModel({
        validationState: "not-run",
        dirty: false,
        validationResult: null,
        showCacheMissingGuidance: false,
        canValidate: true,
        blockedReasons: ["validation 未実行です。"],
        validationStatusLabel: "validation 未実行",
        validationStatusText: "validation 未実行です。入力、runtime、credential を確認して実行してください。"
      })
    )

    render(JobSetupPage, {
      props: {
        createController: () => controller
      }
    })

    await user.selectOptions(screen.getByLabelText("input data"), "41")
    await user.selectOptions(
      screen.getByLabelText("provider / model / execution mode"),
      "anthropic::claude-3-7-sonnet-with-a-very-long-name::sync"
    )
    await user.selectOptions(screen.getByLabelText("credential reference"), "cred-main")
    await user.click(screen.getByRole("button", { name: "validation を実行" }))

    expect(controller.selectInputSource).toHaveBeenCalledWith(41)
    expect(controller.selectRuntime).toHaveBeenCalledWith(
      "anthropic::claude-3-7-sonnet-with-a-very-long-name::sync"
    )
    expect(controller.selectCredentialRef).toHaveBeenCalledWith("cred-main")
    expect(controller.runValidation).toHaveBeenCalledTimes(1)
    expect(controller.createJob).not.toHaveBeenCalled()
  })

  test("create 成功後は read-only summary を表示し create action を隠す", () => {
    const controller = new TranslationJobSetupScreenControllerFake(
      createViewModel({
        summary: createSummary(),
        validationState: "fresh",
        dirty: false,
        validationResult: createValidationResult({
          status: "pass",
          blockingFailureCategory: undefined,
          targetSlices: ["input", "runtime", "credential"],
          canCreate: true,
          passSlices: ["input", "runtime", "credential"]
        }),
        canValidate: false,
        canCreate: false,
        showCacheMissingGuidance: false,
        blockedReasons: []
      })
    )

    render(JobSetupPage, {
      props: {
        createController: () => controller
      }
    })

    expect(screen.getByRole("heading", { level: 3, name: "Ready job summary" })).toBeInTheDocument()
    expect(screen.getByText("501")).toBeInTheDocument()
    expect(screen.getByText("ready")).toBeInTheDocument()
    expect(screen.getByText("/mods/very/long/path/translation/input-review-export.json")).toBeInTheDocument()
    expect(screen.getByText("openai-compatible")).toBeInTheDocument()
    expect(screen.getByText("gpt-4.1-mini-preview-with-a-very-long-name")).toBeInTheDocument()
    expect(screen.getByText("batch")).toBeInTheDocument()
    expect(screen.getAllByText("credential").length).toBeGreaterThan(0)
    expect(screen.queryByRole("button", { name: "ready job を作成" })).not.toBeInTheDocument()
    expect(screen.queryByRole("button", { name: "validation を実行" })).not.toBeInTheDocument()
  })
})