<script lang="ts">
  import type { PhaseDetailItem } from "../../components/phase-panel-types"

  interface Props {
    details: PhaseDetailItem[]
  }

  let { details }: Props = $props()
</script>

<section
  class="phase-card"
  aria-labelledby="personaExecutionHeading"
  data-testid="persona-generation-phase-execution-settings-card"
>
  <div class="section-head">
    <div>
      <p class="eyebrow">AI execution summary</p>
      <h3 id="personaExecutionHeading">実行設定</h3>
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
