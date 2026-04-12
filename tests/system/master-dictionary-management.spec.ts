import path from "node:path"

import { expect, test, type Page } from "@playwright/test"

test.describe.configure({ mode: "serial" })

const dawnguardXmlPath = path.resolve(
  process.cwd(),
  "dictionaries/Dawnguard_english_japanese.xml"
)

async function openMasterDictionary(page: Page): Promise<void> {
  await page.goto("/")
  await page.getByRole("link", { name: "マスター辞書" }).first().click()
  await expect(page.getByRole("heading", { level: 1, name: "マスター辞書" })).toBeVisible()
}

async function clickEditModalSave(page: Page): Promise<void> {
  const saveButton = page.locator("#editModal").getByRole("button", { name: "保存する" })
  await expect(saveButton).toBeVisible()
  await expect(saveButton).toBeEnabled()
  await expect(async () => {
    await saveButton.click()
  }).toPass({ timeout: 10000 })
}

async function importDawnguardXml(page: Page): Promise<void> {
  const xmlFileInput = page.locator("#xmlFileInput")
  const importStatusValue = page.locator("#importStatusValue")
  const importProgressFill = page.locator("#importProgressFill")
  const startImportButton = page.locator("#startImportButton")

  const stageXmlWithResolvedReference = async (usePathInjection: boolean): Promise<void> => {
    if (usePathInjection) {
      await page.evaluate((absolutePath) => {
        const input = document.getElementById("xmlFileInput")
        if (!(input instanceof HTMLInputElement)) {
          return
        }

        const file = new File([""], "Dawnguard_english_japanese.xml", { type: "text/xml" })
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
    } else {
      await xmlFileInput.setInputFiles(dawnguardXmlPath)
    }

    await expect(page.locator("#selectedFileName")).toHaveText("Dawnguard_english_japanese.xml")
    await expect(page.locator("#importBar")).toBeVisible()
    await expect(importStatusValue).toHaveText("取込待ち")
    await expect(importProgressFill).toHaveAttribute("style", /width:\s*0%/)
  }

  for (let attempt = 0; attempt < 3; attempt += 1) {
    await stageXmlWithResolvedReference(attempt > 0)
    await startImportButton.click()

    let observedRunning = false
    let observedStatus = "取込待ち"
    for (let poll = 0; poll < 10; poll += 1) {
      observedStatus = (await importStatusValue.innerText()).trim()
      if (observedStatus === "取込中") {
        observedRunning = true
      }
      if (observedStatus !== "取込待ち") {
        break
      }
      await page.waitForTimeout(250)
    }

    if (!observedRunning) {
      continue
    }

    await expect(importStatusValue).toHaveText("完了", { timeout: 30000 })
    await expect(page.locator("#importResult")).toBeVisible()
    return
  }

  throw new Error("XML 取り込みが完了状態へ到達しませんでした。")
}

test("SCN-MDM-001/002 一覧と検索を同一ページで操作できる", async ({ page }) => {
  await openMasterDictionary(page)

  await expect(page.locator("#listHeading")).toBeVisible()

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
  const rows = page.locator("#listStack .list-row")
  await expect.poll(async () => rows.count()).toBeGreaterThan(0)
  await expect(page.locator("#detailTitle")).toHaveText(sourceText)
  await expect(page.locator("#detailTranslation")).toHaveText(createdTranslation)

  await page.getByRole("button", { name: "更新" }).click()
  const editDialog = page.locator("#editModal")
  await expect(editDialog).toBeVisible()
  await expect(page.locator("#editModalTitle")).toHaveText("更新")
  await editDialog.getByLabel("訳語").fill(updatedTranslation)
  await clickEditModalSave(page)
  await expect(editDialog).toBeHidden()
  await expect(page.locator("#detailTranslation")).toHaveText(updatedTranslation)

  await page.getByRole("button", { name: "削除" }).click()
  const deleteDialog = page.locator("#deleteModal")
  await expect(deleteDialog).toBeVisible()
  await deleteDialog.getByRole("button", { name: "削除する" }).click()
  await expect(deleteDialog).toBeHidden()

  const listStack = page.locator("#listStack")
  await expect.poll(async () => rows.count()).toBe(0)
  await expect.poll(async () => await listStack.innerText()).not.toContain(sourceText)
  await expect.poll(async () => await listStack.innerText()).not.toContain(updatedTranslation)
})

test("SCN-MDM-008/009 XML未選択ゲートと取込バー状態遷移を確認できる", async ({ page }) => {
  await openMasterDictionary(page)

  const importBar = page.locator("#importBar")
  const importStartButton = page.getByRole("button", { name: "この XML を取り込む" })

  await expect(importBar).toBeHidden()
  await expect(importStartButton).toBeHidden()

  await importDawnguardXml(page)

  await expect(page.locator("#searchInput")).toHaveValue("")
  await expect(page.locator("#categorySelect")).toHaveValue("すべて")
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
