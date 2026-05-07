<script lang="ts">
  import type { Snippet } from "svelte"

  interface Props {
    title: string
    titleId: string
    description: string
    reasons: string[]
    emptyText?: string
    primaryLabel: string
    primaryDisabled?: boolean
    onPrimary: () => void
    children?: Snippet
  }

  let {
    title,
    titleId,
    description,
    reasons,
    emptyText = "確認する項目はありません。",
    primaryLabel,
    primaryDisabled = false,
    onPrimary,
    children
  }: Props = $props()

  const firstReason = $derived(reasons[0] ?? "")
  const remainingReasons = $derived(reasons.slice(1))
  const remainingCount = $derived(remainingReasons.length)
  const tooltipId = $derived(`${titleId}-remaining-reasons`)
</script>

<section
  class="sticky-action-footer"
  aria-label={`${title}: ${description}`}
  aria-labelledby={titleId}
>
  <div class="footer-copy">
    <h3 id={titleId}>{title}</h3>
    {#if reasons.length === 0}
      <p class="reason-summary">{emptyText}</p>
    {:else}
      <div class="reason-summary" aria-live="polite">
        <span class="reason-primary">{firstReason}</span>
        {#if remainingCount > 0}
          <span class="reason-more-shell">
            <button
              aria-describedby={tooltipId}
              class="reason-more"
              type="button"
            >
              ほか {remainingCount} 件
            </button>
            <span class="reason-tooltip" id={tooltipId} role="tooltip">
              <span class="tooltip-title">残りの不足</span>
              <ul>
                {#each remainingReasons as reason (reason)}
                  <li>{reason}</li>
                {/each}
              </ul>
            </span>
          </span>
        {/if}
      </div>
    {/if}
    {@render children?.()}
  </div>
  <div class="footer-actions">
    <button
      class="button-primary"
      disabled={primaryDisabled}
      onclick={onPrimary}
      type="button"
    >
      {primaryLabel}
    </button>
  </div>
</section>

<style>
  .sticky-action-footer {
    align-items: end;
    background: rgba(27, 20, 17, 0.95);
    backdrop-filter: blur(12px);
    border: 1px solid rgba(255, 212, 165, 0.18);
    border-radius: 1.25rem;
    bottom: 1rem;
    box-shadow: 0 20px 40px rgba(6, 4, 3, 0.18);
    display: flex;
    gap: 1rem;
    justify-content: space-between;
    padding: 1rem 1.25rem;
    position: sticky;
    z-index: 2;
  }

  .footer-copy,
  .footer-actions {
    display: grid;
    gap: 0.65rem;
  }

  .footer-copy {
    min-width: 0;
  }

  .footer-actions {
    align-self: end;
    justify-items: end;
  }

  .reason-summary {
    color: rgba(252, 241, 232, 0.86);
  }

  .reason-summary {
    align-items: center;
    display: flex;
    flex-wrap: wrap;
    gap: 0.45rem;
    justify-content: flex-start;
    margin: 0;
    min-height: 1.9rem;
  }

  .reason-primary {
    max-width: min(44rem, 68vw);
    overflow: hidden;
    text-overflow: ellipsis;
    white-space: nowrap;
  }

  .reason-more-shell {
    display: inline-flex;
    position: relative;
  }

  .reason-more {
    background: rgba(255, 190, 126, 0.08);
    border: 1px dotted rgba(255, 212, 165, 0.42);
    border-radius: 999px;
    color: #ffd8ae;
    cursor: help;
    display: inline-flex;
    font-size: 0.82rem;
    line-height: 1;
    padding: 0.34rem 0.62rem;
    text-decoration: underline;
    text-decoration-style: dotted;
    text-underline-offset: 0.18rem;
  }

  .reason-more:focus {
    outline: 2px solid rgba(255, 204, 136, 0.72);
    outline-offset: 2px;
  }

  .reason-tooltip {
    background: rgba(22, 19, 18, 0.98);
    border: 1px solid rgba(255, 255, 255, 0.16);
    border-radius: 6px;
    bottom: calc(100% + 8px);
    color: #fef3e8;
    display: none;
    font-size: 0.8rem;
    line-height: 1.5;
    min-width: 18rem;
    padding: 0.75rem 0.85rem;
    position: absolute;
    right: 0;
    text-align: left;
    z-index: 30;
  }

  .reason-more-shell:hover .reason-tooltip,
  .reason-more-shell:focus-within .reason-tooltip {
    display: grid;
    gap: 0.45rem;
  }

  .tooltip-title {
    color: rgba(255, 215, 176, 0.72);
    font-size: 0.76rem;
  }

  .reason-tooltip ul {
    display: grid;
    gap: 0.35rem;
    margin: 0;
    padding-left: 1.1rem;
  }

  .reason-tooltip li {
    overflow-wrap: anywhere;
    word-break: break-word;
  }

  .button-primary {
    background: linear-gradient(135deg, #ff9f5a, #ffcc88);
    border: 1px solid rgba(255, 212, 165, 0.18);
    border-radius: 0.9rem;
    color: #24150d;
    cursor: pointer;
    padding: 0.8rem 1rem;
  }

  .button-primary:disabled {
    cursor: not-allowed;
    opacity: 0.56;
  }

  @media (max-width: 720px) {
    .sticky-action-footer {
      align-items: stretch;
      flex-direction: column;
    }

    .footer-actions {
      justify-items: stretch;
    }

    .reason-primary {
      max-width: 100%;
    }

    .reason-tooltip {
      left: 0;
      min-width: min(18rem, 82vw);
      right: auto;
    }
  }
</style>
