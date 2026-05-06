<script lang="ts">
  import { onMount } from "svelte"

  import type {
    CreateTranslationInputScreenController,
    TranslationInputScreenControllerContract
  } from "@application/contract/translation-input"
  import {
    ERROR_LABELS,
    STATUS_LABELS,
    WARNING_LABELS
  } from "@application/presenter/translation-input"

  import DataLoadHero from "./DataLoadHero.svelte"
  import DataLoadImportPanel from "./DataLoadImportPanel.svelte"
  import LoadedInputDetail from "./LoadedInputDetail.svelte"
  import LoadedInputList from "./LoadedInputList.svelte"

  interface Props {
    createController: CreateTranslationInputScreenController | null
    onOpenJobSetup?: () => void
  }

  let { createController, onOpenJobSetup = undefined }: Props = $props()

  function resolveController(): TranslationInputScreenControllerContract {
    if (!createController) {
      throw new Error(
        "translation input screen controller factory is not provided"
      )
    }

    return createController()
  }

  const controller = resolveController()
  let viewModel = $state(controller.getViewModel())
  let fileInput: HTMLInputElement | null = null

  const unsubscribe = controller.subscribe((nextViewModel) => {
    viewModel = nextViewModel
  })

  function clearJsonFileInput(): void {
    if (fileInput) {
      fileInput.value = ""
    }
  }

  $effect(() => {
    if (!viewModel.hasStagedFile) {
      clearJsonFileInput()
    }
  })

  onMount(() => {
    void controller.mount()

    return () => {
      unsubscribe()
      controller.dispose()
    }
  })

  function chooseJsonFile(): void {
    if (viewModel.isImporting) {
      return
    }

    clearJsonFileInput()
    fileInput?.click()
  }

  function handleJsonSelected(event: Event): void {
    const target = event.currentTarget
    if (!(target instanceof HTMLInputElement)) {
      return
    }

    const file = target.files?.[0] ?? null
    void controller.stageJsonImport(file)
  }

  function resetImportSelection(): void {
    clearJsonFileInput()
    controller.resetImportSelection()
  }

  function selectLoadedInput(localId: string): void {
    controller.selectItem(localId)
  }

  function rebuildSelectedInput(): Promise<void> {
    return controller.rebuildSelected()
  }

  function startImport(): Promise<void> {
    return controller.startImport()
  }

  function formatStatus(localStatus: string): string {
    return (
      STATUS_LABELS[localStatus as keyof typeof STATUS_LABELS] ?? localStatus
    )
  }

  function formatErrorKind(errorKind: string | null): string {
    if (!errorKind) {
      return "-"
    }

    return ERROR_LABELS[errorKind] ?? errorKind
  }

  function formatWarningKind(kind: string): string {
    return WARNING_LABELS[kind] ?? kind
  }

  function formatDate(timestamp: string): string {
    if (!timestamp) {
      return "-"
    }

    const date = new Date(timestamp)
    if (Number.isNaN(date.getTime())) {
      return timestamp
    }

    return date.toLocaleString("ja-JP")
  }

  function localizeUiText(text: string): string {
    return text
      .replaceAll("input review", "読み込み済みデータ")
      .replaceAll("sample field", "サンプル項目")
      .replaceAll("cache", "キャッシュ")
      .replaceAll("error kind", "問題区分")
      .replaceAll("xEdit JSON", "xEdit の JSON")
      .replaceAll("JSON file", "JSON")
  }
</script>

<section class="data-load-shell" id="translationInputReviewView">
  <input
    accept=".json,application/json"
    bind:this={fileInput}
    class="file-input"
    id="translationInputFile"
    onchange={handleJsonSelected}
    type="file"
  />

  <DataLoadHero
    errorMessage={viewModel.errorMessage}
    gatewayStatus={viewModel.gatewayStatus}
    operationStatusLabel={viewModel.operationStatusLabel}
    operationStatusText={localizeUiText(viewModel.operationStatusText)}
  />

  <DataLoadImportPanel
    canImport={viewModel.canImport}
    hasStagedFile={viewModel.hasStagedFile}
    isImporting={viewModel.isImporting}
    stagedFileHash={viewModel.stagedFileHash}
    stagedFileName={viewModel.stagedFileName}
    stagedFilePath={viewModel.stagedFilePath}
    onChooseJson={chooseJsonFile}
    onResetSelection={resetImportSelection}
    onStartImport={startImport}
  />

  <section class="content-grid">
    <LoadedInputList
      emptyStateText={localizeUiText(viewModel.emptyStateText)}
      formatDate={formatDate}
      formatErrorKind={formatErrorKind}
      formatStatus={formatStatus}
      items={viewModel.items}
      selectedItemId={viewModel.selectedItemId}
      totalItemCountLabel={localizeUiText(viewModel.totalItemCountLabel)}
      onSelectItem={selectLoadedInput}
    />

    <LoadedInputDetail
      canRebuildSelected={viewModel.canRebuildSelected}
      formatDate={formatDate}
      formatErrorKind={formatErrorKind}
      formatWarningKind={formatWarningKind}
      isRebuilding={viewModel.isRebuilding}
      latestOutcomeText={localizeUiText(viewModel.latestOutcomeText)}
      latestOutcomeTitle={localizeUiText(viewModel.latestOutcomeTitle)}
      selectedItem={viewModel.selectedItem}
      selectionStatusText={localizeUiText(viewModel.selectionStatusText)}
      onOpenJobSetup={onOpenJobSetup}
      onRebuild={rebuildSelectedInput}
    />
  </section>
</section>

<style>
  .data-load-shell {
    display: grid;
    gap: 1.25rem;
  }

  .file-input {
    display: none;
  }

  .content-grid {
    display: grid;
    gap: 1.25rem;
    grid-template-columns: minmax(0, 0.9fr) minmax(0, 1.1fr);
    align-items: start;
  }

  @media (max-width: 960px) {
    .content-grid {
      grid-template-columns: 1fr;
    }
  }
</style>
