<script lang="ts">
  import type {
    BodyTranslationFieldResultItem,
    BodyTranslationPhaseActionKind,
    BodyTranslationPhaseScreenViewModel,
    BodyTranslationPhaseViewState
  } from "@application/contract/body-translation-phase"

  interface Props {
    viewModel: BodyTranslationPhaseScreenViewModel
    onAction: (actionId: BodyTranslationPhaseActionKind) => void | Promise<void>
  }

  let { viewModel, onAction }: Props = $props()

  function resolveStateToken(viewState: BodyTranslationPhaseViewState): string {
    return viewState
  }

  function rowKey(item: BodyTranslationFieldResultItem, index: number): string {
    return `${item.fieldId}-${item.fieldLabel}-${index}`
  }
</script>

<section class="job-run-shell" id="bodyTranslationPhaseView">
  <section
    class="job-run-card hero-card"
    data-testid="body-translation-phase-body-translation-summary"
  >
    <div class="hero-head">
      <div>
        <p class="eyebrow">translation-management</p>
        <h2>本文翻訳</h2>
      </div>
      <p class="gateway-status">Gateway: {viewModel.gatewayStatus}</p>
    </div>
    <p class="lead">
      現在の翻訳段階、進行状況、項目ごとの結果、失敗情報、出力準備を同じ画面で確認し、開始、中断、再開、回復を判断します。
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
        <span>processed</span>
        <strong>{viewModel.processedCountLabel}</strong>
      </div>
      <div>
        <span>success</span>
        <strong>{viewModel.translatedCountLabel}</strong>
      </div>
      <div>
        <span>failure</span>
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
    aria-labelledby="bodyPhaseActionsHeading"
    data-testid="body-translation-phase-actions"
  >
    <div class="section-head">
      <div>
        <p class="eyebrow">翻訳段階の操作</p>
        <h3 id="bodyPhaseActionsHeading">操作</h3>
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
    <section
      class="job-run-card"
      aria-labelledby="bodyPhaseProgressHeading"
      data-testid="body-translation-phase-progress"
    >
      <div class="section-head">
        <div>
          <p class="eyebrow">翻訳段階の進行状況</p>
          <h3 id="bodyPhaseProgressHeading">進行状況</h3>
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
          <dt>processed count</dt>
          <dd>{viewModel.processedCountLabel}</dd>
        </div>
        <div>
          <dt>success count</dt>
          <dd>{viewModel.translatedCountLabel}</dd>
        </div>
        <div>
          <dt>failure count</dt>
          <dd>{viewModel.failedCountLabel}</dd>
        </div>
        <div>
          <dt>skipped count</dt>
          <dd>{viewModel.skippedCountLabel}</dd>
        </div>
      </dl>
    </section>

    <section
      class="job-run-card"
      aria-labelledby="bodyInputSummaryHeading"
      data-testid="body-translation-phase-input-summary"
    >
      <div class="section-head">
        <div>
          <p class="eyebrow">input summary</p>
          <h3 id="bodyInputSummaryHeading">入力 summary</h3>
        </div>
      </div>
      <dl class="detail-grid compact">
        <div>
          <dt>dictionary digest</dt>
          <dd class="wrap-value">{viewModel.dictionaryDigestLabel}</dd>
        </div>
        <div>
          <dt>persona digest</dt>
          <dd class="wrap-value">{viewModel.personaDigestLabel}</dd>
        </div>
        <div>
          <dt>metadata digest</dt>
          <dd class="wrap-value">{viewModel.metadataDigestLabel}</dd>
        </div>
        <div>
          <dt>prompt digest</dt>
          <dd class="wrap-value">{viewModel.promptDigestLabel}</dd>
        </div>
        <div>
          <dt>input snapshot</dt>
          <dd class="wrap-value">{viewModel.inputSnapshotRefLabel}</dd>
        </div>
        <div>
          <dt>skipped reasons</dt>
          <dd class="wrap-value">{viewModel.skippedReasonsLabel}</dd>
        </div>
      </dl>
    </section>
  </section>

  <section class="summary-grid">
    <section
      class="job-run-card"
      aria-labelledby="bodyExecutionHeading"
      data-testid="body-translation-phase-execution-summary"
    >
      <div class="section-head">
        <div>
          <p class="eyebrow">execution summary</p>
          <h3 id="bodyExecutionHeading">実行 summary</h3>
        </div>
        <span class="mini-text">{viewModel.providerStateLabel}</span>
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
          <dt>request unit count</dt>
          <dd>{viewModel.requestUnitCountLabel}</dd>
        </div>
        <div>
          <dt>provider target count</dt>
          <dd>{viewModel.providerTargetCountLabel}</dd>
        </div>
        <div>
          <dt>exact dictionary excluded</dt>
          <dd>{viewModel.exactDictionaryExclusionCountLabel}</dd>
        </div>
        <div>
          <dt>partial dictionary constrained</dt>
          <dd>{viewModel.partialDictionaryConstraintCountLabel}</dd>
        </div>
        <div>
          <dt>output count</dt>
          <dd>{viewModel.outputCountLabel}</dd>
        </div>
        <div>
          <dt>late response rejected</dt>
          <dd>{viewModel.lateResponseLabel}</dd>
        </div>
      </dl>
    </section>

    <section
      class="job-run-card"
      aria-labelledby="bodyResultSummaryHeading"
      data-testid="body-translation-phase-result-summary"
    >
      <div class="section-head">
        <div>
          <p class="eyebrow">result summary</p>
          <h3 id="bodyResultSummaryHeading">結果 summary</h3>
        </div>
        <span class="mini-text">{viewModel.outputReadinessLabel}</span>
      </div>
      <dl class="detail-grid compact">
        <div>
          <dt>translated count</dt>
          <dd>{viewModel.translatedCountLabel}</dd>
        </div>
        <div>
          <dt>failed count</dt>
          <dd>{viewModel.failedCountLabel}</dd>
        </div>
        <div>
          <dt>skipped count</dt>
          <dd>{viewModel.skippedCountLabel}</dd>
        </div>
        <div>
          <dt>output ready count</dt>
          <dd>{viewModel.outputReadinessCompletedFieldCountLabel}</dd>
        </div>
        <div>
          <dt>result output count</dt>
          <dd>{viewModel.resultOutputCountLabel}</dd>
        </div>
        <div>
          <dt>status consistency</dt>
          <dd>{viewModel.outputReadinessStatusLabel}</dd>
        </div>
      </dl>
    </section>
  </section>

  <section class="summary-grid">
    <section
      class="job-run-card"
      aria-labelledby="bodyFieldResultHeading"
      data-testid="body-translation-phase-field-result-list"
    >
      <div class="section-head">
        <div>
          <p class="eyebrow">field result list</p>
          <h3 id="bodyFieldResultHeading">field result</h3>
        </div>
        <span class="mini-text">{viewModel.fieldResultAvailabilityLabel}</span>
      </div>

      {#if viewModel.fieldResultItems.length === 0}
        <p class="detail-note">
          field result 一覧はまだ返っていません。現在は summary
          だけを表示しています。
        </p>
      {:else}
        <div class="field-result-list">
          {#each viewModel.fieldResultItems as item, index (rowKey(item, index))}
            <article class="field-result-row">
              <div class="field-result-head">
                <strong>{item.fieldLabel}</strong>
                <span>{item.fieldId}</span>
              </div>
              <dl class="field-result-grid">
                <div>
                  <dt>record type</dt>
                  <dd>{item.recordTypeLabel}</dd>
                </div>
                <div>
                  <dt>field type</dt>
                  <dd>{item.fieldTypeLabel}</dd>
                </div>
                <div>
                  <dt>FormID</dt>
                  <dd>{item.formIdLabel}</dd>
                </div>
                <div>
                  <dt>EditorID</dt>
                  <dd>{item.editorIdLabel}</dd>
                </div>
                <div>
                  <dt>source excerpt</dt>
                  <dd>{item.sourceExcerpt}</dd>
                </div>
                <div>
                  <dt>translated text</dt>
                  <dd>{item.translatedText}</dd>
                </div>
                <div>
                  <dt>output status</dt>
                  <dd>{item.outputStatus}</dd>
                </div>
                <div>
                  <dt>protection validation</dt>
                  <dd>{item.protectionValidation}</dd>
                </div>
                <div>
                  <dt>retry count</dt>
                  <dd>{item.retryCountLabel}</dd>
                </div>
              </dl>
            </article>
          {/each}
        </div>
      {/if}
    </section>

    <section
      class="job-run-card"
      aria-labelledby="bodyErrorHeading"
      data-testid="body-translation-phase-failure-information"
    >
      <div class="section-head">
        <div>
          <p class="eyebrow">recovery panel</p>
          <h3 id="bodyErrorHeading">失敗情報</h3>
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
        <div>
          <dt>output readiness blocked reason</dt>
          <dd class="wrap-value">{viewModel.outputReadinessBlockedReason}</dd>
        </div>
      </dl>
    </section>
  </section>

  <section
    class="job-run-card"
    aria-labelledby="bodyReadinessHeading"
    data-testid="body-translation-phase-output-readiness"
  >
    <div class="section-head">
      <div>
        <p class="eyebrow">output readiness</p>
        <h3 id="bodyReadinessHeading">後続出力 readiness</h3>
      </div>
      <span class="mini-text">{viewModel.outputReadinessLabel}</span>
    </div>
    <dl class="detail-grid compact">
      <div>
        <dt>readiness</dt>
        <dd>{viewModel.outputReadinessLabel}</dd>
      </div>
      <div>
        <dt>completed field count</dt>
        <dd>{viewModel.outputReadinessCompletedFieldCountLabel}</dd>
      </div>
      <div>
        <dt>status consistency</dt>
        <dd>{viewModel.outputReadinessStatusLabel}</dd>
      </div>
      <div>
        <dt>blocked reason</dt>
        <dd class="wrap-value">{viewModel.outputReadinessBlockedReason}</dd>
      </div>
    </dl>
  </section>
</section>

<style>
  .job-run-shell {
    display: grid;
    gap: 1.25rem;
    min-width: 0;
  }

  .job-run-card {
    display: grid;
    gap: 1rem;
    padding: 1.4rem;
    border: 1px solid rgba(226, 205, 173, 0.14);
    border-radius: 20px;
    background: rgba(33, 27, 24, 0.88);
    box-shadow: 0 18px 40px rgba(0, 0, 0, 0.22);
    min-width: 0;
  }

  .hero-card {
    gap: 0.85rem;
  }

  .hero-head,
  .section-head,
  .field-result-head {
    display: flex;
    gap: 0.8rem;
    justify-content: space-between;
    align-items: flex-start;
    min-width: 0;
  }

  .eyebrow,
  .mini-text,
  .detail-note,
  .gateway-status {
    color: rgba(236, 223, 205, 0.72);
    font-size: 0.82rem;
  }

  h2,
  h3,
  p {
    margin: 0;
  }

  h2,
  h3,
  .status-title {
    color: #fff6ea;
  }

  .lead,
  .status-copy,
  .progress-copy,
  dd {
    color: #f1e6d6;
  }

  .status-block {
    display: flex;
    gap: 0.9rem;
    align-items: center;
  }

  .state-pill {
    padding: 0.4rem 0.8rem;
    border-radius: 999px;
    background: rgba(255, 246, 234, 0.08);
    color: #fff6ea;
    font-size: 0.82rem;
    white-space: nowrap;
  }

  .state-pill[data-state="ready"],
  .state-pill[data-state="completed"],
  .state-pill[data-state="empty_completed"] {
    background: rgba(155, 205, 96, 0.18);
  }

  .state-pill[data-state="running"] {
    background: rgba(240, 180, 100, 0.2);
  }

  .state-pill[data-state="paused"],
  .state-pill[data-state="recoverable_failed"],
  .state-pill[data-state="validation_failed"] {
    background: rgba(206, 137, 70, 0.2);
  }

  .state-pill[data-state="failed"],
  .state-pill[data-state="canceled"] {
    background: rgba(177, 87, 80, 0.22);
  }

  .counter-grid,
  .action-grid,
  .summary-grid,
  .detail-grid,
  .field-result-grid {
    display: grid;
    gap: 0.8rem;
    min-width: 0;
  }

  .counter-grid {
    grid-template-columns: repeat(auto-fit, minmax(7rem, 1fr));
  }

  .counter-grid div,
  .field-result-row {
    padding: 0.85rem 0.95rem;
    border-radius: 16px;
    background: rgba(255, 255, 255, 0.04);
    min-width: 0;
  }

  .counter-grid span,
  dt {
    color: rgba(236, 223, 205, 0.78);
    font-size: 0.82rem;
  }

  .counter-grid strong {
    display: block;
    margin-top: 0.35rem;
    color: #fff6ea;
    font-size: 1rem;
  }

  .action-grid {
    grid-template-columns: repeat(auto-fit, minmax(10rem, 1fr));
  }

  .action-button {
    min-height: 2.9rem;
    padding: 0.75rem 0.95rem;
    border: 1px solid rgba(233, 213, 186, 0.18);
    border-radius: 14px;
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
    background: rgba(177, 87, 80, 0.22);
  }

  .action-button:disabled {
    opacity: 0.55;
    cursor: not-allowed;
  }

  .action-hints {
    display: grid;
    gap: 0.35rem;
    color: rgba(236, 223, 205, 0.72);
    font-size: 0.82rem;
  }

  .summary-grid {
    grid-template-columns: repeat(2, minmax(0, 1fr));
  }

  .progress-bar {
    position: relative;
    overflow: hidden;
    height: 0.75rem;
    border-radius: 999px;
    background: rgba(255, 255, 255, 0.08);
  }

  .progress-bar span {
    display: block;
    height: 100%;
    border-radius: inherit;
    background: linear-gradient(135deg, #cc8a39 0%, #f0b464 100%);
  }

  .detail-grid.compact {
    grid-template-columns: repeat(auto-fit, minmax(11rem, 1fr));
  }

  .field-result-list {
    display: grid;
    gap: 0.8rem;
    min-width: 0;
  }

  .field-result-grid {
    grid-template-columns: repeat(auto-fit, minmax(12rem, 1fr));
  }

  .detail-grid > div,
  .field-result-grid > div {
    min-width: 0;
  }

  .field-result-head strong,
  .field-result-head span,
  .mini-text,
  .gateway-status,
  .lead,
  .status-copy,
  .progress-copy {
    overflow-wrap: anywhere;
  }

  dd {
    margin: 0.25rem 0 0;
    overflow-wrap: anywhere;
  }

  .wrap-value {
    overflow-wrap: anywhere;
  }

  .error-text {
    color: #ffbf9f;
  }

  @media (max-width: 900px) {
    .summary-grid {
      grid-template-columns: 1fr;
    }

    .action-grid {
      grid-template-columns: 1fr;
    }

    .status-block,
    .hero-head,
    .section-head,
    .field-result-head {
      flex-direction: column;
      align-items: stretch;
    }
  }

  @media (max-width: 480px) {
    .job-run-card {
      padding: 1rem;
      border-radius: 14px;
    }

    .counter-grid,
    .detail-grid.compact,
    .field-result-grid {
      grid-template-columns: minmax(0, 1fr);
    }

    .state-pill {
      white-space: normal;
    }
  }
</style>
