<script lang="ts">
  import type {
    MasterPersonaAISettings,
    MasterPersonaModelOption
  } from "@application/gateway-contract/master-persona/master-persona-gateway-contract"
  import type { ModelSettingsCardViewModel } from "@application/gateway-contract/model-settings-card"
  import AIModelSelectionCard from "@ui/components/AIModelSelectionCard.svelte"

  interface Props {
    aiSettings: MasterPersonaAISettings
    aiSettingsStatusText: string
    aiSettingsWarningText: string
    aiProviderLabel: string
    canSelectModel: boolean
    executionMethodOptions: Array<{ value: string; label: string }>
    modelOptions: MasterPersonaModelOption[]
    modelSettingsCardViewModel?: ModelSettingsCardViewModel
    isAISettingsRefreshing: boolean
    handleAIProviderChange: (event: Event) => void
    handleAIModelChange: (event: Event) => void
    handleAIExecutionMethodChange: (event: Event) => void
    refreshAISettings: () => Promise<void>
    saveAISettings: () => void
  }

  let {
    aiSettings,
    aiSettingsStatusText,
    aiSettingsWarningText,
    aiProviderLabel,
    canSelectModel,
    executionMethodOptions,
    modelOptions,
    modelSettingsCardViewModel,
    isAISettingsRefreshing,
    handleAIProviderChange,
    handleAIModelChange,
    handleAIExecutionMethodChange,
    refreshAISettings,
    saveAISettings
  }: Props = $props()

  function createFallbackModelCard(): ModelSettingsCardViewModel {
    return {
      referenceId: "master-persona",
      provider: aiSettings.provider,
      model: aiSettings.model,
      providerOptions: [
        { value: "gemini", label: "Gemini" },
        { value: "lm_studio", label: "LM Studio" },
        { value: "xai", label: "xAI" }
      ],
      credentialStatusLabel: aiSettingsStatusText,
      credentialStatusTone: aiSettingsWarningText ? "warning" : "success",
      showCredentialStatus: true,
      showCredentialWarning: aiSettingsWarningText !== "",
      credentialWarningText:
        "APIキーが未設定のため、モデル一覧を更新できません。",
      modelListButtonEnabled: true,
      modelListButtonLabel: "モデル一覧を更新",
      modelListButtonAriaLabel: "モデル一覧を更新",
      isModelListRefreshing: isAISettingsRefreshing,
      modelListStatusText: canSelectModel
        ? "使うモデルを選んでください。"
        : "モデル一覧を更新してください。",
      modelOptions,
      modelSelectEnabled: canSelectModel,
      emptyModelLabel: canSelectModel ? "選んでください" : "設定が必要",
      statusLabel: aiProviderLabel || "未選択",
      statusTone: aiSettingsWarningText ? "warning" : "neutral",
      helperText: "ペルソナ作成に使う AI サービスとモデルを選びます。",
      footerMessage: "",
      footerWarningText: "",
      actionButtonDisabled: isAISettingsRefreshing
    }
  }

  const modelCard = $derived(
    modelSettingsCardViewModel ?? createFallbackModelCard()
  )
</script>

<AIModelSelectionCard
  dataTestId="master-persona-ai-settings-card"
  actionButtonDisabled={modelCard.actionButtonDisabled ||
    isAISettingsRefreshing}
  actionButtonId="saveAiSettingsButton"
  actionButtonLabel="設定を保存"
  credentialStatusLabel={modelCard.credentialStatusLabel}
  credentialStatusTone={modelCard.credentialStatusTone}
  credentialWarningText={modelCard.credentialWarningText}
  eyebrow="AI 設定"
  executionDisabled={isAISettingsRefreshing}
  executionOptions={executionMethodOptions}
  executionSelectId="executionMethodSelect"
  executionValue={aiSettings.executionMethod}
  footerMessage={modelCard.footerMessage ||
    "モデル設定はこの画面専用です。必要なら保存できます。"}
  footerWarningText={modelCard.footerWarningText}
  helperText={modelCard.helperText}
  modelDisabled={isAISettingsRefreshing || !modelCard.modelSelectEnabled}
  modelOptions={modelCard.modelOptions}
  modelSelectId="modelInput"
  modelStatusText={modelCard.modelListStatusText}
  modelValue={modelCard.model}
  onAction={saveAISettings}
  onExecutionChange={handleAIExecutionMethodChange}
  onModelChange={handleAIModelChange}
  onProviderChange={handleAIProviderChange}
  onRefresh={() => void refreshAISettings()}
  providerFieldLabel="AI サービス"
  providerOptions={modelCard.providerOptions}
  providerDisabled={isAISettingsRefreshing}
  providerSelectId="providerSelect"
  providerValue={modelCard.provider}
  refreshButtonAriaLabel={modelCard.modelListButtonAriaLabel}
  refreshButtonLabel={modelCard.modelListButtonLabel}
  refreshDisabled={!modelCard.modelListButtonEnabled ||
    isAISettingsRefreshing}
  refreshSpinning={modelCard.isModelListRefreshing || isAISettingsRefreshing}
  secondaryControlMode="execution-select"
  showCredentialStatus={modelCard.showCredentialStatus}
  showCredentialWarning={modelCard.showCredentialWarning}
  statusLabel={modelCard.statusLabel}
  statusTone={modelCard.statusTone}
  title="この画面で使う AI 設定"
  titleId="settingsHeading"
  titleTag="h3"
  emptyModelLabel={modelCard.emptyModelLabel}
/>

<p class="sr-only">
  プロンプトテンプレートは画面入力では変更せず、実装側の説明文として固定しています。
</p>

<style>
  .sr-only {
    border: 0;
    clip: rect(0 0 0 0);
    clip-path: inset(50%);
    height: 1px;
    margin: -1px;
    overflow: hidden;
    padding: 0;
    position: absolute;
    white-space: nowrap;
    width: 1px;
  }
</style>
