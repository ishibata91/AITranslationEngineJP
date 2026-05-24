import { expect, test, type Page } from "@playwright/test"

const PROVIDER_SETTINGS_URL = "/#provider-settings"

async function openProviderSettingsThroughProductionApp(page: Page) {
  await page.goto(PROVIDER_SETTINGS_URL)
  await expect(
    page.getByRole("heading", { level: 1, name: "AIサービス設定" })
  ).toBeVisible()
  await expect(page).toHaveURL(/#provider-settings$/)
}

test("FBC-SC-001 provider settings production factory reaches AppController binding", async ({
  page
}) => {
  // production factory が injected gateway から generated binding を通り、backend controller response を画面初期読込へ反映する。
  await openProviderSettingsThroughProductionApp(page)

  const currentUrl = new URL(page.url())
  expect(currentUrl.searchParams.has("fakeApi")).toBe(false)

  await expect(
    page.getByTestId("provider-settings-screen-summary-region")
  ).toContainText("Gateway: 接続準備済み")
  await expect
    .poll(async () =>
      page.getByTestId("provider-settings-ai-service-row").count()
    )
    .toBeGreaterThan(0)
  await expect(
    page.getByTestId("provider-settings-screen-summary-region")
  ).toContainText(/[1-9]\d* 件の AIサービスを管理します。/)

  const health = await page.evaluate(async () => {
    const wailsWindow = window as typeof window & {
      go?: {
        wails?: {
          AppController?: {
            Health?: () => Promise<{ status?: string }>
          }
        }
      }
    }

    return wailsWindow.go?.wails?.AppController?.Health?.()
  })

  expect(health).toEqual({ status: "ok" })
})
