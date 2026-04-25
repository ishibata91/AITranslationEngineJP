import { render, screen, waitFor } from "@testing-library/svelte"
import userEvent from "@testing-library/user-event"
import { describe, expect, test, vi } from "vitest"

import type { MasterPersonaScreenControllerContract } from "@application/contract/master-persona/master-persona-screen-contract"
import type {
  TranslationInputScreenControllerContract
} from "@application/contract/translation-input"
import type { TranslationInputScreenViewModelListener } from "@application/contract/translation-input/translation-input-screen-contract"
import type {
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

class TranslationInputScreenControllerFake
  implements TranslationInputScreenControllerContract
{
  private readonly viewModel: TranslationInputScreenViewModel

  readonly mount = vi.fn(async () => {})
  readonly dispose = vi.fn(() => {})
  readonly selectItem = vi.fn(() => {})
  readonly stageJsonImport = vi.fn(async () => {})
  readonly resetImportSelection = vi.fn(() => {})
  readonly startImport = vi.fn(async () => {})
  readonly rebuildSelected = vi.fn(async () => {})

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

describe("App translation-input route", () => {
  test("translation-management route で Input Review page を描画する", async () => {
    window.history.replaceState(null, "", "#translation-management")

    const controller = new TranslationInputScreenControllerFake()
    const masterPersonaStub = (() =>
      ({}) as MasterPersonaScreenControllerContract) as () => MasterPersonaScreenControllerContract

    render(App, {
      props: {
        createMasterPersonaScreenController: masterPersonaStub,
        createTranslationInputScreenController: () => controller
      }
    })

    expect(screen.getByRole("heading", { level: 1, name: "翻訳管理" })).toBeInTheDocument()
    expect(screen.getByRole("heading", { level: 2, name: "Input Review" })).toBeInTheDocument()

    await waitFor(() => {
      expect(controller.mount).toHaveBeenCalledTimes(1)
    })
  })

  test("route を離れて戻っても一覧、選択、rebuild 入口を維持する", async () => {
    window.history.replaceState(null, "", "#translation-management")

    const user = userEvent.setup()
    const selectedItem = createItem()
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
        createTranslationInputScreenController
      }
    })

    await waitFor(() => {
      expect(controller.mount).toHaveBeenCalledTimes(1)
    })

    expect(
      screen.getByRole("button", { name: /kept-input.json/ })
    ).toBeInTheDocument()
    expect(screen.getByText("kept-input.json / 登録済み")).toBeInTheDocument()
    expect(
      screen.getByRole("button", { name: "cache を再構築" })
    ).toBeEnabled()

    await user.click(screen.getByRole("link", { name: "ダッシュボード" }))

    expect(
      screen.queryByRole("heading", { level: 2, name: "Input Review" })
    ).not.toBeInTheDocument()
    expect(controller.dispose).toHaveBeenCalledTimes(1)

    await user.click(screen.getByRole("link", { name: "翻訳管理" }))

    expect(screen.getByRole("heading", { level: 2, name: "Input Review" })).toBeInTheDocument()
    expect(
      screen.getByRole("button", { name: /kept-input.json/ })
    ).toBeInTheDocument()
    expect(screen.getByText("kept-input.json / 登録済み")).toBeInTheDocument()
    expect(
      screen.getByRole("button", { name: "cache を再構築" })
    ).toBeEnabled()
    expect(createTranslationInputScreenController).toHaveBeenCalledTimes(2)
    expect(controller.mount).toHaveBeenCalledTimes(2)

    unmount()

    expect(controller.dispose).toHaveBeenCalledTimes(2)
  })
})