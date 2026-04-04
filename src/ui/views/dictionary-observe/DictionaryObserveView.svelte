<script lang="ts">
  import { createEventDispatcher } from "svelte";
  import type {
    DictionaryCandidateGroup,
    DictionaryObserveScreenState,
  } from "@application/usecases/dictionary-observe";

  export let state = undefined as unknown as DictionaryObserveScreenState;

  const dispatch = createEventDispatcher<{
    observe: void;
    refresh: void;
    retry: void;
    select: number;
    updateSourceTexts: string[];
  }>();

  let selectedGroup: DictionaryCandidateGroup | null = null;

  $: selectedGroup =
    state.selection !== null &&
    state.data !== null &&
    state.selection >= 0 &&
    state.selection < state.data.candidateGroups.length
      ? state.data.candidateGroups[state.selection]
      : null;

  function parseSourceTexts(input: string): string[] {
    return input
      .split(/\r?\n/u)
      .filter((sourceText) => sourceText.length > 0);
  }
</script>

<section class="panel">
  <header class="header">
    <div>
      <p class="eyebrow">Foundation Observation</p>
      <h1>Dictionary Observe</h1>
    </div>
    <div class="actions">
      <button
        type="button"
        class="observe"
        disabled={state.loading || state.filters.sourceTexts.length === 0}
        on:click={() => dispatch("observe")}
      >
        {state.loading ? "Observing..." : "Observe"}
      </button>
      <button
        type="button"
        class="refresh"
        disabled={state.loading || state.filters.lastSubmittedRequest === null}
        on:click={() => dispatch("refresh")}
      >
        Refresh
      </button>
    </div>
  </header>

  <label class="input-panel">
    <span>Source Texts</span>
    <textarea
      rows="5"
      value={state.filters.sourceTexts.join("\n")}
      disabled={state.loading}
      on:input={(event) => {
        const target = event.currentTarget as HTMLTextAreaElement;
        dispatch("updateSourceTexts", parseSourceTexts(target.value));
      }}
    ></textarea>
  </label>

  {#if state.error}
    <div class="status error">
      <p>{state.error}</p>
      <button type="button" on:click={() => dispatch("retry")}>Retry</button>
    </div>
  {/if}

  <div class="content">
    <section class="list-panel">
      {#if state.loading && state.data === null}
        <p class="status loading">Observing dictionary...</p>
      {:else if state.data === null}
        <p class="status empty">Run an observation to inspect dictionary candidates.</p>
      {:else if state.data.candidateGroups.length === 0}
        <p class="status empty">No observed requests.</p>
      {:else}
        <ul>
          {#each state.data.candidateGroups as group, index}
            <li>
              <button
                type="button"
                class:selected={state.selection === index}
                class="request-row"
                on:click={() => dispatch("select", index)}
              >
                <span>{group.sourceText}</span>
                <strong>{group.candidates.length} candidates</strong>
              </button>
            </li>
          {/each}
        </ul>
      {/if}
    </section>

    <aside class="summary-panel">
      <p class="summary-label">Selected Request</p>
      {#if selectedGroup}
        <h2>{selectedGroup.sourceText}</h2>
        {#if selectedGroup.candidates.length === 0}
          <p>No candidates found for this request.</p>
        {:else}
          <ul class="candidate-list">
            {#each selectedGroup.candidates as candidate}
              <li>
                <strong>{candidate.sourceText}</strong>
                <span>{candidate.destText}</span>
              </li>
            {/each}
          </ul>
        {/if}
      {:else}
        <p>No request selected.</p>
      {/if}
    </aside>
  </div>
</section>

<style>
  .panel {
    width: min(100%, 68rem);
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
    justify-content: space-between;
    gap: 1rem;
    align-items: flex-start;
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
    font-size: clamp(1.6rem, 3.2vw, 2.3rem);
  }

  .actions {
    display: flex;
    gap: 0.5rem;
    flex-wrap: wrap;
  }

  .observe,
  .refresh {
    padding: 0.7rem 1rem;
    border-radius: 999px;
    border: 0;
    font: inherit;
    cursor: pointer;
  }

  .observe {
    background: #142033;
    color: #fff;
  }

  .refresh {
    background: #e8eff8;
    color: #17263a;
  }

  button:disabled {
    opacity: 0.7;
    cursor: progress;
  }

  .input-panel {
    display: grid;
    gap: 0.5rem;
    font-weight: 600;
  }

  textarea {
    width: 100%;
    resize: vertical;
    min-height: 6rem;
    padding: 0.7rem 0.8rem;
    border-radius: 0.8rem;
    border: 1px solid #c7d4e2;
    font: inherit;
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

  .request-row {
    width: 100%;
    display: flex;
    justify-content: space-between;
    gap: 0.75rem;
    text-align: left;
    padding: 0.75rem 0.9rem;
    border-radius: 0.75rem;
    border: 1px solid #d5dee9;
    background: #fff;
    color: inherit;
    font: inherit;
    cursor: pointer;
  }

  .request-row.selected {
    background: #edf4ff;
    border-color: #142033;
  }

  .candidate-list li {
    display: grid;
    gap: 0.25rem;
    padding: 0.65rem 0.75rem;
    border-radius: 0.65rem;
    border: 1px solid #d5dee9;
    background: #fff;
  }

  @media (min-width: 920px) {
    .content {
      grid-template-columns: minmax(0, 2fr) minmax(18rem, 1fr);
      align-items: start;
    }
  }

  @media (max-width: 720px) {
    .header {
      flex-direction: column;
      align-items: stretch;
    }
  }
</style>
