<script lang="ts">
  import { onMount } from "svelte"

  import { createTranslationJobSetupRuntimeKey } from "@application/contract/translation-job-setup/translation-job-setup-screen-contract"
  import type {
    TranslationJobSetupExtendedViewModel,
    TranslationJobSetupPhaseCardViewModel
  } from "@application/presenter/translation-job-setup/translation-job-setup.presenter"
  import { VALIDATION_LABELS } from "@application/presenter/translation-job-setup"
  import AIModelSelectionCard from "@ui/components/AIModelSelectionCard.svelte"
  import StickyActionFooter from "@ui/components/StickyActionFooter.svelte"

  type CreateTranslationJobSetupScreenController =
    import("@application/contract/translation-job-setup/translation-job-setup-screen-contract").CreateTranslationJobSetupScreenController
  type TranslationJobSetupScreenControllerContract =
    import("@application/contract/translation-job-setup/translation-job-setup-screen-contract").TranslationJobSetupScreenControllerContract

  interface Props {
    createController: CreateTranslationJobSetupScreenController | null
    onReturnToInputReview?: (() => void) | null
  }

  let { createController, onReturnToInputReview = null }: Props = $props()

  function resolveController(): TranslationJobSetupScreenControllerContract {
    if (!createController) {
      throw new Error(
        "translation job setup screen controller factory is not provided"
      )
    }

    return createController()
  }

  function createEmptyViewModel(): TranslationJobSetupExtendedViewModel {
    return {
      phase: "idle",
      options: null,
      selectedInputSourceId: null,
      deletingInputSourceId: null,
      selectedRuntimeKey: null,
      selectedCredentialRef: "",
      phaseRuntimeSelections: [],
      providerModelLists: [],
      validationResult: null,
      validationState: "not-run",
      dirty: false,
      errorMessage: "",
      createErrorKind: null,
      summary: null,
      gatewayStatus: "未接続",
      selectedInputCandidate: null,
      selectedRuntimeOption: null,
      availableCredentialRefs: [],
      phaseValidationResults: [],
      phaseRuntimeSummaries: [],
      selectedInputLabel: "未選択",
      selectedInputSourceKind: "-",
      selectedInputRecordCountLabel: "-",
      selectedInputRegisteredAtLabel: "-",
      existingJobSummary: "既存 job はありません。",
      dictionaryLabels: [],
      personaLabels: [],
      validationStatusLabel: "未完了",
      validationStatusText: "不足があると作成できません。",
      createStatusText: "",
      blockedReasons: [],
      canValidate: false,
      canCreate: false,
      isLoading: false,
      isValidating: false,
      isCreating: false,
      hasExistingJob: false,
      showCacheMissingGuidance: false,
      credentialStateText: "",
      phaseCards: [],
      summaryPhaseCards: [],
      createSectionTitle: "作成前確認",
      createSectionText: "",
      globalBlockedReasons: []
    }
  }

  function normalizeViewModel(
    viewModel:
      | TranslationJobSetupExtendedViewModel
      | TranslationJobSetupScreenControllerContract["getViewModel"]
  ): TranslationJobSetupExtendedViewModel {
    const extendedViewModel = viewModel as TranslationJobSetupExtendedViewModel
    return {
      ...createEmptyViewModel(),
      ...extendedViewModel,
      phaseCards: extendedViewModel.phaseCards ?? [],
      summaryPhaseCards: extendedViewModel.summaryPhaseCards ?? [],
      globalBlockedReasons:
        extendedViewModel.globalBlockedReasons ??
        extendedViewModel.blockedReasons ??
        []
    }
  }

  const controller = resolveController()
  let viewModel = $state(
    normalizeViewModel(
      controller.getViewModel() as TranslationJobSetupExtendedViewModel
    )
  )

  const unsubscribe = controller.subscribe((nextViewModel) => {
    viewModel = normalizeViewModel(
      nextViewModel as TranslationJobSetupExtendedViewModel
    )
  })

  onMount(() => {
    void controller.mount()

    return () => {
      unsubscribe()
      controller.dispose()
    }
  })

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

  function formatRuntimeLabel(
    provider: string,
    model: string,
    mode: string
  ): string {
    return `${provider} / ${model} / ${mode}`
  }

  function resolveValidationLabel(status: string): string {
    return VALIDATION_LABELS[status as keyof typeof VALIDATION_LABELS] ?? status
  }

  function batchSectionText(
    phaseCard: TranslationJobSetupPhaseCardViewModel
  ): string {
    if (phaseCard.showBatchToggle) {
      return phaseCard.batchEnabled ? "有効" : "無効"
    }

    return "対象外"
  }

  function handlePhaseProviderChange(
    phaseId: TranslationJobSetupPhaseCardViewModel["phaseId"],
    event: Event
  ): void {
    const target = event.currentTarget
    if (target instanceof HTMLSelectElement) {
      controller.selectPhaseProvider(phaseId, target.value)
    }
  }

  function handlePhaseModelChange(
    phaseId: TranslationJobSetupPhaseCardViewModel["phaseId"],
    event: Event
  ): void {
    const target = event.currentTarget
    if (target instanceof HTMLSelectElement) {
      controller.selectPhaseModel(phaseId, target.value)
    }
  }

  function handlePhaseBatchChange(
    phaseId: TranslationJobSetupPhaseCardViewModel["phaseId"],
    event: Event
  ): void {
    const target = event.currentTarget
    if (target instanceof HTMLInputElement) {
      controller.togglePhaseBatchMode(phaseId, target.checked)
    }
  }

  function isSelectedInputCard(candidateId: number): boolean {
    return viewModel.selectedInputSourceId === candidateId
  }

  function isDeletingInputCard(candidateId: number): boolean {
    return viewModel.deletingInputSourceId === candidateId
  }
</script>

<section class="job-setup-shell" id="translationJobSetupView">
  <section class="job-setup-card hero-card">
    <h2>翻訳段階ごとの AI 設定</h2>
    <p class="lead">
      入力済みデータを確認し、3 つの翻訳段階で使う AI
      サービスとモデルを選びます。
    </p>
    <p class="error-text" hidden={!viewModel.errorMessage}>
      {viewModel.errorMessage}
    </p>
  </section>

  {#if viewModel.summary}
    <section class="summary-grid">
      <section class="job-setup-card" aria-labelledby="jobSetupSummaryHeading">
        <div class="section-head">
          <div>
            <p class="eyebrow">create result</p>
            <h3
              aria-label={viewModel.summaryPhaseCards.length === 0
                ? "Ready job summary"
                : undefined}
              id="jobSetupSummaryHeading"
            >
              作成済み設定
            </h3>
          </div>
          <span class="status-pill success">{viewModel.summary.jobState}</span>
        </div>
        <dl class="detail-grid compact">
          <div>
            <dt>job id</dt>
            <dd>{viewModel.summary.jobId}</dd>
          </div>
          <div>
            <dt>入力データ</dt>
            <dd class="wrap-value">{viewModel.summary.inputSource}</dd>
          </div>
          {#if viewModel.summaryPhaseCards.length === 0}
            <div>
              <dt>provider</dt>
              <dd class="wrap-value">
                {viewModel.summary.executionSummary.provider}
              </dd>
            </div>
            <div>
              <dt>model</dt>
              <dd class="wrap-value">
                {viewModel.summary.executionSummary.model}
              </dd>
            </div>
            <div>
              <dt>execution mode</dt>
              <dd>{viewModel.summary.executionSummary.executionMode}</dd>
            </div>
          {/if}
        </dl>
      </section>

      <section
        class="job-setup-card"
        aria-labelledby="jobSetupSummaryPhaseHeading"
      >
        <div class="section-head">
          <div>
            <p class="eyebrow">phase settings</p>
            <h3 id="jobSetupSummaryPhaseHeading">翻訳段階ごとの設定</h3>
          </div>
        </div>
        {#if viewModel.summaryPhaseCards.length === 0}
          <div class="tag-list">
            {#each viewModel.summary.validationPassSlices as slice (slice)}
              <span class="tag success">{slice}</span>
            {/each}
          </div>
        {:else}
          <div class="summary-phase-grid">
            {#each viewModel.summaryPhaseCards as summaryPhase (summaryPhase.phaseId)}
              <article class="summary-phase-card">
                <h4>{summaryPhase.phaseLabel}</h4>
                <dl class="detail-grid compact">
                  <div>
                    <dt>AIサービス</dt>
                    <dd>{summaryPhase.providerLabel}</dd>
                  </div>
                  <div>
                    <dt>モデル</dt>
                    <dd class="wrap-value">{summaryPhase.model}</dd>
                  </div>
                  <div>
                    <dt>APIキー状態</dt>
                    <dd>{summaryPhase.credentialStatusLabel}</dd>
                  </div>
                  <div>
                    <dt>一括処理</dt>
                    <dd>{summaryPhase.batchLabel}</dd>
                  </div>
                </dl>
              </article>
            {/each}
          </div>
        {/if}
      </section>
    </section>
  {:else}
    <section class="content-grid">
      <section class="job-setup-card" aria-labelledby="jobSetupInputHeading">
        <div class="section-head">
          <div>
            <h3 id="jobSetupInputHeading">入力データ</h3>
          </div>
        </div>
        <div class="input-card-list" aria-label="input data" role="list">
          {#each viewModel.options?.inputCandidates ?? [] as candidate (candidate.id)}
            <article
              aria-busy={isDeletingInputCard(candidate.id)}
              class:selected={isSelectedInputCard(candidate.id)}
              class="input-card"
              role="listitem"
            >
              <button
                aria-pressed={isSelectedInputCard(candidate.id)}
                class="input-card-select"
                disabled={
                  viewModel.isCreating || isDeletingInputCard(candidate.id)
                }
                onclick={() => controller.selectInputSource(candidate.id)}
                type="button"
              >
                <div class="input-card-head">
                  <strong class="wrap-value">{candidate.label}</strong>
                  {#if isDeletingInputCard(candidate.id)}
                    <span class="status-pill">削除中...</span>
                  {:else if isSelectedInputCard(candidate.id)}
                    <span class="status-pill success">選択中</span>
                  {/if}
                </div>
                <dl class="detail-grid compact">
                  <div>
                    <dt>出自</dt>
                    <dd class="wrap-value">{candidate.sourceKind}</dd>
                  </div>
                  <div>
                    <dt>翻訳レコード件数</dt>
                    <dd>{candidate.recordCount.toLocaleString("ja-JP")} 件</dd>
                  </div>
                  <div>
                    <dt>登録日時</dt>
                    <dd>{formatDate(candidate.registeredAt ?? "")}</dd>
                  </div>
                </dl>
              </button>
              <div class="input-card-actions">
                <button
                  class="button-secondary"
                  disabled={
                    viewModel.isCreating ||
                    viewModel.deletingInputSourceId !== null
                  }
                  onclick={() =>
                    void controller.deleteInputSource(candidate.id)}
                  type="button"
                >
                  {isDeletingInputCard(candidate.id) ? "削除中..." : "削除"}
                </button>
              </div>
            </article>
          {/each}
        </div>
        <dl class="detail-grid compact">
          <div>
            <dt>入力データ名</dt>
            <dd class="wrap-value">{viewModel.selectedInputLabel}</dd>
          </div>
          <div>
            <dt>出自</dt>
            <dd class="wrap-value">{viewModel.selectedInputSourceKind}</dd>
          </div>
          <div>
            <dt>登録日時</dt>
            <dd>{viewModel.selectedInputRegisteredAtLabel}</dd>
          </div>
          <div>
            <dt>翻訳レコード件数</dt>
            <dd>{viewModel.selectedInputRecordCountLabel}</dd>
          </div>
          <div>
            <dt>既存 job 状態</dt>
            <dd class="wrap-value">{viewModel.existingJobSummary}</dd>
          </div>
        </dl>
      </section>

      <section
        class="job-setup-card"
        aria-labelledby="jobSetupFoundationHeading"
      >
        <div class="section-head">
          <div>
            <h3 id="jobSetupFoundationHeading">共通辞書と共通ペルソナ</h3>
          </div>
        </div>
        <div class="foundation-grid split">
          <div class="foundation-table">
            <p class="mini-label">共通辞書</p>
            {#if viewModel.dictionaryLabels.length === 0}
              <p class="empty-text">利用可能な共通辞書はありません。</p>
            {:else}
              <div class="foundation-scroll">
                <ul class="plain-list">
                  {#each viewModel.dictionaryLabels as label (label)}
                    <li>{label}</li>
                  {/each}
                </ul>
              </div>
            {/if}
          </div>
          <div class="foundation-table">
            <p class="mini-label">共通ペルソナ</p>
            {#if viewModel.personaLabels.length === 0}
              <p class="empty-text">利用可能な共通ペルソナはありません。</p>
            {:else}
              <div class="foundation-scroll">
                <ul class="plain-list">
                  {#each viewModel.personaLabels as label (label)}
                    <li>{label}</li>
                  {/each}
                </ul>
              </div>
            {/if}
          </div>
        </div>
      </section>

      {#if viewModel.phaseCards.length === 0}
        <section
          class="job-setup-card"
          aria-labelledby="jobSetupLegacyRuntimeHeading"
        >
          <div class="section-head">
            <div>
              <p class="eyebrow">foundation and runtime</p>
              <h3 id="jobSetupLegacyRuntimeHeading">共通基盤と AI runtime</h3>
            </div>
          </div>
          <label class="field-block" for="jobSetupRuntimeSelect">
            <span>provider / model / execution mode</span>
            <select
              id="jobSetupRuntimeSelect"
              onchange={(event) => {
                const target = event.currentTarget
                if (target instanceof HTMLSelectElement) {
                  controller.selectRuntime(target.value)
                }
              }}
              value={viewModel.selectedRuntimeKey ?? undefined}
            >
              {#each viewModel.options?.aiRuntimeOptions ?? [] as option (createTranslationJobSetupRuntimeKey(option))}
                <option value={createTranslationJobSetupRuntimeKey(option)}>
                  {formatRuntimeLabel(
                    option.provider,
                    option.model,
                    option.mode
                  )}
                </option>
              {/each}
            </select>
          </label>
          <label class="field-block" for="jobSetupCredentialSelect">
            <span>credential reference</span>
            <select
              id="jobSetupCredentialSelect"
              onchange={(event) => {
                const target = event.currentTarget
                if (target instanceof HTMLSelectElement) {
                  controller.selectCredentialRef(target.value)
                }
              }}
              value={viewModel.selectedCredentialRef}
            >
              {#each viewModel.availableCredentialRefs as credential (credential.credentialRef)}
                <option value={credential.credentialRef}>
                  {credential.provider} / {credential.credentialRef}
                </option>
              {/each}
            </select>
          </label>
          <p class="mini-text">{viewModel.credentialStateText}</p>
        </section>

        <section
          class="job-setup-card"
          aria-labelledby="jobSetupValidationHeading"
        >
          <div class="section-head">
            <div>
              <p class="eyebrow">validation</p>
              <h3 id="jobSetupValidationHeading">Validation status</h3>
            </div>
            <button
              class="button-secondary"
              disabled={!viewModel.canValidate}
              onclick={() => void controller.runValidation()}
              type="button"
            >
              validation を実行
            </button>
          </div>
          <dl class="detail-grid compact">
            <div>
              <dt>状態</dt>
              <dd>
                {#if viewModel.validationResult}
                  {resolveValidationLabel(viewModel.validationResult.status)}
                {:else}
                  validation 未実行
                {/if}
              </dd>
            </div>
            <div>
              <dt>validated at</dt>
              <dd>
                {formatDate(viewModel.validationResult?.validatedAt ?? "")}
              </dd>
            </div>
            <div>
              <dt>blocking failure</dt>
              <dd class="wrap-value">
                {viewModel.validationResult?.blockingFailureCategory ?? "-"}
              </dd>
            </div>
            <div>
              <dt>dirty state</dt>
              <dd>{viewModel.dirty ? "dirty" : "clean"}</dd>
            </div>
          </dl>
          <div class="tag-list">
            {#each viewModel.validationResult?.targetSlices ?? [] as slice (slice)}
              <span class="tag warning">{slice}</span>
            {/each}
            {#each viewModel.validationResult?.passSlices ?? [] as slice (slice)}
              <span class="tag success">{slice}</span>
            {/each}
          </div>
        </section>
      {:else}
        <section class="job-setup-card" aria-labelledby="jobSetupPhaseHeading">
          <div class="section-head">
            <div>
              <h3 id="jobSetupPhaseHeading">AI サービスとモデル</h3>
            </div>
          </div>
          <div class="phase-grid">
            {#each viewModel.phaseCards as phaseCard (phaseCard.phaseId)}
              <div class="phase-card-wrap">
                <AIModelSelectionCard
                  helperText={phaseCard.helperText}
                  modelDisabled={!phaseCard.modelSelectEnabled}
                  modelOptions={phaseCard.modelOptions}
                  modelSelectId={`model-${phaseCard.phaseId}`}
                  modelStatusText={phaseCard.modelListStatusText}
                  modelValue={phaseCard.selectedModel}
                  onBatchChange={(event: Event) =>
                    handlePhaseBatchChange(phaseCard.phaseId, event)}
                  onModelChange={(event: Event) =>
                    handlePhaseModelChange(phaseCard.phaseId, event)}
                  onProviderChange={(event: Event) =>
                    handlePhaseProviderChange(phaseCard.phaseId, event)}
                  onRefresh={() =>
                    void controller.refreshPhaseModels(phaseCard.phaseId)}
                  providerDisabled={viewModel.isCreating}
                  providerOptions={phaseCard.providerOptions}
                  providerSelectId={`provider-${phaseCard.phaseId}`}
                  providerValue={phaseCard.provider}
                  refreshButtonAriaLabel={phaseCard.modelListButtonAriaLabel}
                  refreshButtonLabel={phaseCard.modelListButtonLabel}
                  refreshDisabled={!phaseCard.modelListButtonEnabled}
                  refreshSpinning={phaseCard.isModelListRefreshing}
                  secondaryControlMode={phaseCard.showBatchToggle
                    ? "batch-toggle"
                    : "none"}
                  showCredentialStatus={phaseCard.showCredentialStatus}
                  showCredentialWarning={phaseCard.showCredentialWarning}
                  statusLabel={phaseCard.statusLabel}
                  statusTone={phaseCard.statusTone}
                  title={phaseCard.phaseLabel}
                  titleTag="h4"
                  batchChecked={phaseCard.batchEnabled}
                  batchDisabled={viewModel.isCreating}
                  batchHelpId={`batch-help-${phaseCard.phaseId}`}
                  batchHelpText={phaseCard.batchHelpText}
                  credentialStatusLabel={phaseCard.credentialStatusLabel}
                  credentialStatusTone={phaseCard.credentialStatusTone}
                  credentialWarningText={phaseCard.credentialWarningText}
                  emptyModelLabel={phaseCard.showModelSelect
                    ? "選んでください"
                    : "モデル一覧を更新してください"}
                />
                <div class="detail-grid compact phase-detail-grid">
                  <div>
                    <dt>現在の状態</dt>
                    <dd>{phaseCard.statusLabel}</dd>
                  </div>
                  <div>
                    <dt>実行方法</dt>
                    <dd>{batchSectionText(phaseCard)}</dd>
                  </div>
                </div>
              </div>
            {/each}
          </div>
        </section>
      {/if}
    </section>

    <StickyActionFooter
      title="作成前確認"
      titleId="jobSetupCreateHeading"
      description={viewModel.createSectionText}
      reasons={viewModel.globalBlockedReasons}
      emptyText={viewModel.canCreate
        ? "不足はありません。"
        : "作成前確認はまだ未完了です。"}
      primaryLabel="次へ"
      primaryDisabled={!viewModel.canCreate}
      onPrimary={() => void controller.createJob()}
    >
      {#if viewModel.showCacheMissingGuidance}
        <p class="mini-text">
          cache missing は Job Setup で再構築しません。Input Review
          の再構築導線へ戻ってください。
        </p>
      {/if}
      {#if viewModel.showCacheMissingGuidance && onReturnToInputReview}
        <button
          class="button-secondary"
          onclick={() => onReturnToInputReview?.()}
          type="button"
        >
          Input Review へ戻る
        </button>
      {/if}
    </StickyActionFooter>
  {/if}
</section>

<style>
  .job-setup-shell {
    display: grid;
    gap: 1.25rem;
    padding-bottom: 10rem;
  }

  .content-grid,
  .summary-grid {
    display: grid;
    gap: 1rem;
    grid-template-columns: minmax(0, 1fr);
  }

  .job-setup-card {
    display: grid;
    gap: 1rem;
    padding: 1.25rem;
    border: 1px solid rgba(255, 212, 165, 0.18);
    border-radius: 1.25rem;
    background: rgba(34, 26, 23, 0.82);
    box-shadow: 0 20px 40px rgba(6, 4, 3, 0.18);
  }

  .hero-card {
    gap: 0.6rem;
  }

  .section-head {
    display: flex;
    align-items: flex-start;
    justify-content: space-between;
    gap: 1rem;
  }

  .mini-label {
    color: rgba(255, 215, 176, 0.72);
    font-size: 0.78rem;
    letter-spacing: 0.08em;
    text-transform: uppercase;
  }

  .lead,
  .mini-text,
  .empty-text,
  .plain-list {
    color: rgba(252, 241, 232, 0.86);
  }

  .status-pill {
    padding: 0.35rem 0.72rem;
    border-radius: 999px;
    background: rgba(255, 190, 126, 0.14);
    color: #ffd8ae;
    font-size: 0.82rem;
  }

  .status-pill.success {
    background: rgba(145, 208, 134, 0.16);
    color: #b8f0ad;
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

  select {
    width: 100%;
    padding: 0.8rem 0.95rem;
    border: 1px solid rgba(255, 212, 165, 0.18);
    border-radius: 0.9rem;
    background: rgba(18, 13, 11, 0.92);
    color: #fef3e8;
  }

  .detail-grid {
    display: grid;
    gap: 0.75rem;
    grid-template-columns: repeat(auto-fit, minmax(180px, 1fr));
  }

  .detail-grid.compact {
    grid-template-columns: repeat(auto-fit, minmax(150px, 1fr));
  }

  .detail-grid div,
  .foundation-grid,
  .phase-grid,
  .summary-phase-grid,
  .summary-phase-card,
  .foundation-table,
  .phase-card-wrap,
  .input-card,
  .input-card-actions,
  .input-card-select,
  .input-card-list {
    display: grid;
    gap: 0.75rem;
  }

  .input-card-list {
    grid-template-columns: repeat(auto-fit, minmax(260px, 1fr));
  }

  .input-card {
    padding: 0.9rem;
    border-radius: 1rem;
    border: 1px solid rgba(255, 212, 165, 0.12);
    background: rgba(18, 13, 11, 0.52);
  }

  .input-card.selected {
    border-color: rgba(255, 204, 136, 0.72);
    background: rgba(56, 39, 30, 0.78);
    box-shadow: 0 0 0 1px rgba(255, 204, 136, 0.22);
  }

  .input-card-select {
    text-align: left;
    border: 0;
    padding: 0;
    background: transparent;
    color: inherit;
  }

  .input-card-head {
    display: flex;
    align-items: flex-start;
    justify-content: space-between;
    gap: 0.75rem;
  }

  .input-card-actions {
    align-items: start;
  }

  .phase-grid {
    grid-template-columns: repeat(3, minmax(0, 1fr));
  }

  .summary-phase-grid {
    grid-template-columns: repeat(auto-fit, minmax(220px, 1fr));
  }

  .summary-phase-card {
    padding: 1rem;
    border-radius: 1rem;
    background: rgba(18, 13, 11, 0.62);
    border: 1px solid rgba(255, 212, 165, 0.1);
  }

  .foundation-grid.split {
    grid-template-columns: repeat(2, minmax(0, 1fr));
  }

  .foundation-scroll {
    max-height: 10rem;
    overflow: auto;
    padding-right: 0.35rem;
  }

  .plain-list {
    margin: 0;
    padding-left: 1.1rem;
  }

  .button-secondary {
    padding: 0.8rem 1rem;
    border-radius: 0.9rem;
    border: 1px solid rgba(255, 212, 165, 0.18);
    cursor: pointer;
  }

  .button-secondary {
    background: rgba(255, 241, 227, 0.08);
    color: #ffe2bf;
  }

  button:disabled,
  select:disabled {
    opacity: 0.56;
    cursor: not-allowed;
  }

  .wrap-value,
  .plain-list li {
    overflow-wrap: anywhere;
    word-break: break-word;
  }

  .phase-detail-grid {
    padding-top: 0.2rem;
  }

  .error-text {
    color: #ffb4ab;
    overflow-wrap: anywhere;
  }

  @media (max-width: 1080px) {
    .phase-grid {
      grid-template-columns: minmax(0, 1fr);
    }
  }

  @media (max-width: 720px) {
    .section-head {
      flex-direction: column;
    }

    .foundation-grid.split {
      grid-template-columns: minmax(0, 1fr);
    }

  }
</style>
