import type {
  ModelSettingsCardState,
  ModelSettingsCardViewModel,
  ModelSettingsCredentialStatus,
  ModelSettingsModelListState,
  ModelSettingsModelOption,
  ModelSettingsProviderOption,
  ModelSettingsSaveStatus
} from "./model-settings-card-contract"

const PROVIDER_LABELS: Record<string, string> = {
  gemini: "Gemini",
  lm_studio: "LM Studio",
  xai: "xAI"
}

const USER_VISIBLE_PROVIDERS = new Set(Object.keys(PROVIDER_LABELS))

function isUsableModelList(
  status: ModelSettingsModelListState["status"]
): boolean {
  return status === "success" || status === "credential_not_required"
}

function formatProviderLabel(provider: string): string {
  return PROVIDER_LABELS[provider] ?? provider
}

function createModelSettingsModelListState(options: {
  provider: string
  credentialStatus: ModelSettingsCredentialStatus
  status?: ModelSettingsModelListState["status"]
  models?: ModelSettingsModelOption[]
  failureKind?: string
}): ModelSettingsModelListState {
  return {
    provider: options.provider,
    credentialStatus: options.credentialStatus,
    status: options.status ?? "not_updated",
    models: (options.models ?? []).map((model) => ({ ...model })),
    failureKind: options.failureKind
  }
}

export function createModelSettingsCardState(options: {
  referenceId: string
  provider: string
  model: string
  credentialStatus: ModelSettingsCredentialStatus
  modelList?: ModelSettingsModelListState
  saveStatus?: ModelSettingsSaveStatus
  saveMessage?: string
}): ModelSettingsCardState {
  return {
    referenceId: options.referenceId,
    provider: options.provider,
    model: options.model,
    credentialStatus: options.credentialStatus,
    modelList:
      options.modelList ??
      createModelSettingsModelListState({
        provider: options.provider,
        credentialStatus: options.credentialStatus
      }),
    saveStatus: options.saveStatus ?? "clean",
    saveMessage: options.saveMessage ?? ""
  }
}

export function cloneModelSettingsCardState(
  state: ModelSettingsCardState
): ModelSettingsCardState {
  return {
    ...state,
    modelList: {
      ...state.modelList,
      models: state.modelList.models.map((model) => ({ ...model }))
    }
  }
}

export function updateModelSettingsProvider(
  state: ModelSettingsCardState,
  options: {
    provider: string
    credentialStatus: ModelSettingsCredentialStatus
  }
): ModelSettingsCardState {
  const provider = options.provider.trim()
  return {
    ...state,
    provider,
    model: "",
    credentialStatus: options.credentialStatus,
    modelList: createModelSettingsModelListState({
      provider,
      credentialStatus: options.credentialStatus
    }),
    saveStatus: "dirty",
    saveMessage: ""
  }
}

export function startModelSettingsListRefresh(
  state: ModelSettingsCardState
): ModelSettingsCardState {
  return {
    ...state,
    modelList: createModelSettingsModelListState({
      provider: state.provider,
      credentialStatus: state.credentialStatus,
      status: "loading"
    }),
    saveStatus: state.saveStatus === "saved" ? "dirty" : state.saveStatus,
    saveMessage: ""
  }
}

export function applyModelSettingsListResult(
  state: ModelSettingsCardState,
  response: ModelSettingsModelListState
): ModelSettingsCardState {
  if (response.provider !== state.provider) {
    return state
  }

  const models = response.models.map((model) => ({ ...model }))
  const keepsCurrentModel = models.some(
    (model) => model.modelId === state.model
  )
  const nextModel = keepsCurrentModel ? state.model : ""

  return {
    ...state,
    model: isUsableModelList(response.status) ? nextModel : "",
    modelList: {
      ...response,
      models
    },
    saveStatus: "dirty",
    saveMessage: ""
  }
}

export function failModelSettingsListRefresh(
  state: ModelSettingsCardState,
  failureKind = "provider_unreachable"
): ModelSettingsCardState {
  return {
    ...state,
    model: "",
    modelList: createModelSettingsModelListState({
      provider: state.provider,
      credentialStatus: state.credentialStatus,
      status: "failed",
      models: [],
      failureKind
    }),
    saveStatus: "dirty",
    saveMessage: ""
  }
}

export function selectModelSettingsModel(
  state: ModelSettingsCardState,
  model: string
): ModelSettingsCardState {
  const nextModel = model.trim()
  if (!isUsableModelList(state.modelList.status)) {
    return state
  }

  const isKnownModel = state.modelList.models.some(
    (option) => option.modelId === nextModel
  )
  if (nextModel !== "" && !isKnownModel) {
    return state
  }

  return {
    ...state,
    model: nextModel,
    saveStatus: "dirty",
    saveMessage: ""
  }
}

export function markModelSettingsSaving(
  state: ModelSettingsCardState
): ModelSettingsCardState {
  return {
    ...state,
    saveStatus: "saving",
    saveMessage: ""
  }
}

export function markModelSettingsSaved(
  state: ModelSettingsCardState,
  options: {
    provider: string
    model: string
    models?: ModelSettingsModelOption[]
    message: string
  }
): ModelSettingsCardState {
  const provider = options.provider.trim()
  const model = options.model.trim()
  return {
    ...state,
    provider,
    model,
    modelList: createModelSettingsModelListState({
      provider,
      credentialStatus: state.credentialStatus,
      status: "success",
      models:
        options.models ?? (model ? [{ modelId: model, label: model }] : [])
    }),
    saveStatus: "saved",
    saveMessage: options.message
  }
}

export function markModelSettingsSaveFailed(
  state: ModelSettingsCardState,
  message: string
): ModelSettingsCardState {
  return {
    ...state,
    saveStatus: "failed",
    saveMessage: message
  }
}

function buildCredentialStatusLabel(state: ModelSettingsCardState): string {
  if (state.credentialStatus === "configured") {
    return "設定済み"
  }

  if (state.credentialStatus === "not_required") {
    return "不要"
  }

  return "APIキーが未設定です"
}

function buildModelListStatusText(state: ModelSettingsCardState): string {
  if (state.provider === "") {
    return "AIサービスを選んでください。"
  }

  switch (state.modelList.status) {
    case "loading":
      return "モデル一覧を更新しています。"
    case "success":
    case "credential_not_required":
      return state.modelList.models.length === 0
        ? "モデル一覧を取得しました。候補は 0 件です。"
        : "モデル一覧を取得しました。"
    case "credential_missing":
      return "APIキーが未設定のため、モデル一覧を更新できません。"
    case "failed":
      return "モデル一覧を取得できませんでした。もう一度更新してください。"
    case "not_updated":
    default:
      return "モデル一覧を更新してください。"
  }
}

function buildStatus(state: ModelSettingsCardState): {
  label: string
  tone: "neutral" | "warning" | "success"
  helper: string
} {
  if (state.provider === "") {
    return {
      label: "未選択",
      tone: "warning",
      helper: "AIサービスを選んでください。"
    }
  }

  if (state.modelList.status === "loading") {
    return {
      label: "更新中",
      tone: "neutral",
      helper: "モデル一覧を更新しています。"
    }
  }

  if (state.modelList.status === "failed") {
    return {
      label: "取得失敗",
      tone: "warning",
      helper: "モデル一覧を取得できませんでした。"
    }
  }

  if (
    (state.modelList.status === "success" ||
      state.modelList.status === "credential_not_required") &&
    state.modelList.models.length === 0
  ) {
    return {
      label: "候補 0 件",
      tone: "warning",
      helper: "モデル一覧を取得しましたが、候補は 0 件です。"
    }
  }

  if (state.model === "") {
    return {
      label: "モデル未選択",
      tone: "warning",
      helper: "使うモデルを選んでください。"
    }
  }

  if (state.saveStatus === "failed") {
    return {
      label: "保存失敗",
      tone: "warning",
      helper: "未保存の変更があります。もう一度保存してください。"
    }
  }

  if (state.saveStatus === "dirty") {
    return {
      label: "未保存",
      tone: "warning",
      helper: "未保存の変更があります。"
    }
  }

  if (state.saveStatus === "saving") {
    return {
      label: "保存中",
      tone: "neutral",
      helper: "モデル設定を保存しています。"
    }
  }

  return {
    label: state.saveStatus === "saved" ? "保存済み" : "選択済み",
    tone: "success",
    helper: "使うモデルを選択済みです。"
  }
}

export function buildModelSettingsCardViewModel(options: {
  state: ModelSettingsCardState
  providerOptions: ModelSettingsProviderOption[]
  refreshDisabled?: boolean
  actionDisabled?: boolean
  titleLabel?: string
}): ModelSettingsCardViewModel {
  const state = options.state
  const status = buildStatus(state)
  const modelListUsable = isUsableModelList(state.modelList.status)

  return {
    referenceId: state.referenceId,
    provider: state.provider,
    model: state.model,
    providerOptions: options.providerOptions
      .filter((provider) => USER_VISIBLE_PROVIDERS.has(provider.value))
      .map((provider) => ({
        value: provider.value,
        label: provider.label || formatProviderLabel(provider.value)
      })),
    credentialStatusLabel: buildCredentialStatusLabel(state),
    credentialStatusTone:
      state.credentialStatus === "configured"
        ? "success"
        : state.credentialStatus === "not_required"
          ? "neutral"
          : "warning",
    showCredentialStatus: true,
    showCredentialWarning: state.credentialStatus === "missing",
    credentialWarningText:
      "APIキーが未設定です。モデル一覧の取得結果はサービス設定に従います。",
    modelListButtonEnabled:
      !options.refreshDisabled &&
      state.provider !== "" &&
      state.modelList.status !== "loading",
    modelListButtonLabel: "モデル一覧を更新",
    modelListButtonAriaLabel: `${options.titleLabel ?? "モデル設定"}のモデル一覧を更新`,
    isModelListRefreshing: state.modelList.status === "loading",
    modelListStatusText: buildModelListStatusText(state),
    modelOptions: modelListUsable ? state.modelList.models : [],
    modelSelectEnabled:
      modelListUsable &&
      state.modelList.models.length > 0 &&
      state.credentialStatus !== "missing",
    emptyModelLabel: modelListUsable
      ? "選んでください"
      : "モデル一覧を更新してください",
    statusLabel: status.label,
    statusTone: status.tone,
    helperText: status.helper,
    footerMessage: state.saveMessage,
    footerWarningText: state.saveStatus === "failed" ? status.helper : "",
    actionButtonDisabled:
      options.actionDisabled === true ||
      state.provider === "" ||
      state.model === "" ||
      state.saveStatus === "saving" ||
      state.credentialStatus === "missing"
  }
}
