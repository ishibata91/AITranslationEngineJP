<script lang="ts">
  import type { PhaseDetailItem } from "../../components/phase-panel-types"

  interface Props {
    details: PhaseDetailItem[]
  }

  let { details }: Props = $props()
</script>

<section
  class="phase-card"
  aria-labelledby="personaTargetSummaryHeading"
  data-testid="persona-generation-phase-target-summary-card"
>
  <div class="section-head">
    <div>
      <p class="eyebrow">target summary</p>
      <h3 id="personaTargetSummaryHeading">対象 summary</h3>
    </div>
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
</style>
