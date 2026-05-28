import path from "node:path"

import { expect, test, type Page } from "@playwright/test"

import { MasterDictionaryPage } from "./support/system-test-pages"

test.describe.configure({ mode: "serial" })

const dawnguardXmlPath = path.resolve(
  process.cwd(),
  "tests/fixtures/master-dictionary/Dawnguard_english_japanese.xml"
)
const invalidXmlPath = path.resolve(
  process.cwd(),
  "tests/fixtures/master-dictionary/invalid_master_dictionary.xml"
)

async function openMasterDictionary(page: Page): Promise<void> {
  await new MasterDictionaryPage(page).open()
}

async function clickEditModalSave(
  dictionary: MasterDictionaryPage
): Promise<void> {
  const saveButton = dictionary.entrySaveButton
  await expect(saveButton).toBeVisible()
  await expect(saveButton).toBeEnabled()
  await expect(async () => {
    await dictionary.saveEntry()
  }).toPass({ timeout: 10000 })
}

async function importDawnguardXml(page: Page): Promise<void> {
  const dictionary = new MasterDictionaryPage(page)

  const stageXmlWithResolvedReference = async (): Promise<void> => {
    await dictionary.stageXmlFileWithRuntimePath(dawnguardXmlPath)
    await expect(
      dictionary.xmlImportRegion.getByText("Dawnguard_english_japanese.xml")
    ).toBeVisible()
    await expect(
      dictionary.xmlImportProgressPanel.getByText(
        "Dawnguard_english_japanese.xml"
      )
    ).toBeVisible()
    await expect(
      dictionary.xmlImportProgressPanel.getByLabel("XML 取り込みの進行率")
    ).toBeVisible()
    await expect(dictionary.xmlImportProgressPanel).toContainText("取込待ち")
  }

  await stageXmlWithResolvedReference()
  await dictionary.startXmlImport()
  await expect(dictionary.xmlImportProgressPanel).toContainText("完了", {
    timeout: 30000
  })
}

test("E2E-UC-007 master dictionary filters entries by search text", async ({
  page
}) => {
  // 検索入力の利用者操作で、条件に一致する辞書エントリだけが表示されることを証明する。
  await openMasterDictionary(page)
  const dictionary = new MasterDictionaryPage(page)

  await expect(dictionary.operationRegion).toBeVisible()

  await expect(dictionary.entryRows).toHaveCount(30)

  await dictionary.search("Whiterun")

  await expect
    .poll(async () => await dictionary.operationRegion.innerText())
    .toContain("Whiterun")
  await expect(dictionary.operationRegion).not.toContainText("Iron Sword")
})

test("E2E-UC-029 master dictionary shows empty state for no search result", async ({
  page
}) => {
  // 検索結果 0 件の利用者操作で、一覧が該当なし状態になることを証明する。
  await openMasterDictionary(page)
  const dictionary = new MasterDictionaryPage(page)

  await dictionary.search("__no_such_term__")

  await expect(dictionary.emptyRow()).toContainText("処理対象がありません")
})

test("SCN-MDM-002 master dictionary opens selected entry detail", async ({
  page
}) => {
  // 一覧行の利用者操作で、選択した辞書エントリの詳細が表示されることを証明する。
  await openMasterDictionary(page)
  const dictionary = new MasterDictionaryPage(page)

  const secondRowSourceText = await dictionary.entryRows.nth(1).innerText()
  await dictionary.entryRows.nth(1).click()

  await expect(dictionary.detailPanel).toContainText(
    secondRowSourceText.split("\n")[0]
  )
})

test("E2E-UC-008 master dictionary creates a dictionary entry", async ({
  page
}) => {
  // 新規登録の利用者操作で、入力した辞書エントリが一覧と詳細へ反映されることを証明する。
  await openMasterDictionary(page)
  const dictionary = new MasterDictionaryPage(page)

  const sourceText = "System Test Create Entry"
  const createdTranslation = "フェーズ5 作成訳語"

  await dictionary.openCreateModal()
  await expect(dictionary.createEditModal).toContainText("新規登録")
  await dictionary.fillEntry({
    origin: "手動登録",
    source: sourceText,
    translation: createdTranslation
  })
  await clickEditModalSave(dictionary)
  await expect(dictionary.createEditModal).toBeHidden()

  await dictionary.search(sourceText)
  await expect.poll(async () => dictionary.entryRows.count()).toBeGreaterThan(0)
  await expect(dictionary.detailPanel).toContainText(sourceText)
  await expect(dictionary.entryRows.first()).toContainText(createdTranslation)
})

test("E2E-UC-009 master dictionary keeps invalid entry in create modal", async ({
  page
}) => {
  // 必須不足の保存操作で、辞書エントリが登録されず入力 modal が維持されることを証明する。
  await openMasterDictionary(page)
  const dictionary = new MasterDictionaryPage(page)

  await dictionary.openCreateModal()
  await expect(dictionary.createEditModal).toBeVisible()
  await dictionary.fillEntry({ source: "", translation: "" })
  await dictionary.saveEntry()

  await expect(dictionary.createEditModal).toBeVisible()
})

test("E2E-UC-030 master dictionary cancels a create operation", async ({
  page
}) => {
  // 新規登録の取消操作で、入力した候補が一覧へ保存されないことを証明する。
  await openMasterDictionary(page)
  const dictionary = new MasterDictionaryPage(page)

  const sourceText = "System Test Cancel Candidate"

  await dictionary.openCreateModal()
  await dictionary.fillEntry({
    source: sourceText,
    translation: "キャンセル候補"
  })
  await dictionary.closeEntryModal()

  await expect(dictionary.createEditModal).toBeHidden()
  await dictionary.search(sourceText)
  await expect(dictionary.emptyRow()).toContainText("処理対象がありません")
})

test("E2E-UC-010 master dictionary updates a selected entry", async ({
  page
}) => {
  // 更新の利用者操作で、選択した辞書エントリの訳語が更新されることを証明する。
  await openMasterDictionary(page)
  const dictionary = new MasterDictionaryPage(page)

  const sourceText = "System Test Update Entry"
  const updatedTranslation = "フェーズ5 更新訳語"

  await dictionary.openCreateModal()
  await dictionary.fillEntry({
    origin: "手動登録",
    source: sourceText,
    translation: "更新前訳語"
  })
  await clickEditModalSave(dictionary)
  await dictionary.search(sourceText)
  await expect.poll(async () => dictionary.entryRows.count()).toBeGreaterThan(0)

  await dictionary.openEditModal()
  await expect(dictionary.createEditModal).toContainText("更新")
  await dictionary.fillEntry({ translation: updatedTranslation })
  await clickEditModalSave(dictionary)
  await expect(dictionary.createEditModal).toBeHidden()
  await expect(dictionary.entryRows.first()).toContainText(updatedTranslation)
})

test("E2E-UC-011 master dictionary deletes a selected entry", async ({
  page
}) => {
  // 削除確定の利用者操作で、選択した辞書エントリが一覧から消えることを証明する。
  await openMasterDictionary(page)
  const dictionary = new MasterDictionaryPage(page)

  const sourceText = "System Test Delete Entry"

  await dictionary.openCreateModal()
  await dictionary.fillEntry({
    origin: "手動登録",
    source: sourceText,
    translation: "削除候補"
  })
  await clickEditModalSave(dictionary)
  await dictionary.search(sourceText)
  await expect.poll(async () => dictionary.entryRows.count()).toBeGreaterThan(0)

  await dictionary.openDeleteModal()
  await expect(dictionary.deleteModal).toBeVisible()
  await dictionary.confirmDelete()
  await expect(dictionary.deleteModal).toBeHidden()

  await expect.poll(async () => dictionary.entryRows.count()).toBe(0)
  await expect
    .poll(async () => await dictionary.operationRegion.innerText())
    .not.toContain(sourceText)
})

test("E2E-UC-031 master dictionary cancels delete confirmation", async ({
  page
}) => {
  // 削除取消の利用者操作で、選択した辞書エントリが一覧へ残ることを証明する。
  await openMasterDictionary(page)
  const dictionary = new MasterDictionaryPage(page)

  const firstRowSourceText = await dictionary.entryRows.first().innerText()
  await dictionary.entryRows.first().click()
  await dictionary.openDeleteModal()
  await dictionary.cancelDelete()

  await expect(dictionary.deleteModal).toBeHidden()
  await expect(dictionary.operationRegion).toContainText(
    firstRowSourceText.split("\n")[0].trim()
  )
})

test("SCN-MDM-008/009 XML未選択ゲートと取込バー状態遷移を確認できる", async ({
  page
}) => {
  // XML 未選択時は取り込み開始できず、選択後に取り込み状態が完了へ進むことを証明する。
  await openMasterDictionary(page)
  const dictionary = new MasterDictionaryPage(page)

  const importProgressRegion =
    dictionary.xmlImportProgressPanel.getByLabel("XML 取り込みの進行率")

  await expect(importProgressRegion).toBeVisible()
  await expect(dictionary.xmlImportButton).toBeHidden()

  await importDawnguardXml(page)

  await expect(dictionary.searchInput).toHaveValue("")
  await expect(dictionary.categorySelect).toHaveValue("すべて")
  await expect(dictionary.xmlImportProgressPanel).toContainText("完了")
})

test("E2E-UC-012 master dictionary imports allowed records from XML", async ({
  page
}) => {
  // XML 取り込みの利用者操作で、許可された record だけが辞書一覧へ反映されることを証明する。
  await openMasterDictionary(page)
  const dictionary = new MasterDictionaryPage(page)
  await importDawnguardXml(page)

  await dictionary.search("Auriel's Bow")
  await expect(dictionary.entryRows).not.toHaveCount(0)

  await dictionary.search("Crossbow Mount")
  await expect(dictionary.emptyRow()).toContainText("処理対象がありません")

  await dictionary.search("Transform into the vampire lord.")
  await expect(dictionary.emptyRow()).toContainText("処理対象がありません")
})

test("E2E-UC-032 master dictionary keeps existing list when XML import fails", async ({
  page
}) => {
  // 不正 XML の取り込み操作で、失敗表示になり既存一覧が維持されることを証明する。
  await openMasterDictionary(page)
  const dictionary = new MasterDictionaryPage(page)
  const beforeListText = await dictionary.operationRegion.innerText()

  await dictionary.stageXmlFileWithRuntimePath(invalidXmlPath)
  await dictionary.startXmlImport()

  await expect(dictionary.screen).toContainText(/XML syntax error|取り込みに失敗/)
  await expect(dictionary.operationRegion).toContainText(
    beforeListText.split("\n")[0].trim()
  )
})

test("E2E-UC-054 master dictionary disables duplicate XML import while running", async ({
  page
}) => {
  // XML 取り込み中の画面状態で、取り込み開始と選び直しが無効化されることを証明する。
  await openMasterDictionary(page)
  const dictionary = new MasterDictionaryPage(page)

  await dictionary.stageXmlFileWithRuntimePath(dawnguardXmlPath)
  await page.evaluate(() => {
    const controller = globalThis.go?.wails?.AppController
    if (!controller) {
      return
    }
    controller.ImportMasterDictionaryXml = () => new Promise(() => {})
    controller.ImportMasterDictionaryXML = () => new Promise(() => {})
  })
  await dictionary.startXmlImport()

  await expect(dictionary.xmlImportProgressPanel).toContainText(/取込中/)
  await expect(dictionary.xmlImportButton).toBeDisabled()
  await expect(
    dictionary.xmlImportProgressPanel.getByRole("button", {
      name: "選び直す"
    })
  ).toBeDisabled()
})
