<script lang="ts">
  import ActionButton from "./ActionButton.svelte"

  type EmptyStateTone = "neutral" | "warning" | "danger"

  interface Props {
    title: string
    message: string
    tone?: EmptyStateTone
    actionLabel?: string
    disabled?: boolean
    busy?: boolean
    onAction?: (() => void) | null
  }

  let {
    title,
    message,
    tone = "neutral",
    actionLabel = "",
    disabled = false,
    busy = false,
    onAction = null
  }: Props = $props()
</script>

<section
  class={`empty-state-panel empty-state-panel-${tone}`}
  aria-live="polite"
>
  <div class="empty-copy">
    <h3>{title}</h3>
    <p>{message}</p>
  </div>
  {#if actionLabel && onAction}
    <ActionButton
      label={actionLabel}
      variant={tone === "danger" ? "danger" : "secondary"}
      {disabled}
      {busy}
      onClick={onAction}
    />
  {/if}
</section>

<style>
  .empty-state-panel {
    align-items: center;
    border: 1px dashed #cbd5e1;
    border-radius: 0.5rem;
    display: flex;
    gap: 1rem;
    justify-content: space-between;
    min-width: 0;
    padding: 1rem;
  }

  .empty-state-panel-neutral {
    background: #f8fafc;
    color: #243142;
  }

  .empty-state-panel-warning {
    background: #fff7ed;
    border-color: #fdba74;
    color: #7c2d12;
  }

  .empty-state-panel-danger {
    background: #fef2f2;
    border-color: #fca5a5;
    color: #7f1d1d;
  }

  .empty-copy {
    display: grid;
    gap: 0.3rem;
    min-width: 0;
  }

  h3,
  p {
    margin: 0;
    overflow-wrap: anywhere;
  }

  h3 {
    font-size: 1rem;
  }

  @media (max-width: 640px) {
    .empty-state-panel {
      align-items: stretch;
      flex-direction: column;
    }
  }
</style>
