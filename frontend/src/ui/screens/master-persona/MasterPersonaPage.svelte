<script lang="ts">
  import { onMount } from "svelte"

  import type {
    CreateMasterPersonaScreenController,
    MasterPersonaEditableFieldMap,
    MasterPersonaScreenControllerContract
  } from "@application/contract/master-persona/master-persona-screen-contract"
  import PhaseStatusPanel from "@ui/components/PhaseStatusPanel.svelte"
  import type {
    PhaseMetricCounter,
    PhaseStateToken
  } from "@ui/components/phase-panel-types"

  import GenerationSetupPanel from "./GenerationSetupPanel.svelte"
  import MasterPersonaAISettingsPanel from "./MasterPersonaAISettingsPanel.svelte"
  import type {
    GenerationSetupPanelProps,
    PersonaActionModalProps,
    PersonaReviewPanelProps,
    RunStatusPanelProps
  } from "./master-persona-panel-props"
  import PersonaActionModal from "./PersonaActionModal.svelte"
  import PersonaReviewPanel from "./PersonaReviewPanel.svelte"
  import RunStatusPanel from "./RunStatusPanel.svelte"

  interface Props {
    createController: CreateMasterPersonaScreenController | null
  }

  let { createController }: Props = $props()

  function resolveController(): MasterPersonaScreenControllerContract {
    if (!createController) {
      throw new Error(
        "master persona screen controller factory is not provided"
      )
    }
    return createController()
  }

  const controller = resolveController()
  let viewModel = $state(controller.getViewModel())
  let isAISettingsRefreshing = $state(false)

  const unsubscribe = controller.subscribe((nextViewModel) => {
    viewModel = nextViewModel
  })

  onMount(() => {
    void controller.mount()
    return () => {
      unsubscribe()
      controller.dispose()
    }
  })

  const noticeText = $derived(
    viewModel.errorMessage || viewModel.aiSettingsMessage || ""
  )
  const noticeTone = $derived(viewModel.errorMessage ? "error" : "info")
  const generationSetupPanelProps = $derived<GenerationSetupPanelProps>({
    aiSettings: viewModel.aiSettings,
    aiSettingsStatusText: viewModel.aiSettingsStatusText,
    aiSettingsWarningText: viewModel.aiSettingsWarningText,
    aiProviderLabel: viewModel.aiProviderLabel,
    canSelectModel: viewModel.canSelectModel,
    canStartGeneration: viewModel.canStartGeneration,
    executionMethodOptions: viewModel.executionMethodOptions,
    isRunActive: viewModel.isRunActive,
    modelOptions: viewModel.modelOptions,
    modelSettingsCardViewModel: viewModel.modelSettingsCardViewModel,
    preview: viewModel.preview,
    selectedFileName: viewModel.selectedFileName,
    selectedFileReference: viewModel.selectedFileReference,
    isAISettingsRefreshing,
    handleJsonSelected,
    chooseJsonFile,
    resetJsonSelection,
    handleAIProviderChange,
    handleAIModelChange,
    handleAIExecutionMethodChange,
    refreshAISettings,
    startGeneration: () => void controller.executeGeneration(),
    saveAISettings: () => void controller.saveAISettings()
  })
  const aiSettingsPanelProps = $derived({
    aiSettings: viewModel.aiSettings,
    aiSettingsStatusText: viewModel.aiSettingsStatusText,
    aiSettingsWarningText: viewModel.aiSettingsWarningText,
    aiProviderLabel: viewModel.aiProviderLabel,
    canSelectModel: viewModel.canSelectModel,
    executionMethodOptions: viewModel.executionMethodOptions,
    isAISettingsRefreshing,
    modelOptions: viewModel.modelOptions,
    modelSettingsCardViewModel: viewModel.modelSettingsCardViewModel,
    handleAIProviderChange,
    handleAIModelChange,
    handleAIExecutionMethodChange,
    refreshAISettings,
    saveAISettings: () => void controller.saveAISettings()
  })
  const runStatusPanelProps = $derived<RunStatusPanelProps>({
    isRunActive: viewModel.isRunActive,
    progressPercent: viewModel.progressPercent,
    runStatus: viewModel.runStatus,
    cancelGeneration: () => void controller.cancelGeneration(),
    interruptGeneration: () => void controller.interruptGeneration()
  })
  const personaReviewPanelProps = $derived<PersonaReviewPanelProps>({
    canMutate: viewModel.canMutate,
    items: viewModel.items,
    keyword: viewModel.keyword,
    page: viewModel.page,
    pageSize: viewModel.pageSize,
    pluginFilter: viewModel.pluginFilter,
    pluginOptions: viewModel.pluginOptions,
    selectedEntry: viewModel.selectedEntry,
    selectedIdentityKey: viewModel.selectedIdentityKey,
    totalCount: viewModel.totalCount,
    totalPages: viewModel.totalPages,
    editCurrent: () => controller.openEditModal(),
    goToNextPage: () => controller.goToNextPage(),
    goToPrevPage: () => controller.goToPrevPage(),
    openDelete: () => controller.openDeleteModal(),
    selectRow: selectPersonaRow,
    updateKeyword,
    updatePluginFilter
  })
  const personaActionModalProps = $derived<PersonaActionModalProps>({
    editForm: viewModel.editForm,
    errorMessage: viewModel.errorMessage,
    modalState: viewModel.modalState,
    selectedEntry: viewModel.selectedEntry,
    closeDelete: () => controller.closeDeleteModal(),
    closeEdit: () => controller.closeEditModal(),
    deleteCurrentEntry: () => void controller.deleteCurrentEntry(),
    saveCurrentEntry: () => void controller.saveCurrentEntry(),
    setEditFormField: updateEditFormField
  })
  const personaViewState = $derived<PhaseStateToken>(
    viewModel.errorMessage
      ? "failed"
      : viewModel.isRunActive
        ? "running"
        : viewModel.runStatus.runState.includes("完了")
          ? "completed"
          : viewModel.runStatus.runState.includes("中止")
            ? "canceled"
            : viewModel.runStatus.runState.includes("中断")
              ? "paused"
              : "idle_ready"
  )
  const personaStatusTitle = $derived(
    viewModel.isRunActive
      ? "ペルソナを生成中です"
      : viewModel.canStartGeneration
        ? "ペルソナを作成できます"
        : viewModel.runStatus.message
  )
  const personaStatusText = $derived(
    viewModel.hasPreview
      ? `候補 ${viewModel.preview?.candidateCount ?? 0} 件を確認できます。`
      : "JSON と AI 設定を確認してください。"
  )
  const personaMetrics = $derived<PhaseMetricCounter[]>([
    { label: "保存済み", value: String(viewModel.totalCount) },
    {
      label: "候補",
      value: String(viewModel.preview?.candidateCount ?? 0)
    },
    {
      label: "新規作成",
      value: String(viewModel.preview?.newlyAddableCount ?? 0)
    },
    { label: "処理済み", value: String(viewModel.runStatus.processedCount) },
    { label: "作成済み", value: String(viewModel.runStatus.successCount) }
  ])

  function chooseJsonFile(): void {
    const input = document.getElementById("masterPersonaJsonInput")
    if (input instanceof HTMLInputElement) {
      input.click()
    }
  }

  function handleJsonSelected(event: Event): void {
    const target = event.currentTarget
    if (!(target instanceof HTMLInputElement)) {
      return
    }
    controller.stageJsonSelection(target.files?.[0] ?? null)
  }

  function resetJsonSelection(): void {
    const input = document.getElementById("masterPersonaJsonInput")
    if (input instanceof HTMLInputElement) {
      input.value = ""
    }
    controller.resetJsonSelection()
  }

  function handleAIProviderChange(event: Event): void {
    if (isAISettingsRefreshing) {
      return
    }
    controller.setAIProvider(event)
  }

  function handleAIModelChange(event: Event): void {
    if (isAISettingsRefreshing) {
      return
    }
    controller.setAIModel(event)
  }

  function handleAIExecutionMethodChange(event: Event): void {
    if (isAISettingsRefreshing) {
      return
    }
    controller.setAIExecutionMethod(event)
  }

  async function refreshAISettings(): Promise<void> {
    if (!controller.refreshAISettings || isAISettingsRefreshing) {
      return
    }

    isAISettingsRefreshing = true
    try {
      await controller.refreshAISettings()
    } finally {
      isAISettingsRefreshing = false
    }
  }

  function selectPersonaRow(identityKey: string): void {
    void controller.selectRow(identityKey)
  }

  function updateKeyword(event: Event): void {
    controller.handleSearchInput(event)
  }

  function updatePluginFilter(event: Event): void {
    controller.handlePluginFilterChange(event)
  }

  function updateEditFormField(
    field: keyof MasterPersonaEditableFieldMap,
    event: Event
  ): void {
    controller.setEditFormField(field, event)
  }
</script>

<section class="master-persona-shell" id="masterPersonaView">
  <h1 class="sr-only">マスターペルソナ作成</h1>

  {#if noticeText}
    <p
      class:notice-error={noticeTone === "error"}
      class="notice-banner"
      role="status"
    >
      {noticeText}
    </p>
  {/if}

  <PhaseStatusPanel
    eyebrow="基盤データ"
    title="マスターペルソナ"
    gatewayStatus={viewModel.gatewayStatus}
    lead="JSON 入力、AI 設定、生成状態、ペルソナ一覧を同じ画面で操作します。"
    state={personaViewState}
    stateLabel={viewModel.runStatus.runState}
    statusTitle={personaStatusTitle}
    statusText={personaStatusText}
    errorMessage={viewModel.errorMessage}
    showGatewayStatus={false}
    testId="master-persona-screen-status-region"
    statusTestId="master-persona-status-summary-card"
    metrics={personaMetrics}
  />

  <GenerationSetupPanel {...generationSetupPanelProps} />

  <section class="phase-controls-grid">
    <RunStatusPanel {...runStatusPanelProps} />
    <MasterPersonaAISettingsPanel {...aiSettingsPanelProps} />
  </section>

  <PersonaReviewPanel {...personaReviewPanelProps} />

  <PersonaActionModal {...personaActionModalProps} />
</section>

<style>
  .master-persona-shell {
    display: grid;
    gap: 14px;
    min-width: 0;
  }

  .notice-banner {
    border-radius: 20px;
    border: 0.5px solid var(--line);
    box-shadow: var(--shadow);
    min-width: 0;
  }

  p {
    margin: 0;
  }

  .notice-banner {
    align-items: center;
    display: inline-flex;
    overflow-wrap: anywhere;
  }

  .notice-banner {
    background: rgba(255, 255, 255, 0.04);
    color: var(--muted);
    line-height: 1.6;
    padding: 14px 18px;
  }

  .phase-controls-grid {
    align-items: start;
    display: grid;
    gap: 1.25rem;
    grid-template-columns: minmax(0, 0.95fr) minmax(0, 1.05fr);
  }

  .sr-only {
    border: 0;
    clip: rect(0 0 0 0);
    height: 1px;
    margin: -1px;
    overflow: hidden;
    padding: 0;
    position: absolute;
    white-space: nowrap;
    width: 1px;
  }

  .notice-error {
    border-color: rgba(255, 156, 124, 0.35);
    color: #ffd5cb;
  }

  @media (max-width: 900px) {
    .phase-controls-grid {
      grid-template-columns: 1fr;
    }
  }
</style>
