<script lang="ts">
  import type { MasterPersonaScreenViewModel } from "@application/gateway-contract/master-persona"
  import type { MasterPersonaListItem } from "@application/gateway-contract/master-persona/master-persona-gateway-contract"

  interface Props {
    viewModel: MasterPersonaScreenViewModel
    selectRow: (identityKey: string) => void
    updateKeyword: (event: Event) => void
    updatePluginFilter: (event: Event) => void
    goToPrevPage: () => void
    goToNextPage: () => void
    editCurrent: () => void
    openDelete: () => void
  }

  let {
    viewModel,
    selectRow,
    updateKeyword,
    updatePluginFilter,
    goToPrevPage,
    goToNextPage,
    editCurrent,
    openDelete
  }: Props = $props()

  function itemLabel(item: MasterPersonaListItem): string {
    return item.displayName || item.editorId || item.formId
  }

  const pageRangeText = $derived(
    viewModel.totalCount === 0
      ? "0 件"
      : `${(viewModel.page - 1) * viewModel.pageSize + 1}-${Math.min(viewModel.page * viewModel.pageSize, viewModel.totalCount)} / ${viewModel.totalCount.toLocaleString("ja-JP")} 件`
  )
</script>

<section class="review-grid" aria-label="生成結果の確認">
  <section class="panel list-panel" aria-labelledby="listHeading">
    <div class="section-head">
      <div>
        <p class="eyebrow">生成結果</p>
        <h3 id="listHeading">ペルソナ一覧</h3>
      </div>
      <span class="status-pill">{pageRangeText}</span>
    </div>

    <div class="filter-grid">
      <label class="field-group" for="masterPersonaSearchInput">
        <span class="field-label">検索</span>
        <input
          class="text-field"
          id="masterPersonaSearchInput"
          oninput={updateKeyword}
          placeholder="名前またはプラグイン名で検索"
          type="search"
          value={viewModel.keyword}
        />
      </label>

      <label class="field-group" for="masterPersonaPluginSelect">
        <span class="field-label">プラグイン</span>
        <select
          class="text-field"
          id="masterPersonaPluginSelect"
          onchange={updatePluginFilter}
          value={viewModel.pluginFilter}
        >
          {#each viewModel.pluginOptions as option (option.label)}
            <option value={option.value}>{option.label}</option>
          {/each}
        </select>
      </label>
    </div>

    <div class="list-stack" aria-live="polite">
      {#if viewModel.items.length === 0}
        <div class="empty-state">
          生成済みまたは条件に合うペルソナがありません。JSON を選択して作成後に確認できます。
        </div>
      {:else}
        {#each viewModel.items as item (item.identityKey)}
          <button
            aria-pressed={viewModel.selectedIdentityKey === item.identityKey}
            class:is-selected={viewModel.selectedIdentityKey === item.identityKey}
            class="list-row"
            onclick={() => selectRow(item.identityKey)}
            type="button"
          >
            <span class="plugin-name">{item.targetPlugin}</span>
            <strong>{itemLabel(item)}</strong>
          </button>
        {/each}
      {/if}
    </div>

    <nav class="pager-row" aria-label="ペルソナ一覧のページ操作">
      <span class="support-copy">{viewModel.page} / {viewModel.totalPages} ページ</span>
      <div class="pager-actions">
        <button
          class="button-secondary"
          disabled={viewModel.page <= 1}
          id="prevPageButton"
          onclick={goToPrevPage}
          type="button"
        >
          前へ
        </button>
        <button
          class="button-secondary"
          disabled={viewModel.page >= viewModel.totalPages}
          id="nextPageButton"
          onclick={goToNextPage}
          type="button"
        >
          次へ
        </button>
      </div>
    </nav>
  </section>

  <section class="panel detail-panel" aria-labelledby="detailHeading">
    <div class="section-head">
      <div>
        <p class="eyebrow">詳細</p>
        <h3 id="detailHeading">{viewModel.selectedEntry?.displayName || "選択中のペルソナ"}</h3>
      </div>
      <div class="detail-actions">
        <button
          class="button-secondary"
          disabled={!viewModel.canMutate}
          id="editButton"
          onclick={editCurrent}
          type="button"
        >
          編集
        </button>
        <button
          class="button-danger"
          disabled={!viewModel.canMutate}
          id="deleteButton"
          onclick={openDelete}
          type="button"
        >
          削除
        </button>
      </div>
    </div>

    <p class="identity-text" id="detailIdentityText">
      {#if viewModel.selectedEntry}
        FormID {viewModel.selectedEntry.formId} / EditorID {viewModel.selectedEntry.editorId}
      {:else}
        一覧からペルソナを選んでください。
      {/if}
    </p>

    <dl class="identity-grid">
      <div class="detail-card">
        <dt>FormID</dt>
        <dd>{viewModel.selectedEntry?.formId || "-"}</dd>
      </div>
      <div class="detail-card">
        <dt>EditorID</dt>
        <dd>{viewModel.selectedEntry?.editorId || "-"}</dd>
      </div>
      <div class="detail-card">
        <dt>対象プラグイン</dt>
        <dd>{viewModel.selectedEntry?.targetPlugin || "-"}</dd>
      </div>
      <div class="detail-card">
        <dt>元プラグイン</dt>
        <dd>{viewModel.selectedEntry?.sourcePlugin || "-"}</dd>
      </div>
      <div class="detail-card">
        <dt>声</dt>
        <dd>{viewModel.selectedEntry?.voiceType || "-"}</dd>
      </div>
      <div class="detail-card">
        <dt>話し方</dt>
        <dd>{viewModel.selectedEntry?.speechStyle || "未入力"}</dd>
      </div>
    </dl>

    <article class="body-card">
      <span class="field-label">ペルソナ本文</span>
      <p>{viewModel.selectedEntry?.personaBody || "生成後に一覧から選ぶと、本文を確認できます。"}</p>
    </article>
  </section>
</section>

<style>
  .review-grid {
    align-items: start;
    display: grid;
    gap: 12px;
    grid-template-columns: minmax(0, 0.92fr) minmax(0, 1.08fr);
  }

  .panel,
  .detail-card,
  .body-card,
  .list-row,
  .empty-state {
    border-radius: 20px;
  }

  .panel {
    align-content: start;
    background: rgba(17, 13, 12, 0.42);
    border: 0.5px solid var(--line);
    box-shadow: var(--shadow);
    display: grid;
    gap: 6px;
    min-width: 0;
    padding: 12px;
  }

  .section-head,
  .detail-actions,
  .pager-row,
  .pager-actions {
    align-items: flex-start;
    display: flex;
    flex-wrap: wrap;
    gap: 6px;
    justify-content: space-between;
  }

  .eyebrow,
  .field-label,
  .detail-card dt,
  .plugin-name {
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
  .support-copy,
  .list-row strong,
  .detail-card dd,
  .body-card p,
  .identity-text,
  .empty-state {
    overflow-wrap: anywhere;
  }

  .support-copy,
  .identity-text,
  .body-card p,
  .detail-card dd,
  .empty-state {
    color: var(--muted);
    line-height: 1.7;
  }

  .status-pill {
    align-items: center;
    background: rgba(255, 255, 255, 0.04);
    border: 0.5px solid rgba(255, 186, 56, 0.22);
    border-radius: 999px;
    color: var(--text);
    display: inline-flex;
    min-height: 28px;
    padding: 0 10px;
  }

  .filter-grid {
    align-items: end;
    display: grid;
    gap: 4px;
    grid-template-columns: minmax(0, 1fr) minmax(200px, 0.42fr);
  }

  .field-group,
  .identity-grid {
    display: grid;
    gap: 4px;
    min-width: 0;
  }

  .field-label {
    margin: 0;
  }

  .text-field {
    appearance: none;
    background: rgba(255, 255, 255, 0.04);
    border: 0.5px solid var(--line);
    border-radius: 10px;
    color: var(--text);
    height: 34px;
    min-width: 0;
    padding: 0 10px;
    width: 100%;
  }

  .list-stack {
    display: grid;
    gap: 2px;
  }

  .list-row,
  .empty-state {
    background: rgba(255, 255, 255, 0.03);
    border: 0.5px solid rgba(255, 255, 255, 0.08);
    min-width: 0;
    padding: 3px 8px;
    text-align: left;
  }

  .list-row {
    align-items: center;
    cursor: pointer;
    display: grid;
    gap: 3px;
    grid-template-columns: minmax(96px, 0.4fr) minmax(0, 1fr);
    min-height: 28px;
  }

  .list-row.is-selected {
    background: rgba(255, 186, 56, 0.12);
    border-color: rgba(255, 186, 56, 0.28);
  }

  .plugin-name,
  .list-row strong {
    min-width: 0;
  }

  .list-row strong {
    color: var(--text);
  }

  .plugin-name,
  .list-row strong {
    overflow: hidden;
    text-overflow: ellipsis;
    white-space: nowrap;
  }

  .identity-grid {
    grid-template-columns: repeat(2, minmax(0, 1fr));
    margin: 0;
  }

  .detail-card,
  .body-card {
    background: rgba(255, 255, 255, 0.03);
    border: 0.5px solid rgba(255, 255, 255, 0.08);
    min-width: 0;
    padding: 12px;
  }

  .detail-card dt,
  .detail-card dd {
    margin: 0;
  }

  .detail-card dd {
    margin-top: 8px;
  }

  .button-secondary,
  .button-danger {
    border-radius: 999px;
    cursor: pointer;
    min-height: 40px;
    min-width: 0;
    overflow-wrap: anywhere;
    padding: 0 16px;
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

  button:disabled {
    cursor: not-allowed;
    opacity: 0.5;
  }

  @media (max-width: 1080px) {
    .review-grid {
      grid-template-columns: 1fr;
    }
  }

  @media (max-width: 560px) {
    .filter-grid,
    .identity-grid,
    .list-row {
      grid-template-columns: 1fr;
    }

    .plugin-name,
    .list-row strong {
      overflow: visible;
      text-overflow: clip;
      white-space: normal;
    }

    .pager-row,
    .pager-actions,
    .detail-actions {
      width: 100%;
    }

    .pager-actions > *,
    .detail-actions > * {
      width: 100%;
    }
  }
</style>
