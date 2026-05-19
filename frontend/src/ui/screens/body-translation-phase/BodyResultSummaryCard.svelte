<script lang="ts">
  import type { BodyResultSummaryCardProps } from "./body-phase-card-types"

  let { outputReadinessLabel, details }: BodyResultSummaryCardProps = $props()
</script>

<section
  class="phase-card"
  aria-labelledby="bodyResultSummaryHeading"
  data-testid="body-translation-phase-result-summary"
>
  <div class="section-head">
    <div>
      <p class="eyebrow">result summary</p>
      <h3 id="bodyResultSummaryHeading">結果 summary</h3>
    </div>
    <span class="mini-text">{outputReadinessLabel}</span>
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
