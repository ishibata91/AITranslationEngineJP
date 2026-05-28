<script lang="ts">
  import type { TranslationOutputDiffPreviewSnapshot } from "@application/contract/translation-output-artifact/translation-output-artifact-screen-types"

  interface Props {
    compatibilitySummaryText: string
    diffPreview: TranslationOutputDiffPreviewSnapshot | null
    onSelectArtifact: (artifactId: number | null) => void
  }

  let { compatibilitySummaryText, diffPreview, onSelectArtifact }: Props =
    $props()
</script>

<section
  class="output-card diff-card"
  aria-labelledby="outputDiffHeading"
  data-testid="output-management-diff-preview"
>
  <div class="section-head">
    <div>
      <p class="eyebrow">diff preview</p>
      <h3 id="outputDiffHeading">translation unit 差分</h3>
    </div>
    <span class="mini-text">{compatibilitySummaryText}</span>
  </div>
  {#if diffPreview && diffPreview.rows.length > 0}
    <div class="diff-table">
      <div class="diff-table-head">
        <span>Source</span>
        <span>Dest</span>
        <span>Status</span>
        <span>row reflection summary</span>
      </div>
      {#each diffPreview.rows as row (row.rowDigest)}
        <button
          class="diff-row"
          data-row-id={row.rowDigest}
          data-testid="output-management-diff-row"
          onclick={() => onSelectArtifact(diffPreview?.artifactId ?? null)}
          type="button"
        >
          <div>
            <strong>{row.edid}</strong>
            <p>{row.sourceExcerpt || "-"}</p>
          </div>
          <div>
            <strong>{row.formId}</strong>
            <p>{row.destExcerpt || "-"}</p>
          </div>
          <div>
            <p>{row.internalOutputStatus}</p>
            <p>xTranslator: {row.xTranslatorStatus}</p>
          </div>
          <div>
            <p>{row.rowReflectionSummary}</p>
            <p>stale: {row.staleReason || "-"}</p>
            <p>re-output: {row.canRegenerate ? "可" : "不要"}</p>
          </div>
        </button>
      {/each}
    </div>
  {:else}
    <p class="empty-text">
      diff preview は未取得です。artifact がない場合、row count 0 の場合、または
      gateway 未接続の場合は一覧を表示しません。
    </p>
  {/if}
</section>

<style>
  .output-card {
    border: 1px solid var(--line);
    border-radius: var(--radius-md);
    background: rgba(25, 22, 20, 0.82);
    box-shadow: var(--shadow);
    color: var(--text);
    padding: 1.25rem;
  }

  .diff-card,
  .diff-table,
  .diff-row {
    display: grid;
    gap: 0.75rem;
  }

  .section-head {
    align-items: center;
    display: flex;
    gap: 0.75rem;
    justify-content: space-between;
  }

  .eyebrow,
  .mini-text {
    color: var(--muted);
    font-size: 0.85rem;
  }

  .empty-text,
  .diff-row p {
    line-height: 1.6;
  }

  .diff-row {
    width: 100%;
    text-align: left;
    border: 1px solid var(--line);
    border-radius: 14px;
    background: rgba(36, 31, 28, 0.92);
    color: inherit;
    padding: 0.95rem;
  }

  .diff-row:disabled {
    cursor: not-allowed;
    opacity: 0.55;
  }

  .diff-table-head,
  .diff-row {
    display: grid;
    gap: 0.75rem;
    grid-template-columns:
      minmax(0, 1.2fr) minmax(0, 1.2fr) minmax(0, 0.7fr)
      minmax(0, 1fr);
  }

  @media (max-width: 720px) {
    .diff-table-head,
    .diff-row {
      grid-template-columns: 1fr;
    }
  }
</style>
