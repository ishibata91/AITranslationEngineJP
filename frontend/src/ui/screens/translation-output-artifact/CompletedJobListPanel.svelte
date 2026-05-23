<script lang="ts">
  import type { TranslationOutputCompletedJobSummary } from "@application/gateway-contract/translation-output-artifact"

  import {
    formatCount,
    formatDistribution,
    formatStatus
  } from "./output-artifact-formatters"

  interface Props {
    completedJobs: TranslationOutputCompletedJobSummary[]
    selectedJobId: number | null
    refreshDisabled: boolean
    onRefresh: () => void
    onSelectJob: (job: TranslationOutputCompletedJobSummary) => void
  }

  let {
    completedJobs,
    selectedJobId,
    refreshDisabled,
    onRefresh,
    onSelectJob
  }: Props = $props()
</script>

<section
  class="output-card job-list-card"
  aria-labelledby="outputJobListHeading"
  data-testid="output-management-output-candidate-list"
>
  <div class="section-head">
    <div>
      <p class="eyebrow">completed job list</p>
      <h3 id="outputJobListHeading">出力候補</h3>
    </div>
    <button
      class="secondary-button"
      disabled={refreshDisabled}
      onclick={onRefresh}
      type="button"
    >
      更新
    </button>
  </div>
  {#if completedJobs.length === 0}
    <p class="empty-text">
      completed job はありません。target count 0 の job も一覧に出ません。
    </p>
  {:else}
    <div class="job-list">
      {#each completedJobs as job (job.jobId)}
        <button
          class="job-button"
          class:is-selected={job.jobId === selectedJobId}
          onclick={() => onSelectJob(job)}
          type="button"
        >
          <div class="job-button-head">
            <strong>job #{job.jobId}</strong>
            <span class="status-pill">{formatStatus(job.artifactStatus)}</span>
          </div>
          <p>{formatStatus(job.jobStatus)}</p>
          <p>translated: {formatCount(job.translatedCount)}</p>
          <p>status: {formatDistribution(job.outputStatusDistribution)}</p>
        </button>
      {/each}
    </div>
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

  .section-head,
  .job-button-head {
    align-items: center;
    display: flex;
    gap: 0.75rem;
    justify-content: space-between;
  }

  .eyebrow {
    color: var(--muted);
    font-size: 0.85rem;
  }

  .empty-text,
  .job-button p {
    line-height: 1.6;
  }

  .job-list {
    display: grid;
    gap: 0.75rem;
  }

  .job-button {
    width: 100%;
    text-align: left;
    border: 1px solid var(--line);
    border-radius: 14px;
    background: rgba(36, 31, 28, 0.92);
    color: inherit;
    cursor: pointer;
    font: inherit;
    padding: 0.95rem;
  }

  .job-button.is-selected {
    border-color: var(--primary);
    box-shadow: 0 0 0 1px rgba(255, 186, 56, 0.32);
  }

  .status-pill {
    border: 1px solid var(--line-strong);
    border-radius: 999px;
    padding: 0.2rem 0.65rem;
    color: var(--primary);
    font-size: 0.82rem;
  }

  .secondary-button {
    border: 1px solid var(--line);
    border-radius: 12px;
    background: rgba(15, 13, 12, 0.9);
    color: inherit;
    cursor: pointer;
    font: inherit;
    padding: 0.8rem 0.95rem;
  }

  .secondary-button:disabled,
  .job-button:disabled {
    cursor: not-allowed;
    opacity: 0.55;
  }
</style>
