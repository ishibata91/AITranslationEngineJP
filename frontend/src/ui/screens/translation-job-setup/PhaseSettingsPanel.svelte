<script lang="ts">
  import type { TranslationJobSetupPhaseCardViewModel } from "@application/presenter/translation-job-setup/translation-job-setup.presenter"
  import AIModelSelectionCard from "@ui/components/AIModelSelectionCard.svelte"

  interface RuntimeOption {
    provider: string
    model: string
    mode: string
  }

  interface Props {
    isCreating: boolean
    phaseCards: TranslationJobSetupPhaseCardViewModel[]
    runtimeOptions: RuntimeOption[]
    selectedRuntimeKey: string | null
    batchSectionText: (
      phaseCard: TranslationJobSetupPhaseCardViewModel
    ) => string
    createRuntimeKey: (option: RuntimeOption) => string
    formatRuntimeLabel: (
      provider: string,
      model: string,
      mode: string
    ) => string
    onPhaseBatchChange: (
      phaseId: TranslationJobSetupPhaseCardViewModel["phaseId"],
      event: Event
    ) => void
    onPhaseModelChange: (
      phaseId: TranslationJobSetupPhaseCardViewModel["phaseId"],
      event: Event
    ) => void
    onPhaseProviderChange: (
      phaseId: TranslationJobSetupPhaseCardViewModel["phaseId"],
      event: Event
    ) => void
    onRefreshPhaseModels: (
      phaseId: TranslationJobSetupPhaseCardViewModel["phaseId"]
    ) => void
    onSelectRuntime: (runtimeKey: string) => void
  }

  let {
    isCreating,
    phaseCards,
    runtimeOptions,
    selectedRuntimeKey,
    batchSectionText,
    createRuntimeKey,
    formatRuntimeLabel,
    onPhaseBatchChange,
    onPhaseModelChange,
    onPhaseProviderChange,
    onRefreshPhaseModels,
    onSelectRuntime
  }: Props = $props()
</script>

{#if phaseCards.length === 0}
  <section
    class="job-setup-card"
    aria-labelledby="jobSetupLegacyRuntimeHeading"
    data-testid="translation-job-setup-ai-service-model-region"
  >
    <div class="section-head">
      <div>
        <p class="eyebrow">phase settings unavailable</p>
        <h3 id="jobSetupLegacyRuntimeHeading">翻訳段階別設定</h3>
      </div>
    </div>
    <label class="field-block" for="jobSetupRuntimeSelect">
      <span>AIサービス / モデル / 実行方法</span>
      <select
        id="jobSetupRuntimeSelect"
        onchange={(event) => {
          const target = event.currentTarget
          if (target instanceof HTMLSelectElement) {
            onSelectRuntime(target.value)
          }
        }}
        value={selectedRuntimeKey ?? undefined}
      >
        {#each runtimeOptions as option (createRuntimeKey(option))}
          <option value={createRuntimeKey(option)}>
            {formatRuntimeLabel(option.provider, option.model, option.mode)}
          </option>
        {/each}
      </select>
    </label>
    <p class="mini-text">
      翻訳段階別設定を取得できません。AIサービス、モデル、実行方法の選択だけを表示しています。
    </p>
  </section>
{:else}
  <section
    class="job-setup-card"
    aria-labelledby="jobSetupPhaseHeading"
    data-testid="translation-job-setup-ai-service-model-region"
  >
    <div class="section-head">
      <div>
        <h3 id="jobSetupPhaseHeading">AI サービスとモデル</h3>
      </div>
    </div>
    <div class="phase-grid">
      {#each phaseCards as phaseCard (phaseCard.phaseId)}
        <div class="phase-card-wrap">
          <AIModelSelectionCard
            helperText={phaseCard.helperText}
            modelDisabled={!phaseCard.modelSelectEnabled}
            modelOptions={phaseCard.modelOptions}
            modelSelectId={`model-${phaseCard.phaseId}`}
            modelStatusText={phaseCard.modelListStatusText}
            modelValue={phaseCard.selectedModel}
            onBatchChange={(event: Event) =>
              onPhaseBatchChange(phaseCard.phaseId, event)}
            onModelChange={(event: Event) =>
              onPhaseModelChange(phaseCard.phaseId, event)}
            onProviderChange={(event: Event) =>
              onPhaseProviderChange(phaseCard.phaseId, event)}
            onRefresh={() => onRefreshPhaseModels(phaseCard.phaseId)}
            providerDisabled={isCreating}
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
            batchDisabled={isCreating}
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

<style>
  .job-setup-card,
  .phase-grid,
  .phase-card-wrap,
  .detail-grid div {
    display: grid;
    gap: 0.75rem;
  }

  .job-setup-card {
    gap: 1rem;
    padding: 1.25rem;
    border: 1px solid rgba(255, 212, 165, 0.18);
    border-radius: 1.25rem;
    background: rgba(34, 26, 23, 0.82);
    box-shadow: 0 20px 40px rgba(6, 4, 3, 0.18);
    color: var(--text);
  }

  .section-head {
    display: flex;
    align-items: flex-start;
    justify-content: space-between;
    gap: 1rem;
  }

  .eyebrow {
    color: rgba(255, 215, 176, 0.72);
    font-size: 0.78rem;
    letter-spacing: 0.08em;
    text-transform: uppercase;
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

  .mini-text {
    color: rgba(252, 241, 232, 0.86);
  }

  .phase-grid {
    grid-template-columns: repeat(3, minmax(0, 1fr));
  }

  .detail-grid {
    display: grid;
    gap: 0.75rem;
    grid-template-columns: repeat(auto-fit, minmax(180px, 1fr));
  }

  .detail-grid.compact {
    grid-template-columns: repeat(auto-fit, minmax(150px, 1fr));
  }

  .phase-detail-grid {
    padding-top: 0.2rem;
  }

  select:disabled {
    opacity: 0.56;
    cursor: not-allowed;
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
  }
</style>
