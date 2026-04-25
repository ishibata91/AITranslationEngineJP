import { describe, expect, test, vi } from "vitest"

import { TranslationInputUseCase } from "./translation-input.usecase"

type UseCaseGateway = NonNullable<ConstructorParameters<typeof TranslationInputUseCase>[0]>
type StoreLike = ConstructorParameters<typeof TranslationInputUseCase>[1]
type TestState = ReturnType<StoreLike["snapshot"]>
type TestCommandResponse = Awaited<
  ReturnType<UseCaseGateway["importTranslationInput"]>
>
type TestSummary = NonNullable<TestCommandResponse["summary"]>
type TestReviewItem = TestState["items"][number]

function createSummary(
  overrides: Partial<TestSummary> = {}
): TestSummary {
  return {
    input: {
      id: 41,
      sourceFilePath: "/mods/input-review.json",
      sourceTool: "xEdit",
      targetPluginName: "Skyrim.esm",
      targetPluginType: "esm",
      recordCount: 8,
      importedAt: "2026-04-26T09:30:00Z"
    },
    translationRecordCount: 5,
    translationFieldCount: 9,
    categories: [
      {
        category: "NPC",
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
    warnings: [],
    ...overrides
  }
}

function createResponse(
  overrides: Partial<TestCommandResponse> = {}
): TestCommandResponse {
  return {
    accepted: true,
    summary: createSummary(),
    warnings: [],
    ...overrides
  }
}

function createItem(
  overrides: Partial<TestReviewItem> = {}
): TestReviewItem {
  return {
    localId: "input-41",
    inputId: 41,
    fileName: "input-review.json",
    filePath: "/mods/input-review.json",
    fileHash: "hash-41",
    importTimestamp: "2026-04-26T09:30:00Z",
    status: "registered",
    accepted: true,
    canRebuild: true,
    lastAction: "import",
    errorKind: null,
    warnings: [],
    summary: createSummary(),
    ...overrides
  }
}

function createStore(initialState?: Partial<TestState>) {
  let state: TestState = {
    items: [],
    selectedItemId: null,
    stagedFile: null,
    operationState: "idle",
    errorMessage: "",
    latestResponse: null,
    ...initialState
  }

  const store: StoreLike = {
    snapshot() {
      return structuredClone(state)
    },
    update(mutator: (draft: TestState) => void) {
      const draft = structuredClone(state)
      mutator(draft)
      state = draft
    }
  }

  return store
}

function makeGateway(
  partial: Partial<{
    importTranslationInput: (request: { filePath: string }) => Promise<TestCommandResponse>
    rebuildTranslationInputCache: (request: { inputId: number }) => Promise<TestCommandResponse>
  }>
): UseCaseGateway {
  return {
    importTranslationInput: vi.fn(),
    rebuildTranslationInputCache: vi.fn(),
    ...partial
  } as UseCaseGateway
}

describe("TranslationInputUseCase", () => {
  test("startImport は staged file がない時に missing required field を保持する", async () => {
    const store = createStore()
    const useCase = new TranslationInputUseCase(null, store)

    await useCase.startImport()

    expect(store.snapshot().errorMessage).toBe(
      "登録する JSON file を選択してください。"
    )
    expect(store.snapshot().latestResponse?.errorKind).toBe(
      "missing_required_field"
    )
  })

  test("startImport は warning を含む登録結果を review item として保持する", async () => {
    const store = createStore()
    const importTranslationInput = vi.fn().mockResolvedValue(
      createResponse({
        warnings: [
          {
            kind: "unknown_field_definition",
            recordType: "BOOK",
            subrecordType: "DESC",
            message: "unknown description field"
          }
        ],
        summary: createSummary({
          warnings: [
            {
              kind: "unknown_field_definition",
              recordType: "BOOK",
              subrecordType: "DESC",
              message: "unknown description field"
            }
          ]
        })
      })
    )
    const gateway = makeGateway({ importTranslationInput })
    const useCase = new TranslationInputUseCase(gateway, store)

    store.update((draft) => {
      draft.stagedFile = {
        fileName: "input-review.json",
        filePath: "/mods/input-review.json",
        fileHash: "hash-41"
      }
      draft.operationState = "ready"
    })

    await useCase.startImport()

    const state = store.snapshot()

    expect(importTranslationInput).toHaveBeenCalledWith({
      filePath: "/mods/input-review.json"
    })
    expect(state.operationState).toBe("idle")
    expect(state.stagedFile).toBeNull()
    expect(state.items).toHaveLength(1)
    expect(state.items[0]?.status).toBe("warning")
    expect(state.items[0]?.warnings).toHaveLength(1)
    expect(state.selectedItemId).toBe(state.items[0]?.localId)
  })

  test("rebuildSelected は inputId がない時に cache missing を保持する", async () => {
    const store = createStore()
    const useCase = new TranslationInputUseCase(null, store)

    store.update((draft) => {
      draft.items = [createItem({ localId: "missing-cache", inputId: null, canRebuild: false })]
      draft.selectedItemId = "missing-cache"
    })

    await useCase.rebuildSelected()

    expect(store.snapshot().errorMessage).toBe("再構築対象の cache がありません。")
    expect(store.snapshot().latestResponse?.errorKind).toBe("cache_missing")
  })

  test("rebuildSelected は source file missing を rebuild-required として更新する", async () => {
    const store = createStore()
    const rebuildTranslationInputCache = vi.fn().mockResolvedValue(
      createResponse({
        accepted: false,
        errorKind: "source_file_missing",
        summary: undefined,
        warnings: []
      })
    )
    const gateway = makeGateway({ rebuildTranslationInputCache })
    const useCase = new TranslationInputUseCase(gateway, store)

    store.update((draft) => {
      draft.items = [createItem()]
      draft.selectedItemId = "input-41"
    })

    await useCase.rebuildSelected()

    const updatedItem = store.snapshot().items[0]

    expect(rebuildTranslationInputCache).toHaveBeenCalledWith({ inputId: 41 })
    expect(updatedItem?.status).toBe("rebuild-required")
    expect(updatedItem?.errorKind).toBe("source_file_missing")
    expect(updatedItem?.canRebuild).toBe(true)
  })
})