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
    onProcessingTargetSearchInput = undefined,
    onProcessingTargetPreviousPage = undefined,
    onProcessingTargetNextPage = undefined,
    onProcessingTargetPageChange = undefined,
    onAISettingsChange = undefined
  }: Props = $props()
  let processingTargetSearchValue = $state("")
  const phaseActionCards = $derived(
    selectPhaseProgressActions<TermPanelActionKind>(viewModel.actionCards, {
      hiddenActionIds: ["next-phase"],
      runActionIds: ["start", "resume", "retry"],
      alwaysActionIds: ["pause"]
    })
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
  const summaryProcessingTargetItems = $derived<ProcessingTargetListItem[]>([
    {
      id: "term-translation-ai-target",
      name: "原語 / 訳語候補",
      titleParts: [
        { text: `原語候補: ${viewModel.totalTermCountLabel} 件` },
        { text: `AI 訳語候補: ${viewModel.aiTargetCountLabel} 件` }
      ],
      detail: "共通辞書に一致しない用語と固有名詞を AI サービスへ送る対象。",
      metadata: [
        { label: "対象語件数", value: viewModel.totalTermCountLabel },
        { label: "共通辞書一致", value: viewModel.dictionaryHitCountLabel },
        { label: "AI 送信対象", value: viewModel.aiTargetCountLabel },
        { label: "置換対象", value: viewModel.replacementTargetCountLabel },
        { label: "未一致", value: viewModel.unmatchedCountLabel },
        { label: "保存先", value: "翻訳ジョブ内辞書" },
        { label: "スナップショット", value: viewModel.snapshotLabel }
      ]
    }
  ])
  const displayedProcessingTargetItems = $derived(
    processingTargetPageState
      ? processingTargetPageState.items
      : providedProcessingTargetItems &&
          providedProcessingTargetItems.length > 0
        ? providedProcessingTargetItems
        : summaryProcessingTargetItems
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
  data-testid="term-translation-phase-screen"
  id="termTranslationPhaseView"
>
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
      statusTestId="term-translation-phase-ai-model-lock-state"
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

  <ProcessingTargetListWrapper
    items={filteredProcessingTargetItems}
    pageSize={processingTargetPageState?.pageSize ?? 50}
    currentPage={processingTargetPageState?.page}
    totalCount={processingTargetPageState?.totalCount}
    busy={processingTargetPageState?.busy}
    searchId="termPhaseProcessingTargetSearch"
    searchLabel="検索"
    searchPlaceholder="名前・原文・訳語で検索"
    searchValue={processingTargetSearchQuery}
    title="処理対象"
    titleId="termPhaseProcessingTargetsHeading"
    onSearchInput={handleProcessingTargetSearchInput}
    onPreviousPage={onProcessingTargetPreviousPage}
    onNextPage={onProcessingTargetNextPage}
    onPageChange={onProcessingTargetPageChange}
    rowTestId="term-translation-phase-processing-target-row"
  />
</section>

<style>
  .job-run-shell {
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
</style>
