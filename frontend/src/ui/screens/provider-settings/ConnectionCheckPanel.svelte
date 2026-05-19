<script lang="ts">
  import type { ConnectionCheckPanelProps } from "./provider-settings-panel-props"

  let { selectedProvider, validateConnection }: ConnectionCheckPanelProps =
    $props()
</script>

<section class="detail-block" data-testid="provider-settings-connection-check-region">
  <div class="detail-row">
    <div>
      <p class="field-label">接続確認</p>
      <strong>{selectedProvider.validationLabel}</strong>
    </div>
    <button
      class="button-secondary"
      disabled={!selectedProvider.canValidate}
      onclick={validateConnection}
      type="button"
    >
      接続を確認
    </button>
  </div>
  <p class="mini-text">{selectedProvider.validationHelpText}</p>
</section>

<style>
  .detail-block {
    display: grid;
    gap: 0.5rem;
    margin-top: 1rem;
    color: var(--text, rgba(234, 225, 221, 0.92));
  }

  .detail-row {
    align-items: center;
    display: flex;
    flex-wrap: wrap;
    gap: 0.75rem;
    justify-content: space-between;
  }

  .field-label {
    color: var(--muted, rgba(216, 195, 174, 0.92));
    font-size: 0.78rem;
    letter-spacing: 0.08em;
    margin: 0 0 0.35rem;
    text-transform: uppercase;
  }

  .mini-text {
    color: var(--text, rgba(234, 225, 221, 0.92));
    margin: 0.35rem 0 0;
    overflow-wrap: anywhere;
  }

  .button-secondary {
    background: rgba(255, 241, 227, 0.08);
    border: 1px solid rgba(255, 212, 165, 0.18);
    border-radius: 999px;
    color: #ffe2bf;
    cursor: pointer;
    font: inherit;
    min-height: 2.8rem;
    padding: 0.7rem 1.15rem;
  }

  .button-secondary:disabled {
    cursor: not-allowed;
    opacity: 0.6;
  }
</style>
