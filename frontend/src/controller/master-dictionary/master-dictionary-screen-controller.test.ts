import { describe, expect, test, vi } from "vitest"

import type {
  MasterDictionaryScreenState,
  MasterDictionaryScreenViewModel
} from "@application/contract/master-dictionary/master-dictionary-screen-types"

import { MasterDictionaryScreenController } from "./master-dictionary-screen-controller"

function createSelectedEntry() {
  return {
    id: "101",
    source: "Dragon Priest",
    translation: "ドラゴン・プリースト",
    category: "地名",
    origin: "確認待ち",
    rec: "NPC_:FULL",
    edid: "SeedDragonPriest",
    updatedAt: "2026-01-01 00:00",
    note: "note"
  }
}

function createState(
  overrides: Partial<MasterDictionaryScreenState> = {}
): MasterDictionaryScreenState {
  return {
    entries: [],
    selectedEntry: null,
    selectedId: null,
    totalCount: 0,
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
    ...overrides
  }
}

function createViewModel(
  state: MasterDictionaryScreenState
): MasterDictionaryScreenViewModel {
  return {
    ...state,
    gatewayStatus: "接続済み",
    hasStagedFile: state.selectedFileReference !== null,
    isImportRunning: state.importStage === "running",
    importStatusValue: state.importStage,
    importStatusText: state.importStage,
    categoryOptions: ["すべて", "固有名詞", "地名"],
    totalPages: Math.max(1, Math.ceil(state.totalCount / 30)),
    pageStatusText: `${state.page + 1} ページ`,
    listHeadline: `${state.totalCount} 件`,
    selectionStatusText: state.selectedId ?? "none",
    detailSublineText: state.selectedEntry?.updatedAt ?? "none"
  }
}

function createControllerHarness(
  initialState: MasterDictionaryScreenState = createState()
) {
  let state = initialState
  const listeners = new Set<(state: MasterDictionaryScreenState) => void>()

  const store = {
    subscribe: vi.fn((listener: (nextState: MasterDictionaryScreenState) => void) => {
      listeners.add(listener)
      return () => {
        listeners.delete(listener)
      }
    }),
    snapshot: vi.fn(() => state),
    update: vi.fn((mutator: (draft: MasterDictionaryScreenState) => void) => {
      const draft = structuredClone(state)
      mutator(draft)
      state = draft
      for (const listener of listeners) {
        listener(state)
      }
    })
  }

  const presenter = {
    toViewModel: vi.fn(
      (nextState: MasterDictionaryScreenState, isGatewayConnected: boolean) => ({
        ...createViewModel(nextState),
        gatewayStatus: isGatewayConnected ? "接続済み" : "未接続"
      })
    )
  }

  const useCase = {
    loadEntries: vi.fn(async () => {}),
    selectEntry: vi.fn(async () => {}),
    saveCurrentEntry: vi.fn(async () => {}),
    deleteCurrentEntry: vi.fn(async () => {}),
    startStagedXmlImport: vi.fn(async () => {})
  }

  const runtimeEventAdapter = {
    subscribe: vi.fn(() => true),
    detach: vi.fn(() => {})
  }

  const controller = new MasterDictionaryScreenController({
    isGatewayConnected: true,
    store,
    presenter,
    useCase,
    runtimeEventAdapter
  })

  return {
    controller,
    store,
    presenter,
    useCase,
    runtimeEventAdapter,
    getState: () => state
  }
}

describe("MasterDictionaryScreenController", () => {
  test("mount は runtimeEventAdapter.subscribe を呼ぶ", async () => {
    // Arrange
    const harness = createControllerHarness()

    // Act
    await harness.controller.mount()

    // Assert
    expect(harness.runtimeEventAdapter.subscribe).toHaveBeenCalledTimes(1)
  })

  test("mount は useCase.loadEntries を呼ぶ", async () => {
    // Arrange
    const harness = createControllerHarness()

    // Act
    await harness.controller.mount()

    // Assert
    expect(harness.useCase.loadEntries).toHaveBeenCalledTimes(1)
  })

  test("subscribe は store 更新を listener に中継する", () => {
    // Arrange
    const harness = createControllerHarness(createState({ selectedId: "101" }))
    const listener = vi.fn()

    harness.controller.subscribe(listener)

    // Act
    harness.store.update((draft) => {
      draft.selectedId = "102"
    })

    // Assert
    expect(listener).toHaveBeenCalled()
  })

  test("subscribe は presenter.toViewModel を通した値を listener へ渡す", () => {
    // Arrange
    const harness = createControllerHarness(createState({ selectedId: "101" }))
    const listener = vi.fn()

    harness.controller.subscribe(listener)

    // Act
    harness.store.update((draft) => {
      draft.selectedId = "102"
    })

    // Assert
    expect(harness.presenter.toViewModel).toHaveBeenCalled()
  })

  test("getViewModel は最新 state を view model へ変換する", () => {
    // Arrange
    const harness = createControllerHarness(createState({ selectedId: "101" }))

    harness.controller.subscribe(() => {})
    harness.store.update((draft) => {
      draft.selectedId = "102"
    })

    // Act
    const viewModel = harness.controller.getViewModel()

    // Assert
    expect(viewModel.selectionStatusText).toBe("102")
  })

  test("dispose は runtimeEventAdapter.detach を呼ぶ", () => {
    // Arrange
    const harness = createControllerHarness()

    // Act
    harness.controller.dispose()

    // Assert
    expect(harness.runtimeEventAdapter.detach).toHaveBeenCalledTimes(1)
  })

  test("openCreateModal は create modal へ切り替える", () => {
    // Arrange
    const harness = createControllerHarness(
      createState({
        modalState: "edit",
        formSource: "old",
        formCategory: "地名",
        formOrigin: "確認待ち",
        formTranslation: "old",
        errorMessage: "error"
      })
    )

    // Act
    harness.controller.openCreateModal()

    // Assert
    expect(harness.getState().modalState).toBe("create")
  })

  test("openCreateModal は source を空文字へ初期化する", () => {
    // Arrange
    const harness = createControllerHarness(createState({ formSource: "old" }))

    // Act
    harness.controller.openCreateModal()

    // Assert
    expect(harness.getState().formSource).toBe("")
  })

  test("openCreateModal は category を既定値へ初期化する", () => {
    // Arrange
    const harness = createControllerHarness(createState({ formCategory: "地名" }))

    // Act
    harness.controller.openCreateModal()

    // Assert
    expect(harness.getState().formCategory).toBe("固有名詞")
  })

  test("openCreateModal は origin を既定値へ初期化する", () => {
    // Arrange
    const harness = createControllerHarness(createState({ formOrigin: "確認待ち" }))

    // Act
    harness.controller.openCreateModal()

    // Assert
    expect(harness.getState().formOrigin).toBe("手動登録")
  })

  test("openCreateModal は translation を空文字へ初期化する", () => {
    // Arrange
    const harness = createControllerHarness(createState({ formTranslation: "old" }))

    // Act
    harness.controller.openCreateModal()

    // Assert
    expect(harness.getState().formTranslation).toBe("")
  })

  test("openCreateModal は errorMessage をクリアする", () => {
    // Arrange
    const harness = createControllerHarness(createState({ errorMessage: "error" }))

    // Act
    harness.controller.openCreateModal()

    // Assert
    expect(harness.getState().errorMessage).toBe("")
  })

  test("openEditModal は selectedEntry の source をフォームへ反映する", () => {
    // Arrange
    const selectedEntry = createSelectedEntry()
    const harness = createControllerHarness(createState({ selectedEntry }))

    // Act
    harness.controller.openEditModal()

    // Assert
    expect(harness.getState().formSource).toBe("Dragon Priest")
  })

  test("openEditModal は selectedEntry の category をフォームへ反映する", () => {
    // Arrange
    const selectedEntry = createSelectedEntry()
    const harness = createControllerHarness(createState({ selectedEntry }))

    // Act
    harness.controller.openEditModal()

    // Assert
    expect(harness.getState().formCategory).toBe("地名")
  })

  test("openEditModal は selectedEntry の origin をフォームへ反映する", () => {
    // Arrange
    const selectedEntry = createSelectedEntry()
    const harness = createControllerHarness(createState({ selectedEntry }))

    // Act
    harness.controller.openEditModal()

    // Assert
    expect(harness.getState().formOrigin).toBe("確認待ち")
  })

  test("openEditModal は selectedEntry の translation をフォームへ反映する", () => {
    // Arrange
    const selectedEntry = createSelectedEntry()
    const harness = createControllerHarness(createState({ selectedEntry }))

    // Act
    harness.controller.openEditModal()

    // Assert
    expect(harness.getState().formTranslation).toBe("ドラゴン・プリースト")
  })

  test("openEditModal は modalState を edit へ切り替える", () => {
    // Arrange
    const selectedEntry = createSelectedEntry()
    const harness = createControllerHarness(createState({ selectedEntry }))

    // Act
    harness.controller.openEditModal()

    // Assert
    expect(harness.getState().modalState).toBe("edit")
  })

  test("openEditModal は errorMessage をクリアする", () => {
    // Arrange
    const selectedEntry = createSelectedEntry()
    const harness = createControllerHarness(
      createState({ selectedEntry, errorMessage: "error" })
    )

    // Act
    harness.controller.openEditModal()

    // Assert
    expect(harness.getState().errorMessage).toBe("")
  })

  test("openEditModal は selectedEntry が無い時に store を更新しない", () => {
    // Arrange
    const harness = createControllerHarness(createState())

    // Act
    harness.controller.openEditModal()

    // Assert
    expect(harness.store.update).not.toHaveBeenCalled()
  })

  test("openDeleteModal は modalState を delete へ切り替える", () => {
    // Arrange
    const selectedEntry = createSelectedEntry()
    const harness = createControllerHarness(createState({ selectedEntry }))

    // Act
    harness.controller.openDeleteModal()

    // Assert
    expect(harness.getState().modalState).toBe("delete")
  })

  test("openDeleteModal は errorMessage をクリアする", () => {
    // Arrange
    const selectedEntry = createSelectedEntry()
    const harness = createControllerHarness(
      createState({ selectedEntry, errorMessage: "error" })
    )

    // Act
    harness.controller.openDeleteModal()

    // Assert
    expect(harness.getState().errorMessage).toBe("")
  })

  test("openDeleteModal は selectedEntry が無い時に store を更新しない", () => {
    // Arrange
    const harness = createControllerHarness(createState())

    // Act
    harness.controller.openDeleteModal()

    // Assert
    expect(harness.store.update).not.toHaveBeenCalled()
  })

  test("closeEditModal は edit modal を閉じる", () => {
    // Arrange
    const harness = createControllerHarness(createState({ modalState: "edit" }))

    // Act
    harness.controller.closeEditModal()

    // Assert
    expect(harness.getState().modalState).toBeNull()
  })

  test("closeEditModal は create modal を閉じる", () => {
    // Arrange
    const harness = createControllerHarness(createState({ modalState: "create" }))

    // Act
    harness.controller.closeEditModal()

    // Assert
    expect(harness.getState().modalState).toBeNull()
  })

  test("closeEditModal は delete modal を維持する", () => {
    // Arrange
    const harness = createControllerHarness(createState({ modalState: "delete" }))

    // Act
    harness.controller.closeEditModal()

    // Assert
    expect(harness.getState().modalState).toBe("delete")
  })

  test("closeDeleteModal は delete modal を閉じる", () => {
    // Arrange
    const harness = createControllerHarness(createState({ modalState: "delete" }))

    // Act
    harness.controller.closeDeleteModal()

    // Assert
    expect(harness.getState().modalState).toBeNull()
  })

  test("closeDeleteModal は edit modal を維持する", () => {
    // Arrange
    const harness = createControllerHarness(createState({ modalState: "edit" }))

    // Act
    harness.controller.closeDeleteModal()

    // Assert
    expect(harness.getState().modalState).toBe("edit")
  })

  test("selectRow は useCase.selectEntry へ id を委譲する", async () => {
    // Arrange
    const harness = createControllerHarness()

    // Act
    await harness.controller.selectRow("42")

    // Assert
    expect(harness.useCase.selectEntry).toHaveBeenCalledWith("42")
  })

  test("saveCurrentEntry は useCase.saveCurrentEntry を呼ぶ", async () => {
    // Arrange
    const harness = createControllerHarness()

    // Act
    await harness.controller.saveCurrentEntry()

    // Assert
    expect(harness.useCase.saveCurrentEntry).toHaveBeenCalledTimes(1)
  })

  test("deleteCurrentEntry は useCase.deleteCurrentEntry を呼ぶ", async () => {
    // Arrange
    const harness = createControllerHarness()

    // Act
    await harness.controller.deleteCurrentEntry()

    // Assert
    expect(harness.useCase.deleteCurrentEntry).toHaveBeenCalledTimes(1)
  })

  test("startImport は runtime 購読済みなら true を渡す", async () => {
    // Arrange
    const harness = createControllerHarness()
    await harness.controller.mount()

    // Act
    await harness.controller.startImport()

    // Assert
    expect(harness.useCase.startStagedXmlImport).toHaveBeenCalledWith(true)
  })

  test("startImport は runtime 未購読なら false を渡す", async () => {
    // Arrange
    const harness = createControllerHarness()
    harness.runtimeEventAdapter.subscribe.mockReturnValue(false)
    await harness.controller.mount()

    // Act
    await harness.controller.startImport()

    // Assert
    expect(harness.useCase.startStagedXmlImport).toHaveBeenCalledWith(false)
  })

  test("handleSearchInput は query を更新する", () => {
    // Arrange
    const harness = createControllerHarness(
      createState({ query: "old", category: "地名", page: 2, errorMessage: "error" })
    )
    const searchInput = document.createElement("input")
    searchInput.value = "Dragon"
    const event = new Event("input")
    Object.defineProperty(event, "currentTarget", { value: searchInput })

    // Act
    harness.controller.handleSearchInput(event)

    // Assert
    expect(harness.getState().query).toBe("Dragon")
  })

  test("handleSearchInput は page を 0 へ戻す", () => {
    // Arrange
    const harness = createControllerHarness(createState({ page: 2 }))
    const searchInput = document.createElement("input")
    searchInput.value = "Dragon"
    const event = new Event("input")
    Object.defineProperty(event, "currentTarget", { value: searchInput })

    // Act
    harness.controller.handleSearchInput(event)

    // Assert
    expect(harness.getState().page).toBe(0)
  })

  test("handleSearchInput は errorMessage をクリアする", () => {
    // Arrange
    const harness = createControllerHarness(createState({ errorMessage: "error" }))
    const searchInput = document.createElement("input")
    searchInput.value = "Dragon"
    const event = new Event("input")
    Object.defineProperty(event, "currentTarget", { value: searchInput })

    // Act
    harness.controller.handleSearchInput(event)

    // Assert
    expect(harness.getState().errorMessage).toBe("")
  })

  test("handleSearchInput は useCase.loadEntries を呼ぶ", () => {
    // Arrange
    const harness = createControllerHarness()
    const searchInput = document.createElement("input")
    searchInput.value = "Dragon"
    const event = new Event("input")
    Object.defineProperty(event, "currentTarget", { value: searchInput })

    // Act
    harness.controller.handleSearchInput(event)

    // Assert
    expect(harness.useCase.loadEntries).toHaveBeenCalledTimes(1)
  })

  test("handleSearchInput は input 要素以外を無視する", () => {
    // Arrange
    const harness = createControllerHarness()

    // Act
    harness.controller.handleSearchInput(new Event("input"))

    // Assert
    expect(harness.store.update).not.toHaveBeenCalled()
  })

  test("handleCategoryChange は category を更新する", () => {
    // Arrange
    const harness = createControllerHarness(
      createState({ query: "old", category: "地名", page: 2, errorMessage: "error" })
    )
    const categorySelect = document.createElement("select")
    categorySelect.innerHTML = '<option value="固有名詞">固有名詞</option>'
    categorySelect.value = "固有名詞"
    const event = new Event("change")
    Object.defineProperty(event, "currentTarget", { value: categorySelect })

    // Act
    harness.controller.handleCategoryChange(event)

    // Assert
    expect(harness.getState().category).toBe("固有名詞")
  })

  test("handleCategoryChange は page を 0 へ戻す", () => {
    // Arrange
    const harness = createControllerHarness(createState({ page: 2 }))
    const categorySelect = document.createElement("select")
    categorySelect.innerHTML = '<option value="固有名詞">固有名詞</option>'
    categorySelect.value = "固有名詞"
    const event = new Event("change")
    Object.defineProperty(event, "currentTarget", { value: categorySelect })

    // Act
    harness.controller.handleCategoryChange(event)

    // Assert
    expect(harness.getState().page).toBe(0)
  })

  test("handleCategoryChange は errorMessage をクリアする", () => {
    // Arrange
    const harness = createControllerHarness(createState({ errorMessage: "error" }))
    const categorySelect = document.createElement("select")
    categorySelect.innerHTML = '<option value="固有名詞">固有名詞</option>'
    categorySelect.value = "固有名詞"
    const event = new Event("change")
    Object.defineProperty(event, "currentTarget", { value: categorySelect })

    // Act
    harness.controller.handleCategoryChange(event)

    // Assert
    expect(harness.getState().errorMessage).toBe("")
  })

  test("handleCategoryChange は useCase.loadEntries を呼ぶ", () => {
    // Arrange
    const harness = createControllerHarness()
    const categorySelect = document.createElement("select")
    categorySelect.innerHTML = '<option value="固有名詞">固有名詞</option>'
    categorySelect.value = "固有名詞"
    const event = new Event("change")
    Object.defineProperty(event, "currentTarget", { value: categorySelect })

    // Act
    harness.controller.handleCategoryChange(event)

    // Assert
    expect(harness.useCase.loadEntries).toHaveBeenCalledTimes(1)
  })

  test("handleCategoryChange は select 要素以外を無視する", () => {
    // Arrange
    const harness = createControllerHarness()

    // Act
    harness.controller.handleCategoryChange(new Event("change"))

    // Assert
    expect(harness.store.update).not.toHaveBeenCalled()
  })

  test("goToPrevPage は先頭ページで store を更新しない", () => {
    // Arrange
    const harness = createControllerHarness(createState({ page: 0 }))

    // Act
    harness.controller.goToPrevPage()

    // Assert
    expect(harness.store.update).not.toHaveBeenCalled()
  })

  test("goToPrevPage は先頭ページで loadEntries を呼ばない", () => {
    // Arrange
    const harness = createControllerHarness(createState({ page: 0 }))

    // Act
    harness.controller.goToPrevPage()

    // Assert
    expect(harness.useCase.loadEntries).not.toHaveBeenCalled()
  })

  test("goToPrevPage は前ページへ戻す", () => {
    // Arrange
    const harness = createControllerHarness(createState({ page: 2 }))

    // Act
    harness.controller.goToPrevPage()

    // Assert
    expect(harness.getState().page).toBe(1)
  })

  test("goToPrevPage は移動時に loadEntries を呼ぶ", () => {
    // Arrange
    const harness = createControllerHarness(createState({ page: 2 }))

    // Act
    harness.controller.goToPrevPage()

    // Assert
    expect(harness.useCase.loadEntries).toHaveBeenCalledTimes(1)
  })

  test("goToNextPage は最終ページで store を更新しない", () => {
    // Arrange
    const harness = createControllerHarness(createState({ page: 0, totalCount: 30 }))

    // Act
    harness.controller.goToNextPage()

    // Assert
    expect(harness.store.update).not.toHaveBeenCalled()
  })

  test("goToNextPage は最終ページで loadEntries を呼ばない", () => {
    // Arrange
    const harness = createControllerHarness(createState({ page: 0, totalCount: 30 }))

    // Act
    harness.controller.goToNextPage()

    // Assert
    expect(harness.useCase.loadEntries).not.toHaveBeenCalled()
  })

  test("goToNextPage は次ページへ進める", () => {
    // Arrange
    const harness = createControllerHarness(createState({ page: 0, totalCount: 61 }))

    // Act
    harness.controller.goToNextPage()

    // Assert
    expect(harness.getState().page).toBe(1)
  })

  test("goToNextPage は移動時に loadEntries を呼ぶ", () => {
    // Arrange
    const harness = createControllerHarness(createState({ page: 0, totalCount: 61 }))

    // Act
    harness.controller.goToNextPage()

    // Assert
    expect(harness.useCase.loadEntries).toHaveBeenCalledTimes(1)
  })

  test("stageXmlImport は selectedFileName を更新する", () => {
    // Arrange
    const harness = createControllerHarness(createState({ errorMessage: "error" }))
    const file = new File(["<Root />"], "master.xml", { type: "text/xml" })
    Object.defineProperty(file, "path", { value: "", configurable: true })
    Object.defineProperty(file, "webkitRelativePath", {
      value: "mods/master.xml",
      configurable: true
    })

    // Act
    harness.controller.stageXmlImport(file)

    // Assert
    expect(harness.getState().selectedFileName).toBe("master.xml")
  })

  test("stageXmlImport は selectedFileReference を解決する", () => {
    // Arrange
    const harness = createControllerHarness(createState({ errorMessage: "error" }))
    const file = new File(["<Root />"], "master.xml", { type: "text/xml" })
    Object.defineProperty(file, "path", { value: "", configurable: true })
    Object.defineProperty(file, "webkitRelativePath", {
      value: "mods/master.xml",
      configurable: true
    })

    // Act
    harness.controller.stageXmlImport(file)

    // Assert
    expect(harness.getState().selectedFileReference).toBe("mods/master.xml")
  })

  test("stageXmlImport は importStage を ready へ切り替える", () => {
    // Arrange
    const harness = createControllerHarness(createState({ errorMessage: "error" }))
    const file = new File(["<Root />"], "master.xml", { type: "text/xml" })

    // Act
    harness.controller.stageXmlImport(file)

    // Assert
    expect(harness.getState().importStage).toBe("ready")
  })

  test("stageXmlImport は importProgress を 0 に戻す", () => {
    // Arrange
    const harness = createControllerHarness(createState({ importProgress: 80 }))
    const file = new File(["<Root />"], "master.xml", { type: "text/xml" })

    // Act
    harness.controller.stageXmlImport(file)

    // Assert
    expect(harness.getState().importProgress).toBe(0)
  })

  test("stageXmlImport は importSummary を null に戻す", () => {
    // Arrange
    const harness = createControllerHarness(
      createState({
        importSummary: {
          fileName: "before.xml",
          importedCount: 1,
          updatedCount: 0,
          totalCount: 1,
          selectedSource: "Before"
        }
      })
    )
    const file = new File(["<Root />"], "master.xml", { type: "text/xml" })

    // Act
    harness.controller.stageXmlImport(file)

    // Assert
    expect(harness.getState().importSummary).toBeNull()
  })

  test("stageXmlImport は errorMessage をクリアする", () => {
    // Arrange
    const harness = createControllerHarness(createState({ errorMessage: "error" }))
    const file = new File(["<Root />"], "master.xml", { type: "text/xml" })

    // Act
    harness.controller.stageXmlImport(file)

    // Assert
    expect(harness.getState().errorMessage).toBe("")
  })

  test("stageXmlImport は null の時に selectedFileName を未選択へ戻す", () => {
    // Arrange
    const harness = createControllerHarness(createState({ selectedFileName: "master.xml" }))

    // Act
    harness.controller.stageXmlImport(null)

    // Assert
    expect(harness.getState().selectedFileName).toBe("未選択")
  })

  test("stageXmlImport は null の時に selectedFileReference を null へ戻す", () => {
    // Arrange
    const harness = createControllerHarness(
      createState({ selectedFileReference: "master.xml" })
    )

    // Act
    harness.controller.stageXmlImport(null)

    // Assert
    expect(harness.getState().selectedFileReference).toBeNull()
  })

  test("stageXmlImport は null の時に importStage を idle へ戻す", () => {
    // Arrange
    const harness = createControllerHarness(createState({ importStage: "done" }))

    // Act
    harness.controller.stageXmlImport(null)

    // Assert
    expect(harness.getState().importStage).toBe("idle")
  })

  test("resetImportSelection は selectedFileName を未選択へ戻す", () => {
    // Arrange
    const harness = createControllerHarness(
      createState({
        selectedFileName: "master.xml",
        selectedFileReference: "master.xml",
        importStage: "done",
        importProgress: 100,
        importSummary: {
          fileName: "master.xml",
          importedCount: 1,
          updatedCount: 0,
          totalCount: 1,
          selectedSource: "Dragon Priest"
        }
      })
    )

    // Act
    harness.controller.resetImportSelection()

    // Assert
    expect(harness.getState().selectedFileName).toBe("未選択")
  })

  test("resetImportSelection は selectedFileReference を null へ戻す", () => {
    // Arrange
    const harness = createControllerHarness(
      createState({
        selectedFileName: "master.xml",
        selectedFileReference: "master.xml",
        importStage: "done"
      })
    )

    // Act
    harness.controller.resetImportSelection()

    // Assert
    expect(harness.getState().selectedFileReference).toBeNull()
  })

  test("resetImportSelection は importStage を idle へ戻す", () => {
    // Arrange
    const harness = createControllerHarness(
      createState({
        selectedFileName: "master.xml",
        selectedFileReference: "master.xml",
        importStage: "done"
      })
    )

    // Act
    harness.controller.resetImportSelection()

    // Assert
    expect(harness.getState().importStage).toBe("idle")
  })

  test("resetImportSelection は importProgress を 0 へ戻す", () => {
    // Arrange
    const harness = createControllerHarness(
      createState({
        selectedFileName: "master.xml",
        selectedFileReference: "master.xml",
        importStage: "done",
        importProgress: 100
      })
    )

    // Act
    harness.controller.resetImportSelection()

    // Assert
    expect(harness.getState().importProgress).toBe(0)
  })

  test("resetImportSelection は importSummary を null へ戻す", () => {
    // Arrange
    const harness = createControllerHarness(
      createState({
        selectedFileName: "master.xml",
        selectedFileReference: "master.xml",
        importStage: "done",
        importSummary: {
          fileName: "master.xml",
          importedCount: 1,
          updatedCount: 0,
          totalCount: 1,
          selectedSource: "Dragon Priest"
        }
      })
    )

    // Act
    harness.controller.resetImportSelection()

    // Assert
    expect(harness.getState().importSummary).toBeNull()
  })

  test("resetImportSelection は running 中の selectedFileName を維持する", () => {
    // Arrange
    const harness = createControllerHarness(
      createState({
        selectedFileName: "master.xml",
        selectedFileReference: "master.xml",
        importStage: "running"
      })
    )

    // Act
    harness.controller.resetImportSelection()

    // Assert
    expect(harness.getState().selectedFileName).toBe("master.xml")
  })

  test("resetImportSelection は running 中の selectedFileReference を維持する", () => {
    // Arrange
    const harness = createControllerHarness(
      createState({
        selectedFileName: "master.xml",
        selectedFileReference: "master.xml",
        importStage: "running"
      })
    )

    // Act
    harness.controller.resetImportSelection()

    // Assert
    expect(harness.getState().selectedFileReference).toBe("master.xml")
  })

  test("resetImportSelection は running 中の importStage を維持する", () => {
    // Arrange
    const harness = createControllerHarness(
      createState({
        selectedFileName: "master.xml",
        selectedFileReference: "master.xml",
        importStage: "running"
      })
    )

    // Act
    harness.controller.resetImportSelection()

    // Assert
    expect(harness.getState().importStage).toBe("running")
  })

  test("setFormSource は input 値を反映する", () => {
    // Arrange
    const harness = createControllerHarness()
    const sourceInput = document.createElement("input")
    sourceInput.value = "Source"
    const event = new Event("input")
    Object.defineProperty(event, "currentTarget", { value: sourceInput })

    // Act
    harness.controller.setFormSource(event)

    // Assert
    expect(harness.getState().formSource).toBe("Source")
  })

  test("setFormSource は input 以外を無視する", () => {
    // Arrange
    const harness = createControllerHarness()

    // Act
    harness.controller.setFormSource(new Event("input"))

    // Assert
    expect(harness.store.update).not.toHaveBeenCalled()
  })

  test("setFormCategory は select 値を反映する", () => {
    // Arrange
    const harness = createControllerHarness()
    const categorySelect = document.createElement("select")
    categorySelect.innerHTML = '<option value="地名">地名</option>'
    categorySelect.value = "地名"
    const event = new Event("change")
    Object.defineProperty(event, "currentTarget", { value: categorySelect })

    // Act
    harness.controller.setFormCategory(event)

    // Assert
    expect(harness.getState().formCategory).toBe("地名")
  })

  test("setFormCategory は select 以外を無視する", () => {
    // Arrange
    const harness = createControllerHarness()

    // Act
    harness.controller.setFormCategory(new Event("change"))

    // Assert
    expect(harness.store.update).not.toHaveBeenCalled()
  })

  test("setFormOrigin は select 値を反映する", () => {
    // Arrange
    const harness = createControllerHarness()
    const originSelect = document.createElement("select")
    originSelect.innerHTML = '<option value="確認待ち">確認待ち</option>'
    originSelect.value = "確認待ち"
    const event = new Event("change")
    Object.defineProperty(event, "currentTarget", { value: originSelect })

    // Act
    harness.controller.setFormOrigin(event)

    // Assert
    expect(harness.getState().formOrigin).toBe("確認待ち")
  })

  test("setFormOrigin は select 以外を無視する", () => {
    // Arrange
    const harness = createControllerHarness()

    // Act
    harness.controller.setFormOrigin(new Event("change"))

    // Assert
    expect(harness.store.update).not.toHaveBeenCalled()
  })

  test("setFormTranslation は textarea 値を反映する", () => {
    // Arrange
    const harness = createControllerHarness()
    const translationTextArea = document.createElement("textarea")
    translationTextArea.value = "訳語"
    const event = new Event("input")
    Object.defineProperty(event, "currentTarget", { value: translationTextArea })

    // Act
    harness.controller.setFormTranslation(event)

    // Assert
    expect(harness.getState().formTranslation).toBe("訳語")
  })

  test("setFormTranslation は textarea 以外を無視する", () => {
    // Arrange
    const harness = createControllerHarness()

    // Act
    harness.controller.setFormTranslation(new Event("input"))

    // Assert
    expect(harness.store.update).not.toHaveBeenCalled()
  })
})
