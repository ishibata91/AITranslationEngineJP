<script lang="ts">
  interface Props {
    id: string
    label: string
    checked: boolean
    help?: string
    error?: string
    disabled?: boolean
    onChange: (checked: boolean) => void
  }

  let {
    id,
    label,
    checked,
    help = "",
    error = "",
    disabled = false,
    onChange
  }: Props = $props()

  const describedBy = $derived(
    [help ? `${id}-help` : "", error ? `${id}-error` : ""].filter(Boolean).join(" ") || undefined
  )
</script>

<div class:field-disabled={disabled} class="checkbox-field">
  <label class="checkbox-label" for={id}>
    <input
      aria-describedby={describedBy}
      aria-invalid={error ? "true" : undefined}
      class="checkbox-input"
      disabled={disabled}
      id={id}
      onchange={(event) => onChange(event.currentTarget.checked)}
      type="checkbox"
      {checked}
    />
    <span>{label}</span>
  </label>
  {#if help}
    <p class="field-help" id={`${id}-help`}>{help}</p>
  {/if}
  {#if error}
    <p class="field-error" id={`${id}-error`} role="alert">{error}</p>
  {/if}
</div>

<style>
  .checkbox-field {
    display: grid;
    gap: 0.35rem;
    min-width: 0;
  }

  .checkbox-label {
    align-items: center;
    color: #243142;
    display: inline-flex;
    font-weight: 700;
    gap: 0.5rem;
    line-height: 1.35;
  }

  .checkbox-input {
    height: 1.1rem;
    width: 1.1rem;
  }

  .checkbox-input:focus-visible {
    outline: 3px solid rgba(37, 99, 235, 0.28);
    outline-offset: 2px;
  }

  .field-help,
  .field-error {
    font-size: 0.86rem;
    line-height: 1.45;
    margin: 0 0 0 1.6rem;
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
