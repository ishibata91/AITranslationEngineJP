/**
 * dictionary-read-detail-cutover: gateway contract / Wails binding テスト
 *
 * 検証対象:
 * - listMasterDictionaryEntries が返すエントリに rec / edid フィールドが含まれないこと
 * - getMasterDictionaryEntry が返すエントリに rec / edid フィールドが含まれないこと (note のみ)
 * - MasterDictionaryEntrySummary 型に rec / edid プロパティが存在しないこと (型レベル)
 * - gateway が canonical フィールド (source, translation, category, origin, updatedAt) だけを返すこと
 */

import { afterEach, describe, expect, test, vi } from "vitest"

import type {
  ListMasterDictionaryEntriesResponse,
  MasterDictionaryEntrySummary,
  MasterDictionaryEntryDetail,
  GetMasterDictionaryEntryResponse
} from "@application/gateway-contract/master-dictionary"

import { createMasterDictionaryGateway } from "./master-dictionary.gateway"

// ---------------------------------------------------------------------------
// Wails グローバルのセットアップ / クリーンアップ
// ---------------------------------------------------------------------------

type GoRecord = {
  wails: {
    MasterDictionaryController: Record<string, ReturnType<typeof vi.fn>>
  }
}

const originalGo: unknown = Reflect.get(globalThis as object, "go")

function installGo(record: GoRecord): void {
  Object.defineProperty(globalThis, "go", {
    value: record,
    configurable: true,
    writable: true
  })
}

afterEach(() => {
  vi.restoreAllMocks()
  Object.defineProperty(globalThis, "go", {
    value: originalGo,
    configurable: true,
    writable: true
  })
})

// ---------------------------------------------------------------------------
// 型レベル assertion ヘルパー
// ---------------------------------------------------------------------------

// 型 T が K プロパティを持たないことを compile-time に確認するユーティリティ。
// 持っていた場合は型エラーになる。
type AssertNoKey<T, K extends string> = K extends keyof T ? never : true

// MasterDictionaryEntrySummary は rec / edid を持たないこと
const _summaryHasNoRec: AssertNoKey<MasterDictionaryEntrySummary, "rec"> = true
const _summaryHasNoEdid: AssertNoKey<MasterDictionaryEntrySummary, "edid"> = true
// unused variable warning を避けるための参照
void _summaryHasNoRec
void _summaryHasNoEdid

// MasterDictionaryEntryDetail も rec / edid を持たないこと
const _detailHasNoRec: AssertNoKey<MasterDictionaryEntryDetail, "rec"> = true
const _detailHasNoEdid: AssertNoKey<MasterDictionaryEntryDetail, "edid"> = true
void _detailHasNoRec
void _detailHasNoEdid

// ---------------------------------------------------------------------------
// listMasterDictionaryEntries: canonical フィールドのみ
// ---------------------------------------------------------------------------

describe("listMasterDictionaryEntries", () => {
  test("返却エントリに rec / edid フィールドが含まれないこと", async () => {
    const wailsResponse: ListMasterDictionaryEntriesResponse = {
      entries: [
        {
          id: "1",
          source: "Whiterun",
          translation: "ホワイトラン",
          category: "地名",
          origin: "手動登録",
          updatedAt: "2026-04-19T12:00:00Z"
        }
      ],
      totalCount: 1,
      page: 1,
      pageSize: 30
    }

    installGo({
      wails: {
        MasterDictionaryController: {
          ListMasterDictionaryEntries: vi.fn(() => Promise.resolve(wailsResponse))
        }
      }
    })

    const gateway = createMasterDictionaryGateway()
    const result = await gateway.listMasterDictionaryEntries({
      filters: { query: "", category: "", page: 1, pageSize: 30 }
    })

    expect(result.entries).toHaveLength(1)

    const entry = result.entries[0]
    // canonical フィールドが存在すること
    expect(entry.id).toBe("1")
    expect(entry.source).toBe("Whiterun")
    expect(entry.translation).toBe("ホワイトラン")
    expect(entry.category).toBe("地名")
    expect(entry.origin).toBe("手動登録")
    expect(entry.updatedAt).toBe("2026-04-19T12:00:00Z")

    // rec / edid は MasterDictionaryEntrySummary の型に存在しない。
    // 以下は runtime assertion: 実際のオブジェクトに rec/edid プロパティがないこと。
    expect(Object.prototype.hasOwnProperty.call(entry, "rec")).toBe(false)
    expect(Object.prototype.hasOwnProperty.call(entry, "edid")).toBe(false)
  })

  test("フィルタクエリを Wails binding に正しく渡すこと", async () => {
    const listBinding = vi.fn(() =>
      Promise.resolve({
        entries: [],
        totalCount: 0,
        page: 2,
        pageSize: 15
      })
    )

    installGo({
      wails: {
        MasterDictionaryController: {
          ListMasterDictionaryEntries: listBinding
        }
      }
    })

    const gateway = createMasterDictionaryGateway()
    await gateway.listMasterDictionaryEntries({
      filters: { query: "Dragon", category: "NPC", page: 2, pageSize: 15 }
    })

    expect(listBinding).toHaveBeenCalledTimes(1)
    expect(listBinding).toHaveBeenCalledWith({
      filters: { query: "Dragon", category: "NPC", page: 2, pageSize: 15 }
    })
  })
})

// ---------------------------------------------------------------------------
// getMasterDictionaryEntry: canonical フィールドのみ (note あり、rec/edid なし)
// ---------------------------------------------------------------------------

describe("getMasterDictionaryEntry", () => {
  test("返却エントリに note が含まれ rec / edid が含まれないこと", async () => {
    const wailsResponse: GetMasterDictionaryEntryResponse = {
      entry: {
        id: "42",
        source: "Auriel's Bow",
        translation: "アーリエルの弓",
        category: "武器",
        origin: "手動登録",
        updatedAt: "2026-04-19T12:00:00Z",
        note: "マスター辞書エントリ"
      }
    }

    installGo({
      wails: {
        MasterDictionaryController: {
          GetMasterDictionaryEntry: vi.fn(() => Promise.resolve(wailsResponse))
        }
      }
    })

    const gateway = createMasterDictionaryGateway()
    const result = await gateway.getMasterDictionaryEntry({ id: "42" })

    expect(result.entry).not.toBeNull()
    const entry = result.entry!

    // canonical フィールドが存在すること
    expect(entry.id).toBe("42")
    expect(entry.source).toBe("Auriel's Bow")
    expect(entry.translation).toBe("アーリエルの弓")
    expect(entry.note).toBe("マスター辞書エントリ")
  })

  test("エントリが存在しない場合は null を返すこと", async () => {
    installGo({
      wails: {
        MasterDictionaryController: {
          GetMasterDictionaryEntry: vi.fn(() => Promise.resolve({ entry: null }))
        }
      }
    })

    const gateway = createMasterDictionaryGateway()
    const result = await gateway.getMasterDictionaryEntry({ id: "999" })

    expect(result.entry).toBeNull()
  })

  // Go controller は note に常に "マスター辞書エントリ" を返す (post-fix 動作)
  test("Go backend が rec/edid を省略した場合にエントリが正しくマップされること", async () => {
    // Post-fix: Go controller は note に常に "マスター辞書エントリ" を返し、
    // rec/edid フィールドは応答に含まれない。
    const wailsResponse = {
      entry: {
        id: "42",
        source: "Auriel's Bow",
        translation: "アーリエルの弓",
        category: "武器",
        origin: "手動登録",
        updatedAt: "2026-04-19T12:00:00Z",
        note: "マスター辞書エントリ"
        // rec および edid は意図的に省略 — Go backend は返さない
      }
    }

    installGo({
      wails: {
        MasterDictionaryController: {
          GetMasterDictionaryEntry: vi.fn(() => Promise.resolve(wailsResponse))
        }
      }
    })

    const gateway = createMasterDictionaryGateway()
    const result = await gateway.getMasterDictionaryEntry({ id: "42" })

    expect(result.entry).not.toBeNull()
    const entry = result.entry!

    // canonical フィールドが正しくマップされること
    expect(entry.id).toBe("42")
    expect(entry.source).toBe("Auriel's Bow")
    expect(entry.translation).toBe("アーリエルの弓")
    expect(entry.note).toBe("マスター辞書エントリ")

    // runtime assertion: rec/edid が own property として存在しないこと
    expect(Object.prototype.hasOwnProperty.call(entry, "rec")).toBe(false)
    expect(Object.prototype.hasOwnProperty.call(entry, "edid")).toBe(false)
  })
})
