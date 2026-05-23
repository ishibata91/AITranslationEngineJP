<script lang="ts">
  import { onMount } from "svelte"

  import { createTranslationJobSetupRuntimeKey } from "@application/contract/translation-job-setup/translation-job-setup-screen-contract"
  import type { TranslationJobSetupPhaseId } from "@application/gateway-contract/translation-job-setup"
  import type {
    TranslationJobSetupExtendedViewModel,
    TranslationJobSetupPhaseCardViewModel
  } from "@application/presenter/translation-job-setup/translation-job-setup.presenter"
  import { VALIDATION_LABELS } from "@application/presenter/translation-job-setup"
  import StickyActionFooter from "@ui/components/StickyActionFooter.svelte"

  import CompatibilityPrecheckPanel from "./CompatibilityPrecheckPanel.svelte"
  import CreatedJobSummaryPanel from "./CreatedJobSummaryPanel.svelte"
  import FoundationDataPanel from "./FoundationDataPanel.svelte"
  import InputSourcePanel from "./InputSourcePanel.svelte"
  import JobSetupPurposeHeader from "./JobSetupPurposeHeader.svelte"
  import PhaseSettingsPanel from "./PhaseSettingsPanel.svelte"
  import PhaseSettingsSummaryPanel from "./PhaseSettingsSummaryPanel.svelte"

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
</script>

<section class="job-setup-shell" id="translationJobSetupView">
  <JobSetupPurposeHeader errorMessage={viewModel.errorMessage} />

  {#if viewModel.summary}
    <section class="summary-grid">
      <CreatedJobSummaryPanel
        summary={viewModel.summary}
        summaryPhaseCount={viewModel.summaryPhaseCards.length}
      />

      <PhaseSettingsSummaryPanel
        legacyValidationPassSlices={viewModel.summary.validationPassSlices}
        summaryPhaseCards={viewModel.summaryPhaseCards}
      />
    </section>
  {:else}
    <section
      class="content-grid"
      data-testid="translation-job-setup-pre-create-settings-region"
    >
      <InputSourcePanel
        candidates={viewModel.options?.inputCandidates ?? []}
        deletingInputSourceId={viewModel.deletingInputSourceId}
        existingJobSummary={viewModel.existingJobSummary}
        isCreating={viewModel.isCreating}
        selectedInputLabel={viewModel.selectedInputLabel}
        selectedInputRecordCountLabel={viewModel.selectedInputRecordCountLabel}
        selectedInputRegisteredAtLabel={viewModel.selectedInputRegisteredAtLabel}
        selectedInputSourceId={viewModel.selectedInputSourceId}
        selectedInputSourceKind={viewModel.selectedInputSourceKind}
        {formatDate}
        onDeleteInputSource={(candidateId: number) =>
          void controller.deleteInputSource(candidateId)}
        onSelectInputSource={(candidateId: number) =>
          controller.selectInputSource(candidateId)}
      />

      <FoundationDataPanel
        dictionaryLabels={viewModel.dictionaryLabels}
        personaLabels={viewModel.personaLabels}
      />

      {#if viewModel.phaseCards.length === 0}
        <PhaseSettingsPanel
          isCreating={viewModel.isCreating}
          phaseCards={viewModel.phaseCards}
          runtimeOptions={viewModel.options?.aiRuntimeOptions ?? []}
          selectedRuntimeKey={viewModel.selectedRuntimeKey}
          {batchSectionText}
          createRuntimeKey={createTranslationJobSetupRuntimeKey}
          {formatRuntimeLabel}
          onPhaseBatchChange={handlePhaseBatchChange}
          onPhaseModelChange={handlePhaseModelChange}
          onPhaseProviderChange={handlePhaseProviderChange}
          onRefreshPhaseModels={(phaseId: TranslationJobSetupPhaseId) =>
            void controller.refreshPhaseModels(phaseId)}
          onSelectRuntime={(runtimeKey: string) =>
            controller.selectRuntime(runtimeKey)}
        />

        <CompatibilityPrecheckPanel
          canValidate={viewModel.canValidate}
          dirty={viewModel.dirty}
          validationResult={viewModel.validationResult}
          {formatDate}
          {resolveValidationLabel}
          onRunValidation={() => void controller.runValidation()}
        />
      {:else}
        <PhaseSettingsPanel
          isCreating={viewModel.isCreating}
          phaseCards={viewModel.phaseCards}
          runtimeOptions={viewModel.options?.aiRuntimeOptions ?? []}
          selectedRuntimeKey={viewModel.selectedRuntimeKey}
          {batchSectionText}
          createRuntimeKey={createTranslationJobSetupRuntimeKey}
          {formatRuntimeLabel}
          onPhaseBatchChange={handlePhaseBatchChange}
          onPhaseModelChange={handlePhaseModelChange}
          onPhaseProviderChange={handlePhaseProviderChange}
          onRefreshPhaseModels={(phaseId: TranslationJobSetupPhaseId) =>
            void controller.refreshPhaseModels(phaseId)}
          onSelectRuntime={(runtimeKey: string) =>
            controller.selectRuntime(runtimeKey)}
        />
      {/if}
    </section>

    <StickyActionFooter
      dataTestId="translation-job-setup-job-create-sticky-footer"
      title="ジョブの作成確認"
      titleId="jobSetupCreateHeading"
      description="入力データと翻訳設定を確認し、最初の翻訳段階へ進む準備をします。"
      reasons={viewModel.globalBlockedReasons}
      emptyText={viewModel.canCreate
        ? "作成に必要な確認は完了しています。"
        : "作成前に確認が必要な項目があります。"}
      primaryLabel="単語翻訳へ進む"
      primaryDisabled={!viewModel.canCreate}
      onPrimary={() => void controller.createJob()}
    >
      {#if viewModel.showCacheMissingGuidance}
        <p class="mini-text">
          入力データの再構築が必要です。入力データの確認画面に戻ってください。
        </p>
      {/if}
      {#if viewModel.showCacheMissingGuidance && onReturnToInputReview}
        <button
          class="button-secondary"
          onclick={() => onReturnToInputReview?.()}
          type="button"
        >
          入力データの確認へ戻る
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

  .mini-text {
    color: rgba(252, 241, 232, 0.86);
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

  button:disabled {
    opacity: 0.56;
    cursor: not-allowed;
  }
</style>
