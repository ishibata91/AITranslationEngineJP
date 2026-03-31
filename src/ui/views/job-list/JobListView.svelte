<script lang="ts">
  import { createEventDispatcher } from "svelte";
  import type { JobListItem, JobListScreenState } from "@application/usecases/job-list";

  export let state = undefined as unknown as JobListScreenState;

  const dispatch = createEventDispatcher<{
    refresh: void;
    retry: void;
    select: string;
  }>();

  let selectedJob: JobListItem | null = null;

  $: selectedJob =
    state.selection === null || state.data === null
      ? null
      : state.data.jobs.find((job) => job.jobId === state.selection) ?? null;
</script>

<section class="panel">
  <header class="header">
    <div>
      <p class="eyebrow">Phase 1 Observation</p>
      <h1>Job List</h1>
    </div>
    <button type="button" class="refresh" disabled={state.loading} on:click={() => dispatch("refresh")}>
      {state.loading ? "Refreshing..." : "Refresh"}
    </button>
  </header>

  {#if state.error}
    <div class="status error">
      <p>{state.error}</p>
      <button type="button" on:click={() => dispatch("retry")}>Retry</button>
    </div>
  {/if}

  <div class="content">
    <section class="list-panel">
      {#if state.loading && state.data === null}
        <p class="status loading">Loading jobs...</p>
      {:else if state.data && state.data.jobs.length === 0}
        <p class="status empty">No jobs available.</p>
      {:else if state.data}
        <ul>
          {#each state.data.jobs as job}
            <li>
              <button
                type="button"
                class:selected={state.selection === job.jobId}
                class="job-row"
                on:click={() => dispatch("select", job.jobId)}
              >
                <span>{job.jobId}</span>
                <strong>{job.state}</strong>
              </button>
            </li>
          {/each}
        </ul>
      {/if}
    </section>

    <aside class="summary-panel">
      <p class="summary-label">Selected Job</p>
      {#if selectedJob}
        <h2>{selectedJob.jobId}</h2>
        <p>State: {selectedJob.state}</p>
      {:else}
        <p>No job selected.</p>
      {/if}
    </aside>
  </div>
</section>

<style>
  .panel {
    width: min(100%, 62rem);
    display: grid;
    gap: 1rem;
    padding: 1.5rem;
    border: 1px solid #d7e1eb;
    border-radius: 1.25rem;
    background:
      linear-gradient(180deg, rgba(255, 255, 255, 0.96), rgba(245, 249, 253, 0.96)),
      #fff;
    box-shadow: 0 24px 60px rgba(20, 32, 51, 0.12);
  }

  .header {
    display: flex;
    align-items: start;
    justify-content: space-between;
    gap: 1rem;
  }

  .eyebrow {
    margin: 0 0 0.5rem;
    font-size: 0.8rem;
    letter-spacing: 0.14em;
    text-transform: uppercase;
    color: #5b6a7d;
  }

  h1,
  h2,
  p {
    margin: 0;
  }

  h1 {
    font-size: clamp(1.8rem, 4vw, 2.5rem);
  }

  .refresh {
    padding: 0.75rem 1rem;
    border: 0;
    border-radius: 999px;
    background: #142033;
    color: #fff;
    font: inherit;
    cursor: pointer;
  }

  .refresh:disabled {
    opacity: 0.7;
    cursor: progress;
  }

  .status {
    margin: 0;
    padding: 0.85rem 1rem;
    border-radius: 0.85rem;
  }

  .loading {
    background: #edf4ff;
  }

  .empty {
    background: #f1f5f9;
    color: #334155;
  }

  .error {
    display: grid;
    gap: 0.6rem;
    background: #fff1f2;
    color: #9f1239;
  }

  .error button {
    justify-self: start;
  }

  .content {
    display: grid;
    gap: 1rem;
  }

  .list-panel,
  .summary-panel {
    padding: 1rem;
    border-radius: 0.9rem;
    background: #f8fafc;
    border: 1px solid #d7e1eb;
  }

  .summary-label {
    margin-bottom: 0.4rem;
    font-size: 0.8rem;
    letter-spacing: 0.08em;
    text-transform: uppercase;
    color: #526277;
  }

  ul {
    list-style: none;
    margin: 0;
    padding: 0;
    display: grid;
    gap: 0.6rem;
  }

  .job-row {
    width: 100%;
    display: flex;
    align-items: center;
    justify-content: space-between;
    gap: 0.8rem;
    padding: 0.75rem 0.9rem;
    border: 1px solid #d5dee9;
    border-radius: 0.75rem;
    background: #fff;
    color: inherit;
    font: inherit;
    text-align: left;
    cursor: pointer;
  }

  .job-row.selected {
    background: #edf4ff;
    border-color: #142033;
  }

  @media (min-width: 860px) {
    .content {
      grid-template-columns: minmax(0, 2fr) minmax(18rem, 1fr);
      align-items: start;
    }
  }

  @media (max-width: 700px) {
    .header {
      flex-direction: column;
      align-items: stretch;
    }
  }
</style>
