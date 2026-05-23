<script lang="ts">
  import { onMount } from "svelte"

  import type { CreateTranslationJobManagementScreenController } from "@application/contract/translation-job-management"
  import type {
    TranslationJobManagementJobRunTarget,
    TranslationJobManagementJobCardViewModel,
    TranslationJobManagementOperationViewModel,
    TranslationJobManagementScreenViewModel
  } from "@application/contract/translation-job-management/translation-job-management-screen-types"
  import FeedbackPanel from "./FeedbackPanel.svelte"
  import JobListPanel from "./JobListPanel.svelte"
  import JobManagementHeader from "./JobManagementHeader.svelte"
  import TranslationJobManagementDeleteModal from "./TranslationJobManagementDeleteModal.svelte"

  interface Props {
    createController: CreateTranslationJobManagementScreenController | null
    onJobRunTargetChange?: (
      target: TranslationJobManagementJobRunTarget | null
    ) => void
    onOpenInputReview?: () => void
    onOpenJobRun?: (target: TranslationJobManagementJobRunTarget | null) => void
  }

  let {
    createController,
    onJobRunTargetChange = () => undefined,
    onOpenInputReview = () => undefined,
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

  async function handleOperation(
    job: TranslationJobManagementJobCardViewModel,
    operation: TranslationJobManagementOperationViewModel
  ): Promise<void> {
    if (!operation.enabled || operation.busy) {
      return
    }

    if (viewModel.selectedJob?.jobIdLabel !== `ジョブ #${job.jobId}`) {
      await controller.selectJob(job.jobId)
    }

    if (operation.kind === "stop") {
      await controller.requestStop()
      return
    }

    if (operation.kind === "resume") {
      await controller.requestResume()
      const nextTarget = controller.getViewModel().jobRunTarget
      if (job.jobState !== "Ready" && nextTarget) {
        onJobRunTargetChange(nextTarget)
        onOpenJobRun(nextTarget)
      }
      return
    }

    controller.openDeleteConfirmation()
  }

  async function handleOpenJob(
    job: TranslationJobManagementJobCardViewModel
  ): Promise<void> {
    if (!job.canOpenPhase || !job.jobRunTarget) {
      await controller.selectJob(job.jobId)
      return
    }

    onJobRunTargetChange(job.jobRunTarget)
    onOpenJobRun(job.jobRunTarget)
    await controller.selectJob(job.jobId)
  }

  function handleSearchQueryChange(searchQuery: string): void {
    controller.setSearchQuery(searchQuery)
  }

  function handleFilterChange(
    filterId: Parameters<typeof controller.setFilter>[0]
  ): void {
    controller.setFilter(filterId)
  }
</script>

<section
  class="job-management-page"
  data-testid="translation-management-child-screen-region"
>
  <JobManagementHeader
    pageTitle={viewModel.pageTitle}
    pageLead={viewModel.pageLead}
    {onOpenInputReview}
  />

  {#if viewModel.feedback}
    <FeedbackPanel feedback={viewModel.feedback} />
  {/if}

  <JobListPanel
    phase={viewModel.phase}
    jobs={viewModel.jobs}
    searchQuery={viewModel.searchQuery}
    filterChips={viewModel.filterChips}
    listEmptyTitle={viewModel.listEmptyTitle}
    listEmptyDescription={viewModel.listEmptyDescription}
    listErrorTitle={viewModel.listErrorTitle}
    listErrorDescription={viewModel.listErrorDescription}
    onSearchQueryChange={handleSearchQueryChange}
    onFilterChange={handleFilterChange}
    onOpenJob={handleOpenJob}
    onOperation={handleOperation}
  />

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
</style>
