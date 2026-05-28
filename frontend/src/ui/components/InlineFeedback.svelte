<script lang="ts">
  import ActionButton from "./ActionButton.svelte"

  type FeedbackTone = "neutral" | "error" | "warning" | "success"

  interface Props {
    tone?: FeedbackTone
    title?: string
    message: string
    actionLabel?: string
    dataTestId?: string
    onAction?: (() => void) | null
  }

  let {
    tone = "neutral",
    title = "",
    message,
    actionLabel = "",
    dataTestId = undefined,
    onAction = null
  }: Props = $props()
</script>

<section
  class={`inline-feedback inline-feedback-${tone}`}
  data-testid={dataTestId}
  role={tone === "error" ? "alert" : "status"}
>
  <div class="feedback-copy">
    {#if title}
      <p class="feedback-title">{title}</p>
    {/if}
    <p class="feedback-message">{message}</p>
  </div>
  {#if actionLabel && onAction}
    <ActionButton label={actionLabel} onClick={onAction} variant="secondary" />
  {/if}
</section>

<style>
  .inline-feedback {
    align-items: center;
    border: 1px solid #cbd5e1;
    border-radius: 0.5rem;
    display: flex;
    gap: 0.8rem;
    justify-content: space-between;
    min-width: 0;
    padding: 0.8rem 0.9rem;
  }

  .inline-feedback-neutral {
    background: #f8fafc;
    color: #243142;
  }

  .inline-feedback-error {
    background: #fef2f2;
    border-color: #fca5a5;
    color: #7f1d1d;
  }

  .inline-feedback-warning {
    background: #fff7ed;
    border-color: #fdba74;
    color: #7c2d12;
  }

  .inline-feedback-success {
    background: #ecfdf5;
    border-color: #86efac;
    color: #14532d;
  }

  .feedback-copy {
    display: grid;
    gap: 0.25rem;
    min-width: 0;
  }

  .feedback-title,
  .feedback-message {
    margin: 0;
    overflow-wrap: anywhere;
  }

  .feedback-title {
    font-weight: 800;
  }

  @media (max-width: 640px) {
    .inline-feedback {
      align-items: stretch;
      flex-direction: column;
    }
  }
</style>
