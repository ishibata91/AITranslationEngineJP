<script lang="ts">
  import { onMount } from "svelte"

  import type {
    CreateMasterPersonaScreenController,
    MasterPersonaEditableFieldMap,
    MasterPersonaScreenControllerContract
  } from "@application/contract/master-persona/master-persona-screen-contract"
  import type {
    GenerationSetupPanelProps,
    PersonaActionModalProps,
    PersonaReviewPanelProps,
    RunStatusPanelProps
  } from "./master-persona-panel-props"

  import GenerationSetupPanel from "./GenerationSetupPanel.svelte"
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

  <div
    class="status-row"
    aria-label="マスターペルソナ作成状態"
    data-testid="master-persona-screen-status-region"
  >
    <span class="status-label">作成状態</span>
    <span class="status-pill" role="status">{viewModel.runStatus.runState}</span
    >
  </div>

  <GenerationSetupPanel {...generationSetupPanelProps} />

  <RunStatusPanel {...runStatusPanelProps} />

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

  .status-row {
    align-items: center;
    display: flex;
    flex-wrap: wrap;
    gap: 8px;
    min-width: 0;
  }

  .status-label {
    color: var(--muted);
    font-size: 12px;
    letter-spacing: 0.08em;
    text-transform: uppercase;
  }

  .status-pill {
    align-items: center;
    background: rgba(255, 255, 255, 0.04);
    border: 0.5px solid rgba(255, 186, 56, 0.22);
    border-radius: 999px;
    color: var(--text);
    display: inline-flex;
    min-height: 30px;
    padding: 0 12px;
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
</style>
