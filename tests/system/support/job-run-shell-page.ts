import type { Locator, Page } from "@playwright/test"

import { SystemTestPageObject } from "./system-test-page-object"

export class JobRunShellPage extends SystemTestPageObject {
  constructor(page: Page) {
    super(page)
  }

  get shell(): Locator {
    return this.byTestId("job-run-job-run-shell")
  }

  get selectedJobSummary(): Locator {
    return this.byTestId("job-run-selected-job-summary")
  }

  get phaseScreenRegion(): Locator {
    return this.byTestId("job-run-phase-screen-region")
  }

  get nextActionFooter(): Locator {
    return this.byTestId("job-run-next-action-footer")
  }

  get bodyCompleteNextAction(): Locator {
    return this.byTestId("body-translation-phase-complete-next-action")
  }

  get translationCompleteScreen(): Locator {
    return this.byTestId("translation-complete-translation-complete-screen")
  }

  get translationCompleteSummary(): Locator {
    return this.byTestId("translation-complete-completion-summary")
  }

  get phaseActionRegion(): Locator {
    return this.phaseScreenRegion
  }

  get postCompletionNextAction(): Locator {
    return this.byTestId("translation-complete-post-completion-next-action")
  }

  async open(): Promise<void> {
    await this.openHashRoute("/#translation-management/job-run")
    await this.waitFor(this.shell)
  }

  async clickNext(): Promise<void> {
    await this.clickButtonIn(this.nextActionFooter, "次へ進む")
  }

  async clickBodyCompleteNext(): Promise<void> {
    await this.clickButtonIn(this.bodyCompleteNextAction, "次へ進む")
  }

  async retryCurrentPhase(): Promise<void> {
    await this.clickButtonIn(this.phaseActionRegion, "再実行")
  }

  async cancelCurrentPhase(): Promise<void> {
    await this.clickButtonIn(this.phaseActionRegion, "キャンセル")
  }

  async openOutputManagement(): Promise<void> {
    await this.clickButtonIn(this.postCompletionNextAction, "出力管理へ進む")
  }
}
