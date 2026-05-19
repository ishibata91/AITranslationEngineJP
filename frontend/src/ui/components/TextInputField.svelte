<script lang="ts">
  import FormField from "./FormField.svelte"

  type InputType = "text" | "search" | "password" | "email" | "url"

  interface Props {
    id: string
    label: string
    value: string
    type?: InputType
    placeholder?: string
    help?: string
    error?: string
    required?: boolean
    disabled?: boolean
    autocomplete?: string
    onInput: (value: string) => void
  }

  let {
    id,
    label,
    value,
    type = "text",
    placeholder = "",
    help = "",
    error = "",
    required = false,
    disabled = false,
    autocomplete = undefined,
    onInput
  }: Props = $props()

  const describedBy = $derived(
    [help ? `${id}-help` : "", error ? `${id}-error` : ""].filter(Boolean).join(" ") || undefined
  )
</script>

<FormField {id} {label} {help} {error} {required} {disabled}>
  <input
    aria-describedby={describedBy}
    aria-invalid={error ? "true" : undefined}
    autocomplete={autocomplete}
    class="text-input"
    disabled={disabled}
    id={id}
    oninput={(event) => onInput(event.currentTarget.value)}
    placeholder={placeholder}
    required={required}
    {type}
    {value}
  />
</FormField>

<style>
  .text-input {
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

  .text-input:focus {
    border-color: #2563eb;
    outline: 3px solid rgba(37, 99, 235, 0.18);
  }

  .text-input[aria-invalid="true"] {
    border-color: #dc2626;
  }

  .text-input:disabled {
    background: #eef1f5;
    color: #6b7280;
  }
</style>
