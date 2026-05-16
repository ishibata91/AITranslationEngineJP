<script lang="ts">
  import { onMount } from "svelte"

  import type {
    CreateProviderSettingsScreenController,
    ProviderSettingsScreenControllerContract
  } from "@application/contract/provider-settings"

  interface ProviderSettingsPageControllerContract
    extends ProviderSettingsScreenControllerContract {
    updateCredentialInput(nextValue: string): void
    clearCredentialInput(): void
  }

  interface Props {
    createController: CreateProviderSettingsScreenController | null
  }

  let { createController }: Props = $props()

  function resolveController(): ProviderSettingsPageControllerContract {
    if (!createController) {
      throw new Error(
        "provider settings screen controller factory is not provided"
      )
    }

    return createController() as ProviderSettingsPageControllerContract
  }

  const controller = resolveController()
  let viewModel = $state(controller.getViewModel())
  let credentialInputDraft = $state("")

  const unsubscribe = controller.subscribe((nextViewModel) => {
    viewModel = nextViewModel
    if (!viewModel.selectedProvider?.apiKeyPanelOpen) {
      credentialInputDraft = ""
      controller.clearCredentialInput()
    }
  })

  onMount(() => {
    void controller.mount()

    return () => {
      unsubscribe()
      controller.dispose()
    }
  })

  function saveSettings(): void {
    void controller.saveSettings()
  }

  function closeApiKeyPanel(): void {
    credentialInputDraft = ""
    controller.clearCredentialInput()
    controller.closeApiKeyPanel()
  }
</script>

<section
  class="provider-settings-shell"
  data-testid="provider-settings-screen-shell"
  id="providerSettingsView"
>
  <section
    class="provider-settings-card hero-panel"
    data-testid="provider-settings-screen-summary-region"
  >
    <div class="hero-top">
      <div>
        <p class="eyebrow">中央設定</p>
        <h2>{viewModel.pageTitle}</h2>
      </div>
      <div class="status-row">
        <span class="status-pill status-accent">{viewModel.phaseLabel}</span>
        <span class="status-pill">Gateway: {viewModel.gatewayStatus}</span>
      </div>
    </div>
    <p class="lead">{viewModel.pageLead}</p>
    <p class="mini-text">{viewModel.providerCountLabel}</p>
    <p class="notice-text" hidden={!viewModel.saveNotice} role="status">
      {viewModel.saveNotice}
    </p>
    <p class="error-text" hidden={!viewModel.errorMessage}>
      {viewModel.errorMessage}
    </p>
  </section>

  <section class="provider-layout">
    <section
      class="provider-settings-card list-panel"
      data-testid="provider-settings-ai-service-list-region"
    >
      <div class="section-head">
        <div>
          <p class="eyebrow">AIサービス一覧</p>
          <h3>設定状態を比較</h3>
        </div>
      </div>

      <div class="provider-list">
        {#each viewModel.providerList as provider (provider.providerId)}
          <button
            class="provider-row"
            class:is-selected={provider.selected}
            data-testid="provider-settings-ai-service-row"
            onclick={() => controller.selectProvider(provider.providerId)}
            type="button"
          >
            <div class="provider-row-head">
              <strong>{provider.label}</strong>
              <span class={`status-pill tone-${provider.statusTone}`}>
                {provider.statusLabel}
              </span>
            </div>
            <p>{provider.helperText}</p>
          </button>
        {/each}
      </div>
    </section>

    {#if viewModel.selectedProvider}
      <section
        class="provider-settings-card detail-panel"
        data-testid="provider-settings-settings-detail-region"
      >
        <div class="section-head">
          <div>
            <p class="eyebrow">設定詳細</p>
            <h3>{viewModel.selectedProvider.label}</h3>
          </div>
          <span
            class={`status-pill tone-${viewModel.selectedProvider.validationTone}`}
          >
            {viewModel.selectedProvider.validationLabel}
          </span>
        </div>

        <label class="field-group" for="providerEndpointInput">
          <span class="field-label">エンドポイント</span>
          <input
            class="text-field"
            data-testid="provider-settings-endpoint-input"
            id="providerEndpointInput"
            oninput={(event) => controller.updateEndpoint(event)}
            placeholder={viewModel.selectedProvider.endpointPlaceholder}
            type="text"
            value={viewModel.selectedProvider.endpoint}
          />
        </label>
        <p class="mini-text">{viewModel.selectedProvider.endpointHint}</p>

        <section
          class="detail-block"
          data-testid="provider-settings-api-key-status-region"
        >
          <div class="detail-row">
            <div>
              <p class="field-label">APIキー状態</p>
              <strong>{viewModel.selectedProvider.apiKeyStateLabel}</strong>
            </div>
            {#if viewModel.selectedProvider.apiKeyRequired}
              <button
                class="button-secondary"
                onclick={() => {
                  controller.openApiKeyPanel()
                }}
                type="button"
              >
                設定
              </button>
            {/if}
          </div>
          <p class="mini-text">{viewModel.selectedProvider.apiKeyHelpText}</p>

          {#if viewModel.selectedProvider.apiKeyRequired && viewModel.selectedProvider.apiKeyPanelOpen}
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
                  bind:value={credentialInputDraft}
                  oninput={(event) => {
                    const target = event.currentTarget
                    if (target instanceof HTMLInputElement) {
                      controller.updateCredentialInput(target.value)
                    }
                  }}
                  type="password"
                />
              </label>
              <div
                class="inline-actions"
                data-testid="provider-settings-api-key-save-actions"
              >
                <button
                  class="button-primary"
                  onclick={saveSettings}
                  type="button"
                >
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

        <section
          class="detail-block"
          data-testid="provider-settings-connection-check-region"
        >
          <div class="detail-row">
            <div>
              <p class="field-label">接続確認</p>
              <strong>{viewModel.selectedProvider.validationLabel}</strong>
            </div>
            <button
              class="button-secondary"
              disabled={!viewModel.selectedProvider.canValidate}
              onclick={() => void controller.validateConnection()}
              type="button"
            >
              接続を確認
            </button>
          </div>
          <p class="mini-text">
            {viewModel.selectedProvider.validationHelpText}
          </p>
        </section>

        <section
          class="detail-actions"
          data-testid="provider-settings-settings-actions-region"
        >
          <button
            class="button-primary"
            disabled={!viewModel.selectedProvider.canSave}
            onclick={saveSettings}
            type="button"
          >
            設定を保存
          </button>
          <button
            class="button-secondary"
            disabled={!viewModel.selectedProvider.canReset}
            onclick={() => {
              credentialInputDraft = ""
              controller.clearCredentialInput()
              void controller.resetSettings()
            }}
            type="button"
          >
            リセット
          </button>
        </section>
      </section>
    {/if}
  </section>
</section>

<style>
  .provider-settings-shell {
    display: grid;
    gap: 1.5rem;
  }

  .provider-settings-card {
    border: 0.5px solid var(--line, rgba(255, 186, 56, 0.18));
    border-radius: 24px;
    background: var(--surface, rgba(35, 31, 29, 0.78));
    box-shadow: var(--shadow, 0 24px 64px rgba(0, 0, 0, 0.42));
    backdrop-filter: blur(38px);
    padding: 1.5rem;
  }

  .provider-layout {
    display: grid;
    gap: 1.5rem;
    grid-template-columns: minmax(0, 340px) minmax(0, 1fr);
  }

  .hero-top,
  .detail-row,
  .provider-row-head,
  .status-row {
    align-items: center;
    display: flex;
    flex-wrap: wrap;
    gap: 0.75rem;
    justify-content: space-between;
  }

  .section-head {
    margin-bottom: 1rem;
  }

  .eyebrow,
  .field-label {
    color: var(--muted, rgba(216, 195, 174, 0.92));
    font-size: 0.78rem;
    letter-spacing: 0.08em;
    margin: 0 0 0.35rem;
    text-transform: uppercase;
  }

  .lead,
  .mini-text,
  .provider-row p,
  .notice-text,
  .error-text {
    color: var(--text, rgba(234, 225, 221, 0.92));
    margin: 0.35rem 0 0;
    overflow-wrap: anywhere;
  }

  .provider-list {
    display: grid;
    gap: 0.85rem;
  }

  .provider-row {
    background: rgba(255, 255, 255, 0.03);
    border: 1px solid rgba(255, 186, 56, 0.12);
    border-radius: 18px;
    color: inherit;
    cursor: pointer;
    display: grid;
    gap: 0.55rem;
    padding: 1rem;
    text-align: left;
  }

  .provider-row.is-selected {
    border-color: var(--line-strong, rgba(255, 186, 56, 0.32));
    box-shadow: 0 0 0 1px rgba(255, 186, 56, 0.16);
  }

  .status-pill {
    border-radius: 999px;
    display: inline-flex;
    font-size: 0.78rem;
    padding: 0.28rem 0.7rem;
    white-space: nowrap;
  }

  .status-accent,
  .tone-success {
    background: rgba(145, 208, 134, 0.16);
    color: #b8f0ad;
  }

  .tone-warning {
    background: rgba(255, 204, 128, 0.15);
    color: #ffd191;
  }

  .tone-neutral {
    background: rgba(255, 190, 126, 0.14);
    color: #ffd8ae;
  }

  .field-group,
  .detail-block {
    display: grid;
    gap: 0.5rem;
    margin-top: 1rem;
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

  .api-key-panel {
    border: 1px solid rgba(255, 212, 165, 0.18);
    border-radius: 18px;
    display: grid;
    gap: 0.85rem;
    margin-top: 0.6rem;
    padding: 1rem;
    background: rgba(255, 241, 227, 0.05);
  }

  .inline-actions,
  .detail-actions {
    display: flex;
    flex-wrap: wrap;
    gap: 0.75rem;
  }

  .detail-actions {
    margin-top: 1.5rem;
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

  .button-primary:disabled,
  .button-secondary:disabled {
    cursor: not-allowed;
    opacity: 0.6;
  }

  @media (max-width: 840px) {
    .provider-layout {
      grid-template-columns: 1fr;
    }
  }

  @media (max-width: 480px) {
    .provider-settings-card {
      padding: 1rem;
    }

    .button-primary,
    .button-secondary {
      width: 100%;
      justify-content: center;
    }

    .detail-actions {
      display: grid;
    }
  }
</style>
