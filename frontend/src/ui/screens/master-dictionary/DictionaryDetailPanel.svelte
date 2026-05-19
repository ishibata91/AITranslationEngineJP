<script lang="ts">
  import type { DictionaryDetailPanelProps } from "./dictionary-panel-props"

  let {
    detailSublineText,
    selectedEntry,
    openDeleteModal,
    openEditModal
  }: DictionaryDetailPanelProps = $props()
</script>

<section
  class="shell-card detail-panel"
  aria-labelledby="detailHeading"
  data-testid="master-dictionary-detail-panel"
>
  <div class="detail-head">
    <div>
      <h3 id="detailHeading">詳細</h3>
      <p id="detailSubline">{detailSublineText}</p>
    </div>
    <div class="detail-actions">
      <button
        class="button-secondary"
        disabled={!selectedEntry}
        id="editButton"
        onclick={openEditModal}
        type="button"
      >
        更新
      </button>
      <button
        class="button-danger"
        disabled={!selectedEntry}
        id="deleteButton"
        onclick={openDeleteModal}
        type="button"
      >
        削除
      </button>
    </div>
  </div>

  <div class="detail-title">
    <div class="detail-tags" id="detailTags">
      {#if selectedEntry}
        <span class="status-pill">{selectedEntry.category}</span>
        <span class="status-pill">{selectedEntry.origin}</span>
      {/if}
    </div>
    <strong id="detailTitle"
      >{selectedEntry?.source ?? "表示できるエントリがありません"}</strong
    >
    <p id="detailTranslation">
      {selectedEntry?.translation ?? "検索条件を変更してください。"}
    </p>
  </div>

  <div class="detail-grid" id="detailGrid">
    {#if selectedEntry}
      <div class="detail-card detail-meta-card">
        <div class="field-label">ID</div>
        <strong class="detail-meta-value">{selectedEntry.id}</strong>
      </div>
      <div class="detail-card detail-meta-card">
        <div class="field-label">最終更新</div>
        <strong class="detail-meta-value">{selectedEntry.updatedAt}</strong>
      </div>
    {:else}
      <div class="empty-state">
        一覧に表示できるエントリが戻ると、詳細も同じ画面で切り替わります。
      </div>
    {/if}
  </div>
  <p id="detailStatusMessage">
    {selectedEntry
      ? `${selectedEntry.origin} / 最終更新 ${selectedEntry.updatedAt}`
      : "一覧から別のエントリを選択すると、ここも切り替わります。"}
  </p>
  <dl class="detail-list" id="detailList">
    {#if selectedEntry}
      <div>
        <dt>訳語</dt>
        <dd>{selectedEntry.translation}</dd>
      </div>
    {/if}
  </dl>
</section>

<style>
  .shell-card {
    padding: 18px;
    border-radius: 16px;
    border: 1px solid var(--line);
    background: rgba(16, 13, 11, 0.58);
    color: var(--text);
  }

  .detail-panel,
  .detail-title,
  .detail-grid,
  .detail-list {
    display: grid;
    gap: 10px;
  }

  .detail-panel {
    min-width: 0;
    padding: 20px;
    gap: 14px;
  }

  .detail-head,
  .detail-actions {
    display: flex;
    flex-wrap: wrap;
    justify-content: space-between;
    align-items: center;
    gap: 10px;
  }

  #detailSubline,
  #detailStatusMessage,
  p,
  dt {
    color: var(--muted);
  }

  .field-label,
  dt {
    font-size: 12px;
    letter-spacing: 0.08em;
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

  .detail-title {
    width: 100%;
    min-width: 0;
    box-sizing: border-box;
    gap: 8px;
    padding: 16px;
    border: 1px solid var(--line);
    border-radius: 18px;
    background: rgba(255, 255, 255, 0.03);
  }

  .detail-grid {
    min-height: 0;
    padding: 10px;
    gap: 10px;
    grid-template-columns: repeat(2, minmax(0, 1fr));
    align-items: stretch;
    border: 1px solid var(--line);
    border-radius: 12px;
    background: rgba(0, 0, 0, 0.2);
  }

  .detail-card {
    min-width: 0;
    box-sizing: border-box;
    display: grid;
    align-content: start;
    gap: 10px;
    padding: 10px;
    border: 1px solid var(--line);
    border-radius: 8px;
    background: rgba(255, 255, 255, 0.03);
  }

  .detail-meta-value {
    display: flex;
    width: 100%;
    min-height: 56px;
    box-sizing: border-box;
    align-items: center;
    font-size: 0.92rem;
    line-height: 1.35;
    padding: 10px 14px;
    border-radius: 12px;
    background: rgba(255, 255, 255, 0.04);
    overflow-wrap: anywhere;
  }

  .detail-tags {
    display: flex;
    flex-wrap: wrap;
    gap: 8px;
  }

  .status-pill {
    padding: 6px 10px;
    border-radius: 999px;
    border: 1px solid var(--line);
    background: rgba(255, 255, 255, 0.03);
  }

  .empty-state {
    color: var(--muted);
    padding: 8px;
  }

  .detail-list div {
    padding: 14px 16px;
    border-radius: 12px;
    border: 1px solid var(--line);
    background: rgba(255, 255, 255, 0.03);
  }

  .detail-list dd {
    margin: 0;
  }

  @media (max-width: 1180px) {
    .detail-grid {
      grid-template-columns: 1fr;
    }
  }
</style>
