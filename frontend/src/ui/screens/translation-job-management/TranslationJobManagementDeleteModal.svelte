<script lang="ts">
  import type { TranslationJobManagementDeleteConfirmationViewModel } from "@application/contract/translation-job-management/translation-job-management-screen-types"

  interface Props {
    confirmation: TranslationJobManagementDeleteConfirmationViewModel | null
    onClose: () => void
    onConfirm: () => void
  }

  let { confirmation, onClose, onConfirm }: Props = $props()

  const isOpen = $derived(confirmation !== null)
</script>

<div
  aria-hidden={!isOpen}
  class="modal-backdrop"
  class:is-open={isOpen}
  hidden={!isOpen}
  id="jobDeleteModal"
  role="dialog"
>
  {#if confirmation}
    <section aria-labelledby="jobDeleteModalTitle" class="modal-card">
      <div class="section-head">
        <div>
          <p class="eyebrow">削除確認</p>
          <h3 id="jobDeleteModalTitle">{confirmation.title}</h3>
        </div>
        <button class="button-secondary" disabled={confirmation.busy} onclick={onClose} type="button">
          閉じる
        </button>
      </div>

      <p class="support-copy">
        削除対象と残るデータを確認してから実行してください。気づきやすさを優先して、確認はモーダルで表示しています。
      </p>

      <ul class="impact-list">
        {#each confirmation.lines as line (line)}
          <li>{line}</li>
        {/each}
      </ul>

      <div class="button-row">
        <button class="button-secondary" disabled={confirmation.busy} onclick={onClose} type="button">
          {confirmation.cancelLabel}
        </button>
        <button class="button-danger" disabled={confirmation.busy} onclick={onConfirm} type="button">
          {confirmation.busy ? "削除中..." : confirmation.confirmLabel}
        </button>
      </div>
    </section>
  {/if}
</div>

<style>
  .modal-backdrop {
    align-items: center;
    background: rgba(7, 7, 8, 0.72);
    display: none;
    inset: 0;
    justify-content: center;
    padding: 18px;
    position: fixed;
    z-index: 50;
  }

  .modal-backdrop.is-open {
    display: flex;
  }

  .modal-card {
    width: min(100%, 620px);
    display: grid;
    gap: 16px;
    padding: 24px;
    border-radius: 20px;
    border: 1px solid rgba(226, 205, 173, 0.14);
    background: rgba(17, 13, 12, 0.96);
    box-shadow: 0 24px 60px rgba(0, 0, 0, 0.42);
  }

  .section-head,
  .button-row {
    display: flex;
    gap: 12px;
    justify-content: space-between;
    align-items: flex-start;
    flex-wrap: wrap;
  }

  .eyebrow,
  .support-copy,
  .impact-list {
    margin: 0;
    color: rgba(236, 223, 205, 0.78);
  }

  .eyebrow {
    font-size: 0.76rem;
    letter-spacing: 0.12em;
    text-transform: uppercase;
  }

  h3 {
    margin: 0.25rem 0 0;
    color: #fff6ea;
  }

  .impact-list {
    padding-left: 1.2rem;
    display: grid;
    gap: 0.45rem;
  }

  button {
    min-height: 2.8rem;
    padding: 0.65rem 1rem;
    border-radius: 14px;
    border: 1px solid rgba(233, 213, 186, 0.18);
    font: inherit;
    cursor: pointer;
  }

  button:disabled {
    cursor: not-allowed;
    opacity: 0.55;
  }

  .button-secondary {
    background: rgba(255, 255, 255, 0.06);
    color: #fff6ea;
  }

  .button-danger {
    background: linear-gradient(135deg, #df6a4f 0%, #f4a16d 100%);
    color: #1b120c;
  }
</style>
