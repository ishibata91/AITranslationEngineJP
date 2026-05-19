<script lang="ts">
  import type { TranslationJobManagementOperationViewModel } from "@application/contract/translation-job-management/translation-job-management-screen-types"
  import TranslationJobManagementActionButton from "./TranslationJobManagementActionButton.svelte"

  interface Props {
    jobId: number
    canOpenPhase: boolean
    stopOperation: TranslationJobManagementOperationViewModel
    resumeOperation: TranslationJobManagementOperationViewModel
    deleteOperation: TranslationJobManagementOperationViewModel
    onOpenPhase: () => void | Promise<void>
    onOperation: (
      operation: TranslationJobManagementOperationViewModel
    ) => void | Promise<void>
  }

  let {
    jobId,
    canOpenPhase,
    stopOperation,
    resumeOperation,
    deleteOperation,
    onOpenPhase,
    onOperation
  }: Props = $props()

  async function handleOpenPhase(event: MouseEvent): Promise<void> {
    event.stopPropagation()
    await onOpenPhase()
  }
</script>

<div
  class="job-card-actions"
  aria-label={`ジョブ ${jobId} の操作`}
  data-testid="translation-job-management-job-actions"
>
  <button
    class="continue-button"
    disabled={!canOpenPhase}
    onclick={(event) => void handleOpenPhase(event)}
    type="button"
  >
    現在の翻訳段階へ進む
  </button>
  {#each [stopOperation, resumeOperation] as operation (operation.kind)}
    <TranslationJobManagementActionButton
      compact={true}
      {operation}
      onAction={() => onOperation(operation)}
    />
  {/each}
  <TranslationJobManagementActionButton
    compact={true}
    operation={deleteOperation}
    onAction={() => onOperation(deleteOperation)}
    variant="danger"
  />
</div>

<style>
  .job-card-actions {
    display: flex;
    gap: 0.65rem;
    align-items: flex-start;
    justify-content: flex-end;
    flex-wrap: wrap;
    max-width: 22rem;
  }

  .continue-button {
    min-height: 2.35rem;
    border: 1px solid rgba(240, 180, 100, 0.36);
    border-radius: 0.8rem;
    background: rgba(240, 180, 100, 0.12);
    color: #ffe0b8;
    cursor: pointer;
    font: inherit;
    padding: 0.5rem 0.75rem;
  }

  .continue-button:disabled {
    cursor: not-allowed;
    opacity: 0.55;
  }

  @media (max-width: 1120px) {
    .job-card-actions {
      justify-content: flex-start;
      max-width: none;
    }
  }
</style>
