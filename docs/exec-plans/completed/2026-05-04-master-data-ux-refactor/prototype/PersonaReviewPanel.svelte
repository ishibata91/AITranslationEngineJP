<script lang="ts">
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
    personas: Persona[]
    selectedIdentityKey: string
    pageInfo?: {
      currentPage: number
      pageSize: number
      totalCount: number
    }
    openEditAction: () => void
    openDeleteAction: () => void
    changePageAction: (page: number) => void
    selectPersonaAction: (identityKey: string) => void
  }

  let {
    personas,
    selectedIdentityKey,
    pageInfo,
    openEditAction,
    openDeleteAction,
    changePageAction,
    selectPersonaAction
  }: Props = $props()

  const selected = $derived(
    personas.find((persona) => persona.identityKey === selectedIdentityKey) ?? personas[0]
  )
  const resolvedPageInfo = $derived(
    pageInfo ?? {
      currentPage: 1,
      pageSize: 50,
      totalCount: personas.length
    }
  )
  const totalPages = $derived(
    Math.max(1, Math.ceil(resolvedPageInfo.totalCount / resolvedPageInfo.pageSize))
  )
  const pageStart = $derived((resolvedPageInfo.currentPage - 1) * resolvedPageInfo.pageSize + 1)
  const pageEnd = $derived(
    Math.min(resolvedPageInfo.currentPage * resolvedPageInfo.pageSize, resolvedPageInfo.totalCount)
  )
</script>

<section class="review-grid" aria-label="生成結果の確認">
  <section class="panel list-panel" aria-labelledby="personaListHeading">
    <div class="section-head">
      <div>
        <p class="eyebrow">生成結果</p>
        <h2 id="personaListHeading">ペルソナ一覧</h2>
      </div>
      <span class="status-pill">{pageStart}-{pageEnd} / {resolvedPageInfo.totalCount} 件</span>
    </div>

    <div class="filter-row" aria-label="検索と絞り込み">
      <label>
        <span>検索</span>
        <input placeholder="名前、プラグインで検索" type="search" />
      </label>
      <label>
        <span>プラグイン</span>
        <select>
          <option>すべて</option>
          <option>Skyrim.esm</option>
          <option>Update.esm</option>
        </select>
      </label>
    </div>

    <div class="list-stack">
      {#each personas as persona}
        <button
          aria-pressed={persona.identityKey === selectedIdentityKey}
          class:is-selected={persona.identityKey === selectedIdentityKey}
          class="persona-row"
          onclick={() => selectPersonaAction(persona.identityKey)}
          type="button"
        >
          <span>{persona.targetPlugin}</span>
          <strong>{persona.displayName}</strong>
        </button>
      {/each}
    </div>

    <nav class="pagination-row" aria-label="ペルソナ一覧のページ操作">
      <button
        class="button-secondary"
        disabled={resolvedPageInfo.currentPage <= 1}
        onclick={() => changePageAction(resolvedPageInfo.currentPage - 1)}
        type="button"
      >
        前へ
      </button>
      <span>{resolvedPageInfo.currentPage} / {totalPages} ページ</span>
      <button
        class="button-secondary"
        disabled={resolvedPageInfo.currentPage >= totalPages}
        onclick={() => changePageAction(resolvedPageInfo.currentPage + 1)}
        type="button"
      >
        次へ
      </button>
    </nav>
  </section>

  <section class="panel detail-panel" aria-labelledby="personaDetailHeading">
    <div class="section-head">
      <div>
        <p class="eyebrow">詳細</p>
        <h2 id="personaDetailHeading">{selected.displayName}</h2>
      </div>
      <div class="button-row">
        <button
          class="button-secondary"
          onclick={openEditAction}
          type="button"
        >
          編集
        </button>
        <button
          class="button-secondary danger"
          onclick={openDeleteAction}
          type="button"
        >
          削除
        </button>
      </div>
    </div>

    <dl class="identity-list">
      <div>
        <dt>識別情報</dt>
        <dd>{selected.formId} / {selected.editorId}</dd>
      </div>
      <div>
        <dt>声</dt>
        <dd>{selected.voiceType}</dd>
      </div>
      <div>
        <dt>話し方</dt>
        <dd>{selected.speechStyle || "未入力"}</dd>
      </div>
    </dl>

    <article class="body-card">
      <span>ペルソナ本文</span>
      <p>{selected.personaBody || "生成後に一覧と詳細で確認できます。"}</p>
    </article>
  </section>
</section>

<style>
  .review-grid {
    align-items: start;
    display: grid;
    gap: 14px;
    grid-template-columns: minmax(0, 1fr) minmax(320px, 0.82fr);
  }

  .panel {
    background: var(--surface);
    border: 1px solid var(--line);
    border-radius: 8px;
    box-shadow: var(--shadow);
    min-width: 0;
    padding: clamp(16px, 2vw, 22px);
  }

  .section-head,
  .button-row,
  .filter-row,
  .pagination-row,
  .persona-row {
    align-items: flex-start;
    display: flex;
    flex-wrap: wrap;
    gap: 12px;
    justify-content: space-between;
  }

  .filter-row {
    display: grid;
    grid-template-columns: minmax(0, 1fr) minmax(180px, 0.38fr);
    margin-top: 14px;
  }

  label,
  .identity-list div,
  .body-card {
    display: grid;
    gap: 6px;
  }

  label,
  .identity-list div,
  .body-card {
    background: rgba(0, 0, 0, 0.14);
    border: 1px solid rgba(255, 255, 255, 0.08);
    border-radius: 8px;
    padding: 12px;
  }

  input,
  select {
    background: rgba(255, 255, 255, 0.08);
    border: 1px solid var(--line);
    border-radius: 6px;
    color: var(--text);
    min-height: 40px;
    min-width: 0;
    padding: 8px 10px;
    width: 100%;
  }

  .list-stack,
  .identity-list {
    display: grid;
    gap: 8px;
    margin: 14px 0 0;
  }

  .persona-row {
    align-items: center;
    background: rgba(255, 255, 255, 0.04);
    border: 1px solid var(--line);
    border-radius: 8px;
    color: var(--text);
    cursor: pointer;
    display: grid;
    gap: 10px;
    grid-template-columns: minmax(96px, 0.36fr) minmax(0, 1fr);
    min-height: 44px;
    padding: 8px 10px;
    text-align: left;
    width: 100%;
  }

  .persona-row.is-selected {
    background: rgba(255, 186, 56, 0.13);
    border-color: var(--line-strong);
  }

  .persona-row span {
    color: var(--muted);
    font-size: 0.82rem;
    overflow: hidden;
    text-overflow: ellipsis;
    white-space: nowrap;
  }

  .persona-row strong {
    overflow: hidden;
    text-overflow: ellipsis;
    white-space: nowrap;
  }

  .eyebrow,
  label span,
  dt,
  .body-card span {
    color: var(--primary);
    font-size: 0.78rem;
  }

  dd,
  .body-card p {
    color: var(--muted);
    line-height: 1.6;
  }

  dt,
  dd,
  p {
    margin: 0;
  }

  .status-pill {
    border: 1px solid var(--line-strong);
    border-radius: 999px;
    color: var(--accent);
    flex: none;
    padding: 7px 10px;
    white-space: nowrap;
  }

  .button-secondary {
    background: rgba(255, 255, 255, 0.06);
    border: 1px solid var(--line);
    border-radius: 7px;
    color: var(--text);
    cursor: pointer;
    min-height: 40px;
    padding: 9px 13px;
  }

  .button-secondary.danger {
    border-color: rgba(255, 140, 120, 0.5);
  }

  .pagination-row {
    align-items: center;
    color: var(--muted);
    margin-top: 12px;
  }

  h2 {
    margin: 0;
  }

  @media (max-width: 980px) {
    .review-grid,
    .filter-row {
      grid-template-columns: 1fr;
    }
  }
</style>
