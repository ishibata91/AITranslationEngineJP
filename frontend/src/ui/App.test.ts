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

function createXmlFile(contents: string, name = "master-dictionary.xml"): File {
  return new File([contents], name, { type: "text/xml" })
}

describe("App dashboard shell", () => {
  beforeEach(() => {
    window.history.replaceState(null, "", "#")
  })

  test("SCN-DAS-001: 起動時にダッシュボードを既定表示する", () => {
    renderApp()

    expect(
      screen.getByRole("heading", { name: "ダッシュボード" })
    ).toBeInTheDocument()
    expect(window.location.hash).toBe("#dashboard")
  })

  test("invalid hash は dashboard に正規化される", () => {
    window.history.replaceState(null, "", "#not-approved-route")

    renderApp()

    expect(
      screen.getByRole("heading", { name: "ダッシュボード" })
    ).toBeInTheDocument()
    expect(window.location.hash).toBe("#dashboard")
  })

  test("SCN-DAS-002/003: グローバルナビゲーションと入口カードから承認済みルートへ移動できる", async () => {
    const user = userEvent.setup()
    renderApp()

    const globalNavigation = screen.getByRole("navigation", {
      name: "グローバルナビゲーション"
    })

    expect(within(globalNavigation).getAllByRole("link")).toHaveLength(
      DASHBOARD_SHELL_PRIMARY_ROUTES.length
    )

    for (const route of DASHBOARD_ENTRY_ROUTES) {
      await user.click(
        within(globalNavigation).getByRole("link", { name: route.label })
      )
      expect(
        screen.getByRole("heading", { level: 1, name: route.label })
      ).toBeInTheDocument()
      expect(window.location.hash).toBe(`#${route.id}`)

      await user.click(
        within(globalNavigation).getByRole("link", { name: "ダッシュボード" })
      )

      const dashboardCardSection = screen
        .getByRole("heading", { name: "作業を選ぶ" })
        .closest("section")
      if (!dashboardCardSection) {
        throw new Error("ダッシュボード入口カードのセクションが見つかりません")
      }

      const cardHeading = within(dashboardCardSection).getByRole("heading", {
        level: 3,
        name: route.label
      })
      const cardLink = cardHeading.closest("a")
      if (!cardLink) {
        throw new Error(`入口カードのリンクが見つかりません: ${route.label}`)
      }

      await user.click(cardLink)
      expect(
        screen.getByRole("heading", { level: 1, name: route.label })
      ).toBeInTheDocument()
      expect(window.location.hash).toBe(`#${route.id}`)
    }
  })

  test("SCN-DAS-004/005: プレースホルダー画面でも共通シェルを保持して再移動できる", async () => {
    const user = userEvent.setup()
    renderApp()

    const globalNavigation = screen.getByRole("navigation", {
      name: "グローバルナビゲーション"
    })

    await user.click(
      within(globalNavigation).getByRole("link", { name: "翻訳管理" })
    )
    expect(
      screen.getByRole("heading", { level: 1, name: "翻訳管理" })
    ).toBeInTheDocument()
    expect(screen.getByText(PLACEHOLDER_LEAD)).toBeInTheDocument()

    await user.click(
      within(globalNavigation).getByRole("link", { name: "出力管理" })
    )
    expect(
      screen.getByRole("heading", { level: 1, name: "出力管理" })
    ).toBeInTheDocument()
  })
})

describe("App master dictionary screen", () => {
  beforeEach(() => {
    window.history.replaceState(null, "", "#master-dictionary")
  })

  test("contract factory から受け取った view model を描画し mount/dispose を呼ぶ", async () => {
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

    const view = render(App, {
      props: {
        createMasterDictionaryScreenController: () => controller
      }
    })

    await waitFor(() => {
      expect(controller.mount).toHaveBeenCalledTimes(1)
    })
    expect(
      screen.getByRole("heading", { level: 3, name: "辞書一覧" })
    ).toBeInTheDocument()
    expect(document.querySelector("#detailTitle")).toHaveTextContent(
      "Dragon Priest"
    )
    expect(document.querySelector("#importStatusValue")).toHaveTextContent(
      "完了"
    )
    expect(document.querySelector("#importResultSelection")).toHaveTextContent(
      "Dragon Priest"
    )

    view.unmount()

    expect(controller.dispose).toHaveBeenCalledTimes(1)
  })

  test("一覧と主要操作は controller contract を呼び出す", async () => {
    const user = userEvent.setup()
    const controller = renderApp()

    await waitFor(() => {
      expect(controller.mount).toHaveBeenCalledTimes(1)
    })

    await user.click(screen.getByRole("button", { name: /ドラゴン・プリースト/ }))
    await user.click(screen.getByRole("button", { name: "新規登録" }))
    await user.click(screen.getByRole("button", { name: "更新" }))
    await user.click(screen.getByRole("button", { name: "削除" }))

    expect(controller.selectRow).toHaveBeenCalledWith("101")
    expect(controller.openCreateModal).toHaveBeenCalledTimes(1)
    expect(controller.openEditModal).toHaveBeenCalledTimes(1)
    expect(controller.openDeleteModal).toHaveBeenCalledTimes(1)
  })

  test("検索、カテゴリ変更、ページ移動は controller contract を通す", async () => {
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

    await user.type(screen.getByRole("searchbox", { name: "検索" }), "Reach")
    await user.selectOptions(
      screen.getByRole("combobox", { name: "カテゴリ" }),
      "地名"
    )
    await user.click(screen.getByRole("button", { name: "次の30件" }))

    expect(controller.handleSearchInput).toHaveBeenCalled()
    expect(controller.handleCategoryChange).toHaveBeenCalled()
    expect(controller.goToNextPage).toHaveBeenCalledTimes(1)
  })

  test("XML 選択と選び直しは controller contract を通す", async () => {
    const user = userEvent.setup()
    const controller = renderApp()
    const xmlInput = document.querySelector("#xmlFileInput")
    if (!(xmlInput instanceof HTMLInputElement)) {
      throw new Error("xmlFileInput が見つかりません")
    }

    const xmlFile = createXmlFile("<Root />")
    await user.upload(xmlInput, xmlFile)

    expect(controller.stageXmlImport).toHaveBeenCalledWith(xmlFile)

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

    await user.click(screen.getByRole("button", { name: "選び直す" }))

    expect(controller.resetImportSelection).toHaveBeenCalledTimes(1)
  })

  test("import 状態の表示更新は controller の view model push だけで反映する", async () => {
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

    await user.click(screen.getByRole("button", { name: "この XML を取り込む" }))

    expect(controller.startImport).toHaveBeenCalledTimes(1)

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

    await waitFor(() => {
      expect(document.querySelector("#importStatusValue")).toHaveTextContent(
        "取込中"
      )
      expect(document.querySelector("#importProgressFill")).toHaveAttribute(
        "style",
        "width: 78%;"
      )
    })

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

    await waitFor(() => {
      expect(screen.getByText("完了")).toBeInTheDocument()
      expect(screen.getByRole("searchbox", { name: "検索" })).toHaveValue("")
      expect(screen.getByRole("combobox", { name: "カテゴリ" })).toHaveValue(
        "すべて"
      )
      expect(document.querySelector("#detailTitle")).toHaveTextContent(
        "Allowed Book Source"
      )
    })
  })
})
