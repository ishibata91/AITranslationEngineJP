<script lang="ts">
  import type { ProviderListPanelProps } from "./provider-settings-panel-props"

  let { providerList, selectProvider }: ProviderListPanelProps = $props()
</script>

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
    {#each providerList as provider (provider.providerId)}
      <button
        class="provider-row"
        class:is-selected={provider.selected}
        data-testid="provider-settings-ai-service-row"
        data-provider-id={provider.providerId}
        onclick={() => selectProvider(provider.providerId)}
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

<style>
  .provider-settings-card {
    border: 0.5px solid var(--line, rgba(255, 186, 56, 0.18));
    border-radius: 24px;
    background: var(--surface, rgba(35, 31, 29, 0.78));
    box-shadow: var(--shadow, 0 24px 64px rgba(0, 0, 0, 0.42));
    backdrop-filter: blur(38px);
    color: var(--text, rgba(234, 225, 221, 0.92));
    padding: 1.5rem;
  }

  .section-head {
    margin-bottom: 1rem;
  }

  .eyebrow {
    color: var(--muted, rgba(216, 195, 174, 0.92));
    font-size: 0.78rem;
    letter-spacing: 0.08em;
    margin: 0 0 0.35rem;
    text-transform: uppercase;
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

  .provider-row-head {
    align-items: center;
    display: flex;
    flex-wrap: wrap;
    gap: 0.75rem;
    justify-content: space-between;
  }

  .provider-row p {
    color: var(--text, rgba(234, 225, 221, 0.92));
    margin: 0.35rem 0 0;
    overflow-wrap: anywhere;
  }

  .status-pill {
    border-radius: 999px;
    display: inline-flex;
    font-size: 0.78rem;
    padding: 0.28rem 0.7rem;
    white-space: nowrap;
  }

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
</style>
