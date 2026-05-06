<script lang="ts">
  import { onMount } from "svelte"

  import type {
    CreateTranslationJobManagementScreenController,
  } from "@application/contract/translation-job-management"
  import type {
    TranslationJobManagementJobRunTarget,
    TranslationJobManagementOperationViewModel,
    TranslationJobManagementFilterChipViewModel,
    TranslationJobManagementScreenViewModel
  } from "@application/contract/translation-job-management/translation-job-management-screen-types"
  import TranslationJobManagementActionButton from "./TranslationJobManagementActionButton.svelte"
  import TranslationJobManagementDeleteModal from "./TranslationJobManagementDeleteModal.svelte"

  interface Props {
    createController: CreateTranslationJobManagementScreenController | null
    onJobRunTargetChange?: (target: TranslationJobManagementJobRunTarget | null) => void
    onOpenInputReview?: () => void
    onOpenJobSetup?: () => void
    onOpenJobRun?: () => void
  }

  let {
    createController,
    onJobRunTargetChange = () => undefined,
    onOpenInputReview = () => undefined,
    onOpenJobSetup = () => undefined,
    onOpenJobRun = () => undefined
  }: Props = $props()

  function resolveController() {
    if (!createController) {
      throw new Error(
        "translation job management screen controller factory is not provided"
      )
    }

    return createController()
  }

  const controller = resolveController()
  let viewModel = $state<TranslationJobManagementScreenViewModel>(
    controller.getViewModel()
  )

  const unsubscribe = controller.subscribe((nextViewModel) => {
    viewModel = nextViewModel
    onJobRunTargetChange(nextViewModel.jobRunTarget)
  })

  onMount(() => {
    void controller.mount()

    return () => {
      unsubscribe()
      controller.dispose()
    }
  })

  const currentFilter = $derived(
    viewModel.filterChips.find((chip) => chip.selected) ?? viewModel.filterChips[0]
  )

  async function handleOperation(
    jobId: number,
    jobState: string,
    operation: TranslationJobManagementOperationViewModel
  ): Promise<void> {
    if (!operation.enabled || operation.busy) {
      return
    }

    if (viewModel.selectedJob?.jobIdLabel !== `Job #${jobId}`) {
      await controller.selectJob(jobId)
    }

    if (operation.kind === "stop") {
      await controller.requestStop()
      if (jobState === "Ready") {
        onOpenJobSetup()
      } else {
        onOpenJobRun()
      }
      return
    }

    if (operation.kind === "resume") {
      await controller.requestResume()
      if (jobState === "Ready") {
        onOpenJobSetup()
      } else {
        onOpenJobRun()
      }
      return
    }

    controller.openDeleteConfirmation()
  }

  function handleFilterChange(event: Event): void {
    controller.setFilter((event.currentTarget as HTMLSelectElement).value as TranslationJobManagementFilterChipViewModel["id"])
  }

  function toneClass(tone: string): string {
    return `tone-${tone}`
  }

  function handleOpenJobRun(event: MouseEvent): void {
    event.stopPropagation()
    onOpenJobRun()
  }

  async function handleSelectJob(jobId: number): Promise<void> {
    await controller.selectJob(jobId)
  }

  async function handleCardKeydown(
    event: KeyboardEvent,
    jobId: number
  ): Promise<void> {
    if (event.key !== "Enter" && event.key !== " ") {
      return
    }

    event.preventDefault()
    await handleSelectJob(jobId)
  }

</script>

<section class="job-management-page">
  <section class="panel job-management-hero">
    <div class="hero-copy">
      <h2>{viewModel.pageTitle}</h2>
      <p class="hero-lead">{viewModel.pageLead}</p>
    </div>
    <button
      class="ghost-button"
      onclick={onOpenInputReview}
      type="button"
    >
      新規登録
    </button>
  </section>

  {#if viewModel.feedback}
    <section class={`panel feedback-panel ${toneClass(viewModel.feedback.tone)}`}>
      <p class="feedback-title">{viewModel.feedback.title}</p>
      <p class="feedback-message">{viewModel.feedback.message}</p>
    </section>
  {/if}

  <section class="job-management-layout">
    <section class="panel job-list-panel" aria-labelledby="jobManagementListHeading">
      <div class="section-head">
        <div>
          <p class="page-label">未完了 job</p>
          <h3 id="jobManagementListHeading">一覧</h3>
        </div>
      </div>

      <div class="filter-toolbar">
        <label class="search-field" for="jobManagementSearch">
          <span>検索</span>
          <input
            id="jobManagementSearch"
            oninput={(event) =>
              controller.setSearchQuery(
                (event.currentTarget as HTMLInputElement).value
              )}
            placeholder="job id、入力名、phase で検索"
            type="search"
            value={viewModel.searchQuery}
          />
        </label>

        <label class="select-field" for="jobManagementFilter">
          <span>状態フィルタ</span>
          <select
            id="jobManagementFilter"
            onchange={handleFilterChange}
            value={currentFilter?.id ?? "all"}
          >
            {#each viewModel.filterChips as chip (chip.id)}
              <option value={chip.id}>{chip.label} ({chip.count})</option>
            {/each}
          </select>
        </label>
      </div>

      {#if viewModel.phase === "loading"}
        <div class="panel-empty">
          <p class="empty-title">一覧を読み込んでいます</p>
          <p class="empty-description">未完了 job の状態を取得しています。</p>
        </div>
      {:else if viewModel.phase === "error"}
        <div class="panel-empty tone-danger">
          <p class="empty-title">{viewModel.listErrorTitle}</p>
          <p class="empty-description">{viewModel.listErrorDescription}</p>
        </div>
      {:else if viewModel.jobs.length === 0}
        <div class="panel-empty">
          <p class="empty-title">{viewModel.listEmptyTitle}</p>
          <p class="empty-description">{viewModel.listEmptyDescription}</p>
        </div>
      {:else}
        <div class="job-card-list">
          {#each viewModel.jobs as job (job.jobId)}
            <div
              aria-label={`Job ${job.jobId} を選択`}
              class="job-card"
              class:is-selected={job.isSelected}
              onclick={() => void handleSelectJob(job.jobId)}
              onkeydown={(event) => void handleCardKeydown(event, job.jobId)}
              role="button"
              tabindex="0"
            >
              <div class="job-card-main">
                <div class="job-card-title">
                  <p class="job-card-id">Job #{job.jobId}</p>
                  <h4 class="overflow-text">{job.title}</h4>
                </div>
                <div class="job-card-inline">
                  <span>{job.stateDescription}</span>
                  <span>{job.currentPhaseLabel}</span>
                  <span>{job.progressLabel}</span>
                </div>
                <p class="job-card-updated">{job.lastUpdatedLabel}</p>
              </div>

              <div class="job-card-side">
                <div class="job-card-status">
                  <span class={`state-badge ${toneClass(job.stateTone)}`}>
                    {job.stateLabel}
                  </span>
                  {#if job.isSelected}
                    <button
                      class="ghost-button job-run-link"
                      onclick={handleOpenJobRun}
                      type="button"
                    >
                      Job Run
                    </button>
                  {/if}
                </div>

                <div class="job-card-actions" aria-label={`Job ${job.jobId} の操作`}>
                  {#each [job.stopOperation, job.resumeOperation] as operation (operation.kind)}
                    <TranslationJobManagementActionButton
                      compact={true}
                      {operation}
                      onAction={() => handleOperation(job.jobId, job.jobState, operation)}
                    />
                  {/each}
                  <TranslationJobManagementActionButton
                    compact={true}
                    operation={job.deleteOperation}
                    onAction={() =>
                      handleOperation(job.jobId, job.jobState, job.deleteOperation)}
                    variant="danger"
                  />
                </div>

                <div class="job-card-reasons" aria-label={`Job ${job.jobId} の無効理由`}>
                  {#if !job.stopOperation.enabled && job.stopOperation.reasonText}
                    <p class="job-card-reason overflow-text">
                      停止: {job.stopOperation.reasonText}
                    </p>
                  {/if}
                  {#if !job.resumeOperation.enabled && job.resumeOperation.reasonText}
                    <p class="job-card-reason overflow-text">
                      再開: {job.resumeOperation.reasonText}
                    </p>
                  {/if}
                  {#if !job.deleteOperation.enabled && job.deleteOperation.reasonText}
                    <p class="job-card-reason overflow-text danger-text">
                      削除: {job.deleteOperation.reasonText}
                    </p>
                  {/if}
                </div>
              </div>
            </div>
          {/each}
        </div>
      {/if}
    </section>
  </section>

  <TranslationJobManagementDeleteModal
    confirmation={viewModel.deleteConfirmation}
    onClose={() => controller.closeDeleteConfirmation()}
    onConfirm={() => void controller.deleteSelectedJob()}
  />
</section>

<style>
  .job-management-page {
    display: grid;
    gap: 1rem;
  }

  .panel {
    border: 1px solid rgba(226, 205, 173, 0.14);
    border-radius: 20px;
    background: rgba(33, 27, 24, 0.88);
    padding: 1.1rem 1.2rem;
  }

  .job-management-hero,
  .section-head {
    display: flex;
    gap: 1rem;
    justify-content: space-between;
    align-items: flex-start;
    flex-wrap: wrap;
  }

  .hero-lead,
  .page-label,
  .empty-description,
  .feedback-message {
    color: rgba(236, 223, 205, 0.78);
  }

  h2,
  h3,
  h4,
  .feedback-title,
  .empty-title {
    margin: 0;
    color: #fff6ea;
  }

  .hero-lead {
    margin: 0.35rem 0 0;
  }

  .job-management-layout {
    display: grid;
    gap: 1rem;
  }

  .filter-toolbar,
  .search-field,
  .select-field,
  .job-card-list {
    display: grid;
    gap: 0.75rem;
  }

  .job-card-list {
    gap: 1rem;
    margin-top: 1.5rem;
  }

  .filter-toolbar {
    grid-template-columns: minmax(0, 1fr) minmax(220px, 280px);
    margin-top: 1rem;
  }

  .search-field span,
  .select-field span,
  .summary-label,
  .info-heading {
    color: rgba(236, 223, 205, 0.78);
    font-size: 0.82rem;
    letter-spacing: 0.04em;
  }

  .search-field input,
  .select-field select,
  button {
    min-height: 2.8rem;
    border-radius: 14px;
    border: 1px solid rgba(233, 213, 186, 0.18);
    font: inherit;
  }

  .search-field input,
  .select-field select {
    padding: 0.65rem 0.9rem;
    color: #fff6ea;
    background: rgba(255, 255, 255, 0.05);
  }

  button {
    padding: 0.65rem 1rem;
    cursor: pointer;
    background: linear-gradient(135deg, #cc8a39 0%, #f0b464 100%);
    color: #1b120c;
  }

  button:disabled {
    cursor: not-allowed;
    opacity: 0.55;
  }

  .ghost-button {
    background: rgba(255, 255, 255, 0.06);
    color: #fff6ea;
  }

  .danger-button {
    background: linear-gradient(135deg, #df6a4f 0%, #f4a16d 100%);
  }

  .job-card {
    display: grid;
    grid-template-columns: minmax(0, 1fr) auto;
    gap: 1rem;
    align-items: start;
    padding: 1rem 1.05rem;
    border-radius: 16px;
    border: 1px solid rgba(233, 213, 186, 0.18);
    background: rgba(255, 255, 255, 0.04);
    cursor: pointer;
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

  .job-card-title h4 {
    font-size: 1rem;
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

  .job-card-actions {
    display: flex;
    gap: 0.65rem;
    align-items: flex-start;
    justify-content: flex-end;
    flex-wrap: wrap;
    max-width: 22rem;
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

  .job-run-link {
    min-height: 2.2rem;
    padding: 0.5rem 0.8rem;
  }

  .danger-text {
    color: rgba(255, 189, 173, 0.92);
  }

  .panel-empty {
    display: grid;
    gap: 0.35rem;
    padding: 0.8rem 0;
  }

  .feedback-panel,
  .tone-danger,
  .tone-warning {
    border-color: rgba(240, 180, 100, 0.3);
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
    .job-management-layout {
      grid-template-columns: 1fr;
    }

    .job-card {
      grid-template-columns: 1fr;
    }

    .job-card-side {
      justify-items: start;
    }

    .job-card-status {
      justify-content: flex-start;
    }

    .job-card-actions {
      justify-content: flex-start;
      max-width: none;
    }

    .job-card-reasons {
      width: 100%;
    }
  }

  @media (max-width: 760px) {
    .filter-toolbar {
      grid-template-columns: 1fr;
    }
  }
</style>
