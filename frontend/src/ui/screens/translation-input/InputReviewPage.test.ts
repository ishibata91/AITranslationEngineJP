import { render, screen, within } from "@testing-library/svelte"
import userEvent from "@testing-library/user-event"
import { describe, expect, test, vi } from "vitest"

import type {
  TranslationInputReviewItem,
  TranslationInputScreenViewModel
} from "@application/gateway-contract/translation-input"
import type { TranslationInputScreenControllerContract } from "@application/contract/translation-input"
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

class TranslationInputScreenControllerFake implements TranslationInputScreenControllerContract {
  private viewModel: TranslationInputScreenViewModel

  private readonly listeners =
    new Set<TranslationInputScreenViewModelListener>()

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

    expect(
      screen.getByRole("heading", { level: 2, name: "データロード" })
    ).toBeInTheDocument()
    expect(
      screen.queryByRole("heading", { level: 2, name: "Input Review" })
    ).not.toBeInTheDocument()
    expect(screen.getAllByText("input-review.json").length).toBeGreaterThan(0)
    expect(
      screen.getAllByText("/mods/input-review.json").length
    ).toBeGreaterThan(0)
    expect(screen.getAllByText("hash-41").length).toBeGreaterThan(0)
    expect(screen.getByText("登録結果")).toBeInTheDocument()
    expect(screen.getByText("読み込み日時")).toBeInTheDocument()
    expect(screen.getByText("問題区分")).toBeInTheDocument()
    expect(screen.getByText("選択 JSON")).toBeInTheDocument()
    expect(screen.getByText("ロード対象を選ぶ")).toBeInTheDocument()
    expect(screen.getByText("この JSON を登録")).toBeInTheDocument()
    expect(screen.getByText("選び直す")).toBeInTheDocument()
    expect(screen.getAllByText("invalid-timestamp")).toHaveLength(1)
    expect(screen.getAllByText("登録済み").length).toBeGreaterThan(0)
    expect(screen.queryByText("再構築")).not.toBeInTheDocument()
    expect(document.querySelector(".detail-panel")).toBeNull()
    expect(
      screen.queryByTestId("translation-input-review-selected-input-region")
    ).not.toBeInTheDocument()
    expect(screen.queryByText("1 件")).not.toBeInTheDocument()
    expect(document.querySelector(".count-pill")).toBeNull()
    expect(
      screen.queryByRole("button", { name: "ジョブ作成" })
    ).not.toBeInTheDocument()
    expect(
      screen.queryByRole("button", { name: "翻訳開始" })
    ).not.toBeInTheDocument()
    expect(
      screen.queryByRole("button", { name: "出力生成" })
    ).not.toBeInTheDocument()
  })

  test("登録成功後は単語翻訳へ進む footer を表示し、クリックで callback を呼ぶ", async () => {
    const user = userEvent.setup()
    const onOpenJobRun = vi.fn()
    const controller = new TranslationInputScreenControllerFake()

    render(InputReviewPage, {
      props: {
        createController: () => controller,
        onOpenJobRun
      }
    })

    expect(
      screen.getByText("入力データを選択済みです。次に単語翻訳へ進みます。")
    ).toBeInTheDocument()

    const footer = screen.getByRole("region", {
      name: "次の作業"
    })

    await user.click(
      within(footer).getByRole("button", { name: "単語翻訳へ進む" })
    )

    expect(onOpenJobRun).toHaveBeenCalledTimes(1)
  })

  test("failed item では単語翻訳へ進む button を表示しない", () => {
    const failedItem = createItem({
      localId: "failed-item",
      status: "failed",
      accepted: false,
      canRebuild: false,
      errorKind: "invalid_json",
      summary: null
    })
    const controller = new TranslationInputScreenControllerFake(
      createViewModel({
        items: [failedItem],
        selectedItemId: failedItem.localId,
        selectedItem: failedItem,
        canRebuildSelected: false,
        latestOutcomeTitle: "結果: invalid JSON",
        latestOutcomeText:
          "error kind を保持したまま、再試行または別ファイル選択へ戻れます。"
      })
    )

    render(InputReviewPage, {
      props: {
        createController: () => controller
      }
    })

    expect(
      screen.queryByRole("button", { name: "単語翻訳へ進む" })
    ).not.toBeInTheDocument()
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
    expect(screen.getByText("警告あり")).toBeInTheDocument()
    expect(screen.queryByText("unknown field definition")).not.toBeInTheDocument()
    expect(screen.queryByText("再構築が必要")).not.toBeInTheDocument()
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

    await user.upload(input, file)
    await user.click(screen.getByRole("button", { name: "選び直す" }))
    await user.click(screen.getByRole("button", { name: "この JSON を登録" }))

    const list = screen.getByRole("list")
    await user.click(
      within(list).getByRole("button", { name: /input-review.json/ })
    )

    expect(controller.mount).toHaveBeenCalledTimes(1)
    expect(controller.stageJsonImport).toHaveBeenCalledWith(file)
    const uploadedFile = (
      controller.stageJsonImport as unknown as {
        mock: { calls: unknown[][] }
      }
    ).mock.calls[0]?.[0]

    if (!(uploadedFile instanceof File)) {
      throw new Error("uploaded file was not passed to controller")
    }

    expect(uploadedFile.name).toBe("uploaded.json")
    expect((uploadedFile as File & { path?: string }).path).toBeUndefined()
    expect(controller.resetImportSelection).toHaveBeenCalledTimes(1)
    expect(controller.startImport).toHaveBeenCalledTimes(1)
    expect(controller.rebuildSelected).not.toHaveBeenCalled()
    expect(controller.selectItem).toHaveBeenCalledWith("input-41")

    unmount()

    expect(controller.dispose).toHaveBeenCalledTimes(1)
  })
  test("一覧の表示項目を維持し、詳細パネル依存なしで UI が破綻しない", () => {
    const item = createItem({
      warnings: [],
      summary: {
        ...createItem().summary!,
        categories: [],
        sampleFields: [],
        warnings: []
      }
    })
    const controller = new TranslationInputScreenControllerFake(
      createViewModel({
        items: [item],
        selectedItemId: item.localId,
        selectedItem: item,
        canRebuildSelected: true,
        latestOutcomeTitle: "結果: 登録済み",
        latestOutcomeText:
          "翻訳レコード件数、カテゴリ別件数、sample field を確認できます。"
      })
    )

    render(InputReviewPage, {
      props: {
        createController: () => controller
      }
    })

    expect(screen.getAllByText("読み込み済みデータ").length).toBeGreaterThan(0)
    expect(screen.getByText("問題区分")).toBeInTheDocument()
    expect(screen.getByText("-")).toBeInTheDocument()
    expect(
      screen.queryByTestId("translation-input-review-selected-input-region")
    ).not.toBeInTheDocument()
    expect(document.querySelector(".detail-panel")).toBeNull()
  })
})
