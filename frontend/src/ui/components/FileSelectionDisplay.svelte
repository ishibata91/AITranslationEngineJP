<script lang="ts">
  import ActionButton from "./ActionButton.svelte"
  import StatusPill from "./StatusPill.svelte"

  type FileSelectionTone = "neutral" | "success" | "warning" | "danger"

  interface Props {
    title: string
    fileName?: string
    pathLabel?: string
    hashLabel?: string
    statusLabel?: string
    statusTone?: FileSelectionTone
    emptyMessage?: string
    actionLabel?: string
    disabled?: boolean
    busy?: boolean
    error?: string
    onAction?: (() => void) | null
  }

  let {
    title,
    fileName = "",
    pathLabel = "",
    hashLabel = "",
    statusLabel = "",
    statusTone = "neutral",
    emptyMessage = "ファイルは選択されていません",
    actionLabel = "",
    disabled = false,
    busy = false,
    error = "",
    onAction = null
  }: Props = $props()

  const hasFile = $derived(fileName || pathLabel || hashLabel)
</script>

<section class="file-selection-display" aria-label={title}>
  <div class="file-header">
    <div class="file-title-group">
      <h3>{title}</h3>
      {#if statusLabel}
        <StatusPill
          label={statusLabel}
          tone={statusTone === "danger" ? "danger" : statusTone}
        />
      {/if}
    </div>
    {#if actionLabel && onAction}
      <ActionButton
        label={actionLabel}
        variant="secondary"
        {disabled}
        {busy}
        onClick={onAction}
      />
    {/if}
  </div>

  {#if hasFile}
    <dl class="file-metadata">
      {#if fileName}
        <div>
          <dt>ファイル名</dt>
          <dd>{fileName}</dd>
        </div>
      {/if}
      {#if pathLabel}
        <div>
          <dt>場所</dt>
          <dd>{pathLabel}</dd>
        </div>
      {/if}
      {#if hashLabel}
        <div>
          <dt>識別値</dt>
          <dd>{hashLabel}</dd>
        </div>
      {/if}
    </dl>
  {:else}
    <p class="file-empty">{emptyMessage}</p>
  {/if}

  {#if error}
    <p class="file-error" role="alert">{error}</p>
  {/if}
</section>

<style>
  .file-selection-display {
    border: 1px solid #cbd5e1;
    border-radius: 0.5rem;
    display: grid;
    gap: 0.85rem;
    padding: 1rem;
  }

  .file-header,
  .file-title-group {
    align-items: center;
    display: flex;
    gap: 0.7rem;
    justify-content: space-between;
    min-width: 0;
  }

  .file-title-group {
    justify-content: flex-start;
  }

  h3,
  p,
  dl {
    margin: 0;
  }

  h3 {
    color: #172033;
    font-size: 1rem;
  }

  .file-metadata {
    display: grid;
    gap: 0.55rem;
  }

  .file-metadata div {
    display: grid;
    gap: 0.2rem;
  }

  dt {
    color: #64748b;
    font-size: 0.82rem;
    font-weight: 800;
  }

  dd {
    color: #172033;
    margin: 0;
    overflow-wrap: anywhere;
  }

  .file-empty {
    color: #64748b;
  }

  .file-error {
    color: #b91c1c;
    font-weight: 700;
    overflow-wrap: anywhere;
  }

  @media (max-width: 640px) {
    .file-header {
      align-items: stretch;
      flex-direction: column;
    }
  }
</style>
