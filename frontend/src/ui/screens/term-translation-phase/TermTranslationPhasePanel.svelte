<script lang="ts">
  import type {
    TermTranslationPhaseActionKind,
    TermTranslationPhaseScreenViewModel,
    TermTranslationPhaseViewState
  } from "@application/contract/term-translation-phase"

  interface Props {
    viewModel: TermTranslationPhaseScreenViewModel
    onAction: (
      actionId: TermTranslationPhaseActionKind | "next-phase"
    ) => void | Promise<void>
  }

  let { viewModel, onAction }: Props = $props()

  function resolveStateToken(viewState: TermTranslationPhaseViewState): string {
    return viewState
  }
</script>

<section class="job-run-shell" id="termTranslationPhaseView">
  <section class="job-run-card hero-card">
    <div class="hero-head">
      <div>
        <p class="eyebrow">translation-management</p>
        <h2>Job Run</h2>
      </div>
      <p class="gateway-status">Gateway: {viewModel.gatewayStatus}</p>
    </div>
    <p class="lead">
      current phase、progress、phase result、error summary
      を同じ画面で確認し、開始、中断、再開、リトライ、後続 phase
      可否を判断します。
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
    <p class="error-text" hidden={!viewModel.errorMessage}>
      {viewModel.errorMessage}
    </p>
  </section>

  <section
    class="job-run-card action-card"
    aria-labelledby="termPhaseActionsHeading"
  >
    <div class="section-head">
      <div>
        <p class="eyebrow">phase control</p>
        <h3 id="termPhaseActionsHeading">操作</h3>
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
    <section class="job-run-card" aria-labelledby="termPhaseProgressHeading">
      <div class="section-head">
        <div>
          <p class="eyebrow">phase progress</p>
          <h3 id="termPhaseProgressHeading">進行状況</h3>
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
          <dt>対象語件数</dt>
          <dd>{viewModel.totalTermCountLabel}</dd>
        </div>
        <div>
          <dt>共通辞書 hit 件数</dt>
          <dd>{viewModel.dictionaryHitCountLabel}</dd>
        </div>
        <div>
          <dt>AI 実行対象語件数</dt>
          <dd>{viewModel.aiTargetCountLabel}</dd>
        </div>
      </dl>
    </section>

    <section class="job-run-card" aria-labelledby="termPhaseExecutionHeading">
      <div class="section-head">
        <div>
          <p class="eyebrow">execution summary</p>
          <h3 id="termPhaseExecutionHeading">実行設定</h3>
        </div>
        <span class="mini-text">{viewModel.providerSkippedLabel}</span>
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
          <dt>credential reference</dt>
          <dd class="wrap-value">{viewModel.credentialRefLabel}</dd>
        </div>
        <div>
          <dt>snapshot</dt>
          <dd class="wrap-value">{viewModel.snapshotLabel}</dd>
        </div>
      </dl>
    </section>
  </section>

  <section class="summary-grid">
    <section class="job-run-card" aria-labelledby="termPhaseResultHeading">
      <div class="section-head">
        <div>
          <p class="eyebrow">phase result</p>
          <h3 id="termPhaseResultHeading">結果 summary</h3>
        </div>
        <span class="mini-text">{viewModel.nextPhaseStatusLabel}</span>
      </div>
      <dl class="detail-grid compact">
        <div>
          <dt>確定訳語件数</dt>
          <dd>{viewModel.confirmedCountLabel}</dd>
        </div>
        <div>
          <dt>ジョブ内辞書反映件数</dt>
          <dd>{viewModel.jobDictionaryAppliedCountLabel}</dd>
        </div>
        <div>
          <dt>置換対象件数</dt>
          <dd>{viewModel.replacementTargetCountLabel}</dd>
        </div>
        <div>
          <dt>未一致件数</dt>
          <dd>{viewModel.unmatchedCountLabel}</dd>
        </div>
        <div>
          <dt>後続 phase</dt>
          <dd class="wrap-value">
            {viewModel.nextPhaseStatusLabel}
            {#if viewModel.nextPhaseBlockedReason}
              <span class="detail-note">{viewModel.nextPhaseBlockedReason}</span
              >
            {/if}
          </dd>
        </div>
      </dl>
    </section>

    <section class="job-run-card" aria-labelledby="termPhaseErrorHeading">
      <div class="section-head">
        <div>
          <p class="eyebrow">error summary</p>
          <h3 id="termPhaseErrorHeading">失敗情報</h3>
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
  dt {
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
  .state-pill[data-state="blocked"] {
    background: rgba(166, 72, 59, 0.2);
    color: #ffc3b9;
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

  .summary-grid {
    display: grid;
    grid-template-columns: repeat(2, minmax(0, 1fr));
    gap: 1.25rem;
  }

  .detail-grid {
    display: grid;
    grid-template-columns: repeat(2, minmax(0, 1fr));
    gap: 0.9rem 1rem;
  }

  .detail-grid.compact dt {
    margin-bottom: 0.18rem;
  }

  .wrap-value {
    overflow-wrap: anywhere;
    word-break: break-word;
  }

  .action-grid {
    display: grid;
    grid-template-columns: repeat(3, minmax(0, 1fr));
    gap: 0.8rem;
  }

  .action-button {
    min-height: 2.9rem;
    padding: 0.7rem 0.95rem;
    border: 1px solid rgba(233, 213, 186, 0.18);
    border-radius: 14px;
    background: rgba(255, 255, 255, 0.04);
    color: #fff6ea;
    cursor: pointer;
  }

  .action-button.primary {
    background: linear-gradient(135deg, #cc8a39 0%, #f0b464 100%);
    color: #1b120c;
    border-color: transparent;
  }

  .action-button.warning {
    background: rgba(185, 105, 79, 0.18);
  }

  .action-button:disabled {
    cursor: not-allowed;
    opacity: 0.5;
  }

  .action-hints {
    display: grid;
    gap: 0.35rem;
    font-size: 0.9rem;
  }

  .error-text {
    color: #ffbfaf;
    overflow-wrap: anywhere;
    word-break: break-word;
  }

  @media (max-width: 900px) {
    .summary-grid,
    .detail-grid,
    .action-grid {
      grid-template-columns: 1fr;
    }

    .status-block,
    .hero-head,
    .section-head {
      flex-direction: column;
    }
  }
</style>
