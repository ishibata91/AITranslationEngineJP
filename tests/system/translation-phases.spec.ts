import { expect, test, type Page } from "@playwright/test";

import {
  BodyTranslationPhasePage,
  PersonaGenerationPhasePage,
  TermTranslationPhasePage,
  TranslationJobManagementPage,
  type TranslationPhasePage,
} from "./support/system-test-pages";
import { installScenarioWailsMocks } from "./support/scenario-wails-mocks";

test.beforeEach(async ({ page }) => {
  // E2E-UC-051/052/053: 各段階で AI 設定が未完了の状態を再現する。
  // デフォルト execution は provider="-" を返し isExecutionConfigured が truthy になるため
  // pill「固定済み」になってしまう。各 seeded job を未設定 execution で返すオプションを渡す。
  await installScenarioWailsMocks(page, {
    seededTermJobMissingExecution: true,
    seededPersonaJobMissingExecution: true,
    seededBodyJobMissingExecution: true,
  });
});

async function openPhase(
  page: Page,
  jobText: string,
  phase: TranslationPhasePage,
): Promise<void> {
  const management = new TranslationJobManagementPage(page);
  await management.open();
  const card = management.jobCard(jobText);
  await expect(card, `job card is visible: ${jobText}`).toBeVisible();
  await expect(
    management.openCurrentPhaseButton(card),
    `current phase action is available: ${jobText}`,
  ).toBeEnabled();
  await management.openCurrentPhase(card);
  await phase.waitForScreen();
}

test("E2E-UC-051 term translation keeps not-started state when AI settings are incomplete", async ({
  page,
}) => {
  // AI 設定不足のとき、単語翻訳段階の開始ボタンが disabled になることを証明する。
  // 新設計では presenter が projection.aiSettingsConfigured=false から canStart=false を導出し、
  // start ボタンを disabled にする。disabled なボタンは利用者が操作できないため、未開始状態が維持される。
  const phase = new TermTranslationPhasePage(page);
  await openPhase(page, "system-test-term", phase);

  // 開始操作が disabled で阻止されていることを確認する（AI 設定未完了の証拠）。
  await expect(phase.startButton).toBeDisabled();
  await expect(phase.aiModelSelection).toContainText(/設定未完了|認証状態/);
  await expect(phase.statusPanel).toContainText(/未開始|開始待ち/);
  await expect(phase.progress).toContainText(/0/);
});

test("E2E-UC-052 persona generation keeps not-started state when AI settings are incomplete", async ({
  page,
}) => {
  // AI 設定不足のとき、NPC ペルソナ生成段階の開始ボタンが disabled になることを証明する。
  // 新設計では presenter が projection.aiSettingsConfigured=false から canStart=false を導出し、
  // start ボタンを disabled にする。disabled なボタンは利用者が操作できないため、未開始状態が維持される。
  const phase = new PersonaGenerationPhasePage(page);
  await openPhase(page, "system-test-persona", phase);

  // 開始操作が disabled で阻止されていることを確認する（AI 設定未完了の証拠）。
  await expect(phase.startButton).toBeDisabled();
  await expect(phase.aiModelSelection).toContainText(/設定未完了|認証状態/);
  await expect(phase.statusPanel).toContainText(/未開始|開始待ち/);
  await expect(phase.progress).toContainText(/0/);
});

test("E2E-UC-053 body translation keeps not-started state when AI settings are incomplete", async ({
  page,
}) => {
  // AI 設定不足のとき、本文翻訳段階の開始ボタンが disabled になることを証明する。
  // 新設計では presenter が projection.aiSettingsConfigured=false から canStart=false を導出し、
  // start ボタンを disabled にする。disabled なボタンは利用者が操作できないため、未開始状態が維持される。
  const phase = new BodyTranslationPhasePage(page);
  await openPhase(page, "system-test-body-pending", phase);

  // 開始操作が disabled で阻止されていることを確認する（AI 設定未完了の証拠）。
  await expect(phase.startButton).toBeDisabled();
  await expect(phase.aiModelSelection).toContainText(/設定未完了|認証状態/);
  await expect(phase.statusPanel).toContainText(/未開始|開始待ち|失敗/);
  await expect(phase.progress).toContainText(/0/);
});
