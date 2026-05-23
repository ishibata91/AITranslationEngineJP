<script lang="ts">
  import type { Snippet } from "svelte"

  type ButtonVariant = "primary" | "secondary" | "danger"
  type ButtonType = "button" | "submit" | "reset"

  interface Props {
    label: string
    variant?: ButtonVariant
    type?: ButtonType
    disabled?: boolean
    busy?: boolean
    ariaLabel?: string
    describedBy?: string
    onClick?: (() => void) | null
    leading?: Snippet
    trailing?: Snippet
  }

  let {
    label,
    variant = "secondary",
    type = "button",
    disabled = false,
    busy = false,
    ariaLabel = undefined,
    describedBy = undefined,
    onClick = null,
    leading,
    trailing
  }: Props = $props()

  const isDisabled = $derived(disabled || busy)
</script>

<button
  aria-busy={busy}
  aria-describedby={describedBy}
  aria-label={ariaLabel}
  class={`action-button action-button-${variant}`}
  disabled={isDisabled}
  onclick={() => onClick?.()}
  {type}
>
  {#if busy}
    <span class="button-spinner" aria-hidden="true"></span>
  {:else}
    {@render leading?.()}
  {/if}
  <span class="button-label">{label}</span>
  {@render trailing?.()}
</button>

<style>
  .action-button {
    align-items: center;
    border: 1px solid transparent;
    border-radius: 0.5rem;
    cursor: pointer;
    display: inline-flex;
    font: inherit;
    font-weight: 700;
    gap: 0.45rem;
    justify-content: center;
    min-height: 2.5rem;
    min-width: 6.5rem;
    padding: 0.64rem 0.9rem;
    transition:
      background 120ms ease,
      border-color 120ms ease,
      color 120ms ease,
      opacity 120ms ease;
    white-space: normal;
  }

  .action-button-primary {
    background: #f59e0b;
    border-color: #d97706;
    color: #231407;
  }

  .action-button-secondary {
    background: #ffffff;
    border-color: #b7c0cc;
    color: #1f2937;
  }

  .action-button-danger {
    background: #dc2626;
    border-color: #b91c1c;
    color: #ffffff;
  }

  .action-button:hover:not(:disabled) {
    filter: brightness(0.97);
  }

  .action-button:focus-visible {
    outline: 3px solid rgba(37, 99, 235, 0.45);
    outline-offset: 2px;
  }

  .action-button:disabled {
    cursor: not-allowed;
    opacity: 0.56;
  }

  .button-label {
    overflow-wrap: anywhere;
  }

  .button-spinner {
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
