import { describe, expect, test, vi } from "vitest"

import type {
  TranslationInputScreenState,
  TranslationInputScreenViewModel
} from "@application/gateway-contract/translation-input"

import { TranslationInputScreenController } from "./translation-input-screen-controller"

function createState(
  overrides: Partial<TranslationInputScreenState> = {}
): TranslationInputScreenState {
  return {
    items: [],
    selectedItemId: null,
    stagedFile: null,
    operationState: "idle",
    errorMessage: "",
    latestResponse: null,
    ...overrides
  }
}

function createViewModel(
  state: TranslationInputScreenState
): TranslationInputScreenViewModel {
  return {
    ...state,
    selectedItem: null,
    gatewayStatus: "接続準備済み",
    hasStagedFile: state.stagedFile !== null,
    canImport: state.stagedFile !== null && state.operationState === "ready",
    canRebuildSelected: false,
    isImporting: state.operationState === "importing",
    isRebuilding: state.operationState === "rebuilding",
    stagedFileName: state.stagedFile?.fileName ?? "未選択",
    stagedFilePath: state.stagedFile?.filePath ?? "-",
    stagedFileHash: state.stagedFile?.fileHash ?? "-",
    operationStatusLabel: state.operationState,
    operationStatusText: state.operationState,
    latestOutcomeTitle: "-",
    latestOutcomeText: "-",
    selectionStatusText: "-",
    totalItemCountLabel: `${state.items.length} 件`,
    emptyStateText: "empty"
  }
}

function createHarness(initialState: TranslationInputScreenState = createState()) {
  let state = initialState
  const listeners = new Set<(state: TranslationInputScreenState) => void>()

  const store = {
    subscribe: vi.fn((listener: (nextState: TranslationInputScreenState) => void) => {
      listeners.add(listener)
      return () => {
        listeners.delete(listener)
      }
    }),
    snapshot: vi.fn(() => state),
    update: vi.fn((mutator: (draft: TranslationInputScreenState) => void) => {
      const draft = structuredClone(state)
      mutator(draft)
      state = draft
      for (const listener of listeners) {
        listener(state)
      }
    })
  }

  const presenter = {
    toViewModel: vi.fn((nextState: TranslationInputScreenState) =>
      createViewModel(nextState)
    )
  }

  const useCase = {
    startImport: vi.fn(async () => {}),
    rebuildSelected: vi.fn(async () => {})
  }

  const controller = new TranslationInputScreenController({
    isGatewayConnected: true,
    store,
    presenter,
    useCase
  })

  return {
    controller,
    getState: () => state,
    useCase
  }
}

describe("TranslationInputScreenController", () => {
  test("stageJsonImport は file path と digest hash を stagedFile へ保持する", async () => {
    const digestSpy = vi
      .spyOn(globalThis.crypto.subtle, "digest")
      .mockResolvedValue(new Uint8Array([0, 1, 255]).buffer)
    const harness = createHarness()
    const file = new File(["{}"], "input-review.json", {
      type: "application/json"
    })

    Object.defineProperty(file, "path", {
      value: "/mods/input-review.json"
    })
    Object.defineProperty(file, "arrayBuffer", {
      value: vi.fn(() => Promise.resolve(new TextEncoder().encode("{}").buffer))
    })

    await harness.controller.stageJsonImport(file)

    expect(harness.getState().operationState).toBe("ready")
    expect(harness.getState().stagedFile).toEqual({
      fileName: "input-review.json",
      filePath: "/mods/input-review.json",
      fileHash: "0001ff"
    })

    digestSpy.mockRestore()
  })

  test("stageJsonImport は path がない file で bare filename を stagedFile.filePath へ保持する", async () => {
    const digestSpy = vi
      .spyOn(globalThis.crypto.subtle, "digest")
      .mockResolvedValue(new Uint8Array([0xaa, 0xbb]).buffer)
    const harness = createHarness()
    const file = new File(["{}"], "uploaded.json", {
      type: "application/json"
    })

    Object.defineProperty(file, "arrayBuffer", {
      value: vi.fn(() => Promise.resolve(new TextEncoder().encode("{}").buffer))
    })

    await harness.controller.stageJsonImport(file)

    expect(harness.getState().stagedFile?.fileName).toBe("uploaded.json")
    expect(harness.getState().stagedFile?.filePath).toBe("uploaded.json")

    digestSpy.mockRestore()
  })

  test("resetImportSelection は importing 中でない時に stagedFile をクリアする", () => {
    const harness = createHarness(
      createState({
        stagedFile: {
          fileName: "input-review.json",
          filePath: "/mods/input-review.json",
          fileHash: "hash-41"
        },
        operationState: "ready",
        errorMessage: "old error"
      })
    )

    harness.controller.resetImportSelection()

    expect(harness.getState().stagedFile).toBeNull()
    expect(harness.getState().operationState).toBe("idle")
    expect(harness.getState().errorMessage).toBe("")
  })

  test("selectItem は selectedItemId を更新し errorMessage を消す", () => {
    const harness = createHarness(
      createState({
        selectedItemId: null,
        errorMessage: "old error"
      })
    )

    harness.controller.selectItem("item-42")

    expect(harness.getState().selectedItemId).toBe("item-42")
    expect(harness.getState().errorMessage).toBe("")
  })

  test("startImport と rebuildSelected は useCase を委譲する", async () => {
    const harness = createHarness()

    await harness.controller.startImport()
    await harness.controller.rebuildSelected()

    expect(harness.useCase.startImport).toHaveBeenCalledTimes(1)
    expect(harness.useCase.rebuildSelected).toHaveBeenCalledTimes(1)
  })
})