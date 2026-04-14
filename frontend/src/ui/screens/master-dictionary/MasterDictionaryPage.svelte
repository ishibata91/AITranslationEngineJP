<script lang="ts">
  import { onMount } from "svelte"

  import type {
    CreateMasterDictionaryScreenController,
    MasterDictionaryScreenControllerContract
  } from "@application/contract/master-dictionary"

  interface Props {
    createController: CreateMasterDictionaryScreenController | null
  }

  let { createController }: Props = $props()

  function resolveController(): MasterDictionaryScreenControllerContract {
    if (!createController) {
      throw new Error(
        "master dictionary screen controller factory is not provided"
      )
    }

    return createController()
  }

  const controller = resolveController()
  let viewModel = $state(controller.getViewModel())

  const unsubscribe = controller.subscribe((nextViewModel) => {
    viewModel = nextViewModel
  })

  onMount(() => {
    void controller.mount()

    return () => {
      unsubscribe()
      controller.dispose()
    }
  })

  function chooseXmlFile(): void {
    if (viewModel.isImportRunning) {
      return
    }

    const input = document.getElementById("xmlFileInput")
    if (input instanceof HTMLInputElement) {
      input.click()
    }
  }

  function handleXmlSelected(event: Event): void {
    const target = event.currentTarget
    if (!(target instanceof HTMLInputElement)) {
      return
    }

    const file = target.files?.[0] ?? null
    controller.stageXmlImport(file)
  }

  function resetImportSelection(): void {
    if (viewModel.isImportRunning) {
      return
    }

    const input = document.getElementById("xmlFileInput")
    if (input instanceof HTMLInputElement) {
      input.value = ""
    }

    controller.resetImportSelection()
  }
</script>

<section class="master-dictionary-shell" id="masterDictionaryView">
  <section class="shell-card">
    <div class="hero-top">
      <div>
        <p class="eyebrow">基盤データ</p>
        <h2>マスター辞書</h2>
      </div>
      <p class="gateway-status" id="gatewayStatus">
        Gateway: {viewModel.gatewayStatus}
      </p>
    </div>
    <p class="lead">
      一覧、詳細、作成、更新、削除、XML 取り込みを同じ画面で操作できます。
    </p>
    <p
      class="error-text"
      hidden={!viewModel.errorMessage}
      id="masterDictionaryError"
    >
      {viewModel.errorMessage}
    </p>
  </section>

  <section class="shell-card import-shell" aria-labelledby="importHeading">
    <div class="import-top">
      <div>
        <p class="eyebrow">XMLから取り込み</p>
        <h3 id="importHeading">取り込み導線</h3>
      </div>
      <button
        class="button-secondary"
        id="chooseXmlButton"
        onclick={chooseXmlFile}
        type="button"
      >
        ファイルを選択
      </button>
    </div>
    <p class="mini-text" id="importStateText">
      ファイルを選ぶと取込バーが表示されます。
    </p>

    <input
      accept=".xml,text/xml,application/xml"
      class="file-input"
      id="xmlFileInput"
      onchange={handleXmlSelected}
      type="file"
    />

    <div class="file-picker">
      <span class="eyebrow">選択ファイル</span>
      <span class="file-name" id="selectedFileName"
        >{viewModel.selectedFileName}</span
      >
    </div>

    <div class="import-bar" hidden={!viewModel.hasStagedFile} id="importBar">
      <div class="import-bar-head">
        <strong id="importFileTitle">{viewModel.selectedFileName}</strong>
        <div class="import-actions">
          <button
            class="button-primary"
            disabled={viewModel.isImportRunning}
            id="startImportButton"
            onclick={() => void controller.startImport()}
            type="button"
          >
            この XML を取り込む
          </button>
          <button
            class="button-secondary"
            disabled={viewModel.isImportRunning}
            id="resetImportButton"
            onclick={resetImportSelection}
            type="button"
          >
            選び直す
          </button>
        </div>
      </div>
      <div class="status-line">
        <p id="importStatusText">{viewModel.importStatusText}</p>
        <strong id="importStatusValue">{viewModel.importStatusValue}</strong>
      </div>
      <div class="progress-track">
        <div
          class="progress-fill"
          id="importProgressFill"
          style={`width: ${viewModel.importProgress}%;`}
        ></div>
      </div>
      <div
        class="import-result"
        hidden={!viewModel.importSummary}
        id="importResult"
      >
        <div class="import-result-head">
          <strong id="importResultHeadline"
            >XML取り込みを一覧と詳細へ反映しました。</strong
          >
          <span class="status-pill" id="importResultCount"
            >新規取込 {viewModel.importSummary?.importedCount ?? 0} 件</span
          >
        </div>
        <p id="importResultMessage">
          {viewModel.importSummary
            ? `「${viewModel.importSummary.fileName}」の取込を完了し、同じ画面に反映しました。`
            : "-"}
        </p>
        <dl class="result-grid">
          <div>
            <dt>更新件数</dt>
            <dd id="importResultUpdatedCount">
              {viewModel.importSummary?.updatedCount ?? "-"}
            </dd>
          </div>
          <div>
            <dt>取込後の一覧総件数</dt>
            <dd id="importResultListCount">
              {viewModel.importSummary?.totalCount ?? "-"}
            </dd>
          </div>
          <div>
            <dt>選択状態</dt>
            <dd id="importResultSelection">
              {viewModel.importSummary?.selectedSource ?? "-"}
            </dd>
          </div>
          <div>
            <dt>詳細表示</dt>
            <dd id="importResultDetail">
              {viewModel.selectedEntry?.translation ?? "-"}
            </dd>
          </div>
        </dl>
      </div>
    </div>
  </section>

  <section class="content-grid">
    <section class="shell-card" aria-labelledby="listHeading">
      <div class="toolbar-head">
        <div>
          <h3 id="listHeading">辞書一覧</h3>
          <p id="listHeadline">{viewModel.listHeadline}</p>
        </div>
        <div class="toolbar-head-actions">
          <button
            class="button-primary"
            id="createButton"
            onclick={() => controller.openCreateModal()}
            type="button">新規登録</button
          >
          <p class="mini-text" id="pageStatusText">
            {viewModel.pageStatusText}
          </p>
        </div>
      </div>

      <div class="filter-grid">
        <label class="field-label" for="searchInput">検索</label>
        <input
          class="search-field"
          id="searchInput"
          oninput={(event) => controller.handleSearchInput(event)}
          placeholder="原文・訳語・IDで検索"
          type="search"
          value={viewModel.query}
        />

        <label class="field-label" for="categorySelect">カテゴリ</label>
        <select
          class="select-field"
          id="categorySelect"
          onchange={(event) => controller.handleCategoryChange(event)}
          value={viewModel.category}
        >
          {#each viewModel.categoryOptions as option (option)}
            <option value={option}>{option}</option>
          {/each}
        </select>
      </div>

      <div class="list-stack" id="listStack" aria-live="polite">
        {#if viewModel.entries.length === 0}
          <div class="empty-state">一致するエントリがありません</div>
        {:else}
          {#each viewModel.entries as entry (entry.id)}
            <button
              class="list-row"
              class:is-selected={viewModel.selectedId === entry.id}
              onclick={() => void controller.selectRow(entry.id)}
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

      <div class="pager-shell">
        <div class="mini-text" id="selectionStatus">
          {viewModel.selectionStatusText}
        </div>
        <div class="pager-actions">
          <button
            class="button-secondary"
            disabled={viewModel.page === 0}
            id="prevPageButton"
            onclick={() => controller.goToPrevPage()}
            type="button"
          >
            前の30件
          </button>
          <button
            class="button-secondary"
            disabled={viewModel.page + 1 >= viewModel.totalPages}
            id="nextPageButton"
            onclick={() => controller.goToNextPage()}
            type="button"
          >
            次の30件
          </button>
        </div>
      </div>
    </section>

    <section class="shell-card" aria-labelledby="detailHeading">
      <div class="detail-head">
        <div>
          <h3 id="detailHeading">詳細</h3>
          <p id="detailSubline">{viewModel.detailSublineText}</p>
        </div>
        <div class="detail-actions">
          <button
            class="button-secondary"
            disabled={!viewModel.selectedEntry}
            id="editButton"
            onclick={() => controller.openEditModal()}
            type="button"
          >
            更新
          </button>
          <button
            class="button-danger"
            disabled={!viewModel.selectedEntry}
            id="deleteButton"
            onclick={() => controller.openDeleteModal()}
            type="button"
          >
            削除
          </button>
        </div>
      </div>

      <div class="detail-tags" id="detailTags">
        {#if viewModel.selectedEntry}
          <span class="status-pill">{viewModel.selectedEntry.category}</span>
          <span class="status-pill">{viewModel.selectedEntry.origin}</span>
        {/if}
      </div>
      <strong id="detailTitle"
        >{viewModel.selectedEntry?.source ??
          "表示できるエントリがありません"}</strong
      >
      <p id="detailTranslation">
        {viewModel.selectedEntry?.translation ?? "検索条件を変更してください。"}
      </p>
      <div class="detail-grid" id="detailGrid">
        {#if viewModel.selectedEntry}
          <div class="detail-card">
            <div class="field-label">ID</div>
            <strong>{viewModel.selectedEntry.id}</strong>
          </div>
          <div class="detail-card">
            <div class="field-label">最終更新</div>
            <strong>{viewModel.selectedEntry.updatedAt}</strong>
          </div>
        {:else}
          <div class="empty-state">
            一覧に表示できるエントリが戻ると、詳細も同じ画面で切り替わります。
          </div>
        {/if}
      </div>
      <p id="detailStatusMessage">
        {viewModel.selectedEntry
          ? `${viewModel.selectedEntry.origin} / 最終更新 ${viewModel.selectedEntry.updatedAt}`
          : "一覧から別のエントリを選択すると、ここも切り替わります。"}
      </p>
      <dl class="detail-list" id="detailList">
        {#if viewModel.selectedEntry}
          <div>
            <dt>訳語</dt>
            <dd>{viewModel.selectedEntry.translation}</dd>
          </div>
          <div>
            <dt>現在の扱い</dt>
            <dd>{viewModel.selectedEntry.note}</dd>
          </div>
        {/if}
      </dl>
    </section>
  </section>
</section>

<div
  aria-hidden={!(
    viewModel.modalState === "create" || viewModel.modalState === "edit"
  )}
  class="modal-backdrop"
  hidden={!(
    viewModel.modalState === "create" || viewModel.modalState === "edit"
  )}
  id="editModal"
  role="dialog"
>
  <section aria-labelledby="editModalTitle" class="modal-card">
    <div class="eyebrow" id="editModalEyebrow">
      {viewModel.modalState === "create" ? "新規登録" : "更新"}
    </div>
    <h3 id="editModalTitle">
      {viewModel.modalState === "create" ? "新規登録" : "更新"}
    </h3>
    <p id="editModalDescription">
      {viewModel.modalState === "create"
        ? "辞書エントリの内容を入力します。"
        : "選択中の辞書エントリを編集します。"}
    </p>
    <div class="field-grid">
      <label class="field-label" for="formSource">原文</label>
      <input
        class="text-field"
        id="formSource"
        type="text"
        value={viewModel.formSource}
        oninput={(event) => controller.setFormSource(event)}
      />

      <label class="field-label" for="formCategory">カテゴリ</label>
      <select
        class="select-field"
        id="formCategory"
        value={viewModel.formCategory}
        onchange={(event) => controller.setFormCategory(event)}
      >
        {#each viewModel.categoryOptions.filter((item) => item !== "すべて") as option (option)}
          <option value={option}>{option}</option>
        {/each}
      </select>

      <label class="field-label" for="formOrigin">由来</label>
      <select
        class="select-field"
        id="formOrigin"
        value={viewModel.formOrigin}
        onchange={(event) => controller.setFormOrigin(event)}
      >
        <option value="手動登録">手動登録</option>
        <option value="確認待ち">確認待ち</option>
        <option value="XML取込">XML取込</option>
      </select>

      <label class="field-label" for="formTranslation">訳語</label>
      <textarea
        class="textarea-field"
        id="formTranslation"
        value={viewModel.formTranslation}
        oninput={(event) => controller.setFormTranslation(event)}
      ></textarea>
    </div>
    <div class="modal-actions">
      <button
        class="button-secondary"
        id="closeEditModalButton"
        onclick={() => controller.closeEditModal()}
        type="button">閉じる</button
      >
      <button
        class="button-primary"
        id="saveEntryButton"
        onclick={() => void controller.saveCurrentEntry()}
        type="button">保存する</button
      >
    </div>
  </section>
</div>

<div
  aria-hidden={viewModel.modalState !== "delete"}
  class="modal-backdrop"
  hidden={viewModel.modalState !== "delete"}
  id="deleteModal"
  role="dialog"
>
  <section aria-labelledby="deleteModalTitle" class="modal-card">
    <h3 id="deleteModalTitle">削除の確認</h3>
    <p>このエントリを削除すると、一覧から見えなくなります。</p>
    <div class="delete-target">
      <strong id="deleteTargetTitle"
        >{viewModel.selectedEntry?.source ?? "-"}</strong
      >
      <p id="deleteTargetMeta">
        {viewModel.selectedEntry
          ? `${viewModel.selectedEntry.translation} / ID ${viewModel.selectedEntry.id}`
          : "-"}
      </p>
    </div>
    <div class="modal-actions">
      <button
        class="button-secondary"
        id="closeDeleteModalButton"
        onclick={() => controller.closeDeleteModal()}
        type="button">やめる</button
      >
      <button
        class="button-danger"
        id="confirmDeleteButton"
        onclick={() => void controller.deleteCurrentEntry()}
        type="button">削除する</button
      >
    </div>
  </section>
</div>

<style>
  .master-dictionary-shell {
    display: grid;
    gap: 16px;
  }

  .shell-card {
    padding: 18px;
    border-radius: 16px;
    border: 1px solid var(--line);
    background: rgba(16, 13, 11, 0.58);
  }

  .hero-top,
  .import-top,
  .import-bar-head,
  .toolbar-head,
  .pager-shell,
  .detail-head,
  .modal-actions,
  .import-actions,
  .toolbar-head-actions,
  .delete-target,
  .import-result-head {
    display: flex;
    flex-wrap: wrap;
    justify-content: space-between;
    align-items: center;
    gap: 10px;
  }

  .lead,
  .mini-text,
  #listHeadline,
  #detailSubline,
  #detailStatusMessage,
  .gateway-status,
  #importStatusText,
  #importResultMessage,
  dt,
  p {
    color: var(--muted);
  }

  .error-text {
    color: #ffc0ab;
  }

  .error-text[hidden] {
    display: none;
  }

  .eyebrow,
  .field-label,
  dt {
    font-size: 12px;
    letter-spacing: 0.08em;
  }

  .button-primary,
  .button-secondary,
  .button-danger {
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

  .button-danger {
    color: #3d1512;
    background: linear-gradient(135deg, #ffc0ab 0%, #ff9975 100%);
  }

  .import-shell,
  .import-bar,
  .content-grid,
  .filter-grid,
  .field-grid,
  .result-grid,
  .detail-grid,
  .detail-list {
    display: grid;
    gap: 10px;
  }

  .import-bar[hidden],
  .import-result[hidden],
  .modal-backdrop[hidden] {
    display: none !important;
    pointer-events: none;
  }

  .file-picker {
    display: inline-flex;
    gap: 10px;
    align-items: center;
  }

  .file-input {
    position: absolute;
    width: 1px;
    height: 1px;
    margin: -1px;
    padding: 0;
    border: 0;
    overflow: hidden;
    clip: rect(0, 0, 0, 0);
    clip-path: inset(50%);
    white-space: nowrap;
    pointer-events: none;
  }

  .file-name,
  .status-pill {
    padding: 6px 10px;
    border-radius: 999px;
    border: 1px solid var(--line);
    background: rgba(255, 255, 255, 0.03);
  }

  .progress-track {
    height: 10px;
    border-radius: 999px;
    background: rgba(255, 255, 255, 0.08);
    overflow: hidden;
  }

  .progress-fill {
    height: 100%;
    background: linear-gradient(90deg, var(--primary) 0%, #f5ca72 100%);
    transition: width 180ms ease;
  }

  .content-grid {
    grid-template-columns: minmax(0, 1fr) minmax(0, 1fr);
  }

  .search-field,
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

  .list-stack,
  .detail-grid {
    border: 1px solid var(--line);
    border-radius: 10px;
    min-height: 200px;
    padding: 10px;
    background: rgba(0, 0, 0, 0.2);
  }

  .list-row {
    width: 100%;
    display: grid;
    grid-template-columns:
      minmax(0, 1.2fr) minmax(0, 1.2fr) minmax(0, 0.9fr)
      auto;
    gap: 10px;
    align-items: center;
    border: 1px solid rgba(255, 186, 56, 0.12);
    border-radius: 8px;
    background: rgba(255, 255, 255, 0.03);
    color: var(--text);
    padding: 8px 10px;
    text-align: left;
    cursor: pointer;
  }

  .list-row.is-selected {
    border-color: var(--line-strong);
    background: rgba(255, 186, 56, 0.12);
  }

  .row-cell {
    min-width: 0;
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

  .detail-card {
    padding: 10px;
    border-radius: 8px;
    border: 1px solid var(--line);
    background: rgba(255, 255, 255, 0.03);
  }

  .detail-list div {
    padding: 10px;
    border-radius: 8px;
    border: 1px solid var(--line);
    background: rgba(255, 255, 255, 0.03);
  }

  .detail-list dd {
    margin: 0;
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
    display: grid;
    gap: 12px;
  }

  @media (max-width: 980px) {
    .content-grid {
      grid-template-columns: 1fr;
    }

    .list-row {
      grid-template-columns: 1fr;
    }

    .row-id {
      text-align: left;
    }
  }
</style>
