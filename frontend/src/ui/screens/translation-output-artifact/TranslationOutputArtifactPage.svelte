<script lang="ts">
  import { onMount } from "svelte"

  import type {
    CreateTranslationOutputArtifactScreenController,
    TranslationOutputArtifactScreenControllerContract
  } from "@application/contract/translation-output-artifact"
  import type { TranslationOutputCompletedJobSummary } from "@application/gateway-contract/translation-output-artifact"

  import CompletedJobListPanel from "./CompletedJobListPanel.svelte"
  import DiffPreviewPanel from "./DiffPreviewPanel.svelte"
  import LatestOutputResultCard from "./LatestOutputResultCard.svelte"
  import OutputActionPanel from "./OutputActionPanel.svelte"
  import OutputSummaryHeader from "./OutputSummaryHeader.svelte"
  import SelectedJobSummaryCard from "./SelectedJobSummaryCard.svelte"

  interface Props {
    createController: CreateTranslationOutputArtifactScreenController | null
  }

  let { createController }: Props = $props()

  function resolveController(): TranslationOutputArtifactScreenControllerContract {
    if (!createController) {
      throw new Error(
        "translation output artifact screen controller factory is not provided"
      )
    }

    return createController()
  }

  const controller = resolveController()
  let viewModel = $state(controller.getViewModel())

  const unsubscribe = controller.subscribe((nextViewModel) => {
    viewModel = nextViewModel
  })

  onMount(() => {
    void (async () => {
      await controller.mount()
    })()

    return () => {
      unsubscribe()
      controller.dispose()
    }
  })

  async function selectJob(
    job: TranslationOutputCompletedJobSummary
  ): Promise<void> {
    await controller.setJobId(job.jobId)
  }

  function refresh(): void {
    void controller.refresh()
  }

  function selectCompletedJob(job: TranslationOutputCompletedJobSummary): void {
    void selectJob(job)
  }

  function setTargetGame(value: string): void {
    controller.setTargetGame(value)
  }

  function setOutputPath(value: string): void {
    controller.setOutputPath(value)
  }

  function generateArtifact(): void {
    void controller.generateArtifact()
  }

  function regenerateArtifact(): void {
    void controller.regenerateArtifact()
  }

  function selectArtifact(artifactId: number | null): void {
    void controller.setArtifactId(artifactId)
  }
</script>

<section class="output-shell" id="translationOutputArtifactView">
  <OutputSummaryHeader
    gatewayStatus={viewModel.gatewayStatus}
    statusTitle={viewModel.statusTitle}
    statusText={viewModel.statusText}
  />

  <section class="output-grid">
    <CompletedJobListPanel
      completedJobs={viewModel.completedJobs}
      selectedJobId={viewModel.selectedJobId}
      refreshDisabled={viewModel.isLoading || viewModel.isSubmitting}
      onRefresh={refresh}
      onSelectJob={selectCompletedJob}
    />

    <SelectedJobSummaryCard
      review={viewModel.review}
      viewState={viewModel.viewState}
    />

    <div class="action-stack">
      <OutputActionPanel
        targetGame={viewModel.targetGame}
        outputPath={viewModel.outputPath}
        pathState={viewModel.pathState}
        pathReason={viewModel.pathReason}
        canGenerate={viewModel.canGenerate}
        canRegenerate={viewModel.canRegenerate}
        isSubmitting={viewModel.isSubmitting}
        disabledReason={viewModel.disabledReason}
        onTargetGameChange={setTargetGame}
        onOutputPathInput={setOutputPath}
        onGenerate={generateArtifact}
        onRegenerate={regenerateArtifact}
      />
      <LatestOutputResultCard lastCommand={viewModel.lastCommand} />
    </div>
  </section>

  <DiffPreviewPanel
    compatibilitySummaryText={viewModel.compatibilitySummaryText}
    diffPreview={viewModel.diffPreview}
    onSelectArtifact={selectArtifact}
  />
</section>

<style>
  .output-shell {
    display: grid;
    gap: 1.5rem;
  }

  .output-grid,
  .action-stack {
    display: grid;
  }

  .output-grid {
    gap: 1rem;
    grid-template-columns: minmax(0, 1.1fr) minmax(0, 1fr) minmax(18rem, 0.9fr);
  }

  .action-stack {
    gap: 0.75rem;
  }

  @media (max-width: 1080px) {
    .output-grid {
      grid-template-columns: 1fr;
    }
  }
</style>
