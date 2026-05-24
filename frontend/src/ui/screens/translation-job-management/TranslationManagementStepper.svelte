<script lang="ts">
  import type {
    TranslationManagementViewContract,
    TranslationManagementViewId
  } from "@ui/stores/shell-state"

  interface Props {
    currentViewId: TranslationManagementViewId
    views: TranslationManagementViewContract[]
    onSelect: (viewId: TranslationManagementViewId) => void
  }

  let { currentViewId, views, onSelect }: Props = $props()

  type StepState = "complete" | "current" | "upcoming"

  function getStepState(viewId: TranslationManagementViewId): StepState {
    const currentIndex = views.findIndex((view) => view.id === currentViewId)
    const targetIndex = views.findIndex((view) => view.id === viewId)

    if (targetIndex < 0 || currentIndex < 0) {
      return "upcoming"
    }

    if (targetIndex < currentIndex) {
      return "complete"
    }

    if (targetIndex === currentIndex) {
      return "current"
    }

    return "upcoming"
  }

  function getStepStateLabel(
    stepState: StepState,
    directNavigation: boolean
  ): string {
    if (stepState === "current") {
      return "現在の作業"
    }

    return directNavigation ? "ここから開けます" : "順番に進む作業"
  }

  function getStepStateShortLabel(stepState: StepState): string {
    if (stepState === "complete") {
      return "完了"
    }

    if (stepState === "current") {
      return "現在"
    }

    return "次の作業"
  }

  function getTabAriaLabel(
    view: TranslationManagementViewContract,
    stepState: StepState
  ): string {
    const labels = [
      getStepStateLabel(stepState, view.directNavigation),
      view.label
    ]

    labels.push(view.description)
    return labels.join(" ")
  }
</script>

<section class="translation-stepper" aria-label="翻訳管理の全体進捗">
  <ol class="stepper-list" aria-label="翻訳管理の進行状況">
    {#each views as view (view.id)}
      {@const stepState = getStepState(view.id)}
      <li
        class="stepper-item"
        class:is-complete={stepState === "complete"}
        class:is-current={stepState === "current"}
        class:is-reference={!view.directNavigation}
      >
        {#if view.directNavigation}
          <button
            aria-current={view.id === currentViewId ? "step" : undefined}
            aria-label={getTabAriaLabel(view, stepState)}
            class="stepper-card"
            data-step-state={stepState}
            data-testid="translation-management-phase-card"
            onclick={() => onSelect(view.id)}
            role="tab"
            type="button"
          >
            <span class="step-marker" aria-hidden="true">
              {view.stepNumber}
            </span>
            <span class="step-body">
              <span class="step-meta">
                {getStepStateShortLabel(stepState)}
              </span>
              <strong>{view.label}</strong>
            </span>
          </button>
        {:else}
          <div
            aria-current={view.id === currentViewId ? "step" : undefined}
            class="stepper-card"
            data-step-state={stepState}
            data-testid="translation-management-phase-card"
          >
            <span class="step-marker" aria-hidden="true">
              {view.stepNumber}
            </span>
            <span class="step-body">
              <span class="step-meta">
                {getStepStateShortLabel(stepState)}
              </span>
              <strong>{view.label}</strong>
            </span>
          </div>
        {/if}
      </li>
    {/each}
  </ol>
</section>

<style>
  .translation-stepper {
    padding: 12px 14px;
    display: grid;
    color: var(--text);
  }

  .stepper-list {
    list-style: none;
    margin: 0;
    padding: 0;
    display: grid;
    grid-template-columns: repeat(auto-fit, minmax(128px, 1fr));
    gap: 8px;
  }

  .stepper-item {
    min-width: 0;
    position: relative;
  }

  .stepper-card {
    width: 100%;
    min-height: 58px;
    position: relative;
    display: grid;
    grid-template-columns: auto minmax(0, 1fr);
    gap: 10px;
    align-items: center;
    padding: 10px 12px;
    border-radius: 14px;
    border: 0.5px solid rgba(255, 186, 56, 0.12);
    background: rgba(17, 13, 12, 0.34);
    color: var(--text);
    text-align: left;
    transition:
      border-color var(--transition),
      background var(--transition),
      transform var(--transition),
      box-shadow var(--transition);
  }

  button.stepper-card {
    cursor: pointer;
    font: inherit;
  }

  button.stepper-card:hover,
  button.stepper-card:focus-visible {
    border-color: rgba(255, 186, 56, 0.24);
    background: rgba(255, 186, 56, 0.08);
    transform: translateY(-1px);
    box-shadow: 0 16px 40px rgba(0, 0, 0, 0.22);
    outline: none;
  }

  .step-marker {
    display: inline-flex;
    align-items: center;
    justify-content: center;
    width: 30px;
    height: 30px;
    border-radius: 999px;
    border: 1px solid rgba(255, 186, 56, 0.18);
    background: rgba(255, 255, 255, 0.03);
    color: var(--muted);
    font-size: 13px;
    font-weight: 700;
    flex-shrink: 0;
  }

  .step-body {
    min-width: 0;
    display: grid;
    gap: 0.12rem;
  }

  .step-meta {
    font-size: 11px;
    color: var(--muted);
  }

  .step-body strong {
    font-size: 13px;
    overflow-wrap: anywhere;
    line-height: 1.35;
  }

  .stepper-item.is-complete .step-marker {
    border-color: rgba(255, 186, 56, 0.28);
    background: rgba(255, 186, 56, 0.12);
    color: var(--primary);
  }

  .stepper-item.is-reference .stepper-card {
    cursor: default;
  }

  .stepper-item.is-current .stepper-card {
    border-color: rgba(255, 186, 56, 0.26);
    background:
      linear-gradient(
        180deg,
        rgba(255, 186, 56, 0.14),
        rgba(255, 186, 56, 0.06)
      ),
      rgba(17, 13, 12, 0.46);
    box-shadow: inset 0 0 0 1px rgba(255, 186, 56, 0.12);
  }

  .stepper-item.is-current .step-marker {
    border-color: rgba(255, 186, 56, 0.42);
    background: linear-gradient(135deg, var(--primary) 0%, #e8a31d 100%);
    color: var(--bg-strong);
    box-shadow: 0 0 0 4px rgba(255, 186, 56, 0.1);
  }

  .stepper-item.is-current .step-meta {
    color: var(--primary);
  }

  @media (max-width: 720px) {
    .stepper-list {
      grid-template-columns: repeat(2, minmax(0, 1fr));
    }

    .stepper-card {
      grid-template-columns: auto 1fr;
    }
  }
</style>
