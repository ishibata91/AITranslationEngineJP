<script lang="ts">
  import type {
    BodyTranslationPhaseActionKind,
    BodyTranslationPhaseScreenViewModel,
    BodyTranslationPhaseViewState
  } from "@application/contract/body-translation-phase"
  import AIModelSelectionCard from "../../components/AIModelSelectionCard.svelte"
  import PhaseProgressPanel from "../../components/PhaseProgressPanel.svelte"
  import PhaseStatusPanel from "../../components/PhaseStatusPanel.svelte"
  import type {
    PhaseDetailItem,
    PhaseMetricCounter,
    PhaseStateToken
  } from "../../components/phase-panel-types"

  interface Props {
    viewModel: BodyTranslationPhaseScreenViewModel
    onAction: (actionId: BodyTranslationPhaseActionKind) => void | Promise<void>
    onAISettingsChange?: (request: {
      provider: string
      model: string
      executionMode: string
      batchMode: string
    }) => void | Promise<void>
  }

  let { viewModel, onAction, onAISettingsChange = undefined }: Props = $props()

  function resolveStateToken(
    viewState: BodyTranslationPhaseViewState
  ): PhaseStateToken {
    if (viewState === "not_ready" || viewState === "validation_failed") {
      return "blocked"
    }
    return viewState
  }

  const phaseMetrics = $derived<PhaseMetricCounter[]>([
    { label: "対象", value: viewModel.targetCountLabel },
    { label: "処理済み", value: viewModel.processedCountLabel },
    { label: "成功", value: viewModel.translatedCountLabel },
    { label: "失敗", value: viewModel.failedCountLabel },
    { label: "スキップ", value: viewModel.skippedCountLabel }
  ])
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

  const progressDetails = $derived<PhaseDetailItem[]>([
    { label: "対象件数", value: viewModel.targetCountLabel }
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

<section class="job-run-shell" id="bodyTranslationPhaseView">
  <PhaseStatusPanel
    title="本文翻訳"
    eyebrow="translation-management"
    gatewayStatus={viewModel.gatewayStatus}
    lead="現在の翻訳段階、進行状況、AI 設定を同じ画面で確認し、開始、中断、再開、回復を判断します。"
    stateLabel={viewModel.phaseStateLabel}
    state={resolveStateToken(viewModel.viewState)}
    statusTitle={viewModel.statusTitle}
    statusText={viewModel.statusText}
    errorMessage={viewModel.errorMessage}
    testId="body-translation-phase-body-translation-summary"
    metrics={phaseMetrics}
  />

  <PhaseProgressPanel
    headingId="bodyPhaseProgressHeading"
    testId="body-translation-phase-progress"
    eyebrow="翻訳段階の進行状況"
    title="進行状況"
    progressLabel={viewModel.progressLabel}
    progressPercent={viewModel.progressPercent}
    progressDetail={viewModel.progressDetail}
    details={progressDetails}
    currentPhaseLabel={viewModel.currentPhaseLabel}
    actionAriaLabel="翻訳段階の操作"
    actions={viewModel.actionCards}
    {onAction}
  />

  <section class="summary-grid ai-settings-row">
    <AIModelSelectionCard
      dataTestId="body-translation-phase-ai-model-selection"
      ariaLabel="本文翻訳の AI モデル選択"
      eyebrow="本文翻訳"
      title="本文翻訳の AI モデル"
      titleId="bodyPhaseAiModelHeading"
      helperText="本文翻訳を開始する前に使う AI サービス、モデル、処理方式を確認します。"
      statusLabel={aiSettingsStatusLabel}
      statusTone={aiSettingsStatusTone}
      providerSelectId="bodyPhaseProviderSelect"
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
      executionSelectId="bodyPhaseExecutionModeSelect"
      executionValue={viewModel.executionModeLabel}
      executionOptions={[
        {
          value: viewModel.executionModeLabel,
          label: viewModel.executionModeLabel
        }
      ]}
      executionDisabled={!canEditAiSettings}
      onExecutionChange={handleExecutionChange}
      modelSelectId="bodyPhaseModelSelect"
      modelValue={viewModel.modelLabel}
      modelOptions={aiModelOptions}
      modelDisabled={!canEditAiSettings}
      onModelChange={handleModelChange}
      modelStatusText="モデル一覧は本文翻訳の開始前に更新します。"
      refreshDisabled={!canEditAiSettings}
      onRefresh={() => saveAISettings({})}
      footerMessage={canEditAiSettings
        ? `一括処理: ${viewModel.providerStateLabel}。設定は本文翻訳の開始時に固定します。`
        : "実行中は AI 設定を編集できません。"}
      footerWarningText={aiSettingsBlockedReason}
    />
  </section>
</section>

<style>
  .job-run-shell {
    display: grid;
    gap: 1.25rem;
    min-width: 0;
  }

  .summary-grid {
    display: grid;
    gap: 0.8rem;
    min-width: 0;
  }

  .summary-grid {
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
