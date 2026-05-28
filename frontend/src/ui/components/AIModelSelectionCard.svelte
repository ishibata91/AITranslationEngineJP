<script lang="ts">
  type SelectOption = { value: string; label: string }
  type ModelOption = { modelId: string; label: string }
  type Tone = "neutral" | "warning" | "success"
  type SecondaryControlMode = "batch-toggle" | "execution-select" | "none"

  interface Props {
    eyebrow?: string
    title: string
    titleId?: string
    titleTag?: "h3" | "h4"
    dataTestId?: string
    ariaLabel?: string
    helperText?: string
    statusLabel?: string
    statusTestId?: string
    statusTone?: Tone
    providerSelectId: string
    providerSelectTestId?: string
    providerValue: string
    providerOptions: SelectOption[]
    providerFieldLabel?: string
    providerDisabled?: boolean
    onProviderChange: (event: Event) => void
    credentialStatusLabel?: string
    credentialStatusTone?: Tone
    showCredentialStatus?: boolean
    showCredentialWarning?: boolean
    credentialWarningText?: string
    credentialEmptyText?: string
    secondaryControlMode?: SecondaryControlMode
    batchChecked?: boolean
    batchDisabled?: boolean
    batchHelpId?: string
    batchHelpText?: string
    onBatchChange?: (event: Event) => void
    executionSelectId?: string
    executionSelectTestId?: string
    executionValue?: string
    executionOptions?: SelectOption[]
    executionDisabled?: boolean
    onExecutionChange?: (event: Event) => void
    modelStatusText?: string
    refreshButtonLabel?: string
    refreshButtonAriaLabel?: string
    refreshDisabled?: boolean
    refreshSpinning?: boolean
    onRefresh?: (() => void) | null
    modelSelectId: string
    modelSelectTestId?: string
    modelValue: string
    modelOptions: ModelOption[]
    modelDisabled?: boolean
    emptyModelLabel?: string
    onModelChange: (event: Event) => void
    footerMessage?: string
    footerMessageId?: string
    footerWarningText?: string
    actionButtonLabel?: string
    actionButtonId?: string
    actionButtonTestId?: string
    actionButtonDisabled?: boolean
    onAction?: (() => void) | null
  }

  let {
    eyebrow = "",
    title,
    titleId = undefined,
    titleTag = "h4",
    dataTestId = undefined,
    ariaLabel = undefined,
    helperText = "",
    statusLabel = "",
    statusTestId = undefined,
    statusTone = "neutral",
    providerSelectId,
    providerSelectTestId = undefined,
    providerValue,
    providerOptions,
    providerFieldLabel = "AIサービス",
    providerDisabled = false,
    onProviderChange,
    credentialStatusLabel = "",
    credentialStatusTone = "neutral",
    showCredentialStatus = true,
    showCredentialWarning = false,
    credentialWarningText = "AIサービス設定が未完了です。設定が必要です。",
    credentialEmptyText = "登録は不要です。",
    secondaryControlMode = "none",
    batchChecked = false,
    batchDisabled = false,
    batchHelpId = undefined,
    batchHelpText = "",
    onBatchChange = undefined,
    executionSelectId = undefined,
    executionSelectTestId = undefined,
    executionValue = "",
    executionOptions = [],
    executionDisabled = false,
    onExecutionChange = undefined,
    modelStatusText = "",
    refreshButtonLabel = "モデル一覧を更新",
    refreshButtonAriaLabel = "モデル一覧を更新",
    refreshDisabled = false,
    refreshSpinning = false,
    onRefresh = null,
    modelSelectId,
    modelSelectTestId = undefined,
    modelValue,
    modelOptions,
    modelDisabled = false,
    emptyModelLabel = "選んでください",
    onModelChange,
    footerMessage = "",
    footerMessageId = undefined,
    footerWarningText = "",
    actionButtonLabel = "",
    actionButtonId = undefined,
    actionButtonTestId = undefined,
    actionButtonDisabled = false,
    onAction = null
  }: Props = $props()
</script>

<article class="model-card" aria-label={ariaLabel} data-testid={dataTestId}>
  <div class="model-card-head section-head">
    <div class="heading-copy">
      {#if eyebrow}
        <p class="eyebrow">{eyebrow}</p>
      {/if}
      <svelte:element this={titleTag} id={titleId}>{title}</svelte:element>
      {#if helperText}
        <p class="helper-text">{helperText}</p>
      {/if}
    </div>
    {#if statusLabel}
      <span
        class={`status-pill ${
          statusTone === "warning"
            ? "status-alert"
            : statusTone === "success"
              ? "status-success"
              : ""
        }`}
        data-testid={statusTestId}
      >
        {statusLabel}
      </span>
    {/if}
  </div>

  <div class="model-edit-grid">
    <label class="model-field" for={providerSelectId}>
      <span>{providerFieldLabel}</span>
      <select
        aria-label={providerFieldLabel}
        data-testid={providerSelectTestId}
        disabled={providerDisabled}
        id={providerSelectId}
        onchange={onProviderChange}
        value={providerValue}
      >
        {#each providerOptions as option (option.value)}
          <option value={option.value}>{option.label}</option>
        {/each}
      </select>
      {#if showCredentialStatus}
        <p
          class={credentialStatusTone === "warning"
            ? "warning-text"
            : "support-text"}
        >
          {credentialStatusLabel}
        </p>
      {:else}
        <p class="support-text">{credentialEmptyText}</p>
      {/if}
      {#if showCredentialWarning}
        <small class="warning-text">{credentialWarningText}</small>
      {/if}
    </label>

    <div class="model-field">
      <span id={`${modelSelectId}-label`}>モデル</span>
      <div class="model-select-row">
        <select
          aria-labelledby={`${modelSelectId}-label`}
          data-testid={modelSelectTestId}
          disabled={modelDisabled}
          id={modelSelectId}
          onchange={onModelChange}
          value={modelValue}
        >
          <option value="">{emptyModelLabel}</option>
          {#each modelOptions as modelOption (modelOption.modelId)}
            <option value={modelOption.modelId}>{modelOption.label}</option>
          {/each}
        </select>
        {#if onRefresh}
          <button
            aria-label={refreshButtonAriaLabel}
            class="icon-button"
            disabled={refreshDisabled}
            onclick={() => onRefresh?.()}
            type="button"
          >
            <span
              aria-hidden="true"
              class={`refresh-icon ${refreshSpinning ? "is-spinning" : ""}`}
            >
              ↻
            </span>
            <span class="sr-only">{refreshButtonLabel}</span>
          </button>
        {/if}
      </div>
      {#if modelStatusText}
        <small class={modelDisabled ? "warning-text" : "support-text"}>
          {modelStatusText}
        </small>
      {/if}
    </div>

    <div class="model-field">
      {#if secondaryControlMode === "batch-toggle"}
        <span class="field-title-row">
          処理方式
          <button
            aria-describedby={batchHelpId}
            class="help-tooltip"
            type="button"
          >
            ?
            {#if batchHelpText}
              <span class="tooltip-panel" id={batchHelpId} role="tooltip">
                {batchHelpText}
              </span>
            {/if}
          </button>
        </span>
        <label class="checkbox-row">
          <input
            checked={batchChecked}
            disabled={batchDisabled}
            onchange={onBatchChange}
            type="checkbox"
          />
          <span>Batch API</span>
        </label>
      {:else if secondaryControlMode === "execution-select"}
        <label class="field-block" for={executionSelectId}>
          <span>処理方式</span>
          <select
            aria-label="処理方式"
            data-testid={executionSelectTestId}
            disabled={executionDisabled}
            id={executionSelectId}
            onchange={onExecutionChange}
            value={executionValue}
          >
            {#each executionOptions as option (option.value)}
              <option value={option.value}>{option.label}</option>
            {/each}
          </select>
        </label>
      {:else}
        <span>処理方式</span>
        <p class="support-text">
          この AI サービスでは一括処理の切り替えはありません。
        </p>
      {/if}
    </div>
  </div>

  {#if footerMessage || footerWarningText || onAction}
    <div class="footer-row">
      <div class="footer-copy">
        {#if footerMessage}
          <span class="support-text" id={footerMessageId}>{footerMessage}</span>
        {/if}
        {#if footerWarningText}
          <span class="status-pill status-alert">{footerWarningText}</span>
        {/if}
      </div>
      {#if onAction && actionButtonLabel}
        <button
          id={actionButtonId}
          class="button-secondary"
          data-testid={actionButtonTestId}
          disabled={actionButtonDisabled}
          onclick={() => onAction?.()}
          type="button"
        >
          {actionButtonLabel}
        </button>
      {/if}
    </div>
  {/if}
</article>

<style>
  .model-card,
  .field-block,
  .footer-copy,
  .heading-copy {
    display: grid;
    gap: 0.5rem;
  }

  .model-card {
    background: rgba(33, 27, 24, 0.88);
    border: 1px solid rgba(226, 205, 173, 0.14);
    border-radius: 20px;
    box-shadow: 0 18px 40px rgba(0, 0, 0, 0.22);
    display: grid;
    gap: 1rem;
    grid-template-rows: auto 1fr;
    min-width: 0;
    padding: 1.4rem;
  }

  .model-card-head,
  .footer-row,
  .field-title-row {
    align-items: flex-start;
    display: flex;
    gap: 1rem;
    justify-content: space-between;
  }

  .model-card-head :global(h3),
  .model-card-head :global(h4) {
    font-size: 1rem;
    line-height: 1.4;
    margin: 0;
  }

  .helper-text,
  .support-text,
  .warning-text {
    line-height: 1.5;
  }

  .eyebrow,
  .model-field span,
  .field-block span {
    color: var(--primary, #ffba38);
    font-size: 0.75rem;
  }

  .status-pill {
    font-size: 0.82rem;
    padding: 0.35rem 0.72rem;
    border-radius: 999px;
    background: rgba(255, 255, 255, 0.08);
    border: 1px solid rgba(255, 255, 255, 0.12);
    color: var(--text, #fef3e8);
    white-space: nowrap;
  }

  .status-success {
    background: rgba(145, 208, 134, 0.16);
    border-color: rgba(145, 208, 134, 0.22);
    color: #b8f0ad;
  }

  .status-alert {
    background: rgba(255, 204, 128, 0.15);
    border-color: rgba(255, 204, 128, 0.28);
    color: #ffd191;
  }

  .model-edit-grid {
    display: grid;
    gap: 8px;
  }

  .model-field {
    background: rgba(0, 0, 0, 0.14);
    border: 1px solid rgba(255, 255, 255, 0.08);
    border-radius: 8px;
    display: grid;
    gap: 6px;
    min-width: 0;
    padding: 10px;
  }

  select {
    width: 100%;
    min-height: 36px;
    min-width: 0;
    padding: 7px 9px;
    border: 1px solid rgba(255, 255, 255, 0.1);
    border-radius: 6px;
    background: rgba(255, 255, 255, 0.08);
    color: var(--text, #fef3e8);
  }

  .checkbox-row {
    align-items: center;
    color: var(--text, #fef3e8);
    display: inline-flex;
    gap: 8px;
    min-height: 36px;
  }

  .checkbox-row input {
    accent-color: var(--primary, #ffba38);
    height: 16px;
    margin: 0;
    width: 16px;
  }

  .checkbox-row span {
    color: var(--text, #fef3e8);
    font-size: 0.9rem;
  }

  .help-tooltip {
    align-items: center;
    background: transparent;
    border: 1px solid rgba(255, 255, 255, 0.12);
    border-radius: 50%;
    color: rgba(252, 241, 232, 0.72);
    cursor: help;
    display: inline-flex;
    font-size: 0.72rem;
    height: 18px;
    justify-content: center;
    position: relative;
    width: 18px;
  }

  .tooltip-panel {
    background: rgba(22, 19, 18, 0.98);
    border: 1px solid rgba(255, 255, 255, 0.16);
    border-radius: 6px;
    bottom: calc(100% + 8px);
    color: var(--text, #fef3e8);
    display: none;
    font-size: 0.8rem;
    left: 50%;
    line-height: 1.5;
    min-width: 176px;
    padding: 8px 10px;
    position: absolute;
    transform: translateX(-50%);
    z-index: 30;
  }

  .help-tooltip:hover .tooltip-panel,
  .help-tooltip:focus .tooltip-panel {
    display: block;
  }

  .button-secondary {
    background: rgba(255, 255, 255, 0.06);
    border: 1px solid rgba(255, 255, 255, 0.12);
    border-radius: 7px;
    color: var(--text, #fef3e8);
    cursor: pointer;
    min-height: 40px;
    padding: 9px 13px;
  }

  .button-secondary:disabled,
  select:disabled,
  input:disabled {
    opacity: 0.56;
    cursor: not-allowed;
  }

  .icon-button {
    align-items: center;
    background: rgba(255, 255, 255, 0.06);
    border: 1px solid rgba(255, 255, 255, 0.12);
    border-radius: 6px;
    color: var(--text, #fef3e8);
    display: inline-flex;
    font-size: 1rem;
    height: 36px;
    justify-content: center;
    line-height: 1;
    padding: 0;
    width: 36px;
  }

  .refresh-icon {
    display: inline-block;
    font-size: 1.35rem;
    line-height: 1;
  }

  .refresh-icon.is-spinning {
    animation: rotation 0.7s linear infinite;
  }

  .model-select-row {
    align-items: center;
    display: grid;
    gap: 8px;
    grid-template-columns: minmax(0, 1fr) 36px;
  }

  .support-text {
    color: rgba(252, 241, 232, 0.78);
  }

  .warning-text {
    color: #ffd38a;
  }

  .sr-only {
    border: 0;
    clip: rect(0 0 0 0);
    height: 1px;
    margin: -1px;
    overflow: hidden;
    padding: 0;
    position: absolute;
    width: 1px;
  }

  @keyframes rotation {
    from {
      transform: rotate(0deg);
    }

    to {
      transform: rotate(360deg);
    }
  }

  @media (max-width: 720px) {
    .model-card-head,
    .footer-row {
      flex-direction: column;
    }
  }
</style>
