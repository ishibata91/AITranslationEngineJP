<script lang="ts">
  import { createEventDispatcher } from "svelte";
  import type { ExecutionObserveScreenState } from "@application/usecases/execution-observe";

  export let state = undefined as unknown as ExecutionObserveScreenState;

  const dispatch = createEventDispatcher<{
    refresh: void;
  }>();

  let snapshot = null as ExecutionObserveScreenState["snapshot"];
  let shouldPrimaryFailure = false;

  $: snapshot = state.snapshot;
  $: shouldPrimaryFailure =
    snapshot !== null &&
    (snapshot.controlState === "RecoverableFailed" ||
      snapshot.controlState === "Failed" ||
      snapshot.controlState === "Canceled");
</script>

<section class="panel">
  <header class="header">
    <div>
      <p class="eyebrow">Phase 4 Execution Path</p>
      <h1>Execution Observe</h1>
    </div>
    <button type="button" on:click={() => dispatch("refresh")} disabled={state.loading}>
      {state.loading ? "Refreshing..." : "Refresh"}
    </button>
  </header>

  {#if state.error !== null}
    <section class="status error-panel">
      <p>{state.error}</p>
    </section>
  {/if}

  {#if snapshot === null}
    <section class="card">
      <h2>Job Summary</h2>
      <p>No observation snapshot available.</p>
    </section>
  {:else}
    <section class="card">
      <h2>Job Summary</h2>
      <dl class="kv">
        <div>
          <dt>Job Name</dt>
          <dd>{snapshot.summary.jobName}</dd>
        </div>
        <div>
          <dt>Status</dt>
          <dd>{snapshot.summary.statusLabel}</dd>
        </div>
        <div>
          <dt>Current Phase</dt>
          <dd>{snapshot.summary.currentPhase}</dd>
        </div>
        <div>
          <dt>Provider</dt>
          <dd>{snapshot.summary.providerLabel}</dd>
        </div>
        <div>
          <dt>Started At</dt>
          <dd>{snapshot.summary.startedAt}</dd>
        </div>
      </dl>
    </section>

    <section class="card">
      <h2>Failure Summary</h2>
      {#if snapshot.failure === null}
        <p>No failure observed.</p>
      {:else if shouldPrimaryFailure}
        <p class="failure-badge">Primary Failure</p>
        <p>{snapshot.failure.category}</p>
        <p>{snapshot.failure.message}</p>
      {:else}
        <p class="failure-badge">Supplemental Failure</p>
        <p>{snapshot.failure.category}</p>
        <p>{snapshot.failure.message}</p>
      {/if}
    </section>

    <section class="card">
      <h2>Phase Timeline</h2>
      <ul>
        {#each snapshot.phaseTimeline as phase}
          <li class:current={phase.isCurrent}>
            <span>{phase.label}</span>
            <span>{phase.statusLabel}</span>
          </li>
        {/each}
      </ul>
    </section>

    <section class="card">
      <h2>Phase Runs</h2>
      <ul>
        {#each snapshot.phaseRuns as run}
          <li>
            <span>{run.phaseKey}</span>
            <span>{run.statusLabel}</span>
            <span>{run.startedAt}</span>
            <span>{run.endedAt ?? "In progress"}</span>
          </li>
        {/each}
      </ul>
    </section>

    <section class="card">
      <h2>Translation Progress</h2>
      <dl class="kv compact">
        <div>
          <dt>Total Units</dt>
          <dd>{snapshot.translationProgress.totalUnits}</dd>
        </div>
        <div>
          <dt>Completed Units</dt>
          <dd>{snapshot.translationProgress.completedUnits}</dd>
        </div>
        <div>
          <dt>Running Units</dt>
          <dd>{snapshot.translationProgress.runningUnits}</dd>
        </div>
        <div>
          <dt>Queued Units</dt>
          <dd>{snapshot.translationProgress.queuedUnits}</dd>
        </div>
      </dl>
    </section>

    <section class="card">
      <h2>Selected Unit Detail</h2>
      {#if snapshot.selectedUnit === null}
        <p>No unit selected.</p>
      {:else}
        <dl class="kv">
          <div>
            <dt>Form ID</dt>
            <dd>{snapshot.selectedUnit.formId}</dd>
          </div>
          <div>
            <dt>Status</dt>
            <dd>{snapshot.selectedUnit.statusLabel}</dd>
          </div>
          <div>
            <dt>Source Text</dt>
            <dd>{snapshot.selectedUnit.sourceText}</dd>
          </div>
          <div>
            <dt>Dest Text</dt>
            <dd>{snapshot.selectedUnit.destText}</dd>
          </div>
        </dl>
      {/if}
    </section>

    <section class="card">
      <h2>Footer Metadata</h2>
      <dl class="kv">
        <div>
          <dt>Provider Run ID</dt>
          <dd>{snapshot.footerMetadata.providerRunId}</dd>
        </div>
        <div>
          <dt>Run Hash</dt>
          <dd>{snapshot.footerMetadata.runHash}</dd>
        </div>
        <div>
          <dt>Last Event At</dt>
          <dd>{snapshot.footerMetadata.lastEventAt}</dd>
        </div>
        <div>
          <dt>Manual Recovery Guidance</dt>
          <dd>{snapshot.footerMetadata.manualRecoveryGuidance}</dd>
        </div>
      </dl>
    </section>
  {/if}
</section>

<style>
  .panel {
    width: min(100%, 60rem);
    display: grid;
    gap: 1rem;
    padding: 1.5rem;
    border: 1px solid #d7e1eb;
    border-radius: 1.25rem;
    background:
      linear-gradient(180deg, rgba(255, 255, 255, 0.98), rgba(247, 250, 254, 0.98)),
      #fff;
    box-shadow: 0 20px 45px rgba(20, 32, 51, 0.1);
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
    font-size: clamp(1.6rem, 3.1vw, 2.3rem);
  }

  h2 {
    font-size: 1rem;
    margin-bottom: 0.75rem;
  }

  .status {
    margin: 0;
    padding: 0.85rem 1rem;
    border-radius: 0.85rem;
  }

  .error-panel {
    background: #fff1f2;
    color: #9f1239;
    border: 1px solid #fecdd3;
  }

  .card {
    padding: 1rem;
    border-radius: 0.85rem;
    border: 1px solid #dce6f1;
    background: #fff;
    display: grid;
    gap: 0.75rem;
  }

  .kv {
    display: grid;
    gap: 0.45rem;
    margin: 0;
  }

  .kv div {
    display: grid;
    gap: 0.2rem;
  }

  .kv dt {
    font-size: 0.8rem;
    color: #51637a;
  }

  .kv dd {
    margin: 0;
    word-break: break-word;
  }

  .compact {
    grid-template-columns: repeat(2, minmax(0, 1fr));
  }

  ul {
    list-style: none;
    padding: 0;
    margin: 0;
    display: grid;
    gap: 0.55rem;
  }

  li {
    padding: 0.6rem 0.75rem;
    border: 1px solid #dbe5ef;
    border-radius: 0.75rem;
    display: grid;
    gap: 0.25rem;
    background: #fbfdff;
  }

  li.current {
    border-color: #8ab4ff;
    background: #eef5ff;
  }

  button {
    padding: 0.65rem 0.9rem;
    border-radius: 0.75rem;
    border: 1px solid #c4d4e6;
    background: #fff;
    color: #17263a;
    font: inherit;
    font-weight: 600;
    cursor: pointer;
  }

  button:disabled {
    opacity: 0.6;
    cursor: not-allowed;
  }

  .failure-badge {
    display: inline-flex;
    width: fit-content;
    font-size: 0.75rem;
    letter-spacing: 0.08em;
    text-transform: uppercase;
    color: #73510a;
    background: #fff8e8;
    border: 1px solid #f2d08a;
    border-radius: 999px;
    padding: 0.2rem 0.5rem;
  }

  @media (max-width: 720px) {
    .header {
      flex-direction: column;
      align-items: stretch;
    }

    .compact {
      grid-template-columns: 1fr;
    }
  }
</style>
