import type { Locator, Page } from "@playwright/test"

import { SystemTestPageObject } from "./system-test-page-object"

export class DashboardPage extends SystemTestPageObject {
  constructor(page: Page) {
    super(page)
  }

  get content(): Locator {
    return this.byTestId("dashboard-dashboard-content")
  }

  get applicationHeader(): Locator {
    return this.byTestId("dashboard-application-header")
  }

  get entryRegion(): Locator {
    return this.byTestId("dashboard-primary-page-entry-region")
  }

  get globalNavigation(): Locator {
    return this.byTestId("dashboard-global-navigation")
  }

  get currentPageDescription(): Locator {
    return this.byTestId("dashboard-current-page-description")
  }

  get mobileNavigationToggle(): Locator {
    return this.byTestId("dashboard-mobile-navigation-toggle")
  }

  async open(): Promise<void> {
    await this.openHashRoute("/")
    await this.waitFor(this.content)
  }

  card(label: string): Locator {
    return this.byTestId("dashboard-primary-page-card").filter({
      hasText: label
    })
  }

  cardOpenButton(label: string): Locator {
    return this.card(label).getByTestId("dashboard-primary-page-open-button")
  }

  navigationItem(label: string): Locator {
    return this.byTestId("dashboard-global-navigation-item").filter({
      hasText: label
    })
  }

  async openCard(label: string): Promise<void> {
    await this.card(label).click()
  }

  async openMobileNavigation(): Promise<void> {
    await this.mobileNavigationToggle.click()
    await this.waitFor(this.globalNavigation)
  }

  async selectNavigationItem(label: string): Promise<void> {
    await this.navigationItem(label).click()
  }
}
