<script lang="ts">
  import { onMount } from "svelte"

  import type { TranslationJobManagementJobRunTarget } from "@application/contract/translation-job-management/translation-job-management-screen-types"
  import type { TranslationManagementViewId } from "@ui/stores/shell-state"
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
  import PhaseNavigationFooter from "@ui/screens/job-run/PhaseNavigationFooter.svelte"
  import TranslationCompletePage from "@ui/screens/job-run/TranslationCompletePage.svelte"
  import PersonaGenerationPhasePanel from "@ui/screens/persona-generation-phase/PersonaGenerationPhasePanel.svelte"
  import TermTranslationPhasePanel from "@ui/screens/term-translation-phase/TermTranslationPhasePanel.svelte"

  interface Props {
    createBodyController: CreateBodyTranslationPhaseScreenController | null
    createPersonaController: CreatePersonaGenerationPhaseScreenController | null
    createController: CreateTermTranslationPhaseScreenController | null
    selectedJobTarget?: TranslationJobManagementJobRunTarget | null
    onOpenJobManagement?: () => void
    onOpenOutputManagement?: () => void
    onPhaseViewChange?: (viewId: TranslationManagementViewId) => void
  }

  let {
    createBodyController,
    createPersonaController,
    createController,
    selectedJobTarget = null,
    onOpenJobManagement = () => undefined,
    onOpenOutputManagement = () => undefined,
    onPhaseViewChange = () => undefined
  }: Props = $props()

  type PhasePageId = "term" | "persona" | "body" | "complete"

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

  function toTranslationManagementViewId(
    phasePage: PhasePageId
  ): TranslationManagementViewId {
    switch (phasePage) {
      case "persona":
        return "persona-generation"
      case "body":
        return "body-translation"
      case "complete":
        return "translation-complete"
      default:
        return "term-translation"
    }
  }

  function setCurrentPhasePage(phasePage: PhasePageId): void {
    currentPhasePage = phasePage
    onPhaseViewChange(toTranslationManagementViewId(phasePage))
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
      ...(bodyController ? [bodyController.setJobId(selectedJobTarget.jobId)] : [])
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
      case "refresh":
        await controller.refresh()
        return
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
      case "refresh":
        await personaController.refresh()
        return
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
      case "refresh":
        await bodyController.refresh()
        return
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
</script>

<section class="job-run-page">
  {#if selectedJobTarget}
    <section class="job-run-target-summary">
      <div>
        <p class="eyebrow">選択中のジョブ</p>
        <h3>ジョブ #{selectedJobTarget.jobId}</h3>
      </div>
      <dl class="target-summary-grid">
        <div>
          <dt>状態</dt>
          <dd>{selectedJobTarget.stateLabel} / {selectedJobTarget.stateDescription}</dd>
        </div>
        <div>
          <dt>現在の翻訳段階</dt>
          <dd>{selectedJobTarget.currentPhaseLabel}</dd>
        </div>
        <div>
          <dt>進捗</dt>
          <dd>{selectedJobTarget.progressLabel}</dd>
        </div>
        <div>
          <dt>入力</dt>
          <dd>{selectedJobTarget.inputSourceLabel}</dd>
        </div>
      </dl>
      <details class="target-path">
        <summary>入力ファイル path</summary>
        <p>{selectedJobTarget.sourcePath}</p>
      </details>
    </section>

    {#if currentPhasePage === "term"}
      <TermTranslationPhasePanel
        {viewModel}
        onAction={(actionId: TermTranslationPhaseActionCard["id"]) =>
          handleAction(actionId)}
      />
      <PhaseNavigationFooter
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
      />
      <PhaseNavigationFooter
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
      />
      <PhaseNavigationFooter
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
      />
      <PhaseNavigationFooter
        description="翻訳結果を確認した後は、出力管理で出力対象を選びます。"
        onBack={onOpenJobManagement}
        onPrimary={onOpenOutputManagement}
        primaryLabel="出力管理へ進む"
        reasons={[]}
        title="翻訳完了後の次の作業"
        titleId="completeNavigationFooter"
      />
    {/if}
  {:else}
    <section class="job-run-target-summary">
      <div>
        <p class="eyebrow">ジョブ未選択</p>
        <h3>未完了ジョブ一覧でジョブを選んでください</h3>
      </div>
      <p class="selector-copy">
        翻訳段階の画面は選択済みジョブだけを対象にします。ジョブID
        の手入力はこのページでは扱いません。
      </p>
      <button class="secondary-button" onclick={onOpenJobManagement} type="button">
        一覧へ戻る
      </button>
    </section>
  {/if}
</section>

<style>
  .job-run-page {
    display: grid;
    gap: 1.25rem;
    min-width: 0;
    padding-bottom: 7.5rem;
  }

  .job-run-target-summary {
    display: grid;
    gap: 0.8rem;
    padding: 1.15rem 1.3rem;
    border: 1px solid rgba(226, 205, 173, 0.14);
    border-radius: 20px;
    background: rgba(33, 27, 24, 0.88);
    min-width: 0;
  }

  .eyebrow,
  .selector-copy,
  .target-path {
    color: rgba(236, 223, 205, 0.78);
  }

  h3 {
    margin: 0.2rem 0 0;
    color: #fff6ea;
  }

  .secondary-button {
    justify-self: start;
    min-height: 2.8rem;
    padding: 0.65rem 0.9rem;
    border-radius: 14px;
    border: 1px solid rgba(233, 213, 186, 0.18);
    background: rgba(255, 255, 255, 0.05);
    color: #fff6ea;
    cursor: pointer;
  }

  .target-summary-grid {
    display: grid;
    gap: 0.6rem;
    grid-template-columns: repeat(2, minmax(0, 1fr));
    margin: 0;
  }

  .target-summary-grid div {
    display: grid;
    gap: 0.2rem;
    min-width: 0;
  }

  .target-summary-grid dt {
    color: rgba(236, 223, 205, 0.78);
  }

  .target-summary-grid dd {
    margin: 0;
    color: #fff6ea;
    overflow-wrap: anywhere;
  }

  @media (max-width: 900px) {
    .job-run-page {
      padding-bottom: 16rem;
    }

    .target-summary-grid {
      grid-template-columns: 1fr;
    }
  }
</style>
