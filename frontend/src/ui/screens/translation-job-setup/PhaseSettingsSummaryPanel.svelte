<script lang="ts">
  import type { TranslationJobSetupSummaryPhaseViewModel } from "@application/presenter/translation-job-setup/translation-job-setup.presenter"

  interface Props {
    legacyValidationPassSlices: string[]
    summaryPhaseCards: TranslationJobSetupSummaryPhaseViewModel[]
  }

  let { legacyValidationPassSlices, summaryPhaseCards }: Props = $props()
</script>

<section
  class="job-setup-card"
  aria-labelledby="jobSetupSummaryPhaseHeading"
  data-testid="translation-job-setup-phase-settings-summary"
>
  <div class="section-head">
    <div>
      <p class="eyebrow">phase settings</p>
      <h3 id="jobSetupSummaryPhaseHeading">翻訳段階ごとの設定</h3>
    </div>
  </div>
  {#if summaryPhaseCards.length === 0}
    <div class="tag-list">
      {#each legacyValidationPassSlices as slice (slice)}
        <span class="tag success">{slice}</span>
      {/each}
    </div>
  {:else}
    <div class="summary-phase-grid">
      {#each summaryPhaseCards as summaryPhase (summaryPhase.phaseId)}
        <article class="summary-phase-card">
          <h4>{summaryPhase.phaseLabel}</h4>
          <dl class="detail-grid compact">
            <div>
              <dt>AIサービス</dt>
              <dd>{summaryPhase.providerLabel}</dd>
            </div>
            <div>
              <dt>モデル</dt>
              <dd class="wrap-value">{summaryPhase.model}</dd>
            </div>
            <div>
              <dt>APIキー状態</dt>
              <dd>{summaryPhase.credentialStatusLabel}</dd>
            </div>
            <div>
              <dt>一括処理</dt>
              <dd>{summaryPhase.batchLabel}</dd>
            </div>
          </dl>
        </article>
      {/each}
    </div>
  {/if}
</section>

<style>
  .job-setup-card,
  .summary-phase-grid,
  .summary-phase-card,
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

  .summary-phase-grid {
    grid-template-columns: repeat(auto-fit, minmax(220px, 1fr));
  }

  .summary-phase-card {
    padding: 1rem;
    border-radius: 1rem;
    background: rgba(18, 13, 11, 0.62);
    border: 1px solid rgba(255, 212, 165, 0.1);
  }

  .detail-grid {
    display: grid;
    gap: 0.75rem;
    grid-template-columns: repeat(auto-fit, minmax(180px, 1fr));
  }

  .detail-grid.compact {
    grid-template-columns: repeat(auto-fit, minmax(150px, 1fr));
  }

  .tag-list {
    display: flex;
    flex-wrap: wrap;
    gap: 0.5rem;
  }

  .tag {
    padding: 0.25rem 0.6rem;
    border-radius: 999px;
    background: rgba(255, 190, 126, 0.14);
    color: #ffd8ae;
    font-size: 0.82rem;
  }

  .tag.success {
    background: rgba(145, 208, 134, 0.16);
    color: #b8f0ad;
  }

  .wrap-value {
    overflow-wrap: anywhere;
    word-break: break-word;
  }
</style>
