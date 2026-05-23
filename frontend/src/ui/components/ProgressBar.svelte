<script lang="ts">
  type ProgressTone = "neutral" | "success" | "warning" | "danger"

  interface Props {
    label: string
    value: number
    max?: number
    tone?: ProgressTone
    helperText?: string
    showValue?: boolean
  }

  let {
    label,
    value,
    max = 100,
    tone = "neutral",
    helperText = "",
    showValue = true
  }: Props = $props()

  const boundedValue = $derived(Math.min(Math.max(value, 0), max))
  const percent = $derived(max > 0 ? Math.round((boundedValue / max) * 100) : 0)
</script>

<div class="progress-bar">
  <div class="progress-header">
    <span class="progress-label">{label}</span>
    {#if showValue}
      <span class="progress-value">{percent}%</span>
    {/if}
  </div>
  <div
    class="progress-track"
    aria-label={label}
    aria-valuemax={max}
    aria-valuemin="0"
    aria-valuenow={boundedValue}
    role="progressbar"
  >
    <span
      class={`progress-fill progress-fill-${tone}`}
      style={`width: ${percent}%`}
    ></span>
  </div>
  {#if helperText}
    <p>{helperText}</p>
  {/if}
</div>

<style>
  .progress-bar {
    display: grid;
    gap: 0.45rem;
    min-width: 0;
  }

  .progress-header {
    align-items: center;
    display: flex;
    gap: 0.8rem;
    justify-content: space-between;
  }

  .progress-label,
  .progress-value {
    color: #172033;
    font-weight: 800;
    overflow-wrap: anywhere;
  }

  .progress-value {
    color: #475569;
    font-size: 0.9rem;
  }

  .progress-track {
    background: #e2e8f0;
    border-radius: 999px;
    height: 0.72rem;
    overflow: hidden;
  }

  .progress-fill {
    display: block;
    height: 100%;
    min-width: 0;
    transition: width 160ms ease;
  }

  .progress-fill-neutral {
    background: #2563eb;
  }

  .progress-fill-success {
    background: #16a34a;
  }

  .progress-fill-warning {
    background: #f59e0b;
  }

  .progress-fill-danger {
    background: #dc2626;
  }

  p {
    color: #64748b;
    margin: 0;
    overflow-wrap: anywhere;
  }
</style>
