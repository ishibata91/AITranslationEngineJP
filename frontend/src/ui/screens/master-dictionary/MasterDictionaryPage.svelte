<script lang="ts">
  import { onMount } from "svelte"

  import type {
    CreateMasterDictionaryScreenController,
    MasterDictionaryScreenControllerContract
  } from "@application/contract/master-dictionary"

  import DictionaryDeleteModal from "./DictionaryDeleteModal.svelte"
  import DictionaryDetailPanel from "./DictionaryDetailPanel.svelte"
  import DictionaryEditModal from "./DictionaryEditModal.svelte"
  import DictionaryHeader from "./DictionaryHeader.svelte"
  import DictionaryImportPanel from "./DictionaryImportPanel.svelte"
  import DictionaryListPanel from "./DictionaryListPanel.svelte"

  interface Props {
    createController: CreateMasterDictionaryScreenController | null
  }

  let { createController }: Props = $props()

  function resolveController(): MasterDictionaryScreenControllerContract {
    if (!createController) {
      throw new Error(
        "master dictionary screen controller factory is not provided"
      )
    }

    return createController()
  }

  const controller = resolveController()
  let viewModel = $state(controller.getViewModel())

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

  function chooseXmlFile(): void {
    if (viewModel.isImportRunning) {
      return
    }

    const input = document.getElementById("xmlFileInput")
    if (input instanceof HTMLInputElement) {
      input.click()
    }
  }

  function handleXmlSelected(event: Event): void {
    const target = event.currentTarget
    if (!(target instanceof HTMLInputElement)) {
      return
    }

    const file = target.files?.[0] ?? null
    controller.stageXmlImport(file)
  }

  function resetImportSelection(): void {
    if (viewModel.isImportRunning) {
      return
    }

    const input = document.getElementById("xmlFileInput")
    if (input instanceof HTMLInputElement) {
      input.value = ""
    }

    controller.resetImportSelection()
  }

  function goToNextPage(): void {
    controller.goToNextPage()
  }

  function goToPrevPage(): void {
    controller.goToPrevPage()
  }

  function handleCategoryChange(event: Event): void {
    controller.handleCategoryChange(event)
  }

  function handleSearchInput(event: Event): void {
    controller.handleSearchInput(event)
  }

  function selectRow(entryId: string): void {
    void controller.selectRow(entryId)
  }

  function setFormCategory(event: Event): void {
    controller.setFormCategory(event)
  }

  function setFormOrigin(event: Event): void {
    controller.setFormOrigin(event)
  }

  function setFormSource(event: Event): void {
    controller.setFormSource(event)
  }

  function setFormTranslation(event: Event): void {
    controller.setFormTranslation(event)
  }
</script>

<section
  class="master-dictionary-shell"
  data-testid="master-dictionary-master-dictionary-screen"
  id="masterDictionaryView"
>
  <DictionaryHeader
    errorMessage={viewModel.errorMessage}
    gatewayStatus={viewModel.gatewayStatus}
  />

  <DictionaryImportPanel
    hasStagedFile={viewModel.hasStagedFile}
    importProgress={viewModel.importProgress}
    importStatusText={viewModel.importStatusText}
    importStatusValue={viewModel.importStatusValue}
    importSummary={viewModel.importSummary}
    isImportRunning={viewModel.isImportRunning}
    selectedEntry={viewModel.selectedEntry}
    selectedFileName={viewModel.selectedFileName}
    {chooseXmlFile}
    {handleXmlSelected}
    {resetImportSelection}
    startImport={() => void controller.startImport()}
  />

  <section
    class="content-grid"
    data-testid="master-dictionary-dictionary-operation-region"
  >
    <DictionaryListPanel
      category={viewModel.category}
      categoryOptions={viewModel.categoryOptions}
      entries={viewModel.entries}
      listHeadline={viewModel.listHeadline}
      page={viewModel.page}
      pageStatusText={viewModel.pageStatusText}
      query={viewModel.query}
      selectedId={viewModel.selectedId}
      selectionStatusText={viewModel.selectionStatusText}
      totalPages={viewModel.totalPages}
      {goToNextPage}
      {goToPrevPage}
      {handleCategoryChange}
      {handleSearchInput}
      openCreateModal={() => controller.openCreateModal()}
      {selectRow}
    />

    <DictionaryDetailPanel
      detailSublineText={viewModel.detailSublineText}
      selectedEntry={viewModel.selectedEntry}
      openDeleteModal={() => controller.openDeleteModal()}
      openEditModal={() => controller.openEditModal()}
    />
  </section>
</section>

<DictionaryEditModal
  categoryOptions={viewModel.categoryOptions}
  formCategory={viewModel.formCategory}
  formOrigin={viewModel.formOrigin}
  formSource={viewModel.formSource}
  formTranslation={viewModel.formTranslation}
  modalState={viewModel.modalState}
  closeEditModal={() => controller.closeEditModal()}
  saveCurrentEntry={() => void controller.saveCurrentEntry()}
  {setFormCategory}
  {setFormOrigin}
  {setFormSource}
  {setFormTranslation}
/>

<DictionaryDeleteModal
  modalState={viewModel.modalState}
  selectedEntry={viewModel.selectedEntry}
  closeDeleteModal={() => controller.closeDeleteModal()}
  deleteCurrentEntry={() => void controller.deleteCurrentEntry()}
/>

<style>
  .master-dictionary-shell {
    display: grid;
    gap: 16px;
  }

  .content-grid {
    display: grid;
    grid-template-columns: minmax(0, 1.45fr) minmax(320px, 0.8fr);
    gap: 18px;
    align-items: start;
  }

  @media (max-width: 1180px) {
    .content-grid {
      grid-template-columns: 1fr;
    }
  }
</style>
