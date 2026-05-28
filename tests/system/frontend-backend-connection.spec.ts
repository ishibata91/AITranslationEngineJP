import { expect, test, type Page } from "@playwright/test"

import { ProviderSettingsPage } from "./support/system-test-pages"
import { holdWailsBindingPending } from "./support/runtime-boundary-helpers"

const TEST_SECRET = "system-test-secret-value"

test.describe.configure({ mode: "serial" })

async function openProviderSettingsThroughProductionApp(page: Page) {
  await new ProviderSettingsPage(page).open()
  await expect(page).toHaveURL(/#provider-settings$/)
}

test("FBC-SC-001 provider settings production factory reaches AppController binding", async ({
  page
}) => {
  // production factory が injected gateway から generated binding を通り、backend controller response を画面初期読込へ反映する。
  await openProviderSettingsThroughProductionApp(page)

  const currentUrl = new URL(page.url())
  expect(currentUrl.searchParams.has("fakeApi")).toBe(false)
  const providerSettings = new ProviderSettingsPage(page)

  await expect(providerSettings.summaryRegion).toContainText(
    "Gateway: 接続準備済み"
  )
  await expect
    .poll(async () => providerSettings.serviceRows.count())
    .toBeGreaterThan(0)
  await expect(providerSettings.summaryRegion).toContainText(
    /[1-9]\d* 件の AIサービスを管理します。/
  )

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

test("E2E-UC-003 provider settings saves API key without echoing secret", async ({
  page
}) => {
  // APIキー保存の利用者操作で保存済み状態になり、入力した秘密値が画面へ残らないことを証明する。
  await openProviderSettingsThroughProductionApp(page)
  const providerSettings = new ProviderSettingsPage(page)

  await providerSettings.selectService("Gemini")
  await providerSettings.saveApiKey(TEST_SECRET)

  await expect(providerSettings.apiKeyStatusRegion).toContainText("保存済み")
  await expect(providerSettings.apiKeyInputRegion).toHaveCount(0)
  await expect(page.locator("body")).not.toContainText(TEST_SECRET)
})

test("E2E-UC-004 provider settings validates LM Studio connection through fake boundary", async ({
  page
}) => {
  // 接続確認の利用者操作で fake 境界の結果が表示され、設定保存状態は自動更新されないことを証明する。
  await openProviderSettingsThroughProductionApp(page)
  const providerSettings = new ProviderSettingsPage(page)

  await providerSettings.selectService("LM Studio")
  await providerSettings.fillEndpoint("http://127.0.0.1:1234/v1")
  await providerSettings.saveSettings()
  await providerSettings.checkConnection()

  await expect(providerSettings.connectionCheckRegion).toContainText(
    /接続確認成功|接続可能|確認済み/
  )
  await expect(providerSettings.settingsDetailRegion).not.toContainText(
    "保存通知"
  )
})

test("E2E-UC-005 provider settings saves endpoint from UI operation", async ({
  page
}) => {
  // 設定保存の利用者操作で、入力した endpoint が保存済み設定として表示されることを証明する。
  await openProviderSettingsThroughProductionApp(page)
  const providerSettings = new ProviderSettingsPage(page)

  await providerSettings.selectService("LM Studio")
  await providerSettings.fillEndpoint("http://127.0.0.1:1234/v1")
  await providerSettings.saveSettings()

  await expect(providerSettings.summaryRegion).toContainText(/保存|更新/)
  await expect(providerSettings.endpointInput).toHaveValue(
    "http://127.0.0.1:1234/v1"
  )
})

test("E2E-UC-006 provider settings reset clears saved endpoint state", async ({
  page
}) => {
  // リセットの利用者操作で、保存済み endpoint が初期状態へ戻ることを証明する。
  await openProviderSettingsThroughProductionApp(page)
  const providerSettings = new ProviderSettingsPage(page)

  await providerSettings.selectService("LM Studio")
  await providerSettings.fillEndpoint("http://127.0.0.1:1234/v1")
  await providerSettings.saveSettings()
  await expect(providerSettings.resetButton).toBeEnabled()

  await providerSettings.resetSettings()

  await expect(providerSettings.settingsDetailRegion).not.toContainText(
    "http://127.0.0.1:1234/v1"
  )
})

test("E2E-UC-027 provider settings suppresses duplicate connection checks while validating", async ({
  page
}) => {
  // 接続確認中の利用者操作で、接続確認ボタンが無効化され重複実行へ進まないことを証明する。
  await openProviderSettingsThroughProductionApp(page)
  const providerSettings = new ProviderSettingsPage(page)

  await providerSettings.selectService("LM Studio")
  await providerSettings.fillEndpoint("http://127.0.0.1:1234/v1")
  await holdWailsBindingPending(
    page,
    "AppController",
    "ValidateProviderSettings"
  )
  await providerSettings.checkConnection()

  await expect(providerSettings.connectionCheckRegion).toContainText(
    /接続確認中|確認中/
  )
  await expect(providerSettings.connectionCheckButton).toBeDisabled()
})

test("E2E-UC-028 provider settings rejects invalid endpoint input", async ({
  page
}) => {
  // 不正 endpoint の保存操作で、保存済み設定へ反映されないことを証明する。
  await openProviderSettingsThroughProductionApp(page)
  const providerSettings = new ProviderSettingsPage(page)
  const savedEndpoint = "http://127.0.0.1:1234/v1"

  await providerSettings.selectService("LM Studio")
  await providerSettings.fillEndpoint(savedEndpoint)
  await providerSettings.saveSettings()
  await expect(providerSettings.summaryRegion).toContainText(/保存|更新/)
  await page.reload()
  await expect(providerSettings.screen).toBeVisible()
  await providerSettings.selectService("LM Studio")
  await expect(providerSettings.endpointInput).toHaveValue(savedEndpoint)

  await providerSettings.fillEndpoint("invalid-endpoint")
  await expect(providerSettings.endpointInput).toHaveValue("invalid-endpoint")
  await providerSettings.saveSettings()

  await expect(providerSettings.summaryRegion).toContainText(/不正|エラー|入力/)
  await page.reload()
  await expect(providerSettings.screen).toBeVisible()
  await providerSettings.selectService("LM Studio")
  await expect(providerSettings.endpointInput).toHaveValue(savedEndpoint)
  await expect(providerSettings.endpointInput).not.toHaveValue(
    "invalid-endpoint"
  )
})
