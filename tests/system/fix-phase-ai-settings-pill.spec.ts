import { expect, test, type Page } from "@playwright/test";

import {
  BodyTranslationPhasePage,
  PersonaGenerationPhasePage,
  TermTranslationPhasePage,
  TranslationJobManagementPage,
} from "./support/system-test-pages";
import { installScenarioWailsMocks } from "./support/scenario-wails-mocks";

// 単語翻訳: AI 設定未完了状態からモデルを選択して pill「固定済み」と開始ボタン有効化を証明するシナリオ（E2E-UC-056）。
test.describe("E2E-UC-056 term translation: AI 設定保存後に pill が「固定済み」になり開始ボタンが有効化される", () => {
  test.beforeEach(async ({ page }) => {
    await installScenarioWailsMocks(page, { termAISettingsMissing: true });
  });

  async function openTermPhaseMissingAI(page: Page): Promise<TermTranslationPhasePage> {
    // Arrange: 管理画面から単語翻訳 AI 設定未完了ジョブを開く。
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

  test("E2E-UC-056-A 単語翻訳: 未設定状態で pill が「設定未完了」を表示し開始ボタンが disabled になる", async ({
    page,
  }) => {
    // AI 設定未完了状態で、単語翻訳画面の pill が「設定未完了」を表示し、開始ボタンが disabled になることを証明する。

    // Arrange: AI 設定未完了の単語翻訳画面を開く。
    const phase = await openTermPhaseMissingAI(page);

    // Assert: 状態 pill が「設定未完了」を表示する。
    await expect(
      phase.aiSettingsStatusPill,
      "状態 pill が「設定未完了」を表示する",
    ).toContainText("設定未完了");

    // Assert: 開始ボタンが disabled になる。
    await expect(
      phase.startButton,
      "開始ボタンが disabled になる",
    ).toBeDisabled();
  });

  test("E2E-UC-056-B 単語翻訳: モデルを選択して AI 設定を保存すると pill が「固定済み」になり開始ボタンが有効化される", async ({
    page,
  }) => {
    // 未設定状態から AI サービスを選択し「モデル一覧を更新」してモデルを選択することで AI 設定を保存し、
    // pill が「固定済み」に切り替わり、開始ボタンが enabled になることを証明する。
    // これは SaveTermTranslationPhaseAISettings 成功 → GetTermTranslationPhaseSummary 再取得 →
    // isExecutionConfigured=true → pill 反映の経路を証明する。

    // Arrange: AI 設定未完了の単語翻訳画面を開く。
    const phase = await openTermPhaseMissingAI(page);

    // Arrange 前提確認: 操作前は「設定未完了」であることを確認する。
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
    // SaveTermTranslationPhaseAISettings（E2E-UC-056） → GetTermTranslationPhaseSummary 再取得）。
    // ListTranslationJobSetupProviderModels が完了して gemini-test option が現れるまで待機する。
    await phase.aiModelSelect.locator('option[value="gemini-test"]').waitFor({ state: "attached" });
    await phase.aiModelSelect.selectOption("gemini-test");

    // Assert: 状態 pill が「固定済み」になる。
    await expect(
      phase.aiSettingsStatusPill,
      "設定保存後に状態 pill が「固定済み」になる",
    ).toContainText("固定済み");

    // Assert: 開始ボタンが enabled になる（AI 設定保存後に実行設定が構成済みと判定される）。
    await expect(
      phase.startButton,
      "設定保存後に開始ボタンが enabled になる",
    ).toBeEnabled();
  });
});

// NPC ペルソナ生成: AI 設定未完了状態からモデルを選択して pill「固定済み」と開始ボタン有効化を証明するシナリオ（E2E-UC-057）。
test.describe("E2E-UC-057 persona generation: AI 設定保存後に pill が「固定済み」になり開始ボタンが有効化される", () => {
  test.beforeEach(async ({ page }) => {
    await installScenarioWailsMocks(page, { personaAISettingsMissing: true });
  });

  async function openPersonaPhaseMissingAI(page: Page): Promise<PersonaGenerationPhasePage> {
    // Arrange: 管理画面から NPC ペルソナ生成 AI 設定未完了ジョブを開く。
    const management = new TranslationJobManagementPage(page);
    const phase = new PersonaGenerationPhasePage(page);
    await management.open();
    const card = management.jobCard("system-test-persona-ai-missing");
    await expect(card, "NPC ペルソナ生成 AI 設定未完了ジョブのカードが表示される").toBeVisible();
    await expect(
      management.openCurrentPhaseButton(card),
      "現在段階を開くボタンが有効",
    ).toBeEnabled();
    await management.openCurrentPhase(card);
    await phase.waitForScreen();
    return phase;
  }

  test("E2E-UC-057-A NPC ペルソナ生成: 未設定状態で pill が「設定未完了」を表示する", async ({
    page,
  }) => {
    // AI 設定未完了状態で、NPC ペルソナ生成画面の pill が「設定未完了」を表示することを証明する。
    // NPC ペルソナ生成フェーズの presenter は execution 設定状態を canStart に反映しないため、
    // 開始ボタンの enabled/disabled 状態は E2E-UC-057-A の証明対象外とする。

    // Arrange: AI 設定未完了の NPC ペルソナ生成画面を開く。
    const phase = await openPersonaPhaseMissingAI(page);

    // Assert: 状態 pill が「設定未完了」を表示する。
    await expect(
      phase.aiSettingsStatusPill,
      "状態 pill が「設定未完了」を表示する",
    ).toContainText("設定未完了");
  });

  test("E2E-UC-057-B NPC ペルソナ生成: モデルを選択して AI 設定を保存すると pill が「固定済み」になり開始ボタンが有効化される", async ({
    page,
  }) => {
    // 未設定状態から AI サービスを選択し「モデル一覧を更新」してモデルを選択することで AI 設定を保存し、
    // NPC ペルソナ生成段階の pill が「固定済み」に切り替わり、開始ボタンが enabled になることを証明する。
    // 3 phase 共通の SaveXxxPhaseAISettings 成功 → summary 反映 → pill 切り替えの経路を証明する。

    // Arrange: AI 設定未完了の NPC ペルソナ生成画面を開く。
    const phase = await openPersonaPhaseMissingAI(page);

    // Arrange 前提確認: 操作前は「設定未完了」であることを確認する。
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
    // SavePersonaGenerationPhaseAISettings（E2E-UC-057） → GetPersonaGenerationPhaseSummary 再取得）。
    // ListTranslationJobSetupProviderModels が完了して gemini-test option が現れるまで待機する。
    await phase.aiModelSelect.locator('option[value="gemini-test"]').waitFor({ state: "attached" });
    await phase.aiModelSelect.selectOption("gemini-test");

    // Assert: 状態 pill が「固定済み」になる。
    await expect(
      phase.aiSettingsStatusPill,
      "設定保存後に状態 pill が「固定済み」になる",
    ).toContainText("固定済み");

    // Assert: 開始ボタンが enabled になる（AI 設定保存後に実行設定が構成済みと判定される）。
    await expect(
      phase.startButton,
      "設定保存後に開始ボタンが enabled になる",
    ).toBeEnabled();
  });
});

// 本文翻訳: AI 設定未完了状態からモデルを選択して pill「固定済み」と開始ボタン有効化を証明するシナリオ（E2E-UC-058）。
test.describe("E2E-UC-058 body translation: AI 設定保存後に pill が「固定済み」になり開始ボタンが有効化される", () => {
  test.beforeEach(async ({ page }) => {
    await installScenarioWailsMocks(page, { bodyAISettingsMissing: true });
  });

  async function openBodyPhaseMissingAI(page: Page): Promise<BodyTranslationPhasePage> {
    // Arrange: 管理画面から本文翻訳 AI 設定未完了ジョブを開く。
    const management = new TranslationJobManagementPage(page);
    const phase = new BodyTranslationPhasePage(page);
    await management.open();
    const card = management.jobCard("system-test-body-ai-missing");
    await expect(card, "本文翻訳 AI 設定未完了ジョブのカードが表示される").toBeVisible();
    await expect(
      management.openCurrentPhaseButton(card),
      "現在段階を開くボタンが有効",
    ).toBeEnabled();
    await management.openCurrentPhase(card);
    await phase.waitForScreen();
    return phase;
  }

  test("E2E-UC-058-A 本文翻訳: 未設定状態で pill が「設定未完了」を表示する", async ({
    page,
  }) => {
    // AI 設定未完了状態で、本文翻訳画面の pill が「設定未完了」を表示することを証明する。
    // 本文翻訳フェーズの presenter は execution 設定状態を canStart に反映しないため、
    // 開始ボタンの enabled/disabled 状態は E2E-UC-058-A の証明対象外とする。

    // Arrange: AI 設定未完了の本文翻訳画面を開く。
    const phase = await openBodyPhaseMissingAI(page);

    // Assert: 状態 pill が「設定未完了」を表示する。
    await expect(
      phase.aiSettingsStatusPill,
      "状態 pill が「設定未完了」を表示する",
    ).toContainText("設定未完了");
  });

  test("E2E-UC-058-B 本文翻訳: モデルを選択して AI 設定を保存すると pill が「固定済み」になり開始ボタンが有効化される", async ({
    page,
  }) => {
    // 未設定状態から AI サービスを選択し「モデル一覧を更新」してモデルを選択することで AI 設定を保存し、
    // 本文翻訳段階の pill が「固定済み」に切り替わり、開始ボタンが enabled になることを証明する。
    // 3 phase 共通の SaveXxxPhaseAISettings 成功 → summary 反映 → pill 切り替えの経路を証明する。

    // Arrange: AI 設定未完了の本文翻訳画面を開く。
    const phase = await openBodyPhaseMissingAI(page);

    // Arrange 前提確認: 操作前は「設定未完了」であることを確認する。
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
    // SaveBodyTranslationPhaseAISettings（E2E-UC-058） → GetBodyTranslationPhaseSummary 再取得）。
    // ListTranslationJobSetupProviderModels が完了して gemini-test option が現れるまで待機する。
    await phase.aiModelSelect.locator('option[value="gemini-test"]').waitFor({ state: "attached" });
    await phase.aiModelSelect.selectOption("gemini-test");

    // Assert: 状態 pill が「固定済み」になる。
    await expect(
      phase.aiSettingsStatusPill,
      "設定保存後に状態 pill が「固定済み」になる",
    ).toContainText("固定済み");

    // Assert: 開始ボタンが enabled になる（AI 設定保存後に実行設定が構成済みと判定される）。
    await expect(
      phase.startButton,
      "設定保存後に開始ボタンが enabled になる",
    ).toBeEnabled();
  });
});
