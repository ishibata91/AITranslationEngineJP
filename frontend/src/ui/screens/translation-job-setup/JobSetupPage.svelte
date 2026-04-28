<script lang="ts">
  import { onMount } from "svelte"

  import type {
    CreateTranslationJobSetupScreenController,
    TranslationJobSetupScreenControllerContract
  } from "@application/contract/translation-job-setup"
  import { VALIDATION_LABELS } from "@application/presenter/translation-job-setup"
  import { createTranslationJobSetupRuntimeKey } from "@application/gateway-contract/translation-job-setup"

  interface Props {
    createController: CreateTranslationJobSetupScreenController | null
    onReturnToInputReview?: (() => void) | null
  }

  let {
    createController,
    onReturnToInputReview = null
  }: Props = $props()

  function resolveController(): TranslationJobSetupScreenControllerContract {
    if (!createController) {
      throw new Error("translation job setup screen controller factory is not provided")
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

  function formatRuntimeLabel(provider: string, model: string, mode: string): string {
    return `${provider} / ${model} / ${mode}`
  }

  function openInputReview(): void {
    onReturnToInputReview?.()
  }
</script>

<section class="job-setup-shell" id="translationJobSetupView">
  <section class="job-setup-card hero-card">
    <div class="hero-head">
      <div>
        <p class="eyebrow">translation-management</p>
        <h2>Job Setup</h2>
      </div>
      <p class="gateway-status">Gateway: {viewModel.gatewayStatus}</p>
    </div>
    <p class="lead">
      入力、共通基盤、AI runtime、validation、create result を同じ画面で確認し、create 条件を満たした時だけ ready job を作成します。
    </p>
    <p class="status-copy">
      <strong>{viewModel.validationStatusLabel}</strong>
      <span>{viewModel.validationStatusText}</span>
    </p>
    <p class="status-copy">
      <strong>create</strong>
      <span>{viewModel.createStatusText}</span>
    </p>
    <p class="error-text" hidden={!viewModel.errorMessage}>{viewModel.errorMessage}</p>
  </section>

  {#if viewModel.summary}
    <section class="summary-grid">
      <section class="job-setup-card" aria-labelledby="jobSetupSummaryHeading">
        <div class="section-head">
          <div>
            <p class="eyebrow">create result</p>
            <h3 id="jobSetupSummaryHeading">Ready job summary</h3>
          </div>
          <span class="status-pill success">{viewModel.summary.jobState}</span>
        </div>
        <dl class="detail-grid">
          <div>
            <dt>job id</dt>
            <dd>{viewModel.summary.jobId}</dd>
          </div>
          <div>
            <dt>input source</dt>
            <dd class="wrap-value">{viewModel.summary.inputSource}</dd>
          </div>
          <div>
            <dt>provider</dt>
            <dd class="wrap-value">{viewModel.summary.executionSummary.provider}</dd>
          </div>
          <div>
            <dt>model</dt>
            <dd class="wrap-value">{viewModel.summary.executionSummary.model}</dd>
          </div>
          <div>
            <dt>execution mode</dt>
            <dd>{viewModel.summary.executionSummary.executionMode}</dd>
          </div>
        </dl>
      </section>

      <section class="job-setup-card" aria-labelledby="jobSetupPassSlicesHeading">
        <div class="section-head">
          <div>
            <p class="eyebrow">validation pass</p>
            <h3 id="jobSetupPassSlicesHeading">Validation pass slices</h3>
          </div>
        </div>
        {#if viewModel.summary.validationPassSlices.length === 0}
          <p class="empty-text">pass slice はありません。</p>
        {:else}
          <div class="tag-list">
            {#each viewModel.summary.validationPassSlices as slice (slice)}
              <span class="tag success">{slice}</span>
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
            <p class="eyebrow">input</p>
            <h3 id="jobSetupInputHeading">入力データ</h3>
          </div>
          <span class="mini-text">{viewModel.isLoading ? "loading" : "ready"}</span>
        </div>
        <label class="field-block" for="jobSetupInputSelect">
          <span>input data</span>
          <select
            disabled={viewModel.isLoading || viewModel.isValidating || viewModel.isCreating}
            id="jobSetupInputSelect"
            onchange={(event) => {
              const target = event.currentTarget
              if (target instanceof HTMLSelectElement) {
                controller.selectInputSource(Number(target.value))
              }
            }}
            value={viewModel.selectedInputSourceId ?? undefined}
          >
            {#each viewModel.options?.inputCandidates ?? [] as candidate (candidate.id)}
              <option value={candidate.id}>{candidate.label}</option>
            {/each}
          </select>
        </label>
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

      <section class="job-setup-card" aria-labelledby="jobSetupFoundationHeading">
        <div class="section-head">
          <div>
            <p class="eyebrow">foundation and runtime</p>
            <h3 id="jobSetupFoundationHeading">共通基盤と AI runtime</h3>
          </div>
        </div>
        <div class="foundation-grid">
          <div>
            <p class="mini-label">共通辞書</p>
            {#if viewModel.dictionaryLabels.length === 0}
              <p class="empty-text">利用可能な共通辞書はありません。</p>
            {:else}
              <div class="tag-list">
                {#each viewModel.dictionaryLabels as label (label)}
                  <span class="tag">{label}</span>
                {/each}
              </div>
            {/if}
          </div>
          <div>
            <p class="mini-label">共通ペルソナ</p>
            {#if viewModel.personaLabels.length === 0}
              <p class="empty-text">利用可能な共通ペルソナはありません。</p>
            {:else}
              <div class="tag-list">
                {#each viewModel.personaLabels as label (label)}
                  <span class="tag">{label}</span>
                {/each}
              </div>
            {/if}
          </div>
        </div>
        <label class="field-block" for="jobSetupRuntimeSelect">
          <span>provider / model / execution mode</span>
          <select
            disabled={viewModel.isLoading || viewModel.isValidating || viewModel.isCreating}
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
                {formatRuntimeLabel(option.provider, option.model, option.mode)}
              </option>
            {/each}
          </select>
        </label>
        <label class="field-block" for="jobSetupCredentialSelect">
          <span>credential reference</span>
          <select
            disabled={viewModel.isLoading || viewModel.isValidating || viewModel.isCreating}
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

      <section class="job-setup-card" aria-labelledby="jobSetupValidationHeading">
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
                {VALIDATION_LABELS[viewModel.validationResult.status] ?? viewModel.validationResult.status}
              {:else}
                validation 未実行
              {/if}
            </dd>
          </div>
          <div>
            <dt>validated at</dt>
            <dd>{formatDate(viewModel.validationResult?.validatedAt ?? "")}</dd>
          </div>
          <div>
            <dt>blocking failure</dt>
            <dd class="wrap-value">{viewModel.validationResult?.blockingFailureCategory ?? "-"}</dd>
          </div>
          <div>
            <dt>dirty state</dt>
            <dd>{viewModel.dirty ? "dirty" : "clean"}</dd>
          </div>
        </dl>
        <div class="slice-block">
          <p class="mini-label">target slices</p>
          {#if viewModel.validationResult?.targetSlices.length}
            <div class="tag-list">
              {#each viewModel.validationResult?.targetSlices ?? [] as slice (slice)}
                <span class="tag warning">{slice}</span>
              {/each}
            </div>
          {:else}
            <p class="empty-text">target slice はありません。</p>
          {/if}
        </div>
        <div class="slice-block">
          <p class="mini-label">pass slices</p>
          {#if viewModel.validationResult?.passSlices.length}
            <div class="tag-list">
              {#each viewModel.validationResult?.passSlices ?? [] as slice (slice)}
                <span class="tag success">{slice}</span>
              {/each}
            </div>
          {:else}
            <p class="empty-text">pass slice はありません。</p>
          {/if}
        </div>
        {#if viewModel.showCacheMissingGuidance}
          <div class="callout warning">
            <p>cache missing は Job Setup で再構築しません。Input Review の再構築導線へ戻ってください。</p>
            <button class="button-secondary" onclick={openInputReview} type="button">
              Input Review へ戻る
            </button>
          </div>
        {/if}
      </section>

      <section class="job-setup-card" aria-labelledby="jobSetupCreateHeading">
        <div class="section-head">
          <div>
            <p class="eyebrow">create job</p>
            <h3 id="jobSetupCreateHeading">Create ready job</h3>
          </div>
          <button
            class="button-primary"
            disabled={!viewModel.canCreate}
            onclick={() => void controller.createJob()}
            type="button"
          >
            ready job を作成
          </button>
        </div>
        {#if viewModel.blockedReasons.length === 0}
          <p class="empty-text">create 条件を満たしています。</p>
        {:else}
          <ul class="reason-list">
            {#each viewModel.blockedReasons as reason (reason)}
              <li>{reason}</li>
            {/each}
          </ul>
        {/if}
      </section>
    </section>
  {/if}
</section>

<style>
  .job-setup-shell {
    display: grid;
    gap: 1.5rem;
  }

  .content-grid,
  .summary-grid {
    display: grid;
    gap: 1.25rem;
    grid-template-columns: repeat(auto-fit, minmax(280px, 1fr));
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
  .reason-list,
  .callout {
    color: rgba(252, 241, 232, 0.86);
  }

  .gateway-status,
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
    gap: 0.9rem;
    grid-template-columns: repeat(auto-fit, minmax(180px, 1fr));
  }

  .detail-grid.compact {
    grid-template-columns: repeat(auto-fit, minmax(150px, 1fr));
  }

  .detail-grid div,
  .foundation-grid {
    display: grid;
    gap: 0.35rem;
  }

  dd {
    margin: 0;
    color: #fff8f1;
  }

  .tag-list {
    display: flex;
    gap: 0.5rem;
    flex-wrap: wrap;
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

  button:disabled,
  select:disabled {
    opacity: 0.56;
    cursor: not-allowed;
  }

  .wrap-value,
  .reason-list li,
  .callout p {
    overflow-wrap: anywhere;
    word-break: break-word;
  }

  .reason-list {
    margin: 0;
    padding-left: 1.1rem;
    display: grid;
    gap: 0.5rem;
  }

  .callout {
    display: grid;
    gap: 0.75rem;
    padding: 1rem;
    border-radius: 1rem;
    background: rgba(255, 213, 149, 0.1);
  }

  .error-text {
    color: #ffb4ab;
    overflow-wrap: anywhere;
  }

  @media (max-width: 720px) {
    .hero-head,
    .section-head {
      flex-direction: column;
    }

    .job-setup-card {
      padding: 1.2rem;
    }
  }
</style>