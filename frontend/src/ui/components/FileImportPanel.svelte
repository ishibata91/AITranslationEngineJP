<script lang="ts">
  import type { Snippet } from "svelte"

  interface FileImportStat {
    label: string
    value: string | number
    valueId?: string
  }

  interface Props {
    eyebrow: string
    title: string
    titleId: string
    testId: string
    helperText: string
    selectedLabel?: string
    selectedName: string
    primaryActionId: string
    primaryActionLabel: string
    primaryActionDisabled?: boolean
    accept?: string
    inputId?: string
    stats?: FileImportStat[]
    statsId?: string
    onPrimaryAction: () => void
    onFileSelected?: (event: Event) => void
    actions?: Snippet
    children?: Snippet
  }

  let {
    eyebrow,
    title,
    titleId,
    testId,
    helperText,
    selectedLabel = "選択ファイル",
    selectedName,
    primaryActionId,
    primaryActionLabel,
    primaryActionDisabled = false,
    accept = "",
    inputId = undefined,
    stats = [],
    statsId = undefined,
    onPrimaryAction,
    onFileSelected = undefined,
    actions,
    children
  }: Props = $props()
</script>

<section class="phase-card file-import-panel" aria-labelledby={titleId} data-testid={testId}>
  <div class="panel-head">
    <div>
      <p class="eyebrow">{eyebrow}</p>
      <h3 id={titleId}>{title}</h3>
    </div>
    <button
      class="button-secondary"
      disabled={primaryActionDisabled}
      id={primaryActionId}
      onclick={onPrimaryAction}
      type="button"
    >
      {primaryActionLabel}
    </button>
  </div>

  <p class="support-copy">{helperText}</p>

  {#if inputId && onFileSelected}
    <input
      {accept}
      class="file-input"
      id={inputId}
      onchange={onFileSelected}
      type="file"
    />
  {/if}

  <div class="file-picker">
    <span class="eyebrow">{selectedLabel}</span>
    <span class="file-name">{selectedName}</span>
  </div>

  {#if stats.length > 0}
    <dl class="stats-grid" id={statsId}>
      {#each stats as stat (stat.label)}
        <div>
          <dt>{stat.label}</dt>
          <dd id={stat.valueId}>{stat.value}</dd>
        </div>
      {/each}
    </dl>
  {/if}

  {#if actions}
    <div class="button-row">
      {@render actions()}
    </div>
  {/if}

  {#if children}
    <div class="panel-body">
      {@render children()}
    </div>
  {/if}
</section>

<style>
  .phase-card {
    align-content: start;
    background: rgba(33, 27, 24, 0.88);
    border: 1px solid rgba(226, 205, 173, 0.14);
    border-radius: 20px;
    box-shadow: 0 18px 40px rgba(0, 0, 0, 0.22);
    color: var(--text);
    display: grid;
    gap: 1rem;
    min-width: 0;
    padding: 1.4rem;
  }

  .panel-head,
  .button-row {
    align-items: center;
    display: flex;
    flex-wrap: wrap;
    gap: 0.75rem;
    justify-content: space-between;
  }

  .eyebrow,
  dt {
    color: rgba(236, 223, 205, 0.72);
    font-size: 12px;
    letter-spacing: 0.08em;
    margin: 0;
  }

  h3,
  p,
  dl,
  dd {
    margin: 0;
  }

  h3,
  .file-name,
  dd {
    overflow-wrap: anywhere;
  }

  .support-copy,
  .panel-body {
    color: rgba(236, 223, 205, 0.8);
    line-height: 1.7;
  }

  .file-picker {
    align-items: center;
    display: inline-flex;
    flex-wrap: wrap;
    gap: 10px;
  }

  .file-name {
    background: rgba(255, 255, 255, 0.03);
    border: 1px solid var(--line);
    border-radius: 999px;
    padding: 6px 10px;
  }

  .file-input {
    clip: rect(0, 0, 0, 0);
    border: 0;
    clip-path: inset(50%);
    height: 1px;
    margin: -1px;
    overflow: hidden;
    padding: 0;
    pointer-events: none;
    position: absolute;
    white-space: nowrap;
    width: 1px;
  }

  .panel-body {
    display: grid;
    gap: 0.75rem;
  }

  .stats-grid {
    display: grid;
    gap: 10px;
    grid-template-columns: repeat(auto-fit, minmax(140px, 1fr));
  }

  .stats-grid div {
    background: rgba(0, 0, 0, 0.16);
    border: 0.5px solid rgba(255, 255, 255, 0.08);
    border-radius: 16px;
    display: grid;
    gap: 6px;
    min-width: 0;
    padding: 14px;
  }

  .button-secondary {
    background: rgba(255, 255, 255, 0.04);
    border: 0.5px solid var(--line);
    border-radius: 999px;
    color: var(--text);
    cursor: pointer;
    font: inherit;
    min-height: 40px;
    min-width: 0;
    overflow-wrap: anywhere;
    padding: 0 16px;
  }

  button:disabled {
    cursor: not-allowed;
    opacity: 0.5;
  }
</style>
