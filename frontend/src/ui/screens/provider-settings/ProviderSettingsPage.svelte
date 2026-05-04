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

<section class="provider-settings-shell" id="providerSettingsView">
  <section class="provider-settings-card hero-panel">
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
    <section class="provider-settings-card list-panel">
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
      <section class="provider-settings-card detail-panel">
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
            id="providerEndpointInput"
            oninput={(event) => controller.updateEndpoint(event)}
            placeholder={viewModel.selectedProvider.endpointPlaceholder}
            type="text"
            value={viewModel.selectedProvider.endpoint}
          />
        </label>
        <p class="mini-text">{viewModel.selectedProvider.endpointHint}</p>

        <section class="detail-block">
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
            <div class="api-key-panel">
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
              <div class="inline-actions">
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

        <section class="detail-block">
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

        <section class="detail-actions">
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
    border: 1px solid rgba(142, 162, 185, 0.16);
    border-radius: 24px;
    background:
      linear-gradient(180deg, rgba(11, 18, 32, 0.94), rgba(7, 11, 20, 0.98));
    box-shadow: 0 20px 40px rgba(2, 6, 23, 0.18);
    padding: 1.5rem;
  }

  .provider-layout {
    display: grid;
    gap: 1.5rem;
    grid-template-columns: minmax(260px, 340px) minmax(0, 1fr);
  }

  .hero-top,
  .detail-row,
  .provider-row-head,
  .status-row {
    align-items: center;
    display: flex;
    gap: 0.75rem;
    justify-content: space-between;
  }

  .section-head {
    margin-bottom: 1rem;
  }

  .eyebrow,
  .field-label {
    color: rgba(191, 219, 254, 0.84);
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
    color: rgba(226, 232, 240, 0.86);
    margin: 0.35rem 0 0;
  }

  .provider-list {
    display: grid;
    gap: 0.85rem;
  }

  .provider-row {
    background: rgba(15, 23, 42, 0.88);
    border: 1px solid rgba(148, 163, 184, 0.18);
    border-radius: 18px;
    color: inherit;
    cursor: pointer;
    display: grid;
    gap: 0.55rem;
    padding: 1rem;
    text-align: left;
  }

  .provider-row.is-selected {
    border-color: rgba(125, 211, 252, 0.56);
    box-shadow: 0 0 0 1px rgba(125, 211, 252, 0.26);
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
    background: rgba(22, 163, 74, 0.18);
    color: #bbf7d0;
  }

  .tone-warning {
    background: rgba(217, 119, 6, 0.18);
    color: #fde68a;
  }

  .tone-neutral {
    background: rgba(59, 130, 246, 0.16);
    color: #bfdbfe;
  }

  .field-group,
  .detail-block {
    display: grid;
    gap: 0.5rem;
    margin-top: 1rem;
  }

  .text-field {
    background: rgba(15, 23, 42, 0.9);
    border: 1px solid rgba(148, 163, 184, 0.24);
    border-radius: 14px;
    color: #e2e8f0;
    font: inherit;
    min-height: 2.9rem;
    padding: 0.8rem 0.95rem;
    width: 100%;
  }

  .api-key-panel {
    border: 1px solid rgba(148, 163, 184, 0.18);
    border-radius: 18px;
    display: grid;
    gap: 0.85rem;
    margin-top: 0.6rem;
    padding: 1rem;
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
    background: linear-gradient(135deg, #22c55e, #16a34a);
    border: none;
    color: #0b1120;
    font-weight: 700;
  }

  .button-secondary {
    background: transparent;
    border: 1px solid rgba(148, 163, 184, 0.28);
    color: #e2e8f0;
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
</style>
