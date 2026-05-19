<script lang="ts">
  import type {
    BodyTranslationPhaseActionKind,
    BodyTranslationPhaseScreenViewModel,
    BodyTranslationPhaseViewState
  } from "@application/contract/body-translation-phase"
  import PhaseActionPanel from "../../components/PhaseActionPanel.svelte"
  import PhaseProgressPanel from "../../components/PhaseProgressPanel.svelte"
  import PhaseStatusPanel from "../../components/PhaseStatusPanel.svelte"
  import type {
    PhaseDetailItem,
    PhaseMetricCounter,
    PhaseStateToken
  } from "../../components/phase-panel-types"
  import BodyExecutionSummaryCard from "./BodyExecutionSummaryCard.svelte"
  import BodyInputSummaryCard from "./BodyInputSummaryCard.svelte"
  import BodyResultSummaryCard from "./BodyResultSummaryCard.svelte"
  import FieldResultListPanel from "./FieldResultListPanel.svelte"
  import OutputReadinessCard from "./OutputReadinessCard.svelte"

  interface Props {
    viewModel: BodyTranslationPhaseScreenViewModel
    onAction: (actionId: BodyTranslationPhaseActionKind) => void | Promise<void>
  }

  let { viewModel, onAction }: Props = $props()

  function resolveStateToken(
    viewState: BodyTranslationPhaseViewState
  ): PhaseStateToken {
    if (viewState === "not_ready" || viewState === "validation_failed") {
      return "blocked"
    }
    return viewState
  }

  const phaseMetrics = $derived<PhaseMetricCounter[]>([
    { label: "target", value: viewModel.targetCountLabel },
    { label: "processed", value: viewModel.processedCountLabel },
    { label: "success", value: viewModel.translatedCountLabel },
    { label: "failure", value: viewModel.failedCountLabel },
    { label: "skipped", value: viewModel.skippedCountLabel }
  ])

  const progressDetails = $derived<PhaseDetailItem[]>([
    { label: "開始時刻", value: viewModel.startedAtLabel },
    { label: "完了時刻", value: viewModel.finishedAtLabel },
    { label: "target count", value: viewModel.targetCountLabel },
    { label: "processed count", value: viewModel.processedCountLabel },
    { label: "success count", value: viewModel.translatedCountLabel },
    { label: "failure count", value: viewModel.failedCountLabel },
    { label: "skipped count", value: viewModel.skippedCountLabel }
  ])

  const inputDetails = $derived<PhaseDetailItem[]>([
    { label: "dictionary digest", value: viewModel.dictionaryDigestLabel },
    { label: "persona digest", value: viewModel.personaDigestLabel },
    { label: "metadata digest", value: viewModel.metadataDigestLabel },
    { label: "prompt digest", value: viewModel.promptDigestLabel },
    { label: "input snapshot", value: viewModel.inputSnapshotRefLabel },
    { label: "skipped reasons", value: viewModel.skippedReasonsLabel }
  ])

  const executionDetails = $derived<PhaseDetailItem[]>([
    { label: "provider", value: viewModel.providerLabel },
    { label: "model", value: viewModel.modelLabel },
    { label: "execution mode", value: viewModel.executionModeLabel },
    { label: "credential ref", value: viewModel.credentialRefLabel },
    { label: "request unit count", value: viewModel.requestUnitCountLabel },
    {
      label: "provider target count",
      value: viewModel.providerTargetCountLabel
    },
    {
      label: "exact dictionary excluded",
      value: viewModel.exactDictionaryExclusionCountLabel
    },
    {
      label: "partial dictionary constrained",
      value: viewModel.partialDictionaryConstraintCountLabel
    },
    { label: "output count", value: viewModel.outputCountLabel },
    { label: "late response rejected", value: viewModel.lateResponseLabel }
  ])

  const resultDetails = $derived<PhaseDetailItem[]>([
    { label: "translated count", value: viewModel.translatedCountLabel },
    { label: "failed count", value: viewModel.failedCountLabel },
    { label: "skipped count", value: viewModel.skippedCountLabel },
    {
      label: "output ready count",
      value: viewModel.outputReadinessCompletedFieldCountLabel
    },
    { label: "result output count", value: viewModel.resultOutputCountLabel },
    { label: "status consistency", value: viewModel.outputReadinessStatusLabel }
  ])

  const readinessDetails = $derived<PhaseDetailItem[]>([
    { label: "readiness", value: viewModel.outputReadinessLabel },
    {
      label: "completed field count",
      value: viewModel.outputReadinessCompletedFieldCountLabel
    },
    {
      label: "status consistency",
      value: viewModel.outputReadinessStatusLabel
    },
    { label: "blocked reason", value: viewModel.outputReadinessBlockedReason }
  ])
</script>

<section class="job-run-shell" id="bodyTranslationPhaseView">
  <PhaseStatusPanel
    title="本文翻訳"
    eyebrow="translation-management"
    gatewayStatus={viewModel.gatewayStatus}
    lead="現在の翻訳段階、進行状況、項目ごとの結果、失敗情報、出力準備を同じ画面で確認し、開始、中断、再開、回復を判断します。"
    stateLabel={viewModel.phaseStateLabel}
    state={resolveStateToken(viewModel.viewState)}
    statusTitle={viewModel.statusTitle}
    statusText={viewModel.statusText}
    errorMessage={viewModel.errorMessage}
    testId="body-translation-phase-body-translation-summary"
    metrics={phaseMetrics}
  />

  <PhaseActionPanel
    headingId="bodyPhaseActionsHeading"
    testId="body-translation-phase-actions"
    currentPhaseLabel={viewModel.currentPhaseLabel}
    actions={viewModel.actionCards}
    {onAction}
  />

  <section class="summary-grid">
    <PhaseProgressPanel
      headingId="bodyPhaseProgressHeading"
      testId="body-translation-phase-progress"
      eyebrow="翻訳段階の進行状況"
      title="進行状況"
      progressLabel={viewModel.progressLabel}
      progressPercent={viewModel.progressPercent}
      progressDetail={viewModel.progressDetail}
      details={progressDetails}
    />
    <BodyInputSummaryCard
      readinessReason={viewModel.outputReadinessBlockedReason}
      details={inputDetails}
    />
  </section>

  <section class="summary-grid">
    <BodyExecutionSummaryCard
      providerStateLabel={viewModel.providerStateLabel}
      details={executionDetails}
    />
    <BodyResultSummaryCard
      outputReadinessLabel={viewModel.outputReadinessLabel}
      details={resultDetails}
    />
  </section>

  <section class="summary-grid">
    <FieldResultListPanel
      availabilityLabel={viewModel.fieldResultAvailabilityLabel}
      items={viewModel.fieldResultItems}
    />
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

  <OutputReadinessCard
    outputReadinessLabel={viewModel.outputReadinessLabel}
    details={readinessDetails}
  />
</section>

<style>
  .job-run-shell {
    display: grid;
    gap: 1.25rem;
    min-width: 0;
  }

  .job-run-card {
    background: rgba(33, 27, 24, 0.88);
    border: 1px solid rgba(226, 205, 173, 0.14);
    border-radius: 20px;
    box-shadow: 0 18px 40px rgba(0, 0, 0, 0.22);
    display: grid;
    gap: 1rem;
    min-width: 0;
    padding: 1.4rem;
  }

  .section-head {
    align-items: flex-start;
    display: flex;
    gap: 0.8rem;
    justify-content: space-between;
    min-width: 0;
  }

  .eyebrow {
    color: rgba(236, 223, 205, 0.72);
    font-size: 0.82rem;
  }

  h3,
  p {
    margin: 0;
  }

  h3 {
    color: #fff6ea;
  }

  dd {
    color: #f1e6d6;
  }

  .summary-grid,
  .detail-grid {
    display: grid;
    gap: 0.8rem;
    min-width: 0;
  }

  dt {
    color: rgba(236, 223, 205, 0.78);
    font-size: 0.82rem;
  }

  .summary-grid {
    grid-template-columns: repeat(2, minmax(0, 1fr));
  }

  .detail-grid.compact {
    grid-template-columns: repeat(auto-fit, minmax(11rem, 1fr));
  }

  .detail-grid > div {
    min-width: 0;
  }

  dd {
    margin: 0.25rem 0 0;
    overflow-wrap: anywhere;
  }

  .wrap-value {
    overflow-wrap: anywhere;
  }

  @media (max-width: 900px) {
    .summary-grid {
      grid-template-columns: 1fr;
    }

    .section-head {
      flex-direction: column;
      align-items: stretch;
    }
  }

  @media (max-width: 480px) {
    .job-run-card {
      padding: 1rem;
      border-radius: 14px;
    }

    .detail-grid.compact {
      grid-template-columns: minmax(0, 1fr);
    }
  }
</style>
