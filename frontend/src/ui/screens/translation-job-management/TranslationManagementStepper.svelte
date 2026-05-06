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

  function getStepState(
    viewId: TranslationManagementViewId
  ): "complete" | "current" | "upcoming" {
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
    stepState: "complete" | "current" | "upcoming"
  ): string {
    switch (stepState) {
      case "complete":
        return "完了済み"
      case "current":
        return "現在位置"
      default:
        return "次の工程"
    }
  }

  function getLegacyLabel(viewId: TranslationManagementViewId): string | null {
    switch (viewId) {
      case "job-management":
        return "Job Management"
      case "job-setup":
        return "Job Setup"
      case "job-run":
        return "Job Run"
      default:
        return null
    }
  }

  function getTabAriaLabel(
    view: TranslationManagementViewContract,
    stepState: "complete" | "current" | "upcoming"
  ): string {
    const legacyLabel = getLegacyLabel(view.id)
    const labels = [getStepStateLabel(stepState), view.label]

    if (legacyLabel) {
      labels.push(legacyLabel)
    }

    labels.push(view.description)
    return labels.join(" ")
  }
</script>

<section class="translation-stepper">
  <div class="section-head">
    <div>
      <p class="page-label">翻訳セクション</p>
      <h2>進行ステップ</h2>
    </div>
  </div>

  <ol class="stepper-list" role="tablist" aria-label="翻訳管理セクション">
    {#each views as view, index (view.id)}
      {@const stepState = getStepState(view.id)}
      <li
        class="stepper-item"
        class:is-complete={stepState === "complete"}
        class:is-current={stepState === "current"}
      >
        <button
          aria-label={getTabAriaLabel(view, stepState)}
          aria-selected={view.id === currentViewId ? "true" : "false"}
          class="stepper-button"
          data-step-state={stepState}
          onclick={() => onSelect(view.id)}
          role="tab"
          type="button"
        >
          <span class="step-rail" aria-hidden="true">
            <span class="step-marker">{view.stepNumber}</span>
            {#if index < views.length - 1}
              <span class="step-line"></span>
            {/if}
          </span>
          <span class="step-body">
            <span class="step-meta">{getStepStateLabel(stepState)}</span>
            <strong>{view.label}</strong>
            <small>{view.description}</small>
          </span>
        </button>
      </li>
    {/each}
  </ol>
</section>

<style>
  .translation-stepper {
    padding: 24px;
    display: grid;
    gap: 1rem;
  }

  .section-head {
    display: flex;
    justify-content: space-between;
    align-items: flex-start;
    gap: 16px;
    flex-wrap: wrap;
    margin-bottom: 8px;
  }

  .page-label {
    margin: 0;
    font-size: 12px;
    letter-spacing: 0.12em;
    text-transform: uppercase;
    color: var(--muted);
  }

  .section-head h2 {
    margin: 0;
    font-size: 24px;
  }

  .stepper-list {
    list-style: none;
    margin: 0;
    padding: 0;
    display: grid;
    grid-template-columns: repeat(4, minmax(0, 1fr));
    gap: 14px;
  }

  .stepper-item {
    min-width: 0;
  }

  .stepper-button {
    width: 100%;
    min-height: 100%;
    position: relative;
    display: grid;
    grid-template-columns: auto 1fr;
    gap: 14px;
    align-items: start;
    padding: 18px;
    border-radius: var(--radius-md);
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

  .stepper-button:hover,
  .stepper-button:focus-visible {
    border-color: rgba(255, 186, 56, 0.24);
    background: rgba(255, 186, 56, 0.08);
    transform: translateY(-1px);
    box-shadow: 0 16px 40px rgba(0, 0, 0, 0.22);
    outline: none;
  }

  .step-rail {
    display: grid;
    grid-template-rows: auto 1fr;
    justify-items: center;
    gap: 8px;
    min-height: 100%;
  }

  .step-marker {
    display: inline-flex;
    align-items: center;
    justify-content: center;
    width: 36px;
    height: 36px;
    border-radius: 999px;
    border: 1px solid rgba(255, 186, 56, 0.18);
    background: rgba(255, 255, 255, 0.03);
    color: var(--muted);
    font-size: 14px;
    font-weight: 700;
    flex-shrink: 0;
  }

  .step-line {
    width: 2px;
    min-height: 52px;
    border-radius: 999px;
    background: rgba(255, 186, 56, 0.14);
  }

  .step-body {
    min-width: 0;
    display: grid;
    gap: 0.35rem;
  }

  .step-meta {
    font-size: 11px;
    letter-spacing: 0.1em;
    text-transform: uppercase;
    color: var(--muted);
  }

  .step-body strong {
    font-size: 16px;
  }

  .step-body small {
    color: var(--muted);
    overflow-wrap: anywhere;
    line-height: 1.5;
  }

  .stepper-item.is-complete .step-marker {
    border-color: rgba(255, 186, 56, 0.28);
    background: rgba(255, 186, 56, 0.12);
    color: var(--primary);
  }

  .stepper-item.is-complete .step-line {
    background: rgba(255, 186, 56, 0.34);
  }

  .stepper-item.is-current .stepper-button {
    border-color: rgba(255, 186, 56, 0.26);
    background:
      linear-gradient(180deg, rgba(255, 186, 56, 0.14), rgba(255, 186, 56, 0.06)),
      rgba(17, 13, 12, 0.46);
    box-shadow: inset 0 0 0 1px rgba(255, 186, 56, 0.12);
  }

  .stepper-item.is-current .step-marker {
    border-color: rgba(255, 186, 56, 0.42);
    background: linear-gradient(135deg, var(--primary) 0%, #e8a31d 100%);
    color: var(--bg-strong);
    box-shadow: 0 0 0 6px rgba(255, 186, 56, 0.1);
  }

  .stepper-item.is-current .step-meta {
    color: var(--primary);
  }

  @media (max-width: 1120px) {
    .stepper-list {
      grid-template-columns: repeat(2, minmax(0, 1fr));
    }

    .step-line {
      min-height: 36px;
    }
  }

  @media (max-width: 720px) {
    .stepper-list {
      grid-template-columns: 1fr;
    }

    .stepper-button {
      grid-template-columns: auto 1fr;
    }

    .step-line {
      min-height: 28px;
    }
  }
</style>
