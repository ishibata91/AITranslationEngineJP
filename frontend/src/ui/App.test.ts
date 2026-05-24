import { render, screen, waitFor, within } from "@testing-library/svelte"
import userEvent from "@testing-library/user-event"

import type {
  MasterDictionaryScreenControllerContract,
  MasterDictionaryScreenViewModelListener
} from "@application/contract/master-dictionary/master-dictionary-screen-contract"
import type { MasterDictionaryScreenViewModel } from "@application/contract/master-dictionary/master-dictionary-screen-types"
import type {
  MasterPersonaScreenControllerContract,
  MasterPersonaScreenViewModelListener
} from "@application/contract/master-persona/master-persona-screen-contract"
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
import type {
  MasterPersonaScreenViewModel,
  MasterPersonaDetail
} from "@application/gateway-contract/master-persona"
import type { MasterDictionaryEntryDetail } from "@application/gateway-contract/master-dictionary"
import App from "@ui/App.svelte"
import { vi } from "vitest"

const DASHBOARD_SHELL_PRIMARY_ROUTES = [
  {
    id: "dashboard",
    label: "ダッシュボード",
    state: "既定表示",
    description: "主要ページへの入口をまとめて確認します。"
  },
  {
    id: "provider-settings",
    label: "AIサービス設定",
    state: "設定状態を確認",
    description: "AIサービスごとの接続設定を確認します。"
  },
  {
    id: "master-dictionary",
    label: "マスター辞書",
    state: "準備中",
    description: "用語と訳語の基盤データを確認します。"
  },
  {
    id: "master-persona",
    label: "マスターペルソナ",
    state: "準備中",
    description: "翻訳に使うペルソナ設定を確認します。"
  },
  {
    id: "translation-management",
    label: "翻訳管理",
    state: "準備中",
    description: "準備状況と翻訳ジョブの進行をまとめて確認します。"
  },
  {
    id: "output-management",
    label: "出力管理",
    state: "準備中",
    description: "生成物と書き出し結果を確認します。"
  }
] as const

const DASHBOARD_ENTRY_ROUTES = DASHBOARD_SHELL_PRIMARY_ROUTES.filter(
  ({ id }) => id !== "dashboard"
)

const PLACEHOLDER_LEAD =
  "このページはまだ準備中です。上のナビゲーションまたは下の移動から別の主要ページへ進めます。"

function createDefaultSelectedEntry(): MasterDictionaryEntryDetail {
  return {
    id: "101",
    source: "Dragon Priest",
    translation: "ドラゴン・プリースト",
    category: "固有名詞",
    origin: "初期データ",
    note: "マスター辞書エントリ",
    updatedAt: "2026-01-01 00:00"
  }
}

function createImportedBookSelectedEntry(): MasterDictionaryEntryDetail {
  return {
    id: "201",
    source: "Allowed Book Source",
    translation: "許可された本の訳語",
    category: "書籍",
    origin: "XML取込",
    note: "マスター辞書エントリ",
    updatedAt: "2026-04-12 00:00"
  }
}

function buildMasterDictionaryScreenViewModel(
  overrides: Partial<MasterDictionaryScreenViewModel> = {}
): MasterDictionaryScreenViewModel {
  const selectedEntry =
    overrides.selectedEntry === undefined
      ? createDefaultSelectedEntry()
      : overrides.selectedEntry

  return {
    entries: [
      {
        id: "101",
        source: "Dragon Priest",
        translation: "ドラゴン・プリースト",
        category: "固有名詞",
        origin: "初期データ",
        updatedAt: "2026-01-01 00:00"
      },
      {
        id: "102",
        source: "The Reach",
        translation: "リーチ地方",
        category: "地名",
        origin: "初期データ",
        updatedAt: "2026-01-02 00:00"
      }
    ],
    selectedEntry,
    selectedId: selectedEntry?.id ?? null,
    totalCount: 2,
    query: "",
    category: "すべて",
    page: 0,
    errorMessage: "",
    modalState: null,
    formSource: "",
    formCategory: "固有名詞",
    formOrigin: "手動登録",
    formTranslation: "",
    selectedFileName: "未選択",
    selectedFileReference: null,
    importStage: "idle",
    importProgress: 0,
    importSummary: null,
    gatewayStatus: "接続済み",
    hasStagedFile: false,
    isImportRunning: false,
    importStatusValue: "待機中",
    importStatusText: "ファイルを選ぶと取込バーが表示されます。",
    categoryOptions: ["すべて", "固有名詞", "地名", "書籍", "装備"],
    totalPages: 1,
    pageStatusText: "1 - 2 件を表示",
    listHeadline: "2 件のエントリを表示しています。",
    selectionStatusText: selectedEntry
      ? `${selectedEntry.source} を選択中`
      : "一致するエントリがありません",
    detailSublineText: selectedEntry
      ? `${selectedEntry.origin} / 最終更新 ${selectedEntry.updatedAt}`
      : "一覧から別のエントリを選択すると、ここも切り替わります。",
    ...overrides
  }
}

function buildMasterPersonaScreenViewModel(
  overrides: Partial<MasterPersonaScreenViewModel> = {}
): MasterPersonaScreenViewModel {
  return {
    items: [
      {
        identityKey: "TestPersonaPluginA.esp:FE01A812:NPC_",
        targetPlugin: "TestPersonaPluginA.esp",
        formId: "FE01A812",
        recordType: "NPC_",
        editorId: "TEST_NPC_A",
        displayName: "Test NPC A",
        voiceType: "TestVoiceA",
        className: "TestClassA",
        sourcePlugin: "TestPersonaPluginA.esp",
        personaSummary: "乾いた率直さで応じる。",
        updatedAt: "2026-04-15T09:42:00Z"
      }
    ],
    pluginGroups: [{ targetPlugin: "TestPersonaPluginA.esp", count: 1 }],
    selectedIdentityKey: "TestPersonaPluginA.esp:FE01A812:NPC_",
    selectedEntry: {
      identityKey: "TestPersonaPluginA.esp:FE01A812:NPC_",
      targetPlugin: "TestPersonaPluginA.esp",
      formId: "FE01A812",
      recordType: "NPC_",
      editorId: "TEST_NPC_A",
      displayName: "Test NPC A",
      voiceType: "TestVoiceA",
      className: "TestClassA",
      sourcePlugin: "TestPersonaPluginA.esp",
      personaSummary: "乾いた率直さで応じる。",
      updatedAt: "2026-04-15T09:42:00Z",
      personaBody: "短く本音を置く。",
      runLockReason: "更新と削除を行えます"
    } as MasterPersonaDetail,
    keyword: "",
    pluginFilter: "",
    page: 1,
    pageSize: 30,
    totalCount: 1,
    errorMessage: "",
    aiSettings: {
      provider: "gemini",
      model: "gemini-2.5-pro",
      executionMethod: "single_request"
    },
    aiSettingsMessage: "",
    modelOptions: [],
    selectedFileName: "未選択",
    selectedFileReference: null,
    preview: null,
    runStatus: {
      runState: "入力待ち",
      targetPlugin: "",
      processedCount: 0,
      successCount: 0,
      existingSkipCount: 0,
      zeroDialogueSkipCount: 0,
      genericNpcCount: 0,
      currentActorLabel: "",
      message: "入力ファイルを選ぶと状態を表示します。"
    },
    modalState: null,
    editForm: {
      formId: "FE01A812",
      editorId: "TEST_NPC_A",
      displayName: "Test NPC A",
      voiceType: "TestVoiceA",
      className: "TestClassA",
      sourcePlugin: "TestPersonaPluginA.esp",
      personaBody: "短く本音を置く。"
    },
    gatewayStatus: "接続準備済み",
    pluginOptions: [
      { value: "", label: "すべてのプラグイン" },
      { value: "TestPersonaPluginA.esp", label: "TestPersonaPluginA.esp (1)" }
    ],
    totalPages: 1,
    pageStatusText: "1 - 1 件を表示しています。",
    selectionStatusText: "Test NPC A を選択中",
    listHeadline: "1 件から絞り込みます。",
    detailLockText: "更新と削除を行えます",
    detailStatusText: "更新と削除を行えます",
    canStartPreview: false,
    canStartGeneration: false,
    canMutate: true,
    isRunActive: false,
    hasPreview: false,
    aiProviderLabel: "Gemini",
    aiSettingsWarningText: "",
    aiSettingsStatusText: "設定済み",
    canSelectModel: true,
    executionMethodOptions: [{ value: "single_request", label: "通常" }],
    promptTemplateDescription:
      "プロンプトテンプレートは画面入力では変更せず、実装側の説明文として固定しています。",
    progressPercent: 0,
    ...overrides
  } as MasterPersonaScreenViewModel
}

function buildTermTranslationPhaseScreenViewModel(
  overrides: Partial<TermTranslationPhaseScreenViewModel> = {}
): TermTranslationPhaseScreenViewModel {
  return {
    jobId: 101,
    phase: "ready",
    summary: null,
    nextPhaseReadiness: null,
    errorMessage: "",
    pendingAction: null,
    hasLoaded: true,
    gatewayStatus: "接続準備済み",
    viewState: "idle_ready",
    isLoading: false,
    isRefreshing: false,
    isSubmitting: false,
    hasJobSelection: true,
    currentPhaseLabel: "term_translation",
    phaseStateLabel: "開始待ち",
    statusTitle: "Term Translation Phase",
    statusText: "summary を取得できます。",
    progressPercent: 0,
    progressLabel: "0%",
    progressDetail: "対象なし",
    startedAtLabel: "-",
    finishedAtLabel: "-",
    totalTermCountLabel: "-",
    dictionaryHitCountLabel: "-",
    aiTargetCountLabel: "-",
    confirmedCountLabel: "-",
    jobDictionaryAppliedCountLabel: "-",
    replacementTargetCountLabel: "-",
    unmatchedCountLabel: "-",
    providerLabel: "-",
    modelLabel: "-",
    executionModeLabel: "-",
    credentialRefLabel: "-",
    snapshotLabel: "-",
    errorKindLabel: "-",
    errorReasonLabel: "-",
    retryableLabel: "-",
    nextPhaseStatusLabel: "開始不可",
    nextPhaseBlockedReason: "",
    providerSkippedLabel: "-",
    actionCards: [],
    lastErrorSummary: null,
    actionEnablement: null,
    latestProgressSummary: null,
    latestResultSummary: null,
    latestExecutionSummary: null,
    latestErrorKind: null,
    ...overrides
  }
}

function buildPersonaGenerationPhaseScreenViewModel(
  overrides: Partial<PersonaGenerationPhaseScreenViewModel> = {}
): PersonaGenerationPhaseScreenViewModel {
  return {
    jobId: 101,
    phase: "ready",
    summary: null,
    bodyReadiness: null,
    errorMessage: "",
    pendingAction: null,
    hasLoaded: true,
    gatewayStatus: "接続準備済み",
    viewState: "not_started",
    isLoading: false,
    isRefreshing: false,
    isSubmitting: false,
    hasJobSelection: true,
    currentPhaseLabel: "persona_generation",
    phaseStateLabel: "開始待ち",
    statusTitle: "Persona Generation Phase",
    statusText: "summary を取得できます。",
    progressPercent: 0,
    progressLabel: "0%",
    progressDetail: "対象なし",
    startedAtLabel: "-",
    finishedAtLabel: "-",
    targetCountLabel: "-",
    generatedCountLabel: "-",
    failedCountLabel: "-",
    skippedCountLabel: "-",
    npcCountLabel: "-",
    commonPersonaHitCountLabel: "-",
    commonPersonaMissCountLabel: "-",
    skippedReasonsLabel: "-",
    targetSnapshotLabel: "-",
    providerLabel: "-",
    modelLabel: "-",
    executionModeLabel: "-",
    credentialRefLabel: "-",
    inputCountLabel: "-",
    outputCountLabel: "-",
    evidenceRefsLabel: "-",
    promptDigestLabel: "-",
    snapshotLabel: "-",
    snapshotReferenceStatusLabel: "-",
    personaCountLabel: "-",
    missingCountLabel: "-",
    bodyReadinessLabel: "未確認",
    bodyReadinessBlockedReason: "",
    bodyReadinessInputSummaryLabel: "-",
    errorKindLabel: "-",
    errorReasonLabel: "-",
    retryableLabel: "-",
    actionCards: [],
    screenActionEnablement: {
      canRefresh: true,
      canStart: false,
      canPause: false,
      canResume: false,
      canRetry: false,
      canCancel: false,
      canCheckBodyReadiness: false,
      canStartBodyPhase: false
    },
    lastErrorSummary: null,
    actionEnablement: null,
    latestProgressSummary: null,
    latestTargetSummary: null,
    latestResultSummary: null,
    latestExecutionSummary: null,
    latestErrorKind: null,
    latestBodyReadiness: null,
    latestBodyReadinessInputSummary: null,
    ...overrides
  }
}

class MasterDictionaryScreenControllerFake implements MasterDictionaryScreenControllerContract {
  private viewModel: MasterDictionaryScreenViewModel

  private readonly listeners =
    new Set<MasterDictionaryScreenViewModelListener>()

  readonly mount = vi.fn(async () => {})

  readonly dispose = vi.fn(() => {})

  readonly selectRow = vi.fn(async () => {})

  readonly openCreateModal = vi.fn(() => {})

  readonly openEditModal = vi.fn(() => {})

  readonly openDeleteModal = vi.fn(() => {})

  readonly closeEditModal = vi.fn(() => {})

  readonly closeDeleteModal = vi.fn(() => {})

  readonly saveCurrentEntry = vi.fn(async () => {})

  readonly deleteCurrentEntry = vi.fn(async () => {})

  readonly handleSearchInput = vi.fn(() => {})

  readonly handleCategoryChange = vi.fn(() => {})

  readonly goToPrevPage = vi.fn(() => {})

  readonly goToNextPage = vi.fn(() => {})

  readonly stageXmlImport = vi.fn(() => {})

  readonly resetImportSelection = vi.fn(() => {})

  readonly startImport = vi.fn(async () => {})

  readonly setFormSource = vi.fn(() => {})

  readonly setFormCategory = vi.fn(() => {})

  readonly setFormOrigin = vi.fn(() => {})

  readonly setFormTranslation = vi.fn(() => {})

  constructor(initialViewModel = buildMasterDictionaryScreenViewModel()) {
    this.viewModel = initialViewModel
  }

  subscribe(listener: MasterDictionaryScreenViewModelListener): () => void {
    this.listeners.add(listener)
    return () => {
      this.listeners.delete(listener)
    }
  }

  getViewModel(): MasterDictionaryScreenViewModel {
    return this.viewModel
  }

  pushViewModel(nextViewModel: MasterDictionaryScreenViewModel): void {
    this.viewModel = nextViewModel
    for (const listener of this.listeners) {
      listener(nextViewModel)
    }
  }
}

class MasterPersonaScreenControllerFake implements MasterPersonaScreenControllerContract {
  private viewModel: MasterPersonaScreenViewModel

  private readonly listeners = new Set<MasterPersonaScreenViewModelListener>()

  readonly mount = vi.fn(async () => {})
  readonly dispose = vi.fn(() => {})
  readonly selectRow = vi.fn(async () => {})
  readonly handleSearchInput = vi.fn(() => {})
  readonly handlePluginFilterChange = vi.fn(() => {})
  readonly goToPrevPage = vi.fn(() => {})
  readonly goToNextPage = vi.fn(() => {})
  readonly stageJsonSelection = vi.fn(() => {})
  readonly resetJsonSelection = vi.fn(() => {})
  readonly previewGeneration = vi.fn(async () => {})
  readonly executeGeneration = vi.fn(async () => {})
  readonly interruptGeneration = vi.fn(async () => {})
  readonly cancelGeneration = vi.fn(async () => {})
  readonly saveAISettings = vi.fn(async () => {})
  readonly setAIProvider = vi.fn(() => {})
  readonly setAIModel = vi.fn(() => {})
  readonly setAIExecutionMethod = vi.fn(() => {})
  readonly openDialogueModal = vi.fn(async () => {})
  readonly closeDialogueModal = vi.fn(() => {})
  readonly openEditModal = vi.fn(() => {})
  readonly closeEditModal = vi.fn(() => {})
  readonly openDeleteModal = vi.fn(() => {})
  readonly closeDeleteModal = vi.fn(() => {})
  readonly saveCurrentEntry = vi.fn(async () => {})
  readonly deleteCurrentEntry = vi.fn(async () => {})
  readonly setEditFormField = vi.fn(() => {})

  constructor(initialViewModel = buildMasterPersonaScreenViewModel()) {
    this.viewModel = initialViewModel
  }

  subscribe(listener: MasterPersonaScreenViewModelListener): () => void {
    this.listeners.add(listener)
    return () => {
      this.listeners.delete(listener)
    }
  }

  getViewModel(): MasterPersonaScreenViewModel {
    return this.viewModel
  }

  pushViewModel(nextViewModel: MasterPersonaScreenViewModel): void {
    this.viewModel = nextViewModel
    for (const listener of this.listeners) {
      listener(nextViewModel)
    }
  }
}

class TermTranslationPhaseScreenControllerFake implements TermTranslationPhaseScreenControllerContract {
  private viewModel: TermTranslationPhaseScreenViewModel

  private readonly listeners =
    new Set<TermTranslationPhaseScreenViewModelListener>()

  readonly mount = vi.fn(async () => {})
  readonly dispose = vi.fn(() => {})
  readonly setJobId = vi.fn(async () => {})
  readonly refresh = vi.fn(async () => {})
  readonly startPhase = vi.fn(async () => {})
  readonly pausePhase = vi.fn(async () => {})
  readonly resumePhase = vi.fn(async () => {})
  readonly retryPhase = vi.fn(async () => {})

  constructor(initialViewModel = buildTermTranslationPhaseScreenViewModel()) {
    this.viewModel = initialViewModel
  }

  subscribe(listener: TermTranslationPhaseScreenViewModelListener): () => void {
    this.listeners.add(listener)
    return () => {
      this.listeners.delete(listener)
    }
  }

  getViewModel(): TermTranslationPhaseScreenViewModel {
    return this.viewModel
  }
}

class PersonaGenerationPhaseScreenControllerFake implements PersonaGenerationPhaseScreenControllerContract {
  private viewModel: PersonaGenerationPhaseScreenViewModel

  private readonly listeners =
    new Set<PersonaGenerationPhaseScreenViewModelListener>()

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

  constructor(initialViewModel = buildPersonaGenerationPhaseScreenViewModel()) {
    this.viewModel = initialViewModel
  }

  subscribe(
    listener: PersonaGenerationPhaseScreenViewModelListener
  ): () => void {
    this.listeners.add(listener)
    return () => {
      this.listeners.delete(listener)
    }
  }

  getViewModel(): PersonaGenerationPhaseScreenViewModel {
    return this.viewModel
  }
}

function buildTranslationJobManagementScreenViewModel(): TranslationJobManagementScreenViewModel {
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
  private readonly viewModel = buildTranslationJobManagementScreenViewModel()

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

function renderApp(
  controller = new MasterDictionaryScreenControllerFake(),
  masterPersonaController = new MasterPersonaScreenControllerFake()
): MasterDictionaryScreenControllerFake {
  render(App, {
    props: {
      createMasterDictionaryScreenController: () => controller,
      createMasterPersonaScreenController: () => masterPersonaController
    }
  })

  return controller
}

function renderAppView(
  controller = new MasterDictionaryScreenControllerFake(),
  masterPersonaController = new MasterPersonaScreenControllerFake()
): {
  controller: MasterDictionaryScreenControllerFake
  masterPersonaController: MasterPersonaScreenControllerFake
  unmount: () => void
} {
  const view = render(App, {
    props: {
      createMasterDictionaryScreenController: () => controller,
      createMasterPersonaScreenController: () => masterPersonaController
    }
  })

  return {
    controller,
    masterPersonaController,
    unmount: () => view.unmount()
  }
}

function createXmlFile(contents: string, name = "master-dictionary.xml"): File {
  return new File([contents], name, { type: "text/xml" })
}

function getGlobalNavigation(): HTMLElement {
  return screen.getByRole("navigation", {
    name: "グローバルナビゲーション"
  })
}

function getDashboardCardSection(): HTMLElement {
  const dashboardCardSection = screen
    .getByRole("heading", { name: "作業を選ぶ" })
    .closest("section")

  if (!dashboardCardSection) {
    throw new Error("ダッシュボード入口カードのセクションが見つかりません")
  }

  return dashboardCardSection
}

function getDashboardCardLink(routeLabel: string): HTMLAnchorElement {
  const cardHeading = within(getDashboardCardSection()).getByRole("heading", {
    level: 3,
    name: routeLabel
  })
  const cardLink = cardHeading.closest("a")

  if (!(cardLink instanceof HTMLAnchorElement)) {
    throw new Error(`入口カードのリンクが見つかりません: ${routeLabel}`)
  }

  return cardLink
}

describe("App dashboard shell", () => {
  beforeEach(() => {
    window.history.replaceState(null, "", "#")
  })

  test("SCN-DAS-001: 起動時にダッシュボード見出しを表示する", () => {
    // Arrange
    renderApp()

    // Act
    const dashboardHeading = screen.getByRole("heading", {
      name: "ダッシュボード"
    })

    // Assert
    expect(dashboardHeading).toBeInTheDocument()
  })

  test("SCN-DAS-001: 起動時に hash を dashboard へ正規化する", () => {
    // Arrange
    renderApp()

    // Act
    const hash = window.location.hash

    // Assert
    expect(hash).toBe("#dashboard")
  })

  test("invalid hash はダッシュボード見出しを表示する", () => {
    // Arrange
    window.history.replaceState(null, "", "#not-approved-route")

    // Act
    renderApp()

    // Assert
    expect(
      screen.getByRole("heading", { name: "ダッシュボード" })
    ).toBeInTheDocument()
  })

  test("invalid hash は dashboard hash へ正規化する", () => {
    // Arrange
    window.history.replaceState(null, "", "#not-approved-route")
    renderApp()

    // Act
    const hash = window.location.hash

    // Assert
    expect(hash).toBe("#dashboard")
  })

  test("SCN-DAS-002: グローバルナビゲーションに承認済みルートを並べる", () => {
    // Arrange
    renderApp()

    // Act
    const links = within(getGlobalNavigation()).getAllByRole("link")

    // Assert
    expect(links).toHaveLength(DASHBOARD_SHELL_PRIMARY_ROUTES.length)
  })

  test.each(DASHBOARD_ENTRY_ROUTES)(
    "SCN-DAS-002: グローバルナビゲーションから $label 見出しへ移動できる",
    async ({ label }) => {
      // Arrange
      const user = userEvent.setup()
      renderApp()

      // Act
      await user.click(
        within(getGlobalNavigation()).getByRole("link", { name: label })
      )

      // Assert
      expect(
        screen.getByRole("heading", { level: 1, name: label })
      ).toBeInTheDocument()
    }
  )

  test.each(DASHBOARD_ENTRY_ROUTES)(
    "SCN-DAS-002: グローバルナビゲーションから $id hash へ移動できる",
    async ({ id, label }) => {
      // Arrange
      const user = userEvent.setup()
      renderApp()

      // Act
      await user.click(
        within(getGlobalNavigation()).getByRole("link", { name: label })
      )

      // Assert
      expect(window.location.hash).toBe(`#${id}`)
    }
  )

  test.each(DASHBOARD_ENTRY_ROUTES)(
    "SCN-DAS-003: 入口カードから $label 見出しへ移動できる",
    async ({ label }) => {
      // Arrange
      const user = userEvent.setup()
      renderApp()

      // Act
      await user.click(getDashboardCardLink(label))

      // Assert
      expect(
        screen.getByRole("heading", { level: 1, name: label })
      ).toBeInTheDocument()
    }
  )

  test.each(DASHBOARD_ENTRY_ROUTES)(
    "SCN-DAS-003: 入口カードから $id hash へ移動できる",
    async ({ id, label }) => {
      // Arrange
      const user = userEvent.setup()
      renderApp()

      // Act
      await user.click(getDashboardCardLink(label))

      // Assert
      expect(window.location.hash).toBe(`#${id}`)
    }
  )

  test("SCN-DAS-003: 入口カードに確認可能を表示し準備中を表示しない", () => {
    // Arrange
    renderApp()

    // Act
    const dashboardCardSection = getDashboardCardSection()
    const dashboardCardQuery = within(dashboardCardSection)

    // Assert
    expect(dashboardCardQuery.getAllByText("確認可能")).toHaveLength(3)
    expect(dashboardCardQuery.queryByText("準備中")).not.toBeInTheDocument()
  })

  test("SCN-DAS-004: プレースホルダー画面でも共通 lead を表示する", async () => {
    // Arrange
    const user = userEvent.setup()
    renderApp()

    // Act
    await user.click(
      within(getGlobalNavigation()).getByRole("link", { name: "出力管理" })
    )

    // Assert
    expect(screen.getByText(PLACEHOLDER_LEAD)).toBeInTheDocument()
  })

  test("SCN-DAS-005: プレースホルダー画面から別の主要ページへ再移動できる", async () => {
    // Arrange
    const user = userEvent.setup()
    renderApp()

    await user.click(
      within(getGlobalNavigation()).getByRole("link", { name: "翻訳管理" })
    )

    // Act
    await user.click(
      within(getGlobalNavigation()).getByRole("link", { name: "出力管理" })
    )

    // Assert
    expect(
      screen.getByRole("heading", { level: 1, name: "出力管理" })
    ).toBeInTheDocument()
  })
})

describe("App master persona screen", () => {
  beforeEach(() => {
    window.history.replaceState(null, "", "#master-persona")
  })

  test("contract factory から受け取った view model でペルソナ一覧見出しを描画する", () => {
    const masterPersonaController = new MasterPersonaScreenControllerFake()

    render(App, {
      props: {
        createMasterDictionaryScreenController: () =>
          new MasterDictionaryScreenControllerFake(),
        createMasterPersonaScreenController: () => masterPersonaController
      }
    })

    expect(
      screen.getByRole("heading", { level: 3, name: "ペルソナ一覧" })
    ).toBeInTheDocument()
  })

  test("hero title と主要 CTA を公開文言で表示する", () => {
    const masterPersonaController = new MasterPersonaScreenControllerFake(
      buildMasterPersonaScreenViewModel({
        canStartGeneration: true,
        selectedFileReference: "/tmp/master-persona.json"
      })
    )

    render(App, {
      props: {
        createMasterDictionaryScreenController: () =>
          new MasterDictionaryScreenControllerFake(),
        createMasterPersonaScreenController: () => masterPersonaController
      }
    })

    expect(
      screen.getByRole("heading", { level: 1, name: "マスターペルソナ作成" })
    ).toBeInTheDocument()
    expect(screen.getByText("ペルソナを作成")).toBeInTheDocument()
  })

  test("通常表示に Gateway と preview 状態を出さない", () => {
    const masterPersonaController = new MasterPersonaScreenControllerFake()

    render(App, {
      props: {
        createMasterDictionaryScreenController: () =>
          new MasterDictionaryScreenControllerFake(),
        createMasterPersonaScreenController: () => masterPersonaController
      }
    })

    expect(screen.queryByText(/Gateway/i)).not.toBeInTheDocument()
    expect(screen.queryByText("preview 状態")).not.toBeInTheDocument()
  })

  test("一覧初期表示はプラグイン名と NPC 名中心で不要列を行に出さない", () => {
    const masterPersonaController = new MasterPersonaScreenControllerFake(
      buildMasterPersonaScreenViewModel({
        items: [
          {
            identityKey: "TestPersonaPluginA.esp:FE01A812:NPC_",
            targetPlugin: "TestPersonaPluginA.esp",
            formId: "FE01A812",
            recordType: "NPC_",
            editorId: "TEST_NPC_A",
            displayName: "Test NPC A",
            voiceType: "TestVoiceA",
            className: "TestClassA",
            sourcePlugin: "Skyrim.esm",
            personaSummary: "乾いた率直さで応じる。",
            updatedAt: "2026-04-15T09:42:00Z"
          }
        ],
        selectedEntry: {
          identityKey: "TestPersonaPluginA.esp:FE01A812:NPC_",
          targetPlugin: "TestPersonaPluginA.esp",
          formId: "FE01A812",
          recordType: "NPC_",
          editorId: "TEST_NPC_A",
          displayName: "Test NPC A",
          voiceType: "TestVoiceA",
          className: "TestClassA",
          sourcePlugin: "Skyrim.esm",
          personaSummary: "乾いた率直さで応じる。",
          updatedAt: "2026-04-15T09:42:00Z",
          personaBody: "短く本音を置く。",
          runLockReason: "更新と削除を行えます"
        } as MasterPersonaDetail
      })
    )

    render(App, {
      props: {
        createMasterDictionaryScreenController: () =>
          new MasterDictionaryScreenControllerFake(),
        createMasterPersonaScreenController: () => masterPersonaController
      }
    })

    const listHeading = screen.getByRole("heading", {
      level: 3,
      name: "ペルソナ一覧"
    })
    const listSection = listHeading.closest("section")
    if (!(listSection instanceof HTMLElement)) {
      throw new Error("ペルソナ一覧 section が見つかりません")
    }
    const listQuery = within(listSection)
    const targetRow = listQuery.getByRole("row", { name: "Test NPC A" })
    expect(listQuery.getAllByText("TestPersonaPluginA.esp").length).toBeGreaterThan(0)
    expect(targetRow).toBeInTheDocument()
    expect(within(targetRow).queryByText("TestClassA")).not.toBeInTheDocument()
    expect(within(targetRow).queryByText("Skyrim.esm")).not.toBeInTheDocument()
    expect(
      within(targetRow).queryByText("乾いた率直さで応じる。")
    ).not.toBeInTheDocument()
  })

  test("編集モーダルの公開ラベルを維持する", () => {
    const masterPersonaController = new MasterPersonaScreenControllerFake(
      buildMasterPersonaScreenViewModel({
        modalState: "edit"
      })
    )

    render(App, {
      props: {
        createMasterDictionaryScreenController: () =>
          new MasterDictionaryScreenControllerFake(),
        createMasterPersonaScreenController: () => masterPersonaController
      }
    })

    expect(
      screen.getByRole("heading", { level: 3, name: "ペルソナを編集" })
    ).toBeInTheDocument()
    expect(
      screen.getByRole("button", { name: "編集内容を保存" })
    ).toBeInTheDocument()
  })

  test("削除モーダルの公開ラベルを維持する", () => {
    const masterPersonaController = new MasterPersonaScreenControllerFake(
      buildMasterPersonaScreenViewModel({
        modalState: "delete"
      })
    )

    render(App, {
      props: {
        createMasterDictionaryScreenController: () =>
          new MasterDictionaryScreenControllerFake(),
        createMasterPersonaScreenController: () => masterPersonaController
      }
    })

    expect(
      screen.getByRole("heading", { level: 3, name: "ペルソナを削除しますか" })
    ).toBeInTheDocument()
    expect(screen.getByRole("button", { name: "削除する" })).toBeInTheDocument()
  })

  test("render 時に master persona controller.mount を呼ぶ", async () => {
    const masterPersonaController = new MasterPersonaScreenControllerFake()

    render(App, {
      props: {
        createMasterDictionaryScreenController: () =>
          new MasterDictionaryScreenControllerFake(),
        createMasterPersonaScreenController: () => masterPersonaController
      }
    })

    await waitFor(() => {
      expect(masterPersonaController.mount).toHaveBeenCalledTimes(1)
    })
  })

  test("create 導線を出さず prompt template は説明だけを表示する", async () => {
    const masterPersonaController = new MasterPersonaScreenControllerFake()

    render(App, {
      props: {
        createMasterDictionaryScreenController: () =>
          new MasterDictionaryScreenControllerFake(),
        createMasterPersonaScreenController: () => masterPersonaController
      }
    })

    await waitFor(() => {
      expect(screen.queryByRole("button", { name: "新規作成" })).toBeNull()
      expect(screen.queryByLabelText("プロンプトテンプレート")).toBeNull()
      expect(
        screen.getByText(
          "プロンプトテンプレートは画面入力では変更せず、実装側の説明文として固定しています。"
        )
      ).toBeInTheDocument()
    })
  })

  test("AI service dropdown は real provider 3件だけを表示し fake provider を表示しない", async () => {
    const masterPersonaController = new MasterPersonaScreenControllerFake(
      buildMasterPersonaScreenViewModel({
        aiSettings: {
          provider: "gemini",
          model: "gemini-2.5-pro",
          executionMethod: "single_request"
        }
      })
    )

    render(App, {
      props: {
        createMasterDictionaryScreenController: () =>
          new MasterDictionaryScreenControllerFake(),
        createMasterPersonaScreenController: () => masterPersonaController
      }
    })

    const aiServiceSelect = await screen.findByRole("combobox", {
      name: "AI サービス"
    })
    const options = within(aiServiceSelect).getAllByRole("option")
    const optionValues = options.map(
      (option) => option.getAttribute("value") ?? ""
    )
    const optionLabels = options.map(
      (option) => option.textContent?.trim() ?? ""
    )

    expect(optionValues).toEqual(["gemini", "lm_studio", "xai"])
    expect(optionLabels).toEqual(["Gemini", "LM Studio", "xAI"])
    expect(optionValues).not.toContain("fake")
  })

  test("AI service pill は canonical provider ID ではなく表示名を表示する", async () => {
    const masterPersonaController = new MasterPersonaScreenControllerFake(
      buildMasterPersonaScreenViewModel({
        aiSettings: {
          provider: "lm_studio",
          model: "llama3",
          executionMethod: "single_request"
        },
        aiProviderLabel: "LM Studio"
      })
    )

    render(App, {
      props: {
        createMasterDictionaryScreenController: () =>
          new MasterDictionaryScreenControllerFake(),
        createMasterPersonaScreenController: () => masterPersonaController
      }
    })

    const settingsHeading = await screen.findByRole("heading", {
      level: 3,
      name: "この画面で使う AI 設定"
    })
    const settingsHeader = settingsHeading.closest(".section-head")
    if (!(settingsHeader instanceof HTMLElement)) {
      throw new Error("AI 設定 header が見つかりません")
    }

    const providerPill = settingsHeader.querySelector(".status-pill")
    expect(providerPill).toHaveTextContent("LM Studio")
    expect(providerPill).not.toHaveTextContent("lm_studio")
  })

  test("plugin dropdown は plugin filter option だけを表示する", async () => {
    const masterPersonaController = new MasterPersonaScreenControllerFake(
      buildMasterPersonaScreenViewModel({
        pluginOptions: [
          { value: "", label: "すべてのプラグイン" },
          {
            value: "TestPersonaPluginA.esp",
            label: "TestPersonaPluginA.esp (1)"
          },
          {
            value: "TestPersonaPluginB.esp",
            label: "TestPersonaPluginB.esp (1)"
          }
        ],
        runStatus: {
          runState: "生成中",
          targetPlugin: "TestPersonaPluginA.esp",
          processedCount: 1,
          successCount: 0,
          existingSkipCount: 1,
          zeroDialogueSkipCount: 0,
          genericNpcCount: 0,
          currentActorLabel: "Test NPC A",
          message: "ペルソナを作成中"
        }
      })
    )

    render(App, {
      props: {
        createMasterDictionaryScreenController: () =>
          new MasterDictionaryScreenControllerFake(),
        createMasterPersonaScreenController: () => masterPersonaController
      }
    })

    const pluginSelect = await screen.findByRole("combobox", {
      name: "プラグイン"
    })
    const optionTexts = within(pluginSelect)
      .getAllByRole("option")
      .map((option) => option.textContent?.trim() ?? "")

    expect(optionTexts).toEqual([
      "すべてのプラグイン",
      "TestPersonaPluginA.esp (1)",
      "TestPersonaPluginB.esp (1)"
    ])
    expect(optionTexts).not.toContain("生成中")
    expect(optionTexts).not.toContain("完了")
  })

  test("closed-by-default modal は hidden 時に編集削除 modal は非表示で非対話のままにする", () => {
    render(App, {
      props: {
        createMasterDictionaryScreenController: () =>
          new MasterDictionaryScreenControllerFake(),
        createMasterPersonaScreenController: () =>
          new MasterPersonaScreenControllerFake(
            buildMasterPersonaScreenViewModel({
              modalState: null
            })
          )
      }
    })

    const editModal = document.querySelector("#editModal")
    const deleteModal = document.querySelector("#deleteModal")

    expect(editModal).toHaveAttribute("hidden")
    expect(editModal).not.toHaveClass("is-open")
    expect(deleteModal).toHaveAttribute("hidden")
    expect(deleteModal).not.toHaveClass("is-open")
  })

  test("page-local AI settings は他画面へ漏れない", async () => {
    const user = userEvent.setup()
    const masterPersonaController = new MasterPersonaScreenControllerFake(
      buildMasterPersonaScreenViewModel({
        aiSettings: {
          provider: "gemini",
          model: "persona-only-model",
          executionMethod: "single_request"
        },
        modelOptions: [
          { modelId: "persona-only-model", label: "persona-only-model" }
        ]
      })
    )

    render(App, {
      props: {
        createMasterDictionaryScreenController: () =>
          new MasterDictionaryScreenControllerFake(),
        createMasterPersonaScreenController: () => masterPersonaController
      }
    })

    expect(screen.getByDisplayValue("persona-only-model")).toBeInTheDocument()

    await user.click(
      within(getGlobalNavigation()).getByRole("link", { name: "マスター辞書" })
    )

    await waitFor(() => {
      expect(
        screen.getByRole("heading", { level: 1, name: "マスター辞書" })
      ).toBeInTheDocument()
      expect(screen.queryByDisplayValue("persona-only-model")).toBeNull()
      expect(screen.queryByRole("combobox", { name: "AI サービス" })).toBeNull()
    })
  })

  test("persona-read-detail-cutover: plugin filter dropdown を表示する", async () => {
    // Arrange
    const masterPersonaController = new MasterPersonaScreenControllerFake(
      buildMasterPersonaScreenViewModel({
        pluginOptions: [
          { value: "", label: "すべてのプラグイン" },
          {
            value: "TestPersonaPluginA.esp",
            label: "TestPersonaPluginA.esp (1)"
          }
        ]
      })
    )
    render(App, {
      props: {
        createMasterDictionaryScreenController: () =>
          new MasterDictionaryScreenControllerFake(),
        createMasterPersonaScreenController: () => masterPersonaController
      }
    })

    // Assert
    const pluginSelect = await screen.findByRole("combobox", {
      name: "プラグイン"
    })
    expect(pluginSelect).toBeInTheDocument()
    expect(within(pluginSelect).getAllByRole("option")).toHaveLength(2)
  })

  test("persona-read-detail-cutover: selectedEntry の FormID と EditorID を一覧の展開行に表示する", async () => {
    // Arrange
    render(App, {
      props: {
        createMasterDictionaryScreenController: () =>
          new MasterDictionaryScreenControllerFake(),
        createMasterPersonaScreenController: () =>
          new MasterPersonaScreenControllerFake(
            buildMasterPersonaScreenViewModel()
          )
      }
    })

    // Assert
    await waitFor(() => {
      const metadataList = screen.getByLabelText("Test NPC A のメタデータ")
      expect(metadataList).toBeInTheDocument()
      expect(metadataList.textContent).toContain("FE01A812")
      expect(metadataList.textContent).toContain("TEST_NPC_A")
      expect(metadataList.textContent).toContain("短く本音を置く。")
    })
  })

  test("persona-read-detail-cutover: ダイアログ一覧ボタンは表示されない", async () => {
    // Arrange
    render(App, {
      props: {
        createMasterDictionaryScreenController: () =>
          new MasterDictionaryScreenControllerFake(),
        createMasterPersonaScreenController: () =>
          new MasterPersonaScreenControllerFake()
      }
    })

    // Assert
    await waitFor(() => {
      expect(
        screen.queryByRole("button", { name: "ダイアログ一覧" })
      ).toBeNull()
    })
  })

  test("persona-read-detail-cutover: 詳細グリッドに 収録元ファイル ラベルは表示されない", async () => {
    // Arrange
    render(App, {
      props: {
        createMasterDictionaryScreenController: () =>
          new MasterDictionaryScreenControllerFake(),
        createMasterPersonaScreenController: () =>
          new MasterPersonaScreenControllerFake()
      }
    })

    // Assert
    await waitFor(() => {
      expect(screen.queryByText("収録元ファイル")).toBeNull()
    })
  })

  test("persona-read-detail-cutover: 詳細グリッドに ダイアログ数 ラベルは表示されない", async () => {
    // Arrange
    render(App, {
      props: {
        createMasterDictionaryScreenController: () =>
          new MasterDictionaryScreenControllerFake(),
        createMasterPersonaScreenController: () =>
          new MasterPersonaScreenControllerFake()
      }
    })

    // Assert
    await waitFor(() => {
      expect(screen.queryByText("ダイアログ数")).toBeNull()
    })
  })

  test("persona-read-detail-cutover: edit modal は identity / snapshot fields を editable input として表示しない", async () => {
    // Arrange
    render(App, {
      props: {
        createMasterDictionaryScreenController: () =>
          new MasterDictionaryScreenControllerFake(),
        createMasterPersonaScreenController: () =>
          new MasterPersonaScreenControllerFake(
            buildMasterPersonaScreenViewModel({ modalState: "edit" })
          )
      }
    })

    // Assert
    await waitFor(() => {
      expect(document.querySelector("#editFormIdInput")).toBeNull()
      expect(document.querySelector("#editEditorIdInput")).toBeNull()
      expect(document.querySelector("#editRaceInput")).toBeNull()
      expect(document.querySelector("#editSexInput")).toBeNull()
      expect(document.querySelector("#editVoiceTypeInput")).toBeNull()
      expect(document.querySelector("#editClassNameInput")).toBeNull()
      expect(document.querySelector("#editSourcePluginInput")).toBeNull()
    })
  })

  test("persona-json-preview-cutover: previewStats に 会話が見つからない ラベルは表示されない", async () => {
    // Arrange
    render(App, {
      props: {
        createMasterDictionaryScreenController: () =>
          new MasterDictionaryScreenControllerFake(),
        createMasterPersonaScreenController: () =>
          new MasterPersonaScreenControllerFake()
      }
    })

    // Assert
    await waitFor(() => {
      const previewStats = document.querySelector("#previewStats")
      expect(previewStats).toBeInTheDocument()
      expect(
        within(previewStats as HTMLElement).queryByText("会話が見つからない")
      ).toBeNull()
    })
  })

  test("persona-json-preview-cutover: previewStats に 汎用NPC ラベルは表示されない", async () => {
    // Arrange
    render(App, {
      props: {
        createMasterDictionaryScreenController: () =>
          new MasterDictionaryScreenControllerFake(),
        createMasterPersonaScreenController: () =>
          new MasterPersonaScreenControllerFake()
      }
    })

    // Assert
    await waitFor(() => {
      const previewStats = document.querySelector("#previewStats")
      expect(previewStats).toBeInTheDocument()
      expect(
        within(previewStats as HTMLElement).queryByText("汎用NPC")
      ).toBeNull()
    })
  })

  test("persona-generation-cutover: ペルソナを作成 ボタンをクリックすると controller.executeGeneration を呼ぶ", async () => {
    // Arrange
    const user = userEvent.setup()
    const masterPersonaController = new MasterPersonaScreenControllerFake(
      buildMasterPersonaScreenViewModel({
        canStartGeneration: true,
        isRunActive: false
      })
    )
    render(App, {
      props: {
        createMasterDictionaryScreenController: () =>
          new MasterDictionaryScreenControllerFake(),
        createMasterPersonaScreenController: () => masterPersonaController
      }
    })

    // Act
    const generateButton = await screen.findByRole("button", {
      name: "ペルソナを作成"
    })
    await user.click(generateButton)

    // Assert
    expect(masterPersonaController.executeGeneration).toHaveBeenCalledTimes(1)
  })

  test("persona-generation-cutover: isRunActive false のとき 一時停止 と 中止 ボタンは disabled", async () => {
    // Arrange
    const masterPersonaController = new MasterPersonaScreenControllerFake(
      buildMasterPersonaScreenViewModel({ isRunActive: false })
    )
    render(App, {
      props: {
        createMasterDictionaryScreenController: () =>
          new MasterDictionaryScreenControllerFake(),
        createMasterPersonaScreenController: () => masterPersonaController
      }
    })

    // Assert
    await waitFor(() => {
      expect(screen.getByRole("button", { name: "一時停止" })).toBeDisabled()
      expect(screen.getByRole("button", { name: "中止" })).toBeDisabled()
    })
  })

  test("persona-generation-cutover: isRunActive true のとき 一時停止 と 中止 ボタンは enabled", async () => {
    // Arrange
    const masterPersonaController = new MasterPersonaScreenControllerFake(
      buildMasterPersonaScreenViewModel({
        isRunActive: true,
        runStatus: {
          runState: "生成中",
          targetPlugin: "TestPersonaPluginA.esp",
          processedCount: 2,
          successCount: 2,
          existingSkipCount: 0,
          zeroDialogueSkipCount: 0,
          genericNpcCount: 0,
          currentActorLabel: "Test NPC A",
          message: "ペルソナを作成中"
        }
      })
    )
    render(App, {
      props: {
        createMasterDictionaryScreenController: () =>
          new MasterDictionaryScreenControllerFake(),
        createMasterPersonaScreenController: () => masterPersonaController
      }
    })

    // Assert
    await waitFor(() => {
      expect(
        screen.getByRole("button", { name: "一時停止" })
      ).not.toBeDisabled()
      expect(screen.getByRole("button", { name: "中止" })).not.toBeDisabled()
    })
  })

  test("persona-generation-cutover: runStatus.existingSkipCount を 既に作成済み として表示する", async () => {
    // Arrange
    const masterPersonaController = new MasterPersonaScreenControllerFake(
      buildMasterPersonaScreenViewModel({
        runStatus: {
          runState: "完了",
          targetPlugin: "TestPersonaPluginA.esp",
          processedCount: 7,
          successCount: 7,
          existingSkipCount: 3,
          zeroDialogueSkipCount: 0,
          genericNpcCount: 0,
          currentActorLabel: "",
          message: "完了"
        }
      })
    )
    render(App, {
      props: {
        createMasterDictionaryScreenController: () =>
          new MasterDictionaryScreenControllerFake(),
        createMasterPersonaScreenController: () => masterPersonaController
      }
    })

    // Assert: generation never overwrites existing — existingSkipCount が UI に表示される
    await waitFor(() => {
      expect(screen.getByText(/既に作成済み/)).toBeInTheDocument()
    })
  })

  test("persona-edit-delete-cutover: edit modal は displayName を editable input として表示しない", async () => {
    // Arrange
    render(App, {
      props: {
        createMasterDictionaryScreenController: () =>
          new MasterDictionaryScreenControllerFake(),
        createMasterPersonaScreenController: () =>
          new MasterPersonaScreenControllerFake(
            buildMasterPersonaScreenViewModel({ modalState: "edit" })
          )
      }
    })

    // Assert: displayName 入力は edit modal から除去されている (RED: 現在は存在する)
    await waitFor(() => {
      expect(document.querySelector("#editDisplayNameInput")).toBeNull()
    })
  })

  test("persona-edit-delete-cutover: edit modal は personaSummary を editable input として表示する", async () => {
    // Arrange
    render(App, {
      props: {
        createMasterDictionaryScreenController: () =>
          new MasterDictionaryScreenControllerFake(),
        createMasterPersonaScreenController: () =>
          new MasterPersonaScreenControllerFake(
            buildMasterPersonaScreenViewModel({ modalState: "edit" })
          )
      }
    })

    // Assert: personaSummary 入力が edit modal に追加されている (RED: 現在は存在しない)
    await waitFor(() => {
      expect(document.querySelector("#editPersonaSummaryInput")).not.toBeNull()
    })
  })

  test("persona-edit-delete-cutover: edit modal は speechStyle を editable input として表示する", async () => {
    // Arrange
    render(App, {
      props: {
        createMasterDictionaryScreenController: () =>
          new MasterDictionaryScreenControllerFake(),
        createMasterPersonaScreenController: () =>
          new MasterPersonaScreenControllerFake(
            buildMasterPersonaScreenViewModel({ modalState: "edit" })
          )
      }
    })

    // Assert: speechStyle 入力が edit modal に追加されている (RED: 現在は存在しない)
    await waitFor(() => {
      expect(document.querySelector("#editSpeechStyleInput")).not.toBeNull()
    })
  })
})

describe("App term translation phase screen", () => {
  beforeEach(() => {
    window.history.replaceState(null, "", "#translation-management")
  })

  test("未完了ジョブ未選択では翻訳段階 controller factory を起動しない", () => {
    const controller = new TermTranslationPhaseScreenControllerFake()
    const personaController = new PersonaGenerationPhaseScreenControllerFake()
    const jobManagementController =
      new TranslationJobManagementScreenControllerFake()
    const createTermTranslationPhaseScreenController = vi.fn(() => controller)
    const createPersonaGenerationPhaseScreenController = vi.fn(
      () => personaController
    )
    const createTranslationJobManagementScreenController = vi.fn(
      () => jobManagementController
    )

    render(App, {
      props: {
        createMasterDictionaryScreenController: () =>
          new MasterDictionaryScreenControllerFake(),
        createMasterPersonaScreenController: () =>
          new MasterPersonaScreenControllerFake(),
        createPersonaGenerationPhaseScreenController,
        createTermTranslationPhaseScreenController,
        createTranslationJobManagementScreenController
      }
    })

    expect(
      screen.getByRole("heading", { level: 2, name: "未完了ジョブ一覧" })
    ).toBeInTheDocument()
    expect(
      screen.queryByRole("tab", { name: /Job Run|実行/ })
    ).not.toBeInTheDocument()
    expect(createTermTranslationPhaseScreenController).toHaveBeenCalledTimes(0)
    expect(createPersonaGenerationPhaseScreenController).toHaveBeenCalledTimes(
      0
    )
    expect(controller.mount).toHaveBeenCalledTimes(0)
    expect(personaController.mount).toHaveBeenCalledTimes(0)
  })
})

describe("App master dictionary screen", () => {
  beforeEach(() => {
    window.history.replaceState(null, "", "#master-dictionary")
  })

  test("contract factory から受け取った view model で辞書一覧見出しを描画する", () => {
    // Arrange
    const controller = new MasterDictionaryScreenControllerFake(
      buildMasterDictionaryScreenViewModel({
        importSummary: {
          fileName: "master.xml",
          importedCount: 2,
          updatedCount: 1,
          totalCount: 3,
          selectedSource: "Dragon Priest"
        },
        importStage: "done",
        importStatusValue: "完了"
      })
    )

    // Act
    renderApp(controller)

    // Assert
    expect(
      screen.getByRole("heading", { level: 3, name: "辞書一覧" })
    ).toBeInTheDocument()
  })

  test("render 時に controller.mount を呼ぶ", async () => {
    // Arrange
    const controller = new MasterDictionaryScreenControllerFake()

    // Act
    renderApp(controller)

    // Assert
    await waitFor(() => {
      expect(controller.mount).toHaveBeenCalledTimes(1)
    })
  })

  test("unmount 時に controller.dispose を呼ぶ", () => {
    // Arrange
    const controller = new MasterDictionaryScreenControllerFake()
    const view = renderAppView(controller)

    // Act
    view.unmount()

    // Assert
    expect(controller.dispose).toHaveBeenCalledTimes(1)
  })

  test("選択中エントリの detail title を描画する", () => {
    // Arrange
    renderApp(
      new MasterDictionaryScreenControllerFake(
        buildMasterDictionaryScreenViewModel({
          importSummary: {
            fileName: "master.xml",
            importedCount: 2,
            updatedCount: 1,
            totalCount: 3,
            selectedSource: "Dragon Priest"
          },
          importStage: "done",
          importStatusValue: "完了"
        })
      )
    )

    // Act
    const detailTitle = document.querySelector("#detailTitle")

    // Assert
    expect(detailTitle).toHaveTextContent("Dragon Priest")
  })

  test("import 完了 status value を描画する", () => {
    // Arrange
    renderApp(
      new MasterDictionaryScreenControllerFake(
        buildMasterDictionaryScreenViewModel({
          importSummary: {
            fileName: "master.xml",
            importedCount: 2,
            updatedCount: 1,
            totalCount: 3,
            selectedSource: "Dragon Priest"
          },
          importStage: "done",
          importStatusValue: "完了"
        })
      )
    )

    // Act
    const importStatusValue = document.querySelector("#importStatusValue")

    // Assert
    expect(importStatusValue).toHaveTextContent("完了")
  })

  test("import 完了 summary の選択ソースを描画する", () => {
    // Arrange
    renderApp(
      new MasterDictionaryScreenControllerFake(
        buildMasterDictionaryScreenViewModel({
          importSummary: {
            fileName: "master.xml",
            importedCount: 2,
            updatedCount: 1,
            totalCount: 3,
            selectedSource: "Dragon Priest"
          },
          importStage: "done",
          importStatusValue: "完了"
        })
      )
    )

    // Act
    const importResultSelection = document.querySelector(
      "#importResultSelection"
    )

    // Assert
    expect(importResultSelection).toHaveTextContent("Dragon Priest")
  })

  test("一覧行クリックで controller.selectRow を呼ぶ", async () => {
    // Arrange
    const user = userEvent.setup()
    const controller = renderApp()

    await waitFor(() => {
      expect(controller.mount).toHaveBeenCalledTimes(1)
    })

    // Act
    await user.click(
      screen.getByRole("button", { name: /ドラゴン・プリースト/ })
    )

    // Assert
    expect(controller.selectRow).toHaveBeenCalledWith("101")
  })

  test("新規登録ボタンで controller.openCreateModal を呼ぶ", async () => {
    // Arrange
    const user = userEvent.setup()
    const controller = renderApp()

    // Act
    await user.click(screen.getByRole("button", { name: "新規登録" }))

    // Assert
    expect(controller.openCreateModal).toHaveBeenCalledTimes(1)
  })

  test("更新ボタンで controller.openEditModal を呼ぶ", async () => {
    // Arrange
    const user = userEvent.setup()
    const controller = renderApp()

    // Act
    await user.click(screen.getByRole("button", { name: "更新" }))

    // Assert
    expect(controller.openEditModal).toHaveBeenCalledTimes(1)
  })

  test("削除ボタンで controller.openDeleteModal を呼ぶ", async () => {
    // Arrange
    const user = userEvent.setup()
    const controller = renderApp()

    // Act
    await user.click(screen.getByRole("button", { name: "削除" }))

    // Assert
    expect(controller.openDeleteModal).toHaveBeenCalledTimes(1)
  })

  test("検索入力で controller.handleSearchInput を呼ぶ", async () => {
    // Arrange
    const user = userEvent.setup()
    const controller = renderApp(
      new MasterDictionaryScreenControllerFake(
        buildMasterDictionaryScreenViewModel({
          totalCount: 60,
          totalPages: 2,
          pageStatusText: "1 - 30 件を表示"
        })
      )
    )

    // Act
    await user.type(screen.getByRole("searchbox", { name: "検索" }), "Reach")

    // Assert
    expect(controller.handleSearchInput).toHaveBeenCalled()
  })

  test("カテゴリ変更で controller.handleCategoryChange を呼ぶ", async () => {
    // Arrange
    const user = userEvent.setup()
    const controller = renderApp(
      new MasterDictionaryScreenControllerFake(
        buildMasterDictionaryScreenViewModel({
          totalCount: 60,
          totalPages: 2,
          pageStatusText: "1 - 30 件を表示"
        })
      )
    )

    // Act
    await user.selectOptions(
      screen.getByRole("combobox", { name: "カテゴリ" }),
      "地名"
    )

    // Assert
    expect(controller.handleCategoryChange).toHaveBeenCalled()
  })

  test("次ページボタンで controller.goToNextPage を呼ぶ", async () => {
    // Arrange
    const user = userEvent.setup()
    const controller = renderApp(
      new MasterDictionaryScreenControllerFake(
        buildMasterDictionaryScreenViewModel({
          totalCount: 60,
          totalPages: 2,
          pageStatusText: "1 - 30 件を表示"
        })
      )
    )

    // Act
    await user.click(screen.getByRole("button", { name: "次の30件" }))

    // Assert
    expect(controller.goToNextPage).toHaveBeenCalledTimes(1)
  })

  test("XML 選択で controller.stageXmlImport を呼ぶ", async () => {
    // Arrange
    const user = userEvent.setup()
    const controller = renderApp()
    const xmlInput = document.querySelector("#xmlFileInput")
    if (!(xmlInput instanceof HTMLInputElement)) {
      throw new Error("xmlFileInput が見つかりません")
    }
    const xmlFile = createXmlFile("<Root />")

    // Act
    await user.upload(xmlInput, xmlFile)

    // Assert
    expect(controller.stageXmlImport).toHaveBeenCalledWith(xmlFile)
  })

  test("staged file を受けると import 実行ボタンを表示する", async () => {
    // Arrange
    const controller = renderApp()
    const xmlFile = createXmlFile("<Root />")

    // Act
    controller.pushViewModel(
      buildMasterDictionaryScreenViewModel({
        hasStagedFile: true,
        selectedFileName: xmlFile.name,
        selectedFileReference: xmlFile.name,
        importStage: "ready",
        importStatusValue: "取込待ち"
      })
    )

    // Assert
    await waitFor(() => {
      expect(
        screen.getByRole("button", { name: "この XML を取り込む" })
      ).toBeInTheDocument()
    })
  })

  test("選び直す操作で controller.resetImportSelection を呼ぶ", async () => {
    // Arrange
    const user = userEvent.setup()
    const controller = renderApp()
    const xmlFile = createXmlFile("<Root />")

    controller.pushViewModel(
      buildMasterDictionaryScreenViewModel({
        hasStagedFile: true,
        selectedFileName: xmlFile.name,
        selectedFileReference: xmlFile.name,
        importStage: "ready",
        importStatusValue: "取込待ち"
      })
    )

    await waitFor(() => {
      expect(
        screen.getByRole("button", { name: "この XML を取り込む" })
      ).toBeInTheDocument()
    })

    // Act
    await user.click(screen.getByRole("button", { name: "選び直す" }))

    // Assert
    expect(controller.resetImportSelection).toHaveBeenCalledTimes(1)
  })

  test("import 実行操作で controller.startImport を呼ぶ", async () => {
    // Arrange
    const user = userEvent.setup()
    const controller = renderApp(
      new MasterDictionaryScreenControllerFake(
        buildMasterDictionaryScreenViewModel({
          hasStagedFile: true,
          selectedFileName: "master.xml",
          selectedFileReference: "master.xml",
          importStage: "ready",
          importStatusValue: "取込待ち"
        })
      )
    )

    // Act
    await user.click(
      screen.getByRole("button", { name: "この XML を取り込む" })
    )

    // Assert
    expect(controller.startImport).toHaveBeenCalledTimes(1)
  })

  test("running view model push で import status value を更新する", async () => {
    // Arrange
    const controller = renderApp(
      new MasterDictionaryScreenControllerFake(
        buildMasterDictionaryScreenViewModel({
          hasStagedFile: true,
          selectedFileName: "master.xml",
          selectedFileReference: "master.xml",
          importStage: "ready",
          importStatusValue: "取込待ち"
        })
      )
    )

    // Act
    controller.pushViewModel(
      buildMasterDictionaryScreenViewModel({
        hasStagedFile: true,
        selectedFileName: "master.xml",
        selectedFileReference: "master.xml",
        importStage: "running",
        isImportRunning: true,
        importProgress: 78,
        importStatusValue: "取込中"
      })
    )

    // Assert
    await waitFor(() => {
      expect(document.querySelector("#importStatusValue")).toHaveTextContent(
        "取込中"
      )
    })
  })

  test("running view model push で import progress bar 幅を更新する", async () => {
    // Arrange
    const controller = renderApp(
      new MasterDictionaryScreenControllerFake(
        buildMasterDictionaryScreenViewModel({
          hasStagedFile: true,
          selectedFileName: "master.xml",
          selectedFileReference: "master.xml",
          importStage: "ready",
          importStatusValue: "取込待ち"
        })
      )
    )

    // Act
    controller.pushViewModel(
      buildMasterDictionaryScreenViewModel({
        hasStagedFile: true,
        selectedFileName: "master.xml",
        selectedFileReference: "master.xml",
        importStage: "running",
        isImportRunning: true,
        importProgress: 78,
        importStatusValue: "取込中"
      })
    )

    // Assert
    await waitFor(() => {
      expect(document.querySelector("#importProgressFill")).toHaveAttribute(
        "style",
        "width: 78%;"
      )
    })
  })

  test("完了 view model push で完了表示を反映する", async () => {
    // Arrange
    const controller = renderApp()

    // Act
    controller.pushViewModel(
      buildMasterDictionaryScreenViewModel({
        entries: [
          {
            id: "201",
            source: "Allowed Book Source",
            translation: "許可された本の訳語",
            category: "書籍",
            origin: "XML取込",
            updatedAt: "2026-04-12 00:00"
          }
        ],
        selectedEntry: createImportedBookSelectedEntry(),
        selectedId: "201",
        totalCount: 1,
        query: "",
        category: "すべて",
        hasStagedFile: true,
        selectedFileName: "master.xml",
        selectedFileReference: "master.xml",
        importStage: "done",
        importProgress: 100,
        importStatusValue: "完了",
        importSummary: {
          fileName: "master.xml",
          importedCount: 1,
          updatedCount: 0,
          totalCount: 1,
          selectedSource: "Allowed Book Source"
        },
        listHeadline: "1 件のエントリを表示しています。",
        selectionStatusText: "Allowed Book Source を選択中",
        detailSublineText: "XML取込 / 最終更新 2026-04-12 00:00"
      })
    )

    // Assert
    await waitFor(() => {
      expect(screen.getByText("完了")).toBeInTheDocument()
    })
  })

  test("完了 view model push で entry 単位の import 集計ラベルを表示する", async () => {
    // Arrange
    const controller = renderApp()

    // Act
    controller.pushViewModel(
      buildMasterDictionaryScreenViewModel({
        entries: [
          {
            id: "201",
            source: "Allowed Book Source",
            translation: "許可された本の訳語",
            category: "書籍",
            origin: "XML取込",
            updatedAt: "2026-04-12 00:00"
          }
        ],
        selectedEntry: createImportedBookSelectedEntry(),
        selectedId: "201",
        totalCount: 1,
        query: "",
        category: "すべて",
        hasStagedFile: true,
        selectedFileName: "master.xml",
        selectedFileReference: "master.xml",
        importStage: "done",
        importProgress: 100,
        importStatusValue: "完了",
        importSummary: {
          fileName: "master.xml",
          importedCount: 1,
          updatedCount: 947,
          totalCount: 740,
          selectedSource: "Allowed Book Source"
        },
        listHeadline: "740 件のエントリを表示しています。",
        selectionStatusText: "Allowed Book Source を選択中",
        detailSublineText: "XML取込 / 最終更新 2026-04-12 00:00"
      })
    )

    // Assert
    await waitFor(() => {
      expect(screen.getByText("更新済みエントリ件数")).toBeInTheDocument()
      expect(screen.getByText("取込後の保存済み一覧件数")).toBeInTheDocument()
      expect(screen.getByText("新規追加 1 件")).toBeInTheDocument()
      expect(
        screen.getByText(/件数は保存済みエントリ単位で集計しています。/)
      ).toBeInTheDocument()
    })
  })

  test("完了 view model push で検索入力値を初期状態へ反映する", async () => {
    // Arrange
    const controller = renderApp()

    // Act
    controller.pushViewModel(
      buildMasterDictionaryScreenViewModel({
        entries: [
          {
            id: "201",
            source: "Allowed Book Source",
            translation: "許可された本の訳語",
            category: "書籍",
            origin: "XML取込",
            updatedAt: "2026-04-12 00:00"
          }
        ],
        selectedEntry: createImportedBookSelectedEntry(),
        selectedId: "201",
        totalCount: 1,
        query: "",
        category: "すべて",
        hasStagedFile: true,
        selectedFileName: "master.xml",
        selectedFileReference: "master.xml",
        importStage: "done",
        importProgress: 100,
        importStatusValue: "完了",
        importSummary: {
          fileName: "master.xml",
          importedCount: 1,
          updatedCount: 0,
          totalCount: 1,
          selectedSource: "Allowed Book Source"
        },
        listHeadline: "1 件のエントリを表示しています。",
        selectionStatusText: "Allowed Book Source を選択中",
        detailSublineText: "XML取込 / 最終更新 2026-04-12 00:00"
      })
    )

    // Assert
    await waitFor(() => {
      expect(screen.getByRole("searchbox", { name: "検索" })).toHaveValue("")
    })
  })

  test("完了 view model push でカテゴリ選択値を初期状態へ反映する", async () => {
    // Arrange
    const controller = renderApp()

    // Act
    controller.pushViewModel(
      buildMasterDictionaryScreenViewModel({
        entries: [
          {
            id: "201",
            source: "Allowed Book Source",
            translation: "許可された本の訳語",
            category: "書籍",
            origin: "XML取込",
            updatedAt: "2026-04-12 00:00"
          }
        ],
        selectedEntry: createImportedBookSelectedEntry(),
        selectedId: "201",
        totalCount: 1,
        query: "",
        category: "すべて",
        hasStagedFile: true,
        selectedFileName: "master.xml",
        selectedFileReference: "master.xml",
        importStage: "done",
        importProgress: 100,
        importStatusValue: "完了",
        importSummary: {
          fileName: "master.xml",
          importedCount: 1,
          updatedCount: 0,
          totalCount: 1,
          selectedSource: "Allowed Book Source"
        },
        listHeadline: "1 件のエントリを表示しています。",
        selectionStatusText: "Allowed Book Source を選択中",
        detailSublineText: "XML取込 / 最終更新 2026-04-12 00:00"
      })
    )

    // Assert
    await waitFor(() => {
      expect(screen.getByRole("combobox", { name: "カテゴリ" })).toHaveValue(
        "すべて"
      )
    })
  })

  test("完了 view model push で detail title を更新する", async () => {
    // Arrange
    const controller = renderApp()

    // Act
    controller.pushViewModel(
      buildMasterDictionaryScreenViewModel({
        entries: [
          {
            id: "201",
            source: "Allowed Book Source",
            translation: "許可された本の訳語",
            category: "書籍",
            origin: "XML取込",
            updatedAt: "2026-04-12 00:00"
          }
        ],
        selectedEntry: createImportedBookSelectedEntry(),
        selectedId: "201",
        totalCount: 1,
        query: "",
        category: "すべて",
        hasStagedFile: true,
        selectedFileName: "master.xml",
        selectedFileReference: "master.xml",
        importStage: "done",
        importProgress: 100,
        importStatusValue: "完了",
        importSummary: {
          fileName: "master.xml",
          importedCount: 1,
          updatedCount: 0,
          totalCount: 1,
          selectedSource: "Allowed Book Source"
        },
        listHeadline: "1 件のエントリを表示しています。",
        selectionStatusText: "Allowed Book Source を選択中",
        detailSublineText: "XML取込 / 最終更新 2026-04-12 00:00"
      })
    )

    // Assert
    await waitFor(() => {
      expect(document.querySelector("#detailTitle")).toHaveTextContent(
        "Allowed Book Source"
      )
    })
  })
})
