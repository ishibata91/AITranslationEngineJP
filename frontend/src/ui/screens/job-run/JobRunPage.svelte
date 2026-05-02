<script lang="ts">
  import { onMount } from "svelte"

  import type {
    CreateTermTranslationPhaseScreenController,
    TermTranslationPhaseActionCard,
    TermTranslationPhaseScreenControllerContract
  } from "@application/contract/term-translation-phase"
  import TermTranslationPhasePanel from "@ui/screens/term-translation-phase/TermTranslationPhasePanel.svelte"

  interface Props {
    createController: CreateTermTranslationPhaseScreenController | null
  }

  let { createController }: Props = $props()

  function resolveController(): TermTranslationPhaseScreenControllerContract {
    if (!createController) {
      throw new Error("term translation phase screen controller factory is not provided")
    }

    return createController()
  }

  const controller = resolveController()
  let viewModel = $state(controller.getViewModel())
  let jobIdInput = $state("")

  const unsubscribe = controller.subscribe((nextViewModel) => {
    viewModel = nextViewModel
    jobIdInput = nextViewModel.jobId?.toString() ?? jobIdInput
  })

  onMount(() => {
    void controller.mount()

    return () => {
      unsubscribe()
      controller.dispose()
    }
  })

  async function loadJobSummary(): Promise<void> {
    const nextJobId = Number(jobIdInput.trim())
    if (!Number.isInteger(nextJobId) || nextJobId <= 0) {
      await controller.setJobId(null)
      return
    }

    await controller.setJobId(nextJobId)
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
</script>

<section class="job-run-page">
  <section class="job-run-selector">
    <div>
      <p class="eyebrow">job target</p>
      <h2>Job Run</h2>
      <p class="selector-copy">
        対象 job id を指定して summary を取得します。gateway 統合前は mock controller または未接続状態で確認します。
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
      <button onclick={() => void loadJobSummary()} type="button">summary 取得</button>
    </div>
  </section>

  <TermTranslationPhasePanel
    {viewModel}
    onAction={(actionId: TermTranslationPhaseActionCard["id"]) =>
      handleAction(actionId)}
  />
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
