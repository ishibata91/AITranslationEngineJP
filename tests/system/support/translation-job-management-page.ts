import type { Locator, Page } from "@playwright/test"

import { SystemTestPageObject } from "./system-test-page-object"

export class TranslationJobManagementPage extends SystemTestPageObject {
  constructor(page: Page) {
    super(page)
  }

  get childScreenRegion(): Locator {
    return this.byTestId("translation-management-child-screen-region")
  }

  get jobListRegion(): Locator {
    return this.byTestId("translation-job-management-job-list-region")
  }

  get searchField(): Locator {
    return this.byTestId("translation-job-management-search-field")
  }

  get stateFilter(): Locator {
    return this.byTestId("translation-job-management-state-filter")
  }

  get jobCards(): Locator {
    return this.byTestId("translation-job-management-job-card")
  }

  get listStatusDisplay(): Locator {
    return this.byTestId("translation-job-management-list-status-display")
  }

  get feedbackNotification(): Locator {
    return this.byTestId("translation-job-management-feedback-notification")
  }

  get deleteModal(): Locator {
    return this.byTestId("translation-job-management-delete-confirmation-modal")
  }

  get deleteCancelButton(): Locator {
    return this.byTestId("translation-job-management-delete-cancel-button")
  }

  get deleteConfirmButton(): Locator {
    return this.byTestId("translation-job-management-delete-confirm-button")
  }

  async open(): Promise<void> {
    await this.openHashRoute("/#translation-management")
    await this.waitFor(this.jobListRegion)
  }

  jobCard(text: string): Locator {
    return this.jobCards.filter({ hasText: text })
  }

  jobSelectionRegion(card: Locator): Locator {
    return card.getByTestId("translation-job-management-job-selection-region")
  }

  stateLabel(card: Locator): Locator {
    return card.getByTestId("translation-job-management-state-label")
  }

  disabledReason(card: Locator): Locator {
    return card.getByTestId("translation-job-management-disabled-reason")
  }

  cardActions(card: Locator): Locator {
    return card.getByTestId("translation-job-management-job-actions")
  }

  openCurrentPhaseButton(card: Locator): Locator {
    // JobOperationGroup の continue-button は data-current-phase-action-state で識別する。
    // ボタンの表示テキストは「再開」だが、現在の翻訳段階を開く操作専用の属性で区別する。
    return this.cardActions(card).locator("[data-current-phase-action-state]")
  }

  stopButton(card: Locator): Locator {
    return this.cardActions(card).getByRole("button", { name: "停止" })
  }

  resumeButton(card: Locator): Locator {
    // continue-button（翻訳段階を開く「再開」ボタン）も含む。
    // E2E-UC-019 では continue-button を使って再開操作の UI 経路を確認する。
    // SCN-TJM-005/007 と E2E-UC-038 で resume 操作の不在確認には openCurrentPhaseButton を使う。
    return this.cardActions(card).getByRole("button", { name: "再開" })
  }

  deleteButton(card: Locator): Locator {
    return this.cardActions(card).getByRole("button", { name: "削除" })
  }

  async search(text: string): Promise<void> {
    await this.searchField.fill(text)
  }

  async selectState(label: string): Promise<void> {
    const optionValue = await this.stateFilter.evaluate(
      (select, optionLabel) => {
        if (!(select instanceof HTMLSelectElement)) {
          return null
        }
        const option = Array.from(select.options).find((candidate) =>
          candidate.textContent?.trim().startsWith(optionLabel)
        )
        return option?.value ?? null
      },
      label
    )
    if (optionValue === null) {
      throw new Error(`state filter option was not found: ${label}`)
    }
    await this.stateFilter.selectOption(optionValue)
  }

  async openCurrentPhase(card: Locator): Promise<void> {
    await this.openCurrentPhaseButton(card).click()
  }

  async stop(card: Locator): Promise<void> {
    await this.stopButton(card).click()
  }

  async resume(card: Locator): Promise<void> {
    await this.resumeButton(card).click()
  }

  async openDelete(card: Locator): Promise<void> {
    await this.deleteButton(card).click()
    await this.waitFor(this.deleteModal)
  }

  async cancelDelete(): Promise<void> {
    await this.deleteCancelButton.first().click()
  }

  async confirmDelete(): Promise<void> {
    await this.deleteConfirmButton.click()
  }
}
