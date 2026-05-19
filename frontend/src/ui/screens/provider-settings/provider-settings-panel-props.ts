import type {
  ProviderSettingsProviderListItemViewModel,
  ProviderSettingsSelectedProviderViewModel
} from "@application/gateway-contract/provider-settings"

export interface ProviderSettingsSummaryPanelProps {
  gatewayStatus: string
  pageTitle: string
  pageLead: string
  providerCountLabel: string
  phaseLabel: string
  saveNotice: string
  errorMessage: string
}

export interface ProviderListPanelProps {
  providerList: ProviderSettingsProviderListItemViewModel[]
  selectProvider: (providerId: string) => void
}

export interface ProviderDetailPanelProps {
  selectedProvider: ProviderSettingsSelectedProviderViewModel
  updateEndpoint: (event: Event) => void
}

export interface ApiKeyPanelProps {
  selectedProvider: ProviderSettingsSelectedProviderViewModel
  credentialInputDraft: string
  updateCredentialDraft: (nextValue: string) => void
  openApiKeyPanel: () => void
  saveSettings: () => void
  closeApiKeyPanel: () => void
}

export interface ConnectionCheckPanelProps {
  selectedProvider: ProviderSettingsSelectedProviderViewModel
  validateConnection: () => void
}

export interface SettingsActionPanelProps {
  selectedProvider: ProviderSettingsSelectedProviderViewModel
  saveSettings: () => void
  resetSettings: () => void
}
