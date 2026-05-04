import type {
  ProviderSettingsCredentialState,
  ProviderSettingsErrorKind,
  ProviderSettingsProviderId,
  ProviderSettingsSavedState,
  ProviderSettingsValidationState
} from "./provider-settings-gateway-contract"

export interface ProviderSettingsProviderState {
  providerId: ProviderSettingsProviderId
  label: string
  endpointDraft: string
  persistedEndpoint: string
  credentialState: ProviderSettingsCredentialState
  credentialReferenceId?: string
  validationState: ProviderSettingsValidationState
  savedState: ProviderSettingsSavedState
  requestToken: string
  lastFailureKind?: ProviderSettingsErrorKind
}

export interface ProviderSettingsScreenState {
  phase: "idle" | "loading" | "ready" | "saving" | "resetting" | "validating"
  providers: ProviderSettingsProviderState[]
  selectedProviderId: ProviderSettingsProviderId | null
  apiKeyPanelOpen: boolean
  saveNotice: string
  errorMessage: string
}

export interface ProviderSettingsProviderListItemViewModel {
  providerId: ProviderSettingsProviderId
  label: string
  statusLabel: string
  statusTone: "neutral" | "warning" | "success"
  helperText: string
  selected: boolean
}

export interface ProviderSettingsSelectedProviderViewModel {
  providerId: ProviderSettingsProviderId
  label: string
  endpoint: string
  endpointPlaceholder: string
  endpointHint: string
  apiKeyStateLabel: string
  apiKeyHelpText: string
  apiKeyRequired: boolean
  apiKeyPanelOpen: boolean
  validationLabel: string
  validationTone: "neutral" | "warning" | "success"
  validationHelpText: string
  canValidate: boolean
  canSave: boolean
  canReset: boolean
}

export interface ProviderSettingsScreenViewModel {
  gatewayStatus: string
  pageTitle: string
  pageLead: string
  providerCountLabel: string
  phaseLabel: string
  saveNotice: string
  errorMessage: string
  providerList: ProviderSettingsProviderListItemViewModel[]
  selectedProvider: ProviderSettingsSelectedProviderViewModel | null
}
