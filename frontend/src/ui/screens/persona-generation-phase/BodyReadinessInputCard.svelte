<script lang="ts">
  import type { PhaseDetailItem } from "../../components/phase-panel-types"

  interface Props {
    details: PhaseDetailItem[]
  }

  let { details }: Props = $props()
</script>

<section
  class="phase-card"
  aria-labelledby="personaBodyReadinessHeading"
  data-testid="persona-generation-phase-body-translation-input-card"
>
  <div class="section-head">
    <div>
      <p class="eyebrow">本文翻訳の開始可否</p>
      <h3 id="personaBodyReadinessHeading">本文翻訳の入力</h3>
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
