<script lang="ts">
  import type { BodyInputSummaryCardProps } from "./body-phase-card-types"

  let { readinessReason, details }: BodyInputSummaryCardProps = $props()
</script>

<section
  class="phase-card"
  aria-labelledby="bodyInputSummaryHeading"
  data-testid="body-translation-phase-input-summary"
>
  <div class="section-head">
    <div>
      <p class="eyebrow">input summary</p>
      <h3 id="bodyInputSummaryHeading">入力 summary</h3>
    </div>
  </div>
  <p class="detail-note">{readinessReason}</p>
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

  .detail-note {
    color: rgba(236, 223, 205, 0.72);
    font-size: 0.82rem;
    margin: 0 0 0.25rem;
    overflow-wrap: anywhere;
  }
</style>
