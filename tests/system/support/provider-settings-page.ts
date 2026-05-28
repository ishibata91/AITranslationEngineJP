import type { Locator, Page } from "@playwright/test"

import { SystemTestPageObject } from "./system-test-page-object"

export class ProviderSettingsPage extends SystemTestPageObject {
  constructor(page: Page) {
    super(page)
  }

  get screen(): Locator {
    return this.byTestId("provider-settings-screen-shell")
  }

  get summaryRegion(): Locator {
    return this.byTestId("provider-settings-screen-summary-region")
  }

  get apiKeyStatusRegion(): Locator {
    return this.byTestId("provider-settings-api-key-status-region")
  }

  get serviceRows(): Locator {
    return this.byTestId("provider-settings-ai-service-row")
  }

  get apiKeyInputRegion(): Locator {
    return this.byTestId("provider-settings-api-key-input-region")
  }

  get apiKeyInput(): Locator {
    return this.byTestId("provider-settings-api-key-input")
  }

  get apiKeyOpenButton(): Locator {
    return this.byTestId("provider-settings-api-key-open-button")
  }

  get apiKeySaveButton(): Locator {
    return this.byTestId("provider-settings-api-key-save-button")
  }

  get endpointInput(): Locator {
    return this.byTestId("provider-settings-endpoint-input")
  }

  get connectionCheckRegion(): Locator {
    return this.byTestId("provider-settings-connection-check-region")
  }

  get connectionCheckButton(): Locator {
    return this.byTestId("provider-settings-connection-check-button")
  }

  get settingsActionsRegion(): Locator {
    return this.byTestId("provider-settings-settings-actions-region")
  }

  get saveButton(): Locator {
    return this.byTestId("provider-settings-save-button")
  }

  get resetButton(): Locator {
    return this.byTestId("provider-settings-reset-button")
  }

  get settingsDetailRegion(): Locator {
    return this.byTestId("provider-settings-settings-detail-region")
  }

  async open(): Promise<void> {
    await this.openHashRoute("/#provider-settings")
    await this.waitFor(this.screen)
  }

  serviceRow(label: string): Locator {
    return this.byTestId("provider-settings-ai-service-row").filter({
      hasText: label
    })
  }

  async selectService(label: string): Promise<void> {
    await this.serviceRow(label).click()
  }

  async saveApiKey(secret: string): Promise<void> {
    await this.apiKeyOpenButton.click()
    await this.apiKeyInput.fill(secret)
    await this.apiKeySaveButton.click()
  }

  async fillEndpoint(endpoint: string): Promise<void> {
    await this.endpointInput.fill(endpoint)
  }

  async checkConnection(): Promise<void> {
    await this.connectionCheckButton.click()
  }

  async saveSettings(): Promise<void> {
    await this.saveButton.click()
  }

  async resetSettings(): Promise<void> {
    await this.resetButton.click()
  }
}
