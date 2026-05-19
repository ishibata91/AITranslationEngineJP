<script lang="ts">
  import type { BodyExecutionSummaryCardProps } from "./body-phase-card-types"

  let { providerStateLabel, details }: BodyExecutionSummaryCardProps = $props()
</script>

<section
  class="phase-card"
  aria-labelledby="bodyExecutionHeading"
  data-testid="body-translation-phase-execution-summary"
>
  <div class="section-head">
    <div>
      <p class="eyebrow">execution summary</p>
      <h3 id="bodyExecutionHeading">実行 summary</h3>
    </div>
    <span class="mini-text">{providerStateLabel}</span>
  </div>
  <dl class="detail-grid compact">
    {#each details as detail (detail.label)}
      <div>
        <dt>{detail.label}</dt>
        <dd class="wrap-value">{detail.value}</dd>
      </div>
    {/each}
  </dl>
</section>

<style>
  @import "./phase-local-card.css";

  .mini-text {
    color: rgba(236, 223, 205, 0.72);
    font-size: 0.82rem;
    margin: 0 0 0.25rem;
    overflow-wrap: anywhere;
  }
</style>
