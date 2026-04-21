<script lang="ts">
  import { onMount } from "svelte"

  import type {
    CreateMasterPersonaScreenController,
    MasterPersonaScreenControllerContract
  } from "@application/contract/master-persona"

  interface Props {
    createController: CreateMasterPersonaScreenController | null
  }

  let { createController }: Props = $props()

  function resolveController(): MasterPersonaScreenControllerContract {
    if (!createController) {
      throw new Error(
        "master persona screen controller factory is not provided"
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

  function chooseJsonFile(): void {
    const input = document.getElementById("masterPersonaJsonInput")
    if (input instanceof HTMLInputElement) {
      input.click()
    }
  }

  function handleJsonSelected(event: Event): void {
    const target = event.currentTarget
    if (!(target instanceof HTMLInputElement)) {
      return
    }
    controller.stageJsonSelection(target.files?.[0] ?? null)
  }

  function resetJsonSelection(): void {
    const input = document.getElementById("masterPersonaJsonInput")
    if (input instanceof HTMLInputElement) {
      input.value = ""
    }
    controller.resetJsonSelection()
  }
</script>

<section class="master-persona-shell" id="masterPersonaView">
  <section class="master-persona-panel overview-panel">
    <div class="hero-top">
      <div>
        <p class="eyebrow">AI生成</p>
        <h2>JSONから NPC ペルソナを生成</h2>
      </div>
      <div class="status-row">
        <span class="status-pill status-accent">{viewModel.runStatus.runState}</span>
        <span class="status-pill">Gateway: {viewModel.gatewayStatus}</span>
      </div>
    </div>
    <p class="lead">
      extractData.pas JSON を入力にして、未生成のマスターペルソナだけを追加します。
    </p>
    <p class="error-text" hidden={!viewModel.errorMessage} id="masterPersonaError">
      {viewModel.errorMessage}
    </p>
  </section>

  <section class="generator-grid">
    <section class="master-persona-panel" aria-labelledby="settingsHeading">
      <div class="section-head">
        <div>
          <p class="eyebrow">AI 設定</p>
          <h3 id="settingsHeading">この画面で使う設定</h3>
        </div>
        <span class="status-pill">{viewModel.aiProviderLabel}</span>
      </div>

      <label class="field-group" for="providerSelect">
        <span class="field-label">AI サービス</span>
        <select
          class="select-field"
          id="providerSelect"
          onchange={(event) => controller.setAIProvider(event)}
          value={viewModel.aiSettings.provider}
        >
          <option value="gemini">Gemini</option>
          <option value="lm_studio">LM Studio</option>
          <option value="xai">xAI</option>
        </select>
      </label>

      <label class="field-group" for="modelInput">
        <span class="field-label">モデル</span>
        <input
          class="text-field"
          id="modelInput"
          oninput={(event) => controller.setAIModel(event)}
          value={viewModel.aiSettings.model}
        />
      </label>

      <label class="field-group" for="apiKeyInput">
        <span class="field-label">API キー</span>
        <input
          class="text-field"
          id="apiKeyInput"
          oninput={(event) => controller.setAPIKey(event)}
          placeholder="選択した AI サービスの API キーを入力"
          value={viewModel.aiSettings.apiKey}
        />
      </label>

      <div class="prompt-copy" id="promptTemplateDescription">
        {viewModel.promptTemplateDescription}
      </div>

      <div class="inline-actions">
        <span class="mini-text" id="aiSettingsMessage">{viewModel.aiSettingsMessage}</span>
        <button
          class="button-secondary"
          id="saveAiSettingsButton"
          onclick={() => void controller.saveAISettings()}
          type="button"
        >
          設定を保存
        </button>
      </div>
    </section>

    <section class="master-persona-panel" aria-labelledby="previewHeading">
      <div class="section-head">
        <div>
          <p class="eyebrow">入力ファイル</p>
          <h3 id="previewHeading">{viewModel.selectedFileName}</h3>
        </div>
        <div class="inline-actions compact-actions">
          <button
            class="button-secondary"
            id="chooseJsonButton"
            onclick={chooseJsonFile}
            type="button"
          >
            JSON を選ぶ
          </button>
          <button
            class="button-secondary"
            disabled={!viewModel.canStartPreview}
            id="resetJsonButton"
            onclick={resetJsonSelection}
            type="button"
          >
            選び直す
          </button>
        </div>
      </div>

      <input
        accept=".json,application/json"
        class="file-input"
        id="masterPersonaJsonInput"
        onchange={handleJsonSelected}
        type="file"
      />

      <div class="stats-grid" id="previewStats">
        <article class="stat-card">
          <span class="field-label">候補数</span>
          <strong>{viewModel.preview?.candidateCount ?? 0}</strong>
        </article>
        <article class="stat-card">
          <span class="field-label">新規追加可能</span>
          <strong>{viewModel.preview?.newlyAddableCount ?? 0}</strong>
        </article>
        <article class="stat-card">
          <span class="field-label">既存</span>
          <strong>{viewModel.preview?.existingCount ?? 0}</strong>
        </article>
      </div>

      <div class="status-row preview-actions">
        <div class="inline-actions compact-actions">
          <span class="status-pill">作成済みのペルソナはスキップされます</span>
          <span class="status-pill"
            >preview 状態: {viewModel.preview?.status ?? "入力待ち"}</span
          >
        </div>
        <div class="inline-actions compact-actions">
          <button
            class="button-secondary"
            disabled={!viewModel.canStartPreview}
            id="previewButton"
            onclick={() => void controller.previewGeneration()}
            type="button"
          >
            preview を更新
          </button>
          <button
            class="button-primary"
            disabled={!viewModel.canStartGeneration}
            id="executeGenerationButton"
            onclick={() => void controller.executeGeneration()}
            type="button"
          >
            この JSON で生成
          </button>
        </div>
      </div>
    </section>

    <section class="master-persona-panel run-panel" aria-labelledby="runHeading">
      <div class="section-head">
        <div>
          <p class="eyebrow">進行状況</p>
          <h3 id="runHeading">{viewModel.runStatus.message}</h3>
        </div>
        <span class="status-pill status-danger">{viewModel.detailLockText}</span>
      </div>

      <div class="progress-track">
        <div
          class="progress-fill"
          id="runProgressFill"
          style={`width: ${viewModel.progressPercent}%;`}
        ></div>
      </div>

      <div class="run-grid">
        <div class="run-card">
          <span class="field-label">完了件数</span>
          <strong>{viewModel.runStatus.processedCount}</strong>
        </div>
        <div class="run-card">
          <span class="field-label">いま処理中の NPC</span>
          <strong>{viewModel.runStatus.currentActorLabel || "-"}</strong>
        </div>
        <div class="run-card">
          <span class="field-label">現在の状態</span>
          <strong>{viewModel.runStatus.runState}</strong>
        </div>
      </div>

      <div class="status-row">
        <span class="status-pill">作成済み {viewModel.runStatus.successCount}</span>
        <span class="status-pill">既に作成済み {viewModel.runStatus.existingSkipCount}</span>
      </div>

      <div class="inline-actions run-actions">
        <span class="mini-text">一覧と詳細は見続けられます。</span>
        <div class="inline-actions compact-actions">
          <button
            class="button-secondary"
            disabled={!viewModel.isRunActive}
            id="interruptGenerationButton"
            onclick={() => void controller.interruptGeneration()}
            type="button"
          >
            一時停止
          </button>
          <button
            class="button-secondary"
            disabled={!viewModel.isRunActive}
            id="cancelGenerationButton"
            onclick={() => void controller.cancelGeneration()}
            type="button"
          >
            停止
          </button>
        </div>
      </div>
    </section>
  </section>

  <section class="workspace-grid">
    <section class="master-persona-panel" aria-labelledby="listHeading">
      <div class="section-head">
        <div>
          <p class="eyebrow">一覧</p>
          <h3 id="listHeading">ペルソナ一覧</h3>
          <p class="mini-text" id="pageStatusText">{viewModel.pageStatusText}</p>
        </div>
        <span class="status-pill">{viewModel.listHeadline}</span>
      </div>

      <div class="toolbar-grid">
        <label class="field-group" for="masterPersonaSearchInput">
          <span class="field-label">検索</span>
          <input
            class="search-field"
            id="masterPersonaSearchInput"
            oninput={(event) => controller.handleSearchInput(event)}
            placeholder="名前 / FormID / EditorID / 種族 / voice で検索"
            type="search"
            value={viewModel.keyword}
          />
        </label>
        <label class="field-group" for="masterPersonaPluginSelect">
          <span class="field-label">プラグイン</span>
          <select
            class="select-field"
            id="masterPersonaPluginSelect"
            onchange={(event) => controller.handlePluginFilterChange(event)}
            value={viewModel.pluginFilter}
          >
            {#each viewModel.pluginOptions as option (option.label)}
              <option value={option.value}>{option.label}</option>
            {/each}
          </select>
        </label>
      </div>

      <div class="column-row" aria-hidden="true">
        <span>NPC</span>
        <span>識別情報</span>
        <span>ペルソナ要約</span>
        <span>収録先</span>
      </div>

      <div class="list-stack" aria-live="polite">
        {#if viewModel.items.length === 0}
          <div class="empty-state">一致するペルソナがありません</div>
        {:else}
          {#each viewModel.items as item (item.identityKey)}
            <button
              class="list-row"
              class:is-selected={viewModel.selectedIdentityKey === item.identityKey}
              onclick={() => void controller.selectRow(item.identityKey)}
              type="button"
            >
              <div class="row-cell">
                <strong>{item.displayName}</strong>
                <span>{item.voiceType}</span>
              </div>
              <div class="row-cell">
                <span>{item.formId} / {item.editorId}</span>
                <span>クラス: {item.className || "-"}</span>
              </div>
              <div class="row-cell">
                <span>{item.personaSummary}</span>
                <span>{item.race ? item.race : ""}</span>
              </div>
              <div class="row-id">{item.targetPlugin}</div>
            </button>
          {/each}
        {/if}
      </div>

      <div class="pager-shell">
        <span class="mini-text" id="selectionStatusText">{viewModel.selectionStatusText}</span>
        <div class="inline-actions compact-actions">
          <button
            class="button-secondary"
            disabled={viewModel.page <= 1}
            id="prevPageButton"
            onclick={() => controller.goToPrevPage()}
            type="button"
          >
            前の30件
          </button>
          <button
            class="button-secondary"
            disabled={viewModel.page >= viewModel.totalPages}
            id="nextPageButton"
            onclick={() => controller.goToNextPage()}
            type="button"
          >
            次の30件
          </button>
        </div>
      </div>
    </section>

    <section class="master-persona-panel" aria-labelledby="detailHeading">
      <div class="section-head">
        <div>
          <p class="eyebrow">詳細</p>
          <h3 id="detailHeading">選択中のペルソナ</h3>
          <p class="mini-text" id="detailStatusText">{viewModel.detailStatusText}</p>
        </div>
        <div class="inline-actions compact-actions">
          <button
            class="button-secondary"
            disabled={!viewModel.canMutate}
            id="editButton"
            onclick={() => controller.openEditModal()}
            type="button"
          >
            更新
          </button>
          <button
            class="button-danger"
            disabled={!viewModel.canMutate}
            id="deleteButton"
            onclick={() => controller.openDeleteModal()}
            type="button"
          >
            削除
          </button>
        </div>
      </div>

      <div class="detail-title">
        <div class="status-row">
          {#if viewModel.selectedEntry}
            <span class="status-pill">{viewModel.selectedEntry.voiceType}</span>
          {/if}
        </div>
        <strong id="detailTitle"
          >{viewModel.selectedEntry?.displayName ?? "表示できるペルソナがありません"}</strong
        >
        <p class="mini-text" id="detailIdentityText">
          {#if viewModel.selectedEntry}
            FormID {viewModel.selectedEntry.formId} / EditorID {viewModel.selectedEntry.editorId} / {viewModel.selectedEntry.targetPlugin}
          {:else}
            検索条件を変更してください。
          {/if}
        </p>
      </div>

      <div class="detail-grid">
        <article class="detail-card">
          <span class="field-label">voice</span>
          <strong>{viewModel.selectedEntry?.voiceType || "-"}</strong>
        </article>
        <article class="detail-card">
          <span class="field-label">class</span>
          <strong>{viewModel.selectedEntry?.className || "-"}</strong>
        </article>
      </div>

      <dl class="detail-list">
        <div>
          <dt>名前</dt>
          <dd>{viewModel.selectedEntry?.displayName || "-"}</dd>
        </div>
        <div>
          <dt>source</dt>
          <dd>{viewModel.selectedEntry?.sourcePlugin || "-"}</dd>
        </div>
        <div>
          <dt>ペルソナ要約</dt>
          <dd>{viewModel.selectedEntry?.personaSummary || "-"}</dd>
        </div>
        <div>
          <dt>ペルソナ本文</dt>
          <dd>{viewModel.selectedEntry?.personaBody || "-"}</dd>
        </div>
      </dl>
    </section>
  </section>
</section>

<div
  aria-hidden={viewModel.modalState !== "edit"}
  class="modal-backdrop"
  class:is-open={viewModel.modalState === "edit"}
  hidden={viewModel.modalState !== "edit"}
  id="editModal"
  role="dialog"
>
  <section class="modal-card form-modal">
    <div class="section-head">
      <div>
        <p class="eyebrow">更新</p>
        <h3>ペルソナを編集</h3>
      </div>
      <button
        class="button-secondary"
        id="closeEditModalButton"
        onclick={() => controller.closeEditModal()}
        type="button"
      >
        閉じる
      </button>
    </div>

    <div class="form-grid">
      <label class="field-group textarea-group" for="editPersonaSummaryInput">
        <span class="field-label">ペルソナ概要</span>
        <textarea
          class="textarea-field"
          id="editPersonaSummaryInput"
          oninput={(event) => controller.setEditFormField("personaSummary", event)}
          value={viewModel.editForm.personaSummary ?? ""}
        ></textarea>
      </label>
      <label class="field-group" for="editSpeechStyleInput">
        <span class="field-label">話し方</span>
        <input
          class="text-field"
          id="editSpeechStyleInput"
          oninput={(event) => controller.setEditFormField("speechStyle", event)}
          value={viewModel.editForm.speechStyle ?? ""}
        />
      </label>
      <label class="field-group textarea-group" for="editPersonaBodyInput">
        <span class="field-label">ペルソナ本文</span>
        <textarea
          class="textarea-field"
          id="editPersonaBodyInput"
          oninput={(event) => controller.setEditFormField("personaBody", event)}
          value={viewModel.editForm.personaBody}
        ></textarea>
      </label>
    </div>

    <div class="inline-actions compact-actions">
      <button
        class="button-secondary"
        onclick={() => controller.closeEditModal()}
        type="button"
      >
        キャンセル
      </button>
      <button
        class="button-primary"
        id="saveEntryButton"
        onclick={() => void controller.saveCurrentEntry()}
        type="button"
      >
        更新する
      </button>
    </div>
  </section>
</div>

<div
  aria-hidden={viewModel.modalState !== "delete"}
  class="modal-backdrop"
  class:is-open={viewModel.modalState === "delete"}
  hidden={viewModel.modalState !== "delete"}
  id="deleteModal"
  role="dialog"
>
  <section class="modal-card">
    <div class="section-head">
      <div>
        <p class="eyebrow">削除</p>
        <h3>ペルソナを削除しますか</h3>
      </div>
      <button
        class="button-secondary"
        onclick={() => controller.closeDeleteModal()}
        type="button"
      >
        閉じる
      </button>
    </div>

    <dl class="detail-list">
      <div>
        <dt>名前</dt>
        <dd>{viewModel.selectedEntry?.displayName || "-"}</dd>
      </div>
      <div>
        <dt>FormID</dt>
        <dd>{viewModel.selectedEntry?.formId || "-"}</dd>
      </div>
      <div>
        <dt>EditorID</dt>
        <dd>{viewModel.selectedEntry?.editorId || "-"}</dd>
      </div>
    </dl>

    <div class="inline-actions compact-actions">
      <button
        class="button-secondary"
        onclick={() => controller.closeDeleteModal()}
        type="button"
      >
        キャンセル
      </button>
      <button
        class="button-danger"
        id="confirmDeleteButton"
        onclick={() => void controller.deleteCurrentEntry()}
        type="button"
      >
        削除する
      </button>
    </div>
  </section>
</div>

<style>
  .master-persona-shell {
    display: grid;
    gap: 18px;
  }

  .master-persona-panel {
    padding: 20px;
    border-radius: 20px;
    border: 0.5px solid var(--line);
    background: rgba(17, 13, 12, 0.42);
    box-shadow: var(--shadow);
    backdrop-filter: blur(24px);
  }

  .overview-panel,
  .generator-grid,
  .workspace-grid,
  .stats-grid,
  .run-grid,
  .detail-grid,
  .toolbar-grid,
  .form-grid {
    display: grid;
    gap: 14px;
  }

  .generator-grid {
    grid-template-columns: minmax(0, 0.8fr) minmax(0, 1.1fr) minmax(0, 0.9fr);
  }

  .workspace-grid {
    grid-template-columns: minmax(0, 1.3fr) minmax(300px, 0.85fr);
  }

  .stats-grid,
  .detail-grid,
  .run-grid {
    grid-template-columns: repeat(2, minmax(0, 1fr));
  }

  .run-grid {
    grid-template-columns: repeat(3, minmax(0, 1fr));
  }

  .hero-top,
  .section-head,
  .status-row,
  .inline-actions,
  .pager-shell,
  .column-row {
    display: flex;
    gap: 12px;
    align-items: center;
    justify-content: space-between;
    flex-wrap: wrap;
  }

  .eyebrow,
  .field-label,
  .column-row span,
  .detail-list dt,
  .mini-text {
    font-size: 12px;
    letter-spacing: 0.1em;
    text-transform: uppercase;
    color: var(--muted);
  }

  .lead,
  .prompt-copy,
  .detail-list dd,
  .row-cell span,
  .error-text {
    color: var(--muted);
    line-height: 1.6;
  }

  .text-field,
  .search-field,
  .select-field,
  .textarea-field {
    width: 100%;
    min-height: 42px;
    border-radius: 12px;
    border: 0.5px solid var(--line);
    background: rgba(255, 255, 255, 0.04);
    color: var(--text);
    padding: 0 14px;
  }

  .textarea-field {
    min-height: 140px;
    padding: 12px 14px;
    resize: vertical;
  }

  .button-primary,
  .button-secondary,
  .button-danger,
  .status-pill {
    min-height: 38px;
    padding: 0 14px;
    border-radius: 999px;
  }

  .button-primary,
  .button-secondary,
  .button-danger {
    border: 0.5px solid transparent;
    cursor: pointer;
  }

  .button-primary {
    color: #3f2400;
    background: linear-gradient(135deg, var(--primary) 0%, #f0a51f 100%);
  }

  .button-secondary {
    color: var(--text);
    background: rgba(255, 255, 255, 0.04);
    border-color: var(--line);
  }

  .button-danger {
    color: #35150d;
    background: linear-gradient(135deg, #ffc0ab 0%, #ff9c7c 100%);
  }

  .button-primary:disabled,
  .button-secondary:disabled,
  .button-danger:disabled {
    opacity: 0.45;
    cursor: not-allowed;
  }

  .status-pill {
    display: inline-flex;
    align-items: center;
    border: 0.5px solid var(--line);
    background: rgba(255, 255, 255, 0.03);
    color: var(--muted);
  }

  .status-accent {
    color: var(--bg-strong);
    border-color: transparent;
    background: linear-gradient(135deg, var(--primary) 0%, #f0a51f 100%);
  }

  .status-danger {
    color: var(--text);
    background: rgba(255, 156, 124, 0.14);
    border-color: rgba(255, 156, 124, 0.28);
  }

  .field-group,
  .detail-list,
  .detail-title,
  .run-panel,
  .row-cell,
  .pager-shell,
  .form-modal {
    display: grid;
    gap: 10px;
  }

  .textarea-group {
    grid-column: 1 / -1;
  }

  .file-input {
    display: none;
  }

  .stat-card,
  .detail-card,
  .run-card,
  .list-row {
    border-radius: 14px;
    border: 0.5px solid rgba(255, 186, 56, 0.12);
    background: rgba(255, 255, 255, 0.03);
    padding: 14px;
  }

  .stat-card strong,
  .detail-card strong,
  .run-card strong,
  .detail-title strong {
    display: block;
    font-size: 22px;
    overflow-wrap: anywhere;
  }

  .column-row {
    padding: 0 12px;
  }

  .list-stack {
    display: grid;
    gap: 8px;
  }

  .list-row {
    display: grid;
    gap: 10px;
    grid-template-columns: minmax(140px, 1fr) minmax(160px, 1fr) minmax(200px, 1.15fr) minmax(100px, 0.7fr);
    text-align: left;
    color: inherit;
  }

  .list-row.is-selected {
    background: rgba(255, 186, 56, 0.12);
    border-color: rgba(255, 186, 56, 0.28);
  }

  .row-id {
    text-align: right;
    color: var(--muted);
  }

  .detail-list {
    margin: 0;
  }

  .detail-list div {
    padding: 12px 14px;
    border-radius: 14px;
    background: rgba(255, 255, 255, 0.03);
  }

  .detail-list dt,
  .detail-list dd {
    margin: 0;
  }

  .progress-track {
    width: 100%;
    height: 10px;
    overflow: hidden;
    border-radius: 999px;
    background: rgba(255, 255, 255, 0.05);
    border: 0.5px solid var(--line);
  }

  .progress-fill {
    height: 100%;
    background: linear-gradient(135deg, var(--primary) 0%, #f0a51f 100%);
  }

  .modal-backdrop {
    position: fixed;
    inset: 0;
    display: none;
    align-items: center;
    justify-content: center;
    padding: 20px;
    background: rgba(14, 11, 10, 0.68);
    z-index: 12;
    visibility: hidden;
    pointer-events: none;
    opacity: 0;
  }

  .modal-backdrop.is-open {
    display: flex;
    visibility: visible;
    pointer-events: auto;
    opacity: 1;
  }

  .modal-backdrop[hidden] {
    display: none;
    visibility: hidden;
    pointer-events: none;
    opacity: 0;
  }

  .modal-card {
    width: min(860px, 100%);
    max-height: calc(100vh - 40px);
    overflow: auto;
    padding: 20px;
    border-radius: 20px;
    border: 0.5px solid var(--line);
    background: rgba(19, 15, 14, 0.94);
    box-shadow: var(--shadow);
    display: grid;
    gap: 14px;
  }

  .form-grid {
    grid-template-columns: repeat(2, minmax(0, 1fr));
  }

  .empty-state {
    padding: 20px;
    border-radius: 16px;
    background: rgba(255, 255, 255, 0.03);
    color: var(--muted);
  }

  @media (max-width: 1220px) {
    .generator-grid,
    .workspace-grid,
    .stats-grid,
    .run-grid,
    .detail-grid,
    .form-grid {
      grid-template-columns: 1fr;
    }
  }

  @media (max-width: 900px) {
    .list-row {
      grid-template-columns: 1fr;
    }

    .row-id {
      text-align: left;
    }

    .column-row {
      display: none;
    }
  }
</style>
