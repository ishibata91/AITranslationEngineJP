<script lang="ts">
  import type { PhaseDetailItem } from "./phase-panel-types"

  interface Props {
    headingId: string
    testId: string
    eyebrow: string
    title: string
    progressLabel: string
    progressPercent: number
    progressDetail: string
    details: PhaseDetailItem[]
  }

  let {
    headingId,
    testId,
    eyebrow,
    title,
    progressLabel,
    progressPercent,
    progressDetail,
    details
  }: Props = $props()
</script>

<section class="phase-card" aria-labelledby={headingId} data-testid={testId}>
  <div class="section-head">
    <div>
      <p class="eyebrow">{eyebrow}</p>
      <h3 id={headingId}>{title}</h3>
    </div>
    <span class="mini-text">{progressLabel}</span>
  </div>
  <div
    aria-label="progress"
    aria-valuemax="100"
    aria-valuemin="0"
    aria-valuenow={progressPercent}
    class="progress-bar"
    role="progressbar"
  >
    <span style={`width: ${progressPercent}%`}></span>
  </div>
  <p class="progress-copy">{progressDetail}</p>
  <dl class="detail-grid compact">
    {#each details as detail (detail.label)}
      <div>
        <dt>{detail.label}</dt>
        <dd>{detail.value}</dd>
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
  .mini-text {
    color: rgba(236, 223, 205, 0.72);
    font-size: 0.82rem;
    margin: 0 0 0.25rem;
  }

  h3 {
    color: #fff6ea;
    margin: 0;
  }

  .progress-copy,
  dd {
    color: rgba(250, 242, 232, 0.9);
    margin: 0;
  }

  dt {
    color: rgba(236, 223, 205, 0.8);
    margin: 0 0 0.18rem;
  }

  .progress-bar {
    background: rgba(255, 255, 255, 0.09);
    border-radius: 999px;
    height: 0.8rem;
    overflow: hidden;
  }

  .progress-bar span {
    background: linear-gradient(90deg, #d8a95f 0%, #ffcf8b 100%);
    border-radius: inherit;
    display: block;
    height: 100%;
  }

  .detail-grid {
    display: grid;
    gap: 0.9rem 1rem;
    grid-template-columns: repeat(2, minmax(0, 1fr));
    margin: 0;
  }

  .detail-grid div {
    display: grid;
    gap: 0.25rem;
    min-width: 0;
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
