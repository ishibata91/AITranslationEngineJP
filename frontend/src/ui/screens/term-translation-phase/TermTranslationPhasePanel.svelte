<script lang="ts">
  import PhaseActionPanel from "../../components/PhaseActionPanel.svelte"
  import PhaseFailureInfoCard from "../../components/PhaseFailureInfoCard.svelte"
  import PhaseProgressPanel from "../../components/PhaseProgressPanel.svelte"
  import PhaseStatusPanel from "../../components/PhaseStatusPanel.svelte"
  import type { PhaseDetailItem } from "../../components/phase-panel-types"
  import TermExecutionSettingsCard from "./TermExecutionSettingsCard.svelte"
  import TermResultSummaryCard from "./TermResultSummaryCard.svelte"
  import type {
    TermTranslationPhaseActionKind,
    TermTranslationPhaseScreenViewModel
  } from "@application/contract/term-translation-phase"

  type TermPanelActionKind = TermTranslationPhaseActionKind | "next-phase"

  interface Props {
    viewModel: TermTranslationPhaseScreenViewModel
    onAction: (actionId: TermPanelActionKind) => void | Promise<void>
  }

  let { viewModel, onAction }: Props = $props()
  const phaseActionCards = $derived(
    viewModel.actionCards.filter((action) => action.id !== "next-phase")
  )
  const progressDetails = $derived<PhaseDetailItem[]>([
    { label: "開始時刻", value: viewModel.startedAtLabel },
    { label: "完了時刻", value: viewModel.finishedAtLabel },
    { label: "対象語件数", value: viewModel.totalTermCountLabel },
    { label: "共通辞書 hit 件数", value: viewModel.dictionaryHitCountLabel },
    { label: "AI 実行対象語件数", value: viewModel.aiTargetCountLabel }
  ])
  const executionDetails = $derived<PhaseDetailItem[]>([
    { label: "provider", value: viewModel.providerLabel },
    { label: "model", value: viewModel.modelLabel },
    { label: "execution mode", value: viewModel.executionModeLabel },
    { label: "credential reference", value: viewModel.credentialRefLabel },
    { label: "snapshot", value: viewModel.snapshotLabel }
  ])
  const resultDetails = $derived<PhaseDetailItem[]>([
    { label: "確定訳語件数", value: viewModel.confirmedCountLabel },
    {
      label: "ジョブ内辞書反映件数",
      value: viewModel.jobDictionaryAppliedCountLabel
    },
    { label: "置換対象件数", value: viewModel.replacementTargetCountLabel },
    { label: "未一致件数", value: viewModel.unmatchedCountLabel },
    {
      label: "次の翻訳段階",
      value: viewModel.nextPhaseStatusLabel,
      note: viewModel.nextPhaseBlockedReason
    }
  ])
</script>

<section class="job-run-shell" id="termTranslationPhaseView">
  <PhaseStatusPanel
    eyebrow="translation-management"
    title="単語翻訳"
    gatewayStatus={viewModel.gatewayStatus}
    lead="現在の翻訳段階、進行状況、翻訳段階の結果、失敗情報を同じ画面で確認し、開始、中断、再開、リトライ、次の作業へ進めるかを判断します。"
    state={viewModel.viewState}
    stateLabel={viewModel.phaseStateLabel}
    statusTitle={viewModel.statusTitle}
    statusText={viewModel.statusText}
    errorMessage={viewModel.errorMessage}
    testId="term-translation-phase-screen-status-header"
  />

  <PhaseActionPanel
    headingId="termPhaseActionsHeading"
    testId="term-translation-phase-phase-actions-region"
    currentPhaseLabel={viewModel.currentPhaseLabel}
    actions={phaseActionCards}
    columns={3}
    onAction={(actionId: TermTranslationPhaseActionKind) => onAction(actionId)}
  />

  <section class="summary-grid">
    <PhaseProgressPanel
      headingId="termPhaseProgressHeading"
      testId="term-translation-phase-progress-region"
      eyebrow="翻訳段階の進行状況"
      title="進行状況"
      progressLabel={viewModel.progressLabel}
      progressPercent={viewModel.progressPercent}
      progressDetail={viewModel.progressDetail}
      details={progressDetails}
    />
    <TermExecutionSettingsCard
      providerSkippedLabel={viewModel.providerSkippedLabel}
      details={executionDetails}
    />
  </section>

  <section class="summary-grid">
    <TermResultSummaryCard
      nextPhaseStatusLabel={viewModel.nextPhaseStatusLabel}
      details={resultDetails}
    />
    <PhaseFailureInfoCard
      headingId="termPhaseErrorHeading"
      testId="term-translation-phase-failure-information-region"
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
