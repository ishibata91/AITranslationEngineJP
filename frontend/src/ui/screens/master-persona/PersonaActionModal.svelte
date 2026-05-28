<script lang="ts">
  import type { PersonaActionModalProps } from "./master-persona-panel-props"

  let {
    modalState,
    selectedEntry,
    editForm,
    errorMessage = "",
    closeEdit,
    closeDelete,
    saveCurrentEntry,
    deleteCurrentEntry,
    setEditFormField
  }: PersonaActionModalProps = $props()
</script>

<div
  aria-hidden={modalState !== "edit"}
  class="modal-backdrop"
  class:is-open={modalState === "edit"}
  hidden={modalState !== "edit"}
  id="editModal"
  role="dialog"
  data-testid="master-persona-edit-modal"
>
  <section class="modal-card form-modal">
    <div class="section-head">
      <div>
        <p class="eyebrow">編集</p>
        <h3>ペルソナを編集</h3>
      </div>
      <button
        class="button-secondary"
        data-testid="master-persona-edit-cancel-button"
        id="closeEditModalButton"
        onclick={closeEdit}
        type="button"
      >
        閉じる
      </button>
    </div>

    <div class="identity-banner">
      <strong>{selectedEntry?.displayName || "未選択"}</strong>
      <span
        >{selectedEntry?.targetPlugin || "-"} / {selectedEntry?.editorId ||
          "-"}</span
      >
    </div>

    {#if errorMessage}
      <p class="modal-error" role="alert">{errorMessage}</p>
    {/if}

    <div class="form-grid">
      <label class="field-group textarea-group" for="editPersonaSummaryInput">
        <span class="field-label">ペルソナ要約</span>
        <textarea
          class="textarea-field"
          data-testid="master-persona-summary-input"
          id="editPersonaSummaryInput"
          oninput={(event) => setEditFormField("personaSummary", event)}
          value={editForm.personaSummary ?? ""}
        ></textarea>
      </label>

      <label class="field-group" for="editSpeechStyleInput">
        <span class="field-label">話し方</span>
        <input
          class="text-field"
          data-testid="master-persona-speech-style-input"
          id="editSpeechStyleInput"
          oninput={(event) => setEditFormField("speechStyle", event)}
          value={editForm.speechStyle ?? ""}
        />
      </label>

      <label class="field-group textarea-group" for="editPersonaBodyInput">
        <span class="field-label">ペルソナ本文</span>
        <textarea
          class="textarea-field body-field"
          data-testid="master-persona-body-input"
          id="editPersonaBodyInput"
          oninput={(event) => setEditFormField("personaBody", event)}
          value={editForm.personaBody}
        ></textarea>
      </label>
    </div>

    <div class="button-row">
      <button
        class="button-secondary"
        data-testid="master-persona-edit-cancel-button"
        onclick={closeEdit}
        type="button"
      >
        キャンセル
      </button>
      <button
        class="button-primary"
        data-testid="master-persona-edit-save-button"
        id="saveEntryButton"
        onclick={saveCurrentEntry}
        type="button"
      >
        編集内容を保存
      </button>
    </div>
  </section>
</div>

<div
  aria-hidden={modalState !== "delete"}
  class="modal-backdrop"
  class:is-open={modalState === "delete"}
  hidden={modalState !== "delete"}
  id="deleteModal"
  role="dialog"
  data-testid="master-persona-delete-modal"
>
  <section class="modal-card">
    <div class="section-head">
      <div>
        <p class="eyebrow">削除</p>
        <h3>ペルソナを削除しますか</h3>
      </div>
      <button
        class="button-secondary"
        data-testid="master-persona-edit-cancel-button"
        onclick={closeDelete}
        type="button"
      >
        閉じる
      </button>
    </div>

    <p class="support-copy">
      削除すると、選択中のペルソナを一覧から外します。必要なら識別情報を確認してから実行してください。
    </p>

    {#if errorMessage}
      <p class="modal-error" role="alert">{errorMessage}</p>
    {/if}

    <dl class="detail-grid">
      <div class="detail-card">
        <dt>名前</dt>
        <dd>{selectedEntry?.displayName || "-"}</dd>
      </div>
      <div class="detail-card">
        <dt>FormID</dt>
        <dd>{selectedEntry?.formId || "-"}</dd>
      </div>
      <div class="detail-card">
        <dt>EditorID</dt>
        <dd>{selectedEntry?.editorId || "-"}</dd>
      </div>
      <div class="detail-card">
        <dt>対象プラグイン</dt>
        <dd>{selectedEntry?.targetPlugin || "-"}</dd>
      </div>
    </dl>

    <div class="button-row">
      <button
        class="button-secondary"
        data-testid="master-persona-edit-cancel-button"
        onclick={closeDelete}
        type="button"
      >
        キャンセル
      </button>
      <button
        class="button-danger"
        data-testid="master-persona-delete-confirm-button"
        id="confirmDeleteButton"
        onclick={deleteCurrentEntry}
        type="button"
      >
        削除する
      </button>
    </div>
  </section>
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
    z-index: 40;
  }

  .modal-backdrop.is-open {
    display: flex;
  }

  .modal-card,
  .identity-banner,
  .detail-card {
    border-radius: 20px;
  }

  .modal-card {
    background: rgba(17, 13, 12, 0.94);
    border: 0.5px solid var(--line);
    box-shadow: var(--shadow);
    color: var(--text);
    display: grid;
    gap: 16px;
    max-height: min(88vh, 900px);
    max-width: 760px;
    min-width: 0;
    overflow: auto;
    padding: clamp(18px, 3vw, 24px);
    width: min(100%, 760px);
  }

  .section-head,
  .button-row {
    align-items: flex-start;
    display: flex;
    flex-wrap: wrap;
    gap: 12px;
    justify-content: space-between;
  }

  .eyebrow,
  .field-label,
  .detail-card dt,
  .identity-banner span {
    color: var(--muted);
    font-size: 12px;
    letter-spacing: 0.08em;
    text-transform: uppercase;
  }

  h3,
  p,
  dl {
    margin: 0;
  }

  h3,
  strong,
  .support-copy,
  .detail-card dd,
  .identity-banner span {
    overflow-wrap: anywhere;
  }

  .support-copy,
  .modal-error,
  .detail-card dd,
  .identity-banner span {
    color: var(--muted);
    line-height: 1.7;
  }

  .modal-error {
    color: #ffd5cb;
  }

  .identity-banner {
    background: rgba(255, 255, 255, 0.03);
    border: 0.5px solid rgba(255, 255, 255, 0.08);
    display: grid;
    gap: 6px;
    padding: 14px;
  }

  .form-grid,
  .detail-grid {
    display: grid;
    gap: 12px;
  }

  .detail-grid {
    grid-template-columns: repeat(2, minmax(0, 1fr));
  }

  .field-group {
    display: grid;
    gap: 8px;
    min-width: 0;
  }

  .textarea-group {
    grid-column: 1 / -1;
  }

  .text-field,
  .textarea-field {
    background: rgba(255, 255, 255, 0.04);
    border: 0.5px solid var(--line);
    border-radius: 14px;
    color: var(--text);
    min-width: 0;
    padding: 12px 14px;
    width: 100%;
  }

  .text-field {
    min-height: 44px;
  }

  .textarea-field {
    min-height: 120px;
    resize: vertical;
  }

  .body-field {
    min-height: 220px;
  }

  .detail-card {
    background: rgba(255, 255, 255, 0.03);
    border: 0.5px solid rgba(255, 255, 255, 0.08);
    min-width: 0;
    padding: 14px;
  }

  .detail-card dt,
  .detail-card dd {
    margin: 0;
  }

  .detail-card dd {
    margin-top: 8px;
  }

  .button-primary,
  .button-secondary,
  .button-danger {
    border-radius: 999px;
    cursor: pointer;
    min-height: 40px;
    min-width: 0;
    overflow-wrap: anywhere;
    padding: 0 16px;
  }

  .button-primary {
    background: linear-gradient(135deg, var(--primary) 0%, #f0a51f 100%);
    border: 0.5px solid transparent;
    color: #3f2400;
  }

  .button-secondary {
    background: rgba(255, 255, 255, 0.04);
    border: 0.5px solid var(--line);
    color: var(--text);
  }

  .button-danger {
    background: linear-gradient(135deg, #ffc0ab 0%, #ff9c7c 100%);
    border: 0.5px solid transparent;
    color: #35150d;
  }

  @media (max-width: 560px) {
    .detail-grid {
      grid-template-columns: 1fr;
    }

    .button-row > * {
      width: 100%;
    }
  }
</style>
