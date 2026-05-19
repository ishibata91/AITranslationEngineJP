<script lang="ts">
  import type { PhaseMetricCounter } from "./phase-panel-types"

  interface Props {
    metrics: PhaseMetricCounter[]
  }

  let { metrics }: Props = $props()
</script>

<div class="counter-grid">
  {#each metrics as metric (metric.label)}
    <div>
      <span>{metric.label}</span>
      <strong>{metric.value}</strong>
    </div>
  {/each}
</div>

<style>
  .counter-grid {
    display: grid;
    gap: 0.75rem;
    grid-template-columns: repeat(4, minmax(0, 1fr));
  }

  .counter-grid div {
    background: rgba(255, 255, 255, 0.05);
    border-radius: 16px;
    display: grid;
    gap: 0.3rem;
    min-width: 0;
    padding: 0.85rem 0.9rem;
  }

  .counter-grid span {
    color: rgba(236, 223, 205, 0.8);
  }

  .counter-grid strong {
    color: #fff6ea;
    overflow-wrap: anywhere;
  }

  @media (max-width: 900px) {
    .counter-grid {
      grid-template-columns: 1fr;
    }
  }
</style>
