<script lang="ts">
  import { onMount } from "svelte"

  import type { TranslationJobManagementJobRunTarget } from "@application/contract/translation-job-management/translation-job-management-screen-types"
  import type {
    BodyTranslationPhaseActionKind,
    BodyTranslationPhaseScreenViewModel,
    BodyTranslationPhaseScreenControllerContract,
    CreateBodyTranslationPhaseScreenController
  } from "@application/contract/body-translation-phase"
  import type {
    CreatePersonaGenerationPhaseScreenController,
    PersonaGenerationPhaseActionKind,
    PersonaGenerationPhaseScreenControllerContract
  } from "@application/contract/persona-generation-phase"
  import type {
    CreateTermTranslationPhaseScreenController,
    TermTranslationPhaseActionCard,
    TermTranslationPhaseScreenControllerContract
  } from "@application/contract/term-translation-phase"
  import BodyTranslationPhasePanel from "@ui/screens/body-translation-phase/BodyTranslationPhasePanel.svelte"
  import JobRunTargetSummary from "@ui/screens/job-run/JobRunTargetSummary.svelte"
  import JobUnselectedGuidance from "@ui/screens/job-run/JobUnselectedGuidance.svelte"
  import PhaseHost from "@ui/screens/job-run/PhaseHost.svelte"
  import PhaseNavigationFooter from "@ui/screens/job-run/PhaseNavigationFooter.svelte"
  import TranslationCompletePage from "@ui/screens/job-run/TranslationCompletePage.svelte"
  import PersonaGenerationPhasePanel from "@ui/screens/persona-generation-phase/PersonaGenerationPhasePanel.svelte"
  import TermTranslationPhasePanel from "@ui/screens/term-translation-phase/TermTranslationPhasePanel.svelte"
  import type {
    JobRunPhaseStepId,
    ProcessingTargetListItem
  } from "@ui/screens/job-run/job-run-shell-props"

  interface Props {
    createBodyController: CreateBodyTranslationPhaseScreenController | null
    createPersonaController: CreatePersonaGenerationPhaseScreenController | null
    createController: CreateTermTranslationPhaseScreenController | null
    processingTargetItemsByPhase?: Partial<
      Record<PhasePageId, ProcessingTargetListItem[]>
    >
    selectedJobTarget?: TranslationJobManagementJobRunTarget | null
    onOpenJobManagement?: () => void
    onOpenOutputManagement?: () => void
  }

  let {
    createBodyController,
    createPersonaController,
    createController,
    processingTargetItemsByPhase = {},
    selectedJobTarget = null,
    onOpenJobManagement = () => undefined,
    onOpenOutputManagement = () => undefined
  }: Props = $props()

  type PhasePageId = JobRunPhaseStepId

  function resolveController(): TermTranslationPhaseScreenControllerContract {
    if (!createController) {
      throw new Error(
        "term translation phase screen controller factory is not provided"
      )
    }

    return createController()
  }

  function resolvePersonaController(): PersonaGenerationPhaseScreenControllerContract {
    if (!createPersonaController) {
      throw new Error(
        "persona generation phase screen controller factory is not provided"
      )
    }

    return createPersonaController()
  }

  function resolveBodyController(): BodyTranslationPhaseScreenControllerContract | null {
    return createBodyController ? createBodyController() : null
  }

  const controller = resolveController()
  const personaController = resolvePersonaController()
  const bodyController = resolveBodyController()
  let viewModel = $state(controller.getViewModel())
  let personaViewModel = $state(personaController.getViewModel())
  let bodyViewModel = $state<BodyTranslationPhaseScreenViewModel | null>(
    bodyController ? bodyController.getViewModel() : null
  )
  let currentPhasePage = $state<PhasePageId>("term")

  function resolveInitialPhasePage(
    target: TranslationJobManagementJobRunTarget | null
  ): PhasePageId {
    if (target?.currentPhase === "persona_generation") {
      return "persona"
    }
    if (target?.currentPhase === "body_translation") {
      return "body"
    }
    return "term"
  }

  function getBodyProcessingTargetPageState(phase: string) {
    return bodyViewModel?.processingTargetPageStatesByPhase?.[phase]
  }

  function getBodyProcessingTargetPage(phase: string): number {
    return getBodyProcessingTargetPageState(phase)?.page ?? 1
  }

  function setCurrentPhasePage(phasePage: PhasePageId): void {
    currentPhasePage = phasePage
    if (phasePage === "body") {
      void bodyController?.setProcessingTargetPage?.(
        getBodyProcessingTargetPage("body_translation"),
        "body_translation"
      )
    }
    if (phasePage === "complete") {
      void bodyController?.setProcessingTargetPage?.(
        getBodyProcessingTargetPage("translation_complete"),
        "translation_complete"
      )
    }
  }

  $effect(() => {
    if (!selectedJobTarget) {
      setCurrentPhasePage("term")
      void Promise.all([
        controller.setJobId(null),
        personaController.setJobId(null),
        ...(bodyController ? [bodyController.setJobId(null)] : [])
      ])
      return
    }

    setCurrentPhasePage(resolveInitialPhasePage(selectedJobTarget))
    void Promise.all([
      controller.setJobId(selectedJobTarget.jobId),
      personaController.setJobId(selectedJobTarget.jobId),
      ...(bodyController
        ? [bodyController.setJobId(selectedJobTarget.jobId)]
        : [])
    ])
  })

  const unsubscribe = controller.subscribe((nextViewModel) => {
    viewModel = nextViewModel
  })

  const unsubscribePersona = personaController.subscribe((nextViewModel) => {
    personaViewModel = nextViewModel
  })

  const unsubscribeBody = bodyController
    ? bodyController.subscribe((nextViewModel) => {
        bodyViewModel = nextViewModel
      })
    : () => undefined

  onMount(() => {
    void Promise.all([
      controller.mount(),
      personaController.mount(),
      ...(bodyController ? [bodyController.mount()] : [])
    ])

    return () => {
      unsubscribe()
      unsubscribePersona()
      unsubscribeBody()
      controller.dispose()
      personaController.dispose()
      bodyController?.dispose()
    }
  })

  async function handleAction(
    actionId: TermTranslationPhaseActionCard["id"]
  ): Promise<void> {
    if (actionId === "next-phase") {
      return
    }

    switch (actionId) {
      case "start":
        await controller.startPhase()
        return
      case "pause":
        await controller.pausePhase()
        return
      case "resume":
        await controller.resumePhase()
        return
      case "retry":
        await controller.retryPhase()
        return
    }
  }

  async function handlePersonaAction(
    actionId: PersonaGenerationPhaseActionKind
  ): Promise<void> {
    switch (actionId) {
      case "start":
        await personaController.startPhase()
        return
      case "pause":
        await personaController.pausePhase()
        return
      case "resume":
        await personaController.resumePhase()
        return
      case "retry":
        await personaController.retryPhase()
        return
      case "cancel":
        await personaController.cancelPhase()
        return
      case "check-body-readiness":
        await personaController.checkBodyReadiness()
        return
      case "start-body-phase":
        await personaController.startBodyPhase()
        return
    }
  }

  async function handleBodyAction(
    actionId: BodyTranslationPhaseActionKind
  ): Promise<void> {
    if (!bodyController) {
      return
    }

    switch (actionId) {
      case "start":
        await bodyController.startPhase()
        return
      case "pause":
        await bodyController.pausePhase()
        return
      case "resume":
        await bodyController.resumePhase()
        return
      case "retry":
        await bodyController.retryPhase()
        return
      case "cancel":
        await bodyController.cancelPhase()
        return
      case "check-output-readiness":
        await bodyController.checkOutputReadiness()
        return
    }
  }

  type PhaseAISettingsRequest = {
    provider: string
    model: string
    executionMode: string
    batchMode: string
  }

  async function handleTermAISettingsChange(
    request: PhaseAISettingsRequest
  ): Promise<void> {
    await controller.saveAISettings?.(request)
  }

  async function handlePersonaAISettingsChange(
    request: PhaseAISettingsRequest
  ): Promise<void> {
    await personaController.saveAISettings?.(request)
  }

  async function handleBodyAISettingsChange(
    request: PhaseAISettingsRequest
  ): Promise<void> {
    await bodyController?.saveAISettings?.(request)
  }

  function reasonsFrom(...values: Array<string | undefined>): string[] {
    return values.filter((value): value is string => Boolean(value))
  }

  const canOpenPersonaPhase = $derived(
    viewModel.nextPhaseReadiness?.canStartNextPhase === true ||
      viewModel.actionEnablement?.canStartNextPhase === true
  )
  const termFooterReasons = $derived(
    canOpenPersonaPhase
      ? []
      : reasonsFrom(
          viewModel.nextPhaseReadiness?.blockedReason,
          viewModel.nextPhaseBlockedReason,
          "次へ進めません。単語翻訳の完了状況と辞書の参照状態を確認してください。"
        )
  )
  const canOpenBodyPhase = $derived(
    personaViewModel.bodyReadiness?.ready === true ||
      personaViewModel.latestBodyReadiness?.ready === true
  )
  const personaFooterReasons = $derived(
    canOpenBodyPhase
      ? []
      : reasonsFrom(
          personaViewModel.bodyReadiness?.blockedReason,
          personaViewModel.bodyReadinessBlockedReason,
          "次へ進めません。ペルソナ生成の完了状況と参照状態を確認してください。"
        )
  )
  const canOpenCompletePage = $derived(
    bodyViewModel?.latestOutputReadiness?.ready === true &&
      bodyViewModel.viewState === "completed"
  )
  const bodyFooterReasons = $derived(
    canOpenCompletePage
      ? []
      : reasonsFrom(
          bodyViewModel?.latestOutputReadiness?.blockedReason,
          bodyViewModel?.outputReadinessBlockedReason,
          "完了確認へ進めません。本文翻訳の完了状況と翻訳結果を確認してください。"
        )
  )
  const currentProcessingTargetItems = $derived(
    processingTargetItemsByPhase[currentPhasePage] ?? []
  )
  const currentProcessingTargetPageState = $derived.by(() => {
    if (currentProcessingTargetItems.length > 0) {
      return undefined
    }

    switch (currentPhasePage) {
      case "term":
        return viewModel.processingTargetPageState ?? undefined
      case "persona":
        return personaViewModel.processingTargetPageState ?? undefined
      case "body":
        return getBodyProcessingTargetPageState("body_translation")
      case "complete":
        return getBodyProcessingTargetPageState("translation_complete")
      default:
        return undefined
    }
  })

  async function handleProcessingTargetSearchInput(
    event: Event
  ): Promise<void> {
    const target = event.target
    if (!(target instanceof HTMLInputElement)) {
      return
    }

    switch (currentPhasePage) {
      case "term":
        await controller.setProcessingTargetSearchQuery?.(target.value)
        return
      case "persona":
        await personaController.setProcessingTargetSearchQuery?.(target.value)
        return
      case "body":
        await bodyController?.setProcessingTargetSearchQuery?.(target.value)
        return
      case "complete":
        await bodyController?.setProcessingTargetSearchQuery?.(
          target.value,
          "translation_complete"
        )
        return
    }
  }

  async function handleProcessingTargetPageChange(page: number): Promise<void> {
    switch (currentPhasePage) {
      case "term":
        await controller.setProcessingTargetPage?.(page)
        return
      case "persona":
        await personaController.setProcessingTargetPage?.(page)
        return
      case "body":
        await bodyController?.setProcessingTargetPage?.(page)
        return
      case "complete":
        await bodyController?.setProcessingTargetPage?.(
          page,
          "translation_complete"
        )
        return
    }
  }
</script>

<section class="job-run-page" data-testid="job-run-job-run-shell">
  {#if selectedJobTarget}
    <JobRunTargetSummary target={selectedJobTarget} {currentPhasePage} />
    <PhaseHost>
      {#if currentPhasePage === "term"}
        <TermTranslationPhasePanel
          {viewModel}
          onAction={(actionId: TermTranslationPhaseActionCard["id"]) =>
            handleAction(actionId)}
          processingTargetItems={currentProcessingTargetItems}
          processingTargetPageState={currentProcessingTargetPageState}
          onProcessingTargetSearchInput={handleProcessingTargetSearchInput}
          onProcessingTargetPageChange={handleProcessingTargetPageChange}
          onAISettingsChange={handleTermAISettingsChange}
        />
        <PhaseNavigationFooter
          dataTestId="job-run-next-action-footer"
          description="単語翻訳が完了し、辞書を参照できる場合だけ次へ進めます。"
          onBack={onOpenJobManagement}
          onPrimary={() => {
            setCurrentPhasePage("persona")
          }}
          primaryDisabled={!canOpenPersonaPhase}
          reasons={termFooterReasons}
          title="単語翻訳の次の作業"
          titleId="termPhaseNavigationFooter"
        />
      {:else if currentPhasePage === "persona"}
        <PersonaGenerationPhasePanel
          viewModel={personaViewModel}
          onAction={(actionId: PersonaGenerationPhaseActionKind) =>
            handlePersonaAction(actionId)}
          processingTargetItems={currentProcessingTargetItems}
          processingTargetPageState={currentProcessingTargetPageState}
          onProcessingTargetSearchInput={handleProcessingTargetSearchInput}
          onProcessingTargetPageChange={handleProcessingTargetPageChange}
          onAISettingsChange={handlePersonaAISettingsChange}
        />
        <PhaseNavigationFooter
          dataTestId="job-run-next-action-footer"
          description="ペルソナ生成が完了し、生成結果を参照できる場合だけ次へ進めます。"
          onBack={onOpenJobManagement}
          onPrimary={() => {
            setCurrentPhasePage("body")
          }}
          primaryDisabled={!canOpenBodyPhase}
          reasons={personaFooterReasons}
          title="NPC ペルソナ生成の次の作業"
          titleId="personaPhaseNavigationFooter"
        />
      {:else if currentPhasePage === "body" && bodyViewModel}
        <BodyTranslationPhasePanel
          viewModel={bodyViewModel}
          onAction={(actionId: BodyTranslationPhaseActionKind) =>
            handleBodyAction(actionId)}
          processingTargetItems={currentProcessingTargetItems}
          processingTargetPageState={currentProcessingTargetPageState}
          onProcessingTargetSearchInput={handleProcessingTargetSearchInput}
          onProcessingTargetPageChange={handleProcessingTargetPageChange}
          onAISettingsChange={handleBodyAISettingsChange}
        />
        <PhaseNavigationFooter
          dataTestId="body-translation-phase-complete-next-action"
          description="本文翻訳が完了し、翻訳結果を確認できる場合だけ完了確認へ進めます。"
          onBack={onOpenJobManagement}
          onPrimary={() => {
            setCurrentPhasePage("complete")
          }}
          primaryDisabled={!canOpenCompletePage}
          reasons={bodyFooterReasons}
          title="本文翻訳の次の作業"
          titleId="bodyPhaseNavigationFooter"
        />
      {:else if currentPhasePage === "complete" && bodyViewModel}
        <TranslationCompletePage
          jobId={selectedJobTarget.jobId}
          rows={bodyViewModel.fieldResultItems}
          processingTargetPageState={currentProcessingTargetPageState}
          onProcessingTargetPageChange={handleProcessingTargetPageChange}
        />
        <PhaseNavigationFooter
          dataTestId="translation-complete-post-completion-next-action"
          description="翻訳結果を確認した後は、出力管理で出力対象を選びます。"
          onBack={onOpenJobManagement}
          onPrimary={onOpenOutputManagement}
          primaryLabel="出力管理へ進む"
          reasons={[]}
          title="翻訳完了後の次の作業"
          titleId="completeNavigationFooter"
        />
      {/if}
    </PhaseHost>
  {:else}
    <JobUnselectedGuidance {onOpenJobManagement} />
  {/if}
</section>

<style>
  .job-run-page {
    display: grid;
    gap: 1.25rem;
    min-width: 0;
    padding-bottom: 7.5rem;
  }

  @media (max-width: 900px) {
    .job-run-page {
      padding-bottom: 16rem;
    }
  }
</style>
