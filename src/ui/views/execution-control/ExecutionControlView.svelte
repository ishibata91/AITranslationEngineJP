<script lang="ts">
  import { createEventDispatcher } from "svelte";
  import type { ExecutionControlScreenState } from "@application/usecases/execution-control";

  export let state = undefined as unknown as ExecutionControlScreenState;

  const dispatch = createEventDispatcher<{
    cancel: void;
    pause: void;
    resume: void;
    retry: void;
  }>();

  let failure = null as ExecutionControlScreenState["failure"];

  $: failure = state.failure;
</script>

<section class="panel">
  <header class="header">
    <div>
      <p class="eyebrow">Phase 4 Execution Path</p>
      <h1>Execution Control</h1>
    </div>
    <p class="state-chip">{state.controlState}</p>
  </header>

  {#if failure !== null}
    <section class="status failure-panel">
      <p class="status-title">Recoverable Failure Panel</p>
      <p class="status-line">{failure.category}</p>
      <p class="status-line">{failure.message}</p>
    </section>
  {/if}

  {#if state.error}
    <section class="status error-panel">
      <p>{state.error}</p>
    </section>
  {/if}

  <div class="action-row">
    <button type="button" disabled={!state.canPause} on:click={() => dispatch("pause")}>
      {state.pendingAction === "pause" ? "Pausing..." : "Pause"}
    </button>
    <button type="button" disabled={!state.canResume} on:click={() => dispatch("resume")}>
      {state.pendingAction === "resume" ? "Resuming..." : "Resume"}
    </button>
    <button type="button" disabled={!state.canRetry} on:click={() => dispatch("retry")}>
      {state.pendingAction === "retry" ? "Retrying..." : "Retry"}
    </button>
    <button type="button" disabled={!state.canCancel} on:click={() => dispatch("cancel")}>
      {state.pendingAction === "cancel" ? "Cancelling..." : "Cancel"}
    </button>
  </div>
</section>

<style>
  .panel {
    width: min(100%, 48rem);
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
  p {
    margin: 0;
  }

  h1 {
    font-size: clamp(1.6rem, 3.2vw, 2.3rem);
  }

  .state-chip {
    padding: 0.45rem 0.85rem;
    border-radius: 999px;
    background: #e8eff8;
    color: #17263a;
    font-weight: 600;
  }

  .status {
    margin: 0;
    padding: 0.85rem 1rem;
    border-radius: 0.85rem;
  }

  .status-title {
    margin-bottom: 0.3rem;
    font-size: 0.75rem;
    letter-spacing: 0.08em;
    text-transform: uppercase;
    color: #526277;
  }

  .status-line + .status-line {
    margin-top: 0.3rem;
  }

  .failure-panel {
    background: #fff8e8;
    color: #5f4200;
    border: 1px solid #f2d08a;
  }

  .error-panel {
    background: #fff1f2;
    color: #9f1239;
    border: 1px solid #fecdd3;
  }

  .action-row {
    display: grid;
    grid-template-columns: repeat(4, minmax(0, 1fr));
    gap: 0.6rem;
  }

  button {
    padding: 0.7rem 0.9rem;
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

  @media (max-width: 720px) {
    .header {
      flex-direction: column;
      align-items: stretch;
    }

    .action-row {
      grid-template-columns: repeat(2, minmax(0, 1fr));
    }
  }
</style>
