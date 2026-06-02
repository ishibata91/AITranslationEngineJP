import { expect, test, type Page } from "@playwright/test";

import {
  TermTranslationPhasePage,
  TranslationJobManagementPage,
} from "./support/system-test-pages";
import { installScenarioWailsMocks } from "./support/scenario-wails-mocks";

test.beforeEach(async ({ page }) => {
  await installScenarioWailsMocks(page, { termAISettingsMissing: true });
});

/**
 * 前提: AI 設定未完了（execution の provider/model/executionMode が空文字）のジョブを
 * 翻訳実行画面で開く。
 */
async function openTermPhaseMissingAISettings(page: Page): Promise<TermTranslationPhasePage> {
  // Arrange: 管理画面から AI 設定未完了ジョブを開く。
  const management = new TranslationJobManagementPage(page);
  const phase = new TermTranslationPhasePage(page);
  await management.open();
  const card = management.jobCard("system-test-term-ai-missing");
  await expect(card, "AI 設定未完了ジョブのカードが表示される").toBeVisible();
  await expect(
    management.openCurrentPhaseButton(card),
    "現在段階を開くボタンが有効",
  ).toBeEnabled();
  await management.openCurrentPhase(card);
  await phase.waitForScreen();
  return phase;
}

test("E2E-UC-FIX-MODEL-001 未実行ジョブで AI 設定パネルが「設定未完了」を表示する", async ({
  page,
}) => {
  // presenter の空文字正規化（isExecutionConfigured が空文字を未設定と判別する）後、
  // 単語翻訳画面の AI モデル設定パネルが「設定未完了」を表示し、「固定済み」を表示しないことを証明する。

  // Arrange: AI 設定未完了ジョブの単語翻訳画面を開く。
  const phase = await openTermPhaseMissingAISettings(page);

  // Assert: 状態 pill が「設定未完了」を表示する。
  await expect(
    phase.aiSettingsStatusPill,
    "状態 pill が「設定未完了」を表示する",
  ).toContainText("設定未完了");

  // Assert: 「固定済み」pill が表示されない。
  await expect(
    phase.aiSettingsStatusPill,
    "状態 pill が「固定済み」を表示しない",
  ).not.toContainText("固定済み");
});

test("E2E-UC-FIX-MODEL-002 AI 設定未完了状態で開始ボタンが無効化され禁止理由が表示される", async ({
  page,
}) => {
  // presenter 正規化後に aiSettingsBlockedReason が正しい禁止理由を返し、
  // 開始ボタンが disabled を維持し、禁止理由領域に AI 設定不足を示すテキストが表示されることを証明する。

  // Arrange: AI 設定未完了ジョブの単語翻訳画面を開く。
  const phase = await openTermPhaseMissingAISettings(page);

  // Assert: 開始ボタンが disabled を維持する。
  await expect(
    phase.startButton,
    "開始ボタンが disabled を維持する",
  ).toBeDisabled();

  // Assert: 禁止理由領域に AI 設定不足を示すテキストが表示される。
  await expect(
    phase.startBlockedReason,
    "禁止理由領域に AI 設定不足を示すテキストが表示される",
  ).toContainText(/実行設定が未構成のため開始できません/);
});

test("E2E-UC-FIX-MODEL-003 AI 設定未完了から保存操作で「固定済み」へ遷移する", async ({
  page,
}) => {
  // AI 設定未完了状態からユーザーが AI サービス選択 → モデル一覧更新 → モデル選択の操作を行い、
  // saveAISettings が SaveTermTranslationPhaseAISettings を呼び出し、
  // GetTermTranslationPhaseSummary 再取得後に pill が「固定済み」になり、
  // 開始ボタンが有効化されることを証明する。

  // Arrange: AI 設定未完了ジョブの単語翻訳画面を開く。
  const phase = await openTermPhaseMissingAISettings(page);

  // 開始前の状態確認: 設定未完了であることを確認する。
  await expect(
    phase.aiSettingsStatusPill,
    "操作前は状態 pill が「設定未完了」を表示する",
  ).toContainText("設定未完了");

  // Act: AI サービスを選択する（Gemini を選択してモデル一覧更新を可能にする）。
  // ListProviderSettings が非同期で完了して provider options が更新されるまで待機する。
  // <option> 要素は select が閉じた状態では Playwright の visibility 判定が hidden になるため attached で待機する。
  await phase.aiProviderSelect.locator('option[value="gemini"]').waitFor({ state: "attached" });
  await phase.aiProviderSelect.selectOption("gemini");

  // Act: 「モデル一覧を更新」ボタンを押して ListTranslationJobSetupProviderModels を呼び出す。
  await phase.refreshAIModelButton.click();

  // Act: モデル select で gemini-test を選択する（handleModelChange → saveAISettings →
  // SaveTermTranslationPhaseAISettings → GetTermTranslationPhaseSummary 再取得）。
  // ListTranslationJobSetupProviderModels が完了して gemini-test option が現れるまで待機する。
  await phase.aiModelSelect.locator('option[value="gemini-test"]').waitFor({ state: "attached" });
  await phase.aiModelSelect.selectOption("gemini-test");

  // Assert: 状態 pill が「固定済み」になる。
  await expect(
    phase.aiSettingsStatusPill,
    "設定保存後に状態 pill が「固定済み」になる",
  ).toContainText("固定済み");

  // Assert: 開始ボタンが有効化される（禁止理由が解消する）。
  await expect(
    phase.startButton,
    "設定保存後に開始ボタンが有効化される",
  ).toBeEnabled();
});
