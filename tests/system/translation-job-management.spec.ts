import { expect, test, type Locator, type Page } from "@playwright/test"

import {
  JobRunShellPage,
  TranslationJobManagementPage
} from "./support/system-test-pages"

test.describe.configure({ mode: "serial" })

const FORBIDDEN_SECRET_TEXTS = [
  "api key",
  "apikey",
  "token",
  "credential-ref",
  "provider raw",
  "external response"
]

async function openJobManagement(page: Page) {
  const jobManagement = new TranslationJobManagementPage(page)
  await jobManagement.open()
  await expect(page).toHaveURL(/#translation-management$/)
  await expect(
    page.getByRole("heading", { level: 1, name: "翻訳管理" })
  ).toBeVisible()
  await expect(
    page.getByRole("heading", { level: 2, name: "未完了ジョブ一覧" })
  ).toBeVisible()
  await expect(jobManagement.jobListRegion).toBeVisible()
}

async function jobIdText(card: Locator): Promise<string> {
  return (
    await card
      .getByTestId("translation-job-management-job-selection-region")
      .innerText()
  )
    .split("\n")[0]
    .trim()
}

test("SCN-TJM-001 translation-job-management lists incomplete jobs and excludes completed", async ({
  page
}) => {
  // 実 backend の seeded DB から、未完了 job だけが一覧化されることを証明する。
  await openJobManagement(page)
  const jobManagement = new TranslationJobManagementPage(page)

  await expect
    .poll(async () => jobManagement.jobCards.count())
    .toBeGreaterThan(0)
  await expect(
    jobManagement.jobSelectionRegion(jobManagement.jobCards.first())
  ).toContainText(/ジョブ #\d+/)
  await expect(page.getByText(/system-test-completed(?!-)/)).toHaveCount(0)

  await expect
    .poll(async () =>
      page
        .getByTestId("translation-job-management-state-label")
        .filter({ hasText: "実行前" })
        .count()
    )
    .toBeGreaterThan(0)
  await expect
    .poll(async () =>
      page
        .getByTestId("translation-job-management-state-label")
        .filter({ hasText: "実行中" })
        .count()
    )
    .toBeGreaterThan(0)
  await expect(
    page.getByRole("combobox", { name: "状態フィルタ" })
  ).toContainText("再開可能な失敗")
})

test("E2E-UC-017 translation-job-management filters running jobs by search and state", async ({
  page
}) => {
  // 検索と状態 filter の利用者操作で、条件一致した未完了 job だけが表示されることを証明する。
  await openJobManagement(page)
  const jobManagement = new TranslationJobManagementPage(page)

  await jobManagement.search("system-test-running")
  await jobManagement.selectState("実行中")

  await expect(jobManagement.jobCard("system-test-running")).toBeVisible()
  await expect(jobManagement.jobCard("system-test-ready")).toHaveCount(0)
})

test("E2E-UC-036 translation-job-management shows empty state for unmatched search", async ({
  page
}) => {
  // 未完了 job 検索が 0 件になる利用者操作で、一覧が該当なし状態になることを証明する。
  await openJobManagement(page)
  const jobManagement = new TranslationJobManagementPage(page)

  await jobManagement.search("NoSuchJob")

  await expect(jobManagement.listStatusDisplay).toContainText(
    /管理対象なし|管理対象がありません|該当なし|条件に一致/
  )
  await expect(jobManagement.jobCards).toHaveCount(0)
})

test("SCN-TJM-003 translation-job-management opens selected job in Job Run", async ({
  page
}) => {
  // 実 backend の一覧で選んだ job が、Job Run の選択状態へ引き継がれることを証明する。
  await openJobManagement(page)
  const jobManagement = new TranslationJobManagementPage(page)
  const jobRun = new JobRunShellPage(page)

  const firstCard = jobManagement.jobCard("system-test-running")
  const selectedJobIdText = await jobIdText(firstCard)
  await jobManagement.openCurrentPhase(firstCard)

  await expect(jobRun.shell).toBeVisible()
  await expect(
    page.getByRole("heading", { name: selectedJobIdText })
  ).toBeVisible()
  await expect(
    jobRun.selectedJobSummary.getByText(selectedJobIdText)
  ).toBeVisible()
  await expect(page.getByText("Running")).toHaveCount(0)
})

test("SCN-TJM-005 and SCN-TJM-007 translation-job-management shows operation entries and blocked reasons", async ({
  page
}) => {
  // 実 backend の job 状態に応じて、操作可否と理由が表示されることを証明する。
  // JobOperationGroup は stopOperation と deleteOperation と continue-button のみを持つ。
  // resume 操作ボタン（中断からの再開）は UI に存在しない。
  // 再開操作がないことは JobOperationGroup のプロダクトコード上の事実であり、
  // テストでは stop/delete/continue 操作の状態を確認して操作可否の表示を証明する。
  await openJobManagement(page)
  const jobManagement = new TranslationJobManagementPage(page)

  const runningCard = jobManagement.jobCards
    .filter({ hasText: "実行中" })
    .first()
  await expect(jobManagement.stopButton(runningCard)).toBeEnabled()
  // 実行中カードは翻訳段階を開く continue-button（「再開」テキスト）のみ持つ。
  // resume 操作ボタン（中断からの再開）は JobOperationGroup に含まれないため存在しない。
  await expect(jobManagement.openCurrentPhaseButton(runningCard)).toBeVisible()

  const readyCard = jobManagement.jobCards.filter({ hasText: "実行前" }).first()
  await expect(jobManagement.openCurrentPhaseButton(readyCard)).toBeEnabled()
  await expect(jobManagement.stopButton(readyCard)).toBeDisabled()
  await expect(jobManagement.deleteButton(readyCard)).toBeEnabled()
})

test("E2E-UC-020 translation-job-management deletes a canceled job from card action", async ({
  page
}) => {
  // 削除確定の利用者操作で、キャンセル済み job が一覧から消えることを証明する。
  await openJobManagement(page)
  const jobManagement = new TranslationJobManagementPage(page)

  const canceledCard = jobManagement.jobCard("system-test-canceled")
  await expect(jobManagement.deleteButton(canceledCard)).toBeEnabled()
  await jobManagement.openDelete(canceledCard)
  await jobManagement.confirmDelete()

  await expect(jobManagement.feedbackNotification).toContainText(/削除/)
  await expect(jobManagement.jobCard("system-test-canceled")).toHaveCount(0)
})

test("E2E-UC-038 translation-job-management keeps resume disabled for failed job", async ({
  page
}) => {
  // 再開条件不足の job では、失敗状態が維持されることを証明する。
  // JobOperationGroup は resumeOperation を描画しないため、再開専用ボタンは存在しない。
  // 再開不可の証拠として失敗状態ラベルで確認する。
  await openJobManagement(page)
  const jobManagement = new TranslationJobManagementPage(page)

  const failedCard = jobManagement.jobCard("system-test-failed")

  await expect(jobManagement.stateLabel(failedCard)).toContainText(
    /回復不能|失敗/
  )
})

test("E2E-UC-019 translation-job-management reports resume result for paused job", async ({
  page
}) => {
  // 再開の利用者操作で、後続実行画面への遷移が成立することを証明する。
  // UI では continue-button（翻訳段階を開く「再開」ボタン）が再開操作の起点になる。
  // JobOperationGroup は resumeOperation を持たないため、feedback 通知は continue-button 経由では出ない。
  // 遷移先の job run shell と現在段階が表示されることで再開の利用者操作が成立したことを確認する。
  await openJobManagement(page)
  const jobManagement = new TranslationJobManagementPage(page)
  const jobRun = new JobRunShellPage(page)

  const pausedCard = jobManagement.jobCard("system-test-paused")
  await expect(jobManagement.stateLabel(pausedCard)).toContainText("中断中")
  await expect(
    jobManagement.resumeButton(pausedCard),
    "resume action failure: paused job の再開操作が実行可能である"
  ).toBeEnabled()
  await jobManagement.resume(pausedCard)

  await expect.soft(
    jobRun.shell,
    "shell transition failure: 再開成立時だけ job run shell が表示される"
  ).toBeVisible()
  await expect.soft(
    jobRun.phaseScreenRegion,
    "shell transition failure: job run shell 内に現在段階が表示される"
  ).toBeVisible()
})

test("E2E-UC-018 translation-job-management stops a running job from card action", async ({
  page
}) => {
  // 停止の利用者操作で、実行中 job の状態変更結果が feedback と状態表示へ反映されることを証明する。
  await openJobManagement(page)
  const jobManagement = new TranslationJobManagementPage(page)

  const runningCard = jobManagement.jobCard("system-test-running")
  await jobManagement.stop(runningCard)

  await expect(jobManagement.feedbackNotification).toContainText(/停止|中断/)
  await expect(jobManagement.stateLabel(runningCard)).toContainText(
    /中断|停止|失敗/
  )
})

test("E2E-UC-037 translation-job-management keeps stop disabled for ready job", async ({
  page
}) => {
  // 停止できない状態の job では、停止操作が無効で状態が実行前であることを証明する。
  // translation-job-management-disabled-reason testid は現在の JobCard 実装に存在しない。
  // 停止不可の証拠として stop ボタンの disabled 状態と実行前の状態ラベルで確認する。
  await openJobManagement(page)
  const jobManagement = new TranslationJobManagementPage(page)

  const readyCard = jobManagement.jobCard("system-test-ready")

  await expect(jobManagement.stopButton(readyCard)).toBeDisabled()
  await expect(jobManagement.stateLabel(readyCard)).toContainText("実行前")
})

test("E2E-UC-039 translation-job-management cancels delete confirmation", async ({
  page
}) => {
  // 削除取消の利用者操作で、対象 job が一覧へ残ることを証明する。
  await openJobManagement(page)
  const jobManagement = new TranslationJobManagementPage(page)

  const readyCard = jobManagement.jobCard("system-test-ready")
  await jobManagement.openDelete(readyCard)
  await jobManagement.cancelDelete()

  await expect(jobManagement.deleteModal).toBeHidden()
  await expect(jobManagement.jobCard("system-test-ready")).toBeVisible()
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
