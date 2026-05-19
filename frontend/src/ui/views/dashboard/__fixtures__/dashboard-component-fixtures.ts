import type { ShellRouteContract, ShellRouteId } from "@ui/stores/shell-state"
import type {
  AppHeaderProps,
  CurrentPageHeroProps,
  DashboardEntryCardProps,
  DashboardEntryGridProps,
  GlobalNavigationProps
} from "../dashboard-component-props"

const ignoreRoute = (routeId: ShellRouteId): void => {
  void routeId
}

const ignoreAction = (): void => {}

export const dashboardRoutes = [
  {
    id: "dashboard",
    label: "ダッシュボード",
    state: "既定表示",
    lead: "最初に移動したい作業を選びます。",
    description: "主要ページへの入口をまとめて確認します。"
  },
  {
    id: "provider-settings",
    label: "AIサービス設定",
    state: "設定状態を確認",
    lead: "AIサービスごとの接続設定を確認します。",
    description: "エンドポイントと APIキー状態を AIサービスごとに確認します。"
  },
  {
    id: "master-dictionary",
    label: "マスター辞書",
    state: "確認可能",
    lead: "用語と訳語の基盤データを確認します。",
    description: "用語と訳語の基盤データを確認します。"
  },
  {
    id: "master-persona",
    label: "マスターペルソナ",
    state: "確認可能",
    lead: "翻訳前のペルソナを確認します。",
    description: "翻訳に使うペルソナ設定を確認します。"
  },
  {
    id: "translation-management",
    label: "翻訳管理",
    state: "未完了ジョブ入口",
    lead: "未完了ジョブ一覧から翻訳を進めます。",
    description: "未完了ジョブ一覧と翻訳再開を扱います。"
  },
  {
    id: "output-management",
    label: "出力管理",
    state: "確認可能",
    lead: "生成された成果物を確認します。",
    description: "生成物と書き出し結果を確認します。"
  }
] satisfies ShellRouteContract[]

const currentRoute = dashboardRoutes[0]
const entryRoutes = dashboardRoutes.filter((route) => route.id !== "dashboard")

export const appHeaderFixtures = {
  desktop: {
    currentRoute,
    isMobileNavOpen: false,
    routes: dashboardRoutes,
    selectRoute: ignoreRoute,
    toggleMobileNav: ignoreAction
  },
  mobileOpen: {
    currentRoute: dashboardRoutes[1],
    isMobileNavOpen: true,
    routes: dashboardRoutes,
    selectRoute: ignoreRoute,
    toggleMobileNav: ignoreAction
  }
} satisfies Record<string, AppHeaderProps>

export const globalNavigationFixtures = {
  dashboard: {
    currentRoute,
    routes: dashboardRoutes,
    selectRoute: ignoreRoute
  }
} satisfies Record<string, GlobalNavigationProps>

export const currentPageHeroFixtures = {
  dashboard: {
    currentRoute,
    dataTestId: "dashboard-current-page-description"
  },
  providerSettings: {
    currentRoute: dashboardRoutes[1],
    dataTestId: "dashboard-current-page-description"
  }
} satisfies Record<string, CurrentPageHeroProps>

export const dashboardEntryGridFixtures = {
  standard: {
    routes: entryRoutes,
    selectRoute: ignoreRoute
  }
} satisfies Record<string, DashboardEntryGridProps>

export const dashboardEntryCardFixtures = {
  providerSettings: {
    route: dashboardRoutes[1],
    selectRoute: ignoreRoute
  },
  longLabel: {
    route: {
      id: "translation-management",
      label: "翻訳管理",
      state: "未完了ジョブ入口",
      lead: "未完了ジョブ一覧から翻訳を進めます。",
      description:
        "未完了ジョブ一覧、新規翻訳の開始、対象ジョブの現在の翻訳段階への再開を扱います。"
    },
    selectRoute: ignoreRoute
  }
} satisfies Record<string, DashboardEntryCardProps>
