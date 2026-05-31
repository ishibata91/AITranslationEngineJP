import path from "node:path";

import { expect, test, type Page } from "@playwright/test";

import {
  BodyTranslationPhasePage,
  JobRunShellPage,
  PersonaGenerationPhasePage,
  TermTranslationPhasePage,
  TranslationInputReviewPage,
  TranslationJobManagementPage,
  type TranslationPhasePage,
} from "./support/system-test-pages";
import { installScenarioWailsMocks } from "./support/scenario-wails-mocks";

interface ProcessingTargetListScenario {
  expectedProgressTargetDetail: RegExp;
  expectedInitialRows: number;
  expectedInitialTotal: string;
  expectedMatchingTotal: string;
  expectedRowText: RegExp | string;
  expectedScreenTargetMetric: RegExp;
  hiddenAfterMatchText: RegExp | string;
  matchingQuery: string;
  noResultQuery: string;
}

const lucienInputPath = path.resolve(
  process.cwd(),
  "dictionaries/Lucien.esp_Export.json",
);

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

async function openJobRunFromLucienDataLoad(
  page: Page,
): Promise<JobRunShellPage> {
  const inputReview = new TranslationInputReviewPage(page);
  const jobRun = new JobRunShellPage(page);

  await inputReview.open();
  await inputReview.setJsonFile(lucienInputPath);
  await expect(inputReview.registerButton).toBeEnabled({ timeout: 15000 });
  await inputReview.register();
  await expect(inputReview.loadedInputList).toContainText(
    "Lucien.esp_Export.json",
  );
  await expect(inputReview.nextActionFooter).toBeVisible();
  await inputReview.nextActionFooter
    .getByRole("button", { name: "単語翻訳へ進む" })
    .click();
  await expect(jobRun.shell).toBeVisible();
  await expect(jobRun.phaseScreenRegion).toBeVisible();
  return jobRun;
}

async function expectProcessingTargetListScenario(
  phase: TranslationPhasePage,
  scenario: ProcessingTargetListScenario,
): Promise<void> {
  await phase.waitForScreen();
  await expect(phase.processingTargetListRegion).toBeVisible();
  await expect(phase.processingTargetTotalCount).toHaveText(
    scenario.expectedInitialTotal,
  );
  await expect(phase.statusPanel).toContainText(
    scenario.expectedScreenTargetMetric,
  );
  await expect(phase.progress).toContainText(
    scenario.expectedProgressTargetDetail,
  );
  await expect(phase.processingTargetRows).toHaveCount(
    scenario.expectedInitialRows,
  );
  await expect(phase.processingTargetRows.first()).toContainText(
    scenario.expectedRowText,
  );

  await phase.searchProcessingTargets(scenario.matchingQuery);

  await expect(phase.processingTargetTotalCount).toHaveText(
    scenario.expectedMatchingTotal,
  );
  await expect(phase.processingTargetRows).toHaveCount(1);
  await expect(phase.processingTargetRows.first()).toContainText(
    scenario.expectedRowText,
  );
  await expect(phase.processingTargetRows).not.toContainText(
    scenario.hiddenAfterMatchText,
  );

  await phase.searchProcessingTargets(scenario.noResultQuery);

  await expect(phase.processingTargetTotalCount).toHaveText("0 件");
  await expect(phase.processingTargetRows).toHaveCount(0);
  await expect(phase.processingTargetEmptyState).toContainText(
    "処理対象がありません",
  );
  await expect(phase.statusPanel).toContainText(
    scenario.expectedScreenTargetMetric,
  );
  await expect(phase.progress).toContainText(
    scenario.expectedProgressTargetDetail,
  );
}

test("E2E-UC-045 opens term phase through current phase action", async ({
  page,
}) => {
  // 単語翻訳段階で、処理対象一覧の件数、行表示、検索結果を利用者操作から証明する。
  const jobRun = await openJobRun(page, "system-test-term");
  const phase = new TermTranslationPhasePage(page);

  await expect(jobRun.selectedJobSummary).toContainText("ジョブ #7");
  await expect(jobRun.phaseScreenRegion).toContainText(
    /単語翻訳|開始待ち|未開始/,
  );
  await expectProcessingTargetListScenario(phase, {
    expectedProgressTargetDetail: /AI 翻訳対象語件数\s*3 件/,
    expectedInitialRows: 3,
    expectedInitialTotal: "1-3 / 3 件",
    expectedMatchingTotal: "1-1 / 1 件",
    expectedRowText: /Dragonborn|ドラゴンボーン/,
    expectedScreenTargetMetric: /AI 翻訳対象語\s*3/,
    hiddenAfterMatchText: "Whiterun Guard",
    matchingQuery: "Dragonborn",
    noResultQuery: "NoSuchTermTarget",
  });
});

test("E2E-DIFF-LUCIEN-001 keeps processing target list non-empty when term progress target is non-zero", async ({
  page,
}) => {
  // データロードで登録した Lucien の単語翻訳画面が、対象ありの進捗と処理対象一覧を同時に表示することを証明する。
  const jobRun = await openJobRunFromLucienDataLoad(page);
  const phase = new TermTranslationPhasePage(page);

  await expect(jobRun.selectedJobSummary).toContainText(/ジョブ #\d+/);
  await expect(jobRun.phaseScreenRegion).toContainText(
    /単語翻訳|開始待ち|未開始/,
  );
  await phase.waitForScreen();
  await expect(phase.processingTargetListRegion).toBeVisible();
  await expect(phase.progress).toContainText(
    /AI 翻訳対象語件数\s*[1-9]\d*\s*件/,
  );
  await expect(phase.processingTargetTotalCount).toContainText(/[1-9]\d*/);
  await expect(phase.processingTargetEmptyState).toHaveCount(0);
  await expect
    .poll(async () => phase.processingTargetRows.count())
    .toBeGreaterThan(0);
});

test("E2E-UC-046 opens persona generation phase through current phase action", async ({
  page,
}) => {
  // NPC ペルソナ生成段階で、処理対象一覧の件数、行表示、検索結果を利用者操作から証明する。
  const jobRun = await openJobRun(page, "system-test-persona");
  const phase = new PersonaGenerationPhasePage(page);

  await expect(jobRun.selectedJobSummary).toContainText("ジョブ #8");
  await expect(jobRun.phaseScreenRegion).toContainText(
    /NPC ペルソナ生成|開始待ち|未開始/,
  );
  await expectProcessingTargetListScenario(phase, {
    expectedProgressTargetDetail: /対象件数\s*2 件/,
    expectedInitialRows: 2,
    expectedInitialTotal: "1-2 / 2 件",
    expectedMatchingTotal: "1-1 / 1 件",
    expectedRowText: /Lydia|FemaleEvenToned|Warrior/,
    expectedScreenTargetMetric: /対象\s*2/,
    hiddenAfterMatchText: "Hadvar",
    matchingQuery: "FemaleEvenToned",
    noResultQuery: "NoSuchPersonaTarget",
  });
});

test("E2E-UC-047 opens body translation phase through current phase action", async ({
  page,
}) => {
  // 本文翻訳段階で、処理対象一覧の件数、行表示、検索結果を利用者操作から証明する。
  const jobRun = await openJobRun(page, "system-test-body-pending");
  const phase = new BodyTranslationPhasePage(page);

  await expect(jobRun.selectedJobSummary).toContainText("ジョブ #9");
  await expect(jobRun.phaseScreenRegion).toContainText(
    /本文翻訳|開始待ち|未開始/,
  );
  await expectProcessingTargetListScenario(phase, {
    expectedProgressTargetDetail: /AI 送信対象件数\s*4 件/,
    expectedInitialRows: 4,
    expectedInitialTotal: "1-4 / 4 件",
    expectedMatchingTotal: "1-1 / 1 件",
    expectedRowText: /I am sworn to carry your burdens|あなたの荷物/,
    expectedScreenTargetMetric: /AI 送信対象\s*4/,
    hiddenAfterMatchText: "You are finally awake",
    matchingQuery: "burdens",
    noResultQuery: "NoSuchBodyTarget",
  });
});

test("E2E-UC-048 job run shell advances from completed term phase to persona phase", async ({
  page,
}) => {
  // 完了済み単語翻訳段階の次へ進む操作で、NPC ペルソナ生成段階へ遷移することを証明する。
  const jobRun = await openJobRun(page, "system-test-completed-term");
  const termPhase = new TermTranslationPhasePage(page);
  // 段階切り替え時は initialFetchDone=false のオーバーレイが表示される。
  // term 段階の初回取得完了（オーバーレイ消失）を待ってから次へ進む操作を行う。
  await termPhase.waitForScreen();

  await jobRun.clickNext();

  // 遷移先の persona 段階が表示されるまで待つ。
  const personaPhase = new PersonaGenerationPhasePage(page);
  await personaPhase.waitForScreen();
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
  const bodyPhase = new BodyTranslationPhasePage(page);
  // 段階切り替え時は initialFetchDone=false のオーバーレイが表示される。
  // body 段階の初回取得完了（オーバーレイ消失）を待ってから次へ進む操作を行う。
  await bodyPhase.waitForScreen();

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
