<script lang="ts">
  import type { DictionaryEditModalProps } from "./dictionary-panel-props"

  let {
    categoryOptions,
    formCategory,
    formOrigin,
    formSource,
    formTranslation,
    modalState,
    closeEditModal,
    saveCurrentEntry,
    setFormCategory,
    setFormOrigin,
    setFormSource,
    setFormTranslation
  }: DictionaryEditModalProps = $props()

  const isOpen = $derived(modalState === "create" || modalState === "edit")
</script>

<div
  aria-hidden={!isOpen}
  class="modal-backdrop"
  hidden={!isOpen}
  id="editModal"
  role="dialog"
  data-testid="master-dictionary-create-edit-modal"
>
  <section aria-labelledby="editModalTitle" class="modal-card">
    <div class="eyebrow" id="editModalEyebrow">
      {modalState === "create" ? "新規登録" : "更新"}
    </div>
    <h3 id="editModalTitle">
      {modalState === "create" ? "新規登録" : "更新"}
    </h3>
    <p id="editModalDescription">
      {modalState === "create"
        ? "辞書エントリの内容を入力します。"
        : "選択中の辞書エントリを編集します。"}
    </p>
    <p data-testid="master-dictionary-entry-validation-error" hidden></p>
    <div class="field-grid">
      <label class="field-label" for="formSource">原文</label>
      <input
        class="text-field"
        data-testid="master-dictionary-entry-source-input"
        id="formSource"
        type="text"
        value={formSource}
        oninput={setFormSource}
      />

      <label class="field-label" for="formCategory">カテゴリ</label>
      <select
        class="select-field"
        data-testid="master-dictionary-entry-category-select"
        id="formCategory"
        value={formCategory}
        onchange={setFormCategory}
      >
        {#each categoryOptions.filter((item) => item !== "すべて") as option (option)}
          <option value={option}>{option}</option>
        {/each}
      </select>

      <label class="field-label" for="formOrigin">由来</label>
      <select
        class="select-field"
        data-testid="master-dictionary-entry-origin-input"
        id="formOrigin"
        value={formOrigin}
        onchange={setFormOrigin}
      >
        <option value="手動登録">手動登録</option>
        <option value="確認待ち">確認待ち</option>
        <option value="XML取込">XML取込</option>
      </select>

      <label class="field-label" for="formTranslation">訳語</label>
      <textarea
        class="textarea-field"
        data-testid="master-dictionary-entry-translation-input"
        id="formTranslation"
        value={formTranslation}
        oninput={setFormTranslation}
      ></textarea>
    </div>
    <div class="modal-actions">
      <button
        class="button-secondary"
        id="closeEditModalButton"
        onclick={closeEditModal}
        type="button">閉じる</button
      >
      <button
        class="button-primary"
        data-testid="master-dictionary-entry-save-button"
        id="saveEntryButton"
        onclick={saveCurrentEntry}
        type="button">保存する</button
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

  .modal-actions {
    display: flex;
    flex-wrap: wrap;
    justify-content: space-between;
    align-items: center;
    gap: 10px;
  }

  .field-grid {
    display: grid;
    gap: 10px;
  }

  .eyebrow,
  .field-label {
    font-size: 12px;
    letter-spacing: 0.08em;
  }

  p {
    color: var(--muted);
  }

  .button-primary,
  .button-secondary {
    min-height: 36px;
    padding: 0 14px;
    border-radius: 999px;
    border: 1px solid transparent;
    font: inherit;
  }

  .button-primary {
    color: #3a2400;
    background: linear-gradient(135deg, var(--primary) 0%, #ef9d20 100%);
  }

  .button-secondary {
    color: var(--text);
    background: rgba(255, 255, 255, 0.04);
    border-color: var(--line);
  }

  .text-field,
  .select-field,
  .textarea-field {
    width: 100%;
    min-height: 38px;
    border-radius: 10px;
    border: 1px solid var(--line);
    background: rgba(0, 0, 0, 0.24);
    color: var(--text);
    padding: 0 10px;
  }

  .textarea-field {
    min-height: 90px;
    padding: 10px;
  }
</style>
