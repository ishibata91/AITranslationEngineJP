<script lang="ts">
  import FileImportPanel from "@ui/components/FileImportPanel.svelte"
  import type { GenerationSetupPanelProps } from "./master-persona-panel-props"

  type Props = Pick<
    GenerationSetupPanelProps,
    | "canStartGeneration"
    | "isRunActive"
    | "preview"
    | "selectedFileName"
    | "selectedFileReference"
    | "handleJsonSelected"
    | "chooseJsonFile"
    | "resetJsonSelection"
    | "startGeneration"
  >

  let {
    canStartGeneration,
    isRunActive,
    preview,
    selectedFileName,
    selectedFileReference,
    handleJsonSelected,
    chooseJsonFile,
    resetJsonSelection,
    startGeneration
  }: Props = $props()

  const isFileSelected = $derived(selectedFileReference !== null)
  const helperText = $derived(
    preview
      ? `候補 ${preview.candidateCount} 件のうち、新規 ${preview.newlyAddableCount} 件を作成できます。`
      : "ベースゲームや大型 Mod の対象 NPC から、未作成のペルソナだけを追加します。"
  )
  const previewStats = $derived([
    { label: "候補数", value: preview?.candidateCount ?? 0 },
    { label: "新規作成数", value: preview?.newlyAddableCount ?? 0 },
    { label: "既存スキップ数", value: preview?.existingCount ?? 0 }
  ])
</script>

<section
  class="setup-panel"
  aria-labelledby="setupHeading"
  data-testid="master-persona-generation-setup-panel"
>
  <h2 class="sr-only" id="setupHeading">入力ファイルと AI 設定</h2>

  <div class="setup-grid">
    <FileImportPanel
      accept=".json,application/json"
      eyebrow="入力 JSON"
      helperText={helperText}
      inputId="masterPersonaJsonInput"
      inputTestId="master-persona-json-file-input"
      primaryActionId="chooseJsonButton"
      primaryActionLabel="JSON を選択"
      primaryActionDisabled={isRunActive}
      selectedName={selectedFileName}
      stats={previewStats}
      statsId="previewStats"
      testId="master-persona-input-json-panel"
      title="入力ファイル"
      titleId="fileHeading"
      onFileSelected={handleJsonSelected}
      onPrimaryAction={chooseJsonFile}
    >
      {#snippet actions()}
        <button
          class="button-secondary"
          disabled={!isFileSelected}
          id="resetJsonButton"
          onclick={resetJsonSelection}
          type="button"
        >
          選び直す
        </button>
        <button
          class="button-primary button-wide"
          data-testid="master-persona-generate-button"
          disabled={!canStartGeneration}
          id="executeGenerationButton"
          onclick={startGeneration}
          type="button"
        >
          ペルソナを作成
        </button>
      {/snippet}
    </FileImportPanel>
  </div>
</section>

<style>
  .button-primary,
  .button-secondary {
    border-radius: 20px;
  }

  .setup-panel {
    display: grid;
    gap: 14px;
    min-width: 0;
  }

  .setup-grid {
    display: grid;
    gap: 1.25rem;
    grid-template-columns: 1fr;
  }

  h2 {
    margin: 0;
  }

  h2 {
    overflow-wrap: anywhere;
  }

  .button-primary,
  .button-secondary {
    border-radius: 999px;
    cursor: pointer;
    min-height: 40px;
    min-width: 0;
    overflow-wrap: anywhere;
    padding: 0 16px;
  }

  .button-primary {
    background: linear-gradient(135deg, var(--primary) 0%, #f0a51f 100%);
    border: 0.5px solid transparent;
    color: #3f2400;
  }

  .button-secondary {
    background: rgba(255, 255, 255, 0.04);
    border: 0.5px solid var(--line);
    color: var(--text);
  }

  .button-wide {
    margin-left: auto;
  }

  .sr-only {
    border: 0;
    clip: rect(0 0 0 0);
    clip-path: inset(50%);
    height: 1px;
    margin: -1px;
    overflow: hidden;
    padding: 0;
    position: absolute;
    white-space: nowrap;
    width: 1px;
  }

  button:disabled {
    cursor: not-allowed;
    opacity: 0.5;
  }

  @media (max-width: 560px) {
    .button-wide {
      margin-left: 0;
    }
  }
</style>
