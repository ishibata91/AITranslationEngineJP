<script lang="ts">
  import type { TranslationJobManagementOperationViewModel } from "@application/contract/translation-job-management/translation-job-management-screen-types"

  interface Props {
    operation: TranslationJobManagementOperationViewModel
    variant?: "default" | "danger"
    compact?: boolean
    onAction: (
      operation: TranslationJobManagementOperationViewModel
    ) => void | Promise<void>
  }

  let {
    operation,
    variant = "default",
    compact = false,
    onAction
  }: Props = $props()

  const disabledReason = $derived(
    !operation.enabled && operation.reasonText ? operation.reasonText : ""
  )

  async function handleClick(event: MouseEvent): Promise<void> {
    event.stopPropagation()
    if (!operation.enabled || operation.busy) {
      return
    }

    await onAction(operation)
  }
</script>

<div
  class="action-button-shell"
  class:is-compact={compact}
  class:is-disabled={!operation.enabled}
  data-tooltip={disabledReason}
>
  <button
    aria-label={operation.label}
    class="action-button"
    class:is-danger={variant === "danger"}
    disabled={!operation.enabled || operation.busy}
    onclick={(event) => void handleClick(event)}
    type="button"
  >
    {operation.busy ? "処理中..." : operation.label}
  </button>
</div>

<style>
  .action-button-shell {
    position: relative;
    width: 100%;
  }

  .action-button-shell.is-compact {
    width: auto;
    flex: 0 0 auto;
  }

  .action-button {
    width: 100%;
    min-width: 6.5rem;
    min-height: 2.8rem;
    padding: 0.65rem 1rem;
    border-radius: 14px;
    border: 1px solid rgba(233, 213, 186, 0.18);
    font: inherit;
    cursor: pointer;
    background: linear-gradient(135deg, #cc8a39 0%, #f0b464 100%);
    color: #1b120c;
  }

  .action-button-shell.is-compact .action-button {
    width: auto;
  }

  .action-button.is-danger {
    background: linear-gradient(135deg, #df6a4f 0%, #f4a16d 100%);
  }

  .action-button:disabled {
    cursor: not-allowed;
    opacity: 0.55;
  }

  .action-button-shell.is-disabled[data-tooltip]:hover::after,
  .action-button-shell.is-disabled[data-tooltip]:focus-within::after {
    content: attr(data-tooltip);
    position: absolute;
    left: 50%;
    bottom: calc(100% + 8px);
    transform: translateX(-50%);
    width: max-content;
    max-width: min(280px, 90vw);
    padding: 8px 10px;
    border-radius: 10px;
    background: rgba(12, 9, 8, 0.96);
    border: 1px solid rgba(233, 213, 186, 0.18);
    color: #fff6ea;
    font-size: 0.8rem;
    line-height: 1.45;
    white-space: normal;
    box-shadow: 0 18px 36px rgba(0, 0, 0, 0.34);
    z-index: 4;
    pointer-events: none;
  }
</style>
