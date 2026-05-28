<script lang="ts">
  import type { ProviderSettingsSummaryPanelProps } from "./provider-settings-panel-props"

  let {
    gatewayStatus,
    pageTitle,
    pageLead,
    providerCountLabel,
    phaseLabel,
    saveNotice,
    errorMessage
  }: ProviderSettingsSummaryPanelProps = $props()
</script>

<section
  class="provider-settings-card hero-panel"
  data-testid="provider-settings-screen-summary-region"
>
  <div class="hero-top">
    <div>
      <p class="eyebrow">中央設定</p>
      <h2>{pageTitle}</h2>
    </div>
    <div class="status-row">
      <span class="status-pill status-accent">{phaseLabel}</span>
      <span class="status-pill">Gateway: {gatewayStatus}</span>
    </div>
  </div>
  <p class="lead">{pageLead}</p>
  <p class="mini-text">{providerCountLabel}</p>
  <p class="notice-text" hidden={!saveNotice} role="status">{saveNotice}</p>
  <p
    class="error-text"
    data-testid="provider-settings-validation-error"
    hidden={!errorMessage}
  >
    {errorMessage}
  </p>
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

  .hero-top,
  .status-row {
    align-items: center;
    display: flex;
    flex-wrap: wrap;
    gap: 0.75rem;
    justify-content: space-between;
  }

  .eyebrow {
    color: var(--muted, rgba(216, 195, 174, 0.92));
    font-size: 0.78rem;
    letter-spacing: 0.08em;
    margin: 0 0 0.35rem;
    text-transform: uppercase;
  }

  .lead,
  .mini-text,
  .notice-text,
  .error-text {
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

  .status-accent {
    background: rgba(145, 208, 134, 0.16);
    color: #b8f0ad;
  }
</style>
