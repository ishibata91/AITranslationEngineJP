<script lang="ts">
  import type {
    PersonaGenerationPhaseActionKind,
    PersonaGenerationPhaseScreenViewModel,
    PersonaGenerationPhaseViewState
  } from "@application/contract/persona-generation-phase"

  interface Props {
    viewModel: PersonaGenerationPhaseScreenViewModel
    onAction: (
      actionId: PersonaGenerationPhaseActionKind
    ) => void | Promise<void>
  }

  let { viewModel, onAction }: Props = $props()

  function resolveStateToken(
    viewState: PersonaGenerationPhaseViewState
  ): string {
    return viewState
  }
</script>

<section class="job-run-shell" id="personaGenerationPhaseView">
  <section class="job-run-card hero-card">
    <div class="hero-head">
      <div>
        <p class="eyebrow">translation-management</p>
        <h2>NPC ペルソナ生成</h2>
      </div>
      <p class="gateway-status">Gateway: {viewModel.gatewayStatus}</p>
    </div>
    <p class="lead">
      current phase、progress、target summary、phase result、body readiness
      を同じ画面で確認し、開始、中断、再開、リトライ、キャンセルを判断します。
    </p>
    <div class="status-block">
      <span
        class="state-pill"
        data-state={resolveStateToken(viewModel.viewState)}
      >
        {viewModel.phaseStateLabel}
      </span>
      <div>
        <p class="status-title">{viewModel.statusTitle}</p>
        <p class="status-copy">{viewModel.statusText}</p>
      </div>
    </div>
    <div class="counter-grid">
      <div>
        <span>target</span>
        <strong>{viewModel.targetCountLabel}</strong>
      </div>
      <div>
        <span>generated</span>
        <strong>{viewModel.generatedCountLabel}</strong>
      </div>
      <div>
        <span>failed</span>
        <strong>{viewModel.failedCountLabel}</strong>
      </div>
      <div>
        <span>skipped</span>
        <strong>{viewModel.skippedCountLabel}</strong>
      </div>
    </div>
    <p class="error-text" hidden={!viewModel.errorMessage}>
      {viewModel.errorMessage}
    </p>
  </section>

  <section
    class="job-run-card action-card"
    aria-labelledby="personaPhaseActionsHeading"
  >
    <div class="section-head">
      <div>
        <p class="eyebrow">phase control</p>
        <h3 id="personaPhaseActionsHeading">操作</h3>
      </div>
      <span class="mini-text">{viewModel.currentPhaseLabel}</span>
    </div>
    <div class="action-grid">
      {#each viewModel.actionCards as action (action.id)}
        <button
          class="action-button"
          class:primary={action.tone === "primary"}
          class:warning={action.tone === "warning"}
          disabled={action.disabled}
          onclick={() => onAction(action.id)}
          type="button"
        >
          {action.label}
        </button>
      {/each}
    </div>
    <div class="action-hints">
      {#each viewModel.actionCards as action (action.id)}
        {#if action.disabled && action.blockedReason}
          <p>{action.label}: {action.blockedReason}</p>
        {/if}
      {/each}
    </div>
  </section>

  <section class="summary-grid">
    <section class="job-run-card" aria-labelledby="personaPhaseProgressHeading">
      <div class="section-head">
        <div>
          <p class="eyebrow">phase progress</p>
          <h3 id="personaPhaseProgressHeading">進行状況</h3>
        </div>
        <span class="mini-text">{viewModel.progressLabel}</span>
      </div>
      <div
        aria-label="progress"
        aria-valuemax="100"
        aria-valuemin="0"
        aria-valuenow={viewModel.progressPercent}
        class="progress-bar"
        role="progressbar"
      >
        <span style={`width: ${viewModel.progressPercent}%`}></span>
      </div>
      <p class="progress-copy">{viewModel.progressDetail}</p>
      <dl class="detail-grid compact">
        <div>
          <dt>開始時刻</dt>
          <dd>{viewModel.startedAtLabel}</dd>
        </div>
        <div>
          <dt>完了時刻</dt>
          <dd>{viewModel.finishedAtLabel}</dd>
        </div>
        <div>
          <dt>target count</dt>
          <dd>{viewModel.targetCountLabel}</dd>
        </div>
        <div>
          <dt>generated count</dt>
          <dd>{viewModel.generatedCountLabel}</dd>
        </div>
        <div>
          <dt>failed count</dt>
          <dd>{viewModel.failedCountLabel}</dd>
        </div>
        <div>
          <dt>skipped count</dt>
          <dd>{viewModel.skippedCountLabel}</dd>
        </div>
      </dl>
    </section>

    <section class="job-run-card" aria-labelledby="personaTargetSummaryHeading">
      <div class="section-head">
        <div>
          <p class="eyebrow">target summary</p>
          <h3 id="personaTargetSummaryHeading">対象 summary</h3>
        </div>
      </div>
      <dl class="detail-grid compact">
        <div>
          <dt>NPC count</dt>
          <dd>{viewModel.npcCountLabel}</dd>
        </div>
        <div>
          <dt>common persona hit</dt>
          <dd>{viewModel.commonPersonaHitCountLabel}</dd>
        </div>
        <div>
          <dt>common persona miss</dt>
          <dd>{viewModel.commonPersonaMissCountLabel}</dd>
        </div>
        <div>
          <dt>対象外理由</dt>
          <dd class="wrap-value">{viewModel.skippedReasonsLabel}</dd>
        </div>
        <div>
          <dt>target snapshot</dt>
          <dd class="wrap-value">{viewModel.targetSnapshotLabel}</dd>
        </div>
      </dl>
    </section>
  </section>

  <section class="summary-grid">
    <section class="job-run-card" aria-labelledby="personaExecutionHeading">
      <div class="section-head">
        <div>
          <p class="eyebrow">AI execution summary</p>
          <h3 id="personaExecutionHeading">実行設定</h3>
        </div>
      </div>
      <dl class="detail-grid compact">
        <div>
          <dt>provider</dt>
          <dd class="wrap-value">{viewModel.providerLabel}</dd>
        </div>
        <div>
          <dt>model</dt>
          <dd class="wrap-value">{viewModel.modelLabel}</dd>
        </div>
        <div>
          <dt>execution mode</dt>
          <dd class="wrap-value">{viewModel.executionModeLabel}</dd>
        </div>
        <div>
          <dt>credential ref</dt>
          <dd class="wrap-value">{viewModel.credentialRefLabel}</dd>
        </div>
        <div>
          <dt>input count</dt>
          <dd>{viewModel.inputCountLabel}</dd>
        </div>
        <div>
          <dt>output count</dt>
          <dd>{viewModel.outputCountLabel}</dd>
        </div>
        <div>
          <dt>prompt digest</dt>
          <dd class="wrap-value">{viewModel.promptDigestLabel}</dd>
        </div>
        <div>
          <dt>error kind</dt>
          <dd class="wrap-value">{viewModel.errorKindLabel}</dd>
        </div>
      </dl>
    </section>

    <section class="job-run-card" aria-labelledby="personaResultHeading">
      <div class="section-head">
        <div>
          <p class="eyebrow">phase result</p>
          <h3 id="personaResultHeading">結果 summary</h3>
        </div>
        <span class="mini-text">{viewModel.bodyReadinessLabel}</span>
      </div>
      <dl class="detail-grid compact">
        <div>
          <dt>persona snapshot</dt>
          <dd class="wrap-value">{viewModel.snapshotLabel}</dd>
        </div>
        <div>
          <dt>snapshot 参照状態</dt>
          <dd class="wrap-value">{viewModel.snapshotReferenceStatusLabel}</dd>
        </div>
        <div>
          <dt>persona count</dt>
          <dd>{viewModel.personaCountLabel}</dd>
        </div>
        <div>
          <dt>missing count</dt>
          <dd>{viewModel.missingCountLabel}</dd>
        </div>
        <div>
          <dt>body readiness</dt>
          <dd class="wrap-value">
            {viewModel.bodyReadinessLabel}
            {#if viewModel.bodyReadinessBlockedReason}
              <span class="detail-note"
                >{viewModel.bodyReadinessBlockedReason}</span
              >
            {/if}
          </dd>
        </div>
      </dl>
    </section>
  </section>

  <section class="summary-grid">
    <section class="job-run-card" aria-labelledby="personaBodyReadinessHeading">
      <div class="section-head">
        <div>
          <p class="eyebrow">body readiness</p>
          <h3 id="personaBodyReadinessHeading">body phase 入力</h3>
        </div>
      </div>
      <dl class="detail-grid compact">
        <div>
          <dt>入力 summary</dt>
          <dd class="wrap-value">{viewModel.bodyReadinessInputSummaryLabel}</dd>
        </div>
        <div>
          <dt>evidence refs</dt>
          <dd class="wrap-value">{viewModel.evidenceRefsLabel}</dd>
        </div>
      </dl>
    </section>

    <section class="job-run-card" aria-labelledby="personaErrorHeading">
      <div class="section-head">
        <div>
          <p class="eyebrow">error summary</p>
          <h3 id="personaErrorHeading">失敗情報</h3>
        </div>
      </div>
      <dl class="detail-grid compact">
        <div>
          <dt>error kind</dt>
          <dd class="wrap-value">{viewModel.errorKindLabel}</dd>
        </div>
        <div>
          <dt>短い理由</dt>
          <dd class="wrap-value">{viewModel.errorReasonLabel}</dd>
        </div>
        <div>
          <dt>retryable</dt>
          <dd>{viewModel.retryableLabel}</dd>
        </div>
      </dl>
    </section>
  </section>
</section>

<style>
  .job-run-shell {
    display: grid;
    gap: 1.25rem;
  }

  .job-run-card {
    display: grid;
    gap: 1rem;
    padding: 1.4rem;
    border: 1px solid rgba(226, 205, 173, 0.14);
    border-radius: 20px;
    background: rgba(33, 27, 24, 0.88);
    box-shadow: 0 18px 40px rgba(0, 0, 0, 0.22);
  }

  .hero-card {
    gap: 0.85rem;
  }

  .hero-head,
  .section-head {
    display: flex;
    gap: 0.8rem;
    justify-content: space-between;
    align-items: flex-start;
  }

  .eyebrow,
  .mini-text,
  .detail-note {
    color: rgba(236, 223, 205, 0.72);
    font-size: 0.82rem;
  }

  h2,
  h3 {
    margin: 0;
    color: #fff6ea;
  }

  .lead,
  .status-copy,
  .progress-copy,
  .action-hints,
  dd {
    color: rgba(250, 242, 232, 0.9);
  }

  .gateway-status,
  dt,
  .counter-grid span {
    color: rgba(236, 223, 205, 0.8);
  }

  .status-block {
    display: flex;
    gap: 0.9rem;
    align-items: flex-start;
  }

  .status-title {
    margin: 0;
    color: #fff6ea;
    font-weight: 700;
  }

  .state-pill {
    display: inline-flex;
    align-items: center;
    padding: 0.4rem 0.75rem;
    border-radius: 999px;
    background: rgba(198, 155, 82, 0.16);
    color: #ffe2ae;
    white-space: nowrap;
  }

  .state-pill[data-state="completed"],
  .state-pill[data-state="empty_completed"] {
    background: rgba(92, 156, 110, 0.18);
    color: #c8f0bf;
  }

  .state-pill[data-state="recoverable_failed"],
  .state-pill[data-state="failed"],
  .state-pill[data-state="blocked"],
  .state-pill[data-state="snapshot_missing"] {
    background: rgba(166, 72, 59, 0.2);
    color: #ffc3b9;
  }

  .counter-grid {
    display: grid;
    gap: 0.75rem;
    grid-template-columns: repeat(4, minmax(0, 1fr));
  }

  .counter-grid div {
    display: grid;
    gap: 0.3rem;
    min-width: 0;
    padding: 0.85rem 0.9rem;
    border-radius: 16px;
    background: rgba(255, 255, 255, 0.05);
  }

  .counter-grid strong {
    color: #fff6ea;
  }

  .progress-bar {
    overflow: hidden;
    height: 0.8rem;
    border-radius: 999px;
    background: rgba(255, 255, 255, 0.09);
  }

  .progress-bar span {
    display: block;
    height: 100%;
    border-radius: inherit;
    background: linear-gradient(90deg, #d8a95f 0%, #ffcf8b 100%);
  }

  .action-grid {
    display: grid;
    gap: 0.75rem;
    grid-template-columns: repeat(4, minmax(0, 1fr));
  }

  .action-button {
    min-height: 2.8rem;
    padding: 0.75rem 0.95rem;
    border-radius: 14px;
    border: 1px solid rgba(233, 213, 186, 0.18);
    background: rgba(255, 255, 255, 0.04);
    color: #fff6ea;
    font: inherit;
    cursor: pointer;
  }

  .action-button.primary {
    background: linear-gradient(135deg, #cc8a39 0%, #f0b464 100%);
    color: #1b120c;
  }

  .action-button.warning {
    background: rgba(166, 72, 59, 0.22);
    color: #ffd0c7;
  }

  .action-button:disabled {
    cursor: not-allowed;
    opacity: 0.58;
  }

  .action-hints {
    display: grid;
    gap: 0.35rem;
  }

  .summary-grid {
    display: grid;
    gap: 1.25rem;
    grid-template-columns: repeat(2, minmax(0, 1fr));
  }

  .detail-grid {
    display: grid;
    gap: 0.85rem 1rem;
    grid-template-columns: repeat(2, minmax(0, 1fr));
    margin: 0;
  }

  .detail-grid div {
    display: grid;
    gap: 0.25rem;
    min-width: 0;
  }

  dt,
  dd {
    margin: 0;
  }

  .wrap-value {
    overflow-wrap: anywhere;
    word-break: break-word;
  }

  .error-text {
    margin: 0;
    color: #ffc3b9;
  }

  @media (max-width: 900px) {
    .summary-grid,
    .action-grid,
    .counter-grid,
    .detail-grid {
      grid-template-columns: 1fr;
    }

    .hero-head,
    .section-head,
    .status-block {
      flex-direction: column;
    }
  }
</style>
