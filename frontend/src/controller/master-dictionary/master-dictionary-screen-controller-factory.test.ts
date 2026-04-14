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

  test("gateway null は未接続 view model を返す", () => {
    const controller = createMasterDictionaryScreenControllerFactory(null)()

    expect(controller.getViewModel().gatewayStatus).toBe("未接続")
  })

  test("gateway ありの mount は runtime 購読と初期一覧取得を配線する", async () => {
    const { gateway, listMasterDictionaryEntries, getMasterDictionaryEntry } =
      createGateway()
    const detach = vi.fn()
    ;(window as Window & { runtime?: unknown }).runtime = {
      EventsOnMultiple: vi.fn(() => detach)
    }

    const controller = createMasterDictionaryScreenControllerFactory(gateway)()

    expect(controller.getViewModel().gatewayStatus).toBe("接続準備済み")

    await controller.mount()

    expect(listMasterDictionaryEntries).toHaveBeenCalledWith({
      filters: {
        query: "",
        category: "すべて",
        page: 1,
        pageSize: 30
      }
    })
    expect(getMasterDictionaryEntry).toHaveBeenCalledWith({
      id: "101"
    })

    controller.dispose()

    expect(detach).toHaveBeenCalledTimes(2)
  })
})
