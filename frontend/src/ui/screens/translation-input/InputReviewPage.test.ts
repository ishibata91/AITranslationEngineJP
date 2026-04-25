import { render, screen, within } from "@testing-library/svelte"
import userEvent from "@testing-library/user-event"
import { describe, expect, test, vi } from "vitest"

import type {
  TranslationInputReviewItem,
  TranslationInputScreenViewModel
} from "@application/gateway-contract/translation-input"
import type {
  TranslationInputScreenControllerContract
} from "@application/contract/translation-input"
import type { TranslationInputScreenViewModelListener } from "@application/contract/translation-input/translation-input-screen-contract"
import InputReviewPage from "@ui/screens/translation-input/InputReviewPage.svelte"

function createItem(
  overrides: Partial<TranslationInputReviewItem> = {}
): TranslationInputReviewItem {
  return {
    localId: "input-41",
    inputId: 41,
    fileName: "input-review.json",
    filePath: "/mods/input-review.json",
    fileHash: "hash-41",
    importTimestamp: "invalid-timestamp",
    status: "registered",
    accepted: true,
    canRebuild: true,
    lastAction: "import",
    errorKind: null,
    warnings: [],
    summary: {
      input: {
        id: 41,
        sourceFilePath: "/mods/input-review.json",
        sourceTool: "xEdit",
        targetPluginName: "Skyrim.esm",
        targetPluginType: "esm",
        recordCount: 5,
        importedAt: "2026-04-26T09:30:00Z"
      },
      translationRecordCount: 5,
      translationFieldCount: 8,
      categories: [
        {
          category: "NPC",
          recordCount: 2,
          fieldCount: 3
        },
        {
          category: "BOOK",
          recordCount: 3,
          fieldCount: 5
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
  const items = overrides.items ?? [createItem()]
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
    latestOutcomeTitle: "結果: 登録済み",
    latestOutcomeText:
      "翻訳レコード件数、カテゴリ別件数、sample field を確認できます。",
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
  private viewModel: TranslationInputScreenViewModel

  private readonly listeners = new Set<TranslationInputScreenViewModelListener>()

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
    this.listeners.add(listener)
    return () => {
      this.listeners.delete(listener)
    }
  }

  getViewModel(): TranslationInputScreenViewModel {
    return this.viewModel
  }
}

describe("InputReviewPage", () => {
  test("一覧、概要、カテゴリ別件数、sample field を表示し、禁止 action を出さない", () => {
    const controller = new TranslationInputScreenControllerFake()

    render(InputReviewPage, {
      props: {
        createController: () => controller
      }
    })

    expect(screen.getByRole("heading", { level: 2, name: "Input Review" })).toBeInTheDocument()
    expect(screen.getAllByText("input-review.json").length).toBeGreaterThan(0)
    expect(screen.getAllByText("/mods/input-review.json").length).toBeGreaterThan(0)
    expect(screen.getAllByText("hash-41").length).toBeGreaterThan(0)
    expect(screen.getAllByText("invalid-timestamp")).toHaveLength(2)
    expect(screen.getByText("accepted")).toBeInTheDocument()
    expect(screen.getByText("rebuild 可")).toBeInTheDocument()
    expect(screen.getAllByText("5").length).toBeGreaterThan(0)
    expect(screen.getByText("8")).toBeInTheDocument()
    expect(screen.getByText("Skyrim.esm")).toBeInTheDocument()
    expect(screen.getByText("xEdit")).toBeInTheDocument()
    expect(screen.getByText("record 2 / field 3")).toBeInTheDocument()
    expect(screen.getByText("record 3 / field 5")).toBeInTheDocument()
    expect(screen.getByText("NPC_:FULL")).toBeInTheDocument()
    expect(screen.getByText("Hello there")).toBeInTheDocument()
    expect(screen.getByText("00012345")).toBeInTheDocument()
    expect(screen.getByText("SampleNPC")).toBeInTheDocument()
    expect(screen.queryByRole("button", { name: "ジョブ作成" })).not.toBeInTheDocument()
    expect(screen.queryByRole("button", { name: "翻訳開始" })).not.toBeInTheDocument()
    expect(screen.queryByRole("button", { name: "出力生成" })).not.toBeInTheDocument()
  })

  test("error kind と warning kind を区別して表示する", () => {
    const warningItem = createItem({
      localId: "warning",
      fileName: "warning.json",
      filePath: "/mods/warning.json",
      status: "warning",
      warnings: [
        {
          kind: "unknown_field_definition",
          recordType: "BOOK",
          subrecordType: "DESC",
          message: "unknown description field"
        }
      ],
      summary: {
        ...createItem().summary!,
        input: {
          ...createItem().summary!.input,
          id: 51,
          sourceFilePath: "/mods/warning.json"
        },
        warnings: [
          {
            kind: "unknown_field_definition",
            recordType: "BOOK",
            subrecordType: "DESC",
            message: "unknown description field"
          }
        ]
      }
    })
    const items = [
      warningItem,
      createItem({
        localId: "duplicate",
        fileName: "duplicate.json",
        errorKind: "duplicate_input_hash",
        accepted: false,
        canRebuild: false,
        status: "failed",
        summary: null
      }),
      createItem({
        localId: "invalid",
        fileName: "invalid.json",
        errorKind: "invalid_json",
        accepted: false,
        canRebuild: false,
        status: "failed",
        summary: null
      }),
      createItem({
        localId: "shape",
        fileName: "shape.json",
        errorKind: "unsupported_extract_shape",
        accepted: false,
        canRebuild: false,
        status: "failed",
        summary: null
      }),
      createItem({
        localId: "missing-field",
        fileName: "missing-field.json",
        errorKind: "missing_required_field",
        accepted: false,
        canRebuild: false,
        status: "failed",
        summary: null
      }),
      createItem({
        localId: "source-missing",
        fileName: "source-missing.json",
        errorKind: "source_file_missing",
        accepted: false,
        status: "rebuild-required"
      }),
      createItem({
        localId: "cache-missing",
        fileName: "cache-missing.json",
        errorKind: "cache_missing",
        accepted: false,
        canRebuild: false,
        status: "failed",
        summary: null
      })
    ]
    const controller = new TranslationInputScreenControllerFake(
      createViewModel({
        items,
        selectedItemId: "warning",
        selectedItem: warningItem,
        canRebuildSelected: true,
        latestOutcomeTitle: "結果: unknown field definition を含む登録済み",
        latestOutcomeText: "unknown field definition"
      })
    )

    render(InputReviewPage, {
      props: {
        createController: () => controller
      }
    })

    expect(screen.getByText(/重複 input/)).toBeInTheDocument()
    expect(screen.getByText(/invalid JSON/)).toBeInTheDocument()
    expect(screen.getByText(/non-xEdit JSON/)).toBeInTheDocument()
    expect(screen.getByText(/missing required field/)).toBeInTheDocument()
    expect(screen.getByText(/source file missing/)).toBeInTheDocument()
    expect(screen.getByText(/cache missing/)).toBeInTheDocument()
    expect(screen.getAllByText("unknown field definition").length).toBeGreaterThan(0)
    expect(screen.getByText("再構築が必要")).toBeInTheDocument()
  })

  test("JSON upload、選び直し、登録、再構築、一覧選択を controller へ委譲する", async () => {
    const user = userEvent.setup()
    const selectedItem = createItem()
    const controller = new TranslationInputScreenControllerFake(
      createViewModel({
        items: [selectedItem],
        selectedItem,
        selectedItemId: selectedItem.localId,
        stagedFile: {
          fileName: "input-review.json",
          filePath: "/mods/input-review.json",
          fileHash: "hash-41"
        },
        hasStagedFile: true,
        canImport: true,
        stagedFileName: "input-review.json",
        stagedFilePath: "/mods/input-review.json",
        stagedFileHash: "hash-41",
        canRebuildSelected: true
      })
    )
    const { container, unmount } = render(InputReviewPage, {
      props: {
        createController: () => controller
      }
    })
    const input = container.querySelector("#translationInputFile")
    const file = new File(["{}"], "uploaded.json", {
      type: "application/json"
    })

    if (!(input instanceof HTMLInputElement)) {
      throw new Error("translation input file element not found")
    }

    Object.defineProperty(file, "path", {
      value: "/mods/uploaded.json"
    })

    await user.upload(input, file)
    await user.click(screen.getByRole("button", { name: "選び直す" }))
    await user.click(screen.getByRole("button", { name: "この JSON を登録" }))
    await user.click(screen.getByRole("button", { name: "cache を再構築" }))

    const list = screen.getByRole("list")
    await user.click(within(list).getByRole("button", { name: /input-review.json/ }))

    expect(controller.mount).toHaveBeenCalledTimes(1)
    expect(controller.stageJsonImport).toHaveBeenCalledWith(file)
    expect(controller.resetImportSelection).toHaveBeenCalledTimes(1)
    expect(controller.startImport).toHaveBeenCalledTimes(1)
    expect(controller.rebuildSelected).toHaveBeenCalledTimes(1)
    expect(controller.selectItem).toHaveBeenCalledWith("input-41")

    unmount()

    expect(controller.dispose).toHaveBeenCalledTimes(1)
  })
})