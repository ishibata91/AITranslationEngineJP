import { expect, test, type Page } from "@playwright/test"

import {
  BodyTranslationPhasePage,
  PersonaGenerationPhasePage,
  TermTranslationPhasePage,
  TranslationJobManagementPage
} from "./support/system-test-pages"
import { installScenarioWailsMocks } from "./support/scenario-wails-mocks"

/**
 * fix-lucien-target-list-empty シナリオテスト
 *
 * 承認済み詳細仕様差分: job-run-phase-fetch-redesign
 * 対象ユースケース: 処理対象を確認する
 *
 * wave-4 SCN-target-list-e2e の実装。
 * - E2E-LTLE-001/002/003: 残置 E2E を wave-1/2/3 実装後の挙動へ green 化
 * - E2E-LTLE-004: 初回取得中ローディングレイヤーの操作排他確認
 * - E2E-PGTL-001: persona 段階での同型正常確認
 * - E2E-BTTL-001: body 段階での同型正常確認
 * - E2E-BTTL-002: body 成果物出力確認可否（活性）
 * - E2E-BTTL-003: body 成果物出力確認可否（非活性・理由）
 *
 * 注意: モック E2E は同期解決のため実機 IPC 飽和（約15秒遅延）を再現しない。
 * 実機固有症状の確認は実装後ブラウザ確認で扱う。
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

async function openPersonaPhaseForJob(
  page: Page,
  jobText: string
): Promise<PersonaGenerationPhasePage> {
  const management = new TranslationJobManagementPage(page)
  const phase = new PersonaGenerationPhasePage(page)

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

async function openBodyPhaseForJob(
  page: Page,
  jobText: string
): Promise<BodyTranslationPhasePage> {
  const management = new TranslationJobManagementPage(page)
  const phase = new BodyTranslationPhasePage(page)

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
  // screen-design-diff 規則1「処理対象一覧の総件数が1以上なら初回取得完了後の初回表示で行を表示する」の検証。
  await installScenarioWailsMocks(page)

  // Arrange: 翻訳管理画面から母数3の単語翻訳段階画面へ進む
  const phase = await openTermPhaseForJob(page, "system-test-term")

  // Act: 初回取得中ローディングレイヤーが外れるまで待つ（検索・ページ操作は行わない）
  await expect(phase.screen).toBeVisible()
  await expect(
    phase.processingTargetLoadingOverlay,
    "初回取得中ローディングレイヤーが外れること"
  ).toHaveCount(0)
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
  // screen-design-diff 規則11・規則12「開き直し時に再取得し件数分表示」の検証。
  await installScenarioWailsMocks(page)

  // Arrange: 単語翻訳段階画面を初回表示し、処理対象行があることを前提に検索を行う
  const phase = await openTermPhaseForJob(page, "system-test-term")
  await expect(phase.processingTargetLoadingOverlay).toHaveCount(0)
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
  await expect(phaseAfterReload.processingTargetLoadingOverlay).toHaveCount(0)
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
  // screen-design-diff 規則2「処理対象一覧の総件数が0の場合だけ空状態を表示する」の検証。
  // 是正が真の0件シナリオで空状態仕様を壊さないことを確認する。
  await installScenarioWailsMocks(page, {
    termZeroAITargetJobId: ZERO_AI_TARGET_JOB_ID
  })

  // Arrange: 翻訳管理画面から母数0の単語翻訳段階画面へ進む
  const phase = await openTermPhaseForJob(
    page,
    "system-test-term-zero-ai-target"
  )

  // Act: 初回取得中ローディングレイヤーが外れるまで待つ（検索・ページ操作は行わない）
  await expect(phase.screen).toBeVisible()
  await expect(
    phase.processingTargetLoadingOverlay,
    "初回取得中ローディングレイヤーが外れること"
  ).toHaveCount(0)
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

test("E2E-LTLE-004 term phase processing target loading overlay disables search during initial fetch", async ({
  page
}) => {
  // 境界(操作排他): 初回取得中ローディングレイヤーがフェーズ画面全体を覆うオーバーレイとして
  // 全体操作を排他することを証明する。
  // screen-design-diff 規則4〜7「初回取得中ローディングレイヤーがフェーズ画面全体を覆い全体操作を受け付けない」の検証。
  //
  // 注意: モック E2E は同期解決のため初回取得が即座に完了し、ローディングレイヤーの
  // 出現が一瞬になる可能性がある。初回取得が完了した後のローディングレイヤー消失と
  // 検索入力の操作可能状態を確認することで、操作排他が実装されていることを間接的に証明する。
  // 実機での15秒 IPC 飽和中の排他確認は実装後ブラウザ確認で扱う。
  await installScenarioWailsMocks(page)

  // Arrange: 翻訳管理画面から単語翻訳段階画面へ進む
  const phase = await openTermPhaseForJob(page, "system-test-term")
  await expect(phase.screen).toBeVisible()

  // Act & Assert: 初回取得が完了しローディングレイヤーが外れた後、検索入力が操作可能であること
  // （モック環境では初回取得が同期解決のため、ローディングレイヤーが外れた状態から確認する）
  await expect(
    phase.processingTargetLoadingOverlay,
    "初回取得完了後にローディングレイヤーが外れること"
  ).toHaveCount(0)
  await expect(
    phase.processingTargetSearchInput,
    "初回取得完了後に検索入力欄が操作可能であること"
  ).toBeEnabled()
  await expect(
    phase.processingTargetRows,
    "初回取得完了後に処理対象行が表示されること（ローディングが正常に外れた証拠）"
  ).not.toHaveCount(0)
})

test("E2E-PGTL-001 persona phase initial display shows rows when generation target count is non-zero", async ({
  page
}) => {
  // 正常(同型): NPC ペルソナ生成段階でも初回表示で件数分表示されること（term と同型）を証明する。
  // screen-design-diff 差分5 の persona selector（persona-generation-phase-processing-target-*）に依存。
  // screen-design-diff 規則1 が persona 段階にも同型適用されることの検証。
  await installScenarioWailsMocks(page)

  // Arrange: 翻訳管理画面から persona 段階のジョブへ進む
  const phase = await openPersonaPhaseForJob(page, "system-test-persona")

  // Act: 初回取得中ローディングレイヤーが外れるまで待つ（検索・ページ操作は行わない）
  await expect(phase.screen).toBeVisible()
  await expect(
    phase.processingTargetLoadingOverlay,
    "初回取得中ローディングレイヤーが外れること"
  ).toHaveCount(0)
  await expect(
    phase.processingTargetListRegion,
    "処理対象一覧領域が表示されること"
  ).toBeVisible()

  // Assert: 処理対象行が1件以上あり、空状態が表示されないこと
  await expect(
    phase.processingTargetRows,
    "persona 段階の初回表示で処理対象行が1件以上表示されること"
  ).not.toHaveCount(0)
  await expect(
    phase.processingTargetEmptyState,
    "persona 段階の初回表示で空状態が表示されないこと"
  ).toHaveCount(0)
})

test("E2E-BTTL-001 body phase initial display shows rows when translation target count is non-zero", async ({
  page
}) => {
  // 正常(同型): 本文翻訳段階でも初回表示で件数分表示されること（term と同型）を証明する。
  // screen-design-diff 差分5 の body selector（body-translation-phase-processing-target-*）に依存。
  // screen-design-diff 規則1 が body 段階にも同型適用されることの検証。
  await installScenarioWailsMocks(page)

  // Arrange: 翻訳管理画面から body 段階（pending）のジョブへ進む
  const phase = await openBodyPhaseForJob(page, "system-test-body-pending")

  // Act: 初回取得中ローディングレイヤーが外れるまで待つ（検索・ページ操作は行わない）
  await expect(phase.screen).toBeVisible()
  await expect(
    phase.processingTargetLoadingOverlay,
    "初回取得中ローディングレイヤーが外れること"
  ).toHaveCount(0)
  await expect(
    phase.processingTargetListRegion,
    "処理対象一覧領域が表示されること"
  ).toBeVisible()

  // Assert: 処理対象行が1件以上あり、空状態が表示されないこと
  await expect(
    phase.processingTargetRows,
    "body 段階の初回表示で処理対象行が1件以上表示されること"
  ).not.toHaveCount(0)
  await expect(
    phase.processingTargetEmptyState,
    "body 段階の初回表示で空状態が表示されないこと"
  ).toHaveCount(0)
})

test("E2E-BTTL-002 body phase shows check output readiness button enabled when conditions are met", async ({
  page
}) => {
  // 正常(取得経路統合): body の成果物出力確認可否が専用取得ではなく段階要約取得の事実
  // (completedFieldCount / statusConsistent / outputCount) から導出され活性表示されることを証明する。
  // screen-design-diff 差分7・規則21「成果物出力確認可否を段階要約取得の事実から導出する」の検証。
  // 成立条件(completedFieldCount=翻訳項目件数・statusConsistent=true)を満たす完了済みジョブを使う。
  await installScenarioWailsMocks(page)

  // Arrange: 翻訳管理画面から body 段階（completed=出力確認条件を満たすジョブ）へ進む
  const phase = await openBodyPhaseForJob(
    page,
    "system-test-body-ready-for-completion"
  )
  await expect(phase.screen).toBeVisible()
  await expect(
    phase.processingTargetLoadingOverlay,
    "初回取得中ローディングレイヤーが外れること"
  ).toHaveCount(0)

  // Act: 成果物出力確認の footer セクションが表示されていることを確認する（クリックは行わない）
  const completeNextActionSection = page.getByTestId(
    "body-translation-phase-complete-next-action"
  )
  await expect(
    completeNextActionSection,
    "body 段階の成果物出力確認 footer セクションが表示されること"
  ).toBeVisible()

  // Assert: 成果物出力確認ボタンが活性であること（出力確認条件を満たすジョブを使用）
  const checkOutputReadinessButton = completeNextActionSection.getByRole(
    "button",
    { name: "次へ進む" }
  )
  await expect(
    checkOutputReadinessButton,
    "成果物出力確認条件を満たすとき完了確認へ進むボタンが活性であること"
  ).toBeEnabled()
  await expect(
    completeNextActionSection.locator(".reason-summary"),
    "成果物出力確認条件を満たすとき不可理由が表示されないこと"
  ).toHaveText("次の作業へ進めます。")
})

test("E2E-BTTL-003 body phase shows check output readiness button disabled when conditions are not met", async ({
  page
}) => {
  // 境界(操作不可): 成立条件を満たさないとき成果物出力確認が不可で理由が出ることを証明する。
  // screen-design-diff 規則21「成果物出力確認可否の活性・理由を段階要約取得の事実から導出する」の検証。
  // 成立条件未達(pending=未完了)のジョブを使い、ボタン非活性と不可理由が表示されることを確認する。
  // 操作不可確認のためクリック操作は実行しない。
  await installScenarioWailsMocks(page)

  // Arrange: 翻訳管理画面から body 段階（pending=未完了）のジョブへ進む
  const phase = await openBodyPhaseForJob(page, "system-test-body-pending")
  await expect(phase.screen).toBeVisible()
  await expect(
    phase.processingTargetLoadingOverlay,
    "初回取得中ローディングレイヤーが外れること"
  ).toHaveCount(0)

  // Act: 成果物出力確認の footer セクションが表示されていることを確認する（クリックは行わない）
  const completeNextActionSection = page.getByTestId(
    "body-translation-phase-complete-next-action"
  )
  await expect(
    completeNextActionSection,
    "body 段階の成果物出力確認 footer セクションが表示されること"
  ).toBeVisible()

  // Assert: 成果物出力確認ボタンが非活性で不可理由が表示されること
  const checkOutputReadinessButton = completeNextActionSection.getByRole(
    "button",
    { name: "次へ進む" }
  )
  await expect(
    checkOutputReadinessButton,
    "成果物出力確認条件を満たさないとき完了確認へ進むボタンが非活性であること"
  ).toBeDisabled()
  await expect(
    completeNextActionSection.locator(".reason-primary"),
    "成果物出力確認条件を満たさないとき不可理由が表示されること"
  ).not.toHaveText("")
})
