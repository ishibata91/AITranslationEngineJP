<script lang="ts">
  import type { ModelSettingsCardViewModel } from "@application/gateway-contract/model-settings-card"
  import AIModelSelectionCard from "@ui/components/AIModelSelectionCard.svelte"
  import type { GenerationSetupPanelProps } from "./master-persona-panel-props"

  let {
    aiSettings,
    aiSettingsStatusText,
    aiSettingsWarningText,
    aiProviderLabel,
    canSelectModel,
    canStartGeneration,
    executionMethodOptions,
    isRunActive,
    modelOptions,
    modelSettingsCardViewModel,
    preview,
    selectedFileName,
    selectedFileReference,
    isAISettingsRefreshing,
    handleJsonSelected,
    chooseJsonFile,
    resetJsonSelection,
    handleAIProviderChange,
    handleAIModelChange,
    handleAIExecutionMethodChange,
    refreshAISettings,
    startGeneration,
    saveAISettings
  }: GenerationSetupPanelProps = $props()

  const isFileSelected = $derived(selectedFileReference !== null)
  const fileStatusText = $derived(
    isRunActive ? "生成中" : isFileSelected ? "JSON 選択済み" : "入力待ち"
  )
  const helperText = $derived(
    preview
      ? `候補 ${preview.candidateCount} 件のうち、新規 ${preview.newlyAddableCount} 件を作成できます。`
      : "ベースゲームや大型 Mod の対象 NPC から、未作成のペルソナだけを追加します。"
  )

  function createFallbackModelCard(): ModelSettingsCardViewModel {
    return {
      referenceId: "master-persona",
      provider: aiSettings.provider,
      model: aiSettings.model,
      providerOptions: [
        { value: "gemini", label: "Gemini" },
        { value: "lm_studio", label: "LM Studio" },
        { value: "xai", label: "xAI" }
      ],
      credentialStatusLabel: aiSettingsStatusText,
      credentialStatusTone: aiSettingsWarningText ? "warning" : "success",
      showCredentialStatus: true,
      showCredentialWarning: aiSettingsWarningText !== "",
      credentialWarningText:
        "APIキーが未設定のため、モデル一覧を更新できません。",
      modelListButtonEnabled: true,
      modelListButtonLabel: "モデル一覧を更新",
      modelListButtonAriaLabel: "モデル一覧を更新",
      isModelListRefreshing: isAISettingsRefreshing,
      modelListStatusText: canSelectModel
        ? "使うモデルを選んでください。"
        : "モデル一覧を更新してください。",
      modelOptions,
      modelSelectEnabled: canSelectModel,
      emptyModelLabel: canSelectModel ? "選んでください" : "設定が必要",
      statusLabel: aiProviderLabel || "未選択",
      statusTone: aiSettingsWarningText ? "warning" : "neutral",
      helperText: "ペルソナ作成に使う AI サービスとモデルを選びます。",
      footerMessage: "",
      footerWarningText: "",
      actionButtonDisabled: isAISettingsRefreshing
    }
  }

  const modelCard = $derived(
    modelSettingsCardViewModel ?? createFallbackModelCard()
  )
</script>

<section
  class="panel setup-panel"
  aria-labelledby="setupHeading"
  data-testid="master-persona-generation-setup-panel"
>
  <div class="section-head">
    <div>
      <p class="eyebrow">生成準備</p>
      <h2 id="setupHeading">入力ファイルと AI 設定</h2>
    </div>
    <span class="status-pill">{fileStatusText}</span>
  </div>

  <div class="setup-grid">
    <div class="model-panel" data-testid="master-persona-ai-settings-card">
      <AIModelSelectionCard
        actionButtonDisabled={modelCard.actionButtonDisabled ||
          isAISettingsRefreshing}
        actionButtonId="saveAiSettingsButton"
        actionButtonLabel="設定を保存"
        credentialStatusLabel={modelCard.credentialStatusLabel}
        credentialStatusTone={modelCard.credentialStatusTone}
        credentialWarningText={modelCard.credentialWarningText}
        eyebrow="AI 設定"
        executionDisabled={isAISettingsRefreshing}
        executionOptions={executionMethodOptions}
        executionSelectId="executionMethodSelect"
        executionValue={aiSettings.executionMethod}
        footerMessage={modelCard.footerMessage ||
          "モデル設定はこの画面専用です。必要なら保存できます。"}
        footerWarningText={modelCard.footerWarningText}
        helperText={modelCard.helperText}
        modelDisabled={isAISettingsRefreshing || !modelCard.modelSelectEnabled}
        modelOptions={modelCard.modelOptions}
        modelSelectId="modelInput"
        modelStatusText={modelCard.modelListStatusText}
        modelValue={modelCard.model}
        onAction={saveAISettings}
        onExecutionChange={handleAIExecutionMethodChange}
        onModelChange={handleAIModelChange}
        onProviderChange={handleAIProviderChange}
        onRefresh={() => void refreshAISettings()}
        providerFieldLabel="AI サービス"
        providerOptions={modelCard.providerOptions}
        providerDisabled={isAISettingsRefreshing}
        providerSelectId="providerSelect"
        providerValue={modelCard.provider}
        refreshButtonAriaLabel={modelCard.modelListButtonAriaLabel}
        refreshButtonLabel={modelCard.modelListButtonLabel}
        refreshDisabled={!modelCard.modelListButtonEnabled ||
          isAISettingsRefreshing}
        refreshSpinning={modelCard.isModelListRefreshing ||
          isAISettingsRefreshing}
        secondaryControlMode="execution-select"
        showCredentialStatus={modelCard.showCredentialStatus}
        showCredentialWarning={modelCard.showCredentialWarning}
        statusLabel={modelCard.statusLabel}
        statusTone={modelCard.statusTone}
        title="この画面で使う AI 設定"
        titleId="settingsHeading"
        titleTag="h3"
        emptyModelLabel={modelCard.emptyModelLabel}
      />

      <div class="model-footer">
        <p class="support-copy">
          {aiSettingsWarningText ||
            "モデル設定はこの画面専用です。必要なら保存できます。"}
        </p>
        <p class="sr-only">
          プロンプトテンプレートは画面入力では変更せず、実装側の説明文として固定しています。
        </p>
      </div>
    </div>

    <section
      class="file-panel"
      aria-labelledby="fileHeading"
      data-testid="master-persona-input-json-panel"
    >
      <div class="section-head">
        <div>
          <p class="eyebrow">入力 JSON</p>
          <h3 id="fileHeading">{selectedFileName}</h3>
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
          <strong>{preview?.candidateCount ?? 0}</strong>
        </article>
        <article class="stat-card">
          <span>新規作成数</span>
          <strong>{preview?.newlyAddableCount ?? 0}</strong>
        </article>
        <article class="stat-card">
          <span>既存スキップ数</span>
          <strong>{preview?.existingCount ?? 0}</strong>
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
          disabled={!canStartGeneration}
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
    color: var(--text);
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
