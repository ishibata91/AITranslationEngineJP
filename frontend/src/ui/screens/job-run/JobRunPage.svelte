<script lang="ts">
  import { onMount } from "svelte"

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
  import PersonaGenerationPhasePanel from "@ui/screens/persona-generation-phase/PersonaGenerationPhasePanel.svelte"
  import TermTranslationPhasePanel from "@ui/screens/term-translation-phase/TermTranslationPhasePanel.svelte"

  interface Props {
    createBodyController: CreateBodyTranslationPhaseScreenController | null
    createPersonaController: CreatePersonaGenerationPhaseScreenController | null
    createController: CreateTermTranslationPhaseScreenController | null
  }

  let { createBodyController, createPersonaController, createController }: Props = $props()

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
  let jobIdInput = $state("")

  const unsubscribe = controller.subscribe((nextViewModel) => {
    viewModel = nextViewModel
    jobIdInput = nextViewModel.jobId?.toString() ?? jobIdInput
  })

  const unsubscribePersona = personaController.subscribe((nextViewModel) => {
    personaViewModel = nextViewModel
    jobIdInput = nextViewModel.jobId?.toString() ?? jobIdInput
  })

  const unsubscribeBody = bodyController
    ? bodyController.subscribe((nextViewModel) => {
        bodyViewModel = nextViewModel
        jobIdInput = nextViewModel.jobId?.toString() ?? jobIdInput
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

  async function loadJobSummary(): Promise<void> {
    const nextJobId = Number(jobIdInput.trim())
    if (!Number.isInteger(nextJobId) || nextJobId <= 0) {
      await Promise.all([
        controller.setJobId(null),
        personaController.setJobId(null),
        ...(bodyController ? [bodyController.setJobId(null)] : [])
      ])
      return
    }

    await Promise.all([
      controller.setJobId(nextJobId),
      personaController.setJobId(nextJobId),
      ...(bodyController ? [bodyController.setJobId(nextJobId)] : [])
    ])
  }

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
</script>

<section class="job-run-page">
  <section class="job-run-selector">
    <div>
      <p class="eyebrow">job target</p>
      <h2>Job Run</h2>
      <p class="selector-copy">
        対象 job id を指定して summary を取得します。gateway 統合前は mock
        controller または未接続状態で確認します。
      </p>
    </div>
    <div class="selector-form">
      <label class="field-block" for="termPhaseJobIdInput">
        <span>job id</span>
        <input
          id="termPhaseJobIdInput"
          inputmode="numeric"
          onkeydown={(event) => {
            if (event.key === "Enter") {
              void loadJobSummary()
            }
          }}
          bind:value={jobIdInput}
          placeholder="例: 101"
          type="text"
        />
      </label>
      <button onclick={() => void loadJobSummary()} type="button"
        >summary 取得</button
      >
    </div>
  </section>

  <TermTranslationPhasePanel
    {viewModel}
    onAction={(actionId: TermTranslationPhaseActionCard["id"]) =>
      handleAction(actionId)}
  />

  <PersonaGenerationPhasePanel
    viewModel={personaViewModel}
    onAction={(actionId: PersonaGenerationPhaseActionKind) =>
      handlePersonaAction(actionId)}
  />

  {#if bodyViewModel}
    <BodyTranslationPhasePanel
      viewModel={bodyViewModel}
      onAction={(actionId: BodyTranslationPhaseActionKind) =>
        handleBodyAction(actionId)}
    />
  {/if}
</section>

<style>
  .job-run-page {
    display: grid;
    gap: 1.25rem;
  }

  .job-run-selector {
    display: flex;
    gap: 1rem;
    justify-content: space-between;
    align-items: end;
    padding: 1.25rem 1.4rem;
    border: 1px solid rgba(226, 205, 173, 0.14);
    border-radius: 20px;
    background: rgba(33, 27, 24, 0.88);
  }

  .eyebrow,
  .selector-copy,
  .field-block span {
    color: rgba(236, 223, 205, 0.78);
  }

  h2 {
    margin: 0.2rem 0 0;
    color: #fff6ea;
  }

  .selector-form {
    display: flex;
    gap: 0.8rem;
    align-items: end;
  }

  .field-block {
    display: grid;
    gap: 0.35rem;
    min-width: min(16rem, 100%);
  }

  input,
  button {
    min-height: 2.8rem;
    padding: 0.65rem 0.9rem;
    border-radius: 14px;
    border: 1px solid rgba(233, 213, 186, 0.18);
    font: inherit;
  }

  input {
    background: rgba(255, 255, 255, 0.05);
    color: #fff6ea;
  }

  button {
    background: linear-gradient(135deg, #cc8a39 0%, #f0b464 100%);
    color: #1b120c;
    cursor: pointer;
  }

  @media (max-width: 900px) {
    .job-run-selector,
    .selector-form {
      flex-direction: column;
      align-items: stretch;
    }

    .field-block {
      min-width: 0;
    }
  }
</style>
