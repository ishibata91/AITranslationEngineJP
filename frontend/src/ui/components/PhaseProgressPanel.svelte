<script lang="ts">
  import type { PhaseActionItem, PhaseDetailItem } from "./phase-panel-types"

  interface Props<ActionId extends string = string> {
    headingId: string
    testId: string
    eyebrow: string
    title: string
    progressBarTestId?: string
    progressCountsTestId?: string
    startButtonTestId?: string
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
    progressBarTestId = undefined,
    progressCountsTestId = undefined,
    startButtonTestId = undefined,
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
    <span class="progress-count" data-testid={progressCountsTestId}>
      {progressLabel}
    </span>
  </div>
  <p class="support-copy">{progressDetail}</p>
  <div
    aria-label="progress"
    aria-valuemax="100"
    aria-valuemin="0"
    aria-valuenow={progressPercent}
    class="progress-bar"
    data-testid={progressBarTestId}
    role="progressbar"
  >
    <span style={`width: ${progressPercent}%`}></span>
  </div>
  <div class="panel-body">
    <dl class="detail-grid">
      {#each details as detail (detail.label)}
        <div>
          <dt>{detail.label}</dt>
          <dd>{detail.value}</dd>
        </div>
      {/each}
    </dl>
    {#if actions.length > 0 && onAction}
      <div class="action-region" aria-label={actionAriaLabel}>
        {#if currentPhaseLabel}
          <span class="phase-label">{currentPhaseLabel}</span>
        {/if}
        <div class="button-row">
          {#each actions as action (action.id)}
            <button
              class="button-secondary"
              class:primary={action.tone === "primary"}
              class:warning={action.tone === "warning"}
              data-testid={action.id === "start"
                ? startButtonTestId
                : undefined}
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
    align-content: start;
    background: rgba(33, 27, 24, 0.88);
    border: 1px solid rgba(226, 205, 173, 0.14);
    border-radius: 20px;
    box-shadow: 0 18px 40px rgba(0, 0, 0, 0.22);
    display: grid;
    gap: 0.95rem;
    grid-auto-rows: max-content;
    padding: 1.4rem;
  }

  .section-head {
    align-items: flex-start;
    display: flex;
    gap: 0.8rem;
    justify-content: space-between;
  }

  .eyebrow,
  .phase-label {
    color: rgba(236, 223, 205, 0.72);
    font-size: 12px;
    letter-spacing: 0.08em;
    margin: 0;
  }

  h3 {
    color: #fff6ea;
    margin: 0;
  }

  .support-copy,
  dd {
    color: rgba(250, 242, 232, 0.9);
    margin: 0;
  }

  .support-copy {
    line-height: 1.7;
  }

  dt {
    color: rgba(236, 223, 205, 0.8);
    font-size: 12px;
    letter-spacing: 0.08em;
    margin: 0;
  }

  .progress-count {
    background: rgba(198, 155, 82, 0.16);
    border: 1px solid rgba(226, 205, 173, 0.16);
    border-radius: 999px;
    color: #ffe2ae;
    display: inline-flex;
    padding: 0.35rem 0.75rem;
    white-space: nowrap;
  }

  .progress-bar {
    background: rgba(255, 255, 255, 0.09);
    border: 0.5px solid rgba(255, 255, 255, 0.08);
    border-radius: 999px;
    height: 12px;
    overflow: hidden;
  }

  .progress-bar span {
    background: linear-gradient(90deg, var(--primary) 0%, #f5ca72 100%);
    border-radius: inherit;
    display: block;
    height: 100%;
  }

  .panel-body {
    display: grid;
    gap: 1rem;
  }

  .detail-grid {
    display: grid;
    gap: 10px;
    grid-template-columns: repeat(auto-fit, minmax(140px, 1fr));
    margin: 0;
  }

  .detail-grid div {
    background: rgba(255, 255, 255, 0.03);
    border: 0.5px solid rgba(255, 255, 255, 0.08);
    border-radius: 14px;
    display: grid;
    gap: 6px;
    min-width: 0;
    padding: 12px;
  }

  .action-region {
    display: grid;
    gap: 0.75rem;
  }

  .button-row {
    align-items: center;
    display: flex;
    flex-wrap: wrap;
    gap: 0.75rem;
  }

  .button-secondary {
    background: rgba(255, 255, 255, 0.04);
    border: 0.5px solid var(--line);
    border-radius: 999px;
    color: #fff6ea;
    cursor: pointer;
    font: inherit;
    min-height: 40px;
    min-width: 0;
    overflow-wrap: anywhere;
    padding: 0 16px;
  }

  .button-secondary.primary {
    background: linear-gradient(135deg, var(--primary) 0%, #f0a51f 100%);
    border: 0.5px solid transparent;
    color: #3f2400;
  }

  .button-secondary.warning {
    background: rgba(185, 105, 79, 0.18);
    border-color: rgba(255, 156, 124, 0.32);
    color: #ffd0c7;
  }

  .button-secondary:disabled {
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
    .detail-grid {
      grid-template-columns: 1fr;
    }

    .section-head {
      flex-direction: column;
    }

    .button-row > * {
      width: 100%;
    }
  }
</style>
