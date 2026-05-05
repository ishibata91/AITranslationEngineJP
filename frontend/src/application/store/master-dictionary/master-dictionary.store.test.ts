import { describe, expect, it } from "vitest"

import { DEFAULT_CATEGORY, DEFAULT_ORIGIN } from "@application/contract/master-dictionary"

import { MasterDictionaryStore } from "./master-dictionary.store"

describe("MasterDictionaryStore", () => {
  it("subscribe immediately notifies current state", () => {
    const store = new MasterDictionaryStore()
    const snapshots: unknown[] = []

    store.subscribe((state) => {
      snapshots.push(state)
    })

    expect(snapshots).toHaveLength(1)
    expect(store.snapshot()).toMatchObject({
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
      formCategory: DEFAULT_CATEGORY,
      formOrigin: DEFAULT_ORIGIN,
      formTranslation: "",
      selectedFileName: "未選択",
      selectedFileReference: null,
      importStage: "idle",
      importProgress: 0,
      importSummary: null
    })
  })

  it("unsubscribe removes listener from future updates", () => {
    const store = new MasterDictionaryStore()
    let callCount = 0

    const unsubscribe = store.subscribe(() => {
      callCount += 1
    })
    unsubscribe()

    store.update((draft) => {
      draft.query = "changed"
    })

    expect(callCount).toBe(1)
  })

  it("snapshot returns defensive copies", () => {
    const store = new MasterDictionaryStore()

    store.update((draft) => {
      draft.entries = [
        {
          id: "1",
          source: "Test Source A",
          translation: "Test Translation A",
          category: "NPC_",
          origin: "manual",
          updatedAt: "2026-05-05T00:00:00Z"
        }
      ]
      draft.selectedEntry = {
        id: "1",
        source: "Test Source A",
        translation: "Test Translation A",
        category: "NPC_",
        origin: "manual",
        updatedAt: "2026-05-05T00:00:00Z",
        note: "Test Note A"
      }
      draft.importSummary = {
        fileName: "test-dictionary.xml",
        importedCount: 1,
        updatedCount: 2,
        totalCount: 3,
        selectedSource: "Test Source A"
      }
    })

    const snapshot = store.snapshot()
    snapshot.entries.push({
      id: "2",
      source: "Test Source B",
      translation: "Test Translation B",
      category: "NPC_",
      origin: "manual",
      updatedAt: "2026-05-05T00:00:00Z"
    })
    snapshot.selectedEntry!.translation = "changed"
    snapshot.importSummary!.fileName = "changed.xml"

    const nextSnapshot = store.snapshot()
    expect(nextSnapshot.entries).toHaveLength(1)
    expect(nextSnapshot.selectedEntry?.translation).toBe("Test Translation A")
    expect(nextSnapshot.importSummary?.fileName).toBe("test-dictionary.xml")
  })
})
