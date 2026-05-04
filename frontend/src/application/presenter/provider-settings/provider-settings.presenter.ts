import type {
  ProviderSettingsProviderListItemViewModel,
  ProviderSettingsProviderState,
  ProviderSettingsScreenState,
  ProviderSettingsScreenViewModel,
  ProviderSettingsSelectedProviderViewModel
} from "@application/gateway-contract/provider-settings"

function buildPhaseLabel(
  phase: ProviderSettingsScreenState["phase"]
): string {
  switch (phase) {
    case "loading":
      return "読込中"
    case "saving":
      return "保存中"
    case "resetting":
      return "リセット中"
    case "validating":
      return "接続確認中"
    case "ready":
      return "利用可能"
    case "idle":
    default:
      return "待機中"
  }
}

function buildStatus(
  provider: ProviderSettingsProviderState
): Pick<
  ProviderSettingsProviderListItemViewModel,
  "statusLabel" | "statusTone" | "helperText"
> {
  if (provider.validationState === "validated") {
    return {
      statusLabel: "利用可能",
      statusTone: "success",
      helperText: "接続確認済みです。"
    }
  }

  if (provider.validationState === "pending") {
    return {
      statusLabel: "接続確認中",
      statusTone: "neutral",
      helperText: "現在の入力内容で接続確認しています。"
    }
  }

  if (provider.lastFailureKind === "validation_stale") {
    return {
      statusLabel: "再確認が必要",
      statusTone: "warning",
      helperText: "エンドポイント変更後は接続確認をやり直してください。"
    }
  }

  if (provider.lastFailureKind === "credential_missing") {
    return {
      statusLabel: "APIキー未設定",
      statusTone: "warning",
      helperText: "APIキーを設定してから接続確認してください。"
    }
  }

  if (provider.lastFailureKind === "endpoint_missing") {
    return {
      statusLabel: "未設定",
      statusTone: "warning",
      helperText: "エンドポイントを入力してから保存してください。"
    }
  }

  if (provider.validationState === "failed") {
    return {
      statusLabel: "接続失敗",
      statusTone: "warning",
      helperText: "エンドポイントを確認してから再実行してください。"
    }
  }

  if (provider.savedState === "configured") {
    return {
      statusLabel: "保存済み",
      statusTone: "neutral",
      helperText: "保存済みの設定を確認できます。"
    }
  }

  if (provider.savedState === "partial") {
    return {
      statusLabel: "一部設定済み",
      statusTone: "warning",
      helperText: "保存後に接続確認が必要です。"
    }
  }

  return {
    statusLabel: "未設定",
    statusTone: "warning",
    helperText: "エンドポイントと APIキー状態を確認してください。"
  }
}

function buildApiKeyStateLabel(provider: ProviderSettingsProviderState): string {
  if (provider.credentialState === "not_required") {
    return "不要"
  }

  if (provider.credentialState === "configured") {
    return "保存済み"
  }

  return "未保存"
}

function buildValidationLabel(
  provider: ProviderSettingsProviderState
): {
  label: string
  tone: "neutral" | "warning" | "success"
  helpText: string
} {
  if (provider.validationState === "validated") {
    return {
      label: "確認済み",
      tone: "success",
      helpText: "現在の入力内容で接続確認済みです。"
    }
  }

  if (provider.validationState === "pending") {
    return {
      label: "接続確認中",
      tone: "neutral",
      helpText: "確認が終わるまでお待ちください。"
    }
  }

  if (provider.lastFailureKind === "validation_stale") {
    return {
      label: "再確認待ち",
      tone: "warning",
      helpText: "エンドポイント変更後は接続確認をやり直してください。"
    }
  }

  if (provider.lastFailureKind === "credential_missing") {
    return {
      label: "未確認",
      tone: "warning",
      helpText: "APIキーを設定してから接続確認してください。"
    }
  }

  if (provider.lastFailureKind === "endpoint_missing") {
    return {
      label: "未確認",
      tone: "warning",
      helpText: "エンドポイントを入力してから接続確認してください。"
    }
  }

  if (provider.validationState === "failed") {
    return {
      label: "接続失敗",
      tone: "warning",
      helpText: "接続先を確認してから再実行してください。"
    }
  }

  return {
    label: "未確認",
    tone: "neutral",
    helpText: "保存後に接続確認してください。"
  }
}

function buildEndpointHint(provider: ProviderSettingsProviderState): string {
  if (provider.providerId === "lm_studio") {
    return "ローカルの AI サービス接続先を保存します。"
  }

  return "利用する AI サービスの接続先を保存します。"
}

function buildPageLead(selectedProviderCount: number): string {
  return `AIサービスごとのエンドポイントと APIキー状態を管理します。利用可能な AIサービスは ${selectedProviderCount} 件です。`
}

export class ProviderSettingsPresenter {
  toViewModel(
    state: ProviderSettingsScreenState,
    isGatewayConnected: boolean
  ): ProviderSettingsScreenViewModel {
    const providerList = state.providers.map((provider) => {
      const status = buildStatus(provider)
      return {
        providerId: provider.providerId,
        label: provider.label,
        selected: provider.providerId === state.selectedProviderId,
        statusLabel: status.statusLabel,
        statusTone: status.statusTone,
        helperText: status.helperText
      }
    })

    const selectedProviderState =
      state.providers.find(
        (provider) => provider.providerId === state.selectedProviderId
      ) ?? null

    const selectedProvider: ProviderSettingsSelectedProviderViewModel | null =
      selectedProviderState
        ? (() => {
            const validation = buildValidationLabel(selectedProviderState)
            return {
              providerId: selectedProviderState.providerId,
              label: selectedProviderState.label,
              endpoint: selectedProviderState.endpointDraft,
              endpointPlaceholder:
                selectedProviderState.providerId === "lm_studio"
                  ? "http://127.0.0.1:1234/v1"
                  : "https://",
              endpointHint: buildEndpointHint(selectedProviderState),
              apiKeyStateLabel: buildApiKeyStateLabel(selectedProviderState),
              apiKeyHelpText:
                selectedProviderState.credentialState === "not_required"
                  ? "この AIサービスでは APIキーは不要です。"
                  : "APIキーは OS標準の資格情報マネージャーへ保存します。",
              apiKeyRequired:
                selectedProviderState.credentialState !== "not_required",
              apiKeyPanelOpen: state.apiKeyPanelOpen,
              validationLabel: validation.label,
              validationTone: validation.tone,
              validationHelpText: validation.helpText,
              canValidate: state.phase !== "validating",
              canSave: state.phase !== "saving",
              canReset: state.phase !== "resetting"
            }
          })()
        : null

    return {
      gatewayStatus: isGatewayConnected ? "接続準備済み" : "未接続",
      pageTitle: "AIサービス設定",
      pageLead: buildPageLead(state.providers.length),
      providerCountLabel: `${state.providers.length} 件の AIサービスを管理します。`,
      phaseLabel: buildPhaseLabel(state.phase),
      saveNotice: state.saveNotice,
      errorMessage: state.errorMessage,
      providerList,
      selectedProvider
    }
  }
}
