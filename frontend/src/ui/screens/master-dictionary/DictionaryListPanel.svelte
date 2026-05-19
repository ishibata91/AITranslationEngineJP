<script lang="ts">
  import type { DictionaryListPanelProps } from "./dictionary-panel-props"

  let {
    category,
    categoryOptions,
    entries,
    listHeadline,
    page,
    pageStatusText,
    query,
    selectedId,
    selectionStatusText,
    totalPages,
    goToNextPage,
    goToPrevPage,
    handleCategoryChange,
    handleSearchInput,
    openCreateModal,
    selectRow
  }: DictionaryListPanelProps = $props()
</script>

<section
  class="shell-card list-panel"
  aria-labelledby="listHeading"
  data-testid="master-dictionary-dictionary-list-panel"
>
  <div class="toolbar" data-testid="master-dictionary-list-action-bar">
    <div class="toolbar-head">
      <div>
        <h3 id="listHeading">辞書一覧</h3>
        <p id="listHeadline">{listHeadline}</p>
      </div>
      <div class="toolbar-head-actions">
        <button
          class="button-primary"
          id="createButton"
          onclick={openCreateModal}
          type="button">新規登録</button
        >
        <p class="mini-text" id="pageStatusText">{pageStatusText}</p>
      </div>
    </div>

    <div class="toolbar-grid">
      <label class="field-group" for="searchInput">
        <span class="field-label">検索</span>
        <input
          class="search-field"
          id="searchInput"
          oninput={handleSearchInput}
          placeholder="原文・訳語・IDで検索"
          type="search"
          value={query}
        />
      </label>

      <label class="field-group" for="categorySelect">
        <span class="field-label">カテゴリ</span>
        <select
          class="select-field"
          id="categorySelect"
          onchange={handleCategoryChange}
          value={category}
        >
          {#each categoryOptions as option (option)}
            <option value={option}>{option}</option>
          {/each}
        </select>
      </label>
    </div>
  </div>

  <div class="list-shell" data-testid="master-dictionary-dictionary-list">
    <div class="column-row" aria-hidden="true">
      <span>訳語</span>
      <span>原文</span>
      <span>カテゴリ</span>
      <span>ID</span>
    </div>

    <div class="list-stack" id="listStack" aria-live="polite">
      {#if entries.length === 0}
        <div class="empty-state">一致するエントリがありません</div>
      {:else}
        {#each entries as entry (entry.id)}
          <button
            class="list-row"
            class:is-selected={selectedId === entry.id}
            onclick={() => selectRow(entry.id)}
            type="button"
          >
            <div class="row-cell">
              <div class="row-value">{entry.translation}</div>
            </div>
            <div class="row-cell">
              <div class="row-value">{entry.source}</div>
            </div>
            <div class="row-meta">{entry.category} / {entry.origin}</div>
            <div class="row-id">#{entry.id}</div>
          </button>
        {/each}
      {/if}
    </div>

    <div class="pager-shell" data-testid="master-dictionary-pagination-region">
      <div class="mini-text" id="selectionStatus">{selectionStatusText}</div>
      <div class="pager-actions">
        <button
          class="button-secondary"
          disabled={page === 0}
          id="prevPageButton"
          onclick={goToPrevPage}
          type="button"
        >
          前の30件
        </button>
        <button
          class="button-secondary"
          disabled={page + 1 >= totalPages}
          id="nextPageButton"
          onclick={goToNextPage}
          type="button"
        >
          次の30件
        </button>
      </div>
    </div>
  </div>
</section>

<style>
  .shell-card {
    padding: 18px;
    border-radius: 16px;
    border: 1px solid var(--line);
    background: rgba(16, 13, 11, 0.58);
    color: var(--text);
  }

  .list-panel,
  .toolbar,
  .toolbar-grid,
  .field-group,
  .list-shell {
    display: grid;
    gap: 10px;
  }

  .list-panel {
    min-width: 0;
    padding: 20px;
    gap: 14px;
  }

  .toolbar {
    width: 100%;
    min-width: 0;
    box-sizing: border-box;
    position: sticky;
    top: 18px;
    z-index: 1;
    gap: 14px;
    padding: 16px;
    border: 1px solid var(--line);
    border-radius: 18px;
    background: rgba(18, 16, 13, 0.72);
  }

  .toolbar-head,
  .toolbar-head-actions,
  .pager-shell,
  .pager-actions {
    display: flex;
    flex-wrap: wrap;
    justify-content: space-between;
    align-items: center;
    gap: 10px;
  }

  .toolbar-grid {
    width: 100%;
    min-width: 0;
    box-sizing: border-box;
    grid-template-columns: minmax(0, 1.5fr) minmax(220px, 0.7fr);
    gap: 12px;
  }

  .field-group {
    min-width: 0;
    gap: 8px;
  }

  .field-label {
    font-size: 12px;
    letter-spacing: 0.08em;
  }

  .mini-text,
  #listHeadline {
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

  .search-field,
  .select-field {
    width: 100%;
    min-height: 38px;
    border-radius: 10px;
    border: 1px solid var(--line);
    background: rgba(0, 0, 0, 0.24);
    color: var(--text);
    padding: 0 10px;
  }

  .list-shell {
    width: 100%;
    min-width: 0;
    box-sizing: border-box;
    gap: 8px;
    padding: 16px;
    border: 1px solid var(--line);
    border-radius: 18px;
    background: rgba(18, 16, 13, 0.72);
    overflow: hidden;
  }

  .column-row,
  .list-row {
    width: 100%;
    min-width: 0;
    box-sizing: border-box;
    display: grid;
    grid-template-columns:
      minmax(0, 1.4fr) minmax(0, 1.4fr) minmax(0, 0.8fr)
      90px;
    gap: 10px;
    align-items: center;
  }

  .column-row {
    padding: 0 12px;
  }

  .list-stack {
    min-height: 200px;
    padding: 0;
    gap: 6px;
  }

  .list-row {
    width: 100%;
    border: 1px solid rgba(255, 186, 56, 0.12);
    border-radius: 10px;
    background: rgba(255, 255, 255, 0.03);
    color: var(--text);
    padding: 8px 12px;
    text-align: left;
    cursor: pointer;
  }

  .list-row.is-selected {
    border-color: var(--line-strong);
    background: rgba(255, 186, 56, 0.12);
  }

  .row-cell {
    min-width: 0;
    display: grid;
    gap: 2px;
  }

  .row-value {
    overflow: hidden;
    white-space: nowrap;
    text-overflow: ellipsis;
  }

  .row-meta,
  .row-id {
    font-size: 12px;
    color: var(--muted);
  }

  .row-id {
    text-align: right;
  }

  .empty-state {
    color: var(--muted);
    padding: 8px;
  }

  @media (max-width: 1180px) {
    .toolbar-grid {
      grid-template-columns: 1fr;
    }
  }

  @media (max-width: 980px) {
    .toolbar {
      position: static;
    }

    .column-row {
      display: none;
    }

    .list-row {
      grid-template-columns: 1fr;
      gap: 5px;
    }

    .row-id {
      text-align: left;
    }
  }
</style>
