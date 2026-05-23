<script lang="ts">
  import { onMount } from "svelte"

  import type {
    CreateProviderSettingsScreenController,
    ProviderSettingsScreenControllerContract
  } from "@application/contract/provider-settings"
  import ApiKeyPanel from "./ApiKeyPanel.svelte"
  import ConnectionCheckPanel from "./ConnectionCheckPanel.svelte"
  import ProviderDetailPanel from "./ProviderDetailPanel.svelte"
  import ProviderListPanel from "./ProviderListPanel.svelte"
  import ProviderSettingsSummaryPanel from "./ProviderSettingsSummaryPanel.svelte"
  import SettingsActionPanel from "./SettingsActionPanel.svelte"

  interface ProviderSettingsPageControllerContract extends ProviderSettingsScreenControllerContract {
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

  function updateCredentialDraft(nextValue: string): void {
    credentialInputDraft = nextValue
    controller.updateCredentialInput(nextValue)
  }

  function resetSettings(): void {
    credentialInputDraft = ""
    controller.clearCredentialInput()
    void controller.resetSettings()
  }

  function selectProvider(providerId: string): void {
    controller.selectProvider(providerId)
  }

  function updateEndpoint(event: Event): void {
    controller.updateEndpoint(event)
  }

  function openApiKeyPanel(): void {
    controller.openApiKeyPanel()
  }

  function validateConnection(): void {
    void controller.validateConnection()
  }
</script>

<section
  class="provider-settings-shell"
  data-testid="provider-settings-screen-shell"
  id="providerSettingsView"
>
  <ProviderSettingsSummaryPanel
    errorMessage={viewModel.errorMessage}
    gatewayStatus={viewModel.gatewayStatus}
    pageLead={viewModel.pageLead}
    pageTitle={viewModel.pageTitle}
    phaseLabel={viewModel.phaseLabel}
    providerCountLabel={viewModel.providerCountLabel}
    saveNotice={viewModel.saveNotice}
  />

  <section class="provider-layout">
    <ProviderListPanel providerList={viewModel.providerList} {selectProvider} />

    {#if viewModel.selectedProvider}
      <section
        class="provider-settings-card detail-panel"
        data-testid="provider-settings-settings-detail-region"
      >
        <ProviderDetailPanel
          selectedProvider={viewModel.selectedProvider}
          {updateEndpoint}
        />
        <ApiKeyPanel
          {closeApiKeyPanel}
          {credentialInputDraft}
          {openApiKeyPanel}
          {saveSettings}
          selectedProvider={viewModel.selectedProvider}
          {updateCredentialDraft}
        />
        <ConnectionCheckPanel
          selectedProvider={viewModel.selectedProvider}
          {validateConnection}
        />
        <SettingsActionPanel
          {resetSettings}
          {saveSettings}
          selectedProvider={viewModel.selectedProvider}
        />
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
    color: var(--text, rgba(234, 225, 221, 0.92));
    padding: 1.5rem;
  }

  .provider-layout {
    display: grid;
    gap: 1.5rem;
    grid-template-columns: minmax(0, 340px) minmax(0, 1fr);
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
  }
</style>
