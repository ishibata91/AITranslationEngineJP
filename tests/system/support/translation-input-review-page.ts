import type { Locator, Page } from "@playwright/test"

import { SystemTestPageObject } from "./system-test-page-object"

export class TranslationInputReviewPage extends SystemTestPageObject {
  constructor(page: Page) {
    super(page)
  }

  get statusHeader(): Locator {
    return this.byTestId("translation-input-review-screen-status-header")
  }

  get loadPreparationRegion(): Locator {
    return this.byTestId("translation-input-review-load-preparation-region")
  }

  get jsonFileInput(): Locator {
    return this.byTestId("translation-input-review-json-file-input")
  }

  get registerButton(): Locator {
    return this.byTestId("translation-input-review-register-button")
  }

  get registerError(): Locator {
    return this.byTestId("translation-input-review-register-error")
  }

  get loadedInputList(): Locator {
    return this.byTestId("translation-input-review-loaded-input-list")
  }

  get nextActionFooter(): Locator {
    return this.byTestId("translation-input-review-next-action-footer")
  }

  async open(): Promise<void> {
    await this.openHashRoute("/#translation-management")
    await this.page.getByRole("button", { name: "新規翻訳を開始" }).click()
    await this.waitFor(this.loadPreparationRegion)
  }

  async setJsonFile(filePath: string): Promise<void> {
    await this.jsonFileInput.setInputFiles(filePath)
  }

  async register(): Promise<void> {
    await this.registerButton.click()
  }
}
