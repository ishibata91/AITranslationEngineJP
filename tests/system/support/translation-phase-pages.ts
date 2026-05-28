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

  async waitForScreen(): Promise<void> {
    await this.waitFor(this.screen)
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
}
