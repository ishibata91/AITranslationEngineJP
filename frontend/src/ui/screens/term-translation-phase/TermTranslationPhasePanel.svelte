<script lang="ts">
  import AIModelSelectionCard from "../../components/AIModelSelectionCard.svelte"
  import PhaseProgressPanel from "../../components/PhaseProgressPanel.svelte"
  import PhaseStatusPanel from "../../components/PhaseStatusPanel.svelte"
  import type {
    PhaseDetailItem,
    PhaseMetricCounter
  } from "../../components/phase-panel-types"
  import type {
    TermTranslationPhaseActionKind,
    TermTranslationPhaseScreenViewModel
  } from "@application/contract/term-translation-phase"

  type TermPanelActionKind = TermTranslationPhaseActionKind | "next-phase"

  interface Props {
    viewModel: TermTranslationPhaseScreenViewModel
    onAction: (actionId: TermPanelActionKind) => void | Promise<void>
    onAISettingsChange?: (request: {
      provider: string
      model: string
      executionMode: string
      batchMode: string
    }) => void | Promise<void>
  }

  let { viewModel, onAction, onAISettingsChange = undefined }: Props = $props()
  const phaseActionCards = $derived(
    viewModel.actionCards.filter(
      (
        action
      ): action is Extract<
        typeof action,
        { id: TermTranslationPhaseActionKind }
      > => action.id !== "next-phase"
    )
  )
  const canEditAiSettings = $derived(viewModel.viewState !== "running")
  const aiSettingsBlockedReason = $derived(
    viewModel.credentialRefLabel === "-"
      ? "認証状態を確認してください。"
      : viewModel.modelLabel === "-"
        ? "モデルを選択してください。"
        : ""
  )
  const aiSettingsStatusLabel = $derived(
    aiSettingsBlockedReason ? "設定未完了" : "固定済み"
  )
  const aiSettingsStatusTone = $derived(
    aiSettingsBlockedReason ? "warning" : "success"
  )
  const aiModelOptions = $derived([
    { modelId: viewModel.modelLabel, label: viewModel.modelLabel }
  ])
  const phaseMetrics = $derived<PhaseMetricCounter[]>([
    { label: "対象", value: viewModel.totalTermCountLabel },
    { label: "処理済み", value: viewModel.confirmedCountLabel },
    { label: "成功", value: viewModel.confirmedCountLabel },
    { label: "失敗", value: "0" },
    { label: "スキップ", value: "0" }
  ])
  const progressDetails = $derived<PhaseDetailItem[]>([
    { label: "対象語件数", value: viewModel.totalTermCountLabel }
  ])

  function selectedValue(event: Event): string {
    const target = event.currentTarget
    return target instanceof HTMLSelectElement ? target.value : ""
  }

  function saveAISettings(
    overrides: Partial<{
      provider: string
      model: string
      executionMode: string
      batchMode: string
    }>
  ): void {
    void onAISettingsChange?.({
      provider: viewModel.providerLabel === "-" ? "" : viewModel.providerLabel,
      model: viewModel.modelLabel === "-" ? "" : viewModel.modelLabel,
      executionMode:
        viewModel.executionModeLabel === "-"
          ? ""
          : viewModel.executionModeLabel,
      batchMode: "disabled",
      ...overrides
    })
  }

  function handleProviderChange(event: Event): void {
    saveAISettings({ provider: selectedValue(event) })
  }

  function handleExecutionChange(event: Event): void {
    saveAISettings({ executionMode: selectedValue(event) })
  }

  function handleModelChange(event: Event): void {
    saveAISettings({ model: selectedValue(event) })
  }
</script>

<section class="job-run-shell" id="termTranslationPhaseView">
  <PhaseStatusPanel
    eyebrow="translation-management"
    title="単語翻訳"
    gatewayStatus={viewModel.gatewayStatus}
    lead="現在の翻訳段階、進行状況、AI 設定を同じ画面で確認し、開始、中断、再開、リトライを判断します。"
    state={viewModel.viewState}
    stateLabel={viewModel.phaseStateLabel}
    statusTitle={viewModel.statusTitle}
    statusText={viewModel.statusText}
    errorMessage={viewModel.errorMessage}
    testId="term-translation-phase-screen-status-header"
    metrics={phaseMetrics}
  />

  <PhaseProgressPanel
    headingId="termPhaseProgressHeading"
    testId="term-translation-phase-progress-region"
    eyebrow="翻訳段階の進行状況"
    title="進行状況"
    progressLabel={viewModel.progressLabel}
    progressPercent={viewModel.progressPercent}
    progressDetail={viewModel.progressDetail}
    details={progressDetails}
    currentPhaseLabel={viewModel.currentPhaseLabel}
    actionAriaLabel="翻訳段階の操作"
    actions={phaseActionCards}
    onAction={(actionId: TermTranslationPhaseActionKind) => onAction(actionId)}
  />

  <section class="summary-grid ai-settings-row">
    <AIModelSelectionCard
      dataTestId="term-translation-phase-ai-model-selection-region"
      ariaLabel="単語翻訳の AI モデル選択"
      eyebrow="単語翻訳"
      title="単語翻訳の AI モデル"
      titleId="termPhaseAiModelHeading"
      helperText="単語翻訳を開始する前に使う AI サービス、モデル、処理方式を確認します。"
      statusLabel={aiSettingsStatusLabel}
      statusTone={aiSettingsStatusTone}
      providerSelectId="termPhaseProviderSelect"
      providerValue={viewModel.providerLabel}
      providerOptions={[
        { value: viewModel.providerLabel, label: viewModel.providerLabel }
      ]}
      providerDisabled={!canEditAiSettings}
      onProviderChange={handleProviderChange}
      credentialStatusLabel={viewModel.credentialRefLabel}
      credentialStatusTone={aiSettingsBlockedReason ? "warning" : "success"}
      showCredentialWarning={Boolean(aiSettingsBlockedReason)}
      credentialWarningText={aiSettingsBlockedReason}
      secondaryControlMode="execution-select"
      executionSelectId="termPhaseExecutionModeSelect"
      executionValue={viewModel.executionModeLabel}
      executionOptions={[
        {
          value: viewModel.executionModeLabel,
          label: viewModel.executionModeLabel
        }
      ]}
      executionDisabled={!canEditAiSettings}
      onExecutionChange={handleExecutionChange}
      modelSelectId="termPhaseModelSelect"
      modelValue={viewModel.modelLabel}
      modelOptions={aiModelOptions}
      modelDisabled={!canEditAiSettings}
      onModelChange={handleModelChange}
      modelStatusText="モデル一覧は単語翻訳の開始前に更新します。"
      refreshDisabled={!canEditAiSettings}
      onRefresh={() => saveAISettings({})}
      footerMessage={canEditAiSettings
        ? `一括処理: ${viewModel.providerSkippedLabel}。設定は単語翻訳の開始時に固定します。`
        : "実行中は AI 設定を編集できません。"}
      footerWarningText={aiSettingsBlockedReason}
    />
  </section>
</section>

<style>
  .job-run-shell {
    display: grid;
    gap: 1.25rem;
  }

  .summary-grid {
    display: grid;
    gap: 1.25rem;
    grid-template-columns: repeat(2, minmax(0, 1fr));
  }

  .ai-settings-row {
    grid-template-columns: minmax(0, 1fr);
  }

  @media (max-width: 900px) {
    .summary-grid {
      grid-template-columns: 1fr;
    }
  }
</style>
