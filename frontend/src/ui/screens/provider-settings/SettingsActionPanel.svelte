<script lang="ts">
  import type { SettingsActionPanelProps } from "./provider-settings-panel-props"

  let { selectedProvider, saveSettings, resetSettings }: SettingsActionPanelProps =
    $props()
</script>

<section class="detail-actions" data-testid="provider-settings-settings-actions-region">
  <button
    class="button-primary"
    disabled={!selectedProvider.canSave}
    onclick={saveSettings}
    type="button"
  >
    設定を保存
  </button>
  <button
    class="button-secondary"
    disabled={!selectedProvider.canReset}
    onclick={resetSettings}
    type="button"
  >
    リセット
  </button>
</section>

<style>
  .detail-actions {
    display: flex;
    flex-wrap: wrap;
    gap: 0.75rem;
    margin-top: 1.5rem;
    color: var(--text, rgba(234, 225, 221, 0.92));
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
</style>
