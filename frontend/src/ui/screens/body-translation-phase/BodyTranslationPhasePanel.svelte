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
  import { selectPhaseProgressActions } from "../../components/phase-progress-actions"
  import type {
    PhaseDetailItem,
    PhaseMetricCounter,
    PhaseStateToken
  } from "../../components/phase-panel-types"
  import type {
    ProcessingTargetListItem,
    ProcessingTargetListPageState
  } from "../../components/processing-target-list-panel-types"

  interface Props {
    viewModel: BodyTranslationPhaseScreenViewModel
    onAction: (actionId: BodyTranslationPhaseActionKind) => void | Promise<void>
    processingTargetItems?: ProcessingTargetListItem[]
    processingTargetPageState?: ProcessingTargetListPageState
    initialFetchDone?: boolean
    onProcessingTargetSearchInput?: (event: Event) => void
    onProcessingTargetPreviousPage?: () => void
    onProcessingTargetNextPage?: () => void
    onProcessingTargetPageChange?: (page: number) => void
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
    processingTargetPageState = undefined,
    initialFetchDone = true,
    onProcessingTargetSearchInput = undefined,
    onProcessingTargetPreviousPage = undefined,
    onProcessingTargetNextPage = undefined,
    onProcessingTargetPageChange = undefined,
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
    { label: "AI 送信対象", value: viewModel.providerTargetCountLabel },
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
    { label: "AI 送信対象件数", value: viewModel.providerTargetCountLabel }
  ])
  const phaseActionCards = $derived(
    selectPhaseProgressActions<BodyTranslationPhaseActionKind>(
      viewModel.actionCards,
      {
        hiddenActionIds: ["check-output-readiness"],
        runActionIds: ["start", "resume", "retry"],
        alwaysActionIds: ["pause"]
      }
    )
  )
  const displayedProcessingTargetItems = $derived(
    processingTargetPageState
      ? processingTargetPageState.items
      : (providedProcessingTargetItems ?? [])
  )
  const processingTargetSearchQuery = $derived(
    processingTargetPageState?.searchQuery ?? processingTargetSearchValue
  )
  const filteredProcessingTargetItems = $derived(
    processingTargetPageState
      ? displayedProcessingTargetItems
      : filterProcessingTargetItems(
          displayedProcessingTargetItems,
          processingTargetSearchQuery
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
    onProcessingTargetSearchInput?.(event)
    const target = event.currentTarget
    if (target instanceof HTMLInputElement) {
      processingTargetSearchValue = target.value
    }
  }
</script>

<section
  class="job-run-shell"
  data-testid="body-translation-phase-screen"
  id="bodyTranslationPhaseView"
>
  {#if !initialFetchDone}
    <div
      class="phase-loading-overlay"
      data-testid="body-translation-phase-processing-target-loading"
      aria-label="処理対象一覧を取得中"
      aria-busy="true"
    >
      <span class="phase-loading-spinner" aria-hidden="true"></span>
      <span class="loading-text">処理対象一覧を取得中...</span>
    </div>
  {/if}
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
    testId="body-translation-phase-status-panel"
    metrics={phaseMetrics}
  />

  <section class="phase-controls-grid">
    <PhaseProgressPanel
      headingId="bodyPhaseProgressHeading"
      testId="body-translation-phase-progress"
      eyebrow="翻訳段階の進行状況"
      title="進行状況"
      progressLabel={viewModel.progressLabel}
      progressBarTestId="body-translation-phase-progress-bar"
      progressCountsTestId="body-translation-phase-progress-counts"
      progressPercent={viewModel.progressPercent}
      progressDetail={viewModel.progressDetail}
      details={progressDetails}
      currentPhaseLabel={viewModel.currentPhaseLabel}
      actionAriaLabel="翻訳段階の操作"
      actions={phaseActionCards}
      startButtonTestId="body-translation-phase-start-button"
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
      statusTestId="body-translation-phase-ai-model-lock-state"
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

  <div class="processing-target-container">
    <ProcessingTargetListWrapper
      items={filteredProcessingTargetItems}
      pageSize={processingTargetPageState?.pageSize ?? 50}
      currentPage={processingTargetPageState?.page}
      totalCount={processingTargetPageState?.totalCount}
      busy={processingTargetPageState?.busy}
      searchId="bodyPhaseProcessingTargetSearch"
      searchTestId="body-translation-phase-processing-target-search-input"
      searchLabel="検索"
      searchPlaceholder="名前・原文・訳語で検索"
      searchValue={processingTargetSearchQuery}
      title="処理対象"
      titleId="bodyPhaseProcessingTargetsHeading"
      titleColumnLabels={["原文", "訳文"]}
      onSearchInput={!initialFetchDone ? undefined : handleProcessingTargetSearchInput}
      onPreviousPage={!initialFetchDone ? undefined : onProcessingTargetPreviousPage}
      onNextPage={!initialFetchDone ? undefined : onProcessingTargetNextPage}
      onPageChange={!initialFetchDone ? undefined : onProcessingTargetPageChange}
      rowTestId="body-translation-phase-processing-target-row"
      totalCountTestId="body-translation-phase-processing-target-total"
      emptyStateTestId="body-translation-phase-processing-target-empty"
    />
  </div>
</section>

<style>
  .job-run-shell {
    position: relative;
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

  .phase-loading-overlay {
    position: absolute;
    inset: 0;
    z-index: 20;
    display: flex;
    flex-direction: column;
    align-items: center;
    justify-content: center;
    gap: 0.75rem;
    background: rgba(18, 16, 15, 0.72);
    backdrop-filter: blur(2px);
  }

  .phase-loading-spinner {
    width: 2.25rem;
    height: 2.25rem;
    border-radius: 50%;
    border: 3px solid rgba(236, 223, 205, 0.25);
    border-top-color: #ff9c7c;
    animation: phase-loading-spin 0.8s linear infinite;
  }

  .loading-text {
    font-size: 0.875rem;
    color: rgba(236, 223, 205, 0.78);
  }

  @keyframes phase-loading-spin {
    to {
      transform: rotate(360deg);
    }
  }

  @media (prefers-reduced-motion: reduce) {
    .phase-loading-spinner {
      animation-duration: 1.6s;
    }
  }
</style>
