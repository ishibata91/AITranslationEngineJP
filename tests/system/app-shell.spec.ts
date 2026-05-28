import { expect, test } from "@playwright/test"

import {
  DashboardPage,
  MasterDictionaryPage,
  ProviderSettingsPage,
  SYSTEM_TEST_MOBILE_VIEWPORT
} from "./support/system-test-pages"

test("SCN-APP-001 dashboard renders the initial app shell", async ({
  page
}) => {
  // ダッシュボード初期表示が、主要機能選択と共通ナビゲーションの入口を提示することを証明する。
  const dashboard = new DashboardPage(page)
  await dashboard.open()

  await expect(page.getByText("AITranslationEngineJp")).toBeVisible()
  await expect(dashboard.currentPageDescription).toContainText("ダッシュボード")
  await expect(dashboard.globalNavigation).toBeVisible()
  await expect(
    page.getByRole("heading", { level: 1, name: "ダッシュボード" })
  ).toBeVisible()
  await expect(
    page.getByRole("heading", { level: 2, name: "作業を選ぶ" })
  ).toBeVisible()
  await expect(dashboard.navigationItem("ダッシュボード")).toBeVisible()
  await expect(
    page.getByRole("link", { name: "ダッシュボード" }).nth(1)
  ).toHaveCount(0)
  await expect(page.getByText("ジョブ一覧", { exact: true })).toHaveCount(0)
  await expect(page.getByText("進捗サマリ", { exact: true })).toHaveCount(0)
})

test("E2E-UC-001 dashboard card opens provider settings", async ({ page }) => {
  // 主要機能カードの利用者操作で、AIサービス設定画面へ移動できることを証明する。
  const dashboard = new DashboardPage(page)
  const providerSettings = new ProviderSettingsPage(page)
  await dashboard.open()

  await dashboard.openCard("AIサービス設定")

  await expect(providerSettings.screen).toBeVisible()
  await expect(dashboard.navigationItem("AIサービス設定")).toHaveAttribute(
    "aria-current",
    "page"
  )
})

test("E2E-UC-026 selecting current dashboard keeps the dashboard context", async ({
  page
}) => {
  // 現在表示中のダッシュボードを選び直しても、現在地表示が維持されることを証明する。
  const dashboard = new DashboardPage(page)
  await dashboard.open()

  await dashboard.selectNavigationItem("ダッシュボード")

  await expect(dashboard.content).toBeVisible()
  await expect(dashboard.currentPageDescription).toContainText("ダッシュボード")
  await expect(dashboard.navigationItem("ダッシュボード")).toHaveAttribute(
    "aria-current",
    "page"
  )
})

test("E2E-UC-002 mobile navigation opens master dictionary", async ({
  page
}) => {
  // モバイル viewport の利用者操作で、主要ページ menu からマスター辞書へ移動できることを証明する。
  const dashboard = new DashboardPage(page)
  const masterDictionary = new MasterDictionaryPage(page)
  await page.setViewportSize(SYSTEM_TEST_MOBILE_VIEWPORT)
  await dashboard.open()

  await dashboard.openMobileNavigation()
  await dashboard.selectNavigationItem("マスター辞書")

  await expect(masterDictionary.screen).toBeVisible()
})
