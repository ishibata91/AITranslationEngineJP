<script lang="ts">
  import type {
    BodyTranslationPhaseActionKind,
    BodyTranslationPhaseScreenViewModel,
    BodyTranslationPhaseViewState
  } from "@application/contract/body-translation-phase"
  import AIModelSelectionCard from "../../components/AIModelSelectionCard.svelte"
  import PhaseProgressPanel from "../../components/PhaseProgressPanel.svelte"
  import PhaseStatusPanel from "../../components/PhaseStatusPanel.svelte"
  import ProcessingTargetListWrapper from "../../components/ProcessingTargetListWrapper.svelte"
  import type {
    PhaseDetailItem,
    PhaseMetricCounter,
    PhaseStateToken
  } from "../../components/phase-panel-types"
  import type { ProcessingTargetListItem } from "../../components/processing-target-list-panel-types"

  interface Props {
    viewModel: BodyTranslationPhaseScreenViewModel
    onAction: (actionId: BodyTranslationPhaseActionKind) => void | Promise<void>
    processingTargetItems?: ProcessingTargetListItem[]
    onAISettingsChange?: (request: {
      provider: string
      model: string
      executionMode: string
      batchMode: string
    }) => void | Promise<void>
  }

  let {
    viewModel,
    onAction,
    processingTargetItems: providedProcessingTargetItems = undefined,
    onAISettingsChange = undefined
  }: Props = $props()
  let processingTargetSearchValue = $state("")

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
  const summaryProcessingTargetItems = $derived<ProcessingTargetListItem[]>([
    {
      id: "body-translation-provider-target",
      name: "原文 / 訳文",
      titleParts: [
        { text: `原文対象: ${viewModel.providerTargetCountLabel} 件` },
        { text: `訳文出力: ${viewModel.outputCountLabel} 件` }
      ],
      detail: "辞書とペルソナ参照情報を使って本文訳文を作る対象。",
      metadata: [
        { label: "本文翻訳対象", value: viewModel.targetCountLabel },
        { label: "AI 送信対象", value: viewModel.providerTargetCountLabel },
        {
          label: "完全一致辞書除外",
          value: viewModel.exactDictionaryExclusionCountLabel
        },
        {
          label: "部分一致辞書制約",
          value: viewModel.partialDictionaryConstraintCountLabel
        },
        { label: "リクエスト単位", value: viewModel.requestUnitCountLabel },
        { label: "出力済み", value: viewModel.outputCountLabel },
        { label: "参照辞書", value: viewModel.dictionaryDigestLabel },
        { label: "参照ペルソナ", value: viewModel.personaDigestLabel }
      ]
    }
  ])
  const displayedProcessingTargetItems = $derived(
    providedProcessingTargetItems && providedProcessingTargetItems.length > 0
      ? providedProcessingTargetItems
      : summaryProcessingTargetItems
  )
  const filteredProcessingTargetItems = $derived(
    filterProcessingTargetItems(
      displayedProcessingTargetItems,
      processingTargetSearchValue
    )
  )

  function filterProcessingTargetItems(
    items: ProcessingTargetListItem[],
    searchValue: string
  ): ProcessingTargetListItem[] {
    const normalizedSearchValue = searchValue.trim().toLocaleLowerCase("ja-JP")
    if (!normalizedSearchValue) {
      return items
    }

    return items.filter((item) =>
      [item.name, ...(item.titleParts?.map((part) => part.text) ?? [])]
        .join(" ")
        .toLocaleLowerCase("ja-JP")
        .includes(normalizedSearchValue)
    )
  }

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

  function handleProcessingTargetSearchInput(event: Event): void {
    const target = event.currentTarget
    if (target instanceof HTMLInputElement) {
      processingTargetSearchValue = target.value
    }
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

  <section class="phase-controls-grid">
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
      onAction={(actionId) =>
        onAction(actionId as BodyTranslationPhaseActionKind)}
    />

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

  <ProcessingTargetListWrapper
    items={filteredProcessingTargetItems}
    pageSize={50}
    searchId="bodyPhaseProcessingTargetSearch"
    searchLabel="検索"
    searchPlaceholder="名前・原文・訳語で検索"
    searchValue={processingTargetSearchValue}
    title="処理対象"
    titleId="bodyPhaseProcessingTargetsHeading"
    onSearchInput={handleProcessingTargetSearchInput}
  />
</section>

<style>
  .job-run-shell {
    display: grid;
    gap: 1.25rem;
    min-width: 0;
  }

  .phase-controls-grid {
    display: grid;
    gap: 1.25rem;
    grid-template-columns: repeat(2, minmax(0, 1fr));
    min-width: 0;
  }

  @media (max-width: 900px) {
    .phase-controls-grid {
      grid-template-columns: 1fr;
    }
  }
</style>
