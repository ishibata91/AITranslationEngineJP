<script lang="ts">
  import type { PhaseMetricCounter, PhaseStateToken } from "./phase-panel-types"

  interface Props {
    title: string
    eyebrow: string
    gatewayStatus: string
    lead: string
    stateLabel: string
    state: PhaseStateToken
    statusTitle: string
    statusText: string
    errorMessage?: string
    testId: string
    statusTestId?: string
    metrics?: PhaseMetricCounter[]
  }

  let {
    title,
    eyebrow,
    gatewayStatus,
    lead,
    stateLabel,
    state,
    statusTitle,
    statusText,
    errorMessage = "",
    testId,
    statusTestId = undefined,
    metrics = []
  }: Props = $props()
</script>

<section class="phase-card hero-card" data-testid={testId}>
  <div class="hero-head">
    <div>
      <p class="eyebrow">{eyebrow}</p>
      <h2>{title}</h2>
    </div>
    <p class="gateway-status">Gateway: {gatewayStatus}</p>
  </div>
  <p class="lead">{lead}</p>
  <div class="status-block" data-testid={statusTestId}>
    <span class="state-pill" data-state={state}>
      {stateLabel}
    </span>
    <div>
      <p class="status-title">{statusTitle}</p>
      <p class="status-copy">{statusText}</p>
    </div>
  </div>
  {#if metrics.length > 0}
    <div class="counter-grid">
      {#each metrics as metric (metric.label)}
        <div>
          <span>{metric.label}</span>
          <strong>{metric.value}</strong>
        </div>
      {/each}
    </div>
  {/if}
  <p class="error-text" hidden={!errorMessage}>{errorMessage}</p>
</section>

<style>
  .phase-card {
    background: rgba(33, 27, 24, 0.88);
    border: 1px solid rgba(226, 205, 173, 0.14);
    border-radius: 20px;
    box-shadow: 0 18px 40px rgba(0, 0, 0, 0.22);
    display: grid;
    gap: 1rem;
    padding: 1.4rem;
  }

  .hero-card {
    gap: 0.85rem;
  }

  .hero-head,
  .status-block {
    align-items: flex-start;
    display: flex;
    gap: 0.9rem;
  }

  .hero-head {
    justify-content: space-between;
  }

  .eyebrow {
    color: rgba(236, 223, 205, 0.72);
    font-size: 0.82rem;
    margin: 0 0 0.25rem;
  }

  h2 {
    color: #fff6ea;
    margin: 0;
  }

  .gateway-status,
  .counter-grid span {
    color: rgba(236, 223, 205, 0.8);
  }

  .lead,
  .status-copy {
    color: rgba(250, 242, 232, 0.9);
    margin: 0;
  }

  .status-title {
    color: #fff6ea;
    font-weight: 700;
    margin: 0 0 0.2rem;
  }

  .state-pill {
    align-items: center;
    background: rgba(198, 155, 82, 0.16);
    border-radius: 999px;
    color: #ffe2ae;
    display: inline-flex;
    padding: 0.4rem 0.75rem;
    white-space: nowrap;
  }

  .state-pill[data-state="completed"],
  .state-pill[data-state="empty_completed"] {
    background: rgba(92, 156, 110, 0.18);
    color: #c8f0bf;
  }

  .state-pill[data-state="recoverable_failed"],
  .state-pill[data-state="failed"],
  .state-pill[data-state="blocked"],
  .state-pill[data-state="snapshot_missing"],
  .state-pill[data-state="canceled"] {
    background: rgba(166, 72, 59, 0.2);
    color: #ffc3b9;
  }

  .counter-grid {
    display: grid;
    gap: 0.55rem;
    grid-template-columns: repeat(5, minmax(0, 1fr));
  }

  .counter-grid div {
    background: rgba(255, 255, 255, 0.05);
    border-radius: 12px;
    display: grid;
    gap: 0.2rem;
    min-width: 0;
    padding: 0.65rem 0.7rem;
  }

  .counter-grid strong {
    color: #fff6ea;
    overflow-wrap: anywhere;
  }

  .error-text {
    color: #ffbfaf;
    margin: 0;
    overflow-wrap: anywhere;
    word-break: break-word;
  }

  @media (max-width: 900px) {
    .counter-grid {
      grid-template-columns: repeat(2, minmax(0, 1fr));
    }

    .hero-head,
    .status-block {
      flex-direction: column;
    }
  }

  @media (max-width: 520px) {
    .counter-grid {
      grid-template-columns: 1fr;
    }
  }
</style>
