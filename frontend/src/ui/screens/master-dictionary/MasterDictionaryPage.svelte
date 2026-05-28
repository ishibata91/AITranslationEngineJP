<script lang="ts">
  import { onMount } from "svelte"

  import type {
    CreateMasterDictionaryScreenController,
    MasterDictionaryScreenControllerContract
  } from "@application/contract/master-dictionary"
  import type { MasterDictionaryEntrySummary } from "@application/gateway-contract/master-dictionary"
  import PhaseStatusPanel from "@ui/components/PhaseStatusPanel.svelte"
  import ProcessingTargetListWrapper from "@ui/components/ProcessingTargetListWrapper.svelte"
  import type {
    PhaseMetricCounter,
    PhaseStateToken
  } from "@ui/components/phase-panel-types"
  import type { ProcessingTargetListItem } from "@ui/components/processing-target-list-panel-types"

  import DictionaryDeleteModal from "./DictionaryDeleteModal.svelte"
  import DictionaryEditModal from "./DictionaryEditModal.svelte"
  import DictionaryImportPanel from "./DictionaryImportPanel.svelte"
  import DictionaryImportProgressPanel from "./DictionaryImportProgressPanel.svelte"

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

  function buildDictionaryTargetItems(
    entries: MasterDictionaryEntrySummary[]
  ): ProcessingTargetListItem[] {
    return entries.map((entry) => {
      const note =
        viewModel.selectedEntry?.id === entry.id
          ? viewModel.selectedEntry.note
          : ""

      return {
        id: entry.id,
        name: `${entry.source} / ${entry.translation}`,
        titleParts: [
          { text: `原文: ${entry.source}` },
          { text: `訳語: ${entry.translation}` }
        ],
        detail: "",
        metadata: [
          { label: "辞書 ID", value: entry.id },
          { label: "カテゴリ", value: entry.category },
          { label: "登録元", value: entry.origin },
          { label: "最終更新", value: entry.updatedAt },
          { label: "メモ", value: note || "未取得" }
        ]
      }
    })
  }

  const dictionaryTargetItems = $derived(
    buildDictionaryTargetItems(viewModel.entries)
  )
  const dictionaryCategoryOptions = $derived(
    viewModel.categoryOptions.map((option) => ({
      label: option,
      value: option
    }))
  )
  const dictionaryViewState = $derived<PhaseStateToken>(
    viewModel.errorMessage
      ? "failed"
      : viewModel.isImportRunning
        ? "running"
        : viewModel.hasStagedFile
          ? "ready"
          : "idle_ready"
  )
  const dictionaryPhaseTitle = $derived(
    viewModel.isImportRunning
      ? "現在のフェーズ: XML 取り込み"
      : viewModel.hasStagedFile
        ? "現在のフェーズ: XML 取り込み準備"
        : "現在のフェーズ: マスター辞書管理"
  )
  const dictionaryStatusText = $derived(
    viewModel.importSummary
      ? `「${viewModel.importSummary.fileName}」を反映しました。`
      : viewModel.isImportRunning
        ? viewModel.importStatusText
        : viewModel.hasStagedFile
          ? "選択した XML を取り込む前に内容を確認します。"
          : "XML 取り込みは未開始です。"
  )
  const dictionaryStateLabel = $derived(
    viewModel.importStatusValue === "完了"
      ? "取り込み完了"
      : viewModel.importStatusValue
  )
  const dictionaryMetrics = $derived<PhaseMetricCounter[]>([
    { label: "保存済み", value: String(viewModel.totalCount) },
    { label: "表示中", value: viewModel.pageStatusText },
    { label: "カテゴリ", value: viewModel.category },
    { label: "選択", value: viewModel.selectedEntry?.source ?? "-" }
  ])
</script>

<section
  class="master-dictionary-shell"
  data-testid="master-dictionary-master-dictionary-screen"
  id="masterDictionaryView"
>
  <PhaseStatusPanel
    eyebrow="基盤データ"
    title="マスター辞書"
    errorMessage={viewModel.errorMessage}
    gatewayStatus={viewModel.gatewayStatus}
    lead="一覧、作成、更新、削除、XML 取り込みを同じ画面で操作します。"
    state={dictionaryViewState}
    stateLabel={dictionaryStateLabel}
    statusTitle={dictionaryPhaseTitle}
    statusText={dictionaryStatusText}
    testId="master-dictionary-screen-header"
    statusTestId="master-dictionary-status-summary-card"
    metrics={dictionaryMetrics}
  />

  <section class="phase-controls-grid">
    <DictionaryImportPanel
      isImportRunning={viewModel.isImportRunning}
      selectedFileName={viewModel.selectedFileName}
      {chooseXmlFile}
      {handleXmlSelected}
    />
    <DictionaryImportProgressPanel
      hasStagedFile={viewModel.hasStagedFile}
      importProgress={viewModel.importProgress}
      importStatusText={viewModel.importStatusText}
      importStatusValue={viewModel.importStatusValue}
      importSummary={viewModel.importSummary}
      isImportRunning={viewModel.isImportRunning}
      selectedEntry={viewModel.selectedEntry}
      selectedFileName={viewModel.selectedFileName}
      {resetImportSelection}
      startImport={() => void controller.startImport()}
    />
  </section>

  <section
    class="content-grid"
    data-testid="master-dictionary-dictionary-operation-region"
  >
    <ProcessingTargetListWrapper
      countText={viewModel.pageStatusText}
      filterId="categorySelect"
      filterLabel="カテゴリ"
      filterOptions={dictionaryCategoryOptions}
      filterTestId="master-dictionary-category-select"
      filterValue={viewModel.category}
      initialExpandedItemId={viewModel.selectedEntry?.id ?? viewModel.selectedId}
      items={dictionaryTargetItems}
      pageSize={30}
      searchId="searchInput"
      searchLabel="検索"
      searchPlaceholder="原文・訳語・IDで検索"
      searchTestId="master-dictionary-search-input"
      searchValue={viewModel.query}
      supportText={viewModel.listHeadline}
      title="辞書一覧"
      titleId="dictionaryTargetHeading"
      onFilterChange={handleCategoryChange}
      onSearchInput={handleSearchInput}
      onSelectItem={selectRow}
      rowTestId="master-dictionary-entry-row"
    >
      {#snippet actions()}
        <div class="toolbar-actions">
          <button
            class="button-primary"
            data-testid="master-dictionary-create-button"
            id="createButton"
            onclick={() => controller.openCreateModal()}
            type="button">新規登録</button
          >
          <button
            class="button-secondary"
            data-testid="master-dictionary-detail-edit-button"
            disabled={!viewModel.selectedEntry}
            id="editButton"
            onclick={() => controller.openEditModal()}
            type="button">更新</button
          >
          <button
            class="button-danger"
            data-testid="master-dictionary-detail-delete-button"
            disabled={!viewModel.selectedEntry}
            id="deleteButton"
            onclick={() => controller.openDeleteModal()}
            type="button">削除</button
          >
        </div>
      {/snippet}

      {#snippet summary()}
        <div class="selected-entry-summary">
          <span>{viewModel.detailSublineText}</span>
          <strong id="detailTitle">
            {viewModel.selectedEntry?.source ?? "表示できるエントリがありません"}
          </strong>
        </div>
      {/snippet}

      {#snippet footer()}
        <div class="pager-row">
          <span>{viewModel.selectionStatusText}</span>
          <span>{viewModel.pageStatusText}</span>
          <div class="pager-actions">
            <button
              class="button-secondary"
              disabled={viewModel.page === 0}
              id="prevPageButton"
              onclick={goToPrevPage}
              type="button">前の30件</button
            >
            <button
              class="button-secondary"
              disabled={viewModel.page + 1 >= viewModel.totalPages}
              id="nextPageButton"
              onclick={goToNextPage}
              type="button">次の30件</button
            >
          </div>
        </div>
      {/snippet}
    </ProcessingTargetListWrapper>
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
    gap: 1.25rem;
  }

  .phase-controls-grid {
    align-items: start;
    display: grid;
    gap: 1.25rem;
    grid-template-columns: minmax(0, 1fr);
  }

  .content-grid {
    display: grid;
    gap: 18px;
    align-items: start;
  }

  .toolbar-actions,
  .pager-row,
  .pager-actions {
    align-items: center;
    display: flex;
    flex-wrap: wrap;
    gap: 10px;
    justify-content: space-between;
  }

  .pager-row {
    color: var(--muted);
  }

  .selected-entry-summary {
    background: rgba(255, 255, 255, 0.03);
    border: 1px solid var(--line);
    border-radius: 12px;
    display: grid;
    gap: 6px;
    min-width: 0;
    padding: 12px;
  }

  .selected-entry-summary strong {
    color: var(--text);
    overflow-wrap: anywhere;
  }

  .button-primary,
  .button-secondary,
  .button-danger {
    border: 1px solid transparent;
    border-radius: 999px;
    font: inherit;
    min-height: 36px;
    padding: 0 14px;
  }

  .button-primary {
    background: linear-gradient(135deg, var(--primary) 0%, #ef9d20 100%);
    color: #3a2400;
  }

  .button-secondary {
    background: rgba(255, 255, 255, 0.04);
    border-color: var(--line);
    color: var(--text);
  }

  .button-danger {
    background: linear-gradient(135deg, #ffc0ab 0%, #ff9975 100%);
    color: #3d1512;
  }

  button:disabled {
    cursor: not-allowed;
    opacity: 0.5;
  }
</style>
