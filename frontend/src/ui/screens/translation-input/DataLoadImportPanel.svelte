<script lang="ts">
  import FileImportPanel from "@ui/components/FileImportPanel.svelte"

  interface Props {
    stagedFileName: string
    hasStagedFile: boolean
    canImport: boolean
    isImporting: boolean
    onJsonSelected: (event: Event) => void
    onChooseJson: () => void
    onStartImport: () => void | Promise<void>
    onResetSelection: () => void
  }

  let {
    stagedFileName,
    hasStagedFile,
    canImport,
    isImporting,
    onJsonSelected,
    onChooseJson,
    onStartImport,
    onResetSelection
  }: Props = $props()
</script>

<FileImportPanel
  accept=".json,application/json"
  eyebrow="ロード準備"
  helperText="取り込みたい JSON を 1 件選び、内容を確認してから登録する。"
  inputId="translationInputFile"
  inputTestId="translation-input-review-json-file-input"
  primaryActionId="translationInputChooseJsonButton"
  primaryActionLabel="ロード対象を選ぶ"
  primaryActionDisabled={isImporting}
  selectedLabel="選択 JSON"
  selectedName={stagedFileName}
  testId="translation-input-review-load-preparation-region"
  title="入力ファイル"
  titleId="inputReviewImportHeading"
  onFileSelected={onJsonSelected}
  onPrimaryAction={onChooseJson}
>
  {#snippet actions()}
    <button
      class="button-primary"
      data-testid="translation-input-review-register-button"
      disabled={!canImport}
      onclick={() => void onStartImport()}
      type="button"
    >
      この JSON を登録
    </button>
    <button
      class="button-secondary"
      disabled={!hasStagedFile || isImporting}
      onclick={onResetSelection}
      type="button"
    >
      選び直す
    </button>
  {/snippet}
</FileImportPanel>

<style>
  .button-primary,
  .button-secondary {
    font: inherit;
    border-radius: 999px;
    padding: 0.7rem 1.1rem;
    border: 1px solid var(--line-strong);
    cursor: pointer;
  }

  .button-primary {
    background: var(--primary);
    color: #2b1900;
  }

  .button-secondary {
    background: transparent;
    color: var(--text);
  }

  .button-primary:disabled,
  .button-secondary:disabled {
    cursor: not-allowed;
    opacity: 0.55;
  }
</style>
