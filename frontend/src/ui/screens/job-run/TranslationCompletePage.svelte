<script lang="ts">
  import type { BodyTranslationFieldResultItem } from "@application/contract/body-translation-phase"

  interface Props {
    jobId: number
    rows: BodyTranslationFieldResultItem[]
  }

  let { jobId, rows }: Props = $props()
  let pageIndex = $state(0)

  const pageSize = 5
  const pageCount = $derived(Math.max(1, Math.ceil(rows.length / pageSize)))
  const currentPageRows = $derived(
    rows.slice(pageIndex * pageSize, pageIndex * pageSize + pageSize)
  )
  const pageLabel = $derived(`${pageIndex + 1} / ${pageCount}`)

  $effect(() => {
    if (pageIndex > pageCount - 1) {
      pageIndex = pageCount - 1
    }
  })

  function previousPage(): void {
    pageIndex = Math.max(0, pageIndex - 1)
  }

  function nextPage(): void {
    pageIndex = Math.min(pageCount - 1, pageIndex + 1)
  }
</script>

<section
  class="complete-shell"
  data-testid="translation-complete-translation-complete-screen"
  id="translationCompleteView"
>
  <section
    class="complete-card hero-card"
    data-testid="translation-complete-completion-summary"
  >
    <div class="hero-head">
      <div>
        <p class="eyebrow">translation complete</p>
        <h2>翻訳完了</h2>
      </div>
      <span class="status-pill">ジョブ #{jobId}</span>
    </div>
    <p class="lead">
      本文翻訳が完了したジョブの原文と訳文を確認します。成果物の出力処理は出力管理で扱います。
    </p>
  </section>

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

    {#if currentPageRows.length === 0}
      <p
        class="empty-text"
        data-testid="translation-complete-translation-result-empty-state"
      >
        表示できる翻訳結果がありません。
      </p>
    {:else}
      <div class="result-list">
        {#each currentPageRows as row (`${row.fieldId}-${row.fieldLabel}`)}
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
      <button disabled={pageIndex === 0} onclick={previousPage} type="button">
        前へ
      </button>
      <span>{pageLabel}</span>
      <button
        disabled={pageIndex >= pageCount - 1}
        onclick={nextPage}
        type="button"
      >
        次へ
      </button>
    </div>
  </section>
</section>

<style>
  .complete-shell,
  .complete-card,
  .result-list,
  .result-row {
    display: grid;
    gap: 1rem;
  }

  .complete-card {
    border: 1px solid rgba(226, 205, 173, 0.14);
    border-radius: 16px;
    background: rgba(33, 27, 24, 0.88);
    padding: 1.3rem;
  }

  .hero-head,
  .pager,
  .result-meta {
    align-items: center;
    display: flex;
    flex-wrap: wrap;
    gap: 0.75rem;
    justify-content: space-between;
  }

  .eyebrow,
  .mini-text,
  .empty-text,
  .lead {
    color: rgba(236, 223, 205, 0.76);
  }

  h2,
  h3,
  p {
    margin: 0;
  }

  .status-pill {
    border: 1px solid rgba(255, 212, 165, 0.22);
    border-radius: 999px;
    color: #ffd8ae;
    padding: 0.25rem 0.65rem;
  }

  .result-row {
    border: 1px solid rgba(255, 212, 165, 0.14);
    border-radius: 12px;
    background: rgba(18, 16, 15, 0.72);
    padding: 1rem;
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
    border: 1px solid rgba(255, 212, 165, 0.22);
    border-radius: 0.8rem;
    background: rgba(255, 255, 255, 0.05);
    color: #fef3e8;
    cursor: pointer;
    padding: 0.65rem 0.85rem;
  }

  .pager button:disabled {
    cursor: not-allowed;
    opacity: 0.55;
  }
</style>
