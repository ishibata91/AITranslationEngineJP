<script lang="ts">
  type ProviderId = "gemini" | "xai" | "lm_studio"
  type SaveState = "saved" | "dirty" | "saving" | "failed"
  type ValidationState = "not_checked" | "checking" | "ready" | "failed" | "stale"
  type SecretState = "present" | "missing" | "not_required"
  type PrototypePageId = "provider_settings" | "model_cards"
  type ModelCardState = "ready" | "needs_settings" | "disabled"
  type ModelProviderLabel = "Gemini" | "xAI" | "LM Studio"

  interface ProviderSetting {
    id: ProviderId
    label: string
    description: string
    endpoint: string
    savedEndpoint: string
    apiKeyRequired: boolean
    secretState: SecretState
    saveState: SaveState
    validationState: ValidationState
    updatedAt: string
    message: string
  }

  interface PrototypePage {
    id: PrototypePageId
    label: string
    description: string
  }

  interface ModelCard {
    id: string
    title: string
    provider: ModelProviderLabel
    model: string
    batchLabel: string
    state: ModelCardState
  }

  const providerOrder: ProviderId[] = ["gemini", "xai", "lm_studio"]
  const modelProviderOptions: ModelProviderLabel[] = ["Gemini", "xAI", "LM Studio"]
  const modelOptionsByProvider: Record<ModelProviderLabel, string[]> = {
    Gemini: ["gemini-2.5-flash", "gemini-2.5-pro", "gemini-2.0-flash"],
    xAI: ["grok-3-mini", "grok-3", "grok-2-vision"],
    "LM Studio": ["local-qwen3-14b", "local-llama-3.1-8b", "local-mistral-nemo"]
  }
  const prototypePages: PrototypePage[] = [
    {
      id: "provider_settings",
      label: "AIサービス設定",
      description: "エンドポイントと APIキー状態を AIサービスごとに管理します。"
    },
    {
      id: "model_cards",
      label: "モデルカード確認",
      description: "他ページで表示するモデル選択カードを確認します。"
    }
  ]

  let activePrototypePageId = $state<PrototypePageId>("provider_settings")
  let selectedProviderId = $state<ProviderId>("gemini")
  let keyDraft = $state("")
  let showKeyPanel = $state(false)
  let showSavedToast = $state(false)
  let refreshingModelCardIds = $state<string[]>([])
  let activeBatchTooltipId = $state<string | null>(null)

  let modelCards = $state<ModelCard[]>([
    {
      id: "term-translation",
      title: "用語翻訳モデル",
      provider: "Gemini",
      model: "gemini-2.5-flash",
      batchLabel: "Batch API 使用",
      state: "ready"
    },
    {
      id: "body-translation",
      title: "本文翻訳モデル",
      provider: "xAI",
      model: "grok-3-mini",
      batchLabel: "通常 API",
      state: "needs_settings"
    },
    {
      id: "persona-generation",
      title: "ペルソナ生成モデル",
      provider: "LM Studio",
      model: "local-qwen3-14b",
      batchLabel: "通常 API",
      state: "ready"
    }
  ])

  let providers = $state<ProviderSetting[]>([
    {
      id: "gemini",
      label: "Gemini",
      description: "公式エンドポイント",
      endpoint: "https://generativelanguage.googleapis.com",
      savedEndpoint: "https://generativelanguage.googleapis.com",
      apiKeyRequired: true,
      secretState: "present",
      saveState: "saved",
      validationState: "ready",
      updatedAt: "2026/05/04 10:20",
      message: "保存済み"
    },
    {
      id: "xai",
      label: "xAI",
      description: "APIキー設定後に接続確認可能",
      endpoint: "https://api.x.ai/v1",
      savedEndpoint: "https://api.x.ai/v1",
      apiKeyRequired: true,
      secretState: "missing",
      saveState: "dirty",
      validationState: "not_checked",
      updatedAt: "未保存",
      message: "APIキー設定後に接続確認"
    },
    {
      id: "lm_studio",
      label: "LM Studio",
      description: "ローカルエンドポイント。APIキー不要",
      endpoint: "http://127.0.0.1:1234/v1",
      savedEndpoint: "http://127.0.0.1:1234/v1",
      apiKeyRequired: false,
      secretState: "not_required",
      saveState: "saved",
      validationState: "ready",
      updatedAt: "2026/05/04 09:45",
      message: "ローカルエンドポイント設定済み"
    }
  ])

  const selectedProvider = $derived(
    providers.find((provider) => provider.id === selectedProviderId) ?? providers[0]
  )

  const readyCount = $derived(
    providers.filter((provider) => provider.validationState === "ready").length
  )

  const activePrototypePage = $derived(
    prototypePages.find((page) => page.id === activePrototypePageId) ?? prototypePages[0]
  )

  function selectPrototypePage(pageId: PrototypePageId): void {
    activePrototypePageId = pageId
  }

  function selectProvider(providerId: ProviderId): void {
    selectedProviderId = providerId
    showKeyPanel = false
    keyDraft = ""
  }

  function providerById(providerId: ProviderId): ProviderSetting {
    return providers.find((provider) => provider.id === providerId) ?? providers[0]
  }

  function statusLabel(provider: ProviderSetting): string {
    if (provider.saveState === "saving") {
      return "保存中"
    }

    if (provider.saveState === "failed") {
      return "保存失敗"
    }

    if (provider.validationState === "stale") {
      return "再確認が必要"
    }

    if (provider.validationState === "ready") {
      return "利用可能"
    }

    if (provider.secretState === "missing") {
      return "APIキー未設定"
    }

    return "未確認"
  }

  function secretLabel(provider: ProviderSetting): string {
    if (provider.secretState === "not_required") {
      return "不要"
    }

    if (provider.secretState === "present") {
      return "保存済み"
    }

    return "未保存"
  }

  function validationLabel(provider: ProviderSetting): string {
    if (provider.validationState === "ready") {
      return "確認済み"
    }

    if (provider.validationState === "checking") {
      return "接続確認中"
    }

    if (provider.validationState === "failed") {
      return "接続失敗。エンドポイントを確認。"
    }

    if (provider.validationState === "stale") {
      return "設定変更後の再確認が必要。"
    }

    return "未確認"
  }

  function displayEndpoint(value: string): string {
    return value.trim() ? value : "未設定"
  }

  function markDirty(provider: ProviderSetting): void {
    provider.saveState = "dirty"
    provider.validationState = "stale"
    provider.message = "エンドポイント変更済み。接続確認が必要。"
    showSavedToast = false
  }

  function updateEndpoint(provider: ProviderSetting, value: string): void {
    provider.endpoint = value
    markDirty(provider)
  }

  function saveApiKey(provider: ProviderSetting): void {
    if (!keyDraft.trim()) {
      return
    }

    provider.secretState = "present"
    provider.saveState = "dirty"
    provider.validationState = "stale"
    provider.message = "APIキー設定済み。接続確認が必要。"
    keyDraft = ""
    showKeyPanel = false
  }

  function validateConnection(provider: ProviderSetting): void {
    if (provider.apiKeyRequired && provider.secretState !== "present") {
      provider.message = "APIキー設定後に接続確認。"
      provider.validationState = "not_checked"
      return
    }

    if (!provider.endpoint.trim()) {
      provider.validationState = "failed"
      provider.message = "エンドポイント未入力。"
      return
    }

    provider.validationState = "checking"
    provider.message = "接続確認中。"

    window.setTimeout(() => {
      if (provider.id === "xai" && provider.endpoint.includes("invalid")) {
        provider.validationState = "failed"
        provider.message = "接続失敗。APIキー状態は保存済み。"
        return
      }

      provider.validationState = "ready"
      provider.message = "確認済み。"
    }, 420)
  }

  function saveProvider(provider: ProviderSetting): void {
    if (provider.validationState !== "ready") {
      provider.saveState = "failed"
      provider.message = "接続確認後に保存。"
      return
    }

    provider.saveState = "saving"
    provider.message = "保存中。"

    window.setTimeout(() => {
      provider.savedEndpoint = provider.endpoint
      provider.saveState = "saved"
      provider.updatedAt = "2026/05/04 11:05"
      provider.message = "保存済み。"
      showSavedToast = true
    }, 420)
  }

  function resetProvider(provider: ProviderSetting): void {
    provider.endpoint = ""
    provider.savedEndpoint = ""
    provider.secretState = provider.apiKeyRequired ? "missing" : "not_required"
    provider.saveState = "dirty"
    provider.validationState = "not_checked"
    provider.updatedAt = "未保存"
    provider.message =
      "未設定化。エンドポイントと APIキー状態を未設定。secret 本体は削除。"
  }

  function providerIdForModelProvider(providerLabel: ModelProviderLabel): ProviderId {
    if (providerLabel === "xAI") {
      return "xai"
    }

    if (providerLabel === "LM Studio") {
      return "lm_studio"
    }

    return "gemini"
  }

  function providerSettingForModelProvider(providerLabel: ModelProviderLabel): ProviderSetting {
    return providerById(providerIdForModelProvider(providerLabel))
  }

  function canFetchModels(providerLabel: ModelProviderLabel): boolean {
    const provider = providerSettingForModelProvider(providerLabel)

    if (!provider.endpoint.trim()) {
      return false
    }

    if (provider.apiKeyRequired && provider.secretState !== "present") {
      return false
    }

    return true
  }

  function modelFetchWarning(providerLabel: ModelProviderLabel): string {
    const provider = providerSettingForModelProvider(providerLabel)

    if (!provider.endpoint.trim()) {
      return "AIサービス設定が未完了です。エンドポイントを設定してください。"
    }

    if (provider.apiKeyRequired && provider.secretState !== "present") {
      return "AIサービス設定が未完了です。"
    }

    return ""
  }

  function modelOptionsForCard(card: ModelCard): string[] {
    if (!canFetchModels(card.provider)) {
      return [card.model]
    }

    return modelOptionsByProvider[card.provider]
  }

  function updateModelCardProvider(cardId: string, provider: ModelProviderLabel): void {
    const options = canFetchModels(provider) ? modelOptionsByProvider[provider] : []
    updateModelCard(cardId, {
      provider,
      model: options[0] ?? "",
      state: canFetchModels(provider) ? "ready" : "needs_settings"
    })
  }

  function refreshModels(card: ModelCard): void {
    if (!canFetchModels(card.provider) || refreshingModelCardIds.includes(card.id)) {
      return
    }

    refreshingModelCardIds = [...refreshingModelCardIds, card.id]

    window.setTimeout(() => {
      refreshingModelCardIds = refreshingModelCardIds.filter((id) => id !== card.id)
    }, 520)
  }

  function updateModelCard(cardId: string, patch: Partial<ModelCard>): void {
    modelCards = modelCards.map((card) => {
      if (card.id !== cardId) {
        return card
      }

      const nextProvider = patch.provider ?? card.provider
      const nextCanFetchModels = canFetchModels(nextProvider)

      return {
        ...card,
        ...patch,
        state: patch.state ?? (nextCanFetchModels ? "ready" : "needs_settings")
      }
    })
    showSavedToast = false
  }

  function modelCardStatusLabel(card: ModelCard): string {
    if (card.state === "ready") {
      return "利用可能"
    }

    if (card.state === "needs_settings") {
      return "設定が必要"
    }

    return "使用不可"
  }
</script>

<svelte:head>
  <title>AIサービス設定 UI プロトタイプ</title>
</svelte:head>

<main class="shell" data-ui-prototype-sample-data-root>
  <header class="shell-bar">
    <div class="brand">
      <span class="brand-eyebrow">AITranslationEngineJp</span>
      <strong>翻訳エンジン</strong>
    </div>
    <nav aria-label="グローバルナビゲーション" class="global-nav">
      {#each prototypePages as page}
        <button
          aria-current={activePrototypePageId === page.id ? "page" : undefined}
          class:is-active={activePrototypePageId === page.id}
          onclick={() => selectPrototypePage(page.id)}
          type="button"
        >
          {page.label}
        </button>
      {/each}
    </nav>
  </header>

  <section class="page">
    {#if activePrototypePageId === "provider_settings"}
      <section class="panel hero-panel">
        <div class="hero-title-row">
          <div>
            <p class="page-label">現在のページ</p>
            <h1>{activePrototypePage.label}</h1>
          </div>
          <span class="status-pill">{readyCount} / 3 利用可能</span>
        </div>
        <p class="hero-lead">
          {activePrototypePage.description}
        </p>
      </section>

      <section class="layout-grid">
        <aside class="panel provider-list" aria-label="AIサービス一覧">
          <div class="section-head">
            <div>
              <p class="page-label">AIサービス</p>
              <h2>保存状態</h2>
            </div>
          </div>
          <div class="provider-stack">
            {#each providerOrder as providerId}
              {@const provider = providerById(providerId)}
              <button
                aria-pressed={provider.id === selectedProviderId}
                class:is-active={provider.id === selectedProviderId}
                class="provider-row"
                onclick={() => selectProvider(provider.id)}
                type="button"
              >
                <span>
                  <strong>{provider.label}</strong>
                  <small>{statusLabel(provider)}</small>
                </span>
                <span class:ready={provider.validationState === "ready"} class="dot"></span>
              </button>
            {/each}
          </div>
        </aside>

      <section class="panel settings-panel" aria-labelledby="providerSettingsHeading">
        <div class="section-head">
          <div>
            <p class="page-label">{selectedProvider.label}</p>
            <h2 id="providerSettingsHeading">接続設定</h2>
          </div>
          <span class="status-pill">{statusLabel(selectedProvider)}</span>
        </div>

        <div class="setting-grid">
          <label class="field-block" for="endpointInput">
            <span>エンドポイント</span>
            <input
              id="endpointInput"
              oninput={(event) =>
                updateEndpoint(selectedProvider, event.currentTarget.value)}
              value={selectedProvider.endpoint}
            />
            <small>保存済み: {displayEndpoint(selectedProvider.savedEndpoint)}</small>
          </label>

          <div class="field-block">
            <span>APIキー状態</span>
            <div class="key-state">
              <strong>{secretLabel(selectedProvider)}</strong>
              {#if selectedProvider.apiKeyRequired}
                <button
                  class="button-secondary"
                  onclick={() => (showKeyPanel = !showKeyPanel)}
                  type="button"
                >
                  設定
                </button>
              {/if}
            </div>
          </div>

          {#if showKeyPanel && selectedProvider.apiKeyRequired}
            <section class="key-panel" aria-label="APIキー設定">
              <label class="field-block" for="apiKeyInput">
                <span>APIキー</span>
                <input
                  id="apiKeyInput"
                  oninput={(event) => (keyDraft = event.currentTarget.value)}
                  placeholder="保存後非表示"
                  type="password"
                  value={keyDraft}
                />
              </label>
              <p>
                APIキーは OS標準の資格情報マネージャーに保存します。
              </p>
              <div class="button-row">
                <button
                  class="button-primary"
                  disabled={!keyDraft.trim()}
                  onclick={() => saveApiKey(selectedProvider)}
                  type="button"
                >
                  保存
                </button>
                <button
                  class="button-secondary"
                  onclick={() => {
                    showKeyPanel = false
                    keyDraft = ""
                  }}
                  type="button"
                >
                  キャンセル
                </button>
              </div>
            </section>
          {/if}

          <div class="field-block">
            <span>接続確認</span>
            <div class="validation-box">
              <p>{validationLabel(selectedProvider)}</p>
              <button
                class="button-secondary"
                disabled={selectedProvider.validationState === "checking"}
                onclick={() => validateConnection(selectedProvider)}
                type="button"
              >
                接続を確認
              </button>
            </div>
          </div>
        </div>

        <div class="summary-band">
          <div class="button-row">
            <button
              class="button-secondary"
              onclick={() => resetProvider(selectedProvider)}
              type="button"
            >
              リセット
            </button>
            <button
              class="button-primary"
              disabled={selectedProvider.saveState === "saving"}
              onclick={() => saveProvider(selectedProvider)}
              type="button"
            >
              設定を保存
            </button>
          </div>
        </div>
      </section>
      </section>
    {:else}
      <div class="model-settings-grid" aria-label="参照側ページのモデル設定">
        {#each modelCards as card}
          <section class="model-setting" aria-label={`${card.title}のモデル設定`}>
            <div class="model-setting-head">
              <h2>{card.title}</h2>
              <span
                class={card.state === "needs_settings"
                  ? "status-pill status-alert"
                  : "status-pill"}
              >
                {modelCardStatusLabel(card)}
              </span>
            </div>

            <div class="model-edit-grid" aria-label={`${card.title}の設定`}>
              <label class="model-field" for={`${card.id}-provider`}>
                <span>AIサービス</span>
                <select
                  id={`${card.id}-provider`}
                  onchange={(event) =>
                    updateModelCardProvider(
                      card.id,
                      event.currentTarget.value as ModelProviderLabel
                    )}
                >
                  {#each modelProviderOptions as option}
                    <option selected={option === card.provider} value={option}>
                      {option}
                    </option>
                  {/each}
                </select>
              </label>

              <div class="model-field">
                <span id={`${card.id}-model-label`}>モデル</span>
                <div class="model-select-row">
                  <select
                    aria-labelledby={`${card.id}-model-label`}
                    disabled={!canFetchModels(card.provider)}
                    id={`${card.id}-model`}
                    onchange={(event) =>
                      updateModelCard(card.id, { model: event.currentTarget.value })}
                  >
                    {#each modelOptionsForCard(card) as modelOption}
                      <option selected={modelOption === card.model} value={modelOption}>
                        {modelOption || "未取得"}
                      </option>
                    {/each}
                  </select>
                  <button
                    aria-label={`${card.title}のモデル一覧を更新`}
                    class="icon-button"
                    disabled={!canFetchModels(card.provider)}
                    onclick={() => refreshModels(card)}
                    type="button"
                  >
                    <span
                      class="refresh-icon"
                      class:is-spinning={refreshingModelCardIds.includes(card.id)}
                      aria-hidden="true"
                    >
                      ↻
                    </span>
                  </button>
                </div>
                {#if !canFetchModels(card.provider)}
                  <small class="warning-text">{modelFetchWarning(card.provider)}</small>
                {/if}
              </div>

              <div class="model-field">
                <span class="field-title-row">
                  処理方式
                  <button
                    aria-label="Batch API の説明"
                    class="help-tooltip"
                    onblur={() => (activeBatchTooltipId = null)}
                    onfocus={() => (activeBatchTooltipId = card.id)}
                    onmouseenter={() => (activeBatchTooltipId = card.id)}
                    onmouseleave={() => (activeBatchTooltipId = null)}
                    type="button"
                  >
                    ?
                    <span
                      class={activeBatchTooltipId === card.id
                        ? "tooltip-panel is-visible"
                        : "tooltip-panel"}
                      role="tooltip"
                    >
                      Batch API は API利用料が安くなる場合があります。
                    </span>
                  </button>
                </span>
                <label class="checkbox-row" for={`${card.id}-batch`}>
                  <input
                    checked={card.batchLabel === "Batch API 使用"}
                    id={`${card.id}-batch`}
                    onchange={(event) =>
                      updateModelCard(card.id, {
                        batchLabel: event.currentTarget.checked ? "Batch API 使用" : "通常 API"
                      })}
                    type="checkbox"
                  />
                  <span>Batch API</span>
                </label>
              </div>
            </div>
          </section>
        {/each}
      </div>
    {/if}

  </section>

  {#if showSavedToast}
    <aside class="toast" role="status">設定保存済み</aside>
  {/if}
</main>

<style>
  :global(body) {
    margin: 0;
  }

  :global(button),
  :global(input),
  :global(select) {
    font: inherit;
  }

  :global(a) {
    color: inherit;
    text-decoration: none;
  }

  :global(h1),
  :global(h2),
  :global(h3),
  :global(dl),
  :global(p) {
    margin: 0;
  }

  .shell {
    --surface: rgba(34, 30, 29, 0.84);
    --surface-strong: rgba(22, 19, 18, 0.94);
    --line: rgba(255, 186, 56, 0.22);
    --line-strong: rgba(255, 186, 56, 0.38);
    --text: #f3e9e4;
    --muted: #d8c8bc;
    --primary: #ffba38;
    --accent: #8fd5c4;
    color: var(--text);
    min-height: 100vh;
  }

  .shell-bar {
    align-items: center;
    background: rgba(22, 19, 17, 0.88);
    border-bottom: 1px solid var(--line);
    display: flex;
    gap: 24px;
    justify-content: space-between;
    padding: 16px clamp(18px, 4vw, 48px);
    position: sticky;
    top: 0;
    z-index: 10;
  }

  .brand {
    display: grid;
    gap: 3px;
    min-width: 180px;
  }

  .brand-eyebrow,
  .page-label {
    color: var(--primary);
    font-size: 0.72rem;
    letter-spacing: 0;
  }

  .global-nav {
    display: flex;
    flex-wrap: wrap;
    gap: 8px;
    justify-content: flex-end;
  }

  .global-nav button {
    background: transparent;
    border: 1px solid transparent;
    border-radius: 999px;
    color: var(--muted);
    cursor: pointer;
    font: inherit;
    padding: 8px 12px;
  }

  .global-nav button.is-active {
    background: rgba(255, 186, 56, 0.12);
    border-color: var(--line-strong);
    color: var(--text);
  }

  .page {
    display: grid;
    gap: 18px;
    margin: 0 auto;
    max-width: 1240px;
    padding: 24px clamp(16px, 3vw, 36px) 60px;
  }

  .panel {
    background: var(--surface);
    border: 1px solid var(--line);
    border-radius: 8px;
    box-shadow: 0 22px 60px rgba(0, 0, 0, 0.28);
  }

  .hero-panel,
  .settings-panel,
  .provider-list {
    padding: clamp(18px, 3vw, 26px);
  }

  .hero-title-row,
  .section-head,
  .summary-band {
    align-items: flex-start;
    display: flex;
    gap: 16px;
    justify-content: space-between;
  }

  h1 {
    font-size: clamp(2rem, 4vw, 3.3rem);
    font-weight: 700;
    line-height: 1.08;
  }

  h2 {
    font-size: 1.22rem;
    line-height: 1.2;
  }

  h3 {
    font-size: 1.05rem;
    line-height: 1.25;
  }

  .hero-lead {
    color: var(--muted);
    line-height: 1.8;
    margin-top: 12px;
    max-width: 780px;
  }

  .status-pill {
    border: 1px solid var(--line-strong);
    border-radius: 999px;
    color: var(--accent);
    flex: none;
    padding: 7px 10px;
    white-space: nowrap;
  }

  .status-pill.status-alert {
    background: rgba(255, 87, 87, 0.12);
    border-color: rgba(255, 87, 87, 0.56);
    color: #ff9a9a;
  }

  .layout-grid {
    align-items: start;
    display: grid;
    gap: 18px;
    grid-template-columns: minmax(240px, 0.34fr) minmax(0, 1fr);
  }

  .provider-stack {
    display: grid;
    gap: 10px;
    margin-top: 16px;
  }

  .provider-row {
    align-items: center;
    background: rgba(255, 255, 255, 0.04);
    border: 1px solid var(--line);
    border-radius: 8px;
    color: var(--text);
    cursor: pointer;
    display: flex;
    justify-content: space-between;
    min-height: 66px;
    padding: 12px;
    text-align: left;
  }

  .provider-row.is-active {
    background: rgba(255, 186, 56, 0.14);
    border-color: var(--line-strong);
  }

  .provider-row span:first-child {
    display: grid;
    gap: 4px;
  }

  .provider-row small,
  .field-block small {
    color: var(--muted);
    line-height: 1.5;
  }

  .dot {
    background: #8c7563;
    border-radius: 50%;
    height: 10px;
    width: 10px;
  }

  .dot.ready {
    background: var(--accent);
  }

  .setting-grid {
    display: grid;
    gap: 14px;
    grid-template-columns: repeat(2, minmax(0, 1fr));
    margin-top: 18px;
  }

  .field-block {
    background: rgba(0, 0, 0, 0.14);
    border: 1px solid rgba(255, 255, 255, 0.08);
    border-radius: 8px;
    display: grid;
    gap: 8px;
    min-width: 0;
    padding: 14px;
  }

  .field-block > span {
    color: var(--primary);
    font-size: 0.86rem;
  }

  input,
  select {
    background: rgba(255, 255, 255, 0.08);
    border: 1px solid var(--line);
    border-radius: 6px;
    color: var(--text);
    min-height: 42px;
    min-width: 0;
    padding: 9px 10px;
    width: 100%;
  }

  .key-state,
  .validation-box,
  .button-row {
    align-items: center;
    display: flex;
    flex-wrap: wrap;
    gap: 10px;
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

  .key-panel {
    background: rgba(143, 213, 196, 0.08);
    border: 1px solid rgba(143, 213, 196, 0.28);
    border-radius: 8px;
    display: grid;
    gap: 12px;
    grid-column: 1 / -1;
    padding: 14px;
  }

  .key-panel p {
    color: var(--muted);
    line-height: 1.6;
  }

  .summary-band {
    background: var(--surface-strong);
    border: 1px solid var(--line);
    border-radius: 8px;
    margin-top: 18px;
    padding: 14px;
  }

  .model-settings-grid {
    display: grid;
    gap: 14px;
    grid-template-columns: repeat(3, minmax(0, 1fr));
  }

  .model-setting {
    background: rgba(255, 255, 255, 0.05);
    border: 1px solid rgba(255, 255, 255, 0.1);
    border-radius: 8px;
    display: grid;
    gap: 14px;
    grid-template-rows: auto 1fr;
    min-width: 0;
    padding: 16px;
  }

  .model-setting-head {
    align-items: flex-start;
    display: flex;
    gap: 12px;
    justify-content: space-between;
  }

  .model-setting-head h2 {
    font-size: 1rem;
    line-height: 1.4;
    margin: 0;
  }

  .model-edit-grid {
    display: grid;
    gap: 8px;
  }

  .model-field {
    background: rgba(0, 0, 0, 0.14);
    border: 1px solid rgba(255, 255, 255, 0.08);
    border-radius: 8px;
    display: grid;
    gap: 6px;
    min-width: 0;
    padding: 10px;
  }

  .model-field span {
    color: var(--primary);
    font-size: 0.75rem;
  }

  .field-title-row {
    align-items: center;
    display: inline-flex;
    gap: 6px;
  }

  .model-field select {
    min-height: 36px;
    padding: 7px 9px;
  }

  .checkbox-row {
    align-items: center;
    color: var(--text);
    display: inline-flex;
    gap: 8px;
    min-height: 36px;
  }

  .checkbox-row input {
    accent-color: var(--primary);
    height: 16px;
    margin: 0;
    width: 16px;
  }

  .checkbox-row span {
    color: var(--text);
    font-size: 0.9rem;
  }

  .help-tooltip {
    align-items: center;
    background: transparent;
    border: 1px solid var(--line);
    border-radius: 50%;
    color: var(--muted);
    cursor: help;
    display: inline-flex;
    font-size: 0.72rem;
    height: 18px;
    justify-content: center;
    position: relative;
    width: 18px;
  }

  .tooltip-panel {
    background: rgba(22, 19, 18, 0.98);
    border: 1px solid var(--line-strong);
    border-radius: 6px;
    bottom: calc(100% + 8px);
    color: var(--text);
    display: none;
    font-size: 0.8rem;
    left: 50%;
    line-height: 1.5;
    min-width: 176px;
    padding: 8px 10px;
    position: absolute;
    transform: translateX(-50%);
    z-index: 30;
  }

  .help-tooltip:hover .tooltip-panel,
  .help-tooltip:focus .tooltip-panel,
  .tooltip-panel.is-visible {
    display: block;
  }

  .model-select-row {
    align-items: center;
    display: grid;
    gap: 8px;
    grid-template-columns: minmax(0, 1fr) 36px;
  }

  .icon-button {
    align-items: center;
    background: rgba(255, 255, 255, 0.06);
    border: 1px solid var(--line);
    border-radius: 6px;
    color: var(--text);
    cursor: pointer;
    display: inline-flex;
    font-size: 1rem;
    height: 36px;
    justify-content: center;
    line-height: 1;
    padding: 0;
    width: 36px;
  }

  .icon-button .refresh-icon {
    display: inline-block;
    font-size: 1.35rem;
    line-height: 1;
  }

  .icon-button .refresh-icon.is-spinning {
    animation: spin 0.7s linear infinite;
  }

  .warning-text {
    color: #ffd38a;
    line-height: 1.5;
  }

  @keyframes spin {
    to {
      transform: rotate(360deg);
    }
  }

  .toast {
    background: rgba(20, 44, 38, 0.96);
    border: 1px solid rgba(143, 213, 196, 0.42);
    border-radius: 8px;
    bottom: 18px;
    color: var(--text);
    padding: 12px 14px;
    position: fixed;
    right: 18px;
    z-index: 20;
  }

  @media (max-width: 860px) {
    .shell-bar,
    .hero-title-row,
    .section-head,
    .summary-band {
      align-items: stretch;
      flex-direction: column;
    }

    .layout-grid,
    .model-settings-grid,
    .setting-grid {
      grid-template-columns: 1fr;
    }

    .global-nav {
      justify-content: flex-start;
    }

    .status-pill {
      width: fit-content;
    }
  }
</style>
