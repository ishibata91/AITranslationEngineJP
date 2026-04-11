import { expect, test } from "@playwright/test"

test("renders the app shell in browser mode", async ({ page }) => {
  await page.goto("/")

  await expect(
    page.getByRole("heading", { name: "Architecture Skeleton" })
  ).toBeVisible()
  await expect(page.getByText("AITranslationEngineJp")).toBeVisible()
})
