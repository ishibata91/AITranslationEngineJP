import { expect, test } from "@playwright/test"

test("renders dashboard app shell in browser mode", async ({ page }) => {
  await page.goto("/")

  await expect(page.getByText("AITranslationEngineJp")).toBeVisible()
  await expect(
    page.getByRole("navigation", { name: "グローバルナビゲーション" })
  ).toBeVisible()
  await expect(page.getByRole("heading", { level: 1, name: "ダッシュボード" })).toBeVisible()
  await expect(page.getByRole("heading", { level: 2, name: "作業を選ぶ" })).toBeVisible()
  await expect(page.getByText("現在地", { exact: false })).toBeVisible()
  await expect(page.getByText("ジョブ一覧", { exact: true })).toHaveCount(0)
  await expect(page.getByText("進捗サマリ", { exact: true })).toHaveCount(0)
})
