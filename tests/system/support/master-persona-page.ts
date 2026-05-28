import { expect, type Locator, type Page } from "@playwright/test"

import { SystemTestPageObject } from "./system-test-page-object"

export class MasterPersonaPage extends SystemTestPageObject {
  constructor(page: Page) {
    super(page)
  }

  get screenStatusRegion(): Locator {
    return this.byTestId("master-persona-screen-status-region")
  }

  get inputJsonPanel(): Locator {
    return this.byTestId("master-persona-input-json-panel")
  }

  get jsonFileInput(): Locator {
    return this.byTestId("master-persona-json-file-input")
  }

  get generateButton(): Locator {
    return this.byTestId("master-persona-generate-button")
  }

  get aiSettingsCard(): Locator {
    return this.byTestId("master-persona-ai-settings-card")
  }

  get aiServiceSelect(): Locator {
    return this.byTestId("master-persona-ai-service-select")
  }

  get modelSelect(): Locator {
    return this.byTestId("master-persona-model-select")
  }

  get executionModeSelect(): Locator {
    return this.byTestId("master-persona-execution-mode-select")
  }

  get progressPanel(): Locator {
    return this.byTestId("master-persona-progress-panel")
  }

  get resultListPanel(): Locator {
    return this.byTestId("master-persona-generation-result-list-panel")
  }

  get searchInput(): Locator {
    return this.byTestId("master-persona-search-input")
  }

  get pluginSelect(): Locator {
    return this.byTestId("master-persona-plugin-select")
  }

  get personaRows(): Locator {
    return this.byTestId("master-persona-row")
  }

  get editButton(): Locator {
    return this.byTestId("master-persona-edit-button")
  }

  get deleteButton(): Locator {
    return this.byTestId("master-persona-delete-button")
  }

  get editModal(): Locator {
    return this.byTestId("master-persona-edit-modal")
  }

  get deleteModal(): Locator {
    return this.byTestId("master-persona-delete-modal")
  }

  get summaryInput(): Locator {
    return this.byTestId("master-persona-summary-input")
  }

  get speechStyleInput(): Locator {
    return this.byTestId("master-persona-speech-style-input")
  }

  get bodyInput(): Locator {
    return this.byTestId("master-persona-body-input")
  }

  get editSaveButton(): Locator {
    return this.byTestId("master-persona-edit-save-button")
  }

  get editCancelButton(): Locator {
    return this.byTestId("master-persona-edit-cancel-button")
  }

  get deleteConfirmButton(): Locator {
    return this.byTestId("master-persona-delete-confirm-button")
  }

  get detailPanel(): Locator {
    return this.resultListPanel
  }

  async open(): Promise<void> {
    await this.openHashRoute("/#master-persona")
    await this.waitFor(this.screenStatusRegion)
  }

  personaRow(text: string): Locator {
    return this.personaRows.filter({ hasText: text })
  }

  async setJsonFile(filePath: string): Promise<void> {
    await this.jsonFileInput.setInputFiles(filePath)
  }

  async selectAISettings(input: {
    provider?: string
    model?: string
    executionMode?: string
  }): Promise<void> {
    if (input.provider !== undefined) {
      await this.waitForSelectableOption(this.aiServiceSelect, input.provider)
      await this.aiServiceSelect.selectOption({ label: input.provider })
    }
    if (input.model !== undefined) {
      await this.waitForSelectableOption(this.modelSelect, input.model)
      await this.modelSelect.selectOption({ label: input.model })
    }
    if (input.executionMode !== undefined) {
      await this.waitForSelectableOption(
        this.executionModeSelect,
        input.executionMode
      )
      await this.executionModeSelect.selectOption({
        label: input.executionMode
      })
    }
  }

  async waitForSelectableModel(label: string): Promise<void> {
    await this.waitForSelectableOption(this.modelSelect, label)
  }

  async generate(): Promise<void> {
    await this.generateButton.click()
  }

  async search(text: string): Promise<void> {
    await this.searchInput.fill(text)
  }

  async selectPlugin(label: string): Promise<void> {
    await this.waitForSelectableOption(this.pluginSelect, label)
    await this.pluginSelect.selectOption({ value: label })
    await expect(this.pluginSelect).toHaveValue(label)
    if (label !== "すべてのプラグイン") {
      await expect(this.resultListPanel).toContainText(label)
    }
  }

  async selectPersona(text: string): Promise<void> {
    const row = this.personaRow(text).first()
    await expect(row).toBeVisible()

    if ((await row.getAttribute("aria-expanded")) !== "true") {
      await row.click()
    }
    await expect(row).toHaveAttribute("aria-expanded", "true")
  }

  async openEditModal(): Promise<void> {
    await this.editButton.click()
    await this.waitFor(this.editModal)
  }

  async fillEdit(input: {
    summary?: string
    speechStyle?: string
    body?: string
  }): Promise<void> {
    if (input.summary !== undefined) {
      await this.summaryInput.fill(input.summary)
    }
    if (input.speechStyle !== undefined) {
      await this.speechStyleInput.fill(input.speechStyle)
    }
    if (input.body !== undefined) {
      await this.bodyInput.fill(input.body)
    }
  }

  async saveEdit(): Promise<void> {
    await this.editSaveButton.click()
    await expect(this.editModal).toBeHidden()
  }

  async cancelEdit(): Promise<void> {
    await this.editCancelButton.first().click()
  }

  async openDeleteModal(): Promise<void> {
    await this.deleteButton.click()
    await this.waitFor(this.deleteModal)
  }

  async confirmDelete(): Promise<void> {
    await this.deleteConfirmButton.click()
  }

  private async waitForSelectableOption(
    select: Locator,
    label: string
  ): Promise<void> {
    await expect(select).toBeEnabled()
    await expect
      .poll(
        async () =>
          select.evaluate((selectElement, optionLabel) => {
            if (!(selectElement instanceof HTMLSelectElement)) {
              return false
            }
            const option = Array.from(selectElement.options).find(
              (candidate) =>
                candidate.value === optionLabel ||
                candidate.label === optionLabel ||
                candidate.textContent?.trim() === optionLabel ||
                candidate.label.startsWith(`${optionLabel} `) ||
                candidate.textContent?.trim().startsWith(`${optionLabel} `)
            )
            return (
              option !== undefined && !option.disabled && !option.hidden
            )
          }, label),
        { message: `select option should be visible: ${label}` }
      )
      .toBe(true)
  }
}
