<script lang="ts">
  import { onMount } from "svelte"

  import type {
    CreateMasterPersonaScreenController,
    MasterPersonaEditableFieldMap,
    MasterPersonaScreenControllerContract
  } from "@application/contract/master-persona/master-persona-screen-contract"

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
  <section class="hero-panel" aria-labelledby="masterPersonaHeading">
    <div class="hero-copy">
      <p class="eyebrow">生成準備</p>
      <h1 id="masterPersonaHeading">マスターペルソナ作成</h1>
      <p class="lead">
        ベースゲームや大型 Mod の NPC を対象に、翻訳前の準備として
        ペルソナをまとめて作成します。作成後は一覧と詳細で同じ画面から確認できます。
      </p>
    </div>
    <span class="hero-status" role="status">{viewModel.runStatus.runState}</span>
  </section>

  {#if noticeText}
    <p class:notice-error={noticeTone === "error"} class="notice-banner" role="status">
      {noticeText}
    </p>
  {/if}

  <GenerationSetupPanel
    {handleAIExecutionMethodChange}
    {handleAIModelChange}
    {handleAIProviderChange}
    {handleJsonSelected}
    {isAISettingsRefreshing}
    {refreshAISettings}
    {resetJsonSelection}
    {viewModel}
    chooseJsonFile={chooseJsonFile}
    saveAISettings={() => void controller.saveAISettings()}
    startGeneration={() => void controller.executeGeneration()}
  />

  <RunStatusPanel
    {viewModel}
    cancelGeneration={() => void controller.cancelGeneration()}
    interruptGeneration={() => void controller.interruptGeneration()}
  />

  <PersonaReviewPanel
    {viewModel}
    editCurrent={() => controller.openEditModal()}
    goToNextPage={() => controller.goToNextPage()}
    goToPrevPage={() => controller.goToPrevPage()}
    openDelete={() => controller.openDeleteModal()}
    selectRow={selectPersonaRow}
    {updateKeyword}
    {updatePluginFilter}
  />

  <PersonaActionModal
    editForm={viewModel.editForm}
    modalState={viewModel.modalState}
    selectedEntry={viewModel.selectedEntry}
    closeDelete={() => controller.closeDeleteModal()}
    closeEdit={() => controller.closeEditModal()}
    deleteCurrentEntry={() => void controller.deleteCurrentEntry()}
    saveCurrentEntry={() => void controller.saveCurrentEntry()}
    setEditFormField={updateEditFormField}
  />
</section>

<style>
  .master-persona-shell {
    display: grid;
    gap: 18px;
    min-width: 0;
  }

  .hero-panel,
  .notice-banner {
    border-radius: 20px;
    border: 0.5px solid var(--line);
    box-shadow: var(--shadow);
    min-width: 0;
  }

  .hero-panel {
    align-items: start;
    background:
      radial-gradient(circle at top right, rgba(255, 186, 56, 0.16), transparent 38%),
      rgba(17, 13, 12, 0.42);
    backdrop-filter: blur(24px);
    display: flex;
    flex-wrap: wrap;
    gap: 14px;
    justify-content: space-between;
    padding: clamp(20px, 3vw, 28px);
  }

  .hero-copy {
    display: grid;
    gap: 10px;
    min-width: 0;
  }

  .eyebrow {
    color: var(--muted);
    font-size: 12px;
    letter-spacing: 0.1em;
    margin: 0;
    text-transform: uppercase;
  }

  h1,
  p {
    margin: 0;
  }

  h1 {
    font-size: clamp(1.8rem, 2.8vw, 2.4rem);
    line-height: 1.2;
    overflow-wrap: anywhere;
  }

  .lead {
    color: var(--muted);
    line-height: 1.7;
    max-width: 62ch;
  }

  .hero-status,
  .notice-banner {
    align-items: center;
    display: inline-flex;
    overflow-wrap: anywhere;
  }

  .hero-status {
    background: rgba(255, 255, 255, 0.04);
    border: 0.5px solid rgba(255, 186, 56, 0.22);
    border-radius: 999px;
    color: var(--text);
    min-height: 40px;
    padding: 0 16px;
  }

  .notice-banner {
    background: rgba(255, 255, 255, 0.04);
    color: var(--muted);
    line-height: 1.6;
    padding: 14px 18px;
  }

  .notice-error {
    border-color: rgba(255, 156, 124, 0.35);
    color: #ffd5cb;
  }

  @media (max-width: 640px) {
    .hero-status {
      max-width: 100%;
      width: 100%;
    }
  }
</style>
