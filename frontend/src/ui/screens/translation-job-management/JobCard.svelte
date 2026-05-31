<script lang="ts">
  import type {
    TranslationJobManagementJobCardViewModel,
    TranslationJobManagementOperationViewModel
  } from "@application/contract/translation-job-management/translation-job-management-screen-types"
  import JobOperationGroup from "./JobOperationGroup.svelte"

  interface Props {
    job: TranslationJobManagementJobCardViewModel
    onOpenJob: (
      job: TranslationJobManagementJobCardViewModel
    ) => void | Promise<void>
    onOperation: (
      job: TranslationJobManagementJobCardViewModel,
      operation: TranslationJobManagementOperationViewModel
    ) => void | Promise<void>
  }

  let { job, onOpenJob, onOperation }: Props = $props()

  function toneClass(tone: string): string {
    return `tone-${tone}`
  }

  async function handleOpenJob(): Promise<void> {
    await onOpenJob(job)
  }

  async function handleOperation(
    operation: TranslationJobManagementOperationViewModel
  ): Promise<void> {
    await onOperation(job, operation)
  }

  async function handleCardKeydown(event: KeyboardEvent): Promise<void> {
    if (event.key !== "Enter" && event.key !== " ") {
      return
    }

    event.preventDefault()
    await handleOpenJob()
  }

  const openPhaseState = $derived(
    job.canOpenPhase ? (job.jobRunTarget ? "ready" : "missing-target") : "blocked"
  )
</script>

<article
  class="job-card"
  class:is-selected={job.isSelected}
  data-job-id={job.jobId}
  data-open-phase-state={openPhaseState}
  data-testid="translation-job-management-job-card"
>
  <a
    aria-label={`ジョブ ${job.jobId} を選択して再開する`}
    class="job-card-main job-card-select"
    class:is-disabled={!job.canOpenPhase}
    data-current-phase-label={job.currentPhaseLabel}
    data-input-source-label={job.inputSourceLabel}
    data-job-id={job.jobId}
    data-progress-label={job.progressLabel}
    data-source-path={job.sourcePath}
    data-state-description={job.stateDescription}
    data-state-label={job.stateLabel}
    data-testid="translation-job-management-job-selection-region"
    href={openPhaseState === "ready"
      ? "#translation-management/job-run"
      : "#translation-management"}
    onclick={(event) => {
      event.preventDefault()
      void handleOpenJob()
    }}
    onkeydown={(event) => void handleCardKeydown(event)}
  >
    <span class="job-card-title">
      <span class="job-card-id">ジョブ #{job.jobId}</span>
      <span class="job-card-heading overflow-text">{job.title}</span>
    </span>
    <span class="job-card-inline">
      <span>{job.stateDescription}</span>
      <span>{job.currentPhaseLabel}</span>
      <span>{job.progressLabel}</span>
    </span>
    <span class="job-card-updated">{job.lastUpdatedLabel}</span>
  </a>

  <div class="job-card-side">
    <div class="job-card-status">
      <span
        class={`state-badge ${toneClass(job.stateTone)}`}
        data-testid="translation-job-management-state-label"
      >
        {job.stateLabel}
      </span>
    </div>

    <JobOperationGroup
      jobId={job.jobId}
      canOpenPhase={job.canOpenPhase}
      hasJobRunTarget={Boolean(job.jobRunTarget)}
      stopOperation={job.stopOperation}
      deleteOperation={job.deleteOperation}
      onOpenPhase={handleOpenJob}
      onOperation={handleOperation}
    />

  </div>
</article>

<style>
  .job-card {
    display: grid;
    grid-template-columns: minmax(0, 1fr) auto;
    gap: 1rem;
    align-items: start;
    padding: 1rem 1.05rem;
    border-radius: 16px;
    border: 1px solid rgba(233, 213, 186, 0.18);
    background: rgba(255, 255, 255, 0.04);
    color: var(--text);
  }

  .job-card.is-selected {
    border-color: rgba(240, 180, 100, 0.68);
    background: rgba(240, 180, 100, 0.08);
  }

  .job-card-main {
    display: grid;
    gap: 0.7rem;
    text-align: left;
  }

  .job-card-select {
    width: 100%;
    padding: 0;
    border: 0;
    background: transparent;
    color: inherit;
    cursor: pointer;
    font: inherit;
    text-decoration: none;
  }

  .job-card-select.is-disabled {
    cursor: default;
  }

  .job-card-select:focus-visible {
    outline: 3px solid rgba(255, 186, 56, 0.55);
    outline-offset: 4px;
    border-radius: 12px;
  }

  .job-card-title {
    display: grid;
    gap: 0.12rem;
  }

  .job-card-heading {
    display: block;
    font-size: 1rem;
    font-weight: 700;
  }

  .job-card-side {
    display: grid;
    gap: 0.75rem;
    justify-items: end;
    align-content: start;
  }

  .job-card-status {
    display: flex;
    gap: 0.65rem;
    align-items: center;
    flex-wrap: wrap;
    justify-content: flex-end;
  }

  .job-card-id {
    margin: 0 0 0.12rem;
    color: rgba(236, 223, 205, 0.78);
    font-size: 0.82rem;
  }

  .job-card-inline {
    display: flex;
    gap: 0.65rem 1rem;
    flex-wrap: wrap;
    font-size: 0.86rem;
    color: #fff6ea;
  }

  .job-card-updated {
    margin: 0;
    font-size: 0.82rem;
    color: rgba(236, 223, 205, 0.7);
  }

  .job-card-reasons {
    display: grid;
    gap: 0.45rem;
    width: min(100%, 22rem);
  }

  .job-card-reason {
    margin: 0;
    font-size: 0.82rem;
    color: rgba(236, 223, 205, 0.82);
  }

  .danger-text {
    color: rgba(255, 189, 173, 0.92);
  }

  .state-badge {
    display: inline-flex;
    align-items: center;
    justify-content: center;
    min-height: 2rem;
    padding: 0 0.7rem;
    border-radius: 999px;
    border: 1px solid rgba(233, 213, 186, 0.18);
    background: rgba(255, 255, 255, 0.05);
    color: #fff6ea;
    font-size: 0.82rem;
    white-space: nowrap;
  }

  .tone-info {
    border-color: rgba(113, 190, 255, 0.4);
  }

  .tone-warning {
    border-color: rgba(255, 196, 92, 0.4);
  }

  .tone-danger {
    border-color: rgba(255, 128, 102, 0.38);
  }

  .overflow-text {
    overflow-wrap: anywhere;
  }

  @media (max-width: 1120px) {
    .job-card {
      grid-template-columns: 1fr;
    }

    .job-card-side {
      justify-items: start;
    }

    .job-card-status {
      justify-content: flex-start;
    }

    .job-card-reasons {
      width: 100%;
    }
  }
</style>
