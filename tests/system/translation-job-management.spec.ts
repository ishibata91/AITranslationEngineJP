import { expect, test, type Page } from "@playwright/test"

const JOB_MANAGEMENT_URL =
  "/?fakeApi=1&fakeScenario=success#translation-management"

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
  await page
    .getByRole("tab", {
      name: /ジョブ管理/
    })
    .click()
  await expect(page.getByRole("heading", { name: "ジョブ管理" })).toBeVisible()
}

function jobCard(page: Page, jobId: number) {
  return page.getByRole("button", { name: `Job ${jobId} を選択` })
}

test("SCN-TJM-001 translation-job-management lists incomplete jobs and excludes completed", async ({
  page
}) => {
  await openJobManagement(page)

  for (const jobId of [401, 402, 403, 404, 405, 406]) {
    await expect(jobCard(page, jobId).getByText(`Job #${jobId}`)).toBeVisible()
  }
  await expect(page.getByText("Completed")).toHaveCount(0)

  await expect(page.locator(".state-badge", { hasText: "実行前" })).toBeVisible()
  await expect(page.locator(".state-badge", { hasText: "実行中" })).toBeVisible()
  await expect(page.locator(".state-badge", { hasText: "中断中" })).toBeVisible()
  await expect(
    page.locator(".state-badge", { hasText: "再開可能な失敗" })
  ).toBeVisible()
  await expect(
    page.locator(".state-badge", { hasText: "回復できない失敗" })
  ).toBeVisible()
  await expect(
    page.locator(".state-badge", { hasText: "キャンセル済み" })
  ).toBeVisible()
})

test("SCN-TJM-003 translation-job-management opens selected job in Job Run without state mutation", async ({
  page
}) => {
  await openJobManagement(page)

  await page
    .getByLabel("Job 403 の操作")
    .getByRole("button", { name: "再開" })
    .click()

  await expect(page.getByRole("heading", { name: "Job Run" }).first()).toBeVisible()
  await expect(page.getByRole("heading", { name: "Job #403" })).toBeVisible()
  await expect(page.locator("#termPhaseJobIdInput")).toHaveValue("403")
  const targetSummary = page.locator(".job-run-target-summary")
  await expect(targetSummary.getByText("中断中 / 中断中")).toBeVisible()
  await expect(targetSummary.getByText("NPC ペルソナ生成")).toBeVisible()
  await expect(
    targetSummary.getByText("44% / 再開要求を受け付けました")
  ).toBeVisible()
  await expect(page.getByText("Running")).toHaveCount(0)
})

test("SCN-TJM-005 and SCN-TJM-007 translation-job-management shows resume entry and blocked reasons", async ({
  page
}) => {
  await openJobManagement(page)

  await expect(
    page.getByLabel("Job 403 の操作").getByRole("button", { name: "再開" })
  ).toBeEnabled()

  const cacheMissingResume = page
    .getByLabel("Job 404 の操作")
    .getByRole("button", { name: "再開" })
  await expect(cacheMissingResume).toBeDisabled()
  await expect(cacheMissingResume.locator("..")).toHaveAttribute(
    "data-tooltip",
    /入力キャッシュ/
  )

  const projectionFailureResume = page
    .getByLabel("Job 405 の操作")
    .getByRole("button", { name: "再開" })
  await expect(projectionFailureResume).toBeDisabled()
  await expect(projectionFailureResume.locator("..")).toHaveAttribute(
    "data-tooltip",
    /進捗を確認できません/
  )

  const terminalResume = page
    .getByLabel("Job 406 の操作")
    .getByRole("button", { name: "再開" })
  await expect(terminalResume).toBeDisabled()
  await expect(terminalResume.locator("..")).toHaveAttribute(
    "data-tooltip",
    /terminal state/
  )
})

test("SCN-TJM-009 translation-job-management does not expose secret text in UI or console", async ({
  page
}) => {
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
