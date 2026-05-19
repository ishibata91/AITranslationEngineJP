<script lang="ts">
  import type { ProviderDetailPanelProps } from "./provider-settings-panel-props"

  let { selectedProvider, updateEndpoint }: ProviderDetailPanelProps = $props()
</script>

<section class="detail-block">
  <div class="section-head">
    <div>
      <p class="eyebrow">設定詳細</p>
      <h3>{selectedProvider.label}</h3>
    </div>
    <span class={`status-pill tone-${selectedProvider.validationTone}`}>
      {selectedProvider.validationLabel}
    </span>
  </div>

  <label class="field-group" for="providerEndpointInput">
    <span class="field-label">エンドポイント</span>
    <input
      class="text-field"
      data-testid="provider-settings-endpoint-input"
      id="providerEndpointInput"
      oninput={updateEndpoint}
      placeholder={selectedProvider.endpointPlaceholder}
      type="text"
      value={selectedProvider.endpoint}
    />
  </label>
  <p class="mini-text">{selectedProvider.endpointHint}</p>
</section>

<style>
  .detail-block,
  .field-group {
    display: grid;
    gap: 0.5rem;
    color: var(--text, rgba(234, 225, 221, 0.92));
  }

  .section-head {
    align-items: center;
    display: flex;
    flex-wrap: wrap;
    gap: 0.75rem;
    justify-content: space-between;
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

  .mini-text {
    color: var(--text, rgba(234, 225, 221, 0.92));
    margin: 0.35rem 0 0;
    overflow-wrap: anywhere;
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
