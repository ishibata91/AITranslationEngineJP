import { expect, test, type Page } from "@playwright/test";

import {
  JobRunShellPage,
  TranslationJobManagementPage,
} from "./support/system-test-pages";
import { installScenarioWailsMocks } from "./support/scenario-wails-mocks";

test.beforeEach(async ({ page }) => {
  await installScenarioWailsMocks(page);
});

async function openJobRun(
  page: Page,
  jobText: string,
): Promise<JobRunShellPage> {
  const management = new TranslationJobManagementPage(page);
  const jobRun = new JobRunShellPage(page);
  await management.open();
  const card = management.jobCard(jobText);
  await expect(card, `job card is visible: ${jobText}`).toBeVisible();
  await expect(
    management.openCurrentPhaseButton(card),
    `current phase action is available: ${jobText}`,
  ).toBeEnabled();
  await management.openCurrentPhase(card);
  await expect(jobRun.shell).toBeVisible();
  await expect(jobRun.phaseScreenRegion).toBeVisible();
  return jobRun;
}

test("E2E-UC-045 opens term phase through current phase action", async ({
  page,
}) => {
  // 承認済み現行読み替え: 未完了 job から単語翻訳の現在段階を開けることを証明する。
  const jobRun = await openJobRun(page, "system-test-term");

  await expect(jobRun.selectedJobSummary).toContainText("ジョブ #7");
  await expect(jobRun.phaseScreenRegion).toContainText(
    /単語翻訳|開始待ち|未開始/,
  );
  await expect(jobRun.phaseScreenRegion).toContainText(/0|1/);
});

test("E2E-UC-046 opens persona generation phase through current phase action", async ({
  page,
}) => {
  // 承認済み現行読み替え: 未完了 job から NPC ペルソナ生成段階を開けることを証明する。
  const jobRun = await openJobRun(page, "system-test-persona");

  await expect(jobRun.selectedJobSummary).toContainText("ジョブ #8");
  await expect(jobRun.phaseScreenRegion).toContainText(
    /NPC ペルソナ生成|開始待ち|未開始/,
  );
  await expect(jobRun.phaseScreenRegion).toContainText(/0|1/);
});

test("E2E-UC-047 opens body translation phase through current phase action", async ({
  page,
}) => {
  // 承認済み現行読み替え: 未完了 job から本文翻訳段階を開けることを証明する。
  const jobRun = await openJobRun(page, "system-test-body-pending");

  await expect(jobRun.selectedJobSummary).toContainText("ジョブ #9");
  await expect(jobRun.phaseScreenRegion).toContainText(
    /本文翻訳|開始待ち|未開始/,
  );
  await expect(jobRun.phaseScreenRegion).toContainText(/0|1/);
});

test("E2E-UC-048 job run shell advances from completed term phase to persona phase", async ({
  page,
}) => {
  // 完了済み単語翻訳段階の次へ進む操作で、NPC ペルソナ生成段階へ遷移することを証明する。
  const jobRun = await openJobRun(page, "system-test-completed-term");

  await jobRun.clickNext();

  await expect(jobRun.phaseScreenRegion).toContainText("NPC ペルソナ生成");
  await expect(jobRun.selectedJobSummary).toContainText("ジョブ #10");
});

test("E2E-UC-049 job run shell advances from completed body phase to completion page", async ({
  page,
}) => {
  // 完了済み本文翻訳段階の次へ進む操作で、翻訳完了確認へ遷移することを証明する。
  const jobRun = await openJobRun(
    page,
    "system-test-body-ready-for-completion",
  );

  await jobRun.clickBodyCompleteNext();

  await expect(jobRun.translationCompleteScreen).toBeVisible();
  await expect(jobRun.translationCompleteScreen).toContainText(
    /あなたの荷物|I am sworn/,
  );
  await expect(jobRun.postCompletionNextAction).toContainText("出力管理へ進む");
});

test("E2E-UC-050 job run shell keeps phase when next conditions are incomplete", async ({
  page,
}) => {
  // 次段階条件不足では、次へ進む操作が無効で現在段階が維持されることを証明する。
  const jobRun = await openJobRun(page, "system-test-term");

  await expect(
    jobRun.nextActionFooter.getByRole("button", { name: "次へ進む" }),
  ).toBeDisabled();
  await expect(jobRun.nextActionFooter).toContainText(
    /完了していません|次へ進めません/,
  );
  await expect(jobRun.phaseScreenRegion).toContainText(
    /単語翻訳|開始待ち|未開始/,
  );
});
