<script lang="ts">
  import type { TranslationJobManagementJobRunTarget } from "@application/contract/translation-job-management/translation-job-management-screen-types"
  import type { JobRunPhaseStepId } from "@ui/screens/job-run/job-run-shell-props"

  interface JobRunPhaseStep {
    id: JobRunPhaseStepId
    label: string
    stepNumber: number
  }

  interface Props {
    target: TranslationJobManagementJobRunTarget
    currentPhasePage?: JobRunPhaseStepId
  }

  let { target, currentPhasePage }: Props = $props()

  const phaseSteps: JobRunPhaseStep[] = [
    { id: "term", label: "単語", stepNumber: 1 },
    { id: "persona", label: "NPC", stepNumber: 2 },
    { id: "body", label: "本文", stepNumber: 3 },
    { id: "complete", label: "確認", stepNumber: 4 }
  ]

  function resolveTargetPhasePage(
    targetPhase: TranslationJobManagementJobRunTarget["currentPhase"]
  ): JobRunPhaseStepId {
    if (targetPhase === "persona_generation") {
      return "persona"
    }

    if (targetPhase === "body_translation") {
      return "body"
    }

    return "term"
  }

  function getPhaseStepState(
    stepId: JobRunPhaseStepId,
    activePhasePage: JobRunPhaseStepId
  ): "past" | "current" | "upcoming" {
    const activeIndex = phaseSteps.findIndex(
      (step) => step.id === activePhasePage
    )
    const stepIndex = phaseSteps.findIndex((step) => step.id === stepId)

    if (stepIndex < activeIndex) {
      return "past"
    }

    if (stepIndex === activeIndex) {
      return "current"
    }

    return "upcoming"
  }

  const activePhasePage = $derived(
    currentPhasePage ?? resolveTargetPhasePage(target.currentPhase)
  )
</script>

<section
  class="job-run-target-summary"
  data-testid="job-run-selected-job-summary"
>
  <h3>ジョブ #{target.jobId}</h3>
  <ol class="job-phase-rail" aria-label="ジョブ内フェーズ">
    {#each phaseSteps as phaseStep (phaseStep.id)}
      {@const stepState = getPhaseStepState(phaseStep.id, activePhasePage)}
      <li
        class="phase-step"
        class:is-past={stepState === "past"}
        class:is-current={stepState === "current"}
      >
        <span
          class="phase-dot"
          aria-current={stepState === "current" ? "step" : undefined}
        >
          {phaseStep.stepNumber}
        </span>
        <span class="phase-name">{phaseStep.label}</span>
      </li>
    {/each}
  </ol>
</section>

<style>
  .job-run-target-summary {
    align-items: center;
    background: rgba(33, 27, 24, 0.88);
    border: 1px solid rgba(226, 205, 173, 0.14);
    border-radius: 16px;
    display: grid;
    grid-template-columns: max-content minmax(0, 1fr);
    gap: 0.9rem;
    min-width: 0;
    padding: 0.72rem 1rem;
  }

  h3 {
    color: #fff6ea;
    font-size: 1.05rem;
    margin: 0;
    overflow-wrap: anywhere;
  }

  .job-phase-rail {
    list-style: none;
    margin: 0;
    padding: 0;
    display: grid;
    grid-template-columns: repeat(4, minmax(0, 1fr));
    min-width: 0;
  }

  .phase-step {
    position: relative;
    display: grid;
    justify-items: center;
    gap: 0.25rem;
    min-width: 0;
    color: rgba(226, 205, 173, 0.56);
  }

  .phase-step:not(:last-child)::after {
    content: "";
    position: absolute;
    top: 13px;
    left: calc(50% + 18px);
    right: calc(-50% + 18px);
    height: 2px;
    border-radius: 999px;
    background: rgba(226, 205, 173, 0.16);
  }

  .phase-step.is-past:not(:last-child)::after {
    background: rgba(255, 186, 56, 0.52);
  }

  .phase-dot {
    position: relative;
    z-index: 1;
    display: inline-flex;
    align-items: center;
    justify-content: center;
    width: 26px;
    height: 26px;
    border-radius: 999px;
    border: 1px solid rgba(226, 205, 173, 0.18);
    background: #18110e;
    color: rgba(226, 205, 173, 0.72);
    font-size: 0.76rem;
    font-weight: 700;
  }

  .phase-name {
    max-width: 100%;
    overflow: hidden;
    color: inherit;
    font-size: 0.72rem;
    line-height: 1.2;
    text-overflow: ellipsis;
    white-space: nowrap;
  }

  .phase-step.is-past,
  .phase-step.is-current {
    color: #fff6ea;
  }

  .phase-step.is-past .phase-dot {
    border-color: rgba(255, 186, 56, 0.4);
    background: rgba(255, 186, 56, 0.16);
    color: #ffcf7a;
  }

  .phase-step.is-current .phase-dot {
    border-color: rgba(255, 186, 56, 0.9);
    background: #ffba38;
    color: #1a110b;
    box-shadow: 0 0 0 4px rgba(255, 186, 56, 0.13);
  }

  @media (max-width: 720px) {
    .job-run-target-summary {
      grid-template-columns: 1fr;
      gap: 0.7rem;
    }
  }
</style>
