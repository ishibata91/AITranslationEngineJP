<script lang="ts">
  import AIModelSelectionCard from "../../components/AIModelSelectionCard.svelte"
  import PhaseProgressPanel from "../../components/PhaseProgressPanel.svelte"
  import PhaseStatusPanel from "../../components/PhaseStatusPanel.svelte"
  import ProcessingTargetListWrapper from "../../components/ProcessingTargetListWrapper.svelte"
  import { selectPhaseProgressActions } from "../../components/phase-progress-actions"
  import type {
    PhaseDetailItem,
    PhaseMetricCounter
  } from "../../components/phase-panel-types"
  import type {
    ProcessingTargetListItem,
    ProcessingTargetListPageState
  } from "../../components/processing-target-list-panel-types"
  import type {
    TermTranslationPhaseActionKind,
    TermTranslationPhaseScreenViewModel
  } from "@application/contract/term-translation-phase"

  type TermPanelActionKind = TermTranslationPhaseActionKind | "next-phase"

  interface Props {
    viewModel: TermTranslationPhaseScreenViewModel
    onAction: (actionId: TermPanelActionKind) => void | Promise<void>
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
    onRefreshModelList?: (provider: string) => void | Promise<void>
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
    onAISettingsChange = undefined,
    onRefreshModelList = undefined
  }: Props = $props()
  let processingTargetSearchValue = $state("")
  let selectedProviderId = $state("")
  let selectedModelId = $state("")
  let selectedExecutionMode = $state("")
  const resolvedProviderValue = $derived(
    (viewModel.providerOptions ?? []).some((o) => o.value === selectedProviderId)
      ? selectedProviderId
      : ((viewModel.providerOptions ?? []).some((o) => o.value === viewModel.providerLabel)
          ? viewModel.providerLabel
          : ((viewModel.providerOptions ?? [])[0]?.value ?? ""))
  )
  const resolvedModelValue = $derived(
    selectedModelId ||
    ((viewModel.modelOptions ?? []).some((o) => o.value === viewModel.modelLabel)
      ? viewModel.modelLabel
      : "")
  )
  const resolvedExecutionMode = $derived(
    selectedExecutionMode ||
    ((viewModel.executionOptions ?? []).some((o) => o.value === viewModel.executionModeLabel)
      ? viewModel.executionModeLabel
      : ((viewModel.executionOptions ?? [])[0]?.value ?? ""))
  )
  const phaseActionCards = $derived(
    selectPhaseProgressActions<TermPanelActionKind>(viewModel.actionCards, {
      hiddenActionIds: ["next-phase"],
      runActionIds: ["start", "resume", "retry"],
      alwaysActionIds: ["pause"]
    })
  )
  const canEditAiSettings = $derived(viewModel.viewState !== "running")
  const aiSettingsBlockedReason = $derived(
    !viewModel.isExecutionConfigured
      ? "AI 設定が未完了です。"
      : ""
  )
  const aiSettingsStatusLabel = $derived(
    aiSettingsBlockedReason ? "設定未完了" : "固定済み"
  )
  const aiSettingsStatusTone = $derived(
    aiSettingsBlockedReason ? "warning" : "success"
  )
  const aiModelOptions = $derived(
    viewModel.modelOptions.map((o) => ({ modelId: o.value, label: o.label }))
  )
  const phaseMetrics = $derived<PhaseMetricCounter[]>([
    { label: "AI 翻訳対象語", value: viewModel.aiTargetCountLabel },
    { label: "処理済み", value: viewModel.confirmedCountLabel },
    { label: "成功", value: viewModel.confirmedCountLabel },
    { label: "失敗", value: "0" },
    { label: "スキップ", value: "0" }
  ])
  const progressDetails = $derived<PhaseDetailItem[]>([
    { label: "AI 翻訳対象語件数", value: viewModel.aiTargetCountLabel }
  ])
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
    const request = {
      provider: resolvedProviderValue,
      model: resolvedModelValue,
      executionMode: resolvedExecutionMode,
      batchMode: "disabled",
      ...overrides
    }
    if (!request.provider || !request.model || !request.executionMode) {
      return
    }
    void onAISettingsChange?.(request)
  }

  function handleProviderChange(event: Event): void {
    const provider = selectedValue(event)
    selectedProviderId = provider
    selectedModelId = ""
  }

  function handleExecutionChange(event: Event): void {
    const executionMode = selectedValue(event)
    selectedExecutionMode = executionMode
    saveAISettings({ executionMode })
  }

  function handleModelChange(event: Event): void {
    const model = selectedValue(event)
    selectedModelId = model
    saveAISettings({ model })
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
  data-testid="term-translation-phase-screen"
  id="termTranslationPhaseView"
>
  {#if !initialFetchDone}
    <div
      class="phase-loading-overlay"
      data-testid="term-translation-phase-processing-target-loading"
      aria-label="処理対象一覧を取得中"
      aria-busy="true"
    >
      <span class="phase-loading-spinner" aria-hidden="true"></span>
      <span class="loading-text">処理対象一覧を取得中...</span>
    </div>
  {/if}
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
    testId="term-translation-phase-status-panel"
    metrics={phaseMetrics}
  />

  <section class="phase-controls-grid">
    <PhaseProgressPanel
      headingId="termPhaseProgressHeading"
      testId="term-translation-phase-progress-region"
      eyebrow="翻訳段階の進行状況"
      title="進行状況"
      progressLabel={viewModel.progressLabel}
      progressBarTestId="term-translation-phase-progress-bar"
      progressCountsTestId="term-translation-phase-progress-counts"
      progressPercent={viewModel.progressPercent}
      progressDetail={viewModel.progressDetail}
      details={progressDetails}
      currentPhaseLabel={viewModel.currentPhaseLabel}
      actionAriaLabel="翻訳段階の操作"
      actions={phaseActionCards}
      startButtonTestId="term-translation-phase-start-button"
      actionHintsTestId="term-translation-phase-start-blocked-reason"
      onAction={(actionId) =>
        onAction(actionId as TermTranslationPhaseActionKind)}
    />

    <AIModelSelectionCard
      dataTestId="term-translation-phase-ai-model-selection-region"
      ariaLabel="単語翻訳の AI モデル選択"
      eyebrow="単語翻訳"
      title="単語翻訳の AI モデル"
      titleId="termPhaseAiModelHeading"
      helperText="単語翻訳を開始する前に使う AI サービス、モデル、処理方式を確認します。"
      statusLabel={aiSettingsStatusLabel}
      statusTestId="term-translation-phase-ai-settings-status-pill"
      statusTone={aiSettingsStatusTone}
      providerSelectId="termPhaseProviderSelect"
      providerSelectTestId="term-translation-phase-ai-provider-select"
      providerValue={resolvedProviderValue}
      providerOptions={viewModel.providerOptions}
      providerDisabled={!canEditAiSettings}
      onProviderChange={handleProviderChange}
      credentialStatusLabel={viewModel.credentialRefLabel}
      credentialStatusTone={aiSettingsBlockedReason ? "warning" : "success"}
      showCredentialWarning={Boolean(aiSettingsBlockedReason)}
      credentialWarningText={aiSettingsBlockedReason}
      secondaryControlMode="execution-select"
      executionSelectId="termPhaseExecutionModeSelect"
      executionSelectTestId="term-translation-phase-ai-execution-mode-select"
      executionValue={resolvedExecutionMode}
      executionOptions={viewModel.executionOptions}
      executionDisabled={!canEditAiSettings}
      onExecutionChange={handleExecutionChange}
      modelSelectId="termPhaseModelSelect"
      modelSelectTestId="term-translation-phase-ai-model-select"
      modelValue={resolvedModelValue}
      modelOptions={aiModelOptions}
      modelDisabled={!canEditAiSettings}
      onModelChange={handleModelChange}
      modelStatusText="モデル一覧は単語翻訳の開始前に更新します。"
      refreshDisabled={!canEditAiSettings}
      onRefresh={onRefreshModelList ? () => {
        void onRefreshModelList(resolvedProviderValue)
      } : undefined}
      footerMessage={canEditAiSettings
        ? `一括処理: ${viewModel.providerSkippedLabel}。設定は単語翻訳の開始時に固定します。`
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
      searchId="termPhaseProcessingTargetSearch"
      searchTestId="term-translation-phase-processing-target-search-input"
      searchLabel="検索"
      searchPlaceholder="名前・原文・訳語で検索"
      searchValue={processingTargetSearchQuery}
      title="処理対象"
      titleId="termPhaseProcessingTargetsHeading"
      titleColumnLabels={["原語", "訳語", "レコード種別"]}
      onSearchInput={!initialFetchDone ? undefined : handleProcessingTargetSearchInput}
      onPreviousPage={!initialFetchDone ? undefined : onProcessingTargetPreviousPage}
      onNextPage={!initialFetchDone ? undefined : onProcessingTargetNextPage}
      onPageChange={!initialFetchDone ? undefined : onProcessingTargetPageChange}
      rowTestId="term-translation-phase-processing-target-row"
      totalCountTestId="term-translation-phase-processing-target-total"
      emptyStateTestId="term-translation-phase-processing-target-empty"
    />
  </div>
</section>

<style>
  .job-run-shell {
    position: relative;
    display: grid;
    gap: 1.25rem;
  }

  .phase-controls-grid {
    display: grid;
    gap: 1.25rem;
    grid-template-columns: repeat(2, minmax(0, 1fr));
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
