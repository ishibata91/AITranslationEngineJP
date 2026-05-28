import type { Locator, Page } from "@playwright/test"

export const SYSTEM_TEST_MOBILE_VIEWPORT = {
  height: 812,
  width: 375
} as const

export class SystemTestPageObject {
  constructor(protected readonly page: Page) {}

  protected byTestId(testId: string): Locator {
    return this.page.getByTestId(testId)
  }

  protected async openHashRoute(routeHash: string): Promise<void> {
    await this.page.goto(routeHash)
  }

  protected async waitFor(locator: Locator): Promise<void> {
    await locator.waitFor({ state: "visible" })
  }

  protected async clickButtonIn(
    region: Locator,
    accessibleName: string
  ): Promise<void> {
    await region.getByRole("button", { name: accessibleName }).click()
  }
}
