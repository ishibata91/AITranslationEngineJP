import path from "node:path"

import { expect, test, type Page } from "@playwright/test"

test.describe.configure({ mode: "serial" })

const dawnguardXmlPath = path.resolve(
  process.cwd(),
  "tests/fixtures/master-dictionary/Dawnguard_english_japanese.xml"
)

async function openMasterDictionary(page: Page): Promise<void> {
  await page.goto("/")
  await page.getByRole("link", { name: "マスター辞書" }).first().click()
  await expect(
    page.getByRole("heading", { level: 1, name: "マスター辞書" })
  ).toBeVisible()
}

async function clickEditModalSave(page: Page): Promise<void> {
  const saveButton = page
    .locator("#editModal")
    .getByRole("button", { name: "保存する" })
  await expect(saveButton).toBeVisible()
  await expect(saveButton).toBeEnabled()
  await expect(async () => {
    await saveButton.click()
  }).toPass({ timeout: 10000 })
}

async function importDawnguardXml(page: Page): Promise<void> {
  const importStatusValue = page.locator("#importStatusValue")
  const importProgressFill = page.locator("#importProgressFill")
  const startImportButton = page.locator("#startImportButton")
  const importPanel = page.getByTestId("master-dictionary-xml-import-region")
  const progressPanel = page.getByTestId(
    "master-dictionary-import-progress-panel"
  )

  const stageXmlWithResolvedReference = async (): Promise<void> => {
    await page.evaluate((absolutePath) => {
      const input = document.getElementById("xmlFileInput")
      if (!(input instanceof HTMLInputElement)) {
        return
      }

      const file = new File([""], "Dawnguard_english_japanese.xml", {
        type: "text/xml"
      })
      Object.defineProperty(file, "path", {
        value: absolutePath,
        configurable: true
      })

      const transfer = new DataTransfer()
      transfer.items.add(file)
      Object.defineProperty(input, "files", {
        value: transfer.files,
        configurable: true
      })
      input.dispatchEvent(new Event("change", { bubbles: true }))
    }, dawnguardXmlPath)

    await expect(
      importPanel.getByText("Dawnguard_english_japanese.xml")
    ).toBeVisible()
    await expect(
      progressPanel.getByText("Dawnguard_english_japanese.xml")
    ).toBeVisible()
    await expect(progressPanel.locator(".current-file")).toContainText(
      "Dawnguard_english_japanese.xml"
    )
    await expect(
      progressPanel.locator('[aria-label="XML 取り込みの進行率"]')
    ).toBeVisible()
    await expect(importStatusValue).toHaveText("取込待ち")
    await expect(importProgressFill).toHaveAttribute("style", /width:\s*0%/)
  }

  await stageXmlWithResolvedReference()
  await startImportButton.click()
  await expect(importStatusValue).toHaveText("完了", { timeout: 30000 })
  await expect(page.locator("#importResultHeadline")).toBeVisible()
}

function dictionaryRows(page: Page) {
  return page.locator(
    ".target-table tbody tr:not(.target-detail-row):not(.target-empty-row)"
  )
}

function dictionaryEmptyRow(page: Page) {
  return page.locator(".target-empty-row")
}

test("SCN-MDM-001/002 一覧と検索を同一ページで操作できる", async ({ page }) => {
  await openMasterDictionary(page)

  await expect(page.locator("#dictionaryTargetHeading")).toBeVisible()

  const rows = dictionaryRows(page)
  await expect(rows).toHaveCount(30)

  const secondRowSource = rows.nth(1).locator("td").first()
  const secondRowSourceText = await secondRowSource.innerText()
  await rows.nth(1).click()
  await expect(page.locator("#detailTitle")).toHaveText(secondRowSourceText)

  const searchInput = page.getByLabel("検索")
  await searchInput.fill("__no_such_term__")
  await expect(dictionaryEmptyRow(page)).toContainText("処理対象がありません")

  await searchInput.fill("")
  await expect(rows).toHaveCount(30)
})

test("SCN-MDM-003/004/005 新規登録・更新・削除モーダルを完了できる", async ({
  page
}) => {
  await openMasterDictionary(page)

  const sourceText = `Phase5 Source Entry ${Date.now()}`
  const createdTranslation = "フェーズ5 作成訳語"
  const updatedTranslation = "フェーズ5 更新訳語"

  await page.getByRole("button", { name: "新規登録" }).click()
  const createDialog = page.locator("#editModal")
  await expect(createDialog).toBeVisible()
  await expect(page.locator("#editModalTitle")).toHaveText("新規登録")
  await createDialog.getByLabel("原文").fill(sourceText)
  await createDialog.getByLabel("訳語").fill(createdTranslation)
  await createDialog.getByLabel("由来").selectOption("手動登録")
  await clickEditModalSave(page)
  await expect(createDialog).toBeHidden()

  const searchInput = page.getByLabel("検索")
  await searchInput.fill(sourceText)
  const rows = dictionaryRows(page)
  await expect.poll(async () => rows.count()).toBeGreaterThan(0)
  await expect(page.locator("#detailTitle")).toHaveText(sourceText)
  await expect(rows.first().locator("td").nth(1)).toContainText(
    createdTranslation
  )

  await page.getByRole("button", { name: "更新" }).click()
  const editDialog = page.locator("#editModal")
  await expect(editDialog).toBeVisible()
  await expect(page.locator("#editModalTitle")).toHaveText("更新")
  await editDialog.getByLabel("訳語").fill(updatedTranslation)
  await clickEditModalSave(page)
  await expect(editDialog).toBeHidden()
  await expect(rows.first().locator("td").nth(1)).toContainText(
    updatedTranslation
  )

  await page.getByRole("button", { name: "削除" }).click()
  const deleteDialog = page.locator("#deleteModal")
  await expect(deleteDialog).toBeVisible()
  await deleteDialog.getByRole("button", { name: "削除する" }).click()
  await expect(deleteDialog).toBeHidden()

  await expect.poll(async () => rows.count()).toBe(0)
  await expect
    .poll(async () => await page.locator(".target-list-wrapper").innerText())
    .not.toContain(sourceText)
  await expect
    .poll(async () => await page.locator(".target-list-wrapper").innerText())
    .not.toContain(updatedTranslation)
})

test("SCN-MDM-008/009 XML未選択ゲートと取込バー状態遷移を確認できる", async ({
  page
}) => {
  await openMasterDictionary(page)

  const importProgressRegion = page.locator(
    '[aria-label="XML 取り込みの進行率"]'
  )
  const importStartButton = page.getByRole("button", {
    name: "この XML を取り込む"
  })

  await expect(importProgressRegion).toBeVisible()
  await expect(importStartButton).toBeHidden()

  await importDawnguardXml(page)

  await expect(page.locator("#searchInput")).toHaveValue("")
  await expect(page.locator("#categorySelect")).toHaveValue("すべて")
  await expect(page.locator("#importResultHeadline")).toBeVisible()
})

test("SCN-MDM-006 XML取込は許可RECのみを抽出する", async ({ page }) => {
  await openMasterDictionary(page)
  await importDawnguardXml(page)

  const searchInput = page.getByLabel("検索")

  await searchInput.fill("Auriel's Bow")
  await expect(dictionaryRows(page)).not.toHaveCount(0)

  await searchInput.fill("Crossbow Mount")
  await expect(dictionaryEmptyRow(page)).toContainText("処理対象がありません")

  await searchInput.fill("Transform into the vampire lord.")
  await expect(dictionaryEmptyRow(page)).toContainText("処理対象がありません")
})
