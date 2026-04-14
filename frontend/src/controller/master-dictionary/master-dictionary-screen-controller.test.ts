import { describe, expect, test, vi } from "vitest"

import type {
  MasterDictionaryScreenState,
  MasterDictionaryScreenViewModel
} from "@application/contract/master-dictionary/master-dictionary-screen-types"

import { MasterDictionaryScreenController } from "./master-dictionary-screen-controller"

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
  test("mount / dispose / subscribe / getViewModel は依存 seam を使う", async () => {
    const harness = createControllerHarness(
      createState({
        totalCount: 1,
        selectedId: "101"
      })
    )
    const listener = vi.fn()

    const unsubscribe = harness.controller.subscribe(listener)
    await harness.controller.mount()

    expect(harness.runtimeEventAdapter.subscribe).toHaveBeenCalledTimes(1)
    expect(harness.useCase.loadEntries).toHaveBeenCalledTimes(1)

    harness.store.update((draft) => {
      draft.selectedId = "102"
    })

    expect(listener).toHaveBeenCalled()
    expect(harness.presenter.toViewModel).toHaveBeenCalled()
    expect(harness.controller.getViewModel().selectionStatusText).toBe("102")

    unsubscribe()
    harness.controller.dispose()

    expect(harness.runtimeEventAdapter.detach).toHaveBeenCalledTimes(1)
  })

  test("openCreateModal は既定値でフォームを初期化する", () => {
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

    harness.controller.openCreateModal()

    expect(harness.getState()).toMatchObject({
      modalState: "create",
      formSource: "",
      formCategory: "固有名詞",
      formOrigin: "手動登録",
      formTranslation: "",
      errorMessage: ""
    })
  })

  test("openEditModal / openDeleteModal は selectedEntry がある時だけ反映する", () => {
    const selectedEntry = {
      id: "101",
      source: "Dragon Priest",
      translation: "ドラゴン・プリースト",
      category: "地名",
      origin: "確認待ち",
      updatedAt: "2026-01-01 00:00",
      note: "note"
    }
    const harness = createControllerHarness(createState({ selectedEntry }))

    harness.controller.openEditModal()
    expect(harness.getState()).toMatchObject({
      modalState: "edit",
      formSource: "Dragon Priest",
      formCategory: "地名",
      formOrigin: "確認待ち",
      formTranslation: "ドラゴン・プリースト",
      errorMessage: ""
    })

    harness.controller.openDeleteModal()
    expect(harness.getState().modalState).toBe("delete")

    const withoutSelection = createControllerHarness(createState())
    withoutSelection.controller.openEditModal()
    withoutSelection.controller.openDeleteModal()
    expect(withoutSelection.store.update).not.toHaveBeenCalled()
  })

  test("closeEditModal / closeDeleteModal は対象 modal の時だけ閉じる", () => {
    const editHarness = createControllerHarness(createState({ modalState: "edit" }))
    editHarness.controller.closeEditModal()
    expect(editHarness.getState().modalState).toBeNull()

    const createHarness = createControllerHarness(createState({ modalState: "create" }))
    createHarness.controller.closeEditModal()
    expect(createHarness.getState().modalState).toBeNull()

    const untouchedEdit = createControllerHarness(createState({ modalState: "delete" }))
    untouchedEdit.controller.closeEditModal()
    expect(untouchedEdit.getState().modalState).toBe("delete")

    const deleteHarness = createControllerHarness(createState({ modalState: "delete" }))
    deleteHarness.controller.closeDeleteModal()
    expect(deleteHarness.getState().modalState).toBeNull()

    const untouchedDelete = createControllerHarness(createState({ modalState: "edit" }))
    untouchedDelete.controller.closeDeleteModal()
    expect(untouchedDelete.getState().modalState).toBe("edit")
  })

  test("select/save/delete/startImport は useCase へ委譲する", async () => {
    const subscribedHarness = createControllerHarness()
    await subscribedHarness.controller.mount()
    await subscribedHarness.controller.selectRow("42")
    await subscribedHarness.controller.saveCurrentEntry()
    await subscribedHarness.controller.deleteCurrentEntry()
    await subscribedHarness.controller.startImport()

    expect(subscribedHarness.useCase.selectEntry).toHaveBeenCalledWith("42")
    expect(subscribedHarness.useCase.saveCurrentEntry).toHaveBeenCalledTimes(1)
    expect(subscribedHarness.useCase.deleteCurrentEntry).toHaveBeenCalledTimes(1)
    expect(subscribedHarness.useCase.startStagedXmlImport).toHaveBeenCalledWith(true)

    const unsubscribedHarness = createControllerHarness()
    unsubscribedHarness.runtimeEventAdapter.subscribe.mockReturnValue(false)
    await unsubscribedHarness.controller.mount()
    await unsubscribedHarness.controller.startImport()
    expect(unsubscribedHarness.useCase.startStagedXmlImport).toHaveBeenCalledWith(false)
  })

  test("検索とカテゴリ変更は input/select の時だけ state を更新して reload する", () => {
    const harness = createControllerHarness(
      createState({ query: "old", category: "地名", page: 2, errorMessage: "error" })
    )
    const searchInput = document.createElement("input")
    searchInput.value = "Dragon"
    const searchEvent = new Event("input")
    Object.defineProperty(searchEvent, "currentTarget", {
      value: searchInput
    })

    harness.controller.handleSearchInput(searchEvent)
    expect(harness.getState()).toMatchObject({
      query: "Dragon",
      page: 0,
      errorMessage: ""
    })
    expect(harness.useCase.loadEntries).toHaveBeenCalledTimes(1)

    const categorySelect = document.createElement("select")
    categorySelect.innerHTML = '<option value="固有名詞">固有名詞</option>'
    categorySelect.value = "固有名詞"
    const categoryEvent = new Event("change")
    Object.defineProperty(categoryEvent, "currentTarget", {
      value: categorySelect
    })

    harness.controller.handleCategoryChange(categoryEvent)
    expect(harness.getState()).toMatchObject({
      category: "固有名詞",
      page: 0,
      errorMessage: ""
    })
    expect(harness.useCase.loadEntries).toHaveBeenCalledTimes(2)

    const noopHarness = createControllerHarness()
    noopHarness.controller.handleSearchInput(new Event("input"))
    noopHarness.controller.handleCategoryChange(new Event("change"))
    expect(noopHarness.store.update).not.toHaveBeenCalled()
    expect(noopHarness.useCase.loadEntries).not.toHaveBeenCalled()
  })

  test("前後ページ移動は境界外では何もしない", () => {
    const prevHarness = createControllerHarness(createState({ page: 0 }))
    prevHarness.controller.goToPrevPage()
    expect(prevHarness.store.update).not.toHaveBeenCalled()
    expect(prevHarness.useCase.loadEntries).not.toHaveBeenCalled()

    const movablePrev = createControllerHarness(createState({ page: 2 }))
    movablePrev.controller.goToPrevPage()
    expect(movablePrev.getState().page).toBe(1)
    expect(movablePrev.useCase.loadEntries).toHaveBeenCalledTimes(1)

    const nextHarness = createControllerHarness(createState({ page: 0, totalCount: 30 }))
    nextHarness.controller.goToNextPage()
    expect(nextHarness.store.update).not.toHaveBeenCalled()
    expect(nextHarness.useCase.loadEntries).not.toHaveBeenCalled()

    const movableNext = createControllerHarness(createState({ page: 0, totalCount: 61 }))
    movableNext.controller.goToNextPage()
    expect(movableNext.getState().page).toBe(1)
    expect(movableNext.useCase.loadEntries).toHaveBeenCalledTimes(1)
  })

  test("stageXmlImport / resetImportSelection は file と importStage 境界を扱う", () => {
    const harness = createControllerHarness(createState({ errorMessage: "error" }))
    const file = new File(["<Root />"], "master.xml", { type: "text/xml" })
    Object.defineProperty(file, "path", {
      value: "",
      configurable: true
    })
    Object.defineProperty(file, "webkitRelativePath", {
      value: "mods/master.xml",
      configurable: true
    })

    harness.controller.stageXmlImport(file)
    expect(harness.getState()).toMatchObject({
      selectedFileName: "master.xml",
      selectedFileReference: "mods/master.xml",
      importStage: "ready",
      importProgress: 0,
      importSummary: null,
      errorMessage: ""
    })

    harness.controller.stageXmlImport(null)
    expect(harness.getState()).toMatchObject({
      selectedFileName: "未選択",
      selectedFileReference: null,
      importStage: "idle"
    })

    const resetHarness = createControllerHarness(
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
    resetHarness.controller.resetImportSelection()
    expect(resetHarness.getState()).toMatchObject({
      selectedFileName: "未選択",
      selectedFileReference: null,
      importStage: "idle",
      importProgress: 0,
      importSummary: null
    })

    const runningHarness = createControllerHarness(
      createState({
        selectedFileName: "master.xml",
        selectedFileReference: "master.xml",
        importStage: "running"
      })
    )
    runningHarness.controller.resetImportSelection()
    expect(runningHarness.getState()).toMatchObject({
      selectedFileName: "master.xml",
      selectedFileReference: "master.xml",
      importStage: "running"
    })
  })

  test("form setter は正しい要素型の時だけ更新する", () => {
    const harness = createControllerHarness()

    const sourceInput = document.createElement("input")
    sourceInput.value = "Source"
    const sourceEvent = new Event("input")
    Object.defineProperty(sourceEvent, "currentTarget", { value: sourceInput })
    harness.controller.setFormSource(sourceEvent)

    const categorySelect = document.createElement("select")
    categorySelect.innerHTML = '<option value="地名">地名</option>'
    categorySelect.value = "地名"
    const categoryEvent = new Event("change")
    Object.defineProperty(categoryEvent, "currentTarget", { value: categorySelect })
    harness.controller.setFormCategory(categoryEvent)

    const originSelect = document.createElement("select")
    originSelect.innerHTML = '<option value="確認待ち">確認待ち</option>'
    originSelect.value = "確認待ち"
    const originEvent = new Event("change")
    Object.defineProperty(originEvent, "currentTarget", { value: originSelect })
    harness.controller.setFormOrigin(originEvent)

    const translationTextArea = document.createElement("textarea")
    translationTextArea.value = "訳語"
    const translationEvent = new Event("input")
    Object.defineProperty(translationEvent, "currentTarget", {
      value: translationTextArea
    })
    harness.controller.setFormTranslation(translationEvent)

    expect(harness.getState()).toMatchObject({
      formSource: "Source",
      formCategory: "地名",
      formOrigin: "確認待ち",
      formTranslation: "訳語"
    })

    const noopHarness = createControllerHarness()
    noopHarness.controller.setFormSource(new Event("input"))
    noopHarness.controller.setFormCategory(new Event("change"))
    noopHarness.controller.setFormOrigin(new Event("change"))
    noopHarness.controller.setFormTranslation(new Event("input"))
    expect(noopHarness.store.update).not.toHaveBeenCalled()
  })
})
