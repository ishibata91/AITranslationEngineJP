<script lang="ts">
  type StatusTone = "neutral" | "info" | "success" | "warning" | "danger"

  interface Props {
    label: string
    tone?: StatusTone
    ariaLabel?: string
  }

  let { label, tone = "neutral", ariaLabel = undefined }: Props = $props()
</script>

<span class={`status-pill status-pill-${tone}`} aria-label={ariaLabel}>
  {label}
</span>

<style>
  .status-pill {
    align-items: center;
    border: 1px solid #cbd5e1;
    border-radius: 999px;
    display: inline-flex;
    font-size: 0.82rem;
    font-weight: 800;
    line-height: 1.2;
    max-width: 100%;
    min-height: 1.75rem;
    padding: 0.25rem 0.62rem;
    overflow-wrap: anywhere;
  }

  .status-pill-neutral {
    background: #f8fafc;
    color: #334155;
  }

  .status-pill-info {
    background: #eff6ff;
    border-color: #93c5fd;
    color: #1e3a8a;
  }

  .status-pill-success {
    background: #ecfdf5;
    border-color: #86efac;
    color: #14532d;
  }

  .status-pill-warning {
    background: #fff7ed;
    border-color: #fdba74;
    color: #7c2d12;
  }

  .status-pill-danger {
    background: #fef2f2;
    border-color: #fca5a5;
    color: #7f1d1d;
  }
</style>
