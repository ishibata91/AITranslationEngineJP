<script lang="ts">
  import AIModelSelectionCard from "../../components/AIModelSelectionCard.svelte"
  import PhaseProgressPanel from "../../components/PhaseProgressPanel.svelte"
  import PhaseStatusPanel from "../../components/PhaseStatusPanel.svelte"
  import type {
    PhaseDetailItem,
    PhaseMetricCounter
  } from "../../components/phase-panel-types"
  import type {
    PersonaGenerationPhaseActionKind,
    PersonaGenerationPhaseScreenViewModel
  } from "@application/contract/persona-generation-phase"

  interface Props {
    viewModel: PersonaGenerationPhaseScreenViewModel
    onAction: (
      actionId: PersonaGenerationPhaseActionKind
    ) => void | Promise<void>
    onAISettingsChange?: (request: {
      provider: string
      model: string
      executionMode: string
      batchMode: string
    }) => void | Promise<void>
  }

  let { viewModel, onAction, onAISettingsChange = undefined }: Props = $props()
  const statusMetrics = $derived<PhaseMetricCounter[]>([
    { label: "対象", value: viewModel.targetCountLabel },
    { label: "処理済み", value: viewModel.generatedCountLabel },
    { label: "成功", value: viewModel.generatedCountLabel },
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

<section class="job-run-shell" id="personaGenerationPhaseView">
  <PhaseStatusPanel
    eyebrow="translation-management"
    title="NPC ペルソナ生成"
    gatewayStatus={viewModel.gatewayStatus}
    lead="現在の翻訳段階、進行状況、AI 設定を同じ画面で確認し、開始、中断、再開、リトライ、キャンセルを判断します。"
    state={viewModel.viewState}
    stateLabel={viewModel.phaseStateLabel}
    statusTitle={viewModel.statusTitle}
    statusText={viewModel.statusText}
    errorMessage={viewModel.errorMessage}
    testId="persona-generation-phase-persona-generation-phase-screen"
    statusTestId="persona-generation-phase-status-summary-card"
    metrics={statusMetrics}
  />

  <PhaseProgressPanel
    headingId="personaPhaseProgressHeading"
    testId="persona-generation-phase-progress-card"
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
      dataTestId="persona-generation-phase-ai-model-selection-card"
      ariaLabel="NPC ペルソナ生成の AI モデル選択"
      eyebrow="NPC ペルソナ生成"
      title="NPC ペルソナ生成の AI モデル"
      titleId="personaPhaseAiModelHeading"
      helperText="NPC ペルソナ生成を開始する前に使う AI サービス、モデル、処理方式を確認します。"
      statusLabel={aiSettingsStatusLabel}
      statusTone={aiSettingsStatusTone}
      providerSelectId="personaPhaseProviderSelect"
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
      executionSelectId="personaPhaseExecutionModeSelect"
      executionValue={viewModel.executionModeLabel}
      executionOptions={[
        {
          value: viewModel.executionModeLabel,
          label: viewModel.executionModeLabel
        }
      ]}
      executionDisabled={!canEditAiSettings}
      onExecutionChange={handleExecutionChange}
      modelSelectId="personaPhaseModelSelect"
      modelValue={viewModel.modelLabel}
      modelOptions={aiModelOptions}
      modelDisabled={!canEditAiSettings}
      onModelChange={handleModelChange}
      modelStatusText="モデル一覧は NPC ペルソナ生成の開始前に更新します。"
      refreshDisabled={!canEditAiSettings}
      onRefresh={() => saveAISettings({})}
      footerMessage={canEditAiSettings
        ? `一括処理: ${viewModel.outputCountLabel}。設定は NPC ペルソナ生成の開始時に固定します。`
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
