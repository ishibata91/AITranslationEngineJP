<script lang="ts">
  import type { MasterPersonaScreenViewModel } from "@application/gateway-contract/master-persona"
  import AIModelSelectionCard from "@ui/components/AIModelSelectionCard.svelte"

  interface Props {
    viewModel: MasterPersonaScreenViewModel
    isAISettingsRefreshing: boolean
    handleAIProviderChange: (event: Event) => void
    handleAIModelChange: (event: Event) => void
    handleAIExecutionMethodChange: (event: Event) => void
    refreshAISettings: () => Promise<void>
    handleJsonSelected: (event: Event) => void
    chooseJsonFile: () => void
    resetJsonSelection: () => void
    startGeneration: () => void
    saveAISettings: () => void
  }

  let {
    viewModel,
    isAISettingsRefreshing,
    handleAIProviderChange,
    handleAIModelChange,
    handleAIExecutionMethodChange,
    refreshAISettings,
    handleJsonSelected,
    chooseJsonFile,
    resetJsonSelection,
    startGeneration,
    saveAISettings
  }: Props = $props()

  const isFileSelected = $derived(viewModel.selectedFileReference !== null)
  const fileStatusText = $derived(
    viewModel.isRunActive
      ? "生成中"
      : isFileSelected
        ? "JSON 選択済み"
        : "入力待ち"
  )
  const helperText = $derived(
    viewModel.preview
      ? `候補 ${viewModel.preview.candidateCount} 件のうち、新規 ${viewModel.preview.newlyAddableCount} 件を作成できます。`
      : "ベースゲームや大型 Mod の対象 NPC から、未作成のペルソナだけを追加します。"
  )
</script>

<section class="panel setup-panel" aria-labelledby="setupHeading">
  <div class="section-head">
    <div>
      <p class="eyebrow">生成準備</p>
      <h2 id="setupHeading">入力ファイルと AI 設定</h2>
    </div>
    <span class="status-pill">{fileStatusText}</span>
  </div>

  <div class="setup-grid">
    <div class="model-panel">
      <AIModelSelectionCard
        credentialStatusLabel={viewModel.aiSettingsStatusText}
        credentialStatusTone={viewModel.aiSettingsWarningText ? "warning" : "success"}
        eyebrow="AI 設定"
        executionDisabled={isAISettingsRefreshing}
        executionOptions={viewModel.executionMethodOptions}
        executionSelectId="executionMethodSelect"
        executionValue={viewModel.aiSettings.executionMethod}
        helperText="ペルソナ作成に使う AI サービスとモデルを選びます。"
        modelDisabled={isAISettingsRefreshing || !viewModel.canSelectModel}
        modelOptions={viewModel.modelOptions}
        modelSelectId="modelInput"
        modelStatusText={viewModel.canSelectModel
          ? "使うモデルを選んでください。"
          : "モデル一覧を更新してください。"}
        modelValue={viewModel.aiSettings.model}
        onExecutionChange={handleAIExecutionMethodChange}
        onModelChange={handleAIModelChange}
        onProviderChange={handleAIProviderChange}
        onRefresh={() => void refreshAISettings()}
        providerFieldLabel="AI サービス"
        providerOptions={[
          { value: "gemini", label: "Gemini" },
          { value: "lm_studio", label: "LM Studio" },
          { value: "xai", label: "xAI" }
        ]}
        providerDisabled={isAISettingsRefreshing}
        providerSelectId="providerSelect"
        providerValue={viewModel.aiSettings.provider}
        refreshButtonAriaLabel="モデル一覧を更新"
        refreshButtonLabel="モデル一覧を更新"
        refreshDisabled={isAISettingsRefreshing}
        refreshSpinning={isAISettingsRefreshing}
        secondaryControlMode="execution-select"
        showCredentialStatus={true}
        showCredentialWarning={viewModel.aiSettingsWarningText !== ""}
        statusLabel={viewModel.aiProviderLabel || "未選択"}
        statusTone={viewModel.aiSettingsWarningText ? "warning" : "neutral"}
        title="この画面で使う AI 設定"
        titleId="settingsHeading"
        titleTag="h3"
        emptyModelLabel={viewModel.canSelectModel ? "選んでください" : "設定が必要"}
      />

      <div class="model-footer">
        <p class="support-copy">
          {viewModel.aiSettingsWarningText || "モデル設定はこの画面専用です。必要なら保存できます。"}
        </p>
        <p class="sr-only">
          プロンプトテンプレートは画面入力では変更せず、実装側の説明文として固定しています。
        </p>
        <button
          class="button-secondary"
          disabled={isAISettingsRefreshing}
          id="saveAiSettingsButton"
          onclick={saveAISettings}
          type="button"
        >
          設定を保存
        </button>
      </div>
    </div>

    <section class="file-panel" aria-labelledby="fileHeading">
      <div class="section-head">
        <div>
          <p class="eyebrow">入力 JSON</p>
          <h3 id="fileHeading">{viewModel.selectedFileName}</h3>
        </div>
      </div>

      <input
        accept=".json,application/json"
        class="file-input"
        id="masterPersonaJsonInput"
        onchange={handleJsonSelected}
        type="file"
      />

      <p class="support-copy">{helperText}</p>

      <div class="stats-grid" aria-label="対象件数" id="previewStats">
        <article class="stat-card">
          <span>候補数</span>
          <strong>{viewModel.preview?.candidateCount ?? 0}</strong>
        </article>
        <article class="stat-card">
          <span>新規作成数</span>
          <strong>{viewModel.preview?.newlyAddableCount ?? 0}</strong>
        </article>
        <article class="stat-card">
          <span>既存スキップ数</span>
          <strong>{viewModel.preview?.existingCount ?? 0}</strong>
        </article>
      </div>

      <div class="button-row">
        <button
          class="button-secondary"
          id="chooseJsonButton"
          onclick={chooseJsonFile}
          type="button"
        >
          JSON を選択
        </button>
        <button
          class="button-secondary"
          disabled={!isFileSelected}
          id="resetJsonButton"
          onclick={resetJsonSelection}
          type="button"
        >
          選び直す
        </button>
        <button
          class="button-primary button-wide"
          disabled={!viewModel.canStartGeneration}
          id="executeGenerationButton"
          onclick={startGeneration}
          type="button"
        >
          ペルソナを作成
        </button>
      </div>
    </section>
  </div>
</section>

<style>
  .panel,
  .file-panel,
  .model-footer,
  .stat-card {
    border-radius: 20px;
  }

  .panel {
    background: rgba(17, 13, 12, 0.42);
    border: 0.5px solid var(--line);
    box-shadow: var(--shadow);
    min-width: 0;
    padding: clamp(18px, 3vw, 24px);
  }

  .setup-panel,
  .model-panel,
  .file-panel {
    display: grid;
    gap: 14px;
    min-width: 0;
  }

  .setup-grid {
    display: grid;
    gap: 16px;
    grid-template-columns: minmax(0, 1.02fr) minmax(0, 0.98fr);
    margin-top: 14px;
  }

  .section-head,
  .button-row,
  .model-footer {
    align-items: flex-start;
    display: flex;
    flex-wrap: wrap;
    gap: 12px;
    justify-content: space-between;
  }

  .eyebrow,
  .stat-card span {
    color: var(--muted);
    font-size: 12px;
    letter-spacing: 0.1em;
    margin: 0;
    text-transform: uppercase;
  }

  h2,
  h3,
  p {
    margin: 0;
  }

  h2,
  h3,
  strong,
  .support-copy {
    overflow-wrap: anywhere;
  }

  .status-pill {
    align-items: center;
    background: rgba(255, 255, 255, 0.04);
    border: 0.5px solid rgba(255, 186, 56, 0.22);
    border-radius: 999px;
    color: var(--text);
    display: inline-flex;
    min-height: 38px;
    padding: 0 14px;
  }

  .model-footer,
  .file-panel {
    background: rgba(255, 255, 255, 0.03);
    border: 0.5px solid rgba(255, 255, 255, 0.08);
    padding: 16px;
  }

  .support-copy {
    color: var(--muted);
    line-height: 1.7;
  }

  .stats-grid {
    display: grid;
    gap: 10px;
    grid-template-columns: repeat(3, minmax(0, 1fr));
  }

  .stat-card {
    background: rgba(0, 0, 0, 0.16);
    border: 0.5px solid rgba(255, 255, 255, 0.08);
    display: grid;
    gap: 6px;
    min-width: 0;
    padding: 14px;
  }

  .stat-card strong {
    font-size: clamp(1.1rem, 1.8vw, 1.45rem);
  }

  .button-row {
    justify-content: flex-start;
  }

  .button-primary,
  .button-secondary {
    border-radius: 999px;
    cursor: pointer;
    min-height: 40px;
    min-width: 0;
    overflow-wrap: anywhere;
    padding: 0 16px;
  }

  .button-primary {
    background: linear-gradient(135deg, var(--primary) 0%, #f0a51f 100%);
    border: 0.5px solid transparent;
    color: #3f2400;
  }

  .button-secondary {
    background: rgba(255, 255, 255, 0.04);
    border: 0.5px solid var(--line);
    color: var(--text);
  }

  .button-wide {
    margin-left: auto;
  }

  .sr-only {
    border: 0;
    clip: rect(0 0 0 0);
    clip-path: inset(50%);
    height: 1px;
    margin: -1px;
    overflow: hidden;
    padding: 0;
    position: absolute;
    white-space: nowrap;
    width: 1px;
  }

  button:disabled {
    cursor: not-allowed;
    opacity: 0.5;
  }

  .file-input {
    display: none;
  }

  @media (max-width: 980px) {
    .setup-grid {
      grid-template-columns: 1fr;
    }
  }

  @media (max-width: 560px) {
    .stats-grid {
      grid-template-columns: 1fr;
    }

    .button-row > * {
      width: 100%;
    }

    .button-wide {
      margin-left: 0;
    }
  }
</style>
