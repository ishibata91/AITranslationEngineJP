import { describe, expect, test } from "vitest"

import { MasterDictionaryPresenter } from "./master-dictionary.presenter"
import type { MasterDictionaryScreenState } from "./master-dictionary-screen-types"

function createState(overrides: Partial<MasterDictionaryScreenState> = {}): MasterDictionaryScreenState {
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

const EXPECTED_BASE_CATEGORIES = [
  "固有名詞",
  "NPC",
  "地名",
  "装備",
  "アイテム",
  "書籍",
  "設備",
  "シャウト",
  "その他"
]

describe("MasterDictionaryPresenter", () => {
  test("カテゴリ候補はページ外の backend 正規化カテゴリも含む", () => {
    const presenter = new MasterDictionaryPresenter()
    const viewModel = presenter.toViewModel(
      createState({
        entries: [
          {
            id: "1",
            source: "Ancient Vampire",
            translation: "太古の吸血鬼",
            category: "NPC",
            origin: "XML取込",
            updatedAt: "2026-04-12T00:00:00Z"
          },
          {
            id: "2",
            source: "Dawnguard Shield",
            translation: "ドーンガードの盾",
            category: "装備",
            origin: "XML取込",
            updatedAt: "2026-04-12T00:01:00Z"
          }
        ]
      }),
      true
    )

    expect(viewModel.categoryOptions[0]).toBe("すべて")
    for (const category of EXPECTED_BASE_CATEGORIES) {
      expect(viewModel.categoryOptions).toContain(category)
    }

    expect(new Set(viewModel.categoryOptions).size).toBe(viewModel.categoryOptions.length)
  })
})
