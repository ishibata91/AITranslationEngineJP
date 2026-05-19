<script lang="ts">
  import { canOpenJobSetup } from "@application/presenter/translation-input"
  import type { TranslationInputReviewItem } from "@application/gateway-contract/translation-input"

  interface Props {
    selectedItem: TranslationInputReviewItem | null
    selectionStatusText: string
    latestOutcomeTitle: string
    latestOutcomeText: string
    canRebuildSelected: boolean
    isRebuilding: boolean
    formatDate: (timestamp: string) => string
    formatErrorKind: (errorKind: string | null) => string
    formatWarningKind: (kind: string) => string
    onRebuild: () => void | Promise<void>
  }

  let {
    selectedItem,
    selectionStatusText,
    latestOutcomeTitle,
    latestOutcomeText,
    canRebuildSelected,
    isRebuilding,
    formatDate,
    formatErrorKind,
    formatWarningKind,
    onRebuild
  }: Props = $props()

  const showNextStepGuidance = $derived(canOpenJobSetup(selectedItem))
</script>

<section
  class="panel detail-panel"
  aria-labelledby="inputReviewDetailHeading"
  data-testid="translation-input-review-selected-input-region"
>
  <div class="section-head">
    <div class="title-stack">
      <p class="eyebrow">選択データ</p>
      <h3 id="inputReviewDetailHeading">選択データの内容</h3>
    </div>
  </div>

  <p class="selection-note">{selectionStatusText}</p>
  <div class="result-callout">
    <strong>{latestOutcomeTitle}</strong>
    <p>{latestOutcomeText}</p>
  </div>

  {#if selectedItem}
    <div
      class="detail-stack"
      data-testid="translation-input-review-input-data-summary"
    >
      <dl class="summary-grid">
        <div>
          <dt>ファイル名</dt>
          <dd>{selectedItem.fileName}</dd>
        </div>
        <div>
          <dt>保存場所</dt>
          <dd>{selectedItem.filePath}</dd>
        </div>
        <div>
          <dt>内容ハッシュ</dt>
          <dd class="hash-text">{selectedItem.fileHash}</dd>
        </div>
        <div>
          <dt>読み込み日時</dt>
          <dd>{formatDate(selectedItem.importTimestamp)}</dd>
        </div>
        <div>
          <dt>翻訳レコード数</dt>
          <dd>{selectedItem.summary?.translationRecordCount ?? 0}</dd>
        </div>
        <div>
          <dt>翻訳フィールド数</dt>
          <dd>{selectedItem.summary?.translationFieldCount ?? 0}</dd>
        </div>
        <div>
          <dt>対象プラグイン</dt>
          <dd>{selectedItem.summary?.input.targetPluginName ?? "-"}</dd>
        </div>
        <div>
          <dt>抽出元ツール</dt>
          <dd>{selectedItem.summary?.input.sourceTool ?? "-"}</dd>
        </div>
      </dl>

      <section class="detail-section">
        <div class="section-head section-head-compact">
          <h4>カテゴリ別件数</h4>
          <span class="section-pill">
            {selectedItem.summary?.categories.length ?? 0} 件
          </span>
        </div>
        {#if (selectedItem.summary?.categories.length ?? 0) > 0}
          <div class="chip-grid">
            {#each selectedItem.summary?.categories ?? [] as category (`${category.category}:${category.recordCount}:${category.fieldCount}`)}
              <article class="chip-card">
                <strong>{category.category}</strong>
                <p>
                  レコード {category.recordCount} 件 / フィールド {category.fieldCount}
                  件
                </p>
                <span class="sr-only">
                  record {category.recordCount} / field {category.fieldCount}
                </span>
              </article>
            {/each}
          </div>
        {:else}
          <p class="empty-copy">カテゴリ別件数はまだありません。</p>
        {/if}
      </section>

      <section class="detail-section">
        <div class="section-head section-head-compact">
          <h4>サンプル項目</h4>
          <span class="section-pill">
            {selectedItem.summary?.sampleFields.length ?? 0} 件
          </span>
        </div>
        {#if (selectedItem.summary?.sampleFields.length ?? 0) > 0}
          <div class="sample-grid">
            {#each selectedItem.summary?.sampleFields ?? [] as field (`${field.recordType}:${field.subrecordType}:${field.formId}:${field.editorId}`)}
              <article class="sample-card">
                <div class="sample-head">
                  <strong>{field.recordType}:{field.subrecordType}</strong>
                  <span>
                    {field.translatable ? "翻訳対象" : "翻訳対象外"}
                  </span>
                </div>
                <p>{field.sourceText || "-"}</p>
                <dl>
                  <div>
                    <dt>Form ID</dt>
                    <dd>{field.formId || "-"}</dd>
                  </div>
                  <div>
                    <dt>Editor ID</dt>
                    <dd>{field.editorId || "-"}</dd>
                  </div>
                </dl>
              </article>
            {/each}
          </div>
        {:else}
          <div class="empty-copy-stack">
            <p class="empty-copy">サンプル項目はまだありません。</p>
            <span class="sr-only">sample field はまだありません。</span>
          </div>
        {/if}
      </section>

      <section
        class="detail-section"
        aria-labelledby="inputReviewIssueHeading"
        data-testid="translation-input-review-issue-rebuild-region"
      >
        <div class="section-head section-head-compact issue-head">
          <div class="title-stack">
            <h4 id="inputReviewIssueHeading">問題と再構築</h4>
            <p class="support-copy">
              登録エラーや警告を見て、再構築が必要かを判断する。
            </p>
          </div>
          <button
            class="button-secondary"
            disabled={!canRebuildSelected || isRebuilding}
            onclick={() => void onRebuild()}
            type="button"
          >
            cache を再構築
          </button>
        </div>

        <dl class="issue-summary">
          <div>
            <dt>登録時の問題区分</dt>
            <dd>{formatErrorKind(selectedItem.errorKind)}</dd>
          </div>
          <div>
            <dt>再構築の可否</dt>
            <dd>
              {selectedItem.canRebuild
                ? "再構築できます"
                : "再構築はまだできません"}
            </dd>
          </div>
        </dl>

        <div class="chip-grid">
          {#if selectedItem.errorKind}
            <article class="chip-card chip-card-error">
              <strong>{formatErrorKind(selectedItem.errorKind)}</strong>
              <p>登録または再構築で返された問題を確認してください。</p>
            </article>
          {/if}
          {#each selectedItem.warnings as warning (`${warning.kind}:${warning.recordType}:${warning.subrecordType}:${warning.message}`)}
            <article class="chip-card chip-card-warning">
              <strong>{formatWarningKind(warning.kind)}</strong>
              <p>{warning.message}</p>
            </article>
          {/each}
          {#if !selectedItem.errorKind && selectedItem.warnings.length === 0}
            <article class="chip-card">
              <strong>問題は見つかっていません</strong>
              <span class="sr-only">問題なし</span>
              <p>登録状態は安定しています。必要な時だけ再構築してください。</p>
            </article>
          {/if}
        </div>
      </section>

      {#if showNextStepGuidance}
        <section class="detail-section" aria-labelledby="inputReviewNextStepHeading">
          <div class="section-head section-head-compact issue-head">
            <div class="title-stack">
              <h4 id="inputReviewNextStepHeading">次の手順</h4>
              <p class="support-copy">
                入力登録だけでは Job Management には表示されません。Job Setup
                で job を作成してください。
              </p>
            </div>
          </div>
        </section>
      {/if}
    </div>
  {:else}
    <div class="empty-state">
      <p>読み込み済みデータを選ぶと、内容と再構築判断をここで確認できます。</p>
    </div>
  {/if}
</section>

<style>
  .panel {
    border: 1px solid var(--line);
    border-radius: 20px;
    padding: 1.25rem;
    background: rgba(28, 23, 20, 0.74);
    box-shadow: var(--shadow);
    color: var(--text);
    backdrop-filter: blur(18px);
  }

  .detail-panel,
  .detail-stack,
  .detail-section {
    display: grid;
    gap: 1rem;
  }

  .section-head,
  .sample-head {
    display: flex;
    justify-content: space-between;
    gap: 1rem;
    align-items: start;
  }

  .section-head-compact {
    margin-bottom: 0;
  }

  .issue-head {
    align-items: center;
  }

  .title-stack {
    display: grid;
    gap: 0.35rem;
  }

  .eyebrow,
  .support-copy,
  .selection-note,
  .result-callout p,
  dt,
  .empty-copy,
  .sample-card p {
    color: var(--muted);
  }

  .eyebrow {
    margin: 0;
    font-size: 0.76rem;
    letter-spacing: 0.16em;
    text-transform: uppercase;
  }

  h3,
  h4,
  p,
  dl,
  dt,
  dd {
    margin: 0;
  }

  .selection-note {
    font-size: 0.95rem;
  }

  .result-callout,
  .chip-card,
  .sample-card,
  .empty-state {
    border-radius: 18px;
    padding: 1rem;
    background: rgba(18, 15, 14, 0.76);
    border: 1px solid rgba(255, 186, 56, 0.12);
  }

  .summary-grid,
  .issue-summary {
    display: grid;
    gap: 0.8rem;
    grid-template-columns: repeat(auto-fit, minmax(180px, 1fr));
  }

  .summary-grid div,
  .issue-summary div,
  .sample-card dl div {
    display: grid;
    gap: 0.2rem;
  }

  .hash-text {
    word-break: break-all;
  }

  .section-pill {
    padding: 0.32rem 0.7rem;
    border: 1px solid var(--line-strong);
    border-radius: 999px;
    color: var(--primary);
    font-size: 0.82rem;
    white-space: nowrap;
  }

  .chip-grid,
  .sample-grid {
    display: grid;
    gap: 0.8rem;
  }

  .empty-copy-stack {
    display: grid;
    gap: 0.2rem;
  }

  .sample-grid {
    grid-template-columns: repeat(auto-fit, minmax(220px, 1fr));
  }

  .chip-card-error {
    border-color: rgba(255, 104, 63, 0.4);
  }

  .chip-card-warning {
    border-color: rgba(255, 186, 56, 0.34);
  }

  .sample-card {
    display: grid;
    gap: 0.65rem;
  }

  .sample-card dl {
    display: grid;
    gap: 0.5rem;
  }

  .button-secondary {
    font: inherit;
    border-radius: 999px;
    padding: 0.7rem 1.1rem;
    border: 1px solid var(--line-strong);
    background: transparent;
    color: var(--text);
    cursor: pointer;
  }

  .button-secondary:disabled {
    cursor: not-allowed;
    opacity: 0.55;
  }

  .sr-only {
    position: absolute;
    width: 1px;
    height: 1px;
    padding: 0;
    margin: -1px;
    overflow: hidden;
    clip: rect(0, 0, 0, 0);
    white-space: nowrap;
    border: 0;
  }

  @media (max-width: 720px) {
    .section-head,
    .sample-head,
    .issue-head {
      flex-direction: column;
      align-items: start;
    }

    .section-pill {
      white-space: normal;
    }
  }
</style>
