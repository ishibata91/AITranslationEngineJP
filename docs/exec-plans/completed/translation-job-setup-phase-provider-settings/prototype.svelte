<script lang="ts">
  type PhaseId = "terms" | "persona" | "body"
  type ServiceId = "gemini" | "xai" | "openai" | "lmstudio"
  type ModelState = "not_loaded" | "loading" | "ready" | "blocked" | "failed"

  interface AiService {
    id: ServiceId
    label: string
    apiKey: "required" | "not_required"
    apiKeyReady: boolean
    batchAvailable: boolean
    models: string[]
  }

  interface PhaseSetting {
    id: PhaseId
    name: string
    purpose: string
    serviceId: ServiceId
    model: string
    modelState: ModelState
    batch: boolean
    message: string
  }

  let services = $state<AiService[]>([
    {
      id: "gemini",
      label: "Gemini",
      apiKey: "required",
      apiKeyReady: true,
      batchAvailable: true,
      models: ["gemini-2.5-pro", "gemini-2.5-flash"]
    },
    {
      id: "xai",
      label: "xAI",
      apiKey: "required",
      apiKeyReady: false,
      batchAvailable: true,
      models: ["grok-3", "grok-3-mini"]
    },
    {
      id: "openai",
      label: "OpenAI",
      apiKey: "required",
      apiKeyReady: true,
      batchAvailable: false,
      models: ["gpt-5.1", "gpt-5.1-mini"]
    },
    {
      id: "lmstudio",
      label: "LM Studio",
      apiKey: "not_required",
      apiKeyReady: true,
      batchAvailable: false,
      models: ["local-skyrim-ja", "local-fast-draft"]
    }
  ])

  let phases = $state<PhaseSetting[]>([
    {
      id: "terms",
      name: "単語翻訳",
      purpose: "用語と固有名詞を先にそろえる段階です。",
      serviceId: "gemini",
      model: "gemini-2.5-flash",
      modelState: "ready",
      batch: true,
      message: "作成に必要な設定がそろっています。"
    },
    {
      id: "persona",
      name: "NPC ペルソナ生成",
      purpose: "登場人物の話し方を翻訳前に整える段階です。",
      serviceId: "lmstudio",
      model: "local-skyrim-ja",
      modelState: "ready",
      batch: false,
      message: "作成に必要な設定がそろっています。"
    },
    {
      id: "body",
      name: "本文翻訳",
      purpose: "本文をまとめて翻訳する段階です。",
      serviceId: "xai",
      model: "",
      modelState: "blocked",
      batch: false,
      message: "APIキーを登録してからモデル一覧を更新してください。"
    }
  ])

  let created = $state(false)
  let apiKeyModalServiceId = $state<ServiceId | null>(null)
  let apiKeyDraft = $state("")
  let apiKeySaveState = $state<"idle" | "saving" | "saved">("idle")

  function serviceFor(id: ServiceId): AiService {
    return services.find((service) => service.id === id) ?? services[0]
  }

  function modelStateLabel(phase: PhaseSetting): string {
    if (phase.modelState === "loading") {
      return "モデル一覧を更新しています。"
    }

    if (phase.modelState === "blocked") {
      return "モデル一覧を更新できません。APIキーの設定が必要です。"
    }

    if (phase.modelState === "failed") {
      return "モデル一覧を取得できませんでした。時間をおいて再実行してください。"
    }

    if (phase.modelState === "ready") {
      return "取得済みのモデルから選べます。"
    }

    return "モデル一覧を更新してください。"
  }

  function onServiceChange(phase: PhaseSetting, serviceId: ServiceId): void {
    const nextService = serviceFor(serviceId)
    phase.serviceId = serviceId
    phase.model = ""
    phase.batch = false

    if (nextService.apiKey === "required" && !nextService.apiKeyReady) {
      phase.modelState = "blocked"
      phase.message = "APIキーを登録してからモデル一覧を更新してください。"
      return
    }

    phase.modelState = "not_loaded"
    phase.message = "モデル一覧を更新して、使うモデルを選んでください。"
  }

  function refreshModels(phase: PhaseSetting): void {
    const service = serviceFor(phase.serviceId)

    if (service.apiKey === "required" && !service.apiKeyReady) {
      phase.modelState = "blocked"
      phase.message = "APIキーを登録してからモデル一覧を更新してください。"
      return
    }

    phase.modelState = "loading"
    phase.message = "モデル一覧を更新しています。"

    window.setTimeout(() => {
      if (phase.id === "body" && phase.serviceId === "gemini") {
        phase.modelState = "failed"
        phase.model = ""
        phase.message = "モデル一覧を取得できませんでした。もう一度更新してください。"
        return
      }

      phase.modelState = "ready"
      phase.model = service.models[0] ?? ""
      phase.message = phase.model
        ? "作成に必要な設定がそろっています。"
        : "使うモデルを選んでください。"
    }, 450)
  }

  function onModelChange(phase: PhaseSetting, model: string): void {
    phase.model = model
    phase.message = model
      ? "作成に必要な設定がそろっています。"
      : "使うモデルを選んでください。"
  }

  function onBatchChange(phase: PhaseSetting, batch: boolean): void {
    phase.batch = batch
    phase.message = "作成に必要な設定がそろっています。"
  }

  function phaseReady(phase: PhaseSetting): boolean {
    const service = serviceFor(phase.serviceId)
    return (
      phase.modelState === "ready" &&
      Boolean(phase.model) &&
      (service.apiKey === "not_required" || service.apiKeyReady)
    )
  }

  function canCreate(): boolean {
    return phases.every((phase) => phaseReady(phase))
  }

  function readyCount(): number {
    return phases.filter((phase) => phaseReady(phase)).length
  }

  function blockedReasons(): string[] {
    return phases.flatMap((phase) => {
      const service = serviceFor(phase.serviceId)
      const reasons: string[] = []

      if (service.apiKey === "required" && !service.apiKeyReady) {
        reasons.push(`${phase.name}のAPIキーが未設定です`)
      }

      if (!phase.model || phase.modelState !== "ready") {
        reasons.push(`${phase.name}のモデルを選んでください`)
      }

      return reasons
    })
  }

  function createJob(): void {
    if (canCreate()) {
      created = true
    }
  }

  function resetCreate(): void {
    created = false
  }

  function openApiKeyModal(service: AiService): void {
    apiKeyModalServiceId = service.id
    apiKeyDraft = ""
    apiKeySaveState = "idle"
  }

  function closeApiKeyModal(): void {
    apiKeyModalServiceId = null
    apiKeyDraft = ""
    apiKeySaveState = "idle"
  }

  function saveApiKey(): void {
    if (!apiKeyModalServiceId || !apiKeyDraft.trim()) {
      return
    }

    apiKeySaveState = "saving"

    window.setTimeout(() => {
      const service = serviceFor(apiKeyModalServiceId as ServiceId)
      service.apiKeyReady = true
      apiKeySaveState = "saved"

      for (const phase of phases) {
        if (phase.serviceId === service.id && phase.modelState === "blocked") {
          phase.modelState = "not_loaded"
          phase.message = "APIキーを登録しました。モデル一覧を更新してください。"
        }
      }

      apiKeyDraft = ""
    }, 450)
  }
</script>

<svelte:head>
  <title>翻訳ジョブ作成 UI プロトタイプ</title>
</svelte:head>

<section
  class="job-setup-shell"
  data-ui-prototype-sample-data-root
  id="translationJobSetupView"
>
  <section class="job-setup-card hero-card">
    <div class="hero-head">
      <div>
        <p class="eyebrow">翻訳ジョブ作成</p>
        <h2>翻訳段階ごとの AI 設定</h2>
      </div>
    </div>
    <p class="lead">
      入力データ、共通辞書、共通ペルソナを確認し、この画面で選んだ
      AIサービスを設定します。
    </p>
  </section>

  {#if created}
    <section class="summary-grid">
      <section class="job-setup-card" aria-labelledby="jobSetupSummaryHeading">
        <div class="section-head">
          <div>
            <p class="eyebrow">作成結果</p>
            <h3 id="jobSetupSummaryHeading">作成後の設定内容</h3>
          </div>
          <span class="status-pill success">作成済み</span>
        </div>
        <dl class="detail-grid">
          <div>
            <dt>ジョブID</dt>
            <dd>job-local-2026-05-04-001</dd>
          </div>
          <div>
            <dt>入力データ</dt>
            <dd class="wrap-value">Skyrim_dialogue_export_ja_source.json</dd>
          </div>
          <div>
            <dt>APIキー表示</dt>
            <dd>文字列は表示しません</dd>
          </div>
        </dl>
        <button class="button-secondary" onclick={resetCreate} type="button">
          編集へ戻る
        </button>
      </section>

      <section class="job-setup-card" aria-labelledby="jobSetupPhaseSummaryHeading">
        <div class="section-head">
          <div>
            <p class="eyebrow">翻訳段階別</p>
            <h3 id="jobSetupPhaseSummaryHeading">保存された設定</h3>
          </div>
        </div>
        <dl class="detail-grid">
          {#each phases as phase (phase.id)}
            {@const service = serviceFor(phase.serviceId)}
            <div>
              <dt>{phase.name}</dt>
              <dd class="wrap-value">
                {service.label} / {phase.model} / {phase.batch ? "一括処理を使う" : "一括処理を使わない"}
              </dd>
            </div>
          {/each}
        </dl>
      </section>
    </section>
  {:else}
    <section class="content-grid">
      <section class="job-setup-card compact-card" aria-labelledby="jobSetupInputHeading">
        <div class="section-head">
          <div>
            <p class="eyebrow">1 入力</p>
            <h3 id="jobSetupInputHeading">入力データ</h3>
          </div>
          <span class="mini-text">準備済み</span>
        </div>
        <label class="field-block" for="jobSetupInputSelect">
          <span>入力データ</span>
          <select id="jobSetupInputSelect">
            <option>Skyrim_dialogue_export_ja_source.json</option>
            <option>Dawnguard_books_source.json</option>
          </select>
        </label>
        <dl class="detail-grid compact">
          <div>
            <dt>入力データ名</dt>
            <dd class="wrap-value">Skyrim_dialogue_export_ja_source.json</dd>
          </div>
          <div>
            <dt>出自</dt>
            <dd class="wrap-value">翻訳入力レビューで確認済み</dd>
          </div>
          <div>
            <dt>登録日時</dt>
            <dd>2026/05/04 10:30</dd>
          </div>
          <div>
            <dt>翻訳レコード件数</dt>
            <dd>1,248 件</dd>
          </div>
          <div>
            <dt>既存ジョブ状態</dt>
            <dd class="wrap-value">同じ入力から作成された未完了ジョブはありません。</dd>
          </div>
        </dl>
      </section>

      <section class="job-setup-card compact-card" aria-labelledby="jobSetupFoundationHeading">
        <div class="section-head">
          <div>
            <p class="eyebrow">2 共通基盤</p>
            <h3 id="jobSetupFoundationHeading">共通辞書と共通ペルソナ</h3>
          </div>
        </div>
        <div class="foundation-tables">
          <div>
            <p class="mini-label">共通辞書</p>
            <div class="foundation-table-wrap">
              <table class="foundation-table" aria-label="共通辞書一覧">
                <thead>
                  <tr>
                    <th scope="col">名称</th>
                  </tr>
                </thead>
                <tbody>
                  <tr><td>Skyrim 基本語辞書</td></tr>
                  <tr><td>地名と人名の表記統一</td></tr>
                  <tr><td>派閥名の統一</td></tr>
                  <tr><td>魔法効果語彙</td></tr>
                  <tr><td>装備カテゴリ表記</td></tr>
                  <tr><td>クエスト用語</td></tr>
                  <tr><td>種族名の統一</td></tr>
                  <tr><td>祝福と呪い</td></tr>
                </tbody>
              </table>
            </div>
          </div>
          <div>
            <p class="mini-label">共通ペルソナ</p>
            <div class="foundation-table-wrap">
              <table class="foundation-table" aria-label="共通ペルソナ一覧">
                <thead>
                  <tr>
                    <th scope="col">名称</th>
                  </tr>
                </thead>
                <tbody>
                  <tr><td>丁寧な説明口調</td></tr>
                  <tr><td>ファンタジー語彙優先</td></tr>
                  <tr><td>住民会話は自然体</td></tr>
                  <tr><td>固有名詞は原文尊重</td></tr>
                  <tr><td>古風な役職名</td></tr>
                  <tr><td>説明文は簡潔</td></tr>
                  <tr><td>商人は軽い敬語</td></tr>
                  <tr><td>衛兵は短く断定</td></tr>
                </tbody>
              </table>
            </div>
          </div>
        </div>
        <p class="mini-text">翻訳ジョブで参照する共通基盤を確認します。</p>
      </section>

      <section
        class="job-setup-card"
        aria-labelledby="jobSetupPhaseSettingsHeading"
      >
        <div class="section-head">
          <div>
            <p class="eyebrow">3 翻訳段階別設定</p>
            <h3 id="jobSetupPhaseSettingsHeading">AI サービスとモデル</h3>
          </div>
          <span class={canCreate() ? "status-pill success" : "status-pill"}>
            {readyCount()}/3 設定済み
          </span>
        </div>
        <p class="mini-text">
          各翻訳段階で、AIサービス、モデル、実行方法の順に確認します。
        </p>
        <div class="phase-settings-list">
          {#each phases as phase, index (phase.id)}
            {@const service = serviceFor(phase.serviceId)}
            <section
              aria-labelledby={`${phase.id}-phase-heading`}
              class="phase-setting-block"
            >
              <div class="phase-setting-head">
                <div>
                  <p class="mini-label">段階 {index + 1}</p>
                  <h4 id={`${phase.id}-phase-heading`}>{phase.name}</h4>
                  <p class="mini-text">{phase.purpose}</p>
                </div>
                <span
                  class:success={phaseReady(phase)}
                  class:warning={!phaseReady(phase)}
                  class="tag"
                >
                  {phaseReady(phase) ? "設定済み" : "未設定あり"}
                </span>
              </div>
              <div class="phase-control-group compact">
                <div>
                  <p class="mini-label">AIサービス</p>
                </div>
                <label class="field-block" for={`${phase.id}-service`}>
                  <span>AIサービス</span>
                  <select
                    id={`${phase.id}-service`}
                    onchange={(event) => {
                      const target = event.currentTarget
                      if (target instanceof HTMLSelectElement) {
                        onServiceChange(phase, target.value as ServiceId)
                      }
                    }}
                    value={phase.serviceId}
                  >
                    {#each services as option (option.id)}
                      <option value={option.id}>{option.label}</option>
                    {/each}
                  </select>
                </label>
                <p class="status-copy">
                  <strong>APIキー状態</strong>
                  <span>
                    {service.apiKey === "not_required"
                      ? "不要"
                      : service.apiKeyReady
                        ? "設定済み"
                        : "未設定"}
                  </span>
                </p>
                {#if service.apiKey === "required" && !service.apiKeyReady}
                  <button
                    class="button-secondary"
                    onclick={() => openApiKeyModal(service)}
                    type="button"
                  >
                    APIキーを登録
                  </button>
                {/if}
              </div>
              <div class="phase-control-group model-control">
                <div>
                  <p class="mini-label">モデル</p>
                  <p class="mini-text">モデル一覧を更新して、取得済みの候補から選びます。</p>
                </div>
                <button
                  class="button-secondary"
                  disabled={phase.modelState === "loading" ||
                    (service.apiKey === "required" && !service.apiKeyReady)}
                  id={`${phase.id}-refresh-models`}
                  onclick={() => refreshModels(phase)}
                  type="button"
                >
                  モデル一覧を更新
                </button>
                <p class="mini-text">{modelStateLabel(phase)}</p>
                {#if phase.modelState === "ready"}
                  <label class="field-block" for={`${phase.id}-model`}>
                    <span>使うモデル</span>
                    <select
                      id={`${phase.id}-model`}
                      onchange={(event) => {
                        const target = event.currentTarget
                        if (target instanceof HTMLSelectElement) {
                          onModelChange(phase, target.value)
                        }
                      }}
                      value={phase.model}
                    >
                      {#each service.models as model (model)}
                        <option value={model}>{model}</option>
                      {/each}
                    </select>
                  </label>
                {/if}
              </div>
              <div class="phase-control-group compact">
                <div>
                  <p class="mini-label">実行方法</p>
                </div>
                {#if service.batchAvailable}
                  <div class="field-block">
                    <span>
                      一括処理
                      <span class="help-tip">
                        <button
                          aria-describedby={`${phase.id}-batch-help-control`}
                          aria-label="一括処理の説明"
                          class="help-mark"
                          type="button"
                        >
                          ?
                        </button>
                        <span
                          class="help-tooltip"
                          id={`${phase.id}-batch-help-control`}
                          role="tooltip"
                        >
                          バッチAPIを使うと安く済ませられることがあります。
                        </span>
                      </span>
                    </span>
                    <label class="mini-text">
                      <input
                        checked={phase.batch}
                        onchange={(event) => {
                          const target = event.currentTarget
                          if (target instanceof HTMLInputElement) {
                            onBatchChange(phase, target.checked)
                          }
                        }}
                        type="checkbox"
                      />
                      一括処理で実行する
                    </label>
                  </div>
                {:else}
                  <p class="empty-text">
                    この AI サービスでは一括処理の切り替えはありません。
                  </p>
                {/if}
              </div>
              <p class="phase-message">{phase.message}</p>
            </section>
          {/each}
        </div>
      </section>

    </section>
    <div class="global-next-bar" role="region" aria-label="作成前確認">
      {#if canCreate()}
        <p class="empty-text">未設定はありません。</p>
      {:else}
        <ul class="reason-list validation-only-list" aria-label="作成前に必要な対応">
          {#each blockedReasons() as reason (reason)}
            <li>{reason}</li>
          {/each}
        </ul>
      {/if}
      <button
        class="button-primary"
        disabled={!canCreate()}
        onclick={createJob}
        type="button"
      >
        次へ
      </button>
    </div>
  {/if}

  {#if apiKeyModalServiceId}
    {@const modalService = serviceFor(apiKeyModalServiceId)}
    <div class="modal-backdrop" role="presentation">
      <div
        aria-labelledby="apiKeyModalHeading"
        aria-modal="true"
        class="job-setup-card modal-panel"
        role="dialog"
      >
        <div class="section-head">
          <div>
            <p class="eyebrow">APIキー登録</p>
            <h3 id="apiKeyModalHeading">{modalService.label} のAPIキー</h3>
          </div>
        </div>
        <p class="lead">
          APIキーは AIサービスへリクエストを送るために使います。
          モデル一覧の取得と翻訳実行時の接続に必要です。
        </p>
        <p class="status-copy warning-copy">
          <strong>OS認証</strong>
          <span>
            保存時に OS の認証画面が表示されることがあります。APIキーを
            OS の安全な保管場所へ保存するためです。
          </span>
        </p>
        <label class="field-block" for="apiKeyInput">
          <span>APIキー</span>
          <input
            id="apiKeyInput"
            oninput={(event) => {
              const target = event.currentTarget
              if (target instanceof HTMLInputElement) {
                apiKeyDraft = target.value
              }
            }}
            placeholder="APIキーを入力"
            type="password"
            value={apiKeyDraft}
          />
        </label>
        <p class="mini-text">
          保存後は、このモーダルを閉じてモデル一覧を更新してください。
        </p>
        <div class="modal-actions">
          <button class="button-secondary" onclick={closeApiKeyModal} type="button">
            キャンセル
          </button>
          <button
            class="button-primary"
            disabled={!apiKeyDraft.trim() || apiKeySaveState === "saving"}
            onclick={saveApiKey}
            type="button"
          >
            {apiKeySaveState === "saving" ? "保存中" : "APIキーを保存"}
          </button>
        </div>
        {#if apiKeySaveState === "saved"}
          <p class="status-copy">
            <strong>保存済み</strong>
            <span>APIキーを登録しました。次にモデル一覧を更新してください。</span>
          </p>
        {/if}
      </div>
    </div>
  {/if}
</section>

<style>
  .job-setup-shell {
    display: grid;
    gap: 1.5rem;
    padding-bottom: 6rem;
  }

  .content-grid,
  .summary-grid {
    display: grid;
    gap: 1.25rem;
    grid-template-columns: 1fr;
  }

  .job-setup-card {
    display: grid;
    gap: 1rem;
    padding: 1.5rem;
    border: 1px solid rgba(255, 212, 165, 0.18);
    border-radius: 1.25rem;
    background: rgba(34, 26, 23, 0.82);
    box-shadow: 0 20px 40px rgba(6, 4, 3, 0.18);
  }

  .hero-card {
    gap: 0.75rem;
  }

  @media (min-width: 721px) {
    .job-setup-shell {
      gap: 1rem;
    }

    .content-grid,
    .summary-grid {
      gap: 0.9rem;
    }

    .hero-card,
    .compact-card {
      gap: 0.65rem;
      padding: 1rem;
      border-radius: 0.95rem;
    }

    .hero-card .lead,
    .compact-card > .mini-text {
      margin: 0;
    }

    .hero-card h2,
    .hero-card .eyebrow,
    .compact-card h3,
    .compact-card .eyebrow,
    .compact-card .mini-label,
    .compact-card .mini-text,
    .compact-card dt,
    .compact-card dd {
      margin: 0;
    }

    .compact-card .section-head {
      align-items: center;
    }

    .compact-card[aria-labelledby="jobSetupInputHeading"] {
      grid-template-columns: minmax(18rem, 0.55fr) minmax(0, 1fr);
      align-items: start;
    }

    .compact-card[aria-labelledby="jobSetupInputHeading"] .section-head {
      grid-column: 1 / -1;
    }

    .compact-card .field-block {
      gap: 0.3rem;
    }

    .compact-card select {
      padding-block: 0.65rem;
    }

    .compact-card .detail-grid {
      gap: 0.55rem;
    }

    .compact-card .detail-grid.compact {
      grid-template-columns: repeat(auto-fit, minmax(130px, 1fr));
    }

    .compact-card .foundation-table-wrap {
      max-height: 7.4rem;
    }

    .compact-card .foundation-table th,
    .compact-card .foundation-table td {
      padding: 0.34rem 0.6rem;
      font-size: 0.86rem;
    }

  }

  .hero-head,
  .section-head {
    display: flex;
    align-items: flex-start;
    justify-content: space-between;
    gap: 1rem;
  }

  .eyebrow,
  .mini-label {
    color: rgba(255, 215, 176, 0.72);
    font-size: 0.8rem;
    letter-spacing: 0.08em;
    text-transform: uppercase;
  }

  .lead,
  .mini-text,
  .empty-text,
  .reason-list {
    color: rgba(252, 241, 232, 0.86);
  }

  .status-pill {
    padding: 0.4rem 0.75rem;
    border-radius: 999px;
    background: rgba(255, 190, 126, 0.14);
    color: #ffd8ae;
    font-size: 0.85rem;
  }

  .status-pill.success,
  .tag.success {
    background: rgba(145, 208, 134, 0.16);
    color: #b8f0ad;
  }

  .tag.warning {
    background: rgba(255, 204, 128, 0.15);
    color: #ffd191;
  }

  .help-mark {
    display: inline-grid;
    place-items: center;
    width: 1.25rem;
    height: 1.25rem;
    margin-left: 0.35rem;
    padding: 0;
    border: 1px solid rgba(255, 212, 165, 0.28);
    border-radius: 999px;
    background: rgba(255, 241, 227, 0.08);
    color: #ffe2bf;
    font-size: 0.78rem;
    cursor: help;
  }

  .help-tip {
    position: relative;
    display: inline-grid;
    vertical-align: middle;
  }

  .help-tooltip {
    position: absolute;
    z-index: 30;
    bottom: calc(100% + 0.45rem);
    left: 50%;
    width: max-content;
    max-width: min(18rem, calc(100vw - 2rem));
    padding: 0.55rem 0.7rem;
    border: 1px solid rgba(255, 212, 165, 0.28);
    border-radius: 0.65rem;
    background: rgba(18, 13, 11, 0.98);
    color: #fff8f1;
    box-shadow: 0 14px 30px rgba(6, 4, 3, 0.28);
    font-size: 0.82rem;
    line-height: 1.45;
    text-align: left;
    transform: translate(-50%, 0.25rem);
    opacity: 0;
    visibility: hidden;
    pointer-events: none;
    transition:
      opacity 120ms ease,
      transform 120ms ease,
      visibility 120ms ease;
  }

  .help-tooltip::after {
    content: "";
    position: absolute;
    top: 100%;
    left: 50%;
    width: 0.55rem;
    height: 0.55rem;
    border-right: 1px solid rgba(255, 212, 165, 0.28);
    border-bottom: 1px solid rgba(255, 212, 165, 0.28);
    background: rgba(18, 13, 11, 0.98);
    transform: translate(-50%, -50%) rotate(45deg);
  }

  .help-tip:hover .help-tooltip,
  .help-tip:focus-within .help-tooltip {
    transform: translate(-50%, 0);
    opacity: 1;
    visibility: visible;
  }

  .status-copy {
    display: flex;
    gap: 0.75rem;
    align-items: baseline;
    flex-wrap: wrap;
  }

  .field-block {
    display: grid;
    gap: 0.45rem;
  }

  .field-block span,
  dt {
    color: rgba(255, 215, 176, 0.72);
    font-size: 0.9rem;
  }

  select,
  input {
    width: 100%;
    padding: 0.8rem 0.95rem;
    border: 1px solid rgba(255, 212, 165, 0.18);
    border-radius: 0.9rem;
    background: rgba(18, 13, 11, 0.92);
    color: #fef3e8;
  }

  .detail-grid {
    display: grid;
    gap: 0.9rem;
    grid-template-columns: repeat(auto-fit, minmax(180px, 1fr));
  }

  .detail-grid.compact {
    grid-template-columns: repeat(auto-fit, minmax(150px, 1fr));
  }

  .detail-grid div {
    display: grid;
    gap: 0.35rem;
  }

  .foundation-tables {
    display: grid;
    gap: 1rem;
    grid-template-columns: repeat(2, minmax(0, 1fr));
  }

  .foundation-tables > div {
    display: grid;
    gap: 0.35rem;
    min-width: 0;
  }

  .foundation-table-wrap {
    max-height: 11rem;
    overflow: auto;
    border: 1px solid rgba(255, 212, 165, 0.14);
    border-radius: 0.9rem;
    scrollbar-gutter: stable;
  }

  .foundation-table {
    width: 100%;
    border-collapse: collapse;
  }

  .foundation-table th,
  .foundation-table td {
    padding: 0.55rem 0.75rem;
    border-bottom: 1px solid rgba(255, 212, 165, 0.1);
    text-align: left;
    vertical-align: top;
  }

  .foundation-table th {
    position: sticky;
    top: 0;
    z-index: 1;
    background: rgba(18, 13, 11, 0.96);
    color: rgba(255, 215, 176, 0.8);
    font-size: 0.82rem;
  }

  .foundation-table td {
    color: #fff8f1;
    overflow-wrap: anywhere;
  }

  .phase-settings-list {
    display: grid;
    gap: 1rem;
    grid-template-columns: repeat(3, minmax(0, 1fr));
    align-items: stretch;
  }

  .phase-setting-block {
    display: grid;
    gap: 0.75rem;
    grid-template-rows: auto minmax(6.8rem, auto) minmax(12rem, 1fr) minmax(5.8rem, auto) auto;
    height: 100%;
    min-width: 0;
    padding: 1rem;
    border: 1px solid rgba(255, 212, 165, 0.14);
    border-radius: 0.9rem;
    background: rgba(18, 13, 11, 0.18);
  }

  .phase-setting-head {
    display: flex;
    align-items: flex-start;
    justify-content: space-between;
    gap: 1rem;
  }

  .phase-setting-head h4 {
    margin: 0;
    color: #fff8f1;
    font-size: 1rem;
  }

  .phase-control-group {
    display: grid;
    gap: 0.65rem;
    padding: 0.85rem;
    border: 1px solid rgba(255, 212, 165, 0.12);
    border-radius: 0.9rem;
    background: rgba(18, 13, 11, 0.24);
  }

  .phase-control-group.compact {
    align-content: start;
    gap: 0.45rem;
    padding: 0.65rem;
    background: rgba(18, 13, 11, 0.12);
  }

  .phase-control-group.model-control {
    align-content: start;
    min-height: 12rem;
  }

  .phase-message {
    margin: 0;
    padding: 0.75rem;
    border-radius: 0.9rem;
    background: rgba(255, 241, 227, 0.06);
    color: rgba(252, 241, 232, 0.9);
  }

  dd {
    margin: 0;
    color: #fff8f1;
  }

  .tag {
    padding: 0.4rem 0.7rem;
    border-radius: 999px;
    background: rgba(255, 241, 227, 0.1);
    color: #ffe2bf;
    font-size: 0.88rem;
  }

  .button-primary,
  .button-secondary {
    padding: 0.8rem 1rem;
    border-radius: 0.9rem;
    border: 1px solid rgba(255, 212, 165, 0.18);
    cursor: pointer;
  }

  .button-primary {
    background: linear-gradient(135deg, #ff9f5a, #ffcc88);
    color: #24150d;
  }

  .button-secondary {
    background: rgba(255, 241, 227, 0.08);
    color: #ffe2bf;
  }

  .modal-backdrop {
    position: fixed;
    inset: 0;
    z-index: 20;
    display: grid;
    place-items: center;
    padding: 1rem;
    background: rgba(6, 4, 3, 0.76);
    backdrop-filter: blur(14px);
    -webkit-backdrop-filter: blur(14px);
  }

  .modal-panel {
    width: min(560px, 100%);
    background: rgba(30, 22, 18, 0.96);
  }

  .modal-actions {
    display: flex;
    justify-content: flex-end;
    gap: 0.75rem;
    flex-wrap: wrap;
  }

  .warning-copy {
    padding: 0.75rem;
    border-radius: 0.9rem;
    background: rgba(255, 204, 128, 0.12);
  }

  button:disabled,
  select:disabled,
  input:disabled {
    opacity: 0.56;
    cursor: not-allowed;
  }

  .wrap-value,
  .reason-list li {
    overflow-wrap: anywhere;
    word-break: break-word;
  }

  .reason-list {
    margin: 0;
    padding-left: 1.1rem;
    display: grid;
    gap: 0.5rem;
  }

  .validation-only-list {
    padding: 0.75rem 0.75rem 0.75rem 1.75rem;
    border: 1px solid rgba(255, 204, 128, 0.18);
    border-radius: 0.9rem;
    background: rgba(255, 204, 128, 0.08);
  }

  .global-next-bar {
    position: fixed;
    z-index: 18;
    right: max(1rem, calc((100vw - 1440px) / 2 + 1rem));
    bottom: 1rem;
    left: max(1rem, calc((100vw - 1440px) / 2 + 1rem));
    display: flex;
    align-items: center;
    justify-content: space-between;
    gap: 1rem;
    padding: 0.75rem;
    border: 1px solid rgba(255, 212, 165, 0.18);
    border-radius: 0.9rem;
    background: rgba(24, 18, 15, 0.96);
    box-shadow: 0 18px 36px rgba(6, 4, 3, 0.28);
  }

  .global-next-bar .empty-text,
  .global-next-bar .reason-list {
    margin: 0;
  }

  .global-next-bar .reason-list {
    flex: 1;
  }

  @media (max-width: 720px) {
    .hero-head,
    .section-head,
    .phase-setting-head {
      flex-direction: column;
    }

    .job-setup-card {
      padding: 1.2rem;
    }

    .phase-settings-list {
      grid-template-columns: 1fr;
    }

    .foundation-tables {
      grid-template-columns: 1fr;
    }

    .global-next-bar {
      right: 0.75rem;
      bottom: 0.75rem;
      left: 0.75rem;
      flex-direction: column;
      align-items: stretch;
    }
  }
</style>
