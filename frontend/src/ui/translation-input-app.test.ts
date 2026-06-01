import { render, screen, waitFor } from "@testing-library/svelte"
import userEvent from "@testing-library/user-event"
import { describe, expect, test, vi } from "vitest"

import type { MasterPersonaScreenControllerContract } from "@application/contract/master-persona/master-persona-screen-contract"
import type {
  PersonaGenerationPhaseScreenControllerContract,
  PersonaGenerationPhaseScreenViewModel,
  PersonaGenerationPhaseScreenViewModelListener
} from "@application/contract/persona-generation-phase"
import type {
  TermTranslationPhaseScreenControllerContract,
  TermTranslationPhaseScreenViewModel,
  TermTranslationPhaseScreenViewModelListener
} from "@application/contract/term-translation-phase"
import type {
  TranslationJobManagementScreenControllerContract,
  TranslationJobManagementScreenViewModelListener
} from "@application/contract/translation-job-management/translation-job-management-screen-contract"
import type { TranslationJobManagementScreenViewModel } from "@application/contract/translation-job-management/translation-job-management-screen-types"
import type { TranslationInputScreenControllerContract } from "@application/contract/translation-input"
import type { TranslationInputScreenViewModelListener } from "@application/contract/translation-input/translation-input-screen-contract"
import type {
  CreateTranslationJobFromInputResponse,
  TranslationInputReviewItem,
  TranslationInputScreenViewModel
} from "@application/gateway-contract/translation-input"
import App from "@ui/App.svelte"

function createItem(
  overrides: Partial<TranslationInputReviewItem> = {}
): TranslationInputReviewItem {
  return {
    localId: "input-41",
    inputId: 41,
    fileName: "kept-input.json",
    filePath: "/mods/kept-input.json",
    fileHash: "hash-41",
    importTimestamp: "2026-04-26T09:30:00Z",
    status: "registered",
    accepted: true,
    canRebuild: true,
    lastAction: "import",
    errorKind: null,
    warnings: [],
    summary: {
      input: {
        id: 41,
        sourceFilePath: "/mods/kept-input.json",
        sourceTool: "xEdit",
        targetPluginName: "Skyrim.esm",
        targetPluginType: "esm",
        recordCount: 3,
        importedAt: "2026-04-26T09:30:00Z"
      },
      translationRecordCount: 3,
      translationFieldCount: 4,
      categories: [
        {
          category: "NPC",
          recordCount: 3,
          fieldCount: 4
        }
      ],
      sampleFields: [
        {
          recordType: "NPC_",
          subrecordType: "FULL",
          formId: "00012345",
          editorId: "SampleNPC",
          sourceText: "Hello there",
          translatable: true
        }
      ],
      warnings: []
    },
    ...overrides
  }
}

function createViewModel(
  overrides: Partial<TranslationInputScreenViewModel> = {}
): TranslationInputScreenViewModel {
  const items = overrides.items ?? []
  const selectedItemId = overrides.selectedItemId ?? items[0]?.localId ?? null
  const selectedItem =
    overrides.selectedItem ??
    items.find((item) => item.localId === selectedItemId) ??
    null

  const baseViewModel: TranslationInputScreenViewModel = {
    items,
    selectedItemId,
    stagedFile: null,
    operationState: "idle",
    errorMessage: "",
    latestResponse: null,
    selectedItem,
    gatewayStatus: "接続準備済み",
    hasStagedFile: false,
    canImport: false,
    canRebuildSelected: selectedItem?.canRebuild ?? false,
    isImporting: false,
    isRebuilding: false,
    stagedFileName: "未選択",
    stagedFilePath: "-",
    stagedFileHash: "-",
    operationStatusLabel: "待機中",
    operationStatusText:
      "xEdit JSON を 1 件選び、登録結果と再構築状態をここで確認します。",
    latestOutcomeTitle: "登録結果はまだありません。",
    latestOutcomeText: "登録後に選択した入力データの概要をここへ表示します。",
    selectionStatusText: selectedItem
      ? `${selectedItem.fileName} / 登録済み`
      : "一覧から選択すると概要を右側へ表示します。",
    totalItemCountLabel: `${items.length} 件の input review を保持しています。`,
    emptyStateText:
      "まだ入力データがありません。JSON file を登録すると、一覧と sample field がここへ表示されます。"
  }

  return {
    ...baseViewModel,
    ...overrides,
    items,
    selectedItemId,
    selectedItem
  }
}

class TranslationInputScreenControllerFake implements TranslationInputScreenControllerContract {
  private readonly viewModel: TranslationInputScreenViewModel

  readonly mount = vi.fn(async () => {})
  readonly dispose = vi.fn(() => {})
  readonly selectItem = vi.fn(() => {})
  readonly stageJsonImport = vi.fn(async () => {})
  readonly resetImportSelection = vi.fn(() => {})
  readonly startImport = vi.fn(async () => {})
  readonly rebuildSelected = vi.fn(async () => {})
  readonly createTranslationJobFromSelected = vi.fn<
    () => Promise<CreateTranslationJobFromInputResponse | null>
  >(() => Promise.resolve(null))

  constructor(initialViewModel = createViewModel()) {
    this.viewModel = initialViewModel
  }

  subscribe(listener: TranslationInputScreenViewModelListener): () => void {
    void listener
    return () => {}
  }

  getViewModel(): TranslationInputScreenViewModel {
    return this.viewModel
  }
}

function createTranslationJobManagementViewModel(): TranslationJobManagementScreenViewModel {
  return {
    gatewayStatus: "接続済み",
    pageTitle: "未完了ジョブ一覧",
    pageLead:
      "新規翻訳を開始するか、未完了ジョブを選んで現在の翻訳段階から再開します。",
    headerCountLabel: "0 件を表示",
    listEmptyTitle: "管理対象がありません",
    listEmptyDescription: "未完了ジョブはありません。",
    listErrorTitle: "一覧を読み込めません",
    listErrorDescription: "未完了ジョブの一覧取得に失敗しました。",
    detailPlaceholderTitle: "job を選択してください",
    detailPlaceholderDescription: "一覧から 1 件選びます。",
    phase: "empty",
    detailPhase: "idle",
    isReloading: false,
    searchQuery: "",
    filterChips: [{ id: "all", label: "すべて", count: 0, selected: true }],
    jobs: [],
    feedback: null,
    selectedJob: null,
    deleteConfirmation: null,
    jobRunTarget: null
  }
}

class TranslationJobManagementScreenControllerFake implements TranslationJobManagementScreenControllerContract {
  private readonly viewModel = createTranslationJobManagementViewModel()

  readonly mount = vi.fn(async () => {})
  readonly dispose = vi.fn(() => {})
  readonly reload = vi.fn(async () => {})
  readonly selectJob = vi.fn(async () => {})
  readonly setFilter = vi.fn()
  readonly setSearchQuery = vi.fn()
  readonly requestStop = vi.fn(async () => {})
  readonly requestResume = vi.fn(async () => {})
  readonly openDeleteConfirmation = vi.fn()
  readonly closeDeleteConfirmation = vi.fn()
  readonly deleteSelectedJob = vi.fn(async () => {})

  subscribe(
    listener: TranslationJobManagementScreenViewModelListener
  ): () => void {
    void listener
    return () => {}
  }

  getViewModel(): TranslationJobManagementScreenViewModel {
    return this.viewModel
  }
}

class TermTranslationPhaseScreenControllerFake implements TermTranslationPhaseScreenControllerContract {
  readonly mount = vi.fn(async () => {})
  readonly dispose = vi.fn(() => {})
  readonly setJobId = vi.fn(async () => {})
  readonly refresh = vi.fn(async () => {})
  readonly startPhase = vi.fn(async () => {})
  readonly pausePhase = vi.fn(async () => {})
  readonly resumePhase = vi.fn(async () => {})
  readonly retryPhase = vi.fn(async () => {})

  subscribe(listener: TermTranslationPhaseScreenViewModelListener): () => void {
    void listener
    return () => {}
  }

  getViewModel(): TermTranslationPhaseScreenViewModel {
    return {
      jobId: 909,
      phase: "ready",
      summary: null,
      nextPhaseReadiness: null,
      errorMessage: "",
      pendingAction: null,
      hasLoaded: true,
      gatewayStatus: "接続準備済み",
      viewState: "idle_ready",
      hasJobSelection: true,
      statusTitle: "単語翻訳",
      statusText: "単語翻訳を開始できます。",
      phaseStateLabel: "開始待ち",
      progressPercent: 0,
      progressLabel: "0%",
      modelOptions: [],
      actionCards: []
    } as unknown as TermTranslationPhaseScreenViewModel
  }
}

class PersonaGenerationPhaseScreenControllerFake implements PersonaGenerationPhaseScreenControllerContract {
  readonly mount = vi.fn(async () => {})
  readonly dispose = vi.fn(() => {})
  readonly setJobId = vi.fn(async () => {})
  readonly refresh = vi.fn(async () => {})
  readonly startPhase = vi.fn(async () => {})
  readonly pausePhase = vi.fn(async () => {})
  readonly resumePhase = vi.fn(async () => {})
  readonly retryPhase = vi.fn(async () => {})
  readonly cancelPhase = vi.fn(async () => {})
  readonly checkBodyReadiness = vi.fn(async () => {})
  readonly startBodyPhase = vi.fn(async () => {})

  subscribe(
    listener: PersonaGenerationPhaseScreenViewModelListener
  ): () => void {
    void listener
    return () => {}
  }

  getViewModel(): PersonaGenerationPhaseScreenViewModel {
    return {
      jobId: 909,
      phase: "ready",
      summary: null,
      bodyReadiness: null,
      errorMessage: "",
      pendingAction: null,
      hasLoaded: true,
      gatewayStatus: "接続準備済み",
      viewState: "not_started",
      hasJobSelection: true,
      statusTitle: "NPC ペルソナ生成",
      statusText: "NPC ペルソナ生成は待機中です。",
      phaseStateLabel: "開始待ち",
      progressPercent: 0,
      progressLabel: "0%",
      modelOptions: [],
      actionCards: []
    } as unknown as PersonaGenerationPhaseScreenViewModel
  }
}

describe("App translation-input route", () => {
  test("translation-management route は未完了ジョブ一覧から新規開始でデータロードへ進む", async () => {
    window.history.replaceState(null, "", "#translation-management")

    const user = userEvent.setup()
    const controller = new TranslationInputScreenControllerFake()
    const jobManagementController =
      new TranslationJobManagementScreenControllerFake()
    const masterPersonaStub = (() =>
      ({}) as MasterPersonaScreenControllerContract) as () => MasterPersonaScreenControllerContract

    render(App, {
      props: {
        createMasterPersonaScreenController: masterPersonaStub,
        createTranslationJobManagementScreenController: () =>
          jobManagementController,
        createTranslationInputScreenController: () => controller
      }
    })

    expect(
      screen.getByRole("heading", { level: 1, name: "翻訳管理" })
    ).toBeInTheDocument()
    expect(
      screen.getByRole("heading", { level: 2, name: "未完了ジョブ一覧" })
    ).toBeInTheDocument()

    expect(screen.queryByRole("tab", { name: /データロード/ })).toBeNull()

    await user.click(screen.getByRole("button", { name: "新規翻訳を開始" }))

    expect(
      screen.getByRole("heading", { level: 2, name: "データロード" })
    ).toBeInTheDocument()
    expect(
      screen.queryByRole("heading", { level: 2, name: "Input Review" })
    ).not.toBeInTheDocument()

    await waitFor(() => {
      expect(controller.mount).toHaveBeenCalledTimes(1)
    })
  })

  test("route を離れて戻っても一覧と選択状態を維持する", async () => {
    window.history.replaceState(null, "", "#translation-management")

    const user = userEvent.setup()
    const selectedItem = createItem()
    const jobManagementController =
      new TranslationJobManagementScreenControllerFake()
    const controller = new TranslationInputScreenControllerFake(
      createViewModel({
        items: [selectedItem],
        selectedItemId: selectedItem.localId,
        selectedItem,
        canRebuildSelected: true,
        latestOutcomeTitle: "結果: 登録済み",
        latestOutcomeText:
          "翻訳レコード件数、カテゴリ別件数、sample field を確認できます。"
      })
    )
    const createTranslationInputScreenController = vi.fn(() => controller)
    const masterPersonaStub = (() =>
      ({}) as MasterPersonaScreenControllerContract) as () => MasterPersonaScreenControllerContract

    const { unmount } = render(App, {
      props: {
        createMasterPersonaScreenController: masterPersonaStub,
        createTranslationJobManagementScreenController: () =>
          jobManagementController,
        createTranslationInputScreenController
      }
    })

    expect(screen.queryByRole("tab", { name: /データロード/ })).toBeNull()

    await user.click(screen.getByRole("button", { name: "新規翻訳を開始" }))

    await waitFor(() => {
      expect(controller.mount).toHaveBeenCalledTimes(1)
    })

    expect(
      screen.getByRole("button", { name: /kept-input.json/ })
    ).toBeInTheDocument()
    expect(screen.getByText("登録結果")).toBeInTheDocument()
    expect(screen.getByText("問題区分")).toBeInTheDocument()

    await user.click(screen.getByRole("link", { name: "ダッシュボード" }))

    expect(
      screen.queryByRole("heading", { level: 2, name: "データロード" })
    ).not.toBeInTheDocument()
    expect(controller.dispose).toHaveBeenCalledTimes(1)

    await user.click(screen.getByRole("link", { name: "翻訳管理" }))
    expect(
      screen.getByRole("heading", { level: 2, name: "未完了ジョブ一覧" })
    ).toBeInTheDocument()
    expect(screen.queryByRole("tab", { name: /データロード/ })).toBeNull()
    await user.click(screen.getByRole("button", { name: "新規翻訳を開始" }))

    expect(
      screen.getByRole("heading", { level: 2, name: "データロード" })
    ).toBeInTheDocument()
    expect(
      screen.queryByRole("heading", { level: 2, name: "Input Review" })
    ).not.toBeInTheDocument()
    expect(
      screen.getByRole("button", { name: /kept-input.json/ })
    ).toBeInTheDocument()
    expect(screen.getByText("登録結果")).toBeInTheDocument()
    expect(screen.getByText("問題区分")).toBeInTheDocument()
    expect(createTranslationInputScreenController).toHaveBeenCalledTimes(2)
    expect(controller.mount).toHaveBeenCalledTimes(2)

    unmount()

    expect(controller.dispose).toHaveBeenCalledTimes(2)
  })

  test("入力データ確認から単語翻訳へ進むと作成済みジョブの単語翻訳画面を開く", async () => {
    window.history.replaceState(null, "", "#translation-management")

    const user = userEvent.setup()
    const selectedItem = createItem()
    const inputController = new TranslationInputScreenControllerFake(
      createViewModel({
        items: [selectedItem],
        selectedItemId: selectedItem.localId,
        selectedItem,
        canRebuildSelected: true,
        latestOutcomeTitle: "結果: 登録済み",
        latestOutcomeText: "入力データを確認済みです。"
      })
    )
    inputController.createTranslationJobFromSelected.mockResolvedValue({
      accepted: true,
      jobId: 909,
      jobState: "ready",
      currentPhase: "term_translation"
    })
    const jobManagementController =
      new TranslationJobManagementScreenControllerFake()
    const termController = new TermTranslationPhaseScreenControllerFake()
    const personaController = new PersonaGenerationPhaseScreenControllerFake()
    const masterPersonaStub = (() =>
      ({}) as MasterPersonaScreenControllerContract) as () => MasterPersonaScreenControllerContract

    render(App, {
      props: {
        createMasterPersonaScreenController: masterPersonaStub,
        createPersonaGenerationPhaseScreenController: () => personaController,
        createTermTranslationPhaseScreenController: () => termController,
        createTranslationJobManagementScreenController: () =>
          jobManagementController,
        createTranslationInputScreenController: () => inputController
      }
    })

    await user.click(screen.getByRole("button", { name: "新規翻訳を開始" }))
    await user.click(screen.getByRole("button", { name: "単語翻訳へ進む" }))

    await waitFor(() => {
      expect(
        inputController.createTranslationJobFromSelected
      ).toHaveBeenCalledTimes(1)
      expect(termController.setJobId).toHaveBeenCalledWith(909)
    })
    expect(
      screen.getByRole("heading", { level: 2, name: "単語翻訳" })
    ).toBeInTheDocument()
    expect(screen.getByText("ジョブ #909")).toBeInTheDocument()
  })
})
