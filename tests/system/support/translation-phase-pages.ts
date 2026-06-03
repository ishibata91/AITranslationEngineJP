import type { Locator, Page } from "@playwright/test"

import { SystemTestPageObject } from "./system-test-page-object"

type PhasePrefix =
  | "term-translation-phase"
  | "persona-generation-phase"
  | "body-translation-phase"

export class TranslationPhasePage extends SystemTestPageObject {
  constructor(
    page: Page,
    private readonly prefix: PhasePrefix,
    private readonly aiModelSelectionTestId: string,
    private readonly progressTestId: string
  ) {
    super(page)
  }

  get screen(): Locator {
    return this.byTestId(`${this.prefix}-screen`)
  }

  get aiModelSelection(): Locator {
    return this.byTestId(this.aiModelSelectionTestId)
  }

  get aiModelLockState(): Locator {
    return this.byTestId(`${this.prefix}-ai-model-lock-state`)
  }

  get statusPanel(): Locator {
    return this.byTestId(`${this.prefix}-status-panel`)
  }

  get progress(): Locator {
    return this.byTestId(this.progressTestId)
  }

  get progressBar(): Locator {
    return this.byTestId(`${this.prefix}-progress-bar`)
  }

  get progressCounts(): Locator {
    return this.byTestId(`${this.prefix}-progress-counts`)
  }

  get startButton(): Locator {
    return this.byTestId(`${this.prefix}-start-button`)
  }

  /** 中断ボタン（aria-label=中断）。H-2 の有効化条件確認に使う。 */
  get pauseButton(): Locator {
    return this.page.getByRole("button", { name: "中断" })
  }

  /** 再開ボタン（aria-label=再開）。H-3 の有効化条件確認に使う。 */
  get resumeButton(): Locator {
    return this.page.getByRole("button", { name: "再開" })
  }

  /** リトライボタン（aria-label=リトライ）。H-4 の有効化条件確認に使う。 */
  get retryButton(): Locator {
    return this.page.getByRole("button", { name: "リトライ" })
  }

  /** キャンセルボタン（aria-label=キャンセル）。H-10/H-16 の有効化条件確認に使う。 */
  get cancelButton(): Locator {
    return this.page.getByRole("button", { name: "キャンセル" })
  }

  get processingTargetListRegion(): Locator {
    return this.screen.getByRole("region", { name: "処理対象一覧" })
  }

  get processingTargetSearchInput(): Locator {
    return this.byTestId(this.processingTargetTestId("search-input"))
  }

  get processingTargetTotalCount(): Locator {
    return this.byTestId(this.processingTargetTestId("total"))
  }

  get processingTargetEmptyState(): Locator {
    return this.byTestId(this.processingTargetTestId("empty"))
  }

  get processingTargetRows(): Locator {
    return this.byTestId(this.processingTargetTestId("row"))
  }

  get processingTargetLoadingOverlay(): Locator {
    return this.byTestId(this.processingTargetTestId("loading"))
  }

  async waitForScreen(): Promise<void> {
    await this.waitFor(this.screen)
    // 初回取得中は phase-loading-overlay がフェーズ画面全体を覆う。
    // initialFetchDone=true になるとオーバーレイが消えるため、hidden を待つ。
    // オーバーレイが最初から存在しない場合（initialFetchDone=true 初期値）も hidden で通過する。
    await this.processingTargetLoadingOverlay.waitFor({ state: "hidden" })
  }

  async searchProcessingTargets(query: string): Promise<void> {
    await this.processingTargetSearchInput.fill(query)
  }

  async start(): Promise<void> {
    await this.startButton.click()
  }

  private processingTargetTestId(suffix: string): string {
    return `${this.prefix}-processing-target-${suffix}`
  }
}

export class TermTranslationPhasePage extends TranslationPhasePage {
  constructor(page: Page) {
    super(
      page,
      "term-translation-phase",
      "term-translation-phase-ai-model-selection-region",
      "term-translation-phase-progress-region"
    )
  }

  /** AI モデル設定パネルの状態 pill（「固定済み」または「設定未完了」を表示する要素）。 */
  get aiSettingsStatusPill(): Locator {
    return this.byTestId("term-translation-phase-ai-settings-status-pill")
  }

  /** 開始ボタンの禁止理由を表示する領域。 */
  get startBlockedReason(): Locator {
    return this.byTestId("term-translation-phase-start-blocked-reason")
  }

  /** AI サービス select 要素。 */
  get aiProviderSelect(): Locator {
    return this.byTestId("term-translation-phase-ai-provider-select")
  }

  /** モデル select 要素。 */
  get aiModelSelect(): Locator {
    return this.byTestId("term-translation-phase-ai-model-select")
  }

  /** 処理方式 select 要素。 */
  get aiExecutionModeSelect(): Locator {
    return this.byTestId("term-translation-phase-ai-execution-mode-select")
  }

  /** 「モデル一覧を更新」ボタン。AI 設定を保存して summary を再取得する。 */
  get refreshAIModelButton(): Locator {
    return this.aiModelSelection.getByRole("button", { name: "モデル一覧を更新" })
  }
}

export class PersonaGenerationPhasePage extends TranslationPhasePage {
  constructor(page: Page) {
    super(
      page,
      "persona-generation-phase",
      "persona-generation-phase-ai-model-selection-card",
      "persona-generation-phase-progress-card"
    )
  }

  /** AI モデル設定パネルの状態 pill（「固定済み」または「設定未完了」を表示する要素）。 */
  get aiSettingsStatusPill(): Locator {
    return this.byTestId("persona-generation-phase-ai-settings-status-pill")
  }

  /** 開始ボタンの禁止理由を表示する領域。 */
  get startBlockedReason(): Locator {
    return this.byTestId("persona-generation-phase-start-blocked-reason")
  }

  /** 次段階へボタン（persona → body 移行）。H-11 の有効化条件確認に使う。 */
  get nextPhaseButton(): Locator {
    return this.byTestId("persona-generation-phase-start-next-phase-button")
  }

  /** AI サービス select 要素。 */
  get aiProviderSelect(): Locator {
    return this.byTestId("persona-generation-phase-ai-provider-select")
  }

  /** モデル select 要素。 */
  get aiModelSelect(): Locator {
    return this.byTestId("persona-generation-phase-ai-model-select")
  }

  /** 「モデル一覧を更新」ボタン。モデル一覧を取得してモデル選択肢を更新する。 */
  get refreshAIModelButton(): Locator {
    return this.aiModelSelection.getByRole("button", { name: "モデル一覧を更新" })
  }
}

export class BodyTranslationPhasePage extends TranslationPhasePage {
  constructor(page: Page) {
    super(
      page,
      "body-translation-phase",
      "body-translation-phase-ai-model-selection",
      "body-translation-phase-progress"
    )
  }

  /** AI モデル設定パネルの状態 pill（「固定済み」または「設定未完了」を表示する要素）。 */
  get aiSettingsStatusPill(): Locator {
    return this.byTestId("body-translation-phase-ai-settings-status-pill")
  }

  /** 開始ボタンの禁止理由を表示する領域。 */
  get startBlockedReason(): Locator {
    return this.byTestId("body-translation-phase-start-blocked-reason")
  }

  /** AI サービス select 要素。 */
  get aiProviderSelect(): Locator {
    return this.byTestId("body-translation-phase-ai-provider-select")
  }

  /** モデル select 要素。 */
  get aiModelSelect(): Locator {
    return this.byTestId("body-translation-phase-ai-model-select")
  }

  /** 「モデル一覧を更新」ボタン。モデル一覧を取得してモデル選択肢を更新する。 */
  get refreshAIModelButton(): Locator {
    return this.aiModelSelection.getByRole("button", { name: "モデル一覧を更新" })
  }
}
