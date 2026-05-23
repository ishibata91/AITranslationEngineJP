<script lang="ts">
  import ActionButton from "./ActionButton.svelte"
  import ButtonGroup from "./ButtonGroup.svelte"
  import InlineFeedback from "./InlineFeedback.svelte"

  interface Props {
    open: boolean
    title: string
    targetLabel: string
    message: string
    confirmLabel?: string
    cancelLabel?: string
    busy?: boolean
    error?: string
    onConfirm: () => void
    onCancel: () => void
  }

  let {
    open,
    title,
    targetLabel,
    message,
    confirmLabel = "削除する",
    cancelLabel = "キャンセル",
    busy = false,
    error = "",
    onConfirm,
    onCancel
  }: Props = $props()
</script>

{#if open}
  <div class="modal-backdrop" role="presentation">
    <div
      class="confirm-danger-modal"
      role="dialog"
      aria-modal="true"
      aria-labelledby="confirm-danger-title"
    >
      <div class="modal-copy">
        <p class="modal-kicker">危険操作</p>
        <h2 id="confirm-danger-title">{title}</h2>
        <p>{message}</p>
      </div>

      <dl class="target-summary">
        <div>
          <dt>対象</dt>
          <dd>{targetLabel}</dd>
        </div>
      </dl>

      {#if error}
        <InlineFeedback tone="error" title="削除失敗" message={error} />
      {/if}

      <ButtonGroup align="end">
        <ActionButton
          label={cancelLabel}
          variant="secondary"
          disabled={busy}
          onClick={onCancel}
        />
        <ActionButton
          label={confirmLabel}
          variant="danger"
          {busy}
          onClick={onConfirm}
        />
      </ButtonGroup>
    </div>
  </div>
{/if}

<style>
  .modal-backdrop {
    align-items: center;
    background: rgba(15, 23, 42, 0.48);
    display: flex;
    inset: 0;
    justify-content: center;
    padding: 1rem;
    position: fixed;
    z-index: 30;
  }

  .confirm-danger-modal {
    background: #ffffff;
    border: 1px solid #cbd5e1;
    border-radius: 0.5rem;
    box-shadow: 0 24px 70px rgba(15, 23, 42, 0.24);
    display: grid;
    gap: 1rem;
    max-width: 34rem;
    padding: 1.2rem;
    width: min(100%, 34rem);
  }

  .modal-copy {
    display: grid;
    gap: 0.4rem;
  }

  .modal-kicker {
    color: #b91c1c;
    font-size: 0.78rem;
    font-weight: 900;
    text-transform: uppercase;
  }

  h2,
  p,
  dl {
    margin: 0;
    overflow-wrap: anywhere;
  }

  h2 {
    color: #172033;
    font-size: 1.2rem;
  }

  p {
    color: #475569;
  }

  .target-summary {
    background: #f8fafc;
    border: 1px solid #e2e8f0;
    border-radius: 0.5rem;
    padding: 0.8rem;
  }

  .target-summary div {
    display: grid;
    gap: 0.2rem;
  }

  dt {
    color: #64748b;
    font-size: 0.82rem;
    font-weight: 800;
  }

  dd {
    color: #172033;
    font-weight: 800;
    margin: 0;
    overflow-wrap: anywhere;
  }
</style>
