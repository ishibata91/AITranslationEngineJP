<script lang="ts">
  import { createEventDispatcher } from "svelte";
  import type { FeatureTemplateScreenState } from "@application/ports/input/feature-template";

  export let state = undefined as unknown as FeatureTemplateScreenState;

  const dispatch = createEventDispatcher<{
    refresh: void;
    retry: void;
    select: string;
    updateQuery: string;
  }>();
</script>

<section class="panel">
  <header class="header">
    <div>
      <p class="eyebrow">Copyable Feature Template</p>
      <h1>Feature Template</h1>
    </div>
    <button type="button" class="refresh" disabled={state.loading} on:click={() => dispatch("refresh")}>
      {state.loading ? "Refreshing..." : "Refresh"}
    </button>
  </header>

  <div class="filters">
    <label for="query">Filter query</label>
    <input
      id="query"
      type="text"
      value={state.filters.query}
      placeholder="job / dictionary / persona"
      on:change={(event) => dispatch("updateQuery", event.currentTarget.value)}
    />
  </div>

  {#if state.error}
    <div class="status error">
      <p>{state.error}</p>
      <button type="button" on:click={() => dispatch("retry")}>Retry</button>
    </div>
  {/if}

  {#if state.loading && state.data === null}
    <p class="status loading">Loading template data...</p>
  {/if}

  {#if state.data}
    <ul class="grid">
      {#each state.data.items as item}
        <li>
          <button
            type="button"
            class:selected={state.selection === item.id}
            class="card"
            on:click={() => dispatch("select", item.id)}
          >
            <strong>{item.title}</strong>
            <span>{item.status}</span>
            <small>{item.detail}</small>
          </button>
        </li>
      {/each}
    </ul>
  {/if}
</section>

<style>
  .panel {
    display: grid;
    gap: 1rem;
    width: min(100%, 58rem);
    padding: 1.5rem;
    border: 1px solid #cfd8e3;
    border-radius: 1rem;
    background: #fff;
  }

  .header {
    display: flex;
    justify-content: space-between;
    gap: 1rem;
    align-items: start;
  }

  .eyebrow {
    margin: 0 0 0.4rem;
    font-size: 0.75rem;
    letter-spacing: 0.12em;
    text-transform: uppercase;
    color: #5d6b80;
  }

  h1 {
    margin: 0;
    font-size: 1.7rem;
  }

  .refresh {
    border: 0;
    border-radius: 999px;
    padding: 0.6rem 1rem;
    background: #142033;
    color: #fff;
    font: inherit;
    cursor: pointer;
  }

  .filters {
    display: grid;
    gap: 0.4rem;
  }

  input {
    border: 1px solid #c7d2df;
    border-radius: 0.7rem;
    padding: 0.55rem 0.7rem;
    font: inherit;
  }

  .status {
    margin: 0;
    padding: 0.8rem 0.9rem;
    border-radius: 0.7rem;
  }

  .loading {
    background: #ecf3ff;
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

  .grid {
    list-style: none;
    margin: 0;
    padding: 0;
    display: grid;
    gap: 0.7rem;
  }

  .card {
    width: 100%;
    text-align: left;
    border: 1px solid #d5dee9;
    border-radius: 0.7rem;
    padding: 0.8rem;
    background: #f8fafd;
    display: grid;
    gap: 0.35rem;
    font: inherit;
    color: inherit;
    cursor: pointer;
  }

  .card.selected {
    background: #edf4ff;
    border-color: #142033;
  }

  @media (max-width: 700px) {
    .header {
      align-items: stretch;
      flex-direction: column;
    }
  }
</style>
