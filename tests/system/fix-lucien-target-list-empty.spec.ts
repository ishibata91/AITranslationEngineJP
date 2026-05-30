import { expect, test, type Page } from "@playwright/test"

import {
  TermTranslationPhasePage,
  TranslationJobManagementPage
} from "./support/system-test-pages"
import { installScenarioWailsMocks } from "./support/scenario-wails-mocks"

/**
 * fix-lucien-target-list-empty シナリオテスト
 *
 * 確定原因: 単語翻訳段階の初回ロードで取得競合検出の連番不一致により
 * 処理対象取得結果が破棄され、初回表示が0件になる。
 *
 * このファイルは fail-test ベースの修正前テスト追加である。
 * E2E-LTLE-001 と E2E-LTLE-002 は、未修正のプロダクトコードでは
 * 初回0件を検出して期待どおり失敗する。
 * E2E-LTLE-003 は修正前後で通る想定（真の0件は空状態が正しい）。
 */

// jobId=14 を母数0ジョブ専用として予約する（既存の seededPhaseJobs と衝突しない）
const ZERO_AI_TARGET_JOB_ID = 14

async function openTermPhaseForJob(
  page: Page,
  jobText: string
): Promise<TermTranslationPhasePage> {
  const management = new TranslationJobManagementPage(page)
  const phase = new TermTranslationPhasePage(page)

  await management.open()
  const card = management.jobCard(jobText)
  await expect(card, `ジョブカードが表示されること: ${jobText}`).toBeVisible()
  await expect(
    management.openCurrentPhaseButton(card),
    `現在の翻訳段階へ進むボタンが有効なこと: ${jobText}`
  ).toBeEnabled()
  await management.openCurrentPhase(card)
  await phase.waitForScreen()

  return phase
}

test("E2E-LTLE-001 term phase initial display shows rows when ai target count is non-zero", async ({
  page
}) => {
  // 正常: 進捗パネルの AI 翻訳対象語母数が1以上のとき、単語翻訳段階画面の
  // 初回表示（検索・ページ操作なし）で処理対象行が1件以上表示され、
  // 空状態を表示しないことを証明する。
  // 初回ロードで取得結果が破棄される場合、この検証は失敗する（期待どおりの失敗）。
  await installScenarioWailsMocks(page)

  // Arrange: 翻訳管理画面から母数3の単語翻訳段階画面へ進む
  const phase = await openTermPhaseForJob(page, "system-test-term")

  // Act: 初回表示を待つ（検索・ページ操作は行わない）
  await expect(phase.screen).toBeVisible()
  await expect(
    phase.processingTargetListRegion,
    "処理対象一覧領域が表示されること"
  ).toBeVisible()

  // Assert: 処理対象行が1件以上あり、空状態が表示されないこと
  await expect(
    phase.processingTargetRows,
    "初回表示で処理対象行が1件以上表示されること"
  ).not.toHaveCount(0)
  await expect(
    page.getByTestId("term-translation-phase-processing-target-empty"),
    "初回表示で空状態(処理対象がありません)が表示されないこと"
  ).toHaveCount(0)
})

test("E2E-LTLE-002 term phase shows rows after search and reload", async ({
  page
}) => {
  // 境界: 母数1以上 + 一覧表示済みの状態から検索を行い、その後画面をリロードする。
  // リロード後の初回表示で再び0件に戻らず、処理対象行が1件以上表示されることを証明する。
  // 初回ロードで取得結果が破棄される場合、リロード後にも初回0件が再現し失敗する（期待どおりの失敗）。
  await installScenarioWailsMocks(page)

  // Arrange: 単語翻訳段階画面を初回表示し、処理対象行があることを前提に検索を行う
  const phase = await openTermPhaseForJob(page, "system-test-term")
  await expect(phase.processingTargetListRegion).toBeVisible()
  // 前提: 初回表示で行が存在する（E2E-LTLE-001 と同じ前提）
  await expect(
    phase.processingTargetRows,
    "検索前の前提: 処理対象行が1件以上表示されること"
  ).not.toHaveCount(0)

  // Act: 検索語を入力して絞り込み後、ブラウザをリロードして同じ段階画面に再到達する
  // （ハッシュルーティングのため、リロード後は翻訳管理一覧に戻る。
  //   再度同じジョブの段階画面を開くことがリロード後の「初回表示」に相当する）
  await phase.searchProcessingTargets("Dragonborn")
  await expect(
    phase.processingTargetRows,
    "検索後に絞り込み結果が表示されること"
  ).toHaveCount(1)
  await page.reload()
  const phaseAfterReload = await openTermPhaseForJob(page, "system-test-term")
  await expect(phaseAfterReload.processingTargetListRegion).toBeVisible()

  // Assert: リロード後の初回表示で処理対象行が1件以上あり、空状態が表示されないこと
  await expect(
    phaseAfterReload.processingTargetRows,
    "リロード後の初回表示で処理対象行が1件以上表示されること"
  ).not.toHaveCount(0)
  await expect(
    page.getByTestId("term-translation-phase-processing-target-empty"),
    "リロード後の初回表示で空状態(処理対象がありません)が表示されないこと"
  ).toHaveCount(0)
})

test("E2E-LTLE-003 term phase initial display shows empty state when ai target count is zero", async ({
  page
}) => {
  // 境界: 進捗パネルの AI 翻訳対象語母数が0のとき、単語翻訳段階画面の
  // 初回表示で空状態が表示され、処理対象行を表示しないことを証明する。
  // 是正が真の0件シナリオで空状態仕様を壊さないことを確認する。
  await installScenarioWailsMocks(page, {
    termZeroAITargetJobId: ZERO_AI_TARGET_JOB_ID
  })

  // Arrange: 翻訳管理画面から母数0の単語翻訳段階画面へ進む
  const phase = await openTermPhaseForJob(
    page,
    "system-test-term-zero-ai-target"
  )

  // Act: 初回表示を待つ（検索・ページ操作は行わない）
  await expect(phase.screen).toBeVisible()
  await expect(
    phase.processingTargetListRegion,
    "処理対象一覧領域が表示されること"
  ).toBeVisible()

  // Assert: 空状態が表示され、処理対象行がないこと
  await expect(
    page.getByTestId("term-translation-phase-processing-target-empty"),
    "母数0のとき初回表示で空状態(処理対象がありません)が表示されること"
  ).toBeVisible()
  await expect(
    phase.processingTargetRows,
    "母数0のとき初回表示で処理対象行が表示されないこと"
  ).toHaveCount(0)
})
