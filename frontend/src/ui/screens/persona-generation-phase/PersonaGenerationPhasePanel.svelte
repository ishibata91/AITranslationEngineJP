<script lang="ts">
  import PhaseActionPanel from "../../components/PhaseActionPanel.svelte"
  import PhaseFailureInfoCard from "../../components/PhaseFailureInfoCard.svelte"
  import PhaseProgressPanel from "../../components/PhaseProgressPanel.svelte"
  import PhaseStatusPanel from "../../components/PhaseStatusPanel.svelte"
  import type {
    PhaseDetailItem,
    PhaseMetricCounter
  } from "../../components/phase-panel-types"
  import BodyReadinessInputCard from "./BodyReadinessInputCard.svelte"
  import PersonaExecutionSettingsCard from "./PersonaExecutionSettingsCard.svelte"
  import PersonaResultSummaryCard from "./PersonaResultSummaryCard.svelte"
  import PersonaTargetSummaryCard from "./PersonaTargetSummaryCard.svelte"
  import type {
    PersonaGenerationPhaseActionKind,
    PersonaGenerationPhaseScreenViewModel
  } from "@application/contract/persona-generation-phase"

  interface Props {
    viewModel: PersonaGenerationPhaseScreenViewModel
    onAction: (
      actionId: PersonaGenerationPhaseActionKind
    ) => void | Promise<void>
  }

  let { viewModel, onAction }: Props = $props()
  const statusMetrics = $derived<PhaseMetricCounter[]>([
    { label: "target", value: viewModel.targetCountLabel },
    { label: "generated", value: viewModel.generatedCountLabel },
    { label: "failed", value: viewModel.failedCountLabel },
    { label: "skipped", value: viewModel.skippedCountLabel }
  ])
  const progressDetails = $derived<PhaseDetailItem[]>([
    { label: "開始時刻", value: viewModel.startedAtLabel },
    { label: "完了時刻", value: viewModel.finishedAtLabel },
    { label: "target count", value: viewModel.targetCountLabel },
    { label: "generated count", value: viewModel.generatedCountLabel },
    { label: "failed count", value: viewModel.failedCountLabel },
    { label: "skipped count", value: viewModel.skippedCountLabel }
  ])
  const targetDetails = $derived<PhaseDetailItem[]>([
    { label: "NPC count", value: viewModel.npcCountLabel },
    {
      label: "common persona hit",
      value: viewModel.commonPersonaHitCountLabel
    },
    {
      label: "common persona miss",
      value: viewModel.commonPersonaMissCountLabel
    },
    { label: "対象外理由", value: viewModel.skippedReasonsLabel },
    { label: "target snapshot", value: viewModel.targetSnapshotLabel }
  ])
  const executionDetails = $derived<PhaseDetailItem[]>([
    { label: "provider", value: viewModel.providerLabel },
    { label: "model", value: viewModel.modelLabel },
    { label: "execution mode", value: viewModel.executionModeLabel },
    { label: "credential ref", value: viewModel.credentialRefLabel },
    { label: "input count", value: viewModel.inputCountLabel },
    { label: "output count", value: viewModel.outputCountLabel },
    { label: "prompt digest", value: viewModel.promptDigestLabel },
    { label: "error kind", value: viewModel.errorKindLabel }
  ])
  const resultDetails = $derived<PhaseDetailItem[]>([
    { label: "persona snapshot", value: viewModel.snapshotLabel },
    {
      label: "snapshot 参照状態",
      value: viewModel.snapshotReferenceStatusLabel
    },
    { label: "persona count", value: viewModel.personaCountLabel },
    { label: "missing count", value: viewModel.missingCountLabel },
    {
      label: "body readiness",
      value: viewModel.bodyReadinessLabel,
      note: viewModel.bodyReadinessBlockedReason
    }
  ])
  const bodyReadinessDetails = $derived<PhaseDetailItem[]>([
    {
      label: "入力 summary",
      value: viewModel.bodyReadinessInputSummaryLabel
    },
    { label: "evidence refs", value: viewModel.evidenceRefsLabel }
  ])
</script>

<section class="job-run-shell" id="personaGenerationPhaseView">
  <PhaseStatusPanel
    eyebrow="translation-management"
    title="NPC ペルソナ生成"
    gatewayStatus={viewModel.gatewayStatus}
    lead="現在の翻訳段階、進行状況、対象 summary、翻訳段階の結果、本文翻訳の開始可否を同じ画面で確認し、開始、中断、再開、リトライ、キャンセルを判断します。"
    state={viewModel.viewState}
    stateLabel={viewModel.phaseStateLabel}
    statusTitle={viewModel.statusTitle}
    statusText={viewModel.statusText}
    errorMessage={viewModel.errorMessage}
    testId="persona-generation-phase-persona-generation-phase-screen"
    statusTestId="persona-generation-phase-status-summary-card"
    metrics={statusMetrics}
  />

  <PhaseActionPanel
    headingId="personaPhaseActionsHeading"
    testId="persona-generation-phase-action-card"
    currentPhaseLabel={viewModel.currentPhaseLabel}
    actions={viewModel.actionCards}
    columns={4}
    onAction={onAction}
  />

  <section class="summary-grid">
    <PhaseProgressPanel
      headingId="personaPhaseProgressHeading"
      testId="persona-generation-phase-progress-card"
      eyebrow="翻訳段階の進行状況"
      title="進行状況"
      progressLabel={viewModel.progressLabel}
      progressPercent={viewModel.progressPercent}
      progressDetail={viewModel.progressDetail}
      details={progressDetails}
    />
    <PersonaTargetSummaryCard details={targetDetails} />
  </section>

  <section class="summary-grid">
    <PersonaExecutionSettingsCard details={executionDetails} />
    <PersonaResultSummaryCard
      bodyReadinessLabel={viewModel.bodyReadinessLabel}
      details={resultDetails}
    />
  </section>

  <section class="summary-grid">
    <BodyReadinessInputCard details={bodyReadinessDetails} />
    <PhaseFailureInfoCard
      headingId="personaErrorHeading"
      testId="persona-generation-phase-failure-information-card"
      errorKindLabel={viewModel.errorKindLabel}
      errorReasonLabel={viewModel.errorReasonLabel}
      retryableLabel={viewModel.retryableLabel}
    />
  </section>
</section>

<style>
  .job-run-shell {
    display: grid;
    gap: 1.25rem;
  }

  .summary-grid {
    display: grid;
    gap: 1.25rem;
    grid-template-columns: repeat(2, minmax(0, 1fr));
  }

  @media (max-width: 900px) {
    .summary-grid {
      grid-template-columns: 1fr;
    }
  }
</style>
