<script lang="ts">
  import FormField from "./FormField.svelte"

  export type SelectFieldOption = {
    value: string
    label: string
    disabled?: boolean
  }

  interface Props {
    id: string
    label: string
    value: string
    options: SelectFieldOption[]
    help?: string
    error?: string
    required?: boolean
    disabled?: boolean
    placeholder?: string
    onChange: (value: string) => void
  }

  let {
    id,
    label,
    value,
    options,
    help = "",
    error = "",
    required = false,
    disabled = false,
    placeholder = "",
    onChange
  }: Props = $props()

  const describedBy = $derived(
    [help ? `${id}-help` : "", error ? `${id}-error` : ""].filter(Boolean).join(" ") || undefined
  )
</script>

<FormField {id} {label} {help} {error} {required} {disabled}>
  <select
    aria-describedby={describedBy}
    aria-invalid={error ? "true" : undefined}
    class="select-field"
    disabled={disabled}
    id={id}
    onchange={(event) => onChange(event.currentTarget.value)}
    required={required}
    {value}
  >
    {#if placeholder}
      <option value="">{placeholder}</option>
    {/if}
    {#each options as option (option.value)}
      <option value={option.value} disabled={option.disabled}>{option.label}</option>
    {/each}
  </select>
</FormField>

<style>
  .select-field {
    background: #ffffff;
    border: 1px solid #b7c0cc;
    border-radius: 0.5rem;
    color: #172033;
    font: inherit;
    min-height: 2.5rem;
    min-width: 0;
    padding: 0.58rem 0.72rem;
    width: 100%;
  }

  .select-field:focus {
    border-color: #2563eb;
    outline: 3px solid rgba(37, 99, 235, 0.18);
  }

  .select-field[aria-invalid="true"] {
    border-color: #dc2626;
  }

  .select-field:disabled {
    background: #eef1f5;
    color: #6b7280;
  }
</style>
