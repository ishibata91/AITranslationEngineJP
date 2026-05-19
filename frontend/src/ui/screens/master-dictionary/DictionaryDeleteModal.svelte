<script lang="ts">
  import type { DictionaryDeleteModalProps } from "./dictionary-panel-props"

  let {
    modalState,
    selectedEntry,
    closeDeleteModal,
    deleteCurrentEntry
  }: DictionaryDeleteModalProps = $props()
</script>

<div
  aria-hidden={modalState !== "delete"}
  class="modal-backdrop"
  hidden={modalState !== "delete"}
  id="deleteModal"
  role="dialog"
  data-testid="master-dictionary-delete-confirmation-modal"
>
  <section aria-labelledby="deleteModalTitle" class="modal-card">
    <h3 id="deleteModalTitle">削除の確認</h3>
    <p>このエントリを削除すると、一覧から見えなくなります。</p>
    <div class="delete-target">
      <strong id="deleteTargetTitle">{selectedEntry?.source ?? "-"}</strong>
      <p id="deleteTargetMeta">
        {selectedEntry
          ? `${selectedEntry.translation} / ID ${selectedEntry.id}`
          : "-"}
      </p>
    </div>
    <div class="modal-actions">
      <button
        class="button-secondary"
        id="closeDeleteModalButton"
        onclick={closeDeleteModal}
        type="button">やめる</button
      >
      <button
        class="button-danger"
        id="confirmDeleteButton"
        onclick={deleteCurrentEntry}
        type="button">削除する</button
      >
    </div>
  </section>
</div>

<style>
  .modal-backdrop[hidden] {
    display: none !important;
    pointer-events: none;
  }

  .modal-backdrop {
    position: fixed;
    inset: 0;
    display: grid;
    place-items: center;
    padding: 18px;
    background: rgba(0, 0, 0, 0.5);
    z-index: 40;
  }

  .modal-card {
    width: min(560px, 100%);
    padding: 18px;
    border-radius: 14px;
    border: 1px solid var(--line);
    background: rgba(20, 16, 13, 0.96);
    color: var(--text);
    display: grid;
    gap: 12px;
  }

  .modal-actions,
  .delete-target {
    display: flex;
    flex-wrap: wrap;
    justify-content: space-between;
    align-items: center;
    gap: 10px;
  }

  p {
    color: var(--muted);
  }

  .button-secondary,
  .button-danger {
    min-height: 36px;
    padding: 0 14px;
    border-radius: 999px;
    border: 1px solid transparent;
    font: inherit;
  }

  .button-secondary {
    color: var(--text);
    background: rgba(255, 255, 255, 0.04);
    border-color: var(--line);
  }

  .button-danger {
    color: #3d1512;
    background: linear-gradient(135deg, #ffc0ab 0%, #ff9975 100%);
  }
</style>
