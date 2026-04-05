<script lang="ts">
  import { createEventDispatcher } from "svelte";
  import type {
    TranslationPreviewItem,
    TranslationPreviewScreenState,
  } from "@application/usecases/translation-preview";

  export let state = undefined as unknown as TranslationPreviewScreenState;

  const dispatch = createEventDispatcher<{
    observe: void;
    refresh: void;
    retry: void;
    select: string;
    updateJobId: string;
  }>();

  let selectedItem: TranslationPreviewItem | null = null;

  $: selectedItem =
    state.selection !== null && state.data !== null
      ? state.data.items.find((item) => item.unitKey === state.selection) ?? null
      : null;
</script>

<section class="panel">
  <header class="header">
    <div>
      <p class="eyebrow">Translation Flow Observation</p>
      <h1>Translation Preview</h1>
    </div>
    <div class="actions">
      <button
        type="button"
        class="observe"
        disabled={state.loading || state.filters.jobId.length === 0}
        on:click={() => dispatch("observe")}
      >
        {state.loading ? "Loading preview..." : "Observe"}
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
    <span>Job ID</span>
    <input
      type="text"
      value={state.filters.jobId}
      disabled={state.loading}
      on:input={(event) => {
        const target = event.currentTarget as HTMLInputElement;
        dispatch("updateJobId", target.value);
      }}
    />
  </label>

  {#if state.error}
    <div class="status error">
      <p>{state.error}</p>
      <button type="button" on:click={() => dispatch("retry")}>Retry</button>
    </div>
  {/if}

  <div class="metadata">
    <div>
      <p class="summary-label">Observed Job</p>
      <p>{state.data?.jobId ?? "Not observed"}</p>
    </div>
    <div>
      <p class="summary-label">Items</p>
      <p>{state.data?.items.length ?? 0}</p>
    </div>
  </div>

  <div class="content">
    <section class="list-panel">
      {#if state.loading && state.data === null}
        <p class="status loading">Loading translation preview...</p>
      {:else if state.data === null}
        <p class="status empty">Run preview observe to inspect translation details.</p>
      {:else if state.data.items.length === 0}
        <p class="status empty">No preview items found.</p>
      {:else}
        <ul>
          {#each state.data.items as item}
            <li>
              <button
                type="button"
                class:selected={state.selection === item.unitKey}
                class="request-row"
                on:click={() => dispatch("select", item.unitKey)}
              >
                <span>{item.unitKey}</span>
                <strong>{item.translationUnit.sourceEntityType}</strong>
              </button>
            </li>
          {/each}
        </ul>
      {/if}
    </section>

    <aside class="summary-panel">
      <p class="summary-label">Selected Preview</p>
      {#if selectedItem}
        <h2>{selectedItem.unitKey}</h2>

        <div class="detail-block">
          <p class="summary-label">Source Text</p>
          <p>{selectedItem.translationUnit.sourceText}</p>
        </div>

        <div class="detail-block">
          <p class="summary-label">Translated Text</p>
          <p>{selectedItem.translatedText}</p>
        </div>

        <div class="detail-block">
          <p class="summary-label">Reusable Terms</p>
          {#if selectedItem.reusableTerms.length === 0}
            <p>No reusable terms.</p>
          {:else}
            <ul class="compact-list">
              {#each selectedItem.reusableTerms as term}
                <li>
                  <strong>{term.sourceText}</strong>
                  <span>{term.destText}</span>
                </li>
              {/each}
            </ul>
          {/if}
        </div>

        <div class="detail-block">
          <p class="summary-label">Job Persona</p>
          {#if selectedItem.jobPersona === null}
            <p>No job persona.</p>
          {:else}
            <dl>
              <div>
                <dt>NPC FormID</dt>
                <dd>{selectedItem.jobPersona.npcFormId}</dd>
              </div>
              <div>
                <dt>Race</dt>
                <dd>{selectedItem.jobPersona.race}</dd>
              </div>
              <div>
                <dt>Sex</dt>
                <dd>{selectedItem.jobPersona.sex}</dd>
              </div>
              <div>
                <dt>Voice</dt>
                <dd>{selectedItem.jobPersona.voice}</dd>
              </div>
              <div>
                <dt>Persona Text</dt>
                <dd>{selectedItem.jobPersona.personaText}</dd>
              </div>
            </dl>
          {/if}
        </div>

        <div class="detail-block">
          <p class="summary-label">Preserved Embedded Elements</p>
          {#if selectedItem.embeddedElementPolicy.descriptors.length === 0}
            <p>No preserved elements.</p>
          {:else}
            <ul class="compact-list">
              {#each selectedItem.embeddedElementPolicy.descriptors as descriptor}
                <li>
                  <strong>{descriptor.elementId}</strong>
                  <span>{descriptor.rawText}</span>
                </li>
              {/each}
            </ul>
          {/if}
        </div>
      {:else if state.data !== null && state.data.items.length === 0}
        <p>No preview item selected.</p>
      {:else}
        <p>No preview item selected.</p>
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

  input {
    width: 100%;
    padding: 0.7rem 0.8rem;
    border-radius: 0.8rem;
    border: 1px solid #c7d4e2;
    font: inherit;
  }

  .metadata {
    display: grid;
    gap: 0.8rem;
    grid-template-columns: repeat(auto-fit, minmax(12rem, 1fr));
    padding: 1rem;
    border-radius: 0.9rem;
    background: #f8fafc;
    border: 1px solid #d7e1eb;
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

  .detail-block {
    display: grid;
    gap: 0.4rem;
    margin-top: 1rem;
  }

  .compact-list li {
    display: grid;
    gap: 0.2rem;
    padding: 0.65rem 0.75rem;
    border-radius: 0.65rem;
    border: 1px solid #d5dee9;
    background: #fff;
  }

  dl {
    margin: 0;
    display: grid;
    gap: 0.6rem;
  }

  dt {
    font-size: 0.75rem;
    color: #526277;
    text-transform: uppercase;
    letter-spacing: 0.08em;
  }

  dd {
    margin: 0.15rem 0 0;
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
