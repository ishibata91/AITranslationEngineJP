import { render, screen, waitFor, within } from "@testing-library/svelte"
import userEvent from "@testing-library/user-event"
import { vi } from "vitest"

import type {
  CreateMasterDictionaryEntryRequest,
  DeleteMasterDictionaryEntryRequest,
  GetMasterDictionaryEntryRequest,
  ImportMasterDictionaryXmlRequest,
  ImportMasterDictionaryXmlResponse,
  ListMasterDictionaryEntriesRequest,
  MasterDictionaryGatewayContract,
  UpdateMasterDictionaryEntryRequest
} from "@application/gateway-contract/master-dictionary"
import { createTestMasterDictionaryScreenControllerFactory } from "../test/setup"
import App from "@ui/App.svelte"

const DASHBOARD_SHELL_PRIMARY_ROUTES = [
  {
    id: "dashboard",
    label: "ダッシュボード",
    state: "既定表示",
    description: "主要ページへの入口をまとめて確認します。"
  },
  {
    id: "master-dictionary",
    label: "マスター辞書",
    state: "準備中",
    description: "用語と訳語の基盤データを確認します。"
  },
  {
    id: "master-persona",
    label: "マスターペルソナ",
    state: "準備中",
    description: "翻訳に使うペルソナ設定を確認します。"
  },
  {
    id: "translation-management",
    label: "翻訳管理",
    state: "準備中",
    description: "準備状況と翻訳ジョブの進行をまとめて確認します。"
  },
  {
    id: "output-management",
    label: "出力管理",
    state: "準備中",
    description: "生成物と書き出し結果を確認します。"
  }
]

const DASHBOARD_ENTRY_ROUTES = DASHBOARD_SHELL_PRIMARY_ROUTES.filter(
  ({ id }) => id !== "dashboard"
)

const PLACEHOLDER_LEAD =
  "このページはまだ準備中です。上のナビゲーションまたは下の移動から別の主要ページへ進めます。"

const DASHBOARD_EXCLUDED_TEXTS = ["ジョブ一覧", "進捗サマリ"]

const IMPORT_XML_WITH_MIXED_REC = `<?xml version="1.0" encoding="utf-8"?>
<Root>
  <String>
    <REC>BOOK:FULL</REC>
    <Source>Allowed Book Source</Source>
    <Dest>許可された本の訳語</Dest>
  </String>
  <String>
    <REC>ACTI:FULL</REC>
    <Source>Denied Activator Source</Source>
    <Dest>拒否される訳語</Dest>
  </String>
  <String>
    <REC>WEAP:FULL</REC>
    <Source>Allowed Weapon Source</Source>
    <Dest>許可された武器の訳語</Dest>
  </String>
  <String>
    <REC>NPC_:FULL</REC>
    <Source>Skipped Empty Dest Source</Source>
    <Dest></Dest>
  </String>
</Root>`

function createXmlFile(contents: string, name = "mixed-rec.xml"): File {
  const file = new File([contents], name, { type: "text/xml" })
  Object.defineProperty(file, "text", {
    value: () => Promise.resolve(contents),
    configurable: true
  })
  return file
}

const MASTER_DICTIONARY_IMPORT_PROGRESS_EVENT =
  "master-dictionary:import-progress"
const MASTER_DICTIONARY_IMPORT_COMPLETED_EVENT =
  "master-dictionary:import-completed"

type RuntimeEventPayload = Record<string, unknown>

function installRuntimeEventBridgeMock() {
  const callbacks = new Map<string, Array<(...args: unknown[]) => void>>()
  const progressEvents: RuntimeEventPayload[] = []
  const completionEvents: RuntimeEventPayload[] = []
  const previousRuntime = (window as Window & { runtime?: unknown }).runtime

  Object.defineProperty(window, "runtime", {
    configurable: true,
    writable: true,
    value: {
      EventsOnMultiple: (
        eventName: string,
        callback: (...args: unknown[]) => void
      ) => {
        const listeners = callbacks.get(eventName) ?? []
        listeners.push(callback)
        callbacks.set(eventName, listeners)
        return () => {
          const current = callbacks.get(eventName) ?? []
          callbacks.set(
            eventName,
            current.filter((candidate) => candidate !== callback)
          )
        }
      }
    }
  })

  return {
    emitImportProgress: (payload: RuntimeEventPayload) => {
      progressEvents.push(payload)
      for (const callback of callbacks.get(
        MASTER_DICTIONARY_IMPORT_PROGRESS_EVENT
      ) ?? []) {
        callback(payload)
      }
    },
    emitImportCompleted: (payload: RuntimeEventPayload) => {
      completionEvents.push(payload)
      for (const callback of callbacks.get(
        MASTER_DICTIONARY_IMPORT_COMPLETED_EVENT
      ) ?? []) {
        callback(payload)
      }
    },
    progressEventCount: () => progressEvents.length,
    completionEventCount: () => completionEvents.length,
    restore: () => {
      if (previousRuntime === undefined) {
        delete (window as Window & { runtime?: unknown }).runtime
        return
      }
      ;(window as Window & { runtime?: unknown }).runtime = previousRuntime
    }
  }
}

function createPagedSeedEntries(total: number): TestEntry[] {
  return Array.from({ length: total }, (_, index) => {
    const id = String(1000 + index)
    const sequence = String(index + 1).padStart(2, "0")
    const isLocation = index % 2 === 1
    return {
      id,
      source: `Source-${sequence}`,
      translation: `訳語-${sequence}`,
      category: isLocation ? "地名" : "固有名詞",
      origin: "初期データ",
      updatedAt: `2026-04-${String((index % 28) + 1).padStart(2, "0")} 00:00`,
      note: `REC: ${isLocation ? "LCTN:FULL" : "NPC_:FULL"} / EDID: Seed${sequence}`,
      rec: isLocation ? "LCTN:FULL" : "NPC_:FULL",
      edid: `Seed${sequence}`
    }
  })
}

interface TestEntry {
  id: string
  source: string
  translation: string
  category: string
  origin: string
  updatedAt: string
  note: string
  rec: string
  edid: string
}

interface TestGateway extends MasterDictionaryGatewayContract {
  __mocks: {
    importMasterDictionaryXml: ReturnType<typeof vi.fn>
  }
}

function createTestGateway(seedEntries?: TestEntry[]): TestGateway {
  let entries = seedEntries ?? [
    {
      id: "101",
      source: "Dragon Priest",
      translation: "ドラゴン・プリースト",
      category: "固有名詞",
      origin: "初期データ",
      updatedAt: "2026-01-01 00:00",
      note: "REC: NPC_:FULL / EDID: SeedDragonPriest",
      rec: "NPC_:FULL",
      edid: "SeedDragonPriest"
    },
    {
      id: "102",
      source: "The Reach",
      translation: "リーチ地方",
      category: "地名",
      origin: "初期データ",
      updatedAt: "2026-01-02 00:00",
      note: "REC: LCTN:FULL / EDID: SeedReach",
      rec: "LCTN:FULL",
      edid: "SeedReach"
    }
  ]

  let nextId = 1000
  const importEntries: TestEntry[] = [
    {
      id: "201",
      source: "Allowed Book Source",
      translation: "許可された本の訳語",
      category: "書籍",
      origin: "XML取込",
      updatedAt: "2026-04-12 00:00",
      note: "REC: BOOK:FULL / EDID: ImportBook",
      rec: "BOOK:FULL",
      edid: "ImportBook"
    },
    {
      id: "202",
      source: "Allowed Weapon Source",
      translation: "許可された武器の訳語",
      category: "装備",
      origin: "XML取込",
      updatedAt: "2026-04-12 00:01",
      note: "REC: WEAP:FULL / EDID: ImportWeapon",
      rec: "WEAP:FULL",
      edid: "ImportWeapon"
    }
  ]

  const filterEntries = (query: string, category: string): TestEntry[] => {
    const normalizedQuery = query.trim().toLowerCase()
    return entries.filter((entry) => {
      if (
        category !== "" &&
        category !== "すべて" &&
        entry.category !== category
      ) {
        return false
      }
      if (normalizedQuery === "") {
        return true
      }
      return [entry.source, entry.translation, entry.id]
        .join(" ")
        .toLowerCase()
        .includes(normalizedQuery)
    })
  }

  const toPageState = (
    refresh:
      | { query: string; category: string; page: number; pageSize: number }
      | undefined,
    preferredId: string | null
  ) => {
    const query = refresh?.query ?? ""
    const category = refresh?.category ?? ""
    const requestedPage = Math.max(refresh?.page ?? 1, 1)
    const pageSize = Math.max(refresh?.pageSize ?? 30, 1)

    const filtered = filterEntries(query, category)
    const maxPage = Math.max(1, Math.ceil(filtered.length / pageSize))
    const resolvedPage = Math.min(requestedPage, maxPage)
    const start = (resolvedPage - 1) * pageSize
    const pageItems = filtered.slice(start, start + pageSize)

    const selectedCandidate =
      (preferredId
        ? pageItems.find((entry) => entry.id === preferredId)
        : null) ??
      pageItems[0] ??
      null

    return {
      items: pageItems.map((entry) => ({
        id: Number.parseInt(entry.id, 10),
        source: entry.source,
        translation: entry.translation,
        category: entry.category,
        origin: entry.origin,
        rec: entry.rec,
        edid: entry.edid,
        updatedAt: entry.updatedAt
      })),
      totalCount: filtered.length,
      page: resolvedPage,
      pageSize,
      selectedId: selectedCandidate
        ? Number.parseInt(selectedCandidate.id, 10)
        : undefined
    }
  }

  const listMasterDictionaryEntries = vi.fn(
    (request: ListMasterDictionaryEntriesRequest) => {
      const filtered = filterEntries(
        request.filters.query,
        request.filters.category
      )
      const requestPage = Math.max(request.filters.page, 1)
      const pageSize = request.filters.pageSize
      const start = (requestPage - 1) * pageSize
      const pageItems = filtered.slice(start, start + pageSize)
      return Promise.resolve({
        entries: pageItems.map((entry) => ({
          id: entry.id,
          source: entry.source,
          translation: entry.translation,
          category: entry.category,
          origin: entry.origin,
          updatedAt: entry.updatedAt
        })),
        totalCount: filtered.length,
        page: requestPage,
        pageSize
      })
    }
  )

  const getMasterDictionaryEntry = vi.fn(
    (request: GetMasterDictionaryEntryRequest) => {
      const entry =
        entries.find((candidate) => candidate.id === request.id) ?? null
      return Promise.resolve({ entry })
    }
  )

  const createMasterDictionaryEntry = vi.fn(
    (request: CreateMasterDictionaryEntryRequest) => {
      const entry: TestEntry = {
        id: String(nextId++),
        source: request.payload.source,
        translation: request.payload.translation,
        category: request.payload.category,
        origin: request.payload.origin,
        updatedAt: "2026-04-12 01:00",
        note: "マスター辞書エントリ",
        rec: "",
        edid: ""
      }
      entries = [entry, ...entries]
      return Promise.resolve({
        entry,
        refreshTargetId: entry.id,
        page: toPageState(request.refresh, entry.id)
      })
    }
  )

  const updateMasterDictionaryEntry = vi.fn(
    (request: UpdateMasterDictionaryEntryRequest) => {
      entries = entries.map((entry) =>
        entry.id === request.id
          ? {
              ...entry,
              source: request.payload.source,
              translation: request.payload.translation,
              category: request.payload.category,
              origin: request.payload.origin,
              updatedAt: "2026-04-12 01:30"
            }
          : entry
      )
      const entry = entries.find((candidate) => candidate.id === request.id)
      if (!entry) {
        return Promise.reject(new Error("entry not found"))
      }
      return Promise.resolve({
        entry,
        refreshTargetId: request.id,
        page: toPageState(request.refresh, request.id)
      })
    }
  )

  const deleteMasterDictionaryEntry = vi.fn(
    (request: DeleteMasterDictionaryEntryRequest) => {
      entries = entries.filter((entry) => entry.id !== request.id)
      const page = toPageState(request.refresh, null)
      return Promise.resolve({
        deletedId: request.id,
        nextSelectedId: page.selectedId ? String(page.selectedId) : null,
        page
      })
    }
  )

  const importMasterDictionaryXml = vi.fn(
    (request: ImportMasterDictionaryXmlRequest) => {
      if (!request.filePath || !request.fileReference) {
        return Promise.reject(new Error("ファイル参照が不足しています。"))
      }
      entries = [
        ...importEntries,
        ...entries.filter((entry) => entry.origin !== "XML取込")
      ]
      const page = toPageState(request.refresh, importEntries[0]?.id ?? null)
      return Promise.resolve({
        accepted: true,
        summary: {
          filePath: request.fileReference ?? request.filePath,
          fileName: request.filePath,
          importedCount: importEntries.length,
          updatedCount: 0,
          skippedCount: 2,
          selectedRec: ["BOOK:FULL", "WEAP:FULL"],
          lastEntryId: Number.parseInt(importEntries[0]?.id ?? "0", 10)
        },
        page
      })
    }
  )

  return {
    listMasterDictionaryEntries,
    getMasterDictionaryEntry,
    createMasterDictionaryEntry,
    updateMasterDictionaryEntry,
    deleteMasterDictionaryEntry,
    importMasterDictionaryXml,
    __mocks: {
      importMasterDictionaryXml
    }
  }
}

function renderAppWithGateway(
  gateway: MasterDictionaryGatewayContract | null = null
) {
  return render(App, {
    props: {
      createMasterDictionaryScreenController:
        createTestMasterDictionaryScreenControllerFactory(gateway)
    }
  })
}

function expectRouteLinksMatch(
  links: HTMLElement[],
  expectedRoutes = DASHBOARD_SHELL_PRIMARY_ROUTES
): void {
  expect(links).toHaveLength(expectedRoutes.length)

  const actual = links.map((link) => ({
    href: link.getAttribute("href")
  }))
  const expected = expectedRoutes.map((route) => ({
    href: `#${route.id}`
  }))

  expect(actual).toEqual(expected)
}

describe("App dashboard shell", () => {
  beforeEach(() => {
    window.history.replaceState(null, "", "#")
  })

  test("SCN-DAS-001: 起動時にダッシュボードを既定表示する", () => {
    renderAppWithGateway()

    expect(
      screen.getByRole("heading", { name: "ダッシュボード" })
    ).toBeInTheDocument()
    expect(window.location.hash).toBe("#dashboard")
    for (const excludedText of DASHBOARD_EXCLUDED_TEXTS) {
      expect(screen.queryByText(excludedText)).not.toBeInTheDocument()
    }
  })

  test("invalid hash は dashboard に正規化される", () => {
    window.history.replaceState(null, "", "#not-approved-route")

    renderAppWithGateway()

    expect(
      screen.getByRole("heading", { name: "ダッシュボード" })
    ).toBeInTheDocument()
    expect(window.location.hash).toBe("#dashboard")
  })

  test("SCN-DAS-002: グローバルナビゲーションから主要 5 ページへ遷移できる", async () => {
    const user = userEvent.setup()
    renderAppWithGateway()

    const globalNavigation = screen.getByRole("navigation", {
      name: "グローバルナビゲーション"
    })

    for (const route of DASHBOARD_SHELL_PRIMARY_ROUTES) {
      await user.click(
        within(globalNavigation).getByRole("link", { name: route.label })
      )
      expect(
        screen.getByRole("heading", { level: 1, name: route.label })
      ).toBeInTheDocument()
      expect(window.location.hash).toBe(`#${route.id}`)
    }
  })

  test("SCN-DAS-003: ダッシュボード入口カードから主要 5 ページへ遷移できる", async () => {
    const user = userEvent.setup()
    renderAppWithGateway()

    const globalNavigation = screen.getByRole("navigation", {
      name: "グローバルナビゲーション"
    })

    for (const route of DASHBOARD_ENTRY_ROUTES) {
      await user.click(
        within(globalNavigation).getByRole("link", { name: "ダッシュボード" })
      )

      const dashboardCardSectionHeading = screen.getByRole("heading", {
        name: "作業を選ぶ"
      })
      const dashboardCardSection =
        dashboardCardSectionHeading.closest("section")

      if (!dashboardCardSection) {
        throw new Error("ダッシュボード入口カードのセクションが見つかりません")
      }

      const cardHeading = within(dashboardCardSection).getByRole("heading", {
        level: 3,
        name: route.label
      })
      const cardLink = cardHeading.closest("a")

      if (!cardLink) {
        throw new Error(
          `ダッシュボード入口カードのリンクが見つかりません: ${route.label}`
        )
      }

      await user.click(cardLink)
      expect(
        screen.getByRole("heading", { level: 1, name: route.label })
      ).toBeInTheDocument()
      expect(window.location.hash).toBe(`#${route.id}`)
    }
  })

  test("主要導線は承認済み 5 ルートに固定される", () => {
    renderAppWithGateway()

    const globalNavigation = screen.getByRole("navigation", {
      name: "グローバルナビゲーション"
    })
    const globalLinks = within(globalNavigation).getAllByRole("link")
    expectRouteLinksMatch(globalLinks)

    const dashboardCardSectionHeading = screen.getByRole("heading", {
      name: "作業を選ぶ"
    })
    const dashboardCardSection = dashboardCardSectionHeading.closest("section")

    if (!dashboardCardSection) {
      throw new Error("ダッシュボード入口カードのセクションが見つかりません")
    }

    const dashboardLinks = within(dashboardCardSection).getAllByRole("link")
    expectRouteLinksMatch(dashboardLinks, DASHBOARD_ENTRY_ROUTES)

    expect(
      within(dashboardCardSection).queryByRole("link", {
        name: "ダッシュボード"
      })
    ).not.toBeInTheDocument()

    for (const route of DASHBOARD_ENTRY_ROUTES) {
      expect(
        within(dashboardCardSection).getAllByText(route.state).length
      ).toBeGreaterThan(0)
      expect(
        within(dashboardCardSection).getByText(route.description)
      ).toBeInTheDocument()
    }
  })

  test("mobile nav トグルの表示契約を持つ", () => {
    renderAppWithGateway()

    const mobileToggle = screen.getByRole("button", { name: "主要ページ" })
    expect(mobileToggle).toHaveAttribute("aria-controls", "globalNav")
    expect(mobileToggle).toHaveAttribute("aria-expanded", "false")
  })

  test("非ダッシュボード遷移先ではダッシュボード専用領域を表示しない", async () => {
    const user = userEvent.setup()
    renderAppWithGateway()

    const globalNavigation = screen.getByRole("navigation", {
      name: "グローバルナビゲーション"
    })

    for (const route of DASHBOARD_SHELL_PRIMARY_ROUTES.filter(
      ({ id }) => id !== "dashboard"
    )) {
      await user.click(
        within(globalNavigation).getByRole("link", { name: route.label })
      )
      expect(
        screen.getByRole("heading", { level: 1, name: route.label })
      ).toBeInTheDocument()
      expect(
        screen.queryByRole("heading", { name: "作業を選ぶ" })
      ).not.toBeInTheDocument()

      if (route.id === "master-dictionary") {
        expect(
          screen.getByRole("heading", { level: 3, name: "辞書一覧" })
        ).toBeInTheDocument()
        expect(screen.queryByText(PLACEHOLDER_LEAD)).not.toBeInTheDocument()
        continue
      }

      expect(screen.getByText(PLACEHOLDER_LEAD)).toBeInTheDocument()
    }
  })

  test("マスター辞書ルートは contract shell の主要IDを表示する", async () => {
    const user = userEvent.setup()
    renderAppWithGateway(createTestGateway())

    const globalNavigation = screen.getByRole("navigation", {
      name: "グローバルナビゲーション"
    })

    await user.click(
      within(globalNavigation).getByRole("link", { name: "マスター辞書" })
    )

    expect(
      screen.getByRole("heading", { level: 3, name: "辞書一覧" })
    ).toBeInTheDocument()
    expect(document.querySelector("#xmlFileInput")).toBeInTheDocument()
    expect(document.querySelector("#listStack")).toBeInTheDocument()
    expect(document.querySelector("#detailTitle")).toBeInTheDocument()
    expect(document.querySelector("#createButton")).toBeInTheDocument()
    await waitFor(() => {
      expect(document.querySelector("#detailTitle")).toHaveTextContent(
        "Dragon Priest"
      )
    })
  })

  test("SCN-MDM-001: 一覧初期表示で30件ページングと詳細同期を維持する", async () => {
    const user = userEvent.setup()
    const gateway = createTestGateway(createPagedSeedEntries(31))

    renderAppWithGateway(gateway)

    const globalNavigation = screen.getByRole("navigation", {
      name: "グローバルナビゲーション"
    })
    await user.click(
      within(globalNavigation).getByRole("link", { name: "マスター辞書" })
    )

    await waitFor(() => {
      expect(screen.getByText("1 - 30 件を表示")).toBeInTheDocument()
      expect(document.querySelector("#detailTitle")).toHaveTextContent(
        "Source-01"
      )
      expect(document.querySelector("#detailTranslation")).toHaveTextContent(
        "訳語-01"
      )
    })

    await user.click(screen.getByRole("button", { name: /訳語-02/ }))

    await waitFor(() => {
      expect(document.querySelector("#detailTitle")).toHaveTextContent(
        "Source-02"
      )
      expect(document.querySelector("#detailTranslation")).toHaveTextContent(
        "訳語-02"
      )
    })
  })

  test("SCN-MDM-002/007: 検索とカテゴリ絞り込みを適用し、0件状態から復帰できる", async () => {
    const user = userEvent.setup()
    renderAppWithGateway(createTestGateway())

    const globalNavigation = screen.getByRole("navigation", {
      name: "グローバルナビゲーション"
    })
    await user.click(
      within(globalNavigation).getByRole("link", { name: "マスター辞書" })
    )

    await user.type(screen.getByRole("searchbox", { name: "検索" }), "Reach")
    await waitFor(() => {
      const listStack = document.querySelector("#listStack")
      if (!(listStack instanceof HTMLElement)) {
        throw new Error("listStack が見つかりません")
      }
      expect(within(listStack).getByText("リーチ地方")).toBeInTheDocument()
      expect(
        within(listStack).queryByText("ドラゴン・プリースト")
      ).not.toBeInTheDocument()
      expect(document.querySelector("#detailTitle")).toHaveTextContent(
        "The Reach"
      )
    })

    await user.selectOptions(
      screen.getByRole("combobox", { name: "カテゴリ" }),
      "地名"
    )
    await waitFor(() => {
      const listStack = document.querySelector("#listStack")
      if (!(listStack instanceof HTMLElement)) {
        throw new Error("listStack が見つかりません")
      }
      expect(within(listStack).getByText("リーチ地方")).toBeInTheDocument()
      expect(document.querySelector("#selectionStatus")).toHaveTextContent(
        "The Reach"
      )
    })

    const searchBox = screen.getByRole("searchbox", { name: "検索" })
    await user.clear(searchBox)
    await user.type(searchBox, "not-found-entry")
    await waitFor(() => {
      expect(
        screen.getByText("一致するエントリがありません")
      ).toBeInTheDocument()
      expect(document.querySelector("#detailTitle")).toHaveTextContent(
        "表示できるエントリがありません"
      )
    })

    await user.clear(searchBox)
    await user.selectOptions(
      screen.getByRole("combobox", { name: "カテゴリ" }),
      "すべて"
    )
    await waitFor(() => {
      const listStack = document.querySelector("#listStack")
      if (!(listStack instanceof HTMLElement)) {
        throw new Error("listStack が見つかりません")
      }
      expect(
        within(listStack).getByText("ドラゴン・プリースト")
      ).toBeInTheDocument()
      expect(document.querySelector("#detailTitle")).not.toHaveTextContent(
        "表示できるエントリがありません"
      )
    })
  })

  test("SCN-MDM-003/004/010: 新規登録と更新モーダルを開閉し、保存後は同一ページで一覧と詳細へ反映する", async () => {
    const user = userEvent.setup()
    renderAppWithGateway(createTestGateway())

    const globalNavigation = screen.getByRole("navigation", {
      name: "グローバルナビゲーション"
    })
    await user.click(
      within(globalNavigation).getByRole("link", { name: "マスター辞書" })
    )

    await waitFor(() => {
      expect(document.querySelector("#detailTitle")).toHaveTextContent(
        "Dragon Priest"
      )
    })

    await user.click(screen.getByRole("button", { name: "新規登録" }))
    expect(
      screen.getByRole("heading", { name: "新規登録" })
    ).toBeInTheDocument()
    await user.click(screen.getByRole("button", { name: "閉じる" }))
    expect(
      screen.queryByRole("heading", { name: "新規登録" })
    ).not.toBeInTheDocument()
    expect(document.querySelector("#detailTitle")).toHaveTextContent(
      "Dragon Priest"
    )

    await user.click(screen.getByRole("button", { name: "更新" }))
    expect(screen.getByRole("heading", { name: "更新" })).toBeInTheDocument()
    await user.click(screen.getByRole("button", { name: "閉じる" }))
    expect(
      screen.queryByRole("heading", { name: "更新" })
    ).not.toBeInTheDocument()
    expect(document.querySelector("#detailTitle")).toHaveTextContent(
      "Dragon Priest"
    )

    await user.click(screen.getByRole("button", { name: "新規登録" }))
    const editModal = document.querySelector("#editModal")
    if (!(editModal instanceof HTMLElement)) {
      throw new Error("editModal が見つかりません")
    }
    await user.type(within(editModal).getByLabelText("原文"), "New Source")
    await user.type(within(editModal).getByLabelText("訳語"), "新規訳語")
    await user.selectOptions(
      within(editModal).getByLabelText("カテゴリ"),
      "地名"
    )
    await user.selectOptions(
      within(editModal).getByLabelText("由来"),
      "確認待ち"
    )
    await user.click(
      within(editModal).getByRole("button", { name: "保存する" })
    )

    await waitFor(() => {
      const listStack = document.querySelector("#listStack")
      if (!(listStack instanceof HTMLElement)) {
        throw new Error("listStack が見つかりません")
      }
      expect(within(listStack).getByText("新規訳語")).toBeInTheDocument()
      expect(document.querySelector("#detailTitle")).toHaveTextContent(
        "New Source"
      )
      expect(document.querySelector("#detailTranslation")).toHaveTextContent(
        "新規訳語"
      )
    })

    await user.click(screen.getByRole("button", { name: "更新" }))
    const editModalForUpdate = document.querySelector("#editModal")
    if (!(editModalForUpdate instanceof HTMLElement)) {
      throw new Error("editModal が見つかりません")
    }
    const translationInput = within(editModalForUpdate).getByLabelText("訳語")
    await user.clear(translationInput)
    await user.type(translationInput, "更新後訳語")
    await user.click(
      within(editModalForUpdate).getByRole("button", { name: "保存する" })
    )

    await waitFor(() => {
      const listStack = document.querySelector("#listStack")
      if (!(listStack instanceof HTMLElement)) {
        throw new Error("listStack が見つかりません")
      }
      expect(within(listStack).getByText("更新後訳語")).toBeInTheDocument()
      expect(document.querySelector("#detailTranslation")).toHaveTextContent(
        "更新後訳語"
      )
      expect(window.location.hash).toBe("#master-dictionary")
    })
  })

  test("SCN-MDM-005/011: 削除確認モーダルで削除を確定し、同一ページで詳細を次対象または空状態へ切替える", async () => {
    const user = userEvent.setup()
    renderAppWithGateway(createTestGateway())

    const globalNavigation = screen.getByRole("navigation", {
      name: "グローバルナビゲーション"
    })
    await user.click(
      within(globalNavigation).getByRole("link", { name: "マスター辞書" })
    )

    await waitFor(() => {
      expect(document.querySelector("#detailTitle")).toHaveTextContent(
        "Dragon Priest"
      )
    })

    await user.click(screen.getByRole("button", { name: "削除" }))
    expect(
      screen.getByRole("heading", { name: "削除の確認" })
    ).toBeInTheDocument()
    expect(document.querySelector("#deleteTargetTitle")).toHaveTextContent(
      "Dragon Priest"
    )
    await user.click(screen.getByRole("button", { name: "削除する" }))

    await waitFor(() => {
      expect(screen.queryByText("ドラゴン・プリースト")).not.toBeInTheDocument()
      expect(document.querySelector("#detailTitle")).toHaveTextContent(
        "The Reach"
      )
      expect(window.location.hash).toBe("#master-dictionary")
    })

    await user.click(screen.getByRole("button", { name: "削除" }))
    await user.click(screen.getByRole("button", { name: "削除する" }))

    await waitFor(() => {
      expect(
        screen.getByText("一致するエントリがありません")
      ).toBeInTheDocument()
      expect(document.querySelector("#detailTitle")).toHaveTextContent(
        "表示できるエントリがありません"
      )
      expect(
        screen.getByRole("heading", { name: "取り込み導線" })
      ).toBeInTheDocument()
      expect(
        screen.getByRole("heading", { name: "辞書一覧" })
      ).toBeInTheDocument()
      expect(screen.getByRole("heading", { name: "詳細" })).toBeInTheDocument()
    })
  })

  test("SCN-MDM-008: XML未選択時は取込バーと開始操作を表示せず、選択後のみ表示する", async () => {
    const user = userEvent.setup()
    renderAppWithGateway(createTestGateway())

    const globalNavigation = screen.getByRole("navigation", {
      name: "グローバルナビゲーション"
    })
    await user.click(
      within(globalNavigation).getByRole("link", { name: "マスター辞書" })
    )

    const importBar = document.querySelector("#importBar")
    const xmlInput = document.querySelector("#xmlFileInput")
    if (!(importBar instanceof HTMLDivElement)) {
      throw new Error("importBar が見つかりません")
    }
    if (!(xmlInput instanceof HTMLInputElement)) {
      throw new Error("xmlFileInput が見つかりません")
    }

    expect(importBar).toHaveAttribute("hidden")
    expect(
      screen.queryByRole("button", { name: "この XML を取り込む" })
    ).not.toBeInTheDocument()

    const xmlFile = createXmlFile(
      IMPORT_XML_WITH_MIXED_REC,
      "Dawnguard_english_japanese.xml"
    )
    await user.upload(xmlInput, xmlFile)

    expect(importBar).not.toHaveAttribute("hidden")
    expect(
      screen.getByRole("button", { name: "この XML を取り込む" })
    ).toBeInTheDocument()
    expect(screen.getByText("取込待ち")).toBeInTheDocument()
    expect(document.querySelector("#importProgressFill")).toHaveAttribute(
      "style",
      "width: 0%;"
    )

    await user.click(screen.getByRole("button", { name: "選び直す" }))

    expect(importBar).toHaveAttribute("hidden")
    expect(
      screen.queryByRole("button", { name: "この XML を取り込む" })
    ).not.toBeInTheDocument()
  })

  test("path が空文字の file でも fileReference は file.name へ fallback して取込開始する", async () => {
    const gateway = createTestGateway()
    const user = userEvent.setup()
    renderAppWithGateway(gateway)

    const globalNavigation = screen.getByRole("navigation", {
      name: "グローバルナビゲーション"
    })
    await user.click(
      within(globalNavigation).getByRole("link", { name: "マスター辞書" })
    )

    const xmlInput = document.querySelector("#xmlFileInput")
    if (!(xmlInput instanceof HTMLInputElement)) {
      throw new Error("xmlFileInput が見つかりません")
    }

    const xmlFile = createXmlFile(
      IMPORT_XML_WITH_MIXED_REC,
      "Dawnguard_english_japanese.xml"
    )
    Object.defineProperty(xmlFile, "path", {
      value: "",
      configurable: true
    })
    Object.defineProperty(xmlFile, "webkitRelativePath", {
      value: "",
      configurable: true
    })
    await user.upload(xmlInput, xmlFile)
    await user.click(
      screen.getByRole("button", { name: "この XML を取り込む" })
    )

    await waitFor(() => {
      expect(gateway.__mocks.importMasterDictionaryXml).toHaveBeenCalledTimes(1)
    })
    expect(gateway.__mocks.importMasterDictionaryXml).toHaveBeenCalledWith({
      filePath: "Dawnguard_english_japanese.xml",
      fileReference: "Dawnguard_english_japanese.xml",
      refresh: {
        query: "",
        category: "",
        page: 1,
        pageSize: 30
      }
    })
  })

  test("SCN-MDM-006/009: XML取込完了時にbackend応答で再同期し、検索とカテゴリを初期化する", async () => {
    const gateway = createTestGateway()
    const user = userEvent.setup()

    renderAppWithGateway(gateway)

    const globalNavigation = screen.getByRole("navigation", {
      name: "グローバルナビゲーション"
    })
    await user.click(
      within(globalNavigation).getByRole("link", { name: "マスター辞書" })
    )

    await user.type(screen.getByRole("searchbox", { name: "検索" }), "Dragon")
    await user.selectOptions(
      screen.getByRole("combobox", { name: "カテゴリ" }),
      "固有名詞"
    )

    const xmlInput = document.querySelector("#xmlFileInput")
    if (!(xmlInput instanceof HTMLInputElement)) {
      throw new Error("xmlFileInput が見つかりません")
    }

    const xmlFile = createXmlFile(
      IMPORT_XML_WITH_MIXED_REC,
      "Dawnguard_english_japanese.xml"
    )
    await user.upload(xmlInput, xmlFile)
    await user.click(
      screen.getByRole("button", { name: "この XML を取り込む" })
    )

    await waitFor(() => {
      expect(screen.getByText("完了")).toBeInTheDocument()
      expect(screen.getByRole("searchbox", { name: "検索" })).toHaveValue("")
      expect(screen.getByRole("combobox", { name: "カテゴリ" })).toHaveValue(
        "すべて"
      )
    })

    expect(
      screen.queryByText("Denied Activator Source")
    ).not.toBeInTheDocument()
    expect(gateway.__mocks.importMasterDictionaryXml).toHaveBeenCalledWith({
      filePath: "Dawnguard_english_japanese.xml",
      fileReference: "Dawnguard_english_japanese.xml",
      refresh: {
        query: "",
        category: "",
        page: 1,
        pageSize: 30
      }
    })
  })

  test("import 開始表示を描画してから XML import binding を呼ぶ", async () => {
    const gateway = createTestGateway()
    const baseImport = gateway.importMasterDictionaryXml.bind(gateway)
    const importDeferred: {
      resolve: (value: ImportMasterDictionaryXmlResponse) => void
    } = {
      resolve: () => {
        throw new Error("import resolve handler が見つかりません")
      }
    }
    const importGatewaySpy = vi.fn(
      () =>
        new Promise<ImportMasterDictionaryXmlResponse>((resolve) => {
          importDeferred.resolve = resolve
        })
    )
    gateway.importMasterDictionaryXml = importGatewaySpy

    const user = userEvent.setup()

    renderAppWithGateway(gateway)

    const globalNavigation = screen.getByRole("navigation", {
      name: "グローバルナビゲーション"
    })
    await user.click(
      within(globalNavigation).getByRole("link", { name: "マスター辞書" })
    )

    const xmlInput = document.querySelector("#xmlFileInput")
    if (!(xmlInput instanceof HTMLInputElement)) {
      throw new Error("xmlFileInput が見つかりません")
    }

    const xmlFile = createXmlFile(
      IMPORT_XML_WITH_MIXED_REC,
      "Dawnguard_english_japanese.xml"
    )
    await user.upload(xmlInput, xmlFile)
    await user.click(
      screen.getByRole("button", { name: "この XML を取り込む" })
    )

    expect(document.querySelector("#importStatusValue")).toHaveTextContent(
      "取込中"
    )

    await waitFor(() => {
      expect(importGatewaySpy).toHaveBeenCalledTimes(1)
    })

    importDeferred.resolve(
      await baseImport({
        filePath: "Dawnguard_english_japanese.xml",
        fileReference: "Dawnguard_english_japanese.xml",
        refresh: {
          query: "",
          category: "",
          page: 1,
          pageSize: 30
        }
      })
    )

    await waitFor(() => {
      expect(screen.getByText("完了")).toBeInTheDocument()
    })
  })

  test("runtime completion event 到達まで import 完了へ遷移しない", async () => {
    const gateway = createTestGateway()
    const capturedImportResponse: {
      value: ImportMasterDictionaryXmlResponse | null
    } = {
      value: null
    }
    const baseImport = gateway.importMasterDictionaryXml.bind(gateway)
    const importSpy = vi.fn(
      async (request: ImportMasterDictionaryXmlRequest) => {
        const payload = await baseImport(request)
        capturedImportResponse.value = payload
        return payload
      }
    )
    gateway.importMasterDictionaryXml = importSpy

    const runtimeBridge = installRuntimeEventBridgeMock()
    const user = userEvent.setup()

    try {
      renderAppWithGateway(gateway)

      const globalNavigation = screen.getByRole("navigation", {
        name: "グローバルナビゲーション"
      })
      await user.click(
        within(globalNavigation).getByRole("link", { name: "マスター辞書" })
      )

      const xmlInput = document.querySelector("#xmlFileInput")
      if (!(xmlInput instanceof HTMLInputElement)) {
        throw new Error("xmlFileInput が見つかりません")
      }

      const xmlFile = createXmlFile(
        IMPORT_XML_WITH_MIXED_REC,
        "Dawnguard_english_japanese.xml"
      )
      await user.upload(xmlInput, xmlFile)
      await user.click(
        screen.getByRole("button", { name: "この XML を取り込む" })
      )

      await waitFor(() => {
        expect(importSpy).toHaveBeenCalledTimes(1)
      })

      expect(document.querySelector("#importResult")).toHaveAttribute("hidden")
      expect(document.querySelector("#importStatusValue")).toHaveTextContent(
        "取込中"
      )

      runtimeBridge.emitImportProgress({ progress: 78 })
      await waitFor(() => {
        expect(document.querySelector("#importStatusValue")).toHaveTextContent(
          "取込中"
        )
        expect(runtimeBridge.progressEventCount()).toBe(1)
      })

      const responsePayload = capturedImportResponse.value
      if (!responsePayload) {
        throw new Error("import 応答 payload を取得できませんでした")
      }

      runtimeBridge.emitImportCompleted({
        page: responsePayload.page,
        summary: responsePayload.summary
      })

      await waitFor(() => {
        expect(screen.getByText("完了")).toBeInTheDocument()
        expect(screen.getByRole("searchbox", { name: "検索" })).toHaveValue("")
        expect(screen.getByRole("combobox", { name: "カテゴリ" })).toHaveValue(
          "すべて"
        )
        expect(runtimeBridge.completionEventCount()).toBe(1)
      })
    } finally {
      runtimeBridge.restore()
    }
  })

  test("page.selectedId がない import payload は summary.lastEntryId で同一ページ再選択する", async () => {
    const gateway = createTestGateway()
    const baseImport = gateway.importMasterDictionaryXml.bind(gateway)
    gateway.importMasterDictionaryXml = vi.fn(
      async (request: ImportMasterDictionaryXmlRequest) => {
        const payload: ImportMasterDictionaryXmlResponse =
          await baseImport(request)
        return {
          ...payload,
          page: payload.page
            ? {
                ...payload.page,
                selectedId: undefined
              }
            : undefined
        }
      }
    )

    const user = userEvent.setup()
    renderAppWithGateway(gateway)

    const globalNavigation = screen.getByRole("navigation", {
      name: "グローバルナビゲーション"
    })
    await user.click(
      within(globalNavigation).getByRole("link", { name: "マスター辞書" })
    )

    const xmlInput = document.querySelector("#xmlFileInput")
    if (!(xmlInput instanceof HTMLInputElement)) {
      throw new Error("xmlFileInput が見つかりません")
    }

    const xmlFile = createXmlFile(
      IMPORT_XML_WITH_MIXED_REC,
      "Dawnguard_english_japanese.xml"
    )
    await user.upload(xmlInput, xmlFile)
    await user.click(
      screen.getByRole("button", { name: "この XML を取り込む" })
    )

    await waitFor(() => {
      expect(screen.getByText("完了")).toBeInTheDocument()
      expect(document.querySelector("#detailTitle")).toHaveTextContent(
        "Allowed Book Source"
      )
      expect(
        document.querySelector("#importResultSelection")
      ).toHaveTextContent("Allowed Book Source")
    })
  })

  test("backend failure 時は成功表示を出さずにエラーを表示する", async () => {
    const gateway = createTestGateway()
    gateway.createMasterDictionaryEntry = vi.fn(() => {
      throw new Error("create failed")
    })
    const user = userEvent.setup()
    renderAppWithGateway(gateway)

    const globalNavigation = screen.getByRole("navigation", {
      name: "グローバルナビゲーション"
    })
    await user.click(
      within(globalNavigation).getByRole("link", { name: "マスター辞書" })
    )
    await user.click(screen.getByRole("button", { name: "新規登録" }))
    await user.type(screen.getByLabelText("原文"), "Failure Source")
    await user.type(screen.getByLabelText("訳語"), "失敗訳語")
    await user.click(screen.getByRole("button", { name: "保存する" }))

    await waitFor(() => {
      expect(screen.getByText("create failed")).toBeInTheDocument()
    })
    expect(document.querySelector("#detailTitle")).not.toHaveTextContent(
      "Failure Source"
    )
    expect(document.querySelector("#importResult")).toHaveAttribute("hidden")
  })

  test("SCN-DAS-004: 非ダッシュボード遷移先はプレースホルダーで導線切れを起こさない", async () => {
    const user = userEvent.setup()
    renderAppWithGateway()

    const globalNavigation = screen.getByRole("navigation", {
      name: "グローバルナビゲーション"
    })
    await user.click(
      within(globalNavigation).getByRole("link", { name: "翻訳管理" })
    )
    expect(
      screen.getByRole("heading", { level: 1, name: "翻訳管理" })
    ).toBeInTheDocument()
    expect(screen.getByText(PLACEHOLDER_LEAD)).toBeInTheDocument()
    expect(
      screen.queryByRole("heading", { name: "作業を選ぶ" })
    ).not.toBeInTheDocument()

    const placeholderLead = screen.getByText(PLACEHOLDER_LEAD)
    const placeholderCard = placeholderLead.closest("section")

    if (!placeholderCard) {
      throw new Error("プレースホルダー導線のセクションが見つかりません")
    }

    await user.click(
      within(placeholderCard).getByRole("link", { name: "出力管理" })
    )
    expect(
      screen.getByRole("heading", { level: 1, name: "出力管理" })
    ).toBeInTheDocument()
    expect(screen.getByText(PLACEHOLDER_LEAD)).toBeInTheDocument()
  })

  test("SCN-DAS-005: プレースホルダー表示中も共通シェルを保持して再移動できる", async () => {
    const user = userEvent.setup()
    renderAppWithGateway()

    const globalNavigation = screen.getByRole("navigation", {
      name: "グローバルナビゲーション"
    })

    await user.click(
      within(globalNavigation).getByRole("link", { name: "翻訳管理" })
    )
    expect(
      screen.getByRole("heading", { level: 1, name: "翻訳管理" })
    ).toBeInTheDocument()

    await user.click(
      within(globalNavigation).getByRole("link", { name: "ダッシュボード" })
    )
    expect(
      screen.getByRole("heading", { level: 1, name: "ダッシュボード" })
    ).toBeInTheDocument()

    await user.click(
      within(globalNavigation).getByRole("link", { name: "出力管理" })
    )
    expect(
      screen.getByRole("heading", { level: 1, name: "出力管理" })
    ).toBeInTheDocument()
  })
})
