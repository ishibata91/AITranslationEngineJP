<script lang="ts">
  import type { BodyTranslationFieldResultItem } from "@application/contract/body-translation-phase"

  interface Props {
    rows: BodyTranslationFieldResultItem[]
    pageIndex: number
    pageCount: number
    pageLabel: string
    onPreviousPage: () => void
    onNextPage: () => void
  }

  let {
    rows,
    pageIndex,
    pageCount,
    pageLabel,
    onPreviousPage,
    onNextPage
  }: Props = $props()
</script>

<section
  class="complete-card"
  aria-labelledby="completeResultHeading"
  data-testid="translation-complete-source-translation-region"
>
  <div class="hero-head">
    <div>
      <p class="eyebrow">field result</p>
      <h3 id="completeResultHeading">原文と訳文</h3>
    </div>
    <span class="mini-text">{pageLabel}</span>
  </div>

  {#if rows.length === 0}
    <p
      class="empty-text"
      data-testid="translation-complete-translation-result-empty-state"
    >
      表示できる翻訳結果がありません。
    </p>
  {:else}
    <div class="result-list">
      {#each rows as row (`${row.fieldId}-${row.fieldLabel}`)}
        <article
          class="result-row"
          data-testid="translation-complete-translation-result-row"
        >
          <div class="result-meta">
            <strong>{row.recordTypeLabel}</strong>
            <span>{row.fieldLabel}</span>
            <span>{row.outputStatus}</span>
          </div>
          <details open>
            <summary>原文</summary>
            <p>{row.sourceExcerpt}</p>
          </details>
          <details open>
            <summary>訳文</summary>
            <p>{row.translatedText}</p>
          </details>
        </article>
      {/each}
    </div>
  {/if}

  <div
    class="pager"
    aria-label="翻訳結果ページング"
    data-testid="translation-complete-translation-result-pagination"
  >
    <button disabled={pageIndex === 0} onclick={onPreviousPage} type="button">
      前へ
    </button>
    <span>{pageLabel}</span>
    <button
      disabled={pageIndex >= pageCount - 1}
      onclick={onNextPage}
      type="button"
    >
      次へ
    </button>
  </div>
</section>

<style>
  @import "./translation-complete-card.css";

  h3,
  p {
    margin: 0;
  }

  h3 {
    color: #fff6ea;
  }

  .empty-text {
    color: rgba(236, 223, 205, 0.76);
  }

  .mini-text {
    color: rgba(236, 223, 205, 0.76);
  }

  .result-list,
  .result-row {
    display: grid;
    gap: 1rem;
  }

  .result-row {
    background: rgba(18, 16, 15, 0.72);
    border: 1px solid rgba(255, 212, 165, 0.14);
    border-radius: 12px;
    padding: 1rem;
  }

  .result-meta,
  .pager {
    align-items: center;
    display: flex;
    flex-wrap: wrap;
    gap: 0.75rem;
  }

  .result-meta {
    justify-content: flex-start;
  }

  details {
    overflow-wrap: anywhere;
  }

  summary {
    color: rgba(236, 223, 205, 0.78);
    cursor: pointer;
    margin-bottom: 0.35rem;
  }

  .pager {
    justify-content: center;
  }

  .pager button {
    background: rgba(255, 255, 255, 0.05);
    border: 1px solid rgba(255, 212, 165, 0.22);
    border-radius: 0.8rem;
    color: #fef3e8;
    cursor: pointer;
    padding: 0.65rem 0.85rem;
  }

  .pager button:disabled {
    cursor: not-allowed;
    opacity: 0.55;
  }
</style>
