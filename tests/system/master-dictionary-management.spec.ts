import path from "node:path"

import { expect, test, type Page } from "@playwright/test"

const dawnguardXmlPath = path.resolve(
  process.cwd(),
  "dictionaries/Dawnguard_english_japanese.xml"
)

async function openMasterDictionary(page: Page): Promise<void> {
  await page.goto("/")
  await page.getByRole("link", { name: "マスター辞書" }).first().click()
  await expect(page.getByRole("heading", { level: 1, name: "マスター辞書" })).toBeVisible()
}

async function importDawnguardXml(page: Page): Promise<void> {
  const xmlFileInput = page.locator("#xmlFileInput")
  const importStatusValue = page.locator("#importStatusValue")

  await xmlFileInput.setInputFiles(dawnguardXmlPath)
  await expect(page.locator("#selectedFileName")).toHaveText("Dawnguard_english_japanese.xml")
  await expect(page.locator("#importBar")).toBeVisible()
  await expect(importStatusValue).toHaveText("取込待ち")

  await page.getByRole("button", { name: "この XML を取り込む" }).click()
  await expect(importStatusValue).toHaveText("取込中")
  await expect(importStatusValue).toHaveText("完了", { timeout: 30000 })
}

test("SCN-MDM-001/002 一覧と検索を同一ページで操作できる", async ({ page }) => {
  await openMasterDictionary(page)

  await expect(page.getByRole("heading", { level: 2, name: "辞書一覧" })).toBeVisible()

  const rows = page.locator("#listStack .list-row")
  await expect(rows).toHaveCount(30)

  const secondRowSource = rows.nth(1).locator(".row-cell .row-value").nth(1)
  const secondRowSourceText = await secondRowSource.innerText()
  await rows.nth(1).click()
  await expect(page.locator("#detailTitle")).toHaveText(secondRowSourceText)

  const searchInput = page.getByLabel("検索")
  await searchInput.fill("__no_such_term__")
  await expect(page.locator("#listStack .empty-state")).toContainText(
    "一致するエントリがありません"
  )

  await searchInput.fill("")
  await expect(rows).toHaveCount(30)
})

test("SCN-MDM-003/004/005 新規登録・更新・削除モーダルを完了できる", async ({ page }) => {
  await openMasterDictionary(page)

  const sourceText = "Phase5 Source Entry"
  const createdTranslation = "フェーズ5 作成訳語"
  const updatedTranslation = "フェーズ5 更新訳語"

  await page.getByRole("button", { name: "新規登録" }).click()
  const createDialog = page.getByRole("dialog", { name: "新規登録" })
  await expect(createDialog).toBeVisible()
  await createDialog.getByLabel("原文").fill(sourceText)
  await createDialog.getByLabel("訳語").fill(createdTranslation)
  await createDialog.getByLabel("カテゴリ").selectOption("固有名詞")
  await createDialog.getByLabel("由来").selectOption("手動登録")
  await createDialog.getByRole("button", { name: "保存する" }).click()
  await expect(createDialog).toBeHidden()

  const searchInput = page.getByLabel("検索")
  await searchInput.fill(sourceText)
  const rows = page.locator("#listStack .list-row")
  await expect(rows).toHaveCount(1)
  await expect(page.locator("#detailTitle")).toHaveText(sourceText)
  await expect(page.locator("#detailTranslation")).toHaveText(createdTranslation)

  await page.getByRole("button", { name: "更新" }).click()
  const editDialog = page.getByRole("dialog", { name: "更新" })
  await expect(editDialog).toBeVisible()
  await editDialog.getByLabel("訳語").fill(updatedTranslation)
  await editDialog.getByRole("button", { name: "保存する" }).click()
  await expect(editDialog).toBeHidden()
  await expect(page.locator("#detailTranslation")).toHaveText(updatedTranslation)

  await page.getByRole("button", { name: "削除" }).click()
  const deleteDialog = page.getByRole("dialog", { name: "削除の確認" })
  await expect(deleteDialog).toBeVisible()
  await deleteDialog.getByRole("button", { name: "削除する" }).click()
  await expect(deleteDialog).toBeHidden()
  await expect(page.locator("#listStack .empty-state")).toContainText(
    "一致するエントリがありません"
  )
})

test("SCN-MDM-008/009 XML未選択ゲートと取込バー状態遷移を確認できる", async ({ page }) => {
  await openMasterDictionary(page)

  const importBar = page.locator("#importBar")
  const importStartButton = page.getByRole("button", { name: "この XML を取り込む" })

  await expect(importBar).toBeHidden()
  await expect(importStartButton).toBeHidden()

  await importDawnguardXml(page)

  await expect(page.getByLabel("検索")).toHaveValue("")
  await expect(page.getByLabel("カテゴリ")).toHaveValue("すべて")
  await expect(page.locator("#importResult")).toBeVisible()
})

test("SCN-MDM-006 XML取込は許可RECのみを抽出する", async ({ page }) => {
  await openMasterDictionary(page)
  await importDawnguardXml(page)

  const searchInput = page.getByLabel("検索")

  await searchInput.fill("Auriel's Bow")
  await expect(page.locator("#listStack .list-row")).not.toHaveCount(0)

  await searchInput.fill("Crossbow Mount")
  await expect(page.locator("#listStack .empty-state")).toContainText(
    "一致するエントリがありません"
  )

  await searchInput.fill("Transform into the vampire lord.")
  await expect(page.locator("#listStack .empty-state")).toContainText(
    "一致するエントリがありません"
  )
})
