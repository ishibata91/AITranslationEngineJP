<script lang="ts">
  import type { DictionaryImportProgressPanelProps } from "./dictionary-panel-props"

  let {
    hasStagedFile,
    importProgress,
    importStatusText,
    importStatusValue,
    importSummary,
    isImportRunning,
    selectedEntry,
    selectedFileName,
    resetImportSelection,
    startImport
  }: DictionaryImportProgressPanelProps = $props()

  const resultStats = $derived([
    {
      label: "更新済みエントリ件数",
      value: importSummary?.updatedCount ?? "-",
      valueId: "importResultUpdatedCount"
    },
    {
      label: "取込後の保存済み一覧件数",
      value: importSummary?.totalCount ?? "-",
      valueId: "importResultListCount"
    },
    {
      label: "選択状態",
      value: importSummary?.selectedSource ?? "-",
      valueId: "importResultSelection"
    },
    {
      label: "選択中の訳語",
      value: selectedEntry?.translation ?? "-",
      valueId: "importResultDetail"
    }
  ])
</script>

<section
  class="phase-card import-progress-panel"
  aria-labelledby="importProgressHeading"
  data-testid="master-dictionary-import-progress-panel"
>
  <div class="panel-head">
    <div>
      <p class="eyebrow">進行状況</p>
      <h3 id="importProgressHeading">XML 取り込みの進行状況</h3>
    </div>
    <span class="state-pill" id="importStatusValue">{importStatusValue}</span>
  </div>

  <p class="support-copy" id="importStatusText">{importStatusText}</p>

  <div class="progress-region" aria-label="XML 取り込みの進行率">
    <div
      class="progress-fill"
      id="importProgressFill"
      style={`width: ${importProgress}%;`}
    ></div>
  </div>

  <div class="current-file">
    <span class="eyebrow">対象ファイル</span>
    <strong>{selectedFileName}</strong>
  </div>

  {#if hasStagedFile}
    <div class="button-row">
      <button
        class="button-primary"
        data-testid="master-dictionary-xml-import-button"
        disabled={isImportRunning}
        id="startImportButton"
        onclick={startImport}
        type="button"
      >
        この XML を取り込む
      </button>
      <button
        class="button-secondary"
        disabled={isImportRunning}
        id="resetImportButton"
        onclick={resetImportSelection}
        type="button"
      >
        選び直す
      </button>
    </div>
  {/if}

  {#if importSummary}
    <section class="result-panel" aria-labelledby="importResultHeadline">
      <div class="result-head">
        <strong id="importResultHeadline">
          XML取り込みを一覧へ反映しました。
        </strong>
        <span class="result-badge">
          新規追加 {importSummary.importedCount} 件
        </span>
      </div>
      <p id="importResultMessage">
        「{importSummary.fileName}」の取込を完了し、同じ画面に反映しました。件数は保存済みエントリ単位で集計しています。
      </p>
      <dl class="result-grid">
        {#each resultStats as stat (stat.label)}
          <div>
            <dt>{stat.label}</dt>
            <dd id={stat.valueId}>{stat.value}</dd>
          </div>
        {/each}
      </dl>
    </section>
  {/if}
</section>

<style>
  .phase-card {
    align-content: start;
    background: rgba(33, 27, 24, 0.88);
    border: 1px solid rgba(226, 205, 173, 0.14);
    border-radius: 20px;
    box-shadow: 0 18px 40px rgba(0, 0, 0, 0.22);
    color: var(--text);
    display: grid;
    gap: 1rem;
    min-width: 0;
    padding: 1.4rem;
  }

  .panel-head,
  .button-row,
  .result-head {
    align-items: center;
    display: flex;
    flex-wrap: wrap;
    gap: 0.75rem;
    justify-content: space-between;
  }

  .eyebrow,
  dt {
    color: rgba(236, 223, 205, 0.72);
    font-size: 12px;
    letter-spacing: 0.08em;
    margin: 0;
  }

  h3,
  p,
  dl,
  dd {
    margin: 0;
  }

  h3,
  strong,
  dd {
    overflow-wrap: anywhere;
  }

  .support-copy,
  #importResultMessage {
    color: rgba(236, 223, 205, 0.8);
    line-height: 1.7;
  }

  .state-pill,
  .result-badge {
    background: rgba(198, 155, 82, 0.16);
    border: 1px solid rgba(226, 205, 173, 0.16);
    border-radius: 999px;
    color: #ffe2ae;
    display: inline-flex;
    padding: 0.35rem 0.75rem;
    white-space: nowrap;
  }

  .progress-region {
    background: rgba(255, 255, 255, 0.08);
    border-radius: 999px;
    height: 10px;
    overflow: hidden;
  }

  .progress-fill {
    background: linear-gradient(90deg, var(--primary) 0%, #f5ca72 100%);
    height: 100%;
    transition: width 180ms ease;
  }

  .current-file {
    background: rgba(255, 255, 255, 0.03);
    border: 1px solid var(--line);
    border-radius: 14px;
    display: grid;
    gap: 0.35rem;
    padding: 0.85rem;
  }

  .button-primary,
  .button-secondary {
    border-radius: 999px;
    cursor: pointer;
    font: inherit;
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

  button:disabled {
    cursor: not-allowed;
    opacity: 0.5;
  }

  .result-panel {
    background: rgba(0, 0, 0, 0.16);
    border: 1px solid rgba(255, 255, 255, 0.08);
    border-radius: 16px;
    display: grid;
    gap: 0.85rem;
    padding: 1rem;
  }

  .result-grid {
    display: grid;
    gap: 10px;
    grid-template-columns: repeat(auto-fit, minmax(140px, 1fr));
  }

  .result-grid div {
    background: rgba(255, 255, 255, 0.04);
    border: 0.5px solid rgba(255, 255, 255, 0.08);
    border-radius: 14px;
    display: grid;
    gap: 6px;
    min-width: 0;
    padding: 12px;
  }
</style>
