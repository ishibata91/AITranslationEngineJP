import type { ProviderSettingsScreenViewModel } from "./provider-settings-screen-types"

export type ProviderSettingsScreenViewModelListener = (
  viewModel: ProviderSettingsScreenViewModel
) => void

export interface ProviderSettingsScreenControllerContract {
  mount(): Promise<void>
  dispose(): void
  subscribe(listener: ProviderSettingsScreenViewModelListener): () => void
  getViewModel(): ProviderSettingsScreenViewModel
  selectProvider(providerId: string): void
  updateEndpoint(event: Event): void
  openApiKeyPanel(): void
  closeApiKeyPanel(): void
  saveSettings(): Promise<void>
  resetSettings(): Promise<void>
  validateConnection(): Promise<void>
}

export type CreateProviderSettingsScreenController =
  () => ProviderSettingsScreenControllerContract
