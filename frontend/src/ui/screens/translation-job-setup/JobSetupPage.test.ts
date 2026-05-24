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
import type { TranslationJobSetupExtendedViewModel } from "@application/presenter/translation-job-setup/translation-job-setup.presenter"
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
    ...overrides
  }
}

function createPhaseOptions(
  overrides: Partial<TranslationJobSetupOptionsResponse> = {}
): TranslationJobSetupOptionsResponse {
  return createOptions({
    aiRuntimeOptions: [],
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
      },
      {
        provider: "openai",
        credentialRequirement: "required",
        supportedExecutionModes: ["sync"],
        supportsBatchMode: false
      }
    ],
    phaseRuntimeDrafts: [
      {
        phaseId: "word_translation",
        provider: "gemini",
        model: "gemini-model-core",
        credentialStatus: "configured",
        executionMode: "sync",
        batchMode: "enabled"
      },
      {
        phaseId: "npc_persona_generation",
        provider: "xai",
        model: "xai-persona-model",
        credentialStatus: "configured",
        executionMode: "sync",
        batchMode: "disabled"
      },
      {
        phaseId: "text_translation",
        provider: "lm_studio",
        model: "local-text-model",
        credentialStatus: "not_required",
        executionMode: "sync",
        batchMode: "unsupported"
      }
    ],
    ...overrides
  })
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
    deletingInputSourceId: null,
    selectedRuntimeKey: selectedRuntimeOption
      ? createTranslationJobSetupRuntimeKey(selectedRuntimeOption)
      : null,
    selectedCredentialRef: "",
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

function createPresentedPhaseViewModel(
  overrides: Partial<TranslationJobSetupScreenState> = {}
): TranslationJobSetupExtendedViewModel {
  const options = overrides.options ?? createPhaseOptions()
  const phaseRuntimeSelections =
    overrides.phaseRuntimeSelections ??
    options.phaseRuntimeDrafts?.map((draft) => ({ ...draft })) ??
    []
  const providerModelLists =
    overrides.providerModelLists ??
    phaseRuntimeSelections.map((selection) => ({
      phaseId: selection.phaseId,
      provider: selection.provider,
      credentialStatus: selection.credentialStatus,
      requestToken: "",
      sourceToken: "",
      status: "success",
      models: selection.model
        ? [{ modelId: selection.model, label: selection.model }]
        : []
    }))
  const presenter = new TranslationJobSetupPresenter()
  return presenter.toViewModel(
    {
      phase: "ready",
      options,
      selectedInputSourceId: options.inputCandidates[0]?.id ?? null,
      deletingInputSourceId: null,
      selectedRuntimeKey: null,
      selectedCredentialRef: "",
      phaseRuntimeSelections,
      providerModelLists,
      validationResult: createValidationResult({
        status: "pass",
        blockingFailureCategory: undefined,
        targetSlices: ["input", "phase-runtime"],
        validatedAt: "2026-05-04T02:30:00Z",
        canCreate: true,
        passSlices: ["input", "phase-runtime"]
      }),
      validationState: "fresh",
      dirty: false,
      errorMessage: "",
      createErrorKind: null,
      summary: null,
      ...overrides
    },
    true
  )
}

function createViewModel(
  overrides: Partial<TranslationJobSetupScreenViewModel> = {}
): TranslationJobSetupScreenViewModel {
  const options = overrides.options ?? createOptions()
  const validationResult =
    overrides.validationResult ?? createValidationResult()
  const summary = overrides.summary ?? null

  return {
    phase: summary ? "summary" : "ready",
    options,
    selectedInputSourceId: 41,
    deletingInputSourceId: null,
    selectedRuntimeKey:
      "openai-compatible::gpt-4.1-mini-preview-with-a-very-long-name::batch",
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
    availableCredentialRefs: [],
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

class TranslationJobSetupScreenControllerFake implements TranslationJobSetupScreenControllerContract {
  private viewModel: TranslationJobSetupScreenViewModel

  private readonly listeners =
    new Set<TranslationJobSetupScreenViewModelListener>()

  readonly mount = vi.fn(async () => {})
  readonly dispose = vi.fn(() => {})
  readonly selectInputSource = vi.fn(() => {})
  readonly deleteInputSource = vi.fn(async () => {})
  readonly selectRuntime = vi.fn(() => {})
  readonly selectCredentialRef = vi.fn(() => {})
  readonly runValidation = vi.fn(async () => {})
  readonly createJob = vi.fn(async () => {})
  readonly selectPhaseProvider = vi.fn(() => {})
  readonly refreshPhaseModels = vi.fn(async () => {})
  readonly selectPhaseModel = vi.fn(() => {})
  readonly togglePhaseBatchMode = vi.fn(() => {})
  readonly acknowledgeCredentialConfigured = vi.fn(() => {})

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
  test("SCN-TJSPPS-004: 入力データ領域は承認済み項目を表示する", () => {
    const controller = new TranslationJobSetupScreenControllerFake(
      createPresentedPhaseViewModel()
    )

    render(JobSetupPage, {
      props: {
        createController: () => controller
      }
    })

    expect(
      screen.getByRole("region", { name: "入力データ" })
    ).toBeInTheDocument()
    expect(screen.getAllByText("出自").length).toBeGreaterThan(0)
    expect(screen.getAllByText("翻訳レコード件数").length).toBeGreaterThan(0)
    expect(screen.getAllByText("登録日時").length).toBeGreaterThan(0)
    expect(screen.getAllByText("選択中").length).toBeGreaterThan(0)
    expect(screen.getByText("既存 job 状態")).toBeInTheDocument()
    expect(screen.getByText("既存 job はありません。")).toBeInTheDocument()
  })

  test("SCN-TJSPPS-005: ジョブ作成固定フッターは承認済み文言を表示する", () => {
    const controller = new TranslationJobSetupScreenControllerFake(
      createPresentedPhaseViewModel()
    )

    render(JobSetupPage, {
      props: {
        createController: () => controller
      }
    })

    expect(
      screen.getByRole("heading", { name: "ジョブの作成確認" })
    ).toBeInTheDocument()
    expect(
      screen.getByRole("region", {
        name: "ジョブの作成確認"
      })
    ).toBeInTheDocument()
    expect(
      screen.getByText("作成に必要な確認は完了しています。")
    ).toBeInTheDocument()
    expect(
      screen.getByRole("button", { name: "単語翻訳へ進む" })
    ).toBeEnabled()
    expect(
      screen.queryByRole("button", { name: "入力データの確認へ戻る" })
    ).not.toBeInTheDocument()
    expect(
      screen.queryByText("作成前に確認が必要な項目があります。")
    ).not.toBeInTheDocument()
  })

  test("SCN-PSJD-003: 作成不可時フッターは不足理由と戻り導線を表示する", async () => {
    const user = userEvent.setup()
    const onReturnToInputReview = vi.fn()
    const controller = new TranslationJobSetupScreenControllerFake(
      createViewModel({
        validationState: "fresh",
        canCreate: false,
        showCacheMissingGuidance: true,
        blockedReasons: [
          "入力データの確認が必要です。",
          "blocking failure を解消するまで create できません。"
        ]
      })
    )

    render(JobSetupPage, {
      props: {
        createController: () => controller,
        onReturnToInputReview
      }
    })

    expect(
      screen.getByRole("button", { name: "単語翻訳へ進む" })
    ).toBeDisabled()
    expect(
      screen.getByText("入力データの確認が必要です。入力データの確認画面に戻ってください。")
    ).toBeInTheDocument()
    expect(
      screen.getByRole("button", { name: "入力データの確認へ戻る" })
    ).toBeInTheDocument()

    await user.click(
      screen.getByRole("button", { name: "入力データの確認へ戻る" })
    )

    expect(onReturnToInputReview).toHaveBeenCalledTimes(1)
    expect(controller.createJob).not.toHaveBeenCalled()
  })

  test("SCN-PSJD-002: Job Setup の Ready job 作成後要約は選択値だけを表示する", async () => {
    const user = userEvent.setup()
    const controller = new TranslationJobSetupScreenControllerFake(
      createPresentedPhaseViewModel()
    )

    render(JobSetupPage, {
      props: {
        createController: () => controller
      }
    })

    await user.click(screen.getByRole("button", { name: "単語翻訳へ進む" }))

    controller.pushViewModel(
      createPresentedPhaseViewModel({
        summary: createSummary({
          phaseRuntimeSummaries: [
            {
              phaseId: "word_translation",
              provider: "gemini",
              model: "gemini-ready-model",
              credentialStatus: "configured",
              executionMode: "sync",
              batchMode: "enabled"
            },
            {
              phaseId: "npc_persona_generation",
              provider: "xai",
              model: "xai-ready-model",
              credentialStatus: "configured",
              executionMode: "sync",
              batchMode: "disabled"
            },
            {
              phaseId: "text_translation",
              provider: "lm_studio",
              model: "local-ready-model",
              credentialStatus: "not_required",
              executionMode: "sync",
              batchMode: "unsupported"
            }
          ]
        })
      })
    )

    expect(
      await screen.findByRole("heading", { name: "翻訳段階ごとの設定" })
    ).toBeInTheDocument()
    expect(screen.getByText("gemini-ready-model")).toBeInTheDocument()
    expect(screen.getByText("xai-ready-model")).toBeInTheDocument()
    expect(screen.getByText("local-ready-model")).toBeInTheDocument()
    expect(screen.getAllByText("設定済み").length).toBeGreaterThan(0)
    expect(screen.getByText("APIキー不要")).toBeInTheDocument()
    expect(document.body).not.toHaveTextContent(
      "credential-ref-must-stay-hidden"
    )
    expect(document.body).not.toHaveTextContent(
      "xai-secret-ref-must-stay-hidden"
    )
    expect(document.body).not.toHaveTextContent(
      "model-list-token-must-stay-hidden"
    )
    expect(document.body).not.toHaveTextContent("endpoint")
    expect(document.body).not.toHaveTextContent("raw request")
  })

  test("input metadata の registeredAt supplied 値を表示する", () => {
    const registeredAt = "2026-04-27T00:30:00.000Z"
    const controller = new TranslationJobSetupScreenControllerFake(
      createPresentedViewModel({
        options: createOptions({
          inputCandidates: [
            {
              id: 41,
              label:
                "/mods/very/long/path/translation/input-review-export.json",
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

    expect(
      screen.getAllByText(new Date(registeredAt).toLocaleString("ja-JP")).length
    ).toBeGreaterThan(0)
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

    expect(screen.getByRole("heading", { name: "入力データ" })).toBeInTheDocument()
    expect(
      screen.getAllByText(
        "/mods/very/long/path/translation/input-review-export.json"
      ).length
    ).toBeGreaterThan(0)
    expect(screen.getAllByText("xEdit extract").length).toBeGreaterThan(0)
    expect(screen.getAllByText("2026/4/27 9:30:00").length).toBeGreaterThan(0)
    expect(screen.getAllByText("128 件").length).toBeGreaterThan(0)
    expect(screen.getByText("既存 job はありません。")).toBeInTheDocument()
    expect(screen.getByText("ジョブの作成確認")).toBeInTheDocument()
    expect(
      screen.getByText("入力データの確認が必要です。入力データの確認画面に戻ってください。")
    ).toBeInTheDocument()
    expect(
      screen.getByRole("button", { name: "単語翻訳へ進む" })
    ).toBeDisabled()
    expect(
      screen.getByText("validation が失効しています。")
    ).toBeInTheDocument()
    expect(
      screen.getByText("blocking failure を解消するまで create できません。")
    ).toBeInTheDocument()
    expect(screen.getByText("validation が失効しています。")).toBeInTheDocument()
    expect(
      screen.queryByRole("heading", { name: "共通辞書と共通ペルソナ" })
    ).not.toBeInTheDocument()
    expect(
      screen.queryByRole("heading", { name: "作成前確認" })
    ).not.toBeInTheDocument()

    await user.click(
      screen.getByRole("button", { name: "入力データの確認へ戻る" })
    )

    expect(onReturnToInputReview).toHaveBeenCalledTimes(1)
    await waitFor(() => {
      expect(controller.mount).toHaveBeenCalledTimes(1)
    })
  })

  test("入力カード選択、削除、入力候補確認の操作を controller へ委譲する", async () => {
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
        validationStatusText:
          "validation 未実行です。入力、runtime、credential を確認して実行してください。"
      })
    )

    render(JobSetupPage, {
      props: {
        createController: () => controller
      }
    })

    await user.click(
      screen.getByRole("button", {
        name: /input-review-export\.json/
      })
    )
    await user.click(screen.getByRole("button", { name: "削除" }))
    await user.click(screen.getByRole("button", { name: "選択候補を確認" }))

    expect(controller.selectInputSource).toHaveBeenCalledWith(41)
    expect(controller.deleteInputSource).toHaveBeenCalledWith(41)
    expect(controller.selectInputSource).toHaveBeenCalledTimes(2)
    expect(controller.selectRuntime).not.toHaveBeenCalled()
    expect(controller.selectCredentialRef).not.toHaveBeenCalled()
    expect(controller.runValidation).not.toHaveBeenCalled()
    expect(controller.createJob).not.toHaveBeenCalled()
  })

  test("削除中カードだけ 削除中 表示と選択不可を出し、他カード削除を止める", () => {
    const controller = new TranslationJobSetupScreenControllerFake(
      createViewModel({
        options: createOptions({
          inputCandidates: [
            {
              id: 41,
              label: "/mods/input-a.json",
              sourceKind: "xEdit extract",
              recordCount: 128
            },
            {
              id: 52,
              label: "/mods/input-b.json",
              sourceKind: "xEdit extract",
              recordCount: 64
            }
          ]
        }),
        selectedInputSourceId: 41,
        deletingInputSourceId: 41
      })
    )

    render(JobSetupPage, {
      props: {
        createController: () => controller
      }
    })

    const deleteButtons = screen.getAllByRole("button")
    expect(screen.getAllByText("削除中...")).toHaveLength(2)
    expect(
      screen.getByRole("button", { name: /\/mods\/input-a\.json/ })
    ).toBeDisabled()
    expect(
      screen.getByRole("button", { name: /\/mods\/input-b\.json/ })
    ).not.toBeDisabled()
    expect(
      deleteButtons.filter((button) => button.textContent === "削除中...")
    ).toHaveLength(1)
    expect(screen.getByRole("button", { name: "削除中..." })).toBeDisabled()
    expect(screen.getByRole("button", { name: "削除" })).toBeDisabled()
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

    expect(
      screen.getByRole("heading", { level: 3, name: "Ready job summary" })
    ).toBeInTheDocument()
    expect(screen.getByText("501")).toBeInTheDocument()
    expect(screen.getByText("ready")).toBeInTheDocument()
    expect(
      screen.getByText(
        "/mods/very/long/path/translation/input-review-export.json"
      )
    ).toBeInTheDocument()
    expect(screen.getByText("openai-compatible")).toBeInTheDocument()
    expect(
      screen.getByText("gpt-4.1-mini-preview-with-a-very-long-name")
    ).toBeInTheDocument()
    expect(screen.getByText("batch")).toBeInTheDocument()
    expect(screen.getAllByText("credential").length).toBeGreaterThan(0)
    expect(
      screen.queryByRole("button", { name: "ready job を作成" })
    ).not.toBeInTheDocument()
    expect(
      screen.queryByRole("button", { name: "validation を実行" })
    ).not.toBeInTheDocument()
  })
})
