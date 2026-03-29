<script lang="ts">
  import { createEventDispatcher } from "svelte";
  import type {
    BootstrapStatusField,
    BootstrapStatusScreenState
  } from "@application/ports/input/bootstrap-status";

  export let state = undefined as unknown as BootstrapStatusScreenState;

  type StatusEntry = {
    key: BootstrapStatusField;
    label: string;
    value: string;
  };

  const dispatch = createEventDispatcher<{
    refresh: void;
    retry: void;
    select: BootstrapStatusField;
  }>();

  let statusEntries: StatusEntry[] = [];
  let selectedEntry: StatusEntry | null = null;

  $: statusEntries = state.data
    ? [
        {
          key: "backendVersion",
          label: "Backend Version",
          value: state.data.backendVersion
        },
        {
          key: "boundaryReady",
          label: "Boundary Ready",
          value: state.data.boundaryReady ? "true" : "false"
        },
        {
          key: "frontendEntry",
          label: "Frontend Entry",
          value: state.data.frontendEntry
        }
      ]
    : [];

  $: selectedEntry =
    state.selection === null ? null : statusEntries.find((entry) => entry.key === state.selection) ?? null;
</script>

<section class="panel">
  <div class="header">
    <div>
      <p class="eyebrow">Feature Screen Template</p>
      <h1>Bootstrap Status</h1>
    </div>
    <button type="button" class="refresh" disabled={state.loading} on:click={() => dispatch("refresh")}>
      {state.loading ? "Refreshing..." : "Refresh"}
    </button>
  </div>

  <p class="lede">
    This screen shows the standard split between `screen`, `view`, `store`, `usecase`, and `gateway`.
  </p>

  {#if state.error}
    <div class="status error">
      <p>{state.error}</p>
      <button type="button" on:click={() => dispatch("retry")}>Retry</button>
    </div>
  {/if}

  {#if state.loading && state.data === null}
    <p class="status loading">Loading bootstrap status...</p>
  {/if}

  {#if state.data}
    <div class="content">
      <div class="status-grid">
        {#each statusEntries as entry}
          <button
            type="button"
            class:selected={state.selection === entry.key}
            class="card"
            on:click={() => dispatch("select", entry.key)}
          >
            <span>{entry.label}</span>
            <strong>{entry.value}</strong>
          </button>
        {/each}
      </div>

      <aside class="selection">
        <p class="selection-label">Selection</p>
        {#if selectedEntry}
          <h2>{selectedEntry.label}</h2>
          <p>{selectedEntry.value}</p>
        {:else}
          <p>Select a field to inspect how screen-local selection state is held in the store.</p>
        {/if}
      </aside>
    </div>
  {/if}
</section>

<style>
  .panel {
    width: min(100%, 60rem);
    padding: 2rem;
    border: 1px solid #c8d2e0;
    border-radius: 1.25rem;
    background: rgba(255, 255, 255, 0.88);
    box-shadow: 0 24px 60px rgba(20, 32, 51, 0.12);
  }

  .header {
    display: flex;
    align-items: start;
    justify-content: space-between;
    gap: 1rem;
  }

  .eyebrow {
    margin: 0 0 0.75rem;
    font-size: 0.8rem;
    letter-spacing: 0.18em;
    text-transform: uppercase;
    color: #5b6a7d;
  }

  h1,
  h2,
  p {
    margin: 0;
  }

  h1 {
    font-size: clamp(2rem, 4vw, 3rem);
  }

  .lede {
    margin: 1rem 0 1.5rem;
    line-height: 1.6;
    color: #334155;
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
    margin-bottom: 1rem;
    padding: 1rem;
    border-radius: 0.9rem;
  }

  .loading {
    background: #edf4ff;
  }

  .error {
    display: grid;
    gap: 0.75rem;
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

  .status-grid {
    display: grid;
    gap: 1rem;
  }

  .card {
    display: grid;
    gap: 0.4rem;
    padding: 1rem;
    border: 1px solid #d7e0ec;
    border-radius: 0.9rem;
    background: #f8fafc;
    color: inherit;
    font: inherit;
    text-align: left;
    cursor: pointer;
  }

  .card.selected {
    border-color: #142033;
    background: #edf4ff;
  }

  .card span {
    font-size: 0.85rem;
    color: #526277;
  }

  .card strong {
    font-size: 1.1rem;
  }

  .selection {
    padding: 1rem;
    border-radius: 0.9rem;
    background: #f3f6fb;
  }

  .selection-label {
    margin-bottom: 0.5rem;
    font-size: 0.8rem;
    letter-spacing: 0.08em;
    text-transform: uppercase;
    color: #526277;
  }

  @media (min-width: 800px) {
    .content {
      grid-template-columns: minmax(0, 2fr) minmax(18rem, 1fr);
      align-items: start;
    }
  }
</style>
