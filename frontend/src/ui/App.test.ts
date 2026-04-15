import { render, screen, waitFor, within } from "@testing-library/svelte"
import userEvent from "@testing-library/user-event"

import type {
  MasterDictionaryScreenControllerContract,
  MasterDictionaryScreenViewModelListener
} from "@application/contract/master-dictionary/master-dictionary-screen-contract"
import type { MasterDictionaryScreenViewModel } from "@application/contract/master-dictionary/master-dictionary-screen-types"
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
    updatedAt: "2026-01-01 00:00",
    note: "REC: NPC_:FULL / EDID: SeedDragonPriest"
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

class MasterDictionaryScreenControllerFake
  implements MasterDictionaryScreenControllerContract
{
  private viewModel: MasterDictionaryScreenViewModel

  private readonly listeners = new Set<MasterDictionaryScreenViewModelListener>()

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

function renderApp(
  controller = new MasterDictionaryScreenControllerFake()
): MasterDictionaryScreenControllerFake {
  render(App, {
    props: {
      createMasterDictionaryScreenController: () => controller
    }
  })

  return controller
}

function renderAppView(
  controller = new MasterDictionaryScreenControllerFake()
): {
  controller: MasterDictionaryScreenControllerFake
  unmount: () => void
} {
  const view = render(App, {
    props: {
      createMasterDictionaryScreenController: () => controller
    }
  })

  return {
    controller,
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
    const dashboardHeading = screen.getByRole("heading", { name: "ダッシュボード" })

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
    expect(screen.getByRole("heading", { name: "ダッシュボード" })).toBeInTheDocument()
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
      await user.click(within(getGlobalNavigation()).getByRole("link", { name: label }))

      // Assert
      expect(screen.getByRole("heading", { level: 1, name: label })).toBeInTheDocument()
    }
  )

  test.each(DASHBOARD_ENTRY_ROUTES)(
    "SCN-DAS-002: グローバルナビゲーションから $id hash へ移動できる",
    async ({ id, label }) => {
      // Arrange
      const user = userEvent.setup()
      renderApp()

      // Act
      await user.click(within(getGlobalNavigation()).getByRole("link", { name: label }))

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
      expect(screen.getByRole("heading", { level: 1, name: label })).toBeInTheDocument()
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

  test("SCN-DAS-004: プレースホルダー画面でも共通 lead を表示する", async () => {
    // Arrange
    const user = userEvent.setup()
    renderApp()

    // Act
    await user.click(
      within(getGlobalNavigation()).getByRole("link", { name: "翻訳管理" })
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
    expect(screen.getByRole("heading", { level: 3, name: "辞書一覧" })).toBeInTheDocument()
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
    const importResultSelection = document.querySelector("#importResultSelection")

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
    await user.click(screen.getByRole("button", { name: /ドラゴン・プリースト/ }))

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
    await user.click(screen.getByRole("button", { name: "この XML を取り込む" }))

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
      expect(document.querySelector("#importStatusValue")).toHaveTextContent("取込中")
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
        selectedEntry: {
          id: "201",
          source: "Allowed Book Source",
          translation: "許可された本の訳語",
          category: "書籍",
          origin: "XML取込",
          updatedAt: "2026-04-12 00:00",
          note: "REC: BOOK:FULL / EDID: ImportBook"
        },
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
        selectedEntry: {
          id: "201",
          source: "Allowed Book Source",
          translation: "許可された本の訳語",
          category: "書籍",
          origin: "XML取込",
          updatedAt: "2026-04-12 00:00",
          note: "REC: BOOK:FULL / EDID: ImportBook"
        },
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
      expect(screen.getByText(/件数は保存済みエントリ単位で集計しています。/)).toBeInTheDocument()
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
        selectedEntry: {
          id: "201",
          source: "Allowed Book Source",
          translation: "許可された本の訳語",
          category: "書籍",
          origin: "XML取込",
          updatedAt: "2026-04-12 00:00",
          note: "REC: BOOK:FULL / EDID: ImportBook"
        },
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
        selectedEntry: {
          id: "201",
          source: "Allowed Book Source",
          translation: "許可された本の訳語",
          category: "書籍",
          origin: "XML取込",
          updatedAt: "2026-04-12 00:00",
          note: "REC: BOOK:FULL / EDID: ImportBook"
        },
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
      expect(screen.getByRole("combobox", { name: "カテゴリ" })).toHaveValue("すべて")
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
        selectedEntry: {
          id: "201",
          source: "Allowed Book Source",
          translation: "許可された本の訳語",
          category: "書籍",
          origin: "XML取込",
          updatedAt: "2026-04-12 00:00",
          note: "REC: BOOK:FULL / EDID: ImportBook"
        },
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
