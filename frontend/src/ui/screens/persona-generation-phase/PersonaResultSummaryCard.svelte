<script lang="ts">
  import type { PhaseDetailItem } from "../../components/phase-panel-types"

  interface Props {
    bodyReadinessLabel: string
    details: PhaseDetailItem[]
  }

  let { bodyReadinessLabel, details }: Props = $props()
</script>

<section
  class="phase-card"
  aria-labelledby="personaResultHeading"
  data-testid="persona-generation-phase-result-summary-card"
>
  <div class="section-head">
    <div>
      <p class="eyebrow">翻訳段階の結果</p>
      <h3 id="personaResultHeading">結果 summary</h3>
    </div>
    <span class="mini-text">{bodyReadinessLabel}</span>
  </div>
  <dl class="detail-grid compact">
    {#each details as detail (detail.label)}
      <div>
        <dt>{detail.label}</dt>
        <dd class="wrap-value">
          {detail.value}
          {#if detail.note}
            <span class="detail-note">{detail.note}</span>
          {/if}
        </dd>
      </div>
    {/each}
  </dl>
</section>

<style>
  @import "./phase-local-card.css";

  .mini-text,
  .detail-note {
    color: rgba(236, 223, 205, 0.72);
    font-size: 0.82rem;
    margin: 0 0 0.25rem;
  }

  .detail-note {
    display: block;
    margin-top: 0.25rem;
  }
</style>
