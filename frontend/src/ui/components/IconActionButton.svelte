<script lang="ts">
  import type { Snippet } from "svelte"

  type IconButtonVariant = "secondary" | "danger"

  interface Props {
    ariaLabel: string
    variant?: IconButtonVariant
    disabled?: boolean
    busy?: boolean
    title?: string
    onClick?: (() => void) | null
    icon?: Snippet
  }

  let {
    ariaLabel,
    variant = "secondary",
    disabled = false,
    busy = false,
    title = undefined,
    onClick = null,
    icon
  }: Props = $props()

  const isDisabled = $derived(disabled || busy)
</script>

<button
  aria-busy={busy}
  aria-label={ariaLabel}
  class={`icon-action-button icon-action-button-${variant}`}
  disabled={isDisabled}
  onclick={() => onClick?.()}
  title={title ?? ariaLabel}
  type="button"
>
  {#if busy}
    <span class="icon-spinner" aria-hidden="true"></span>
  {:else if icon}
    {@render icon()}
  {:else}
    <span aria-hidden="true">↻</span>
  {/if}
</button>

<style>
  .icon-action-button {
    align-items: center;
    background: #ffffff;
    border: 1px solid #b7c0cc;
    border-radius: 0.5rem;
    color: #1f2937;
    cursor: pointer;
    display: inline-flex;
    font: inherit;
    height: 2.5rem;
    justify-content: center;
    padding: 0;
    width: 2.5rem;
  }

  .icon-action-button-danger {
    border-color: #dc2626;
    color: #b91c1c;
  }

  .icon-action-button:focus-visible {
    outline: 3px solid rgba(37, 99, 235, 0.45);
    outline-offset: 2px;
  }

  .icon-action-button:disabled {
    cursor: not-allowed;
    opacity: 0.56;
  }

  .icon-spinner {
    border: 2px solid currentColor;
    border-right-color: transparent;
    border-radius: 999px;
    height: 1rem;
    width: 1rem;
    animation: spin 760ms linear infinite;
  }

  @keyframes spin {
    to {
      transform: rotate(1turn);
    }
  }
</style>
