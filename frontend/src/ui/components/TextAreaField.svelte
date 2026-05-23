<script lang="ts">
  import FormField from "./FormField.svelte"

  interface Props {
    id: string
    label: string
    value: string
    rows?: number
    placeholder?: string
    help?: string
    error?: string
    required?: boolean
    disabled?: boolean
    onInput: (value: string) => void
  }

  let {
    id,
    label,
    value,
    rows = 4,
    placeholder = "",
    help = "",
    error = "",
    required = false,
    disabled = false,
    onInput
  }: Props = $props()

  const describedBy = $derived(
    [help ? `${id}-help` : "", error ? `${id}-error` : ""]
      .filter(Boolean)
      .join(" ") || undefined
  )
</script>

<FormField {id} {label} {help} {error} {required} {disabled}>
  <textarea
    aria-describedby={describedBy}
    aria-invalid={error ? "true" : undefined}
    class="text-area"
    {disabled}
    {id}
    oninput={(event) => onInput(event.currentTarget.value)}
    {placeholder}
    {required}
    {rows}
    {value}
  ></textarea>
</FormField>

<style>
  .text-area {
    background: #ffffff;
    border: 1px solid #b7c0cc;
    border-radius: 0.5rem;
    color: #172033;
    font: inherit;
    line-height: 1.5;
    min-width: 0;
    padding: 0.58rem 0.72rem;
    resize: vertical;
    width: 100%;
  }

  .text-area:focus {
    border-color: #2563eb;
    outline: 3px solid rgba(37, 99, 235, 0.18);
  }

  .text-area[aria-invalid="true"] {
    border-color: #dc2626;
  }

  .text-area:disabled {
    background: #eef1f5;
    color: #6b7280;
  }
</style>
