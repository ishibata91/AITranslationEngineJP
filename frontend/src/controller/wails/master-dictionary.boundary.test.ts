// @vitest-environment node
/**
 * 境界結合テスト: master-dictionary (frontend 側)
 *
 * 本ファイルは BCT-MDC-001..010 の観点に従い、wails runtime を mock した上で
 * 共有 golden を mock 応答として供給し、frontend gateway / 解釈層の field 値が
 * golden と矛盾しないことを assert する。
 *
 * 両側（backend / frontend）が同一の物理 golden ファイルを参照するため、
 * 片側のみ golden を書き換えると本テスト群が失敗する検出経路が成立する。
 *
 * 対象 golden ディレクトリ: internal/apitest/testdata/boundary/master_dictionary/
 * golden 読み込み helper: boundary-golden-loader.ts
 *
 * 【片側書き換え検出の設計】
 * golden.expected を wails mock 応答として登録する。
 * assert 期待値はテストソースコードにハードコードされたリテラルで記述する。
 * この形にすることで、golden の field 値を書き換えると mock 応答が変わり、
 * ハードコードされた期待値との assert で test が落ちる。
 *
 * @vitest-environment node の理由:
 * boundary-golden-loader.ts が node:fs と node:url を使い import.meta.url を
 * file: スキームとして解決するため、Node.js 環境が必要。
 * jsdom 環境（デフォルト）では import.meta.url のスキームが一致しない。
 */

import { afterEach, describe, expect, it, vi } from "vitest"

import type {
  CreateMasterDictionaryEntryResponse,
  DeleteMasterDictionaryEntryResponse,
  GetMasterDictionaryEntryResponse,
  ImportMasterDictionaryXmlResponse,
  ListMasterDictionaryEntriesResponse,
  MasterDictionaryEntryDetail,
  MasterDictionaryEntrySummary,
  UpdateMasterDictionaryEntryResponse
} from "@application/gateway-contract/master-dictionary"

import { loadBoundaryGolden } from "./boundary-golden-loader"
import { createMasterDictionaryGateway } from "./master-dictionary.gateway"

// ---------------------------------------------------------------------------
// Wails グローバルのセットアップ / クリーンアップ
// gateway test と同じ mock 機構を流用する
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
// BCT-MDC-001: ListMasterDictionaryEntries / 正常系
// ---------------------------------------------------------------------------

describe("boundary > master-dictionary > ListMasterDictionaryEntries", () => {
  it("normal: golden の entries[0] 各 field が gateway 応答と一致すること (BCT-MDC-001)", async () => {
    // golden が定義する「backend が返すべき DTO field 値」を frontend 側でも受け取れることを証明する。
    // golden.expected を mock 応答として登録し、assert 期待値はソースコードにリテラルで固定する。
    // golden を書き換えると mock 応答が変わり、リテラル期待値との assert が失敗する。
    const golden =
      loadBoundaryGolden<ListMasterDictionaryEntriesResponse>("list_normal.golden.json")

    // Arrange: golden の expected を wails mock 応答として登録する
    installGo({
      wails: {
        MasterDictionaryController: {
          ListMasterDictionaryEntries: vi.fn(() =>
            Promise.resolve(golden.expected)
          )
        }
      }
    })

    // Act: gateway を経由して一覧を取得する
    const gateway = createMasterDictionaryGateway()
    const result = await gateway.listMasterDictionaryEntries({
      filters: { query: "", category: "", page: 1, pageSize: 30 }
    })

    // Assert: gateway が返す field 値が golden から定まるリテラル期待値と一致すること
    // 期待値リテラルは list_normal.golden.json の expected から決定論的に固定している
    expect(result.entries).toHaveLength(1)
    expect(result.totalCount).toBe(1)
    expect(result.page).toBe(1)
    expect(result.pageSize).toBe(30)

    const entry: MasterDictionaryEntrySummary = result.entries[0]
    expect(entry.id).toBe("1")
    expect(entry.source).toBe("Whiterun")
    expect(entry.translation).toBe("ホワイトラン")
    expect(entry.category).toBe("地名")
    expect(entry.origin).toBe("手動登録")
    expect(entry.updatedAt).toBe("2026-01-01T00:00:00Z")
  })

  // ---------------------------------------------------------------------------
  // BCT-MDC-002: ListMasterDictionaryEntries / 空集合
  // ---------------------------------------------------------------------------

  it("empty: golden の entries が空配列で totalCount=0 であること (BCT-MDC-002)", async () => {
    // エントリ 0 件の golden を消費し、gateway が空配列と totalCount=0 を返すことを証明する。
    const golden =
      loadBoundaryGolden<ListMasterDictionaryEntriesResponse>("list_empty.golden.json")

    installGo({
      wails: {
        MasterDictionaryController: {
          ListMasterDictionaryEntries: vi.fn(() =>
            Promise.resolve(golden.expected)
          )
        }
      }
    })

    const gateway = createMasterDictionaryGateway()
    const result = await gateway.listMasterDictionaryEntries({
      filters: { query: "", category: "", page: 1, pageSize: 30 }
    })

    // Assert: 空配列かつ totalCount=0 であること（list_empty.golden.json の expected 値）
    expect(result.entries).toHaveLength(0)
    expect(result.totalCount).toBe(0)
    expect(result.page).toBe(1)
    expect(result.pageSize).toBe(30)
  })
})

// ---------------------------------------------------------------------------
// BCT-MDC-003: GetMasterDictionaryEntry / 正常系
// ---------------------------------------------------------------------------

describe("boundary > master-dictionary > GetMasterDictionaryEntry", () => {
  it("normal: golden の entry 各 field (note 含む) が gateway 応答と一致すること (BCT-MDC-003)", async () => {
    // GetMasterDictionaryEntry の正常系。entry.note が固定文言「マスター辞書エントリ」であることを証明する。
    // get_after_create.golden.json は Create → Get 連鎖の Get 側 golden であり、
    // 単独の正常系参照の代表 case として兼用する。
    const golden =
      loadBoundaryGolden<GetMasterDictionaryEntryResponse>("get_after_create.golden.json")

    installGo({
      wails: {
        MasterDictionaryController: {
          GetMasterDictionaryEntry: vi.fn(() =>
            Promise.resolve(golden.expected)
          )
        }
      }
    })

    const gateway = createMasterDictionaryGateway()
    const result = await gateway.getMasterDictionaryEntry({ id: "1" })

    // Assert: entry が non-null で各 field が get_after_create.golden.json の期待値と一致すること
    expect(result.entry).not.toBeNull()
    const entry: MasterDictionaryEntryDetail = result.entry!
    expect(entry.id).toBe("1")
    expect(entry.source).toBe("Whiterun")
    expect(entry.translation).toBe("ホワイトラン")
    expect(entry.category).toBe("地名")
    expect(entry.origin).toBe("手動登録")
    expect(entry.updatedAt).toBe("2026-01-01T00:00:00Z")
    expect(entry.note).toBe("マスター辞書エントリ")
  })

  // ---------------------------------------------------------------------------
  // BCT-MDC-004: GetMasterDictionaryEntry / 不在（null 応答）
  // ---------------------------------------------------------------------------

  it("absent: golden の entry が null であること (BCT-MDC-004)", async () => {
    // 存在しない ID への応答が entry: null であることを証明する。
    const golden =
      loadBoundaryGolden<GetMasterDictionaryEntryResponse>("get_absent.golden.json")

    installGo({
      wails: {
        MasterDictionaryController: {
          GetMasterDictionaryEntry: vi.fn(() =>
            Promise.resolve(golden.expected)
          )
        }
      }
    })

    const gateway = createMasterDictionaryGateway()
    const result = await gateway.getMasterDictionaryEntry({ id: "999999" })

    // Assert: entry が null であること（get_absent.golden.json の expected.entry: null）
    expect(result.entry).toBeNull()
  })
})

// ---------------------------------------------------------------------------
// BCT-MDC-005: CreateMasterDictionaryEntry / 正常系
// ---------------------------------------------------------------------------

describe("boundary > master-dictionary > CreateMasterDictionaryEntry", () => {
  it("normal: golden の entry 各 field と refreshTargetId が gateway 応答と一致すること (BCT-MDC-005)", async () => {
    // 作成直後の応答 DTO field semantic を証明する。
    // entry.note が固定値であること、refreshTargetId が entry.id と同じ値であることを確認する。
    const golden =
      loadBoundaryGolden<CreateMasterDictionaryEntryResponse>("create_response.golden.json")

    installGo({
      wails: {
        MasterDictionaryController: {
          CreateMasterDictionaryEntry: vi.fn(() =>
            Promise.resolve(golden.expected)
          )
        }
      }
    })

    const gateway = createMasterDictionaryGateway()
    const result = await gateway.createMasterDictionaryEntry({
      payload: {
        source: "Whiterun",
        translation: "ホワイトラン",
        category: "地名",
        origin: "手動登録"
      }
    })

    // Assert: entry 各 field が create_response.golden.json の期待値と一致し、
    //         refreshTargetId が entry.id と同じであること
    const entry: MasterDictionaryEntryDetail = result.entry
    expect(entry.id).toBe("1")
    expect(entry.source).toBe("Whiterun")
    expect(entry.translation).toBe("ホワイトラン")
    expect(entry.category).toBe("地名")
    expect(entry.origin).toBe("手動登録")
    expect(entry.note).toBe("マスター辞書エントリ")
    expect(result.refreshTargetId).toBe("1")
    expect(result.refreshTargetId).toBe(entry.id)
  })
})

// ---------------------------------------------------------------------------
// BCT-MDC-006 / BCT-MDC-007 / BCT-MDC-008: 状態遷移 3 連鎖
// ---------------------------------------------------------------------------

describe("boundary > master-dictionary > 状態遷移 3 連鎖", () => {
  it("create→get: Create 応答と Get 応答の entry.id / entry.source が一致すること (BCT-MDC-006)", async () => {
    // Create → Get の field 値 semantic 整合を証明する。
    // Create 応答の entry.id と Get 応答の entry.id が golden 上で同一値であることを確認する。
    const createGolden =
      loadBoundaryGolden<CreateMasterDictionaryEntryResponse>("create_response.golden.json")
    const getGolden =
      loadBoundaryGolden<GetMasterDictionaryEntryResponse>("get_after_create.golden.json")

    const createMock = vi.fn(() => Promise.resolve(createGolden.expected))
    const getMock = vi.fn(() => Promise.resolve(getGolden.expected))
    installGo({
      wails: {
        MasterDictionaryController: {
          CreateMasterDictionaryEntry: createMock,
          GetMasterDictionaryEntry: getMock
        }
      }
    })

    const gateway = createMasterDictionaryGateway()

    // Act: Create を呼び出す
    const createResult = await gateway.createMasterDictionaryEntry({
      payload: {
        source: "Whiterun",
        translation: "ホワイトラン",
        category: "地名",
        origin: "手動登録"
      }
    })

    // Act: Get を呼び出す（mock 応答は get_after_create golden に差し替え済み）
    const getResult = await gateway.getMasterDictionaryEntry({
      id: createResult.entry.id
    })

    // Assert: Create 応答の id が create_response.golden.json と一致すること
    expect(createResult.entry.id).toBe("1")

    // Assert: Get 応答の field が get_after_create.golden.json と一致すること
    expect(getResult.entry).not.toBeNull()
    const getEntry = getResult.entry!
    expect(getEntry.id).toBe("1")
    expect(getEntry.source).toBe("Whiterun")
    expect(getEntry.note).toBe("マスター辞書エントリ")

    // Assert: Create 応答の id と Get 応答の id が同じであること（golden 上での整合確認）
    expect(createResult.entry.id).toBe(getEntry.id)
  })

  it("update→get: Update 応答の translation と Get 応答の translation が一致し source は不変であること (BCT-MDC-007)", async () => {
    // 更新で変化する field（translation）と変化しない field（source）の semantic を証明する。
    const updateGolden =
      loadBoundaryGolden<UpdateMasterDictionaryEntryResponse>("update_response.golden.json")
    const getGolden =
      loadBoundaryGolden<GetMasterDictionaryEntryResponse>("get_after_update.golden.json")

    const updateMock = vi.fn(() => Promise.resolve(updateGolden.expected))
    const getMock = vi.fn(() => Promise.resolve(getGolden.expected))
    installGo({
      wails: {
        MasterDictionaryController: {
          UpdateMasterDictionaryEntry: updateMock,
          GetMasterDictionaryEntry: getMock
        }
      }
    })

    const gateway = createMasterDictionaryGateway()

    // Act: Update を呼び出す
    const updateResult = await gateway.updateMasterDictionaryEntry({
      id: "1",
      payload: {
        source: "Whiterun",
        translation: "白の城下町",
        category: "地名",
        origin: "手動登録"
      }
    })

    // Act: Get を呼び出す（mock 応答は get_after_update golden に差し替え済み）
    const getResult = await gateway.getMasterDictionaryEntry({ id: "1" })

    // Assert: Update 後の translation が update_response.golden.json と一致すること
    expect(updateResult.entry.translation).toBe("白の城下町")
    // source は変化しないこと
    expect(updateResult.entry.source).toBe("Whiterun")

    // Assert: Get 応答の translation が get_after_update.golden.json と一致すること
    expect(getResult.entry).not.toBeNull()
    const getEntry = getResult.entry!
    expect(getEntry.translation).toBe("白の城下町")
    expect(getEntry.source).toBe("Whiterun")

    // Assert: Update 応答と Get 応答で translation が同じ値であること（golden 上での整合確認）
    expect(updateResult.entry.translation).toBe(getEntry.translation)
    expect(updateResult.entry.source).toBe(getEntry.source)
  })

  it("delete→get: Delete 応答の deletedId が golden と一致し、Get 応答が null 遷移すること (BCT-MDC-008)", async () => {
    // 削除後の参照可否（null 遷移）という状態遷移後 semantic を証明する。
    const deleteGolden =
      loadBoundaryGolden<DeleteMasterDictionaryEntryResponse>("delete_response.golden.json")
    const getGolden =
      loadBoundaryGolden<GetMasterDictionaryEntryResponse>("get_after_delete.golden.json")

    const deleteMock = vi.fn(() => Promise.resolve(deleteGolden.expected))
    const getMock = vi.fn(() => Promise.resolve(getGolden.expected))
    installGo({
      wails: {
        MasterDictionaryController: {
          DeleteMasterDictionaryEntry: deleteMock,
          GetMasterDictionaryEntry: getMock
        }
      }
    })

    const gateway = createMasterDictionaryGateway()

    // Act: Delete を呼び出す
    const deleteResult = await gateway.deleteMasterDictionaryEntry({ id: "1" })

    // Act: Get を呼び出す（mock 応答は get_after_delete golden に差し替え済み）
    const getResult = await gateway.getMasterDictionaryEntry({
      id: deleteResult.deletedId
    })

    // Assert: deletedId が delete_response.golden.json と一致すること
    expect(deleteResult.deletedId).toBe("1")

    // Assert: 削除後の Get 応答が null 遷移すること（get_after_delete.golden.json の expected.entry: null）
    expect(getResult.entry).toBeNull()
  })
})

// ---------------------------------------------------------------------------
// BCT-MDC-009: ImportMasterDictionaryXml / 13 種別内 REC 正常取り込み
// ---------------------------------------------------------------------------

describe("boundary > master-dictionary > ImportMasterDictionaryXml", () => {
  it("rec_in_scope: accepted=true / importedCount>0 / skippedCount=0 が golden と一致すること (BCT-MDC-009)", async () => {
    // 13 種別内 REC のみを含む XML の取り込み応答 field semantic を証明する。
    // skippedCount=0 であることを境界 DTO で確認する。
    const golden =
      loadBoundaryGolden<ImportMasterDictionaryXmlResponse>("import_rec_in_scope.golden.json")

    installGo({
      wails: {
        MasterDictionaryController: {
          ImportMasterDictionaryXml: vi.fn(() =>
            Promise.resolve(golden.expected)
          )
        }
      }
    })

    const gateway = createMasterDictionaryGateway()
    const result = await gateway.importMasterDictionaryXml({
      filePath: "/tmp/test_in_scope.xml"
    })

    // Assert: accepted / summary 各 field が import_rec_in_scope.golden.json の期待値と一致すること
    expect(result.accepted).toBe(true)
    expect(result.summary).toBeDefined()
    const summary = result.summary!
    expect(summary.importedCount).toBe(1)
    expect(summary.updatedCount).toBe(0)
    expect(summary.skippedCount).toBe(0)
  })

  // ---------------------------------------------------------------------------
  // BCT-MDC-010: ImportMasterDictionaryXml / 13 種別外 REC 取り込み（境界）
  // ---------------------------------------------------------------------------

  it("rec_out_of_scope: importedCount=0 / skippedCount>0 が golden と一致すること (BCT-MDC-010)", async () => {
    // 13 種別外 REC を含む XML を取り込んだ場合、skippedCount が増加することを証明する。
    // IsTermTarget の判定ロジックそのものは backend 単体テストの責務（Q-007）。
    // 本テストは「応答 DTO の skippedCount field に件数が乗ること」の semantic を証明する。
    const golden =
      loadBoundaryGolden<ImportMasterDictionaryXmlResponse>("import_rec_out_of_scope.golden.json")

    installGo({
      wails: {
        MasterDictionaryController: {
          ImportMasterDictionaryXml: vi.fn(() =>
            Promise.resolve(golden.expected)
          )
        }
      }
    })

    const gateway = createMasterDictionaryGateway()
    const result = await gateway.importMasterDictionaryXml({
      filePath: "/tmp/test_out_of_scope.xml"
    })

    // Assert: importedCount=0 / skippedCount>0 が import_rec_out_of_scope.golden.json の期待値と一致すること
    expect(result.accepted).toBe(true)
    expect(result.summary).toBeDefined()
    const summary = result.summary!
    expect(summary.importedCount).toBe(0)
    expect(summary.updatedCount).toBe(0)
    expect(summary.skippedCount).toBe(1)
    // 13 種別外のみなので skippedCount が 0 超えであること
    expect(summary.skippedCount).toBeGreaterThan(0)
  })
})
