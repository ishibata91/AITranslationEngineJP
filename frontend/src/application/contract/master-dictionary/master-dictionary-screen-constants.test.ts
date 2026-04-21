import { describe, expect, test } from "vitest"

import type { MasterDictionaryScreenState } from "./master-dictionary-screen-types"
import { buildUpsertPayload } from "./master-dictionary-screen-constants"

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

describe("buildUpsertPayload", () => {
  test("source の先頭と末尾の空白をトリムする", () => {
    // Arrange
    const state = createState({
      formSource: "  Dragon Priest  ",
      formTranslation: "ドラゴン・プリースト"
    })

    // Act
    const payload = buildUpsertPayload(state)

    // Assert
    expect(payload.source).toBe("Dragon Priest")
  })

  test("translation の先頭と末尾の空白をトリムする", () => {
    // Arrange
    const state = createState({
      formSource: "Dragon Priest",
      formTranslation: "  ドラゴン・プリースト  "
    })

    // Act
    const payload = buildUpsertPayload(state)

    // Assert
    expect(payload.translation).toBe("ドラゴン・プリースト")
  })

  test("source が空白のみの場合は空文字を返す", () => {
    // Arrange
    const state = createState({
      formSource: "   ",
      formTranslation: "ドラゴン・プリースト"
    })

    // Act
    const payload = buildUpsertPayload(state)

    // Assert
    expect(payload.source).toBe("")
  })

  test("translation が空白のみの場合は空文字を返す", () => {
    // Arrange
    const state = createState({
      formSource: "Dragon Priest",
      formTranslation: "   "
    })

    // Act
    const payload = buildUpsertPayload(state)

    // Assert
    expect(payload.translation).toBe("")
  })

  test("category と origin はトリムせずそのまま返す", () => {
    // Arrange
    const state = createState({
      formSource: "Dragon",
      formTranslation: "ドラゴン",
      formCategory: "地名",
      formOrigin: "確認待ち"
    })

    // Act
    const payload = buildUpsertPayload(state)

    // Assert
    expect(payload.category).toBe("地名")
    expect(payload.origin).toBe("確認待ち")
  })

  test("両フィールドが空文字の場合は両方空文字を返す", () => {
    // Arrange
    const state = createState({ formSource: "", formTranslation: "" })

    // Act
    const payload = buildUpsertPayload(state)

    // Assert
    expect(payload.source).toBe("")
    expect(payload.translation).toBe("")
  })
})
