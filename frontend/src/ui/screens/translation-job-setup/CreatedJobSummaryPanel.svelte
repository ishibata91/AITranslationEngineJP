<script lang="ts">
  import type { TranslationJobSetupSummaryResponse } from "@application/gateway-contract/translation-job-setup"

  interface Props {
    summary: TranslationJobSetupSummaryResponse
    summaryPhaseCount: number
  }

  let { summary, summaryPhaseCount }: Props = $props()
</script>

<section
  class="job-setup-card"
  aria-labelledby="jobSetupSummaryHeading"
  data-testid="translation-job-setup-created-summary-region"
>
  <div class="section-head">
    <div>
      <p class="eyebrow">create result</p>
      <h3
        aria-label={summaryPhaseCount === 0 ? "Ready job summary" : undefined}
        id="jobSetupSummaryHeading"
      >
        作成済み設定
      </h3>
    </div>
    <span class="status-pill success">{summary.jobState}</span>
  </div>
  <dl
    class="detail-grid compact"
    data-testid="translation-job-setup-created-settings-region"
  >
    <div>
      <dt>job id</dt>
      <dd>{summary.jobId}</dd>
    </div>
    <div>
      <dt>入力データ</dt>
      <dd class="wrap-value">{summary.inputSource}</dd>
    </div>
    {#if summaryPhaseCount === 0}
      <div>
        <dt>AIサービス</dt>
        <dd class="wrap-value">{summary.executionSummary.provider}</dd>
      </div>
      <div>
        <dt>モデル</dt>
        <dd class="wrap-value">{summary.executionSummary.model}</dd>
      </div>
      <div>
        <dt>実行方法</dt>
        <dd>{summary.executionSummary.executionMode}</dd>
      </div>
    {/if}
  </dl>
</section>

<style>
  .job-setup-card,
  .detail-grid div {
    display: grid;
    gap: 0.75rem;
  }

  .job-setup-card {
    gap: 1rem;
    padding: 1.25rem;
    border: 1px solid rgba(255, 212, 165, 0.18);
    border-radius: 1.25rem;
    background: rgba(34, 26, 23, 0.82);
    box-shadow: 0 20px 40px rgba(6, 4, 3, 0.18);
    color: var(--text);
  }

  .section-head {
    display: flex;
    align-items: flex-start;
    justify-content: space-between;
    gap: 1rem;
  }

  .eyebrow,
  dt {
    color: rgba(255, 215, 176, 0.72);
    font-size: 0.9rem;
  }

  .eyebrow {
    font-size: 0.78rem;
    letter-spacing: 0.08em;
    text-transform: uppercase;
  }

  .detail-grid {
    display: grid;
    gap: 0.75rem;
    grid-template-columns: repeat(auto-fit, minmax(180px, 1fr));
  }

  .detail-grid.compact {
    grid-template-columns: repeat(auto-fit, minmax(150px, 1fr));
  }

  .status-pill {
    padding: 0.35rem 0.72rem;
    border-radius: 999px;
    background: rgba(255, 190, 126, 0.14);
    color: #ffd8ae;
    font-size: 0.82rem;
  }

  .status-pill.success {
    background: rgba(145, 208, 134, 0.16);
    color: #b8f0ad;
  }

  .wrap-value {
    overflow-wrap: anywhere;
    word-break: break-word;
  }
</style>
