import * as MasterPersonaGateway from "@application/gateway-contract/master-persona"
import { buildModelSettingsCardViewModel } from "@application/gateway-contract/model-settings-card"

type MasterPersonaScreenState = MasterPersonaGateway.MasterPersonaScreenState
type MasterPersonaScreenViewModel =
  MasterPersonaGateway.MasterPersonaScreenViewModel
type ModelSettingsCardState = NonNullable<
  MasterPersonaScreenState["modelSettingsCard"]
>

const AI_PROVIDER_LABEL_BY_ID: Record<string, string> = {
  gemini: "Gemini",
  lm_studio: "LM Studio",
  xai: "xAI"
}

const MASTER_PERSONA_PROVIDER_OPTIONS = [
  { value: "gemini", label: "Gemini", credentialStatus: "configured" as const },
  {
    value: "lm_studio",
    label: "LM Studio",
    credentialStatus: "not_required" as const
  },
  { value: "xai", label: "xAI", credentialStatus: "configured" as const }
]

function credentialStatusForProvider(
  state: MasterPersonaScreenState
): ModelSettingsCardState["credentialStatus"] {
  const provider = state.aiSettings.provider.trim()
  const found = state.providerOptions.find(
    (option) => option.value === provider
  )
  if (found) {
    return found.credentialStatus
  }
  return provider === "lm_studio" ? "not_required" : "missing"
}

function createFallbackModelSettingsCard(
  state: MasterPersonaScreenState
): ModelSettingsCardState {
  const credentialStatus = credentialStatusForProvider(state)
  return {
    referenceId: "master-persona",
    provider: state.aiSettings.provider,
    model: state.aiSettings.model,
    credentialStatus,
    modelList: {
      provider: state.aiSettings.provider,
      credentialStatus,
      status: state.modelOptions.length > 0 ? "success" : "not_updated",
      models: state.modelOptions.map((model) => ({ ...model }))
    },
    saveStatus: "clean",
    saveMessage: ""
  }
}

function buildPluginOptions(state: MasterPersonaScreenState): Array<{
  value: string
  label: string
}> {
  const options = state.pluginGroups.map((group) => ({
    value: group.targetPlugin,
    label: `${group.targetPlugin} (${group.count})`
  }))

  return [{ value: "", label: "すべてのプラグイン" }, ...options]
}

function buildPageStatusText(state: MasterPersonaScreenState): string {
  if (state.totalCount === 0) {
    return "1 - 0 件を表示しています。"
  }

  const start = (state.page - 1) * state.pageSize + 1
  const end = Math.min(state.page * state.pageSize, state.totalCount)
  return `${start} - ${end} 件を表示しています。`
}

function buildSelectionStatusText(state: MasterPersonaScreenState): string {
  if (!state.selectedEntry) {
    return "選択中のペルソナはありません。"
  }
  return `${state.selectedEntry.displayName} を選択中`
}

function buildDetailStatusText(state: MasterPersonaScreenState): string {
  if (!state.selectedEntry) {
    return "一覧からペルソナを選ぶと、詳細を同じ画面で確認できます。"
  }
  return state.selectedEntry.runLockReason
}

function isRunActive(runState: string): boolean {
  return runState === "生成中"
}

function normalizeProviderId(provider: string): string {
  return provider.trim().toLowerCase()
}

function buildAIProviderLabel(provider: string): string {
  const providerId = normalizeProviderId(provider)
  if (providerId === "") {
    return ""
  }
  return AI_PROVIDER_LABEL_BY_ID[providerId] ?? provider.trim()
}

function buildAISettingsWarningText(state: MasterPersonaScreenState): string {
  const providerId = normalizeProviderId(state.aiSettings.provider)
  if (providerId === "") {
    return "AIサービスを選んでください。"
  }

  if (state.modelOptions.length === 0) {
    return "モデル一覧を更新後に選べる状態で接続します。"
  }

  return ""
}

function isAISettingsComplete(state: MasterPersonaScreenState): boolean {
  const providerId = normalizeProviderId(state.aiSettings.provider)
  const hasProvider = providerId !== ""
  const hasModel = state.aiSettings.model.trim() !== ""

  return hasProvider && hasModel
}

function buildExecutionMethodOptions(provider: string) {
  const providerId = normalizeProviderId(provider)
  if (providerId === "gemini" || providerId === "xai") {
    return [
      { value: "single_request", label: "通常" },
      { value: "batch", label: "Batch API" }
    ]
  }

  return [{ value: "single_request", label: "通常" }]
}

function buildProgressPercent(state: MasterPersonaScreenState): number {
  const processed = state.runStatus.processedCount
  const total = processed + state.runStatus.existingSkipCount
  if (total <= 0) {
    return state.runStatus.runState === "完了" ? 100 : 0
  }
  return Math.max(0, Math.min(100, Math.round((processed / total) * 100)))
}

export class MasterPersonaPresenter {
  toViewModel(
    state: MasterPersonaScreenState,
    isGatewayConnected: boolean
  ): MasterPersonaScreenViewModel {
    const activeRun = isRunActive(state.runStatus.runState)
    const totalPages = Math.max(
      1,
      Math.ceil(
        state.totalCount / MasterPersonaGateway.MASTER_PERSONA_PAGE_SIZE
      )
    )
    const hasPreview = state.preview !== null
    const aiSettingsWarningText = buildAISettingsWarningText(state)
    const executionMethodOptions = buildExecutionMethodOptions(
      state.aiSettings.provider
    )
    const modelSettingsCard =
      state.modelSettingsCard ?? createFallbackModelSettingsCard(state)

    return {
      ...state,
      gatewayStatus: isGatewayConnected ? "接続準備済み" : "未接続",
      pluginOptions: buildPluginOptions(state),
      totalPages,
      pageStatusText: buildPageStatusText(state),
      selectionStatusText: buildSelectionStatusText(state),
      listHeadline: `${state.totalCount.toLocaleString("ja-JP")} 件から絞り込みます。`,
      detailLockText: activeRun
        ? "更新と削除を行えません"
        : "更新と削除を行えます",
      detailStatusText: buildDetailStatusText(state),
      canStartPreview: state.selectedFileReference !== null,
      canStartGeneration:
        isAISettingsComplete(state) &&
        state.preview !== null &&
        state.preview.status === "生成可能" &&
        !activeRun,
      canMutate: !activeRun && state.selectedEntry !== null,
      isRunActive: activeRun,
      hasPreview,
      aiProviderLabel: buildAIProviderLabel(state.aiSettings.provider),
      aiSettingsWarningText,
      aiSettingsStatusText:
        aiSettingsWarningText === "" ? "設定済み" : "設定が必要",
      modelSettingsCardViewModel: buildModelSettingsCardViewModel({
        state: modelSettingsCard,
        providerOptions:
          state.providerOptions.length > 0
            ? state.providerOptions
            : MASTER_PERSONA_PROVIDER_OPTIONS,
        actionDisabled: activeRun,
        titleLabel: "マスターペルソナ"
      }),
      canSelectModel:
        state.modelOptions.length > 0 && aiSettingsWarningText === "",
      executionMethodOptions,
      promptTemplateDescription:
        MasterPersonaGateway.MASTER_PERSONA_PROMPT_TEMPLATE_DESCRIPTION,
      progressPercent: buildProgressPercent(state),
      page: Math.min(state.page, totalPages)
    }
  }
}
