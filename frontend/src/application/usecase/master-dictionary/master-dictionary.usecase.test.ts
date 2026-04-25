import { describe, expect, test, vi } from "vitest"

import type { MasterDictionaryGatewayContract } from "@application/gateway-contract/master-dictionary"
import { MasterDictionaryStore } from "@application/store/master-dictionary"

import { MasterDictionaryUseCase } from "./master-dictionary.usecase"

function makeGateway(
  partial: Partial<MasterDictionaryGatewayContract>
): MasterDictionaryGatewayContract {
  return {
    listMasterDictionaryEntries: vi.fn(),
    getMasterDictionaryEntry: vi.fn(),
    createMasterDictionaryEntry: vi.fn(),
    updateMasterDictionaryEntry: vi.fn(),
    deleteMasterDictionaryEntry: vi.fn(),
    importMasterDictionaryXml: vi.fn(),
    ...partial
  } as unknown as MasterDictionaryGatewayContract
}

describe("MasterDictionaryUseCase", () => {
  test("handleImportCompleted は imported / updated / total を importSummary に保持する", async () => {
    const store = new MasterDictionaryStore()
    const useCase = new MasterDictionaryUseCase(null, store)

    store.update((draft) => {
      draft.selectedFileName = "Dawnguard_english_japanese.xml"
      draft.importStage = "running"
      draft.totalCount = 40
    })

    await useCase.handleImportCompleted({
      page: {
        items: [
          {
            id: 740,
            source: "Ancient Vampire",
            translation: "太古の吸血鬼",
            category: "NPC",
            origin: "XML取込",
            updatedAt: "2026-04-12T00:00:00Z"
          }
        ],
        totalCount: 740,
        page: 1,
        pageSize: 30,
        selectedId: 740
      },
      summary: {
        filePath: "Dawnguard_english_japanese.xml",
        fileName: "Dawnguard_english_japanese.xml",
        importedCount: 700,
        updatedCount: 947,
        skippedCount: 7235,
        lastEntryId: 740
      }
    })

    const state = store.snapshot()

    expect(state.importStage).toBe("done")
    expect(state.importSummary).not.toBeNull()
    expect(state.importSummary?.importedCount).toBe(700)
    expect(state.importSummary?.updatedCount).toBe(947)
    expect(state.importSummary?.totalCount).toBe(740)
    expect(state.importSummary?.selectedSource).toBe("Ancient Vampire")
  })
})

describe("loadEntryDetail", () => {
  test("gateway がない時は selectedEntry を null にする", async () => {
    // Arrange
    const store = new MasterDictionaryStore()
    const useCase = new MasterDictionaryUseCase(null, store)
    store.update((draft) => {
      draft.selectedEntry = {
        id: "1",
        source: "Dragonborn",
        translation: "ドラゴンボーン",
        category: "固有名詞",
        origin: "手動登録",
        note: "マスター辞書エントリ",
        updatedAt: "2026-01-01T00:00:00Z"
      }
    })

    // Act
    await useCase.loadEntryDetail("1")

    // Assert
    expect(store.snapshot().selectedEntry).toBeNull()
  })

  test("id が null の時は gateway を呼ばずに selectedEntry を null にする", async () => {
    // Arrange
    const store = new MasterDictionaryStore()
    const getMasterDictionaryEntry = vi.fn()
    const gateway = makeGateway({ getMasterDictionaryEntry })
    const useCase = new MasterDictionaryUseCase(gateway, store)

    // Act
    await useCase.loadEntryDetail(null)

    // Assert
    expect(store.snapshot().selectedEntry).toBeNull()
    expect(getMasterDictionaryEntry).not.toHaveBeenCalled()
  })

  test("gateway が null エントリを返す時は selectedEntry を null にする", async () => {
    // Arrange
    const store = new MasterDictionaryStore()
    const getMasterDictionaryEntry = vi.fn().mockResolvedValue({ entry: null })
    const gateway = makeGateway({ getMasterDictionaryEntry })
    const useCase = new MasterDictionaryUseCase(gateway, store)

    // Act
    await useCase.loadEntryDetail("42")

    // Assert
    expect(store.snapshot().selectedEntry).toBeNull()
  })

  test("gateway が有効なエントリを返す時は selectedEntry を設定する", async () => {
    // Arrange
    const store = new MasterDictionaryStore()
    const entry = {
      id: "42",
      source: "Whiterun",
      translation: "ホワイトラン",
      category: "地名",
      origin: "XML取込",
      note: "マスター辞書エントリ",
      updatedAt: "2026-04-10T00:00:00Z"
    }
    const getMasterDictionaryEntry = vi.fn().mockResolvedValue({ entry })
    const gateway = makeGateway({ getMasterDictionaryEntry })
    const useCase = new MasterDictionaryUseCase(gateway, store)

    // Act
    await useCase.loadEntryDetail("42")

    // Assert
    expect(store.snapshot().selectedEntry?.source).toBe("Whiterun")
    expect(store.snapshot().selectedEntry?.note).toBe("マスター辞書エントリ")
    const entryRecordA = store.snapshot().selectedEntry as Record<string, unknown> | null
    expect(entryRecordA?.["rec"]).toBeUndefined()
    expect(entryRecordA?.["edid"]).toBeUndefined()
    expect(getMasterDictionaryEntry).toHaveBeenCalledWith({ id: "42" })
  })

  test("gateway がエラーをスローした時は selectedEntry を null にして errorMessage を設定する", async () => {
    // Arrange
    const store = new MasterDictionaryStore()
    const getMasterDictionaryEntry = vi
      .fn()
      .mockRejectedValue(new Error("サーバーエラー"))
    const gateway = makeGateway({ getMasterDictionaryEntry })
    const useCase = new MasterDictionaryUseCase(gateway, store)

    // Act
    await useCase.loadEntryDetail("99")

    // Assert
    expect(store.snapshot().selectedEntry).toBeNull()
    expect(store.snapshot().errorMessage).toBe("サーバーエラー")
  })
})

describe("selectEntry", () => {
  test("selectedId を設定して gateway からエントリを取得する", async () => {
    // Arrange
    const store = new MasterDictionaryStore()
    const entry = {
      id: "55",
      source: "Ulfric Stormcloak",
      translation: "ウルフリック・ストームクローク",
      category: "NPC",
      origin: "XML取込",
      note: "マスター辞書エントリ",
      updatedAt: "2026-03-10T00:00:00Z"
    }
    const getMasterDictionaryEntry = vi.fn().mockResolvedValue({ entry })
    const gateway = makeGateway({ getMasterDictionaryEntry })
    const useCase = new MasterDictionaryUseCase(gateway, store)

    // Act
    await useCase.selectEntry("55")

    // Assert
    expect(store.snapshot().selectedId).toBe("55")
    expect(store.snapshot().selectedEntry?.source).toBe("Ulfric Stormcloak")
    const entryRecordB = store.snapshot().selectedEntry as Record<string, unknown> | null
    expect(entryRecordB?.["rec"]).toBeUndefined()
    expect(entryRecordB?.["edid"]).toBeUndefined()
    expect(getMasterDictionaryEntry).toHaveBeenCalledWith({ id: "55" })
  })
})

describe("loadEntries - null gateway", () => {
  test("gateway がない時は entries と selectedId と selectedEntry を空にする", async () => {
    // Arrange
    const store = new MasterDictionaryStore()
    const useCase = new MasterDictionaryUseCase(null, store)
    store.update((draft) => {
      draft.entries = [
        {
          id: "1",
          source: "Skyrim",
          translation: "スカイリム",
          category: "地名",
          origin: "手動登録",
          updatedAt: "2026-01-01T00:00:00Z"
        }
      ]
      draft.selectedId = "1"
    })

    // Act
    await useCase.loadEntries()

    // Assert
    expect(store.snapshot().entries).toHaveLength(0)
    expect(store.snapshot().totalCount).toBe(0)
    expect(store.snapshot().selectedId).toBeNull()
    expect(store.snapshot().selectedEntry).toBeNull()
  })
})

// ---------------------------------------------------------------------------
// saveCurrentEntry - create
// ---------------------------------------------------------------------------

function makeCreatedEntry(source = "Dragon Priest") {
  return {
    id: "99",
    source,
    translation: "ドラゴン・プリースト",
    category: "固有名詞",
    origin: "手動登録",
    note: "マスター辞書エントリ",
    updatedAt: "2026-04-21T00:00:00Z"
  }
}

function makePageState(
  items: { id: number; source: string; translation: string; category: string; origin: string; updatedAt: string }[] = [],
  selectedId?: number
) {
  return {
    items,
    totalCount: items.length,
    page: 1,
    pageSize: 30,
    selectedId
  }
}

describe("saveCurrentEntry - create", () => {
  test("gateway がない時は createMasterDictionaryEntry を呼ばない", async () => {
    // Arrange
    const store = new MasterDictionaryStore()
    const useCase = new MasterDictionaryUseCase(null, store)
    const createMasterDictionaryEntry = vi.fn()
    store.update((draft) => {
      draft.modalState = "create"
      draft.formSource = "Dragon Priest"
      draft.formTranslation = "ドラゴン・プリースト"
    })

    // Act
    await useCase.saveCurrentEntry()

    // Assert
    expect(createMasterDictionaryEntry).not.toHaveBeenCalled()
  })

  test("modalState が null の時は createMasterDictionaryEntry を呼ばない", async () => {
    // Arrange
    const store = new MasterDictionaryStore()
    const createMasterDictionaryEntry = vi.fn()
    const gateway = makeGateway({ createMasterDictionaryEntry })
    const useCase = new MasterDictionaryUseCase(gateway, store)
    store.update((draft) => {
      draft.modalState = null
      draft.formSource = "Dragon Priest"
      draft.formTranslation = "ドラゴン・プリースト"
    })

    // Act
    await useCase.saveCurrentEntry()

    // Assert
    expect(createMasterDictionaryEntry).not.toHaveBeenCalled()
  })

  test("source が空の時は errorMessage を設定する", async () => {
    // Arrange
    const store = new MasterDictionaryStore()
    const createMasterDictionaryEntry = vi.fn()
    const gateway = makeGateway({ createMasterDictionaryEntry })
    const useCase = new MasterDictionaryUseCase(gateway, store)
    store.update((draft) => {
      draft.modalState = "create"
      draft.formSource = ""
      draft.formTranslation = "ドラゴン・プリースト"
    })

    // Act
    await useCase.saveCurrentEntry()

    // Assert
    expect(store.snapshot().errorMessage).toBe("原文と訳語を入力してください。")
    expect(createMasterDictionaryEntry).not.toHaveBeenCalled()
  })

  test("source が空白のみの時は trim 後に空判定となり errorMessage を設定する", async () => {
    // Arrange
    const store = new MasterDictionaryStore()
    const createMasterDictionaryEntry = vi.fn()
    const gateway = makeGateway({ createMasterDictionaryEntry })
    const useCase = new MasterDictionaryUseCase(gateway, store)
    store.update((draft) => {
      draft.modalState = "create"
      draft.formSource = "   "
      draft.formTranslation = "ドラゴン・プリースト"
    })

    // Act
    await useCase.saveCurrentEntry()

    // Assert
    expect(store.snapshot().errorMessage).toBe("原文と訳語を入力してください。")
    expect(createMasterDictionaryEntry).not.toHaveBeenCalled()
  })

  test("translation が空の時は errorMessage を設定する", async () => {
    // Arrange
    const store = new MasterDictionaryStore()
    const createMasterDictionaryEntry = vi.fn()
    const gateway = makeGateway({ createMasterDictionaryEntry })
    const useCase = new MasterDictionaryUseCase(gateway, store)
    store.update((draft) => {
      draft.modalState = "create"
      draft.formSource = "Dragon Priest"
      draft.formTranslation = ""
    })

    // Act
    await useCase.saveCurrentEntry()

    // Assert
    expect(store.snapshot().errorMessage).toBe("原文と訳語を入力してください。")
    expect(createMasterDictionaryEntry).not.toHaveBeenCalled()
  })

  test("create 成功時は createMasterDictionaryEntry に trimmed source が渡る", async () => {
    // Arrange
    const store = new MasterDictionaryStore()
    const entry = makeCreatedEntry()
    const pageItems = [{ id: 99, source: "Dragon Priest", translation: "ドラゴン・プリースト", category: "固有名詞", origin: "手動登録", updatedAt: "2026-04-21T00:00:00Z" }]
    const createMasterDictionaryEntry = vi.fn().mockResolvedValue({
      entry,
      refreshTargetId: "99",
      page: makePageState(pageItems, 99)
    })
    const gateway = makeGateway({ createMasterDictionaryEntry })
    const useCase = new MasterDictionaryUseCase(gateway, store)
    store.update((draft) => {
      draft.modalState = "create"
      draft.formSource = "  Dragon Priest  "
      draft.formTranslation = "ドラゴン・プリースト"
      draft.formCategory = "固有名詞"
      draft.formOrigin = "手動登録"
    })

    // Act
    await useCase.saveCurrentEntry()

    // Assert
    expect(createMasterDictionaryEntry).toHaveBeenCalledTimes(1)
    const callArg = createMasterDictionaryEntry.mock.calls[0][0] as { payload: { source: string } }
    expect(callArg.payload.source).toBe("Dragon Priest")
  })

  test("create 成功時は modalState が null になる", async () => {
    // Arrange
    const store = new MasterDictionaryStore()
    const entry = makeCreatedEntry()
    const pageItems = [{ id: 99, source: "Dragon Priest", translation: "ドラゴン・プリースト", category: "固有名詞", origin: "手動登録", updatedAt: "2026-04-21T00:00:00Z" }]
    const gateway = makeGateway({
      createMasterDictionaryEntry: vi.fn().mockResolvedValue({
        entry,
        refreshTargetId: "99",
        page: makePageState(pageItems, 99)
      })
    })
    const useCase = new MasterDictionaryUseCase(gateway, store)
    store.update((draft) => {
      draft.modalState = "create"
      draft.formSource = "Dragon Priest"
      draft.formTranslation = "ドラゴン・プリースト"
    })

    // Act
    await useCase.saveCurrentEntry()

    // Assert
    expect(store.snapshot().modalState).toBeNull()
  })

  test("create 失敗時は errorMessage が設定されて modalState を維持する", async () => {
    // Arrange
    const store = new MasterDictionaryStore()
    const gateway = makeGateway({
      createMasterDictionaryEntry: vi.fn().mockRejectedValue(
        new Error("duplicate_entry: trim(source_term)+translated_term が重複しています")
      )
    })
    const useCase = new MasterDictionaryUseCase(gateway, store)
    store.update((draft) => {
      draft.modalState = "create"
      draft.formSource = "Dragon Priest"
      draft.formTranslation = "ドラゴン・プリースト"
    })

    // Act
    await useCase.saveCurrentEntry()

    // Assert
    expect(store.snapshot().errorMessage).toContain("duplicate_entry")
    expect(store.snapshot().modalState).toBe("create")
  })
})

// ---------------------------------------------------------------------------
// saveCurrentEntry - edit
// ---------------------------------------------------------------------------

describe("saveCurrentEntry - edit", () => {
  test("selectedId がない時は errorMessage を設定して updateMasterDictionaryEntry を呼ばない", async () => {
    // Arrange
    const store = new MasterDictionaryStore()
    const updateMasterDictionaryEntry = vi.fn()
    const gateway = makeGateway({ updateMasterDictionaryEntry })
    const useCase = new MasterDictionaryUseCase(gateway, store)
    store.update((draft) => {
      draft.modalState = "edit"
      draft.selectedId = null
      draft.formSource = "Dragon Priest"
      draft.formTranslation = "ドラゴン・プリースト"
    })

    // Act
    await useCase.saveCurrentEntry()

    // Assert
    expect(store.snapshot().errorMessage).toBe("更新対象が選択されていません。")
    expect(updateMasterDictionaryEntry).not.toHaveBeenCalled()
  })

  test("edit 成功時は updateMasterDictionaryEntry に trimmed source と translation が渡る", async () => {
    // Arrange
    const store = new MasterDictionaryStore()
    const updatedEntry = {
      id: "1",
      source: "Dragon Priest",
      translation: "ドラゴン・プリースト",
      category: "固有名詞",
      origin: "手動登録",
      note: "マスター辞書エントリ",
      updatedAt: "2026-04-21T00:00:00Z"
    }
    const pageItems = [{ id: 1, source: "Dragon Priest", translation: "ドラゴン・プリースト", category: "固有名詞", origin: "手動登録", updatedAt: "2026-04-21T00:00:00Z" }]
    const updateMasterDictionaryEntry = vi.fn().mockResolvedValue({
      entry: updatedEntry,
      refreshTargetId: "1",
      page: makePageState(pageItems, 1)
    })
    const gateway = makeGateway({ updateMasterDictionaryEntry })
    const useCase = new MasterDictionaryUseCase(gateway, store)
    store.update((draft) => {
      draft.modalState = "edit"
      draft.selectedId = "1"
      draft.formSource = "  Dragon Priest  "
      draft.formTranslation = "  ドラゴン・プリースト  "
      draft.formCategory = "固有名詞"
      draft.formOrigin = "手動登録"
    })

    // Act
    await useCase.saveCurrentEntry()

    // Assert
    expect(updateMasterDictionaryEntry).toHaveBeenCalledTimes(1)
    const callArg = updateMasterDictionaryEntry.mock.calls[0][0] as { payload: { source: string; translation: string } }
    expect(callArg.payload.source).toBe("Dragon Priest")
    expect(callArg.payload.translation).toBe("ドラゴン・プリースト")
  })

  test("edit 成功時は modalState が null になる", async () => {
    // Arrange
    const store = new MasterDictionaryStore()
    const updatedEntry = {
      id: "1",
      source: "Dragon Priest",
      translation: "ドラゴン・プリースト",
      category: "固有名詞",
      origin: "手動登録",
      note: "マスター辞書エントリ",
      updatedAt: "2026-04-21T00:00:00Z"
    }
    const pageItems = [{ id: 1, source: "Dragon Priest", translation: "ドラゴン・プリースト", category: "固有名詞", origin: "手動登録", updatedAt: "2026-04-21T00:00:00Z" }]
    const gateway = makeGateway({
      updateMasterDictionaryEntry: vi.fn().mockResolvedValue({
        entry: updatedEntry,
        refreshTargetId: "1",
        page: makePageState(pageItems, 1)
      })
    })
    const useCase = new MasterDictionaryUseCase(gateway, store)
    store.update((draft) => {
      draft.modalState = "edit"
      draft.selectedId = "1"
      draft.formSource = "Dragon Priest"
      draft.formTranslation = "ドラゴン・プリースト"
    })

    // Act
    await useCase.saveCurrentEntry()

    // Assert
    expect(store.snapshot().modalState).toBeNull()
  })

  test("edit 失敗時は errorMessage が設定されて modalState を維持する", async () => {
    // Arrange
    const store = new MasterDictionaryStore()
    const gateway = makeGateway({
      updateMasterDictionaryEntry: vi.fn().mockRejectedValue(
        new Error("duplicate_entry: trim(source_term)+translated_term が重複しています")
      )
    })
    const useCase = new MasterDictionaryUseCase(gateway, store)
    store.update((draft) => {
      draft.modalState = "edit"
      draft.selectedId = "1"
      draft.formSource = "Dragon Priest"
      draft.formTranslation = "ドラゴン・プリースト"
    })

    // Act
    await useCase.saveCurrentEntry()

    // Assert
    expect(store.snapshot().errorMessage).toContain("duplicate_entry")
    expect(store.snapshot().modalState).toBe("edit")
  })
})

// ---------------------------------------------------------------------------
// deleteCurrentEntry
// ---------------------------------------------------------------------------

describe("deleteCurrentEntry", () => {
  test("gateway がない時は deleteMasterDictionaryEntry を呼ばない", async () => {
    // Arrange
    const store = new MasterDictionaryStore()
    const useCase = new MasterDictionaryUseCase(null, store)
    const deleteMasterDictionaryEntry = vi.fn()
    store.update((draft) => {
      draft.selectedId = "1"
      draft.modalState = "delete"
    })

    // Act
    await useCase.deleteCurrentEntry()

    // Assert
    expect(deleteMasterDictionaryEntry).not.toHaveBeenCalled()
  })

  test("selectedId がない時は deleteMasterDictionaryEntry を呼ばない", async () => {
    // Arrange
    const store = new MasterDictionaryStore()
    const deleteMasterDictionaryEntry = vi.fn()
    const gateway = makeGateway({ deleteMasterDictionaryEntry })
    const useCase = new MasterDictionaryUseCase(gateway, store)
    store.update((draft) => {
      draft.selectedId = null
    })

    // Act
    await useCase.deleteCurrentEntry()

    // Assert
    expect(deleteMasterDictionaryEntry).not.toHaveBeenCalled()
  })

  test("delete 成功時は deleteMasterDictionaryEntry に selectedId が渡る", async () => {
    // Arrange
    const store = new MasterDictionaryStore()
    const deleteMasterDictionaryEntry = vi.fn().mockResolvedValue({
      deletedId: "1",
      nextSelectedId: null,
      page: makePageState()
    })
    const gateway = makeGateway({ deleteMasterDictionaryEntry })
    const useCase = new MasterDictionaryUseCase(gateway, store)
    store.update((draft) => {
      draft.selectedId = "1"
      draft.modalState = "delete"
    })

    // Act
    await useCase.deleteCurrentEntry()

    // Assert
    expect(deleteMasterDictionaryEntry).toHaveBeenCalledTimes(1)
    const callArg = deleteMasterDictionaryEntry.mock.calls[0][0] as { id: string }
    expect(callArg.id).toBe("1")
  })

  test("delete 成功時は modalState が null になる", async () => {
    // Arrange
    const store = new MasterDictionaryStore()
    const gateway = makeGateway({
      deleteMasterDictionaryEntry: vi.fn().mockResolvedValue({
        deletedId: "1",
        nextSelectedId: null,
        page: makePageState()
      })
    })
    const useCase = new MasterDictionaryUseCase(gateway, store)
    store.update((draft) => {
      draft.selectedId = "1"
      draft.modalState = "delete"
    })

    // Act
    await useCase.deleteCurrentEntry()

    // Assert
    expect(store.snapshot().modalState).toBeNull()
  })

  test("delete 失敗時は errorMessage が設定されて modalState を維持する", async () => {
    // Arrange
    const store = new MasterDictionaryStore()
    const gateway = makeGateway({
      deleteMasterDictionaryEntry: vi.fn().mockRejectedValue(
        new Error("削除対象が見つかりません")
      )
    })
    const useCase = new MasterDictionaryUseCase(gateway, store)
    store.update((draft) => {
      draft.selectedId = "1"
      draft.modalState = "delete"
    })

    // Act
    await useCase.deleteCurrentEntry()

    // Assert
    expect(store.snapshot().errorMessage).toBe("削除対象が見つかりません")
    expect(store.snapshot().modalState).toBe("delete")
  })
})

describe("startStagedXmlImport", () => {
  test("import 失敗時は plain object の message を errorMessage に表示する", async () => {
    // Arrange
    const store = new MasterDictionaryStore()
    const gateway = makeGateway({
      importMasterDictionaryXml: vi.fn().mockRejectedValue({
        message: "import failed from plain object"
      })
    })
    const useCase = new MasterDictionaryUseCase(gateway, store)
    store.update((draft) => {
      draft.selectedFileName = "broken.xml"
      draft.selectedFileReference = "/tmp/broken.xml"
      draft.importStage = "ready"
      draft.importProgress = 10
      draft.importSummary = {
        fileName: "old.xml",
        importedCount: 1,
        updatedCount: 2,
        totalCount: 3,
        selectedSource: "Old"
      }
    })

    // Act
    await useCase.startStagedXmlImport(false)

    // Assert
    expect(store.snapshot().importStage).toBe("ready")
    expect(store.snapshot().importProgress).toBe(0)
    expect(store.snapshot().importSummary).toBeNull()
    expect(store.snapshot().errorMessage).toBe("import failed from plain object")
  })

  test("import 失敗時は string reject を errorMessage に表示する", async () => {
    // Arrange
    const store = new MasterDictionaryStore()
    const gateway = makeGateway({
      importMasterDictionaryXml: vi.fn().mockRejectedValue(
        "import failed from string"
      )
    })
    const useCase = new MasterDictionaryUseCase(gateway, store)
    store.update((draft) => {
      draft.selectedFileName = "broken.xml"
      draft.selectedFileReference = "/tmp/broken.xml"
      draft.importStage = "ready"
    })

    // Act
    await useCase.startStagedXmlImport(false)

    // Assert
    expect(store.snapshot().importStage).toBe("ready")
    expect(store.snapshot().importProgress).toBe(0)
    expect(store.snapshot().importSummary).toBeNull()
    expect(store.snapshot().errorMessage).toBe("import failed from string")
  })
})
