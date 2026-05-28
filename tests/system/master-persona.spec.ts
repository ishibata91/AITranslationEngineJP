import path from "node:path"

import { expect, test } from "@playwright/test"

import { MasterPersonaPage } from "./support/system-test-pages"
import { installScenarioWailsMocks } from "./support/scenario-wails-mocks"

test.describe.configure({ mode: "serial" })

const personaJsonPath = path.resolve(
  process.cwd(),
  "tests/fixtures/master-persona/system-test-persona.json"
)

test.beforeEach(async ({ page }) => {
  await installScenarioWailsMocks(page)
})

test("E2E-UC-013 master persona generates personas from selected JSON", async ({
  page
}) => {
  // JSON 選択と AI 設定の利用者操作で、fake 境界の生成結果が一覧と詳細へ表示されることを証明する。
  const masterPersona = new MasterPersonaPage(page)
  await masterPersona.open()

  await masterPersona.setJsonFile(personaJsonPath)
  await masterPersona.selectAISettings({
    provider: "Gemini",
    model: "gemini-test",
    executionMode: "通常"
  })
  await masterPersona.waitForSelectableModel("gemini-test")
  await expect(masterPersona.modelSelect).toBeEnabled()
  await expect(masterPersona.generateButton).toBeEnabled()
  await masterPersona.generate()

  await expect(masterPersona.progressPanel).toContainText(/完了|生成完了/)
  await expect(masterPersona.resultListPanel).toContainText("Lydia")
  await expect(masterPersona.resultListPanel).toContainText("000A2C8E")
  await expect(masterPersona.resultListPanel).toContainText("丁寧な口調")
  await expect(masterPersona.detailPanel).toContainText(/従士|ホワイトラン/)
})

test("E2E-UC-014 master persona filters existing personas", async ({
  page
}) => {
  // 検索とプラグイン選択の利用者操作で、条件一致したペルソナ詳細を参照できることを証明する。
  const masterPersona = new MasterPersonaPage(page)
  await masterPersona.open()

  await masterPersona.search("Lydia")
  await masterPersona.selectPlugin("Skyrim.esm")
  await masterPersona.selectPersona("Lydia")

  await expect(masterPersona.resultListPanel).toContainText("Lydia")
  await expect(masterPersona.resultListPanel).not.toContainText("Serana")
  await expect(masterPersona.detailPanel).toContainText("000A2C8E")
  await expect(masterPersona.detailPanel).toContainText(/従士|ホワイトラン/)
})

test("E2E-UC-015 master persona updates selected persona", async ({
  page
}) => {
  // 編集 modal の利用者操作で、選択したペルソナ本文と metadata が更新されることを証明する。
  const masterPersona = new MasterPersonaPage(page)
  await masterPersona.open()

  await masterPersona.selectPersona("Lydia")
  await masterPersona.openEditModal()
  await masterPersona.fillEdit({
    summary: "忠実な従士",
    speechStyle: "落ち着いた口調",
    body: "新本文"
  })
  await masterPersona.saveEdit()

  await expect(masterPersona.editModal).toBeHidden()
  await expect(masterPersona.detailPanel).toContainText("新本文")
  await expect(masterPersona.resultListPanel).toContainText("落ち着いた口調")
})

test("E2E-UC-016 master persona deletes selected persona", async ({
  page
}) => {
  // 削除確定の利用者操作で、選択したペルソナが一覧から消えることを証明する。
  const masterPersona = new MasterPersonaPage(page)
  await masterPersona.open()

  await masterPersona.search("Lydia")
  await masterPersona.selectPersona("Lydia")
  await masterPersona.openDeleteModal()
  await masterPersona.confirmDelete()

  await expect(masterPersona.deleteModal).toBeHidden()
  await expect(masterPersona.resultListPanel).not.toContainText("Lydia")
})

test("E2E-UC-033 master persona keeps generation disabled when AI settings are incomplete", async ({
  page
}) => {
  // AI 設定不足の前提では、JSON 選択後も生成操作が無効で未実行状態を維持することを証明する。
  await installScenarioWailsMocks(page, { masterPersonaAISettings: "missing" })
  const masterPersona = new MasterPersonaPage(page)
  await masterPersona.open()

  await masterPersona.setJsonFile(personaJsonPath)

  await expect(masterPersona.aiSettingsCard).toContainText(/APIキーが未設定/)
  await expect(masterPersona.generateButton).toBeDisabled()
  await expect(masterPersona.progressPanel).toContainText(/入力待ち|未実行/)
})

test("E2E-UC-034 master persona shows empty state for unmatched search", async ({
  page
}) => {
  // ペルソナ検索が 0 件になる利用者操作で、一覧が該当なし状態になることを証明する。
  const masterPersona = new MasterPersonaPage(page)
  await masterPersona.open()

  await masterPersona.search("NoSuchPersona")

  await expect(masterPersona.resultListPanel).toContainText(/0 件|該当|ありません/)
})

test("E2E-UC-035 master persona cancels edit without changing detail", async ({
  page
}) => {
  // 編集取消の利用者操作で、入力した本文が詳細へ反映されないことを証明する。
  const masterPersona = new MasterPersonaPage(page)
  await masterPersona.open()

  await masterPersona.selectPersona("Lydia")
  await masterPersona.openEditModal()
  await masterPersona.fillEdit({ body: "キャンセルする本文" })
  await masterPersona.cancelEdit()

  await expect(masterPersona.editModal).toBeHidden()
  await expect(masterPersona.detailPanel).not.toContainText("キャンセルする本文")
  await expect(masterPersona.detailPanel).toContainText(/従士|ホワイトラン/)
})
