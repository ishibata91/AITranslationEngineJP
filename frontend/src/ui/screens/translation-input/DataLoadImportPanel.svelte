<script lang="ts">
  interface Props {
    stagedFileName: string
    stagedFilePath: string
    stagedFileHash: string
    hasStagedFile: boolean
    canImport: boolean
    isImporting: boolean
    onChooseJson: () => void
    onStartImport: () => void | Promise<void>
    onResetSelection: () => void
  }

  let {
    stagedFileName,
    stagedFilePath,
    stagedFileHash,
    hasStagedFile,
    canImport,
    isImporting,
    onChooseJson,
    onStartImport,
    onResetSelection
  }: Props = $props()
</script>

<section class="panel import-panel" aria-labelledby="inputReviewImportHeading">
  <div class="section-head">
    <div class="title-stack">
      <p class="eyebrow">ロード準備</p>
      <h3 id="inputReviewImportHeading">ロード準備</h3>
      <p class="support-copy">
        取り込みたい JSON を 1 件選び、内容を確認してから登録する。
      </p>
    </div>
    <button class="button-secondary" onclick={onChooseJson} type="button">
      ロード対象を選ぶ
    </button>
  </div>

  <div class="prep-grid">
    <article class="prep-card">
      <span class="prep-label">選択状態</span>
      <strong>{hasStagedFile ? "登録前の JSON を確認中" : "まだ選択していません"}</strong>
      <p>
        {hasStagedFile
          ? "選択した内容を確認してから登録できます。"
          : "まずは JSON を選ぶと、登録前の情報がここに表示されます。"}
      </p>
    </article>

    <dl class="file-grid">
      <div>
        <dt>ファイル名</dt>
        <dd>{stagedFileName}</dd>
      </div>
      <div>
        <dt>保存場所</dt>
        <dd>{stagedFilePath}</dd>
      </div>
      <div>
        <dt>内容ハッシュ</dt>
        <dd class="hash-text">{stagedFileHash}</dd>
      </div>
    </dl>
  </div>

  <div class="action-row">
    <button
      class="button-primary"
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
  </div>
</section>

<style>
  .panel {
    border: 1px solid var(--line);
    border-radius: 20px;
    padding: 1.25rem;
    background: rgba(28, 23, 20, 0.74);
    box-shadow: var(--shadow);
    backdrop-filter: blur(18px);
  }

  .import-panel {
    display: grid;
    gap: 1rem;
  }

  .section-head {
    display: flex;
    justify-content: space-between;
    gap: 1rem;
    align-items: start;
  }

  .title-stack {
    display: grid;
    gap: 0.35rem;
  }

  .eyebrow,
  .support-copy,
  dt,
  .prep-label {
    color: var(--muted);
  }

  .eyebrow {
    margin: 0;
    font-size: 0.76rem;
    letter-spacing: 0.16em;
    text-transform: uppercase;
  }

  h3,
  p,
  dl,
  dt,
  dd {
    margin: 0;
  }

  .prep-grid {
    display: grid;
    gap: 0.9rem;
    grid-template-columns: minmax(0, 0.95fr) minmax(0, 1.05fr);
  }

  .prep-card,
  .file-grid div {
    display: grid;
    gap: 0.3rem;
    padding: 1rem;
    border-radius: 18px;
    background: rgba(18, 15, 14, 0.76);
    border: 1px solid rgba(255, 186, 56, 0.12);
  }

  .file-grid {
    display: grid;
    gap: 0.8rem;
    grid-template-columns: repeat(auto-fit, minmax(180px, 1fr));
  }

  dt {
    font-size: 0.8rem;
  }

  dd {
    font-size: 0.95rem;
  }

  .hash-text {
    word-break: break-all;
  }

  .action-row {
    display: flex;
    gap: 0.75rem;
    flex-wrap: wrap;
  }

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

  @media (max-width: 860px) {
    .prep-grid {
      grid-template-columns: 1fr;
    }

    .section-head {
      flex-direction: column;
    }
  }
</style>
