<script lang="ts">
  import AIModelSelectionCard from "/@fs/Users/iorishibata/Repositories/AITranslationEngineJP/frontend/src/ui/components/AIModelSelectionCard.svelte"

  type SelectOption = { value: string; label: string }
  type ModelOption = { modelId: string; label: string }
  type PreviewState = "empty" | "ready" | "running" | "complete" | "error"

  interface Props {
    aiProviders: SelectOption[]
    executionMethods: SelectOption[]
    modelOptions: ModelOption[]
    provider: string
    model: string
    executionMethod: string
    previewState: PreviewState
    changeProviderAction: (event: Event) => void
    changeModelAction: (event: Event) => void
    changeExecutionAction: (event: Event) => void
    refreshModelsAction: () => void
    chooseFileAction: () => void
    startGenerationAction: () => void
  }

  let {
    aiProviders,
    executionMethods,
    modelOptions,
    provider,
    model,
    executionMethod,
    previewState,
    changeProviderAction,
    changeModelAction,
    changeExecutionAction,
    refreshModelsAction,
    chooseFileAction,
    startGenerationAction,
  }: Props = $props()

  const fileSelected = $derived(previewState !== "empty" && previewState !== "error")
  const canGenerate = $derived(previewState === "ready" || previewState === "complete")
</script>

<section class="panel setup-panel" aria-labelledby="setupHeading">
  <div class="section-head">
    <div>
      <p class="eyebrow">生成準備</p>
      <h2 id="setupHeading">入力ファイルと AI 設定</h2>
    </div>
    <span class="status-pill">{fileSelected ? "JSON 選択済み" : "入力待ち"}</span>
  </div>

  <div class="setup-grid">
    <AIModelSelectionCard
      credentialStatusLabel="APIキー保存済み"
      credentialStatusTone="success"
      eyebrow="AI 設定"
      executionOptions={executionMethods}
      executionSelectId="personaExecutionMethod"
      executionValue={executionMethod}
      helperText="ペルソナ生成に使う AI サービスとモデルを選びます。"
      modelOptions={modelOptions}
      modelSelectId="personaModel"
      modelStatusText="取得済みのモデルから選びます。"
      modelValue={model}
      onExecutionChange={changeExecutionAction}
      onModelChange={changeModelAction}
      onProviderChange={changeProviderAction}
      onRefresh={refreshModelsAction}
      providerFieldLabel="AI サービス"
      providerOptions={aiProviders}
      providerSelectId="personaProvider"
      providerValue={provider}
      refreshButtonAriaLabel="モデル一覧を更新"
      refreshButtonLabel="モデル一覧を更新"
      secondaryControlMode="execution-select"
      showCredentialStatus={true}
      statusLabel="利用可能"
      statusTone="success"
      title="ペルソナ生成モデル"
      titleId="personaModelCardHeading"
      titleTag="h3"
    />

    <div class="file-card">
      <div>
        <p class="eyebrow">入力 JSON</p>
        <h3>JSON を選択</h3>
        <p class="support-copy">
          ベースゲームや大型 Mod の対象 NPC から、未作成のペルソナだけを追加します。
        </p>
      </div>

      <div class="file-summary">
        <span>選択中</span>
        <strong>{previewState === "empty" ? "未選択" : "whiterun_dialogue.json"}</strong>
      </div>

      <div class="stat-grid">
        <article>
          <span>候補</span>
          <strong>{fileSelected ? 42 : 0}</strong>
        </article>
        <article>
          <span>新規</span>
          <strong>{fileSelected ? 18 : 0}</strong>
        </article>
        <article>
          <span>既存</span>
          <strong>{fileSelected ? 24 : 0}</strong>
        </article>
      </div>

      {#if previewState === "error"}
        <p class="error-text">JSON を読み取れませんでした。ファイル形式を確認してください。</p>
      {:else if previewState === "complete"}
        <p class="success-text">生成が完了しました。下の一覧と詳細で確認できます。</p>
      {:else}
        <p class="support-copy">
          既に作成済みのペルソナは生成対象から外します。
        </p>
      {/if}

      <div class="button-row">
        <button
          class="button-secondary"
          onclick={chooseFileAction}
          type="button"
        >
          JSON を選択
        </button>
        <button
          class="button-primary"
          disabled={!canGenerate}
          onclick={startGenerationAction}
          type="button"
        >
          ペルソナを作成
        </button>
      </div>
    </div>
  </div>
</section>

<style>
  .panel {
    background: var(--surface);
    border: 1px solid var(--line);
    border-radius: 8px;
    box-shadow: var(--shadow);
    min-width: 0;
  }

  .setup-panel {
    padding: clamp(16px, 2vw, 22px);
  }

  .section-head,
  .button-row {
    align-items: flex-start;
    display: flex;
    flex-wrap: wrap;
    gap: 12px;
    justify-content: space-between;
  }

  .setup-grid {
    display: grid;
    gap: 14px;
    grid-template-columns: minmax(280px, 0.92fr) minmax(0, 1.08fr);
    margin-top: 14px;
  }

  .file-card {
    background: rgba(255, 255, 255, 0.05);
    border: 1px solid rgba(255, 255, 255, 0.1);
    border-radius: 8px;
    display: grid;
    gap: 14px;
    min-width: 0;
    padding: 16px;
  }

  .file-summary,
  .stat-grid article {
    background: rgba(0, 0, 0, 0.14);
    border: 1px solid rgba(255, 255, 255, 0.08);
    border-radius: 8px;
    display: grid;
    gap: 6px;
    padding: 12px;
  }

  .stat-grid {
    display: grid;
    gap: 8px;
    grid-template-columns: repeat(3, minmax(0, 1fr));
  }

  .eyebrow,
  .file-summary span,
  .stat-grid span {
    color: var(--primary);
    font-size: 0.78rem;
  }

  .support-copy,
  .error-text,
  .success-text {
    color: var(--muted);
    line-height: 1.6;
  }

  .error-text {
    color: #ffb1a1;
  }

  .success-text {
    color: #a9e8d4;
  }

  .status-pill {
    border: 1px solid var(--line-strong);
    border-radius: 999px;
    color: var(--accent);
    flex: none;
    padding: 7px 10px;
    white-space: nowrap;
  }

  .button-primary,
  .button-secondary {
    border-radius: 7px;
    cursor: pointer;
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

  button:disabled {
    cursor: not-allowed;
    opacity: 0.55;
  }

  h2,
  h3,
  p {
    margin: 0;
  }

  @media (max-width: 980px) {
    .setup-grid,
    .stat-grid {
      grid-template-columns: 1fr;
    }
  }
</style>
