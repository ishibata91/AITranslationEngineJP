<script lang="ts">
  import { onMount } from "svelte"

  import type {
    CreateTranslationInputScreenController,
    TranslationInputScreenControllerContract
  } from "@application/contract/translation-input"
  import type { TranslationJobManagementJobRunTarget } from "@application/contract/translation-job-management/translation-job-management-screen-types"
  import {
    canOpenJobSetup,
    ERROR_LABELS,
    STATUS_LABELS
  } from "@application/presenter/translation-input"
  import StickyActionFooter from "@ui/components/StickyActionFooter.svelte"

  import DataLoadHero from "./DataLoadHero.svelte"
  import DataLoadImportPanel from "./DataLoadImportPanel.svelte"
  import LoadedInputList from "./LoadedInputList.svelte"

  interface Props {
    createController: CreateTranslationInputScreenController | null
    onOpenJobRun?: (target: TranslationJobManagementJobRunTarget) => void
  }

  let { createController, onOpenJobRun = undefined }: Props = $props()

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
  const showNextActionFooter = $derived(canOpenJobSetup(viewModel.selectedItem))
  const selectedInputHasExistingJob = $derived(
    viewModel.selectedItem?.warnings.some((warning) =>
      warning.message.includes("既存")
    ) ?? false
  )
  const footerReasons = $derived(
    showNextActionFooter
      ? []
      : ["登録済みまたは警告ありの入力データを選択してください。"]
  )

  const unsubscribe = controller.subscribe((nextViewModel) => {
    viewModel = nextViewModel
  })

  function clearJsonFileInput(): void {
    const input = document.getElementById("translationInputFile")
    if (input instanceof HTMLInputElement) {
      input.value = ""
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
    const input = document.getElementById("translationInputFile")
    if (input instanceof HTMLInputElement) {
      input.click()
    }
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

  async function openTermTranslation(): Promise<void> {
    const selectedItem = viewModel.selectedItem
    if (!selectedItem || selectedItem.inputId === null) {
      return
    }

    const response = await controller.createTranslationJobFromSelected?.()
    if (response && (!response.accepted || response.jobId === undefined)) {
      return
    }

    onOpenJobRun?.({
      jobId: response?.jobId ?? selectedItem.inputId,
      stateLabel: response?.jobState ?? "作成済み",
      stateDescription: "選択した入力データから翻訳ジョブを作成しました。",
      currentPhase: response?.currentPhase ?? "term_translation",
      currentPhaseLabel: "単語翻訳",
      progressLabel: "未開始",
      inputSourceLabel: selectedItem.fileName,
      sourcePath: selectedItem.filePath
    })
  }
</script>

<section class="data-load-shell" id="translationInputReviewView">
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
    stagedFileName={viewModel.stagedFileName}
    onJsonSelected={handleJsonSelected}
    onChooseJson={chooseJsonFile}
    onResetSelection={resetImportSelection}
    onStartImport={startImport}
  />

  <section class="content-grid">
    <LoadedInputList
      emptyStateText={localizeUiText(viewModel.emptyStateText)}
      {formatDate}
      {formatErrorKind}
      {formatStatus}
      items={viewModel.items}
      selectedItemId={viewModel.selectedItemId}
      onSelectItem={selectLoadedInput}
    />
  </section>

  {#if showNextActionFooter}
    <StickyActionFooter
      dataTestId="translation-input-review-next-action-footer"
      title="次の作業"
      titleId="translationInputNextNavigationHeading"
      description="選択した入力データで翻訳ジョブを作成し、単語翻訳へ進みます。"
      reasons={footerReasons}
      emptyText={selectedInputHasExistingJob
        ? "既存の翻訳ジョブがあります。必要に応じて未完了ジョブ一覧から再開できます。"
        : "入力データを選択済みです。次に単語翻訳へ進みます。"}
      primaryLabel="単語翻訳へ進む"
      onPrimary={openTermTranslation}
    />
  {/if}
</section>

<style>
  .data-load-shell {
    display: grid;
    gap: 1.25rem;
    padding-bottom: 10rem;
  }

  .content-grid {
    display: grid;
    gap: 1.25rem;
    grid-template-columns: 1fr;
    align-items: start;
  }

  @media (max-width: 960px) {
    .content-grid {
      grid-template-columns: 1fr;
    }
  }

  @media (max-width: 720px) {
    .data-load-shell {
      padding-bottom: 14rem;
    }
  }
</style>
