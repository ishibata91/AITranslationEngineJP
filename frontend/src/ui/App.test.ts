import { render, screen, within } from "@testing-library/svelte"
import userEvent from "@testing-library/user-event"

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
    render(App)

    expect(screen.getByRole("heading", { name: "ダッシュボード" })).toBeInTheDocument()
    expect(window.location.hash).toBe("#dashboard")
    for (const excludedText of DASHBOARD_EXCLUDED_TEXTS) {
      expect(screen.queryByText(excludedText)).not.toBeInTheDocument()
    }
  })

  test("invalid hash は dashboard に正規化される", () => {
    window.history.replaceState(null, "", "#not-approved-route")

    render(App)

    expect(screen.getByRole("heading", { name: "ダッシュボード" })).toBeInTheDocument()
    expect(window.location.hash).toBe("#dashboard")
  })

  test("SCN-DAS-002: グローバルナビゲーションから主要 5 ページへ遷移できる", async () => {
    const user = userEvent.setup()
    render(App)

    const globalNavigation = screen.getByRole("navigation", {
      name: "グローバルナビゲーション"
    })

    for (const route of DASHBOARD_SHELL_PRIMARY_ROUTES) {
      await user.click(within(globalNavigation).getByRole("link", { name: route.label }))
      expect(
        screen.getByRole("heading", { level: 1, name: route.label })
      ).toBeInTheDocument()
      expect(window.location.hash).toBe(`#${route.id}`)
    }
  })

  test("SCN-DAS-003: ダッシュボード入口カードから主要 5 ページへ遷移できる", async () => {
    const user = userEvent.setup()
    render(App)

    const globalNavigation = screen.getByRole("navigation", {
      name: "グローバルナビゲーション"
    })

    for (const route of DASHBOARD_ENTRY_ROUTES) {
      await user.click(within(globalNavigation).getByRole("link", { name: "ダッシュボード" }))

      const dashboardCardSectionHeading = screen.getByRole("heading", { name: "作業を選ぶ" })
      const dashboardCardSection = dashboardCardSectionHeading.closest("section")

      if (!dashboardCardSection) {
        throw new Error("ダッシュボード入口カードのセクションが見つかりません")
      }

      const cardHeading = within(dashboardCardSection).getByRole("heading", {
        level: 3,
        name: route.label
      })
      const cardLink = cardHeading.closest("a")

      if (!cardLink) {
        throw new Error(`ダッシュボード入口カードのリンクが見つかりません: ${route.label}`)
      }

      await user.click(
        cardLink
      )
      expect(
        screen.getByRole("heading", { level: 1, name: route.label })
      ).toBeInTheDocument()
      expect(window.location.hash).toBe(`#${route.id}`)
    }
  })

  test("主要導線は承認済み 5 ルートに固定される", () => {
    render(App)

    const globalNavigation = screen.getByRole("navigation", {
      name: "グローバルナビゲーション"
    })
    const globalLinks = within(globalNavigation).getAllByRole("link")
    expectRouteLinksMatch(globalLinks)

    const dashboardCardSectionHeading = screen.getByRole("heading", { name: "作業を選ぶ" })
    const dashboardCardSection = dashboardCardSectionHeading.closest("section")

    if (!dashboardCardSection) {
      throw new Error("ダッシュボード入口カードのセクションが見つかりません")
    }

    const dashboardLinks = within(dashboardCardSection).getAllByRole("link")
    expectRouteLinksMatch(dashboardLinks, DASHBOARD_ENTRY_ROUTES)

    expect(within(dashboardCardSection).queryByRole("link", { name: "ダッシュボード" })).not.toBeInTheDocument()

    for (const route of DASHBOARD_ENTRY_ROUTES) {
      expect(within(dashboardCardSection).getAllByText(route.state).length).toBeGreaterThan(0)
      expect(within(dashboardCardSection).getByText(route.description)).toBeInTheDocument()
    }
  })

  test("mobile nav トグルの表示契約を持つ", () => {
    render(App)

    const mobileToggle = screen.getByRole("button", { name: "主要ページ" })
    expect(mobileToggle).toHaveAttribute("aria-controls", "globalNav")
    expect(mobileToggle).toHaveAttribute("aria-expanded", "false")
  })

  test("非ダッシュボード遷移先ではダッシュボード専用領域を表示しない", async () => {
    const user = userEvent.setup()
    render(App)

    const globalNavigation = screen.getByRole("navigation", {
      name: "グローバルナビゲーション"
    })

    for (const route of DASHBOARD_SHELL_PRIMARY_ROUTES.filter(
      ({ id }) => id !== "dashboard"
    )) {
      await user.click(within(globalNavigation).getByRole("link", { name: route.label }))
      expect(
        screen.getByRole("heading", { level: 1, name: route.label })
      ).toBeInTheDocument()
      expect(screen.queryByRole("heading", { name: "作業を選ぶ" })).not.toBeInTheDocument()
      expect(screen.getByText(PLACEHOLDER_LEAD)).toBeInTheDocument()
    }
  })

  test("SCN-DAS-004: 非ダッシュボード遷移先はプレースホルダーで導線切れを起こさない", async () => {
    const user = userEvent.setup()
    render(App)

    const globalNavigation = screen.getByRole("navigation", {
      name: "グローバルナビゲーション"
    })
    await user.click(within(globalNavigation).getByRole("link", { name: "マスター辞書" }))
    expect(
      screen.getByRole("heading", { level: 1, name: "マスター辞書" })
    ).toBeInTheDocument()
    expect(screen.getByText(PLACEHOLDER_LEAD)).toBeInTheDocument()
    expect(screen.queryByRole("heading", { name: "作業を選ぶ" })).not.toBeInTheDocument()

    const placeholderLead = screen.getByText(PLACEHOLDER_LEAD)
    const placeholderCard = placeholderLead.closest("section")

    if (!placeholderCard) {
      throw new Error("プレースホルダー導線のセクションが見つかりません")
    }

    await user.click(within(placeholderCard).getByRole("link", { name: "出力管理" }))
    expect(
      screen.getByRole("heading", { level: 1, name: "出力管理" })
    ).toBeInTheDocument()
    expect(screen.getByText(PLACEHOLDER_LEAD)).toBeInTheDocument()
  })

  test("SCN-DAS-005: プレースホルダー表示中も共通シェルを保持して再移動できる", async () => {
    const user = userEvent.setup()
    render(App)

    const globalNavigation = screen.getByRole("navigation", {
      name: "グローバルナビゲーション"
    })

    await user.click(within(globalNavigation).getByRole("link", { name: "翻訳管理" }))
    expect(
      screen.getByRole("heading", { level: 1, name: "翻訳管理" })
    ).toBeInTheDocument()

    await user.click(within(globalNavigation).getByRole("link", { name: "ダッシュボード" }))
    expect(
      screen.getByRole("heading", { level: 1, name: "ダッシュボード" })
    ).toBeInTheDocument()

    await user.click(within(globalNavigation).getByRole("link", { name: "出力管理" }))
    expect(
      screen.getByRole("heading", { level: 1, name: "出力管理" })
    ).toBeInTheDocument()
  })
})
