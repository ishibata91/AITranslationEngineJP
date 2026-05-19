import type {
  ApiKeyPanelProps,
  ConnectionCheckPanelProps,
  ProviderDetailPanelProps,
  ProviderListPanelProps,
  ProviderSettingsSummaryPanelProps,
  SettingsActionPanelProps
} from "../provider-settings-panel-props"

const ignoreAction = (): void => {}
const ignoreEvent = (event: Event): void => {
  void event
}
const ignoreDraft = (nextValue: string): void => {
  void nextValue
}
const ignoreProvider = (providerId: string): void => {
  void providerId
}

const selectedProvider = {
  providerId: "gemini",
  label: "Sample Provider",
  endpoint: "synthetic-endpoint-reference",
  endpointPlaceholder: "synthetic endpoint label",
  endpointHint: "接続先は設定済みの参照名として確認します。",
  apiKeyStateLabel: "設定済み",
  apiKeyHelpText: "APIキーは保存済みの参照状態だけを表示します。",
  apiKeyRequired: true,
  apiKeyPanelOpen: false,
  validationLabel: "確認済み",
  validationTone: "success",
  validationHelpText: "接続確認は synthetic provider の状態だけを表示します。",
  canValidate: true,
  canSave: true,
  canReset: true
} as const

const providerList = [
  {
    providerId: "gemini",
    label: "Sample Provider",
    statusLabel: "設定済み",
    statusTone: "success",
    helperText: "保存済みの資格情報参照で利用できます。",
    selected: true
  },
  {
    providerId: "lm_studio",
    label: "Local Provider",
    statusLabel: "APIキー不要",
    statusTone: "neutral",
    helperText: "ローカル実行用の参照状態です。",
    selected: false
  },
  {
    providerId: "xai",
    label: "Review Provider",
    statusLabel: "設定が必要",
    statusTone: "warning",
    helperText: "資格情報参照を設定すると確認できます。",
    selected: false
  }
] as const

export const providerSettingsSummaryPanelFixtures = {
  ready: {
    gatewayStatus: "利用可能",
    pageTitle: "AIサービス設定",
    pageLead:
      "AIサービスごとのエンドポイントと APIキー状態を管理します。利用可能な AIサービスは 3 件です。",
    providerCountLabel: "3 件の AIサービスを確認できます。",
    phaseLabel: "設定確認",
    saveNotice: "",
    errorMessage: ""
  },
  saved: {
    gatewayStatus: "利用可能",
    pageTitle: "AIサービス設定",
    pageLead:
      "AIサービスごとのエンドポイントと APIキー状態を管理します。利用可能な AIサービスは 3 件です。",
    providerCountLabel: "3 件の AIサービスを確認できます。",
    phaseLabel: "保存済み",
    saveNotice: "設定を保存しました。",
    errorMessage: ""
  },
  failed: {
    gatewayStatus: "確認不可",
    pageTitle: "AIサービス設定",
    pageLead:
      "AIサービスごとのエンドポイントと APIキー状態を管理します。利用可能な AIサービスは 3 件です。",
    providerCountLabel: "3 件の AIサービスを確認できます。",
    phaseLabel: "確認失敗",
    saveNotice: "",
    errorMessage: "設定を確認できませんでした。"
  }
} satisfies Record<string, ProviderSettingsSummaryPanelProps>

export const providerListPanelFixtures = {
  mixed: {
    providerList: [...providerList],
    selectProvider: ignoreProvider
  },
  empty: {
    providerList: [],
    selectProvider: ignoreProvider
  }
} satisfies Record<string, ProviderListPanelProps>

export const providerDetailPanelFixtures = {
  selected: {
    selectedProvider,
    updateEndpoint: ignoreEvent
  },
  warning: {
    selectedProvider: {
      ...selectedProvider,
      validationLabel: "未確認",
      validationTone: "warning",
      validationHelpText: "設定を保存してから接続を確認します。",
      canValidate: false
    },
    updateEndpoint: ignoreEvent
  }
} satisfies Record<string, ProviderDetailPanelProps>

export const apiKeyPanelFixtures = {
  maskedState: {
    selectedProvider,
    credentialInputDraft: "",
    updateCredentialDraft: ignoreDraft,
    openApiKeyPanel: ignoreAction,
    saveSettings: ignoreAction,
    closeApiKeyPanel: ignoreAction
  },
  draftOpen: {
    selectedProvider: {
      ...selectedProvider,
      apiKeyPanelOpen: true,
      apiKeyStateLabel: "入力中"
    },
    credentialInputDraft: "",
    updateCredentialDraft: ignoreDraft,
    openApiKeyPanel: ignoreAction,
    saveSettings: ignoreAction,
    closeApiKeyPanel: ignoreAction
  },
  notRequired: {
    selectedProvider: {
      ...selectedProvider,
      apiKeyRequired: false,
      apiKeyStateLabel: "APIキー不要",
      apiKeyHelpText: "ローカル provider は APIキーを使いません。"
    },
    credentialInputDraft: "",
    updateCredentialDraft: ignoreDraft,
    openApiKeyPanel: ignoreAction,
    saveSettings: ignoreAction,
    closeApiKeyPanel: ignoreAction
  }
} satisfies Record<string, ApiKeyPanelProps>

export const connectionCheckPanelFixtures = {
  ready: {
    selectedProvider,
    validateConnection: ignoreAction
  },
  disabled: {
    selectedProvider: {
      ...selectedProvider,
      validationLabel: "未確認",
      validationTone: "warning",
      validationHelpText: "APIキー状態を確認してから接続確認できます。",
      canValidate: false
    },
    validateConnection: ignoreAction
  }
} satisfies Record<string, ConnectionCheckPanelProps>

export const settingsActionPanelFixtures = {
  enabled: {
    selectedProvider,
    saveSettings: ignoreAction,
    resetSettings: ignoreAction
  },
  disabled: {
    selectedProvider: {
      ...selectedProvider,
      canSave: false,
      canReset: false
    },
    saveSettings: ignoreAction,
    resetSettings: ignoreAction
  }
} satisfies Record<string, SettingsActionPanelProps>
