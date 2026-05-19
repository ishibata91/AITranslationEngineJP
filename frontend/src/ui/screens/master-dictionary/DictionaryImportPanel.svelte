<script lang="ts">
  import type { DictionaryImportPanelProps } from "./dictionary-panel-props"

  let {
    hasStagedFile,
    importProgress,
    importStatusText,
    importStatusValue,
    importSummary,
    isImportRunning,
    selectedEntry,
    selectedFileName,
    chooseXmlFile,
    handleXmlSelected,
    resetImportSelection,
    startImport
  }: DictionaryImportPanelProps = $props()
</script>

<section
  class="shell-card import-shell"
  aria-labelledby="importHeading"
  data-testid="master-dictionary-xml-import-region"
>
  <div class="import-top">
    <div>
      <p class="eyebrow">XMLから取り込み</p>
      <h3 id="importHeading">取り込み導線</h3>
    </div>
    <button
      class="button-secondary"
      id="chooseXmlButton"
      onclick={chooseXmlFile}
      type="button"
    >
      ファイルを選択
    </button>
  </div>
  <p class="mini-text" id="importStateText">
    ファイルを選ぶと取込バーが表示されます。
  </p>

  <input
    accept=".xml,text/xml,application/xml"
    class="file-input"
    id="xmlFileInput"
    onchange={handleXmlSelected}
    type="file"
  />

  <div class="file-picker">
    <span class="eyebrow">選択ファイル</span>
    <span class="file-name" id="selectedFileName">{selectedFileName}</span>
  </div>

  <div
    class="import-bar"
    data-testid="master-dictionary-xml-import-bar"
    hidden={!hasStagedFile}
    id="importBar"
  >
    <div class="import-bar-head">
      <strong id="importFileTitle">{selectedFileName}</strong>
      <div class="import-actions">
        <button
          class="button-primary"
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
    </div>
    <div class="status-line">
      <p id="importStatusText">{importStatusText}</p>
      <strong id="importStatusValue">{importStatusValue}</strong>
    </div>
    <div class="progress-track">
      <div
        class="progress-fill"
        id="importProgressFill"
        style={`width: ${importProgress}%;`}
      ></div>
    </div>
    <div class="import-result" hidden={!importSummary} id="importResult">
      <div class="import-result-head">
        <strong id="importResultHeadline"
          >XML取り込みを一覧と詳細へ反映しました。</strong
        >
        <span class="status-pill" id="importResultCount"
          >新規追加 {importSummary?.importedCount ?? 0} 件</span
        >
      </div>
      <p id="importResultMessage">
        {importSummary
          ? `「${importSummary.fileName}」の取込を完了し、同じ画面に反映しました。件数は保存済みエントリ単位で集計しています。`
          : "-"}
      </p>
      <dl class="result-grid">
        <div>
          <dt>更新済みエントリ件数</dt>
          <dd id="importResultUpdatedCount">
            {importSummary?.updatedCount ?? "-"}
          </dd>
        </div>
        <div>
          <dt>取込後の保存済み一覧件数</dt>
          <dd id="importResultListCount">{importSummary?.totalCount ?? "-"}</dd>
        </div>
        <div>
          <dt>選択状態</dt>
          <dd id="importResultSelection">
            {importSummary?.selectedSource ?? "-"}
          </dd>
        </div>
        <div>
          <dt>詳細表示</dt>
          <dd id="importResultDetail">{selectedEntry?.translation ?? "-"}</dd>
        </div>
      </dl>
    </div>
  </div>
</section>

<style>
  .shell-card {
    padding: 18px;
    border-radius: 16px;
    border: 1px solid var(--line);
    background: rgba(16, 13, 11, 0.58);
    color: var(--text);
  }

  .import-shell,
  .import-bar,
  .result-grid {
    display: grid;
    gap: 10px;
  }

  .import-top,
  .import-bar-head,
  .import-actions,
  .import-result-head {
    display: flex;
    flex-wrap: wrap;
    justify-content: space-between;
    align-items: center;
    gap: 10px;
  }

  .mini-text,
  #importStatusText,
  #importResultMessage,
  dt {
    color: var(--muted);
  }

  .eyebrow,
  dt {
    font-size: 12px;
    letter-spacing: 0.08em;
  }

  .button-primary,
  .button-secondary {
    min-height: 36px;
    padding: 0 14px;
    border-radius: 999px;
    border: 1px solid transparent;
    font: inherit;
  }

  .button-primary {
    color: #3a2400;
    background: linear-gradient(135deg, var(--primary) 0%, #ef9d20 100%);
  }

  .button-secondary {
    color: var(--text);
    background: rgba(255, 255, 255, 0.04);
    border-color: var(--line);
  }

  .import-bar[hidden],
  .import-result[hidden] {
    display: none !important;
    pointer-events: none;
  }

  .file-picker {
    display: inline-flex;
    gap: 10px;
    align-items: center;
  }

  .file-input {
    position: absolute;
    width: 1px;
    height: 1px;
    margin: -1px;
    padding: 0;
    border: 0;
    overflow: hidden;
    clip: rect(0, 0, 0, 0);
    clip-path: inset(50%);
    white-space: nowrap;
    pointer-events: none;
  }

  .file-name,
  .status-pill {
    padding: 6px 10px;
    border-radius: 999px;
    border: 1px solid var(--line);
    background: rgba(255, 255, 255, 0.03);
  }

  .progress-track {
    height: 10px;
    border-radius: 999px;
    background: rgba(255, 255, 255, 0.08);
    overflow: hidden;
  }

  .progress-fill {
    height: 100%;
    background: linear-gradient(90deg, var(--primary) 0%, #f5ca72 100%);
    transition: width 180ms ease;
  }
</style>
