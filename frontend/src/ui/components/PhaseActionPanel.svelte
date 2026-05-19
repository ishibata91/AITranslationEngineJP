<script lang="ts" generics="ActionId extends string">
  import type { PhaseActionItem } from "./phase-panel-types"

  interface Props<ActionId extends string> {
    headingId: string
    testId: string
    currentPhaseLabel: string
    actions: PhaseActionItem<ActionId>[]
    columns?: 3 | 4
    onAction: (actionId: ActionId) => void | Promise<void>
  }

  let {
    headingId,
    testId,
    currentPhaseLabel,
    actions,
    columns = 3,
    onAction
  }: Props<ActionId> = $props()
</script>

<section
  class="phase-card action-card"
  aria-labelledby={headingId}
  data-testid={testId}
>
  <div class="section-head">
    <div>
      <p class="eyebrow">翻訳段階の操作</p>
      <h3 id={headingId}>操作</h3>
    </div>
    <span class="mini-text">{currentPhaseLabel}</span>
  </div>
  <div class="action-grid" data-columns={columns}>
    {#each actions as action (action.id)}
      <button
        class="action-button"
        class:primary={action.tone === "primary"}
        class:warning={action.tone === "warning"}
        disabled={action.disabled}
        onclick={() => onAction(action.id)}
        type="button"
      >
        {action.label}
      </button>
    {/each}
  </div>
  <div class="action-hints">
    {#each actions as action (action.id)}
      {#if action.disabled && action.blockedReason}
        <p>{action.label}: {action.blockedReason}</p>
      {/if}
    {/each}
  </div>
</section>

<style>
  .phase-card {
    background: rgba(33, 27, 24, 0.88);
    border: 1px solid rgba(226, 205, 173, 0.14);
    border-radius: 20px;
    box-shadow: 0 18px 40px rgba(0, 0, 0, 0.22);
    display: grid;
    gap: 1rem;
    padding: 1.4rem;
  }

  .section-head {
    align-items: flex-start;
    display: flex;
    gap: 0.8rem;
    justify-content: space-between;
  }

  .eyebrow,
  .mini-text {
    color: rgba(236, 223, 205, 0.72);
    font-size: 0.82rem;
    margin: 0 0 0.25rem;
  }

  h3 {
    color: #fff6ea;
    margin: 0;
  }

  .action-grid {
    display: grid;
    gap: 0.8rem;
    grid-template-columns: repeat(3, minmax(0, 1fr));
  }

  .action-grid[data-columns="4"] {
    grid-template-columns: repeat(4, minmax(0, 1fr));
  }

  .action-button {
    background: rgba(255, 255, 255, 0.04);
    border: 1px solid rgba(233, 213, 186, 0.18);
    border-radius: 14px;
    color: #fff6ea;
    cursor: pointer;
    font: inherit;
    min-height: 2.9rem;
    padding: 0.7rem 0.95rem;
  }

  .action-button.primary {
    background: linear-gradient(135deg, #cc8a39 0%, #f0b464 100%);
    border-color: transparent;
    color: #1b120c;
  }

  .action-button.warning {
    background: rgba(185, 105, 79, 0.18);
    color: #ffd0c7;
  }

  .action-button:disabled {
    cursor: not-allowed;
    opacity: 0.5;
  }

  .action-hints {
    color: rgba(250, 242, 232, 0.9);
    display: grid;
    font-size: 0.9rem;
    gap: 0.35rem;
  }

  .action-hints p {
    margin: 0;
  }

  @media (max-width: 900px) {
    .action-grid,
    .action-grid[data-columns="4"] {
      grid-template-columns: 1fr;
    }

    .section-head {
      flex-direction: column;
    }
  }
</style>
