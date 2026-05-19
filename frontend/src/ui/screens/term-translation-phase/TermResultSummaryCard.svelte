<script lang="ts">
  import type { PhaseDetailItem } from "../../components/phase-panel-types"

  interface Props {
    nextPhaseStatusLabel: string
    details: PhaseDetailItem[]
  }

  let { nextPhaseStatusLabel, details }: Props = $props()
</script>

<section
  class="phase-card"
  aria-labelledby="termPhaseResultHeading"
  data-testid="term-translation-phase-result-region"
>
  <div class="section-head">
    <div>
      <p class="eyebrow">翻訳段階の結果</p>
      <h3 id="termPhaseResultHeading">結果 summary</h3>
    </div>
    <span class="mini-text">{nextPhaseStatusLabel}</span>
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
  .phase-card {
    background: rgba(33, 27, 24, 0.88);
    border: 1px solid rgba(226, 205, 173, 0.14);
    border-radius: 20px;
    box-shadow: 0 18px 40px rgba(0, 0, 0, 0.22);
    display: grid;
    gap: 1rem;
    padding: 1.4rem;
  }

  .section-head {
    align-items: flex-start;
    display: flex;
    gap: 0.8rem;
    justify-content: space-between;
  }

  .eyebrow,
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

  h3 {
    color: #fff6ea;
    margin: 0;
  }

  dt {
    color: rgba(236, 223, 205, 0.8);
    margin: 0 0 0.18rem;
  }

  dd {
    color: rgba(250, 242, 232, 0.9);
    margin: 0;
  }

  .detail-grid {
    display: grid;
    gap: 0.9rem 1rem;
    grid-template-columns: repeat(2, minmax(0, 1fr));
    margin: 0;
  }

  .detail-grid div {
    min-width: 0;
  }

  .wrap-value {
    overflow-wrap: anywhere;
    word-break: break-word;
  }

  @media (max-width: 900px) {
    .detail-grid {
      grid-template-columns: 1fr;
    }

    .section-head {
      flex-direction: column;
    }
  }
</style>
