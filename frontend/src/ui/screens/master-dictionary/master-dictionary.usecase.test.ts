import { describe, expect, test } from "vitest"

import { MasterDictionaryStore } from "./master-dictionary.store"
import { MasterDictionaryUseCase } from "./master-dictionary.usecase"

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
            rec: "NPC_:FULL",
            edid: "DLC1VampireLordAncient",
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
        selectedRec: ["CONT:FULL", "BOOK:FULL", "NPC_:FULL"],
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
