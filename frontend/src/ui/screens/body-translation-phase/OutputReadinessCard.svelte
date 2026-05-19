<script lang="ts">
  import type { OutputReadinessCardProps } from "./body-phase-card-types"

  let { outputReadinessLabel, details }: OutputReadinessCardProps = $props()
</script>

<section
  class="phase-card"
  aria-labelledby="bodyReadinessHeading"
  data-testid="body-translation-phase-output-readiness"
>
  <div class="section-head">
    <div>
      <p class="eyebrow">output readiness</p>
      <h3 id="bodyReadinessHeading">後続出力 readiness</h3>
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
