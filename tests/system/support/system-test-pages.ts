import type { Locator, Page } from "@playwright/test"

import { DashboardPage } from "./dashboard-page"
import { ProviderSettingsPage } from "./provider-settings-page"
import {
  SYSTEM_TEST_MOBILE_VIEWPORT,
  SystemTestPageObject
} from "./system-test-page-object"
import { TranslationJobManagementPage } from "./translation-job-management-page"

export { DashboardPage } from "./dashboard-page"
export { JobRunShellPage } from "./job-run-shell-page"
export { MasterDictionaryPage } from "./master-dictionary-page"
export { MasterPersonaPage } from "./master-persona-page"
export { OutputManagementPage } from "./output-management-page"
export { ProviderSettingsPage } from "./provider-settings-page"
export {
  BodyTranslationPhasePage,
  PersonaGenerationPhasePage,
  TermTranslationPhasePage,
  TranslationPhasePage
} from "./translation-phase-pages"
export { TranslationInputReviewPage } from "./translation-input-review-page"
export { TranslationJobManagementPage } from "./translation-job-management-page"
export { SYSTEM_TEST_MOBILE_VIEWPORT, SystemTestPageObject }

export async function openDashboard(page: Page): Promise<void> {
  await new DashboardPage(page).open()
}

export async function openRouteByHash(
  page: Page,
  routeHash: string,
  screen: Locator
): Promise<void> {
  await page.goto(routeHash)
  await screen.waitFor({ state: "visible" })
}

export function dashboardCard(page: Page, label: string): Locator {
  return new DashboardPage(page).card(label)
}

export function navigationItem(page: Page, label: string): Locator {
  return new DashboardPage(page).navigationItem(label)
}

export function providerRow(page: Page, label: string): Locator {
  return new ProviderSettingsPage(page).serviceRow(label)
}

export function jobCard(page: Page, text: string): Locator {
  return new TranslationJobManagementPage(page).jobCard(text)
}
