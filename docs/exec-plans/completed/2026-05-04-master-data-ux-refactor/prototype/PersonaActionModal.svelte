<script lang="ts">
  type PersonaModalMode = "edit" | "delete" | null

  interface Persona {
    identityKey: string
    displayName: string
    targetPlugin: string
    formId: string
    editorId: string
    voiceType: string
    personaSummary: string
    personaBody: string
    speechStyle: string
  }

  interface Props {
    mode: PersonaModalMode
    persona: Persona | null
    closeAction: () => void
    saveAction: () => void
    deleteAction: () => void
  }

  let { mode, persona, closeAction, saveAction, deleteAction }: Props = $props()
</script>

{#if mode && persona}
  <div class="modal-backdrop" role="presentation">
    <div
      aria-labelledby="personaActionModalHeading"
      aria-modal="true"
      class="modal-card"
      role="dialog"
    >
      <div class="section-head">
        <div>
          <p class="eyebrow">{mode === "edit" ? "編集" : "削除"}</p>
          <h2 id="personaActionModalHeading">
            {mode === "edit" ? "ペルソナを編集" : "ペルソナを削除"}
          </h2>
        </div>
        <button
          class="button-secondary"
          onclick={closeAction}
          type="button"
        >
          閉じる
        </button>
      </div>

      {#if mode === "edit"}
        <div class="form-grid">
          <label class="field-block" for="personaSummaryInput">
            <span>ペルソナ要約</span>
            <textarea id="personaSummaryInput" value={persona.personaSummary}></textarea>
          </label>
          <label class="field-block" for="speechStyleInput">
            <span>話し方</span>
            <input id="speechStyleInput" value={persona.speechStyle} />
          </label>
          <label class="field-block wide" for="personaBodyInput">
            <span>ペルソナ本文</span>
            <textarea id="personaBodyInput" value={persona.personaBody}></textarea>
          </label>
        </div>

        <div class="button-row">
          <button
            class="button-secondary"
            onclick={closeAction}
            type="button"
          >
            キャンセル
          </button>
          <button
            class="button-primary"
            onclick={saveAction}
            type="button"
          >
            保存
          </button>
        </div>
      {:else}
        <p class="support-copy">
          このペルソナを一覧から削除します。削除後は元に戻せません。
        </p>

        <dl class="delete-summary">
          <div>
            <dt>名前</dt>
            <dd>{persona.displayName}</dd>
          </div>
          <div>
            <dt>識別情報</dt>
            <dd>{persona.formId} / {persona.editorId}</dd>
          </div>
          <div>
            <dt>プラグイン</dt>
            <dd>{persona.targetPlugin}</dd>
          </div>
        </dl>

        <div class="button-row">
          <button
            class="button-secondary"
            onclick={closeAction}
            type="button"
          >
            キャンセル
          </button>
          <button
            class="button-danger"
            onclick={deleteAction}
            type="button"
          >
            削除
          </button>
        </div>
      {/if}
    </div>
  </div>
{/if}

<style>
  .modal-backdrop {
    align-items: center;
    background: rgba(8, 8, 8, 0.68);
    display: flex;
    inset: 0;
    justify-content: center;
    padding: 18px;
    position: fixed;
    z-index: 20;
  }

  .modal-card {
    background: rgba(26, 25, 24, 0.98);
    border: 1px solid var(--line);
    border-radius: 8px;
    box-shadow: var(--shadow);
    color: var(--text);
    display: grid;
    font-size: 1rem;
    gap: 16px;
    max-height: calc(100vh - 36px);
    max-width: 840px;
    overflow: auto;
    padding: clamp(16px, 3vw, 24px);
    width: min(100%, 840px);
  }

  .section-head,
  .button-row {
    align-items: flex-start;
    display: flex;
    flex-wrap: wrap;
    gap: 12px;
    justify-content: space-between;
  }

  .form-grid {
    display: grid;
    gap: 12px;
    grid-template-columns: minmax(0, 1fr) minmax(220px, 0.48fr);
  }

  .field-block,
  .delete-summary div {
    background: rgba(0, 0, 0, 0.16);
    border: 1px solid rgba(255, 255, 255, 0.08);
    border-radius: 8px;
    display: grid;
    gap: 8px;
    padding: 12px;
  }

  .field-block.wide {
    grid-column: 1 / -1;
  }

  input,
  textarea {
    background: rgba(255, 255, 255, 0.08);
    border: 1px solid var(--line);
    border-radius: 6px;
    color: var(--text);
    font: inherit;
    line-height: 1.6;
    min-height: 40px;
    min-width: 0;
    padding: 8px 10px;
    width: 100%;
  }

  textarea {
    min-height: 112px;
    resize: vertical;
  }

  .delete-summary {
    display: grid;
    gap: 8px;
    margin: 0;
  }

  .eyebrow,
  .field-block span,
  dt {
    color: var(--primary);
    font-size: 0.9rem;
    line-height: 1.4;
  }

  .support-copy,
  dd {
    color: var(--muted);
    line-height: 1.6;
  }

  dd,
  dt,
  h2,
  p {
    margin: 0;
  }

  h2 {
    font-size: 1.5rem;
    line-height: 1.2;
  }

  .button-primary,
  .button-secondary,
  .button-danger {
    border-radius: 7px;
    cursor: pointer;
    font: inherit;
    min-height: 40px;
    padding: 9px 13px;
  }

  .button-primary {
    background: var(--primary);
    border: 1px solid var(--primary);
    color: #201309;
  }

  .button-secondary {
    background: rgba(255, 255, 255, 0.06);
    border: 1px solid var(--line);
    color: var(--text);
  }

  .button-danger {
    background: rgba(255, 140, 120, 0.18);
    border: 1px solid rgba(255, 140, 120, 0.58);
    color: #ffd4cb;
  }

  @media (max-width: 720px) {
    .form-grid {
      grid-template-columns: 1fr;
    }
  }
</style>
