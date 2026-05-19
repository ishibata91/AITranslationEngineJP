<script lang="ts">
  import type { Snippet } from "svelte"

  interface Props {
    id: string
    label: string
    help?: string
    error?: string
    required?: boolean
    disabled?: boolean
    children?: Snippet
  }

  let {
    id,
    label,
    help = "",
    error = "",
    required = false,
    disabled = false,
    children
  }: Props = $props()

  const helpId = $derived(help ? `${id}-help` : undefined)
  const errorId = $derived(error ? `${id}-error` : undefined)
</script>

<div class:field-disabled={disabled} class="form-field">
  <label class="field-label" for={id}>
    <span>{label}</span>
    {#if required}
      <span class="required-mark" aria-label="必須">*</span>
    {/if}
  </label>
  {@render children?.()}
  {#if help}
    <p class="field-help" id={helpId}>{help}</p>
  {/if}
  {#if error}
    <p class="field-error" id={errorId} role="alert">{error}</p>
  {/if}
</div>

<style>
  .form-field {
    display: grid;
    gap: 0.38rem;
    min-width: 0;
  }

  .field-label {
    align-items: center;
    color: #243142;
    display: inline-flex;
    font-weight: 700;
    gap: 0.25rem;
    line-height: 1.3;
  }

  .required-mark {
    color: #dc2626;
  }

  .field-help,
  .field-error {
    font-size: 0.86rem;
    line-height: 1.45;
    margin: 0;
    overflow-wrap: anywhere;
  }

  .field-help {
    color: #526173;
  }

  .field-error {
    color: #b91c1c;
  }

  .field-disabled {
    opacity: 0.72;
  }
</style>
