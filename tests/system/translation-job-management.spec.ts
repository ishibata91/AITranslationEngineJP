import { expect, test, type Locator, type Page } from "@playwright/test"

const JOB_MANAGEMENT_URL = "/#translation-management"

const FORBIDDEN_SECRET_TEXTS = [
  "api key",
  "apikey",
  "token",
  "credential-ref",
  "provider raw",
  "external response"
]

async function openJobManagement(page: Page) {
  await page.goto(JOB_MANAGEMENT_URL)
  await expect(
    page.getByRole("heading", { level: 1, name: "翻訳管理" })
  ).toBeVisible()
  await expect(
    page.getByRole("heading", { level: 2, name: "未完了ジョブ一覧" })
  ).toBeVisible()
  await expect(
    page.getByTestId("translation-job-management-job-list-region")
  ).toBeVisible()
}

function jobCards(page: Page) {
  return page.getByTestId("translation-job-management-job-card")
}

async function jobIdText(card: Locator): Promise<string> {
  return (await card.locator(".job-card-id").innerText()).trim()
}

function cardActions(card: Locator) {
  return card.getByTestId("translation-job-management-job-actions")
}

test("SCN-TJM-001 translation-job-management lists incomplete jobs and excludes completed", async ({
  page
}) => {
  // 実 backend の seeded DB から、未完了 job だけが一覧化されることを証明する。
  await openJobManagement(page)

  const cards = jobCards(page)
  await expect.poll(async () => cards.count()).toBeGreaterThan(0)
  await expect(cards.first().locator(".job-card-id")).toContainText(
    /ジョブ #\d+/
  )
  await expect(page.getByText("Completed")).toHaveCount(0)

  await expect
    .poll(async () =>
      page.locator(".state-badge", { hasText: "実行前" }).count()
    )
    .toBeGreaterThan(0)
  await expect
    .poll(async () =>
      page.locator(".state-badge", { hasText: "実行中" }).count()
    )
    .toBeGreaterThan(0)
  await expect(
    page.getByRole("combobox", { name: "状態フィルタ" })
  ).toContainText("再開可能な失敗")
})

test("SCN-TJM-003 translation-job-management opens selected job in Job Run", async ({
  page
}) => {
  // 実 backend の一覧で選んだ job が、Job Run の選択状態へ引き継がれることを証明する。
  await openJobManagement(page)

  const firstCard = jobCards(page).first()
  const selectedJobIdText = await jobIdText(firstCard)
  await cardActions(firstCard)
    .getByRole("button", { name: "現在の翻訳段階へ進む" })
    .click()

  await expect(page.getByTestId("job-run-job-run-shell")).toBeVisible()
  await expect(
    page.getByRole("heading", { name: selectedJobIdText })
  ).toBeVisible()
  await expect(
    page
      .getByTestId("job-run-selected-job-summary")
      .getByText(selectedJobIdText)
  ).toBeVisible()
  await expect(page.getByText("Running")).toHaveCount(0)
})

test("SCN-TJM-005 and SCN-TJM-007 translation-job-management shows operation entries and blocked reasons", async ({
  page
}) => {
  // 実 backend の job 状態に応じて、操作可否と理由が表示されることを証明する。
  await openJobManagement(page)

  const runningCard = jobCards(page).filter({ hasText: "実行中" }).first()
  await expect(
    cardActions(runningCard).getByRole("button", { name: "停止" })
  ).toBeEnabled()
  await expect(
    cardActions(runningCard).getByRole("button", { name: "再開" })
  ).toBeDisabled()
  await expect(
    runningCard.getByTestId("translation-job-management-disabled-reason")
  ).toContainText("再開:")

  const readyCard = jobCards(page).filter({ hasText: "実行前" }).first()
  await expect(
    cardActions(readyCard).getByRole("button", {
      name: "現在の翻訳段階へ進む"
    })
  ).toBeEnabled()
  await expect(
    cardActions(readyCard).getByRole("button", { name: "停止" })
  ).toBeDisabled()
  await expect(
    cardActions(readyCard).getByRole("button", { name: "再開" })
  ).toBeDisabled()
  await expect(
    readyCard.getByTestId("translation-job-management-disabled-reason")
  ).toContainText("停止:")
  await expect(
    readyCard.getByTestId("translation-job-management-disabled-reason")
  ).toContainText("再開:")
  await expect(readyCard.getByRole("button", { name: "削除" })).toBeEnabled()
})

test("SCN-TJM-009 translation-job-management does not expose secret text in UI or console", async ({
  page
}) => {
  // 実 backend の read model が、秘密値に相当する文字列を画面と console へ出さないことを証明する。
  const consoleMessages: string[] = []
  page.on("console", (message) => {
    consoleMessages.push(message.text())
  })

  await openJobManagement(page)

  const visibleText = (await page.locator("body").innerText()).toLowerCase()
  const consoleText = consoleMessages.join("\n").toLowerCase()

  for (const forbiddenText of FORBIDDEN_SECRET_TEXTS) {
    expect(visibleText).not.toContain(forbiddenText)
    expect(consoleText).not.toContain(forbiddenText)
  }
})
