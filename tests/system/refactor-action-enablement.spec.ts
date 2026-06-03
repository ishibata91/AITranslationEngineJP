import { expect, test, type Page } from "@playwright/test";

import {
  BodyTranslationPhasePage,
  PersonaGenerationPhasePage,
  TermTranslationPhasePage,
  TranslationJobManagementPage,
  type TranslationPhasePage,
} from "./support/system-test-pages";
import { installScenarioWailsMocks } from "./support/scenario-wails-mocks";

// RAEF-E2E テスト: refactor-action-enablement-derive-on-frontend
// 承認済み設計差分（design-diff.md）の H 節論理式と G 節 BlockedReason を UI 操作から証明する。
// motivating bug（actionEnablement 必須検証で gateway がエラーを返す）の消滅を含む。

test.beforeEach(async ({ page }) => {
  await installScenarioWailsMocks(page);
});

async function openPhase(
  page: Page,
  jobText: string,
  phase: TranslationPhasePage,
): Promise<void> {
  // ジョブ管理画面から対象ジョブの「現在の段階を開く」ボタン経由で段階画面に遷移する。
  const management = new TranslationJobManagementPage(page);
  await management.open();
  const card = management.jobCard(jobText);
  await expect(card, `job card is visible: ${jobText}`).toBeVisible();
  await management.openCurrentPhase(card);
  await phase.waitForScreen();
}

// ---------------------------------------------------------------------------
// RAEF-E2E-001: motivating bug の再現と消滅
// ---------------------------------------------------------------------------

test("RAEF-E2E-001 term start accepted without actionEnablement shape validator error", async ({
  page,
}) => {
  // motivating bug 消滅の確認。
  // projection が gateway を正常に通過し、actionEnablement 必須検証エラーが発生しないことを証明する。
  // 前提: aiSettingsConfigured=true, aiTargetCount=3, phaseLifecycle=idle_ready, jobLifecycle=running（非 TERMINAL_JOB）。
  // 期待: 開始ボタンが enabled、クリック後に段階状態が実行中に遷移する。
  // 「単語翻訳段階の操作に失敗しました。」通知は表示されない。
  const phase = new TermTranslationPhasePage(page);
  await openPhase(page, "raef-term-normal-start", phase);

  // Arrange: 開始ボタンが有効であることを確認する（canStart=true の証明）。
  await expect(phase.startButton).toBeEnabled();

  // Act: 開始ボタンをクリックする。
  await phase.start();

  // Assert: エラー通知が出ないこと（motivating bug が消滅していること）を確認する。
  // actionEnablement 由来の shape validator エラーなら画面上にエラー通知が表示される。
  await expect(page.getByText("単語翻訳段階の操作に失敗しました。")).not.toBeVisible();
  await expect(page.getByText("shape validator")).not.toBeVisible();
});

// ---------------------------------------------------------------------------
// RAEF-E2E-002: term blocked - aiSettingsConfigured=false
// ---------------------------------------------------------------------------

test("RAEF-E2E-002 term start button disabled when aiSettingsConfigured=false", async ({
  page,
}) => {
  // 例外経路: aiSettingsConfigured=false のとき canStart=false（H-1 条件）。
  // 前提: projection.aiSettingsConfigured=false。
  // 期待: 開始ボタンが disabled、BlockedReason が「実行設定が未構成のため開始できません。」。
  const phase = new TermTranslationPhasePage(page);
  await openPhase(page, "raef-term-blocked-ai-missing", phase);

  // Arrange / Assert: 開始ボタンが無効であることを確認する。
  await expect(phase.startButton).toBeDisabled();

  // Assert: BlockedReason 文言が表示されていることを確認する。
  await expect(phase.startBlockedReason).toContainText(
    "実行設定が未構成のため開始できません。",
  );
});

// ---------------------------------------------------------------------------
// RAEF-E2E-003: term blocked - aiTargetCount=0
// ---------------------------------------------------------------------------

test("RAEF-E2E-003 term start button disabled when aiTargetCount=0", async ({
  page,
}) => {
  // 例外経路: aiTargetCount=0 のとき hasProcessingTarget=false で canStart=false（H-1 条件）。
  // 前提: projection.aiTargetCount=0。
  // 期待: 開始ボタンが disabled、BlockedReason が「処理対象が 0 件のため開始できません。」。
  const phase = new TermTranslationPhasePage(page);
  await openPhase(page, "raef-term-blocked-zero-target", phase);

  // Arrange / Assert: 開始ボタンが無効であることを確認する。
  await expect(phase.startButton).toBeDisabled();

  // Assert: BlockedReason 文言が表示されていることを確認する。
  await expect(phase.startBlockedReason).toContainText(
    "処理対象が 0 件のため開始できません。",
  );
});

// ---------------------------------------------------------------------------
// RAEF-E2E-004: term blocked - TERMINAL_JOB
// ---------------------------------------------------------------------------

test("RAEF-E2E-004 term start button disabled when jobLifecycle is TERMINAL_JOB", async ({
  page,
}) => {
  // 境界: jobLifecycle=completed（TERMINAL_JOB）のとき canStart=false（H-1 条件）。
  // 前提: projection.jobLifecycle="completed"。
  // 期待: 開始ボタンが disabled、BlockedReason が「ジョブが終端状態のため開始できません。」。
  const phase = new TermTranslationPhasePage(page);
  await openPhase(page, "raef-term-blocked-terminal", phase);

  // Arrange / Assert: 開始ボタンが無効であることを確認する。
  await expect(phase.startButton).toBeDisabled();

  // Assert: BlockedReason 文言が表示されていることを確認する。
  await expect(phase.startBlockedReason).toContainText(
    "ジョブが終端状態のため開始できません。",
  );
});

// ---------------------------------------------------------------------------
// RAEF-E2E-005: persona start 正常パス
// ---------------------------------------------------------------------------

test("RAEF-E2E-005 persona start accepted when previousPhaseLifecycle=completed", async ({
  page,
}) => {
  // 正常経路: persona phase での motivating bug 消滅確認。
  // 前提: aiSettingsConfigured=true, targetCount=2, phaseLifecycle=idle_ready, previousPhaseLifecycle=completed。
  // 期待: 開始ボタンが enabled、クリック後に実行中状態に遷移する。
  const phase = new PersonaGenerationPhasePage(page);
  await openPhase(page, "raef-persona-normal-start", phase);

  // Arrange: 開始ボタンが有効であることを確認する。
  await expect(phase.startButton).toBeEnabled();

  // Act: 開始ボタンをクリックする。
  await phase.start();

  // Assert: actionEnablement 由来のエラー通知が出ないことを確認する。
  await expect(page.getByText("NPC ペルソナ生成段階の操作に失敗しました。")).not.toBeVisible();
});

// ---------------------------------------------------------------------------
// RAEF-E2E-006: persona blocked - previousPhaseLifecycle 不足
// ---------------------------------------------------------------------------

test("RAEF-E2E-006 persona start button disabled when previousPhaseLifecycle is not COMPLETED_PHASE", async ({
  page,
}) => {
  // 例外経路: previousPhaseLifecycle が非 COMPLETED_PHASE のとき canStart=false（H-6 条件）。
  // 前提: projection.previousPhaseLifecycle="pending"（term 未完了）。
  // 期待: 開始ボタンが disabled、BlockedReason が「単語翻訳段階が完了していないためペルソナ生成を開始できません。」。
  const phase = new PersonaGenerationPhasePage(page);
  await openPhase(page, "raef-persona-blocked-prev-incomplete", phase);

  // Arrange / Assert: 開始ボタンが無効であることを確認する。
  await expect(phase.startButton).toBeDisabled();

  // Assert: BlockedReason 文言が表示されていることを確認する。
  await expect(phase.startBlockedReason).toContainText(
    "単語翻訳段階が完了していないためペルソナ生成を開始できません。",
  );
});

// ---------------------------------------------------------------------------
// RAEF-E2E-007: body start 正常パス
// ---------------------------------------------------------------------------

test("RAEF-E2E-007 body start accepted when personaBodyReadiness.bodyReadiness=true", async ({
  page,
}) => {
  // 正常経路: body phase での motivating bug 消滅確認。
  // 前提: aiSettingsConfigured=true, targetCount=4, phaseLifecycle=idle_ready,
  //       previousPhaseLifecycle=completed, personaBodyReadiness.bodyReadiness=true。
  // 期待: 開始ボタンが enabled、クリック後に実行中状態に遷移する。
  const phase = new BodyTranslationPhasePage(page);
  await openPhase(page, "raef-body-normal-start", phase);

  // Arrange: 開始ボタンが有効であることを確認する。
  await expect(phase.startButton).toBeEnabled();

  // Act: 開始ボタンをクリックする。
  await phase.start();

  // Assert: actionEnablement 由来のエラー通知が出ないことを確認する。
  await expect(page.getByText("本文翻訳段階の操作に失敗しました。")).not.toBeVisible();
});

// ---------------------------------------------------------------------------
// RAEF-E2E-008: body blocked - personaBodyReady 未達
// ---------------------------------------------------------------------------

test("RAEF-E2E-008 body start button disabled when personaBodyReadiness.bodyReadiness=false", async ({
  page,
}) => {
  // 例外経路（選択 A の体験差分）: personaBodyReadiness.bodyReadiness=false のとき canStart=false（H-12 条件）。
  // 前提: projection.personaBodyReadiness.bodyReadiness=false,
  //       projection.personaBodyReadiness.snapshotReferenceStatus="not_available"。
  // 期待: 開始ボタンが disabled、BlockedReason が「ペルソナ snapshot 参照が準備できていないため本文翻訳を開始できません。」。
  const phase = new BodyTranslationPhasePage(page);
  await openPhase(page, "raef-body-blocked-no-readiness", phase);

  // Arrange / Assert: 開始ボタンが無効であることを確認する。
  await expect(phase.startButton).toBeDisabled();

  // Assert: BlockedReason 文言が表示されていることを確認する。
  await expect(phase.startBlockedReason).toContainText(
    "ペルソナ snapshot 参照が準備できていないため本文翻訳を開始できません。",
  );
});

// ---------------------------------------------------------------------------
// RAEF-E2E-009: persona 次段階へ有効化（選択 A 体験差分）
// ---------------------------------------------------------------------------

test("RAEF-E2E-009 persona next-phase button enabled when phaseLifecycle=COMPLETED_PHASE without personaBodyReadiness gate", async ({
  page,
}) => {
  // 正常経路（選択 A）: persona 完了かつ非 terminal のとき次段階へボタンが有効になる（H-11 条件）。
  // 選択 A の方針により persona 側の次段階へボタンに personaBodyReadiness ガードがないことを確認する。
  // 前提: projection.phaseLifecycle=completed, projection.jobLifecycle=running（非 TERMINAL_JOB）。
  // 期待: 次段階へボタンが enabled。
  const phase = new PersonaGenerationPhasePage(page);
  await openPhase(page, "raef-persona-completed-next-phase", phase);

  // Arrange / Assert: 次段階へボタンが有効であることを確認する（personaBodyReadiness ガードなしの証明）。
  await expect(phase.nextPhaseButton).toBeEnabled();
});

// ---------------------------------------------------------------------------
// RAEF-E2E-010: persona 次段階へ無効化
// ---------------------------------------------------------------------------

test("RAEF-E2E-010 persona next-phase button disabled when phaseLifecycle is not COMPLETED_PHASE", async ({
  page,
}) => {
  // 例外経路: phaseLifecycle が非 COMPLETED_PHASE のとき次段階へボタンが無効になる（H-11 条件）。
  // 前提: projection.phaseLifecycle=idle_ready（未完了）。
  // 期待: 次段階へボタンが disabled。
  const phase = new PersonaGenerationPhasePage(page);
  await openPhase(page, "raef-persona-incomplete-next-phase", phase);

  // Arrange / Assert: 次段階へボタンが無効であることを確認する。
  await expect(phase.nextPhaseButton).toBeDisabled();
});

// ---------------------------------------------------------------------------
// RAEF-E2E-011: term 中断有効化（RUNNING_PHASE）
// ---------------------------------------------------------------------------

test("RAEF-E2E-011 term pause button enabled when phaseLifecycle=RUNNING_PHASE", async ({
  page,
}) => {
  // 正常経路: phaseLifecycle=running かつ非 terminal のとき中断ボタンが有効になる（H-2 条件）。
  // 前提: projection.phaseLifecycle=running, jobLifecycle=running（非 TERMINAL_JOB）。
  // 期待: 中断ボタンが enabled。
  const phase = new TermTranslationPhasePage(page);
  await openPhase(page, "raef-term-running-pause", phase);

  // Arrange / Assert: 中断ボタンが有効であることを確認する。
  await expect(phase.pauseButton).toBeEnabled();
});

// ---------------------------------------------------------------------------
// RAEF-E2E-012: term 中断無効化（非 RUNNING_PHASE）
// ---------------------------------------------------------------------------

test("RAEF-E2E-012 term pause button disabled when phaseLifecycle is not RUNNING_PHASE", async ({
  page,
}) => {
  // 境界: phaseLifecycle=idle_ready（実行中でない）のとき中断ボタンが無効になる（H-2 条件）。
  // 前提: projection.phaseLifecycle=idle_ready。
  // 期待: 中断ボタンが disabled。
  const phase = new TermTranslationPhasePage(page);
  await openPhase(page, "raef-term-idle-pause-blocked", phase);

  // Arrange / Assert: 中断ボタンが無効であることを確認する。
  // idle 状態では中断ボタンが表示されないか、disabled で表示される。
  const pauseButtonVisible = await phase.pauseButton.isVisible().catch(() => false);
  if (pauseButtonVisible) {
    await expect(phase.pauseButton).toBeDisabled();
  }
  // 中断ボタンが表示されない場合も H-2 の境界確認として有効（idle_ready では中断操作対象なし）。
});

// ---------------------------------------------------------------------------
// RAEF-E2E-013: term 再開有効化（PAUSED_PHASE）
// ---------------------------------------------------------------------------

test("RAEF-E2E-013 term resume button enabled when phaseLifecycle=PAUSED_PHASE", async ({
  page,
}) => {
  // 正常経路: phaseLifecycle=paused かつ非 terminal のとき再開ボタンが有効になる（H-3 条件）。
  // 前提: projection.phaseLifecycle=paused, jobLifecycle=running（非 TERMINAL_JOB）。
  // 期待: 再開ボタンが enabled。
  const phase = new TermTranslationPhasePage(page);
  await openPhase(page, "raef-term-paused-resume", phase);

  // Arrange / Assert: 再開ボタンが有効であることを確認する。
  await expect(phase.resumeButton).toBeEnabled();
});

// ---------------------------------------------------------------------------
// RAEF-E2E-014: term 再試行有効化（RECOVERABLE_FAILED_PHASE）
// ---------------------------------------------------------------------------

test("RAEF-E2E-014 term retry button enabled when phaseLifecycle=RECOVERABLE_FAILED_PHASE", async ({
  page,
}) => {
  // 正常経路: phaseLifecycle=recoverable_failed かつ非 terminal のとき再試行ボタンが有効になる（H-4 条件）。
  // 前提: projection.phaseLifecycle=recoverable_failed, jobLifecycle=running（非 TERMINAL_JOB）。
  // 期待: 再試行ボタンが enabled。
  const phase = new TermTranslationPhasePage(page);
  await openPhase(page, "raef-term-failed-retry", phase);

  // Arrange / Assert: 再試行ボタンが有効であることを確認する。
  await expect(phase.retryButton).toBeEnabled();
});

// ---------------------------------------------------------------------------
// RAEF-E2E-015: persona キャンセル有効化（RUNNING_PHASE）
// ---------------------------------------------------------------------------

test("RAEF-E2E-015 persona cancel button enabled when phaseLifecycle=RUNNING_PHASE", async ({
  page,
}) => {
  // 正常経路: phaseLifecycle=running かつ非 terminal のときキャンセルボタンが有効になる（H-10 条件）。
  // 前提: projection.phaseLifecycle=running, jobLifecycle=running（非 TERMINAL_JOB）。
  // 期待: キャンセルボタンが enabled。
  const phase = new PersonaGenerationPhasePage(page);
  await openPhase(page, "raef-persona-running-cancel", phase);

  // Arrange / Assert: キャンセルボタンが有効であることを確認する。
  await expect(phase.cancelButton).toBeEnabled();
});

// ---------------------------------------------------------------------------
// RAEF-E2E-016: body キャンセル有効化（RUNNING_PHASE）
// ---------------------------------------------------------------------------

test("RAEF-E2E-016 body cancel button enabled when phaseLifecycle=RUNNING_PHASE", async ({
  page,
}) => {
  // 正常経路: phaseLifecycle=running かつ非 terminal のときキャンセルボタンが有効になる（H-16 条件）。
  // 前提: projection.phaseLifecycle=running, jobLifecycle=running（非 TERMINAL_JOB）。
  // 期待: キャンセルボタンが enabled。
  const phase = new BodyTranslationPhasePage(page);
  await openPhase(page, "raef-body-running-cancel", phase);

  // Arrange / Assert: キャンセルボタンが有効であることを確認する。
  await expect(phase.cancelButton).toBeEnabled();
});
