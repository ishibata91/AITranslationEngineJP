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
  await installScenarioWailsMocks(page);
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
  // AI 設定不足の開始操作で、単語翻訳段階が未開始状態を維持することを証明する。
  const phase = new TermTranslationPhasePage(page);
  await openPhase(page, "system-test-term", phase);

  await phase.start();

  await expect(phase.aiModelSelection).toContainText(/設定未完了|認証状態/);
  await expect(phase.statusPanel).toContainText(/未開始|開始待ち/);
  await expect(phase.progress).toContainText(/0/);
});

test("E2E-UC-052 persona generation keeps not-started state when AI settings are incomplete", async ({
  page,
}) => {
  // AI 設定不足の開始操作で、NPC ペルソナ生成段階が未開始状態を維持することを証明する。
  const phase = new PersonaGenerationPhasePage(page);
  await openPhase(page, "system-test-persona", phase);

  await phase.start();

  await expect(phase.aiModelSelection).toContainText(/設定未完了|認証状態/);
  await expect(phase.statusPanel).toContainText(/未開始|開始待ち/);
  await expect(phase.progress).toContainText(/0/);
});

test("E2E-UC-053 body translation keeps not-started state when AI settings are incomplete", async ({
  page,
}) => {
  // AI 設定不足の開始操作で、本文翻訳段階が未開始状態を維持することを証明する。
  const phase = new BodyTranslationPhasePage(page);
  await openPhase(page, "system-test-body-pending", phase);

  await phase.start();

  await expect(phase.aiModelSelection).toContainText(/設定未完了|認証状態/);
  await expect(phase.statusPanel).toContainText(/未開始|開始待ち|失敗/);
  await expect(phase.progress).toContainText(/0/);
});
