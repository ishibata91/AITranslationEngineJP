export type ModelSettingsCredentialStatus =
  | "configured"
  | "missing"
  | "not_required"

export type ModelSettingsModelListStatus =
  | "not_updated"
  | "loading"
  | "success"
  | "failed"
  | "credential_missing"
  | "credential_not_required"

export type ModelSettingsSaveStatus =
  | "clean"
  | "dirty"
  | "saving"
  | "saved"
  | "failed"

export interface ModelSettingsModelOption {
  modelId: string
  label: string
}

export interface ModelSettingsProviderOption {
  value: string
  label: string
  credentialStatus: ModelSettingsCredentialStatus
}

export interface ModelSettingsModelListState {
  provider: string
  credentialStatus: ModelSettingsCredentialStatus
  status: ModelSettingsModelListStatus
  models: ModelSettingsModelOption[]
  failureKind?: string
}

export interface ModelSettingsCardState {
  referenceId: string
  provider: string
  model: string
  credentialStatus: ModelSettingsCredentialStatus
  modelList: ModelSettingsModelListState
  saveStatus: ModelSettingsSaveStatus
  saveMessage: string
}

export interface ModelSettingsCardViewModel {
  referenceId: string
  provider: string
  model: string
  providerOptions: Array<{ value: string; label: string }>
  credentialStatusLabel: string
  credentialStatusTone: "neutral" | "warning" | "success"
  showCredentialStatus: boolean
  showCredentialWarning: boolean
  credentialWarningText: string
  modelListButtonEnabled: boolean
  modelListButtonLabel: string
  modelListButtonAriaLabel: string
  isModelListRefreshing: boolean
  modelListStatusText: string
  modelOptions: ModelSettingsModelOption[]
  modelSelectEnabled: boolean
  emptyModelLabel: string
  statusLabel: string
  statusTone: "neutral" | "warning" | "success"
  helperText: string
  footerMessage: string
  footerWarningText: string
  actionButtonDisabled: boolean
}
