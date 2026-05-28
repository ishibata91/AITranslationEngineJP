import type { Locator, Page } from "@playwright/test"

import { SystemTestPageObject } from "./system-test-page-object"

export class MasterDictionaryPage extends SystemTestPageObject {
  constructor(page: Page) {
    super(page)
  }

  get screen(): Locator {
    return this.byTestId("master-dictionary-master-dictionary-screen")
  }

  get operationRegion(): Locator {
    return this.byTestId("master-dictionary-dictionary-operation-region")
  }

  get detailPanel(): Locator {
    return this.operationRegion
  }

  get searchInput(): Locator {
    return this.byTestId("master-dictionary-search-input")
  }

  get categorySelect(): Locator {
    return this.byTestId("master-dictionary-category-select")
  }

  get createButton(): Locator {
    return this.byTestId("master-dictionary-create-button")
  }

  get entryRows(): Locator {
    return this.byTestId("master-dictionary-entry-row")
  }

  get detailEditButton(): Locator {
    return this.byTestId("master-dictionary-detail-edit-button")
  }

  get detailDeleteButton(): Locator {
    return this.byTestId("master-dictionary-detail-delete-button")
  }

  get createEditModal(): Locator {
    return this.byTestId("master-dictionary-create-edit-modal")
  }

  get sourceInput(): Locator {
    return this.byTestId("master-dictionary-entry-source-input")
  }

  get categoryInput(): Locator {
    return this.byTestId("master-dictionary-entry-category-select")
  }

  get originInput(): Locator {
    return this.byTestId("master-dictionary-entry-origin-input")
  }

  get translationInput(): Locator {
    return this.byTestId("master-dictionary-entry-translation-input")
  }

  get entrySaveButton(): Locator {
    return this.byTestId("master-dictionary-entry-save-button")
  }

  get entryValidationError(): Locator {
    return this.byTestId("master-dictionary-entry-validation-error")
  }

  get deleteModal(): Locator {
    return this.byTestId("master-dictionary-delete-confirmation-modal")
  }

  get deleteCancelButton(): Locator {
    return this.byTestId("master-dictionary-delete-cancel-button")
  }

  get deleteConfirmButton(): Locator {
    return this.byTestId("master-dictionary-delete-confirm-button")
  }

  get xmlImportRegion(): Locator {
    return this.byTestId("master-dictionary-xml-import-region")
  }

  get xmlFileInput(): Locator {
    return this.byTestId("master-dictionary-xml-file-input")
  }

  get xmlImportProgressPanel(): Locator {
    return this.byTestId("master-dictionary-import-progress-panel")
  }

  get xmlImportButton(): Locator {
    return this.byTestId("master-dictionary-xml-import-button")
  }

  async open(): Promise<void> {
    await this.openHashRoute("/#master-dictionary")
    await this.waitFor(this.screen)
  }

  entryRow(text: string): Locator {
    return this.entryRows.filter({ hasText: text })
  }

  emptyRow(): Locator {
    return this.operationRegion.getByText("処理対象がありません")
  }

  async search(text: string): Promise<void> {
    await this.searchInput.fill(text)
  }

  async selectCategory(value: string): Promise<void> {
    await this.categorySelect.selectOption(value)
  }

  async openCreateModal(): Promise<void> {
    await this.createButton.click()
    await this.waitFor(this.createEditModal)
  }

  async fillEntry(input: {
    source?: string
    category?: string
    origin?: string
    translation?: string
  }): Promise<void> {
    if (input.source !== undefined) {
      await this.sourceInput.fill(input.source)
    }
    if (input.category !== undefined) {
      await this.categoryInput.selectOption(input.category)
    }
    if (input.origin !== undefined) {
      await this.originInput.selectOption(input.origin)
    }
    if (input.translation !== undefined) {
      await this.translationInput.fill(input.translation)
    }
  }

  async saveEntry(): Promise<void> {
    await this.entrySaveButton.click()
  }

  async closeEntryModal(): Promise<void> {
    await this.clickButtonIn(this.createEditModal, "閉じる")
  }

  async openEditModal(): Promise<void> {
    await this.detailEditButton.click()
    await this.waitFor(this.createEditModal)
  }

  async openDeleteModal(): Promise<void> {
    await this.detailDeleteButton.click()
    await this.waitFor(this.deleteModal)
  }

  async confirmDelete(): Promise<void> {
    await this.deleteConfirmButton.click()
  }

  async cancelDelete(): Promise<void> {
    await this.deleteCancelButton.first().click()
  }

  async setXmlFile(filePath: string): Promise<void> {
    await this.xmlFileInput.setInputFiles(filePath)
  }

  async stageXmlFileWithRuntimePath(filePath: string): Promise<void> {
    await this.page.evaluate((absolutePath) => {
      const input = document.querySelector(
        "[data-testid='master-dictionary-xml-file-input']"
      )
      if (!(input instanceof HTMLInputElement)) {
        return
      }

      const fileName = absolutePath.split(/[\\/]/).at(-1) ?? "fixture.xml"
      const file = new File([""], fileName, { type: "text/xml" })
      Object.defineProperty(file, "path", {
        value: absolutePath,
        configurable: true
      })

      const transfer = new DataTransfer()
      transfer.items.add(file)
      Object.defineProperty(input, "files", {
        value: transfer.files,
        configurable: true
      })
      input.dispatchEvent(new Event("change", { bubbles: true }))
    }, filePath)
  }

  async startXmlImport(): Promise<void> {
    await this.xmlImportButton.click()
  }
}
