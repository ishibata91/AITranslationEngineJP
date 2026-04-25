<script lang="ts">
  import { onMount } from "svelte"

  import type {
    CreateTranslationInputScreenController,
    TranslationInputScreenControllerContract
  } from "@application/contract/translation-input"
  import {
    ERROR_LABELS,
    STATUS_LABELS,
    WARNING_LABELS
  } from "@application/presenter/translation-input"

  interface Props {
    createController: CreateTranslationInputScreenController | null
  }

  let { createController }: Props = $props()

  function resolveController(): TranslationInputScreenControllerContract {
    if (!createController) {
      throw new Error("translation input screen controller factory is not provided")
    }

    return createController()
  }

  const controller = resolveController()
  let viewModel = $state(controller.getViewModel())
  let fileInput: HTMLInputElement | null = null

  const unsubscribe = controller.subscribe((nextViewModel) => {
    viewModel = nextViewModel
  })

  function clearJsonFileInput(): void {
    if (fileInput) {
      fileInput.value = ""
    }
  }

  $effect(() => {
    if (!viewModel.hasStagedFile) {
      clearJsonFileInput()
    }
  })

  onMount(() => {
    void controller.mount()

    return () => {
      unsubscribe()
      controller.dispose()
    }
  })

  function chooseJsonFile(): void {
    if (viewModel.isImporting) {
      return
    }

    clearJsonFileInput()
    fileInput?.click()
  }

  function handleJsonSelected(event: Event): void {
    const target = event.currentTarget
    if (!(target instanceof HTMLInputElement)) {
      return
    }

    const file = target.files?.[0] ?? null
    void controller.stageJsonImport(file)
  }

  function resetImportSelection(): void {
    clearJsonFileInput()
    controller.resetImportSelection()
  }

  function formatStatus(localStatus: string): string {
    return STATUS_LABELS[localStatus as keyof typeof STATUS_LABELS] ?? localStatus
  }

  function formatErrorKind(errorKind: string | null): string {
    if (!errorKind) {
      return "-"
    }

    return ERROR_LABELS[errorKind] ?? errorKind
  }

  function formatWarningKind(kind: string): string {
    return WARNING_LABELS[kind] ?? kind
  }

  function formatDate(timestamp: string): string {
    if (!timestamp) {
      return "-"
    }

    const date = new Date(timestamp)
    if (Number.isNaN(date.getTime())) {
      return timestamp
    }

    return date.toLocaleString("ja-JP")
  }
</script>

<section class="review-shell" id="translationInputReviewView">
  <section class="review-card hero-card">
    <div class="hero-head">
      <div>
        <p class="eyebrow">translation-management</p>
        <h2>Input Review</h2>
      </div>
      <p class="gateway-status">Gateway: {viewModel.gatewayStatus}</p>
    </div>
    <p class="lead">
      1 JSON file を登録し、一覧、概要、sample field、error kind、retry / rebuild 状態を同じページで確認します。
    </p>
    <p class="status-copy">
      <strong>{viewModel.operationStatusLabel}</strong>
      <span>{viewModel.operationStatusText}</span>
    </p>
    <p class="error-text" hidden={!viewModel.errorMessage}>
      {viewModel.errorMessage}
    </p>
  </section>

  <section class="review-card import-card" aria-labelledby="inputReviewImportHeading">
    <div class="section-head">
      <div>
        <p class="eyebrow">register input</p>
        <h3 id="inputReviewImportHeading">JSON file 登録</h3>
      </div>
      <button class="button-secondary" onclick={chooseJsonFile} type="button">
        JSON file を選択
      </button>
    </div>
    <input
      accept=".json,application/json"
      bind:this={fileInput}
      class="file-input"
      id="translationInputFile"
      onchange={handleJsonSelected}
      type="file"
    />
    <dl class="import-grid">
      <div>
        <dt>file name</dt>
        <dd>{viewModel.stagedFileName}</dd>
      </div>
      <div>
        <dt>file path</dt>
        <dd>{viewModel.stagedFilePath}</dd>
      </div>
      <div>
        <dt>file hash</dt>
        <dd class="hash-text">{viewModel.stagedFileHash}</dd>
      </div>
    </dl>
    <div class="action-row">
      <button
        class="button-primary"
        disabled={!viewModel.canImport}
        onclick={() => void controller.startImport()}
        type="button"
      >
        この JSON を登録
      </button>
      <button
        class="button-secondary"
        disabled={!viewModel.hasStagedFile || viewModel.isImporting}
        onclick={resetImportSelection}
        type="button"
      >
        選び直す
      </button>
    </div>
  </section>

  <section class="content-grid">
    <section class="review-card list-card" aria-labelledby="inputReviewListHeading">
      <div class="section-head">
        <div>
          <p class="eyebrow">review list</p>
          <h3 id="inputReviewListHeading">入力ファイル一覧</h3>
        </div>
        <p class="mini-text">{viewModel.totalItemCountLabel}</p>
      </div>

      {#if viewModel.items.length === 0}
        <div class="empty-state">
          <p>{viewModel.emptyStateText}</p>
        </div>
      {:else}
        <div class="review-list" role="list">
          {#each viewModel.items as item (item.localId)}
            <div role="listitem">
              <button
                aria-pressed={item.localId === viewModel.selectedItemId ? "true" : "false"}
                class="review-item"
                class:is-selected={item.localId === viewModel.selectedItemId}
                onclick={() => controller.selectItem(item.localId)}
                type="button"
              >
                <div class="review-item-head">
                  <div>
                    <strong>{item.fileName}</strong>
                    <p>{item.filePath}</p>
                  </div>
                  <span class="status-pill">{formatStatus(item.status)}</span>
                </div>
                <dl class="review-meta">
                  <div>
                    <dt>file hash</dt>
                    <dd class="hash-text">{item.fileHash}</dd>
                  </div>
                  <div>
                    <dt>import timestamp</dt>
                    <dd>{formatDate(item.importTimestamp)}</dd>
                  </div>
                  <div>
                    <dt>登録状態</dt>
                    <dd>{item.accepted ? "accepted" : "rejected"}</dd>
                  </div>
                  <div>
                    <dt>再構築可否</dt>
                    <dd>{item.canRebuild ? "rebuild 可" : "不可"}</dd>
                  </div>
                </dl>
                <p class="mini-text">
                  error kind: {formatErrorKind(item.errorKind)}
                </p>
              </button>
            </div>
          {/each}
        </div>
      {/if}
    </section>

    <section class="review-card detail-card" aria-labelledby="inputReviewDetailHeading">
      <div class="section-head">
        <div>
          <p class="eyebrow">selected input</p>
          <h3 id="inputReviewDetailHeading">概要と sample field</h3>
        </div>
        <button
          class="button-secondary"
          disabled={!viewModel.canRebuildSelected || viewModel.isRebuilding}
          onclick={() => void controller.rebuildSelected()}
          type="button"
        >
          cache を再構築
        </button>
      </div>

      <p class="mini-text">{viewModel.selectionStatusText}</p>
      <div class="result-callout">
        <strong>{viewModel.latestOutcomeTitle}</strong>
        <p>{viewModel.latestOutcomeText}</p>
      </div>

      {#if viewModel.selectedItem}
        <div class="detail-stack">
          <dl class="summary-grid">
            <div>
              <dt>file name</dt>
              <dd>{viewModel.selectedItem.fileName}</dd>
            </div>
            <div>
              <dt>file path</dt>
              <dd>{viewModel.selectedItem.filePath}</dd>
            </div>
            <div>
              <dt>file hash</dt>
              <dd class="hash-text">{viewModel.selectedItem.fileHash}</dd>
            </div>
            <div>
              <dt>import timestamp</dt>
              <dd>{formatDate(viewModel.selectedItem.importTimestamp)}</dd>
            </div>
            <div>
              <dt>translation record count</dt>
              <dd>{viewModel.selectedItem.summary?.translationRecordCount ?? 0}</dd>
            </div>
            <div>
              <dt>translation field count</dt>
              <dd>{viewModel.selectedItem.summary?.translationFieldCount ?? 0}</dd>
            </div>
            <div>
              <dt>target plugin</dt>
              <dd>{viewModel.selectedItem.summary?.input.targetPluginName ?? "-"}</dd>
            </div>
            <div>
              <dt>source tool</dt>
              <dd>{viewModel.selectedItem.summary?.input.sourceTool ?? "-"}</dd>
            </div>
          </dl>

          <section class="detail-section">
            <div class="section-head section-head-compact">
              <h4>カテゴリ別件数</h4>
              <span class="mini-text">
                {viewModel.selectedItem.summary?.categories.length ?? 0} 件
              </span>
            </div>
            {#if (viewModel.selectedItem.summary?.categories.length ?? 0) > 0}
              <div class="chip-grid">
                {#each viewModel.selectedItem.summary?.categories ?? [] as category (`${category.category}:${category.recordCount}:${category.fieldCount}`)}
                  <article class="chip-card">
                    <strong>{category.category}</strong>
                    <p>record {category.recordCount} / field {category.fieldCount}</p>
                  </article>
                {/each}
              </div>
            {:else}
              <p class="mini-text">カテゴリ別件数はまだありません。</p>
            {/if}
          </section>

          <section class="detail-section">
            <div class="section-head section-head-compact">
              <h4>sample field</h4>
              <span class="mini-text">
                {viewModel.selectedItem.summary?.sampleFields.length ?? 0} 件
              </span>
            </div>
            {#if (viewModel.selectedItem.summary?.sampleFields.length ?? 0) > 0}
              <div class="sample-grid">
                {#each viewModel.selectedItem.summary?.sampleFields ?? [] as field (`${field.recordType}:${field.subrecordType}:${field.formId}:${field.editorId}`)}
                  <article class="sample-card">
                    <div class="sample-head">
                      <strong>{field.recordType}:{field.subrecordType}</strong>
                      <span>{field.translatable ? "translatable" : "non-translatable"}</span>
                    </div>
                    <p>{field.sourceText || "-"}</p>
                    <dl>
                      <div>
                        <dt>formId</dt>
                        <dd>{field.formId || "-"}</dd>
                      </div>
                      <div>
                        <dt>editorId</dt>
                        <dd>{field.editorId || "-"}</dd>
                      </div>
                    </dl>
                  </article>
                {/each}
              </div>
            {:else}
              <p class="mini-text">sample field はまだありません。</p>
            {/if}
          </section>

          <section class="detail-section">
            <div class="section-head section-head-compact">
              <h4>error / warning</h4>
              <span class="mini-text">retry / rebuild 判断用</span>
            </div>
            <div class="chip-grid">
              {#if viewModel.selectedItem.errorKind}
                <article class="chip-card chip-card-error">
                  <strong>{formatErrorKind(viewModel.selectedItem.errorKind)}</strong>
                  <p>登録または再構築で返された error kind</p>
                </article>
              {/if}
              {#each viewModel.selectedItem.warnings as warning (`${warning.kind}:${warning.recordType}:${warning.subrecordType}:${warning.message}`)}
                <article class="chip-card chip-card-warning">
                  <strong>{formatWarningKind(warning.kind)}</strong>
                  <p>{warning.message}</p>
                </article>
              {/each}
              {#if !viewModel.selectedItem.errorKind && viewModel.selectedItem.warnings.length === 0}
                <article class="chip-card">
                  <strong>問題なし</strong>
                  <p>登録状態は正常です。</p>
                </article>
              {/if}
            </div>
          </section>
        </div>
      {:else}
        <div class="empty-state">
          <p>選択中の入力データはありません。</p>
        </div>
      {/if}
    </section>
  </section>
</section>

<style>
  .review-shell {
    display: grid;
    gap: 1.25rem;
  }

  .review-card {
    border: 1px solid var(--line);
    border-radius: 20px;
    padding: 1.25rem;
    background: rgba(28, 23, 20, 0.74);
    box-shadow: var(--shadow);
    backdrop-filter: blur(18px);
  }

  .hero-card,
  .import-card {
    display: grid;
    gap: 0.9rem;
  }

  .hero-head,
  .section-head {
    display: flex;
    justify-content: space-between;
    gap: 1rem;
    align-items: start;
  }

  .section-head-compact {
    margin-bottom: 0.75rem;
  }

  .eyebrow {
    margin-bottom: 0.2rem;
    color: var(--muted);
    font-size: 0.76rem;
    letter-spacing: 0.16em;
    text-transform: uppercase;
  }

  h2,
  h3,
  h4,
  p,
  dl,
  dt,
  dd {
    margin: 0;
  }

  .lead,
  .status-copy,
  .mini-text,
  .review-item p,
  .sample-card p,
  .chip-card p {
    color: var(--muted);
  }

  .status-copy {
    display: grid;
    gap: 0.2rem;
  }

  .gateway-status,
  .status-pill {
    padding: 0.32rem 0.7rem;
    border: 1px solid var(--line-strong);
    border-radius: 999px;
    color: var(--primary);
    font-size: 0.82rem;
    white-space: nowrap;
  }

  .error-text {
    padding: 0.8rem 1rem;
    border-radius: 14px;
    background: rgba(255, 104, 63, 0.16);
    color: #ffd8c5;
  }

  .file-input {
    display: none;
  }

  .import-grid,
  .summary-grid,
  .review-meta {
    display: grid;
    grid-template-columns: repeat(auto-fit, minmax(180px, 1fr));
    gap: 0.8rem;
  }

  .import-grid div,
  .summary-grid div,
  .review-meta div,
  .sample-card dl div {
    display: grid;
    gap: 0.2rem;
  }

  dt {
    color: var(--muted);
    font-size: 0.8rem;
  }

  dd {
    font-size: 0.95rem;
  }

  .hash-text {
    word-break: break-all;
  }

  .action-row {
    display: flex;
    gap: 0.75rem;
    flex-wrap: wrap;
  }

  .button-primary,
  .button-secondary,
  .review-item {
    font: inherit;
  }

  .button-primary,
  .button-secondary {
    border-radius: 999px;
    padding: 0.7rem 1.1rem;
    border: 1px solid var(--line-strong);
    cursor: pointer;
  }

  .button-primary {
    background: var(--primary);
    color: #2b1900;
  }

  .button-secondary {
    background: transparent;
    color: var(--text);
  }

  .button-primary:disabled,
  .button-secondary:disabled {
    cursor: not-allowed;
    opacity: 0.55;
  }

  .content-grid {
    display: grid;
    gap: 1.25rem;
    grid-template-columns: minmax(0, 1fr) minmax(0, 1.2fr);
  }

  .list-card,
  .detail-card,
  .detail-stack,
  .detail-section {
    display: grid;
    gap: 0.9rem;
  }

  .review-list,
  .sample-grid,
  .chip-grid {
    display: grid;
    gap: 0.8rem;
  }

  .review-item {
    display: grid;
    gap: 0.8rem;
    text-align: left;
    padding: 1rem;
    border-radius: 18px;
    border: 1px solid rgba(255, 186, 56, 0.18);
    background: rgba(22, 18, 17, 0.78);
    color: var(--text);
    cursor: pointer;
  }

  .review-item.is-selected {
    border-color: var(--primary);
    box-shadow: 0 0 0 1px rgba(255, 186, 56, 0.18);
  }

  .review-item-head,
  .sample-head {
    display: flex;
    justify-content: space-between;
    gap: 0.75rem;
    align-items: start;
  }

  .result-callout,
  .empty-state,
  .chip-card,
  .sample-card {
    border-radius: 18px;
    padding: 1rem;
    background: rgba(18, 15, 14, 0.76);
    border: 1px solid rgba(255, 186, 56, 0.12);
  }

  .chip-card-error {
    border-color: rgba(255, 104, 63, 0.4);
  }

  .chip-card-warning {
    border-color: rgba(255, 186, 56, 0.34);
  }

  .sample-grid {
    grid-template-columns: repeat(auto-fit, minmax(220px, 1fr));
  }

  .sample-card {
    display: grid;
    gap: 0.65rem;
  }

  .sample-card dl {
    display: grid;
    gap: 0.5rem;
  }

  @media (max-width: 960px) {
    .content-grid {
      grid-template-columns: 1fr;
    }
  }
</style>