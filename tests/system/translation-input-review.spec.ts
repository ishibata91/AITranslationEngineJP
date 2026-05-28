import path from "node:path"

import { expect, test } from "@playwright/test"

import { TranslationInputReviewPage } from "./support/system-test-pages"

test.describe.configure({ mode: "serial" })

const validInputPath = path.resolve(
  process.cwd(),
  "tests/fixtures/translation-input/system-test-dialogue.json"
)
const invalidInputPath = path.resolve(
  process.cwd(),
  "tests/fixtures/translation-input/invalid-dialogue.json"
)

test("E2E-UC-021 translation input review registers selected JSON", async ({
  page
}) => {
  // JSON 選択と登録の利用者操作で、登録済み入力が一覧と次操作導線へ反映されることを証明する。
  const inputReview = new TranslationInputReviewPage(page)
  await inputReview.open()

  await inputReview.setJsonFile(validInputPath)
  await expect(inputReview.registerButton).toBeEnabled()
  await inputReview.register()

  await expect(inputReview.loadedInputList).toContainText(
    "system-test-dialogue.json"
  )
  await expect(inputReview.nextActionFooter).toContainText(/翻訳設定|次/)
})

test("E2E-UC-040 translation input review rejects invalid JSON shape", async ({
  page
}) => {
  // 不正 JSON の登録操作で、登録失敗理由が表示され次操作導線が有効化されないことを証明する。
  const inputReview = new TranslationInputReviewPage(page)
  await inputReview.open()

  await inputReview.setJsonFile(invalidInputPath)
  await expect(inputReview.registerButton).toBeEnabled()
  await inputReview.register()

  await expect(inputReview.loadedInputList).toContainText("登録失敗")
  await expect(inputReview.loadedInputList).toContainText(
    /missing required field|必須/
  )
  await expect(inputReview.nextActionFooter).toHaveCount(0)
})
