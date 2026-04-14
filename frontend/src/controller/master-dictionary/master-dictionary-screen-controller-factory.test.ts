import { afterEach, describe, expect, test, vi } from "vitest"

import type { MasterDictionaryGatewayContract } from "@application/gateway-contract/master-dictionary"

import { createMasterDictionaryScreenControllerFactory } from "./master-dictionary-screen-controller-factory"

function createGateway(): {
  gateway: MasterDictionaryGatewayContract
  listMasterDictionaryEntries: ReturnType<typeof vi.fn>
  getMasterDictionaryEntry: ReturnType<typeof vi.fn>
} {
  const listMasterDictionaryEntries = vi.fn(() =>
    Promise.resolve({
      entries: [
        {
          id: "101",
          source: "Dragon Priest",
          translation: "ドラゴン・プリースト",
          category: "固有名詞",
          origin: "初期データ",
          updatedAt: "2026-01-01 00:00"
        }
      ],
      totalCount: 1,
      page: 1,
      pageSize: 30
    })
  )
  const getMasterDictionaryEntry = vi.fn(() =>
    Promise.resolve({
      entry: {
        id: "101",
        source: "Dragon Priest",
        translation: "ドラゴン・プリースト",
        category: "固有名詞",
        origin: "初期データ",
        updatedAt: "2026-01-01 00:00",
        note: "REC: NPC_:FULL / EDID: SeedDragonPriest"
      }
    })
  )

  return {
    gateway: {
      listMasterDictionaryEntries,
      getMasterDictionaryEntry,
      createMasterDictionaryEntry: vi.fn(() => Promise.reject(new Error("not used"))),
      updateMasterDictionaryEntry: vi.fn(() => Promise.reject(new Error("not used"))),
      deleteMasterDictionaryEntry: vi.fn(() => Promise.reject(new Error("not used"))),
      importMasterDictionaryXml: vi.fn(() => Promise.reject(new Error("not used")))
    },
    listMasterDictionaryEntries,
    getMasterDictionaryEntry
  }
}

describe("createMasterDictionaryScreenControllerFactory", () => {
  afterEach(() => {
    delete (window as Window & { runtime?: unknown }).runtime
  })

  test("gateway null は未接続 status を返す", () => {
    // Arrange
    const controller = createMasterDictionaryScreenControllerFactory(null)()

    // Act
    const viewModel = controller.getViewModel()

    // Assert
    expect(viewModel.gatewayStatus).toBe("未接続")
  })

  test("gateway ありの controller は接続準備済み status を返す", () => {
    // Arrange
    const { gateway } = createGateway()
    const controller = createMasterDictionaryScreenControllerFactory(gateway)()

    // Act
    const viewModel = controller.getViewModel()

    // Assert
    expect(viewModel.gatewayStatus).toBe("接続準備済み")
  })

  test("mount は初期一覧取得に既定 filter を渡す", async () => {
    // Arrange
    const { gateway, listMasterDictionaryEntries } = createGateway()
    ;(window as Window & { runtime?: unknown }).runtime = {
      EventsOnMultiple: vi.fn(() => vi.fn())
    }
    const controller = createMasterDictionaryScreenControllerFactory(gateway)()

    // Act
    await controller.mount()

    // Assert
    expect(listMasterDictionaryEntries).toHaveBeenCalledWith({
      filters: {
        query: "",
        category: "すべて",
        page: 1,
        pageSize: 30
      }
    })
  })

  test("mount は初期選択エントリの詳細取得を呼ぶ", async () => {
    // Arrange
    const { gateway, getMasterDictionaryEntry } = createGateway()
    ;(window as Window & { runtime?: unknown }).runtime = {
      EventsOnMultiple: vi.fn(() => vi.fn())
    }
    const controller = createMasterDictionaryScreenControllerFactory(gateway)()

    // Act
    await controller.mount()

    // Assert
    expect(getMasterDictionaryEntry).toHaveBeenCalledWith({
      id: "101"
    })
  })

  test("mount は runtime listener を登録する", async () => {
    // Arrange
    const { gateway } = createGateway()
    const eventsOnMultiple = vi.fn(() => vi.fn())
    ;(window as Window & { runtime?: unknown }).runtime = {
      EventsOnMultiple: eventsOnMultiple
    }
    const controller = createMasterDictionaryScreenControllerFactory(gateway)()

    // Act
    await controller.mount()

    // Assert
    expect(eventsOnMultiple).toHaveBeenCalledTimes(2)
  })

  test("dispose は登録済み runtime listener を解除する", async () => {
    // Arrange
    const { gateway } = createGateway()
    const detach = vi.fn()
    ;(window as Window & { runtime?: unknown }).runtime = {
      EventsOnMultiple: vi.fn(() => detach)
    }
    const controller = createMasterDictionaryScreenControllerFactory(gateway)()
    await controller.mount()

    // Act
    controller.dispose()

    // Assert
    expect(detach).toHaveBeenCalledTimes(2)
  })
})
