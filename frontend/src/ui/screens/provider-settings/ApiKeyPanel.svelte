<script lang="ts">
  import type { ApiKeyPanelProps } from "./provider-settings-panel-props"

  let {
    selectedProvider,
    credentialInputDraft,
    updateCredentialDraft,
    openApiKeyPanel,
    saveSettings,
    closeApiKeyPanel
  }: ApiKeyPanelProps = $props()

  function handleCredentialInput(event: Event): void {
    const target = event.currentTarget
    if (target instanceof HTMLInputElement) {
      updateCredentialDraft(target.value)
    }
  }
</script>

<section
  class="detail-block"
  data-testid="provider-settings-api-key-status-region"
>
  <div class="detail-row">
    <div>
      <p class="field-label">APIキー状態</p>
      <strong>{selectedProvider.apiKeyStateLabel}</strong>
    </div>
    {#if selectedProvider.apiKeyRequired}
      <button class="button-secondary" onclick={openApiKeyPanel} type="button">
        設定
      </button>
    {/if}
  </div>
  <p class="mini-text">{selectedProvider.apiKeyHelpText}</p>

  {#if selectedProvider.apiKeyRequired && selectedProvider.apiKeyPanelOpen}
    <div
      class="api-key-panel"
      data-testid="provider-settings-api-key-input-region"
    >
      <label class="field-group" for="providerApiKeyInput">
        <span class="field-label">APIキー</span>
        <input
          autocomplete="off"
          class="text-field"
          id="providerApiKeyInput"
          oninput={handleCredentialInput}
          type="password"
          value={credentialInputDraft}
        />
      </label>
      <div
        class="inline-actions"
        data-testid="provider-settings-api-key-save-actions"
      >
        <button class="button-primary" onclick={saveSettings} type="button">
          保存
        </button>
        <button
          class="button-secondary"
          onclick={closeApiKeyPanel}
          type="button"
        >
          キャンセル
        </button>
      </div>
    </div>
  {/if}
</section>

<style>
  .detail-block,
  .field-group {
    display: grid;
    gap: 0.5rem;
    margin-top: 1rem;
    color: var(--text, rgba(234, 225, 221, 0.92));
  }

  .detail-row,
  .inline-actions {
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

  .api-key-panel {
    background: rgba(255, 241, 227, 0.05);
    border: 1px solid rgba(255, 212, 165, 0.18);
    border-radius: 18px;
    display: grid;
    gap: 0.85rem;
    margin-top: 0.6rem;
    padding: 1rem;
  }

  .text-field {
    background: rgba(18, 13, 11, 0.92);
    border: 1px solid rgba(255, 212, 165, 0.18);
    border-radius: 14px;
    color: #fef3e8;
    font: inherit;
    min-height: 2.9rem;
    padding: 0.8rem 0.95rem;
    width: 100%;
  }

  .button-primary,
  .button-secondary {
    border-radius: 999px;
    cursor: pointer;
    font: inherit;
    min-height: 2.8rem;
    padding: 0.7rem 1.15rem;
  }

  .button-primary {
    background: linear-gradient(135deg, var(--primary, #ffba38), #f3a114);
    border: none;
    color: #432c00;
    font-weight: 700;
  }

  .button-secondary {
    background: rgba(255, 241, 227, 0.08);
    border: 1px solid rgba(255, 212, 165, 0.18);
    color: #ffe2bf;
  }
</style>
