<script lang="ts">
  import type { PhaseActionItem, PhaseDetailItem } from "./phase-panel-types"

  interface Props<ActionId extends string = string> {
    headingId: string
    testId: string
    eyebrow: string
    title: string
    progressLabel: string
    progressPercent: number
    progressDetail: string
    details: PhaseDetailItem[]
    currentPhaseLabel?: string
    actionAriaLabel?: string
    actions?: PhaseActionItem<ActionId>[]
    onAction?: (actionId: ActionId) => void | Promise<void>
  }

  let {
    headingId,
    testId,
    eyebrow,
    title,
    progressLabel,
    progressPercent,
    progressDetail,
    details,
    currentPhaseLabel = "",
    actionAriaLabel = "",
    actions = [],
    onAction
  }: Props = $props()
</script>

<section class="phase-card" aria-labelledby={headingId} data-testid={testId}>
  <div class="section-head">
    <div>
      <p class="eyebrow">{eyebrow}</p>
      <h3 id={headingId}>{title}</h3>
    </div>
    <span class="mini-text">{progressLabel}</span>
  </div>
  <div
    aria-label="progress"
    aria-valuemax="100"
    aria-valuemin="0"
    aria-valuenow={progressPercent}
    class="progress-bar"
    role="progressbar"
  >
    <span style={`width: ${progressPercent}%`}></span>
  </div>
  <p class="progress-copy">{progressDetail}</p>
  <div class:progress-body={actions.length > 0 && onAction}>
    <dl class="detail-grid compact">
      {#each details as detail (detail.label)}
        <div>
          <dt>{detail.label}</dt>
          <dd>{detail.value}</dd>
        </div>
      {/each}
    </dl>
    {#if actions.length > 0 && onAction}
      <div class="embedded-actions" aria-label={actionAriaLabel}>
        <div class="embedded-actions-head">
          <span class="mini-text">{currentPhaseLabel}</span>
        </div>
        <div class="action-grid">
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
      </div>
    {/if}
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

  .progress-copy,
  dd {
    color: rgba(250, 242, 232, 0.9);
    margin: 0;
  }

  dt {
    color: rgba(236, 223, 205, 0.8);
    margin: 0 0 0.18rem;
  }

  .progress-bar {
    background: rgba(255, 255, 255, 0.09);
    border-radius: 999px;
    height: 0.8rem;
    overflow: hidden;
  }

  .progress-bar span {
    background: linear-gradient(90deg, #d8a95f 0%, #ffcf8b 100%);
    border-radius: inherit;
    display: block;
    height: 100%;
  }

  .progress-body {
    align-items: start;
    display: grid;
    gap: 1.25rem;
    grid-template-columns: minmax(10rem, 0.42fr) minmax(16rem, 0.58fr);
  }

  .detail-grid {
    display: grid;
    gap: 0.9rem 1rem;
    grid-template-columns: repeat(2, minmax(0, 1fr));
    margin: 0;
  }

  .detail-grid div {
    display: grid;
    gap: 0.25rem;
    min-width: 0;
  }

  .progress-body .detail-grid {
    grid-template-columns: 1fr;
  }

  .embedded-actions {
    border-left: 1px solid rgba(226, 205, 173, 0.12);
    display: grid;
    gap: 0.75rem;
    padding-left: 1.25rem;
  }

  .embedded-actions-head {
    align-items: flex-start;
    display: flex;
    gap: 0.8rem;
    justify-content: space-between;
  }

  .action-grid {
    display: flex;
    flex-wrap: wrap;
    gap: 0.75rem;
  }

  .action-button {
    background: rgba(255, 255, 255, 0.04);
    border: 1px solid rgba(233, 213, 186, 0.18);
    border-radius: 14px;
    color: #fff6ea;
    cursor: pointer;
    font: inherit;
    min-height: 2.4rem;
    min-width: 5.5rem;
    padding: 0.45rem 0.8rem;
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
    .progress-body {
      grid-template-columns: 1fr;
    }

    .detail-grid {
      grid-template-columns: 1fr;
    }

    .embedded-actions {
      border-left: 0;
      border-top: 1px solid rgba(226, 205, 173, 0.12);
      padding-left: 0;
      padding-top: 0.9rem;
    }

    .section-head,
    .embedded-actions-head {
      flex-direction: column;
    }
  }
</style>
