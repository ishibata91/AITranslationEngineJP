<script lang="ts">
  import type { TranslationOutputReviewSnapshot } from "@application/contract/translation-output-artifact/translation-output-artifact-screen-types"

  import { formatCount, formatStatus } from "./output-artifact-formatters"

  interface Props {
    review: TranslationOutputReviewSnapshot | null
    viewState: string
  }

  let { review, viewState }: Props = $props()
</script>

<section
  class="output-card summary-card"
  aria-labelledby="outputSummaryHeading"
  data-testid="output-management-selected-job"
>
  <div class="section-head">
    <div>
      <p class="eyebrow">selected job summary</p>
      <h3 id="outputSummaryHeading">
        {review ? "選択中 job" : "job 選択待ち"}
      </h3>
    </div>
    <span class="status-pill" data-view-state={viewState}>
      {formatStatus(viewState)}
    </span>
  </div>
  {#if review}
    <dl class="detail-grid">
      <div>
        <dt>job id</dt>
        <dd>{review.selectedJobId}</dd>
      </div>
      <div>
        <dt>job status</dt>
        <dd>{formatStatus(review.selectedJobStatus)}</dd>
      </div>
      <div>
        <dt>body phase</dt>
        <dd>{formatStatus(review.bodyPhaseStatus)}</dd>
      </div>
      <div>
        <dt>readiness</dt>
        <dd>{review.readiness ? "ready" : "not ready"}</dd>
      </div>
      <div>
        <dt>translated count</dt>
        <dd>{formatCount(review.translatedCount)}</dd>
      </div>
      <div>
        <dt>row count</dt>
        <dd>{formatCount(review.rowCount)}</dd>
      </div>
      <div>
        <dt>artifact status</dt>
        <dd>{formatStatus(review.artifactStatus)}</dd>
      </div>
      <div>
        <dt>current version</dt>
        <dd>{review.currentVersion ? "yes" : "stale"}</dd>
      </div>
      <div>
        <dt>input provenance</dt>
        <dd class="wrap-value">
          {review.inputSnapshotDigest} / {review.sourceFileDigest}
        </dd>
      </div>
    </dl>
    {#if review.rejectionReasons.length > 0}
      <div class="notice-block warning">
        <h4>拒否理由</h4>
        <ul>
          {#each review.rejectionReasons as reason (`${reason.errorKind}-${reason.reason}`)}
            <li>{reason.errorKind}: {reason.reason}</li>
          {/each}
        </ul>
      </div>
    {/if}
  {:else}
    <p class="empty-text">
      出力候補から completed job を選ぶと、summary と出力準備を表示します。
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

  .section-head {
    align-items: center;
    display: flex;
    gap: 0.75rem;
    justify-content: space-between;
  }

  .eyebrow {
    color: var(--muted);
    font-size: 0.85rem;
  }

  .empty-text {
    line-height: 1.6;
  }

  .detail-grid,
  .notice-block {
    display: grid;
    gap: 0.75rem;
  }

  .detail-grid {
    grid-template-columns: repeat(2, minmax(0, 1fr));
  }

  dt {
    color: var(--muted);
    font-size: 0.82rem;
    margin-bottom: 0.2rem;
  }

  dd {
    margin: 0;
  }

  .wrap-value {
    overflow-wrap: anywhere;
  }

  .status-pill {
    border: 1px solid var(--line-strong);
    border-radius: 999px;
    padding: 0.2rem 0.65rem;
    color: var(--primary);
    font-size: 0.82rem;
  }

  .status-pill[data-view-state="failed"],
  .notice-block.warning {
    color: #ff9f7f;
  }

  .notice-block {
    border: 1px solid var(--line);
    border-radius: 12px;
    padding: 0.85rem;
    background: rgba(18, 16, 15, 0.88);
    color: var(--text);
  }

  .notice-block h4 {
    font-weight: 700;
  }

  .notice-block ul {
    margin: 0;
    padding-left: 1.2rem;
  }

  @media (max-width: 720px) {
    .detail-grid {
      grid-template-columns: 1fr;
    }
  }
</style>
